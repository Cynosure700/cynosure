## MODIFIED Requirements

### Requirement: Agent runtime loads enabled skills from the database for the current user
The system SHALL construct the runtime skill set for each agent conversation turn from two sources: builtin skills loaded from the deployed `workspace/skills` catalog at service startup, and the current user's enabled custom skills from persistent storage. The merged skill set SHALL be made available to the model through the runtime skill-loading mechanism, and builtin skills SHALL remain available even when the user has no enabled custom skills.

#### Scenario: Builtin workspace skills are injected for every user
- **WHEN** the service has loaded one or more builtin skills from `workspace/skills` and an authenticated user sends a message
- **THEN** the runtime includes those builtin skills in the conversation turn

#### Scenario: Enabled user skills are merged with workspace builtin skills
- **WHEN** a user sends a message and has one or more enabled custom skills
- **THEN** the runtime includes those enabled custom skills alongside the builtin skills loaded from `workspace/skills`

#### Scenario: User without custom skills still receives an answer
- **WHEN** a user sends a message and has no enabled custom skills
- **THEN** the runtime answers using the builtin workspace skills and base agent capabilities without failing

### Requirement: Runtime can invoke registered tools while serving web requests
The system SHALL allow the agent runtime to invoke only registered tools that are enabled for the web deployment, and SHALL execute those tools with the deployed workspace environment injected so that tool calls can resolve `workspace/bin`, `workspace/cmd`, and the workspace working directory deterministically.

#### Scenario: Tool invocation completes within a chat turn
- **WHEN** the model requests a registered tool that is enabled for the current web runtime
- **THEN** the runtime executes the tool with the configured workspace environment, feeds the result back to the model, and continues the turn until a final assistant response is produced

## ADDED Requirements

### Requirement: Runtime validates packaged workspace assets before serving requests
The system SHALL resolve builtin skill and command artifact locations from the deployed workspace root during service startup and SHALL fail initialization if required runtime asset directories cannot be found or prepared.

#### Scenario: Startup checks packaged skill and command directories
- **WHEN** the web service starts in a deployed environment
- **THEN** it validates the configured `workspace/skills` and command artifact directories before accepting requests
