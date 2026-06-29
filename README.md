# Cynosure

Cynosure 是一个面向本地代码工作的终端 TUI Agent。它以当前目录作为工作区，支持读取文件、搜索代码、执行命令、修改代码、调用 Skills、接入 MCP 工具，并把会话、记忆和大工具结果保存在本机。

详细设计和完整能力说明见 [cynosure/README.md](./cynosure/README.md)。

## 功能特性

- **终端 TUI 交互**：在终端中进行多轮对话，支持流式输出、Markdown 渲染、状态栏、快捷命令和历史会话恢复。
- **本地代码工具**：内置文件读取、写入、编辑、搜索、目录浏览、命令执行、网页抓取、任务清单、子 Agent 等工具。
- **权限审批**：变更类操作在执行前触发交互式审批，可选择单次允许、记住规则或拒绝。
- **Skill 系统**：加载内置、用户级和工作区级 Skills，让 Agent 按需获得任务专用能力。
- **MCP 集成**：读取工作区 `.cynosure/.mcp.json`，自动连接并暴露 MCP 工具。
- **项目级记忆**：在 `~/.cynosure/memory/` 下按工作区维护长期记忆和会话记忆。
- **上下文压缩**：自动压缩长会话，并把大工具结果落盘后通过 `read_persisted_output` 分段读取。
- **纯本地运行**：TUI 模式不依赖 Web 服务、数据库、Redis 或 Elasticsearch。

## 系统要求

Cynosure 是一个运行在终端中的 TUI 程序，并依赖系统的 `bash` 来执行命令工具，因此主要适用于类 Unix 系统：

- **macOS**：原生支持。
- **Linux**：原生支持。
- **Windows**：不直接支持。命令执行工具依赖 `bash`，需在 WSL 或 Git Bash 等提供 `bash` 的环境中运行。

此外还需要：

- Go 1.26.1 或更高版本（用于构建和安装）。
- 系统已安装 `bash`（`bash` 工具直接调用系统 `bash`）。
- 一个 OpenAI 兼容的 LLM 服务（通过 `~/.cynosure/settings.json` 配置）。

## 安装

### 安装 Go

Cynosure 使用 Go 构建和安装。请先安装 Go 1.26.1 或更高版本：

- macOS 可使用 Homebrew：`brew install go`
- 其他系统可从 [go.dev/dl](https://go.dev/dl/) 下载安装包

安装完成后确认 Go 可用：

```bash
go version
```

### 安装 Cynosure

执行以下命令安装：

```bash
go install github.com/example/cynosure/cynosure@latest
```

安装后，`cynosure` 会进入 `GOBIN`，默认通常是 `~/go/bin`。请确保该目录已加入 `PATH`。

## 配置

Cynosure 启动时从 `~/.cynosure/settings.json` 读取模型配置。最小配置示例：

```json
{
  "env": {
    "open_auth_token": "your-api-key",
    "open_model": "your-model",
    "open_base_url": "https://api.example.com"
  }
}
```

本地运行会自动创建 `~/.cynosure/` 下的会话、日志、记忆和工具输出目录。工作区级配置放在当前项目的 `.cynosure/` 目录下，例如 `.cynosure/.mcp.json` 和 `.cynosure/settings.json`。

## 使用

在任意项目目录中启动：

```bash
cd /path/to/your/project
cynosure
```

显式指定工作区：

```bash
cynosure --cwd /path/to/your/project
```

查看命令帮助：

```bash
cynosure help
```

常用 TUI 命令：

- `/help`：显示可用斜杠命令。
- `/clear`：开启全新对话并清空当前上下文。
- `/cwd`：显示当前工作区。
- `/skills`：显示已加载的 Skill。
- `/mcp`：显示 MCP server 状态。
- `/resume`：恢复当前工作区的历史会话。

## 项目导览

```text
.
├── cynosure/              # Go module，TUI Agent 主体
│   ├── main.go            # 程序入口
│   ├── assets/            # go:embed 嵌入资源、系统提示词和内置 Skills
│   └── internal/          # Agent runtime、工具、TUI、配置、记忆和会话实现
├── docs/                  # 设计文档
├── cn/                    # 参考项目与工具定义
└── README.md              # 当前项目首页
```

核心代码位于 `cynosure/internal/`：

- `agent/`：Agent runtime、MCP 与存储模型。
- `assistant/`：system prompt 拼装。
- `cli/`：命令行参数解析。
- `config/`：配置与本地路径。
- `local/`：启动装配与本地 Store。
- `tools/`：内置工具定义、参数校验和执行。
- `tui/`：Bubble Tea TUI 界面与事件处理。

更多运行机制、工具清单、权限模型、记忆系统和上下文压缩细节见 [cynosure/README.md](./cynosure/README.md)。
