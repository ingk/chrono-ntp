package audio

import (
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
