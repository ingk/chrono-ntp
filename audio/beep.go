package audio

import (
	"io"
	"math"
	"runtime"
	"time"

	"github.com/ebitengine/oto/v3"
)

type BeepStrategy string

const (
	BeepStrategyGreenwich  BeepStrategy = "greenwich"
	BeepStrategyContinuous BeepStrategy = "continuous"
)

var AllowedBeepPatterns = [...]string{string(BeepStrategyGreenwich), string(BeepStrategyContinuous)}

var (
	shortBeep     []byte
	longBeep      []byte
	currentSecond int = -1
)

func init() {
	shortBeep = makeSineWaveTable(shortMs)
	longBeep = makeSineWaveTable(longMs)
}

func BeepTick(ctx *oto.Context, now time.Time, strategy BeepStrategy) {
	switch strategy {
	case BeepStrategyGreenwich:
		beepGreenwich(ctx, now)
	case BeepStrategyContinuous:
		beepContinuous(ctx, now)
	}
}

func beepGreenwich(ctx *oto.Context, now time.Time) {
	shouldBeep := now.Second() >= 55 || now.Second() == 0
	if !shouldBeep || currentSecond == now.Second() {
		return
	}

	currentSecond = now.Second()
	if currentSecond == 0 {
		go playBeep(ctx, longBeep, longMs)
	} else {
		go playBeep(ctx, shortBeep, shortMs)
	}
}

func beepContinuous(ctx *oto.Context, now time.Time) {
	if currentSecond == -1 {
		currentSecond = now.Second()
		return
	}
	if currentSecond == now.Second() {
		return
	}

	currentSecond = now.Second()
	go playBeep(ctx, shortBeep, shortMs)
}

func playBeep(ctx *oto.Context, data []byte, durationMs int) {
	reader := &beepReader{data: data}
	player := ctx.NewPlayer(reader)
	player.Play()
	time.Sleep(time.Duration(durationMs) * time.Millisecond)
	runtime.KeepAlive(player)
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
