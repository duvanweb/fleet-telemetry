# Architecture Overview

## System Diagram

```
                        ┌─────────────────────────────────────────┐
                        │              Client Layer                │
                        │                                         │
                        │   React Web App          React Native   │
                        │   (Vite + TS)            (Driver App)   │
                        │   localhost:5173                        │
                        └──────┬───────────────────────┬─────────┘
                               │ HTTP / SSE             │ HTTP
                               ▼                        ▼
        ┌──────────────────────────────────────────────────────────────┐
        │                    Microservices Layer                        │
        │                                                              │
        │  ┌─────────────────┐  ┌──────────────────┐  ┌────────────┐  │
        │  │ vehicle-service │  │telemetry-service  │  │alert-      │  │
        │  │  :8081          │  │  :8082            │  │service     │  │
        │  │                 │  │                   │  │  :8083     │  │
        │  │ CRUD vehicles   │  │ Ingest GPS        │  │            │  │
        │  │ vehicle_db      │  │ Dedup (Redis)     │  │ Stopped    │  │
        │  │                 │  │ Outbox -> MQ      │  │ detection  │  │
        │  │                 │  │ SSE /api/events   │  │ alert_db   │  │
        │  └────────┬────────┘  └────────┬──────────┘  └─────┬──────┘  │
        │           │                   │                     │        │
        │  ┌────────┴──────────────────────────────────────┐  │        │
        │  │              Simulator  :8090                 │  │        │
        │  │  5 scenarios, per-vehicle goroutines          │  │        │
        │  └───────────────────────────────────────────────┘  │        │
        └──────────────────────────────────────────────────────────────┘
                               │                     │
                               ▼                     ▼
        ┌──────────────────────────────────────────────────────────────┐
        │                  Infrastructure Layer                        │
        │                                                              │
        │  PostgreSQL (one instance, 4 databases)                      │
        │    vehicle_db  │  telemetry_db  │  alert_db  │  simulator_db │
        │                                                              │
        │  Redis                                                        │
        │    telemetry:dedup:{hash}    TTL 60s                         │
        │    vehicle:{id}:last-position TTL 10m                        │
        │                                                              │
        │  RabbitMQ  exchange: fleet.events (topic, durable)           │
        │    telemetry.received  |  alert.created  |  alert.resolved  │
        └──────────────────────────────────────────────────────────────┘
```

## Data Flow: GPS Ingest to Dashboard

```
GPS Device / Simulator
        │
        │  POST /api/telemetry
        ▼
telemetry-service
  1. Validate lat/lon/timestamp
  2. Check vehicle exists (HTTP → vehicle-service)
  3. Check Redis dedup key → 409 if duplicate
  4. Set Redis dedup key
  5. BEGIN transaction
     INSERT telemetry_points
     INSERT outbox_events (status=pending)
  6. COMMIT
  7. Update Redis last-position (if newer timestamp)
        │
        │  (background, 500ms poll)
        ▼
  OutboxWorker
  8. SELECT pending events FOR UPDATE SKIP LOCKED
  9. Publish to RabbitMQ fleet.events/telemetry.received
  10. UPDATE outbox_events SET status='published'
        │
        ├──────────────────────────────────────┐
        │                                      │
        ▼                                      ▼
alert-service consumer               telemetry-service SSE hub
  11. Decode telemetry event           hub.Broadcast(event)
  12. Evaluate vehicle state                  │
      (same position >= 1 min?)               │  text/event-stream
  13. If stopped → CreateAlert        ◄───────┘
  14. If moved + open alert →
      ResolveAlert                    React Dashboard
  15. Publish alert.created /           ├─ vehicleStates Map updated
      alert.resolved to MQ              ├─ Alerts panel refreshed
                                        └─ FleetMap markers move
```

## Service Ports

| Service | Port | Database |
|---|---|---|
| vehicle-service | 8081 | vehicle_db |
| telemetry-service | 8082 | telemetry_db |
| alert-service | 8083 | alert_db |
| simulator | 8090 | — |
| React web app (dev) | 5173 | — |

## Key Design Principles

- **No API gateway**: Frontend connects directly to each service on its port. CORS is
  enabled on all services (`AllowedOrigins: *`).
- **Database isolation**: No service ever queries another service's database. Cross-service
  data flows only via HTTP calls or RabbitMQ events.
- **Hexagonal architecture**: Domain logic in `internal/core/` depends only on port
  interfaces. Infrastructure adapters in `internal/infrastructure/` implement those interfaces.
- **Transactional outbox**: Guarantees that no telemetry event is silently dropped between
  the DB write and the MQ publish.
