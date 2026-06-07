package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"api_key"`
	ModelID string `json:"model_id"`
}

type fileConfig struct {
	Config
	AppHome                    string `json:"app_home"`
	BuiltinSkillsDir           string `json:"builtin_skills_dir"`
	CommandBinDir              string `json:"command_bin_dir"`
	CommandScriptDir           string `json:"command_script_dir"`
	SystemPromptPath           string `json:"system_prompt_path"`
	WorkspaceRoot              string `json:"workspace_root"`
	WebAllowedTools            string `json:"web_allowed_tools"`
	BashAllowOutsideWorkspace  bool   `json:"bash_allow_outside_workspace"`
	BashAllowDangerousCommands bool   `json:"bash_allow_dangerous_commands"`
}

type AppConfig struct {
	LLM                        Config
	ServerAddr                 string
	AllowedOrigin              string
	DatabaseURL                string
	RedisAddr                  string
	RedisPassword              string
	RedisDB                    int
	JWTSecret                  string
	AppHome                    string
	BuiltinSkillsDir           string
	CommandBinDir              string
	CommandScriptDir           string
	SystemPromptPath           string
	WorkspaceRoot              string
	LogsDir                    string
	WebAllowedTools            []string
	BashAllowOutsideWorkspace  bool
	BashAllowDangerousCommands bool
	CookieName                 string
	SessionTTLMinutes          int
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
	if appHome := strings.TrimSpace(os.Getenv("APP_HOME")); appHome != "" {
		return filepath.Join(appHome, "config.json")
	}
	return "config.json"
}

func loadLLMConfig() (Config, error) {
	cfg := Config{
		BaseURL: strings.TrimSpace(os.Getenv("OPENAI_BASE_URL")),
		APIKey:  strings.TrimSpace(os.Getenv("OPENAI_API_KEY")),
		ModelID: strings.TrimSpace(os.Getenv("MODEL_ID")),
	}

	fileCfg, err := loadConfigFile()
	if err != nil && !os.IsNotExist(err) {
		return Config{}, err
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = fileCfg.Config.BaseURL
	}
	if cfg.APIKey == "" {
		cfg.APIKey = fileCfg.Config.APIKey
	}
	if cfg.ModelID == "" {
		cfg.ModelID = fileCfg.Config.ModelID
	}

	if cfg.BaseURL == "" || cfg.APIKey == "" || cfg.ModelID == "" {
		return Config{}, fmt.Errorf("missing LLM config; set OPENAI_BASE_URL, OPENAI_API_KEY, MODEL_ID or provide WORKSPACE_ROOT/config.json")
	}

	return cfg, nil
}
