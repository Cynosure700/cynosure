package hooks

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/idgen"
	"nano_cc/internal/safety"
	"nano_cc/internal/textutil"
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
	h.State.History = append(h.State.History, storage.Message{ID: h.State.NextMessageID(), ConversationID: h.State.Conversation.ID, UserID: h.State.User.ID, Role: "tool", Content: content, ToolCallID: h.ToolCall.ID})
	return nil
}

func newToolCallID() string { return "tc_" + idgen.Hex() }

func classifyCommandArtifactSource(workspaceRoot, commandArtifactPath string) string {
	if strings.TrimSpace(commandArtifactPath) == "" {
		return ""
	}
	cleanArtifact := filepath.Clean(commandArtifactPath)
	cleanWorkspace := strings.TrimSpace(workspaceRoot)
	if cleanWorkspace != "" {
		cleanWorkspace = filepath.Clean(cleanWorkspace)
		if safety.Contains(cleanWorkspace, cleanArtifact) {
			return "workspace"
		}
	}
	return "custom"
}

func resolveCommandPaths(toolName, rawArgs string, roots ...string) (string, string) {
	if toolName != "bash" {
		return "", ""
	}
	var args map[string]any
	if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
		return "", ""
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
		cleanResolved := filepath.Clean(resolved)
		for _, root := range roots {
			if root == "" {
				continue
			}
			resolvedRoot, err := filepath.Abs(root)
			if err != nil {
				continue
			}
			cleanRoot := filepath.Clean(resolvedRoot)
			if safety.Contains(cleanRoot, cleanResolved) {
				return cleanResolved, cleanResolved
			}
		}
		return cleanResolved, ""
	}
	return "", ""
}
