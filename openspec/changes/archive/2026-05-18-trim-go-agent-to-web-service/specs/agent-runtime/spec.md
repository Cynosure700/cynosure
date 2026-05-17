## MODIFIED Requirements

### Requirement: Agent runtime loads enabled skills from the database for the current user
The system SHALL allow the web runtime to start and serve browser chat requests without depending on CLI REPL components, and SHALL continue to treat skill loading as runtime behavior that can coexist with the deployment-aware workspace layout used by the web service.

#### Scenario: Service starts with deployment-aware runtime workspace
- **WHEN** the web service starts and resolves its runtime workspace root
- **THEN** the runtime initializes successfully using the configured web-service dependencies and the resolved runtime asset layout

#### Scenario: User without enabled skills still receives an answer
- **WHEN** a user sends a message and no enabled skills are available for that turn
- **THEN** the runtime answers using the base web assistant behavior without failing startup or request handling

### Requirement: Runtime can invoke registered tools while serving web requests
The system SHALL complete browser chat turns without requiring CLI REPL interaction, and any authorized tool execution during a web request SHALL come only from tools explicitly registered by the web service and SHALL continue to use the resolved runtime workspace assets, including `workspace/bin` and `workspace/cmd`, when command execution is needed.

#### Scenario: Standard browser turn completes without CLI interaction
- **WHEN** the runtime handles a normal browser chat turn that does not need external side effects
- **THEN** it produces the assistant response without depending on CLI agent code or terminal interaction

#### Scenario: Authorized command uses runtime workspace assets
- **WHEN** the runtime invokes an authorized command-oriented tool during a browser chat turn
- **THEN** that tool is resolved from the web service's registered capability set and may execute binaries or scripts from the resolved `workspace/bin` or `workspace/cmd` paths

## ADDED Requirements

### Requirement: Runtime validates web-service dependencies together with required runtime workspace assets
The system SHALL validate the dependencies required to start and serve the web application, including configuration, storage connectivity, and the runtime workspace assets required by the enabled web agent capabilities, and SHALL not treat CLI-only modules as startup requirements.

#### Scenario: Required runtime assets are validated for command-capable service
- **WHEN** the configured web agent capabilities rely on binaries or scripts published in the runtime workspace
- **THEN** startup validates the resolved `workspace/bin` and `workspace/cmd` paths required by those capabilities before serving requests

#### Scenario: Missing CLI-only modules do not block startup
- **WHEN** historical CLI or REPL-only modules are absent from the service entry path
- **THEN** startup still proceeds as long as the required web-service configuration, storage dependencies, and runtime assets are valid
