# web-chat Specification

## Purpose
Define the browser chat experience and APIs delivered through the single web-service startup path, including assistant-first responses and secondary runtime events.

## Requirements
### Requirement: Authenticated users can chat with the agent in a browser conversation
The system SHALL provide a web chat interface and corresponding APIs through a single web-service startup path that allow an authenticated user to create a conversation, send messages, and receive assistant responses without depending on historical CLI bootstrap flows.

#### Scenario: User starts a new conversation after standard service startup
- **WHEN** an authenticated user accesses a normally started web-service instance and creates a new conversation
- **THEN** the system creates a conversation bound to that user and returns its identifier through the web API contract

#### Scenario: User sends a message after simplified deployment
- **WHEN** an authenticated user sends a chat message in one of their conversations to a locally started or cloud-deployed service instance
- **THEN** the system forwards the message to the agent runtime and returns an assistant response through the same browser API contract

### Requirement: Users can view their own conversation history
The system SHALL allow an authenticated user to list and inspect only their own conversations and messages.

#### Scenario: Conversation list shown
- **WHEN** an authenticated user requests their conversation list
- **THEN** the system returns only conversations owned by that user

#### Scenario: Cross-user conversation access blocked
- **WHEN** a user requests a conversation owned by another user
- **THEN** the system denies access to that conversation and its messages

### Requirement: Chat responses include assistant output and runtime events
The system SHALL prioritize assistant responses in a format that allows the web client to display final answers as the primary chat output while runtime events such as tool execution or skill loading remain secondary information.

#### Scenario: Primary chat flow works with optional runtime details
- **WHEN** a user receives a browser chat response from the refactored service
- **THEN** the assistant output remains sufficient to complete the main workflow even when runtime details are present only as secondary information

#### Scenario: Tool event displayed as secondary response detail
- **WHEN** the agent runtime invokes a tool while answering a user's message
- **THEN** the response stream includes a runtime event that the client can render alongside, but secondary to, the assistant response
