## ADDED Requirements

### Requirement: Micro compact replaces old tool results with placeholders
The system SHALL replace tool result messages (role: "tool") older than the 3 most recent with a placeholder string `[Previous: used <tool_name>]`, preserving the most recent 3 tool results intact.

#### Scenario: More than 3 tool results exist
- **WHEN** the message history contains more than 3 tool result messages
- **THEN** the system SHALL keep the 3 most recent tool results unchanged and replace all older tool results with placeholder strings

#### Scenario: Three or fewer tool results exist
- **WHEN** the message history contains 3 or fewer tool result messages
- **THEN** the system SHALL leave all tool results unchanged

### Requirement: Auto compact summarizes conversation when token threshold exceeded
The system SHALL save the full conversation to a transcript file, call the LLM to generate a summary, and replace the message history with the summary when the estimated token count exceeds 50,000.

#### Scenario: Token threshold is exceeded
- **WHEN** the estimated token count of messages exceeds 50,000
- **THEN** the system SHALL save the full conversation to `.transcripts/transcript_<timestamp>.jsonl`, call the LLM for a summary, and replace messages with the summary result

#### Scenario: Token threshold is not exceeded
- **WHEN** the estimated token count of messages is 50,000 or below
- **THEN** the system SHALL NOT trigger auto compaction

### Requirement: Token estimation uses character count divided by 4
The system SHALL estimate token count as `len(JSON.stringify(messages)) / 4`.

#### Scenario: Token estimation is calculated
- **WHEN** the system needs to check the token threshold
- **THEN** it SHALL serialize messages to JSON, divide the byte length by 4, and use the result as the estimated token count

### Requirement: Compact tool triggers manual compaction
The system SHALL provide a `compact` tool that, when invoked by the model, triggers the same auto_compact behavior after the current round completes.

#### Scenario: Model invokes compact tool
- **WHEN** the model calls the `compact` tool
- **THEN** the system SHALL execute auto_compact after processing all tool results in the current round

### Requirement: Transcript files are saved in JSONL format
The system SHALL save full conversation transcripts as JSONL files (one JSON object per line) in the `.transcripts/` directory with timestamped filenames.

#### Scenario: Auto compact saves transcript
- **WHEN** auto_compact is triggered
- **THEN** the system SHALL create the `.transcripts/` directory if it does not exist and write the full message history as a JSONL file