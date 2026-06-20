package runtime

import (
	"context"
	"fmt"

	"nano_cc/internal/agent/runtime/compression"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/logger"
	agenttools "nano_cc/internal/tools"
)

// loadModelHistory returns the reusable "model history" persisted from the
// previous turn's compressed request history. When no row exists, decoding
// fails, or the store errors, it falls back to a clone of the full display
// history (current behavior).
func (s *Service) loadModelHistory(ctx context.Context, conversationID string, displayHistory []storage.Message) []storage.Message {
	modelHistory, ok, err := s.Store.GetConversationModelHistory(ctx, conversationID)
	if err != nil {
		logger.Warn(fmt.Sprintf("model history: load failed conversation=%s: %v", conversationID, err))
	}
	if !ok || len(modelHistory) == 0 {
		return cloneMessages(displayHistory)
	}
	return modelHistory
}

// compressContextBeforeLLM deep-copies the model history and runs the
// compression pipeline, returning the request-only history for this round.
func (s *Service) compressContextBeforeLLM(ctx context.Context, state *LoopState) ([]storage.Message, error) {
	requestHistory := cloneMessages(state.ModelHistory)
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

		DisplayHistory:               cloneMessages(state.History),
		ConversationMemoryBreakpoint: s.loadConversationMemoryBreakpoint(ctx, state.Conversation.ID),
	}
	if s.Tools != nil {
		req.ToolMaxResultSizeChars = s.Tools.MaxResultSizeChars
	}
	if err := compressor.Compress(ctx, req); err != nil {
		return nil, err
	}
	return req.RequestHistory, nil
}

// reactiveCompact runs the aggressive ReactiveCompactStrategy out-of-band when
// the LLM rejects a request with HTTP 413 (context overflow). On success it
// updates both state.Messages (effective this round) and state.ModelHistory
// (the new baseline reused by later rounds), but never state.History (the
// verbatim display history). On failure state is left untouched.
func (s *Service) reactiveCompact(ctx context.Context, state *LoopState) error {
	store, ok := s.Store.(compression.Store)
	if !ok {
		return fmt.Errorf("reactive compact: store does not support compression artifacts")
	}
	requestHistory := cloneMessages(state.ModelHistory)
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
	if s.Tools != nil {
		req.ToolMaxResultSizeChars = s.Tools.MaxResultSizeChars
	}
	if err := (&compression.ReactiveCompactStrategy{}).Apply(ctx, req); err != nil {
		return err
	}
	state.ModelHistory = req.RequestHistory
	state.Messages = buildOpenAIMessages(state.SystemPrompt, req.RequestHistory)
	return nil
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
