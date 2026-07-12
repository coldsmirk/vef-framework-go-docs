---
sidebar_position: 4
---

# 事件总线

VEF 把事件系统做成了一个可插拔的多 transport 平台，对外通过单一的 `event.Bus` 暴露。FX 会负责启动总线，应用只管定义事件、发布、订阅，剩下的由路由配置决定。

> 事件系统在 v0.21 被重写为可插拔多 transport 平台（`feat(event)!: rewrite event system as pluggable multi-transport platform`），之后不断加固：类型化订阅（`SubscribeTyped`）、能力感知路由、事务性 outbox、Redis Streams transport、Inbox 去重、路由检查器、明确的错误 sentinel。本页描述当前公开 API surface；老版本快照与现在已经不再兼容。

## FX 自动注入的内容

框架启动后，可通过依赖注入获得：

| 接口 | 用途 |
| --- | --- |
| `event.Bus` | 发布 + 订阅入口 |
| `event.RouteInspector` | 路由只读查询 —— 用于在 OnStart 做快速失败检查 |
| `transport.Transport`（FX group `vef:event:transports`） | 所有已注册 transport（memory、outbox、redis-stream、自定义） |
| `event.ErrorSink` | 处理非同步发布失败（默认按 error 级别打日志） |

Bus 生命周期由 FX 驱动。`fx.Provide` 阶段调用 Publish 会返回 `event.ErrBusNotStarted`，应该改成 `fx.Invoke` 或 OnStart 钩子。

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
- 在 `Start` 之前调用 `Publish` 或 `PublishBatch` 会返回 `event.ErrBusNotStarted`；在 `Start` 之前调用 `Subscribe` 会被缓冲，启动时统一刷出，所以 FX 装配顺序无所谓。
- `PublishBatch` 不是跨 transport 的原子性保证。非事务 transport 可能在返回错误前已经接受了前面的 frame；使用 `WithTx` 时，事务性 transport 会加入调用方事务，随外层 commit / rollback 生效。

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

PublishOption 从左到右合成；多个选项设置同一字段时，后面的值覆盖前面的值：

| 选项 | 作用 |
| --- | --- |
| `event.WithTx(tx orm.DB)` | 只走 capability 为 transactional 的 transport。事务性 outbox 必备；若路由里没有 `TxTransport` 则返回 `event.ErrTxRequired`。 |
| `event.WithAsync()` | 把发布丢到 bus 的异步队列并立即返回；异步 publish 错误通过 `ErrorSink` 报告，不返回给调用方。如果队列已满，bus 会去掉 `WithAsync` 后回退为同步 publish；只有回退 publish 也失败时才报告 `event.ErrAsyncQueueFull`。已入队任务使用脱离请求取消的 context，所以原请求结束不会中断已经接受的异步事件。 |
| `event.WithSource(name)` | 覆盖 `Envelope.Source`（默认是 `vef.app.name`）。 |
| `event.WithOccurredAt(t)` | 覆盖 `Envelope.OccurredAt`（默认 `time.Now`）。 |
| `event.WithCorrelationID(id)` | 自定义关联 ID。省略时 bus 会继承 `contextx.RequestID(ctx)`。 |
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
| `event.WithGroup(name)` | 消费者组。**当可订阅路由命中任何 at-least-once transport（例如 Redis Streams 或其他 durable sink）时必须显式提供**，否则返回 `event.ErrGroupRequired`。该 group 同时是 Inbox 去重作用域和 Redis Streams XGROUP，重启期间必须保持稳定。 |
| `event.WithConcurrency(n)` | 单个订阅的 worker 数量，默认 1；非正数会被忽略。大于 1 时用顺序换吞吐——多个 worker 争抢同一条订阅 feed，即使在 `Ordered` transport 上 handler 的执行也会交错；顺序敏感的订阅者保持 1。 |

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

T 可以是指针类型（推荐）也可以是其指针实现了 `Event` 的值类型。bus 同时兼容进程内投递（payload 已是 T）和跨进程投递（`RawPayload`，canonical JSON body）。`SubscribeTyped[event.Event]` 会返回 `event.ErrNilTypeParameter`。

## Envelope 与 Frame

每个 `Event` 会被包进 `Envelope`，携带 transport 级元数据：

| 字段 | 契约 |
| --- | --- |
| `ID` | 框架生成的消息 ID；重试期间保持稳定，也是 Inbox 去重键 |
| `Type` | `Event.EventType()`，用于路由和分发 |
| `Source` | 来自 `WithSource(...)`；省略时使用 `vef.app.name` |
| `OccurredAt` | 来自 `WithOccurredAt(...)`；省略时使用 publish 时间 |
| `PublishedAt` | 框架第一次接受 publish 的时间 |
| `TraceID` / `SpanID` | tracing 中间件启用时写入 |
| `CorrelationID` | 来自 `WithCorrelationID(...)`；省略时使用 `contextx.RequestID(ctx)` |
| `Headers` | `WithHeaders(...)` 和中间件合并出的调用方元数据 |
| `Payload` | 进程内投递时是原始 `Event`；跨进程后是 `RawPayload` |

跨进程 transport 会把 body 序列化进 `transport.Frame`（Body 为 canonical JSON），bus 在收到时解码回 `Envelope.Payload = RawPayload{...}`，再交给 `SubscribeTyped[T]` 反序列化。

`CorrelationID` 会跨越所有 transport 边界，包括 outbox 和 Redis Streams。如果你的部署里 request ID 属于敏感数据，应注册 publish middleware，在持久化 transport 看到 frame 之前清空或替换 `Envelope.CorrelationID`。

Envelope 大小限制也是公开 publish 契约：

| 限制 | 值 | 错误 |
| --- | --- | --- |
| JSON frame body | 1 MiB | `event.ErrPayloadTooLarge` |
| header 条目数 | 32 | `event.ErrPayloadTooLarge` |
| header key | 128 bytes | `event.ErrPayloadTooLarge` |
| header value | 1024 bytes | `event.ErrPayloadTooLarge` |

## Transport

`transport.Transport` 是可插拔的投递后端，每个实现都通过 `Capabilities` 声明自己的能力：

| 能力 | 含义 |
| --- | --- |
| `Durable` | 消息能在进程重启后存活 |
| `Transactional` | 实现 `TxTransport`，可被 `WithTx` 选中 |
| `Ordered` | 消息按发布顺序*交付*给订阅。该保证只对串行消费成立：`WithConcurrency(n > 1)` 时多个 worker 从同一条有序 feed 拉取，handler 执行交错，可观察的处理顺序即告丢失 |
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

memory transport 的队列策略值是 `error`、`block` 和 `drop_oldest`。`publish_timeout` 只在 memory transport 使用 `block` 策略时限制 `Publish` 的等待时间；默认 `error` 策略会直接返回 `event.ErrQueueFull`，不会等待队列腾出空间。

公开 Go 包 `event/transport/memory.Config` 暴露 `QueueSize`、`FullPolicy`
和 `PublishTimeout`。`QueueSize` 默认是 `1024`；`FullPolicy` 默认是
`FullPolicyError`。

## 路由

路由是声明式的 —— TOML 里的 `[vef.event.routing]` 规则按从上到下顺序、使用 `path.Match` 语义（`*`、`?`、`[abc]`）匹配。第一条匹配命中即停止；扇出由列出多个 transport 来表达。

```toml
[vef.event]
default_transport = "memory"

[[vef.event.routing]]
pattern    = "user.*"
transports = ["memory"]

[[vef.event.routing]]
pattern    = "approval.*"
transports = ["outbox", "redis_stream"]
```

没有规则匹配时会使用 `default_transport`。缺少专门的 routing rule
本身不会让 publish 失败。`event.ErrNoRouteMatched` 只用于 route 解析不到任何
transport，或者订阅时解析结果只包含 publish-only transport 的情况。

如果某条 route 列出了 `outbox` 且还需要订阅者收到消息，就必须把配置的
outbox `sink` 也放进同一个 transports 列表。框架会在启动时校验这一点，
避免事件发布到一个 sink、订阅者却挂到另一个 transport 上。

只包含 `["outbox"]` 的 route 允许用于没有订阅者的纯发布流程。route 引用了未知 transport 时会在框架启动阶段失败；不要把 `ErrTransportNotFound` 当成唯一可能的启动错误形状。

## 路由检查（启动期快速失败）

依赖特定投递语义的模块应该在 OnStart 阶段就断言路由，而不是等到第一次 Publish/Subscribe 才失败：

```go
type RouteInspector interface {
    HasTransactionalRoute(eventType string) bool
    HasSubscribableTransport(eventType string) bool
}
```

- `HasTransactionalRoute`：使用 `WithTx`（事务性 outbox 模式）的模块必须确认路由里有 `Transactional` transport，否则第一次 `WithTx` 发布就会 `ErrTxRequired`。
- `HasSubscribableTransport`：框架侧自行订阅事件的模块（集成 handler、事件驱动的 projection）必须确认路由里有可订阅 transport，否则路由若只解析到 publish-only transport，应用能启动，但每次 Subscribe 都会 `ErrNoRouteMatched`。v0.25.0 新增。

Approval 模块在启动时对它的 `approval.*` 事件断言 `HasTransactionalRoute`（其业务投影从持久化状态收敛，不再订阅事件）—— 参考 [Approval 模块](../approval)。

## 中间件

Bus 会运行 FX group `vef:event:publish-middlewares` 和 `vef:event:consume-middlewares` 上的中间件。内置中间件（通过 `[vef.event.middleware]` 切换）：

| 中间件 | 作用 | 激活方式 |
| --- | --- | --- |
| logging | 包住 publish/consume 的结构化日志 | `logging` 开关 |
| tracing | W3C trace/span 跨 transport 传播 | `tracing` 开关；在信任边界处把 `tracing_strict = true` 打开，把跨进程进来的 TraceID 当作不可信 |
| metrics | 通过 `PublishObserved` 和 `ConsumeObserved` 记录计数和延迟直方图 | 可自定义 `MetricsRecorder` |
| recover | 把 panic 包装为 `event.ErrHandlerPanic` | `recover` 开关 |
| inbox | 消费侧幂等去重 | `inbox` 开关；**只**在 transport 的 `Capabilities.AtLeastOnce = true` 时挂上 |

公开 supporting API：

| 包 | 公开 surface |
| --- | --- |
| `event` | `AsEvents`, `ApplyPublishOptions`, `ApplySubscribeOptions`, `PublishConfig`, `SubscribeConfig`, `RawPayload`, `MetricsRecorder`, `ErrorSink`, `Unsubscribe`, `TypedHandler` |
| `event/middleware` | `PublishHandler`, `ConsumeHandler`, `PublishMiddleware`, `ConsumeMiddleware`, `ChainPublish`, `ChainConsume`, `TraceIDFromContext`, `IncomingTraceIDFromContext`, `WithTraceID`, `WithIncomingTraceID`，以及顺序常量（`OrderLogging`, `OrderTracing`, `OrderMetrics`, `OrderRecover`, `OrderInbox`） |
| `event/transport` | `Frame`, `Delivery`, `Capabilities`, `SubscribeConfig`, `Transport`, `TxTransport`, `ConsumeFunc`, `Unsubscribe`, `ErrSubscribeUnsupported`, `EventTypePattern` |
| `event/inbox` | `Status`, `StatusProcessing`, `StatusCompleted`, `AcquireResult`, `AcquireResultAcquired`, `AcquireResultCompleted`, `AcquireResultInProgress`, `Record`, `Repository`, `ErrInProgress`, `ErrLockLost`, `ErrMissingLockID`, `ErrUnknownAcquireResult` |
| `event/transport/memory` | `Name`, `Config`, `FullPolicy`, `FullPolicyError`, `FullPolicyBlock`, `FullPolicyDropOldest` |
| `event/transport/outbox` | `Name`, `Config`, `Status`, `StatusPending`, `StatusProcessing`, `StatusCompleted`, `StatusFailed`, `StatusDead`, `Record`, `Repository` |
| `event/transport/redisstream` | `Name` 和 `Config` |

`ChainPublish` 和 `ChainConsume` 会按升序 `Order` 排列 middleware；相同
order 保留注册顺序。内置 cron job 名称是 `vef:event:outbox:relay`、
`vef:event:outbox:cleanup` 和 `vef:event:inbox:cleanup`。

### Wire 值与记录字段

Inbox 的 status 和 acquire-result 都按字符串持久化：

| API | Wire 值 | 含义 |
| --- | --- | --- |
| `inbox.StatusProcessing` | `processing` | delivery 当前被某个 consumer 租约持有 |
| `inbox.StatusCompleted` | `completed` | handler 已完成；重复 delivery 会直接 ack，不会重跑业务 |
| `inbox.AcquireResultAcquired` | `acquired` | 调用方拥有 delivery，应执行 handler |
| `inbox.AcquireResultCompleted` | `completed` | delivery 此前已完成 |
| `inbox.AcquireResultInProgress` | `in_progress` | 另一个 consumer 仍持有未过期租约 |

`event/inbox.Record` 除了嵌入的 ORM 模型字段外，还暴露 JSON 字段 `eventId`、`consumerGroup`、`status`、`lockId`、`lockedUntil`、`completedAt`。

Outbox status 也按字符串持久化：

| API | Wire 值 | 含义 |
| --- | --- | --- |
| `outbox.StatusPending` | `pending` | 等待第一次 dispatch |
| `outbox.StatusProcessing` | `processing` | 当前被 relay worker 租约持有 |
| `outbox.StatusCompleted` | `completed` | 下游 sink 已接受 frame |
| `outbox.StatusFailed` | `failed` | 最近一次 dispatch 失败，记录已排入重试 |
| `outbox.StatusDead` | `dead` | 重试预算用尽，且 DLQ 转发已成功 |

`event/transport/outbox.Record` 除了嵌入的 ORM 模型字段外，还暴露 JSON 字段 `eventId`、`eventType`、`source`、`traceId`、`spanId`、`correlationId`、`headers`、`payload`、`status`、`retryCount`、`lastError`、`processedAt`、`retryAfter`、`occurredAt`。

## Inbox 幂等

对 at-least-once transport，Inbox 中间件按 `(envelope_id, consumer_group)` 持久化记录，已完成的消息再来时直接 ack 而不会重跑业务。

```toml
[vef.event.inbox]
retention        = "168h"   # 默认 7 天
processing_lease = "10m"    # 一个 worker 持有未完成 claim 的最长时间
cleanup_interval = "1h"
```

Bus 在启动时会校验 `inbox.retention` 是否大于 outbox 最坏指数退避时长 —— 见 `config.ErrInboxRetentionTooShort`。如果 retention 小于退避总长，重试到达时去重记录可能已经被清掉，会发生重复执行。

当重复 delivery 到达时，如果另一个 consumer 仍持有 processing lease，Inbox 中间件会返回 `inbox.ErrInProgress`，让 at-least-once transport 稍后重试。已完成的重复 delivery 会被 ack，不会再次执行 handler。

## 事务性 Outbox

启用 `[vef.event.transports.outbox]` 后，outbox transport 会在业务写入的同一事务里持久化记录，relay goroutine 再把它们转发给 `sink` transport（通常是 `memory` 或 `redis_stream`）。

```toml
[vef.event.transports.outbox]
enabled          = true
sink             = "redis_stream"
relay_interval   = "10s"
max_retries      = 10
batch_size       = 100
lease_multiplier = 4
min_lease        = "15s"
cleanup_interval = "1h"
completed_ttl    = "168h"
```

省略 `sink` 时默认使用 `memory`。如果事件需要跨进程投递，就把 sink
配置为 `redis_stream`。

公开 Go 包 `event/transport/outbox.Config` 只包含 `RelayInterval`、
`MaxRetries`、`BatchSize`、`LeaseMultiplier`、`MinLease` 和 `SinkName`。
框架 TOML 块还包含 `cleanup_interval` 和 `completed_ttl`；它们驱动框架的
cleanup cron job，不是 transport 包 `Config` 的字段。

生产者：

```go
err := bus.Publish(ctx, evt, event.WithTx(tx))
```

事务提交后 relay 最终会转发；事务回滚则记录消失。

Relay 失败后使用指数退避重试（`2^retryCount` 秒，最高 1h）。重试预算用尽后，relay 会把原 frame 转发一次到 DLQ topic `vef-dlq.<eventType>`，并加上 header `vef.dlq=1`。

- 如果 DLQ 转发失败，记录保持 `failed` 且仍可被 claim，后续继续重试 DLQ hand-off。
- 如果 DLQ 转发成功，记录变为 `dead`，保留给诊断使用。
- cleanup job 只删除早于 `completed_ttl` 的 completed 行；dead 行会保留。
- 持久化的 `lastError` 会清理常见凭据片段，并截断到 256 bytes。

Bus 生成的是 canonical JSON frame，outbox 可以直接持久化。自定义代码如果直接 publish 到 outbox transport，必须提供 JSON 形状的 frame body。已经携带 `vef.dlq` header 的 frame 会被 outbox loop guard 拒绝，避免 sink 配错时无限把自己的 DLQ 流量重新持久化。

## Redis Streams Transport

```toml
[vef.event.transports.redis_stream]
enabled          = true
stream_prefix    = "vef:events:"
max_len_approx   = 100000
block_timeout    = "5s"
claim_idle       = "60s"
claim_interval   = "30s"
claim_batch_size = 64
reaper_concurrency = 4
handler_timeout = "0s"
setup_timeout   = "5s"
consumer_id      = "vef"
start_id         = "0"
```

订阅者必须提供稳定的 `WithGroup` —— 它会成为 Redis XGROUP，并需要在重启间稳定。

Redis Streams 契约细节：

- 只有 `enabled = true` 且 Redis client 可用时，transport 才会被注册；`Start` 会用 `PING` 校验连接。
- Stream key 是 `stream_prefix + eventType`；`max_len_approx` 使用 Redis `XADD MAXLEN ~`。
- `start_id = "0"` 表示新建 group 会消费已有 backlog。对 fire-and-forget topic，可设为 `"$"`，让新 group 跳过历史消息。
- `consumer_id` 只是可读前缀；运行时会追加 UUID 后缀，避免多个副本在同一 group 内重名。
- 缺失、非字符串、超大或 JSON 无效的 frame 会被视为 poison message：transport 记录日志、`XACK` 并丢弃。at-least-once 只保证格式正确的 frame。
- Handler 失败会让消息留在 pending；reaper 会根据 `claim_idle`、`claim_interval` 和 `claim_batch_size` 定期 `XCLAIM` 空闲 pending 条目。
- `reaper_concurrency` 限制每轮可并行 reclaim 的 subscription 数量。`handler_timeout` 限制 fresh delivery 和 reaper redelivery 的单次 handler 执行时间；`0s` 表示禁用 deadline。`setup_timeout` 限制 `Subscribe` 期间创建 consumer group 的等待时间。

公开 Go 包 `event/transport/redisstream.Config` 包含 `StreamPrefix`、
`MaxLenApprox`、`BlockTimeout`、`ClaimIdle`、`ClaimInterval`、
`ClaimBatchSize`、`ReaperConcurrency`、`HandlerTimeout`、`SetupTimeout`、
`ConsumerID`、`StartID`、`IdleGroupRetention` 和 `IdleGroupSweepInterval`。
`StreamPrefix` 默认是 `vef:events:`，`BlockTimeout` 默认 `5s`，`ClaimIdle`
默认 `60s`，`ClaimInterval` 默认 `30s`，`ClaimBatchSize` 默认 `64`，
`ReaperConcurrency` 默认 `4`，`SetupTimeout` 默认 `5s`，`StartID` 默认
`0`。`HandlerTimeout` 默认 `0`，表示不启用单次 handler deadline。
`IdleGroupRetention` 默认 `0`，表示不启用孤儿 group 清理；
`IdleGroupSweepInterval`（通过 `Config.EffectiveIdleGroupSweepInterval()`
读取）默认 `10m`。

### 观测 Stream 与回收孤儿 Group

一个被移除或改名的 subscriber 会把它的 consumer group 遗留在 Redis
上：不再有 consumer 读取它，lag 会无限增长。`event.StreamInspector`
接口让这一状态可观测，`IdleGroupRetention` 让它可回收：

```go
type StreamInspector interface {
    Streams(ctx context.Context) ([]StreamInfo, error)
}

type StreamInfo struct {
    Stream string
    Length int64
    Groups []StreamGroupInfo
}

type StreamGroupInfo struct {
    Name            string
    Consumers       int64
    Pending         int64
    Lag             int64
    LastDeliveredID string
}
```

`redis_stream` transport 在启用时提供该实现；这个依赖在 DI 中是可选的
（`optional:"true"`），调用方必须容忍 nil inspector。
`sys/monitor.get_event_streams` 把 `Streams` 暴露成一个 API 端点——列出
每个 stream 及其 group（consumer 数 / pending / lag / last-delivered）；
响应字段表见 [监控](./monitor) 页。

将 `idle_group_retention` 设为非零值即可开启回收。只有当一个 group
没有 pending 条目、不是当前进程的活跃订阅、且它的每一条 consumer
记录都已空闲超过 `idle_group_retention` 时才会被销毁——有 pending
条目的 group，或完全没有 consumer 记录的 group，永远不会被处理。
清理按 `idle_group_sweep_interval` 周期运行（默认 `10m`），并且总是
跳过当前进程自己的订阅。

```toml
[vef.event.transports.redis_stream]
idle_group_retention     = "24h"
idle_group_sweep_interval = "10m"
```

## 错误 Sentinel

| 错误 | 含义 |
| --- | --- |
| `event.ErrBusNotStarted` | 在 `Start` 之前 publish |
| `event.ErrBusAlreadyStarted` | 重复 `Start` |
| `event.ErrTxRequired` | `WithTx` 但路由里没有 `TxTransport` |
| `event.ErrTransportNotFound` | 框架边界上的 transport lookup 失败 |
| `event.ErrAsyncQueueFull` | 异步队列满，且回退同步 publish 也失败；通过 `ErrorSink` 报告，不返回给调用方 |
| `event.ErrQueueFull` | transport 在非阻塞策略下拒绝 publish |
| `event.ErrHandlerPanic` | 订阅 handler 内 panic |
| `event.ErrShutdownTimeout` | 优雅停机超时 |
| `event.ErrNoRouteMatched` | route 解析不到任何 transport，或 Subscribe 只解析到 publish-only transport |
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
| `vef.storage.delete.dead_letter` | 删除 worker 重试用尽；队列行已退役，事件携带排查信息 |
| `vef.approval.task.created` | 审批任务创建（v0.25，所有任务创建路径统一发出） |

启用 approval 模块后，还会发布更多审批域事件 —— 参考 [Approval 模块](../approval)。

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

继续阅读 [缓存](./cache) 看如何把事件发布连到缓存失效流程，或读 [Approval 模块](../approval) 来看事务性 outbox 的典型用法。
