package main

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
	result := FormatDate(inputTime)
	if result != "2023-10-01" {
		t.Errorf("Expected '2023-10-01', got '%s'", result)
	}
}

func TestFormatTime(t *testing.T) {
	inputTime := time.Date(2023, 10, 1, 15, 16, 17, 0, time.UTC)

	result := FormatTime(inputTime, sPtr("ISO8601"))
	if result != "15:16:17" {
		t.Errorf("Expected '15:16:17' for ISO8601, got '%s'", result)
	}

	result = FormatTime(inputTime, sPtr("12h"))
	if result != "03:16:17" {
		t.Errorf("Expected '03:16:17' for 12h, got '%s'", result)
	}

	// .beat test: inputTime is 2023-10-01 20:13:14 UTC
	// UTC+1 = 21:13:14, seconds since midnight = 21*3600+13*60+14 = 76414
	// .beat = 76414 / 86.4 = 884.5 -> @884
	result = FormatTime(inputTime, sPtr(".beat"))
	if result != "@677" {
		t.Errorf("Expected '@677' for .beat, got '%s'", result)
	}

	// Test .beat with zero padding
	result = FormatTime(time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), sPtr(".beat"))
	if result != "@041" {
		t.Errorf("Expected '@041', got '%s'", result)
	}

	result = FormatTime(inputTime, sPtr("septimal"))
	if result != "43 11 52" {
		t.Errorf("Expected '43 11 52' for septimal, got '%s'", result)
	}

	result = FormatTime(inputTime, sPtr("mars"))
	if result != "23:42:49" {
		t.Errorf("Expected '23:42:49' for mars, got '%s'", result)
	}

	result = FormatTime(inputTime, sPtr("lunar"))
	if result != "393:56:10" {
		t.Errorf("Expected '393:56:10' for lunar, got '%s'", result)
	}
}
