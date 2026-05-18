## MODIFIED Requirements

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
