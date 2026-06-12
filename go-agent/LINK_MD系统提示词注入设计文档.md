# LINK.MD 系统提示词注入设计文档

## 1. 背景与目标

当前 go-agent 已迁移到本地 TUI 启动链路：`main.go` 调用 CLI，CLI 解析 `--cwd` 后进入本地启动流程，`local.Bootstrap` 初始化配置、runtime、skills、MCP 与会话，runtime 在每轮请求前构造系统提示词。

现状关键链路：

- CLI 入口解析当前工作区并启动 TUI：`main.go:9`、`internal/cli/root.go:50`、`internal/cli/root.go:58`。
- 本地启动时通过 `config.LoadLocalConfig(cwd)` 解析工作区，并读取基础系统提示词：`internal/local/bootstrap.go:34`、`internal/local/bootstrap.go:59`。
- 系统提示词由 `assistant.BuildSystemPrompt` 统一拼接 identity、workspace、tools、skills、memory 等动态段落：`internal/assistant/prompt.go:35`。
- runtime 每轮调用模型前通过 `buildSystemPromptWithMemory` 传入基础提示词、工作目录、skills、memory 与工具列表：`internal/agent/runtime/prompt_builder.go:24`。

本次目标：

1. 启动 TUI 进程时自动读取用户级 `~/.link/LINK.MD`。
2. 启动 TUI 进程时自动读取当前工作目录 `<cwd>/.link/LINK.MD`。
3. 将读到的内容注入系统提示词，使模型在回答用户问题时可参考这些上下文。
4. 路径必须运行时解析：用户家目录使用 `os.UserHomeDir()`，当前工作目录使用 CLI 传入并校验后的 `cfg.WorkspaceRoot`，不能硬编码示例中的绝对路径。
5. 缺失任一 `LINK.MD` 不应阻断 TUI 启动；存在但读取失败时给出明确错误。

> 说明：需求中的 `~./.link/目录` 按常见 home 目录语义理解为 `~/.link/目录`。

## 2. 设计原则

- **启动期读取，一次性注入**：`LINK.MD` 属于进程级上下文，随 TUI 启动读取；不在每轮对话重复访问文件系统，避免会话中内容漂移。
- **统一提示词出口**：不在 TUI 或 LLM client 中拼字符串，继续通过 `assistant.BuildSystemPrompt` 统一构造最终系统提示词。
- **用户级先于工作区级**：提示词中先展示全局个人说明，再展示当前项目说明，符合“全局默认 + 项目补充/覆盖”的阅读顺序。
- **缺失容忍，错误明确**：文件不存在视为没有额外上下文；权限错误、路径为目录等异常返回包含路径的错误，帮助用户定位。
- **最小改动**：复用已有 `config` 路径解析与 `assistant` 提示词构造模式，不引入新的配置文件格式或后台 watcher。

## 3. 方案比较

### 方案 A：读取后直接拼到 `basePrompt`（不推荐）

在 `local.Bootstrap` 中读取 `LINK.MD`，把内容直接追加到 `basePrompt` 后，再调用 `runtimeService.SetBasePrompt`。

- 优点：改动最少，不需要调整 `PromptOptions`。
- 缺点：`LINK.MD` 上下文与基础身份提示词混在 `<identity>` 标签里，测试和后续维护不清晰；未来如果需要在系统提示词中调整位置，只能继续操作大字符串。

### 方案 B：新增 `PromptOptions.LinkContext` 独立段落（推荐）

新增一个结构化的 `LinkContext`，由启动流程读取后传给 runtime，最终由 `assistant.BuildSystemPrompt` 渲染为独立 `<system-reminder>` 段落。

- 优点：边界清晰，和 identity/workspace/tools/skills/memory 一样是独立动态段落；便于单测验证；后续扩展更多上下文来源时不会污染基础提示词。
- 缺点：需要给 runtime service 增加字段和 setter，并补充单测。

### 方案 C：作为 memory 注入（不推荐）

把 `LINK.MD` 内容合并到现有 `MemorySection`。

- 优点：不需要新增系统提示词段落。
- 缺点：`LINK.MD` 是用户/项目说明，不是会话记忆；当前 TUI 启动时 `EnableMemory=false`，走 memory 会导致功能不可用或语义混乱。

推荐采用 **方案 B**：它最符合现有系统提示词的模块化结构，改动范围小且语义清楚。

## 4. 文件路径与读取规则

### 4.1 目标文件

- 用户级：`~/.link/LINK.MD`
- 工作区级：`<cwd>/.link/LINK.MD`

其中：

- `~` 通过 `os.UserHomeDir()` 解析。
- `<cwd>` 使用 `config.LoadLocalConfig(cwd)` 已解析和校验后的 `cfg.WorkspaceRoot`。
- 文件名严格为 `LINK.MD`，本次不扩展 `link.md`、`Link.md` 等大小写变体，避免超出需求。

### 4.2 缺失与错误处理

- 文件不存在：忽略，不产生提示词段落中的对应内容。
- 两个文件都不存在：不注入额外 `system-reminder` 段落，系统行为与当前一致。
- 路径存在但不是普通文件、权限不足、读取失败：启动失败，并返回 `read LINK.MD <path>: <reason>` 形式的错误。
- 内容只包含空白：视为不存在，不注入。

### 4.3 内容处理

- 使用 `strings.TrimSpace` 去除首尾空白。
- 不解析 frontmatter，不执行模板，不改写 Markdown 内容。
- 注入时保留原始 Markdown，以便用户在 `LINK.MD` 中维护说明、约束、项目背景等内容。

## 5. 提示词格式设计

当至少读取到一个 `LINK.MD` 时，在系统提示词中追加一个独立段落，建议位置在 `<workspace>` 之后、`<tools>` 之前。原因是这类上下文与当前运行环境强相关，应早于工具/技能说明进入模型上下文。

目标渲染格式：

```xml
<system-reminder>
在回答用户问题时，你可以参考以下上下文：
# linkMd
下面展示了用户与代码库说明。请务必遵循这些说明。重要：这些说明将覆盖任何默认行为，你必须严格按其文字要求执行。

/Users/bytedance/.link/LINK.MD 的内容（用户为所有项目配置的私人全局说明）：

{{user link content}}

/Users/bytedance/project/.link/LINK.MD 的内容（项目说明，已提交到代码库或工作区）：

{{workspace link content}}

# 重要指令提醒

只做被要求的事：不多不少。
除非为达成目标绝对必要，切勿创建新文件。
能修改现有文件，绝不新建文件。
不要主动创建文档文件（*.md）或README。仅当用户明确要求时才创建文档。

重要：这些上下文可能与当前任务相关，也可能无关。除非与任务高度相关，否则不要对其作出回应。
</system-reminder>
```

实际路径由运行环境决定，不使用示例绝对路径。若只存在其中一个文件，只渲染该文件对应小节。

## 6. 模块设计

### 6.1 `internal/config`：路径解析与上下文加载

新增轻量结构体：

```go
type LinkMarkdownContext struct {
    UserPath         string
    UserContent      string
    WorkspacePath    string
    WorkspaceContent string
}
```

新增函数建议：

- `LinkMarkdownPath() (string, error)`：返回 `~/.link/LINK.MD`。
- `WorkspaceLinkMarkdownPath(workspaceRoot string) string`：返回 `<workspaceRoot>/.link/LINK.MD`。
- `LoadLinkMarkdownContext(workspaceRoot string) (LinkMarkdownContext, error)`：读取两级文件并 trim 内容。

实现细节：

- 复用 `LinkSkillsDir()` 和 `WorkspaceLinkSkillsDir()` 的 home/workspace 路径风格：`internal/config/local_config.go:77`、`internal/config/local_config.go:85`。
- 新增内部 helper `readOptionalMarkdown(path string) (string, error)`：不存在返回空串；目录或读取失败返回带路径错误。
- 保持 `config.LoadLocalConfig(cwd)` 只负责配置解析，不在其中读取 `LINK.MD` 内容；实际读取放在 `local.Bootstrap`，因为该上下文属于 runtime prompt 资源，而不是 AppConfig 基础配置。

### 6.2 `internal/local`：启动期读取并注入 runtime

在 `local.Bootstrap` 中，在 `cfg` 解析和 layout 校验完成后读取：

1. 调用 `config.LoadLinkMarkdownContext(cfg.WorkspaceRoot)`。
2. 读取失败则返回 `load LINK.MD context: ...`。
3. 成功后调用 `runtimeService.SetLinkMarkdownContext(linkCtx)`。

放置位置建议在 `runtimeService := runtime.NewService(...)` 之后、`SetBasePrompt` 附近，使所有 prompt 相关初始化集中在一起：`internal/local/bootstrap.go:63` 至 `internal/local/bootstrap.go:68`。

### 6.3 `internal/agent/runtime`：持有上下文并传入 PromptOptions

在 `Service` 中新增字段：

```go
LinkMarkdownContext config.LinkMarkdownContext
```

新增 setter：

```go
func (s *Service) SetLinkMarkdownContext(ctx config.LinkMarkdownContext)
```

在 `buildSystemPromptWithMemory` 调用 `assistant.BuildSystemPrompt` 时传入新字段，位置与现有 `WorkingDirectory`、`ToolNames` 并列：`internal/agent/runtime/prompt_builder.go:40`。

### 6.4 `internal/assistant`：渲染独立系统提醒段落

扩展 `PromptOptions`：

```go
LinkMarkdownContext LinkMarkdownContext
```

也可以为了避免 `assistant` 反向依赖 `config`，在 `assistant` 包内定义只含路径和内容的同名/相近结构，由 `runtime` 做字段映射。推荐避免 `assistant -> config` 依赖，保持 prompt 包纯粹。

新增渲染函数：

- `renderLinkMarkdownContext(ctx LinkMarkdownContext) string`
- `appendMarkdownFileSection(builder, path, label, content)`

`BuildSystemPrompt` 中：

1. 先渲染 `<identity>`。
2. 再渲染 `<workspace>`。
3. 若 link context 非空，追加 `renderTag("system-reminder", body)`。
4. 后续保持 `<tools>`、`<skills>`、`<memory>` 现有顺序。

这样不会影响工具、技能、记忆逻辑，也符合示例提示词的 `system-reminder` 表达。

## 7. 测试设计

### 7.1 `internal/config` 单测

新增或扩展 `internal/config/local_config_test.go`：

- 用户级和工作区级 `LINK.MD` 同时存在时均读取，并返回绝对路径与 trim 后内容。
- 只存在用户级时只返回用户内容。
- 只存在工作区级时只返回工作区内容。
- 两者都不存在时返回空 context 且无错误。
- 路径存在但为目录时返回错误，错误包含路径。

如测试需要控制 home 目录，使用 `t.Setenv("HOME", tempDir)`；在 darwin/linux 下 `os.UserHomeDir()` 会使用该环境变量。

### 7.2 `internal/assistant` 单测

扩展 `internal/assistant/prompt_test.go`：

- `BuildSystemPrompt` 在 link context 非空时包含 `<system-reminder>`。
- 输出包含用户级与工作区级路径、标签和原始 Markdown 内容。
- link context 为空时不输出 `# linkMd`，避免无意义噪音。
- 验证段落顺序：`</workspace>` 出现在 `<system-reminder>` 之前，`</system-reminder>` 出现在 `<tools>` 之前。

### 7.3 `internal/local` 或 runtime 单测

根据现有测试可行性选择其一：

- 优先在 runtime 层测试 `SetLinkMarkdownContext` 后 `buildSystemPromptWithMemory` 能把 context 传入最终系统提示词。
- 如 `local.Bootstrap` 集成测试成本可控，再补充临时工作区 `.link/LINK.MD` 到 runtime prompt 的端到端验证。

## 8. 实施步骤建议

1. 在 `internal/config` 增加 `LINK.MD` 路径函数与可选读取函数，并补充单测。
2. 扩展 `assistant.PromptOptions` 与 `BuildSystemPrompt`，新增 `system-reminder` 渲染，并补充单测。
3. 扩展 runtime `Service` 字段与 setter，在 prompt_builder 中传递上下文。
4. 修改 `local.Bootstrap` 启动流程，读取 `LINK.MD` 并设置到 runtime。
5. 运行相关 Go 测试：至少覆盖 `./internal/config`、`./internal/assistant`、`./internal/agent/runtime`、`./internal/local`。

## 9. 风险与边界

- **大文件风险**：本设计不额外限制 `LINK.MD` 大小，保持与需求一致；如后续发现上下文过大影响 token，可再引入大小限制或摘要策略。
- **内容优先级风险**：提示词中明确说明这些说明会覆盖默认行为，但仍不应覆盖系统/安全约束；实现只负责注入内容，不额外解释或改写用户说明。
- **启动后变更不生效**：TUI 进程启动后修改 `LINK.MD` 不会自动刷新，需要重启进程。这符合“启动进程时自动读取”的需求。
- **工作区解析**：工作区级文件以 `--cwd` 或进程启动 cwd 解析后的 `cfg.WorkspaceRoot` 为准，不以 go-agent 程序所在目录为准。

## 10. 自检结论

- 本设计无待定项；`~./.link` 已明确解释为 `~/.link`。
- 设计只覆盖 `LINK.MD` 启动读取与系统提示词注入，不包含 watcher、配置热更新、大小写兼容等额外功能。
- 推荐方案与现有架构一致：配置包负责路径与读取，local 启动期装配，runtime 持有上下文，assistant 统一渲染系统提示词。
