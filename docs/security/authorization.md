---
sidebar_position: 2
---

# Authorization

Authentication tells VEF who the caller is. Authorization decides what that caller is allowed to do.

## Permission Checks In Operations

The most common authorization entry point is `PermToken` on an operation:

```go
crud.NewUpdate[User, UserParams]().
	PermToken("sys:user:update")
```

When the operation runs, the API auth middleware extracts the permission token and asks the configured permission checker whether the current principal is allowed.

## The Main Interfaces

The important application-owned dependency is usually:

- `security.RolePermissionsLoader`

The security module already constructs a default RBAC-style `security.PermissionChecker`. You normally provide your own `PermissionChecker` only when you want to replace that behavior entirely.

Applications commonly provide:

- `security.RolePermissionsLoader`

The built-in checker depends on the role-permission loader. In practice, that means:

- if your operations use `PermToken(...)`
- and you rely on the default RBAC checker
- then you must provide a working `security.RolePermissionsLoader`

Treat that loader as required unless you intentionally replace the default checker. Without it, the default RBAC permission path does not have a valid permission source behind it.

## Resource-Level Meaning

Permission tokens should describe the action from the application's point of view, not the transport shape.

Good examples:

- `sys:user:query`
- `sys:user:create`
- `approval:delegation:update`

These tokens stay stable even if the exact request payload changes.

## What Happens On Failure

If permission checking fails, VEF returns an access-denied response. The framework preserves the structured result shape and maps the failure to the correct authorization error code.

## Practical Advice

- define permission tokens per business action
- attach them at the operation level
- keep authentication and authorization separate in your mental model
- let handlers assume authorization has already happened

## Next Step

Read [Data Permissions](./data-permissions) for row-level filtering and request-scoped data access control.
