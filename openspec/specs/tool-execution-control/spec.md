# tool-execution-control Specification

## Purpose
Define the authorization, isolation, and audit rules that govern tool execution in the web agent platform.
## Requirements
### Requirement: Runtime only executes tools that are registered and authorized
The system SHALL execute only tools that are explicitly registered by the platform, authorized for the current web runtime environment, and bound to the configured deployment workspace environment.

#### Scenario: Registered tool allowed
- **WHEN** the model requests a tool that is registered and permitted for the current runtime
- **THEN** the system executes that tool with the configured workspace runtime paths and platform controls applied

#### Scenario: Unknown tool rejected
- **WHEN** the model requests a tool that is not registered by the platform
- **THEN** the system rejects the request and does not execute arbitrary code or commands

### Requirement: Tool execution is restricted by user scope
The system SHALL resolve terminal and file tool paths relative to the single resolved runtime workspace root, SHALL use that root as the default current working directory, and SHALL reject tool requests that escape that workspace boundary regardless of whether the service is running in packaged deployment mode or local debug mode.

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
The system SHALL record each attempted or completed tool execution with enough metadata to support debugging and audit review, including the user identity, conversation identity, resolved working directory, tool name, status, and outcome summary or denial reason.

#### Scenario: Successful execution logged
- **WHEN** a tool completes successfully
- **THEN** the system stores an audit record with the user identity, conversation identity, resolved working directory, tool name, and outcome summary

#### Scenario: Rejected execution logged
- **WHEN** a tool execution is rejected by policy or validation
- **THEN** the system stores an audit record showing that the tool request was denied and why it was denied

### Requirement: Deployment command resources are validated before use
The system SHALL validate that deployment-provided command resources referenced by skills or tools exist in the configured command artifact roots before attempting execution.

#### Scenario: Missing binary artifact is rejected before execution
- **WHEN** a skill or tool references a deployment-provided binary path that does not exist in the configured artifact directory
- **THEN** the system rejects the execution request and records the missing artifact reason in audit metadata

#### Scenario: Missing script artifact is rejected before execution
- **WHEN** a skill or tool references a deployment-provided script path that does not exist in the configured artifact directory
- **THEN** the system rejects the execution request and records the missing artifact reason in audit metadata

### Requirement: Authorized tools can resolve packaged commands from workspace paths
The system SHALL allow authorized skills and tools to invoke approved binaries from `workspace/bin` and approved helper scripts from the packaged workspace command directories using stable deployment paths instead of relying on the host PATH or process current working directory.

#### Scenario: Packaged command binary is invoked through workspace path
- **WHEN** an authorized tool call references a command binary published in `workspace/bin`
- **THEN** the system resolves and executes that binary from the packaged workspace path

### Requirement: Tool execution SHALL use deployment-aware command asset resolution
The system SHALL resolve approved command binaries and helper scripts from directories derived from the resolved runtime workspace root instead of relying on the process cwd or host PATH.

#### Scenario: Packaged command binary is resolved from deployment workspace
- **WHEN** the active runtime workspace root is `<app-home>/output/workspace` and an authorized tool or skill invokes a packaged command
- **THEN** the system resolves that command from `<app-home>/output/workspace/bin` or `<app-home>/output/workspace/cmd`

#### Scenario: Local debug command binary is resolved from source workspace
- **WHEN** the active runtime workspace root is `<app-home>/workspace` and an authorized tool or skill invokes a packaged command
- **THEN** the system resolves that command from `<app-home>/workspace/bin` or `<app-home>/workspace/cmd`

