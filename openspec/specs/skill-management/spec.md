# skill-management Specification

## Purpose
Define per-user skill lifecycle management and runtime eligibility rules for the web agent platform.

## Requirements
### Requirement: Users can manage their own skills
The system SHALL allow each authenticated user to create, view, update, enable, disable, and delete only their own skills through the web platform.

#### Scenario: User creates a skill
- **WHEN** an authenticated user submits a valid skill name, description, and content
- **THEN** the system stores the skill under that user's ownership

#### Scenario: User updates a skill
- **WHEN** an authenticated user edits one of their existing skills
- **THEN** the system saves the updated skill content and metadata for that same user

#### Scenario: User deletes a skill
- **WHEN** an authenticated user deletes one of their own skills
- **THEN** the system removes the skill from that user's active skill list

### Requirement: Skill ownership is isolated per user
The system SHALL enforce skill ownership so that one user cannot read, modify, enable, disable, or delete another user's skills.

#### Scenario: Cross-user access blocked
- **WHEN** a user requests a skill resource owned by a different user
- **THEN** the system denies access to that skill resource

### Requirement: Only enabled skills are eligible for runtime loading
The system SHALL expose skill status and SHALL load only skills marked as enabled into the agent runtime.

#### Scenario: Enabled skill available to runtime
- **WHEN** a user has a skill with enabled status
- **THEN** the agent runtime includes that skill in the user's active skill set

#### Scenario: Disabled skill excluded from runtime
- **WHEN** a user's skill is marked disabled
- **THEN** the agent runtime excludes that skill from the user's active skill set
