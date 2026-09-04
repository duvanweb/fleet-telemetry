# Fleet Telemetry & Monitoring System — Specification

See the full specification in `specs/01-fleet-telemetry-system.md` in the parent review directory,
or refer to the CLAUDE.md plan for the implementation roadmap.

## Summary

A system that:
- Receives GPS coordinates from multiple vehicles via HTTP
- Validates and deduplicates telemetry
- Maintains last position in Redis
- Persists history in PostgreSQL
- Detects stopped vehicles (>1 min same position) and generates alerts
- Publishes events via RabbitMQ (transactional outbox)
- Streams realtime updates to browsers via SSE
- Provides a React dashboard with map visualization
- Includes a telemetry simulator with deterministic scenarios
- Includes a React Native driver app (offline-first)
