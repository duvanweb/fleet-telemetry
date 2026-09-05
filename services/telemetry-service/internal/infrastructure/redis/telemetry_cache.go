package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"fleet/shared/pkg/redisclient"
	"fleet/telemetry-service/internal/core/domain"
)

type lastPositionEntry struct {
	Latitude        float64   `json:"lat"`
	Longitude       float64   `json:"lon"`
	DeviceTimestamp time.Time `json:"ts"`
}

// TelemetryCache implements resources.TelemetryCache using Redis.
type TelemetryCache struct {
	cache redisclient.Cache
}

// NewTelemetryCache creates and returns a new TelemetryCache.
func NewTelemetryCache(cache redisclient.Cache) *TelemetryCache {
	return &TelemetryCache{cache: cache}
}

// CheckDedup returns true if the deduplication key exists in Redis.
func (c *TelemetryCache) CheckDedup(ctx context.Context, key string) (bool, error) {
	_, err := c.cache.Get(ctx, dedupKey(key))
	if errors.Is(err, goredis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("checking dedup key: %w", err)
	}
	return true, nil
}

// SetDedup stores the deduplication key with the given TTL.
func (c *TelemetryCache) SetDedup(ctx context.Context, key string, ttl time.Duration) error {
	if err := c.cache.Set(ctx, dedupKey(key), "1", ttl); err != nil {
		return fmt.Errorf("setting dedup key: %w", err)
	}
	return nil
}

// GetLastPosition retrieves the last known position for the given vehicle.
func (c *TelemetryCache) GetLastPosition(ctx context.Context, vehicleID string) (domain.TelemetryPoint, bool, error) {
	raw, err := c.cache.Get(ctx, lastPositionKey(vehicleID))
	if errors.Is(err, goredis.Nil) {
		return domain.TelemetryPoint{}, false, nil
	}
	if err != nil {
		return domain.TelemetryPoint{}, false, fmt.Errorf("getting last position: %w", err)
	}

	var entry lastPositionEntry
	if err := json.Unmarshal([]byte(raw), &entry); err != nil {
		return domain.TelemetryPoint{}, false, fmt.Errorf("unmarshalling last position: %w", err)
	}

	return domain.TelemetryPoint{
		VehicleID:       vehicleID,
		Latitude:        entry.Latitude,
		Longitude:       entry.Longitude,
		DeviceTimestamp: entry.DeviceTimestamp,
	}, true, nil
}

// SetLastPosition stores the given point as the last known position for the vehicle.
func (c *TelemetryCache) SetLastPosition(ctx context.Context, vehicleID string, point domain.TelemetryPoint, ttl time.Duration) error {
	entry := lastPositionEntry{
		Latitude:        point.Latitude,
		Longitude:       point.Longitude,
		DeviceTimestamp: point.DeviceTimestamp,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshalling last position: %w", err)
	}

	if err := c.cache.Set(ctx, lastPositionKey(vehicleID), string(data), ttl); err != nil {
		return fmt.Errorf("setting last position: %w", err)
	}

	return nil
}

func dedupKey(hash string) string {
	return fmt.Sprintf("telemetry:dedup:%s", hash)
}

func lastPositionKey(vehicleID string) string {
	return fmt.Sprintf("vehicle:%s:last-position", vehicleID)
}
