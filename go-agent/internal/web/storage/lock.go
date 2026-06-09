package storage

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// lockRenewScript 续期脚本：仅当锁仍由当前 token 持有时才刷新 TTL。
var lockRenewScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("PEXPIRE", KEYS[1], ARGV[2])
else
    return 0
end`)

// lockReleaseScript 释放脚本：仅当锁仍由当前 token 持有时才删除。
var lockReleaseScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`)

// AcquireConversationLock 阻塞式获取会话锁。锁被占用时按固定间隔轮询重试，
// 直到获取成功、等待超时或 ctx 取消。成功返回 (true, nil)；等待超时返回
// (false, nil)；ctx 取消或 Redis 异常返回 (false, err)。
func (s *Store) AcquireConversationLock(ctx context.Context, conversationID, token string, ttl, waitTimeout time.Duration) (bool, error) {
	key := conversationLockKey(conversationID)
	deadline := time.Now().Add(waitTimeout)
	const pollInterval = 100 * time.Millisecond
	for {
		ok, err := s.Redis.SetNX(ctx, key, token, ttl).Result()
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
		if time.Now().After(deadline) {
			return false, nil
		}
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(pollInterval):
		}
	}
}

// RenewConversationLock 为当前 token 持有的会话锁续期。返回 true 表示续期成功；
// false 表示锁已不属于该 token（已过期或被他人持有）。
func (s *Store) RenewConversationLock(ctx context.Context, conversationID, token string, ttl time.Duration) (bool, error) {
	key := conversationLockKey(conversationID)
	res, err := lockRenewScript.Run(ctx, s.Redis, []string{key}, token, ttl.Milliseconds()).Int64()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

// ReleaseConversationLock 释放当前 token 持有的会话锁；非持有者调用为空操作。
func (s *Store) ReleaseConversationLock(ctx context.Context, conversationID, token string) error {
	key := conversationLockKey(conversationID)
	return lockReleaseScript.Run(ctx, s.Redis, []string{key}, token).Err()
}
