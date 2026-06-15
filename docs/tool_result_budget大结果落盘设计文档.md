# tool_result_budget 大结果本地落盘设计文档

## 1. 背景

当前 `cynosure` 已迁移到 TUI，本地会话已经具备以下基础能力：

- 每个 TUI 会话拥有独立 `session_id`，会话创建点在 `internal/local/bootstrap.go:86`。
- 展示历史与模型历史会写入 `~/.cynosure/session/{session_id}/history` 与 `~/.cynosure/session/{session_id}/model_history`，写入入口在 `internal/local/store.go:121` 和 `internal/local/store.go:351`。
- 上下文压缩链路已经有 `ToolResultCompressionStrategy`：当最近一轮 tool_result 总字节数超过 `200 * 1024` 时，会按结果大小从大到小替换为 `<persisted-output>` 标记与前 2000 字符预览，核心逻辑在 `internal/agent/runtime/compression/tool_result_compression.go:76`。
- 当前大结果保存仍依赖 `Store.CreatePersistedOutput` 的本地内存 map，进程退出后不可恢复，见 `internal/local/store.go:185`。
- `read_persisted_output` 已存在，可按 id 分段读取 persisted output，处理入口在 `internal/tools/persisted_output.go:42`，读取适配器在 `internal/agent/runtime/context_compression.go:114`。
- 工具执行结果目前只追加到内存中的 `state.History` / `state.ModelHistory`，追加点在 `internal/agent/runtime/hooks/tool.go:42`，尚未写入用户目录的 `~/.cynosure/task_outputs/{session_id}/tools.md`。

本需求要求把 tool_result_budget 触发的大结果落盘到用户目录 `~/.cynosure/task_outputs/`，并把所有工具执行结果追加写入该目录下的 Markdown 文件，便于 TUI 本地使用、会话排查和跨进程恢复。

## 2. 需求与边界

### 2.1 目标

1. 当 tool_result_budget 触发大结果压缩时，将完整工具结果落盘到用户目录：
   - 目录：`~/.cynosure/task_outputs/tool-results/`
   - 策略：从最大的 tool_result 开始落盘，直到当前上下文内 tool_result 总字节数降到预算以内。
   - 上下文保留：`<persisted-output>` 标记 + 前 2000 字符预览。
2. 工具每次执行完成后，将工具执行结果追加写入：
   - 文件：`~/.cynosure/task_outputs/{session_id}/tools.md`
   - 写入方式：追加写，保留历史工具执行记录。
3. `read_persisted_output` 继续可用：模型看到 `<persisted-output>` 后，仍能通过 id 分段读取完整结果。

### 2.2 非目标

1. 不改变工具本身的执行逻辑、权限模型和 TUI 展示样式。
2. 不改变 `~/.cynosure/session/{session_id}/history` / `model_history` 的会话历史格式。
3. 不新增远端存储或数据库依赖。
4. 不把 `~/.cynosure/task_outputs/{session_id}/tools.md` 的内容自动注入 LLM 上下文；上下文仍只通过模型历史与压缩策略控制。

### 2.3 明确假设

1. TUI 仍使用 `workspaceRoot` / `CWD` 标识当前项目与校验会话归属，但工具输出文件统一写入用户目录 `~/.cynosure/task_outputs/`。
2. 需求中的 `task_outputs /{session_id}/tools.md` 视为路径书写空格误差，实际实现为 `~/.cynosure/task_outputs/{session_id}/tools.md`。
3. `~/.cynosure/task_outputs/tool-results/` 不再额外按 session 分目录；文件名中包含 `session_id` 与 persisted output id，避免冲突并便于定位。
4. 大结果文件和 `tools.md` 均写入 `~/.cynosure/task_outputs/`；这与历史会话文件仍写在 `~/.cynosure/session` 下互不替代。

## 3. 方案对比

### 方案 A：只扩展 `local.Store.CreatePersistedOutput`（推荐）

在本地 `Store` 的 persisted output 创建/查询路径中增加文件落盘：

- `CreatePersistedOutput` 同时写内存 map 和 `~/.cynosure/task_outputs/tool-results/{session_id}-{id}.txt`。
- `GetPersistedOutputForConversation` 内存命中则直接返回；内存未命中时从文件恢复内容。
- `GetPersistedOutputByMessageHash` 继续使用内存索引；如需跨进程恢复，可读取 sidecar metadata 重建索引。

优点：改动集中，复用现有压缩策略、marker 与 `read_persisted_output`。缺点：需要为文件读取补一层 metadata 或稳定文件名查找。

### 方案 B：把文件落盘逻辑放入 `ToolResultCompressionStrategy`

压缩策略直接知道工作区路径并写文件。

优点：压缩和落盘绑定最直接。缺点：会让 compression 层依赖本地文件系统与 TUI 工作区，破坏当前 `compression.Store` 抽象，影响测试和未来非本地运行形态。

### 方案 C：新增独立 `PersistedOutputStore` 组件并注入 runtime

将本地文件 persisted output 作为单独组件，runtime 和 compression 都依赖它。

优点：边界最清晰，长期演进更好。缺点：当前需求较小，会引入较多接口迁移和注入改造。

### 推荐

采用方案 A，并在 `local` 包内部新增小型文件辅助类型。理由：当前 `ToolResultCompressionStrategy` 已经按“大到小、预算以内、marker + 2000 字符预览”完成核心语义；需求缺口主要是本地持久化介质和工具结果审计文件。将能力补在 `local.Store` 能最小化改动，并保持压缩层只依赖 `compression.Store`。

## 4. 设计

### 4.1 目录结构

用户目录下新增：

```text
~/.cynosure/task_outputs/
├── tool-results/
│   ├── {session_id}-{persisted_output_id}.txt
│   └── {session_id}-{persisted_output_id}.json
└── {session_id}/
    └── tools.md
```

说明：

- `.txt` 保存完整工具输出内容。
- `.json` 保存 metadata，用于跨进程读取与校验。
- `{session_id}` 必须复用现有 `validSessionID` 规则，拒绝路径穿越。
- `{persisted_output_id}` 仅允许安全文件名字符；当前 id 格式如 `po_...`，仍需在写文件前做防御性校验。

### 4.2 persisted output 文件格式

Metadata 文件建议格式：

```json
{
  "version": 1,
  "id": "po_xxx",
  "session_id": "041581e7-c3e7-46c8-afe7-7cdcc671e80e",
  "conversation_id": "conv_xxx",
  "user_id": "local-user",
  "message_id": "msg_xxx",
  "tool_call_id": "call_xxx",
  "kind": "tool_result",
  "strategy": "tool_result_compression",
  "original_bytes": 123456,
  "content_sha256": "...",
  "preview": "前 2000 字符预览",
  "content_file": "041581e7-c3e7-46c8-afe7-7cdcc671e80e-po_xxx.txt",
  "created_at": "2026-06-12T00:00:00Z"
}
```

写入顺序：

1. 先原子写 `.txt` 完整内容。
2. 再原子写 `.json` metadata。
3. 最后更新内存 map 和 hash 索引。

如果 metadata 写失败，应返回错误，让压缩策略失败并避免上下文里出现不可读取的 marker。

### 4.3 marker 与上下文内容

保持现有 marker 格式不变：

```text
<persisted-output id="po_xxx" kind="tool_result" original_bytes="123456" preview_chars="2000" retrieval_tool="read_persisted_output">
完整输出已持久化；如需更多内容，请调用 read_persisted_output(id="po_xxx", offset=0, limit=20000) 分段读取。

前 2000 字符预览
</persisted-output>
```

原因：

- 现有 `isCompactedResult` 通过 `<persisted-output` 判断已压缩，见 `internal/agent/runtime/compression/tool_result_compression.go:45`。
- `read_persisted_output` 的入参仍是 id，不需要模型理解本地文件路径，避免把工作区绝对路径暴露进上下文。

### 4.4 “从最大的开始落盘”策略

保持 `ToolResultCompressionStrategy.Apply` 当前算法：

1. 只处理最近一个 user turn 之后的 tool messages。
2. 计算所有未压缩 tool_result 的总字节数。
3. 若总字节数不超过 `toolResultByteThreshold`，不处理。
4. 若超过阈值，按单个 result 字节数倒序排序。
5. 从最大结果开始逐个调用 `persistAndBuildMarker`，直到剩余总字节数不超过阈值。

该逻辑已经存在于 `internal/agent/runtime/compression/tool_result_compression.go:76` 到 `internal/agent/runtime/compression/tool_result_compression.go:115`，实现阶段只需要确保 `CreatePersistedOutput` 真正落盘。

### 4.5 `read_persisted_output` 读取流程

`GetPersistedOutputForConversation` 增加 fallback：

1. 先按当前内存 map 查找，并校验 `user_id` 与 `conversation_id`。
2. 未命中时，通过 `conversation_id` 找到当前 `Conversation`，取得 `session_id`。
3. 在 `~/.cynosure/task_outputs/tool-results/` 中读取 `{session_id}-{id}.json`。
4. 校验 metadata：
   - id 一致。
   - conversation_id 一致。
   - user_id 一致。
   - kind 为 `tool_result`。
5. 读取 `.txt` 内容并计算 sha256，与 metadata 中的 `content_sha256` 对比。
6. 返回 `storage.PersistedOutput`。

这样 TUI 重启、恢复历史会话后，只要 `~/.cynosure/task_outputs/tool-results/` 仍在，模型就可以继续通过 marker id 读取完整输出。

### 4.6 `~/.cynosure/task_outputs/{session_id}/tools.md` 追加写

在工具执行后 hook 中追加写入：

- 当前追加工具消息的入口是 `appendToolMessageHook`，见 `internal/agent/runtime/hooks/tool.go:42`。
- 该 hook 目前只更新 `state.Messages`、`state.History`、`state.ModelHistory`。
- 设计新增一个 Store 能力，例如 `AppendToolResultLog(ctx, conversationID, userID, toolCallID, toolName, rawArgs, outcome)`，由本地 Store 负责写入 `~/.cynosure/task_outputs/`。

追加格式建议：

```markdown
## 2026-06-12T15:04:05Z · bash · success

- conversation_id: conv_xxx
- session_id: 041581e7-c3e7-46c8-afe7-7cdcc671e80e
- tool_call_id: call_xxx

### Arguments

```json
{"command":"go test ./..."}
```

### Result

```text
完整工具结果或错误信息
```
```

注意：Markdown 代码围栏需根据内容动态选择足够长的反引号长度，避免结果中包含 ``` 时破坏文件结构。

追加策略：

- 使用 `os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)`。
- 写入前 `os.MkdirAll(filepath.Dir(path), 0o755)`。
- 单次追加在 `local.Store` 的互斥锁内完成，避免同一进程并发工具调用交错写入。
- 如果追加失败，不阻断工具结果进入模型上下文；记录 warn 日志即可。原因是 tools.md 是审计/排查产物，不应改变工具执行语义。

### 4.7 接口变化

建议在 runtime 层新增窄接口，而不是让 hooks 直接依赖 `local.Store`：

```go
type toolResultLogStore interface {
    AppendToolResultLog(ctx context.Context, entry storage.ToolResultLogEntry) error
}
```

`storage.ToolResultLogEntry` 包含：

- ConversationID
- SessionID（可选；Store 可自行从 conversation 映射推导）
- UserID
- ToolCallID
- ToolName
- RawArgs
- Status
- Result
- AuditSummary
- CreatedAt

`appendToolMessageHook` 中判断 `state.Store` 是否实现该接口；实现则追加，不实现则跳过，保持测试 fake store 与未来其他 Store 兼容。

### 4.8 与现有历史持久化的关系

- `history`：继续保存展示历史，包含完整工具结果或压缩后的 marker，取决于当轮写入时机。
- `model_history`：继续保存压缩后的模型历史，用于下轮上下文。
- `~/.cynosure/task_outputs/tool-results/`：保存被 tool_result_budget 移出上下文的大结果全文。
- `~/.cynosure/task_outputs/{session_id}/tools.md`：保存每次工具执行的结果日志，不参与上下文压缩。

## 5. 错误处理

1. 大结果 `.txt` 或 `.json` 写失败：返回错误，中止本轮压缩，避免 marker 指向不可读内容。
2. 读取 persisted output 时文件不存在：维持现有 `persisted output "id" not found` 错误语义。
3. sha256 校验失败：返回错误，不返回可能损坏的内容。
4. `tools.md` 追加失败：记录日志，不影响工具消息追加到模型上下文。
5. session_id 或 output id 非法：拒绝写入/读取，避免路径穿越。
6. 用户目录下 `~/.cynosure/task_outputs/` 不存在：自动创建；若创建失败按上述策略处理。

## 6. 测试计划

### 6.1 单元测试

1. `ToolResultCompressionStrategy` 现有测试继续覆盖：
   - 未超阈值不落盘。
   - 超阈值后从最大结果开始落盘。
   - marker 保留前 2000 字符预览。
2. `local.Store.CreatePersistedOutput` 新增测试：
   - 创建 `~/.cynosure/task_outputs/tool-results/{session_id}-{id}.txt/.json`。
   - 文件内容等于完整 result。
   - metadata 字段完整，sha256 正确。
3. `local.Store.GetPersistedOutputForConversation` 新增测试：
   - 内存命中正常返回。
   - 清空内存后可从文件恢复。
   - user_id / conversation_id 不匹配时拒绝。
   - sha256 不匹配时报错。
4. `AppendToolResultLog` 新增测试：
   - 首次创建 `~/.cynosure/task_outputs/{session_id}/tools.md`。
   - 多次调用追加多段记录。
   - 结果中包含反引号时 Markdown 仍结构完整。
5. hook 测试：
   - 工具成功/失败后均调用追加日志。
   - 日志追加失败不影响 `state.History` / `state.ModelHistory`。

### 6.2 集成验证

1. 在临时工作区触发一个超过 200KB 的工具输出。
2. 确认模型上下文中只剩 `<persisted-output>` + 2000 字符预览。
3. 确认用户目录生成 `~/.cynosure/task_outputs/tool-results/*.txt` 与 `.json`。
4. 调用 `read_persisted_output(id, offset, limit)` 能读取完整内容分片。
5. 确认 `~/.cynosure/task_outputs/{session_id}/tools.md` 追加了工具名、参数、状态与结果。
6. 重启 TUI 并 `/resume` 同一会话后，再调用 `read_persisted_output` 仍可读取文件中的完整输出。

## 7. 实施步骤建议

1. 在 `internal/local` 中新增 persisted output 文件辅助逻辑：路径构造、安全文件名校验、metadata 结构、原子写、文件读取与 sha256 校验。
2. 修改 `local.Store.CreatePersistedOutput`：保持内存写入，同时写入 `~/.cynosure/task_outputs/tool-results/`。
3. 修改 `local.Store.GetPersistedOutputForConversation`：增加文件 fallback。
4. 新增 `storage.ToolResultLogEntry` 与 `AppendToolResultLog` 能力。
5. 在工具 post hook / append hook 后追加写 `~/.cynosure/task_outputs/{session_id}/tools.md`。
6. 补充单元测试和集成测试。

## 8. 验收标准

1. 超过 tool_result_budget 时，大结果从最大的开始落盘到 `~/.cynosure/task_outputs/tool-results/`。
2. 模型上下文中被落盘的大结果只保留 `<persisted-output>` 标记与前 2000 字符预览。
3. `read_persisted_output` 能读取本轮和恢复会话后的 persisted output。
4. 每次工具执行后，`~/.cynosure/task_outputs/{session_id}/tools.md` 都会追加对应结果。
5. `go test ./...` 通过。
