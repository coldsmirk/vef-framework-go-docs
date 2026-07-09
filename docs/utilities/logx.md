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

## See Also

- [Extension Points — Logging](../reference/extension-points#logging) covers
  the same `vef.NamedLogger` helper alongside the framework's other
  extension points.
