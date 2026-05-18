## ADDED Requirements

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
