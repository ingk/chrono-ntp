package timesource

import "testing"

func TestOfflineTimeSource_Name(t *testing.T) {
	timeSource := NewOfflineTimeSource()
	if timeSource.Name() != "Offline" {
		t.Errorf("Expected name 'Offline', got '%s'", timeSource.Name())
	}
}

func TestOfflineTimeSource_TimeStatus(t *testing.T) {
	timeSource := NewOfflineTimeSource()

	status := timeSource.TimeStatus()
	if status.Source != "Offline" {
		t.Errorf("Expected source 'Offline', got '%s'", status.Source)
	}
	if status.State != StateLocked {
		t.Errorf("Expected state 'StateLocked', got '%s'", status.State)
	}
	if status.Offset != 0 {
		t.Errorf("Expected offset 0, got '%d'", status.Offset)
	}
}

func TestOfflineTimeSource_Refresh(t *testing.T) {
	timeSource := NewOfflineTimeSource()

	// Call Refresh and ensure no panic or error occurs
	timeSource.Refresh()
}
