# agent-runtime Specification

## Purpose
Define how the web agent runtime starts as part of the web service, loads runtime capabilities for each user turn, preserves per-conversation context, and executes authorized tools within the deployment-aware workspace layout.

## Requirements
### Requirement: Agent runtime loads enabled skills from the database for the current user
The system SHALL allow the web runtime to start and serve browser chat requests without depending on CLI or REPL components, SHALL load the current user's enabled skills from persistent storage for each agent conversation turn, and SHALL keep any shared builtin runtime capabilities loaded from the resolved runtime workspace available for that turn.

#### Scenario: Service starts with deployment-aware runtime workspace
- **WHEN** the web service starts and resolves a runtime workspace root
- **THEN** the runtime initializes successfully using the configured web-service dependencies and the resolved runtime asset layout

#### Scenario: Enabled user skills are merged with builtin runtime skills
- **WHEN** a user sends a message and has one or more enabled skills
- **THEN** the runtime makes both the enabled user skills and the builtin runtime capabilities from the resolved runtime workspace available during that turn

#### Scenario: User without custom skills still receives an answer
- **WHEN** a user sends a message and has no enabled custom skills
- **THEN** the runtime answers using the shared builtin runtime capabilities and base agent behavior without failing

### Requirement: Runtime preserves conversation context per conversation
The system SHALL construct agent context from the targeted conversation's stored message history and SHALL keep different conversations isolated from one another.

#### Scenario: Existing conversation context reused
- **WHEN** a user sends a follow-up message in an existing conversation
- **THEN** the runtime builds the agent context using the prior messages from that same conversation

#### Scenario: Separate conversations isolated
- **WHEN** a user has multiple conversations
- **THEN** the runtime does not mix messages or runtime state across those conversations

### Requirement: Runtime can invoke registered tools while serving web requests
The system SHALL complete browser chat turns without requiring CLI or REPL interaction, SHALL allow the runtime to invoke only tools explicitly registered and authorized for the current service instance, SHALL return each tool execution result back into the runtime loop, and SHALL use the resolved runtime workspace assets, including `workspace/bin` and `workspace/cmd`, when command execution is needed.

#### Scenario: Standard browser turn completes without CLI interaction
- **WHEN** the runtime handles a normal browser chat turn that does not need external side effects
- **THEN** it produces the assistant response without depending on CLI agent code or terminal interaction

#### Scenario: Tool invocation completes within a chat turn
- **WHEN** the model requests a registered and authorized tool during a user's chat turn
- **THEN** the runtime executes the tool against the resolved runtime workspace, feeds the result back to the model, and continues the turn until a final assistant response is produced

#### Scenario: Authorized command uses runtime workspace assets
- **WHEN** the runtime invokes an authorized command-oriented tool during a browser chat turn
- **THEN** that tool may execute binaries or scripts from the resolved `workspace/bin` or `workspace/cmd` paths

### Requirement: Runtime exposes deployment command resources through stable paths
The system SHALL provide the runtime with the fixed deployment application root and stable command resource paths so that retained builtin capabilities and authorized tools can invoke deployment-provided binaries or scripts without depending on the process current working directory.

#### Scenario: Skill invokes compiled command artifact
- **WHEN** a retained builtin capability or authorized custom skill needs to invoke a deployment-provided binary built from the `cmd` source tree
- **THEN** the runtime resolves that binary from the configured deployment command artifact path instead of inferring it from the process current directory

#### Scenario: Script-based helper uses fixed deployment path
- **WHEN** a skill or tool needs to invoke a helper script shipped with the deployment
- **THEN** the runtime resolves the script path from the fixed deployment application root and keeps the execution working directory under the allowed user scope

### Requirement: Runtime validates web-service dependencies together with required runtime workspace assets
The system SHALL validate the dependencies required to start and serve the web application, including configuration, storage connectivity, and the runtime workspace assets required by enabled web-agent capabilities, and SHALL not treat CLI-only modules as startup requirements.

#### Scenario: Required runtime assets are validated for command-capable service
- **WHEN** the configured web-agent capabilities rely on binaries or scripts published in the runtime workspace
- **THEN** startup validates the resolved `workspace/bin` and `workspace/cmd` paths required by those capabilities before serving requests

#### Scenario: Missing CLI-only modules do not block startup
- **WHEN** historical CLI or REPL-only modules are absent from the service entry path
- **THEN** startup still proceeds as long as the required web-service configuration, storage dependencies, and runtime assets are valid

### Requirement: Runtime SHALL expose deployment-aware workspace paths to skills and tools
The system SHALL expose the resolved application home, runtime workspace root, command binary directory, and helper script directory to authorized skills and tools without requiring them to infer those locations from the process current working directory, and those exposed runtime asset paths SHALL stay consistent with the single resolved workspace root used by the active service instance.

#### Scenario: Runtime environment reflects packaged deployment paths
- **WHEN** the active runtime workspace root is resolved to `<app-home>/output/workspace`
- **THEN** the runtime-provided environment for skills and tools references the packaged deployment paths under that workspace

#### Scenario: Runtime environment reflects local debug paths
- **WHEN** the active runtime workspace root falls back to `<app-home>/workspace`
- **THEN** the runtime-provided environment for skills and tools references the local debug paths under that workspace

#### Scenario: Runtime environment exposes canonical derived asset paths
- **WHEN** the service finishes loading web configuration for a given runtime workspace root
- **THEN** the runtime environment exposed to skills and tools uses the canonical asset paths already derived from that workspace root instead of re-deriving alternate defaults later in the execution path

### Requirement: Runtime defaults to a general-purpose assistant role
The system SHALL instruct the runtime to behave as a general-purpose agent assistant that can support broad conversational tasks, rather than assuming the user is asking for software engineering help.

#### Scenario: User asks a non-coding question
- **WHEN** a user asks for general analysis, writing help, planning, or everyday assistance
- **THEN** the runtime answers within the same conversation product without redirecting the user to a coding-specific workflow

#### Scenario: User asks a coding question
- **WHEN** a user asks for software engineering help
- **THEN** the runtime can still assist, but does so as one supported domain within a broader assistant role

### Requirement: Runtime prefers direct answers before optional tool use
The system SHALL prefer producing a direct conversational answer when the request can be satisfied from model reasoning and available conversation context, and SHALL invoke tools only when they materially improve the result.

#### Scenario: No tool needed for simple request
- **WHEN** a user's request can be answered without external tool execution
- **THEN** the runtime returns an answer without invoking a tool

#### Scenario: Tool used only when beneficial
- **WHEN** a user's request needs external information, skill content, or controlled side effects
- **THEN** the runtime may invoke an authorized tool and continue the turn until a final assistant response is produced

### Requirement: Runtime does not rely on shell or local workspace operations in browser chat
The system SHALL serve the browser chat experience without assuming shell execution, local filesystem mutation, or user-directory access as part of normal runtime behavior.

#### Scenario: General browser turn proceeds without workspace tool assumptions
- **WHEN** the runtime handles a browser chat turn
- **THEN** it constructs and completes the response flow without depending on shell commands or operations against a user's local directory
