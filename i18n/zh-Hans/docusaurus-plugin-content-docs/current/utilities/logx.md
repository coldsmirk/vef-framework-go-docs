---
sidebar_position: 9
---

# logx

`logx` 包定义了框架统一使用的结构化日志契约——框架中所有会打日志的组件（API、
security、event bus、CRUD、approval 等）都依赖 `logx.Logger`，从不直接依赖某个
具体的日志库。

## API 参考

| API | Contract | 作用 |
| --- | --- | --- |
| `logx.Level` | `type Level int8` | 日志优先级；数值越高优先级越高。 |
| `logx.LevelDebug = 1` | debug level constant | 通常较多，生产环境一般关闭。 |
| `logx.LevelInfo = 2` | info level constant | 默认日志优先级。 |
| `logx.LevelWarn = 3` | warning level constant | 比 info 更重要，但不一定需要逐条人工处理。 |
| `logx.LevelError = 4` | error level constant | 表示非预期应用行为的高优先级日志。 |
| `logx.LevelPanic = 5` | panic level constant | 记录消息后触发 panic。 |
| `logx.Logger` | logging interface | 由框架提供的 logger 和自定义 logger 实现的契约。 |
| `logx.LoggerConfigurable[T]` | generic interface | 由 immutable component 实现，通过 `WithLogger` 返回配置了 logger 的副本。 |

`Level.String() string` 返回 `debug`、`info`、`warn`、`error` 或 `panic`。
未知 level 值，包括 zero value，返回 `unknown`。

## Logger 接口

| Method | Contract |
| --- | --- |
| `Logger.Named(name string) logx.Logger` | 返回带 namespace 的 child logger |
| `Logger.WithCallerSkip(skip int) logx.Logger` | 返回调整 caller stack-frame reporting 的 logger |
| `Logger.Enabled(level logx.Level) bool` | 报告指定 level 是否启用 |
| `Logger.Sync()` | flush buffered log entries；接口不返回 error |
| `Logger.Debug(message string)` | 记录 debug message |
| `Logger.Debugf(template string, args ...any)` | 记录 formatted debug message |
| `Logger.Info(message string)` | 记录 info message |
| `Logger.Infof(template string, args ...any)` | 记录 formatted info message |
| `Logger.Warn(message string)` | 记录 warning message |
| `Logger.Warnf(template string, args ...any)` | 记录 formatted warning message |
| `Logger.Error(message string)` | 记录 error message |
| `Logger.Errorf(template string, args ...any)` | 记录 formatted error message |
| `Logger.Panic(message string)` | 记录 panic message 后触发 panic |
| `Logger.Panicf(template string, args ...any)` | 记录 formatted panic message 后触发 panic |
| `LoggerConfigurable[T].WithLogger(logger logx.Logger) T` | 返回配置了指定 logger 的 component 副本 |

`LoggerConfigurable[T]` 是给 immutable、由 DI 构造的 component 用的模式——组件
构造完成后需要注入 logger 时，`WithLogger` 返回一个新的已配置副本，而不是原地
修改 receiver。

## 在 DI 之外使用框架 Logger

当集成代码需要在 dependency injection 之外拿框架 logger 时，可以调用
`vef.NamedLogger(name string) logx.Logger`。这个 helper 从 root `vef`
package 导出，不是 `logx` 里的额外 top-level symbol。

```go
import "github.com/coldsmirk/vef-framework-go"

logger := vef.NamedLogger("myjob")
logger.Infof("processed %d records", count)
```

在由 FX 构造的组件内部，优先直接注入 `logx.Logger`（对 immutable component
则注入 `logx.LoggerConfigurable[T]`），而不是使用 `vef.NamedLogger`。

## 参见

- [Extension Points — Logging](../reference/extension-points#日志) 涵盖了同样的
  `vef.NamedLogger` helper，并放在框架其他 extension points 旁边一起讲。
