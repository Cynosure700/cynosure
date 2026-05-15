## Why

当前 go-agent 已经支持服务端 workspace 与部署产物目录，但目录布局仍偏向“应用根目录 + 用户工作区”分离模型，和目标部署方式不一致。现在需要把运行时约定收敛为单一 `workspace/` 根目录：服务启动后从 `workspace/skills` 加载内置 skills，从 `workspace/bin` 调用由 `cmd` 编译出的二进制，并确保 agent 的终端与文件工具只能在该 workspace 范围内运行。

## What Changes

- 将部署产物目录布局调整为统一的 `workspace/` 根目录，而不是分别依赖应用根目录与独立 workspace 根目录
- 约定服务启动时从 `workspace/skills` 加载平台内置 skills，并将其作为运行时默认 skill catalog
- 约定 `cmd` 下的 Go 入口在构建阶段编译到 `workspace/bin`，供服务端 runtime 与 skills 稳定调用
- 明确 agent 的终端命令与文件类工具只能作用于受控的 `workspace/` 范围，不得直接以宿主机用户目录作为默认操作范围
- 增加标准化 `build.sh` 打包脚本，统一构建二进制、复制 skills、配置与启动脚本到部署输出目录

## Capabilities

### New Capabilities
- `workspace-package-layout`: 定义部署输出中的 `workspace/` 根目录结构，以及 `workspace/skills`、`workspace/bin` 等运行时资产目录的组织方式

### Modified Capabilities
- `agent-runtime`: 调整运行时资产定位规则，使服务初始化从统一 workspace 根目录加载 skills 与命令产物
- `tool-execution-control`: 调整工具执行边界，使终端与文件工具默认在服务端 workspace 内运行，而不是访问宿主机任意目录
- `skill-management`: 调整内置 skill 的来源与加载位置，统一从部署后的 `workspace/skills` 目录读取

## Impact

- 影响 `go-agent/internal/config/` 中运行时路径配置和绝对路径解析方式
- 影响 `go-agent/internal/deploy/` 与构建脚本，需要生成符合新目录约定的部署产物
- 影响 `go-agent/internal/web/runtime/`、`go-agent/internal/tools/` 与 `go-agent/internal/sessions/skill.go` 中的 skill 加载、命令路径注入与工具 cwd 限制
- 影响部署流程，需要新增或更新标准 `build.sh` 以构建 `workspace/bin` 并同步 `workspace/skills`
