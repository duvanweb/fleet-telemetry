package redisclient

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Configuration holds the Redis connection settings.
type Configuration struct {
	URL string `env:"REDIS_URL" envDefault:"redis://localhost:6379"`
}

type redisCache struct {
	client *redis.Client
}

// NewClient creates and returns a Cache backed by a Redis client.
func NewClient(cfg *Configuration) (Cache, error) {
	opts, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	return &redisCache{client: client}, nil
}

// Delete removes one or more keys from Redis.
func (r *redisCache) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// Get retrieves a string value for the given key.
func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Set stores a value with the given TTL.
func (r *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// SetNX stores a value only if the key does not exist, returning true if the key was set.
func (r *redisCache) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, ttl).Result()
}
