package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type Config struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"api_key"`
	ModelID string `json:"model_id"`
}

type fileConfig struct {
	Config
	AppHome          string `json:"app_home"`
	BuiltinSkillsDir string `json:"builtin_skills_dir"`
	CommandBinDir    string `json:"command_bin_dir"`
	CommandScriptDir string `json:"command_script_dir"`
	WorkspaceRoot    string `json:"workspace_root"`
	WebAllowedTools  string `json:"web_allowed_tools"`
}

type AppConfig struct {
	LLM               Config
	ServerAddr        string
	AllowedOrigin     string
	DatabaseURL       string
	RedisAddr         string
	RedisPassword     string
	RedisDB           int
	JWTSecret         string
	AppHome           string
	BuiltinSkillsDir  string
	CommandBinDir     string
	CommandScriptDir  string
	WorkspaceRoot     string
	WebAllowedTools   []string
	CookieName        string
	SessionTTLMinutes int
}

type LLMClient interface {
	CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

var (
	Client  LLMClient
	ModelID string
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

func LoadWebConfig() (AppConfig, error) {
	fileCfg, err := loadConfigFile()
	if err != nil && !os.IsNotExist(err) {
		return AppConfig{}, err
	}

	appHome, err := resolveAppHome(fileCfg)
	if err != nil {
		return AppConfig{}, err
	}

	llm, err := loadLLMConfig()
	if err != nil {
		return AppConfig{}, err
	}

	builtinSkillsDir, err := resolvePath(appHome, getenv("BUILTIN_SKILLS_DIR", firstNonEmpty(fileCfg.BuiltinSkillsDir, filepath.Join("workspaces", "skills"))))
	if err != nil {
		return AppConfig{}, fmt.Errorf("resolve builtin skills dir: %w", err)
	}
	commandBinDir, err := resolvePath(appHome, getenv("COMMAND_BIN_DIR", firstNonEmpty(fileCfg.CommandBinDir, filepath.Join("workspaces", "bin"))))
	if err != nil {
		return AppConfig{}, fmt.Errorf("resolve command bin dir: %w", err)
	}
	commandScriptDir, err := resolvePath(appHome, getenv("COMMAND_SCRIPT_DIR", firstNonEmpty(fileCfg.CommandScriptDir, filepath.Join("workspaces", "cmd"))))
	if err != nil {
		return AppConfig{}, fmt.Errorf("resolve command script dir: %w", err)
	}
	workspaceRoot, err := resolvePath(appHome, getenv("WORKSPACE_ROOT", firstNonEmpty(fileCfg.WorkspaceRoot, "workspace")))
	if err != nil {
		return AppConfig{}, fmt.Errorf("resolve workspace root: %w", err)
	}
	dbURL := getenv("DATABASE_URL", buildDefaultMySQLDSN())
	redisDB, err := getenvInt("REDIS_DB", 0)
	if err != nil {
		return AppConfig{}, fmt.Errorf("invalid REDIS_DB: %w", err)
	}

	return AppConfig{
		LLM:               llm,
		ServerAddr:        getenv("SERVER_ADDR", ":8080"),
		AllowedOrigin:     getenv("ALLOWED_ORIGIN", "http://localhost:5173"),
		DatabaseURL:       dbURL,
		RedisAddr:         getenv("REDIS_ADDR", "1.12.217.28:6379"),
		RedisPassword:     getenv("REDIS_PASSWORD", "213140"),
		RedisDB:           redisDB,
		JWTSecret:         getenv("JWT_SECRET", "nano-cc-local-secret"),
		AppHome:           appHome,
		BuiltinSkillsDir:  builtinSkillsDir,
		CommandBinDir:     commandBinDir,
		CommandScriptDir:  commandScriptDir,
		WorkspaceRoot:     workspaceRoot,
		WebAllowedTools:   parseCSVList(getenv("WEB_ALLOWED_TOOLS", firstNonEmpty(fileCfg.WebAllowedTools, "load_skill"))),
		CookieName:        getenv("SESSION_COOKIE_NAME", "nano_cc_session"),
		SessionTTLMinutes: getenvIntOrDefault("SESSION_TTL_MINUTES", 60*24*7),
	}, nil
}

func EnsureAppLayout(cfg AppConfig) error {
	paths := []struct {
		label string
		path  string
	}{
		{label: "app home", path: cfg.AppHome},
		{label: "logs", path: filepath.Join(cfg.AppHome, "logs")},
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

func buildDefaultMySQLDSN() string {
	host := getenv("MYSQL_HOST", "1.12.217.28")
	port := getenv("MYSQL_PORT", "3306")
	user := getenv("MYSQL_USER", getenv("DB_USER", "root"))
	password := getenv("MYSQL_PASSWORD", getenv("DB_PASSWORD", "213140"))
	database := getenv("MYSQL_DATABASE", getenv("DB_NAME", "vibe_coding"))
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&multiStatements=true&loc=Local", user, password, host, port, database)
}

func getenv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func getenvInt(key string, fallback int) (int, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

func getenvIntOrDefault(key string, fallback int) int {
	v, err := getenvInt(key, fallback)
	if err != nil {
		return fallback
	}
	return v
}

func parseCSVList(value string) []string {
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item == "" {
			continue
		}
		if _, exists := seen[item]; exists {
			continue
		}
		seen[item] = struct{}{}
		items = append(items, item)
	}
	return items
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func resolveAppHome(fileCfg fileConfig) (string, error) {
	appHome := getenv("APP_HOME", firstNonEmpty(fileCfg.AppHome, "."))
	resolved, err := filepath.Abs(appHome)
	if err != nil {
		return "", fmt.Errorf("resolve APP_HOME: %w", err)
	}
	return filepath.Clean(resolved), nil
}

func configFilePath() string {
	if appHome := strings.TrimSpace(os.Getenv("APP_HOME")); appHome != "" {
		return filepath.Join(appHome, "config.json")
	}
	return "config.json"
}

func resolvePath(appHome, pathValue string) (string, error) {
	if strings.TrimSpace(pathValue) == "" {
		pathValue = "."
	}
	if !filepath.IsAbs(pathValue) {
		pathValue = filepath.Join(appHome, pathValue)
	}
	resolved, err := filepath.Abs(pathValue)
	if err != nil {
		return "", err
	}
	return filepath.Clean(resolved), nil
}
