# `/clear` 开启全新对话 — 设计文档

## 1. 需求

让 cynosure TUI 的 `/clear` 命令在「清屏」基础上，额外具备：

1. 清除当前对话的上下文历史，开启一个全新的对话（新 session）。
2. 不重启进程。

参考 Claude Code 的机制：`/clear` 等价于「在同一进程内重新开始一段全新会话」，老会话仍以历史文件形式保留在磁盘上（可被 `/resume` 找回），不做物理删除。

## 2. 现状分析

### 2.1 `/clear` 当前行为

`internal/tui/app.go:319-345` 的 `handleSlashCommand`：

```go
case "/clear":
    m.messages = nil          // 只清空 TUI 显示缓冲
    m.resumeSelecting = false
    m.resumeCandidates = nil
    m.appendMessage("system", "已清空当前 TUI 显示上下文")
    return true
```

问题：只清空了 TUI 展示用的 `m.messages`，**没有**触碰后端会话状态。下一轮发话时：

- `RespondToConversation` 仍用 `m.session.Conversation`（同一个 `conversation.ID` / `SessionID`）；
- `Store` 会从内存 map（`messages` / `modelHistory`）或磁盘 `~/.cynosure/session/<SessionID>/{history,model_history}` 把旧历史读回来塞给模型。

所以"上下文"并未真正清除，只是视觉上消失了。

### 2.2 会话状态的存放位置

| 层级 | 位置 | 索引键 |
|---|---|---|
| TUI 显示 | `Model.messages` (`app.go:89`) | — |
| 当前会话引用 | `Model.session.Conversation` (`app.go:28`) | — |
| 内存历史 | `Store.messages` / `Store.cache` / `Store.modelHistory` (`store.go:18,19,24`) | `conversation.ID` |
| 磁盘历史 | `~/.cynosure/session/<SessionID>/{history,model_history}` | `SessionID` |

会话由 `bootstrap.go:86` 在启动时创建，`SessionID = idgen.UUID()`，`ID = idgen.New("conv")`，一次进程一份。

### 2.3 可复用的最佳参考：`/resume`

`/resume` 已经演示了「**在不重启进程的前提下切换当前会话**」的完整模式（`app.go:372-402`）：

```go
conv, history, err := m.session.Resumer.ResumeSession(...)
m.session.Conversation = conv          // 切换当前会话
m.messages = messagesForDisplay(history)
```

`SessionResumer` 接口（`app.go:38-41`）由 `Store` 实现并通过 `root.go:68` 注入。这正是我们要套用的设计方向：**新增一个"开启新会话"的能力，让 TUI 通过接口调用 Store 创建新 Conversation，再切换 `m.session.Conversation`。**

## 3. 设计方案

核心思路：**`/clear` = 在 Store 中创建一个全新的 Conversation（新 conv ID + 新 SessionID），把 TUI 的当前会话指针切到它，并重置相关 UI 状态。** 进程、Runtime、MCP、Store 实例全部复用，不重启。

### 3.1 后端：Store 新增创建新会话的方法

在 `internal/local/store.go` 新增方法：

```go
// StartNewConversation 在当前进程内创建一个全新的空会话（新 SessionID），
// 用于 /clear。不删除任何已有会话的内存或磁盘数据。
func (s *Store) StartNewConversation(ctx context.Context, user storage.User) (storage.Conversation, error) {
    conv := storage.Conversation{
        ID:            idgen.New("conv"),
        SessionID:     idgen.UUID(),
        UserID:        user.ID,
        RootMessageID: idgen.New("msg"),
        Title:         "TUI 会话",
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }
    if err := s.CreateConversation(ctx, conv); err != nil {
        return storage.Conversation{}, err
    }
    return conv, nil
}
```

说明：
- 与 `bootstrap.go:86` 创建初始会话的逻辑保持一致（同样的字段、同样的 idgen 用法）。
- 新会话的内存历史天然为空（map 中没有该 key），磁盘文件在第一次产生历史时才会写入，因此新会话起始即"零上下文"。
- **不删除**旧会话数据：旧 SessionID 的磁盘文件保留，使其后续仍可 `/resume`，与 Claude Code 行为一致。

### 3.2 接口：扩展 TUI 依赖的会话接口

`SessionResumer` 接口目前只含 resume 能力。为表达"会话生命周期管理"，将其扩展（在 `app.go`）：

```go
type SessionResumer interface {
    ListResumableSessions(ctx context.Context, workspaceRoot string) ([]storage.ResumableSession, error)
    ResumeSession(ctx context.Context, sessionID, currentWorkspace string, user storage.User) (storage.Conversation, []storage.Message, error)
    StartNewConversation(ctx context.Context, user storage.User) (storage.Conversation, error)  // 新增
}
```

- `Store` 已经被注入为 `Resumer`（`root.go:68`），新增方法后自动满足接口。
- 需要同步更新测试里的 `fakeSessionResumer`（`events_test.go:216`）实现该新方法，避免编译失败。

> 备选方案：不扩展 `SessionResumer`，而是新增一个独立的小接口（如 `SessionStarter`）并在 `SessionInfo` 上新增一个字段。但这会增加注入点与字段数量；考虑到 `Store` 已是这些会话能力的统一实现方，扩展现有接口更内聚。**推荐扩展现有接口。**

### 3.3 TUI：重写 `/clear` 分支

将 `handleSlashCommand` 的 `case "/clear"` 改为调用新方法并重置全部相关状态。提取为一个 `startNewSession` 辅助方法，便于复用与测试：

```go
case "/clear":
    m.startNewSession()
    return true
```

```go
func (m *Model) startNewSession() {
    // 1) 正在生成时禁止清除（与 /resume 行为一致，避免状态竞争）
    if m.running {
        m.appendMessage("system", "当前正在生成，请先 Ctrl+C 中断后再执行 /clear")
        return
    }
    // 2) 若运行环境不支持（如无 Resumer 注入），退化为仅清屏
    if m.session.Resumer == nil {
        m.messages = nil
        m.resumeSelecting = false
        m.resumeCandidates = nil
        m.appendMessage("system", "已清空当前 TUI 显示上下文")
        return
    }
    // 3) 创建全新会话
    conv, err := m.session.Resumer.StartNewConversation(context.Background(), m.session.User)
    if err != nil {
        m.appendMessage("error", err.Error())
        return
    }
    // 4) 切换当前会话 + 重置 UI 状态
    m.session.Conversation = conv
    m.messages = nil
    m.resumeSelecting = false
    m.resumeCandidates = nil
    m.toolCallCount = 0
    m.contextTokens = 0
    m.appendMessage("system", "已开启全新对话，上下文已清空")
}
```

要点：
- `m.session.Conversation = conv` 是关键：之后 `respond`（`app.go:293`）把新 Conversation 传给 `RespondToConversation`，由于新 conv 没有任何内存/磁盘历史，模型从空上下文开始。
- 重置 `toolCallCount` / `contextTokens`，让底部状态栏正确反映新会话（`contextBudget` 是模型固定容量，保留）。
- `m.running` 保护：正在流式生成时不允许切换会话，避免事件 generation 错乱（沿用 `/resume` 在 `startResumeSelection` 中的同款保护）。

### 3.4 `/resume` 选择态里的 `/clear` 分支

`handleResumeSelection`（`app.go:372-384`）中也处理了 `/clear`，目前同样只清屏。为保持一致，改为：先退出 resume 选择态，再调用 `m.startNewSession()`。

```go
if text == "/cancel" || text == "/clear" {
    m.resumeSelecting = false
    m.resumeCandidates = nil
    if text == "/clear" {
        m.startNewSession()
    } else {
        m.appendMessage("system", "已取消恢复历史会话")
    }
    return true
}
```

## 4. 改动文件清单

| 文件 | 改动 |
|---|---|
| `internal/local/store.go` | 新增 `StartNewConversation` 方法；补 `idgen` import |
| `internal/tui/app.go` | 扩展 `SessionResumer` 接口；`/clear` 改调 `startNewSession`；新增 `startNewSession` 方法；调整 `handleResumeSelection` 的 `/clear` 分支 |
| `internal/tui/events_test.go` | 为 `fakeSessionResumer` 补 `StartNewConversation` 实现 |

> 不改动：`bootstrap.go`、`RespondToConversation` 流程、磁盘存储格式、`/resume` 主流程。改动面小且聚焦。

## 5. 行为对照

| 场景 | 改前 | 改后 |
|---|---|---|
| `/clear` 后界面 | 清屏 | 清屏 |
| `/clear` 后模型上下文 | 仍含旧历史 | 空（全新会话） |
| `/clear` 后 SessionID | 不变 | 新 UUID |
| 进程 | 不重启 | 不重启 |
| 旧会话可否 `/resume` | — | 可以（磁盘文件保留） |
| 生成中执行 `/clear` | 直接清屏 | 提示先 Ctrl+C 中断 |

## 6. 验证方式

1. 启动 TUI，发一句话建立上下文（如"记住我叫张三"）。
2. 执行 `/clear`，确认提示"已开启全新对话"。
3. 再问"我叫什么"，模型应**不知道**（上下文已清）。
4. 执行 `/resume`，确认上一段会话仍在列表中、可恢复。
5. 全程进程 PID 不变（无重启）。
6. `go build ./...` 与现有测试（`internal/tui`）通过。

## 7. 风险与边界

- **并发**：`m.running` 守卫避免与流式生成竞争。
- **内存增长**：每次 `/clear` 新增一条 Conversation 记录在 `Store.conversations` map 中（仅元数据，极小），可接受；旧历史磁盘文件保留是有意为之（供 resume）。
- **接口破坏**：扩展 `SessionResumer` 需同步更新测试 fake，已列入清单。
