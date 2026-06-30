package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"cynosure/internal/local"
)

type textEventWriter struct {
	out io.Writer
}

func (w textEventWriter) Event(name string, data any) error {
	if w.out == nil {
		return nil
	}
	switch name {
	case "assistant_delta":
		if content := eventString(data, "content"); content != "" {
			_, _ = fmt.Fprint(w.out, content)
		}
	case "assistant":
		if content := eventString(data, "content"); content != "" {
			_, _ = fmt.Fprintln(w.out, content)
		}
	case "tool_call_start":
		_, _ = fmt.Fprintf(w.out, "\n[tool:start] %s\n", eventString(data, "name"))
	case "tool_call_done":
		_, _ = fmt.Fprintf(w.out, "[tool:done] %s %s\n", eventString(data, "name"), eventString(data, "status"))
	case "error":
		msg := eventString(data, "message")
		if msg == "" {
			msg = fmt.Sprint(data)
		}
		_, _ = fmt.Fprintf(w.out, "[error] %s\n", msg)
	}
	return nil
}

func eventString(data any, key string) string {
	if m, ok := data.(map[string]any); ok {
		if value, ok := m[key].(string); ok {
			return value
		}
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	var decoded map[string]any
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return ""
	}
	if value, ok := decoded[key].(string); ok {
		return value
	}
	return ""
}

func main() {
	if err := run(context.Background(), os.Args[1:], os.Stdout, os.Stderr); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("cynosure-eval-agent", flag.ContinueOnError)
	fs.SetOutput(stderr)
	cwd := fs.String("cwd", "", "workspace directory")
	prompt := fs.String("prompt", "", "task prompt")
	promptFile := fs.String("prompt-file", "", "path to task prompt file")
	if err := fs.Parse(args); err != nil {
		return err
	}
	promptText := strings.TrimSpace(*prompt)
	if strings.TrimSpace(*promptFile) != "" {
		data, err := os.ReadFile(*promptFile)
		if err != nil {
			return fmt.Errorf("read prompt file %s: %w", *promptFile, err)
		}
		promptText = strings.TrimSpace(string(data))
	}
	if promptText == "" {
		return fmt.Errorf("prompt is required")
	}
	if strings.TrimSpace(*cwd) == "" {
		return fmt.Errorf("cwd is required")
	}
	info, err := os.Stat(*cwd)
	if err != nil {
		return fmt.Errorf("cwd %s: %w", *cwd, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("cwd %s is not a directory", *cwd)
	}

	bundle, err := local.Bootstrap(ctx, *cwd)
	if err != nil {
		return err
	}
	defer bundle.Close()

	_, err = bundle.Runtime.RespondToConversation(ctx, bundle.Conversation, bundle.User, promptText, textEventWriter{out: stdout})
	return err
}
