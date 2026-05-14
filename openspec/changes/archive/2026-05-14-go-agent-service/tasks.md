## 1. Project Setup

- [x] 1.1 Initialize Go module (`go mod init nano_cc`) and create directory structure (`internal/{agent,config,tools,sessions,safety,logger}`)
- [x] 1.2 Add dependencies: `go-openai` for LLM client, `gopkg.in/yaml.v3` for skill frontmatter parsing
- [x] 1.3 Create `main.go` entry point that calls into the REPL

## 2. Configuration & Logging

- [x] 2.1 Implement `internal/config/config.go`: read `OPENAI_BASE_URL`, `OPENAI_API_KEY`, `MODEL_ID` from environment, initialize OpenAI client with custom base URL
- [x] 2.2 Implement `internal/logger/logger.go`: colored terminal output with ANSI codes (info, warn, error, success, tool, assistant, user prefixes)

## 3. Path Safety

- [x] 3.1 Implement `internal/safety/path.go`: `SafePath(p string) (string, error)` — resolve relative to working directory, reject paths escaping the workspace

## 4. Tool Implementations

- [x] 4.1 Implement `internal/tools/bash.go`: execute shell commands via `bash -c` with 120s timeout, dangerous command blacklist, 50000 char output truncation
- [x] 4.2 Implement `internal/tools/read.go`: read file with optional line limit, path safety check, 50000 char truncation
- [x] 4.3 Implement `internal/tools/write.go`: write file with auto parent directory creation, path safety check
- [x] 4.4 Implement `internal/tools/edit.go`: find and replace first occurrence of old text with new text, path safety check, error if text not found
- [x] 4.5 Implement `internal/tools/todo.go`: TodoManager with max 20 items, single in_progress enforcement, auto-id generation, status markers rendering

## 5. Tool Registry

- [x] 5.1 Implement `internal/tools/registry.go`: define `ToolDef` and `ToolHandler` types, register all 8 tools with OpenAI function-calling JSON Schema, create parent (8 tools) and child (6 tools, no task/todo) tool sets, build handler dispatch map

## 6. Sessions: Subagent

- [x] 6.1 Implement `internal/sessions/subagent.go`: subagent loop with independent message history, child tool set, max 30 rounds, fixed system prompt, return final text summary only

## 7. Sessions: Skill Loading

- [x] 7.1 Implement `internal/sessions/skill.go`: scan `skills/` directory for `.md` files, parse YAML frontmatter (`---` delimited), build skill index
- [x] 7.2 Implement `GetDescriptions()` for system prompt injection and `GetContent(name)` for on-demand loading with `<skill>` tag wrapping

## 8. Sessions: Context Compaction

- [x] 8.1 Implement `internal/sessions/compact.go`: `estimateTokens` (JSON length / 4), `microCompact` (replace old tool results with placeholders, keep last 3), `shouldCompact` (threshold 50000)
- [x] 8.2 Implement `autoCompact`: save full transcript to `.transcripts/transcript_<timestamp>.jsonl`, call LLM for summary, return compressed messages

## 9. Agent Loop

- [x] 9.1 Implement `internal/agent/loop.go`: core agent loop with tool-calling while loop, todo reminder injection (3 rounds), integration of all three compaction layers (micro every round, auto on threshold, manual on compact tool call)

## 10. CLI REPL

- [x] 10.1 Implement `internal/agent/repl.go`: stdin/stdout REPL loop, build system prompt with skill descriptions, call agent_loop per user input, handle exit/quit commands

## 11. Integration & Polish

- [x] 11.1 End-to-end manual test: start REPL, verify tool calling, subagent delegation, skill loading, and context compaction work correctly
- [x] 11.2 Verify edge cases: empty skills directory, missing config, path escape attempts, dangerous commands, timeout handling