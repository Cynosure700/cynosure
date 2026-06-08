package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"-"`
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
	ServerAddr                 string `json:"server_addr"`
	AllowedOrigin              string `json:"allowed_origin"`
	MySQLHost                  string `json:"mysql_host"`
	MySQLPort                  string `json:"mysql_port"`
	MySQLUser                  string `json:"mysql_user"`
	MySQLDatabase              string `json:"mysql_database"`
	RedisAddr                  string `json:"redis_addr"`
	RedisDB                    int    `json:"redis_db"`
	ESAddresses                string `json:"es_addresses"`
	ESUsername                 string `json:"es_username"`
	SessionCookieName          string `json:"session_cookie_name"`
	SessionTTLMinutes          int    `json:"session_ttl_minutes"`
}

type AppConfig struct {
	LLM                        Config
	ServerAddr                 string
	AllowedOrigin              string
	DatabaseURL                string
	RedisAddr                  string
	RedisPassword              string
	RedisDB                    int
	ESAddresses                []string
	ESUsername                 string
	ESPassword                 string
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
	return "config.json"
}

func loadLLMConfig(fileCfg fileConfig) (Config, error) {
	cfg := Config{
		BaseURL: strings.TrimSpace(fileCfg.Config.BaseURL),
		APIKey:  firstNonEmpty(getenv("OPENAI_API_KEY"), "sk-06bad5d0b8bc47409c9473aa74615a21"),
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
