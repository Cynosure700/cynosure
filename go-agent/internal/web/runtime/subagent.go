package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/logger"
	agenttools "nano_cc/internal/tools"
	"nano_cc/internal/web/storage"
)

const defaultSubagentMaxRounds = 20

type spawnSubagentArgs struct {
	Task string `json:"task"`
	CWD  string `json:"cwd"`
}

type subagentContextKey string

const subagentDepthKey subagentContextKey = "subagent_depth"

type subagentTrace struct {
	store            conversationStore
	runID            string
	parentToolCallID string
	conversationID   string
	userID           string
	sequenceNo       int
}

func (s *Service) runSubagent(ctx context.Context, parent ToolContext, args spawnSubagentArgs, audit toolExecutionAudit) (string, error) {
	if depth, _ := ctx.Value(subagentDepthKey).(int); depth > 0 {
		return "", fmt.Errorf("spawn_subagent cannot be called from a subagent")
	}
	task := strings.TrimSpace(args.Task)
	if task == "" {
		return "", fmt.Errorf("task is required")
	}
	resolvedCWD, err := resolveSubagentCWD(s.Cfg.WorkspaceRoot, args.CWD)
	if err != nil {
		return "", err
	}
	runID := newID("subagent")
	parentToolCallID := strings.TrimSpace(parent.ParentToolCallID)
	if parentToolCallID == "" {
		parentToolCallID = "spawn_subagent"
	}
	trace := &subagentTrace{store: s.Store, runID: runID, parentToolCallID: parentToolCallID, conversationID: parent.Conversation.ID, userID: parent.User.ID}
	childTools := NewChildToolRegistry(s.Cfg, resolvedCWD)
	childState := s.newLoopState(parent.Conversation, parent.User, task, nil, nil)
	childState.SkillSnapshot = parent.Skills
	childState.SystemPrompt = s.buildSubagentSystemPrompt(parent.User, parent.Skills)
	childState.Messages = []openai.ChatCompletionMessage{{Role: "system", Content: childState.SystemPrompt}, {Role: "user", Content: task}}
	if err := trace.record(ctx, storage.Message{ID: childState.NextMessageID(), ConversationID: parent.Conversation.ID, UserID: parent.User.ID, Role: "user", Content: task}); err != nil {
		return "", err
	}
	childState.ToolRuntimeEnv = childTools.runtimeEnv
	childCtx := context.WithValue(ctx, subagentDepthKey, 1)
	msg, err := s.runSubagentLoop(childCtx, childState, childTools, parent, trace, defaultSubagentMaxRounds)
	if err != nil {
		return "", err
	}
	return "Subagent completed.\n\nSummary:\n" + fallbackAssistantContent(msg.Content), nil
}

func (s *Service) buildSubagentSystemPrompt(user storage.User, snapshot *agenttools.SkillSnapshot) string {
	base := s.buildSystemPrompt(user, snapshot)
	return base + "\n\n---\n\n## Child Agent Context\n\nYou are a child agent spawned by `spawn_subagent`.\n\nRules:\n- You cannot see the parent conversation history.\n- Work only from the current task and workspace files.\n- Do not call `spawn_subagent`.\n- When finished, output only a concise summary of what you did, key findings, and unresolved items."
}

func (s *Service) runSubagentLoop(ctx context.Context, state *LoopState, tools *ToolRegistry, parent ToolContext, trace *subagentTrace, maxRounds int) (openai.ChatCompletionMessage, error) {
	roundsSinceTodoWrite := 0
	for round := 1; round <= maxRounds; round++ {
		req := openai.ChatCompletionRequest{Model: s.Cfg.LLM.ModelID, Messages: state.Messages, Tools: tools.Definitions()}
		reqBody, _ := json.Marshal(req)
		msg, finishReason, err := s.runModelRoundStream(ctx, state, req)
		respBody, _ := json.Marshal(msg)
		logger.LogLLMRound(round, fmt.Sprintf("subagent run=%s parent_tool_call=%s conversation=%s", trace.runID, trace.parentToolCallID, parent.Conversation.ID), reqBody, respBody, err)
		if err != nil {
			return openai.ChatCompletionMessage{}, err
		}
		state.Messages = append(state.Messages, msg)
		if toolCallsInclude(msg.ToolCalls, todoWriteToolName) {
			roundsSinceTodoWrite = 0
		} else {
			roundsSinceTodoWrite++
		}
		storedAssistant := storage.Message{ID: state.NextMessageID(), ConversationID: parent.Conversation.ID, UserID: parent.User.ID, Role: "assistant", Content: msg.Content, ReasoningContent: msg.ReasoningContent, ToolCalls: openAIToolCallsToStorage(msg.ToolCalls)}
		if err := trace.record(ctx, storedAssistant); err != nil {
			return openai.ChatCompletionMessage{}, err
		}
		if finishReason != "tool_calls" || len(msg.ToolCalls) == 0 {
			return msg, nil
		}
		state.History = append(state.History, storedAssistant)
		for _, tc := range msg.ToolCalls {
			toolCtx := &ToolUseContext{State: state, ToolCall: tc, Name: tc.Function.Name, RawArgs: tc.Function.Arguments}
			if err := s.hookManager().RunPreToolUse(ctx, toolCtx); err != nil {
				return openai.ChatCompletionMessage{}, err
			}
			toolCtx.Outcome = s.executeChildToolCall(ctx, tools, parent, tc.Function.Name, tc.Function.Arguments, toolCtx.Outcome.Audit)
			if toolCtx.Name == todoWriteToolName && toolCtx.Outcome.Status == "success" {
				state.Todos = append([]agenttools.TodoItem(nil), toolCtx.Outcome.Todos...)
			}
			if err := s.hookManager().RunPostToolUse(ctx, toolCtx); err != nil {
				return openai.ChatCompletionMessage{}, err
			}
			if len(state.History) > 0 {
				if err := trace.record(ctx, state.History[len(state.History)-1]); err != nil {
					return openai.ChatCompletionMessage{}, err
				}
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

func (t *subagentTrace) record(ctx context.Context, msg storage.Message) error {
	if t == nil || t.store == nil {
		return nil
	}
	t.sequenceNo++
	return t.store.CreateSubagentMessage(ctx, storage.SubagentMessage{ID: newID("submsg"), RunID: t.runID, ParentToolCallID: t.parentToolCallID, ConversationID: t.conversationID, UserID: t.userID, SequenceNo: t.sequenceNo, Role: msg.Role, Content: msg.Content, ReasoningContent: msg.ReasoningContent, ToolCallID: msg.ToolCallID, ToolCalls: msg.ToolCalls})
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
	if resolved != root && !strings.HasPrefix(resolved, root+string(filepath.Separator)) {
		return "", fmt.Errorf("subagent cwd escapes workspace: %s", cwd)
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
