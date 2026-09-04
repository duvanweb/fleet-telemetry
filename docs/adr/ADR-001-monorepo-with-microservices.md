# ADR-001: Monorepo with Go Workspace

## Status
Accepted

## Context
The fleet telemetry platform requires multiple independent services (vehicle management,
telemetry ingestion, alert processing, simulator) that share common infrastructure code
(logger, env config, Redis client, retry, circuit breaker). Each service must be deployable
independently and have its own dependency graph.

## Decision
Use a Go workspace monorepo (`go.work`) with one `go.mod` per service module. Shared code
lives in `fleet/shared` and is consumed by all services via workspace resolution — no
version pinning required within the repo.

Directory layout:
```
go.work
shared/           module: fleet/shared
services/
  vehicle-service/  module: fleet/vehicle-service
  telemetry-service/
  alert-service/
apps/
  simulator/        module: fleet/simulator
  web/              (Node/React)
  mobile/           (React Native)
```

## Alternatives Considered

**Polyrepo**: Each service in its own repository. Eliminates workspace complexity but
makes cross-cutting changes (e.g. updating the shared logger interface) require coordinated
PRs across multiple repos.

**Single binary**: All services compiled into one executable with feature flags. Simpler
deployment, but couples release cycles and prevents independent scaling.

## Consequences

- `go work sync` must run in CI before any `go build` or `go test` step.
- `fleet/shared` changes are immediately visible to all services without version bumps.
- Services never import each other — cross-service communication is HTTP or RabbitMQ only.
- Each service has its own database; no cross-service SQL joins.
