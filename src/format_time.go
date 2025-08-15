package main

import (
	"strings"
	"time"
)

func formatDate(time time.Time) string {
	return time.Format("2006-01-02")
}

func formatTime(time time.Time) string {
	return time.Format("15:04:05")
}

func normalizeTimezoneName(location *time.Location) string {
	// Replace underscores with spaces for better readability
	return strings.ReplaceAll(location.String(), "_", " ")
}
