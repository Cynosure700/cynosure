package runtime

import (
	"strings"
	"testing"

	"cynosure/internal/agent/storage"
)

func TestSessionMemoryProgress_ShouldUpdate(t *testing.T) {
	t.Run("initial extraction gated at 10K", func(t *testing.T) {
		p := &sessionMemoryProgress{}
		if p.shouldUpdate(sessionMemoryInitialTokens-1, false) {
			t.Fatalf("expected no initial extraction below 10K")
		}
		if !p.shouldUpdate(sessionMemoryInitialTokens, false) {
			t.Fatalf("expected initial extraction at 10K")
		}
	})

	t.Run("initial extraction at turn end too", func(t *testing.T) {
		p := &sessionMemoryProgress{}
		if !p.shouldUpdate(sessionMemoryInitialTokens+100, true) {
			t.Fatalf("expected initial extraction at turn end when >=10K")
		}
	})

	t.Run("mid-loop needs growth and tool calls", func(t *testing.T) {
		p := &sessionMemoryProgress{extracted: true, baselineTokens: 20000}
		// growth below threshold
		if p.shouldUpdate(20000+sessionMemoryTokenGrowth-1, false) {
			t.Fatalf("expected no update below growth threshold")
		}
		// growth met but not enough tool calls
		p.toolCallsSinceBase = sessionMemoryToolCallsMin - 1
		if p.shouldUpdate(20000+sessionMemoryTokenGrowth, false) {
			t.Fatalf("expected no mid-loop update with too few tool calls")
		}
		// growth + tool calls met
		p.toolCallsSinceBase = sessionMemoryToolCallsMin
		if !p.shouldUpdate(20000+sessionMemoryTokenGrowth, false) {
			t.Fatalf("expected mid-loop update with growth + tool calls")
		}
	})

	t.Run("turn end needs only growth", func(t *testing.T) {
		p := &sessionMemoryProgress{extracted: true, baselineTokens: 20000}
		if p.shouldUpdate(20000+sessionMemoryTokenGrowth-1, true) {
			t.Fatalf("expected no turn-end update below growth threshold")
		}
		p.toolCallsSinceBase = 0
		if !p.shouldUpdate(20000+sessionMemoryTokenGrowth, true) {
			t.Fatalf("expected turn-end update on growth alone (no tool calls needed)")
		}
	})

	t.Run("baseline drops to compressed low point", func(t *testing.T) {
		p := &sessionMemoryProgress{extracted: true, baselineTokens: 50000}
		// context fell to 30000 after compression; should not update yet but baseline lowers.
		if p.shouldUpdate(30000, true) {
			t.Fatalf("expected no update right after baseline drop")
		}
		if p.baselineTokens != 30000 {
			t.Fatalf("expected baseline lowered to 30000, got %d", p.baselineTokens)
		}
		// now growth measured from the new low point.
		if !p.shouldUpdate(30000+sessionMemoryTokenGrowth, true) {
			t.Fatalf("expected update once growth from new low point is met")
		}
	})
}

func TestRenderModelHistoryForMemory_IncludesToolCallsAndResults(t *testing.T) {
	history := []storage.Message{
		{Role: "user", Content: "fix the bug"},
		{Role: "assistant", Content: "looking", ToolCalls: []storage.MessageToolCall{
			{ID: "c1", Type: "function", Function: storage.MessageFunctionCall{Name: "bash", Arguments: `{"command":"go test"}`}},
		}},
		{Role: "tool", ToolCallID: "c1", Content: `{"status":"success","result":"ok PASS"}`},
		{Role: "tool", ToolCallID: "c2", Content: "raw not json"},
	}
	out := renderModelHistoryForMemory(history)
	for _, want := range []string{"[user] fix the bug", "[assistant] looking", "[tool_call] bash(", "go test", "[tool_result] success: ok PASS", "[tool_result] raw not json"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected transcript to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderModelHistoryForMemory_Truncates(t *testing.T) {
	var history []storage.Message
	for i := 0; i < 5000; i++ {
		history = append(history, storage.Message{Role: "user", Content: "padding line"})
	}
	out := renderModelHistoryForMemory(history)
	if len(out) > maxMemoryTranscriptChars+200 {
		t.Fatalf("expected transcript capped near %d, got %d", maxMemoryTranscriptChars, len(out))
	}
}
