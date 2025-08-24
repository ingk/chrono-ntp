package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/beevik/ntp"
	"github.com/gdamore/tcell/v2"
	"github.com/pelletier/go-toml/v2"
)

const defaultNtpServer = "time.google.com"
const defaultTimeFormat = "ISO8601"
const defaultTimeZone = "Local"

var allowedTimeFormats = []string{"ISO8601", "12h", "12h_AM_PM", ".beat"}

type Config struct {
	Server        string `toml:"server"`
	TimeZone      string `toml:"time-zone"`
	HideStatusbar bool   `toml:"hide-statusbar"`
	HideDate      bool   `toml:"hide-date"`
	ShowTimeZone  bool   `toml:"show-time-zone"`
	TimeFormat    string `toml:"time-format"`
	Beeps         bool   `toml:"beeps"`
}

func main() {
	config := loadConfig()

	ntpServer := flag.String("server", config.Server, "NTP server to synchronize time from")
	timeZone := flag.String("time-zone", config.TimeZone, "Time zone name (e.g., 'America/New_York')")
	debug := flag.Bool("debug", false, "Show debug information (e.g., offset from NTP server) and exit")
	hideStatusbar := flag.Bool("hide-statusbar", config.HideStatusbar, "Hide the status bar")
	hideDate := flag.Bool("hide-date", config.HideDate, "Hide the date display")
	showTimeZone := flag.Bool("show-time-zone", config.ShowTimeZone, "Show the time zone")
	timeFormat := flag.String("time-format", config.TimeFormat, fmt.Sprintf("Time display format (%s)", strings.Join(allowedTimeFormats, ", ")))
	beeps := flag.Bool("beeps", config.Beeps, "Play 6 beeps at the end of each minute, with the sixth beep at second 0 (emulates the Greenwich Time Signal)")
	flag.Parse()

	if !slices.Contains(allowedTimeFormats, *timeFormat) {
		log.Fatalf("Error: invalid time format '%s'. Allowed values: %s", *timeFormat, strings.Join(allowedTimeFormats, ", "))
	}

	ntpTime, err := ntp.Time(*ntpServer)
	if err != nil {
		log.Fatalf("Failed to get time from NTP server %s: %v", *ntpServer, err)
	}
	offset := time.Since(ntpTime)

	if *debug {
		log.Printf("NTP server: %s", *ntpServer)
		log.Printf("Offset: %s", offset.String())
		return
	}

	timeZoneLocation, err := time.LoadLocation(*timeZone)
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Failed to create screen: %v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("Failed to initialize screen: %v", err)
	}
	defer screen.Fini()

	screen.Clear()

	quit := false
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	audioContext := InitializeAudioContext()

	go func() {
		for {
			ev := screen.PollEvent()
			switch tev := ev.(type) {
			case *tcell.EventKey:
				if tev.Key() == tcell.KeyCtrlC || slices.Contains([]rune{'q', 'Q'}, tev.Rune()) {
					quit = true
					return
				}
			case *tcell.EventResize:
				screen.Sync()
			}
		}
	}()

	for !quit {
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

		// Beep logic
		if *beeps {
			sec := now.Second()
			// If last 5 seconds of the minute
			if sec >= 55 || sec == 0 {
				// Only beep once per second
				go func(s int) {
					if s == 0 {
						PlayBeep(audioContext, 1*time.Second) // last beep, 1 second
					} else {
						PlayBeep(audioContext, 100*time.Millisecond) // short beep
					}
				}(sec)
			}
		}

		screen.Show()

		<-ticker.C
	}
}

func loadConfig() Config {
	config := Config{
		Server:        defaultNtpServer,
		TimeZone:      defaultTimeZone,
		HideStatusbar: false,
		HideDate:      false,
		ShowTimeZone:  true,
		TimeFormat:    defaultTimeFormat,
		Beeps:         false,
	}

	configPath := filepath.Join(os.Getenv("HOME"), ".chrono-ntp.toml")
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatalf("Failed to read config file %s: %v", configPath, err)
		}
		if err := toml.Unmarshal(data, &config); err != nil {
			log.Fatalf("Failed to parse config file %s: %v", configPath, err)
		}
	}

	return config
}
