## ADDED Requirements

### Requirement: Core go-agent modules SHALL be split by responsibility without changing supported behavior
The system SHALL organize core `go-agent` implementation hotspots into responsibility-focused modules so that HTTP entry handling, runtime orchestration, tool infrastructure, configuration parsing, storage access, and skill loading are no longer concentrated in a few oversized files, while continuing to preserve the currently supported APIs, runtime boundaries, deployment layout, and configuration behavior.

#### Scenario: Refactor preserves supported runtime behavior
- **WHEN** maintainers modularize the existing `go-agent` codebase
- **THEN** the supported Web API behavior, conversation runtime flow, tool allowlist semantics, workspace isolation rules, and deployment artifact layout remain unchanged

### Requirement: Web and runtime layers SHALL expose clear ownership boundaries
The system SHALL separate module ownership so that the web application layer focuses on server assembly, routing, and request handling, while the runtime layer focuses on conversation orchestration and web-runtime-specific tool flow, and the shared tools layer continues to own generic tool definitions, runtime environment handling, and path-safe tool execution.

#### Scenario: Tooling concerns do not remain duplicated across layers
- **WHEN** maintainers inspect the modularized implementation
- **THEN** web runtime specific responsibilities such as `load_skill` orchestration and SSE flow live under the runtime layer, and generic file or command tool definitions and path-guarded handlers live under the shared tools layer

### Requirement: Project documentation SHALL describe the modularized structure accurately
The system SHALL keep `go-agent/README.md` aligned with the implemented module layout so that developers can understand the main directory structure, core module responsibilities, and recommended verification commands after the refactor.

#### Scenario: README matches modularized project structure
- **WHEN** a developer reads `go-agent/README.md` after the refactor
- **THEN** the documented directory structure and module descriptions reflect the reorganized code layout and still describe how to validate the service without implying behavior changes
