---
sidebar_position: 2
---

# Application Lifecycle

This page explains what happens between `vef.Run(...)` and a live HTTP server.

## Boot order

The current framework boot sequence is defined in `bootstrap.go`:

`config -> database -> orm -> middleware -> api -> security -> event -> cqrs -> cron -> redis -> mold -> storage -> sequence -> schema -> monitor -> mcp -> app`

That order matters because later modules depend on earlier ones:

- config comes before everything
- database and ORM come before API handlers that need `orm.DB`
- security comes before authenticated API requests
- storage, monitor, schema, and MCP are registered before the app starts listening

## What `vef.Run(...)` actually does

At a high level:

1. builds the framework module list
2. appends your own FX options
3. adds `startApp`
4. creates the FX app
5. runs it

The default start timeout is `30s`. The default stop timeout is `60s`.

Minimal boot example:

```go
func main() {
  vef.Run(
    ivef.Module,
    auth.Module,
    sys.Module,
    web.Module,
  )
}
```

## App startup

The application module creates a Fiber app, applies “before” middleware, mounts the API engine, then applies “after” middleware.

That is why middleware ordering in VEF has two layers:

- app-level middleware order around the whole Fiber app
- API-level middleware order inside the API engine

## App-level middleware order

From the current middleware module, the common order is:

- compression
- headers
- CORS
- content type checks
- request ID
- request logger binding
- panic recovery
- request record logging
- API routes
- MCP endpoint middleware
- storage proxy routes
- SPA fallback middleware

The important consequence is that app middleware can run even for requests that never reach the API engine.

## API-level middleware order

Inside the API engine, the request chain is sorted by middleware order and currently runs as:

- auth
- contextual setup
- data permission resolution
- rate limit
- audit
- handler

This order is what gives handlers access to:

- the authenticated principal
- request-scoped `orm.DB`
- request-scoped logger
- resolved data permission applier

## Startup hooks from built-in modules

Some modules also use lifecycle hooks:

- database pings the connection and logs the database version
- event bus starts its in-memory dispatcher
- storage initializes providers that need startup work
- app starts the HTTP server and registers a stop hook

So a successful boot means more than “Fiber started”. It means the runtime dependencies were initialized too.

## Why this matters when you debug

If a VEF app fails to start, the problem is often one of these:

- config file not found
- invalid database config
- unsupported provider configuration
- constructor registration error in FX
- handler factory resolution failure

Looking at the boot sequence tells you which layer probably failed first.

## Next step

Continue to [Routing](../guide/routing) to see how registered resources become actual HTTP endpoints.
