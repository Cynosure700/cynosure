package tools

import (
	"context"
	"encoding/json"
	"fmt"
)

// PersistedOutput is the runtime-agnostic view of a persisted tool output that
// the read_persisted_output handler returns to the model.
type PersistedOutput struct {
	ID            string
	Kind          string
	ToolCallID    string
	OriginalBytes int
	ContentSHA256 string
	Content       string
}

// PersistedOutputReader retrieves a persisted output for the current
// conversation. Implementations must enforce conversation/user scoping.
type PersistedOutputReader interface {
	ReadPersistedOutput(ctx context.Context, id string) (PersistedOutput, error)
}

type persistedOutputReaderContextKey struct{}

func WithPersistedOutputReader(ctx context.Context, reader PersistedOutputReader) context.Context {
	return context.WithValue(ctx, persistedOutputReaderContextKey{}, reader)
}

func PersistedOutputReaderFromContext(ctx context.Context) (PersistedOutputReader, bool) {
	reader, ok := ctx.Value(persistedOutputReaderContextKey{}).(PersistedOutputReader)
	return reader, ok && reader != nil
}

const (
	defaultPersistedOutputReadLimit = 20000
	maxPersistedOutputReadLimit     = 100000
)

func handleReadPersistedOutput(ctx context.Context, args map[string]any) (string, error) {
	id, _ := args["id"].(string)
	if id == "" {
		return "", fmt.Errorf("id is required")
	}
	offset := 0
	if v, ok := args["offset"].(float64); ok && int(v) > 0 {
		offset = int(v)
	}
	limit := defaultPersistedOutputReadLimit
	if v, ok := args["limit"].(float64); ok && int(v) > 0 {
		limit = int(v)
	}
	if limit > maxPersistedOutputReadLimit {
		limit = maxPersistedOutputReadLimit
	}

	reader, ok := PersistedOutputReaderFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("read_persisted_output is not available in this context")
	}
	output, err := reader.ReadPersistedOutput(ctx, id)
	if err != nil {
		return "", fmt.Errorf("persisted output %q not found", id)
	}

	runes := []rune(output.Content)
	total := len(runes)
	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	chunk := string(runes[offset:end])
	hasMore := end < total
	nextOffset := end

	payload := map[string]any{
		"id":             output.ID,
		"kind":           output.Kind,
		"tool_call_id":   output.ToolCallID,
		"offset":         offset,
		"limit":          limit,
		"returned_chars": end - offset,
		"total_chars":    total,
		"original_bytes": output.OriginalBytes,
		"next_offset":    nextOffset,
		"has_more":       hasMore,
		"content_sha256": output.ContentSHA256,
		"content":        chunk,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal persisted output: %w", err)
	}
	return string(data), nil
}
