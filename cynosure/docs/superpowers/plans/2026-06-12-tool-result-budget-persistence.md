# Tool Result Budget Persistence Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Persist oversized `tool_result` outputs under `~/.cynosure/task_outputs/{workspace}/{session_id}/tool-results/`, keep only `<persisted-output>` plus a 2000-character preview in model context, and append every tool execution result to `~/.cynosure/task_outputs/{workspace}/{session_id}/tools.md`.

**Architecture:** Reuse the existing compression algorithm and move durability into the local store. Add focused local file helpers for persisted-output metadata/content and tool-result Markdown logging, then wire the default tool hook through an optional store interface.

**Tech Stack:** Go, local filesystem storage, existing runtime hooks, existing compression pipeline, Go unit tests.

---

### Task 1: Failing tests

**Files:**
- Modify: `internal/local/store_test.go`
- Modify: `internal/agent/runtime/runtime_test.go`

- [ ] Add tests for persisted-output file creation, fallback read after memory loss, sha256 validation, and `tools.md` append behavior.
- [ ] Add runtime hook test assertion that a store implementing tool-result logging receives the tool result.
- [ ] Run targeted tests and confirm they fail because methods/types are missing.

### Task 2: Local persisted-output files

**Files:**
- Create: `internal/local/persisted_output_files.go`
- Modify: `internal/local/store.go`

- [ ] Add path-safe output id validation and metadata struct.
- [ ] Write `{id}.txt` and `{id}.json` under `~/.cynosure/task_outputs/{workspace}/{session_id}/tool-results/` before updating in-memory indexes.
- [ ] Add file fallback in `GetPersistedOutputForConversation` with user/conversation validation and sha256 check.
- [ ] Run local package tests.

### Task 3: Tool result Markdown log

**Files:**
- Modify: `internal/agent/storage/models.go`
- Create: `internal/local/tool_result_log.go`
- Modify: `internal/agent/runtime/hooks/tool.go`

- [ ] Add `storage.ToolResultLogEntry`.
- [ ] Implement `local.Store.AppendToolResultLog` to append Markdown under `~/.cynosure/task_outputs/{workspace}/{session_id}/tools.md`.
- [ ] Call optional logging interface from the tool append hook; warn and continue on errors.
- [ ] Run runtime and local package tests.

### Task 4: README and verification

**Files:**
- Modify: `README.md`

- [ ] Document `~/.cynosure/task_outputs/{workspace}/{session_id}/tool-results/` and `~/.cynosure/task_outputs/{workspace}/{session_id}/tools.md`.
- [ ] Run `gofmt` on changed Go files.
- [ ] Run targeted tests, then `go test ./...` if feasible.
