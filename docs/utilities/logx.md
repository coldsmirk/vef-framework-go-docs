---
sidebar_position: 9
---

# logx

The `logx` package defines the structured-logging contract used across the
framework — every framework component that logs (API, security, event bus,
CRUD, approval, and so on) depends on `logx.Logger`, never on a concrete
logging library directly.

## API Reference

| API | Contract | Purpose |
| --- | --- | --- |
| `logx.Level` | `type Level int8` | Logging priority; higher levels are more important. |
| `logx.LevelDebug = 1` | debug level constant | Voluminous logs, usually disabled in production. |
| `logx.LevelInfo = 2` | info level constant | Default logging priority. |
| `logx.LevelWarn = 3` | warning level constant | More important than info but not necessarily human-reviewed one by one. |
| `logx.LevelError = 4` | error level constant | High-priority logs for unexpected application behavior. |
| `logx.LevelPanic = 5` | panic level constant | Logs a message and then panics. |
| `logx.Logger` | logging interface | Contract implemented by framework-provided and custom loggers. |
| `logx.LoggerConfigurable[T]` | generic interface | Implemented by immutable components that return a logger-configured copy from `WithLogger`. |

`Level.String() string` returns `debug`, `info`, `warn`, `error`, or `panic`.
Unknown level values, including the zero value, return `unknown`.

## Logger Interface

| Method | Contract |
| --- | --- |
| `Logger.Named(name string) logx.Logger` | returns a child logger with the given namespace |
| `Logger.WithCallerSkip(skip int) logx.Logger` | returns a logger that adjusts caller stack-frame reporting |
| `Logger.Enabled(level logx.Level) bool` | reports whether the given level is enabled |
| `Logger.Sync()` | flushes buffered log entries; the interface does not return an error |
| `Logger.Debug(message string)` | logs a debug message |
| `Logger.Debugf(template string, args ...any)` | logs a formatted debug message |
| `Logger.Info(message string)` | logs an info message |
| `Logger.Infof(template string, args ...any)` | logs a formatted info message |
| `Logger.Warn(message string)` | logs a warning message |
| `Logger.Warnf(template string, args ...any)` | logs a formatted warning message |
| `Logger.Error(message string)` | logs an error message |
| `Logger.Errorf(template string, args ...any)` | logs a formatted error message |
| `Logger.Panic(message string)` | logs a panic message and then panics |
| `Logger.Panicf(template string, args ...any)` | logs a formatted panic message and then panics |
| `LoggerConfigurable[T].WithLogger(logger logx.Logger) T` | returns a copy of the component configured with the given logger |

`LoggerConfigurable[T]` is the pattern used by immutable, DI-constructed
components that need a logger injected after construction: `WithLogger`
returns a new configured copy rather than mutating the receiver in place.

## Default Implementation

Out of the box, every `logx.Logger` is backed by a `zap.SugaredLogger`
(`go.uber.org/zap`) built inside `internal/logx`:

| Aspect | Behavior |
| --- | --- |
| Encoding | console — single-line, human-readable text; there is no JSON option at v0.37.0 |
| Destination | log entries go to stdout; zap's own internal errors go to stderr |
| Line layout | dimmed `2006-01-02T15:04:05.000` timestamp, capitalized level, `[name]` namespace, trimmed caller path, message |
| Color | the level column is always ANSI-colorized (`zapcore.CapitalColorLevelEncoder`); timestamp, caller, and name styling is applied only when stdout supports it (termenv detection: disabled for non-TTY output and under `NO_COLOR`) |
| Stack traces | disabled |

The root logger is a package-level value constructed once at process start;
every logger the framework hands out — component loggers, `vef.NamedLogger`,
and the request-scoped logger — is a `Named` child of that single root. FX
lifecycle events and config-loader (viper) messages funnel into the same
logger through an internal `log/slog` bridge, filtered to warn and above.

There is no configuration surface for encoding or destination: to ship logs
elsewhere, capture stdout at the deployment level.

## Log Levels

The active threshold comes from the `VEF_LOG_LEVEL` environment variable
(`config.EnvLogLevel`), read once at process start:

| `VEF_LOG_LEVEL` | Threshold |
| --- | --- |
| `debug` | `logx.LevelDebug` |
| `warn` | `logx.LevelWarn` |
| `error` | `logx.LevelError` |
| `panic` | `logx.LevelPanic` |
| anything else, including unset and `info` | `logx.LevelInfo` (the default) |

Matching is case-insensitive. The threshold is a single shared
`zap.AtomicLevel` that governs every named logger at once — there is no
per-component level configuration, and no public runtime API to change the
level after startup (an internal setter exists but is not exported). Use
`Logger.Enabled(level)` to guard expensive log-message construction.

## Request-Scoped Logging

Two built-in app middlewares cooperate to give every HTTP request its own
logger:

1. The request-ID middleware (order `-650`) reuses a valid incoming
   `X-Request-ID` header, or generates a UUID v7 via `id.GenerateUUID`, and
   echoes the value in the response header.
2. The logger middleware (order `-600`) derives a child logger named
   `request_id:<id>` from the root logger and stores it with
   `contextx.SetLogger` — into both the fiber locals and the embedded
   `context.Context`, so fiber handlers and plain-`context.Context` consumers
   resolve the same logger.

Downstream, the request logger travels three ways:

- `contextx.Logger(ctx, fallbacks...)` retrieves it anywhere the request
  context flows; it returns `nil` when nothing is stored and no non-empty
  fallback is given.
- An API handler parameter of type `logx.Logger` is resolved automatically
  from the request context.
- When a handler parameter is served from a resource struct field whose type
  has a `WithLogger(logx.Logger)` method (the `LoggerConfigurable` shape),
  the injected value is a per-request copy configured with the request
  logger.

The request ID is carried as a logger *namespace* — log lines show
`[request_id:<uuid>]` — not as a structured field. `logx.Logger` has no
`With(key, value)` field API at v0.37.0; attach context by deriving
`Named(...)` children or by formatting values into the message.

## Replacing Or Wrapping The Logger

The concrete implementation is **not replaceable** at v0.37.0. The root zap
logger is a package-level value in `internal/logx`, framework packages
capture their named children at package initialization, and `logx.Logger` is
not a DI-provided type — so there is no `fx.Decorate` point and no exported
replacement hook. Framework-internal output (library, format, destination)
is fixed.

What an application can do instead:

- Implement `logx.Logger` for its own code — the framework never forces
  application code through the built-in implementation.
- Override the request-scoped logger downstream of a point it controls by
  calling `contextx.SetLogger(ctx, custom)`; `contextx.Logger` lookups and
  handler-parameter injection then resolve the custom logger, while
  framework-internal component logs still use the built-in one.
- Pass a custom logger into components it constructs via
  `LoggerConfigurable[T].WithLogger`.
- Redirect or collect stdout at the deployment level for log shipping.

## Using The Framework Logger Outside DI

Application integration code can call `vef.NamedLogger(name string) logx.Logger`
when it needs the framework logger outside dependency injection. That helper
is exported from the root `vef` package; it is not an additional top-level
symbol in `logx`.

```go
import "github.com/coldsmirk/vef-framework-go"

logger := vef.NamedLogger("myjob")
logger.Infof("processed %d records", count)
```

Inside FX-constructed components, prefer injecting `logx.Logger` directly
(or `logx.LoggerConfigurable[T]` for immutable components) over reaching for
`vef.NamedLogger`.

## Practical Advice

- Set `VEF_LOG_LEVEL=error` in CI and test runs to silence framework chatter;
  the framework's own CI does exactly this.
- Inside request handling, prefer `contextx.Logger(ctx)` or a `logx.Logger`
  handler parameter so entries correlate by request ID; reserve
  `vef.NamedLogger` for background and integration code.
- Guard voluminous debug logging with `logger.Enabled(logx.LevelDebug)`.
- Expect ANSI escape codes in the level column when parsing captured output;
  the console encoding is designed for humans, not machines.
- `Panic`/`Panicf` really panic — reserve them for unrecoverable
  initialization failures.

## See Also

- [Extension Points — Logging](../reference/extension-points#logging) covers
  the same `vef.NamedLogger` helper alongside the framework's other
  extension points.
- [Lifecycle — App-level middleware order](../core-concepts/lifecycle#app-level-middleware-order)
  shows where the request-ID and logger middlewares sit in the chain.
- [Extending Handler Parameters](../advanced/extending-parameters) documents
  the `contextx` helpers and the built-in handler-parameter resolvers.
- [Configuration — Environment overrides](../getting-started/configuration#environment-overrides)
  lists `VEF_LOG_LEVEL` among the framework's environment variables.
