## MODIFIED Requirements

### Requirement: Runtime only executes tools that are registered and authorized
The system SHALL execute only tools that are explicitly registered by the platform, authorized for the current web runtime environment, and bound to the configured deployment workspace environment.

#### Scenario: Registered tool allowed
- **WHEN** the model requests a tool that is registered and permitted for the current runtime
- **THEN** the system executes that tool with the configured workspace runtime paths and platform controls applied

#### Scenario: Unknown tool rejected
- **WHEN** the model requests a tool that is not registered by the platform
- **THEN** the system rejects the request and does not execute arbitrary code or commands

### Requirement: Tool execution is restricted by user scope
The system SHALL resolve terminal and file tool paths relative to the configured deployment workspace and SHALL reject tool requests that escape that workspace boundary or directly target host-user filesystem locations outside the allowed workspace scope.

#### Scenario: Relative path request uses deployment workspace
- **WHEN** a tool request contains a relative path or relies on the default working directory
- **THEN** the system resolves that path from the configured deployment workspace before executing the tool

#### Scenario: Host filesystem path outside workspace is rejected
- **WHEN** a tool request targets an absolute path outside the configured deployment workspace
- **THEN** the system rejects the execution request

## ADDED Requirements

### Requirement: Authorized tools can resolve packaged commands from workspace paths
The system SHALL allow authorized skills and tools to invoke approved binaries from `workspace/bin` and approved helper scripts from the packaged workspace command directories using stable deployment paths instead of relying on the host PATH or process current working directory.

#### Scenario: Packaged command binary is invoked through workspace path
- **WHEN** an authorized tool call references a command binary published in `workspace/bin`
- **THEN** the system resolves and executes that binary from the packaged workspace path
