package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const defaultLocalAllowedTools = "load_skill,bash,read_file,write_file,edit_file,todo_write,spawn_subagent"

type linkSettings struct {
	Env linkSettingsEnv `json:"env"`
}

type linkSettingsEnv struct {
	OpenAuthToken string `json:"open_auth_token"`
	OpenModel     string `json:"open_model"`
	OpenBaseURL   string `json:"open_base_url"`
}

type LinkMarkdownContext struct {
	UserPath         string
	UserContent      string
	WorkspacePath    string
	WorkspaceContent string
}

func LoadLocalConfig(cwd string) (AppConfig, error) {
	fileCfg, err := loadConfigFile()
	if err != nil {
		return AppConfig{}, err
	}
	appHome, err := resolveAppHome(fileCfg)
	if err != nil {
		return AppConfig{}, err
	}
	llm, err := loadLocalLLMConfig(fileCfg)
	if err != nil {
		return AppConfig{}, err
	}
	runtimeDirs, err := resolveRuntimePaths(appHome, fileCfg)
	if err != nil {
		return AppConfig{}, err
	}
	workspaceRoot, err := resolveLocalWorkspaceRoot(cwd)
	if err != nil {
		return AppConfig{}, err
	}
	allowedTools := firstNonEmpty(fileCfg.AllowedTools, defaultLocalAllowedTools)
	return AppConfig{
		LLM:                         llm,
		AppHome:                     appHome,
		BuiltinSkillsDir:            runtimeDirs.builtinSkillsDir,
		CommandBinDir:               runtimeDirs.commandBinDir,
		CommandScriptDir:            runtimeDirs.commandScriptDir,
		SystemPromptPath:            runtimeDirs.systemPromptPath,
		WorkspaceRoot:               workspaceRoot,
		LogsDir:                     runtimeDirs.logsDir,
		AllowedTools:                parseCSVList(allowedTools),
		BashAllowOutsideWorkspace:   fileCfg.BashAllowOutsideWorkspace,
		BashAllowDangerousCommands:  fileCfg.BashAllowDangerousCommands,
		ConversationLockTTL:         time.Duration(intOrDefault(fileCfg.ConversationLockTTLSeconds, 30)) * time.Second,
		MemoryWorkTimeout:           time.Duration(intOrDefault(fileCfg.MemoryWorkTimeoutSeconds, 110)) * time.Second,
		ConversationLockWaitTimeout: time.Duration(intOrDefault(fileCfg.ConversationLockWaitTimeoutSeconds, 130)) * time.Second,
	}, nil
}

func loadLocalLLMConfig(_ fileConfig) (Config, error) {
	return loadLinkLLMConfig()
}

func LinkSettingsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	return filepath.Join(home, ".link", "settings.json"), nil
}

func LinkSkillsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	return filepath.Join(home, ".link", "skills"), nil
}

func LinkSessionDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	return filepath.Join(home, "link", "session"), nil
}

func LinkMarkdownPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	return filepath.Join(home, ".link", "LINK.MD"), nil
}

func WorkspaceLinkSkillsDir(workspaceRoot string) string {
	return filepath.Join(strings.TrimSpace(workspaceRoot), ".link", "skills")
}

func WorkspaceLinkMarkdownPath(workspaceRoot string) string {
	return filepath.Join(strings.TrimSpace(workspaceRoot), ".link", "LINK.MD")
}

func WorkspaceMCPConfigPath(workspaceRoot string) string {
	return filepath.Join(strings.TrimSpace(workspaceRoot), ".link", ".mcp.json")
}

func loadLinkLLMConfig() (Config, error) {
	path, err := LinkSettingsPath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, fmt.Errorf("missing LLM settings; create %s with env.open_auth_token, env.open_model, env.open_base_url", path)
		}
		return Config{}, fmt.Errorf("read LLM settings %s: %w", path, err)
	}
	var settings linkSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return Config{}, fmt.Errorf("parse LLM settings %s: %w", path, err)
	}
	cfg := Config{
		BaseURL: strings.TrimSpace(settings.Env.OpenBaseURL),
		APIKey:  strings.TrimSpace(settings.Env.OpenAuthToken),
		ModelID: strings.TrimSpace(settings.Env.OpenModel),
	}
	missing := make([]string, 0, 3)
	if cfg.APIKey == "" {
		missing = append(missing, "env.open_auth_token")
	}
	if cfg.ModelID == "" {
		missing = append(missing, "env.open_model")
	}
	if cfg.BaseURL == "" {
		missing = append(missing, "env.open_base_url")
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return Config{}, fmt.Errorf("missing LLM settings in %s: %s", path, strings.Join(missing, ", "))
	}
	return cfg, nil
}

func LoadLinkMarkdownContext(workspaceRoot string) (LinkMarkdownContext, error) {
	userPath, err := LinkMarkdownPath()
	if err != nil {
		return LinkMarkdownContext{}, err
	}
	workspacePath := WorkspaceLinkMarkdownPath(workspaceRoot)
	userContent, err := readOptionalMarkdown(userPath)
	if err != nil {
		return LinkMarkdownContext{}, err
	}
	workspaceContent, err := readOptionalMarkdown(workspacePath)
	if err != nil {
		return LinkMarkdownContext{}, err
	}
	return LinkMarkdownContext{UserPath: userPath, UserContent: userContent, WorkspacePath: workspacePath, WorkspaceContent: workspaceContent}, nil
}

func readOptionalMarkdown(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("read LINK.MD %s: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("read LINK.MD %s: not a regular file", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read LINK.MD %s: %w", path, err)
	}
	return strings.TrimSpace(string(data)), nil
}

func resolveLocalWorkspaceRoot(cwd string) (string, error) {
	if strings.TrimSpace(cwd) == "" {
		cwd = "."
	}
	resolved, err := filepath.Abs(cwd)
	if err != nil {
		return "", fmt.Errorf("resolve cwd: %w", err)
	}
	clean := filepath.Clean(resolved)
	info, err := os.Stat(clean)
	if err != nil {
		return "", fmt.Errorf("stat cwd: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("cwd is not a directory: %s", clean)
	}
	return clean, nil
}
