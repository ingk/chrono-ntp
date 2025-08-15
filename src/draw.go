package main

import (
	"github.com/gdamore/tcell/v2"
)

func drawStatusbar(screen tcell.Screen) {
	statusbarQuitLabel := "Quit"
	statusbarQuitShortcut := "Q, <C-c>"
	width, height := screen.Size()

	x := (width - len(statusbarQuitLabel+statusbarQuitShortcut) + 1) / 2
	for i, r := range statusbarQuitLabel {
		screen.SetContent(x+i, height-1, r, nil, tcell.StyleDefault.Bold(true))
	}
	for i, r := range " " + statusbarQuitShortcut {
		screen.SetContent(x+4+i, height-1, r, nil, tcell.StyleDefault)
	}
}

func drawTextCentered(s tcell.Screen, y int, text string, style tcell.Style) {
	w, _ := s.Size()
	x := (w - len(text)) / 2
	for i, r := range text {
		s.SetContent(x+i, y, r, nil, style)
	}
}
