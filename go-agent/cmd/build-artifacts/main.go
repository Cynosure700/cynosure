package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"nano_cc/internal/config"
	"nano_cc/internal/deploy"
)

func main() {
	appHome := flag.String("app-home", ".", "deployment application root")
	commandSource := flag.String("command-source", "", "source cmd directory (default: <app-home>/cmd)")
	commandBinDir := flag.String("command-bin-dir", "", "binary output directory (default: <app-home>/workspaces/bin)")
	commandScriptDir := flag.String("command-script-dir", "", "script asset directory (default: <app-home>/workspaces/cmd)")
	goBinary := flag.String("go-binary", "go", "go executable used to build commands")
	flag.Parse()

	resolvedAppHome, err := filepath.Abs(*appHome)
	if err != nil {
		fmt.Fprintf(os.Stderr, "resolve app home: %v\n", err)
		os.Exit(1)
	}
	resolvedAppHome = filepath.Clean(resolvedAppHome)

	if *commandSource == "" {
		*commandSource = filepath.Join(resolvedAppHome, "cmd")
	}
	if *commandBinDir == "" {
		*commandBinDir = filepath.Join(resolvedAppHome, "workspaces", "bin")
	}
	if *commandScriptDir == "" {
		*commandScriptDir = filepath.Join(resolvedAppHome, "workspaces", "cmd")
	}
	if err := config.EnsureAppLayout(config.AppConfig{
		AppHome:          resolvedAppHome,
		BuiltinSkillsDir: filepath.Join(resolvedAppHome, "workspaces", "skills"),
		CommandBinDir:    *commandBinDir,
		CommandScriptDir: *commandScriptDir,
		WorkspaceRoot:    filepath.Join(resolvedAppHome, "workspace"),
	}); err != nil {
		fmt.Fprintf(os.Stderr, "ensure app layout: %v\n", err)
		os.Exit(1)
	}

	if err := deploy.BuildCommandArtifacts(deploy.BuildOptions{
		AppHome:          resolvedAppHome,
		CommandSource:    *commandSource,
		CommandBinDir:    *commandBinDir,
		CommandScriptDir: *commandScriptDir,
		GoBinary:         *goBinary,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "build command artifacts: %v\n", err)
		os.Exit(1)
	}
}
