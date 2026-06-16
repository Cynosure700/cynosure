package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	runtimehooks "nano_cc/internal/agent/runtime/hooks"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/config"
	agenttools "nano_cc/internal/tools"
)

type toolExecutionOutcome = runtimehooks.ToolExecutionOutcome
type toolExecutionAudit = runtimehooks.ToolExecutionAudit

type ToolContext struct {
	User                  storage.User
	Conversation          storage.Conversation
	Skills                *agenttools.SkillSnapshot
	ParentToolCallID      string
	PersistedOutputReader agenttools.PersistedOutputReader
}

type ToolExecutionResult struct {
	Output string
	Todos  []agenttools.TodoItem
}

type ToolRegistry struct {
	definitions []openai.Tool
	baseEnv     agenttools.RuntimeEnv
}

const defaultAllowedTool = "load_skill"

func NewToolRegistry(cfg config.AppConfig) *ToolRegistry {
	allowed := loadAllowedToolNames(cfg)
	definitions := appendPersistedOutputTool(buildToolDefinitions(allowed))
	return &ToolRegistry{
		definitions: definitions,
		baseEnv:     runtimeEnvFromConfig(cfg),
	}
}

// appendPersistedOutputTool always exposes read_persisted_output to the main
// agent so that <persisted-output> markers produced by context compression
// remain readable regardless of the configured AllowedTools list.
func appendPersistedOutputTool(defs []openai.Tool) []openai.Tool {
	for _, def := range defs {
		if def.Function != nil && def.Function.Name == agenttools.ReadPersistedOutputToolName {
			return defs
		}
	}
	return append(defs, agenttools.ReadPersistedOutputToolDef)
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

// toolDefinitionsForUser 返回内置工具定义，并在 MCP 启用时合并该用户已连接 MCP 服务器的工具。
func (s *Service) toolDefinitionsForUser(ctx context.Context, userID string) []openai.Tool {
	defs := s.Tools.Definitions()
	if s.MCP == nil {
		return defs
	}
	s.MCP.EnsureBuiltinSessions(ctx)
	s.MCP.EnsureWorkspaceSessions(ctx)
	mcpTools := s.MCP.ToolsForUser(userID)
	if len(mcpTools) == 0 {
		return defs
	}
	return append(defs, mcpTools...)
}

func (r *ToolRegistry) Execute(ctx context.Context, toolCtx ToolContext, name string, rawArgs string) (ToolExecutionResult, error) {
	var args map[string]any
	if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
		return ToolExecutionResult{}, fmt.Errorf("invalid tool arguments: %w", err)
	}
	if !r.isAllowed(name) {
		return ToolExecutionResult{}, fmt.Errorf("tool %s is not registered for local runtime", name)
	}
	if def, ok := r.lookupDefinition(name); ok && def.Function != nil {
		if err := agenttools.ValidateToolArgs(name, agenttools.RawSchemaFromParameters(def.Function.Parameters), args); err != nil {
			return ToolExecutionResult{}, err
		}
	}
	ctx = agenttools.WithRuntimeEnv(ctx, r.runtimeEnv())
	ctx = agenttools.WithSkillSnapshot(ctx, toolCtx.Skills)
	if toolCtx.PersistedOutputReader != nil {
		ctx = agenttools.WithPersistedOutputReader(ctx, toolCtx.PersistedOutputReader)
	}
	execResult, err := agenttools.Dispatch(ctx, name, args)
	if err != nil {
		return ToolExecutionResult{}, err
	}
	return ToolExecutionResult{Output: execResult.Output, Todos: execResult.Todos}, nil
}

func (r *ToolRegistry) runtimeEnv() agenttools.RuntimeEnv {
	env := r.baseEnv
	workspaceRoot := strings.TrimSpace(env.WorkspaceRoot)
	currentWorkingDir := strings.TrimSpace(env.CurrentWorkingDir)
	if currentWorkingDir == "" {
		currentWorkingDir = workspaceRoot
	}
	return agenttools.RuntimeEnv{
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

func (r *ToolRegistry) lookupDefinition(name string) (openai.Tool, bool) {
	for _, tool := range r.definitions {
		if tool.Function != nil && tool.Function.Name == name {
			return tool, true
		}
	}
	return openai.Tool{}, false
}

func loadAllowedToolNames(cfg config.AppConfig) []string {
	configured := cfg.AllowedTools
	if len(configured) == 0 {
		configured = []string{defaultAllowedTool}
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
		WorkspaceRoot:          strings.TrimSpace(cfg.WorkspaceRoot),
		AllowOutsideWorkspace:  cfg.BashAllowOutsideWorkspace,
		AllowDangerousCommands: cfg.BashAllowDangerousCommands,
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
	for _, tool := range agenttools.AllToolDefs {
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

func (s *Service) executeToolCall(ctx context.Context, toolCtx ToolContext, name string, rawArgs string, audit toolExecutionAudit) toolExecutionOutcome {
	if strings.HasPrefix(name, "mcp__") {
		if s.MCP == nil {
			return toolExecutionOutcome{Status: "rejected", Result: "Error: MCP is not enabled", Audit: audit}
		}
		output, err := s.MCP.CallTool(ctx, toolCtx.User.ID, name, rawArgs)
		if err != nil {
			return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Error: %v", err), Audit: audit}
		}
		return toolExecutionOutcome{Status: "success", Result: output, Audit: audit}
	}
	if name == "spawn_subagent" {
		if s.Tools == nil || !s.Tools.isAllowed(name) {
			return toolExecutionOutcome{Status: "rejected", Result: "Error: tool spawn_subagent is not registered for local runtime", Audit: audit}
		}
		var rawMap map[string]any
		if err := json.Unmarshal([]byte(rawArgs), &rawMap); err != nil {
			return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Error: invalid spawn_subagent arguments: %v", err), Audit: audit}
		}
		if def, ok := s.Tools.lookupDefinition(name); ok && def.Function != nil {
			if err := agenttools.ValidateToolArgs(name, agenttools.RawSchemaFromParameters(def.Function.Parameters), rawMap); err != nil {
				return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Error: %v", err), Audit: audit}
			}
		}
		var args spawnSubagentArgs
		if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
			return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Error: invalid spawn_subagent arguments: %v", err), Audit: audit}
		}
		result, err := s.runSubagent(ctx, toolCtx, args, audit)
		if err != nil {
			return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Subagent failed: %v", err), Audit: audit}
		}
		return toolExecutionOutcome{Status: "success", Result: result, Audit: audit}
	}
	execResult, err := s.Tools.Execute(s.withWebProcessor(ctx), toolCtx, name, rawArgs)
	if err != nil {
		return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Error: %v", err), Audit: audit}
	}
	return toolExecutionOutcome{Status: "success", Result: execResult.Output, Audit: audit, Todos: execResult.Todos}
}

// withWebProcessor injects an LLM-backed processor used by the web_fetch tool
// to summarize/extract from fetched content. When no LLM is configured the
// context is returned unchanged and web_fetch falls back to raw text.
func (s *Service) withWebProcessor(ctx context.Context) context.Context {
	if s.LLM == nil {
		return ctx
	}
	processor := func(ctx context.Context, prompt, content string) (string, error) {
		messages := []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "You analyze web page content and answer the user's prompt using only the provided content."},
			{Role: openai.ChatMessageRoleUser, Content: fmt.Sprintf("%s\n\nWeb page content:\n%s", prompt, content)},
		}
		resp, err := s.LLM.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model:    s.Cfg.LLM.ModelID,
			Messages: messages,
		})
		if err != nil {
			return "", err
		}
		if len(resp.Choices) == 0 {
			return "", fmt.Errorf("empty response from model")
		}
		return resp.Choices[0].Message.Content, nil
	}
	return agenttools.WithWebProcessor(ctx, processor)
}
