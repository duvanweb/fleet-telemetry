package redisclient

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Configuration holds Redis connection settings.
type Configuration struct {
	URL string `env:"REDIS_URL" envDefault:"redis://localhost:6379"`
}

type client struct {
	rdb *redis.Client
}

// NewClient creates and returns a Cache backed by Redis.
func NewClient(cfg *Configuration) (Cache, error) {
	opts, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parsing redis URL: %w", err)
	}

	rdb := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("connecting to redis: %w", err)
	}

	return &client{rdb: rdb}, nil
}

// Delete removes one or more keys.
func (c *client) Delete(ctx context.Context, keys ...string) error {
	return c.rdb.Del(ctx, keys...).Err()
}

// Get retrieves the string value at key.
func (c *client) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

// Set stores value at key with the given TTL.
func (c *client) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return c.rdb.Set(ctx, key, value, ttl).Err()
}

// SetNX sets key to value only if the key does not exist. Returns true if the key was set.
func (c *client) SetNX(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	return c.rdb.SetNX(ctx, key, value, ttl).Result()
}
