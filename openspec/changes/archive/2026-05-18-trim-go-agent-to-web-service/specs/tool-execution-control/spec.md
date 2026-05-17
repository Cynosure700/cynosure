## MODIFIED Requirements

### Requirement: Runtime only executes tools that are registered and authorized
The system SHALL expose only the tool set explicitly retained for the web service, SHALL not execute historical CLI/REPL-oriented tools by default, and SHALL continue to allow authorized web-agent tools to execute approved binaries or scripts from the resolved runtime workspace assets when those tools are part of the supported web runtime contract.

#### Scenario: Default web runtime omits CLI-only tools
- **WHEN** the web service starts under the refactored default configuration
- **THEN** browser requests are served without exposing deprecated CLI or REPL-only tools as part of the default runtime contract

#### Scenario: Authorized workspace-backed tool remains available
- **WHEN** an operator keeps a command-capable tool enabled after the refactor
- **THEN** that tool is available only if it is explicitly registered and authorized for the web service and it executes through the resolved runtime workspace assets

### Requirement: Tool execution is restricted by user scope
The system SHALL continue to resolve terminal and file tool paths relative to the single resolved runtime workspace root, SHALL use that root as the default current working directory when appropriate, and SHALL reject tool requests that escape that workspace boundary.

#### Scenario: Workspace-backed command uses resolved runtime root
- **WHEN** an authorized tool request relies on the default working directory
- **THEN** the system executes the tool with the resolved runtime workspace root as its default cwd

#### Scenario: Absolute path outside the active workspace is rejected
- **WHEN** a tool request targets an absolute path outside the resolved runtime workspace root
- **THEN** the system rejects the execution request

### Requirement: Tool execution is logged for auditability
The system SHALL continue to record audit information for retained tool execution in the web service, including the workspace-aware metadata that remains relevant after CLI/REPL-specific paths are removed.

#### Scenario: Retained tool execution is still auditable
- **WHEN** a retained web-service tool executes successfully or is rejected
- **THEN** the system stores an audit record with the supported metadata for that tool request, including workspace-related execution context when applicable

### Requirement: Deployment command resources are validated before use
The system SHALL validate deployment-provided command resources referenced by authorized web-agent tools before attempting execution.

#### Scenario: Missing binary artifact is rejected before execution
- **WHEN** a tool references a binary path that does not exist in the resolved runtime workspace artifact directory
- **THEN** the system rejects the execution request and records the missing artifact reason in audit metadata

#### Scenario: Missing script artifact is rejected before execution
- **WHEN** a tool references a script path that does not exist in the resolved runtime workspace artifact directory
- **THEN** the system rejects the execution request and records the missing artifact reason in audit metadata

### Requirement: Authorized tools can resolve packaged commands from workspace paths
The system SHALL allow authorized tools to invoke approved binaries from the resolved `workspace/bin` directory and approved helper scripts from the resolved `workspace/cmd` directory using stable runtime paths.

#### Scenario: Packaged command binary is invoked through workspace path
- **WHEN** an authorized tool call references a command binary published in the resolved runtime workspace
- **THEN** the system resolves and executes that binary from the active workspace `bin` path

### Requirement: Tool execution SHALL use deployment-aware command asset resolution
The system SHALL resolve approved command binaries and helper scripts from directories derived from the resolved runtime workspace root instead of relying on CLI entrypoints or ad hoc process cwd assumptions.

#### Scenario: Deployment workspace command is resolved from packaged assets
- **WHEN** the active runtime workspace root is the packaged deployment workspace and an authorized tool invokes a packaged command
- **THEN** the system resolves that command from the packaged workspace `bin` or `cmd` path

#### Scenario: Local debug command is resolved from source workspace
- **WHEN** the active runtime workspace root falls back to the source workspace and an authorized tool invokes a packaged command
- **THEN** the system resolves that command from the source workspace `bin` or `cmd` path
