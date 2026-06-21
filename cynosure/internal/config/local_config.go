package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const defaultLocalAllowedTools = "load_skill,bash,read_file,write_file,edit_file,multi_edit,grep,glob,ls,web_fetch,todo_write,todo_list,spawn_subagent,update_memory,delete_memory"

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

var workspaceRegistryMu sync.Mutex

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
		LLM:                            llm,
		AppHome:                        appHome,
		SystemPromptPath:               systemPromptPath,
		WorkspaceRoot:                  workspaceRoot,
		LogsDir:                        "",
		AllowedTools:                   parseCSVList(allowedTools),
		ConversationLockTTL:            time.Duration(intOrDefault(fileCfg.ConversationLockTTLSeconds, 30)) * time.Second,
		MemoryWorkTimeout:              time.Duration(intOrDefault(fileCfg.MemoryWorkTimeoutSeconds, 110)) * time.Second,
		ConversationLockWaitTimeout:    time.Duration(intOrDefault(fileCfg.ConversationLockWaitTimeoutSeconds, 130)) * time.Second,
		MemoryConsolidationInterval:    time.Duration(intOrDefault(fileCfg.MemoryConsolidationIntervalHours, 24)) * time.Hour,
		MemoryConsolidationMinSessions: intOrDefault(fileCfg.MemoryConsolidationMinSessions, 5),
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

func CynosureSessionDir(workspaceRoot string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	workspace, err := WorkspaceName(workspaceRoot)
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cynosure", "session", workspace), nil
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

// WorkspaceName 返回用于 ~/.cynosure/ 下隔离工作区运行期数据的目录名。
// 默认使用工作区目录名；同名工作区按首次登记顺序追加 _1、_2 等后缀。
func WorkspaceName(workspaceRoot string) (string, error) {
	workspaceRegistryMu.Lock()
	defer workspaceRegistryMu.Unlock()

	abs, err := filepath.Abs(strings.TrimSpace(workspaceRoot))
	if err != nil {
		return "", fmt.Errorf("resolve workspace path: %w", err)
	}
	abs = filepath.Clean(abs)
	base := workspaceBaseName(abs)
	path, err := workspaceRegistryPath()
	if err != nil {
		return "", err
	}
	registry, err := readWorkspaceRegistry(path)
	if err != nil {
		return "", err
	}
	if name := registry.Workspaces[abs]; name != "" {
		return name, nil
	}
	used := make(map[string]struct{}, len(registry.Workspaces))
	for _, name := range registry.Workspaces {
		if name != "" {
			used[name] = struct{}{}
		}
	}
	name := nextWorkspaceName(base, used)
	if registry.Workspaces == nil {
		registry.Workspaces = make(map[string]string)
	}
	registry.Workspaces[abs] = name
	if err := writeWorkspaceRegistry(path, registry); err != nil {
		return "", err
	}
	return name, nil
}

type workspaceRegistry struct {
	Workspaces map[string]string `json:"workspaces"`
}

func workspaceBaseName(workspaceRoot string) string {
	base := sanitizePathSegment(filepath.Base(filepath.Clean(workspaceRoot)))
	if base == "" {
		base = "project"
	}
	return base
}

func workspaceRegistryPath() (string, error) {
	home, err := cynosureHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "workspaces.json"), nil
}

func readWorkspaceRegistry(path string) (workspaceRegistry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return workspaceRegistry{Workspaces: make(map[string]string)}, nil
		}
		return workspaceRegistry{}, fmt.Errorf("read workspace registry %s: %w", path, err)
	}
	var registry workspaceRegistry
	if err := json.Unmarshal(data, &registry); err != nil {
		return workspaceRegistry{}, fmt.Errorf("parse workspace registry %s: %w", path, err)
	}
	if registry.Workspaces == nil {
		registry.Workspaces = make(map[string]string)
	}
	return registry, nil
}

func writeWorkspaceRegistry(path string, registry workspaceRegistry) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create workspace registry dir: %w", err)
	}
	data, err := json.MarshalIndent(registry, "", "  ")
	if err != nil {
		return fmt.Errorf("encode workspace registry: %w", err)
	}
	data = append(data, '\n')
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write workspace registry temp file: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("replace workspace registry: %w", err)
	}
	return nil
}

func nextWorkspaceName(base string, used map[string]struct{}) string {
	if _, ok := used[base]; !ok {
		return base
	}
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s_%d", base, i)
		if _, ok := used[candidate]; !ok {
			return candidate
		}
	}
}

// CynosureSessionLogsDir 返回按工作区与会话隔离的日志目录：
// ~/.cynosure/logs/<workspace>/<session_id>
func CynosureSessionLogsDir(workspaceRoot, sessionID string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	session := sanitizePathSegment(sessionID)
	if session == "" {
		session = "session"
	}
	workspace, err := WorkspaceName(workspaceRoot)
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cynosure", "logs", workspace, session), nil
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
