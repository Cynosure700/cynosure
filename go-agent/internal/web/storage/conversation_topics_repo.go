package storage

import (
	"context"
	"database/sql"
)

func (s *Store) UpsertConversationTopics(ctx context.Context, topics ConversationTopics) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO conversation_topics (conversation_id, user_id, topics_json, updated_at)
		VALUES (?, ?, ?, NOW(6))
		ON DUPLICATE KEY UPDATE topics_json = VALUES(topics_json)
	`, topics.ConversationID, topics.UserID, topics.TopicsJSON)
	return err
}

func (s *Store) ListRecentTopicsByUser(ctx context.Context, userID, excludeConversationID string, limit int) ([]ConversationTopics, error) {
	if limit <= 0 {
		limit = 5
	}
	rows, err := s.DB.QueryContext(ctx, `
		SELECT conversation_id, user_id, topics_json, updated_at
		FROM conversation_topics
		WHERE user_id = ? AND conversation_id <> ?
		ORDER BY updated_at DESC
		LIMIT ?
	`, userID, excludeConversationID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanConversationTopics(rows)
}

func scanConversationTopics(rows *sql.Rows) ([]ConversationTopics, error) {
	var result []ConversationTopics
	for rows.Next() {
		var t ConversationTopics
		if err := rows.Scan(&t.ConversationID, &t.UserID, &t.TopicsJSON, &t.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, rows.Err()
}
