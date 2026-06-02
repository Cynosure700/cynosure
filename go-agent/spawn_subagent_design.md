# spawn_subagent 工具设计文档

## 1. 背景

当前 Web Runtime 的主 Agent 循环位于 `internal/web/runtime/conversation_flow.go`：会话入口会加载历史消息、构造 `LoopState`，再以 `state.Messages` 调用模型；当模型返回 tool calls 时，逐个执行工具并把工具结果追加回本轮 messages，直到模型停止。

现有工具体系由以下部分组成：

- 工具定义：`internal/tools/definitions.go`
- 工具 handler 注册：`internal/tools/handlers.go`
- Web Runtime 工具白名单与执行入口：`internal/web/runtime/tool_catalog.go`、`internal/web/runtime/tool_registry.go`
- 工具运行时环境：`internal/tools/runtime_env.go`
- workspace 路径约束：`internal/tools/path_guard.go`
- 工具 hook：`internal/web/runtime/hooks/tool.go`、`internal/web/runtime/hooks/manager.go`

本次目标是新增一个 `spawn_subagent` 工具：主 Agent 调用它时，Runtime 启动一个子 Agent。子 Agent 使用全新的 `messages[]` 和独立循环执行任务，结束后只把摘要文本作为工具结果返回给主 Agent；子 Agent 的对话上下文不进入主 Agent 上下文，但文件系统副作用保留在同一个 workspace 中。

## 2. 目标

1. 新增工具名：`spawn_subagent`。
2. 主 Agent 调用 `spawn_subagent` 后，Runtime 在服务端同步启动一个子 Agent 循环。
3. 子 Agent 使用全新的 `messages[]`：不继承主对话 history，也不把主 Agent 当前 messages 透传给子 Agent。
4. 子 Agent 能调用常规工具完成任务；它对文件系统、命令执行产生的副作用保留在 workspace 中。
5. 子 Agent 结束后，只把最终摘要文本返回给主 Agent，作为 `spawn_subagent` 的工具结果。
6. 子 Agent 的工作目录只能位于 `WorkspaceRoot` 下。
7. 子 Agent 不能递归调用 `spawn_subagent`。
8. 子 Agent 的工具调用仍经过现有 hook 流程，尤其是工具使用前/后的审计、权限与记录逻辑。

## 3. 非目标

1. 不新增前端交互面板或专门的子 Agent UI。
2. 不持久化子 Agent 的完整消息列表为主会话 history。
3. 不支持子 Agent 并发后台运行；`spawn_subagent` 是同步工具调用，结束后返回文本。
4. 不设计多级 Agent DAG；明确禁止递归 spawn。
5. 不改变现有 `bash/read_file/write_file/edit_file/load_skill` 的参数语义。

## 4. 现状分析

### 4.1 主 Agent 循环

`RespondToConversation` 当前负责完整对话循环：

1. 加载 conversation history。
2. 运行 `UserPromptSubmit` hooks，把用户输入追加进 state history。
3. 构造 skill snapshot 和 system prompt。
4. 用 `buildOpenAIMessages(systemPrompt, history)` 生成 OpenAI messages。
5. 循环调用模型；如果返回 tool calls，则执行工具并通过 hooks 写回 tool message。
6. 如果模型停止，则通过 Stop hook 持久化助手回复。

这一套逻辑适合主会话，但不能直接复用于子 Agent，因为子 Agent 不应复用主 conversation history，也不应把自己的完整消息列表持久化回主会话。

### 4.2 工具执行与 hook

工具执行入口是 `Service.executeToolCall`，实际调用 `ToolRegistry.Execute`。工具执行前后在主循环中显式运行：

- `RunPreToolUse`
- `executeToolCall`
- `RunPostToolUse`

默认 hook 当前会：

- 记录工具审计信息。
- 持久化 `tool_calls` 记录。
- 向 SSE writer 发工具事件。
- 把 tool result 追加进当前 `LoopState.Messages` / `LoopState.History`。

子 Agent 要“仍经过权限 hook”，因此不能绕开这条执行链。设计上应把“执行一轮模型 + 工具 hooks”的循环抽成可复用内部方法，让主 Agent 与子 Agent 共用同一工具执行路径。

### 4.3 workspace 约束

现有 `RuntimeEnv` 包含：

- `WorkspaceRoot`
- `CurrentWorkingDir`
- `AllowOutsideWorkspace`
- `AllowDangerousCommands`

路径 guard 已在工具层验证读写路径和 bash path token 是否逃逸 workspace。`ToolRegistry.runtimeEnv()` 目前默认把 `CurrentWorkingDir` 设为 `WorkspaceRoot`。`spawn_subagent` 需要允许传入子 Agent 工作目录，但必须在进入子 Agent 循环前解析并验证该目录位于 workspace 下。

## 5. 方案对比

### 方案 A：把 `spawn_subagent` 做成普通 `internal/tools` handler

做法：在 `internal/tools/handlers.go` 新增 `handleSpawnSubagent`，handler 内直接调用模型和工具。

优点：符合现有工具 handler 注册形式，工具名可直接进入 `Handlers` map。

缺点：`internal/tools` 层目前不依赖 Web Runtime；如果 handler 内调用 `Service`，会形成包依赖倒置或循环依赖。并且 handler 很难复用现有 hook manager 与 conversation 上下文。

结论：不推荐。

### 方案 B：把 `spawn_subagent` 作为 Runtime 内建工具拦截执行【推荐】

做法：工具定义仍注册在工具目录中，方便白名单和模型工具定义统一；但执行时由 Web Runtime 在 `executeToolCall` 或 `ToolRegistry.Execute` 之前识别 `spawn_subagent`，调用 `Service.runSubagent(...)`。普通工具仍走 `ToolRegistry.Execute`。

优点：

- 能直接访问 `Service`、`Cfg`、`Tools`、`Hooks`、skill snapshot、conversation/user 上下文。
- 子 Agent 循环可以复用主 Runtime 的模型调用和 hook 链路。
- 避免 `internal/tools` 反向依赖 `internal/web/runtime`。
- 更容易保证子 Agent 不能递归 spawn：构造 child tool definitions 时排除 `spawn_subagent`。

缺点：工具执行路径里会出现一个 Runtime 特殊分支，需要在代码中明确注释边界。

结论：推荐采用。

### 方案 C：创建隐藏 conversation 作为子 Agent 会话

做法：每次 spawn 创建一个隐藏 conversation，把子 Agent 的完整消息历史持久化为独立会话。

优点：调试时可以回看子 Agent 全量过程。

缺点：违背“对话上下文被丢弃，只把摘要文本回传”的要求；还会引入隐藏会话生命周期和清理问题。

结论：不采用。

## 6. 推荐设计

采用方案 B：**`spawn_subagent` 是 Web Runtime 内建的特殊工具，子 Agent 循环复用主 Runtime 的模型调用、工具执行与 hook 机制，但使用独立 messages、独立 tool registry 和受限 runtime env。**

## 7. 工具接口设计

### 7.1 工具定义

新增工具：

```json
{
  "name": "spawn_subagent",
  "description": "Spawn a child agent with a fresh message list to complete an isolated task. The child agent may use workspace tools, but it cannot spawn another subagent. Only its final summary is returned.",
  "parameters": {
    "type": "object",
    "properties": {
      "task": {
        "type": "string",
        "description": "The task for the child agent to complete. Include all context it needs because parent conversation history is not shared."
      },
      "cwd": {
        "type": "string",
        "description": "Optional working directory for the child agent. Relative paths are resolved under the workspace root; absolute paths must remain inside the workspace."
      }
    },
    "required": ["task"]
  }
}
```

### 7.2 参数语义

- `task`：必填。主 Agent 必须把子 Agent 需要的背景写进这个字段，因为子 Agent 不会继承主 messages。
- `cwd`：可选。为空时使用 `WorkspaceRoot`；非空时解析到绝对路径，并验证目录存在且位于 workspace 内。

不建议在第一版加入 `max_rounds`、`model`、`tools` 等参数，避免让模型自行突破运行边界。轮数和工具集合由 Runtime 固定控制。

### 7.3 返回语义

`spawn_subagent` 的工具结果应是结构化文本，至少包含：

```text
Subagent completed.

Summary:
<子 Agent 最终回答文本>
```

如果子 Agent 失败：

```text
Subagent failed: <错误原因>
```

工具调用对主 Agent 来说仍是普通 tool result；主 Agent 可以基于该摘要继续回答用户。

## 8. 子 Agent 循环设计

### 8.1 抽取通用循环

新增内部方法，供主会话和子 Agent 复用：

```go
type agentRunOptions struct {
    Label        string
    Conversation storage.Conversation
    User         storage.User
    State        *LoopState
    Tools        *ToolRegistry
    ToolContext  ToolContext
    PersistStop  bool
    EmitStream   bool
    MaxRounds    int
}

func (s *Service) runAgentLoop(ctx context.Context, opts agentRunOptions) (openai.ChatCompletionMessage, error)
```

主 Agent 使用：

- `PersistStop = true`
- `EmitStream = true`
- `Tools = s.Tools`
- `State` 使用现有 conversation history

子 Agent 使用：

- `PersistStop = false`
- `EmitStream = false`
- `Tools = childTools`
- `State.Messages` 为全新列表
- `State.History` 仅用于子循环内部串联工具结果，不写入主 conversation history

实现时抽取该通用循环，避免复制两份完整循环逻辑。

### 8.2 子 Agent messages 初始化

子 Agent 的 messages 应只包含：

1. 子 Agent system prompt。
2. 一条 user message，内容来自 `task`。

示例：

```go
messages := []openai.ChatCompletionMessage{
    {Role: "system", Content: childSystemPrompt},
    {Role: "user", Content: task},
}
```

`childSystemPrompt` 基于现有 `buildSystemPrompt`，并追加子 Agent 专属约束：

- 你是被 `spawn_subagent` 启动的子 Agent。
- 你看不到父 Agent 的对话历史；只根据当前 task 工作。
- 你可以修改 workspace 中的文件，副作用会保留。
- 你不能调用 `spawn_subagent`；如果需要更多上下文，只能使用其他允许工具读取 workspace。
- 完成后只输出简洁总结，说明做了什么、关键发现和是否有未完成事项。

### 8.3 子 Agent 不持久化对话上下文

子 Agent 的完整 `messages[]` 不写入 `SetConversationHistory`，也不进入主 `state.Messages`。主 Agent 只会看到 `spawn_subagent` tool message，其中内容是子 Agent 最终摘要。

实现方式：

- 子 Agent stop 时不运行会持久化 assistant history 的 Stop hook，或运行一个仅构造返回值、不落库的 stop 流程。
- 子 Agent 的工具 post hook 继续把 tool result 追加到子 Agent 自己的 `State.Messages`，以便子 Agent 后续轮次读取工具结果；但这个 `State.History` 是内存态，循环结束后丢弃。

### 8.4 子 Agent 工具 hook

子 Agent 工具调用仍按顺序执行：

1. `RunPreToolUse`
2. 执行工具
3. `RunPostToolUse`

为避免子 Agent 内部过程污染用户前端主聊天流：

- 子 Agent `LoopState.Writer` 设为 `nil`，因此默认 SSE tool event 不会发给前端。
- `persistToolCallHook` 继续写入 `tool_calls` 表，用于审计；记录仍关联父 conversation ID。第一版不新增字段区分主/子 Agent 工具调用，避免数据库或审计结构扩展。
- `appendToolMessageHook` 仍可运行，因为它只追加到子 Agent 内存态 `State.Messages` / `State.History`，循环结束后丢弃。

如果后续需要区分主 Agent 与子 Agent 工具调用，可在 `ToolExecutionAudit` 增加 `agent_scope` / `parent_tool_call_id` 字段；这不作为第一版范围。

## 9. 工作目录与安全设计

### 9.1 cwd 解析规则

新增 helper：

```go
func resolveSubagentCWD(workspaceRoot, cwd string) (string, error)
```

规则：

1. `workspaceRoot` 必须存在且为目录。
2. `cwd == ""` 时返回 `workspaceRoot`。
3. 相对路径用 `filepath.Join(workspaceRoot, cwd)` 解析。
4. 绝对路径直接 clean/abs。
5. 解析结果必须等于 `workspaceRoot` 或以 `workspaceRoot + pathSeparator` 为前缀。
6. 解析结果必须存在且为目录。

### 9.2 子 Agent RuntimeEnv

子 Agent 使用独立 `RuntimeEnv`：

```go
childEnv := parentEnv
childEnv.CurrentWorkingDir = resolvedCWD
childEnv.AllowOutsideWorkspace = false
```

即使主配置允许 bash 访问 workspace 外部，子 Agent 也强制不能访问 workspace 外部，满足“子agent的工作目录，只能在workspace下”的要求。

`AllowDangerousCommands` 继承主配置，因为它是全局安全策略；但路径逃逸能力必须强制关闭。

### 9.3 不能递归 spawn

通过工具定义层和执行层双重保护：

1. 子 Agent 使用 `childTools`，其 definitions 中不包含 `spawn_subagent`。
2. `runSubagent` 在 context 中设置 `subagentDepth = 1`。
3. 如果 `spawn_subagent` 执行入口发现当前 context 已是子 Agent，则直接返回错误：`spawn_subagent cannot be called from a subagent`。

这样即使模型伪造工具名或白名单配置错误，也不会递归启动。

## 10. 工具白名单与注册设计

### 10.1 工具定义集合

当前 `ChildToolDefs = baseToolDefs` 命名容易混淆。实现时调整为：

```go
var baseToolDefs = []openai.Tool{...}
var parentToolDefs = append(baseToolDefs, spawnSubagentToolDef)
var childToolDefs = baseToolDefs
```

或者保留导出名但明确：

- `AllToolDefs`：包含 `spawn_subagent`。
- `ChildToolDefs`：不包含 `spawn_subagent`。

### 10.2 WebAllowedTools

`spawn_subagent` 只有在 `WebAllowedTools` 包含它时才暴露给主 Agent。第一版采用“保守默认”：工具已注册，但不加入 `LoadWebConfig` 的默认 `WEB_ALLOWED_TOOLS`；只有配置 `WEB_ALLOWED_TOOLS=...,spawn_subagent` 才暴露。原因是它会消耗额外模型调用并可能产生文件系统副作用，应由部署方显式开启。

子 Agent 的工具集合从同一份 `WebAllowedTools` 派生，但强制移除 `spawn_subagent`。例如主 Agent 配置为 `bash,read_file,write_file,edit_file,load_skill,spawn_subagent` 时，子 Agent 可用 `bash/read_file/write_file/edit_file/load_skill`；如果主 Agent 只配置了 `spawn_subagent`，子 Agent 将没有其他工具，只能直接用模型回答 task。

## 11. 错误处理与限制

### 11.1 MaxRounds

为避免子 Agent 无限循环，新增固定上限：

```go
const defaultSubagentMaxRounds = 20
```

超过上限返回错误：`subagent exceeded max rounds`。

主 Agent 收到的是失败工具结果，可向用户解释或继续处理。

### 11.2 Context 与取消

子 Agent 继承父工具调用的 `ctx`。如果用户断开、请求取消或服务超时，子 Agent 模型调用和工具调用应一起取消。

### 11.3 模型空响应

复用现有 `runModelRoundStream` 的空 stream 检测；子 Agent 遇到空响应直接失败。

### 11.4 工具执行失败

普通工具执行失败仍以 tool result 的形式返回给子 Agent，让子 Agent 自行总结或修复。只有 Runtime 级错误（例如 cwd 非法、模型调用失败、超过轮数）才让 `spawn_subagent` 返回失败。

## 12. 测试设计

### 12.1 单元测试

1. `spawn_subagent` 工具定义存在
   - 配置 `WebAllowedTools = []string{"spawn_subagent"}`。
   - 断言主 `ToolRegistry.Definitions()` 包含 `spawn_subagent`。

2. 子 Agent 工具定义不包含 `spawn_subagent`
   - 构造 child registry。
   - 断言 definitions 中没有 `spawn_subagent`。

3. 子 Agent 使用全新 messages
   - 主 conversation 预置 history。
   - fake LLM 捕获子 Agent 首次请求。
   - 断言子请求只包含 system + task user message，不包含主历史内容。

4. 子 Agent 最终摘要作为工具结果返回
   - fake LLM：主 Agent 第一次请求调用 `spawn_subagent`，子 Agent 返回 `done summary`，主 Agent 第二次请求基于工具结果返回最终回答。
   - 断言主 Agent 的 tool message 包含 `done summary`。

5. 子 Agent cwd 必须在 workspace 内
   - `cwd = "subdir"` 且目录存在：成功。
   - `cwd = "../outside"`：失败。
   - `cwd = "/tmp"`：失败。

6. 子 Agent 不能递归 spawn
   - fake LLM 让子 Agent 返回 `spawn_subagent` tool call。
   - 断言执行结果为 rejected / error，且没有启动第二个子 Agent。

7. 子 Agent 工具调用经过 hooks
   - 注入自定义 `PreToolUse` / `PostToolUse` hook 计数。
   - 子 Agent 调用 `read_file` 或 `bash`。
   - 断言 hook 被调用。

### 12.2 集成测试

1. 文件系统副作用保留
   - 子 Agent 调用 `write_file` 写入 workspace 文件。
   - `spawn_subagent` 返回后，测试读取文件并断言内容存在。

2. 子 Agent 对话上下文不持久化
   - 子 Agent 执行多轮工具调用。
   - 最终 `SetConversationHistory` 只包含主会话用户消息、主 assistant/tool 轨迹和最终 assistant，不包含子 Agent 内部 user/assistant/tool 明细。

3. 子 Agent 无 SSE 输出
   - 主会话传入 writer。
   - 子 Agent 内部工具调用不直接发 `tool` event；前端只看到主 Agent 的 `spawn_subagent` 工具事件和最终 assistant。

## 13. 实施步骤

1. 在 `internal/tools/definitions.go` 增加 `spawn_subagent` 工具定义，并区分 parent/child tool definitions。
2. 调整 `internal/web/runtime/tool_catalog.go`，主 registry 可识别 `spawn_subagent`，child registry 明确排除它。
3. 在 Runtime 层新增 `runSubagent` 和 `resolveSubagentCWD`。
4. 抽取或复用主 Agent loop，使主 Agent 与子 Agent 都走 `RunPreToolUse -> execute -> RunPostToolUse`。
5. 在 `executeToolCall` 中拦截 `spawn_subagent`，解析参数并启动子 Agent。
6. 为子 Agent 构造独立 `LoopState`、独立 `messages[]`、独立 runtime env、无 writer。
7. 添加最大轮数与递归保护。
8. 补充单元测试和集成测试。

## 14. 兼容性

1. 现有工具名和参数不变。
2. 未配置 `spawn_subagent` 时，现有行为不变。
3. 主 conversation history 仍按现有格式持久化。
4. `spawn_subagent` 的文件系统副作用与普通工具一致，受 workspace guard 约束。
5. 子 Agent 不继承主对话上下文，因此主 Agent 必须在 `task` 中提供足够信息；这是预期行为。

## 15. 设计决策

1. `spawn_subagent` 第一版不默认加入 `WEB_ALLOWED_TOOLS`，显式配置后使用。
2. 子 Agent 内部工具调用第一版不直接展示在 UI 中，只保留现有审计记录。
3. `defaultSubagentMaxRounds` 第一版取 20。
