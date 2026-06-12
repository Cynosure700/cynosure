# go-agent / nano_cc

`go-agent` 是 `nano_cc` 的浏览器优先 Agent 后端，面向 Web Chat 场景提供登录鉴权、多会话管理、SSE 流式对话、Skill 管理、受控工具调用、MCP 工具接入、子 Agent 编排、上下文压缩、长期记忆、工具审计与可扩展 Hook 能力。

前端项目位于本目录的兄弟目录 `../web`，基于 Vite + React + TypeScript 实现。后端只提供 API，当前不托管前端静态资源。

## 核心能力

- **Web 聊天服务**：注册、登录、登出、当前用户查询、多会话列表/创建/重命名/删除、SSE 流式聊天。
- **流式输出**：支持 `assistant_delta`、`reasoning_delta`、`meta`、最终 `assistant`、`error`、`done` 等 SSE 事件；推理过程 `reasoning_content` 会持久化并下发。
- **Skill 系统**：支持服务端内置 Skill 与用户数据库 Skill；每轮对话生成 Skill snapshot，保证 system prompt 与 `load_skill` 工具读取一致。
- **内置工具**：支持 `load_skill`、`bash`、`read_file`、`write_file`、`edit_file`、`todo_write`、`spawn_subagent`、`read_persisted_output` 等工具，按 `web_allowed_tools` 控制暴露范围。
- **MCP 工具接入**：支持服务端内置 stdio MCP 与用户配置的远程 `sse` / `streamable` MCP；发现到的 MCP 工具会以 `mcp__{server}__{tool}` 形式动态加入模型工具列表。
- **子 Agent**：`spawn_subagent` 可派生隔离子任务 Agent，使用独立消息列表和工具集，禁止嵌套，最终只把 summary 返回父 Agent。
- **上下文压缩**：请求模型前自动压缩超长历史，支持大工具结果落盘、消息窗口裁剪、最近工具结果保留、会话记忆替换、全量历史摘要与 413 后激进压缩。
- **记忆系统**：每轮结束后异步抽取长期记忆，支持 `episodic_memory`、`user_preference`、`semantic` 三类；对话开始前选择相关记忆注入 prompt，并支持用户级记忆开关。
- **模型历史复用**：压缩后的模型历史可持久化到 `conversation_model_histories`，减少后续对话重复压缩成本。
- **安全与审计**：工具默认限制在服务端 `workspace_root` 内运行，危险命令和越权路径默认拒绝；工具调用、拒绝原因、cwd、命令路径等信息会被审计记录。
- **Hook 编排**：在用户输入、工具执行前后、循环停止前提供 Hook 扩展点，默认 Hook 承载持久化、审计、SSE 与标题推断等流程。
- **存储集成**：MySQL 持久化用户、会话、消息、Skill、MCP 配置、工具调用、子 Agent 消息、压缩产物与记忆；Redis 用于缓存和锁；Elasticsearch 提供通用文档存储 API。

## 项目结构

```text
go-agent/
├── main.go                  # 后端入口
├── build.sh                 # 后端发布包构建脚本
├── config.json              # 后端配置模板，构建时复制到 output/config.json
├── system_prompt.md         # 基础 identity prompt，启动时必须存在
├── skills/                  # 服务端内置 Skill
├── workspace/               # 本地 runtime workspace
├── output/                  # build.sh 生成的后端部署产物
└── internal/
    ├── assistant/           # system prompt 拼装
    ├── config/              # 配置与 runtime 路径
    ├── idgen/               # ID 生成
    ├── llm/                 # OpenAI 兼容 LLM 客户端与错误分类
    ├── logger/              # 日志
    ├── safety/              # 路径安全
    ├── sessions/            # Skill loader / merge / render
    ├── textutil/            # 文本工具
    ├── tools/               # 内置工具定义、参数校验与执行
    └── web/
        ├── app/             # HTTP 路由与 handler
        ├── auth/            # Cookie/JWT 鉴权服务
        ├── mcp/             # MCP 配置加载、连接管理与 transport
        ├── runtime/         # Agent loop、Hook、子 Agent、压缩、记忆
        │   ├── compression/ # 上下文压缩策略
        │   └── hooks/       # Hook 子包
        └── storage/         # MySQL / Redis / Elasticsearch 存储
```

前端目录：

```text
../web/
├── package.json
├── vite.config.ts
└── src/
    ├── App.tsx              # 主界面：聊天、会话、Skill、MCP、记忆开关
    ├── api.ts               # API 封装与 SSE 解析
    ├── main.tsx
    └── styles.css
```

## 环境要求

- Go 1.26.1（见 `go.mod`）
- Node.js / npm（前端开发与构建）
- OpenAI 兼容模型服务
- MySQL
- Redis
- Elasticsearch

## 配置

### 配置文件读取规则

- 后端固定从进程当前工作目录读取 `config.json`。
- `config.json` 中必须配置 `base_url` 和 `model_id`，否则启动失败。
- `app_home` 默认是 `.`，会解析为绝对路径；运行时目录均以 `app_home` 为基准。
- `system_prompt_path` 默认是 `app_home/system_prompt.md`，当前代码要求该文件存在；文件缺失会导致启动失败。
- `builtin_skills_dir`、`command_bin_dir`、`command_script_dir` 若显式配置，必须分别解析到 `app_home/skills`、`app_home/bin`、`app_home/cmd`。
- `workspace_root` 是工具执行与文件访问的根目录，也是默认越权边界。

启动时会自动创建：`app_home`、`logs/`、`skills/`、`bin/`、`cmd/`、`workspace/`。

### 敏感配置

LLM 的 `base_url`、`model_id` 从 `config.json` 读取；`api_key` 从环境变量 `OPENAI_API_KEY` 读取。MySQL、Redis、Elasticsearch、JWT 等敏感信息建议通过环境变量传入，生产环境不要依赖代码内置的本地开发默认值。

| 变量 | 用途 |
| --- | --- |
| `OPENAI_API_KEY` | LLM API Key |
| `DATABASE_URL` | 直接覆盖 MySQL DSN；设置后忽略 `mysql_*` 字段 |
| `MYSQL_PASSWORD` | MySQL 密码 |
| `REDIS_PASSWORD` | Redis 密码 |
| `ES_ADDRESSES` | Elasticsearch 地址，覆盖 `es_addresses` |
| `ES_USERNAME` / `ES_PASSWORD` | Elasticsearch 用户名 / 密码 |
| `JWT_SECRET` | JWT 签名密钥，生产环境必须覆盖 |

常用环境变量示例：

```bash
export OPENAI_API_KEY=your-api-key
export MYSQL_PASSWORD=your-mysql-password
export REDIS_PASSWORD=your-redis-password
export ES_PASSWORD=your-es-password
export JWT_SECRET=replace-with-your-own-secret
```

### `config.json` 示例

```json
{
  "base_url": "https://api.deepseek.com",
  "model_id": "deepseek-v4-flash",
  "app_home": ".",
  "system_prompt_path": "system_prompt.md",
  "workspace_root": "workspace",
  "builtin_skills_dir": "skills",
  "command_bin_dir": "bin",
  "command_script_dir": "cmd",
  "web_allowed_tools": "load_skill,bash,read_file,write_file,edit_file,spawn_subagent",
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

主要默认值：`server_addr=:8080`、`allowed_origin=http://localhost:5173`、`session_cookie_name=nano_cc_session`、`session_ttl_minutes=10080`（7 天）。

## 后端本地启动

```bash
cd go-agent

export OPENAI_API_KEY=your-api-key
export MYSQL_PASSWORD=your-mysql-password
export REDIS_PASSWORD=your-redis-password
export ES_PASSWORD=your-es-password
export JWT_SECRET=replace-with-your-own-secret

# base_url / model_id / mysql_* / redis_addr / es_addresses 等在 config.json 中配置
go run .
```

默认后端地址：`http://localhost:8080`。

启动过程会连接 MySQL、Redis、Elasticsearch 并执行迁移；请先确保这些依赖可用。

## 前端 Web

前端位于 `../web`，技术栈为 Vite + React + TypeScript。

```bash
cd ../web
npm install
npm run dev
```

默认前端地址：`http://localhost:5173`。

前端默认请求 `http://localhost:8080`，可通过 `VITE_API_BASE` 覆盖：

```bash
VITE_API_BASE=http://localhost:8080 npm run dev
```

前端所有 API 请求都会带 `credentials: "include"`，依赖后端 Cookie 会话。跨域开发时，后端 `config.json` 的 `allowed_origin` 必须与前端地址一致，否则登录态 Cookie 或 CORS 请求会失败。

前端能力概览：

- 登录 / 注册 / 退出。
- 多会话列表、新建、切换、重命名、删除。
- SSE 流式聊天，展示推理过程、生成耗时、工具调用次数、上下文 token 用量。
- 基础 Markdown-like 渲染：标题、列表、代码块、表格、行内代码、加粗等。
- 用户 Skill 管理：创建、编辑、启用/禁用、删除。
- 用户 MCP Server 管理：创建、编辑、启用/禁用、删除、测试连接。
- 用户级记忆开关。

## Skill 系统

- **内置 Skill**：服务启动时从 `builtin_skills_dir` 加载，默认 `app_home/skills`。
- **用户 Skill**：通过 `/api/skills` 管理，存储在 MySQL；每轮对话前只读取当前用户 `enabled` 状态的 Skill。
- **Skill snapshot**：每轮响应会生成一份 snapshot，system prompt 与 `load_skill` 工具都使用同一份 snapshot，避免单轮内 Skill 描述与正文不一致。
- **加载优先级**：`load_skill` 先查当前用户已启用的 DB Skill，再回退到本地内置 Skill。
- **冲突校验**：用户 Skill slug 不能与内置 Skill 冲突。

注意：当前 `/api/skills` 面向用户 DB Skill 管理，不会把内置 Skill 作为只读项返回给前端列表。

相关代码：

- `internal/web/runtime/prompt_builder.go`
- `internal/tools/load_skill.go`
- `internal/web/app/skill_handlers.go`
- `internal/web/storage/skills_repo.go`

## 工具系统

### 内置工具

| 工具 | 说明 | 暴露方式 |
| --- | --- | --- |
| `load_skill` | 加载 Skill 正文与元数据 | 由 `web_allowed_tools` 控制 |
| `bash` | 在 workspace 内执行 shell 命令，受越权路径与危险命令限制 | 由 `web_allowed_tools` 控制 |
| `read_file` | 读取 workspace 内文件 | 由 `web_allowed_tools` 控制 |
| `write_file` | 写入 workspace 内文件 | 由 `web_allowed_tools` 控制 |
| `edit_file` | 按精确文本替换编辑文件 | 由 `web_allowed_tools` 控制 |
| `todo_write` | 维护多步骤任务清单 | 由 `web_allowed_tools` 控制 |
| `spawn_subagent` | 派生隔离子 Agent | 由 `web_allowed_tools` 控制 |
| `read_persisted_output` | 读取上下文压缩落盘的大工具结果 | 主 Agent 自动追加，无需配置 |

当前仓库 `config.json` 默认暴露：

```text
load_skill,bash,read_file,write_file,edit_file,spawn_subagent
```

如果 `web_allowed_tools` 配置为空，代码层面的默认值是：

```text
load_skill,bash,read_file,write_file,edit_file,todo_write
```

`read_persisted_output` 会自动追加给主 Agent；子 Agent 工具集会移除 `spawn_subagent`，禁止嵌套派生。

相关代码：

- 工具定义：`internal/tools/definitions.go`
- 参数校验：`internal/tools/validation.go`
- 工具分发：`internal/tools/handlers.go`
- runtime 注册与特殊分发：`internal/web/runtime/tool_registry.go`

### 安全边界

- 文件读写、编辑和命令执行都以服务端 `workspace_root` 为默认边界。
- `bash_allow_outside_workspace=false` 时拒绝访问 workspace 外路径。
- `bash_allow_dangerous_commands=false` 时拒绝危险命令。
- 浏览器聊天不是用户本地终端代理；所有工具都运行在服务端环境中。

## MCP

MCP 分为两类：

1. **服务端内置 MCP**：通过 `app_home/mcp_config.json` 配置，当前只支持 stdio。
2. **用户 MCP**：通过前端或 `/api/mcp-servers` 管理，当前支持远程 `sse` / `streamable` transport。

### 内置 MCP 配置

服务启动时会尝试读取：

```text
{app_home}/mcp_config.json
```

文件不存在时会跳过内置 MCP，不影响启动。示例：

```json
{
  "mcp_servers": [
    {
      "name": "demo",
      "command": "/path/to/mcp-server",
      "args": ["--flag"],
      "env": { "TOKEN": "xxx" },
      "enabled": true
    }
  ]
}
```

内置 MCP 不允许配置 `transport`、`url`、`headers`，这些字段只用于用户远程 MCP。

注意：当前 `build.sh` 不会复制 `mcp_config.json` 到 `output/`。如果部署需要内置 MCP，请手动把该文件放到部署环境的 `app_home` 下。

### 用户 MCP

用户 MCP 通过 API 存储在 `mcp_servers` 表：

- `transport` 取值：`sse` 或 `streamable`。
- `url` 必填。
- `headers` 可选，用于鉴权等 HTTP 请求头。
- 更新、删除、启用/禁用后会使该用户的 MCP 连接失效，下次对话按最新配置重连。
- `POST /api/mcp-servers/{id}/test` 会临时连接并返回发现到的工具名。

MCP Manager 行为：

- 用户维度懒加载和缓存连接。
- 配置变更自动重连。
- 单个 MCP 连接失败不会影响其他服务器。
- 调用超时 120 秒。
- 空闲 10 分钟后回收连接。
- 内置 stdio MCP 调用失败时会尝试重连后重试。

发现到的工具会转换为 OpenAI function tool，并按以下格式命名：

```text
mcp__{sanitized_server_name}__{original_tool_name}
```

相关代码：

- `internal/web/mcp/config.go`
- `internal/web/mcp/manager.go`
- `internal/web/mcp/transport.go`
- `internal/web/app/mcp_handlers.go`
- `internal/web/runtime/tool_registry.go`

## 子 Agent

`spawn_subagent` 用于把可隔离的复杂任务交给独立子 Agent：

- 子 Agent 拥有全新消息列表，不共享父对话历史。
- 子 Agent 继承当前轮 Skill snapshot。
- 子 Agent 工具集中移除 `spawn_subagent`，禁止嵌套。
- 子 Agent 最多执行 20 轮。
- `cwd` 必须解析到 workspace 内。
- 子 Agent 的消息写入 `subagent_messages` 表。
- 最终只把 summary 返回父 Agent。

相关代码：`internal/web/runtime/subagent.go`、`internal/web/storage/subagent_messages_repo.go`。

## 上下文压缩

请求模型前，runtime 会估算上下文并按策略链压缩：

1. **大工具结果落盘**：超过阈值的 tool result 写入 `persisted_outputs`，上下文只保留预览和 `<persisted-output>` 标记。
2. **消息窗口裁剪**：历史消息过多时保留开头和尾部窗口。
3. **最近工具结果保留**：只保留最近若干条完整工具结果。
4. **当前会话记忆替换**：在可用时使用 `conversation_memories` 作为长期对话主干信息。
5. **全量历史摘要**：对历史整体做摘要，写入 `context_summaries`。
6. **Reactive Compact**：遇到 413 / `prompt_too_long` 时触发更激进的压缩恢复。

落盘工具结果可由模型通过 `read_persisted_output` 分块读取。压缩后的模型历史可写入 `conversation_model_histories`，用于后续轮次复用。

相关代码：

- `internal/web/runtime/context_compression.go`
- `internal/web/runtime/compression/`
- `internal/web/storage/persisted_outputs_repo.go`
- `internal/web/storage/conversation_model_histories_repo.go`

## 记忆系统

记忆系统由全局开关和用户开关共同控制：runtime 默认开启，用户可通过 `/api/me/memory` 切换自己的 `memory_enabled`。

长期记忆分三类：

| 类型 | 含义 | 作用域 |
| --- | --- | --- |
| `episodic_memory` | 与用户相关的具体事件、经历或会话摘要 | 当前用户 |
| `user_preference` | 用户稳定偏好、约束或习惯 | 当前用户 |
| `semantic` | 去身份化、可复用的通用知识 | 系统级，跨用户候选 |

流程：

1. 对话开始前，从 `memories` 中读取用户相关记忆和系统级语义记忆。
2. 调用 LLM 从候选记忆中选择最多 10 条相关记忆。
3. 将选中的记忆渲染为 `<memory>` 段注入 system prompt。
4. 对话结束后异步抽取新记忆，并按类型做合并、裁剪或去重。
5. 同步维护 `conversation_memories`，作为当前会话主干信息；它不直接注入 prompt，主要用于上下文压缩替代。

相关代码：

- `internal/web/runtime/memory.go`
- `internal/web/runtime/conversation_memory.go`
- `internal/web/runtime/compression/conversation_memory_strategy.go`
- `internal/web/storage/memories_repo.go`
- `internal/web/storage/conversation_memories_repo.go`

## System Prompt

最终 system prompt 由 `assistant.BuildSystemPrompt` 拼装，包含：

- `<identity>`：从 `system_prompt_path` 指向的文件读取的基础人设。
- `<workspace>`：当前 surface 与 working directory。
- `<tools>`：本轮可用工具名；包含 `read_persisted_output` 时会追加落盘产物使用说明。
- `<skills>`：当前 Skill snapshot 的技能摘要与 `load_skill` 使用规则。
- `<memory>`：本轮选择出的相关长期记忆。

为空的动态段不会写入最终 prompt。修改 `system_prompt.md` 只会替换 `<identity>` 段，其余段落由 runtime 按当前用户、工具、Skill 与记忆动态生成。

相关代码：

- `internal/assistant/prompt.go`
- `internal/config/paths.go`
- `internal/web/runtime/prompt_builder.go`

## Agent Loop 与 Hook

每轮 Web 对话的主流程：

1. 获取会话锁，避免同一会话并发写入。
2. 读取会话历史，触发 `UserPromptSubmit` Hook。
3. 构建当前用户 Skill snapshot。
4. 选择相关记忆并构建 system prompt。
5. 合并内置工具与 MCP 工具定义。
6. 请求模型前执行上下文压缩。
7. 调用 OpenAI 兼容模型并流式下发正文和推理增量。
8. 如模型发起工具调用，触发 `PreToolUse` → 执行工具 → `PostToolUse`。
9. 模型停止工具调用后，触发 `Stop` Hook，持久化最终 assistant 消息并通过 SSE 下发。
10. 对话结束后异步维护模型历史和记忆。

默认 Hook：

- `UserPromptSubmit`：追加用户消息、推断或刷新会话标题/活动时间。
- `PreToolUse`：预解析工具审计信息，例如 cwd、命令路径、产物路径与来源。
- `PostToolUse`：记录执行结果或拒绝原因，持久化工具调用，追加 tool message 给下一轮模型。
- `Stop`：持久化 assistant 消息，下发最终 assistant 事件，包含 `content`、`reasoning_content` 与 meta。

相关代码：

- `internal/web/runtime/conversation_flow.go`
- `internal/web/runtime/hooks/manager.go`
- `internal/web/runtime/hooks/user_prompt.go`
- `internal/web/runtime/hooks/tool.go`
- `internal/web/runtime/hooks/stop.go`
- `internal/web/runtime/hook_bridge.go`

## API 简表

### 公开接口

- `GET /api/health`
- `POST /api/auth/register`
- `POST /api/auth/login`

### 需要鉴权的接口

- `POST /api/auth/logout`
- `GET /api/me`
- `PATCH /api/me/memory`：更新当前用户记忆开关。

Skill：

- `GET /api/skills`
- `POST /api/skills`
- `GET /api/skills/{id}`
- `PUT /api/skills/{id}`
- `PATCH /api/skills/{id}`：仅更新状态。
- `DELETE /api/skills/{id}`

MCP：

- `GET /api/mcp-servers`
- `POST /api/mcp-servers`
- `GET /api/mcp-servers/{id}`
- `PUT /api/mcp-servers/{id}`
- `PATCH /api/mcp-servers/{id}`：启用/禁用。
- `DELETE /api/mcp-servers/{id}`
- `POST /api/mcp-servers/{id}/test`：临时连接并返回工具名。

会话：

- `GET /api/conversations`
- `POST /api/conversations`
- `GET /api/conversations/{id}`
- `PATCH /api/conversations/{id}`：改标题。
- `DELETE /api/conversations/{id}`
- `POST /api/conversations/{id}/stream`：SSE 流式对话。

### SSE 事件

`POST /api/conversations/{id}/stream` 返回 `text/event-stream`。当前前端主要消费：

| event | 说明 |
| --- | --- |
| `assistant_delta` | assistant 正文增量 |
| `reasoning_delta` | 推理过程增量 |
| `meta` | 工具调用次数、上下文 token / budget 等元信息 |
| `assistant` | 最终 assistant 消息 |
| `error` | 错误信息 |
| `done` | 流式响应结束 |

后端还会发送 `conversation` 事件；当前前端不依赖该事件。当前没有独立工具调用明细 SSE 事件，前端主要通过 `meta.tool_call_count` 展示工具调用次数。

相关代码：`internal/web/app/routes.go`、`internal/web/app/*_handlers.go`、`../web/src/api.ts`。

## 部署

构建后端发布包：

```bash
cd go-agent
./build.sh
```

`build.sh` 会清空并重建 `output/`，编译主程序到 `output/bin/go-agent`，编译根目录可选 `cmd/` 下各子命令到 `output/bin/`，并复制 `config.json`、`system_prompt.md`、`skills/`、`cmd/` 资源（源不存在则跳过）。产物结构：

```text
output/
├── bin/                 # go-agent 主二进制（+ 各 cmd 子命令二进制，若有）
├── cmd/                 # cmd 源码副本（若有）
├── skills/              # 内置 Skill 副本
├── logs/                # 空目录
├── workspace/           # 空目录
├── config.json
└── system_prompt.md
```

云端启动示例：

```bash
cd /srv/go-agent/output    # 包含 config.json、bin/、skills/、system_prompt.md 等

export OPENAI_API_KEY=your-api-key
export MYSQL_PASSWORD=your-mysql-password
export REDIS_PASSWORD=your-redis-password
export ES_PASSWORD=your-es-password
export JWT_SECRET=replace-with-production-secret

./bin/go-agent
```

如果需要内置 MCP，请把 `mcp_config.json` 放到 `app_home` 下，例如 `output/mcp_config.json`。

前端需要单独构建与发布：

```bash
cd ../web
npm install
npm run build
```

`go-agent/build.sh` 不会构建或复制 `web/dist`；当前 Go 后端也未托管前端静态资源。请使用 Nginx、CDN 或其他静态资源服务发布 `web/dist`。

## 开发验证

后端：

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

- 运行前确保 MySQL、Redis、Elasticsearch、LLM 服务可用。
- `config.json` 固定从进程当前目录读取；启动目录需包含该文件。
- `system_prompt.md` 当前必须存在；部署时需与 `config.json` 保持路径一致。
- `OPENAI_API_KEY`、数据库密码、Redis 密码、ES 密码、`JWT_SECRET` 建议全部通过环境变量提供。
- 日志默认写入 `app_home/logs`，部署时需确保可写。
- 前端跨域失败时检查后端 `allowed_origin` 是否等于前端地址。
- 登录后仍返回 401 时检查 Cookie 策略、域名、SameSite 与 `credentials: include`。
- 浏览器聊天不是用户本地终端代理；工具运行在服务端 `workspace_root`。
- 生产环境不建议开启 `bash_allow_outside_workspace` 或 `bash_allow_dangerous_commands`。
- 如果使用内置 MCP，记得单独部署 `app_home/mcp_config.json`；当前构建脚本不会自动复制该文件。
