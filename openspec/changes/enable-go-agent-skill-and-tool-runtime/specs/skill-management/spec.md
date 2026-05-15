## MODIFIED Requirements

### Requirement: Users can manage their own skills
The system SHALL allow each authenticated user to create, view, update, enable, disable, and delete only their own database-backed custom skills through the web platform. Platform builtin skills loaded from `go-agent/skills` SHALL be shared runtime capabilities and SHALL NOT be mutable through user skill management APIs.

#### Scenario: User creates a custom skill
- **WHEN** an authenticated user submits a valid custom skill name, description, and content
- **THEN** the system stores that custom skill under that user's ownership

#### Scenario: User updates a custom skill
- **WHEN** an authenticated user edits one of their existing custom skills
- **THEN** the system saves the updated custom skill content and metadata for that same user

#### Scenario: User cannot mutate a builtin skill through custom skill APIs
- **WHEN** an authenticated user attempts to update or delete a builtin skill through the user skill management APIs
- **THEN** the system denies the mutation request because builtin skills are not user-owned resources

### Requirement: Only enabled skills are eligible for runtime loading
The system SHALL expose custom skill status and SHALL load only custom skills marked as enabled into the user-specific portion of the agent runtime. Shared builtin skills SHALL remain available independently of any per-user enablement state.

#### Scenario: Enabled custom skill available to runtime
- **WHEN** a user has a custom skill with enabled status
- **THEN** the agent runtime includes that custom skill alongside the shared builtin skills in the user's active skill set

#### Scenario: Disabled custom skill excluded from runtime
- **WHEN** a user's custom skill is marked disabled
- **THEN** the agent runtime excludes that custom skill from the user's active skill set while keeping shared builtin skills available

## ADDED Requirements

### Requirement: Custom skill identifiers cannot shadow builtin skills
The system SHALL reject creation or update of a custom skill whose runtime identifier conflicts with a shared builtin skill identifier.

#### Scenario: Conflicting custom skill identifier is rejected
- **WHEN** a user creates or updates a custom skill with the same runtime identifier as a builtin skill
- **THEN** the system rejects the request and preserves the builtin skill identifier as reserved
