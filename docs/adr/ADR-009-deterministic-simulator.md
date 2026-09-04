# ADR-009: Deterministic Telemetry Simulator

## Status
Accepted

## Context
End-to-end validation of the pipeline (ingest -> dedup -> outbox -> RabbitMQ -> stopped
detection -> alert -> SSE -> dashboard) requires a controlled data source. Real GPS traces
are non-deterministic; we need scenarios that reliably trigger specific code paths (e.g.
always triggering a VEHICLE_STOPPED alert after a known duration).

## Decision
The simulator service (`apps/simulator`, port 8090) defines five named scenarios with
hardcoded coordinate arrays:

| Scenario | Points | Purpose |
|---|---|---|
| `normal` | 20 waypoints, Bogota city route | Baseline moving vehicle |
| `stopped` | 8 identical coordinates at (7.1193, -73.1227) | Triggers VEHICLE_STOPPED |
| `moving` | 10 diverging waypoints from stopped point | Resolves VEHICLE_STOPPED |
| `duplicate` | 8 waypoints, high duplicate rate | Tests deduplication |
| `invalid` | Mix of valid and out-of-range coords | Tests validation |

`duplicate_rate` and `invalid_rate` (0.0–1.0) are runtime parameters applied per-send
using a seeded `math/rand`. Each vehicle runs in its own goroutine managed by an `errgroup`.
The scenario loops indefinitely until `Stop()` is called, which cancels the context.

## Alternatives Considered

**GPS trace replay**: Replay recorded NMEA files. Realistic but not deterministic — a
specific alert may or may not fire depending on the recorded vehicle's behaviour. Harder
to script test assertions against.

**Random walk**: Generates infinite coordinate variations. Cannot reliably test threshold-
based logic like the 1-minute stopped detection window.

## Consequences

- The `stopped` scenario always emits the same coordinate; with a 10s interval and an
  8-point loop, it triggers the 1-minute stopped threshold reliably after ~70 seconds.
- The simulator does not persist state; restarting resets all vehicle positions.
- Duplicate and invalid rates are applied probabilistically — test assertions must tolerate
  some variance rather than exact counts.
