package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"

	"fleet/shared/pkg/env"
	"fleet/shared/pkg/logger"
)

func newConnection(cfg *Configuration) (*amqp.Connection, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("connecting to RabbitMQ: %w", err)
	}
	return conn, nil
}

func registerConnectionHooks(lc fx.Lifecycle, conn *amqp.Connection, log logger.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			if err := conn.Close(); err != nil {
				log.Warnw(context.Background(), "error closing RabbitMQ connection", "error", err)
			}
			return nil
		},
	})
}

func registerWorkerHooks(lc fx.Lifecycle, worker *OutboxWorker) {
	ctx, cancel := context.WithCancel(context.Background())
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go worker.start(ctx)
			return nil
		},
		OnStop: func(_ context.Context) error {
			cancel()
			return nil
		},
	})
}

// Module wires the RabbitMQ connection, publisher, and outbox worker into the FX dependency graph.
func Module() fx.Option {
	return fx.Module(
		"rabbitmq",
		fx.Provide(
			env.LoadEnvConfiguration[Configuration],
			newConnection,
			newPublisher,
			newOutboxWorker,
		),
		fx.Invoke(
			registerConnectionHooks,
			registerWorkerHooks,
		),
	)
}
