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

`vef.Run(...)` wires the FX app in this order:

1. installs the framework FX logger with `fx.WithLogger(newFxLogger)`
2. adds the internal config module
3. adds the internal datasource module
4. appends every option returned by `bootmodules.Core()`
5. appends the user-provided `options...`
6. appends `fx.Invoke(startApp)`
7. appends `fx.StartTimeout(defaultTimeout)`
8. appends `fx.StopTimeout(defaultTimeout*2)`
9. creates the app with `fx.New(opts...)`
10. runs it with `app.Run()`

`defaultTimeout` is `30 * time.Second`, so the default start timeout is `30s`
and the default stop timeout is `60s`.

Because user options are appended after `bootmodules.Core()`, application
modules can append group members through helpers such as
`vef.ProvideAPIResource(...)`. Replacing a core-provided singleton usually
requires `vef.Decorate(...)`, `vef.Replace(...)`, or a framework replacement
helper such as `vef.SupplyFileACL(...)`; registering a second plain
`vef.Provide(...)` for the same service is not an override.

Advanced modules can receive `vef.Lifecycle` and call `Lifecycle.Append(...)`
to register an `fx.Hook` directly; `vef.StartHook`, `vef.StopHook`, and
`vef.StartStopHook` are convenience constructors for those hooks.
The exact `Lifecycle.Append` signature is tracked in the public API index.

The internal `startApp` invoke appends the HTTP server lifecycle hook after the
module graph has been constructed. Its `OnStart` waits for
`application.Start()` or the start context timeout, and its `OnStop` calls
`application.Stop()`.

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
