package hooks

import (
	"context"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/web/storage"
)

func toolAuditPreHook(ctx context.Context, h *ToolUseContext) error {
	runtimeEnv := h.State.RuntimeEnv()
	resolvedCommandPath, commandArtifactPath := resolveCommandPaths(h.Name, h.RawArgs, runtimeEnv.CommandBinDir, runtimeEnv.CommandScriptDir)
	h.Outcome.Audit.ResolvedCWD = runtimeEnv.CurrentWorkingDir
	h.Outcome.Audit.ResolvedCommandPath = resolvedCommandPath
	h.Outcome.Audit.CommandArtifactPath = commandArtifactPath
	h.Outcome.Audit.CommandArtifactSource = classifyCommandArtifactSource(runtimeEnv.WorkspaceRoot, commandArtifactPath)
	return nil
}

func toolAuditPostHook(ctx context.Context, h *ToolUseContext) error {
	if h.Outcome.Status == "success" {
		h.Outcome.Audit.OutcomeSummary = truncate(h.Outcome.Result, 500)
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

func emitToolEventHook(ctx context.Context, h *ToolUseContext) error {
	if h.State.Writer == nil {
		return nil
	}
	_ = h.State.Writer.Event("tool", map[string]any{"name": h.Name, "status": h.Outcome.Status, "result": h.Outcome.Result})
	return nil
}

func appendToolMessageHook(ctx context.Context, h *ToolUseContext) error {
	h.State.Messages = append(h.State.Messages, openai.ChatCompletionMessage{Role: "tool", ToolCallID: h.ToolCall.ID, Content: h.Outcome.MessageContent()})
	return nil
}

func newToolCallID() string { return "tc_" + randomID() }
