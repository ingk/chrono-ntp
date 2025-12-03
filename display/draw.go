package display

import (
	"github.com/gdamore/tcell/v2"
)

var (
	statusbarQuitLabel    = "Quit"
	statusbarQuitShortcut = "Q, <C-c>"
)

func drawStatusbar(screen tcell.Screen) {
	_, height := screen.Size()

	for i, r := range statusbarQuitShortcut {
		screen.SetContent(i, height-1, r, nil, tcell.StyleDefault.Bold(true).Reverse(true))
	}
	x := len(statusbarQuitShortcut) + 1
	for i, r := range " " + statusbarQuitLabel {
		screen.SetContent(x+i, height-1, r, nil, tcell.StyleDefault)
	}
}

func drawTextCentered(s tcell.Screen, y int, text string, style tcell.Style) {
	w, _ := s.Size()
	x := (w - len(text)) / 2
	for i, r := range text {
		s.SetContent(x+i, y, r, nil, style)
	}
}
