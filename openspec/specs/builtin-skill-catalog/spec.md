# builtin-skill-catalog Specification

## Purpose
TBD - created by archiving change enable-go-agent-skill-and-tool-runtime. Update Purpose after archive.
## Requirements
### Requirement: Builtin skills can reference deployment command artifacts through stable conventions
The system SHALL allow builtin skills to reference deployment-provided command artifacts through stable deployment path conventions so that those skills remain portable across environments.

#### Scenario: Builtin skill references compiled binary by deployment convention
- **WHEN** a builtin skill needs a helper command that is compiled from the `cmd` source tree
- **THEN** the skill can rely on the documented deployment command artifact convention rather than assuming a repository-relative current working directory

