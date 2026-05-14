# agent-runtime Specification

## Purpose
Define how the web agent runtime loads enabled user skills, preserves per-conversation context, and executes registered tools during web chat turns.

## Requirements
### Requirement: Agent runtime loads enabled skills from the database for the current user
The system SHALL load the current user's enabled skills from persistent storage for each agent conversation turn and make them available to the model through the runtime skill-loading mechanism.

#### Scenario: Enabled user skills injected into runtime
- **WHEN** a user sends a message and has one or more enabled skills
- **THEN** the runtime loads those skills for that user and makes them available to the model during that turn

#### Scenario: User without custom skills still receives an answer
- **WHEN** a user sends a message and has no enabled custom skills
- **THEN** the runtime answers using the base agent capabilities without failing

### Requirement: Runtime preserves conversation context per conversation
The system SHALL construct agent context from the targeted conversation's stored message history and SHALL keep different conversations isolated from one another.

#### Scenario: Existing conversation context reused
- **WHEN** a user sends a follow-up message in an existing conversation
- **THEN** the runtime builds the agent context using the prior messages from that same conversation

#### Scenario: Separate conversations isolated
- **WHEN** a user has multiple conversations
- **THEN** the runtime does not mix messages or runtime state across those conversations

### Requirement: Runtime can invoke registered tools while serving web requests
The system SHALL allow the agent runtime to invoke registered tools during a web chat turn and SHALL return the tool execution results back into the runtime loop.

#### Scenario: Tool invocation completes within a chat turn
- **WHEN** the model requests a registered tool during a user's chat turn
- **THEN** the runtime executes the tool, feeds the result back to the model, and continues the turn until a final assistant response is produced
