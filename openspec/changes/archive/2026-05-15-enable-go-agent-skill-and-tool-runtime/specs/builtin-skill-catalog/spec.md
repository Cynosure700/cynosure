## ADDED Requirements

### Requirement: Platform loads builtin skills from the repository catalog
The system SHALL load builtin skills from the configured `go-agent/skills` directory when the Web service starts and SHALL keep them in a shared runtime catalog that is available to every chat turn.

#### Scenario: Service startup loads builtin skill catalog
- **WHEN** the service starts with one or more valid builtin skill definitions in the configured catalog directory
- **THEN** the system loads those skills into a shared in-memory catalog before serving chat requests

### Requirement: Builtin skills are available to every authenticated user
The system SHALL include builtin skills in each authenticated user's runtime skill set regardless of whether that user has any custom skills stored in the database.

#### Scenario: User without custom skills can still call builtin skills
- **WHEN** an authenticated user has no enabled custom skills in the database
- **THEN** the runtime still exposes the shared builtin skill catalog for that user's chat turn

#### Scenario: User with custom skills receives merged catalog
- **WHEN** an authenticated user has one or more enabled custom skills
- **THEN** the runtime exposes both the shared builtin skills and that user's enabled custom skills in the same skill set

## ADDED Requirements

### Requirement: Builtin skills can reference deployment command artifacts through stable conventions
The system SHALL allow builtin skills to reference deployment-provided command artifacts through stable deployment path conventions so that those skills remain portable across environments.

#### Scenario: Builtin skill references compiled binary by deployment convention
- **WHEN** a builtin skill needs a helper command that is compiled from the `cmd` source tree
- **THEN** the skill can rely on the documented deployment command artifact convention rather than assuming a repository-relative current working directory
