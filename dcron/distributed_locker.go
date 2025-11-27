package dcron

import (
	"context"
	"time"
)

// DistributedLocker 分布式锁接口
type DistributedLocker interface {
	Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error)
	Release(ctx context.Context, key string) error
	Refresh(ctx context.Context, key string, ttl time.Duration) error
}
