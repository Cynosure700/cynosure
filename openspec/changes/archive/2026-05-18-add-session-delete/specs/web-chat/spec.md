## ADDED Requirements

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
