# reasoning_content 全轮次流式下发设计文档

## 1. 背景

当前 `go-agent` 已经把普通模型回答 `content` 收敛为“仅最后一轮不再调用工具时才流式下发”。这可以避免工具调用 round 中的中间回答泄漏到前端。

新的需求是将 `reasoning_content` 与普通 `content` 的策略分离：

- `content`：继续保持现状；工具调用 round 不流式下发，只在最后一轮不再调用工具时下发。
- `reasoning_content`：不论当前 round 是否会调用工具，都统一通过 `reasoning_delta` 流式下发给前端。

## 2. 当前实现分析

后端模型流读取集中在 `runModelRoundStream`：`go-agent/internal/web/runtime/conversation_flow.go:85`。

当前逻辑会同时收集：

- 普通回答 `content`：`go-agent/internal/web/runtime/conversation_flow.go:113`
- 推理内容 `reasoning_content`：`go-agent/internal/web/runtime/conversation_flow.go:118`
- 工具调用 `tool_calls`：`go-agent/internal/web/runtime/conversation_flow.go:123`

目前 `content` 和 `reasoning_content` 都会进入同一个 `bufferedDeltas`，并在 round 结束后统一经过 `shouldEmitModelDeltas` 判断：`go-agent/internal/web/runtime/conversation_flow.go:135`。

这意味着只要当前 round 是工具调用 round，`content` 和 `reasoning_content` 都不会下发。

前端已经分别消费两种事件：

- `assistant_delta` 追加到临时 assistant 消息的 `content`：`web/src/App.tsx:305`
- `reasoning_delta` 追加到临时 assistant 消息的 `reasoning_content`：`web/src/App.tsx:308`
- SSE 事件分发在 `web/src/api.ts:154`

因此本次主要改后端事件发送策略，不需要新增前端事件类型。

## 3. 目标与非目标

### 3.1 目标

1. 工具调用 round 中的 `reasoning_content` 也实时下发给前端。
2. 最后一轮不调用工具的 `reasoning_content` 继续实时下发给前端。
3. 工具调用 round 中的普通 `content` 继续不下发给前端。
4. 最后一轮不调用工具的普通 `content` 继续下发给前端。
5. 不改变 SSE 事件名，不新增前端协议字段，不新增数据库字段。

### 3.2 非目标

1. 不改变工具事件 `tool` 的展示策略。
2. 不改变工具调用 round 的普通 `content` 持久化与模型上下文行为。
3. 不调整前端推理区域 UI 样式。
4. 不在本次合并多轮 `reasoning_content` 到最终 assistant message；最终 assistant message 仍代表最后一轮模型输出。

## 4. 方案对比

### 方案 A：继续统一缓存，round 结束后分类回放

保留当前 `bufferedDeltas`，round 结束后：

- `reasoning_delta` 全部回放。
- `assistant_delta` 仅最终 round 回放。

优点：改动少。

缺点：工具调用 round 的 `reasoning_content` 仍需要等整个 round 结束后才能下发，不是真正实时流式。

结论：不推荐。

### 方案 B：`reasoning_content` 即时下发，`content` 继续缓冲【推荐】

调整 `runModelRoundStream`：

- 遇到 `choice.Delta.ReasoningContent` 时，写入 builder 并立即发送 `reasoning_delta`。
- 遇到 `choice.Delta.Content` 时，写入 builder 并缓存为 `assistant_delta`。
- round 结束后，如果确认不调用工具，则回放缓存的 `assistant_delta`。
- round 结束后，如果确认调用工具，则丢弃缓存的普通 `content` delta。

优点：

- `reasoning_content` 真正做到所有 round 统一实时流式下发。
- `content` 仍保持“只有最终轮下发”的保护策略。
- 前端协议不变。
- 改动范围集中在 runtime 流式逻辑和测试。

缺点：

- 前端临时 assistant 消息会累计多轮 reasoning，但最终 `assistant` 事件当前会用最后一轮 `reasoning_content` 覆盖临时内容。若产品希望最终保留前端累计的多轮 reasoning，需要后续单独设计前端合并策略。

结论：推荐。

### 方案 C：新增 `tool_reasoning_delta` 事件

工具调用 round 的 reasoning 使用新事件名，最终 round 仍使用 `reasoning_delta`。

优点：前端可区分 reasoning 来源。

缺点：需要新增前端协议和状态分支；当前需求是“统一给前端”，不要求区分来源。

结论：不推荐。

## 5. 推荐设计

采用方案 B：`reasoning_content` 即时下发，`content` 继续按最终轮策略缓冲回放。

### 5.1 后端流式策略

`runModelRoundStream` 中保留两个 builder：

- `content`：构建当前 round 的完整普通内容。
- `reasoningContent`：构建当前 round 的完整推理内容。

调整后的逻辑：

1. `choice.Delta.Content` 非空时：
   - 写入 `content`。
   - 追加到 `bufferedDeltas`，等待 round 结束后判断是否回放。
2. `choice.Delta.ReasoningContent` 非空时：
   - 写入 `reasoningContent`。
   - 如果 `state.Writer != nil`，立即发送 `reasoning_delta`。
3. round 结束后：
   - 计算 `calls := toolCalls.Calls()`。
   - 如果当前 round 不调用工具，回放 `bufferedDeltas` 中的普通 `assistant_delta`。
   - 如果当前 round 调用工具，不回放普通 `assistant_delta`。

### 5.2 判定函数命名

当前 `shouldEmitModelDeltas` 不再准确，因为它只控制普通 `content` delta，不再控制 `reasoning_delta`。

建议改名为：`shouldEmitAssistantContentDeltas`。

语义：

- `reasoning_delta`：不经过该函数，所有 round 即时发送。
- `assistant_delta`：经过该函数，仅最终 round 发送。

判定条件保持不变：`finishReason != tool_calls && len(toolCalls) == 0`。

### 5.3 SSE 事件顺序

不调用工具的直接回答：

1. `conversation`
2. `reasoning_delta`，如果模型输出 reasoning
3. `assistant_delta`，round 结束确认不调用工具后回放
4. `assistant`
5. `done`

先调用工具再最终回答：

1. `conversation`
2. `reasoning_delta`，工具调用 round 的 reasoning，实时下发
3. `tool`，工具执行事件
4. `reasoning_delta`，最终 round 的 reasoning，实时下发
5. `assistant_delta`，最终 round 的普通 content
6. `assistant`
7. `done`

工具调用 round 中的普通 `content` 不会出现 `assistant_delta`。

### 5.4 前端行为

前端无需新增事件类型。现有 `onReasoningDelta` 会把所有 `reasoning_delta` 追加到临时 assistant 消息：`web/src/App.tsx:308`。

注意：最终 `assistant` 事件当前会覆盖临时消息的 `reasoning_content`：`web/src/App.tsx:311`。本设计暂不调整该行为，因为用户当前要求是“流式输出统一给前端”，不是“最终消息保留所有轮次 reasoning”。如果审查时认为最终前端展示也必须保留所有轮次 reasoning，需要增加前端合并逻辑作为额外范围。

## 6. 测试设计

### 6.1 修改工具调用 round 测试

修改现有测试：`go-agent/internal/web/runtime/runtime_test.go:496`。

建议重命名为：`TestRespondToConversation_StreamsReasoningForToolRoundsButNotContent`。

断言：

- 工具调用 round 的 `reasoning_delta` 必须出现。
- 工具调用 round 的 `assistant_delta` 仍不能出现。
- 最终 round 的 `reasoning_delta` 必须出现。
- 最终 round 的 `assistant_delta` 必须出现。

### 6.2 修改 helper 测试

将 `TestShouldEmitModelDeltas` 改名为 `TestShouldEmitAssistantContentDeltas`。

表驱动用例保持不变，因为判定条件仍用于普通 `content`：

- `stop` 且无 tool calls：返回 true。
- `tool_calls` 且有 tool calls：返回 false。
- `stop` 但存在 tool calls：返回 false。
- 空 finish reason 且无 tool calls：返回 true。

### 6.3 保持既有测试

以下测试语义应保持：

- `TestRespondToConversation_StreamsAssistantContentDeltas`
- `TestRespondToConversation_StreamsReasoningContentDeltas`
- `TestRespondToConversation_StreamsToolCallsAcrossChunks`

## 7. 风险与兼容性

### 7.1 reasoning 会暴露工具调用前推理

这是本次需求的目标行为。普通 `content` 仍不会在工具调用 round 中下发，因此不会把中间回答展示成最终回答。

### 7.2 最终 assistant 事件覆盖临时 reasoning

当前前端收到最终 `assistant` 事件时，会用 payload 覆盖临时 assistant 消息。由于 payload 的 `reasoning_content` 是最终 round 的 reasoning，前端最终展示可能只保留最终 round reasoning。

本设计建议先不改前端；如果审查希望“最终展示也保留工具调用 round reasoning”，需要在实施中增加前端合并策略。

### 7.3 事件顺序变化

`reasoning_delta` 会比 `assistant_delta` 更早到达。尤其最终 round 中，reasoning 会实时下发，而普通 content 仍需等 round 结束确认不调用工具后回放。这与当前前端按字段分别追加的实现兼容。

## 8. 实施步骤

1. 修改 `go-agent/internal/web/runtime/conversation_flow.go`：
   - `ReasoningContent` chunk 即时发送 `reasoning_delta`。
   - `Content` chunk 继续缓存。
   - round 结束后只回放普通 `assistant_delta`。
   - 将 `shouldEmitModelDeltas` 改名为 `shouldEmitAssistantContentDeltas`。
2. 修改 `go-agent/internal/web/runtime/runtime_test.go`：
   - 更新工具调用 round 测试断言。
   - 更新 helper 测试名称。
3. 运行验证：`cd /Users/bytedance/golang_pro/nano_cc/go-agent && go test ./internal/web/runtime`。
4. 必要时运行全量验证：`cd /Users/bytedance/golang_pro/nano_cc/go-agent && go test ./...`。

## 9. 验收标准

1. 工具调用 round 中，模型返回 `reasoning_content` 时，前端收到 `reasoning_delta`。
2. 工具调用 round 中，模型返回普通 `content` 时，前端不收到 `assistant_delta`。
3. 最终 round 中，模型返回 `reasoning_content` 时，前端收到 `reasoning_delta`。
4. 最终 round 中，模型返回普通 `content` 时，前端收到 `assistant_delta`。
5. SSE 协议不新增事件名，前端 `api.ts` 无需调整。
6. `go test ./internal/web/runtime` 通过。
