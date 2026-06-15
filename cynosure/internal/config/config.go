package config

import (
	"encoding/json"
	"fmt"
	"os"
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
	BuiltinSkillsDir                   string `json:"builtin_skills_dir"`
	CommandBinDir                      string `json:"command_bin_dir"`
	CommandScriptDir                   string `json:"command_script_dir"`
	SystemPromptPath                   string `json:"system_prompt_path"`
	WorkspaceRoot                      string `json:"workspace_root"`
	AllowedTools                       string `json:"allowed_tools"`
	BashAllowOutsideWorkspace          bool   `json:"bash_allow_outside_workspace"`
	BashAllowDangerousCommands         bool   `json:"bash_allow_dangerous_commands"`
	ConversationLockTTLSeconds         int    `json:"conversation_lock_ttl_seconds"`
	MemoryWorkTimeoutSeconds           int    `json:"memory_work_timeout_seconds"`
	ConversationLockWaitTimeoutSeconds int    `json:"conversation_lock_wait_timeout_seconds"`
}

type AppConfig struct {
	LLM                         Config
	AppHome                     string
	BuiltinSkillsDir            string
	CommandBinDir               string
	CommandScriptDir            string
	SystemPromptPath            string
	WorkspaceRoot               string
	LogsDir                     string
	AllowedTools                []string
	BashAllowOutsideWorkspace   bool
	BashAllowDangerousCommands  bool
	ConversationLockTTL         time.Duration
	MemoryWorkTimeout           time.Duration
	ConversationLockWaitTimeout time.Duration
}

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
	return "config.json"
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
