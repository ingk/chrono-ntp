package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

type DisplayState struct {
	Now           time.Time
	TimeFormat    string
	HideDate      bool
	ShowTimeZone  bool
	HideStatusbar bool
	TimeZone      *time.Location
}

type Display struct {
	screen tcell.Screen
}

func NewDisplay(screen tcell.Screen) *Display {
	return &Display{screen: screen}
}

func (d *Display) SetInitText(text string) {
	_, height := d.screen.Size()
	centerY := height/2 - 1

	d.screen.Clear()
	drawTextCentered(d.screen, centerY, text, tcell.StyleDefault.Bold(true))
	d.screen.Show()
}

func (d *Display) Update(state DisplayState) {
	d.screen.Clear()

	_, height := d.screen.Size()
	centerY := height/2 - 1

	drawTextCentered(d.screen, centerY, formatTime(state.Now, &state.TimeFormat), tcell.StyleDefault.Bold(true))

	if !state.HideDate {
		drawTextCentered(d.screen, centerY-1, formatDate(state.Now), tcell.StyleDefault)
	}

	if state.ShowTimeZone {
		drawTextCentered(d.screen, centerY+1, normalizeTimezoneName(state.TimeZone), tcell.StyleDefault)
	}

	if !state.HideStatusbar {
		drawStatusbar(d.screen)
	}

	d.screen.Show()
}
