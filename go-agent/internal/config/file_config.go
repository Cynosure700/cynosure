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
		return fileConfig{}, fmt.Errorf("failed to parse workspace config.json: %w", err)
	}
	return cfg, nil
}

func configFilePath() string {
	if workspaceRoot := strings.TrimSpace(os.Getenv("WORKSPACE_ROOT")); workspaceRoot != "" {
		if appHome := strings.TrimSpace(os.Getenv("APP_HOME")); appHome != "" && !filepath.IsAbs(workspaceRoot) {
			workspaceRoot = filepath.Join(appHome, workspaceRoot)
		}
		return filepath.Join(workspaceRoot, "config.json")
	}
	if appHome := strings.TrimSpace(os.Getenv("APP_HOME")); appHome != "" {
		return filepath.Join(appHome, "workspace", "config.json")
	}
	return filepath.Join("workspace", "config.json")
}
