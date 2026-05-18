## 1. Web 与运行时模块拆分

- [x] 1.1 拆分 `go-agent/internal/web/app/server.go`，将路由注册、鉴权 handler、Skill handler、Conversation handler 与 HTTP 辅助逻辑整理到职责清晰的同包文件中。
- [x] 1.2 拆分 `go-agent/internal/web/runtime/runtime.go`，分离会话编排、tool 执行、prompt 构造、SSE 输出等逻辑，同时保持现有运行时行为不变。
- [x] 1.3 梳理 `go-agent/internal/web/runtime/tools.go` 与 `go-agent/internal/tools/registry.go` 的边界，把 Web runtime 专属逻辑与通用工具定义 / handler 拆开。

## 2. 配置、存储与 Skill 加载整理

- [x] 2.1 拆分 `go-agent/internal/config/config.go`，按 LLM 配置、Web 配置、runtime 路径、layout 校验和环境变量辅助整理文件结构。
- [x] 2.2 拆分 `go-agent/internal/web/storage/store.go`，按 store 初始化、migrations、users/sessions/skills/conversations/messages/toolcalls/cache 组织仓储代码。
- [x] 2.3 拆分 `go-agent/internal/sessions/skill.go`，分离 Skill loader、frontmatter 解析、merge 与渲染逻辑，保持现有调用方式兼容。

## 3. 回归验证与文档同步

- [x] 3.1 运行并修复 `go-agent` 相关测试，确保重构后 Web API、聊天 runtime、工具执行、配置解析和 Skill 加载行为保持兼容。
- [ ] 3.2 运行构建或打包流程，确认部署产物结构与现有方式一致。
- [ ] 3.3 更新 `go-agent/README.md` 的目录结构、模块职责与开发验证说明，使其与重构后的实现一致。
