package storage

import (
	"context"
	"fmt"
)

func (s *Store) RunMigrations(ctx context.Context) error {
	data, err := migrationFiles.ReadFile("migrations/001_init.sql")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	if _, err := s.DB.ExecContext(ctx, string(data)); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
