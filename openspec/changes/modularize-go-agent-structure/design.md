## Context

`go-agent` 当前已经收敛为一个浏览器优先的 Web Agent 服务，但多个核心文件在同一处同时承载启动装配、HTTP 处理、运行时编排、工具注册、配置解析、存储访问和辅助函数，导致单文件体量偏大、包边界逐渐模糊。典型热点包括 `internal/web/app/server.go`、`internal/web/runtime/runtime.go`、`internal/tools/registry.go`、`internal/config/config.go`、`internal/web/storage/store.go` 与 `internal/sessions/skill.go`。

这次变更的核心约束有三点：

1. 不改变现有对外行为，包括 Web API、聊天流程、工具调用边界、workspace 约束和部署方式。
2. 以保守重构为主，优先在现有包内拆文件和整理职责，而不是一次性大规模调整 import 路径或引入新框架。
3. README 需要同步到重构后的真实结构，确保文档仍能指导开发者理解模块边界与验证方式。

## Goals / Non-Goals

**Goals:**
- 将 `go-agent` 中职责混杂的大文件拆分为更聚焦的模块文件，降低单文件复杂度。
- 明确 `app`、`runtime`、`tools`、`config`、`storage`、`sessions` 等层的职责边界，减少跨层辅助逻辑混写。
- 保持现有 API、runtime 行为、配置语义、部署产物结构与测试基线不变。
- 更新 `go-agent/README.md`，使目录结构、模块说明和开发验证步骤与重构后的实现一致。

**Non-Goals:**
- 不新增用户可见功能，不修改已有 HTTP 接口契约。
- 不调整数据库 schema、Redis key 语义或会话存储模型。
- 不引入新的外部依赖或新的应用架构范式。
- 不在本次重构中强制完成所有包重命名，例如 `internal/sessions` 改名为 `internal/skills` 之类的高扰动调整。

## Decisions

### 1. 优先采用“包内拆文件”而不是“跨包迁移”

**选择：** 对热点模块优先保持原有 package 名称和对外构造方式不变，通过新增同包文件把装配逻辑、handler、辅助函数、仓储方法、路径解析和工具定义按职责拆开。

**原因：** 这能最大程度降低行为回归风险，也能减少一次性修改 import、测试桩和初始化路径带来的扰动，符合“结构优化但功能不变”的目标。

**备选方案：**
- 直接重命名包并重构目录层级——可读性提升更彻底，但影响范围过大，不适合本次保守重构。
- 完全维持当前结构，仅补充注释——无法解决单文件过大与职责混杂问题。

### 2. 以垂直职责拆分 `internal/web/app`

**选择：** 让 `server.go` 保留 `Server` 构造和启动相关逻辑，将路由注册、鉴权 handler、Skill handler、Conversation handler、HTTP 通用辅助函数拆到独立文件。

**原因：** `internal/web/app/server.go` 当前既是装配入口，也是主要 handler 实现载体。按 HTTP 责任面拆分后，可以在不改动路由行为的前提下降低阅读和测试成本。

**备选方案：**
- 按资源类型再拆新 package，例如 `handlers/skills`、`handlers/conversations`——后续可行，但当前会引入更多跨包依赖与可见性问题。

### 3. 将 `runtime` 与 `tools` 的边界收敛为“编排层”和“通用工具层”

**选择：** `internal/web/runtime` 负责会话编排、tool 调用流程、SSE 输出与 `load_skill` 这类 Web runtime 专属逻辑；`internal/tools` 负责通用 tool schema、运行时环境、路径安全和基础文件/命令工具 handler。

**原因：** 当前 `internal/web/runtime/tools.go` 与 `internal/tools/registry.go` 存在概念重叠。先明确两层角色，可以在不重写工具系统的前提下减少职责重复。

**备选方案：**
- 将所有工具逻辑统一并到单一包中——理论上更整齐，但本次会造成过大改动面。

### 4. 将 `config`、`storage`、`sessions` 的重构限定为“文件分层”

**选择：**
- `internal/config` 按 LLM 配置、Web 配置、runtime 路径、layout 校验、环境变量辅助拆分。
- `internal/web/storage` 按 store 初始化、migration、users/sessions/skills/conversations/messages/toolcalls/cache 拆分。
- `internal/sessions` 按 loader、frontmatter、merge、render 拆分。

**原因：** 这些模块对外 API 已经被多处依赖。先拆分内部文件结构，能改善维护性，同时避免触发更高风险的跨模块重构。

**备选方案：**
- 直接改包名或建立 repository 接口层——收益更大，但会显著提高本次变更复杂度。

### 5. README 跟随“职责边界”而不是“实现细节”更新

**选择：** README 重点更新目录结构、核心模块职责说明与开发验证章节，不把文档写成逐文件索引，而是反映稳定的模块边界。

**原因：** 如果文档过度依赖细碎文件名，后续继续整理时很快又会过期；按模块边界说明更稳定，也更符合对外认知。

**备选方案：**
- 仅更新目录树不补充职责说明——可读性不足，无法真正解释这次模块化重构的价值。

## Risks / Trade-offs

- **[Risk] 拆文件过程中出现遗漏的私有辅助函数或初始化顺序变化** → 通过 `go test ./...` 和构建流程校验行为等价，并保持公开构造器与主要调用顺序不变。
- **[Risk] runtime 与 tools 拆分后责任边界仍有灰区** → 以“不改变 tool 对外定义与执行结果”为前提，先完成最小职责收敛，避免一次性过度抽象。
- **[Risk] README 更新后仍与真实目录轻微漂移** → 在完成代码重构后再统一核对目录树与章节描述，确保文档基于最终实现。
- **[Trade-off] 保守拆分会保留部分历史命名，如 `internal/sessions`** → 接受短期命名不完美，换取更低的回归风险；后续如有需要可单独提案继续演进。

## Migration Plan

1. 先在现有包内完成文件拆分，保持对外构造器、函数签名和包路径稳定。
2. 运行现有 Go 测试，修复因私有符号移动、依赖注入方式或测试桩导致的问题。
3. 执行构建或打包流程，确认部署产物与运行行为没有变化。
4. 基于最终结构更新 `go-agent/README.md`。
5. 如出现回归，可按模块粒度回滚对应拆分文件；本次变更不涉及数据迁移。

## Open Questions

- 当前不计划在本次重构中改动 `internal/sessions` 的包名；若实施中发现命名已经明显阻碍代码组织，可在保持对外行为不变的前提下评估是否单独补充后续提案。
