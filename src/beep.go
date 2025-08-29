package main

import (
	"io"
	"math"
	"runtime"
	"time"

	"github.com/ebitengine/oto/v3"
)

const (
	sampleRate   = 44100
	channelCount = 1
	freq         = 1000.0
	amplitude    = 0.3
	maxInt16     = 32767
	shortMs      = 100
	longMs       = 1000
)

var (
	shortBeep []byte
	longBeep  []byte
)

func init() {
	shortBeep = makeSineWaveTable(shortMs)
	longBeep = makeSineWaveTable(longMs)
}

func makeSineWaveTable(durationMs int) []byte {
	numSamples := sampleRate * durationMs / 1000
	buf := make([]byte, numSamples*2) // 2 bytes per sample
	for i := range numSamples {
		t := float64(i) / float64(sampleRate)
		v := int16(math.Sin(2*math.Pi*freq*t) * amplitude * maxInt16)
		buf[2*i] = byte(v)
		buf[2*i+1] = byte(v >> 8)
	}
	return buf
}

type beepReader struct {
	data []byte
	pos  int
}

func (r *beepReader) Read(buf []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n := copy(buf, r.data[r.pos:])
	r.pos += n
	return n, nil
}

func playBeep(ctx *oto.Context, data []byte, durationMs int) {
	reader := &beepReader{data: data}
	player := ctx.NewPlayer(reader)
	player.Play()
	time.Sleep(time.Duration(durationMs) * time.Millisecond)
	runtime.KeepAlive(player)
}

func InitializeAudioContext() (*oto.Context, error) {
	op := &oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: channelCount,
		Format:       oto.FormatSignedInt16LE,
	}
	ctx, ready, err := oto.NewContext(op)
	if err != nil {
		return nil, err
	}
	<-ready
	return ctx, nil
}

func PlayShortBeep(ctx *oto.Context) {
	playBeep(ctx, shortBeep, shortMs)
}

func PlayLongBeep(ctx *oto.Context) {
	playBeep(ctx, longBeep, longMs)
}
