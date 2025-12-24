package configuration

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseConfiguration_Defaults(t *testing.T) {
	config, _ := parseConfiguration(nil)

	if config.Server != "time.google.com" {
		t.Errorf("expected Server %q, got %q", "time.google.com", config.Server)
	}
	if config.TimeZone != "Local" {
		t.Errorf("expected TimeZone %q, got %q", "Local", config.TimeZone)
	}
	if config.HideStatusBar != false {
		t.Errorf("expected HideStatusBar false, got %v", config.HideStatusBar)
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
	if config.Offline != false {
		t.Errorf("expected Offline false, got %v", config.Offline)
	}
}

func TestParseConfiguration_Content(t *testing.T) {
	tomlContent := `
server = "pool.example-time-server.org"
time-zone = "Europe/Berlin"
hide-status-bar = true
hide-date = true
show-time-zone = true
time-format = "12h_AM_PM"
beeps = true
offline = true
`
	config, _ := parseConfiguration([]byte(tomlContent))

	if config.Server != "pool.example-time-server.org" {
		t.Errorf("expected Server 'pool.example-time-server.org', got %q", config.Server)
	}
	if config.TimeZone != "Europe/Berlin" {
		t.Errorf("expected TimeZone 'Europe/Berlin', got %q", config.TimeZone)
	}
	if config.HideStatusBar != true {
		t.Errorf("expected HideStatusBar true, got %v", config.HideStatusBar)
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
	if config.Offline != true {
		t.Errorf("expected Offline true, got %v", config.Offline)
	}
}

func TestLoadConfiguration(t *testing.T) {
	// Create a temporary directory to act as HOME
	tempDir, err := os.MkdirTemp("", "chrono-ntp-test-home")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set HOME to the temp directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	// Write a config file in the temp HOME
	configPath := filepath.Join(tempDir, ".chrono-ntp.toml")
	tomlContent := `
server = "mocked.server"
time-zone = "UTC"
`
	if err := os.WriteFile(configPath, []byte(tomlContent), 0644); err != nil {
		t.Fatalf("Failed to write mock config: %v", err)
	}

	config, err := LoadConfiguration()

	if err != nil {
		t.Errorf("unexpected error', got %v", err)
	}
	if config.Server != "mocked.server" {
		t.Errorf("expected Server 'mocked.server', got %q", config.Server)
	}
	if config.TimeZone != "UTC" {
		t.Errorf("expected TimeZone 'UTC', got %q", config.TimeZone)
	}
}

func TestLoadConfiguration_Error(t *testing.T) {
	// Create a temporary directory to act as HOME
	tempDir, err := os.MkdirTemp("", "chrono-ntp-test-home")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set HOME to the temp directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	// Write a config file in the temp HOME
	configPath := filepath.Join(tempDir, ".chrono-ntp.toml")
	tomlContent := `invalid-toml`
	if err := os.WriteFile(configPath, []byte(tomlContent), 0644); err != nil {
		t.Fatalf("Failed to write mock config: %v", err)
	}

	config, err := LoadConfiguration()

	if err == nil {
		t.Fatalf("expected error, got %v", config)
	}
}

func TestLoadConfiguration_MissingFile(t *testing.T) {
	// Create a temporary directory to act as HOME
	tempDir, err := os.MkdirTemp("", "chrono-ntp-test-home")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set HOME to the temp directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	_, configErr := LoadConfiguration()

	if configErr != nil {
		t.Errorf("did not expect error, got %v", configErr)
	}
}

func TestWriteConfiguration(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "chrono-ntp-test-home")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	configPath := filepath.Join(tempDir, ".chrono-ntp.toml")
	config := Configuration{
		Server:        "write.test.server",
		TimeZone:      "Mars/Colony",
		HideStatusBar: true,
		HideDate:      true,
		ShowTimeZone:  false,
		TimeFormat:    "mars",
		Beeps:         true,
		Offline:       true,
	}

	configPathResult, err := WriteConfiguration(config)

	if configPathResult != configPath {
		t.Fatalf("expected config path %q, got %q", configPath, configPathResult)
	}

	if err != nil {
		t.Fatalf("Failed to write configuration: %v", err)
	}

	loadedConfig, err := LoadConfiguration()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	if !reflect.DeepEqual(config, loadedConfig) {
		t.Fatalf("Written configuration does not equal configuration: %v %v", config, loadedConfig)
	}
}

func TestWriteConfiguration_WriteError(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "chrono-ntp-test-home")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	configPath := filepath.Join(tempDir, ".chrono-ntp.toml")
	if err := os.Mkdir(configPath, 0755); err != nil {
		t.Fatalf("Failed to create directory to block config file: %v", err)
	}

	config := Configuration{}
	configPathResult, err := WriteConfiguration(config)

	if configPathResult != configPath {
		t.Fatalf("expected config path %q, got %q", configPath, configPathResult)
	}

	if err == nil {
		t.Fatalf("expected error, got %v", err)
	}
}
