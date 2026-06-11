# MCP 连接空闲回收机制设计

## 1. 背景与目标

当前 go-agent 的 MCP 客户端连接由 `internal/web/mcp.Manager` 按用户和 MCP server 配置维度缓存。对话构建工具定义时会调用 `EnsureUserSessions` 建立或复用连接，工具调用时通过 `CallTool` 路由到对应 `ClientSession`。

现有连接关闭路径主要包括：

- MCP 配置新增、更新、启停、删除后调用 `Invalidate` 关闭该用户全部 MCP 连接，见 `internal/web/app/mcp_handlers.go:92`、`internal/web/app/mcp_handlers.go:150`、`internal/web/app/mcp_handlers.go:165`、`internal/web/app/mcp_handlers.go:172`。
- 服务退出时可调用 `Manager.Close` 关闭全部连接，见 `internal/web/mcp/manager.go:267`。
- 配置变更或启用列表变化时，`EnsureUserSessions` 会关闭旧连接或禁用连接，见 `internal/web/mcp/manager.go:108`、`internal/web/mcp/manager.go:120`。

但当前没有基于空闲时间的连接回收机制：只要用户不修改 MCP 配置，已建立的远程 `sse` / `streamable` 会话会长期驻留在内存中。

本次目标：

1. 为 MCP 连接建立空闲回收机制。
2. MCP 连接空闲超过 10 分钟后自动关闭并从缓存移除。
3. 不影响现有 MCP 配置 CRUD、工具发现、工具调用、测试连接、非 MCP 工具与对话主流程。

## 2. 现状分析

### 2.1 MCP 连接缓存结构

`Manager` 当前使用互斥锁保护一个二级 map：

```go
// userID -> (serverID -> session)
sessions map[string]map[string]*serverSession
```

定义位置见 `internal/web/mcp/manager.go:39`。

每个 `serverSession` 保存 MCP server 配置、配置指纹、`ClientSession`、已发现工具定义和工具名前缀映射，见 `internal/web/mcp/manager.go:30`。

### 2.2 连接建立与复用

对话准备工具定义时会调用：

- `Service.toolDefinitionsForUser`：先获取内置工具，再调用 MCP manager，见 `internal/web/runtime/tool_registry.go:74`。
- `Manager.EnsureUserSessions`：读取启用的 MCP server 配置，连接缺失或配置变更的 server，并复用配置未变的连接，见 `internal/web/mcp/manager.go:83`。

### 2.3 工具调用

MCP 工具名以 `mcp__` 开头时，`executeToolCall` 会转发给 `Manager.CallTool`，见 `internal/web/runtime/tool_registry.go:219`。

`CallTool` 当前在锁内找到目标 `serverSession` 后释放锁，再调用 `target.session.CallTool`，见 `internal/web/mcp/manager.go:203`。这避免了工具调用期间长期持有全局锁，但也意味着未来新增后台回收时，必须避免回收 goroutine 在工具调用进行中关闭同一个 session。

### 2.4 传输类型

MCP transport 当前只支持：

- `sse`
- `streamable`

实现见 `internal/web/mcp/transport.go:37`。本次回收机制不改变 transport 选择与构造逻辑。

## 3. 设计原则

1. **最小侵入**：回收逻辑集中在 `internal/web/mcp` 包内，不改动 MCP 配置存储结构、不改动 HTTP API、不改动非 MCP 工具。
2. **按真实使用更新空闲时间**：连接被用于工具定义聚合或工具调用时，视为活跃。
3. **不打断正在执行的工具调用**：回收只关闭没有 in-flight 调用的 session。
4. **懒加载保持不变**：连接被回收后，下次对话构建工具定义时仍由 `EnsureUserSessions` 自动重连。
5. **失败降级保持不变**：重连失败时仍只记录日志并跳过该 MCP server，不阻断对话。

## 4. 可选方案

### 方案 A：只在 `EnsureUserSessions` 中顺手回收

每次对话准备工具定义时，检查当前用户已有 session 的 `lastUsedAt`，超过 10 分钟则关闭。

优点：

- 实现最简单，不需要后台 goroutine。
- 不引入 Manager 生命周期新增方法。

缺点：

- 只有用户再次发起对话时才触发回收；如果用户彻底不再使用，连接仍会一直驻留。
- 不满足“空闲十分钟后回收”的严格语义。

### 方案 B：后台 ticker 定期扫描并回收（推荐）

`Manager` 创建时启动后台 goroutine，每隔一段时间扫描所有 session，关闭 `now - lastUsedAt >= 10m` 且无 in-flight 调用的连接。

优点：

- 语义清晰，能在空闲达到阈值后自动回收。
- 回收逻辑集中在 Manager 内部，不影响调用方。
- 与现有懒加载兼容，回收后下次自然重连。

缺点：

- 需要为 Manager 增加后台任务生命周期控制。
- 需要补充并发保护，避免工具调用期间关闭 session。

### 方案 C：为每个 session 维护独立 timer

每个连接创建一个 `time.Timer`，每次使用后 reset，timer 到期后关闭连接。

优点：

- 回收时间更精确。
- 不需要定期全量扫描。

缺点：

- 每个 session 都有独立 timer，状态更复杂。
- 配置变更、Invalidate、Close、工具调用并发时更容易遗漏 timer stop/reset。
- 当前连接数量预期不大，精确到秒的回收收益有限。

## 5. 推荐方案

采用 **方案 B：后台 ticker 定期扫描并回收**。

推荐原因：

- 能满足“空闲十分钟后回收”的需求。
- 实现复杂度低于 per-session timer。
- 与当前 `Manager` 集中维护 session map 的结构匹配。
- 可以通过 `activeCalls` 避免关闭正在执行的 MCP 工具调用。

建议参数：

```go
const (
    idleTimeout     = 10 * time.Minute
    cleanupInterval = 1 * time.Minute
)
```

实际回收时间为 `10m ~ 11m` 之间，取决于扫描周期。该误差可接受；如果后续需要更精确，可将扫描周期调整为 30 秒，但本次不建议引入额外配置项。

## 6. 详细设计

### 6.1 `serverSession` 增加运行态字段

在 `internal/web/mcp/manager.go` 的 `serverSession` 中新增字段：

```go
lastUsedAt time.Time
activeCalls int
```

含义：

- `lastUsedAt`：最近一次被使用的时间。
- `activeCalls`：当前正在执行的 MCP 工具调用数量。

初始化：

- `connectAndDiscover` 成功返回前设置 `lastUsedAt = time.Now()`。
- `activeCalls` 默认为 0。

### 6.2 活跃时间更新规则

以下行为应刷新 `lastUsedAt`：

1. 新连接创建成功。
2. `EnsureUserSessions` 复用已有连接时。
3. `ToolsForUser` 返回某用户 MCP 工具定义时。
4. `CallTool` 成功找到目标 session 并准备发起调用时。
5. `CallTool` 调用结束时再次刷新，避免长调用刚结束就被下一轮扫描回收。

这里将“对话构建工具列表”视为使用，是因为模型只有在收到工具定义后才可能调用 MCP 工具；只要用户仍在发起对话，MCP 连接就不应被视为空闲。

### 6.3 工具调用并发保护

当前 `CallTool` 在锁内找到 `target` 后释放锁，然后执行 `target.session.CallTool`。新增后台回收后，如果不做保护，可能出现：

1. `CallTool` 找到 target 并释放锁。
2. cleanup goroutine 获得锁，发现 target 空闲并关闭 session。
3. `CallTool` 使用已关闭 session 发起请求。

因此 `CallTool` 应调整为：

1. 在锁内找到 target。
2. 同步执行：`target.activeCalls++`，`target.lastUsedAt = now`。
3. 解锁后执行 MCP `CallTool`。
4. defer 中重新加锁：`target.activeCalls--`，`target.lastUsedAt = time.Now()`。

cleanup 扫描时只回收：

```go
sess.activeCalls == 0 && now.Sub(sess.lastUsedAt) >= idleTimeout
```

这样不会关闭正在调用的 session。

### 6.4 后台回收生命周期

`Manager` 增加字段：

```go
done chan struct{}
```

`NewManager` 创建 Manager 后启动后台 goroutine：

```go
go manager.cleanupLoop()
```

`cleanupLoop` 使用 `time.NewTicker(cleanupInterval)`，每次 tick 调用 `cleanupIdleSessions(time.Now())`。

`Close` 需要先关闭 `done`，再关闭所有 session 并清空 map。为避免重复 Close 导致 panic，可增加 `sync.Once`：

```go
closeOnce sync.Once
```

`Close` 中：

```go
m.closeOnce.Do(func() { close(m.done) })
```

然后保持现有关闭全部 session 的行为。

### 6.5 回收执行方式

建议 `cleanupIdleSessions` 在锁内从 map 删除待回收 session，但把实际 `session.Close()` 放到锁外执行，避免关闭远端连接时长时间持有 `Manager.mu`。

流程：

1. 加锁。
2. 遍历 `sessions[userID][serverID]`。
3. 找到满足空闲阈值且 `activeCalls == 0` 的 session。
4. 从 map 删除。
5. 如果某个用户的 server map 已空，删除该 userID key。
6. 解锁。
7. 逐个调用 `sess.session.Close()`。
8. 记录 debug/info 级日志，例如：`mcp: idle session closed user=<id> server=<name>`。

### 6.6 与现有流程的兼容性

#### 对话工具定义

`Service.toolDefinitionsForUser` 仍调用 `EnsureUserSessions` 与 `ToolsForUser`，见 `internal/web/runtime/tool_registry.go:74`。连接被回收后，该路径会在下一次对话时重新建立连接，不需要改 runtime 层。

#### MCP 配置变更

现有 `Invalidate(userID)` 仍立即关闭并删除该用户所有连接。cleanup goroutine 如果同时扫描，需要通过同一把 `Manager.mu` 串行化，避免重复删除。即便 close 同一个底层 session 返回错误，也沿用现有 `_ = session.Close()` 的容错风格。

#### 测试连接

`TestServer` 当前临时连接并 defer close，见 `internal/web/mcp/manager.go:251`。它不写入 `Manager.sessions`，不参与空闲回收。

#### 服务退出

当前 `Server.Run` 直接调用 `http.ListenAndServe`，见 `internal/web/app/server.go:107`，没有显式优雅退出流程。本次设计只保证 `Manager.Close` 被调用时能停止后台 goroutine；是否在进程退出路径接入统一关闭不作为本次必要范围。

## 7. 测试设计

建议新增 `internal/web/mcp/manager_test.go`，覆盖纯 Manager 行为。为了便于测试，可将回收逻辑拆为可直接调用的非导出方法：

```go
func (m *Manager) cleanupIdleSessions(now time.Time)
```

测试用例：

1. **未超时不回收**
   - 构造 `lastUsedAt = now - 9m` 的 session。
   - 调用 `cleanupIdleSessions(now)`。
   - 断言 session 仍存在。

2. **超过 10 分钟回收**
   - 构造 `lastUsedAt = now - 10m - 1s` 的 session。
   - 调用 cleanup。
   - 断言 session 从 map 删除。

3. **有 active call 不回收**
   - 构造超时 session，但 `activeCalls = 1`。
   - 调用 cleanup。
   - 断言 session 仍存在。

4. **ToolsForUser 刷新活跃时间**
   - 构造已连接 session。
   - 调用 `ToolsForUser(userID)`。
   - 断言 `lastUsedAt` 被更新。

5. **Close 停止后台任务且清空连接**
   - 调用 `Close`。
   - 断言 session map 被清空。
   - 重复调用 `Close` 不 panic。

由于 `serverSession.session` 是 SDK 的具体类型 `*mcpsdk.ClientSession`，单元测试中如果难以 mock close 行为，可以优先断言 map 删除和状态变化；关闭错误沿用现有忽略策略，不作为核心断言。

## 8. 实施范围

预计只需要改动：

1. `internal/web/mcp/manager.go`
   - 增加 idle timeout 常量。
   - 增加 `lastUsedAt` / `activeCalls`。
   - 增加 cleanup goroutine 和 `cleanupIdleSessions`。
   - 调整 `ToolsForUser`、`CallTool`、`Close` 的状态维护。

2. `internal/web/mcp/manager_test.go`
   - 新增上述单元测试。

不计划改动：

- 数据库表结构。
- MCP HTTP API。
- 前端页面。
- `transport.go`。
- 非 MCP 工具注册与执行流程。

## 9. 风险与规避

| 风险 | 影响 | 规避方式 |
|---|---|---|
| 回收 goroutine 关闭正在调用的 session | MCP 工具调用失败 | 引入 `activeCalls`，cleanup 跳过 in-flight session |
| 后台 goroutine 泄漏 | 测试或服务生命周期中残留任务 | `Manager.Close` 关闭 `done`，并用 `sync.Once` 防止重复关闭 |
| 连接刚被回收后用户继续使用 | 下一次工具定义构建时重新连接 | 保持 `EnsureUserSessions` 懒加载逻辑不变 |
| 回收扫描持锁过久 | 阻塞工具定义或工具调用 | 锁内只删除 map，锁外执行 `session.Close()` |
| 10 分钟阈值不可配置 | 未来不同环境需求不同 | 本次按需求硬编码，避免引入未要求的配置；后续可扩展配置项 |

## 10. 验收标准

1. MCP session 空闲 10 分钟后会被后台回收，从 `Manager.sessions` 中移除并调用 `ClientSession.Close()`。
2. 正在执行 MCP 工具调用的 session 不会被回收。
3. 回收后的 MCP server 在下一次对话准备工具定义时可按现有 `EnsureUserSessions` 逻辑重新连接。
4. MCP 配置 CRUD 的 `Invalidate` 行为保持不变。
5. `TestServer` 临时连接行为保持不变。
6. 非 MCP 工具执行路径不受影响。
7. 新增或调整的 Go 单元测试通过。

## 11. 本次文档阶段不执行的事项

按照需求，本阶段只生成设计文档供审阅，不进行代码改造。待设计确认后，再进入实现与测试阶段。
