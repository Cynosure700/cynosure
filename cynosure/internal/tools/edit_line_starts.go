package tools

import (
	"encoding/json"
	"os"
	"strings"

	"cynosure/internal/safety"
)

// EditFileLineStarts 计算编辑类工具（edit_file/multi_edit）每个文件、每处改动
// 的 new_string 在真实文件中的 1-based 起始行，返回 [文件][改动] 的二维行号。
//
// 它是行号计算的唯一权威实现，供运行时在工具执行后（文件刚写盘、内容最新）
// 计算并持久化，也供 TUI 在缺少持久化值时回退计算。任何无法确定的位置
// （非编辑工具、解析失败、路径越界、文件读取失败、new_string 未命中）置 0，
// 由展示层回退为 1。
//
// 仅用于 TUI diff 行号展示，不参与模型上下文。
func EditFileLineStarts(workspaceRoot, toolName, rawArgs string) [][]int {
	files, ok := editFilesFromRawArgs(toolName, rawArgs)
	if !ok {
		return nil
	}
	starts := make([][]int, len(files))
	for fileIdx, file := range files {
		starts[fileIdx] = make([]int, len(file.newStrings))
		content, ok := readWorkspaceFile(workspaceRoot, file.path)
		if !ok {
			continue
		}
		for changeIdx, newStr := range file.newStrings {
			starts[fileIdx][changeIdx] = lineStartOfSubstring(content, newStr)
		}
	}
	return starts
}

// editLineStartsFile 描述单个被编辑文件及其各处改动的 new_string，
// 顺序与 TUI parseMultiEditDisplayFiles 解析出的文件/改动顺序一致。
type editLineStartsFile struct {
	path       string
	newStrings []string
}

// editFilesFromRawArgs 把 edit_file / multi_edit 的原始参数解析为待定位的
// 文件与每处改动的 new_string；非编辑工具或解析失败返回 ok=false。
func editFilesFromRawArgs(toolName, rawArgs string) ([]editLineStartsFile, bool) {
	switch normalizeToolName(toolName) {
	case "editfile":
		var args struct {
			FilePath string `json:"file_path"`
			NewText  string `json:"new_text"`
		}
		if err := json.Unmarshal([]byte(strings.TrimSpace(rawArgs)), &args); err != nil || strings.TrimSpace(args.FilePath) == "" {
			return nil, false
		}
		return []editLineStartsFile{{path: args.FilePath, newStrings: []string{args.NewText}}}, true
	case "multiedit":
		return multiEditFilesForLineStarts(rawArgs)
	default:
		return nil, false
	}
}

func multiEditFilesForLineStarts(rawArgs string) ([]editLineStartsFile, bool) {
	type editEntry struct {
		NewString string `json:"new_string"`
	}
	type fileEntry struct {
		FilePath string      `json:"file_path"`
		Edits    []editEntry `json:"edits"`
	}
	var args struct {
		FilePath string      `json:"file_path"`
		Edits    []editEntry `json:"edits"`
		Files    []fileEntry `json:"files"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(rawArgs)), &args); err != nil {
		return nil, false
	}
	toFile := func(path string, edits []editEntry) (editLineStartsFile, bool) {
		if strings.TrimSpace(path) == "" || len(edits) == 0 {
			return editLineStartsFile{}, false
		}
		newStrings := make([]string, 0, len(edits))
		for _, edit := range edits {
			newStrings = append(newStrings, edit.NewString)
		}
		return editLineStartsFile{path: path, newStrings: newStrings}, true
	}
	if len(args.Files) > 0 {
		files := make([]editLineStartsFile, 0, len(args.Files))
		for _, file := range args.Files {
			parsed, ok := toFile(file.FilePath, file.Edits)
			if !ok {
				return nil, false
			}
			files = append(files, parsed)
		}
		return files, true
	}
	parsed, ok := toFile(args.FilePath, args.Edits)
	if !ok {
		return nil, false
	}
	return []editLineStartsFile{parsed}, true
}

// readWorkspaceFile 在 workspace 边界内读取文件当前内容；失败返回 ok=false。
func readWorkspaceFile(workspaceRoot, filePath string) (string, bool) {
	resolved, err := safety.SafePathFromRoot(workspaceRoot, filePath)
	if err != nil {
		return "", false
	}
	data, err := os.ReadFile(resolved)
	if err != nil {
		return "", false
	}
	return string(data), true
}

// lineStartOfSubstring 返回 needle 在 content 中首次出现位置的 1-based 起始行，
// 未命中或 needle 为空时返回 0（表示未知）。
func lineStartOfSubstring(content, needle string) int {
	if needle == "" {
		return 0
	}
	idx := strings.Index(content, needle)
	if idx < 0 {
		return 0
	}
	return 1 + strings.Count(content[:idx], "\n")
}

func normalizeToolName(name string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), "_", ""))
}
