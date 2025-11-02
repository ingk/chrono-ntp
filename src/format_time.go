package main

import (
	"fmt"
	"strings"
	"time"
)

func formatDate(t time.Time) string {
	return t.Format(time.DateOnly)
}

func formatTime(t time.Time, timeFormat *string) string {
	switch *timeFormat {
	case ".beat":
		return formatBeatTime(t)
	case "septimal":
		return formatSeptimalTime(t)
	default:
		timeFormatMap := map[string]string{
			"ISO8601":   "15:04:05",
			"12h":       "03:04:05",
			"12h_AM_PM": "03:04:05 PM",
		}
		return t.Format(timeFormatMap[*timeFormat])
	}
}

// formatSeptimalTime returns the time in septimal format (base-7 pairs)
// See: http://the-light.com/cal/veseptimal.html
func formatSeptimalTime(t time.Time) string {
	// Get time since midnight in local time
	year, month, day := t.Date()
	midnight := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	msSinceMidnight := float64(t.Sub(midnight).Milliseconds())

	// Use the original JS divisors for accuracy
	sep1 := int(msSinceMidnight / 12342857.14)
	sep2 := int(msSinceMidnight/1763265.306) % 7
	sep3 := int(msSinceMidnight/251895.0437) % 7
	sep4 := int(msSinceMidnight/35985.00625) % 7
	sep5 := int(msSinceMidnight/5140.715178) % 7
	sep6 := int((msSinceMidnight / 734.3878826)) % 7

	return fmt.Sprintf("%d%d %d%d %d%d", sep1, sep2, sep3, sep4, sep5, sep6)
}

// formatBeatTime returns Swatch Internet Time (.beat)
// @see https://en.wikipedia.org/wiki/Swatch_Internet_Time
func formatBeatTime(t time.Time) string {
	// Convert time to UTC+1 (Biel Mean Time)
	bmt := t.UTC().Add(1 * time.Hour)
	seconds := bmt.Hour()*3600 + bmt.Minute()*60 + bmt.Second()
	beat := float64(seconds) / 86.4
	return fmt.Sprintf("@%s", formatBeat(beat))
}

func formatBeat(beat float64) string {
	// .beat is always 3 digits, rounded down
	return leftPadInt(int(beat), 3)
}

func leftPadInt(n int, width int) string {
	s := fmt.Sprintf("%d", n)
	for len(s) < width {
		s = "0" + s
	}
	return s
}

func normalizeTimezoneName(location *time.Location) string {
	// Replace underscores with spaces for better readability
	return strings.ReplaceAll(location.String(), "_", " ")
}
