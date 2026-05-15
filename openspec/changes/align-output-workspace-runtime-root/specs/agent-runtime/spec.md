## MODIFIED Requirements

### Requirement: Runtime can invoke registered tools while serving web requests
The system SHALL allow the agent runtime to invoke registered tools during a web chat turn, SHALL return each tool execution result back into the runtime loop, and SHALL bind that execution to the single resolved runtime workspace root for the current service instance.

#### Scenario: Tool invocation completes within a chat turn
- **WHEN** the model requests a registered tool during a user's chat turn
- **THEN** the runtime executes the tool against the resolved runtime workspace, feeds the result back to the model, and continues the turn until a final assistant response is produced

### Requirement: Agent runtime loads enabled skills from the database for the current user
The system SHALL load the current user's enabled skills from persistent storage for each agent conversation turn, SHALL make them available through the runtime skill-loading mechanism, and SHALL also keep the shared builtin skill catalog loaded from the resolved runtime workspace available for that turn.

#### Scenario: Builtin skills are loaded from the resolved runtime workspace
- **WHEN** the service initializes and resolves a runtime workspace root
- **THEN** the runtime loads the shared builtin skill catalog from that workspace before serving chat turns

#### Scenario: Enabled user skills are merged with builtin runtime skills
- **WHEN** a user sends a message and has one or more enabled skills
- **THEN** the runtime makes both the enabled user skills and the builtin skills from the resolved runtime workspace available during that turn

#### Scenario: User without custom skills still receives an answer
- **WHEN** a user sends a message and has no enabled custom skills
- **THEN** the runtime answers using the shared builtin skills and base agent capabilities without failing

## ADDED Requirements

### Requirement: Runtime SHALL expose deployment-aware workspace paths to skills and tools
The system SHALL expose the resolved application home, runtime workspace root, command binary directory, and helper script directory to authorized skills and tools without requiring them to infer those locations from the process current working directory.

#### Scenario: Runtime environment reflects packaged deployment paths
- **WHEN** the active runtime workspace root is resolved to `<app-home>/output/workspace`
- **THEN** the runtime-provided environment for skills and tools references the packaged deployment paths under that workspace

#### Scenario: Runtime environment reflects local debug paths
- **WHEN** the active runtime workspace root falls back to `<app-home>/workspace`
- **THEN** the runtime-provided environment for skills and tools references the local debug paths under that workspace
