package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"cynosure/internal/local"
	"cynosure/internal/tui"
)

type Options struct {
	CWD string
}

type Runner struct {
	RunTUI func(context.Context, Options) error
}

func Run(ctx context.Context, args []string, defaultCWD string, out io.Writer, runner Runner) error {
	if out == nil {
		out = io.Discard
	}
	if runner.RunTUI == nil {
		runner.RunTUI = runTUI
	}
	if len(args) > 0 && args[0] == "tui" {
		args = args[1:]
	}
	if len(args) > 0 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		printHelp(out)
		return nil
	}

	fs := flag.NewFlagSet("cynosure", flag.ContinueOnError)
	fs.SetOutput(out)
	cwd := fs.String("cwd", defaultCWD, "工作区目录，默认是当前目录")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("unknown command: %s", strings.Join(fs.Args(), " "))
	}
	return runner.RunTUI(ctx, Options{CWD: *cwd})
}

func Main() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	return Run(context.Background(), os.Args[1:], cwd, os.Stdout, Runner{})
}

func runTUI(ctx context.Context, opts Options) error {
	bundle, err := local.Bootstrap(ctx, opts.CWD)
	if err != nil {
		return err
	}
	defer bundle.Close()
	return tui.Run(ctx, bundle.Runtime, tui.SessionInfo{
		User:         bundle.User,
		Conversation: bundle.Conversation,
		CWD:          bundle.CWD,
		Resumer:      bundle.Store,
		Skills:       bundle.Skills,
		MCPServers:   bundle.MCPServers,
		SkillCount:   bundle.SkillCount,
		MCPToolCount: bundle.MCPToolCount,
		ModelID:      bundle.Runtime.Cfg.LLM.ModelID,
	})
}

func printHelp(out io.Writer) {
	_, _ = fmt.Fprintln(out, `cynosure - 本地 TUI 代码助手

在任意项目目录下直接运行 cynosure 即可对当前目录启动 agent。

用法：
  cynosure [--cwd <path>]   启动 TUI，默认工作区为当前目录
  cynosure tui [--cwd <path>] 同上
  cynosure help             显示帮助`)
}
