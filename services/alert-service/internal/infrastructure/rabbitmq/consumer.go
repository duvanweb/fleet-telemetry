package rabbitmq

import (
	"context"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	amqp "github.com/rabbitmq/amqp091-go"

	"fleet/alert-service/internal/core/domain"
	"fleet/alert-service/internal/core/ports/services"
	"fleet/shared/pkg/logger"
)

const (
	fleetExchange    = "fleet.events"
	consumerQueue    = "alert-service.telemetry.received"
	telemetryRouting = "telemetry.received"
)

var jsonConsumer = jsoniter.ConfigCompatibleWithStandardLibrary

// inboundMessage is the envelope format published by telemetry-service.
type inboundMessage struct {
	EventID    string                `json:"event_id"`
	EventType  string                `json:"event_type"`
	OccurredAt time.Time             `json:"occurred_at"`
	Payload    domain.TelemetryEvent `json:"payload"`
}

// Consumer subscribes to telemetry.received events and drives stopped-vehicle detection.
type Consumer struct {
	logger    logger.Logger
	channel   *amqp.Channel
	detection services.StoppedDetectionService
}

// newConsumer creates a Consumer, declares the exchange/queue/binding.
func newConsumer(conn *amqp.Connection, log logger.Logger, detection services.StoppedDetectionService) (*Consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("opening RabbitMQ channel: %w", err)
	}

	if err := ch.ExchangeDeclare(fleetExchange, "topic", true, false, false, false, nil); err != nil {
		return nil, fmt.Errorf("declaring exchange: %w", err)
	}

	if _, err := ch.QueueDeclare(consumerQueue, true, false, false, false, nil); err != nil {
		return nil, fmt.Errorf("declaring queue: %w", err)
	}

	if err := ch.QueueBind(consumerQueue, telemetryRouting, fleetExchange, false, nil); err != nil {
		return nil, fmt.Errorf("binding queue: %w", err)
	}

	return &Consumer{logger: log, channel: ch, detection: detection}, nil
}

func (c *Consumer) start(ctx context.Context) {
	deliveries, err := c.channel.Consume(consumerQueue, "alert-service", false, false, false, false, nil)
	if err != nil {
		c.logger.Errorw(ctx, "failed to start consuming", "queue", consumerQueue, "error", err)
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
			c.handle(ctx, delivery)
		}
	}
}

func (c *Consumer) handle(ctx context.Context, delivery amqp.Delivery) {
	var msg inboundMessage
	if err := jsonConsumer.Unmarshal(delivery.Body, &msg); err != nil {
		c.logger.Errorw(ctx, "failed to unmarshal telemetry message", "error", err)
		_ = delivery.Nack(false, false)
		return
	}

	if err := c.detection.Evaluate(ctx, msg.Payload); err != nil {
		c.logger.Errorw(ctx, "failed to evaluate telemetry event", "vehicle_id", msg.Payload.VehicleID, "error", err)
		_ = delivery.Nack(false, true)
		return
	}

	_ = delivery.Ack(false)
}
