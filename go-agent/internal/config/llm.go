package config

import (
	"fmt"
	"os"
	"strings"
)

func Init() {
	llm, err := loadLLMConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load LLM config: %v\n", err)
		os.Exit(1)
	}
	InitLLM(llm)
}

func InitLLM(cfg Config) {
	if cfg.BaseURL == "" || cfg.APIKey == "" || cfg.ModelID == "" {
		fmt.Fprintln(os.Stderr, "LLM config requires base_url, api_key, model_id")
		os.Exit(1)
	}

	ModelID = cfg.ModelID
	Client = newDeepseekClient(cfg.BaseURL, cfg.APIKey)
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
		return Config{}, fmt.Errorf("missing LLM config; set OPENAI_BASE_URL, OPENAI_API_KEY, MODEL_ID or provide config.json")
	}

	return cfg, nil
}
