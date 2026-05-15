package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"nano_cc/internal/config"
	"nano_cc/internal/deploy"
)

type buildConfig struct {
	AppHome          string
	CommandSource    string
	CommandBinDir    string
	CommandScriptDir string
	WorkspaceRoot    string
	BuiltinSkillsDir string
}

func main() {
	appHome := flag.String("app-home", ".", "deployment application root")
	commandSource := flag.String("command-source", "", "source cmd directory (default: <app-home>/cmd)")
	commandBinDir := flag.String("command-bin-dir", "", "binary output directory (default: <app-home>/workspace/bin)")
	commandScriptDir := flag.String("command-script-dir", "", "script asset directory (default: <app-home>/workspace/cmd)")
	goBinary := flag.String("go-binary", "go", "go executable used to build commands")
	flag.Parse()

	built, err := resolveBuildConfig(*appHome, *commandSource, *commandBinDir, *commandScriptDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "resolve build config: %v\n", err)
		os.Exit(1)
	}
	if err := config.EnsureAppLayout(config.AppConfig{
		AppHome:          built.AppHome,
		BuiltinSkillsDir: built.BuiltinSkillsDir,
		CommandBinDir:    built.CommandBinDir,
		CommandScriptDir: built.CommandScriptDir,
		WorkspaceRoot:    built.WorkspaceRoot,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "ensure app layout: %v\n", err)
		os.Exit(1)
	}

	if err := deploy.BuildCommandArtifacts(deploy.BuildOptions{
		AppHome:          built.AppHome,
		CommandSource:    built.CommandSource,
		CommandBinDir:    built.CommandBinDir,
		CommandScriptDir: built.CommandScriptDir,
		GoBinary:         *goBinary,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "build command artifacts: %v\n", err)
		os.Exit(1)
	}
}

func resolveBuildConfig(appHome, commandSource, commandBinDir, commandScriptDir string) (buildConfig, error) {
	resolvedAppHome, err := filepath.Abs(appHome)
	if err != nil {
		return buildConfig{}, err
	}
	resolvedAppHome = filepath.Clean(resolvedAppHome)

	if commandSource == "" {
		commandSource = filepath.Join(resolvedAppHome, "cmd")
	}
	if commandBinDir == "" {
		commandBinDir = filepath.Join(resolvedAppHome, "workspace", "bin")
	}
	if commandScriptDir == "" {
		commandScriptDir = filepath.Join(resolvedAppHome, "workspace", "cmd")
	}

	return buildConfig{
		AppHome:          resolvedAppHome,
		CommandSource:    commandSource,
		CommandBinDir:    commandBinDir,
		CommandScriptDir: commandScriptDir,
		WorkspaceRoot:    filepath.Join(resolvedAppHome, "workspace"),
		BuiltinSkillsDir: filepath.Join(resolvedAppHome, "workspace", "skills"),
	}, nil
}
