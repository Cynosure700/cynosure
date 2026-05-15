# nano_cc (Go)

`nano_cc` 现已收敛为一个**可部署的浏览器优先 Agent 后端**：

- 默认入口是 Web 服务，而不是 CLI REPL
- 默认角色是通用聊天助手，但运行时支持 Skill / Tool 调用
- 支持登录、会话、流式聊天、内置/用户 Skill 管理
- 支持在**服务端隔离 workspace** 中运行授权工具，而不是访问用户本地机器
- 支持在部署阶段编译 `cmd/` 下的命令产物，并将脚本资源发布到固定目录

它通过 OpenAI 兼容 API 驱动聊天运行时，合并加载仓库内置 Skill 与数据库中的用户 Skill，并在需要时把工具执行限制在服务端统一 workspace 与只读部署产物目录内，再通过 SSE 把响应流式返回给前端页面。

---

## 当前定位

这是一个面向浏览器聊天产品的 Go 后端，核心目标是提供类似 ChatGPT 的对话体验，同时把 Agent 所需的 Skill、Tool、部署产物、用户工作区统一收敛到服务端运行时中，而不是继续围绕本地 CLI 编码代理构建产品。

当前正式使用方式：

1. **Go Web 服务**：提供鉴权、会话、聊天、能力管理 API
2. **React 前端**：提供 conversation-first 的网页聊天界面

仓库里仍保留了一部分旧的 CLI / 本地工具相关代码，主要用于历史兼容或被 Web runtime 复用，但它们**不再是正式产品入口**。

---

## 功能概览

### 浏览器聊天能力

- 用户注册 / 登录 / 登出
- 基于 Cookie + JWT 的鉴权
- 多会话聊天
- SSE 流式返回 assistant 输出
- conversation-first 的网页聊天体验
- 通用问答、写作、规划、分析、代码协助等对话能力

### Skill 管理

- 创建、编辑、启用、禁用、删除个人 Skill
- 启动时加载 `APP_HOME/workspaces/skills` 下的内置 Skill catalog
- 运行时合并“共享内置 Skill + 当前用户已启用 Skill”
- 内置 Skill 通过 API 以只读条目暴露，不能被用户修改或删除
- 按用户隔离 Skill 数据

### 平台能力与边界

- 会话、消息、工具调用记录持久化
- Redis 缓存活跃会话上下文
- 多用户数据隔离
- 浏览器端仅暴露显式允许的已注册工具
- 默认 Web runtime 当前仅暴露 `load_skill`，可通过 `WEB_ALLOWED_TOOLS` 扩展
- 工具默认运行在服务端配置的统一 `WORKSPACE_ROOT` 下，并拒绝越权路径
- 部署命令产物位于 `APP_HOME/workspaces/bin` 与 `APP_HOME/workspaces/cmd` 只读目录，不落入用户可写 workspace
- 工具调用会记录审计摘要，包括 cwd、命令产物路径、结果摘要或拒绝原因
- 浏览器端不会访问**用户本地机器**的 shell、目录或文件；如果启用工具，访问的也是服务端隔离 workspace

---

## 目录结构

```text
go-agent/
├── main.go                     # 默认入口：启动 Web 服务
├── cmd/
│   ├── build-artifacts/
│   │   └── main.go             # 部署阶段构建 cmd 产物与脚本资源
│   └── web/
│       └── main.go             # Web 服务备用入口（等价启动方式）
├── config.json                 # LLM 配置文件（可选）
├── internal/
│   ├── agent/                  # 旧的 agent/REPL 相关实现（非正式产品入口）
│   ├── assistant/              # 通用 assistant system prompt 构造
│   ├── config/                 # LLM / Web 配置
│   ├── deploy/                 # 部署产物构建逻辑
│   ├── logger/                 # 日志
│   ├── safety/                 # 路径与访问安全辅助
│   ├── sessions/               # memory / skill / subagent / compact
│   ├── tools/                  # 共享工具注册与执行逻辑（Web runtime 复用）
│   └── web/
│       ├── app/                # HTTP Server 与路由
│       ├── auth/               # 注册 / 登录 / Session / JWT
│       ├── runtime/            # Web 聊天 runtime / tool registry / SSE
│       └── storage/            # MySQL / Redis / migrations / repository
├── logs/                       # 服务日志目录
├── workspace/                  # 服务端统一共享 workspace
└── workspaces/
    ├── bin/                    # 部署阶段编译后的命令产物
    ├── cmd/                    # 部署阶段发布的脚本资源
    └── skills/                 # 平台内置 Skill catalog（所有用户共享）
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

后端会优先读取环境变量；如果 LLM 相关环境变量未设置，则回退到 `config.json`。Web 配置会在启动时把 `APP_HOME`、builtin skills 目录、命令产物目录和 workspace 根目录解析为绝对路径；日志文件会写入 `APP_HOME/logs/`。

### `config.json` 示例

```json
{
  "base_url": "https://api.deepseek.com",
  "api_key": "your-api-key",
  "model_id": "deepseek-chat",
  "builtin_skills_dir": "workspaces/skills",
  "command_bin_dir": "workspaces/bin",
  "command_script_dir": "workspaces/cmd",
  "workspace_root": "workspace",
  "web_allowed_tools": "load_skill"
}
```

### 常用环境变量

```bash
OPENAI_BASE_URL=https://api.deepseek.com
OPENAI_API_KEY=your-api-key
MODEL_ID=deepseek-chat

SERVER_ADDR=:8080
ALLOWED_ORIGIN=http://localhost:5173
APP_HOME=/path/to/go-agent

BUILTIN_SKILLS_DIR=workspaces/skills
COMMAND_BIN_DIR=workspaces/bin
COMMAND_SCRIPT_DIR=workspaces/cmd
WORKSPACE_ROOT=workspace
WEB_ALLOWED_TOOLS=load_skill

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

- `APP_HOME` 默认是当前工作目录，其他路径配置会相对它解析
- `BUILTIN_SKILLS_DIR` 默认是 `APP_HOME/workspaces/skills`
- `COMMAND_BIN_DIR` 默认是 `APP_HOME/workspaces/bin`
- `COMMAND_SCRIPT_DIR` 默认是 `APP_HOME/workspaces/cmd`
- `WORKSPACE_ROOT` 默认是 `APP_HOME/workspace`，工具会以它作为服务端统一工作目录根路径
- 日志默认写入 `APP_HOME/logs/session_<timestamp>.log`
- `WEB_ALLOWED_TOOLS` 默认是 `load_skill`；只有显式允许且已注册的工具才会暴露给模型
- 即使启用 `bash` / `read_file` / `write_file` / `edit_file`，它们访问的也是服务端 workspace，而不是用户本地电脑
- 如果未提供 MySQL / Redis 环境变量，程序会使用代码中的默认值构造连接

---

## 启动方式

### 1）构建部署命令产物（推荐在部署阶段执行）

```bash
cd go-agent
go run ./cmd/build-artifacts --app-home .
```

这一步会：

- 发现 `cmd/*/main.go` 并编译到 `workspaces/bin/`
- 复制 `.py` / `.sh` 等脚本资源到 `workspaces/cmd/`

### 2）启动 Go 后端（默认入口）

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

### 3）启动前端

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

### Skill 管理

- `GET /api/skills`  # 返回 builtin + custom skills；builtin 条目含 `source=builtin`、`readonly=true`
- `POST /api/skills`
- `GET /api/skills/:id`
- `PUT /api/skills/:id`
- `PATCH /api/skills/:id`
- `DELETE /api/skills/:id`

说明：

- builtin skill 的 ID 形式为 `builtin:<skill-name>`
- builtin skill 允许读取，但不允许通过用户 API 修改或删除

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
- 不默认假设用户本地 shell、本地目录、本地文件访问能力

### 3. Skill 运行时合并加载

每次用户在网页中发送消息时，runtime 会：

1. 加载共享 builtin Skill catalog
2. 从数据库读取该用户所有 `enabled` 状态的 Skill
3. 合并成当前对话的运行时 Skill Loader
4. 将合并后的能力描述注入 system prompt
5. 在模型请求 `load_skill` 时返回对应能力正文，并附带部署路径提示

### 4. 多用户数据与 workspace 隔离

- 用户只能访问自己的 Skill
- 用户只能访问自己的 Conversation / Message
- 工具调用记录按用户与会话隔离存储
- 如果启用工具，所有用户都共享服务端配置的统一 workspace
- 相对路径与默认 cwd 都会解析到该统一 workspace
- 访问 workspace 外部路径或越过 `APP_HOME/workspaces/skills`、`APP_HOME/workspaces/bin`、`APP_HOME/workspaces/cmd` 等只读部署目录会被拒绝

### 5. 工具暴露与审计

当前 Web runtime：

- 默认只暴露 `load_skill`
- 可以通过 `WEB_ALLOWED_TOOLS` 暴露更多已注册工具
- 每次工具执行都会记录状态与审计摘要
- 对于 `bash` 等工具，会附带解析后的 cwd、命令产物路径、成功摘要或拒绝原因

需要注意的是：这里的工具执行发生在**服务端部署环境**中，而不是用户本地浏览器所在机器上。如果用户请求访问“本地 shell / 本地目录 / 本地文件”，系统仍会返回清晰边界说明，并继续提供替代帮助，例如：

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
- 浏览器模式下健康接口默认只暴露 `load_skill`
- builtin Skill 加载、builtin+user Skill 合并、工具白名单、workspace 隔离、命令产物构建、工具审计记录均已补充测试
- shell / 本地文件 / 用户目录请求会返回明确的能力边界说明

---

## 注意事项

1. 运行前请确保 MySQL、Redis、LLM 服务可用
2. 如果前端跨域访问失败，请检查 `ALLOWED_ORIGIN`
3. 如果登录后接口仍返回 401，请检查 Cookie 是否被浏览器拦截
4. 浏览器聊天模式不是用户本地终端代理；即使启用工具，也是在服务端隔离 workspace 中运行
5. 如果你要让 Skill 调用部署命令，建议先执行 `go run ./cmd/build-artifacts --app-home .`
6. 默认日志文件位于 `APP_HOME/logs/`，默认共享 workspace 位于 `APP_HOME/workspace/`

---

## 后续可继续增强的方向

1. 增加密码重置与用户资料管理
2. 增加 Skill 版本管理
3. 增加会话分页与消息分页
4. 增加更细粒度的能力面板与权限控制
5. 增加部署脚本 / Docker Compose
