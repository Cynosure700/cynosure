## Why

当前 go-agent 的目录相关变量在配置加载、runtime 环境注入和 README 说明中分散表达，同一组路径关系被多处重复推导，增加了理解成本，也让后续维护更容易出现命名和解释不一致的问题。现在需要在不改变现有运行行为的前提下，把目录变量收敛为更清晰的主从关系，并同步更新文档。

## What Changes

- 收敛 go-agent 中目录变量的解析入口，明确哪些是根目录变量，哪些是从属派生目录。
- 精简 runtime 与工具层对目录变量的重复拼装逻辑，保留现有默认值、回退顺序和边界校验。
- 更新 `go-agent/README.md`，用更简洁的方式说明目录变量、默认解析规则和推荐配置方式。

## Capabilities

### New Capabilities

### Modified Capabilities
- `deployment-runtime-layout`: 明确目录变量的主从关系，要求 runtime 资产目录继续从单一的 workspace 根目录一致派生。
- `agent-runtime`: 明确 runtime 暴露给技能和工具的目录变量必须与已解析的 workspace 根目录保持一致，避免多处重复推导。

## Impact

- 影响 `go-agent/internal/config/config.go` 中的目录配置解析与校验代码。
- 影响 `go-agent/internal/web/runtime/tools.go` 以及可能相关的 runtime 环境组装代码，减少重复目录推导。
- 影响 `go-agent/README.md` 中对 `APP_HOME`、`WORKSPACE_ROOT`、`BUILTIN_SKILLS_DIR`、`COMMAND_BIN_DIR`、`COMMAND_SCRIPT_DIR` 的说明。
- 需要补充或更新配置与 runtime 相关测试，确保行为保持兼容。
