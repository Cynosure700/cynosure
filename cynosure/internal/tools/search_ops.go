package tools

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ignoredDirNames 是 grep/glob 遍历时跳过的目录，用于在遍历大型目录树时
// 避免噪声和性能问题。
var ignoredDirNames = map[string]struct{}{
	".git":         {},
	"node_modules": {},
	"vendor":       {},
}

const defaultSearchHeadLimit = 100

// RunGrepFromRoot 使用 Go 正则在受工作区约束的路径下搜索文件内容。
// 它支持三种输出模式：files_with_matches（默认）、content 和 count。
func RunGrepFromRoot(root, searchPath, pattern, globPattern, outputMode string, ignoreCase, showLineNumbers bool, headLimit int) (string, error) {
	if strings.TrimSpace(pattern) == "" {
		return "", fmt.Errorf("pattern is required")
	}
	if headLimit <= 0 {
		headLimit = defaultSearchHeadLimit
	}
	expr := pattern
	if ignoreCase {
		expr = "(?i)" + expr
	}
	re, err := regexp.Compile(expr)
	if err != nil {
		return "", fmt.Errorf("invalid pattern: %w", err)
	}

	info, err := os.Stat(searchPath)
	if err != nil {
		return "", fmt.Errorf("failed to stat path: %w", err)
	}

	files := make([]string, 0)
	if info.IsDir() {
		err = filepath.WalkDir(searchPath, func(p string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return nil
			}
			if d.IsDir() {
				if _, skip := ignoredDirNames[d.Name()]; skip {
					return filepath.SkipDir
				}
				return nil
			}
			if globPattern != "" && !matchGlob(globPattern, d.Name()) {
				return nil
			}
			files = append(files, p)
			return nil
		})
		if err != nil {
			return "", fmt.Errorf("walk path: %w", err)
		}
	} else {
		files = append(files, searchPath)
	}
	sort.Strings(files)

	switch outputMode {
	case "", "files_with_matches":
		return grepFilesWithMatches(re, root, files, headLimit), nil
	case "content":
		return grepContent(re, root, files, showLineNumbers, headLimit), nil
	case "count":
		return grepCount(re, root, files, headLimit), nil
	default:
		return "", fmt.Errorf("invalid output_mode: %s", outputMode)
	}
}

func grepFilesWithMatches(re *regexp.Regexp, root string, files []string, headLimit int) string {
	matched := make([]string, 0)
	for _, f := range files {
		data, ok := readTextFile(f)
		if !ok {
			continue
		}
		if re.Match(data) {
			matched = append(matched, relDisplayPath(root, f))
			if len(matched) >= headLimit {
				break
			}
		}
	}
	if len(matched) == 0 {
		return "No matches found"
	}
	return strings.Join(matched, "\n")
}

func grepContent(re *regexp.Regexp, root string, files []string, showLineNumbers bool, headLimit int) string {
	lines := make([]string, 0)
	for _, f := range files {
		data, ok := readTextFile(f)
		if !ok {
			continue
		}
		display := relDisplayPath(root, f)
		scanner := bufio.NewScanner(bytes.NewReader(data))
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		lineNo := 0
		for scanner.Scan() {
			lineNo++
			text := scanner.Text()
			if !re.MatchString(text) {
				continue
			}
			if showLineNumbers {
				lines = append(lines, fmt.Sprintf("%s:%d:%s", display, lineNo, text))
			} else {
				lines = append(lines, fmt.Sprintf("%s:%s", display, text))
			}
			if len(lines) >= headLimit {
				goto done
			}
		}
	}
done:
	if len(lines) == 0 {
		return "No matches found"
	}
	return strings.Join(lines, "\n")
}

func grepCount(re *regexp.Regexp, root string, files []string, headLimit int) string {
	lines := make([]string, 0)
	for _, f := range files {
		data, ok := readTextFile(f)
		if !ok {
			continue
		}
		count := len(re.FindAllIndex(data, -1))
		if count == 0 {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s:%d", relDisplayPath(root, f), count))
		if len(lines) >= headLimit {
			break
		}
	}
	if len(lines) == 0 {
		return "No matches found"
	}
	return strings.Join(lines, "\n")
}

// RunGlobFromRoot 返回 searchPath 下匹配 glob 模式的文件，
// 按修改时间排序（最新的在前）。
func RunGlobFromRoot(searchPath, pattern string, headLimit int) (string, error) {
	if strings.TrimSpace(pattern) == "" {
		return "", fmt.Errorf("pattern is required")
	}
	if headLimit <= 0 {
		headLimit = defaultSearchHeadLimit
	}
	type entry struct {
		path    string
		modTime int64
	}
	matches := make([]entry, 0)
	err := filepath.WalkDir(searchPath, func(p string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if d.IsDir() {
			if _, skip := ignoredDirNames[d.Name()]; skip {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(searchPath, p)
		if err != nil {
			return nil
		}
		if !matchGlob(pattern, rel) {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		matches = append(matches, entry{path: p, modTime: info.ModTime().UnixNano()})
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("walk path: %w", err)
	}
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].modTime > matches[j].modTime
	})
	if len(matches) == 0 {
		return "No files found", nil
	}
	truncated := false
	if len(matches) > headLimit {
		matches = matches[:headLimit]
		truncated = true
	}
	paths := make([]string, 0, len(matches))
	for _, m := range matches {
		paths = append(paths, m.path)
	}
	out := strings.Join(paths, "\n")
	if truncated {
		out += fmt.Sprintf("\n... (truncated to %d results)", headLimit)
	}
	return out, nil
}

// RunLsFromRoot 列出某个目录的直接子项，跳过其基名匹配任一 ignore glob 的条目。
func RunLsFromRoot(dirPath string, ignore []string) (string, error) {
	info, err := os.Stat(dirPath)
	if err != nil {
		return "", fmt.Errorf("failed to stat path: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("path is not a directory: %s", dirPath)
	}
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		name := e.Name()
		if matchesAnyGlob(ignore, name) {
			continue
		}
		if e.IsDir() {
			name += "/"
		}
		names = append(names, name)
	}
	sort.Strings(names)
	if len(names) == 0 {
		return "(empty directory)", nil
	}
	return strings.Join(names, "\n"), nil
}

func matchesAnyGlob(patterns []string, name string) bool {
	for _, p := range patterns {
		if strings.TrimSpace(p) == "" {
			continue
		}
		if matchGlob(p, name) {
			return true
		}
	}
	return false
}

// readTextFile 读取文件并判断它是否看起来像文本（不含 NUL 字节）。
// 二进制文件和无法读取的文件会被跳过（ok=false）。
func readTextFile(path string) ([]byte, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	if bytes.IndexByte(data, 0) >= 0 {
		return nil, false
	}
	return data, true
}

func relDisplayPath(root, path string) string {
	if root == "" {
		return path
	}
	if rel, err := filepath.Rel(root, path); err == nil && !strings.HasPrefix(rel, "..") {
		return rel
	}
	return path
}

// matchGlob 将一个以斜杠分隔的相对路径与 glob 模式进行匹配，
// 支持 **、*、? 和 [...] 片段。** 可跨目录分隔符匹配；
// * 和 ? 仅在单个片段内匹配。
func matchGlob(pattern, name string) bool {
	pattern = filepath.ToSlash(pattern)
	name = filepath.ToSlash(name)
	// 当候选路径带有目录部分时，不含斜杠的模式将与其基名进行匹配
	//（例如 "*.go" 对 "a/b/c.go"）。
	if !strings.Contains(pattern, "/") && strings.Contains(name, "/") {
		name = name[strings.LastIndex(name, "/")+1:]
	}
	pSegs := strings.Split(pattern, "/")
	nSegs := strings.Split(name, "/")
	return matchSegments(pSegs, nSegs)
}

// matchSegments 递归地将模式片段与名称片段进行匹配。
// "**" 模式片段可匹配零个或多个名称片段。
func matchSegments(pat, name []string) bool {
	if len(pat) == 0 {
		return len(name) == 0
	}
	if pat[0] == "**" {
		// 匹配零个名称片段，或消费一个名称片段后再重试。
		for i := 0; i <= len(name); i++ {
			if matchSegments(pat[1:], name[i:]) {
				return true
			}
		}
		return false
	}
	if len(name) == 0 {
		return false
	}
	if ok, err := filepath.Match(pat[0], name[0]); err != nil || !ok {
		return false
	}
	return matchSegments(pat[1:], name[1:])
}
