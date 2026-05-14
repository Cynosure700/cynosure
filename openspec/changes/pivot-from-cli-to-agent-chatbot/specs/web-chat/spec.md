## MODIFIED Requirements

### Requirement: Authenticated users can chat with the agent in a browser conversation
The system SHALL provide a browser-first chat interface and corresponding APIs that allow an authenticated user to create a conversation, send messages, and receive responses from a general-purpose assistant without depending on CLI workflows.

#### Scenario: User starts a new conversation
- **WHEN** an authenticated user creates a new conversation from the web interface
- **THEN** the system creates a conversation bound to that user and returns its identifier

#### Scenario: User sends a message
- **WHEN** an authenticated user sends a chat message in one of their conversations
- **THEN** the system forwards the message to the agent runtime and begins producing an assistant response for that conversation

### Requirement: Chat responses include assistant output and optional runtime details
The system SHALL return chat responses in a format that lets the web client prioritize the assistant's answer while optionally exposing runtime details such as tool execution or skill loading.

#### Scenario: Assistant answer is shown as the primary output
- **WHEN** an authenticated user receives a response in a conversation
- **THEN** the client can render the assistant answer directly without requiring the user to inspect developer-oriented runtime events

#### Scenario: Runtime detail is available when relevant
- **WHEN** the agent runtime invokes a tool or loads a skill while answering a user's message
- **THEN** the response stream includes runtime detail that the client may expose as secondary information

## ADDED Requirements

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
