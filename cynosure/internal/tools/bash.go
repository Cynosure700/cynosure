package tools

import (
	"fmt"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"
)

var dangerousCommandNames = []string{
	"rm",
	"sudo",
	"shutdown",
	"reboot",
	"halt",
	"poweroff",
	"mkfs",
	"dd",
}

var dangerousSnippets = []string{
	"rm -rf /",
	"> /dev/",
}

const maxOutputLen = 50000
const terminalToolTimeout = 120 * time.Second

func RunBash(command string) (string, error) {
	return RunBashInDir(command, "")
}

func RunBashInDir(command, dir string) (string, error) {
	return RunBashInDirWithOptions(command, dir, false)
}

func RunBashInDirWithOptions(command, dir string, allowDangerous bool) (string, error) {
	if strings.TrimSpace(command) == "" {
		return "", fmt.Errorf("command is required")
	}
	if !allowDangerous {
		if pattern, ok := dangerousCommandPattern(command); ok {
			return "", fmt.Errorf("dangerous command blocked: contains %q", pattern)
		}
	}

	cmd := exec.Command("bash", "-c", command)
	if dir != "" {
		cmd.Dir = dir
	}

	var timedOut atomic.Bool
	timer := time.AfterFunc(terminalToolTimeout, func() {
		timedOut.Store(true)
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	})
	output, err := cmd.CombinedOutput()
	timer.Stop()

	if timedOut.Load() {
		return "", fmt.Errorf("command timed out after %v", terminalToolTimeout)
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

func dangerousCommandPattern(command string) (string, bool) {
	for _, snippet := range dangerousSnippets {
		if strings.Contains(command, snippet) {
			return snippet, true
		}
	}
	fields := splitShellFields(command)
	expectCommand := true
	for _, field := range fields {
		if field == ";" || field == "&" || field == "|" {
			expectCommand = true
			continue
		}
		if !expectCommand {
			continue
		}
		if strings.Contains(field, "=") && !strings.HasPrefix(field, "/") {
			continue
		}
		name := filepathBase(field)
		for _, dangerous := range dangerousCommandNames {
			if name == dangerous {
				return dangerous, true
			}
		}
		expectCommand = false
	}
	return "", false
}

func filepathBase(commandName string) string {
	commandName = strings.TrimSpace(commandName)
	if idx := strings.LastIndex(commandName, "/"); idx >= 0 {
		return commandName[idx+1:]
	}
	return commandName
}
