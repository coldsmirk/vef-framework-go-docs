---
sidebar_position: 2
---

# Event Bus

VEF ships a pluggable multi-transport event platform behind a single `event.Bus` facade. The bus is auto-wired by FX — applications register events, publish, subscribe, and let the configured routing decide where each frame goes.

> The event system was rewritten in v0.21 as a pluggable platform (`feat(event)!: rewrite event system as pluggable multi-transport platform`) and has been hardened repeatedly since: typed delivery (`SubscribeTyped`), capability-aware routing, transactional outbox, Redis Streams transport, Inbox dedupe, route inspection, and explicit error sentinels. This page describes the current public surface; older snapshots are not API-compatible.

## What FX Wires For You

Once the framework starts, the following are available for injection:

| Interface | Purpose |
| --- | --- |
| `event.Bus` | combined publish/subscribe entry point |
| `event.RouteInspector` | read-only routing queries — used to fail fast at OnStart |
| `transport.Transport` (FX group `vef:event:transports`) | every registered transport (memory, outbox, redis-stream, custom) |
| `event.ErrorSink` | sink for out-of-band publish errors (default logs at error) |

The bus lifecycle is driven by FX. Publishing during `fx.Provide` returns `event.ErrBusNotStarted` — move the call into a `fx.Invoke` or lifecycle hook.

## The Bus Interface

```go
type Bus interface {
    Publish(ctx context.Context, evt Event, opts ...PublishOption) error
    PublishBatch(ctx context.Context, evts []Event, opts ...PublishOption) error
    Subscribe(eventType string, h Handler, opts ...SubscribeOption) (Unsubscribe, error)
}
```

Notes:

- The return value reflects whether the *transport* accepted the frame, not whether downstream handlers succeeded.
- `Publish` and `PublishBatch` before `Start` return `event.ErrBusNotStarted`; `Subscribe` registrations made before `Start` are buffered and flushed during boot, so order of FX wiring does not matter.
- `PublishBatch` is not a cross-transport atomicity guarantee. Non-transactional transports may accept earlier frames before returning an error; transactional transports participate in the caller's transaction when `WithTx` is used.

## Defining Events

Any value implementing `EventType()` is publishable. The receiver must be safe on a zero value of `T` because `SubscribeTyped[T]` derives the topic from `var zero T`.

```go
type UserCreatedEvent struct {
    UserID string `json:"userId"`
    Email  string `json:"email"`
}

func (*UserCreatedEvent) EventType() string { return "user.created" }
```

Event types are constrained to `^[a-zA-Z0-9._-]+$` (`transport.EventTypePattern`). The bus enforces this at both `Publish` and `Subscribe` entry points — any character outside the alphabet returns `event.ErrInvalidEventType`.

## Publishing

```go
err := bus.Publish(ctx, &UserCreatedEvent{UserID: "u-1"})
```

Publish-time options compose left-to-right; later options win when they set the same field:

| Option | Effect |
| --- | --- |
| `event.WithTx(tx orm.DB)` | route only through transports whose capabilities are transactional. Required for the transactional outbox pattern; returns `event.ErrTxRequired` if the resolved route has no `TxTransport`. |
| `event.WithAsync()` | hand the publish to the bus async fan-in queue and return immediately. Async publish errors are reported to `ErrorSink`, not returned to the caller. If the queue is full, the bus strips `WithAsync` and falls back to synchronous publish; only a failed fallback reports `event.ErrAsyncQueueFull`. Enqueued jobs use a cancellation-detached context so request teardown does not abort already accepted async work. |
| `event.WithSource(name)` | override `Envelope.Source` (defaults to `vef.app.name`). |
| `event.WithOccurredAt(t)` | override `Envelope.OccurredAt` (defaults to `time.Now`). |
| `event.WithCorrelationID(id)` | caller-controlled correlation key. When omitted, the bus inherits `contextx.RequestID(ctx)`. |
| `event.WithHeaders(map)` | merge arbitrary headers into the envelope. |

`WithTx` and `WithAsync` are mutually exclusive — combining them returns `event.ErrTxAsyncMutex`.

## Subscribing

```go
unsub, err := bus.Subscribe(
    "user.created",
    func(ctx context.Context, env event.Envelope) error {
        return nil
    },
    event.WithGroup("user-projection"),
)
```

Subscribe options:

| Option | Effect |
| --- | --- |
| `event.WithGroup(name)` | consumer group. **Required** when the subscribable route resolves to any at-least-once transport, such as Redis Streams or another durable sink — otherwise `event.ErrGroupRequired`. The group is the dedupe scope for the Inbox middleware and the XGROUP for Redis Streams; it must stay stable across restarts. |
| `event.WithConcurrency(n)` | worker count per subscription. Defaults to 1; non-positive values are ignored. |

### Typed subscriptions

`SubscribeTyped[T]` decodes the wire payload into your concrete type:

```go
unsub, err := event.SubscribeTyped(bus,
    func(ctx context.Context, evt *UserCreatedEvent, env event.Envelope) error {
        return projection.Apply(ctx, evt)
    },
    event.WithGroup("user-projection"),
)
```

T can be a pointer type (recommended) or a value type whose pointer also satisfies `Event`. The bus accepts both in-process delivery (payload already typed) and cross-process delivery (`RawPayload` with canonical JSON body). `SubscribeTyped[event.Event]` is rejected with `event.ErrNilTypeParameter`.

## Envelope and Frame

The framework wraps each `Event` in an `Envelope` carrying transport metadata:

| Field | Contract |
| --- | --- |
| `ID` | framework-generated message ID, stable across retries and used by Inbox dedupe |
| `Type` | `Event.EventType()`, used for routing and dispatch |
| `Source` | `WithSource(...)`, or `vef.app.name` when omitted |
| `OccurredAt` | `WithOccurredAt(...)`, or publish time when omitted |
| `PublishedAt` | time the framework first accepted the publish |
| `TraceID` / `SpanID` | set by tracing middleware when enabled |
| `CorrelationID` | `WithCorrelationID(...)`, or `contextx.RequestID(ctx)` when omitted |
| `Headers` | caller metadata merged by `WithHeaders(...)` and middleware |
| `Payload` | original `Event` for in-process delivery, `RawPayload` after crossing process boundaries |

Cross-process transports serialize the body into a `transport.Frame` whose `Body` is canonical JSON; the bus decodes it back to `Envelope.Payload = RawPayload{...}` so `SubscribeTyped[T]` can deserialize.

`CorrelationID` crosses every transport boundary, including outbox and Redis Streams. If request IDs are sensitive in your deployment, register a publish middleware that clears or replaces `Envelope.CorrelationID` before persistent transports see the frame.

Envelope size limits are part of the public publish contract:

| Limit | Value | Error |
| --- | --- | --- |
| JSON frame body | 1 MiB | `event.ErrPayloadTooLarge` |
| header entries | 32 | `event.ErrPayloadTooLarge` |
| header key | 128 bytes | `event.ErrPayloadTooLarge` |
| header value | 1024 bytes | `event.ErrPayloadTooLarge` |

## Transports

A `transport.Transport` is the pluggable backend. Each declares `Capabilities`:

| Capability | Meaning |
| --- | --- |
| `Durable` | messages survive process restart |
| `Transactional` | implements `TxTransport` — `WithTx` routes through it |
| `Ordered` | per-partition order preserved |
| `AtLeastOnce` | delivery may be duplicated → Inbox middleware activates |
| `SupportsGroups` | `WithGroup` affects load balancing |
| `PublishOnly` | accepts publishes but cannot deliver (e.g. the transactional outbox itself) |

Built-in transports:

| Package | Name | Capabilities |
| --- | --- | --- |
| `event/transport/memory` | `memory` | in-process, ordered, at-most-once. Default fallback. |
| `event/transport/outbox` | `outbox` | persistent, transactional, durable, at-least-once, **publish-only**. A relay drains records into a sink transport. |
| `event/transport/redisstream` | `redis_stream` | durable, at-least-once, supports groups. Cross-process fan-out via Redis Streams. |

`PublishOnly` is important: routing to the outbox alone makes events publishable but not deliverable. Subscribers attach to the sink transport (the outbox `sink` setting). The bus filters out publish-only transports when resolving Subscribe targets.

Memory transport queue policy values are `error`, `block`, and `drop_oldest`. `publish_timeout` only caps `Publish` when the memory transport is using the `block` policy; the default `error` policy returns `event.ErrQueueFull` instead of waiting.

The public Go package `event/transport/memory.Config` exposes `QueueSize`,
`FullPolicy`, and `PublishTimeout`. `QueueSize` defaults to `1024`;
`FullPolicy` defaults to `FullPolicyError`.

## Routing

Routing is declarative — `[vef.event.routing]` rules in TOML are matched top-to-bottom using `path.Match` semantics (`*`, `?`, `[abc]`). The first matching rule wins; fan-out is expressed by listing multiple transports.

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

When no rule matches, `default_transport` is used. A missing routing rule does
not fail publishing by itself. `event.ErrNoRouteMatched` is reserved for a route
that resolves to no transports, or for subscriptions whose resolved route only
contains publish-only transports.

If a route lists `outbox` and also needs subscribers, include the configured
outbox `sink` in the same transport list. The framework validates this at
startup so a route cannot silently publish to one sink while subscribers attach
to another.

An `["outbox"]`-only route is allowed for publisher-only flows with no subscribers. A route that references an unknown transport fails during the framework start lifecycle; do not rely on `ErrTransportNotFound` as the only possible startup error shape.

## Route Inspection (Fail-Fast Wiring)

Modules that depend on specific delivery semantics should assert routing at OnStart instead of failing at the first Publish / Subscribe:

```go
type RouteInspector interface {
    HasTransactionalRoute(eventType string) bool
    HasSubscribableTransport(eventType string) bool
}
```

- `HasTransactionalRoute`: required by modules that publish with `WithTx` (transactional outbox pattern). Without it, the first `WithTx` publish fails with `ErrTxRequired`.
- `HasSubscribableTransport`: required by modules whose framework-side code subscribes (binding listeners, projections, integration handlers). Without it, a route resolving only to publish-only transports lets the app start but every Subscribe fails with `ErrNoRouteMatched`. Added in v0.25.0.

The approval module's binding listener and outbox publisher rely on both — see [Approval module](../modules/approval).

## Middleware

The bus runs publish-side and consume-side middlewares from the FX groups `vef:event:publish-middlewares` and `vef:event:consume-middlewares`. Built-in middlewares (toggled in `[vef.event.middleware]`):

| Middleware | Purpose | Activation |
| --- | --- | --- |
| logging | structured logs around publish/consume | toggled by `logging` |
| tracing | W3C trace/span propagation across transports | toggled by `tracing`; set `tracing_strict = true` to treat incoming TraceIDs as untrusted at trust boundaries |
| metrics | counters and latency histograms via `PublishObserved` and `ConsumeObserved` | pluggable `MetricsRecorder` |
| recover | converts panics into `event.ErrHandlerPanic` | toggled by `recover` |
| inbox | consume-side idempotency dedupe | toggled by `inbox`; attaches **only** when the transport's `Capabilities.AtLeastOnce` is true |

Public supporting APIs:

| Package | Public surface |
| --- | --- |
| `event` | `AsEvents`, `ApplyPublishOptions`, `ApplySubscribeOptions`, `PublishConfig`, `SubscribeConfig`, `RawPayload`, `MetricsRecorder`, `ErrorSink`, `Unsubscribe`, `TypedHandler` |
| `event/middleware` | `PublishHandler`, `ConsumeHandler`, `PublishMiddleware`, `ConsumeMiddleware`, `ChainPublish`, `ChainConsume`, `TraceIDFromContext`, `IncomingTraceIDFromContext`, `WithTraceID`, `WithIncomingTraceID`, and order constants (`OrderLogging`, `OrderTracing`, `OrderMetrics`, `OrderRecover`, `OrderInbox`) |
| `event/transport` | `Frame`, `Delivery`, `Capabilities`, `SubscribeConfig`, `Transport`, `TxTransport`, `ConsumeFunc`, `Unsubscribe`, `ErrSubscribeUnsupported`, `EventTypePattern` |
| `event/inbox` | `Status`, `StatusProcessing`, `StatusCompleted`, `AcquireResult`, `AcquireResultAcquired`, `AcquireResultCompleted`, `AcquireResultInProgress`, `Record`, `Repository`, `ErrInProgress`, `ErrLockLost`, `ErrMissingLockID`, `ErrUnknownAcquireResult` |
| `event/transport/memory` | `Name`, `Config`, `FullPolicy`, `FullPolicyError`, `FullPolicyBlock`, `FullPolicyDropOldest` |
| `event/transport/outbox` | `Name`, `Config`, `Status`, `StatusPending`, `StatusProcessing`, `StatusCompleted`, `StatusFailed`, `StatusDead`, `Record`, `Repository` |
| `event/transport/redisstream` | `Name` and `Config` |

`ChainPublish` and `ChainConsume` sort middleware by ascending `Order`; equal
orders preserve registration order. Built-in cron jobs are named
`vef:event:outbox:relay`, `vef:event:outbox:cleanup`, and
`vef:event:inbox:cleanup`.

### Wire Values and Records

Inbox status and acquire-result values are persisted as strings:

| API | Wire value | Meaning |
| --- | --- | --- |
| `inbox.StatusProcessing` | `processing` | delivery is currently leased by a consumer |
| `inbox.StatusCompleted` | `completed` | handler completed; duplicate deliveries are acknowledged without rerunning business code |
| `inbox.AcquireResultAcquired` | `acquired` | caller owns the delivery and should run the handler |
| `inbox.AcquireResultCompleted` | `completed` | delivery was already completed |
| `inbox.AcquireResultInProgress` | `in_progress` | another consumer still owns a non-expired lease |

`event/inbox.Record` exposes JSON fields `eventId`, `consumerGroup`, `status`, `lockId`, `lockedUntil`, and `completedAt`, in addition to embedded ORM model fields.

Outbox status values are also persisted as strings:

| API | Wire value | Meaning |
| --- | --- | --- |
| `outbox.StatusPending` | `pending` | awaiting first dispatch |
| `outbox.StatusProcessing` | `processing` | currently leased by a relay worker |
| `outbox.StatusCompleted` | `completed` | downstream sink accepted the frame |
| `outbox.StatusFailed` | `failed` | most recent dispatch failed and the row is scheduled for retry |
| `outbox.StatusDead` | `dead` | retry budget was exhausted and DLQ forwarding succeeded |

`event/transport/outbox.Record` exposes JSON fields `eventId`, `eventType`, `source`, `traceId`, `spanId`, `correlationId`, `headers`, `payload`, `status`, `retryCount`, `lastError`, `processedAt`, `retryAfter`, and `occurredAt`, in addition to embedded ORM model fields.

## Inbox Idempotency

For at-least-once transports, the Inbox middleware persists a record per `(envelope_id, consumer_group)` and short-circuits future deliveries that have already completed.

```toml
[vef.event.inbox]
retention        = "168h"   # default 7 days
processing_lease = "10m"    # how long a worker may hold an in-flight claim
cleanup_interval = "1h"
```

The bus validates `inbox.retention` against the worst-case outbox backoff horizon at boot — see `config.ErrInboxRetentionTooShort`. A retention window shorter than the exponential-backoff sum could prune a dedupe record before its last retry arrives.

When a duplicate arrives while another consumer still holds the processing lease, the Inbox middleware returns `inbox.ErrInProgress` so the at-least-once transport retries later. A completed duplicate is acknowledged without running the handler again.

## Transactional Outbox

When `[vef.event.transports.outbox]` is enabled, the outbox transport persists records in the same transaction as the business write, and a relay goroutine forwards them to the configured `sink` transport (typically `memory` or `redis_stream`).

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

When `sink` is omitted, it defaults to `memory`. Use `redis_stream` as the sink
when events must cross process boundaries.

The public Go package `event/transport/outbox.Config` contains `RelayInterval`,
`MaxRetries`, `BatchSize`, `LeaseMultiplier`, `MinLease`, and `SinkName`.
The framework TOML block also has `cleanup_interval` and `completed_ttl`; those
drive the framework cleanup cron job and are not fields on the transport package
`Config`.

Producers call:

```go
err := bus.Publish(ctx, evt, event.WithTx(tx))
```

If the transaction commits, the relay eventually forwards the frame; if it rolls back, the row disappears.

Relay failures retry with exponential backoff (`2^retryCount` seconds, capped at 1h). Once the retry budget is exhausted, the relay forwards the original frame once to DLQ topic `vef-dlq.<eventType>` with header `vef.dlq=1`.

- If DLQ forwarding fails, the row remains `failed` and claimable so the DLQ hand-off can be retried.
- If DLQ forwarding succeeds, the row becomes `dead` and is kept for diagnostics.
- The cleanup job deletes completed rows older than `completed_ttl`; dead rows are retained.
- The persisted `lastError` is scrubbed for common credential fragments and truncated to 256 bytes.

The bus emits canonical JSON frames, which the outbox can persist directly. Custom code that publishes directly into an outbox transport must supply JSON-shaped frame bodies. Frames already carrying the `vef.dlq` header are refused by the outbox loop guard so a misconfigured sink cannot re-persist its own DLQ traffic indefinitely.

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

Subscribers must supply a stable `WithGroup` — the group becomes the Redis XGROUP and survives restarts.

Redis Streams contract details:

- The transport is registered only when `enabled = true` and a Redis client is available; `Start` verifies the connection with `PING`.
- Stream keys are `stream_prefix + eventType`; `max_len_approx` uses Redis `XADD MAXLEN ~`.
- `start_id = "0"` means a newly created group receives existing backlog. Set `start_id = "$"` for fire-and-forget topics where a new group should skip old messages.
- `consumer_id` is a human-readable prefix only; the runtime appends a UUID suffix so replicas do not collide.
- Missing, non-string, oversized, or invalid-JSON frames are treated as poison messages: the transport logs, `XACK`s, and drops them. At-least-once delivery applies to well-formed frames.
- Handler failures leave the message pending. The reaper periodically `XCLAIM`s idle pending entries according to `claim_idle`, `claim_interval`, and `claim_batch_size`.
- `reaper_concurrency` bounds how many subscriptions are reclaimed in parallel per cycle. `handler_timeout` bounds each fresh delivery and reaper redelivery; `0s` disables the deadline. `setup_timeout` bounds consumer-group creation during `Subscribe`.

The public Go package `event/transport/redisstream.Config` contains
`StreamPrefix`, `MaxLenApprox`, `BlockTimeout`, `ClaimIdle`, `ClaimInterval`,
`ClaimBatchSize`, `ReaperConcurrency`, `HandlerTimeout`, `SetupTimeout`,
`ConsumerID`, and `StartID`. `StreamPrefix` defaults to `vef:events:`,
`BlockTimeout` to `5s`, `ClaimIdle` to `60s`, `ClaimInterval` to `30s`,
`ClaimBatchSize` to `64`, `ReaperConcurrency` to `4`, `SetupTimeout` to
`5s`, and `StartID` to `0`. `HandlerTimeout` defaults to `0`, which disables
the per-handler deadline.

## Error Sentinels

| Error | Meaning |
| --- | --- |
| `event.ErrBusNotStarted` | publish before `Start` |
| `event.ErrBusAlreadyStarted` | duplicate `Start` |
| `event.ErrTxRequired` | `WithTx` used but no `TxTransport` in the resolved route |
| `event.ErrTransportNotFound` | transport lookup failed at a framework boundary |
| `event.ErrAsyncQueueFull` | async fan-in queue was full and the fallback synchronous publish also failed; reported via `ErrorSink`, not returned |
| `event.ErrQueueFull` | transport rejected publish under non-blocking policy |
| `event.ErrHandlerPanic` | recovered panic in a subscriber |
| `event.ErrShutdownTimeout` | bus did not drain within the graceful deadline |
| `event.ErrNoRouteMatched` | route resolved to no transports, or Subscribe resolved only to publish-only transports |
| `event.ErrUnknownPayload` | `SubscribeTyped` got an envelope it could not decode into `T` |
| `event.ErrPayloadTooLarge` | payload or headers exceed framework size limits |
| `event.ErrInvalidEventType` | event type contains characters outside `EventTypePattern` |
| `event.ErrNilTypeParameter` | `SubscribeTyped` instantiated with a nil-interface type parameter |
| `event.ErrGroupRequired` | at-least-once subscription missing `WithGroup` |
| `event.ErrTxAsyncMutex` | `WithTx` and `WithAsync` combined |

## Built-In Framework Event Types

| Event type | Source |
| --- | --- |
| `vef.api.request.audit` | API audit |
| `vef.security.login` | login flow |
| `vef.security.role_permissions.changed` | role-permission cache invalidation |
| `vef.storage.file.claimed` | a pending upload claim was adopted by a business transaction |
| `vef.storage.file.deleted` | the storage delete worker drained a file from the backend |
| `vef.storage.delete.dead_letter` | the delete worker exhausted retries; the queue row was retired and the event carries the investigation details |
| `vef.approval.task.created` | approval task created (v0.25 — emitted on every task-creation path) |

If the approval module is enabled, additional approval-domain events are published on top of these — see [Approval module](../modules/approval).

## Typical Wiring Pattern

Subscribers are typically registered through an integration module:

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

## When To Pick Which Transport

| Need | Transport |
| --- | --- |
| in-process fanout, low latency, fire-and-forget | `memory` |
| publish must commit with business write | `outbox` (→ sink) |
| cross-process delivery, at-least-once | `redis_stream` |
| reliable approval / saga events | `outbox` + `redis_stream` sink |

## Next Step

Read [Cache](./cache) to connect event publishing with invalidation flows, or [Approval module](../modules/approval) for the canonical transactional outbox example.
