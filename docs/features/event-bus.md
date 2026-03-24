---
sidebar_position: 2
---

# Event Bus

VEF boots an in-memory event bus and exposes it through the public `event` package.

## What The Module Provides Automatically

The event module registers one in-memory bus that is exposed as:

| Interface |
| --- |
| `event.Bus` |
| `event.Publisher` |
| `event.Subscriber` |

The bus is started and stopped through the FX lifecycle automatically.

## Core Event Interfaces

### `event.Event`

Custom events should implement:

| Method | Meaning |
| --- | --- |
| `ID()` | unique event instance ID |
| `Type()` | event type string |
| `Source()` | source that produced the event |
| `Time()` | occurrence time |
| `Meta()` | metadata map |

### Publish and subscribe interfaces

| Interface | Method |
| --- | --- |
| `event.Publisher` | `Publish(event)` |
| `event.Subscriber` | `Subscribe(eventType, handler)` |
| `event.Bus` | combines publisher, subscriber, `Start()`, and `Shutdown(ctx)` |

### Middleware interfaces

| Interface or type | Purpose |
| --- | --- |
| `event.Middleware` | intercept event delivery |
| `event.MiddlewareFunc` | next function in the middleware chain |

## `event.BaseEvent`

Most custom events embed `event.BaseEvent`:

```go
type UserCreatedEvent struct {
  event.BaseEvent

  UserID string `json:"userId"`
}
```

Create the base part with:

```go
&UserCreatedEvent{
  BaseEvent: event.NewBaseEvent(
    "user.created",
    event.WithSource("user-service"),
    event.WithMeta("scope", "admin"),
  ),
  UserID: "user-1001",
}
```

Base event helpers:

| Helper | Meaning |
| --- | --- |
| `event.NewBaseEvent(type, opts...)` | creates a new base event |
| `event.WithSource(source)` | sets the source field |
| `event.WithMeta(key, value)` | adds metadata |

## Publish and Subscribe Example

```go
package userevents

import (
  "context"

  "github.com/coldsmirk/vef-framework-go/event"
)

func PublishUserCreated(publisher event.Publisher, userID string) {
  publisher.Publish(&UserCreatedEvent{
    BaseEvent: event.NewBaseEvent("user.created"),
    UserID:    userID,
  })
}

func RegisterUserCreatedHandler(subscriber event.Subscriber) event.UnsubscribeFunc {
  return subscriber.Subscribe("user.created", func(ctx context.Context, evt event.Event) {
    _ = evt
  })
}
```

## Event Middleware

The bus supports event middleware through the FX group:

```text
vef:event:middlewares
```

This is the place for cross-cutting concerns such as:

- event logging
- tracing
- filtering
- event mutation before delivery

## Built-In Framework Event Types

The framework itself currently publishes these core event types:

| Event type | Source |
| --- | --- |
| `vef.api.request.audit` | API audit event |
| `vef.security.login` | login flow event |
| `vef.storage.file.promoted` | storage promoter event |
| `vef.storage.file.deleted` | storage promoter event |
| `vef.security.role_permissions.changed` | role-permission cache invalidation event |

If the approval module is enabled, it also publishes additional approval-domain events on top of these framework-level events.

## Typical Wiring Pattern

In larger apps, subscribers are often registered through an integration module instead of inside a resource:

```go
var Module = vef.Module(
  "app:event",
  vef.Invoke(registerEventSubscribers),
)
```

That `registerEventSubscribers` function can subscribe to framework events and register lifecycle cleanup if needed.

## When To Use It

The built-in event bus is a good fit when:

- producer and consumer live in the same application
- asynchronous decoupling is enough
- you do not need an external broker yet

If you later need cross-process messaging, you can still keep the internal bus as the application-level abstraction.

## Next Step

Read [Cache](./cache) if you want to connect event publishing with invalidation or async refresh flows.
