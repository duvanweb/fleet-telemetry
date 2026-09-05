# AI-Assisted Development Report

## Tools Used

- **Claude Sonnet 4.6** via Claude Code CLI (interactive terminal session)
- Session total: ~21 pull requests, covering full monorepo scaffolding through mobile prototype

## Tasks Where AI Was Most Helpful

### 1. Hexagonal Architecture Boilerplate
Generating the full port/adapter/service/module skeleton for each new domain consistently
followed the project conventions (fx.In structs, fx.Annotate, fx.As, fx.Lifecycle hooks)
without drift between services.

### 2. Table-Driven Tests With mockery Mocks
The AI applied the project's strict test pattern (one test per execution branch, reusable
success mocks, inline closure mocks for error paths, AssertExpectations ordering) correctly
after seeing the first example, avoiding manual repetition across 9 test packages.

### 3. RabbitMQ Transactional Outbox
Designing the `CreateWithOutbox` method that BEGIN TX → INSERT telemetry → INSERT outbox →
COMMIT with proper `defer tx.Rollback()` and pq unique constraint error handling required
careful orchestration that the AI produced correctly on the first attempt.

### 4. SSE Hub with Non-Blocking Broadcast
The in-memory fan-out hub using sync.RWMutex + buffered channels (64-slot) + non-blocking
select/default for slow clients was generated correctly without race conditions.

### 5. React Dashboard with SSE-Driven State
Wiring `EventSource` events to a `Map<string, VehicleState>` via `useState` + `useCallback`,
combined with TanStack Query invalidation for alerts, was designed cleanly.

## Real AI Errors and Corrections

### Error 1: Clock Sequencing in Stopped Detection Tests (stopped_detection/service_test.go)

**What happened:** The AI initially wrote the stopped-detection tests using a `callCount`
variable inside the clock closure to return different times based on call order. This produced
incorrect state: after 2 setup events (posA → posB), `FirstSamePositionAt` was still zero
because neither event was at the "same position" relative to the previous one.

**Root cause:** The state machine requires exactly **3 calls at the same position** to trigger
an alert:
1. Call at posA (establishes last known position)
2. Call at posB (samePos=false, updates last position to posB)
3. Call at posB again (samePos=true, FirstSamePositionAt.IsZero()=true → sets FirstSamePositionAt=t0)
4. Call at posB again (samePos=true, now-t0 >= threshold → CreateAlert)

The AI's original design used only 2 setup events before the threshold check event, so
`FirstSamePositionAt` was never set and the test never triggered the alert branch.

**Correction:** Rewrote all test cases to use explicit `clockTimes []time.Time` slices with 4
events: `[posA, posB, posB (sets FirstSamePositionAt=t0), posB (threshold check)]`. The clock
closure reads sequentially from the slice, falling back to the last value when exhausted.

### Error 2: Unused Import in fleet-map.tsx

**What happened:** The AI imported `useEffect` and `useRef` from React in fleet-map.tsx but
never used them in the final implementation.

**Correction:** Removed unused imports; TypeScript compiler caught this immediately.

### Error 3: Simulator Module Wiring

**What happened:** The AI wrote `simulator.go` controller with a duplicate `var json = ...`
declaration, which conflicted with the existing `var json` in `health.go` (same package).

**Correction:** Removed the duplicate declaration. The `json` variable is declared once per
package in `health.go` per project conventions.

## AI Limitations Observed

- **IDE diagnostic noise**: The AI occasionally reported "imported and not used" diagnostics
  from the IDE that did not reflect actual build state — actual `go build ./...` was always
  the authoritative check.
- **State machine reasoning**: Multi-step stateful logic (like the clock sequencing above)
  required human review of execution paths before tests were correct.
- **React Native environment**: Cannot run `npm install` for native modules in a code-only
  environment; type shims were needed as a workaround for CI typecheck.

## Productivity Impact

Estimated 70-80% reduction in time for boilerplate-heavy tasks (FX modules, sqlmock tests,
migration files, DTO/service/controller triples). Domain logic and test path analysis still
required human review on each PR.
