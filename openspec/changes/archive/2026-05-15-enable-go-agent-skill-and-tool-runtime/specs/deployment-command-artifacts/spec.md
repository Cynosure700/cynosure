## ADDED Requirements

### Requirement: Deployment defines a fixed application root for runtime assets
The system SHALL define a fixed application root for each deployed go-agent instance and SHALL resolve builtin skills, deployment command resources, and runtime data directories relative to that root instead of relying on the process current working directory.

#### Scenario: Service startup resolves runtime asset roots from application root
- **WHEN** the service starts in a deployed environment
- **THEN** it resolves the builtin skills directory, deployment command roots, and workspace root from the configured application root before accepting requests

### Requirement: Go commands under cmd are built into deployment binaries
The system SHALL compile supported Go command entrypoints from the source `cmd` tree during the deployment build process and SHALL publish the resulting binaries into a stable deployment artifact directory.

#### Scenario: Deployment build publishes compiled command binary
- **WHEN** the deployment pipeline processes a supported Go command entrypoint under the source `cmd` tree
- **THEN** it produces a binary artifact in the configured deployment binary directory for runtime use

### Requirement: Script helpers under cmd are published as runtime assets
The system SHALL publish script-based helper files, including `.py` files, from the source `cmd` tree into a stable deployment script directory when those files are part of the approved runtime command set.

#### Scenario: Deployment build publishes approved script helper
- **WHEN** the deployment pipeline processes an approved script helper under the source `cmd` tree
- **THEN** it copies that script into the configured deployment command asset directory for runtime use

### Requirement: Deployment command artifacts are read-only runtime resources
The system SHALL treat deployment-provided binaries and scripts as read-only runtime resources that are managed by the platform build and release process rather than by end users.

#### Scenario: User cannot mutate deployed command artifact
- **WHEN** a runtime request attempts to overwrite or delete a deployment-provided binary or script
- **THEN** the system rejects the operation because deployment command artifacts are platform-managed read-only resources
