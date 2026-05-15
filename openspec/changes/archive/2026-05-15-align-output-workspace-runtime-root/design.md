## Context

当前 go-agent 的路径模型已经同时存在源码态 `workspace/` 与打包态 `output/workspace/` 两套目录，但运行时没有稳定区分“部署包根目录”和“源码根目录”。结果是 Web runtime、builtin skill 加载与 bash/file 工具可能继续把源码根下的 `workspace/` 当作默认 cwd，导致云端部署时无法命中 `output/workspace/skills` 与 `output/workspace/bin`，而本地调试又依赖进程启动位置是否恰好在源码根目录。

这次变更需要把目录解析逻辑从“依赖 cwd”改成“依赖显式 runtime root 解析”，并让云端部署与本地调试共用一套优先级明确的路径模型：部署态优先使用打包产物根目录下的 `workspace`，本地未打包时再回退到源码目录下的 `workspace`。

## Goals / Non-Goals

**Goals:**
- 让 runtime 能稳定识别部署态的 `output/workspace`，并把它作为 agent 默认 cwd
- 让 builtin skills、命令二进制与脚本命令共享同一套 runtime root 解析结果
- 保持本地源码调试可用，在未生成部署包时自动回退到源码 `workspace`
- 减少模块各自拼路径的行为，把路径决策收敛到配置层或共享解析函数
- 为部署优先级、回退顺序和安全边界补充可回归测试

**Non-Goals:**
- 不重做整套多租户用户 workspace 隔离模型
- 不引入容器化沙箱、远程执行器或新的命令分发系统
- 不改变数据库中用户自定义 skill 的 CRUD 语义
- 不扩展新的工具种类，只修正已有 runtime/tool/skill 对目录的绑定方式

## Decisions

### 1. 引入统一的 runtime workspace 解析优先级

**选择：** 将 `AppHome` 视为服务运行根，并在其下按优先级解析 runtime workspace：
1. 若配置或环境变量显式给出 `WORKSPACE_ROOT`，优先使用该值；
2. 否则若 `AppHome/output/workspace` 存在且包含部署资产，使用 `output/workspace`；
3. 否则回退到 `AppHome/workspace` 作为本地调试目录。

`BuiltinSkillsDir`、`CommandBinDir`、`CommandScriptDir` 的默认值都从解析后的 runtime workspace 推导，而不是分别独立相对 `AppHome` 拼接。

**备选：**
- 始终固定使用 `AppHome/workspace`
- 始终固定使用 `AppHome/output/workspace`
- 继续依赖进程 cwd 推断路径

**理由：** 这样能兼容云端部署包和本地源码调试，同时避免 runtime 不同模块各自选了不同根目录。

### 2. 将 builtin skill 加载与工具 env 注入绑定到同一个 runtime root

**选择：** Server 启动时先解析一次 runtime workspace，并把结果同时用于 builtin skill loader、runtime env 注入、默认 cwd 与命令路径暴露。后续对话回合不再自行猜测 skill 根目录或命令目录。

**备选：**
- builtin skill loader 单独用 `BuiltinSkillsDir`，工具系统单独用 `WorkspaceRoot`
- 每次工具执行时临时重新扫描目录

**理由：** skill 加载、bash/file 工具和命令调用如果绑定不同根目录，就会再次出现“能看到 skill 但执行不到 cmd”或“cwd 对了但 skill 目录不对”的问题。

### 3. 使用“部署优先、本地回退”的资产识别规则

**选择：** 当 runtime workspace 指向 `output/workspace` 时，默认认为 `skills/`、`bin/`、`cmd/` 都来自打包产物；当回退到源码 `workspace` 时，保留当前本地调试体验。路径识别依赖目录是否存在及是否包含期望资产，而不是仅凭字符串名称判断。

**备选：**
- 仅根据路径名中是否含有 `output` 决定模式
- 本地调试也强制要求先执行打包

**理由：** 大厂常见服务部署方式都是“产物目录优先、源码目录只用于开发态”，目录存在性检测比约定启动位置更稳健。

### 4. 保持工具边界始终受 runtime workspace 约束

**选择：** 无论命中部署态还是本地调试态，bash/file 等工具看到的默认 cwd 都是解析后的 runtime workspace，并继续沿用现有越权校验，不允许越过该根目录访问宿主机其他路径。

**备选：**
- 部署态限制在 `output/workspace`，本地调试态放宽到源码根
- 允许 skill 或 bash 自己再切换 cwd 到任意绝对路径

**理由：** 兼容本地调试不应以放弃安全边界为代价；切换模式只影响根目录选择，不影响边界 enforcement。

## Risks / Trade-offs

- **[Risk] `output/workspace` 目录存在但产物不完整，导致误判为部署态** → 通过目录探测时同时校验关键子目录，如 `skills`、`bin` 或配置中显式覆盖
- **[Risk] 历史测试夹具默认假设 `workspace` 固定在源码根** → 同步更新测试夹具，分别覆盖部署优先与本地回退两种分支
- **[Risk] 配置项之间出现双重语义，导致 `BuiltinSkillsDir` 与 `WorkspaceRoot` 不一致** → 将默认值统一改为从解析后的 runtime workspace 派生，仅保留显式覆盖能力
- **[Trade-off] 自动探测提升了兼容性，但增加了少量启动期路径判断逻辑** → 把逻辑集中在配置层，避免在 runtime 各处复制判断

## Migration Plan

1. 在配置层新增或重构共享解析函数，统一得到 runtime workspace 根目录
2. 调整 `LoadWebConfig` 默认值推导，让 skills/bin/cmd 目录都从该根目录派生
3. 更新 server 启动逻辑和 runtime env 注入逻辑，统一消费解析后的路径
4. 更新 builtin skill 加载与工具测试，验证部署态优先命中 `output/workspace`
5. 回滚时恢复原有默认路径拼接策略即可；由于保留显式环境变量覆盖，不需要数据迁移

## Open Questions

- 当前是否需要把“部署态判定”进一步暴露为显式配置开关；本次先采用目录探测 + 环境变量覆盖，后续若上线环境有歧义再补充
