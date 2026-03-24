---
sidebar_position: 3
---

# Context Helpers

Most of the time, you should prefer handler parameter injection over manually reading values from context.

## The public helper package

The `contextx` package provides helpers for request-scoped values such as:

- request ID
- principal
- logger
- `orm.DB`
- data permission applier

## When to use handler injection instead

If you are writing an API handler, these are usually better as direct parameters:

- `orm.DB`
- `log.Logger`
- `*security.Principal`
- `fiber.Ctx`

This keeps the handler signature honest and avoids hidden dependencies.

## When direct context access makes sense

`contextx` is more useful when:

- you are inside code that does not participate in handler parameter injection
- you are writing helper libraries used below the handler layer
- you are working with contexts outside the immediate resource method signature

## Minimal example

```go
package auditctx

import (
  "context"

  "github.com/coldsmirk/vef-framework-go/contextx"
)

func AuditContext(ctx context.Context) string {
  return contextx.RequestID(ctx)
}
```

## Request-scoped DB

The API middleware chain injects a request-scoped `orm.DB` that includes operator information. That is why reading `contextx.DB(ctx)` is not the same as using a raw global DB instance.

## Request-scoped logger

The logger stored in context is also enriched with request identity and, for API requests, operation identity.

That makes it useful in deeper service layers when you need request-correlated logs without threading a logger parameter through every function yourself.

## Next step

Read [Custom Param Resolvers](./custom-param-resolvers) if you want to avoid repeated context access by extending handler injection directly.
