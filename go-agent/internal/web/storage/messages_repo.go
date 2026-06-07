package storage

import (
	"context"
	"database/sql"
)

func (s *Store) ListMessagesByConversation(ctx context.Context, conversationID string, limit int) ([]Message, error) {
	if limit <= 0 {
		limit = 100
	}
	conversation, err := s.GetConversationByID(ctx, conversationID)
	if err == nil {
		messages, err := DecodeConversationHistory(conversation.HistoryJSON)
		if err != nil {
			return nil, err
		}
		if len(messages) > limit {
			return messages[:limit], nil
		}
		return messages, nil
	}
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, conversation_id, user_id, role, content, reasoning_content, created_at
		FROM messages WHERE conversation_id = ? ORDER BY created_at ASC LIMIT ?
	`, conversationID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []Message
	for rows.Next() {
		var m Message
		var reasoningContent sql.NullString
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.UserID, &m.Role, &m.Content, &reasoningContent, &m.CreatedAt); err != nil {
			return nil, err
		}
		m.ReasoningContent = reasoningContent.String
		messages = append(messages, m)
	}
	return messages, rows.Err()
}
