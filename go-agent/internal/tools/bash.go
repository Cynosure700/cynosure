package tools

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

var dangerousPatterns = []string{
	"rm -rf /",
	"sudo",
	"shutdown",
	"reboot",
	"> /dev/",
}

const maxOutputLen = 50000
const defaultTimeout = 120 * time.Second

func RunBash(command string) (string, error) {
	return RunBashInDir(command, "")
}

func RunBashInDir(command, dir string) (string, error) {
	for _, pattern := range dangerousPatterns {
		if strings.Contains(command, pattern) {
			return "", fmt.Errorf("dangerous command blocked: contains %q", pattern)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	if dir != "" {
		cmd.Dir = dir
	}
	output, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("command timed out after %v", defaultTimeout)
	}

	result := strings.TrimSpace(string(output))
	if result == "" {
		return "(no output)", nil
	}

	if len(result) > maxOutputLen {
		result = result[:maxOutputLen]
	}

	return result, err
}
