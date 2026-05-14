## 1. Memory Loading Module

- [x] 1.1 Implement `internal/sessions/memory.go`: create `LoadProjectMemory()` that reads `AGENTS.md` from working directory, returns content string (empty if file not found)
- [x] 1.2 Implement `LoadUserMemory()` that reads `~/.link/AGENTS.md`, returns content string (empty if file not found)
- [x] 1.3 Implement `BuildPersistentMemorySection()` that formats loaded memory into `<project_memory>` and `<user_memory>` XML blocks for system prompt injection
- [x] 1.4 Implement `SessionMemory` struct with `Add(entry string)` and `GetSection()` methods, thread-safe with `sync.RWMutex`

## 2. System Prompt Integration

- [x] 2.1 Modify `internal/agent/repl.go`: call memory loading functions during system prompt construction, append persistent memory section after skill descriptions
- [x] 2.2 Create global `SessionMemory` instance in `internal/sessions/memory.go`, initialize in REPL startup
- [x] 2.3 Verify persistent memory section appears correctly in system prompt when files exist, and is absent when files don't exist

## 3. Session Memory Injection

- [x] 3.1 Modify `internal/agent/loop.go`: after user message is added to messages array, inject `<session_memory>` block if session memory is non-empty
- [x] 3.2 Ensure session memory injection happens before micro compaction in the agent loop
- [x] 3.3 Verify session memory is preserved across context compaction (not stripped by micro/auto/manual compact)

## 4. update_memory Tool

- [x] 4.1 Implement `internal/tools/memory.go`: create `handleUpdateMemory` handler that supports `scope` (session/project), `action` (append/replace for project scope), and `content` parameters
- [x] 4.2 For `scope: "session"`: append content to global `SessionMemory` instance
- [x] 4.3 For `scope: "project"`: add path safety check via `safety.SafePath()`, then append or replace `AGENTS.md`
- [x] 4.4 Register `update_memory` tool in `internal/tools/registry.go`: add tool definition with JSON Schema for scope/action/content params, add to `ParentToolDefs` only (not child tools)
- [x] 4.5 Return confirmation message on success, error message on failure (invalid scope, invalid action, path safety violation)

## 5. Integration & Verification

- [x] 5.1 End-to-end test: start REPL with existing `AGENTS.md`, verify persistent memory appears in system prompt
- [x] 5.2 Test `update_memory` session scope: call tool, verify session memory injected in next turn
- [x] 5.3 Test `update_memory` project append: call tool, verify content appended to `AGENTS.md`
- [x] 5.4 Test `update_memory` project replace: call tool, verify file overwritten
- [x] 5.5 Test edge cases: missing memory files (no error), empty memory files, path traversal rejection, session memory cleared on exit