## MODIFIED Requirements

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
