## ADDED Requirements

### Requirement: Subagent executes tasks with isolated context
The system SHALL create subagents with an independent message history, a restricted tool set (child tools only), and a maximum of 30 rounds, returning only the final text summary to the parent agent.

#### Scenario: Subagent completes task within round limit
- **WHEN** a subagent is created with a task prompt and completes within 30 rounds
- **THEN** the system SHALL return the subagent's final text content as the result

#### Scenario: Subagent reaches round limit
- **WHEN** a subagent reaches 30 rounds without the LLM returning a final response
- **THEN** the system SHALL stop the loop and return the last text content produced by the subagent

#### Scenario: Subagent context is isolated from parent
- **WHEN** a subagent is created
- **THEN** the subagent's message history SHALL be independent and SHALL NOT include the parent agent's conversation history

### Requirement: Subagent uses fixed system prompt
The system SHALL use a fixed system prompt for subagents that describes them as coding subagents and instructs them to complete the given task and summarize findings.

#### Scenario: Subagent system prompt is applied
- **WHEN** a subagent is initialized
- **THEN** the first message in its history SHALL be a system message with the fixed subagent prompt

### Requirement: Subagent cannot spawn further subagents
The system SHALL prevent subagents from calling the `task` tool by using the child tool definition set which excludes `task`.

#### Scenario: Subagent tool set excludes task
- **WHEN** a subagent's tool definitions are passed to the LLM
- **THEN** the `task` tool SHALL NOT be present in the definitions