package tools

import "strings"

// approvalRequiredTools 是除 bash 外、调用即需审批的工具集合（写/改类）。
var approvalRequiredTools = map[string]struct{}{
	"write_file": {},
	"edit_file":  {},
	"multi_edit": {},
}

// readOnlyBashCommands 是确定无副作用的查询/搜索类命令名集合，命中即免审批。
// 不含 git（含 commit/push 写子命令）与 sed（含 -i 原地写），它们一律走审批。
var readOnlyBashCommands = map[string]struct{}{
	"ls": {}, "pwd": {}, "echo": {}, "cat": {}, "head": {}, "tail": {},
	"grep": {}, "rg": {}, "fgrep": {}, "egrep": {}, "find": {}, "wc": {},
	"stat": {}, "file": {}, "which": {}, "type": {}, "whoami": {}, "id": {},
	"date": {}, "env": {}, "printenv": {}, "tree": {}, "du": {}, "df": {},
	"basename": {}, "dirname": {}, "realpath": {}, "diff": {}, "cmp": {},
	"sort": {}, "uniq": {}, "cut": {}, "awk": {}, "less": {}, "more": {},
	"column": {}, "uname": {}, "hostname": {},	"cd":{},"sed":{},
}

// RequiresApproval 按工具维度判定该调用是否需要审批，并给出放行规则候选。
//
//   - bash：查询/搜索类命令（ls/grep/cat 等只读命令）免审批；其余命令需审批，
//     rule 为命令首 token（命令名）+ " *"，如 "curl *"；命令名无法解析时回退为 "bash *"。
//   - write_file / edit_file / multi_edit：一律需要审批，rule 为 "<tool> *"。
//   - 其余工具（查询/搜索类专用工具）：无需审批。
func RequiresApproval(toolName string, args map[string]any) (need bool, rule string) {
	if toolName == "bash" {
		command, _ := args["command"].(string)
		if isReadOnlyBashCommand(command) {
			return false, ""
		}
		return true, bashApprovalRule(command)
	}
	if _, ok := approvalRequiredTools[toolName]; ok {
		return true, toolName + " *"
	}
	return false, ""
}

// isReadOnlyBashCommand 判断 bash 命令是否完全由只读命令组成（无需审批）。
// 命令由 ; && | 连接多个子命令时，任一子命令不是只读命令即返回 false；
// 含写重定向（> 或 >>）也视为变更类，返回 false。
func isReadOnlyBashCommand(command string) bool {
	if strings.TrimSpace(command) == "" {
		return false
	}
	if strings.Contains(command, ">") {
		return false
	}
	expectCommand := true
	sawCommand := false
	for _, field := range splitShellFields(command) {
		if field == ";" || field == "&" || field == "|" {
			expectCommand = true
			continue
		}
		if !expectCommand {
			continue
		}
		// 跳过形如 FOO=bar 的前置环境变量赋值。
		if strings.Contains(field, "=") && !strings.HasPrefix(field, "/") {
			continue
		}
		name := filepathBase(field)
		if _, ok := readOnlyBashCommands[name]; !ok {
			return false
		}
		sawCommand = true
		expectCommand = false
	}
	return sawCommand
}

// bashApprovalRule 取命令的首个有效 token（命令名）生成放行规则，如 "curl *"。
func bashApprovalRule(command string) string {
	for _, field := range splitShellFields(command) {
		if field == ";" || field == "&" || field == "|" {
			continue
		}
		// 跳过形如 FOO=bar 的前置环境变量赋值。
		if strings.Contains(field, "=") && !strings.HasPrefix(field, "/") {
			continue
		}
		name := filepathBase(field)
		if strings.TrimSpace(name) == "" {
			continue
		}
		return name + " *"
	}
	return "bash *"
}