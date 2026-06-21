# 支持在任意项目目录运行 Agent 设计文档

## 1. 背景与目标

当前 `cynosure` 必须在 `cynosure/` 源码目录下执行 `go run .` 才能启动。根本原因是程序自身的**运行资源**（`config.json`、`system_prompt.md`、内置 `skills/`）通过配置项 `app_home="."` 硬编码，依赖**进程启动目录**读取。

而用户要操作的**工作区**（目标项目目录）其实已经通过 `os.Getwd()` / `--cwd` 解耦，本身没有问题。

**目标**：像 Claude Code 一样，用户安装后在任意项目目录直接执行 `cynosure`，即可对当前目录启动 agent，无需关心源码位置或资源文件的摆放。

### 核心决策
- **资源定位方案**：用 `go:embed` 把 `system_prompt.md` 与内置 `skills/` 嵌入二进制；`config.json` 默认值内置为代码常量。这样二进制自包含，移动到任何位置都能运行。
- **分发方式**：`go install` 安装单一二进制，命令名 `cynosure`，在任意项目目录直接调用。

### 约束（注意事项）
- 不破坏现有功能：工作区解耦逻辑、双层 Skills、MCP、记忆系统、会话恢复、安全边界保持不变。
- 本设计文档先交用户审阅，审阅通过后再实施代码改动。

---

## 2. 现状分析

| 关注点 | 代码位置 | 现状 |
| --- | --- | --- |
| 入口 / cwd | `internal/cli/root.go:50-56` | `os.Getwd()` 作默认 cwd，`--cwd` 可覆盖（**已 OK**） |
| config.json 读取 | `internal/config/config.go:50-65` | `configFilePath()` 硬编码返回 `"config.json"`（相对启动目录） |
| app_home 解析 | `internal/config/paths.go:25-32` | 默认 `"."`，所有运行资源以此为基准 |
| system_prompt | `internal/assistant/prompt.go:35-41` | 从 `cfg.SystemPromptPath` 读文件，缺失即启动失败 |
| 内置 skills | `internal/local/bootstrap.go:48-58` | 仅读用户级 `~/.cynosure/skills` + 工作区级，**实际未加载** app_home 下的内置 skills |
| logs | `paths.go:73`、`bootstrap.go:45` | 写到 `<app_home>/logs` |
| bin/cmd 资产 | `paths.go:78-95`、`tools/runtime_env.go:57-69` | `<app_home>/bin`、`<app_home>/cmd`，供 bash 工具白名单 |
| 用户级目录 | `internal/config/local_config.go:76-122` | settings/skills/memory/session/task_outputs 已统一在 `~/.cynosure/`（**已 OK**） |

**关键结论**：耦合点集中在 `app_home`。`config.json` 的 `app_home / builtin_skills_dir / command_bin_dir / command_script_dir` 字段，以及 `EnsureAppLayout / ValidateAppLayout` 的目录创建与校验，是"必须在源码目录运行"的根因。

---

## 3. 设计方案

### 3.1 嵌入静态资源（新增 `assets` 包）

`go:embed` 只能嵌入 `.go` 文件所在目录及其子目录，不能使用 `..` 向上引用。当前 `system_prompt.md` 与 `skills/` 位于 `cynosure/` 根目录（与 `main.go` 同级）。

方案：新建目录 `cynosure/assets/`，将 `system_prompt.md` 与 `skills/` 移入，新增 `assets/embed.go` 执行嵌入。

```go
package assets

import "embed"

//go:embed system_prompt.md
var systemPrompt string

//go:embed skills
var skillsFS embed.FS

func SystemPrompt() string { return systemPrompt }
func BuiltinSkillsFS() fs.FS { /* 返回 skills 子树 */ }
```

需同步更新引用路径：`build.sh`、`README.md`、`architecture_test.go`（其 `WalkDir` 遍历源码树）。

### 3.2 config 默认值内置，去掉对 config.json 文件的硬依赖

- 在 `internal/config/config.go` 定义内置默认 `fileConfig`。注意：真实 LLM 配置（key/model/baseURL）仍来自 `~/.cynosure/settings.json`，此处不变，仅托管运行行为类默认值（如 `allowed_tools`、安全开关）。
- `loadConfigFile()` 改为：
  - 优先读 `~/.cynosure/config.json`（用户可选覆盖）；
  - 文件不存在则返回内置默认值，不再报错；
  - 移除 `configFilePath()` 中相对路径 `"config.json"`。
- `app_home` 概念退化：不再用于定位 system_prompt / skills（改用 embed）。

### 3.3 运行期可写目录统一到 `~/.cynosure/`

- `logs`：从 `<app_home>/logs` 改为按工作区 + 会话隔离的 **`~/.cynosure/logs/<workspace>/{session_id}`**。
  - `<workspace>` 与记忆目录保持一致，默认使用工作区目录名；若同名工作区已存在，则按首次登记顺序追加 `_1`、`_2` 等后缀，并由 `config.WorkspaceName(workspaceRoot)` 统一分配。
  - `{session_id}` 即会话 UUID。由于 logger 初始化早于 session_id 生成，需在 `bootstrap.go` 中**将 session_id（`idgen.UUID()`）的生成提前到 logger 初始化之前**，再用该 id 计算 logs 目录。
  - 新增 `config.CynosureSessionLogsDir(workspaceRoot, sessionID)` 返回上述路径。
- `bin`、`cmd`（bash 工具系统资产白名单）：保留能力，默认基准从 app_home 改为 `~/.cynosure/bin`、`~/.cynosure/cmd`（新增对应路径函数）；目录为空不影响运行。
- `EnsureAppLayout` 仅创建：`~/.cynosure/`（含 bin/cmd 按需）+ 当前 workspaceRoot；logs 目录在 logger 初始化时按 session 路径创建；移除对 app_home / builtin_skills_dir 的创建与校验。

### 3.4 内置 skills 接入加载链

- `bootstrap.go` 加载 skills 时新增"内置（builtin/embedded）"来源，从 `assets.BuiltinSkillsFS()` 加载，与用户级、工作区级合并。
- `internal/sessions/skill.go` 当前用 `filepath.WalkDir` + `os.ReadFile`。新增基于 `fs.FS` 的加载函数（`fs.WalkDir` + `fs.ReadFile`），用于嵌入资源。
- 覆盖优先级：**workspace > user > builtin**（与现有合并语义一致，后加载者覆盖同名）。

### 3.5 system_prompt 改为读嵌入内容

- `assistant.LoadBaseSystemPrompt`：若无外部覆盖文件，则使用 `assets.SystemPrompt()`。
- 保留用户可选覆盖：`~/.cynosure/system_prompt.md` 存在时优先使用。
- `bootstrap.go:59` 相应调整：优先用户覆盖文件 → 否则用嵌入内容；不再因缺文件而启动失败。

### 3.6 CLI / 帮助 / 安装

- `internal/cli/root.go`：`--cwd` 行为不变；更新 `printHelp` 文案，说明可在任意目录运行。
- `build.sh`：不再需要拷贝 `system_prompt.md` / `skills/`（已 embed），保留 config 示例拷贝。
- `README.md`：新增 `go install`（命令名 `cynosure`）安装说明；更新"启动 TUI"章节为在任意项目目录直接执行 `cynosure`。

---

## 4. 需要改动的文件

| 文件 | 改动 |
| --- | --- |
| `cynosure/assets/embed.go`（新增） | 嵌入 system_prompt 与 skills |
| `cynosure/assets/system_prompt.md`、`cynosure/assets/skills/`（移动） | 从根目录移入 assets 包 |
| `internal/config/config.go` | 内置默认值 + 用户级 config.json 可选覆盖 |
| `internal/config/local_config.go` | 新增 `CynosureLogsDir/BinDir/ScriptDir` 等路径函数 |
| `internal/config/paths.go` | app_home 解耦，运行期目录改到 `~/.cynosure/` |
| `internal/local/bootstrap.go` | 接入 embedded skills + embedded system_prompt + 新 logs 路径 |
| `internal/assistant/prompt.go` | system_prompt 默认走 embed |
| `internal/sessions/skill.go` | 支持从 `fs.FS` 加载内置 skills |
| `internal/cli/root.go` | 帮助文案 |
| `build.sh`、`README.md` | 分发与启动说明更新 |
| 相关测试 | `internal/config`、`internal/local`、`internal/sessions`、`architecture_test.go` 同步更新 |

---

## 5. 风险与兼容性

- **embed 路径**：资源文件必须位于包目录内，移动 `system_prompt.md` / `skills/` 后需同步所有引用与构建脚本，避免遗漏导致编译失败。
- **旧 config.json 字段弱化**：`app_home / builtin_skills_dir` 等字段保留解析以兼容旧文件，但不再强制要求；行为以内置默认值为准。
- **行为保持不变**：工作区安全边界、记忆 / 会话 / MCP 加载逻辑完全不动。

---

## 6. 验证方式

1. `cd cynosure && go build ./...` 通过；`go test ./...` 全绿（重点 config / local / sessions / architecture）。
2. `go run .`（在 cynosure 目录）仍能启动，行为不变。
3. 模拟安装：`go install .` 后，`cd /任意项目目录 && cynosure`，确认：
   - 启动成功，工作区 = 当前目录；
   - system prompt、内置 skills 正常加载（`/skills` 可见 builtin 来源）；
   - 文件 / bash 工具边界仍限制在当前目录；
   - logs 写入 `~/.cynosure/logs/<workspace>/{session_id}`。
4. `cynosure --cwd /other/project` 指定其他目录仍正常。

---

## 7. 交付节奏

1. **当前阶段**：本设计文档交用户审阅。
2. **审阅通过后**：按第 3、4 节实施代码改动，并按第 6 节验证。
