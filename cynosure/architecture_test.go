package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNoLegacyServicePackageOrImports(t *testing.T) {
	legacyDir := filepath.Join("internal", "web")
	if _, err := os.Stat(legacyDir); !os.IsNotExist(err) {
		t.Fatalf("legacy service package must be removed after TUI migration; stat error: %v", err)
	}

	err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		forbiddenImport := "github.com/Cynosure700/cynosure/cynosure/internal/" + "web/"
		if strings.Contains(string(data), forbiddenImport) {
			t.Fatalf("%s still imports legacy service package", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk source tree: %v", err)
	}
}
