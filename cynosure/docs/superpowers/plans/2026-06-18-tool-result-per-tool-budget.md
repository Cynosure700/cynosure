# Tool Result Per-Tool Budget Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Change tool result compression from a global latest-turn byte budget to per-tool `maxResultSizeChars`, defaulting to 50,000 characters, while reusing the existing persisted output directory.

**Architecture:** Keep tool execution and persisted output storage boundaries intact. Add result-size metadata next to built-in tool definitions, pass a resolver into request-only compression, and make `ToolResultCompressionStrategy` persist each oversized individual result.

**Tech Stack:** Go, existing `go-openai` tool definitions, `internal/tools`, `internal/agent/runtime`, `internal/agent/runtime/compression`, existing local persisted output files under `~/.cynosure/task_outputs/{workspace}/{session_id}/tool-results/`.

## Global Constraints

- Default `maxResultSizeChars` is exactly `50,000` characters.
- Unknown tools and MCP tools use the default `50,000` character limit.
- Existing persisted output directory `~/.cynosure/task_outputs/{workspace}/{session_id}/tool-results/` stays unchanged.
- Existing marker format and `read_persisted_output` behavior stay unchanged.
- Other compression strategies and their order stay unchanged.
- Use TDD: write failing tests before implementation changes.

---

### Task 1: Tool Metadata

**Files:**
- Modify: `internal/tools/definitions.go`
- Test: `internal/tools/definitions_test.go`

**Interfaces:**
- Produces: `tools.DefaultMaxResultSizeChars int`
- Produces: `tools.ToolSpec`
- Produces: `tools.AllToolSpecs []ToolSpec`
- Produces: `tools.MaxResultSizeCharsForTool(name string) int`
- Preserves: `tools.AllToolDefs []openai.Tool`

- [ ] **Step 1: Write failing tests**

Add tests that assert:

```go
func TestAllToolSpecsExposeDefaultMaxResultSizeChars(t *testing.T) {
	if DefaultMaxResultSizeChars != 50000 {
		t.Fatalf("DefaultMaxResultSizeChars = %d, want 50000", DefaultMaxResultSizeChars)
	}
	if MaxResultSizeCharsForTool("bash") != 50000 {
		t.Fatalf("bash max result chars = %d, want 50000", MaxResultSizeCharsForTool("bash"))
	}
	if MaxResultSizeCharsForTool("mcp__unknown__tool") != 50000 {
		t.Fatalf("unknown max result chars = %d, want 50000", MaxResultSizeCharsForTool("mcp__unknown__tool"))
	}
}
```

Also assert that `AllToolDefs` and `AllToolSpecs` have the same tool-name set.

- [ ] **Step 2: Run red test**

Run: `go test ./internal/tools -run 'TestAllToolSpecs|TestAllToolDefinitions' -count=1`

Expected: FAIL because `DefaultMaxResultSizeChars`, `AllToolSpecs`, or `MaxResultSizeCharsForTool` is undefined.

- [ ] **Step 3: Implement metadata**

Convert internal tool definition construction from `toolDef(...) openai.Tool` to `toolSpec(...).Definition`, derive `AllToolDefs` from specs, and keep existing exported tool defs for `web_search`, `spawn_subagent`, and `read_persisted_output`.

- [ ] **Step 4: Run green test**

Run: `go test ./internal/tools -count=1`

Expected: PASS.

### Task 2: Compression Uses Per-Tool Limit

**Files:**
- Modify: `internal/agent/runtime/compression/compression.go`
- Modify: `internal/agent/runtime/compression/tool_result_compression.go`
- Test: `internal/agent/runtime/compression/compression_test.go`

**Interfaces:**
- Consumes: `tools.DefaultMaxResultSizeChars`
- Produces: `compression.Request.ToolMaxResultSizeChars func(toolName string) int`
- Preserves: persisted output marker and store interfaces

- [ ] **Step 1: Write failing compression tests**

Add tests for:

```go
func TestToolResultCompression_PersistsSingleResultOverToolLimit(t *testing.T)
func TestToolResultCompression_DoesNotCompressByCombinedTotal(t *testing.T)
func TestToolResultCompression_UsesDefaultLimitWhenToolNameMissing(t *testing.T)
func TestToolResultCompression_CountsRunesForLimit(t *testing.T)
```

The first uses resolver `bash -> 10` and expects an 11-character result to persist. The combined-total test uses two 120 KB results with per-tool limit 200 KB and expects no persistence, proving old global total budget no longer applies.

- [ ] **Step 2: Run red test**

Run: `go test ./internal/agent/runtime/compression -run 'TestToolResultCompression' -count=1`

Expected: FAIL because compression still uses global byte total and has no resolver field.

- [ ] **Step 3: Implement per-tool strategy**

Add `ToolMaxResultSizeChars` to `Request`, build `toolCallID -> toolName`, use rune count for limit comparison, and persist only candidates whose individual result length exceeds their tool limit.

- [ ] **Step 4: Run green test**

Run: `go test ./internal/agent/runtime/compression -count=1`

Expected: PASS.

### Task 3: Runtime Passes Limit Resolver

**Files:**
- Modify: `internal/agent/runtime/tool_registry.go`
- Modify: `internal/agent/runtime/context_compression.go`
- Test: `internal/agent/runtime/context_compression_test.go`
- Test: `internal/agent/runtime/runtime_test.go`

**Interfaces:**
- Consumes: `tools.MaxResultSizeCharsForTool(name string) int`
- Produces: `(*ToolRegistry).MaxResultSizeChars(name string) int`

- [ ] **Step 1: Write failing runtime tests**

Add or adjust tests so request compression uses a tool-specific low limit and the request history gets a marker while stored display/model histories remain complete.

- [ ] **Step 2: Run red test**

Run: `go test ./internal/agent/runtime -run 'TestCompressRequestHistory|TestToolRegistry' -count=1`

Expected: FAIL because runtime does not pass a result limit resolver.

- [ ] **Step 3: Implement resolver wiring**

Expose `ToolRegistry.MaxResultSizeChars`, and pass it into `compression.Request` wherever request compression is constructed.

- [ ] **Step 4: Run green test**

Run: `go test ./internal/agent/runtime -count=1`

Expected: PASS.

### Task 4: Remove Silent Output Truncation

**Files:**
- Modify: `internal/tools/bash.go`
- Modify: `internal/tools/search_ops.go`
- Modify: `internal/tools/web_ops.go`
- Test: relevant tests under `internal/tools`

**Interfaces:**
- Preserves: timeouts, `webFetchMaxBodySize`, `head_limit`, `read_file` limits
- Removes: return-before-compression truncation for bash, grep content, web_fetch cleaned text

- [ ] **Step 1: Write failing tests**

Add tests confirming bash/search/web_fetch can return more than 50,000 characters when the underlying tool produces that much output within existing resource limits.

- [ ] **Step 2: Run red tests**

Run: `go test ./internal/tools -run 'Bash|Grep|WebFetch' -count=1`

Expected: FAIL where output is still truncated to 50,000.

- [ ] **Step 3: Remove truncation**

Remove only the `maxOutputLen` truncation from bash, grep content output, and web_fetch cleaned text. Keep all resource limits.

- [ ] **Step 4: Run green tests**

Run: `go test ./internal/tools -count=1`

Expected: PASS.

### Task 5: Final Verification

**Files:**
- Review all modified files

- [ ] **Step 1: Run focused tests**

Run: `go test ./internal/tools ./internal/agent/runtime/... ./internal/local -count=1`

Expected: PASS.

- [ ] **Step 2: Run full test suite**

Run: `go test ./... -count=1`

Expected: PASS.

- [ ] **Step 3: Requirement checklist**

Check the final diff against the design doc:

- per-tool declaration exists;
- default is 50,000 chars;
- oversized individual results persist to existing directory;
- marker/read tool unchanged;
- global total tool_result_budget no longer triggers this strategy;
- other compression strategies unchanged;
- silent truncation removed only where required.
