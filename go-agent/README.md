# nano_cc (Go)

基于 Go 实现的编码智能体，当前同时支持两种使用方式：

1. **CLI 模式**：本地 REPL 交互式编码助手
2. **Web 平台模式**：多用户登录、Skill 管理、网页聊天、MySQL/Redis 持久化

它通过 OpenAI 兼容 API 驱动 agent loop，支持工具调用、Skill 加载、上下文管理，以及在 Web 模式下按用户从数据库动态加载 Skill。

---

## 功能概览

### CLI 模式

- 交互式 REPL
- 工具调用循环（tool calling）
- 文件读写与编辑
- Shell 命令执行
- Todo 管理
- 子智能体委派
- 本地 Skill 文件加载
- 上下文压缩

### Web 平台模式

- 用户注册 / 登录 / 登出
- 基于 Cookie + JWT 的鉴权
- 用户独立 Skill 创建、编辑、启停、删除
- 会话与消息持久化
- Agent 从 MySQL 中动态加载用户已启用 Skill
- SSE 流式返回 assistant 输出与 tool event
- Redis 会话缓存
- 多用户数据隔离
- 工具白名单与用户工作区隔离

---

## 目录结构

```text
go-agent/
├── main.go                     # CLI 入口
├── cmd/
│   └── web/
│       └── main.go             # Web 后端入口
├── config.json                 # CLI / LLM 配置文件（可选）
├── internal/
│   ├── agent/                  # CLI agent loop 与 REPL
│   ├── config/                 # LLM / Web 配置
│   ├── logger/                 # 日志
│   ├── safety/                 # 路径安全
│   ├── sessions/               # memory / skill / subagent / compact
│   ├── tools/                  # CLI 工具系统
│   └── web/
│       ├── app/                # HTTP Server 与路由
│       ├── auth/               # 注册 / 登录 / Session / JWT
│       ├── runtime/            # Web agent runtime / tool registry / SSE
│       └── storage/            # MySQL / Redis / migrations / repository
└── skills/                     # 本地文件 Skill（CLI 模式）
```

前端页面位于仓库根目录下的 `web/`。

---

## 环境要求

- Go 1.21+
- Node.js / npm（用于前端 Web 页面）
- OpenAI 兼容模型服务
- MySQL（Web 模式）
- Redis（Web 模式）

---

## 一、CLI 模式

### 启动方式

```bash
cd go-agent
go run .
```

### LLM 配置

CLI 模式会优先读取环境变量；如果未设置，则回退到 `config.json`。

#### 方式 1：使用 `config.json`

```json
{
  "base_url": "https://api.deepseek.com",
  "api_key": "your-api-key",
  "model_id": "deepseek-chat"
}
```

#### 方式 2：使用环境变量

```bash
export OPENAI_BASE_URL=https://api.deepseek.com
export OPENAI_API_KEY=your-api-key
export MODEL_ID=deepseek-chat
```

### CLI 示例

```text
You: 用 Go 写一个 HTTP 服务，监听 8080 端口，返回 hello world
Assistant: 我来创建这个 HTTP 服务。
```

退出：

```text
You: exit
```

---

## 二、Web 平台模式

Web 模式由两部分组成：

1. **Go 后端**：`go-agent/cmd/web/main.go`
2. **React 前端**：仓库根目录 `web/`

### Web 模式默认基础设施

当前默认值如下：

- MySQL: `1.12.217.28:3306`
- Redis: `1.12.217.28:6379`

如果不传环境变量，后端会默认按上述地址构造连接。

### Web 后端配置

后端同样需要 LLM 配置，且支持以下环境变量：

```bash
OPENAI_BASE_URL=https://api.deepseek.com
OPENAI_API_KEY=your-api-key
MODEL_ID=deepseek-chat

SERVER_ADDR=:8080
ALLOWED_ORIGIN=http://localhost:5173

MYSQL_HOST=1.12.217.28
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=213140
MYSQL_DATABASE=vibe_coding

REDIS_ADDR=1.12.217.28:6379
REDIS_PASSWORD=213140
REDIS_DB=0

JWT_SECRET=nano-cc-local-secret
SESSION_COOKIE_NAME=nano_cc_session
SESSION_TTL_MINUTES=10080

WORKSPACE_ROOT=data/workspaces
```

### 启动 Web 后端

```bash
cd go-agent

export OPENAI_BASE_URL=https://api.deepseek.com
export OPENAI_API_KEY=your-api-key
export MODEL_ID=deepseek-chat

export MYSQL_HOST=1.12.217.28
export MYSQL_PORT=3306
export MYSQL_USER=root
export MYSQL_PASSWORD=213140
export MYSQL_DATABASE=vibe_coding

export REDIS_ADDR=1.12.217.28:6379
export REDIS_PASSWORD=213140

go run ./cmd/web
```

### 启动前端

```bash
cd web
npm install
npm run dev
```

默认前端地址：

- `http://localhost:5173`

默认后端地址：

- `http://localhost:8080`

---

## Web API 简述

### 鉴权相关

- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/logout`
- `GET /api/me`

### Skill 管理

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

## Web 模式的核心行为

### 1. 用户 Skill 动态加载

每次用户在 Web 页面发消息时，runtime 会：

1. 从 MySQL 读取该用户所有 `enabled` 状态的 Skill
2. 构造成运行时 Skill Loader
3. 将 Skill 描述注入 system prompt
4. 在模型请求 `load_skill` 时返回对应 Skill 内容

### 2. 多用户数据隔离

- 用户只能访问自己的 Skill
- 用户只能访问自己的 Conversation / Message
- 文件类工具只允许操作自己的工作区

### 3. 工具执行控制

Web 模式下默认只暴露受控工具集合，不开放任意 shell。

当前 Web runtime 注册的工具包括：

- `read_file`
- `write_file`
- `edit_file`
- `load_skill`（仅当当前用户存在 enabled skill 时可用）

---

## CLI 工具列表

CLI 模式下支持的工具：

| 工具 | 说明 |
|------|------|
| `bash` | 执行 shell 命令（带危险命令拦截） |
| `read_file` | 读取文件 |
| `write_file` | 写入文件 |
| `edit_file` | 精确替换文本 |
| `todo` | 管理任务列表 |
| `task` | 委派子智能体 |
| `load_skill` | 加载 Skill |
| `compact` | 手动触发上下文压缩 |

---

## Skill 文件格式（CLI 模式）

CLI 模式会从 `skills/` 目录扫描 Markdown Skill 文件。

格式示例：

```markdown
---
description: Git workflow helpers
tags: git, version-control
---

# Git Workflow

...正文...
```

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

当前实现已完成以下验证：

- Go 后端编译通过
- 前端 typecheck / build 通过
- 用户注册 / 登录 / 登出流程通过
- 未登录访问受保护 API 返回 401
- 多用户 Skill / Conversation 隔离通过
- enabled Skill 可在运行时动态加载
- SSE 可返回 tool event + assistant 内容
- Redis cache miss 时可从 MySQL 回退
- 禁用 Skill / 请求未注册工具等异常路径行为符合预期

---

## 注意事项

1. Web 模式依赖可用的 MySQL、Redis、LLM 服务
2. 如果前端跨域访问失败，请检查 `ALLOWED_ORIGIN`
3. 如果登录后接口仍返回 401，请检查 Cookie 是否被浏览器拦截
4. Web 模式下文件工具会限制在用户工作区目录内

---

## 后续建议

如果后续要继续增强这个平台，优先建议：

1. 增加密码重置与用户资料管理
2. 增加 Skill 版本管理
3. 增加会话分页与消息分页
4. 增加更细粒度的工具权限控制
5. 增加部署脚本 / Docker Compose
