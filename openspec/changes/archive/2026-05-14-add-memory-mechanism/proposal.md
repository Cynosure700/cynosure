## Why

当前 go-agent 每次启动都是"空白大脑"，无法记住用户偏好、项目约定或历史决策。Claude Code 通过 `CLAUDE.md` / `AGENTS.md` 文件实现了持久化记忆，使 Agent 在多次会话间保持一致的上下文。go-agent 需要同样的能力，让用户可以通过 Markdown 文件定义项目级和个人级的持久化指令，Agent 启动时自动加载并注入到 system prompt 中。

此外，Agent 在单次会话中也需要"工作记忆"——在对话过程中记住用户提到的关键信息（如"我喜欢用蛇形命名"、"这个项目用 Redis 做缓存"），并在后续对话中持续引用。Claude Code 的 `update_memory` 工具同时支持持久化和会话级记忆，go-agent 也应具备此能力。

## What Changes

- 新增记忆文件加载机制：启动时自动扫描并加载 `AGENTS.md`（项目级）和 `~/.link/AGENTS.md`（用户级）两个记忆文件
- 记忆内容注入 system prompt：加载的记忆内容以结构化方式追加到 system prompt 中，作为 Agent 的持久化行为指引
- 新增 `update_memory` 工具：Agent 可在对话中主动更新项目级记忆文件，实现跨会话的知识积累
- 新增会话级记忆机制：Agent 可通过 `update_memory` 工具写入会话级记忆（仅当前会话有效），内容注入到后续每轮对话的上下文中
- 记忆文件优先级：项目级记忆覆盖用户级记忆中的冲突项，就近原则

## Capabilities

### New Capabilities
- `memory-loading`: 启动时自动加载项目级 (`AGENTS.md`) 和用户级 (`~/.link/AGENTS.md`) 记忆文件，解析并注入 system prompt
- `memory-update`: 提供 `update_memory` 工具，支持持久化记忆（写入文件）和会话级记忆（仅当前会话有效），允许 Agent 在对话中追加或修改记忆内容
- `session-memory`: 会话级记忆管理，Agent 写入的会话记忆在每轮对话中作为上下文注入，会话结束后自动清除

### Modified Capabilities
<!-- 无现有 spec 需要修改 -->

## Impact

- 新增文件：`internal/sessions/memory.go` — 持久化记忆加载与会话记忆管理
- 新增文件：`internal/tools/memory.go` — `update_memory` 工具处理器
- 修改文件：`internal/agent/repl.go` — system prompt 构建时注入记忆内容
- 修改文件：`internal/agent/loop.go` — 每轮对话注入会话记忆上下文
- 修改文件：`internal/tools/registry.go` — 注册 `update_memory` 工具
- 无外部依赖变更