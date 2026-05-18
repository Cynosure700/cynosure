package storage

import "context"

func (s *Store) CreateMessage(ctx context.Context, message Message) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO messages (id, conversation_id, user_id, role, content, created_at)
		VALUES (?, ?, ?, ?, ?, NOW())
	`, message.ID, message.ConversationID, message.UserID, message.Role, message.Content)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) ListMessagesByConversation(ctx context.Context, conversationID string, limit int) ([]Message, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, conversation_id, user_id, role, content, created_at
		FROM messages WHERE conversation_id = ? ORDER BY created_at ASC LIMIT ?
	`, conversationID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.UserID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}
