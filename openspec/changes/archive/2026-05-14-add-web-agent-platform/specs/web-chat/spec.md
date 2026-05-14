## ADDED Requirements

### Requirement: Authenticated users can chat with the agent in a browser conversation
The system SHALL provide a web chat interface and corresponding APIs that allow an authenticated user to create a conversation, send messages, and receive assistant responses.

#### Scenario: User starts a new conversation
- **WHEN** an authenticated user creates a new conversation from the web interface
- **THEN** the system creates a conversation bound to that user and returns its identifier

#### Scenario: User sends a message
- **WHEN** an authenticated user sends a chat message in one of their conversations
- **THEN** the system forwards the message to the agent runtime and begins producing an assistant response for that conversation

### Requirement: Users can view their own conversation history
The system SHALL allow an authenticated user to list and inspect only their own conversations and messages.

#### Scenario: Conversation list shown
- **WHEN** an authenticated user requests their conversation list
- **THEN** the system returns only conversations owned by that user

#### Scenario: Cross-user conversation access blocked
- **WHEN** a user requests a conversation owned by another user
- **THEN** the system denies access to that conversation and its messages

### Requirement: Chat responses include assistant output and runtime events
The system SHALL return assistant responses in a format that allows the web client to display final answers and runtime events such as tool execution or skill loading.

#### Scenario: Tool event displayed during response
- **WHEN** the agent runtime invokes a tool while answering a user's message
- **THEN** the response stream includes a runtime event that the client can render alongside the assistant response

