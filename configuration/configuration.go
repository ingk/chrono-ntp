package configuration

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

var (
    userConfigDir = os.UserConfigDir
    tomlMarshal  = toml.Marshal
)

const (
	defaultNtpServer      = "time.google.com"
	defaultTimeFormat     = "ISO8601"
	defaultTimeZone       = "Local"
	defaultConfigFileName = ".chrono-ntp.toml"
	appConfigDirName      = "chrono-ntp"
	appConfigFileName     = "chrono-ntp.toml"
)

type Configuration struct {
	Server        string `toml:"server"`
	TimeZone      string `toml:"time-zone"`
	HideStatusBar bool   `toml:"hide-status-bar"`
	HideDate      bool   `toml:"hide-date"`
	ShowTimeZone  bool   `toml:"show-time-zone"`
	TimeFormat    string `toml:"time-format"`
	BeepPattern   string `toml:"beep-pattern"`
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
		HideStatusBar: false,
		HideDate:      false,
		ShowTimeZone:  true,
		TimeFormat:    defaultTimeFormat,
		BeepPattern:   "",
		Offline:       false,
	}

	err := toml.Unmarshal(data, &config)
	if err != nil {
		return Configuration{}, err
	}
	return config, nil
}

func canonicalizeDir(path string) (string, error) {
	if path == "" {
		return "", errors.New("empty path")
	}

	cleanPath := filepath.Clean(path)
	if !filepath.IsAbs(cleanPath) {
		return "", errors.New("path is not absolute")
	}
	return cleanPath, nil
}

func configurationSearchPaths() []string {
	paths := make([]string, 0, 2)

	if configDir, err := userConfigDir(); err == nil && configDir != "" {
		if canonicalConfigDir, err := canonicalizeDir(configDir); err == nil {
			paths = append(paths, filepath.Join(canonicalConfigDir, appConfigDirName, appConfigFileName))
		}
	}

	if home := os.Getenv("HOME"); home != "" {
		if canonicalHome, err := canonicalizeDir(home); err == nil {
			paths = append(paths, filepath.Join(canonicalHome, defaultConfigFileName))
		}
	}

	return paths
}

func configWritePath() (string, error) {
	if configDir, err := userConfigDir(); err == nil && configDir != "" {
		if canonicalConfigDir, err := canonicalizeDir(configDir); err == nil {
			appDir := filepath.Join(canonicalConfigDir, appConfigDirName)
			if err := os.MkdirAll(appDir, 0700); err != nil {
				return "", err
			}
			return filepath.Join(appDir, appConfigFileName), nil
		}
	}

	if home := os.Getenv("HOME"); home != "" {
		if canonicalHome, err := canonicalizeDir(home); err == nil {
			configPath := filepath.Join(canonicalHome, defaultConfigFileName)
			if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
				return "", err
			}
			return configPath, nil
		}
	}

	return "", errors.New("configuration directory unavailable")
}

func LoadConfiguration() (Configuration, error) {
	for _, configPath := range configurationSearchPaths() {
		data, err := getConfigurationContents(configPath)
		if err != nil {
			return Configuration{}, err
		}
		if data != nil {
			return parseConfiguration(data)
		}
	}

	return parseConfiguration(nil)
}

func WriteConfiguration(config Configuration) (string, error) {
	configPath, err := configWritePath()
	if err != nil {
		return "", err
	}

	data, err := tomlMarshal(config)
	if err != nil {
		return configPath, err
	}
	return configPath, os.WriteFile(configPath, data, 0600)
}
