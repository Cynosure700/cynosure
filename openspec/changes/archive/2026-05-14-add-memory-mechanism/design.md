## Context

go-agent 当前是无状态的：每次启动都是全新会话，不保留任何跨会话信息。Claude Code 通过 `CLAUDE.md`（项目级）和 `~/.claude/CLAUDE.md`（用户级）实现了持久化记忆，同时通过会话级记忆让 Agent 在单次对话中记住关键信息。go-agent 已有 `AGENTS.md` 文件存在于项目根目录，但未被加载使用。

现有架构中，system prompt 在 `repl.go` 中构建，仅包含工作目录路径和技能描述。记忆机制需要在此处注入持久化记忆，同时在 `loop.go` 的每轮对话中注入会话级记忆。

## Goals / Non-Goals

**Goals:**
- 启动时自动加载项目级 `AGENTS.md` 和用户级 `~/.link/AGENTS.md`
- 持久化记忆以结构化方式注入 system prompt
- 提供 `update_memory` 工具，支持两种模式：
  - 持久化模式：追加/替换项目级 `AGENTS.md` 文件
  - 会话模式：写入会话级记忆，仅当前会话有效
- 会话级记忆在每轮对话中作为上下文注入，会话结束后自动清除
- 记忆文件为纯 Markdown，用户可直接编辑

**Non-Goals:**
- 不支持向量化记忆或 RAG 检索（保持简单）
- 不支持记忆的自动总结/压缩（由 Agent 通过工具自行管理）
- 不修改子代理的 system prompt（子代理不需要记忆上下文）
- 不引入新的外部依赖

## Decisions

### Decision 1: 两层持久化记忆 + 一层会话记忆

**选择**：`AGENTS.md`（项目根目录）+ `~/.link/AGENTS.md`（用户目录）+ 会话级内存记忆

**理由**：
- 项目级记忆随代码版本控制，团队共享
- 用户级记忆跨项目生效，存放个人偏好（如语言、代码风格）
- 会话级记忆存放本次对话中产生的临时知识，不需要持久化
- 项目级优先：当两者冲突时，项目级覆盖用户级

**替代方案**：全部持久化到文件 — 会话级临时信息（如"用户刚才说要重构 X 模块"）不应污染持久化记忆文件。

### Decision 2: 记忆注入位置

**选择**：持久化记忆注入 system prompt 末尾，会话记忆作为独立 user 消息注入每轮对话

```
<!-- system prompt 末尾 -->
<project_memory>
{AGENTS.md 内容}
</project_memory>

<user_memory>
{~/.link/AGENTS.md 内容}
</user_memory>

<!-- 每轮对话中，在用户消息之后追加 -->
<session_memory>
{会话记忆内容}
</session_memory>
```

**理由**：
- 持久化记忆放在 system prompt 中，作为全局行为约束
- 会话记忆放在每轮对话中，确保模型始终能"看到"最新的会话上下文
- XML 标签格式与现有 skill 加载格式（`<skill name="...">`）一致

**替代方案**：全部放 system prompt — 会话记忆频繁变化，每次都重建 system prompt 成本高。

### Decision 3: update_memory 工具设计

**选择**：提供 `update_memory` 工具，接受 `scope`（session/project）、`action`（append/replace）和 `content` 参数

- `scope: "session"` — 写入会话级记忆（内存中，会话结束清除）
- `scope: "project"` + `action: "append"` — 追加到 `AGENTS.md`
- `scope: "project"` + `action: "replace"` — 覆盖 `AGENTS.md`

**理由**：
- 一个工具两种模式，减少工具数量，降低模型选择负担
- 会话记忆默认行为：追加到已有会话记忆中
- 仅项目级记忆可持久化到文件，用户级记忆由用户手动维护

**替代方案**：两个独立工具（`update_memory` + `remember`）— 增加模型选择复杂度。

### Decision 4: 会话记忆存储结构

**选择**：使用简单的字符串拼接，每条记忆用换行分隔

```go
type SessionMemory struct {
    mu      sync.RWMutex
    entries []string
}
```

**理由**：
- 简单直接，无需复杂的数据结构
- 模型自行管理记忆内容的格式和去重
- 线程安全（Agent 循环和工具调用可能并发）

**替代方案**：Key-Value 存储 — 过于复杂，模型难以通过工具正确使用。

### Decision 5: 记忆文件缺失处理

**选择**：文件不存在时静默跳过，不报错

**理由**：
- 记忆是可选的增强功能，不应阻塞启动
- 用户可随时创建文件，下次启动自动生效

## Risks / Trade-offs

- **[风险] 记忆文件过大导致 token 消耗过高** → 建议用户保持记忆文件精简（< 2000 字），Agent 也可通过 `update_memory` 自行整理
- **[风险] Agent 可能写入错误或冗余的记忆** → `update_memory` 的 `replace` 操作会覆盖全部内容，用户可随时手动修正
- **[风险] 会话记忆无限增长导致上下文爆炸** → 会话记忆在 context compact 时一并被压缩，不会无限增长
- **[权衡] 子代理不加载记忆** → 子代理专注于单一任务，记忆上下文可能干扰其判断。如后续需要，可选择性注入