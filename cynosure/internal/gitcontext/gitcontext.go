package gitcontext

import (
	"os/exec"
	"strings"
	"sync/atomic"
	"time"
)

// gitCommandTimeout 单条 Git 命令的看门狗超时。Git 本地只读操作通常毫秒级返回，
// 5s 足以兜底，超时即视为该命令失败并降级。沿用项目既有的 time.AfterFunc 看门狗风格，
// 不使用 context.WithTimeout。
const gitCommandTimeout = 5 * time.Second

// MaxStatusChars 工作区变更状态的最大字符数（按 rune 计），超过则截断。
const MaxStatusChars = 2000

const statusTruncationNotice = "... (truncated because it exceeds 2k characters. If you need more information, run \"git status\" using the bash tool)"

// Status 是对话开始时采集的 Git 快照原始字段。
type Status struct {
	IsRepo        bool   // 是否为 Git 仓库（整段是否渲染的总开关）
	Branch        string // 当前分支
	MainBranch    string // 主分支
	WorkTree      string // 工作区变更状态（git status --short 原始输出，未截断）
	RecentCommits string // 最近提交（git log --oneline -n 5 原始输出）
	UserName      string // Git 用户名
}

// CollectAndFormat 是 Collect+Format 的便捷组合，供启动流程调用。
// 非 Git 仓库或采集失败时返回空串。
func CollectAndFormat(workspaceRoot string) string {
	return Format(Collect(workspaceRoot))
}

// Collect 在 workspaceRoot 下采集 Git 快照。非 Git 仓库返回 IsRepo=false。
// 仓库判定失败（含 Git 未安装、目录非仓库）直接返回 IsRepo=false，不再执行后续命令。
// 其余字段为最佳努力采集：任一命令失败/超时仅令对应字段为空，不影响整段渲染。
func Collect(workspaceRoot string) Status {
	if out, ok := runGitCommand(workspaceRoot, "rev-parse", "--is-inside-work-tree"); !ok || out != "true" {
		return Status{IsRepo: false}
	}

	status := Status{IsRepo: true}
	if out, ok := runGitCommand(workspaceRoot, "branch", "--show-current"); ok {
		status.Branch = out
	}
	if out, ok := runGitCommand(workspaceRoot, "symbolic-ref", "--short", "refs/remotes/origin/HEAD"); ok {
		status.MainBranch = strings.TrimPrefix(out, "origin/")
	}
	if out, ok := runGitCommand(workspaceRoot, "--no-optional-locks", "status", "--short"); ok {
		status.WorkTree = out
	}
	if out, ok := runGitCommand(workspaceRoot, "--no-optional-locks", "log", "--oneline", "-n", "5"); ok {
		status.RecentCommits = out
	}
	if out, ok := runGitCommand(workspaceRoot, "config", "user.name"); ok {
		status.UserName = out
	}
	return status
}

// Format 将快照格式化为注入提示词的文本块（不含外层标签）。非 Git 仓库返回空串。
func Format(s Status) string {
	if !s.IsRepo {
		return ""
	}

	sections := []string{
		"This is the git status at the start of the conversation. Note that this status is a snapshot in time, and will not update during the conversation.",
	}
	if branch := strings.TrimSpace(s.Branch); branch != "" {
		sections = append(sections, "Current branch: "+branch)
	}
	if mainBranch := strings.TrimSpace(s.MainBranch); mainBranch != "" {
		sections = append(sections, "Main branch (you will usually use this for PRs): "+mainBranch)
	}
	if userName := strings.TrimSpace(s.UserName); userName != "" {
		sections = append(sections, "Git user: "+userName)
	}

	workTree := strings.TrimSpace(s.WorkTree)
	if workTree == "" {
		workTree = "(clean)"
	} else {
		workTree = truncateStatus(workTree)
	}
	sections = append(sections, "Status:\n"+workTree)

	if commits := strings.TrimSpace(s.RecentCommits); commits != "" {
		sections = append(sections, "Recent commits:\n"+commits)
	}

	return strings.Join(sections, "\n\n")
}

// truncateStatus 按 rune 截断工作区状态，超过 MaxStatusChars 时追加截断提示。
func truncateStatus(workTree string) string {
	runes := []rune(workTree)
	if len(runes) <= MaxStatusChars {
		return workTree
	}
	return string(runes[:MaxStatusChars]) + "\n" + statusTruncationNotice
}

// runGitCommand 在 workspaceRoot 下执行只读 Git 命令，返回 trim 后的输出与成功标志。
// 命令带看门狗超时；超时或执行失败返回 (空串, false)。
func runGitCommand(workspaceRoot string, args ...string) (string, bool) {
	cmd := exec.Command("git", args...)
	if strings.TrimSpace(workspaceRoot) != "" {
		cmd.Dir = workspaceRoot
	}

	var timedOut atomic.Bool
	timer := time.AfterFunc(gitCommandTimeout, func() {
		timedOut.Store(true)
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	})
	output, err := cmd.Output()
	timer.Stop()

	if timedOut.Load() || err != nil {
		return "", false
	}
	return strings.TrimSpace(string(output)), true
}
