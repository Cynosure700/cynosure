package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *Store) SetConversationCache(ctx context.Context, conversationID string, messages []Message) error {
	encoded, err := json.Marshal(messages)
	if err != nil {
		return err
	}
	return s.Redis.Set(ctx, conversationCacheKey(conversationID), encoded, 30*time.Minute).Err()
}

func (s *Store) GetConversationCache(ctx context.Context, conversationID string) ([]Message, bool, error) {
	value, err := s.Redis.Get(ctx, conversationCacheKey(conversationID)).Result()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	var messages []Message
	if err := json.Unmarshal([]byte(value), &messages); err != nil {
		return nil, false, err
	}
	return messages, true, nil
}

func sessionRedisKey(sessionID string) string {
	return "session:" + sessionID
}

func conversationCacheKey(conversationID string) string {
	return "conversation-cache:" + conversationID
}
