package tools

import (
	"fmt"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"
)

const terminalToolTimeout = 120 * time.Second

func RunBash(command string) (string, error) {
	return RunBashInDir(command, "")
}

func RunBashInDir(command, dir string) (string, error) {
	if strings.TrimSpace(command) == "" {
		return "", fmt.Errorf("command is required")
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

	return result, err
}

// filepathBase 返回命令 token 的尾名（去掉路径前缀），供命令分类使用。
func filepathBase(commandName string) string {
	commandName = strings.TrimSpace(commandName)
	if idx := strings.LastIndex(commandName, "/"); idx >= 0 {
		return commandName[idx+1:]
	}
	return commandName
}
