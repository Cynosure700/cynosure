---
name: "go-agent项目的技术栈与架构"
description: "go-agent是一个Go语言编写的TUI通用智能体，使用Deepseek、Bubble Tea、SQLite等。"
metadata:
  node_type: memory
  type: "project_fact"
  project: "go-agent"
  originSessionId: "mem_556e56df2d0459031874409e"
---

go-agent由Go编写，LLM调用Deepseek（通过OpenAI兼容API），TUI基于Bubble Tea+Lipgloss，存储用SQLite，配置为JSON。核心运行流程：用户输入→CLI初始化→TUI→Runtime循环（构建prompt→调用LLM→执行工具→循环直到纯文本回复）。内置工具包括bash、文件读写、编辑、技能加载、子智能体等，并支持MCP集成和技能系统。
