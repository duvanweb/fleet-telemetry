package rabbitmq

import (
	"context"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	amqp "github.com/rabbitmq/amqp091-go"

	"fleet/shared/pkg/logger"
)

const fleetExchange = "fleet.events"

var jsonPub = jsoniter.ConfigCompatibleWithStandardLibrary

// messageEnvelope is the standard wrapper for all events published to the exchange.
type messageEnvelope struct {
	EventID    string    `json:"event_id"`
	EventType  string    `json:"event_type"`
	OccurredAt time.Time `json:"occurred_at"`
	Payload    any       `json:"payload"`
}

// Publisher sends AMQP messages to the fleet.events exchange.
type Publisher struct {
	logger  logger.Logger
	channel *amqp.Channel
}

// newPublisher creates a Publisher and declares the exchange.
func newPublisher(conn *amqp.Connection, log logger.Logger) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("opening RabbitMQ channel: %w", err)
	}

	if err := ch.ExchangeDeclare(
		fleetExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("declaring exchange %s: %w", fleetExchange, err)
	}

	return &Publisher{logger: log, channel: ch}, nil
}

// Publish wraps payload in the standard envelope and sends it to the exchange.
func (p *Publisher) Publish(ctx context.Context, eventID, eventType string, payload []byte) error {
	var rawPayload any
	if err := jsonPub.Unmarshal(payload, &rawPayload); err != nil {
		return fmt.Errorf("unmarshalling event payload: %w", err)
	}

	envelope := messageEnvelope{
		EventID:    eventID,
		EventType:  eventType,
		OccurredAt: time.Now().UTC(),
		Payload:    rawPayload,
	}

	body, err := jsonPub.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("marshalling message envelope: %w", err)
	}

	err = p.channel.PublishWithContext(ctx,
		fleetExchange,
		eventType,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("publishing message: %w", err)
	}

	return nil
}
