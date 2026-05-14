## Why

当前 nano_cc 基于 TypeScript + Bun 实现，而团队主要技术栈为 Go。使用 Go 重写 agent 服务可以：统一技术栈降低维护成本、编译为单一二进制便于部署、利用 Go 的并发模型提升子智能体并行执行效率、获得更强的类型安全保障。

## What Changes

- 新建 Go 项目，实现与 nano_cc 等价的 agent 服务
- 实现核心 agent 循环（OpenAI 兼容 API 工具调用模式）
- 实现工具系统：bash 执行、文件读写/编辑、todo 管理、子智能体委派
- 实现技能按需加载机制（两层注入架构）
- 实现三层上下文压缩管道
- 实现文件系统安全沙箱和命令执行安全检查
- 支持通过环境变量配置 LLM 后端

## Capabilities

### New Capabilities
- `agent-loop`: 核心智能体循环，基于 OpenAI 兼容 API 的工具调用 while 循环，集成上下文压缩和 todo 提醒
- `tool-system`: 工具注册、定义和分发执行，支持父子智能体分级工具集
- `subagent`: 子智能体创建与上下文隔离，限制轮次，仅返回文本摘要
- `skill-loading`: 技能按需加载，两层注入架构（名称描述注入 system prompt + 调用时返回完整内容）
- `context-compact`: 三层上下文压缩（micro_compact / auto_compact / compact 工具）
- `file-operations`: 文件读写与精确字符串替换，路径沙箱保护
- `bash-execution`: Shell 命令执行，危险命令黑名单拦截，超时控制
- `task-management`: 任务列表管理，支持创建、更新、依赖跟踪

### Modified Capabilities
<!-- 无现有 spec 需要修改 -->

## Impact

- 新增 Go 项目目录（与现有 TypeScript 项目并存，不影响现有代码）
- 依赖：Go 标准库 + OpenAI Go SDK + YAML 解析库
- 部署：编译为单一二进制，通过环境变量配置
- 现有 nano_cc TypeScript 项目保持不变