package hooks

import "context"

func NewDefaultHookManager() *HookManager {
	return &HookManager{
		UserPromptSubmit: []UserPromptSubmitHook{appendUserMessageHook, conversationActivityHook},
		PreToolUse:       []PreToolUseHook{toolAuditPreHook},
		PostToolUse:      []PostToolUseHook{toolAuditPostHook, persistToolCallHook, emitToolEventHook, appendToolMessageHook},
		Stop:             []StopHook{persistAssistantStopHook, emitAssistantStopHook},
	}
}

func (m *HookManager) RunUserPromptSubmit(ctx context.Context, h *UserPromptSubmitContext) error {
	if m == nil {
		return nil
	}
	for _, hook := range m.UserPromptSubmit {
		if hook == nil {
			continue
		}
		if err := hook(ctx, h); err != nil {
			return err
		}
	}
	return nil
}

func (m *HookManager) RunPreToolUse(ctx context.Context, h *ToolUseContext) error {
	if m == nil {
		return nil
	}
	for _, hook := range m.PreToolUse {
		if hook == nil {
			continue
		}
		if err := hook(ctx, h); err != nil {
			return err
		}
	}
	return nil
}

func (m *HookManager) RunPostToolUse(ctx context.Context, h *ToolUseContext) error {
	if m == nil {
		return nil
	}
	for _, hook := range m.PostToolUse {
		if hook == nil {
			continue
		}
		if err := hook(ctx, h); err != nil {
			return err
		}
	}
	return nil
}

func (m *HookManager) RunStop(ctx context.Context, h *StopContext) error {
	if m == nil {
		return nil
	}
	for _, hook := range m.Stop {
		if hook == nil {
			continue
		}
		if err := hook(ctx, h); err != nil {
			return err
		}
	}
	return nil
}
