package runtime

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func classifyCommandArtifactSource(workspaceRoot, commandArtifactPath string) string {
	if strings.TrimSpace(commandArtifactPath) == "" {
		return ""
	}
	cleanArtifact := filepath.Clean(commandArtifactPath)
	cleanWorkspace := strings.TrimSpace(workspaceRoot)
	if cleanWorkspace != "" {
		cleanWorkspace = filepath.Clean(cleanWorkspace)
		if cleanArtifact == cleanWorkspace || strings.HasPrefix(cleanArtifact, cleanWorkspace+string(os.PathSeparator)) {
			return "workspace"
		}
	}
	return "custom"
}

func resolveCommandPaths(toolName, rawArgs string, roots ...string) (string, string) {
	if toolName != "bash" {
		return "", ""
	}
	var args map[string]any
	if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
		return "", ""
	}
	command, _ := args["command"].(string)
	for _, token := range strings.Fields(command) {
		candidate := strings.Trim(token, "\"'`;,()[]{}")
		if candidate == "" || !filepath.IsAbs(candidate) {
			continue
		}
		resolved, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}
		cleanResolved := filepath.Clean(resolved)
		for _, root := range roots {
			if root == "" {
				continue
			}
			resolvedRoot, err := filepath.Abs(root)
			if err != nil {
				continue
			}
			cleanRoot := filepath.Clean(resolvedRoot)
			if cleanResolved == cleanRoot || strings.HasPrefix(cleanResolved, cleanRoot+string(os.PathSeparator)) {
				return cleanResolved, cleanResolved
			}
		}
		return cleanResolved, ""
	}
	return "", ""
}
