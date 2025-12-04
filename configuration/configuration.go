package configuration

import (
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

func getConfigurationContents(path string) ([]byte, error) {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return nil, nil
	} else if err == nil {
		return os.ReadFile(path)
	}
	return nil, err
}

func parseConfiguration(data []byte) (Configuration, error) {
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

	err := toml.Unmarshal(data, &config)
	if err != nil {
		return Configuration{}, err
	}
	return config, nil
}

func LoadConfiguration() (Configuration, error) {
	configPath := filepath.Join(os.Getenv("HOME"), ".chrono-ntp.toml")
	data, err := getConfigurationContents(configPath)
	if err != nil {
		return Configuration{}, err
	}
	return parseConfiguration(data)
}

func WriteConfiguration(config Configuration) (string, error) {
	configPath := filepath.Join(os.Getenv("HOME"), ".chrono-ntp.toml")
	data, err := toml.Marshal(config)
	if err != nil {
		return configPath, err
	}
	return configPath, os.WriteFile(configPath, data, 0644)
}
