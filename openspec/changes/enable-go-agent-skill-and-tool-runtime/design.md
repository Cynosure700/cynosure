## Context

当前默认入口已经是 Web 服务，`main.go` 和 `cmd/web/main.go` 都直接启动 Web server，而不是进入 REPL。与此同时，真正在线的 Web runtime 只暴露 `load_skill`，对 shell、本地文件和工作区请求会直接返回能力边界说明；旧版 CLI agent loop 虽然保留了 bash、文件读写、todo、task、compact 等工具能力，但并没有接到 Web runtime 上。

另外，skill 来源目前是割裂的：`go-agent/skills` 目录下的默认 skill 只会在旧 REPL 路径中加载，Web runtime 则仅在每轮对话中从数据库查询当前用户已启用的自定义 skill。配置中虽然已有 `WorkspaceRoot`，但尚未真正接线到 Web runtime，因此当前服务没有一个稳定、可部署的服务端工作目录模型。更重要的是，skill 后续可能需要调用部署时生成的命令产物，例如 `cmd` 目录下编译出的二进制文件或随包发布的 `.py` 脚本；如果没有固定的应用根目录与命令产物目录，这类引用会高度依赖进程 cwd，无法可靠部署。

本次变更要把这些分散能力收敛为一个统一的服务端 agent：浏览器聊天继续作为主入口，但运行时要能合并“内置 skill + 用户 skill”，并在用户隔离的 workspace 中调用平台允许的 tool。

## Goals / Non-Goals

**Goals:**
- 让 Web runtime 能同时使用 `go-agent/skills` 下的内置 skill 与数据库中的用户自定义 skill
- 让所有用户都能调用内置 skill，同时继续保持用户自定义 skill 的所有权隔离与启停控制
- 将现有 go-agent 工具注册表接入 Web runtime，使 agent 能在服务端调用已授权 tool
- 为部署实例定义固定的应用根目录，并将 skills、cmd 命令产物、workspace 等路径都锚定到该根目录下
- 为每个用户定义稳定且隔离的服务端 workspace，并为工具执行提供默认工作目录
- 明确 go-agent 的部署方式为长期运行的 Web 服务，而不是依赖 CLI REPL 的临时进程 cwd
- 明确 `cmd` 目录中的 Go 程序在部署时需要预编译为二进制产物，脚本类文件需要按固定路径随包发布

**Non-Goals:**
- 不引入分布式远程执行器或多机任务调度
- 不实现用户自定义上传本地代码仓库或绑定任意外部目录的能力
- 不重做整套 skill 存储模型，数据库中的用户 skill 仍保持当前 owner + enabled 状态管理方式
- 不新增 skill 市场、共享社区或权限细分体系
- 不在首版中支持用户动态上传新的服务端可执行文件；`cmd` 产物仍由平台构建和发布

## Decisions

### 1. 统一 skill 目录模型：启动时加载内置 skill，请求时加载用户 skill，再做合并

**选择**：新增一个面向 Web runtime 的组合 skill catalog。

- 服务启动时一次性扫描并加载 `go-agent/skills` 下的内置 skill，作为全局只读 catalog 常驻内存
- 每次对话 turn 时再查询当前用户已启用的数据库 skill，构造成用户私有 catalog
- runtime 对模型暴露的是“合并后的 skill 集合”，描述统一注入 prompt，正文仍通过 `load_skill` 惰性读取
- 内置 skill 标识符保留为保留名，数据库中的自定义 skill 不允许与其冲突

**备选**：
- 完全复用旧 `LoadAll()`，让 Web runtime 也依赖进程 cwd 扫描本地目录
- 将数据库 skill 落地到文件系统，再走同一套文件扫描逻辑

**理由**：启动时缓存内置 skill 可以避免每轮重复扫盘；数据库 skill 继续按用户动态查询，可以保持多租户隔离；组合 catalog 模式能最大程度复用现有 `load_skill` 交互方式，而不用改写模型侧心智。

### 2. Web runtime 复用现有工具注册表，但增加显式的 Web 允许列表

**选择**：保留 `go-agent/internal/tools` 里的工具定义与 handler，实现一个 Web runtime adapter，把“已注册工具”映射为“Web 部署允许的工具子集”。

- 工具定义优先复用已有 parent tool defs 和 handlers
- Web runtime 在构造工具列表时只暴露显式允许的工具名称
- 工具执行时统一注入 user、conversation、workspace 等上下文，再调用底层 handler
- 对拒绝执行的工具，不做静默失败，而是把拒绝结果作为 tool result 回填 runtime loop，并写入审计记录

**备选**：
- 在 `internal/web/runtime` 中重新实现一套新的工具系统
- 继续保持 `load_skill` 为唯一工具，其余能力全部通过 prompt 解释代替

**理由**：复用现有注册表可以减少重复实现，也便于与 CLI agent 保持工具语义一致；显式 Web allowlist 则能防止“注册了就默认全开”的风险。

### 3. 工作目录模型采用“每用户一个 workspace 根目录，默认 cwd 等于该根目录”

**选择**：区分“部署应用根目录”和“用户 workspace 根目录”。以配置项 `AppHome` 作为部署实例的固定根目录，在其下解析 skills、cmd 命令产物和数据目录；以 `WorkspaceRoot` 为用户工作区总根目录，为每个用户分配独立子目录 `WorkspaceRoot/<user_id>/`，并把它作为该用户所有 Web 会话的默认工作目录。

- 新用户在第一次触发需要 workspace 的对话或工具调用时创建目录
- shell / 文件类工具中的相对路径都相对于该默认 cwd 解析
- 任何显式路径或 cwd 解析结果都必须落在该用户 workspace 根目录之下
- skill 如需调用平台随部署发布的命令产物，必须使用 `AppHome/bin` 或 `AppHome/cmd` 下的固定绝对路径，而不是依赖共享 cwd 寻址
- MVP 不引入“每个会话一个不同 cwd”的额外状态；如后续需要，可在 conversation metadata 中扩展

**备选**：
- 沿用当前进程 cwd 作为所有用户共享目录
- 为每个 conversation 都创建独立目录并持久化 cwd

**理由**：按用户建 workspace 足以满足隔离要求，且不会引入额外数据模型复杂度；同时把部署资源从用户工作区中拆开，能避免 skill 调用共享二进制或脚本时污染用户目录，也避免部署目录与用户 cwd 混用。

### 4. 部署模型统一为长期运行的单服务 Web agent，并在启动时解析绝对路径配置

**选择**：go-agent 继续以 Web 服务形式部署，启动时将应用根目录、builtin skills 目录、cmd 命令产物目录和 workspace 根目录全部解析成绝对路径，并在 server 初始化阶段完成依赖检查。

- `main.go` / `cmd/web/main.go` 继续作为等价入口，均启动同一套 Web 服务
- 应用根目录默认指向部署包中的 `go-agent/` 根目录
- builtin skills 目录默认指向 `AppHome/skills`
- workspace 根目录默认指向 `AppHome/data/workspaces`
- `cmd` 目录中的 Go 命令在构建阶段编译到 `AppHome/bin`，`.py` 等脚本类辅助命令复制到 `AppHome/cmd`
- 启动阶段若关键目录配置不可用或命令产物缺失，应直接报错，避免服务进入半可用状态

**备选**：
- 继续依赖“从 go-agent 根目录启动”的隐式 cwd 约定
- 恢复 REPL 为主入口，再由 REPL 间接复用 Web runtime 组件

**理由**：将关键路径在启动时归一化成绝对路径，能降低部署时对 cwd 的隐式依赖；把 Go 命令预编译成二进制后统一放到部署目录，也能让 skill 在运行时稳定调用，不必依赖源码环境或临时编译。

### 5. 命令产物采用“平台构建、只读发布、技能按固定路径调用”

**选择**：将 `go-agent/cmd` 视为部署命令源目录，其中可运行的 Go 子命令在构建阶段编译为二进制文件，脚本类资源按原文件发布；运行时只允许 skill 或 tool 调用部署包中的只读命令产物。

- Go 命令产物统一输出到 `AppHome/bin/<name>`
- Python 或其他脚本文件保留在 `AppHome/cmd/...`，由运行环境提供解释器执行
- skill 文本中引用这些命令时，应使用平台约定的绝对路径或由 runtime 注入的只读命令根目录变量
- 不允许用户把自定义可执行文件写入共享命令目录

**备选**：
- 运行时首次调用时再临时编译 `cmd` 目录
- 把用户 workspace 同时作为平台命令安装目录

**理由**：预编译能减少运行时不确定性和延迟；把共享命令目录保持为只读资源，有助于防止用户互相污染部署产物，也更符合可审计的服务端发布流程。

### 6. 工具审计记录要补充 workspace 与命令产物维度，便于排查越权与路径问题

**选择**：扩展现有工具审计语义，至少记录 user id、conversation id、tool name、resolved cwd、命令产物路径（如有）、状态、摘要/拒绝原因。

**备选**：
- 继续只记录工具名与简单结果

**理由**：引入服务端 workspace 和共享命令产物目录后，很多问题都会与路径解析、产物缺失和越权校验相关；不记录 cwd、命令路径和拒绝原因将很难排查线上问题。

## Risks / Trade-offs

- **[Risk] 旧工具 handler 默认依赖进程 cwd，直接接入 Web 可能绕过用户隔离** → 在 Web adapter 层统一做 workspace 解析和路径重写，禁止底层直接使用共享 cwd
- **[Risk] 内置 skill 与数据库 skill 命名冲突会导致 `load_skill` 行为不确定** → 为内置 skill 预留标识符，创建或更新自定义 skill 时做冲突校验
- **[Risk] 将更多工具开放给 Web runtime 会扩大攻击面** → 仅暴露显式 allowlist 中的工具，并在每次调用前执行 scope 校验和审计
- **[Risk] `cmd` 产物在部署阶段未成功构建或脚本解释器缺失，skill 会在运行时失败** → 启动前增加构建/打包校验，启动时校验关键命令产物存在性
- **[Trade-off] 启动时预加载内置 skill 会增加少量启动成本** → 换取对话时更低延迟与更稳定的 catalog 行为
- **[Trade-off] 采用“每用户一个默认 cwd”简化了 MVP，但会限制会话级目录切换灵活性** → 先保证隔离与可部署性，后续如有需要再增加 conversation-level cwd 状态
- **[Trade-off] 命令产物由平台统一构建，会增加部署流水线步骤** → 换取稳定的运行时路径和更可控的安全边界

## Migration Plan

1. 在配置层补充并规范应用根目录、builtin skills 目录、cmd 产物目录与 workspace 根目录的绝对路径解析
2. 为 `cmd` 目录建立构建与打包步骤：编译 Go 子命令到 `AppHome/bin`，复制脚本类文件到 `AppHome/cmd`
3. 抽象内置 skill catalog，并在 Web server 启动时加载 `AppHome/skills`
4. 在 runtime 中引入“内置 skill + enabled 用户 skill”的合并装配逻辑，保持 `load_skill` 的调用方式不变
5. 为 Web runtime 增加工具适配层，将现有注册工具映射到 Web allowlist，并接入用户 workspace 与命令产物目录上下文
6. 为文件/命令类工具补充路径校验和默认 cwd 解析，确保用户访问落在 `WorkspaceRoot/<user_id>/`，共享命令访问落在只读的 `AppHome/bin` 或 `AppHome/cmd`
7. 扩充工具审计记录与相关测试，覆盖成功调用、拒绝调用、skill 合并、内置 skill 冲突、命令产物缺失等场景
8. 部署时继续使用 Web 服务入口；若回滚，可退回当前仅支持 `load_skill` 的 Web runtime，而不影响数据库中的用户 skill 数据

## Open Questions

- Web 首版是否默认开放全部现有注册工具，还是先从文件类 + `load_skill` + `todo/task` + 少量受控命令调用开始分阶段放开？
- 是否需要在 API 层显式展示“内置 skill 列表”，还是仅在 runtime 中作为隐式共享能力注入？
- 未来如果要支持 conversation-level cwd 切换，是否要把 cwd 持久化到 `conversations` 表，而不是仅依赖用户根目录默认值？
