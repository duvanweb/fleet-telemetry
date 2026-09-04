# Fleet Telemetry & Monitoring System

A distributed fleet telemetry platform built as a monorepo with microservices.

## Architecture

- **vehicle-service** (port 8081) — Vehicle CRUD, owns `vehicle_db`
- **telemetry-service** (port 8082) — GPS ingestion, deduplication, outbox, SSE; owns `telemetry_db`
- **alert-service** (port 8083) — Stopped-vehicle detection, alert lifecycle; owns `alert_db`
- **simulator** (port 8090) — Deterministic telemetry generator
- **web** (port 3000) — React dashboard + simulator UI

No API Gateway — the frontend connects directly to each microservice.

## Quick Start

```bash
# Copy env config
cp .env.example .env

# Start infrastructure (postgres, redis, rabbitmq)
make docker-up

# Run migrations for each service
cd services/vehicle-service && go run ./cmd/migrate/...
cd services/telemetry-service && go run ./cmd/migrate/...
cd services/alert-service && go run ./cmd/migrate/...

# Start all services locally
make docker-all

# Start the web dev server
make web-dev
```

## Development

```bash
# Build all Go modules
make build

# Run all unit tests with race detector
make test

# Lint
make lint

# Format code
make fmt
```

## Stack

**Backend:** Go 1.26, Chi, Uber FX, PostgreSQL, Redis, RabbitMQ
**Frontend:** React 18, TypeScript, Vite, TailwindCSS, shadcn/ui, Leaflet
**Mobile:** React Native, TypeScript
**Infra:** Docker Compose, GitHub Actions

## Ports

| Service | Port |
|---|---|
| vehicle-service | 8081 |
| telemetry-service | 8082 |
| alert-service | 8083 |
| simulator | 8090 |
| web | 3000 |
| PostgreSQL | 5432 |
| Redis | 6379 |
| RabbitMQ AMQP | 5672 |
| RabbitMQ UI | 15672 |
