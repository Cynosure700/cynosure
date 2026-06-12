# 记忆系统 Markdown 存储改造设计文档

## 1. 背景与目标

当前 go-agent 已迁移到 TUI 形态，启动入口位于 `internal/local/bootstrap.go`，运行时仍复用 `internal/agent/runtime` 的对话、提示词、压缩与记忆逻辑。现有记忆相关代码主要沿用服务端/数据库模型：

- 长期记忆：`runtime.extractMemories` 在每轮结束后抽取 `storage.Memory`，通过 `Store.InsertMemory` 持久化；启动/请求前由 `runtime.selectRelevantMemories` 从 `Store.ListRelevantMemories` 获取候选记忆，再让模型选择要注入的记忆。
- 当前会话记忆：`runtime.updateConversationMemory` 在每轮结束后读取并替换 `storage.ConversationMemory` 列表，目前以 `conversation.ID` 作为会话维度；改造后会显式引入随机 UUID 格式的 `session_id`，用于当前会话记忆文件和会话锁。
- TUI 本地 Store：`internal/local/store.go` 当前是内存实现，`ListRelevantMemories`、项目事实记忆替换、长期记忆替换等方法尚未真正落盘。
- TUI 启动时在 `Bootstrap` 中将 `runtimeService.EnableMemory = false`，本地用户 `MemoryEnabled = false`，因此记忆功能实际未启用。

本次改造目标：

1. 将记忆系统存储迁移到当前工作目录下的 `.link/memory/` 目录。
2. 每条长期记忆是一个独立 `.md` 文件，并由 `memory/memory.md` 作为索引文件维护文件位置、名称和描述。
3. 系统启动/构建提示词时先加载 `memory/memory.md`，将候选记忆索引交给大模型判断哪些记忆有用，再按需读取/注入记忆内容或摘要到系统提示词。
4. 当前会话记忆使用 `session_id` 标识；同一会话只维护同一个 `.md` 文件，每轮结束后更新该文件，不生成多个会话记忆文件。
5. 保持现有“每轮结束异步更新记忆、失败不影响用户响应”的运行特性。
6. 所有记忆只对当前项目下的会话有效；打开其他项目时只能读取该项目自己的 `.link/memory/`，不得复用当前项目记忆。

## 2. 需求拆解

### 2.1 长期记忆文件结构

工作目录下新增 `.link/memory/`：

```text
<cwd>/.link/memory/
  memory.md
  language-preference.md
  project-style.md
  sessions/
    <session_id>.md
```

其中 `memory/memory.md` 是长期记忆索引，示例：

```markdown
# Memory Index

- [语言偏好：仅使用 Go](language-preference.md) — 项目必须使用 Go 语言，不能用其他语言
```

长期记忆文件使用 Markdown + YAML front matter：

```markdown
---
name: language-preference
description: 用户语言偏好：项目必须使用 Go 语言
metadata:
  node_type: memory
  type: user_preference
  project: nano_cc
  originSessionId: 041581e7-c3e7-46c8-afe7-7cdcc671e80e
---

用户在 nano_cc 项目中只使用 Go 语言，不能使用其他编程语言。

**Why:** 用户的技术栈和偏好明确限定为 Go。

**How to apply:** 所有代码实现、脚本、工具必须用 Go 编写。
```

说明：

- `metadata.node_type` 固定为 `memory`。
- `metadata.type` 与改造后的 `storage.Memory.Type` 对齐，可为 `user_preference`、`episodic_memory`、`project_fact`。
- `metadata.project` 记录当前项目名，用于审计和提示词说明；实际隔离以 `<cwd>/.link/memory` 的项目目录边界为准。
- `originSessionId` 记录生成该长期记忆的会话来源。
- 索引中的链接路径相对 `.link/memory/` 目录，禁止越界到工作目录外。

### 2.2 当前会话记忆文件结构

当前会话记忆写入：

```text
<cwd>/.link/memory/sessions/<session_id>.md
```

文件内容同样使用 front matter，但 `node_type` 使用 `session_memory`：

```markdown
---
name: current-session
description: 当前会话主干信息
metadata:
  node_type: session_memory
  project: nano_cc
  session_id: 041581e7-c3e7-46c8-afe7-7cdcc671e80e
  originSessionId: 041581e7-c3e7-46c8-afe7-7cdcc671e80e
---

## 当前目标

- 记录本场 TUI 会话中用户当前要完成的主要任务。

## 关键决策

- 记录本场会话已经确认的实现方案、约束和取舍。

## 已完成

- 记录本场会话已经产出的代码、文档或配置变更。

## 待办

- 记录下一轮继续工作时需要承接的事项。
```

同一个 `session_id` 在每轮会话结束后覆盖更新同一个文件，不拆成多个 `.md`。

## 3. 方案选择

### 方案 A：直接把现有 Store 方法改成 Markdown 文件读写（推荐）

在 `internal/local` 中实现一个 Markdown 记忆存储层，让本地 `Store` 的记忆接口真正读写 `.link/memory/` 文件；运行时 `runtime.memory.go` 和 `runtime.conversation_memory.go` 尽量复用现有抽取、选择、合并逻辑。

优点：改动集中，复用现有运行时接口、测试和异步调度逻辑；TUI 与未来可能存在的数据库 Store 保持接口隔离。缺点：需要在本地 Store 内补齐索引解析、Markdown 序列化和原子写入。

### 方案 B：绕过 Store，在 runtime 中直接读写 Markdown

让 `runtime.selectRelevantMemories`、`extractMemories`、`updateConversationMemory` 直接操作文件系统。

优点：实现路径短。缺点：runtime 与本地文件布局强耦合，不利于继续保留数据库 Store；测试替身和服务端逻辑也会被污染。

### 方案 C：新增独立 MemoryRepository，并让 Store 组合它

新增 `internal/memory` 或 `internal/local/memory` 仓储，Store 委托它实现所有记忆接口。

优点：边界最清晰，便于独立测试。缺点：比方案 A 多一层类型与构造注入，短期改动稍大。

推荐采用 **方案 C 的边界 + 方案 A 的集成方式**：新增本地 Markdown 记忆仓储，但仍通过 `local.Store` 暴露现有 `conversationStore` 接口。这样 runtime 基本不感知文件系统，TUI 本地实现也具备清晰单测边界。

## 4. 总体架构

### 4.1 新增组件

新增 `internal/local/memory_store.go`，提供 `MarkdownMemoryStore`：

- `rootDir`: `<WorkspaceRoot>/.link/memory`
- `indexPath`: `<WorkspaceRoot>/.link/memory/memory.md`
- `sessionsDir`: `<WorkspaceRoot>/.link/memory/sessions`

核心能力：

1. `EnsureLayout()`：创建 `.link/memory/`、`.link/memory/sessions/` 和空 `memory.md`。
2. `ListMemoryIndex()`：解析 `memory.md`，得到候选记忆的文件名、展示名称、描述。
3. `ReadMemoryFile(path)`：读取并解析单个长期记忆 Markdown。
4. `WriteMemoryFile(memory)`：为新记忆生成稳定 slug 文件名，写入 front matter + body，并更新 `memory.md`。
5. `ReplaceMemories(type/user scope)`：用于现有合并逻辑，将某类记忆对应的文件和索引替换为模型整理后的完整列表。
6. `ReadSessionMemory(sessionID)` / `WriteSessionMemory(sessionID, items)`：读写 `sessions/<session_id>.md`。

### 4.2 Store 集成

调整 `internal/local.Store`：

- 增加 `memory *MarkdownMemoryStore` 字段。
- `NewStore` 增加一个接收 `workspaceRoot` 的构造路径，或新增 `NewStoreWithMemory(workspaceRoot string)`，由 `Bootstrap` 使用。
- 现有内存字段可继续用于测试，但 TUI Bootstrap 路径应优先初始化 Markdown 存储。

本地 Store 的记忆接口改为：

- `ListRelevantMemories(ctx, userID)`：读取 `memory.md` 索引并返回候选 `storage.Memory`，候选可先只含 `Name`、`Description`、`Type`、`Body` 的轻量内容；如果后续选择阶段需要完整正文，再按文件读取。
- `InsertMemory(ctx, storage.Memory)`：写入一个长期记忆 `.md`，并更新索引。
- `ListMemoriesByUserAndType` / `ListProjectFactMemories`：从索引和文件 front matter 过滤。
- `ReplaceMemoriesByUserAndType` / `ReplaceProjectFactMemories`：重写对应类型长期记忆文件与索引。
- `ListConversationMemories(ctx, conversationID)`：通过本地会话上下文或 Store 映射解析出 UUID `session_id`，读取 `.link/memory/sessions/<session_id>.md` 并解析为 `[]storage.ConversationMemory`。
- `ReplaceConversationMemories(ctx, conversationID, userID, items)`：通过相同映射覆盖写入 `.link/memory/sessions/<session_id>.md`，确保同一会话只更新同一个文件。

## 5. 提示词与加载流程

### 5.1 启动时启用记忆

`Bootstrap` 当前显式关闭记忆。改造后应：

- `runtimeService.EnableMemory = true`
- 本地用户 `MemoryEnabled = true`
- 初始化 `MarkdownMemoryStore` 并确保 `.link/memory/` 布局存在。

### 5.2 候选记忆选择

当前 `selectRelevantMemories` 会把所有候选 `storage.Memory` 交给模型选择索引。改造后保持该模式，但候选来源变成 `memory.md`：

1. `ListRelevantMemories` 读取 `memory.md`。
2. 将每条索引项渲染为候选列表：`[i] (type) name: description`。
3. 模型返回有用的候选索引。
4. Store 或 runtime 根据选中结果读取对应 `.md` 正文。
5. 注入系统提示词时不要只注入索引描述；应注入选中记忆的 `name`、`description` 和正文，以满足“按需修改提示词”。
6. 注入段落必须明确这些记忆只适用于当前项目；模型不得把当前项目记忆迁移到其他项目会话中使用。

为降低改动量，可以先把 `ListRelevantMemories` 返回的 `Body` 填成 Markdown 正文，再调整 `renderMemorySection` 允许注入 body。注入格式建议：

```markdown
### 当前项目记忆（仅适用于 nano_cc 项目）

以下记忆只能用于当前项目；不要迁移到其他项目会话。

#### 用户偏好：语言偏好：仅使用 Go

项目必须使用 Go 语言，不能用其他语言。

用户在 nano_cc 项目中只使用 Go 语言，不能使用 Python、Node.js、Rust 等其他语言。

#### 项目事实：TUI 本地记忆目录

本项目记忆存储在当前工作目录的 memory/ 下。

启动其他项目时不能读取 nano_cc/memory 下的记忆。
```

### 5.3 记忆抽取/选择提示词调整

现有 `runtime.memory.go` 的记忆提示词需要同步调整：第三类长期记忆改成“项目事实记忆 / project_fact”。

长期记忆抽取提示词应表达为：

```text
You are a project-scoped long-term memory extraction engine for a personal assistant named Link.
All memories are valid ONLY for the current project. Do not create memories that should be reused in other projects.

From the dialogue, extract durable memories worth keeping for the current project. Use the "type" field for three kinds:
- "episodic_memory": a concrete event/experience in this project session. Preserve factual integrity and temporal order.
- "user_preference": a stable user preference, constraint, or recurring habit that applies to this project.
- "project_fact": a reusable fact about the current project, such as architecture, commands, conventions, dependencies, known constraints, or implementation decisions. It is NOT general world knowledge and must NOT be treated as valid outside this project.

Rules:
- Only extract NEW information not already covered by "Existing memories".
- Do not store one-off, trivial, or sensitive private data.
- Do not extract facts about other projects.
- "name": short title (<=80 chars). "description": one-sentence gist (<=300 chars). "body": supporting detail (<=2000 chars).
- Output ONLY a JSON array: [{"name","type","description","body"}].
- If nothing new or everything is already covered, output exactly [].
```

记忆选择提示词应补充当前项目作用域：

```text
You are a project-scoped memory retrieval engine. Given the current project, current conversation context, and a numbered list of candidate memories from this project's memory index, select the ones RELEVANT and USEFUL for answering the user right now.
- Candidate memories are valid only for the current project.
- Select at most 10.
- Prefer specific, on-topic memories; ignore unrelated ones.
- Output ONLY a JSON array of the selected memory indices, e.g. [0,3,7]. If none, output [].
```

合并提示词也应使用 `project_fact` 类型标签，并强调只整理当前项目的项目事实，不生成跨项目知识。

### 5.4 记忆抽取与索引维护

每轮结束后仍由 `scheduleMemoryWork` 异步调用：

1. `extractMemories` 抽取长期记忆。
2. `InsertMemory` 将每条记忆写成单独 `.md`。
3. 写入成功后更新 `memory.md`。
4. 达到阈值时沿用现有合并/裁剪逻辑，但落盘为“替换对应类型文件集合 + 重建索引”。

索引更新必须保持确定性排序，建议按：`type`、`name`、`file_path` 排序，减少无意义 diff。

## 6. session_id 与会话锁设计

代码中目前没有独立 `session_id` 字段，TUI 启动时每次创建一个 `storage.Conversation{ID: idgen.New("conv")}`。改造后应 **为每个 TUI 会话随机生成 UUID 格式的 `session_id`**，并将其作为当前会话记忆与会话锁的稳定标识。

设计要求：

- `session_id` 采用随机 UUID 生成，例如 `041581e7-c3e7-46c8-afe7-7cdcc671e80e`。
- 当前会话记忆文件固定为 `.link/memory/sessions/<session_id>.md`。
- 同一会话每轮结束后覆盖更新该文件，不因为轮次变化生成多个 `.md`。
- `conversation.ID` 可继续作为现有 runtime 历史、模型历史、事件流等内部会话 ID；但当前会话记忆落盘和记忆收尾锁使用显式 `session_id`。
- `session_id` 需要随 `storage.Conversation` 或 TUI 本地会话上下文传入 runtime；若不修改公共结构体，也可在本地 Store 维护 `conversation.ID -> session_id` 的映射，但落盘文件名必须使用 UUID `session_id`。

### 6.1 会话锁

仅对“当前会话记忆/收尾工作”进行加锁，不引入长期记忆文件级锁或跨项目全局锁。

会话锁 key 规则：

```text
<project_name>:<session_id>
```

其中：

- `project_name` 来自当前工作目录名，例如 `/Users/bytedance/golang_pro/nano_cc` 对应 `nano_cc`；若目录名为空，则使用清洗后的 workspace root 兜底。
- `session_id` 为上述随机 UUID。
- key 中的项目名和 session_id 均做安全清洗，仅保留 `[A-Za-z0-9._-]`，其他字符替换为 `-`。

该锁用于保证同一项目同一 TUI 会话的收尾更新串行执行，避免上一轮记忆更新尚未完成时下一轮覆盖同一个 `sessions/<session_id>.md`。不同项目或不同 `session_id` 的会话互不阻塞。

文件名安全规则：

- `session_id` 必须是 UUID；落盘前仍按 `[A-Za-z0-9._-]` 白名单清洗。
- 长期记忆 slug 同样做字符白名单与重名后缀处理。
- 索引中的相对路径必须 `filepath.Clean` 后仍位于 `.link/memory/` 下。

## 7. 文件格式与解析策略

### 7.1 Front matter

为避免新增重量级依赖，front matter 可用简单解析器处理：

- 文件以 `---\n` 开头时，读取到下一个单独一行 `---`。
- 支持当前需求中的 `name`、`description`、`metadata.node_type`、`metadata.type`、`metadata.project`、`metadata.originSessionId`、`metadata.session_id`。
- 写文件时使用固定格式输出；读历史/手写文件时容忍未知字段。

### 7.2 memory.md 索引解析

索引以 Markdown 列表为主协议：

```markdown
- [显示名](relative-file.md) — 描述
```

解析规则：

- 只解析符合该格式的列表项；其他说明文字忽略。
- 描述分隔符支持 `—`，兼容 `-` 作为降级。
- 文件路径必须是相对路径且不可包含 `..` 越界。

### 7.3 原子写入

所有 `.md` 写入采用：

1. 写入同目录临时文件。
2. `fsync` 或至少 `Close` 后 `os.Rename` 覆盖目标。
3. 索引最后写入，避免索引指向未写完文件。

长期记忆与索引写入采用进程内 `sync.Mutex` 串行化；当前会话记忆写入额外受“项目名 + session_id”会话锁保护。

## 8. 错误处理

- 读取 `memory.md` 不存在：创建空索引，返回空候选，不报错阻断启动。
- 某条索引指向文件不存在或解析失败：跳过该条并记录 warning。
- 写长期记忆失败：记录 warning，继续处理其他记忆，不影响用户响应。
- 写会话记忆失败：记录 warning，不影响用户响应。
- 模型选择记忆失败：保持现有行为，返回空 memory section。

## 9. 测试计划

### 9.1 单元测试

新增或更新以下测试：

1. Markdown front matter 解析与序列化。
2. `memory.md` 索引解析：正常项、空索引、非法路径、缺失文件。
3. `InsertMemory`：写入单个 `.md`，并更新索引。
4. `ListRelevantMemories`：只从当前项目的 `memory.md` 加载候选，必要时包含 body，不读取其他项目目录。
5. `ReplaceConversationMemories`：随机 UUID `session_id` 多次更新只覆盖同一个 `sessions/<session_id>.md`，会话锁 key 为 `项目名 + session_id`。
6. `selectRelevantMemories`：模型选中索引后注入正文而不只是 name/description，注入段落明确“仅当前项目有效”。
7. 记忆抽取/选择/合并提示词：只包含 `user_preference`、`episodic_memory`、`project_fact` 三类，并拒绝跨项目复用。
8. Bootstrap：TUI 本地模式记忆开关启用，当前项目 `.link/memory/` 布局被创建。

### 9.2 集成验证

运行：

```bash
go test ./internal/local ./internal/agent/runtime ./internal/config ./internal/cli
```

如改动影响公共接口，再运行：

```bash
go test ./...
```

## 10. 实施步骤

1. 新增 Markdown 记忆仓储与文件格式工具，覆盖解析、序列化、原子写入。
2. 调整 `local.Store` 初始化与记忆接口实现，让本地 TUI 使用 `<cwd>/.link/memory`。
3. 修改 `Bootstrap`，启用本地记忆能力并创建目录布局。
4. 修改 `runtime.memory.go` 中的抽取、选择、合并提示词：第三类长期记忆改为 `project_fact` 项目事实记忆，并强调记忆仅当前项目有效。
5. 调整记忆注入渲染逻辑，使模型从 `memory.md` 索引选择后能将选中 `.md` 的有效内容注入系统提示词，并在注入段落注明“仅当前项目有效”。
6. 调整当前会话记忆读写，为 TUI 会话生成随机 UUID `session_id`，覆盖写入 `.link/memory/sessions/<session_id>.md`，并使用 `项目名 + session_id` 作为会话锁 key。
7. 补齐单元测试与必要的集成测试。
8. 运行测试并修复失败。

## 11. 非目标

- 不实现跨工作目录共享记忆；本次记忆只存放在当前工作目录的 `.link/memory/` 下。
- 不允许当前项目记忆影响其他项目会话；其他项目只能读取各自工作目录下的 `.link/memory/`。
- 不实现长期记忆文件级锁或跨项目全局锁；本次仅对当前会话使用 `项目名 + session_id` 维度加锁。
- 不做历史数据库记忆迁移；TUI 本地 Store 当前没有真实数据库落盘，直接启用 Markdown 存储。
- 不新增用户配置项控制记忆目录；目录固定为 `<cwd>/.link/memory`，符合需求。

## 12. 审阅关注点

请重点确认以下设计决策是否符合预期：

1. 是否接受 `session_id` 使用随机 UUID，并与 `conversation.ID` 分离。
2. 长期记忆 `.md` 是否需要注入完整 body；本设计建议模型先从索引选择，再注入选中记忆正文。
3. `.link/memory/sessions/<session_id>.md` 是否可作为当前会话记忆文件路径。
4. 是否接受本次仅对当前会话使用 `项目名 + session_id` 加锁，不做长期记忆文件级锁或跨项目全局锁。
5. 是否接受第三类长期记忆命名为 `project_fact`，用于记录当前项目事实且不跨项目复用。
