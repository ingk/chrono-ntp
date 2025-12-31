package timesource

import "testing"

func TestNtpTimeSource_Name(t *testing.T) {
	timeSource := NewNtpTimeSource("ntp.example.org")
	if timeSource.Name() != "NTP" {
		t.Errorf("Expected name 'NTP', got '%s'", timeSource.Name())
	}
	if timeSource.lastStatus.State != StateSyncing {
		t.Errorf("Expected initial state 'StateSyncing', got '%s'", timeSource.lastStatus.State)
	}
}

func TestNtpTimeSource_TimeStatus(t *testing.T) {
	timeSource := NewNtpTimeSource("ntp.example.org")

	status := timeSource.TimeStatus()
	if status.Source != "NTP" {
		t.Errorf("Expected source 'NTP', got '%s'", status.Source)
	}
	if status.State != StateSyncing {
		t.Errorf("Expected state 'StateSyncing', got '%s'", status.State)
	}
	if status.Offset != 0 {
		t.Errorf("Expected offset 0, got '%d'", status.Offset)
	}
}
