package runtime

import (
	"context"
	"encoding/json"
	"fmt"
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

func (r *ToolRegistry) Execute(ctx context.Context, toolCtx ToolContext, name string, rawArgs string) (string, error) {
	var args map[string]any
	if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
		return "", fmt.Errorf("invalid tool arguments: %w", err)
	}
	if !r.isAllowed(name) {
		return "", fmt.Errorf("tool %s is not registered for web runtime", name)
	}
	// 走特殊处理，加载技能内容
	if name == "load_skill" {
		if toolCtx.Loader == nil {
			return "", fmt.Errorf("no capabilities are available in this conversation")
		}
		skillName, _ := args["name"].(string)
		return r.loadSkillContent(toolCtx.Loader, skillName)
	}
	ctx = agenttools.WithRuntimeEnv(ctx, r.runtimeEnv())
	handler, ok := agenttools.Handlers[name]
	if !ok || handler == nil {
		return "", fmt.Errorf("tool %s has no handler", name)
	}
	return handler(ctx, args)
}

func (r *ToolRegistry) loadSkillContent(loader *sessions.SkillLoader, skillName string) (string, error) {
	content, err := loader.GetContent(skillName)
	if err != nil {
		return "", err
	}
	envNote := formatRuntimeEnvNote(r.runtimeEnv())
	if envNote == "" {
		return content, nil
	}
	return content + "\n\n" + envNote, nil
}

func (r *ToolRegistry) runtimeEnv() agenttools.RuntimeEnv {
	env := r.baseEnv
	workspaceRoot := strings.TrimSpace(env.WorkspaceRoot)
	commandBinDir := strings.TrimSpace(env.CommandBinDir)
	if commandBinDir == "" && workspaceRoot != "" {
		commandBinDir = filepath.Join(workspaceRoot, "bin")
	}
	commandScriptDir := strings.TrimSpace(env.CommandScriptDir)
	if commandScriptDir == "" && workspaceRoot != "" {
		commandScriptDir = filepath.Join(workspaceRoot, "cmd")
	}
	return agenttools.RuntimeEnv{
		AppHome:          env.AppHome,
		CommandBinDir:    commandBinDir,
		CommandScriptDir: commandScriptDir,
		WorkspaceRoot:    workspaceRoot,
	}
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
	lines := make([]string, 0, 4)
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
	if len(lines) == 0 {
		return ""
	}
	return "<runtime-paths>\n" + strings.Join(lines, "\n") + "\n</runtime-paths>"
}
