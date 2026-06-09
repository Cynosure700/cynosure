package storage

import (
	"context"
	"database/sql"
)

// GetConversationModelHistory returns the persisted "model history" (the
// previous turn's compressed request history) for the given conversation. The
// bool is false when no row exists or the stored JSON fails to decode, so the
// caller can fall back to the full display history.
func (s *Store) GetConversationModelHistory(ctx context.Context, conversationID string) ([]Message, bool, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT history_json FROM conversation_model_histories WHERE conversation_id = ?
	`, conversationID)
	var historyJSON string
	if err := row.Scan(&historyJSON); err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, err
	}
	messages, err := DecodeConversationHistory(historyJSON)
	if err != nil {
		return nil, false, nil
	}
	if len(messages) == 0 {
		return nil, false, nil
	}
	return messages, true, nil
}

// UpsertConversationModelHistory inserts or overwrites the model history for the
// given conversation (one row per conversation).
func (s *Store) UpsertConversationModelHistory(ctx context.Context, conversationID, userID string, messages []Message) error {
	historyJSON, err := EncodeConversationHistory(messages)
	if err != nil {
		return err
	}
	_, err = s.DB.ExecContext(ctx, `
		INSERT INTO conversation_model_histories (conversation_id, user_id, history_json, estimated_tokens, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(6), NOW(6))
		ON DUPLICATE KEY UPDATE
			user_id = VALUES(user_id),
			history_json = VALUES(history_json),
			estimated_tokens = VALUES(estimated_tokens),
			updated_at = NOW(6)
	`, conversationID, userID, historyJSON, len(historyJSON))
	return err
}
