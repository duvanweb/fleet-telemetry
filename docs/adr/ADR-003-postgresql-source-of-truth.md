# ADR-003: PostgreSQL as Source of Truth

## Status
Accepted

## Context
Telemetry points and alerts require durable, ACID-compliant storage. The transactional
outbox pattern (ADR-007) requires atomic writes across two tables in the same transaction,
which mandates a relational database. Schema evolution must be versioned and reversible.

## Decision
Use PostgreSQL via the standard `database/sql` package with the `github.com/lib/pq` driver.
Schema changes are managed with `github.com/golang-migrate/migrate/v4` using sequential
numbered SQL files (`000001_*.up.sql` / `000001_*.down.sql`) under `migrations/` per service.

No ORM is used. All queries are SQL string constants defined in `sql/queries.go` files
within each repository package. Queries are referenced via `regexp.QuoteMeta` in sqlmock
tests to avoid regex escaping issues.

Connection pool settings per service:
```go
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(5 * time.Minute)
db.SetConnMaxIdleTime(2 * time.Minute)
```

## Alternatives Considered

**MongoDB**: Document model fits telemetry points but lacks multi-document ACID transactions
needed for the outbox pattern without using change streams.

**SQLite**: Adequate for local development but not suitable for concurrent writes from
multiple service instances or horizontal scaling.

**GORM / sqlx**: ORMs add reflection overhead and hide query plans. Raw `database/sql`
keeps queries explicit and testable with `go-sqlmock`.

## Consequences

- Each service owns exactly one database (`vehicle_db`, `telemetry_db`, `alert_db`,
  `simulator_db`). Cross-service joins are not possible by design.
- Migrations must always include a `.down.sql` that exactly reverses the `.up.sql`.
- All SQL queries use `*Context` variants (`QueryContext`, `ExecContext`) to propagate
  cancellation and timeouts.
- `go-sqlmock` is the only approved mock for `*sql.DB` in unit tests — no mockery mocks
  for database interfaces.
