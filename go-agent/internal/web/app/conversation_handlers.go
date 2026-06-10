package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"nano_cc/internal/logger"
	"nano_cc/internal/web/auth"
	"nano_cc/internal/web/runtime"
	"nano_cc/internal/web/storage"
)

func (s *Server) handleConversations(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromContext(r.Context())
	switch r.Method {
	case http.MethodGet:
		conversations, err := s.store.ListConversationsByUser(r.Context(), user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"conversations": conversations})
	case http.MethodPost:
		var body struct{ Title string }
		_ = json.NewDecoder(r.Body).Decode(&body)
		conversation := storage.Conversation{ID: newID("conv"), UserID: user.ID, RootMessageID: newID("msg"), Title: defaultConversationTitle(body.Title)}
		if err := s.store.CreateConversation(r.Context(), conversation); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"conversation": conversation})
	default:
		methodNotAllowed(w)
	}
}

func (s *Server) handleConversationByID(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromContext(r.Context())
	path := strings.TrimPrefix(r.URL.Path, "/api/conversations/")
	parts := strings.Split(path, "/")
	conversationID := parts[0]
	conversation, err := s.store.GetConversationByID(r.Context(), conversationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if conversation.UserID != user.ID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodGet {
		messages, err := s.store.ListMessagesByConversation(r.Context(), conversation.ID, 100)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"conversation": conversation, "messages": displayConversationMessages(messages)})
		return
	}

	if len(parts) == 1 && r.Method == http.MethodDelete {
		if err := s.store.DeleteConversation(r.Context(), conversation.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
		return
	}

	if len(parts) == 1 && r.Method == http.MethodPatch {
		var body struct{ Title string }
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			badRequest(w, err)
			return
		}
		title := strings.TrimSpace(body.Title)
		if title == "" {
			http.Error(w, "conversation title is required", http.StatusBadRequest)
			return
		}
		if err := s.store.UpdateConversationTitle(r.Context(), conversation.ID, title); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		conversation.Title = title
		writeJSON(w, http.StatusOK, map[string]any{"conversation": conversation})
		return
	}

	if len(parts) == 2 && parts[1] == "stream" && r.Method == http.MethodPost {
		var body struct{ Content string }
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			badRequest(w, err)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		writer := runtime.SSEWriter{W: w}
		_ = writer.Event("conversation", map[string]any{"id": conversation.ID, "title": conversation.Title})
		// panic 兜底：任意未捕获 panic 也下发统一兜底文案，避免连接异常断开导致前端无提示。
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error(fmt.Sprintf("conversation stream panic conversation=%s: %v", conversation.ID, rec))
				_ = writer.Event("error", map[string]any{"message": runtime.FallbackErrorMessage})
			}
		}()
		assistant, err := s.runtime.RespondToConversation(r.Context(), conversation, user, body.Content, writer)
		if err != nil {
			// 原始错误仅记录日志，便于排查；对用户统一下发兜底文案。
			logger.Warn(fmt.Sprintf("conversation stream failed conversation=%s user=%s: %v", conversation.ID, user.ID, err))
			// 客户端主动断开（context.Canceled）属正常现象，连接已断，不下发 error 事件。
			if !errors.Is(err, context.Canceled) {
				_ = writer.Event("error", map[string]any{"message": runtime.FallbackErrorMessage})
			}
			return
		}
		_ = writer.Event("done", map[string]any{"message_id": assistant.ID})
		return
	}

	methodNotAllowed(w)
}

func displayConversationMessages(messages []storage.Message) []storage.Message {
	displayMessages := make([]storage.Message, 0, len(messages))
	for _, message := range messages {
		if message.Role == "user" || (message.Role == "assistant" && len(message.ToolCalls) == 0) {
			displayMessages = append(displayMessages, message)
		}
	}
	return displayMessages
}
