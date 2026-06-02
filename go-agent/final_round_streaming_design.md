# 模型最终轮输出流式下发设计文档

## 1. 背景

当前 `go-agent` 在模型流式响应过程中，会把每个模型 round 中的 `content` 与 `reasoning_content` 增量直接通过 SSE 下发给前端。由于 agent loop 允许模型先输出一段文本/推理，再以 `tool_calls` 结束当前 round，因此前端会看到“模型准备调用工具时的中间内容”。

本次需求要求：

1. 模型调用工具的 round 中，模型输出的 `content` 不要流式输出给前端。
2. 模型调用工具的 round 中，模型输出的 `reasoning_content` 也不要流式输出给前端。
3. 只有最后一轮模型不再调用工具时，才将该轮模型输出流式输出给前端。
4. 工具执行事件本身仍可按现有 `tool` SSE 事件展示，需求只限制模型文本/推理内容的流式下发。

## 2. 现状分析

### 2.1 后端模型 round 循环

`RespondToConversation` 负责 agent loop。每次循环调用一次模型：

- 构造 `ChatCompletionRequest`：`go-agent/internal/web/runtime/conversation_flow.go:38`
- 调用流式模型 round：`go-agent/internal/web/runtime/conversation_flow.go:44`
- 如果 `finishReason != "tool_calls"` 或没有工具调用，则进入 `Stop` hook 并返回最终 assistant 消息：`go-agent/internal/web/runtime/conversation_flow.go:53`
- 如果模型要求工具调用，则持久化中间 assistant tool-call message，再执行工具：`go-agent/internal/web/runtime/conversation_flow.go:60`

也就是说，**是否是最终 round，需要等当前 round 完整读完后才能根据 `finishReason` 和 `ToolCalls` 判断**。

### 2.2 当前泄漏点

`runModelRoundStream` 在读取每个模型流 chunk 时立即下发 SSE：

- `content` 增量写入 `assistant_delta`：`go-agent/internal/web/runtime/conversation_flow.go:102`
- `reasoning_content` 增量写入 `reasoning_delta`：`go-agent/internal/web/runtime/conversation_flow.go:109`
- 工具调用增量只进入 accumulator：`go-agent/internal/web/runtime/conversation_flow.go:116`
- round 结束后返回完整 `ChatCompletionMessage` 与 `finishReason`：`go-agent/internal/web/runtime/conversation_flow.go:128`

问题在于：当 `content/reasoning_content` chunk 到达时，后端还不知道当前 round 最终是否会以 `tool_calls` 结束，因此当前实现会提前把中间内容推给前端。

### 2.3 最终 assistant 事件

最终答案通过 Stop hook 持久化并发送：

- 持久化最终 assistant message：`go-agent/internal/web/runtime/hooks/stop.go:9`
- 发送 `assistant` SSE 事件，包含完整 `content` 和 `reasoning_content`：`go-agent/internal/web/runtime/hooks/stop.go:16`

该事件只在 `RespondToConversation` 判断当前 round 不再调用工具后触发，因此它天然符合“最终轮”的语义。

### 2.4 前端消费逻辑

前端当前会消费三类模型输出事件：

- `assistant_delta` 追加到临时 assistant 消息的 `content`：`web/src/App.tsx:305`
- `reasoning_delta` 追加到临时 assistant 消息的 `reasoning_content`：`web/src/App.tsx:308`
- `assistant` 用最终完整消息覆盖临时消息：`web/src/App.tsx:311`

如果后端保证工具 round 不发送 `assistant_delta` / `reasoning_delta`，前端无需改动即可停止展示工具调用前的模型中间内容。

## 3. 目标与非目标

### 3.1 目标

1. 工具调用 round 中的模型 `content` 不发送 `assistant_delta`。
2. 工具调用 round 中的模型 `reasoning_content` 不发送 `reasoning_delta`。
3. 最终 round 中的模型 `content` 和 `reasoning_content` 仍通过现有 delta 事件下发给前端。
4. 最终 `assistant` SSE 事件保持不变，仍包含完整 `content` 与 `reasoning_content`。
5. 工具调用、工具结果、历史持久化、后续模型上下文不受影响。
6. 尽量只改后端 runtime 层，不扩大到前端协议或数据库结构。

### 3.2 非目标

1. 不隐藏工具调用事件 `tool`。
2. 不改变 `tool_calls` 的执行、审计、持久化逻辑。
3. 不改变最终 assistant 消息中的 `reasoning_content` 是否展示；本需求仅约束工具调用 round 的流式输出。
4. 不新增 SSE 事件类型。
5. 不新增数据库迁移。

## 4. 方案对比

### 方案 A：前端忽略工具调用前的 delta

前端收到 `assistant_delta` / `reasoning_delta` 后先缓存，等收到 `assistant` 再决定是否展示。

优点：后端改动少。

缺点：

- 前端并不知道某些 delta 属于工具调用 round 还是最终 round。
- 后端已经把不应暴露的 `reasoning_content` 发到了浏览器，不满足“不要流式输出给前端”。
- 多轮工具调用时前端状态机复杂。

结论：不推荐。

### 方案 B：模型 round 读取时完全不发送 delta，只发送最终 `assistant`

删除 `assistant_delta` / `reasoning_delta` 下发逻辑，只保留最终 assistant 事件。

优点：最简单，绝不会泄漏工具 round 内容。

缺点：

- 破坏“最终轮仍流式输出”的体验。
- 前端现有流式展示能力会退化。

结论：不推荐。

### 方案 C：后端按 round 缓冲模型 delta，确认是最终 round 后再回放【推荐】

`runModelRoundStream` 读取模型流时不立即发送 `assistant_delta` / `reasoning_delta`，而是：

1. 按顺序缓存模型文本 delta：`assistant_delta` 对应 `content` chunk，`reasoning_delta` 对应 `reasoning_content` chunk。
2. 继续累积完整 `content`、`reasoning_content`、`tool_calls`。
3. 当前 round 读完后，根据 `finishReason` 和累积到的 `toolCalls` 判断是否为最终 round。
4. 如果当前 round 是工具调用 round，则丢弃缓存的模型文本 delta，不发送给前端。
5. 如果当前 round 是最终 round，则按缓存顺序将 `assistant_delta` / `reasoning_delta` 通过 SSE 回放给前端。
6. `Stop` hook 仍发送最终 `assistant` 事件。

优点：

- 后端不再把工具调用 round 的 `content/reasoning_content` 发给浏览器，满足需求。
- SSE 协议和前端处理逻辑不需要变化。
- 最终 round 仍以 delta 事件形式输出，只是会在确认该 round 不调用工具后下发。
- 改动范围集中在 `go-agent/internal/web/runtime/conversation_flow.go` 与测试。

缺点：最终 round 的 delta 会延迟到该 round 完整结束后再下发，不再是边生成边下发。

结论：推荐采用。

## 5. 推荐设计

采用方案 C：**模型 round 内先缓冲所有文本/推理 delta，round 结束后仅在确认当前 round 不调用工具时回放给前端。**

### 5.1 最终 round 判定

新增一个小 helper，避免判定条件散落：

```go
func shouldEmitModelDeltas(finishReason openai.FinishReason, toolCalls []openai.ToolCall) bool {
	return finishReason != openai.FinishReasonToolCalls && len(toolCalls) == 0
}
```

判定依据：

- `finishReason == openai.FinishReasonToolCalls`：明确是工具调用 round，不发送模型 delta。
- `len(toolCalls) > 0`：即使 provider 的 finishReason 异常，也视为工具调用 round，不发送模型 delta。
- 其余情况视为最终 round，允许发送模型 delta。

### 5.2 缓冲事件结构

在 runtime 层新增私有结构体即可，不影响外部接口：

```go
type bufferedModelDelta struct {
	Event   string
	Content string
}
```

建议同时增加事件名常量：

```go
const (
	assistantDeltaEvent = "assistant_delta"
	reasoningDeltaEvent = "reasoning_delta"
)
```

缓存必须保留 chunk 到达顺序，因为部分模型可能交错输出 reasoning 与 content。

### 5.3 `runModelRoundStream` 调整

目标文件：`go-agent/internal/web/runtime/conversation_flow.go`

当前逻辑中，`choice.Delta.Content` 和 `choice.Delta.ReasoningContent` 一边写 builder，一边写 SSE。调整为：

1. 仍写入 `content` / `reasoningContent` builder，保证最终消息与持久化不变。
2. 不立即调用 `state.Writer.Event(...)`。
3. 将 delta 追加到 `bufferedDeltas`。
4. stream 结束后拿到 `calls := toolCalls.Calls()`。
5. 如果 `shouldEmitModelDeltas(finishReason, calls)` 为 true，再按顺序发送缓存 delta。
6. 返回 `ChatCompletionMessage{Content, ReasoningContent, ToolCalls: calls}`。

伪代码如下：

```go
var bufferedDeltas []bufferedModelDelta

if choice.Delta.Content != "" {
	seenOutput = true
	content.WriteString(choice.Delta.Content)
	bufferedDeltas = append(bufferedDeltas, bufferedModelDelta{Event: assistantDeltaEvent, Content: choice.Delta.Content})
}

if choice.Delta.ReasoningContent != "" {
	seenOutput = true
	reasoningContent.WriteString(choice.Delta.ReasoningContent)
	bufferedDeltas = append(bufferedDeltas, bufferedModelDelta{Event: reasoningDeltaEvent, Content: choice.Delta.ReasoningContent})
}

calls := toolCalls.Calls()
if state.Writer != nil && shouldEmitModelDeltas(finishReason, calls) {
	for _, delta := range bufferedDeltas {
		_ = state.Writer.Event(delta.Event, map[string]any{"content": delta.Content})
	}
}

return openai.ChatCompletionMessage{
	Role:             "assistant",
	Content:          content.String(),
	ReasoningContent: reasoningContent.String(),
	ToolCalls:        calls,
}, finishReason, nil
```

### 5.4 Agent loop 与持久化行为

`RespondToConversation` 的主流程无需改变：

- 工具调用 round 仍追加到 `state.Messages`，供下一轮模型看到 assistant tool-call message。
- 工具调用 round 仍写入 `state.History`，保留完整模型上下文：`go-agent/internal/web/runtime/conversation_flow.go:60`
- 最终 round 仍走 `RunStop`，持久化并下发最终 assistant message：`go-agent/internal/web/runtime/conversation_flow.go:54`

这意味着：

- 只是停止“向前端流式展示”工具 round 的 `content/reasoning_content`。
- 不影响模型后续推理可见的上下文。
- 不影响历史会话过滤逻辑；带 `ToolCalls` 的中间 assistant message 当前不会作为聊天气泡返回给前端：`go-agent/internal/web/app/conversation_handlers.go:124`

### 5.5 SSE 协议

无需新增或删除事件类型，仍保留：

- `assistant_delta`：仅最终 round 的 `content` chunk。
- `reasoning_delta`：仅最终 round 的 `reasoning_content` chunk。
- `tool`：工具执行展示事件，保持现状。
- `assistant`：最终完整 assistant 消息。
- `done`：流结束。

事件顺序示例：

#### 场景一：不调用工具，直接最终回答

```text
conversation
reasoning_delta   # 如果模型输出 reasoning_content
assistant_delta   # 如果模型输出 content
assistant         # final=true，完整 content/reasoning_content
done
```

#### 场景二：先调用工具，再最终回答

```text
conversation
tool              # 工具执行完成事件
reasoning_delta   # 仅最终 round 的 reasoning_content
assistant_delta   # 仅最终 round 的 content
assistant         # final=true，完整最终答案
done
```

工具调用 round 中即使模型返回了 `content` 或 `reasoning_content`，也不会出现对应 delta。

## 6. 测试设计

### 6.1 新增测试：工具调用 round 的 content/reasoning 不下发

文件：`go-agent/internal/web/runtime/runtime_test.go`

新增用例建议命名：

```go
func TestRespondToConversation_DoesNotStreamToolRoundContentOrReasoning(t *testing.T)
```

测试构造：

1. 第一轮模型流式返回：
   - `ReasoningContent: "先思考要不要调用工具"`
   - `Content: "我先查一下"`
   - `ToolCalls: load_skill`，`FinishReason: tool_calls`
2. 第二轮模型流式返回：
   - `ReasoningContent: "工具已返回，组织最终答案"`
   - `Content: "最终答案"`
   - `FinishReason: stop`
3. 断言：
   - writer events 中不存在内容为 `"先思考要不要调用工具"` 的 `reasoning_delta`。
   - writer events 中不存在内容为 `"我先查一下"` 的 `assistant_delta`。
   - writer events 中存在最终 round 的 `reasoning_delta` 与 `assistant_delta`。
   - 最终 message 的 `Content` 和 `ReasoningContent` 正确。

### 6.2 调整既有流式测试预期

现有测试中直接最终回答的场景仍应通过：

- `TestRespondToConversation_StreamsAssistantContentDeltas`：`go-agent/internal/web/runtime/runtime_test.go:359`
- `TestRespondToConversation_StreamsReasoningContentDeltas`：`go-agent/internal/web/runtime/runtime_test.go:425`

现有工具调用流式测试需要确认事件顺序仍合理：

- `TestRespondToConversation_StreamsToolCallsAcrossChunks`：`go-agent/internal/web/runtime/runtime_test.go:459`

该测试当前已经期望工具事件后跟最终答案 delta，推荐保留该语义；如果新增第一轮工具调用 content/reasoning，需要明确断言这些中间 delta 不出现。

### 6.3 Helper 单元测试

可给 `shouldEmitModelDeltas` 增加轻量表驱动测试：

```go
func TestShouldEmitModelDeltas(t *testing.T)
```

覆盖：

| finishReason | toolCalls | expected |
| --- | --- | --- |
| `stop` | nil | true |
| `tool_calls` | 1 个 | false |
| `stop` | 1 个 | false |
| `""` | nil | true |

## 7. 风险与兼容性

### 7.1 最终答案首字延迟

由于需要等 round 结束后才能确认是否调用工具，最终 round 的 delta 会在模型完整返回后再回放给前端。用户仍会收到 delta 事件，但不再是严格实时逐 token。

这是满足“不把工具调用 round 的内容发给前端”的必要权衡。若未来模型 API 能在开始阶段可靠声明本 round 是否会调用工具，可再优化为真正实时。

### 7.2 内存占用

每个 round 会在内存里额外保存一份 delta 列表。当前本来已经用 `strings.Builder` 保存完整 `content` 和 `reasoning_content`，新增列表主要保存事件边界与字符串内容，影响可控。

如后续担心超长输出，可改为只缓存 `{event, start, end}` 或最终 round 不回放 delta、只发 assistant；本次不做过度设计。

### 7.3 前端兼容

前端事件协议不变：`web/src/api.ts:154` 仍按事件名分发，`web/src/App.tsx:305` / `web/src/App.tsx:308` 仍追加 delta。

因为后端不再发送工具 round 的 delta，前端无需判断 round 类型。

## 8. 实施步骤

1. 在 `go-agent/internal/web/runtime/conversation_flow.go` 增加 `bufferedModelDelta` 与事件名常量。
2. 将 `runModelRoundStream` 中即时发送 `assistant_delta` / `reasoning_delta` 的逻辑改为缓存。
3. 增加 `shouldEmitModelDeltas`，在 stream 结束后决定是否回放缓存 delta。
4. 保持 `RespondToConversation`、Stop hook、工具 hook 不变。
5. 在 `go-agent/internal/web/runtime/runtime_test.go` 新增工具调用 round 不下发 content/reasoning 的测试。
6. 运行验证：

```bash
cd /Users/bytedance/golang_pro/nano_cc/go-agent && go test ./internal/web/runtime
```

必要时再运行：

```bash
cd /Users/bytedance/golang_pro/nano_cc/go-agent && go test ./...
```

## 9. 验收标准

1. 当模型 round 最终以 `tool_calls` 结束时：
   - 前端不会收到该 round 的 `assistant_delta`。
   - 前端不会收到该 round 的 `reasoning_delta`。
2. 当模型最终 round 不再调用工具时：
   - 前端会收到该 round 的 `assistant_delta` / `reasoning_delta`。
   - 前端最终会收到 `assistant` 事件，且包含完整 `content` 与 `reasoning_content`。
3. 工具调用事件 `tool` 仍正常发送。
4. 会话历史中仍保留工具调用所需上下文，模型后续轮次能继续基于工具结果回答。
5. `go test ./internal/web/runtime` 通过。
