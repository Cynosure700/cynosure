package compression

import (
	"context"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	agenttools "nano_cc/internal/tools"
)

const (
	// toolResultPreviewRunes is how many leading runes of an oversized
	// tool_result remain inline as a preview.
	toolResultPreviewRunes = 2000
	// messageWindowLimit triggers head/tail trimming when the latest user turn
	// exceeds it (message count of the current turn, excluding earlier history).
	messageWindowLimit = 50
	messageWindowHead  = 3
	messageWindowTail  = 46
	// recentToolResultRetentionThreshold triggers micro compaction only when
	// full inline tool results exceed this count.
	recentToolResultRetentionThreshold = 20
	// recentToolResultRetention keeps the most recent N full inline tool
	// results once micro compaction is triggered.
	recentToolResultRetention = 5

	earlierToolResultPlaceholder = "[Earlier result compacted. Re-run if needed]"
	// PersistedOutputMarkerPrefix marks a tool_result whose full content was
	// persisted out-of-band and replaced inline with a preview.
	PersistedOutputMarkerPrefix = "<persisted-output"
)

// Request carries the request-only history copy and dependencies that
// strategies use. Strategies mutate RequestHistory in place.
type Request struct {
	Conversation   storage.Conversation
	User           storage.User
	RequestHistory []storage.Message
	SystemPrompt   string
	Tools          []openai.Tool
	Store          Store
	Estimator      TokenEstimator
	Summarizer     HistorySummarizer
	// ToolMaxResultSizeChars resolves the per-tool inline result limit. A nil
	// resolver or non-positive result falls back to the default tool limit.
	ToolMaxResultSizeChars func(toolName string) int
}

func (r *Request) maxResultSizeChars(toolName string) int {
	if r != nil && r.ToolMaxResultSizeChars != nil {
		if limit := r.ToolMaxResultSizeChars(toolName); limit > 0 {
			return limit
		}
	}
	return agenttools.DefaultMaxResultSizeChars
}

// Store is the minimal storage surface required by strategies.
type Store interface {
	CreatePersistedOutput(ctx context.Context, output storage.PersistedOutput) error
	GetPersistedOutputByMessageHash(ctx context.Context, conversationID, userID, messageID, toolCallID, strategy, contentSHA256 string) (storage.PersistedOutput, error)
	CreateContextSummary(ctx context.Context, summary storage.ContextSummary) error
	GetContextSummaryByHistoryHash(ctx context.Context, conversationID, userID, sourceHistorySHA256 string) (storage.ContextSummary, error)
	ListConversationMemories(ctx context.Context, conversationID string) ([]storage.ConversationMemory, error)
}

// Strategy applies one compression layer to RequestHistory.
type Strategy interface {
	Name() string
	Apply(ctx context.Context, req *Request) error
}

// Compressor runs the registered strategies in order.
type Compressor struct {
	strategies []Strategy
}

func NewDefaultCompressor() *Compressor {
	return &Compressor{strategies: []Strategy{
		&ToolResultCompressionStrategy{},
		&MessageWindowCompressionStrategy{},
		&RecentToolResultRetentionStrategy{},
		&ConversationMemoryStrategy{},
		&FullHistorySummarizationStrategy{},
	}}
}

func (c *Compressor) Compress(ctx context.Context, req *Request) error {
	for _, strategy := range c.strategies {
		if err := strategy.Apply(ctx, req); err != nil {
			return err
		}
	}
	return nil
}
