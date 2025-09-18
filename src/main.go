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

var appName = "chrono-ntp"
var appVersion = "dev"

var allowedTimeFormats = []string{"ISO8601", "12h", "12h_AM_PM", ".beat"}

func main() {
	config := LoadConfiguration()

	ntpServer := flag.String("server", config.Server, "NTP server to synchronize time from")
	timeZone := flag.String("time-zone", config.TimeZone, "Time zone name (e.g., 'America/New_York')")
	debug := flag.Bool("debug", false, "Show debug information (e.g., offset from NTP server) and exit")
	hideStatusbar := flag.Bool("hide-statusbar", config.HideStatusbar, "Hide the status bar")
	hideDate := flag.Bool("hide-date", config.HideDate, "Hide the date display")
	showTimeZone := flag.Bool("show-time-zone", config.ShowTimeZone, "Show the time zone")
	timeFormat := flag.String("time-format", config.TimeFormat, fmt.Sprintf("Time display format (%s)", strings.Join(allowedTimeFormats, ", ")))
	beeps := flag.Bool("beeps", config.Beeps, "Play 6 beeps at the end of each minute, with the sixth beep at second 0 (emulates the Greenwich Time Signal)")
	version := flag.Bool("version", false, "Show version and exit")
	flag.Parse()

	if *version {
		fmt.Println(appName, appVersion)
		return
	}

	if !slices.Contains(allowedTimeFormats, *timeFormat) {
		log.Fatalf("Error: invalid time format '%s'. Allowed values: %s", *timeFormat, strings.Join(allowedTimeFormats, ", "))
	}

	if *debug {
		log.Printf("NTP server: %s", *ntpServer)
		return
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

	_, height := screen.Size()
	centerY := height/2 - 1

	screen.Clear()
	drawTextCentered(screen, centerY, "Querying NTP server for time...", tcell.StyleDefault.Bold(true))
	screen.Show()

	ntpTime, err := ntp.Time(*ntpServer)
	if err != nil {
		log.Fatalf("Failed to get time from NTP server %s: %v", *ntpServer, err)
	}
	offset := time.Since(ntpTime)

	timeZoneLocation, err := time.LoadLocation(*timeZone)
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}

	audioContext, err := InitializeAudioContext()
	if err != nil {
		log.Fatalf("Failed to initialize audio context: %v", err)
	}

	quitChan := make(chan struct{})

	go func() {
		for {
			ev := screen.PollEvent()
			switch tev := ev.(type) {
			case *tcell.EventKey:
				if tev.Key() == tcell.KeyCtrlC || slices.Contains([]rune{'q', 'Q'}, tev.Rune()) {
					quitChan <- struct{}{}
					return
				}
			case *tcell.EventResize:
				screen.Sync()
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			screen.Clear()

			now := time.Now().Add(-offset).In(timeZoneLocation)

			_, height := screen.Size()
			centerY := height/2 - 1

			drawTextCentered(screen, centerY, formatTime(now, timeFormat), tcell.StyleDefault.Bold(true))

			if !*hideDate {
				drawTextCentered(screen, centerY-1, formatDate(now), tcell.StyleDefault)
			}

			if *showTimeZone {
				drawTextCentered(screen, centerY+1, normalizeTimezoneName(timeZoneLocation), tcell.StyleDefault)
			}

			if !*hideStatusbar {
				drawStatusbar(screen)
			}

			if *beeps {
				sec := now.Second()
				if sec >= 55 || sec == 0 {
					go func(s int) {
						if s == 0 {
							PlayLongBeep(audioContext)
						} else {
							PlayShortBeep(audioContext)
						}
					}(sec)
				}
			}

			screen.Show()
		case <-quitChan:
			return
		}
	}
}
