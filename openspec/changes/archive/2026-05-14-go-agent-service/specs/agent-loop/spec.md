## ADDED Requirements

### Requirement: Agent loop processes user input with tool-calling cycle
The system SHALL accept a system prompt and user message, then enter a tool-calling loop that repeatedly invokes the LLM until it returns a final text response without tool calls.

#### Scenario: Single-turn response without tools
- **WHEN** the LLM returns a message with `finish_reason` other than `tool_calls` and no `tool_calls` array
- **THEN** the system SHALL return the message content as the final response and exit the loop

#### Scenario: Multi-turn tool-calling cycle
- **WHEN** the LLM returns a message with `finish_reason` equal to `tool_calls` containing one or more tool calls
- **THEN** the system SHALL execute each tool handler in order, append the results as `role: "tool"` messages, and continue the loop

### Requirement: Agent loop injects todo reminder after inactivity
The system SHALL track the number of consecutive rounds without a `todo` tool call and inject a reminder message when the count reaches 3.

#### Scenario: Reminder after 3 rounds without todo update
- **WHEN** the agent has completed 3 consecutive rounds without calling the `todo` tool
- **THEN** the system SHALL append a user message containing a reminder to update the todo list before the next LLM call

#### Scenario: Reminder counter resets on todo call
- **WHEN** the agent calls the `todo` tool
- **THEN** the system SHALL reset the rounds-since-todo counter to 0

### Requirement: Agent loop integrates context compaction pipeline
The system SHALL execute the three-layer context compaction pipeline within the agent loop: micro_compact every round, auto_compact when token threshold is exceeded, and manual compact when the model invokes the `compact` tool.

#### Scenario: Micro compact runs every round
- **WHEN** the agent loop begins a new iteration
- **THEN** the system SHALL call micro_compact on the messages array before making the LLM API call

#### Scenario: Auto compact triggers on token threshold
- **WHEN** the estimated token count of messages exceeds 50,000
- **THEN** the system SHALL call auto_compact to summarize and replace the message history before making the LLM API call

#### Scenario: Manual compact after compact tool call
- **WHEN** the agent invokes the `compact` tool
- **THEN** the system SHALL execute auto_compact immediately after processing all tool results in that round