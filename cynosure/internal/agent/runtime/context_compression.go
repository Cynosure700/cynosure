package runtime

import (
	"context"
	"fmt"

	"nano_cc/internal/agent/runtime/compression"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/logger"
	agenttools "nano_cc/internal/tools"
)

// loadModelHistory 返回从上一回合压缩后的请求历史中持久化下来、可复用的
// “模型历史”。当不存在对应记录、解码失败或存储出错时，回退为完整展示历史
// 的克隆（即当前行为）。
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

// compressContextBeforeLLM 深拷贝模型历史并运行压缩流水线，返回本轮压缩后的
// 真实消息历史。调用方会把结果赋回 state.ModelHistory，使这条唯一的真实消息
// 历史始终反映最新的压缩输出（内存态 == 发送态 == 落库态）。逐字的展示历史
// state.History 不会被改动（它只会被克隆进 req.DisplayHistory，以便摘要/会话
// 记忆策略从原始消息中保留各自的尾部）。
func (s *Service) compressContextBeforeLLM(ctx context.Context, state *LoopState) ([]storage.Message, error) {
	requestHistory := cloneMessages(state.ModelHistory)
	store, ok := s.Store.(compression.Store)
	if !ok {
		// 存储不支持压缩产物；静默跳过。
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

// reactiveCompact 在 LLM 因 HTTP 413（上下文溢出）拒绝请求时，带外执行激进的
// ReactiveCompactStrategy。成功时会同时更新 state.Messages（本轮生效）与
// state.ModelHistory（后续各轮复用的新基线），但绝不更新 state.History（逐字
// 展示历史）。失败时保持 state 不变。
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

// persistedOutputReader 将存储层适配为面向工具的读取器，并强制按会话/用户进行
// 作用域隔离。
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
