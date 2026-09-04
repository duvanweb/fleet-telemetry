# ADR-010: No Authentication for This Scope

## Status
Accepted (scope-limited)

## Context
The project scope is a fleet telemetry platform demonstrating GPS ingestion, event
propagation, realtime updates, and a React dashboard. Implementing a full auth layer
(token issuance, refresh, RBAC, device provisioning) would double the implementation
effort without adding value to the telemetry pipeline demonstration.

## Decision
No authentication or authorisation on any endpoint. All services allow `*` CORS origins
and accept all requests without credentials. The React frontend communicates directly to
each service port without bearer tokens or session cookies.

## Alternatives Considered

**JWT with a shared secret**: Simple stateless tokens. Would require a token issuance flow
(login endpoint, device enrollment) and token validation middleware on every service.
Adds ~200 lines of infrastructure per service for no pipeline benefit.

**API keys per device**: Suitable for production GPS devices. Requires a key management
service and per-request key lookup, adding latency to the hot path of telemetry ingestion.

**OAuth2 / OIDC**: Industry standard for user-facing applications. Requires an external
identity provider and significantly complicates the local development setup.

## Consequences

- **Not production-ready.** Any unauthenticated user with network access can read all
  vehicle data and inject arbitrary telemetry.
- Acceptable for a controlled demo environment running on localhost or an internal network.
- Adding authentication later is straightforward: add `middleware.Auth` to chi routes and
  inject a token validator into the FX graph. No architectural changes are needed.
- CORS set to `AllowedOrigins: []string{"*"}` must be tightened before any public exposure.
