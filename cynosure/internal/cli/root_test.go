package cli

import (
	"bytes"
	"context"
	"testing"
)

func TestRunDefaultsToTUI(t *testing.T) {
	var called bool
	runner := Runner{
		RunTUI: func(ctx context.Context, opts Options) error {
			called = true
			if opts.CWD != "/tmp/project" {
				t.Fatalf("CWD = %q, want /tmp/project", opts.CWD)
			}
			return nil
		},
	}

	if err := Run(context.Background(), []string{}, "/tmp/project", &bytes.Buffer{}, runner); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if !called {
		t.Fatal("RunTUI was not called")
	}
}

func TestRunRejectsUnknownCommand(t *testing.T) {
	var called bool
	out := &bytes.Buffer{}
	runner := Runner{RunTUI: func(ctx context.Context, opts Options) error {
		called = true
		return nil
	}}

	err := Run(context.Background(), []string{"serve"}, "/tmp/project", out, runner)
	if err == nil {
		t.Fatal("Run returned nil error for unknown command")
	}
	if called {
		t.Fatal("RunTUI should not be called for unknown command")
	}
}

func TestRunUsesExplicitCWD(t *testing.T) {
	var got string
	runner := Runner{RunTUI: func(ctx context.Context, opts Options) error {
		got = opts.CWD
		return nil
	}}

	if err := Run(context.Background(), []string{"--cwd", "/tmp/other"}, "/tmp/project", &bytes.Buffer{}, runner); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if got != "/tmp/other" {
		t.Fatalf("CWD = %q, want /tmp/other", got)
	}
}
