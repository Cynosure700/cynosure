package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/idgen"
	"nano_cc/internal/logger"
	agenttools "nano_cc/internal/tools"
)

const defaultSubagentMaxRounds = 300

// subAgentTurnTimeout 限定单个子 Agent 回合的耗时上限，它是在各轮之间进行检查的
// 软边界。
const subAgentTurnTimeout = 1 * time.Hour

type spawnSubagentArgs struct {
	SubType string `json:"sub_type"`
	Task    string `json:"task"`
	CWD     string `json:"cwd"`
}

type subagentType string

const (
	subagentTypeGeneral subagentType = "general"
	subagentTypeExplore subagentType = "explore"
)

type subagentProfile struct {
	Type         subagentType
	SystemPrompt string
	ToolRegistry *ToolRegistry
	MaxRounds    int
}

type subagentContextKey string

const subagentDepthKey subagentContextKey = "subagent_depth"

func (s *Service) runSubagent(ctx context.Context, parent ToolContext, args spawnSubagentArgs, audit toolExecutionAudit) (string, error) {
	if depth, _ := ctx.Value(subagentDepthKey).(int); depth > 0 {
		return "", fmt.Errorf("spawn_subagent cannot be called from a subagent")
	}
	task := strings.TrimSpace(args.Task)
	if task == "" {
		return "", fmt.Errorf("task is required")
	}
	kind, err := parseSubagentType(args.SubType)
	if err != nil {
		return "", err
	}
	resolvedCWD, err := resolveSubagentCWD(s.Cfg.WorkspaceRoot, args.CWD)
	if err != nil {
		return "", err
	}
	runID := idgen.New("subagent")
	profile := s.buildSubagentProfile(kind, parent.User, parent.Skills, resolvedCWD)
	childTools := profile.ToolRegistry
	childState := s.newLoopState(parent.Conversation, parent.User, task, nil, nil, nil)
	childState.SkillSnapshot = parent.Skills
	childState.SystemPrompt = profile.SystemPrompt
	childState.UserMessage = storage.Message{ID: childState.NextMessageID(), ConversationID: parent.Conversation.ID, UserID: parent.User.ID, Role: "user", Content: task}
	childState.History = []storage.Message{childState.UserMessage}
	childState.ModelHistory = cloneMessages(childState.History)
	childState.Messages = buildOpenAIMessages(childState.SystemPrompt, childState.ModelHistory)
	childState.ToolRuntimeEnv = childTools.runtimeEnv
	childCtx := context.WithValue(ctx, subagentDepthKey, 1)
	msg, err := s.runSubagentLoop(childCtx, childState, childTools, parent, runID, profile.MaxRounds)
	if err != nil {
		return "", err
	}
	return "Subagent completed.\n\nSummary:\n" + fallbackAssistantContent(msg.Content), nil
}

func parseSubagentType(value string) (subagentType, error) {
	switch kind := subagentType(strings.TrimSpace(value)); kind {
	case subagentTypeGeneral, subagentTypeExplore:
		return kind, nil
	case "":
		return "", fmt.Errorf("sub_type is required")
	default:
		return "", fmt.Errorf("unsupported sub_type %q; allowed values: general, explore", value)
	}
}

func (s *Service) buildSubagentProfile(kind subagentType, user storage.User, snapshot *agenttools.SkillSnapshot, cwd string) subagentProfile {
	switch kind {
	case subagentTypeExplore:
		return subagentProfile{
			Type:         kind,
			SystemPrompt: buildExploreSubagentSystemPrompt(cwd),
			ToolRegistry: NewExploreToolRegistry(s.Cfg, cwd),
			MaxRounds:    defaultSubagentMaxRounds,
		}
	default:
		return subagentProfile{
			Type:         subagentTypeGeneral,
			SystemPrompt: s.buildSubagentSystemPrompt(user, snapshot),
			ToolRegistry: NewChildToolRegistry(s.Cfg, cwd),
			MaxRounds:    defaultSubagentMaxRounds,
		}
	}
}

func (s *Service) buildSubagentSystemPrompt(user storage.User, snapshot *agenttools.SkillSnapshot) string {
	base := s.buildSystemPromptWithMemory(user, snapshot, "")
	return base + "\n\n<subagent>\n你是由 `spawn_subagent` 派生出来的 general 子智能体。\n\n规则：\n- 你看不到父对话的历史记录。\n- 只能依据当前任务和工作区文件来工作。\n- 不要调用 `spawn_subagent`。\n- 搜索、文件定位、代码探索、实现梳理、证据收集等搜索相关任务必须交给 explore 子智能体，不应由 general 子智能体承担。\n- 完成后，只输出一段简洁的摘要，说明你做了什么、关键发现以及尚未解决的问题。\n</subagent>"
}

func buildExploreSubagentSystemPrompt(cwd string) string {
	currentWorkingDir := strings.TrimSpace(cwd)
	if currentWorkingDir == "" {
		currentWorkingDir = "."
	}
	prompt := `You are Cynosure's explore subagent, a read-only codebase search specialist.

=== READ-ONLY MODE ===

You must only inspect existing files and report findings.

You must not create, modify, delete, move, copy, install, download, persist, or overwrite files.

You must not change repository, workspace, system, network, dependency, package-manager, environment, configuration, cache, database, service, process, or runtime state.

Your job:

- Rapidly locate relevant files, symbols, configuration, tests, documentation, and implementation details.
- Read only the files required to answer the caller's search request.
- Prefer the minimum number of file reads necessary to establish evidence.
- Return a concise report with file paths, important line references when available, and confidence or gaps.

Path verification rules:

- Never assume a file or directory exists.
- Before reading any file, first verify the path exists.
- Prefer grep, glob, or listing a known parent directory to confirm existence.
- Do not call read_file on speculative, inferred, guessed, or unverified paths.
- If a path provided by the user cannot be verified, report it as "path not found" instead of attempting to read it.
- If multiple matching files exist, identify the candidates first and then read the most relevant ones.
- Only read files whose existence has been confirmed.
- Do not construct synthetic paths from naming conventions without verifying them.

Search strategy:

1. Search broadly.
2. Identify candidate files.
3. Verify file existence.
4. Read the smallest set of high-signal files.
5. Stop once sufficient evidence is collected.

Tool rules:

- Prefer grep for content search.
- Prefer glob for filename pattern matching.
- Use ls only on known existing directories.
- Use read_file only after the target path has been verified to exist.
- Read the minimum amount of content needed.
- Avoid reading large files unless necessary.
- Do not use write_file.
- Do not use edit_file.
- Do not use multi_edit.
- Do not use delete_file.
- Do not use move_file.
- Do not use copy_file.
- Do not use create_file.
- Do not use update_memory.
- Do not use delete_memory.
- Do not use spawn_subagent.
- Do not use package managers.
- Do not use dependency installation commands.
- Do not use git mutation commands.
- Do not use shell commands that modify state.
- Do not perform network operations that change state.
- If bash is unavailable, complete the search using only the available read-only tools.

Allowed behavior:

- Search.
- Inspect.
- Read.
- Analyze.
- Summarize.
- Cross-reference findings.
- Report evidence.

Forbidden behavior:

- Any filesystem modification.
- Any repository modification.
- Any configuration change.
- Any dependency change.
- Any environment change.
- Any cache change.
- Any database write.
- Any service restart.
- Any process manipulation.
- Any network-side mutation.
- Any operation whose effect persists after execution.

Environment:

- Current working directory: {{current_working_directory}}
- Treat relative paths as relative to the current working directory unless an absolute path is provided.
- Prefer workspace-relative or absolute paths consistently.
- Include sufficient path context for direct navigation.
- Parent conversation history is unavailable.
- Rely only on this task, the current working directory, and verified files.

Efficiency:

- Search broadly first.
- Read narrowly second.
- Run independent searches in parallel when supported.
- Avoid duplicate reads.
- Stop when enough evidence exists to answer the request.
- Do not perform unrelated exploration.

Evidence requirements:

- Every substantive claim should be backed by inspected files.
- Cite file paths for findings.
- Include line numbers when available.
- Distinguish confirmed facts from assumptions.
- Explicitly state uncertainty when evidence is incomplete.

Final response:

- Reply in normal text.
- Do not create files.
- Do not modify files.
- Do not suggest changes unless explicitly requested.
- Include:
  - Key findings
  - Evidence paths
  - Relevant line references
  - Confidence level
  - Unresolved gaps or missing files`
	return strings.TrimSpace(strings.ReplaceAll(prompt, "{{current_working_directory}}", currentWorkingDir))
}

func (s *Service) runSubagentLoop(ctx context.Context, state *LoopState, tools *ToolRegistry, parent ToolContext, runID string, maxRounds int) (openai.ChatCompletionMessage, error) {
	roundsSinceTodoWrite := 0
	turnStart := time.Now()
	for round := 1; round <= maxRounds; round++ {
		if time.Since(turnStart) > subAgentTurnTimeout {
			return openai.ChatCompletionMessage{}, fmt.Errorf("subagent turn timed out after %v", subAgentTurnTimeout)
		}
		modelHistory, err := s.compressSubagentContextBeforeLLM(ctx, state, tools)
		if err != nil {
			return openai.ChatCompletionMessage{}, err
		}
		state.ModelHistory = modelHistory
		state.Messages = buildOpenAIMessages(state.SystemPrompt, state.ModelHistory)
		req := openai.ChatCompletionRequest{Model: s.Cfg.LLM.ModelID, Messages: state.Messages, Tools: tools.Definitions(), MaxTokens: defaultMaxTokens}
		reqBody, _ := json.Marshal(req)
		msg, finishReason, err := s.runModelRoundWithRecovery(ctx, state, req)
		respBody, _ := json.Marshal(msg)
		logger.LogLLMRound(round, fmt.Sprintf("subagent run=%s conversation=%s", runID, parent.Conversation.ID), reqBody, respBody, err)
		if err != nil {
			return openai.ChatCompletionMessage{}, err
		}
		state.Messages = append(state.Messages, msg)
		if toolCallsInclude(msg.ToolCalls, agenttools.TodoWriteToolName) {
			roundsSinceTodoWrite = 0
		} else {
			roundsSinceTodoWrite++
		}
		storedAssistant := storage.Message{ID: state.NextMessageID(), ConversationID: parent.Conversation.ID, UserID: parent.User.ID, Role: "assistant", Content: msg.Content, ReasoningContent: msg.ReasoningContent, ToolCalls: openAIToolCallsToStorage(msg.ToolCalls)}
		if finishReason != "tool_calls" || len(msg.ToolCalls) == 0 {
			return msg, nil
		}
		state.History = append(state.History, storedAssistant)
		state.ModelHistory = append(state.ModelHistory, storedAssistant)
		for _, tc := range msg.ToolCalls {
			toolCtx := &ToolUseContext{State: state, ToolCall: tc, Name: tc.Function.Name, RawArgs: tc.Function.Arguments}
			if err := s.hookManager().RunPreToolUse(ctx, toolCtx); err != nil {
				return openai.ChatCompletionMessage{}, err
			}
			if approved, _ := s.approveToolCall(ctx, tc); !approved {
				return openai.ChatCompletionMessage{}, fmt.Errorf("subagent tool %s rejected by user", tc.Function.Name)
			}
			childToolContext := parent
			childToolContext.Todos = state.Todos
			toolCtx.Outcome = s.executeChildToolCall(ctx, tools, childToolContext, tc.Function.Name, tc.Function.Arguments, toolCtx.Outcome.Audit)
			if toolCtx.Name == agenttools.TodoWriteToolName && toolCtx.Outcome.Status == "success" {
				state.Todos = append([]agenttools.TodoItem(nil), toolCtx.Outcome.Todos...)
			}
			if err := s.hookManager().RunPostToolUse(ctx, toolCtx); err != nil {
				return openai.ChatCompletionMessage{}, err
			}
		}
		roundsSinceTodoWrite = maybeAppendTodoWriteReminder(state, tools, roundsSinceTodoWrite)
	}
	return openai.ChatCompletionMessage{}, fmt.Errorf("subagent exceeded max rounds")
}

func (s *Service) executeChildToolCall(ctx context.Context, tools *ToolRegistry, toolCtx ToolContext, name string, rawArgs string, audit toolExecutionAudit) toolExecutionOutcome {
	if name == "spawn_subagent" {
		return toolExecutionOutcome{Status: "rejected", Result: "Error: spawn_subagent cannot be called from a subagent", Audit: audit}
	}
	execResult, err := tools.Execute(ctx, toolCtx, name, rawArgs)
	if err != nil {
		return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Error: %v", err), Audit: audit}
	}
	return toolExecutionOutcome{Status: "success", Result: execResult.Output, Audit: audit, Todos: execResult.Todos}
}

func resolveSubagentCWD(workspaceRoot, cwd string) (string, error) {
	root, err := filepath.Abs(strings.TrimSpace(workspaceRoot))
	if err != nil || strings.TrimSpace(workspaceRoot) == "" {
		return "", fmt.Errorf("workspace root is required")
	}
	root = filepath.Clean(root)
	resolved := root
	if strings.TrimSpace(cwd) != "" {
		resolved = strings.TrimSpace(cwd)
		if !filepath.IsAbs(resolved) {
			resolved = filepath.Join(root, resolved)
		}
		resolved, err = filepath.Abs(resolved)
		if err != nil {
			return "", fmt.Errorf("resolve subagent cwd: %w", err)
		}
		resolved = filepath.Clean(resolved)
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", fmt.Errorf("subagent cwd is unavailable: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("subagent cwd is not a directory: %s", cwd)
	}
	return resolved, nil
}
