package main

import (
	"flag"
	"log"
	"time"

	"github.com/beevik/ntp"
	"github.com/gdamore/tcell/v2"
)

func drawTextCentered(s tcell.Screen, y int, text string, style tcell.Style) {
	w, _ := s.Size()
	x := (w - len(text)) / 2
	for i, r := range text {
		s.SetContent(x+i, y, r, nil, style)
	}
}

func main() {
	ntpServer := flag.String("server", "time.google.com", "NTP server to sync time from")
	timezone := flag.String("timezone", "Local", "NTP server to sync time from")
	hideStatusbar := flag.Bool("hide-statusbar", false, "Hide the status bar")
	flag.Parse()

	ntpTime, err := ntp.Time(*ntpServer)
	if err != nil {
		log.Fatalf("Failed to get time from NTP server %s: %v", *ntpServer, err)
	}
	offset := time.Since(ntpTime)

	timezoneLocation, err := time.LoadLocation(*timezone)
	if err != nil {
		log.Fatalf("failed to load location: %v", err)
	}

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Failed to create screen: %v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("Failed to initialize screen: %v", err)
	}
	defer s.Fini()

	s.Clear()

	quit := false
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	go func() {
		for {
			ev := s.PollEvent()
			switch tev := ev.(type) {
			case *tcell.EventKey:
				if tev.Key() == tcell.KeyCtrlC || tev.Rune() == 'q' || tev.Rune() == 'Q' {
					quit = true
					return
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	for !quit {
		now := time.Now().Add(-offset).In(timezoneLocation)
		s.Clear()

		w, h := s.Size()

		dateStr := now.Format("2006-01-02")
		timeStr := now.Format("15:04:05")
		dateStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)
		boldStyle := tcell.StyleDefault.Bold(true).Foreground(tcell.ColorWhite)
		centerY := h/2 - 1
		drawTextCentered(s, centerY, dateStr, dateStyle)
		drawTextCentered(s, centerY+1, timeStr, boldStyle)

		if !*hideStatusbar {
			x := (w - len("Quit Q, <C-c>")) / 2
			for i, r := range "Quit" {
				s.SetContent(x+i, h-2, r, nil, tcell.StyleDefault.Bold(true).Foreground(tcell.ColorWhite))
			}
			for i, r := range " Q, <C-c>" {
				s.SetContent(x+4+i, h-2, r, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			}
		}

		s.Show()

		<-ticker.C
	}
}
