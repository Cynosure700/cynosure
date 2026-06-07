package runtime

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"nano_cc/internal/config"
	"nano_cc/internal/textutil"
	"nano_cc/internal/web/runtime/compression"
	"nano_cc/internal/web/storage"
)

func compToolMsg(id, status, result string) storage.Message {
	content, _ := json.Marshal(struct {
		Status string `json:"status"`
		Result string `json:"result"`
	}{Status: status, Result: result})
	return storage.Message{ID: "m_" + id, Role: "tool", ToolCallID: id, Content: string(content)}
}

func compAssistantToolCallMsg(callID string) storage.Message {
	return storage.Message{Role: "assistant", ToolCalls: []storage.MessageToolCall{{ID: callID, Type: "function", Function: storage.MessageFunctionCall{Name: "bash", Arguments: "{}"}}}}
}

func compResultOf(t *testing.T, content string) string {
	t.Helper()
	_, result, _ := textutil.ParseToolResult(content)
	return result
}

// --- compressContextBeforeLLM does not mutate state.History ---

func TestCompressContextBeforeLLM_DoesNotMutateDisplayHistory(t *testing.T) {
	store := &fakeStore{}
	cfg := config.AppConfig{LLM: config.Config{ModelID: "m"}}
	service := &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}

	big := strings.Repeat("b", 300*1024)
	history := []storage.Message{
		{Role: "user", Content: "go"},
		compAssistantToolCallMsg("c1"),
		compToolMsg("c1", "success", big),
	}
	state := &LoopState{Conversation: storage.Conversation{ID: "c"}, User: storage.User{ID: "u"}, History: history, SystemPrompt: "sys"}

	requestHistory, err := service.compressContextBeforeLLM(context.Background(), state)
	if err != nil {
		t.Fatalf("compress: %v", err)
	}
	// Display history untouched.
	if compResultOf(t, state.History[2].Content) != big {
		t.Fatalf("expected display history tool result untouched")
	}
	// Request history compacted.
	if !strings.Contains(compResultOf(t, requestHistory[2].Content), compression.PersistedOutputMarkerPrefix) {
		t.Fatalf("expected request history compacted to marker")
	}
}
