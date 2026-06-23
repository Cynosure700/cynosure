# Sub-Agent Tool Status Display Design

## Goal

Optimize the TUI display for child tool calls emitted by `spawn_subagent`: show only the latest child tool call under the running status area, align it with the existing status column, and display every tool argument without relying on shortened previews.

## Current Behavior

Sub-agent tool events arrive in the parent TUI as normal tool messages with `scope: "subagent"`, `ephemeral_group_id`, `suppress_result`, `raw_args`, and `args_preview`. The renderer already gives these messages muted gray styling, no leading blue bullet, and status-column alignment. Each child tool start currently appends another temporary message, so multiple child tool calls can be visible until the group is cleared.

## Design

Keep the runtime event contract unchanged. In the TUI layer, when a new `tool_call_start` event has `scope: "subagent"` and a non-empty `ephemeral_group_id`, remove existing child tool messages from that same group before appending the new one. This keeps the message model and rendered output to one latest child tool line for each active sub-agent run.

Render sub-agent child tool rows as compact tool-call information only: `tool(args...)`. Do not render the tool icon, checkmark, running line, success line, or result preview. Attach the child row directly below the parent running block with a single newline and align it to the same status column.

If a `tool_call_done` event later arrives for an older child tool that was already replaced, ignore it instead of appending a completed status row. This prevents finished child calls from reappearing below the latest running child tool.

For argument display, render tool calls from `raw_args` when it contains a JSON object. Format all top-level parameters in a stable key order as `key: value`, joining parameters with `, `. Preserve the existing `args_preview` fallback for malformed or non-object payloads so unusual tool events still render.

For status metadata, sub-agent event forwarding must not propagate child `tool_call_count` into the parent UI. Context token metadata may still be forwarded, but the parent tool count should reflect parent-context tool calls only, such as the `spawn_subagent` call itself.

For concurrent sub-agents, child tool events must carry the parent `spawn_subagent` tool-call id as `parent_tool_call_id`. The TUI inserts each child tool row immediately after the matching parent row by that id. Render adjacency alone is not reliable because multiple sub-agents can run and emit child tool events out of order.

## Constraints

- Do not change runtime event emission, result suppression, or group clearing.
- Preserve sub-agent gray styling, no leading blue bullet, and whole-block status-column alignment.
- Do not show sub-agent child tool status icons or status/result lines.
- Preserve existing display-name mappings such as `read_file` -> `read`, `grep` -> `grep`, and `spawn_subagent` subtype names.
- Keep changes scoped to the TUI renderer and focused tests.

## Testing

Add TUI tests that first fail under current behavior:

- Starting a second child tool in the same `ephemeral_group_id` replaces the first one.
- A late `tool_call_done` for the replaced child tool does not re-add it.
- A sub-agent child tool row attaches directly under the parent running block and shows only tool-call information.
- Concurrent sub-agent child rows attach to their matching parent `spawn_subagent`, not the last rendered sub-agent.
- A tool with multiple raw JSON parameters renders all parameters, including values omitted from a shortened `args_preview`.
- Sub-agent meta forwarding strips child `tool_call_count`.
- Existing sub-agent rendering still hides result previews and keeps status-column alignment.
