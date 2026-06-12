# TUI 化改造设计文档

## 1. 目标

将 `go-agent` 改造为本地 TUI 代码助手，默认进入终端交互界面，面向用户当前工作区提供代码阅读、代码修改、命令执行、Skill 加载与内置 MCP 工具接入能力。

核心目标：

1. `go run .` 默认启动 TUI。
2. 工作区默认为启动时的当前目录，也可通过 `--cwd` 指定。
3. 支持 `read_file`、`write_file`、`edit_file`、`bash`、`todo_write`、`load_skill`、`spawn_subagent`、`read_persisted_output` 等本地工具。
4. 启动时加载本地 Skills，并在每轮对话前生成一致的 Skill snapshot。
5. 启动时读取 `<app_home>/mcp_config.json`，自动连接内置 stdio MCP。
6. 使用 Bubble Tea、Lipgloss 与 Glamour 构建可交互、可读性好的终端界面。
7. 移除旧服务入口、鉴权、路由和外部持久化依赖，TUI 初版使用本地内存 Store。

## 2. 最终架构

```text
main.go
  └── internal/cli.Run
        ├── 解析 --cwd
        ├── internal/local.Bootstrap
        │     ├── config.LoadLocalConfig
        │     ├── sessions.LoadBuiltinSkillsFromDir
        │     ├── assistant.LoadBaseSystemPrompt
        │     ├── mcp.LoadBuiltinConfig + mcp.Manager
        │     ├── local.Store
        │     └── runtime.Service
        └── internal/tui.Run
              ├── Bubble Tea Model / Update / View
              ├── textarea 输入区
              ├── viewport 输出区
              ├── markdown 渲染
              └── EventWriter 桥接 runtime 事件
```

代码分层：

- `internal/cli`：命令入口和参数解析。
- `internal/config`：TUI 本地配置加载与路径解析。
- `internal/local`：本地启动装配和内存 Store。
- `internal/tui`：终端界面、快捷命令、事件桥接和取消控制。
- `internal/agent/runtime`：Agent loop、工具注册、上下文压缩、子 Agent、Hook。
- `internal/agent/mcp`：内置 MCP 配置、连接管理、工具发现与调用。
- `internal/agent/storage`：runtime 复用的数据模型与可选存储类型。

## 3. 启动流程

1. `main.go` 调用 `cli.Main()`。
2. `cli.Run` 解析参数：
   - 无子命令：启动 TUI。
   - `tui`：启动 TUI。
   - `help` / `--help` / `-h`：输出帮助。
   - 其他参数：返回未知命令错误。
3. `local.Bootstrap(ctx, cwd)` 完成本地依赖装配：
   - 读取 `config.json`。
   - 从 `OPENAI_API_KEY` 读取模型密钥。
   - 将 workspace 覆盖为启动 cwd 或 `--cwd`。
   - 加载 system prompt 与内置 Skills。
   - 读取内置 MCP 配置并连接 stdio MCP。
   - 创建本地内存 Store、用户与会话。
   - 构建 `runtime.Service`。
4. `tui.Run` 启动 Bubble Tea 程序并进入交互循环。

## 4. 配置设计

TUI 配置只保留本地运行所需字段：

```json
{
  "base_url": "https://api.deepseek.com",
  "model_id": "deepseek-v4-flash",
  "app_home": ".",
  "system_prompt_path": "system_prompt.md",
  "builtin_skills_dir": "skills",
  "command_bin_dir": "bin",
  "command_script_dir": "cmd",
  "allowed_tools": "load_skill,bash,read_file,write_file,edit_file,todo_write,spawn_subagent",
  "bash_allow_outside_workspace": false,
  "bash_allow_dangerous_commands": false
}
```

规则：

- `base_url`、`model_id` 必填。
- `OPENAI_API_KEY` 必须来自环境变量。
- `app_home` 解析为绝对路径。
- `system_prompt_path`、`builtin_skills_dir`、`command_bin_dir`、`command_script_dir` 以 `app_home` 为基准解析。
- `workspace_root` 在 TUI 模式下由启动 cwd 或 `--cwd` 覆盖，且目标目录必须存在。
- `allowed_tools` 为空时使用本地默认工具集。

## 5. TUI 设计

界面由三部分组成：

1. 顶部状态栏：显示应用名、当前工作区、加载的 Skill 数量、内置 MCP 工具数量和运行状态。
2. 中间输出区：展示用户消息、助手回复、工具事件、错误和系统提示；助手消息使用 Markdown 渲染。
3. 底部输入区：支持多行输入，`Enter` 发送。

快捷键与命令：

- `Enter`：发送当前输入。
- `Ctrl+C`：请求生成中时取消当前请求；空闲时退出。
- `/help`：显示帮助。
- `/clear`：清空当前 TUI 显示上下文。
- `/cwd`：显示当前工作区。
- `/skills`：显示启动时加载的 Skill 数量。
- `/mcp`：显示内置 MCP 工具数量。

## 6. Runtime 事件桥接

`runtime.EventWriter` 的接口保持为：

```go
type EventWriter interface {
    Event(name string, data any)
}
```

TUI 通过 `internal/tui.EventWriter` 将 runtime 事件转换为 Bubble Tea message：

- `assistant_delta`：追加助手流式文本。
- `reasoning_delta`：可选展示推理片段。
- `tool`：展示工具调用摘要。
- `meta`：更新 token、耗时、工具调用次数等状态。
- `assistant`：落最终助手消息。
- `error`：展示错误。
- `done`：恢复空闲状态。

为避免取消后的旧事件污染新请求，事件中携带 generation ID；`Ctrl+C` 取消时递增 generation，`Update` 只处理当前 generation 的事件。

## 7. 本地 Store

TUI 初版使用 `internal/local.Store` 实现 runtime 所需的 `conversationStore` 接口：

- 用户、会话、消息、工具调用、子 Agent 消息保存在内存中。
- 上下文压缩需要的 persisted output 和 summary 保存在内存中。
- 记忆相关接口返回空结果，`EnableMemory=false`。
- 进程退出后不保留历史。

这样 TUI 可以在没有外部服务的情况下独立运行。

## 8. Skill 加载

启动时从 `builtin_skills_dir` 读取本地 Skills。每轮对话前生成 Skill snapshot，并同时用于：

- system prompt 中的 Skill 摘要注入。
- `load_skill` 工具读取完整说明。

这保证同一轮对话内模型看到的 Skill 摘要与工具读取到的正文一致。

## 9. MCP 加载

启动时读取 `<app_home>/mcp_config.json`：

```json
{
  "mcp_servers": [
    {
      "name": "filesystem",
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/path/to/project"],
      "env": {},
      "enabled": true
    }
  ]
}
```

约束：

- 内置 MCP 只支持 stdio。
- 文件不存在时跳过 MCP 初始化。
- 文件存在但格式错误时启动失败。
- 成功发现的 MCP 工具统一以 `mcp__{server}__{tool}` 命名加入工具列表。

## 10. 验收标准

1. `go run .` 默认进入 TUI。
2. `go run . --cwd /path/to/project` 以指定目录作为工作区。
3. TUI 中可读取、编辑工作区内文件并执行允许的 shell 命令。
4. 启动时能加载本地 Skills。
5. 启动时能读取并连接内置 stdio MCP。
6. 取消生成后旧事件不会污染下一次请求。
7. 仓库中不存在旧服务端目录，也不存在对旧服务端包路径的 import。
8. `go test ./...` 通过。
