## MODIFIED Requirements

### Requirement: Agent runtime loads enabled skills from the database for the current user
The system SHALL load a merged runtime skill set for each agent conversation turn, consisting of the shared builtin skills from the configured `go-agent/skills` catalog and the current user's enabled custom skills from persistent storage. The merged skill set SHALL be made available to the model through the runtime skill-loading mechanism, and builtin skills SHALL remain available even when the user has no enabled custom skills.

#### Scenario: Builtin skills are injected for every user
- **WHEN** an authenticated user sends a message
- **THEN** the runtime loads the shared builtin skills and makes them available during that turn

#### Scenario: Enabled user skills are merged with builtin skills
- **WHEN** a user sends a message and has one or more enabled custom skills
- **THEN** the runtime includes those enabled custom skills alongside the shared builtin skills for that user's turn

#### Scenario: User without custom skills still receives an answer
- **WHEN** a user sends a message and has no enabled custom skills
- **THEN** the runtime answers using the shared builtin skills and base agent capabilities without failing

### Requirement: Runtime can invoke registered tools while serving web requests
The system SHALL allow the agent runtime to invoke only tools that are registered by the platform and explicitly enabled for the Web deployment, and SHALL return each tool execution result or policy rejection back into the runtime loop.

#### Scenario: Tool invocation completes within a chat turn
- **WHEN** the model requests a registered tool that is enabled for the current Web runtime
- **THEN** the runtime executes the tool, feeds the result back to the model, and continues the turn until a final assistant response is produced

#### Scenario: Tool rejection is returned into the runtime loop
- **WHEN** the model requests a tool that is registered but rejected by policy, scope validation, or deployment configuration
- **THEN** the runtime returns the rejection result into the tool-calling loop and continues safely without executing unauthorized operations

## ADDED Requirements

### Requirement: Runtime exposes deployment command resources through stable paths
The system SHALL provide the runtime with the fixed deployment application root and stable command resource paths so that builtin skills and authorized tools can invoke deployment-provided binaries or scripts without depending on the process current working directory.

#### Scenario: Skill invokes compiled command artifact
- **WHEN** a builtin or authorized custom skill needs to invoke a deployment-provided binary built from the `cmd` source tree
- **THEN** the runtime resolves that binary from the configured deployment command artifact path instead of inferring it from the process current directory

#### Scenario: Script-based helper uses fixed deployment path
- **WHEN** a skill or tool needs to invoke a `.py` helper shipped with the deployment
- **THEN** the runtime resolves the script path from the fixed deployment application root and keeps the execution working directory under the allowed user scope
