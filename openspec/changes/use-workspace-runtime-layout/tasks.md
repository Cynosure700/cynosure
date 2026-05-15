## 1. Workspace path model

- [x] 1.1 调整配置默认值与路径解析逻辑，统一将 builtin skills、命令二进制和脚本目录收敛到 `AppHome/workspace/...`
- [x] 1.2 更新应用启动阶段的目录校验，确保 `workspace/skills`、`workspace/bin` 与相关 workspace 子目录在服务启动前可用

## 2. Runtime skill and command loading

- [x] 2.1 更新内置 skill 加载逻辑，使 Web runtime 在初始化时从 `workspace/skills` 扫描共享 skills
- [ ] 2.2 更新 runtime 环境注入逻辑，使工具与 skill 能稳定解析 `workspace/bin` 和 `workspace/cmd` 路径
- [ ] 2.3 为用户自定义 skill 与 builtin skill 的标识符冲突增加校验，避免覆盖 `workspace/skills` 中的部署技能

## 3. Tool execution boundary

- [ ] 3.1 调整 Web runtime 默认 cwd 与相对路径解析逻辑，使终端和文件工具默认在部署 workspace 内运行
- [ ] 3.2 增加 workspace 越权校验，拒绝访问或执行落在 deployment workspace 之外的宿主机路径
- [ ] 3.3 更新工具审计或错误回填，确保 workspace 路径拒绝与命令路径解析结果可观测

## 4. Build packaging and verification

- [ ] 4.1 新增或更新标准 `build.sh`，生成包含 `output/bin` 与 `output/workspace/...` 的部署产物结构
- [ ] 4.2 调整 `cmd/build-artifacts` 或共享构建逻辑，使其与 `build.sh` 使用一致的 workspace 目录规则
- [ ] 4.3 补充测试，覆盖 workspace 目录默认值、builtin skills 加载路径、命令产物输出路径与越权拒绝场景
