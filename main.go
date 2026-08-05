package main

import (
	"flag"
	"fmt"
	"log"
	"slices"
	"strings"

	"chrono-ntp/app"
	"chrono-ntp/audio"
	"chrono-ntp/configuration"
	"chrono-ntp/display"
)

const (
	appName    = "chrono-ntp"
	appVersion = "dev"
)

var allowedTimeFormats = display.AllowedTimeFormats[:]
var allowedDateFormats = display.AllowedDateFormats[:]
var allowedBeepPatterns = audio.AllowedBeepPatterns[:]

func main() {
	config, err := configuration.LoadConfiguration()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ntpServer := flag.String("server", config.Server, "NTP server to synchronize time from")
	timeZone := flag.String("time-zone", config.TimeZone, "Time zone name (e.g., 'America/New_York')")
	debug := flag.Bool("debug", false, "Show debug information (e.g., offset from NTP server) and exit")
	hideStatusBar := flag.Bool("hide-status-bar", config.HideStatusBar, "Hide the status bar")
	hideDate := flag.Bool("hide-date", config.HideDate, "Hide the date display")
	showTimeZone := flag.Bool("show-time-zone", config.ShowTimeZone, "Show the time zone")
	dateFormat := flag.String("date-format", "YYYY-MM-DD", fmt.Sprintf("Date display format (%s)", strings.Join(allowedDateFormats, ", ")))
	timeFormat := flag.String("time-format", config.TimeFormat, fmt.Sprintf("Time display format (%s)", strings.Join(allowedTimeFormats, ", ")))
	beepPattern := flag.String("beep-pattern", config.BeepPattern, fmt.Sprintf("Beep pattern (%s)", strings.Join(allowedBeepPatterns, ", ")))
	version := flag.Bool("version", false, "Show version and exit")
	offline := flag.Bool("offline", false, "Run in offline mode (use system time, ignore NTP server)")
	writeConfig := flag.Bool("write-config", false, "Write configuration file (merged from existing configuration file and flags)")
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

	if !slices.Contains(allowedDateFormats, *dateFormat) {
		log.Fatalf("Error: invalid date format '%s'. Allowed values: %s", *dateFormat, strings.Join(allowedDateFormats, ", "))
	}

	if !slices.Contains(allowedTimeFormats, *timeFormat) {
		log.Fatalf("Error: invalid time format '%s'. Allowed values: %s", *timeFormat, strings.Join(allowedTimeFormats, ", "))
	}

	if *beepPattern != "" && !slices.Contains(allowedBeepPatterns, *beepPattern) {
		log.Fatalf("Error: invalid beep pattern '%s'. Allowed patterns: %s", *beepPattern, strings.Join(allowedBeepPatterns, ", "))
	}
	beepEnabled := *beepPattern != "" && slices.Contains([]string{"ISO8601", "12h", "12h_AM_PM"}, *timeFormat)
	if !beepEnabled {
		*beepPattern = "silence"
	}

	if *writeConfig {
		mergedConfig := configuration.Configuration{
			Server:        *ntpServer,
			TimeZone:      *timeZone,
			HideStatusBar: *hideStatusBar,
			HideDate:      *hideDate,
			ShowTimeZone:  *showTimeZone,
			TimeFormat:    *timeFormat,
			BeepPattern:   *beepPattern,
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

	audioContext, disp, timeSource, timeZoneLocation, beepStrategy := app.Init(app.InitOptions{
		Server:      *ntpServer,
		TimeZone:    *timeZone,
		Offline:     *offline,
		BeepPattern: *beepPattern,
	})

	app.Run(app.RunOptions{
		AudioContext:  audioContext,
		Display:       disp,
		TimeSource:    timeSource,
		Location:      timeZoneLocation,
		BeepStrategy:  beepStrategy,
		HideStatusBar: *hideStatusBar,
		HideDate:      *hideDate,
		ShowTimeZone:  *showTimeZone,
		DateFormat:    *dateFormat,
		TimeFormat:    *timeFormat,
		Offline:       *offline,
	})
}
