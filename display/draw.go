package display

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
)

var (
	statusbarQuitLabel    = "Quit"
	statusbarQuitShortcut = "Q, <C-c>"
	statusbarOffsetLabel  = "Offset"
)

func drawStatusbar(screen tcell.Screen, state DisplayState) {
	_, height := screen.Size()
	y := height - 1

	for i, r := range statusbarQuitShortcut {
		screen.SetContent(i, y, r, nil, tcell.StyleDefault.Bold(true).Reverse(true))
	}
	x := len(statusbarQuitShortcut) + 1
	for i, r := range statusbarQuitLabel {
		screen.SetContent(x+i, y, r, nil, tcell.StyleDefault)
	}

	x = x + len(statusbarQuitLabel) + 4
	for i, r := range statusbarOffsetLabel {
		screen.SetContent(x+i, y, r, nil, tcell.StyleDefault.Bold(true).Reverse(true))
	}

	x = x + len(statusbarOffsetLabel) + 1
	offset := strconv.FormatInt(state.Offset.Milliseconds(), 10) + "ms"
	for i, r := range offset {
		screen.SetContent(x+i, y, r, nil, tcell.StyleDefault)
	}
}

func drawTextCentered(s tcell.Screen, y int, text string, style tcell.Style) {
	w, _ := s.Size()
	x := (w - len(text)) / 2
	for i, r := range text {
		s.SetContent(x+i, y, r, nil, style)
	}
}
