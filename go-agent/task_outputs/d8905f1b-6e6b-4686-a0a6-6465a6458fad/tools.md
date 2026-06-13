## 2026-06-13T09:05:18Z · bash · rejected

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_00_nvemiLtsLlnsD0OLkQKz2273

### Arguments

```json
{"command": "ls -1 ~/.link/skills 2>/dev/null; echo \"---\"; ls -1 /Users/bytedance/golang_pro/nano_cc/go-agent/.link/skills 2>/dev/null"}
```

### Result

```text
Error: exit status 1
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","resolved_command_path":"/Users/bytedance/golang_pro/nano_cc/go-agent/.link/skills","denial_reason":"Error: exit status 1"}
```

## 2026-06-13T09:05:20Z · bash · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_00_VqnaAwVSoj7QLRmcuYDw1876

### Arguments

```json
{"command": "ls -1 ~/.link/skills/ 2>&1 || echo \"NO_USER_SKILLS\""}
```

### Result

```text
ls: /Users/bytedance/.link/skills/: No such file or directory
NO_USER_SKILLS
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"ls: /Users/bytedance/.link/skills/: No such file or directory\nNO_USER_SKILLS"}
```

## 2026-06-13T09:05:22Z · bash · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_00_fFU83L4gqpxYCoXdBthw1562

### Arguments

```json
{"command": "ls -1 .link/skills/ 2>&1 || echo \"NO_WORKSPACE_SKILLS\""}
```

### Result

```text
ls: .link/skills/: No such file or directory
NO_WORKSPACE_SKILLS
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"ls: .link/skills/: No such file or directory\nNO_WORKSPACE_SKILLS"}
```

## 2026-06-13T09:06:45Z · bash · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_00_myCKm4cDeqePacLVsWpO9958

### Arguments

```json
{"command": "find . -type f -not -path './.git/*' -not -path '*/node_modules/*' | head -100"}
```

### Result

```text
./TUI化改造设计文档.md
./LINK_MD系统提示词注入设计文档.md
./go.mod
./.DS_Store
./TUI显示与事件循环修复设计文档.md
./task_outputs/a56df192-a2a8-4e64-8132-15c6ba71a5b7/tools.md.lock
./task_outputs/a56df192-a2a8-4e64-8132-15c6ba71a5b7/tools.md
./task_outputs/d8905f1b-6e6b-4686-a0a6-6465a6458fad/tools.md.lock
./task_outputs/d8905f1b-6e6b-4686-a0a6-6465a6458fad/tools.md
./config.json
./workspace/pdf-test-intro.pdf
./workspace/english-morning-reading/evals.json
./workspace/english-morning-reading/SKILL.md
./workspace/english-morning-reading-workspace/iteration-1/eval-2/eval_metadata.json
./workspace/english-morning-reading-workspace/iteration-1/eval-2/without_skill/outputs/result.md
./workspace/english-morning-reading-workspace/iteration-1/eval-2/with_skill/outputs/result.md
./workspace/english-morning-reading-workspace/iteration-1/eval-3/eval_metadata.json
./workspace/english-morning-reading-workspace/iteration-1/eval-3/without_skill/outputs/result.md
./workspace/english-morning-reading-workspace/iteration-1/eval-3/with_skill/outputs/result.md
./workspace/english-morning-reading-workspace/iteration-1/eval-1/eval_metadata.json
./workspace/english-morning-reading-workspace/iteration-1/eval-1/without_skill/outputs/result.md
./workspace/english-morning-reading-workspace/iteration-1/eval-1/with_skill/outputs/result.md
./workspace/config.json
./workspace/logs/session_20260602_214755.log
./workspace/logs/session_20260602_220208.log
./workspace/logs/session_20260602_210121.log
./workspace/logs/session_20260602_235310.log
./workspace/logs/session_20260608_022707.log
./workspace/logs/session_20260602_205800.log
./workspace/logs/session_20260602_174854.log
./workspace/logs/session_20260602_194455.log
./workspace/logs/session_20260602_211557.log
./workspace/basic.go
./workspace/skills/english-morning-reading/SKILL.md
./workspace/skills/pdf-generator/evals/evals.json
./workspace/skills/pdf-generator/scripts/generate_pdf.py
./workspace/skills/pdf-generator/SKILL.md
./workspace/skills/weekly-report-generator/evals/evals.json
./workspace/skills/weekly-report-generator/references/best-practices.md
./workspace/skills/weekly-report-generator/scripts/generate_report.py
./workspace/skills/weekly-report-generator/SKILL.md
./workspace/skills/code-assistant/SKILL.md
./workspace/skills/skill-creator/eval-viewer/generate_review.py
./workspace/skills/skill-creator/eval-viewer/viewer.html
./workspace/skills/skill-creator/references/schemas.md
./workspace/skills/skill-creator/agents/grader.md
./workspace/skills/skill-creator/agents/comparator.md
./workspace/skills/skill-creator/agents/analyzer.md
./workspace/skills/skill-creator/scripts/run_eval.py
./workspace/skills/skill-creator/scripts/package_skill.py
./workspace/skills/skill-creator/scripts/quick_validate.py
./workspace/skills/skill-creator/scripts/improve_description.py
./workspace/skills/skill-creator/scripts/aggregate_benchmark.py
./workspace/skills/skill-creator/scripts/__init__.py
./workspace/skills/skill-creator/scripts/run_loop.py
./workspace/skills/skill-creator/scripts/generate_report.py
./workspace/skills/skill-creator/scripts/utils.py
./workspace/skills/skill-creator/SKILL.md
./workspace/skills/skill-creator/LICENSE.txt
./workspace/skills/skill-creator/assets/eval_review.html
./workspace/skills/writing-assistant/SKILL.md
./workspace/system_prompt.md
./internal/assistant/prompt_test.go
./internal/assistant/prompt.go
./internal/tools/definitions.go
./internal/tools/persisted_output_test.go
./internal/tools/file_ops.go
./internal/tools/handlers.go
./internal/tools/bash.go
./internal/tools/definitions_test.go
./internal/tools/todo_write.go
./internal/tools/persisted_output.go
./internal/tools/load_skill.go
./internal/tools/runtime_env.go
./internal/tools/validation.go
./internal/tools/runtime_test.go
./internal/tools/todo_write_test.go
./internal/tools/path_guard.go
./internal/tools/validation_test.go
./internal/llm/client.go
./internal/llm/errors.go
./internal/logger/logger.go
./internal/logger/logger_test.go
./internal/idgen/idgen.go
./internal/config/paths.go
./internal/config/config.go
./internal/config/local_config.go
./internal/config/local_config_test.go
./internal/config/helpers.go
./internal/tui/events_test.go
./internal/tui/events.go
./internal/tui/app.go
./internal/agent/runtime/memory.go
./internal/agent/runtime/prompt_builder.go
./internal/agent/runtime/context_compression.go
./internal/agent/runtime/subagent.go
./internal/agent/runtime/runtime.go
./internal/agent/runtime/hook_bridge.go
./internal/agent/runtime/tool_registry.go
./internal/agent/runtime/summarizer.go
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"./TUI化改造设计文档.md\n./LINK_MD系统提示词注入设计文档.md\n./go.mod\n./.DS_Store\n./TUI显示与事件循环修复设计文档.md\n./task_outputs/a56df192-a2a8-4e64-8132-15c6ba71a5b7/tools.md.lock\n./task_outputs/a56df192-a2a8-4e64-8132-15c6ba71a5b7/tools.md\n./task_outputs/d8905f1b-6e6b-4686-a0a6-6465a6458fad/tools.md.lock\n./task_outputs/d8905f1b-6e6b-4686-a0a6-6465a6458fad/tools.md\n./config.json\n./workspace/pdf-test-intro.pdf\n./workspace/english-morning-reading/evals.json\n./workspa"}
```

## 2026-06-13T09:06:48Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_00_tOlnt2IFwD3Zd7KCAf1v6054

### Arguments

```json
{"path": "./go.mod"}
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

## 2026-06-13T09:06:48Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_01_ut8CdsK7kKk2wPMtFu6v6829

### Arguments

```json
{"path": "./config.json"}
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

## 2026-06-13T09:06:48Z · bash · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_02_4EvTArdYiZaQOq4ENsnj1966

### Arguments

```json
{"command": "find . -maxdepth 1 -name '*.go' -o -name 'main.go' | head -10"}
```

### Result

```text
./architecture_test.go
./main.go
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"./architecture_test.go\n./main.go"}
```

## 2026-06-13T09:06:51Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_00_QjfMvqiGDczr40hGBdiB9760

### Arguments

```json
{"path": "./main.go"}
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

## 2026-06-13T09:06:51Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_01_phtbmDgKKHvguJ4RzX892037

### Arguments

```json
{"path": "./architecture_test.go"}
```

### Result

```text
package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNoLegacyServicePackageOrImports(t *testing.T) {
	legacyDir := filepath.Join("internal", "web")
	if _, err := os.Stat(legacyDir); !os.IsNotExist(err) {
		t.Fatalf("legacy service package must be removed after TUI migration; stat error: %v", err)
	}

	err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		forbiddenImport := "nano_cc/internal/" + "web/"
		if strings.Contains(string(data), forbiddenImport) {
			t.Fatalf("%s still imports legacy service package", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk source tree: %v", err)
	}
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package main\n\nimport (\n\t\"os\"\n\t\"path/filepath\"\n\t\"strings\"\n\t\"testing\"\n)\n\nfunc TestNoLegacyServicePackageOrImports(t *testing.T) {\n\tlegacyDir := filepath.Join(\"internal\", \"web\")\n\tif _, err := os.Stat(legacyDir); !os.IsNotExist(err) {\n\t\tt.Fatalf(\"legacy service package must be removed after TUI migration; stat error: %v\", err)\n\t}\n\n\terr := filepath.WalkDir(\".\", func(path string, d os.DirEntry, err error) error {\n\t\tif err != nil || d.IsDir() || !strings.HasSuffix(path, \".go\") {\n\t\t\treturn err\n\t\t}\n\t\tdat"}
```

## 2026-06-13T09:06:51Z · bash · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_02_EKD3YPEqytJM5i5YtHYG9675

### Arguments

```json
{"command": "find . -type f -name '*.go' -not -path './.git/*' | sort"}
```

### Result

```text
./architecture_test.go
./internal/agent/mcp/config_test.go
./internal/agent/mcp/config.go
./internal/agent/mcp/manager_test.go
./internal/agent/mcp/manager.go
./internal/agent/mcp/transport_test.go
./internal/agent/mcp/transport.go
./internal/agent/runtime/compression/compression_test.go
./internal/agent/runtime/compression/compression.go
./internal/agent/runtime/compression/conversation_memory_strategy.go
./internal/agent/runtime/compression/full_history_summarization.go
./internal/agent/runtime/compression/message_window_compression.go
./internal/agent/runtime/compression/reactive_compact.go
./internal/agent/runtime/compression/recent_tool_result_retention.go
./internal/agent/runtime/compression/token_estimator.go
./internal/agent/runtime/compression/tool_result_compression.go
./internal/agent/runtime/context_compression_test.go
./internal/agent/runtime/context_compression.go
./internal/agent/runtime/conversation_flow.go
./internal/agent/runtime/conversation_memory.go
./internal/agent/runtime/helpers.go
./internal/agent/runtime/hook_bridge.go
./internal/agent/runtime/hooks/manager.go
./internal/agent/runtime/hooks/stop.go
./internal/agent/runtime/hooks/tool.go
./internal/agent/runtime/hooks/types.go
./internal/agent/runtime/hooks/user_prompt.go
./internal/agent/runtime/memory_test.go
./internal/agent/runtime/memory.go
./internal/agent/runtime/prompt_builder.go
./internal/agent/runtime/runtime_test.go
./internal/agent/runtime/runtime.go
./internal/agent/runtime/subagent.go
./internal/agent/runtime/summarizer.go
./internal/agent/runtime/todo_reminder.go
./internal/agent/runtime/tool_registry.go
./internal/agent/storage/conversation_history_test.go
./internal/agent/storage/conversation_history.go
./internal/agent/storage/models.go
./internal/assistant/prompt_test.go
./internal/assistant/prompt.go
./internal/cli/root_test.go
./internal/cli/root.go
./internal/config/config.go
./internal/config/helpers.go
./internal/config/local_config_test.go
./internal/config/local_config.go
./internal/config/paths.go
./internal/idgen/idgen.go
./internal/llm/client.go
./internal/llm/errors.go
./internal/local/bootstrap.go
./internal/local/memory_store.go
./internal/local/persisted_output_files.go
./internal/local/session_history_store.go
./internal/local/store_test.go
./internal/local/store.go
./internal/local/tool_result_log.go
./internal/logger/logger_test.go
./internal/logger/logger.go
./internal/safety/path.go
./internal/sessions/render_test.go
./internal/sessions/skill_test.go
./internal/sessions/skill.go
./internal/textutil/textutil.go
./internal/tools/bash.go
./internal/tools/definitions_test.go
./internal/tools/definitions.go
./internal/tools/file_ops.go
./internal/tools/handlers.go
./internal/tools/load_skill.go
./internal/tools/path_guard.go
./internal/tools/persisted_output_test.go
./internal/tools/persisted_output.go
./internal/tools/runtime_env.go
./internal/tools/runtime_test.go
./internal/tools/todo_write_test.go
./internal/tools/todo_write.go
./internal/tools/validation_test.go
./internal/tools/validation.go
./internal/tui/app.go
./internal/tui/events_test.go
./internal/tui/events.go
./main.go
./workspace/basic.go
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"./architecture_test.go\n./internal/agent/mcp/config_test.go\n./internal/agent/mcp/config.go\n./internal/agent/mcp/manager_test.go\n./internal/agent/mcp/manager.go\n./internal/agent/mcp/transport_test.go\n./internal/agent/mcp/transport.go\n./internal/agent/runtime/compression/compression_test.go\n./internal/agent/runtime/compression/compression.go\n./internal/agent/runtime/compression/conversation_memory_strategy.go\n./internal/agent/runtime/compression/full_history_summarization.go\n./internal/agent/runtim"}
```

## 2026-06-13T09:06:53Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_00_CmFoSvpVk4qlhRb0TfeS9776

### Arguments

```json
{"path": "./internal/cli/root.go"}
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

## 2026-06-13T09:06:53Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_01_lXYjmtphpf4zBuqrmFhC5201

### Arguments

```json
{"path": "./internal/config/config.go"}
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

## 2026-06-13T09:06:53Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_02_s10FexfsWtBPPGWEJ9p46577

### Arguments

```json
{"path": "./internal/tui/app.go"}
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
	program := tea.NewProgram(NewModel(runtimeService, session), tea.WithContext(ctx), tea.WithAltScreen())
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

## 2026-06-13T09:06:53Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_03_1PdpYXa2tbOviZalskk80112

### Arguments

```json
{"path": "./internal/agent/runtime/runtime.go"}
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

## 2026-06-13T09:06:57Z · read_file · rejected

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_00_8gZpdPSLnp0P2pzB7UPF5146

### Arguments

```json
{"path": "./internal/agent/runtime/respond.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/respond.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/respond.go: no such file or directory"}
```

## 2026-06-13T09:06:57Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_01_4f5r68YNb5CTxyByNnxA6562

### Arguments

```json
{"path": "./go.sum"}
```

### Result

```text
github.com/MakeNowJust/heredoc v1.0.0 h1:cXCdzVdstXyiTqTvfqk9SDHpKNjxuom+DOlyEeQ4pzQ=
github.com/MakeNowJust/heredoc v1.0.0/go.mod h1:mG5amYoWBHf8vpLOuehzbGGw0EHxpZZ6lCpQ4fNJ8LE=
github.com/alecthomas/assert/v2 v2.11.0 h1:2Q9r3ki8+JYXvGsDyBXwH3LcJ+WK5D0gc5E8vS6K3D0=
github.com/alecthomas/assert/v2 v2.11.0/go.mod h1:Bze95FyfUr7x34QZrjL+XP+0qgp/zg8yS+TtBj1WA3k=
github.com/alecthomas/chroma/v2 v2.20.0 h1:sfIHpxPyR07/Oylvmcai3X/exDlE8+FA820NTz+9sGw=
github.com/alecthomas/chroma/v2 v2.20.0/go.mod h1:e7tViK0xh/Nf4BYHl00ycY6rV7b8iXBksI9E359yNmA=
github.com/alecthomas/repr v0.5.1 h1:E3G4t2QbHTSNpPKBgMTln5KLkZHLOcU7r37J4pXBuIg=
github.com/alecthomas/repr v0.5.1/go.mod h1:Fr0507jx4eOXV7AlPV6AVZLYrLIuIeSOWtW57eE/O/4=
github.com/atotto/clipboard v0.1.4 h1:EH0zSVneZPSuFR11BlR9YppQTVDbh5+16AmcJi4g1z4=
github.com/atotto/clipboard v0.1.4/go.mod h1:ZY9tmq7sm5xIbd9bOK4onWV4S6X0u6GY7Vn0Yu86PYI=
github.com/aymanbagabas/go-osc52/v2 v2.0.1 h1:HwpRHbFMcZLEVr42D4p7XBqjyuxQH5SMiErDT4WkJ2k=
github.com/aymanbagabas/go-osc52/v2 v2.0.1/go.mod h1:uYgXzlJ7ZpABp8OJ+exZzJJhRNQ2ASbcXHWsFqH8hp8=
github.com/aymanbagabas/go-udiff v0.3.1 h1:LV+qyBQ2pqe0u42ZsUEtPiCaUoqgA9gYRDs3vj1nolY=
github.com/aymanbagabas/go-udiff v0.3.1/go.mod h1:G0fsKmG+P6ylD0r6N/KgQD/nWzgfnl8ZBcNLgcbrw8E=
github.com/aymerick/douceur v0.2.0 h1:Mv+mAeH1Q+n9Fr+oyamOlAkUNPWPlA8PPGR0QAaYuPk=
github.com/aymerick/douceur v0.2.0/go.mod h1:wlT5vV2O3h55X9m7iVYN0TBM0NH/MmbLnd30/FjWUq4=
github.com/charmbracelet/bubbles v1.0.0 h1:12J8/ak/uCZEMQ6KU7pcfwceyjLlWsDLAxB5fXonfvc=
github.com/charmbracelet/bubbles v1.0.0/go.mod h1:9d/Zd5GdnauMI5ivUIVisuEm3ave1XwXtD1ckyV6r3E=
github.com/charmbracelet/bubbletea v1.3.10 h1:otUDHWMMzQSB0Pkc87rm691KZ3SWa4KUlvF9nRvCICw=
github.com/charmbracelet/bubbletea v1.3.10/go.mod h1:ORQfo0fk8U+po9VaNvnV95UPWA1BitP1E0N6xJPlHr4=
github.com/charmbracelet/colorprofile v0.4.1 h1:a1lO03qTrSIRaK8c3JRxJDZOvhvIeSco3ej+ngLk1kk=
github.com/charmbracelet/colorprofile v0.4.1/go.mod h1:U1d9Dljmdf9DLegaJ0nGZNJvoXAhayhmidOdcBwAvKk=
github.com/charmbracelet/glamour v1.0.0 h1:AWMLOVFHTsysl4WV8T8QgkQ0s/ZNZo7CiE4WKhk8l08=
github.com/charmbracelet/glamour v1.0.0/go.mod h1:DSdohgOBkMr2ZQNhw4LZxSGpx3SvpeujNoXrQyH2hxo=
github.com/charmbracelet/lipgloss v1.1.1-0.20250404203927-76690c660834 h1:ZR7e0ro+SZZiIZD7msJyA+NjkCNNavuiPBLgerbOziE=
github.com/charmbracelet/lipgloss v1.1.1-0.20250404203927-76690c660834/go.mod h1:aKC/t2arECF6rNOnaKaVU6y4t4ZeHQzqfxedE/VkVhA=
github.com/charmbracelet/x/ansi v0.11.6 h1:GhV21SiDz/45W9AnV2R61xZMRri5NlLnl6CVF7ihZW8=
github.com/charmbracelet/x/ansi v0.11.6/go.mod h1:2JNYLgQUsyqaiLovhU2Rv/pb8r6ydXKS3NIttu3VGZQ=
github.com/charmbracelet/x/cellbuf v0.0.15 h1:ur3pZy0o6z/R7EylET877CBxaiE1Sp1GMxoFPAIztPI=
github.com/charmbracelet/x/cellbuf v0.0.15/go.mod h1:J1YVbR7MUuEGIFPCaaZ96KDl5NoS0DAWkskup+mOY+Q=
github.com/charmbracelet/x/exp/golden v0.0.0-20241011142426-46044092ad91 h1:payRxjMjKgx2PaCWLZ4p3ro9y97+TVLZNaRZgJwSVDQ=
github.com/charmbracelet/x/exp/golden v0.0.0-20241011142426-46044092ad91/go.mod h1:wDlXFlCrmJ8J+swcL/MnGUuYnqgQdW9rhSD61oNMb6U=
github.com/charmbracelet/x/exp/slice v0.0.0-20250327172914-2fdc97757edf h1:rLG0Yb6MQSDKdB52aGX55JT1oi0P0Kuaj7wi1bLUpnI=
github.com/charmbracelet/x/exp/slice v0.0.0-20250327172914-2fdc97757edf/go.mod h1:B3UgsnsBZS/eX42BlaNiJkD1pPOUa+oF1IYC6Yd2CEU=
github.com/charmbracelet/x/term v0.2.2 h1:xVRT/S2ZcKdhhOuSP4t5cLi5o+JxklsoEObBSgfgZRk=
github.com/charmbracelet/x/term v0.2.2/go.mod h1:kF8CY5RddLWrsgVwpw4kAa6TESp6EB5y3uxGLeCqzAI=
github.com/clipperhouse/displaywidth v0.9.0 h1:Qb4KOhYwRiN3viMv1v/3cTBlz3AcAZX3+y9OLhMtAtA=
github.com/clipperhouse/displaywidth v0.9.0/go.mod h1:aCAAqTlh4GIVkhQnJpbL0T/WfcrJXHcj8C0yjYcjOZA=
github.com/clipperhouse/stringish v0.1.1 h1:+NSqMOr3GR6k1FdRhhnXrLfztGzuG+VuFDfatpWHKCs=
github.com/clipperhouse/stringish v0.1.1/go.mod h1:v/WhFtE1q0ovMta2+m+UbpZ+2/HEXNWYXQgCt4hdOzA=
github.com/clipperhouse/uax29/v2 v2.5.0 h1:x7T0T4eTHDONxFJsL94uKNKPHrclyFI0lm7+w94cO8U=
github.com/clipperhouse/uax29/v2 v2.5.0/go.mod h1:Wn1g7MK6OoeDT0vL+Q0SQLDz/KpfsVRgg6W7ihQeh4g=
github.com/dlclark/regexp2 v1.11.5 h1:Q/sSnsKerHeCkc/jSTNq1oCm7KiVgUMZRDUoRu0JQZQ=
github.com/dlclark/regexp2 v1.11.5/go.mod h1:DHkYz0B9wPfa6wondMfaivmHpzrQ3v9q8cnmRbL6yW8=
github.com/erikgeiser/coninput v0.0.0-20211004153227-1c3628e74d0f h1:Y/CXytFA4m6baUTXGLOoWe4PQhGxaX0KpnayAqC48p4=
github.com/erikgeiser/coninput v0.0.0-20211004153227-1c3628e74d0f/go.mod h1:vw97MGsxSvLiUE2X8qFplwetxpGLQrlU1Q9AUEIzCaM=
github.com/golang-jwt/jwt/v5 v5.3.1 h1:kYf81DTWFe7t+1VvL7eS+jKFVWaUnK9cB1qbwn63YCY=
github.com/golang-jwt/jwt/v5 v5.3.1/go.mod h1:fxCRLWMO43lRc8nhHWY6LGqRcf+1gQWArsqaEUEa5bE=
github.com/google/go-cmp v0.7.0 h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=
github.com/google/go-cmp v0.7.0/go.mod h1:pXiqmnSA92OHEEa9HXL2W4E7lf9JzCmGVUdgjX3N/iU=
github.com/google/jsonschema-go v0.4.3 h1:/DBOLZTfDow7pe2GmaJNhltueGTtDKICi8V8p+DQPd0=
github.com/google/jsonschema-go v0.4.3/go.mod h1:r5quNTdLOYEz95Ru18zA0ydNbBuYoo9tgaYcxEYhJVE=
github.com/gorilla/css v1.0.1 h1:ntNaBIghp6JmvWnxbZKANoLyuXTPZ4cAMlo6RyhlbO8=
github.com/gorilla/css v1.0.1/go.mod h1:BvnYkspnSzMmwRK+b8/xgNPLiIuNZr6vbZBTPQ2A3b0=
github.com/hexops/gotextdiff v1.0.3 h1:gitA9+qJrrTCsiCl7+kh75nPqQt1cx4ZkudSTLoUqJM=
github.com/hexops/gotextdiff v1.0.3/go.mod h1:pSWU5MAI3yDq+fZBTazCSJysOMbxWL1BSow5/V2vxeg=
github.com/lucasb-eyer/go-colorful v1.3.0 h1:2/yBRLdWBZKrf7gB40FoiKfAWYQ0lqNcbuQwVHXptag=
github.com/lucasb-eyer/go-colorful v1.3.0/go.mod h1:R4dSotOR9KMtayYi1e77YzuveK+i7ruzyGqttikkLy0=
github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
github.com/mattn/go-localereader v0.0.1 h1:ygSAOl7ZXTx4RdPYinUpg6W99U8jWvWi9Ye2JC/oIi4=
github.com/mattn/go-localereader v0.0.1/go.mod h1:8fBrzywKY7BI3czFoHkuzRoWE9C+EiG4R1k4Cjx5p88=
github.com/mattn/go-runewidth v0.0.12/go.mod h1:RAqKPSqVFrSLVXbA8x7dzmKdmGzieGRCM46jaSJTDAk=
github.com/mattn/go-runewidth v0.0.19 h1:v++JhqYnZuu5jSKrk9RbgF5v4CGUjqRfBm05byFGLdw=
github.com/mattn/go-runewidth v0.0.19/go.mod h1:XBkDxAl56ILZc9knddidhrOlY5R/pDhgLpndooCuJAs=
github.com/microcosm-cc/bluemonday v1.0.27 h1:MpEUotklkwCSLeH+Qdx1VJgNqLlpY2KXwXFM08ygZfk=
github.com/microcosm-cc/bluemonday v1.0.27/go.mod h1:jFi9vgW+H7c3V0lb6nR74Ib/DIB5OBs92Dimizgw2cA=
github.com/modelcontextprotocol/go-sdk v1.6.1 h1:0zOSupjKUxPKSocPT1Wtago+mUHU2/uZ4xSOY0FGReU=
github.com/modelcontextprotocol/go-sdk v1.6.1/go.mod h1:kzm3kzFL1/+AziGOE0nUs3gvPoNxMCvkxokMkuFapXQ=
github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6 h1:ZK8zHtRHOkbHy6Mmr5D264iyp3TiX5OmNcI5cIARiQI=
github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6/go.mod h1:CJlz5H+gyd6CUWT45Oy4q24RdLyn7Md9Vj2/ldJBSIo=
github.com/muesli/cancelreader v0.2.2 h1:3I4Kt4BQjOR54NavqnDogx/MIoWBFa0StPA8ELUXHmA=
github.com/muesli/cancelreader v0.2.2/go.mod h1:3XuTXfFS2VjM+HTLZY9Ak0l6eUKfijIfMUZ4EgX0QYo=
github.com/muesli/reflow v0.3.0 h1:IFsN6K9NfGtjeggFP+68I4chLZV2yIKsXJFNZ+eWh6s=
github.com/muesli/reflow v0.3.0/go.mod h1:pbwTDkVPibjO2kyvBQRBxTWEEGDGq0FlB1BIKtnHY/8=
github.com/muesli/termenv v0.16.0 h1:S5AlUN9dENB57rsbnkPyfdGuWIlkmzJjbFf0Tf5FWUc=
github.com/muesli/termenv v0.16.0/go.mod h1:ZRfOIKPFDYQoDFF4Olj7/QJbW60Ol/kL1pU3VfY/Cnk=
github.com/rivo/uniseg v0.1.0/go.mod h1:J6wj4VEh+S6ZtnVlnTBMWIodfgj8LQOQFoIToxlJtxc=
github.com/rivo/uniseg v0.2.0/go.mod h1:J6wj4VEh+S6ZtnVlnTBMWIodfgj8LQOQFoIToxlJtxc=
github.com/rivo/uniseg v0.4.7 h1:WUdvkW8uEhrYfLC4ZzdpI2ztxP1I582+49Oc5Mq64VQ=
github.com/rivo/uniseg v0.4.7/go.mod h1:FN3SvrM+Zdj16jyLfmOkMNblXMcoc8DfTHruCPUcx88=
github.com/sashabaranov/go-openai v1.41.2 h1:vfPRBZNMpnqu8ELsclWcAvF19lDNgh1t6TVfFFOPiSM=
github.com/sashabaranov/go-openai v1.41.2/go.mod h1:lj5b/K+zjTSFxVLijLSTDZuP7adOgerWeFyZLUhAKRg=
github.com/segmentio/asm v1.1.3 h1:WM03sfUOENvvKexOLp+pCqgb/WDjsi7EK8gIsICtzhc=
github.com/segmentio/asm v1.1.3/go.mod h1:Ld3L4ZXGNcSLRg4JBsZ3//1+f/TjYl0Mzen/DQy1EJg=
github.com/segmentio/encoding v0.5.4 h1:OW1VRern8Nw6ITAtwSZ7Idrl3MXCFwXHPgqESYfvNt0=
github.com/segmentio/encoding v0.5.4/go.mod h1:HS1ZKa3kSN32ZHVZ7ZLPLXWvOVIiZtyJnO1gPH1sKt0=
github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e h1:JVG44RsyaB9T2KIHavMF/ppJZNG9ZpyihvCd0w101no=
github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e/go.mod h1:RbqR21r5mrJuqunuUZ/Dhy/avygyECGrLceyNeo4LiM=
github.com/yosida95/uritemplate/v3 v3.0.2 h1:Ed3Oyj9yrmi9087+NczuL5BwkIc4wvTb5zIM+UJPGz4=
github.com/yosida95/uritemplate/v3 v3.0.2/go.mod h1:ILOh0sOhIJR3+L/8afwt/kE++YT040gmv5BQTMR2HP4=
github.com/yuin/goldmark v1.7.13 h1:GPddIs617DnBLFFVJFgpo1aBfe/4xcvMc3SB5t/D0pA=
github.com/yuin/goldmark v1.7.13/go.mod h1:ip/1k0VRfGynBgxOz0yCqHrbZXhcjxyuS66Brc7iBKg=
github.com/yuin/goldmark-emoji v1.0.6 h1:QWfF2FYaXwL74tfGOW5izeiZepUDroDJfWubQI9HTHs=
github.com/yuin/goldmark-emoji v1.0.6/go.mod h1:ukxJDKFpdFb5x0a5HqbdlcKtebh086iJpI31LTKmWuA=
golang.org/x/exp v0.0.0-20231006140011-7918f672742d h1:jtJma62tbqLibJ5sFQz8bKtEM8rJBtfilJ2qTU199MI=
golang.org/x/exp v0.0.0-20231006140011-7918f672742d/go.mod h1:ldy0pHrwJyGW56pPQzzkH36rKxoZW1tw7ZJpeKx+hdo=
golang.org/x/net v0.38.0 h1:vRMAPTMaeGqVhG5QyLJHqNDwecKTomGeqbnfZyKlBI8=
golang.org/x/net v0.38.0/go.mod h1:ivrbrMbzFq5J41QOQh0siUuly180yBYtLp+CKbEaFx8=
golang.org/x/oauth2 v0.35.0 h1:Mv2mzuHuZuY2+bkyWXIHMfhNdJAdwW3FuWeCPYN5GVQ=
golang.org/x/oauth2 v0.35.0/go.mod h1:lzm5WQJQwKZ3nwavOZ3IS5Aulzxi68dUSgRHujetwEA=
golang.org/x/sys v0.0.0-20210809222454-d867a43fc93e/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.41.0 h1:Ivj+2Cp/ylzLiEU89QhWblYnOE9zerudt9Ftecq2C6k=
golang.org/x/sys v0.41.0/go.mod h1:OgkHotnGiDImocRcuBABYBEXf8A9a87e/uXjp9XT3ks=
golang.org/x/term v0.36.0 h1:zMPR+aF8gfksFprF/Nc/rd1wRS1EI6nDBGyWAvDzx2Q=
golang.org/x/term v0.36.0/go.mod h1:Qu394IJq6V6dCBRgwqshf3mPF85AqzYEzofzRdZkWss=
golang.org/x/text v0.30.0 h1:yznKA/E9zq54KzlzBEAWn1NXSQ8DIp/NYMy88xJjl4k=
golang.org/x/text v0.30.0/go.mod h1:yDdHFIX9t+tORqspjENWgzaCVXgk0yYnYuSZ8UzzBVM=
golang.org/x/tools v0.42.0 h1:uNgphsn75Tdz5Ji2q36v/nsFSfR/9BRFvqhGBaJGd5k=
golang.org/x/tools v0.42.0/go.mod h1:Ma6lCIwGZvHK6XtgbswSoWroEkhugApmsXyrUmBhfr0=

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"github.com/MakeNowJust/heredoc v1.0.0 h1:cXCdzVdstXyiTqTvfqk9SDHpKNjxuom+DOlyEeQ4pzQ=\ngithub.com/MakeNowJust/heredoc v1.0.0/go.mod h1:mG5amYoWBHf8vpLOuehzbGGw0EHxpZZ6lCpQ4fNJ8LE=\ngithub.com/alecthomas/assert/v2 v2.11.0 h1:2Q9r3ki8+JYXvGsDyBXwH3LcJ+WK5D0gc5E8vS6K3D0=\ngithub.com/alecthomas/assert/v2 v2.11.0/go.mod h1:Bze95FyfUr7x34QZrjL+XP+0qgp/zg8yS+TtBj1WA3k=\ngithub.com/alecthomas/chroma/v2 v2.20.0 h1:sfIHpxPyR07/Oylvmcai3X/exDlE8+FA820NTz+9sGw=\ngithub.com/alecthomas/chroma/v2 v2.20.0/go.mod h1:"}
```

## 2026-06-13T09:06:57Z · bash · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_02_pQBcDIrboxeRP46u12Z31702

### Arguments

```json
{"command": "find . -type f -name '*.go' -not -path './.git/*' | sort"}
```

### Result

```text
./architecture_test.go
./internal/agent/mcp/config_test.go
./internal/agent/mcp/config.go
./internal/agent/mcp/manager_test.go
./internal/agent/mcp/manager.go
./internal/agent/mcp/transport_test.go
./internal/agent/mcp/transport.go
./internal/agent/runtime/compression/compression_test.go
./internal/agent/runtime/compression/compression.go
./internal/agent/runtime/compression/conversation_memory_strategy.go
./internal/agent/runtime/compression/full_history_summarization.go
./internal/agent/runtime/compression/message_window_compression.go
./internal/agent/runtime/compression/reactive_compact.go
./internal/agent/runtime/compression/recent_tool_result_retention.go
./internal/agent/runtime/compression/token_estimator.go
./internal/agent/runtime/compression/tool_result_compression.go
./internal/agent/runtime/context_compression_test.go
./internal/agent/runtime/context_compression.go
./internal/agent/runtime/conversation_flow.go
./internal/agent/runtime/conversation_memory.go
./internal/agent/runtime/helpers.go
./internal/agent/runtime/hook_bridge.go
./internal/agent/runtime/hooks/manager.go
./internal/agent/runtime/hooks/stop.go
./internal/agent/runtime/hooks/tool.go
./internal/agent/runtime/hooks/types.go
./internal/agent/runtime/hooks/user_prompt.go
./internal/agent/runtime/memory_test.go
./internal/agent/runtime/memory.go
./internal/agent/runtime/prompt_builder.go
./internal/agent/runtime/runtime_test.go
./internal/agent/runtime/runtime.go
./internal/agent/runtime/subagent.go
./internal/agent/runtime/summarizer.go
./internal/agent/runtime/todo_reminder.go
./internal/agent/runtime/tool_registry.go
./internal/agent/storage/conversation_history_test.go
./internal/agent/storage/conversation_history.go
./internal/agent/storage/models.go
./internal/assistant/prompt_test.go
./internal/assistant/prompt.go
./internal/cli/root_test.go
./internal/cli/root.go
./internal/config/config.go
./internal/config/helpers.go
./internal/config/local_config_test.go
./internal/config/local_config.go
./internal/config/paths.go
./internal/idgen/idgen.go
./internal/llm/client.go
./internal/llm/errors.go
./internal/local/bootstrap.go
./internal/local/memory_store.go
./internal/local/persisted_output_files.go
./internal/local/session_history_store.go
./internal/local/store_test.go
./internal/local/store.go
./internal/local/tool_result_log.go
./internal/logger/logger_test.go
./internal/logger/logger.go
./internal/safety/path.go
./internal/sessions/render_test.go
./internal/sessions/skill_test.go
./internal/sessions/skill.go
./internal/textutil/textutil.go
./internal/tools/bash.go
./internal/tools/definitions_test.go
./internal/tools/definitions.go
./internal/tools/file_ops.go
./internal/tools/handlers.go
./internal/tools/load_skill.go
./internal/tools/path_guard.go
./internal/tools/persisted_output_test.go
./internal/tools/persisted_output.go
./internal/tools/runtime_env.go
./internal/tools/runtime_test.go
./internal/tools/todo_write_test.go
./internal/tools/todo_write.go
./internal/tools/validation_test.go
./internal/tools/validation.go
./internal/tui/app.go
./internal/tui/events_test.go
./internal/tui/events.go
./main.go
./workspace/basic.go
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"./architecture_test.go\n./internal/agent/mcp/config_test.go\n./internal/agent/mcp/config.go\n./internal/agent/mcp/manager_test.go\n./internal/agent/mcp/manager.go\n./internal/agent/mcp/transport_test.go\n./internal/agent/mcp/transport.go\n./internal/agent/runtime/compression/compression_test.go\n./internal/agent/runtime/compression/compression.go\n./internal/agent/runtime/compression/conversation_memory_strategy.go\n./internal/agent/runtime/compression/full_history_summarization.go\n./internal/agent/runtim"}
```

## 2026-06-13T09:07:00Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_00_S2xrYygUorxkX1nxeHrs4093

### Arguments

```json
{"path": "./main.go"}
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

## 2026-06-13T09:07:00Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_01_9n7b1VfPNrmDTCKPEHAo1181

### Arguments

```json
{"path": "./internal/agent/runtime/conversation_flow.go"}
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

## 2026-06-13T09:07:00Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_02_v3eEluudypCF3hqmPxIP0236

### Arguments

```json
{"path": "./internal/llm/client.go"}
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

## 2026-06-13T09:07:00Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_03_4zYS8aEzB82H2dPsisgG2917

### Arguments

```json
{"path": "./internal/tools/definitions.go"}
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

## 2026-06-13T09:07:03Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_00_7RiTmVIWqTTkcJHI91Vk9005

### Arguments

```json
{"path": "./go.mod"}
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

## 2026-06-13T09:07:03Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_01_3Q2aXMFHyaSDgMC24ven6501

### Arguments

```json
{"path": "./config.json"}
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

## 2026-06-13T09:07:03Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_02_R8CQsopiLBc8064K2Wrc2646

### Arguments

```json
{"path": "./internal/cli/root.go"}
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

## 2026-06-13T09:07:03Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_03_9dyfQ4DE7sJkZoy1RwU05068

### Arguments

```json
{"path": "./internal/config/config.go"}
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

## 2026-06-13T09:07:05Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_00_aVqi9OpLH67dey9mmPzV9102

### Arguments

```json
{"path": "./internal/agent/runtime/runtime.go"}
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

## 2026-06-13T09:07:05Z · read_file · rejected

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_01_exB32jsStQEk2PsNZmVq9077

### Arguments

```json
{"path": "./internal/agent/runtime/respond.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/respond.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/respond.go: no such file or directory"}
```

## 2026-06-13T09:07:05Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_02_3L9SycA7RmGVeIC3Urtn7036

### Arguments

```json
{"path": "./internal/tui/app.go"}
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
	program := tea.NewProgram(NewModel(runtimeService, session), tea.WithContext(ctx), tea.WithAltScreen())
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

## 2026-06-13T09:07:08Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_00_UdAejM8BqDCn9IOZzROb6028

### Arguments

```json
{"path": "./internal/agent/runtime/conversation_flow.go"}
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

## 2026-06-13T09:07:08Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_01_yaOIsyrISHbbaka7p7D76981

### Arguments

```json
{"path": "./internal/llm/client.go"}
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

## 2026-06-13T09:07:08Z · read_file · rejected

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_02_8CD4mxYOZ3o0KX9wbWnc1295

### Arguments

```json
{"path": "./internal/agent/storage/storage.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/storage.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/storage.go: no such file or directory"}
```

## 2026-06-13T09:07:08Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_03_m9ksj3sn6wwZkjZRm3iC2204

### Arguments

```json
{"path": "./architecture_test.go"}
```

### Result

```text
package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNoLegacyServicePackageOrImports(t *testing.T) {
	legacyDir := filepath.Join("internal", "web")
	if _, err := os.Stat(legacyDir); !os.IsNotExist(err) {
		t.Fatalf("legacy service package must be removed after TUI migration; stat error: %v", err)
	}

	err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		forbiddenImport := "nano_cc/internal/" + "web/"
		if strings.Contains(string(data), forbiddenImport) {
			t.Fatalf("%s still imports legacy service package", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk source tree: %v", err)
	}
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package main\n\nimport (\n\t\"os\"\n\t\"path/filepath\"\n\t\"strings\"\n\t\"testing\"\n)\n\nfunc TestNoLegacyServicePackageOrImports(t *testing.T) {\n\tlegacyDir := filepath.Join(\"internal\", \"web\")\n\tif _, err := os.Stat(legacyDir); !os.IsNotExist(err) {\n\t\tt.Fatalf(\"legacy service package must be removed after TUI migration; stat error: %v\", err)\n\t}\n\n\terr := filepath.WalkDir(\".\", func(path string, d os.DirEntry, err error) error {\n\t\tif err != nil || d.IsDir() || !strings.HasSuffix(path, \".go\") {\n\t\t\treturn err\n\t\t}\n\t\tdat"}
```

## 2026-06-13T09:07:11Z · bash · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_00_K2F30Hx3QVcp1NUmLK6s8144

### Arguments

```json
{"command": "find . -type f -name '*.go' -not -path './.git/*' | sort"}
```

### Result

```text
./architecture_test.go
./internal/agent/mcp/config_test.go
./internal/agent/mcp/config.go
./internal/agent/mcp/manager_test.go
./internal/agent/mcp/manager.go
./internal/agent/mcp/transport_test.go
./internal/agent/mcp/transport.go
./internal/agent/runtime/compression/compression_test.go
./internal/agent/runtime/compression/compression.go
./internal/agent/runtime/compression/conversation_memory_strategy.go
./internal/agent/runtime/compression/full_history_summarization.go
./internal/agent/runtime/compression/message_window_compression.go
./internal/agent/runtime/compression/reactive_compact.go
./internal/agent/runtime/compression/recent_tool_result_retention.go
./internal/agent/runtime/compression/token_estimator.go
./internal/agent/runtime/compression/tool_result_compression.go
./internal/agent/runtime/context_compression_test.go
./internal/agent/runtime/context_compression.go
./internal/agent/runtime/conversation_flow.go
./internal/agent/runtime/conversation_memory.go
./internal/agent/runtime/helpers.go
./internal/agent/runtime/hook_bridge.go
./internal/agent/runtime/hooks/manager.go
./internal/agent/runtime/hooks/stop.go
./internal/agent/runtime/hooks/tool.go
./internal/agent/runtime/hooks/types.go
./internal/agent/runtime/hooks/user_prompt.go
./internal/agent/runtime/memory_test.go
./internal/agent/runtime/memory.go
./internal/agent/runtime/prompt_builder.go
./internal/agent/runtime/runtime_test.go
./internal/agent/runtime/runtime.go
./internal/agent/runtime/subagent.go
./internal/agent/runtime/summarizer.go
./internal/agent/runtime/todo_reminder.go
./internal/agent/runtime/tool_registry.go
./internal/agent/storage/conversation_history_test.go
./internal/agent/storage/conversation_history.go
./internal/agent/storage/models.go
./internal/assistant/prompt_test.go
./internal/assistant/prompt.go
./internal/cli/root_test.go
./internal/cli/root.go
./internal/config/config.go
./internal/config/helpers.go
./internal/config/local_config_test.go
./internal/config/local_config.go
./internal/config/paths.go
./internal/idgen/idgen.go
./internal/llm/client.go
./internal/llm/errors.go
./internal/local/bootstrap.go
./internal/local/memory_store.go
./internal/local/persisted_output_files.go
./internal/local/session_history_store.go
./internal/local/store_test.go
./internal/local/store.go
./internal/local/tool_result_log.go
./internal/logger/logger_test.go
./internal/logger/logger.go
./internal/safety/path.go
./internal/sessions/render_test.go
./internal/sessions/skill_test.go
./internal/sessions/skill.go
./internal/textutil/textutil.go
./internal/tools/bash.go
./internal/tools/definitions_test.go
./internal/tools/definitions.go
./internal/tools/file_ops.go
./internal/tools/handlers.go
./internal/tools/load_skill.go
./internal/tools/path_guard.go
./internal/tools/persisted_output_test.go
./internal/tools/persisted_output.go
./internal/tools/runtime_env.go
./internal/tools/runtime_test.go
./internal/tools/todo_write_test.go
./internal/tools/todo_write.go
./internal/tools/validation_test.go
./internal/tools/validation.go
./internal/tui/app.go
./internal/tui/events_test.go
./internal/tui/events.go
./main.go
./workspace/basic.go
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"./architecture_test.go\n./internal/agent/mcp/config_test.go\n./internal/agent/mcp/config.go\n./internal/agent/mcp/manager_test.go\n./internal/agent/mcp/manager.go\n./internal/agent/mcp/transport_test.go\n./internal/agent/mcp/transport.go\n./internal/agent/runtime/compression/compression_test.go\n./internal/agent/runtime/compression/compression.go\n./internal/agent/runtime/compression/conversation_memory_strategy.go\n./internal/agent/runtime/compression/full_history_summarization.go\n./internal/agent/runtim"}
```

## 2026-06-13T09:07:13Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_00_x6YAlTIElgbPuzBWNiLN1445

### Arguments

```json
{"path": "./internal/agent/runtime/subagent.go"}
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

## 2026-06-13T09:07:13Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_01_h8FbrLi0k67LMAczRlhV0310

### Arguments

```json
{"path": "./internal/agent/runtime/prompt_builder.go"}
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

## 2026-06-13T09:07:13Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_02_T6wz5umORLD4CMp8xQLL7323

### Arguments

```json
{"path": "./internal/agent/runtime/tool_registry.go"}
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

## 2026-06-13T09:07:16Z · read_file · rejected

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_00_w1f8EjFmVdIgzbnViG7o6530

### Arguments

```json
{"path": "./internal/assistant/prompts.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/assistant/prompts.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/assistant/prompts.go: no such file or directory"}
```

## 2026-06-13T09:07:16Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_01_Rn8DsiR9yeTNdlBkkfwu5634

### Arguments

```json
{"path": "./internal/safety/path.go"}
```

### Result

```text
package safety

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func SafePath(p string) (string, error) {
	return SafePathFromRoot("", p)
}

// Contains reports whether the cleaned absolute path resolved is the root
// itself or located within it. Both arguments are expected to be cleaned,
// absolute paths.
func Contains(root, resolved string) bool {
	return resolved == root || strings.HasPrefix(resolved, root+string(os.PathSeparator))
}

func SafePathFromRoot(root, p string) (string, error) {
	workdir, err := os.Getwd()
	if strings.TrimSpace(root) != "" {
		workdir = root
	} else {
		if err != nil {
			return "", fmt.Errorf("failed to get working directory: %w", err)
		}
	}

	workdir, err = filepath.Abs(workdir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve workspace root: %w", err)
	}

	resolvedPath := p
	if !filepath.IsAbs(resolvedPath) {
		resolvedPath = filepath.Join(workdir, resolvedPath)
	}
	resolved, err := filepath.Abs(resolvedPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %w", err)
	}

	workdir = filepath.Clean(workdir)
	resolved = filepath.Clean(resolved)
	if resolved != workdir && !strings.HasPrefix(resolved, workdir+string(os.PathSeparator)) {
		return "", fmt.Errorf("path escapes workspace: %s", p)
	}

	return resolved, nil
}

```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"package safety\n\nimport (\n\t\"fmt\"\n\t\"os\"\n\t\"path/filepath\"\n\t\"strings\"\n)\n\nfunc SafePath(p string) (string, error) {\n\treturn SafePathFromRoot(\"\", p)\n}\n\n// Contains reports whether the cleaned absolute path resolved is the root\n// itself or located within it. Both arguments are expected to be cleaned,\n// absolute paths.\nfunc Contains(root, resolved string) bool {\n\treturn resolved == root || strings.HasPrefix(resolved, root+string(os.PathSeparator))\n}\n\nfunc SafePathFromRoot(root, p string) (string, erro"}
```

## 2026-06-13T09:07:16Z · read_file · rejected

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_02_GapjKar79yH2O3xiZXhw3398

### Arguments

```json
{"path": "./internal/llm/config.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/llm/config.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/llm/config.go: no such file or directory"}
```

## 2026-06-13T09:07:16Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_03_99AtLemvV7zpQMvPGRvD3329

### Arguments

```json
{"path": "./internal/config/config.go"}
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

## 2026-06-13T09:07:19Z · read_file · rejected

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_00_Y9ZQWvHAU8ViJ8vNcoD71248

### Arguments

```json
{"path": "./internal/agent/runtime/hooks/hooks.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/hooks.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/runtime/hooks/hooks.go: no such file or directory"}
```

## 2026-06-13T09:07:19Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_01_yBFreCxJH0mNLN10VgkD4874

### Arguments

```json
{"path": "./internal/agent/runtime/memory.go"}
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

## 2026-06-13T09:07:19Z · read_file · rejected

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_02_GRA18Hj2Bg9mQHe1wCQ93720

### Arguments

```json
{"path": "./internal/agent/storage/storage.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/storage.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/storage.go: no such file or directory"}
```

## 2026-06-13T09:07:19Z · read_file · rejected

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_03_DHlP4U8dER9EXnlWFvlF9937

### Arguments

```json
{"path": "./internal/agent/storage/db.go"}
```

### Result

```text
Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/db.go: no such file or directory
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","denial_reason":"Error: failed to read file: open /Users/bytedance/golang_pro/nano_cc/go-agent/internal/agent/storage/db.go: no such file or directory"}
```

## 2026-06-13T09:07:21Z · bash · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_00_ZNryv45KFAOrXiMusSTE9598

### Arguments

```json
{"command": "find . -type f -name '*.go' -not -path './.git/*' | head -80"}
```

### Result

```text
./workspace/basic.go
./internal/assistant/prompt_test.go
./internal/assistant/prompt.go
./internal/tools/definitions.go
./internal/tools/persisted_output_test.go
./internal/tools/file_ops.go
./internal/tools/handlers.go
./internal/tools/bash.go
./internal/tools/definitions_test.go
./internal/tools/todo_write.go
./internal/tools/persisted_output.go
./internal/tools/load_skill.go
./internal/tools/runtime_env.go
./internal/tools/validation.go
./internal/tools/runtime_test.go
./internal/tools/todo_write_test.go
./internal/tools/path_guard.go
./internal/tools/validation_test.go
./internal/llm/client.go
./internal/llm/errors.go
./internal/logger/logger.go
./internal/logger/logger_test.go
./internal/idgen/idgen.go
./internal/config/paths.go
./internal/config/config.go
./internal/config/local_config.go
./internal/config/local_config_test.go
./internal/config/helpers.go
./internal/tui/events_test.go
./internal/tui/events.go
./internal/tui/app.go
./internal/agent/runtime/memory.go
./internal/agent/runtime/prompt_builder.go
./internal/agent/runtime/context_compression.go
./internal/agent/runtime/subagent.go
./internal/agent/runtime/runtime.go
./internal/agent/runtime/hook_bridge.go
./internal/agent/runtime/tool_registry.go
./internal/agent/runtime/summarizer.go
./internal/agent/runtime/hooks/types.go
./internal/agent/runtime/hooks/user_prompt.go
./internal/agent/runtime/hooks/stop.go
./internal/agent/runtime/hooks/manager.go
./internal/agent/runtime/hooks/tool.go
./internal/agent/runtime/compression/compression_test.go
./internal/agent/runtime/compression/full_history_summarization.go
./internal/agent/runtime/compression/conversation_memory_strategy.go
./internal/agent/runtime/compression/token_estimator.go
./internal/agent/runtime/compression/compression.go
./internal/agent/runtime/compression/message_window_compression.go
./internal/agent/runtime/compression/tool_result_compression.go
./internal/agent/runtime/compression/reactive_compact.go
./internal/agent/runtime/compression/recent_tool_result_retention.go
./internal/agent/runtime/conversation_memory.go
./internal/agent/runtime/memory_test.go
./internal/agent/runtime/context_compression_test.go
./internal/agent/runtime/runtime_test.go
./internal/agent/runtime/helpers.go
./internal/agent/runtime/todo_reminder.go
./internal/agent/runtime/conversation_flow.go
./internal/agent/mcp/transport.go
./internal/agent/mcp/config.go
./internal/agent/mcp/manager_test.go
./internal/agent/mcp/transport_test.go
./internal/agent/mcp/manager.go
./internal/agent/mcp/config_test.go
./internal/agent/storage/models.go
./internal/agent/storage/conversation_history_test.go
./internal/agent/storage/conversation_history.go
./internal/cli/root_test.go
./internal/cli/root.go
./internal/local/store_test.go
./internal/local/store.go
./internal/local/bootstrap.go
./internal/local/persisted_output_files.go
./internal/local/memory_store.go
./internal/local/session_history_store.go
./internal/local/tool_result_log.go
./internal/safety/path.go
./internal/sessions/skill.go
```

### Audit

```json
{"resolved_cwd":"/Users/bytedance/golang_pro/nano_cc/go-agent","outcome_summary":"./workspace/basic.go\n./internal/assistant/prompt_test.go\n./internal/assistant/prompt.go\n./internal/tools/definitions.go\n./internal/tools/persisted_output_test.go\n./internal/tools/file_ops.go\n./internal/tools/handlers.go\n./internal/tools/bash.go\n./internal/tools/definitions_test.go\n./internal/tools/todo_write.go\n./internal/tools/persisted_output.go\n./internal/tools/load_skill.go\n./internal/tools/runtime_env.go\n./internal/tools/validation.go\n./internal/tools/runtime_test.go\n./internal/tools/todo_w"}
```

## 2026-06-13T09:07:24Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_00_IAJTTaRSrfHxp1XYTwUS9500

### Arguments

```json
{"path": "./internal/agent/mcp/manager.go"}
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

## 2026-06-13T09:07:24Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_01_iGsTzSbZM3cAMK1bwFjp6351

### Arguments

```json
{"path": "./internal/agent/storage/models.go"}
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

## 2026-06-13T09:07:24Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_02_qvX6fINZKJXx7PlsJxv87988

### Arguments

```json
{"path": "./internal/agent/storage/conversation_history.go"}
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

## 2026-06-13T09:07:24Z · read_file · success

- conversation_id: conv_0c489398919e8fcd8fb7cd91
- session_id: d8905f1b-6e6b-4686-a0a6-6465a6458fad
- tool_call_id: call_03_Iq6d1BF2d1jtt4L10MYt8035

### Arguments

```json
{"path": "./internal/local/bootstrap.go"}
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

