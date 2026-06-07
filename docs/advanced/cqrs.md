---
sidebar_position: 5
---

# CQRS

VEF includes a lightweight CQRS bus with typed handlers and behavior middleware.

Audit note: this page covers 26 public CQRS entries, including 8 grouped CQRS method entries across 8 CQRS receiver/type families. The grouped CQRS surface contains 0 exported CQRS field entries and 8 exported CQRS method entries.

## Public API

The public CQRS package mainly exposes:

- `cqrs.BaseCommand`
- `cqrs.BaseQuery`
- `cqrs.Register(...)`
- `cqrs.Send(...)`
- `cqrs.Behavior`
- `cqrs.BehaviorFunc`
- `vef.ProvideCQRSBehavior(...)`

Supporting public APIs:

| API | Purpose |
| --- | --- |
| `cqrs.NewBus(behaviors)` | create a standalone bus, mostly useful in tests or custom wiring |
| `cqrs.Bus` | handler registry and dispatcher abstraction |
| `cqrs.Action` / `ActionKind` | action contract; `Command` and `Query` are the exported kind constants |
| `cqrs.Handler[TAction, TResult]` / `HandlerFunc[...]` | typed handler contracts |
| `cqrs.Behavior` / `BehaviorFunc` | command/query execution pipeline |
| `cqrs.Ordered` | optional behavior ordering hook |
| `cqrs.Unit` | empty result type for commands |
| `ErrHandlerNotFound` | no handler registered for the action type |
| `ErrResultTypeMismatch` | handler result could not be converted to the requested result type |

### Action Kind Contract

`Action.Kind()` returns the action discriminator. `BaseCommand.Kind()` returns
`Command` (`0`), and `BaseQuery.Kind()` returns `Query` (`1`).

## The bus model

Handlers are registered by action type, and sends are type-safe:

```go
package useractions

import (
  "context"

  "github.com/coldsmirk/vef-framework-go/cqrs"
)

type CreateUser struct {
  cqrs.BaseCommand
  Name string
}

type CreateUserHandler struct{}

func (CreateUserHandler) Handle(ctx context.Context, cmd CreateUser) (cqrs.Unit, error) {
  return cqrs.Unit{}, nil
}

func RegisterHandlers(bus cqrs.Bus) {
  cqrs.Register(bus, CreateUserHandler{})
}

func Run(ctx context.Context, bus cqrs.Bus) error {
  _, err := cqrs.Send[CreateUser, cqrs.Unit](ctx, bus, CreateUser{Name: "alice"})
  return err
}
```

`Register` uses the concrete `TAction` type as the registry key and panics if
another handler is already registered for the same action type. `Send` dispatches
by that same action type. If no handler exists, it returns an error matching
`ErrHandlerNotFound` and a zero-value result. Handler errors are propagated
unchanged.

`HandlerFunc` is a function adapter for `Handler`; its `Handle` method simply
calls the wrapped function.

## Behavior pipeline

The CQRS bus supports middleware-like behaviors around command/query execution.

Use `vef.ProvideCQRSBehavior(...)` to register them into the runtime.

This is the right place for:

- tracing
- logging
- metrics
- cross-cutting validation

Minimal behavior example:

```go
func NewLoggingBehavior() cqrs.Behavior {
  return cqrs.BehaviorFunc(func(ctx context.Context, action cqrs.Action, next func(context.Context) (any, error)) (any, error) {
    return next(ctx)
  })
}
```

`BehaviorFunc` is a function adapter for `Behavior`; its `Handle` method simply
calls the wrapped function. A behavior receives the original action and can
short-circuit without calling `next`. If it short-circuits with `nil`, `Send`
returns the zero value for `TResult`. If it short-circuits with a non-nil value
whose concrete type is not `TResult`, `Send` returns an error matching
`ErrResultTypeMismatch` instead of panicking.

Behaviors are sorted once when the bus is built. `Ordered.Order()` controls
wrapping order: lower values wrap outside higher values. Behaviors that do not
implement `Ordered` default to order `0`; equal orders preserve the input order
for a standalone `NewBus`, but FX value-group ordering is not stable. The
framework reserves these conventional bands:

| Order band | Use |
| --- | --- |
| `0..99` | transactional / contextual setup that must wrap everything |
| `100..199` | audit / collector lifecycle |
| `200..299` | event publish / outbox side effects |
| `1000+` | custom host behaviors |

## When to use it

CQRS is optional. It is most useful when you want:

- explicit command/query boundaries
- type-safe dispatch
- a pipeline around application actions outside the HTTP resource layer

## Next step

Read [Transactions](./transactions) if your behaviors need to wrap command execution in database boundaries.
