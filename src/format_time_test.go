package main

import (
	"testing"
	"time"
)

func sPtr(s string) *string { return &s }

func TestNormalizeTimezoneName(t *testing.T) {
	loc, _ := time.LoadLocation("America/New_York")
	result := normalizeTimezoneName(loc)
	if result != "America/New York" {
		t.Errorf("Expected 'America/New York', got '%s'", result)
	}

	loc, _ = time.LoadLocation("UTC")
	result = normalizeTimezoneName(loc)
	if result != "UTC" {
		t.Errorf("Expected 'UTC', got '%s'", result)
	}
}

func TestFormatDate(t *testing.T) {
	inputTime := time.Date(2023, 10, 1, 12, 13, 14, 0, time.UTC)
	result := formatDate(inputTime)
	if result != "2023-10-01" {
		t.Errorf("Expected '2023-10-01', got '%s'", result)
	}
}

func TestFormatTime(t *testing.T) {
	inputTime := time.Date(2023, 10, 1, 20, 13, 14, 0, time.UTC)

	result := formatTime(inputTime, sPtr("ISO8601"))
	if result != "20:13:14" {
		t.Errorf("Expected '20:13:14', got '%s'", result)
	}

	result = formatTime(inputTime, sPtr("12h"))
	if result != "08:13:14" {
		t.Errorf("Expected '08:13:14', got '%s'", result)
	}
}
