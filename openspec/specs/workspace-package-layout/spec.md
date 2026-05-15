# workspace-package-layout Specification

## Purpose
TBD - created by archiving change use-workspace-runtime-layout. Update Purpose after archive.
## Requirements
### Requirement: Deployment publishes a unified workspace runtime layout
The deployment artifact SHALL include a single runtime workspace root under the service package, and that workspace root SHALL contain the runtime skill catalog directory `workspace/skills`, the compiled command directory `workspace/bin`, and any approved helper-command directory required by the service.

#### Scenario: Deployment output includes required workspace asset directories
- **WHEN** the service build completes successfully
- **THEN** the deployment output contains `workspace/skills` and `workspace/bin` under the packaged workspace root

### Requirement: Standard build script packages workspace runtime assets
The repository SHALL provide a standard `build.sh` packaging script that cleans the output directory, compiles the main service binary, compiles supported `cmd` entrypoints into `output/workspace/bin`, and copies builtin skills plus required bootstrap and config assets into the deployment output.

#### Scenario: Running build script produces deployable workspace layout
- **WHEN** an operator runs `build.sh` from the repository root
- **THEN** the script produces a deployment output tree containing the main service binary and the packaged workspace runtime assets in their expected directories

### Requirement: Runtime command binaries are built before deployment startup
The system SHALL require runtime command binaries referenced by skills or authorized tools to be built during packaging and placed in `workspace/bin` before the deployed service starts serving traffic.

#### Scenario: Startup uses prebuilt workspace command artifacts
- **WHEN** the deployed service initializes runtime dependencies
- **THEN** it resolves command binaries from `workspace/bin` and does not attempt to compile `cmd` sources on demand

