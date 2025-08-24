package main

import (
	"testing"
)

func TestParseConfiguration_Defaults(t *testing.T) {
	config := parseConfiguration(nil)

	if config.Server != "time.google.com" {
		t.Errorf("expected Server %q, got %q", "time.google.com", config.Server)
	}
	if config.TimeZone != "Local" {
		t.Errorf("expected TimeZone %q, got %q", "Local", config.TimeZone)
	}
	if config.HideStatusbar != false {
		t.Errorf("expected HideStatusbar false, got %v", config.HideStatusbar)
	}
	if config.HideDate != false {
		t.Errorf("expected HideDate false, got %v", config.HideDate)
	}
	if config.ShowTimeZone != true {
		t.Errorf("expected ShowTimeZone true, got %v", config.ShowTimeZone)
	}
	if config.TimeFormat != "ISO8601" {
		t.Errorf("expected TimeFormat %q, got %q", "ISO8601", config.TimeFormat)
	}
	if config.Beeps != false {
		t.Errorf("expected Beeps false, got %v", config.Beeps)
	}
}

func TestParseConfiguration_Content(t *testing.T) {
	tomlContent := `
server = "pool.example-time-server.org"
time-zone = "Europe/Berlin"
hide-statusbar = true
hide-date = true
show-time-zone = true
time-format = "12h_AM_PM"
beeps = true
`
	config := parseConfiguration([]byte(tomlContent))

	if config.Server != "pool.example-time-server.org" {
		t.Errorf("expected Server 'pool.example-time-server.org', got %q", config.Server)
	}
	if config.TimeZone != "Europe/Berlin" {
		t.Errorf("expected TimeZone 'Europe/Berlin', got %q", config.TimeZone)
	}
	if config.HideStatusbar != true {
		t.Errorf("expected HideStatusbar true, got %v", config.HideStatusbar)
	}
	if config.HideDate != true {
		t.Errorf("expected HideDate true, got %v", config.HideDate)
	}
	if config.ShowTimeZone != true {
		t.Errorf("expected ShowTimeZone true, got %v", config.ShowTimeZone)
	}
	if config.TimeFormat != "12h_AM_PM" {
		t.Errorf("expected TimeFormat '12h_AM_PM', got %q", config.TimeFormat)
	}
	if config.Beeps != true {
		t.Errorf("expected Beeps true, got %v", config.Beeps)
	}
}
