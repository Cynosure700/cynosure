# web-chat Specification

## Purpose
Define the browser chat experience and APIs delivered through the single web-service startup path, including assistant-first responses and secondary runtime events.

## Requirements
### Requirement: Authenticated users can chat with the agent in a browser conversation
The system SHALL provide a browser-first chat interface and corresponding APIs that allow an authenticated user to create a conversation, send messages, and receive responses from a general-purpose assistant without depending on CLI workflows.

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

### Requirement: Authenticated users can delete their own conversations from the web chat UI
The system SHALL allow an authenticated user to delete a conversation they own from the web chat experience and SHALL update the visible conversation state so the deleted conversation is no longer selectable.

#### Scenario: User deletes a non-active conversation
- **WHEN** an authenticated user deletes one of their conversations that is not currently active
- **THEN** the system removes that conversation from the user's conversation list
- **AND** the currently active conversation remains unchanged

#### Scenario: User deletes the active conversation
- **WHEN** an authenticated user deletes the conversation currently open in the chat panel
- **THEN** the system removes that conversation from the user's conversation list
- **AND** the system switches the chat panel to another remaining conversation or an empty-state view when no conversations remain

#### Scenario: Cross-user conversation deletion blocked
- **WHEN** an authenticated user requests deletion of a conversation owned by another user
- **THEN** the system denies the deletion request
- **AND** the other user's conversation remains available to its owner

### Requirement: Chat responses include assistant output and optional runtime details
The system SHALL return chat responses in a format that lets the web client prioritize the assistant's answer while optionally exposing runtime details such as tool execution or skill loading.

#### Scenario: Assistant answer is shown as the primary output
- **WHEN** an authenticated user receives a response in a conversation
- **THEN** the client can render the assistant answer directly without requiring the user to inspect developer-oriented runtime events

#### Scenario: Runtime detail is available when relevant
- **WHEN** the agent runtime invokes a tool or loads a skill while answering a user's message
- **THEN** the response stream includes runtime detail that the client may expose as secondary information

### Requirement: Chat interface is conversation-first
The system SHALL organize the main web chat experience around conversation history, message composition, and assistant replies, with advanced controls presented outside the primary interaction flow.

#### Scenario: User focuses on conversation flow
- **WHEN** a user opens the main chat screen
- **THEN** the primary visible actions are conversation selection, message reading, and message sending rather than developer-centric controls

### Requirement: Chat interface follows a low-friction chatbot interaction model
The system SHALL present the browser experience in a low-friction chatbot style similar to ChatGPT, where the assistant reply is the focal point and implementation-oriented controls do not dominate the primary screen.

#### Scenario: Primary screen emphasizes conversation over tooling
- **WHEN** a user views the main chat product
- **THEN** the interface emphasizes the assistant conversation and composer area rather than tool panels, workspace actions, or skill editing controls
