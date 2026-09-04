package redisclient

import (
	"context"
	"time"
)

// Cache defines the minimal Redis operations used across services.
type Cache interface {
	Delete(ctx context.Context, keys ...string) error
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	SetNX(ctx context.Context, key string, value any, ttl time.Duration) (bool, error)
}
