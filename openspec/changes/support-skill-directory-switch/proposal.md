## Why

当前 go-agent 在执行 `Skill` 后只把技能正文注入模型，但没有把该技能所在目录暴露为后续工具调用的默认工作目录。结果是技能里引用相对脚本、资源文件或同目录辅助产物时，仍会落回共享 workspace 根目录执行，导致技能目录内的实现无法按预期工作。

## What Changes

- 为技能加载结果补充可消费的技能目录上下文，让 runtime 能识别当前回合激活技能的源目录。
- 调整工具执行上下文解析：当调用由已加载技能驱动且技能目录可用时，默认切换到该技能目录执行相对路径操作。
- 保留现有 workspace 边界校验与审计记录，并补充技能目录切换相关的拒绝与成功场景覆盖。

## Capabilities

### New Capabilities

### Modified Capabilities
- `agent-runtime`: 变更技能加载后的运行时上下文，使当前激活技能可以为后续工具调用提供默认执行目录。
- `tool-execution-control`: 变更工具默认 cwd 解析规则，在技能驱动的工具调用中允许优先使用技能目录，同时继续保证边界校验与审计可见性。

## Impact

- 影响 `go-agent/internal/web/runtime/runtime.go`、`go-agent/internal/web/runtime/tools.go` 中的技能上下文装配与工具执行流程。
- 影响 `go-agent/internal/tools/registry.go`、`go-agent/internal/tools/bash.go` 相关的运行时环境注入与默认目录选择。
- 需要补充 `go-agent/internal/web/runtime/runtime_test.go` 及相关工具测试，验证技能目录切换、相对路径解析和审计行为。
