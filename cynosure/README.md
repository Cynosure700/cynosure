# cynosure / nano_cc

`cynosure` 是面向本地代码工作的纯本地 TUI Agent。它默认在终端中启动，不依赖 Web 服务、数据库、Redis 或 Elasticsearch；能够读取当前项目文件、执行 shell 命令、修改代码、加载 `cynosure` Skills，并自动连接工作区 MCP 工具。

## 核心能力

- **TUI 聊天界面**：终端内多轮对话、Claude Code 风格布局、流式输出、Markdown 渲染、实时状态栏和基础 slash commands。
- **本地代码工具**：支持 `bash`、`read_file`、`write_file`、`edit_file`、`multi_edit`、`grep`、`glob`、`ls`、`web_fetch`、`web_search`、`todo_write`、`todo_list`、`load_skill`、`spawn_subagent`、`read_persisted_output`。
- **工作区安全边界**：默认工作区是启动命令所在目录或 `--cwd` 指定目录；文件和 shell 工具默认不能越过工作区。
- **命令权限审批**：查询/搜索类工具直接执行；`bash`、`write_file`、`edit_file`、`multi_edit` 等变更类操作在执行前向用户申请权限，用户可选择放行一次、放行并记住规则、或拒绝；拒绝时立即结束本轮且不做记忆与历史收尾。
- **Skill 系统**：启动时读取 `~/.cynosure/skills` 与 `<cwd>/.cynosure/skills`；模型可通过 `load_skill` 按需加载正文。
- **工作区 MCP**：启动时读取 `<cwd>/.cynosure/.mcp.json` 并自动连接；发现到的工具以 `mcp__{server}__{tool}` 形式加入模型工具列表。
- **项目级记忆**：启动时读取 `~/.cynosure/memory/<workspace-key>/memory.md`，由模型按当前对话筛选有用记忆；每条长期记忆是一个独立 Markdown 文件，仅对当前项目有效。
- **历史会话恢复**：对话历史持久化到 `~/.cynosure/session/{session_id}/`，可在同一项目目录通过 `/resume` 选择恢复。
- **上下文压缩**：请求模型前自动压缩超长历史，支持大工具结果落盘、消息窗口裁剪、最近工具结果保留与 413 后激进压缩；消息窗口裁剪在消息总数超过 50 条时，从队首保留到“用户最新消息 + 其后 2 条”并补足尾部至合计 50 条，避免裁掉用户最新提问；全量摘要采用结构化提示词（用户最新问题、改动文件路径、关键决策、进展与待办、报错与命令结论），并保留最近 5 条原文，最终结构为 `[系统提示 + 摘要 + 最近 5 条]`；上下文摘要仅保存在运行期内存中，不做持久化。

## 项目结构

```text
cynosure/
├── main.go                  # TUI 默认入口
├── build.sh                 # 发布包构建脚本
├── assets/                  # go:embed 嵌入资源
│   ├── embed.go             # 嵌入声明
│   ├── system_prompt.md     # 基础 identity prompt（嵌入二进制）
│   └── skills/              # 内置 Skill（嵌入二进制）
└── internal/
    ├── agent/               # Agent runtime、MCP 与存储模型
    ├── assistant/           # system prompt 拼装
    ├── cli/                 # 命令行参数解析
    ├── config/              # 配置与 runtime 路径
    ├── idgen/               # ID 生成
    ├── llm/                 # OpenAI 兼容 LLM 客户端与错误分类
    ├── local/               # TUI bootstrap 与本地内存 Store
    ├── logger/              # 日志
    ├── safety/              # 路径安全
    ├── sessions/            # Skill loader / merge / render
    ├── textutil/            # 文本工具
    ├── tools/               # 内置工具定义、参数校验与执行
    └── tui/                 # Bubble Tea TUI 界面与事件桥接
```

## 环境要求

- Go 1.26.1（见 `go.mod`）
- OpenAI 兼容模型服务
- 用户级 `~/.cynosure/settings.json`

TUI 模式不包含 MySQL、Redis、Elasticsearch 或 Web 服务依赖。

## 配置

### 读取规则

- 运行配置默认值内置在代码中；若存在 `~/.cynosure/config.json`，其内容会覆盖默认值（可选）。配置用于工具白名单和安全开关等非敏感项。
- LLM 密钥、模型和 Base URL 从 `~/.cynosure/settings.json` 读取。
- `system_prompt.md` 与内置 skills 通过 `go:embed` 嵌入二进制；若存在 `~/.cynosure/system_prompt.md`，则优先使用该用户覆盖文件。
- 内置 skills 来源为 `builtin`，与用户级 `~/.cynosure/skills` 和工作区级 `<cwd>/.cynosure/skills` 合并，优先级 workspace > user > builtin。
- TUI 的工作区为启动 cwd，或由 `--cwd` 指定。
- 本地配置不包含数据库、Redis、Elasticsearch、JWT 或 Cookie 会话字段。

启动时会自动创建：`~/.cynosure/`、按会话隔离的日志目录和当前工作区目录。

### `~/.cynosure/settings.json` 示例

```json
{
  "env": {
    "open_auth_token": "your-api-key",
    "open_model": "deepseek-v4-flash",
    "open_base_url": "https://api.deepseek.com"
  }
}
```

`env.open_auth_token`、`env.open_model`、`env.open_base_url` 任一缺失都会导致 TUI 启动失败。密钥不会写入日志或 TUI 展示。

### 本地存储边界

- **项目级文件**：`<cwd>/.cynosure/skills` 与 `<cwd>/.cynosure/.mcp.json` 分别提供工作区 Skills 与 MCP 配置；`<cwd>/.cynosure/settings.json` 保存工作区命令权限配置（`permissions.defaultMode` 与 `permissions.allowedRules`）。
- **用户级文件**：`~/.cynosure/settings.json` 保存 LLM 配置，`~/.cynosure/skills` 保存用户级 Skills，`~/.cynosure/memory/` 保存项目记忆，`~/.cynosure/task_outputs/` 保存工具输出日志和大结果落盘文件，`~/.cynosure/session/{session_id}/` 保存历史会话。
- **运行期内存**：本地 Store 使用内存 map 和锁维护当前进程状态；上下文摘要只保存在内存中，进程退出后不恢复。

### `~/.cynosure/config.json` 示例（可选）

不创建该文件时使用内置默认值。需要自定义工具白名单或安全开关时才创建：

```json
{
  "allowed_tools": "load_skill,bash,read_file,write_file,edit_file,multi_edit,grep,glob,ls,web_fetch,todo_write,todo_list,spawn_subagent",
  "bash_allow_outside_workspace": false,
  "bash_allow_dangerous_commands": false
}
```

## 项目级记忆

TUI 本地模式会在用户目录创建并维护 `~/.cynosure/memory/<workspace-key>/` 目录。`<workspace-key>` 由当前工作区目录名和工作区绝对路径 hash 组成；记忆只对当前项目下的会话有效，切换到其他项目时不会复用当前项目记忆。

```text
~/.cynosure/memory/<workspace-key>/
├── memory.md                 # 长期记忆索引
├── <memory-name>.md          # 单条长期记忆
└── sessions/
    └── <session_id>.md       # 当前会话记忆，按随机 UUID session_id 覆盖更新
```

- `memory.md` 只保存记忆文件位置、名称和描述；模型会先基于索引判断哪些记忆对当前轮对话有用，再注入选中记忆正文。
- 长期记忆分为四类：`preference`（用户长期偏好、习惯与项目相关描述）、`feedback`（对 Agent 行为的纠正与肯定，body 含 `Why:` 与 `How to apply:`）、`project`（项目进展、决策、截止日期等不可从代码推导的动态，相对日期转为绝对日期）、`reference`（外部系统中信息的定位信息）。抽取时不会记录可从代码、Git 历史、调试上下文或 `CYNOSURE.MD` 直接获得的内容。
- 当前会话记忆使用随机 UUID `session_id` 标识，同一会话每轮结束后覆盖更新 `~/.cynosure/memory/<workspace-key>/sessions/<session_id>.md`，不会按轮次生成多个文件。
- 记忆提取（长期记忆与会话记忆）基于**模型线**而非纯文本对话：渲染时包含完整交互——用户与助手文本、助手发起的工具调用（名称与参数）、工具结果（状态与内容），使记忆扎根于真实发生的交互。轮末提取与模型历史落库共用本轮压缩后的请求线（`lastRequestHistory` + 本轮最终 assistant）；循环中途提取则取 `state.ModelHistory`（会话循环内逐字 lockstep 追加的 user / assistant / tool 消息）。
- 会话收尾更新只对当前项目当前会话加锁，锁 key 为 `项目名 + session_id`。

## 历史会话

TUI 会将每个会话的展示历史和模型历史分别写入用户目录下的 `~/.cynosure/session/{session_id}/`：

```text
~/.cynosure/session/
└── <session_id>/
    ├── history         # 完整展示历史，用于恢复用户与 agent 的交互记录
    └── model_history   # 压缩后的模型上下文，恢复后继续作为 LLM 上下文基线
```

- `history` 保存完整会话消息，恢复展示时默认只渲染 user / assistant / system / error 消息，不直接展开 tool 大结果。
- `model_history` 保存上一轮压缩后的上下文；缺失时会回退到 `history`，保持现有上下文加载逻辑。
- 上下文摘要不单独写入历史文件；恢复会话时不会加载旧摘要。
- 历史文件内记录 `workspace_root`，`/resume` 只展示当前工作区对应的历史会话，避免跨项目误恢复。

## 工具输出落盘

TUI 会在用户目录维护 `~/.cynosure/task_outputs/`，用于保存工具执行审计记录和被上下文压缩移出模型上下文的大工具结果：

```text
~/.cynosure/task_outputs/
├── tool-results/
│   ├── <session_id>-<persisted_output_id>.txt   # 大工具结果全文
│   └── <session_id>-<persisted_output_id>.json  # 结果 metadata 与 sha256
└── <session_id>/
    └── tools.md                                 # 当前会话工具执行结果追加日志
```

- 当最近一轮 `tool_result` 总量超过预算时，系统会从最大的结果开始落盘到 `~/.cynosure/task_outputs/tool-results/`，模型上下文中只保留 `<persisted-output>` 标记和前 2000 字符预览。
- 模型如需读取完整结果，会通过 `read_persisted_output(id, offset, limit)` 分段读取，读取时会校验会话、用户、conversation 和 sha256。
- 每次工具执行完成后，都会向 `~/.cynosure/task_outputs/{session_id}/tools.md` 追加工具名、参数、状态、审计摘要和结果内容，便于本地排查；该文件不会自动注入模型上下文。

## 安装与启动

`cynosure` 是自包含的单一二进制：`system_prompt.md` 与内置 skills 通过 `go:embed` 嵌入二进制，运行配置默认值内置在代码中，无需随包分发资源文件。安装后可在任意项目目录直接运行。

### 安装（推荐）

```bash
go install nano_cc@latest    # 上传 GitHub 后用对应 module 路径
# 或在源码目录：
cd cynosure && go install .
```

安装后 `cynosure` 进入 `GOBIN`（默认 `~/go/bin`，需在 `PATH` 中）。

### 启动 TUI

```bash
# 在任意项目目录下，默认工作区为当前目录
cd /path/to/your/project
cynosure

# 显式指定用户项目目录作为工作区
cynosure --cwd /path/to/project
```

### 开发期运行

```bash
cd cynosure
go run .                      # 默认工作区为当前 shell 所在目录
go run . --cwd /path/to/project
```

## TUI 界面与实时状态

TUI 主界面由顶部状态栏、会话视窗和输入区组成：

- **顶部状态栏**：展示当前运行状态、当前回复已使用工具数、上下文 token 使用比例、当前工作区、Skill 数量和 MCP 工具数量。
- **会话视窗**：按角色区分用户、Agent、系统和错误消息；Agent 回复支持 Markdown 渲染。
- **思考过程**：模型流式返回的 reasoning 会以灰色 `✦ 思考` 块实时展示，降低视觉权重但保留执行透明度。
- **流式回答**：最终对用户可见的 assistant 内容会继续按增量输出；涉及工具调用的中间 assistant 文本不会提前刷屏，避免干扰阅读。
- **输入区提示**：底部固定展示 `Enter`、`Ctrl+C`、`/help` 等常用操作提示。

状态栏示例：

```text
✦ cynosure  generating  工具 3  上下文 72% · 72k/100k
cwd /path/to/project · skills 6 · mcp tools 4
```

其中 `工具` 是当前回复累计工具调用次数，`上下文` 是最近一次发送给模型的上下文估算 token 与预算占比。

## TUI 快捷命令

- `Enter`：发送当前输入。
- `Ctrl+C`：生成中断当前请求；空闲时退出。
- `/help`：显示帮助。
- `/clear`：清空当前 TUI 显示上下文。
- `/cwd`：显示当前工作区。
- `/skills`：显示启动时加载的 Skill 名称、来源、描述和路径。
- `/mcp`：显示工作区 MCP server 状态、连接信息、工具数量和错误摘要。
- `/resume`：展示当前工作区可恢复的历史会话列表，输入序号后恢复所选会话；输入 `/cancel` 取消选择。

命令权限审批面板出现时，使用 `↑/↓` 移动光标、`Enter` 确认所选项，或直接按 `1`/`2`/`3` 选择放行/记住放行/拒绝，`Ctrl+C`、`Esc` 等同于拒绝。

## 系统提示词

每轮请求模型前，`assistant.BuildSystemPrompt` 会把基础 identity 提示词与运行期动态段落拼成最终 system prompt：

- `<identity>`：基础提示词，来自嵌入二进制的 `assets/system_prompt.md`，若存在 `~/.cynosure/system_prompt.md` 则优先使用该覆盖文件。内容参考 `cn/system.md` 的分域设计，并按 Cynosure 当前能力划分为 `定义角色`、`安全性声明`、`帮助文档`、`输出风格`、`任务管理`、`工具调用`、`环境信息` 七类；其中工具与环境部分只描述 Cynosure 当前真实存在的能力和运行期动态注入边界。
- `<workspace>`：注入 Surface 与当前工作目录。
- `<system-reminder>`：在存在 `CYNOSURE.MD` 时注入用户级与工作区级说明，以及重要指令提醒。
- `<tools>`：注入本次会话实际可用的工具名清单（含 MCP 工具）。
- `<skills>`：注入加载到的 Skill 摘要，模型按需 `load_skill` 加载正文。
- `<memory>`：开启记忆时注入按当前对话筛选出的项目记忆。

基础提示词只描述行为约束，不写死工具列表、工作目录、Skill 等动态信息——这些由上述段落在运行期注入。`任务管理` 章节强调在复杂或多步骤软件工程任务中频繁使用 `todo_write` 规划、拆解、逐项更新和验证，并在上下文裁剪/压缩后用 `todo_list` 查询当前待办状态；`环境信息` 章节会要求模型以运行期 `<workspace>`、`<tools>`、`<skills>`、`<memory>` 与 `<system-reminder>` 为准，避免把静态提示词当成固定工作区配置。相关代码见 `internal/assistant/prompt.go` 与 `internal/agent/runtime/prompt_builder.go`。

基础提示词会要求模型只使用 `<tools>` 中实际列出的工具：内容搜索优先 `grep`、文件名匹配优先 `glob`、目录浏览使用 `ls`，文件读取和修改分别使用 `read_file`、`edit_file`、`multi_edit`、`write_file`；`bash` 仅用于确需 shell 的操作并遵循审批与工作区边界；`todo_write` 用于维护任务清单，`todo_list` 用于只读查询当前任务状态；`spawn_subagent` 仅用于隔离子任务；`read_persisted_output` 仅用于读取上下文压缩产生的落盘结果。

## 工具系统

| 工具 | 说明 |
| --- | --- |
| `load_skill` | 加载 Skill 正文与元数据 |
| `bash` | 在工作区内执行 shell 命令，受越权路径与危险命令限制，执行前需用户审批 |
| `read_file` | 读取工作区内文件 |
| `write_file` | 写入工作区内文件，执行前需用户审批 |
| `edit_file` | 按精确文本替换编辑文件，执行前需用户审批 |
| `multi_edit` | 对单个文件一次性顺序执行多处查找替换，全部成功才写盘，执行前需用户审批 |
| `grep` | 纯 Go 正则内容搜索，支持 `files_with_matches`/`content`/`count` 三种输出模式 |
| `glob` | 文件名 glob 模式匹配，结果按修改时间倒序返回 |
| `ls` | 列出指定绝对路径下的文件与目录，支持 `ignore` glob 过滤 |
| `web_fetch` | 抓取 URL，将 HTML 转为文本后用 LLM 按 prompt 处理；未配置 LLM 时返回清洗文本 |
| `web_search` | 联网搜索（占位）；未配置搜索后端时返回提示，默认不启用 |
| `todo_write` | 维护多步骤任务清单 |
| `todo_list` | 查询当前任务清单状态，不修改任务 |
| `spawn_subagent` | 派生隔离子 Agent |
| `read_persisted_output` | 自动暴露，用于读取上下文压缩落盘的大工具结果 |

安全边界：

- 文件读写、编辑、内容/文件名检索、目录列举和命令执行都以 TUI 工作区为默认边界。
- `grep`/`glob` 的 `path` 默认是当前工作目录，会被约束在工作区内；`ls` 要求传入绝对路径且同样受工作区约束。
- `web_fetch` 会将 `http://` 升级为 `https://`，限制响应体大小并复用普通工具的 60 秒看门狗；`web_search` 当前为占位实现。
- `bash_allow_outside_workspace=false` 时拒绝访问工作区外路径。
- `bash_allow_dangerous_commands=false` 时拒绝危险命令。

超时机制：

- **主 Agent 单轮会话**：上限 24 小时，在会话循环每轮入口做截止时间检查的软边界，超时后结束本轮会话。
- **子 Agent 单轮会话**：上限 1 小时，同样为循环间的截止时间检查。
- **终端工具（`bash`）**：上限 120 秒，超时通过 `time.AfterFunc` + `Process.Kill()` 看门狗终止子进程。
- **其他普通工具**（`read_file`/`write_file`/`edit_file`/`multi_edit`/`grep`/`glob`/`ls`/`web_fetch`/`web_search`/`load_skill`/`todo_write`/`todo_list`/`read_persisted_output`）：上限 60 秒，在 `Dispatch` 层以 goroutine + `time.After` 看门狗兜底。
- 以上阈值均为代码内硬编码常量，超时机制不依赖 `context`，不影响现有 `ctx` 透传链路。

## 命令权限审批

为避免变更类操作未经确认即执行，TUI 在工具执行前引入交互式权限审批：

- **免审批（查询/搜索）**：`read_file`、`grep`、`glob`、`ls`、`web_search`、`web_fetch`、`load_skill`、`todo_write`、`todo_list`、`read_persisted_output`、`spawn_subagent` 等工具直接执行，无需确认。
- **需审批（变更类）**：`bash`（含 `curl`、写入、删除等任意命令）以及 `write_file`、`edit_file`、`multi_edit` 在执行前弹出审批面板，提供三个选项：
  - `1. Yes`：本次放行。
  - `2. Yes, and don't ask again for: <rule>`：本次放行并把放行规则写入工作区 `settings.json`，后续同类操作不再询问。`bash` 规则按命令名通配（如 `curl *`），写工具规则按工具名通配（如 `write_file *`）。
  - `3. No`：拒绝。该轮对话立即结束，且**不进行记忆与对话历史的收尾落库**。
- **判定优先级**：是否需审批 → `bypassPermissions` 全放行 → 命中已放行规则 → 弹出审批面板。权限配置在**每次审批时实时读取**工作区 `settings.json`，运行中修改即时生效。
- **子 Agent 一致**：`spawn_subagent` 派生的子 Agent 执行变更类工具时走同一套审批逻辑，共用工作区放行规则与 bypass 模式；子 Agent 中被拒绝会使该次子任务失败并结束父轮。

### `<cwd>/.cynosure/settings.json` 示例（可选）

```json
{
  "permissions": {
    "defaultMode": "bypassPermissions",
    "allowedRules": ["curl *", "write_file *"]
  }
}
```

- `permissions.defaultMode` 为 `bypassPermissions` 时，一切命令直接放行、不再审批；其它取值或缺省时走正常审批流程。
- `permissions.allowedRules` 是已放行规则列表，命中则免审批；选择审批选项 2 时会自动追加去重写入。
- 不创建该文件或缺省 `permissions` 时，对 `bash` 与写工具均会审批（默认安全行为）。

## Skill 系统

- 用户级 Skill 从 `~/.cynosure/skills` 加载。
- 工作区级 Skill 从 `<cwd>/.cynosure/skills` 加载。
- 支持 `skill.md` 和 `SKILL.md` 作为 Skill 入口文件。
- 工作区级同名 Skill 覆盖用户级 Skill；同一级目录内重复 Skill 名称会导致启动失败。
- 每轮对话生成 Skill snapshot，system prompt 与 `load_skill` 工具使用同一份 snapshot。
- system prompt 只注入 Skill 摘要；模型需要正文时通过 `load_skill` 按需加载。

示例：

```text
~/.cynosure/skills/code-review/skill.md
<cwd>/.cynosure/skills/project-helper/SKILL.md
```

相关代码：

- `internal/sessions/skill.go`
- `internal/agent/runtime/prompt_builder.go`
- `internal/tools/load_skill.go`

## MCP

TUI 启动时会读取当前工作区 `<cwd>/.cynosure/.mcp.json`。文件不存在时跳过工作区 MCP；文件存在但格式错误时启动失败。单个 MCP server 连接失败不会阻断启动，可通过 `/mcp` 查看失败原因。

支持 `mcp_servers` 数组格式：

```json
{
  "mcp_servers": [
    {
      "name": "filesystem",
      "transport": "stdio",
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "."],
      "env": {},
      "enabled": true
    }
  ]
}
```

也支持 `mcpServers` map 格式：

```json
{
  "mcpServers": {
    "docs": {
      "transport": "sse",
      "url": "https://example.com/sse",
      "headers": {
        "Authorization": "Bearer token"
      },
      "enabled": true
    }
  }
}
```

规则：

- `transport` 支持 `stdio`、`sse`、`streamable`。
- `transport` 为空且配置了 `command` 时，默认为 `stdio`。
- `enabled` 为空时默认为 `true`。
- MCP 工具对 TUI 会话可用，工具名格式为 `mcp__{server}__{tool}`。

## 开发与测试

```bash
cd cynosure
go test ./...
```

重点测试包：

- `internal/cli`：入口分发。
- `internal/config`：本地配置加载。
- `internal/local`：本地内存 Store。
- `architecture_test.go`：防止重新引入 `internal/web` 和数据库 / Redis / Elasticsearch 依赖。
- `internal/tui`：TUI 事件桥接与 slash commands。

## 设计文档

TUI 化设计文档位于：

```text
TUI化改造设计文档.md
```

超时机制设计文档位于：

```text
docs/超时机制设计文档.md
```

命令权限审批机制设计文档位于：

```text
docs/终端命令权限审批机制设计文档.md
```

上下文压缩优化设计文档位于：

```text
docs/上下文压缩优化设计文档.md
```

会话记忆提取与更新机制（含提取数据源切换为模型线）设计文档位于：

```text
docs/会话记忆提取与更新机制调整设计文档.md
```
