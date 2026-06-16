package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"-"`
	ModelID string `json:"model_id"`
}

type fileConfig struct {
	Config
	AppHome                            string `json:"app_home"`
	SystemPromptPath                   string `json:"system_prompt_path"`
	WorkspaceRoot                      string `json:"workspace_root"`
	AllowedTools                       string `json:"allowed_tools"`
	ConversationLockTTLSeconds         int    `json:"conversation_lock_ttl_seconds"`
	MemoryWorkTimeoutSeconds           int    `json:"memory_work_timeout_seconds"`
	ConversationLockWaitTimeoutSeconds int    `json:"conversation_lock_wait_timeout_seconds"`
}

type AppConfig struct {
	LLM                         Config
	AppHome                     string
	SystemPromptPath            string
	WorkspaceRoot               string
	LogsDir                     string
	AllowedTools                []string
	ConversationLockTTL         time.Duration
	MemoryWorkTimeout           time.Duration
	ConversationLockWaitTimeout time.Duration
}

func defaultFileConfig() fileConfig {
	return fileConfig{
		AppHome:          ".",
		SystemPromptPath: "system_prompt.md",
		WorkspaceRoot:    "workspace",
		AllowedTools:     defaultLocalAllowedTools,
	}
}

// loadConfigFile 返回运行配置。默认值内置在代码中，使二进制可在任意目录运行；
// 若存在 ~/.cynosure/config.json，则用其内容覆盖默认值（可选）。
func loadConfigFile() (fileConfig, error) {
	cfg := defaultFileConfig()
	path, err := configFilePath()
	if err != nil {
		return fileConfig{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return fileConfig{}, fmt.Errorf("read config %s: %w", path, err)
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fileConfig{}, fmt.Errorf("failed to parse config %s: %w", path, err)
	}
	return cfg, nil
}

func configFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	return filepath.Join(home, ".cynosure", "config.json"), nil
}

func loadLLMConfig(fileCfg fileConfig) (Config, error) {
	cfg := Config{
		BaseURL: strings.TrimSpace(fileCfg.Config.BaseURL),
		APIKey:  strings.TrimSpace(getenv("OPENAI_API_KEY")),
		ModelID: strings.TrimSpace(fileCfg.Config.ModelID),
	}

	if cfg.BaseURL == "" || cfg.ModelID == "" {
		return Config{}, fmt.Errorf("missing LLM config; set base_url, model_id in config.json")
	}
	if cfg.APIKey == "" {
		return Config{}, fmt.Errorf("missing LLM api key; set OPENAI_API_KEY environment variable")
	}

	return cfg, nil
}
