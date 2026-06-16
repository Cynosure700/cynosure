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

// ignoredDirNames are directories skipped by grep/glob to avoid noise and
// performance problems when traversing large trees.
var ignoredDirNames = map[string]struct{}{
	".git":         {},
	"node_modules": {},
	"vendor":       {},
}

const defaultSearchHeadLimit = 100

// RunGrepFromRoot searches file contents under a workspace-constrained path
// using a Go regexp. It supports three output modes: files_with_matches
// (default), content and count.
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
	out := strings.Join(lines, "\n")
	if len(out) > maxOutputLen {
		out = out[:maxOutputLen]
	}
	return out
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

// RunGlobFromRoot returns files matching a glob pattern under searchPath,
// sorted by modification time (newest first).
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

// RunLsFromRoot lists immediate entries of a directory, skipping any entry
// whose basename matches an ignore glob.
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

// readTextFile reads a file and reports whether it looks like text (no NUL
// byte). Binary files and unreadable files are skipped (ok=false).
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

// matchGlob matches a relative slash-separated path against a glob pattern
// supporting **, *, ? and [...] segments. ** matches across directory
// separators; * and ? match within a single segment.
func matchGlob(pattern, name string) bool {
	pattern = filepath.ToSlash(pattern)
	name = filepath.ToSlash(name)
	// A pattern without a slash matches against the basename when the
	// candidate carries directory components (e.g. "*.go" vs "a/b/c.go").
	if !strings.Contains(pattern, "/") && strings.Contains(name, "/") {
		name = name[strings.LastIndex(name, "/")+1:]
	}
	pSegs := strings.Split(pattern, "/")
	nSegs := strings.Split(name, "/")
	return matchSegments(pSegs, nSegs)
}

// matchSegments recursively matches pattern segments against name segments.
// A "**" pattern segment matches zero or more name segments.
func matchSegments(pat, name []string) bool {
	if len(pat) == 0 {
		return len(name) == 0
	}
	if pat[0] == "**" {
		// Match zero name segments, or consume one and retry.
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
