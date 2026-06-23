package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	runtimehooks "cynosure/internal/agent/runtime/hooks"
	"cynosure/internal/agent/storage"
	"cynosure/internal/config"
	agenttools "cynosure/internal/tools"
)

type toolExecutionOutcome = runtimehooks.ToolExecutionOutcome
type toolExecutionAudit = runtimehooks.ToolExecutionAudit

type ToolContext struct {
	User                  storage.User
	Conversation          storage.Conversation
	Skills                *agenttools.SkillSnapshot
	PersistedOutputReader agenttools.PersistedOutputReader
	Todos                 []agenttools.TodoItem
	Writer                EventWriter
	ParentToolCallID      string
}

type ToolExecutionResult struct {
	Output string
	Todos  []agenttools.TodoItem
}

type ToolRegistry struct {
	definitions        []openai.Tool
	maxResultSizeChars map[string]int
	baseEnv            agenttools.RuntimeEnv
	bashPolicy         bashExecutionPolicy
}

const defaultAllowedTool = "load_skill"

type bashExecutionPolicy int

const (
	bashPolicyDefault bashExecutionPolicy = iota
	bashPolicyExploreReadOnly
)

func NewToolRegistry(cfg config.AppConfig) *ToolRegistry {
	allowed := loadAllowedToolNames(cfg)
	definitions := appendPersistedOutputTool(buildToolDefinitions(allowed))
	return &ToolRegistry{
		definitions:        definitions,
		maxResultSizeChars: buildMaxResultSizeMap(definitions),
		baseEnv:            runtimeEnvFromConfig(cfg),
	}
}

// appendPersistedOutputTool 始终向主 Agent 暴露 read_persisted_output，使得上下文
// 压缩产生的 <persisted-output> 标记无论配置的 AllowedTools 列表如何，都保持可读。
func appendPersistedOutputTool(defs []openai.Tool) []openai.Tool {
	for _, def := range defs {
		if def.Function != nil && def.Function.Name == agenttools.ReadPersistedOutputToolName {
			return defs
		}
	}
	return append(defs, agenttools.ReadPersistedOutputToolDef)
}

func NewChildToolRegistry(cfg config.AppConfig) *ToolRegistry {
	allowed := withoutTool(loadAllowedToolNames(cfg), "spawn_subagent")
	definitions := buildToolDefinitions(allowed)
	return &ToolRegistry{definitions: definitions, maxResultSizeChars: buildMaxResultSizeMap(definitions), baseEnv: runtimeEnvFromConfig(cfg)}
}

func NewExploreToolRegistry(cfg config.AppConfig) *ToolRegistry {
	allowed := intersectTools(loadAllowedToolNames(cfg), []string{"bash", "read_file", "grep", "glob", "ls"})
	definitions := appendPersistedOutputTool(buildToolDefinitions(allowed))
	return &ToolRegistry{definitions: definitions, maxResultSizeChars: buildMaxResultSizeMap(definitions), baseEnv: runtimeEnvFromConfig(cfg), bashPolicy: bashPolicyExploreReadOnly}
}

func (r *ToolRegistry) Definitions() []openai.Tool {
	return append([]openai.Tool(nil), r.definitions...)
}

func (r *ToolRegistry) MaxResultSizeChars(name string) int {
	if r != nil && r.maxResultSizeChars != nil {
		if limit := r.maxResultSizeChars[name]; limit > 0 {
			return limit
		}
	}
	return agenttools.MaxResultSizeCharsForTool(name)
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
	if name == "bash" {
		if err := r.validateBashExecution(args); err != nil {
			return ToolExecutionResult{}, err
		}
	}
	if def, ok := r.lookupDefinition(name); ok && def.Function != nil {
		if err := agenttools.ValidateToolArgs(name, agenttools.RawSchemaFromParameters(def.Function.Parameters), args); err != nil {
			return ToolExecutionResult{}, err
		}
	}
	ctx = agenttools.WithRuntimeEnv(ctx, r.runtimeEnv())
	ctx = agenttools.WithSkillSnapshot(ctx, toolCtx.Skills)
	ctx = agenttools.WithTodoSnapshot(ctx, toolCtx.Todos)
	if toolCtx.PersistedOutputReader != nil {
		ctx = agenttools.WithPersistedOutputReader(ctx, toolCtx.PersistedOutputReader)
	}
	execResult, err := agenttools.Dispatch(ctx, name, args)
	if err != nil {
		return ToolExecutionResult{}, err
	}
	return ToolExecutionResult{Output: execResult.Output, Todos: execResult.Todos}, nil
}

func (r *ToolRegistry) validateBashExecution(args map[string]any) error {
	if r == nil || r.bashPolicy != bashPolicyExploreReadOnly {
		return nil
	}
	command, _ := args["command"].(string)
	if !isExploreReadOnlyBashCommand(command) {
		return fmt.Errorf("bash command is not allowed for explore subagent; use only read-only operations: ls, git status, git log, git diff, find, cat, head, tail")
	}
	return nil
}

func (r *ToolRegistry) runtimeEnv() agenttools.RuntimeEnv {
	return agenttools.RuntimeEnv{
		WorkspaceRoot: strings.TrimSpace(r.baseEnv.WorkspaceRoot),
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

func intersectTools(names []string, allowed []string) []string {
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, name := range allowed {
		allowedSet[name] = struct{}{}
	}
	filtered := make([]string, 0, len(names))
	for _, name := range names {
		if _, ok := allowedSet[name]; ok {
			filtered = append(filtered, name)
		}
	}
	return filtered
}

func isExploreReadOnlyBashCommand(command string) bool {
	fields := strings.Fields(strings.TrimSpace(command))
	if len(fields) == 0 {
		return false
	}
	if strings.Contains(command, ">") || strings.Contains(command, ";") || strings.Contains(command, "&") || strings.Contains(command, "|") {
		return false
	}
	switch filepath.Base(fields[0]) {
	case "ls", "find", "cat", "head", "tail", "sed":
		return !containsExploreBashDeniedArg(fields[1:])
	case "git":
		if len(fields) < 2 {
			return false
		}
		switch fields[1] {
		case "status", "log", "diff":
			return !containsExploreBashDeniedArg(fields[2:])
		default:
			return false
		}
	default:
		return false
	}
}

func containsExploreBashDeniedArg(args []string) bool {
	for _, arg := range args {
		switch {
		case arg == "-delete" || arg == "-exec" || arg == "-execdir":
			return true
		case arg == "-o" || arg == "--output":
			return true
		case strings.HasPrefix(arg, "--output="):
			return true
		}
	}
	return false
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
		WorkspaceRoot: strings.TrimSpace(cfg.WorkspaceRoot),
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

func buildMaxResultSizeMap(definitions []openai.Tool) map[string]int {
	limits := make(map[string]int, len(definitions))
	for _, def := range definitions {
		if def.Function == nil {
			continue
		}
		limits[def.Function.Name] = agenttools.MaxResultSizeCharsForTool(def.Function.Name)
	}
	return limits
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
			return toolExecutionOutcome{Status: "failed", Result: fmt.Sprintf("Error: %v", err), Audit: audit}
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
			return toolExecutionOutcome{Status: "failed", Result: fmt.Sprintf("Subagent failed: %v", err), Audit: audit}
		}
		return toolExecutionOutcome{Status: "success", Result: result, Audit: audit}
	}
	if name == agenttools.UpdateMemoryToolName || name == agenttools.DeleteMemoryToolName {
		if s.Tools == nil || !s.Tools.isAllowed(name) {
			return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Error: tool %s is not registered for local runtime", name), Audit: audit}
		}
		if def, ok := s.Tools.lookupDefinition(name); ok && def.Function != nil {
			var rawMap map[string]any
			if err := json.Unmarshal([]byte(rawArgs), &rawMap); err != nil {
				return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Error: invalid %s arguments: %v", name, err), Audit: audit}
			}
			if err := agenttools.ValidateToolArgs(name, agenttools.RawSchemaFromParameters(def.Function.Parameters), rawMap); err != nil {
				return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Error: %v", err), Audit: audit}
			}
		}
		output, err := s.executeMemoryTool(ctx, name, rawArgs)
		if err != nil {
			return toolExecutionOutcome{Status: "failed", Result: fmt.Sprintf("Error: %v", err), Audit: audit}
		}
		return toolExecutionOutcome{Status: "success", Result: output, Audit: audit}
	}
	execResult, err := s.Tools.Execute(s.withWebProcessor(ctx), toolCtx, name, rawArgs)
	if err != nil {
		return toolExecutionOutcome{Status: "failed", Result: fmt.Sprintf("Error: %v", err), Audit: audit}
	}
	return toolExecutionOutcome{Status: "success", Result: execResult.Output, Audit: audit, Todos: execResult.Todos}
}

// withWebProcessor 注入一个由 LLM 支撑的处理器，供 web_fetch 工具对抓取到的内容
// 进行摘要/提取。当未配置 LLM 时，原样返回 context，web_fetch 回退到原始文本。
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
