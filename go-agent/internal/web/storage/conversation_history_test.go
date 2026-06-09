package storage

import (
	"strings"
	"testing"
)

func TestConversationHistoryJSONRoundTrip(t *testing.T) {
	messages := []Message{
		{ID: "msg_user", ConversationID: "conv_1", UserID: "usr_1", Role: "user", Content: "hello"},
		{ID: "msg_assistant", ConversationID: "conv_1", UserID: "usr_1", Role: "assistant", Content: "hi", ReasoningContent: "thinking", ToolCalls: []MessageToolCall{{ID: "tool_1", Type: "function", Function: MessageFunctionCall{Name: "load_skill", Arguments: `{"name":"builtin-skill"}`}}}},
		{ID: "msg_tool", ConversationID: "conv_1", UserID: "usr_1", Role: "tool", Content: `{"status":"success","result":"loaded"}`, ToolCallID: "tool_1"},
	}

	encoded, err := EncodeConversationHistory(messages)
	if err != nil {
		t.Fatalf("encode conversation history: %v", err)
	}
	decoded, err := DecodeConversationHistory(encoded)
	if err != nil {
		t.Fatalf("decode conversation history: %v", err)
	}

	if len(decoded) != 3 {
		t.Fatalf("expected 3 decoded messages, got %d", len(decoded))
	}
	if decoded[0].ID != "msg_user" || decoded[0].Content != "hello" || decoded[1].Role != "assistant" || decoded[1].ReasoningContent != "thinking" {
		t.Fatalf("unexpected decoded messages: %#v", decoded)
	}
	if len(decoded[1].ToolCalls) != 1 || decoded[1].ToolCalls[0].ID != "tool_1" || decoded[2].Role != "tool" || decoded[2].ToolCallID != "tool_1" {
		t.Fatalf("expected decoded tool call and result messages, got %#v", decoded)
	}
}

func TestDecodeConversationHistoryRejectsInvalidJSON(t *testing.T) {
	_, err := DecodeConversationHistory(`{"messages":`)
	if err == nil {
		t.Fatalf("expected invalid history json to return an error")
	}
}

func TestMessageMetaToolCallCountZeroIsSerialized(t *testing.T) {
	messages := []Message{
		{ID: "msg_assistant", Role: "assistant", Content: "纯文本回复", Meta: &MessageMeta{ToolCallCount: 0, ContextTokens: 50, ContextBudget: 100}},
	}

	encoded, err := EncodeConversationHistory(messages)
	if err != nil {
		t.Fatalf("encode conversation history: %v", err)
	}
	if !strings.Contains(encoded, `"tool_call_count":0`) {
		t.Fatalf("expected tool_call_count to be serialized even when zero, got %s", encoded)
	}

	decoded, err := DecodeConversationHistory(encoded)
	if err != nil {
		t.Fatalf("decode conversation history: %v", err)
	}
	if decoded[0].Meta == nil || decoded[0].Meta.ToolCallCount != 0 {
		t.Fatalf("expected decoded meta with zero tool call count, got %#v", decoded[0].Meta)
	}
}
