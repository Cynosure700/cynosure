package storage

import (
	"context"
	"database/sql"
)

func (s *Store) CreatePersistedOutput(ctx context.Context, output PersistedOutput) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO persisted_outputs (id, conversation_id, user_id, message_id, tool_call_id, kind, strategy, original_bytes, content_sha256, content, preview, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(6))
	`, output.ID, output.ConversationID, output.UserID, output.MessageID, output.ToolCallID, output.Kind, output.Strategy, output.OriginalBytes, output.ContentSHA256, output.Content, output.Preview)
	return err
}

func (s *Store) GetPersistedOutputForConversation(ctx context.Context, id, userID, conversationID string) (PersistedOutput, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, conversation_id, user_id, message_id, tool_call_id, kind, strategy, original_bytes, content_sha256, content, preview, created_at
		FROM persisted_outputs WHERE id = ? AND user_id = ? AND conversation_id = ?
	`, id, userID, conversationID)
	return scanPersistedOutput(row)
}

func (s *Store) GetPersistedOutputByMessageHash(ctx context.Context, conversationID, userID, messageID, toolCallID, strategy, contentSHA256 string) (PersistedOutput, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, conversation_id, user_id, message_id, tool_call_id, kind, strategy, original_bytes, content_sha256, content, preview, created_at
		FROM persisted_outputs
		WHERE conversation_id = ? AND user_id = ? AND message_id = ? AND tool_call_id = ? AND strategy = ? AND content_sha256 = ?
	`, conversationID, userID, messageID, toolCallID, strategy, contentSHA256)
	return scanPersistedOutput(row)
}

func scanPersistedOutput(row *sql.Row) (PersistedOutput, error) {
	var o PersistedOutput
	if err := row.Scan(&o.ID, &o.ConversationID, &o.UserID, &o.MessageID, &o.ToolCallID, &o.Kind, &o.Strategy, &o.OriginalBytes, &o.ContentSHA256, &o.Content, &o.Preview, &o.CreatedAt); err != nil {
		return PersistedOutput{}, err
	}
	return o, nil
}

func (s *Store) CreateContextSummary(ctx context.Context, summary ContextSummary) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO context_summaries (id, conversation_id, user_id, source_history_sha256, strategy, estimated_tokens_before, estimated_tokens_after, summary, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(6))
	`, summary.ID, summary.ConversationID, summary.UserID, summary.SourceHistorySHA256, summary.Strategy, summary.EstimatedTokensBefore, summary.EstimatedTokensAfter, summary.Summary)
	return err
}

func (s *Store) GetContextSummaryByHistoryHash(ctx context.Context, conversationID, userID, sourceHistorySHA256 string) (ContextSummary, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, conversation_id, user_id, source_history_sha256, strategy, estimated_tokens_before, estimated_tokens_after, summary, created_at
		FROM context_summaries
		WHERE conversation_id = ? AND user_id = ? AND source_history_sha256 = ?
		ORDER BY created_at DESC LIMIT 1
	`, conversationID, userID, sourceHistorySHA256)
	var c ContextSummary
	if err := row.Scan(&c.ID, &c.ConversationID, &c.UserID, &c.SourceHistorySHA256, &c.Strategy, &c.EstimatedTokensBefore, &c.EstimatedTokensAfter, &c.Summary, &c.CreatedAt); err != nil {
		return ContextSummary{}, err
	}
	return c, nil
}
