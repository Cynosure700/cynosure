## MODIFIED Requirements

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

## ADDED Requirements

### Requirement: Tool execution SHALL use deployment-aware command asset resolution
The system SHALL resolve approved command binaries and helper scripts from directories derived from the resolved runtime workspace root instead of relying on the process cwd or host PATH.

#### Scenario: Packaged command binary is resolved from deployment workspace
- **WHEN** the active runtime workspace root is `<app-home>/output/workspace` and an authorized tool or skill invokes a packaged command
- **THEN** the system resolves that command from `<app-home>/output/workspace/bin` or `<app-home>/output/workspace/cmd`

#### Scenario: Local debug command binary is resolved from source workspace
- **WHEN** the active runtime workspace root is `<app-home>/workspace` and an authorized tool or skill invokes a packaged command
- **THEN** the system resolves that command from `<app-home>/workspace/bin` or `<app-home>/workspace/cmd`
