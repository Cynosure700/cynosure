package compression

import (
	"context"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	agenttools "nano_cc/internal/tools"
)

const (
	// toolResultPreviewRunes 表示超大 tool_result 在内联保留时，作为预览的前导 rune 数量。
	toolResultPreviewRunes = 2000
	// messageWindowLimit 是触发首尾裁剪的阈值：当最近一个用户回合的消息数量超过它时触发
	//（仅统计当前回合的消息数，不含更早的历史）。
	messageWindowLimit = 50
	messageWindowHead  = 3
	messageWindowTail  = 46
	// recentToolResultRetentionThreshold 表示仅当完整内联的 tool_result 数量超过该值时，
	// 才触发微压缩。
	recentToolResultRetentionThreshold = 30
	// recentToolResultRetention 表示一旦触发微压缩，保留最近 N 个完整内联的 tool_result。
	recentToolResultRetention = 15

	earlierToolResultPlaceholder = "[Earlier result compacted. Re-run if needed]"
	// PersistedOutputMarkerPrefix 用于标记某个 tool_result：其完整内容已在带外持久化，
	// 并在内联处替换为预览。
	PersistedOutputMarkerPrefix = "<persisted-output"
)

// Request 携带仅供请求使用的历史副本，以及各策略所需的依赖。
// 各策略会就地修改 RequestHistory。
type Request struct {
	Conversation   storage.Conversation
	User           storage.User
	RequestHistory []storage.Message
	SystemPrompt   string
	Tools          []openai.Tool
	Store          Store
	Estimator      TokenEstimator
	Summarizer     HistorySummarizer
	// ToolMaxResultSizeChars 用于解析每个工具的内联结果上限。当解析器为 nil 或返回非正数时，
	// 回退到默认的工具上限。
	ToolMaxResultSizeChars func(toolName string) int
	// ConversationMemoryBreakpoint 是最后一条已折叠进会话记忆的消息 ID（作为锚点被包含保留）。
	// ConversationMemoryStrategy 会把从该消息（含）到末尾的尾部作为“未压缩”消息保留。
	// 为空表示未知——此时该策略让位于全量历史摘要兜底。它从会话记忆文件加载，
	// 而非从内存状态读取。
	ConversationMemoryBreakpoint string
	// DisplayHistory 是逐字展示历史（state.History）的克隆。
	// ConversationMemoryStrategy 与 FullHistorySummarizationStrategy 从【这个列表】保留尾部
	//（而非 RequestHistory / 模型线），因此保留的最近消息是原始、未经压缩的版本。
	DisplayHistory []storage.Message
}

func (r *Request) maxResultSizeChars(toolName string) int {
	if r != nil && r.ToolMaxResultSizeChars != nil {
		if limit := r.ToolMaxResultSizeChars(toolName); limit > 0 {
			return limit
		}
	}
	return agenttools.DefaultMaxResultSizeChars
}

// Store 是各策略所需的最小存储接口。
type Store interface {
	CreatePersistedOutput(ctx context.Context, output storage.PersistedOutput) error
	GetPersistedOutputByMessageHash(ctx context.Context, conversationID, userID, messageID, toolCallID, strategy, contentSHA256 string) (storage.PersistedOutput, error)
	ListConversationMemories(ctx context.Context, conversationID string) ([]storage.ConversationMemory, error)
}

// Strategy 对 RequestHistory 施加一层压缩。
type Strategy interface {
	Name() string
	Apply(ctx context.Context, req *Request) error
}

// Compressor 按顺序运行已注册的各个策略。
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

// NewSubagentCompressor returns the four-layer request compression chain used
// by child agents. It intentionally excludes ConversationMemoryStrategy so a
// fresh subagent context never receives parent conversation memories.
func NewSubagentCompressor() *Compressor {
	return &Compressor{strategies: []Strategy{
		&ToolResultCompressionStrategy{},
		&MessageWindowCompressionStrategy{},
		&RecentToolResultRetentionStrategy{},
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
