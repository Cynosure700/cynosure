## Context

当前实现已经具备服务端 workspace、内置 skill 加载和命令产物构建能力，但部署目录约定仍然是分散的：配置层同时维护 `AppHome`、`BuiltinSkillsDir`、`CommandBinDir`、`CommandScriptDir` 与 `WorkspaceRoot`，默认路径还混用了 `workspaces/` 与 `workspace/`。这会让部署产物结构不稳定，也让 runtime 对 skill、命令二进制和工具 cwd 的定位依赖多个目录约定。

用户期望的部署模型更简单：服务发布后，所有运行时资产都放在统一的 `workspace/` 根目录下；内置 skills 位于 `workspace/skills`，`cmd` 下构建出的二进制位于 `workspace/bin`，服务初始化时直接从这些目录加载；agent 拥有终端命令与文件类工具，但执行边界必须限定在服务端 workspace，而不是宿主机用户目录。

## Goals / Non-Goals

**Goals:**
- 将运行时资产目录收敛为统一的 `workspace/` 根目录约定
- 明确服务启动时从 `workspace/skills` 加载内置 skills，并把其作为默认共享 catalog
- 明确 `cmd` 下的 Go 入口在构建阶段编译到 `workspace/bin`
- 让终端命令与文件类工具默认在 workspace 内执行，并拒绝落到宿主机任意目录
- 提供标准 `build.sh`，统一生成部署输出目录与 workspace 资产结构

**Non-Goals:**
- 不改变 Web 聊天、会话存储或数据库中用户自定义 skill 的基本数据模型
- 不引入新的远程执行器、容器沙箱或多机分发能力
- 不支持用户在运行时动态上传新的平台级二进制到 `workspace/bin`
- 不在本次变更中扩展新的工具类型，只调整已有工具的目录边界与路径注入

## Decisions

### 1. 以单一 deployment workspace 作为运行时资产根目录

**选择**：保留 `AppHome` 作为部署包根目录，但将 runtime 可见资产统一锚定到 `AppHome/workspace`。`BuiltinSkillsDir` 默认解析为 `workspace/skills`，`CommandBinDir` 默认解析为 `workspace/bin`，脚本类辅助命令默认解析为 `workspace/cmd`。

**备选：**
- 继续保留 `AppHome`、`workspaces/skills`、`workspaces/bin`、`workspace/` 并行存在的目录模型
- 将所有路径完全硬编码为进程 cwd 相对目录

**理由：** 单一 workspace 根目录更符合部署后的运行时心智，也能避免不同模块对 `workspaces` / `workspace` 命名不一致而产生路径偏差。

### 2. 服务启动时一次性加载 `workspace/skills`

**选择**：Web server 初始化阶段从部署后的 `workspace/skills` 扫描并加载内置 skills；每轮对话继续动态查询用户已启用的数据库 skills，再做合并。

**备选：**
- 每轮对话时都重新扫描 `workspace/skills`
- 只保留数据库 skills，不再支持部署时内置 skills

**理由：** 启动时加载共享 skills 可以降低每轮请求的 I/O 开销，同时保持用户自定义 skills 的动态性。

### 3. 将 `cmd` 产物统一编译到 `workspace/bin`

**选择**：标准构建流程将服务主二进制输出到 `output/bin/`，将运行时会被 skill 或工具调用的 `cmd` 子命令编译到 `output/workspace/bin/`，脚本类资源复制到 `output/workspace/cmd/`。

**备选：**
- 继续把命令二进制输出到 `workspaces/bin`
- 运行时首次使用时再从源码即时编译 `cmd` 目录

**理由：** 运行时依赖的命令应当是可预测的只读部署资产，预编译后放在固定 workspace 路径下，技能文本和工具逻辑才能稳定引用。

### 4. 终端与文件工具统一绑定到 workspace 边界

**选择**：工具运行时注入的默认 cwd 固定为部署 `workspace/` 根目录；所有相对路径解析、bash 执行和文件访问都必须先归一化到该 workspace 边界内，再交给底层 handler。`workspace/skills` 与 `workspace/bin` 作为保留运行时目录，允许 agent 读取，并允许 agent 执行 `workspace/bin` 下的标准命令产物，但不作为普通工作文件的默认写入目标。

**备选：**
- 继续允许工具默认使用进程 cwd
- 允许工具直接访问宿主机用户目录，只在 prompt 中提醒风险

**理由：** 用户已经明确默认 cwd 就是 `workspace/` 根目录，而不是再细分 `workspace/run`。同时，内置 skills 和标准命令产物本身就是部署资源的一部分，因此需要允许 agent 读取、并允许执行标准二进制，但仍要避免把这些目录当作普通可写工作区使用。将边界下沉到 runtime/handler 层，比仅靠 prompt 约束更可靠。

### 5. 提供标准 `build.sh` 作为唯一部署打包入口

**选择**：在仓库根目录提供标准 `build.sh`，完成清理输出目录、编译主服务、编译 `cmd` 二进制、复制 configs/bootstrap/skills 到输出目录，并最终形成统一的部署树；现有 `cmd/build-artifacts` 由 `build.sh` 直接替代，不再保留并行的打包入口。

**备选：**
- 继续保留 `cmd/build-artifacts` 与 `build.sh` 双入口共存
- 只保留 Go 内部构建程序，不提供 shell 打包入口

**理由：** 用户明确要求部署脚本格式，并给出了参考样式，同时已经确认 `build.sh` 应直接替代 `cmd/build-artifacts`。保留唯一打包入口可以避免行为漂移，也更便于流水线与人工打包使用。

## Risks / Trade-offs

- **[Risk] 现有配置和测试仍引用 `workspaces/` 目录名** → 通过集中修改默认值与测试夹具，统一迁移到 `workspace/` 布局
- **[Risk] 工具若仍隐式依赖进程 cwd，可能绕过新边界** → 在 runtime 注入与路径校验层统一重写 cwd，并补充越权测试
- **[Risk] `workspace/skills` 与 `workspace/bin` 被误当作普通写目录，导致运行时资产被污染** → 将这两个目录视为保留运行时目录，仅允许读取及对标准二进制的执行，不作为默认工作文件写入目标
- **[Trade-off] 统一到单一 workspace 根目录后，部署资产和工作目录语义更紧密** → 通过保留 `workspace/skills`、`workspace/bin` 等保留子目录，继续区分运行时资产与普通工作文件

## Migration Plan

1. 调整配置默认值与路径解析逻辑，使 builtin skills、命令二进制和脚本目录都默认落到 `AppHome/workspace/...`
2. 更新 runtime 初始化逻辑，从 `workspace/skills` 加载内置 skills，并把 `workspace/bin` / `workspace/cmd` 注入工具运行环境
3. 更新工具边界校验，使终端与文件类工具默认在 workspace 中工作，并拒绝访问 workspace 外路径
4. 新增标准 `build.sh` 并替代 `cmd/build-artifacts`，生成 `output/bin` 与 `output/workspace/...` 的部署结构
5. 更新测试与部署校验，覆盖新目录布局、skills 加载、命令产物路径、保留目录访问规则和越权拒绝场景
6. 回滚时可恢复到旧的分散目录约定，但需要同步回退配置默认值、默认 cwd 规则与构建脚本入口
