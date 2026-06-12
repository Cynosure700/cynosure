# TUI Local Cleanup Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remove Web/DB/Redis/ES runtime remnants while preserving TUI local memory, compression, history, skills, and MCP behavior.

**Architecture:** Keep `internal/local` as the only runtime store for TUI. Keep `internal/agent/storage` only as a pure model/serialization package. Remove DB-backed skills and MCP user sessions from runtime paths, and preserve context summaries as in-memory cache only.

**Tech Stack:** Go 1.26, local JSON/Markdown files, Bubble Tea TUI, MCP Go SDK, OpenAI-compatible LLM client.

---

## File Map

- Modify `internal/agent/runtime/runtime.go`: remove DB skill method from runtime store interface.
- Modify `internal/agent/runtime/prompt_builder.go`: build skill snapshot only from local loader.
- Modify `internal/tools/definitions.go`: remove database skill wording.
- Modify `internal/agent/mcp/manager.go`: remove DB user MCP store/session logic.
- Modify `internal/agent/runtime/tool_registry.go`: stop ensuring DB user MCP sessions.
- Modify `internal/local/store.go`: remove no-op DB skill/MCP methods and add persisted output hash recovery.
- Modify `internal/local/persisted_output_files.go`: add metadata scanning helper.
- Modify tests under `internal/agent/runtime`, `internal/agent/mcp`, `internal/local`, `internal/tools`.
- Delete DB implementation files under `internal/agent/storage`, keeping `models.go`, `conversation_history.go`, and tests.
- Modify `internal/config/config.go`, `config.json`, `go.mod`, `go.sum`: remove DB/Redis/ES/JWT configuration and dependencies.
- Modify `README.md`: document local-only TUI storage.

## Tasks

### Task 1: Keep context summary in memory and add persisted output hash recovery

- [ ] Add a failing test in `internal/local/store_test.go` that writes a persisted output, creates a new store for the same workspace/session, calls `GetPersistedOutputByMessageHash`, and expects the existing file-backed output.
- [ ] Run `go test ./internal/local -run TestStoreRestoresPersistedOutputByMessageHashFromWorkspace -count=1`; expected failure is `sql: no rows in result set`.
- [ ] Add metadata scanning in `internal/local/persisted_output_files.go` and call it from `GetPersistedOutputByMessageHash` in `internal/local/store.go`.
- [ ] Re-run the same test and expect PASS.

### Task 2: Remove DB skill runtime path

- [ ] Update runtime skill tests to expect local user/workspace loaders only, without `storage.Skill` DB entries.
- [ ] Run focused runtime/tool tests and confirm failures reference `ListEnabledSkillsByUser` or old database wording.
- [ ] Remove `ListEnabledSkillsByUser` from `conversationStore`, remove `buildDBSkillLoader`, and build snapshots directly from `BuiltinSkills`.
- [ ] Update `load_skill` description to local user/workspace wording.
- [ ] Re-run focused tests and expect PASS.

### Task 3: Remove DB MCP user session path

- [ ] Update MCP manager/runtime tests to validate builtin + workspace sessions without fake DB store.
- [ ] Run focused MCP/runtime tests and confirm failures reference user DB sessions or fake store interface.
- [ ] Remove `mcp.Store`, `Manager.sessions`, and `EnsureUserSessions`; make `NewManager` require no store.
- [ ] Update bootstrap and tool registry to use `mcp.NewManager()` and stop calling `EnsureUserSessions`.
- [ ] Re-run focused tests and expect PASS.

### Task 4: Delete DB/Redis/ES storage implementation

- [ ] Delete DB-backed files in `internal/agent/storage`, keeping `models.go`, `conversation_history.go`, and `conversation_history_test.go`.
- [ ] Run `go test ./...`; expected failures reveal remaining references.
- [ ] Remove remaining references to deleted storage methods/types that are no longer needed.
- [ ] Re-run tests until compile passes.

### Task 5: Clean config and dependencies

- [ ] Remove DB/Redis/ES/JWT fields from `internal/config/config.go` and `config.json`.
- [ ] Delete `internal/config/env.go` if unused.
- [ ] Run `go mod tidy`.
- [ ] Verify `go.mod` no longer lists MySQL, Redis, Elasticsearch, or JWT direct dependencies.

### Task 6: README and verification

- [ ] Update `README.md` with local-only storage locations: memory, session history, model history, skills, MCP, tool outputs.
- [ ] Run `gofmt` on changed Go files.
- [ ] Run `go test ./...`; expected PASS.
- [ ] Report changed files and verification output.
