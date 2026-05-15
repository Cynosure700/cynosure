# user-workspace Specification

## Purpose
TBD - created by archiving change enable-go-agent-skill-and-tool-runtime. Update Purpose after archive.
## Requirements
### Requirement: Deployment command resources are separated from user workspaces
The system SHALL keep deployment-provided command resources under a fixed application root that is separate from user workspaces, and SHALL treat those command resources as read-only runtime assets rather than user-owned files.

#### Scenario: Shared command artifact path does not overlap user workspace
- **WHEN** the system resolves the path of a deployment-provided binary or script
- **THEN** it resolves that path under the configured application root and not under the current user's writable workspace

#### Scenario: User tool execution cannot overwrite shared command resources
- **WHEN** a user-triggered tool call attempts to write into the shared deployment command directory
- **THEN** the system rejects the request because deployment command resources are read-only assets

