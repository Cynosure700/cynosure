package textutil

import (
	"encoding/json"
)

// Truncate 返回 text 的前 max 个字节（不足时返回全部）。
func Truncate(text string, max int) string {
	if len(text) <= max {
		return text
	}
	return text[:max]
}

// ParseToolResult 从工具消息 Content 中提取 status 与 result。
// 若 content 不是合法 JSON，则 isJSON 为 false，并把整段内容作为 result 返回。
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
