package display

import (
	"slices"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

type DisplayState struct {
	Now           time.Time
	DateFormat    string
	TimeFormat    string
	HideDate      bool
	ShowTimeZone  bool
	HideStatusBar bool
	TimeZone      *time.Location
	Offset        time.Duration
	Offline       bool
}

type Display struct {
	screen tcell.Screen
}

func NewDisplay() (*Display, error) {
	screen, err := tcell.NewScreen()
	return &Display{screen: screen}, err
}

func (d *Display) Init() error {
	return d.screen.Init()
}

func (d *Display) Finalize() {
	d.screen.Fini()
}

func (d *Display) PollEvents(quitChan chan<- struct{}) {
	for {
		ev := d.screen.PollEvent()
		switch tev := ev.(type) {
		case *tcell.EventKey:
			if tev.Key() == tcell.KeyCtrlC || slices.Contains([]rune{'q', 'Q'}, tev.Rune()) {
				quitChan <- struct{}{}
				return
			}
		case *tcell.EventResize:
			d.screen.Sync()
		}
	}
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

	drawTextCentered(d.screen, centerY, FormatTime(state.Now, &state.TimeFormat), tcell.StyleDefault.Bold(true))

	if !state.HideDate {
		drawTextCentered(d.screen, centerY-1, FormatDate(state.Now, &state.DateFormat), tcell.StyleDefault)
	}

	if state.ShowTimeZone {
		var timeZoneLabel string
		switch state.TimeFormat {
		case "mars":
			timeZoneLabel = "Coordinated Mars Time"
		case "lunar":
			timeZoneLabel = "Coordinated Lunar Time"
		default:
			timeZoneLabel = normalizeTimeZoneName(state.TimeZone)
		}
		drawTextCentered(d.screen, centerY+1, timeZoneLabel, tcell.StyleDefault)
	}

	if !state.HideStatusBar {
		drawStatusBar(d.screen, state)
	}

	d.screen.Show()
}

func normalizeTimeZoneName(location *time.Location) string {
	// Replace underscores with spaces for better readability
	return strings.ReplaceAll(location.String(), "_", " ")
}
