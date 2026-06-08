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
	if err := s.ensureContextCompressionTables(ctx); err != nil {
		return fmt.Errorf("migrate context compression tables: %w", err)
	}
	if err := s.ensureMemoryTables(ctx); err != nil {
		return fmt.Errorf("migrate memory tables: %w", err)
	}
	if err := s.ensureMemoriesTable(ctx); err != nil {
		return fmt.Errorf("migrate memories table: %w", err)
	}
	return nil
}

func (s *Store) ensureMemoriesTable(ctx context.Context) error {
	if _, err := s.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS memories (
			id VARCHAR(64) PRIMARY KEY,
			user_id VARCHAR(64) NULL,
			type VARCHAR(32) NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			body LONGTEXT NOT NULL,
			created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
			updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
			KEY idx_memories_user_type (user_id, type, updated_at),
			KEY idx_memories_type (type, updated_at),
			CONSTRAINT fk_memories_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`); err != nil {
		return fmt.Errorf("create memories: %w", err)
	}
	return nil
}

func (s *Store) ensureMemoryTables(ctx context.Context) error {
	if _, err := s.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS user_profiles (
			user_id VARCHAR(64) PRIMARY KEY,
			profile_json LONGTEXT NOT NULL,
			updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
			CONSTRAINT fk_user_profiles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`); err != nil {
		return fmt.Errorf("create user_profiles: %w", err)
	}
	if _, err := s.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS conversation_topics (
			conversation_id VARCHAR(64) PRIMARY KEY,
			user_id VARCHAR(64) NOT NULL,
			topics_json LONGTEXT NOT NULL,
			updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
			KEY idx_conversation_topics_user (user_id, updated_at),
			CONSTRAINT fk_conversation_topics_conv FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
			CONSTRAINT fk_conversation_topics_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`); err != nil {
		return fmt.Errorf("create conversation_topics: %w", err)
	}
	return nil
}

func (s *Store) ensureContextCompressionTables(ctx context.Context) error {
	if _, err := s.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS persisted_outputs (
			id VARCHAR(64) PRIMARY KEY,
			conversation_id VARCHAR(64) NOT NULL,
			user_id VARCHAR(64) NOT NULL,
			message_id VARCHAR(64) NOT NULL,
			tool_call_id VARCHAR(128) NOT NULL DEFAULT '',
			kind VARCHAR(64) NOT NULL,
			strategy VARCHAR(128) NOT NULL,
			original_bytes INT NOT NULL,
			content_sha256 CHAR(64) NOT NULL,
			content LONGTEXT NOT NULL,
			preview TEXT NOT NULL,
			created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
			UNIQUE KEY uniq_po_message_hash (conversation_id, user_id, message_id, tool_call_id, strategy, content_sha256),
			KEY idx_po_conversation (conversation_id, created_at),
			CONSTRAINT fk_po_conversation FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
			CONSTRAINT fk_po_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`); err != nil {
		return fmt.Errorf("create persisted_outputs: %w", err)
	}
	if _, err := s.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS context_summaries (
			id VARCHAR(64) PRIMARY KEY,
			conversation_id VARCHAR(64) NOT NULL,
			user_id VARCHAR(64) NOT NULL,
			source_history_sha256 CHAR(64) NOT NULL,
			strategy VARCHAR(128) NOT NULL,
			estimated_tokens_before INT NOT NULL,
			estimated_tokens_after INT NOT NULL,
			summary LONGTEXT NOT NULL,
			created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
			KEY idx_cs_lookup (conversation_id, user_id, source_history_sha256),
			CONSTRAINT fk_cs_conversation FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
			CONSTRAINT fk_cs_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`); err != nil {
		return fmt.Errorf("create context_summaries: %w", err)
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
