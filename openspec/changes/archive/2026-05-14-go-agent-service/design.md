## Context

nano_cc 是一个基于 TypeScript + Bun 的 CLI 编码智能体，核心架构为：CLI REPL → Agent 循环（OpenAI 兼容 API 工具调用 while 循环）→ 工具处理器。当前项目需要 Go 语言版本，以统一团队技术栈、简化部署（单一二进制）、利用 Go 并发模型提升子智能体并行执行效率。

## Goals / Non-Goals

**Goals:**
- 用 Go 实现与 nano_cc 功能等价的 agent 服务
- 保持相同的工具集、技能系统、上下文压缩策略
- 编译为单一二进制，通过环境变量配置
- 支持 OpenAI 兼容 API（GLM-5 等）

**Non-Goals:**
- 不实现 HTTP/gRPC 服务端（仅 CLI REPL）
- 不实现会话持久化
- 不实现 MCP 协议支持
- 不修改现有 TypeScript 项目

## Decisions

### 1. 包结构：扁平化 vs 分层

**选择**：扁平化包结构，按功能域划分。

```
go-agent/
├── main.go              # 入口
├── internal/
│   ├── agent/           # 核心 Agent 循环
│   ├── config/          # 配置 & OpenAI 客户端
│   ├── tools/           # 工具注册 + 各工具实现
│   ├── sessions/        # 子 Agent、技能、压缩
│   ├── safety/          # 路径安全
│   └── logger/          # 彩色日志
└── skills/              # 技能 Markdown 文件
```

**理由**：Go 标准项目布局，`internal/` 防止外部导入。每个包职责单一，与 TypeScript 版本的文件组织一一对应。

### 2. OpenAI SDK：官方 SDK vs 手写 HTTP 客户端

**选择**：使用 `github.com/sashabaranov/go-openai` 官方社区 SDK。

**备选**：手写 HTTP 客户端更灵活但增加维护成本。

**理由**：该 SDK 支持自定义 `BaseURL`，兼容任何 OpenAI 兼容 API。类型定义完善，减少手写 JSON Schema 的工作量。

### 3. 工具系统：接口 vs 函数映射

**选择**：函数映射 `map[string]ToolHandler`，与 TypeScript 版本一致。

```go
type ToolDef struct {
    Type     string       `json:"type"`
    Function FunctionDef  `json:"function"`
}

type FunctionDef struct {
    Name        string          `json:"name"`
    Description string          `json:"description"`
    Parameters  json.RawMessage `json:"parameters"`
}

type ToolHandler func(ctx context.Context, args map[string]any) (string, error)
```

**理由**：简单直接，与现有架构一致。工具数量有限（8 个），不需要接口抽象。

### 4. 并发模型：Goroutine vs 单线程

**选择**：主循环保持单线程顺序执行（与 TS 版本一致），子 Agent 可并行执行。

**理由**：Agent 循环的工具调用有严格顺序依赖（前一个工具的输出影响后续决策），不适合并行。但多个 `task` 调用可以并行委派给子 Agent，利用 Go 的 goroutine 提升效率。

### 5. 上下文压缩：原地修改 vs 不可变

**选择**：`microCompact` 原地修改切片，`autoCompact` 返回新切片。

**理由**：与 TypeScript 版本行为一致。`microCompact` 是高频操作（每轮执行），原地修改避免分配。`autoCompact` 低频（token 超阈值时），返回新切片更安全。

### 6. 技能文件格式：保持 Markdown + YAML frontmatter

**选择**：保持与 TypeScript 版本完全相同的技能文件格式。

**理由**：技能文件是共享资源，不应因语言切换而改变格式。Go 使用 `gopkg.in/yaml.v3` 解析 frontmatter。

### 7. 日志输出：ANSI 转义码 vs 第三方库

**选择**：手写 ANSI 转义码。

**理由**：避免引入额外依赖。日志需求简单（颜色 + 前缀），手写足够。

## Risks / Trade-offs

- **[Risk] OpenAI Go SDK 与 TypeScript SDK 行为差异** → 使用相同的 API 参数，通过集成测试验证
- **[Risk] Go 错误处理风格与 TS 不同（error 返回值 vs try/catch）** → 统一在 handler 层捕获 panic，返回错误字符串给 LLM
- **[Risk] 子 Agent 并行执行可能导致文件竞争** → 初期保持顺序执行，后续按需引入文件锁
- **[Trade-off] 单一二进制部署方便但失去热重载** → 提供 `--watch` 模式（可选，非 MVP）

## Open Questions

- 是否需要支持流式输出（SSE）？TypeScript 版本未实现，Go 版本初期也不实现
- 技能文件目录路径：硬编码 `./skills/` 还是可配置？初期硬编码，后续可加环境变量