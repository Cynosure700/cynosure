# LLM 调用前上下文压缩机制设计

## 背景

`go-agent` 的主 Agent 循环在 `RespondToConversation` 中加载历史、追加当前 user 消息、构造 system prompt，并在每轮 LLM 调用前通过 `compressContextBeforeLLM` 执行上下文压缩，再用压缩后的请求态上下文构造 OpenAI Chat messages，见 `internal/web/runtime/conversation_flow.go:49`、`internal/web/runtime/conversation_flow.go:52`。

工具执行结果当前以 `role=tool` 的消息追加到内存态 `LoopState.Messages` 和 `LoopState.History`，内容来自 `ToolExecutionOutcome.MessageContent()`，格式为 JSON：`{"status":"...","result":"..."}`，见 `internal/web/runtime/hooks/tool.go:71`、`internal/web/runtime/hooks/types.go:118`。会话最终通过 `SetConversationHistory` 写入 `conversations.history_json`，见 `internal/web/runtime/hooks/stop.go:9`、`internal/web/storage/conversations_repo.go:72`。

当工具输出或长会话历史很大时，这些内容会进入下一轮 LLM 请求，容易造成上下文膨胀、费用增加或触发模型上下文限制。因此需要在每次调用 LLM 前维护可扩展的上下文压缩流水线。压缩必须只影响“发给 LLM 的请求态上下文”，不能改变前端会话展示所依赖的真实历史；也就是说，大 tool_result 落库 marker、消息窗口裁剪、旧 tool_result 占位符、LLM 全量摘要都不能写回 `state.History` 或 `conversations.history_json`。

## 目标

1. 每次调用 LLM 前执行上下文压缩，而不是仅在会话结束时处理。
2. 压缩流水线执行顺序固定为：大 tool_result 落库压缩 → 消息数量窗口裁剪 → 旧 tool_result 占位压缩 → token 超限后的 LLM 全量摘要；四层压缩全部只作用于 `RequestHistory`。
3. 先统计“最后一条 user 消息之后、下一次 LLM 调用前”所有 `role=tool` 消息中 `result` 字段的总大小。
4. 当总大小超过 200KB 时：
   - 按单个 tool_result 大小从大到小排序。
   - 从最大的开始把完整输出落库。
   - 上下文中只保留 `<persisted-output>` 标记和前 2000 字符预览。
   - 逐个压缩，直到剩余内联 tool_result 总大小不超过 200KB；如果仍有超大单项，则继续压缩该单项。
5. 大 tool_result 落库压缩完成后，当历史消息数超过 50 条时，执行消息窗口裁剪：保留头部 3 条消息作为初始上下文，保留尾部 47 条消息作为当前工作上下文，中间消息从本次 LLM 请求上下文中裁掉。
6. 消息数量窗口裁剪完成后，只保留最近 3 条 `role=tool` 消息的完整 `result` 内容；更早的 tool_result 替换为一行占位符 `[Earlier result compacted. Re-run if needed]`。
7. 保持消息格式兼容：`tool` 消息仍然是 `{"status":"...","result":"..."}` JSON，只替换其中的 `result` 字符串。
8. 被大 tool_result 落库压缩裁剪的完整内容仍可由大模型按需读取：模型看到 `<persisted-output>` 标记后，可以调用运行时提供的读取工具从数据库取回指定范围的原文。
9. 如果上述三层请求态压缩后，发送给 LLM 的请求仍超过 token 限定值，则触发一次 LLM 全量摘要，把当前请求态 history 转换为“摘要 + 最近工作窗口”的请求态上下文。
10. 所有压缩都不能影响会话内容展示：压缩结果只用于本轮 LLM 请求，不替换前端展示所依赖的原始会话消息，不写入 `conversations.history_json` 作为用户可见历史。
11. 支持未来增加其他上下文压缩策略，例如旧轮次工具输出压缩、图片/文件引用替换等。

## 非目标

1. 本次消息数量裁剪只影响发给 LLM 的请求态上下文，不影响前端展示和最终持久化的真实会话历史；被消息数量裁掉的内容仍保留在 `state.History` / `conversations.history_json` 中供前端展示。
2. 本次旧 tool_result 占位压缩不做原文落库，也不提供按需读取能力；占位符明确提示需要时重新运行工具。
3. 本次不设计面向前端用户的“按 ID 查看落库原文”接口；读取能力只面向大模型工具调用链路。
4. 本次不改变前端工具结果展示逻辑。前端仍从真实 `tool` 消息解析 `status/result`，展示原始工具输出；`<persisted-output>` marker 和 `[Earlier result compacted. Re-run if needed]` 只出现在 LLM 请求态上下文中。
5. 本次 LLM 全量摘要只作为三层请求态压缩后仍超 token 的兜底策略，不替代前三层确定性压缩。
6. 本次不把任何压缩产物作为前端用户可见消息展示，也不改变历史消息列表的展示顺序和原始内容。

## 当前链路梳理

### LLM 调用链路

- `RespondToConversation` 负责主 Agent 循环：加载历史、追加当前 user 消息、构造 system prompt、构建 OpenAI messages，并在循环内调用 LLM，见 `internal/web/runtime/conversation_flow.go:27`。
- 初始 OpenAI messages 由 `buildOpenAIMessages` 从 `storage.Message` 转换，见 `internal/web/runtime/prompt_builder.go:67`。
- 每一轮 LLM 请求当前直接使用 `state.Messages`，见 `internal/web/runtime/conversation_flow.go:49`。

### tool_result 表示

- 模型返回 tool calls 后，会先把 assistant tool call 消息追加到 `state.History`，见 `internal/web/runtime/conversation_flow.go:72`。
- 工具执行完成后，`appendToolMessageHook` 会追加 `role=tool` 消息，`ToolCallID` 对应 assistant tool call，内容为 `ToolExecutionOutcome.MessageContent()` 生成的 JSON，见 `internal/web/runtime/hooks/tool.go:71`。
- 历史消息最终整体编码进 `conversations.history_json`，现有 schema 使用 `LONGTEXT`，见 `internal/web/storage/migrations/001_init.sql:36`。

## 设计方案

### 方案对比

#### 方案 A：在 `buildOpenAIMessages` 内直接压缩

优点：改动点集中，所有 LLM 请求都会经过转换函数。  
缺点：`buildOpenAIMessages` 当前是纯转换函数；在其中落库会引入副作用、需要传入 `context/store/user/conversation`，职责会变重，也不利于未来增加多种压缩策略。

#### 方案 B：在 tool post hook 中立即压缩

优点：工具结果产生时即可处理，避免大内容进入内存后续链路。  
缺点：需求是统计“最后一条 user 消息里所有 tool_result 的总大小”，单个 hook 只能看到当前工具结果，无法方便地基于本轮所有 tool_result 的总量做排序和选择；并行或多工具调用时策略不清晰。

#### 方案 C：在每轮 LLM 请求前执行压缩流水线（推荐）

优点：准确满足“调用 LLM 前”与“统计最新用户轮次所有 tool_result 总大小”；可以在一个位置集中执行多条策略；策略可以基于真实 `state.History` 的请求态副本做全局决策；压缩后再重建 `state.Messages`，避免真实展示态和请求态混用。  
缺点：需要新增压缩抽象、落库表和 repository 方法，初始改动略多。

推荐采用方案 C。

#### 方案 D：消息数量裁剪作为大输出落库之后的窗口策略（推荐）

优点：与方案 C 的“LLM 请求前统一压缩”一致，能在一个入口明确策略顺序；大 tool_result 先落库，确保最新用户轮次中的超大输出在后续窗口裁剪和旧结果占位前已有可恢复记录；随后把历史消息窗口裁到最多 50 条。  
缺点：会丢失中间历史的直接上下文；如果被保留窗口边界切在 assistant tool call / tool result 对中间，需要额外做 OpenAI tool 消息合法性修复，避免请求失败。

推荐在默认压缩器中按以下顺序执行：

1. `ToolResultCompressionStrategy`
2. `MessageWindowCompressionStrategy`
3. `RecentToolResultRetentionStrategy`

#### 方案 E：仅保留最近 3 条 tool_result 完整内容（新增，推荐）

优点：相比按字节阈值压缩，它能稳定限制工具结果数量，避免很多小型旧 tool_result 累积占用上下文；替换为固定一行占位符，不需要额外数据库写入，也不会引入恢复接口复杂度。  
缺点：被替换为占位符的更早 tool_result 原文不再可由模型读取，只能根据占位符提示重新运行工具；如果旧结果仍然重要，模型需要重新调用对应工具或用户重新提供信息。已先被大 tool_result 策略替换成 `<persisted-output>` 的内容仍按 marker 语义读取。

推荐把该策略放在大 tool_result 落库压缩和消息数量裁剪之后：

1. 先对最新用户轮次内超过 200KB 阈值的大 tool_result 执行落库压缩。
2. 再由消息窗口策略裁掉中间历史。
3. 最后在保留窗口内只保留最近 3 条 tool_result 的完整 result，更旧的替换成 `[Earlier result compacted. Re-run if needed]`。

#### 方案 F：三层确定性压缩后仍超限时，生成请求态 LLM 全量摘要（新增，推荐）

优点：只在确定性压缩不足以满足模型上下文限制时触发，避免平时额外 LLM 调用；摘要覆盖三层压缩后的完整 active context，能保留任务目标、关键决策、文件路径、工具调用结论和未完成事项；摘要只参与 LLM 请求态上下文，不写回会话展示历史，满足“不影响会话内容展示”。  
缺点：需要额外一次 LLM 调用，摘要存在遗漏或偏差风险；需要新增 token 估算、摘要缓存/失效和请求态上下文分离，避免摘要消息污染 `state.History`。

推荐把该策略放在三层压缩之后：

1. 先执行大 tool_result 落库、消息窗口裁剪、旧 tool_result 占位三层确定性压缩。
2. 用 token 估算器评估即将发送给主模型的 `system prompt + messages + tool definitions` 是否超过预算。
3. 若未超限，直接使用三层压缩后的 `RequestHistory` 构造请求。
4. 若超限，调用摘要模型对三层压缩后的 `RequestHistory` 做全量摘要，生成请求态 synthetic context。
5. 主 LLM 请求使用 synthetic context；`state.History` 仍保持未压缩的真实展示消息，不被摘要或前三层压缩结果替换。

#### 方案 G：四层压缩全部 request-only，不写回展示历史（新增，推荐）

优点：彻底满足“压缩结果不影响前端展示”。前端和 `conversations.history_json` 始终保留真实 user/assistant/tool 消息；LLM 请求前从真实历史复制一份 `RequestHistory`，在副本上执行四层压缩。这样用户仍能看到完整工具结果和完整对话，模型仍能获得压缩后的可控上下文。  
缺点：真实历史会继续占用数据库空间；由于真实历史不保存 marker/占位符，下一轮请求需要重新从真实历史派生请求态压缩结果。因此需要给大 tool_result 落库和摘要增加按 hash 的缓存查询，避免重复落库或重复摘要。

推荐采用方案 G 作为整体语义：

1. `state.History` 永远表示“真实展示历史”，只由用户提交、模型响应、工具执行和 Stop hook 修改。
2. `compressContextBeforeLLM` 不再直接修改 `state.History`，而是返回 `CompressionOutput{RequestHistory, Artifacts}`。
3. 四层压缩只修改 `RequestHistory` 副本。
4. `SetConversationHistory` 仍只写真实 `state.History`，因此前端展示不出现 marker、占位符、窗口裁剪缺口或摘要 synthetic message。
5. 压缩副作用只允许写内部 artifact 表，例如 `persisted_outputs`、`context_summaries`；这些表不参与前端消息列表展示。

## 核心架构

### 新增压缩流水线

在 `internal/web/runtime` 下新增上下文压缩模块，建议文件：

- `context_compression.go`：压缩器入口、策略接口、公共配置。
- `message_window_compression.go`：消息数量窗口裁剪策略。
- `recent_tool_result_retention.go`：仅保留最近 3 条 tool_result 完整内容的占位压缩策略。
- `tool_result_compression.go`：本次 tool_result 压缩策略。
- `full_history_summarization.go`：三层压缩后仍超 token 时的 LLM 全量摘要策略。
- `token_estimator.go`：请求 token 预算估算，避免直接依赖具体模型 SDK 的内部实现。

核心接口建议：

```go
type ContextCompressionStrategy interface {
    Name() string
    Apply(ctx context.Context, req *ContextCompressionRequest) (ContextCompressionResult, error)
}

type ContextCompressionRequest struct {
    Conversation storage.Conversation
    User         storage.User
    History      []storage.Message
    Store        contextCompressionStore
    Config       ContextCompressionConfig
}

type ContextCompressionResult struct {
    DisplayHistory []storage.Message // 真实展示历史；压缩策略不得修改
    RequestHistory []storage.Message // 仅用于本轮 LLM 请求；四层压缩只修改这里
    Mutated        bool
    RequestOnly    bool
    Compressed     []CompressedOutputSummary
}
```

其中：

- `DisplayHistory` 表示真实会话历史，会在 `Stop` hook 中持久化，继续作为前端展示和后续会话加载的基础；压缩策略不得修改它。
- `RequestHistory` 表示本轮请求给 LLM 的上下文。默认是 `DisplayHistory` 的深拷贝；四层压缩都只修改这个副本。
- `RequestOnly=true` 是所有压缩策略的固定语义，表示压缩结果不应写回 `state.History` 或 `conversations.history_json`。

`Service` 持有一个默认压缩器：

```go
type Service struct {
    // existing fields...
    ContextCompressor *ContextCompressor
}
```

如果字段为空，则使用默认压缩器，默认按顺序注册 `ToolResultCompressionStrategy`、`MessageWindowCompressionStrategy`、`RecentToolResultRetentionStrategy`、`FullHistorySummarizationStrategy`。未来新增策略只需要实现 `ContextCompressionStrategy` 并加入默认策略列表。

### 调用位置

在主 Agent 循环内、构造 `ChatCompletionRequest` 前执行：

```go
for {
    round++
    requestHistory, err := s.compressContextBeforeLLM(ctx, state)
    if err != nil {
        return storage.Message{}, err
    }
    state.Messages = buildOpenAIMessages(state.SystemPrompt, requestHistory)
    req := openai.ChatCompletionRequest{...}
    ...
}
```

这样：

1. 第一轮请求前通常没有 tool_result，策略快速 no-op。
2. 工具执行完成并追加 tool 消息后，下一轮请求前会先执行大 tool_result 落库压缩。
3. 当历史消息数超过 50 条时，再裁掉中间历史，保留头 3 + 尾 47。
4. 在裁剪后的窗口中，最后只保留最近 3 条 tool_result 的完整 result，更旧 tool_result 替换为固定占位符。
5. 前三层压缩只修改 `requestHistory` 副本，不修改 `state.History`。
6. 如果三层压缩后仍超过 token 预算，摘要策略继续只改写 `requestHistory`；`state.History` 不写入摘要消息，避免影响会话展示。
7. `state.Messages` 使用 `requestHistory` 重建；Stop hook 和前端展示继续使用真实 `state.History`。
8. 如果模型认为 `<persisted-output>` 预览不足以回答问题，可在后续轮次调用 `read_persisted_output` 工具按需读取数据库中的原文片段；如果看到 `[Earlier result compacted. Re-run if needed]`，则需要重新运行对应工具。

### 展示历史与请求历史分离

新增明确边界：

```text
state.History / conversations.history_json / 前端展示
    = DisplayHistory，真实消息，不包含任何压缩产物

LLM ChatCompletionRequest.Messages
    = system prompt + RequestHistory，压缩副本，可能包含 marker / 占位符 / 摘要 synthetic message
```

`compressContextBeforeLLM` 的职责改为：

1. 深拷贝 `state.History` 得到 `requestHistory`。
2. 在 `requestHistory` 上执行大 tool_result 落库 marker 替换。
3. 在 `requestHistory` 上执行消息窗口裁剪和 tool call 边界修复。
4. 在 `requestHistory` 上执行旧 tool_result 占位替换。
5. 必要时把 `requestHistory` 转换为摘要 synthetic context。
6. 返回 `requestHistory` 给主 LLM 请求。
7. 不修改 `state.History`。

为避免浅拷贝修改 `ToolCalls` slice 等嵌套结构时污染展示历史，拷贝必须是深拷贝：每条 `storage.Message`、`ToolCalls` slice 和内部 function call 字段都要复制。

## MessageWindowCompressionStrategy 细节

### 触发条件

当请求态副本 `RequestHistory` 的消息数量 `len(requestHistory) > 50` 时触发。这里统计的是业务历史消息，不包含运行时额外插入的 system prompt；system prompt 由 `buildOpenAIMessages` 单独加在最前面。

### 裁剪规则

1. 保留 `history[:3]` 作为头部初始上下文。
2. 保留 `history[len(history)-47:]` 作为尾部当前工作上下文。
3. 丢弃中间区间 `history[3 : len(history)-47]`。
4. 当 `len(history) <= 50` 时 no-op。

裁剪后的基础结果最多 50 条消息，顺序保持不变：

```text
head[0:3] + tail[last 47]
```

### 与 OpenAI tool message 约束的关系

OpenAI chat messages 对工具调用有结构要求：带 `tool_calls` 的 assistant 消息需要后续存在对应 `role=tool` 消息；`role=tool` 消息也需要有前置 assistant tool call。简单拼接头 3 + 尾 47 时，切口可能落在一组 assistant tool call / tool result 中间，导致请求被模型 API 拒绝。

因此策略需要在裁剪后做一次“边界合法性修复”：

1. 删除保留窗口中找不到前置 assistant tool call 的孤儿 `role=tool` 消息。
2. 对保留窗口中找不到对应 tool result 的 assistant 消息，清空该 assistant 的 `ToolCalls` 字段；如果该 assistant 同时没有 `Content` 和 `ReasoningContent`，则删除该 assistant 消息。
3. 修复只处理裁剪产生的窗口边界问题，不改变完整保留区间内部已有的合法 tool call / tool result 对。

修复后最终消息数可能少于 50 条。这个取舍优先保证 LLM 请求合法性；如果严格要求一定保留 50 条且不做修复，则存在 API 拒绝风险。

### 持久化语义

`MessageWindowCompressionStrategy` 只修改 `RequestHistory`，不修改 `state.History`。后续 `Stop` hook 通过 `SetConversationHistory` 写回的仍是真实展示历史，因此被请求窗口裁掉的中间消息仍会保留在 `conversations.history_json` 中，前端展示不受影响。

如果后续希望模型也能主动恢复“请求态窗口裁掉但展示历史仍保留”的中间消息，可以另行设计 `read_conversation_history` 一类模型工具；本次不新增该工具。

### 与前后 tool_result 压缩的顺序

消息窗口裁剪在大 tool_result 落库压缩之后、旧 tool_result 占位压缩之前执行：

1. 先由 `ToolResultCompressionStrategy` 定位最新用户轮次 tool_result、统计总大小，并对超过阈值的大输出做落库 + marker 替换。
2. 再把 history 从 N 条裁到最多 50 条，并修复 OpenAI tool call / tool result 边界。
3. 最后在裁剪后的窗口内执行“只保留最近 3 条 tool_result 完整内容”的占位压缩。

这样可以保证最新用户轮次中的超大 tool_result 优先获得可恢复的 persisted output 记录，同时仍在最终发送给 LLM 前收敛消息数量和旧 tool_result 内容。

## RecentToolResultRetentionStrategy 细节

### 触发条件

在大 tool_result 落库压缩、消息窗口裁剪和 tool call 边界修复之后，统计当前 `history` 中所有 `Role == "tool"` 的消息数量。如果 tool 消息数量大于 3，则触发；否则 no-op。

这里的 “tool_result” 对应当前存储模型中的 `role=tool` 消息，不区分工具名称、状态或所属 user 轮次。

### 保留规则

1. 从 history 尾部向前扫描 `Role == "tool"` 的消息。
2. 最近的 3 条 tool 消息保留完整 `result` 内容，不做修改。
3. 第 4 条及更早的 tool 消息视为旧 tool_result；如果其 `result` 仍是完整原文，则替换为固定一行占位符：

```text
[Earlier result compacted. Re-run if needed]
```

4. 替换只影响 tool 消息的结果内容，不删除消息本身，也不删除对应 assistant tool call 元数据，避免破坏 OpenAI tool call / tool result 对。

### JSON 兼容

现有 tool 消息通常是 `ToolExecutionOutcome.MessageContent()` 生成的 JSON：

```json
{"status":"success","result":"原始工具输出"}
```

旧 tool_result 占位压缩应保持 JSON 结构不变，只替换 `result` 字段：

```json
{"status":"success","result":"[Earlier result compacted. Re-run if needed]"}
```

如果遇到历史异常数据，tool message content 不是合法 JSON，则只在 `RequestHistory` 中把整个 `Content` 替换为占位符。这样不会新增请求态 JSON 解析失败风险；真实展示历史保留原始异常内容不变。

### 幂等规则

如果 tool result 已经是以下任一压缩形态，则视为已压缩，不重复处理：

1. `result == "[Earlier result compacted. Re-run if needed]"`
2. `result` 包含 `<persisted-output` marker

最近 3 条的判定按 tool 消息本身的位置计算，而不是按“仍未压缩的 tool result”计算。也就是说，如果最近 3 条中有一条已经是 `<persisted-output>`，它仍然占用最近 3 条名额；策略不会为了寻找 3 条未压缩结果而向更早历史扩张。

### 与大 tool_result 落库压缩的顺序

该策略必须在 `ToolResultCompressionStrategy` 和 `MessageWindowCompressionStrategy` 之后执行：

1. `ToolResultCompressionStrategy` 先处理最新用户轮次内超过 200KB 阈值的大输出，并生成可恢复的 `<persisted-output>` marker。
2. `MessageWindowCompressionStrategy` 再裁剪消息窗口，保证后续只处理最终发送给 LLM 的保留窗口。
3. 旧 tool_result 占位策略最后运行，将保留窗口中早于最近 3 条的完整 tool result 替换为一行占位符。

这符合整体执行顺序“大 tool_result 落库压缩 → 消息数量窗口裁剪 → 旧 tool_result 占位压缩”。

### 与 read_persisted_output 的关系

`[Earlier result compacted. Re-run if needed]` 不包含 persisted output ID，因此不能通过 `read_persisted_output` 读取原文。模型若需要该旧结果，应重新调用对应工具，或请求用户重新提供必要信息。若旧 tool result 在该策略运行前已经被大 tool_result 落库策略替换为 `<persisted-output>`，则仍可按 marker 中的 ID 调用 `read_persisted_output`。

## ToolResultCompressionStrategy 细节

### 最新用户轮次定义

当前存储模型中没有“user 消息内部 content blocks”的结构，tool_result 表示为独立 `role=tool` 消息。因此本设计把“最后一条 user 消息里所有 tool_result”映射为：

> 从 `RequestHistory` 中最后一条 `Role == "user"` 消息开始，到当前请求态历史末尾之间的所有 `Role == "tool"` 消息。

该定义与当前主循环一致：当前 user 消息之后会出现 assistant tool call 消息和对应 tool 消息；在下一次 LLM 调用前，这些 tool 消息就是最新用户轮次内的工具结果。

### 大小统计

1. 只统计 `role=tool` 消息中 JSON content 的 `result` 字段。
2. 大小按 UTF-8 字节数计算，即 `len([]byte(result))`，阈值为 `200 * 1024`。
3. 如果 tool content 不是合法 JSON，则把整个 `Content` 当作 result 统计和替换，兼容历史异常数据。
4. 已经包含 `<persisted-output` 标记的 result 视为已压缩，不再重复落库。

### 压缩选择

1. 找到最新用户轮次内所有未压缩 tool_result。
2. 计算总字节数 `totalBytes`。
3. 如果 `totalBytes <= 200KB`，直接 no-op。
4. 如果超过阈值，按单个 result 字节数降序排序。
5. 从最大项开始逐个落库并替换上下文，直到“未压缩内联 result 总大小”不超过 200KB。

### 上下文替换格式

为了保持 tool 消息仍可被解析为 JSON，只替换 `result` 字段，格式示例：

```text
<persisted-output id="po_abc123" kind="tool_result" original_bytes="524288" preview_chars="2000" retrieval_tool="read_persisted_output">
完整输出已持久化；如需更多内容，请调用 read_persisted_output(id="po_abc123", offset=0, limit=20000) 分段读取。

这里是原始 tool_result 的前 2000 字符预览……
</persisted-output>
```

说明：

- `id`：落库记录 ID，便于未来查询原文。
- `kind`：当前固定为 `tool_result`，未来可扩展。
- `original_bytes`：原始 result 的 UTF-8 字节数。
- `preview_chars`：预览字符数，本次固定最多 2000 个 Unicode rune。
- `retrieval_tool`：提示模型使用哪个工具读取完整内容。
- 预览使用原始输出前 2000 字符，不做摘要，不额外调用 LLM。
- 标记正文第一段显式告诉模型：原文没有丢失，只是移入数据库；如果预览不足，应按需调用读取工具。

替换后的完整 tool 消息仍类似：

```json
{
  "status": "success",
  "result": "<persisted-output id=\"po_abc123\" kind=\"tool_result\" original_bytes=\"524288\" preview_chars=\"2000\" retrieval_tool=\"read_persisted_output\">\n完整输出已持久化；如需更多内容，请调用 read_persisted_output(id=\"po_abc123\", offset=0, limit=20000) 分段读取。\n\n...\n</persisted-output>"
}
```

## FullHistorySummarizationStrategy 细节

### 触发条件

该策略是兜底策略，只在前三层确定性压缩全部执行完成后触发判断：

1. 使用三层压缩后的 `RequestHistory` 构造候选 `ChatCompletionRequest`。
2. 估算 `system prompt + history messages + tool definitions` 的总 token 数。
3. 如果估算值 `<= ContextTokenBudget`，策略 no-op，继续使用三层压缩后的 `RequestHistory`。
4. 如果估算值 `> ContextTokenBudget`，触发 LLM 全量摘要。

`ContextTokenBudget` 建议按配置推导：

```text
ContextTokenBudget = ModelContextLimit - MaxResponseTokens - SafetyMarginTokens
```

默认值建议：

- `ModelContextLimit`：优先从配置读取；如果未配置，使用保守默认值，例如 128k。
- `MaxResponseTokens`：优先从 LLM 配置读取；如果当前项目没有显式字段，先用 8k 作为保守预留。
- `SafetyMarginTokens`：默认 4k，避免估算偏差导致真实请求仍超限。

### token 估算

本项目当前没有模型 tokenizer 依赖，第一版建议使用保守估算：

```text
estimatedTokens = ceil(utf8Bytes / 3)
```

估算对象包括：

1. system prompt。
2. 每条请求消息的 `role/content/reasoning_content/tool_call_id/tool_calls` JSON 表示。
3. 当前暴露给模型的 tool definitions JSON 表示。

该估算偏保守，可能更早触发摘要，但可以降低超上下文失败风险。未来如果引入模型精确 tokenizer，只替换 `TokenEstimator` 实现，不改变压缩策略接口。

### 摘要输入

摘要输入使用三层压缩后的完整 `RequestHistory`，而不是原始未压缩展示历史：

1. 已落库的大 tool_result 以 `<persisted-output>` marker 表示，摘要需要保留 persisted output ID、工具名称线索和预览结论。
2. 已窗口裁剪掉的消息不再进入摘要输入；它们不属于本轮请求态 active context，但仍保留在真实展示历史中。
3. 已被旧 tool_result 占位的内容只保留占位语义，摘要应说明“更早工具结果已压缩，需要时重新运行工具”。

摘要 prompt 必须要求模型输出结构化 Markdown，至少包含：

- 当前用户目标与最新请求。
- 已完成操作和关键结论。
- 重要文件路径、函数名、命令、错误信息和决策。
- 工具调用结果摘要；遇到 `<persisted-output>` 时必须保留 ID 和读取方式。
- 未完成事项、待验证项和下一步建议。
- 不确定或被占位压缩的信息，明确标注需要重新读取或重新运行工具。

### 请求态上下文构造

摘要完成后，不把摘要写入 `state.History`，而是生成仅用于 LLM 请求的 `RequestHistory`：

```text
[
  system: "以下是为了满足上下文窗口限制生成的会话摘要，仅用于本次模型推理，不是用户发送的真实消息...",
  user:   "<conversation-summary>\n...摘要内容...\n</conversation-summary>",
  ...recentTailMessages
]
```

其中 `recentTailMessages` 用于保留最新交互的精确信息，建议按 token 预算从尾部向前选取：

1. 预留 `SummaryTargetTokens` 给摘要，默认 8k。
2. 预留 `RecentTailMinTokens` 给最近消息，默认 16k 或预算的 25%，取较小值。
3. 从三层压缩后的 `RequestHistory` 尾部向前累加消息，直到接近剩余预算。
4. 对选出的 tail 再运行一次 tool call 边界修复，避免请求非法。

如果 `summary + recentTailMessages` 仍超过预算，则继续减少 `recentTailMessages`；如果只剩摘要仍超过预算，则要求摘要模型用更小 `SummaryTargetTokens` 重试一次。重试后仍超限则返回错误，中止本轮 LLM 调用，避免发送必然失败的请求。

### 不影响会话展示的保证

为了保证全量摘要不影响会话内容展示，必须遵守以下边界：

1. 摘要消息不调用 `CreateMessage`，不追加到 `state.History`。
2. 摘要消息不写入 `conversations.history_json`。
3. `Stop` hook 仍然只持久化未压缩的真实 `state.History`。
4. 前端继续展示真实 user/assistant/tool 消息；摘要仅存在于本轮 LLM 请求和可选的内部摘要缓存中。
5. 如果需要缓存摘要，只能写入独立的内部表或内存字段，并且该表不参与会话消息列表查询。

建议新增内部缓存表 `context_summaries`，避免同一份超限请求态历史在多轮中反复摘要；具体 schema 见“落库设计”。缓存命中条件为 `conversation_id + user_id + source_history_sha256` 完全一致。只要三层压缩后的 `RequestHistory` 发生变化，hash 变化，缓存自然失效。

### 摘要 LLM 调用方式

摘要调用不应走主 Agent 的工具循环，避免递归触发工具调用或再次进入上下文压缩。建议新增内部方法：

```go
func (s *Service) summarizeHistoryForContext(ctx context.Context, req SummaryRequest) (SummaryResult, error)
```

调用约束：

1. 使用同一个 `config.Client`，但 `ChatCompletionRequest.Tools` 为空。
2. 使用专门的 summary system prompt，禁止模型调用工具或编造未提供的信息。
3. 不写 LLM round 日志中的完整原文；只记录请求/响应大小、hash、token 估算和错误状态，避免日志再次膨胀。
4. 摘要失败时返回错误并中止本轮主 LLM 调用；不要退回到超限原文请求。

### 幂等与一致性

1. 对同一份三层压缩后 `RequestHistory`，`source_history_sha256` 相同，优先复用已有摘要。
2. 摘要缓存只影响请求态上下文，不改变真实历史；即使缓存丢失，也只会导致重新摘要，不会影响前端展示。
3. 如果摘要生成成功但主 LLM 调用失败，摘要缓存仍可保留，供下一次相同 `RequestHistory` 使用。
4. 如果摘要生成失败，`state.History` 不应被修改。
5. 摘要策略不能覆盖 `<persisted-output>` marker 或 `[Earlier result compacted. Re-run if needed]` 的语义，必须在摘要中保留恢复方式或重新运行提示。

## 按需读取机制

### 设计目标

上下文压缩只减少 LLM 请求中的内联内容，不代表工具原始输出被丢弃。模型后续如果发现预览不足，应能通过工具读取 `persisted_outputs` 表里的原文。该能力必须满足：

1. **按需读取**：默认不把完整原文重新塞回上下文，只有模型显式调用工具时才读取。
2. **分段读取**：避免一次读取又把大内容全部放回上下文。
3. **权限隔离**：只能读取当前用户、当前会话下的 persisted output，不能凭 ID 越权读取其他会话数据。
4. **可追踪**：读取动作仍走现有 tool call 事件、审计和历史追加链路。

### 新增工具：read_persisted_output

新增工具 `read_persisted_output`，提供给主 Agent 使用。工具定义与 handler 放在 `internal/tools`，runtime 负责在执行工具前向 context 注入当前 user/conversation 范围内的 `PersistedOutputReader`：

```json
{
  "name": "read_persisted_output",
  "description": "Read a chunk of a persisted tool output by id when a <persisted-output> marker preview is insufficient. Only outputs from the current conversation are accessible.",
  "parameters": {
    "type": "object",
    "properties": {
      "id": {
        "type": "string",
        "description": "The persisted output id from the <persisted-output> marker, for example po_abc123."
      },
      "offset": {
        "type": "integer",
        "description": "Zero-based character offset to start reading from. Defaults to 0."
      },
      "limit": {
        "type": "integer",
        "description": "Maximum characters to return. Defaults to 20000 and is capped by the runtime."
      }
    },
    "required": ["id"]
  }
}
```

建议把该工具作为“压缩机制配套工具”自动加入主 Agent 的工具定义，而不是依赖 `WebAllowedTools` 配置。原因是：一旦上下文中出现 `<persisted-output>` 标记，读取工具就是恢复语义所必需的能力；如果被配置漏掉，模型只能看到预览，无法按需补取完整信息。

实现上按现有 tools 目录模式组织：

- 工具定义加入 `internal/tools/definitions.go`。
- handler 加入 `internal/tools/handlers.go`，具体读取逻辑放在 `internal/tools/persisted_output.go`。
- handler 不直接依赖 `internal/web/storage`，只依赖一个抽象 `PersistedOutputReader`，从 `context.Context` 获取。
- runtime 在 `ToolRegistry.Execute` 前把当前 `ToolContext` 中的读取器注入 context；读取器内部调用 `GetPersistedOutputForConversation(ctx, id, userID, conversationID)`，完成会话级权限校验。
- 子 Agent 是否暴露该工具需要单独控制。推荐第一阶段只暴露给主 Agent；子 Agent 没有父会话完整历史，且子 Agent 输出最终会摘要回主 Agent，通常不需要读取主会话 persisted output。

### 工具参数与返回值

读取参数：

- `id`：必填，来自 `<persisted-output id="...">`。
- `offset`：可选，按 Unicode 字符偏移，默认 0。
- `limit`：可选，按 Unicode 字符数限制，默认 20000；运行时设置硬上限，例如 100000，防止模型一次取回超大内容。

工具返回仍使用现有 tool result JSON 包装，即最终进入上下文的是 `ToolExecutionOutcome{Status, Result}`。其中 `Result` 建议是可读 JSON 字符串：

```json
{
  "id": "po_abc123",
  "kind": "tool_result",
  "tool_call_id": "call_xxx",
  "offset": 0,
  "limit": 20000,
  "returned_chars": 20000,
  "total_chars": 180000,
  "original_bytes": 524288,
  "next_offset": 20000,
  "has_more": true,
  "content_sha256": "...",
  "content": "本次读取到的原文片段..."
}
```

如果 `has_more=true`，模型可以继续调用：

```json
{"id":"po_abc123","offset":20000,"limit":20000}
```

### 权限与安全

`read_persisted_output` 必须使用当前 `ToolContext` 做强校验：

1. 根据 `id` 查询 `persisted_outputs`。
2. 校验 `output.UserID == toolCtx.User.ID`。
3. 校验 `output.ConversationID == toolCtx.Conversation.ID`。
4. 不允许通过参数指定任意 conversation_id/user_id。
5. 查询失败、越权或 ID 不存在时返回 rejected/error tool result，不泄露记录是否属于其他用户或其他会话。

### 对模型的提示

仅在 marker 中提示还不够，建议在 system prompt 的工具使用说明中补充一条通用规则：

> 当你在历史消息中看到 `<persisted-output ...>` 标记时，说明完整工具输出已存入数据库，当前上下文只包含预览。如果预览不足以完成用户请求，请调用 `read_persisted_output` 按 ID 和 offset 分段读取，不要臆测被省略的内容。

这条提示可以在默认 system prompt 构建时根据工具是否启用动态加入；也可以写入 `read_persisted_output` 的 tool description。推荐两者都做：marker 负责局部提示，tool description/system prompt 负责全局行为约束。

### 避免读取工具造成二次膨胀

读取工具会把片段作为新的 `role=tool` 消息追加到历史，因此它本身也可能成为后续压缩对象。建议：

1. `read_persisted_output` 的单次返回硬限制默认不超过 100000 字符。
2. 如果读取结果导致最新用户轮次 tool_result 总大小再次超过 200KB，下一次 LLM 调用前仍由同一压缩策略处理。
3. 对读取工具返回的内容也允许再次落库，但要在 metadata 中记录 `source_persisted_output_id`，避免排查时看不出二次持久化来源。
4. 第一阶段可以先不做跨记录去重；后续如果二次压缩频繁，再用 `content_sha256` 做复用。

## 落库设计

### 新表：persisted_outputs

建议新增表，用于保存被上下文压缩替换掉的原始输出：

```sql
CREATE TABLE IF NOT EXISTS persisted_outputs (
    id VARCHAR(64) PRIMARY KEY,
    conversation_id VARCHAR(64) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    message_id VARCHAR(64) NOT NULL,
    tool_call_id VARCHAR(128) NOT NULL DEFAULT '',
    source_persisted_output_id VARCHAR(64) NOT NULL DEFAULT '',
    kind VARCHAR(64) NOT NULL,
    strategy VARCHAR(128) NOT NULL,
    original_bytes INT NOT NULL,
    content_sha256 CHAR(64) NOT NULL,
    content LONGTEXT NOT NULL,
    preview TEXT NOT NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    KEY idx_persisted_outputs_conversation (conversation_id, created_at),
    KEY idx_persisted_outputs_message (message_id),
    KEY idx_persisted_outputs_tool_call (tool_call_id),
    KEY idx_persisted_outputs_source (source_persisted_output_id),
    CONSTRAINT fk_persisted_outputs_conversation FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    CONSTRAINT fk_persisted_outputs_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

`message_id` 对应被替换的 `storage.Message.ID`；`tool_call_id` 对应该 tool 消息的 `ToolCallID`。`source_persisted_output_id` 用于记录“读取工具返回内容被再次压缩”这类二次持久化来源，普通工具结果为空。`content_sha256` 用于请求态压缩复用已有落库记录，避免真实历史不写回 marker 后重复落库。

### 新表：context_summaries

建议新增表，用于缓存 LLM 全量摘要结果。该表是内部上下文优化数据，不参与会话消息展示：

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
    KEY idx_context_summaries_lookup (conversation_id, user_id, source_history_sha256),
    KEY idx_context_summaries_conversation (conversation_id, created_at),
    CONSTRAINT fk_context_summaries_conversation FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    CONSTRAINT fk_context_summaries_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

`source_history_sha256` 使用三层压缩后的 `RequestHistory` 规范化 JSON 计算，确保同一份请求上下文可以复用摘要；真实历史新增或变化后会重新派生请求态历史，hash 随之变化，摘要自然失效。

### 请求态 persisted output 复用

由于大 tool_result marker 不再写回真实历史，同一条真实 tool 消息在多次 LLM 请求前可能反复被识别为“需要落库压缩”。为避免重复写入 `persisted_outputs`，需要新增按消息和内容 hash 的查询：

```go
GetPersistedOutputByMessageHash(ctx context.Context, conversationID, userID, messageID, toolCallID, strategy, contentSHA256 string) (storage.PersistedOutput, error)
```

大 tool_result 压缩流程调整为：

1. 计算原始 result 的 `content_sha256`。
2. 先按 `conversation_id + user_id + message_id + tool_call_id + strategy + content_sha256` 查询已有 persisted output。
3. 命中则复用已有 ID 生成请求态 marker，不重复写库。
4. 未命中才创建新的 `persisted_outputs` 记录。

为提升查询效率，建议给 `persisted_outputs` 增加索引：

```sql
UNIQUE KEY uniq_persisted_outputs_message_hash (conversation_id, user_id, message_id, tool_call_id, strategy, content_sha256)
```

### storage model 与 repository

新增模型：

```go
type PersistedOutput struct {
    ID             string
    ConversationID string
    UserID         string
    MessageID      string
    ToolCallID     string
    SourcePersistedOutputID string
    Kind           string
    Strategy       string
    OriginalBytes  int
    ContentSHA256  string
    Content        string
    Preview        string
    CreatedAt      time.Time
}

type ContextSummary struct {
    ID                    string
    ConversationID        string
    UserID                string
    SourceHistorySHA256   string
    Strategy              string
    EstimatedTokensBefore int
    EstimatedTokensAfter  int
    Summary               string
    CreatedAt             time.Time
}
```

新增 repository 方法：

- `CreatePersistedOutput(ctx context.Context, output storage.PersistedOutput) error`
- `GetPersistedOutput(ctx context.Context, id string) (storage.PersistedOutput, error)`（内部基础查询）
- `GetPersistedOutputForConversation(ctx context.Context, id, userID, conversationID string) (storage.PersistedOutput, error)`（读取工具使用的权限收敛查询，推荐优先实现）
- `GetPersistedOutputByMessageHash(ctx context.Context, conversationID, userID, messageID, toolCallID, strategy, contentSHA256 string) (storage.PersistedOutput, error)`（请求态压缩复用已有落库记录，避免不写回 marker 后重复落库）
- `CreateContextSummary(ctx context.Context, summary storage.ContextSummary) error`
- `GetContextSummaryByHistoryHash(ctx context.Context, conversationID, userID, sourceHistorySHA256 string) (storage.ContextSummary, error)`

同时扩展 runtime 的 store interface，使压缩策略和读取工具可以只依赖最小接口。

## 幂等与一致性

1. 所有策略只修改 `RequestHistory`，不修改 `state.History`。
2. 大 tool_result 策略在请求态 marker 替换前，先查询 `persisted_outputs` 是否已有同一真实消息和同一内容 hash 的记录；命中则复用，未命中才写库。
3. 如果落库失败，返回错误并中止本轮 LLM 调用，避免请求态上下文里出现无法追溯的 marker。
4. `SetConversationHistory` 始终写入真实展示历史，不写入 marker、占位符、窗口裁剪结果或摘要 synthetic message。
5. 如果进程在落库成功后、本轮 LLM 请求完成前异常退出，可能存在一条本轮未使用的 `persisted_outputs` 记录；这是可接受的孤儿记录，后续可通过清理任务处理。本次不做清理任务。
6. LLM 全量摘要只写 `context_summaries` 内部缓存，不写 `conversations.history_json`，不追加可见消息。
7. 摘要缓存以三层压缩后 `RequestHistory` 的 SHA256 为幂等键；命中缓存时不重复调用摘要 LLM。

## 可扩展性

压缩策略按顺序运行，后续可以新增：

1. `OldToolResultCompressionStrategy`：压缩非最新用户轮次中的旧工具结果。
2. `IncrementalSummaryRefreshStrategy`：在已有全量摘要基础上增量合并新消息，减少反复全量摘要成本。
3. `ReasoningContentCompressionStrategy`：对历史 reasoning_content 做截断或移除。
4. `ArtifactReferenceStrategy`：把大文件、图片、二进制内容替换为 artifact 引用。

为避免策略互相覆盖，约定：

- 策略只修改自己识别的消息片段。
- 每条策略返回 `Mutated` 和压缩摘要，便于日志记录和测试断言。
- 策略执行顺序固定写在默认 compressor 中；当前顺序为大 tool_result 落库压缩在前、消息数量窗口裁剪居中、旧 tool_result 占位压缩随后、LLM 全量摘要兜底在最后。未来若需要配置化再扩展，不在本次范围内。

读取能力也按同样原则扩展：当前 `read_persisted_output` 只处理 `kind=tool_result` 的文本内容；未来如果压缩图片、文件或二进制 artifact，可以新增 `kind` 和对应读取/下载工具，但 marker 仍保持统一的 `<persisted-output id="..." kind="...">` 引用形态。

## 日志与观测

建议在每次发生压缩时记录结构化日志，内容包括：

- conversation_id
- strategy
- compressed_count
- total_before_bytes
- inline_after_bytes
- persisted_output_ids

消息数量窗口裁剪发生时额外记录：

- message_count_before
- message_count_after
- preserved_head_count
- preserved_tail_count
- removed_middle_count
- repaired_orphan_tool_count
- repaired_incomplete_tool_call_count

旧 tool_result 占位压缩发生时额外记录：

- tool_result_count_before
- full_tool_result_retained_count
- earlier_tool_result_compacted_count
- placeholder

LLM 全量摘要发生时额外记录：

- token_budget
- estimated_tokens_before_summary
- estimated_tokens_after_summary
- source_history_sha256
- context_summary_id
- summary_cache_hit
- summary_chars
- retained_tail_message_count
- request_only

不要把原始 tool_result 内容或摘要正文写入日志，避免日志膨胀或泄露敏感信息。

读取工具也需要记录轻量审计信息：`conversation_id`、`persisted_output_id`、`offset`、`limit`、`returned_chars`、`has_more`、`status`。同样不要把读取到的 `content` 写入日志。

## 测试计划

### 单元测试

1. `len(history) <= 50`：消息窗口策略 no-op，history 不变。
2. `len(history) == 51`：保留前 3 条和后 47 条，中间裁掉 1 条，顺序保持为 `head + tail`。
3. `len(history) > 50`：保留前 3 条和后 47 条，中间裁掉 `len(history)-50` 条。
4. 裁剪切口产生孤儿 `role=tool`：删除该孤儿 tool 消息，避免 OpenAI 请求非法。
5. 裁剪切口产生 assistant tool call 缺少对应 tool result：清空该 assistant 的 `ToolCalls`；若该 assistant 无文本内容则删除。
6. 大 tool_result 落库策略先于消息窗口策略：构造 60 条历史，其中最新用户轮次内将被窗口裁掉的中间 tool_result 很大，断言该大结果会先写入 `persisted_outputs`，随后再由窗口策略裁掉对应消息。
7. tool 消息数量 `<= 3`：旧 tool_result 占位策略 no-op，history 不变。
8. tool 消息数量 `> 3`：只保留最近 3 条 tool 消息的完整 result，更早 tool result 替换为 `[Earlier result compacted. Re-run if needed]`。
9. 旧 tool_result 占位策略保持合法 JSON tool content 的 `status` 字段不变，只替换 `result` 字段。
10. 非 JSON tool content 被旧 tool_result 占位策略替换为纯占位符文本。
11. 已是占位符或 `<persisted-output>` 的旧 tool result 不重复处理；最近 3 条中已压缩结果仍占用保留名额。
12. 策略顺序：构造 5 条 tool result，其中第 1 条很大、第 5 条很小，断言第 1 条会先进入大 tool_result 落库表，随后旧 tool_result 占位策略再处理保留窗口内早于最近 3 条的完整结果。
13. `total <= 200KB`：大 tool_result 不压缩，不落库，history 不变。
14. `total > 200KB` 且多个 tool_result：按大小降序压缩，直到剩余内联总量不超过 200KB。
15. 单个超大 tool_result：落库完整内容，上下文 result 替换为 marker + 2000 字符预览。
16. 已压缩 marker：不重复落库。
17. 非法 JSON tool content：按整个 content 处理并替换为 marker。
18. 最新 user 之前的 tool_result：不参与大 tool_result 策略统计。
19. 落库失败：返回错误，history 不应被部分替换。
20. `read_persisted_output` 正常读取：按 `id/offset/limit` 返回片段、`next_offset`、`has_more`。
21. `read_persisted_output` 权限校验：其他 user 或 conversation 的记录不可读取。
22. `read_persisted_output` limit 上限：超过硬上限时被截断到运行时最大值。
23. token 估算 `<= ContextTokenBudget`：全量摘要策略 no-op，继续使用三层压缩后的 `RequestHistory`。
24. token 估算 `> ContextTokenBudget`：调用摘要 LLM，返回 `summary + recentTailMessages` 请求态上下文。
25. 四层策略均不修改 `state.History`，不调用 `CreateMessage`，不影响最终 `SetConversationHistory` 的真实历史。
26. 摘要缓存命中：相同 `source_history_sha256` 复用 `context_summaries.summary`，不重复调用摘要 LLM。
27. 摘要失败：返回错误并中止本轮主 LLM 调用，history 不被部分修改。
28. 摘要输出仍超预算：缩小摘要目标重试一次；重试失败则返回错误。

### 集成测试

1. 在 `RespondToConversation` 的 fake LLM 测试中预置超过 50 条历史，触发下一轮 LLM 请求前压缩，断言请求 messages 中业务历史最多保留 50 条，且保留前 3 + 后 47；同时断言最终 `storedHistory` 仍保留完整真实历史。
2. 断言大 tool_result 落库压缩先于消息窗口裁剪：超过阈值的大结果先写入 `persisted_outputs`，随后请求 messages 中业务历史最多保留 50 条；真实历史中的 tool result 仍是原文。
3. 断言消息窗口裁剪后再执行旧 tool_result 占位压缩：裁剪后请求窗口内只保留最近 3 条 tool result 完整内容；真实历史不出现占位符。
4. 在 `RespondToConversation` 的 fake LLM 测试中模拟：第一轮返回 tool call，工具返回大输出，第二轮 LLM 请求前断言请求 messages 中 tool result 已被替换为 marker。
5. 断言最终 `storedHistory` 中保存的是真实未压缩展示历史，而 `persisted_outputs` fake store 收到了完整原文。
6. 模拟模型在下一轮看到 marker 后调用 `read_persisted_output`，断言工具返回数据库中原文片段，且后续模型请求包含该读取结果。
7. 断言小输出路径不影响现有历史工具结果传递测试，例如 `TestBuildOpenAIMessagesCarriesHistoricalToolResult` 的语义保持不变。
8. 在 `RespondToConversation` 的 fake LLM 测试中预置三层压缩后仍超预算的历史，断言主 LLM 收到的是摘要请求态上下文，而 `store.historyUpdates` 中没有摘要 synthetic message。
9. 断言同一份三层压缩后 `RequestHistory` 第二次触发摘要时命中 `context_summaries` 缓存，不重复调用摘要 LLM。
10. 断言摘要请求不携带 tools，避免摘要过程递归调用工具或进入主 Agent loop。
11. 断言同一条大 tool result 在多次请求态压缩中复用同一条 `persisted_outputs`，不因 marker 不写回真实历史而重复落库。

## 实施步骤建议

1. 新增 `PersistedOutput` 与 `ContextSummary` model、migration、repository 方法，并扩展相关 store interface。
2. 新增 `ContextCompressor`、策略接口、`DisplayHistory/RequestHistory` 分离返回结构和默认策略注册。
3. 实现 `ToolResultCompressionStrategy`，包含最新 user 轮次定位、JSON 解析、大小统计、排序、落库、marker 替换，并作为默认流水线第一步执行。
4. 实现 `MessageWindowCompressionStrategy`，包含 `>50` 触发、保留头 3 + 尾 47、裁掉中间消息、OpenAI tool 边界合法性修复，并作为默认流水线第二步执行。
5. 实现 `RecentToolResultRetentionStrategy`，包含 tool 消息倒序扫描、最近 3 条保留、更早 result 占位替换、JSON/非 JSON 兼容、幂等处理，并作为默认流水线第三步执行。
6. 新增 `read_persisted_output` 工具定义与执行逻辑，支持权限校验和分段读取；工具定义和 handler 放在 `internal/tools`，runtime 只注入当前 user/conversation 范围内的读取器。
7. 在 system prompt / marker / tool description 中提示模型：看到 marker 且预览不足时应调用读取工具，不要臆测省略内容；看到 `[Earlier result compacted. Re-run if needed]` 时需要重新运行对应工具。
8. 实现 `TokenEstimator` 和 `FullHistorySummarizationStrategy`，包含预算判断、摘要 LLM 调用、摘要缓存、recent tail 选取和 request-only 上下文生成。
9. 在 `RespondToConversation` 每轮 LLM 请求前调用压缩器，并使用返回的 `RequestHistory` 重建 `state.Messages`；四层压缩结果都不写入 `state.History`。
10. 补充大 tool_result 落库压缩、消息窗口裁剪、旧 tool_result 占位压缩、读取工具、全量摘要兜底、真实展示历史不变的单元测试与主循环集成测试。
11. 运行 `go test ./...` 验证。

## 待确认点

1. 本设计选择“压缩到剩余内联 tool_result 总量不超过 200KB”为停止条件；如果希望一旦超过阈值就压缩本轮所有 tool_result，需要在实施前调整。
2. 本设计保留前 2000 个 Unicode 字符；如果严格要求 2000 字节，需要改为按 UTF-8 字节安全截断。
3. 本设计把“最后一条 user 消息里的 tool_result”映射为“最后一条 user 消息之后的 tool 消息”；这是基于当前项目消息模型的实现口径。
4. 本设计建议 `read_persisted_output` 不受 `WebAllowedTools` 配置影响、自动随压缩机制暴露给主 Agent；如果希望所有工具都必须显式配置，需要调整 marker 策略或在压缩启用时强制校验配置。
5. 本设计按 Unicode 字符 offset/limit 分段读取；如果希望与 `original_bytes` 完全一致，可以改为字节 offset，但实现时需要额外保证 UTF-8 安全切片。
6. 本设计在消息数量裁剪后允许因 tool call 边界合法性修复导致最终保留消息数少于 50 条；如果必须严格保留 50 条，需要接受潜在 OpenAI tool 消息结构非法风险，或引入更复杂的边界扩张策略。
7. 本设计不为请求态消息窗口裁掉的中间历史提供模型读取工具；这些中间历史仍保留在前端展示历史中，但模型不能主动读取。如果希望模型后续也能恢复这部分内容，需要新增 `read_conversation_history` 等模型工具。
8. 本设计不为旧 tool_result 占位符提供数据库恢复能力；如果希望旧结果也可恢复，应改为复用 `persisted_outputs` 落库机制，而不是固定占位符。
9. 本设计把“最近 3 条 tool_result”定义为裁剪后 history 中按消息位置倒序的最近 3 条 `Role == "tool"` 消息，而不是最近 3 条未压缩或成功状态的 tool result。
10. 本设计第一版使用 `ceil(utf8Bytes/3)` 做保守 token 估算；如果需要严格贴合具体模型 tokenizer，应在实施前引入精确 tokenizer。
11. 本设计默认把摘要缓存写入独立 `context_summaries` 表但不做前端展示；如果完全不希望持久化摘要缓存，可以改为仅进程内缓存，但会增加重复摘要成本。
12. 本设计默认摘要失败时中止本轮 LLM 调用；如果希望降级为继续裁剪 recent tail 或提示用户新开会话，需要调整错误处理策略。
