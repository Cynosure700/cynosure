## 1. Runtime root resolution

- [x] 1.1 在 `internal/config` 中收敛 runtime workspace 解析逻辑，支持显式配置覆盖、`output/workspace` 优先和源码 `workspace` 回退
- [x] 1.2 让 builtin skills、command bin、command scripts 的默认目录统一从解析后的 runtime workspace 派生
- [x] 1.3 更新 `config.json` 与 `output/config.json` 的默认约定，避免继续依赖进程 cwd 推断路径

## 2. Runtime and skill loading

- [x] 2.1 更新 `internal/web/app/server.go`，在启动阶段使用解析后的 runtime workspace 加载 builtin skills
- [x] 2.2 更新 `internal/web/runtime` 的环境注入逻辑，确保会话执行、tool runtime env 和审计字段共享同一 runtime root
- [x] 2.3 校正 skill 与命令路径注入，确保部署态命中 `output/workspace/*`，本地态命中源码 `workspace/*`

## 3. Tool execution boundary

- [x] 3.1 更新 bash/file 工具默认 cwd 选择逻辑，使其始终绑定解析后的 runtime workspace
- [x] 3.2 保持越权校验基于活动 runtime workspace 生效，拒绝访问该根目录之外的绝对路径
- [ ] 3.3 更新命令路径解析与审计输出，确保能区分部署态与本地调试态的命令资产来源

## 4. Verification

- [ ] 4.1 为配置解析补充测试，覆盖显式覆盖、部署优先和本地回退三种分支
- [ ] 4.2 为 runtime/skill 加载补充测试，验证 builtin skills 与 command 目录在部署态、本地态下都能正确注入
- [ ] 4.3 为工具执行补充测试，验证默认 cwd、绝对路径拒绝和命令资产解析行为
