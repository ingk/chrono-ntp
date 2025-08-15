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
const defaultTimezone = "Local"

var allowedTimeFormats = []string{"ISO8601", "12h", "12h_AM_PM", ".beat"}

type Config struct {
	Server        string `toml:"server"`
	Timezone      string `toml:"timezone"`
	Debug         bool   `toml:"debug"`
	HideStatusbar bool   `toml:"hide-statusbar"`
	HideDate      bool   `toml:"hide-date"`
	ShowTimezone  bool   `toml:"show-timezone"`
	TimeFormat    string `toml:"time-format"`
}

func main() {
	config := loadConfig()

	ntpServer := flag.String("server", config.Server, "NTP server to sync time from")
	timezone := flag.String("timezone", config.Timezone, "Name of the timezone (e.g., 'America/New_York')")
	debug := flag.Bool("debug", config.Debug, "Show debug information (e.g. offset from NTP server), then exit")
	hideStatusbar := flag.Bool("hide-statusbar", config.HideStatusbar, "Hide the status bar")
	hideDate := flag.Bool("hide-date", config.HideDate, "Hide the current date")
	showTimezone := flag.Bool("show-timezone", config.ShowTimezone, "Show the timezone")
	timeFormat := flag.String("time-format", config.TimeFormat, fmt.Sprintf("Format for displaying time (%s)", strings.Join(allowedTimeFormats, ", ")))
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

	timezoneLocation, err := time.LoadLocation(*timezone)
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
		now := time.Now().Add(-offset).In(timezoneLocation)

		_, height := screen.Size()
		centerY := height/2 - 1

		drawTextCentered(screen, centerY, formatTime(now, timeFormat), tcell.StyleDefault.Bold(true))

		if !*hideDate {
			drawTextCentered(screen, centerY-1, formatDate(now), tcell.StyleDefault)
		}

		if *showTimezone {
			drawTextCentered(screen, centerY+1, normalizeTimezoneName(timezoneLocation), tcell.StyleDefault)
		}

		if !*hideStatusbar {
			drawStatusbar(screen)
		}

		screen.Show()

		<-ticker.C
	}
}

func loadConfig() Config {
	config := Config{
		Server:        defaultNtpServer,
		Timezone:      defaultTimezone,
		Debug:         false,
		HideStatusbar: false,
		HideDate:      false,
		ShowTimezone:  true,
		TimeFormat:    defaultTimeFormat,
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
