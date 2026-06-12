package runtime

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"nano_cc/internal/agent/runtime/compression"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/config"
	"nano_cc/internal/textutil"
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
	state := &LoopState{Conversation: storage.Conversation{ID: "c"}, User: storage.User{ID: "u"}, History: history, ModelHistory: cloneMessages(history), SystemPrompt: "sys"}

	requestHistory, err := service.compressContextBeforeLLM(context.Background(), state)
	if err != nil {
		t.Fatalf("compress: %v", err)
	}
	// Display history untouched.
	if compResultOf(t, state.History[2].Content) != big {
		t.Fatalf("expected display history tool result untouched")
	}
	// Model history untouched (compression works on a clone).
	if compResultOf(t, state.ModelHistory[2].Content) != big {
		t.Fatalf("expected model history tool result untouched")
	}
	// Request history compacted.
	if !strings.Contains(compResultOf(t, requestHistory[2].Content), compression.PersistedOutputMarkerPrefix) {
		t.Fatalf("expected request history compacted to marker")
	}
}

// --- loadModelHistory ---

func TestLoadModelHistory_UsesStoredModelHistory(t *testing.T) {
	stored := []storage.Message{{Role: "user", Content: "<conversation-summary>prev</conversation-summary>"}, {Role: "assistant", Content: "prev answer"}}
	store := &fakeStore{modelHistory: stored, modelHistoryExists: true}
	cfg := config.AppConfig{LLM: config.Config{ModelID: "m"}}
	service := &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}

	display := []storage.Message{{Role: "user", Content: "full original 1"}, {Role: "assistant", Content: "full original 2"}}
	got := service.loadModelHistory(context.Background(), "c", display)
	if len(got) != 2 || got[0].Content != stored[0].Content {
		t.Fatalf("expected stored model history reused, got %#v", got)
	}
}

func TestLoadModelHistory_FallsBackToDisplayWhenAbsent(t *testing.T) {
	store := &fakeStore{} // no model history row
	cfg := config.AppConfig{LLM: config.Config{ModelID: "m"}}
	service := &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}

	display := []storage.Message{{Role: "user", Content: "orig"}}
	got := service.loadModelHistory(context.Background(), "c", display)
	if len(got) != 1 || got[0].Content != "orig" {
		t.Fatalf("expected fallback to display history, got %#v", got)
	}
}

func TestLoadModelHistory_FallsBackOnError(t *testing.T) {
	store := &fakeStore{modelHistoryErr: errors.New("db down")}
	cfg := config.AppConfig{LLM: config.Config{ModelID: "m"}}
	service := &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}

	display := []storage.Message{{Role: "user", Content: "orig"}}
	got := service.loadModelHistory(context.Background(), "c", display)
	if len(got) != 1 || got[0].Content != "orig" {
		t.Fatalf("expected fallback to display history on error, got %#v", got)
	}
}
