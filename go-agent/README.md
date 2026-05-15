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
- 启动时从**已解析的 runtime workspace** 下加载内置 Skill catalog：优先 `APP_HOME/output/workspace/skills`
- 运行时合并“共享内置 Skill + 当前用户已启用 Skill”
- 内置 Skill 通过 API 以只读条目暴露，不能被用户修改或删除
- 按用户隔离 Skill 数据

### 平台能力与边界

- 会话、消息、工具调用记录持久化
- Redis 缓存活跃会话上下文
- 多用户数据隔离
- 浏览器端仅暴露显式允许的已注册工具
- `WORKSPACE_ROOT` 未显式配置时，runtime 会优先使用 `APP_HOME/output/workspace`；若部署产物不存在，则回退到 `APP_HOME/workspace`
- 默认 Web runtime 当前会暴露 `load_skill,bash,read_file,write_file,edit_file`；也可以通过 `WEB_ALLOWED_TOOLS` 显式覆盖
- 工具默认运行在服务端解析后的统一 `WORKSPACE_ROOT` 下，并拒绝越权路径
- 命令产物目录会从解析后的 runtime workspace 派生，即 `<resolved-workspace>/bin` 与 `<resolved-workspace>/cmd`
- 工具调用会记录审计摘要，包括 cwd、命令产物路径、命令产物来源（deployment/local/custom）、结果摘要或拒绝原因
- 浏览器端不会访问**用户本地机器**的 shell、目录或文件；如果启用工具，访问的也是服务端隔离 workspace

---

## 目录结构

```text
go-agent/
├── main.go                     # 默认入口：启动 Web 服务
├── cmd/
│   ├── build-artifacts/
│   │   └── main.go             # 兼容入口：构建 workspace 命令产物与脚本资源
│   └── web/
│       └── main.go             # Web 服务备用入口（等价启动方式）
├── build.sh                    # 标准部署打包脚本
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
├── workspace/                  # 本地调试时使用的 runtime workspace（源码侧回退目录）
│   ├── bin/                    # 本地调试态命令产物
│   ├── cmd/                    # 本地调试态脚本资源
│   └── skills/                 # 本地调试态 builtin Skill catalog
└── output/                     # build.sh 生成的部署输出目录
    ├── bin/
    └── workspace/              # 部署态 runtime workspace（优先使用）
        ├── bin/
        ├── cmd/
        └── skills/
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

后端会优先读取环境变量；如果 LLM 相关环境变量未设置，则回退到 `config.json`。Web 配置会在启动时把 `APP_HOME`、builtin skills 目录、命令产物目录和 workspace 根目录解析为绝对路径；日志文件会写入 `APP_HOME/logs/`。无论是执行部署构建命令还是直接启动 Web 服务，程序都会自动创建这套目录结构。

运行时 workspace 的默认解析顺序：

1. 如果显式设置 `WORKSPACE_ROOT`，优先使用该值
2. 否则如果 `APP_HOME/output/workspace` 存在，优先使用它

`BUILTIN_SKILLS_DIR`、`COMMAND_BIN_DIR`、`COMMAND_SCRIPT_DIR` 在未显式配置时，都会从上面解析出的 runtime workspace 自动派生。

### `config.json` 示例

```json
{
  "base_url": "https://api.deepseek.com",
  "api_key": "your-api-key",
  "model_id": "deepseek-chat",
  "web_allowed_tools": "load_skill,bash,read_file,write_file,edit_file"
}
```

如果你希望覆盖默认路径解析，也可以额外显式配置：

```json
{
  "workspace_root": "custom/runtime-workspace",
  "builtin_skills_dir": "custom/runtime-workspace/skills",
  "command_bin_dir": "custom/runtime-workspace/bin",
  "command_script_dir": "custom/runtime-workspace/cmd"
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

BUILTIN_SKILLS_DIR=
COMMAND_BIN_DIR=
COMMAND_SCRIPT_DIR=
WORKSPACE_ROOT=
WEB_ALLOWED_TOOLS=load_skill,bash,read_file,write_file,edit_file

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
- `WORKSPACE_ROOT` 如果未显式设置：在源码目录运行时会优先解析到 `APP_HOME/output/workspace`，否则回退到 `APP_HOME/workspace`
- `BUILTIN_SKILLS_DIR` 默认从解析后的 `WORKSPACE_ROOT/skills` 派生
- `COMMAND_BIN_DIR` 默认从解析后的 `WORKSPACE_ROOT/bin` 派生
- `COMMAND_SCRIPT_DIR` 默认从解析后的 `WORKSPACE_ROOT/cmd` 派生
- 日志默认写入 `APP_HOME/logs/session_<timestamp>.log`
- `WEB_ALLOWED_TOOLS` 默认是 `load_skill,bash,read_file,write_file,edit_file`；只有显式允许且已注册的工具才会暴露给模型
- 即使启用 `bash` / `read_file` / `write_file` / `edit_file`，它们访问的也是服务端 workspace，而不是用户本地电脑
- 如果未提供 MySQL / Redis 环境变量，程序会使用代码中的默认值构造连接

---

## 启动方式

### 1）执行标准部署打包脚本（推荐在部署阶段执行）

```bash
cd go-agent
./build.sh
```

这一步会：

- 清理并重新生成 `output/`
- 编译主服务到 `output/bin/go-agent`
- 发现 `cmd/*/main.go` 并编译到 `output/workspace/bin/`
- 复制 `.py` / `.sh` 等脚本资源到 `output/workspace/cmd/`
- 复制内置 skills 到 `output/workspace/skills/`
- 复制 `config.json` 到 `output/`（如果存在）

如果你只是想在当前目录补齐 `workspace/bin` 和 `workspace/cmd`，也可以继续使用兼容入口：

```bash
cd go-agent
go run ./cmd/build-artifacts --app-home .
```

### 2）启动 Go 后端（默认入口）

```bash
cd go-agent
go run .
```

启动时也会自动补齐 `APP_HOME/logs/` 以及当前解析到的 runtime workspace 目录树。

如果你是在源码目录本地调试，通常会落到：

- `APP_HOME/workspace/`

如果你是在部署产物目录运行，并且存在 `output/workspace/`，则默认会落到：

- `APP_HOME/output/workspace/`

### 2.1）云端部署推荐方式（生产环境）

推荐把 `go-agent` 当作一个**长期运行的 Web 服务**部署，而不是在云端直接以源码目录 `go run .` 的方式启动。

推荐流程：

1. **在构建机或 CI 中执行打包**

```bash
cd go-agent
./build.sh
```

2. **把整个 `output/` 目录发布到云端机器**，例如发布到：

```text
/srv/go-agent/
├── bin/
│   └── go-agent
└── workspace/
    ├── bin/
    ├── cmd/
    └── skills/
```

3. **在云端以部署产物目录作为 `APP_HOME` 启动服务**

```bash
export APP_HOME=/srv/go-agent
export SERVER_ADDR=:8080
export ALLOWED_ORIGIN=https://your-frontend.example.com

export OPENAI_BASE_URL=https://api.deepseek.com
export OPENAI_API_KEY=your-api-key
export MODEL_ID=deepseek-chat

export MYSQL_HOST=your-mysql-host
export MYSQL_PORT=3306
export MYSQL_USER=your-mysql-user
export MYSQL_PASSWORD=your-mysql-password
export MYSQL_DATABASE=vibe_coding

export REDIS_ADDR=your-redis-host:6379
export REDIS_PASSWORD=
export REDIS_DB=0

export JWT_SECRET=replace-with-production-secret

"${APP_HOME}/bin/go-agent"
```

在这个部署模型下：

- `APP_HOME` 应该指向**部署包根目录**，例如 `/srv/go-agent`
- 云端推荐的最终目录结构是：
  - `APP_HOME/bin/go-agent`
  - `APP_HOME/workspace/bin`
  - `APP_HOME/workspace/cmd`
  - `APP_HOME/workspace/skills`
- 也就是说，云端运行时真正使用的是部署包中的 `APP_HOME/workspace/`，而不是源码仓库里的 `go-agent/workspace/`
- README 前面提到的 `APP_HOME/output/workspace` 优先规则，主要是为了兼容**源码目录本地调试**场景

可以简单理解为：

- **本地源码运行**：`APP_HOME=go-agent/`，runtime 可能命中 `go-agent/output/workspace` 或 `go-agent/workspace`
- **云端部署运行**：`APP_HOME=/srv/go-agent`，runtime 直接命中部署包中的 `APP_HOME/workspace`

### 2.2）systemd 示例（推荐）

如果你在 Linux 云主机上部署，推荐使用 `systemd` 托管服务：

```ini
[Unit]
Description=go-agent web service
After=network.target

[Service]
Type=simple
WorkingDirectory=/srv/go-agent
Environment=APP_HOME=/srv/go-agent
Environment=SERVER_ADDR=:8080
Environment=ALLOWED_ORIGIN=https://your-frontend.example.com
Environment=OPENAI_BASE_URL=https://api.deepseek.com
Environment=OPENAI_API_KEY=your-api-key
Environment=MODEL_ID=deepseek-chat
Environment=MYSQL_HOST=your-mysql-host
Environment=MYSQL_PORT=3306
Environment=MYSQL_USER=your-mysql-user
Environment=MYSQL_PASSWORD=your-mysql-password
Environment=MYSQL_DATABASE=vibe_coding
Environment=REDIS_ADDR=your-redis-host:6379
Environment=REDIS_PASSWORD=
Environment=REDIS_DB=0
Environment=JWT_SECRET=replace-with-production-secret
ExecStart=/srv/go-agent/bin/go-agent
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
```

部署后常用命令：

```bash
sudo systemctl daemon-reload
sudo systemctl enable go-agent
sudo systemctl restart go-agent
sudo systemctl status go-agent
sudo journalctl -u go-agent -f
```

### 2.3）Nginx / 负载均衡建议

如果前端与后端分域部署，建议在 Nginx 或网关层处理：

- HTTPS 终止
- `/api/*` 反向代理到 `go-agent`
- SSE 长连接超时调优（`/api/conversations/:id/stream`）
- 同域 Cookie / CORS 头配置

至少需要注意：

- 反向代理不要过早切断 SSE 连接
- `ALLOWED_ORIGIN` 要与实际前端域名一致
- 生产环境不要继续使用示例里的默认数据库密码、JWT secret

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
5. 在模型请求 `load_skill` 时返回对应能力正文，并附带解析后的 deployment/runtime 路径提示

### 4. 多用户数据与 workspace 隔离

- 用户只能访问自己的 Skill
- 用户只能访问自己的 Conversation / Message
- 工具调用记录按用户与会话隔离存储
- 如果启用工具，所有用户都共享服务端解析后的统一 runtime workspace
- 相对路径与默认 cwd 都会解析到该统一 workspace
- 访问 workspace 外部路径会被拒绝；运行时命令与脚本解析固定指向当前 runtime workspace 下的 `bin/`、`cmd/`

### 5. 工具暴露与审计

当前 Web runtime：

- 默认暴露 `load_skill,bash,read_file,write_file,edit_file`
- 可以通过 `WEB_ALLOWED_TOOLS` 覆盖或收缩已注册工具集合
- 每次工具执行都会记录状态与审计摘要
- 对于 `bash` 等工具，会附带解析后的 cwd、命令产物路径、命令产物来源、成功摘要或拒绝原因

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
- runtime workspace 解析已覆盖显式覆盖、deployment 优先、本地回退三种分支
- builtin Skill 加载、builtin+user Skill 合并、工具白名单、workspace 隔离、命令产物构建、命令产物来源审计均已补充测试
- shell / 本地文件 / 用户目录请求会返回明确的能力边界说明

---

## 注意事项

1. 运行前请确保 MySQL、Redis、LLM 服务可用
2. 如果前端跨域访问失败，请检查 `ALLOWED_ORIGIN`
3. 如果登录后接口仍返回 401，请检查 Cookie 是否被浏览器拦截
4. 浏览器聊天模式不是用户本地终端代理；即使启用工具，也是在服务端隔离 workspace 中运行
5. 如果你要生成完整部署产物，建议优先执行 `./build.sh`；如果只需在当前 `APP_HOME` 下补齐运行时命令目录，可执行 `go run ./cmd/build-artifacts --app-home .`
6. 默认日志文件位于 `APP_HOME/logs/`；源码目录运行时共享 workspace 会优先解析到 `APP_HOME/output/workspace/`，云端部署时则通常直接使用部署包里的 `APP_HOME/workspace/`

---

## 后续可继续增强的方向

1. 增加密码重置与用户资料管理
2. 增加 Skill 版本管理
3. 增加会话分页与消息分页
4. 增加更细粒度的能力面板与权限控制
5. 增加部署脚本 / Docker Compose
