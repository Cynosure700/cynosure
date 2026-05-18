package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type toolExecutionOutcome struct {
	Status string             `json:"status"`
	Result string             `json:"result"`
	Audit  toolExecutionAudit `json:"-"`
}

type toolExecutionAudit struct {
	ResolvedCWD           string `json:"resolved_cwd,omitempty"`
	ResolvedCommandPath   string `json:"resolved_command_path,omitempty"`
	CommandArtifactPath   string `json:"command_artifact_path,omitempty"`
	CommandArtifactSource string `json:"command_artifact_source,omitempty"`
	OutcomeSummary        string `json:"outcome_summary,omitempty"`
	DenialReason          string `json:"denial_reason,omitempty"`
}

func (s *Service) executeToolCall(ctx context.Context, toolCtx ToolContext, name string, rawArgs string) toolExecutionOutcome {
	runtimeEnv := s.Tools.runtimeEnv()
	resolvedCommandPath, commandArtifactPath := resolveCommandPaths(name, rawArgs, runtimeEnv.CommandBinDir, runtimeEnv.CommandScriptDir)
	audit := toolExecutionAudit{
		ResolvedCWD:           strings.TrimSpace(runtimeEnv.CurrentWorkingDir),
		ResolvedCommandPath:   resolvedCommandPath,
		CommandArtifactPath:   commandArtifactPath,
		CommandArtifactSource: classifyCommandArtifactSource(runtimeEnv.WorkspaceRoot, commandArtifactPath),
	}
	execResult, err := s.Tools.Execute(ctx, toolCtx, name, rawArgs)
	if err != nil {
		audit.DenialReason = err.Error()
		return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Error: %v", err), Audit: audit}
	}
	audit.OutcomeSummary = truncate(execResult.Output, 500)
	return toolExecutionOutcome{Status: "success", Result: execResult.Output, Audit: audit}
}

func (o toolExecutionOutcome) MessageContent() string {
	data, err := json.Marshal(o)
	if err != nil {
		return fmt.Sprintf(`{"status":%q,"result":%q}`, o.Status, o.Result)
	}
	return string(data)
}

func (o toolExecutionOutcome) AuditSummary() string {
	data, err := json.Marshal(o.Audit)
	if err != nil {
		return fmt.Sprintf(`{"resolved_cwd":%q,"resolved_command_path":%q,"command_artifact_path":%q,"command_artifact_source":%q,"outcome_summary":%q,"denial_reason":%q}`,
			o.Audit.ResolvedCWD,
			o.Audit.ResolvedCommandPath,
			o.Audit.CommandArtifactPath,
			o.Audit.CommandArtifactSource,
			o.Audit.OutcomeSummary,
			o.Audit.DenialReason,
		)
	}
	return string(data)
}
