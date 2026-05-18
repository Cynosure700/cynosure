package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func EnsureAppLayout(cfg AppConfig) error {
	paths := []struct {
		label string
		path  string
	}{
		{label: "app home", path: cfg.AppHome},
		{label: "logs", path: filepath.Join(cfg.WorkspaceRoot, "logs")},
		{label: "builtin skills dir", path: cfg.BuiltinSkillsDir},
		{label: "command bin dir", path: cfg.CommandBinDir},
		{label: "command script dir", path: cfg.CommandScriptDir},
		{label: "workspace root", path: cfg.WorkspaceRoot},
	}
	for _, item := range paths {
		if strings.TrimSpace(item.path) == "" {
			continue
		}
		if err := os.MkdirAll(item.path, 0o755); err != nil {
			return fmt.Errorf("create %s: %w", item.label, err)
		}
	}
	return nil
}

func ValidateAppLayout(cfg AppConfig) error {
	paths := []struct {
		label string
		path  string
	}{
		{label: "app home", path: cfg.AppHome},
		{label: "builtin skills dir", path: cfg.BuiltinSkillsDir},
		{label: "command bin dir", path: cfg.CommandBinDir},
		{label: "command script dir", path: cfg.CommandScriptDir},
		{label: "workspace root", path: cfg.WorkspaceRoot},
	}
	for _, item := range paths {
		if strings.TrimSpace(item.path) == "" {
			return fmt.Errorf("%s is required", item.label)
		}
		info, err := os.Stat(item.path)
		if err != nil {
			return fmt.Errorf("stat %s: %w", item.label, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("%s is not a directory", item.label)
		}
	}
	return nil
}
