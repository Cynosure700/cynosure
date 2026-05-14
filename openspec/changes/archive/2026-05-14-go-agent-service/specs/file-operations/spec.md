## ADDED Requirements

### Requirement: Read file returns file contents with optional line limit
The system SHALL read a file at the given path (relative to the working directory), apply path safety checks, and return its contents. If a line limit is specified, only the first N lines SHALL be returned with a truncation notice.

#### Scenario: Read entire file
- **WHEN** the `read_file` tool is called with a valid path and no limit
- **THEN** the system SHALL return the complete file contents

#### Scenario: Read file with line limit
- **WHEN** the `read_file` tool is called with a valid path and a limit of N
- **THEN** the system SHALL return the first N lines followed by a notice indicating how many lines were omitted

#### Scenario: Read empty file
- **WHEN** the `read_file` tool is called on an empty file
- **THEN** the system SHALL return a message indicating the file is empty

#### Scenario: Read non-existent file
- **WHEN** the `read_file` tool is called on a path that does not exist
- **THEN** the system SHALL return an error message

### Requirement: Write file creates or overwrites a file
The system SHALL write content to a file at the given path, creating parent directories as needed, and return a confirmation with the byte count.

#### Scenario: Write to new file
- **WHEN** the `write_file` tool is called with a path that does not exist
- **THEN** the system SHALL create the file and all necessary parent directories, write the content, and return the number of bytes written

#### Scenario: Overwrite existing file
- **WHEN** the `write_file` tool is called with a path that already exists
- **THEN** the system SHALL overwrite the file with the new content and return the number of bytes written

### Requirement: Edit file replaces exact text match
The system SHALL find the first occurrence of the given old text in the file and replace it with the new text. If the old text is not found, an error SHALL be returned.

#### Scenario: Successful text replacement
- **WHEN** the `edit_file` tool is called with old text that exists exactly once in the file
- **THEN** the system SHALL replace that occurrence with the new text and return a success message

#### Scenario: Old text not found
- **WHEN** the `edit_file` tool is called with old text that does not exist in the file
- **THEN** the system SHALL return an error message indicating the text was not found

### Requirement: All file operations enforce path sandbox
The system SHALL resolve all file paths relative to the working directory and reject any path that resolves outside the working directory.

#### Scenario: Path within workspace is allowed
- **WHEN** a file operation is called with a path that resolves within the working directory
- **THEN** the system SHALL proceed with the operation

#### Scenario: Path escaping workspace is rejected
- **WHEN** a file operation is called with a path containing `../` that resolves outside the working directory
- **THEN** the system SHALL return an error message and SHALL NOT perform the operation

### Requirement: File output is truncated at 50000 characters
The system SHALL truncate file read output at 50,000 characters to prevent excessive token consumption.

#### Scenario: Large file is read
- **WHEN** a file read returns more than 50,000 characters
- **THEN** the system SHALL return only the first 50,000 characters