package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"fleet/shared/pkg/logger"
	"fleet/telemetry-service/internal/infrastructure/sse"
)

const (
	realtimeQueue    = "telemetry-service.realtime"
	realtimeRouteAll = "#"
)

// RealtimeConsumer subscribes to all fleet.events and forwards them to the SSE hub.
type RealtimeConsumer struct {
	logger  logger.Logger
	channel *amqp.Channel
	hub     *sse.Hub
}

// newRealtimeConsumer creates a RealtimeConsumer, declares the fanout queue.
func newRealtimeConsumer(conn *amqp.Connection, log logger.Logger, hub *sse.Hub) (*RealtimeConsumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("opening realtime channel: %w", err)
	}

	if _, err := ch.QueueDeclare(realtimeQueue, false, true, false, false, nil); err != nil {
		return nil, fmt.Errorf("declaring realtime queue: %w", err)
	}

	if err := ch.QueueBind(realtimeQueue, realtimeRouteAll, fleetExchange, false, nil); err != nil {
		return nil, fmt.Errorf("binding realtime queue: %w", err)
	}

	return &RealtimeConsumer{logger: log, channel: ch, hub: hub}, nil
}

func (c *RealtimeConsumer) start(ctx context.Context) {
	deliveries, err := c.channel.Consume(realtimeQueue, "telemetry-service-realtime", true, false, false, false, nil)
	if err != nil {
		c.logger.Errorw(ctx, "failed to start realtime consuming", "queue", realtimeQueue, "error", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case delivery, ok := <-deliveries:
			if !ok {
				return
			}
			c.hub.Broadcast(sse.Event{
				Type:    delivery.RoutingKey,
				Payload: delivery.Body,
			})
		}
	}
}
