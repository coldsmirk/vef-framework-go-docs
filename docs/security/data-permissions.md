---
sidebar_position: 3
---

# Data Permissions

VEF supports request-scoped data permissions so that authorization can narrow the rows a user may read or mutate, not just whether the endpoint is callable at all.

## The Main Pieces

The data-permission system revolves around:

- `security.DataPermissionResolver`
- `security.DataScope`
- `security.DataPermissionApplier`

Minimal scope example:

```go
scope := security.NewSelfDataScope(orm.ColumnCreatedBy)
```

During request processing, the API middleware asks the resolver for the current principal's data scope for the current permission token. It then stores a request-scoped applier in context.

## How CRUD Uses It

Many CRUD operations can automatically apply data-permission filtering:

- read operations apply it while building queries
- update/delete operations apply it before mutating records

That means handlers often do not need to know the details of row-level permission enforcement.

In practice, automatic filtering still depends on your application providing a working data-permission source. If the default RBAC resolver has no usable permission loader behind it, there may be no scope to apply.

If you rely on the default RBAC data-permission resolver, treat `security.RolePermissionsLoader` as the practical prerequisite for meaningful row-level filtering.

## Disabling Automatic Data Permission

Some CRUD builders expose `DisableDataPerm()` for cases where automatic filtering is not appropriate.

Use it carefully. If you disable data permission on a privileged endpoint, you are taking responsibility for enforcing the correct data boundary elsewhere.

## Why This Matters

Regular authorization answers:

> can this principal call this operation?

Data permission answers:

> which records can this principal touch once the operation is allowed?

These are different layers, and VEF keeps them distinct.

## Typical Use Cases

Examples of data scopes include:

- all rows for a system admin
- only self-created records
- department-scoped access
- tenant-scoped filtering

The built-in `SelfDataScope` is the smallest example of this pattern. Broader organization- or department-based scopes are usually application-defined.

The exact policy is application-owned through the resolver and scope implementation.

## Practical Advice

- keep data scope rules outside handlers
- let CRUD apply row filters whenever possible
- use `DisableDataPerm()` only when you can clearly justify the alternative

## Next Step

Move to [Features](../features/cache) or [Advanced](../advanced/transactions) depending on whether you want platform capabilities or deeper runtime extension points next.
