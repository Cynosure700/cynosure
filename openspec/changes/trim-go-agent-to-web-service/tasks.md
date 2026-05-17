## 1. 收敛 Web 服务启动链路

- [x] 1.1 梳理 `main.go`、`cmd/web`、`internal/web/app` 的调用关系并删除重复启动入口
- [x] 1.2 从 Web 服务依赖链中移除 CLI/REPL 相关模块与无入口历史 session 能力
- [x] 1.3 调整包引用与测试，使 Web 服务仅依赖配置、日志、存储、鉴权、聊天 runtime 及其所需 runtime workspace 资产解析链路

## 2. 简化运行时与部署布局

- [x] 2.1 调整配置解析与启动校验逻辑，明确 `output/workspace` / `workspace` 是 Web runtime 资产目录，并保留 `workspace/bin`、`workspace/cmd` 作为 agent 命令执行所需路径
- [x] 2.2 清理这些 runtime 资产目录与 CLI/REPL 历史结构之间的错误耦合，使服务仅校验 Web 服务与 agent runtime 实际需要的依赖
- [ ] 2.3 处理数据库迁移和其他服务必需资源的加载方式，确保本地运行与部署运行一致

## 3. 清理构建与发布结构

- [ ] 3.1 精简或替换 `build.sh`、`cmd/build-artifacts`、`internal/deploy` 等构建逻辑，只保留 Web 服务正式构建路径，并继续正确产出 `workspace/bin`、`workspace/cmd` 等 runtime 资产
- [ ] 3.2 清理与 Web 服务启动无关的目录约定、产物说明和历史兼容脚本引用，但保留 Web agent 必需的 workspace 资产目录约定
- [ ] 3.3 验证精简后的构建产物能够以单一正式服务发布包方式启动，且该发布包包含 Web agent 所需的 runtime workspace 资产

## 4. 更新文档与验证

- [ ] 4.1 重写 `README.md` 的项目结构说明，只描述保留后的 Web 服务目录与职责
- [ ] 4.2 重写 `README.md` 的本地部署步骤，提供最短可运行路径
- [ ] 4.3 重写 `README.md` 的云服务器部署步骤，提供构建、上传、启动与运维建议
- [ ] 4.4 运行 Go 测试与最小启动验证，确认核心 Web 服务流程在重构后仍可用，且 runtime workspace 资产解析与命令执行链路未被破坏
