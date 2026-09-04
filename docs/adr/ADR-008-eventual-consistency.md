# ADR-008: Accepting Eventual Consistency Across Services

## Status
Accepted

## Context
The telemetry-service must verify that a vehicle exists and is not deleted before accepting
GPS data. The alert-service must know the current alert state for a vehicle to avoid
duplicate alerts. Both pieces of information originate in other services (vehicle-service,
alert-service's own DB). Achieving strong consistency would require distributed transactions.

## Decision
Accept eventual consistency in cross-service interactions:

1. **Vehicle existence check**: `telemetry-service` calls `vehicle-service` via HTTP
   (`GET /vehicles/:id`) at ingestion time. If the vehicle was deleted milliseconds before
   the request, the check returns active and the telemetry is accepted — a brief inconsistency
   window that closes on the next request.

2. **Alert state**: `alert-service` maintains its own alert table. The stopped detection
   service holds an in-memory `map[vehicleID]*VehicleState` that reflects the last known
   alert state. On service restart, the map is rebuilt from the DB
   (`GetOpenByVehicle` on first event per vehicle).

No distributed transactions (2PC), no saga orchestration.

## Alternatives Considered

**Saga pattern**: Choreographed compensating transactions. Significantly more complex to
implement and reason about; the failure modes (partial compensation) can be worse than
the inconsistency window we accept here.

**Distributed transactions (2PC)**: Requires a transaction coordinator and locks resources
across services. Kills throughput and is notoriously difficult to operate correctly.

**Event sourcing + CQRS**: Would give a consistent read model but adds substantial
infrastructure and cognitive overhead not justified by this use case.

## Consequences

- A race condition exists: a vehicle deleted while a GPS device is in the middle of a
  transmission may have its last few points accepted. This is acceptable for fleet telemetry
  where the data is historical and the anomaly self-corrects.
- The in-memory vehicle state map is lost on restart. The first telemetry event after
  restart initialises state from DB correctly; no alert is missed.
- The vehicle-service HTTP call adds ~5ms latency to each ingest. A circuit breaker
  (ADR-012) protects against cascading failure if vehicle-service is unavailable.
