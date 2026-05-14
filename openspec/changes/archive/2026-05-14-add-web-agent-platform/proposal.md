## Why

当前 nano_cc / go-agent 仅支持本地 CLI 交互，无法让多个用户通过网页登录后独立使用 agent，也无法将用户自定义 skill 作为平台能力进行持久化和隔离管理。现在需要把 agent 平台化，让用户能在浏览器中登录、对话、创建自己的 skill，并由后端从数据库动态加载 skill 与工具能力，形成一个可实际部署和使用的多租户 Web Agent 系统。

## What Changes

- 新增基于网页的多用户 Agent 平台，包含前端 Web 应用、后端 API 服务、Agent Runtime 服务
- 新增用户注册、登录、会话鉴权能力，确保不同用户的数据、skill、会话彼此隔离
- 新增用户级 skill 管理能力，支持创建、编辑、启用、禁用、删除 skill，并持久化到数据库
- 新增基于浏览器的聊天界面，用户可在网页中与 agent 持续对话并查看工具/skill 执行结果
- 新增后端 Agent Runtime，将现有 agent 的 tool calling / skill loading 机制扩展为按用户上下文从数据库加载 skill
- 新增会话存储、消息持久化、Redis 缓存/任务协调，数据库与 Redis 默认使用本机 IP 和默认端口
- 新增平台级安全边界：用户鉴权、skill 归属校验、工具执行权限控制、会话访问控制

## Capabilities

### New Capabilities
- `user-authentication`: 用户注册、登录、鉴权、会话续期与登出，保证多用户隔离访问
- `skill-management`: 用户独立创建、编辑、发布、启停和删除 skill，skill 元数据与正文持久化到数据库
- `web-chat`: 浏览器聊天界面、会话列表、消息收发与 agent 响应展示
- `agent-runtime`: 面向 Web 请求的 agent 执行服务，支持工具调用、用户 skill 注入、对话上下文管理
- `conversation-storage`: 持久化保存聊天会话、消息、工具调用记录，并通过 Redis 支撑会话状态与短期缓存
- `tool-execution-control`: 平台内工具执行的授权、审计与隔离规则，保证不同用户无法越权访问彼此数据或工作区

### Modified Capabilities
<!-- 无现有 spec 需要修改 -->

## Impact

- 新增前端 Web 应用（登录页、聊天页、skill 管理页）
- 新增后端服务层、鉴权模块、Agent Runtime 接口、数据库与 Redis 集成
- 新增数据模型：users、sessions、skills、conversations、messages、tool_calls 等
- 默认依赖本机数据库与 Redis（local IP + default port），同时需要新增环境变量用于模型配置、JWT 密钥和工作目录
- 现有 CLI agent 可作为底层能力复用，但需要补充面向 Web 平台的接口和多租户隔离机制
