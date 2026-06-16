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

const defaultSubagentMaxRounds = 20

// subAgentTurnTimeout bounds a single subagent turn. It is a soft boundary
// checked between rounds.
const subAgentTurnTimeout = 1 * time.Hour

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
	runID := idgen.New("subagent")
	parentToolCallID := strings.TrimSpace(parent.ParentToolCallID)
	if parentToolCallID == "" {
		parentToolCallID = "spawn_subagent"
	}
	trace := &subagentTrace{store: s.Store, runID: runID, parentToolCallID: parentToolCallID, conversationID: parent.Conversation.ID, userID: parent.User.ID}
	childTools := NewChildToolRegistry(s.Cfg, resolvedCWD)
	childState := s.newLoopState(parent.Conversation, parent.User, task, nil, nil, nil)
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
	base := s.buildSystemPromptWithMemory(user, snapshot, "")
	return base + "\n\n<subagent>\n你是由 `spawn_subagent` 派生出来的子智能体。\n\n规则：\n- 你看不到父对话的历史记录。\n- 只能依据当前任务和工作区文件来工作。\n- 不要调用 `spawn_subagent`。\n- 完成后，只输出一段简洁的摘要，说明你做了什么、关键发现以及尚未解决的问题。\n</subagent>"
}

func (s *Service) runSubagentLoop(ctx context.Context, state *LoopState, tools *ToolRegistry, parent ToolContext, trace *subagentTrace, maxRounds int) (openai.ChatCompletionMessage, error) {
	roundsSinceTodoWrite := 0
	turnStart := time.Now()
	for round := 1; round <= maxRounds; round++ {
		if time.Since(turnStart) > subAgentTurnTimeout {
			return openai.ChatCompletionMessage{}, fmt.Errorf("subagent turn timed out after %v", subAgentTurnTimeout)
		}
		req := openai.ChatCompletionRequest{Model: s.Cfg.LLM.ModelID, Messages: state.Messages, Tools: tools.Definitions(), MaxTokens: defaultMaxTokens}
		reqBody, _ := json.Marshal(req)
		msg, finishReason, err := s.runModelRoundWithRecovery(ctx, state, req)
		respBody, _ := json.Marshal(msg)
		logger.LogLLMRound(round, fmt.Sprintf("subagent run=%s parent_tool_call=%s conversation=%s", trace.runID, trace.parentToolCallID, parent.Conversation.ID), reqBody, respBody, err)
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
			if approved, _ := s.approveToolCall(ctx, tc); !approved {
				return openai.ChatCompletionMessage{}, fmt.Errorf("subagent tool %s rejected by user", tc.Function.Name)
			}
			toolCtx.Outcome = s.executeChildToolCall(ctx, tools, parent, tc.Function.Name, tc.Function.Arguments, toolCtx.Outcome.Audit)
			if toolCtx.Name == agenttools.TodoWriteToolName && toolCtx.Outcome.Status == "success" {
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
	return t.store.CreateSubagentMessage(ctx, storage.SubagentMessage{ID: idgen.New("submsg"), RunID: t.runID, ParentToolCallID: t.parentToolCallID, ConversationID: t.conversationID, UserID: t.userID, SequenceNo: t.sequenceNo, Role: msg.Role, Content: msg.Content, ReasoningContent: msg.ReasoningContent, ToolCallID: msg.ToolCallID, ToolCalls: msg.ToolCalls})
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
