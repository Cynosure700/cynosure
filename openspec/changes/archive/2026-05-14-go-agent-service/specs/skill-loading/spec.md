## ADDED Requirements

### Requirement: Skill loader scans skills directory at startup
The system SHALL scan the `skills/` directory for `.md` files at startup, parse their YAML frontmatter, and build an in-memory index of available skills.

#### Scenario: Skills directory contains markdown files
- **WHEN** the system starts and the `skills/` directory contains `.md` files with valid YAML frontmatter
- **THEN** the system SHALL load all skills into memory with their metadata and body content

#### Scenario: Skills directory is empty or missing
- **WHEN** the system starts and the `skills/` directory does not exist or contains no `.md` files
- **THEN** the system SHALL initialize with an empty skill index and continue without error

### Requirement: Skill descriptions are injected into system prompt
The system SHALL generate a compact summary of available skills (name and description) and include it in the system prompt for the agent.

#### Scenario: System prompt includes skill list
- **WHEN** the system prompt is constructed
- **THEN** it SHALL include a list of available skill names and their descriptions from the YAML frontmatter

### Requirement: Skills are loaded on demand via load_skill tool
The system SHALL return the full content of a skill only when the `load_skill` tool is called with the skill name, wrapping it in XML-style tags.

#### Scenario: Valid skill name is requested
- **WHEN** the `load_skill` tool is called with a name matching a loaded skill
- **THEN** the system SHALL return the skill body wrapped in `<skill name="...">...</skill>` tags

#### Scenario: Unknown skill name is requested
- **WHEN** the `load_skill` tool is called with a name not matching any loaded skill
- **THEN** the system SHALL return an error message listing the available skill names

### Requirement: Skill files use YAML frontmatter format
The system SHALL parse skill files with the format: `---` delimited YAML frontmatter followed by Markdown body content.

#### Scenario: Valid frontmatter is parsed
- **WHEN** a skill file contains `---\nkey: value\n---\nbody content`
- **THEN** the system SHALL extract the frontmatter as metadata and the remaining content as the skill body

#### Scenario: File without frontmatter
- **WHEN** a skill file does not start with `---`
- **THEN** the system SHALL treat the entire file content as the skill body with empty metadata