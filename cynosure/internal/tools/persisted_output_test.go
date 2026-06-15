package tools

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

type fakePersistedReader struct {
	output PersistedOutput
	err    error
}

func (r fakePersistedReader) ReadPersistedOutput(ctx context.Context, id string) (PersistedOutput, error) {
	if r.err != nil {
		return PersistedOutput{}, r.err
	}
	if id != r.output.ID {
		return PersistedOutput{}, errors.New("not found")
	}
	return r.output, nil
}

func TestReadPersistedOutput_ReturnsChunkWithPagination(t *testing.T) {
	content := strings.Repeat("a", 50000)
	reader := fakePersistedReader{output: PersistedOutput{ID: "po_1", Kind: "tool_result", ToolCallID: "tc1", OriginalBytes: 50000, Content: content}}
	ctx := WithPersistedOutputReader(context.Background(), reader)

	out, err := handleReadPersistedOutput(ctx, map[string]any{"id": "po_1", "offset": float64(0), "limit": float64(20000)})
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if payload["returned_chars"].(float64) != 20000 {
		t.Fatalf("expected 20000 returned chars, got %v", payload["returned_chars"])
	}
	if payload["has_more"].(bool) != true {
		t.Fatalf("expected has_more true")
	}
	if payload["next_offset"].(float64) != 20000 {
		t.Fatalf("expected next_offset 20000, got %v", payload["next_offset"])
	}
}

func TestReadPersistedOutput_CapsLimit(t *testing.T) {
	content := strings.Repeat("b", 300000)
	reader := fakePersistedReader{output: PersistedOutput{ID: "po_1", Content: content}}
	ctx := WithPersistedOutputReader(context.Background(), reader)

	out, err := handleReadPersistedOutput(ctx, map[string]any{"id": "po_1", "limit": float64(999999)})
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var payload map[string]any
	_ = json.Unmarshal([]byte(out), &payload)
	if payload["returned_chars"].(float64) != float64(maxPersistedOutputReadLimit) {
		t.Fatalf("expected limit capped to %d, got %v", maxPersistedOutputReadLimit, payload["returned_chars"])
	}
}

func TestReadPersistedOutput_DeniesWhenReaderRejects(t *testing.T) {
	reader := fakePersistedReader{err: errors.New("forbidden")}
	ctx := WithPersistedOutputReader(context.Background(), reader)
	if _, err := handleReadPersistedOutput(ctx, map[string]any{"id": "po_other"}); err == nil {
		t.Fatalf("expected error when reader rejects access")
	}
}

func TestReadPersistedOutput_RequiresReaderInContext(t *testing.T) {
	if _, err := handleReadPersistedOutput(context.Background(), map[string]any{"id": "po_1"}); err == nil {
		t.Fatalf("expected error when reader is absent")
	}
}
