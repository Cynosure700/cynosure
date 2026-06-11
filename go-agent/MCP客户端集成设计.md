# MCP 客户端集成设计文档

## 1. 需求与目标

为 go-agent（nano_cc）新增 MCP（Model Context Protocol）配置能力：

- 用户可在前端自行配置 MCP 服务器，支持两种连接方式：
  - **sse**：基于 HTTP Server-Sent Events 的远程连接（配置 url + headers）。
  - **streamable**：Streamable HTTP 传输（配置 url + headers）。
- 后台服务为每个启用的 MCP 配置创建 MCP client，建立与 MCP 服务器的连接。
- 实现 MCP 标准能力：**工具发现（tools/list）**、**工具调用（tools/call）**、**结果返回**。
- 发现到的 MCP 工具自动注入到对话的工具列表中，LLM 可像调用内置工具一样调用它们。

### 核心约束（注意事项）

> **不影响现有功能。** 所有改动均为增量：新增表、新增 repo、新增 handler、新增 manager，对现有内置工具系统、对话主流程、压缩/记忆等机制做**最小侵入**改动。

---

## 2. 现状架构分析

### 2.1 工具系统（关键契约）

- 工具定义统一使用 `openai.Tool`（`internal/tools/definitions.go`），全局静态来源是 `agenttools.AllToolDefs`。
- 工具执行三层：`ToolRegistry.Execute`（注入 ctx）→ `agenttools.Dispatch`（按名分发）→ `Handlers[name]`。
  - 特例旁路：`todo_write`（结构化返回）、`spawn_subagent`（子 agent 流程）不走标准 `Handlers`。
- `ToolRegistry`（`internal/web/runtime/tool_registry.go`）是 **Service 级单例**，启动时由 `NewToolRegistry(cfg)` 构建，工具定义编译期固定。
- 对话主循环 `RespondToConversation`（`internal/web/runtime/conversation_flow.go:29`）中：
  - 第 100 行 `Tools: s.Tools.Definitions()` 把工具定义传给 LLM。
  - 第 94 行用 `s.Tools.Definitions()` 估算 token。
  - 第 92 行 `maybeAppendTodoWriteReminder(state, s.Tools, ...)`。
  - 第 146 行 `s.executeToolCall(...)` 执行工具调用。

> **关键洞察**：主循环与 LLM client 只依赖两个稳定契约——`[]openai.Tool`（定义）与 `ToolExecutionResult`/outcome（结果）。MCP 工具只要产出这两种结构，即可零侵入接入主循环。

### 2.2 配置与存储

- 静态配置 `config.json` 为**只读、全局、需重启**，**不适合**承载用户级、可增删改的 MCP 配置。
- 用户级动态数据的标准模式是 **MySQL 表 + Repo + CRUD API + 前端面板**，`skills` 是完整可复刻的模板：
  - 迁移：`storage/migrations.go` 的 `ensureXxxTable`（幂等 `CREATE TABLE IF NOT EXISTS`）。
  - 模型：`storage/models.go` 的 `Skill` 结构体。
  - Repo：`storage/skills_repo.go`（手写 SQL，无 ORM）。
  - 接口：`app/server.go` 的 `serverStore` 接口。
  - Handler：`app/skill_handlers.go`（集合 + 单资源双 handler，含归属校验 `UserID != user.ID → 403`）。
  - 路由：`app/routes.go`（`/api/skills` + `/api/skills/`，`AuthenticateRequest` 包裹）。
  - 前端：`web/src/api.ts` + `web/src/App.tsx`（侧边栏能力面板的列表卡 + 表单）。

### 2.3 数据库

MySQL（主存储）+ Redis（锁/缓存）+ ES（检索）。MCP 配置只需用 MySQL。

---

## 3. 总体设计

```
┌─────────────┐   配置 CRUD    ┌──────────────────┐
│  前端 App    │ ────────────> │ /api/mcp-servers  │  (mcp_handlers.go)
│ MCP 配置面板 │ <──────────── │   (CRUD + test)   │
└─────────────┘                └────────┬─────────┘
                                         │ 读写
                                   ┌─────▼──────┐
                                   │ mcp_servers │  (MySQL 表)
                                   └─────┬──────┘
                                         │ 加载启用配置
                                ┌────────▼─────────┐
                                │   MCPManager      │  (新增: internal/web/mcp)
                                │ - 按用户管理 client│
                                │ - 连接/发现/调用   │
                                │ - 工具定义缓存     │
                                └────────┬─────────┘
                  Definitions / Call     │
        ┌────────────────────────────────▼────────────────────┐
        │  RespondToConversation 对话主循环 (conversation_flow) │
        │  - 合并 静态工具 + 该用户 MCP 工具 → 传给 LLM         │
        │  - executeToolCall: MCP 工具 → 转发 MCPManager        │
        └──────────────────────────────────────────────────────┘
```

### 3.1 依赖选型

引入官方 Go SDK：`github.com/modelcontextprotocol/go-sdk/mcp`（Apache-2.0，Google 联合维护）。本项目使用其中两种 transport：

- `SSEClientTransport` → sse
- `StreamableClientTransport` → streamable HTTP

> SDK 同时内置 `CommandTransport`（stdio，启动子进程），本项目出于安全与部署考虑不启用该选型。

> go.mod 当前为 `go 1.26.1`，满足官方 SDK 的 Go 版本要求。新增依赖为纯增量，不影响现有依赖。

### 3.2 数据模型

新增 MySQL 表 `mcp_servers`（用户级，外键级联删除，仿 `skills`）：

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | VARCHAR(64) PK | `newID("mcp")` |
| `user_id` | VARCHAR(64) | 外键 → users(id) ON DELETE CASCADE |
| `name` | VARCHAR(255) | 显示名称（用于工具命名空间前缀） |
| `transport` | VARCHAR(32) | `sse` / `streamable` |
| `command` | VARCHAR(1024) | 保留列（历史 stdio 字段，当前不使用） |
| `args` | JSON/TEXT | 保留列（历史 stdio 字段，当前不使用） |
| `env` | JSON/TEXT | 保留列（历史 stdio 字段，当前不使用） |
| `url` | VARCHAR(1024) | sse/streamable：服务地址 |
| `headers` | JSON/TEXT | sse/streamable：自定义请求头（JSON 序列化） |
| `enabled` | TINYINT(1) | 是否启用 |
| `created_at` / `updated_at` | DATETIME(6) | 时间戳 |

- 唯一约束：`UNIQUE KEY uniq_mcp_user_name (user_id, name)`，避免同一用户重名导致工具命名冲突。
- 索引：`KEY idx_mcp_user_id (user_id)`。

对应 `storage/models.go` 新增 `MCPServer` 结构体（带 json tag，args/env/headers 在 Go 侧以 `[]string` / `map[string]string` 表示，repo 层做 JSON 编解码）。

### 3.3 工具命名空间

MCP 工具名统一加前缀：`mcp__<sanitized_server_name>__<tool_name>`（与 Claude Code 约定一致）。

- 避免与内置工具（bash/read_file…）及不同 server 间重名。
- 在 `executeToolCall` 中通过 `strings.HasPrefix(name, "mcp__")` 快速判定路由到 MCPManager。
- `<sanitized_server_name>` 对 name 做 `[a-zA-Z0-9_]` 归一化。

### 3.4 MCPManager（核心新增组件）

新增包 `internal/web/mcp`，提供 `Manager`：

```go
type Manager struct {
    store mcpStore          // 读取用户 mcp_servers 配置
    mu    sync.RWMutex
    // key: userID -> 已连接会话集合（按 server id 索引）
    sessions map[string]map[string]*serverSession
}

type serverSession struct {
    server   storage.MCPServer
    session  *mcp.ClientSession   // 官方 SDK 会话
    tools    []openai.Tool        // 已发现并转换后的工具定义（带 mcp__ 前缀）
    toolMap  map[string]string    // 前缀名 -> 原始 MCP 工具名
}
```

职责：

1. **连接管理**：
   - `EnsureUserSessions(ctx, userID)`：懒加载——读取该用户 `enabled=1` 的配置，对未连接/配置变更的 server 建立连接，对已删除/禁用的 server 关闭连接。
   - 根据 `transport` 字段创建对应 transport 并 `client.Connect(ctx, transport)`。
2. **工具发现**：连接成功后调用 `session.ListTools`，把每个 MCP tool 的 `InputSchema`（标准 JSON Schema）映射为 `openai.Tool.Function.Parameters`，名称加 `mcp__` 前缀。
3. **工具定义聚合**：`ToolsForUser(ctx, userID) []openai.Tool` 返回该用户所有 MCP 工具定义（供合并进 LLM 请求）。
4. **工具调用**：`CallTool(ctx, userID, prefixedName, args) (string, error)`：
   - 解析前缀定位 server 与原始工具名 → `session.CallTool` → 将返回的 content blocks（text/json 等）序列化为字符串（契合现有 `ExecResult.Output` 字符串契约）。
   - MCP 返回 `IsError` 时转为 error，由主循环现有的 `rejected` 分支处理。
5. **连接探活/测试**：`TestServer(ctx, server)` 供前端"测试连接"按钮调用（连接 + ListTools，返回工具数量或错误）。
6. **生命周期**：`Close()` 在服务退出时关闭全部会话；配置 CRUD 变更后调用使该用户缓存失效（下次对话重连）。

> **健壮性原则**：单个 MCP server 连接失败/发现失败**不阻断对话**，仅记录日志并跳过该 server 的工具（降级）。这与现有"工具白名单查不到则跳过"的容错风格一致。

### 3.5 接入对话主流程（最小侵入）

在 `runtime.Service` 新增字段 `MCP *mcp.Manager`，通过 `SetMCPManager()` 注入（与 `SetBuiltinSkills` 同风格，保持 `NewService` 签名不变）。

改动点（仅 `conversation_flow.go` / `tool_registry.go` 内的几处）：

1. 新增辅助方法 `toolDefinitionsForUser(ctx, userID) []openai.Tool`：
   - 返回 `s.Tools.Definitions()`；若 `s.MCP != nil`，先 `EnsureUserSessions` 再 `append` MCP 工具定义。
2. 主循环把 `s.Tools.Definitions()` 的 **3 处** 调用替换为 `defs := s.toolDefinitionsForUser(ctx, user.ID)` 并复用：
   - token 估算（:94）、`req.Tools`（:100）、最终用量估算（:124）。
   - `maybeAppendTodoWriteReminder` 改为接收已合并的定义判断（或保持仅判断内置 todo_write，不影响逻辑）。
3. `executeToolCall`（`tool_registry.go:191`）开头新增分支：

```go
if strings.HasPrefix(name, "mcp__") {
    if s.MCP == nil {
        return outcome{Status: "rejected", Result: "Error: MCP not enabled", Audit: audit}
    }
    out, err := s.MCP.CallTool(ctx, toolCtx.User.ID, name, rawArgs)
    if err != nil {
        return outcome{Status: "rejected", Result: fmt.Sprintf("Error: %v", err), Audit: audit}
    }
    return outcome{Status: "success", Result: out, Audit: audit}
}
```

> 结果写回、SSE 下发、持久化等全部复用现有 hook 链（`appendToolMessageHook` 等），无需改动。
> 子 agent 流程（`subagent.go`）本期**不注入** MCP 工具（保持子 agent 工具集最小化），后续可按需扩展。

### 3.6 HTTP API

新增 `app/mcp_handlers.go`，复刻 skills 双 handler 模式（全部 `AuthenticateRequest` 鉴权 + 归属校验）：

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/mcp-servers` | 列出当前用户的 MCP 配置 |
| POST | `/api/mcp-servers` | 创建配置（`normalizeTransport` 校验枚举 + 字段校验） |
| GET | `/api/mcp-servers/{id}` | 获取详情 |
| PUT | `/api/mcp-servers/{id}` | 全量更新 |
| PATCH | `/api/mcp-servers/{id}` | 切换 enabled |
| DELETE | `/api/mcp-servers/{id}` | 删除 |
| POST | `/api/mcp-servers/{id}/test` | 测试连接（返回发现到的工具数/错误） |

- 校验：`transport` 必须为两枚举之一；sse/streamable 必填 `url`。
- 任何写操作成功后调用 `mcpManager.Invalidate(userID)` 使缓存失效。
- `serverStore` 接口（`app/server.go:19`）新增 5+1 个方法签名（CRUD + GetByID）。

### 3.7 存储层

新增 `storage/mcp_repo.go`（复刻 `skills_repo.go`）：

```go
func (s *Store) CreateMCPServer(ctx, m MCPServer) error
func (s *Store) ListMCPServersByUser(ctx, userID) ([]MCPServer, error)
func (s *Store) ListEnabledMCPServersByUser(ctx, userID) ([]MCPServer, error)
func (s *Store) GetMCPServerByID(ctx, id) (MCPServer, error)
func (s *Store) UpdateMCPServer(ctx, m MCPServer) error
func (s *Store) DeleteMCPServer(ctx, id) error
```

迁移：`migrations.go` 新增 `ensureMCPServersTable(ctx)`，并在 `RunMigrations` 末尾追加调用（**不改动 `001_init.sql`**，保持幂等增量迁移风格）。

### 3.8 前端

`web/src/api.ts`：新增 `MCPServer` 类型与 API 方法（`listMCPServers / createMCPServer / updateMCPServer / patchMCPServerEnabled / deleteMCPServer / testMCPServer`），复刻 skills 段。

`web/src/App.tsx`：
- `SidePanel` 类型增加 `"mcp"`，在能力面板旁新增"MCP 服务器"入口按钮。
- 复刻 skill 面板的「列表卡片 + 创建/编辑表单」：
  - `transport` 用 `<select>`（sse/streamable）。
  - 按 `transport` 渲染 url/headers 字段。
  - 每张卡：启用/禁用、编辑、删除、测试连接（显示发现的工具数或错误）。
- `refreshAll` 中并行加载 MCP 配置列表。

---

## 4. 改动文件清单

| 层 | 文件 | 动作 |
|---|---|---|
| 依赖 | `go-agent/go.mod` / `go.sum` | 新增 `modelcontextprotocol/go-sdk` |
| 迁移 | `internal/web/storage/migrations.go` | 新增 `ensureMCPServersTable` + 调用 |
| 模型 | `internal/web/storage/models.go` | 新增 `MCPServer` 结构体 |
| Repo | `internal/web/storage/mcp_repo.go` ✨新建 | CRUD |
| Manager | `internal/web/mcp/manager.go` ✨新建 | 连接/发现/调用/缓存 |
| Manager | `internal/web/mcp/transport.go` ✨新建 | transport 构建（sse/streamable）+ content 序列化 |
| 接口 | `internal/web/app/server.go` | `serverStore` 增方法；装配 Manager |
| Handler | `internal/web/app/mcp_handlers.go` ✨新建 | CRUD + test |
| 校验 | `internal/web/app/mcp_handlers.go` | `normalizeTransport` 等 |
| 路由 | `internal/web/app/routes.go` | 注册 2 条路由 |
| Runtime | `internal/web/runtime/runtime.go` | `Service` 增 `MCP` 字段 + `SetMCPManager` |
| Runtime | `internal/web/runtime/conversation_flow.go` | `toolDefinitionsForUser` + 替换 3 处定义调用 |
| Runtime | `internal/web/runtime/tool_registry.go` | `executeToolCall` 增 MCP 分支 |
| 前端 | `web/src/api.ts` | 类型 + API 方法 |
| 前端 | `web/src/App.tsx` | MCP 配置面板 |

> ✨ = 新建文件，其余为增量编辑。

---

## 5. 关键设计决策

1. **配置存 MySQL 而非 config.json**：MCP 配置是用户级、动态、可增删改数据，与 skills 同构，不应进静态全局配置。
2. **工具命名空间 `mcp__server__tool`**：避免冲突、便于路由判定、对齐业界约定。
3. **懒加载 + 按用户缓存连接**：首轮对话时按需连接，配置变更使缓存失效后重连；避免启动时为所有用户建连。
4. **降级容错**：单个 MCP server 故障不影响对话与其他工具，仅跳过其工具。
5. **最小侵入**：复用 `openai.Tool` 与现有 hook/结果写回链，主循环仅替换"工具定义来源"和新增一个执行分支。
6. **本期范围**：MCP 仅提供 **Tools**（工具）能力；Resources/Prompts 暂不接入（可后续扩展）。子 agent 暂不注入 MCP 工具。

---

## 6. 不影响现有功能的保障

- 数据库：仅新增表，幂等迁移，不触碰既有表与 `001_init.sql`。
- 工具系统：`AllToolDefs`、`Handlers`、`Dispatch` 完全不动；MCP 工具走独立路由分支。
- 主流程：未配置 MCP 时 `s.MCP == nil` 或用户无启用配置，`toolDefinitionsForUser` 返回值与 `s.Tools.Definitions()` 完全一致，行为零变化。
- 压缩/记忆/SSE/持久化：复用现有 hook，无改动。

---

## 7. 风险与缓解

| 风险 | 缓解 |
|---|---|
| 远程 MCP 连接阻塞对话 | 连接/发现/调用均设置 context 超时；失败降级 |
| MCP 工具过多撑大上下文 | token 估算已纳入合并后的定义；可后续加每用户工具数上限 |
| SDK 版本/API 变动 | 锁定具体版本；transport 构建集中在 `transport.go` 便于适配 |
| 工具名冲突 | `mcp__` 前缀 + server 名唯一约束 |

---

## 8. 验证计划

1. `go build ./...` 通过；新增 repo/manager 编译无误。
2. 迁移在已有库上幂等执行，`mcp_servers` 表创建成功。
3. 单测：`mcp_repo` CRUD；`manager` 工具名前缀/content 序列化；`normalizeTransport` 枚举校验。
4. 端到端：
   - 前端创建一个 sse/streamable 配置（填写远程 MCP 服务地址 url）→ 测试连接显示工具数。
   - 发起对话，确认 LLM 可调用 `mcp__*` 工具并返回结果。
   - 删除/禁用配置后，工具从对话中消失。
5. 回归：未配置 MCP 时，对话与内置工具行为与改动前一致。
