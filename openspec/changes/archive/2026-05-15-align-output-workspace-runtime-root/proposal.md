## Why

当前 go-agent 虽然已经引入 `workspace` 与 `output/workspace` 的部署产物布局，但实际 agent 运行时默认终端 cwd 仍可能落在源码侧的 `workspace/`，而不是部署包侧的 `output/workspace/`。这会导致运行时无法稳定识别随部署发布的 `skills` 和 `cmd` 命令产物，在云端部署时尤其容易失效，同时也缺少一套对本地调试友好的目录解析策略。

## What Changes

- 调整运行时根目录解析策略，优先把部署包中的 `output/workspace/` 识别为 agent 的默认工作目录，而不是源码根下的 `workspace/`
- 为内置 skills、命令二进制和脚本命令增加统一的“部署优先、本地回退”路径解析规则，避免依赖进程启动 cwd
- 明确区分服务包根目录、运行时 workspace 根目录和源码工作区目录，支持云端部署与本地调试两种模式使用同一套配置模型
- 调整 Web runtime 与工具运行环境注入逻辑，确保终端、文件工具、skill 加载与命令执行共享同一套解析后的 runtime root
- 增加测试覆盖，验证部署模式优先命中 `output/workspace`，本地开发模式在未打包时仍可回退到源码 `workspace`

## Capabilities

### New Capabilities
- `deployment-runtime-layout`: 定义 go-agent 在部署包与本地源码环境下如何解析运行时根目录、workspace 目录以及随包发布的 skills/cmd 产物目录

### Modified Capabilities
- `agent-runtime`: 调整会话运行时的 workspace 解析与注入规则，要求聊天回合、skill 装载与工具调用使用统一的 runtime root
- `tool-execution-control`: 调整终端与文件工具的默认 cwd 和路径边界规则，要求优先绑定部署 workspace，并在本地调试时安全回退
- `skill-management`: 调整平台内置 skill 的加载来源，要求优先从部署 runtime workspace 的 skill catalog 读取，而不是隐式依赖源码相对路径

## Impact

- 影响 `go-agent/internal/config/config.go` 中应用根目录、workspace 根目录和运行时资产目录的解析规则
- 影响 `go-agent/internal/web/runtime/runtime.go` 与 `go-agent/internal/web/runtime/tools.go` 中 runtime 环境注入、cwd 选择与命令路径计算
- 影响 `go-agent/internal/web/app/server.go` 与 `go-agent/internal/sessions/skill.go` 中 builtin skills 的启动加载路径
- 影响 `go-agent/config.json`、`go-agent/output/config.json` 及相关测试，需确保部署模式和本地模式都能通过同一套约定启动
