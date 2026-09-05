# ADR-004: Redis for Caching and Deduplication

## Status
Accepted

## Context
GPS devices may retransmit the same telemetry point due to network retries. The ingestion
pipeline must reject duplicates quickly (before hitting the database) while also maintaining
the last known position of each vehicle for the SSE realtime feed.

## Decision
Use Redis (`github.com/redis/go-redis/v9`) with two key patterns:

| Key pattern | TTL | Purpose |
|---|---|---|
| `telemetry:dedup:{sha256_hash}` | 60s | Idempotency guard — reject duplicate ingest requests |
| `vehicle:{vehicleID}:last-position` | 10m | Cached last known GPS point for realtime display |

The deduplication hash is derived from:
```
sha256(vehicleId + fmt.Sprintf("%.7f", lat) + fmt.Sprintf("%.7f", lon) + deviceTimestamp.RFC3339Nano)
```

Ingestion flow:
1. Check Redis dedup key — return `ErrDuplicateTelemetry` if present.
2. Set dedup key (optimistic, before DB write).
3. Insert into PostgreSQL (UNIQUE constraint on `deduplication_key` is a second defense).
4. Update last-position cache only if `deviceTimestamp >= lastKnown.DeviceTimestamp` (out-of-order protection).

## Alternatives Considered

**In-memory map**: Fast but not shared across service instances and lost on restart.

**Database unique constraint only**: Works but requires a full DB round-trip for every
duplicate. Redis check is ~0.5ms vs ~5ms for a DB round-trip under load.

**Bloom filter**: Probabilistic and space-efficient but allows false negatives after TTL
expiry; the exact Redis key approach is preferable for correctness.

## Consequences

- A Redis outage causes `CheckDedup` to fail open or closed depending on error handling.
  Current implementation fails open (skips cache, proceeds to DB). The DB unique constraint
  remains as the authoritative deduplication fence.
- TTL of 60s means duplicates arriving after 60s are not caught by Redis and must rely on
  the DB constraint.
- Last-position cache has a 10-minute TTL; vehicles inactive longer than this will return
  stale or empty position data until the next telemetry event.
