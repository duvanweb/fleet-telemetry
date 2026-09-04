package stopped_detection

import (
	"time"

	"go.uber.org/fx"

	"fleet/alert-service/internal/core/ports/services"
	"fleet/shared/pkg/env"
	"fleet/shared/pkg/logger"
)

// Configuration holds stopped-detection tuning parameters.
type Configuration struct {
	StoppedThreshold time.Duration `env:"STOPPED_THRESHOLD" envDefault:"1m"`
}

func newService(log logger.Logger, alerts services.AlertService, cfg *Configuration) services.StoppedDetectionService {
	return NewService(log, alerts, cfg.StoppedThreshold, time.Now)
}

// Module wires the stopped detection service into the FX dependency graph.
func Module() fx.Option {
	return fx.Module(
		"stopped_detection",
		fx.Provide(
			env.LoadEnvConfiguration[Configuration],
			newService,
		),
	)
}
