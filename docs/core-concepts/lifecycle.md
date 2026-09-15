---
sidebar_position: 2
---

# Application Lifecycle

This page explains what happens between `vef.Run(...)` and a live HTTP server.

## Boot order

This is the canonical statement of the VEF boot pipeline. It is assembled in
`bootstrap.go` (`vef.Run`) and `internal/bootmodules/bootmodules.go` (the
`Assemble` function shared by `vef.Run` and the `internal/apptest` test
harness, so the two graphs cannot drift):

`config -> datasource -> middleware -> api -> security -> event -> expression -> js -> cqrs -> cron -> redis -> lock -> mold -> storage -> sequence -> tx_memory -> outbox -> redis_stream -> inbox -> schema -> monitor -> mcp -> push -> app`

`datasource` is a single step: it connects `*sql.DB` (via `internal/database`)
and wraps it into `orm.DB` (via `internal/orm`) in one module — there is no
separate `database` or `orm` boot step. `tx_memory`, `outbox`, `redis_stream`,
and `inbox` are the event transport submodules — the tx_memory transport
module, the outbox transport module, the redis_stream transport module, and
the inbox module — registered after `sequence` and before `schema`. `js` is
the shared JS engine module and `push` is the WebSocket push module.

Note the list order is for readability: FX resolves the actual construction
order from declared dependencies, so what matters is the dependency shape —

- config feeds everything
- datasource feeds API handlers that need `orm.DB`
- security feeds authenticated API requests
- the event transport submodules (tx_memory, outbox, redis-stream, inbox)
  build on the core `event` module
- storage, monitor, schema, MCP, and push are all wired before the app starts
  listening

## What `vef.Run(...)` actually does

`vef.Run(...)` wires the FX app in this order:

1. installs the framework FX logger with `fx.WithLogger(newFxLogger)`
2. adds the internal config module
3. adds the internal datasource module
4. appends every option returned by `bootmodules.Assemble`
5. appends the user-provided `options...`
6. appends `fx.Invoke(cron.StartScheduler)`
7. appends `fx.Invoke(startApp)`
8. appends `fx.StartTimeout(defaultTimeout)`
9. appends `fx.StopTimeout(defaultTimeout*2)`
10. creates the app with `fx.New(opts...)`
11. runs it with `app.Run()`

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

- compression (`-1000`)
- headers (`-900`)
- CORS (`-800`)
- body-encoding (`-750`, decodes an opt-in `X-Body-Encoding` request body, or enforces and applies protected bidirectional body transport when enabled, before the guard/parse, scoped to `/api`)
- content type (`-700`, JSON/multipart guard, scoped to `/api`)
- request ID (`-650`)
- request logger binding (`-600`)
- panic recovery (`-500`)
- request record logging (`-100`)
- API routes
- push WebSocket endpoint (order `450`; the integration inbound gateway sits at `400` when that module is enabled)
- MCP endpoint middleware (order `500`)
- storage proxy routes (order `900`)
- SPA fallback middleware (order `1000`)

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

- datasource pings the primary connection and logs the database version
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

Continue to [Routing](../building-apis/routing) to see how registered resources become actual HTTP endpoints.
