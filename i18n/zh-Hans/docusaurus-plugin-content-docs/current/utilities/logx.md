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

## 默认实现

开箱即用的 `logx.Logger` 由 `zap.SugaredLogger`（`go.uber.org/zap`）支撑，
在 `internal/logx` 内部构造：

| 方面 | 行为 |
| --- | --- |
| Encoding | console——单行、面向人的文本；没有 JSON 选项 |
| 输出目标 | 日志条目写到 stdout；zap 自身的内部错误写到 stderr |
| 行结构 | 弱化显示的 `2006-01-02T15:04:05.000` 时间戳、大写 level、`[name]` namespace、截短的 caller 路径、message |
| 颜色 | level 列始终带 ANSI 颜色（`zapcore.CapitalColorLevelEncoder`）；时间戳、caller 和 name 的样式只在 stdout 支持时应用（termenv 检测：非 TTY 输出和 `NO_COLOR` 下禁用） |
| Stack traces | 禁用 |

root logger 是进程启动时构造一次的 package-level 值；框架交出的每个
logger——组件 logger、`vef.NamedLogger`、request-scoped logger——都是这个
唯一 root 的 `Named` child。FX lifecycle 事件和配置加载器（viper）的消息
也通过内部的 `log/slog` bridge 汇入同一个 logger，过滤到 warn 及以上。

encoding 和输出目标没有配置入口：要把日志送到别处，请在部署层面收集
stdout。

## 日志级别

生效阈值来自 `VEF_LOG_LEVEL` 环境变量（`config.EnvLogLevel`），进程启动时
读取一次：

| `VEF_LOG_LEVEL` | 阈值 |
| --- | --- |
| `debug` | `logx.LevelDebug` |
| `warn` | `logx.LevelWarn` |
| `error` | `logx.LevelError` |
| `panic` | `logx.LevelPanic` |
| 其他任何值，包括未设置和 `info` | `logx.LevelInfo`（默认） |

匹配不区分大小写。阈值是一个共享的 `zap.AtomicLevel`，同时约束所有 named
logger——没有 per-component 级别配置，也没有公开的运行时 API 在启动后修改
级别（内部有 setter 但未导出）。用 `Logger.Enabled(level)` 来 guard 开销
较大的日志消息构造。

## Request-Scoped 日志

两个内置的 app middleware 协作，让每个 HTTP 请求拿到自己的 logger：

1. request-ID middleware（order `-650`）复用合法的传入 `X-Request-ID`
   header，否则通过 `id.GenerateUUID` 生成 UUID v7，并把该值回写到响应
   header。
2. logger middleware（order `-600`）从 root logger 派生名为
   `request_id:<id>` 的 child logger，用 `contextx.SetLogger` 同时存入
   fiber locals 和内嵌的 `context.Context`，因此 fiber handler 和普通
   `context.Context` 消费方解析到同一个 logger。

在下游，request logger 有三条传递路径：

- `contextx.Logger(ctx, fallbacks...)` 在请求 context 流经的任何位置取回
  它；没有存入 logger 且没有非空 fallback 时返回 `nil`。
- 类型为 `logx.Logger` 的 API handler 参数会自动从请求 context 解析。
- 当 handler 参数由 resource struct field 提供、且该 field 类型带有
  `WithLogger(logx.Logger)` 方法（即 `LoggerConfigurable` 形态）时，注入的
  是配置了 request logger 的 per-request 副本。

request ID 作为 logger 的 *namespace* 携带——日志行显示
`[request_id:<uuid>]`——而不是 structured field。`logx.Logger`
没有 `With(key, value)` 字段 API；要附加上下文，请派生 `Named(...)` child
或把值格式化进 message。

## 替换或包装 Logger

具体实现**不可替换**。root zap logger 是 `internal/logx` 里的
package-level 值，框架各 package 在包初始化时就捕获了各自的 named child，
而且 `logx.Logger` 不是 DI 提供的类型——因此没有 `fx.Decorate` 切入点，也
没有导出的替换钩子。框架内部的输出（库、格式、目标）是固定的。

应用可以做的替代方案：

- 为自己的代码实现 `logx.Logger`——框架从不强迫应用代码走内置实现。
- 在自己控制的位置调用 `contextx.SetLogger(ctx, custom)`，覆盖下游的
  request-scoped logger；之后 `contextx.Logger` 查找和 handler 参数注入都
  会解析到自定义 logger，但框架内部组件的日志仍走内置 logger。
- 通过 `LoggerConfigurable[T].WithLogger` 把自定义 logger 传给自己构造的
  组件。
- 在部署层面重定向或收集 stdout 来做日志采集。

## 在 DI 之外使用框架 Logger

当集成代码需要在 dependency injection 之外拿框架 logger 时，可以调用
`vef.NamedLogger(name string) logx.Logger`。这个 helper 从 root `vef`
package 导出，不是 `logx` 里的额外 top-level symbol。

```go
import "github.com/coldsmirk/vef-framework-go"

logger := vef.NamedLogger("myjob")
logger.Infof("processed %d records", count)
```

注意 `logx.Logger` **不在** DI 图中——把它声明为 FX 构造器参数会导致启动
失败。在由 FX 构造的组件内部，应在构造器体内调用 `vef.NamedLogger(...)`
（或当框架配置你的组件时经 `logx.LoggerConfigurable[T]` 接收 logger）；
请求 handler 直接声明 `logx.Logger` 参数即可，由 API 引擎按请求注入。

## 实用建议

- 在 CI 和测试运行中设置 `VEF_LOG_LEVEL=error`，压掉框架自身的日志噪声；
  框架自己的 CI 正是这么做的。
- 在请求处理内部，优先用 `contextx.Logger(ctx)` 或 `logx.Logger` handler
  参数，让日志条目按 request ID 关联；`vef.NamedLogger` 留给后台任务和集成
  代码。
- 用 `logger.Enabled(logx.LevelDebug)` guard 大量的 debug 日志。
- 解析采集到的输出时，注意 level 列会带 ANSI 转义码；console encoding 是
  给人看的，不是给机器解析的。
- `Panic`/`Panicf` 真的会 panic——只用于不可恢复的初始化失败。

## 参见

- [Extension Points — Logging](../reference/extension-points#日志) 涵盖了同样的
  `vef.NamedLogger` helper，并放在框架其他 extension points 旁边一起讲。
- [Lifecycle — App 级中间件顺序](../core-concepts/lifecycle#app-级中间件顺序)
  展示了 request-ID 和 logger middleware 在链路中的位置。
- [扩展 Handler 参数](../advanced/extending-parameters) 记录了 `contextx`
  辅助函数和内置的 handler 参数解析器。
- [配置 — 环境变量](../getting-started/configuration#环境变量)
  把 `VEF_LOG_LEVEL` 列在框架环境变量之中。
