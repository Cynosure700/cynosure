# memory-update Specification

## Purpose
TBD - created by archiving change add-memory-mechanism. Update Purpose after archive.
## Requirements
### Requirement: Agent can append to project memory file
The system SHALL provide an `update_memory` tool that accepts `scope: "project"`, `action: "append"`, and `content` parameters. When invoked, the system SHALL append the content to the end of the project-level `AGENTS.md` file, followed by a newline.

#### Scenario: Append new memory entry to file
- **WHEN** the agent calls `update_memory` with `scope: "project"`, `action: "append"`, and `content: "Always use snake_case for Go functions"`
- **THEN** the content is appended to `AGENTS.md` with a trailing newline

#### Scenario: Append to non-existent memory file
- **WHEN** the agent calls `update_memory` with `scope: "project"`, `action: "append"`, and `AGENTS.md` does not exist
- **THEN** a new `AGENTS.md` file is created with the content

### Requirement: Agent can replace project memory file entirely
The system SHALL provide an `update_memory` tool that accepts `scope: "project"`, `action: "replace"`, and `content` parameters. When invoked, the system SHALL overwrite the entire `AGENTS.md` file with the provided content.

#### Scenario: Replace entire memory file content
- **WHEN** the agent calls `update_memory` with `scope: "project"`, `action: "replace"`, and `content: "New project conventions:\n- Use tabs for indentation"`
- **THEN** `AGENTS.md` is overwritten with exactly the provided content

### Requirement: Agent can write session-level memory
The system SHALL provide an `update_memory` tool that accepts `scope: "session"` and `content` parameters. When invoked, the system SHALL store the content in session memory (in-memory only, cleared when session ends). Session memory SHALL be injected into subsequent conversation turns.

#### Scenario: Write session memory
- **WHEN** the agent calls `update_memory` with `scope: "session"` and `content: "User prefers snake_case naming"`
- **THEN** the content is stored in session memory and injected into subsequent conversation turns

#### Scenario: Multiple session memory writes
- **WHEN** the agent calls `update_memory` with `scope: "session"` multiple times
- **THEN** all entries are accumulated and injected together in subsequent turns

### Requirement: update_memory project scope enforces path safety
The system SHALL validate that the target file path resolves within the project working directory before performing any write operation for `scope: "project"`.

#### Scenario: Path traversal attempt
- **WHEN** the agent calls `update_memory` with `scope: "project"` and content that attempts to write outside the project directory
- **THEN** the operation is rejected and an error message is returned

### Requirement: update_memory only modifies project-level file
The system SHALL restrict `update_memory` with `scope: "project"` to only modify the project-level `AGENTS.md` file. User-level memory (`~/.link/AGENTS.md`) SHALL NOT be modifiable through this tool.

#### Scenario: Attempt to modify user memory
- **WHEN** the agent calls `update_memory` with a path outside the project directory
- **THEN** the operation is rejected with a safety error

### Requirement: update_memory returns confirmation
The system SHALL return a confirmation message after successfully updating memory, indicating the scope and action performed.

#### Scenario: Successful project append
- **WHEN** `update_memory` with `scope: "project"`, `action: "append"` completes successfully
- **THEN** the tool returns a message like "Memory updated: appended to AGENTS.md"

#### Scenario: Successful session write
- **WHEN** `update_memory` with `scope: "session"` completes successfully
- **THEN** the tool returns a message like "Session memory updated"

