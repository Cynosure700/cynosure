---
name: "go-agent项目详细目录结构与核心架构"
description: "助手在总结中列出了go-agent项目的完整文件树、核心架构流程和关键特性"
metadata:
  node_type: memory
  type: "project_fact"
  project: "go-agent"
  originSessionId: "mem_ba8ed86c6f1561633559134b"
---

项目结构：main.go → internal/（cli/、tui/、config/、llm/、agent/（runtime/、storage/、mcp/）、tools/、local/、sessions/、logger/、safety/、idgen/、assistant/） → config.json → skills/ → workspace/ → docs/ → .link/memory/。核心流程：用户输入→CLI初始化→TUI→Runtime主循环（构建提示词→LLM→解析响应→工具调用或文本回复→循环）→会话持久化（SQLite）。关键特性包括：工具系统（bash、文件、技能加载等）、MCP集成、技能系统、上下文压缩（多种策略）、记忆系统（Markdown文件）、子智能体、钩子系统、TUI（Bubble Tea）、路径安全。配置包含LLM端点/模型、工具开关、MCP配置、路径白名单。
