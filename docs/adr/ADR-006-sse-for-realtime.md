# ADR-006: Server-Sent Events for Realtime Updates

## Status
Accepted

## Context
The React dashboard needs to display live vehicle positions, status changes, and alert
notifications without polling. The update frequency can be high (every 5 seconds per
vehicle) when the simulator is running. The communication is unidirectional: server pushes
to browser.

## Decision
Expose a Server-Sent Events (SSE) endpoint on telemetry-service:

```
GET /api/events
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive
X-Accel-Buffering: no
```

An in-process `Hub` struct maintains a `map[clientID]chan Event` (protected by
`sync.RWMutex`). A RabbitMQ consumer subscribes to `fleet.events/#` and calls
`hub.Broadcast` for each received message. The SSE handler subscribes to the hub,
streams events to the browser, and unsubscribes when the request context is cancelled.

Non-blocking broadcast: slow clients are skipped with a `select { default: }` to prevent
one slow browser from blocking all others. Each client channel is buffered at 64 slots.

SSE event format:
```
event: telemetry.received
data: {"vehicleId":"...","latitude":7.1,"longitude":-73.1,"timestamp":"..."}

```

## Alternatives Considered

**WebSockets**: Bidirectional, but the dashboard only needs server-to-client push.
SSE is simpler to implement, works over HTTP/1.1, and reconnects automatically in
browsers.

**Long polling**: Each poll creates a new HTTP request. Higher latency and more server
overhead than a persistent SSE connection.

**GraphQL subscriptions**: Adds a GraphQL layer not otherwise used in the project.
Overkill for a straightforward event stream.

## Consequences

- Each open SSE connection holds a goroutine in the handler. A single telemetry-service
  instance can handle hundreds of concurrent connections; horizontal scaling requires a
  shared pub/sub backend (e.g. Redis Pub/Sub) which is not implemented in this version.
- The `X-Accel-Buffering: no` header is required to prevent nginx from buffering the stream.
- The React client uses the native `EventSource` API; no library dependency needed.
