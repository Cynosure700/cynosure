# LLM 调用前上下文压缩机制设计

## 背景

`go-agent` 的主 Agent 循环 `RespondToConversation` 在每一轮循环中直接用 `state.Messages` 构造 `ChatCompletionRequest` 并调用 LLM，目前不做任何上下文压缩，见 `internal/web/runtime/conversation_flow.go:49`。

工具执行结果以 `role=tool` 消息追加到 `state.Messages` 和 `state.History`，内容来自 `ToolExecutionOutcome.MessageContent()`，格式为 JSON `{"status":"...","result":"..."}`，见 `internal/web/runtime/hooks/types.go:120`。会话最终通过 `SetConversationHistory` 写入 `conversations.history_json`，见 `internal/web/storage/conversations_repo.go`。

当工具输出或会话历史很大时，这些内容会进入下一轮 LLM 请求，造成上下文膨胀、费用增加，甚至触发模型上下文上限。需要在每次调用 LLM 前执行压缩，且压缩只能影响“发给 LLM 的请求态上下文”，不能改变前端展示所依赖的真实历史。

## 目标

1. 每次调用 LLM 前执行上下文压缩，而不是仅在会话结束时处理。
2. 压缩流水线顺序固定：**大 tool_result 落库压缩 → 消息数量窗口裁剪 → 旧 tool_result 占位压缩 → token 超限后的 LLM 全量摘要**；四层全部只作用于 `RequestHistory`。
3. 先统计“最后一条 user 消息之后、下一次 LLM 调用前”所有 `role=tool` 消息中 `result` 字段的总大小。
4. 总大小超过 200KB 时按单项从大到小落库，上下文只保留 `<persisted-output>` 标记 + 前 2000 字符预览，逐个压缩到剩余内联总量 ≤ 200KB。
5. 历史消息数 > 50 条时裁剪：保留头部 3 条 + 尾部 47 条，中间从本次请求上下文裁掉。
6. 只保留最近 3 条 `role=tool` 消息的完整 `result`，更早的替换为占位符 `[Earlier result compacted. Re-run if needed]`。
7. 保持消息格式兼容：`tool` 消息仍是 `{"status":"...","result":"..."}` JSON，只替换其中 `result` 字符串。
8. 被落库压缩的完整内容可由模型按需读取：模型看到 `<persisted-output>` 后调用运行时工具从数据库取回原文。
9. 三层请求态压缩后仍超 token 限定值时，触发一次 LLM 全量摘要，把请求态 history 转换为“摘要 + 最近工作窗口”。

## 核心原则：展示历史与请求历史分离

```text
state.History / conversations.history_json / 前端展示
    = DisplayHistory，真实消息，不含任何压缩产物

LLM ChatCompletionRequest.Messages
    = system prompt + RequestHistory，压缩副本，可能含 marker / 占位符 / 摘要
```

- `state.History` 永远是真实展示历史，只由用户提交、模型响应、工具执行、Stop hook 修改。
- 每轮 LLM 请求前深拷贝 `state.History` 得到 `RequestHistory`，四层压缩只改这个副本。
- `SetConversationHistory` 始终只写真实 `state.History`，因此前端不出现 marker、占位符、窗口缺口或摘要消息。
- 压缩副作用只允许写内部表 `persisted_outputs`、`context_summaries`，这些表不参与前端消息列表展示。

> 深拷贝必须是深拷贝：每条 `storage.Message`、`ToolCalls` slice 及内部 function call 字段都要复制，避免污染展示历史。

## 调用位置

在主循环内、构造 `ChatCompletionRequest` 前执行压缩，见 `internal/web/runtime/conversation_flow.go:49`：

```go
for {
    round++
    requestHistory, err := s.compressContextBeforeLLM(ctx, state)
    if err != nil {
        return storage.Message{}, err
    }
    state.Messages = buildOpenAIMessages(state.SystemPrompt, requestHistory)
    req := openai.ChatCompletionRequest{
        Model:    s.Cfg.LLM.ModelID,
        Messages: state.Messages,
        Tools:    s.Tools.Definitions(),
    }
    // ...runModelRoundStream
}
```

`compressContextBeforeLLM` 职责：深拷贝 `state.History` → 依次执行四层策略 → 返回 `RequestHistory`，不修改 `state.History`。

## 压缩流水线

新增模块 `internal/web/runtime`，建议文件：`context_compression.go`（入口 + 策略接口）、`tool_result_compression.go`、`message_window_compression.go`、`recent_tool_result_retention.go`、`full_history_summarization.go`、`token_estimator.go`。

策略接口：

```go
type ContextCompressionStrategy interface {
    Name() string
    Apply(ctx context.Context, req *CompressionRequest) error // 原地修改 req.RequestHistory
}
```

默认按顺序注册四层策略；未来新增策略只需实现接口并加入列表。

### 第 1 层：ToolResultCompressionStrategy（大 tool_result 落库）

**最新用户轮次定义**：从 `RequestHistory` 中最后一条 `Role=="user"` 消息开始，到末尾之间的所有 `Role=="tool"` 消息。

**统计与选择**：
1. 只统计 `role=tool` 消息 JSON content 的 `result` 字段，大小按 UTF-8 字节 `len([]byte(result))`，阈值 `200*1024`。
2. content 非合法 JSON 时，把整个 `Content` 当作 result 统计与替换。
3. 已含 `<persisted-output` 标记的 result 视为已压缩，跳过。
4. `total <= 200KB` 直接 no-op；超限则按单项字节降序，从最大项开始逐个落库替换，直到剩余内联总量 ≤ 200KB。

**落库去重**：marker 不写回真实历史，同一真实消息可能多轮被识别为需落库。落库前按 `conversation_id + user_id + message_id + tool_call_id + strategy + content_sha256` 查已有记录，命中复用 ID，未命中才写库。

**替换格式**（只替换 `result` 字段，保持 JSON 可解析）：

```text
<persisted-output id="po_abc123" kind="tool_result" original_bytes="524288" preview_chars="2000" retrieval_tool="read_persisted_output">
完整输出已持久化；如需更多内容，请调用 read_persisted_output(id="po_abc123", offset=0, limit=20000) 分段读取。

这里是原始 tool_result 的前 2000 字符预览……
</persisted-output>
```

预览取原文前 2000 个 Unicode 字符，不做摘要、不额外调 LLM。落库失败返回错误并中止本轮 LLM 调用。

### 第 2 层：MessageWindowCompressionStrategy（消息窗口裁剪）

**触发**：`len(RequestHistory) > 50`（仅业务历史，不含 system prompt）。

**裁剪规则**：保留 `history[:3]` + `history[len-47:]`，丢弃中间，结果最多 50 条，顺序不变。

**OpenAI tool 边界修复**：拼接头尾后切口可能落在 assistant tool call / tool result 对中间，导致 API 拒绝，需修复：
1. 删除窗口中找不到前置 assistant tool call 的孤儿 `role=tool` 消息。
2. 对找不到对应 tool result 的 assistant 消息，清空其 `ToolCalls`；若该 assistant 同时无 `Content` 和 `ReasoningContent`，则删除该消息。
3. 仅处理边界，不动完整保留区内已合法的 tool call/result 对。

修复后最终条数可能少于 50，优先保证请求合法性。

### 第 3 层：RecentToolResultRetentionStrategy（旧 tool_result 占位）

**触发**：窗口裁剪后 `Role=="tool"` 消息数 > 3。

**保留规则**：从尾部倒序，最近 3 条 tool 消息保留完整 `result`；第 4 条及更早，把 `result` 替换为：

```text
[Earlier result compacted. Re-run if needed]
```

只替换 `result` 字段保持 JSON 结构：

```json
{"status":"success","result":"[Earlier result compacted. Re-run if needed]"}
```

content 非合法 JSON 时整体替换为纯占位符文本。

**幂等**：result 已是占位符或已含 `<persisted-output` 标记的视为已压缩，不重复处理；最近 3 条按消息位置计算（已压缩结果仍占名额，不向前扩张）。

**与读取的关系**：占位符不含 ID，不能 `read_persisted_output`，模型需要时应重新运行工具；若该结果已被第 1 层替换为 `<persisted-output>`，仍可按 marker 中 ID 读取。

### 第 4 层：FullHistorySummarizationStrategy（兜底 LLM 全量摘要）

**触发**：三层压缩后用 `RequestHistory` 估算 `system prompt + messages + tool definitions` 总 token，`> ContextTokenBudget` 才触发。

```text
ContextTokenBudget = ModelContextLimit - MaxResponseTokens - SafetyMarginTokens
```

默认：`ModelContextLimit` 优先读配置否则 128k；`MaxResponseTokens` 默认 8k；`SafetyMarginTokens` 默认 4k。

**token 估算**（第一版保守，无 tokenizer 依赖）：`estimatedTokens = ceil(utf8Bytes / 3)`，覆盖 system prompt、每条消息 JSON、tool definitions JSON。未来可只替换 `TokenEstimator` 实现。

**摘要输入**：三层压缩后的 `RequestHistory`。摘要 prompt 要求输出结构化 Markdown，至少含：当前用户目标、已完成操作与关键结论、重要文件/函数/命令/错误、工具调用结论（遇 `<persisted-output>` 必须保留 ID 和读取方式）、未完成事项与下一步、被占位压缩需重新运行的信息。

**请求态上下文构造**：

```text
[
  system: "以下是为满足上下文窗口限制生成的会话摘要，仅用于本次推理，不是用户真实消息...",
  user:   "<conversation-summary>...摘要...</conversation-summary>",
  ...recentTailMessages
]
```

`recentTailMessages` 从三层压缩后 history 尾部按 token 预算向前选取（预留摘要 `SummaryTargetTokens` 默认 8k，最近消息 `RecentTailMinTokens` 默认 16k 或预算 25% 取小），选完再做一次 tool 边界修复。仍超限则缩小重试一次，再失败则返回错误中止本轮调用。

**摘要 LLM 调用约束**：用同一 `config.Client` 但 `Tools` 为空、专用 summary system prompt 禁止调用工具、日志只记大小/hash/token 不记原文、失败返回错误且不修改 `state.History`。

**缓存**：以三层压缩后 `RequestHistory` 规范化 JSON 的 `source_history_sha256` 为幂等键，命中 `context_summaries` 则复用，不重复调用摘要 LLM。

## 按需读取：read_persisted_output

新增工具供主 Agent 使用，工具定义/handler 放 `internal/tools`（参考 `definitions.go`、`handlers.go`，读取逻辑放 `internal/tools/persisted_output.go`），runtime 在执行前向 context 注入当前 user/conversation 范围的 `PersistedOutputReader`。

```json
{
  "name": "read_persisted_output",
  "description": "Read a chunk of a persisted tool output by id when a <persisted-output> marker preview is insufficient. Only outputs from the current conversation are accessible.",
  "parameters": {
    "type": "object",
    "properties": {
      "id": {"type": "string", "description": "persisted output id from the marker, e.g. po_abc123"},
      "offset": {"type": "integer", "description": "zero-based char offset, default 0"},
      "limit": {"type": "integer", "description": "max chars, default 20000, capped by runtime"}
    },
    "required": ["id"]
  }
}
```

- 该工具作为压缩配套工具自动随压缩启用而暴露给主 Agent，不依赖 `WebAllowedTools` 配置（否则模型只能看预览无法补取）。
- 按 Unicode 字符 offset/limit 分段读取，单次硬上限默认 100000 字符。
- 返回仍是 tool result JSON 包装，`result` 为可读 JSON：含 `id/offset/limit/returned_chars/total_chars/next_offset/has_more/content`。
- 第一阶段只暴露给主 Agent，不暴露给子 Agent。

**权限校验**（强制）：按 `id` 查 `persisted_outputs`，校验 `UserID==toolCtx.User.ID` 且 `ConversationID==toolCtx.Conversation.ID`；不允许参数指定任意 conversation/user；越权或不存在返回 rejected，不泄露归属。

**系统提示**：在 system prompt 工具说明中加一条：看到 `<persisted-output>` 且预览不足时调用 `read_persisted_output` 分段读取，不要臆测被省略内容；看到 `[Earlier result compacted. Re-run if needed]` 时需重新运行对应工具。

## 落库设计

### 表 persisted_outputs

```sql
CREATE TABLE IF NOT EXISTS persisted_outputs (
    id VARCHAR(64) PRIMARY KEY,
    conversation_id VARCHAR(64) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    message_id VARCHAR(64) NOT NULL,
    tool_call_id VARCHAR(128) NOT NULL DEFAULT '',
    kind VARCHAR(64) NOT NULL,
    strategy VARCHAR(128) NOT NULL,
    original_bytes INT NOT NULL,
    content_sha256 CHAR(64) NOT NULL,
    content LONGTEXT NOT NULL,
    preview TEXT NOT NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    UNIQUE KEY uniq_po_message_hash (conversation_id, user_id, message_id, tool_call_id, strategy, content_sha256),
    KEY idx_po_conversation (conversation_id, created_at),
    CONSTRAINT fk_po_conversation FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    CONSTRAINT fk_po_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 表 context_summaries

```sql
CREATE TABLE IF NOT EXISTS context_summaries (
    id VARCHAR(64) PRIMARY KEY,
    conversation_id VARCHAR(64) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    source_history_sha256 CHAR(64) NOT NULL,
    strategy VARCHAR(128) NOT NULL,
    estimated_tokens_before INT NOT NULL,
    estimated_tokens_after INT NOT NULL,
    summary LONGTEXT NOT NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    KEY idx_cs_lookup (conversation_id, user_id, source_history_sha256),
    CONSTRAINT fk_cs_conversation FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    CONSTRAINT fk_cs_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### model 与 repository

新增 `storage.PersistedOutput`、`storage.ContextSummary` model，及 repository 方法：

- `CreatePersistedOutput(ctx, output) error`
- `GetPersistedOutputForConversation(ctx, id, userID, conversationID) (PersistedOutput, error)` — 读取工具用，含权限收敛
- `GetPersistedOutputByMessageHash(ctx, conversationID, userID, messageID, toolCallID, strategy, contentSHA256) (PersistedOutput, error)` — 落库去重
- `CreateContextSummary(ctx, summary) error`
- `GetContextSummaryByHistoryHash(ctx, conversationID, userID, sourceHistorySHA256) (ContextSummary, error)`

同时扩展 runtime store interface，使压缩策略和读取工具只依赖最小接口。

## 幂等与一致性

1. 四层策略只修改 `RequestHistory`，不修改 `state.History`，不调用 `CreateMessage`。
2. 大 tool_result 落库前按内容 hash 查重，命中复用，未命中才写。
3. 落库失败返回错误并中止本轮 LLM 调用，避免出现无法追溯的 marker。
4. `SetConversationHistory` 始终只写真实展示历史。
5. 进程异常退出可能留下未被使用的 `persisted_outputs` 孤儿记录，可接受，本次不做清理任务。
6. 摘要只写 `context_summaries` 缓存，以三层压缩后 `RequestHistory` 的 SHA256 为幂等键。

## 日志与观测

每次压缩记录结构化日志：`conversation_id`、`strategy`、`compressed_count`、`total_before_bytes`、`inline_after_bytes`、`persisted_output_ids`；窗口裁剪额外记 `message_count_before/after`、`removed_middle_count`、`repaired_*`；占位压缩记 `full_tool_result_retained_count`、`earlier_tool_result_compacted_count`；摘要记 `token_budget`、`estimated_tokens_before/after`、`source_history_sha256`、`summary_cache_hit`、`retained_tail_message_count`。**不写原始 tool_result 内容或摘要正文**。读取工具记 `persisted_output_id/offset/limit/returned_chars/has_more/status`，不记 content。

## 实施步骤

1. 新增 `PersistedOutput`/`ContextSummary` model、migration、repository，扩展 store interface。
2. 新增 `ContextCompressor`、策略接口、`compressContextBeforeLLM`（深拷贝 + 四层）。
3. 实现 `ToolResultCompressionStrategy`（定位轮次 / JSON 解析 / 统计排序 / 落库去重 / marker 替换）。
4. 实现 `MessageWindowCompressionStrategy`（>50 触发 / 头 3 尾 47 / tool 边界修复）。
5. 实现 `RecentToolResultRetentionStrategy`（倒序扫描 / 保留最近 3 / 占位替换 / JSON 兼容 / 幂等）。
6. 新增 `read_persisted_output` 工具定义、handler、`PersistedOutputReader` 注入与权限校验。
7. 在 system prompt / marker / tool description 加入读取/重跑提示。
8. 实现 `TokenEstimator` + `FullHistorySummarizationStrategy`（预算判断 / 摘要 LLM / 缓存 / tail 选取 / request-only）。
9. 在 `RespondToConversation` 每轮请求前调用压缩器，用 `RequestHistory` 重建 `state.Messages`。
10. 补单元测试 + 主循环集成测试，运行 `go test ./...` 验证。

## 测试要点

- 窗口：`<=50` no-op；`51` 裁中间 1 条；切口孤儿 tool 删除；assistant 缺 result 清空 ToolCalls。
- 大 tool_result：`<=200KB` no-op；超限降序压缩到 ≤200KB；单超大项落库 + marker + 2000 预览；已 marker 不重复落库；非法 JSON 整体替换；最新 user 之前的 tool 不统计；落库失败不部分替换；同一结果多轮复用同一记录。
- 占位：tool 数 `<=3` no-op；`>3` 仅保留最近 3 条完整；保持 status 只换 result；非法 JSON 替换为纯占位符；已压缩不重复处理。
- 摘要：估算 `<=budget` no-op；`>budget` 调摘要返回 `summary + tail`；缓存命中不重复调用；摘要请求不带 tools；失败中止不改 history；超预算缩小重试。
- 分离：四层均不改 `state.History`，最终 `SetConversationHistory` 仍为真实历史；`read_persisted_output` 正常分段读取 + 权限校验 + limit 上限。

## 待确认点

1. 停止条件为“剩余内联 tool_result 总量 ≤ 200KB”；若希望一旦超阈值就压缩本轮全部 tool_result，需调整。
2. 预览按 2000 个 Unicode 字符；若要求严格 2000 字节需改按 UTF-8 字节安全截断。
3. “最后一条 user 消息里的 tool_result”映射为“最后一条 user 之后的 tool 消息”（基于当前消息模型）。
4. 窗口裁剪后允许因边界修复导致最终少于 50 条；若必须严格 50 条需接受 API 非法风险或更复杂扩张策略。
5. token 估算第一版用 `ceil(utf8Bytes/3)`；需贴合具体 tokenizer 时再替换实现。
6. 摘要失败默认中止本轮调用；若希望降级（继续裁 tail / 提示新开会话）需调整。
