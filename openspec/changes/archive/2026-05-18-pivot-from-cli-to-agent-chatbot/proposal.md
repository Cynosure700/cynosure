## Why

当前项目虽然已经具备 Web 对话链路，但默认入口与默认 agent 定位仍然明显偏向 CLI 编码助手：根入口仍直接启动 REPL，系统提示词也将产品定义为 coding agent。现在需要把产品形态彻底收敛为类似 ChatGPT 的通用 Agent 聊天机器人，让用户通过浏览器对话即可完成日常问答、规划、分析与任务协助，而不是继续围绕命令行编码场景构建。

## What Changes

- **BREAKING** 移除 CLI 作为默认和主要使用方式，将 Web 服务 + 浏览器聊天界面调整为唯一面向用户的产品入口
- 将默认 assistant 定位从“coding agent”改为“通用 agent 聊天机器人”，优先直接对话，只有在确有帮助时才调用技能或工具
- 将 Web 端体验调整为 conversation-first 的聊天产品，弱化开发者导向的信息密度与操作心智，使其更接近 ChatGPT 的会话体验
- 调整 runtime 的系统提示词、能力暴露方式与工具使用策略，使回答不再预设编码任务、工作区操作或文件编辑需求
- 保留技能与受控工具能力，但改为可选增强项，并对高权限工具采用更严格的安全边界与默认关闭策略
- 为现有 CLI 相关入口、启动方式与文案提供迁移路径，避免产品和实现继续围绕终端使用模型演化

## Capabilities

### New Capabilities
- `application-entrypoint`: 定义服务化、浏览器优先的产品入口，并退役 CLI REPL 作为正式用户交互方式

### Modified Capabilities
- `web-chat`: 将网页聊天体验调整为通用对话优先、低心智负担的聊天机器人界面
- `agent-runtime`: 将运行时默认角色从编码助手改为通用 agent，并调整技能/工具调用策略
- `tool-execution-control`: 将网页聊天模式下的工具暴露改为安全默认、按需升级，而不是默认围绕开发操作设计

## Impact

- 影响应用入口与启动方式，包括 `go-agent/main.go`、`go-agent/cmd/web/main.go` 及相关运行脚本
- 影响默认系统提示词与 agent 运行时行为，包括 `go-agent/internal/agent/` 与 `go-agent/internal/web/runtime/`
- 影响 Web 聊天产品交互与页面组织，包括 `web/src/App.tsx`、前端 API 调用层与样式表现
- 影响工具注册、权限控制与默认暴露范围，包括 `go-agent/internal/tools/` 与 `go-agent/internal/web/runtime/tools.go`
- 需要梳理 CLI 兼容、启动迁移和已有开发者使用方式的回退路径
