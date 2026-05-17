## Context

当前 `go-agent` 已经把默认入口切到 Web 服务，但代码结构仍保留大量历史兼容层：根入口与 `cmd/web` 重复、`internal/agent` 与 `internal/sessions` 中仍包含 REPL / memory / compact / subagent 等 CLI 时代能力。与此同时，`cmd/build-artifacts`、`internal/deploy` 以及 `output/workspace` / `workspace` 相关布局并不等同于 CLI 历史包袱，其中 `workspace/bin` 与 `workspace/cmd` 仍然是 Web agent 执行二进制命令和脚本能力所依赖的运行时路径；当前真正的问题是这些运行时资产与 CLI 兼容结构、重复入口和文档叙事混杂在一起，导致职责边界不清。

这些结构不是“Web 服务可用”的必要条件，却持续提高了目录复杂度、部署复杂度和后续维护成本。尤其是当前云端部署说明依赖 `output/` 打包模型，但数据库迁移文件仍从源码路径读取，这说明现有部署约定已经出现服务本体与历史打包逻辑不一致的问题。

## Goals / Non-Goals

**Goals:**
- 将 `go-agent` 收敛为单一 Web 服务程序，删除与 CLI/REPL 直接相关的代码结构
- 保留浏览器聊天服务运行所需的最小后端模块：启动入口、HTTP server、配置、鉴权、存储、日志、聊天 runtime，以及 Web agent 所需的 runtime workspace 资产与其解析链路
- 简化目录职责，避免重复入口、历史 session 能力与真正服务 Web agent 的运行时命令资产继续混在同一层叙事中
- 简化部署方式，提供清晰的本地运行和云服务器运行两条最短路径
- 更新 `README.md`，使其反映新的目录结构、运行约定与部署步骤

**Non-Goals:**
- 不重写现有登录、会话、数据库表结构或前端产品
- 不在本次设计中引入 Docker、Kubernetes、Terraform 等新的部署体系
- 不追求保留所有历史 skill/tool/runtime 兼容行为，只保留 Web 服务实际需要的能力
- 不要求在本次变更中新增新的产品功能；重点是删除冗余结构并收敛已有能力

## Decisions

### 1. 保留单一 Web 服务入口，删除重复启动路径

**选择：** 保留一个正式启动入口作为唯一推荐方式，并移除等价的重复入口与围绕它形成的文档分叉。

**备选：**
- 同时保留 `main.go` 与 `cmd/web/main.go`
- 新增第三套 `cmd/server/main.go` 再做迁移

**理由：** 用户目标是“只保留启动 Web 服务所需代码”，因此应先消除最直接的结构重复。继续保留双入口只会让 README、构建脚本和部署说明继续重复；新增第三入口则会增加无意义迁移成本。

### 2. 以“服务本体必需依赖”为裁剪标准

**选择：** 仅保留 Web 服务启动链路真实使用到的模块：配置、日志、存储、鉴权、Web runtime 以及这些模块的直接依赖；删除 CLI agent、REPL、memory/compact/subagent 等不参与 Web 主链路的代码，但保留 Web agent 执行命令所需的运行时命令资产目录与其构建/解析链路。

**备选：**
- 继续保留历史代码，只在 README 中弱化
- 仅移动目录位置，不真正删除无关模块

**理由：** 本次重构目标是减少维护面，而不是做表面整理。只有以“是否被 Web 服务启动链路直接依赖”作为裁剪标准，才能真正收敛项目复杂度。

### 3. 保留 deployment-aware workspace 布局，但统一其作为 Web runtime 资产目录的语义

**选择：** 保留 `output/workspace` 与 `workspace` 的 deployment-aware 布局，用于承载 Web agent 所需的 `workspace/bin`、`workspace/cmd` 与相关运行时资产；但需要把它们明确为 Web runtime 的正式依赖，而不是 CLI 兼容残留。同时收敛入口和文档表达，避免把“本地调试回退目录”和“部署产物目录”解释成两套产品形态。数据库迁移等服务必需资源仍应直接随服务二进制或服务包稳定交付，例如使用嵌入资源或稳定的服务内路径。

**备选：**
- 完全移除 `output/workspace` / `workspace` 双轨解析，改为单一路径布局
- 保留当前解析逻辑，但不澄清这些目录对 Web agent 的实际职责

**理由：** `workspace/bin` 和 `workspace/cmd` 并不是可随意删除的历史包袱，而是 agent 在 Web 服务内执行二进制命令/脚本能力的运行时依赖。真正需要收敛的是它们的语义和装配方式：保留这套目录，但把它们从 CLI 兼容叙事中解耦，并让本地与云端部署都围绕同一套 runtime 资产模型来说明。

### 4. Web runtime 保留 agent 所需命令产物能力，但与 CLI/REPL 结构解耦

**选择：** 浏览器聊天默认走“直接回答优先”的运行时，但当 agent 需要执行受控二进制命令或脚本时，仍应通过 `workspace/bin` 与 `workspace/cmd` 中的运行时资产完成。这些目录属于 Web runtime 的正式能力边界，而不是 CLI/REPL 的副产品；需要剥离的是旧的终端交互结构，而不是 agent 命令执行能力本身。

**备选：**
- 将 `workspace/bin`、`workspace/cmd` 重构为其他等价的 runtime 资产目录命名或布局，但继续保留 agent 命令执行能力
- 保留现有命令目录，但继续让其实现与 CLI 结构、历史兼容路径强耦合

**理由：** 用户要求是“仅保留启动 Web 服务所需代码”，而 `workspace/bin`、`workspace/cmd` 恰恰属于 Web agent 的有效运行时依赖，不应被误删。需要移除的是 CLI/REPL 入口、无入口历史 session 能力和错误的产品叙事，而不是 Web agent 为执行命令所需的命令资产路径。

### 5. README 以“本地运行”和“云服务器部署”两条主线重写

**选择：** README 聚焦两种标准场景：
1. 本地开发：准备依赖、配置环境变量、直接启动服务
2. 云端部署：构建单一正式服务发布包（包含 Web 服务二进制与所需 runtime workspace 资产）、上传到服务器、配置环境变量、通过 systemd 或等价方式运行

**备选：**
- 保留 build-artifacts、源码态 workspace、部署态 workspace 的并列说明
- 继续在同一文档中混合介绍调试入口、兼容入口和历史脚本

**理由：** 用户明确要求优化本地与云端部署步骤，因此文档必须围绕“最短成功路径”组织，而不是围绕兼容性历史解释。

## Risks / Trade-offs

- **[Risk] 删除历史模块后，少量依赖 skill/tool 的隐式调用路径可能暴露出来** → 先以 Web 服务真实依赖链为基准梳理引用，再补充编译与测试验证，确保删除的是无入口死路径而非仍被 runtime 间接依赖的代码
- **[Risk] 现有部署脚本或运维习惯依赖 `output/` 目录** → 在 README 中明确迁移到单一正式服务发布包模型，但说明该发布包仍需包含 `output/workspace` 或等价 runtime workspace 资产，并在必要时保留一个过渡期兼容脚本但不再作为主路径
- **[Risk] 迁移文件与静态资源在新布局下找不到** → 将服务必需资源切换为嵌入式或稳定的服务内相对路径，避免依赖源码工作目录
- **[Trade-off] 牺牲历史 CLI/REPL 能力，但保留 agent 命令资产链路** → 接受该取舍，因为本次目标是让仓库服务于 Web 服务本体，同时保留 Web agent 必需的命令执行能力
- **[Trade-off] 短期内需要较多删除和移动操作** → 通过先定义保留边界、后按模块清理的方式降低重构风险

## Migration Plan

1. 明确 Web 服务最小依赖链（包含 runtime workspace 资产依赖），并删除重复入口与无入口历史模块
2. 调整配置与启动逻辑，明确 `output/workspace` / `workspace` 是 Web runtime 资产目录，并保持其解析规则与部署说明一致
3. 处理数据库迁移和其他服务必需资源的装载方式，确保本地与部署运行一致
4. 简化构建/发布脚本，仅保留 Web 服务实际需要的产物，并继续正确产出 runtime workspace 资产
5. 运行 Go 测试与最小启动验证，确认登录、建会话、发送消息等核心流程不受影响
6. 更新 `README.md`，提供新的目录说明、本地部署步骤与云服务器部署步骤
7. 若需要回滚，可恢复旧入口与旧部署脚本，但不恢复已确认无依赖的 CLI/REPL 模块

## Open Questions

- 是否完全移除 builtin skill 目录与相关加载逻辑，还是仅将其降级为非默认、显式配置的可选能力？当前设计倾向于“非默认且不阻塞启动”
- 构建产物是保留精简版 `build.sh` 统一生成 Web 服务二进制与 runtime 资产，还是拆分为更明确的构建步骤？当前设计要求无论采用哪种方式，都必须继续正确产出并打包 `workspace/bin` 与 `workspace/cmd`
