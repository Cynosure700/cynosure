## 1. Runtime & Config Foundation

- [x] 1.1 扩展配置加载逻辑，显式支持应用根目录、builtin skills 目录、cmd 产物目录与 workspace 根目录，并在启动时解析为绝对路径
- [x] 1.2 在 Web server 初始化阶段加载 `go-agent/skills` 下的内置 skill catalog，并将其注入 runtime 依赖
- [x] 1.3 为内置 skill 定义稳定的运行时标识符和描述提取逻辑，避免继续依赖 REPL 路径下的隐式命名方式
- [x] 1.4 建立 `cmd` 目录的部署构建步骤：编译 Go 子命令为二进制产物，并复制 `.py` 等脚本资源到部署目录

## 2. Skill Source Merging

- [x] 2.1 抽象组合 skill loader，合并“共享内置 skill”与“当前用户已启用的数据库 skill”
- [x] 2.2 调整 Web runtime 的 system prompt / `load_skill` 装配逻辑，使每轮对话都能看到合并后的 skill 集合
- [x] 2.3 为用户自定义 skill 增加与内置 skill 的标识符冲突校验，防止 shadow builtin skill

## 3. Web Tool Runtime Integration

- [x] 3.1 为 Web runtime 增加工具适配层，复用 `go-agent/internal/tools` 中的已注册工具定义与 handlers
- [x] 3.2 定义并实现 Web 部署允许的工具白名单，确保只有显式允许的已注册工具会暴露给模型
- [x] 3.3 将工具执行结果和拒绝结果统一回填到 runtime loop，保持 tool-calling 闭环
- [x] 3.4 为 skill / tool 注入共享命令产物根目录，使其可稳定调用部署包中的二进制或脚本文件

## 4. User Workspace Isolation

- [x] 4.1 为每个用户创建并解析独立的 workspace 根目录，默认路径落在配置的 `WorkspaceRoot/<user_id>/`
- [ ] 4.2 让 shell / 文件类工具在 Web 场景下使用用户 workspace 作为默认 cwd，并正确解析相对路径
- [ ] 4.3 区分用户 workspace 与只读部署命令目录，确保共享命令产物不落入用户可写路径
- [ ] 4.4 增加越权路径与非法 cwd 校验，拒绝访问用户 workspace 之外的资源以及未授权的部署目录

## 5. Audit, API, and Verification

- [ ] 5.1 扩展工具审计记录，保存 user、conversation、tool、resolved cwd、命令产物路径、状态与拒绝原因/结果摘要
- [ ] 5.2 调整 skill 管理接口或服务层校验，明确 builtin skill 只读且不能通过用户 API 修改
- [ ] 5.3 补充测试，覆盖内置 skill 加载、内置+用户 skill 合并、工具白名单、workspace 隔离、命令产物构建与审计记录
