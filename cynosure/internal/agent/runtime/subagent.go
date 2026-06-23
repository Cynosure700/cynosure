package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"cynosure/internal/agent/storage"
	"cynosure/internal/idgen"
	"cynosure/internal/logger"
	agenttools "cynosure/internal/tools"
)

const (
	defaultSubagentMaxRounds = 50
	exploreSubagentMaxRounds = 30
	subagentSummaryRound     = 0

	subagentMaxRoundsSummaryPrompt = `请基于历史会话进行总结输出。不要继续调用工具，直接总结已完成的观察、关键结论、仍未解决的问题，以及建议父 Agent 下一步如何处理。`
)

// subAgentTurnTimeout 限定单个子 Agent 回合的耗时上限，它是在各轮之间进行检查的
// 软边界。
const subAgentTurnTimeout = 1 * time.Hour

type spawnSubagentArgs struct {
	SubType string `json:"sub_type"`
	Task    string `json:"task"`
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
	runID := idgen.New("subagent")
	profile := s.buildSubagentProfile(kind, parent.User, parent.Skills)
	childTools := profile.ToolRegistry
	childState := s.newLoopState(parent.Conversation, parent.User, task, nil, nil, newSubagentEventWriter(parent, runID))
	childState.SkillSnapshot = parent.Skills
	childState.SystemPrompt = profile.SystemPrompt
	childState.UserMessage = storage.Message{ID: childState.NextMessageID(), ConversationID: parent.Conversation.ID, UserID: parent.User.ID, Role: "user", Content: task}
	childState.History = []storage.Message{childState.UserMessage}
	childState.ModelHistory = cloneMessages(childState.History)
	childState.Messages = buildOpenAIMessages(childState.SystemPrompt, childState.ModelHistory)
	childState.ToolRuntimeEnv = childTools.runtimeEnv
	childCtx := context.WithValue(ctx, subagentDepthKey, 1)
	msg, err := s.runSubagentLoop(childCtx, childState, childTools, parent, runID, profile.MaxRounds)
	emitToolCallGroupClear(parent, runID)
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

func (s *Service) buildSubagentProfile(kind subagentType, user storage.User, snapshot *agenttools.SkillSnapshot) subagentProfile {
	switch kind {
	case subagentTypeExplore:
		return subagentProfile{
			Type:         kind,
			SystemPrompt: s.buildExploreSubagentSystemPrompt(),
			ToolRegistry: NewExploreToolRegistry(s.Cfg),
			MaxRounds:    exploreSubagentMaxRounds,
		}
	default:
		return subagentProfile{
			Type:         subagentTypeGeneral,
			SystemPrompt: s.buildSubagentSystemPrompt(user, snapshot),
			ToolRegistry: NewChildToolRegistry(s.Cfg),
			MaxRounds:    defaultSubagentMaxRounds,
		}
	}
}

func (s *Service) buildSubagentSystemPrompt(user storage.User, snapshot *agenttools.SkillSnapshot) string {
	base := s.buildSystemPromptWithMemory(user, snapshot, "")
	return strings.TrimSpace(base) + "\n\n" + strings.TrimSpace(s.Prompts.withDefaults().GeneralSubagent)
}

func (s *Service) buildExploreSubagentSystemPrompt() string {
	workspaceRoot := strings.TrimSpace(s.Cfg.WorkspaceRoot)
	if workspaceRoot == "" {
		workspaceRoot = "."
	}
	prompt := s.Prompts.withDefaults().ExploreSubagent
	return strings.TrimSpace(strings.ReplaceAll(prompt, "{{workspace_root}}", workspaceRoot))
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
		logger.LogLLMRound(round, fmt.Sprintf("subagent run=%s conversation=%s", runID, parent.Conversation.ID), reqBody, respBody, string(finishReason), err)
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
		childToolContext := parent
		childToolContext.Todos = state.Todos
		batch, err := s.executeToolCallBatch(ctx, state, msg.ToolCalls, toolBatchOptions{
			ToolContext:      childToolContext,
			Registry:         tools,
			UseChildRegistry: true,
			SuppressResultUI: true,
			EphemeralGroupID: runID,
		})
		if err != nil {
			return openai.ChatCompletionMessage{}, err
		}
		if batch.Rejected {
			return openai.ChatCompletionMessage{}, fmt.Errorf("subagent tool rejected by user")
		}
		roundsSinceTodoWrite = maybeAppendTodoWriteReminder(state, tools, roundsSinceTodoWrite)
	}
	msg, err := s.summarizeSubagentAtMaxRounds(ctx, state, tools, parent, runID)
	if err != nil {
		return openai.ChatCompletionMessage{}, fmt.Errorf("subagent exceeded max rounds and summary fallback failed: %w", err)
	}
	return msg, nil
}

func (s *Service) summarizeSubagentAtMaxRounds(ctx context.Context, state *LoopState, tools *ToolRegistry, parent ToolContext, runID string) (openai.ChatCompletionMessage, error) {
	appendInternalUserPrompt(state, subagentMaxRoundsSummaryPrompt)
	modelHistory, err := s.compressSubagentContextBeforeLLM(ctx, state, tools)
	if err != nil {
		return openai.ChatCompletionMessage{}, err
	}
	state.ModelHistory = modelHistory
	state.Messages = buildOpenAIMessages(state.SystemPrompt, state.ModelHistory)
	req := openai.ChatCompletionRequest{Model: s.Cfg.LLM.ModelID, Messages: state.Messages, MaxTokens: defaultMaxTokens}
	reqBody, _ := json.Marshal(req)
	msg, finishReason, err := s.runModelRoundWithRecovery(ctx, state, req)
	respBody, _ := json.Marshal(msg)
	logger.LogLLMRound(subagentSummaryRound, fmt.Sprintf("subagent run=%s conversation=%s max-round-summary", runID, parent.Conversation.ID), reqBody, respBody, string(finishReason), err)
	if err != nil {
		return openai.ChatCompletionMessage{}, err
	}
	return msg, nil
}

func (s *Service) executeChildToolCall(ctx context.Context, tools *ToolRegistry, toolCtx ToolContext, name string, rawArgs string, audit toolExecutionAudit) toolExecutionOutcome {
	if name == "spawn_subagent" {
		return toolExecutionOutcome{Status: "rejected", Result: "Error: spawn_subagent cannot be called from a subagent", Audit: audit}
	}
	execResult, err := tools.Execute(ctx, toolCtx, name, rawArgs)
	if err != nil {
		return toolExecutionOutcome{Status: "failed", Result: fmt.Sprintf("Error: %v", err), Audit: audit}
	}
	return toolExecutionOutcome{Status: "success", Result: execResult.Output, Audit: audit, Todos: execResult.Todos}
}

type subagentEventWriter struct {
	parent           EventWriter
	runID            string
	parentToolCallID string
}

func newSubagentEventWriter(parent ToolContext, runID string) EventWriter {
	if parent.Writer == nil {
		return nil
	}
	return subagentEventWriter{parent: parent.Writer, runID: runID, parentToolCallID: parent.ParentToolCallID}
}

func (w subagentEventWriter) Event(name string, data any) error {
	if w.parent == nil {
		return nil
	}
	if name != toolCallStartEvent && name != toolCallDoneEvent && name != "meta" {
		return nil
	}
	payload, ok := cloneEventMap(data)
	if !ok {
		return nil
	}
	if name == "meta" {
		delete(payload, "tool_call_count")
		return w.parent.Event(name, payload)
	}
	if id, _ := payload["tool_call_id"].(string); id != "" {
		payload["tool_call_id"] = w.runID + ":" + id
	}
	payload["scope"] = "subagent"
	payload["subagent_run_id"] = w.runID
	payload["ephemeral_group_id"] = w.runID
	if w.parentToolCallID != "" {
		payload["parent_tool_call_id"] = w.parentToolCallID
	}
	payload["suppress_result"] = true
	delete(payload, "result_preview")
	return w.parent.Event(name, payload)
}

func cloneEventMap(data any) (map[string]any, bool) {
	src, ok := data.(map[string]any)
	if !ok {
		return nil, false
	}
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst, true
}

func emitToolCallGroupClear(parent ToolContext, groupID string) {
	if parent.Writer == nil || strings.TrimSpace(groupID) == "" {
		return
	}
	_ = parent.Writer.Event(toolCallGroupClearEvent, map[string]any{"ephemeral_group_id": groupID})
}
