package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"cynosure/internal/agent/storage"
	"cynosure/internal/idgen"
	"cynosure/internal/logger"
	"cynosure/internal/textutil"
)

func toolAuditPreHook(ctx context.Context, h *ToolUseContext) error {
	runtimeEnv := h.State.RuntimeEnv()
	h.Outcome.Audit.ResolvedCWD = runtimeEnv.CurrentWorkingDir
	h.Outcome.Audit.ResolvedCommandPath = resolveCommandPath(h.Name, h.RawArgs)
	return nil
}

func toolAuditPostHook(ctx context.Context, h *ToolUseContext) error {
	if h.Outcome.Status == "success" {
		h.Outcome.Audit.OutcomeSummary = textutil.Truncate(h.Outcome.Result, 500)
		return nil
	}
	h.Outcome.Audit.DenialReason = h.Outcome.Result
	return nil
}

func persistToolCallHook(ctx context.Context, h *ToolUseContext) error {
	state := h.State
	_ = state.Store.CreateToolCall(ctx, storage.ToolCall{ID: newToolCallID(), ConversationID: state.Conversation.ID, UserID: state.User.ID, ToolName: h.Name, Status: h.Outcome.Status, Summary: h.Outcome.AuditSummary()})
	return nil
}

func appendToolMessageHook(ctx context.Context, h *ToolUseContext) error {
	content := h.Outcome.MessageContent()
	h.State.Messages = append(h.State.Messages, openai.ChatCompletionMessage{Role: "tool", ToolCallID: h.ToolCall.ID, Content: content})
	toolMessage := storage.Message{ID: h.State.NextMessageID(), ConversationID: h.State.Conversation.ID, UserID: h.State.User.ID, Role: "tool", Content: content, ToolCallID: h.ToolCall.ID}
	h.State.History = append(h.State.History, toolMessage)
	h.State.ModelHistory = append(h.State.ModelHistory, toolMessage)
	appendToolResultLog(ctx, h)
	return nil
}

type toolResultLogStore interface {
	AppendToolResultLog(ctx context.Context, entry storage.ToolResultLogEntry) error
}

func appendToolResultLog(ctx context.Context, h *ToolUseContext) {
	if h == nil || h.State == nil || h.State.Store == nil {
		return
	}
	store, ok := h.State.Store.(toolResultLogStore)
	if !ok {
		return
	}
	entry := storage.ToolResultLogEntry{
		ConversationID: h.State.Conversation.ID,
		SessionID:      h.State.Conversation.SessionID,
		UserID:         h.State.User.ID,
		ToolCallID:     h.ToolCall.ID,
		ToolName:       h.Name,
		RawArgs:        h.RawArgs,
		Status:         h.Outcome.Status,
		Result:         h.Outcome.Result,
		AuditSummary:   h.Outcome.AuditSummary(),
		CreatedAt:      time.Now(),
	}
	if err := store.AppendToolResultLog(ctx, entry); err != nil {
		logger.Warn(fmt.Sprintf("tool result log: append failed conversation=%s tool_call=%s: %v", entry.ConversationID, entry.ToolCallID, err))
	}
}

func newToolCallID() string { return "tc_" + idgen.Hex() }

func resolveCommandPath(toolName, rawArgs string) string {
	if toolName != "bash" {
		return ""
	}
	var args map[string]any
	if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
		return ""
	}
	command, _ := args["command"].(string)
	for _, token := range strings.Fields(command) {
		candidate := strings.Trim(token, "\"'`;,()[]{}")
		if candidate == "" || !filepath.IsAbs(candidate) {
			continue
		}
		resolved, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}
		return filepath.Clean(resolved)
	}
	return ""
}
