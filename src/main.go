package main

import (
	"flag"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/beevik/ntp"
	"github.com/gdamore/tcell/v2"
)

func drawStatusbar(screen tcell.Screen) {
	statusbarQuitLabel := "Quit"
	statusbarQuitShortcut := "Q, <C-c>"
	width, height := screen.Size()

	x := (width - len(statusbarQuitLabel+statusbarQuitShortcut) + 1) / 2
	for i, r := range statusbarQuitLabel {
		screen.SetContent(x+i, height-2, r, nil, tcell.StyleDefault.Bold(true))
	}
	for i, r := range " " + statusbarQuitShortcut {
		screen.SetContent(x+4+i, height-2, r, nil, tcell.StyleDefault)
	}
}

func drawTextCentered(s tcell.Screen, y int, text string, style tcell.Style) {
	w, _ := s.Size()
	x := (w - len(text)) / 2
	for i, r := range text {
		s.SetContent(x+i, y, r, nil, style)
	}
}

func normalizeTimezoneName(location *time.Location) string {
	// Replace underscores with spaces for better readability
	return strings.ReplaceAll(location.String(), "_", " ")
}

func main() {
	ntpServer := flag.String("server", "time.google.com", "NTP server to sync time from")
	timezone := flag.String("timezone", "Local", "NTP server to sync time from")
	hideStatusbar := flag.Bool("hide-statusbar", false, "Hide the status bar")
	hideDate := flag.Bool("hide-date", false, "Hide the current date")
	showTimezone := flag.Bool("show-timezone", false, "Show the timezone")
	flag.Parse()

	ntpTime, err := ntp.Time(*ntpServer)
	if err != nil {
		log.Fatalf("Failed to get time from NTP server %s: %v", *ntpServer, err)
	}
	offset := time.Since(ntpTime)

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

		timeStr := now.Format("15:04:05")
		drawTextCentered(screen, centerY, timeStr, tcell.StyleDefault.Bold(true))

		if !*hideDate {
			dateStr := now.Format("2006-01-02")
			drawTextCentered(screen, centerY-1, dateStr, tcell.StyleDefault)
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
