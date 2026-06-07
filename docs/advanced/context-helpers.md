---
sidebar_position: 3
---

# Context Helpers

Most of the time, prefer handler parameter injection over manually reading values from context. The `contextx` package is for lower-level code that only has a `context.Context`, or for integration code that needs to pass request-scoped framework values across a context boundary.

## Overview

The public `contextx` surface has 9 exported constants and 16 exported functions. It has no exported types, no exported fields, and no exported methods.

The exported key constants use an unexported key type. Application code can pass those constants to APIs such as `ctx.Value(...)`, but cannot construct additional keys of the same type.

## API Reference

### Context Keys

| Key | Stored value | Accessors |
| --- | --- | --- |
| `KeyRequest` | request container value, if a caller chooses to store one | No Request or SetRequest accessor exists in `contextx`; the package itself does not read or write this key. |
| `KeyRequestID` | `string` | `RequestID`, `SetRequestID` |
| `KeyRequestIP` | `string` | `RequestIP`, `SetRequestIP` |
| `KeyPrincipal` | `*security.Principal` | `Principal`, `SetPrincipal` |
| `KeyLogger` | `logx.Logger` | `Logger`, `SetLogger` |
| `KeyDB` | `orm.DB` | `DB`, `SetDB` |
| `KeyDataPermApplier` | `security.DataPermissionApplier` | `DataPermApplier`, `SetDataPermApplier` |
| `KeyRequestMethod` | `string` | `RequestMethod`, `SetRequestMethod` |
| `KeyRequestPath` | `string` | `RequestPath`, `SetRequestPath` |

The constant values are stable in this order: `KeyRequest = 0`, `KeyRequestID = 1`, `KeyRequestIP = 2`, `KeyPrincipal = 3`, `KeyLogger = 4`, `KeyDB = 5`, `KeyDataPermApplier = 6`, `KeyRequestMethod = 7`, and `KeyRequestPath = 8`.

### Functions

| Function | Signature | Missing or wrong-type value |
| --- | --- | --- |
| `RequestID` | `contextx.RequestID(ctx context.Context) string` | Returns `""`. |
| `SetRequestID` | `contextx.SetRequestID(ctx context.Context, requestID string) context.Context` | Stores `requestID`. |
| `RequestIP` | `contextx.RequestIP(ctx context.Context) string` | Returns `""`. |
| `SetRequestIP` | `contextx.SetRequestIP(ctx context.Context, ip string) context.Context` | Stores `ip`. |
| `RequestMethod` | `contextx.RequestMethod(ctx context.Context) string` | Returns `""`. |
| `SetRequestMethod` | `contextx.SetRequestMethod(ctx context.Context, method string) context.Context` | Stores `method`. |
| `RequestPath` | `contextx.RequestPath(ctx context.Context) string` | Returns `""`. |
| `SetRequestPath` | `contextx.SetRequestPath(ctx context.Context, path string) context.Context` | Stores `path`. |
| `Principal` | `contextx.Principal(ctx context.Context) *security.Principal` | Returns `nil`. |
| `SetPrincipal` | `contextx.SetPrincipal(ctx context.Context, principal *security.Principal) context.Context` | Stores `principal`. |
| `Logger` | `contextx.Logger(ctx context.Context, fallbacks ...logx.Logger) logx.Logger` | Uses fallbacks, then returns `nil`. |
| `SetLogger` | `contextx.SetLogger(ctx context.Context, logger logx.Logger) context.Context` | Stores `logger`. |
| `DB` | `contextx.DB(ctx context.Context, fallbacks ...orm.DB) orm.DB` | Uses fallbacks, then returns `nil`. |
| `SetDB` | `contextx.SetDB(ctx context.Context, db orm.DB) context.Context` | Stores `db`. |
| `DataPermApplier` | `contextx.DataPermApplier(ctx context.Context) security.DataPermissionApplier` | Returns `nil`. |
| `SetDataPermApplier` | `contextx.SetDataPermApplier(ctx context.Context, applier security.DataPermissionApplier) context.Context` | Stores `applier`. |

String getters return the zero value `""` when the value is unset or stored with the wrong type. They cannot distinguish "unset" from "explicitly set to an empty string".

`Principal` and `DataPermApplier` return `nil` when unset or stored with the wrong type.

`Logger` and `DB` first return a correctly typed context value. Only when the context does not contain that type do they inspect fallbacks. Fallbacks are scanned left to right and the first `reflectx.IsNotEmpty(...)` value is returned, so nil and typed nil fallbacks are skipped. This filtering applies only to fallbacks: a typed nil value already stored in the context still wins because the type assertion succeeds.

### Request Identity

```go
id := contextx.RequestID(ctx) // Returns "" if not set
ctx = contextx.SetRequestID(ctx, "req-abc-123")

ip := contextx.RequestIP(ctx) // Returns "" if not set
ctx = contextx.SetRequestIP(ctx, "192.168.1.1")

method := contextx.RequestMethod(ctx) // e.g. "GET"; "" if not set
path := contextx.RequestPath(ctx)     // e.g. "/api/users"; "" if not set

ctx = contextx.SetRequestMethod(ctx, "POST")
ctx = contextx.SetRequestPath(ctx, "/api/orders")
```

### Principal (Current User)

```go
principal := contextx.Principal(ctx) // Returns nil if not authenticated
ctx = contextx.SetPrincipal(ctx, principal)

if p := contextx.Principal(ctx); p != nil {
    userID := p.ID
    roles := p.Roles
    _ = userID
    _ = roles
}
```

### Logger

```go
logger := contextx.Logger(ctx)

// The context value wins over fallbacks.
logger := contextx.Logger(ctx, fallbackLogger1, fallbackLogger2)

ctx = contextx.SetLogger(ctx, logger)
```

Framework middleware stores request-scoped loggers that include request identity and operation identity, making them useful for request-correlated logging in deeper service layers.

### Database (orm.DB)

```go
db := contextx.DB(ctx)

// The context value wins over fallbacks.
db := contextx.DB(ctx, globalDB)

ctx = contextx.SetDB(ctx, db)
```

> **Important**: The request-scoped `orm.DB` is different from a raw global DB instance. It includes operator information (current user) used for automatic audit field population (`created_by`, `updated_by`).

### Data Permission Applier

```go
applier := contextx.DataPermApplier(ctx) // Returns nil if not set
ctx = contextx.SetDataPermApplier(ctx, applier)
```

## Fiber / stdlib Transparency

All setter functions return a `context.Context`, but the write behavior depends on the concrete context:

| Context Type | Get | Set |
| --- | --- | --- |
| `fiber.Ctx` | `ctx.Value(key)` reads Fiber request locals, then the getter type-asserts the value. | Writes with `ctx.Locals(key, value)` and returns the same Fiber context. |
| Standard `context.Context` | `ctx.Value(key)` reads context values, then the getter type-asserts the value. | Returns `context.WithValue(ctx, key, value)`. The original context is unchanged. |

Always keep the returned context when using a standard context:

```go
ctx = contextx.SetRequestID(ctx, "req-abc-123")
```

For `fiber.Ctx`, the setter mutates request locals in place, but keeping the return value is still harmless and keeps call sites consistent.

## When to Use Handler Injection Instead

For API handlers, prefer direct parameter injection:

```go
func (r *UserResource) FindPage(ctx fiber.Ctx, db orm.DB, principal *security.Principal) error {
    // ...
}

func (r *UserResource) FindPage(ctx fiber.Ctx) error {
    db := contextx.DB(ctx)
    principal := contextx.Principal(ctx)
    // ...
}
```

## When `contextx` Makes Sense

Use `contextx` when:

- You are inside service code reused by multiple entry points
- You are writing helper libraries below the handler layer
- You only have `context.Context`, not the full handler signature
- You need request-correlated logging in deep call stacks

## Who Sets These Values

The framework middleware chain populates context values automatically:

| Value | Set By |
| --- | --- |
| Request ID | Logger middleware writes both Fiber locals and the embedded standard context. |
| Logger | Logger middleware writes both paths; contextual middleware may replace it with an operation-scoped logger. |
| Request IP | Auth middleware writes the embedded standard context before authenticators run. Signature auth also uses it for IP whitelist checks. |
| Request method/path | Auth middleware writes the embedded standard context before authenticators run. Signature auth binds both values into signature verification. |
| Principal | Auth middleware writes both Fiber locals and the embedded standard context. |
| DB | Contextual middleware writes both Fiber locals and the embedded standard context. |
| DataPermApplier | Data permission middleware writes both Fiber locals and the embedded standard context. |

## Next Step

Read [Custom Param Resolvers](./custom-param-resolvers) if you want to avoid repeated context access by extending handler injection directly.
