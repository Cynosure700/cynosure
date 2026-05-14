## ADDED Requirements

### Requirement: Todo manager maintains a list of tasks with status
The system SHALL maintain a list of todo items, each with an id, text, and status (one of `pending`, `in_progress`, `completed`), and render them in a human-readable format.

#### Scenario: Todo list is rendered
- **WHEN** the todo list contains items
- **THEN** the system SHALL render each item with a status marker (`[ ]` for pending, `[>]` for in_progress, `[x]` for completed), its id, and text, followed by a completion summary

#### Scenario: Empty todo list is rendered
- **WHEN** the todo list is empty
- **THEN** the system SHALL return a message indicating no tasks exist

### Requirement: Todo manager enforces maximum of 20 items
The system SHALL reject updates that contain more than 20 todo items.

#### Scenario: Update with more than 20 items
- **WHEN** the `todo` tool is called with more than 20 items
- **THEN** the system SHALL return an error message and SHALL NOT update the todo list

### Requirement: Todo manager enforces single in-progress task
The system SHALL allow at most one task to have `in_progress` status at any time.

#### Scenario: Multiple in-progress tasks are rejected
- **WHEN** the `todo` tool is called with more than one item having `in_progress` status
- **THEN** the system SHALL return an error message and SHALL NOT update the todo list

### Requirement: Todo items require non-empty text
The system SHALL reject todo items with empty or whitespace-only text.

#### Scenario: Item with empty text is rejected
- **WHEN** the `todo` tool is called with an item whose text is empty or only whitespace
- **THEN** the system SHALL return an error message and SHALL NOT update the todo list

### Requirement: Todo items auto-generate id if not provided
The system SHALL assign an auto-incrementing numeric id to todo items that do not have an explicit id.

#### Scenario: Item without id gets auto-generated id
- **WHEN** the `todo` tool is called with an item that has no `id` field
- **THEN** the system SHALL assign an id based on the item's position in the list (1-based)