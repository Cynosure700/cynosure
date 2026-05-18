## 1. 目录解析模型收敛

- [x] 1.1 收敛 `go-agent/internal/config/config.go` 中目录变量解析逻辑，确保 `WORKSPACE_ROOT` 成为运行时资产目录的单一规范来源。
- [x] 1.2 保留 `BUILTIN_SKILLS_DIR`、`COMMAND_BIN_DIR`、`COMMAND_SCRIPT_DIR` 的显式配置兼容性，并继续校验它们只能指向 `WORKSPACE_ROOT` 下的规范派生目录。

## 2. Runtime 目录使用简化

- [x] 2.1 清理 `go-agent/internal/web/runtime/tools.go` 及相关 runtime 代码中的重复目录默认推导，改为直接消费配置层输出的规范路径。
- [x] 2.2 确认工具与技能运行时暴露的目录变量仍与当前解析出的 workspace 根目录保持一致，不改变现有行为与边界约束。

## 3. 回归验证与文档同步

- [x] 3.1 更新或补充 `go-agent/internal/config/config_test.go`、`go-agent/internal/web/runtime/*_test.go` 中与目录变量相关的回归测试，覆盖默认解析、显式配置与兼容行为。
- [ ] 3.2 更新 `go-agent/README.md`，用更简洁的方式说明目录变量关系、默认推导顺序和推荐配置方式。
