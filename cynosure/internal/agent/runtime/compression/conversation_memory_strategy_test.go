package compression

import (
	"context"
	"strings"
	"testing"

	"cynosure/internal/agent/storage"
)

// conversationMemoryItems builds count memories sized so the rendered text lands
// within [conversationMemoryMinTokens, conversationMemoryMaxTokens].
func conversationMemoryItems(conversationID string, count, bodyBytes int) []storage.ConversationMemory {
	items := make([]storage.ConversationMemory, 0, count)
	for i := 0; i < count; i++ {
		items = append(items, storage.ConversationMemory{
			ConversationID: conversationID,
			Name:           "topic",
			Body:           strings.Repeat("m", bodyBytes),
		})
	}
	return items
}

func TestConversationMemory_NoopUnderBudget(t *testing.T) {
	store := &fakeStore{conversationMemories: conversationMemoryItems("c", 5, 8*1024)}
	history := []storage.Message{
		{ID: "m_a", Role: "user", Content: "hi"},
		{ID: "m_b", Role: "assistant", Content: "hello"},
	}
	req := &Request{
		User: storage.User{ID: "u"}, Conversation: storage.Conversation{ID: "c"},
		RequestHistory:               history,
		Store:                        store,
		Estimator:                    DefaultTokenEstimator{},
		ConversationMemoryBreakpoint: "m_a",
	}
	if err := (&ConversationMemoryStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if len(req.RequestHistory) != 2 {
		t.Fatalf("expected no-op under budget, got %#v", req.RequestHistory)
	}
}

func TestConversationMemory_ReplacesWithUnfoldedTail(t *testing.T) {
	store := &fakeStore{conversationMemories: conversationMemoryItems("c", 5, 8*1024)}
	huge := storage.Message{ID: "m_huge", Role: "user", Content: strings.Repeat("z", 600*1024)}
	// Breakpoint at m_bp; only m_bp and after are "uncompressed" tail.
	bp := storage.Message{ID: "m_bp", Role: "assistant", Content: "folded boundary"}
	t1 := storage.Message{ID: "m_t1", Role: "user", Content: "after breakpoint 1"}
	t2 := storage.Message{ID: "m_t2", Role: "assistant", Content: "after breakpoint 2"}
	req := &Request{
		User: storage.User{ID: "u"}, Conversation: storage.Conversation{ID: "c"},
		RequestHistory:               []storage.Message{huge, bp, t1, t2},
		DisplayHistory:               []storage.Message{huge, bp, t1, t2},
		Store:                        store,
		Estimator:                    DefaultTokenEstimator{},
		ConversationMemoryBreakpoint: "m_bp",
	}
	if err := (&ConversationMemoryStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	// Expect: [system preamble, user <conversation-memory>, m_bp, m_t1, m_t2].
	if len(req.RequestHistory) != 5 {
		t.Fatalf("expected preamble+memory+3 tail, got %d: %#v", len(req.RequestHistory), req.RequestHistory)
	}
	if req.RequestHistory[0].Role != "system" || req.RequestHistory[1].Role != "user" ||
		!strings.Contains(req.RequestHistory[1].Content, "<conversation-memory>") {
		t.Fatalf("expected memory preamble block, got %#v", req.RequestHistory[:2])
	}
	gotIDs := []string{req.RequestHistory[2].ID, req.RequestHistory[3].ID, req.RequestHistory[4].ID}
	if gotIDs[0] != "m_bp" || gotIDs[1] != "m_t1" || gotIDs[2] != "m_t2" {
		t.Fatalf("expected tail from breakpoint inclusive, got %#v", gotIDs)
	}
	// The folded huge message must be gone.
	for _, msg := range req.RequestHistory {
		if msg.ID == "m_huge" {
			t.Fatalf("expected folded message dropped")
		}
	}
}

func TestConversationMemory_NoopWhenBreakpointUnknown(t *testing.T) {
	store := &fakeStore{conversationMemories: conversationMemoryItems("c", 5, 8*1024)}
	huge := storage.Message{ID: "m_huge", Role: "user", Content: strings.Repeat("z", 600*1024)}
	last := storage.Message{ID: "m_last", Role: "user", Content: "latest"}
	req := &Request{
		User: storage.User{ID: "u"}, Conversation: storage.Conversation{ID: "c"},
		RequestHistory:               []storage.Message{huge, last},
		DisplayHistory:               []storage.Message{huge, last},
		Store:                        store,
		Estimator:                    DefaultTokenEstimator{},
		ConversationMemoryBreakpoint: "", // unknown → defer to fallback
	}
	if err := (&ConversationMemoryStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if len(req.RequestHistory) != 2 || req.RequestHistory[0].ID != "m_huge" {
		t.Fatalf("expected no-op when breakpoint unknown, got %#v", req.RequestHistory)
	}
}

func TestConversationMemory_NoopWhenBreakpointNotFound(t *testing.T) {
	store := &fakeStore{conversationMemories: conversationMemoryItems("c", 5, 8*1024)}
	huge := storage.Message{ID: "m_huge", Role: "user", Content: strings.Repeat("z", 600*1024)}
	last := storage.Message{ID: "m_last", Role: "user", Content: "latest"}
	req := &Request{
		User: storage.User{ID: "u"}, Conversation: storage.Conversation{ID: "c"},
		RequestHistory:               []storage.Message{huge, last},
		DisplayHistory:               []storage.Message{huge, last},
		Store:                        store,
		Estimator:                    DefaultTokenEstimator{},
		ConversationMemoryBreakpoint: "m_missing", // not in history → defer to fallback
	}
	if err := (&ConversationMemoryStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if len(req.RequestHistory) != 2 {
		t.Fatalf("expected no-op when breakpoint not found, got %#v", req.RequestHistory)
	}
}

func TestConversationMemory_NoopWhenUnfoldedTailOverBudget(t *testing.T) {
	store := &fakeStore{conversationMemories: conversationMemoryItems("c", 5, 8*1024)}
	huge := storage.Message{ID: "m_huge", Role: "user", Content: strings.Repeat("z", 600*1024)}
	// Breakpoint is the huge message itself → unfolded tail still over budget.
	last := storage.Message{ID: "m_last", Role: "user", Content: "latest"}
	req := &Request{
		User: storage.User{ID: "u"}, Conversation: storage.Conversation{ID: "c"},
		RequestHistory:               []storage.Message{huge, last},
		DisplayHistory:               []storage.Message{huge, last},
		Store:                        store,
		Estimator:                    DefaultTokenEstimator{},
		ConversationMemoryBreakpoint: "m_huge",
	}
	if err := (&ConversationMemoryStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if len(req.RequestHistory) != 2 {
		t.Fatalf("expected no-op when unfolded tail over budget, got %#v", req.RequestHistory)
	}
}

func TestSelectUnfoldedTail(t *testing.T) {
	history := []storage.Message{
		{ID: "a"}, {ID: "b"}, {ID: "c"},
	}
	if _, ok := selectUnfoldedTail(history, ""); ok {
		t.Fatalf("expected empty breakpoint to return ok=false")
	}
	if _, ok := selectUnfoldedTail(history, "missing"); ok {
		t.Fatalf("expected missing breakpoint to return ok=false")
	}
	tail, ok := selectUnfoldedTail(history, "b")
	if !ok || len(tail) != 2 || tail[0].ID != "b" || tail[1].ID != "c" {
		t.Fatalf("expected inclusive tail from b, got ok=%v %#v", ok, tail)
	}
}

func TestFullHistorySummarization_KeepsTailFromDisplayHistory(t *testing.T) {
	store := &fakeStore{}
	summarizer := func(ctx context.Context, r SummaryRequest) (SummaryResult, error) {
		return SummaryResult{Summary: "SUMMARY"}, nil
	}
	// RequestHistory (model line) is over budget and uses placeholder-ish content;
	// DisplayHistory holds the verbatim recent messages that must be preserved.
	huge := storage.Message{ID: "r_huge", Role: "user", Content: strings.Repeat("z", 600*1024)}
	reqLast := storage.Message{ID: "r_last", Role: "assistant", Content: "model-line last"}
	d1 := storage.Message{ID: "d1", Role: "user", Content: "display one"}
	d2 := storage.Message{ID: "d2", Role: "assistant", Content: "display two"}
	req := &Request{
		User: storage.User{ID: "u"}, Conversation: storage.Conversation{ID: "c"},
		RequestHistory: []storage.Message{huge, reqLast},
		DisplayHistory: []storage.Message{huge, d1, d2},
		Store:          store,
		Estimator:      DefaultTokenEstimator{},
		Summarizer:     summarizer,
	}
	if err := (&FullHistorySummarizationStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	// Preamble + summary + tail from DisplayHistory (d1,d2); never the model-line r_last.
	final := req.RequestHistory[len(req.RequestHistory)-1]
	if final.ID != "d2" {
		t.Fatalf("expected last kept message from display history, got %#v", final)
	}
	for _, m := range req.RequestHistory {
		if m.ID == "r_last" || m.ID == "r_huge" {
			t.Fatalf("model-line message leaked into kept tail: %#v", m)
		}
	}
}
