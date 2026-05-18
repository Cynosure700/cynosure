## Why

当前 `go-agent` 已经演进出较完整的 Web Agent 服务能力，但部分核心文件持续膨胀，出现单文件职责过多、跨层逻辑混写、包边界不够清晰的问题。这会提高后续维护和测试成本，也让 README 中的项目结构说明逐渐难以准确反映真实实现，因此需要在不改变现有功能与对外行为的前提下完成一次保守的模块化整理。

## What Changes

- 对 `go-agent` 的核心热点包进行保守重构，优先拆分超长文件，按职责收敛为更清晰的模块边界。
- 保持现有 Web API、聊天运行时、工具白名单、workspace 约束、部署产物与配置解析行为不变，重构仅以结构优化和可维护性提升为目标。
- 规范 `internal/web/app`、`internal/web/runtime`、`internal/tools`、`internal/config`、`internal/web/storage`、`internal/sessions` 等目录下的文件职责，减少单文件混合装配、业务编排、辅助函数和基础设施细节的情况。
- 同步更新 `go-agent/README.md`，让目录结构、核心模块职责、开发验证方式与重构后的实现保持一致。

## Capabilities

### New Capabilities
- `go-agent-module-layout`: 约束 go-agent 的核心模块应按职责清晰拆分，并要求重构后继续保持现有服务行为、运行时边界与部署方式不变。

### Modified Capabilities

## Impact

- 影响 `go-agent/internal/web/app/server.go` 及相关 HTTP 路由 / handler 组织方式。
- 影响 `go-agent/internal/web/runtime/runtime.go`、`go-agent/internal/web/runtime/tools.go` 与 `go-agent/internal/tools/registry.go` 的职责划分。
- 影响 `go-agent/internal/config/config.go`、`go-agent/internal/web/storage/store.go`、`go-agent/internal/sessions/skill.go` 等大文件的内部结构拆分。
- 影响 `go-agent/README.md` 中的目录结构、模块说明与开发验证章节。
- 需要通过现有 Go 测试与构建流程验证功能等价，确保本次改动不引入行为回归。
