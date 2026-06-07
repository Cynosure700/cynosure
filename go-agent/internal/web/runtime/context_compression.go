package runtime

import (
	"context"

	agenttools "nano_cc/internal/tools"
	"nano_cc/internal/web/runtime/compression"
	"nano_cc/internal/web/storage"
)

// compressContextBeforeLLM deep-copies the display history and runs the
// compression pipeline, returning the request-only history for this round.
func (s *Service) compressContextBeforeLLM(ctx context.Context, state *LoopState) ([]storage.Message, error) {
	requestHistory := cloneMessages(state.History)
	store, ok := s.Store.(compression.Store)
	if !ok {
		// Store does not support compression artifacts; skip silently.
		return requestHistory, nil
	}
	compressor := s.ContextCompressor
	if compressor == nil {
		compressor = compression.NewDefaultCompressor()
	}
	req := &compression.Request{
		Conversation:   state.Conversation,
		User:           state.User,
		RequestHistory: requestHistory,
		SystemPrompt:   state.SystemPrompt,
		Tools:          s.Tools.Definitions(),
		Store:          store,
		Estimator:      compression.DefaultTokenEstimator{},
		Summarizer:     s.summarizeHistoryForContext,
	}
	if err := compressor.Compress(ctx, req); err != nil {
		return nil, err
	}
	return req.RequestHistory, nil
}

func cloneMessages(messages []storage.Message) []storage.Message {
	if len(messages) == 0 {
		return nil
	}
	cloned := make([]storage.Message, len(messages))
	for i, msg := range messages {
		cloned[i] = msg
		if len(msg.ToolCalls) > 0 {
			calls := make([]storage.MessageToolCall, len(msg.ToolCalls))
			copy(calls, msg.ToolCalls)
			cloned[i].ToolCalls = calls
		}
	}
	return cloned
}

// persistedOutputReader adapts the storage layer to the tool-facing reader,
// enforcing conversation/user scoping.
type persistedOutputReader struct {
	store          compressionReaderStore
	userID         string
	conversationID string
}

type compressionReaderStore interface {
	GetPersistedOutputForConversation(ctx context.Context, id, userID, conversationID string) (storage.PersistedOutput, error)
}

func (r persistedOutputReader) ReadPersistedOutput(ctx context.Context, id string) (agenttools.PersistedOutput, error) {
	output, err := r.store.GetPersistedOutputForConversation(ctx, id, r.userID, r.conversationID)
	if err != nil {
		return agenttools.PersistedOutput{}, err
	}
	return agenttools.PersistedOutput{
		ID:            output.ID,
		Kind:          output.Kind,
		ToolCallID:    output.ToolCallID,
		OriginalBytes: output.OriginalBytes,
		ContentSHA256: output.ContentSHA256,
		Content:       output.Content,
	}, nil
}

func (s *Service) newPersistedOutputReader(conversationID, userID string) agenttools.PersistedOutputReader {
	store, ok := s.Store.(compressionReaderStore)
	if !ok {
		return nil
	}
	return persistedOutputReader{store: store, userID: userID, conversationID: conversationID}
}
