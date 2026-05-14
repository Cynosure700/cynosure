## ADDED Requirements

### Requirement: Tool registry provides parent and child tool sets
The system SHALL maintain two tool definition sets: a parent set containing all 8 tools (bash, read_file, write_file, edit_file, load_skill, todo, task, compact) and a child set containing 6 tools (bash, read_file, write_file, edit_file, load_skill, compact) excluding `task` and `todo` to prevent recursive subagent spawning.

#### Scenario: Parent agent has full tool set
- **WHEN** the parent agent loop starts
- **THEN** the system SHALL pass the parent tool definitions (8 tools) to the LLM API call

#### Scenario: Child agent has restricted tool set
- **WHEN** a subagent is created via the `task` tool
- **THEN** the system SHALL pass the child tool definitions (6 tools, excluding `task` and `todo`) to the LLM API call

### Requirement: Tool handlers dispatch by name
The system SHALL maintain a mapping from tool name to handler function and dispatch tool calls by looking up the handler by name.

#### Scenario: Known tool is dispatched
- **WHEN** the LLM returns a tool call with a name registered in the handler map
- **THEN** the system SHALL invoke the corresponding handler with the parsed arguments

#### Scenario: Unknown tool returns error
- **WHEN** the LLM returns a tool call with a name not registered in the handler map
- **THEN** the system SHALL return an error message indicating the tool is unknown

### Requirement: Tool definitions follow OpenAI function-calling schema
The system SHALL define each tool using the OpenAI function-calling format with `type: "function"`, a `function` object containing `name`, `description`, and `parameters` as a JSON Schema object.

#### Scenario: Tool definition is valid JSON Schema
- **WHEN** a tool definition is serialized for the LLM API call
- **THEN** the `parameters` field SHALL be a valid JSON Schema object describing the expected arguments