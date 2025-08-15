package main

import (
	"strings"
	"time"
)

func formatDate(t time.Time) string {
	return t.Format(time.DateOnly)
}

func formatTime(t time.Time, timeFormat *string) string {
	timeFormatMap := map[string]string{
		"ISO8601":   "15:04:05",
		"12h":       "03:04:05",
		"12h_AM_PM": "03:04:05 PM",
	}

	return t.Format(timeFormatMap[*timeFormat])
}

func normalizeTimezoneName(location *time.Location) string {
	// Replace underscores with spaces for better readability
	return strings.ReplaceAll(location.String(), "_", " ")
}
