## 1. Runtime skill context

- [x] 1.1 Extend skill loading/output handling so a successfully loaded filesystem-backed skill exposes its source directory as active runtime context.
- [x] 1.2 Update the web runtime tool-call loop to carry the active skill directory across subsequent tool invocations within the same response turn.

## 2. Tool working directory resolution

- [x] 2.1 Update runtime environment injection and default cwd resolution so skill-driven tool calls prefer the active skill directory and otherwise fall back to the workspace root.
- [x] 2.2 Preserve workspace-boundary validation for skill-driven cwd switching, including fallback or rejection behavior for non-filesystem or out-of-scope skills.
- [x] 2.3 Ensure tool audit metadata records the resolved cwd produced by skill-directory execution.

## 3. Regression coverage

- [ ] 3.1 Add runtime tests covering filesystem-backed skill activation, multi-tool turn behavior, and DB-skill fallback to workspace root.
- [ ] 3.2 Add tool execution tests covering skill-directory cwd selection, workspace escape rejection, and audit visibility for skill-driven executions.
