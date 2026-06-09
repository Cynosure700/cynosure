package storage

import (
	"context"
	"database/sql"
)

const conversationMemoryColumns = "id, conversation_id, user_id, name, description, body, created_at, updated_at"

func scanConversationMemories(rows *sql.Rows) ([]ConversationMemory, error) {
	var result []ConversationMemory
	for rows.Next() {
		var m ConversationMemory
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.UserID, &m.Name, &m.Description, &m.Body, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, rows.Err()
}

// ListConversationMemories returns the conversation memory entries of the given
// conversation, oldest first.
func (s *Store) ListConversationMemories(ctx context.Context, conversationID string) ([]ConversationMemory, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT `+conversationMemoryColumns+`
		FROM conversation_memories
		WHERE conversation_id = ?
		ORDER BY updated_at ASC
	`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanConversationMemories(rows)
}

// ReplaceConversationMemories deletes all memory entries of the given
// conversation, then inserts the provided items in a single transaction.
func (s *Store) ReplaceConversationMemories(ctx context.Context, conversationID, userID string, items []ConversationMemory) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `DELETE FROM conversation_memories WHERE conversation_id = ?`, conversationID); err != nil {
		return err
	}
	for _, m := range items {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO conversation_memories (id, conversation_id, user_id, name, description, body, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, NOW(6), NOW(6))
		`, m.ID, conversationID, userID, m.Name, m.Description, m.Body); err != nil {
			return err
		}
	}
	return tx.Commit()
}
