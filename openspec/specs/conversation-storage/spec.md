# conversation-storage Specification

## Purpose
Define durable conversation, message, and tool execution storage for the web agent platform, with Redis used only as an active-state cache.

## Requirements
### Requirement: Conversations and messages are persisted
The system SHALL persist conversations and chat messages so that users can resume prior chats after the current request or process ends.

#### Scenario: Conversation persisted after creation
- **WHEN** a new conversation is created
- **THEN** the system stores the conversation metadata in persistent storage

#### Scenario: Message persisted after send
- **WHEN** a user message or assistant reply is accepted by the system
- **THEN** the system stores that message under the corresponding conversation

### Requirement: Tool execution records are auditable
The system SHALL persist tool execution records linked to the originating user and conversation.

#### Scenario: Tool call record saved
- **WHEN** the runtime executes a tool during a conversation turn
- **THEN** the system stores a tool call record containing the user, conversation, tool name, and execution result summary

### Requirement: Active conversation state can be cached in Redis
The system SHALL support caching active conversation state in Redis without making Redis the sole source of truth.

#### Scenario: Active conversation cached
- **WHEN** a conversation becomes active during runtime processing
- **THEN** the system may cache recent context in Redis while keeping the durable copy in persistent storage

#### Scenario: Redis miss falls back to database
- **WHEN** active conversation state is not present in Redis
- **THEN** the runtime rebuilds the required context from persistent storage and continues processing

### Requirement: Deleted conversations are removed with related persisted state
The system SHALL remove a deleted conversation from durable storage and SHALL clear any related persisted or cached runtime state tied to that conversation.

#### Scenario: Conversation delete removes related records
- **WHEN** a user-owned conversation is deleted
- **THEN** the system removes the conversation record
- **AND** the system removes persisted messages and tool execution records linked to that conversation

#### Scenario: Conversation delete clears active cache
- **WHEN** the system deletes a conversation that has cached active context in Redis
- **THEN** the system clears the cache entry associated with that conversation before reporting success

#### Scenario: Deleted conversation cannot be loaded again
- **WHEN** a client requests a conversation after it has been successfully deleted
- **THEN** the system reports that the conversation no longer exists for that user
