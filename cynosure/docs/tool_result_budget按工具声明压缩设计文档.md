# tool_result_budget 按工具声明压缩设计文档

## 1. 背景

当前 `cynosure` 已有 tool result 压缩与本地落盘能力：

- `ToolResultCompressionStrategy` 在请求 LLM 前处理最近一轮 user turn 之后的 tool messages，当前触发条件是最近一轮未压缩 tool_result 总字节数超过 `200 * 1024`。
- 超阈值后按结果大小从大到小持久化，内联内容替换为 `<persisted-output>` 标记和前 2000 字符预览。
- 完整结果已经通过 `local.Store.CreatePersistedOutput` 写入现有目录 `~/.cynosure/task_outputs/{workspace}/{session_id}/tool-results/`，并可由 `read_persisted_output` 按 id 分段读取。
- 工具执行日志已经追加到 `~/.cynosure/task_outputs/{workspace}/{session_id}/tools.md`。
- 其他压缩策略包括 `MessageWindowCompressionStrategy`、`RecentToolResultRetentionStrategy`、`ConversationMemoryStrategy` 和 `FullHistorySummarizationStrategy`。

本次需求是调整 `tool_result_budget` 压缩策略：由全局总量预算改为每个工具声明 `maxResultSizeChars`，默认 `50,000` 字符；单个工具结果超过该声明值时，将完整结果落盘到项目现有落盘目录，并在上下文中保留可读 marker 与预览。

## 2. 目标与边界

### 2.1 目标

1. 每个工具拥有一个结果内联上限声明：`maxResultSizeChars`。
2. 默认上限为 `50,000` 字符；未显式声明的工具，包括 MCP 工具，使用默认值。
3. 当单个 tool result 的字符数超过对应工具的 `maxResultSizeChars` 时：
   - 完整结果写入现有 `~/.cynosure/task_outputs/{workspace}/{session_id}/tool-results/`；
   - 进入 LLM 请求上下文的 tool message 改为 `<persisted-output>` marker；
   - marker 继续包含前 2000 字符预览；
   - 模型仍通过 `read_persisted_output` 读取完整结果。
4. 保持其他上下文压缩策略、marker 格式、读取工具、工具日志格式和 TUI 展示语义不变。
5. 避免工具内部在返回前截断可恢复结果，确保 runtime 压缩层能拿到完整结果并负责落盘。

### 2.2 非目标

1. 不改变工具调用权限、审批、超时、参数校验和调度逻辑。
2. 不改变 `read_persisted_output` 的参数和返回格式。
3. 不改变 `read_persisted_output` 读取 persisted output 的行为。
4. 不改变 `RecentToolResultRetentionStrategy` 的“只保留最近 N 个完整 tool results”行为。
5. 不把 `tools.md` 自动注入模型上下文。
6. 不为用户配置新增可编辑项；本次只做代码层工具声明。

### 2.3 明确约定

1. `maxResultSizeChars` 按字符数判断，Go 实现使用 rune 数量；metadata 中既有 `original_bytes` 字段继续记录 UTF-8 字节数。
2. `maxResultSizeChars <= 0` 视为未声明，使用默认 `50,000`。
3. 如果无法从 tool message 反查工具名，使用默认 `50,000`。
4. 内联预览长度继续使用现有 `toolResultPreviewRunes = 2000`，不随 `maxResultSizeChars` 改变。
5. 工具输出落盘仍发生在 request-only 压缩阶段：display history 和 model history 保存完整 tool result，发给 LLM 的请求副本被压缩。这沿用当前设计，避免改变历史恢复与审计语义。

## 3. 当前实现核查结论

### 3.1 需要改造的现有逻辑

`internal/agent/runtime/compression/compression.go` 当前定义：

```go
toolResultByteThreshold = 200 * 1024
```

`internal/agent/runtime/compression/tool_result_compression.go` 当前逻辑是：

1. 收集最近一个 user turn 后的 tool messages。
2. 解析每个 tool message 中的 result。
3. 计算未压缩 tool_result 总字节数。
4. 总字节数超过 `toolResultByteThreshold` 后，按单个 result 字节数倒序持久化，直到剩余总字节数回到阈值内。

这与“每个工具声明 `maxResultSizeChars`，单个结果超出即落盘”的新语义不同，需要替换触发判断和候选选择方式。

### 3.2 需要保留的现有能力

以下能力已经存在，设计上继续复用：

- `persistAndBuildMarker` 的 hash 去重、`storage.PersistedOutput` 写入和 marker 生成。
- `local/persisted_output_files.go` 下的 `.txt` / `.json` 文件落盘。
- `GetPersistedOutputForConversation` 和 `GetPersistedOutputByMessageHash` 的本地文件 fallback。
- `read_persisted_output` 工具。
- `tools.md` 工具执行日志。
- `MessageWindowCompressionStrategy`、`RecentToolResultRetentionStrategy`、`ConversationMemoryStrategy`、`FullHistorySummarizationStrategy` 的注册顺序和行为。

### 3.3 需要处理的兼容问题

部分工具当前在工具内部截断结果：

- `bash` 通过 `maxOutputLen = 50000` 截断命令输出。
- `grep content` 复用 `maxOutputLen` 截断匹配内容。
- `web_fetch` 在无 LLM processor 时截断清洗后的网页文本。

如果保留这些截断，压缩层只能拿到已被截断的结果，无法把完整结果落盘。因此本次实现应取消这些“返回前静默截断”，改由统一的 per-tool result budget 在压缩层处理。资源保护类限制仍保留，例如 `webFetchMaxBodySize`、`read_file` 的读取上限、搜索 `head_limit`、命令超时等，因为这些限制属于工具自身的工作量边界，不是 tool_result_budget 压缩策略。

## 4. 方案对比

### 方案 A：在工具定义旁增加本地 metadata，并由压缩策略按 tool call 反查上限（推荐）

为内置工具增加轻量声明结构：

```go
type ToolSpec struct {
    Definition             openai.Tool
    MaxResultSizeChars     int
}
```

`AllToolDefs` 继续对外暴露 `[]openai.Tool`，runtime 额外通过工具名查询 `MaxResultSizeChars`。压缩策略根据 tool message 的 `ToolCallID` 反查前序 assistant tool call 的工具名，再读取对应上限。

优点：

- 不污染 OpenAI tool schema，避免把 runtime 内部预算暴露给模型。
- 改动集中在 tools registry 和 compression 层。
- MCP 工具天然可走默认值。
- 保留现有 `openai.Tool` 使用方式，降低回归风险。

缺点：

- 需要在 compression request 中携带或可查询工具结果上限。
- 需要为 tool_call_id 到 tool name 建立反查逻辑。

### 方案 B：把 `maxResultSizeChars` 注入 OpenAI tool schema

在每个 tool definition 的 JSON schema 或 description 中加入 `maxResultSizeChars`。

优点：

- “每个工具声明”在模型可见的定义中也能看到。

缺点：

- JSON schema 中加入非参数字段容易被误解为工具入参。
- description 加预算会增加 prompt 噪音，且模型不需要知道该实现细节。
- MCP 工具 schema 不一定可控。

### 方案 C：工具执行后立即落盘并把 history 中的结果替换为 marker

在 `executeToolCall` 或 hook 中根据工具声明直接替换 `Outcome.Result`。

优点：

- 单次工具执行后立即完成落盘，后续历史里已经是 marker。

缺点：

- 改变 display history 和 model history 现有语义，用户审计日志可能不再保存完整结果。
- 与当前 request-only 压缩模型不一致，影响历史恢复和测试面更大。
- 需要额外保证 `tools.md` 保存完整结果。

### 推荐

采用方案 A。它最贴近当前架构：工具定义仍负责“声明能力”，compression 层仍负责“进入 LLM 上下文前压缩”，local Store 仍负责“落盘和读取”。这能把变更范围限制在 tool metadata、runtime/compression 请求参数和少量工具内部截断清理上。

## 5. 详细设计

### 5.1 工具声明结构

在 `internal/tools` 中新增工具声明 metadata，保留现有 `openai.Tool` 对外形态。

建议结构：

```go
const DefaultMaxResultSizeChars = 50000

type ToolSpec struct {
    Definition         openai.Tool
    MaxResultSizeChars int
}
```

构造函数：

```go
func toolSpec(name, desc string, params any, opts ...ToolSpecOption) ToolSpec
```

默认行为：

- `toolSpec(...)` 默认设置 `MaxResultSizeChars: DefaultMaxResultSizeChars`。
- 如未来某个工具需要特殊值，可用 option 显式声明，例如 `withMaxResultSizeChars(100000)`。
- `AllToolDefs` 从 `AllToolSpecs` 派生，保持现有 runtime 和测试对 `[]openai.Tool` 的依赖。

查询函数：

```go
func MaxResultSizeCharsForTool(name string) int
```

规则：

- 找到内置声明且值大于 0，返回声明值。
- 找不到或值无效，返回 `DefaultMaxResultSizeChars`。

MCP 工具不在 `AllToolSpecs` 中声明，统一走默认值。

### 5.2 runtime registry 携带结果上限

`ToolRegistry` 内部维护工具名到结果上限的映射：

```go
type ToolRegistry struct {
    definitions        []openai.Tool
    maxResultSizeChars map[string]int
    baseEnv            agenttools.RuntimeEnv
}
```

`NewToolRegistry` 和 `NewChildToolRegistry` 构造时从 `internal/tools` 的声明中填充映射。`ToolRegistry` 提供只读方法：

```go
func (r *ToolRegistry) MaxResultSizeChars(name string) int
```

`Service.toolDefinitionsForUser` 合并 MCP 工具时，不需要为 MCP 单独注册映射；压缩层查询不到时使用默认值。

### 5.3 compression request 增加工具结果上限解析能力

在 `compression.Request` 中增加一个窄接口或函数字段：

```go
ToolResultLimit ToolResultLimitResolver

type ToolResultLimitResolver interface {
    MaxResultSizeChars(toolName string) int
}
```

也可使用函数：

```go
ToolMaxResultSizeChars func(toolName string) int
```

推荐使用函数字段，避免 compression 包依赖 runtime 的具体类型。构造压缩请求时传入 `s.Tools.MaxResultSizeChars`；如果为空，compression 层使用 `DefaultMaxResultSizeChars` 的本地常量或请求字段默认值。

### 5.4 tool_call_id 到 tool name 的反查

`ToolResultCompressionStrategy` 处理 tool message 时，通过 `ToolCallID` 找到对应工具名：

1. 遍历当前 `RequestHistory`。
2. 对每个 assistant message 的 `ToolCalls` 建立 `map[toolCallID]toolName`。
3. 对每个 tool message，用 `ToolCallID` 查询工具名。
4. 查询不到时使用空工具名，最终走默认 `50,000`。

这符合现有 OpenAI tool call 结构，也能覆盖一轮中多个工具调用的场景。

### 5.5 压缩触发逻辑

`ToolResultCompressionStrategy.Apply` 改为 per-tool 判断：

1. 继续只处理最近一个 user turn 之后的 tool messages。
2. 继续跳过已经 compacted 的 result。
3. 对每个 tool message：
   - 解析 status/result/isJSON；
   - 根据 `ToolCallID` 找到 tool name；
   - 获取 `maxResultSizeChars`；
   - 用 rune 数计算 result 字符数；
   - 如果字符数大于上限，则调用现有 `persistAndBuildMarker`；
   - 用现有 `rebuildToolResult` 保留 JSON wrapper 和 status。
4. 不再计算最近一轮 tool_result 总字节数。
5. 不再按大小排序，也不再“持久化到总量低于 200KB 为止”。

示例语义：

```text
bash.maxResultSizeChars = 50000

tool result A: 40000 chars -> 保持内联
tool result B: 50000 chars -> 保持内联
tool result C: 50001 chars -> 落盘并替换为 marker
```

多个工具结果相互独立，不会因为总和超过某个全局阈值而压缩较小结果。

### 5.6 落盘目录与 marker

落盘目录使用当前 workspace 与 session 隔离：

```text
~/.cynosure/task_outputs/{workspace}/{session_id}/tool-results/
├── {persisted_output_id}.txt
└── {persisted_output_id}.json
```

marker 格式保持不变：

```text
<persisted-output id="po_xxx" kind="tool_result" original_bytes="123456" preview_chars="2000" retrieval_tool="read_persisted_output">
完整输出已持久化；如需更多内容，请调用 read_persisted_output(id="po_xxx", offset=0, limit=20000) 分段读取。

前 2000 字符预览
</persisted-output>
```

`original_bytes` 继续表示完整 UTF-8 内容字节数；`preview_chars` 继续表示预览字符数。

### 5.7 工具内部截断清理

为保证“超过声明值后完整落盘”，应清理会丢弃完整结果的内部截断：

- `bash`：移除 `RunBashInDir` 中 `if len(result) > maxOutputLen { result = result[:maxOutputLen] }`。
- `grep content`：移除 `grepContent` 中对最终 `out` 的 `maxOutputLen` 截断。
- `web_fetch`：移除 `fetchAndCleanText` 对清洗文本的 `maxOutputLen` 截断。

保留以下非 tool_result_budget 限制：

- `terminalToolTimeout` 和 `normalToolTimeout`。
- `webFetchMaxBodySize`。
- `grep` / `glob` 的 `head_limit`。
- `read_file` 的 `limit` 和 `maxReadLen`，因为这是读取工具自身的输入/资源边界，不是输出压缩策略。

如果后续希望 `read_file` 也可读取完整大文件并通过 tool_result_budget 落盘，需要另起设计；本次不扩大范围。

### 5.8 与其他压缩策略的关系

策略注册顺序保持：

```go
ToolResultCompressionStrategy
MessageWindowCompressionStrategy
RecentToolResultRetentionStrategy
ConversationMemoryStrategy
FullHistorySummarizationStrategy
```

影响说明：

- `ToolResultCompressionStrategy` 仍先运行，优先把超出单工具声明值的结果替换为 marker。
- `MessageWindowCompressionStrategy` 仍按消息数量裁剪。
- `RecentToolResultRetentionStrategy` 仍只保留最近 N 个完整 tool results；已经是 `<persisted-output>` 的结果会被 `isCompactedResult` 跳过。
- memory 和 full history summarization 继续基于前面策略处理后的 request history。

因此“其他策略不变”在实现上体现为不修改策略顺序、不修改各策略常量、不修改其他策略的触发条件。

## 6. 数据流

### 6.1 工具调用后

1. 工具执行返回完整 `ExecResult.Output`。
2. runtime 构造 `toolExecutionOutcome.Result`。
3. hook 将完整 tool result 写入 `state.History` 和 `state.ModelHistory`。
4. hook 将完整 tool result 追加到 `~/.cynosure/task_outputs/{workspace}/{session_id}/tools.md`。

### 6.2 下一次请求 LLM 前

1. runtime 复制 model history 得到 request history。
2. runtime 构造 `compression.Request`，传入 tool result limit resolver。
3. `ToolResultCompressionStrategy` 查找每个 tool result 对应工具的 `maxResultSizeChars`。
4. 超出上限的结果通过现有 Store 写入 `~/.cynosure/task_outputs/{workspace}/{session_id}/tool-results/`。
5. request history 中该 tool message 替换为 `<persisted-output>` marker。
6. LLM 如需完整内容，调用 `read_persisted_output`。

## 7. 错误处理

1. 落盘失败：沿用当前策略，`ToolResultCompressionStrategy.Apply` 返回错误，避免向上下文写入不可读取 marker。
2. 工具名反查失败：不报错，使用默认 `50,000` 字符。
3. 工具声明值非法：不报错，使用默认 `50,000` 字符。
4. 已压缩 result：继续跳过，避免重复落盘。
5. hash 命中已有 persisted output：继续复用已有 id。

## 8. 测试计划

### 8.1 tools metadata 测试

- `AllToolSpecs` 覆盖所有现有内置工具。
- `AllToolDefs` 与 `AllToolSpecs` 的工具名集合一致。
- 未显式配置的工具 `MaxResultSizeCharsForTool(name)` 返回 `50000`。
- 未知工具返回 `50000`。

### 8.2 compression 单测

新增或调整 `internal/agent/runtime/compression/compression_test.go`：

- 单个工具结果 `50000` 字符：不落盘。
- 单个工具结果 `50001` 字符：落盘并替换为 marker。
- 两个工具结果各 `40000` 字符：即使总量超过旧的全局阈值判断边界，也不因总量触发压缩。
- 一个工具结果超过阈值、一个未超过阈值：只压缩超限结果。
- 自定义 resolver 返回特殊上限，例如 `bash=10`：验证按工具名生效。
- 找不到 tool name：使用默认 `50000`。
- JSON wrapper 和 status 仍保持。
- 已包含 `<persisted-output>` 的 result 不重复落盘。
- CJK 内容按 rune 计数，避免用字节数误触发。

### 8.3 runtime 集成测试

新增或调整 runtime 测试：

- `Service` 构造压缩请求时传入 `ToolRegistry` 的 result limit resolver。
- request-only 压缩后，`state.History` / `state.ModelHistory` 仍保存完整结果，发给 LLM 的 request history 使用 marker。
- MCP 风格未知工具名使用默认 `50000`。

### 8.4 工具截断回归测试

- `bash` 返回超过 `50000` 字符的输出时，工具执行层不截断；压缩层负责落盘。
- `grep content` 结果超过 `50000` 字符时，工具执行层不截断；压缩层负责落盘。
- `web_fetch` 无 processor 时清洗文本超过 `50000` 字符，工具执行层不截断；`webFetchMaxBodySize` 仍生效。

## 9. 实施步骤建议

用户审阅通过后再实施代码，建议按以下顺序：

1. 新增 `internal/tools` 工具 metadata 和默认 `maxResultSizeChars` 查询。
2. 调整 `ToolRegistry` 持有并暴露工具结果上限 resolver。
3. 调整 `compression.Request` 和 runtime 构造压缩请求的入口，传入 resolver。
4. 改造 `ToolResultCompressionStrategy` 为 per-tool result limit。
5. 移除 `bash`、`grep content`、`web_fetch` 的返回前静默截断，保留资源边界限制。
6. 更新单测和集成测试。
7. 运行 `go test ./internal/tools ./internal/agent/runtime/... ./internal/local`。

## 10. 验收标准

1. 每个内置工具都有 `maxResultSizeChars` 声明路径，默认值为 `50,000`。
2. 未知工具或 MCP 工具默认按 `50,000` 字符处理。
3. 单个 tool result 超过对应上限时，完整结果落盘到现有 `~/.cynosure/task_outputs/{workspace}/{session_id}/tool-results/`。
4. LLM 请求上下文中只保留 `<persisted-output>` marker 和 2000 字符预览。
5. 单个 tool result 未超过对应上限时，不因多个工具结果总量过大而触发本策略压缩。
6. `read_persisted_output` 能读取落盘完整结果。
7. `MessageWindowCompressionStrategy`、`RecentToolResultRetentionStrategy`、memory 和 summarization 策略行为不变。
8. 工具执行日志 `~/.cynosure/task_outputs/{workspace}/{session_id}/tools.md` 继续追加完整工具结果。
9. 现有测试通过，新增测试覆盖 per-tool 阈值、默认值、未知工具、CJK 字符计数和内部截断清理。
