# nano_cc (Go)

`nano_cc` 现已收敛为一个**浏览器优先的通用聊天机器人后端**：

- 默认入口是 Web 服务，而不是 CLI REPL
- 默认角色是通用聊天助手，而不是编码助手
- 支持登录、会话、流式聊天、个人能力（Skill）管理
- 网页端**不执行 shell 命令，不访问用户本地目录，不读写本地文件**

它通过 OpenAI 兼容 API 驱动聊天运行时，在需要时可从数据库按用户加载已启用的能力内容，并通过 SSE 把响应流式返回给前端页面。

---

## 当前定位

这是一个面向浏览器聊天产品的 Go 后端，核心目标是提供类似 ChatGPT 的对话体验，而不是继续围绕本地 CLI 编码代理构建产品。

当前正式使用方式：

1. **Go Web 服务**：提供鉴权、会话、聊天、能力管理 API
2. **React 前端**：提供 conversation-first 的网页聊天界面

仓库里仍保留了一部分旧的 CLI / 本地工具相关代码，主要用于历史兼容或内部复用，但它们**不再是正式产品入口**。

---

## 功能概览

### 浏览器聊天能力

- 用户注册 / 登录 / 登出
- 基于 Cookie + JWT 的鉴权
- 多会话聊天
- SSE 流式返回 assistant 输出
- conversation-first 的网页聊天体验
- 通用问答、写作、规划、分析、代码协助等对话能力

### 用户能力（Skill）管理

- 创建、编辑、启用、禁用、删除个人能力
- 按用户隔离能力数据
- 运行时动态加载当前用户已启用能力

### 平台能力与边界

- 会话、消息、工具调用记录持久化
- Redis 缓存活跃会话上下文
- 多用户数据隔离
- 浏览器端仅暴露安全且与网页聊天模型兼容的能力
- 默认浏览器模式下当前仅保留 `load_skill`
- 明确拒绝 shell、本地目录、本地文件相关请求，并返回清晰解释

---

## 目录结构

```text
go-agent/
├── main.go                     # 默认入口：启动 Web 服务
├── cmd/
│   └── web/
│       └── main.go             # Web 服务备用入口（等价启动方式）
├── config.json                 # LLM 配置文件（可选）
├── internal/
│   ├── agent/                  # 旧的 agent/REPL 相关实现（非正式产品入口）
│   ├── assistant/              # 通用 assistant system prompt 构造
│   ├── config/                 # LLM / Web 配置
│   ├── logger/                 # 日志
│   ├── safety/                 # 安全辅助
│   ├── sessions/               # memory / skill / subagent / compact
│   ├── tools/                  # 旧工具系统（主要供非 Web 路径复用）
│   └── web/
│       ├── app/                # HTTP Server 与路由
│       ├── auth/               # 注册 / 登录 / Session / JWT
│       ├── runtime/            # Web 聊天 runtime / tool registry / SSE
│       └── storage/            # MySQL / Redis / migrations / repository
└── skills/                     # 本地能力文件（历史兼容用途）
```

前端页面位于仓库根目录下的 `web/`。

---

## 环境要求

- Go 1.21+
- Node.js / npm（用于前端开发与构建）
- OpenAI 兼容模型服务
- MySQL
- Redis

---

## 配置说明

后端会优先读取环境变量；如果 LLM 相关环境变量未设置，则回退到 `config.json`。

### `config.json` 示例

```json
{
  "base_url": "https://api.deepseek.com",
  "api_key": "your-api-key",
  "model_id": "deepseek-chat"
}
```

### 常用环境变量

```bash
OPENAI_BASE_URL=https://api.deepseek.com
OPENAI_API_KEY=your-api-key
MODEL_ID=deepseek-chat

SERVER_ADDR=:8080
ALLOWED_ORIGIN=http://localhost:5173

MYSQL_HOST=127.0.0.1
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=your-password
MYSQL_DATABASE=vibe_coding

REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=
REDIS_DB=0

JWT_SECRET=replace-with-your-own-secret
SESSION_COOKIE_NAME=nano_cc_session
SESSION_TTL_MINUTES=10080
```

说明：

- `WORKSPACE_ROOT` 仍然在配置结构中保留，但**当前浏览器聊天主流程不依赖本地工作区能力**
- 如果未提供 MySQL / Redis 环境变量，程序会使用代码中的默认值构造连接

---

## 启动方式

### 1）启动 Go 后端（默认入口）

```bash
cd go-agent
go run .
```

也可以使用备用入口：

```bash
cd go-agent
go run ./cmd/web
```

默认后端地址：

- `http://localhost:8080`

### 2）启动前端

```bash
cd web
npm install
npm run dev
```

默认前端地址：

- `http://localhost:5173`

---

## Web API 简述

### 鉴权相关

- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/logout`
- `GET /api/me`

### 能力管理

- `GET /api/skills`
- `POST /api/skills`
- `GET /api/skills/:id`
- `PUT /api/skills/:id`
- `PATCH /api/skills/:id`
- `DELETE /api/skills/:id`

### 会话与聊天

- `GET /api/conversations`
- `POST /api/conversations`
- `GET /api/conversations/:id`
- `POST /api/conversations/:id/stream`

### 健康检查

- `GET /api/health`

---

## 核心行为

### 1. 默认入口是 Web 服务

`go-agent/main.go` 现在默认启动 Web 服务，而不是进入本地 REPL。

### 2. 默认角色是通用聊天助手

系统提示词已统一为通用 assistant 基线：

- 支持通用问答、分析、规划、写作、编码协助
- 优先直接回答，而不是先假设要调用工具
- 不默认假设 shell、本地目录、本地文件访问能力

### 3. 用户能力动态加载

每次用户在网页中发送消息时，runtime 会：

1. 从数据库读取该用户所有 `enabled` 状态的能力
2. 构造成运行时 Skill Loader
3. 将能力描述注入 system prompt
4. 在模型请求 `load_skill` 时返回对应能力正文

### 4. 多用户数据隔离

- 用户只能访问自己的 Skill
- 用户只能访问自己的 Conversation / Message
- 工具调用记录按用户与会话隔离存储

### 5. 浏览器能力边界

当前浏览器聊天模式下：

- **不会暴露 shell 命令执行能力**
- **不会暴露本地目录浏览能力**
- **不会暴露本地文件读写能力**
- 当前 Web runtime 默认只暴露：
  - `load_skill`

如果用户请求“执行 shell / 读取本地文件 / 浏览用户目录”，系统会直接返回清晰说明，并继续提供替代帮助，例如：

- 解释命令含义
- 生成可手动执行的命令或脚本
- 基于用户贴出的报错 / 文件内容继续分析

---

## 开发与验证

### Go 侧

```bash
cd go-agent

gofmt -w ./...
go test ./...
```

### 前端

```bash
cd web

npm install
npm run typecheck
npm run build
```

---

## 已验证内容

当前实现已经完成以下验证：

- Go 测试通过
- 前端 `typecheck` / `build` 通过
- 默认入口 `go run .` 可直接启动 Web 服务
- 健康检查接口可用
- 用户注册 / 建会话 / SSE 发消息主流程已真实走通
- 通用问答、写作、规划类请求可正常返回，不依赖 CLI 心智
- 浏览器模式下健康接口只暴露 `load_skill`
- shell / 本地文件 / 用户目录请求会返回明确的能力边界说明

---

## 注意事项

1. 运行前请确保 MySQL、Redis、LLM 服务可用
2. 如果前端跨域访问失败，请检查 `ALLOWED_ORIGIN`
3. 如果登录后接口仍返回 401，请检查 Cookie 是否被浏览器拦截
4. 浏览器聊天模式不是本地终端代理，不支持替用户执行本地命令或访问本地目录

---

## 后续可继续增强的方向

1. 增加密码重置与用户资料管理
2. 增加能力版本管理
3. 增加会话分页与消息分页
4. 增加更细粒度的能力面板与权限控制
5. 增加部署脚本 / Docker Compose
