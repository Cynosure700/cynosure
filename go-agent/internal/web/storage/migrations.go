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
	if err := s.ensureConversationHistoryColumns(ctx); err != nil {
		return fmt.Errorf("migrate conversation history columns: %w", err)
	}
	if err := s.ensureMessageReasoningContentColumn(ctx); err != nil {
		return fmt.Errorf("migrate message reasoning content column: %w", err)
	}
	if err := s.ensureSubagentMessagesTable(ctx); err != nil {
		return fmt.Errorf("migrate subagent messages table: %w", err)
	}
	return nil
}

func (s *Store) ensureConversationHistoryColumns(ctx context.Context) error {
	hasRootMessageID, err := s.columnExists(ctx, "conversations", "root_message_id")
	if err != nil {
		return err
	}
	if !hasRootMessageID {
		if _, err := s.DB.ExecContext(ctx, `ALTER TABLE conversations ADD COLUMN root_message_id VARCHAR(64) NOT NULL DEFAULT '' AFTER user_id`); err != nil {
			return fmt.Errorf("add root_message_id: %w", err)
		}
	}

	hasHistoryJSON, err := s.columnExists(ctx, "conversations", "history_json")
	if err != nil {
		return err
	}
	if !hasHistoryJSON {
		if _, err := s.DB.ExecContext(ctx, `ALTER TABLE conversations ADD COLUMN history_json LONGTEXT NULL AFTER title`); err != nil {
			return fmt.Errorf("add history_json: %w", err)
		}
	}

	if _, err := s.DB.ExecContext(ctx, `UPDATE conversations SET root_message_id = id WHERE root_message_id = ''`); err != nil {
		return fmt.Errorf("backfill root_message_id: %w", err)
	}
	if _, err := s.DB.ExecContext(ctx, `UPDATE conversations SET history_json = '[]' WHERE history_json IS NULL OR history_json = ''`); err != nil {
		return fmt.Errorf("backfill history_json: %w", err)
	}
	if _, err := s.DB.ExecContext(ctx, `ALTER TABLE conversations MODIFY history_json LONGTEXT NOT NULL`); err != nil {
		return fmt.Errorf("require history_json: %w", err)
	}

	hasRootMessageIndex, err := s.indexExists(ctx, "conversations", "uniq_conversations_user_root_message")
	if err != nil {
		return err
	}
	if !hasRootMessageIndex {
		if _, err := s.DB.ExecContext(ctx, `ALTER TABLE conversations ADD UNIQUE KEY uniq_conversations_user_root_message (user_id, root_message_id)`); err != nil {
			return fmt.Errorf("add conversation root message index: %w", err)
		}
	}
	return nil
}

func (s *Store) ensureMessageReasoningContentColumn(ctx context.Context) error {
	hasReasoningContent, err := s.columnExists(ctx, "messages", "reasoning_content")
	if err != nil {
		return err
	}
	if hasReasoningContent {
		return nil
	}
	if _, err := s.DB.ExecContext(ctx, `ALTER TABLE messages ADD COLUMN reasoning_content LONGTEXT NULL AFTER content`); err != nil {
		return fmt.Errorf("add reasoning_content: %w", err)
	}
	return nil
}

func (s *Store) ensureSubagentMessagesTable(ctx context.Context) error {
	_, err := s.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS subagent_messages (
			id VARCHAR(64) PRIMARY KEY,
			run_id VARCHAR(64) NOT NULL,
			parent_tool_call_id VARCHAR(128) NOT NULL,
			conversation_id VARCHAR(64) NOT NULL,
			user_id VARCHAR(64) NOT NULL,
			sequence_no INT NOT NULL,
			role VARCHAR(32) NOT NULL,
			content LONGTEXT NOT NULL,
			reasoning_content LONGTEXT NULL,
			tool_call_id VARCHAR(128) NOT NULL DEFAULT '',
			tool_calls_json LONGTEXT NOT NULL,
			created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
			UNIQUE KEY uniq_subagent_messages_run_sequence (run_id, sequence_no),
			KEY idx_subagent_messages_conversation (conversation_id, created_at),
			CONSTRAINT fk_subagent_messages_conversation FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
			CONSTRAINT fk_subagent_messages_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`)
	return err
}

func (s *Store) columnExists(ctx context.Context, tableName, columnName string) (bool, error) {
	var count int
	err := s.DB.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = ?
	`, tableName, columnName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("check column %s.%s: %w", tableName, columnName, err)
	}
	return count > 0, nil
}

func (s *Store) indexExists(ctx context.Context, tableName, indexName string) (bool, error) {
	var count int
	err := s.DB.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND INDEX_NAME = ?
	`, tableName, indexName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("check index %s.%s: %w", tableName, indexName, err)
	}
	return count > 0, nil
}
