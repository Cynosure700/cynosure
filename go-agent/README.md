# nano_cc

`nano_cc` 是一个浏览器优先的 Go Agent 后端，提供登录、会话、流式聊天、Skill 管理、受控工具调用、子 Agent 编排、上下文压缩、跨会话记忆、工具审计与可扩展 Hook 能力。

## 核心能力

- Web 聊天服务：注册、登录、多会话、SSE 流式响应，支持推理过程 `reasoning_content` 持久化与下发
- Skill：内置 Skill + 用户自定义 Skill；每轮对话生成 Skill snapshot，保证 prompt 与工具读取一致
- Tool：默认启用 `load_skill`，可通过白名单开放 `bash`、`read_file`、`write_file`、`edit_file`、`todo_write`、`spawn_subagent`、`update_memory`；`read_persisted_output` 始终自动注入给主 Agent
- 子 Agent：通过 `spawn_subagent` 派生隔离的子任务 Agent，使用独立消息列表与工具集（剔除 `spawn_subagent`，禁止嵌套），仅把最终 summary 返回父 Agent
- 上下文压缩：超长历史在请求模型前自动压缩（大工具结果落盘、消息窗口裁剪、保留最近工具结果、全量历史摘要），并可通过 `read_persisted_output` 回取落盘内容
- 记忆系统：跨会话沉淀「用户档案卡」与「近期话题」，在 system prompt 中动态注入
- Workspace：运行时配置、system prompt、日志、Skill 和命令资源都围绕服务端 `WORKSPACE_ROOT` 组织
- 安全与审计：工具默认拒绝越权路径和危险命令，并记录工具调用、解析后的 cwd、命令路径、产物来源与拒绝原因
- Hook 编排：在用户输入、工具执行前后、循环停止前提供 Hook 扩展点，默认 Hook 承载持久化、审计、SSE 与标题推断
- 存储：MySQL 持久化用户、会话、消息、Skill、工具调用、子 Agent 消息、压缩产物与记忆；Redis 缓存会话上下文；Elasticsearch 提供通用文档存储

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

## 工具集

| 工具 | 说明 | 默认白名单 |
| --- | --- | --- |
| `load_skill` | 加载 Skill 正文与元数据 | 是（白名单为空时也会兜底保留） |
| `bash` | 在 workspace 内执行命令，受越权路径与危险命令限制 | 是 |
| `read_file` / `write_file` / `edit_file` | 受控文件读写 | 是 |
| `todo_write` | 维护任务清单 | 见下方说明 |
| `spawn_subagent` | 派生子 Agent 完成隔离子任务 | 是 |
| `update_memory` | 写入/覆盖当前用户档案卡 | 是 |
| `read_persisted_output` | 回取被压缩落盘的工具产物 | 自动注入，无需白名单 |

工具定义见 `internal/tools/definitions.go`，执行分发见 `internal/web/runtime/tool_registry.go`（`spawn_subagent`、`update_memory` 由 Service 特殊处理，其余走 `internal/tools/handlers.go` 的 Dispatch）。

> 白名单由 `web_allowed_tools` 控制：仓库内 `config.json` 默认值为 `load_skill,bash,read_file,write_file,edit_file,spawn_subagent,update_memory`；当该字段为空时，代码兜底默认值为 `load_skill,bash,read_file,write_file,edit_file,todo_write,update_memory`。两者均可按需调整。

## 子 Agent

`spawn_subagent` 用于把一段可隔离的复杂任务交给独立子 Agent 执行：

- 子 Agent 拥有全新的消息列表与独立工具注册表（不含 `spawn_subagent`），最多 20 轮
- 禁止嵌套：子 Agent 再次调用 `spawn_subagent` 会被拒绝
- 子 Agent 的每条消息写入 `subagent_messages` 表，最终只把 summary 返回父 Agent

相关代码：`internal/web/runtime/subagent.go`、`internal/web/storage/subagent_messages_repo.go`。

## 上下文压缩

请求模型前，运行时会按顺序应用多个压缩策略，控制上下文长度：

1. 大工具结果落盘：超阈值的 tool_result 持久化到 `persisted_outputs`，上下文仅保留预览与 `<persisted-output>` 标记
2. 消息窗口裁剪：消息过多时保留首尾、裁剪中段
3. 最近工具结果保留：仅保留最近若干条完整工具结果
4. 全量历史摘要：对历史整体做摘要，结果写入 `context_summaries`

落盘的工具产物可由模型通过 `read_persisted_output` 回取。

相关代码：`internal/web/runtime/context_compression.go`、`internal/web/runtime/compression/`、`internal/web/storage/persisted_outputs_repo.go`。

## 记忆系统

跨会话沉淀两类记忆并注入 system prompt：

- 用户档案卡：`update_memory` 工具写入 `user_profiles` 表（每用户一条，整体覆盖）
- 近期话题：每轮对话后用一次 LLM 调用把历史提炼为若干话题短语，写入 `conversation_topics` 表

`prompt_builder.go` 的 `buildMemorySection` 会聚合用户档案卡与最近若干会话的话题，生成 `### 用户档案卡` / `### 近期聊过的话题` 段落注入 prompt（档案注入有字节上限）。

相关代码：`internal/web/runtime/topic_memory.go`、`internal/web/storage/user_profiles_repo.go`、`internal/web/storage/conversation_topics_repo.go`。

## Elasticsearch 通用文档存储

`Store` 集成了 Elasticsearch 客户端，提供与具体表无关的通用文档存储 API：

- 索引：`EnsureESIndex`、`ESIndexExists`
- 写入：`IndexDocument`（按 ID upsert）、`BulkIndexDocuments`（批量）
- 查询：`GetDocument`（按 ID）、`SearchDocuments`（Query DSL）、`MatchDocuments`（单字段匹配便捷封装）
- 删除：`DeleteDocument`

客户端在 `NewStore` 初始化，`HealthCheck` 会对 ES 执行 Ping。地址默认 `http://1.12.217.28:9200`，可通过 `config.json` 的 `es_addresses` 或环境变量 `ES_ADDRESSES` 覆盖；用户名/密码见配置一节。

相关代码：`internal/web/storage/es_repo.go`、`internal/web/storage/store.go`。

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

相关代码：`internal/config/paths.go`。

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
- prompt 文件路径解析：`internal/config/paths.go`
- system prompt 构建：`internal/web/runtime/prompt_builder.go`

## Agent Loop 与 Hook

Web runtime 的每轮对话按以下流程执行：

1. 读取会话历史，触发 `UserPromptSubmit` Hook
2. 构建当前用户 Skill snapshot 与 system prompt（含记忆段落）
3. 请求模型前对上下文执行压缩
4. 请求 OpenAI 兼容模型
5. 如模型发起工具调用，依次触发 `PreToolUse` → 执行工具 → `PostToolUse`
6. 如模型停止工具调用，触发 `Stop` Hook，持久化 assistant 消息并通过 SSE 下发

默认 Hook：

- `UserPromptSubmit`：追加用户消息、推断或刷新会话活动时间
- `PreToolUse`：预解析工具审计信息，例如 cwd、命令路径、命令产物路径与来源
- `PostToolUse`：补充执行结果或拒绝原因，持久化工具调用，发送 SSE 工具事件，追加 tool message 给下一轮模型
- `Stop`：持久化 assistant 消息，发送 SSE assistant 事件，包含 `content` 和 `reasoning_content`

相关代码：

- Hook 管理与默认注册：`internal/web/runtime/hooks/manager.go`
- 用户输入 Hook：`internal/web/runtime/hooks/user_prompt.go`
- 工具 Hook：`internal/web/runtime/hooks/tool.go`
- 停止 Hook：`internal/web/runtime/hooks/stop.go`
- Hook 类型与桥接：`internal/web/runtime/hooks/types.go`、`internal/web/runtime/hook_bridge.go`

## 目录结构

```text
go-agent/
├── main.go                  # Web 服务入口
├── build.sh                 # 构建发布包脚本
├── cmd/                     # 可构建到 workspace/cmd 的 runtime 命令资源
├── config.json              # 根目录配置模板，构建时复制到 output/workspace/config.json
├── system_prompt.md         # 根目录 prompt 模板，构建时复制到 output/workspace/system_prompt.md
├── internal/
│   ├── assistant/           # 默认 system prompt
│   ├── config/              # 配置与 runtime 路径（config.go / web_config.go / env.go / paths.go）
│   ├── idgen/               # ID 生成
│   ├── llm/                 # LLM 客户端
│   ├── logger/              # 日志
│   ├── safety/              # 路径安全
│   ├── sessions/            # Skill loader / merge / render
│   ├── textutil/            # 文本工具
│   ├── tools/               # 工具定义与执行
│   └── web/
│       ├── app/             # HTTP 路由与 handler
│       ├── auth/            # 鉴权服务
│       ├── runtime/         # Agent runtime、Hook、子 Agent、压缩、记忆
│       │   ├── hooks/       # Hook 子包
│       │   └── compression/ # 上下文压缩策略
│       └── storage/         # MySQL / Redis / Elasticsearch 存储
├── skills/                  # 源码内置 Skill
├── workspace/               # 本地 runtime workspace：config、system prompt、logs、skills、cmd、bin
└── output/                  # build.sh 生成的部署产物
```

前端位于仓库根目录的 `web/`。

## 环境要求

- Go 1.26+
- Node.js / npm
- OpenAI 兼容模型服务
- MySQL
- Redis
- Elasticsearch

## 配置

LLM 的 `base_url`、`model_id` 从 `config.json` 读取，`api_key` 从环境变量 `OPENAI_API_KEY` 读取；密钥与密码类敏感信息一律走环境变量，不写入 `config.json`。

环境变量清单：

| 变量 | 用途 |
| --- | --- |
| `OPENAI_API_KEY` | LLM API Key（必需） |
| `MYSQL_PASSWORD` | MySQL 密码 |
| `DATABASE_URL` | 直接覆盖 MySQL DSN（设置后忽略 `mysql_*` 字段） |
| `REDIS_PASSWORD` | Redis 密码 |
| `ES_ADDRESSES` | Elasticsearch 地址，覆盖 `es_addresses` |
| `ES_USERNAME` / `ES_PASSWORD` | Elasticsearch 用户名 / 密码 |
| `JWT_SECRET` | JWT 密钥（默认 `nano-cc-local-secret`，生产务必覆盖） |

常用环境变量示例：

```bash
OPENAI_API_KEY=your-api-key

MYSQL_PASSWORD=your-password
REDIS_PASSWORD=
ES_PASSWORD=

JWT_SECRET=replace-with-your-own-secret
```

`config.json` 示例：

```json
{
  "base_url": "https://api.deepseek.com",
  "model_id": "deepseek-chat",
  "app_home": ".",
  "workspace_root": "workspace",
  "system_prompt_path": "system_prompt.md",
  "builtin_skills_dir": "skills",
  "command_bin_dir": "bin",
  "command_script_dir": "cmd",
  "web_allowed_tools": "load_skill,bash,read_file,write_file,edit_file,spawn_subagent,update_memory",
  "bash_allow_outside_workspace": false,
  "bash_allow_dangerous_commands": false,
  "server_addr": ":8080",
  "allowed_origin": "http://localhost:5173",
  "mysql_host": "1.12.217.28",
  "mysql_port": "3306",
  "mysql_user": "root",
  "mysql_database": "vibe_coding",
  "redis_addr": "1.12.217.28:6379",
  "redis_db": 0,
  "es_addresses": "http://1.12.217.28:9200",
  "es_username": "",
  "session_cookie_name": "nano_cc_session",
  "session_ttl_minutes": 10080
}
```

主要默认值：`server_addr=:8080`、`allowed_origin=http://localhost:5173`、`redis_addr=1.12.217.28:6379`、`es_addresses=http://1.12.217.28:9200`、`session_cookie_name=nano_cc_session`、`session_ttl_minutes=10080`（7 天）。MySQL DSN 默认 `root@tcp(1.12.217.28:3306)/vibe_coding`，密码由 `MYSQL_PASSWORD` 提供。

## 本地启动

```bash
cd go-agent

export OPENAI_API_KEY=your-api-key
export MYSQL_PASSWORD=your-password
export JWT_SECRET=replace-with-your-own-secret

# base_url / model_id / mysql_* / redis_addr / es_addresses 等在 config.json 中配置
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

`build.sh` 会清空并重建 `output/`，编译主程序到 `output/bin/go-agent`，编译 `cmd/` 下各子命令到 `output/bin/`，并把 `config.json`、`system_prompt.md`、`skills/`、`cmd/` 资源复制到对应位置。产物结构：

```text
output/
├── bin/                 # go-agent 主二进制 + 各 cmd 子命令二进制
├── cmd/                 # cmd 源码副本
├── skills/              # 技能文件副本
├── logs/                # 空目录
├── workspace/           # 空目录
├── config.json
└── system_prompt.md
```

云端启动示例：

```bash
export APP_HOME=/srv/go-agent
export WORKSPACE_ROOT=${APP_HOME}/workspace
export OPENAI_API_KEY=your-api-key
export MYSQL_PASSWORD=your-mysql-password
export ES_PASSWORD=your-es-password
export JWT_SECRET=replace-with-production-secret

${APP_HOME}/bin/go-agent
```

部署时建议保留 `output/` 内的 `bin/` 与 `workspace/` 相对布局，并把 `output/config.json`、`output/system_prompt.md` 放入运行时 `WORKSPACE_ROOT`。

## API 简表

- `GET /api/health`
- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/logout`
- `GET /api/me`
- `GET /api/skills`
- `POST /api/skills`
- `GET /api/skills/{id}`
- `PUT /api/skills/{id}`
- `DELETE /api/skills/{id}`
- `GET /api/conversations`
- `POST /api/conversations`
- `GET /api/conversations/{id}`
- `PATCH /api/conversations/{id}`（改标题）
- `DELETE /api/conversations/{id}`
- `POST /api/conversations/{id}/stream`（SSE 流式对话）

除 `register`、`login`、`health` 外的接口均需鉴权。`POST /api/conversations/{id}/stream` 通过 SSE 返回对话事件：工具调用事件由 `PostToolUse` Hook 发送，assistant 事件由 `Stop` Hook 发送。

相关代码：`internal/web/app/routes.go`、`internal/web/app/server.go`。

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

- 运行前确保 MySQL、Redis、Elasticsearch、LLM 服务可用
- 运行时 `config.json` 和 `system_prompt.md` 默认都在 `WORKSPACE_ROOT` 下，不再从项目根目录直接读取
- LLM 的 `base_url`、`model_id` 在 `config.json` 配置，`OPENAI_API_KEY` 必须通过环境变量提供
- 日志默认写入 `WORKSPACE_ROOT/logs`，部署时需确保该目录可写
- 前端跨域失败时检查 `ALLOWED_ORIGIN`
- 登录后仍返回 401 时检查 Cookie 策略
- 浏览器聊天不是用户本地终端代理；工具运行在服务端 workspace
- 生产环境不建议开启 `bash_allow_outside_workspace` 或 `bash_allow_dangerous_commands`
