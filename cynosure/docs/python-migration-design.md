# Cynosure (nano_cc) Python 迁移设计文档

> **项目**：Cynosure — 纯本地、终端原生的通用 AI Agent（原名 nano_cc）
> **当前实现**：Go 1.26.1，约 21,000 行（16,727 行 Go + 4,503 行 Markdown 文档）
> **目标**：完整迁移至 Python，保留全部功能与架构设计

---

## 1. 迁移动机

- **生态优势**：Python 拥有更丰富的 AI/ML 生态（LangChain、LlamaIndex、Transformers），且在 MCP SDK、JSON Schema 等方面有原生支持
- **维护便利**：团队成员 Python 经验更丰富，降低维护门槛
- **快速迭代**：Python 无编译等待，开发反馈周期短
- **社区集成**：大量第三方工具和库可直接复用

---

## 2. 总体架构（保持原样）

迁移后整体架构不变，仅语言从 Go 切换为 Python：

```
cynosure/
├── pyproject.toml              # Python 项目配置（Poetry/uv）
├── src/cynosure/
│   ├── __init__.py
│   ├── __main__.py             # 入口（替换 main.go）
│   ├── agent/                  # 核心 Agent 包
│   │   ├── runtime/            # 运行时（主循环、服务、工具注册）
│   │   │   ├── compression/    # 上下文压缩策略
│   │   │   └── hooks/          # 插件式钩子
│   │   ├── mcp/                # MCP 集成
│   │   └── storage/            # 会话存储
│   ├── tui/                    # 终端 UI（替换 Bubble Tea → Textual）
│   ├── tools/                  # 工具定义与调度
│   ├── config/                 # 配置管理
│   ├── llm/                    # LLM 客户端封装
│   ├── sessions/               # 会话管理
│   └── ...                     # 其他辅助包
└── docs/                       # 设计文档
```

---

## 3. 核心依赖映射

| 用途 | Go 依赖 | Python 替代 |
|------|---------|-------------|
| LLM 客户端 | `sashabaranov/go-openai` | `openai` Python SDK |
| TUI 框架 | `charmbracelet/bubbletea` | **Textual**（最成熟、事件驱动） |
| 终端样式 | `charmbracelet/lipgloss` | Textual 内置 CSS 样式 |
| Markdown 渲染 | `charmbracelet/glamour` | Textual `Markdown` widget 或 `rich.markdown` |
| MCP 协议 | 自定义实现（~8 个文件） | `mcp` Python SDK（pip install mcp） |
| JSON Schema 验证 | 自定义 | `jsonschema` 或 `pydantic` |
| 日志 | 自定义 logger | `loguru` 或标准库 `logging` |

---

## 4. 分层迁移方案

### 第 1 层：基础库（无外部依赖）

**覆盖包**：`config`, `idgen`, `logger`, `textutil`, `safety`

这些包基本是纯逻辑、无外部依赖，直接逐文件翻译即可。
- `config` → Pydantic models + `json` 文件读取
- `idgen` → `uuid4()` 生成
- `logger` → `loguru` 或标准 `logging`
- `textutil` → Python `re` / `textwrap` 标准库
- `safety` → 目录白名单、路径检查（`pathlib`）

### 第 2 层：LLM 客户端与存储

**覆盖包**：`llm`, `agent/storage`

#### `llm` 包
Go 中 `sashabaranov/go-openai` 的使用模式：
```go
client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
    Model:    modelID,
    Messages: messages,
    Tools:    tools,
})
```

Python 等价：
```python
from openai import OpenAI
client = OpenAI()
response = client.chat.completions.create(
    model=model_id,
    messages=messages,
    tools=tools,
)
```

需要封装一层 `LLMClient`，支持 OpenAI 兼容 API（如 Ollama、vLLM）：
```python
class LLMClient:
    def __init__(self, base_url: str, api_key: str, model_id: str):
        self.client = OpenAI(base_url=base_url, api_key=api_key)
        self.model_id = model_id

    def chat(self, messages, tools=None, **kwargs) -> ChatResponse:
        ...
```

#### `agent/storage` 包
Go 中的会话持久化（`~/.cynosure/session/<id>/`）：
- 存储格式：JSON 文件
- 关键接口：`ConversationHistory.{Append, Load, Save}`

Python 实现：
```python
# 用 Pydantic 定义消息模型
class Message(BaseModel):
    role: str
    content: str
    tool_calls: list[ToolCall] | None = None
    ...

class ConversationHistory:
    def __init__(self, session_dir: Path):
        self.session_dir = session_dir

    def append(self, message: Message): ...
    def load(self) -> list[Message]: ...
    def save(self): ...
```

### 第 3 层：工具系统

**覆盖包**：`tools`

Go 中通过 `openai.Tool` 定义工具，通过 `Dispatch` 函数调度。

Python 实现方案：

```python
# tools/definitions.py
from pydantic import BaseModel
from openai.types.chat import ChatCompletionToolParam

class ToolDefinition(BaseModel):
    name: str
    description: str
    parameters: dict  # JSON Schema
    handler: Callable  # 实际执行函数

REGISTERED_TOOLS: dict[str, ToolDefinition] = {}

def register(func):
    """装饰器注册工具"""
    ...

# tools/bash.py
@register
def bash(command: str) -> str:
    """Execute a shell command via subprocess."""
    ...

# tools/read_file.py
@register
def read_file(path: str, limit: int | None = None) -> str:
    """Read a file from the filesystem."""
    ...
```

**关键变化**：Go 中使用泛型函数 `Dispatch(ctx, name, args)` → Python 使用函数注册表 + `inspect.signature` 做参数校验。

### 第 4 层：MCP 集成

**覆盖包**：`agent/mcp`

Go 中自定义了 MCP Client、Transport、Discovery 等 8 个文件。

Python 可以直接使用官方 `mcp` SDK：
```python
from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client

class MCPServer:
    def __init__(self, name: str, config: MCPServerConfig):
        self.name = name
        self.config = config
        self.session: ClientSession | None = None

    async def connect(self):
        if self.config.transport == "stdio":
            params = StdioServerParameters(
                command=self.config.command,
                args=self.config.args,
                env=self.config.env,
            )
            async with stdio_client(params) as (read, write):
                async with ClientSession(read, write) as session:
                    self.session = session
                    await session.initialize()
        elif self.config.transport in ("sse", "streamable"):
            ...

    async def list_tools(self):
        result = await self.session.list_tools()
        return result.tools

    async def call_tool(self, name: str, args: dict) -> str:
        result = await self.session.call_tool(name, arguments=args)
        return result.content
```

**MCP 配置解析**（`config.py`）：
```python
class MCPServerConfig(BaseModel):
    name: str
    transport: Literal["stdio", "sse", "streamable"]
    command: str = ""
    args: list[str] = []
    env: dict[str, str] = {}
    url: str = ""
    headers: dict[str, str] = {}
    enabled: bool = True
```

支持两种配置格式（`mcp_servers` 数组 和 `mcpServers` map）。

### 第 5 层：核心运行时

**覆盖包**：`agent/runtime/`

这是迁移的核心。Go 中的主循环在 `handler.go`（已重命名为 `conversation_flow.go`）。

```python
# agent/runtime/service.py
class RuntimeService:
    def __init__(self, config: AppConfig):
        self.config = config
        self.llm = LLMClient(...)
        self.storage = ConversationHistory(...)
        self.tools = ToolRegistry(config)
        self.mcp = MCPManager(config)
        self.memory = MemoryManager(...)
        self.compression = CompressionManager(...)
        self.hooks = HookManager(...)

    async def respond(
        self,
        conversation: Conversation,
        user: User,
        text: str,
        events: EventWriter,
    ) -> None:
        """主循环：用户输入 → LLM → 工具 → 循环直到完成"""
        ...

    async def _loop(
        self,
        ctx: Context,
        messages: list[Message],
        ...
    ) -> None:
        """内部循环，处理 LLM 响应和执行工具调用"""
        ...
```

**主循环逻辑**（等价于 `conversation_flow.go`）：
```python
async def _loop(self, ...):
    while True:
        response = await self.llm.chat(
            messages=messages,
            tools=self.tools.definitions(),
        )
        msg = response.choices[0].message

        if msg.tool_calls:
            for call in msg.tool_calls:
                # 审批
                if self.needs_approval(call.function.name):
                    decision = await self.request_approval(...)
                    if decision == "no":
                        continue

                # 执行工具 → hooks
                result = await self.tools.execute(
                    name=call.function.name,
                    args=call.function.arguments,
                )

                # 大结果自动落盘（>200KB）
                if len(result) > 200_000:
                    persisted_id = self.persist_output(result)
                    result = f"<persisted-output id={persisted_id}>..."

                messages.append(response_message)
                messages.append(tool_result_message)
        else:
            # 纯文本回复
            events.emit("assistant", msg.content)
            break

        # 上下文压缩检查
        if self.compression.should_compress(messages):
            messages = await self.compression.compress(messages)
```

#### 子包：压缩策略 (`compression/`)

```python
# agent/runtime/compression/strategies.py
class CompressionStrategy(ABC):
    @abstractmethod
    async def compress(self, messages: list[Message]) -> list[Message]: ...

class WindowTrimStrategy(CompressionStrategy):
    """保留最近 N 轮对话"""
    def __init__(self, window_size: int = 20):
        self.window_size = window_size

class PersistedOutputStrategy(CompressionStrategy):
    """单轮结果 >200KB 时写入文件系统"""

class SummaryStrategy(CompressionStrategy):
    """对较早上下文进行 LLM 摘要压缩"""

class ReactiveCompactStrategy(CompressionStrategy):
    """HTTP 413 时的激进压缩"""
```

#### 子包：钩子系统 (`hooks/`)

```python
# agent/runtime/hooks/types.py
class ToolFinishHook:
    async def after_tool_execution(self, ctx, name, args, result): ...

class ModelInterceptorHook:
    async def after_model_response(self, ctx, response): ...

class ToolPersistHook:
    async def on_large_result(self, ctx, result): ...

class CompressionHook:
    async def on_context_overflow(self, ctx): ...
```

**事件驱动方式**：Go 中通过 `EventWriter` 通道向 TUI 发送事件；Python 中可以使用 `asyncio.Queue` 或回调机制。

```python
@dataclass
class Event:
    generation: int
    name: str
    content: str
    data: dict | None = None

class EventWriter:
    def __init__(self, queue: asyncio.Queue, generation: int = 0):
        self.queue = queue
        self.generation = generation

    async def emit(self, name: str, data: Any = None):
        await self.queue.put(Event(
            generation=self.generation,
            name=name,
            content=event_content(data),
            data=data,
        ))
```

事件类型（与 Go 一致）：
- `assistant_delta` — 流式文本增量
- `reasoning_delta` — 推理内容增量
- `tool_call_start` — 工具调用开始
- `tool_call_done` — 工具调用完成
- `approval_request` — 审批请求
- `error` — 错误
- `done` — 完成
- `meta` — 元信息（token 计数等）

### 第 6 层：TUI（终端界面）

**覆盖包**：`tui/`

Go 使用 Bubble Tea（Elm 架构）：
```
Model { state, view, update }
  ↓
tea.NewProgram(model).Run()
```

**Python 选择：Textual**

Textual 是 Python 最成熟的 TUI 框架，完全事件驱动，支持 CSS 样式、实时更新、鼠标事件、异步操作。

```python
# tui/app.py
from textual.app import App, ComposeResult
from textual.widgets import Header, Footer, Input, RichLog
from textual.containers import Vertical

class CynosureApp(App):
    """Cynosure TUI 主应用"""

    CSS = """
    Screen {
        layout: vertical;
    }
    #conversation {
        height: 1fr;
    }
    #input-box {
        dock: bottom;
        height: 3;
    }
    """

    def compose(self) -> ComposeResult:
        yield Header()
        yield RichLog(id="conversation", markup=True, highlight=True)
        yield Input(id="input-box", placeholder="问 cynosure 一件事...")

    async def on_input_submitted(self, message: Input.Submitted) -> None:
        """用户按 Enter 后的处理"""
        text = message.value
        # 显示用户消息
        conv = self.query_one("#conversation", RichLog)
        conv.write(f"[bold green]› {text}[/bold green]")

        # 异步调用 runtime
        asyncio.create_task(self._handle_response(text))
```

**关键 UI 元素映射**：

| Go (Bubble Tea) | Python (Textual) |
|---|---|
| `textarea.Model` | `Input` widget |
| `viewport.Model` | `RichLog` widget + `ScrollView` |
| `glamour.TermRenderer` | `RichLog` 内置 Markdown 支持 |
| `lipgloss` 样式 | CSS inline styles / `rich.style` |
| `tea.Program.Run()` | `App.run()` |
| `tea.Msg` 消息循环 | `asyncio` 事件 + 消息处理器 |
| `update()` 方法 | `on_*` 事件处理器 |

**TUI 事件消费者**：Textual 使用 `asyncio`，可以直接 `await queue.get()` 消费 Runtime 发出的事件：
```python
async def _consume_events(self):
    while True:
        event = await self.event_queue.get()
        match event.name:
            case "assistant_delta":
                self.append_assistant_delta(event.content)
            case "tool_call_start":
                self.show_tool_call(event.data)
            case "approval_request":
                self.show_approval_dialog(event.data)
            ...
```

### 第 7 层：其他辅助包

| 包 | Go | Python |
|---|---|---|
| `sessions` | 会话管理，技能摘要 | 直接翻译 |
| `cli` | 命令行入口 | `click` 或 `argparse` |
| `assistant` | 辅助功能 | 直接翻译 |
| `local` | 本地文件操作 | `pathlib` 标准库 |
| `safety` | 安全检查、脱敏 | 直接翻译 |

---

## 5. 关键设计差异与注意事项

### 5.1 并发模型

- **Go**：goroutine + channel（CSP 模型）
- **Python**：`asyncio` + `asyncio.Queue`（async/await 协程）

Runtime 与 TUI 之间的通信通道（Go 的 `chan Event`）→ Python 的 `asyncio.Queue`。

### 5.2 工具调度

Go 通过类型系统做函数注册与参数校验。Python 推荐用 `inspect` 模块 + Pydantic 做参数校验：

```python
def execute(func: Callable, args: dict) -> str:
    sig = inspect.signature(func)
    bound = sig.bind(**args)
    return func(*bound.args, **bound.kwargs)
```

### 5.3 配置结构

Go 使用 `config.AppConfig` 结构体 + JSON 文件。Python 直接使用 Pydantic：

```python
class AppConfig(BaseModel):
    workspace_root: Path
    llm: LLMConfig
    allowed_tools: list[str] = ["load_skill"]
    bash_allow_outside_workspace: bool = False
    bash_allow_dangerous_commands: bool = False
    ...
```

### 5.4 JSON Schema 验证

Go 中自定义了 `ValidateToolArgs`。Python 直接用 `jsonschema`：

```python
import jsonschema

def validate_tool_args(name: str, schema: dict, args: dict):
    jsonschema.validate(instance=args, schema=schema)
```

### 5.5 记忆系统

Python 中可考虑使用专门的向量存储（如 ChromaDB）来增强记忆能力，而非仅依赖文件系统。但为了保持纯本地、零外部基础设施的定位，初始迁移可以先保持文件系统实现。

---

## 6. 迁移步骤建议

### 阶段一：基础设施（预计 1-2 天）
1. 创建项目骨架（`pyproject.toml`、目录结构）
2. 迁移 `config`, `idgen`, `logger`, `textutil`, `safety`
3. 设置依赖（`openai`, `textual`, `mcp`, `pydantic`, `jsonschema`）

### 阶段二：核心服务（预计 3-5 天）
4. 迁移 `llm` 包（OpenAI 客户端封装）
5. 迁移 `agent/storage`（会话持久化）
6. 迁移 `tools` 包（工具定义与调度）

### 阶段三：运行时（预计 5-7 天）
7. 迁移 `agent/mcp`（MCP 集成）
8. 迁移 `agent/runtime` 主循环
9. 迁移 `compression` 和 `hooks` 子包

### 阶段四：TUI 与集成（预计 3-5 天）
10. 用 Textual 重写 `tui` 包
11. 编写 `__main__.py` 入口
12. 集成测试与端到端验证

### 阶段五：打磨（预计 2-3 天）
13. 调试、修复边缘情况
14. 性能优化（尤其是上下文压缩和流式渲染）
15. 补充单元测试

**总计预估**：14-22 天（视经验而定）

---

## 7. 风险与缓解

| 风险 | 缓解 |
|------|------|
| Textual 的流式渲染性能 | 使用 `RichLog` 增量写入而非每次重建 |
| asyncio 事件循环阻塞 | 将 CPU 密集型操作（如 JSON 序列化）移至线程池 |
| MCP 子进程管理 | 使用 `asyncio.create_subprocess_exec` + 生命周期 hook |
| 大文件读取/写入性能 | 沿用 Go 中的分段读取策略（`limit` 参数） |
| 上下文压缩精度损失 | 使用与 Go 版本相同的 LLM 编写摘要 prompt |

---

## 8. 不迁移的内容

以下功能在 Python 迁移中保持不变：
- 文件存储路径：`~/.cynosure/`（与 Go 版本兼容）
- 配置格式：`.json`（`config.json`, `.mcp.json`）
- MCP 工具命名格式：`mcp__{server}__{tool}`
- 大结果落盘路径：`~/.cynosure/task_outputs/tool-results/`
- 会话存储路径：`~/.cynosure/session/<id>/`
- 技能目录：`~/.cynosure/skills/` + `./.cynosure/skills/`

---

## 9. 总结

Cynosure 的 Python 迁移是一次**逐层直译 + 平台适配**的工作：

- **架构不变**：13 个内部包的职责划分完全保留
- **主循环不变**：用户输入 → LLM → 工具 → 循环的流程不变
- **存储格式不变**：JSON 文件路径、格式完全兼容
- **TUI 重写**：Bubble Tea → Textual（同为消息驱动架构）
- **MCP 简化**：自定义实现 → 官方 SDK
- **工具系统简化**：泛型函数调度 → 装饰器注册

迁移后，Python 版本将保留全部功能，同时获得更好的生态集成能力和维护便利性。
