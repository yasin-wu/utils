package dcron

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisLocker struct {
	prefix string
	client *redis.Client
}

func NewRedisLocker(prefix string, client *redis.Client) *RedisLocker {
	return &RedisLocker{
		client: client,
		prefix: prefix,
	}
}

func (r *RedisLocker) buildKey(key string) string {
	return fmt.Sprintf("%s:%s", r.prefix, key)
}

func (r *RedisLocker) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	actualKey := r.buildKey(key)
	result, err := r.client.SetNX(ctx, actualKey, "1", ttl).Result()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}
	return result, nil
}

func (r *RedisLocker) Release(ctx context.Context, key string) error {
	actualKey := r.buildKey(key)
	_, err := r.client.Del(ctx, actualKey).Result()
	if err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}
	return nil
}

func (r *RedisLocker) Refresh(ctx context.Context, key string, ttl time.Duration) error {
	actualKey := r.buildKey(key)
	_, err := r.client.Expire(ctx, actualKey, ttl).Result()
	if err != nil {
		return fmt.Errorf("failed to refresh lock: %w", err)
	}
	return nil
}
