# Sub-Agent Tool Status Display Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Show only the latest sub-agent child tool call in the parent TUI and render complete tool arguments from raw JSON.

**Architecture:** Keep the runtime event stream unchanged and implement the behavior at the TUI message/rendering layer. `appendToolCallStart` owns ephemeral sub-agent message replacement, while `renderToolMessage` owns complete argument display.

**Tech Stack:** Go 1.26.1, Bubble Tea, lipgloss, existing `internal/tui` tests.

## Global Constraints

- Other strategies remain unchanged.
- Do not change runtime event emission, result suppression, or group clearing.
- Preserve sub-agent gray styling, no leading blue bullet, and status-column alignment.
- Render sub-agent child tool rows as compact `tool(args...)` information only, with no icon or status/result line.
- Display every top-level JSON argument when `raw_args` is available.
- Sub-agent internal tool calls must not change the parent UI tool-call count.
- Concurrent sub-agent child rows must attach to the matching parent `spawn_subagent` by `parent_tool_call_id`.

---

### Task 1: Latest Sub-Agent Tool Message

**Files:**
- Modify: `internal/tui/events_test.go`
- Modify: `internal/tui/app.go`

**Interfaces:**
- Consumes: `ToolCallView.Scope`, `ToolCallView.EphemeralGroupID`, `Model.appendToolCallStart(data any)`
- Produces: `appendToolCallStart` removes older messages in the same sub-agent ephemeral group before appending a new child tool start.

- [ ] **Step 1: Write the failing test**

Add `TestModelShowsOnlyLatestSubagentToolStatus` to `internal/tui/events_test.go`. It should start two sub-agent tool calls with the same `ephemeral_group_id` and assert that the first tool is gone while the second remains.

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/tui -run TestModelShowsOnlyLatestSubagentToolStatus -count=1`
Expected: FAIL because current behavior keeps both child tool messages.

- [ ] **Step 3: Write minimal implementation**

In `Model.appendToolCallStart`, after defaulting status and before appending, if `tool.Scope == "subagent"` and `tool.EphemeralGroupID` is non-empty, filter `m.messages` to remove existing tool messages with the same `EphemeralGroupID` and `Scope == "subagent"`.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/tui -run TestModelShowsOnlyLatestSubagentToolStatus -count=1`
Expected: PASS.

### Task 2: Complete Raw Argument Rendering

**Files:**
- Modify: `internal/tui/events_test.go`
- Modify: `internal/tui/app.go`

**Interfaces:**
- Consumes: `ToolCallView.RawArgs`, `ToolCallView.ArgsPreview`
- Produces: `toolArgsDisplay(rawArgs, fallback string) string`, returning a stable complete argument display from top-level JSON object keys or falling back to `args_preview`.

- [ ] **Step 1: Write the failing test**

Add `TestToolMessageRendersAllRawArgs` to `internal/tui/events_test.go`. Use a tool with `raw_args` containing at least `pattern`, `path`, and `output_mode`, but `args_preview` containing only `pattern`. Assert all raw parameters appear in rendered output.

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/tui -run TestToolMessageRendersAllRawArgs -count=1`
Expected: FAIL because only `args_preview` is rendered.

- [ ] **Step 3: Write minimal implementation**

Add helpers in `internal/tui/app.go`:

```go
func toolArgsDisplay(rawArgs, fallback string) string
func formatToolArgValue(value any) string
```

Use `json.Unmarshal` into `map[string]any`, sort keys with `sort.Strings`, format scalar values with `fmt.Sprint`, and marshal arrays/objects back to compact JSON. Update `renderToolMessage` to call `toolArgsDisplay(tool.RawArgs, tool.ArgsPreview)`.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/tui -run TestToolMessageRendersAllRawArgs -count=1`
Expected: PASS.

### Task 3: Regression Verification

**Files:**
- Test: `internal/tui/events_test.go`
- Modify: `internal/agent/runtime/subagent.go`
- Test: `internal/agent/runtime/runtime_test.go`
- Test: project build

**Interfaces:**
- Consumes: completed Task 1 and Task 2 behavior.
- Produces: verified narrow change with existing behavior preserved, including stale child done filtering and parent tool-count isolation.

- [ ] **Step 1: Add stale child done and meta-count regression tests**

Add tests proving:

- Sub-agent child tool rows attach directly under the parent running block and show no checkmark, running line, success line, or result preview.
- A replaced sub-agent child tool is not re-added by a later `tool_call_done`.
- `subagentEventWriter` does not forward child `tool_call_count` in `meta` payloads.
- `subagentEventWriter` forwards `parent_tool_call_id` on child tool lifecycle events, and the TUI inserts child rows after the matching parent tool row even when multiple sub-agents run concurrently.

- [ ] **Step 2: Implement stale done and meta-count filtering**

In `updateToolCallDone`, ignore unmatched sub-agent done events that still carry an `ephemeral_group_id`.

In `subagentEventWriter.Event`, delete `tool_call_count` from forwarded `meta` payloads.

- [ ] **Step 3: Run focused TUI tests**

Run: `go test ./internal/tui -count=1`
Expected: PASS.

- [ ] **Step 4: Run focused runtime tests**

Run: `go test ./internal/agent/runtime -run 'TestSubagentEventWriterDoesNotForwardChildToolCallCount|TestRespondToConversation_SubagentForwardsToolStatusWithoutResultAndClearsGroup' -count=1`
Expected: PASS.

- [ ] **Step 5: Run full verification**

Run: `go test ./...`
Expected: PASS.

Run: `go build ./...`
Expected: PASS.
