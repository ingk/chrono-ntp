package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

const defaultNtpServer = "time.google.com"
const defaultTimeFormat = "ISO8601"
const defaultTimeZone = "Local"

type Configuration struct {
	Server        string `toml:"server"`
	TimeZone      string `toml:"time-zone"`
	HideStatusbar bool   `toml:"hide-statusbar"`
	HideDate      bool   `toml:"hide-date"`
	ShowTimeZone  bool   `toml:"show-time-zone"`
	TimeFormat    string `toml:"time-format"`
	Beeps         bool   `toml:"beeps"`
	Offline       bool   `toml:"offline"`
}

func getConfigurationContents(path string) []byte {
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read config file %s: %v", path, err)
		}
		return data
	}

	return nil
}

func parseConfiguration(data []byte) Configuration {
	config := Configuration{
		Server:        defaultNtpServer,
		TimeZone:      defaultTimeZone,
		HideStatusbar: false,
		HideDate:      false,
		ShowTimeZone:  true,
		TimeFormat:    defaultTimeFormat,
		Beeps:         false,
		Offline:       false,
	}

	if err := toml.Unmarshal(data, &config); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	return config
}

func LoadConfiguration() Configuration {
	configPath := filepath.Join(os.Getenv("HOME"), ".chrono-ntp.toml")
	data := getConfigurationContents(configPath)
	return parseConfiguration(data)
}
