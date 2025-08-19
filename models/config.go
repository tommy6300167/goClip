package models

import (
	"os"
	"path/filepath"
)

const (
	DefaultMaxHistoryItems  = 100
	DefaultMaxDisplayLength = 30
	DefaultPollingInterval  = 800
	DefaultAuthApiEndpoint = "https://api.superclip.com/auth/validate"
)

type AppConfig struct {
	MaxHistoryItems  int
	MaxDisplayLength int
	PollingInterval  int
	LogDirPath       string
	LogFilePath      string
	ImageDirPath     string
	ConfigFilePath   string
	AuthFilePath     string
	AuthApiEndpoint  string
}

func NewAppConfig() *AppConfig {
	home, _ := os.UserHomeDir()
	logDir := filepath.Join(home, "Library", "Application Support", "ClipMini")

	return &AppConfig{
		MaxHistoryItems:  DefaultMaxHistoryItems,
		MaxDisplayLength: DefaultMaxDisplayLength,
		PollingInterval:  DefaultPollingInterval,
		LogDirPath:       logDir,
		LogFilePath:      filepath.Join(logDir, "history.txt"),
		ImageDirPath:     filepath.Join(logDir, "images"),
		ConfigFilePath:   filepath.Join(logDir, "config.txt"),
		AuthFilePath:     filepath.Join(logDir, "auth.json"),
		AuthApiEndpoint:  DefaultAuthApiEndpoint,
	}
}
