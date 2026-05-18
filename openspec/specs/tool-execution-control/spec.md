# tool-execution-control Specification

## Purpose
Define the authorization, isolation, validation, and audit rules that govern the reduced tool surface exposed by the web agent platform.

## Requirements
### Requirement: Runtime only executes tools that are registered and authorized
The system SHALL expose and execute only the tool set explicitly retained and registered for the web service, SHALL not execute historical CLI or REPL-oriented tools by default, and SHALL allow authorized web-agent tools to execute approved binaries or scripts from the resolved runtime workspace assets only when those tools are part of the supported web runtime contract.

#### Scenario: Registered web-service tool allowed
- **WHEN** the model requests a tool that is registered and permitted for the current web runtime
- **THEN** the system executes that tool with the configured workspace runtime paths and platform controls applied

#### Scenario: Default web runtime omits CLI-only tools
- **WHEN** the web service starts under the refactored default configuration
- **THEN** browser requests are served without exposing deprecated CLI or REPL-only tools as part of the default runtime contract

#### Scenario: Unknown tool rejected
- **WHEN** the model requests a tool that is not registered by the platform
- **THEN** the system rejects the request and does not execute arbitrary code or commands

### Requirement: Tool execution is restricted by user scope
The system SHALL continue to resolve terminal and file tool paths relative to the single resolved runtime workspace root, SHALL use that root as the default current working directory when appropriate, and SHALL reject tool requests that escape that workspace boundary regardless of whether the service is running in packaged deployment mode or local debug mode.

#### Scenario: Deployment workspace becomes the default tool cwd
- **WHEN** the resolved runtime workspace root is `<app-home>/output/workspace` and a tool request relies on the default working directory
- **THEN** the system executes the tool with `<app-home>/output/workspace` as its default cwd

#### Scenario: Local debug workspace becomes the default tool cwd
- **WHEN** the resolved runtime workspace root is `<app-home>/workspace` and a tool request relies on the default working directory
- **THEN** the system executes the tool with `<app-home>/workspace` as its default cwd

#### Scenario: Absolute path outside the active workspace is rejected
- **WHEN** a tool request targets an absolute path outside the resolved runtime workspace root
- **THEN** the system rejects the execution request

### Requirement: Tool execution is logged for auditability
The system SHALL continue to record each attempted or completed retained tool execution with enough metadata to support debugging and audit review, including the user identity, conversation identity, resolved working directory, tool name, status, and outcome summary or denial reason.

#### Scenario: Successful execution logged
- **WHEN** a retained tool completes successfully
- **THEN** the system stores an audit record with the user identity, conversation identity, resolved working directory, tool name, and outcome summary

#### Scenario: Rejected execution logged
- **WHEN** a retained tool execution is rejected by policy or validation
- **THEN** the system stores an audit record showing that the tool request was denied and why it was denied

### Requirement: Deployment command resources are validated before use
The system SHALL validate that deployment-provided command resources referenced by authorized web-agent tools exist in the configured command artifact roots before attempting execution.

#### Scenario: Missing binary artifact is rejected before execution
- **WHEN** an authorized tool references a deployment-provided binary path that does not exist in the configured artifact directory
- **THEN** the system rejects the execution request and records the missing artifact reason in audit metadata

#### Scenario: Missing script artifact is rejected before execution
- **WHEN** an authorized tool references a deployment-provided script path that does not exist in the configured artifact directory
- **THEN** the system rejects the execution request and records the missing artifact reason in audit metadata

### Requirement: Authorized tools can resolve packaged commands from workspace paths
The system SHALL allow authorized tools to invoke approved binaries from the resolved `workspace/bin` directory and approved helper scripts from the resolved `workspace/cmd` directory using stable runtime paths.

#### Scenario: Packaged command binary is invoked through workspace path
- **WHEN** an authorized tool call references a command binary published in the resolved runtime workspace
- **THEN** the system resolves and executes that binary from the active workspace `bin` path

#### Scenario: Packaged helper script is invoked through workspace path
- **WHEN** an authorized tool call references a helper script published in the resolved runtime workspace
- **THEN** the system resolves and executes that script from the active workspace `cmd` path

### Requirement: Tool execution SHALL use deployment-aware command asset resolution
The system SHALL resolve approved command binaries and helper scripts from directories derived from the resolved runtime workspace root instead of relying on CLI entrypoints or ad hoc process cwd assumptions.

#### Scenario: Packaged command binary is resolved from deployment workspace
- **WHEN** the active runtime workspace root is `<app-home>/output/workspace` and an authorized tool invokes a packaged command
- **THEN** the system resolves that command from `<app-home>/output/workspace/bin` or `<app-home>/output/workspace/cmd`

#### Scenario: Local debug command binary is resolved from source workspace
- **WHEN** the active runtime workspace root is `<app-home>/workspace` and an authorized tool invokes a packaged command
- **THEN** the system resolves that command from `<app-home>/workspace/bin` or `<app-home>/workspace/cmd`

### Requirement: General chat mode exposes only a safe default tool set
The system SHALL expose only the subset of registered tools that are compatible with a browser-first chatbot product, and SHALL exclude shell execution, arbitrary file mutation, directory browsing, and other tools that assume access to a user's local workspace.

#### Scenario: Standard chat session omits shell and local workspace tools
- **WHEN** a user interacts with the default browser chat product
- **THEN** the runtime does not expose shell execution, arbitrary file mutation, directory browsing, or similar local-workspace tools

### Requirement: Restricted tool availability is communicated gracefully
The system SHALL handle requests for unavailable browser-incompatible tools by responding with a clear capability boundary while continuing to assist conversationally where possible.

#### Scenario: User asks for action requiring unsupported local operation
- **WHEN** a request would require shell execution, local file mutation, or access to a user's local directory
- **THEN** the assistant explains that the browser product does not support that operation and provides the best available non-local help instead of failing silently
