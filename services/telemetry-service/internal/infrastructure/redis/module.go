package redis

import (
	"go.uber.org/fx"

	"fleet/telemetry-service/internal/core/ports/resources"
	"fleet/shared/pkg/env"
	"fleet/shared/pkg/redisclient"
)

func newTelemetryCache(cache redisclient.Cache) resources.TelemetryCache {
	return NewTelemetryCache(cache)
}

// Module wires the Redis client and TelemetryCache into the FX dependency graph.
func Module() fx.Option {
	return fx.Module(
		"redis",
		fx.Provide(
			env.LoadEnvConfiguration[redisclient.Configuration],
			redisclient.NewClient,
			newTelemetryCache,
		),
	)
}
