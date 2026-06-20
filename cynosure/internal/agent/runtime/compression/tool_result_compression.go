package compression

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/idgen"
	"nano_cc/internal/textutil"
)

const toolResultCompressionStrategyName = "tool_result_compression"

// ToolResultCompressionStrategy 将最近用户回合中超大的 tool_result 输出持久化，
// 并把它们的内联内容替换为一个 <persisted-output> 标记加一段简短预览。
type ToolResultCompressionStrategy struct{}

func (s *ToolResultCompressionStrategy) Name() string { return toolResultCompressionStrategyName }

// toolResultMessageContent 是 tool 消息所使用的 JSON 包装结构。
type toolResultMessageContent struct {
	Status string `json:"status"`
	Result string `json:"result"`
}

// rebuildToolResult 重新组装 tool 消息内容，当原始内容为 JSON 时保留 JSON 包装与 status。
func rebuildToolResult(status, result string, isJSON bool) string {
	if !isJSON {
		return result
	}
	data, err := json.Marshal(toolResultMessageContent{Status: status, Result: result})
	if err != nil {
		return result
	}
	return string(data)
}

func isCompactedResult(result string) bool {
	return result == earlierToolResultPlaceholder || strings.Contains(result, PersistedOutputMarkerPrefix)
}

// latestUserTurnToolIndexes 返回出现在最后一条用户消息之后的 tool 消息的下标。
func latestUserTurnToolIndexes(history []storage.Message) []int {
	lastUser := -1
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Role == "user" {
			lastUser = i
			break
		}
	}
	var indexes []int
	for i := lastUser + 1; i < len(history); i++ {
		if history[i].Role == "tool" {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

type toolResultCandidate struct {
	index    int
	status   string
	result   string
	isJSON   bool
	bytes    int
	chars    int
	toolName string
}

func (s *ToolResultCompressionStrategy) Apply(ctx context.Context, req *Request) error {
	history := req.RequestHistory
	indexes := latestUserTurnToolIndexes(history)
	if len(indexes) == 0 {
		return nil
	}

	toolNames := toolNamesByCallID(history)
	for _, idx := range indexes {
		status, result, isJSON := textutil.ParseToolResult(history[idx].Content)
		if isCompactedResult(result) {
			continue
		}
		toolName := toolNames[history[idx].ToolCallID]
		chars := len([]rune(result))
		if chars <= req.maxResultSizeChars(toolName) {
			continue
		}
		candidate := toolResultCandidate{index: idx, status: status, result: result, isJSON: isJSON, bytes: len([]byte(result)), chars: chars, toolName: toolName}
		marker, err := s.persistAndBuildMarker(ctx, req, history[candidate.index], candidate)
		if err != nil {
			return err
		}
		history[candidate.index].Content = rebuildToolResult(candidate.status, marker, candidate.isJSON)
	}
	return nil
}

func toolNamesByCallID(history []storage.Message) map[string]string {
	names := make(map[string]string)
	for _, msg := range history {
		if msg.Role != "assistant" {
			continue
		}
		for _, call := range msg.ToolCalls {
			if call.ID == "" {
				continue
			}
			names[call.ID] = call.Function.Name
		}
	}
	return names
}

func (s *ToolResultCompressionStrategy) persistAndBuildMarker(ctx context.Context, req *Request, msg storage.Message, candidate toolResultCandidate) (string, error) {
	sum := sha256.Sum256([]byte(candidate.result))
	contentSHA := hex.EncodeToString(sum[:])

	existing, err := req.Store.GetPersistedOutputByMessageHash(ctx, req.Conversation.ID, req.User.ID, msg.ID, msg.ToolCallID, toolResultCompressionStrategyName, contentSHA)
	var id string
	if err == nil && existing.ID != "" {
		id = existing.ID
	} else {
		id = "po_" + idgen.Hex()
		preview := previewRunes(candidate.result, toolResultPreviewRunes)
		record := storage.PersistedOutput{
			ID:             id,
			ConversationID: req.Conversation.ID,
			UserID:         req.User.ID,
			MessageID:      msg.ID,
			ToolCallID:     msg.ToolCallID,
			Kind:           "tool_result",
			Strategy:       toolResultCompressionStrategyName,
			OriginalBytes:  candidate.bytes,
			ContentSHA256:  contentSHA,
			Content:        candidate.result,
			Preview:        preview,
		}
		if err := req.Store.CreatePersistedOutput(ctx, record); err != nil {
			return "", fmt.Errorf("persist tool output: %w", err)
		}
	}

	preview := previewRunes(candidate.result, toolResultPreviewRunes)
	return buildPersistedOutputMarker(id, candidate.bytes, preview), nil
}

func buildPersistedOutputMarker(id string, originalBytes int, preview string) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<persisted-output id=%q kind="tool_result" original_bytes="%d" preview_chars="%d" retrieval_tool="read_persisted_output">`, id, originalBytes, toolResultPreviewRunes)
	b.WriteString("\n")
	fmt.Fprintf(&b, "完整输出已持久化；如需更多内容，请调用 read_persisted_output(id=%q, offset=0, limit=20000) 分段读取。\n\n", id)
	b.WriteString(preview)
	b.WriteString("\n</persisted-output>")
	return b.String()
}

func previewRunes(text string, max int) string {
	runes := []rune(text)
	if len(runes) <= max {
		return text
	}
	return string(runes[:max])
}
