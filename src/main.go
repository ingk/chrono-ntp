package main

import (
	"flag"
	"fmt"
	"log"

	"slices"
	"strings"
	"time"

	"github.com/beevik/ntp"
	"github.com/gdamore/tcell/v2"
)

const (
	appName                  = "chrono-ntp"
	appVersion               = "dev"
	ntpOffsetRefreshInterval = 15 * time.Minute
)

var allowedTimeFormats = [...]string{"ISO8601", "12h", "12h_AM_PM", ".beat", "septimal"}

func main() {
	config, err := LoadConfiguration()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ntpServer := flag.String("server", config.Server, "NTP server to synchronize time from")
	timeZone := flag.String("time-zone", config.TimeZone, "Time zone name (e.g., 'America/New_York')")
	debug := flag.Bool("debug", false, "Show debug information (e.g., offset from NTP server) and exit")
	hideStatusbar := flag.Bool("hide-statusbar", config.HideStatusbar, "Hide the status bar")
	hideDate := flag.Bool("hide-date", config.HideDate, "Hide the date display")
	showTimeZone := flag.Bool("show-time-zone", config.ShowTimeZone, "Show the time zone")
	timeFormat := flag.String("time-format", config.TimeFormat, fmt.Sprintf("Time display format (%s)", strings.Join(allowedTimeFormats[:], ", ")))
	beeps := flag.Bool("beeps", config.Beeps, "Play 6 beeps at the end of each minute, with the sixth beep at second 0 (emulates the Greenwich Time Signal)")
	version := flag.Bool("version", false, "Show version and exit")
	offline := flag.Bool("offline", false, "Run in offline mode (use system time, ignore NTP server)")
	flag.Parse()

	if *debug {
		fmt.Printf("Version: %s\n", appVersion)
		fmt.Printf("Configuration: %+v\n", config)
		return
	}

	if *version {
		fmt.Println(appName, appVersion)
		return
	}

	if !slices.Contains(allowedTimeFormats[:], *timeFormat) {
		log.Fatalf("Error: invalid time format '%s'. Allowed values: %s", *timeFormat, strings.Join(allowedTimeFormats[:], ", "))
	}

	timeZoneLocation, err := time.LoadLocation(*timeZone)
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}

	audioContext, err := InitializeAudioContext()
	if err != nil {
		log.Fatalf("Failed to initialize audio context: %v", err)
	}

	// Initialize screen early to show loading message
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Failed to create screen: %v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("Failed to initialize screen: %v", err)
	}
	defer screen.Fini()

	display := NewDisplay(screen)

	if *offline {
		display.SetInitText("Offline mode: using system time")
	} else {
		display.SetInitText("Querying NTP server for time...")
	}

	var offset time.Duration
	if *offline {
		offset = 0
	} else {
		ntpTime, err := ntp.Time(*ntpServer)
		if err != nil {
			log.Fatalf("Failed to get time from NTP server %s: %v", *ntpServer, err)
		}
		offset = time.Since(ntpTime)

		go func() {
			ticker := time.NewTicker(ntpOffsetRefreshInterval)
			defer ticker.Stop()
			for range ticker.C {
				ntpTime, err := ntp.Time(*ntpServer)
				if err == nil {
					offset = time.Since(ntpTime)
				}
				// If error, ignore and keep previous offset
			}
		}()
	}

	quitChan := make(chan struct{})
	go display.PollEvents(quitChan)

	displayTicker := time.NewTicker(100 * time.Millisecond)
	defer displayTicker.Stop()

	for {
		select {
		case <-displayTicker.C:
			displayState := DisplayState{
				Now:           time.Now().Add(-offset).In(timeZoneLocation),
				TimeFormat:    *timeFormat,
				HideDate:      *hideDate,
				ShowTimeZone:  *showTimeZone,
				HideStatusbar: *hideStatusbar,
				TimeZone:      timeZoneLocation,
			}
			display.Update(displayState)

			if *beeps {
				now := time.Now().Add(-offset).In(timeZoneLocation)
				BeepTick(audioContext, now)
			}
		case <-quitChan:
			return
		}
	}
}
