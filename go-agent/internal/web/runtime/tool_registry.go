package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/config"
	agenttools "nano_cc/internal/tools"
	"nano_cc/internal/web/storage"
)

type ToolContext struct {
	User             storage.User
	Conversation     storage.Conversation
	Skills           *agenttools.SkillSnapshot
	ParentToolCallID string
}

type ToolExecutionResult struct {
	Output string
	Todos  []agenttools.TodoItem
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

func NewChildToolRegistry(cfg config.AppConfig, cwd string) *ToolRegistry {
	allowed := withoutTool(loadAllowedToolNames(cfg), "spawn_subagent")
	env := runtimeEnvFromConfig(cfg)
	env.CurrentWorkingDir = strings.TrimSpace(cwd)
	env.AllowOutsideWorkspace = false
	return &ToolRegistry{definitions: buildToolDefinitions(allowed), baseEnv: env}
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
	ctx = agenttools.WithRuntimeEnv(ctx, r.runtimeEnv())
	ctx = agenttools.WithSkillSnapshot(ctx, toolCtx.Skills)
	if name == todoWriteToolName {
		result, err := agenttools.ExecuteTodoWrite(ctx, args)
		if err != nil {
			return ToolExecutionResult{}, err
		}
		return ToolExecutionResult{Output: result.Output, Todos: result.Todos}, nil
	}
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

func (r *ToolRegistry) runtimeEnv() agenttools.RuntimeEnv {
	env := r.baseEnv
	workspaceRoot := strings.TrimSpace(env.WorkspaceRoot)
	currentWorkingDir := strings.TrimSpace(env.CurrentWorkingDir)
	if currentWorkingDir == "" {
		currentWorkingDir = workspaceRoot
	}
	return agenttools.RuntimeEnv{
		AppHome:                env.AppHome,
		CommandBinDir:          strings.TrimSpace(env.CommandBinDir),
		CommandScriptDir:       strings.TrimSpace(env.CommandScriptDir),
		WorkspaceRoot:          workspaceRoot,
		CurrentWorkingDir:      currentWorkingDir,
		AllowOutsideWorkspace:  env.AllowOutsideWorkspace,
		AllowDangerousCommands: env.AllowDangerousCommands,
	}
}

func withoutTool(names []string, excluded string) []string {
	filtered := make([]string, 0, len(names))
	for _, name := range names {
		if name == excluded {
			continue
		}
		filtered = append(filtered, name)
	}
	return filtered
}

func (r *ToolRegistry) isAllowed(name string) bool {
	for _, tool := range r.definitions {
		if tool.Function != nil && tool.Function.Name == name {
			return true
		}
	}
	return false
}
