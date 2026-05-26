---
sidebar_position: 2
---

# Event Bus

VEF ships a pluggable multi-transport event platform behind a single `event.Bus` facade. The bus is auto-wired by FX — applications register events, publish, subscribe, and let the configured routing decide where each frame goes.

> The event system was rewritten in v0.21 as a pluggable platform (`feat(event)!: rewrite event system as pluggable multi-transport platform`) and has been hardened repeatedly since: typed delivery (`SubscribeTyped`), capability-aware routing, transactional outbox, Redis Streams transport, Inbox dedupe, route inspection, and explicit error sentinels. This page reflects the v0.26 surface.

## What FX Wires For You

Once the framework starts, the following are available for injection:

| Interface | Purpose |
| --- | --- |
| `event.Bus` | combined publish/subscribe entry point |
| `event.RouteInspector` | read-only routing queries — used to fail fast at OnStart |
| `transport.Transport` (FX group `vef:event:transports`) | every registered transport (memory, outbox, redis-stream, custom) |
| `event.ErrorSink` | sink for out-of-band publish errors (default logs at error) |

`Bus.Start` / `Bus.Shutdown` are driven by FX. Publishing during `fx.Provide` returns `event.ErrBusNotStarted` — move the call into a `fx.Invoke` or lifecycle hook.

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
- `Subscribe` registrations made before `Start` are buffered and flushed during boot, so order of FX wiring does not matter.

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

Publish-time options compose left-to-right:

| Option | Effect |
| --- | --- |
| `event.WithTx(tx orm.DB)` | route through a `TxTransport`. Required for the transactional outbox pattern; returns `event.ErrTxRequired` if the resolved route has no `TxTransport`. |
| `event.WithAsync()` | hand the publish to the bus async fan-in queue and return immediately. Errors flow to `ErrorSink`, not the caller. |
| `event.WithSource(name)` | override `Envelope.Source` (defaults to `vef.app.name`). |
| `event.WithOccurredAt(t)` | override `Envelope.OccurredAt` (defaults to `time.Now`). |
| `event.WithCorrelationID(id)` | caller-controlled correlation key. |
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
| `event.WithGroup(name)` | consumer group. **Required** when the route resolves to any at-least-once transport (outbox-sink, Redis Streams) — otherwise `event.ErrGroupRequired`. The group is the dedupe scope for the Inbox middleware and the XGROUP for Redis Streams; it must stay stable across restarts. |
| `event.WithConcurrency(n)` | worker count per subscription. Defaults to 1. |

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

T can be a pointer type (recommended) or a value type whose pointer also satisfies `Event`. The bus accepts both in-process delivery (payload already typed) and cross-process delivery (`RawPayload` with canonical JSON body).

## Envelope and Frame

The framework wraps each `Event` in an `Envelope` carrying transport metadata: `ID`, `Type`, `Source`, `OccurredAt`, `PublishedAt`, `TraceID` / `SpanID`, `CorrelationID`, `Headers`, and `Payload`. Cross-process transports serialize the body into a `transport.Frame` whose `Body` is canonical JSON; the bus decodes it back to `Envelope.Payload = RawPayload{...}` so `SubscribeTyped[T]` can deserialize.

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

## Routing

Routing is declarative — `[vef.event.routing]` rules in TOML are matched top-to-bottom using `path.Match` semantics (`*`, `?`, `[abc]`). The first matching rule wins; fan-out is expressed by listing multiple transports.

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

When no rule matches, `default_transport` is used. Publishing an unrouted type returns `event.ErrNoRouteMatched`.

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
| metrics | counters and latency histograms | pluggable `MetricsRecorder` |
| recover | converts panics into `event.ErrHandlerPanic` | toggled by `recover` |
| inbox | consume-side idempotency dedupe | toggled by `inbox`; attaches **only** when the transport's `Capabilities.AtLeastOnce` is true |

## Inbox Idempotency

For at-least-once transports, the Inbox middleware persists a record per `(envelope_id, consumer_group)` and short-circuits future deliveries that have already completed.

```toml
[vef.event.inbox]
retention        = "168h"   # default 7 days
processing_lease = "10m"    # how long a worker may hold an in-flight claim
cleanup_interval = "1h"
```

The bus validates `inbox.retention` against the worst-case outbox backoff horizon at boot — see `config.ErrInboxRetentionTooShort`. A retention window shorter than the exponential-backoff sum could prune a dedupe record before its last retry arrives.

## Transactional Outbox

When `[vef.event.transports.outbox]` is enabled, the outbox transport persists records in the same transaction as the business write, and a relay goroutine forwards them to the configured `sink` transport (typically `memory` or `redis_stream`).

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

Producers call:

```go
err := bus.Publish(ctx, evt, event.WithTx(tx))
```

If the transaction commits, the relay eventually forwards the frame; if it rolls back, the row disappears.

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

Subscribers must supply a stable `WithGroup` — the group becomes the Redis XGROUP and survives restarts.

## Error Sentinels

| Error | Meaning |
| --- | --- |
| `event.ErrBusNotStarted` | publish before `Start` |
| `event.ErrBusAlreadyStarted` | duplicate `Start` |
| `event.ErrTxRequired` | `WithTx` used but no `TxTransport` in the resolved route |
| `event.ErrTransportNotFound` | routing referenced an unknown transport name |
| `event.ErrAsyncQueueFull` | async fan-in queue full (reported via `ErrorSink`) |
| `event.ErrQueueFull` | transport rejected publish under non-blocking policy |
| `event.ErrHandlerPanic` | recovered panic in a subscriber |
| `event.ErrShutdownTimeout` | bus did not drain within the graceful deadline |
| `event.ErrNoRouteMatched` | no routing rule (or default) matched the event type |
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
| `vef.storage.delete.dead_letter` | the delete worker exhausted retries; the row is parked for manual investigation |
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
