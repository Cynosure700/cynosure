package config

import (
	"fmt"
	"os"
	"strings"
)

// EnsureAppLayout 创建运行期需要的可写目录。所有运行期目录都在 ~/.cynosure/ 下，
// 与二进制位置和启动目录无关，从而支持在任意项目目录运行。
func EnsureAppLayout(cfg AppConfig) error {
	paths := []struct {
		label string
		path  string
	}{
		{label: "app home", path: cfg.AppHome},
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

// ValidateAppLayout 校验运行期关键目录存在且为目录。
func ValidateAppLayout(cfg AppConfig) error {
	paths := []struct {
		label string
		path  string
	}{
		{label: "app home", path: cfg.AppHome},
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
