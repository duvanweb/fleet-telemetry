# ADR-002: go-chi/chi and go.uber.org/fx

## Status
Accepted

## Context
Each Go service needs an HTTP router and a mechanism to wire dependencies (logger, config,
DB, repositories, services, controllers) at startup without manual constructor chaining.
The wiring must support lifecycle hooks (graceful start/stop) for the HTTP server, RabbitMQ
consumers, and background workers.

## Decision
Use `github.com/go-chi/chi/v5` as the HTTP router and `go.uber.org/fx` as the dependency
injection framework.

- **chi**: idiomatic `net/http` compatibility, URL parameter extraction via `chi.URLParam`,
  middleware composable via `chi.Router.Use`.
- **fx**: constructor-based DI via `fx.Provide`, interface binding via `fx.Annotate` +
  `fx.As`, grouped struct injection via `fx.In`, lifecycle management via `fx.Lifecycle` +
  `fx.Hook{OnStart, OnStop}`.

Each service exposes a top-level `Module() fx.Option` that composes all domain and
infrastructure sub-modules, wired in `cmd/api/module.go`.

## Alternatives Considered

**gorilla/mux**: Older, no active development since 2022. chi is a drop-in replacement
with better middleware and active maintenance.

**google/wire**: Compile-time code generation. More verbose setup; less useful for lifecycle
management and optional/conditional bindings.

**uber/dig**: The underlying DI container for fx. Using it directly loses the lifecycle
and module abstractions fx provides on top.

## Consequences

- All exported constructors must follow the pattern `NewX(log logger.Logger, deps X.Dependencies) *X`.
- `fx.In` structs group dependencies by role (Repositories, Resources, Services).
- Router module always registered last in `fx.Options` to ensure domain modules provide
  their service interfaces before the router tries to inject controllers.
- `fx.Shutdowner` is injected into server lifecycle hooks to trigger clean shutdown on
  `http.ErrServerClosed`.
