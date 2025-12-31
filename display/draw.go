package display

import (
	"chrono-ntp/timesource"
	"strconv"

	"github.com/gdamore/tcell/v2"
)

var (
	statusBarQuitLabel    = "Quit"
	statusBarQuitShortcut = "Q, <C-c>"
	statusBarOffsetLabel  = "Offset"
)

func drawStatusBar(screen tcell.Screen, state DisplayState) {
	_, height := screen.Size()
	y := height - 1

	for i, r := range statusBarQuitShortcut {
		screen.SetContent(i, y, r, nil, tcell.StyleDefault.Bold(true).Reverse(true))
	}
	x := len(statusBarQuitShortcut) + 1
	for i, r := range statusBarQuitLabel {
		screen.SetContent(x+i, y, r, nil, tcell.StyleDefault)
	}

	x = x + len(statusBarQuitLabel) + 4
	for i, r := range statusBarOffsetLabel {
		screen.SetContent(x+i, y, r, nil, tcell.StyleDefault.Bold(true).Reverse(true))
	}

	x = x + len(statusBarOffsetLabel) + 1
	offset := strconv.FormatInt(state.TimeStatus.Offset.Milliseconds(), 10) + "ms"
	if state.Offline {
		offset = "(offline)"
	} else if state.TimeStatus.State == timesource.StateSyncing {
		offset = "(syncing)"
	}
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
