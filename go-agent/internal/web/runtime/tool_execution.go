package runtime

import (
	"context"
	"fmt"

	runtimehooks "nano_cc/internal/web/runtime/hooks"
)

type toolExecutionOutcome = runtimehooks.ToolExecutionOutcome
type toolExecutionAudit = runtimehooks.ToolExecutionAudit

func (s *Service) executeToolCall(ctx context.Context, toolCtx ToolContext, name string, rawArgs string, audit toolExecutionAudit) toolExecutionOutcome {
	execResult, err := s.Tools.Execute(ctx, toolCtx, name, rawArgs)
	if err != nil {
		return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Error: %v", err), Audit: audit}
	}
	return toolExecutionOutcome{Status: "success", Result: execResult.Output, Audit: audit}
}
