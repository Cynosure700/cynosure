## 2026-06-13T09:15:36Z · bash · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_HsSfgrXn9zX3lthChQhe6387

### Arguments

```json
{"command": "find /Users/bytedance/golang_pro/nano_cc/go-agent -type f -name \"*.go\" -o -name \"*.json\" -o -name \"*.yaml\" -o -name \"*.yml\" -o -name \"*.md\" -o -name \"go.mod\" -o -name \"go.sum\" | head -80"}
```

### Result

```text
/Users/bytedance/golang_pro/nano_cc/go-agent/TUI化改造设计文档.md
/Users/bytedance/golang_pro/nano_cc/go-agent/LINK_MD系统提示词注入设计文档.md
/Users/bytedance/golang_pro/nano_cc/go-agent/go.mod
/Users/bytedance/golang_pro/nano_cc/go-agent/TUI显示与事件循环修复设计文档.md
/Users/bytedance/golang_pro/nano_cc/go-agent/task_outputs/a56df192-a2a8-4e64-8132-15c6ba71a5b7/tools.md
/Users/bytedance/golang_pro/nano_cc/go-agent/task_outputs/d8905f1b-6e6b-4686-a0a6-6465a6458fad/tools.md
/Users/bytedance/golang_pro/nano_cc/go-agent/config.json
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/english-morning-reading/evals.json
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/english-morning-reading/SKILL.md
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/english-morning-reading-workspace/iteration-1/eval-2/eval_metadata.json
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/english-morning-reading-workspace/iteration-1/eval-2/without_skill/outputs/result.md
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/english-morning-reading-workspace/iteration-1/eval-2/with_skill/outputs/result.md
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/english-morning-reading-workspace/iteration-1/eval-3/eval_metadata.json
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/english-morning-reading-workspace/iteration-1/eval-3/without_skill/outputs/result.md
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/english-morning-reading-workspace/iteration-1/eval-3/with_skill/outputs/result.md
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/english-morning-reading-workspace/iteration-1/eval-1/eval_metadata.json
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/english-morning-reading-workspace/iteration-1/eval-1/without_skill/outputs/result.md
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/english-morning-reading-workspace/iteration-1/eval-1/with_skill/outputs/result.md
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/config.json
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/basic.go
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/skills/english-morning-reading/SKILL.md
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/skills/pdf-generator/evals/evals.json
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/skills/pdf-generator/SKILL.md
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/skills/weekly-report-generator/evals/evals.json
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/skills/weekly-report-generator/references/best-practices.md
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/skills/weekly-report-generator/SKILL.md
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/skills/code-assistant/SKILL.md
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/skills/skill-creator/references/schemas.md
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/skills/skill-creator/agents/grader.md
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/skills/skill-creator/agents/comparator.md
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/skills/skill-creator/agents/analyzer.md
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/skills/skill-creator/SKILL.md
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/skills/writing-assistant/SKILL.md
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/system_prompt.md
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/assistant/prompt_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/assistant/prompt.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/definitions.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/persisted_output_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/file_ops.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/handlers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/bash.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/definitions_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/todo_write.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/persisted_output.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/load_skill.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/runtime_env.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/validation.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/runtime_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/todo_write_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/path_guard.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/validation_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/llm/client.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/llm/errors.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/logger/logger.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/logger/logger_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/idgen/idgen.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/paths.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/local_config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/local_config_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/helpers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/events_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/events.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/app.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/memory.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/prompt_builder.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/context_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/subagent.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/runtime.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hook_bridge.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/tool_registry.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/summarizer.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/types.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/user_prompt.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/stop.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/manager.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/tool.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/compression_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/full_history_summarization.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/conversation_memory_strategy.go
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","resolved_command_path":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"/Users/bytedance/golang_pro/nano_cc/go-agent/TUI化改造设计文档.md\n/Users/bytedance/golang_pro/nano_cc/go-agent/LINK_MD系统提示词注入设计文档.md\n/Users/bytedance/golang_pro/nano_cc/go-agent/go.mod\n/Users/bytedance/golang_pro/nano_cc/go-agent/TUI显示与事件循环修复设计文档.md\n/Users/bytedance/golang_pro/nano_cc/go-agent/task_outputs/a56df192-a2a8-4e64-8132-15c6ba71a5b7/tools.md\n/Users/bytedance/golang_pro/nano_cc/go-agent/task_outputs/d8905f1b-6e6b-4686-a0a6-6465a6458"}
```

## 2026-06-13T09:15:39Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_N9tJb2UyYr4zpabo2QjW7514

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/main.go"}
```

### Result

```text
package main

import (
	"log"

	"nano_cc/internal/cli"
)

func main() {
	if err := cli.Main(); err != nil {
		log.Fatal(err)
	}
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package main\n\nimport (\n\t\"log\"\n\n\t\"nano_cc/internal/cli\"\n)\n\nfunc main() {\n\tif err := cli.Main(); err != nil {\n\t\tlog.Fatal(err)\n\t}\n}\n"}
```

## 2026-06-13T09:15:39Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_DTrXicBm2ADzVG9tEmdG1233

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/go.mod"}
```

### Result

```text
module nano_cc

go 1.26.1

require (
	github.com/charmbracelet/bubbles v1.0.0
	github.com/charmbracelet/bubbletea v1.3.10
	github.com/charmbracelet/glamour v1.0.0
	github.com/charmbracelet/lipgloss v1.1.1-0.20250404203927-76690c660834
	github.com/google/jsonschema-go v0.4.3
	github.com/modelcontextprotocol/go-sdk v1.6.1
	github.com/sashabaranov/go-openai v1.41.2
)

require (
	github.com/alecthomas/chroma/v2 v2.20.0 // indirect
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/charmbracelet/colorprofile v0.4.1 // indirect
	github.com/charmbracelet/x/ansi v0.11.6 // indirect
	github.com/charmbracelet/x/cellbuf v0.0.15 // indirect
	github.com/charmbracelet/x/exp/slice v0.0.0-20250327172914-2fdc97757edf // indirect
	github.com/charmbracelet/x/term v0.2.2 // indirect
	github.com/clipperhouse/displaywidth v0.9.0 // indirect
	github.com/clipperhouse/stringish v0.1.1 // indirect
	github.com/clipperhouse/uax29/v2 v2.5.0 // indirect
	github.com/dlclark/regexp2 v1.11.5 // indirect
	github.com/erikgeiser/coninput v0.0.0-20211004153227-1c3628e74d0f // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.3.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.19 // indirect
	github.com/microcosm-cc/bluemonday v1.0.27 // indirect
	github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/muesli/termenv v0.16.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/segmentio/asm v1.1.3 // indirect
	github.com/segmentio/encoding v0.5.4 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	github.com/yuin/goldmark v1.7.13 // indirect
	github.com/yuin/goldmark-emoji v1.0.6 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/oauth2 v0.35.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/term v0.36.0 // indirect
	golang.org/x/text v0.30.0 // indirect
)

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"module nano_cc\n\ngo 1.26.1\n\nrequire (\n\tgithub.com/charmbracelet/bubbles v1.0.0\n\tgithub.com/charmbracelet/bubbletea v1.3.10\n\tgithub.com/charmbracelet/glamour v1.0.0\n\tgithub.com/charmbracelet/lipgloss v1.1.1-0.20250404203927-76690c660834\n\tgithub.com/google/jsonschema-go v0.4.3\n\tgithub.com/modelcontextprotocol/go-sdk v1.6.1\n\tgithub.com/sashabaranov/go-openai v1.41.2\n)\n\nrequire (\n\tgithub.com/alecthomas/chroma/v2 v2.20.0 // indirect\n\tgithub.com/atotto/clipboard v0.1.4 // indirect\n\tgithub.com/aymanbaga"}
```

## 2026-06-13T09:15:39Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_pE7Z53kb3s3U4ersMYXZ3442

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/config.json"}
```

### Result

```text
{
    "base_url": "https://api.deepseek.com",
    "model_id": "deepseek-v4-flash",
    "app_home": ".",
    "system_prompt_path": "system_prompt.md",
    "workspace_root": "workspace",
    "builtin_skills_dir": "skills",
    "command_bin_dir": "bin",
    "command_script_dir": "cmd",
    "allowed_tools": "load_skill,bash,read_file,write_file,edit_file,spawn_subagent",
    "bash_allow_outside_workspace": false,
    "bash_allow_dangerous_commands": false
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"{\n    \"base_url\": \"https://api.deepseek.com\",\n    \"model_id\": \"deepseek-v4-flash\",\n    \"app_home\": \".\",\n    \"system_prompt_path\": \"system_prompt.md\",\n    \"workspace_root\": \"workspace\",\n    \"builtin_skills_dir\": \"skills\",\n    \"command_bin_dir\": \"bin\",\n    \"command_script_dir\": \"cmd\",\n    \"allowed_tools\": \"load_skill,bash,read_file,write_file,edit_file,spawn_subagent\",\n    \"bash_allow_outside_workspace\": false,\n    \"bash_allow_dangerous_commands\": false\n}\n"}
```

## 2026-06-13T09:15:40Z · bash · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_lFHIPPgWXHlZPM0qztQN8747

### Arguments

```json
{"command": "find /Users/bytedance/golang_pro/nano_cc/go-agent -type f -name \"*.go\" | sort"}
```

### Result

```text
/Users/bytedance/golang_pro/nano_cc/go-agent/architecture_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/transport_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/transport.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/compression_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/conversation_memory_strategy.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/full_history_summarization.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/message_window_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/reactive_compact.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/recent_tool_result_retention.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/token_estimator.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/tool_result_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/context_compression_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/context_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/conversation_flow.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/conversation_memory.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/helpers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hook_bridge.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/manager.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/stop.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/tool.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/types.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/user_prompt.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/memory_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/memory.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/prompt_builder.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/runtime_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/runtime.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/subagent.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/summarizer.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/todo_reminder.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/tool_registry.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/conversation_history_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/conversation_history.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/models.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/assistant/prompt_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/assistant/prompt.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/cli/root_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/cli/root.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/helpers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/local_config_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/local_config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/paths.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/idgen/idgen.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/llm/client.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/llm/errors.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/bootstrap.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/memory_store.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/persisted_output_files.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/session_history_store.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/store_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/store.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/tool_result_log.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/logger/logger_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/logger/logger.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/safety/path.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/sessions/render_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/sessions/skill_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/sessions/skill.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/textutil/textutil.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/bash.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/definitions_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/definitions.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/file_ops.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/handlers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/load_skill.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/path_guard.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/persisted_output_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/persisted_output.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/runtime_env.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/runtime_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/todo_write_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/todo_write.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/validation_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/validation.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/app.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/events_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/events.go
/Users/bytedance/golang_pro/nano_cc/go-agent/main.go
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/basic.go
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","resolved_command_path":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"/Users/bytedance/golang_pro/nano_cc/go-agent/architecture_test.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config_test.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager_test.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/transport_test.go\n/Users/bytedance/golang_pro/nano_cc/go-agent"}
```

## 2026-06-13T09:15:43Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_8bK7qMDJmbf8MpLNzj1Y1893

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/runtime.go"}
```

### Result

```text
package runtime

import (
	"context"
	"time"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/runtime/compression"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/config"
	"nano_cc/internal/llm"
	"nano_cc/internal/sessions"
)

type EventWriter interface {
	Event(name string, data any) error
}

type conversationStore interface {
	UpdateConversationTitle(ctx context.Context, conversationID, title string) error
	TouchConversationActivity(ctx context.Context, conversationID string) error
	SetConversationHistory(ctx context.Context, conversationID string, messages []storage.Message) error
	SetConversationCache(ctx context.Context, conversationID string, messages []storage.Message) error
	GetConversationCache(ctx context.Context, conversationID string) ([]storage.Message, bool, error)
	ListMessagesByConversation(ctx context.Context, conversationID string, limit int) ([]storage.Message, error)
	CreateToolCall(ctx context.Context, tc storage.ToolCall) error
	CreateSubagentMessage(ctx context.Context, message storage.SubagentMessage) error
	CreatePersistedOutput(ctx context.Context, output storage.PersistedOutput) error
	GetPersistedOutputForConversation(ctx context.Context, id, userID, conversationID string) (storage.PersistedOutput, error)
	GetPersistedOutputByMessageHash(ctx context.Context, conversationID, userID, messageID, toolCallID, strategy, contentSHA256 string) (storage.PersistedOutput, error)
	CreateContextSummary(ctx context.Context, summary storage.ContextSummary) error
	GetContextSummaryByHistoryHash(ctx context.Context, conversationID, userID, sourceHistorySHA256 string) (storage.ContextSummary, error)
	ListRelevantMemories(ctx context.Context, userID string) ([]storage.Memory, error)
	ListMemoriesByUserAndType(ctx context.Context, userID, memType string) ([]storage.Memory, error)
	ListProjectFactMemories(ctx context.Context, userID string) ([]storage.Memory, error)
	InsertMemory(ctx context.Context, m storage.Memory) error
	CountMemoriesByUserAndType(ctx context.Context, userID, memType string) (int, error)
	CountProjectFactMemories(ctx context.Context, userID string) (int, error)
	DeleteOldestMemories(ctx context.Context, userID, memType string, n int) error
	ReplaceMemoriesByUserAndType(ctx context.Context, userID, memType string, items []storage.Memory) error
	ReplaceProjectFactMemories(ctx context.Context, userID string, items []storage.Memory) error
	ListConversationMemories(ctx context.Context, conversationID string) ([]storage.ConversationMemory, error)
	ReplaceConversationMemories(ctx context.Context, conversationID, userID string, items []storage.ConversationMemory) error
	GetConversationModelHistory(ctx context.Context, conversationID string) ([]storage.Message, bool, error)
	UpsertConversationModelHistory(ctx context.Context, conversationID, userID string, messages []storage.Message) error
	AcquireConversationLock(ctx context.Context, conversationID, token string, ttl, waitTimeout time.Duration) (bool, error)
	RenewConversationLock(ctx context.Context, conversationID, token string, ttl time.Duration) (bool, error)
	ReleaseConversationLock(ctx context.Context, conversationID, token string) error
}

type Service struct {
	Store             conversationStore
	Cfg               config.AppConfig
	LLM               llm.Client
	Tools             *ToolRegistry
	BuiltinSkills     *sessions.SkillLoader
	BasePrompt        string
	LinkMarkdown      config.LinkMarkdownContext
	Hooks             *HookManager
	ContextCompressor *compression.Compressor
	EnableMemory      bool
	MCP               *mcp.Manager
}

func NewService(store conversationStore, cfg config.AppConfig, client llm.Client) *Service {
	return &Service{Store: store, Cfg: cfg, LLM: client, Tools: NewToolRegistry(cfg), Hooks: NewDefaultHookManager(), EnableMemory: true}
}

func (s *Service) hookManager() *HookManager {
	if s.Hooks == nil {
		s.Hooks = NewDefaultHookManager()
	}
	return s.Hooks
}

func (s *Service) SetBuiltinSkills(loader *sessions.SkillLoader) {
	s.BuiltinSkills = loader
}

func (s *Service) SetBasePrompt(prompt string) {
	s.BasePrompt = prompt
}

func (s *Service) SetLinkMarkdownContext(ctx config.LinkMarkdownContext) {
	s.LinkMarkdown = ctx
}

func (s *Service) SetMCPManager(manager *mcp.Manager) {
	s.MCP = manager
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"context\"\n\t\"time\"\n\n\t\"nano_cc/internal/agent/mcp\"\n\t\"nano_cc/internal/agent/runtime/compression\"\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/config\"\n\t\"nano_cc/internal/llm\"\n\t\"nano_cc/internal/sessions\"\n)\n\ntype EventWriter interface {\n\tEvent(name string, data any) error\n}\n\ntype conversationStore interface {\n\tUpdateConversationTitle(ctx context.Context, conversationID, title string) error\n\tTouchConversationActivity(ctx context.Context, conversationID string) error"}
```

## 2026-06-13T09:15:43Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_K4Dhp0JOzYb9cJwjjj0k2928

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/app.go"}
```

### Result

```text
package tui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/runtime"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/logger"
	"nano_cc/internal/sessions"
)

type SessionInfo struct {
	User         storage.User
	Conversation storage.Conversation
	CWD          string
	Resumer      SessionResumer
	Skills       []sessions.SkillSummary
	MCPServers   []mcp.ServerStatus
	SkillCount   int
	MCPToolCount int
}

type SessionResumer interface {
	ListResumableSessions(ctx context.Context, workspaceRoot string) ([]storage.ResumableSession, error)
	ResumeSession(ctx context.Context, sessionID, currentWorkspace string, user storage.User) (storage.Conversation, []storage.Message, error)
}

type Message struct {
	Role             string
	Content          string
	ReasoningContent string
}

type palette struct {
	ink      lipgloss.Color
	muted    lipgloss.Color
	panel    lipgloss.Color
	cyan     lipgloss.Color
	mint     lipgloss.Color
	lavender lipgloss.Color
	butter   lipgloss.Color
	coral    lipgloss.Color
}

var tuiPalette = palette{
	ink:      lipgloss.Color("255"),
	muted:    lipgloss.Color("244"),
	panel:    lipgloss.Color("238"),
	cyan:     lipgloss.Color("81"),
	mint:     lipgloss.Color("120"),
	lavender: lipgloss.Color("183"),
	butter:   lipgloss.Color("229"),
	coral:    lipgloss.Color("210"),
}

const inputCursor = "█"

type Model struct {
	runtime          *runtime.Service
	session          SessionInfo
	messages         []Message
	input            textarea.Model
	viewport         viewport.Model
	width            int
	height           int
	running          bool
	events           chan Event
	cancel           context.CancelFunc
	renderer         *glamour.TermRenderer
	generation       int64
	resumeSelecting  bool
	resumeCandidates []storage.ResumableSession
	toolCallCount    int
	contextTokens    int
	contextBudget    int
}

func NewModel(runtimeService *runtime.Service, session SessionInfo) Model {
	input := textarea.New()
	input.Placeholder = "问 go-agent 一件事..."
	input.Focus()
	input.SetHeight(2)
	input.ShowLineNumbers = false
	vp := viewport.New(100, 20)
	renderer := newMarkdownRenderer(100)
	return Model{runtime: runtimeService, session: session, input: input, viewport: vp, width: 100, height: 20, events: make(chan Event, 128), renderer: renderer}
}

func Run(ctx context.Context, runtimeService *runtime.Service, session SessionInfo) error {
	previousConsole := logger.SetConsoleEnabled(false)
	defer logger.SetConsoleEnabled(previousConsole)
	program := tea.NewProgram(NewModel(runtimeService, session), tea.WithContext(ctx), tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := program.Run()
	return err
}

func (m Model) Init() tea.Cmd { return textarea.Blink }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = m.viewportHeight()
		m.input.SetWidth(max(20, msg.Width-4))
		m.renderer = newMarkdownRenderer(m.messageWidth())
		m.refreshViewport()
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.running && m.cancel != nil {
				m.cancel()
				m.generation++
				m.running = false
				m.appendMessage("system", "已中断当前生成")
				return m, nil
			}
			return m, tea.Quit
		case "enter":
			text := strings.TrimSpace(m.input.Value())
			if text == "" || m.running {
				return m, nil
			}
			m.input.Reset()
			if m.resumeSelecting && m.handleResumeSelection(text) {
				m.refreshViewport()
				return m, nil
			}
			if strings.HasPrefix(text, "/") && m.handleSlashCommand(text) {
				m.refreshViewport()
				return m, nil
			}
			m.appendMessage("user", text)
			m.running = true
			m.toolCallCount = 0
			m.generation++
			generation := m.generation
			ctx, cancel := context.WithCancel(context.Background())
			m.cancel = cancel
			return m, tea.Batch(m.waitEvent(), m.respond(ctx, text, generation))
		}
		if isTerminalProbeResponseInput(msg) {
			return m, nil
		}
		if cmd, ok := m.updateViewportScroll(msg); ok {
			return m, cmd
		}
	case tea.MouseMsg:
		if cmd, ok := m.updateViewportScroll(msg); ok {
			return m, cmd
		}
	case Event:
		if msg.Generation != 0 && msg.Generation != m.generation {
			if m.running {
				return m, m.waitEvent()
			}
			return m, nil
		}
		switch msg.Name {
		case "assistant_delta":
			m.appendAssistantDelta(msg.Content)
		case "reasoning_delta":
			m.appendThinkingDelta(msg.Content)
		case "assistant":
			m.updateMetaFromData(msg.Data)
			content := msg.Content
			if content == "" && msg.Data != nil {
				content = eventContent(msg.Data)
			}
			if content != "" {
				m.replaceLastAssistant(content, eventString(msg.Data, "reasoning_content"))
			}
		case "meta":
			m.updateMetaFromData(msg.Data)
		case "error":
			m.appendMessage("error", msg.Content)
			m.running = false
		case "done":
			m.running = false
		}
		m.refreshViewport()
		if m.running {
			return m, m.waitEvent()
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.refreshViewport()
	return m, cmd
}

func (m *Model) updateViewportScroll(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyPgUp, tea.KeyPgDown, tea.KeyUp, tea.KeyDown:
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return cmd, true
		}
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress {
			switch msg.Button {
			case tea.MouseButtonWheelUp, tea.MouseButtonWheelDown, tea.MouseButtonWheelLeft, tea.MouseButtonWheelRight:
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				return cmd, true
			}
		}
	}
	return nil, false
}

func (m Model) View() string {
	return m.viewport.View()
}

func (m Model) respond(ctx context.Context, text string, generation int64) tea.Cmd {
	return func() tea.Msg {
		if m.runtime == nil {
			m.events <- Event{Generation: generation, Name: "error", Content: "runtime 未初始化"}
			return nil
		}
		_, err := m.runtime.RespondToConversation(ctx, m.session.Conversation, m.session.User, text, NewEventWriter(m.events, generation))
		if err != nil {
			m.events <- Event{Generation: generation, Name: "error", Content: err.Error()}
			return nil
		}
		m.events <- Event{Generation: generation, Name: "done"}
		return nil
	}
}

func (m Model) waitEvent() tea.Cmd {
	return func() tea.Msg { return <-m.events }
}

func (m *Model) handleSlashCommand(text string) bool {
	switch strings.TrimSpace(text) {
	case "/help":
		m.appendMessage("system", "命令：/help /clear /cwd /skills /mcp /resume。Enter 发送，Ctrl+C 中断或退出。")
		return true
	case "/clear":
		m.messages = nil
		m.resumeSelecting = false
		m.resumeCandidates = nil
		m.appendMessage("system", "已清空当前 TUI 显示上下文")
		return true
	case "/cwd":
		m.appendMessage("system", "当前工作区："+m.session.CWD)
		return true
	case "/skills":
		m.appendMessage("system", renderSkillDetails(m.session.Skills, m.session.SkillCount))
		return true
	case "/mcp":
		m.appendMessage("system", renderMCPDetails(m.session.MCPServers, m.session.MCPToolCount))
		return true
	case "/resume":
		m.startResumeSelection()
		return true
	}
	m.appendMessage("system", "未知命令："+text)
	return true
}

func (m *Model) startResumeSelection() {
	if m.running {
		m.appendMessage("system", "当前正在生成，请先 Ctrl+C 中断后再执行 /resume")
		return
	}
	if m.session.Resumer == nil {
		m.appendMessage("system", "当前运行环境不支持 /resume")
		return
	}
	sessions, err := m.session.Resumer.ListResumableSessions(context.Background(), m.session.CWD)
	if err != nil {
		m.appendMessage("error", err.Error())
		return
	}
	if len(sessions) == 0 {
		m.resumeSelecting = false
		m.resumeCandidates = nil
		m.appendMessage("system", "当前目录暂无可恢复的历史会话")
		return
	}
	m.resumeSelecting = true
	m.resumeCandidates = sessions
	m.appendMessage("system", renderResumeCandidates(sessions))
}

func (m *Model) handleResumeSelection(text string) bool {
	text = strings.TrimSpace(text)
	if text == "/cancel" || text == "/clear" {
		m.resumeSelecting = false
		m.resumeCandidates = nil
		if text == "/clear" {
			m.messages = nil
			m.appendMessage("system", "已清空当前 TUI 显示上下文")
		} else {
			m.appendMessage("system", "已取消恢复历史会话")
		}
		return true
	}
	idx, err := strconv.Atoi(text)
	if err != nil || idx < 1 || idx > len(m.resumeCandidates) {
		m.appendMessage("system", fmt.Sprintf("请输入 1-%d 之间的序号，或输入 /cancel 取消", len(m.resumeCandidates)))
		return true
	}
	candidate := m.resumeCandidates[idx-1]
	conv, history, err := m.session.Resumer.ResumeSession(context.Background(), candidate.SessionID, m.session.CWD, m.session.User)
	if err != nil {
		m.appendMessage("error", err.Error())
		return true
	}
	m.session.Conversation = conv
	m.resumeSelecting = false
	m.resumeCandidates = nil
	m.messages = messagesForDisplay(history)
	m.appendMessage("system", fmt.Sprintf("已恢复历史会话：%s", conv.SessionID))
	return true
}

func renderResumeCandidates(sessions []storage.ResumableSession) string {
	var b strings.Builder
	b.WriteString("可恢复的历史会话：")
	for i, session := range sessions {
		title := strings.TrimSpace(session.Title)
		if title == "" {
			title = "TUI 会话"
		}
		updated := "unknown"
		if !session.UpdatedAt.IsZero() {
			updated = session.UpdatedAt.Local().Format(time.RFC3339)
		}
		b.WriteString(fmt.Sprintf("\n%d. %s | %s | 消息:%d | %s", i+1, updated, title, session.MessageCount, session.SessionID))
	}
	b.WriteString("\n输入序号恢复，或输入 /cancel 取消。")
	return b.String()
}

func messagesForDisplay(history []storage.Message) []Message {
	if len(history) == 0 {
		return nil
	}
	messages := make([]Message, 0, len(history))
	for _, msg := range history {
		switch msg.Role {
		case "user", "assistant", "system", "error":
			messages = append(messages, Message{Role: msg.Role, Content: msg.Content, ReasoningContent: msg.ReasoningContent})
		}
	}
	return messages
}

func renderSkillDetails(skills []sessions.SkillSummary, fallbackCount int) string {
	count := len(skills)
	if count == 0 && fallbackCount > 0 {
		count = fallbackCount
	}
	if len(skills) == 0 {
		return fmt.Sprintf("已加载 Skills：%d 个", count)
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("已加载 Skills：%d 个", len(skills)))
	for _, skill := range skills {
		b.WriteString("\n- ")
		b.WriteString(skill.Name)
		if strings.TrimSpace(skill.Source) != "" {
			b.WriteString(" [")
			b.WriteString(skill.Source)
			b.WriteString("]")
		}
		if strings.TrimSpace(skill.Description) != "" {
			b.WriteString(" ")
			b.WriteString(skill.Description)
		}
		if strings.TrimSpace(skill.Path) != "" {
			b.WriteString("\n  path: ")
			b.WriteString(skill.Path)
		}
	}
	return b.String()
}

func renderMCPDetails(servers []mcp.ServerStatus, toolCount int) string {
	if len(servers) == 0 {
		return fmt.Sprintf("MCP Servers：0 个，工具：%d 个", toolCount)
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("MCP Servers：%d 个，工具：%d 个", len(servers), toolCount))
	for _, server := range servers {
		state := "disabled"
		if server.Connected {
			state = "connected"
		} else if server.LastError != "" {
			state = "failed"
		} else if server.Enabled {
			state = "pending"
		}
		b.WriteString("\n- ")
		b.WriteString(server.Name)
		b.WriteString(" [")
		b.WriteString(server.Transport)
		b.WriteString("] ")
		b.WriteString(state)
		b.WriteString(fmt.Sprintf(", tools: %d", server.ToolCount))
		if server.Command != "" {
			b.WriteString("\n  command: ")
			b.WriteString(strings.Join(append([]string{server.Command}, server.Args...), " "))
		}
		if server.URL != "" {
			b.WriteString("\n  url: ")
			b.WriteString(server.URL)
		}
		if server.LastError != "" {
			b.WriteString("\n  error: ")
			b.WriteString(server.LastError)
		}
	}
	return b.String()
}

func (m *Model) appendMessage(role, content string) {
	m.messages = append(m.messages, Message{Role: role, Content: content})
}

func (m *Model) appendAssistantDelta(delta string) {
	if len(m.messages) == 0 || !isLiveAssistantRole(m.messages[len(m.messages)-1].Role) {
		m.appendMessage("assistant", delta)
		return
	}
	if m.messages[len(m.messages)-1].Role == "thinking" {
		m.messages[len(m.messages)-1].Role = "assistant"
	}
	m.messages[len(m.messages)-1].Content += delta
}

func (m *Model) appendThinkingDelta(delta string) {
	if delta == "" {
		return
	}
	if len(m.messages) == 0 || !isLiveAssistantRole(m.messages[len(m.messages)-1].Role) {
		m.messages = append(m.messages, Message{Role: "assistant", ReasoningContent: delta})
		return
	}
	m.messages[len(m.messages)-1].ReasoningContent += delta
}

func (m *Model) replaceLastAssistant(content, reasoning string) {
	if len(m.messages) == 0 || !isLiveAssistantRole(m.messages[len(m.messages)-1].Role) {
		m.messages = append(m.messages, Message{Role: "assistant", Content: content, ReasoningContent: reasoning})
		return
	}
	m.messages[len(m.messages)-1].Role = "assistant"
	m.messages[len(m.messages)-1].Content = content
	m.messages[len(m.messages)-1].ReasoningContent = reasoning
}

func isLiveAssistantRole(role string) bool {
	return role == "assistant" || role == "thinking"
}

func (m *Model) refreshViewport() {
	m.viewport.Width = max(10, m.width)
	m.viewport.Height = m.viewportHeight()
	m.viewport.SetContent(m.renderTranscript())
	m.viewport.GotoBottom()
}

func (m Model) viewportHeight() int {
	if m.height <= 0 {
		return 20
	}
	return max(1, m.height)
}

func (m Model) renderTranscript() string {
	var b strings.Builder
	b.WriteString(m.renderWelcome())
	if len(m.messages) > 0 {
		b.WriteString("\n\n")
		b.WriteString(m.renderMessages())
	}
	b.WriteString("\n")
	b.WriteString(m.renderInput())
	b.WriteString("\n")
	b.WriteString(subtleStyle().Render(m.renderLiveStatus()))
	return b.String()
}

func (m Model) renderMessages() string {
	var b strings.Builder
	for _, msg := range m.messages {
		b.WriteString(m.renderMessage(msg))
		b.WriteString("\n\n")
	}
	return b.String()
}

func (m Model) renderMessage(msg Message) string {
	switch msg.Role {
	case "user":
		return promptLineStyle().Render("›") + " " + userStyle().Render(wrapText(msg.Content, m.messageWidth()-2))
	case "assistant":
		content := msg.Content
		if m.renderer != nil {
			if rendered, err := m.renderer.Render(content); err == nil {
				content = wrapText(strings.TrimSpace(rendered), m.messageWidth())
			}
		} else {
			content = wrapText(content, m.messageWidth())
		}
		if m.running && strings.TrimSpace(msg.ReasoningContent) != "" {
			content = thinkingStyle().Render("✽ 思考中\n"+wrapText(strings.TrimSpace(msg.ReasoningContent), m.messageWidth()-2)) + "\n" + content
		}
		return assistantLeadStyle().Render("go-agent") + "\n" + content
	case "thinking":
		return thinkingStyle().Render("✽ 思考中\n" + wrapText(msg.Content, m.messageWidth()-2))
	case "system":
		return systemStyle().Render("• " + wrapText(msg.Content, m.messageWidth()-2))
	case "error":
		return errorStyle().Render("✗ " + wrapText(msg.Content, m.messageWidth()-2))
	default:
		return roleLabel(msg.Role, lipgloss.Color("245")) + "\n" + wrapText(msg.Content, m.messageWidth())
	}
}

func (m Model) messageWidth() int {
	return max(10, m.width)
}

func newMarkdownRenderer(width int) *glamour.TermRenderer {
	renderer, _ := glamour.NewTermRenderer(glamour.WithStandardStyle("dark"), glamour.WithWordWrap(max(10, width)))
	return renderer
}

func wrapText(text string, width int) string {
	width = max(1, width)
	return ansi.Hardwrap(text, width, true)
}

func (m Model) renderWelcome() string {
	art := strings.Join([]string{
		`   /\_/\\`,
		`  ( o.o )   nano, but cozy`,
		`   > ^ <    ask · think · build`,
	}, "\n")
	intro := "像 Claude Code 一样一问一答：在终端直接提问，回答会流式显示在下方。"
	quick := "快捷键：Enter 发送 · Ctrl+C 中断/退出 · /resume 恢复 · /clear 清屏"
	stats := fmt.Sprintf("Skills %d · MCP tools %d", m.session.SkillCount, m.session.MCPToolCount)
	return startupPanelStyle().Width(max(20, m.width-2)).Render(accentArtStyle().Render(art) + "\n\n" + intro + "\n" + subtleStyle().Render(quick) + "\n" + subtleStyle().Render(stats))
}

func (m Model) renderHeader() string {
	width := max(20, m.width)
	status := runningText(m.running)
	contextText := "上下文 --"
	if m.contextBudget > 0 {
		contextText = fmt.Sprintf("上下文 %d%% · %s/%s", min(100, m.contextTokens*100/m.contextBudget), compactNumber(m.contextTokens), compactNumber(m.contextBudget))
	}
	left := titleStyle().Render("✦ go-agent")
	right := headerMetaStyle().Render(fmt.Sprintf("%s  ·  本轮工具 %d  ·  %s", status, m.toolCallCount, contextText))
	gap := max(1, width-lipgloss.Width(left)-lipgloss.Width(right)-4)
	line := left + strings.Repeat(" ", gap) + right
	workspace := subtleStyle().Render(fmt.Sprintf("cwd %s · skills %d · mcp tools %d", m.session.CWD, m.session.SkillCount, m.session.MCPToolCount))
	return headerStyle().Width(width).Render(line + "\n" + workspace)
}

func (m Model) renderConversationFrame() string {
	return conversationStyle().Width(max(10, m.width)).Render(m.viewport.View())
}

func (m Model) renderInput() string {
	prompt := inputPromptStyle().Render("›")
	text := strings.TrimSpace(m.input.Value())
	if text == "" {
		text = inputPromptStyle().Render(inputCursor) + " " + subtleStyle().Render(m.input.Placeholder)
	} else {
		text = userStyle().Render(text) + inputPromptStyle().Render(inputCursor)
	}
	return inputLineStyle().Width(max(10, m.width)).Render(prompt + " " + text)
}

func isTerminalProbeResponseInput(msg tea.KeyMsg) bool {
	if msg.Type != tea.KeyRunes {
		return false
	}
	value := string(msg.Runes)
	return strings.Contains(value, ";rgb:") && strings.Contains(value, "/")
}

func (m Model) renderLiveStatus() string {
	parts := []string{"Enter 发送", "Ctrl+C 中断/退出", "/help", fmt.Sprintf("工具 %d", m.toolCallCount)}
	if m.contextBudget > 0 {
		parts = append(parts, fmt.Sprintf("上下文 %d%% · %s/%s", min(100, m.contextTokens*100/m.contextBudget), compactNumber(m.contextTokens), compactNumber(m.contextBudget)))
	} else {
		parts = append(parts, "上下文 --")
	}
	return strings.Join(parts, " · ")
}

func (m *Model) updateMetaFromData(data any) {
	if count, ok := eventInt(data, "tool_call_count"); ok {
		m.toolCallCount = count
	}
	if tokens, ok := eventInt(data, "context_tokens"); ok {
		m.contextTokens = tokens
	}
	if budget, ok := eventInt(data, "context_budget"); ok {
		m.contextBudget = budget
	}
}

func eventString(data any, key string) string {
	m, ok := data.(map[string]any)
	if !ok {
		return ""
	}
	value, _ := m[key].(string)
	return value
}

func eventInt(data any, key string) (int, bool) {
	m, ok := data.(map[string]any)
	if !ok {
		return 0, false
	}
	switch value := m[key].(type) {
	case int:
		return value, true
	case int64:
		return int(value), true
	case float64:
		return int(value), true
	default:
		return 0, false
	}
}

func roleLabel(label string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Bold(true).Foreground(color).Render(label)
}

func headerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(tuiPalette.panel)
}

func titleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(tuiPalette.mint)
}

func headerMetaStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
}

func subtleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.muted)
}

func conversationStyle() lipgloss.Style {
	return lipgloss.NewStyle().Padding(1, 2, 0, 2)
}

func startupPanelStyle() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(tuiPalette.coral).Padding(1, 2).Margin(1, 0)
}

func userStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.ink)
}

func thinkingStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Italic(true).PaddingLeft(2)
}

func systemStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.butter)
}

func errorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.coral)
}

func promptLineStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.cyan).Bold(true)
}

func assistantLeadStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.mint).Bold(true)
}

func accentArtStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.lavender).Bold(true)
}

func inputPromptStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.mint).Bold(true)
}

func inputLineStyle() lipgloss.Style {
	return lipgloss.NewStyle().Background(lipgloss.Color("238")).Foreground(tuiPalette.ink)
}

func compactNumber(n int) string {
	if n >= 1000 {
		return fmt.Sprintf("%dk", n/1000)
	}
	return strconv.Itoa(n)
}

func runningText(running bool) string {
	if running {
		return "generating"
	}
	return "ready"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package tui\n\nimport (\n\t\"context\"\n\t\"fmt\"\n\t\"strconv\"\n\t\"strings\"\n\t\"time\"\n\n\t\"github.com/charmbracelet/bubbles/textarea\"\n\t\"github.com/charmbracelet/bubbles/viewport\"\n\ttea \"github.com/charmbracelet/bubbletea\"\n\t\"github.com/charmbracelet/glamour\"\n\t\"github.com/charmbracelet/lipgloss\"\n\t\"github.com/charmbracelet/x/ansi\"\n\n\t\"nano_cc/internal/agent/mcp\"\n\t\"nano_cc/internal/agent/runtime\"\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/logger\"\n\t\"nano_cc/internal/sessions\"\n)\n\ntype SessionInfo struct {\n\tUser"}
```

## 2026-06-13T09:15:43Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_DQ9G85MUxaknzamXt5Yf6308

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/cli/root.go"}
```

### Result

```text
package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"nano_cc/internal/local"
	"nano_cc/internal/tui"
)

type Options struct {
	CWD string
}

type Runner struct {
	RunTUI func(context.Context, Options) error
}

func Run(ctx context.Context, args []string, defaultCWD string, out io.Writer, runner Runner) error {
	if out == nil {
		out = io.Discard
	}
	if runner.RunTUI == nil {
		runner.RunTUI = runTUI
	}
	if len(args) > 0 && args[0] == "tui" {
		args = args[1:]
	}
	if len(args) > 0 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		printHelp(out)
		return nil
	}

	fs := flag.NewFlagSet("go-agent", flag.ContinueOnError)
	fs.SetOutput(out)
	cwd := fs.String("cwd", defaultCWD, "工作区目录，默认是当前目录")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("unknown command: %s", strings.Join(fs.Args(), " "))
	}
	return runner.RunTUI(ctx, Options{CWD: *cwd})
}

func Main() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	return Run(context.Background(), os.Args[1:], cwd, os.Stdout, Runner{})
}

func runTUI(ctx context.Context, opts Options) error {
	bundle, err := local.Bootstrap(ctx, opts.CWD)
	if err != nil {
		return err
	}
	defer bundle.Close()
	return tui.Run(ctx, bundle.Runtime, tui.SessionInfo{
		User:         bundle.User,
		Conversation: bundle.Conversation,
		CWD:          bundle.CWD,
		Resumer:      bundle.Store,
		Skills:       bundle.Skills,
		MCPServers:   bundle.MCPServers,
		SkillCount:   bundle.SkillCount,
		MCPToolCount: bundle.MCPToolCount,
	})
}

func printHelp(out io.Writer) {
	_, _ = fmt.Fprintln(out, `go-agent - 本地 TUI 代码助手

用法：
  go-agent [--cwd <path>]   启动 TUI，默认工作区为当前目录
  go-agent tui [--cwd <path>] 同上
  go-agent help             显示帮助`)
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package cli\n\nimport (\n\t\"context\"\n\t\"flag\"\n\t\"fmt\"\n\t\"io\"\n\t\"os\"\n\t\"strings\"\n\n\t\"nano_cc/internal/local\"\n\t\"nano_cc/internal/tui\"\n)\n\ntype Options struct {\n\tCWD string\n}\n\ntype Runner struct {\n\tRunTUI func(context.Context, Options) error\n}\n\nfunc Run(ctx context.Context, args []string, defaultCWD string, out io.Writer, runner Runner) error {\n\tif out == nil {\n\t\tout = io.Discard\n\t}\n\tif runner.RunTUI == nil {\n\t\trunner.RunTUI = runTUI\n\t}\n\tif len(args) \u003e 0 \u0026\u0026 args[0] == \"tui\" {\n\t\targs = args[1:]\n\t}\n\tif len(args"}
```

## 2026-06-13T09:15:43Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_03_ebbnornZmxuxnaHqcp6C6538

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/config.go"}
```

### Result

```text
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"-"`
	ModelID string `json:"model_id"`
}

type fileConfig struct {
	Config
	AppHome                            string `json:"app_home"`
	BuiltinSkillsDir                   string `json:"builtin_skills_dir"`
	CommandBinDir                      string `json:"command_bin_dir"`
	CommandScriptDir                   string `json:"command_script_dir"`
	SystemPromptPath                   string `json:"system_prompt_path"`
	WorkspaceRoot                      string `json:"workspace_root"`
	AllowedTools                       string `json:"allowed_tools"`
	BashAllowOutsideWorkspace          bool   `json:"bash_allow_outside_workspace"`
	BashAllowDangerousCommands         bool   `json:"bash_allow_dangerous_commands"`
	ConversationLockTTLSeconds         int    `json:"conversation_lock_ttl_seconds"`
	MemoryWorkTimeoutSeconds           int    `json:"memory_work_timeout_seconds"`
	ConversationLockWaitTimeoutSeconds int    `json:"conversation_lock_wait_timeout_seconds"`
}

type AppConfig struct {
	LLM                         Config
	AppHome                     string
	BuiltinSkillsDir            string
	CommandBinDir               string
	CommandScriptDir            string
	SystemPromptPath            string
	WorkspaceRoot               string
	LogsDir                     string
	AllowedTools                []string
	BashAllowOutsideWorkspace   bool
	BashAllowDangerousCommands  bool
	ConversationLockTTL         time.Duration
	MemoryWorkTimeout           time.Duration
	ConversationLockWaitTimeout time.Duration
}

func loadConfigFile() (fileConfig, error) {
	data, err := os.ReadFile(configFilePath())
	if err != nil {
		return fileConfig{}, err
	}

	var cfg fileConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fileConfig{}, fmt.Errorf("failed to parse workspace config.json: %w", err)
	}
	return cfg, nil
}

func configFilePath() string {
	return "config.json"
}

func loadLLMConfig(fileCfg fileConfig) (Config, error) {
	cfg := Config{
		BaseURL: strings.TrimSpace(fileCfg.Config.BaseURL),
		APIKey:  strings.TrimSpace(getenv("OPENAI_API_KEY")),
		ModelID: strings.TrimSpace(fileCfg.Config.ModelID),
	}

	if cfg.BaseURL == "" || cfg.ModelID == "" {
		return Config{}, fmt.Errorf("missing LLM config; set base_url, model_id in config.json")
	}
	if cfg.APIKey == "" {
		return Config{}, fmt.Errorf("missing LLM api key; set OPENAI_API_KEY environment variable")
	}

	return cfg, nil
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package config\n\nimport (\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"os\"\n\t\"strings\"\n\t\"time\"\n)\n\ntype Config struct {\n\tBaseURL string `json:\"base_url\"`\n\tAPIKey  string `json:\"-\"`\n\tModelID string `json:\"model_id\"`\n}\n\ntype fileConfig struct {\n\tConfig\n\tAppHome                            string `json:\"app_home\"`\n\tBuiltinSkillsDir                   string `json:\"builtin_skills_dir\"`\n\tCommandBinDir                      string `json:\"command_bin_dir\"`\n\tCommandScriptDir                   string `json:\"command_script_dir\"`\n"}
```

## 2026-06-13T09:15:43Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_04_bVmAMuShFGLnBALjfwBj1113

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/llm/client.go"}
```

### Result

```text
package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// Client is the minimal chat-completion interface the runtime depends on.
type Client interface {
	CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
	CreateChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (ChatCompletionStream, error)
}

// ChatCompletionStream is a streamed chat-completion response.
type ChatCompletionStream interface {
	Recv() (openai.ChatCompletionStreamResponse, error)
	Close() error
}

type deepseekClient struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

// NewDeepseekClient returns a Client backed by the DeepSeek chat-completions API.
func NewDeepseekClient(baseURL, apiKey string) Client {
	return &deepseekClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		http:    &http.Client{},
	}
}

// buildRequestBody assembles the DeepSeek request body, passing through
// max_tokens when set so callers can control/override the output budget.
func buildRequestBody(req openai.ChatCompletionRequest, stream bool) map[string]any {
	body := map[string]any{
		"model":            req.Model,
		"messages":         req.Messages,
		"thinking":         map[string]string{"type": "enabled"},
		"reasoning_effort": "high",
		"stream":           stream,
	}
	if len(req.Tools) > 0 {
		body["tools"] = req.Tools
	}
	if req.MaxTokens > 0 {
		body["max_tokens"] = req.MaxTokens
	}
	return body
}

// retryTransient runs do until it succeeds, returns a non-retryable error, or
// the transient-retry budget is exhausted. 429/529 are retried with
// exponential backoff plus jitter; ctx cancellation aborts immediately.
func retryTransient[T any](ctx context.Context, do func() (T, error)) (T, error) {
	var zero T
	for attempt := 0; ; attempt++ {
		result, err := do()
		if err == nil {
			return result, nil
		}
		var apiErr *APIError
		if !errors.As(err, &apiErr) || !isRetryableStatus(apiErr.StatusCode) || attempt >= maxTransientRetries {
			return zero, err
		}
		select {
		case <-time.After(backoffDelay(attempt)):
		case <-ctx.Done():
			return zero, ctx.Err()
		}
	}
}

func (c *deepseekClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	return retryTransient(ctx, func() (openai.ChatCompletionResponse, error) {
		return c.createChatCompletionOnce(ctx, req)
	})
}

func (c *deepseekClient) createChatCompletionOnce(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	jsonBody, err := json.Marshal(buildRequestBody(req, false))
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("marshal request: %w", err)
	}

	url := c.baseURL + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return openai.ChatCompletionResponse{}, &APIError{StatusCode: resp.StatusCode, Body: string(respBody)}
	}

	var result openai.ChatCompletionResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("unmarshal response: %w", err)
	}

	return result, nil
}

func (c *deepseekClient) CreateChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (ChatCompletionStream, error) {
	return retryTransient(ctx, func() (ChatCompletionStream, error) {
		return c.createChatCompletionStreamOnce(ctx, req)
	})
}

func (c *deepseekClient) createChatCompletionStreamOnce(ctx context.Context, req openai.ChatCompletionRequest) (ChatCompletionStream, error) {
	jsonBody, err := json.Marshal(buildRequestBody(req, true))
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := c.baseURL + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		respBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("read response: %w", readErr)
		}
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(respBody)}
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	return &deepseekChatStream{body: resp.Body, scanner: scanner}, nil
}

type deepseekChatStream struct {
	body    io.Closer
	scanner *bufio.Scanner
}

func (s *deepseekChatStream) Recv() (openai.ChatCompletionStreamResponse, error) {
	for s.scanner.Scan() {
		line := strings.TrimSpace(s.scanner.Text())
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			return openai.ChatCompletionStreamResponse{}, io.EOF
		}
		var chunk openai.ChatCompletionStreamResponse
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			return openai.ChatCompletionStreamResponse{}, fmt.Errorf("unmarshal stream response: %w", err)
		}
		return chunk, nil
	}
	if err := s.scanner.Err(); err != nil {
		return openai.ChatCompletionStreamResponse{}, err
	}
	return openai.ChatCompletionStreamResponse{}, io.EOF
}

func (s *deepseekChatStream) Close() error {
	return s.body.Close()
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package llm\n\nimport (\n\t\"bufio\"\n\t\"bytes\"\n\t\"context\"\n\t\"encoding/json\"\n\t\"errors\"\n\t\"fmt\"\n\t\"io\"\n\t\"net/http\"\n\t\"strings\"\n\t\"time\"\n\n\topenai \"github.com/sashabaranov/go-openai\"\n)\n\n// Client is the minimal chat-completion interface the runtime depends on.\ntype Client interface {\n\tCreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)\n\tCreateChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (ChatCompletionStream, error)\n}\n"}
```

## 2026-06-13T09:15:43Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_05_WvAmdODA5cpFcZKwtUgr7577

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/definitions.go"}
```

### Result

```text
package tools

import (
	"encoding/json"

	openai "github.com/sashabaranov/go-openai"
)

func mustMarshal(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func toolDef(name, desc string, params any) openai.Tool {
	return openai.Tool{
		Type: "function",
		Function: &openai.FunctionDefinition{
			Name:        name,
			Description: desc,
			Parameters:  mustMarshal(params),
		},
	}
}

func strParam(desc string) map[string]any {
	return map[string]any{
		"type":        "string",
		"description": desc,
	}
}

func intParam(desc string) map[string]any {
	return map[string]any{
		"type":        "integer",
		"description": desc,
	}
}

var baseToolDefs = []openai.Tool{
	toolDef("bash", "Execute a shell command via bash -c. Relative path arguments are interpreted under the workspace root; absolute paths outside the workspace and dangerous commands are rejected unless explicitly allowed by configuration.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": strParam("The shell command to execute"),
		},
		"required": []string{"command"},
	}),
	toolDef("read_file", "Read a file from the filesystem", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":  strParam("Path to the file to read"),
			"limit": intParam("Maximum number of lines to read"),
		},
		"required": []string{"path"},
	}),
	toolDef("write_file", "Write content to a file", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":    strParam("Path to the file to write"),
			"content": strParam("Content to write to the file"),
		},
		"required": []string{"path", "content"},
	}),
	toolDef("edit_file", "Replace text in a file by exact match", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":     strParam("Path to the file to edit"),
			"old_text": strParam("Exact text to find and replace"),
			"new_text": strParam("Text to replace it with"),
		},
		"required": []string{"path", "old_text", "new_text"},
	}),
	toolDef("load_skill", "Load the full instructions of a local skill by exact name before using or following that skill. Skills are loaded from the user's ~/.link/skills and the workspace .link/skills directories, with workspace skills taking precedence.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": strParam("Name of the skill to load"),
		},
		"required": []string{"name"},
	}),
	toolDef("todo_write", "Create or update the current task plan. Use this tool to track progress on multi-step tasks.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"todos": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"id":      map[string]any{"type": "string"},
						"content": map[string]any{"type": "string"},
						"status":  map[string]any{"type": "string", "enum": []string{TodoStatusPending, TodoStatusInProgress, TodoStatusCompleted}},
					},
					"required": []string{"id", "content", "status"},
				},
			},
		},
		"required": []string{"todos"},
	}),
}

var spawnSubagentToolDef = toolDef("spawn_subagent", "Spawn a child agent with a fresh message list to complete an isolated task. The child agent may use workspace tools, but it cannot spawn another subagent. Only its final summary is returned to the parent agent.", map[string]any{
	"type": "object",
	"properties": map[string]any{
		"task": strParam("The task for the child agent to complete. Include all context it needs because parent conversation history is not shared."),
		"cwd":  strParam("Optional working directory for the child agent. Relative paths are resolved under the workspace root; absolute paths must remain inside the workspace."),
	},
	"required": []string{"task"},
})

// ReadPersistedOutputToolName is exposed automatically alongside context
// compression so the model can fetch the full content behind a
// <persisted-output> marker when the inline preview is insufficient.
const ReadPersistedOutputToolName = "read_persisted_output"

var ReadPersistedOutputToolDef = toolDef(ReadPersistedOutputToolName, "Read a chunk of a persisted tool output by id when a <persisted-output> marker preview is insufficient. Only outputs from the current conversation are accessible.", map[string]any{
	"type": "object",
	"properties": map[string]any{
		"id":     strParam("The persisted output id from the <persisted-output> marker, for example po_abc123."),
		"offset": intParam("Zero-based character offset to start reading from. Defaults to 0."),
		"limit":  intParam("Maximum characters to return. Defaults to 20000 and is capped by the runtime."),
	},
	"required": []string{"id"},
})

var AllToolDefs = append(append([]openai.Tool(nil), baseToolDefs...), spawnSubagentToolDef, ReadPersistedOutputToolDef)
var ChildToolDefs = baseToolDefs

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package tools\n\nimport (\n\t\"encoding/json\"\n\n\topenai \"github.com/sashabaranov/go-openai\"\n)\n\nfunc mustMarshal(v any) json.RawMessage {\n\tb, err := json.Marshal(v)\n\tif err != nil {\n\t\tpanic(err)\n\t}\n\treturn b\n}\n\nfunc toolDef(name, desc string, params any) openai.Tool {\n\treturn openai.Tool{\n\t\tType: \"function\",\n\t\tFunction: \u0026openai.FunctionDefinition{\n\t\t\tName:        name,\n\t\t\tDescription: desc,\n\t\t\tParameters:  mustMarshal(params),\n\t\t},\n\t}\n}\n\nfunc strParam(desc string) map[string]any {\n\treturn map[string]an"}
```

## 2026-06-13T09:15:46Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_snDZEEyJci4KRPUbE0U11021

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/runtime.go"}
```

### Result

```text
package runtime

import (
	"context"
	"time"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/runtime/compression"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/config"
	"nano_cc/internal/llm"
	"nano_cc/internal/sessions"
)

type EventWriter interface {
	Event(name string, data any) error
}

type conversationStore interface {
	UpdateConversationTitle(ctx context.Context, conversationID, title string) error
	TouchConversationActivity(ctx context.Context, conversationID string) error
	SetConversationHistory(ctx context.Context, conversationID string, messages []storage.Message) error
	SetConversationCache(ctx context.Context, conversationID string, messages []storage.Message) error
	GetConversationCache(ctx context.Context, conversationID string) ([]storage.Message, bool, error)
	ListMessagesByConversation(ctx context.Context, conversationID string, limit int) ([]storage.Message, error)
	CreateToolCall(ctx context.Context, tc storage.ToolCall) error
	CreateSubagentMessage(ctx context.Context, message storage.SubagentMessage) error
	CreatePersistedOutput(ctx context.Context, output storage.PersistedOutput) error
	GetPersistedOutputForConversation(ctx context.Context, id, userID, conversationID string) (storage.PersistedOutput, error)
	GetPersistedOutputByMessageHash(ctx context.Context, conversationID, userID, messageID, toolCallID, strategy, contentSHA256 string) (storage.PersistedOutput, error)
	CreateContextSummary(ctx context.Context, summary storage.ContextSummary) error
	GetContextSummaryByHistoryHash(ctx context.Context, conversationID, userID, sourceHistorySHA256 string) (storage.ContextSummary, error)
	ListRelevantMemories(ctx context.Context, userID string) ([]storage.Memory, error)
	ListMemoriesByUserAndType(ctx context.Context, userID, memType string) ([]storage.Memory, error)
	ListProjectFactMemories(ctx context.Context, userID string) ([]storage.Memory, error)
	InsertMemory(ctx context.Context, m storage.Memory) error
	CountMemoriesByUserAndType(ctx context.Context, userID, memType string) (int, error)
	CountProjectFactMemories(ctx context.Context, userID string) (int, error)
	DeleteOldestMemories(ctx context.Context, userID, memType string, n int) error
	ReplaceMemoriesByUserAndType(ctx context.Context, userID, memType string, items []storage.Memory) error
	ReplaceProjectFactMemories(ctx context.Context, userID string, items []storage.Memory) error
	ListConversationMemories(ctx context.Context, conversationID string) ([]storage.ConversationMemory, error)
	ReplaceConversationMemories(ctx context.Context, conversationID, userID string, items []storage.ConversationMemory) error
	GetConversationModelHistory(ctx context.Context, conversationID string) ([]storage.Message, bool, error)
	UpsertConversationModelHistory(ctx context.Context, conversationID, userID string, messages []storage.Message) error
	AcquireConversationLock(ctx context.Context, conversationID, token string, ttl, waitTimeout time.Duration) (bool, error)
	RenewConversationLock(ctx context.Context, conversationID, token string, ttl time.Duration) (bool, error)
	ReleaseConversationLock(ctx context.Context, conversationID, token string) error
}

type Service struct {
	Store             conversationStore
	Cfg               config.AppConfig
	LLM               llm.Client
	Tools             *ToolRegistry
	BuiltinSkills     *sessions.SkillLoader
	BasePrompt        string
	LinkMarkdown      config.LinkMarkdownContext
	Hooks             *HookManager
	ContextCompressor *compression.Compressor
	EnableMemory      bool
	MCP               *mcp.Manager
}

func NewService(store conversationStore, cfg config.AppConfig, client llm.Client) *Service {
	return &Service{Store: store, Cfg: cfg, LLM: client, Tools: NewToolRegistry(cfg), Hooks: NewDefaultHookManager(), EnableMemory: true}
}

func (s *Service) hookManager() *HookManager {
	if s.Hooks == nil {
		s.Hooks = NewDefaultHookManager()
	}
	return s.Hooks
}

func (s *Service) SetBuiltinSkills(loader *sessions.SkillLoader) {
	s.BuiltinSkills = loader
}

func (s *Service) SetBasePrompt(prompt string) {
	s.BasePrompt = prompt
}

func (s *Service) SetLinkMarkdownContext(ctx config.LinkMarkdownContext) {
	s.LinkMarkdown = ctx
}

func (s *Service) SetMCPManager(manager *mcp.Manager) {
	s.MCP = manager
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"context\"\n\t\"time\"\n\n\t\"nano_cc/internal/agent/mcp\"\n\t\"nano_cc/internal/agent/runtime/compression\"\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/config\"\n\t\"nano_cc/internal/llm\"\n\t\"nano_cc/internal/sessions\"\n)\n\ntype EventWriter interface {\n\tEvent(name string, data any) error\n}\n\ntype conversationStore interface {\n\tUpdateConversationTitle(ctx context.Context, conversationID, title string) error\n\tTouchConversationActivity(ctx context.Context, conversationID string) error"}
```

## 2026-06-13T09:15:46Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_5KlMMmluIBBSlxTSwI4m6711

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/app.go"}
```

### Result

```text
package tui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/runtime"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/logger"
	"nano_cc/internal/sessions"
)

type SessionInfo struct {
	User         storage.User
	Conversation storage.Conversation
	CWD          string
	Resumer      SessionResumer
	Skills       []sessions.SkillSummary
	MCPServers   []mcp.ServerStatus
	SkillCount   int
	MCPToolCount int
}

type SessionResumer interface {
	ListResumableSessions(ctx context.Context, workspaceRoot string) ([]storage.ResumableSession, error)
	ResumeSession(ctx context.Context, sessionID, currentWorkspace string, user storage.User) (storage.Conversation, []storage.Message, error)
}

type Message struct {
	Role             string
	Content          string
	ReasoningContent string
}

type palette struct {
	ink      lipgloss.Color
	muted    lipgloss.Color
	panel    lipgloss.Color
	cyan     lipgloss.Color
	mint     lipgloss.Color
	lavender lipgloss.Color
	butter   lipgloss.Color
	coral    lipgloss.Color
}

var tuiPalette = palette{
	ink:      lipgloss.Color("255"),
	muted:    lipgloss.Color("244"),
	panel:    lipgloss.Color("238"),
	cyan:     lipgloss.Color("81"),
	mint:     lipgloss.Color("120"),
	lavender: lipgloss.Color("183"),
	butter:   lipgloss.Color("229"),
	coral:    lipgloss.Color("210"),
}

const inputCursor = "█"

type Model struct {
	runtime          *runtime.Service
	session          SessionInfo
	messages         []Message
	input            textarea.Model
	viewport         viewport.Model
	width            int
	height           int
	running          bool
	events           chan Event
	cancel           context.CancelFunc
	renderer         *glamour.TermRenderer
	generation       int64
	resumeSelecting  bool
	resumeCandidates []storage.ResumableSession
	toolCallCount    int
	contextTokens    int
	contextBudget    int
}

func NewModel(runtimeService *runtime.Service, session SessionInfo) Model {
	input := textarea.New()
	input.Placeholder = "问 go-agent 一件事..."
	input.Focus()
	input.SetHeight(2)
	input.ShowLineNumbers = false
	vp := viewport.New(100, 20)
	renderer := newMarkdownRenderer(100)
	return Model{runtime: runtimeService, session: session, input: input, viewport: vp, width: 100, height: 20, events: make(chan Event, 128), renderer: renderer}
}

func Run(ctx context.Context, runtimeService *runtime.Service, session SessionInfo) error {
	previousConsole := logger.SetConsoleEnabled(false)
	defer logger.SetConsoleEnabled(previousConsole)
	program := tea.NewProgram(NewModel(runtimeService, session), tea.WithContext(ctx), tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := program.Run()
	return err
}

func (m Model) Init() tea.Cmd { return textarea.Blink }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = m.viewportHeight()
		m.input.SetWidth(max(20, msg.Width-4))
		m.renderer = newMarkdownRenderer(m.messageWidth())
		m.refreshViewport()
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.running && m.cancel != nil {
				m.cancel()
				m.generation++
				m.running = false
				m.appendMessage("system", "已中断当前生成")
				return m, nil
			}
			return m, tea.Quit
		case "enter":
			text := strings.TrimSpace(m.input.Value())
			if text == "" || m.running {
				return m, nil
			}
			m.input.Reset()
			if m.resumeSelecting && m.handleResumeSelection(text) {
				m.refreshViewport()
				return m, nil
			}
			if strings.HasPrefix(text, "/") && m.handleSlashCommand(text) {
				m.refreshViewport()
				return m, nil
			}
			m.appendMessage("user", text)
			m.running = true
			m.toolCallCount = 0
			m.generation++
			generation := m.generation
			ctx, cancel := context.WithCancel(context.Background())
			m.cancel = cancel
			return m, tea.Batch(m.waitEvent(), m.respond(ctx, text, generation))
		}
		if isTerminalProbeResponseInput(msg) {
			return m, nil
		}
		if cmd, ok := m.updateViewportScroll(msg); ok {
			return m, cmd
		}
	case tea.MouseMsg:
		if cmd, ok := m.updateViewportScroll(msg); ok {
			return m, cmd
		}
	case Event:
		if msg.Generation != 0 && msg.Generation != m.generation {
			if m.running {
				return m, m.waitEvent()
			}
			return m, nil
		}
		switch msg.Name {
		case "assistant_delta":
			m.appendAssistantDelta(msg.Content)
		case "reasoning_delta":
			m.appendThinkingDelta(msg.Content)
		case "assistant":
			m.updateMetaFromData(msg.Data)
			content := msg.Content
			if content == "" && msg.Data != nil {
				content = eventContent(msg.Data)
			}
			if content != "" {
				m.replaceLastAssistant(content, eventString(msg.Data, "reasoning_content"))
			}
		case "meta":
			m.updateMetaFromData(msg.Data)
		case "error":
			m.appendMessage("error", msg.Content)
			m.running = false
		case "done":
			m.running = false
		}
		m.refreshViewport()
		if m.running {
			return m, m.waitEvent()
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.refreshViewport()
	return m, cmd
}

func (m *Model) updateViewportScroll(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyPgUp, tea.KeyPgDown, tea.KeyUp, tea.KeyDown:
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return cmd, true
		}
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress {
			switch msg.Button {
			case tea.MouseButtonWheelUp, tea.MouseButtonWheelDown, tea.MouseButtonWheelLeft, tea.MouseButtonWheelRight:
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				return cmd, true
			}
		}
	}
	return nil, false
}

func (m Model) View() string {
	return m.viewport.View()
}

func (m Model) respond(ctx context.Context, text string, generation int64) tea.Cmd {
	return func() tea.Msg {
		if m.runtime == nil {
			m.events <- Event{Generation: generation, Name: "error", Content: "runtime 未初始化"}
			return nil
		}
		_, err := m.runtime.RespondToConversation(ctx, m.session.Conversation, m.session.User, text, NewEventWriter(m.events, generation))
		if err != nil {
			m.events <- Event{Generation: generation, Name: "error", Content: err.Error()}
			return nil
		}
		m.events <- Event{Generation: generation, Name: "done"}
		return nil
	}
}

func (m Model) waitEvent() tea.Cmd {
	return func() tea.Msg { return <-m.events }
}

func (m *Model) handleSlashCommand(text string) bool {
	switch strings.TrimSpace(text) {
	case "/help":
		m.appendMessage("system", "命令：/help /clear /cwd /skills /mcp /resume。Enter 发送，Ctrl+C 中断或退出。")
		return true
	case "/clear":
		m.messages = nil
		m.resumeSelecting = false
		m.resumeCandidates = nil
		m.appendMessage("system", "已清空当前 TUI 显示上下文")
		return true
	case "/cwd":
		m.appendMessage("system", "当前工作区："+m.session.CWD)
		return true
	case "/skills":
		m.appendMessage("system", renderSkillDetails(m.session.Skills, m.session.SkillCount))
		return true
	case "/mcp":
		m.appendMessage("system", renderMCPDetails(m.session.MCPServers, m.session.MCPToolCount))
		return true
	case "/resume":
		m.startResumeSelection()
		return true
	}
	m.appendMessage("system", "未知命令："+text)
	return true
}

func (m *Model) startResumeSelection() {
	if m.running {
		m.appendMessage("system", "当前正在生成，请先 Ctrl+C 中断后再执行 /resume")
		return
	}
	if m.session.Resumer == nil {
		m.appendMessage("system", "当前运行环境不支持 /resume")
		return
	}
	sessions, err := m.session.Resumer.ListResumableSessions(context.Background(), m.session.CWD)
	if err != nil {
		m.appendMessage("error", err.Error())
		return
	}
	if len(sessions) == 0 {
		m.resumeSelecting = false
		m.resumeCandidates = nil
		m.appendMessage("system", "当前目录暂无可恢复的历史会话")
		return
	}
	m.resumeSelecting = true
	m.resumeCandidates = sessions
	m.appendMessage("system", renderResumeCandidates(sessions))
}

func (m *Model) handleResumeSelection(text string) bool {
	text = strings.TrimSpace(text)
	if text == "/cancel" || text == "/clear" {
		m.resumeSelecting = false
		m.resumeCandidates = nil
		if text == "/clear" {
			m.messages = nil
			m.appendMessage("system", "已清空当前 TUI 显示上下文")
		} else {
			m.appendMessage("system", "已取消恢复历史会话")
		}
		return true
	}
	idx, err := strconv.Atoi(text)
	if err != nil || idx < 1 || idx > len(m.resumeCandidates) {
		m.appendMessage("system", fmt.Sprintf("请输入 1-%d 之间的序号，或输入 /cancel 取消", len(m.resumeCandidates)))
		return true
	}
	candidate := m.resumeCandidates[idx-1]
	conv, history, err := m.session.Resumer.ResumeSession(context.Background(), candidate.SessionID, m.session.CWD, m.session.User)
	if err != nil {
		m.appendMessage("error", err.Error())
		return true
	}
	m.session.Conversation = conv
	m.resumeSelecting = false
	m.resumeCandidates = nil
	m.messages = messagesForDisplay(history)
	m.appendMessage("system", fmt.Sprintf("已恢复历史会话：%s", conv.SessionID))
	return true
}

func renderResumeCandidates(sessions []storage.ResumableSession) string {
	var b strings.Builder
	b.WriteString("可恢复的历史会话：")
	for i, session := range sessions {
		title := strings.TrimSpace(session.Title)
		if title == "" {
			title = "TUI 会话"
		}
		updated := "unknown"
		if !session.UpdatedAt.IsZero() {
			updated = session.UpdatedAt.Local().Format(time.RFC3339)
		}
		b.WriteString(fmt.Sprintf("\n%d. %s | %s | 消息:%d | %s", i+1, updated, title, session.MessageCount, session.SessionID))
	}
	b.WriteString("\n输入序号恢复，或输入 /cancel 取消。")
	return b.String()
}

func messagesForDisplay(history []storage.Message) []Message {
	if len(history) == 0 {
		return nil
	}
	messages := make([]Message, 0, len(history))
	for _, msg := range history {
		switch msg.Role {
		case "user", "assistant", "system", "error":
			messages = append(messages, Message{Role: msg.Role, Content: msg.Content, ReasoningContent: msg.ReasoningContent})
		}
	}
	return messages
}

func renderSkillDetails(skills []sessions.SkillSummary, fallbackCount int) string {
	count := len(skills)
	if count == 0 && fallbackCount > 0 {
		count = fallbackCount
	}
	if len(skills) == 0 {
		return fmt.Sprintf("已加载 Skills：%d 个", count)
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("已加载 Skills：%d 个", len(skills)))
	for _, skill := range skills {
		b.WriteString("\n- ")
		b.WriteString(skill.Name)
		if strings.TrimSpace(skill.Source) != "" {
			b.WriteString(" [")
			b.WriteString(skill.Source)
			b.WriteString("]")
		}
		if strings.TrimSpace(skill.Description) != "" {
			b.WriteString(" ")
			b.WriteString(skill.Description)
		}
		if strings.TrimSpace(skill.Path) != "" {
			b.WriteString("\n  path: ")
			b.WriteString(skill.Path)
		}
	}
	return b.String()
}

func renderMCPDetails(servers []mcp.ServerStatus, toolCount int) string {
	if len(servers) == 0 {
		return fmt.Sprintf("MCP Servers：0 个，工具：%d 个", toolCount)
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("MCP Servers：%d 个，工具：%d 个", len(servers), toolCount))
	for _, server := range servers {
		state := "disabled"
		if server.Connected {
			state = "connected"
		} else if server.LastError != "" {
			state = "failed"
		} else if server.Enabled {
			state = "pending"
		}
		b.WriteString("\n- ")
		b.WriteString(server.Name)
		b.WriteString(" [")
		b.WriteString(server.Transport)
		b.WriteString("] ")
		b.WriteString(state)
		b.WriteString(fmt.Sprintf(", tools: %d", server.ToolCount))
		if server.Command != "" {
			b.WriteString("\n  command: ")
			b.WriteString(strings.Join(append([]string{server.Command}, server.Args...), " "))
		}
		if server.URL != "" {
			b.WriteString("\n  url: ")
			b.WriteString(server.URL)
		}
		if server.LastError != "" {
			b.WriteString("\n  error: ")
			b.WriteString(server.LastError)
		}
	}
	return b.String()
}

func (m *Model) appendMessage(role, content string) {
	m.messages = append(m.messages, Message{Role: role, Content: content})
}

func (m *Model) appendAssistantDelta(delta string) {
	if len(m.messages) == 0 || !isLiveAssistantRole(m.messages[len(m.messages)-1].Role) {
		m.appendMessage("assistant", delta)
		return
	}
	if m.messages[len(m.messages)-1].Role == "thinking" {
		m.messages[len(m.messages)-1].Role = "assistant"
	}
	m.messages[len(m.messages)-1].Content += delta
}

func (m *Model) appendThinkingDelta(delta string) {
	if delta == "" {
		return
	}
	if len(m.messages) == 0 || !isLiveAssistantRole(m.messages[len(m.messages)-1].Role) {
		m.messages = append(m.messages, Message{Role: "assistant", ReasoningContent: delta})
		return
	}
	m.messages[len(m.messages)-1].ReasoningContent += delta
}

func (m *Model) replaceLastAssistant(content, reasoning string) {
	if len(m.messages) == 0 || !isLiveAssistantRole(m.messages[len(m.messages)-1].Role) {
		m.messages = append(m.messages, Message{Role: "assistant", Content: content, ReasoningContent: reasoning})
		return
	}
	m.messages[len(m.messages)-1].Role = "assistant"
	m.messages[len(m.messages)-1].Content = content
	m.messages[len(m.messages)-1].ReasoningContent = reasoning
}

func isLiveAssistantRole(role string) bool {
	return role == "assistant" || role == "thinking"
}

func (m *Model) refreshViewport() {
	m.viewport.Width = max(10, m.width)
	m.viewport.Height = m.viewportHeight()
	m.viewport.SetContent(m.renderTranscript())
	m.viewport.GotoBottom()
}

func (m Model) viewportHeight() int {
	if m.height <= 0 {
		return 20
	}
	return max(1, m.height)
}

func (m Model) renderTranscript() string {
	var b strings.Builder
	b.WriteString(m.renderWelcome())
	if len(m.messages) > 0 {
		b.WriteString("\n\n")
		b.WriteString(m.renderMessages())
	}
	b.WriteString("\n")
	b.WriteString(m.renderInput())
	b.WriteString("\n")
	b.WriteString(subtleStyle().Render(m.renderLiveStatus()))
	return b.String()
}

func (m Model) renderMessages() string {
	var b strings.Builder
	for _, msg := range m.messages {
		b.WriteString(m.renderMessage(msg))
		b.WriteString("\n\n")
	}
	return b.String()
}

func (m Model) renderMessage(msg Message) string {
	switch msg.Role {
	case "user":
		return promptLineStyle().Render("›") + " " + userStyle().Render(wrapText(msg.Content, m.messageWidth()-2))
	case "assistant":
		content := msg.Content
		if m.renderer != nil {
			if rendered, err := m.renderer.Render(content); err == nil {
				content = wrapText(strings.TrimSpace(rendered), m.messageWidth())
			}
		} else {
			content = wrapText(content, m.messageWidth())
		}
		if m.running && strings.TrimSpace(msg.ReasoningContent) != "" {
			content = thinkingStyle().Render("✽ 思考中\n"+wrapText(strings.TrimSpace(msg.ReasoningContent), m.messageWidth()-2)) + "\n" + content
		}
		return assistantLeadStyle().Render("go-agent") + "\n" + content
	case "thinking":
		return thinkingStyle().Render("✽ 思考中\n" + wrapText(msg.Content, m.messageWidth()-2))
	case "system":
		return systemStyle().Render("• " + wrapText(msg.Content, m.messageWidth()-2))
	case "error":
		return errorStyle().Render("✗ " + wrapText(msg.Content, m.messageWidth()-2))
	default:
		return roleLabel(msg.Role, lipgloss.Color("245")) + "\n" + wrapText(msg.Content, m.messageWidth())
	}
}

func (m Model) messageWidth() int {
	return max(10, m.width)
}

func newMarkdownRenderer(width int) *glamour.TermRenderer {
	renderer, _ := glamour.NewTermRenderer(glamour.WithStandardStyle("dark"), glamour.WithWordWrap(max(10, width)))
	return renderer
}

func wrapText(text string, width int) string {
	width = max(1, width)
	return ansi.Hardwrap(text, width, true)
}

func (m Model) renderWelcome() string {
	art := strings.Join([]string{
		`   /\_/\\`,
		`  ( o.o )   nano, but cozy`,
		`   > ^ <    ask · think · build`,
	}, "\n")
	intro := "像 Claude Code 一样一问一答：在终端直接提问，回答会流式显示在下方。"
	quick := "快捷键：Enter 发送 · Ctrl+C 中断/退出 · /resume 恢复 · /clear 清屏"
	stats := fmt.Sprintf("Skills %d · MCP tools %d", m.session.SkillCount, m.session.MCPToolCount)
	return startupPanelStyle().Width(max(20, m.width-2)).Render(accentArtStyle().Render(art) + "\n\n" + intro + "\n" + subtleStyle().Render(quick) + "\n" + subtleStyle().Render(stats))
}

func (m Model) renderHeader() string {
	width := max(20, m.width)
	status := runningText(m.running)
	contextText := "上下文 --"
	if m.contextBudget > 0 {
		contextText = fmt.Sprintf("上下文 %d%% · %s/%s", min(100, m.contextTokens*100/m.contextBudget), compactNumber(m.contextTokens), compactNumber(m.contextBudget))
	}
	left := titleStyle().Render("✦ go-agent")
	right := headerMetaStyle().Render(fmt.Sprintf("%s  ·  本轮工具 %d  ·  %s", status, m.toolCallCount, contextText))
	gap := max(1, width-lipgloss.Width(left)-lipgloss.Width(right)-4)
	line := left + strings.Repeat(" ", gap) + right
	workspace := subtleStyle().Render(fmt.Sprintf("cwd %s · skills %d · mcp tools %d", m.session.CWD, m.session.SkillCount, m.session.MCPToolCount))
	return headerStyle().Width(width).Render(line + "\n" + workspace)
}

func (m Model) renderConversationFrame() string {
	return conversationStyle().Width(max(10, m.width)).Render(m.viewport.View())
}

func (m Model) renderInput() string {
	prompt := inputPromptStyle().Render("›")
	text := strings.TrimSpace(m.input.Value())
	if text == "" {
		text = inputPromptStyle().Render(inputCursor) + " " + subtleStyle().Render(m.input.Placeholder)
	} else {
		text = userStyle().Render(text) + inputPromptStyle().Render(inputCursor)
	}
	return inputLineStyle().Width(max(10, m.width)).Render(prompt + " " + text)
}

func isTerminalProbeResponseInput(msg tea.KeyMsg) bool {
	if msg.Type != tea.KeyRunes {
		return false
	}
	value := string(msg.Runes)
	return strings.Contains(value, ";rgb:") && strings.Contains(value, "/")
}

func (m Model) renderLiveStatus() string {
	parts := []string{"Enter 发送", "Ctrl+C 中断/退出", "/help", fmt.Sprintf("工具 %d", m.toolCallCount)}
	if m.contextBudget > 0 {
		parts = append(parts, fmt.Sprintf("上下文 %d%% · %s/%s", min(100, m.contextTokens*100/m.contextBudget), compactNumber(m.contextTokens), compactNumber(m.contextBudget)))
	} else {
		parts = append(parts, "上下文 --")
	}
	return strings.Join(parts, " · ")
}

func (m *Model) updateMetaFromData(data any) {
	if count, ok := eventInt(data, "tool_call_count"); ok {
		m.toolCallCount = count
	}
	if tokens, ok := eventInt(data, "context_tokens"); ok {
		m.contextTokens = tokens
	}
	if budget, ok := eventInt(data, "context_budget"); ok {
		m.contextBudget = budget
	}
}

func eventString(data any, key string) string {
	m, ok := data.(map[string]any)
	if !ok {
		return ""
	}
	value, _ := m[key].(string)
	return value
}

func eventInt(data any, key string) (int, bool) {
	m, ok := data.(map[string]any)
	if !ok {
		return 0, false
	}
	switch value := m[key].(type) {
	case int:
		return value, true
	case int64:
		return int(value), true
	case float64:
		return int(value), true
	default:
		return 0, false
	}
}

func roleLabel(label string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Bold(true).Foreground(color).Render(label)
}

func headerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(tuiPalette.panel)
}

func titleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(tuiPalette.mint)
}

func headerMetaStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
}

func subtleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.muted)
}

func conversationStyle() lipgloss.Style {
	return lipgloss.NewStyle().Padding(1, 2, 0, 2)
}

func startupPanelStyle() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(tuiPalette.coral).Padding(1, 2).Margin(1, 0)
}

func userStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.ink)
}

func thinkingStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Italic(true).PaddingLeft(2)
}

func systemStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.butter)
}

func errorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.coral)
}

func promptLineStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.cyan).Bold(true)
}

func assistantLeadStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.mint).Bold(true)
}

func accentArtStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.lavender).Bold(true)
}

func inputPromptStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.mint).Bold(true)
}

func inputLineStyle() lipgloss.Style {
	return lipgloss.NewStyle().Background(lipgloss.Color("238")).Foreground(tuiPalette.ink)
}

func compactNumber(n int) string {
	if n >= 1000 {
		return fmt.Sprintf("%dk", n/1000)
	}
	return strconv.Itoa(n)
}

func runningText(running bool) string {
	if running {
		return "generating"
	}
	return "ready"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package tui\n\nimport (\n\t\"context\"\n\t\"fmt\"\n\t\"strconv\"\n\t\"strings\"\n\t\"time\"\n\n\t\"github.com/charmbracelet/bubbles/textarea\"\n\t\"github.com/charmbracelet/bubbles/viewport\"\n\ttea \"github.com/charmbracelet/bubbletea\"\n\t\"github.com/charmbracelet/glamour\"\n\t\"github.com/charmbracelet/lipgloss\"\n\t\"github.com/charmbracelet/x/ansi\"\n\n\t\"nano_cc/internal/agent/mcp\"\n\t\"nano_cc/internal/agent/runtime\"\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/logger\"\n\t\"nano_cc/internal/sessions\"\n)\n\ntype SessionInfo struct {\n\tUser"}
```

## 2026-06-13T09:15:46Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_wbLrlVoxBX9FDsFK26661007

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/cli/root.go"}
```

### Result

```text
package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"nano_cc/internal/local"
	"nano_cc/internal/tui"
)

type Options struct {
	CWD string
}

type Runner struct {
	RunTUI func(context.Context, Options) error
}

func Run(ctx context.Context, args []string, defaultCWD string, out io.Writer, runner Runner) error {
	if out == nil {
		out = io.Discard
	}
	if runner.RunTUI == nil {
		runner.RunTUI = runTUI
	}
	if len(args) > 0 && args[0] == "tui" {
		args = args[1:]
	}
	if len(args) > 0 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		printHelp(out)
		return nil
	}

	fs := flag.NewFlagSet("go-agent", flag.ContinueOnError)
	fs.SetOutput(out)
	cwd := fs.String("cwd", defaultCWD, "工作区目录，默认是当前目录")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("unknown command: %s", strings.Join(fs.Args(), " "))
	}
	return runner.RunTUI(ctx, Options{CWD: *cwd})
}

func Main() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	return Run(context.Background(), os.Args[1:], cwd, os.Stdout, Runner{})
}

func runTUI(ctx context.Context, opts Options) error {
	bundle, err := local.Bootstrap(ctx, opts.CWD)
	if err != nil {
		return err
	}
	defer bundle.Close()
	return tui.Run(ctx, bundle.Runtime, tui.SessionInfo{
		User:         bundle.User,
		Conversation: bundle.Conversation,
		CWD:          bundle.CWD,
		Resumer:      bundle.Store,
		Skills:       bundle.Skills,
		MCPServers:   bundle.MCPServers,
		SkillCount:   bundle.SkillCount,
		MCPToolCount: bundle.MCPToolCount,
	})
}

func printHelp(out io.Writer) {
	_, _ = fmt.Fprintln(out, `go-agent - 本地 TUI 代码助手

用法：
  go-agent [--cwd <path>]   启动 TUI，默认工作区为当前目录
  go-agent tui [--cwd <path>] 同上
  go-agent help             显示帮助`)
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package cli\n\nimport (\n\t\"context\"\n\t\"flag\"\n\t\"fmt\"\n\t\"io\"\n\t\"os\"\n\t\"strings\"\n\n\t\"nano_cc/internal/local\"\n\t\"nano_cc/internal/tui\"\n)\n\ntype Options struct {\n\tCWD string\n}\n\ntype Runner struct {\n\tRunTUI func(context.Context, Options) error\n}\n\nfunc Run(ctx context.Context, args []string, defaultCWD string, out io.Writer, runner Runner) error {\n\tif out == nil {\n\t\tout = io.Discard\n\t}\n\tif runner.RunTUI == nil {\n\t\trunner.RunTUI = runTUI\n\t}\n\tif len(args) \u003e 0 \u0026\u0026 args[0] == \"tui\" {\n\t\targs = args[1:]\n\t}\n\tif len(args"}
```

## 2026-06-13T09:15:46Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_03_a49SltJ0gVe3l8Yyqim04453

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/storage.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/storage.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/storage.go: no such file or directory"}
```

## 2026-06-13T09:15:48Z · bash · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_CXbWHZW8rnyLzCjmKO512350

### Arguments

```json
{"command": "find /Users/bytedance/golang_pro/nano_cc/go-agent -type f -name \"*.go\" | sort"}
```

### Result

```text
/Users/bytedance/golang_pro/nano_cc/go-agent/architecture_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/transport_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/transport.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/compression_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/conversation_memory_strategy.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/full_history_summarization.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/message_window_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/reactive_compact.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/recent_tool_result_retention.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/token_estimator.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/tool_result_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/context_compression_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/context_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/conversation_flow.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/conversation_memory.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/helpers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hook_bridge.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/manager.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/stop.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/tool.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/types.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/user_prompt.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/memory_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/memory.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/prompt_builder.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/runtime_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/runtime.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/subagent.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/summarizer.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/todo_reminder.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/tool_registry.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/conversation_history_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/conversation_history.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/models.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/assistant/prompt_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/assistant/prompt.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/cli/root_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/cli/root.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/helpers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/local_config_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/local_config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/paths.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/idgen/idgen.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/llm/client.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/llm/errors.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/bootstrap.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/memory_store.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/persisted_output_files.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/session_history_store.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/store_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/store.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/tool_result_log.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/logger/logger_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/logger/logger.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/safety/path.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/sessions/render_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/sessions/skill_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/sessions/skill.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/textutil/textutil.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/bash.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/definitions_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/definitions.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/file_ops.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/handlers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/load_skill.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/path_guard.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/persisted_output_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/persisted_output.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/runtime_env.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/runtime_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/todo_write_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/todo_write.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/validation_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/validation.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/app.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/events_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/events.go
/Users/bytedance/golang_pro/nano_cc/go-agent/main.go
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/basic.go
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","resolved_command_path":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"/Users/bytedance/golang_pro/nano_cc/go-agent/architecture_test.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config_test.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager_test.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/transport_test.go\n/Users/bytedance/golang_pro/nano_cc/go-agent"}
```

## 2026-06-13T09:15:51Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_pWT0Vkd6d6vlLOsdqBhB4853

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/runtime.go"}
```

### Result

```text
package runtime

import (
	"context"
	"time"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/runtime/compression"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/config"
	"nano_cc/internal/llm"
	"nano_cc/internal/sessions"
)

type EventWriter interface {
	Event(name string, data any) error
}

type conversationStore interface {
	UpdateConversationTitle(ctx context.Context, conversationID, title string) error
	TouchConversationActivity(ctx context.Context, conversationID string) error
	SetConversationHistory(ctx context.Context, conversationID string, messages []storage.Message) error
	SetConversationCache(ctx context.Context, conversationID string, messages []storage.Message) error
	GetConversationCache(ctx context.Context, conversationID string) ([]storage.Message, bool, error)
	ListMessagesByConversation(ctx context.Context, conversationID string, limit int) ([]storage.Message, error)
	CreateToolCall(ctx context.Context, tc storage.ToolCall) error
	CreateSubagentMessage(ctx context.Context, message storage.SubagentMessage) error
	CreatePersistedOutput(ctx context.Context, output storage.PersistedOutput) error
	GetPersistedOutputForConversation(ctx context.Context, id, userID, conversationID string) (storage.PersistedOutput, error)
	GetPersistedOutputByMessageHash(ctx context.Context, conversationID, userID, messageID, toolCallID, strategy, contentSHA256 string) (storage.PersistedOutput, error)
	CreateContextSummary(ctx context.Context, summary storage.ContextSummary) error
	GetContextSummaryByHistoryHash(ctx context.Context, conversationID, userID, sourceHistorySHA256 string) (storage.ContextSummary, error)
	ListRelevantMemories(ctx context.Context, userID string) ([]storage.Memory, error)
	ListMemoriesByUserAndType(ctx context.Context, userID, memType string) ([]storage.Memory, error)
	ListProjectFactMemories(ctx context.Context, userID string) ([]storage.Memory, error)
	InsertMemory(ctx context.Context, m storage.Memory) error
	CountMemoriesByUserAndType(ctx context.Context, userID, memType string) (int, error)
	CountProjectFactMemories(ctx context.Context, userID string) (int, error)
	DeleteOldestMemories(ctx context.Context, userID, memType string, n int) error
	ReplaceMemoriesByUserAndType(ctx context.Context, userID, memType string, items []storage.Memory) error
	ReplaceProjectFactMemories(ctx context.Context, userID string, items []storage.Memory) error
	ListConversationMemories(ctx context.Context, conversationID string) ([]storage.ConversationMemory, error)
	ReplaceConversationMemories(ctx context.Context, conversationID, userID string, items []storage.ConversationMemory) error
	GetConversationModelHistory(ctx context.Context, conversationID string) ([]storage.Message, bool, error)
	UpsertConversationModelHistory(ctx context.Context, conversationID, userID string, messages []storage.Message) error
	AcquireConversationLock(ctx context.Context, conversationID, token string, ttl, waitTimeout time.Duration) (bool, error)
	RenewConversationLock(ctx context.Context, conversationID, token string, ttl time.Duration) (bool, error)
	ReleaseConversationLock(ctx context.Context, conversationID, token string) error
}

type Service struct {
	Store             conversationStore
	Cfg               config.AppConfig
	LLM               llm.Client
	Tools             *ToolRegistry
	BuiltinSkills     *sessions.SkillLoader
	BasePrompt        string
	LinkMarkdown      config.LinkMarkdownContext
	Hooks             *HookManager
	ContextCompressor *compression.Compressor
	EnableMemory      bool
	MCP               *mcp.Manager
}

func NewService(store conversationStore, cfg config.AppConfig, client llm.Client) *Service {
	return &Service{Store: store, Cfg: cfg, LLM: client, Tools: NewToolRegistry(cfg), Hooks: NewDefaultHookManager(), EnableMemory: true}
}

func (s *Service) hookManager() *HookManager {
	if s.Hooks == nil {
		s.Hooks = NewDefaultHookManager()
	}
	return s.Hooks
}

func (s *Service) SetBuiltinSkills(loader *sessions.SkillLoader) {
	s.BuiltinSkills = loader
}

func (s *Service) SetBasePrompt(prompt string) {
	s.BasePrompt = prompt
}

func (s *Service) SetLinkMarkdownContext(ctx config.LinkMarkdownContext) {
	s.LinkMarkdown = ctx
}

func (s *Service) SetMCPManager(manager *mcp.Manager) {
	s.MCP = manager
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"context\"\n\t\"time\"\n\n\t\"nano_cc/internal/agent/mcp\"\n\t\"nano_cc/internal/agent/runtime/compression\"\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/config\"\n\t\"nano_cc/internal/llm\"\n\t\"nano_cc/internal/sessions\"\n)\n\ntype EventWriter interface {\n\tEvent(name string, data any) error\n}\n\ntype conversationStore interface {\n\tUpdateConversationTitle(ctx context.Context, conversationID, title string) error\n\tTouchConversationActivity(ctx context.Context, conversationID string) error"}
```

## 2026-06-13T09:15:51Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_llBMNQQzYb0xp0Ki3DMQ6421

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/prompt_builder.go"}
```

### Result

```text
package runtime

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/assistant"
	"nano_cc/internal/sessions"
	agenttools "nano_cc/internal/tools"
)

func (s *Service) buildSystemPrompt(ctx context.Context, conversation storage.Conversation, user storage.User, snapshot *agenttools.SkillSnapshot, history []storage.Message, memoryOn bool) string {
	memorySection := ""
	if memoryOn {
		memorySection = s.selectRelevantMemories(ctx, user, history)
	}
	return s.buildSystemPromptWithMemory(user, snapshot, memorySection)
}

func (s *Service) buildSystemPromptWithMemory(user storage.User, snapshot *agenttools.SkillSnapshot, memorySection string) string {
	toolNames := []string(nil)
	if s.Tools != nil {
		toolDefs := s.Tools.Definitions()
		toolNames = make([]string, 0, len(toolDefs))
		for _, def := range toolDefs {
			if def.Function == nil || strings.TrimSpace(def.Function.Name) == "" {
				continue
			}
			toolNames = append(toolNames, def.Function.Name)
		}
	}
	loader := sessions.NewSkillLoader()
	if snapshot != nil && snapshot.Merged != nil {
		loader = snapshot.Merged
	}
	return assistant.BuildSystemPrompt(assistant.PromptOptions{
		BasePrompt:        s.BasePrompt,
		Surface:           fmt.Sprintf("the local TUI session for user %s", user.Username),
		SkillDescriptions: loader.GetDescriptions(),
		MemorySection:     memorySection,
		WorkingDirectory:  strings.TrimSpace(s.Cfg.WorkspaceRoot),
		LinkMarkdown: assistant.LinkMarkdownContext{
			UserPath:         s.LinkMarkdown.UserPath,
			UserContent:      s.LinkMarkdown.UserContent,
			WorkspacePath:    s.LinkMarkdown.WorkspacePath,
			WorkspaceContent: s.LinkMarkdown.WorkspaceContent,
		},
		ToolNames: toolNames,
	})
}

func (s *Service) buildSkillSnapshot(ctx context.Context, userID string) (*agenttools.SkillSnapshot, error) {
	return agenttools.NewSkillSnapshot(nil, s.BuiltinSkills), nil
}

func buildOpenAIMessages(systemPrompt string, history []storage.Message) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{{Role: "system", Content: systemPrompt}}
	for _, msg := range history {
		messages = append(messages, openai.ChatCompletionMessage{Role: msg.Role, Content: msg.Content, ReasoningContent: msg.ReasoningContent, ToolCallID: msg.ToolCallID, ToolCalls: storageToolCallsToOpenAI(msg.ToolCalls)})
	}
	return messages
}

func storageToolCallsToOpenAI(toolCalls []storage.MessageToolCall) []openai.ToolCall {
	if len(toolCalls) == 0 {
		return nil
	}
	converted := make([]openai.ToolCall, 0, len(toolCalls))
	for _, toolCall := range toolCalls {
		converted = append(converted, openai.ToolCall{ID: toolCall.ID, Type: openai.ToolType(toolCall.Type), Function: openai.FunctionCall{Name: toolCall.Function.Name, Arguments: toolCall.Function.Arguments}})
	}
	return converted
}

func openAIToolCallsToStorage(toolCalls []openai.ToolCall) []storage.MessageToolCall {
	if len(toolCalls) == 0 {
		return nil
	}
	converted := make([]storage.MessageToolCall, 0, len(toolCalls))
	for _, toolCall := range toolCalls {
		converted = append(converted, storage.MessageToolCall{ID: toolCall.ID, Type: string(toolCall.Type), Function: storage.MessageFunctionCall{Name: toolCall.Function.Name, Arguments: toolCall.Function.Arguments}})
	}
	return converted
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"context\"\n\t\"fmt\"\n\t\"strings\"\n\n\topenai \"github.com/sashabaranov/go-openai\"\n\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/assistant\"\n\t\"nano_cc/internal/sessions\"\n\tagenttools \"nano_cc/internal/tools\"\n)\n\nfunc (s *Service) buildSystemPrompt(ctx context.Context, conversation storage.Conversation, user storage.User, snapshot *agenttools.SkillSnapshot, history []storage.Message, memoryOn bool) string {\n\tmemorySection := \"\"\n\tif memoryOn {\n\t\tmemorySection = s.selectReleva"}
```

## 2026-06-13T09:15:51Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_MUF70fnJ1zpjsL5Y3Dvq6251

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/memory.go"}
```

### Result

```text
package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/idgen"
	"nano_cc/internal/logger"
)

const (
	MemoryTypeEpisodicMemory = "episodic_memory"
	MemoryTypeUserPreference = "user_preference"
	MemoryTypeProjectFact    = "project_fact"
)

const (
	maxPreferenceMemories    = 20
	maxSessionSummaries      = 20
	sessionSummaryPruneCount = 10
	maxProjectFactMemories   = 50
	maxInjectedMemories      = 10
	maxMemoryDialogueChars   = 4000
	maxMemoryNameRunes       = 80
	maxMemoryDescRunes       = 300
	maxMemoryBodyRunes       = 2000
)

func validMemoryType(t string) bool {
	switch t {
	case MemoryTypeEpisodicMemory, MemoryTypeUserPreference, MemoryTypeProjectFact:
		return true
	}
	return false
}

const memoryExtractionSystemPrompt = `You are a project-scoped long-term memory extraction engine for a personal assistant named Link.
All memories are valid ONLY for the current project. Do not create memories that should be reused in other projects.
From the dialogue, extract durable memories worth keeping for the current project. Use the "type" field for three kinds:
- "episodic_memory": a concrete event/experience in this project session. Preserve factual integrity and temporal order; summarize what happened, not raw messages.
- "user_preference": a stable user preference, constraint, or recurring habit that applies to this project.
- "project_fact": a reusable fact about the current project, such as architecture, commands, conventions, dependencies, known constraints, or implementation decisions. It is not general world knowledge and must not be treated as valid outside this project.

Rules:
- Only extract NEW information not already covered by "Existing memories".
- Do not store one-off, trivial, or sensitive private data (passwords, payment, health).
- Do not extract facts about other projects.
- "name": short title (<=80 chars). "description": one-sentence gist (<=300 chars). "body": supporting detail (<=2000 chars).
- Output ONLY a JSON array: [{"name","type","description","body"}].
- If nothing new or everything is already covered, output exactly [].`

const memorySelectionSystemPrompt = `You are a project-scoped memory retrieval engine. Given the current project conversation context and a numbered list of candidate memories from this project's memory index, select the ones RELEVANT and USEFUL for answering the user right now.
- Candidate memories are valid only for the current project.
- Select at most 10.
- Prefer specific, on-topic memories; ignore unrelated ones.
- Output ONLY a JSON array of the selected memory indices, e.g. [0,3,7]. If none, output [].`

func memoryConsolidationSystemPrompt(typeLabel, typeValue string) string {
	return fmt.Sprintf(`You are a project-scoped memory consolidation engine. You are given the FULL current list of "%s" memories for the current project.
Rewrite them into a clean, minimal set:
- Merge duplicates and near-duplicates into a single entry.
- Reconcile contradictions, keeping the most recent / most reliable statement.
- Drop outdated or superseded memories.
- Keep only facts, preferences, and events that are valid for the current project; do not create cross-project memories.
- Never invent new facts; only reorganize what is given.
- Keep limits: name <=80, description <=300, body <=2000 chars.
Output ONLY a JSON array [{"name","type","description","body"}] representing the COMPLETE refined list (it fully replaces the old list). All entries must have type "%s".`, typeLabel, typeValue)
}

type extractedMemory struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Body        string `json:"body"`
}

// extractMemories runs once at the end of a conversation turn: it asks the LLM
// to extract the three kinds of memories from the dialogue, persists them, and
// triggers consolidation/pruning by type. It is best-effort: failures are
// logged and swallowed so the user-facing response is never affected.
func (s *Service) extractMemories(ctx context.Context, user storage.User, history []storage.Message) {
	if s.LLM == nil {
		return
	}
	dialogue := renderDialogueForMemory(history)
	if strings.TrimSpace(dialogue) == "" {
		return
	}
	existing, err := s.Store.ListRelevantMemories(ctx, user.ID)
	if err != nil {
		logger.Warn(fmt.Sprintf("memory: load existing memories failed: %v", err))
	}
	resp, err := s.LLM.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.Cfg.LLM.ModelID,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: memoryExtractionSystemPrompt},
			{Role: "user", Content: buildExtractionUserPrompt(existing, dialogue)},
		},
	})
	if err != nil {
		logger.Warn(fmt.Sprintf("memory: extraction failed: %v", err))
		return
	}
	if len(resp.Choices) == 0 {
		return
	}
	items := parseExtractedMemories(resp.Choices[0].Message.Content)
	if len(items) == 0 {
		return
	}

	var touchedPref, touchedSession, touchedProjectFact bool
	for _, it := range items {
		uid := user.ID
		if err := s.Store.InsertMemory(ctx, storage.Memory{
			ID:          idgen.New("mem"),
			UserID:      uid,
			Type:        it.Type,
			Name:        it.Name,
			Description: it.Description,
			Body:        it.Body,
		}); err != nil {
			logger.Warn(fmt.Sprintf("memory: insert failed: %v", err))
			continue
		}
		switch it.Type {
		case MemoryTypeUserPreference:
			touchedPref = true
		case MemoryTypeEpisodicMemory:
			touchedSession = true
		case MemoryTypeProjectFact:
			touchedProjectFact = true
		}
	}
	if touchedPref {
		s.maybeConsolidateUserPreferences(ctx, user.ID)
	}
	if touchedSession {
		s.maybePruneSessionSummaries(ctx, user.ID)
	}
	if touchedProjectFact {
		s.maybeConsolidateProjectFacts(ctx, user.ID)
	}
}

func (s *Service) maybeConsolidateUserPreferences(ctx context.Context, userID string) {
	items, err := s.Store.ListMemoriesByUserAndType(ctx, userID, MemoryTypeUserPreference)
	if err != nil || len(items) < maxPreferenceMemories {
		return
	}
	refined := s.consolidateViaLLM(ctx, "user_preference", MemoryTypeUserPreference, items)
	if refined == nil {
		return
	}
	if err := s.Store.ReplaceMemoriesByUserAndType(ctx, userID, MemoryTypeUserPreference, refined); err != nil {
		logger.Warn(fmt.Sprintf("memory: replace user preferences failed: %v", err))
	}
}

func (s *Service) maybeConsolidateProjectFacts(ctx context.Context, userID string) {
	items, err := s.Store.ListProjectFactMemories(ctx, userID)
	if err != nil || len(items) < maxProjectFactMemories {
		return
	}
	refined := s.consolidateViaLLM(ctx, "project_fact", MemoryTypeProjectFact, items)
	if refined == nil {
		return
	}
	if err := s.Store.ReplaceProjectFactMemories(ctx, userID, refined); err != nil {
		logger.Warn(fmt.Sprintf("memory: replace project fact memories failed: %v", err))
	}
}

func (s *Service) maybePruneSessionSummaries(ctx context.Context, userID string) {
	n, err := s.Store.CountMemoriesByUserAndType(ctx, userID, MemoryTypeEpisodicMemory)
	if err != nil || n < maxSessionSummaries {
		return
	}
	if err := s.Store.DeleteOldestMemories(ctx, userID, MemoryTypeEpisodicMemory, sessionSummaryPruneCount); err != nil {
		logger.Warn(fmt.Sprintf("memory: prune session summaries failed: %v", err))
	}
}

// consolidateViaLLM feeds the full memory list to the model and parses the
// refined complete list. Returns nil on failure so the caller leaves data
// untouched.
func (s *Service) consolidateViaLLM(ctx context.Context, typeLabel, typeValue string, items []storage.Memory) []storage.Memory {
	if s.LLM == nil {
		return nil
	}
	resp, err := s.LLM.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.Cfg.LLM.ModelID,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: memoryConsolidationSystemPrompt(typeLabel, typeValue)},
			{Role: "user", Content: buildConsolidationUserPrompt(items)},
		},
	})
	if err != nil {
		logger.Warn(fmt.Sprintf("memory: consolidation failed: %v", err))
		return nil
	}
	if len(resp.Choices) == 0 {
		return nil
	}
	parsed := parseExtractedMemories(resp.Choices[0].Message.Content)
	refined := make([]storage.Memory, 0, len(parsed))
	for _, it := range parsed {
		if it.Type != typeValue {
			continue
		}
		refined = append(refined, storage.Memory{
			ID:          idgen.New("mem"),
			Type:        typeValue,
			Name:        it.Name,
			Description: it.Description,
			Body:        it.Body,
		})
	}
	return refined
}

// selectRelevantMemories runs once before the conversation loop: it asks the
// LLM to pick the memories most relevant to the current context and renders
// them into the MemorySection. Best-effort: returns "" on failure or no data.
func (s *Service) selectRelevantMemories(ctx context.Context, user storage.User, history []storage.Message) string {
	if !s.EnableMemory || s.LLM == nil {
		return ""
	}
	all, err := s.Store.ListRelevantMemories(ctx, user.ID)
	if err != nil {
		logger.Warn(fmt.Sprintf("memory: load relevant memories failed: %v", err))
		return ""
	}
	if len(all) == 0 {
		return ""
	}
	resp, err := s.LLM.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.Cfg.LLM.ModelID,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: memorySelectionSystemPrompt},
			{Role: "user", Content: buildSelectionUserPrompt(all, renderDialogueForMemory(history))},
		},
	})
	if err != nil {
		logger.Warn(fmt.Sprintf("memory: selection failed: %v", err))
		return ""
	}
	if len(resp.Choices) == 0 {
		return ""
	}
	indices := parseSelectedIDs(resp.Choices[0].Message.Content)
	selected := pickMemoriesByIndex(all, indices, maxInjectedMemories)
	return renderMemorySection(selected)
}

// parseExtractedMemories extracts a JSON array of memory objects from model
// output, tolerating surrounding code fences or prose, then validates types
// and truncates fields.
func parseExtractedMemories(raw string) []extractedMemory {
	trimmed := extractJSONArray(raw)
	if trimmed == "" {
		return nil
	}
	var parsed []extractedMemory
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return nil
	}
	result := make([]extractedMemory, 0, len(parsed))
	for _, it := range parsed {
		it.Type = strings.TrimSpace(it.Type)
		it.Name = strings.TrimSpace(it.Name)
		it.Description = strings.TrimSpace(it.Description)
		it.Body = strings.TrimSpace(it.Body)
		if !validMemoryType(it.Type) || it.Name == "" || it.Body == "" {
			continue
		}
		it.Name = truncateRunes(it.Name, maxMemoryNameRunes)
		it.Description = truncateRunes(it.Description, maxMemoryDescRunes)
		it.Body = truncateRunes(it.Body, maxMemoryBodyRunes)
		result = append(result, it)
	}
	return result
}

// parseSelectedIDs parses a JSON array of integer indices from model output.
func parseSelectedIDs(raw string) []int {
	trimmed := extractJSONArray(raw)
	if trimmed == "" {
		return nil
	}
	var parsed []int
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return nil
	}
	return parsed
}

// pickMemoriesByIndex resolves indices into memories, dropping out-of-range and
// duplicate indices, capped at max.
func pickMemoriesByIndex(all []storage.Memory, indices []int, max int) []storage.Memory {
	seen := make(map[int]struct{})
	result := make([]storage.Memory, 0, len(indices))
	for _, idx := range indices {
		if idx < 0 || idx >= len(all) {
			continue
		}
		if _, exists := seen[idx]; exists {
			continue
		}
		seen[idx] = struct{}{}
		result = append(result, all[idx])
		if len(result) >= max {
			break
		}
	}
	return result
}

// renderMemorySection groups selected project-scoped memories by type and
// renders a Markdown block. Returns "" when empty.
func renderMemorySection(memories []storage.Memory) string {
	if len(memories) == 0 {
		return ""
	}
	var userLines, projectFactLines []string
	for _, m := range memories {
		switch m.Type {
		case MemoryTypeUserPreference:
			userLines = append(userLines, memoryBlock("(喜好) ", m))
		case MemoryTypeEpisodicMemory:
			userLines = append(userLines, memoryBlock("(经历) ", m))
		case MemoryTypeProjectFact:
			projectFactLines = append(projectFactLines, memoryBlock("", m))
		}
	}
	sections := make([]string, 0, 3)
	sections = append(sections, "### 当前项目记忆\n以下记忆仅适用于当前项目；不要迁移到其他项目会话。")
	if len(userLines) > 0 {
		sections = append(sections, "#### 关于用户的长期记忆\n"+strings.Join(userLines, "\n"))
	}
	if len(projectFactLines) > 0 {
		sections = append(sections, "#### 当前项目事实\n"+strings.Join(projectFactLines, "\n"))
	}
	return strings.Join(sections, "\n\n")
}

func memoryBlock(prefix string, m storage.Memory) string {
	line := "- " + prefix + memoryLine(m)
	body := strings.TrimSpace(m.Body)
	if body == "" {
		return line
	}
	return line + "\n  " + strings.ReplaceAll(body, "\n", "\n  ")
}

func memoryLine(m storage.Memory) string {
	if strings.TrimSpace(m.Description) == "" {
		return m.Name
	}
	return m.Name + "：" + m.Description
}

func buildExtractionUserPrompt(existing []storage.Memory, dialogue string) string {
	return "Existing memories:\n" + renderMemoryListForPrompt(existing) + "\n\nDialogue:\n" + dialogue
}

func buildConsolidationUserPrompt(items []storage.Memory) string {
	var b strings.Builder
	for i, m := range items {
		b.WriteString(fmt.Sprintf("[%d] (%s) %s: %s\n", i, m.Type, m.Name, m.Description))
		if strings.TrimSpace(m.Body) != "" {
			b.WriteString("    body: " + m.Body + "\n")
		}
	}
	return b.String()
}

func buildSelectionUserPrompt(all []storage.Memory, dialogue string) string {
	var b strings.Builder
	b.WriteString("Recent conversation:\n")
	b.WriteString(dialogue)
	b.WriteString("\n\nCandidate memories:\n")
	for i, m := range all {
		b.WriteString(fmt.Sprintf("[%d] (%s) %s: %s\n", i, m.Type, m.Name, m.Description))
	}
	return b.String()
}

func renderMemoryListForPrompt(items []storage.Memory) string {
	if len(items) == 0 {
		return "(none)"
	}
	lines := make([]string, 0, len(items))
	for _, m := range items {
		lines = append(lines, fmt.Sprintf("- [%s] %s: %s", m.Type, m.Name, m.Description))
	}
	return strings.Join(lines, "\n")
}

// renderDialogueForMemory builds a plain-text transcript of user and assistant
// messages only, dropping tool noise to keep the LLM calls cheap.
func renderDialogueForMemory(history []storage.Message) string {
	var b strings.Builder
	for _, msg := range history {
		var content string
		switch msg.Role {
		case "user":
			content = strings.TrimSpace(msg.Content)
			if content == "" {
				continue
			}
			b.WriteString("[user] ")
		case "assistant":
			content = strings.TrimSpace(msg.Content)
			if content == "" {
				continue
			}
			b.WriteString("[assistant] ")
		default:
			continue
		}
		b.WriteString(content)
		b.WriteString("\n\n")
		if b.Len() > maxMemoryDialogueChars {
			break
		}
	}
	return b.String()
}

func extractJSONArray(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if start := strings.IndexByte(trimmed, '['); start >= 0 {
		if end := strings.LastIndexByte(trimmed, ']'); end > start {
			return trimmed[start : end+1]
		}
	}
	return ""
}

func truncateRunes(s string, max int) string {
	if runes := []rune(s); len(runes) > max {
		return string(runes[:max])
	}
	return s
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"context\"\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"strings\"\n\n\topenai \"github.com/sashabaranov/go-openai\"\n\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/idgen\"\n\t\"nano_cc/internal/logger\"\n)\n\nconst (\n\tMemoryTypeEpisodicMemory = \"episodic_memory\"\n\tMemoryTypeUserPreference = \"user_preference\"\n\tMemoryTypeProjectFact    = \"project_fact\"\n)\n\nconst (\n\tmaxPreferenceMemories    = 20\n\tmaxSessionSummaries      = 20\n\tsessionSummaryPruneCount = 10\n\tmaxProjectFactMemories   = 50\n\tmaxInjectedMe"}
```

## 2026-06-13T09:15:51Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_03_cMc1zK1Kj7JPCHn8oYR48077

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/models.go"}
```

### Result

```text
package storage

import "time"

type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	MemoryEnabled bool      `json:"memory_enabled"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// MCPServer 是用户配置的一个 MCP 服务器连接。Transport 取值 stdio/sse/streamable。
// Args/Env/Headers 从本地 .link/.mcp.json 读取，在 Go 侧以强类型表示。
type MCPServer struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Name      string            `json:"name"`
	Transport string            `json:"transport"`
	Command   string            `json:"command"`
	Args      []string          `json:"args"`
	Env       map[string]string `json:"env"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Enabled   bool              `json:"enabled"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type Conversation struct {
	ID            string    `json:"id"`
	SessionID     string    `json:"session_id,omitempty"`
	UserID        string    `json:"user_id"`
	RootMessageID string    `json:"root_message_id"`
	Title         string    `json:"title"`
	HistoryJSON   string    `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ResumableSession struct {
	SessionID      string    `json:"session_id"`
	ConversationID string    `json:"conversation_id"`
	WorkspaceRoot  string    `json:"workspace_root"`
	Title          string    `json:"title"`
	MessageCount   int       `json:"message_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// 该表已经不使用，作为存储的结构体，用于存储消息历史
type Message struct {
	ID               string            `json:"id"`
	ConversationID   string            `json:"conversation_id"`
	UserID           string            `json:"user_id"`
	Role             string            `json:"role"`
	Content          string            `json:"content"`
	ReasoningContent string            `json:"reasoning_content,omitempty"`
	ToolCallID       string            `json:"tool_call_id,omitempty"`
	ToolCalls        []MessageToolCall `json:"tool_calls,omitempty"`
	Meta             *MessageMeta      `json:"meta,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
}

// MessageMeta 记录助手回复的元信息，仅对最终 assistant 消息填充。
type MessageMeta struct {
	ToolCallCount int `json:"tool_call_count"`          // 调用工具次数（0 也需序列化，保证历史展示一致）
	ContextTokens int `json:"context_tokens,omitempty"` // 当前上下文估算 token
	ContextBudget int `json:"context_budget,omitempty"` // 上下文预算（用于算占比）
}

type MessageToolCall struct {
	ID       string              `json:"id"`
	Type     string              `json:"type"`
	Function MessageFunctionCall `json:"function"`
}

type MessageFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolCall struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	ToolName       string    `json:"tool_name"`
	Status         string    `json:"status"`
	Summary        string    `json:"summary"`
	CreatedAt      time.Time `json:"created_at"`
}

type PersistedOutput struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	MessageID      string    `json:"message_id"`
	ToolCallID     string    `json:"tool_call_id"`
	Kind           string    `json:"kind"`
	Strategy       string    `json:"strategy"`
	OriginalBytes  int       `json:"original_bytes"`
	ContentSHA256  string    `json:"content_sha256"`
	Content        string    `json:"content"`
	Preview        string    `json:"preview"`
	CreatedAt      time.Time `json:"created_at"`
}

type ToolResultLogEntry struct {
	ConversationID string    `json:"conversation_id"`
	SessionID      string    `json:"session_id,omitempty"`
	UserID         string    `json:"user_id"`
	ToolCallID     string    `json:"tool_call_id"`
	ToolName       string    `json:"tool_name"`
	RawArgs        string    `json:"raw_args"`
	Status         string    `json:"status"`
	Result         string    `json:"result"`
	AuditSummary   string    `json:"audit_summary,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type ContextSummary struct {
	ID                    string    `json:"id"`
	ConversationID        string    `json:"conversation_id"`
	UserID                string    `json:"user_id"`
	SourceHistorySHA256   string    `json:"source_history_sha256"`
	Strategy              string    `json:"strategy"`
	EstimatedTokensBefore int       `json:"estimated_tokens_before"`
	EstimatedTokensAfter  int       `json:"estimated_tokens_after"`
	Summary               string    `json:"summary"`
	CreatedAt             time.Time `json:"created_at"`
}

type SubagentMessage struct {
	ID               string            `json:"id"`
	RunID            string            `json:"run_id"`
	ParentToolCallID string            `json:"parent_tool_call_id"`
	ConversationID   string            `json:"conversation_id"`
	UserID           string            `json:"user_id"`
	SequenceNo       int               `json:"sequence_no"`
	Role             string            `json:"role"`
	Content          string            `json:"content"`
	ReasoningContent string            `json:"reasoning_content,omitempty"`
	ToolCallID       string            `json:"tool_call_id,omitempty"`
	ToolCalls        []MessageToolCall `json:"tool_calls,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
}

type Memory struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ConversationMemory 是单个会话维度、随每轮对话增量维护的“当前会话主干信息”
// 条目。它不注入 system prompt，仅在上下文压缩触发全量摘要时作为替代品使用。
type ConversationMemory struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package storage\n\nimport \"time\"\n\ntype User struct {\n\tID            string    `json:\"id\"`\n\tEmail         string    `json:\"email\"`\n\tUsername      string    `json:\"username\"`\n\tMemoryEnabled bool      `json:\"memory_enabled\"`\n\tCreatedAt     time.Time `json:\"created_at\"`\n\tUpdatedAt     time.Time `json:\"updated_at\"`\n}\n\n// MCPServer 是用户配置的一个 MCP 服务器连接。Transport 取值 stdio/sse/streamable。\n// Args/Env/Headers 从本地 .link/.mcp.json 读取，在 Go 侧以强类型表示。"}
```

## 2026-06-13T09:15:51Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_04_MPk8Zb50fkPoOJnmkXCM8904

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/conversation_memory.go"}
```

### Result

```text
package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/idgen"
	"nano_cc/internal/logger"
)

const conversationMemorySystemPrompt = `你是“当前会话记忆”维护引擎。给你“已有会话记忆条目”和“最新一轮对话”，
请输出更新后的【完整】会话记忆条目列表，使其准确反映本场会话至今的主干信息：
- 覆盖：当前用户目标、关键决策与结论、已完成/已产出的内容、重要约束与上下文、待办与下一步。
- 合并重复、用新信息更新旧条目、删除已过时或被推翻的内容；只重组已知信息，不编造。
- name：短标题(<=80字，可用 [目标]/[决策]/[产出]/[待办] 等前缀)。description：一句话要点(<=300字)。body：支撑细节(<=2000字)。
- 仅输出 JSON 数组：[{"name","description","body"}]。无可记录时输出 []。`

type extractedConversationMemory struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Body        string `json:"body"`
}

// updateConversationMemory runs once at the end of a conversation turn: it asks
// the LLM to rewrite the full conversation memory list based on existing entries
// plus the latest dialogue, then replaces the stored entries for this
// conversation. It is best-effort: failures are logged and swallowed so the
// user-facing response is never affected. The result is never injected into the
// system prompt; it is only consumed by the context compression pipeline.
func (s *Service) updateConversationMemory(ctx context.Context, conversation storage.Conversation, user storage.User, history []storage.Message) {
	if s.LLM == nil {
		return
	}
	dialogue := renderDialogueForMemory(history)
	if strings.TrimSpace(dialogue) == "" {
		return
	}
	existing, err := s.Store.ListConversationMemories(ctx, conversation.ID)
	if err != nil {
		logger.Warn(fmt.Sprintf("conversation memory: load existing failed: %v", err))
	}
	resp, err := s.LLM.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.Cfg.LLM.ModelID,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: conversationMemorySystemPrompt},
			{Role: "user", Content: buildConversationMemoryUserPrompt(existing, dialogue)},
		},
	})
	if err != nil {
		logger.Warn(fmt.Sprintf("conversation memory: update failed: %v", err))
		return
	}
	if len(resp.Choices) == 0 {
		return
	}
	items := parseConversationMemories(resp.Choices[0].Message.Content)
	refined := make([]storage.ConversationMemory, 0, len(items))
	for _, it := range items {
		refined = append(refined, storage.ConversationMemory{
			ID:          idgen.New("cm"),
			Name:        it.Name,
			Description: it.Description,
			Body:        it.Body,
		})
	}
	if err := s.Store.ReplaceConversationMemories(ctx, conversation.ID, user.ID, refined); err != nil {
		logger.Warn(fmt.Sprintf("conversation memory: replace failed: %v", err))
	}
}

// parseConversationMemories extracts a JSON array of memory objects from model
// output, tolerating surrounding code fences or prose, then truncates fields.
func parseConversationMemories(raw string) []extractedConversationMemory {
	trimmed := extractJSONArray(raw)
	if trimmed == "" {
		return nil
	}
	var parsed []extractedConversationMemory
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return nil
	}
	result := make([]extractedConversationMemory, 0, len(parsed))
	for _, it := range parsed {
		it.Name = strings.TrimSpace(it.Name)
		it.Description = strings.TrimSpace(it.Description)
		it.Body = strings.TrimSpace(it.Body)
		if it.Name == "" || it.Body == "" {
			continue
		}
		it.Name = truncateRunes(it.Name, maxMemoryNameRunes)
		it.Description = truncateRunes(it.Description, maxMemoryDescRunes)
		it.Body = truncateRunes(it.Body, maxMemoryBodyRunes)
		result = append(result, it)
	}
	return result
}

func buildConversationMemoryUserPrompt(existing []storage.ConversationMemory, dialogue string) string {
	var b strings.Builder
	b.WriteString("Existing conversation memory entries:\n")
	if len(existing) == 0 {
		b.WriteString("(none)")
	} else {
		for i, m := range existing {
			b.WriteString(fmt.Sprintf("[%d] %s: %s\n", i, m.Name, m.Description))
			if strings.TrimSpace(m.Body) != "" {
				b.WriteString("    body: " + m.Body + "\n")
			}
		}
	}
	b.WriteString("\n\nLatest dialogue:\n")
	b.WriteString(dialogue)
	return b.String()
}

// scheduleMemoryWork 在一轮对话结束后异步执行收尾操作（模型历史持久化 + 记忆提取
// + 会话记忆更新）。
// 它接管入口持有的会话锁（token）：在独立的 background context 中执行，期间持续
// 续期，完成后停止续期并释放锁。返回 true 表示已接管锁所有权（调用方应跳过 defer
// 释放）；返回 false 表示未持锁（已降级），调用方按原逻辑处理。
// memoryOn 仅控制记忆提取与会话记忆更新；锁释放与模型历史持久化始终执行。
func (s *Service) scheduleMemoryWork(conv storage.Conversation, user storage.User, history []storage.Message, modelHistory []storage.Message, token string, stopRenew func(), memoryOn bool) bool {
	if token == "" {
		// 入口未持锁（已降级）→ 跳过收尾，不接管锁。
		return false
	}
	// 停止请求期看门狗，收尾 goroutine 内重新启动一个。
	if stopRenew != nil {
		stopRenew()
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Warn(fmt.Sprintf("memory work: panic recovered conversation=%s: %v", conv.ID, r))
			}
		}()
		defer s.Store.ReleaseConversationLock(context.Background(), conv.ID, token)
		stop := s.startLockRenewer(conv.ID, token)
		defer stop()

		ctx, cancel := context.WithTimeout(context.Background(), s.Cfg.MemoryWorkTimeout)
		defer cancel()
		// 先落库模型历史：即使后续记忆相关的 LLM 调用超时，也不丢失本轮压缩成果。
		if len(modelHistory) > 0 {
			if err := s.Store.UpsertConversationModelHistory(ctx, conv.ID, user.ID, modelHistory); err != nil {
				logger.Warn(fmt.Sprintf("model history: persist failed conversation=%s: %v", conv.ID, err))
			}
		}
		if memoryOn {
			s.extractMemories(ctx, user, history)
			s.updateConversationMemory(ctx, conv, user, history)
		}
	}()
	return true
}

// startLockRenewer 启动一个后台看门狗，按 TTL/3 周期为会话锁续期，返回的 stop
// 函数用于停止续期（幂等）。续期失败（锁已不属于当前 token）时记录告警并停止。
func (s *Service) startLockRenewer(conversationID, token string) func() {
	ttl := s.Cfg.ConversationLockTTL
	interval := ttl / 3
	if interval <= 0 {
		interval = ttl
	}
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				ok, err := s.Store.RenewConversationLock(context.Background(), conversationID, token, ttl)
				if err != nil {
					logger.Warn(fmt.Sprintf("conversation lock: renew failed conversation=%s: %v", conversationID, err))
					return
				}
				if !ok {
					logger.Warn(fmt.Sprintf("conversation lock: renew lost ownership conversation=%s", conversationID))
					return
				}
			}
		}
	}()
	var once sync.Once
	return func() {
		once.Do(func() { close(done) })
	}
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"context\"\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"strings\"\n\t\"sync\"\n\t\"time\"\n\n\topenai \"github.com/sashabaranov/go-openai\"\n\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/idgen\"\n\t\"nano_cc/internal/logger\"\n)\n\nconst conversationMemorySystemPrompt = `你是“当前会话记忆”维护引擎。给你“已有会话记忆条目”和“最新一轮对话”，\n请输出更新后的【完整】会话记忆条目列表，使其准确反映本场会话至今的主干信息：\n- 覆盖\ufffd"}
```

## 2026-06-13T09:15:55Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_MPk3r3t4hC4vNDQgWP4D0747

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/runtime.go"}
```

### Result

```text
package runtime

import (
	"context"
	"time"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/runtime/compression"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/config"
	"nano_cc/internal/llm"
	"nano_cc/internal/sessions"
)

type EventWriter interface {
	Event(name string, data any) error
}

type conversationStore interface {
	UpdateConversationTitle(ctx context.Context, conversationID, title string) error
	TouchConversationActivity(ctx context.Context, conversationID string) error
	SetConversationHistory(ctx context.Context, conversationID string, messages []storage.Message) error
	SetConversationCache(ctx context.Context, conversationID string, messages []storage.Message) error
	GetConversationCache(ctx context.Context, conversationID string) ([]storage.Message, bool, error)
	ListMessagesByConversation(ctx context.Context, conversationID string, limit int) ([]storage.Message, error)
	CreateToolCall(ctx context.Context, tc storage.ToolCall) error
	CreateSubagentMessage(ctx context.Context, message storage.SubagentMessage) error
	CreatePersistedOutput(ctx context.Context, output storage.PersistedOutput) error
	GetPersistedOutputForConversation(ctx context.Context, id, userID, conversationID string) (storage.PersistedOutput, error)
	GetPersistedOutputByMessageHash(ctx context.Context, conversationID, userID, messageID, toolCallID, strategy, contentSHA256 string) (storage.PersistedOutput, error)
	CreateContextSummary(ctx context.Context, summary storage.ContextSummary) error
	GetContextSummaryByHistoryHash(ctx context.Context, conversationID, userID, sourceHistorySHA256 string) (storage.ContextSummary, error)
	ListRelevantMemories(ctx context.Context, userID string) ([]storage.Memory, error)
	ListMemoriesByUserAndType(ctx context.Context, userID, memType string) ([]storage.Memory, error)
	ListProjectFactMemories(ctx context.Context, userID string) ([]storage.Memory, error)
	InsertMemory(ctx context.Context, m storage.Memory) error
	CountMemoriesByUserAndType(ctx context.Context, userID, memType string) (int, error)
	CountProjectFactMemories(ctx context.Context, userID string) (int, error)
	DeleteOldestMemories(ctx context.Context, userID, memType string, n int) error
	ReplaceMemoriesByUserAndType(ctx context.Context, userID, memType string, items []storage.Memory) error
	ReplaceProjectFactMemories(ctx context.Context, userID string, items []storage.Memory) error
	ListConversationMemories(ctx context.Context, conversationID string) ([]storage.ConversationMemory, error)
	ReplaceConversationMemories(ctx context.Context, conversationID, userID string, items []storage.ConversationMemory) error
	GetConversationModelHistory(ctx context.Context, conversationID string) ([]storage.Message, bool, error)
	UpsertConversationModelHistory(ctx context.Context, conversationID, userID string, messages []storage.Message) error
	AcquireConversationLock(ctx context.Context, conversationID, token string, ttl, waitTimeout time.Duration) (bool, error)
	RenewConversationLock(ctx context.Context, conversationID, token string, ttl time.Duration) (bool, error)
	ReleaseConversationLock(ctx context.Context, conversationID, token string) error
}

type Service struct {
	Store             conversationStore
	Cfg               config.AppConfig
	LLM               llm.Client
	Tools             *ToolRegistry
	BuiltinSkills     *sessions.SkillLoader
	BasePrompt        string
	LinkMarkdown      config.LinkMarkdownContext
	Hooks             *HookManager
	ContextCompressor *compression.Compressor
	EnableMemory      bool
	MCP               *mcp.Manager
}

func NewService(store conversationStore, cfg config.AppConfig, client llm.Client) *Service {
	return &Service{Store: store, Cfg: cfg, LLM: client, Tools: NewToolRegistry(cfg), Hooks: NewDefaultHookManager(), EnableMemory: true}
}

func (s *Service) hookManager() *HookManager {
	if s.Hooks == nil {
		s.Hooks = NewDefaultHookManager()
	}
	return s.Hooks
}

func (s *Service) SetBuiltinSkills(loader *sessions.SkillLoader) {
	s.BuiltinSkills = loader
}

func (s *Service) SetBasePrompt(prompt string) {
	s.BasePrompt = prompt
}

func (s *Service) SetLinkMarkdownContext(ctx config.LinkMarkdownContext) {
	s.LinkMarkdown = ctx
}

func (s *Service) SetMCPManager(manager *mcp.Manager) {
	s.MCP = manager
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"context\"\n\t\"time\"\n\n\t\"nano_cc/internal/agent/mcp\"\n\t\"nano_cc/internal/agent/runtime/compression\"\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/config\"\n\t\"nano_cc/internal/llm\"\n\t\"nano_cc/internal/sessions\"\n)\n\ntype EventWriter interface {\n\tEvent(name string, data any) error\n}\n\ntype conversationStore interface {\n\tUpdateConversationTitle(ctx context.Context, conversationID, title string) error\n\tTouchConversationActivity(ctx context.Context, conversationID string) error"}
```

## 2026-06-13T09:15:55Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_jt3TmqWObpqQODujVB7M7030

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/config.go"}
```

### Result

```text
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"-"`
	ModelID string `json:"model_id"`
}

type fileConfig struct {
	Config
	AppHome                            string `json:"app_home"`
	BuiltinSkillsDir                   string `json:"builtin_skills_dir"`
	CommandBinDir                      string `json:"command_bin_dir"`
	CommandScriptDir                   string `json:"command_script_dir"`
	SystemPromptPath                   string `json:"system_prompt_path"`
	WorkspaceRoot                      string `json:"workspace_root"`
	AllowedTools                       string `json:"allowed_tools"`
	BashAllowOutsideWorkspace          bool   `json:"bash_allow_outside_workspace"`
	BashAllowDangerousCommands         bool   `json:"bash_allow_dangerous_commands"`
	ConversationLockTTLSeconds         int    `json:"conversation_lock_ttl_seconds"`
	MemoryWorkTimeoutSeconds           int    `json:"memory_work_timeout_seconds"`
	ConversationLockWaitTimeoutSeconds int    `json:"conversation_lock_wait_timeout_seconds"`
}

type AppConfig struct {
	LLM                         Config
	AppHome                     string
	BuiltinSkillsDir            string
	CommandBinDir               string
	CommandScriptDir            string
	SystemPromptPath            string
	WorkspaceRoot               string
	LogsDir                     string
	AllowedTools                []string
	BashAllowOutsideWorkspace   bool
	BashAllowDangerousCommands  bool
	ConversationLockTTL         time.Duration
	MemoryWorkTimeout           time.Duration
	ConversationLockWaitTimeout time.Duration
}

func loadConfigFile() (fileConfig, error) {
	data, err := os.ReadFile(configFilePath())
	if err != nil {
		return fileConfig{}, err
	}

	var cfg fileConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fileConfig{}, fmt.Errorf("failed to parse workspace config.json: %w", err)
	}
	return cfg, nil
}

func configFilePath() string {
	return "config.json"
}

func loadLLMConfig(fileCfg fileConfig) (Config, error) {
	cfg := Config{
		BaseURL: strings.TrimSpace(fileCfg.Config.BaseURL),
		APIKey:  strings.TrimSpace(getenv("OPENAI_API_KEY")),
		ModelID: strings.TrimSpace(fileCfg.Config.ModelID),
	}

	if cfg.BaseURL == "" || cfg.ModelID == "" {
		return Config{}, fmt.Errorf("missing LLM config; set base_url, model_id in config.json")
	}
	if cfg.APIKey == "" {
		return Config{}, fmt.Errorf("missing LLM api key; set OPENAI_API_KEY environment variable")
	}

	return cfg, nil
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package config\n\nimport (\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"os\"\n\t\"strings\"\n\t\"time\"\n)\n\ntype Config struct {\n\tBaseURL string `json:\"base_url\"`\n\tAPIKey  string `json:\"-\"`\n\tModelID string `json:\"model_id\"`\n}\n\ntype fileConfig struct {\n\tConfig\n\tAppHome                            string `json:\"app_home\"`\n\tBuiltinSkillsDir                   string `json:\"builtin_skills_dir\"`\n\tCommandBinDir                      string `json:\"command_bin_dir\"`\n\tCommandScriptDir                   string `json:\"command_script_dir\"`\n"}
```

## 2026-06-13T09:15:55Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_6ErLZ3xcduSXE93AdJx05868

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/llm/client.go"}
```

### Result

```text
package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// Client is the minimal chat-completion interface the runtime depends on.
type Client interface {
	CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
	CreateChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (ChatCompletionStream, error)
}

// ChatCompletionStream is a streamed chat-completion response.
type ChatCompletionStream interface {
	Recv() (openai.ChatCompletionStreamResponse, error)
	Close() error
}

type deepseekClient struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

// NewDeepseekClient returns a Client backed by the DeepSeek chat-completions API.
func NewDeepseekClient(baseURL, apiKey string) Client {
	return &deepseekClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		http:    &http.Client{},
	}
}

// buildRequestBody assembles the DeepSeek request body, passing through
// max_tokens when set so callers can control/override the output budget.
func buildRequestBody(req openai.ChatCompletionRequest, stream bool) map[string]any {
	body := map[string]any{
		"model":            req.Model,
		"messages":         req.Messages,
		"thinking":         map[string]string{"type": "enabled"},
		"reasoning_effort": "high",
		"stream":           stream,
	}
	if len(req.Tools) > 0 {
		body["tools"] = req.Tools
	}
	if req.MaxTokens > 0 {
		body["max_tokens"] = req.MaxTokens
	}
	return body
}

// retryTransient runs do until it succeeds, returns a non-retryable error, or
// the transient-retry budget is exhausted. 429/529 are retried with
// exponential backoff plus jitter; ctx cancellation aborts immediately.
func retryTransient[T any](ctx context.Context, do func() (T, error)) (T, error) {
	var zero T
	for attempt := 0; ; attempt++ {
		result, err := do()
		if err == nil {
			return result, nil
		}
		var apiErr *APIError
		if !errors.As(err, &apiErr) || !isRetryableStatus(apiErr.StatusCode) || attempt >= maxTransientRetries {
			return zero, err
		}
		select {
		case <-time.After(backoffDelay(attempt)):
		case <-ctx.Done():
			return zero, ctx.Err()
		}
	}
}

func (c *deepseekClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	return retryTransient(ctx, func() (openai.ChatCompletionResponse, error) {
		return c.createChatCompletionOnce(ctx, req)
	})
}

func (c *deepseekClient) createChatCompletionOnce(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	jsonBody, err := json.Marshal(buildRequestBody(req, false))
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("marshal request: %w", err)
	}

	url := c.baseURL + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return openai.ChatCompletionResponse{}, &APIError{StatusCode: resp.StatusCode, Body: string(respBody)}
	}

	var result openai.ChatCompletionResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("unmarshal response: %w", err)
	}

	return result, nil
}

func (c *deepseekClient) CreateChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (ChatCompletionStream, error) {
	return retryTransient(ctx, func() (ChatCompletionStream, error) {
		return c.createChatCompletionStreamOnce(ctx, req)
	})
}

func (c *deepseekClient) createChatCompletionStreamOnce(ctx context.Context, req openai.ChatCompletionRequest) (ChatCompletionStream, error) {
	jsonBody, err := json.Marshal(buildRequestBody(req, true))
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := c.baseURL + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		respBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("read response: %w", readErr)
		}
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(respBody)}
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	return &deepseekChatStream{body: resp.Body, scanner: scanner}, nil
}

type deepseekChatStream struct {
	body    io.Closer
	scanner *bufio.Scanner
}

func (s *deepseekChatStream) Recv() (openai.ChatCompletionStreamResponse, error) {
	for s.scanner.Scan() {
		line := strings.TrimSpace(s.scanner.Text())
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			return openai.ChatCompletionStreamResponse{}, io.EOF
		}
		var chunk openai.ChatCompletionStreamResponse
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			return openai.ChatCompletionStreamResponse{}, fmt.Errorf("unmarshal stream response: %w", err)
		}
		return chunk, nil
	}
	if err := s.scanner.Err(); err != nil {
		return openai.ChatCompletionStreamResponse{}, err
	}
	return openai.ChatCompletionStreamResponse{}, io.EOF
}

func (s *deepseekChatStream) Close() error {
	return s.body.Close()
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package llm\n\nimport (\n\t\"bufio\"\n\t\"bytes\"\n\t\"context\"\n\t\"encoding/json\"\n\t\"errors\"\n\t\"fmt\"\n\t\"io\"\n\t\"net/http\"\n\t\"strings\"\n\t\"time\"\n\n\topenai \"github.com/sashabaranov/go-openai\"\n)\n\n// Client is the minimal chat-completion interface the runtime depends on.\ntype Client interface {\n\tCreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)\n\tCreateChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (ChatCompletionStream, error)\n}\n"}
```

## 2026-06-13T09:15:55Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_03_F18GlEv4Mhoo7FD8sqSh0326

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/app.go"}
```

### Result

```text
package tui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/runtime"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/logger"
	"nano_cc/internal/sessions"
)

type SessionInfo struct {
	User         storage.User
	Conversation storage.Conversation
	CWD          string
	Resumer      SessionResumer
	Skills       []sessions.SkillSummary
	MCPServers   []mcp.ServerStatus
	SkillCount   int
	MCPToolCount int
}

type SessionResumer interface {
	ListResumableSessions(ctx context.Context, workspaceRoot string) ([]storage.ResumableSession, error)
	ResumeSession(ctx context.Context, sessionID, currentWorkspace string, user storage.User) (storage.Conversation, []storage.Message, error)
}

type Message struct {
	Role             string
	Content          string
	ReasoningContent string
}

type palette struct {
	ink      lipgloss.Color
	muted    lipgloss.Color
	panel    lipgloss.Color
	cyan     lipgloss.Color
	mint     lipgloss.Color
	lavender lipgloss.Color
	butter   lipgloss.Color
	coral    lipgloss.Color
}

var tuiPalette = palette{
	ink:      lipgloss.Color("255"),
	muted:    lipgloss.Color("244"),
	panel:    lipgloss.Color("238"),
	cyan:     lipgloss.Color("81"),
	mint:     lipgloss.Color("120"),
	lavender: lipgloss.Color("183"),
	butter:   lipgloss.Color("229"),
	coral:    lipgloss.Color("210"),
}

const inputCursor = "█"

type Model struct {
	runtime          *runtime.Service
	session          SessionInfo
	messages         []Message
	input            textarea.Model
	viewport         viewport.Model
	width            int
	height           int
	running          bool
	events           chan Event
	cancel           context.CancelFunc
	renderer         *glamour.TermRenderer
	generation       int64
	resumeSelecting  bool
	resumeCandidates []storage.ResumableSession
	toolCallCount    int
	contextTokens    int
	contextBudget    int
}

func NewModel(runtimeService *runtime.Service, session SessionInfo) Model {
	input := textarea.New()
	input.Placeholder = "问 go-agent 一件事..."
	input.Focus()
	input.SetHeight(2)
	input.ShowLineNumbers = false
	vp := viewport.New(100, 20)
	renderer := newMarkdownRenderer(100)
	return Model{runtime: runtimeService, session: session, input: input, viewport: vp, width: 100, height: 20, events: make(chan Event, 128), renderer: renderer}
}

func Run(ctx context.Context, runtimeService *runtime.Service, session SessionInfo) error {
	previousConsole := logger.SetConsoleEnabled(false)
	defer logger.SetConsoleEnabled(previousConsole)
	program := tea.NewProgram(NewModel(runtimeService, session), tea.WithContext(ctx), tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := program.Run()
	return err
}

func (m Model) Init() tea.Cmd { return textarea.Blink }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = m.viewportHeight()
		m.input.SetWidth(max(20, msg.Width-4))
		m.renderer = newMarkdownRenderer(m.messageWidth())
		m.refreshViewport()
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.running && m.cancel != nil {
				m.cancel()
				m.generation++
				m.running = false
				m.appendMessage("system", "已中断当前生成")
				return m, nil
			}
			return m, tea.Quit
		case "enter":
			text := strings.TrimSpace(m.input.Value())
			if text == "" || m.running {
				return m, nil
			}
			m.input.Reset()
			if m.resumeSelecting && m.handleResumeSelection(text) {
				m.refreshViewport()
				return m, nil
			}
			if strings.HasPrefix(text, "/") && m.handleSlashCommand(text) {
				m.refreshViewport()
				return m, nil
			}
			m.appendMessage("user", text)
			m.running = true
			m.toolCallCount = 0
			m.generation++
			generation := m.generation
			ctx, cancel := context.WithCancel(context.Background())
			m.cancel = cancel
			return m, tea.Batch(m.waitEvent(), m.respond(ctx, text, generation))
		}
		if isTerminalProbeResponseInput(msg) {
			return m, nil
		}
		if cmd, ok := m.updateViewportScroll(msg); ok {
			return m, cmd
		}
	case tea.MouseMsg:
		if cmd, ok := m.updateViewportScroll(msg); ok {
			return m, cmd
		}
	case Event:
		if msg.Generation != 0 && msg.Generation != m.generation {
			if m.running {
				return m, m.waitEvent()
			}
			return m, nil
		}
		switch msg.Name {
		case "assistant_delta":
			m.appendAssistantDelta(msg.Content)
		case "reasoning_delta":
			m.appendThinkingDelta(msg.Content)
		case "assistant":
			m.updateMetaFromData(msg.Data)
			content := msg.Content
			if content == "" && msg.Data != nil {
				content = eventContent(msg.Data)
			}
			if content != "" {
				m.replaceLastAssistant(content, eventString(msg.Data, "reasoning_content"))
			}
		case "meta":
			m.updateMetaFromData(msg.Data)
		case "error":
			m.appendMessage("error", msg.Content)
			m.running = false
		case "done":
			m.running = false
		}
		m.refreshViewport()
		if m.running {
			return m, m.waitEvent()
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.refreshViewport()
	return m, cmd
}

func (m *Model) updateViewportScroll(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyPgUp, tea.KeyPgDown, tea.KeyUp, tea.KeyDown:
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return cmd, true
		}
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress {
			switch msg.Button {
			case tea.MouseButtonWheelUp, tea.MouseButtonWheelDown, tea.MouseButtonWheelLeft, tea.MouseButtonWheelRight:
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				return cmd, true
			}
		}
	}
	return nil, false
}

func (m Model) View() string {
	return m.viewport.View()
}

func (m Model) respond(ctx context.Context, text string, generation int64) tea.Cmd {
	return func() tea.Msg {
		if m.runtime == nil {
			m.events <- Event{Generation: generation, Name: "error", Content: "runtime 未初始化"}
			return nil
		}
		_, err := m.runtime.RespondToConversation(ctx, m.session.Conversation, m.session.User, text, NewEventWriter(m.events, generation))
		if err != nil {
			m.events <- Event{Generation: generation, Name: "error", Content: err.Error()}
			return nil
		}
		m.events <- Event{Generation: generation, Name: "done"}
		return nil
	}
}

func (m Model) waitEvent() tea.Cmd {
	return func() tea.Msg { return <-m.events }
}

func (m *Model) handleSlashCommand(text string) bool {
	switch strings.TrimSpace(text) {
	case "/help":
		m.appendMessage("system", "命令：/help /clear /cwd /skills /mcp /resume。Enter 发送，Ctrl+C 中断或退出。")
		return true
	case "/clear":
		m.messages = nil
		m.resumeSelecting = false
		m.resumeCandidates = nil
		m.appendMessage("system", "已清空当前 TUI 显示上下文")
		return true
	case "/cwd":
		m.appendMessage("system", "当前工作区："+m.session.CWD)
		return true
	case "/skills":
		m.appendMessage("system", renderSkillDetails(m.session.Skills, m.session.SkillCount))
		return true
	case "/mcp":
		m.appendMessage("system", renderMCPDetails(m.session.MCPServers, m.session.MCPToolCount))
		return true
	case "/resume":
		m.startResumeSelection()
		return true
	}
	m.appendMessage("system", "未知命令："+text)
	return true
}

func (m *Model) startResumeSelection() {
	if m.running {
		m.appendMessage("system", "当前正在生成，请先 Ctrl+C 中断后再执行 /resume")
		return
	}
	if m.session.Resumer == nil {
		m.appendMessage("system", "当前运行环境不支持 /resume")
		return
	}
	sessions, err := m.session.Resumer.ListResumableSessions(context.Background(), m.session.CWD)
	if err != nil {
		m.appendMessage("error", err.Error())
		return
	}
	if len(sessions) == 0 {
		m.resumeSelecting = false
		m.resumeCandidates = nil
		m.appendMessage("system", "当前目录暂无可恢复的历史会话")
		return
	}
	m.resumeSelecting = true
	m.resumeCandidates = sessions
	m.appendMessage("system", renderResumeCandidates(sessions))
}

func (m *Model) handleResumeSelection(text string) bool {
	text = strings.TrimSpace(text)
	if text == "/cancel" || text == "/clear" {
		m.resumeSelecting = false
		m.resumeCandidates = nil
		if text == "/clear" {
			m.messages = nil
			m.appendMessage("system", "已清空当前 TUI 显示上下文")
		} else {
			m.appendMessage("system", "已取消恢复历史会话")
		}
		return true
	}
	idx, err := strconv.Atoi(text)
	if err != nil || idx < 1 || idx > len(m.resumeCandidates) {
		m.appendMessage("system", fmt.Sprintf("请输入 1-%d 之间的序号，或输入 /cancel 取消", len(m.resumeCandidates)))
		return true
	}
	candidate := m.resumeCandidates[idx-1]
	conv, history, err := m.session.Resumer.ResumeSession(context.Background(), candidate.SessionID, m.session.CWD, m.session.User)
	if err != nil {
		m.appendMessage("error", err.Error())
		return true
	}
	m.session.Conversation = conv
	m.resumeSelecting = false
	m.resumeCandidates = nil
	m.messages = messagesForDisplay(history)
	m.appendMessage("system", fmt.Sprintf("已恢复历史会话：%s", conv.SessionID))
	return true
}

func renderResumeCandidates(sessions []storage.ResumableSession) string {
	var b strings.Builder
	b.WriteString("可恢复的历史会话：")
	for i, session := range sessions {
		title := strings.TrimSpace(session.Title)
		if title == "" {
			title = "TUI 会话"
		}
		updated := "unknown"
		if !session.UpdatedAt.IsZero() {
			updated = session.UpdatedAt.Local().Format(time.RFC3339)
		}
		b.WriteString(fmt.Sprintf("\n%d. %s | %s | 消息:%d | %s", i+1, updated, title, session.MessageCount, session.SessionID))
	}
	b.WriteString("\n输入序号恢复，或输入 /cancel 取消。")
	return b.String()
}

func messagesForDisplay(history []storage.Message) []Message {
	if len(history) == 0 {
		return nil
	}
	messages := make([]Message, 0, len(history))
	for _, msg := range history {
		switch msg.Role {
		case "user", "assistant", "system", "error":
			messages = append(messages, Message{Role: msg.Role, Content: msg.Content, ReasoningContent: msg.ReasoningContent})
		}
	}
	return messages
}

func renderSkillDetails(skills []sessions.SkillSummary, fallbackCount int) string {
	count := len(skills)
	if count == 0 && fallbackCount > 0 {
		count = fallbackCount
	}
	if len(skills) == 0 {
		return fmt.Sprintf("已加载 Skills：%d 个", count)
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("已加载 Skills：%d 个", len(skills)))
	for _, skill := range skills {
		b.WriteString("\n- ")
		b.WriteString(skill.Name)
		if strings.TrimSpace(skill.Source) != "" {
			b.WriteString(" [")
			b.WriteString(skill.Source)
			b.WriteString("]")
		}
		if strings.TrimSpace(skill.Description) != "" {
			b.WriteString(" ")
			b.WriteString(skill.Description)
		}
		if strings.TrimSpace(skill.Path) != "" {
			b.WriteString("\n  path: ")
			b.WriteString(skill.Path)
		}
	}
	return b.String()
}

func renderMCPDetails(servers []mcp.ServerStatus, toolCount int) string {
	if len(servers) == 0 {
		return fmt.Sprintf("MCP Servers：0 个，工具：%d 个", toolCount)
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("MCP Servers：%d 个，工具：%d 个", len(servers), toolCount))
	for _, server := range servers {
		state := "disabled"
		if server.Connected {
			state = "connected"
		} else if server.LastError != "" {
			state = "failed"
		} else if server.Enabled {
			state = "pending"
		}
		b.WriteString("\n- ")
		b.WriteString(server.Name)
		b.WriteString(" [")
		b.WriteString(server.Transport)
		b.WriteString("] ")
		b.WriteString(state)
		b.WriteString(fmt.Sprintf(", tools: %d", server.ToolCount))
		if server.Command != "" {
			b.WriteString("\n  command: ")
			b.WriteString(strings.Join(append([]string{server.Command}, server.Args...), " "))
		}
		if server.URL != "" {
			b.WriteString("\n  url: ")
			b.WriteString(server.URL)
		}
		if server.LastError != "" {
			b.WriteString("\n  error: ")
			b.WriteString(server.LastError)
		}
	}
	return b.String()
}

func (m *Model) appendMessage(role, content string) {
	m.messages = append(m.messages, Message{Role: role, Content: content})
}

func (m *Model) appendAssistantDelta(delta string) {
	if len(m.messages) == 0 || !isLiveAssistantRole(m.messages[len(m.messages)-1].Role) {
		m.appendMessage("assistant", delta)
		return
	}
	if m.messages[len(m.messages)-1].Role == "thinking" {
		m.messages[len(m.messages)-1].Role = "assistant"
	}
	m.messages[len(m.messages)-1].Content += delta
}

func (m *Model) appendThinkingDelta(delta string) {
	if delta == "" {
		return
	}
	if len(m.messages) == 0 || !isLiveAssistantRole(m.messages[len(m.messages)-1].Role) {
		m.messages = append(m.messages, Message{Role: "assistant", ReasoningContent: delta})
		return
	}
	m.messages[len(m.messages)-1].ReasoningContent += delta
}

func (m *Model) replaceLastAssistant(content, reasoning string) {
	if len(m.messages) == 0 || !isLiveAssistantRole(m.messages[len(m.messages)-1].Role) {
		m.messages = append(m.messages, Message{Role: "assistant", Content: content, ReasoningContent: reasoning})
		return
	}
	m.messages[len(m.messages)-1].Role = "assistant"
	m.messages[len(m.messages)-1].Content = content
	m.messages[len(m.messages)-1].ReasoningContent = reasoning
}

func isLiveAssistantRole(role string) bool {
	return role == "assistant" || role == "thinking"
}

func (m *Model) refreshViewport() {
	m.viewport.Width = max(10, m.width)
	m.viewport.Height = m.viewportHeight()
	m.viewport.SetContent(m.renderTranscript())
	m.viewport.GotoBottom()
}

func (m Model) viewportHeight() int {
	if m.height <= 0 {
		return 20
	}
	return max(1, m.height)
}

func (m Model) renderTranscript() string {
	var b strings.Builder
	b.WriteString(m.renderWelcome())
	if len(m.messages) > 0 {
		b.WriteString("\n\n")
		b.WriteString(m.renderMessages())
	}
	b.WriteString("\n")
	b.WriteString(m.renderInput())
	b.WriteString("\n")
	b.WriteString(subtleStyle().Render(m.renderLiveStatus()))
	return b.String()
}

func (m Model) renderMessages() string {
	var b strings.Builder
	for _, msg := range m.messages {
		b.WriteString(m.renderMessage(msg))
		b.WriteString("\n\n")
	}
	return b.String()
}

func (m Model) renderMessage(msg Message) string {
	switch msg.Role {
	case "user":
		return promptLineStyle().Render("›") + " " + userStyle().Render(wrapText(msg.Content, m.messageWidth()-2))
	case "assistant":
		content := msg.Content
		if m.renderer != nil {
			if rendered, err := m.renderer.Render(content); err == nil {
				content = wrapText(strings.TrimSpace(rendered), m.messageWidth())
			}
		} else {
			content = wrapText(content, m.messageWidth())
		}
		if m.running && strings.TrimSpace(msg.ReasoningContent) != "" {
			content = thinkingStyle().Render("✽ 思考中\n"+wrapText(strings.TrimSpace(msg.ReasoningContent), m.messageWidth()-2)) + "\n" + content
		}
		return assistantLeadStyle().Render("go-agent") + "\n" + content
	case "thinking":
		return thinkingStyle().Render("✽ 思考中\n" + wrapText(msg.Content, m.messageWidth()-2))
	case "system":
		return systemStyle().Render("• " + wrapText(msg.Content, m.messageWidth()-2))
	case "error":
		return errorStyle().Render("✗ " + wrapText(msg.Content, m.messageWidth()-2))
	default:
		return roleLabel(msg.Role, lipgloss.Color("245")) + "\n" + wrapText(msg.Content, m.messageWidth())
	}
}

func (m Model) messageWidth() int {
	return max(10, m.width)
}

func newMarkdownRenderer(width int) *glamour.TermRenderer {
	renderer, _ := glamour.NewTermRenderer(glamour.WithStandardStyle("dark"), glamour.WithWordWrap(max(10, width)))
	return renderer
}

func wrapText(text string, width int) string {
	width = max(1, width)
	return ansi.Hardwrap(text, width, true)
}

func (m Model) renderWelcome() string {
	art := strings.Join([]string{
		`   /\_/\\`,
		`  ( o.o )   nano, but cozy`,
		`   > ^ <    ask · think · build`,
	}, "\n")
	intro := "像 Claude Code 一样一问一答：在终端直接提问，回答会流式显示在下方。"
	quick := "快捷键：Enter 发送 · Ctrl+C 中断/退出 · /resume 恢复 · /clear 清屏"
	stats := fmt.Sprintf("Skills %d · MCP tools %d", m.session.SkillCount, m.session.MCPToolCount)
	return startupPanelStyle().Width(max(20, m.width-2)).Render(accentArtStyle().Render(art) + "\n\n" + intro + "\n" + subtleStyle().Render(quick) + "\n" + subtleStyle().Render(stats))
}

func (m Model) renderHeader() string {
	width := max(20, m.width)
	status := runningText(m.running)
	contextText := "上下文 --"
	if m.contextBudget > 0 {
		contextText = fmt.Sprintf("上下文 %d%% · %s/%s", min(100, m.contextTokens*100/m.contextBudget), compactNumber(m.contextTokens), compactNumber(m.contextBudget))
	}
	left := titleStyle().Render("✦ go-agent")
	right := headerMetaStyle().Render(fmt.Sprintf("%s  ·  本轮工具 %d  ·  %s", status, m.toolCallCount, contextText))
	gap := max(1, width-lipgloss.Width(left)-lipgloss.Width(right)-4)
	line := left + strings.Repeat(" ", gap) + right
	workspace := subtleStyle().Render(fmt.Sprintf("cwd %s · skills %d · mcp tools %d", m.session.CWD, m.session.SkillCount, m.session.MCPToolCount))
	return headerStyle().Width(width).Render(line + "\n" + workspace)
}

func (m Model) renderConversationFrame() string {
	return conversationStyle().Width(max(10, m.width)).Render(m.viewport.View())
}

func (m Model) renderInput() string {
	prompt := inputPromptStyle().Render("›")
	text := strings.TrimSpace(m.input.Value())
	if text == "" {
		text = inputPromptStyle().Render(inputCursor) + " " + subtleStyle().Render(m.input.Placeholder)
	} else {
		text = userStyle().Render(text) + inputPromptStyle().Render(inputCursor)
	}
	return inputLineStyle().Width(max(10, m.width)).Render(prompt + " " + text)
}

func isTerminalProbeResponseInput(msg tea.KeyMsg) bool {
	if msg.Type != tea.KeyRunes {
		return false
	}
	value := string(msg.Runes)
	return strings.Contains(value, ";rgb:") && strings.Contains(value, "/")
}

func (m Model) renderLiveStatus() string {
	parts := []string{"Enter 发送", "Ctrl+C 中断/退出", "/help", fmt.Sprintf("工具 %d", m.toolCallCount)}
	if m.contextBudget > 0 {
		parts = append(parts, fmt.Sprintf("上下文 %d%% · %s/%s", min(100, m.contextTokens*100/m.contextBudget), compactNumber(m.contextTokens), compactNumber(m.contextBudget)))
	} else {
		parts = append(parts, "上下文 --")
	}
	return strings.Join(parts, " · ")
}

func (m *Model) updateMetaFromData(data any) {
	if count, ok := eventInt(data, "tool_call_count"); ok {
		m.toolCallCount = count
	}
	if tokens, ok := eventInt(data, "context_tokens"); ok {
		m.contextTokens = tokens
	}
	if budget, ok := eventInt(data, "context_budget"); ok {
		m.contextBudget = budget
	}
}

func eventString(data any, key string) string {
	m, ok := data.(map[string]any)
	if !ok {
		return ""
	}
	value, _ := m[key].(string)
	return value
}

func eventInt(data any, key string) (int, bool) {
	m, ok := data.(map[string]any)
	if !ok {
		return 0, false
	}
	switch value := m[key].(type) {
	case int:
		return value, true
	case int64:
		return int(value), true
	case float64:
		return int(value), true
	default:
		return 0, false
	}
}

func roleLabel(label string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Bold(true).Foreground(color).Render(label)
}

func headerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(tuiPalette.panel)
}

func titleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(tuiPalette.mint)
}

func headerMetaStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
}

func subtleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.muted)
}

func conversationStyle() lipgloss.Style {
	return lipgloss.NewStyle().Padding(1, 2, 0, 2)
}

func startupPanelStyle() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(tuiPalette.coral).Padding(1, 2).Margin(1, 0)
}

func userStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.ink)
}

func thinkingStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Italic(true).PaddingLeft(2)
}

func systemStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.butter)
}

func errorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.coral)
}

func promptLineStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.cyan).Bold(true)
}

func assistantLeadStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.mint).Bold(true)
}

func accentArtStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.lavender).Bold(true)
}

func inputPromptStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.mint).Bold(true)
}

func inputLineStyle() lipgloss.Style {
	return lipgloss.NewStyle().Background(lipgloss.Color("238")).Foreground(tuiPalette.ink)
}

func compactNumber(n int) string {
	if n >= 1000 {
		return fmt.Sprintf("%dk", n/1000)
	}
	return strconv.Itoa(n)
}

func runningText(running bool) string {
	if running {
		return "generating"
	}
	return "ready"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package tui\n\nimport (\n\t\"context\"\n\t\"fmt\"\n\t\"strconv\"\n\t\"strings\"\n\t\"time\"\n\n\t\"github.com/charmbracelet/bubbles/textarea\"\n\t\"github.com/charmbracelet/bubbles/viewport\"\n\ttea \"github.com/charmbracelet/bubbletea\"\n\t\"github.com/charmbracelet/glamour\"\n\t\"github.com/charmbracelet/lipgloss\"\n\t\"github.com/charmbracelet/x/ansi\"\n\n\t\"nano_cc/internal/agent/mcp\"\n\t\"nano_cc/internal/agent/runtime\"\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/logger\"\n\t\"nano_cc/internal/sessions\"\n)\n\ntype SessionInfo struct {\n\tUser"}
```

## 2026-06-13T09:15:55Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_04_S73VjNbZgl8k8jbvIOwh3758

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/style.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/style.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/style.go: no such file or directory"}
```

## 2026-06-13T09:15:55Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_05_BacQcXMHXMscjU7PbkKM9223

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/cli/root.go"}
```

### Result

```text
package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"nano_cc/internal/local"
	"nano_cc/internal/tui"
)

type Options struct {
	CWD string
}

type Runner struct {
	RunTUI func(context.Context, Options) error
}

func Run(ctx context.Context, args []string, defaultCWD string, out io.Writer, runner Runner) error {
	if out == nil {
		out = io.Discard
	}
	if runner.RunTUI == nil {
		runner.RunTUI = runTUI
	}
	if len(args) > 0 && args[0] == "tui" {
		args = args[1:]
	}
	if len(args) > 0 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		printHelp(out)
		return nil
	}

	fs := flag.NewFlagSet("go-agent", flag.ContinueOnError)
	fs.SetOutput(out)
	cwd := fs.String("cwd", defaultCWD, "工作区目录，默认是当前目录")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("unknown command: %s", strings.Join(fs.Args(), " "))
	}
	return runner.RunTUI(ctx, Options{CWD: *cwd})
}

func Main() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	return Run(context.Background(), os.Args[1:], cwd, os.Stdout, Runner{})
}

func runTUI(ctx context.Context, opts Options) error {
	bundle, err := local.Bootstrap(ctx, opts.CWD)
	if err != nil {
		return err
	}
	defer bundle.Close()
	return tui.Run(ctx, bundle.Runtime, tui.SessionInfo{
		User:         bundle.User,
		Conversation: bundle.Conversation,
		CWD:          bundle.CWD,
		Resumer:      bundle.Store,
		Skills:       bundle.Skills,
		MCPServers:   bundle.MCPServers,
		SkillCount:   bundle.SkillCount,
		MCPToolCount: bundle.MCPToolCount,
	})
}

func printHelp(out io.Writer) {
	_, _ = fmt.Fprintln(out, `go-agent - 本地 TUI 代码助手

用法：
  go-agent [--cwd <path>]   启动 TUI，默认工作区为当前目录
  go-agent tui [--cwd <path>] 同上
  go-agent help             显示帮助`)
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package cli\n\nimport (\n\t\"context\"\n\t\"flag\"\n\t\"fmt\"\n\t\"io\"\n\t\"os\"\n\t\"strings\"\n\n\t\"nano_cc/internal/local\"\n\t\"nano_cc/internal/tui\"\n)\n\ntype Options struct {\n\tCWD string\n}\n\ntype Runner struct {\n\tRunTUI func(context.Context, Options) error\n}\n\nfunc Run(ctx context.Context, args []string, defaultCWD string, out io.Writer, runner Runner) error {\n\tif out == nil {\n\t\tout = io.Discard\n\t}\n\tif runner.RunTUI == nil {\n\t\trunner.RunTUI = runTUI\n\t}\n\tif len(args) \u003e 0 \u0026\u0026 args[0] == \"tui\" {\n\t\targs = args[1:]\n\t}\n\tif len(args"}
```

## 2026-06-13T09:15:59Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_18BnILRh9J5qACjHWvVT9597

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/storage.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/storage.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/storage.go: no such file or directory"}
```

## 2026-06-13T09:15:59Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_RE12Z34vt51eyqEFPAhd7000

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/sqlite.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/sqlite.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/sqlite.go: no such file or directory"}
```

## 2026-06-13T09:15:59Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_FSOcrcOioGBJvmgeYbfG2838

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/models.go"}
```

### Result

```text
package storage

import "time"

type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	MemoryEnabled bool      `json:"memory_enabled"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// MCPServer 是用户配置的一个 MCP 服务器连接。Transport 取值 stdio/sse/streamable。
// Args/Env/Headers 从本地 .link/.mcp.json 读取，在 Go 侧以强类型表示。
type MCPServer struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Name      string            `json:"name"`
	Transport string            `json:"transport"`
	Command   string            `json:"command"`
	Args      []string          `json:"args"`
	Env       map[string]string `json:"env"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Enabled   bool              `json:"enabled"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type Conversation struct {
	ID            string    `json:"id"`
	SessionID     string    `json:"session_id,omitempty"`
	UserID        string    `json:"user_id"`
	RootMessageID string    `json:"root_message_id"`
	Title         string    `json:"title"`
	HistoryJSON   string    `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ResumableSession struct {
	SessionID      string    `json:"session_id"`
	ConversationID string    `json:"conversation_id"`
	WorkspaceRoot  string    `json:"workspace_root"`
	Title          string    `json:"title"`
	MessageCount   int       `json:"message_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// 该表已经不使用，作为存储的结构体，用于存储消息历史
type Message struct {
	ID               string            `json:"id"`
	ConversationID   string            `json:"conversation_id"`
	UserID           string            `json:"user_id"`
	Role             string            `json:"role"`
	Content          string            `json:"content"`
	ReasoningContent string            `json:"reasoning_content,omitempty"`
	ToolCallID       string            `json:"tool_call_id,omitempty"`
	ToolCalls        []MessageToolCall `json:"tool_calls,omitempty"`
	Meta             *MessageMeta      `json:"meta,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
}

// MessageMeta 记录助手回复的元信息，仅对最终 assistant 消息填充。
type MessageMeta struct {
	ToolCallCount int `json:"tool_call_count"`          // 调用工具次数（0 也需序列化，保证历史展示一致）
	ContextTokens int `json:"context_tokens,omitempty"` // 当前上下文估算 token
	ContextBudget int `json:"context_budget,omitempty"` // 上下文预算（用于算占比）
}

type MessageToolCall struct {
	ID       string              `json:"id"`
	Type     string              `json:"type"`
	Function MessageFunctionCall `json:"function"`
}

type MessageFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolCall struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	ToolName       string    `json:"tool_name"`
	Status         string    `json:"status"`
	Summary        string    `json:"summary"`
	CreatedAt      time.Time `json:"created_at"`
}

type PersistedOutput struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	MessageID      string    `json:"message_id"`
	ToolCallID     string    `json:"tool_call_id"`
	Kind           string    `json:"kind"`
	Strategy       string    `json:"strategy"`
	OriginalBytes  int       `json:"original_bytes"`
	ContentSHA256  string    `json:"content_sha256"`
	Content        string    `json:"content"`
	Preview        string    `json:"preview"`
	CreatedAt      time.Time `json:"created_at"`
}

type ToolResultLogEntry struct {
	ConversationID string    `json:"conversation_id"`
	SessionID      string    `json:"session_id,omitempty"`
	UserID         string    `json:"user_id"`
	ToolCallID     string    `json:"tool_call_id"`
	ToolName       string    `json:"tool_name"`
	RawArgs        string    `json:"raw_args"`
	Status         string    `json:"status"`
	Result         string    `json:"result"`
	AuditSummary   string    `json:"audit_summary,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type ContextSummary struct {
	ID                    string    `json:"id"`
	ConversationID        string    `json:"conversation_id"`
	UserID                string    `json:"user_id"`
	SourceHistorySHA256   string    `json:"source_history_sha256"`
	Strategy              string    `json:"strategy"`
	EstimatedTokensBefore int       `json:"estimated_tokens_before"`
	EstimatedTokensAfter  int       `json:"estimated_tokens_after"`
	Summary               string    `json:"summary"`
	CreatedAt             time.Time `json:"created_at"`
}

type SubagentMessage struct {
	ID               string            `json:"id"`
	RunID            string            `json:"run_id"`
	ParentToolCallID string            `json:"parent_tool_call_id"`
	ConversationID   string            `json:"conversation_id"`
	UserID           string            `json:"user_id"`
	SequenceNo       int               `json:"sequence_no"`
	Role             string            `json:"role"`
	Content          string            `json:"content"`
	ReasoningContent string            `json:"reasoning_content,omitempty"`
	ToolCallID       string            `json:"tool_call_id,omitempty"`
	ToolCalls        []MessageToolCall `json:"tool_calls,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
}

type Memory struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ConversationMemory 是单个会话维度、随每轮对话增量维护的“当前会话主干信息”
// 条目。它不注入 system prompt，仅在上下文压缩触发全量摘要时作为替代品使用。
type ConversationMemory struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package storage\n\nimport \"time\"\n\ntype User struct {\n\tID            string    `json:\"id\"`\n\tEmail         string    `json:\"email\"`\n\tUsername      string    `json:\"username\"`\n\tMemoryEnabled bool      `json:\"memory_enabled\"`\n\tCreatedAt     time.Time `json:\"created_at\"`\n\tUpdatedAt     time.Time `json:\"updated_at\"`\n}\n\n// MCPServer 是用户配置的一个 MCP 服务器连接。Transport 取值 stdio/sse/streamable。\n// Args/Env/Headers 从本地 .link/.mcp.json 读取，在 Go 侧以强类型表示。"}
```

## 2026-06-13T09:15:59Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_03_6c67dxX4j3SgdM9nKpcY9525

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/bootstrap.go"}
```

### Result

```text
package local

import (
	"context"
	"fmt"
	"time"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/runtime"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/assistant"
	"nano_cc/internal/config"
	"nano_cc/internal/idgen"
	"nano_cc/internal/llm"
	"nano_cc/internal/logger"
	"nano_cc/internal/sessions"
)

const LocalUserID = "local-user"

type Bundle struct {
	Runtime      *runtime.Service
	Store        *Store
	MCP          *mcp.Manager
	User         storage.User
	Conversation storage.Conversation
	CWD          string
	SkillCount   int
	Skills       []sessions.SkillSummary
	MCPToolCount int
	MCPServers   []mcp.ServerStatus
}

func Bootstrap(ctx context.Context, cwd string) (*Bundle, error) {
	cfg, err := config.LoadLocalConfig(cwd)
	if err != nil {
		return nil, err
	}
	if err := config.EnsureAppLayout(cfg); err != nil {
		return nil, err
	}
	if err := config.ValidateAppLayout(cfg); err != nil {
		return nil, err
	}
	if err := logger.InitFileLoggerAt(cfg.LogsDir); err != nil {
		logger.Warn(fmt.Sprintf("failed to init file logger: %v", err))
	}
	userSkillsDir, err := config.LinkSkillsDir()
	if err != nil {
		return nil, err
	}
	builtinSkills, err := sessions.LoadSkillsFromDirs([]sessions.SkillDir{
		{Path: userSkillsDir, Source: "user"},
		{Path: config.WorkspaceLinkSkillsDir(cfg.WorkspaceRoot), Source: "workspace"},
	})
	if err != nil {
		return nil, fmt.Errorf("load .link skills: %w", err)
	}
	basePrompt, err := assistant.LoadBaseSystemPrompt(cfg.SystemPromptPath)
	if err != nil {
		return nil, fmt.Errorf("load system prompt: %w", err)
	}
	linkMarkdown, err := config.LoadLinkMarkdownContext(cfg.WorkspaceRoot)
	if err != nil {
		return nil, fmt.Errorf("load LINK.MD context: %w", err)
	}
	store, err := NewStoreWithMemory(cfg.WorkspaceRoot)
	if err != nil {
		return nil, fmt.Errorf("init memory store: %w", err)
	}
	client := llm.NewDeepseekClient(cfg.LLM.BaseURL, cfg.LLM.APIKey)
	runtimeService := runtime.NewService(store, cfg, client)
	runtimeService.EnableMemory = true
	runtimeService.SetBuiltinSkills(builtinSkills)
	runtimeService.SetBasePrompt(basePrompt)
	runtimeService.SetLinkMarkdownContext(linkMarkdown)
	mcpManager := mcp.NewManager()
	workspaceMCPServers, err := mcp.LoadWorkspaceConfig(config.WorkspaceMCPConfigPath(cfg.WorkspaceRoot))
	if err != nil {
		mcpManager.Close()
		return nil, fmt.Errorf("load workspace mcp config: %w", err)
	}
	mcpManager.SetWorkspaceServers(ctx, workspaceMCPServers)
	runtimeService.SetMCPManager(mcpManager)
	user := storage.User{ID: LocalUserID, Email: "local@go-agent", Username: "local", MemoryEnabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	conversation := storage.Conversation{ID: idgen.New("conv"), SessionID: idgen.UUID(), UserID: user.ID, RootMessageID: idgen.New("msg"), Title: "TUI 会话", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := store.CreateConversation(ctx, conversation); err != nil {
		mcpManager.Close()
		return nil, err
	}
	mcpManager.EnsureWorkspaceSessions(ctx)
	mcpTools := mcpManager.ToolsForUser(user.ID)
	mcpSnapshot := mcpManager.Snapshot(user.ID)
	skills := builtinSkills.Summaries()
	return &Bundle{Runtime: runtimeService, Store: store, MCP: mcpManager, User: user, Conversation: conversation, CWD: cfg.WorkspaceRoot, SkillCount: len(skills), Skills: skills, MCPToolCount: len(mcpTools), MCPServers: mcpSnapshot.Servers}, nil
}

func (b *Bundle) Close() {
	if b != nil && b.MCP != nil {
		b.MCP.Close()
	}
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package local\n\nimport (\n\t\"context\"\n\t\"fmt\"\n\t\"time\"\n\n\t\"nano_cc/internal/agent/mcp\"\n\t\"nano_cc/internal/agent/runtime\"\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/assistant\"\n\t\"nano_cc/internal/config\"\n\t\"nano_cc/internal/idgen\"\n\t\"nano_cc/internal/llm\"\n\t\"nano_cc/internal/logger\"\n\t\"nano_cc/internal/sessions\"\n)\n\nconst LocalUserID = \"local-user\"\n\ntype Bundle struct {\n\tRuntime      *runtime.Service\n\tStore        *Store\n\tMCP          *mcp.Manager\n\tUser         storage.User\n\tConversation storage.Co"}
```

## 2026-06-13T09:15:59Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_04_QjK6CkKIJoobjuwZgnRX9581

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/config.go"}
```

### Result

```text
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"-"`
	ModelID string `json:"model_id"`
}

type fileConfig struct {
	Config
	AppHome                            string `json:"app_home"`
	BuiltinSkillsDir                   string `json:"builtin_skills_dir"`
	CommandBinDir                      string `json:"command_bin_dir"`
	CommandScriptDir                   string `json:"command_script_dir"`
	SystemPromptPath                   string `json:"system_prompt_path"`
	WorkspaceRoot                      string `json:"workspace_root"`
	AllowedTools                       string `json:"allowed_tools"`
	BashAllowOutsideWorkspace          bool   `json:"bash_allow_outside_workspace"`
	BashAllowDangerousCommands         bool   `json:"bash_allow_dangerous_commands"`
	ConversationLockTTLSeconds         int    `json:"conversation_lock_ttl_seconds"`
	MemoryWorkTimeoutSeconds           int    `json:"memory_work_timeout_seconds"`
	ConversationLockWaitTimeoutSeconds int    `json:"conversation_lock_wait_timeout_seconds"`
}

type AppConfig struct {
	LLM                         Config
	AppHome                     string
	BuiltinSkillsDir            string
	CommandBinDir               string
	CommandScriptDir            string
	SystemPromptPath            string
	WorkspaceRoot               string
	LogsDir                     string
	AllowedTools                []string
	BashAllowOutsideWorkspace   bool
	BashAllowDangerousCommands  bool
	ConversationLockTTL         time.Duration
	MemoryWorkTimeout           time.Duration
	ConversationLockWaitTimeout time.Duration
}

func loadConfigFile() (fileConfig, error) {
	data, err := os.ReadFile(configFilePath())
	if err != nil {
		return fileConfig{}, err
	}

	var cfg fileConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fileConfig{}, fmt.Errorf("failed to parse workspace config.json: %w", err)
	}
	return cfg, nil
}

func configFilePath() string {
	return "config.json"
}

func loadLLMConfig(fileCfg fileConfig) (Config, error) {
	cfg := Config{
		BaseURL: strings.TrimSpace(fileCfg.Config.BaseURL),
		APIKey:  strings.TrimSpace(getenv("OPENAI_API_KEY")),
		ModelID: strings.TrimSpace(fileCfg.Config.ModelID),
	}

	if cfg.BaseURL == "" || cfg.ModelID == "" {
		return Config{}, fmt.Errorf("missing LLM config; set base_url, model_id in config.json")
	}
	if cfg.APIKey == "" {
		return Config{}, fmt.Errorf("missing LLM api key; set OPENAI_API_KEY environment variable")
	}

	return cfg, nil
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package config\n\nimport (\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"os\"\n\t\"strings\"\n\t\"time\"\n)\n\ntype Config struct {\n\tBaseURL string `json:\"base_url\"`\n\tAPIKey  string `json:\"-\"`\n\tModelID string `json:\"model_id\"`\n}\n\ntype fileConfig struct {\n\tConfig\n\tAppHome                            string `json:\"app_home\"`\n\tBuiltinSkillsDir                   string `json:\"builtin_skills_dir\"`\n\tCommandBinDir                      string `json:\"command_bin_dir\"`\n\tCommandScriptDir                   string `json:\"command_script_dir\"`\n"}
```

## 2026-06-13T09:16:03Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_06Ojm3sZiX6gY3OqOQuA3563

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager.go"}
```

### Result

```text
package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/logger"
)

const (
	toolNamePrefix  = "mcp__"
	connectTimeout  = 30 * time.Second
	callTimeout     = 120 * time.Second
	idleTimeout     = 10 * time.Minute
	cleanupInterval = 1 * time.Minute
)

// serverSession 持有一个已连接 MCP 服务器的会话与其工具元信息。
type serverSession struct {
	server      storage.MCPServer
	signature   string // 配置指纹，用于判断配置是否变更需重连
	session     *mcpsdk.ClientSession
	tools       []openai.Tool     // 已转换并加前缀的工具定义
	toolNames   map[string]string // 前缀名 -> 原始 MCP 工具名
	lastUsedAt  time.Time
	activeCalls int
	closing     bool
}

type ServerStatus struct {
	ID        string
	Name      string
	Scope     string
	Transport string
	Command   string
	Args      []string
	URL       string
	Enabled   bool
	Connected bool
	ToolCount int
	LastError string
}

type Snapshot struct {
	Servers    []ServerStatus
	ToolCount  int
	ErrorCount int
}

// Manager 管理本地 MCP 客户端连接，提供工具发现与调用能力。
type Manager struct {
	mu                sync.Mutex
	done              chan struct{}
	closeOnce         sync.Once
	builtinServers    map[string]storage.MCPServer
	builtinSessions   map[string]*serverSession
	workspaceServers  map[string]storage.MCPServer
	workspaceSessions map[string]*serverSession
	workspaceErrors   map[string]string
}

func NewManager() *Manager {
	manager := &Manager{
		done:              make(chan struct{}),
		builtinServers:    make(map[string]storage.MCPServer),
		builtinSessions:   make(map[string]*serverSession),
		workspaceServers:  make(map[string]storage.MCPServer),
		workspaceSessions: make(map[string]*serverSession),
		workspaceErrors:   make(map[string]string),
	}
	go manager.cleanupLoop()
	return manager
}

func serverSignature(s storage.MCPServer) string {
	data, _ := json.Marshal(struct {
		Transport string
		Command   string
		Args      []string
		Env       map[string]string
		URL       string
		Headers   map[string]string
	}{s.Transport, s.Command, s.Args, s.Env, s.URL, s.Headers})
	return string(data)
}

func sanitizeName(name string) string {
	var sb strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_':
			sb.WriteRune(r)
		default:
			sb.WriteRune('_')
		}
	}
	if sb.Len() == 0 {
		return "server"
	}
	return sb.String()
}

func prefixedToolName(serverName, toolName string) string {
	return toolNamePrefix + sanitizeName(serverName) + "__" + toolName
}

func (m *Manager) SetBuiltinServers(ctx context.Context, servers []storage.MCPServer) {
	m.mu.Lock()
	old := m.builtinSessions
	m.builtinServers = make(map[string]storage.MCPServer, len(servers))
	m.builtinSessions = make(map[string]*serverSession)
	for _, server := range servers {
		m.builtinServers[server.ID] = server
	}
	m.mu.Unlock()

	for _, sess := range old {
		closeSession(sess)
	}
	m.EnsureBuiltinSessions(ctx)
}

func (m *Manager) EnsureBuiltinSessions(ctx context.Context) {
	m.mu.Lock()
	missing := make([]storage.MCPServer, 0)
	for id, server := range m.builtinServers {
		if !server.Enabled {
			continue
		}
		if _, ok := m.builtinSessions[id]; !ok {
			missing = append(missing, server)
		}
	}
	m.mu.Unlock()

	for _, server := range missing {
		sess, err := connectBuiltinAndDiscover(ctx, server)
		if err != nil {
			logger.Warn(fmt.Sprintf("mcp: builtin stdio connect failed server=%s: %v", server.Name, err))
			continue
		}
		m.mu.Lock()
		if old := m.builtinSessions[server.ID]; old != nil {
			closeSession(old)
		}
		m.builtinSessions[server.ID] = sess
		m.mu.Unlock()
		logger.Info(fmt.Sprintf("mcp: builtin stdio connected server=%s tools=%d", server.Name, len(sess.tools)))
	}
}

func (m *Manager) SetWorkspaceServers(ctx context.Context, servers []storage.MCPServer) {
	m.mu.Lock()
	old := m.workspaceSessions
	m.workspaceServers = make(map[string]storage.MCPServer, len(servers))
	m.workspaceSessions = make(map[string]*serverSession)
	m.workspaceErrors = make(map[string]string)
	for _, server := range servers {
		m.workspaceServers[server.ID] = server
	}
	m.mu.Unlock()

	for _, sess := range old {
		closeSession(sess)
	}
	m.EnsureWorkspaceSessions(ctx)
}

func (m *Manager) EnsureWorkspaceSessions(ctx context.Context) {
	m.mu.Lock()
	missing := make([]storage.MCPServer, 0)
	for id, server := range m.workspaceServers {
		if !server.Enabled {
			continue
		}
		if _, ok := m.workspaceSessions[id]; !ok {
			missing = append(missing, server)
		}
	}
	m.mu.Unlock()

	for _, server := range missing {
		sess, err := connectAndDiscover(ctx, server)
		m.mu.Lock()
		if err != nil {
			m.workspaceErrors[server.ID] = err.Error()
			m.mu.Unlock()
			logger.Warn(fmt.Sprintf("mcp: workspace connect failed server=%s: %v", server.Name, err))
			continue
		}
		delete(m.workspaceErrors, server.ID)
		if old := m.workspaceSessions[server.ID]; old != nil {
			closeSession(old)
		}
		m.workspaceSessions[server.ID] = sess
		m.mu.Unlock()
		logger.Info(fmt.Sprintf("mcp: workspace connected server=%s tools=%d", server.Name, len(sess.tools)))
	}
}

func connectAndDiscover(ctx context.Context, server storage.MCPServer) (*serverSession, error) {
	transport, err := buildTransport(server)
	if err != nil {
		return nil, err
	}
	connCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	client := mcpsdk.NewClient(&mcpsdk.Implementation{Name: "nano_cc", Version: "1.0.0"}, nil)
	session, err := client.Connect(connCtx, transport, nil)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	tools, names, err := discoverTools(connCtx, session, server.Name)
	if err != nil {
		_ = session.Close()
		return nil, fmt.Errorf("list tools: %w", err)
	}
	return &serverSession{
		server:     server,
		signature:  serverSignature(server),
		session:    session,
		tools:      tools,
		toolNames:  names,
		lastUsedAt: time.Now(),
	}, nil
}

func connectBuiltinAndDiscover(ctx context.Context, server storage.MCPServer) (*serverSession, error) {
	transport, err := buildBuiltinStdioTransport(server)
	if err != nil {
		return nil, err
	}
	return connectWithTransport(ctx, server, transport)
}

func connectWithTransport(ctx context.Context, server storage.MCPServer, transport mcpsdk.Transport) (*serverSession, error) {
	connCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	client := mcpsdk.NewClient(&mcpsdk.Implementation{Name: "nano_cc", Version: "1.0.0"}, nil)
	session, err := client.Connect(connCtx, transport, nil)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	tools, names, err := discoverTools(connCtx, session, server.Name)
	if err != nil {
		_ = session.Close()
		return nil, fmt.Errorf("list tools: %w", err)
	}
	return &serverSession{
		server:     server,
		signature:  serverSignature(server),
		session:    session,
		tools:      tools,
		toolNames:  names,
		lastUsedAt: time.Now(),
	}, nil
}

func discoverTools(ctx context.Context, session *mcpsdk.ClientSession, serverName string) ([]openai.Tool, map[string]string, error) {
	result, err := session.ListTools(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	tools := make([]openai.Tool, 0, len(result.Tools))
	names := make(map[string]string, len(result.Tools))
	for _, tool := range result.Tools {
		prefixed := prefixedToolName(serverName, tool.Name)
		params, err := json.Marshal(tool.InputSchema)
		if err != nil {
			logger.Warn(fmt.Sprintf("mcp: marshal schema failed server=%s tool=%s: %v", serverName, tool.Name, err))
			continue
		}
		tools = append(tools, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        prefixed,
				Description: tool.Description,
				Parameters:  json.RawMessage(params),
			},
		})
		names[prefixed] = tool.Name
	}
	return tools, names, nil
}

// ToolsForUser 返回所有已连接本地 MCP 服务器发现到的工具定义（带前缀）。
func (m *Manager) ToolsForUser(userID string) []openai.Tool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.builtinSessions) == 0 && len(m.workspaceSessions) == 0 {
		return nil
	}
	tools := make([]openai.Tool, 0)
	now := time.Now()
	for _, sess := range m.builtinSessions {
		sess.lastUsedAt = now
		tools = append(tools, sess.tools...)
	}
	for _, sess := range m.workspaceSessions {
		sess.lastUsedAt = now
		tools = append(tools, sess.tools...)
	}
	sort.Slice(tools, func(i, j int) bool {
		return tools[i].Function.Name < tools[j].Function.Name
	})
	return tools
}

// CallTool 把带前缀的工具名路由到对应 MCP 服务器执行，返回文本结果。
func (m *Manager) CallTool(ctx context.Context, userID, prefixedName, rawArgs string) (string, error) {
	m.mu.Lock()
	var target *serverSession
	var originalName string
	var builtinID string
	for _, sess := range m.workspaceSessions {
		if name, ok := sess.toolNames[prefixedName]; ok {
			target = sess
			originalName = name
			target.activeCalls++
			break
		}
	}
	if target == nil {
		for id, sess := range m.builtinSessions {
			if name, ok := sess.toolNames[prefixedName]; ok {
				target = sess
				originalName = name
				builtinID = id
				target.activeCalls++
				break
			}
		}
	}
	m.mu.Unlock()

	if target == nil {
		return "", fmt.Errorf("mcp tool %s not found", prefixedName)
	}
	defer m.finishCall(target)

	var args map[string]any
	if strings.TrimSpace(rawArgs) != "" {
		if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
			return "", fmt.Errorf("invalid tool arguments: %w", err)
		}
	}
	if builtinID != "" {
		return m.callBuiltinTool(ctx, builtinID, target, originalName, args)
	}

	callCtx, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()

	result, err := target.session.CallTool(callCtx, &mcpsdk.CallToolParams{Name: originalName, Arguments: args})
	if err != nil {
		return "", err
	}
	output := serializeContent(result)
	if result.IsError {
		return "", fmt.Errorf("mcp tool error: %s", output)
	}
	return output, nil
}

func (m *Manager) callBuiltinTool(ctx context.Context, builtinID string, target *serverSession, originalName string, args map[string]any) (string, error) {
	output, err := callSessionTool(ctx, target, originalName, args)
	if err == nil {
		return output, nil
	}
	logger.Warn(fmt.Sprintf("mcp: builtin stdio call failed, reconnecting server=%s: %v", target.server.Name, err))
	m.reconnectBuiltin(ctx, builtinID, target)

	m.mu.Lock()
	retryTarget := m.builtinSessions[builtinID]
	if retryTarget != nil {
		retryTarget.activeCalls++
	}
	m.mu.Unlock()
	if retryTarget == nil {
		return "", err
	}
	defer m.finishCall(retryTarget)
	return callSessionTool(ctx, retryTarget, originalName, args)
}

func callSessionTool(ctx context.Context, target *serverSession, originalName string, args map[string]any) (string, error) {
	callCtx, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	result, err := target.session.CallTool(callCtx, &mcpsdk.CallToolParams{Name: originalName, Arguments: args})
	if err != nil {
		return "", err
	}
	output := serializeContent(result)
	if result.IsError {
		return "", fmt.Errorf("mcp tool error: %s", output)
	}
	return output, nil
}

func (m *Manager) reconnectBuiltin(ctx context.Context, builtinID string, old *serverSession) {
	m.mu.Lock()
	server, ok := m.builtinServers[builtinID]
	if current := m.builtinSessions[builtinID]; current == old {
		delete(m.builtinSessions, builtinID)
	}
	m.mu.Unlock()
	if !ok || !server.Enabled {
		return
	}
	closeSession(old)
	sess, err := connectBuiltinAndDiscover(ctx, server)
	if err != nil {
		logger.Warn(fmt.Sprintf("mcp: builtin stdio reconnect failed server=%s: %v", server.Name, err))
		return
	}
	m.mu.Lock()
	m.builtinSessions[builtinID] = sess
	m.mu.Unlock()
	logger.Info(fmt.Sprintf("mcp: builtin stdio reconnected server=%s tools=%d", server.Name, len(sess.tools)))
}

// Invalidate 保留兼容入口。本地 TUI 不再维护用户 DB MCP 连接，因此该方法不影响 builtin/workspace 会话。
func (m *Manager) Invalidate(userID string) {
}

// TestServer 临时连接一个配置并发现工具，返回工具名列表，供连接测试使用。
func (m *Manager) TestServer(ctx context.Context, server storage.MCPServer) ([]string, error) {
	sess, err := connectAndDiscover(ctx, server)
	if err != nil {
		return nil, err
	}
	defer sess.session.Close()
	names := make([]string, 0, len(sess.toolNames))
	for _, original := range sess.toolNames {
		names = append(names, original)
	}
	sort.Strings(names)
	return names, nil
}

// Close 关闭全部连接，用于服务退出。
func (m *Manager) Close() {
	m.closeOnce.Do(func() { close(m.done) })
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, sess := range m.builtinSessions {
		m.closeWhenIdleLocked(sess)
	}
	for _, sess := range m.workspaceSessions {
		m.closeWhenIdleLocked(sess)
	}
	m.builtinSessions = make(map[string]*serverSession)
	m.workspaceSessions = make(map[string]*serverSession)
}

func (m *Manager) Snapshot(userID string) Snapshot {
	m.mu.Lock()
	defer m.mu.Unlock()
	statuses := make([]ServerStatus, 0, len(m.builtinServers)+len(m.workspaceServers))
	addConfigured := func(scope string, server storage.MCPServer, sess *serverSession, lastErr string) {
		status := ServerStatus{ID: server.ID, Name: server.Name, Scope: scope, Transport: server.Transport, Command: server.Command, Args: append([]string(nil), server.Args...), URL: server.URL, Enabled: server.Enabled, Connected: sess != nil, LastError: lastErr}
		if sess != nil {
			status.ToolCount = len(sess.tools)
		}
		statuses = append(statuses, status)
	}
	for id, server := range m.builtinServers {
		addConfigured("builtin", server, m.builtinSessions[id], "")
	}
	for id, server := range m.workspaceServers {
		addConfigured("workspace", server, m.workspaceSessions[id], m.workspaceErrors[id])
	}
	sort.Slice(statuses, func(i, j int) bool {
		if statuses[i].Scope != statuses[j].Scope {
			return statuses[i].Scope < statuses[j].Scope
		}
		return statuses[i].Name < statuses[j].Name
	})
	snapshot := Snapshot{Servers: statuses}
	for _, status := range statuses {
		snapshot.ToolCount += status.ToolCount
		if status.LastError != "" {
			snapshot.ErrorCount++
		}
	}
	return snapshot
}

func (m *Manager) finishCall(sess *serverSession) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if sess.activeCalls > 0 {
		sess.activeCalls--
	}
	sess.lastUsedAt = time.Now()
	if sess.closing && sess.activeCalls == 0 {
		closeSession(sess)
	}
}

func (m *Manager) closeWhenIdleLocked(sess *serverSession) {
	if sess == nil {
		return
	}
	sess.closing = true
	if sess.activeCalls == 0 {
		closeSession(sess)
	}
}

func (m *Manager) cleanupLoop() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			m.cleanupIdleSessions(time.Now())
		case <-m.done:
			return
		}
	}
}

func (m *Manager) cleanupIdleSessions(now time.Time) {
	type idleSession struct {
		scope string
		id    string
		sess  *serverSession
	}
	var idle []idleSession

	m.mu.Lock()
	for serverID, sess := range m.workspaceSessions {
		if sess.activeCalls > 0 || now.Sub(sess.lastUsedAt) < idleTimeout {
			continue
		}
		idle = append(idle, idleSession{scope: "workspace", id: serverID, sess: sess})
		delete(m.workspaceSessions, serverID)
	}
	m.mu.Unlock()

	for _, item := range idle {
		closeSession(item.sess)
		logger.Info(fmt.Sprintf("mcp: idle session closed scope=%s id=%s server=%s", item.scope, item.id, item.sess.server.Name))
	}
}

func closeSession(sess *serverSession) {
	if sess == nil || sess.session == nil {
		return
	}
	_ = sess.session.Close()
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package mcp\n\nimport (\n\t\"context\"\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"sort\"\n\t\"strings\"\n\t\"sync\"\n\t\"time\"\n\n\tmcpsdk \"github.com/modelcontextprotocol/go-sdk/mcp\"\n\topenai \"github.com/sashabaranov/go-openai\"\n\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/logger\"\n)\n\nconst (\n\ttoolNamePrefix  = \"mcp__\"\n\tconnectTimeout  = 30 * time.Second\n\tcallTimeout     = 120 * time.Second\n\tidleTimeout     = 10 * time.Minute\n\tcleanupInterval = 1 * time.Minute\n)\n\n// serverSession 持有一个已连接 MCP 服务器的会话\ufffd\ufffd"}
```

## 2026-06-13T09:16:03Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_9LvwycSygTopiOHvmqQt6430

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/assistant/prompt.go"}
```

### Result

```text
package assistant

import (
	"os"
	"strings"
)

const persistedOutputGuidance = "当较早的消息中出现 `<persisted-output ...>` 标记时，表示完整的工具输出已存入本地文件，内联的只是预览。" +
	"如果预览不足以完成任务，请用标记中的 id 和偏移量调用 `read_persisted_output` 分块读取更多内容，不要猜测被省略的部分。" +
	"当看到 `[Earlier result compacted. Re-run if needed]` 时，请重新执行相关工具以再次获取该结果。"

const DefaultBaseSystemPrompt = "你是 Link，一个运行在本地 TUI 中的通用型智能体（general-purpose agent），而不是只能聊天的助手。\n\n" +
	"帮助用户处理日常问答、分析、规划、写作，以及阅读代码、修改代码、生成文件等任务。优先给出直接、有用的回答；当工具或技能能帮助你完成任务时，主动使用它们。\n\n" +
	"你为用户修改代码或生成文件时，可能会在运行时的工作区里实际操作，但工作区只是对用户不可见的临时执行环境。完成后必须把修改后的完整代码、生成的文件内容原封不动地返回给用户，不要省略、截断或改写；返回之后工作区里的对应文件即可删除。\n\n" +
	"绝不向用户透传任何终端 / 命令行相关信息：包括但不限于\"执行 Shell 命令\"\"运行 bash\"\"终端操作\"\"命令行\"等表述，以及具体命令、脚本、工作目录、文件路径、运行环境、工具名称、技能调用方式等内部执行机制。描述你的工作时用\"我帮你阅读/修改了代码\"\"我为你生成了文件\"这类面向结果的说法。\n\n" +
	"当用户询问关于你自身的问题时，只介绍你能为用户提供的能力与价值（如阅读代码、修改代码、生成文件等）；即使运行时提供了上述信息，也仅供内部决策使用，绝不复述给用户，除非用户为完成具体任务确实需要。"

type PromptOptions struct {
	BasePrompt        string
	Surface           string
	SkillDescriptions string
	MemorySection     string
	WorkingDirectory  string
	LinkMarkdown      LinkMarkdownContext
	ToolNames         []string
}

type LinkMarkdownContext struct {
	UserPath         string
	UserContent      string
	WorkspacePath    string
	WorkspaceContent string
}

func LoadBaseSystemPrompt(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

func BuildSystemPrompt(opts PromptOptions) string {
	surface := strings.TrimSpace(opts.Surface)
	if surface == "" {
		surface = "this conversation"
	}

	basePrompt := strings.TrimSpace(opts.BasePrompt)
	if basePrompt == "" {
		basePrompt = DefaultBaseSystemPrompt
	}

	sections := []string{renderTag("identity", basePrompt)}

	workspaceLines := []string{"Surface: " + surface}
	if workingDirectory := strings.TrimSpace(opts.WorkingDirectory); workingDirectory != "" {
		workspaceLines = append(workspaceLines,
			"Working directory: "+workingDirectory,
			"除非运行时另有说明，默认以工作目录作为运行时文件与 Shell 操作的根目录。",
		)
	}
	sections = append(sections, renderTag("workspace", strings.Join(workspaceLines, "\n")))
	if linkContext := renderLinkMarkdownContext(opts.LinkMarkdown); linkContext != "" {
		sections = append(sections, renderTag("system-reminder", linkContext))
	}

	toolNames := make([]string, 0, len(opts.ToolNames))
	for _, name := range opts.ToolNames {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			continue
		}
		toolNames = append(toolNames, trimmed)
	}
	if len(toolNames) > 0 {
		toolBody := "本次会话可用的工具如下：\n\n" + renderList(toolNames)
		if toolNamesContain(toolNames, "read_persisted_output") {
			toolBody += "\n\n" + persistedOutputGuidance
		}
		sections = append(sections, renderTag("tools", toolBody))
	}

	if descriptions := strings.TrimSpace(opts.SkillDescriptions); descriptions != "" {
		skillBody := strings.Join([]string{
			"以下技能只提供摘要。",
			"重要规则：\n" + renderList([]string{
				"每个技能都有名称和描述。",
				"使用或遵循某个技能前，先用 `load_skill` 以精确的技能名加载其完整说明。",
				"不要仅凭摘要臆测完整的工作流。",
				"若多个技能看起来都相关，先加载最匹配、最具体的那个。",
			}),
			"可用技能：\n\n" + descriptions,
		}, "\n\n")
		sections = append(sections, renderTag("skills", skillBody))
	}

	if memory := strings.TrimSpace(opts.MemorySection); memory != "" {
		sections = append(sections, renderTag("memory", memory))
	}

	return strings.Join(sections, "\n\n")
}

func renderLinkMarkdownContext(ctx LinkMarkdownContext) string {
	userContent := strings.TrimSpace(ctx.UserContent)
	workspaceContent := strings.TrimSpace(ctx.WorkspaceContent)
	if userContent == "" && workspaceContent == "" {
		return ""
	}
	parts := []string{
		"在回答用户问题时，你可以参考以下上下文：",
		"# linkMd",
		"下面展示了用户与代码库说明。请务必遵循这些说明。重要：这些说明将覆盖任何默认行为，你必须严格按其文字要求执行。",
	}
	if userContent != "" {
		parts = append(parts, strings.Join([]string{
			strings.TrimSpace(ctx.UserPath) + " 的内容（用户为所有项目配置的私人全局说明）：",
			userContent,
		}, "\n\n"))
	}
	if workspaceContent != "" {
		parts = append(parts, strings.Join([]string{
			strings.TrimSpace(ctx.WorkspacePath) + " 的内容（项目说明，已提交到代码库或工作区）：",
			workspaceContent,
		}, "\n\n"))
	}
	parts = append(parts, strings.Join([]string{
		"# 重要指令提醒",
		"只做被要求的事：不多不少。",
		"除非为达成目标绝对必要，切勿创建新文件。",
		"能修改现有文件，绝不新建文件。",
		"不要主动创建文档文件（*.md）或README。仅当用户明确要求时才创建文档。",
		"重要：这些上下文可能与当前任务相关，也可能无关。除非与任务高度相关，否则不要对其作出回应。",
	}, "\n"))
	return strings.Join(parts, "\n\n")
}

func renderTag(tag, body string) string {
	tag = strings.TrimSpace(tag)
	body = strings.TrimSpace(body)
	if body == "" {
		return "<" + tag + ">\n</" + tag + ">"
	}
	return "<" + tag + ">\n" + body + "\n</" + tag + ">"
}

func toolNamesContain(names []string, target string) bool {
	for _, name := range names {
		if name == target {
			return true
		}
	}
	return false
}

func renderList(items []string) string {
	lines := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		lines = append(lines, "- "+trimmed)
	}
	return strings.Join(lines, "\n")
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package assistant\n\nimport (\n\t\"os\"\n\t\"strings\"\n)\n\nconst persistedOutputGuidance = \"当较早的消息中出现 `\u003cpersisted-output ...\u003e` 标记时，表示完整的工具输出已存入本地文件，内联的只是预览。\" +\n\t\"如果预览不足以完成任务，请用标记中的 id 和偏移量调用 `read_persisted_output` 分块读取更多内容，不要猜测被省略的部分。\" +\n\t\"当看到 `[Earlier result compacted. Re-run if needed]` 时，请重新执行相关工具以再次获取\ufffd"}
```

## 2026-06-13T09:16:03Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_aRSN9A35XDsbgga48rJv3265

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/sessions/manager.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/sessions/manager.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/sessions/manager.go: no such file or directory"}
```

## 2026-06-13T09:16:03Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_03_baiz9DYdHnb2NnKqQbMp8459

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/main.go"}
```

### Result

```text
package main

import (
	"log"

	"nano_cc/internal/cli"
)

func main() {
	if err := cli.Main(); err != nil {
		log.Fatal(err)
	}
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package main\n\nimport (\n\t\"log\"\n\n\t\"nano_cc/internal/cli\"\n)\n\nfunc main() {\n\tif err := cli.Main(); err != nil {\n\t\tlog.Fatal(err)\n\t}\n}\n"}
```

## 2026-06-13T09:16:03Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_04_GkKL5I5Gmk7KQrhEqj111537

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/handler.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/handler.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/handler.go: no such file or directory"}
```

## 2026-06-13T09:16:06Z · bash · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_5wNba99inG7P7jSphKx49558

### Arguments

```json
{"command": "find /Users/bytedance/golang_pro/nano_cc/go-agent/internal -type f -name \"*.go\" | sort"}
```

### Result

```text
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/transport_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/transport.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/compression_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/conversation_memory_strategy.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/full_history_summarization.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/message_window_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/reactive_compact.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/recent_tool_result_retention.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/token_estimator.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/tool_result_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/context_compression_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/context_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/conversation_flow.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/conversation_memory.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/helpers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hook_bridge.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/manager.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/stop.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/tool.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/types.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/user_prompt.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/memory_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/memory.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/prompt_builder.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/runtime_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/runtime.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/subagent.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/summarizer.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/todo_reminder.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/tool_registry.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/conversation_history_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/conversation_history.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/models.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/assistant/prompt_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/assistant/prompt.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/cli/root_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/cli/root.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/helpers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/local_config_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/local_config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/paths.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/idgen/idgen.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/llm/client.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/llm/errors.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/bootstrap.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/memory_store.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/persisted_output_files.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/session_history_store.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/store_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/store.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/tool_result_log.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/logger/logger_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/logger/logger.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/safety/path.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/sessions/render_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/sessions/skill_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/sessions/skill.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/textutil/textutil.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/bash.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/definitions_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/definitions.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/file_ops.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/handlers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/load_skill.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/path_guard.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/persisted_output_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/persisted_output.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/runtime_env.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/runtime_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/todo_write_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/todo_write.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/validation_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/validation.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/app.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/events_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/events.go
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","resolved_command_path":"/Users/bytedance/golang_pro/nano_cc/go-agent/internal","outcome_summary":"/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config_test.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager_test.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/transport_test.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/transport.go\n/Users/bytedance/golang_pro/nano_"}
```

## 2026-06-13T09:16:08Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_zNDdkjnJLMJYnl9a26TA8694

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/app.go"}
```

### Result

```text
package tui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/runtime"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/logger"
	"nano_cc/internal/sessions"
)

type SessionInfo struct {
	User         storage.User
	Conversation storage.Conversation
	CWD          string
	Resumer      SessionResumer
	Skills       []sessions.SkillSummary
	MCPServers   []mcp.ServerStatus
	SkillCount   int
	MCPToolCount int
}

type SessionResumer interface {
	ListResumableSessions(ctx context.Context, workspaceRoot string) ([]storage.ResumableSession, error)
	ResumeSession(ctx context.Context, sessionID, currentWorkspace string, user storage.User) (storage.Conversation, []storage.Message, error)
}

type Message struct {
	Role             string
	Content          string
	ReasoningContent string
}

type palette struct {
	ink      lipgloss.Color
	muted    lipgloss.Color
	panel    lipgloss.Color
	cyan     lipgloss.Color
	mint     lipgloss.Color
	lavender lipgloss.Color
	butter   lipgloss.Color
	coral    lipgloss.Color
}

var tuiPalette = palette{
	ink:      lipgloss.Color("255"),
	muted:    lipgloss.Color("244"),
	panel:    lipgloss.Color("238"),
	cyan:     lipgloss.Color("81"),
	mint:     lipgloss.Color("120"),
	lavender: lipgloss.Color("183"),
	butter:   lipgloss.Color("229"),
	coral:    lipgloss.Color("210"),
}

const inputCursor = "█"

type Model struct {
	runtime          *runtime.Service
	session          SessionInfo
	messages         []Message
	input            textarea.Model
	viewport         viewport.Model
	width            int
	height           int
	running          bool
	events           chan Event
	cancel           context.CancelFunc
	renderer         *glamour.TermRenderer
	generation       int64
	resumeSelecting  bool
	resumeCandidates []storage.ResumableSession
	toolCallCount    int
	contextTokens    int
	contextBudget    int
}

func NewModel(runtimeService *runtime.Service, session SessionInfo) Model {
	input := textarea.New()
	input.Placeholder = "问 go-agent 一件事..."
	input.Focus()
	input.SetHeight(2)
	input.ShowLineNumbers = false
	vp := viewport.New(100, 20)
	renderer := newMarkdownRenderer(100)
	return Model{runtime: runtimeService, session: session, input: input, viewport: vp, width: 100, height: 20, events: make(chan Event, 128), renderer: renderer}
}

func Run(ctx context.Context, runtimeService *runtime.Service, session SessionInfo) error {
	previousConsole := logger.SetConsoleEnabled(false)
	defer logger.SetConsoleEnabled(previousConsole)
	program := tea.NewProgram(NewModel(runtimeService, session), tea.WithContext(ctx), tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := program.Run()
	return err
}

func (m Model) Init() tea.Cmd { return textarea.Blink }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = m.viewportHeight()
		m.input.SetWidth(max(20, msg.Width-4))
		m.renderer = newMarkdownRenderer(m.messageWidth())
		m.refreshViewport()
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.running && m.cancel != nil {
				m.cancel()
				m.generation++
				m.running = false
				m.appendMessage("system", "已中断当前生成")
				return m, nil
			}
			return m, tea.Quit
		case "enter":
			text := strings.TrimSpace(m.input.Value())
			if text == "" || m.running {
				return m, nil
			}
			m.input.Reset()
			if m.resumeSelecting && m.handleResumeSelection(text) {
				m.refreshViewport()
				return m, nil
			}
			if strings.HasPrefix(text, "/") && m.handleSlashCommand(text) {
				m.refreshViewport()
				return m, nil
			}
			m.appendMessage("user", text)
			m.running = true
			m.toolCallCount = 0
			m.generation++
			generation := m.generation
			ctx, cancel := context.WithCancel(context.Background())
			m.cancel = cancel
			return m, tea.Batch(m.waitEvent(), m.respond(ctx, text, generation))
		}
		if isTerminalProbeResponseInput(msg) {
			return m, nil
		}
		if cmd, ok := m.updateViewportScroll(msg); ok {
			return m, cmd
		}
	case tea.MouseMsg:
		if cmd, ok := m.updateViewportScroll(msg); ok {
			return m, cmd
		}
	case Event:
		if msg.Generation != 0 && msg.Generation != m.generation {
			if m.running {
				return m, m.waitEvent()
			}
			return m, nil
		}
		switch msg.Name {
		case "assistant_delta":
			m.appendAssistantDelta(msg.Content)
		case "reasoning_delta":
			m.appendThinkingDelta(msg.Content)
		case "assistant":
			m.updateMetaFromData(msg.Data)
			content := msg.Content
			if content == "" && msg.Data != nil {
				content = eventContent(msg.Data)
			}
			if content != "" {
				m.replaceLastAssistant(content, eventString(msg.Data, "reasoning_content"))
			}
		case "meta":
			m.updateMetaFromData(msg.Data)
		case "error":
			m.appendMessage("error", msg.Content)
			m.running = false
		case "done":
			m.running = false
		}
		m.refreshViewport()
		if m.running {
			return m, m.waitEvent()
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.refreshViewport()
	return m, cmd
}

func (m *Model) updateViewportScroll(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyPgUp, tea.KeyPgDown, tea.KeyUp, tea.KeyDown:
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return cmd, true
		}
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress {
			switch msg.Button {
			case tea.MouseButtonWheelUp, tea.MouseButtonWheelDown, tea.MouseButtonWheelLeft, tea.MouseButtonWheelRight:
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				return cmd, true
			}
		}
	}
	return nil, false
}

func (m Model) View() string {
	return m.viewport.View()
}

func (m Model) respond(ctx context.Context, text string, generation int64) tea.Cmd {
	return func() tea.Msg {
		if m.runtime == nil {
			m.events <- Event{Generation: generation, Name: "error", Content: "runtime 未初始化"}
			return nil
		}
		_, err := m.runtime.RespondToConversation(ctx, m.session.Conversation, m.session.User, text, NewEventWriter(m.events, generation))
		if err != nil {
			m.events <- Event{Generation: generation, Name: "error", Content: err.Error()}
			return nil
		}
		m.events <- Event{Generation: generation, Name: "done"}
		return nil
	}
}

func (m Model) waitEvent() tea.Cmd {
	return func() tea.Msg { return <-m.events }
}

func (m *Model) handleSlashCommand(text string) bool {
	switch strings.TrimSpace(text) {
	case "/help":
		m.appendMessage("system", "命令：/help /clear /cwd /skills /mcp /resume。Enter 发送，Ctrl+C 中断或退出。")
		return true
	case "/clear":
		m.messages = nil
		m.resumeSelecting = false
		m.resumeCandidates = nil
		m.appendMessage("system", "已清空当前 TUI 显示上下文")
		return true
	case "/cwd":
		m.appendMessage("system", "当前工作区："+m.session.CWD)
		return true
	case "/skills":
		m.appendMessage("system", renderSkillDetails(m.session.Skills, m.session.SkillCount))
		return true
	case "/mcp":
		m.appendMessage("system", renderMCPDetails(m.session.MCPServers, m.session.MCPToolCount))
		return true
	case "/resume":
		m.startResumeSelection()
		return true
	}
	m.appendMessage("system", "未知命令："+text)
	return true
}

func (m *Model) startResumeSelection() {
	if m.running {
		m.appendMessage("system", "当前正在生成，请先 Ctrl+C 中断后再执行 /resume")
		return
	}
	if m.session.Resumer == nil {
		m.appendMessage("system", "当前运行环境不支持 /resume")
		return
	}
	sessions, err := m.session.Resumer.ListResumableSessions(context.Background(), m.session.CWD)
	if err != nil {
		m.appendMessage("error", err.Error())
		return
	}
	if len(sessions) == 0 {
		m.resumeSelecting = false
		m.resumeCandidates = nil
		m.appendMessage("system", "当前目录暂无可恢复的历史会话")
		return
	}
	m.resumeSelecting = true
	m.resumeCandidates = sessions
	m.appendMessage("system", renderResumeCandidates(sessions))
}

func (m *Model) handleResumeSelection(text string) bool {
	text = strings.TrimSpace(text)
	if text == "/cancel" || text == "/clear" {
		m.resumeSelecting = false
		m.resumeCandidates = nil
		if text == "/clear" {
			m.messages = nil
			m.appendMessage("system", "已清空当前 TUI 显示上下文")
		} else {
			m.appendMessage("system", "已取消恢复历史会话")
		}
		return true
	}
	idx, err := strconv.Atoi(text)
	if err != nil || idx < 1 || idx > len(m.resumeCandidates) {
		m.appendMessage("system", fmt.Sprintf("请输入 1-%d 之间的序号，或输入 /cancel 取消", len(m.resumeCandidates)))
		return true
	}
	candidate := m.resumeCandidates[idx-1]
	conv, history, err := m.session.Resumer.ResumeSession(context.Background(), candidate.SessionID, m.session.CWD, m.session.User)
	if err != nil {
		m.appendMessage("error", err.Error())
		return true
	}
	m.session.Conversation = conv
	m.resumeSelecting = false
	m.resumeCandidates = nil
	m.messages = messagesForDisplay(history)
	m.appendMessage("system", fmt.Sprintf("已恢复历史会话：%s", conv.SessionID))
	return true
}

func renderResumeCandidates(sessions []storage.ResumableSession) string {
	var b strings.Builder
	b.WriteString("可恢复的历史会话：")
	for i, session := range sessions {
		title := strings.TrimSpace(session.Title)
		if title == "" {
			title = "TUI 会话"
		}
		updated := "unknown"
		if !session.UpdatedAt.IsZero() {
			updated = session.UpdatedAt.Local().Format(time.RFC3339)
		}
		b.WriteString(fmt.Sprintf("\n%d. %s | %s | 消息:%d | %s", i+1, updated, title, session.MessageCount, session.SessionID))
	}
	b.WriteString("\n输入序号恢复，或输入 /cancel 取消。")
	return b.String()
}

func messagesForDisplay(history []storage.Message) []Message {
	if len(history) == 0 {
		return nil
	}
	messages := make([]Message, 0, len(history))
	for _, msg := range history {
		switch msg.Role {
		case "user", "assistant", "system", "error":
			messages = append(messages, Message{Role: msg.Role, Content: msg.Content, ReasoningContent: msg.ReasoningContent})
		}
	}
	return messages
}

func renderSkillDetails(skills []sessions.SkillSummary, fallbackCount int) string {
	count := len(skills)
	if count == 0 && fallbackCount > 0 {
		count = fallbackCount
	}
	if len(skills) == 0 {
		return fmt.Sprintf("已加载 Skills：%d 个", count)
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("已加载 Skills：%d 个", len(skills)))
	for _, skill := range skills {
		b.WriteString("\n- ")
		b.WriteString(skill.Name)
		if strings.TrimSpace(skill.Source) != "" {
			b.WriteString(" [")
			b.WriteString(skill.Source)
			b.WriteString("]")
		}
		if strings.TrimSpace(skill.Description) != "" {
			b.WriteString(" ")
			b.WriteString(skill.Description)
		}
		if strings.TrimSpace(skill.Path) != "" {
			b.WriteString("\n  path: ")
			b.WriteString(skill.Path)
		}
	}
	return b.String()
}

func renderMCPDetails(servers []mcp.ServerStatus, toolCount int) string {
	if len(servers) == 0 {
		return fmt.Sprintf("MCP Servers：0 个，工具：%d 个", toolCount)
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("MCP Servers：%d 个，工具：%d 个", len(servers), toolCount))
	for _, server := range servers {
		state := "disabled"
		if server.Connected {
			state = "connected"
		} else if server.LastError != "" {
			state = "failed"
		} else if server.Enabled {
			state = "pending"
		}
		b.WriteString("\n- ")
		b.WriteString(server.Name)
		b.WriteString(" [")
		b.WriteString(server.Transport)
		b.WriteString("] ")
		b.WriteString(state)
		b.WriteString(fmt.Sprintf(", tools: %d", server.ToolCount))
		if server.Command != "" {
			b.WriteString("\n  command: ")
			b.WriteString(strings.Join(append([]string{server.Command}, server.Args...), " "))
		}
		if server.URL != "" {
			b.WriteString("\n  url: ")
			b.WriteString(server.URL)
		}
		if server.LastError != "" {
			b.WriteString("\n  error: ")
			b.WriteString(server.LastError)
		}
	}
	return b.String()
}

func (m *Model) appendMessage(role, content string) {
	m.messages = append(m.messages, Message{Role: role, Content: content})
}

func (m *Model) appendAssistantDelta(delta string) {
	if len(m.messages) == 0 || !isLiveAssistantRole(m.messages[len(m.messages)-1].Role) {
		m.appendMessage("assistant", delta)
		return
	}
	if m.messages[len(m.messages)-1].Role == "thinking" {
		m.messages[len(m.messages)-1].Role = "assistant"
	}
	m.messages[len(m.messages)-1].Content += delta
}

func (m *Model) appendThinkingDelta(delta string) {
	if delta == "" {
		return
	}
	if len(m.messages) == 0 || !isLiveAssistantRole(m.messages[len(m.messages)-1].Role) {
		m.messages = append(m.messages, Message{Role: "assistant", ReasoningContent: delta})
		return
	}
	m.messages[len(m.messages)-1].ReasoningContent += delta
}

func (m *Model) replaceLastAssistant(content, reasoning string) {
	if len(m.messages) == 0 || !isLiveAssistantRole(m.messages[len(m.messages)-1].Role) {
		m.messages = append(m.messages, Message{Role: "assistant", Content: content, ReasoningContent: reasoning})
		return
	}
	m.messages[len(m.messages)-1].Role = "assistant"
	m.messages[len(m.messages)-1].Content = content
	m.messages[len(m.messages)-1].ReasoningContent = reasoning
}

func isLiveAssistantRole(role string) bool {
	return role == "assistant" || role == "thinking"
}

func (m *Model) refreshViewport() {
	m.viewport.Width = max(10, m.width)
	m.viewport.Height = m.viewportHeight()
	m.viewport.SetContent(m.renderTranscript())
	m.viewport.GotoBottom()
}

func (m Model) viewportHeight() int {
	if m.height <= 0 {
		return 20
	}
	return max(1, m.height)
}

func (m Model) renderTranscript() string {
	var b strings.Builder
	b.WriteString(m.renderWelcome())
	if len(m.messages) > 0 {
		b.WriteString("\n\n")
		b.WriteString(m.renderMessages())
	}
	b.WriteString("\n")
	b.WriteString(m.renderInput())
	b.WriteString("\n")
	b.WriteString(subtleStyle().Render(m.renderLiveStatus()))
	return b.String()
}

func (m Model) renderMessages() string {
	var b strings.Builder
	for _, msg := range m.messages {
		b.WriteString(m.renderMessage(msg))
		b.WriteString("\n\n")
	}
	return b.String()
}

func (m Model) renderMessage(msg Message) string {
	switch msg.Role {
	case "user":
		return promptLineStyle().Render("›") + " " + userStyle().Render(wrapText(msg.Content, m.messageWidth()-2))
	case "assistant":
		content := msg.Content
		if m.renderer != nil {
			if rendered, err := m.renderer.Render(content); err == nil {
				content = wrapText(strings.TrimSpace(rendered), m.messageWidth())
			}
		} else {
			content = wrapText(content, m.messageWidth())
		}
		if m.running && strings.TrimSpace(msg.ReasoningContent) != "" {
			content = thinkingStyle().Render("✽ 思考中\n"+wrapText(strings.TrimSpace(msg.ReasoningContent), m.messageWidth()-2)) + "\n" + content
		}
		return assistantLeadStyle().Render("go-agent") + "\n" + content
	case "thinking":
		return thinkingStyle().Render("✽ 思考中\n" + wrapText(msg.Content, m.messageWidth()-2))
	case "system":
		return systemStyle().Render("• " + wrapText(msg.Content, m.messageWidth()-2))
	case "error":
		return errorStyle().Render("✗ " + wrapText(msg.Content, m.messageWidth()-2))
	default:
		return roleLabel(msg.Role, lipgloss.Color("245")) + "\n" + wrapText(msg.Content, m.messageWidth())
	}
}

func (m Model) messageWidth() int {
	return max(10, m.width)
}

func newMarkdownRenderer(width int) *glamour.TermRenderer {
	renderer, _ := glamour.NewTermRenderer(glamour.WithStandardStyle("dark"), glamour.WithWordWrap(max(10, width)))
	return renderer
}

func wrapText(text string, width int) string {
	width = max(1, width)
	return ansi.Hardwrap(text, width, true)
}

func (m Model) renderWelcome() string {
	art := strings.Join([]string{
		`   /\_/\\`,
		`  ( o.o )   nano, but cozy`,
		`   > ^ <    ask · think · build`,
	}, "\n")
	intro := "像 Claude Code 一样一问一答：在终端直接提问，回答会流式显示在下方。"
	quick := "快捷键：Enter 发送 · Ctrl+C 中断/退出 · /resume 恢复 · /clear 清屏"
	stats := fmt.Sprintf("Skills %d · MCP tools %d", m.session.SkillCount, m.session.MCPToolCount)
	return startupPanelStyle().Width(max(20, m.width-2)).Render(accentArtStyle().Render(art) + "\n\n" + intro + "\n" + subtleStyle().Render(quick) + "\n" + subtleStyle().Render(stats))
}

func (m Model) renderHeader() string {
	width := max(20, m.width)
	status := runningText(m.running)
	contextText := "上下文 --"
	if m.contextBudget > 0 {
		contextText = fmt.Sprintf("上下文 %d%% · %s/%s", min(100, m.contextTokens*100/m.contextBudget), compactNumber(m.contextTokens), compactNumber(m.contextBudget))
	}
	left := titleStyle().Render("✦ go-agent")
	right := headerMetaStyle().Render(fmt.Sprintf("%s  ·  本轮工具 %d  ·  %s", status, m.toolCallCount, contextText))
	gap := max(1, width-lipgloss.Width(left)-lipgloss.Width(right)-4)
	line := left + strings.Repeat(" ", gap) + right
	workspace := subtleStyle().Render(fmt.Sprintf("cwd %s · skills %d · mcp tools %d", m.session.CWD, m.session.SkillCount, m.session.MCPToolCount))
	return headerStyle().Width(width).Render(line + "\n" + workspace)
}

func (m Model) renderConversationFrame() string {
	return conversationStyle().Width(max(10, m.width)).Render(m.viewport.View())
}

func (m Model) renderInput() string {
	prompt := inputPromptStyle().Render("›")
	text := strings.TrimSpace(m.input.Value())
	if text == "" {
		text = inputPromptStyle().Render(inputCursor) + " " + subtleStyle().Render(m.input.Placeholder)
	} else {
		text = userStyle().Render(text) + inputPromptStyle().Render(inputCursor)
	}
	return inputLineStyle().Width(max(10, m.width)).Render(prompt + " " + text)
}

func isTerminalProbeResponseInput(msg tea.KeyMsg) bool {
	if msg.Type != tea.KeyRunes {
		return false
	}
	value := string(msg.Runes)
	return strings.Contains(value, ";rgb:") && strings.Contains(value, "/")
}

func (m Model) renderLiveStatus() string {
	parts := []string{"Enter 发送", "Ctrl+C 中断/退出", "/help", fmt.Sprintf("工具 %d", m.toolCallCount)}
	if m.contextBudget > 0 {
		parts = append(parts, fmt.Sprintf("上下文 %d%% · %s/%s", min(100, m.contextTokens*100/m.contextBudget), compactNumber(m.contextTokens), compactNumber(m.contextBudget)))
	} else {
		parts = append(parts, "上下文 --")
	}
	return strings.Join(parts, " · ")
}

func (m *Model) updateMetaFromData(data any) {
	if count, ok := eventInt(data, "tool_call_count"); ok {
		m.toolCallCount = count
	}
	if tokens, ok := eventInt(data, "context_tokens"); ok {
		m.contextTokens = tokens
	}
	if budget, ok := eventInt(data, "context_budget"); ok {
		m.contextBudget = budget
	}
}

func eventString(data any, key string) string {
	m, ok := data.(map[string]any)
	if !ok {
		return ""
	}
	value, _ := m[key].(string)
	return value
}

func eventInt(data any, key string) (int, bool) {
	m, ok := data.(map[string]any)
	if !ok {
		return 0, false
	}
	switch value := m[key].(type) {
	case int:
		return value, true
	case int64:
		return int(value), true
	case float64:
		return int(value), true
	default:
		return 0, false
	}
}

func roleLabel(label string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Bold(true).Foreground(color).Render(label)
}

func headerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(tuiPalette.panel)
}

func titleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(tuiPalette.mint)
}

func headerMetaStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
}

func subtleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.muted)
}

func conversationStyle() lipgloss.Style {
	return lipgloss.NewStyle().Padding(1, 2, 0, 2)
}

func startupPanelStyle() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(tuiPalette.coral).Padding(1, 2).Margin(1, 0)
}

func userStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.ink)
}

func thinkingStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Italic(true).PaddingLeft(2)
}

func systemStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.butter)
}

func errorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.coral)
}

func promptLineStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.cyan).Bold(true)
}

func assistantLeadStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.mint).Bold(true)
}

func accentArtStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.lavender).Bold(true)
}

func inputPromptStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.mint).Bold(true)
}

func inputLineStyle() lipgloss.Style {
	return lipgloss.NewStyle().Background(lipgloss.Color("238")).Foreground(tuiPalette.ink)
}

func compactNumber(n int) string {
	if n >= 1000 {
		return fmt.Sprintf("%dk", n/1000)
	}
	return strconv.Itoa(n)
}

func runningText(running bool) string {
	if running {
		return "generating"
	}
	return "ready"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package tui\n\nimport (\n\t\"context\"\n\t\"fmt\"\n\t\"strconv\"\n\t\"strings\"\n\t\"time\"\n\n\t\"github.com/charmbracelet/bubbles/textarea\"\n\t\"github.com/charmbracelet/bubbles/viewport\"\n\ttea \"github.com/charmbracelet/bubbletea\"\n\t\"github.com/charmbracelet/glamour\"\n\t\"github.com/charmbracelet/lipgloss\"\n\t\"github.com/charmbracelet/x/ansi\"\n\n\t\"nano_cc/internal/agent/mcp\"\n\t\"nano_cc/internal/agent/runtime\"\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/logger\"\n\t\"nano_cc/internal/sessions\"\n)\n\ntype SessionInfo struct {\n\tUser"}
```

## 2026-06-13T09:16:08Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_SJKXTzuuXELugCErUmwh5331

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/events.go"}
```

### Result

```text
package tui

import "fmt"

type Event struct {
	Generation int64
	Name       string
	Content    string
	Data       any
}

type EventWriter struct {
	ch         chan<- Event
	generation int64
}

func NewEventWriter(ch chan<- Event, generation ...int64) EventWriter {
	var gen int64
	if len(generation) > 0 {
		gen = generation[0]
	}
	return EventWriter{ch: ch, generation: gen}
}

func (w EventWriter) Event(name string, data any) error {
	if w.ch == nil {
		return nil
	}
	w.ch <- Event{Generation: w.generation, Name: name, Content: eventContent(data), Data: data}
	return nil
}

func eventContent(data any) string {
	if m, ok := data.(map[string]any); ok {
		if value, ok := m["content"].(string); ok {
			return value
		}
		if value, ok := m["message"].(string); ok {
			return value
		}
	}
	return fmt.Sprint(data)
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package tui\n\nimport \"fmt\"\n\ntype Event struct {\n\tGeneration int64\n\tName       string\n\tContent    string\n\tData       any\n}\n\ntype EventWriter struct {\n\tch         chan\u003c- Event\n\tgeneration int64\n}\n\nfunc NewEventWriter(ch chan\u003c- Event, generation ...int64) EventWriter {\n\tvar gen int64\n\tif len(generation) \u003e 0 {\n\t\tgen = generation[0]\n\t}\n\treturn EventWriter{ch: ch, generation: gen}\n}\n\nfunc (w EventWriter) Event(name string, data any) error {\n\tif w.ch == nil {\n\t\treturn nil\n\t}\n\tw.ch \u003c- Event{Generation: w"}
```

## 2026-06-13T09:16:08Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_ew96BGxTeLKC8csYBtN28110

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/runtime.go"}
```

### Result

```text
package runtime

import (
	"context"
	"time"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/runtime/compression"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/config"
	"nano_cc/internal/llm"
	"nano_cc/internal/sessions"
)

type EventWriter interface {
	Event(name string, data any) error
}

type conversationStore interface {
	UpdateConversationTitle(ctx context.Context, conversationID, title string) error
	TouchConversationActivity(ctx context.Context, conversationID string) error
	SetConversationHistory(ctx context.Context, conversationID string, messages []storage.Message) error
	SetConversationCache(ctx context.Context, conversationID string, messages []storage.Message) error
	GetConversationCache(ctx context.Context, conversationID string) ([]storage.Message, bool, error)
	ListMessagesByConversation(ctx context.Context, conversationID string, limit int) ([]storage.Message, error)
	CreateToolCall(ctx context.Context, tc storage.ToolCall) error
	CreateSubagentMessage(ctx context.Context, message storage.SubagentMessage) error
	CreatePersistedOutput(ctx context.Context, output storage.PersistedOutput) error
	GetPersistedOutputForConversation(ctx context.Context, id, userID, conversationID string) (storage.PersistedOutput, error)
	GetPersistedOutputByMessageHash(ctx context.Context, conversationID, userID, messageID, toolCallID, strategy, contentSHA256 string) (storage.PersistedOutput, error)
	CreateContextSummary(ctx context.Context, summary storage.ContextSummary) error
	GetContextSummaryByHistoryHash(ctx context.Context, conversationID, userID, sourceHistorySHA256 string) (storage.ContextSummary, error)
	ListRelevantMemories(ctx context.Context, userID string) ([]storage.Memory, error)
	ListMemoriesByUserAndType(ctx context.Context, userID, memType string) ([]storage.Memory, error)
	ListProjectFactMemories(ctx context.Context, userID string) ([]storage.Memory, error)
	InsertMemory(ctx context.Context, m storage.Memory) error
	CountMemoriesByUserAndType(ctx context.Context, userID, memType string) (int, error)
	CountProjectFactMemories(ctx context.Context, userID string) (int, error)
	DeleteOldestMemories(ctx context.Context, userID, memType string, n int) error
	ReplaceMemoriesByUserAndType(ctx context.Context, userID, memType string, items []storage.Memory) error
	ReplaceProjectFactMemories(ctx context.Context, userID string, items []storage.Memory) error
	ListConversationMemories(ctx context.Context, conversationID string) ([]storage.ConversationMemory, error)
	ReplaceConversationMemories(ctx context.Context, conversationID, userID string, items []storage.ConversationMemory) error
	GetConversationModelHistory(ctx context.Context, conversationID string) ([]storage.Message, bool, error)
	UpsertConversationModelHistory(ctx context.Context, conversationID, userID string, messages []storage.Message) error
	AcquireConversationLock(ctx context.Context, conversationID, token string, ttl, waitTimeout time.Duration) (bool, error)
	RenewConversationLock(ctx context.Context, conversationID, token string, ttl time.Duration) (bool, error)
	ReleaseConversationLock(ctx context.Context, conversationID, token string) error
}

type Service struct {
	Store             conversationStore
	Cfg               config.AppConfig
	LLM               llm.Client
	Tools             *ToolRegistry
	BuiltinSkills     *sessions.SkillLoader
	BasePrompt        string
	LinkMarkdown      config.LinkMarkdownContext
	Hooks             *HookManager
	ContextCompressor *compression.Compressor
	EnableMemory      bool
	MCP               *mcp.Manager
}

func NewService(store conversationStore, cfg config.AppConfig, client llm.Client) *Service {
	return &Service{Store: store, Cfg: cfg, LLM: client, Tools: NewToolRegistry(cfg), Hooks: NewDefaultHookManager(), EnableMemory: true}
}

func (s *Service) hookManager() *HookManager {
	if s.Hooks == nil {
		s.Hooks = NewDefaultHookManager()
	}
	return s.Hooks
}

func (s *Service) SetBuiltinSkills(loader *sessions.SkillLoader) {
	s.BuiltinSkills = loader
}

func (s *Service) SetBasePrompt(prompt string) {
	s.BasePrompt = prompt
}

func (s *Service) SetLinkMarkdownContext(ctx config.LinkMarkdownContext) {
	s.LinkMarkdown = ctx
}

func (s *Service) SetMCPManager(manager *mcp.Manager) {
	s.MCP = manager
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"context\"\n\t\"time\"\n\n\t\"nano_cc/internal/agent/mcp\"\n\t\"nano_cc/internal/agent/runtime/compression\"\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/config\"\n\t\"nano_cc/internal/llm\"\n\t\"nano_cc/internal/sessions\"\n)\n\ntype EventWriter interface {\n\tEvent(name string, data any) error\n}\n\ntype conversationStore interface {\n\tUpdateConversationTitle(ctx context.Context, conversationID, title string) error\n\tTouchConversationActivity(ctx context.Context, conversationID string) error"}
```

## 2026-06-13T09:16:08Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_03_AB3aeCgJG48mk8dDZRyK5121

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/bootstrap.go"}
```

### Result

```text
package local

import (
	"context"
	"fmt"
	"time"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/runtime"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/assistant"
	"nano_cc/internal/config"
	"nano_cc/internal/idgen"
	"nano_cc/internal/llm"
	"nano_cc/internal/logger"
	"nano_cc/internal/sessions"
)

const LocalUserID = "local-user"

type Bundle struct {
	Runtime      *runtime.Service
	Store        *Store
	MCP          *mcp.Manager
	User         storage.User
	Conversation storage.Conversation
	CWD          string
	SkillCount   int
	Skills       []sessions.SkillSummary
	MCPToolCount int
	MCPServers   []mcp.ServerStatus
}

func Bootstrap(ctx context.Context, cwd string) (*Bundle, error) {
	cfg, err := config.LoadLocalConfig(cwd)
	if err != nil {
		return nil, err
	}
	if err := config.EnsureAppLayout(cfg); err != nil {
		return nil, err
	}
	if err := config.ValidateAppLayout(cfg); err != nil {
		return nil, err
	}
	if err := logger.InitFileLoggerAt(cfg.LogsDir); err != nil {
		logger.Warn(fmt.Sprintf("failed to init file logger: %v", err))
	}
	userSkillsDir, err := config.LinkSkillsDir()
	if err != nil {
		return nil, err
	}
	builtinSkills, err := sessions.LoadSkillsFromDirs([]sessions.SkillDir{
		{Path: userSkillsDir, Source: "user"},
		{Path: config.WorkspaceLinkSkillsDir(cfg.WorkspaceRoot), Source: "workspace"},
	})
	if err != nil {
		return nil, fmt.Errorf("load .link skills: %w", err)
	}
	basePrompt, err := assistant.LoadBaseSystemPrompt(cfg.SystemPromptPath)
	if err != nil {
		return nil, fmt.Errorf("load system prompt: %w", err)
	}
	linkMarkdown, err := config.LoadLinkMarkdownContext(cfg.WorkspaceRoot)
	if err != nil {
		return nil, fmt.Errorf("load LINK.MD context: %w", err)
	}
	store, err := NewStoreWithMemory(cfg.WorkspaceRoot)
	if err != nil {
		return nil, fmt.Errorf("init memory store: %w", err)
	}
	client := llm.NewDeepseekClient(cfg.LLM.BaseURL, cfg.LLM.APIKey)
	runtimeService := runtime.NewService(store, cfg, client)
	runtimeService.EnableMemory = true
	runtimeService.SetBuiltinSkills(builtinSkills)
	runtimeService.SetBasePrompt(basePrompt)
	runtimeService.SetLinkMarkdownContext(linkMarkdown)
	mcpManager := mcp.NewManager()
	workspaceMCPServers, err := mcp.LoadWorkspaceConfig(config.WorkspaceMCPConfigPath(cfg.WorkspaceRoot))
	if err != nil {
		mcpManager.Close()
		return nil, fmt.Errorf("load workspace mcp config: %w", err)
	}
	mcpManager.SetWorkspaceServers(ctx, workspaceMCPServers)
	runtimeService.SetMCPManager(mcpManager)
	user := storage.User{ID: LocalUserID, Email: "local@go-agent", Username: "local", MemoryEnabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	conversation := storage.Conversation{ID: idgen.New("conv"), SessionID: idgen.UUID(), UserID: user.ID, RootMessageID: idgen.New("msg"), Title: "TUI 会话", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := store.CreateConversation(ctx, conversation); err != nil {
		mcpManager.Close()
		return nil, err
	}
	mcpManager.EnsureWorkspaceSessions(ctx)
	mcpTools := mcpManager.ToolsForUser(user.ID)
	mcpSnapshot := mcpManager.Snapshot(user.ID)
	skills := builtinSkills.Summaries()
	return &Bundle{Runtime: runtimeService, Store: store, MCP: mcpManager, User: user, Conversation: conversation, CWD: cfg.WorkspaceRoot, SkillCount: len(skills), Skills: skills, MCPToolCount: len(mcpTools), MCPServers: mcpSnapshot.Servers}, nil
}

func (b *Bundle) Close() {
	if b != nil && b.MCP != nil {
		b.MCP.Close()
	}
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package local\n\nimport (\n\t\"context\"\n\t\"fmt\"\n\t\"time\"\n\n\t\"nano_cc/internal/agent/mcp\"\n\t\"nano_cc/internal/agent/runtime\"\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/assistant\"\n\t\"nano_cc/internal/config\"\n\t\"nano_cc/internal/idgen\"\n\t\"nano_cc/internal/llm\"\n\t\"nano_cc/internal/logger\"\n\t\"nano_cc/internal/sessions\"\n)\n\nconst LocalUserID = \"local-user\"\n\ntype Bundle struct {\n\tRuntime      *runtime.Service\n\tStore        *Store\n\tMCP          *mcp.Manager\n\tUser         storage.User\n\tConversation storage.Co"}
```

## 2026-06-13T09:16:11Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_nYXbSOWy9D9KTw86J0cj7899

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/app.go"}
```

### Result

```text
package tui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/runtime"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/logger"
	"nano_cc/internal/sessions"
)

type SessionInfo struct {
	User         storage.User
	Conversation storage.Conversation
	CWD          string
	Resumer      SessionResumer
	Skills       []sessions.SkillSummary
	MCPServers   []mcp.ServerStatus
	SkillCount   int
	MCPToolCount int
}

type SessionResumer interface {
	ListResumableSessions(ctx context.Context, workspaceRoot string) ([]storage.ResumableSession, error)
	ResumeSession(ctx context.Context, sessionID, currentWorkspace string, user storage.User) (storage.Conversation, []storage.Message, error)
}

type Message struct {
	Role             string
	Content          string
	ReasoningContent string
}

type palette struct {
	ink      lipgloss.Color
	muted    lipgloss.Color
	panel    lipgloss.Color
	cyan     lipgloss.Color
	mint     lipgloss.Color
	lavender lipgloss.Color
	butter   lipgloss.Color
	coral    lipgloss.Color
}

var tuiPalette = palette{
	ink:      lipgloss.Color("255"),
	muted:    lipgloss.Color("244"),
	panel:    lipgloss.Color("238"),
	cyan:     lipgloss.Color("81"),
	mint:     lipgloss.Color("120"),
	lavender: lipgloss.Color("183"),
	butter:   lipgloss.Color("229"),
	coral:    lipgloss.Color("210"),
}

const inputCursor = "█"

type Model struct {
	runtime          *runtime.Service
	session          SessionInfo
	messages         []Message
	input            textarea.Model
	viewport         viewport.Model
	width            int
	height           int
	running          bool
	events           chan Event
	cancel           context.CancelFunc
	renderer         *glamour.TermRenderer
	generation       int64
	resumeSelecting  bool
	resumeCandidates []storage.ResumableSession
	toolCallCount    int
	contextTokens    int
	contextBudget    int
}

func NewModel(runtimeService *runtime.Service, session SessionInfo) Model {
	input := textarea.New()
	input.Placeholder = "问 go-agent 一件事..."
	input.Focus()
	input.SetHeight(2)
	input.ShowLineNumbers = false
	vp := viewport.New(100, 20)
	renderer := newMarkdownRenderer(100)
	return Model{runtime: runtimeService, session: session, input: input, viewport: vp, width: 100, height: 20, events: make(chan Event, 128), renderer: renderer}
}

func Run(ctx context.Context, runtimeService *runtime.Service, session SessionInfo) error {
	previousConsole := logger.SetConsoleEnabled(false)
	defer logger.SetConsoleEnabled(previousConsole)
	program := tea.NewProgram(NewModel(runtimeService, session), tea.WithContext(ctx), tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := program.Run()
	return err
}

func (m Model) Init() tea.Cmd { return textarea.Blink }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = m.viewportHeight()
		m.input.SetWidth(max(20, msg.Width-4))
		m.renderer = newMarkdownRenderer(m.messageWidth())
		m.refreshViewport()
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.running && m.cancel != nil {
				m.cancel()
				m.generation++
				m.running = false
				m.appendMessage("system", "已中断当前生成")
				return m, nil
			}
			return m, tea.Quit
		case "enter":
			text := strings.TrimSpace(m.input.Value())
			if text == "" || m.running {
				return m, nil
			}
			m.input.Reset()
			if m.resumeSelecting && m.handleResumeSelection(text) {
				m.refreshViewport()
				return m, nil
			}
			if strings.HasPrefix(text, "/") && m.handleSlashCommand(text) {
				m.refreshViewport()
				return m, nil
			}
			m.appendMessage("user", text)
			m.running = true
			m.toolCallCount = 0
			m.generation++
			generation := m.generation
			ctx, cancel := context.WithCancel(context.Background())
			m.cancel = cancel
			return m, tea.Batch(m.waitEvent(), m.respond(ctx, text, generation))
		}
		if isTerminalProbeResponseInput(msg) {
			return m, nil
		}
		if cmd, ok := m.updateViewportScroll(msg); ok {
			return m, cmd
		}
	case tea.MouseMsg:
		if cmd, ok := m.updateViewportScroll(msg); ok {
			return m, cmd
		}
	case Event:
		if msg.Generation != 0 && msg.Generation != m.generation {
			if m.running {
				return m, m.waitEvent()
			}
			return m, nil
		}
		switch msg.Name {
		case "assistant_delta":
			m.appendAssistantDelta(msg.Content)
		case "reasoning_delta":
			m.appendThinkingDelta(msg.Content)
		case "assistant":
			m.updateMetaFromData(msg.Data)
			content := msg.Content
			if content == "" && msg.Data != nil {
				content = eventContent(msg.Data)
			}
			if content != "" {
				m.replaceLastAssistant(content, eventString(msg.Data, "reasoning_content"))
			}
		case "meta":
			m.updateMetaFromData(msg.Data)
		case "error":
			m.appendMessage("error", msg.Content)
			m.running = false
		case "done":
			m.running = false
		}
		m.refreshViewport()
		if m.running {
			return m, m.waitEvent()
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.refreshViewport()
	return m, cmd
}

func (m *Model) updateViewportScroll(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyPgUp, tea.KeyPgDown, tea.KeyUp, tea.KeyDown:
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return cmd, true
		}
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress {
			switch msg.Button {
			case tea.MouseButtonWheelUp, tea.MouseButtonWheelDown, tea.MouseButtonWheelLeft, tea.MouseButtonWheelRight:
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				return cmd, true
			}
		}
	}
	return nil, false
}

func (m Model) View() string {
	return m.viewport.View()
}

func (m Model) respond(ctx context.Context, text string, generation int64) tea.Cmd {
	return func() tea.Msg {
		if m.runtime == nil {
			m.events <- Event{Generation: generation, Name: "error", Content: "runtime 未初始化"}
			return nil
		}
		_, err := m.runtime.RespondToConversation(ctx, m.session.Conversation, m.session.User, text, NewEventWriter(m.events, generation))
		if err != nil {
			m.events <- Event{Generation: generation, Name: "error", Content: err.Error()}
			return nil
		}
		m.events <- Event{Generation: generation, Name: "done"}
		return nil
	}
}

func (m Model) waitEvent() tea.Cmd {
	return func() tea.Msg { return <-m.events }
}

func (m *Model) handleSlashCommand(text string) bool {
	switch strings.TrimSpace(text) {
	case "/help":
		m.appendMessage("system", "命令：/help /clear /cwd /skills /mcp /resume。Enter 发送，Ctrl+C 中断或退出。")
		return true
	case "/clear":
		m.messages = nil
		m.resumeSelecting = false
		m.resumeCandidates = nil
		m.appendMessage("system", "已清空当前 TUI 显示上下文")
		return true
	case "/cwd":
		m.appendMessage("system", "当前工作区："+m.session.CWD)
		return true
	case "/skills":
		m.appendMessage("system", renderSkillDetails(m.session.Skills, m.session.SkillCount))
		return true
	case "/mcp":
		m.appendMessage("system", renderMCPDetails(m.session.MCPServers, m.session.MCPToolCount))
		return true
	case "/resume":
		m.startResumeSelection()
		return true
	}
	m.appendMessage("system", "未知命令："+text)
	return true
}

func (m *Model) startResumeSelection() {
	if m.running {
		m.appendMessage("system", "当前正在生成，请先 Ctrl+C 中断后再执行 /resume")
		return
	}
	if m.session.Resumer == nil {
		m.appendMessage("system", "当前运行环境不支持 /resume")
		return
	}
	sessions, err := m.session.Resumer.ListResumableSessions(context.Background(), m.session.CWD)
	if err != nil {
		m.appendMessage("error", err.Error())
		return
	}
	if len(sessions) == 0 {
		m.resumeSelecting = false
		m.resumeCandidates = nil
		m.appendMessage("system", "当前目录暂无可恢复的历史会话")
		return
	}
	m.resumeSelecting = true
	m.resumeCandidates = sessions
	m.appendMessage("system", renderResumeCandidates(sessions))
}

func (m *Model) handleResumeSelection(text string) bool {
	text = strings.TrimSpace(text)
	if text == "/cancel" || text == "/clear" {
		m.resumeSelecting = false
		m.resumeCandidates = nil
		if text == "/clear" {
			m.messages = nil
			m.appendMessage("system", "已清空当前 TUI 显示上下文")
		} else {
			m.appendMessage("system", "已取消恢复历史会话")
		}
		return true
	}
	idx, err := strconv.Atoi(text)
	if err != nil || idx < 1 || idx > len(m.resumeCandidates) {
		m.appendMessage("system", fmt.Sprintf("请输入 1-%d 之间的序号，或输入 /cancel 取消", len(m.resumeCandidates)))
		return true
	}
	candidate := m.resumeCandidates[idx-1]
	conv, history, err := m.session.Resumer.ResumeSession(context.Background(), candidate.SessionID, m.session.CWD, m.session.User)
	if err != nil {
		m.appendMessage("error", err.Error())
		return true
	}
	m.session.Conversation = conv
	m.resumeSelecting = false
	m.resumeCandidates = nil
	m.messages = messagesForDisplay(history)
	m.appendMessage("system", fmt.Sprintf("已恢复历史会话：%s", conv.SessionID))
	return true
}

func renderResumeCandidates(sessions []storage.ResumableSession) string {
	var b strings.Builder
	b.WriteString("可恢复的历史会话：")
	for i, session := range sessions {
		title := strings.TrimSpace(session.Title)
		if title == "" {
			title = "TUI 会话"
		}
		updated := "unknown"
		if !session.UpdatedAt.IsZero() {
			updated = session.UpdatedAt.Local().Format(time.RFC3339)
		}
		b.WriteString(fmt.Sprintf("\n%d. %s | %s | 消息:%d | %s", i+1, updated, title, session.MessageCount, session.SessionID))
	}
	b.WriteString("\n输入序号恢复，或输入 /cancel 取消。")
	return b.String()
}

func messagesForDisplay(history []storage.Message) []Message {
	if len(history) == 0 {
		return nil
	}
	messages := make([]Message, 0, len(history))
	for _, msg := range history {
		switch msg.Role {
		case "user", "assistant", "system", "error":
			messages = append(messages, Message{Role: msg.Role, Content: msg.Content, ReasoningContent: msg.ReasoningContent})
		}
	}
	return messages
}

func renderSkillDetails(skills []sessions.SkillSummary, fallbackCount int) string {
	count := len(skills)
	if count == 0 && fallbackCount > 0 {
		count = fallbackCount
	}
	if len(skills) == 0 {
		return fmt.Sprintf("已加载 Skills：%d 个", count)
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("已加载 Skills：%d 个", len(skills)))
	for _, skill := range skills {
		b.WriteString("\n- ")
		b.WriteString(skill.Name)
		if strings.TrimSpace(skill.Source) != "" {
			b.WriteString(" [")
			b.WriteString(skill.Source)
			b.WriteString("]")
		}
		if strings.TrimSpace(skill.Description) != "" {
			b.WriteString(" ")
			b.WriteString(skill.Description)
		}
		if strings.TrimSpace(skill.Path) != "" {
			b.WriteString("\n  path: ")
			b.WriteString(skill.Path)
		}
	}
	return b.String()
}

func renderMCPDetails(servers []mcp.ServerStatus, toolCount int) string {
	if len(servers) == 0 {
		return fmt.Sprintf("MCP Servers：0 个，工具：%d 个", toolCount)
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("MCP Servers：%d 个，工具：%d 个", len(servers), toolCount))
	for _, server := range servers {
		state := "disabled"
		if server.Connected {
			state = "connected"
		} else if server.LastError != "" {
			state = "failed"
		} else if server.Enabled {
			state = "pending"
		}
		b.WriteString("\n- ")
		b.WriteString(server.Name)
		b.WriteString(" [")
		b.WriteString(server.Transport)
		b.WriteString("] ")
		b.WriteString(state)
		b.WriteString(fmt.Sprintf(", tools: %d", server.ToolCount))
		if server.Command != "" {
			b.WriteString("\n  command: ")
			b.WriteString(strings.Join(append([]string{server.Command}, server.Args...), " "))
		}
		if server.URL != "" {
			b.WriteString("\n  url: ")
			b.WriteString(server.URL)
		}
		if server.LastError != "" {
			b.WriteString("\n  error: ")
			b.WriteString(server.LastError)
		}
	}
	return b.String()
}

func (m *Model) appendMessage(role, content string) {
	m.messages = append(m.messages, Message{Role: role, Content: content})
}

func (m *Model) appendAssistantDelta(delta string) {
	if len(m.messages) == 0 || !isLiveAssistantRole(m.messages[len(m.messages)-1].Role) {
		m.appendMessage("assistant", delta)
		return
	}
	if m.messages[len(m.messages)-1].Role == "thinking" {
		m.messages[len(m.messages)-1].Role = "assistant"
	}
	m.messages[len(m.messages)-1].Content += delta
}

func (m *Model) appendThinkingDelta(delta string) {
	if delta == "" {
		return
	}
	if len(m.messages) == 0 || !isLiveAssistantRole(m.messages[len(m.messages)-1].Role) {
		m.messages = append(m.messages, Message{Role: "assistant", ReasoningContent: delta})
		return
	}
	m.messages[len(m.messages)-1].ReasoningContent += delta
}

func (m *Model) replaceLastAssistant(content, reasoning string) {
	if len(m.messages) == 0 || !isLiveAssistantRole(m.messages[len(m.messages)-1].Role) {
		m.messages = append(m.messages, Message{Role: "assistant", Content: content, ReasoningContent: reasoning})
		return
	}
	m.messages[len(m.messages)-1].Role = "assistant"
	m.messages[len(m.messages)-1].Content = content
	m.messages[len(m.messages)-1].ReasoningContent = reasoning
}

func isLiveAssistantRole(role string) bool {
	return role == "assistant" || role == "thinking"
}

func (m *Model) refreshViewport() {
	m.viewport.Width = max(10, m.width)
	m.viewport.Height = m.viewportHeight()
	m.viewport.SetContent(m.renderTranscript())
	m.viewport.GotoBottom()
}

func (m Model) viewportHeight() int {
	if m.height <= 0 {
		return 20
	}
	return max(1, m.height)
}

func (m Model) renderTranscript() string {
	var b strings.Builder
	b.WriteString(m.renderWelcome())
	if len(m.messages) > 0 {
		b.WriteString("\n\n")
		b.WriteString(m.renderMessages())
	}
	b.WriteString("\n")
	b.WriteString(m.renderInput())
	b.WriteString("\n")
	b.WriteString(subtleStyle().Render(m.renderLiveStatus()))
	return b.String()
}

func (m Model) renderMessages() string {
	var b strings.Builder
	for _, msg := range m.messages {
		b.WriteString(m.renderMessage(msg))
		b.WriteString("\n\n")
	}
	return b.String()
}

func (m Model) renderMessage(msg Message) string {
	switch msg.Role {
	case "user":
		return promptLineStyle().Render("›") + " " + userStyle().Render(wrapText(msg.Content, m.messageWidth()-2))
	case "assistant":
		content := msg.Content
		if m.renderer != nil {
			if rendered, err := m.renderer.Render(content); err == nil {
				content = wrapText(strings.TrimSpace(rendered), m.messageWidth())
			}
		} else {
			content = wrapText(content, m.messageWidth())
		}
		if m.running && strings.TrimSpace(msg.ReasoningContent) != "" {
			content = thinkingStyle().Render("✽ 思考中\n"+wrapText(strings.TrimSpace(msg.ReasoningContent), m.messageWidth()-2)) + "\n" + content
		}
		return assistantLeadStyle().Render("go-agent") + "\n" + content
	case "thinking":
		return thinkingStyle().Render("✽ 思考中\n" + wrapText(msg.Content, m.messageWidth()-2))
	case "system":
		return systemStyle().Render("• " + wrapText(msg.Content, m.messageWidth()-2))
	case "error":
		return errorStyle().Render("✗ " + wrapText(msg.Content, m.messageWidth()-2))
	default:
		return roleLabel(msg.Role, lipgloss.Color("245")) + "\n" + wrapText(msg.Content, m.messageWidth())
	}
}

func (m Model) messageWidth() int {
	return max(10, m.width)
}

func newMarkdownRenderer(width int) *glamour.TermRenderer {
	renderer, _ := glamour.NewTermRenderer(glamour.WithStandardStyle("dark"), glamour.WithWordWrap(max(10, width)))
	return renderer
}

func wrapText(text string, width int) string {
	width = max(1, width)
	return ansi.Hardwrap(text, width, true)
}

func (m Model) renderWelcome() string {
	art := strings.Join([]string{
		`   /\_/\\`,
		`  ( o.o )   nano, but cozy`,
		`   > ^ <    ask · think · build`,
	}, "\n")
	intro := "像 Claude Code 一样一问一答：在终端直接提问，回答会流式显示在下方。"
	quick := "快捷键：Enter 发送 · Ctrl+C 中断/退出 · /resume 恢复 · /clear 清屏"
	stats := fmt.Sprintf("Skills %d · MCP tools %d", m.session.SkillCount, m.session.MCPToolCount)
	return startupPanelStyle().Width(max(20, m.width-2)).Render(accentArtStyle().Render(art) + "\n\n" + intro + "\n" + subtleStyle().Render(quick) + "\n" + subtleStyle().Render(stats))
}

func (m Model) renderHeader() string {
	width := max(20, m.width)
	status := runningText(m.running)
	contextText := "上下文 --"
	if m.contextBudget > 0 {
		contextText = fmt.Sprintf("上下文 %d%% · %s/%s", min(100, m.contextTokens*100/m.contextBudget), compactNumber(m.contextTokens), compactNumber(m.contextBudget))
	}
	left := titleStyle().Render("✦ go-agent")
	right := headerMetaStyle().Render(fmt.Sprintf("%s  ·  本轮工具 %d  ·  %s", status, m.toolCallCount, contextText))
	gap := max(1, width-lipgloss.Width(left)-lipgloss.Width(right)-4)
	line := left + strings.Repeat(" ", gap) + right
	workspace := subtleStyle().Render(fmt.Sprintf("cwd %s · skills %d · mcp tools %d", m.session.CWD, m.session.SkillCount, m.session.MCPToolCount))
	return headerStyle().Width(width).Render(line + "\n" + workspace)
}

func (m Model) renderConversationFrame() string {
	return conversationStyle().Width(max(10, m.width)).Render(m.viewport.View())
}

func (m Model) renderInput() string {
	prompt := inputPromptStyle().Render("›")
	text := strings.TrimSpace(m.input.Value())
	if text == "" {
		text = inputPromptStyle().Render(inputCursor) + " " + subtleStyle().Render(m.input.Placeholder)
	} else {
		text = userStyle().Render(text) + inputPromptStyle().Render(inputCursor)
	}
	return inputLineStyle().Width(max(10, m.width)).Render(prompt + " " + text)
}

func isTerminalProbeResponseInput(msg tea.KeyMsg) bool {
	if msg.Type != tea.KeyRunes {
		return false
	}
	value := string(msg.Runes)
	return strings.Contains(value, ";rgb:") && strings.Contains(value, "/")
}

func (m Model) renderLiveStatus() string {
	parts := []string{"Enter 发送", "Ctrl+C 中断/退出", "/help", fmt.Sprintf("工具 %d", m.toolCallCount)}
	if m.contextBudget > 0 {
		parts = append(parts, fmt.Sprintf("上下文 %d%% · %s/%s", min(100, m.contextTokens*100/m.contextBudget), compactNumber(m.contextTokens), compactNumber(m.contextBudget)))
	} else {
		parts = append(parts, "上下文 --")
	}
	return strings.Join(parts, " · ")
}

func (m *Model) updateMetaFromData(data any) {
	if count, ok := eventInt(data, "tool_call_count"); ok {
		m.toolCallCount = count
	}
	if tokens, ok := eventInt(data, "context_tokens"); ok {
		m.contextTokens = tokens
	}
	if budget, ok := eventInt(data, "context_budget"); ok {
		m.contextBudget = budget
	}
}

func eventString(data any, key string) string {
	m, ok := data.(map[string]any)
	if !ok {
		return ""
	}
	value, _ := m[key].(string)
	return value
}

func eventInt(data any, key string) (int, bool) {
	m, ok := data.(map[string]any)
	if !ok {
		return 0, false
	}
	switch value := m[key].(type) {
	case int:
		return value, true
	case int64:
		return int(value), true
	case float64:
		return int(value), true
	default:
		return 0, false
	}
}

func roleLabel(label string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Bold(true).Foreground(color).Render(label)
}

func headerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(tuiPalette.panel)
}

func titleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(tuiPalette.mint)
}

func headerMetaStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
}

func subtleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.muted)
}

func conversationStyle() lipgloss.Style {
	return lipgloss.NewStyle().Padding(1, 2, 0, 2)
}

func startupPanelStyle() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(tuiPalette.coral).Padding(1, 2).Margin(1, 0)
}

func userStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.ink)
}

func thinkingStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Italic(true).PaddingLeft(2)
}

func systemStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.butter)
}

func errorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.coral)
}

func promptLineStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.cyan).Bold(true)
}

func assistantLeadStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.mint).Bold(true)
}

func accentArtStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.lavender).Bold(true)
}

func inputPromptStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.mint).Bold(true)
}

func inputLineStyle() lipgloss.Style {
	return lipgloss.NewStyle().Background(lipgloss.Color("238")).Foreground(tuiPalette.ink)
}

func compactNumber(n int) string {
	if n >= 1000 {
		return fmt.Sprintf("%dk", n/1000)
	}
	return strconv.Itoa(n)
}

func runningText(running bool) string {
	if running {
		return "generating"
	}
	return "ready"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package tui\n\nimport (\n\t\"context\"\n\t\"fmt\"\n\t\"strconv\"\n\t\"strings\"\n\t\"time\"\n\n\t\"github.com/charmbracelet/bubbles/textarea\"\n\t\"github.com/charmbracelet/bubbles/viewport\"\n\ttea \"github.com/charmbracelet/bubbletea\"\n\t\"github.com/charmbracelet/glamour\"\n\t\"github.com/charmbracelet/lipgloss\"\n\t\"github.com/charmbracelet/x/ansi\"\n\n\t\"nano_cc/internal/agent/mcp\"\n\t\"nano_cc/internal/agent/runtime\"\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/logger\"\n\t\"nano_cc/internal/sessions\"\n)\n\ntype SessionInfo struct {\n\tUser"}
```

## 2026-06-13T09:16:11Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_CTRO0KlessvlmADi7cJg5111

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/handler.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/handler.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/handler.go: no such file or directory"}
```

## 2026-06-13T09:16:11Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_pv1ClWL0oNKjC7kZviz85748

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/sessions/manager.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/sessions/manager.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/sessions/manager.go: no such file or directory"}
```

## 2026-06-13T09:16:11Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_03_3COSefcglJPBWlQQhLIO0385

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager.go"}
```

### Result

```text
package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/logger"
)

const (
	toolNamePrefix  = "mcp__"
	connectTimeout  = 30 * time.Second
	callTimeout     = 120 * time.Second
	idleTimeout     = 10 * time.Minute
	cleanupInterval = 1 * time.Minute
)

// serverSession 持有一个已连接 MCP 服务器的会话与其工具元信息。
type serverSession struct {
	server      storage.MCPServer
	signature   string // 配置指纹，用于判断配置是否变更需重连
	session     *mcpsdk.ClientSession
	tools       []openai.Tool     // 已转换并加前缀的工具定义
	toolNames   map[string]string // 前缀名 -> 原始 MCP 工具名
	lastUsedAt  time.Time
	activeCalls int
	closing     bool
}

type ServerStatus struct {
	ID        string
	Name      string
	Scope     string
	Transport string
	Command   string
	Args      []string
	URL       string
	Enabled   bool
	Connected bool
	ToolCount int
	LastError string
}

type Snapshot struct {
	Servers    []ServerStatus
	ToolCount  int
	ErrorCount int
}

// Manager 管理本地 MCP 客户端连接，提供工具发现与调用能力。
type Manager struct {
	mu                sync.Mutex
	done              chan struct{}
	closeOnce         sync.Once
	builtinServers    map[string]storage.MCPServer
	builtinSessions   map[string]*serverSession
	workspaceServers  map[string]storage.MCPServer
	workspaceSessions map[string]*serverSession
	workspaceErrors   map[string]string
}

func NewManager() *Manager {
	manager := &Manager{
		done:              make(chan struct{}),
		builtinServers:    make(map[string]storage.MCPServer),
		builtinSessions:   make(map[string]*serverSession),
		workspaceServers:  make(map[string]storage.MCPServer),
		workspaceSessions: make(map[string]*serverSession),
		workspaceErrors:   make(map[string]string),
	}
	go manager.cleanupLoop()
	return manager
}

func serverSignature(s storage.MCPServer) string {
	data, _ := json.Marshal(struct {
		Transport string
		Command   string
		Args      []string
		Env       map[string]string
		URL       string
		Headers   map[string]string
	}{s.Transport, s.Command, s.Args, s.Env, s.URL, s.Headers})
	return string(data)
}

func sanitizeName(name string) string {
	var sb strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_':
			sb.WriteRune(r)
		default:
			sb.WriteRune('_')
		}
	}
	if sb.Len() == 0 {
		return "server"
	}
	return sb.String()
}

func prefixedToolName(serverName, toolName string) string {
	return toolNamePrefix + sanitizeName(serverName) + "__" + toolName
}

func (m *Manager) SetBuiltinServers(ctx context.Context, servers []storage.MCPServer) {
	m.mu.Lock()
	old := m.builtinSessions
	m.builtinServers = make(map[string]storage.MCPServer, len(servers))
	m.builtinSessions = make(map[string]*serverSession)
	for _, server := range servers {
		m.builtinServers[server.ID] = server
	}
	m.mu.Unlock()

	for _, sess := range old {
		closeSession(sess)
	}
	m.EnsureBuiltinSessions(ctx)
}

func (m *Manager) EnsureBuiltinSessions(ctx context.Context) {
	m.mu.Lock()
	missing := make([]storage.MCPServer, 0)
	for id, server := range m.builtinServers {
		if !server.Enabled {
			continue
		}
		if _, ok := m.builtinSessions[id]; !ok {
			missing = append(missing, server)
		}
	}
	m.mu.Unlock()

	for _, server := range missing {
		sess, err := connectBuiltinAndDiscover(ctx, server)
		if err != nil {
			logger.Warn(fmt.Sprintf("mcp: builtin stdio connect failed server=%s: %v", server.Name, err))
			continue
		}
		m.mu.Lock()
		if old := m.builtinSessions[server.ID]; old != nil {
			closeSession(old)
		}
		m.builtinSessions[server.ID] = sess
		m.mu.Unlock()
		logger.Info(fmt.Sprintf("mcp: builtin stdio connected server=%s tools=%d", server.Name, len(sess.tools)))
	}
}

func (m *Manager) SetWorkspaceServers(ctx context.Context, servers []storage.MCPServer) {
	m.mu.Lock()
	old := m.workspaceSessions
	m.workspaceServers = make(map[string]storage.MCPServer, len(servers))
	m.workspaceSessions = make(map[string]*serverSession)
	m.workspaceErrors = make(map[string]string)
	for _, server := range servers {
		m.workspaceServers[server.ID] = server
	}
	m.mu.Unlock()

	for _, sess := range old {
		closeSession(sess)
	}
	m.EnsureWorkspaceSessions(ctx)
}

func (m *Manager) EnsureWorkspaceSessions(ctx context.Context) {
	m.mu.Lock()
	missing := make([]storage.MCPServer, 0)
	for id, server := range m.workspaceServers {
		if !server.Enabled {
			continue
		}
		if _, ok := m.workspaceSessions[id]; !ok {
			missing = append(missing, server)
		}
	}
	m.mu.Unlock()

	for _, server := range missing {
		sess, err := connectAndDiscover(ctx, server)
		m.mu.Lock()
		if err != nil {
			m.workspaceErrors[server.ID] = err.Error()
			m.mu.Unlock()
			logger.Warn(fmt.Sprintf("mcp: workspace connect failed server=%s: %v", server.Name, err))
			continue
		}
		delete(m.workspaceErrors, server.ID)
		if old := m.workspaceSessions[server.ID]; old != nil {
			closeSession(old)
		}
		m.workspaceSessions[server.ID] = sess
		m.mu.Unlock()
		logger.Info(fmt.Sprintf("mcp: workspace connected server=%s tools=%d", server.Name, len(sess.tools)))
	}
}

func connectAndDiscover(ctx context.Context, server storage.MCPServer) (*serverSession, error) {
	transport, err := buildTransport(server)
	if err != nil {
		return nil, err
	}
	connCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	client := mcpsdk.NewClient(&mcpsdk.Implementation{Name: "nano_cc", Version: "1.0.0"}, nil)
	session, err := client.Connect(connCtx, transport, nil)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	tools, names, err := discoverTools(connCtx, session, server.Name)
	if err != nil {
		_ = session.Close()
		return nil, fmt.Errorf("list tools: %w", err)
	}
	return &serverSession{
		server:     server,
		signature:  serverSignature(server),
		session:    session,
		tools:      tools,
		toolNames:  names,
		lastUsedAt: time.Now(),
	}, nil
}

func connectBuiltinAndDiscover(ctx context.Context, server storage.MCPServer) (*serverSession, error) {
	transport, err := buildBuiltinStdioTransport(server)
	if err != nil {
		return nil, err
	}
	return connectWithTransport(ctx, server, transport)
}

func connectWithTransport(ctx context.Context, server storage.MCPServer, transport mcpsdk.Transport) (*serverSession, error) {
	connCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	client := mcpsdk.NewClient(&mcpsdk.Implementation{Name: "nano_cc", Version: "1.0.0"}, nil)
	session, err := client.Connect(connCtx, transport, nil)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	tools, names, err := discoverTools(connCtx, session, server.Name)
	if err != nil {
		_ = session.Close()
		return nil, fmt.Errorf("list tools: %w", err)
	}
	return &serverSession{
		server:     server,
		signature:  serverSignature(server),
		session:    session,
		tools:      tools,
		toolNames:  names,
		lastUsedAt: time.Now(),
	}, nil
}

func discoverTools(ctx context.Context, session *mcpsdk.ClientSession, serverName string) ([]openai.Tool, map[string]string, error) {
	result, err := session.ListTools(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	tools := make([]openai.Tool, 0, len(result.Tools))
	names := make(map[string]string, len(result.Tools))
	for _, tool := range result.Tools {
		prefixed := prefixedToolName(serverName, tool.Name)
		params, err := json.Marshal(tool.InputSchema)
		if err != nil {
			logger.Warn(fmt.Sprintf("mcp: marshal schema failed server=%s tool=%s: %v", serverName, tool.Name, err))
			continue
		}
		tools = append(tools, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        prefixed,
				Description: tool.Description,
				Parameters:  json.RawMessage(params),
			},
		})
		names[prefixed] = tool.Name
	}
	return tools, names, nil
}

// ToolsForUser 返回所有已连接本地 MCP 服务器发现到的工具定义（带前缀）。
func (m *Manager) ToolsForUser(userID string) []openai.Tool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.builtinSessions) == 0 && len(m.workspaceSessions) == 0 {
		return nil
	}
	tools := make([]openai.Tool, 0)
	now := time.Now()
	for _, sess := range m.builtinSessions {
		sess.lastUsedAt = now
		tools = append(tools, sess.tools...)
	}
	for _, sess := range m.workspaceSessions {
		sess.lastUsedAt = now
		tools = append(tools, sess.tools...)
	}
	sort.Slice(tools, func(i, j int) bool {
		return tools[i].Function.Name < tools[j].Function.Name
	})
	return tools
}

// CallTool 把带前缀的工具名路由到对应 MCP 服务器执行，返回文本结果。
func (m *Manager) CallTool(ctx context.Context, userID, prefixedName, rawArgs string) (string, error) {
	m.mu.Lock()
	var target *serverSession
	var originalName string
	var builtinID string
	for _, sess := range m.workspaceSessions {
		if name, ok := sess.toolNames[prefixedName]; ok {
			target = sess
			originalName = name
			target.activeCalls++
			break
		}
	}
	if target == nil {
		for id, sess := range m.builtinSessions {
			if name, ok := sess.toolNames[prefixedName]; ok {
				target = sess
				originalName = name
				builtinID = id
				target.activeCalls++
				break
			}
		}
	}
	m.mu.Unlock()

	if target == nil {
		return "", fmt.Errorf("mcp tool %s not found", prefixedName)
	}
	defer m.finishCall(target)

	var args map[string]any
	if strings.TrimSpace(rawArgs) != "" {
		if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
			return "", fmt.Errorf("invalid tool arguments: %w", err)
		}
	}
	if builtinID != "" {
		return m.callBuiltinTool(ctx, builtinID, target, originalName, args)
	}

	callCtx, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()

	result, err := target.session.CallTool(callCtx, &mcpsdk.CallToolParams{Name: originalName, Arguments: args})
	if err != nil {
		return "", err
	}
	output := serializeContent(result)
	if result.IsError {
		return "", fmt.Errorf("mcp tool error: %s", output)
	}
	return output, nil
}

func (m *Manager) callBuiltinTool(ctx context.Context, builtinID string, target *serverSession, originalName string, args map[string]any) (string, error) {
	output, err := callSessionTool(ctx, target, originalName, args)
	if err == nil {
		return output, nil
	}
	logger.Warn(fmt.Sprintf("mcp: builtin stdio call failed, reconnecting server=%s: %v", target.server.Name, err))
	m.reconnectBuiltin(ctx, builtinID, target)

	m.mu.Lock()
	retryTarget := m.builtinSessions[builtinID]
	if retryTarget != nil {
		retryTarget.activeCalls++
	}
	m.mu.Unlock()
	if retryTarget == nil {
		return "", err
	}
	defer m.finishCall(retryTarget)
	return callSessionTool(ctx, retryTarget, originalName, args)
}

func callSessionTool(ctx context.Context, target *serverSession, originalName string, args map[string]any) (string, error) {
	callCtx, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	result, err := target.session.CallTool(callCtx, &mcpsdk.CallToolParams{Name: originalName, Arguments: args})
	if err != nil {
		return "", err
	}
	output := serializeContent(result)
	if result.IsError {
		return "", fmt.Errorf("mcp tool error: %s", output)
	}
	return output, nil
}

func (m *Manager) reconnectBuiltin(ctx context.Context, builtinID string, old *serverSession) {
	m.mu.Lock()
	server, ok := m.builtinServers[builtinID]
	if current := m.builtinSessions[builtinID]; current == old {
		delete(m.builtinSessions, builtinID)
	}
	m.mu.Unlock()
	if !ok || !server.Enabled {
		return
	}
	closeSession(old)
	sess, err := connectBuiltinAndDiscover(ctx, server)
	if err != nil {
		logger.Warn(fmt.Sprintf("mcp: builtin stdio reconnect failed server=%s: %v", server.Name, err))
		return
	}
	m.mu.Lock()
	m.builtinSessions[builtinID] = sess
	m.mu.Unlock()
	logger.Info(fmt.Sprintf("mcp: builtin stdio reconnected server=%s tools=%d", server.Name, len(sess.tools)))
}

// Invalidate 保留兼容入口。本地 TUI 不再维护用户 DB MCP 连接，因此该方法不影响 builtin/workspace 会话。
func (m *Manager) Invalidate(userID string) {
}

// TestServer 临时连接一个配置并发现工具，返回工具名列表，供连接测试使用。
func (m *Manager) TestServer(ctx context.Context, server storage.MCPServer) ([]string, error) {
	sess, err := connectAndDiscover(ctx, server)
	if err != nil {
		return nil, err
	}
	defer sess.session.Close()
	names := make([]string, 0, len(sess.toolNames))
	for _, original := range sess.toolNames {
		names = append(names, original)
	}
	sort.Strings(names)
	return names, nil
}

// Close 关闭全部连接，用于服务退出。
func (m *Manager) Close() {
	m.closeOnce.Do(func() { close(m.done) })
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, sess := range m.builtinSessions {
		m.closeWhenIdleLocked(sess)
	}
	for _, sess := range m.workspaceSessions {
		m.closeWhenIdleLocked(sess)
	}
	m.builtinSessions = make(map[string]*serverSession)
	m.workspaceSessions = make(map[string]*serverSession)
}

func (m *Manager) Snapshot(userID string) Snapshot {
	m.mu.Lock()
	defer m.mu.Unlock()
	statuses := make([]ServerStatus, 0, len(m.builtinServers)+len(m.workspaceServers))
	addConfigured := func(scope string, server storage.MCPServer, sess *serverSession, lastErr string) {
		status := ServerStatus{ID: server.ID, Name: server.Name, Scope: scope, Transport: server.Transport, Command: server.Command, Args: append([]string(nil), server.Args...), URL: server.URL, Enabled: server.Enabled, Connected: sess != nil, LastError: lastErr}
		if sess != nil {
			status.ToolCount = len(sess.tools)
		}
		statuses = append(statuses, status)
	}
	for id, server := range m.builtinServers {
		addConfigured("builtin", server, m.builtinSessions[id], "")
	}
	for id, server := range m.workspaceServers {
		addConfigured("workspace", server, m.workspaceSessions[id], m.workspaceErrors[id])
	}
	sort.Slice(statuses, func(i, j int) bool {
		if statuses[i].Scope != statuses[j].Scope {
			return statuses[i].Scope < statuses[j].Scope
		}
		return statuses[i].Name < statuses[j].Name
	})
	snapshot := Snapshot{Servers: statuses}
	for _, status := range statuses {
		snapshot.ToolCount += status.ToolCount
		if status.LastError != "" {
			snapshot.ErrorCount++
		}
	}
	return snapshot
}

func (m *Manager) finishCall(sess *serverSession) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if sess.activeCalls > 0 {
		sess.activeCalls--
	}
	sess.lastUsedAt = time.Now()
	if sess.closing && sess.activeCalls == 0 {
		closeSession(sess)
	}
}

func (m *Manager) closeWhenIdleLocked(sess *serverSession) {
	if sess == nil {
		return
	}
	sess.closing = true
	if sess.activeCalls == 0 {
		closeSession(sess)
	}
}

func (m *Manager) cleanupLoop() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			m.cleanupIdleSessions(time.Now())
		case <-m.done:
			return
		}
	}
}

func (m *Manager) cleanupIdleSessions(now time.Time) {
	type idleSession struct {
		scope string
		id    string
		sess  *serverSession
	}
	var idle []idleSession

	m.mu.Lock()
	for serverID, sess := range m.workspaceSessions {
		if sess.activeCalls > 0 || now.Sub(sess.lastUsedAt) < idleTimeout {
			continue
		}
		idle = append(idle, idleSession{scope: "workspace", id: serverID, sess: sess})
		delete(m.workspaceSessions, serverID)
	}
	m.mu.Unlock()

	for _, item := range idle {
		closeSession(item.sess)
		logger.Info(fmt.Sprintf("mcp: idle session closed scope=%s id=%s server=%s", item.scope, item.id, item.sess.server.Name))
	}
}

func closeSession(sess *serverSession) {
	if sess == nil || sess.session == nil {
		return
	}
	_ = sess.session.Close()
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package mcp\n\nimport (\n\t\"context\"\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"sort\"\n\t\"strings\"\n\t\"sync\"\n\t\"time\"\n\n\tmcpsdk \"github.com/modelcontextprotocol/go-sdk/mcp\"\n\topenai \"github.com/sashabaranov/go-openai\"\n\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/logger\"\n)\n\nconst (\n\ttoolNamePrefix  = \"mcp__\"\n\tconnectTimeout  = 30 * time.Second\n\tcallTimeout     = 120 * time.Second\n\tidleTimeout     = 10 * time.Minute\n\tcleanupInterval = 1 * time.Minute\n)\n\n// serverSession 持有一个已连接 MCP 服务器的会话\ufffd\ufffd"}
```

## 2026-06-13T09:16:15Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_hpAdrk48Grqcl8o5ODJq7335

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/storage.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/storage.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/storage.go: no such file or directory"}
```

## 2026-06-13T09:16:15Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_MdDWitU5XDEljgDx6t468739

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/sqlite.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/sqlite.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/sqlite.go: no such file or directory"}
```

## 2026-06-13T09:16:15Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_ca0wOrgzoNtpSba1O1hp5037

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/models.go"}
```

### Result

```text
package storage

import "time"

type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	MemoryEnabled bool      `json:"memory_enabled"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// MCPServer 是用户配置的一个 MCP 服务器连接。Transport 取值 stdio/sse/streamable。
// Args/Env/Headers 从本地 .link/.mcp.json 读取，在 Go 侧以强类型表示。
type MCPServer struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Name      string            `json:"name"`
	Transport string            `json:"transport"`
	Command   string            `json:"command"`
	Args      []string          `json:"args"`
	Env       map[string]string `json:"env"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Enabled   bool              `json:"enabled"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type Conversation struct {
	ID            string    `json:"id"`
	SessionID     string    `json:"session_id,omitempty"`
	UserID        string    `json:"user_id"`
	RootMessageID string    `json:"root_message_id"`
	Title         string    `json:"title"`
	HistoryJSON   string    `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ResumableSession struct {
	SessionID      string    `json:"session_id"`
	ConversationID string    `json:"conversation_id"`
	WorkspaceRoot  string    `json:"workspace_root"`
	Title          string    `json:"title"`
	MessageCount   int       `json:"message_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// 该表已经不使用，作为存储的结构体，用于存储消息历史
type Message struct {
	ID               string            `json:"id"`
	ConversationID   string            `json:"conversation_id"`
	UserID           string            `json:"user_id"`
	Role             string            `json:"role"`
	Content          string            `json:"content"`
	ReasoningContent string            `json:"reasoning_content,omitempty"`
	ToolCallID       string            `json:"tool_call_id,omitempty"`
	ToolCalls        []MessageToolCall `json:"tool_calls,omitempty"`
	Meta             *MessageMeta      `json:"meta,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
}

// MessageMeta 记录助手回复的元信息，仅对最终 assistant 消息填充。
type MessageMeta struct {
	ToolCallCount int `json:"tool_call_count"`          // 调用工具次数（0 也需序列化，保证历史展示一致）
	ContextTokens int `json:"context_tokens,omitempty"` // 当前上下文估算 token
	ContextBudget int `json:"context_budget,omitempty"` // 上下文预算（用于算占比）
}

type MessageToolCall struct {
	ID       string              `json:"id"`
	Type     string              `json:"type"`
	Function MessageFunctionCall `json:"function"`
}

type MessageFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolCall struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	ToolName       string    `json:"tool_name"`
	Status         string    `json:"status"`
	Summary        string    `json:"summary"`
	CreatedAt      time.Time `json:"created_at"`
}

type PersistedOutput struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	MessageID      string    `json:"message_id"`
	ToolCallID     string    `json:"tool_call_id"`
	Kind           string    `json:"kind"`
	Strategy       string    `json:"strategy"`
	OriginalBytes  int       `json:"original_bytes"`
	ContentSHA256  string    `json:"content_sha256"`
	Content        string    `json:"content"`
	Preview        string    `json:"preview"`
	CreatedAt      time.Time `json:"created_at"`
}

type ToolResultLogEntry struct {
	ConversationID string    `json:"conversation_id"`
	SessionID      string    `json:"session_id,omitempty"`
	UserID         string    `json:"user_id"`
	ToolCallID     string    `json:"tool_call_id"`
	ToolName       string    `json:"tool_name"`
	RawArgs        string    `json:"raw_args"`
	Status         string    `json:"status"`
	Result         string    `json:"result"`
	AuditSummary   string    `json:"audit_summary,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type ContextSummary struct {
	ID                    string    `json:"id"`
	ConversationID        string    `json:"conversation_id"`
	UserID                string    `json:"user_id"`
	SourceHistorySHA256   string    `json:"source_history_sha256"`
	Strategy              string    `json:"strategy"`
	EstimatedTokensBefore int       `json:"estimated_tokens_before"`
	EstimatedTokensAfter  int       `json:"estimated_tokens_after"`
	Summary               string    `json:"summary"`
	CreatedAt             time.Time `json:"created_at"`
}

type SubagentMessage struct {
	ID               string            `json:"id"`
	RunID            string            `json:"run_id"`
	ParentToolCallID string            `json:"parent_tool_call_id"`
	ConversationID   string            `json:"conversation_id"`
	UserID           string            `json:"user_id"`
	SequenceNo       int               `json:"sequence_no"`
	Role             string            `json:"role"`
	Content          string            `json:"content"`
	ReasoningContent string            `json:"reasoning_content,omitempty"`
	ToolCallID       string            `json:"tool_call_id,omitempty"`
	ToolCalls        []MessageToolCall `json:"tool_calls,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
}

type Memory struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ConversationMemory 是单个会话维度、随每轮对话增量维护的“当前会话主干信息”
// 条目。它不注入 system prompt，仅在上下文压缩触发全量摘要时作为替代品使用。
type ConversationMemory struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package storage\n\nimport \"time\"\n\ntype User struct {\n\tID            string    `json:\"id\"`\n\tEmail         string    `json:\"email\"`\n\tUsername      string    `json:\"username\"`\n\tMemoryEnabled bool      `json:\"memory_enabled\"`\n\tCreatedAt     time.Time `json:\"created_at\"`\n\tUpdatedAt     time.Time `json:\"updated_at\"`\n}\n\n// MCPServer 是用户配置的一个 MCP 服务器连接。Transport 取值 stdio/sse/streamable。\n// Args/Env/Headers 从本地 .link/.mcp.json 读取，在 Go 侧以强类型表示。"}
```

## 2026-06-13T09:16:15Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_03_encTf2ibNzjOBSudbJ3F3440

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/runtime.go"}
```

### Result

```text
package runtime

import (
	"context"
	"time"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/runtime/compression"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/config"
	"nano_cc/internal/llm"
	"nano_cc/internal/sessions"
)

type EventWriter interface {
	Event(name string, data any) error
}

type conversationStore interface {
	UpdateConversationTitle(ctx context.Context, conversationID, title string) error
	TouchConversationActivity(ctx context.Context, conversationID string) error
	SetConversationHistory(ctx context.Context, conversationID string, messages []storage.Message) error
	SetConversationCache(ctx context.Context, conversationID string, messages []storage.Message) error
	GetConversationCache(ctx context.Context, conversationID string) ([]storage.Message, bool, error)
	ListMessagesByConversation(ctx context.Context, conversationID string, limit int) ([]storage.Message, error)
	CreateToolCall(ctx context.Context, tc storage.ToolCall) error
	CreateSubagentMessage(ctx context.Context, message storage.SubagentMessage) error
	CreatePersistedOutput(ctx context.Context, output storage.PersistedOutput) error
	GetPersistedOutputForConversation(ctx context.Context, id, userID, conversationID string) (storage.PersistedOutput, error)
	GetPersistedOutputByMessageHash(ctx context.Context, conversationID, userID, messageID, toolCallID, strategy, contentSHA256 string) (storage.PersistedOutput, error)
	CreateContextSummary(ctx context.Context, summary storage.ContextSummary) error
	GetContextSummaryByHistoryHash(ctx context.Context, conversationID, userID, sourceHistorySHA256 string) (storage.ContextSummary, error)
	ListRelevantMemories(ctx context.Context, userID string) ([]storage.Memory, error)
	ListMemoriesByUserAndType(ctx context.Context, userID, memType string) ([]storage.Memory, error)
	ListProjectFactMemories(ctx context.Context, userID string) ([]storage.Memory, error)
	InsertMemory(ctx context.Context, m storage.Memory) error
	CountMemoriesByUserAndType(ctx context.Context, userID, memType string) (int, error)
	CountProjectFactMemories(ctx context.Context, userID string) (int, error)
	DeleteOldestMemories(ctx context.Context, userID, memType string, n int) error
	ReplaceMemoriesByUserAndType(ctx context.Context, userID, memType string, items []storage.Memory) error
	ReplaceProjectFactMemories(ctx context.Context, userID string, items []storage.Memory) error
	ListConversationMemories(ctx context.Context, conversationID string) ([]storage.ConversationMemory, error)
	ReplaceConversationMemories(ctx context.Context, conversationID, userID string, items []storage.ConversationMemory) error
	GetConversationModelHistory(ctx context.Context, conversationID string) ([]storage.Message, bool, error)
	UpsertConversationModelHistory(ctx context.Context, conversationID, userID string, messages []storage.Message) error
	AcquireConversationLock(ctx context.Context, conversationID, token string, ttl, waitTimeout time.Duration) (bool, error)
	RenewConversationLock(ctx context.Context, conversationID, token string, ttl time.Duration) (bool, error)
	ReleaseConversationLock(ctx context.Context, conversationID, token string) error
}

type Service struct {
	Store             conversationStore
	Cfg               config.AppConfig
	LLM               llm.Client
	Tools             *ToolRegistry
	BuiltinSkills     *sessions.SkillLoader
	BasePrompt        string
	LinkMarkdown      config.LinkMarkdownContext
	Hooks             *HookManager
	ContextCompressor *compression.Compressor
	EnableMemory      bool
	MCP               *mcp.Manager
}

func NewService(store conversationStore, cfg config.AppConfig, client llm.Client) *Service {
	return &Service{Store: store, Cfg: cfg, LLM: client, Tools: NewToolRegistry(cfg), Hooks: NewDefaultHookManager(), EnableMemory: true}
}

func (s *Service) hookManager() *HookManager {
	if s.Hooks == nil {
		s.Hooks = NewDefaultHookManager()
	}
	return s.Hooks
}

func (s *Service) SetBuiltinSkills(loader *sessions.SkillLoader) {
	s.BuiltinSkills = loader
}

func (s *Service) SetBasePrompt(prompt string) {
	s.BasePrompt = prompt
}

func (s *Service) SetLinkMarkdownContext(ctx config.LinkMarkdownContext) {
	s.LinkMarkdown = ctx
}

func (s *Service) SetMCPManager(manager *mcp.Manager) {
	s.MCP = manager
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"context\"\n\t\"time\"\n\n\t\"nano_cc/internal/agent/mcp\"\n\t\"nano_cc/internal/agent/runtime/compression\"\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/config\"\n\t\"nano_cc/internal/llm\"\n\t\"nano_cc/internal/sessions\"\n)\n\ntype EventWriter interface {\n\tEvent(name string, data any) error\n}\n\ntype conversationStore interface {\n\tUpdateConversationTitle(ctx context.Context, conversationID, title string) error\n\tTouchConversationActivity(ctx context.Context, conversationID string) error"}
```

## 2026-06-13T09:16:15Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_04_xmDkMwuz4jSg6iPldkgx4058

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/prompt_builder.go"}
```

### Result

```text
package runtime

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/assistant"
	"nano_cc/internal/sessions"
	agenttools "nano_cc/internal/tools"
)

func (s *Service) buildSystemPrompt(ctx context.Context, conversation storage.Conversation, user storage.User, snapshot *agenttools.SkillSnapshot, history []storage.Message, memoryOn bool) string {
	memorySection := ""
	if memoryOn {
		memorySection = s.selectRelevantMemories(ctx, user, history)
	}
	return s.buildSystemPromptWithMemory(user, snapshot, memorySection)
}

func (s *Service) buildSystemPromptWithMemory(user storage.User, snapshot *agenttools.SkillSnapshot, memorySection string) string {
	toolNames := []string(nil)
	if s.Tools != nil {
		toolDefs := s.Tools.Definitions()
		toolNames = make([]string, 0, len(toolDefs))
		for _, def := range toolDefs {
			if def.Function == nil || strings.TrimSpace(def.Function.Name) == "" {
				continue
			}
			toolNames = append(toolNames, def.Function.Name)
		}
	}
	loader := sessions.NewSkillLoader()
	if snapshot != nil && snapshot.Merged != nil {
		loader = snapshot.Merged
	}
	return assistant.BuildSystemPrompt(assistant.PromptOptions{
		BasePrompt:        s.BasePrompt,
		Surface:           fmt.Sprintf("the local TUI session for user %s", user.Username),
		SkillDescriptions: loader.GetDescriptions(),
		MemorySection:     memorySection,
		WorkingDirectory:  strings.TrimSpace(s.Cfg.WorkspaceRoot),
		LinkMarkdown: assistant.LinkMarkdownContext{
			UserPath:         s.LinkMarkdown.UserPath,
			UserContent:      s.LinkMarkdown.UserContent,
			WorkspacePath:    s.LinkMarkdown.WorkspacePath,
			WorkspaceContent: s.LinkMarkdown.WorkspaceContent,
		},
		ToolNames: toolNames,
	})
}

func (s *Service) buildSkillSnapshot(ctx context.Context, userID string) (*agenttools.SkillSnapshot, error) {
	return agenttools.NewSkillSnapshot(nil, s.BuiltinSkills), nil
}

func buildOpenAIMessages(systemPrompt string, history []storage.Message) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{{Role: "system", Content: systemPrompt}}
	for _, msg := range history {
		messages = append(messages, openai.ChatCompletionMessage{Role: msg.Role, Content: msg.Content, ReasoningContent: msg.ReasoningContent, ToolCallID: msg.ToolCallID, ToolCalls: storageToolCallsToOpenAI(msg.ToolCalls)})
	}
	return messages
}

func storageToolCallsToOpenAI(toolCalls []storage.MessageToolCall) []openai.ToolCall {
	if len(toolCalls) == 0 {
		return nil
	}
	converted := make([]openai.ToolCall, 0, len(toolCalls))
	for _, toolCall := range toolCalls {
		converted = append(converted, openai.ToolCall{ID: toolCall.ID, Type: openai.ToolType(toolCall.Type), Function: openai.FunctionCall{Name: toolCall.Function.Name, Arguments: toolCall.Function.Arguments}})
	}
	return converted
}

func openAIToolCallsToStorage(toolCalls []openai.ToolCall) []storage.MessageToolCall {
	if len(toolCalls) == 0 {
		return nil
	}
	converted := make([]storage.MessageToolCall, 0, len(toolCalls))
	for _, toolCall := range toolCalls {
		converted = append(converted, storage.MessageToolCall{ID: toolCall.ID, Type: string(toolCall.Type), Function: storage.MessageFunctionCall{Name: toolCall.Function.Name, Arguments: toolCall.Function.Arguments}})
	}
	return converted
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"context\"\n\t\"fmt\"\n\t\"strings\"\n\n\topenai \"github.com/sashabaranov/go-openai\"\n\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/assistant\"\n\t\"nano_cc/internal/sessions\"\n\tagenttools \"nano_cc/internal/tools\"\n)\n\nfunc (s *Service) buildSystemPrompt(ctx context.Context, conversation storage.Conversation, user storage.User, snapshot *agenttools.SkillSnapshot, history []storage.Message, memoryOn bool) string {\n\tmemorySection := \"\"\n\tif memoryOn {\n\t\tmemorySection = s.selectReleva"}
```

## 2026-06-13T09:16:19Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_X0WBHQWSJS9Fl8RVQGfl0602

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/process.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/process.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/process.go: no such file or directory"}
```

## 2026-06-13T09:16:19Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_zQPfWPAOwwL48IAgmwWp9316

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/tools.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/tools.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/tools.go: no such file or directory"}
```

## 2026-06-13T09:16:19Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_5O5M7YpyCLb57g7AjWym6245

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/tool_registry.go"}
```

### Result

```text
package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	runtimehooks "nano_cc/internal/agent/runtime/hooks"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/config"
	agenttools "nano_cc/internal/tools"
)

type toolExecutionOutcome = runtimehooks.ToolExecutionOutcome
type toolExecutionAudit = runtimehooks.ToolExecutionAudit

type ToolContext struct {
	User                  storage.User
	Conversation          storage.Conversation
	Skills                *agenttools.SkillSnapshot
	ParentToolCallID      string
	PersistedOutputReader agenttools.PersistedOutputReader
}

type ToolExecutionResult struct {
	Output string
	Todos  []agenttools.TodoItem
}

type ToolRegistry struct {
	definitions []openai.Tool
	baseEnv     agenttools.RuntimeEnv
}

const defaultAllowedTool = "load_skill"

func NewToolRegistry(cfg config.AppConfig) *ToolRegistry {
	allowed := loadAllowedToolNames(cfg)
	definitions := appendPersistedOutputTool(buildToolDefinitions(allowed))
	return &ToolRegistry{
		definitions: definitions,
		baseEnv:     runtimeEnvFromConfig(cfg),
	}
}

// appendPersistedOutputTool always exposes read_persisted_output to the main
// agent so that <persisted-output> markers produced by context compression
// remain readable regardless of the configured AllowedTools list.
func appendPersistedOutputTool(defs []openai.Tool) []openai.Tool {
	for _, def := range defs {
		if def.Function != nil && def.Function.Name == agenttools.ReadPersistedOutputToolName {
			return defs
		}
	}
	return append(defs, agenttools.ReadPersistedOutputToolDef)
}

func NewChildToolRegistry(cfg config.AppConfig, cwd string) *ToolRegistry {
	allowed := withoutTool(loadAllowedToolNames(cfg), "spawn_subagent")
	env := runtimeEnvFromConfig(cfg)
	env.CurrentWorkingDir = strings.TrimSpace(cwd)
	env.AllowOutsideWorkspace = false
	return &ToolRegistry{definitions: buildToolDefinitions(allowed), baseEnv: env}
}

func (r *ToolRegistry) Definitions() []openai.Tool {
	return append([]openai.Tool(nil), r.definitions...)
}

// toolDefinitionsForUser 返回内置工具定义，并在 MCP 启用时合并该用户已连接 MCP 服务器的工具。
func (s *Service) toolDefinitionsForUser(ctx context.Context, userID string) []openai.Tool {
	defs := s.Tools.Definitions()
	if s.MCP == nil {
		return defs
	}
	s.MCP.EnsureBuiltinSessions(ctx)
	s.MCP.EnsureWorkspaceSessions(ctx)
	mcpTools := s.MCP.ToolsForUser(userID)
	if len(mcpTools) == 0 {
		return defs
	}
	return append(defs, mcpTools...)
}

func (r *ToolRegistry) Execute(ctx context.Context, toolCtx ToolContext, name string, rawArgs string) (ToolExecutionResult, error) {
	var args map[string]any
	if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
		return ToolExecutionResult{}, fmt.Errorf("invalid tool arguments: %w", err)
	}
	if !r.isAllowed(name) {
		return ToolExecutionResult{}, fmt.Errorf("tool %s is not registered for local runtime", name)
	}
	if def, ok := r.lookupDefinition(name); ok && def.Function != nil {
		if err := agenttools.ValidateToolArgs(name, agenttools.RawSchemaFromParameters(def.Function.Parameters), args); err != nil {
			return ToolExecutionResult{}, err
		}
	}
	ctx = agenttools.WithRuntimeEnv(ctx, r.runtimeEnv())
	ctx = agenttools.WithSkillSnapshot(ctx, toolCtx.Skills)
	if toolCtx.PersistedOutputReader != nil {
		ctx = agenttools.WithPersistedOutputReader(ctx, toolCtx.PersistedOutputReader)
	}
	execResult, err := agenttools.Dispatch(ctx, name, args)
	if err != nil {
		return ToolExecutionResult{}, err
	}
	return ToolExecutionResult{Output: execResult.Output, Todos: execResult.Todos}, nil
}

func (r *ToolRegistry) runtimeEnv() agenttools.RuntimeEnv {
	env := r.baseEnv
	workspaceRoot := strings.TrimSpace(env.WorkspaceRoot)
	currentWorkingDir := strings.TrimSpace(env.CurrentWorkingDir)
	if currentWorkingDir == "" {
		currentWorkingDir = workspaceRoot
	}
	return agenttools.RuntimeEnv{
		AppHome:                env.AppHome,
		CommandBinDir:          strings.TrimSpace(env.CommandBinDir),
		CommandScriptDir:       strings.TrimSpace(env.CommandScriptDir),
		WorkspaceRoot:          workspaceRoot,
		CurrentWorkingDir:      currentWorkingDir,
		AllowOutsideWorkspace:  env.AllowOutsideWorkspace,
		AllowDangerousCommands: env.AllowDangerousCommands,
	}
}

func withoutTool(names []string, excluded string) []string {
	filtered := make([]string, 0, len(names))
	for _, name := range names {
		if name == excluded {
			continue
		}
		filtered = append(filtered, name)
	}
	return filtered
}

func (r *ToolRegistry) isAllowed(name string) bool {
	for _, tool := range r.definitions {
		if tool.Function != nil && tool.Function.Name == name {
			return true
		}
	}
	return false
}

func (r *ToolRegistry) lookupDefinition(name string) (openai.Tool, bool) {
	for _, tool := range r.definitions {
		if tool.Function != nil && tool.Function.Name == name {
			return tool, true
		}
	}
	return openai.Tool{}, false
}

func loadAllowedToolNames(cfg config.AppConfig) []string {
	configured := cfg.AllowedTools
	if len(configured) == 0 {
		configured = []string{defaultAllowedTool}
	}

	names := make([]string, 0, len(configured))
	seen := make(map[string]struct{}, len(configured))
	for _, name := range configured {
		if _, ok := lookupRegisteredTool(name); !ok {
			continue
		}
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}
		names = append(names, name)
	}
	return names
}

func runtimeEnvFromConfig(cfg config.AppConfig) agenttools.RuntimeEnv {
	return agenttools.RuntimeEnv{
		AppHome:                strings.TrimSpace(cfg.AppHome),
		CommandBinDir:          strings.TrimSpace(cfg.CommandBinDir),
		CommandScriptDir:       strings.TrimSpace(cfg.CommandScriptDir),
		WorkspaceRoot:          strings.TrimSpace(cfg.WorkspaceRoot),
		AllowOutsideWorkspace:  cfg.BashAllowOutsideWorkspace,
		AllowDangerousCommands: cfg.BashAllowDangerousCommands,
	}
}

func buildToolDefinitions(allowed []string) []openai.Tool {
	toolDefs := make([]openai.Tool, 0, len(allowed))
	for _, name := range allowed {
		def, ok := lookupRegisteredTool(name)
		if !ok {
			continue
		}
		toolDefs = append(toolDefs, def)
	}
	return toolDefs
}

func lookupRegisteredTool(name string) (openai.Tool, bool) {
	for _, tool := range agenttools.AllToolDefs {
		if tool.Function != nil && tool.Function.Name == name {
			return tool, true
		}
	}
	return openai.Tool{}, false
}

func RegisteredTools(cfg config.AppConfig) []string {
	names := loadAllowedToolNames(cfg)
	sort.Strings(names)
	return names
}

func (s *Service) executeToolCall(ctx context.Context, toolCtx ToolContext, name string, rawArgs string, audit toolExecutionAudit) toolExecutionOutcome {
	if strings.HasPrefix(name, "mcp__") {
		if s.MCP == nil {
			return toolExecutionOutcome{Status: "rejected", Result: "Error: MCP is not enabled", Audit: audit}
		}
		output, err := s.MCP.CallTool(ctx, toolCtx.User.ID, name, rawArgs)
		if err != nil {
			return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Error: %v", err), Audit: audit}
		}
		return toolExecutionOutcome{Status: "success", Result: output, Audit: audit}
	}
	if name == "spawn_subagent" {
		if s.Tools == nil || !s.Tools.isAllowed(name) {
			return toolExecutionOutcome{Status: "rejected", Result: "Error: tool spawn_subagent is not registered for local runtime", Audit: audit}
		}
		var rawMap map[string]any
		if err := json.Unmarshal([]byte(rawArgs), &rawMap); err != nil {
			return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Error: invalid spawn_subagent arguments: %v", err), Audit: audit}
		}
		if def, ok := s.Tools.lookupDefinition(name); ok && def.Function != nil {
			if err := agenttools.ValidateToolArgs(name, agenttools.RawSchemaFromParameters(def.Function.Parameters), rawMap); err != nil {
				return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Error: %v", err), Audit: audit}
			}
		}
		var args spawnSubagentArgs
		if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
			return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Error: invalid spawn_subagent arguments: %v", err), Audit: audit}
		}
		result, err := s.runSubagent(ctx, toolCtx, args, audit)
		if err != nil {
			return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Subagent failed: %v", err), Audit: audit}
		}
		return toolExecutionOutcome{Status: "success", Result: result, Audit: audit}
	}
	execResult, err := s.Tools.Execute(ctx, toolCtx, name, rawArgs)
	if err != nil {
		return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Error: %v", err), Audit: audit}
	}
	return toolExecutionOutcome{Status: "success", Result: execResult.Output, Audit: audit, Todos: execResult.Todos}
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"context\"\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"sort\"\n\t\"strings\"\n\n\topenai \"github.com/sashabaranov/go-openai\"\n\n\truntimehooks \"nano_cc/internal/agent/runtime/hooks\"\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/config\"\n\tagenttools \"nano_cc/internal/tools\"\n)\n\ntype toolExecutionOutcome = runtimehooks.ToolExecutionOutcome\ntype toolExecutionAudit = runtimehooks.ToolExecutionAudit\n\ntype ToolContext struct {\n\tUser                  storage.User\n\tConversation          storage.Conver"}
```

## 2026-06-13T09:16:19Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_03_k7qBPBH8IeS17bonDkuT9487

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks.go: no such file or directory"}
```

## 2026-06-13T09:16:19Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_04_6KEHCdP3eApICjGgWe6X8403

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/conversation_memory.go"}
```

### Result

```text
package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/idgen"
	"nano_cc/internal/logger"
)

const conversationMemorySystemPrompt = `你是“当前会话记忆”维护引擎。给你“已有会话记忆条目”和“最新一轮对话”，
请输出更新后的【完整】会话记忆条目列表，使其准确反映本场会话至今的主干信息：
- 覆盖：当前用户目标、关键决策与结论、已完成/已产出的内容、重要约束与上下文、待办与下一步。
- 合并重复、用新信息更新旧条目、删除已过时或被推翻的内容；只重组已知信息，不编造。
- name：短标题(<=80字，可用 [目标]/[决策]/[产出]/[待办] 等前缀)。description：一句话要点(<=300字)。body：支撑细节(<=2000字)。
- 仅输出 JSON 数组：[{"name","description","body"}]。无可记录时输出 []。`

type extractedConversationMemory struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Body        string `json:"body"`
}

// updateConversationMemory runs once at the end of a conversation turn: it asks
// the LLM to rewrite the full conversation memory list based on existing entries
// plus the latest dialogue, then replaces the stored entries for this
// conversation. It is best-effort: failures are logged and swallowed so the
// user-facing response is never affected. The result is never injected into the
// system prompt; it is only consumed by the context compression pipeline.
func (s *Service) updateConversationMemory(ctx context.Context, conversation storage.Conversation, user storage.User, history []storage.Message) {
	if s.LLM == nil {
		return
	}
	dialogue := renderDialogueForMemory(history)
	if strings.TrimSpace(dialogue) == "" {
		return
	}
	existing, err := s.Store.ListConversationMemories(ctx, conversation.ID)
	if err != nil {
		logger.Warn(fmt.Sprintf("conversation memory: load existing failed: %v", err))
	}
	resp, err := s.LLM.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.Cfg.LLM.ModelID,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: conversationMemorySystemPrompt},
			{Role: "user", Content: buildConversationMemoryUserPrompt(existing, dialogue)},
		},
	})
	if err != nil {
		logger.Warn(fmt.Sprintf("conversation memory: update failed: %v", err))
		return
	}
	if len(resp.Choices) == 0 {
		return
	}
	items := parseConversationMemories(resp.Choices[0].Message.Content)
	refined := make([]storage.ConversationMemory, 0, len(items))
	for _, it := range items {
		refined = append(refined, storage.ConversationMemory{
			ID:          idgen.New("cm"),
			Name:        it.Name,
			Description: it.Description,
			Body:        it.Body,
		})
	}
	if err := s.Store.ReplaceConversationMemories(ctx, conversation.ID, user.ID, refined); err != nil {
		logger.Warn(fmt.Sprintf("conversation memory: replace failed: %v", err))
	}
}

// parseConversationMemories extracts a JSON array of memory objects from model
// output, tolerating surrounding code fences or prose, then truncates fields.
func parseConversationMemories(raw string) []extractedConversationMemory {
	trimmed := extractJSONArray(raw)
	if trimmed == "" {
		return nil
	}
	var parsed []extractedConversationMemory
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return nil
	}
	result := make([]extractedConversationMemory, 0, len(parsed))
	for _, it := range parsed {
		it.Name = strings.TrimSpace(it.Name)
		it.Description = strings.TrimSpace(it.Description)
		it.Body = strings.TrimSpace(it.Body)
		if it.Name == "" || it.Body == "" {
			continue
		}
		it.Name = truncateRunes(it.Name, maxMemoryNameRunes)
		it.Description = truncateRunes(it.Description, maxMemoryDescRunes)
		it.Body = truncateRunes(it.Body, maxMemoryBodyRunes)
		result = append(result, it)
	}
	return result
}

func buildConversationMemoryUserPrompt(existing []storage.ConversationMemory, dialogue string) string {
	var b strings.Builder
	b.WriteString("Existing conversation memory entries:\n")
	if len(existing) == 0 {
		b.WriteString("(none)")
	} else {
		for i, m := range existing {
			b.WriteString(fmt.Sprintf("[%d] %s: %s\n", i, m.Name, m.Description))
			if strings.TrimSpace(m.Body) != "" {
				b.WriteString("    body: " + m.Body + "\n")
			}
		}
	}
	b.WriteString("\n\nLatest dialogue:\n")
	b.WriteString(dialogue)
	return b.String()
}

// scheduleMemoryWork 在一轮对话结束后异步执行收尾操作（模型历史持久化 + 记忆提取
// + 会话记忆更新）。
// 它接管入口持有的会话锁（token）：在独立的 background context 中执行，期间持续
// 续期，完成后停止续期并释放锁。返回 true 表示已接管锁所有权（调用方应跳过 defer
// 释放）；返回 false 表示未持锁（已降级），调用方按原逻辑处理。
// memoryOn 仅控制记忆提取与会话记忆更新；锁释放与模型历史持久化始终执行。
func (s *Service) scheduleMemoryWork(conv storage.Conversation, user storage.User, history []storage.Message, modelHistory []storage.Message, token string, stopRenew func(), memoryOn bool) bool {
	if token == "" {
		// 入口未持锁（已降级）→ 跳过收尾，不接管锁。
		return false
	}
	// 停止请求期看门狗，收尾 goroutine 内重新启动一个。
	if stopRenew != nil {
		stopRenew()
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Warn(fmt.Sprintf("memory work: panic recovered conversation=%s: %v", conv.ID, r))
			}
		}()
		defer s.Store.ReleaseConversationLock(context.Background(), conv.ID, token)
		stop := s.startLockRenewer(conv.ID, token)
		defer stop()

		ctx, cancel := context.WithTimeout(context.Background(), s.Cfg.MemoryWorkTimeout)
		defer cancel()
		// 先落库模型历史：即使后续记忆相关的 LLM 调用超时，也不丢失本轮压缩成果。
		if len(modelHistory) > 0 {
			if err := s.Store.UpsertConversationModelHistory(ctx, conv.ID, user.ID, modelHistory); err != nil {
				logger.Warn(fmt.Sprintf("model history: persist failed conversation=%s: %v", conv.ID, err))
			}
		}
		if memoryOn {
			s.extractMemories(ctx, user, history)
			s.updateConversationMemory(ctx, conv, user, history)
		}
	}()
	return true
}

// startLockRenewer 启动一个后台看门狗，按 TTL/3 周期为会话锁续期，返回的 stop
// 函数用于停止续期（幂等）。续期失败（锁已不属于当前 token）时记录告警并停止。
func (s *Service) startLockRenewer(conversationID, token string) func() {
	ttl := s.Cfg.ConversationLockTTL
	interval := ttl / 3
	if interval <= 0 {
		interval = ttl
	}
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				ok, err := s.Store.RenewConversationLock(context.Background(), conversationID, token, ttl)
				if err != nil {
					logger.Warn(fmt.Sprintf("conversation lock: renew failed conversation=%s: %v", conversationID, err))
					return
				}
				if !ok {
					logger.Warn(fmt.Sprintf("conversation lock: renew lost ownership conversation=%s", conversationID))
					return
				}
			}
		}
	}()
	var once sync.Once
	return func() {
		once.Do(func() { close(done) })
	}
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"context\"\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"strings\"\n\t\"sync\"\n\t\"time\"\n\n\topenai \"github.com/sashabaranov/go-openai\"\n\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/idgen\"\n\t\"nano_cc/internal/logger\"\n)\n\nconst conversationMemorySystemPrompt = `你是“当前会话记忆”维护引擎。给你“已有会话记忆条目”和“最新一轮对话”，\n请输出更新后的【完整】会话记忆条目列表，使其准确反映本场会话至今的主干信息：\n- 覆盖\ufffd"}
```

## 2026-06-13T09:16:19Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_05_rHhGBdFFI2eNeMssIBAC7023

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/memory.go"}
```

### Result

```text
package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/idgen"
	"nano_cc/internal/logger"
)

const (
	MemoryTypeEpisodicMemory = "episodic_memory"
	MemoryTypeUserPreference = "user_preference"
	MemoryTypeProjectFact    = "project_fact"
)

const (
	maxPreferenceMemories    = 20
	maxSessionSummaries      = 20
	sessionSummaryPruneCount = 10
	maxProjectFactMemories   = 50
	maxInjectedMemories      = 10
	maxMemoryDialogueChars   = 4000
	maxMemoryNameRunes       = 80
	maxMemoryDescRunes       = 300
	maxMemoryBodyRunes       = 2000
)

func validMemoryType(t string) bool {
	switch t {
	case MemoryTypeEpisodicMemory, MemoryTypeUserPreference, MemoryTypeProjectFact:
		return true
	}
	return false
}

const memoryExtractionSystemPrompt = `You are a project-scoped long-term memory extraction engine for a personal assistant named Link.
All memories are valid ONLY for the current project. Do not create memories that should be reused in other projects.
From the dialogue, extract durable memories worth keeping for the current project. Use the "type" field for three kinds:
- "episodic_memory": a concrete event/experience in this project session. Preserve factual integrity and temporal order; summarize what happened, not raw messages.
- "user_preference": a stable user preference, constraint, or recurring habit that applies to this project.
- "project_fact": a reusable fact about the current project, such as architecture, commands, conventions, dependencies, known constraints, or implementation decisions. It is not general world knowledge and must not be treated as valid outside this project.

Rules:
- Only extract NEW information not already covered by "Existing memories".
- Do not store one-off, trivial, or sensitive private data (passwords, payment, health).
- Do not extract facts about other projects.
- "name": short title (<=80 chars). "description": one-sentence gist (<=300 chars). "body": supporting detail (<=2000 chars).
- Output ONLY a JSON array: [{"name","type","description","body"}].
- If nothing new or everything is already covered, output exactly [].`

const memorySelectionSystemPrompt = `You are a project-scoped memory retrieval engine. Given the current project conversation context and a numbered list of candidate memories from this project's memory index, select the ones RELEVANT and USEFUL for answering the user right now.
- Candidate memories are valid only for the current project.
- Select at most 10.
- Prefer specific, on-topic memories; ignore unrelated ones.
- Output ONLY a JSON array of the selected memory indices, e.g. [0,3,7]. If none, output [].`

func memoryConsolidationSystemPrompt(typeLabel, typeValue string) string {
	return fmt.Sprintf(`You are a project-scoped memory consolidation engine. You are given the FULL current list of "%s" memories for the current project.
Rewrite them into a clean, minimal set:
- Merge duplicates and near-duplicates into a single entry.
- Reconcile contradictions, keeping the most recent / most reliable statement.
- Drop outdated or superseded memories.
- Keep only facts, preferences, and events that are valid for the current project; do not create cross-project memories.
- Never invent new facts; only reorganize what is given.
- Keep limits: name <=80, description <=300, body <=2000 chars.
Output ONLY a JSON array [{"name","type","description","body"}] representing the COMPLETE refined list (it fully replaces the old list). All entries must have type "%s".`, typeLabel, typeValue)
}

type extractedMemory struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Body        string `json:"body"`
}

// extractMemories runs once at the end of a conversation turn: it asks the LLM
// to extract the three kinds of memories from the dialogue, persists them, and
// triggers consolidation/pruning by type. It is best-effort: failures are
// logged and swallowed so the user-facing response is never affected.
func (s *Service) extractMemories(ctx context.Context, user storage.User, history []storage.Message) {
	if s.LLM == nil {
		return
	}
	dialogue := renderDialogueForMemory(history)
	if strings.TrimSpace(dialogue) == "" {
		return
	}
	existing, err := s.Store.ListRelevantMemories(ctx, user.ID)
	if err != nil {
		logger.Warn(fmt.Sprintf("memory: load existing memories failed: %v", err))
	}
	resp, err := s.LLM.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.Cfg.LLM.ModelID,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: memoryExtractionSystemPrompt},
			{Role: "user", Content: buildExtractionUserPrompt(existing, dialogue)},
		},
	})
	if err != nil {
		logger.Warn(fmt.Sprintf("memory: extraction failed: %v", err))
		return
	}
	if len(resp.Choices) == 0 {
		return
	}
	items := parseExtractedMemories(resp.Choices[0].Message.Content)
	if len(items) == 0 {
		return
	}

	var touchedPref, touchedSession, touchedProjectFact bool
	for _, it := range items {
		uid := user.ID
		if err := s.Store.InsertMemory(ctx, storage.Memory{
			ID:          idgen.New("mem"),
			UserID:      uid,
			Type:        it.Type,
			Name:        it.Name,
			Description: it.Description,
			Body:        it.Body,
		}); err != nil {
			logger.Warn(fmt.Sprintf("memory: insert failed: %v", err))
			continue
		}
		switch it.Type {
		case MemoryTypeUserPreference:
			touchedPref = true
		case MemoryTypeEpisodicMemory:
			touchedSession = true
		case MemoryTypeProjectFact:
			touchedProjectFact = true
		}
	}
	if touchedPref {
		s.maybeConsolidateUserPreferences(ctx, user.ID)
	}
	if touchedSession {
		s.maybePruneSessionSummaries(ctx, user.ID)
	}
	if touchedProjectFact {
		s.maybeConsolidateProjectFacts(ctx, user.ID)
	}
}

func (s *Service) maybeConsolidateUserPreferences(ctx context.Context, userID string) {
	items, err := s.Store.ListMemoriesByUserAndType(ctx, userID, MemoryTypeUserPreference)
	if err != nil || len(items) < maxPreferenceMemories {
		return
	}
	refined := s.consolidateViaLLM(ctx, "user_preference", MemoryTypeUserPreference, items)
	if refined == nil {
		return
	}
	if err := s.Store.ReplaceMemoriesByUserAndType(ctx, userID, MemoryTypeUserPreference, refined); err != nil {
		logger.Warn(fmt.Sprintf("memory: replace user preferences failed: %v", err))
	}
}

func (s *Service) maybeConsolidateProjectFacts(ctx context.Context, userID string) {
	items, err := s.Store.ListProjectFactMemories(ctx, userID)
	if err != nil || len(items) < maxProjectFactMemories {
		return
	}
	refined := s.consolidateViaLLM(ctx, "project_fact", MemoryTypeProjectFact, items)
	if refined == nil {
		return
	}
	if err := s.Store.ReplaceProjectFactMemories(ctx, userID, refined); err != nil {
		logger.Warn(fmt.Sprintf("memory: replace project fact memories failed: %v", err))
	}
}

func (s *Service) maybePruneSessionSummaries(ctx context.Context, userID string) {
	n, err := s.Store.CountMemoriesByUserAndType(ctx, userID, MemoryTypeEpisodicMemory)
	if err != nil || n < maxSessionSummaries {
		return
	}
	if err := s.Store.DeleteOldestMemories(ctx, userID, MemoryTypeEpisodicMemory, sessionSummaryPruneCount); err != nil {
		logger.Warn(fmt.Sprintf("memory: prune session summaries failed: %v", err))
	}
}

// consolidateViaLLM feeds the full memory list to the model and parses the
// refined complete list. Returns nil on failure so the caller leaves data
// untouched.
func (s *Service) consolidateViaLLM(ctx context.Context, typeLabel, typeValue string, items []storage.Memory) []storage.Memory {
	if s.LLM == nil {
		return nil
	}
	resp, err := s.LLM.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.Cfg.LLM.ModelID,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: memoryConsolidationSystemPrompt(typeLabel, typeValue)},
			{Role: "user", Content: buildConsolidationUserPrompt(items)},
		},
	})
	if err != nil {
		logger.Warn(fmt.Sprintf("memory: consolidation failed: %v", err))
		return nil
	}
	if len(resp.Choices) == 0 {
		return nil
	}
	parsed := parseExtractedMemories(resp.Choices[0].Message.Content)
	refined := make([]storage.Memory, 0, len(parsed))
	for _, it := range parsed {
		if it.Type != typeValue {
			continue
		}
		refined = append(refined, storage.Memory{
			ID:          idgen.New("mem"),
			Type:        typeValue,
			Name:        it.Name,
			Description: it.Description,
			Body:        it.Body,
		})
	}
	return refined
}

// selectRelevantMemories runs once before the conversation loop: it asks the
// LLM to pick the memories most relevant to the current context and renders
// them into the MemorySection. Best-effort: returns "" on failure or no data.
func (s *Service) selectRelevantMemories(ctx context.Context, user storage.User, history []storage.Message) string {
	if !s.EnableMemory || s.LLM == nil {
		return ""
	}
	all, err := s.Store.ListRelevantMemories(ctx, user.ID)
	if err != nil {
		logger.Warn(fmt.Sprintf("memory: load relevant memories failed: %v", err))
		return ""
	}
	if len(all) == 0 {
		return ""
	}
	resp, err := s.LLM.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.Cfg.LLM.ModelID,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: memorySelectionSystemPrompt},
			{Role: "user", Content: buildSelectionUserPrompt(all, renderDialogueForMemory(history))},
		},
	})
	if err != nil {
		logger.Warn(fmt.Sprintf("memory: selection failed: %v", err))
		return ""
	}
	if len(resp.Choices) == 0 {
		return ""
	}
	indices := parseSelectedIDs(resp.Choices[0].Message.Content)
	selected := pickMemoriesByIndex(all, indices, maxInjectedMemories)
	return renderMemorySection(selected)
}

// parseExtractedMemories extracts a JSON array of memory objects from model
// output, tolerating surrounding code fences or prose, then validates types
// and truncates fields.
func parseExtractedMemories(raw string) []extractedMemory {
	trimmed := extractJSONArray(raw)
	if trimmed == "" {
		return nil
	}
	var parsed []extractedMemory
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return nil
	}
	result := make([]extractedMemory, 0, len(parsed))
	for _, it := range parsed {
		it.Type = strings.TrimSpace(it.Type)
		it.Name = strings.TrimSpace(it.Name)
		it.Description = strings.TrimSpace(it.Description)
		it.Body = strings.TrimSpace(it.Body)
		if !validMemoryType(it.Type) || it.Name == "" || it.Body == "" {
			continue
		}
		it.Name = truncateRunes(it.Name, maxMemoryNameRunes)
		it.Description = truncateRunes(it.Description, maxMemoryDescRunes)
		it.Body = truncateRunes(it.Body, maxMemoryBodyRunes)
		result = append(result, it)
	}
	return result
}

// parseSelectedIDs parses a JSON array of integer indices from model output.
func parseSelectedIDs(raw string) []int {
	trimmed := extractJSONArray(raw)
	if trimmed == "" {
		return nil
	}
	var parsed []int
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return nil
	}
	return parsed
}

// pickMemoriesByIndex resolves indices into memories, dropping out-of-range and
// duplicate indices, capped at max.
func pickMemoriesByIndex(all []storage.Memory, indices []int, max int) []storage.Memory {
	seen := make(map[int]struct{})
	result := make([]storage.Memory, 0, len(indices))
	for _, idx := range indices {
		if idx < 0 || idx >= len(all) {
			continue
		}
		if _, exists := seen[idx]; exists {
			continue
		}
		seen[idx] = struct{}{}
		result = append(result, all[idx])
		if len(result) >= max {
			break
		}
	}
	return result
}

// renderMemorySection groups selected project-scoped memories by type and
// renders a Markdown block. Returns "" when empty.
func renderMemorySection(memories []storage.Memory) string {
	if len(memories) == 0 {
		return ""
	}
	var userLines, projectFactLines []string
	for _, m := range memories {
		switch m.Type {
		case MemoryTypeUserPreference:
			userLines = append(userLines, memoryBlock("(喜好) ", m))
		case MemoryTypeEpisodicMemory:
			userLines = append(userLines, memoryBlock("(经历) ", m))
		case MemoryTypeProjectFact:
			projectFactLines = append(projectFactLines, memoryBlock("", m))
		}
	}
	sections := make([]string, 0, 3)
	sections = append(sections, "### 当前项目记忆\n以下记忆仅适用于当前项目；不要迁移到其他项目会话。")
	if len(userLines) > 0 {
		sections = append(sections, "#### 关于用户的长期记忆\n"+strings.Join(userLines, "\n"))
	}
	if len(projectFactLines) > 0 {
		sections = append(sections, "#### 当前项目事实\n"+strings.Join(projectFactLines, "\n"))
	}
	return strings.Join(sections, "\n\n")
}

func memoryBlock(prefix string, m storage.Memory) string {
	line := "- " + prefix + memoryLine(m)
	body := strings.TrimSpace(m.Body)
	if body == "" {
		return line
	}
	return line + "\n  " + strings.ReplaceAll(body, "\n", "\n  ")
}

func memoryLine(m storage.Memory) string {
	if strings.TrimSpace(m.Description) == "" {
		return m.Name
	}
	return m.Name + "：" + m.Description
}

func buildExtractionUserPrompt(existing []storage.Memory, dialogue string) string {
	return "Existing memories:\n" + renderMemoryListForPrompt(existing) + "\n\nDialogue:\n" + dialogue
}

func buildConsolidationUserPrompt(items []storage.Memory) string {
	var b strings.Builder
	for i, m := range items {
		b.WriteString(fmt.Sprintf("[%d] (%s) %s: %s\n", i, m.Type, m.Name, m.Description))
		if strings.TrimSpace(m.Body) != "" {
			b.WriteString("    body: " + m.Body + "\n")
		}
	}
	return b.String()
}

func buildSelectionUserPrompt(all []storage.Memory, dialogue string) string {
	var b strings.Builder
	b.WriteString("Recent conversation:\n")
	b.WriteString(dialogue)
	b.WriteString("\n\nCandidate memories:\n")
	for i, m := range all {
		b.WriteString(fmt.Sprintf("[%d] (%s) %s: %s\n", i, m.Type, m.Name, m.Description))
	}
	return b.String()
}

func renderMemoryListForPrompt(items []storage.Memory) string {
	if len(items) == 0 {
		return "(none)"
	}
	lines := make([]string, 0, len(items))
	for _, m := range items {
		lines = append(lines, fmt.Sprintf("- [%s] %s: %s", m.Type, m.Name, m.Description))
	}
	return strings.Join(lines, "\n")
}

// renderDialogueForMemory builds a plain-text transcript of user and assistant
// messages only, dropping tool noise to keep the LLM calls cheap.
func renderDialogueForMemory(history []storage.Message) string {
	var b strings.Builder
	for _, msg := range history {
		var content string
		switch msg.Role {
		case "user":
			content = strings.TrimSpace(msg.Content)
			if content == "" {
				continue
			}
			b.WriteString("[user] ")
		case "assistant":
			content = strings.TrimSpace(msg.Content)
			if content == "" {
				continue
			}
			b.WriteString("[assistant] ")
		default:
			continue
		}
		b.WriteString(content)
		b.WriteString("\n\n")
		if b.Len() > maxMemoryDialogueChars {
			break
		}
	}
	return b.String()
}

func extractJSONArray(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if start := strings.IndexByte(trimmed, '['); start >= 0 {
		if end := strings.LastIndexByte(trimmed, ']'); end > start {
			return trimmed[start : end+1]
		}
	}
	return ""
}

func truncateRunes(s string, max int) string {
	if runes := []rune(s); len(runes) > max {
		return string(runes[:max])
	}
	return s
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"context\"\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"strings\"\n\n\topenai \"github.com/sashabaranov/go-openai\"\n\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/idgen\"\n\t\"nano_cc/internal/logger\"\n)\n\nconst (\n\tMemoryTypeEpisodicMemory = \"episodic_memory\"\n\tMemoryTypeUserPreference = \"user_preference\"\n\tMemoryTypeProjectFact    = \"project_fact\"\n)\n\nconst (\n\tmaxPreferenceMemories    = 20\n\tmaxSessionSummaries      = 20\n\tsessionSummaryPruneCount = 10\n\tmaxProjectFactMemories   = 50\n\tmaxInjectedMe"}
```

## 2026-06-13T09:16:23Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_FJ9Pd5gaVF9VFQ16Mri66887

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/process.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/process.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/process.go: no such file or directory"}
```

## 2026-06-13T09:16:23Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_j0pR1bWLPYuwT1Osy1xE5656

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/service.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/service.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/service.go: no such file or directory"}
```

## 2026-06-13T09:16:23Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_caUI35t2uYzQc3QIKqHx1518

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/subagent.go"}
```

### Result

```text
package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/idgen"
	"nano_cc/internal/logger"
	"nano_cc/internal/safety"
	agenttools "nano_cc/internal/tools"
)

const defaultSubagentMaxRounds = 20

type spawnSubagentArgs struct {
	Task string `json:"task"`
	CWD  string `json:"cwd"`
}

type subagentContextKey string

const subagentDepthKey subagentContextKey = "subagent_depth"

type subagentTrace struct {
	store            conversationStore
	runID            string
	parentToolCallID string
	conversationID   string
	userID           string
	sequenceNo       int
}

func (s *Service) runSubagent(ctx context.Context, parent ToolContext, args spawnSubagentArgs, audit toolExecutionAudit) (string, error) {
	if depth, _ := ctx.Value(subagentDepthKey).(int); depth > 0 {
		return "", fmt.Errorf("spawn_subagent cannot be called from a subagent")
	}
	task := strings.TrimSpace(args.Task)
	if task == "" {
		return "", fmt.Errorf("task is required")
	}
	resolvedCWD, err := resolveSubagentCWD(s.Cfg.WorkspaceRoot, args.CWD)
	if err != nil {
		return "", err
	}
	runID := idgen.New("subagent")
	parentToolCallID := strings.TrimSpace(parent.ParentToolCallID)
	if parentToolCallID == "" {
		parentToolCallID = "spawn_subagent"
	}
	trace := &subagentTrace{store: s.Store, runID: runID, parentToolCallID: parentToolCallID, conversationID: parent.Conversation.ID, userID: parent.User.ID}
	childTools := NewChildToolRegistry(s.Cfg, resolvedCWD)
	childState := s.newLoopState(parent.Conversation, parent.User, task, nil, nil, nil)
	childState.SkillSnapshot = parent.Skills
	childState.SystemPrompt = s.buildSubagentSystemPrompt(parent.User, parent.Skills)
	childState.Messages = []openai.ChatCompletionMessage{{Role: "system", Content: childState.SystemPrompt}, {Role: "user", Content: task}}
	if err := trace.record(ctx, storage.Message{ID: childState.NextMessageID(), ConversationID: parent.Conversation.ID, UserID: parent.User.ID, Role: "user", Content: task}); err != nil {
		return "", err
	}
	childState.ToolRuntimeEnv = childTools.runtimeEnv
	childCtx := context.WithValue(ctx, subagentDepthKey, 1)
	msg, err := s.runSubagentLoop(childCtx, childState, childTools, parent, trace, defaultSubagentMaxRounds)
	if err != nil {
		return "", err
	}
	return "Subagent completed.\n\nSummary:\n" + fallbackAssistantContent(msg.Content), nil
}

func (s *Service) buildSubagentSystemPrompt(user storage.User, snapshot *agenttools.SkillSnapshot) string {
	base := s.buildSystemPromptWithMemory(user, snapshot, "")
	return base + "\n\n<subagent>\n你是由 `spawn_subagent` 派生出来的子智能体。\n\n规则：\n- 你看不到父对话的历史记录。\n- 只能依据当前任务和工作区文件来工作。\n- 不要调用 `spawn_subagent`。\n- 完成后，只输出一段简洁的摘要，说明你做了什么、关键发现以及尚未解决的问题。\n</subagent>"
}

func (s *Service) runSubagentLoop(ctx context.Context, state *LoopState, tools *ToolRegistry, parent ToolContext, trace *subagentTrace, maxRounds int) (openai.ChatCompletionMessage, error) {
	roundsSinceTodoWrite := 0
	for round := 1; round <= maxRounds; round++ {
		req := openai.ChatCompletionRequest{Model: s.Cfg.LLM.ModelID, Messages: state.Messages, Tools: tools.Definitions(), MaxTokens: defaultMaxTokens}
		reqBody, _ := json.Marshal(req)
		msg, finishReason, err := s.runModelRoundWithRecovery(ctx, state, req)
		respBody, _ := json.Marshal(msg)
		logger.LogLLMRound(round, fmt.Sprintf("subagent run=%s parent_tool_call=%s conversation=%s", trace.runID, trace.parentToolCallID, parent.Conversation.ID), reqBody, respBody, err)
		if err != nil {
			return openai.ChatCompletionMessage{}, err
		}
		state.Messages = append(state.Messages, msg)
		if toolCallsInclude(msg.ToolCalls, agenttools.TodoWriteToolName) {
			roundsSinceTodoWrite = 0
		} else {
			roundsSinceTodoWrite++
		}
		storedAssistant := storage.Message{ID: state.NextMessageID(), ConversationID: parent.Conversation.ID, UserID: parent.User.ID, Role: "assistant", Content: msg.Content, ReasoningContent: msg.ReasoningContent, ToolCalls: openAIToolCallsToStorage(msg.ToolCalls)}
		if err := trace.record(ctx, storedAssistant); err != nil {
			return openai.ChatCompletionMessage{}, err
		}
		if finishReason != "tool_calls" || len(msg.ToolCalls) == 0 {
			return msg, nil
		}
		state.History = append(state.History, storedAssistant)
		for _, tc := range msg.ToolCalls {
			toolCtx := &ToolUseContext{State: state, ToolCall: tc, Name: tc.Function.Name, RawArgs: tc.Function.Arguments}
			if err := s.hookManager().RunPreToolUse(ctx, toolCtx); err != nil {
				return openai.ChatCompletionMessage{}, err
			}
			toolCtx.Outcome = s.executeChildToolCall(ctx, tools, parent, tc.Function.Name, tc.Function.Arguments, toolCtx.Outcome.Audit)
			if toolCtx.Name == agenttools.TodoWriteToolName && toolCtx.Outcome.Status == "success" {
				state.Todos = append([]agenttools.TodoItem(nil), toolCtx.Outcome.Todos...)
			}
			if err := s.hookManager().RunPostToolUse(ctx, toolCtx); err != nil {
				return openai.ChatCompletionMessage{}, err
			}
			if len(state.History) > 0 {
				if err := trace.record(ctx, state.History[len(state.History)-1]); err != nil {
					return openai.ChatCompletionMessage{}, err
				}
			}
		}
		roundsSinceTodoWrite = maybeAppendTodoWriteReminder(state, tools, roundsSinceTodoWrite)
	}
	return openai.ChatCompletionMessage{}, fmt.Errorf("subagent exceeded max rounds")
}

func (s *Service) executeChildToolCall(ctx context.Context, tools *ToolRegistry, toolCtx ToolContext, name string, rawArgs string, audit toolExecutionAudit) toolExecutionOutcome {
	if name == "spawn_subagent" {
		return toolExecutionOutcome{Status: "rejected", Result: "Error: spawn_subagent cannot be called from a subagent", Audit: audit}
	}
	execResult, err := tools.Execute(ctx, toolCtx, name, rawArgs)
	if err != nil {
		return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Error: %v", err), Audit: audit}
	}
	return toolExecutionOutcome{Status: "success", Result: execResult.Output, Audit: audit, Todos: execResult.Todos}
}

func (t *subagentTrace) record(ctx context.Context, msg storage.Message) error {
	if t == nil || t.store == nil {
		return nil
	}
	t.sequenceNo++
	return t.store.CreateSubagentMessage(ctx, storage.SubagentMessage{ID: idgen.New("submsg"), RunID: t.runID, ParentToolCallID: t.parentToolCallID, ConversationID: t.conversationID, UserID: t.userID, SequenceNo: t.sequenceNo, Role: msg.Role, Content: msg.Content, ReasoningContent: msg.ReasoningContent, ToolCallID: msg.ToolCallID, ToolCalls: msg.ToolCalls})
}

func resolveSubagentCWD(workspaceRoot, cwd string) (string, error) {
	root, err := filepath.Abs(strings.TrimSpace(workspaceRoot))
	if err != nil || strings.TrimSpace(workspaceRoot) == "" {
		return "", fmt.Errorf("workspace root is required")
	}
	root = filepath.Clean(root)
	resolved := root
	if strings.TrimSpace(cwd) != "" {
		resolved = strings.TrimSpace(cwd)
		if !filepath.IsAbs(resolved) {
			resolved = filepath.Join(root, resolved)
		}
		resolved, err = filepath.Abs(resolved)
		if err != nil {
			return "", fmt.Errorf("resolve subagent cwd: %w", err)
		}
		resolved = filepath.Clean(resolved)
	}
	if !safety.Contains(root, resolved) {
		return "", fmt.Errorf("subagent cwd escapes workspace: %s", cwd)
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", fmt.Errorf("subagent cwd is unavailable: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("subagent cwd is not a directory: %s", cwd)
	}
	return resolved, nil
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"context\"\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"os\"\n\t\"path/filepath\"\n\t\"strings\"\n\n\topenai \"github.com/sashabaranov/go-openai\"\n\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/idgen\"\n\t\"nano_cc/internal/logger\"\n\t\"nano_cc/internal/safety\"\n\tagenttools \"nano_cc/internal/tools\"\n)\n\nconst defaultSubagentMaxRounds = 20\n\ntype spawnSubagentArgs struct {\n\tTask string `json:\"task\"`\n\tCWD  string `json:\"cwd\"`\n}\n\ntype subagentContextKey string\n\nconst subagentDepthKey subagentContextKey = \"sub"}
```

## 2026-06-13T09:16:23Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_03_AAOx0aqMzAv4kk4oSeed5444

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/style.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/style.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/style.go: no such file or directory"}
```

## 2026-06-13T09:16:23Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_04_uIIh9qhSANC8vGNcnIhU4145

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/handler.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/handler.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/handler.go: no such file or directory"}
```

## 2026-06-13T09:16:27Z · bash · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_9puCCeMbVJrBuUiQnzyI2560

### Arguments

```json
{"command": "find /Users/bytedance/golang_pro/nano_cc/go-agent -type f -name \"*.go\" -not -path \"*/vendor/*\" | sort"}
```

### Result

```text
/Users/bytedance/golang_pro/nano_cc/go-agent/architecture_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/transport_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/transport.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/compression_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/conversation_memory_strategy.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/full_history_summarization.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/message_window_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/reactive_compact.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/recent_tool_result_retention.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/token_estimator.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/tool_result_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/context_compression_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/context_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/conversation_flow.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/conversation_memory.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/helpers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hook_bridge.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/manager.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/stop.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/tool.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/types.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/user_prompt.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/memory_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/memory.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/prompt_builder.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/runtime_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/runtime.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/subagent.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/summarizer.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/todo_reminder.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/tool_registry.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/conversation_history_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/conversation_history.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/models.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/assistant/prompt_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/assistant/prompt.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/cli/root_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/cli/root.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/helpers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/local_config_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/local_config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/paths.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/idgen/idgen.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/llm/client.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/llm/errors.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/bootstrap.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/memory_store.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/persisted_output_files.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/session_history_store.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/store_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/store.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/tool_result_log.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/logger/logger_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/logger/logger.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/safety/path.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/sessions/render_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/sessions/skill_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/sessions/skill.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/textutil/textutil.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/bash.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/definitions_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/definitions.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/file_ops.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/handlers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/load_skill.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/path_guard.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/persisted_output_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/persisted_output.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/runtime_env.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/runtime_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/todo_write_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/todo_write.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/validation_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/validation.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/app.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/events_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/events.go
/Users/bytedance/golang_pro/nano_cc/go-agent/main.go
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/basic.go
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","resolved_command_path":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"/Users/bytedance/golang_pro/nano_cc/go-agent/architecture_test.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config_test.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager_test.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/transport_test.go\n/Users/bytedance/golang_pro/nano_cc/go-agent"}
```

## 2026-06-13T09:16:31Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_hlqJvdKRETZbfq9JttW88082

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/conversation_flow.go"}
```

### Result

```text
package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/runtime/compression"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/idgen"
	"nano_cc/internal/llm"
	"nano_cc/internal/logger"
	agenttools "nano_cc/internal/tools"
)

const (
	assistantDeltaEvent = "assistant_delta"
	reasoningDeltaEvent = "reasoning_delta"
	maxRound            = 50

	// defaultMaxTokens is the per-request output budget; on the first
	// truncation it is upgraded to truncationMaxTokens.
	defaultMaxTokens = 8000
	// truncationMaxTokens is the upgraded output budget after the first
	// truncation (8x the default).
	truncationMaxTokens = 64 * 1024
	// maxResumeAttempts is how many continuation requests are issued after the
	// upgraded budget still truncates.
	maxResumeAttempts = 3

	// truncationResumePrompt is injected as a user message to make the model
	// continue an output that was cut off by the token limit.
	truncationResumePrompt = `Output token limit hit. Resume directly — no apology, no recap of what you were doing. Pick up mid-thought if that is where the cut happened. Break remaining work into smaller pieces.`
)

type bufferedModelDelta struct {
	Event   string
	Content string
}

func (s *Service) RespondToConversation(ctx context.Context, conversation storage.Conversation, user storage.User, userMessage string, writer EventWriter) (storage.Message, error) {
	// 记忆业务开关：系统级能力开启 且 用户个人偏好开启 时才注入/提取记忆。
	// 注意：它只控制记忆注入与提取，不控制会话锁与模型历史持久化。
	memoryOn := s.EnableMemory && user.MemoryEnabled
	// 入口获取会话锁：上一轮收尾未完成时阻塞等待，直到拿到锁或等待超时。
	// 获取失败（本地锁异常 / 等待超时）则降级放行，跳过本轮收尾。
	var lockToken string
	var stopRenew func()
	handedOff := false
	if s.EnableMemory {
		token := idgen.New("lock")
		ok, err := s.Store.AcquireConversationLock(ctx, conversation.ID, token, s.Cfg.ConversationLockTTL, s.Cfg.ConversationLockWaitTimeout)
		if err != nil {
			logger.Warn(fmt.Sprintf("conversation lock: acquire failed conversation=%s: %v", conversation.ID, err))
		} else if !ok {
			logger.Warn(fmt.Sprintf("conversation lock: acquire timed out conversation=%s", conversation.ID))
		} else {
			lockToken = token
			stopRenew = s.startLockRenewer(conversation.ID, token)
		}
	}
	// defer 兜底：任何提前 return（出错路径）都停止续期并释放锁；
	// 正常走到异步收尾时通过 handedOff 转交所有权，避免被误释放。
	defer func() {
		if handedOff || lockToken == "" {
			return
		}
		if stopRenew != nil {
			stopRenew()
		}
		_ = s.Store.ReleaseConversationLock(context.Background(), conversation.ID, lockToken)
	}()

	history, err := s.loadConversationMessages(ctx, conversation.ID)
	if err != nil {
		return storage.Message{}, err
	}
	modelHistory := s.loadModelHistory(ctx, conversation.ID, history)
	state := s.newLoopState(conversation, user, userMessage, history, modelHistory, writer)
	if err := s.hookManager().RunUserPromptSubmit(ctx, &UserPromptSubmitContext{State: state}); err != nil {
		return storage.Message{}, err
	}

	snapshot, err := s.buildSkillSnapshot(ctx, user.ID)
	if err != nil {
		return storage.Message{}, err
	}
	state.SkillSnapshot = snapshot
	state.SystemPrompt = s.buildSystemPrompt(ctx, conversation, user, snapshot, state.History, memoryOn)
	state.Messages = buildOpenAIMessages(state.SystemPrompt, state.ModelHistory)
	toolDefs := s.toolDefinitionsForUser(ctx, user.ID)
	round := 0
	roundsSinceTodoWrite := 0
	var cumulativeReasoning strings.Builder
	var lastRequestHistory []storage.Message

	for {
		round++
		if round > maxRound {
			return storage.Message{}, fmt.Errorf("conversation round limit %d exceeded", maxRound)
		}
		requestHistory, err := s.compressContextBeforeLLM(ctx, state)
		if err != nil {
			return storage.Message{}, err
		}
		lastRequestHistory = requestHistory
		state.Messages = buildOpenAIMessages(state.SystemPrompt, requestHistory)
		roundsSinceTodoWrite = maybeAppendTodoWriteReminder(state, s.Tools, roundsSinceTodoWrite)
		estimator := compression.DefaultTokenEstimator{}
		state.LastContextTokens = estimator.EstimateRequestTokens(state.SystemPrompt, requestHistory, toolDefs)
		state.LastContextBudget = estimator.ContextTokenBudget()
		emitMeta(state)
		req := openai.ChatCompletionRequest{
			Model:     s.Cfg.LLM.ModelID,
			Messages:  state.Messages,
			Tools:     toolDefs,
			MaxTokens: defaultMaxTokens,
		}
		reqBody, _ := json.Marshal(req)
		msg, finishReason, err := s.runModelRoundWithRecovery(ctx, state, req)
		respBody, _ := json.Marshal(msg)
		logger.LogLLMRound(round, fmt.Sprintf("main-agent conversation=%s", conversation.ID), reqBody, respBody, err)
		if err != nil {
			return storage.Message{}, err
		}
		cumulativeReasoning.WriteString(msg.ReasoningContent)
		requestMsg := msg
		state.Messages = append(state.Messages, requestMsg)
		if toolCallsInclude(msg.ToolCalls, agenttools.TodoWriteToolName) {
			roundsSinceTodoWrite = 0
		} else {
			roundsSinceTodoWrite++
		}

		if finishReason != "tool_calls" || len(msg.ToolCalls) == 0 {
			// 会话结束：把本轮模型最终回复纳入上下文后重算，得到会话结束时的最终 token 用量，
			// 覆盖请求前的估算值，确保存储与下发的都是最终用量。
			finalAssistant := storage.Message{Role: "assistant", Content: msg.Content, ReasoningContent: msg.ReasoningContent, ToolCalls: openAIToolCallsToStorage(msg.ToolCalls)}
			finalHistoryForEstimate := append(cloneMessages(lastRequestHistory), finalAssistant)
			finalEstimator := compression.DefaultTokenEstimator{}
			state.LastContextTokens = finalEstimator.EstimateRequestTokens(state.SystemPrompt, finalHistoryForEstimate, toolDefs)
			state.LastContextBudget = finalEstimator.ContextTokenBudget()
			stopCtx := &StopContext{State: state, ModelMessage: msg, Content: fallbackAssistantContent(msg.Content), ReasoningContent: cumulativeReasoning.String()}
			if err := s.hookManager().RunStop(ctx, stopCtx); err != nil {
				return storage.Message{}, err
			}
			if s.EnableMemory {
				finalHistory := append(state.History, stopCtx.AssistantMessage)
				finalModelHistory := append(cloneMessages(lastRequestHistory), stopCtx.AssistantMessage)
				handedOff = s.scheduleMemoryWork(conversation, user, finalHistory, finalModelHistory, lockToken, stopRenew, memoryOn)
			}
			return stopCtx.AssistantMessage, nil
		}
		assistantMessage := storage.Message{ID: state.NextMessageID(), ConversationID: conversation.ID, UserID: user.ID, Role: "assistant", Content: msg.Content, ReasoningContent: msg.ReasoningContent, ToolCalls: openAIToolCallsToStorage(msg.ToolCalls)}
		state.History = append(state.History, assistantMessage)
		state.ModelHistory = append(state.ModelHistory, assistantMessage)

		for _, tc := range msg.ToolCalls {
			toolCtx := &ToolUseContext{State: state, ToolCall: tc, Name: tc.Function.Name, RawArgs: tc.Function.Arguments}
			if err := s.hookManager().RunPreToolUse(ctx, toolCtx); err != nil {
				return storage.Message{}, err
			}
			toolCtx.Outcome = s.executeToolCall(ctx, ToolContext{User: user, Conversation: conversation, Skills: snapshot, ParentToolCallID: tc.ID, PersistedOutputReader: s.newPersistedOutputReader(conversation.ID, user.ID)}, tc.Function.Name, tc.Function.Arguments, toolCtx.Outcome.Audit)
			if toolCtx.Name == agenttools.TodoWriteToolName && toolCtx.Outcome.Status == "success" {
				state.Todos = append([]agenttools.TodoItem(nil), toolCtx.Outcome.Todos...)
			}
			if err := s.hookManager().RunPostToolUse(ctx, toolCtx); err != nil {
				return storage.Message{}, err
			}
			state.ToolCallCount++
			emitMeta(state)
		}
	}
}

// emitMeta 实时下发当前累计的回复元信息。
func emitMeta(state *LoopState) {
	if state == nil || state.Writer == nil {
		return
	}
	_ = state.Writer.Event("meta", map[string]any{
		"tool_call_count": state.ToolCallCount,
		"context_tokens":  state.LastContextTokens,
		"context_budget":  state.LastContextBudget,
	})
}

// runModelRoundWithRecovery wraps runModelRoundStream with truncation and
// context-overflow recovery. Transient 429/529 retries are handled inside the
// LLM client and are transparent here.
//
//   - Truncation (finish_reason == length): the first time, the output budget is
//     upgraded to truncationMaxTokens and the same request is retried WITHOUT
//     touching messages (the partial output is discarded). If it still
//     truncates, up to maxResumeAttempts continuation requests are issued, each
//     appending the truncated text plus a resume prompt; segments are streamed
//     and concatenated into one assistant message.
//   - Context overflow (HTTP 413): reactiveCompact is run once and the request
//     rebuilt from the compacted history, then retried. A second overflow is
//     returned to the caller (handled by the existing fallback boundary).
func (s *Service) runModelRoundWithRecovery(ctx context.Context, state *LoopState, req openai.ChatCompletionRequest) (openai.ChatCompletionMessage, openai.FinishReason, error) {
	cur := req
	upgraded := false
	compacted := false
	resumeAttempts := 0
	var accumulated strings.Builder
	var accumulatedReasoning strings.Builder

	for {
		msg, finishReason, deltas, err := s.runModelRoundStream(ctx, state, cur)
		if err != nil {
			if llm.IsContextOverflow(err) && !compacted {
				if compactErr := s.reactiveCompact(ctx, state); compactErr != nil {
					logger.Warn(fmt.Sprintf("reactive compact failed conversation=%s: %v", state.Conversation.ID, compactErr))
					return openai.ChatCompletionMessage{}, "", err
				}
				compacted = true
				cur.Messages = state.Messages
				cur.MaxTokens = defaultMaxTokens
				upgraded = false
				resumeAttempts = 0
				accumulated.Reset()
				accumulatedReasoning.Reset()
				continue
			}
			return openai.ChatCompletionMessage{}, "", err
		}

		resuming := accumulated.Len() > 0
		if finishReason != openai.FinishReasonLength {
			if resuming {
				flushContentDeltas(state, deltas)
				return openai.ChatCompletionMessage{
					Role:             "assistant",
					Content:          accumulated.String() + msg.Content,
					ReasoningContent: accumulatedReasoning.String() + msg.ReasoningContent,
					ToolCalls:        msg.ToolCalls,
				}, finishReason, nil
			}
			if state.Writer != nil && shouldEmitAssistantContentDeltas(finishReason, msg.ToolCalls) {
				flushContentDeltas(state, deltas)
			}
			return msg, finishReason, nil
		}

		// Truncated output.
		if !upgraded {
			// Discard the partial output, upgrade the budget, retry unchanged.
			upgraded = true
			cur.MaxTokens = truncationMaxTokens
			continue
		}
		if resumeAttempts >= maxResumeAttempts {
			// Best-effort: emit and return what we have so far.
			flushContentDeltas(state, deltas)
			logger.Warn(fmt.Sprintf("truncation resume exhausted conversation=%s after %d attempts", state.Conversation.ID, resumeAttempts))
			return openai.ChatCompletionMessage{
				Role:             "assistant",
				Content:          accumulated.String() + msg.Content,
				ReasoningContent: accumulatedReasoning.String() + msg.ReasoningContent,
			}, finishReason, nil
		}
		// Keep this partial segment and ask the model to continue.
		flushContentDeltas(state, deltas)
		accumulated.WriteString(msg.Content)
		accumulatedReasoning.WriteString(msg.ReasoningContent)
		resumeAttempts++
		cur.Messages = append(cur.Messages,
			openai.ChatCompletionMessage{Role: "assistant", Content: msg.Content},
			openai.ChatCompletionMessage{Role: "user", Content: truncationResumePrompt},
		)
	}
}

func flushContentDeltas(state *LoopState, deltas []bufferedModelDelta) {
	if state.Writer == nil {
		return
	}
	for _, delta := range deltas {
		_ = state.Writer.Event(delta.Event, map[string]any{"content": delta.Content})
	}
}

func (s *Service) runModelRoundStream(ctx context.Context, state *LoopState, req openai.ChatCompletionRequest) (openai.ChatCompletionMessage, openai.FinishReason, []bufferedModelDelta, error) {
	stream, err := s.LLM.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return openai.ChatCompletionMessage{}, "", nil, err
	}
	defer stream.Close()

	var content strings.Builder
	var reasoningContent strings.Builder
	var finishReason openai.FinishReason
	toolCalls := &streamedToolCallAccumulator{}
	var bufferedContentDeltas []bufferedModelDelta
	seenChoice := false
	seenOutput := false

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return openai.ChatCompletionMessage{}, "", nil, err
		}
		if len(chunk.Choices) == 0 {
			continue
		}
		seenChoice = true
		choice := chunk.Choices[0]
		if choice.Delta.Content != "" {
			seenOutput = true
			content.WriteString(choice.Delta.Content)
			bufferedContentDeltas = append(bufferedContentDeltas, bufferedModelDelta{Event: assistantDeltaEvent, Content: choice.Delta.Content})
		}
		if choice.Delta.ReasoningContent != "" {
			seenOutput = true
			reasoningContent.WriteString(choice.Delta.ReasoningContent)
			if state.Writer != nil {
				_ = state.Writer.Event(reasoningDeltaEvent, map[string]any{"content": choice.Delta.ReasoningContent})
			}
		}
		if len(choice.Delta.ToolCalls) > 0 {
			seenOutput = true
			toolCalls.Add(choice.Delta.ToolCalls)
		}
		if choice.FinishReason != "" {
			finishReason = choice.FinishReason
		}
	}
	if !seenChoice || (!seenOutput && finishReason == "") {
		return openai.ChatCompletionMessage{}, "", nil, fmt.Errorf("model stream returned no choices")
	}
	calls := toolCalls.Calls()

	return openai.ChatCompletionMessage{Role: "assistant", Content: content.String(), ReasoningContent: reasoningContent.String(), ToolCalls: calls}, finishReason, bufferedContentDeltas, nil
}

func shouldEmitAssistantContentDeltas(finishReason openai.FinishReason, toolCalls []openai.ToolCall) bool {
	return finishReason != openai.FinishReasonToolCalls && len(toolCalls) == 0
}

type streamedToolCallAccumulator struct {
	calls []openai.ToolCall
}

func (a *streamedToolCallAccumulator) Add(deltas []openai.ToolCall) {
	for deltaPosition, delta := range deltas {
		index := len(a.calls) - 1
		if delta.Index != nil {
			index = *delta.Index
		} else if len(deltas) > 1 {
			index = deltaPosition
		}
		if index < 0 {
			index = 0
		}
		for len(a.calls) <= index {
			a.calls = append(a.calls, openai.ToolCall{})
		}

		call := a.calls[index]
		if delta.ID != "" {
			call.ID = delta.ID
		}
		if delta.Type != "" {
			call.Type = delta.Type
		}
		if delta.Function.Name != "" {
			call.Function.Name = delta.Function.Name
		}
		if delta.Function.Arguments != "" {
			call.Function.Arguments += delta.Function.Arguments
		}
		a.calls[index] = call
	}
}

func (a *streamedToolCallAccumulator) Calls() []openai.ToolCall {
	if len(a.calls) == 0 {
		return nil
	}
	return append([]openai.ToolCall(nil), a.calls...)
}

func (s *Service) loadConversationMessages(ctx context.Context, conversationID string) ([]storage.Message, error) {
	if cached, ok, err := s.Store.GetConversationCache(ctx, conversationID); err == nil && ok {
		return cached, nil
	}
	messages, err := s.Store.ListMessagesByConversation(ctx, conversationID, 100)
	if err != nil {
		return nil, err
	}
	if err := s.Store.SetConversationCache(ctx, conversationID, messages); err != nil {
		_ = err
	}
	return messages, nil
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"context\"\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"io\"\n\t\"strings\"\n\n\topenai \"github.com/sashabaranov/go-openai\"\n\n\t\"nano_cc/internal/agent/runtime/compression\"\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/idgen\"\n\t\"nano_cc/internal/llm\"\n\t\"nano_cc/internal/logger\"\n\tagenttools \"nano_cc/internal/tools\"\n)\n\nconst (\n\tassistantDeltaEvent = \"assistant_delta\"\n\treasoningDeltaEvent = \"reasoning_delta\"\n\tmaxRound            = 50\n\n\t// defaultMaxTokens is the per-request output budget; on the "}
```

## 2026-06-13T09:16:31Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_ThrxVZf0c75Bz6PwXvdf7311

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/helpers.go"}
```

### Result

```text
package runtime

import (
	"strings"

	"nano_cc/internal/idgen"
)

func newMessageID() string { return idgen.New("msg") }

func fallbackAssistantContent(content string) string {
	if strings.TrimSpace(content) == "" {
		return "(no response)"
	}
	return content
}

func shouldInferConversationTitle(currentTitle string) bool {
	trimmed := strings.TrimSpace(currentTitle)
	return trimmed == "" || trimmed == "新对话"
}

func inferConversationTitle(userMessage string) string {
	trimmed := strings.TrimSpace(userMessage)
	if len([]rune(trimmed)) > 30 {
		return string([]rune(trimmed)[:30])
	}
	if trimmed == "" {
		return "新对话"
	}
	return trimmed
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"strings\"\n\n\t\"nano_cc/internal/idgen\"\n)\n\nfunc newMessageID() string { return idgen.New(\"msg\") }\n\nfunc fallbackAssistantContent(content string) string {\n\tif strings.TrimSpace(content) == \"\" {\n\t\treturn \"(no response)\"\n\t}\n\treturn content\n}\n\nfunc shouldInferConversationTitle(currentTitle string) bool {\n\ttrimmed := strings.TrimSpace(currentTitle)\n\treturn trimmed == \"\" || trimmed == \"新对话\"\n}\n\nfunc inferConversationTitle(userMessage string) string {\n\ttrimmed := strings.Tr"}
```

## 2026-06-13T09:16:31Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_6NXCNanfqWAN5RuYKbHO6671

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/summarizer.go"}
```

### Result

```text
package runtime

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/runtime/compression"
	"nano_cc/internal/agent/storage"
)

const summarizerSystemPrompt = `You are a context compaction engine. Summarize the conversation so far into a structured Markdown brief that another AI agent can use to continue the task without losing critical information.

Rules:
- Do NOT call tools or invent information that is not present.
- Preserve concrete details: current user goal, completed actions and key conclusions, important file paths, function names, commands, errors, and decisions.
- For any <persisted-output ...> marker you see, keep its id and note that read_persisted_output can fetch the full content.
- If a tool result shows "[Earlier result compacted. Re-run if needed]", note that it must be re-run if needed.
- List unfinished items, things to verify, and next steps.

Output only the Markdown summary.`

// summarizeHistoryForContext runs a single non-tool LLM call to summarize the
// request-state history. It never mutates state.History.
func (s *Service) summarizeHistoryForContext(ctx context.Context, req compression.SummaryRequest) (compression.SummaryResult, error) {
	if s.LLM == nil {
		return compression.SummaryResult{}, fmt.Errorf("llm client is not configured")
	}
	messages := []openai.ChatCompletionMessage{
		{Role: "system", Content: summarizerSystemPrompt},
		{Role: "user", Content: renderHistoryForSummary(req.History)},
	}
	resp, err := s.LLM.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    s.Cfg.LLM.ModelID,
		Messages: messages,
	})
	if err != nil {
		return compression.SummaryResult{}, err
	}
	if len(resp.Choices) == 0 {
		return compression.SummaryResult{}, fmt.Errorf("summary response had no choices")
	}
	summary := strings.TrimSpace(resp.Choices[0].Message.Content)
	if summary == "" {
		return compression.SummaryResult{}, fmt.Errorf("summary response was empty")
	}
	return compression.SummaryResult{Summary: summary}, nil
}

func renderHistoryForSummary(history []storage.Message) string {
	var b strings.Builder
	b.WriteString("Summarize the following conversation history:\n\n")
	for _, msg := range history {
		role := msg.Role
		switch role {
		case "tool":
			b.WriteString("[tool result] ")
			b.WriteString(msg.Content)
		case "assistant":
			b.WriteString("[assistant] ")
			if strings.TrimSpace(msg.Content) != "" {
				b.WriteString(msg.Content)
			}
			for _, call := range msg.ToolCalls {
				fmt.Fprintf(&b, "\n[tool call] %s(%s)", call.Function.Name, call.Function.Arguments)
			}
		default:
			fmt.Fprintf(&b, "[%s] %s", role, msg.Content)
		}
		b.WriteString("\n\n")
	}
	return b.String()
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"context\"\n\t\"fmt\"\n\t\"strings\"\n\n\topenai \"github.com/sashabaranov/go-openai\"\n\n\t\"nano_cc/internal/agent/runtime/compression\"\n\t\"nano_cc/internal/agent/storage\"\n)\n\nconst summarizerSystemPrompt = `You are a context compaction engine. Summarize the conversation so far into a structured Markdown brief that another AI agent can use to continue the task without losing critical information.\n\nRules:\n- Do NOT call tools or invent information that is not present.\n- Preserve concrete d"}
```

## 2026-06-13T09:16:31Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_03_V3tuscVlVWje7uF0RMVW1602

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/todo_reminder.go"}
```

### Result

```text
package runtime

import (
	openai "github.com/sashabaranov/go-openai"

	agenttools "nano_cc/internal/tools"
)

const (
	todoWriteReminderThreshold = 3
	todoWriteReminderText      = "<system-reminder>\nYou have not called todo_write for 3 consecutive model rounds. If the task is multi-step or your plan has changed, call todo_write to create or update the current task plan before continuing. If todo_write is unnecessary for this simple step, continue normally.\n</system-reminder>"
)

func toolCallsInclude(toolCalls []openai.ToolCall, name string) bool {
	for _, call := range toolCalls {
		if call.Function.Name == name {
			return true
		}
	}
	return false
}

func todoWriteReminderMessage() openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{Role: "system", Content: todoWriteReminderText}
}

func maybeAppendTodoWriteReminder(state *LoopState, tools *ToolRegistry, roundsSinceTodoWrite int) int {
	if roundsSinceTodoWrite < todoWriteReminderThreshold || tools == nil || !tools.isAllowed(agenttools.TodoWriteToolName) {
		return roundsSinceTodoWrite
	}
	state.Messages = append(state.Messages, todoWriteReminderMessage())
	return 0
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\topenai \"github.com/sashabaranov/go-openai\"\n\n\tagenttools \"nano_cc/internal/tools\"\n)\n\nconst (\n\ttodoWriteReminderThreshold = 3\n\ttodoWriteReminderText      = \"\u003csystem-reminder\u003e\\nYou have not called todo_write for 3 consecutive model rounds. If the task is multi-step or your plan has changed, call todo_write to create or update the current task plan before continuing. If todo_write is unnecessary for this simple step, continue normally.\\n\u003c/system-reminder\u003e\"\n)\n\nfunc toolCall"}
```

## 2026-06-13T09:16:31Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_04_ciojf0wzUx0FeARhmMAq9829

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/context_compression.go"}
```

### Result

```text
package runtime

import (
	"context"
	"fmt"

	"nano_cc/internal/agent/runtime/compression"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/logger"
	agenttools "nano_cc/internal/tools"
)

// loadModelHistory returns the reusable "model history" persisted from the
// previous turn's compressed request history. When no row exists, decoding
// fails, or the store errors, it falls back to a clone of the full display
// history (current behavior).
func (s *Service) loadModelHistory(ctx context.Context, conversationID string, displayHistory []storage.Message) []storage.Message {
	modelHistory, ok, err := s.Store.GetConversationModelHistory(ctx, conversationID)
	if err != nil {
		logger.Warn(fmt.Sprintf("model history: load failed conversation=%s: %v", conversationID, err))
	}
	if !ok || len(modelHistory) == 0 {
		return cloneMessages(displayHistory)
	}
	return modelHistory
}

// compressContextBeforeLLM deep-copies the model history and runs the
// compression pipeline, returning the request-only history for this round.
func (s *Service) compressContextBeforeLLM(ctx context.Context, state *LoopState) ([]storage.Message, error) {
	requestHistory := cloneMessages(state.ModelHistory)
	store, ok := s.Store.(compression.Store)
	if !ok {
		// Store does not support compression artifacts; skip silently.
		return requestHistory, nil
	}
	compressor := s.ContextCompressor
	if compressor == nil {
		compressor = compression.NewDefaultCompressor()
	}
	req := &compression.Request{
		Conversation:   state.Conversation,
		User:           state.User,
		RequestHistory: requestHistory,
		SystemPrompt:   state.SystemPrompt,
		Tools:          s.Tools.Definitions(),
		Store:          store,
		Estimator:      compression.DefaultTokenEstimator{},
		Summarizer:     s.summarizeHistoryForContext,
	}
	if err := compressor.Compress(ctx, req); err != nil {
		return nil, err
	}
	return req.RequestHistory, nil
}

// reactiveCompact runs the aggressive ReactiveCompactStrategy out-of-band when
// the LLM rejects a request with HTTP 413 (context overflow). On success it
// updates both state.Messages (effective this round) and state.ModelHistory
// (the new baseline reused by later rounds), but never state.History (the
// verbatim display history). On failure state is left untouched.
func (s *Service) reactiveCompact(ctx context.Context, state *LoopState) error {
	store, ok := s.Store.(compression.Store)
	if !ok {
		return fmt.Errorf("reactive compact: store does not support compression artifacts")
	}
	requestHistory := cloneMessages(state.ModelHistory)
	req := &compression.Request{
		Conversation:   state.Conversation,
		User:           state.User,
		RequestHistory: requestHistory,
		SystemPrompt:   state.SystemPrompt,
		Tools:          s.Tools.Definitions(),
		Store:          store,
		Estimator:      compression.DefaultTokenEstimator{},
		Summarizer:     s.summarizeHistoryForContext,
	}
	if err := (&compression.ReactiveCompactStrategy{}).Apply(ctx, req); err != nil {
		return err
	}
	state.ModelHistory = req.RequestHistory
	state.Messages = buildOpenAIMessages(state.SystemPrompt, req.RequestHistory)
	return nil
}

func cloneMessages(messages []storage.Message) []storage.Message {
	if len(messages) == 0 {
		return nil
	}
	cloned := make([]storage.Message, len(messages))
	for i, msg := range messages {
		cloned[i] = msg
		if len(msg.ToolCalls) > 0 {
			calls := make([]storage.MessageToolCall, len(msg.ToolCalls))
			copy(calls, msg.ToolCalls)
			cloned[i].ToolCalls = calls
		}
	}
	return cloned
}

// persistedOutputReader adapts the storage layer to the tool-facing reader,
// enforcing conversation/user scoping.
type persistedOutputReader struct {
	store          compressionReaderStore
	userID         string
	conversationID string
}

type compressionReaderStore interface {
	GetPersistedOutputForConversation(ctx context.Context, id, userID, conversationID string) (storage.PersistedOutput, error)
}

func (r persistedOutputReader) ReadPersistedOutput(ctx context.Context, id string) (agenttools.PersistedOutput, error) {
	output, err := r.store.GetPersistedOutputForConversation(ctx, id, r.userID, r.conversationID)
	if err != nil {
		return agenttools.PersistedOutput{}, err
	}
	return agenttools.PersistedOutput{
		ID:            output.ID,
		Kind:          output.Kind,
		ToolCallID:    output.ToolCallID,
		OriginalBytes: output.OriginalBytes,
		ContentSHA256: output.ContentSHA256,
		Content:       output.Content,
	}, nil
}

func (s *Service) newPersistedOutputReader(conversationID, userID string) agenttools.PersistedOutputReader {
	store, ok := s.Store.(compressionReaderStore)
	if !ok {
		return nil
	}
	return persistedOutputReader{store: store, userID: userID, conversationID: conversationID}
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"context\"\n\t\"fmt\"\n\n\t\"nano_cc/internal/agent/runtime/compression\"\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/logger\"\n\tagenttools \"nano_cc/internal/tools\"\n)\n\n// loadModelHistory returns the reusable \"model history\" persisted from the\n// previous turn's compressed request history. When no row exists, decoding\n// fails, or the store errors, it falls back to a clone of the full display\n// history (current behavior).\nfunc (s *Service) loadModelHistory(ctx context.Co"}
```

## 2026-06-13T09:16:31Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_05_nhfSEHf7tizelyHCjqZ28731

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/service.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/service.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/service.go: no such file or directory"}
```

## 2026-06-13T09:16:42Z · bash · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_hMXqYw0YTEmuK0a2XjXv6456

### Arguments

```json
{"command": "ls -la /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/"}
```

### Result

```text
total 400
drwxr-xr-x@ 19 bytedance  staff    608  6 12 13:34 .
drwxr-xr-x@  5 bytedance  staff    160  6 12 12:42 ..
drwxr-xr-x@ 11 bytedance  staff    352  6 12 12:50 compression
-rw-r--r--@  1 bytedance  staff   4023  6 12 12:50 context_compression_test.go
-rw-r--r--@  1 bytedance  staff   4740  6 12 12:50 context_compression.go
-rw-r--r--@  1 bytedance  staff  15481  6 12 16:18 conversation_flow.go
-rw-r--r--@  1 bytedance  staff   7397  6 12 12:50 conversation_memory.go
-rw-r--r--@  1 bytedance  staff    656  6 12 16:05 helpers.go
-rw-r--r--@  1 bytedance  staff   1583  6 12 12:50 hook_bridge.go
drwxr-xr-x@  7 bytedance  staff    224  6 12 12:50 hooks
-rw-r--r--@  1 bytedance  staff   8511  6 12 14:26 memory_test.go
-rw-r--r--@  1 bytedance  staff  15294  6 12 14:28 memory.go
-rw-r--r--@  1 bytedance  staff   3337  6 12 16:06 prompt_builder.go
-rw-r--r--@  1 bytedance  staff  89861  6 12 16:04 runtime_test.go
-rw-r--r--@  1 bytedance  staff   4294  6 12 16:05 runtime.go
-rw-r--r--@  1 bytedance  staff   7665  6 12 12:50 subagent.go
-rw-r--r--@  1 bytedance  staff   2729  6 12 12:42 summarizer.go
-rw-r--r--@  1 bytedance  staff   1157  6  6 00:45 todo_reminder.go
-rw-r--r--@  1 bytedance  staff   8754  6 12 16:10 tool_registry.go
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","resolved_command_path":"/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime","outcome_summary":"total 400\ndrwxr-xr-x@ 19 bytedance  staff    608  6 12 13:34 .\ndrwxr-xr-x@  5 bytedance  staff    160  6 12 12:42 ..\ndrwxr-xr-x@ 11 bytedance  staff    352  6 12 12:50 compression\n-rw-r--r--@  1 bytedance  staff   4023  6 12 12:50 context_compression_test.go\n-rw-r--r--@  1 bytedance  staff   4740  6 12 12:50 context_compression.go\n-rw-r--r--@  1 bytedance  staff  15481  6 12 16:18 conversation_flow.go\n-rw-r--r--@  1 bytedance  staff   7397  6 12 12:50 conversation_memory.go\n-rw-r--r--@  1 byteda"}
```

## 2026-06-13T09:16:46Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_E4HAb3I83J8arvYNGORt1627

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/process.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/process.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/process.go: no such file or directory"}
```

## 2026-06-13T09:16:46Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_RHq2FJM7O2oupRzeS4hP7498

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/conversation_flow.go"}
```

### Result

```text
package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/runtime/compression"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/idgen"
	"nano_cc/internal/llm"
	"nano_cc/internal/logger"
	agenttools "nano_cc/internal/tools"
)

const (
	assistantDeltaEvent = "assistant_delta"
	reasoningDeltaEvent = "reasoning_delta"
	maxRound            = 50

	// defaultMaxTokens is the per-request output budget; on the first
	// truncation it is upgraded to truncationMaxTokens.
	defaultMaxTokens = 8000
	// truncationMaxTokens is the upgraded output budget after the first
	// truncation (8x the default).
	truncationMaxTokens = 64 * 1024
	// maxResumeAttempts is how many continuation requests are issued after the
	// upgraded budget still truncates.
	maxResumeAttempts = 3

	// truncationResumePrompt is injected as a user message to make the model
	// continue an output that was cut off by the token limit.
	truncationResumePrompt = `Output token limit hit. Resume directly — no apology, no recap of what you were doing. Pick up mid-thought if that is where the cut happened. Break remaining work into smaller pieces.`
)

type bufferedModelDelta struct {
	Event   string
	Content string
}

func (s *Service) RespondToConversation(ctx context.Context, conversation storage.Conversation, user storage.User, userMessage string, writer EventWriter) (storage.Message, error) {
	// 记忆业务开关：系统级能力开启 且 用户个人偏好开启 时才注入/提取记忆。
	// 注意：它只控制记忆注入与提取，不控制会话锁与模型历史持久化。
	memoryOn := s.EnableMemory && user.MemoryEnabled
	// 入口获取会话锁：上一轮收尾未完成时阻塞等待，直到拿到锁或等待超时。
	// 获取失败（本地锁异常 / 等待超时）则降级放行，跳过本轮收尾。
	var lockToken string
	var stopRenew func()
	handedOff := false
	if s.EnableMemory {
		token := idgen.New("lock")
		ok, err := s.Store.AcquireConversationLock(ctx, conversation.ID, token, s.Cfg.ConversationLockTTL, s.Cfg.ConversationLockWaitTimeout)
		if err != nil {
			logger.Warn(fmt.Sprintf("conversation lock: acquire failed conversation=%s: %v", conversation.ID, err))
		} else if !ok {
			logger.Warn(fmt.Sprintf("conversation lock: acquire timed out conversation=%s", conversation.ID))
		} else {
			lockToken = token
			stopRenew = s.startLockRenewer(conversation.ID, token)
		}
	}
	// defer 兜底：任何提前 return（出错路径）都停止续期并释放锁；
	// 正常走到异步收尾时通过 handedOff 转交所有权，避免被误释放。
	defer func() {
		if handedOff || lockToken == "" {
			return
		}
		if stopRenew != nil {
			stopRenew()
		}
		_ = s.Store.ReleaseConversationLock(context.Background(), conversation.ID, lockToken)
	}()

	history, err := s.loadConversationMessages(ctx, conversation.ID)
	if err != nil {
		return storage.Message{}, err
	}
	modelHistory := s.loadModelHistory(ctx, conversation.ID, history)
	state := s.newLoopState(conversation, user, userMessage, history, modelHistory, writer)
	if err := s.hookManager().RunUserPromptSubmit(ctx, &UserPromptSubmitContext{State: state}); err != nil {
		return storage.Message{}, err
	}

	snapshot, err := s.buildSkillSnapshot(ctx, user.ID)
	if err != nil {
		return storage.Message{}, err
	}
	state.SkillSnapshot = snapshot
	state.SystemPrompt = s.buildSystemPrompt(ctx, conversation, user, snapshot, state.History, memoryOn)
	state.Messages = buildOpenAIMessages(state.SystemPrompt, state.ModelHistory)
	toolDefs := s.toolDefinitionsForUser(ctx, user.ID)
	round := 0
	roundsSinceTodoWrite := 0
	var cumulativeReasoning strings.Builder
	var lastRequestHistory []storage.Message

	for {
		round++
		if round > maxRound {
			return storage.Message{}, fmt.Errorf("conversation round limit %d exceeded", maxRound)
		}
		requestHistory, err := s.compressContextBeforeLLM(ctx, state)
		if err != nil {
			return storage.Message{}, err
		}
		lastRequestHistory = requestHistory
		state.Messages = buildOpenAIMessages(state.SystemPrompt, requestHistory)
		roundsSinceTodoWrite = maybeAppendTodoWriteReminder(state, s.Tools, roundsSinceTodoWrite)
		estimator := compression.DefaultTokenEstimator{}
		state.LastContextTokens = estimator.EstimateRequestTokens(state.SystemPrompt, requestHistory, toolDefs)
		state.LastContextBudget = estimator.ContextTokenBudget()
		emitMeta(state)
		req := openai.ChatCompletionRequest{
			Model:     s.Cfg.LLM.ModelID,
			Messages:  state.Messages,
			Tools:     toolDefs,
			MaxTokens: defaultMaxTokens,
		}
		reqBody, _ := json.Marshal(req)
		msg, finishReason, err := s.runModelRoundWithRecovery(ctx, state, req)
		respBody, _ := json.Marshal(msg)
		logger.LogLLMRound(round, fmt.Sprintf("main-agent conversation=%s", conversation.ID), reqBody, respBody, err)
		if err != nil {
			return storage.Message{}, err
		}
		cumulativeReasoning.WriteString(msg.ReasoningContent)
		requestMsg := msg
		state.Messages = append(state.Messages, requestMsg)
		if toolCallsInclude(msg.ToolCalls, agenttools.TodoWriteToolName) {
			roundsSinceTodoWrite = 0
		} else {
			roundsSinceTodoWrite++
		}

		if finishReason != "tool_calls" || len(msg.ToolCalls) == 0 {
			// 会话结束：把本轮模型最终回复纳入上下文后重算，得到会话结束时的最终 token 用量，
			// 覆盖请求前的估算值，确保存储与下发的都是最终用量。
			finalAssistant := storage.Message{Role: "assistant", Content: msg.Content, ReasoningContent: msg.ReasoningContent, ToolCalls: openAIToolCallsToStorage(msg.ToolCalls)}
			finalHistoryForEstimate := append(cloneMessages(lastRequestHistory), finalAssistant)
			finalEstimator := compression.DefaultTokenEstimator{}
			state.LastContextTokens = finalEstimator.EstimateRequestTokens(state.SystemPrompt, finalHistoryForEstimate, toolDefs)
			state.LastContextBudget = finalEstimator.ContextTokenBudget()
			stopCtx := &StopContext{State: state, ModelMessage: msg, Content: fallbackAssistantContent(msg.Content), ReasoningContent: cumulativeReasoning.String()}
			if err := s.hookManager().RunStop(ctx, stopCtx); err != nil {
				return storage.Message{}, err
			}
			if s.EnableMemory {
				finalHistory := append(state.History, stopCtx.AssistantMessage)
				finalModelHistory := append(cloneMessages(lastRequestHistory), stopCtx.AssistantMessage)
				handedOff = s.scheduleMemoryWork(conversation, user, finalHistory, finalModelHistory, lockToken, stopRenew, memoryOn)
			}
			return stopCtx.AssistantMessage, nil
		}
		assistantMessage := storage.Message{ID: state.NextMessageID(), ConversationID: conversation.ID, UserID: user.ID, Role: "assistant", Content: msg.Content, ReasoningContent: msg.ReasoningContent, ToolCalls: openAIToolCallsToStorage(msg.ToolCalls)}
		state.History = append(state.History, assistantMessage)
		state.ModelHistory = append(state.ModelHistory, assistantMessage)

		for _, tc := range msg.ToolCalls {
			toolCtx := &ToolUseContext{State: state, ToolCall: tc, Name: tc.Function.Name, RawArgs: tc.Function.Arguments}
			if err := s.hookManager().RunPreToolUse(ctx, toolCtx); err != nil {
				return storage.Message{}, err
			}
			toolCtx.Outcome = s.executeToolCall(ctx, ToolContext{User: user, Conversation: conversation, Skills: snapshot, ParentToolCallID: tc.ID, PersistedOutputReader: s.newPersistedOutputReader(conversation.ID, user.ID)}, tc.Function.Name, tc.Function.Arguments, toolCtx.Outcome.Audit)
			if toolCtx.Name == agenttools.TodoWriteToolName && toolCtx.Outcome.Status == "success" {
				state.Todos = append([]agenttools.TodoItem(nil), toolCtx.Outcome.Todos...)
			}
			if err := s.hookManager().RunPostToolUse(ctx, toolCtx); err != nil {
				return storage.Message{}, err
			}
			state.ToolCallCount++
			emitMeta(state)
		}
	}
}

// emitMeta 实时下发当前累计的回复元信息。
func emitMeta(state *LoopState) {
	if state == nil || state.Writer == nil {
		return
	}
	_ = state.Writer.Event("meta", map[string]any{
		"tool_call_count": state.ToolCallCount,
		"context_tokens":  state.LastContextTokens,
		"context_budget":  state.LastContextBudget,
	})
}

// runModelRoundWithRecovery wraps runModelRoundStream with truncation and
// context-overflow recovery. Transient 429/529 retries are handled inside the
// LLM client and are transparent here.
//
//   - Truncation (finish_reason == length): the first time, the output budget is
//     upgraded to truncationMaxTokens and the same request is retried WITHOUT
//     touching messages (the partial output is discarded). If it still
//     truncates, up to maxResumeAttempts continuation requests are issued, each
//     appending the truncated text plus a resume prompt; segments are streamed
//     and concatenated into one assistant message.
//   - Context overflow (HTTP 413): reactiveCompact is run once and the request
//     rebuilt from the compacted history, then retried. A second overflow is
//     returned to the caller (handled by the existing fallback boundary).
func (s *Service) runModelRoundWithRecovery(ctx context.Context, state *LoopState, req openai.ChatCompletionRequest) (openai.ChatCompletionMessage, openai.FinishReason, error) {
	cur := req
	upgraded := false
	compacted := false
	resumeAttempts := 0
	var accumulated strings.Builder
	var accumulatedReasoning strings.Builder

	for {
		msg, finishReason, deltas, err := s.runModelRoundStream(ctx, state, cur)
		if err != nil {
			if llm.IsContextOverflow(err) && !compacted {
				if compactErr := s.reactiveCompact(ctx, state); compactErr != nil {
					logger.Warn(fmt.Sprintf("reactive compact failed conversation=%s: %v", state.Conversation.ID, compactErr))
					return openai.ChatCompletionMessage{}, "", err
				}
				compacted = true
				cur.Messages = state.Messages
				cur.MaxTokens = defaultMaxTokens
				upgraded = false
				resumeAttempts = 0
				accumulated.Reset()
				accumulatedReasoning.Reset()
				continue
			}
			return openai.ChatCompletionMessage{}, "", err
		}

		resuming := accumulated.Len() > 0
		if finishReason != openai.FinishReasonLength {
			if resuming {
				flushContentDeltas(state, deltas)
				return openai.ChatCompletionMessage{
					Role:             "assistant",
					Content:          accumulated.String() + msg.Content,
					ReasoningContent: accumulatedReasoning.String() + msg.ReasoningContent,
					ToolCalls:        msg.ToolCalls,
				}, finishReason, nil
			}
			if state.Writer != nil && shouldEmitAssistantContentDeltas(finishReason, msg.ToolCalls) {
				flushContentDeltas(state, deltas)
			}
			return msg, finishReason, nil
		}

		// Truncated output.
		if !upgraded {
			// Discard the partial output, upgrade the budget, retry unchanged.
			upgraded = true
			cur.MaxTokens = truncationMaxTokens
			continue
		}
		if resumeAttempts >= maxResumeAttempts {
			// Best-effort: emit and return what we have so far.
			flushContentDeltas(state, deltas)
			logger.Warn(fmt.Sprintf("truncation resume exhausted conversation=%s after %d attempts", state.Conversation.ID, resumeAttempts))
			return openai.ChatCompletionMessage{
				Role:             "assistant",
				Content:          accumulated.String() + msg.Content,
				ReasoningContent: accumulatedReasoning.String() + msg.ReasoningContent,
			}, finishReason, nil
		}
		// Keep this partial segment and ask the model to continue.
		flushContentDeltas(state, deltas)
		accumulated.WriteString(msg.Content)
		accumulatedReasoning.WriteString(msg.ReasoningContent)
		resumeAttempts++
		cur.Messages = append(cur.Messages,
			openai.ChatCompletionMessage{Role: "assistant", Content: msg.Content},
			openai.ChatCompletionMessage{Role: "user", Content: truncationResumePrompt},
		)
	}
}

func flushContentDeltas(state *LoopState, deltas []bufferedModelDelta) {
	if state.Writer == nil {
		return
	}
	for _, delta := range deltas {
		_ = state.Writer.Event(delta.Event, map[string]any{"content": delta.Content})
	}
}

func (s *Service) runModelRoundStream(ctx context.Context, state *LoopState, req openai.ChatCompletionRequest) (openai.ChatCompletionMessage, openai.FinishReason, []bufferedModelDelta, error) {
	stream, err := s.LLM.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return openai.ChatCompletionMessage{}, "", nil, err
	}
	defer stream.Close()

	var content strings.Builder
	var reasoningContent strings.Builder
	var finishReason openai.FinishReason
	toolCalls := &streamedToolCallAccumulator{}
	var bufferedContentDeltas []bufferedModelDelta
	seenChoice := false
	seenOutput := false

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return openai.ChatCompletionMessage{}, "", nil, err
		}
		if len(chunk.Choices) == 0 {
			continue
		}
		seenChoice = true
		choice := chunk.Choices[0]
		if choice.Delta.Content != "" {
			seenOutput = true
			content.WriteString(choice.Delta.Content)
			bufferedContentDeltas = append(bufferedContentDeltas, bufferedModelDelta{Event: assistantDeltaEvent, Content: choice.Delta.Content})
		}
		if choice.Delta.ReasoningContent != "" {
			seenOutput = true
			reasoningContent.WriteString(choice.Delta.ReasoningContent)
			if state.Writer != nil {
				_ = state.Writer.Event(reasoningDeltaEvent, map[string]any{"content": choice.Delta.ReasoningContent})
			}
		}
		if len(choice.Delta.ToolCalls) > 0 {
			seenOutput = true
			toolCalls.Add(choice.Delta.ToolCalls)
		}
		if choice.FinishReason != "" {
			finishReason = choice.FinishReason
		}
	}
	if !seenChoice || (!seenOutput && finishReason == "") {
		return openai.ChatCompletionMessage{}, "", nil, fmt.Errorf("model stream returned no choices")
	}
	calls := toolCalls.Calls()

	return openai.ChatCompletionMessage{Role: "assistant", Content: content.String(), ReasoningContent: reasoningContent.String(), ToolCalls: calls}, finishReason, bufferedContentDeltas, nil
}

func shouldEmitAssistantContentDeltas(finishReason openai.FinishReason, toolCalls []openai.ToolCall) bool {
	return finishReason != openai.FinishReasonToolCalls && len(toolCalls) == 0
}

type streamedToolCallAccumulator struct {
	calls []openai.ToolCall
}

func (a *streamedToolCallAccumulator) Add(deltas []openai.ToolCall) {
	for deltaPosition, delta := range deltas {
		index := len(a.calls) - 1
		if delta.Index != nil {
			index = *delta.Index
		} else if len(deltas) > 1 {
			index = deltaPosition
		}
		if index < 0 {
			index = 0
		}
		for len(a.calls) <= index {
			a.calls = append(a.calls, openai.ToolCall{})
		}

		call := a.calls[index]
		if delta.ID != "" {
			call.ID = delta.ID
		}
		if delta.Type != "" {
			call.Type = delta.Type
		}
		if delta.Function.Name != "" {
			call.Function.Name = delta.Function.Name
		}
		if delta.Function.Arguments != "" {
			call.Function.Arguments += delta.Function.Arguments
		}
		a.calls[index] = call
	}
}

func (a *streamedToolCallAccumulator) Calls() []openai.ToolCall {
	if len(a.calls) == 0 {
		return nil
	}
	return append([]openai.ToolCall(nil), a.calls...)
}

func (s *Service) loadConversationMessages(ctx context.Context, conversationID string) ([]storage.Message, error) {
	if cached, ok, err := s.Store.GetConversationCache(ctx, conversationID); err == nil && ok {
		return cached, nil
	}
	messages, err := s.Store.ListMessagesByConversation(ctx, conversationID, 100)
	if err != nil {
		return nil, err
	}
	if err := s.Store.SetConversationCache(ctx, conversationID, messages); err != nil {
		_ = err
	}
	return messages, nil
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"context\"\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"io\"\n\t\"strings\"\n\n\topenai \"github.com/sashabaranov/go-openai\"\n\n\t\"nano_cc/internal/agent/runtime/compression\"\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/idgen\"\n\t\"nano_cc/internal/llm\"\n\t\"nano_cc/internal/logger\"\n\tagenttools \"nano_cc/internal/tools\"\n)\n\nconst (\n\tassistantDeltaEvent = \"assistant_delta\"\n\treasoningDeltaEvent = \"reasoning_delta\"\n\tmaxRound            = 50\n\n\t// defaultMaxTokens is the per-request output budget; on the "}
```

## 2026-06-13T09:16:46Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_W3earZp2crKYkhMKtwMq2417

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/helpers.go"}
```

### Result

```text
package runtime

import (
	"strings"

	"nano_cc/internal/idgen"
)

func newMessageID() string { return idgen.New("msg") }

func fallbackAssistantContent(content string) string {
	if strings.TrimSpace(content) == "" {
		return "(no response)"
	}
	return content
}

func shouldInferConversationTitle(currentTitle string) bool {
	trimmed := strings.TrimSpace(currentTitle)
	return trimmed == "" || trimmed == "新对话"
}

func inferConversationTitle(userMessage string) string {
	trimmed := strings.TrimSpace(userMessage)
	if len([]rune(trimmed)) > 30 {
		return string([]rune(trimmed)[:30])
	}
	if trimmed == "" {
		return "新对话"
	}
	return trimmed
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"strings\"\n\n\t\"nano_cc/internal/idgen\"\n)\n\nfunc newMessageID() string { return idgen.New(\"msg\") }\n\nfunc fallbackAssistantContent(content string) string {\n\tif strings.TrimSpace(content) == \"\" {\n\t\treturn \"(no response)\"\n\t}\n\treturn content\n}\n\nfunc shouldInferConversationTitle(currentTitle string) bool {\n\ttrimmed := strings.TrimSpace(currentTitle)\n\treturn trimmed == \"\" || trimmed == \"新对话\"\n}\n\nfunc inferConversationTitle(userMessage string) string {\n\ttrimmed := strings.Tr"}
```

## 2026-06-13T09:16:46Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_03_zs80MpySZXuhFbdW8qHI3187

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hook_bridge.go"}
```

### Result

```text
package runtime

import (
	runtimehooks "nano_cc/internal/agent/runtime/hooks"
	"nano_cc/internal/agent/storage"
	agenttools "nano_cc/internal/tools"
)

type UserPromptSubmitHook = runtimehooks.UserPromptSubmitHook
type PreToolUseHook = runtimehooks.PreToolUseHook
type PostToolUseHook = runtimehooks.PostToolUseHook
type StopHook = runtimehooks.StopHook
type HookManager = runtimehooks.HookManager
type LoopState = runtimehooks.LoopState
type UserPromptSubmitContext = runtimehooks.UserPromptSubmitContext
type ToolUseContext = runtimehooks.ToolUseContext
type StopContext = runtimehooks.StopContext

func NewDefaultHookManager() *HookManager {
	return runtimehooks.NewDefaultHookManager()
}

func (s *Service) newLoopState(conversation storage.Conversation, user storage.User, userInput string, history []storage.Message, modelHistory []storage.Message, writer EventWriter) *LoopState {
	return &LoopState{
		Store:                        s.Store,
		NewMessageID:                 newMessageID,
		ToolRuntimeEnv:               s.toolRuntimeEnv,
		ShouldInferConversationTitle: shouldInferConversationTitle,
		InferConversationTitle:       inferConversationTitle,
		Conversation:                 conversation,
		User:                         user,
		UserInput:                    userInput,
		History:                      history,
		ModelHistory:                 modelHistory,
		Writer:                       writer,
	}
}

func (s *Service) toolRuntimeEnv() agenttools.RuntimeEnv {
	if s == nil || s.Tools == nil {
		return agenttools.RuntimeEnv{}
	}
	return s.Tools.runtimeEnv()
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\truntimehooks \"nano_cc/internal/agent/runtime/hooks\"\n\t\"nano_cc/internal/agent/storage\"\n\tagenttools \"nano_cc/internal/tools\"\n)\n\ntype UserPromptSubmitHook = runtimehooks.UserPromptSubmitHook\ntype PreToolUseHook = runtimehooks.PreToolUseHook\ntype PostToolUseHook = runtimehooks.PostToolUseHook\ntype StopHook = runtimehooks.StopHook\ntype HookManager = runtimehooks.HookManager\ntype LoopState = runtimehooks.LoopState\ntype UserPromptSubmitContext = runtimehooks.UserPromptSubmitC"}
```

## 2026-06-13T09:16:46Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_04_HpyoAmOn0dVU3EjhyuHk1114

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/prompt_builder.go"}
```

### Result

```text
package runtime

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/assistant"
	"nano_cc/internal/sessions"
	agenttools "nano_cc/internal/tools"
)

func (s *Service) buildSystemPrompt(ctx context.Context, conversation storage.Conversation, user storage.User, snapshot *agenttools.SkillSnapshot, history []storage.Message, memoryOn bool) string {
	memorySection := ""
	if memoryOn {
		memorySection = s.selectRelevantMemories(ctx, user, history)
	}
	return s.buildSystemPromptWithMemory(user, snapshot, memorySection)
}

func (s *Service) buildSystemPromptWithMemory(user storage.User, snapshot *agenttools.SkillSnapshot, memorySection string) string {
	toolNames := []string(nil)
	if s.Tools != nil {
		toolDefs := s.Tools.Definitions()
		toolNames = make([]string, 0, len(toolDefs))
		for _, def := range toolDefs {
			if def.Function == nil || strings.TrimSpace(def.Function.Name) == "" {
				continue
			}
			toolNames = append(toolNames, def.Function.Name)
		}
	}
	loader := sessions.NewSkillLoader()
	if snapshot != nil && snapshot.Merged != nil {
		loader = snapshot.Merged
	}
	return assistant.BuildSystemPrompt(assistant.PromptOptions{
		BasePrompt:        s.BasePrompt,
		Surface:           fmt.Sprintf("the local TUI session for user %s", user.Username),
		SkillDescriptions: loader.GetDescriptions(),
		MemorySection:     memorySection,
		WorkingDirectory:  strings.TrimSpace(s.Cfg.WorkspaceRoot),
		LinkMarkdown: assistant.LinkMarkdownContext{
			UserPath:         s.LinkMarkdown.UserPath,
			UserContent:      s.LinkMarkdown.UserContent,
			WorkspacePath:    s.LinkMarkdown.WorkspacePath,
			WorkspaceContent: s.LinkMarkdown.WorkspaceContent,
		},
		ToolNames: toolNames,
	})
}

func (s *Service) buildSkillSnapshot(ctx context.Context, userID string) (*agenttools.SkillSnapshot, error) {
	return agenttools.NewSkillSnapshot(nil, s.BuiltinSkills), nil
}

func buildOpenAIMessages(systemPrompt string, history []storage.Message) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{{Role: "system", Content: systemPrompt}}
	for _, msg := range history {
		messages = append(messages, openai.ChatCompletionMessage{Role: msg.Role, Content: msg.Content, ReasoningContent: msg.ReasoningContent, ToolCallID: msg.ToolCallID, ToolCalls: storageToolCallsToOpenAI(msg.ToolCalls)})
	}
	return messages
}

func storageToolCallsToOpenAI(toolCalls []storage.MessageToolCall) []openai.ToolCall {
	if len(toolCalls) == 0 {
		return nil
	}
	converted := make([]openai.ToolCall, 0, len(toolCalls))
	for _, toolCall := range toolCalls {
		converted = append(converted, openai.ToolCall{ID: toolCall.ID, Type: openai.ToolType(toolCall.Type), Function: openai.FunctionCall{Name: toolCall.Function.Name, Arguments: toolCall.Function.Arguments}})
	}
	return converted
}

func openAIToolCallsToStorage(toolCalls []openai.ToolCall) []storage.MessageToolCall {
	if len(toolCalls) == 0 {
		return nil
	}
	converted := make([]storage.MessageToolCall, 0, len(toolCalls))
	for _, toolCall := range toolCalls {
		converted = append(converted, storage.MessageToolCall{ID: toolCall.ID, Type: string(toolCall.Type), Function: storage.MessageFunctionCall{Name: toolCall.Function.Name, Arguments: toolCall.Function.Arguments}})
	}
	return converted
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"context\"\n\t\"fmt\"\n\t\"strings\"\n\n\topenai \"github.com/sashabaranov/go-openai\"\n\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/assistant\"\n\t\"nano_cc/internal/sessions\"\n\tagenttools \"nano_cc/internal/tools\"\n)\n\nfunc (s *Service) buildSystemPrompt(ctx context.Context, conversation storage.Conversation, user storage.User, snapshot *agenttools.SkillSnapshot, history []storage.Message, memoryOn bool) string {\n\tmemorySection := \"\"\n\tif memoryOn {\n\t\tmemorySection = s.selectReleva"}
```

## 2026-06-13T09:16:48Z · bash · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_i96W3YBpmEcbSwVl2L4m5983

### Arguments

```json
{"command": "find /Users/bytedance/golang_pro/nano_cc/go-agent -type f -name \"*.go\" -not -path \"*/vendor/*\" -not -name \"*_test.go\" | sort"}
```

### Result

```text
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/transport.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/conversation_memory_strategy.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/full_history_summarization.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/message_window_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/reactive_compact.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/recent_tool_result_retention.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/token_estimator.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/tool_result_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/context_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/conversation_flow.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/conversation_memory.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/helpers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hook_bridge.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/manager.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/stop.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/tool.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/types.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/user_prompt.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/memory.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/prompt_builder.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/runtime.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/subagent.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/summarizer.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/todo_reminder.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/tool_registry.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/conversation_history.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/models.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/assistant/prompt.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/cli/root.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/helpers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/local_config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/paths.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/idgen/idgen.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/llm/client.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/llm/errors.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/bootstrap.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/memory_store.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/persisted_output_files.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/session_history_store.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/store.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/tool_result_log.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/logger/logger.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/safety/path.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/sessions/skill.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/textutil/textutil.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/bash.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/definitions.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/file_ops.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/handlers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/load_skill.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/path_guard.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/persisted_output.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/runtime_env.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/todo_write.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/validation.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/app.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/events.go
/Users/bytedance/golang_pro/nano_cc/go-agent/main.go
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/basic.go
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","resolved_command_path":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/transport.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/compression.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/conversation_memory_strategy.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime"}
```

## 2026-06-13T09:16:52Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_ct5cKvLSkO8Mp33Hp8kr4187

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/definitions.go"}
```

### Result

```text
package tools

import (
	"encoding/json"

	openai "github.com/sashabaranov/go-openai"
)

func mustMarshal(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func toolDef(name, desc string, params any) openai.Tool {
	return openai.Tool{
		Type: "function",
		Function: &openai.FunctionDefinition{
			Name:        name,
			Description: desc,
			Parameters:  mustMarshal(params),
		},
	}
}

func strParam(desc string) map[string]any {
	return map[string]any{
		"type":        "string",
		"description": desc,
	}
}

func intParam(desc string) map[string]any {
	return map[string]any{
		"type":        "integer",
		"description": desc,
	}
}

var baseToolDefs = []openai.Tool{
	toolDef("bash", "Execute a shell command via bash -c. Relative path arguments are interpreted under the workspace root; absolute paths outside the workspace and dangerous commands are rejected unless explicitly allowed by configuration.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": strParam("The shell command to execute"),
		},
		"required": []string{"command"},
	}),
	toolDef("read_file", "Read a file from the filesystem", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":  strParam("Path to the file to read"),
			"limit": intParam("Maximum number of lines to read"),
		},
		"required": []string{"path"},
	}),
	toolDef("write_file", "Write content to a file", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":    strParam("Path to the file to write"),
			"content": strParam("Content to write to the file"),
		},
		"required": []string{"path", "content"},
	}),
	toolDef("edit_file", "Replace text in a file by exact match", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":     strParam("Path to the file to edit"),
			"old_text": strParam("Exact text to find and replace"),
			"new_text": strParam("Text to replace it with"),
		},
		"required": []string{"path", "old_text", "new_text"},
	}),
	toolDef("load_skill", "Load the full instructions of a local skill by exact name before using or following that skill. Skills are loaded from the user's ~/.link/skills and the workspace .link/skills directories, with workspace skills taking precedence.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": strParam("Name of the skill to load"),
		},
		"required": []string{"name"},
	}),
	toolDef("todo_write", "Create or update the current task plan. Use this tool to track progress on multi-step tasks.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"todos": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"id":      map[string]any{"type": "string"},
						"content": map[string]any{"type": "string"},
						"status":  map[string]any{"type": "string", "enum": []string{TodoStatusPending, TodoStatusInProgress, TodoStatusCompleted}},
					},
					"required": []string{"id", "content", "status"},
				},
			},
		},
		"required": []string{"todos"},
	}),
}

var spawnSubagentToolDef = toolDef("spawn_subagent", "Spawn a child agent with a fresh message list to complete an isolated task. The child agent may use workspace tools, but it cannot spawn another subagent. Only its final summary is returned to the parent agent.", map[string]any{
	"type": "object",
	"properties": map[string]any{
		"task": strParam("The task for the child agent to complete. Include all context it needs because parent conversation history is not shared."),
		"cwd":  strParam("Optional working directory for the child agent. Relative paths are resolved under the workspace root; absolute paths must remain inside the workspace."),
	},
	"required": []string{"task"},
})

// ReadPersistedOutputToolName is exposed automatically alongside context
// compression so the model can fetch the full content behind a
// <persisted-output> marker when the inline preview is insufficient.
const ReadPersistedOutputToolName = "read_persisted_output"

var ReadPersistedOutputToolDef = toolDef(ReadPersistedOutputToolName, "Read a chunk of a persisted tool output by id when a <persisted-output> marker preview is insufficient. Only outputs from the current conversation are accessible.", map[string]any{
	"type": "object",
	"properties": map[string]any{
		"id":     strParam("The persisted output id from the <persisted-output> marker, for example po_abc123."),
		"offset": intParam("Zero-based character offset to start reading from. Defaults to 0."),
		"limit":  intParam("Maximum characters to return. Defaults to 20000 and is capped by the runtime."),
	},
	"required": []string{"id"},
})

var AllToolDefs = append(append([]openai.Tool(nil), baseToolDefs...), spawnSubagentToolDef, ReadPersistedOutputToolDef)
var ChildToolDefs = baseToolDefs

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package tools\n\nimport (\n\t\"encoding/json\"\n\n\topenai \"github.com/sashabaranov/go-openai\"\n)\n\nfunc mustMarshal(v any) json.RawMessage {\n\tb, err := json.Marshal(v)\n\tif err != nil {\n\t\tpanic(err)\n\t}\n\treturn b\n}\n\nfunc toolDef(name, desc string, params any) openai.Tool {\n\treturn openai.Tool{\n\t\tType: \"function\",\n\t\tFunction: \u0026openai.FunctionDefinition{\n\t\t\tName:        name,\n\t\t\tDescription: desc,\n\t\t\tParameters:  mustMarshal(params),\n\t\t},\n\t}\n}\n\nfunc strParam(desc string) map[string]any {\n\treturn map[string]an"}
```

## 2026-06-13T09:16:52Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_mYbcQtSDW4fslraGC4DL4640

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/handlers.go"}
```

### Result

```text
package tools

import (
	"context"
	"fmt"
)

// TodoWriteToolName is the tool that updates the task plan and returns
// structured todo items in addition to a textual summary.
const TodoWriteToolName = "todo_write"

// ExecResult is the unified result of executing a stateless tool. Todos is
// populated only by the todo_write tool.
type ExecResult struct {
	Output string
	Todos  []TodoItem
}

// Handlers maps stateless tool names to their textual handlers. todo_write is
// dispatched separately because it returns structured todos (see Dispatch).
var Handlers = map[string]ToolHandler{
	"bash":                  handleBash,
	"read_file":             handleRead,
	"write_file":            handleWrite,
	"edit_file":             handleEdit,
	"load_skill":            handleLoadSkill,
	"read_persisted_output": handleReadPersistedOutput,
}

// Dispatch is the single entry point for executing a stateless tool by name.
// It is the authority for tool execution semantics, including todo_write's
// structured output.
func Dispatch(ctx context.Context, name string, args map[string]any) (ExecResult, error) {
	if name == TodoWriteToolName {
		result, err := ExecuteTodoWrite(ctx, args)
		if err != nil {
			return ExecResult{}, err
		}
		return ExecResult{Output: result.Output, Todos: result.Todos}, nil
	}
	handler, ok := Handlers[name]
	if !ok || handler == nil {
		return ExecResult{}, fmt.Errorf("tool %s has no handler", name)
	}
	output, err := handler(ctx, args)
	if err != nil {
		return ExecResult{}, err
	}
	return ExecResult{Output: output}, nil
}

func handleBash(ctx context.Context, args map[string]any) (string, error) {
	cmd, _ := args["command"].(string)
	if cmd == "" {
		return "", fmt.Errorf("command is required")
	}
	root, err := validatedWorkspaceRootFromContext(ctx)
	if err != nil {
		return "", err
	}
	workingDir, err := validatedCurrentWorkingDirFromContext(ctx)
	if err != nil {
		return "", err
	}
	if err := validateBashCommandPaths(root, cmd, allowOutsideWorkspaceFromContext(ctx), systemAssetDirsFromContext(ctx)...); err != nil {
		return "", err
	}
	return RunBashInDirWithOptions(cmd, workingDir, allowDangerousCommandsFromContext(ctx))
}

func handleRead(ctx context.Context, args map[string]any) (string, error) {
	path, _ := args["path"].(string)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	limit := 0
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}
	root, resolvedPath, err := resolvePathFromContext(ctx, path)
	if err != nil {
		return "", err
	}
	return RunReadFromRoot(root, resolvedPath, limit)
}

func handleWrite(ctx context.Context, args map[string]any) (string, error) {
	path, _ := args["path"].(string)
	content, _ := args["content"].(string)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	root, resolvedPath, err := resolvePathFromContext(ctx, path)
	if err != nil {
		return "", err
	}
	return RunWriteFromRoot(root, resolvedPath, content)
}

func handleEdit(ctx context.Context, args map[string]any) (string, error) {
	path, _ := args["path"].(string)
	oldText, _ := args["old_text"].(string)
	newText, _ := args["new_text"].(string)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	root, resolvedPath, err := resolvePathFromContext(ctx, path)
	if err != nil {
		return "", err
	}
	return RunEditFromRoot(root, resolvedPath, oldText, newText)
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package tools\n\nimport (\n\t\"context\"\n\t\"fmt\"\n)\n\n// TodoWriteToolName is the tool that updates the task plan and returns\n// structured todo items in addition to a textual summary.\nconst TodoWriteToolName = \"todo_write\"\n\n// ExecResult is the unified result of executing a stateless tool. Todos is\n// populated only by the todo_write tool.\ntype ExecResult struct {\n\tOutput string\n\tTodos  []TodoItem\n}\n\n// Handlers maps stateless tool names to their textual handlers. todo_write is\n// dispatched separately "}
```

## 2026-06-13T09:16:52Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_Lr8qMIBlI6vSV2s0Dboj5215

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/assistant/prompt.go"}
```

### Result

```text
package assistant

import (
	"os"
	"strings"
)

const persistedOutputGuidance = "当较早的消息中出现 `<persisted-output ...>` 标记时，表示完整的工具输出已存入本地文件，内联的只是预览。" +
	"如果预览不足以完成任务，请用标记中的 id 和偏移量调用 `read_persisted_output` 分块读取更多内容，不要猜测被省略的部分。" +
	"当看到 `[Earlier result compacted. Re-run if needed]` 时，请重新执行相关工具以再次获取该结果。"

const DefaultBaseSystemPrompt = "你是 Link，一个运行在本地 TUI 中的通用型智能体（general-purpose agent），而不是只能聊天的助手。\n\n" +
	"帮助用户处理日常问答、分析、规划、写作，以及阅读代码、修改代码、生成文件等任务。优先给出直接、有用的回答；当工具或技能能帮助你完成任务时，主动使用它们。\n\n" +
	"你为用户修改代码或生成文件时，可能会在运行时的工作区里实际操作，但工作区只是对用户不可见的临时执行环境。完成后必须把修改后的完整代码、生成的文件内容原封不动地返回给用户，不要省略、截断或改写；返回之后工作区里的对应文件即可删除。\n\n" +
	"绝不向用户透传任何终端 / 命令行相关信息：包括但不限于\"执行 Shell 命令\"\"运行 bash\"\"终端操作\"\"命令行\"等表述，以及具体命令、脚本、工作目录、文件路径、运行环境、工具名称、技能调用方式等内部执行机制。描述你的工作时用\"我帮你阅读/修改了代码\"\"我为你生成了文件\"这类面向结果的说法。\n\n" +
	"当用户询问关于你自身的问题时，只介绍你能为用户提供的能力与价值（如阅读代码、修改代码、生成文件等）；即使运行时提供了上述信息，也仅供内部决策使用，绝不复述给用户，除非用户为完成具体任务确实需要。"

type PromptOptions struct {
	BasePrompt        string
	Surface           string
	SkillDescriptions string
	MemorySection     string
	WorkingDirectory  string
	LinkMarkdown      LinkMarkdownContext
	ToolNames         []string
}

type LinkMarkdownContext struct {
	UserPath         string
	UserContent      string
	WorkspacePath    string
	WorkspaceContent string
}

func LoadBaseSystemPrompt(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

func BuildSystemPrompt(opts PromptOptions) string {
	surface := strings.TrimSpace(opts.Surface)
	if surface == "" {
		surface = "this conversation"
	}

	basePrompt := strings.TrimSpace(opts.BasePrompt)
	if basePrompt == "" {
		basePrompt = DefaultBaseSystemPrompt
	}

	sections := []string{renderTag("identity", basePrompt)}

	workspaceLines := []string{"Surface: " + surface}
	if workingDirectory := strings.TrimSpace(opts.WorkingDirectory); workingDirectory != "" {
		workspaceLines = append(workspaceLines,
			"Working directory: "+workingDirectory,
			"除非运行时另有说明，默认以工作目录作为运行时文件与 Shell 操作的根目录。",
		)
	}
	sections = append(sections, renderTag("workspace", strings.Join(workspaceLines, "\n")))
	if linkContext := renderLinkMarkdownContext(opts.LinkMarkdown); linkContext != "" {
		sections = append(sections, renderTag("system-reminder", linkContext))
	}

	toolNames := make([]string, 0, len(opts.ToolNames))
	for _, name := range opts.ToolNames {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			continue
		}
		toolNames = append(toolNames, trimmed)
	}
	if len(toolNames) > 0 {
		toolBody := "本次会话可用的工具如下：\n\n" + renderList(toolNames)
		if toolNamesContain(toolNames, "read_persisted_output") {
			toolBody += "\n\n" + persistedOutputGuidance
		}
		sections = append(sections, renderTag("tools", toolBody))
	}

	if descriptions := strings.TrimSpace(opts.SkillDescriptions); descriptions != "" {
		skillBody := strings.Join([]string{
			"以下技能只提供摘要。",
			"重要规则：\n" + renderList([]string{
				"每个技能都有名称和描述。",
				"使用或遵循某个技能前，先用 `load_skill` 以精确的技能名加载其完整说明。",
				"不要仅凭摘要臆测完整的工作流。",
				"若多个技能看起来都相关，先加载最匹配、最具体的那个。",
			}),
			"可用技能：\n\n" + descriptions,
		}, "\n\n")
		sections = append(sections, renderTag("skills", skillBody))
	}

	if memory := strings.TrimSpace(opts.MemorySection); memory != "" {
		sections = append(sections, renderTag("memory", memory))
	}

	return strings.Join(sections, "\n\n")
}

func renderLinkMarkdownContext(ctx LinkMarkdownContext) string {
	userContent := strings.TrimSpace(ctx.UserContent)
	workspaceContent := strings.TrimSpace(ctx.WorkspaceContent)
	if userContent == "" && workspaceContent == "" {
		return ""
	}
	parts := []string{
		"在回答用户问题时，你可以参考以下上下文：",
		"# linkMd",
		"下面展示了用户与代码库说明。请务必遵循这些说明。重要：这些说明将覆盖任何默认行为，你必须严格按其文字要求执行。",
	}
	if userContent != "" {
		parts = append(parts, strings.Join([]string{
			strings.TrimSpace(ctx.UserPath) + " 的内容（用户为所有项目配置的私人全局说明）：",
			userContent,
		}, "\n\n"))
	}
	if workspaceContent != "" {
		parts = append(parts, strings.Join([]string{
			strings.TrimSpace(ctx.WorkspacePath) + " 的内容（项目说明，已提交到代码库或工作区）：",
			workspaceContent,
		}, "\n\n"))
	}
	parts = append(parts, strings.Join([]string{
		"# 重要指令提醒",
		"只做被要求的事：不多不少。",
		"除非为达成目标绝对必要，切勿创建新文件。",
		"能修改现有文件，绝不新建文件。",
		"不要主动创建文档文件（*.md）或README。仅当用户明确要求时才创建文档。",
		"重要：这些上下文可能与当前任务相关，也可能无关。除非与任务高度相关，否则不要对其作出回应。",
	}, "\n"))
	return strings.Join(parts, "\n\n")
}

func renderTag(tag, body string) string {
	tag = strings.TrimSpace(tag)
	body = strings.TrimSpace(body)
	if body == "" {
		return "<" + tag + ">\n</" + tag + ">"
	}
	return "<" + tag + ">\n" + body + "\n</" + tag + ">"
}

func toolNamesContain(names []string, target string) bool {
	for _, name := range names {
		if name == target {
			return true
		}
	}
	return false
}

func renderList(items []string) string {
	lines := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		lines = append(lines, "- "+trimmed)
	}
	return strings.Join(lines, "\n")
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package assistant\n\nimport (\n\t\"os\"\n\t\"strings\"\n)\n\nconst persistedOutputGuidance = \"当较早的消息中出现 `\u003cpersisted-output ...\u003e` 标记时，表示完整的工具输出已存入本地文件，内联的只是预览。\" +\n\t\"如果预览不足以完成任务，请用标记中的 id 和偏移量调用 `read_persisted_output` 分块读取更多内容，不要猜测被省略的部分。\" +\n\t\"当看到 `[Earlier result compacted. Re-run if needed]` 时，请重新执行相关工具以再次获取\ufffd"}
```

## 2026-06-13T09:16:52Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_03_IBvX2UFT6bMsMwziK8Ye8631

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/compression.go"}
```

### Result

```text
package compression

import (
	"context"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
)

const (
	// toolResultByteThreshold is the inline tool_result budget for the latest
	// user turn before large outputs get persisted (200KB).
	toolResultByteThreshold = 200 * 1024
	// toolResultPreviewRunes is how many leading runes of an oversized
	// tool_result remain inline as a preview.
	toolResultPreviewRunes = 2000
	// messageWindowLimit triggers head/tail trimming when exceeded.
	messageWindowLimit = 50
	messageWindowHead  = 3
	messageWindowTail  = 47
	// recentToolResultRetention keeps the most recent N full tool results.
	recentToolResultRetention = 3

	earlierToolResultPlaceholder = "[Earlier result compacted. Re-run if needed]"
	// PersistedOutputMarkerPrefix marks a tool_result whose full content was
	// persisted out-of-band and replaced inline with a preview.
	PersistedOutputMarkerPrefix = "<persisted-output"
)

// Request carries the request-only history copy and dependencies that
// strategies use. Strategies mutate RequestHistory in place.
type Request struct {
	Conversation   storage.Conversation
	User           storage.User
	RequestHistory []storage.Message
	SystemPrompt   string
	Tools          []openai.Tool
	Store          Store
	Estimator      TokenEstimator
	Summarizer     HistorySummarizer
}

// Store is the minimal storage surface required by strategies.
type Store interface {
	CreatePersistedOutput(ctx context.Context, output storage.PersistedOutput) error
	GetPersistedOutputByMessageHash(ctx context.Context, conversationID, userID, messageID, toolCallID, strategy, contentSHA256 string) (storage.PersistedOutput, error)
	CreateContextSummary(ctx context.Context, summary storage.ContextSummary) error
	GetContextSummaryByHistoryHash(ctx context.Context, conversationID, userID, sourceHistorySHA256 string) (storage.ContextSummary, error)
	ListConversationMemories(ctx context.Context, conversationID string) ([]storage.ConversationMemory, error)
}

// Strategy applies one compression layer to RequestHistory.
type Strategy interface {
	Name() string
	Apply(ctx context.Context, req *Request) error
}

// Compressor runs the registered strategies in order.
type Compressor struct {
	strategies []Strategy
}

func NewDefaultCompressor() *Compressor {
	return &Compressor{strategies: []Strategy{
		&ToolResultCompressionStrategy{},
		&MessageWindowCompressionStrategy{},
		&RecentToolResultRetentionStrategy{},
		&ConversationMemoryStrategy{},
		&FullHistorySummarizationStrategy{},
	}}
}

func (c *Compressor) Compress(ctx context.Context, req *Request) error {
	for _, strategy := range c.strategies {
		if err := strategy.Apply(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package compression\n\nimport (\n\t\"context\"\n\n\topenai \"github.com/sashabaranov/go-openai\"\n\n\t\"nano_cc/internal/agent/storage\"\n)\n\nconst (\n\t// toolResultByteThreshold is the inline tool_result budget for the latest\n\t// user turn before large outputs get persisted (200KB).\n\ttoolResultByteThreshold = 200 * 1024\n\t// toolResultPreviewRunes is how many leading runes of an oversized\n\t// tool_result remain inline as a preview.\n\ttoolResultPreviewRunes = 2000\n\t// messageWindowLimit triggers head/tail trimming w"}
```

## 2026-06-13T09:16:52Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_04_nOPROUiblk4SdGckHBgr9952

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/runtime.go"}
```

### Result

```text
package runtime

import (
	"context"
	"time"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/runtime/compression"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/config"
	"nano_cc/internal/llm"
	"nano_cc/internal/sessions"
)

type EventWriter interface {
	Event(name string, data any) error
}

type conversationStore interface {
	UpdateConversationTitle(ctx context.Context, conversationID, title string) error
	TouchConversationActivity(ctx context.Context, conversationID string) error
	SetConversationHistory(ctx context.Context, conversationID string, messages []storage.Message) error
	SetConversationCache(ctx context.Context, conversationID string, messages []storage.Message) error
	GetConversationCache(ctx context.Context, conversationID string) ([]storage.Message, bool, error)
	ListMessagesByConversation(ctx context.Context, conversationID string, limit int) ([]storage.Message, error)
	CreateToolCall(ctx context.Context, tc storage.ToolCall) error
	CreateSubagentMessage(ctx context.Context, message storage.SubagentMessage) error
	CreatePersistedOutput(ctx context.Context, output storage.PersistedOutput) error
	GetPersistedOutputForConversation(ctx context.Context, id, userID, conversationID string) (storage.PersistedOutput, error)
	GetPersistedOutputByMessageHash(ctx context.Context, conversationID, userID, messageID, toolCallID, strategy, contentSHA256 string) (storage.PersistedOutput, error)
	CreateContextSummary(ctx context.Context, summary storage.ContextSummary) error
	GetContextSummaryByHistoryHash(ctx context.Context, conversationID, userID, sourceHistorySHA256 string) (storage.ContextSummary, error)
	ListRelevantMemories(ctx context.Context, userID string) ([]storage.Memory, error)
	ListMemoriesByUserAndType(ctx context.Context, userID, memType string) ([]storage.Memory, error)
	ListProjectFactMemories(ctx context.Context, userID string) ([]storage.Memory, error)
	InsertMemory(ctx context.Context, m storage.Memory) error
	CountMemoriesByUserAndType(ctx context.Context, userID, memType string) (int, error)
	CountProjectFactMemories(ctx context.Context, userID string) (int, error)
	DeleteOldestMemories(ctx context.Context, userID, memType string, n int) error
	ReplaceMemoriesByUserAndType(ctx context.Context, userID, memType string, items []storage.Memory) error
	ReplaceProjectFactMemories(ctx context.Context, userID string, items []storage.Memory) error
	ListConversationMemories(ctx context.Context, conversationID string) ([]storage.ConversationMemory, error)
	ReplaceConversationMemories(ctx context.Context, conversationID, userID string, items []storage.ConversationMemory) error
	GetConversationModelHistory(ctx context.Context, conversationID string) ([]storage.Message, bool, error)
	UpsertConversationModelHistory(ctx context.Context, conversationID, userID string, messages []storage.Message) error
	AcquireConversationLock(ctx context.Context, conversationID, token string, ttl, waitTimeout time.Duration) (bool, error)
	RenewConversationLock(ctx context.Context, conversationID, token string, ttl time.Duration) (bool, error)
	ReleaseConversationLock(ctx context.Context, conversationID, token string) error
}

type Service struct {
	Store             conversationStore
	Cfg               config.AppConfig
	LLM               llm.Client
	Tools             *ToolRegistry
	BuiltinSkills     *sessions.SkillLoader
	BasePrompt        string
	LinkMarkdown      config.LinkMarkdownContext
	Hooks             *HookManager
	ContextCompressor *compression.Compressor
	EnableMemory      bool
	MCP               *mcp.Manager
}

func NewService(store conversationStore, cfg config.AppConfig, client llm.Client) *Service {
	return &Service{Store: store, Cfg: cfg, LLM: client, Tools: NewToolRegistry(cfg), Hooks: NewDefaultHookManager(), EnableMemory: true}
}

func (s *Service) hookManager() *HookManager {
	if s.Hooks == nil {
		s.Hooks = NewDefaultHookManager()
	}
	return s.Hooks
}

func (s *Service) SetBuiltinSkills(loader *sessions.SkillLoader) {
	s.BuiltinSkills = loader
}

func (s *Service) SetBasePrompt(prompt string) {
	s.BasePrompt = prompt
}

func (s *Service) SetLinkMarkdownContext(ctx config.LinkMarkdownContext) {
	s.LinkMarkdown = ctx
}

func (s *Service) SetMCPManager(manager *mcp.Manager) {
	s.MCP = manager
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package runtime\n\nimport (\n\t\"context\"\n\t\"time\"\n\n\t\"nano_cc/internal/agent/mcp\"\n\t\"nano_cc/internal/agent/runtime/compression\"\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/config\"\n\t\"nano_cc/internal/llm\"\n\t\"nano_cc/internal/sessions\"\n)\n\ntype EventWriter interface {\n\tEvent(name string, data any) error\n}\n\ntype conversationStore interface {\n\tUpdateConversationTitle(ctx context.Context, conversationID, title string) error\n\tTouchConversationActivity(ctx context.Context, conversationID string) error"}
```

## 2026-06-13T09:16:52Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_05_JQdPZ1cQxZSORR7kgEQA2617

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/bootstrap.go"}
```

### Result

```text
package local

import (
	"context"
	"fmt"
	"time"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/runtime"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/assistant"
	"nano_cc/internal/config"
	"nano_cc/internal/idgen"
	"nano_cc/internal/llm"
	"nano_cc/internal/logger"
	"nano_cc/internal/sessions"
)

const LocalUserID = "local-user"

type Bundle struct {
	Runtime      *runtime.Service
	Store        *Store
	MCP          *mcp.Manager
	User         storage.User
	Conversation storage.Conversation
	CWD          string
	SkillCount   int
	Skills       []sessions.SkillSummary
	MCPToolCount int
	MCPServers   []mcp.ServerStatus
}

func Bootstrap(ctx context.Context, cwd string) (*Bundle, error) {
	cfg, err := config.LoadLocalConfig(cwd)
	if err != nil {
		return nil, err
	}
	if err := config.EnsureAppLayout(cfg); err != nil {
		return nil, err
	}
	if err := config.ValidateAppLayout(cfg); err != nil {
		return nil, err
	}
	if err := logger.InitFileLoggerAt(cfg.LogsDir); err != nil {
		logger.Warn(fmt.Sprintf("failed to init file logger: %v", err))
	}
	userSkillsDir, err := config.LinkSkillsDir()
	if err != nil {
		return nil, err
	}
	builtinSkills, err := sessions.LoadSkillsFromDirs([]sessions.SkillDir{
		{Path: userSkillsDir, Source: "user"},
		{Path: config.WorkspaceLinkSkillsDir(cfg.WorkspaceRoot), Source: "workspace"},
	})
	if err != nil {
		return nil, fmt.Errorf("load .link skills: %w", err)
	}
	basePrompt, err := assistant.LoadBaseSystemPrompt(cfg.SystemPromptPath)
	if err != nil {
		return nil, fmt.Errorf("load system prompt: %w", err)
	}
	linkMarkdown, err := config.LoadLinkMarkdownContext(cfg.WorkspaceRoot)
	if err != nil {
		return nil, fmt.Errorf("load LINK.MD context: %w", err)
	}
	store, err := NewStoreWithMemory(cfg.WorkspaceRoot)
	if err != nil {
		return nil, fmt.Errorf("init memory store: %w", err)
	}
	client := llm.NewDeepseekClient(cfg.LLM.BaseURL, cfg.LLM.APIKey)
	runtimeService := runtime.NewService(store, cfg, client)
	runtimeService.EnableMemory = true
	runtimeService.SetBuiltinSkills(builtinSkills)
	runtimeService.SetBasePrompt(basePrompt)
	runtimeService.SetLinkMarkdownContext(linkMarkdown)
	mcpManager := mcp.NewManager()
	workspaceMCPServers, err := mcp.LoadWorkspaceConfig(config.WorkspaceMCPConfigPath(cfg.WorkspaceRoot))
	if err != nil {
		mcpManager.Close()
		return nil, fmt.Errorf("load workspace mcp config: %w", err)
	}
	mcpManager.SetWorkspaceServers(ctx, workspaceMCPServers)
	runtimeService.SetMCPManager(mcpManager)
	user := storage.User{ID: LocalUserID, Email: "local@go-agent", Username: "local", MemoryEnabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	conversation := storage.Conversation{ID: idgen.New("conv"), SessionID: idgen.UUID(), UserID: user.ID, RootMessageID: idgen.New("msg"), Title: "TUI 会话", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := store.CreateConversation(ctx, conversation); err != nil {
		mcpManager.Close()
		return nil, err
	}
	mcpManager.EnsureWorkspaceSessions(ctx)
	mcpTools := mcpManager.ToolsForUser(user.ID)
	mcpSnapshot := mcpManager.Snapshot(user.ID)
	skills := builtinSkills.Summaries()
	return &Bundle{Runtime: runtimeService, Store: store, MCP: mcpManager, User: user, Conversation: conversation, CWD: cfg.WorkspaceRoot, SkillCount: len(skills), Skills: skills, MCPToolCount: len(mcpTools), MCPServers: mcpSnapshot.Servers}, nil
}

func (b *Bundle) Close() {
	if b != nil && b.MCP != nil {
		b.MCP.Close()
	}
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package local\n\nimport (\n\t\"context\"\n\t\"fmt\"\n\t\"time\"\n\n\t\"nano_cc/internal/agent/mcp\"\n\t\"nano_cc/internal/agent/runtime\"\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/assistant\"\n\t\"nano_cc/internal/config\"\n\t\"nano_cc/internal/idgen\"\n\t\"nano_cc/internal/llm\"\n\t\"nano_cc/internal/logger\"\n\t\"nano_cc/internal/sessions\"\n)\n\nconst LocalUserID = \"local-user\"\n\ntype Bundle struct {\n\tRuntime      *runtime.Service\n\tStore        *Store\n\tMCP          *mcp.Manager\n\tUser         storage.User\n\tConversation storage.Co"}
```

## 2026-06-13T09:16:57Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_1h1T9KmhHFrIRJl4smcc1156

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/handlers.go"}
```

### Result

```text
package tools

import (
	"context"
	"fmt"
)

// TodoWriteToolName is the tool that updates the task plan and returns
// structured todo items in addition to a textual summary.
const TodoWriteToolName = "todo_write"

// ExecResult is the unified result of executing a stateless tool. Todos is
// populated only by the todo_write tool.
type ExecResult struct {
	Output string
	Todos  []TodoItem
}

// Handlers maps stateless tool names to their textual handlers. todo_write is
// dispatched separately because it returns structured todos (see Dispatch).
var Handlers = map[string]ToolHandler{
	"bash":                  handleBash,
	"read_file":             handleRead,
	"write_file":            handleWrite,
	"edit_file":             handleEdit,
	"load_skill":            handleLoadSkill,
	"read_persisted_output": handleReadPersistedOutput,
}

// Dispatch is the single entry point for executing a stateless tool by name.
// It is the authority for tool execution semantics, including todo_write's
// structured output.
func Dispatch(ctx context.Context, name string, args map[string]any) (ExecResult, error) {
	if name == TodoWriteToolName {
		result, err := ExecuteTodoWrite(ctx, args)
		if err != nil {
			return ExecResult{}, err
		}
		return ExecResult{Output: result.Output, Todos: result.Todos}, nil
	}
	handler, ok := Handlers[name]
	if !ok || handler == nil {
		return ExecResult{}, fmt.Errorf("tool %s has no handler", name)
	}
	output, err := handler(ctx, args)
	if err != nil {
		return ExecResult{}, err
	}
	return ExecResult{Output: output}, nil
}

func handleBash(ctx context.Context, args map[string]any) (string, error) {
	cmd, _ := args["command"].(string)
	if cmd == "" {
		return "", fmt.Errorf("command is required")
	}
	root, err := validatedWorkspaceRootFromContext(ctx)
	if err != nil {
		return "", err
	}
	workingDir, err := validatedCurrentWorkingDirFromContext(ctx)
	if err != nil {
		return "", err
	}
	if err := validateBashCommandPaths(root, cmd, allowOutsideWorkspaceFromContext(ctx), systemAssetDirsFromContext(ctx)...); err != nil {
		return "", err
	}
	return RunBashInDirWithOptions(cmd, workingDir, allowDangerousCommandsFromContext(ctx))
}

func handleRead(ctx context.Context, args map[string]any) (string, error) {
	path, _ := args["path"].(string)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	limit := 0
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}
	root, resolvedPath, err := resolvePathFromContext(ctx, path)
	if err != nil {
		return "", err
	}
	return RunReadFromRoot(root, resolvedPath, limit)
}

func handleWrite(ctx context.Context, args map[string]any) (string, error) {
	path, _ := args["path"].(string)
	content, _ := args["content"].(string)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	root, resolvedPath, err := resolvePathFromContext(ctx, path)
	if err != nil {
		return "", err
	}
	return RunWriteFromRoot(root, resolvedPath, content)
}

func handleEdit(ctx context.Context, args map[string]any) (string, error) {
	path, _ := args["path"].(string)
	oldText, _ := args["old_text"].(string)
	newText, _ := args["new_text"].(string)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	root, resolvedPath, err := resolvePathFromContext(ctx, path)
	if err != nil {
		return "", err
	}
	return RunEditFromRoot(root, resolvedPath, oldText, newText)
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package tools\n\nimport (\n\t\"context\"\n\t\"fmt\"\n)\n\n// TodoWriteToolName is the tool that updates the task plan and returns\n// structured todo items in addition to a textual summary.\nconst TodoWriteToolName = \"todo_write\"\n\n// ExecResult is the unified result of executing a stateless tool. Todos is\n// populated only by the todo_write tool.\ntype ExecResult struct {\n\tOutput string\n\tTodos  []TodoItem\n}\n\n// Handlers maps stateless tool names to their textual handlers. todo_write is\n// dispatched separately "}
```

## 2026-06-13T09:16:57Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_K10N3nYtTHLfi2DzkAnJ6293

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/context.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/context.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/context.go: no such file or directory"}
```

## 2026-06-13T09:16:57Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_GDn0f0tPrJLgpiYSgC6J5940

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/dispatch.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/dispatch.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/dispatch.go: no such file or directory"}
```

## 2026-06-13T09:16:57Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_03_zXRDJijmW0Jb0Be6bqAs6695

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/types.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/types.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/types.go: no such file or directory"}
```

## 2026-06-13T09:16:57Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_04_Y364Zp8Z6QFyAG05Obib7344

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/validate.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/validate.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/validate.go: no such file or directory"}
```

## 2026-06-13T09:16:57Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_05_3w1z9hWvL53EtT7en7kX3361

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager.go"}
```

### Result

```text
package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/logger"
)

const (
	toolNamePrefix  = "mcp__"
	connectTimeout  = 30 * time.Second
	callTimeout     = 120 * time.Second
	idleTimeout     = 10 * time.Minute
	cleanupInterval = 1 * time.Minute
)

// serverSession 持有一个已连接 MCP 服务器的会话与其工具元信息。
type serverSession struct {
	server      storage.MCPServer
	signature   string // 配置指纹，用于判断配置是否变更需重连
	session     *mcpsdk.ClientSession
	tools       []openai.Tool     // 已转换并加前缀的工具定义
	toolNames   map[string]string // 前缀名 -> 原始 MCP 工具名
	lastUsedAt  time.Time
	activeCalls int
	closing     bool
}

type ServerStatus struct {
	ID        string
	Name      string
	Scope     string
	Transport string
	Command   string
	Args      []string
	URL       string
	Enabled   bool
	Connected bool
	ToolCount int
	LastError string
}

type Snapshot struct {
	Servers    []ServerStatus
	ToolCount  int
	ErrorCount int
}

// Manager 管理本地 MCP 客户端连接，提供工具发现与调用能力。
type Manager struct {
	mu                sync.Mutex
	done              chan struct{}
	closeOnce         sync.Once
	builtinServers    map[string]storage.MCPServer
	builtinSessions   map[string]*serverSession
	workspaceServers  map[string]storage.MCPServer
	workspaceSessions map[string]*serverSession
	workspaceErrors   map[string]string
}

func NewManager() *Manager {
	manager := &Manager{
		done:              make(chan struct{}),
		builtinServers:    make(map[string]storage.MCPServer),
		builtinSessions:   make(map[string]*serverSession),
		workspaceServers:  make(map[string]storage.MCPServer),
		workspaceSessions: make(map[string]*serverSession),
		workspaceErrors:   make(map[string]string),
	}
	go manager.cleanupLoop()
	return manager
}

func serverSignature(s storage.MCPServer) string {
	data, _ := json.Marshal(struct {
		Transport string
		Command   string
		Args      []string
		Env       map[string]string
		URL       string
		Headers   map[string]string
	}{s.Transport, s.Command, s.Args, s.Env, s.URL, s.Headers})
	return string(data)
}

func sanitizeName(name string) string {
	var sb strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_':
			sb.WriteRune(r)
		default:
			sb.WriteRune('_')
		}
	}
	if sb.Len() == 0 {
		return "server"
	}
	return sb.String()
}

func prefixedToolName(serverName, toolName string) string {
	return toolNamePrefix + sanitizeName(serverName) + "__" + toolName
}

func (m *Manager) SetBuiltinServers(ctx context.Context, servers []storage.MCPServer) {
	m.mu.Lock()
	old := m.builtinSessions
	m.builtinServers = make(map[string]storage.MCPServer, len(servers))
	m.builtinSessions = make(map[string]*serverSession)
	for _, server := range servers {
		m.builtinServers[server.ID] = server
	}
	m.mu.Unlock()

	for _, sess := range old {
		closeSession(sess)
	}
	m.EnsureBuiltinSessions(ctx)
}

func (m *Manager) EnsureBuiltinSessions(ctx context.Context) {
	m.mu.Lock()
	missing := make([]storage.MCPServer, 0)
	for id, server := range m.builtinServers {
		if !server.Enabled {
			continue
		}
		if _, ok := m.builtinSessions[id]; !ok {
			missing = append(missing, server)
		}
	}
	m.mu.Unlock()

	for _, server := range missing {
		sess, err := connectBuiltinAndDiscover(ctx, server)
		if err != nil {
			logger.Warn(fmt.Sprintf("mcp: builtin stdio connect failed server=%s: %v", server.Name, err))
			continue
		}
		m.mu.Lock()
		if old := m.builtinSessions[server.ID]; old != nil {
			closeSession(old)
		}
		m.builtinSessions[server.ID] = sess
		m.mu.Unlock()
		logger.Info(fmt.Sprintf("mcp: builtin stdio connected server=%s tools=%d", server.Name, len(sess.tools)))
	}
}

func (m *Manager) SetWorkspaceServers(ctx context.Context, servers []storage.MCPServer) {
	m.mu.Lock()
	old := m.workspaceSessions
	m.workspaceServers = make(map[string]storage.MCPServer, len(servers))
	m.workspaceSessions = make(map[string]*serverSession)
	m.workspaceErrors = make(map[string]string)
	for _, server := range servers {
		m.workspaceServers[server.ID] = server
	}
	m.mu.Unlock()

	for _, sess := range old {
		closeSession(sess)
	}
	m.EnsureWorkspaceSessions(ctx)
}

func (m *Manager) EnsureWorkspaceSessions(ctx context.Context) {
	m.mu.Lock()
	missing := make([]storage.MCPServer, 0)
	for id, server := range m.workspaceServers {
		if !server.Enabled {
			continue
		}
		if _, ok := m.workspaceSessions[id]; !ok {
			missing = append(missing, server)
		}
	}
	m.mu.Unlock()

	for _, server := range missing {
		sess, err := connectAndDiscover(ctx, server)
		m.mu.Lock()
		if err != nil {
			m.workspaceErrors[server.ID] = err.Error()
			m.mu.Unlock()
			logger.Warn(fmt.Sprintf("mcp: workspace connect failed server=%s: %v", server.Name, err))
			continue
		}
		delete(m.workspaceErrors, server.ID)
		if old := m.workspaceSessions[server.ID]; old != nil {
			closeSession(old)
		}
		m.workspaceSessions[server.ID] = sess
		m.mu.Unlock()
		logger.Info(fmt.Sprintf("mcp: workspace connected server=%s tools=%d", server.Name, len(sess.tools)))
	}
}

func connectAndDiscover(ctx context.Context, server storage.MCPServer) (*serverSession, error) {
	transport, err := buildTransport(server)
	if err != nil {
		return nil, err
	}
	connCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	client := mcpsdk.NewClient(&mcpsdk.Implementation{Name: "nano_cc", Version: "1.0.0"}, nil)
	session, err := client.Connect(connCtx, transport, nil)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	tools, names, err := discoverTools(connCtx, session, server.Name)
	if err != nil {
		_ = session.Close()
		return nil, fmt.Errorf("list tools: %w", err)
	}
	return &serverSession{
		server:     server,
		signature:  serverSignature(server),
		session:    session,
		tools:      tools,
		toolNames:  names,
		lastUsedAt: time.Now(),
	}, nil
}

func connectBuiltinAndDiscover(ctx context.Context, server storage.MCPServer) (*serverSession, error) {
	transport, err := buildBuiltinStdioTransport(server)
	if err != nil {
		return nil, err
	}
	return connectWithTransport(ctx, server, transport)
}

func connectWithTransport(ctx context.Context, server storage.MCPServer, transport mcpsdk.Transport) (*serverSession, error) {
	connCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	client := mcpsdk.NewClient(&mcpsdk.Implementation{Name: "nano_cc", Version: "1.0.0"}, nil)
	session, err := client.Connect(connCtx, transport, nil)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	tools, names, err := discoverTools(connCtx, session, server.Name)
	if err != nil {
		_ = session.Close()
		return nil, fmt.Errorf("list tools: %w", err)
	}
	return &serverSession{
		server:     server,
		signature:  serverSignature(server),
		session:    session,
		tools:      tools,
		toolNames:  names,
		lastUsedAt: time.Now(),
	}, nil
}

func discoverTools(ctx context.Context, session *mcpsdk.ClientSession, serverName string) ([]openai.Tool, map[string]string, error) {
	result, err := session.ListTools(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	tools := make([]openai.Tool, 0, len(result.Tools))
	names := make(map[string]string, len(result.Tools))
	for _, tool := range result.Tools {
		prefixed := prefixedToolName(serverName, tool.Name)
		params, err := json.Marshal(tool.InputSchema)
		if err != nil {
			logger.Warn(fmt.Sprintf("mcp: marshal schema failed server=%s tool=%s: %v", serverName, tool.Name, err))
			continue
		}
		tools = append(tools, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        prefixed,
				Description: tool.Description,
				Parameters:  json.RawMessage(params),
			},
		})
		names[prefixed] = tool.Name
	}
	return tools, names, nil
}

// ToolsForUser 返回所有已连接本地 MCP 服务器发现到的工具定义（带前缀）。
func (m *Manager) ToolsForUser(userID string) []openai.Tool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.builtinSessions) == 0 && len(m.workspaceSessions) == 0 {
		return nil
	}
	tools := make([]openai.Tool, 0)
	now := time.Now()
	for _, sess := range m.builtinSessions {
		sess.lastUsedAt = now
		tools = append(tools, sess.tools...)
	}
	for _, sess := range m.workspaceSessions {
		sess.lastUsedAt = now
		tools = append(tools, sess.tools...)
	}
	sort.Slice(tools, func(i, j int) bool {
		return tools[i].Function.Name < tools[j].Function.Name
	})
	return tools
}

// CallTool 把带前缀的工具名路由到对应 MCP 服务器执行，返回文本结果。
func (m *Manager) CallTool(ctx context.Context, userID, prefixedName, rawArgs string) (string, error) {
	m.mu.Lock()
	var target *serverSession
	var originalName string
	var builtinID string
	for _, sess := range m.workspaceSessions {
		if name, ok := sess.toolNames[prefixedName]; ok {
			target = sess
			originalName = name
			target.activeCalls++
			break
		}
	}
	if target == nil {
		for id, sess := range m.builtinSessions {
			if name, ok := sess.toolNames[prefixedName]; ok {
				target = sess
				originalName = name
				builtinID = id
				target.activeCalls++
				break
			}
		}
	}
	m.mu.Unlock()

	if target == nil {
		return "", fmt.Errorf("mcp tool %s not found", prefixedName)
	}
	defer m.finishCall(target)

	var args map[string]any
	if strings.TrimSpace(rawArgs) != "" {
		if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
			return "", fmt.Errorf("invalid tool arguments: %w", err)
		}
	}
	if builtinID != "" {
		return m.callBuiltinTool(ctx, builtinID, target, originalName, args)
	}

	callCtx, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()

	result, err := target.session.CallTool(callCtx, &mcpsdk.CallToolParams{Name: originalName, Arguments: args})
	if err != nil {
		return "", err
	}
	output := serializeContent(result)
	if result.IsError {
		return "", fmt.Errorf("mcp tool error: %s", output)
	}
	return output, nil
}

func (m *Manager) callBuiltinTool(ctx context.Context, builtinID string, target *serverSession, originalName string, args map[string]any) (string, error) {
	output, err := callSessionTool(ctx, target, originalName, args)
	if err == nil {
		return output, nil
	}
	logger.Warn(fmt.Sprintf("mcp: builtin stdio call failed, reconnecting server=%s: %v", target.server.Name, err))
	m.reconnectBuiltin(ctx, builtinID, target)

	m.mu.Lock()
	retryTarget := m.builtinSessions[builtinID]
	if retryTarget != nil {
		retryTarget.activeCalls++
	}
	m.mu.Unlock()
	if retryTarget == nil {
		return "", err
	}
	defer m.finishCall(retryTarget)
	return callSessionTool(ctx, retryTarget, originalName, args)
}

func callSessionTool(ctx context.Context, target *serverSession, originalName string, args map[string]any) (string, error) {
	callCtx, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	result, err := target.session.CallTool(callCtx, &mcpsdk.CallToolParams{Name: originalName, Arguments: args})
	if err != nil {
		return "", err
	}
	output := serializeContent(result)
	if result.IsError {
		return "", fmt.Errorf("mcp tool error: %s", output)
	}
	return output, nil
}

func (m *Manager) reconnectBuiltin(ctx context.Context, builtinID string, old *serverSession) {
	m.mu.Lock()
	server, ok := m.builtinServers[builtinID]
	if current := m.builtinSessions[builtinID]; current == old {
		delete(m.builtinSessions, builtinID)
	}
	m.mu.Unlock()
	if !ok || !server.Enabled {
		return
	}
	closeSession(old)
	sess, err := connectBuiltinAndDiscover(ctx, server)
	if err != nil {
		logger.Warn(fmt.Sprintf("mcp: builtin stdio reconnect failed server=%s: %v", server.Name, err))
		return
	}
	m.mu.Lock()
	m.builtinSessions[builtinID] = sess
	m.mu.Unlock()
	logger.Info(fmt.Sprintf("mcp: builtin stdio reconnected server=%s tools=%d", server.Name, len(sess.tools)))
}

// Invalidate 保留兼容入口。本地 TUI 不再维护用户 DB MCP 连接，因此该方法不影响 builtin/workspace 会话。
func (m *Manager) Invalidate(userID string) {
}

// TestServer 临时连接一个配置并发现工具，返回工具名列表，供连接测试使用。
func (m *Manager) TestServer(ctx context.Context, server storage.MCPServer) ([]string, error) {
	sess, err := connectAndDiscover(ctx, server)
	if err != nil {
		return nil, err
	}
	defer sess.session.Close()
	names := make([]string, 0, len(sess.toolNames))
	for _, original := range sess.toolNames {
		names = append(names, original)
	}
	sort.Strings(names)
	return names, nil
}

// Close 关闭全部连接，用于服务退出。
func (m *Manager) Close() {
	m.closeOnce.Do(func() { close(m.done) })
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, sess := range m.builtinSessions {
		m.closeWhenIdleLocked(sess)
	}
	for _, sess := range m.workspaceSessions {
		m.closeWhenIdleLocked(sess)
	}
	m.builtinSessions = make(map[string]*serverSession)
	m.workspaceSessions = make(map[string]*serverSession)
}

func (m *Manager) Snapshot(userID string) Snapshot {
	m.mu.Lock()
	defer m.mu.Unlock()
	statuses := make([]ServerStatus, 0, len(m.builtinServers)+len(m.workspaceServers))
	addConfigured := func(scope string, server storage.MCPServer, sess *serverSession, lastErr string) {
		status := ServerStatus{ID: server.ID, Name: server.Name, Scope: scope, Transport: server.Transport, Command: server.Command, Args: append([]string(nil), server.Args...), URL: server.URL, Enabled: server.Enabled, Connected: sess != nil, LastError: lastErr}
		if sess != nil {
			status.ToolCount = len(sess.tools)
		}
		statuses = append(statuses, status)
	}
	for id, server := range m.builtinServers {
		addConfigured("builtin", server, m.builtinSessions[id], "")
	}
	for id, server := range m.workspaceServers {
		addConfigured("workspace", server, m.workspaceSessions[id], m.workspaceErrors[id])
	}
	sort.Slice(statuses, func(i, j int) bool {
		if statuses[i].Scope != statuses[j].Scope {
			return statuses[i].Scope < statuses[j].Scope
		}
		return statuses[i].Name < statuses[j].Name
	})
	snapshot := Snapshot{Servers: statuses}
	for _, status := range statuses {
		snapshot.ToolCount += status.ToolCount
		if status.LastError != "" {
			snapshot.ErrorCount++
		}
	}
	return snapshot
}

func (m *Manager) finishCall(sess *serverSession) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if sess.activeCalls > 0 {
		sess.activeCalls--
	}
	sess.lastUsedAt = time.Now()
	if sess.closing && sess.activeCalls == 0 {
		closeSession(sess)
	}
}

func (m *Manager) closeWhenIdleLocked(sess *serverSession) {
	if sess == nil {
		return
	}
	sess.closing = true
	if sess.activeCalls == 0 {
		closeSession(sess)
	}
}

func (m *Manager) cleanupLoop() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			m.cleanupIdleSessions(time.Now())
		case <-m.done:
			return
		}
	}
}

func (m *Manager) cleanupIdleSessions(now time.Time) {
	type idleSession struct {
		scope string
		id    string
		sess  *serverSession
	}
	var idle []idleSession

	m.mu.Lock()
	for serverID, sess := range m.workspaceSessions {
		if sess.activeCalls > 0 || now.Sub(sess.lastUsedAt) < idleTimeout {
			continue
		}
		idle = append(idle, idleSession{scope: "workspace", id: serverID, sess: sess})
		delete(m.workspaceSessions, serverID)
	}
	m.mu.Unlock()

	for _, item := range idle {
		closeSession(item.sess)
		logger.Info(fmt.Sprintf("mcp: idle session closed scope=%s id=%s server=%s", item.scope, item.id, item.sess.server.Name))
	}
}

func closeSession(sess *serverSession) {
	if sess == nil || sess.session == nil {
		return
	}
	_ = sess.session.Close()
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package mcp\n\nimport (\n\t\"context\"\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"sort\"\n\t\"strings\"\n\t\"sync\"\n\t\"time\"\n\n\tmcpsdk \"github.com/modelcontextprotocol/go-sdk/mcp\"\n\topenai \"github.com/sashabaranov/go-openai\"\n\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/logger\"\n)\n\nconst (\n\ttoolNamePrefix  = \"mcp__\"\n\tconnectTimeout  = 30 * time.Second\n\tcallTimeout     = 120 * time.Second\n\tidleTimeout     = 10 * time.Minute\n\tcleanupInterval = 1 * time.Minute\n)\n\n// serverSession 持有一个已连接 MCP 服务器的会话\ufffd\ufffd"}
```

## 2026-06-13T09:16:57Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_06_oa2vbYl6aMJDUglv1tg80945

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/server.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/server.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/server.go: no such file or directory"}
```

## 2026-06-13T09:16:59Z · bash · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_yQSxVk2vpJVYbSXATpJB8944

### Arguments

```json
{"command": "ls -la /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/"}
```

### Result

```text
total 96
drwxr-xr-x@ 8 bytedance  staff    256  6 12 16:08 .
drwxr-xr-x@ 5 bytedance  staff    160  6 12 12:42 ..
-rw-r--r--@ 1 bytedance  staff   6005  6 12 13:09 config_test.go
-rw-r--r--@ 1 bytedance  staff   6549  6 12 13:14 config.go
-rw-r--r--@ 1 bytedance  staff   6432  6 12 16:08 manager_test.go
-rw-r--r--@ 1 bytedance  staff  15809  6 12 16:10 manager.go
-rw-r--r--@ 1 bytedance  staff   1831  6 12 13:09 transport_test.go
-rw-r--r--@ 1 bytedance  staff   3167  6 12 13:10 transport.go
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","resolved_command_path":"/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp","outcome_summary":"total 96\ndrwxr-xr-x@ 8 bytedance  staff    256  6 12 16:08 .\ndrwxr-xr-x@ 5 bytedance  staff    160  6 12 12:42 ..\n-rw-r--r--@ 1 bytedance  staff   6005  6 12 13:09 config_test.go\n-rw-r--r--@ 1 bytedance  staff   6549  6 12 13:14 config.go\n-rw-r--r--@ 1 bytedance  staff   6432  6 12 16:08 manager_test.go\n-rw-r--r--@ 1 bytedance  staff  15809  6 12 16:10 manager.go\n-rw-r--r--@ 1 bytedance  staff   1831  6 12 13:09 transport_test.go\n-rw-r--r--@ 1 bytedance  staff   3167  6 12 13:10 transport.go"}
```

## 2026-06-13T09:17:02Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_jkT3c8qU8YJQF0lWTVpg2134

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config.go"}
```

### Result

```text
package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/logger"
)

type builtinConfig struct {
	Servers *[]builtinServerConfig `json:"mcp_servers"`
}

type builtinServerConfig struct {
	Name      string            `json:"name"`
	Command   string            `json:"command"`
	Args      []string          `json:"args"`
	Env       map[string]string `json:"env"`
	Enabled   *bool             `json:"enabled"`
	Transport *string           `json:"transport"`
	URL       *string           `json:"url"`
	Headers   map[string]string `json:"headers"`
}

type workspaceConfig struct {
	ServersArray *[]workspaceServerConfig         `json:"mcp_servers"`
	ServersMap   map[string]workspaceServerConfig `json:"mcpServers"`
}

type workspaceServerConfig struct {
	Name      string            `json:"name"`
	Command   string            `json:"command"`
	Args      []string          `json:"args"`
	Env       map[string]string `json:"env"`
	Enabled   *bool             `json:"enabled"`
	Transport string            `json:"transport"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
}

func LoadBuiltinConfig(path string) ([]storage.MCPServer, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Info(fmt.Sprintf("mcp: builtin config not found path=%s, skip builtin servers", path))
			return nil, nil
		}
		return nil, err
	}

	var cfg builtinConfig
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("parse mcp_config.json: %w", err)
	}
	if cfg.Servers == nil {
		return nil, fmt.Errorf("mcp_config.json requires mcp_servers")
	}

	servers := make([]storage.MCPServer, 0, len(*cfg.Servers))
	seen := make(map[string]struct{}, len(*cfg.Servers))
	for i, item := range *cfg.Servers {
		server, key, err := builtinServer(i, item)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[key]; ok {
			return nil, fmt.Errorf("builtin mcp server %d has duplicate sanitized name %q", i, key)
		}
		seen[key] = struct{}{}
		servers = append(servers, server)
	}
	logger.Info(fmt.Sprintf("mcp: loaded builtin stdio servers count=%d", len(servers)))
	return servers, nil
}

func builtinServer(index int, item builtinServerConfig) (storage.MCPServer, string, error) {
	if item.Transport != nil {
		return storage.MCPServer{}, "", fmt.Errorf("builtin mcp server %d must not set transport; builtin mcp only supports stdio", index)
	}
	if item.URL != nil {
		return storage.MCPServer{}, "", fmt.Errorf("builtin mcp server %d must not set url; builtin mcp only supports stdio", index)
	}
	if item.Headers != nil {
		return storage.MCPServer{}, "", fmt.Errorf("builtin mcp server %d must not set headers; builtin mcp only supports stdio", index)
	}

	name := strings.TrimSpace(item.Name)
	if name == "" {
		return storage.MCPServer{}, "", fmt.Errorf("builtin mcp server %d name is required", index)
	}
	command := strings.TrimSpace(item.Command)
	if command == "" {
		return storage.MCPServer{}, "", fmt.Errorf("builtin mcp server %s command is required", name)
	}

	enabled := true
	if item.Enabled != nil {
		enabled = *item.Enabled
	}
	key := sanitizeName(name)
	return storage.MCPServer{
		ID:        "builtin:" + key,
		UserID:    "",
		Name:      "builtin_" + name,
		Transport: "stdio",
		Command:   command,
		Args:      item.Args,
		Env:       item.Env,
		Enabled:   enabled,
	}, key, nil
}

func LoadWorkspaceConfig(path string) ([]storage.MCPServer, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Info(fmt.Sprintf("mcp: workspace config not found path=%s, skip workspace servers", path))
			return nil, nil
		}
		return nil, err
	}

	var cfg workspaceConfig
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("parse workspace mcp config %s: %w", path, err)
	}

	items := make([]workspaceServerConfig, 0)
	if cfg.ServersArray != nil {
		items = append(items, (*cfg.ServersArray)...)
	}
	if cfg.ServersMap != nil {
		keys := make([]string, 0, len(cfg.ServersMap))
		for key := range cfg.ServersMap {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			item := cfg.ServersMap[key]
			if strings.TrimSpace(item.Name) == "" {
				item.Name = key
			}
			items = append(items, item)
		}
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("workspace mcp config %s requires mcp_servers or mcpServers", path)
	}

	servers := make([]storage.MCPServer, 0, len(items))
	seen := make(map[string]struct{}, len(items))
	for i, item := range items {
		server, key, err := workspaceServer(i, item)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[key]; ok {
			return nil, fmt.Errorf("workspace mcp server %d has duplicate sanitized name %q", i, key)
		}
		seen[key] = struct{}{}
		servers = append(servers, server)
	}
	logger.Info(fmt.Sprintf("mcp: loaded workspace servers count=%d", len(servers)))
	return servers, nil
}

func workspaceServer(index int, item workspaceServerConfig) (storage.MCPServer, string, error) {
	name := strings.TrimSpace(item.Name)
	if name == "" {
		return storage.MCPServer{}, "", fmt.Errorf("workspace mcp server %d name is required", index)
	}
	transport := strings.TrimSpace(item.Transport)
	if transport == "" && strings.TrimSpace(item.Command) != "" {
		transport = "stdio"
	}
	if transport == "" {
		return storage.MCPServer{}, "", fmt.Errorf("workspace mcp server %s transport is required", name)
	}
	enabled := true
	if item.Enabled != nil {
		enabled = *item.Enabled
	}
	command := strings.TrimSpace(item.Command)
	url := strings.TrimSpace(item.URL)
	switch transport {
	case "stdio":
		if command == "" {
			return storage.MCPServer{}, "", fmt.Errorf("workspace mcp server %s command is required for stdio", name)
		}
	case "sse", "streamable":
		if url == "" {
			return storage.MCPServer{}, "", fmt.Errorf("workspace mcp server %s url is required for %s", name, transport)
		}
	default:
		return storage.MCPServer{}, "", fmt.Errorf("workspace mcp server %s unsupported transport %q", name, transport)
	}
	key := sanitizeName(name)
	return storage.MCPServer{
		ID:        "workspace:" + key,
		Name:      name,
		Transport: transport,
		Command:   command,
		Args:      item.Args,
		Env:       item.Env,
		URL:       url,
		Headers:   item.Headers,
		Enabled:   enabled,
	}, key, nil
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package mcp\n\nimport (\n\t\"bytes\"\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"os\"\n\t\"sort\"\n\t\"strings\"\n\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/logger\"\n)\n\ntype builtinConfig struct {\n\tServers *[]builtinServerConfig `json:\"mcp_servers\"`\n}\n\ntype builtinServerConfig struct {\n\tName      string            `json:\"name\"`\n\tCommand   string            `json:\"command\"`\n\tArgs      []string          `json:\"args\"`\n\tEnv       map[string]string `json:\"env\"`\n\tEnabled   *bool             `json:\"enabled\"`\n\tTransport *strin"}
```

## 2026-06-13T09:17:02Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_JJ1CMnN8V4tSxfn3Ytgp0771

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/transport.go"}
```

### Result

```text
package mcp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"nano_cc/internal/agent/storage"
)

// headerRoundTripper 在每个 HTTP 请求上附加自定义请求头，用于 sse/streamable 传输鉴权。
type headerRoundTripper struct {
	base    http.RoundTripper
	headers map[string]string
}

func (t *headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range t.headers {
		if strings.TrimSpace(key) == "" {
			continue
		}
		req.Header.Set(key, value)
	}
	return t.base.RoundTrip(req)
}

func httpClientWithHeaders(headers map[string]string) *http.Client {
	if len(headers) == 0 {
		return nil
	}
	return &http.Client{Transport: &headerRoundTripper{base: http.DefaultTransport, headers: headers}}
}

// buildTransport 根据配置的 transport 类型构造对应的 MCP 客户端传输层。
func buildTransport(server storage.MCPServer) (mcpsdk.Transport, error) {
	switch server.Transport {
	case "stdio":
		return buildStdioTransport(server)
	case "sse":
		if strings.TrimSpace(server.URL) == "" {
			return nil, fmt.Errorf("sse transport requires a url")
		}
		return &mcpsdk.SSEClientTransport{Endpoint: server.URL, HTTPClient: httpClientWithHeaders(server.Headers)}, nil
	case "streamable":
		if strings.TrimSpace(server.URL) == "" {
			return nil, fmt.Errorf("streamable transport requires a url")
		}
		return &mcpsdk.StreamableClientTransport{Endpoint: server.URL, HTTPClient: httpClientWithHeaders(server.Headers)}, nil
	default:
		return nil, fmt.Errorf("unsupported transport %q", server.Transport)
	}
}

func buildBuiltinStdioTransport(server storage.MCPServer) (mcpsdk.Transport, error) {
	return buildStdioTransport(server)
}

func buildStdioTransport(server storage.MCPServer) (mcpsdk.Transport, error) {
	command := strings.TrimSpace(server.Command)
	if command == "" {
		return nil, fmt.Errorf("stdio transport requires a command")
	}
	cmd := exec.Command(command, server.Args...)
	cmd.Env = append(os.Environ(), envPairs(server.Env)...)
	return &mcpsdk.CommandTransport{Command: cmd}, nil
}

func envPairs(env map[string]string) []string {
	if len(env) == 0 {
		return nil
	}
	pairs := make([]string, 0, len(env))
	for key, value := range env {
		if strings.TrimSpace(key) == "" {
			continue
		}
		pairs = append(pairs, key+"="+value)
	}
	return pairs
}

// serializeContent 把 MCP 工具调用返回的 content blocks 拼成纯文本，契合现有工具结果的字符串契约。
func serializeContent(result *mcpsdk.CallToolResult) string {
	if result == nil {
		return ""
	}
	var sb strings.Builder
	for _, content := range result.Content {
		switch c := content.(type) {
		case *mcpsdk.TextContent:
			if sb.Len() > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(c.Text)
		default:
			data, err := json.Marshal(content)
			if err != nil {
				continue
			}
			if sb.Len() > 0 {
				sb.WriteString("\n")
			}
			sb.Write(data)
		}
	}
	if sb.Len() == 0 && result.StructuredContent != nil {
		if data, err := json.Marshal(result.StructuredContent); err == nil {
			sb.Write(data)
		}
	}
	return sb.String()
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package mcp\n\nimport (\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"net/http\"\n\t\"os\"\n\t\"os/exec\"\n\t\"strings\"\n\n\tmcpsdk \"github.com/modelcontextprotocol/go-sdk/mcp\"\n\n\t\"nano_cc/internal/agent/storage\"\n)\n\n// headerRoundTripper 在每个 HTTP 请求上附加自定义请求头，用于 sse/streamable 传输鉴权。\ntype headerRoundTripper struct {\n\tbase    http.RoundTripper\n\theaders map[string]string\n}\n\nfunc (t *headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {\n\tfor key, value := range t.headers {\n"}
```

## 2026-06-13T09:17:02Z · bash · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_31uG4WhvmRfJG5JFQYgQ4476

### Arguments

```json
{"command": "find /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage -type f -name \"*.go\" -not -name \"*_test.go\" | sort"}
```

### Result

```text
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/conversation_history.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/models.go
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","resolved_command_path":"/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage","outcome_summary":"/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/conversation_history.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/models.go"}
```

## 2026-06-13T09:17:07Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_Y8fl1PQPI7Su0a4739mT1227

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/models.go"}
```

### Result

```text
package storage

import "time"

type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	MemoryEnabled bool      `json:"memory_enabled"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// MCPServer 是用户配置的一个 MCP 服务器连接。Transport 取值 stdio/sse/streamable。
// Args/Env/Headers 从本地 .link/.mcp.json 读取，在 Go 侧以强类型表示。
type MCPServer struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Name      string            `json:"name"`
	Transport string            `json:"transport"`
	Command   string            `json:"command"`
	Args      []string          `json:"args"`
	Env       map[string]string `json:"env"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Enabled   bool              `json:"enabled"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type Conversation struct {
	ID            string    `json:"id"`
	SessionID     string    `json:"session_id,omitempty"`
	UserID        string    `json:"user_id"`
	RootMessageID string    `json:"root_message_id"`
	Title         string    `json:"title"`
	HistoryJSON   string    `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ResumableSession struct {
	SessionID      string    `json:"session_id"`
	ConversationID string    `json:"conversation_id"`
	WorkspaceRoot  string    `json:"workspace_root"`
	Title          string    `json:"title"`
	MessageCount   int       `json:"message_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// 该表已经不使用，作为存储的结构体，用于存储消息历史
type Message struct {
	ID               string            `json:"id"`
	ConversationID   string            `json:"conversation_id"`
	UserID           string            `json:"user_id"`
	Role             string            `json:"role"`
	Content          string            `json:"content"`
	ReasoningContent string            `json:"reasoning_content,omitempty"`
	ToolCallID       string            `json:"tool_call_id,omitempty"`
	ToolCalls        []MessageToolCall `json:"tool_calls,omitempty"`
	Meta             *MessageMeta      `json:"meta,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
}

// MessageMeta 记录助手回复的元信息，仅对最终 assistant 消息填充。
type MessageMeta struct {
	ToolCallCount int `json:"tool_call_count"`          // 调用工具次数（0 也需序列化，保证历史展示一致）
	ContextTokens int `json:"context_tokens,omitempty"` // 当前上下文估算 token
	ContextBudget int `json:"context_budget,omitempty"` // 上下文预算（用于算占比）
}

type MessageToolCall struct {
	ID       string              `json:"id"`
	Type     string              `json:"type"`
	Function MessageFunctionCall `json:"function"`
}

type MessageFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolCall struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	ToolName       string    `json:"tool_name"`
	Status         string    `json:"status"`
	Summary        string    `json:"summary"`
	CreatedAt      time.Time `json:"created_at"`
}

type PersistedOutput struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	MessageID      string    `json:"message_id"`
	ToolCallID     string    `json:"tool_call_id"`
	Kind           string    `json:"kind"`
	Strategy       string    `json:"strategy"`
	OriginalBytes  int       `json:"original_bytes"`
	ContentSHA256  string    `json:"content_sha256"`
	Content        string    `json:"content"`
	Preview        string    `json:"preview"`
	CreatedAt      time.Time `json:"created_at"`
}

type ToolResultLogEntry struct {
	ConversationID string    `json:"conversation_id"`
	SessionID      string    `json:"session_id,omitempty"`
	UserID         string    `json:"user_id"`
	ToolCallID     string    `json:"tool_call_id"`
	ToolName       string    `json:"tool_name"`
	RawArgs        string    `json:"raw_args"`
	Status         string    `json:"status"`
	Result         string    `json:"result"`
	AuditSummary   string    `json:"audit_summary,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type ContextSummary struct {
	ID                    string    `json:"id"`
	ConversationID        string    `json:"conversation_id"`
	UserID                string    `json:"user_id"`
	SourceHistorySHA256   string    `json:"source_history_sha256"`
	Strategy              string    `json:"strategy"`
	EstimatedTokensBefore int       `json:"estimated_tokens_before"`
	EstimatedTokensAfter  int       `json:"estimated_tokens_after"`
	Summary               string    `json:"summary"`
	CreatedAt             time.Time `json:"created_at"`
}

type SubagentMessage struct {
	ID               string            `json:"id"`
	RunID            string            `json:"run_id"`
	ParentToolCallID string            `json:"parent_tool_call_id"`
	ConversationID   string            `json:"conversation_id"`
	UserID           string            `json:"user_id"`
	SequenceNo       int               `json:"sequence_no"`
	Role             string            `json:"role"`
	Content          string            `json:"content"`
	ReasoningContent string            `json:"reasoning_content,omitempty"`
	ToolCallID       string            `json:"tool_call_id,omitempty"`
	ToolCalls        []MessageToolCall `json:"tool_calls,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
}

type Memory struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ConversationMemory 是单个会话维度、随每轮对话增量维护的“当前会话主干信息”
// 条目。它不注入 system prompt，仅在上下文压缩触发全量摘要时作为替代品使用。
type ConversationMemory struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package storage\n\nimport \"time\"\n\ntype User struct {\n\tID            string    `json:\"id\"`\n\tEmail         string    `json:\"email\"`\n\tUsername      string    `json:\"username\"`\n\tMemoryEnabled bool      `json:\"memory_enabled\"`\n\tCreatedAt     time.Time `json:\"created_at\"`\n\tUpdatedAt     time.Time `json:\"updated_at\"`\n}\n\n// MCPServer 是用户配置的一个 MCP 服务器连接。Transport 取值 stdio/sse/streamable。\n// Args/Env/Headers 从本地 .link/.mcp.json 读取，在 Go 侧以强类型表示。"}
```

## 2026-06-13T09:17:07Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_vRztTvptZomak6DZiI3P9789

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/conversation_history.go"}
```

### Result

```text
package storage

import (
	"encoding/json"
	"strings"
)

const emptyConversationHistoryJSON = "[]"

func EncodeConversationHistory(messages []Message) (string, error) {
	if messages == nil {
		messages = []Message{}
	}
	data, err := json.Marshal(messages)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func DecodeConversationHistory(historyJSON string) ([]Message, error) {
	if strings.TrimSpace(historyJSON) == "" {
		historyJSON = emptyConversationHistoryJSON
	}
	var messages []Message
	if err := json.Unmarshal([]byte(historyJSON), &messages); err != nil {
		return nil, err
	}
	if messages == nil {
		messages = []Message{}
	}
	return messages, nil
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package storage\n\nimport (\n\t\"encoding/json\"\n\t\"strings\"\n)\n\nconst emptyConversationHistoryJSON = \"[]\"\n\nfunc EncodeConversationHistory(messages []Message) (string, error) {\n\tif messages == nil {\n\t\tmessages = []Message{}\n\t}\n\tdata, err := json.Marshal(messages)\n\tif err != nil {\n\t\treturn \"\", err\n\t}\n\treturn string(data), nil\n}\n\nfunc DecodeConversationHistory(historyJSON string) ([]Message, error) {\n\tif strings.TrimSpace(historyJSON) == \"\" {\n\t\thistoryJSON = emptyConversationHistoryJSON\n\t}\n\tvar messages ["}
```

## 2026-06-13T09:17:07Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_9HVtLHSyA6useuiMYE5o7946

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/config.go"}
```

### Result

```text
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"-"`
	ModelID string `json:"model_id"`
}

type fileConfig struct {
	Config
	AppHome                            string `json:"app_home"`
	BuiltinSkillsDir                   string `json:"builtin_skills_dir"`
	CommandBinDir                      string `json:"command_bin_dir"`
	CommandScriptDir                   string `json:"command_script_dir"`
	SystemPromptPath                   string `json:"system_prompt_path"`
	WorkspaceRoot                      string `json:"workspace_root"`
	AllowedTools                       string `json:"allowed_tools"`
	BashAllowOutsideWorkspace          bool   `json:"bash_allow_outside_workspace"`
	BashAllowDangerousCommands         bool   `json:"bash_allow_dangerous_commands"`
	ConversationLockTTLSeconds         int    `json:"conversation_lock_ttl_seconds"`
	MemoryWorkTimeoutSeconds           int    `json:"memory_work_timeout_seconds"`
	ConversationLockWaitTimeoutSeconds int    `json:"conversation_lock_wait_timeout_seconds"`
}

type AppConfig struct {
	LLM                         Config
	AppHome                     string
	BuiltinSkillsDir            string
	CommandBinDir               string
	CommandScriptDir            string
	SystemPromptPath            string
	WorkspaceRoot               string
	LogsDir                     string
	AllowedTools                []string
	BashAllowOutsideWorkspace   bool
	BashAllowDangerousCommands  bool
	ConversationLockTTL         time.Duration
	MemoryWorkTimeout           time.Duration
	ConversationLockWaitTimeout time.Duration
}

func loadConfigFile() (fileConfig, error) {
	data, err := os.ReadFile(configFilePath())
	if err != nil {
		return fileConfig{}, err
	}

	var cfg fileConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fileConfig{}, fmt.Errorf("failed to parse workspace config.json: %w", err)
	}
	return cfg, nil
}

func configFilePath() string {
	return "config.json"
}

func loadLLMConfig(fileCfg fileConfig) (Config, error) {
	cfg := Config{
		BaseURL: strings.TrimSpace(fileCfg.Config.BaseURL),
		APIKey:  strings.TrimSpace(getenv("OPENAI_API_KEY")),
		ModelID: strings.TrimSpace(fileCfg.Config.ModelID),
	}

	if cfg.BaseURL == "" || cfg.ModelID == "" {
		return Config{}, fmt.Errorf("missing LLM config; set base_url, model_id in config.json")
	}
	if cfg.APIKey == "" {
		return Config{}, fmt.Errorf("missing LLM api key; set OPENAI_API_KEY environment variable")
	}

	return cfg, nil
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package config\n\nimport (\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"os\"\n\t\"strings\"\n\t\"time\"\n)\n\ntype Config struct {\n\tBaseURL string `json:\"base_url\"`\n\tAPIKey  string `json:\"-\"`\n\tModelID string `json:\"model_id\"`\n}\n\ntype fileConfig struct {\n\tConfig\n\tAppHome                            string `json:\"app_home\"`\n\tBuiltinSkillsDir                   string `json:\"builtin_skills_dir\"`\n\tCommandBinDir                      string `json:\"command_bin_dir\"`\n\tCommandScriptDir                   string `json:\"command_script_dir\"`\n"}
```

## 2026-06-13T09:17:07Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_03_Z3149kGX4Vy0mXYueJRJ4716

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/hooks/hooks.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/hooks/hooks.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/hooks/hooks.go: no such file or directory"}
```

## 2026-06-13T09:17:07Z · bash · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_04_FnKDuIgGySyk9oIkn02J1058

### Arguments

```json
{"command": "find /Users/bytedance/golang_pro/nano_cc/go-agent/internal -type f -name \"*.go\" -not -name \"*_test.go\" -not -path \"*/vendor/*\" -not -path \"*/agent/runtime/compression/strategies/*\" | sort"}
```

### Result

```text
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/transport.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/conversation_memory_strategy.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/full_history_summarization.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/message_window_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/reactive_compact.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/recent_tool_result_retention.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/token_estimator.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/tool_result_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/context_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/conversation_flow.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/conversation_memory.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/helpers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hook_bridge.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/manager.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/stop.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/tool.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/types.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/user_prompt.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/memory.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/prompt_builder.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/runtime.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/subagent.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/summarizer.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/todo_reminder.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/tool_registry.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/conversation_history.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/models.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/assistant/prompt.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/cli/root.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/helpers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/local_config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/paths.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/idgen/idgen.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/llm/client.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/llm/errors.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/bootstrap.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/memory_store.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/persisted_output_files.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/session_history_store.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/store.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/tool_result_log.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/logger/logger.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/safety/path.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/sessions/skill.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/textutil/textutil.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/bash.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/definitions.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/file_ops.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/handlers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/load_skill.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/path_guard.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/persisted_output.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/runtime_env.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/todo_write.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/validation.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/app.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/events.go
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","resolved_command_path":"/Users/bytedance/golang_pro/nano_cc/go-agent/internal","outcome_summary":"/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/transport.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/compression.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/conversation_memory_strategy.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime"}
```

## 2026-06-13T09:17:11Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_RWsw3Qx66GfruvrgMR8S2511

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/models.go"}
```

### Result

```text
package storage

import "time"

type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	MemoryEnabled bool      `json:"memory_enabled"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// MCPServer 是用户配置的一个 MCP 服务器连接。Transport 取值 stdio/sse/streamable。
// Args/Env/Headers 从本地 .link/.mcp.json 读取，在 Go 侧以强类型表示。
type MCPServer struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Name      string            `json:"name"`
	Transport string            `json:"transport"`
	Command   string            `json:"command"`
	Args      []string          `json:"args"`
	Env       map[string]string `json:"env"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Enabled   bool              `json:"enabled"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type Conversation struct {
	ID            string    `json:"id"`
	SessionID     string    `json:"session_id,omitempty"`
	UserID        string    `json:"user_id"`
	RootMessageID string    `json:"root_message_id"`
	Title         string    `json:"title"`
	HistoryJSON   string    `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ResumableSession struct {
	SessionID      string    `json:"session_id"`
	ConversationID string    `json:"conversation_id"`
	WorkspaceRoot  string    `json:"workspace_root"`
	Title          string    `json:"title"`
	MessageCount   int       `json:"message_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// 该表已经不使用，作为存储的结构体，用于存储消息历史
type Message struct {
	ID               string            `json:"id"`
	ConversationID   string            `json:"conversation_id"`
	UserID           string            `json:"user_id"`
	Role             string            `json:"role"`
	Content          string            `json:"content"`
	ReasoningContent string            `json:"reasoning_content,omitempty"`
	ToolCallID       string            `json:"tool_call_id,omitempty"`
	ToolCalls        []MessageToolCall `json:"tool_calls,omitempty"`
	Meta             *MessageMeta      `json:"meta,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
}

// MessageMeta 记录助手回复的元信息，仅对最终 assistant 消息填充。
type MessageMeta struct {
	ToolCallCount int `json:"tool_call_count"`          // 调用工具次数（0 也需序列化，保证历史展示一致）
	ContextTokens int `json:"context_tokens,omitempty"` // 当前上下文估算 token
	ContextBudget int `json:"context_budget,omitempty"` // 上下文预算（用于算占比）
}

type MessageToolCall struct {
	ID       string              `json:"id"`
	Type     string              `json:"type"`
	Function MessageFunctionCall `json:"function"`
}

type MessageFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolCall struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	ToolName       string    `json:"tool_name"`
	Status         string    `json:"status"`
	Summary        string    `json:"summary"`
	CreatedAt      time.Time `json:"created_at"`
}

type PersistedOutput struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	MessageID      string    `json:"message_id"`
	ToolCallID     string    `json:"tool_call_id"`
	Kind           string    `json:"kind"`
	Strategy       string    `json:"strategy"`
	OriginalBytes  int       `json:"original_bytes"`
	ContentSHA256  string    `json:"content_sha256"`
	Content        string    `json:"content"`
	Preview        string    `json:"preview"`
	CreatedAt      time.Time `json:"created_at"`
}

type ToolResultLogEntry struct {
	ConversationID string    `json:"conversation_id"`
	SessionID      string    `json:"session_id,omitempty"`
	UserID         string    `json:"user_id"`
	ToolCallID     string    `json:"tool_call_id"`
	ToolName       string    `json:"tool_name"`
	RawArgs        string    `json:"raw_args"`
	Status         string    `json:"status"`
	Result         string    `json:"result"`
	AuditSummary   string    `json:"audit_summary,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type ContextSummary struct {
	ID                    string    `json:"id"`
	ConversationID        string    `json:"conversation_id"`
	UserID                string    `json:"user_id"`
	SourceHistorySHA256   string    `json:"source_history_sha256"`
	Strategy              string    `json:"strategy"`
	EstimatedTokensBefore int       `json:"estimated_tokens_before"`
	EstimatedTokensAfter  int       `json:"estimated_tokens_after"`
	Summary               string    `json:"summary"`
	CreatedAt             time.Time `json:"created_at"`
}

type SubagentMessage struct {
	ID               string            `json:"id"`
	RunID            string            `json:"run_id"`
	ParentToolCallID string            `json:"parent_tool_call_id"`
	ConversationID   string            `json:"conversation_id"`
	UserID           string            `json:"user_id"`
	SequenceNo       int               `json:"sequence_no"`
	Role             string            `json:"role"`
	Content          string            `json:"content"`
	ReasoningContent string            `json:"reasoning_content,omitempty"`
	ToolCallID       string            `json:"tool_call_id,omitempty"`
	ToolCalls        []MessageToolCall `json:"tool_calls,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
}

type Memory struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ConversationMemory 是单个会话维度、随每轮对话增量维护的“当前会话主干信息”
// 条目。它不注入 system prompt，仅在上下文压缩触发全量摘要时作为替代品使用。
type ConversationMemory struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package storage\n\nimport \"time\"\n\ntype User struct {\n\tID            string    `json:\"id\"`\n\tEmail         string    `json:\"email\"`\n\tUsername      string    `json:\"username\"`\n\tMemoryEnabled bool      `json:\"memory_enabled\"`\n\tCreatedAt     time.Time `json:\"created_at\"`\n\tUpdatedAt     time.Time `json:\"updated_at\"`\n}\n\n// MCPServer 是用户配置的一个 MCP 服务器连接。Transport 取值 stdio/sse/streamable。\n// Args/Env/Headers 从本地 .link/.mcp.json 读取，在 Go 侧以强类型表示。"}
```

## 2026-06-13T09:17:11Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_CvC8fxMaeTyNEK6nGEpu8714

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/conversation_history.go"}
```

### Result

```text
package storage

import (
	"encoding/json"
	"strings"
)

const emptyConversationHistoryJSON = "[]"

func EncodeConversationHistory(messages []Message) (string, error) {
	if messages == nil {
		messages = []Message{}
	}
	data, err := json.Marshal(messages)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func DecodeConversationHistory(historyJSON string) ([]Message, error) {
	if strings.TrimSpace(historyJSON) == "" {
		historyJSON = emptyConversationHistoryJSON
	}
	var messages []Message
	if err := json.Unmarshal([]byte(historyJSON), &messages); err != nil {
		return nil, err
	}
	if messages == nil {
		messages = []Message{}
	}
	return messages, nil
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package storage\n\nimport (\n\t\"encoding/json\"\n\t\"strings\"\n)\n\nconst emptyConversationHistoryJSON = \"[]\"\n\nfunc EncodeConversationHistory(messages []Message) (string, error) {\n\tif messages == nil {\n\t\tmessages = []Message{}\n\t}\n\tdata, err := json.Marshal(messages)\n\tif err != nil {\n\t\treturn \"\", err\n\t}\n\treturn string(data), nil\n}\n\nfunc DecodeConversationHistory(historyJSON string) ([]Message, error) {\n\tif strings.TrimSpace(historyJSON) == \"\" {\n\t\thistoryJSON = emptyConversationHistoryJSON\n\t}\n\tvar messages ["}
```

## 2026-06-13T09:17:11Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_fxLV2x68a870DW6crOvA1505

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/types.go"}
```

### Result

```text
package hooks

import (
	"context"
	"encoding/json"
	"fmt"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	agenttools "nano_cc/internal/tools"
)

type Store interface {
	UpdateConversationTitle(ctx context.Context, conversationID, title string) error
	TouchConversationActivity(ctx context.Context, conversationID string) error
	SetConversationHistory(ctx context.Context, conversationID string, messages []storage.Message) error
	CreateToolCall(ctx context.Context, tc storage.ToolCall) error
}

type EventWriter interface {
	Event(name string, data any) error
}

type UserPromptSubmitHook func(ctx context.Context, h *UserPromptSubmitContext) error
type PreToolUseHook func(ctx context.Context, h *ToolUseContext) error
type PostToolUseHook func(ctx context.Context, h *ToolUseContext) error
type StopHook func(ctx context.Context, h *StopContext) error

type HookManager struct {
	UserPromptSubmit []UserPromptSubmitHook
	PreToolUse       []PreToolUseHook
	PostToolUse      []PostToolUseHook
	Stop             []StopHook
}

type LoopState struct {
	Store                        Store
	NewMessageID                 func() string
	ToolRuntimeEnv               func() agenttools.RuntimeEnv
	ShouldInferConversationTitle func(string) bool
	InferConversationTitle       func(string) string

	Conversation storage.Conversation
	User         storage.User
	UserInput    string

	History     []storage.Message
	UserMessage storage.Message

	// ModelHistory is the "model line": the history actually fed to the
	// compression pipeline / model. It starts from the previous turn's
	// compressed request history (loaded from storage) and is appended in
	// lockstep with History during the loop. History stays full and verbatim
	// for display/persistence; ModelHistory may already carry compression
	// artifacts (summary preamble, placeholders, trimmed messages).
	ModelHistory []storage.Message

	SkillSnapshot *agenttools.SkillSnapshot
	SystemPrompt  string
	Messages      []openai.ChatCompletionMessage
	Todos         []agenttools.TodoItem

	ToolCallCount     int // 累计工具调用次数
	LastContextTokens int // 最近一轮发送给模型的上下文 token
	LastContextBudget int // 上下文 token 预算

	Writer EventWriter
}

type UserPromptSubmitContext struct {
	State *LoopState
}

type ToolUseContext struct {
	State    *LoopState
	ToolCall openai.ToolCall
	Name     string
	RawArgs  string
	Outcome  ToolExecutionOutcome
}

type StopContext struct {
	State            *LoopState
	ModelMessage     openai.ChatCompletionMessage
	Content          string
	ReasoningContent string
	AssistantMessage storage.Message
}

type ToolExecutionOutcome struct {
	Status string                `json:"status"`
	Result string                `json:"result"`
	Audit  ToolExecutionAudit    `json:"-"`
	Todos  []agenttools.TodoItem `json:"-"`
}

type ToolExecutionAudit struct {
	ResolvedCWD           string `json:"resolved_cwd,omitempty"`
	ResolvedCommandPath   string `json:"resolved_command_path,omitempty"`
	CommandArtifactPath   string `json:"command_artifact_path,omitempty"`
	CommandArtifactSource string `json:"command_artifact_source,omitempty"`
	OutcomeSummary        string `json:"outcome_summary,omitempty"`
	DenialReason          string `json:"denial_reason,omitempty"`
}

func (s *LoopState) NextMessageID() string {
	if s == nil || s.NewMessageID == nil {
		return ""
	}
	return s.NewMessageID()
}

func (s *LoopState) RuntimeEnv() agenttools.RuntimeEnv {
	if s == nil || s.ToolRuntimeEnv == nil {
		return agenttools.RuntimeEnv{}
	}
	return s.ToolRuntimeEnv()
}

func (s *LoopState) ShouldInferTitle(title string) bool {
	return s != nil && s.ShouldInferConversationTitle != nil && s.ShouldInferConversationTitle(title)
}

func (s *LoopState) InferTitle(input string) string {
	if s == nil || s.InferConversationTitle == nil {
		return input
	}
	return s.InferConversationTitle(input)
}

func (o ToolExecutionOutcome) MessageContent() string {
	data, err := json.Marshal(o)
	if err != nil {
		return fmt.Sprintf(`{"status":%q,"result":%q}`, o.Status, o.Result)
	}
	return string(data)
}

func (o ToolExecutionOutcome) AuditSummary() string {
	data, err := json.Marshal(o.Audit)
	if err != nil {
		return fmt.Sprintf(`{"resolved_cwd":%q,"resolved_command_path":%q,"command_artifact_path":%q,"command_artifact_source":%q,"outcome_summary":%q,"denial_reason":%q}`,
			o.Audit.ResolvedCWD,
			o.Audit.ResolvedCommandPath,
			o.Audit.CommandArtifactPath,
			o.Audit.CommandArtifactSource,
			o.Audit.OutcomeSummary,
			o.Audit.DenialReason,
		)
	}
	return string(data)
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package hooks\n\nimport (\n\t\"context\"\n\t\"encoding/json\"\n\t\"fmt\"\n\n\topenai \"github.com/sashabaranov/go-openai\"\n\n\t\"nano_cc/internal/agent/storage\"\n\tagenttools \"nano_cc/internal/tools\"\n)\n\ntype Store interface {\n\tUpdateConversationTitle(ctx context.Context, conversationID, title string) error\n\tTouchConversationActivity(ctx context.Context, conversationID string) error\n\tSetConversationHistory(ctx context.Context, conversationID string, messages []storage.Message) error\n\tCreateToolCall(ctx context.Context, "}
```

## 2026-06-13T09:17:11Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_03_2EdIasUxV4uDrsJno5Qk1915

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/manager.go"}
```

### Result

```text
package hooks

import "context"

func NewDefaultHookManager() *HookManager {
	return &HookManager{
		UserPromptSubmit: []UserPromptSubmitHook{appendUserMessageHook, conversationActivityHook},
		PreToolUse:       []PreToolUseHook{toolAuditPreHook},
		PostToolUse:      []PostToolUseHook{toolAuditPostHook, persistToolCallHook, appendToolMessageHook},
		Stop:             []StopHook{persistAssistantStopHook, emitAssistantStopHook},
	}
}

func (m *HookManager) RunUserPromptSubmit(ctx context.Context, h *UserPromptSubmitContext) error {
	if m == nil {
		return nil
	}
	for _, hook := range m.UserPromptSubmit {
		if hook == nil {
			continue
		}
		if err := hook(ctx, h); err != nil {
			return err
		}
	}
	return nil
}

func (m *HookManager) RunPreToolUse(ctx context.Context, h *ToolUseContext) error {
	if m == nil {
		return nil
	}
	for _, hook := range m.PreToolUse {
		if hook == nil {
			continue
		}
		if err := hook(ctx, h); err != nil {
			return err
		}
	}
	return nil
}

func (m *HookManager) RunPostToolUse(ctx context.Context, h *ToolUseContext) error {
	if m == nil {
		return nil
	}
	for _, hook := range m.PostToolUse {
		if hook == nil {
			continue
		}
		if err := hook(ctx, h); err != nil {
			return err
		}
	}
	return nil
}

func (m *HookManager) RunStop(ctx context.Context, h *StopContext) error {
	if m == nil {
		return nil
	}
	for _, hook := range m.Stop {
		if hook == nil {
			continue
		}
		if err := hook(ctx, h); err != nil {
			return err
		}
	}
	return nil
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package hooks\n\nimport \"context\"\n\nfunc NewDefaultHookManager() *HookManager {\n\treturn \u0026HookManager{\n\t\tUserPromptSubmit: []UserPromptSubmitHook{appendUserMessageHook, conversationActivityHook},\n\t\tPreToolUse:       []PreToolUseHook{toolAuditPreHook},\n\t\tPostToolUse:      []PostToolUseHook{toolAuditPostHook, persistToolCallHook, appendToolMessageHook},\n\t\tStop:             []StopHook{persistAssistantStopHook, emitAssistantStopHook},\n\t}\n}\n\nfunc (m *HookManager) RunUserPromptSubmit(ctx context.Context, "}
```

## 2026-06-13T09:17:11Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_04_cS3JiTPLy2wwaD3vDkVb4219

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/tool.go"}
```

### Result

```text
package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/idgen"
	"nano_cc/internal/logger"
	"nano_cc/internal/safety"
	"nano_cc/internal/textutil"
)

func toolAuditPreHook(ctx context.Context, h *ToolUseContext) error {
	runtimeEnv := h.State.RuntimeEnv()
	resolvedCommandPath, commandArtifactPath := resolveCommandPaths(h.Name, h.RawArgs, runtimeEnv.CommandBinDir, runtimeEnv.CommandScriptDir)
	h.Outcome.Audit.ResolvedCWD = runtimeEnv.CurrentWorkingDir
	h.Outcome.Audit.ResolvedCommandPath = resolvedCommandPath
	h.Outcome.Audit.CommandArtifactPath = commandArtifactPath
	h.Outcome.Audit.CommandArtifactSource = classifyCommandArtifactSource(runtimeEnv.WorkspaceRoot, commandArtifactPath)
	return nil
}

func toolAuditPostHook(ctx context.Context, h *ToolUseContext) error {
	if h.Outcome.Status == "success" {
		h.Outcome.Audit.OutcomeSummary = textutil.Truncate(h.Outcome.Result, 500)
		return nil
	}
	h.Outcome.Audit.DenialReason = h.Outcome.Result
	return nil
}

func persistToolCallHook(ctx context.Context, h *ToolUseContext) error {
	state := h.State
	_ = state.Store.CreateToolCall(ctx, storage.ToolCall{ID: newToolCallID(), ConversationID: state.Conversation.ID, UserID: state.User.ID, ToolName: h.Name, Status: h.Outcome.Status, Summary: h.Outcome.AuditSummary()})
	return nil
}

func appendToolMessageHook(ctx context.Context, h *ToolUseContext) error {
	content := h.Outcome.MessageContent()
	h.State.Messages = append(h.State.Messages, openai.ChatCompletionMessage{Role: "tool", ToolCallID: h.ToolCall.ID, Content: content})
	toolMessage := storage.Message{ID: h.State.NextMessageID(), ConversationID: h.State.Conversation.ID, UserID: h.State.User.ID, Role: "tool", Content: content, ToolCallID: h.ToolCall.ID}
	h.State.History = append(h.State.History, toolMessage)
	h.State.ModelHistory = append(h.State.ModelHistory, toolMessage)
	appendToolResultLog(ctx, h)
	return nil
}

type toolResultLogStore interface {
	AppendToolResultLog(ctx context.Context, entry storage.ToolResultLogEntry) error
}

func appendToolResultLog(ctx context.Context, h *ToolUseContext) {
	if h == nil || h.State == nil || h.State.Store == nil {
		return
	}
	store, ok := h.State.Store.(toolResultLogStore)
	if !ok {
		return
	}
	entry := storage.ToolResultLogEntry{
		ConversationID: h.State.Conversation.ID,
		SessionID:      h.State.Conversation.SessionID,
		UserID:         h.State.User.ID,
		ToolCallID:     h.ToolCall.ID,
		ToolName:       h.Name,
		RawArgs:        h.RawArgs,
		Status:         h.Outcome.Status,
		Result:         h.Outcome.Result,
		AuditSummary:   h.Outcome.AuditSummary(),
		CreatedAt:      time.Now(),
	}
	if err := store.AppendToolResultLog(ctx, entry); err != nil {
		logger.Warn(fmt.Sprintf("tool result log: append failed conversation=%s tool_call=%s: %v", entry.ConversationID, entry.ToolCallID, err))
	}
}

func newToolCallID() string { return "tc_" + idgen.Hex() }

func classifyCommandArtifactSource(workspaceRoot, commandArtifactPath string) string {
	if strings.TrimSpace(commandArtifactPath) == "" {
		return ""
	}
	cleanArtifact := filepath.Clean(commandArtifactPath)
	cleanWorkspace := strings.TrimSpace(workspaceRoot)
	if cleanWorkspace != "" {
		cleanWorkspace = filepath.Clean(cleanWorkspace)
		if safety.Contains(cleanWorkspace, cleanArtifact) {
			return "workspace"
		}
	}
	return "custom"
}

func resolveCommandPaths(toolName, rawArgs string, roots ...string) (string, string) {
	if toolName != "bash" {
		return "", ""
	}
	var args map[string]any
	if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
		return "", ""
	}
	command, _ := args["command"].(string)
	for _, token := range strings.Fields(command) {
		candidate := strings.Trim(token, "\"'`;,()[]{}")
		if candidate == "" || !filepath.IsAbs(candidate) {
			continue
		}
		resolved, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}
		cleanResolved := filepath.Clean(resolved)
		for _, root := range roots {
			if root == "" {
				continue
			}
			resolvedRoot, err := filepath.Abs(root)
			if err != nil {
				continue
			}
			cleanRoot := filepath.Clean(resolvedRoot)
			if safety.Contains(cleanRoot, cleanResolved) {
				return cleanResolved, cleanResolved
			}
		}
		return cleanResolved, ""
	}
	return "", ""
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package hooks\n\nimport (\n\t\"context\"\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"path/filepath\"\n\t\"strings\"\n\t\"time\"\n\n\topenai \"github.com/sashabaranov/go-openai\"\n\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/idgen\"\n\t\"nano_cc/internal/logger\"\n\t\"nano_cc/internal/safety\"\n\t\"nano_cc/internal/textutil\"\n)\n\nfunc toolAuditPreHook(ctx context.Context, h *ToolUseContext) error {\n\truntimeEnv := h.State.RuntimeEnv()\n\tresolvedCommandPath, commandArtifactPath := resolveCommandPaths(h.Name, h.RawArgs, runtimeEnv.CommandBinDi"}
```

## 2026-06-13T09:17:11Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_05_OlBayywn56Q8iXulBx1Y7965

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/user_prompt.go"}
```

### Result

```text
package hooks

import (
	"context"

	"nano_cc/internal/agent/storage"
)

func appendUserMessageHook(ctx context.Context, h *UserPromptSubmitContext) error {
	state := h.State
	state.UserMessage = storage.Message{ID: state.NextMessageID(), ConversationID: state.Conversation.ID, UserID: state.User.ID, Role: "user", Content: state.UserInput}
	state.History = append(state.History, state.UserMessage)
	state.ModelHistory = append(state.ModelHistory, state.UserMessage)
	return nil
}

func conversationActivityHook(ctx context.Context, h *UserPromptSubmitContext) error {
	state := h.State
	if state.ShouldInferTitle(state.Conversation.Title) {
		return state.Store.UpdateConversationTitle(ctx, state.Conversation.ID, state.InferTitle(state.UserInput))
	}
	return state.Store.TouchConversationActivity(ctx, state.Conversation.ID)
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package hooks\n\nimport (\n\t\"context\"\n\n\t\"nano_cc/internal/agent/storage\"\n)\n\nfunc appendUserMessageHook(ctx context.Context, h *UserPromptSubmitContext) error {\n\tstate := h.State\n\tstate.UserMessage = storage.Message{ID: state.NextMessageID(), ConversationID: state.Conversation.ID, UserID: state.User.ID, Role: \"user\", Content: state.UserInput}\n\tstate.History = append(state.History, state.UserMessage)\n\tstate.ModelHistory = append(state.ModelHistory, state.UserMessage)\n\treturn nil\n}\n\nfunc conversationA"}
```

## 2026-06-13T09:17:11Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_06_fdYop3iZZ5e8d77qrCZK4620

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/stop.go"}
```

### Result

```text
package hooks

import (
	"context"

	"nano_cc/internal/agent/storage"
)

func persistAssistantStopHook(ctx context.Context, h *StopContext) error {
	state := h.State
	h.AssistantMessage = storage.Message{ID: state.NextMessageID(), ConversationID: state.Conversation.ID, UserID: state.User.ID, Role: "assistant", Content: h.Content, ReasoningContent: h.ReasoningContent, Meta: stopMessageMeta(state)}
	updatedHistory := append(state.History, h.AssistantMessage)
	return state.Store.SetConversationHistory(ctx, state.Conversation.ID, updatedHistory)
}

func emitAssistantStopHook(ctx context.Context, h *StopContext) error {
	if h.State.Writer == nil {
		return nil
	}
	meta := stopMessageMeta(h.State)
	_ = h.State.Writer.Event("assistant", map[string]any{
		"message_id":        h.AssistantMessage.ID,
		"content":           h.AssistantMessage.Content,
		"reasoning_content": h.AssistantMessage.ReasoningContent,
		"final":             true,
		"tool_call_count":   meta.ToolCallCount,
		"context_tokens":    meta.ContextTokens,
		"context_budget":    meta.ContextBudget,
	})
	return nil
}

func stopMessageMeta(state *LoopState) *storage.MessageMeta {
	return &storage.MessageMeta{
		ToolCallCount: state.ToolCallCount,
		ContextTokens: state.LastContextTokens,
		ContextBudget: state.LastContextBudget,
	}
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package hooks\n\nimport (\n\t\"context\"\n\n\t\"nano_cc/internal/agent/storage\"\n)\n\nfunc persistAssistantStopHook(ctx context.Context, h *StopContext) error {\n\tstate := h.State\n\th.AssistantMessage = storage.Message{ID: state.NextMessageID(), ConversationID: state.Conversation.ID, UserID: state.User.ID, Role: \"assistant\", Content: h.Content, ReasoningContent: h.ReasoningContent, Meta: stopMessageMeta(state)}\n\tupdatedHistory := append(state.History, h.AssistantMessage)\n\treturn state.Store.SetConversationHist"}
```

## 2026-06-13T09:17:14Z · bash · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_0DABXNugBatyaE5xNPy17851

### Arguments

```json
{"command": "find /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression -type f -name \"*.go\" | sort"}
```

### Result

```text
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/compression_test.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/conversation_memory_strategy.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/full_history_summarization.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/message_window_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/reactive_compact.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/recent_tool_result_retention.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/token_estimator.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/tool_result_compression.go
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","resolved_command_path":"/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression","outcome_summary":"/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/compression_test.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/compression.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/conversation_memory_strategy.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/full_history_summarization.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/mes"}
```

## 2026-06-13T09:17:14Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_3TeKinB8MfLq3DwVnvt93840

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/strategies.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/strategies.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/strategies.go: no such file or directory"}
```

## 2026-06-13T09:17:14Z · bash · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_eRZizYZ2yt3zkGQx41xD3090

### Arguments

```json
{"command": "find /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/strategies -type f -name \"*.go\" 2>/dev/null | sort"}
```

### Result

```text
(no output)
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","resolved_command_path":"/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/strategies","outcome_summary":"(no output)"}
```

## 2026-06-13T09:17:18Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_kbciG6eubD9FZKGZ6LpB2643

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/tool_result_compression.go"}
```

### Result

```text
package compression

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/idgen"
	"nano_cc/internal/textutil"
)

const toolResultCompressionStrategyName = "tool_result_compression"

// ToolResultCompressionStrategy persists oversized tool_result outputs from the
// latest user turn and replaces their inline content with a <persisted-output>
// marker plus a short preview.
type ToolResultCompressionStrategy struct{}

func (s *ToolResultCompressionStrategy) Name() string { return toolResultCompressionStrategyName }

// toolResultMessageContent is the JSON wrapper used by tool messages.
type toolResultMessageContent struct {
	Status string `json:"status"`
	Result string `json:"result"`
}

// rebuildToolResult reassembles a tool message content, preserving the JSON
// wrapper and status when the original was JSON.
func rebuildToolResult(status, result string, isJSON bool) string {
	if !isJSON {
		return result
	}
	data, err := json.Marshal(toolResultMessageContent{Status: status, Result: result})
	if err != nil {
		return result
	}
	return string(data)
}

func isCompactedResult(result string) bool {
	return result == earlierToolResultPlaceholder || strings.Contains(result, PersistedOutputMarkerPrefix)
}

// latestUserTurnToolIndexes returns the indexes of tool messages that appear
// after the last user message.
func latestUserTurnToolIndexes(history []storage.Message) []int {
	lastUser := -1
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Role == "user" {
			lastUser = i
			break
		}
	}
	var indexes []int
	for i := lastUser + 1; i < len(history); i++ {
		if history[i].Role == "tool" {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

type toolResultCandidate struct {
	index  int
	status string
	result string
	isJSON bool
	bytes  int
}

func (s *ToolResultCompressionStrategy) Apply(ctx context.Context, req *Request) error {
	history := req.RequestHistory
	indexes := latestUserTurnToolIndexes(history)
	if len(indexes) == 0 {
		return nil
	}

	candidates := make([]toolResultCandidate, 0, len(indexes))
	totalBytes := 0
	for _, idx := range indexes {
		status, result, isJSON := textutil.ParseToolResult(history[idx].Content)
		if isCompactedResult(result) {
			continue
		}
		size := len([]byte(result))
		candidates = append(candidates, toolResultCandidate{index: idx, status: status, result: result, isJSON: isJSON, bytes: size})
		totalBytes += size
	}
	if totalBytes <= toolResultByteThreshold {
		return nil
	}

	// Largest first.
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].bytes > candidates[j].bytes
	})

	remaining := totalBytes
	for _, candidate := range candidates {
		if remaining <= toolResultByteThreshold {
			break
		}
		marker, err := s.persistAndBuildMarker(ctx, req, history[candidate.index], candidate)
		if err != nil {
			return err
		}
		history[candidate.index].Content = rebuildToolResult(candidate.status, marker, candidate.isJSON)
		remaining -= candidate.bytes
	}
	return nil
}

func (s *ToolResultCompressionStrategy) persistAndBuildMarker(ctx context.Context, req *Request, msg storage.Message, candidate toolResultCandidate) (string, error) {
	sum := sha256.Sum256([]byte(candidate.result))
	contentSHA := hex.EncodeToString(sum[:])

	existing, err := req.Store.GetPersistedOutputByMessageHash(ctx, req.Conversation.ID, req.User.ID, msg.ID, msg.ToolCallID, toolResultCompressionStrategyName, contentSHA)
	var id string
	if err == nil && existing.ID != "" {
		id = existing.ID
	} else {
		id = "po_" + idgen.Hex()
		preview := previewRunes(candidate.result, toolResultPreviewRunes)
		record := storage.PersistedOutput{
			ID:             id,
			ConversationID: req.Conversation.ID,
			UserID:         req.User.ID,
			MessageID:      msg.ID,
			ToolCallID:     msg.ToolCallID,
			Kind:           "tool_result",
			Strategy:       toolResultCompressionStrategyName,
			OriginalBytes:  candidate.bytes,
			ContentSHA256:  contentSHA,
			Content:        candidate.result,
			Preview:        preview,
		}
		if err := req.Store.CreatePersistedOutput(ctx, record); err != nil {
			return "", fmt.Errorf("persist tool output: %w", err)
		}
	}

	preview := previewRunes(candidate.result, toolResultPreviewRunes)
	return buildPersistedOutputMarker(id, candidate.bytes, preview), nil
}

func buildPersistedOutputMarker(id string, originalBytes int, preview string) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<persisted-output id=%q kind="tool_result" original_bytes="%d" preview_chars="%d" retrieval_tool="read_persisted_output">`, id, originalBytes, toolResultPreviewRunes)
	b.WriteString("\n")
	fmt.Fprintf(&b, "完整输出已持久化；如需更多内容，请调用 read_persisted_output(id=%q, offset=0, limit=20000) 分段读取。\n\n", id)
	b.WriteString(preview)
	b.WriteString("\n</persisted-output>")
	return b.String()
}

func previewRunes(text string, max int) string {
	runes := []rune(text)
	if len(runes) <= max {
		return text
	}
	return string(runes[:max])
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package compression\n\nimport (\n\t\"context\"\n\t\"crypto/sha256\"\n\t\"encoding/hex\"\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"sort\"\n\t\"strings\"\n\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/idgen\"\n\t\"nano_cc/internal/textutil\"\n)\n\nconst toolResultCompressionStrategyName = \"tool_result_compression\"\n\n// ToolResultCompressionStrategy persists oversized tool_result outputs from the\n// latest user turn and replaces their inline content with a \u003cpersisted-output\u003e\n// marker plus a short preview.\ntype ToolResultCompressionSt"}
```

## 2026-06-13T09:17:18Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_ijLHKhE2S1omBDcLipwE7280

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/message_window_compression.go"}
```

### Result

```text
package compression

import (
	"context"
	"strings"

	"nano_cc/internal/agent/storage"
)

const messageWindowCompressionStrategyName = "message_window_compression"

// MessageWindowCompressionStrategy trims middle history when the request copy
// exceeds the window limit, keeping the head and tail, then repairs any
// dangling OpenAI tool_call / tool_result pairs created by the cut.
type MessageWindowCompressionStrategy struct{}

func (s *MessageWindowCompressionStrategy) Name() string {
	return messageWindowCompressionStrategyName
}

func (s *MessageWindowCompressionStrategy) Apply(ctx context.Context, req *Request) error {
	history := req.RequestHistory
	if len(history) <= messageWindowLimit {
		return nil
	}
	head := history[:messageWindowHead]
	tail := history[len(history)-messageWindowTail:]
	windowed := make([]storage.Message, 0, messageWindowHead+messageWindowTail)
	windowed = append(windowed, head...)
	windowed = append(windowed, tail...)
	req.RequestHistory = repairToolCallBoundaries(windowed)
	return nil
}

// repairToolCallBoundaries removes orphan tool messages (no preceding assistant
// tool_call) and clears assistant tool_calls that lost their tool results.
func repairToolCallBoundaries(history []storage.Message) []storage.Message {
	// Collect tool_call ids that the assistant messages still expose.
	assistantCallIDs := make(map[string]struct{})
	for _, msg := range history {
		if msg.Role == "assistant" {
			for _, call := range msg.ToolCalls {
				assistantCallIDs[call.ID] = struct{}{}
			}
		}
	}
	// Collect tool result ids present in the window.
	toolResultIDs := make(map[string]struct{})
	for _, msg := range history {
		if msg.Role == "tool" && msg.ToolCallID != "" {
			toolResultIDs[msg.ToolCallID] = struct{}{}
		}
	}

	repaired := make([]storage.Message, 0, len(history))
	for _, msg := range history {
		switch msg.Role {
		case "tool":
			// Drop orphan tool messages without a matching assistant call.
			if _, ok := assistantCallIDs[msg.ToolCallID]; !ok {
				continue
			}
			repaired = append(repaired, msg)
		case "assistant":
			if len(msg.ToolCalls) > 0 {
				kept := msg.ToolCalls[:0:0]
				for _, call := range msg.ToolCalls {
					if _, ok := toolResultIDs[call.ID]; ok {
						kept = append(kept, call)
					}
				}
				msg.ToolCalls = kept
				if len(kept) == 0 && strings.TrimSpace(msg.Content) == "" && strings.TrimSpace(msg.ReasoningContent) == "" {
					continue
				}
			}
			repaired = append(repaired, msg)
		default:
			repaired = append(repaired, msg)
		}
	}
	return repaired
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package compression\n\nimport (\n\t\"context\"\n\t\"strings\"\n\n\t\"nano_cc/internal/agent/storage\"\n)\n\nconst messageWindowCompressionStrategyName = \"message_window_compression\"\n\n// MessageWindowCompressionStrategy trims middle history when the request copy\n// exceeds the window limit, keeping the head and tail, then repairs any\n// dangling OpenAI tool_call / tool_result pairs created by the cut.\ntype MessageWindowCompressionStrategy struct{}\n\nfunc (s *MessageWindowCompressionStrategy) Name() string {\n\treturn"}
```

## 2026-06-13T09:17:18Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_9uAtQ0IH7iUZ6euhH1rg3877

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/conversation_memory_strategy.go"}
```

### Result

```text
package compression

import (
	"context"
	"strings"

	"nano_cc/internal/agent/storage"
)

const conversationMemoryStrategyName = "conversation_memory"

const (
	conversationMemoryMinTokens = 10 * 1024
	conversationMemoryMaxTokens = 40 * 1024
	conversationMemoryMinCount  = 5

	conversationMemorySystemPreamble = "以下是本会话的记忆，由系统随对话持续维护，仅用于本次模型推理，不是用户发送的真实消息。请把它当作已发生对话的可靠浓缩。"
)

// ConversationMemoryStrategy runs just before FullHistorySummarizationStrategy.
// When the request still exceeds the token budget but a sufficiently rich,
// already-maintained conversation memory exists, it rebuilds RequestHistory from
// that memory (plus a recent tail and the always-kept last message), so the
// fallback full-history summarization is skipped entirely. Otherwise it is a
// no-op and the fallback summarizer handles compression.
type ConversationMemoryStrategy struct{}

func (s *ConversationMemoryStrategy) Name() string {
	return conversationMemoryStrategyName
}

func (s *ConversationMemoryStrategy) Apply(ctx context.Context, req *Request) error {
	if req.Estimator == nil || req.Store == nil {
		return nil
	}
	budget := req.Estimator.ContextTokenBudget()
	before := req.Estimator.EstimateRequestTokens(req.SystemPrompt, req.RequestHistory, req.Tools)
	if before <= budget {
		return nil
	}
	if len(req.RequestHistory) < 2 {
		return nil
	}

	items, err := req.Store.ListConversationMemories(ctx, req.Conversation.ID)
	if err != nil || len(items) < conversationMemoryMinCount {
		return nil
	}
	memoryText := renderConversationMemory(items)
	memTokens := estimateTokensFromBytes(len([]byte(memoryText)))
	if memTokens < conversationMemoryMinTokens || memTokens > conversationMemoryMaxTokens {
		return nil
	}

	lastMessage := req.RequestHistory[len(req.RequestHistory)-1]
	summarizable := req.RequestHistory[:len(req.RequestHistory)-1]
	tail := selectRecentTail(req, summarizable, memoryText, budget)
	kept := repairToolCallBoundaries(append(append([]storage.Message{}, tail...), lastMessage))
	req.RequestHistory = buildConversationMemoryHistory(memoryText, kept)
	return nil
}

func renderConversationMemory(items []storage.ConversationMemory) string {
	var b strings.Builder
	for _, m := range items {
		b.WriteString("- ")
		b.WriteString(m.Name)
		if strings.TrimSpace(m.Description) != "" {
			b.WriteString("：")
			b.WriteString(m.Description)
		}
		b.WriteByte('\n')
		if strings.TrimSpace(m.Body) != "" {
			b.WriteString("  ")
			b.WriteString(m.Body)
			b.WriteByte('\n')
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

func buildConversationMemoryHistory(memoryText string, tail []storage.Message) []storage.Message {
	messages := []storage.Message{
		{Role: "system", Content: conversationMemorySystemPreamble},
		{Role: "user", Content: "<conversation-memory>\n" + memoryText + "\n</conversation-memory>"},
	}
	return append(messages, tail...)
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package compression\n\nimport (\n\t\"context\"\n\t\"strings\"\n\n\t\"nano_cc/internal/agent/storage\"\n)\n\nconst conversationMemoryStrategyName = \"conversation_memory\"\n\nconst (\n\tconversationMemoryMinTokens = 10 * 1024\n\tconversationMemoryMaxTokens = 40 * 1024\n\tconversationMemoryMinCount  = 5\n\n\tconversationMemorySystemPreamble = \"以下是本会话的记忆，由系统随对话持续维护，仅用于本次模型推理，不是用户发送的真实消息。请把它当作已发生对话的可靠浓缩。\"\n)\n\n// Co"}
```

## 2026-06-13T09:17:18Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_03_XQpsTOs9C5sGV8XfppnQ2878

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/full_history_summarization.go"}
```

### Result

```text
package compression

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/idgen"
)

const fullHistorySummarizationStrategyName = "full_history_summarization"

const (
	summaryTargetTokens   = 8 * 1024
	recentTailMinTokens   = 16 * 1024
	summarySystemPreamble = "以下是为满足上下文窗口限制生成的会话摘要，仅用于本次模型推理，不是用户发送的真实消息。请把它当作已发生对话的可靠浓缩。"
)

// HistorySummarizer summarizes a request-state history into a markdown summary.
type HistorySummarizer func(ctx context.Context, req SummaryRequest) (SummaryResult, error)

type SummaryRequest struct {
	Conversation storage.Conversation
	User         storage.User
	History      []storage.Message
	TargetTokens int
}

type SummaryResult struct {
	Summary string
}

// FullHistorySummarizationStrategy is the fallback layer: when the deterministic
// three layers still exceed the token budget, it replaces RequestHistory with a
// summary plus a recent tail window. Request-only; never written to display history.
type FullHistorySummarizationStrategy struct{}

func (s *FullHistorySummarizationStrategy) Name() string {
	return fullHistorySummarizationStrategyName
}

func (s *FullHistorySummarizationStrategy) Apply(ctx context.Context, req *Request) error {
	if req.Estimator == nil || req.Summarizer == nil {
		return nil
	}
	budget := req.Estimator.ContextTokenBudget()
	before := req.Estimator.EstimateRequestTokens(req.SystemPrompt, req.RequestHistory, req.Tools)
	if before <= budget {
		return nil
	}
	// The last message is always preserved verbatim; only the earlier history
	// is summarized. With a single message there is nothing to summarize.
	if len(req.RequestHistory) < 2 {
		return nil
	}

	lastMessage := req.RequestHistory[len(req.RequestHistory)-1]
	summarizable := req.RequestHistory[:len(req.RequestHistory)-1]
	sourceHash := historyHash(summarizable)

	// Try cache first.
	summary := ""
	if cached, err := req.Store.GetContextSummaryByHistoryHash(ctx, req.Conversation.ID, req.User.ID, sourceHash); err == nil && cached.Summary != "" {
		summary = cached.Summary
	} else {
		result, err := req.Summarizer(ctx, SummaryRequest{
			Conversation: req.Conversation,
			User:         req.User,
			History:      summarizable,
			TargetTokens: summaryTargetTokens,
		})
		if err != nil {
			return fmt.Errorf("summarize history: %w", err)
		}
		summary = result.Summary
		after := req.Estimator.EstimateRequestTokens(req.SystemPrompt, buildSummaryHistory(summary, nil), req.Tools)
		_ = req.Store.CreateContextSummary(ctx, storage.ContextSummary{
			ID:                    "cs_" + idgen.Hex(),
			ConversationID:        req.Conversation.ID,
			UserID:                req.User.ID,
			SourceHistorySHA256:   sourceHash,
			Strategy:              fullHistorySummarizationStrategyName,
			EstimatedTokensBefore: before,
			EstimatedTokensAfter:  after,
			Summary:               summary,
		})
	}

	tail := selectRecentTail(req, summarizable, summary, budget)
	// Tail plus the always-kept last message are repaired together so the last
	// message is never dropped as an orphan tool result.
	kept := repairToolCallBoundaries(append(append([]storage.Message{}, tail...), lastMessage))
	req.RequestHistory = buildSummaryHistory(summary, kept)
	return nil
}

// selectRecentTail picks messages from the end of the summarizable history
// within the budget remaining after reserving space for the summary. The
// always-kept last message is handled separately by the caller.
func selectRecentTail(req *Request, summarizable []storage.Message, summary string, budget int) []storage.Message {
	summaryTokens := req.Estimator.EstimateRequestTokens(req.SystemPrompt, buildSummaryHistory(summary, nil), req.Tools)
	tailBudget := budget - summaryTokens
	if tailBudget > recentTailMinTokens {
		// Cap reserved tail space to keep summary authoritative.
		tailBudget = recentTailMinTokens
	}
	if tailBudget <= 0 {
		return nil
	}

	var selected []storage.Message
	used := 0
	for i := len(summarizable) - 1; i >= 0; i-- {
		msgTokens := estimateTokensFromBytes(messageBytes(summarizable[i]))
		if used+msgTokens > tailBudget {
			break
		}
		used += msgTokens
		selected = append([]storage.Message{summarizable[i]}, selected...)
	}
	return selected
}

func buildSummaryHistory(summary string, tail []storage.Message) []storage.Message {
	messages := []storage.Message{
		{Role: "system", Content: summarySystemPreamble},
		{Role: "user", Content: "<conversation-summary>\n" + summary + "\n</conversation-summary>"},
	}
	return append(messages, tail...)
}

func messageBytes(msg storage.Message) int {
	if data, err := json.Marshal(msg); err == nil {
		return len(data)
	}
	return len(msg.Content) + len(msg.ReasoningContent)
}

func historyHash(history []storage.Message) string {
	data, err := json.Marshal(history)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package compression\n\nimport (\n\t\"context\"\n\t\"crypto/sha256\"\n\t\"encoding/hex\"\n\t\"encoding/json\"\n\t\"fmt\"\n\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/idgen\"\n)\n\nconst fullHistorySummarizationStrategyName = \"full_history_summarization\"\n\nconst (\n\tsummaryTargetTokens   = 8 * 1024\n\trecentTailMinTokens   = 16 * 1024\n\tsummarySystemPreamble = \"以下是为满足上下文窗口限制生成的会话摘要，仅用于本次模型推理，不是用户发送的真实消息。请把它当作已发生对话\ufffd"}
```

## 2026-06-13T09:17:18Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_04_RfRa8L4pLfi6YWXV5kJZ3405

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/recent_tool_result_retention.go"}
```

### Result

```text
package compression

import (
	"context"

	"nano_cc/internal/textutil"
)

const recentToolResultRetentionStrategyName = "recent_tool_result_retention"

// RecentToolResultRetentionStrategy keeps only the most recent N tool results
// fully inline; earlier tool results are replaced with a one-line placeholder.
type RecentToolResultRetentionStrategy struct{}

func (s *RecentToolResultRetentionStrategy) Name() string {
	return recentToolResultRetentionStrategyName
}

func (s *RecentToolResultRetentionStrategy) Apply(ctx context.Context, req *Request) error {
	history := req.RequestHistory

	// Collect tool message indexes in order.
	var toolIndexes []int
	for i := range history {
		if history[i].Role == "tool" {
			toolIndexes = append(toolIndexes, i)
		}
	}
	if len(toolIndexes) <= recentToolResultRetention {
		return nil
	}

	// The last N tool messages (by position) keep their full result.
	cutoff := len(toolIndexes) - recentToolResultRetention
	for _, idx := range toolIndexes[:cutoff] {
		status, result, isJSON := textutil.ParseToolResult(history[idx].Content)
		if isCompactedResult(result) {
			continue
		}
		history[idx].Content = rebuildToolResult(status, earlierToolResultPlaceholder, isJSON)
	}
	return nil
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package compression\n\nimport (\n\t\"context\"\n\n\t\"nano_cc/internal/textutil\"\n)\n\nconst recentToolResultRetentionStrategyName = \"recent_tool_result_retention\"\n\n// RecentToolResultRetentionStrategy keeps only the most recent N tool results\n// fully inline; earlier tool results are replaced with a one-line placeholder.\ntype RecentToolResultRetentionStrategy struct{}\n\nfunc (s *RecentToolResultRetentionStrategy) Name() string {\n\treturn recentToolResultRetentionStrategyName\n}\n\nfunc (s *RecentToolResultRetent"}
```

## 2026-06-13T09:17:18Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_05_HLHQoxf7lgNJBsEpZWTc5758

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/reactive_compact.go"}
```

### Result

```text
package compression

import (
	"context"
	"fmt"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/idgen"
)

const reactiveCompactStrategyName = "reactive_compact"

const (
	// reactiveSummaryTargetTokens is the summary budget for reactive compaction;
	// more aggressive than FullHistorySummarization's 8K.
	reactiveSummaryTargetTokens = 2 * 1024
	// reactiveTailMinTokens caps how much recent history is kept verbatim;
	// more aggressive than FullHistorySummarization's 16K.
	reactiveTailMinTokens = 4 * 1024
)

// ReactiveCompactStrategy is invoked out-of-band when the LLM rejects a request
// with HTTP 413 (context overflow). It is structurally similar to
// FullHistorySummarizationStrategy but more aggressive: it compacts
// unconditionally (no budget threshold check), targets a smaller summary, keeps
// a shorter recent tail, and drops reasoning_content from kept messages.
type ReactiveCompactStrategy struct{}

func (s *ReactiveCompactStrategy) Name() string {
	return reactiveCompactStrategyName
}

func (s *ReactiveCompactStrategy) Apply(ctx context.Context, req *Request) error {
	if req.Estimator == nil || req.Summarizer == nil {
		return nil
	}
	// The last message is always preserved verbatim; only the earlier history
	// is summarized. With a single message there is nothing to summarize.
	if len(req.RequestHistory) < 2 {
		return nil
	}

	before := req.Estimator.EstimateRequestTokens(req.SystemPrompt, req.RequestHistory, req.Tools)
	lastMessage := req.RequestHistory[len(req.RequestHistory)-1]
	summarizable := req.RequestHistory[:len(req.RequestHistory)-1]
	sourceHash := historyHash(summarizable)

	summary := ""
	if cached, err := req.Store.GetContextSummaryByHistoryHash(ctx, req.Conversation.ID, req.User.ID, sourceHash); err == nil && cached.Summary != "" {
		summary = cached.Summary
	} else {
		result, err := req.Summarizer(ctx, SummaryRequest{
			Conversation: req.Conversation,
			User:         req.User,
			History:      summarizable,
			TargetTokens: reactiveSummaryTargetTokens,
		})
		if err != nil {
			return fmt.Errorf("reactive compact summarize history: %w", err)
		}
		summary = result.Summary
		after := req.Estimator.EstimateRequestTokens(req.SystemPrompt, buildSummaryHistory(summary, nil), req.Tools)
		_ = req.Store.CreateContextSummary(ctx, storage.ContextSummary{
			ID:                    "cs_" + idgen.Hex(),
			ConversationID:        req.Conversation.ID,
			UserID:                req.User.ID,
			SourceHistorySHA256:   sourceHash,
			Strategy:              reactiveCompactStrategyName,
			EstimatedTokensBefore: before,
			EstimatedTokensAfter:  after,
			Summary:               summary,
		})
	}

	tail := selectReactiveTail(req, summarizable, summary)
	kept := repairToolCallBoundaries(append(append([]storage.Message{}, tail...), lastMessage))
	dropReasoningContent(kept)
	req.RequestHistory = buildSummaryHistory(summary, kept)
	return nil
}

// selectReactiveTail picks messages from the end of the summarizable history
// within an aggressive tail budget reserved after the summary.
func selectReactiveTail(req *Request, summarizable []storage.Message, summary string) []storage.Message {
	summaryTokens := req.Estimator.EstimateRequestTokens(req.SystemPrompt, buildSummaryHistory(summary, nil), req.Tools)
	tailBudget := req.Estimator.ContextTokenBudget() - summaryTokens
	if tailBudget > reactiveTailMinTokens {
		tailBudget = reactiveTailMinTokens
	}
	if tailBudget <= 0 {
		return nil
	}

	var selected []storage.Message
	used := 0
	for i := len(summarizable) - 1; i >= 0; i-- {
		msgTokens := estimateTokensFromBytes(messageBytes(summarizable[i]))
		if used+msgTokens > tailBudget {
			break
		}
		used += msgTokens
		selected = append([]storage.Message{summarizable[i]}, selected...)
	}
	return selected
}

// dropReasoningContent strips reasoning_content from kept messages in place to
// reclaim context space; it is never needed for continuation.
func dropReasoningContent(messages []storage.Message) {
	for i := range messages {
		messages[i].ReasoningContent = ""
	}
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package compression\n\nimport (\n\t\"context\"\n\t\"fmt\"\n\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/idgen\"\n)\n\nconst reactiveCompactStrategyName = \"reactive_compact\"\n\nconst (\n\t// reactiveSummaryTargetTokens is the summary budget for reactive compaction;\n\t// more aggressive than FullHistorySummarization's 8K.\n\treactiveSummaryTargetTokens = 2 * 1024\n\t// reactiveTailMinTokens caps how much recent history is kept verbatim;\n\t// more aggressive than FullHistorySummarization's 16K.\n\treactiveTailMinTok"}
```

## 2026-06-13T09:17:18Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_06_PRK1TNT8aXFSOkPvS37C2778

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/token_estimator.go"}
```

### Result

```text
package compression

import (
	"encoding/json"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
)

const (
	// 128KB 上下文限制
	defaultModelContextLimit = 128 * 1024
	// 8KB 最大响应 token 数
	defaultMaxResponseTokens = 8 * 1024
	// 4KB 安全余量
	defaultSafetyMargin = 4 * 1024
)

// TokenEstimator estimates the token footprint of an outgoing request.
type TokenEstimator interface {
	EstimateRequestTokens(systemPrompt string, history []storage.Message, tools []openai.Tool) int
	ContextTokenBudget() int
}

// DefaultTokenEstimator uses a conservative ceil(utf8Bytes/3) approximation.
type DefaultTokenEstimator struct{}

func (DefaultTokenEstimator) ContextTokenBudget() int {
	return defaultModelContextLimit - defaultMaxResponseTokens - defaultSafetyMargin
}

func (e DefaultTokenEstimator) EstimateRequestTokens(systemPrompt string, history []storage.Message, tools []openai.Tool) int {
	bytes := len([]byte(systemPrompt))
	for _, msg := range history {
		if data, err := json.Marshal(msg); err == nil {
			bytes += len(data)
		} else {
			bytes += len(msg.Content) + len(msg.ReasoningContent)
		}
	}
	if len(tools) > 0 {
		if data, err := json.Marshal(tools); err == nil {
			bytes += len(data)
		}
	}
	return estimateTokensFromBytes(bytes)
}

func estimateTokensFromBytes(bytes int) int {
	if bytes <= 0 {
		return 0
	}
	return (bytes + 2) / 3
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package compression\n\nimport (\n\t\"encoding/json\"\n\n\topenai \"github.com/sashabaranov/go-openai\"\n\n\t\"nano_cc/internal/agent/storage\"\n)\n\nconst (\n\t// 128KB 上下文限制\n\tdefaultModelContextLimit = 128 * 1024\n\t// 8KB 最大响应 token 数\n\tdefaultMaxResponseTokens = 8 * 1024\n\t// 4KB 安全余量\n\tdefaultSafetyMargin = 4 * 1024\n)\n\n// TokenEstimator estimates the token footprint of an outgoing request.\ntype TokenEstimator interface {\n\tEstimateRequestTokens(systemPrompt string, history []storage.Message"}
```

## 2026-06-13T09:17:22Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_yyWXkAQPppVYTlKI6wRr4879

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/config.go"}
```

### Result

```text
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"-"`
	ModelID string `json:"model_id"`
}

type fileConfig struct {
	Config
	AppHome                            string `json:"app_home"`
	BuiltinSkillsDir                   string `json:"builtin_skills_dir"`
	CommandBinDir                      string `json:"command_bin_dir"`
	CommandScriptDir                   string `json:"command_script_dir"`
	SystemPromptPath                   string `json:"system_prompt_path"`
	WorkspaceRoot                      string `json:"workspace_root"`
	AllowedTools                       string `json:"allowed_tools"`
	BashAllowOutsideWorkspace          bool   `json:"bash_allow_outside_workspace"`
	BashAllowDangerousCommands         bool   `json:"bash_allow_dangerous_commands"`
	ConversationLockTTLSeconds         int    `json:"conversation_lock_ttl_seconds"`
	MemoryWorkTimeoutSeconds           int    `json:"memory_work_timeout_seconds"`
	ConversationLockWaitTimeoutSeconds int    `json:"conversation_lock_wait_timeout_seconds"`
}

type AppConfig struct {
	LLM                         Config
	AppHome                     string
	BuiltinSkillsDir            string
	CommandBinDir               string
	CommandScriptDir            string
	SystemPromptPath            string
	WorkspaceRoot               string
	LogsDir                     string
	AllowedTools                []string
	BashAllowOutsideWorkspace   bool
	BashAllowDangerousCommands  bool
	ConversationLockTTL         time.Duration
	MemoryWorkTimeout           time.Duration
	ConversationLockWaitTimeout time.Duration
}

func loadConfigFile() (fileConfig, error) {
	data, err := os.ReadFile(configFilePath())
	if err != nil {
		return fileConfig{}, err
	}

	var cfg fileConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fileConfig{}, fmt.Errorf("failed to parse workspace config.json: %w", err)
	}
	return cfg, nil
}

func configFilePath() string {
	return "config.json"
}

func loadLLMConfig(fileCfg fileConfig) (Config, error) {
	cfg := Config{
		BaseURL: strings.TrimSpace(fileCfg.Config.BaseURL),
		APIKey:  strings.TrimSpace(getenv("OPENAI_API_KEY")),
		ModelID: strings.TrimSpace(fileCfg.Config.ModelID),
	}

	if cfg.BaseURL == "" || cfg.ModelID == "" {
		return Config{}, fmt.Errorf("missing LLM config; set base_url, model_id in config.json")
	}
	if cfg.APIKey == "" {
		return Config{}, fmt.Errorf("missing LLM api key; set OPENAI_API_KEY environment variable")
	}

	return cfg, nil
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package config\n\nimport (\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"os\"\n\t\"strings\"\n\t\"time\"\n)\n\ntype Config struct {\n\tBaseURL string `json:\"base_url\"`\n\tAPIKey  string `json:\"-\"`\n\tModelID string `json:\"model_id\"`\n}\n\ntype fileConfig struct {\n\tConfig\n\tAppHome                            string `json:\"app_home\"`\n\tBuiltinSkillsDir                   string `json:\"builtin_skills_dir\"`\n\tCommandBinDir                      string `json:\"command_bin_dir\"`\n\tCommandScriptDir                   string `json:\"command_script_dir\"`\n"}
```

## 2026-06-13T09:17:22Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_PGgThv4Kn7oGVQARALHy4389

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/validate.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/validate.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/validate.go: no such file or directory"}
```

## 2026-06-13T09:17:22Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_MdF6qf6hyqwgRbg2FAhg8902

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/defaults.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/defaults.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/defaults.go: no such file or directory"}
```

## 2026-06-13T09:17:22Z · read_file · rejected

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_03_LrECYDwEiFmRfaE9EPce5427

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/loader.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/loader.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/loader.go: no such file or directory"}
```

## 2026-06-13T09:17:22Z · bash · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_04_IidPn72IGkQs5xZb8WfV1180

### Arguments

```json
{"command": "find /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage -type f -name \"*.go\" -not -name \"*_test.go\" | sort"}
```

### Result

```text
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/conversation_history.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/models.go
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","resolved_command_path":"/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage","outcome_summary":"/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/conversation_history.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/models.go"}
```

## 2026-06-13T09:17:25Z · bash · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_wCeEQATrDHrst7I8fsEP7872

### Arguments

```json
{"command": "find /Users/bytedance/golang_pro/nano_cc/go-agent -type f -name \"*.go\" -not -name \"*_test.go\" -not -path \"*/vendor/*\" | sort"}
```

### Result

```text
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/transport.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/conversation_memory_strategy.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/full_history_summarization.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/message_window_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/reactive_compact.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/recent_tool_result_retention.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/token_estimator.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/tool_result_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/context_compression.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/conversation_flow.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/conversation_memory.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/helpers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hook_bridge.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/manager.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/stop.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/tool.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/types.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/user_prompt.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/memory.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/prompt_builder.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/runtime.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/subagent.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/summarizer.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/todo_reminder.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/tool_registry.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/conversation_history.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/models.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/assistant/prompt.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/cli/root.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/helpers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/local_config.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/paths.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/idgen/idgen.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/llm/client.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/llm/errors.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/bootstrap.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/memory_store.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/persisted_output_files.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/session_history_store.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/store.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/tool_result_log.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/logger/logger.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/safety/path.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/sessions/skill.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/textutil/textutil.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/bash.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/definitions.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/file_ops.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/handlers.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/load_skill.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/path_guard.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/persisted_output.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/runtime_env.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/todo_write.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tools/validation.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/app.go
/Users/bytedance/golang_pro/nano_cc/go-agent/internal/tui/events.go
/Users/bytedance/golang_pro/nano_cc/go-agent/main.go
/Users/bytedance/golang_pro/nano_cc/go-agent/workspace/basic.go
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","resolved_command_path":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/config.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/manager.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/mcp/transport.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/compression.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/compression/conversation_memory_strategy.go\n/Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime"}
```

## 2026-06-13T09:17:29Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_VyEg4C0TROje1q6Zhbd45390

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/config.go"}
```

### Result

```text
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"-"`
	ModelID string `json:"model_id"`
}

type fileConfig struct {
	Config
	AppHome                            string `json:"app_home"`
	BuiltinSkillsDir                   string `json:"builtin_skills_dir"`
	CommandBinDir                      string `json:"command_bin_dir"`
	CommandScriptDir                   string `json:"command_script_dir"`
	SystemPromptPath                   string `json:"system_prompt_path"`
	WorkspaceRoot                      string `json:"workspace_root"`
	AllowedTools                       string `json:"allowed_tools"`
	BashAllowOutsideWorkspace          bool   `json:"bash_allow_outside_workspace"`
	BashAllowDangerousCommands         bool   `json:"bash_allow_dangerous_commands"`
	ConversationLockTTLSeconds         int    `json:"conversation_lock_ttl_seconds"`
	MemoryWorkTimeoutSeconds           int    `json:"memory_work_timeout_seconds"`
	ConversationLockWaitTimeoutSeconds int    `json:"conversation_lock_wait_timeout_seconds"`
}

type AppConfig struct {
	LLM                         Config
	AppHome                     string
	BuiltinSkillsDir            string
	CommandBinDir               string
	CommandScriptDir            string
	SystemPromptPath            string
	WorkspaceRoot               string
	LogsDir                     string
	AllowedTools                []string
	BashAllowOutsideWorkspace   bool
	BashAllowDangerousCommands  bool
	ConversationLockTTL         time.Duration
	MemoryWorkTimeout           time.Duration
	ConversationLockWaitTimeout time.Duration
}

func loadConfigFile() (fileConfig, error) {
	data, err := os.ReadFile(configFilePath())
	if err != nil {
		return fileConfig{}, err
	}

	var cfg fileConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fileConfig{}, fmt.Errorf("failed to parse workspace config.json: %w", err)
	}
	return cfg, nil
}

func configFilePath() string {
	return "config.json"
}

func loadLLMConfig(fileCfg fileConfig) (Config, error) {
	cfg := Config{
		BaseURL: strings.TrimSpace(fileCfg.Config.BaseURL),
		APIKey:  strings.TrimSpace(getenv("OPENAI_API_KEY")),
		ModelID: strings.TrimSpace(fileCfg.Config.ModelID),
	}

	if cfg.BaseURL == "" || cfg.ModelID == "" {
		return Config{}, fmt.Errorf("missing LLM config; set base_url, model_id in config.json")
	}
	if cfg.APIKey == "" {
		return Config{}, fmt.Errorf("missing LLM api key; set OPENAI_API_KEY environment variable")
	}

	return cfg, nil
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package config\n\nimport (\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"os\"\n\t\"strings\"\n\t\"time\"\n)\n\ntype Config struct {\n\tBaseURL string `json:\"base_url\"`\n\tAPIKey  string `json:\"-\"`\n\tModelID string `json:\"model_id\"`\n}\n\ntype fileConfig struct {\n\tConfig\n\tAppHome                            string `json:\"app_home\"`\n\tBuiltinSkillsDir                   string `json:\"builtin_skills_dir\"`\n\tCommandBinDir                      string `json:\"command_bin_dir\"`\n\tCommandScriptDir                   string `json:\"command_script_dir\"`\n"}
```

## 2026-06-13T09:17:29Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_6KJ1wv4bbzjA3FTLJO9J6281

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/helpers.go"}
```

### Result

```text
package config

import (
	"os"
	"strings"
)

func getenv(key string) string {
	return os.Getenv(key)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func parseCSVList(value string) []string {
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			items = append(items, part)
		}
	}
	return items
}

func intOrDefault(value, fallback int) int {
	if value == 0 {
		return fallback
	}
	return value
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package config\n\nimport (\n\t\"os\"\n\t\"strings\"\n)\n\nfunc getenv(key string) string {\n\treturn os.Getenv(key)\n}\n\nfunc firstNonEmpty(values ...string) string {\n\tfor _, value := range values {\n\t\tif strings.TrimSpace(value) != \"\" {\n\t\t\treturn value\n\t\t}\n\t}\n\treturn \"\"\n}\n\nfunc parseCSVList(value string) []string {\n\tparts := strings.Split(value, \",\")\n\titems := make([]string, 0, len(parts))\n\tfor _, part := range parts {\n\t\tpart = strings.TrimSpace(part)\n\t\tif part != \"\" {\n\t\t\titems = append(items, part)\n\t\t}\n\t}\n\tretu"}
```

## 2026-06-13T09:17:29Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_mTfjQoQoat9utJn9K1sV9263

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/local_config.go"}
```

### Result

```text
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const defaultLocalAllowedTools = "load_skill,bash,read_file,write_file,edit_file,todo_write,spawn_subagent"

type linkSettings struct {
	Env linkSettingsEnv `json:"env"`
}

type linkSettingsEnv struct {
	OpenAuthToken string `json:"open_auth_token"`
	OpenModel     string `json:"open_model"`
	OpenBaseURL   string `json:"open_base_url"`
}

type LinkMarkdownContext struct {
	UserPath         string
	UserContent      string
	WorkspacePath    string
	WorkspaceContent string
}

func LoadLocalConfig(cwd string) (AppConfig, error) {
	fileCfg, err := loadConfigFile()
	if err != nil {
		return AppConfig{}, err
	}
	appHome, err := resolveAppHome(fileCfg)
	if err != nil {
		return AppConfig{}, err
	}
	llm, err := loadLocalLLMConfig(fileCfg)
	if err != nil {
		return AppConfig{}, err
	}
	runtimeDirs, err := resolveRuntimePaths(appHome, fileCfg)
	if err != nil {
		return AppConfig{}, err
	}
	workspaceRoot, err := resolveLocalWorkspaceRoot(cwd)
	if err != nil {
		return AppConfig{}, err
	}
	allowedTools := firstNonEmpty(fileCfg.AllowedTools, defaultLocalAllowedTools)
	return AppConfig{
		LLM:                         llm,
		AppHome:                     appHome,
		BuiltinSkillsDir:            runtimeDirs.builtinSkillsDir,
		CommandBinDir:               runtimeDirs.commandBinDir,
		CommandScriptDir:            runtimeDirs.commandScriptDir,
		SystemPromptPath:            runtimeDirs.systemPromptPath,
		WorkspaceRoot:               workspaceRoot,
		LogsDir:                     runtimeDirs.logsDir,
		AllowedTools:                parseCSVList(allowedTools),
		BashAllowOutsideWorkspace:   fileCfg.BashAllowOutsideWorkspace,
		BashAllowDangerousCommands:  fileCfg.BashAllowDangerousCommands,
		ConversationLockTTL:         time.Duration(intOrDefault(fileCfg.ConversationLockTTLSeconds, 30)) * time.Second,
		MemoryWorkTimeout:           time.Duration(intOrDefault(fileCfg.MemoryWorkTimeoutSeconds, 110)) * time.Second,
		ConversationLockWaitTimeout: time.Duration(intOrDefault(fileCfg.ConversationLockWaitTimeoutSeconds, 130)) * time.Second,
	}, nil
}

func loadLocalLLMConfig(_ fileConfig) (Config, error) {
	return loadLinkLLMConfig()
}

func LinkSettingsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	return filepath.Join(home, ".link", "settings.json"), nil
}

func LinkSkillsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	return filepath.Join(home, ".link", "skills"), nil
}

func LinkSessionDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	return filepath.Join(home, "link", "session"), nil
}

func LinkMarkdownPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	return filepath.Join(home, ".link", "LINK.MD"), nil
}

func WorkspaceLinkSkillsDir(workspaceRoot string) string {
	return filepath.Join(strings.TrimSpace(workspaceRoot), ".link", "skills")
}

func WorkspaceLinkMarkdownPath(workspaceRoot string) string {
	return filepath.Join(strings.TrimSpace(workspaceRoot), ".link", "LINK.MD")
}

func WorkspaceMCPConfigPath(workspaceRoot string) string {
	return filepath.Join(strings.TrimSpace(workspaceRoot), ".link", ".mcp.json")
}

func loadLinkLLMConfig() (Config, error) {
	path, err := LinkSettingsPath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, fmt.Errorf("missing LLM settings; create %s with env.open_auth_token, env.open_model, env.open_base_url", path)
		}
		return Config{}, fmt.Errorf("read LLM settings %s: %w", path, err)
	}
	var settings linkSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return Config{}, fmt.Errorf("parse LLM settings %s: %w", path, err)
	}
	cfg := Config{
		BaseURL: strings.TrimSpace(settings.Env.OpenBaseURL),
		APIKey:  strings.TrimSpace(settings.Env.OpenAuthToken),
		ModelID: strings.TrimSpace(settings.Env.OpenModel),
	}
	missing := make([]string, 0, 3)
	if cfg.APIKey == "" {
		missing = append(missing, "env.open_auth_token")
	}
	if cfg.ModelID == "" {
		missing = append(missing, "env.open_model")
	}
	if cfg.BaseURL == "" {
		missing = append(missing, "env.open_base_url")
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return Config{}, fmt.Errorf("missing LLM settings in %s: %s", path, strings.Join(missing, ", "))
	}
	return cfg, nil
}

func LoadLinkMarkdownContext(workspaceRoot string) (LinkMarkdownContext, error) {
	userPath, err := LinkMarkdownPath()
	if err != nil {
		return LinkMarkdownContext{}, err
	}
	workspacePath := WorkspaceLinkMarkdownPath(workspaceRoot)
	userContent, err := readOptionalMarkdown(userPath)
	if err != nil {
		return LinkMarkdownContext{}, err
	}
	workspaceContent, err := readOptionalMarkdown(workspacePath)
	if err != nil {
		return LinkMarkdownContext{}, err
	}
	return LinkMarkdownContext{UserPath: userPath, UserContent: userContent, WorkspacePath: workspacePath, WorkspaceContent: workspaceContent}, nil
}

func readOptionalMarkdown(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("read LINK.MD %s: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("read LINK.MD %s: not a regular file", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read LINK.MD %s: %w", path, err)
	}
	return strings.TrimSpace(string(data)), nil
}

func resolveLocalWorkspaceRoot(cwd string) (string, error) {
	if strings.TrimSpace(cwd) == "" {
		cwd = "."
	}
	resolved, err := filepath.Abs(cwd)
	if err != nil {
		return "", fmt.Errorf("resolve cwd: %w", err)
	}
	clean := filepath.Clean(resolved)
	info, err := os.Stat(clean)
	if err != nil {
		return "", fmt.Errorf("stat cwd: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("cwd is not a directory: %s", clean)
	}
	return clean, nil
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package config\n\nimport (\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"os\"\n\t\"path/filepath\"\n\t\"sort\"\n\t\"strings\"\n\t\"time\"\n)\n\nconst defaultLocalAllowedTools = \"load_skill,bash,read_file,write_file,edit_file,todo_write,spawn_subagent\"\n\ntype linkSettings struct {\n\tEnv linkSettingsEnv `json:\"env\"`\n}\n\ntype linkSettingsEnv struct {\n\tOpenAuthToken string `json:\"open_auth_token\"`\n\tOpenModel     string `json:\"open_model\"`\n\tOpenBaseURL   string `json:\"open_base_url\"`\n}\n\ntype LinkMarkdownContext struct {\n\tUserPath         string"}
```

## 2026-06-13T09:17:29Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_03_t1PAUmGCtRs9tf1XgkZO3913

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/paths.go"}
```

### Result

```text
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type runtimePaths struct {
	workspaceRoot    string
	builtinSkillsDir string
	commandBinDir    string
	commandScriptDir string
	systemPromptPath string
	logsDir          string
}

type runtimeAssetSpec struct {
	label     string
	fileValue string
	subdir    string
}

func resolveAppHome(fileCfg fileConfig) (string, error) {
	appHome := firstNonEmpty(fileCfg.AppHome, ".")
	resolved, err := filepath.Abs(appHome)
	if err != nil {
		return "", fmt.Errorf("resolve app_home: %w", err)
	}
	return filepath.Clean(resolved), nil
}

func resolvePath(appHome, pathValue string) (string, error) {
	if strings.TrimSpace(pathValue) == "" {
		pathValue = "."
	}
	if !filepath.IsAbs(pathValue) {
		pathValue = filepath.Join(appHome, pathValue)
	}
	resolved, err := filepath.Abs(pathValue)
	if err != nil {
		return "", err
	}
	return filepath.Clean(resolved), nil
}

func resolveWorkspaceRoot(appHome, configured string) (string, error) {
	configured = strings.TrimSpace(configured)
	if configured != "" {
		return resolvePath(appHome, configured)
	}
	return resolvePath(appHome, "workspace")
}

func resolveSystemPromptPath(appHome string, fileCfg fileConfig) (string, error) {
	configured := strings.TrimSpace(fileCfg.SystemPromptPath)
	if configured == "" {
		configured = "system_prompt.md"
	}
	return resolvePath(appHome, configured)
}

func resolveRuntimePaths(appHome string, fileCfg fileConfig) (runtimePaths, error) {
	workspaceRoot, err := resolveWorkspaceRoot(appHome, strings.TrimSpace(fileCfg.WorkspaceRoot))
	if err != nil {
		return runtimePaths{}, fmt.Errorf("resolve workspace root: %w", err)
	}
	systemPromptPath, err := resolveSystemPromptPath(appHome, fileCfg)
	if err != nil {
		return runtimePaths{}, fmt.Errorf("resolve system prompt path: %w", err)
	}
	logsDir, err := resolvePath(appHome, "logs")
	if err != nil {
		return runtimePaths{}, fmt.Errorf("resolve logs dir: %w", err)
	}
	paths := runtimePaths{workspaceRoot: workspaceRoot, systemPromptPath: systemPromptPath, logsDir: logsDir}
	for _, spec := range []runtimeAssetSpec{
		{label: "builtin skills dir", fileValue: fileCfg.BuiltinSkillsDir, subdir: "skills"},
		{label: "command bin dir", fileValue: fileCfg.CommandBinDir, subdir: "bin"},
		{label: "command script dir", fileValue: fileCfg.CommandScriptDir, subdir: "cmd"},
	} {
		resolved, err := resolveRuntimeAssetFromSpec(appHome, spec)
		if err != nil {
			return runtimePaths{}, err
		}
		switch spec.subdir {
		case "skills":
			paths.builtinSkillsDir = resolved
		case "bin":
			paths.commandBinDir = resolved
		case "cmd":
			paths.commandScriptDir = resolved
		}
	}
	return paths, nil
}

func resolveRuntimeAssetFromSpec(appHome string, spec runtimeAssetSpec) (string, error) {
	configured := strings.TrimSpace(spec.fileValue)
	resolved, err := resolveRuntimeAssetDir(appHome, configured, spec.subdir)
	if err != nil {
		return "", fmt.Errorf("resolve %s: %w", spec.label, err)
	}
	return resolved, nil
}

func resolveRuntimeAssetDir(appHome, configured, subdir string) (string, error) {
	expected, err := resolvePath(appHome, subdir)
	if err != nil {
		return "", err
	}

	configured = strings.TrimSpace(configured)
	if configured == "" {
		return expected, nil
	}

	resolved, err := resolvePath(appHome, configured)
	if err != nil {
		return "", err
	}
	if resolved != expected {
		return "", fmt.Errorf("runtime asset dir must stay under app home: expected %q", expected)
	}
	return resolved, nil
}

func EnsureAppLayout(cfg AppConfig) error {
	paths := []struct {
		label string
		path  string
	}{
		{label: "app home", path: cfg.AppHome},
		{label: "logs", path: cfg.LogsDir},
		{label: "builtin skills dir", path: cfg.BuiltinSkillsDir},
		{label: "command bin dir", path: cfg.CommandBinDir},
		{label: "command script dir", path: cfg.CommandScriptDir},
		{label: "workspace root", path: cfg.WorkspaceRoot},
	}
	for _, item := range paths {
		if strings.TrimSpace(item.path) == "" {
			continue
		}
		if err := os.MkdirAll(item.path, 0o755); err != nil {
			return fmt.Errorf("create %s: %w", item.label, err)
		}
	}
	return nil
}

func ValidateAppLayout(cfg AppConfig) error {
	paths := []struct {
		label string
		path  string
	}{
		{label: "app home", path: cfg.AppHome},
		{label: "builtin skills dir", path: cfg.BuiltinSkillsDir},
		{label: "command bin dir", path: cfg.CommandBinDir},
		{label: "command script dir", path: cfg.CommandScriptDir},
		{label: "workspace root", path: cfg.WorkspaceRoot},
	}
	for _, item := range paths {
		if strings.TrimSpace(item.path) == "" {
			return fmt.Errorf("%s is required", item.label)
		}
		info, err := os.Stat(item.path)
		if err != nil {
			return fmt.Errorf("stat %s: %w", item.label, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("%s is not a directory", item.label)
		}
	}
	return nil
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package config\n\nimport (\n\t\"fmt\"\n\t\"os\"\n\t\"path/filepath\"\n\t\"strings\"\n)\n\ntype runtimePaths struct {\n\tworkspaceRoot    string\n\tbuiltinSkillsDir string\n\tcommandBinDir    string\n\tcommandScriptDir string\n\tsystemPromptPath string\n\tlogsDir          string\n}\n\ntype runtimeAssetSpec struct {\n\tlabel     string\n\tfileValue string\n\tsubdir    string\n}\n\nfunc resolveAppHome(fileCfg fileConfig) (string, error) {\n\tappHome := firstNonEmpty(fileCfg.AppHome, \".\")\n\tresolved, err := filepath.Abs(appHome)\n\tif err != nil {\n"}
```

## 2026-06-13T09:17:29Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_04_pmqDZRZ2JDZKvpk89Poj2337

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/store.go"}
```

### Result

```text
package local

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"nano_cc/internal/agent/storage"
)

type Store struct {
	mu              sync.RWMutex
	conversations   map[string]storage.Conversation
	messages        map[string][]storage.Message
	cache           map[string][]storage.Message
	toolCalls       []storage.ToolCall
	persisted       map[string]storage.PersistedOutput
	persistedByHash map[string]string
	summaries       map[string]storage.ContextSummary
	modelHistory    map[string][]storage.Message
	memories        []storage.Memory
	convMemories    map[string][]storage.ConversationMemory
	memory          *MarkdownMemoryStore
	sessionHistory  *SessionHistoryStore
	workspaceRoot   string
	locks           map[string]string
}

func NewStore() *Store {
	return &Store{
		conversations:   make(map[string]storage.Conversation),
		messages:        make(map[string][]storage.Message),
		cache:           make(map[string][]storage.Message),
		persisted:       make(map[string]storage.PersistedOutput),
		persistedByHash: make(map[string]string),
		summaries:       make(map[string]storage.ContextSummary),
		modelHistory:    make(map[string][]storage.Message),
		convMemories:    make(map[string][]storage.ConversationMemory),
		locks:           make(map[string]string),
	}
}

func NewStoreWithMemory(workspaceRoot string) (*Store, error) {
	store := NewStore()
	memory, err := NewMarkdownMemoryStore(workspaceRoot)
	if err != nil {
		return nil, err
	}
	sessionHistory, err := NewSessionHistoryStore()
	if err != nil {
		return nil, err
	}
	store.memory = memory
	store.sessionHistory = sessionHistory
	store.workspaceRoot = workspaceRoot
	return store, nil
}

func (s *Store) CreateConversation(ctx context.Context, conversation storage.Conversation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if conversation.SessionID == "" {
		conversation.SessionID = conversation.ID
	}
	now := time.Now()
	if conversation.CreatedAt.IsZero() {
		conversation.CreatedAt = now
	}
	conversation.UpdatedAt = now
	s.conversations[conversation.ID] = conversation
	return nil
}

func (s *Store) GetConversationByID(ctx context.Context, conversationID string) (storage.Conversation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	conv, ok := s.conversations[conversationID]
	if !ok {
		return storage.Conversation{}, sql.ErrNoRows
	}
	return conv, nil
}

func (s *Store) UpdateConversationTitle(ctx context.Context, conversationID, title string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	conv, ok := s.conversations[conversationID]
	if !ok {
		return sql.ErrNoRows
	}
	conv.Title = title
	conv.UpdatedAt = time.Now()
	s.conversations[conversationID] = conv
	return nil
}

func (s *Store) TouchConversationActivity(ctx context.Context, conversationID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	conv, ok := s.conversations[conversationID]
	if !ok {
		return nil
	}
	conv.UpdatedAt = time.Now()
	s.conversations[conversationID] = conv
	return nil
}

func (s *Store) SetConversationHistory(ctx context.Context, conversationID string, messages []storage.Message) error {
	s.mu.Lock()
	conv := s.conversations[conversationID]
	s.messages[conversationID] = cloneMessages(messages)
	s.cache[conversationID] = cloneMessages(messages)
	s.mu.Unlock()
	if s.sessionHistory != nil && conv.ID != "" {
		return s.sessionHistory.SaveHistory(ctx, s.workspaceRoot, conv, messages)
	}
	return nil
}

func (s *Store) SetConversationCache(ctx context.Context, conversationID string, messages []storage.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache[conversationID] = cloneMessages(messages)
	return nil
}

func (s *Store) GetConversationCache(ctx context.Context, conversationID string) ([]storage.Message, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	messages, ok := s.cache[conversationID]
	return cloneMessages(messages), ok, nil
}

func (s *Store) ListMessagesByConversation(ctx context.Context, conversationID string, limit int) ([]storage.Message, error) {
	s.mu.RLock()
	conv := s.conversations[conversationID]
	messages := cloneMessages(s.messages[conversationID])
	s.mu.RUnlock()
	if len(messages) == 0 && s.sessionHistory != nil && conv.SessionID != "" {
		doc, err := s.sessionHistory.LoadHistory(ctx, conv.SessionID)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		if err == nil {
			messages = cloneMessages(doc.Messages)
			s.mu.Lock()
			s.messages[conversationID] = cloneMessages(messages)
			s.cache[conversationID] = cloneMessages(messages)
			s.mu.Unlock()
		}
	}
	if limit > 0 && len(messages) > limit {
		messages = messages[len(messages)-limit:]
	}
	return messages, nil
}

func (s *Store) CreateToolCall(ctx context.Context, tc storage.ToolCall) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if tc.CreatedAt.IsZero() {
		tc.CreatedAt = time.Now()
	}
	s.toolCalls = append(s.toolCalls, tc)
	return nil
}

func (s *Store) CreateSubagentMessage(ctx context.Context, message storage.SubagentMessage) error {
	return nil
}

func (s *Store) CreatePersistedOutput(ctx context.Context, output storage.PersistedOutput) error {
	s.mu.RLock()
	conv := s.conversations[output.ConversationID]
	workspaceRoot := s.workspaceRoot
	s.mu.RUnlock()
	if workspaceRoot != "" && conv.SessionID != "" {
		persisted, err := persistOutputToWorkspace(ctx, workspaceRoot, conv.SessionID, output)
		if err != nil {
			return err
		}
		output = persisted
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if output.CreatedAt.IsZero() {
		output.CreatedAt = time.Now()
	}
	s.persisted[output.ID] = output
	s.persistedByHash[persistedHashKey(output.ConversationID, output.UserID, output.MessageID, output.ToolCallID, output.Strategy, output.ContentSHA256)] = output.ID
	return nil
}

func (s *Store) GetPersistedOutputForConversation(ctx context.Context, id, userID, conversationID string) (storage.PersistedOutput, error) {
	s.mu.RLock()
	output, ok := s.persisted[id]
	if ok && output.UserID == userID && output.ConversationID == conversationID {
		s.mu.RUnlock()
		return output, nil
	}
	conv := s.conversations[conversationID]
	workspaceRoot := s.workspaceRoot
	s.mu.RUnlock()
	if workspaceRoot == "" || conv.SessionID == "" {
		return storage.PersistedOutput{}, sql.ErrNoRows
	}
	output, err := loadPersistedOutputFromWorkspace(ctx, workspaceRoot, conv.SessionID, id)
	if err != nil {
		return storage.PersistedOutput{}, err
	}
	if output.UserID != userID || output.ConversationID != conversationID {
		return storage.PersistedOutput{}, sql.ErrNoRows
	}
	s.mu.Lock()
	s.persisted[output.ID] = output
	s.persistedByHash[persistedHashKey(output.ConversationID, output.UserID, output.MessageID, output.ToolCallID, output.Strategy, output.ContentSHA256)] = output.ID
	s.mu.Unlock()
	return output, nil
}

func (s *Store) GetPersistedOutputByMessageHash(ctx context.Context, conversationID, userID, messageID, toolCallID, strategy, contentSHA256 string) (storage.PersistedOutput, error) {
	s.mu.RLock()
	id, ok := s.persistedByHash[persistedHashKey(conversationID, userID, messageID, toolCallID, strategy, contentSHA256)]
	if ok {
		output := s.persisted[id]
		s.mu.RUnlock()
		return output, nil
	}
	conv := s.conversations[conversationID]
	workspaceRoot := s.workspaceRoot
	s.mu.RUnlock()
	if workspaceRoot == "" || conv.SessionID == "" {
		return storage.PersistedOutput{}, sql.ErrNoRows
	}
	output, err := findPersistedOutputByMessageHashInWorkspace(ctx, workspaceRoot, conv.SessionID, conversationID, userID, messageID, toolCallID, strategy, contentSHA256)
	if err != nil {
		return storage.PersistedOutput{}, err
	}
	s.mu.Lock()
	s.persisted[output.ID] = output
	s.persistedByHash[persistedHashKey(output.ConversationID, output.UserID, output.MessageID, output.ToolCallID, output.Strategy, output.ContentSHA256)] = output.ID
	s.mu.Unlock()
	return output, nil
}

func (s *Store) CreateContextSummary(ctx context.Context, summary storage.ContextSummary) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.summaries[summary.ConversationID+"\x00"+summary.UserID+"\x00"+summary.SourceHistorySHA256] = summary
	return nil
}

func (s *Store) GetContextSummaryByHistoryHash(ctx context.Context, conversationID, userID, sourceHistorySHA256 string) (storage.ContextSummary, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	summary, ok := s.summaries[conversationID+"\x00"+userID+"\x00"+sourceHistorySHA256]
	if !ok {
		return storage.ContextSummary{}, sql.ErrNoRows
	}
	return summary, nil
}

func (s *Store) ListRelevantMemories(ctx context.Context, userID string) ([]storage.Memory, error) {
	if s.memory != nil {
		return s.memory.ListRelevantMemories(ctx, userID)
	}
	return nil, nil
}

func (s *Store) ListMemoriesByUserAndType(ctx context.Context, userID, memType string) ([]storage.Memory, error) {
	if s.memory != nil {
		return s.memory.ListMemoriesByUserAndType(ctx, userID, memType)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []storage.Memory
	for _, memory := range s.memories {
		if memory.UserID == userID && memory.Type == memType {
			result = append(result, memory)
		}
	}
	return result, nil
}

func (s *Store) ListProjectFactMemories(ctx context.Context, userID string) ([]storage.Memory, error) {
	if s.memory != nil {
		return s.memory.ListProjectFactMemories(ctx, userID)
	}
	return s.ListMemoriesByUserAndType(ctx, userID, "project_fact")
}

func (s *Store) InsertMemory(ctx context.Context, m storage.Memory) error {
	if s.memory != nil {
		return s.memory.InsertMemory(ctx, m)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.memories = append(s.memories, m)
	return nil
}

func (s *Store) CountMemoriesByUserAndType(ctx context.Context, userID, memType string) (int, error) {
	if s.memory != nil {
		return s.memory.CountMemoriesByUserAndType(ctx, userID, memType)
	}
	items, err := s.ListMemoriesByUserAndType(ctx, userID, memType)
	return len(items), err
}

func (s *Store) CountProjectFactMemories(ctx context.Context, userID string) (int, error) {
	if s.memory != nil {
		return s.memory.CountProjectFactMemories(ctx, userID)
	}
	items, err := s.ListProjectFactMemories(ctx, userID)
	return len(items), err
}
func (s *Store) DeleteOldestMemories(ctx context.Context, userID, memType string, n int) error {
	if s.memory != nil {
		return s.memory.DeleteOldestMemories(ctx, userID, memType, n)
	}
	return nil
}
func (s *Store) ReplaceMemoriesByUserAndType(ctx context.Context, userID, memType string, items []storage.Memory) error {
	if s.memory != nil {
		return s.memory.ReplaceMemoriesByUserAndType(ctx, userID, memType, items)
	}
	return nil
}
func (s *Store) ReplaceProjectFactMemories(ctx context.Context, userID string, items []storage.Memory) error {
	if s.memory != nil {
		return s.memory.ReplaceProjectFactMemories(ctx, userID, items)
	}
	return nil
}

func (s *Store) ListConversationMemories(ctx context.Context, conversationID string) ([]storage.ConversationMemory, error) {
	if s.memory != nil {
		sessionID := s.sessionIDForConversation(conversationID)
		return s.memory.ListConversationMemories(ctx, sessionID)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := append([]storage.ConversationMemory(nil), s.convMemories[conversationID]...)
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.Before(items[j].CreatedAt) })
	return items, nil
}

func (s *Store) ReplaceConversationMemories(ctx context.Context, conversationID, userID string, items []storage.ConversationMemory) error {
	if s.memory != nil {
		sessionID := s.sessionIDForConversation(conversationID)
		return s.memory.ReplaceConversationMemories(ctx, sessionID, userID, items)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.convMemories[conversationID] = append([]storage.ConversationMemory(nil), items...)
	return nil
}

func (s *Store) GetConversationModelHistory(ctx context.Context, conversationID string) ([]storage.Message, bool, error) {
	s.mu.RLock()
	conv := s.conversations[conversationID]
	messages, ok := s.modelHistory[conversationID]
	s.mu.RUnlock()
	if (!ok || len(messages) == 0) && s.sessionHistory != nil && conv.SessionID != "" {
		doc, err := s.sessionHistory.LoadModelHistory(ctx, conv.SessionID)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, false, nil
			}
			return nil, false, err
		}
		messages = doc.Messages
		ok = len(messages) > 0
		s.mu.Lock()
		s.modelHistory[conversationID] = cloneMessages(messages)
		s.mu.Unlock()
	}
	return cloneMessages(messages), ok, nil
}

func (s *Store) UpsertConversationModelHistory(ctx context.Context, conversationID, userID string, messages []storage.Message) error {
	s.mu.Lock()
	conv := s.conversations[conversationID]
	s.modelHistory[conversationID] = cloneMessages(messages)
	s.mu.Unlock()
	if s.sessionHistory != nil && conv.ID != "" {
		if conv.UserID == "" {
			conv.UserID = userID
		}
		return s.sessionHistory.SaveModelHistory(ctx, s.workspaceRoot, conv, messages)
	}
	return nil
}

func (s *Store) ListResumableSessions(ctx context.Context, currentWorkspace string) ([]storage.ResumableSession, error) {
	if s.sessionHistory == nil {
		return nil, nil
	}
	return s.sessionHistory.ListSessions(ctx, currentWorkspace)
}

func (s *Store) ResumeSession(ctx context.Context, sessionID, currentWorkspace string, user storage.User) (storage.Conversation, []storage.Message, error) {
	if s.sessionHistory == nil {
		return storage.Conversation{}, nil, fmt.Errorf("session history store is not initialized")
	}
	doc, err := s.sessionHistory.LoadHistory(ctx, sessionID)
	if err != nil {
		return storage.Conversation{}, nil, err
	}
	current, err := normalizeWorkspacePath(currentWorkspace)
	if err != nil {
		return storage.Conversation{}, nil, err
	}
	stored, err := normalizeWorkspacePath(doc.WorkspaceRoot)
	if err != nil {
		return storage.Conversation{}, nil, err
	}
	if current != stored {
		return storage.Conversation{}, nil, fmt.Errorf("历史会话属于 %s，请 cd 到该目录或使用 --cwd %s 后再执行 /resume", doc.WorkspaceRoot, doc.WorkspaceRoot)
	}
	conv := storage.Conversation{ID: doc.ConversationID, SessionID: doc.SessionID, UserID: user.ID, Title: doc.Title, CreatedAt: doc.CreatedAt, UpdatedAt: doc.UpdatedAt}
	if conv.ID == "" {
		conv.ID = doc.SessionID
	}
	if conv.UserID == "" {
		conv.UserID = doc.UserID
	}
	if conv.UserID == "" {
		conv.UserID = LocalUserID
	}
	history := cloneMessages(doc.Messages)
	var modelHistory []storage.Message
	if modelDoc, err := s.sessionHistory.LoadModelHistory(ctx, sessionID); err == nil {
		modelHistory = cloneMessages(modelDoc.Messages)
	} else if !os.IsNotExist(err) {
		return storage.Conversation{}, nil, err
	}
	s.mu.Lock()
	s.conversations[conv.ID] = conv
	s.messages[conv.ID] = cloneMessages(history)
	s.cache[conv.ID] = cloneMessages(history)
	if len(modelHistory) > 0 {
		s.modelHistory[conv.ID] = cloneMessages(modelHistory)
	} else {
		delete(s.modelHistory, conv.ID)
	}
	s.mu.Unlock()
	return conv, history, nil
}

func (s *Store) AcquireConversationLock(ctx context.Context, conversationID, token string, ttl, waitTimeout time.Duration) (bool, error) {
	key := conversationID
	if s.memory != nil {
		key = s.memory.ProjectLockKey(s.sessionIDForConversation(conversationID))
	}
	deadline := time.Now().Add(waitTimeout)
	for {
		s.mu.Lock()
		if s.locks[key] == "" {
			s.locks[key] = token
			s.mu.Unlock()
			return true, nil
		}
		s.mu.Unlock()
		if waitTimeout <= 0 || time.Now().After(deadline) {
			return false, nil
		}
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(10 * time.Millisecond):
		}
	}
}
func (s *Store) RenewConversationLock(ctx context.Context, conversationID, token string, ttl time.Duration) (bool, error) {
	key := conversationID
	if s.memory != nil {
		key = s.memory.ProjectLockKey(s.sessionIDForConversation(conversationID))
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.locks[key] == token, nil
}
func (s *Store) ReleaseConversationLock(ctx context.Context, conversationID, token string) error {
	key := conversationID
	if s.memory != nil {
		key = s.memory.ProjectLockKey(s.sessionIDForConversation(conversationID))
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.locks[key] == token {
		delete(s.locks, key)
	}
	return nil
}

func (s *Store) sessionIDForConversation(conversationID string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if conv, ok := s.conversations[conversationID]; ok && conv.SessionID != "" {
		return conv.SessionID
	}
	return conversationID
}

func cloneMessages(messages []storage.Message) []storage.Message {
	if len(messages) == 0 {
		return nil
	}
	cloned := make([]storage.Message, len(messages))
	copy(cloned, messages)
	for i := range cloned {
		cloned[i].ToolCalls = append([]storage.MessageToolCall(nil), cloned[i].ToolCalls...)
	}
	return cloned
}

func persistedHashKey(parts ...string) string {
	key := ""
	for _, part := range parts {
		key += part + "\x00"
	}
	return key
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package local\n\nimport (\n\t\"context\"\n\t\"database/sql\"\n\t\"fmt\"\n\t\"os\"\n\t\"sort\"\n\t\"sync\"\n\t\"time\"\n\n\t\"nano_cc/internal/agent/storage\"\n)\n\ntype Store struct {\n\tmu              sync.RWMutex\n\tconversations   map[string]storage.Conversation\n\tmessages        map[string][]storage.Message\n\tcache           map[string][]storage.Message\n\ttoolCalls       []storage.ToolCall\n\tpersisted       map[string]storage.PersistedOutput\n\tpersistedByHash map[string]string\n\tsummaries       map[string]storage.ContextSummary\n\tmodelHis"}
```

## 2026-06-13T09:17:29Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_05_D9NZSiZgnqzxmy7gvPrW4531

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/bootstrap.go"}
```

### Result

```text
package local

import (
	"context"
	"fmt"
	"time"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/runtime"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/assistant"
	"nano_cc/internal/config"
	"nano_cc/internal/idgen"
	"nano_cc/internal/llm"
	"nano_cc/internal/logger"
	"nano_cc/internal/sessions"
)

const LocalUserID = "local-user"

type Bundle struct {
	Runtime      *runtime.Service
	Store        *Store
	MCP          *mcp.Manager
	User         storage.User
	Conversation storage.Conversation
	CWD          string
	SkillCount   int
	Skills       []sessions.SkillSummary
	MCPToolCount int
	MCPServers   []mcp.ServerStatus
}

func Bootstrap(ctx context.Context, cwd string) (*Bundle, error) {
	cfg, err := config.LoadLocalConfig(cwd)
	if err != nil {
		return nil, err
	}
	if err := config.EnsureAppLayout(cfg); err != nil {
		return nil, err
	}
	if err := config.ValidateAppLayout(cfg); err != nil {
		return nil, err
	}
	if err := logger.InitFileLoggerAt(cfg.LogsDir); err != nil {
		logger.Warn(fmt.Sprintf("failed to init file logger: %v", err))
	}
	userSkillsDir, err := config.LinkSkillsDir()
	if err != nil {
		return nil, err
	}
	builtinSkills, err := sessions.LoadSkillsFromDirs([]sessions.SkillDir{
		{Path: userSkillsDir, Source: "user"},
		{Path: config.WorkspaceLinkSkillsDir(cfg.WorkspaceRoot), Source: "workspace"},
	})
	if err != nil {
		return nil, fmt.Errorf("load .link skills: %w", err)
	}
	basePrompt, err := assistant.LoadBaseSystemPrompt(cfg.SystemPromptPath)
	if err != nil {
		return nil, fmt.Errorf("load system prompt: %w", err)
	}
	linkMarkdown, err := config.LoadLinkMarkdownContext(cfg.WorkspaceRoot)
	if err != nil {
		return nil, fmt.Errorf("load LINK.MD context: %w", err)
	}
	store, err := NewStoreWithMemory(cfg.WorkspaceRoot)
	if err != nil {
		return nil, fmt.Errorf("init memory store: %w", err)
	}
	client := llm.NewDeepseekClient(cfg.LLM.BaseURL, cfg.LLM.APIKey)
	runtimeService := runtime.NewService(store, cfg, client)
	runtimeService.EnableMemory = true
	runtimeService.SetBuiltinSkills(builtinSkills)
	runtimeService.SetBasePrompt(basePrompt)
	runtimeService.SetLinkMarkdownContext(linkMarkdown)
	mcpManager := mcp.NewManager()
	workspaceMCPServers, err := mcp.LoadWorkspaceConfig(config.WorkspaceMCPConfigPath(cfg.WorkspaceRoot))
	if err != nil {
		mcpManager.Close()
		return nil, fmt.Errorf("load workspace mcp config: %w", err)
	}
	mcpManager.SetWorkspaceServers(ctx, workspaceMCPServers)
	runtimeService.SetMCPManager(mcpManager)
	user := storage.User{ID: LocalUserID, Email: "local@go-agent", Username: "local", MemoryEnabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	conversation := storage.Conversation{ID: idgen.New("conv"), SessionID: idgen.UUID(), UserID: user.ID, RootMessageID: idgen.New("msg"), Title: "TUI 会话", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := store.CreateConversation(ctx, conversation); err != nil {
		mcpManager.Close()
		return nil, err
	}
	mcpManager.EnsureWorkspaceSessions(ctx)
	mcpTools := mcpManager.ToolsForUser(user.ID)
	mcpSnapshot := mcpManager.Snapshot(user.ID)
	skills := builtinSkills.Summaries()
	return &Bundle{Runtime: runtimeService, Store: store, MCP: mcpManager, User: user, Conversation: conversation, CWD: cfg.WorkspaceRoot, SkillCount: len(skills), Skills: skills, MCPToolCount: len(mcpTools), MCPServers: mcpSnapshot.Servers}, nil
}

func (b *Bundle) Close() {
	if b != nil && b.MCP != nil {
		b.MCP.Close()
	}
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package local\n\nimport (\n\t\"context\"\n\t\"fmt\"\n\t\"time\"\n\n\t\"nano_cc/internal/agent/mcp\"\n\t\"nano_cc/internal/agent/runtime\"\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/assistant\"\n\t\"nano_cc/internal/config\"\n\t\"nano_cc/internal/idgen\"\n\t\"nano_cc/internal/llm\"\n\t\"nano_cc/internal/logger\"\n\t\"nano_cc/internal/sessions\"\n)\n\nconst LocalUserID = \"local-user\"\n\ntype Bundle struct {\n\tRuntime      *runtime.Service\n\tStore        *Store\n\tMCP          *mcp.Manager\n\tUser         storage.User\n\tConversation storage.Co"}
```

## 2026-06-13T09:17:29Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_06_Ioz5IQJ7hi3aYECMAidz3787

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/memory_store.go"}
```

### Result

```text
package local

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/idgen"
)

const memoryIndexHeader = "# Memory Index\n"

type MarkdownMemoryStore struct {
	mu          sync.Mutex
	rootDir     string
	indexPath   string
	sessionsDir string
	projectName string
}

func NewMarkdownMemoryStore(workspaceRoot string) (*MarkdownMemoryStore, error) {
	root := filepath.Join(strings.TrimSpace(workspaceRoot), ".link", "memory")
	store := &MarkdownMemoryStore{
		rootDir:     root,
		indexPath:   filepath.Join(root, "memory.md"),
		sessionsDir: filepath.Join(root, "sessions"),
		projectName: sanitizeName(filepath.Base(filepath.Clean(workspaceRoot))),
	}
	if store.projectName == "" {
		store.projectName = "project"
	}
	if err := store.EnsureLayout(); err != nil {
		return nil, err
	}
	return store, nil
}

func (m *MarkdownMemoryStore) EnsureLayout() error {
	if err := os.MkdirAll(m.sessionsDir, 0o755); err != nil {
		return err
	}
	if _, err := os.Stat(m.indexPath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	return atomicWriteFile(m.indexPath, []byte(memoryIndexHeader+"\n"), 0o644)
}

func (m *MarkdownMemoryStore) ListRelevantMemories(ctx context.Context, userID string) ([]storage.Memory, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.listRelevantMemoriesLocked()
}

func (m *MarkdownMemoryStore) InsertMemory(ctx context.Context, mem storage.Memory) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if mem.ID == "" {
		mem.ID = idgen.New("mem")
	}
	if mem.Type == "" {
		mem.Type = "user_preference"
	}
	if mem.CreatedAt.IsZero() {
		mem.CreatedAt = time.Now()
	}
	mem.UpdatedAt = time.Now()
	filename := m.uniqueMemoryFilenameLocked(mem.Name)
	if err := m.writeMemoryFileLocked(filename, mem); err != nil {
		return err
	}
	return m.rewriteIndexLocked()
}

func (m *MarkdownMemoryStore) ListMemoriesByUserAndType(ctx context.Context, userID, memType string) ([]storage.Memory, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	items, err := m.listRelevantMemoriesLocked()
	if err != nil {
		return nil, err
	}
	return filterMemories(items, memType), nil
}

func (m *MarkdownMemoryStore) ListProjectFactMemories(ctx context.Context, userID string) ([]storage.Memory, error) {
	return m.ListMemoriesByUserAndType(ctx, userID, "project_fact")
}

func (m *MarkdownMemoryStore) CountMemoriesByUserAndType(ctx context.Context, userID, memType string) (int, error) {
	items, err := m.ListMemoriesByUserAndType(ctx, userID, memType)
	return len(items), err
}

func (m *MarkdownMemoryStore) CountProjectFactMemories(ctx context.Context, userID string) (int, error) {
	items, err := m.ListProjectFactMemories(ctx, userID)
	return len(items), err
}

func (m *MarkdownMemoryStore) DeleteOldestMemories(ctx context.Context, userID, memType string, n int) error {
	if n <= 0 {
		return nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	entries, err := m.indexEntriesLocked()
	if err != nil {
		return err
	}
	deleted := 0
	kept := entries[:0]
	for _, entry := range entries {
		mem, err := m.readMemoryFileLocked(entry.Path)
		if err == nil && mem.Type == memType && deleted < n {
			_ = os.Remove(filepath.Join(m.rootDir, entry.Path))
			deleted++
			continue
		}
		kept = append(kept, entry)
	}
	return m.writeIndexEntriesLocked(kept)
}

func (m *MarkdownMemoryStore) ReplaceMemoriesByUserAndType(ctx context.Context, userID, memType string, items []storage.Memory) error {
	return m.replaceMemoriesByType(memType, items)
}

func (m *MarkdownMemoryStore) ReplaceProjectFactMemories(ctx context.Context, userID string, items []storage.Memory) error {
	return m.replaceMemoriesByType("project_fact", items)
}

func (m *MarkdownMemoryStore) ListConversationMemories(ctx context.Context, sessionID string) ([]storage.ConversationMemory, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	path := m.sessionPath(sessionID)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	_, body := splitFrontMatter(string(data))
	return parseConversationMemoryBody(sessionID, body), nil
}

func (m *MarkdownMemoryStore) ReplaceConversationMemories(ctx context.Context, sessionID, userID string, items []storage.ConversationMemory) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.writeSessionMemoryLocked(sessionID, items)
}

func (m *MarkdownMemoryStore) ProjectLockKey(sessionID string) string {
	return m.projectName + ":" + sanitizeName(sessionID)
}

func (m *MarkdownMemoryStore) replaceMemoriesByType(memType string, items []storage.Memory) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	entries, err := m.indexEntriesLocked()
	if err != nil {
		return err
	}
	kept := entries[:0]
	for _, entry := range entries {
		mem, err := m.readMemoryFileLocked(entry.Path)
		if err == nil && mem.Type == memType {
			_ = os.Remove(filepath.Join(m.rootDir, entry.Path))
			continue
		}
		kept = append(kept, entry)
	}
	if err := m.writeIndexEntriesLocked(kept); err != nil {
		return err
	}
	for _, item := range items {
		if item.ID == "" {
			item.ID = idgen.New("mem")
		}
		if item.Type == "" {
			item.Type = memType
		}
		filename := m.uniqueMemoryFilenameLocked(item.Name)
		if err := m.writeMemoryFileLocked(filename, item); err != nil {
			return err
		}
	}
	return m.rewriteIndexLocked()
}

type memoryIndexEntry struct {
	Name        string
	Path        string
	Description string
}

func (m *MarkdownMemoryStore) listRelevantMemoriesLocked() ([]storage.Memory, error) {
	entries, err := m.indexEntriesLocked()
	if err != nil {
		return nil, err
	}
	result := make([]storage.Memory, 0, len(entries))
	for _, entry := range entries {
		mem, err := m.readMemoryFileLocked(entry.Path)
		if err != nil {
			continue
		}
		result = append(result, mem)
	}
	return result, nil
}

func (m *MarkdownMemoryStore) indexEntriesLocked() ([]memoryIndexEntry, error) {
	data, err := os.ReadFile(m.indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var entries []memoryIndexEntry
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "- [") {
			continue
		}
		closeName := strings.Index(line, "](")
		closePath := strings.Index(line, ")")
		if closeName < 3 || closePath <= closeName+2 {
			continue
		}
		name := line[3:closeName]
		path := line[closeName+2 : closePath]
		desc := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line[closePath+1:]), "—"))
		if !safeRelativeMemoryPath(path) {
			continue
		}
		entries = append(entries, memoryIndexEntry{Name: name, Path: filepath.ToSlash(filepath.Clean(path)), Description: desc})
	}
	return entries, scanner.Err()
}

func (m *MarkdownMemoryStore) rewriteIndexLocked() error {
	entries, err := m.indexEntriesLocked()
	if err != nil {
		return err
	}
	seen := make(map[string]memoryIndexEntry, len(entries))
	for _, entry := range entries {
		if mem, err := m.readMemoryFileLocked(entry.Path); err == nil {
			entry.Name = mem.Name
			entry.Description = mem.Description
			seen[entry.Path] = entry
		}
	}
	files, err := filepath.Glob(filepath.Join(m.rootDir, "*.md"))
	if err != nil {
		return err
	}
	for _, file := range files {
		if filepath.Clean(file) == filepath.Clean(m.indexPath) {
			continue
		}
		rel, err := filepath.Rel(m.rootDir, file)
		if err != nil {
			continue
		}
		rel = filepath.ToSlash(rel)
		mem, err := m.readMemoryFileLocked(rel)
		if err != nil {
			continue
		}
		seen[rel] = memoryIndexEntry{Name: mem.Name, Path: rel, Description: mem.Description}
	}
	entries = entries[:0]
	for _, entry := range seen {
		entries = append(entries, entry)
	}
	return m.writeIndexEntriesLocked(entries)
}

func (m *MarkdownMemoryStore) writeIndexEntriesLocked(entries []memoryIndexEntry) error {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Name == entries[j].Name {
			return entries[i].Path < entries[j].Path
		}
		return entries[i].Name < entries[j].Name
	})
	var b strings.Builder
	b.WriteString(memoryIndexHeader)
	b.WriteString("\n")
	for _, entry := range entries {
		b.WriteString(fmt.Sprintf("- [%s](%s) — %s\n", entry.Name, entry.Path, entry.Description))
	}
	return atomicWriteFile(m.indexPath, []byte(b.String()), 0o644)
}

func (m *MarkdownMemoryStore) uniqueMemoryFilenameLocked(name string) string {
	base := sanitizeName(name)
	if base == "" {
		base = "memory"
	}
	path := base + ".md"
	if _, err := os.Stat(filepath.Join(m.rootDir, path)); os.IsNotExist(err) {
		return path
	}
	return base + "-" + idgen.Hex()[:8] + ".md"
}

func (m *MarkdownMemoryStore) writeMemoryFileLocked(path string, mem storage.Memory) error {
	if !safeRelativeMemoryPath(path) {
		return fmt.Errorf("unsafe memory path: %s", path)
	}
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("name: " + yamlScalar(mem.Name) + "\n")
	b.WriteString("description: " + yamlScalar(mem.Description) + "\n")
	b.WriteString("metadata:\n")
	b.WriteString("  node_type: memory\n")
	b.WriteString("  type: " + yamlScalar(mem.Type) + "\n")
	b.WriteString("  project: " + yamlScalar(m.projectName) + "\n")
	b.WriteString("  originSessionId: " + yamlScalar(mem.ID) + "\n")
	b.WriteString("---\n\n")
	b.WriteString(strings.TrimSpace(mem.Body))
	b.WriteString("\n")
	return atomicWriteFile(filepath.Join(m.rootDir, path), []byte(b.String()), 0o644)
}

func (m *MarkdownMemoryStore) readMemoryFileLocked(path string) (storage.Memory, error) {
	if !safeRelativeMemoryPath(path) {
		return storage.Memory{}, fmt.Errorf("unsafe memory path: %s", path)
	}
	data, err := os.ReadFile(filepath.Join(m.rootDir, path))
	if err != nil {
		return storage.Memory{}, err
	}
	meta, body := splitFrontMatter(string(data))
	return storage.Memory{ID: strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)), Type: meta["metadata.type"], Name: meta["name"], Description: meta["description"], Body: strings.TrimSpace(body)}, nil
}

func (m *MarkdownMemoryStore) sessionPath(sessionID string) string {
	return filepath.Join(m.sessionsDir, sanitizeName(sessionID)+".md")
}

func (m *MarkdownMemoryStore) writeSessionMemoryLocked(sessionID string, items []storage.ConversationMemory) error {
	sessionID = sanitizeName(sessionID)
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("name: current-session\n")
	b.WriteString("description: 当前会话主干信息\n")
	b.WriteString("metadata:\n")
	b.WriteString("  node_type: session_memory\n")
	b.WriteString("  project: " + yamlScalar(m.projectName) + "\n")
	b.WriteString("  session_id: " + yamlScalar(sessionID) + "\n")
	b.WriteString("  originSessionId: " + yamlScalar(sessionID) + "\n")
	b.WriteString("---\n\n")
	for _, item := range items {
		b.WriteString("## " + strings.TrimSpace(item.Name) + "\n\n")
		if strings.TrimSpace(item.Description) != "" {
			b.WriteString(strings.TrimSpace(item.Description) + "\n\n")
		}
		b.WriteString(strings.TrimSpace(item.Body) + "\n\n")
	}
	return atomicWriteFile(m.sessionPath(sessionID), []byte(b.String()), 0o644)
}

func parseConversationMemoryBody(sessionID, body string) []storage.ConversationMemory {
	sections := strings.Split(body, "\n## ")
	var result []storage.ConversationMemory
	for _, section := range sections {
		section = strings.TrimSpace(strings.TrimPrefix(section, "## "))
		if section == "" {
			continue
		}
		parts := strings.SplitN(section, "\n", 2)
		name := strings.TrimSpace(parts[0])
		content := ""
		if len(parts) > 1 {
			content = strings.TrimSpace(parts[1])
		}
		result = append(result, storage.ConversationMemory{ConversationID: sessionID, Name: name, Body: content})
	}
	return result
}

func splitFrontMatter(raw string) (map[string]string, string) {
	meta := make(map[string]string)
	if !strings.HasPrefix(raw, "---\n") {
		return meta, raw
	}
	end := strings.Index(raw[4:], "\n---")
	if end < 0 {
		return meta, raw
	}
	front := raw[4 : 4+end]
	body := raw[4+end+5:]
	var prefix string
	scanner := bufio.NewScanner(strings.NewReader(front))
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if strings.HasSuffix(trimmed, ":") {
			prefix = strings.TrimSuffix(trimmed, ":") + "."
			continue
		}
		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		if strings.HasPrefix(line, "  ") {
			key = prefix + key
		}
		meta[key] = strings.Trim(strings.TrimSpace(parts[1]), `"`)
	}
	return meta, body
}

func filterMemories(items []storage.Memory, memType string) []storage.Memory {
	var result []storage.Memory
	for _, item := range items {
		if item.Type == memType {
			result = append(result, item)
		}
	}
	return result
}

func sanitizeName(s string) string {
	var b strings.Builder
	lastDash := false
	for _, r := range strings.TrimSpace(s) {
		allowed := r == '.' || r == '_' || r == '-' || unicode.IsLetter(r) || unicode.IsDigit(r)
		if allowed {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-._")
}

func safeRelativeMemoryPath(path string) bool {
	if strings.TrimSpace(path) == "" || filepath.IsAbs(path) {
		return false
	}
	clean := filepath.Clean(path)
	return clean == filepath.Base(clean) && !strings.HasPrefix(clean, "..") && strings.HasSuffix(clean, ".md") && clean != "memory.md"
}

func yamlScalar(s string) string {
	s = strings.ReplaceAll(strings.TrimSpace(s), "\"", "\\\"")
	return "\"" + s + "\""
}

func atomicWriteFile(path string, data []byte, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpName, perm); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package local\n\nimport (\n\t\"bufio\"\n\t\"context\"\n\t\"fmt\"\n\t\"os\"\n\t\"path/filepath\"\n\t\"sort\"\n\t\"strings\"\n\t\"sync\"\n\t\"time\"\n\t\"unicode\"\n\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/idgen\"\n)\n\nconst memoryIndexHeader = \"# Memory Index\\n\"\n\ntype MarkdownMemoryStore struct {\n\tmu          sync.Mutex\n\trootDir     string\n\tindexPath   string\n\tsessionsDir string\n\tprojectName string\n}\n\nfunc NewMarkdownMemoryStore(workspaceRoot string) (*MarkdownMemoryStore, error) {\n\troot := filepath.Join(strings.TrimSpace(worksp"}
```

## 2026-06-13T09:17:29Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_07_aEtdiTAlDOEsTIlu7VzX0098

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/persisted_output_files.go"}
```

### Result

```text
package local

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"nano_cc/internal/agent/storage"
)

type persistedOutputMetadata struct {
	Version        int       `json:"version"`
	ID             string    `json:"id"`
	SessionID      string    `json:"session_id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	MessageID      string    `json:"message_id"`
	ToolCallID     string    `json:"tool_call_id"`
	Kind           string    `json:"kind"`
	Strategy       string    `json:"strategy"`
	OriginalBytes  int       `json:"original_bytes"`
	ContentSHA256  string    `json:"content_sha256"`
	Preview        string    `json:"preview"`
	ContentFile    string    `json:"content_file"`
	CreatedAt      time.Time `json:"created_at"`
}

func persistOutputToWorkspace(ctx context.Context, workspaceRoot, sessionID string, output storage.PersistedOutput) (storage.PersistedOutput, error) {
	if err := ctx.Err(); err != nil {
		return output, err
	}
	if !validSessionID(sessionID) {
		return output, fmt.Errorf("invalid session_id: %q", sessionID)
	}
	if !validStoredOutputID(output.ID) {
		return output, fmt.Errorf("invalid persisted output id: %q", output.ID)
	}
	if output.CreatedAt.IsZero() {
		output.CreatedAt = time.Now()
	}
	if output.OriginalBytes == 0 {
		output.OriginalBytes = len([]byte(output.Content))
	}
	sum := sha256.Sum256([]byte(output.Content))
	computedSHA := hex.EncodeToString(sum[:])
	if output.ContentSHA256 == "" {
		output.ContentSHA256 = computedSHA
	}

	contentFile := persistedOutputContentFileName(sessionID, output.ID)
	metadataFile := persistedOutputMetadataFileName(sessionID, output.ID)
	dir := persistedOutputDir(workspaceRoot)
	metadata := persistedOutputMetadata{
		Version:        1,
		ID:             output.ID,
		SessionID:      sessionID,
		ConversationID: output.ConversationID,
		UserID:         output.UserID,
		MessageID:      output.MessageID,
		ToolCallID:     output.ToolCallID,
		Kind:           output.Kind,
		Strategy:       output.Strategy,
		OriginalBytes:  output.OriginalBytes,
		ContentSHA256:  output.ContentSHA256,
		Preview:        output.Preview,
		ContentFile:    contentFile,
		CreatedAt:      output.CreatedAt,
	}
	metadataBytes, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return output, err
	}
	if err := atomicWriteFile(filepath.Join(dir, contentFile), []byte(output.Content), 0o644); err != nil {
		return output, err
	}
	if err := atomicWriteFile(filepath.Join(dir, metadataFile), metadataBytes, 0o644); err != nil {
		return output, err
	}
	return output, nil
}

func loadPersistedOutputFromWorkspace(ctx context.Context, workspaceRoot, sessionID, id string) (storage.PersistedOutput, error) {
	if err := ctx.Err(); err != nil {
		return storage.PersistedOutput{}, err
	}
	if !validSessionID(sessionID) {
		return storage.PersistedOutput{}, fmt.Errorf("invalid session_id: %q", sessionID)
	}
	if !validStoredOutputID(id) {
		return storage.PersistedOutput{}, fmt.Errorf("invalid persisted output id: %q", id)
	}
	dir := persistedOutputDir(workspaceRoot)
	metadataBytes, err := os.ReadFile(filepath.Join(dir, persistedOutputMetadataFileName(sessionID, id)))
	if err != nil {
		return storage.PersistedOutput{}, err
	}
	var metadata persistedOutputMetadata
	if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
		return storage.PersistedOutput{}, err
	}
	if metadata.ID != id || metadata.SessionID != sessionID || metadata.Kind != "tool_result" {
		return storage.PersistedOutput{}, fmt.Errorf("persisted output metadata mismatch for %q", id)
	}
	if filepath.Base(metadata.ContentFile) != metadata.ContentFile || metadata.ContentFile == "" {
		return storage.PersistedOutput{}, fmt.Errorf("invalid persisted output content file: %q", metadata.ContentFile)
	}
	contentBytes, err := os.ReadFile(filepath.Join(dir, metadata.ContentFile))
	if err != nil {
		return storage.PersistedOutput{}, err
	}
	sum := sha256.Sum256(contentBytes)
	actualSHA := hex.EncodeToString(sum[:])
	if metadata.ContentSHA256 != "" && actualSHA != metadata.ContentSHA256 {
		return storage.PersistedOutput{}, fmt.Errorf("persisted output %q sha256 mismatch", id)
	}
	return storage.PersistedOutput{
		ID:             metadata.ID,
		ConversationID: metadata.ConversationID,
		UserID:         metadata.UserID,
		MessageID:      metadata.MessageID,
		ToolCallID:     metadata.ToolCallID,
		Kind:           metadata.Kind,
		Strategy:       metadata.Strategy,
		OriginalBytes:  metadata.OriginalBytes,
		ContentSHA256:  metadata.ContentSHA256,
		Content:        string(contentBytes),
		Preview:        metadata.Preview,
		CreatedAt:      metadata.CreatedAt,
	}, nil
}

func findPersistedOutputByMessageHashInWorkspace(ctx context.Context, workspaceRoot, sessionID, conversationID, userID, messageID, toolCallID, strategy, contentSHA256 string) (storage.PersistedOutput, error) {
	if err := ctx.Err(); err != nil {
		return storage.PersistedOutput{}, err
	}
	if !validSessionID(sessionID) {
		return storage.PersistedOutput{}, fmt.Errorf("invalid session_id: %q", sessionID)
	}
	dir := persistedOutputDir(workspaceRoot)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return storage.PersistedOutput{}, sql.ErrNoRows
		}
		return storage.PersistedOutput{}, err
	}
	prefix := sessionID + "-"
	for _, entry := range entries {
		if ctx.Err() != nil {
			return storage.PersistedOutput{}, ctx.Err()
		}
		name := entry.Name()
		if entry.IsDir() || !strings.HasPrefix(name, prefix) || !strings.HasSuffix(name, ".json") {
			continue
		}
		metadataBytes, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return storage.PersistedOutput{}, err
		}
		var metadata persistedOutputMetadata
		if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
			return storage.PersistedOutput{}, err
		}
		if metadata.SessionID != sessionID || metadata.ConversationID != conversationID || metadata.UserID != userID || metadata.MessageID != messageID || metadata.ToolCallID != toolCallID || metadata.Strategy != strategy || metadata.ContentSHA256 != contentSHA256 {
			continue
		}
		return loadPersistedOutputFromWorkspace(ctx, workspaceRoot, sessionID, metadata.ID)
	}
	return storage.PersistedOutput{}, sql.ErrNoRows
}

func persistedOutputDir(workspaceRoot string) string {
	return filepath.Join(workspaceRoot, "task_outputs", "tool-results")
}

func persistedOutputContentFileName(sessionID, id string) string {
	return sessionID + "-" + id + ".txt"
}

func persistedOutputMetadataFileName(sessionID, id string) string {
	return sessionID + "-" + id + ".json"
}

func validStoredOutputID(id string) bool {
	id = strings.TrimSpace(id)
	if id == "" || id == "." || id == ".." || strings.Contains(id, "..") {
		return false
	}
	for _, r := range id {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == '.' {
			continue
		}
		return false
	}
	return true
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package local\n\nimport (\n\t\"context\"\n\t\"crypto/sha256\"\n\t\"database/sql\"\n\t\"encoding/hex\"\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"os\"\n\t\"path/filepath\"\n\t\"strings\"\n\t\"time\"\n\t\"unicode\"\n\n\t\"nano_cc/internal/agent/storage\"\n)\n\ntype persistedOutputMetadata struct {\n\tVersion        int       `json:\"version\"`\n\tID             string    `json:\"id\"`\n\tSessionID      string    `json:\"session_id\"`\n\tConversationID string    `json:\"conversation_id\"`\n\tUserID         string    `json:\"user_id\"`\n\tMessageID      string    `json:\"message_i"}
```

## 2026-06-13T09:17:29Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_08_OQFkbnSvYL73W4a5lYmx9849

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/session_history_store.go"}
```

### Result

```text
package local

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/config"
)

type SessionHistoryStore struct {
	dir string
}

type SessionHistoryDocument struct {
	Version        int               `json:"version"`
	SessionID      string            `json:"session_id"`
	ConversationID string            `json:"conversation_id"`
	UserID         string            `json:"user_id"`
	WorkspaceRoot  string            `json:"workspace_root"`
	Title          string            `json:"title"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	Messages       []storage.Message `json:"messages"`
}

func NewSessionHistoryStore() (*SessionHistoryStore, error) {
	dir, err := config.LinkSessionDir()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create session dir: %w", err)
	}
	return &SessionHistoryStore{dir: dir}, nil
}

func (s *SessionHistoryStore) SaveHistory(ctx context.Context, workspaceRoot string, conversation storage.Conversation, messages []storage.Message) error {
	return s.save(ctx, workspaceRoot, conversation, messages, "history")
}

func (s *SessionHistoryStore) LoadHistory(ctx context.Context, sessionID string) (SessionHistoryDocument, error) {
	return s.load(ctx, sessionID, "history")
}

func (s *SessionHistoryStore) SaveModelHistory(ctx context.Context, workspaceRoot string, conversation storage.Conversation, messages []storage.Message) error {
	return s.save(ctx, workspaceRoot, conversation, messages, "model_history")
}

func (s *SessionHistoryStore) LoadModelHistory(ctx context.Context, sessionID string) (SessionHistoryDocument, error) {
	return s.load(ctx, sessionID, "model_history")
}

func (s *SessionHistoryStore) ListSessions(ctx context.Context, workspaceRoot string) ([]storage.ResumableSession, error) {
	currentWorkspace, err := normalizeWorkspacePath(workspaceRoot)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var result []storage.ResumableSession
	for _, entry := range entries {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if !entry.IsDir() || !validSessionID(entry.Name()) {
			continue
		}
		doc, err := s.LoadHistory(ctx, entry.Name())
		if err != nil {
			if os.IsNotExist(err) || err == sql.ErrNoRows {
				continue
			}
			return nil, err
		}
		docWorkspace, err := normalizeWorkspacePath(doc.WorkspaceRoot)
		if err != nil || docWorkspace != currentWorkspace {
			continue
		}
		result = append(result, storage.ResumableSession{
			SessionID:      doc.SessionID,
			ConversationID: doc.ConversationID,
			WorkspaceRoot:  doc.WorkspaceRoot,
			Title:          doc.Title,
			MessageCount:   len(doc.Messages),
			CreatedAt:      doc.CreatedAt,
			UpdatedAt:      doc.UpdatedAt,
		})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].UpdatedAt.After(result[j].UpdatedAt) })
	return result, nil
}

func (s *SessionHistoryStore) save(ctx context.Context, workspaceRoot string, conversation storage.Conversation, messages []storage.Message, name string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !validSessionID(conversation.SessionID) {
		return fmt.Errorf("invalid session_id: %q", conversation.SessionID)
	}
	now := time.Now()
	createdAt := conversation.CreatedAt
	if createdAt.IsZero() {
		createdAt = now
	}
	updatedAt := conversation.UpdatedAt
	if updatedAt.IsZero() || updatedAt.Before(now) {
		updatedAt = now
	}
	doc := SessionHistoryDocument{
		Version:        1,
		SessionID:      conversation.SessionID,
		ConversationID: conversation.ID,
		UserID:         conversation.UserID,
		WorkspaceRoot:  workspaceRoot,
		Title:          conversation.Title,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		Messages:       cloneMessages(messages),
	}
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	return atomicWriteFile(s.path(conversation.SessionID, name), data, 0o644)
}

func (s *SessionHistoryStore) load(ctx context.Context, sessionID, name string) (SessionHistoryDocument, error) {
	if err := ctx.Err(); err != nil {
		return SessionHistoryDocument{}, err
	}
	if !validSessionID(sessionID) {
		return SessionHistoryDocument{}, fmt.Errorf("invalid session_id: %q", sessionID)
	}
	data, err := os.ReadFile(s.path(sessionID, name))
	if err != nil {
		return SessionHistoryDocument{}, err
	}
	var doc SessionHistoryDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return SessionHistoryDocument{}, fmt.Errorf("parse session %s %s: %w", sessionID, name, err)
	}
	if doc.SessionID == "" {
		doc.SessionID = sessionID
	}
	if doc.SessionID != sessionID {
		return SessionHistoryDocument{}, fmt.Errorf("session_id mismatch: path=%s document=%s", sessionID, doc.SessionID)
	}
	return doc, nil
}

func (s *SessionHistoryStore) path(sessionID, name string) string {
	return filepath.Join(s.dir, sessionID, name)
}

func validSessionID(sessionID string) bool {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" || sessionID == "." || sessionID == ".." || strings.Contains(sessionID, "..") {
		return false
	}
	for _, r := range sessionID {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == '.' {
			continue
		}
		return false
	}
	return true
}

func normalizeWorkspacePath(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		path = "."
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	clean := filepath.Clean(abs)
	if evaluated, err := filepath.EvalSymlinks(clean); err == nil {
		clean = evaluated
	}
	return clean, nil
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package local\n\nimport (\n\t\"context\"\n\t\"database/sql\"\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"os\"\n\t\"path/filepath\"\n\t\"sort\"\n\t\"strings\"\n\t\"time\"\n\t\"unicode\"\n\n\t\"nano_cc/internal/agent/storage\"\n\t\"nano_cc/internal/config\"\n)\n\ntype SessionHistoryStore struct {\n\tdir string\n}\n\ntype SessionHistoryDocument struct {\n\tVersion        int               `json:\"version\"`\n\tSessionID      string            `json:\"session_id\"`\n\tConversationID string            `json:\"conversation_id\"`\n\tUserID         string            `json:\"user_id\""}
```

## 2026-06-13T09:17:29Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_09_S8HEfugverMRbH3FqPcB1783

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/local/tool_result_log.go"}
```

### Result

````text
package local

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"nano_cc/internal/agent/storage"
)

func (s *Store) AppendToolResultLog(ctx context.Context, entry storage.ToolResultLogEntry) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.workspaceRoot == "" {
		return nil
	}
	conv := s.conversations[entry.ConversationID]
	sessionID := entry.SessionID
	if sessionID == "" {
		sessionID = conv.SessionID
	}
	if !validSessionID(sessionID) {
		return fmt.Errorf("invalid session_id: %q", sessionID)
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now()
	}
	path := filepath.Join(s.workspaceRoot, "task_outputs", sessionID, "tools.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	unlock, err := lockFile(path + ".lock")
	if err != nil {
		return err
	}
	defer unlock()
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(renderToolResultLogEntry(sessionID, entry))
	return err
}

func lockFile(path string) (func(), error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		_ = file.Close()
		return nil, err
	}
	return func() {
		_ = syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
		_ = file.Close()
	}, nil
}

func renderToolResultLogEntry(sessionID string, entry storage.ToolResultLogEntry) string {
	var b strings.Builder
	b.WriteString("## ")
	b.WriteString(entry.CreatedAt.UTC().Format(time.RFC3339))
	b.WriteString(" · ")
	b.WriteString(entry.ToolName)
	b.WriteString(" · ")
	b.WriteString(entry.Status)
	b.WriteString("\n\n")
	b.WriteString("- conversation_id: ")
	b.WriteString(entry.ConversationID)
	b.WriteString("\n- session_id: ")
	b.WriteString(sessionID)
	b.WriteString("\n- tool_call_id: ")
	b.WriteString(entry.ToolCallID)
	b.WriteString("\n\n")
	b.WriteString("### Arguments\n\n")
	b.WriteString(markdownFence("json", entry.RawArgs))
	b.WriteString("\n\n")
	b.WriteString("### Result\n\n")
	b.WriteString(markdownFence("text", entry.Result))
	if strings.TrimSpace(entry.AuditSummary) != "" {
		b.WriteString("\n\n### Audit\n\n")
		b.WriteString(markdownFence("json", entry.AuditSummary))
	}
	b.WriteString("\n\n")
	return b.String()
}

func markdownFence(language, body string) string {
	fence := strings.Repeat("`", longestBacktickRun(body)+1)
	if len(fence) < 3 {
		fence = "```"
	}
	return fence + language + "\n" + body + "\n" + fence
}

func longestBacktickRun(body string) int {
	longest := 0
	current := 0
	for _, r := range body {
		if r == '`' {
			current++
			if current > longest {
				longest = current
			}
			continue
		}
		current = 0
	}
	return longest
}

````

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package local\n\nimport (\n\t\"context\"\n\t\"fmt\"\n\t\"os\"\n\t\"path/filepath\"\n\t\"strings\"\n\t\"syscall\"\n\t\"time\"\n\n\t\"nano_cc/internal/agent/storage\"\n)\n\nfunc (s *Store) AppendToolResultLog(ctx context.Context, entry storage.ToolResultLogEntry) error {\n\tif err := ctx.Err(); err != nil {\n\t\treturn err\n\t}\n\ts.mu.Lock()\n\tdefer s.mu.Unlock()\n\tif s.workspaceRoot == \"\" {\n\t\treturn nil\n\t}\n\tconv := s.conversations[entry.ConversationID]\n\tsessionID := entry.SessionID\n\tif sessionID == \"\" {\n\t\tsessionID = conv.SessionID\n\t}\n\tif !val"}
```

## 2026-06-13T09:17:32Z · read_persisted_output · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_00_sij6jKjQDkDRizUFxqWA9047

### Arguments

```json
{"id": "po_9e84cc723d4c231d302d6445", "offset": 2000, "limit": 20000}
```

### Result

```text
{"content":"on:\"role\"`\n\tContent          string            `json:\"content\"`\n\tReasoningContent string            `json:\"reasoning_content,omitempty\"`\n\tToolCallID       string            `json:\"tool_call_id,omitempty\"`\n\tToolCalls        []MessageToolCall `json:\"tool_calls,omitempty\"`\n\tMeta             *MessageMeta      `json:\"meta,omitempty\"`\n\tCreatedAt        time.Time         `json:\"created_at\"`\n}\n\n// MessageMeta 记录助手回复的元信息，仅对最终 assistant 消息填充。\ntype MessageMeta struct {\n\tToolCallCount int `json:\"tool_call_count\"`          // 调用工具次数（0 也需序列化，保证历史展示一致）\n\tContextTokens int `json:\"context_tokens,omitempty\"` // 当前上下文估算 token\n\tContextBudget int `json:\"context_budget,omitempty\"` // 上下文预算（用于算占比）\n}\n\ntype MessageToolCall struct {\n\tID       string              `json:\"id\"`\n\tType     string              `json:\"type\"`\n\tFunction MessageFunctionCall `json:\"function\"`\n}\n\ntype MessageFunctionCall struct {\n\tName      string `json:\"name\"`\n\tArguments string `json:\"arguments\"`\n}\n\ntype ToolCall struct {\n\tID             string    `json:\"id\"`\n\tConversationID string    `json:\"conversation_id\"`\n\tUserID         string    `json:\"user_id\"`\n\tToolName       string    `json:\"tool_name\"`\n\tStatus         string    `json:\"status\"`\n\tSummary        string    `json:\"summary\"`\n\tCreatedAt      time.Time `json:\"created_at\"`\n}\n\ntype PersistedOutput struct {\n\tID             string    `json:\"id\"`\n\tConversationID string    `json:\"conversation_id\"`\n\tUserID         string    `json:\"user_id\"`\n\tMessageID      string    `json:\"message_id\"`\n\tToolCallID     string    `json:\"tool_call_id\"`\n\tKind           string    `json:\"kind\"`\n\tStrategy       string    `json:\"strategy\"`\n\tOriginalBytes  int       `json:\"original_bytes\"`\n\tContentSHA256  string    `json:\"content_sha256\"`\n\tContent        string    `json:\"content\"`\n\tPreview        string    `json:\"preview\"`\n\tCreatedAt      time.Time `json:\"created_at\"`\n}\n\ntype ToolResultLogEntry struct {\n\tConversationID string    `json:\"conversation_id\"`\n\tSessionID      string    `json:\"session_id,omitempty\"`\n\tUserID         string    `json:\"user_id\"`\n\tToolCallID     string    `json:\"tool_call_id\"`\n\tToolName       string    `json:\"tool_name\"`\n\tRawArgs        string    `json:\"raw_args\"`\n\tStatus         string    `json:\"status\"`\n\tResult         string    `json:\"result\"`\n\tAuditSummary   string    `json:\"audit_summary,omitempty\"`\n\tCreatedAt      time.Time `json:\"created_at\"`\n}\n\ntype ContextSummary struct {\n\tID                    string    `json:\"id\"`\n\tConversationID        string    `json:\"conversation_id\"`\n\tUserID                string    `json:\"user_id\"`\n\tSourceHistorySHA256   string    `json:\"source_history_sha256\"`\n\tStrategy              string    `json:\"strategy\"`\n\tEstimatedTokensBefore int       `json:\"estimated_tokens_before\"`\n\tEstimatedTokensAfter  int       `json:\"estimated_tokens_after\"`\n\tSummary               string    `json:\"summary\"`\n\tCreatedAt             time.Time `json:\"created_at\"`\n}\n\ntype SubagentMessage struct {\n\tID               string            `json:\"id\"`\n\tRunID            string            `json:\"run_id\"`\n\tParentToolCallID string            `json:\"parent_tool_call_id\"`\n\tConversationID   string            `json:\"conversation_id\"`\n\tUserID           string            `json:\"user_id\"`\n\tSequenceNo       int               `json:\"sequence_no\"`\n\tRole             string            `json:\"role\"`\n\tContent          string            `json:\"content\"`\n\tReasoningContent string            `json:\"reasoning_content,omitempty\"`\n\tToolCallID       string            `json:\"tool_call_id,omitempty\"`\n\tToolCalls        []MessageToolCall `json:\"tool_calls,omitempty\"`\n\tCreatedAt        time.Time         `json:\"created_at\"`\n}\n\ntype Memory struct {\n\tID          string    `json:\"id\"`\n\tUserID      string    `json:\"user_id\"`\n\tType        string    `json:\"type\"`\n\tName        string    `json:\"name\"`\n\tDescription string    `json:\"description\"`\n\tBody        string    `json:\"body\"`\n\tCreatedAt   time.Time `json:\"created_at\"`\n\tUpdatedAt   time.Time `json:\"updated_at\"`\n}\n\n// ConversationMemory 是单个会话维度、随每轮对话增量维护的“当前会话主干信息”\n// 条目。它不注入 system prompt，仅在上下文压缩触发全量摘要时作为替代品使用。\ntype ConversationMemory struct {\n\tID             string    `json:\"id\"`\n\tConversationID string    `json:\"conversation_id\"`\n\tUserID         string    `json:\"user_id\"`\n\tName           string    `json:\"name\"`\n\tDescription    string    `json:\"description\"`\n\tBody           string    `json:\"body\"`\n\tCreatedAt      time.Time `json:\"created_at\"`\n\tUpdatedAt      time.Time `json:\"updated_at\"`\n}\n","content_sha256":"6efed0813606063b9e48e518fb33d8856dc5d5fe166ad9267307173d86357c98","has_more":false,"id":"po_9e84cc723d4c231d302d6445","kind":"tool_result","limit":20000,"next_offset":6474,"offset":2000,"original_bytes":6826,"returned_chars":4474,"tool_call_id":"call_00_Y8fl1PQPI7Su0a4739mT1227","total_chars":6474}
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"{\"content\":\"on:\\\"role\\\"`\\n\\tContent          string            `json:\\\"content\\\"`\\n\\tReasoningContent string            `json:\\\"reasoning_content,omitempty\\\"`\\n\\tToolCallID       string            `json:\\\"tool_call_id,omitempty\\\"`\\n\\tToolCalls        []MessageToolCall `json:\\\"tool_calls,omitempty\\\"`\\n\\tMeta             *MessageMeta      `json:\\\"meta,omitempty\\\"`\\n\\tCreatedAt        time.Time         `json:\\\"created_at\\\"`\\n}\\n\\n// MessageMeta 记录助手回复的元信息，仅对最终 assistan"}
```

## 2026-06-13T09:17:32Z · read_persisted_output · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_01_jdbjXPuS0kKjjklHWjja0508

### Arguments

```json
{"id": "po_6f296a3f0f3c75fb6bed6f5d", "offset": 2000, "limit": 20000}
```

### Result

```text
{"content":"on:\"role\"`\n\tContent          string            `json:\"content\"`\n\tReasoningContent string            `json:\"reasoning_content,omitempty\"`\n\tToolCallID       string            `json:\"tool_call_id,omitempty\"`\n\tToolCalls        []MessageToolCall `json:\"tool_calls,omitempty\"`\n\tMeta             *MessageMeta      `json:\"meta,omitempty\"`\n\tCreatedAt        time.Time         `json:\"created_at\"`\n}\n\n// MessageMeta 记录助手回复的元信息，仅对最终 assistant 消息填充。\ntype MessageMeta struct {\n\tToolCallCount int `json:\"tool_call_count\"`          // 调用工具次数（0 也需序列化，保证历史展示一致）\n\tContextTokens int `json:\"context_tokens,omitempty\"` // 当前上下文估算 token\n\tContextBudget int `json:\"context_budget,omitempty\"` // 上下文预算（用于算占比）\n}\n\ntype MessageToolCall struct {\n\tID       string              `json:\"id\"`\n\tType     string              `json:\"type\"`\n\tFunction MessageFunctionCall `json:\"function\"`\n}\n\ntype MessageFunctionCall struct {\n\tName      string `json:\"name\"`\n\tArguments string `json:\"arguments\"`\n}\n\ntype ToolCall struct {\n\tID             string    `json:\"id\"`\n\tConversationID string    `json:\"conversation_id\"`\n\tUserID         string    `json:\"user_id\"`\n\tToolName       string    `json:\"tool_name\"`\n\tStatus         string    `json:\"status\"`\n\tSummary        string    `json:\"summary\"`\n\tCreatedAt      time.Time `json:\"created_at\"`\n}\n\ntype PersistedOutput struct {\n\tID             string    `json:\"id\"`\n\tConversationID string    `json:\"conversation_id\"`\n\tUserID         string    `json:\"user_id\"`\n\tMessageID      string    `json:\"message_id\"`\n\tToolCallID     string    `json:\"tool_call_id\"`\n\tKind           string    `json:\"kind\"`\n\tStrategy       string    `json:\"strategy\"`\n\tOriginalBytes  int       `json:\"original_bytes\"`\n\tContentSHA256  string    `json:\"content_sha256\"`\n\tContent        string    `json:\"content\"`\n\tPreview        string    `json:\"preview\"`\n\tCreatedAt      time.Time `json:\"created_at\"`\n}\n\ntype ToolResultLogEntry struct {\n\tConversationID string    `json:\"conversation_id\"`\n\tSessionID      string    `json:\"session_id,omitempty\"`\n\tUserID         string    `json:\"user_id\"`\n\tToolCallID     string    `json:\"tool_call_id\"`\n\tToolName       string    `json:\"tool_name\"`\n\tRawArgs        string    `json:\"raw_args\"`\n\tStatus         string    `json:\"status\"`\n\tResult         string    `json:\"result\"`\n\tAuditSummary   string    `json:\"audit_summary,omitempty\"`\n\tCreatedAt      time.Time `json:\"created_at\"`\n}\n\ntype ContextSummary struct {\n\tID                    string    `json:\"id\"`\n\tConversationID        string    `json:\"conversation_id\"`\n\tUserID                string    `json:\"user_id\"`\n\tSourceHistorySHA256   string    `json:\"source_history_sha256\"`\n\tStrategy              string    `json:\"strategy\"`\n\tEstimatedTokensBefore int       `json:\"estimated_tokens_before\"`\n\tEstimatedTokensAfter  int       `json:\"estimated_tokens_after\"`\n\tSummary               string    `json:\"summary\"`\n\tCreatedAt             time.Time `json:\"created_at\"`\n}\n\ntype SubagentMessage struct {\n\tID               string            `json:\"id\"`\n\tRunID            string            `json:\"run_id\"`\n\tParentToolCallID string            `json:\"parent_tool_call_id\"`\n\tConversationID   string            `json:\"conversation_id\"`\n\tUserID           string            `json:\"user_id\"`\n\tSequenceNo       int               `json:\"sequence_no\"`\n\tRole             string            `json:\"role\"`\n\tContent          string            `json:\"content\"`\n\tReasoningContent string            `json:\"reasoning_content,omitempty\"`\n\tToolCallID       string            `json:\"tool_call_id,omitempty\"`\n\tToolCalls        []MessageToolCall `json:\"tool_calls,omitempty\"`\n\tCreatedAt        time.Time         `json:\"created_at\"`\n}\n\ntype Memory struct {\n\tID          string    `json:\"id\"`\n\tUserID      string    `json:\"user_id\"`\n\tType        string    `json:\"type\"`\n\tName        string    `json:\"name\"`\n\tDescription string    `json:\"description\"`\n\tBody        string    `json:\"body\"`\n\tCreatedAt   time.Time `json:\"created_at\"`\n\tUpdatedAt   time.Time `json:\"updated_at\"`\n}\n\n// ConversationMemory 是单个会话维度、随每轮对话增量维护的“当前会话主干信息”\n// 条目。它不注入 system prompt，仅在上下文压缩触发全量摘要时作为替代品使用。\ntype ConversationMemory struct {\n\tID             string    `json:\"id\"`\n\tConversationID string    `json:\"conversation_id\"`\n\tUserID         string    `json:\"user_id\"`\n\tName           string    `json:\"name\"`\n\tDescription    string    `json:\"description\"`\n\tBody           string    `json:\"body\"`\n\tCreatedAt      time.Time `json:\"created_at\"`\n\tUpdatedAt      time.Time `json:\"updated_at\"`\n}\n","content_sha256":"6efed0813606063b9e48e518fb33d8856dc5d5fe166ad9267307173d86357c98","has_more":false,"id":"po_6f296a3f0f3c75fb6bed6f5d","kind":"tool_result","limit":20000,"next_offset":6474,"offset":2000,"original_bytes":6826,"returned_chars":4474,"tool_call_id":"call_00_RWsw3Qx66GfruvrgMR8S2511","total_chars":6474}
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"{\"content\":\"on:\\\"role\\\"`\\n\\tContent          string            `json:\\\"content\\\"`\\n\\tReasoningContent string            `json:\\\"reasoning_content,omitempty\\\"`\\n\\tToolCallID       string            `json:\\\"tool_call_id,omitempty\\\"`\\n\\tToolCalls        []MessageToolCall `json:\\\"tool_calls,omitempty\\\"`\\n\\tMeta             *MessageMeta      `json:\\\"meta,omitempty\\\"`\\n\\tCreatedAt        time.Time         `json:\\\"created_at\\\"`\\n}\\n\\n// MessageMeta 记录助手回复的元信息，仅对最终 assistan"}
```

## 2026-06-13T09:17:32Z · read_file · success

- conversation_id: conv_a1ebf86b6d0ab87bbaa61928
- session_id: 99f6c277-f695-47f9-8a6e-8c91573f6083
- tool_call_id: call_02_8jnqOGKr9dt9E2GjS20Z3114

### Arguments

```json
{"path": "/Users/bytedance/golang_pro/nano_cc/go-agent/internal/config/config.go", "limit": 200}
```

### Result

```text
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"-"`
	ModelID string `json:"model_id"`
}

type fileConfig struct {
	Config
	AppHome                            string `json:"app_home"`
	BuiltinSkillsDir                   string `json:"builtin_skills_dir"`
	CommandBinDir                      string `json:"command_bin_dir"`
	CommandScriptDir                   string `json:"command_script_dir"`
	SystemPromptPath                   string `json:"system_prompt_path"`
	WorkspaceRoot                      string `json:"workspace_root"`
	AllowedTools                       string `json:"allowed_tools"`
	BashAllowOutsideWorkspace          bool   `json:"bash_allow_outside_workspace"`
	BashAllowDangerousCommands         bool   `json:"bash_allow_dangerous_commands"`
	ConversationLockTTLSeconds         int    `json:"conversation_lock_ttl_seconds"`
	MemoryWorkTimeoutSeconds           int    `json:"memory_work_timeout_seconds"`
	ConversationLockWaitTimeoutSeconds int    `json:"conversation_lock_wait_timeout_seconds"`
}

type AppConfig struct {
	LLM                         Config
	AppHome                     string
	BuiltinSkillsDir            string
	CommandBinDir               string
	CommandScriptDir            string
	SystemPromptPath            string
	WorkspaceRoot               string
	LogsDir                     string
	AllowedTools                []string
	BashAllowOutsideWorkspace   bool
	BashAllowDangerousCommands  bool
	ConversationLockTTL         time.Duration
	MemoryWorkTimeout           time.Duration
	ConversationLockWaitTimeout time.Duration
}

func loadConfigFile() (fileConfig, error) {
	data, err := os.ReadFile(configFilePath())
	if err != nil {
		return fileConfig{}, err
	}

	var cfg fileConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fileConfig{}, fmt.Errorf("failed to parse workspace config.json: %w", err)
	}
	return cfg, nil
}

func configFilePath() string {
	return "config.json"
}

func loadLLMConfig(fileCfg fileConfig) (Config, error) {
	cfg := Config{
		BaseURL: strings.TrimSpace(fileCfg.Config.BaseURL),
		APIKey:  strings.TrimSpace(getenv("OPENAI_API_KEY")),
		ModelID: strings.TrimSpace(fileCfg.Config.ModelID),
	}

	if cfg.BaseURL == "" || cfg.ModelID == "" {
		return Config{}, fmt.Errorf("missing LLM config; set base_url, model_id in config.json")
	}
	if cfg.APIKey == "" {
		return Config{}, fmt.Errorf("missing LLM api key; set OPENAI_API_KEY environment variable")
	}

	return cfg, nil
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package config\n\nimport (\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"os\"\n\t\"strings\"\n\t\"time\"\n)\n\ntype Config struct {\n\tBaseURL string `json:\"base_url\"`\n\tAPIKey  string `json:\"-\"`\n\tModelID string `json:\"model_id\"`\n}\n\ntype fileConfig struct {\n\tConfig\n\tAppHome                            string `json:\"app_home\"`\n\tBuiltinSkillsDir                   string `json:\"builtin_skills_dir\"`\n\tCommandBinDir                      string `json:\"command_bin_dir\"`\n\tCommandScriptDir                   string `json:\"command_script_dir\"`\n"}
```

