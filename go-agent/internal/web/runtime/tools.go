package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/config"
	"nano_cc/internal/sessions"
	agenttools "nano_cc/internal/tools"
	"nano_cc/internal/web/storage"
)

type ToolContext struct {
	User         storage.User
	Conversation storage.Conversation
	Loader       *sessions.SkillLoader // 用来加载技能信息的
	ActiveSkillDir string
}

type ToolExecutionResult struct {
	Output         string
	ActiveSkillDir string
}

type ToolRegistry struct {
	definitions []openai.Tool
	baseEnv     agenttools.RuntimeEnv
}

const defaultWebAllowedTool = "load_skill"

func NewToolRegistry(cfg config.AppConfig) *ToolRegistry {
	allowed := loadAllowedToolNames(cfg)
	return &ToolRegistry{
		definitions: buildToolDefinitions(allowed),
		baseEnv:     runtimeEnvFromConfig(cfg),
	}
}

func (r *ToolRegistry) Definitions() []openai.Tool {
	return append([]openai.Tool(nil), r.definitions...)
}

func (r *ToolRegistry) Execute(ctx context.Context, toolCtx ToolContext, name string, rawArgs string) (ToolExecutionResult, error) {
	var args map[string]any
	if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
		return ToolExecutionResult{}, fmt.Errorf("invalid tool arguments: %w", err)
	}
	if !r.isAllowed(name) {
		return ToolExecutionResult{}, fmt.Errorf("tool %s is not registered for web runtime", name)
	}
	// 走特殊处理，加载技能内容
	if name == "load_skill" {
		if toolCtx.Loader == nil {
			return ToolExecutionResult{}, fmt.Errorf("no capabilities are available in this conversation")
		}
		skillName, _ := args["name"].(string)
		return r.loadSkillContent(toolCtx.Loader, skillName)
	}
	ctx = agenttools.WithRuntimeEnv(ctx, r.runtimeEnv(toolCtx))
	handler, ok := agenttools.Handlers[name]
	if !ok || handler == nil {
		return ToolExecutionResult{}, fmt.Errorf("tool %s has no handler", name)
	}
	output, err := handler(ctx, args)
	if err != nil {
		return ToolExecutionResult{}, err
	}
	return ToolExecutionResult{Output: output}, nil
}

// load_Skill执行函数
func (r *ToolRegistry) loadSkillContent(loader *sessions.SkillLoader, skillName string) (ToolExecutionResult, error) {
	entry, err := loader.GetEntry(skillName)
	if err != nil {
		return ToolExecutionResult{}, err
	}
	content := fmt.Sprintf("<skill name=\"%s\">\n%s\n</skill>", skillName, entry.Body)
	activeSkillDir := resolveSkillWorkingDir(entry.Path)
	envNote := formatRuntimeEnvNote(r.runtimeEnv(ToolContext{ActiveSkillDir: activeSkillDir}))
	if envNote == "" {
		return ToolExecutionResult{Output: content, ActiveSkillDir: activeSkillDir}, nil
	}
	return ToolExecutionResult{Output: content + "\n\n" + envNote, ActiveSkillDir: activeSkillDir}, nil
}

func resolveSkillWorkingDir(skillPath string) string {
	skillPath = strings.TrimSpace(skillPath)
	if skillPath == "" || strings.Contains(skillPath, "://") {
		return ""
	}
	if filepath.Base(skillPath) != "SKILL.md" {
		return ""
	}
	resolved, err := filepath.Abs(skillPath)
	if err != nil {
		return filepath.Clean(filepath.Dir(skillPath))
	}
	return filepath.Clean(filepath.Dir(resolved))
}

func (r *ToolRegistry) runtimeEnv(toolCtx ...ToolContext) agenttools.RuntimeEnv {
	env := r.baseEnv
	workspaceRoot := strings.TrimSpace(env.WorkspaceRoot)
	currentWorkingDir := workspaceRoot
	if len(toolCtx) > 0 {
		// 默认切换到技能目录，只加载一个skill目录
		candidate := normalizeActiveSkillDir(workspaceRoot, toolCtx[0].ActiveSkillDir)
		if candidate != "" {
			currentWorkingDir = candidate
		}
	}
	return agenttools.RuntimeEnv{
		AppHome:          env.AppHome,
		CommandBinDir:    strings.TrimSpace(env.CommandBinDir),
		CommandScriptDir: strings.TrimSpace(env.CommandScriptDir),
		WorkspaceRoot:    workspaceRoot,
		CurrentWorkingDir: currentWorkingDir,
	}
}

func normalizeActiveSkillDir(workspaceRoot, activeSkillDir string) string {
	workspaceRoot = strings.TrimSpace(workspaceRoot)
	activeSkillDir = strings.TrimSpace(activeSkillDir)
	if workspaceRoot == "" || activeSkillDir == "" {
		return ""
	}
	resolvedWorkspace, err := filepath.Abs(workspaceRoot)
	if err != nil {
		return ""
	}
	resolvedSkillDir, err := filepath.Abs(activeSkillDir)
	if err != nil {
		return ""
	}
	resolvedWorkspace = filepath.Clean(resolvedWorkspace)
	resolvedSkillDir = filepath.Clean(resolvedSkillDir)
	if resolvedSkillDir != resolvedWorkspace && !strings.HasPrefix(resolvedSkillDir, resolvedWorkspace+string(filepath.Separator)) {
		return ""
	}
	info, err := os.Stat(resolvedSkillDir)
	if err != nil || !info.IsDir() {
		return ""
	}
	return resolvedSkillDir
}

func (r *ToolRegistry) isAllowed(name string) bool {
	for _, tool := range r.definitions {
		if tool.Function != nil && tool.Function.Name == name {
			return true
		}
	}
	return false
}

func loadAllowedToolNames(cfg config.AppConfig) []string {
	configured := cfg.WebAllowedTools
	if len(configured) == 0 {
		configured = []string{defaultWebAllowedTool}
	}

	names := make([]string, 0, len(configured))
	seen := make(map[string]struct{}, len(configured))
	for _, name := range configured {
		if _, ok := lookupRegisteredTool(name); !ok {
			continue
		}
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}
		names = append(names, name)
	}
	return names
}

func runtimeEnvFromConfig(cfg config.AppConfig) agenttools.RuntimeEnv {
	return agenttools.RuntimeEnv{
		AppHome:          strings.TrimSpace(cfg.AppHome),
		CommandBinDir:    strings.TrimSpace(cfg.CommandBinDir),
		CommandScriptDir: strings.TrimSpace(cfg.CommandScriptDir),
		WorkspaceRoot:    strings.TrimSpace(cfg.WorkspaceRoot),
	}
}

func buildToolDefinitions(allowed []string) []openai.Tool {
	toolDefs := make([]openai.Tool, 0, len(allowed))
	for _, name := range allowed {
		def, ok := lookupRegisteredTool(name)
		if !ok {
			continue
		}
		toolDefs = append(toolDefs, def)
	}
	return toolDefs
}

func lookupRegisteredTool(name string) (openai.Tool, bool) {
	for _, tool := range agenttools.ChildToolDefs {
		if tool.Function != nil && tool.Function.Name == name {
			return tool, true
		}
	}
	return openai.Tool{}, false
}

func RegisteredTools(cfg config.AppConfig) []string {
	names := loadAllowedToolNames(cfg)
	sort.Strings(names)
	return names
}

func formatRuntimeEnvNote(env agenttools.RuntimeEnv) string {
	lines := make([]string, 0, 5)
	if env.AppHome != "" {
		lines = append(lines, "APP_HOME="+env.AppHome)
	}
	if env.CommandBinDir != "" {
		lines = append(lines, "COMMAND_BIN_DIR="+env.CommandBinDir)
	}
	if env.CommandScriptDir != "" {
		lines = append(lines, "COMMAND_SCRIPT_DIR="+env.CommandScriptDir)
	}
	if env.WorkspaceRoot != "" {
		lines = append(lines, "WORKSPACE_ROOT="+env.WorkspaceRoot)
	}
	if env.CurrentWorkingDir != "" {
		lines = append(lines, "CURRENT_WORKING_DIR="+env.CurrentWorkingDir)
	}
	if len(lines) == 0 {
		return ""
	}
	return "<runtime-paths>\n" + strings.Join(lines, "\n") + "\n</runtime-paths>"
}
