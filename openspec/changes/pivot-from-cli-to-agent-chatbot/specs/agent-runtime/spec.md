## ADDED Requirements

### Requirement: Runtime defaults to a general-purpose assistant role
The system SHALL instruct the runtime to behave as a general-purpose agent assistant that can support broad conversational tasks, rather than assuming the user is asking for software engineering help.

#### Scenario: User asks a non-coding question
- **WHEN** a user asks for general analysis, writing help, planning, or everyday assistance
- **THEN** the runtime answers within the same conversation product without redirecting the user to a coding-specific workflow

#### Scenario: User asks a coding question
- **WHEN** a user asks for software engineering help
- **THEN** the runtime can still assist, but does so as one supported domain within a broader assistant role

### Requirement: Runtime prefers direct answers before optional tool use
The system SHALL prefer producing a direct conversational answer when the request can be satisfied from model reasoning and available conversation context, and SHALL invoke tools only when they materially improve the result.

#### Scenario: No tool needed for simple request
- **WHEN** a user's request can be answered without external tool execution
- **THEN** the runtime returns an answer without invoking a tool

#### Scenario: Tool used only when beneficial
- **WHEN** a user's request needs external information, skill content, or controlled side effects
- **THEN** the runtime may invoke an authorized tool and continue the turn until a final assistant response is produced

### Requirement: Runtime does not rely on shell or local workspace operations in browser chat
The system SHALL serve the browser chat experience without assuming shell execution, local filesystem mutation, or user-directory access as part of normal runtime behavior.

#### Scenario: General browser turn proceeds without workspace tool assumptions
- **WHEN** the runtime handles a browser chat turn
- **THEN** it constructs and completes the response flow without depending on shell commands or operations against a user's local directory
