# Functional Prompts

This directory contains non-identity prompts used by Cynosure runtime features.
The base system prompt remains in `../system_prompt.md`.

- `context_summary.md`: full-history context compression summarizer.
- `memory_extraction.md`: project-scoped long-term memory extraction. It uses the `{{current_date}}` placeholder (replaced at runtime with today's date) so the model can convert relative dates to absolute dates. Extracts four kinds: preference, feedback, project, reference.
- `memory_selection.md`: project-scoped memory selection before prompt injection.
- `memory_consolidation.md`: project-scoped memory cleanup template. It uses `{{type_label}}` and `{{type_value}}` placeholders.
- `conversation_memory_update.md`: current-session memory maintenance.
- `general_subagent.md`: appended rules for the `general` child agent.
- `explore_subagent.md`: standalone system prompt for the read-only `explore` child agent. It uses the `{{workspace_root}}` placeholder.

Short structural messages that are tightly coupled to request shaping can remain hardcoded in Go.
