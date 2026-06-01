package config

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
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
	WorkspaceRoot              string
	WebAllowedTools            []string
	BashAllowOutsideWorkspace  bool
	BashAllowDangerousCommands bool
	CookieName                 string
	SessionTTLMinutes          int
}

type LLMClient interface {
	CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

var (
	Client  LLMClient
	ModelID string
)
