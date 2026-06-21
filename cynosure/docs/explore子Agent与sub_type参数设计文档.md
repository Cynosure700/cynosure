# explore 子 Agent 与 sub_type 参数设计文档

## 1. 背景

当前 `spawn_subagent` 工具只有 `task` 与 `cwd` 两个参数。运行时在 `internal/agent/runtime/subagent.go` 中为子 Agent 创建 fresh message list，并通过 `buildSubagentSystemPrompt` 复用主 Agent 的系统提示词，再追加一段通用 `<subagent>` 规则。子 Agent 工具集由 `NewChildToolRegistry` 创建，当前策略是在主 Agent 可用工具基础上移除 `spawn_subagent`。

这套机制已经具备以下既有行为：

1. 子 Agent 看不到父对话历史，只接收自己的任务消息。
2. 子 Agent 不能再派生子 Agent。
3. 子 Agent 复用审批、超时、工具结果压缩、TUI 展示规约。
4. `cwd` 可传相对或绝对路径，最终解析为可用工作目录。

本需求要在不改变上述策略的前提下，为 `spawn_subagent` 增加 `sub_type` 参数，并新增专门用于搜索和代码库探索的 `explore` 子 Agent。`explore` 子 Agent 需要单独系统提示词，且应从运行时工具层面限制为只读检索/读取能力，不能只依赖提示词约束。

## 2. 目标

1. `spawn_subagent` 新增 `sub_type` 参数，用来声明子 Agent 类型。
2. 支持两个定义好的类型：
   - `general`：保留当前通用子 Agent 行为，但系统提示词必须明确标注它不是搜索专用子 Agent；搜索、文件定位、代码探索、实现梳理、证据收集等任务必须交给 `explore`。
   - `explore`：新增的代码搜索/文件探索子 Agent。
3. 对 `sub_type` 做严格参数校验，只允许创建已定义类型。
4. `explore` 使用单独系统提示词，不复用主 Agent 的完整身份/任务管理/写入能力提示词。
5. `explore` 的工具集限制为只读探索工具，不能暴露写入、编辑、记忆维护、任务管理或再次派生子 Agent 的工具。
6. 调整主 Agent 系统提示词，使主 Agent 知道必须为 `spawn_subagent` 选择明确的 `sub_type`，并且搜索相关任务必须委派给 `explore`，不得交给 `general`。
7. 删除或收敛因类型分支引入的冗余代码，保持子 Agent 入口集中、策略清晰。

## 3. 非目标

1. 不新增用户配置项。
2. 不改变父 Agent 的默认工具列表、审批策略、上下文压缩策略、记忆策略和 TUI 展示行为。
3. 不引入多层子 Agent 或子 Agent 并发调度框架。
4. 不改变 `general` 子 Agent 的任务执行能力和已有压缩、超时、审批语义。
5. 不实现复杂的权限沙箱或命令黑名单；`explore` 的只读约束通过工具集收窄和系统提示词共同完成。

## 4. 方案比较

### 方案 A：只在提示词中要求 sub_type 与只读

做法是在工具描述和主/子提示词里说明 `sub_type=explore` 只读，但运行时仍给 `explore` 暴露当前子 Agent 工具集。

优点是改动最小。缺点是不能满足“只用于搜索”的真实边界：模型仍可能调用 `write_file`、`edit_file`、`multi_edit`、`todo_write` 等工具，只是被提示词劝阻。该方案不推荐。

### 方案 B：sub_type 枚举 + 类型化 prompt + 类型化工具集

做法是把子 Agent 类型定义为运行时枚举；`spawn_subagent` schema 使用 `enum` 校验 `sub_type`；运行时根据类型选择系统提示词和工具 registry。`general` 走现有逻辑，`explore` 使用只读工具集与专属系统提示词。

优点是行为边界清晰，参数校验与执行能力一致，最符合需求。缺点是需要新增少量类型分发代码和测试。推荐采用该方案。

### 方案 C：把每类子 Agent 拆成独立工具

做法是保留 `spawn_subagent`，再新增 `spawn_explore_subagent`。优点是工具级语义直观。缺点是用户明确要求调整 `spawn_subagent` 参数新增 `sub_type`，且会扩大工具表面积，不利于后续类型扩展。该方案不推荐。

## 5. 推荐设计

采用方案 B。

### 5.1 子 Agent 类型定义

在 `internal/agent/runtime/subagent.go` 或新增的同包小文件中定义类型常量：

```go
type subagentType string

const (
	subagentTypeGeneral subagentType = "general"
	subagentTypeExplore subagentType = "explore"
)
```

同时提供集中校验函数：

```go
func parseSubagentType(value string) (subagentType, error)
```

校验规则：

1. `value` 必须 trim 后非空。
2. 只接受 `general` 与 `explore`。
3. 其他值返回清晰错误，例如 `unsupported sub_type "review"; allowed values: general, explore`。

### 5.2 `spawn_subagent` 参数调整

在 `internal/tools/definitions.go` 的 `spawnSubagentToolSpec` 中新增 `sub_type`：

```json
{
  "sub_type": {
    "type": "string",
    "enum": ["general", "explore"],
    "description": "Type of child agent to spawn. Use general for isolated implementation or analysis tasks; use explore for read-only codebase search and file exploration."
  }
}
```

`required` 调整为：

```json
["sub_type", "task"]
```

`spawnSubagentArgs` 调整为：

```go
type spawnSubagentArgs struct {
	SubType string `json:"sub_type"`
	Task    string `json:"task"`
	CWD     string `json:"cwd"`
}
```

由于当前 `executeToolCall` 已经读取工具 schema 并调用 `ValidateToolArgs`，schema enum 会先拦截非法类型；运行时仍保留 `parseSubagentType` 作为第二道校验，避免工具定义变化或测试直接调用 `runSubagent` 时绕过 schema。

### 5.3 子 Agent 配置对象

为避免 `runSubagent` 内部出现多处分支，新增内部配置结构：

```go
type subagentProfile struct {
	Type          subagentType
	SystemPrompt  string
	ToolRegistry  *ToolRegistry
	MaxRounds     int
}
```

新增构造函数：

```go
func (s *Service) buildSubagentProfile(kind subagentType, user storage.User, snapshot *agenttools.SkillSnapshot, cwd string) subagentProfile
```

职责：

1. `general`：
   - `SystemPrompt` 使用现有 `buildSubagentSystemPrompt` 的语义。
   - `ToolRegistry` 使用现有 `NewChildToolRegistry(s.Cfg, cwd)`。
   - `MaxRounds` 使用 `defaultSubagentMaxRounds`。
2. `explore`：
   - `SystemPrompt` 使用新的 `buildExploreSubagentSystemPrompt`。
   - `ToolRegistry` 使用新的 `NewExploreToolRegistry(s.Cfg, cwd)`。
   - `MaxRounds` 仍使用 `defaultSubagentMaxRounds`，不新增超时或轮次配置。

`runSubagent` 流程保持单入口：

1. 校验 `sub_type`。
2. 校验 `task`。
3. 解析 `cwd`。
4. 根据类型构建 profile。
5. 用 profile 初始化 child state、tools、system prompt。
6. 调用现有 `runSubagentLoop`。

### 5.4 `general` 子 Agent 系统提示词

现有 `buildSubagentSystemPrompt` 继续保留，用于 `general`。内容需要做轻微收敛，增加类型名并明确搜索任务边界：

```text
<subagent type="general">
你是由 `spawn_subagent` 派生出来的 general 子智能体。
你用于处理需要隔离上下文的综合分析或执行型子任务。搜索、文件定位、代码探索、实现梳理、证据收集等任务必须交给 explore 子智能体，不应由 general 子智能体承担。
...
</subagent>
```

原有规则保留：

1. 看不到父对话历史。
2. 只能依据当前任务和工作区文件工作。
3. 不要调用 `spawn_subagent`。
4. 不承担搜索专用任务；如任务本质是搜索或代码探索，应报告该任务应由 `explore` 执行。
5. 完成后输出简洁摘要。

### 5.5 `explore` 子 Agent 系统提示词

新增 `buildExploreSubagentSystemPrompt(cwd string) string`，不调用 `s.buildSystemPromptWithMemory`，避免继承主 Agent 的通用执行、写作、任务管理、记忆维护和编辑能力说明。

推荐提示词如下，实施时可作为 Go 字符串常量，或放入 `assets/prompts/explore_subagent.md` 并由 `FunctionalPrompts` 加载。考虑到这是子 Agent 类型的固定系统提示词，推荐放在代码中集中维护，避免额外引入功能 prompt 加载路径；如果后续类型增多，再迁移为嵌入式 prompt 文件。

```text
You are Cynosure's explore subagent, a read-only codebase search specialist.

=== READ-ONLY MODE ===
You must only inspect existing files and report findings. You must not create, modify, delete, move, copy, install, or persist files. You must not change repository, workspace, system, network, dependency, or package-manager state.

Your job:
- Rapidly locate relevant files, symbols, configuration, tests, docs, and implementation details.
- Read only the files needed to answer the caller's search request.
- Return a concise report with file paths, important line references when available, and confidence or gaps.

Tool rules:
- Prefer grep for content search.
- Prefer glob for filename pattern matching.
- Use ls only for known absolute directories.
- Use read_file when you already know the specific file path.
- Do not use write_file, edit_file, multi_edit, todo_write, update_memory, delete_memory, spawn_subagent, package managers, git mutation commands, or any state-changing operation.
- If bash is unavailable, do not ask for it; complete the search with grep, glob, ls, and read_file.

Environment:
- Current working directory: <resolved absolute cwd>
- Treat relative paths in the user task as relative to the current working directory unless the task gives an absolute path.
- Prefer reporting absolute paths or workspace-root-relative paths consistently; include enough path context for the parent agent to jump directly to the evidence.
- The parent conversation history is not available. Rely only on this task, the current working directory, and files you inspect.

Efficiency:
- Search broadly first, then read the smallest set of high-signal files.
- Run independent searches in parallel whenever the runtime supports it.
- Stop when you have enough evidence to answer the request; do not perform unrelated exploration.

Final response:
- Reply directly in normal text.
- Do not create files.
- Include key findings, evidence paths, and unresolved gaps.
```

项目适配点：

1. 不写“Claude Code”或外部产品身份，改为 Cynosure。
2. 使用项目真实工具名：`grep`、`glob`、`ls`、`read_file`。
3. 强调 `ls` 需要绝对路径，符合当前工具定义。
4. 明确当前工作目录由运行时注入为已解析绝对路径。
5. 不承诺不存在的 Glob/Grep/Read 大写工具名。
6. 不说“没有编辑工具”作为唯一约束，而是同时通过工具集限制保证。

### 5.6 `explore` 工具集

新增：

```go
func NewExploreToolRegistry(cfg config.AppConfig, cwd string) *ToolRegistry
```

该 registry 不从 `cfg.AllowedTools` 全量继承，而是取当前配置与 allowlist 的交集，确保用户配置没有启用某工具时不会被 `explore` 额外打开。

推荐允许工具：

1. `read_file`
2. `grep`
3. `glob`
4. `ls`
5. `web_fetch`（仅当用户当前配置允许；用于读取用户明确提供的 URL 内容时仍是只读）
6. `web_search`（仅当用户当前配置允许；如果本地未启用搜索后端，保持现有行为）
7. `load_skill`（不推荐给 explore，默认不允许；搜索型子 Agent 不应加载流程技能扩大行为面）
8. `read_persisted_output`（仅当后续子 Agent 压缩向 child 暴露该工具时再加入；本次不额外改变该行为）

本期实际推荐只允许前 4 个本地文件探索工具；`web_fetch` 与 `web_search` 不加入 `explore` 默认工具集，避免“文件搜索子 Agent”变成联网研究 Agent。若用户任务明确需要联网，仍由 `general` 或主 Agent 处理。

实现细节：

```go
var exploreToolNames = []string{"read_file", "grep", "glob", "ls"}
```

`NewExploreToolRegistry` 使用 `loadAllowedToolNames(cfg)` 得到当前配置启用的工具，再与 `exploreToolNames` 求交集，最后用 `buildToolDefinitions` 构造 definitions。这样：

1. 默认配置中这些工具可用时，`explore` 能搜索读取。
2. 用户显式移除某个工具时，`explore` 不会绕过配置重新启用。
3. 写入类工具、任务管理工具、记忆工具、`bash` 和 `spawn_subagent` 永远不会暴露给 `explore`。

### 5.7 主 Agent 系统提示词调整

需要同步修改：

1. `cynosure/assets/system_prompt.md`
2. `cynosure/internal/assistant/prompt.go` 中的 `DefaultBaseSystemPrompt`

调整重点放在“工具调用”域，保持七域结构不变：

1. 将现有 `spawn_subagent` 描述改为明确使用 `sub_type`。
2. 说明 `sub_type=explore` 用于只读代码搜索、文件定位、实现梳理、证据收集；搜索相关任务必须使用 `explore`。
3. 说明 `sub_type=general` 用于需要综合分析、执行多步骤工具链或可能涉及修改的隔离任务，但不得用于搜索相关任务。
4. 强调子 Agent 仍只返回最终摘要，不能再派生子 Agent。
5. 不把复杂类型定义写进静态 prompt 之外的地方；工具 schema 仍是参数真实性来源，静态 prompt 只教模型如何选择。

推荐文案：

```text
- `spawn_subagent` 必须提供 `sub_type` 与 `task`：搜索、文件定位、代码探索、实现梳理、证据收集等搜索相关任务必须使用 `sub_type=explore`；`sub_type=general` 仅用于需要隔离上下文的综合分析或执行型子任务，不得用于搜索相关任务。子智能体只返回最终摘要，不能再派生子智能体。
```

同时更新 README 中工具说明与系统提示词说明，避免文档仍描述旧参数。

### 5.8 参数校验与错误行为

校验层次：

1. JSON Schema：`sub_type` 必填，enum 为 `general`、`explore`。
2. 运行时：`parseSubagentType` 二次校验。
3. 工具注册：`spawn_subagent` 仍必须在主 Agent `AllowedTools` 中启用，否则返回现有未注册错误。

错误返回示例：

1. 缺少 `sub_type`：`invalid arguments for tool "spawn_subagent": ...`
2. 非法 `sub_type`：`invalid arguments for tool "spawn_subagent": ... value must be one of general, explore`
3. 运行时非法值：`Subagent failed: unsupported sub_type "x"; allowed values: general, explore`

### 5.9 冗余代码清理

本次调整会将子 Agent 类型分支集中到 profile 构造函数中，避免以下冗余：

1. 不在 `runSubagent` 多处散落 `if sub_type == ...`。
2. 不复制 `runSubagentLoop`；`general` 与 `explore` 共享同一 LLM/tool loop、审批、压缩、超时逻辑。
3. 不新增 `spawn_explore_subagent` 工具。
4. 不为 `explore` 单独复制工具执行器；只新增 registry 构造函数。
5. 不新增配置字段表达子 Agent 类型；类型是工具调用参数和代码内枚举。

## 6. 涉及文件

预计修改：

1. `cynosure/internal/tools/definitions.go`
   - 为 `spawn_subagent` 增加 `sub_type` schema、enum 和 required。
2. `cynosure/internal/agent/runtime/subagent.go`
   - 扩展 `spawnSubagentArgs`。
   - 新增子 Agent 类型常量和校验。
   - 新增 profile 构造。
   - 新增 `buildExploreSubagentSystemPrompt`。
   - 调整 `runSubagent` 根据类型选择 prompt 和 tools。
3. `cynosure/internal/agent/runtime/tool_registry.go`
   - 新增 `NewExploreToolRegistry`。
   - 增加工具名交集辅助函数，或复用现有 `withoutTool` 风格新增小函数。
4. `cynosure/internal/agent/runtime/runtime_test.go`
   - 更新旧的 `spawn_subagent` 参数测试，补充 `sub_type`。
   - 新增 `explore` prompt 与工具集测试。
   - 新增非法 `sub_type` 拒绝测试。
5. `cynosure/internal/tools/definitions_test.go` 或 `cynosure/internal/tools/validation_test.go`
   - 覆盖 `spawn_subagent` schema 中 `sub_type` 必填和 enum 校验。
6. `cynosure/assets/system_prompt.md`
   - 更新主 Agent 对 `spawn_subagent` 的使用说明。
7. `cynosure/internal/assistant/prompt.go`
   - 同步 `DefaultBaseSystemPrompt`。
8. `cynosure/internal/assistant/prompt_test.go`
   - 如现有 prompt 测试检查关键文案，补充或更新断言。
9. `cynosure/README.md`
   - 更新工具列表或 `spawn_subagent` 描述。

## 7. 测试策略

### 7.1 单元测试

1. `TestSpawnSubagentSchemaRequiresSubType`
   - 构造缺少 `sub_type` 的参数，调用 `ValidateToolArgs`，期望失败。
2. `TestSpawnSubagentSchemaRejectsUnknownSubType`
   - `sub_type=review`，期望 enum 校验失败。
3. `TestParseSubagentTypeAcceptsDefinedTypesOnly`
   - `general`、`explore` 成功；空字符串和未知值失败。
4. `TestBuildGeneralSubagentSystemPromptKeepsExistingRules`
   - 确认旧 `<subagent>` 核心规则仍存在。
   - 确认 `general` prompt 明确说明搜索相关任务必须交给 `explore`。
5. `TestBuildExploreSubagentSystemPromptIsReadOnlyAndSearchFocused`
   - 确认 prompt 包含 read-only、`grep`、`glob`、`read_file`、绝对路径/当前工作目录说明。
   - 确认 prompt 不包含主 Agent 的通用身份文案和记忆维护规则。
6. `TestNewExploreToolRegistryAllowsOnlyReadOnlySearchTools`
   - 当配置允许所有默认工具时，explore definitions 只包含 `read_file`、`grep`、`glob`、`ls`。
   - 确认不包含 `bash`、`write_file`、`edit_file`、`multi_edit`、`todo_write`、`update_memory`、`delete_memory`、`spawn_subagent`。
7. `TestRespondToConversation_SpawnExploreSubagentUsesFreshMessagesAndExploreTools`
   - 主 Agent 调用 `spawn_subagent`，参数为 `{"sub_type":"explore","task":"inspect workspace only","cwd":"."}`。
   - 子 Agent LLM request 只有 system + task，不泄漏父历史。
   - 子 Agent tools 只包含 explore 工具。
8. `TestRespondToConversation_SpawnGeneralSubagentPreservesExistingBehavior`
   - 旧的 fresh message、trace、无前端 tool event、压缩行为在 `sub_type=general` 下保持。

### 7.2 验证命令

实施后运行：

```bash
go -C cynosure test ./internal/tools ./internal/assistant ./internal/agent/runtime
go -C cynosure test ./...
go -C cynosure build ./...
```

如果 `go test ./...` 触发项目既有无关失败，需要记录失败用例与是否为基线问题；但本次改动涉及 runtime、tools、assistant，至少定向测试必须通过。

## 8. 风险与缓解

1. 风险：把 `sub_type` 设为 required 后，模型仍按旧格式调用 `spawn_subagent`。
   - 缓解：同步更新工具 schema、主系统提示词和相关测试；schema 错误会把模型拉回修正参数。
2. 风险：`explore` 只读限制只写在 prompt 中，实际仍能编辑。
   - 缓解：新增 `NewExploreToolRegistry`，运行时不暴露写入类工具。
3. 风险：`explore` 不复用主 prompt 后缺少环境路径说明。
   - 缓解：`buildExploreSubagentSystemPrompt` 动态注入 resolved cwd，并明确相对路径解释和 `ls` 绝对路径要求。
4. 风险：新增类型分支破坏 `general` 既有行为。
   - 缓解：`general` 继续使用现有 prompt 与 child registry；用回归测试覆盖现有 fresh message、压缩、禁止嵌套行为。
5. 风险：工具配置被绕过。
   - 缓解：`NewExploreToolRegistry` 使用当前配置允许工具与 explore allowlist 的交集，不额外打开用户禁用的工具。

## 9. 需求覆盖自检

1. 新增 `sub_type` 参数：已覆盖，schema 与 Go args 均调整。
2. 新增 `explore` 子 Agent：已覆盖，包含类型、prompt、工具集。
3. `explore` 单独系统提示词：已覆盖，不复用主 Agent 完整 prompt。
4. 结合项目改造和优化参考 prompt：已覆盖，替换为 Cynosure 身份、真实工具名、绝对路径和 env 说明。
5. 参数校验只允许定义的 `sub_type`：已覆盖，schema enum + runtime parser 双层校验。
6. 调整主 Agent 系统提示词：已覆盖，要求同步 `assets/system_prompt.md` 与 `DefaultBaseSystemPrompt`。
7. `general` 可保留但搜索任务必须交给 `explore`：已覆盖，主提示词和 `general` 子提示词都会明确该边界。
8. 其他策略不变：已覆盖，`general` 保持现有执行行为，审批/压缩/超时/禁止嵌套共用原逻辑。
9. 删除冗余代码：已覆盖，集中 profile 构造，不新增重复 loop 或新工具。
