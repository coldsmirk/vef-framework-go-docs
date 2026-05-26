---
sidebar_position: 2
---

# 事件总线

VEF 把事件系统做成了一个可插拔的多 transport 平台，对外通过单一的 `event.Bus` 暴露。FX 会负责启动总线，应用只管定义事件、发布、订阅，剩下的由路由配置决定。

> 事件系统在 v0.21 被重写为可插拔多 transport 平台（`feat(event)!: rewrite event system as pluggable multi-transport platform`），之后不断加固：类型化订阅（`SubscribeTyped`）、能力感知路由、事务性 outbox、Redis Streams transport、Inbox 去重、路由检查器、明确的错误 sentinel。本页对应 v0.26 状态。

## FX 自动注入的内容

框架启动后，可通过依赖注入获得：

| 接口 | 用途 |
| --- | --- |
| `event.Bus` | 发布 + 订阅入口 |
| `event.RouteInspector` | 路由只读查询 —— 用于在 OnStart 做快速失败检查 |
| `transport.Transport`（FX group `vef:event:transports`） | 所有已注册 transport（memory、outbox、redis-stream、自定义） |
| `event.ErrorSink` | 处理非同步发布失败（默认按 error 级别打日志） |

`Bus.Start` / `Bus.Shutdown` 都由 FX 驱动。`fx.Provide` 阶段调用 Publish 会返回 `event.ErrBusNotStarted`，应该改成 `fx.Invoke` 或 OnStart 钩子。

## Bus 接口

```go
type Bus interface {
    Publish(ctx context.Context, evt Event, opts ...PublishOption) error
    PublishBatch(ctx context.Context, evts []Event, opts ...PublishOption) error
    Subscribe(eventType string, h Handler, opts ...SubscribeOption) (Unsubscribe, error)
}
```

要点：

- 返回值反映 transport 是否接受了 frame，不代表订阅端 handler 是否成功。
- 在 `Start` 之前调用 Subscribe 会被缓冲，启动时统一刷出，所以 FX 装配顺序无所谓。

## 定义事件

任何实现了 `EventType()` 的值都可以发布。该方法必须在 `T` 的零值上安全可调用，因为 `SubscribeTyped[T]` 是用 `var zero T` 推断主题的。

```go
type UserCreatedEvent struct {
    UserID string `json:"userId"`
    Email  string `json:"email"`
}

func (*UserCreatedEvent) EventType() string { return "user.created" }
```

事件类型字符必须落在 `^[a-zA-Z0-9._-]+$` 范围（`transport.EventTypePattern`）。Bus 在 Publish 和 Subscribe 入口都强制校验，超出范围会返回 `event.ErrInvalidEventType`。

## 发布

```go
err := bus.Publish(ctx, &UserCreatedEvent{UserID: "u-1"})
```

PublishOption 从左到右合成：

| 选项 | 作用 |
| --- | --- |
| `event.WithTx(tx orm.DB)` | 走 `TxTransport`，事务提交后事件才可见。事务性 outbox 必备；若路由里没有 `TxTransport` 则返回 `event.ErrTxRequired`。 |
| `event.WithAsync()` | 把发布丢到 bus 的异步队列并立即返回；错误走 `ErrorSink`，不会返回给调用者。 |
| `event.WithSource(name)` | 覆盖 `Envelope.Source`（默认是 `vef.app.name`）。 |
| `event.WithOccurredAt(t)` | 覆盖 `Envelope.OccurredAt`（默认 `time.Now`）。 |
| `event.WithCorrelationID(id)` | 自定义关联 ID。 |
| `event.WithHeaders(map)` | 合并任意 header 进 envelope。 |

`WithTx` 和 `WithAsync` 互斥 —— 同时使用返回 `event.ErrTxAsyncMutex`。

## 订阅

```go
unsub, err := bus.Subscribe(
    "user.created",
    func(ctx context.Context, env event.Envelope) error {
        return nil
    },
    event.WithGroup("user-projection"),
)
```

SubscribeOption：

| 选项 | 作用 |
| --- | --- |
| `event.WithGroup(name)` | 消费者组。**当路由命中任何 at-least-once transport（outbox sink、Redis Streams）时必须显式提供**，否则返回 `event.ErrGroupRequired`。该 group 同时是 Inbox 去重作用域和 Redis Streams XGROUP，重启期间必须保持稳定。 |
| `event.WithConcurrency(n)` | 单个订阅的 worker 数量，默认 1。 |

### 类型化订阅

`SubscribeTyped[T]` 会自动把消息解码成你期望的具体类型：

```go
unsub, err := event.SubscribeTyped(bus,
    func(ctx context.Context, evt *UserCreatedEvent, env event.Envelope) error {
        return projection.Apply(ctx, evt)
    },
    event.WithGroup("user-projection"),
)
```

T 可以是指针类型（推荐）也可以是其指针实现了 `Event` 的值类型。bus 同时兼容进程内投递（payload 已是 T）和跨进程投递（`RawPayload`，canonical JSON body）。

## Envelope 与 Frame

每个 `Event` 会被包进 `Envelope`，携带 transport 级元数据：`ID`、`Type`、`Source`、`OccurredAt`、`PublishedAt`、`TraceID` / `SpanID`、`CorrelationID`、`Headers`、`Payload`。跨进程 transport 会把 body 序列化进 `transport.Frame`（Body 为 canonical JSON），bus 在收到时解码回 `Envelope.Payload = RawPayload{...}`，再交给 `SubscribeTyped[T]` 反序列化。

## Transport

`transport.Transport` 是可插拔的投递后端，每个实现都通过 `Capabilities` 声明自己的能力：

| 能力 | 含义 |
| --- | --- |
| `Durable` | 消息能在进程重启后存活 |
| `Transactional` | 实现 `TxTransport`，可被 `WithTx` 选中 |
| `Ordered` | 同 partition 内保持发布顺序 |
| `AtLeastOnce` | 可能重复投递 → Inbox 中间件会自动挂上 |
| `SupportsGroups` | `WithGroup` 影响投递语义（负载均衡） |
| `PublishOnly` | 只接受发布但不会投递（典型如事务性 outbox） |

内置 transport：

| 包 | 名称 | 能力 |
| --- | --- | --- |
| `event/transport/memory` | `memory` | 进程内、ordered、at-most-once。默认回退。 |
| `event/transport/outbox` | `outbox` | 持久化、事务、durable、at-least-once、**publish-only**。一个 relay 把记录推到 sink transport。 |
| `event/transport/redisstream` | `redis_stream` | durable、at-least-once、支持 group。跨进程扇出。 |

`PublishOnly` 很重要：只把路由配到 outbox 上，事件能发出去但永远没人能订阅到。订阅者要挂在 sink transport 上（outbox 的 `sink` 配置）。Bus 在解析 Subscribe 路由时会自动过滤掉 publish-only 的 transport。

## 路由

路由是声明式的 —— TOML 里的 `[vef.event.routing]` 规则按从上到下顺序、使用 `path.Match` 语义（`*`、`?`、`[abc]`）匹配。第一条匹配命中即停止；扇出由列出多个 transport 来表达。

```toml
[vef.event]
default_transport = "memory"

[[vef.event.routing]]
pattern    = "user.*"
transports = ["memory", "outbox"]

[[vef.event.routing]]
pattern    = "approval.*"
transports = ["outbox"]
```

没有规则匹配时使用 `default_transport`。完全无路由命中时 Publish 返回 `event.ErrNoRouteMatched`。

## 路由检查（启动期快速失败）

依赖特定投递语义的模块应该在 OnStart 阶段就断言路由，而不是等到第一次 Publish/Subscribe 才失败：

```go
type RouteInspector interface {
    HasTransactionalRoute(eventType string) bool
    HasSubscribableTransport(eventType string) bool
}
```

- `HasTransactionalRoute`：使用 `WithTx`（事务性 outbox 模式）的模块必须确认路由里有 `Transactional` transport，否则第一次 `WithTx` 发布就会 `ErrTxRequired`。
- `HasSubscribableTransport`：框架侧自行订阅事件的模块（binding listener、projection、集成 handler）必须确认路由里有可订阅 transport，否则路由若只解析到 publish-only transport，应用能启动，但每次 Subscribe 都会 `ErrNoRouteMatched`。v0.25.0 新增。

Approval 模块的 binding listener 和 outbox 发布都依赖这两个检查 —— 参考 [Approval 模块](../modules/approval)。

## 中间件

Bus 会运行 FX group `vef:event:publish-middlewares` 和 `vef:event:consume-middlewares` 上的中间件。内置中间件（通过 `[vef.event.middleware]` 切换）：

| 中间件 | 作用 | 激活方式 |
| --- | --- | --- |
| logging | 包住 publish/consume 的结构化日志 | `logging` 开关 |
| tracing | W3C trace/span 跨 transport 传播 | `tracing` 开关；在信任边界处把 `tracing_strict = true` 打开，把跨进程进来的 TraceID 当作不可信 |
| metrics | 计数和延迟直方图 | 可自定义 `MetricsRecorder` |
| recover | 把 panic 包装为 `event.ErrHandlerPanic` | `recover` 开关 |
| inbox | 消费侧幂等去重 | `inbox` 开关；**只**在 transport 的 `Capabilities.AtLeastOnce = true` 时挂上 |

## Inbox 幂等

对 at-least-once transport，Inbox 中间件按 `(envelope_id, consumer_group)` 持久化记录，已完成的消息再来时直接 ack 而不会重跑业务。

```toml
[vef.event.inbox]
retention        = "168h"   # 默认 7 天
processing_lease = "10m"    # 一个 worker 持有未完成 claim 的最长时间
cleanup_interval = "1h"
```

Bus 在启动时会校验 `inbox.retention` 是否大于 outbox 最坏指数退避时长 —— 见 `config.ErrInboxRetentionTooShort`。如果 retention 小于退避总长，重试到达时去重记录可能已经被清掉，会发生重复执行。

## 事务性 Outbox

启用 `[vef.event.transports.outbox]` 后，outbox transport 会在业务写入的同一事务里持久化记录，relay goroutine 再把它们转发给 `sink` transport（通常是 `memory` 或 `redis_stream`）。

```toml
[vef.event.transports.outbox]
enabled          = true
sink             = "redis_stream"
relay_interval   = "5s"
max_retries      = 10
batch_size       = 100
lease_multiplier = 4
min_lease        = "1s"
cleanup_interval = "1h"
completed_ttl    = "168h"
```

生产者：

```go
err := bus.Publish(ctx, evt, event.WithTx(tx))
```

事务提交后 relay 最终会转发；事务回滚则记录消失。

## Redis Streams Transport

```toml
[vef.event.transports.redis_stream]
enabled          = true
stream_prefix    = "vef:events:"
max_len_approx   = 100000
block_timeout    = "5s"
claim_idle       = "30s"
claim_interval   = "10s"
claim_batch_size = 64
start_id         = "0"
```

订阅者必须提供稳定的 `WithGroup` —— 它会成为 Redis XGROUP，并需要在重启间稳定。

## 错误 Sentinel

| 错误 | 含义 |
| --- | --- |
| `event.ErrBusNotStarted` | 在 `Start` 之前 publish |
| `event.ErrBusAlreadyStarted` | 重复 `Start` |
| `event.ErrTxRequired` | `WithTx` 但路由里没有 `TxTransport` |
| `event.ErrTransportNotFound` | 路由引用了未注册的 transport |
| `event.ErrAsyncQueueFull` | 异步队列满（走 `ErrorSink`） |
| `event.ErrQueueFull` | transport 在非阻塞策略下拒绝 publish |
| `event.ErrHandlerPanic` | 订阅 handler 内 panic |
| `event.ErrShutdownTimeout` | 优雅停机超时 |
| `event.ErrNoRouteMatched` | 没有任何路由匹配 |
| `event.ErrUnknownPayload` | `SubscribeTyped` 无法解码到 T |
| `event.ErrPayloadTooLarge` | payload/header 超过框架限制 |
| `event.ErrInvalidEventType` | 事件类型含非法字符 |
| `event.ErrNilTypeParameter` | `SubscribeTyped` 使用了 nil 接口类型参数 |
| `event.ErrGroupRequired` | at-least-once 订阅缺 `WithGroup` |
| `event.ErrTxAsyncMutex` | `WithTx` 和 `WithAsync` 同时使用 |

## 内置框架事件类型

| 事件类型 | 来源 |
| --- | --- |
| `vef.api.request.audit` | API 审计 |
| `vef.security.login` | 登录流 |
| `vef.security.role_permissions.changed` | 角色权限缓存失效 |
| `vef.storage.file.claimed` | 业务事务采纳了一个 pending 的上传 claim |
| `vef.storage.file.deleted` | storage 删除 worker 把文件从后端删除完成 |
| `vef.storage.delete.dead_letter` | 删除 worker 重试用尽，该行被 park 起来供人工排查 |
| `vef.approval.task.created` | 审批任务创建（v0.25，所有任务创建路径统一发出） |

启用 approval 模块后，还会发布更多审批域事件 —— 参考 [Approval 模块](../modules/approval)。

## 典型装配模式

订阅器通常通过集成模块挂上：

```go
var Module = vef.Module(
    "app:event",
    vef.Invoke(registerUserProjections),
)

func registerUserProjections(lc fx.Lifecycle, bus event.Bus, inspector event.RouteInspector) error {
    if !inspector.HasSubscribableTransport("user.created") {
        return fmt.Errorf("user.created has no subscribable transport in routing")
    }

    unsub, err := event.SubscribeTyped(bus,
        func(ctx context.Context, evt *UserCreatedEvent, env event.Envelope) error {
            return nil
        },
        event.WithGroup("user-projection"),
    )
    if err != nil {
        return err
    }

    lc.Append(fx.Hook{OnStop: func(context.Context) error { unsub(); return nil }})
    return nil
}
```

## Transport 选择建议

| 需求 | Transport |
| --- | --- |
| 进程内扇出、低延迟、fire-and-forget | `memory` |
| 发布要和业务写入一起提交 | `outbox`（→ sink） |
| 跨进程投递，at-least-once | `redis_stream` |
| 可靠审批 / saga 事件 | `outbox` + `redis_stream` sink |

## 下一步

继续阅读 [缓存](./cache) 看如何把事件发布连到缓存失效流程，或读 [Approval 模块](../modules/approval) 来看事务性 outbox 的典型用法。
