## MODIFIED Requirements

### Requirement: Runtime workspace root SHALL prefer packaged deployment assets
The system SHALL continue to resolve a deployment-aware runtime workspace root for the web service, SHALL prefer packaged deployment assets under `output/workspace` when present, and SHALL fall back to the source `workspace` layout for local debugging when packaged assets are absent.

#### Scenario: Local source run uses source workspace assets
- **WHEN** a developer starts the service from the repository checkout without explicit path overrides and no packaged deployment workspace is present
- **THEN** the service resolves the source `workspace` layout as the active runtime workspace root

#### Scenario: Deployed binary run prefers packaged workspace assets
- **WHEN** an operator starts the deployed web-service binary from its installation directory and `output/workspace` is present
- **THEN** the service resolves `output/workspace` as the active runtime workspace root for command binaries, scripts, and other runtime assets

### Requirement: Runtime asset directories SHALL derive from the resolved workspace root
The system SHALL derive the runtime asset directories required by the web agent from the resolved workspace root, including `workspace/bin` and `workspace/cmd`, while continuing to derive core service paths from the application root or explicit configuration.

#### Scenario: Core service directories come from application root
- **WHEN** the service initializes its configuration, logs, or persistent data paths
- **THEN** those paths derive from the application root or explicit configuration rather than from the runtime workspace root

#### Scenario: Command asset directories come from runtime workspace root
- **WHEN** the web agent needs runtime command binaries or helper scripts
- **THEN** the service resolves those assets from the active runtime workspace directories under `bin` and `cmd`
