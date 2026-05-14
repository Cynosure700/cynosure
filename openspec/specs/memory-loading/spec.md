# memory-loading Specification

## Purpose
TBD - created by archiving change add-memory-mechanism. Update Purpose after archive.
## Requirements
### Requirement: Agent loads project-level memory on startup
The system SHALL automatically load the `AGENTS.md` file from the project root directory at startup. If the file exists, its content SHALL be injected into the system prompt within `<project_memory>` tags. If the file does not exist, the system SHALL proceed without error.

#### Scenario: Project memory file exists
- **WHEN** the agent starts and `AGENTS.md` exists in the project root directory
- **THEN** the file content is read and injected into the system prompt as `<project_memory>...content...</project_memory>`

#### Scenario: Project memory file does not exist
- **WHEN** the agent starts and `AGENTS.md` does not exist in the project root directory
- **THEN** the system prompt is built without project memory, and no error is raised

#### Scenario: Project memory file is empty
- **WHEN** the agent starts and `AGENTS.md` exists but is empty
- **THEN** the system prompt includes `<project_memory></project_memory>` with no content

### Requirement: Agent loads user-level memory on startup
The system SHALL automatically load the `AGENTS.md` file from `~/.link/AGENTS.md` at startup. If the file exists, its content SHALL be injected into the system prompt within `<user_memory>` tags. If the file does not exist, the system SHALL proceed without error.

#### Scenario: User memory file exists
- **WHEN** the agent starts and `~/.link/AGENTS.md` exists
- **THEN** the file content is read and injected into the system prompt as `<user_memory>...content...</user_memory>`

#### Scenario: User memory file does not exist
- **WHEN** the agent starts and `~/.link/AGENTS.md` does not exist
- **THEN** the system prompt is built without user memory, and no error is raised

### Requirement: Memory content is placed at the end of system prompt
The system SHALL place persistent memory content after all other system prompt content (cwd, skill descriptions), ensuring the model sees memory instructions last for maximum attention.

#### Scenario: Both memory files exist
- **WHEN** both project and user memory files exist
- **THEN** the system prompt ends with `<project_memory>` block followed by `<user_memory>` block

#### Scenario: Only project memory exists
- **WHEN** only `AGENTS.md` exists
- **THEN** the system prompt ends with `<project_memory>` block only

### Requirement: Memory loading is read-only at startup
The system SHALL NOT modify memory files during the loading phase. Memory files are only modified through the `update_memory` tool during an active session.

#### Scenario: Memory file is read but not written
- **WHEN** the agent loads memory files at startup
- **THEN** the files are opened in read-only mode and never written to during loading

