## ADDED Requirements

### Requirement: General chat mode exposes only a safe default tool set
The system SHALL expose only the subset of registered tools that are compatible with a browser-first chatbot product, and SHALL exclude shell execution, arbitrary file mutation, directory browsing, and other tools that assume access to a user's local workspace.

#### Scenario: Standard chat session omits shell and local workspace tools
- **WHEN** a user interacts with the default browser chat product
- **THEN** the runtime does not expose shell execution, arbitrary file mutation, directory browsing, or similar local-workspace tools

### Requirement: Restricted tool availability is communicated gracefully
The system SHALL handle requests for unavailable browser-incompatible tools by responding with a clear capability boundary while continuing to assist conversationally where possible.

#### Scenario: User asks for action requiring unsupported local operation
- **WHEN** a request would require shell execution, local file mutation, or access to a user's local directory
- **THEN** the assistant explains that the browser product does not support that operation and provides the best available non-local help instead of failing silently
