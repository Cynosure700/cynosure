# deployment-runtime-layout Specification

## Purpose
TBD - created by archiving change align-output-workspace-runtime-root. Update Purpose after archive.
## Requirements
### Requirement: Runtime workspace root SHALL prefer packaged deployment assets
The system SHALL resolve a single runtime workspace root for the go-agent service and SHALL prefer the packaged deployment workspace under the service package when both packaged and source workspaces are present. If no explicit workspace root is configured and `output/workspace` is available under the application root, the runtime SHALL use that directory; otherwise it SHALL fall back to the source `workspace` directory for local debugging.

#### Scenario: Deployment workspace is preferred when packaged assets exist
- **WHEN** the service starts with both `<app-home>/output/workspace` and `<app-home>/workspace` present and no explicit workspace override is configured
- **THEN** the runtime resolves `<app-home>/output/workspace` as the active workspace root

#### Scenario: Source workspace is used for local debugging when deployment workspace is absent
- **WHEN** the service starts without a packaged `output/workspace` directory and no explicit workspace override is configured
- **THEN** the runtime resolves `<app-home>/workspace` as the active workspace root

### Requirement: Runtime asset directories SHALL derive from the resolved workspace root
The system SHALL derive builtin skill, command binary, and helper script directories from the resolved runtime workspace root unless an explicit override is configured.

#### Scenario: Skills and command directories follow deployment workspace
- **WHEN** the resolved runtime workspace root is `<app-home>/output/workspace`
- **THEN** builtin skills resolve from `<app-home>/output/workspace/skills`, command binaries resolve from `<app-home>/output/workspace/bin`, and helper scripts resolve from `<app-home>/output/workspace/cmd`

#### Scenario: Skills and command directories follow source workspace
- **WHEN** the resolved runtime workspace root is `<app-home>/workspace`
- **THEN** builtin skills resolve from `<app-home>/workspace/skills`, command binaries resolve from `<app-home>/workspace/bin`, and helper scripts resolve from `<app-home>/workspace/cmd`

