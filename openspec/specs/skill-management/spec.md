# skill-management Specification

## Purpose
Define per-user skill lifecycle management and runtime eligibility rules for the web agent platform.
## Requirements
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

### Requirement: Skill ownership is isolated per user
The system SHALL enforce skill ownership so that one user cannot read, modify, enable, disable, or delete another user's skills.

#### Scenario: Cross-user access blocked
- **WHEN** a user requests a skill resource owned by a different user
- **THEN** the system denies access to that skill resource

### Requirement: Only enabled skills are eligible for runtime loading
The system SHALL expose skill status for user-managed skills, SHALL load only user skills marked as enabled into the agent runtime, and SHALL keep the platform builtin skill catalog loaded from the resolved runtime workspace available independently of user-managed skill status.

#### Scenario: Enabled user skill is available together with builtin runtime skills
- **WHEN** a user's custom skill is marked enabled
- **THEN** the agent runtime includes that enabled user skill together with the builtin skills from the resolved runtime workspace

#### Scenario: Disabled user skill is excluded while builtin runtime skills remain available
- **WHEN** a user's custom skill is marked disabled
- **THEN** the agent runtime excludes that user skill but still keeps the builtin skills from the resolved runtime workspace available

### Requirement: Custom skill identifiers cannot shadow builtin skills
The system SHALL reject creation or update of a custom skill whose runtime identifier conflicts with a shared builtin skill identifier.

#### Scenario: Conflicting custom skill identifier is rejected
- **WHEN** a user creates or updates a custom skill with the same runtime identifier as a builtin skill
- **THEN** the system rejects the request and preserves the builtin skill identifier as reserved

### Requirement: Users cannot override deployed builtin skill identifiers
The system SHALL reject creation or update of a user-managed skill whose runtime identifier conflicts with a builtin skill loaded from the deployed `workspace/skills` catalog.

#### Scenario: Conflicting custom skill slug is rejected
- **WHEN** a user attempts to create or rename a custom skill to an identifier already present in `workspace/skills`
- **THEN** the system rejects the request and preserves the deployed builtin skill mapping

### Requirement: Builtin skill catalog SHALL prefer the resolved deployment workspace
The system SHALL load platform builtin skills from the builtin skill directory derived from the resolved runtime workspace root and SHALL not require the process current working directory to match the source repository layout.

#### Scenario: Builtin skills are loaded from packaged deployment workspace
- **WHEN** the resolved runtime workspace root is `<app-home>/output/workspace`
- **THEN** the platform builtin skill catalog loads from `<app-home>/output/workspace/skills`

#### Scenario: Builtin skills are loaded from source workspace during local debugging
- **WHEN** the resolved runtime workspace root falls back to `<app-home>/workspace`
- **THEN** the platform builtin skill catalog loads from `<app-home>/workspace/skills`

