package storage

import (
	"context"
	"database/sql"
	"strings"
)

const memoryColumns = "id, user_id, type, name, description, body, created_at, updated_at"

func scanMemories(rows *sql.Rows) ([]Memory, error) {
	var result []Memory
	for rows.Next() {
		var m Memory
		var userID sql.NullString
		if err := rows.Scan(&m.ID, &userID, &m.Type, &m.Name, &m.Description, &m.Body, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		if userID.Valid {
			m.UserID = userID.String
		}
		result = append(result, m)
	}
	return result, rows.Err()
}

func nullableUserID(userID string) any {
	if strings.TrimSpace(userID) == "" {
		return nil
	}
	return userID
}

// ListRelevantMemories returns the user's own memories (episodic_memory and
// user_preference) together with all system-level semantic memories.
func (s *Store) ListRelevantMemories(ctx context.Context, userID string) ([]Memory, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT `+memoryColumns+`
		FROM memories
		WHERE (user_id = ? AND type IN ('episodic_memory', 'user_preference'))
		   OR (user_id IS NULL AND type = 'semantic')
		ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMemories(rows)
}

func (s *Store) ListMemoriesByUserAndType(ctx context.Context, userID, memType string) ([]Memory, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT `+memoryColumns+`
		FROM memories
		WHERE user_id = ? AND type = ?
		ORDER BY updated_at DESC
	`, userID, memType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMemories(rows)
}

func (s *Store) ListSemanticMemories(ctx context.Context) ([]Memory, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT `+memoryColumns+`
		FROM memories
		WHERE user_id IS NULL AND type = 'semantic'
		ORDER BY updated_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMemories(rows)
}

func (s *Store) InsertMemory(ctx context.Context, m Memory) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO memories (id, user_id, type, name, description, body, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW(6), NOW(6))
	`, m.ID, nullableUserID(m.UserID), m.Type, m.Name, m.Description, m.Body)
	return err
}

func (s *Store) CountMemoriesByUserAndType(ctx context.Context, userID, memType string) (int, error) {
	var count int
	err := s.DB.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM memories WHERE user_id = ? AND type = ?
	`, userID, memType).Scan(&count)
	return count, err
}

func (s *Store) CountSemanticMemories(ctx context.Context) (int, error) {
	var count int
	err := s.DB.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM memories WHERE user_id IS NULL AND type = 'semantic'
	`).Scan(&count)
	return count, err
}

// DeleteOldestMemories removes the n oldest memories (by updated_at ASC) of the
// given user and type. The inner subquery is wrapped so MySQL allows LIMIT in a
// DELETE ... WHERE id IN (...) against the same table.
func (s *Store) DeleteOldestMemories(ctx context.Context, userID, memType string, n int) error {
	if n <= 0 {
		return nil
	}
	_, err := s.DB.ExecContext(ctx, `
		DELETE FROM memories WHERE id IN (
			SELECT id FROM (
				SELECT id FROM memories
				WHERE user_id = ? AND type = ?
				ORDER BY updated_at ASC
				LIMIT ?
			) AS t
		)
	`, userID, memType, n)
	return err
}

// ReplaceMemoriesByUserAndType deletes all memories of the given user and type,
// then inserts the provided items in a single transaction.
func (s *Store) ReplaceMemoriesByUserAndType(ctx context.Context, userID, memType string, items []Memory) error {
	return s.replaceMemories(ctx, `DELETE FROM memories WHERE user_id = ? AND type = ?`, []any{userID, memType}, userID, memType, items)
}

// ReplaceSemanticMemories deletes all system-level semantic memories, then
// inserts the provided items in a single transaction.
func (s *Store) ReplaceSemanticMemories(ctx context.Context, items []Memory) error {
	return s.replaceMemories(ctx, `DELETE FROM memories WHERE user_id IS NULL AND type = 'semantic'`, nil, "", MemoryTypeSemanticValue, items)
}

// MemoryTypeSemanticValue mirrors the runtime semantic type constant so the
// storage layer can normalize replaced semantic memories without importing the
// runtime package.
const MemoryTypeSemanticValue = "semantic"

func (s *Store) replaceMemories(ctx context.Context, deleteSQL string, deleteArgs []any, userID, memType string, items []Memory) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, deleteSQL, deleteArgs...); err != nil {
		return err
	}
	for _, m := range items {
		uid := userID
		if memType == MemoryTypeSemanticValue {
			uid = ""
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO memories (id, user_id, type, name, description, body, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, NOW(6), NOW(6))
		`, m.ID, nullableUserID(uid), memType, m.Name, m.Description, m.Body); err != nil {
			return err
		}
	}
	return tx.Commit()
}
