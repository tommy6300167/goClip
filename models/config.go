package models

import (
	"os"
	"path/filepath"
)

const (
	DefaultMaxHistoryItems  = 30
	DefaultMaxDisplayLength = 50
	DefaultPollingInterval  = 800 // milliseconds
)

type AppConfig struct {
	MaxHistoryItems  int
	MaxDisplayLength int
	PollingInterval  int
	LogDirPath       string
	LogFilePath      string
	ImageDirPath     string
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
	}
}