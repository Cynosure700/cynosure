# Cynosure (nano_cc) Go → Python 迁移设计方案

> **项目名称**: Cynosure / nano_cc  
> **当前语言**: Go 1.26.1（~16,727 行 Go 代码，~21,000+ 总计）  
> **目标语言**: Python 3.12+  
> **文档版本**: v1.0  
> **日期**: 2026-02

---

## 1. 概述与迁移目标

### 1.1 为什么迁移到 Python

| 因素 | 说明 |
|------|------|
| **生态丰富度** | Python 拥有更成熟的 AI/LLM 生态（OpenAI SDK、LangChain、LlamaIndex）、MCP SDK 和 TUI 框架（Textual、Rich） |
| **社区与协作** | 降低贡献门槛，吸引更多 AI Agent 开发者参与 |
| **快速迭代** | Python 的动态特性使原型开发和功能实验更敏捷 |
| **依赖管理** | Go 在 LLM/MCP 领域的 SDK 成熟度不如 Python |

### 1.2 迁移原则

1. **架构对等** — 保持 Go 版本的分层架构和核心数据流不变
2. **功能完整** — 覆盖所有内置工具、技能系统、MCP 集成、会话持久化
3. **渐进替换** — 按模块独立迁移，避免大爆炸式重写
4. **行为兼容** — 相同输入产生相同输出，工具接口语义不变
5. **能力不降级** — TUI 体验、安全性、上下文压缩等核心特性不弱于 Go 版本

---

## 2. 整体架构

### 2.1 Python 项目结构

```
cynosure/
├── pyproject.toml              # 项目元数据与依赖声明
├── README.md
├── cynosure/
│   ├── __init__.py
│   ├── __main__.py             # python -m cynosure 入口
│   │
│   ├── cli/                    # CLI 入口与参数解析
│   │   ├── __init__.py
│   │   └── root.py             # ArgumentParser + 启动逻辑
│   │
│   ├── config/                 # 配置管理
│   │   ├── __init__.py
│   │   ├── settings.py         # settings.json 加载
│   │   └── config_model.py     # config.json 加载与安全边界
│   │
│   ├── idgen/                  # ID 生成器
│   │   ├── __init__.py
│   │   └── generator.py        # nanoid / ulid 实现
│   │
│   ├── logger/                 # 日志系统
│   │   ├── __init__.py
│   │   └── logger.py
│   │
│   ├── llm/                    # LLM 客户端抽象
│   │   ├── __init__.py
│   │   ├── client.py           # OpenAI 兼容 API 客户端
│   │   └── types.py            # 消息/工具调用类型定义
│   │
│   ├── textutil/               # 文本工具
│   │   ├── __init__.py
│   │   └── utils.py
│   │
│   ├── safety/                 # 安全边界
│   │   ├── __init__.py
│   │   ├── guardrails.py       # 路径验证、命令白名单
│   │   └── redaction.py        # 敏感信息脱敏
│   │
│   ├── tools/                  # 内置工具 (对应 Go internal/tools)
│   │   ├── __init__.py
│   │   ├── base.py             # 工具基类与注册器
│   │   ├── bash.py             # bash 执行
│   │   ├── file_ops.py         # read_file, write_file, edit_file
│   │   ├── search.py           # grep, glob
│   │   ├── web.py              # web_fetch
│   │   ├── todo.py             # todo_write
│   │   ├── subagent.py         # spawn_subagent
│   │   ├── skill_loader.py     # load_skill
│   │   ├── runtime_env.py      # RuntimeEnv 上下文
│   │   └── path_guard.py       # 路径安全检查
│   │
│   ├── sessions/               # 会话与技能管理
│   │   ├── __init__.py
│   │   ├── session.py          # 会话加载/保存
│   │   ├── skill.py            # 技能发现与加载
│   │   └── skill_creator.py    # 技能创建/编辑/评估
│   │
│   ├── assistant/              # System Prompt 构建
│   │   ├── __init__.py
│   │   └── prompt_builder.py   # 系统提示词组装
│   │
│   ├── agent/                  # 核心 Agent 运行时
│   │   ├── __init__.py
│   │   ├── runtime/
│   │   │   ├── __init__.py
│   │   │   ├── service.py      # 服务主入口
│   │   │   ├── handler.py      # 主循环 (用户输入→LLM→工具→存储)
│   │   │   ├── memory.py       # 长期记忆管理
│   │   │   ├── summarizer.py   # 上下文摘要
│   │   │   ├── tool_manager.py # 工具注册与分发
│   │   │   ├── mcp_provider.py # MCP 工具桥接
│   │   │   ├── approval.py     # 用户审批流程
│   │   │   └── todo_reminder.py# 待办提醒
│   │   │   ├── compression/
│   │   │   │   ├── __init__.py
│   │   │   │   ├── compressor.py       # 压缩策略基类
│   │   │   │   ├── window.py           # 窗口裁剪
│   │   │   │   ├── summary.py          # 摘要压缩
│   │   │   │   ├── reactive.py         # 413 激进压缩
│   │   │   │   └── helpers.py          # 压缩工具函数
│   │   │   └── hooks/
│   │   │       ├── __init__.py
│   │   │       ├── types.py            # Hook 类型定义
│   │   │       ├── tool_finish.py      # 工具调用后处理
│   │   │       ├── tool_persist.py     # 工具结果持久化
│   │   │       ├── compression.py      # 上下文压缩钩子
│   │   │       └── model_interceptor.py# 模型调用拦截
│   │   ├── mcp/
│   │   │   ├── __init__.py
│   │   │   ├── client.py       # MCP 客户端
│   │   │   ├── discovery.py    # 自动发现 MCP 服务
│   │   │   ├── transport.py    # stdio / sse / streamable 传输
│   │   │   ├── tools.py        # MCP 工具适配
│   │   │   ├── config.py       # .mcp.json 解析
│   │   │   ├── registry.py     # MCP 工具注册表
│   │   │   └── types.py        # MCP 类型定义
│   │   └── storage/
│   │       ├── __init__.py
│   │       ├── models.py       # 数据模型 (Conversation, Message, ToolCall)
│   │       └── history.py      # 会话持久化与恢复
│   │
│   ├── tui/                    # 终端界面
│   │   ├── __init__.py
│   │   ├── app.py              # Textual 应用主入口
│   │   ├── widgets/            # 自定义组件
│   │   │   ├── __init__.py
│   │   │   ├── chat_panel.py   # 对话面板 (Markdown 渲染)
│   │   │   ├── input_bar.py    # 输入栏
│   │   │   ├── status_bar.py   # 状态栏
│   │   │   ├── tool_output.py  # 工具结果面板
│   │   │   └── todo_panel.py   # 待办面板
│   │   └── styles.py           # 主题与样式 (Rich)
│   │
│   └── local/                  # 启动引导
│       ├── __init__.py
│       └── bootstrap.py        # 初始化目录、配置、技能发现
│
├── tests/                      # 测试套件
│   ├── __init__.py
│   ├── test_tools/
│   ├── test_agent/
│   ├── test_mcp/
│   ├── test_safety/
│   └── conftest.py
│
└── .cynosure/                  # 用户配置目录 (运行时自动创建)
    ├── settings.json
    ├── config.json
    └── skills/
```

### 2.2 核心数据流（保持不变）

```
用户输入
    │
    ▼
TUI (Textual) ───→ Handler 主循环
                        │
                        ▼
                Prompt Builder ───→ 组装 System Prompt + 记忆
                        │
                        ▼
                LLM Client ───→ OpenAI 兼容 API
                        │
                        ▼
                解析响应 (文本 / Tool Calls)
                        │
                ┌───────┴───────┐
                ▼               ▼
           文本回复         工具执行
                              │
                              ▼
                        Tool Manager
                        ├─ 内置工具 (bash, read_file, ...)
                        └─ MCP Provider → MCP 工具
                              │
                              ▼
                        Hooks 处理
                        ├─ tool_finish (审计日志)
                        ├─ tool_persist (大结果落盘)
                        └─ compression (上下文压缩)
                              │
                              ▼
                        Storage ───→ 持久化到 ~/.cynosure/session/
                              │
                              ▼
                        Memory 更新
                              │
                              ▼
                        循环 (回到 Handler)
```

---

## 3. 关键模块迁移方案

### 3.1 依赖映射

| Go 依赖 | Python 替代 | 说明 |
|---------|------------|------|
| `github.com/sashabaranov/go-openai` | `openai` (官方 SDK) | OpenAI 兼容 API 客户端 |
| `charmbracelet/bubbletea` | `textual` | TUI 框架 |
| `charmbracelet/lipgloss` | `rich` | 终端样式 |
| `charmbracelet/glamour` | `rich.markdown` | Markdown 渲染 |
| `MCP SDK` | `mcp` (Python SDK) | MCP 协议实现 |
| `JSON Schema` | `jsonschema` | JSON Schema 验证 |
| `go-nanoid` | `nanoid` | ID 生成 |
| Go `encoding/json` | `json` (stdlib) | JSON 序列化 |
| Go `net/http` | `httpx` / `aiohttp` | HTTP 请求 |
| Go `os/exec` | `subprocess` | Shell 命令执行 |

### 3.2 模块迁移优先级与阶段

```
阶段 1 (核心基础设施)     阶段 2 (Agent 运行时)     阶段 3 (TUI 与集成)
┌──────────────────┐     ┌──────────────────┐     ┌──────────────────┐
│  config           │  →  │  agent/runtime/  │  →  │  tui/            │
│  logger           │     │    handler       │     │  (Textual)       │
│  idgen            │     │    service       │     │                  │
│  textutil         │     │    memory        │     │  sessions/skill  │
│  safety           │     │    summarizer    │     │  bootstrap       │
│  llm/client       │     │    tool_manager  │     │                  │
│  tools/ (基础)     │     │    mcp_provider  │     │  agent/mcp/      │
│  tools/ (文件)     │     │    approval      │     │  agent/storage/  │
└──────────────────┘     │    compression   │     └──────────────────┘
                          │    hooks          │
                          │  agent/mcp/       │
                          │  agent/storage/   │
                          └──────────────────┘
```

### 3.3 模块详情

#### 3.3.1 `config/` — 配置管理

- 加载 `~/.cynosure/settings.json`（必需：`open_auth_token`, `open_model`, `open_base_url`）
- 加载 `~/.cynosure/config.json`（可选：安全边界配置）
- 使用 `pydantic` 进行 schema 验证

```python
# cynosure/config/settings.py
from pydantic import BaseModel

class Settings(BaseModel):
    open_auth_token: str
    open_model: str
    open_base_url: str

class SecurityConfig(BaseModel):
    bash_allow_outside_workspace: bool = False
    bash_allow_dangerous_commands: list[str] = []
```

#### 3.3.2 `llm/` — LLM 客户端

- 封装 `openai.OpenAI` 同步/异步客户端
- 支持 streaming（Textual 需要流式响应更新）
- 错误处理：重试逻辑、413 检测、超时

```python
# cynosure/llm/client.py
import openai
from openai.types.chat import ChatCompletionMessage

class LLMClient:
    def __init__(self, settings: Settings):
        self.client = openai.OpenAI(
            api_key=settings.open_auth_token,
            base_url=settings.open_base_url,
        )
        self.model = settings.open_model

    def chat_completion(
        self,
        messages: list[dict],
        tools: list[dict] | None = None,
        stream: bool = False,
    ) -> ChatCompletionMessage:
        ...
```

#### 3.3.3 `tools/` — 内置工具

Go 版本有 8 个内置工具 + `spawn_subagent`，全部迁移：

| 工具 | Python 实现要点 |
|------|----------------|
| `load_skill` | 扫描两级 skill 目录，加载 SKILL.md |
| `bash` | `subprocess.run()` + 路径/命令安全检查 |
| `read_file` | `pathlib.Path.read_text()` + 路径限制 |
| `write_file` | `pathlib.Path.write_text()` + 审批 |
| `edit_file` | `difflib` 精确匹配替换 |
| `todo_write` | 内存中任务列表，支持 CRUD |
| `spawn_subagent` | 子进程或 asyncio 子任务隔离执行 |
| `read_persisted_output` | 从 `~/.cynosure/task_outputs/` 分块读取 |

**工具注册器** — Python 版本引入装饰器式注册：

```python
# cynosure/tools/base.py
from dataclasses import dataclass
from typing import Any, Callable

@dataclass
class ToolDef:
    name: str
    description: str
    parameters: dict
    handler: Callable

class ToolRegistry:
    _tools: dict[str, ToolDef] = {}

    @classmethod
    def register(cls, name: str, description: str, parameters: dict):
        def decorator(fn):
            cls._tools[name] = ToolDef(name, description, parameters, fn)
            return fn
        return decorator

    @classmethod
    def definitions(cls) -> list[dict]:
        return [
            {
                "type": "function",
                "function": {
                    "name": t.name,
                    "description": t.description,
                    "parameters": t.parameters,
                }
            }
            for t in cls._tools.values()
        ]
```

#### 3.3.4 `safety/` — 安全边界

- **路径守卫**：使用 `pathlib.Path.resolve()` 验证路径不出工作区
- **命令白名单**：危险命令（`rm`, `dd`, `mkfs` 等）需白名单允许
- **脱敏**：检测并替换 `settings.json` 中的 token/key
- **审计日志**：记录每次工具调用的完整上下文

#### 3.3.5 `agent/runtime/` — Agent 运行时核心

这是迁移中最关键的模块，包含主循环逻辑。

**Handler 主循环 (handler.py)**：

```python
async def run_loop(service: "Service", user_input: str):
    # 1. 构建 System Prompt (含记忆、技能、工具定义)
    # 2. 调用 LLM
    # 3. 解析响应
    #   a. 文本 → 输出到 TUI
    #   b. ToolCalls → 执行工具
    # 4. 执行 Pre/Post Tool Hooks
    # 5. 大结果落盘 (persisted-output)
    # 6. 存储消息到 conversation_history
    # 7. 更新长期记忆
    # 8. 压缩上下文
    # 9. 循环
```

**上下文压缩 (compression/)**：

| 策略 | 说明 | Python 实现 |
|------|------|------------|
| 窗口裁剪 | 保留最近 N 轮 | `deque(maxlen=N)` |
| 大结果落盘 | >200KB 写入文件，替换为标记 | `aiofiles` + 标记替换 |
| 摘要压缩 | 对早期上下文做 LLM 摘要 | 调用 LLM 生成 |
| 413 激进压缩 | 逐层裁剪直到通过 | 递归/循环降级策略 |

**Hooks 系统 (hooks/)**：

Python 版本使用类装饰器模式替代 Go 的函数切片：

```python
from typing import Protocol

class UserPromptSubmitHook(Protocol):
    async def __call__(self, ctx: "UserPromptSubmitContext") -> None: ...

class PreToolUseHook(Protocol):
    async def __call__(self, ctx: "ToolUseContext") -> None: ...

class PostToolUseHook(Protocol):
    async def __call__(self, ctx: "ToolUseContext") -> None: ...

class StopHook(Protocol):
    async def __call__(self, ctx: "StopContext") -> None: ...

class HookManager:
    user_prompt_submit: list[UserPromptSubmitHook] = []
    pre_tool_use: list[PreToolUseHook] = []
    post_tool_use: list[PostToolUseHook] = []
    stop: list[StopHook] = []
```

#### 3.3.6 `agent/mcp/` — MCP 集成

- **自动发现**：扫描工作区 `.cynosure/.mcp.json`，支持 `mcp_servers` (数组) 和 `mcpServers` (map) 格式
- **三种传输方式**：stdio / sse / streamable
- **MCP 工具命名**：`mcp__{server}__{tool}` 格式注入到模型工具列表

Python MCP SDK (`mcp`) 已原生支持这些传输方式，迁移工作量较小。

#### 3.3.7 `agent/storage/` — 会话持久化

- 路径：`~/.cynosure/session/{session_id}/`
- 文件：`history` (完整展示历史) + `model_history` (压缩后的模型上下文)
- 恢复：`/resume` 命令读取 history 重建上下文
- 使用 `json` 或 `msgpack` 序列化

#### 3.3.8 `tui/` — 终端用户界面

**技术选型：Textual**

Textual 是 Python 最成熟的 TUI 框架，提供：

| 能力 | Textual 对应 |
|------|-------------|
| 组件树 | `Widget` 继承体系 |
| 布局 | CSS 布局引擎 |
| Markdown 渲染 | `RichMarkdown` widget |
| 异步事件循环 | asyncio 原生 |
| 键盘快捷键 | `@on(Key)` 绑定 |
| 主题 | CSS 变量 + 动态主题 |

**界面布局**：

```
┌─────────────────────────────────────────────┐
│  Status Bar         会话: xxx  模型: xxx    │
├──────────────────────┬──────────────────────┤
│                      │                      │
│  对话面板 (Chat)      │  工具输出面板        │
│  · Markdown 渲染     │  · 实时工具结果      │
│  · 用户/助手消息      │  · 代码块语法高亮   │
│  · 代码块折叠         │  · 文件 diff 展示    │
│                      │                      │
│                      │                      │
├──────────────────────┴──────────────────────┤
│  输入栏 (Input Bar)    [Ctrl+Enter 发送]    │
├─────────────────────────────────────────────┤
│  待办面板 (Todo)  |  帮助命令               │
└─────────────────────────────────────────────┘
```

**快捷键映射**：

| 快捷键 | 功能 |
|--------|------|
| `Ctrl+C` | 中断当前响应 / 退出 |
| `Ctrl+L` | 清屏 |
| `Ctrl+S` | 保存会话 |
| `Ctrl+[/Esc` | 关闭侧面板 |
| `Tab` | 焦点切换 |
| `Ctrl+P` | 打开待办面板 |
| `/resume` | 恢复历史会话命令 |
| `/help` | 帮助命令 |

---

## 4. 迁移策略与执行计划

### 4.1 阶段划分

| 阶段 | 模块 | 预估工时 | 里程碑 |
|------|------|---------|--------|
| **P0 基础设施** | config, logger, idgen, textutil, safety, llm/client, tools/base | 1 周 | 工具框架可运行 |
| **P1 核心工具** | tools/bash, file_ops, search, web, todo, path_guard, skill_loader | 1 周 | 内置工具全可用 |
| **P2 Agent 运行时** | agent/runtime/service, handler, tool_manager, mcp_provider, approval | 2 周 | 基本对话循环 |
| **P3 高级运行时** | memory, summarizer, compression, hooks, todo_reminder | 1.5 周 | 完整 Agent 能力 |
| **P4 存储与 MCP** | agent/storage, agent/mcp (所有文件) | 1 周 | 会话持久化 + MCP |
| **P5 TUI 界面** | tui/ 全部 (app, widgets, styles) | 2 周 | 可用图形界面 |
| **P6 集成与测试** | bootstrap, cli, sessions/skill, 集成测试 | 1 周 | 端到端可用 |
| **P7 打磨** | 性能优化、错误处理、文档、技能创建工具 | 1 周 | 生产级质量 |

**总计预估：约 10.5 周（全职）**

### 4.2 并行策略

```
周 1-2:     P0 ──→ P1
              ↘
周 3-4:         P2 ──→ P3
                         ↘
周 5-6:              P4 ──→ P5
                                ↘
周 7-8:                     P6 ──→ P7
```

- P0–P1 可单人完成
- P2–P3 和 P4 可并行（Agent 运行时 + 存储/MCP 由不同人负责）
- P5 TUI 可在 P2 基本就绪后立即启动（用 mock 数据先行开发界面）
- P6集成测试持续进行

### 4.3 风险与缓解

| 风险 | 影响 | 缓解措施 |
|------|------|---------|
| Python asyncio 与 Go goroutine 并发模型差异 | 工具执行/LLM 调用的并发行为可能不同 | 用 `asyncio.gather` 替代 goroutine + channel；充分测试并发场景 |
| Textual 渲染性能不如 Bubble Tea | 大量日志/长文本时卡顿 | 使用虚拟化渲染（`VirtualDocument`）、分页加载、懒加载 |
| MCP SDK 版本兼容性 | 与目标 MCP 服务器不兼容 | 锁定 `mcp>=1.0.0` 版本；保留 stdio 回退 |
| 上下文压缩行为差异 | 413 处理效果不同 | 压缩策略单独测试；保留 Go 版压缩结果作为 benchmark |
| subprocess 安全限制 | bash 工具行为差异 | `shlex.quote()` + 白名单 + 超时控制 |

### 4.4 Python 独有的增强机会

迁移过程中可顺势引入的改进：

1. **插件系统** — 利用 Python 的 `importlib.metadata` 实现 entry point 插件
2. **异步工具执行** — 工具可声明 `async def` 实现并行执行
3. **更好的 streaming** — Textual 原生支持异步流式更新，对话体验更流畅
4. **类型安全** — 全项目 `pyright` / `mypy` 严格模式
5. **测试覆盖率** — `pytest` + `pytest-asyncio` + `pytest-cov` 确保 >80%
6. **包发布** — `pyproject.toml` 支持 `pip install cynosure` 一键安装

---

## 5. 关键设计决策

### 5.1 同步 vs 异步

- **全面异步** — 所有 I/O 操作（LLM 调用、文件读写、子进程、MCP 通信）均使用 `asyncio`
- TUI 框架 Textual 原生基于 `asyncio`
- 工具执行：同步命令（bash）用 `asyncio.create_subprocess_exec` 包装，异步命令直接 `await`

### 5.2 数据模型

使用 `pydantic` 替代 Go struct，确保序列化/反序列化兼容：

```python
from pydantic import BaseModel
from typing import Optional

class MessageToolCall(BaseModel):
    id: str
    type: str = "function"
    function: MessageFunctionCall

class Message(BaseModel):
    id: str
    role: str  # system, user, assistant, tool
    content: str
    reasoning_content: Optional[str] = None
    tool_call_id: Optional[str] = None
    tool_calls: list[MessageToolCall] = []
    created_at: str  # ISO 8601

class Conversation(BaseModel):
    id: str
    title: str
    created_at: str
    updated_at: str
    messages: list[Message] = []
```

### 5.3 配置文件兼容

Python 版本保持与 Go 版本完全相同的 JSON schema：

- `~/.cynosure/settings.json` — 字段名、结构不变
- `~/.cynosure/config.json` — 字段名、结构不变
- `.cynosure/.mcp.json` — 兼容 `mcp_servers` (数组) 和 `mcpServers` (map) 两种格式

### 5.4 会话文件兼容

- `history` 文件使用 JSON Lines 格式（与 Go 版本兼容）
- 已存在的 Go 版本会话数据可被 Python 版本读取
- `/resume` 命令恢复逻辑相同

### 5.5 技能系统兼容

- 两级目录：`~/.cynosure/skills/{name}/skill.md` 和 `<cwd>/.cynosure/skills/{name}/SKILL.md`
- 文件名大小写区分（用户级小写，工作区级大写）保持不变
- 摘要预注入 + `load_skill` 按需加载机制不变

---

## 6. 测试策略

### 6.1 测试金字塔

```
     /\
    /  \        E2E 集成测试 (5%)
   /    \
  /______\      模块集成测试 (15%)
 /        \
/__________\   单元测试 (80%)
```

### 6.2 测试工具

| 用途 | 工具 |
|------|------|
| 单元测试 | `pytest` + `pytest-asyncio` |
| Mock | `unittest.mock` / `pytest-mock` |
| LLM Mock | 自定义 mock server 或 `responses` |
| 覆盖率 | `pytest-cov` (目标 >80%) |
| 代码质量 | `ruff` (lint) + `mypy` / `pyright` (type) |

### 6.3 关键测试场景

- 工具执行边界（路径越权、危险命令、超时）
- 上下文压缩（窗口裁剪、大结果落盘、413 激进压缩）
- 会话持久化与恢复（包括 Go 版本遗留数据）
- MCP 自动发现与连接（stdio/sse/streamable）
- TUI 交互（快捷键、面板切换、Markdown 渲染）

---

## 7. 性能考量

| 关注点 | 方案 |
|--------|------|
| LLM 调用延迟 | 异步非阻塞 + 可选 streaming |
| 大文件读取 | 分块读取 + persisted-output 机制 |
| 上下文压缩 | 增量压缩而非全量压缩 |
| 会话加载 | 懒加载消息历史（分页） |
| TUI 渲染 | 虚拟化列表（`Textual VirtualScroll`） |
| 内存占用 | 大文本结果及时落盘，内存中只保留引用 |

---

## 8. 附录

### 8.1 Go → Python API 对照

| Go 概念 | Python 对应 |
|---------|------------|
| `struct` | `pydantic.BaseModel` / `dataclass` |
| `interface` | `Protocol` / `ABC` |
| `goroutine + channel` | `asyncio.Task + asyncio.Queue` |
| `sync.Mutex` | `asyncio.Lock` |
| `context.Context` | `contextvars.ContextVar` / 显式传递 |
| `defer` | `try/finally` / `contextlib.ExitStack` |
| `go:embed` | `importlib.resources` |
| `json.RawMessage` | `dict` / `Any` |
| `io.Reader/Writer` | `IO[bytes]` / `IO[str]` |
| `net/http` | `httpx` / `aiohttp` |
| `os/exec` | `asyncio.create_subprocess_exec` |

### 8.2 依赖清单 (pyproject.toml)

```toml
[project]
name = "cynosure"
version = "0.1.0"
requires-python = ">=3.12"
dependencies = [
    "openai>=1.0.0",
    "textual>=1.0.0",
    "rich>=13.0.0",
    "mcp>=1.0.0",
    "pydantic>=2.0.0",
    "httpx>=0.27.0",
    "aiofiles>=24.0.0",
    "nanoid>=2.0.0",
    "jsonschema>=4.0.0",
    "pyyaml>=6.0",
]

[project.optional-dependencies]
dev = [
    "pytest>=8.0",
    "pytest-asyncio>=0.24",
    "pytest-cov>=5.0",
    "ruff>=0.4.0",
    "mypy>=1.10",
    "pyright>=1.1.360",
]

[tool.ruff]
line-length = 100
target-version = "py312"

[tool.mypy]
strict = true
python_version = "3.12"
```

---

> **文档维护者**: Cynosure 团队  
> **下一步**: 启动阶段 P0，从 `config/`、`logger/`、`safety/` 和工具框架开始迁移。
