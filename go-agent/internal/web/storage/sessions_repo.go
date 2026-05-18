package storage

import (
	"context"
	"fmt"
	"time"
)

func (s *Store) CreateSession(ctx context.Context, session AuthSession) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO auth_sessions (id, user_id, status, expires_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`, session.ID, session.UserID, session.Status, session.ExpiresAt)
	if err != nil {
		return fmt.Errorf("insert session: %w", err)
	}
	return s.Redis.Set(ctx, sessionRedisKey(session.ID), session.UserID, time.Until(session.ExpiresAt)).Err()
}

func (s *Store) GetSession(ctx context.Context, sessionID string) (AuthSession, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, user_id, status, expires_at, revoked_at, created_at, updated_at
		FROM auth_sessions WHERE id = ?
	`, sessionID)
	var session AuthSession
	if err := row.Scan(&session.ID, &session.UserID, &session.Status, &session.ExpiresAt, &session.RevokedAt, &session.CreatedAt, &session.UpdatedAt); err != nil {
		return AuthSession{}, err
	}
	return session, nil
}

func (s *Store) RevokeSession(ctx context.Context, sessionID string) error {
	_, err := s.DB.ExecContext(ctx, `
		UPDATE auth_sessions
		SET status = 'revoked', revoked_at = NOW(), updated_at = NOW()
		WHERE id = ?
	`, sessionID)
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	return s.Redis.Del(ctx, sessionRedisKey(sessionID)).Err()
}
