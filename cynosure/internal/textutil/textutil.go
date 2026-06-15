package textutil

import (
	"encoding/json"
	"strings"
)

// Truncate returns at most the first max bytes of text.
func Truncate(text string, max int) string {
	if len(text) <= max {
		return text
	}
	return text[:max]
}

// ToolResultPreview builds a short, line- and rune-bounded preview of a tool
// result for display. It returns the preview and whether it was truncated.
func ToolResultPreview(result string) (string, bool) {
	trimmed := strings.TrimSpace(result)
	if trimmed == "" {
		return "(无输出)", false
	}
	lines := strings.Split(trimmed, "\n")
	truncated := false
	if len(lines) > 6 {
		lines = lines[:6]
		truncated = true
	}
	preview := strings.Join(lines, "\n")
	if runes := []rune(preview); len(runes) > 300 {
		preview = string(runes[:300])
		truncated = true
	}
	if truncated {
		preview += "…"
	}
	return preview, truncated
}

// ParseToolResult extracts the status and result from a tool message Content.
// If the content is not valid JSON, isJSON is false and the whole content is
// returned as result.
func ParseToolResult(content string) (status, result string, isJSON bool) {
	var parsed struct {
		Status string `json:"status"`
		Result string `json:"result"`
	}
	if err := json.Unmarshal([]byte(content), &parsed); err == nil {
		return parsed.Status, parsed.Result, true
	}
	return "", content, false
}
