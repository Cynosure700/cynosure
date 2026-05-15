## MODIFIED Requirements

### Requirement: Runtime only executes tools that are registered and authorized
The system SHALL execute only tools that are present in the platform registry and explicitly enabled for the current Web runtime deployment. Tools that are not registered or not enabled for the Web deployment MUST be rejected without execution.

#### Scenario: Registered tool allowed
- **WHEN** the model requests a tool that is registered and permitted for the current Web runtime
- **THEN** the system executes that tool under platform control

#### Scenario: Unknown tool rejected
- **WHEN** the model requests a tool that is not registered by the platform
- **THEN** the system rejects the request and does not execute arbitrary code or commands

#### Scenario: Registered but disabled Web tool rejected
- **WHEN** the model requests a tool that exists in the platform registry but is not enabled for the current Web deployment
- **THEN** the system rejects the request without invoking that tool handler

### Requirement: Tool execution is restricted by user scope
The system SHALL restrict tool execution so that every file path, shell working directory, skill reference, conversation resource, or persisted artifact reference resolves within the current user's allowed scope and current conversation context.

#### Scenario: User workspace isolation enforced
- **WHEN** a tool request targets a file path or working directory outside the current user's allowed workspace
- **THEN** the system rejects the tool execution request

#### Scenario: Relative path uses conversation working directory
- **WHEN** a tool request includes a relative path or shell command that depends on a current directory
- **THEN** the system resolves it against the conversation's default working directory before checking whether it remains inside the user's allowed scope

#### Scenario: Deployment command path is allowed only from approved roots
- **WHEN** a tool or skill requests execution of a deployment-provided binary or script
- **THEN** the system permits the command only if its resolved path is under the configured read-only deployment command roots and the tool is authorized for that operation

### Requirement: Tool execution is logged for auditability
The system SHALL record each attempted or completed tool execution with enough metadata to support debugging and audit review, including the user identity, conversation identity, resolved working directory, tool name, status, and outcome summary or denial reason.

#### Scenario: Successful execution logged
- **WHEN** a tool completes successfully
- **THEN** the system stores an audit record with the user identity, conversation identity, resolved working directory, tool name, and outcome summary

#### Scenario: Rejected execution logged
- **WHEN** a tool execution is rejected by policy or validation
- **THEN** the system stores an audit record showing that the tool request was denied and why it was denied

## ADDED Requirements

### Requirement: Deployment command resources are validated before use
The system SHALL validate that deployment-provided command resources referenced by skills or tools exist in the configured command artifact roots before attempting execution.

#### Scenario: Missing binary artifact is rejected before execution
- **WHEN** a skill or tool references a deployment-provided binary path that does not exist in the configured artifact directory
- **THEN** the system rejects the execution request and records the missing artifact reason in audit metadata

#### Scenario: Missing script artifact is rejected before execution
- **WHEN** a skill or tool references a deployment-provided script path that does not exist in the configured artifact directory
- **THEN** the system rejects the execution request and records the missing artifact reason in audit metadata
