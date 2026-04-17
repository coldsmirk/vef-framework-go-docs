---
sidebar_position: 3
---

# Context Helpers

Most of the time, you should prefer handler parameter injection over manually reading values from context. The `contextx` package is designed for lower-level code that only has access to `context.Context`.

## Overview

The `contextx` package provides type-safe getters and setters for request-scoped values. It transparently handles both `fiber.Ctx` (via `Locals`) and standard `context.Context` (via `context.WithValue`).

## API Reference

### Context Keys

| Key | Type | Description |
| --- | --- | --- |
| `KeyRequestID` | `string` | Unique request identifier |
| `KeyRequestIP` | `string` | Client IP address |
| `KeyPrincipal` | `*security.Principal` | Authenticated user principal |
| `KeyLogger` | `logx.Logger` | Request-scoped logger |
| `KeyDB` | `orm.DB` | Request-scoped database connection |
| `KeyDataPermApplier` | `security.DataPermissionApplier` | Data permission applier |

### Request ID

```go
// Get the request ID
id := contextx.RequestID(ctx) // Returns "" if not set

// Set the request ID
ctx = contextx.SetRequestID(ctx, "req-abc-123")
```

### Request IP

```go
// Get the client IP address
ip := contextx.RequestIP(ctx) // Returns "" if not set

// Set the client IP
ctx = contextx.SetRequestIP(ctx, "192.168.1.1")
```

### Principal (Current User)

```go
// Get the authenticated principal
principal := contextx.Principal(ctx) // Returns nil if not authenticated

// Set the principal
ctx = contextx.SetPrincipal(ctx, principal)

// Common usage
if p := contextx.Principal(ctx); p != nil {
    userID := p.ID
    tenantID := p.TenantID
}
```

### Logger

```go
// Get the request-scoped logger
logger := contextx.Logger(ctx)

// With fallback loggers (returns the first non-nil)
logger := contextx.Logger(ctx, fallbackLogger1, fallbackLogger2)

// Set the logger
ctx = contextx.SetLogger(ctx, logger)
```

The logger stored in context is enriched with request identity (request ID, user ID) and operation identity, making it useful for request-correlated logging in deeper service layers.

### Database (orm.DB)

```go
// Get the request-scoped DB
db := contextx.DB(ctx)

// With fallback (useful when context DB may not be set)
db := contextx.DB(ctx, globalDB)

// Set the DB
ctx = contextx.SetDB(ctx, db)
```

> **Important**: The request-scoped `orm.DB` is different from a raw global DB instance. It includes operator information (current user) used for automatic audit field population (`created_by`, `updated_by`).

### Data Permission Applier

```go
// Get the data permission applier
applier := contextx.DataPermApplier(ctx) // Returns nil if not set

// Set the applier
ctx = contextx.SetDataPermApplier(ctx, applier)
```

## Fiber / stdlib Transparency

The `contextx` package handles both Fiber and standard contexts transparently:

| Context Type | Get | Set |
| --- | --- | --- |
| `fiber.Ctx` | `ctx.Value(key)` → type assert | `ctx.Locals(key, value)` |
| `context.Context` | `ctx.Value(key)` → type assert | `context.WithValue(ctx, key, value)` |

This means the same `contextx` calls work whether you're in a Fiber handler or in a service layer using standard `context.Context`.

## When to Use Handler Injection Instead

For API handlers, prefer direct parameter injection:

```go
// ✅ Preferred: explicit dependencies in signature
func (r *UserResource) FindPage(ctx fiber.Ctx, db orm.DB, principal *security.Principal) error {
    // ...
}

// ❌ Avoid: hidden dependencies via context
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
| Request ID | App middleware (before routing) |
| Request IP | App middleware (before routing) |
| Logger | App middleware + contextual middleware |
| Principal | Auth middleware |
| DB | Contextual middleware |
| DataPermApplier | Data permission middleware |

## Next Step

Read [Custom Param Resolvers](./custom-param-resolvers) if you want to avoid repeated context access by extending handler injection directly.
