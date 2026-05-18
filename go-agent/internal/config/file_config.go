package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func loadConfigFile() (fileConfig, error) {
	data, err := os.ReadFile(configFilePath())
	if err != nil {
		return fileConfig{}, err
	}

	var cfg fileConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fileConfig{}, fmt.Errorf("failed to parse config.json: %w", err)
	}
	return cfg, nil
}

func configFilePath() string {
	if appHome := strings.TrimSpace(os.Getenv("APP_HOME")); appHome != "" {
		return filepath.Join(appHome, "config.json")
	}
	return "config.json"
}
