package storage

import "testing"

func TestConversationHistoryJSONRoundTrip(t *testing.T) {
	messages := []Message{
		{ID: "msg_user", ConversationID: "conv_1", UserID: "usr_1", Role: "user", Content: "hello"},
		{ID: "msg_assistant", ConversationID: "conv_1", UserID: "usr_1", Role: "assistant", Content: "hi"},
	}

	encoded, err := EncodeConversationHistory(messages)
	if err != nil {
		t.Fatalf("encode conversation history: %v", err)
	}
	decoded, err := DecodeConversationHistory(encoded)
	if err != nil {
		t.Fatalf("decode conversation history: %v", err)
	}

	if len(decoded) != 2 {
		t.Fatalf("expected 2 decoded messages, got %d", len(decoded))
	}
	if decoded[0].ID != "msg_user" || decoded[0].Content != "hello" || decoded[1].Role != "assistant" {
		t.Fatalf("unexpected decoded messages: %#v", decoded)
	}
}

func TestDecodeConversationHistoryRejectsInvalidJSON(t *testing.T) {
	_, err := DecodeConversationHistory(`{"messages":`)
	if err == nil {
		t.Fatalf("expected invalid history json to return an error")
	}
}
