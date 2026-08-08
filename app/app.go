package app

import (
	"chrono-ntp/audio"
	"chrono-ntp/display"
	"chrono-ntp/timesource"
	"log"
	"time"

	"github.com/ebitengine/oto/v3"
)

const (
	timeSourceRefreshInterval = 15 * time.Minute
)

type InitOptions struct {
	Server      string
	TimeZone    string
	Offline     bool
	BeepPattern string
}

type RunOptions struct {
	AudioContext  *oto.Context
	Display       *display.Display
	TimeSource    timesource.TimeSource
	Location      *time.Location
	BeepStrategy  audio.BeepStrategy
	HideStatusBar bool
	HideDate      bool
	ShowTimeZone  bool
	DateFormat    string
	TimeFormat    string
	Offline       bool
}

func Init(options InitOptions) (*oto.Context, *display.Display, timesource.TimeSource, *time.Location, audio.BeepStrategy) {
	timeZoneLocation, err := time.LoadLocation(options.TimeZone)
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}

	audioContext, err := audio.InitializeAudioContext()
	if err != nil {
		log.Fatalf("Failed to initialize audio context: %v", err)
	}

	// Initialize display early to show loading message
	disp, err := display.NewDisplay()
	if err != nil {
		log.Fatalf("Failed to create display: %v", err)
	}
	if err := disp.Init(); err != nil {
		log.Fatalf("Failed to initialize display: %v", err)
	}

	var timeSource timesource.TimeSource
	if options.Offline {
		timeSource = timesource.NewOfflineTimeSource()
	} else {
		timeSource = timesource.NewNtpTimeSource(options.Server)
	}

	beepStrategy := audio.BeepStrategy(options.BeepPattern)

	return audioContext, disp, timeSource, timeZoneLocation, beepStrategy
}

func Run(options RunOptions) {
	defer options.Display.Finalize()

	go func() {
		options.TimeSource.Refresh()

		ticker := time.NewTicker(timeSourceRefreshInterval)
		defer ticker.Stop()
		for range ticker.C {
			options.TimeSource.Refresh()
		}
	}()

	quit := make(chan struct{})
	go options.Display.PollEvents(quit)

	displayTicker := time.NewTicker(100 * time.Millisecond)
	defer displayTicker.Stop()

	for {
		select {
		case <-displayTicker.C:
			timeStatus := options.TimeSource.TimeStatus()

			options.Display.Update(display.DisplayState{
				TimeStatus:    timeStatus,
				TimeZone:      options.Location,
				DateFormat:    options.DateFormat,
				TimeFormat:    options.TimeFormat,
				HideDate:      options.HideDate,
				ShowTimeZone:  options.ShowTimeZone,
				HideStatusBar: options.HideStatusBar,
				Offline:       options.Offline,
			})

			audio.BeepTick(options.AudioContext, timeStatus.ReferenceTime, options.BeepStrategy)
		case <-quit:
			return
		}
	}
}
