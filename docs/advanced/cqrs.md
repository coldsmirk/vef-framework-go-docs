---
sidebar_position: 5
---

# CQRS

VEF includes a lightweight CQRS bus with typed handlers and behavior middleware.

## Public API

The public CQRS package mainly exposes:

- `cqrs.Register(...)`
- `cqrs.Send(...)`
- `cqrs.Behavior`
- `cqrs.BehaviorFunc`

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

## When to use it

CQRS is optional. It is most useful when you want:

- explicit command/query boundaries
- type-safe dispatch
- a pipeline around application actions outside the HTTP resource layer

## Next step

Read [Transactions](./transactions) if your behaviors need to wrap command execution in database boundaries.
