# ADR-007: Transactional Outbox Pattern

## Status
Accepted

## Context
When a telemetry point is successfully ingested, a `telemetry.received` event must be
published to RabbitMQ so the alert-service can run stopped vehicle detection. A naive
dual-write (insert to DB, then publish to MQ) is not atomic: if the process crashes
between the two operations, the event is lost and the alert pipeline misses the point.

## Decision
Use the transactional outbox pattern:

1. `IngestTelemetry` opens a PostgreSQL transaction.
2. Inserts the `telemetry_points` row.
3. Inserts an `outbox_events` row with `status = 'pending'` in the same transaction.
4. Commits. Either both rows land or neither does.

A background worker (`OutboxWorker`) polls every 500ms:
```sql
SELECT id, event_type, payload FROM outbox_events
WHERE status = 'pending'
ORDER BY created_at ASC
LIMIT 50
FOR UPDATE SKIP LOCKED
```
For each row, it publishes to RabbitMQ and marks the row `published`.

`FOR UPDATE SKIP LOCKED` ensures multiple worker instances (in different replicas) do not
double-process the same event.

## Alternatives Considered

**Dual write without atomicity**: Simple but risks lost events on crash or MQ unavailability.
Unacceptable for the alert pipeline.

**Kafka Connect / Debezium CDC**: Robust but requires running an additional Kafka ecosystem.
The outbox pattern achieves the same guarantee with only PostgreSQL.

**Synchronous publish in the HTTP handler**: Blocks the HTTP response on MQ availability.
Any MQ hiccup adds latency to telemetry ingestion.

## Consequences

- The `outbox_events` table grows without a cleanup job. A periodic DELETE for rows older
  than 24h with `status = 'published'` should be added before production use.
- Publishing is at-least-once. Consumers (alert-service) must be idempotent.
- Worker retry uses exponential backoff from `shared/pkg/retry` and the gobreaker circuit
  breaker from `shared/pkg/circuitbreaker` to avoid hammering a recovering MQ.
