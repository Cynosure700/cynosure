package runtime

import (
	runtimehooks "github.com/Cynosure700/cynosure/cynosure/internal/agent/runtime/hooks"
	"github.com/Cynosure700/cynosure/cynosure/internal/agent/storage"
	agenttools "github.com/Cynosure700/cynosure/cynosure/internal/tools"
)

type UserPromptSubmitHook = runtimehooks.UserPromptSubmitHook
type PreToolUseHook = runtimehooks.PreToolUseHook
type PostToolUseHook = runtimehooks.PostToolUseHook
type StopHook = runtimehooks.StopHook
type HookManager = runtimehooks.HookManager
type LoopState = runtimehooks.LoopState
type UserPromptSubmitContext = runtimehooks.UserPromptSubmitContext
type ToolUseContext = runtimehooks.ToolUseContext
type StopContext = runtimehooks.StopContext

func NewDefaultHookManager() *HookManager {
	return runtimehooks.NewDefaultHookManager()
}

func (s *Service) newLoopState(conversation storage.Conversation, user storage.User, userInput string, history []storage.Message, modelHistory []storage.Message, writer EventWriter) *LoopState {
	return &LoopState{
		Store:                        s.Store,
		NewMessageID:                 newMessageID,
		ToolRuntimeEnv:               s.toolRuntimeEnv,
		ShouldInferConversationTitle: shouldInferConversationTitle,
		InferConversationTitle:       inferConversationTitle,
		Conversation:                 conversation,
		User:                         user,
		UserInput:                    userInput,
		History:                      history,
		ModelHistory:                 modelHistory,
		Writer:                       writer,
	}
}

func (s *Service) toolRuntimeEnv() agenttools.RuntimeEnv {
	if s == nil || s.Tools == nil {
		return agenttools.RuntimeEnv{}
	}
	return s.Tools.runtimeEnv()
}
