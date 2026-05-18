## MODIFIED Requirements

### Requirement: Runtime can invoke registered tools while serving web requests
The system SHALL complete browser chat turns without requiring CLI or REPL interaction, SHALL allow the runtime to invoke only tools explicitly registered and authorized for the current service instance, SHALL return each tool execution result back into the runtime loop, SHALL use the resolved runtime workspace assets, including `workspace/bin` and `workspace/cmd`, when command execution is needed, and SHALL preserve the currently active filesystem-backed skill directory so that subsequent tool calls in the same runtime loop can execute relative to that skill when appropriate.

#### Scenario: Standard browser turn completes without CLI interaction
- **WHEN** the runtime handles a normal browser chat turn that does not need external side effects
- **THEN** it produces the assistant response without depending on CLI agent code or terminal interaction

#### Scenario: Tool invocation completes within a chat turn
- **WHEN** the model requests a registered and authorized tool during a user's chat turn
- **THEN** the runtime executes the tool against the resolved runtime workspace or active skill directory, feeds the result back to the model, and continues the turn until a final assistant response is produced

#### Scenario: Authorized command uses runtime workspace assets
- **WHEN** the runtime invokes an authorized command-oriented tool during a browser chat turn
- **THEN** that tool may execute binaries or scripts from the resolved `workspace/bin` or `workspace/cmd` paths

#### Scenario: Loaded filesystem skill becomes the active tool context
- **WHEN** the runtime successfully loads a skill from a `SKILL.md` file on disk and the model issues a subsequent tool call in the same runtime loop
- **THEN** the runtime makes that skill's parent directory available as the active default execution directory for that tool call
