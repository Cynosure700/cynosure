# Tool Call Terminal Display Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在 go-agent TUI 中实时展示工具调用开始、运行中、完成/失败摘要，使用户能在终端看到类似 Claude Code 的工具调用流程。

**Architecture:** Runtime 在工具执行前后通过现有 EventWriter 下发 `tool_call_start` / `tool_call_done` 事件；TUI 将这些事件合并成工具过程卡片并渲染。事件只面向 UI，不进入模型上下文，不改变工具执行、审计日志、历史持久化和上下文压缩。

**Tech Stack:** Go、Bubble Tea、Lipgloss、现有 runtime hooks、go test。

---

## Files

- Modify: `internal/agent/runtime/conversation_flow.go` — 工具事件发送与摘要函数。
- Modify: `internal/agent/runtime/runtime_test.go` — Runtime 工具事件顺序与失败事件测试。
- Modify: `internal/tui/app.go` — TUI 工具消息结构、事件处理、渲染样式。
- Modify: `internal/tui/events_test.go` or `internal/tui/render_test.go` — TUI 工具事件更新与渲染测试。

## Task 1: Runtime tool event tests

- [ ] Add tests in `internal/agent/runtime/runtime_test.go` that run a tool-calling response and assert the emitted event sequence contains `tool_call_start`, `tool_call_done`, then updated `meta`.
- [ ] Add a failure-path test where a bad tool call emits `tool_call_done` with `status=rejected` and a non-empty result preview.
- [ ] Run: `go test ./internal/agent/runtime -run 'Tool.*Event|Respond'` and verify the new tests fail before implementation.

## Task 2: Runtime implementation

- [ ] Add constants for `tool_call_start` and `tool_call_done` next to `assistantDeltaEvent` / `reasoningDeltaEvent`.
- [ ] Add `emitToolCallStart(state, toolCtx)` after `RunPreToolUse` and before `executeToolCall`.
- [ ] Add `emitToolCallDone(state, toolCtx)` after `RunPostToolUse`.
- [ ] Add helpers to build `args_preview` and `result_preview` with bounded lengths.
- [ ] Run: `go test ./internal/agent/runtime` and verify pass.

## Task 3: TUI tests

- [ ] Add tests that create a `Model`, feed `tool_call_start`, and assert one `tool` role message exists.
- [ ] Add tests that feed matching `tool_call_done` and assert the same message updates to success/rejected instead of appending duplicate messages.
- [ ] Add a render test asserting output includes a Claude Code-like tool line and `⎿` result line.
- [ ] Run: `go test ./internal/tui` and verify tests fail before implementation.

## Task 4: TUI implementation

- [ ] Extend `Message` with `ToolCall *ToolCallView` and add `ToolCallView` fields: ID, Name, ArgsPreview, Status, ResultPreview.
- [ ] Handle `tool_call_start` / `tool_call_done` in `Model.Update`.
- [ ] Add `appendToolCallStart`, `updateToolCallDone`, and event parsing helpers.
- [ ] Add `renderToolMessage` and styles for running/success/rejected states.
- [ ] Run: `go test ./internal/tui` and verify pass.

## Task 5: Regression and similarity evaluation

- [ ] Run: `go test ./internal/agent/runtime ./internal/tui ./internal/agent/runtime/hooks`.
- [ ] Run: `go test ./...` if targeted tests pass.
- [ ] Compare against Claude Code style on five dimensions: lifecycle visibility, compact syntax, status clarity, bounded output, non-interference with final answer. If score < 95, refine render wording/style and rerun tests.

## Self-review

- Spec coverage: lifecycle display, built-in/MCP/subagent support, bounded args/result previews, no model-history changes, tests are covered.
- Placeholder scan: no TODO/TBD placeholders.
- Type consistency: event names and TUI fields are defined before use.
