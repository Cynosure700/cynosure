package storage

import (
	"context"
	"fmt"
	"strings"
)

func (s *Store) CreateUser(ctx context.Context, user User) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO users (id, email, username, password_hash, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`, user.ID, strings.ToLower(user.Email), user.Username, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (s *Store) GetUserByEmailOrUsername(ctx context.Context, login string) (User, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users
		WHERE LOWER(email) = LOWER(?) OR username = ?
	`, login, login)
	return scanUser(row)
}

func (s *Store) GetUserByID(ctx context.Context, userID string) (User, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users WHERE id = ?
	`, userID)
	return scanUser(row)
}
