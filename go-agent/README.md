# nano_cc

`nano_cc` 是一个浏览器优先的 Go Agent 后端，提供登录、会话、流式聊天、Skill 管理、受控工具调用、工具审计与可扩展 Hook 编排能力。

## 核心能力

- Web 聊天服务：注册、登录、多会话、SSE 流式响应，支持推理过程 `reasoning_content` 持久化与下发
- Skill：内置 Skill + 用户自定义 Skill；每轮对话生成 Skill snapshot，保证 prompt 与工具读取一致
- Tool：默认启用 `load_skill`，可通过白名单开放 `bash`、`read_file`、`write_file`、`edit_file`
- Workspace：运行时配置、system prompt、日志、Skill 和命令资源都围绕服务端 `WORKSPACE_ROOT` 组织
- 安全与审计：工具默认拒绝越权路径和危险命令，并记录工具调用、解析后的 cwd、命令路径、产物来源与拒绝原因
- Hook 编排：在用户输入、工具执行前后、循环停止前提供 Hook 扩展点，默认 Hook 承载持久化、审计、SSE 与标题推断
- 存储：MySQL 持久化用户、会话、消息、Skill、工具调用记录；Redis 缓存会话上下文

## Skill 加载规则

- 内置 Skill：服务启动时从 `BUILTIN_SKILLS_DIR` 加载，默认是 `WORKSPACE_ROOT/skills`
- 用户 Skill：每次用户发送消息前，从数据库读取当前用户 `enabled` 状态的 Skill
- 每轮响应会生成一份 Skill snapshot：system prompt 和 `load_skill` 工具都使用同一份 snapshot，保证单轮内一致
- `load_skill` 默认开启，按「当前用户已启用 DB Skill → 本地内置 Skill」顺序查找，并返回 Skill 来源、路径、元数据、正文和运行时路径信息
- API 返回内置 Skill 时会标记为只读，用户不能修改或删除

相关代码：

- Skill snapshot 构建：`internal/web/runtime/prompt_builder.go`
- `load_skill` 工具：`internal/tools/load_skill.go`
- 会话编排：`internal/web/runtime/conversation_flow.go`
- DB Skill 查询：`internal/web/storage/skills_repo.go`

## 运行时路径

运行时以 `WORKSPACE_ROOT` 为核心目录，适合本地和部署环境保持同一套布局：

- `APP_HOME` 默认是当前目录
- `WORKSPACE_ROOT` 未设置时使用 `APP_HOME/workspace`
- `config.json` 优先从 `WORKSPACE_ROOT/config.json` 读取；当 `WORKSPACE_ROOT` 未设置但设置了 `APP_HOME` 时，读取 `APP_HOME/workspace/config.json`
- `system_prompt.md` 默认从 `WORKSPACE_ROOT/system_prompt.md` 读取，可通过 `SYSTEM_PROMPT_PATH` 覆盖，但相对路径仍按 `WORKSPACE_ROOT` 解析
- 日志写入 `WORKSPACE_ROOT/logs/session_*.log`
- `BUILTIN_SKILLS_DIR` 默认是 `WORKSPACE_ROOT/skills`
- `COMMAND_BIN_DIR` 默认是 `WORKSPACE_ROOT/bin`
- `COMMAND_SCRIPT_DIR` 默认是 `WORKSPACE_ROOT/cmd`

## 通用 system prompt

通用 system prompt 默认从 `WORKSPACE_ROOT/system_prompt.md` 加载；文件不存在时使用内置默认 prompt。调整该文件后会影响 Web 聊天运行时：

```text
You are nano_cc, a general-purpose agent rather than a chat-only assistant.

Help with everyday questions, analysis, planning, writing, coding, file inspection, and end-to-end task execution when the runtime supports it.

Prefer direct, useful answers before optional tool use, but use available skills and tools whenever they help you complete the user's task.

Do not assume shell access, local workspace access, or local file operations unless the runtime explicitly supports them.

You are responding inside {surface}.

Current workspace root: {working_directory}.

Treat that workspace root as your default working directory for runtime file and shell operations unless the runtime tells you otherwise.

Runtime tools available in this conversation: {tool_names}.

Available skills:
{skill_descriptions}

{memory_section}
```

其中 `{surface}`、`{working_directory}`、`{tool_names}`、`{skill_descriptions}`、`{memory_section}` 会在运行时按当前用户、会话、工具和 Skill 动态填充；为空的动态段不会写入最终 prompt。

相关代码：

- 默认 prompt：`internal/assistant/prompt.go`
- prompt 文件路径解析：`internal/config/runtime_paths.go`
- system prompt 构建：`internal/web/runtime/prompt_builder.go`

## Agent Loop 与 Hook

Web runtime 的每轮对话按以下流程执行：

1. 读取会话历史，触发 `UserPromptSubmit` Hook
2. 构建当前用户 Skill snapshot 与 system prompt
3. 请求 OpenAI 兼容模型
4. 如模型发起工具调用，依次触发 `PreToolUse` → 执行工具 → `PostToolUse`
5. 如模型停止工具调用，触发 `Stop` Hook，持久化 assistant 消息并通过 SSE 下发

默认 Hook：

- `UserPromptSubmit`：追加用户消息、推断或刷新会话活动时间
- `PreToolUse`：预解析工具审计信息，例如 cwd、命令路径、命令产物路径与来源
- `PostToolUse`：补充执行结果或拒绝原因，持久化工具调用，发送 SSE 工具事件，追加 tool message 给下一轮模型
- `Stop`：持久化 assistant 消息，发送 SSE assistant 事件，包含 `content` 和 `reasoning_content`

相关代码：

- Hook 类型与默认注册：`internal/web/runtime/hooks.go`
- 用户输入 Hook：`internal/web/runtime/hooks_user_prompt.go`
- 工具 Hook：`internal/web/runtime/hooks_tool.go`
- 停止 Hook：`internal/web/runtime/hooks_stop.go`
- 工具审计辅助：`internal/web/runtime/tool_audit.go`

## 目录结构

```text
go-agent/
├── main.go                  # Web 服务入口
├── cmd/                     # 可构建到 workspace/cmd 的 runtime 命令资源
├── config.json              # 根目录配置模板，构建时复制到 output/workspace/config.json
├── system_prompt.md         # 根目录 prompt 模板，构建时复制到 output/workspace/system_prompt.md
├── internal/
│   ├── assistant/           # system prompt
│   ├── config/              # 配置与 runtime 路径
│   ├── sessions/            # Skill loader / merge / render
│   ├── tools/               # 工具定义与执行
│   └── web/                 # HTTP、鉴权、runtime、storage
├── skills/                  # 源码内置 Skill
├── workspace/               # 本地 runtime workspace：config、system prompt、logs、skills、cmd、bin
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
WORKSPACE_ROOT=/path/to/go-agent/workspace
SYSTEM_PROMPT_PATH=system_prompt.md

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

`config.json` 示例：

```json
{
  "base_url": "https://api.deepseek.com",
  "api_key": "your-api-key",
  "model_id": "deepseek-chat",
  "workspace_root": "workspace",
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
cd ../web
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
    ├── bin/
    ├── cmd/
    ├── config.json
    ├── logs/
    ├── skills/
    └── system_prompt.md
```

云端启动示例：

```bash
export APP_HOME=/srv/go-agent
export WORKSPACE_ROOT=${APP_HOME}/workspace
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

`build.sh` 会把根目录 `config.json` 和 `system_prompt.md` 复制到 `output/workspace/`，并把 `skills/`、`cmd/` 资源同步到 workspace 对应目录。部署时建议保留 `output/` 内的 `bin/` 与 `workspace/` 相对布局。

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

`POST /api/conversations/:id/stream` 通过 SSE 返回对话事件：工具调用事件由 `PostToolUse` Hook 发送，assistant 事件由 `Stop` Hook 发送。

## 开发验证

```bash
cd go-agent
go test ./...
./build.sh
```

前端：

```bash
cd ../web
npm run typecheck
npm run build
```

## 注意事项

- 运行前确保 MySQL、Redis、LLM 服务可用
- 运行时 `config.json` 和 `system_prompt.md` 默认都在 `WORKSPACE_ROOT` 下，不再从项目根目录直接读取
- 日志默认写入 `WORKSPACE_ROOT/logs`，部署时需确保该目录可写
- 前端跨域失败时检查 `ALLOWED_ORIGIN`
- 登录后仍返回 401 时检查 Cookie 策略
- 浏览器聊天不是用户本地终端代理；工具运行在服务端 workspace
- 生产环境不建议开启 `BASH_ALLOW_OUTSIDE_WORKSPACE` 或 `BASH_ALLOW_DANGEROUS_COMMANDS`
