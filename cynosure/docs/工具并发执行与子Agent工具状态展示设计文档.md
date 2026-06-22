# 工具并发执行与子Agent工具状态展示设计文档

## 1. 背景与目标

当前 runtime 在同一轮 `tool_calls` 中按模型返回顺序串行执行工具。多个互不依赖的工具调用会被依次阻塞，整体耗时等于所有工具耗时之和。当前 TUI 已支持主 Agent 工具调用的 `tool_call_start` / `tool_call_done` 展示，但子 Agent 内部执行工具时不会把工具过程展示到父 TUI。用户只能看到 `spawn_subagent` 结束后的最终结果，无法知道子 Agent 当前在调用哪些工具、是否仍在运行。

本设计目标：

1. 同一轮对话中模型一次返回的多个工具调用并发执行，并在全部完成后再把工具结果按原始顺序送回模型。
2. 子 Agent 执行期间把内部工具调用状态展示到 TUI，但不展示内部工具结果内容。
3. 子 Agent 完成后移除这批内部工具状态展示，避免留下空白或噪音，然后按现有方式展示 `spawn_subagent` 工具的最终调用结果。
4. 保持审批、安全、工具结果入历史、上下文压缩、记忆更新、超时等现有策略不变。

## 2. 当前实现梳理

主 Agent 工具执行链路在 `internal/agent/runtime/conversation_flow.go`：

- 模型返回 `tool_calls` 后，先把 assistant tool-call 消息追加到 `state.History` 与 `state.ModelHistory`。
- 随后 `for _, tc := range msg.ToolCalls` 串行执行每个工具。
- 每个工具的顺序是 `RunPreToolUse` -> `approveToolCall` -> `emitToolCallStart` -> `executeToolCall` -> `RunPostToolUse` -> `emitToolCallDone` -> `ToolCallCount++` -> `emitMeta`。
- 工具结果由 `RunPostToolUse` 中的 hook 追加到 `state.Messages`、`state.History`、`state.ModelHistory`，并写入工具结果日志。

子 Agent 工具执行链路在 `internal/agent/runtime/subagent.go`：

- `runSubagent` 为子 Agent 构造独立的 `LoopState`，当前 `childState.Writer` 为空，因此子 Agent 内部事件不会进入 TUI。
- `runSubagentLoop` 与主循环类似，也是串行遍历 `msg.ToolCalls`。
- 子 Agent 的工具执行仍走相同 hook 与审批逻辑，所以工具结果会进入子 Agent 自己的 `state.ModelHistory`，用于后续子 Agent 模型轮次。
- `spawn_subagent` 完成后，父工具返回 `"Subagent completed.\n\nSummary:\n..."`，父 TUI 展示的是 `spawn_subagent` 这一个工具调用的结果摘要。

TUI 展示链路在 `internal/tui/app.go`：

- `tool_call_start` 会追加一条 `Role="tool"` 的消息。
- `tool_call_done` 会按 `tool_call_id` 从后向前更新同一条工具消息。
- `renderToolMessage` 默认展示工具名、参数摘要、状态和结果摘要。
- `shouldHideMessageAt` 会在后续 assistant 内容出现后隐藏较早的 tool/thinking 消息，用于减少历史噪音。

现有测试中 `TestRespondToConversation_SpawnSubagentUsesFreshMessagesStoresTraceAndDoesNotEmitToolEvents` 明确断言 `spawn_subagent` 不向前端发工具事件。本次需求需要更新这条旧约束。

## 3. 范围

### 3.1 纳入范围

1. 主 Agent 同一轮 `tool_calls` 并发执行。
2. 子 Agent 同一轮 `tool_calls` 使用同一套并发执行机制。
3. 工具结果仍按模型返回的 tool-call 顺序追加到 OpenAI 消息历史，保证 `assistant.tool_calls` 与后续 `tool` 消息合法配对。
4. 工具审批仍逐个工具执行，审批 UI 和允许规则不改变。
5. TUI 显示子 Agent 内部工具调用状态，包括工具名、参数摘要、运行中/完成/失败状态。
6. 子 Agent 内部工具状态不展示 `result_preview`。
7. 子 Agent 完成后隐藏或删除内部工具状态消息，并正常展示父级 `spawn_subagent` 工具的完成结果。

### 3.2 不纳入范围

1. 不新增用户配置开关。
2. 不改变工具 schema、工具权限模型、审批规则和允许列表。
3. 不改变 `spawn_subagent` 的模型可见结果格式。
4. 不把子 Agent 内部工具结果全文展示给 TUI。
5. 不做跨轮工具并发；只并发同一次模型响应中的 `msg.ToolCalls`。
6. 不允许子 Agent 再调用 `spawn_subagent`，保留现有深度限制。

## 4. 推荐方案

### 4.1 统一工具批执行器

新增一个 runtime 内部函数，用于执行一轮模型返回的工具批：

```go
type toolBatchOptions struct {
	ToolContext        ToolContext
	Registry           *ToolRegistry
	UseChildRegistry   bool
	SuppressResultUI   bool
	EphemeralGroupID   string
	OnSuccessTodoWrite func([]agenttools.TodoItem)
}

type toolBatchResult struct {
	Contexts []*ToolUseContext
	Rejected bool
}
```

函数职责：

1. 为每个 `openai.ToolCall` 创建 `ToolUseContext`。
2. 先按原始顺序执行 `RunPreToolUse` 与 `approveToolCall`。
3. 审批全部通过后，为每个工具发送 `tool_call_start` 事件。
4. 使用 goroutine 并发执行工具本体。
5. 等待所有 goroutine 完成。
6. 按原始顺序执行 `RunPostToolUse`，从而按原始顺序追加工具结果消息。
7. 按原始顺序发送 `tool_call_done`、更新 `ToolCallCount`、发送 `meta`。

这样可以并发耗时最大的工具执行阶段，同时保持历史顺序、hook 顺序、TUI 完成事件顺序稳定。

### 4.2 审批策略

审批仍然在工具真正执行前逐个处理。原因：

- 现有审批 UI 是阻塞式交互，一次只处理一个请求。
- 并发弹出多个审批会破坏 TUI 输入劫持与 `beginApproval` 语义。
- 审批不是主要耗时来源，保持串行审批风险最低。

如果任意工具被拒绝：

- 已通过审批但尚未启动的工具不执行。
- 被拒绝工具发送 start/done，状态为 `rejected`，与当前主 Agent 行为一致。
- 本轮直接结束，不进入后续并发执行，也不做记忆/历史收尾。
- 已经开始执行的工具不存在，因为本设计在全部审批完成后才启动并发执行。

### 4.3 并发执行与顺序一致性

并发只覆盖工具本体执行：

```text
PreToolUse/Approval:  顺序执行
Tool start event:     顺序发送
Tool execution:       并发执行
PostToolUse/history:  顺序执行
Tool done event/meta: 顺序发送
```

这样满足两个约束：

- 用户看到多个工具同时进入 running 状态。
- OpenAI 请求历史仍保持 `assistant(tool_calls)` 后紧跟相同顺序的 `tool` 结果消息，避免破坏上下文压缩中的 tool-call 边界修复逻辑。

并发执行时每个 goroutine 只写自己的 `ToolUseContext.Outcome`，不直接写 `state.History`、`state.ModelHistory`、`state.Messages`。共享状态写入集中在等待完成后的顺序阶段完成。

### 4.4 Todo 状态处理

`todo_write` 成功后仍需要更新 `state.Todos`。并发后可能同一批里出现多个 `todo_write`，处理规则为：

- 工具执行阶段不直接写 `state.Todos`。
- 顺序收尾阶段按模型返回顺序处理成功的 `todo_write`。
- 后一个成功的 `todo_write` 覆盖前一个快照，与串行执行的最终效果一致。

传入工具执行的 `ToolContext.Todos` 使用本轮开始时的快照，避免 goroutine 并发读写。该行为与“同一批工具由同一个模型响应生成，批内不依赖前一个工具结果”的工具调用语义一致。

### 4.5 子 Agent 工具状态事件

为子 Agent 的 `childState.Writer` 设置一个包装 writer，而不是让子 Agent 直接复用父 writer 的原始事件：

```go
type subagentEventWriter struct {
	parent       EventWriter
	runID        string
	parentCallID string
}
```

包装规则：

- 只转发 `tool_call_start`、`tool_call_done`、`meta` 中需要的工具状态信息。
- 给内部工具事件增加：
  - `scope: "subagent"`
  - `subagent_run_id`
  - `parent_tool_call_id`
  - `ephemeral_group_id`
  - `suppress_result: true`
- `tool_call_id` 使用 `subagent_run_id + ":" + 原始 tool_call_id`，避免与父 Agent 工具 ID 冲突。
- `tool_call_done` 不携带内部工具 `result_preview`，只携带状态与参数摘要。

子 Agent 仍使用自己的 `state.ModelHistory` 接收真实工具结果，供子 Agent 后续轮次推理。隐藏结果只影响 TUI。

### 4.6 TUI 覆盖与隐藏策略

新增一个轻量事件：

```text
tool_call_group_clear
```

事件字段：

```text
ephemeral_group_id string
```

TUI 的 `ToolCallView` 增加：

```go
Scope            string
EphemeralGroupID string
SuppressResult   bool
Hidden           bool
```

处理规则：

1. 子 Agent 内部工具 start/done 事件渲染为普通工具状态行，但 `SuppressResult=true` 时不渲染结果摘要。
2. `tool_call_group_clear` 到达后，TUI 从 `m.messages` 中删除该 `EphemeralGroupID` 对应的 tool 消息，而不是渲染成空字符串。
3. 删除后立即 `refreshViewport()`，避免原位置留下大量空行。
4. 父级 `spawn_subagent` 的工具 done 事件随后按现有逻辑更新同一条父工具消息，展示最终摘要。

删除消息比 `shouldHideMessageAt` 更适合本需求，因为用户明确要求“覆盖掉展示工具调用状态的信息，不然会出现很多空行”。`shouldHideMessageAt` 只能跳过渲染，但如果调用点仍无条件追加分隔换行，容易留下空白；直接从消息切片删除更确定。

### 4.7 事件时序

子 Agent 执行时推荐的 TUI 事件顺序：

```text
父 spawn_subagent tool_call_start
子 read_file tool_call_start(scope=subagent, group=G)
子 grep tool_call_start(scope=subagent, group=G)
子 read_file tool_call_done(scope=subagent, suppress_result=true)
子 grep tool_call_done(scope=subagent, suppress_result=true)
tool_call_group_clear(group=G)
父 spawn_subagent tool_call_done(result_preview=Subagent completed...)
meta
```

如果子 Agent 失败：

- 仍发送 `tool_call_group_clear(group=G)`。
- 父 `spawn_subagent` 工具 done 状态为 `rejected`，展示错误摘要。
- 不展示子 Agent 内部工具结果。

## 5. 错误与取消处理

1. 任一工具执行返回错误时，`executeToolCall` / `executeChildToolCall` 仍把错误包装成 `toolExecutionOutcome{Status:"rejected"}`，批处理器等待其他已启动工具结束后再顺序收尾。
2. 如果 `ctx` 被取消，工具 handler 仍按现有 context 行为停止；批处理器等待 goroutine 返回，并把取消错误向上返回。
3. 子 Agent 超时策略不变，仍使用 `subAgentTurnTimeout` 软边界；本次不引入 `context.WithTimeout`。
4. 批内并发不抢占已运行工具，不额外 kill 工具；底层工具已有的超时机制保持不变。

## 6. 文件改动计划

预计修改：

- `internal/agent/runtime/conversation_flow.go`
  - 提取主 Agent 工具批执行逻辑。
  - 主循环从串行 `for` 改为调用批执行器。
  - 扩展工具事件字段：`scope`、`ephemeral_group_id`、`suppress_result`。
- `internal/agent/runtime/subagent.go`
  - 子 Agent loop 使用同一批执行器。
  - 为子 Agent 设置包装 writer。
  - 子 Agent 完成或失败时发送 `tool_call_group_clear`。
- `internal/tui/app.go`
  - `ToolCallView` 增加临时分组与结果隐藏字段。
  - 处理 `tool_call_group_clear`。
  - `renderToolMessage` 在 `SuppressResult` 时只展示状态，不展示结果内容。
- `internal/tui/events_test.go`
  - 增加子 Agent 临时工具状态展示和清理测试。
- `internal/agent/runtime/runtime_test.go`
  - 增加同轮工具并发执行测试。
  - 更新子 Agent 工具事件测试，删除“不发送工具事件”的旧断言。

不预计修改：

- `internal/tools/*` 工具实现。
- 工具 schema 和允许列表。
- prompt 静态资产。
- README。此变更是运行机制与 TUI 行为调整，若实现阶段发现 README 已描述串行工具执行或子 Agent 不展示过程，再同步修正。

## 7. 测试计划

### 7.1 Runtime 测试

1. `TestRespondToConversation_ExecutesSameRoundToolsConcurrently`
   - 构造同一轮两个 fake 工具。
   - 第一个工具阻塞等待第二个工具开始。
   - 若串行执行会超时；并发执行应通过。
   - 断言最终写入模型历史的 tool 消息仍按模型返回顺序排列。

2. `TestRespondToConversation_ToolBatchEmitsStartBeforeWaitingForDone`
   - 同一轮两个工具。
   - 断言两个 `tool_call_start` 都出现在任意 `tool_call_done` 之前。

3. `TestRespondToConversation_RejectedApprovalDoesNotStartParallelBatch`
   - 模拟第二个工具审批拒绝。
   - 断言第一个已审批工具没有被执行。
   - 断言拒绝事件与现有行为一致。

4. `TestRespondToConversation_SubagentForwardsToolStatusWithoutResultAndClearsGroup`
   - `spawn_subagent` 内部调用工具。
   - 断言 writer 收到子工具 start/done，且带 `scope=subagent`、`suppress_result=true`。
   - 断言子工具 done 不包含内部 `result_preview`。
   - 断言父 `spawn_subagent` done 前收到 `tool_call_group_clear`。

### 7.2 TUI 测试

1. `TestModelDisplaysSubagentToolStatusWithoutResultPreview`
   - 发送子 Agent 工具 start/done。
   - 断言渲染包含工具名、参数摘要、状态。
   - 断言不包含 `result_preview`。

2. `TestModelClearsSubagentToolGroupWithoutBlankLines`
   - 追加多个同组子 Agent tool 消息。
   - 发送 `tool_call_group_clear`。
   - 断言 `m.messages` 中没有该组消息。
   - 断言渲染文本不包含连续多段空白工具占位。

3. `TestModelKeepsParentSpawnSubagentResultAfterClearingChildGroup`
   - 先显示父 `spawn_subagent` running。
   - 显示并清理子工具组。
   - 发送父 `spawn_subagent` done。
   - 断言最终只保留父工具结果展示。

### 7.3 回归验证

实现后运行：

```bash
go test ./internal/agent/runtime ./internal/tui
go test ./...
```

如果全量测试出现已知 TUI header/ASCII-art 基线失败，需要先确认是否与本次改动相关，再记录结果。

## 8. 风险与约束

1. 工具并发可能暴露工具 handler 的共享全局状态问题。实现阶段应避免在 goroutine 内写 runtime 共享状态，必要时仅把具体 handler 调用并发化。
2. 同一批工具理论上不应依赖彼此结果，因为模型是在看到任何结果前一次性发出的多个 tool calls。按此语义并发执行合理。
3. `load_skill` 与其他工具同批出现时，其他工具不会看到本批 `load_skill` 新增的 skill snapshot；这与 OpenAI 多工具调用语义一致。如果模型需要使用刚加载的 skill，应在下一轮工具调用中进行。
4. 子 Agent 内部工具状态是临时 UI，不进入父 Agent 模型历史，也不进入会话持久化。
5. 审批串行处理会让“需要审批的工具批”不能完全同时开始，但工具执行阶段仍并发。该取舍保证现有交互安全。

## 9. 自检结论

- 本设计没有引入新的权限策略或配置项。
- 本设计保留工具结果按原始顺序进入模型历史。
- 本设计明确子 Agent 内部工具结果只进入子 Agent 上下文，不展示到 TUI。
- 本设计通过删除临时消息解决子 Agent 工具状态完成后的空行问题。
- 本设计只要求实现阶段修改 runtime/TUI/测试，不要求同步修改 prompt。
