## ADDED Requirements

### Requirement: Session memory is injected into each conversation turn
The system SHALL inject the current session memory content into every conversation turn after the user message, wrapped in `<session_memory>` tags. If session memory is empty, no injection SHALL occur.

#### Scenario: Session memory with content
- **WHEN** session memory contains entries and the user sends a message
- **THEN** a `<session_memory>...content...</session_memory>` block is appended after the user message before sending to the LLM

#### Scenario: Empty session memory
- **WHEN** session memory is empty and the user sends a message
- **THEN** no `<session_memory>` block is injected

### Requirement: Session memory is cleared when REPL exits
The system SHALL clear all session memory entries when the REPL session ends (user types exit/quit or process terminates).

#### Scenario: Normal exit clears memory
- **WHEN** the user types "exit" or "quit" in the REPL
- **THEN** all session memory entries are discarded

#### Scenario: Process termination clears memory
- **WHEN** the go-agent process is terminated (SIGINT, SIGTERM)
- **THEN** all session memory entries are discarded (in-memory only, no persistence)

### Requirement: Session memory survives context compaction
The system SHALL preserve session memory content across context compaction operations. Session memory SHALL NOT be removed or summarized during micro, auto, or manual compaction.

#### Scenario: Auto compaction preserves session memory
- **WHEN** auto compaction is triggered and session memory exists
- **THEN** the session memory entries remain intact and continue to be injected in subsequent turns

### Requirement: Session memory is not available to subagents
The system SHALL NOT inject session memory into subagent conversations. Subagents operate with independent context.

#### Scenario: Subagent does not receive session memory
- **WHEN** a subagent is spawned via the `task` tool
- **THEN** the subagent's messages do not include `<session_memory>` blocks