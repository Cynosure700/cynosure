# 子Agent四层压缩机制设计文档

## 1. 背景

主 Agent 在每轮 LLM 请求前会对 `state.ModelHistory` 执行上下文压缩，并把压缩结果写回 `state.ModelHistory`，使内存态、发送态和后续落库态保持一致。当前默认策略链包含：

1. `ToolResultCompressionStrategy`：大工具结果落盘。
2. `MessageWindowCompressionStrategy`：头尾窗口裁剪，即 `snip_compact`。
3. `RecentToolResultRetentionStrategy`：只保留最近完整工具结果，即 `micro_compact`。
4. `ConversationMemoryStrategy`：会话记忆注入。
5. `FullHistorySummarizationStrategy`：全量摘要兜底，即 `compact_history`。

子 Agent 由 `spawn_subagent` 派生，使用 fresh message list，不读取父对话历史。当前 `runSubagentLoop` 直接使用 `state.Messages` 构造请求，没有在每轮请求前执行压缩，因此长工具链或大工具结果会持续进入子 Agent 上下文。

## 2. 目标

为子 Agent 增加与主 Agent 一致的四层压缩机制：

1. `tool_result_budget`：单个工具结果超过工具声明上限时落盘，并在上下文中替换为 `<persisted-output>` marker。
2. `snip_compact`：当轮消息窗口超过阈值时保留用户最新消息及尾部消息。
3. `micro_compact`：完整 inline tool result 超过 20 条时，仅保留最近 5 条。
4. `compact_history`：前三层后仍超出上下文预算时生成摘要并保留最近 5 条逐字消息。

压缩结果必须写回子 Agent 的 `state.ModelHistory`，并用它重建 `state.Messages` 后再请求 LLM。

## 3. 边界

1. 不让子 Agent 启用 `ConversationMemoryStrategy`。子 Agent 看不到父历史，也不应注入父会话记忆；本需求中的“四层压缩”不包含会话记忆层。
2. 不改变父 Agent 默认压缩链、记忆提取、历史落库、TUI 展示、审批、超时、禁止嵌套 subagent 等既有行为。
3. 不改变任何压缩策略内部阈值、marker 格式、占位符文案或 `read_persisted_output` 行为。
4. 不新增用户配置项。

## 4. 设计

### 4.1 子 Agent 专用压缩链

在 `internal/agent/runtime/compression` 中新增构造函数：

```go
func NewSubagentCompressor() *Compressor
```

它按顺序注册四个策略：

1. `ToolResultCompressionStrategy`
2. `MessageWindowCompressionStrategy`
3. `RecentToolResultRetentionStrategy`
4. `FullHistorySummarizationStrategy`

主 Agent 继续使用 `NewDefaultCompressor()`，包含原有五层策略。

### 4.2 Runtime 压缩入口

在 `internal/agent/runtime/context_compression.go` 增加：

```go
func (s *Service) compressSubagentContextBeforeLLM(ctx context.Context, state *LoopState, tools *ToolRegistry) ([]storage.Message, error)
```

该函数与主 Agent 的 `compressContextBeforeLLM` 保持一致的核心语义：

- 从 `state.ModelHistory` 克隆出 `RequestHistory`。
- `DisplayHistory` 使用 `state.History` 的克隆。
- `Conversation`、`User` 使用当前子 Agent state 中的父会话与用户，确保落盘、hash 去重和读取作用域与现有存储兼容。
- `Tools` 和 `ToolMaxResultSizeChars` 使用子 Agent 的 `childTools`，避免使用主 Agent registry 的预算。
- 使用 `NewSubagentCompressor()` 执行四层压缩。

当 Store 不支持 compression.Store 时，保持现有降级语义，直接返回克隆后的 `state.ModelHistory`。

### 4.3 子 Agent 循环接入

在 `runSubagentLoop` 每轮请求 LLM 之前：

1. 调用 `compressSubagentContextBeforeLLM(ctx, state, tools)`。
2. 将返回值写回 `state.ModelHistory`。
3. 使用 `buildOpenAIMessages(state.SystemPrompt, state.ModelHistory)` 重建 `state.Messages`。
4. 保持后续 `maybeAppendTodoWriteReminder`、LLM 请求、工具执行、审批与现有逻辑一致。

子 Agent 的初始 `ModelHistory` 必须包含用户任务消息。否则压缩入口没有可压缩的真实消息线，后续工具调用历史也无法作为统一基线。

## 5. 测试策略

1. 新增单元测试验证子 Agent 压缩使用子工具 registry 的 `MaxResultSizeChars`：把子工具 `bash` 上限设为 10，构造超过上限的 tool result，调用子 Agent 压缩入口后应产生 `<persisted-output>` marker，且 `state.History` 保持原始结果。
2. 新增单元测试验证子 Agent 压缩不注入 `ConversationMemoryStrategy`：在 Store 中预置会话记忆，并让摘要预算不触发；压缩后上下文不应出现 `<conversation-memory>`。
3. 运行 runtime 定向测试与 compression 包测试，确认主 Agent 压缩策略未被破坏。

## 6. 风险与缓解

1. 风险：子 Agent 压缩使用主工具 registry，导致预算与实际可用工具不一致。缓解：压缩入口显式接收 `*ToolRegistry` 参数，并使用子 registry。
2. 风险：子 Agent 注入父会话记忆，破坏 fresh message list 语义。缓解：子 Agent 使用专用四层 compressor，不包含 `ConversationMemoryStrategy`。
3. 风险：只压缩 `state.Messages` 而不写回 `state.ModelHistory`，下一轮又恢复旧上下文。缓解：与主 Agent 一样写回 `state.ModelHistory`，再重建 `state.Messages`。
