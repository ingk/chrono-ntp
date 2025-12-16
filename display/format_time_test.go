package display

import (
	"testing"
	"time"
)

func sPtr(s string) *string { return &s }

func TestNormalizeTimezoneName(t *testing.T) {
	loc, _ := time.LoadLocation("America/New_York")
	result := normalizeTimeZoneName(loc)
	if result != "America/New York" {
		t.Errorf("Expected 'America/New York', got '%s'", result)
	}

	loc, _ = time.LoadLocation("UTC")
	result = normalizeTimeZoneName(loc)
	if result != "UTC" {
		t.Errorf("Expected 'UTC', got '%s'", result)
	}
}

func TestFormatDate(t *testing.T) {
	inputTime := time.Date(2023, 10, 1, 12, 13, 14, 0, time.UTC)

	tests := []struct {
		format   string
		expected string
	}{
		{"unknown-format", "2023-10-01"},
		{"YYYY-MM-DD", "2023-10-01"},
		{"DD/MM/YYYY", "01/10/2023"},
		{"MM/DD/YYYY", "10/01/2023"},
		{"DD.MM.YYYY", "01.10.2023"},
	}

	for _, tt := range tests {
		got := FormatDate(inputTime, &tt.format)
		if got != tt.expected {
			t.Errorf("FormatDate(%q): expected '%s', got '%s'", tt.format, tt.expected, got)
		}
	}
}

func TestFormatTime(t *testing.T) {
	inputTime := time.Date(2023, 10, 1, 15, 16, 17, 0, time.UTC)
	tests := []struct {
		format   string
		expected string
	}{
		{"ISO8601", "15:16:17"},
		{"12h", "03:16:17"},
		{"12h_AM_PM", "03:16:17 PM"},
		{".beat", "@677.97"},
		{"septimal", "43 11 52"},
		{"mars", "23:42:49"},
		{"lunar", "393:56:10"},
		{"unix", "1696173377"},
	}

	for _, tt := range tests {
		format := tt.format
		result := FormatTime(inputTime, &format)
		if result != tt.expected {
			t.Errorf("FormatTime(%s): got %s, want %s", tt.format, result, tt.expected)
		}
	}
}

func TestFormatTime_BeatZeroPadding(t *testing.T) {
	result := FormatTime(time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), sPtr(".beat"))
	if result != "@041.66" {
		t.Errorf("Expected '@041.66', got '%s'", result)
	}
}
