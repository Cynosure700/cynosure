package runtime

import (
	"context"
	"encoding/json"
	"fmt"

	runtimehooks "nano_cc/internal/web/runtime/hooks"
)

type toolExecutionOutcome = runtimehooks.ToolExecutionOutcome
type toolExecutionAudit = runtimehooks.ToolExecutionAudit

func (s *Service) executeToolCall(ctx context.Context, toolCtx ToolContext, name string, rawArgs string, audit toolExecutionAudit) toolExecutionOutcome {
	if name == "spawn_subagent" {
		if s.Tools == nil || !s.Tools.isAllowed(name) {
			return toolExecutionOutcome{Status: "rejected", Result: "Error: tool spawn_subagent is not registered for web runtime", Audit: audit}
		}
		var args spawnSubagentArgs
		if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
			return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Error: invalid spawn_subagent arguments: %v", err), Audit: audit}
		}
		result, err := s.runSubagent(ctx, toolCtx, args, audit)
		if err != nil {
			return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Subagent failed: %v", err), Audit: audit}
		}
		return toolExecutionOutcome{Status: "success", Result: result, Audit: audit}
	}
	execResult, err := s.Tools.Execute(ctx, toolCtx, name, rawArgs)
	if err != nil {
		return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Error: %v", err), Audit: audit}
	}
	return toolExecutionOutcome{Status: "success", Result: execResult.Output, Audit: audit, Todos: execResult.Todos}
}
