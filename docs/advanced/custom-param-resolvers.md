---
sidebar_position: 6
---

# Custom Parameter Resolvers

If the built-in handler injection surface is not enough, VEF lets you extend it.

## Two extension groups

You can add:

- request-time parameter resolvers
- startup-time factory parameter resolvers

The relevant DI groups are:

- `vef:api:handler_param_resolvers`
- `vef:api:factory_param_resolvers`

## What a resolver does

A handler parameter resolver tells the framework:

- which Go type it handles
- how to resolve a value of that type

The same idea applies to factory parameter resolvers, except those run at startup time rather than per request.

## When you need one

Custom resolvers are useful when:

- you want a domain-specific request-scoped object injected directly
- you want to inject a service wrapper derived from context
- you need a reusable handler contract across many resources

## Minimal example

```go
package tenantresolver

import (
  "reflect"

  "github.com/gofiber/fiber/v3"

  "github.com/coldsmirk/vef-framework-go/api"
)

type TenantContext struct {
  ID string
}

type TenantResolver struct{}

func (*TenantResolver) Type() reflect.Type {
  return reflect.TypeFor[TenantContext]()
}

func (*TenantResolver) Resolve(ctx fiber.Ctx) (reflect.Value, error) {
  tenant := TenantContext{ID: ctx.Get("X-Tenant-ID")}
  return reflect.ValueOf(tenant), nil
}
```

Registration example:

```go
fx.Provide(
  fx.Annotate(
    func() api.HandlerParamResolver { return &TenantResolver{} },
    fx.ResultTags(`group:"vef:api:handler_param_resolvers"`),
  ),
)
```

## Recommendation

Use custom resolvers for cross-cutting conventions, not for one-off shortcuts. If only one handler needs a value, plain function calls are usually simpler.

## Next step

Read [Parameters And Metadata](../guide/params-and-meta) if you first need to understand what the built-in request decoding layer already injects for you.
