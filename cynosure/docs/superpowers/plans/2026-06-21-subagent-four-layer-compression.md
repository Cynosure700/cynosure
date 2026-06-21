# Subagent Four-Layer Compression Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add the same four request-context compression layers to subagents without changing parent-agent memory or existing runtime behavior.

**Architecture:** Introduce a subagent-specific compressor chain that reuses the existing compression strategies except `ConversationMemoryStrategy`. Wire `runSubagentLoop` to compress `state.ModelHistory` before each LLM request, write the result back, and rebuild OpenAI messages from the compressed model history.

**Tech Stack:** Go 1.26.1, `go-openai`, existing `internal/agent/runtime/compression` package, existing runtime fake store tests.

## Global Constraints

- Subagents must keep fresh message-list semantics and must not read or inject parent conversation history.
- Subagents use exactly four compression layers: `ToolResultCompressionStrategy`, `MessageWindowCompressionStrategy`, `RecentToolResultRetentionStrategy`, `FullHistorySummarizationStrategy`.
- Do not change main-agent `NewDefaultCompressor()` behavior.
- Use the child `ToolRegistry` when resolving tool definitions and per-tool result limits.
- Compression must not mutate display history `state.History`.
- Do not change approval, timeout, trace, tool event, or nested-subagent behavior.

---

### Task 1: Add Subagent Compressor and Compression Entry

**Files:**
- Modify: `internal/agent/runtime/compression/compression.go`
- Modify: `internal/agent/runtime/context_compression.go`
- Test: `internal/agent/runtime/context_compression_test.go`

**Interfaces:**
- Produces: `compression.NewSubagentCompressor() *compression.Compressor`
- Produces: `func (s *Service) compressSubagentContextBeforeLLM(ctx context.Context, state *LoopState, tools *ToolRegistry) ([]storage.Message, error)`

- [ ] **Step 1: Write the failing tests**

Add tests in `internal/agent/runtime/context_compression_test.go`:

```go
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
		Type:           "fact",
		Name:           "parent-memory",
		Content:        "parent memory must not enter child context",
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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go -C cynosure test ./internal/agent/runtime -run 'TestCompressSubagentContextBeforeLLM' -count=1`

Expected: FAIL because `compressSubagentContextBeforeLLM` is undefined.

- [ ] **Step 3: Add minimal implementation**

Add `NewSubagentCompressor()` to `compression.go` and add `compressSubagentContextBeforeLLM` to `context_compression.go`, mirroring `compressContextBeforeLLM` while using `compression.NewSubagentCompressor()` and the passed child registry.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go -C cynosure test ./internal/agent/runtime -run 'TestCompressSubagentContextBeforeLLM' -count=1`

Expected: PASS.

### Task 2: Wire Compression Into Subagent Loop

**Files:**
- Modify: `internal/agent/runtime/subagent.go`
- Test: `internal/agent/runtime/runtime_test.go`

**Interfaces:**
- Consumes: `compressSubagentContextBeforeLLM(ctx, state, tools)`

- [ ] **Step 1: Write the failing test**

Add a runtime test proving the subagent request sent after a large tool result receives a compressed marker instead of the full result.

- [ ] **Step 2: Run test to verify it fails**

Run: `go -C cynosure test ./internal/agent/runtime -run 'TestRespondToConversation_SubagentCompressesContextBeforeNextRound' -count=1`

Expected: FAIL because the second subagent LLM request still contains the full result.

- [ ] **Step 3: Wire compression in `runSubagentLoop`**

Before constructing the OpenAI request each round, call the subagent compression entry, write the result to `state.ModelHistory`, rebuild `state.Messages`, and keep existing reminder, request, logging, approval, and tool execution behavior.

- [ ] **Step 4: Run focused runtime tests**

Run: `go -C cynosure test ./internal/agent/runtime -run 'TestRespondToConversation_(SpawnSubagentUsesFreshMessagesStoresTraceAndDoesNotEmitToolEvents|SubagentCompressesContextBeforeNextRound)|TestCompressSubagentContextBeforeLLM' -count=1`

Expected: PASS.

### Task 3: Final Verification

**Files:**
- No additional files.

- [ ] **Step 1: Format Go code**

Run: `go -C cynosure fmt ./internal/agent/runtime ./internal/agent/runtime/compression`

- [ ] **Step 2: Run compression tests**

Run: `go -C cynosure test ./internal/agent/runtime/compression -count=1`

Expected: PASS.

- [ ] **Step 3: Run runtime tests**

Run: `go -C cynosure test ./internal/agent/runtime -count=1`

Expected: PASS or only known pre-existing unrelated failures must be documented with exact output.

- [ ] **Step 4: Run build**

Run: `go -C cynosure build ./...`

Expected: exit 0.
