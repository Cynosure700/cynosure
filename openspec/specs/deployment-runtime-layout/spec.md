# deployment-runtime-layout Specification

## Purpose
Define how the web service resolves its deployment-aware runtime workspace root and separates application-root service paths from workspace-derived runtime command assets.

## Requirements
### Requirement: Runtime workspace root SHALL prefer packaged deployment assets
The system SHALL continue to resolve a deployment-aware runtime workspace root for the web service, SHALL prefer packaged deployment assets under `output/workspace` when present, and SHALL fall back to the source `workspace` layout for local debugging when packaged assets are absent.

#### Scenario: Local source run uses source workspace assets
- **WHEN** a developer starts the service from the repository checkout without explicit path overrides and no packaged deployment workspace is present
- **THEN** the service resolves the source `workspace` layout as the active runtime workspace root

#### Scenario: Deployed binary run prefers packaged workspace assets
- **WHEN** an operator starts the deployed web-service binary from its installation directory and `output/workspace` is present
- **THEN** the service resolves `output/workspace` as the active runtime workspace root for command binaries, scripts, and other runtime assets

### Requirement: Runtime asset directories SHALL derive from the resolved workspace root
The system SHALL derive the runtime asset directories required by the web agent from the resolved workspace root, including `workspace/skills`, `workspace/bin`, and `workspace/cmd`, while continuing to derive core service paths from the application root or explicit configuration. When explicit configuration is provided for those runtime asset directories, the configured values SHALL still resolve to the canonical directories under the active workspace root rather than introducing alternate runtime asset roots.

#### Scenario: Core service directories come from application root
- **WHEN** the service initializes its configuration, logs, or persistent data paths
- **THEN** those paths derive from the application root or explicit configuration rather than from the runtime workspace root

#### Scenario: Command asset directories come from runtime workspace root
- **WHEN** the web agent needs runtime command binaries or helper scripts
- **THEN** the service resolves those assets from the active runtime workspace directories under `bin` and `cmd`

#### Scenario: Explicit runtime asset overrides stay canonical
- **WHEN** an operator explicitly configures builtin skill, command binary, or command script directories for the web service
- **THEN** each configured directory resolves to the canonical `skills`, `bin`, or `cmd` path under the active runtime workspace root
