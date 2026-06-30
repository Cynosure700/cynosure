package main

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
)

func TestRunRequiresCWDAndPrompt(t *testing.T) {
	var stderr bytes.Buffer

	err := run(context.Background(), []string{"--cwd", t.TempDir()}, nil, &stderr)
	if err == nil {
		t.Fatal("run returned nil error without prompt")
	}
	if !strings.Contains(err.Error(), "prompt") {
		t.Fatalf("error = %q, want prompt mention", err.Error())
	}
}

func TestRunRejectsMissingCWD(t *testing.T) {
	var stderr bytes.Buffer

	err := run(context.Background(), []string{"--cwd", "/path/does/not/exist", "--prompt", "hello"}, nil, &stderr)
	if err == nil {
		t.Fatal("run returned nil error for missing cwd")
	}
	if !strings.Contains(err.Error(), "cwd") {
		t.Fatalf("error = %q, want cwd mention", err.Error())
	}
}

func TestRunAcceptsPromptFileForValidation(t *testing.T) {
	tmp := t.TempDir()
	promptFile := tmp + "/prompt.txt"
	if err := os.WriteFile(promptFile, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	var stderr bytes.Buffer

	err := run(context.Background(), []string{"--cwd", "/path/does/not/exist", "--prompt-file", promptFile}, nil, &stderr)
	if err == nil {
		t.Fatal("run returned nil error for missing cwd")
	}
	if strings.Contains(err.Error(), "prompt") {
		t.Fatalf("error = %q, prompt file should satisfy prompt validation", err.Error())
	}
}

func TestTextEventWriterCapturesAssistantAndErrors(t *testing.T) {
	var out bytes.Buffer
	writer := textEventWriter{out: &out}

	if err := writer.Event("assistant_delta", map[string]any{"content": "hello"}); err != nil {
		t.Fatalf("assistant_delta event returned error: %v", err)
	}
	if err := writer.Event("error", map[string]any{"message": "bad"}); err != nil {
		t.Fatalf("error event returned error: %v", err)
	}

	got := out.String()
	for _, want := range []string{"hello", "[error] bad"} {
		if !strings.Contains(got, want) {
			t.Fatalf("output = %q, want %q", got, want)
		}
	}
}
