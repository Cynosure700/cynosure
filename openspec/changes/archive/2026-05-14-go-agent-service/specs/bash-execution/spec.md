## ADDED Requirements

### Requirement: Bash tool executes shell commands with timeout
The system SHALL execute shell commands via `bash -c` with a default timeout of 120 seconds, capturing both stdout and stderr.

#### Scenario: Command completes within timeout
- **WHEN** a shell command is executed and completes within 120 seconds
- **THEN** the system SHALL return the combined stdout and stderr output

#### Scenario: Command exceeds timeout
- **WHEN** a shell command runs longer than 120 seconds
- **THEN** the system SHALL terminate the process and return a timeout error message

#### Scenario: Command produces no output
- **WHEN** a shell command completes successfully but produces no stdout or stderr
- **THEN** the system SHALL return a message indicating no output was produced

### Requirement: Bash tool blocks dangerous commands
The system SHALL reject commands that contain dangerous patterns: `rm -rf /`, `sudo`, `shutdown`, `reboot`, `> /dev/`.

#### Scenario: Dangerous command is blocked
- **WHEN** a command contains any of the dangerous patterns
- **THEN** the system SHALL return an error message and SHALL NOT execute the command

#### Scenario: Safe command is allowed
- **WHEN** a command does not contain any dangerous patterns
- **THEN** the system SHALL proceed with execution

### Requirement: Bash output is truncated at 50000 characters
The system SHALL truncate command output at 50,000 characters to prevent excessive token consumption.

#### Scenario: Command produces large output
- **WHEN** a command produces more than 50,000 characters of output
- **THEN** the system SHALL return only the first 50,000 characters