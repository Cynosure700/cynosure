## MODIFIED Requirements

### Requirement: Deployment publishes a unified workspace runtime layout
The deployment artifact SHALL publish the web-service binary together with the runtime assets directly required by that service, including the deployment-aware workspace layout used by the web agent for command execution such as `workspace/bin` and `workspace/cmd` when those capabilities are supported.

#### Scenario: Minimal deployment package starts the service with required runtime assets
- **WHEN** the build or release process publishes the refactored web service
- **THEN** the resulting package contains the server binary plus the configuration, migration, static assets, and runtime workspace assets required by the supported web-agent capabilities

### Requirement: Standard build script packages workspace runtime assets
The repository SHALL provide at most one standard build/deploy path for the web service, and that path SHALL build the single supported server startup target while continuing to package the runtime workspace assets required by the web agent without generating duplicate web-server binaries.

#### Scenario: Standard build path produces server artifact and runtime assets
- **WHEN** an operator runs the documented build path for the service
- **THEN** the build output contains the single supported web-service executable together with the required runtime workspace assets under the expected workspace layout

### Requirement: Runtime command binaries are built before deployment startup
The system SHALL require runtime command binaries referenced by supported web-agent capabilities to be built during packaging and placed in the runtime workspace before the deployed service starts serving traffic.

#### Scenario: Startup uses prebuilt workspace command artifacts
- **WHEN** the deployed service initializes runtime dependencies for command-capable tools
- **THEN** it resolves command binaries from the packaged workspace asset directories and does not compile those command sources on demand during request handling
