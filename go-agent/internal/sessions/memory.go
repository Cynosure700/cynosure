package sessions

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"nano_cc/internal/safety"
	"nano_cc/internal/tools"
)

// LoadProjectMemory reads AGENTS.md from the working directory.
// Returns empty string if the file does not exist.
func LoadProjectMemory() string {
	data, err := os.ReadFile("AGENTS.md")
	if err != nil {
		return ""
	}
	return string(data)
}

// LoadUserMemory reads AGENTS.md from ~/.link/AGENTS.md.
// Returns empty string if the file or home directory is not accessible.
func LoadUserMemory() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	data, err := os.ReadFile(filepath.Join(home, ".link", "AGENTS.md"))
	if err != nil {
		return ""
	}
	return string(data)
}

// BuildPersistentMemorySection formats project and user memory into XML blocks
// for system prompt injection. Returns empty string if both are empty.
func BuildPersistentMemorySection(projectMemory, userMemory string) string {
	var section string
	if projectMemory != "" {
		section += "<project_memory>\n" + projectMemory + "\n</project_memory>"
	}
	if userMemory != "" {
		if section != "" {
			section += "\n\n"
		}
		section += "<user_memory>\n" + userMemory + "\n</user_memory>"
	}
	return section
}

func handleUpdateMemory(ctx context.Context, args map[string]any) (string, error) {
	action, _ := args["action"].(string)
	if action == "" {
		action = "append"
	}
	content, _ := args["content"].(string)
	if content == "" {
		return "", fmt.Errorf("content is required")
	}

	safePath, err := safety.SafePath("AGENTS.md")
	if err != nil {
		return "", fmt.Errorf("path safety error: %w", err)
	}

	switch action {
	case "append":
		f, err := os.OpenFile(safePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return "", fmt.Errorf("failed to open AGENTS.md: %w", err)
		}
		defer f.Close()
		if _, err := f.WriteString(content + "\n"); err != nil {
			return "", fmt.Errorf("failed to write AGENTS.md: %w", err)
		}
		return "Memory updated: appended to AGENTS.md", nil

	case "replace":
		if err := os.WriteFile(safePath, []byte(content), 0o644); err != nil {
			return "", fmt.Errorf("failed to write AGENTS.md: %w", err)
		}
		return "Memory updated: replaced AGENTS.md", nil

	default:
		return "", fmt.Errorf("invalid action: %s (use append or replace)", action)
	}
}

func init() {
	tools.SetHandler("update_memory", handleUpdateMemory)
}
