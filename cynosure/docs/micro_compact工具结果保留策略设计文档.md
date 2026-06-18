# micro_compact 工具结果保留策略设计文档

## 1. 背景

当前 `cynosure` 的请求上下文压缩采用策略链，在每轮请求 LLM 前对 `ModelHistory` 的拷贝执行压缩，不直接修改 TUI 展示历史或已保存的完整模型历史。默认策略顺序为：

1. `ToolResultCompressionStrategy`：单个大 `tool_result` 超出工具声明上限时落盘，并在上下文中替换为 `<persisted-output>` marker 与预览。
2. `MessageWindowCompressionStrategy`：按消息窗口裁剪上下文。
3. `RecentToolResultRetentionStrategy`：把较早的完整 `tool_result` 替换为一行占位符。
4. `ConversationMemoryStrategy`：注入会话记忆。
5. `FullHistorySummarizationStrategy`：上下文仍超预算时做全量摘要兜底。

本需求中的 `micro_compact` 对应当前代码中的 `RecentToolResultRetentionStrategy`。代码内没有字面量 `micro_compact`，本设计沿用现有策略名，并把用户侧语义称为 `micro_compact`。

## 2. 当前实现核查结论

### 2.1 现有行为

`RecentToolResultRetentionStrategy` 当前逻辑位于 `internal/agent/runtime/compression/recent_tool_result_retention.go`：

1. 遍历 `RequestHistory` 中所有 `Role == "tool"` 的消息，按历史顺序收集索引。
2. 如果 tool message 总数 `<= recentToolResultRetention`，直接返回。
3. 当前 `recentToolResultRetention = 3`。
4. 超过 3 条时，仅保留最后 3 条 tool message 的完整内容。
5. 更早的 tool message 会解析出 `result`，若结果尚未被压缩，则替换为一行占位符：`[Earlier result compacted. Re-run if needed]`。
6. 若原始 tool message content 是 JSON wrapper，则保留 `status` 字段，只替换 `result` 字段；非 JSON content 直接替换为纯占位符。

### 2.2 已有可复用能力

以下能力应继续复用，避免引入重复判断：

- `textutil.ParseToolResult(content)`：解析 tool message 的 JSON wrapper，得到 `status`、`result`、`isJSON`。
- `isCompactedResult(result)`：判断结果是否已经是占位符或 `<persisted-output>` marker。
- `rebuildToolResult(status, result, isJSON)`：替换结果时保留原 JSON wrapper 与状态。
- `earlierToolResultPlaceholder`：现有一行占位符文案。

### 2.3 需要调整的问题

当前策略按所有 tool message 数量触发，即使其中很多已经是 `<persisted-output>` 或占位符，也会参与阈值判断。新需求要求：

> 如果本次会话工具调用数超过了 20 条【只统计工具未落盘且没有占位符的标识】，只保留最近 5 条 `tool_result` 的完整内容，更旧的替换为一行占位符。

因此触发判断和保留数量都需要调整：

- 触发阈值从 `> 3` 改为 `> 20`。
- 保留完整结果数量从最近 `3` 条改为最近 `5` 条。
- 计数对象从“所有 tool message”改为“未落盘且没有占位符的完整 inline tool result”。

## 3. 目标与边界

### 3.1 目标

1. `micro_compact` 只统计仍完整内联的 tool result：
   - 角色必须是 `Role == "tool"`；
   - 解析后的 `result` 不能包含 `<persisted-output` marker；
   - 解析后的 `result` 不能等于现有一行占位符。
2. 当上述可统计 tool result 数量 `<= 20` 时，策略不做任何替换。
3. 当上述可统计 tool result 数量 `> 20` 时：
   - 只保留最近 5 条可统计 tool result 的完整内容；
   - 更旧的可统计 tool result 替换为现有一行占位符；
   - 已经落盘或已经是占位符的 tool result 保持原样，不重复处理。
4. 其他压缩策略的顺序、触发条件和行为不变。
5. 不改变展示历史、模型历史持久化、工具执行日志、`read_persisted_output`、大结果落盘 marker 格式和工具审批流程。

### 3.2 非目标

1. 不调整 `ToolResultCompressionStrategy` 的单工具结果落盘策略。
2. 不调整 `MessageWindowCompressionStrategy` 的消息窗口裁剪策略。
3. 不调整 `ConversationMemoryStrategy`、`FullHistorySummarizationStrategy` 或 `ReactiveCompactStrategy`。
4. 不新增用户配置项。本次阈值按代码常量实现。
5. 不修改占位符文案，继续使用现有 `earlierToolResultPlaceholder`。
6. 不为已经占位或已经落盘的结果反向恢复内容。

## 4. 关键定义

### 4.1 可统计完整工具结果

本设计把“工具未落盘且没有占位符的标识”定义为：

```go
msg.Role == "tool" && !isCompactedResult(parsedResult)
```

其中 `parsedResult` 来自：

```go
status, result, isJSON := textutil.ParseToolResult(msg.Content)
```

判定说明：

- `<persisted-output ...>` marker 表示完整结果已落盘，不参与 20 条计数。
- `[Earlier result compacted. Re-run if needed]` 表示已经被 `micro_compact` 替换为占位符，不参与 20 条计数。
- 非 JSON 的 tool message content 按整体内容作为 `result`；只要不是上述 marker 或占位符，就参与计数。
- 空字符串如果来自 tool message 且未被压缩，也视为一个完整 inline tool result，参与计数。

### 4.2 本次会话范围

本设计中的“本次会话”沿用当前 `RecentToolResultRetentionStrategy` 的作用范围：当前请求的 `RequestHistory` 全量消息，而不是仅最后一个 user turn。

原因：

- 现有 `RecentToolResultRetentionStrategy` 已经按整个 `RequestHistory` 收集 tool message。
- `RequestHistory` 是当前会话模型历史的请求副本，已经过前置策略处理。
- 这样可以继续约束跨多轮会话累计的完整 inline tool result 数量，避免旧轮次工具结果持续占用上下文。

如果产品语义希望只统计“最后一轮 user 后的工具调用”，需要另行调整为复用 `latestUserTurnToolIndexes`，但这会改变当前策略的跨轮压缩语义，不作为本次推荐方案。

## 5. 方案对比

### 方案 A：在现有 `RecentToolResultRetentionStrategy` 内改造统计对象与阈值（推荐）

做法：

- 将常量调整为：
  - `recentToolResultRetentionThreshold = 20`
  - `recentToolResultRetention = 5`
- 遍历 `RequestHistory` 时，只收集未落盘且无占位符的完整 inline tool result。
- 仅当可统计数量超过 20 时，把除最近 5 条外的更旧完整结果替换为占位符。

优点：

- 改动最小，完全复用现有策略位置、helper 和测试结构。
- 不影响其他策略，也不改变大结果落盘与读取链路。
- 与“其他策略不变”的约束最一致。

缺点：

- `micro_compact` 名称仍不会出现在代码里，除非额外重命名策略；但重命名会扩大回归面。

### 方案 B：新增一个独立 `MicroCompactStrategy`

做法：

- 新增策略类型，注册在当前 `RecentToolResultRetentionStrategy` 的位置。
- 保留旧策略代码但不再注册，或删除旧策略。

优点：

- 策略名可以与产品语义一致。

缺点：

- 增加重复代码或重命名成本。
- 需要同步更多测试与文档，且没有行为收益。
- 可能误导为新增策略，而需求实际是调整当前策略。

### 方案 C：把计数逻辑合并进 `ToolResultCompressionStrategy`

做法：

- 大结果落盘后，在同一个策略中继续做数量阈值压缩。

优点：

- 落盘 marker 与占位符判断都集中在 tool result compression 入口。

缺点：

- 混合了“大结果可恢复落盘”和“旧结果不可恢复占位”两种不同策略。
- 会改变现有策略边界和测试职责。
- 不符合“其他策略不变”的约束。

### 推荐

采用方案 A。它只调整当前 `micro_compact` 对应策略的触发条件和保留数量，保留策略链顺序、marker 语义、占位符语义和 request-only 压缩边界。

## 6. 详细设计

### 6.1 常量调整

在 `internal/agent/runtime/compression/compression.go` 中调整或新增常量：

```go
const (
    // recentToolResultRetentionThreshold triggers micro compaction only when
    // full inline tool results exceed this count.
    recentToolResultRetentionThreshold = 20

    // recentToolResultRetention keeps the most recent N full inline tool results
    // once micro compaction is triggered.
    recentToolResultRetention = 5
)
```

原有 `recentToolResultRetention = 3` 改为 `5`。新增 `recentToolResultRetentionThreshold`，避免用“保留数量”同时表达“触发阈值”。

### 6.2 收集可统计 tool result

在 `RecentToolResultRetentionStrategy.Apply` 中把原来的 `toolIndexes` 改为 `inlineToolResults`，只收集未压缩结果：

```go
type inlineToolResult struct {
    index  int
    status string
    result string
    isJSON bool
}
```

遍历逻辑：

```go
var candidates []inlineToolResult
for i := range history {
    if history[i].Role != "tool" {
        continue
    }
    status, result, isJSON := textutil.ParseToolResult(history[i].Content)
    if isCompactedResult(result) {
        continue
    }
    candidates = append(candidates, inlineToolResult{
        index:  i,
        status: status,
        result: result,
        isJSON: isJSON,
    })
}
```

这样计数天然排除：

- 已由 `ToolResultCompressionStrategy` 替换为 `<persisted-output>` 的结果；
- 已由历史 `micro_compact` 替换为一行占位符的结果；
- 任何后续复跑压缩时已经处理过的结果。

### 6.3 触发判断

触发判断改为：

```go
if len(candidates) <= recentToolResultRetentionThreshold {
    return nil
}
```

这意味着：

- 20 条完整 inline tool result：不触发，全部保留。
- 21 条完整 inline tool result：触发，只保留最近 5 条，压缩前 16 条。
- 如果 `RequestHistory` 中总共有 30 条 tool message，但其中 12 条已经落盘或占位，则可统计数量为 18，不触发。

### 6.4 保留最近 5 条完整内容

触发后计算：

```go
cutoff := len(candidates) - recentToolResultRetention
for _, candidate := range candidates[:cutoff] {
    history[candidate.index].Content = rebuildToolResult(
        candidate.status,
        earlierToolResultPlaceholder,
        candidate.isJSON,
    )
}
```

效果：

- `candidates[cutoff:]` 是最近 5 条可统计完整结果，保持原文。
- 更旧的可统计完整结果替换为一行占位符。
- 不在 `candidates` 里的 tool message 保持原样，包括 `<persisted-output>` marker 和已有占位符。

### 6.5 与 JSON wrapper 的兼容

继续使用 `rebuildToolResult`：

- JSON wrapper 输入：输出仍是 `{"status":"...","result":"[Earlier result compacted. Re-run if needed]"}`。
- 非 JSON 输入：输出为 `[Earlier result compacted. Re-run if needed]`。

这保持现有模型请求格式兼容，不改变 tool result status 的保留方式。

### 6.6 与其他压缩策略的关系

策略顺序保持不变：

```go
ToolResultCompressionStrategy
MessageWindowCompressionStrategy
RecentToolResultRetentionStrategy
ConversationMemoryStrategy
FullHistorySummarizationStrategy
```

影响说明：

- `ToolResultCompressionStrategy` 仍先运行；被它落盘的结果进入 `<persisted-output>` marker，不参与 `micro_compact` 的 20 条计数。
- `MessageWindowCompressionStrategy` 仍可能先裁剪消息；`micro_compact` 只处理裁剪后仍在 `RequestHistory` 中的 tool message。
- `ConversationMemoryStrategy` 和 `FullHistorySummarizationStrategy` 继续基于前面策略处理后的 request history。
- `ReactiveCompactStrategy` 不使用 `RecentToolResultRetentionStrategy`，不受本次改动影响。

## 7. 测试设计

### 7.1 更新现有单测

修改 `internal/agent/runtime/compression/compression_test.go` 中 `RecentToolResultRetentionStrategy` 相关测试：

1. 原 `TestRecentToolRetention_NoopAtOrBelowThree` 改为覆盖 20 条阈值：
   - 构造 20 条完整 inline tool result；
   - 执行策略后全部保持原文；
   - 验证没有任何结果被替换为占位符。
2. 原 `TestRecentToolRetention_CompactsOlderKeepsRecentThree` 改为覆盖 21 条触发：
   - 构造 21 条完整 inline tool result；
   - 执行策略后前 16 条为占位符；
   - 最近 5 条保持完整；
   - JSON wrapper 的 `status` 保持不变。
3. 原 `TestRecentToolRetention_NonJSONReplacedWithPlainPlaceholder` 保留并调整数量：
   - 构造超过 20 条，其中最旧一条为非 JSON；
   - 执行策略后非 JSON 旧结果被替换为纯占位符。

### 7.2 新增单测

新增测试覆盖本需求的关键边界：

1. `TestRecentToolRetention_CountsOnlyFullInlineToolResults`
   - 构造 25 条 tool message，其中 10 条是 `<persisted-output>` 或已有占位符，15 条是完整 inline result；
   - 因可统计数量为 15，不触发压缩；
   - 验证完整结果保持原文，marker 和占位符保持原样。
2. `TestRecentToolRetention_SkipsPersistedAndPlaceholderWhenKeepingRecentFive`
   - 构造超过 20 条完整 inline result，并夹杂若干 marker 与占位符；
   - 验证最近 5 条“完整 inline result”保留；
   - 验证 marker 与已有占位符不被改写，也不占用最近 5 条名额。
3. `TestRecentToolRetention_NoopAtExactlyTwentyFullInlineResults`
   - 明确覆盖边界：正好 20 条不触发。
4. `TestRecentToolRetention_CompactsAtTwentyOneFullInlineResults`
   - 明确覆盖边界：21 条触发，保留 5 条。

### 7.3 验证命令

实现阶段建议按以下顺序验证：

```bash
go test ./internal/agent/runtime/compression -run 'TestRecentToolRetention' -count=1
go test ./internal/agent/runtime/compression -count=1
go test ./internal/agent/runtime/... -count=1
go test ./... -count=1
```

## 8. 验收标准

1. 当当前请求历史中可统计完整 inline tool result 数量不超过 20 条时，`RecentToolResultRetentionStrategy` 不替换任何新内容。
2. 当可统计完整 inline tool result 数量超过 20 条时，只保留最近 5 条完整内容。
3. 更旧的完整 inline tool result 被替换为现有一行占位符。
4. `<persisted-output>` marker 不参与计数，不被替换。
5. 已有一行占位符不参与计数，不被重复替换。
6. JSON wrapper tool result 的 `status` 字段保持不变。
7. 非 JSON tool result 仍按现有方式替换为纯占位符。
8. 压缩策略顺序不变，其他策略测试不需要因行为变化而调整。
9. TUI 展示历史、模型历史持久化和工具输出日志不受影响。

## 9. 风险与规避

1. **“本次会话”范围理解偏差**
   - 风险：如果期望只统计最后一轮 user 后的工具结果，按全 `RequestHistory` 统计会更积极地压缩跨轮旧结果。
   - 规避：本设计明确沿用当前策略的全 `RequestHistory` 语义；若审阅后要求按最新 user turn，可在实现前调整设计。
2. **已落盘 marker 被误替换**
   - 风险：如果直接按 `Role == "tool"` 计数，`<persisted-output>` 可能占用名额或被替换。
   - 规避：统一通过 `textutil.ParseToolResult` 和 `isCompactedResult` 判断，只收集完整 inline result。
3. **已有占位符重复处理**
   - 风险：重复压缩可能改变已有占位符或导致边界判断错误。
   - 规避：已有占位符由 `isCompactedResult` 排除。
4. **测试只覆盖数量不覆盖 marker**
   - 风险：实现满足 20/5 数字但不满足“只统计未落盘且无占位符”。
   - 规避：新增 marker 和占位符混合场景测试。

## 10. 建议实施步骤

待本设计审阅通过后，再进入代码实现。建议实施顺序：

1. 先更新 `RecentToolResultRetentionStrategy` 相关测试，覆盖 20 条阈值、21 条触发、最近 5 条保留、marker/占位符排除计数。
2. 修改 `compression.go` 中相关常量，新增触发阈值常量。
3. 修改 `recent_tool_result_retention.go`，把统计对象改为完整 inline tool result candidates。
4. 运行 compression 包测试。
5. 运行 runtime 压缩相关测试和全量测试。

## 11. 设计自检

- 未包含待补充占位内容。
- 范围只覆盖 `micro_compact` 对应的 `RecentToolResultRetentionStrategy`。
- 明确排除了 `<persisted-output>` 和已有一行占位符的计数。
- 明确保留其他压缩策略、落盘读取、历史持久化和 TUI 展示行为。
- 明确了 20 条不触发、21 条触发、触发后最近 5 条完整保留的边界。
