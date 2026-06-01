# nano_cc

`nano_cc` 是一个浏览器优先的 Go Agent 后端，提供登录、会话、流式聊天、Skill 管理与受控工具调用能力。

## 核心能力

- Web 聊天服务：注册、登录、多会话、SSE 流式响应
- Skill：内置 Skill + 用户自定义 Skill
- Tool：按白名单暴露 `load_skill`、`bash`、`read_file`、`write_file`、`edit_file`
- Workspace：工具只访问服务端 `WORKSPACE_ROOT`，默认拒绝越权路径和危险命令
- 存储：MySQL 持久化用户、会话历史、Skill、工具调用记录；Redis 缓存会话上下文

## 会话历史存储

- 会话以 `user_id + root_message_id` 作为用户内唯一锚点
- `conversations.history_json` 保存完整消息历史 JSON，是新会话历史的事实来源
- 每轮对话成功完成后，会把 user + assistant 的完整历史一次性写回 `history_json`
- `GET /api/conversations/:id` 返回的 `messages` 从 `history_json` 解析得到
- `messages` 表仅作为 legacy 兼容保留，新逻辑不再逐条写入

## Skill 加载规则

- 内置 Skill：服务启动时从 `BUILTIN_SKILLS_DIR` 加载，默认是 `WORKSPACE_ROOT/skills`
- 用户 Skill：每次用户发送消息前，从数据库读取当前用户 `enabled` 状态的 Skill
- 每轮响应会生成一份 Skill snapshot：system prompt 和 `load_skill` 工具都使用同一份 snapshot，保证单轮内一致
- API 返回内置 Skill 时会标记为只读，用户不能修改或删除

相关代码：

- Skill snapshot 构建：`internal/web/runtime/prompt_builder.go`
- 会话编排：`internal/web/runtime/conversation_flow.go`
- DB Skill 查询：`internal/web/storage/skills_repo.go`

## 通用 system prompt

通用 system prompt 的基础内容运行时位于 `WORKSPACE_ROOT/system_prompt.md`。服务启动时读取该文件；开启 agent loop 时，再拼接当前用户、workspace、工具、Skill 等动态段，组成最终 system prompt。

默认路径是 `WORKSPACE_ROOT/system_prompt.md`，也可以通过 `SYSTEM_PROMPT_PATH` 或 `WORKSPACE_ROOT/config.json` 中的 `system_prompt_path` 覆盖。

运行时会额外拼接的动态段包括：

- 当前响应场景：`surface`
- 当前 workspace：`WORKSPACE_ROOT`
- 本轮可用工具列表
- 本轮 Skill snapshot 描述
- 可选 memory 段

通常只需要编辑 `WORKSPACE_ROOT/system_prompt.md` 中的基础部分；动态段为空时不会写入最终 prompt。

## 目录结构

```text
go-agent/
├── main.go                  # Web 服务入口
├── config.json              # 构建时同步到 output/workspace/config.json，可选
├── system_prompt.md         # 构建时同步到 output/workspace/system_prompt.md
├── cmd/                     # build.sh 会逐个编译其直接子目录
├── internal/
│   ├── assistant/           # system prompt
│   ├── config/              # 配置与 runtime 路径
│   ├── sessions/            # Skill loader / merge / render
│   ├── tools/               # 工具定义与执行
│   └── web/                 # HTTP、鉴权、runtime、storage
├── skills/                  # 源码内置 Skill
├── workspace/               # 本地 runtime workspace
│   ├── config.json          # 运行时配置，可选
│   ├── system_prompt.md     # 运行时读取的基础 system prompt
│   ├── logs/                # 会话日志、调试日志
│   ├── cmd/                 # 构建后的 runtime 命令
│   └── skills/              # 内置 Skill
└── output/                  # build.sh 生成的部署产物
```

前端位于仓库根目录的 `web/`。

## 环境要求

- Go 1.21+
- Node.js / npm
- OpenAI 兼容模型服务
- MySQL
- Redis

## 配置

后端优先读取环境变量；LLM 配置缺失时可回退到 `WORKSPACE_ROOT/config.json`。

最常用配置：

```bash
OPENAI_BASE_URL=https://api.deepseek.com
OPENAI_API_KEY=your-api-key
MODEL_ID=deepseek-chat

SERVER_ADDR=:8080
ALLOWED_ORIGIN=http://localhost:5173
APP_HOME=/path/to/go-agent
WORKSPACE_ROOT=
SYSTEM_PROMPT_PATH=

WEB_ALLOWED_TOOLS=load_skill,bash,read_file,write_file,edit_file
BASH_ALLOW_OUTSIDE_WORKSPACE=false
BASH_ALLOW_DANGEROUS_COMMANDS=false

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

路径规则：

- `APP_HOME` 默认是当前目录
- `WORKSPACE_ROOT` 未设置时使用 `APP_HOME/workspace`
- `SYSTEM_PROMPT_PATH` 未设置时使用 `WORKSPACE_ROOT/system_prompt.md`
- `BUILTIN_SKILLS_DIR` 默认是 `WORKSPACE_ROOT/skills`
- `COMMAND_BIN_DIR` 默认是 `WORKSPACE_ROOT/bin`
- `COMMAND_SCRIPT_DIR` 默认是 `WORKSPACE_ROOT/cmd`
- 日志固定写入 `WORKSPACE_ROOT/logs`

`WORKSPACE_ROOT/config.json` 示例：

```json
{
  "base_url": "https://api.deepseek.com",
  "api_key": "your-api-key",
  "model_id": "deepseek-chat",
  "system_prompt_path": "system_prompt.md",
  "web_allowed_tools": "load_skill,bash,read_file,write_file,edit_file",
  "bash_allow_outside_workspace": false,
  "bash_allow_dangerous_commands": false
}
```

## 本地启动

```bash
cd go-agent

export OPENAI_BASE_URL=https://api.deepseek.com
export OPENAI_API_KEY=your-api-key
export MODEL_ID=deepseek-chat

export MYSQL_HOST=127.0.0.1
export MYSQL_PORT=3306
export MYSQL_USER=root
export MYSQL_PASSWORD=your-password
export MYSQL_DATABASE=vibe_coding

export REDIS_ADDR=127.0.0.1:6379
export JWT_SECRET=replace-with-your-own-secret

go run .
```

默认后端地址：`http://localhost:8080`。

启动前端：

```bash
cd web
npm install
npm run dev
```

默认前端地址：`http://localhost:5173`。

## 部署

构建发布包：

```bash
cd go-agent
./build.sh
```

构建结果位于 `output/`：

```text
output/
├── bin/go-agent
└── workspace/
    ├── cmd/                 # cmd/ 下每个直接子目录编译出的命令
    ├── config.json          # 从根目录 config.json 同步，可选
    ├── logs/                # 日志目录
    ├── skills/              # 从 skills/ 同步的内置 Skill
    └── system_prompt.md     # 从根目录 system_prompt.md 同步
```

`build.sh` 会直接完成全部构建：主服务二进制输出到 `output/bin/`，`cmd/` 下每个直接子目录编译到 `output/workspace/cmd/`，`skills/` 同步到 `output/workspace/skills/`，根目录 `system_prompt.md` 与可选的 `config.json` 同步到 `output/workspace/`。

云端启动示例：

```bash
export APP_HOME=/srv/go-agent
export SERVER_ADDR=:8080
export ALLOWED_ORIGIN=https://your-frontend.example.com
export OPENAI_BASE_URL=https://api.deepseek.com
export OPENAI_API_KEY=your-api-key
export MODEL_ID=deepseek-chat
export MYSQL_HOST=your-mysql-host
export MYSQL_USER=your-mysql-user
export MYSQL_PASSWORD=your-mysql-password
export MYSQL_DATABASE=vibe_coding
export REDIS_ADDR=your-redis-host:6379
export JWT_SECRET=replace-with-production-secret

${APP_HOME}/bin/go-agent
```

## API 简表

- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/logout`
- `GET /api/me`
- `GET /api/skills`
- `POST /api/skills`
- `GET /api/skills/:id`
- `PUT /api/skills/:id`
- `PATCH /api/skills/:id`
- `DELETE /api/skills/:id`
- `GET /api/conversations`
- `POST /api/conversations`
- `GET /api/conversations/:id`
- `POST /api/conversations/:id/stream`
- `GET /api/health`

## 开发验证

```bash
cd go-agent
go test ./...
./build.sh
```

前端：

```bash
cd web
npm run typecheck
npm run build
```

## 注意事项

- 运行前确保 MySQL、Redis、LLM 服务可用
- 前端跨域失败时检查 `ALLOWED_ORIGIN`
- 登录后仍返回 401 时检查 Cookie 策略
- 浏览器聊天不是用户本地终端代理；工具运行在服务端 workspace
- 生产环境不建议开启 `BASH_ALLOW_OUTSIDE_WORKSPACE` 或 `BASH_ALLOW_DANGEROUS_COMMANDS`
