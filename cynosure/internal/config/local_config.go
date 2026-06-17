package config

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const defaultLocalAllowedTools = "load_skill,bash,read_file,write_file,edit_file,multi_edit,grep,glob,ls,web_fetch,todo_write,todo_list,spawn_subagent"

type linkSettings struct {
	Env linkSettingsEnv `json:"env"`
}

type linkSettingsEnv struct {
	OpenAuthToken string `json:"open_auth_token"`
	OpenModel     string `json:"open_model"`
	OpenBaseURL   string `json:"open_base_url"`
}

type CynosureMarkdownContext struct {
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
	llm, err := loadLocalLLMConfig(fileCfg)
	if err != nil {
		return AppConfig{}, err
	}
	appHome, err := cynosureHomeDir()
	if err != nil {
		return AppConfig{}, err
	}
	systemPromptPath, err := userSystemPromptOverridePath()
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
		SystemPromptPath:            systemPromptPath,
		WorkspaceRoot:               workspaceRoot,
		LogsDir:                     "",
		AllowedTools:                parseCSVList(allowedTools),
		ConversationLockTTL:         time.Duration(intOrDefault(fileCfg.ConversationLockTTLSeconds, 30)) * time.Second,
		MemoryWorkTimeout:           time.Duration(intOrDefault(fileCfg.MemoryWorkTimeoutSeconds, 110)) * time.Second,
		ConversationLockWaitTimeout: time.Duration(intOrDefault(fileCfg.ConversationLockWaitTimeoutSeconds, 130)) * time.Second,
	}, nil
}

func cynosureHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	return filepath.Join(home, ".cynosure"), nil
}

// userSystemPromptOverridePath 返回用户可选覆盖文件 ~/.cynosure/system_prompt.md；
// 文件不存在时返回空字符串，表示使用内置（embedded）system prompt。
func userSystemPromptOverridePath() (string, error) {
	home, err := cynosureHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(home, "system_prompt.md")
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("stat system prompt override %s: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return "", nil
	}
	return path, nil
}

func loadLocalLLMConfig(_ fileConfig) (Config, error) {
	return loadLinkLLMConfig()
}

func CynosureSettingsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	return filepath.Join(home, ".cynosure", "settings.json"), nil
}

func CynosureSkillsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	return filepath.Join(home, ".cynosure", "skills"), nil
}

func CynosureMemoryDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	return filepath.Join(home, ".cynosure", "memory"), nil
}

func CynosureTaskOutputsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	return filepath.Join(home, ".cynosure", "task_outputs"), nil
}

func CynosureSessionDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	return filepath.Join(home, ".cynosure", "session"), nil
}

func CynosureMarkdownPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	return filepath.Join(home, ".cynosure", "CYNOSURE.MD"), nil
}

func WorkspaceCynosureSkillsDir(workspaceRoot string) string {
	return filepath.Join(strings.TrimSpace(workspaceRoot), ".cynosure", "skills")
}

func WorkspaceCynosureMarkdownPath(workspaceRoot string) string {
	return filepath.Join(strings.TrimSpace(workspaceRoot), ".cynosure", "CYNOSURE.MD")
}

func WorkspaceMCPConfigPath(workspaceRoot string) string {
	return filepath.Join(strings.TrimSpace(workspaceRoot), ".cynosure", ".mcp.json")
}

// WorkspaceKey 由工作区目录名与其绝对路径的 sha256 前 8 位组成，用于在
// ~/.cynosure/ 下隔离不同项目的运行期数据（记忆、日志等）。
func WorkspaceKey(workspaceRoot string) string {
	base := sanitizePathSegment(filepath.Base(filepath.Clean(workspaceRoot)))
	if base == "" {
		base = "project"
	}
	abs, err := filepath.Abs(workspaceRoot)
	if err != nil {
		abs = filepath.Clean(workspaceRoot)
	}
	sum := sha256.Sum256([]byte(abs))
	return base + "-" + hex.EncodeToString(sum[:])[:8]
}

// CynosureSessionLogsDir 返回按工作区与会话隔离的日志目录：
// ~/.cynosure/<workspace-key>/<session_id>/logs
func CynosureSessionLogsDir(workspaceRoot, sessionID string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	session := sanitizePathSegment(sessionID)
	if session == "" {
		session = "session"
	}
	return filepath.Join(home, ".cynosure", WorkspaceKey(workspaceRoot), session, "logs"), nil
}

// sanitizePathSegment 把任意字符串收敛为安全的单层目录名。
func sanitizePathSegment(s string) string {
	var b strings.Builder
	lastDash := false
	for _, r := range strings.TrimSpace(s) {
		allowed := r == '.' || r == '_' || r == '-' ||
			(r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
		if allowed {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-._")
}

func loadLinkLLMConfig() (Config, error) {
	path, err := CynosureSettingsPath()
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

func LoadCynosureMarkdownContext(workspaceRoot string) (CynosureMarkdownContext, error) {
	userPath, err := CynosureMarkdownPath()
	if err != nil {
		return CynosureMarkdownContext{}, err
	}
	workspacePath := WorkspaceCynosureMarkdownPath(workspaceRoot)
	userContent, err := readOptionalMarkdown(userPath)
	if err != nil {
		return CynosureMarkdownContext{}, err
	}
	workspaceContent, err := readOptionalMarkdown(workspacePath)
	if err != nil {
		return CynosureMarkdownContext{}, err
	}
	return CynosureMarkdownContext{UserPath: userPath, UserContent: userContent, WorkspacePath: workspacePath, WorkspaceContent: workspaceContent}, nil
}

func readOptionalMarkdown(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("read CYNOSURE.MD %s: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("read CYNOSURE.MD %s: not a regular file", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read CYNOSURE.MD %s: %w", path, err)
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
