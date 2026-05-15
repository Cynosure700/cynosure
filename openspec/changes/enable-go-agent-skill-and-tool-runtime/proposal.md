## Why

当前 `go-agent` 默认以 Web 服务方式运行，但真正在线的 runtime 只支持 `load_skill`，既不会把旧版 go-agent 的 tool 调用能力接入 Web，也不会把 `go-agent/skills` 目录下的默认 skill 暴露给所有用户。现在需要把它收敛成一个可部署的服务端 agent：统一加载平台内置 skill 与用户自定义 skill，允许 agent 在受控工作目录中调用已注册 tool，并保持多用户隔离。同时，部署后的运行根目录必须固定下来，因为 skill 可能需要调用随服务一起发布的二进制命令或 `.py` 脚本，这些产物需要从 `cmd` 目录构建并随部署一并交付。

## What Changes

- 将当前 Web runtime 从“仅能懒加载数据库 skill 的聊天机器人”扩展为“可调用 skill 和 tool 的 go-agent runtime”
- 新增共享内置 skill 目录机制，服务启动时加载 `go-agent/skills` 下的默认 skill，并让所有用户在对话时都可调用
- 保留并强化数据库中的用户自定义 skill 能力，仅为所属用户加载已启用的自定义 skill，并与内置 skill 合并注入 runtime
- 将现有 go-agent 工具注册表接入 Web runtime，允许 agent 在服务端调用受控的 tool 能力，而不是只暴露 `load_skill`
- 新增服务端工作目录策略：为部署实例定义固定的应用根目录，并为每个用户分配独立 workspace 根目录，在同一部署模型下区分“应用资源目录”和“用户工作目录”
- 新增部署命令产物目录机制：约定 `cmd` 目录中的 Go 入口在部署时编译成二进制产物，`.py` 等脚本文件作为只读辅助命令随服务发布，供 skill 通过固定路径调用
- 明确 go-agent 的部署模型为长期运行的 Web 服务，CLI REPL 不再作为主运行时能力来源，服务启动时需要校验应用根目录、skills 目录、cmd 产物目录和 workspace 根目录

## Capabilities

### New Capabilities
- `builtin-skill-catalog`: 定义 `go-agent/skills` 目录中的平台内置 skill 如何被加载、命名、暴露给所有用户，以及与用户自定义 skill 的合并规则
- `user-workspace`: 定义服务端用户工作目录的创建、默认路径选择、会话工作目录绑定和跨用户隔离规则
- `deployment-command-artifacts`: 定义部署根目录、`cmd` 目录命令产物的构建与发布方式，以及 skill 如何稳定引用这些只读命令资源

### Modified Capabilities
- `agent-runtime`: 调整运行时行为，使其在每轮对话中合并内置 skill 与用户自定义 skill，并能够调用平台允许的已注册 tool
- `skill-management`: 调整 skill 体系边界，明确“平台内置共享 skill”与“数据库中的用户自定义 skill”两类来源及其运行时可见性
- `tool-execution-control`: 调整工具执行规则，明确 Web runtime 下的已注册工具白名单、工作目录约束、权限校验与审计要求

## Impact

- 影响 `go-agent/internal/web/runtime/` 的主循环、skill 注入、tool 注册与执行流程
- 影响 `go-agent/internal/sessions/skill.go` 与 `go-agent/skills/` 的加载方式，需要把本地默认 skill 接入 Web runtime
- 影响 `go-agent/internal/tools/` 与 `go-agent/internal/safety/`，需要将既有工具能力纳入 Web 场景并增加工作目录隔离
- 影响配置与部署约定，需要定义应用根目录、cmd 产物目录、workspace 根目录、默认工作目录和服务启动方式
- 影响数据库 skill 的运行时装配逻辑，但不改变用户自定义 skill 的所有权隔离原则
- 影响部署流水线，需要增加 `cmd` 目录下 Go 命令的构建步骤，并将脚本类辅助命令一并打包到部署产物中
