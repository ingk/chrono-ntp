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
	if *timeFormat == ".beat" {
		return formatBeatTime(t)
	}

	timeFormatMap := map[string]string{
		"ISO8601":   "15:04:05",
		"12h":       "03:04:05",
		"12h_AM_PM": "03:04:05 PM",
	}

	return t.Format(timeFormatMap[*timeFormat])
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
