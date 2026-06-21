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

func TestCompressContextBeforeLLM_UsesToolRegistryResultLimit(t *testing.T) {
	store := &fakeStore{}
	cfg := config.AppConfig{LLM: config.Config{ModelID: "m"}}
	tools := NewToolRegistry(cfg)
	tools.maxResultSizeChars["bash"] = 10
	service := &Service{Store: store, Cfg: cfg, Tools: tools}

	result := "12345678901"
	history := []storage.Message{
		{Role: "user", Content: "go"},
		compAssistantToolCallMsg("c1"),
		compToolMsg("c1", "success", result),
	}
	state := &LoopState{Conversation: storage.Conversation{ID: "c"}, User: storage.User{ID: "u"}, History: history, ModelHistory: cloneMessages(history), SystemPrompt: "sys"}

	requestHistory, err := service.compressContextBeforeLLM(context.Background(), state)
	if err != nil {
		t.Fatalf("compress: %v", err)
	}
	if compResultOf(t, state.History[2].Content) != result {
		t.Fatalf("expected display history tool result untouched")
	}
	if !strings.Contains(compResultOf(t, requestHistory[2].Content), compression.PersistedOutputMarkerPrefix) {
		t.Fatalf("expected request history compacted by tool registry limit")
	}
}

func TestCompressSubagentContextBeforeLLM_UsesChildToolRegistryResultLimit(t *testing.T) {
	store := &fakeStore{}
	cfg := config.AppConfig{LLM: config.Config{ModelID: "m"}}
	parentTools := NewToolRegistry(cfg)
	childTools := NewToolRegistry(cfg)
	childTools.maxResultSizeChars["bash"] = 10
	service := &Service{Store: store, Cfg: cfg, Tools: parentTools}

	result := "12345678901"
	history := []storage.Message{
		{Role: "user", Content: "inspect"},
		compAssistantToolCallMsg("c1"),
		compToolMsg("c1", "success", result),
	}
	state := &LoopState{Conversation: storage.Conversation{ID: "c"}, User: storage.User{ID: "u"}, History: history, ModelHistory: cloneMessages(history), SystemPrompt: "sys"}

	requestHistory, err := service.compressSubagentContextBeforeLLM(context.Background(), state, childTools)
	if err != nil {
		t.Fatalf("compress subagent: %v", err)
	}
	if compResultOf(t, state.History[2].Content) != result {
		t.Fatalf("expected display history tool result untouched")
	}
	if !strings.Contains(compResultOf(t, requestHistory[2].Content), compression.PersistedOutputMarkerPrefix) {
		t.Fatalf("expected subagent request history compacted by child tool registry limit")
	}
}

func TestCompressSubagentContextBeforeLLM_DoesNotInjectConversationMemory(t *testing.T) {
	store := &fakeStore{conversationMemories: []storage.ConversationMemory{{
		ID:             "mem_1",
		ConversationID: "c",
		UserID:         "u",
		Name:           "parent-memory",
		Body:           "parent memory must not enter child context",
	}}}
	cfg := config.AppConfig{LLM: config.Config{ModelID: "m"}}
	service := &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}
	state := &LoopState{
		Conversation: storage.Conversation{ID: "c"},
		User:         storage.User{ID: "u"},
		History:      []storage.Message{{Role: "user", Content: "child task"}},
		ModelHistory: []storage.Message{{Role: "user", Content: "child task"}},
		SystemPrompt: "sys",
	}

	requestHistory, err := service.compressSubagentContextBeforeLLM(context.Background(), state, service.Tools)
	if err != nil {
		t.Fatalf("compress subagent: %v", err)
	}
	for _, msg := range requestHistory {
		if strings.Contains(msg.Content, "<conversation-memory>") || strings.Contains(msg.Content, "parent memory must not enter child context") {
			t.Fatalf("subagent compression injected conversation memory: %#v", requestHistory)
		}
	}
}

// --- ModelHistory carries compression output forward across rounds ---

// TestModelHistory_CompressionWriteBackIsStableAcrossRounds verifies the
// unified-history invariant: once a round's compression output is written back
// into state.ModelHistory, re-running compression on that same line (the next
// round's seed) keeps the persisted-output marker in place (idempotent) and
// never resurrects the original oversized result, while the verbatim display
// history stays untouched.
func TestModelHistory_CompressionWriteBackIsStableAcrossRounds(t *testing.T) {
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

	// Round 1: compress and write back into ModelHistory (the main loop's job).
	round1, err := service.compressContextBeforeLLM(context.Background(), state)
	if err != nil {
		t.Fatalf("compress round 1: %v", err)
	}
	state.ModelHistory = round1
	if !strings.Contains(compResultOf(t, state.ModelHistory[2].Content), compression.PersistedOutputMarkerPrefix) {
		t.Fatalf("expected round 1 to compact tool result into a marker")
	}

	// Round 2: the next round seeds from the (already compressed) ModelHistory.
	round2, err := service.compressContextBeforeLLM(context.Background(), state)
	if err != nil {
		t.Fatalf("compress round 2: %v", err)
	}
	state.ModelHistory = round2
	got := compResultOf(t, state.ModelHistory[2].Content)
	if !strings.Contains(got, compression.PersistedOutputMarkerPrefix) {
		t.Fatalf("expected round 2 to preserve the marker, got %q", got)
	}
	if strings.Contains(got, big) {
		t.Fatalf("expected round 2 not to resurrect the original oversized result")
	}
	// Display history stays verbatim across rounds.
	if compResultOf(t, state.History[2].Content) != big {
		t.Fatalf("expected display history to remain untouched")
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
