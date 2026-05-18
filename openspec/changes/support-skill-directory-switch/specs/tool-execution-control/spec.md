## MODIFIED Requirements

### Requirement: Tool execution is restricted by user scope
The system SHALL continue to resolve terminal and file tool paths relative to the single resolved runtime workspace root, SHALL use that root as the default current working directory when no active skill directory is available, SHALL allow an active filesystem-backed skill directory inside that workspace to override the default current working directory for subsequent skill-driven tool calls, and SHALL reject tool requests that escape that workspace boundary regardless of whether the service is running in packaged deployment mode or local debug mode.

#### Scenario: Deployment workspace becomes the default tool cwd
- **WHEN** the resolved runtime workspace root is `<app-home>/output/workspace` and a tool request relies on the default working directory without any active skill directory
- **THEN** the system executes the tool with `<app-home>/output/workspace` as its default cwd

#### Scenario: Local debug workspace becomes the default tool cwd
- **WHEN** the resolved runtime workspace root is `<app-home>/workspace` and a tool request relies on the default working directory without any active skill directory
- **THEN** the system executes the tool with `<app-home>/workspace` as its default cwd

#### Scenario: Filesystem-backed skill directory overrides the default tool cwd
- **WHEN** a loaded skill was sourced from `<workspace-root>/skills/<skill-name>/SKILL.md` and a subsequent tool request relies on the default working directory
- **THEN** the system executes the tool with `<workspace-root>/skills/<skill-name>` as its default cwd

#### Scenario: Skill without a local directory falls back to workspace root
- **WHEN** the active skill was loaded from a non-filesystem source and a subsequent tool request relies on the default working directory
- **THEN** the system executes the tool with the resolved runtime workspace root as its default cwd

#### Scenario: Absolute path outside the active workspace is rejected
- **WHEN** a tool request targets an absolute path outside the resolved runtime workspace root
- **THEN** the system rejects the execution request

### Requirement: Tool execution is logged for auditability
The system SHALL continue to record each attempted or completed retained tool execution with enough metadata to support debugging and audit review, including the user identity, conversation identity, resolved working directory, tool name, status, and outcome summary or denial reason, and the resolved working directory SHALL reflect the active skill directory when a skill-driven tool call changes the default cwd.

#### Scenario: Successful execution logged
- **WHEN** a retained tool completes successfully
- **THEN** the system stores an audit record with the user identity, conversation identity, resolved working directory, tool name, and outcome summary

#### Scenario: Rejected execution logged
- **WHEN** a retained tool execution is rejected by policy or validation
- **THEN** the system stores an audit record showing that the tool request was denied and why it was denied

#### Scenario: Skill-driven cwd is visible in audit metadata
- **WHEN** a retained tool executes using an active filesystem-backed skill directory as its default cwd
- **THEN** the stored audit record includes that skill directory as the resolved working directory
