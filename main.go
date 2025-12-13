package main

import (
	"flag"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"chrono-ntp/audio"
	"chrono-ntp/configuration"
	"chrono-ntp/display"
	"chrono-ntp/ntp"
)

const (
	appName                  = "chrono-ntp"
	appVersion               = "dev"
	ntpOffsetRefreshInterval = 15 * time.Minute
)

var allowedTimeFormats = display.AllowedTimeFormats[:]
var allowedDateFormats = display.AllowedDateFormats[:]

func main() {
	config, err := configuration.LoadConfiguration()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ntpServer := flag.String("server", config.Server, "NTP server to synchronize time from")
	timeZone := flag.String("time-zone", config.TimeZone, "Time zone name (e.g., 'America/New_York')")
	debug := flag.Bool("debug", false, "Show debug information (e.g., offset from NTP server) and exit")
	hideStatusbar := flag.Bool("hide-statusbar", config.HideStatusbar, "Hide the status bar")
	hideDate := flag.Bool("hide-date", config.HideDate, "Hide the date display")
	showTimeZone := flag.Bool("show-time-zone", config.ShowTimeZone, "Show the time zone")
	dateFormat := flag.String("date-format", "YYYY-MM-DD", fmt.Sprintf("Date display format (%s)", strings.Join(allowedDateFormats, ", ")))
	timeFormat := flag.String("time-format", config.TimeFormat, fmt.Sprintf("Time display format (%s)", strings.Join(allowedTimeFormats, ", ")))
	beeps := flag.Bool("beeps", config.Beeps, "Play 6 beeps at the end of each minute, with the sixth beep at second 0 (emulates the Greenwich Time Signal)")
	version := flag.Bool("version", false, "Show version and exit")
	offline := flag.Bool("offline", false, "Run in offline mode (use system time, ignore NTP server)")
	writeConfig := flag.Bool("write-config", false, "Write configuration file (merged from existing configuration file and flags)")
	flag.Parse()

	beepsEnabled := *beeps && !slices.Contains([]string{".beat", "septimal", "lunar", "mars"}, *timeFormat)

	if *debug {
		fmt.Printf("Version: %s\n", appVersion)
		fmt.Printf("Configuration: %+v\n", config)
		return
	}

	if *version {
		fmt.Println(appName, appVersion)
		return
	}

	if !slices.Contains(allowedTimeFormats, *timeFormat) {
		log.Fatalf("Error: invalid time format '%s'. Allowed values: %s", *timeFormat, strings.Join(allowedTimeFormats, ", "))
	}

	if *writeConfig {
		mergedConfig := configuration.Configuration{
			Server:        *ntpServer,
			TimeZone:      *timeZone,
			HideStatusbar: *hideStatusbar,
			HideDate:      *hideDate,
			ShowTimeZone:  *showTimeZone,
			TimeFormat:    *timeFormat,
			Beeps:         *beeps,
			Offline:       *offline,
		}
		configPath, err := configuration.WriteConfiguration(mergedConfig)
		if err == nil {
			fmt.Printf("Configuration written to %s\n", configPath)
		} else {
			log.Fatalf("Failed to write configuration (%s): %v", configPath, err)
		}
		return
	}

	timeZoneLocation, err := time.LoadLocation(*timeZone)
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}

	audioContext, err := audio.InitializeAudioContext()
	if err != nil {
		log.Fatalf("Failed to initialize audio context: %v", err)
	}

	// Initialize display early to show loading message
	d, err := display.NewDisplay()
	if err != nil {
		log.Fatalf("Failed to create display: %v", err)
	}
	if err := d.Init(); err != nil {
		log.Fatalf("Failed to initialize display: %v", err)
	}
	defer d.Finalize()

	var offset time.Duration

	if *offline {
		offset = 0
	} else {
		d.SetInitText("Querying NTP server for time...")

		ntpClient, err := ntp.NewNtp(*ntpServer)
		if err != nil {
			log.Fatalf("Failed to get time from NTP server %s: %v", *ntpServer, err)
		}
		offset = ntpClient.Offset()

		go func() {
			ticker := time.NewTicker(ntpOffsetRefreshInterval)
			defer ticker.Stop()
			for range ticker.C {
				if err := ntpClient.Refresh(); err == nil {
					offset = ntpClient.Offset()
				}
				// If error, ignore and keep previous offset
			}
		}()
	}

	quitChan := make(chan struct{})
	go d.PollEvents(quitChan)

	displayTicker := time.NewTicker(100 * time.Millisecond)
	defer displayTicker.Stop()

	for {
		select {
		case <-displayTicker.C:
			now := time.Now().Add(-offset).In(timeZoneLocation)

			displayState := &display.DisplayState{
				Now:           now,
				DateFormat:    *dateFormat,
				TimeFormat:    *timeFormat,
				HideDate:      *hideDate,
				ShowTimeZone:  *showTimeZone,
				HideStatusbar: *hideStatusbar,
				TimeZone:      timeZoneLocation,
				Offset:        offset,
				Offline:       *offline,
			}
			d.Update(*displayState)

			if beepsEnabled {
				audio.BeepTick(audioContext, now)
			}
		case <-quitChan:
			return
		}
	}
}
