package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	cyan   = "\033[36m"
)

var (
	logFile     *os.File
	logFilePath string
	logMu       sync.Mutex
)

func InitFileLogger() error {
	return InitFileLoggerAt("logs")
}

func InitFileLoggerUnderWorkspaceRoot(workspaceRoot string) error {
	base := strings.TrimSpace(workspaceRoot)
	if base == "" {
		return InitFileLogger()
	}
	return InitFileLoggerAt(filepath.Join(base, "logs"))
}

func InitFileLoggerAt(dir string) error {
	logMu.Lock()
	defer logMu.Unlock()

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create log dir: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	path := filepath.Join(dir, fmt.Sprintf("session_%s.log", timestamp))

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create log file: %w", err)
	}
	if logFile != nil {
		_ = logFile.Close()
	}
	logFile = f
	logFilePath = path
	return nil
}

func LogFilePath() string {
	logMu.Lock()
	defer logMu.Unlock()
	return logFilePath
}

func writeDebugLog(level, msg string) {
	logMu.Lock()
	defer logMu.Unlock()
	if logFile == nil {
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(logFile, "[%s] %s %s\n", now, level, msg)
	_ = logFile.Sync()
}

func LogLLMRound(round int, source string, reqBody []byte, respBody []byte, err error) {
	logMu.Lock()
	defer logMu.Unlock()

	if logFile == nil {
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	fmt.Fprintf(logFile, "\n%s\n", strings.Repeat("=", 80))
	fmt.Fprintf(logFile, "[%s] Round %d | Source: %s\n", now, round, source)
	fmt.Fprintf(logFile, "%s\n", strings.Repeat("-", 80))

	// Request
	fmt.Fprintf(logFile, ">>> REQUEST:\n")
	var reqPretty bytes.Buffer
	if json.Indent(&reqPretty, reqBody, "", "  ") == nil {
		fmt.Fprintf(logFile, "%s\n", reqPretty.String())
	} else {
		fmt.Fprintf(logFile, "%s\n", string(reqBody))
	}

	fmt.Fprintf(logFile, "%s\n", strings.Repeat("-", 80))

	// Response
	if err != nil {
		fmt.Fprintf(logFile, "<<< ERROR: %v\n", err)
	} else {
		fmt.Fprintf(logFile, "<<< RESPONSE:\n")
		var respPretty bytes.Buffer
		if json.Indent(&respPretty, respBody, "", "  ") == nil {
			fmt.Fprintf(logFile, "%s\n", respPretty.String())
		} else {
			fmt.Fprintf(logFile, "%s\n", string(respBody))
		}
	}

	fmt.Fprintf(logFile, "%s\n", strings.Repeat("=", 80))
	_ = logFile.Sync()
}

func Info(msg string) {
	fmt.Printf("%s%s%s\n", blue, msg, reset)
	writeDebugLog("INFO", msg)
}

func Warn(msg string) {
	fmt.Printf("%s%s%s\n", yellow, msg, reset)
	writeDebugLog("WARN", msg)
}

func Error(msg string) {
	fmt.Printf("%s%s%s\n", red, msg, reset)
	writeDebugLog("ERROR", msg)
}

func Success(msg string) {
	fmt.Printf("%s%s%s\n", green, msg, reset)
	writeDebugLog("SUCCESS", msg)
}

func Tool(name, msg string) {
	fmt.Printf("%s[%s]%s %s\n", cyan, name, reset, msg)
	writeDebugLog("TOOL", fmt.Sprintf("[%s] %s", name, msg))
}

func Assistant(msg string) {
	fmt.Printf("%sAssistant:%s %s\n", green, reset, msg)
	writeDebugLog("ASSISTANT", msg)
}

func User(msg string) {
	fmt.Printf("%sYou:%s %s\n", blue, reset, msg)
	writeDebugLog("USER", msg)
}
