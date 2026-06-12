package storage

import (
	"encoding/json"
	"strings"
)

const emptyConversationHistoryJSON = "[]"

func EncodeConversationHistory(messages []Message) (string, error) {
	if messages == nil {
		messages = []Message{}
	}
	data, err := json.Marshal(messages)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func DecodeConversationHistory(historyJSON string) ([]Message, error) {
	if strings.TrimSpace(historyJSON) == "" {
		historyJSON = emptyConversationHistoryJSON
	}
	var messages []Message
	if err := json.Unmarshal([]byte(historyJSON), &messages); err != nil {
		return nil, err
	}
	if messages == nil {
		messages = []Message{}
	}
	return messages, nil
}
