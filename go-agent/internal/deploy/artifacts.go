package deploy

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type BuildOptions struct {
	AppHome        string
	CommandSource  string
	CommandBinDir  string
	CommandScriptDir string
	GoBinary       string
}

func BuildCommandArtifacts(opts BuildOptions) error {
	if opts.GoBinary == "" {
		opts.GoBinary = "go"
	}
	if opts.AppHome == "" {
		return fmt.Errorf("app home is required")
	}
	if opts.CommandSource == "" {
		return fmt.Errorf("command source is required")
	}
	if opts.CommandBinDir == "" {
		return fmt.Errorf("command bin dir is required")
	}
	if opts.CommandScriptDir == "" {
		return fmt.Errorf("command script dir is required")
	}

	if err := os.MkdirAll(opts.CommandBinDir, 0o755); err != nil {
		return fmt.Errorf("mkdir command bin dir: %w", err)
	}
	if err := os.MkdirAll(opts.CommandScriptDir, 0o755); err != nil {
		return fmt.Errorf("mkdir command script dir: %w", err)
	}

	commands, err := DiscoverGoCommands(opts.CommandSource)
	if err != nil {
		return err
	}
	for _, cmdName := range commands {
		pkgDir := filepath.Join(opts.CommandSource, cmdName)
		relPkgDir, err := filepath.Rel(opts.AppHome, pkgDir)
		if err != nil {
			return fmt.Errorf("resolve command package %s: %w", cmdName, err)
		}
		if relPkgDir == ".." || strings.HasPrefix(relPkgDir, ".."+string(os.PathSeparator)) {
			return fmt.Errorf("command source %q must stay under app home %q", pkgDir, opts.AppHome)
		}
		pkgPath := "./" + filepath.ToSlash(relPkgDir)
		output := filepath.Join(opts.CommandBinDir, cmdName)
		cmd := exec.Command(opts.GoBinary, "build", "-o", output, pkgPath)
		cmd.Dir = opts.AppHome
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("build command %s: %w", cmdName, err)
		}
	}

	if err := CopyScriptAssets(opts.CommandSource, opts.CommandScriptDir); err != nil {
		return err
	}

	return nil
}

func DiscoverGoCommands(sourceRoot string) ([]string, error) {
	entries, err := os.ReadDir(sourceRoot)
	if err != nil {
		return nil, fmt.Errorf("read command source: %w", err)
	}
	commands := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		mainPath := filepath.Join(sourceRoot, entry.Name(), "main.go")
		if _, err := os.Stat(mainPath); err == nil {
			commands = append(commands, entry.Name())
		}
	}
	sort.Strings(commands)
	return commands, nil
}

func CopyScriptAssets(sourceRoot, targetRoot string) error {
	return filepath.WalkDir(sourceRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !isScriptAsset(path) {
			return nil
		}

		relPath, err := filepath.Rel(sourceRoot, path)
		if err != nil {
			return fmt.Errorf("resolve script relative path: %w", err)
		}
		destination := filepath.Join(targetRoot, relPath)

		if sameFilePath(path, destination) {
			return nil
		}
		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return fmt.Errorf("mkdir script destination: %w", err)
		}
		if err := copyFile(path, destination); err != nil {
			return fmt.Errorf("copy script asset %s: %w", relPath, err)
		}
		return nil
	})
}

func isScriptAsset(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".py", ".sh", ".rb", ".pl":
		return true
	default:
		return false
	}
}

func sameFilePath(left, right string) bool {
	leftClean := filepath.Clean(left)
	rightClean := filepath.Clean(right)
	return leftClean == rightClean
}

func copyFile(source, destination string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return err
	}

	out, err := os.OpenFile(destination, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return nil
}
