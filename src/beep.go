package main

import (
	"fmt"
	"io"
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/ebitengine/oto/v3"
)

type SineWave struct {
	freq   float64
	length int64
	pos    int64

	channelCount int
	format       oto.Format

	remaining []byte
}

var (
	sampleRate    = 44100
	channelCount  = 1
	sineTableSize = 1024
	sineTable     = make([]float64, sineTableSize)
)

func init() {
	for i := range sineTableSize {
		sineTable[i] = math.Sin(2 * math.Pi * float64(i) / float64(sineTableSize))
	}
}

func formatByteLength(format oto.Format) int {
	switch format {
	case oto.FormatFloat32LE:
		return 4
	case oto.FormatUnsignedInt8:
		return 1
	case oto.FormatSignedInt16LE:
		return 2
	default:
		panic(fmt.Sprintf("unexpected format: %d", format))
	}
}

func InitializeAudioContext() *oto.Context {
	op := &oto.NewContextOptions{}
	op.SampleRate = 44100
	op.ChannelCount = 1
	op.Format = oto.FormatSignedInt16LE
	c, ready, err := oto.NewContext(op)
	if err != nil {
		panic(err)
	}
	<-ready
	return c
}

func newSineWave(freq float64, duration time.Duration, channelCount int, format oto.Format) *SineWave {
	l := int64(channelCount) * int64(formatByteLength(format)) * 44100 * int64(duration) / int64(time.Second)
	l = l / 4 * 4
	return &SineWave{
		freq:         freq,
		length:       l,
		channelCount: channelCount,
		format:       format,
	}
}

func (s *SineWave) Read(buf []byte) (int, error) {
	if len(s.remaining) > 0 {
		n := copy(buf, s.remaining)
		copy(s.remaining, s.remaining[n:])
		s.remaining = s.remaining[:len(s.remaining)-n]
		return n, nil
	}

	if s.pos == s.length {
		return 0, io.EOF
	}

	eof := false
	if s.pos+int64(len(buf)) > s.length {
		buf = buf[:s.length-s.pos]
		eof = true
	}

	var origBuf []byte
	if len(buf)%4 > 0 {
		origBuf = buf
		buf = make([]byte, len(origBuf)+4-len(origBuf)%4)
	}

	length := float64(sampleRate) / float64(s.freq)
	num := formatByteLength(s.format) * s.channelCount
	p := s.pos / int64(num)

	for i := 0; i < len(buf)/num; i++ {
		idx := int((float64(p)/length)*float64(sineTableSize)) % sineTableSize
		const max = 32767
		b := int16(sineTable[idx] * 0.3 * max)
		buf[num*i] = byte(b)
		buf[num*i+1] = byte(b >> 8)
		p++
	}

	s.pos += int64(len(buf))

	n := len(buf)
	if origBuf != nil {
		n = copy(origBuf, buf)
		s.remaining = buf[n:]
	}

	if eof {
		return n, io.EOF
	}
	return n, nil
}

func playerPlay(context *oto.Context, freq float64, duration time.Duration) *oto.Player {
	p := context.NewPlayer(newSineWave(freq, duration, channelCount, oto.FormatSignedInt16LE))
	p.Play()
	return p
}

func PlayBeep(c *oto.Context, duration time.Duration) {
	var wg sync.WaitGroup
	var players []*oto.Player
	var m sync.Mutex

	wg.Add(1)
	go func() {
		defer wg.Done()
		player := playerPlay(c, 1000, duration)
		m.Lock()
		players = append(players, player)
		m.Unlock()
		time.Sleep(1 * time.Second)
	}()

	wg.Wait()

	// Pin the players not to GC the players.
	runtime.KeepAlive(players)
}
