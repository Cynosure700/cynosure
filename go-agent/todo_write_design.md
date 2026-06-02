# todo_write 工具设计文档

## 1. 背景

当前 go-agent 已经支持 `bash`、`read_file`、`write_file`、`edit_file`、`load_skill`、`spawn_subagent` 等工具。工具定义集中在 `internal/tools/definitions.go`，Web 运行时会基于 `WebAllowedTools` 构建实际暴露给模型的工具列表，并在每一轮模型请求中通过 `ChatCompletionRequest.Tools` 传给模型。

本次需求是新增一个 `todo_write` 工具：

1. 它不提供文件、Shell、网络、子 Agent 等执行能力。
2. 它只让模型显式维护“当前任务计划”，提升多步骤任务的规划与进度跟踪能力。
3. 如果模型连续 3 个模型轮次没有调用 `todo_write`，运行时注入一条提醒，促使模型在合适场景更新计划。

相关现有位置：

- 工具定义：`go-agent/internal/tools/definitions.go:42`
- 工具 Handler 注册：`go-agent/internal/tools/handlers.go:8`
- Web 工具 allowlist 与过滤：`go-agent/internal/web/runtime/tool_catalog.go:13`
- 每轮模型请求与工具调用循环：`go-agent/internal/web/runtime/conversation_flow.go:47`
- 工具调用执行入口：`go-agent/internal/web/runtime/tool_execution.go:14`
- Hook 中追加 tool message：`go-agent/internal/web/runtime/hooks/tool.go:71`
- 子 Agent 工具注册：`go-agent/internal/web/runtime/subagent.go:56`
- 系统提示构建中的工具名提示：`go-agent/internal/web/runtime/prompt_builder.go:16`

## 2. 目标

本次改造目标：

1. 新增 OpenAI function tool：`todo_write`。
2. 工具参数 schema 与需求保持一致：
   - 入参只包含必填字段 `todos`。
   - 每个 todo 包含必填字段 `id`、`content`、`status`。
   - `status` 枚举为 `pending`、`in_progress`、`completed`。
3. 工具执行不产生外部副作用：
   - 不读写工作区文件。
   - 不执行命令。
   - 不调用外部服务。
   - 不新增数据库迁移。
4. 在运行时内存中保存当前任务计划，方便后续 hook、日志或 UI 展示复用。
5. 在主 Agent 和子 Agent 中都可用，除非配置显式不暴露该工具。
6. 连续 3 个模型轮次没有调用 `todo_write` 时，在下一次请求模型前注入一条临时提醒消息。
7. 保持已有工具调用、SSE 展示、历史持久化语义基本不变。

## 3. 非目标

本次不做以下事情：

1. 不新增前端 TODO 面板。
2. 不新增单独的 todo 数据表或迁移。
3. 不跨会话持久化结构化 todo 状态；todo 状态仅在本次运行时循环内存中保存。
4. 不强制所有任务必须调用 `todo_write`。
5. 不改变模型是否能调用其他执行型工具的权限模型。
6. 不把 `todo_write` 设计成任务调度器、执行器或审批器。

## 4. 当前架构分析

### 4.1 工具定义与暴露

当前基础工具定义在 `baseToolDefs` 中，`AllToolDefs` 在基础工具之外追加 `spawn_subagent`，`ChildToolDefs` 当前等于 `baseToolDefs`。对应代码在 `go-agent/internal/tools/definitions.go:42`。

Web 运行时通过 `loadAllowedToolNames` 读取配置中的 `WebAllowedTools`，并调用 `lookupRegisteredTool` 从 `agenttools.AllToolDefs` 中过滤合法工具。也就是说，新工具若只加入定义但未加入默认 allowlist，则默认不会暴露给模型。

### 4.2 工具执行与消息循环

主 Agent 的模型循环在 `RespondToConversation` 中完成。每轮请求模型后：

1. 如果模型没有返回工具调用，则进入 stop hook 并结束。
2. 如果返回工具调用，则逐个执行工具。
3. 工具结果以 `tool` role message 追加到 `state.Messages` 和 `state.History`。
4. 循环进入下一轮模型请求。

相关逻辑在 `go-agent/internal/web/runtime/conversation_flow.go:47` 和 `go-agent/internal/web/runtime/hooks/tool.go:71`。

这意味着 `todo_write` 的调用和结果天然会进入模型上下文，模型后续轮次可以看到自己刚更新的计划。只要工具结果足够简洁，就能避免明显增加上下文噪音。

### 4.3 子 Agent

子 Agent 使用 `NewChildToolRegistry`，该函数会移除 `spawn_subagent`，避免递归创建子 Agent。由于 `todo_write` 不增加执行能力，可以允许子 Agent 使用它进行独立任务规划。相关位置是 `go-agent/internal/web/runtime/subagent.go:56`。

## 5. 方案对比

### 方案 A：只在系统提示中要求模型维护 TODO

在系统提示中新增一段说明，要求模型多步骤任务时维护 TODO，但不提供工具。

优点：

- 改动最少。
- 没有工具调用与测试成本。

缺点：

- 模型没有结构化动作，难以检测“连续 3 轮没更新计划”。
- 无法通过工具调用历史明确知道模型是否更新过计划。
- 不满足用户明确要求的 `todo_write` function schema。

结论：不采用。

### 方案 B：新增 `todo_write` 工具，并在运行时内存保存当前 TODO【推荐】

新增一个无副作用工具。Handler 校验参数并产出结构化 todo 列表，运行时将其写入当前 `LoopState` 的内存字段，再返回简洁确认文本，例如：

```text
Todo list updated: 3 items (1 in_progress, 2 pending, 0 completed).
```

计划内容会同时存在于两处：

1. 模型发起的 tool call arguments 中，随 assistant 中间消息进入 conversation history，便于历史回放。
2. 当前运行时 `LoopState` 的内存字段中，便于本次循环内后续 hook、SSE 展示或调试日志直接读取。

工具结果只负责确认更新成功与提供摘要，不回显完整 todo 内容。

优点：

- 完全符合“不给 Agent 增加执行能力”的约束。
- 不需要新增数据库结构。
- 与现有工具执行框架兼容。
- 模型可通过历史 tool call arguments 回看最新计划。
- 后端运行时也能直接读取当前最新计划，不必反解析 tool call arguments。
- 易于测试。

缺点：

- 内存状态只在当前响应流程内有效，服务重启或重新打开历史会话后需要从 history 派生。
- 如果未来要做跨会话 TODO 面板，仍需要从 history 派生或新增存储。

结论：本次采用。

### 方案 C：新增 `todo_write` 工具并持久化最新 TODO 状态

新增数据库字段或表保存每个会话的最新 TODO。

优点：

- 前端可直接展示最新计划。
- 后续可做任务状态统计。

缺点：

- 超出当前“只增加规划能力”的范围。
- 需要迁移和更多 API 设计。
- 会引入并发、历史回放、跨轮覆盖语义等额外问题。

结论：不作为本次方案。

## 6. 推荐设计

### 6.1 工具定义

在 `internal/tools/definitions.go` 新增 `todo_write` 定义。建议将其放入基础工具列表，使主 Agent 和子 Agent 都能使用。

工具 schema：

```json
{
  "name": "todo_write",
  "description": "Create or update the current task plan. Use this tool to track progress on multi-step tasks.",
  "parameters": {
    "type": "object",
    "properties": {
      "todos": {
        "type": "array",
        "items": {
          "type": "object",
          "properties": {
            "id": { "type": "string" },
            "content": { "type": "string" },
            "status": {
              "type": "string",
              "enum": ["pending", "in_progress", "completed"]
            }
          },
          "required": ["id", "content", "status"]
        }
      }
    },
    "required": ["todos"]
  }
}
```

为了让模型更容易正确使用，可在 description 中保持克制，不加入过多策略；具体“何时使用”通过提醒消息完成。

### 6.2 默认暴露策略

建议把 `todo_write` 加入默认 `WEB_ALLOWED_TOOLS`：

```text
load_skill,bash,read_file,write_file,edit_file,todo_write
```

原因：

1. 它不增加执行能力，默认开启风险低。
2. 如果不默认开启，“连续 3 轮提醒”在默认配置下无法发挥作用。
3. 仍然保留 `WebAllowedTools` 配置能力，部署方可以显式移除。

同时需要更新配置测试中对默认工具列表的断言。

### 6.3 Handler 行为

在 `internal/tools/handlers.go` 注册 `todo_write` handler，并新增独立实现文件，例如 `internal/tools/todo_write.go`。

Handler 职责：

1. 从 `args["todos"]` 解析数组。
2. 校验每个 todo：
   - `id` 非空。
   - `content` 非空。
   - `status` 必须是 `pending`、`in_progress`、`completed` 之一。
3. 不限制 todo 数量，但为了避免上下文膨胀，可在工具结果中只返回统计摘要，不回显完整 todo 内容。
4. 返回确认文本。

同时建议定义复用的结构体，供 handler、runtime 内存字段和测试共享：

```go
type TodoItem struct {
    ID      string `json:"id"`
    Content string `json:"content"`
    Status  string `json:"status"`
}
```

建议结果格式：

```text
Todo list updated: 4 items (pending: 2, in_progress: 1, completed: 1).
```

错误示例：

```text
Error: todos[0].status must be one of pending, in_progress, completed
```

错误返回沿用现有工具执行模型：handler 返回 error，`executeToolCall` 会把工具调用标记为 `rejected`，并将错误内容包装进 tool result。

### 6.4 内存状态字段

在 `go-agent/internal/web/runtime/hooks/types.go:37` 的 `LoopState` 中新增字段，用于保存当前运行时最新 todo 状态：

```go
type LoopState struct {
    // existing fields...
    Todos []agenttools.TodoItem
}
```

字段语义：

1. `Todos` 表示当前响应流程内最近一次成功 `todo_write` 提交的完整任务列表。
2. 每次 `todo_write` 成功执行后，用新列表整体覆盖旧列表，而不是增量 patch。
3. 如果模型传入空数组，则 `Todos` 被更新为空数组，表示清空当前计划。
4. 如果 `todo_write` 参数校验失败，则不更新 `Todos`。
5. `Todos` 不直接持久化到数据库；历史仍由 assistant tool call arguments 和 tool message 保存。

为了让工具 handler 能把解析后的结构化结果传回 runtime，建议将 `ToolExecutionResult` 扩展一个可选字段：

```go
type ToolExecutionResult struct {
    Output string
    Todos  []agenttools.TodoItem
}
```

同时将 runtime hook 使用的 `ToolExecutionOutcome` 扩展一个不进入 tool message JSON 的内存字段：

```go
type ToolExecutionOutcome struct {
    Status string             `json:"status"`
    Result string             `json:"result"`
    Audit  ToolExecutionAudit `json:"-"`
    Todos  []agenttools.TodoItem `json:"-"`
}
```

只有 `todo_write` 会填充 `Todos`；其他工具保持零值。`ToolRegistry.Execute` 返回后，`executeToolCall` 将 `execResult.Todos` 拷贝到 `toolExecutionOutcome.Todos`。工具调用循环在 `toolCtx.Outcome` 返回后根据工具名和状态判断：当 `name == "todo_write"` 且执行成功时，将 `toolCtx.Outcome.Todos` 整体写入 `state.Todos`。

这样设计有两个好处：

1. `todo_write` handler 仍然只做参数解析与校验，不直接依赖 Web runtime 的 `LoopState`。
2. runtime 可以在统一入口控制内存状态更新，避免把运行时状态写入底层 tools 包。

### 6.5 “连续 3 轮没调 todo_write”提醒机制

#### 6.5.1 轮次定义

“一轮”定义为一次模型请求与响应，即 `RespondToConversation` 循环中一次 `runModelRoundStream` 调用。该循环当前在 `go-agent/internal/web/runtime/conversation_flow.go:47`。

#### 6.5.2 计数规则

在单次 `RespondToConversation` 调用内维护计数器：

- `roundsSinceTodoWrite`：当前响应流程中，距离上次 `todo_write` 工具调用已经过去的模型轮次数。
- 初始值为 0。
- 每轮模型返回后：
  - 如果 `msg.ToolCalls` 中包含 `todo_write`，计数器重置为 0。
  - 否则计数器加 1。
- 如果计数器达到 3，则在下一次请求模型前注入提醒，并将计数器重置为 0，避免每轮重复注入。

#### 6.5.3 注入时机

只在“还会继续请求下一轮模型”时注入提醒。

具体来说：

1. 如果当前轮没有工具调用，模型即将结束，没必要注入提醒。
2. 如果当前轮有工具调用并执行完成，运行时会进入下一轮模型请求；此时如果计数达到 3，就在下一轮请求前追加提醒消息。

这样可以避免在最终回答已经产生时额外污染会话。

#### 6.5.4 注入内容

建议使用临时的 `system` role message 追加到 `state.Messages`，不写入 `state.History`：

```text
<system-reminder>
You have not called todo_write for 3 consecutive model rounds. If the task is multi-step or your plan has changed, call todo_write to create or update the current task plan before continuing. If todo_write is unnecessary for this simple step, continue normally.
</system-reminder>
```

不写入 `state.History` 的原因：

- 这是运行时策略提醒，不是用户消息或模型真实输出。
- 避免重新打开会话后重复积累提醒。
- 避免影响长期上下文质量。

如果后续发现目标模型不接受中途追加 `system` role，可退化为 `user` role 的同样内容；本次设计优先使用 `system`，因为它最准确表达“运行时提醒”。

#### 6.5.5 工具未启用时的行为

如果当前 `ToolRegistry` 中没有暴露 `todo_write`，不启用提醒机制。否则会提醒模型调用一个不可用工具，造成无效行为。

需要在 runtime 中增加一个小方法，例如：

```go
func (r *ToolRegistry) hasTool(name string) bool
```

或者复用现有 `isAllowed`。

### 6.6 与子 Agent 的关系

子 Agent 使用独立消息列表与独立循环。建议复用同一套提醒逻辑：

1. 子 Agent 可调用 `todo_write`。
2. 子 Agent 的“连续 3 轮”计数独立于父 Agent。
3. 子 Agent 的 `LoopState.Todos` 独立于父 Agent，不回写父 Agent 的 `LoopState.Todos`。
4. 子 Agent 的提醒同样只写入子 Agent 的临时 `state.Messages`，不写入主 conversation history，也不写入 subagent trace。

原因：子 Agent 往往处理隔离任务，也可能需要规划；但提醒属于运行时策略，不应作为实际执行轨迹持久化。

### 6.7 SSE 与工具展示

现有 `emitToolEventHook` 会对非 `spawn_subagent` 工具发出 `tool` SSE 事件。`todo_write` 默认也可以展示为普通工具调用：

- name：`todo_write`
- status：`success` / `rejected`
- result：摘要文本

这符合“规划能力是显式动作”的产品语义。由于 result 不回显完整 todo 列表，不会显著增加 UI 噪音。

如未来前端要做专门 TODO 面板，本轮 SSE 可直接从 `LoopState.Todos` 或 `todo_write` 的结构化结果读取；重新打开历史会话时仍需要从 tool call arguments 中派生最新 todos，或另行设计持久化存储。

## 7. 数据流

### 7.1 正常调用 todo_write

```text
模型返回 tool_calls: todo_write(args.todos)
  -> runtime 执行 handleTodoWrite
  -> handler 校验 todos，返回摘要与结构化 TodoItem 列表
  -> runtime 将 TodoItem 列表整体覆盖到 state.Todos
  -> hook 记录 tool call / 发送 SSE / 追加 tool message
  -> 下一轮模型可从上下文看到 todo_write 调用和确认结果
```

### 7.2 三轮未调用后的提醒

```text
round 1: 模型调用 bash，未调用 todo_write，计数 = 1
round 2: 模型调用 read_file，未调用 todo_write，计数 = 2
round 3: 模型调用 edit_file，未调用 todo_write，计数 = 3
工具执行完成后，发现还要进入下一轮
  -> state.Messages 追加临时 system reminder
  -> 计数重置为 0
round 4: 请求模型时带上 reminder
```

如果第 4 轮模型调用 `todo_write`，计数继续保持 0；如果没有调用且还继续工具循环，则重新开始计数。

## 8. 实现范围

### 8.1 新增/修改文件

预计修改：

1. `go-agent/internal/tools/definitions.go`
   - 新增 `todo_write` 工具定义。
   - 将其加入基础工具集合。

2. `go-agent/internal/tools/handlers.go`
   - 注册 `todo_write` handler。

3. `go-agent/internal/tools/todo_write.go`
   - 新增 `TodoItem` 结构体、handler 与参数校验逻辑。

4. `go-agent/internal/web/runtime/tool_registry.go`
   - 扩展 `ToolExecutionResult`，增加可选 `Todos []agenttools.TodoItem` 字段。

5. `go-agent/internal/web/runtime/hooks/types.go`
   - 在 `LoopState` 中新增 `Todos []agenttools.TodoItem` 字段，用于保存当前运行时 todo 状态。
   - 在 `ToolExecutionOutcome` 中新增 `Todos []agenttools.TodoItem \`json:"-"\``，用于在 runtime 内部传递结构化 todo 状态。

6. `go-agent/internal/web/runtime/tool_execution.go`
   - `todo_write` 成功执行后，将 `ToolExecutionResult.Todos` 拷贝到 `ToolExecutionOutcome.Todos`。

7. `go-agent/internal/web/runtime/conversation_flow.go`
   - 在主 Agent 工具调用循环中，当 `todo_write` 成功后将 `toolCtx.Outcome.Todos` 写入当前 `LoopState.Todos`。

8. `go-agent/internal/config/web_config.go`
   - 默认 `WEB_ALLOWED_TOOLS` 增加 `todo_write`。

9. `go-agent/internal/web/runtime/conversation_flow.go`
   - 在主 Agent 循环中加入三轮提醒计数与临时消息注入。

10. `go-agent/internal/web/runtime/subagent.go`
   - 在子 Agent 循环中加入同样的提醒计数与临时消息注入。
   - 在子 Agent 工具调用循环中，当 `todo_write` 成功后写入子 Agent 自己的 `LoopState.Todos`。

11. `go-agent/internal/web/runtime/runtime_test.go`
   - 覆盖工具暴露、内存状态更新、提醒注入、计数重置等 runtime 行为。

12. `go-agent/internal/tools/runtime_test.go` 或新增 `todo_write_test.go`
   - 覆盖 handler 参数校验与摘要输出。

13. `go-agent/internal/config/config_test.go`
   - 更新默认工具列表断言。

### 8.2 建议抽取的小函数

为避免主 Agent 和子 Agent 复制逻辑，可新增 runtime 内部小工具：

```go
const todoWriteToolName = "todo_write"
const todoWriteReminderThreshold = 3

func toolCallsInclude(toolCalls []openai.ToolCall, name string) bool
func appendTodoWriteReminderIfNeeded(state *LoopState, tools *ToolRegistry, roundsSinceTodoWrite int) int
func todoWriteReminderMessage() openai.ChatCompletionMessage
```

也可以将计数器封装成结构体：

```go
type todoWriteReminderTracker struct {
    roundsSinceTodoWrite int
}
```

推荐先使用简单函数，避免过度抽象。

## 9. 边界与错误处理

1. `todos` 不是数组：返回参数错误。
2. todo item 不是对象：返回参数错误。
3. `id` 为空：返回参数错误。
4. `content` 为空：返回参数错误。
5. `status` 不在枚举内：返回参数错误。
6. `todos` 为空数组：允许。表示清空/重置当前计划，返回 `Todo list updated: 0 items ...`。
7. 同一个 `id` 重复：本次不拒绝，避免把工具变成复杂状态管理器；可在摘要中不特别处理。
8. 多个 `in_progress`：本次不拒绝，因为用户给定 schema 未要求限制；是否只允许一个进行中任务交给模型策略处理。
9. 参数校验失败：不更新 `LoopState.Todos`，保留上一次成功提交的内存状态。
10. 工具未启用：不注入提醒。
11. 当前轮已经是最终回答：不注入提醒。

## 10. 测试计划

### 10.1 单元测试：工具定义与 handler

1. `todo_write` 出现在注册工具列表中。
2. schema 中 `todos` 为 required。
3. `status` 枚举包含且只包含 `pending`、`in_progress`、`completed`。
4. 合法 todos 返回统计摘要。
5. 合法 todos 返回结构化 `TodoItem` 列表，供 runtime 写入内存状态。
6. 非法 status 返回错误。
7. 缺失 id/content/status 返回错误。
8. 空数组合法。

### 10.2 Runtime 测试：主 Agent 提醒

1. 连续 3 轮工具调用均不是 `todo_write`，第 4 轮请求前包含 reminder。
2. 中间调用 `todo_write` 会重置计数，不注入 reminder。
3. 如果工具列表没有 `todo_write`，即使连续 3 轮也不注入 reminder。
4. 如果第 3 轮之后模型结束而不是继续工具循环，不注入 reminder。
5. `todo_write` 成功后，`LoopState.Todos` 被更新为最新列表。
6. `todo_write` 失败后，`LoopState.Todos` 不被覆盖。

### 10.3 Runtime 测试：子 Agent

1. 子 Agent 工具列表包含 `todo_write`，不包含 `spawn_subagent`。
2. 子 Agent 成功调用 `todo_write` 后，写入子 Agent 自己的 `LoopState.Todos`。
3. 子 Agent 连续 3 轮未调用 `todo_write` 时，在下一轮临时注入 reminder。

### 10.4 配置测试

1. 默认 Web 工具列表包含 `todo_write`。
2. 显式 `WEB_ALLOWED_TOOLS` 仍然可以覆盖默认列表。

## 11. 风险与缓解

### 11.1 提醒过度干扰模型

风险：简单任务或短链路任务中提醒可能不必要。

缓解：只在工具循环持续超过 3 轮且仍会继续请求模型时注入；最终回答前不注入。

### 11.2 上下文膨胀

风险：模型频繁调用 `todo_write`，tool call arguments 会进入历史，内存中也会保存完整列表。

缓解：handler 结果不回显完整 todos；系统提示与提醒文案只建议多步骤任务使用，不强制每轮使用。内存状态只保存最新一次完整列表，不累计历史版本。

### 11.3 中途 system message 兼容性

风险：部分 OpenAI 兼容服务可能不接受非首位 system message。

缓解：如果测试或联调发现不兼容，将 reminder role 调整为 `user`，内容保持 `<system-reminder>` 包裹；设计中该点可作为实现时的兼容开关。

### 11.4 与工具 allowlist 不一致

风险：提醒模型调用未启用工具。

缓解：提醒注入前检查 `ToolRegistry` 是否允许 `todo_write`。

## 12. 验收标准

1. 默认 Web 工具列表中包含 `todo_write`。
2. 模型可以按指定 schema 调用 `todo_write`。
3. `todo_write` 不执行任何外部操作，只返回计划更新摘要。
4. `todo_write` 成功后，当前运行时内存字段保存最新 todo 状态。
5. `todo_write` 参数校验失败时，不覆盖已有内存 todo 状态。
6. 连续 3 个模型轮次未调用 `todo_write` 且仍需进入下一轮时，会注入提醒。
7. 调用 `todo_write` 后提醒计数重置。
8. 未启用 `todo_write` 时不会注入提醒。
9. 主 Agent 与子 Agent 行为一致，且子 Agent 仍不能调用 `spawn_subagent`。
10. 相关 Go 单元测试通过。

## 13. 建议实施顺序

1. 先新增 `todo_write` 工具定义、`TodoItem` 结构体、handler 和工具层测试。
2. 扩展 runtime `ToolExecutionResult` 与 `LoopState.Todos`，实现成功调用后的内存状态更新。
3. 更新默认 `WEB_ALLOWED_TOOLS` 与配置测试。
4. 在主 Agent runtime 中实现提醒计数与注入，并补测试。
5. 将相同逻辑应用到子 Agent loop，并补测试。
6. 运行 `go test ./...` 验证。
