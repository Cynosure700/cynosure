## MODIFIED Requirements

### Requirement: Only enabled skills are eligible for runtime loading
The system SHALL load only user-managed skills marked as enabled into the current user's runtime skill set and SHALL merge them with the builtin skill catalog loaded from the deployed `workspace/skills` directory. Builtin skills SHALL be shared across all users and SHALL not depend on per-user ownership records or user-managed enablement state.

#### Scenario: Enabled skill available to runtime
- **WHEN** a user has a custom skill with enabled status
- **THEN** the agent runtime includes that skill in the user's active skill set together with the shared builtin workspace skills

#### Scenario: Disabled skill excluded from runtime
- **WHEN** a user's custom skill is marked disabled
- **THEN** the agent runtime excludes that custom skill from the user's active skill set while continuing to expose the shared builtin workspace skills

#### Scenario: Builtin workspace skill is available without user ownership
- **WHEN** an authenticated user starts a conversation and the deployed `workspace/skills` catalog is populated
- **THEN** the runtime exposes those builtin skills even if the user has no skill records in persistent storage

## ADDED Requirements

### Requirement: Users cannot override deployed builtin skill identifiers
The system SHALL reject creation or update of a user-managed skill whose runtime identifier conflicts with a builtin skill loaded from the deployed `workspace/skills` catalog.

#### Scenario: Conflicting custom skill slug is rejected
- **WHEN** a user attempts to create or rename a custom skill to an identifier already present in `workspace/skills`
- **THEN** the system rejects the request and preserves the deployed builtin skill mapping
