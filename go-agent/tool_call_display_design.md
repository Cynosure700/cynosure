# 工具调用简洁展示与会话重载设计文档

## 1. 背景

当前 Web 前端已经能在本轮 SSE 中收到工具调用事件，并在右侧“工具调用”面板中展示工具名、状态和完整执行结果。现状存在两个问题：

1. 展示过重：`web/src/App.tsx` 当前直接把 `event.result` 放进 `<pre>`，长命令输出会占据大量空间，不符合“只需要展示调用了什么工具以及部分执行结果”的目标。
2. 重新打开会话后工具信息丢失：`GET /api/conversations/{id}` 只返回过滤后的用户/助手消息，`displayConversationMessages` 会过滤掉 `tool` 消息和带 `ToolCalls` 的中间 assistant 消息，因此前端无法恢复历史工具调用信息。

相关现有位置：

- 前端工具事件类型与 SSE 消费：`web/src/api.ts`
- 前端右侧工具面板展示：`web/src/App.tsx`
- 当前会话详情接口：`go-agent/internal/web/app/conversation_handlers.go`
- 工具调用 SSE 事件与历史 tool message 写入：`go-agent/internal/web/runtime/hooks/tool.go`
- 消息模型中的 `ToolCallID` 和 `ToolCalls`：`go-agent/internal/web/storage/models.go`
- `tool_calls` 表当前仅写入，不提供查询接口：`go-agent/internal/web/storage/tool_calls_repo.go`

## 2. 目标

本次改造目标：

1. 工具调用面板保持简洁，只展示：
   - 工具名
   - 执行状态
   - 执行结果预览，不展示完整原始输出
2. 重新打开历史会话时，工具调用面板能够加载该会话历史工具调用信息。
3. 保持当前主聊天消息流简洁，不把 tool message 混入聊天气泡。
4. 不新增数据库迁移，优先复用已持久化的 conversation history。
5. 不影响模型使用完整工具结果进行后续推理；截断仅用于 UI 展示。

## 3. 当前数据来源分析

工具调用相关数据有两个来源：

### 3.1 conversation history 中的 assistant/tool 消息

runtime 在工具调用流程中会写入：

- assistant 中间消息：包含 `ToolCalls`，其中有工具调用 ID、工具名、参数。
- tool 消息：包含 `ToolCallID` 和工具执行结果 JSON，内容来自 `ToolExecutionOutcome.MessageContent()`，格式包括 `status` 和 `result`。

优点：

- 已经存在于会话历史中，重新打开会话时可通过现有 `ListMessagesByConversation` 读取。
- `tool` 消息中保留真实执行结果，适合生成 UI 预览。
- 不需要新增表字段或迁移。

缺点：

- 需要后端从完整 history 中派生一个专门给 UI 的 `tool_events` 数组。
- 需要通过 `ToolCallID` 将 tool 消息关联回工具名。

### 3.2 `tool_calls` 表

当前 `CreateToolCall` 会写入 `tool_calls` 表，字段包括工具名、状态和 summary。

优点：

- 表结构天然代表“工具调用记录”。
- 后续如果需要审计列表，可以继续扩展。

缺点：

- 当前没有 List 查询方法。
- `summary` 当前是 audit JSON，不是用户友好的执行结果预览。
- 若只用该表，无法展示“部分执行结果”，只能展示 audit summary。

## 4. 方案对比

### 方案 A：直接把历史 tool 消息放回 messages

后端不新增字段，`displayConversationMessages` 不再过滤 tool 消息，前端在 messages 中识别 tool 消息并展示。

优点：

- 改动最少。
- 不需要新增 API 字段。

缺点：

- 主聊天消息会混入工具内部消息，破坏聊天区简洁性。
- 前端需要在消息列表里处理 tool role，但当前 `ChatMessage.role` 只支持 `user | assistant`。
- 工具调用与主回答耦合，后续 UI 维护成本高。

结论：不推荐。

### 方案 B：查询 `tool_calls` 表返回工具列表

新增 `ListToolCallsByConversation`，会话详情接口返回 `tool_calls` 表数据。

优点：

- 数据结构清晰，工具调用列表与消息历史分离。
- 查询成本低，按 conversation_id 索引即可。

缺点：

- 当前 `summary` 不是展示友好的 result preview。
- 要满足“部分执行结果”，需要改写 `summary` 写入语义，或新增字段；这会引入更大变更。
- 历史已写入数据仍然可能只有 audit JSON。

结论：可作为后续审计能力扩展，但不作为本次首选。

### 方案 C：从 conversation history 派生 `tool_events`，并返回给前端【推荐】

后端读取完整会话 history 后，继续用 `displayConversationMessages` 返回主聊天消息，同时新增 `displayConversationToolEvents` 派生工具事件数组：

```json
{
  "conversation": { ... },
  "messages": [ ... ],
  "tool_events": [
    {
      "id": "tool_1",
      "name": "bash",
      "status": "success",
      "result": "pwd 输出预览...",
      "truncated": false
    }
  ]
}
```

优点：

- 不新增数据库迁移。
- 能拿到真实执行结果，并只给 UI 返回截断预览。
- 主聊天消息仍然保持简洁。
- SSE 实时工具事件与历史加载返回同一前端类型，前端状态逻辑简单。

缺点：

- 后端需要从 assistant ToolCalls 与 tool message 做一次关联。
- 如果极早期历史缺失 ToolCalls 元数据，工具名可能只能降级显示为 `tool` 或 ToolCallID。

结论：推荐采用。

## 5. 推荐设计

采用方案 C：**从 conversation history 派生工具展示事件，SSE 与历史 API 统一返回简洁 ToolEvent。**

### 5.1 后端 API 响应

修改 `GET /api/conversations/{id}` 响应，新增 `tool_events`：

```json
{
  "conversation": { ... },
  "messages": [ ... ],
  "tool_events": [
    {
      "id": "tool_1",
      "name": "bash",
      "status": "success",
      "result": "/Users/bytedance/golang_pro/nano_cc\n",
      "truncated": false
    }
  ]
}
```

字段含义：

- `id`：优先使用 `tool_call_id`，方便 React key 和排查。
- `name`：工具名，如 `bash`、`read`、`write`。
- `status`：`success` / `rejected` / `error` 等现有状态。
- `result`：展示用预览文本，不是完整工具结果。
- `truncated`：后端是否截断过结果。前端可在预览末尾展示“已截断”。

### 5.2 后端派生逻辑

新增 helper：

```go
func displayConversationToolEvents(messages []storage.Message) []toolEventPayload
```

处理流程：

1. 遍历完整 messages。
2. 遇到 assistant 且 `ToolCalls` 非空时，建立映射：
   - `toolCallID -> toolName`
3. 遇到 role 为 `tool` 的消息时：
   - 解析 `message.Content` 中的 JSON：`status`、`result`。
   - 用 `message.ToolCallID` 找工具名。
   - 对 result 做展示截断。
   - 追加到 `tool_events`。
4. 如果解析失败，使用 `message.Content` 作为 result 预览，status 设为 `unknown`。
5. 如果找不到工具名，降级为 `tool`。

### 5.3 结果预览规则

后端统一生成 preview，避免前端处理超长输出。

建议规则：

- 去掉首尾空白。
- 最多保留 6 行。
- 最多保留 300 个 rune。
- 超出时追加 `…`，并设置 `truncated: true`。
- 空结果显示为 `(无输出)`。

该规则只影响 UI 展示，不改变 tool message 中给模型使用的完整结果。

### 5.4 SSE 实时工具事件

当前 SSE `tool` 事件 payload 是：

```json
{"name":"bash","status":"success","result":"完整输出"}
```

调整为与历史 API 一致：

```json
{"id":"tool_1","name":"bash","status":"success","result":"预览输出","truncated":true}
```

注意：这里的 `id` 可以使用 OpenAI tool call id，即 `h.ToolCall.ID`，不使用 `tool_calls` 表的 `tc_` ID。这样 SSE 实时事件与 history 中的 `ToolCallID` 能保持一致。

### 5.5 前端类型与加载流程

扩展 `ToolEvent` 类型：

```ts
type ToolEvent = {
  id?: string;
  name: string;
  status: string;
  result: string;
  truncated?: boolean;
};
```

扩展 `api.getConversation` 返回：

```ts
{
  conversation: Conversation;
  messages: ChatMessage[];
  tool_events: ToolEvent[];
}
```

会话加载时同时设置：

```ts
setMessages(result.messages)
setToolEvents(result.tool_events)
```

这样重新打开会话、切换会话时右侧工具面板都能恢复历史工具调用信息。

### 5.6 前端简洁展示

右侧工具调用面板从“完整 pre 输出”改成简洁卡片：

```tsx
<div className="tool-event compact">
  <div className="tool-event-head">
    <strong>{event.name}</strong>
    <span className={`status ${event.status}`}>{event.status}</span>
  </div>
  <p className="tool-event-preview">
    {event.result}
    {event.truncated ? " …已截断" : ""}
  </p>
</div>
```

展示原则：

- 默认只展示工具名 + 状态 + 预览结果。
- 不展示完整 JSON。
- 不展示工具参数，除非后续有明确需求。
- 如果没有工具调用，文案调整为“当前会话暂无工具调用”。

## 6. 需要修改的文件

后端：

1. `go-agent/internal/web/app/conversation_handlers.go`
   - 新增 `toolEventPayload` 类型。
   - 新增 `displayConversationToolEvents`。
   - `GET /api/conversations/{id}` 响应增加 `tool_events`。
   - 增加 result preview helper。

2. `go-agent/internal/web/runtime/hooks/tool.go`
   - SSE `tool` 事件改为发送展示用 preview。
   - payload 增加 `id` 和 `truncated`。

3. `go-agent/internal/web/app/server_test.go` 或相关测试文件
   - 覆盖会话详情接口返回 `tool_events`。
   - 覆盖长结果截断。
   - 覆盖 tool message JSON 解析失败降级。

4. `go-agent/internal/web/runtime/runtime_test.go`
   - 更新工具事件断言，确认 SSE `tool` 事件返回预览而非完整结果。

前端：

1. `web/src/api.ts`
   - 导出 `ToolEvent` 类型，或在 App 内保持一致类型。
   - `getConversation` 解析 `tool_events`。
   - `streamConversation.onTool` payload 类型增加 `id` 与 `truncated`。

2. `web/src/App.tsx`
   - 会话加载时设置 `toolEvents`。
   - 工具调用面板改为简洁展示。
   - 切换/删除/新建会话时保持现有清空逻辑。

3. `web/src/styles.css`
   - 调整 `.tool-event`、`.tool-event-head`、`.tool-event-preview`，弱化 `<pre>` 大块展示。

## 7. 测试计划

### 7.1 后端测试

运行：

```bash
cd /Users/bytedance/golang_pro/nano_cc/go-agent && go test ./...
```

建议新增测试点：

1. `GET /api/conversations/{id}` 返回 `tool_events`。
2. tool event 能从 assistant ToolCalls + tool message 关联出工具名。
3. 长工具结果被截断，`truncated=true`。
4. 空工具结果显示 `(无输出)`。
5. 非法 tool message JSON 不导致接口失败。

### 7.2 前端验证

运行：

```bash
cd /Users/bytedance/golang_pro/nano_cc/web && npm run build
```

验证点：

1. 新会话工具调用时，右侧面板实时出现简洁工具卡片。
2. 工具结果很长时，只显示预览，不撑开面板。
3. 切换到历史会话时，工具调用面板恢复历史工具记录。
4. 无工具调用时显示“当前会话暂无工具调用”。

## 8. 边界与取舍

1. 本次不新增工具调用详情弹窗，也不支持展开完整输出。
2. 本次不新增数据库迁移，避免扩大改动面。
3. `tool_calls` 表继续保留现有审计用途；本次 UI 历史展示优先从 conversation history 派生。
4. 如果后续需要跨会话搜索/审计工具调用，再单独设计 `tool_calls` 查询接口和更明确的 summary/result 字段。

## 9. 实施顺序建议

1. 后端先添加 `tool_events` 派生逻辑和测试。
2. 调整 SSE `tool` 事件为 preview payload，并更新 runtime 测试。
3. 前端扩展 API 类型和会话加载逻辑。
4. 前端简化工具调用 UI 样式。
5. 运行 `go test ./...` 和 `npm run build`。

## 10. 自审结论

- 没有引入新数据库表或迁移，符合最小改动原则。
- 设计覆盖了实时工具调用展示和历史会话重载两个需求。
- UI 展示明确限定为工具名、状态、部分执行结果，避免完整输出造成界面噪音。
- conversation history 中的完整工具结果仍保留给模型后续推理使用，UI 截断不会影响 agent 行为。
