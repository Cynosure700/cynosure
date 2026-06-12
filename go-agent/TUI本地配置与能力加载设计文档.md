# TUI 本地配置与能力加载设计文档

## 1. 背景与目标

当前 go-agent 已迁移到本地 TUI 启动链路，入口经 `main.go` 调用 CLI，再由 `local.Bootstrap` 初始化配置、runtime、skills、MCP 与会话。现状中：

- LLM 的 `base_url`、`model_id` 来自项目根目录 `config.json`，API Key 来自环境变量 `OPENAI_API_KEY`，对应代码在 `internal/config/local_config.go:53`。
- TUI 启动只加载 `config.json` 中 `builtin_skills_dir` 指向的单一路径，加载发生在 `internal/local/bootstrap.go:47`。
- skill 文件扫描固定识别 `SKILL.md`，逻辑在 `internal/sessions/skill.go:24` 与 `internal/sessions/skill.go:51`。
- `load_skill` 工具通过 `SkillSnapshot` 从用户侧 loader 与本地 loader 中按需读取 skill 全文，逻辑在 `internal/tools/load_skill.go:44`。
- MCP 当前支持内置配置 `app_home/mcp_config.json`，启动时读取并建立内置连接，逻辑在 `internal/local/bootstrap.go:61` 与 `internal/agent/mcp/config.go:29`。
- TUI 已有 `/skills`、`/mcp` 命令，但只展示数量，逻辑在 `internal/tui/app.go:170`。

本次改造目标：

1. LLM 密钥、模型、BASE_URL 从用户级 `~/.link/settings.json` 的 `env.open_auth_token`、`env.open_model`、`env.open_base_url` 读取。
2. 启动时自动扫描用户级 `~/.link/skills` 与当前工作目录 `.link/skills` 下的 `skill.md` / `SKILL.md`，只把描述注入系统提示，模型按需通过 `load_skill` 加载全文。
3. 启动时自动读取当前工作目录 `.link/.mcp.json`，解析 MCP 配置并自动建立连接。
4. 用户输入 `/skills` 与 `/mcp` 时展示已加载 skill 与 MCP 连接/工具的详细信息。

> 说明：需求中的 `~./.link/settings.json` 按常见 home 目录语义理解为 `~/.link/settings.json`。实现时会统一使用 `os.UserHomeDir()` 拼接，避免相对路径歧义。

## 2. 设计原则

- **本地优先**：TUI 模式不再依赖 Web 时代的用户数据库配置，优先读取用户 home 与当前工作目录下的 `.link` 配置。
- **按需加载**：系统提示只注入 skill name、description、source、path 等摘要信息，不注入 skill 正文；正文只在模型调用 `load_skill` 时读取。
- **双层覆盖**：用户级配置提供全局默认能力，工作区级配置提供项目定制能力。相同 skill 名称以工作区级覆盖用户级。
- **容错启动**：缺失的 `~/.link/skills`、`<cwd>/.link/skills`、`<cwd>/.link/.mcp.json` 不应导致启动失败；存在但格式错误的配置应返回明确错误。
- **信息可见**：TUI 状态栏保留数量概览，`/skills` 与 `/mcp` 输出足够定位问题的详情。

## 3. 方案比较

### 方案 A：在现有 `BuiltinSkillsDir` 与内置 MCP 上叠加新路径

继续保留 `config.json` 的 `builtin_skills_dir` 与 `app_home/mcp_config.json`，额外加载 `~/.link/skills`、`<cwd>/.link/skills` 与 `<cwd>/.link/.mcp.json`。

- 优点：改动较少，兼容现有 `skills` 与内置 MCP。
- 缺点：TUI 本地语义会混杂旧 app home 配置；同名覆盖关系更复杂，`/skills` 来源展示也更难解释。

### 方案 B：TUI 模式切换为 `.link` 本地能力源（推荐）

TUI 启动仍保留 `config.json` 中非敏感运行时配置（如 system prompt、允许工具、bash 安全开关），但 LLM 配置、skills、工作区 MCP 全部从 `.link` 约定读取：

- LLM：`~/.link/settings.json`。
- Skills：`~/.link/skills` + `<cwd>/.link/skills`。
- MCP：`<cwd>/.link/.mcp.json`。

- 优点：符合用户需求，边界清晰；本地 TUI 与 Web/数据库能力解耦；便于 `/skills`、`/mcp` 解释来源。
- 缺点：需要新增本地配置解析与本地 MCP 配置解析，测试覆盖较多。

### 方案 C：把 `.link` 配置导入内存 Store，再复用现有用户技能/MCP Store 接口

启动时把 `.link` skill 与 MCP 转换为 `storage.Skill`、`storage.MCPServer` 注入 `local.Store`，runtime 继续通过 `ListEnabledSkillsByUser`、`ListEnabledMCPServersByUser` 获取。

- 优点：runtime 改动少，复用现有用户维度懒加载接口。
- 缺点：skill 文件路径、source、优先级在转换后不直观；`load_skill` 若只拿到 DB 形态内容，会弱化“按需从文件加载”的语义。

推荐采用 **方案 B**：它最贴近 TUI 本地化后的约定，且能明确地区分用户级与工作区级能力；MCP 连接由 manager 直接管理 `.link` 配置，不再经过本地内存 Store 伪装成数据库用户配置。

## 4. 目标目录与配置格式

### 4.1 用户设置文件

路径：`~/.link/settings.json`

目标字段：

```json
{
  "env": {
    "open_auth_token": "sk-...",
    "open_model": "deepseek-v4-flash",
    "open_base_url": "https://api.deepseek.com"
  }
}
```

读取规则：

- `env.open_auth_token` -> `config.Config.APIKey`
- `env.open_model` -> `config.Config.ModelID`
- `env.open_base_url` -> `config.Config.BaseURL`
- 三者任一为空时，TUI 启动失败，并提示缺失字段与文件路径。
- 不把 token 写入日志、TUI、错误详情或测试快照。

### 4.2 Skill 目录

用户级路径：`~/.link/skills`

工作区级路径：`<cwd>/.link/skills`

支持结构：

```text
~/.link/skills/
  code-review/
    skill.md
  writer/
    SKILL.md

<cwd>/.link/skills/
  project-helper/
    skill.md
```

读取规则：

- 扫描 `skill.md` 与 `SKILL.md`，大小写兼容。
- frontmatter 中 `name` 作为 skill 名称；没有 `name` 时使用父目录名。
- frontmatter 中 `description` 用于系统提示与 `/skills` 展示；没有时使用默认描述。
- 同名冲突：工作区级覆盖用户级；同一级内出现重复 skill 名称时启动失败，避免隐式覆盖。
- 加载后系统提示只展示描述列表；`load_skill` 按名称读取正文。

### 4.3 工作区 MCP 配置

路径：`<cwd>/.link/.mcp.json`

兼容两种常见格式：

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

以及：

```json
{
  "mcpServers": {
    "filesystem": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "."],
      "env": {}
    }
  }
}
```

归一化规则：

- `transport` 为空且存在 `command` 时默认为 `stdio`。
- `transport` 为 `sse` 或 `streamable` 时必须存在 `url`。
- `enabled` 为空时默认为 `true`。
- `stdio` 支持 `command`、`args`、`env`；实现时扩展 `buildTransport` 支持 stdio，并复用 stdio command 构造 helper。
- server ID 使用 `workspace:<sanitized-name>`，server name 使用 `.mcp.json` 中的原始名称，工具名前缀继续遵循现有 `mcp__<server>__<tool>` 规则。
- 单个 MCP 连接失败不阻断 TUI 启动，但 `/mcp` 应展示失败状态和错误摘要。

## 5. 模块设计

### 5.1 配置读取

新增或扩展 `internal/config`：

- 新增 `LinkSettingsPath()`：返回 `~/.link/settings.json`。
- 新增 `loadLinkLLMConfig()`：解析 `settings.json` 并映射到 `Config`。
- 修改 `LoadLocalConfig(cwd string)`：LLM 配置从 `loadLocalLLMConfig(fileCfg)` 切换为 `loadLinkLLMConfig()`；保留 `config.json` 中非敏感运行时字段。
- 单测覆盖：文件缺失、JSON 非法、字段缺失、字段 trim、token 不出现在错误中。

涉及代码：`internal/config/local_config.go:13`、`internal/config/local_config.go:53`。

### 5.2 Skill Loader

扩展 `internal/sessions`：

- 将 `skillEntryFileName` 从单一常量改为大小写兼容判断函数，例如 `isSkillEntryFile(name string)`。
- 新增 `LoadSkillsFromDirs([]SkillDir)`，其中 `SkillDir` 包含 `Path`、`Source`、`Priority`。
- `SkillEntry` 增加 `Source string` 字段，便于 `/skills` 展示来源。
- 合并顺序：先用户级，再工作区级，使工作区级覆盖用户级。

启动时在 `local.Bootstrap` 中替换现有单目录加载：

1. 解析 `userSkillsDir = ~/.link/skills`。
2. 解析 `workspaceSkillsDir = cfg.WorkspaceRoot/.link/skills`。
3. 加载两个目录并合并为 runtime 的本地 skill loader。
4. 保存 skill 摘要到 `Bundle`，传入 TUI。

涉及代码：`internal/sessions/skill.go:24`、`internal/sessions/skill.go:45`、`internal/local/bootstrap.go:47`、`internal/tools/load_skill.go:27`。

### 5.3 `load_skill` 路径与优先级

现有 `SkillSnapshot` 字段命名为 `UserSkills` 与 `LocalSkills`，其中 `UserSkills` 实际对应 DB skills，`LocalSkills` 对应本地内置 skills。TUI 新语义下调整为：

- `UserSkills`：用户级 `~/.link/skills`。
- `WorkspaceSkills`：当前工作区 `.link/skills`。
- `Merged`：按“用户级 -> 工作区级”合并后的 loader。

`LoadSkill(name)` 推荐直接从 `Merged` 读取，再根据 entry 的 `Source` 展示来源，避免查找顺序与合并顺序不一致。当前代码在 `internal/tools/load_skill.go:52` 先查 `UserSkills` 再查 `LocalSkills`，这会导致用户级覆盖工作区级，不符合本次设计，需要调整。

### 5.4 MCP 加载与连接

扩展 `internal/agent/mcp`：

- 新增 `LoadWorkspaceConfig(path string) ([]storage.MCPServer, []ServerLoadStatus, error)`，负责解析 `<cwd>/.link/.mcp.json`。
- 支持 `mcp_servers` 数组与 `mcpServers` map 两种格式。
- 扩展 transport 构造：普通用户/工作区 MCP 也支持 `stdio`，不再只有内置 MCP 支持 stdio。现有 `buildTransport` 只支持 `sse`、`streamable`，位置在 `internal/agent/mcp/transport.go:40`。
- `Manager` 增加可查询状态的方法，例如 `Snapshot(userID string) MCPSnapshot`，返回 server 名称、transport、enabled、connected、tool count、last error。

启动流程：

1. `local.Bootstrap` 读取 `<cwd>/.link/.mcp.json`。
2. 将 workspace MCP servers 设置到 manager。
3. 自动连接并发现工具。
4. 把 MCP 摘要与工具数量写入 `Bundle`，传给 TUI。

现有 `local.Store.ListEnabledMCPServersByUser` 返回空，位置在 `internal/local/store.go:40`。本次不把 workspace MCP 注入 store，而是在 `Manager` 中增加 `SetWorkspaceServers`，避免把本地文件配置伪装成数据库用户配置。

### 5.5 TUI 展示

扩展 `tui.SessionInfo`：

```go
type SessionInfo struct {
    User storage.User
    Conversation storage.Conversation
    CWD string
    Skills []SkillInfo
    MCPServers []MCPServerInfo
    SkillCount int
    MCPToolCount int
}
```

`/skills` 输出格式：

```text
已加载 Skills：3 个
- code-review [workspace] 代码审查助手
  path: /path/to/project/.link/skills/code-review/skill.md
- writer [user] 写作助手
  path: /Users/xxx/.link/skills/writer/skill.md
```

`/mcp` 输出格式：

```text
MCP Servers：2 个，工具：5 个
- filesystem [stdio] connected, tools: 3
  command: npx -y @modelcontextprotocol/server-filesystem .
- docs [sse] failed
  url: https://example.com/sse
  error: connect timeout
```

涉及代码：`internal/tui/app.go:18`、`internal/tui/app.go:133`、`internal/tui/app.go:170`。

## 6. 启动数据流

```text
main.go
  -> cli.Main / cli.Run
    -> local.Bootstrap(cwd)
      -> config.LoadLocalConfig(cwd)
        -> ~/.link/settings.json 读取 LLM env
        -> config.json 读取运行时非敏感字段
      -> 加载 ~/.link/skills
      -> 加载 <cwd>/.link/skills
      -> runtime.SetBuiltinSkills(mergedLinkSkills)
      -> mcp.LoadWorkspaceConfig(<cwd>/.link/.mcp.json)
      -> mcpManager.SetWorkspaceServers(...)
      -> mcpManager.EnsureWorkspaceSessions(...)
      -> 返回 Bundle{Runtime, SkillInfos, MCPInfos}
    -> tui.Run(...SessionInfo...)
```

对话时：

```text
用户输入普通消息
  -> runtime.RespondToConversation
    -> buildSkillSnapshot
    -> buildSystemPrompt 注入 skills 摘要
    -> toolDefinitionsForUser 合并内置工具与 MCP 工具
    -> 模型按需调用 load_skill(name)
      -> SkillSnapshot.LoadSkill(name)
      -> 从合并后的 .link skill loader 返回正文
```

## 7. 错误处理

- `~/.link/settings.json` 不存在：启动失败，提示创建文件与必要字段。
- `~/.link/settings.json` JSON 非法：启动失败，提示 parse 错误与文件路径。
- LLM 三字段缺失：启动失败，逐项列出缺失字段，不输出 token 内容。
- skill 目录不存在：忽略，数量为 0。
- skill 文件读取失败：记录 warning，跳过单个文件；若 frontmatter 明显非法，跳过并在 `/skills` 中展示 warning 计数。
- 同级 skill 名称重复：启动失败，避免模型看到不确定能力。
- `.link/.mcp.json` 不存在：忽略，`/mcp` 展示“未配置工作区 MCP”。
- `.link/.mcp.json` JSON 非法或字段非法：启动失败，提示具体 server 与字段。
- 单个 MCP 连接失败：不阻断启动，记录状态，`/mcp` 展示失败原因；可用 MCP 继续提供工具。

## 8. 测试计划

### 8.1 单元测试

- `internal/config/local_config_test.go`
  - 从临时 HOME 的 `.link/settings.json` 读取 LLM 配置。
  - 缺失 `env.open_auth_token/open_model/open_base_url` 时返回明确错误。
  - 原 `OPENAI_API_KEY` 不再是 TUI LLM 的必要条件。

- `internal/sessions/skill_test.go`
  - 同时识别 `skill.md` 与 `SKILL.md`。
  - 用户级与工作区级合并，工作区级同名覆盖用户级。
  - `load_skill` 返回被覆盖后的工作区版本。

- `internal/agent/mcp/config_test.go`
  - 解析 `mcp_servers` 数组格式。
  - 解析 `mcpServers` map 格式。
  - `stdio` 默认 transport 与 `enabled` 默认 true。
  - 非法 transport、缺失 command/url、重复名称报错。

- `internal/agent/mcp/transport_test.go`
  - 普通 workspace stdio server 可构造 `CommandTransport`。
  - sse/streamable headers 保持原行为。

- `internal/tui/events_test.go` 或新增测试
  - `/skills` 展示名称、来源、描述、路径。
  - `/mcp` 展示 server 状态、工具数与错误摘要。

### 8.2 集成验证

在临时目录构造：

```text
$HOME/.link/settings.json
$HOME/.link/skills/user-skill/skill.md
<cwd>/.link/skills/project-skill/skill.md
<cwd>/.link/.mcp.json
```

执行：

```bash
go test ./...
go run . --cwd <temp-project>
```

手动验证：

- 状态栏 skills/mcp 数量正确。
- `/skills` 列出用户级与工作区级 skill。
- `/mcp` 列出配置的 server、连接状态与工具数量。
- 普通对话 system prompt 只包含 skill 摘要，不包含 skill 正文。
- 模型调用 `load_skill` 时能读取 `.link/skills` 下对应正文。

## 9. 实施步骤

1. 配置层：实现 `~/.link/settings.json` 解析并替换 TUI LLM 来源。
2. Skill 层：支持 `skill.md`，实现用户级 + 工作区级扫描、合并、摘要输出，并修正 `load_skill` 优先级。
3. MCP 层：实现 `<cwd>/.link/.mcp.json` 解析，支持 workspace stdio/sse/streamable，记录连接状态。
4. TUI 层：扩展 `SessionInfo`，完善 `/skills` 与 `/mcp` 展示。
5. 测试与回归：补齐单元测试，运行 `go test ./...`。

## 10. 待用户确认点

- 本设计将需求中的 `~./.link/settings.json` 解释为 `~/.link/settings.json`。
- `.mcp.json` 将同时兼容 `mcp_servers` 与 `mcpServers` 两种格式。
- 同名 skill 的覆盖规则为“工作区级覆盖用户级；同一级重复名称启动失败”。
- MCP 单个 server 连接失败不阻断 TUI 启动，但在 `/mcp` 中展示失败原因。
