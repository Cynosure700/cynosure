## Why

当前 `go-agent` 仓库仍混杂了 Web 服务、CLI/REPL、技能运行时、部署产物构建器与历史 workspace 结构，导致代码入口分散、目录职责不清，README 中本地与云端部署步骤也偏重兼容路径而不是最小可用路径。其中 `workspace/bin` 与 `workspace/cmd` 实际上仍是 Web agent 执行二进制命令和脚本所依赖的运行时资产，不应被当作 CLI 遗留物误删；现在需要把项目收敛为“只负责启动和提供 Web 服务”的后端，同时保留并理顺这些 runtime 资产的职责，以减少维护面、降低部署复杂度，并让开发者能够更直接地完成本地和云服务器部署。

## What Changes

- **BREAKING** 将 `go-agent` 收敛为仅提供 Web 服务的后端程序，移除 CLI REPL 与不再服务于 Web 启动链路的历史代码结构，但保留 Web agent 所需的 runtime workspace 资产目录与命令执行能力
- **BREAKING** 重组项目目录，只保留 Web 服务启动、HTTP 应用、认证、存储、基础配置、日志以及 Web agent 运行所需的 runtime 资产装配链路，并清理与 Web 服务无关的历史兼容结构
- 精简启动入口和运行时初始化路径，统一默认启动方式，避免 `main.go`、重复 `cmd/` 入口以及混杂的历史兼容层继续扩散
- 简化部署模型，明确“本地开发启动”和“云服务器部署运行”两条最短路径，并把 `output/workspace` / `workspace` 统一解释为 Web runtime 资产布局，而不是 CLI 的双重形态
- 更新 `README.md`，将文档改为围绕 Web 服务的项目结构说明、环境配置、快速启动、本地部署与云端部署指引

## Capabilities

### New Capabilities

### Modified Capabilities
- `agent-runtime`: 运行时职责收敛为服务 Web 聊天请求所需的最小后端能力，移除 CLI agent 与 REPL 结构，但保留 `workspace/bin`、`workspace/cmd` 等 Web agent 所需 runtime 资产
- `deployment-runtime-layout`: 部署布局继续围绕 Web 服务进程与 runtime workspace 资产组织，但需要澄清 `output/workspace` 与源码 workspace 的职责和解析规则
- `web-chat`: 浏览器聊天主链路继续保留，但其后端依赖必须来自精简后的 Web 服务结构与统一启动方式
- `tool-execution-control`: 浏览器运行时只保留 Web 服务实际需要的受控能力，并明确命令执行通过 `workspace/bin`、`workspace/cmd` 这类 runtime 资产完成，而不是通过 CLI 入口完成
- `workspace-package-layout`: 构建与发布物布局收敛为 Web 服务部署真正需要的最小集合，同时继续正确产出 Web agent 所需的 workspace runtime 资产

## Impact

- 影响应用入口与项目结构，包括 `go-agent/main.go`、`go-agent/cmd/`、`go-agent/internal/` 下的模块划分与保留范围
- 影响运行时与工具边界，包括 `go-agent/internal/web/runtime/`、`go-agent/internal/tools/`、`go-agent/internal/agent/`、`go-agent/internal/sessions/`
- 影响部署与打包路径，包括 `go-agent/build.sh`、`go-agent/output/`、`go-agent/workspace/` 与相关目录约定
- 影响开发与运维文档，包括 `go-agent/README.md` 中的结构说明、本地启动步骤、云服务器部署步骤与运维建议
