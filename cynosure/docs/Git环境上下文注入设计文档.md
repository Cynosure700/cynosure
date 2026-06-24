# Git 环境上下文注入设计文档

## 1. 背景与目标

当前 cynosure 在每轮请求前由 `assistant.BuildSystemPrompt` 统一拼接系统提示词，动态段落依次为 `identity`、`workspace`、`system-reminder`（CYNOSURE.MD）、`tools`、`skills`、`memory`（`internal/assistant/prompt.go:147`、`internal/assistant/prompt.go:202`）。启动期通过 `local.Bootstrap` 读取进程级上下文（如 CYNOSURE.MD），再注入 runtime service（`internal/local/bootstrap.go:71`、`internal/local/bootstrap.go:85`）。

本次目标：为系统提示词新增 **Git 环境上下文** 段落，让模型在回答时了解当前工作区的 Git 状态。

具体要求：

1. 采集当前工作区的 Git 信息：当前分支、主分支、工作区变更状态、最近提交、Git 用户名。
2. 按固定格式拼接成一段文本，作为对话开始时的 Git 快照（snapshot），会话过程中不刷新。
3. `status` 输出超过 2000 字符时截断，并追加截断提示。
4. 将该段落注入到 `memory` 段落之前。
5. 若本地不是 Git 仓库 / 无 Git 信息 / 关键 Git 命令执行失败，则**整段忽略**，不影响其余提示词。
6. 其他策略不变，不影响现有功能。

## 2. 设计原则

- **启动期采集、一次性注入**：Git 状态属于"对话开始时的快照"，随 TUI 启动采集一次，注入 runtime；不在每轮对话重复执行 Git 命令，避免会话中状态漂移，也避免每轮启动子进程的开销。这与 CYNOSURE.MD 的注入方式一致，也符合提示词文案"will not update during the conversation"。
- **统一提示词出口**：不在 TUI 或 LLM client 中拼字符串，继续通过 `assistant.BuildSystemPrompt` 统一构造，新增独立 `<git-status>` 段落。
- **缺失容忍、失败降级**：不是 Git 仓库 / Git 未安装 / 关键命令失败 → 整段不渲染；非关键字段（主分支、提交、用户名）缺失 → 只省略对应行，仍保留段落。
- **只读、不修改仓库**：只执行只读 Git 命令（`rev-parse`、`branch`、`status`、`log`、`config`、`symbolic-ref`），绝不修改仓库状态。
- **不阻塞启动**：每条 Git 命令带看门狗超时（复用 `internal/tools/bash.go` 的 `time.AfterFunc` 模式，不引入 `context.WithTimeout`，与项目既有超时风格一致）；超时即视为该命令失败并降级。
- **最小改动**：复用现有"启动期装配 → runtime 持有 → assistant 渲染"链路，新增一个独立采集包，不引入配置文件、watcher 或热更新。

## 3. 方案比较

### 方案 A：每轮请求时实时采集（不推荐）

在 `buildSystemPromptWithMemory` 每轮执行 Git 命令。

- 优点：状态实时。
- 缺点：与"对话开始快照"语义矛盾（提示词文案承诺 snapshot）；每轮启动多个子进程，拖慢响应；与 CYNOSURE.MD 的启动期注入风格不一致。

### 方案 B：启动期采集，独立段落注入（推荐）

启动期在 `local.Bootstrap` 采集 Git 快照并格式化为文本，注入 runtime；`BuildSystemPrompt` 渲染为独立 `<git-status>` 段落，位置在 `<memory>` 之前。

- 优点：边界清晰、与 CYNOSURE.MD/workspace 等动态段落一致；快照语义自洽；启动只执行一次；易于单测。
- 缺点：需新增采集包、runtime 字段与 setter，并补充单测（改动可控）。

### 方案 C：并入 CYNOSURE.MD 的 `<system-reminder>`（不推荐）

把 Git 文本拼进现有 system-reminder 段落。

- 缺点：Git 状态与"用户/项目说明"语义不同，混在一起边界不清；需求明确要求注入到 `memory` 段落之前，单独成段更直观。

**采用方案 B。**

## 4. 数据采集设计

新增采集包 `internal/gitcontext`，对外暴露：

```go
// Status 是对话开始时采集的 Git 快照原始字段。
type Status struct {
    IsRepo        bool   // 是否为 Git 仓库（整段是否渲染的总开关）
    Branch        string // 当前分支
    MainBranch    string // 主分支
    WorkTree      string // 工作区变更状态（git status --short 原始输出，未截断）
    RecentCommits string // 最近提交（git log --oneline -n 5 原始输出）
    UserName      string // Git 用户名
}

// Collect 在 workspaceRoot 下采集 Git 快照。非 Git 仓库返回 IsRepo=false。
func Collect(workspaceRoot string) Status

// Format 将快照格式化为注入提示词的文本块（不含外层 XML 标签）。
// 非 Git 仓库返回空串。
func Format(s Status) string

// CollectAndFormat 是 Collect+Format 的便捷组合，供启动流程调用。
func CollectAndFormat(workspaceRoot string) string
```

### 4.1 采集命令

所有命令在 `workspaceRoot` 目录下执行（`cmd.Dir = workspaceRoot`），均为只读，并带看门狗超时。

| 字段 | 命令 | 说明 |
|------|------|------|
| 仓库判定 | `git rev-parse --is-inside-work-tree` | 输出非 `true` 或执行失败 → 非仓库，整段忽略 |
| `Branch` | `git branch --show-current` | 当前分支；detached HEAD 时为空 |
| `MainBranch` | `git symbolic-ref --short refs/remotes/origin/HEAD` | 主分支（去掉 `origin/` 前缀），见 4.2 |
| `WorkTree` | `git --no-optional-locks status --short` | 工作区变更状态 |
| `RecentCommits` | `git --no-optional-locks log --oneline -n 5` | 最近 5 条提交 |
| `UserName` | `git config user.name` | Git 用户名 |

实现要点：

- 先执行 `rev-parse` 判定仓库；失败（含 Git 未安装、目录非仓库）则直接返回 `Status{IsRepo:false}`，后续命令不再执行。
- 其余命令为最佳努力（best-effort）：任一命令失败/超时，仅令对应字段为空，不影响整段渲染。
- 各命令输出统一 `strings.TrimSpace`。
- 采用与 `internal/tools/bash.go:17` 相同的 `time.AfterFunc` 看门狗超时（不使用 `context.WithTimeout`，遵循项目既有超时风格），单条命令超时常量在采集包内硬编码（建议 5s，足以覆盖本地 Git 操作）。
- 命令顺序执行即可：均为本地只读操作，总耗时通常 <100ms；不引入并发以保持实现简单（参考来源使用并发，本设计认为本地命令无并发必要，作为后续可选优化）。

### 4.2 主分支命令的取舍（已确认采用离线方案）

需求参考来源使用 `git remote show origin` 解析 HEAD 分支获取主分支。但该命令会**发起网络请求**访问远端，在远端不可达 / 网络慢 / 鉴权交互时可能在 **TUI 启动期阻塞或挂起**，与"不阻塞启动"原则冲突。

本设计采用离线命令 `git symbolic-ref --short refs/remotes/origin/HEAD`：

- 它读取本地 `refs/remotes/origin/HEAD` 引用（`git clone` 时已建立），完全离线、毫秒级返回。
- 结果形如 `origin/main`，去掉 `origin/` 前缀即为主分支。
- 若本地仓库无该引用（如纯本地仓库、未设置 origin/HEAD），命令失败 → `MainBranch` 为空 → 仅省略主分支行。

> 已采纳：用离线 `symbolic-ref` 替代参考来源的 `git remote show origin`，避免 TUI 启动期网络阻塞，结果一致。

## 5. 文本格式设计

`Format` 输出的文本块，各段以 `\n\n` 连接，字段为空则省略对应段（`WorkTree` 例外，见下）：

```
This is the git status at the start of the conversation. Note that this status is a snapshot in time, and will not update during the conversation.

Current branch: <Branch>

Main branch (you will usually use this for PRs): <MainBranch>

Git user: <UserName>

Status:
<WorkTree 或 (clean)>

Recent commits:
<RecentCommits>
```

规则：

- 第一段固定前言始终存在（只要 `IsRepo=true`）。
- `Current branch` / `Main branch` / `Git user` / `Recent commits`：对应字段为空时整段省略。
- `Status:` 段始终存在；`WorkTree` 为空（工作区干净）时渲染为 `(clean)`，与启动期会话快照风格一致。
- 段内多行内容（status、log）保留原始换行。

### 5.1 截断逻辑

- 常量 `MaxStatusChars = 2000`，按 **rune（字符）** 计数，避免多字节路径被截断成乱码。
- 当 `WorkTree` 的 rune 数 > 2000 时，截断到前 2000 个 rune，并在末尾追加：

```
... (truncated because it exceeds 2k characters. If you need more information, run "git status" using the bash tool)
```

> 说明：参考来源提示语为 `... run "git status" using BashTool`。本项目终端命令工具名为 `bash`，故将提示语适配为 `using the bash tool`，使指引指向真实可用工具。

## 6. 模块设计

### 6.1 `internal/gitcontext`（新增包）

- `gitcontext.go`：`Status` 结构、`Collect`、`Format`、`CollectAndFormat`，以及内部 `runGitCommand(workspaceRoot string, args ...string) (string, bool)` 看门狗执行器。
- 不依赖 `assistant` / `runtime` / `config`，保持采集逻辑独立可测。

### 6.2 `internal/assistant`：渲染独立段落

- 扩展 `PromptOptions`（`internal/assistant/prompt.go:122`）新增字段：

```go
GitStatus string // 已格式化的 Git 快照文本块（不含外层标签），为空则不渲染
```

- 在 `BuildSystemPrompt` 中，于 `<memory>` 段落**之前**插入：

```go
if gitStatus := strings.TrimSpace(opts.GitStatus); gitStatus != "" {
    sections = append(sections, renderTag("git-status", gitStatus))
}
```

- 段落顺序变为：`identity` → `workspace` → `system-reminder`(CYNOSURE.MD) → `tools` → `skills` → `git-status` → `memory`。满足"注入到记忆段落的前面"。
- 该段独立于 memory：即使 memory 段为空，只要 `GitStatus` 非空仍会渲染。

### 6.3 `internal/agent/runtime`：持有并传递

- `Service` 新增字段（`internal/agent/runtime/runtime.go:57`）：

```go
GitStatus string
```

- 新增 setter（与 `SetCynosureMarkdownContext` 并列，`runtime.go:109`）：

```go
func (s *Service) SetGitStatusContext(text string) { s.GitStatus = text }
```

- `buildSystemPromptWithMemory`（`internal/agent/runtime/prompt_builder.go:41`）调用 `assistant.BuildSystemPrompt` 时传入 `GitStatus: s.GitStatus`，与现有 `WorkingDirectory`、`ToolNames` 并列。

### 6.4 `internal/local`：启动期采集并注入

- 在 `local.Bootstrap`（`internal/local/bootstrap.go:71` 一带，与 CYNOSURE.MD 装配相邻）新增：

```go
gitStatus := gitcontext.CollectAndFormat(cfg.WorkspaceRoot)
...
runtimeService.SetGitStatusContext(gitStatus)
```

- 采集失败/非仓库时 `CollectAndFormat` 返回空串，`SetGitStatusContext("")` → 段落不渲染，**不阻断启动、不返回错误**（与"失败降级"原则一致）。

### 6.5 范围边界

- 仅作用于**主 Agent** 系统提示词。子 Agent（如 explore）使用独立提示词构造，本次不注入 Git 段落（与 CYNOSURE.MD 同样仅主 Agent），避免扩大范围。
- `/clear` 在进程内开启新会话时不重新采集 Git（与 CYNOSURE.MD 行为一致，保持"进程启动期快照"语义）。如需会话级刷新，作为后续独立需求。

## 7. 测试设计

### 7.1 `internal/gitcontext` 单测

- `Format`（不依赖真实 Git，构造 `Status`）：
  - 全字段齐全时包含前言、各行与正确顺序。
  - `WorkTree` 为空渲染 `(clean)`。
  - `MainBranch`/`UserName`/`RecentCommits` 为空时省略对应行。
  - `WorkTree` 超过 2000 rune 时被截断且追加截断提示；恰好 2000 不截断。
  - `IsRepo=false` 返回空串。
- `Collect`（在 `t.TempDir()` 内 `git init` + 配置 user + 提交）：
  - 正常仓库返回 `IsRepo=true`、`Branch` 非空、`RecentCommits` 含提交、`UserName` 正确。
  - 非 Git 目录返回 `IsRepo=false`。
  - 工作区有改动时 `WorkTree` 非空。
  - 测试用 `t.Setenv` 隔离 `GIT_CONFIG_*` / user 配置，避免污染全局 Git 配置；测试若依赖 `git` 可执行文件，缺失时 `t.Skip`。

### 7.2 `internal/assistant` 单测（扩展 `prompt_test.go`）

- `GitStatus` 非空时输出包含 `<git-status>` 及其内容。
- 段落顺序：`</skills>`（或 `</tools>`）在 `<git-status>` 之前，`</git-status>` 在 `<memory>` 之前。
- `GitStatus` 为空时不输出 `<git-status>`。
- 即使 `MemorySection` 为空，`GitStatus` 非空仍渲染 `<git-status>`。

### 7.3 `internal/agent/runtime` 单测

- `SetGitStatusContext` 后，`buildSystemPromptWithMemory`/`buildSystemPrompt` 产出的系统提示词包含该 Git 文本。

## 8. 实施步骤

1. 新增 `internal/gitcontext` 包（`Collect`/`Format`/`CollectAndFormat` + 看门狗执行器）并补充单测。
2. 扩展 `assistant.PromptOptions` 与 `BuildSystemPrompt`，在 `<memory>` 前渲染 `<git-status>`，补充单测。
3. `runtime.Service` 新增 `GitStatus` 字段与 `SetGitStatusContext`，在 `prompt_builder` 传递；补充单测。
4. `local.Bootstrap` 启动期调用 `gitcontext.CollectAndFormat` 并 `SetGitStatusContext`。
5. 运行 `go build ./...` 与相关测试：`go test ./internal/gitcontext/ ./internal/assistant/ ./internal/agent/runtime/ ./internal/local/ -count=1`，最后 `go test ./... -count=1` 回归。
6. 同步更新 `README.md` 与本设计文档中涉及的能力说明（如 README 描述系统提示词注入能力，则补充 Git 段落）。

## 9. 风险与边界

- **网络阻塞**：通过离线 `symbolic-ref` 获取主分支规避；所有命令带 5s 看门狗超时兜底。
- **大仓库 status 过大**：通过 2000 rune 截断与提示控制 token。
- **Git 未安装 / 非仓库**：整段忽略，行为与未启用一致，不报错、不阻断启动。
- **快照不刷新**：进程启动后仓库变更不反映到提示词，符合"对话开始快照"语义；提示词文案已明示。
- **测试环境无 git**：`Collect` 相关用例缺 `git` 时 `t.Skip`，不阻塞 CI。

## 10. 自检结论

- 段落位置（memory 之前）、采集字段、截断阈值、失败降级均与需求一致。
- 主分支命令已确认采用离线 `symbolic-ref`（理由见 4.2）。
- 设计仅覆盖主 Agent 启动期 Git 快照注入，不含 watcher、热更新、会话级刷新、子 Agent 注入等额外功能。
