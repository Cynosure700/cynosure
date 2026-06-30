package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Cynosure700/cynosure/cynosure/internal/safety"
)

// validatedWorkspaceRootFromContext 返回工作区根目录，即 Agent 启动时的工作目录。
// 它是所有工具进行路径解析的唯一基准。
func validatedWorkspaceRootFromContext(ctx context.Context) (string, error) {
	root := strings.TrimSpace(workspaceRootFromContext(ctx))
	if root == "" {
		return "", fmt.Errorf("workspace root is required")
	}
	resolved, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolve workspace root: %w", err)
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", fmt.Errorf("workspace root is unavailable: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("workspace root is not a directory")
	}
	return filepath.Clean(resolved), nil
}

// resolvePathFromContext 解析工具入参中的 path：相对路径基于工作区根目录拼接，
// 绝对路径原样使用。返回工作区根目录与解析后的绝对路径。
func resolvePathFromContext(ctx context.Context, path string) (string, string, error) {
	workspaceRoot, err := validatedWorkspaceRootFromContext(ctx)
	if err != nil {
		return "", "", err
	}
	resolvedPath, err := safety.SafePathFromRoot(workspaceRoot, path)
	if err != nil {
		return "", "", err
	}
	return workspaceRoot, resolvedPath, nil
}

func splitShellFields(command string) []string {
	fields := make([]string, 0)
	var current strings.Builder
	inSingle := false
	inDouble := false
	escaped := false
	for _, r := range command {
		if escaped {
			current.WriteRune(r)
			escaped = false
			continue
		}
		if r == '\\' && !inSingle {
			escaped = true
			continue
		}
		switch r {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
				continue
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
				continue
			}
		case ' ', '\t', '\n', '\r':
			if !inSingle && !inDouble {
				appendShellField(&fields, &current)
				continue
			}
		case ';', '&', '|':
			if !inSingle && !inDouble {
				appendShellField(&fields, &current)
				fields = append(fields, string(r))
				continue
			}
		}
		current.WriteRune(r)
	}
	appendShellField(&fields, &current)
	return fields
}

func appendShellField(fields *[]string, current *strings.Builder) {
	if current.Len() == 0 {
		return
	}
	*fields = append(*fields, current.String())
	current.Reset()
}
