## MODIFIED Requirements

### Requirement: Only enabled skills are eligible for runtime loading
The system SHALL expose skill status for user-managed skills, SHALL load only user skills marked as enabled into the agent runtime, and SHALL keep the platform builtin skill catalog loaded from the resolved runtime workspace available independently of user-managed skill status.

#### Scenario: Enabled user skill is available together with builtin runtime skills
- **WHEN** a user's custom skill is marked enabled
- **THEN** the agent runtime includes that enabled user skill together with the builtin skills from the resolved runtime workspace

#### Scenario: Disabled user skill is excluded while builtin runtime skills remain available
- **WHEN** a user's custom skill is marked disabled
- **THEN** the agent runtime excludes that user skill but still keeps the builtin skills from the resolved runtime workspace available

## ADDED Requirements

### Requirement: Builtin skill catalog SHALL prefer the resolved deployment workspace
The system SHALL load platform builtin skills from the builtin skill directory derived from the resolved runtime workspace root and SHALL not require the process current working directory to match the source repository layout.

#### Scenario: Builtin skills are loaded from packaged deployment workspace
- **WHEN** the resolved runtime workspace root is `<app-home>/output/workspace`
- **THEN** the platform builtin skill catalog loads from `<app-home>/output/workspace/skills`

#### Scenario: Builtin skills are loaded from source workspace during local debugging
- **WHEN** the resolved runtime workspace root falls back to `<app-home>/workspace`
- **THEN** the platform builtin skill catalog loads from `<app-home>/workspace/skills`
