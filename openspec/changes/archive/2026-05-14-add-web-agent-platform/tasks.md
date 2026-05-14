## 1. Web Platform Scaffold

- [x] 1.1 创建前端 Web 应用目录与基础脚手架（React + TypeScript + Bun），配置路由与基础页面骨架
- [x] 1.2 在 Go 侧创建 Web API/Runtime 模块结构，拆分 auth、skills、chat、runtime、storage 子包
- [x] 1.3 增加统一配置读取，支持 PostgreSQL 默认 `127.0.0.1:5432`、Redis 默认 `127.0.0.1:6379`、JWT 密钥、模型配置与用户工作区根目录

## 2. Database & Redis Foundation

- [x] 2.1 实现数据库连接初始化与健康检查
- [x] 2.2 实现 Redis 客户端初始化与健康检查
- [x] 2.3 创建数据库迁移，新增 `users`、`auth_sessions`、`skills`、`conversations`、`messages`、`tool_calls` 表
- [x] 2.4 为 skill、conversation、message、tool_call 建立基础 repository 层

## 3. User Authentication

- [x] 3.1 实现用户注册接口，校验唯一邮箱/用户名并安全存储密码哈希
- [x] 3.2 实现用户登录接口，签发带 `user_id` 和 `session_id` 的认证会话
- [x] 3.3 实现鉴权中间件，保护 skill、conversation、chat runtime 相关 API
- [x] 3.4 实现登出接口，使当前认证会话失效

## 4. Skill Management

- [x] 4.1 实现用户 skill 的创建、列表、详情、更新、删除接口
- [x] 4.2 实现 skill 启用/禁用状态切换接口
- [x] 4.3 在服务层加入 skill ownership 校验，防止跨用户访问或修改
- [x] 4.4 为 skill 数据模型补充 slug、description、content、status 等字段映射与校验

## 5. Conversation Storage & APIs

- [x] 5.1 实现会话创建与会话列表接口，仅返回当前用户数据
- [x] 5.2 实现消息持久化逻辑，保存用户消息与 assistant 回复
- [x] 5.3 实现工具调用审计持久化逻辑，记录工具名、状态、摘要和所属用户/会话
- [x] 5.4 实现 Redis 会话缓存，用于保存活跃会话最近上下文并支持缓存失效回退数据库

## 6. Agent Runtime Integration

- [x] 6.1 将现有 `go-agent` skill loading 抽象为可从数据库记录构造的 runtime loader
- [x] 6.2 实现面向 Web 请求的 agent turn 执行入口，按 conversation 组装上下文消息
- [x] 6.3 在每次对话前按用户加载所有 enabled skills，并注入到 agent runtime
- [x] 6.4 打通 runtime 中的工具调用闭环：模型请求工具、后端执行工具、结果回填模型上下文

## 7. Tool Execution Control

- [x] 7.1 建立 Web 平台工具 registry，仅暴露允许的工具集合
- [x] 7.2 为文件/命令类工具增加用户工作区隔离和参数校验
- [x] 7.3 为工具执行增加授权校验与拒绝日志
- [x] 7.4 将每次工具执行结果写入 `tool_calls` 审计记录

## 8. Web Chat Experience

- [x] 8.1 实现登录页与注册页，完成认证态保存与路由守卫
- [x] 8.2 实现聊天主页面，支持会话列表、消息列表、输入框和发送动作
- [x] 8.3 实现基于 SSE 的 assistant 响应流展示，并在界面中展示工具事件
- [x] 8.4 实现 skill 管理页面，支持 skill 的创建、编辑、启停和删除

## 9. End-to-End Verification

- [x] 9.1 验证新用户可注册、登录、登出，且未登录无法访问受保护 API
- [x] 9.2 验证两个不同用户的 skill、conversation、message 数据完全隔离
- [x] 9.3 验证用户创建并启用 skill 后，agent 能从数据库加载该 skill 并参与回答
- [x] 9.4 验证 agent 在网页聊天中可以调用已注册工具，并将工具事件/结果返回前端
- [x] 9.5 验证 Redis 缓存未命中、skill 禁用、越权访问、非法工具调用等异常路径
