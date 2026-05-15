## ADDED Requirements

### Requirement: Each user receives an isolated server workspace
The system SHALL provision a filesystem workspace for each user under the configured workspace root and SHALL ensure that different users do not share the same workspace path.

#### Scenario: Workspace is provisioned for a user
- **WHEN** a user starts a chat turn or tool execution that requires a workspace and no workspace exists for that user yet
- **THEN** the system creates that user's workspace directory under the configured workspace root

#### Scenario: Different users resolve different workspace roots
- **WHEN** two different authenticated users access the runtime
- **THEN** the system resolves two different workspace directories and does not map both users to the same filesystem root

### Requirement: Agent turns run with a deterministic default working directory
The system SHALL resolve a default working directory for each Web conversation under the current user's workspace and SHALL use that directory when executing tools that rely on relative paths or shell commands.

#### Scenario: New conversation uses user workspace root as default cwd
- **WHEN** a new conversation starts without an explicit working directory override
- **THEN** the runtime uses that user's workspace root as the default working directory for tool execution

#### Scenario: Relative path tool request uses conversation cwd
- **WHEN** a tool request contains a relative path
- **THEN** the system resolves that path relative to the conversation's default working directory before applying further validation

## ADDED Requirements

### Requirement: Deployment command resources are separated from user workspaces
The system SHALL keep deployment-provided command resources under a fixed application root that is separate from user workspaces, and SHALL treat those command resources as read-only runtime assets rather than user-owned files.

#### Scenario: Shared command artifact path does not overlap user workspace
- **WHEN** the system resolves the path of a deployment-provided binary or script
- **THEN** it resolves that path under the configured application root and not under the current user's writable workspace

#### Scenario: User tool execution cannot overwrite shared command resources
- **WHEN** a user-triggered tool call attempts to write into the shared deployment command directory
- **THEN** the system rejects the request because deployment command resources are read-only assets
