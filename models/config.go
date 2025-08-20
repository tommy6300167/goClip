package models

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	DefaultMaxHistoryItems   = 100
	DefaultMaxDisplayLength  = 30
	DefaultPollingInterval   = 800
	DefaultUndoCountdownSecs = 3
	DefaultWindowWidth       = 820
	DefaultWindowHeight      = 520
	DefaultAuthApiEndpoint   = "https://api.superclip.com/auth/validate"
)

type AppConfig struct {
	MaxHistoryItems   int
	MaxDisplayLength  int
	PollingInterval   int
	UndoCountdownSecs int
	WindowWidth       int
	WindowHeight      int
	LogDirPath        string
	LogFilePath       string
	ImageDirPath      string
	ConfigFilePath    string
	AuthFilePath      string
	AuthApiEndpoint   string
}

func NewAppConfig() *AppConfig {
	home, _ := os.UserHomeDir()
	logDir := filepath.Join(home, "Library", "Application Support", "ClipMini")

	config := &AppConfig{
		MaxHistoryItems:   DefaultMaxHistoryItems,
		MaxDisplayLength:  DefaultMaxDisplayLength,
		PollingInterval:   DefaultPollingInterval,
		UndoCountdownSecs: DefaultUndoCountdownSecs,
		WindowWidth:       DefaultWindowWidth,
		WindowHeight:      DefaultWindowHeight,
		LogDirPath:        logDir,
		LogFilePath:       filepath.Join(logDir, "history.txt"),
		ImageDirPath:      filepath.Join(logDir, "images"),
		ConfigFilePath:    filepath.Join(logDir, "config.json"),
		AuthFilePath:      filepath.Join(logDir, "auth.json"),
		AuthApiEndpoint:   DefaultAuthApiEndpoint,
	}
	config.LoadConfig()

	return config
}

func (c *AppConfig) SaveConfig() error {
	if err := os.MkdirAll(c.LogDirPath, 0755); err != nil {
		return err
	}

	configData := map[string]interface{}{
		"MaxHistoryItems":   c.MaxHistoryItems,
		"MaxDisplayLength":  c.MaxDisplayLength,
		"PollingInterval":   c.PollingInterval,
		"UndoCountdownSecs": c.UndoCountdownSecs,
		"WindowWidth":       c.WindowWidth,
		"WindowHeight":      c.WindowHeight,
	}

	data, err := json.MarshalIndent(configData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.ConfigFilePath, data, 0644)
}

func (c *AppConfig) LoadConfig() error {

	if _, err := os.Stat(c.ConfigFilePath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(c.ConfigFilePath)
	if err != nil {
		return err
	}

	var configData map[string]interface{}
	if err := json.Unmarshal(data, &configData); err != nil {
		return err
	}

	if val, ok := configData["MaxHistoryItems"].(float64); ok {
		c.MaxHistoryItems = int(val)
	}
	if val, ok := configData["MaxDisplayLength"].(float64); ok {
		c.MaxDisplayLength = int(val)
	}
	if val, ok := configData["PollingInterval"].(float64); ok {
		c.PollingInterval = int(val)
	}
	if val, ok := configData["UndoCountdownSecs"].(float64); ok {
		c.UndoCountdownSecs = int(val)
	}
	if val, ok := configData["WindowWidth"].(float64); ok {
		c.WindowWidth = int(val)
	}
	if val, ok := configData["WindowHeight"].(float64); ok {
		c.WindowHeight = int(val)
	}

	return nil
}
