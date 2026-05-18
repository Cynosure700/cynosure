package storage

import (
	"context"
	"fmt"
)

func (s *Store) CreateConversation(ctx context.Context, conversation Conversation) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO conversations (id, user_id, title, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
	`, conversation.ID, conversation.UserID, conversation.Title)
	return err
}

func (s *Store) ListConversationsByUser(ctx context.Context, userID string) ([]Conversation, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, user_id, title, created_at, updated_at
		FROM conversations WHERE user_id = ? ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var conversations []Conversation
	for rows.Next() {
		var c Conversation
		if err := rows.Scan(&c.ID, &c.UserID, &c.Title, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		conversations = append(conversations, c)
	}
	return conversations, rows.Err()
}

func (s *Store) GetConversationByID(ctx context.Context, conversationID string) (Conversation, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, user_id, title, created_at, updated_at
		FROM conversations WHERE id = ?
	`, conversationID)
	var c Conversation
	if err := row.Scan(&c.ID, &c.UserID, &c.Title, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return Conversation{}, err
	}
	return c, nil
}

func (s *Store) TouchConversation(ctx context.Context, conversationID, title string) error {
	_, err := s.DB.ExecContext(ctx, `
		UPDATE conversations SET title = COALESCE(NULLIF(?, ''), title), updated_at = NOW() WHERE id = ?
	`, title, conversationID)
	return err
}

func (s *Store) UpdateConversationTitle(ctx context.Context, conversationID, title string) error {
	_, err := s.DB.ExecContext(ctx, `
		UPDATE conversations SET title = ?, updated_at = NOW() WHERE id = ?
	`, title, conversationID)
	return err
}

func (s *Store) TouchConversationActivity(ctx context.Context, conversationID string) error {
	_, err := s.DB.ExecContext(ctx, `
		UPDATE conversations SET updated_at = NOW() WHERE id = ?
	`, conversationID)
	return err
}

func (s *Store) DeleteConversation(ctx context.Context, conversationID string) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM conversations WHERE id = ?`, conversationID)
	if err != nil {
		return fmt.Errorf("delete conversation: %w", err)
	}
	if err := s.Redis.Del(ctx, conversationCacheKey(conversationID)).Err(); err != nil {
		return fmt.Errorf("clear conversation cache: %w", err)
	}
	return nil
}
