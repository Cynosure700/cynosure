# 内置 MCP 连接设计文档

## 1. 背景与目标

当前 go-agent 已支持用户级 MCP 配置：用户通过接口维护自己的 MCP server，运行时由 `internal/web/mcp.Manager` 按 `userID -> serverID -> session` 缓存连接。对话准备工具定义时，`Service.toolDefinitionsForUser` 会调用 `EnsureUserSessions` 懒加载该用户启用的 MCP server，并将发现到的 MCP 工具合并进模型工具列表，见 `internal/web/runtime/tool_registry.go:74`。MCP 工具调用通过 `executeToolCall` 中的 `mcp__` 前缀分支路由到 `Manager.CallTool`，见 `internal/web/runtime/tool_registry.go:219`。

新需求已调整为：

1. 为系统新增**内置 MCP 连接**。
2. 内置 MCP 配置来自 `<app_home>/mcp_config.json`。
3. 内置 MCP **仅支持 stdio**，不支持 `sse` / `streamable`。
4. 服务启动时创建内置 MCP 连接。
5. 内置 MCP 工具对所有用户可用。
6. 内置 MCP 连接**不会被空闲回收**。
7. 如果内置 MCP 连接中断，需要自动重新建立连接。
8. 用户级 MCP 能力、非 MCP 工具、现有对话流程不受影响。

本阶段只修改设计文档，待审阅确认后再进行代码改造。

## 2. 现状分析

### 2.1 启动流程

服务启动入口 `NewServer` 的关键流程在 `internal/web/app/server.go:52`：

1. `config.LoadWebConfig()` 加载 `config.json` 与环境变量，见 `internal/config/web_config.go:8`。
2. `config.EnsureAppLayout` / `config.ValidateAppLayout` 准备运行目录，见 `internal/web/app/server.go:57`。
3. 初始化存储、迁移、日志、内置 skills、base prompt。
4. 创建 runtime service。
5. 创建 MCP manager：`mcp.NewManager(store)`，见 `internal/web/app/server.go:92`。
6. 将 MCP manager 注入 runtime：`runtimeService.SetMCPManager(mcpManager)`，见 `internal/web/app/server.go:93`。

内置 MCP 的配置加载与启动连接适合放在 `mcp.NewManager(store)` 之后、`runtimeService.SetMCPManager(mcpManager)` 之前执行。这样 runtime 接收到的 manager 已经包含启动时尝试建立的内置连接。

### 2.2 配置文件路径

当前 `config.json` 固定通过 `configFilePath()` 读取，路径为当前工作目录下的 `config.json`，见 `internal/config/config.go:74` 与 `internal/config/config.go:87`。

`AppConfig.AppHome` 会解析为绝对路径，见 `internal/config/paths.go:25`。运行时资源目录均按 `app_home` 解析，例如 `system_prompt.md`、`skills`、`workspace`。

本设计将内置 MCP 配置文件固定为：

```text
<app_home>/mcp_config.json
```

理由：

- 与 `system_prompt.md`、`skills`、`workspace` 等运行时资源保持同一归属。
- 避免依赖进程启动时的当前工作目录。
- 与部署时“应用目录内放置配置文件”的方式一致。

### 2.3 MCP Manager 当前能力

`internal/web/mcp/manager.go` 当前提供：

- `EnsureUserSessions(ctx, userID)`：读取用户启用 MCP 配置，建立或复用连接，关闭用户禁用或删除的连接，见 `internal/web/mcp/manager.go:92`。
- `ToolsForUser(userID)`：返回某用户所有已连接 MCP 工具定义，见 `internal/web/mcp/manager.go:195`。
- `CallTool(ctx, userID, prefixedName, rawArgs)`：按带前缀的工具名找到对应 session 并调用 MCP tool，见 `internal/web/mcp/manager.go:215`。
- `Invalidate(userID)`：用户 MCP 配置变更后关闭并清理该用户连接，见 `internal/web/mcp/manager.go:258`。
- `TestServer(ctx, server)`：临时连接单个配置并发现工具，见 `internal/web/mcp/manager.go:268`。
- 用户级空闲回收：session 记录 `lastUsedAt` / `activeCalls`，后台每分钟回收空闲超过 10 分钟且无 in-flight 调用的连接，见 `internal/web/mcp/manager.go:318` 与 `internal/web/mcp/manager.go:331`。

本次改造需要保留用户级 MCP 的现有行为，但内置 MCP 不参与空闲回收。

### 2.4 stdio 支持现状

当前 `buildTransport` 只支持用户级远程 transport，见 `internal/web/mcp/transport.go:37`。但项目依赖的 `github.com/modelcontextprotocol/go-sdk/mcp` 已提供 stdio 客户端 transport：

```go
&mcp.CommandTransport{Command: exec.Command(command, args...)}
```

SDK 中 `CommandTransport` 会启动子进程，并通过 stdin/stdout 与 MCP server 通信。内置 MCP 可单独实现 stdio transport 构造，不需要改变用户级 `buildTransport` 的远程 transport 行为。

## 3. 设计原则

1. **内置与用户级隔离**：内置 MCP 不写入 `mcp_servers` 数据表，不出现在用户 MCP CRUD 列表中，不被用户修改或删除。
2. **仅内置支持 stdio**：stdio 只用于系统内置 MCP；用户级 MCP 不新增 stdio 能力。
3. **所有用户可见**：任意用户对话构建工具列表时，都合并内置 MCP 工具。
4. **服务启动预连接**：启动时读取 `mcp_config.json` 并尝试连接启用的内置 MCP server。
5. **不空闲回收**：内置 MCP 是系统级长连接，不参与 10 分钟空闲回收。
6. **中断后恢复**：内置 MCP 连接中断后，后续获取工具定义或调用工具时应重新建立连接。
7. **最小侵入**：复用 `serverSession`、工具发现、工具命名前缀、工具调用结果序列化等现有逻辑；新增 stdio 构造只服务于内置 MCP。
8. **失败隔离**：单个内置 MCP 连接失败不阻塞其他内置 MCP，也不阻塞用户级 MCP；配置文件格式错误应在启动时暴露。

## 4. 可选方案

### 方案 A：启动时把内置 MCP 写入每个用户的 `mcp_servers` 表

启动时读取 `mcp_config.json`，为每个用户插入对应 MCP 配置。

优点：

- 能复用现有用户级 `EnsureUserSessions`。

缺点：

- 会污染用户数据，用户列表中会出现系统内置配置。
- 新用户注册后还需要补写，逻辑复杂。
- 用户可能误删或修改内置配置。
- 用户级 MCP 当前不支持 stdio；为复用用户级逻辑而恢复用户 stdio 会扩大需求范围。
- “所有用户均可使用”的系统级语义被拆散到每个用户数据中。

结论：不推荐。

### 方案 B：Manager 内新增全局 builtin stdio session 池（推荐）

`Manager` 同时维护两类连接：

- 用户级连接：现有 `sessions map[userID]map[serverID]*serverSession`，保持现有远程 transport 与空闲回收行为。
- 内置连接：新增 `builtinSessions map[serverID]*serverSession` 与 `builtinServers map[serverID]storage.MCPServer`，只使用 stdio，且不参与空闲回收。

启动时读取 `mcp_config.json`，调用 Manager 的内置 MCP 初始化方法建立 stdio 连接。每次对话构建工具定义时，确保内置连接仍可用；如果缺失或中断，则重连。最后合并内置工具与用户工具。

优点：

- 系统级语义清晰。
- 不污染数据库与用户 CRUD。
- stdio 能力只暴露给内置 MCP，不影响用户级安全边界。
- 可复用现有工具发现、工具调用、结果序列化逻辑。
- 能实现“不回收但中断重连”的独立生命周期。

缺点：

- `Manager` 需要新增 builtin 状态、stdio transport 构造和重连逻辑。
- `ToolsForUser` 和 `CallTool` 需要同时考虑 builtin 与 user session。

结论：推荐采用。

### 方案 C：独立 BuiltinManager，与用户级 Manager 并行

新增一个独立的 builtin MCP manager，runtime 层同时持有用户级 manager 与 builtin manager。

优点：

- 内置与用户级代码物理隔离。

缺点：

- 大量重复工具发现、调用、结果序列化、并发保护逻辑。
- runtime 层需要感知两个 manager，侵入更大。
- 后续 MCP 行为变更需要维护两套实现。

结论：不推荐。

## 5. 推荐总体方案

采用 **方案 B：Manager 内新增全局 builtin stdio session 池**。

总体流程：

```text
启动 NewServer
  ├─ LoadWebConfig() 读取 config.json
  ├─ 解析 cfg.AppHome
  ├─ mcp.NewManager(store)
  ├─ 读取 cfg.AppHome/mcp_config.json
  ├─ 校验并转换为 []storage.MCPServer
  │    ├─ Transport 固定为 stdio
  │    ├─ Command/Args/Env 来自配置文件
  │    └─ UserID 为空，ID 为 builtin:<sanitized_name>
  ├─ mcpManager.SetBuiltinServers(ctx, servers)
  │    ├─ 对 enabled=true 的内置 server 启动子进程
  │    ├─ 通过 CommandTransport 建立 MCP session
  │    ├─ ListTools 并缓存工具定义
  │    └─ 单个连接失败只记录 warn 并跳过
  └─ runtimeService.SetMCPManager(mcpManager)

用户发起对话
  ├─ toolDefinitionsForUser(ctx, userID)
  │    ├─ EnsureBuiltinSessions(ctx) 检查内置连接，缺失或失效则重连
  │    ├─ EnsureUserSessions(ctx, userID) 保持用户级 MCP 逻辑
  │    └─ ToolsForUser(userID) 返回 builtin tools + user tools
  └─ executeToolCall(ctx, name)
       └─ mcpManager.CallTool(ctx, userID, name, args)
            ├─ 先查用户级 session
            └─ 再查 builtin session；调用失败疑似连接中断时重连并重试一次
```

## 6. `mcp_config.json` 设计

### 6.1 文件路径

默认路径：

```text
<app_home>/mcp_config.json
```

其中 `app_home` 来自 `config.json` 的 `app_home`，解析逻辑见 `internal/config/paths.go:25`。

本次不新增 `config.json` 字段来配置路径，避免引入额外配置复杂度。后续如有多环境需求，可再增加 `mcp_config_path`。

### 6.2 文件不存在时的行为

如果 `<app_home>/mcp_config.json` 不存在：

- 记录 info 日志。
- 启动继续。
- 内置 MCP 列表为空。

理由：保持本地开发和未配置内置 MCP 的部署兼容，不让新能力变成强制依赖。

### 6.3 文件格式错误时的行为

如果文件存在但 JSON 解析失败、字段校验失败或出现重复 server name：

- `NewServer` 返回错误，服务启动失败。

理由：配置文件存在说明部署方期望启用内置 MCP；格式错误应 fail fast，避免服务启动后静默缺失系统能力。

### 6.4 配置 schema

配置只描述 stdio MCP server。建议结构：

```json
{
  "mcp_servers": [
    {
      "name": "filesystem",
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"],
      "env": {
        "NODE_ENV": "production"
      },
      "enabled": true
    },
    {
      "name": "memory",
      "command": "/usr/local/bin/mcp-memory-server",
      "args": [],
      "env": {},
      "enabled": true
    }
  ]
}
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `mcp_servers` | array | 是 | 内置 MCP server 列表 |
| `name` | string | 是 | 内置 server 名称；同一文件内唯一 |
| `command` | string | 是 | stdio MCP server 启动命令 |
| `args` | string array | 否 | 启动参数，缺省为空数组 |
| `env` | object | 否 | 追加或覆盖子进程环境变量，缺省为空对象 |
| `enabled` | bool | 否 | 缺省为 `true`；为 `false` 时跳过连接 |

不支持字段：

- `transport`：内置 MCP 固定为 stdio，不需要配置 transport。
- `url` / `headers`：内置 MCP 不支持远程 HTTP/SSE 连接。
- 环境变量展开：本次直接读取 JSON 字面值；如果需要 `${ENV_NAME}` 展开能力，后续单独设计。

### 6.5 配置转换

加载后转换为 `storage.MCPServer`，用于复用现有 `serverSession` 与签名逻辑：

```go
storage.MCPServer{
    ID:        "builtin:" + sanitizeName(name),
    UserID:    "",
    Name:      "builtin_" + name,
    Transport: "stdio",
    Command:   command,
    Args:      args,
    Env:       env,
    Enabled:   enabled,
}
```

注意：`Name` 增加 `builtin_` 前缀是为了复用现有 `prefixedToolName`，让内置工具对外名称自然变成 `mcp__builtin_<server>__<tool>`。

### 6.6 校验规则

1. `name` 去空格后不能为空。
2. 同一 `mcp_config.json` 内，`sanitizeName(name)` 后不能重复。
3. `command` 去空格后不能为空。
4. `args` 如果出现必须是字符串数组。
5. `env` 如果出现必须是字符串键值对象。
6. `enabled` 缺省为 `true`。
7. 出现 `transport`、`url`、`headers` 字段时应报错，避免误以为内置 MCP 支持远程连接。

## 7. 命名与冲突处理

现有 MCP 工具名格式为：

```text
mcp__<sanitized_server_name>__<tool_name>
```

由 `prefixedToolName` 生成，见 `internal/web/mcp/manager.go:88`。

内置 MCP 的 server name 在加载阶段统一加 `builtin_` 前缀，因此对外工具名为：

```text
mcp__builtin_<sanitized_server_name>__<tool_name>
```

例如 `mcp_config.json` 中：

```json
{"name":"filesystem","command":"npx","args":["-y","@modelcontextprotocol/server-filesystem","/tmp"]}
```

发现到 MCP tool `read_file` 后，对外工具名为：

```text
mcp__builtin_filesystem__read_file
```

这样即使用户自己配置了名为 `filesystem` 的 MCP server，其工具名仍是：

```text
mcp__filesystem__read_file
```

两者不会冲突。

## 8. Manager 详细设计

### 8.1 新增字段

在 `internal/web/mcp/manager.go` 的 `Manager` 中新增：

```go
// 内置 MCP 配置，来自 mcp_config.json。key: serverID
builtinServers map[string]storage.MCPServer

// 内置 MCP 已连接 session。key: serverID
builtinSessions map[string]*serverSession
```

其中内置 `serverID` 固定为：

```text
builtin:<sanitized_server_name>
```

### 8.2 新增配置加载文件

新增 `internal/web/mcp/config.go`，职责：

1. 定义 `builtinConfig` / `builtinServerConfig`。
2. 提供 `LoadBuiltinConfig(path string) ([]storage.MCPServer, error)`。
3. 处理文件不存在、JSON 解析、字段校验、默认值和 ID/Name 生成。

不把该逻辑放到 `internal/config` 的原因：

- `mcp_config.json` 的 schema 与 MCP 领域模型强相关。
- 可复用 `storage.MCPServer` 和 mcp 包内的 `sanitizeName` 规则。
- 避免让全局 config 包依赖 MCP 领域细节。

### 8.3 stdio transport 构造

新增内置专用函数，不修改用户级 `buildTransport`：

```go
func buildBuiltinStdioTransport(server storage.MCPServer) (mcpsdk.Transport, error)
```

职责：

1. 校验 `server.Command` 非空。
2. 使用 `exec.Command(server.Command, server.Args...)` 创建子进程。
3. 合并环境变量：`os.Environ()` + `server.Env`。
4. 返回 `&mcpsdk.CommandTransport{Command: cmd}`。

说明：

- 使用 `exec.Command` 而不是 `exec.CommandContext`，避免连接阶段 context 取消后误杀已建立的长期 stdio 子进程。
- 子进程生命周期由 MCP session 的 `Close()` 管理，SDK 的 `CommandTransport` 会关闭 stdin 并等待进程退出，必要时终止进程。

### 8.4 内置连接建立

新增内置专用连接函数：

```go
func connectBuiltinAndDiscover(ctx context.Context, server storage.MCPServer) (*serverSession, error)
```

它与现有 `connectAndDiscover` 类似，但使用 `buildBuiltinStdioTransport`。

启动时方法：

```go
func (m *Manager) SetBuiltinServers(ctx context.Context, servers []storage.MCPServer)
```

职责：

- 保存 `builtinServers` 配置。
- 对 `Enabled=true` 的 server 建立 stdio 连接并发现工具。
- 单个 server 连接失败只记录 warn，不阻塞其他 server，也不阻塞服务启动。

### 8.5 内置连接保活与重连

新增：

```go
func (m *Manager) EnsureBuiltinSessions(ctx context.Context)
```

触发时机：

1. 服务启动时 `SetBuiltinServers` 会立即连接。
2. 每次用户对话构建工具定义前调用。
3. 内置工具调用失败且疑似连接中断时调用。

重连策略：

- 如果 `builtinSessions` 中缺少某个启用 server 的 session，则立即重连。
- 如果工具调用返回错误，且目标来自内置 session，则：
  1. 将该 builtin session 从 `builtinSessions` 移除并关闭。
  2. 使用 `builtinServers` 中的原始配置重连。
  3. 重连成功后重试本次工具调用一次。
  4. 重试仍失败则返回错误给现有 tool rejected 分支。

不做后台心跳：

- 本次不新增定期 ping/list-tools 心跳，避免引入额外请求和进程噪声。
- 通过“构建工具定义时 ensure + 调用失败时重连一次”满足连接中断后的恢复。

### 8.6 工具定义合并

当前 `toolDefinitionsForUser` 逻辑为：

```go
s.MCP.EnsureUserSessions(ctx, userID)
mcpTools := s.MCP.ToolsForUser(userID)
```

设计改为：

```go
s.MCP.EnsureBuiltinSessions(ctx)
s.MCP.EnsureUserSessions(ctx, userID)
mcpTools := s.MCP.ToolsForUser(userID)
```

`ToolsForUser(userID)` 返回：

1. 内置 MCP tools。
2. 当前用户 MCP tools。

然后按工具名排序，保持输出稳定。

内置 MCP 不参与空闲回收，因此 `ToolsForUser` 不需要因为内置工具读取而更新内置 session 的 `lastUsedAt`。用户级 session 保持现有刷新逻辑。

### 8.7 工具调用路由

`CallTool(ctx, userID, prefixedName, rawArgs)` 查找顺序：

1. 当前用户 session。
2. 内置 session。

由于内置工具名使用 `mcp__builtin_...` 命名空间，正常不会与用户工具冲突。保留“用户优先”的查找顺序是为了维持现有语义，不引入破坏性变化。

内置工具调用失败时执行一次重连重试；用户级 MCP 调用保持现有行为，不增加自动重连。

### 8.8 空闲回收

内置 MCP session **不参与**现有 10 分钟空闲回收：

- `cleanupIdleSessions` 只扫描用户级 `sessions`。
- `builtinSessions` 不按 `lastUsedAt` 回收。
- 内置 session 只会在以下情况下关闭：
  1. 服务退出 `Manager.Close()`。
  2. 内置连接被判定中断后，重连前关闭旧 session。
  3. 后续重新加载内置配置时被替换。

这样满足“服务启动时创建连接，并且不会被回收”的要求。

### 8.9 Invalidate 与 Close

`Invalidate(userID)` 只影响用户级 session，不影响内置 MCP。

`Close()` 同时关闭：

- 全部用户级 session。
- 全部内置 session。
- 后台 cleanup goroutine。

## 9. Server 接入设计

在 `internal/web/app/server.go` 中，创建 MCP manager 后增加：

```go
mcpConfigPath := filepath.Join(cfg.AppHome, "mcp_config.json")
builtinMCPServers, err := mcp.LoadBuiltinConfig(mcpConfigPath)
if err != nil {
    return nil, fmt.Errorf("load builtin mcp config: %w", err)
}

builtinCtx, builtinCancel := context.WithTimeout(context.Background(), 35*time.Second)
defer builtinCancel()
mcpManager.SetBuiltinServers(builtinCtx, builtinMCPServers)
```

需要新增 `path/filepath` import。

说明：

- `SetBuiltinServers` 不因单个 server 连接失败返回错误，因此不会因为单个内置 MCP 不可用阻塞服务启动。
- 35 秒 context 只用于连接和工具发现阶段；不绑定 stdio 子进程生命周期。
- `connectTimeout = 30s` 仍用于 MCP 初始化与 ListTools 的超时边界。

## 10. HTTP API 与前端

本次不新增前端管理页面，不修改现有 `/api/mcp-servers` 行为。

原因：

- 需求指定配置来源为 `mcp_config.json`。
- 内置 MCP 是系统级能力，不应暴露给普通用户修改。

现有用户级 MCP CRUD 仍只读写数据库中的用户配置。

## 11. 日志与可观测性

建议日志：

1. 文件不存在：`mcp: builtin config not found path=<path>, skip builtin servers`。
2. 成功加载：`mcp: loaded builtin stdio servers count=<n>`。
3. 单个连接成功：`mcp: builtin stdio connected server=<name> tools=<n>`。
4. 单个连接失败：`mcp: builtin stdio connect failed server=<name>: <err>`。
5. 内置调用失败并准备重连：`mcp: builtin stdio call failed, reconnecting server=<name>: <err>`。
6. 重连成功：`mcp: builtin stdio reconnected server=<name> tools=<n>`。
7. 重连失败：`mcp: builtin stdio reconnect failed server=<name>: <err>`。

日志不得打印 `env` 中的敏感值。

## 12. 测试设计

### 12.1 配置加载测试

新增 `internal/web/mcp/config_test.go`：

1. 文件不存在返回空列表且无错误。
2. 合法 `mcp_config.json` 解析为 `[]storage.MCPServer`。
3. `enabled` 缺省为 `true`。
4. `enabled=false` 被保留但后续连接跳过。
5. 空 name 报错。
6. 空 command 报错。
7. sanitize 后重名报错，例如 `docs-api` 与 `docs_api`。
8. 出现 `transport` 字段时报错。
9. 出现 `url` 或 `headers` 字段时报错。
10. `args` 必须为字符串数组。
11. `env` 必须为字符串键值对象。

### 12.2 stdio transport 测试

新增或扩展 `internal/web/mcp/transport_test.go`：

1. `buildBuiltinStdioTransport` 对空 command 报错。
2. `buildBuiltinStdioTransport` 正确合并环境变量。
3. 用户级 `buildTransport` 不接受 stdio，确保用户级 stdio 未被打开。

### 12.3 Manager 行为测试

扩展 `internal/web/mcp/manager_test.go`：

1. `ToolsForUser` 返回内置 tools + 用户 tools。
2. 内置 tools 使用 `mcp__builtin_<server>__<tool>` 命名空间。
3. `Invalidate(userID)` 不删除内置 session。
4. `cleanupIdleSessions` 不回收内置 session。
5. `Close()` 会关闭内置 session。
6. 内置 session 缺失时，`EnsureBuiltinSessions` 会按 `builtinServers` 重连。
7. 内置工具调用失败时，会关闭旧 session、重连并重试一次。

如果受 `*mcpsdk.ClientSession` 具体类型限制，调用失败重连测试可通过抽象最小 session 接口或注入连接函数实现；不为了测试重写整个 MCP manager。

### 12.4 Server 启动接入测试

如现有 `app.NewServer` 集成测试环境不便启动真实 stdio MCP，可优先在 config loader 与 manager 层覆盖。Server 层可通过小范围单元测试验证路径拼接和文件缺失不影响启动；真实 stdio 连接行为用 MCP manager 测试覆盖。

### 12.5 回归测试

实施后至少运行：

```bash
go test ./internal/web/mcp
go test ./internal/web/...
```

若全量测试仍有既有失败，需要在结果中明确区分本次 MCP 改动相关测试与既有失败。

## 13. 风险与规避

| 风险 | 影响 | 规避方式 |
|---|---|---|
| 内置 stdio 子进程长期驻留 | 资源占用增加 | 内置 MCP 不回收是需求；仅配置必要 server，服务退出时 `Close()` 统一关闭 |
| 内置工具名与用户工具名冲突 | 模型调用路由不确定 | 内置 server 工具名前缀统一加 `builtin_` |
| 内置 MCP 连接失败阻塞服务 | 服务不可用 | 单个连接失败只 warn 并跳过；配置解析错误才 fail fast |
| 内置连接中断后工具不可用 | 所有用户无法使用该内置工具 | 构建工具定义时 ensure；调用失败时重连并重试一次 |
| 用户修改 MCP 配置误影响内置 MCP | 系统能力被删除或修改 | 内置配置不写 DB，`Invalidate` 只影响用户级 session |
| 配置文件中 env 泄露 | 安全风险 | 文件权限由部署控制；日志不打印 env 内容 |
| stdio 命令执行风险 | 启动非预期进程 | 仅从部署侧 `mcp_config.json` 读取；普通用户无法写入；用户级不开放 stdio |
| 启动耗时增加 | 服务启动变慢 | 单连接使用现有 30s 超时；多个 server 可先串行实现，必要时后续并发优化 |

## 14. 实施范围

预计改动文件：

1. `internal/web/mcp/config.go`
   - 新增 `mcp_config.json` 加载与校验。

2. `internal/web/mcp/config_test.go`
   - 新增配置解析测试。

3. `internal/web/mcp/transport.go`
   - 新增内置专用 stdio transport 构造函数。
   - 保持用户级 `buildTransport` 不支持 stdio。

4. `internal/web/mcp/transport_test.go`
   - 补充 stdio transport 构造和用户级 stdio 禁用测试。

5. `internal/web/mcp/manager.go`
   - 增加 builtin 配置与 session 池。
   - 增加启动连接、缺失重连、调用失败重连重试、工具合并、调用路由支持。
   - 保持内置 session 不参与空闲回收。

6. `internal/web/mcp/manager_test.go`
   - 补充内置 MCP session 行为测试。

7. `internal/web/runtime/tool_registry.go`
   - 在用户 MCP ensure 前增加内置 MCP ensure。

8. `internal/web/app/server.go`
   - 启动时读取 `<app_home>/mcp_config.json` 并初始化内置 MCP。

9. `mcp_config.json`
   - 可选新增示例文件；若用户已有实际配置，则不创建示例以避免误启进程。

不计划改动：

- 数据库结构与迁移。
- `/api/mcp-servers` HTTP API。
- 前端页面。
- 非 MCP 工具注册与执行流程。
- 用户级 MCP transport 范围。
- 内置 MCP 远程 HTTP/SSE 能力。

## 15. 验收标准

1. 服务启动时会读取 `<app_home>/mcp_config.json`。
2. 文件不存在时服务正常启动，内置 MCP 为空。
3. 文件格式错误或字段非法时服务启动失败并返回明确错误。
4. 文件中启用的内置 MCP server 会在启动时通过 stdio 启动子进程、建立连接并发现工具。
5. 任意用户对话时都能获得内置 MCP 工具定义。
6. 内置 MCP 工具名带 `mcp__builtin_` 命名空间，避免与用户 MCP 冲突。
7. 用户级 MCP 配置 CRUD 不展示、不修改、不删除内置 MCP。
8. `Invalidate(userID)` 不影响内置 MCP session。
9. 内置 MCP session 不会被 10 分钟空闲回收机制回收。
10. 内置 MCP 连接中断后，后续构建工具定义或工具调用会触发重连；工具调用场景重连成功后会重试一次。
11. 用户级 MCP 仍不支持 stdio。
12. 非 MCP 工具与用户级 MCP 原有行为不受影响。

## 16. 本次文档阶段不执行的事项

按照需求，本阶段只修改设计文档供审阅，不进行代码改造。待设计确认后，再进入实现与测试阶段。
