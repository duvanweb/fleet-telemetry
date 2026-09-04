# ADR-005: RabbitMQ as Event Bus

## Status
Accepted

## Context
Services need to react to events produced by other services without direct coupling.
The telemetry-service must notify the alert-service when telemetry is received so stopped
vehicle detection can run. The alert-service must broadcast alert lifecycle changes so
the SSE hub and the React dashboard can reflect them in realtime.

## Decision
Use RabbitMQ (`github.com/rabbitmq/amqp091-go`) with a single topic exchange named
`fleet.events` (durable). Routing keys follow the pattern `<domain>.<event>`:

| Routing key | Producer | Consumers |
|---|---|---|
| `telemetry.received` | telemetry-service | alert-service, telemetry-service SSE hub |
| `alert.created` | alert-service | telemetry-service SSE hub |
| `alert.resolved` | alert-service | telemetry-service SSE hub |
| `vehicle.updated` | vehicle-service | telemetry-service SSE hub |

Each consumer declares its own queue (e.g. `alert-service.telemetry`) and binds with its
routing key pattern. The SSE realtime consumer uses an auto-delete queue bound to
`fleet.events/#` to receive all event types.

## Alternatives Considered

**Apache Kafka**: Higher throughput and durable log, but operationally heavier for a demo
platform. RabbitMQ is simpler to run locally and in Docker Compose.

**HTTP webhooks**: Direct service-to-service HTTP calls. Creates tight coupling and
requires retry logic at the caller; RabbitMQ provides built-in buffering and retry.

**Database polling**: Simple but adds latency proportional to poll interval and increases
DB load. The transactional outbox (ADR-007) already uses polling; the MQ step removes
the need for consumers to poll the DB directly.

## Consequences

- At-least-once delivery: consumers must be idempotent. `CreateAlert` uses
  `GetOpenByVehicle` as an idempotency guard before inserting.
- Messages are JSON-encoded using the envelope format:
  `{eventId, eventType, occurredAt, correlationId, payload}`.
- If RabbitMQ is unavailable, the outbox worker retries publishing with exponential
  backoff (ADR-012) so no events are lost.
