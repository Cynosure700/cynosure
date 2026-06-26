package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cynosure/internal/safety"
)

const maxReadLen = 50000

func RunRead(path string, limit int) (string, error) {
	return RunReadFromRoot("", path, 1, limit)
}

func RunReadFromRoot(root, path string, offset, limit int) (string, error) {
	resolved, err := safety.SafePathFromRoot(root, path)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(resolved)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) == 0 {
		return "(empty file)", nil
	}

	lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	if offset < 1 {
		offset = 1
	}
	start := offset - 1
	if start > len(lines) {
		return "", nil
	}
	end := len(lines)
	if limit > 0 && start+limit < end {
		end = start + limit
	}
	numbered := make([]string, 0, end-start)
	for i, line := range lines[start:end] {
		numbered = append(numbered, fmt.Sprintf("%d\t%s", start+i+1, line))
	}

	result := strings.Join(numbered, "\n")
	if len(result) > maxReadLen {
		result = result[:maxReadLen]
	}

	return result, nil
}

func RunWrite(path, content string) (string, error) {
	return RunWriteFromRoot("", path, content)
}

func RunWriteFromRoot(root, path, content string) (string, error) {
	resolved, err := safety.SafePathFromRoot(root, path)
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(resolved)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create parent directory: %w", err)
	}

	if err := os.WriteFile(resolved, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return fmt.Sprintf("Wrote %d bytes to %s", len(content), path), nil
}

func RunEdit(path, oldText, newText string) (string, error) {
	return RunEditFromRoot("", path, oldText, newText)
}

func RunEditFromRoot(root, path, oldText, newText string) (string, error) {
	resolved, err := safety.SafePathFromRoot(root, path)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(resolved)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	content := string(data)
	if !strings.Contains(content, oldText) {
		return "", fmt.Errorf("text not found in %s", path)
	}

	updated := strings.Replace(content, oldText, newText, 1)
	if err := os.WriteFile(resolved, []byte(updated), 0o644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return fmt.Sprintf("Edited %s", path), nil
}

// Edit 是供 RunMultiEditFromRoot 使用的单次查找替换操作。
type Edit struct {
	OldString  string
	NewString  string
	ReplaceAll bool
}

// RunMultiEditFromRoot 以原子方式对单个文件应用多处编辑：所有编辑按顺序
// 应用到内存中的副本上，只有当每一处编辑都成功时才将文件写回。
// 任何一处失败都会中止操作，且不触碰磁盘。
func RunMultiEditFromRoot(root, path string, edits []Edit) (string, error) {
	if len(edits) == 0 {
		return "", fmt.Errorf("edits is required")
	}
	resolved, err := safety.SafePathFromRoot(root, path)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(resolved)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	content := string(data)
	for i, e := range edits {
		updated, err := applyEdit(content, e.OldString, e.NewString, e.ReplaceAll)
		if err != nil {
			return "", fmt.Errorf("edit %d: %w", i+1, err)
		}
		content = updated
	}

	if err := os.WriteFile(resolved, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return fmt.Sprintf("Applied %d edits to %s", len(edits), path), nil
}

// applyEdit 在内存中执行单次查找替换；当 oldText 为空、等于 newText、
// 或在 content 中不存在时返回错误。
func applyEdit(content, oldText, newText string, replaceAll bool) (string, error) {
	if oldText == "" {
		return "", fmt.Errorf("old_string is required")
	}
	if oldText == newText {
		return "", fmt.Errorf("old_string and new_string must differ")
	}
	if !strings.Contains(content, oldText) {
		return "", fmt.Errorf("text not found: %q", oldText)
	}
	if replaceAll {
		return strings.ReplaceAll(content, oldText, newText), nil
	}
	return strings.Replace(content, oldText, newText, 1), nil
}
