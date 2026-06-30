package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"github.com/Cynosure700/cynosure/cynosure/internal/agent/storage"
	"github.com/Cynosure700/cynosure/cynosure/internal/idgen"
	"github.com/Cynosure700/cynosure/cynosure/internal/logger"
	"github.com/Cynosure700/cynosure/cynosure/internal/textutil"
)

func toolAuditPreHook(ctx context.Context, h *ToolUseContext) error {
	runtimeEnv := h.State.RuntimeEnv()
	h.Outcome.Audit.ResolvedCWD = runtimeEnv.WorkspaceRoot
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
	messageID := h.State.NextMessageID()
	// 展示历史额外携带 edit_file/multi_edit 的 diff 真实行号（exec 时计算），随
	// 会话历史持久化，供 /resume 在文件后续被改动或进程重启后仍能还原准确行号。
	// 该字段只留在展示历史 History，不进入 ModelHistory，也不进入发给模型的 openai
	// 消息（buildOpenAIMessages 不拷贝它），确保不会泄露给大模型。
	displayMessage := storage.Message{ID: messageID, ConversationID: h.State.Conversation.ID, UserID: h.State.User.ID, Role: "tool", Content: content, ToolCallID: h.ToolCall.ID, EditLineStarts: h.Outcome.EditLineStarts}
	modelMessage := storage.Message{ID: messageID, ConversationID: h.State.Conversation.ID, UserID: h.State.User.ID, Role: "tool", Content: content, ToolCallID: h.ToolCall.ID}
	h.State.History = append(h.State.History, displayMessage)
	h.State.ModelHistory = append(h.State.ModelHistory, modelMessage)
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
