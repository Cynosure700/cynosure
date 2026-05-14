## Context

当前仓库已有两个 agent 方向：TypeScript/Bun 版本的 `nano_cc` 与 Go 版本的 `go-agent`。两者都以 CLI 形态运行，具备工具调用、技能加载、上下文管理等基础能力，但缺少一个可供最终用户通过浏览器访问的产品化入口，也没有用户体系、数据库持久化、Redis 会话缓存和多租户 skill 隔离。

本次变更需要设计一个完整可用的 Web Agent 平台：
- 前端提供登录、聊天、skill 管理界面
- 后端提供用户鉴权、会话管理、skill CRUD、聊天 API
- Agent Runtime 复用现有 agent 核心能力，但必须支持按用户从数据库动态加载 skill
- 数据库与 Redis 均默认连接本机 IP 与默认端口，以便本地开发和单机部署

## Goals / Non-Goals

**Goals:**
- 提供用户注册、登录、登出与鉴权能力
- 提供用户独立的 skill 管理能力，skill 持久化存储在数据库
- 提供网页聊天界面，用户可创建会话并与 agent 对话
- 让 Agent Runtime 在每次对话时按用户上下文动态加载 DB 中已启用的 skill
- 复用现有 tool calling / skill loading 思路，确保 agent 仍可调用工具和技能
- 持久化保存会话、消息、工具调用审计记录，并使用 Redis 管理热点会话状态
- 默认采用 `127.0.0.1:5432`（PostgreSQL）与 `127.0.0.1:6379`（Redis）

**Non-Goals:**
- 不实现企业级权限模型（如组织、团队、RBAC）
- 不实现支付、订阅、邀请码等商业化能力
- 不实现跨机器分布式工具执行调度
- 不实现复杂的 skill 市场或 skill 分享社区
- 不要求 agent 具备“代操作第三方系统”的完整产品闭环，重点是 agent、skill、tool 的平台化承载

## Decisions

### 1. 总体架构：前后端分离 + 单独 Agent Runtime 层

**选择**：采用三层架构。

1. **Web Frontend（Bun + React）**：登录页、聊天页、skill 管理页
2. **API Server（Go）**：鉴权、业务 API、会话控制、数据库访问、Redis 访问
3. **Agent Runtime（Go，复用 go-agent 内核）**：消息编排、tool calling、skill 动态加载、工具执行

**备选**：
- 将前端直接内嵌进 Go 模板服务
- 让 API Server 直接承担所有 runtime 逻辑，不做分层

**理由**：前后端分离便于后续独立演进；将 Agent Runtime 作为清晰边界可避免 Web API 与 LLM/tool 循环强耦合，也更容易复用现有 `go-agent` 模块能力。

### 2. 技术栈：Go 后端 + React 前端 + PostgreSQL + Redis

**选择**：
- 后端：Go（优先扩展 `go-agent/`）
- 前端：React + TypeScript，使用 Bun 作为包管理/开发运行时
- 主数据库：PostgreSQL，默认 `127.0.0.1:5432`
- 缓存与会话态：Redis，默认 `127.0.0.1:6379`

**备选**：
- MySQL 作为主数据库
- 仅用 PostgreSQL，不引入 Redis

**理由**：PostgreSQL 适合存储结构化数据和 JSONB 扩展字段（如 skill metadata、tool traces 摘要），Redis 适合缓存活跃会话、流式响应状态与短期 session/token 索引。Go 已是仓库中的服务端方向，前端使用 React 更适合构建聊天和编辑器界面。

### 3. 鉴权模型：JWT Access Token + Redis Session 索引

**选择**：
- 注册/登录后签发短期 access token
- token 中包含 `user_id`、`session_id`
- Redis 保存 session 索引、撤销态和续期窗口
- 前端通过 HttpOnly Cookie 或 Bearer Token 调用 API（MVP 优先 HttpOnly Cookie）

**备选**：
- 纯服务端 session
- Access/Refresh 双 token 完整体系

**理由**：JWT 能降低 API Server 的无状态压力，Redis 保留登出、撤销和续期能力。对浏览器产品来说，HttpOnly Cookie 更适合降低 token 泄漏风险。

### 4. 用户 skill 模型：数据库持久化 + 运行时按用户加载

**选择**：建立 `skills` 表，核心字段包括：
- `id`
- `user_id`
- `name`
- `slug`
- `description`
- `content`
- `status`（draft/enabled/disabled）
- `created_at` / `updated_at`

Agent Runtime 每次处理用户消息时：
1. 查询该用户所有 `enabled` skill
2. 转换为与现有 skill loader 兼容的内存结构
3. 将 skill 描述注入 system prompt
4. 在模型调用 skill 时按名称返回该 skill 正文

**备选**：
- 将用户 skill 写回文件系统再由原 loader 扫描
- 将 skill 嵌入 conversation 上下文，不走统一加载器

**理由**：DB 原生存储更符合多租户平台；避免把用户数据散落到文件系统；运行时适配器模式可最大程度复用现有 skill loading 逻辑。

### 5. 会话与消息模型：DB 为事实来源，Redis 为热点缓存

**选择**：
- PostgreSQL 保存 `conversations`、`messages`、`tool_calls`
- Redis 保存活跃会话的最近消息窗口、流式响应状态、幂等键
- Agent Runtime 处理时优先从 Redis 取近期上下文，未命中则回落 DB

**备选**：
- 仅保存数据库
- 仅保存 Redis，不落库

**理由**：聊天产品需要可恢复和可审计，数据库不能省；同时 LLM 会话是高频读写场景，引入 Redis 能降低重复拼接上下文的成本。

### 6. 工具执行模型：受控工具白名单 + 用户工作区隔离

**选择**：
- 平台只开放注册在后端 registry 中的工具
- 每次工具调用前执行鉴权、参数校验、用户作用域校验
- 文件类工具仅允许访问用户专属工作区，如 `data/workspaces/<user_id>/`
- Bash 类工具默认关闭高风险命令，沿用现有危险命令黑名单思路
- 所有工具调用写入 `tool_calls` 审计表

**备选**：
- 完全复用 CLI 权限模型
- 不做工作区隔离，由 prompt 约束模型行为

**理由**：Web 平台是多用户环境，不能依赖 prompt 约束。必须在服务端建立明确的强制隔离和审计边界。

### 7. API 形态：REST + SSE

**选择**：
- 普通业务（登录、skill CRUD、会话列表）使用 REST
- 聊天响应使用 SSE，逐步推送 assistant token、tool event、final answer

**备选**：
- 全部 REST 轮询
- 直接使用 WebSocket

**理由**：SSE 对单向流式输出足够简单，适合 LLM 输出和工具事件展示，实现成本低于完整 WebSocket 双向协议。

## Risks / Trade-offs

- **[Risk] 现有 go-agent 更偏 CLI，抽离为可复用 runtime 时可能牵涉较多重构** → 先增加 runtime adapter 层，保持 CLI 入口与 Web 入口共用核心 loop
- **[Risk] 用户自定义 skill 可能注入恶意提示，诱导工具越权** → skill 仅作为内容来源，真正权限由服务端 registry、ACL 和工作区隔离决定
- **[Risk] Redis 与 DB 双写可能造成状态短暂不一致** → DB 作为事实来源，Redis 仅缓存；写入顺序固定为先 DB 后 Redis 刷新
- **[Risk] 聊天消息和工具调用增长较快** → 设计归档/分页能力，前端会话详情分页读取，Redis 仅缓存最近窗口
- **[Trade-off] 采用 PostgreSQL + Redis 增加部署组件数量** → 换取可恢复消息、低延迟上下文访问和会话控制能力
- **[Trade-off] SSE 简单但不适合双向复杂实时协议** → 对本次以单向输出为主的聊天产品足够，后续如需协同编辑再考虑 WebSocket

## Migration Plan

1. 在仓库中新增 Web 前端与 Web API/Runtime 模块
2. 引入数据库迁移脚本，创建 users / auth_sessions / skills / conversations / messages / tool_calls 表
3. 配置本地 PostgreSQL 与 Redis 默认地址
4. 先打通用户注册登录、skill CRUD、聊天会话创建
5. 再接入 Agent Runtime，使其能从 DB 加载用户 skill 并通过 SSE 返回响应
6. 完成前端页面联调后再补充审计、缓存和错误恢复
7. 若回滚，关闭 Web 路由并回退数据库迁移；CLI agent 保持可独立运行

## Open Questions

- 用户 skill 是否允许上传附件或只支持纯文本/Markdown？MVP 先限制为纯文本
- 是否在首版开放 Bash 工具给网页用户？MVP 建议默认关闭高风险工具，只开放安全白名单工具
- 前端是否需要展示逐 token 流式输出与工具调用时间线？设计中预留 SSE 事件类型，MVP 可先实现基础版本
