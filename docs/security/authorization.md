---
sidebar_position: 3
---

# Authorization

Authentication tells VEF who the caller is. Authorization decides what that caller is allowed to do.

## Permission Checks In Operations

The most common authorization entry point is `RequiredPermission` on an operation (renamed from `PermToken` in v0.24):

```go
crud.NewUpdate[User, UserParams]().
	RequiredPermission("sys:user:update")
```

When the operation runs, the API auth middleware extracts the permission token and asks the configured permission checker whether the current principal is allowed.

## The Main Interfaces

The important application-owned dependency is usually:

- `security.RolePermissionsLoader`

The security module already constructs a default RBAC-style `security.PermissionChecker`. You normally provide your own `PermissionChecker` only when you want to replace that behavior entirely.

Applications commonly provide:

- `security.RolePermissionsLoader`

The built-in checker depends on the role-permission loader. In practice, that means:

- if your operations use `RequiredPermission(...)`
- and you rely on the default RBAC checker
- then you must provide a working `security.RolePermissionsLoader`

Treat that loader as required unless you intentionally replace the default checker. Without it, the default RBAC permission path does not have a valid permission source behind it.

## Public Authorization APIs

| API group | Public surface |
| --- | --- |
| permissions | `PermissionChecker`, `RolePermissionsLoader`, `CachedRolePermissionsLoader`, `NewCachedRolePermissionsLoader` |
| cache invalidation | `RolePermissionsChangedEvent`, `PublishRolePermissionsChangedEvent` |
| user info | `UserInfo`, `UserInfoLoader`, `UserMenu`, `UserMenuType`, `Gender` |
| login audit events | `LoginEvent`, `LoginEventParams`, `NewLoginEvent`, `SubscribeLoginEvent` |
| auth failures | `ErrPrincipalInvalid(...)`, `ErrCredentialsInvalid(...)`, `ErrUnauthenticated`, `ErrCodePrincipalInvalid`, `ErrCodeCredentialsInvalid`, and access-denied results from permission checks |

`CachedRolePermissionsLoader` listens for `vef.security.role_permissions.changed`
events. Publish that event when role-permission assignments change so the
default RBAC checker can refresh cached grants.
`RolePermissionsChangedEvent` serializes as JSON `roles`; an empty `roles`
array means all cached role grants are invalidated. `NewCachedRolePermissionsLoader`
panics if it cannot subscribe to the invalidation event bus.

The default RBAC permission checker returns false when the principal is nil, has
no roles, or no `RolePermissionsLoader` is configured. It checks roles
sequentially and grants access when any role's permission map contains the
operation's permission token.

The default RBAC data-permission resolver also loads roles sequentially. When
multiple roles provide the same permission token with different data scopes,
the scope with the highest `DataScope.Priority()` value wins.

`LoginEvent` publishes the event type `vef.security.login`. Its JSON fields are
`authType`, `userId`, `username`, `loginIp`, `userAgent`, `traceId`, `isOk`,
`failReason`, and `errorCode`. `SubscribeLoginEvent` registers a typed handler
for that event and returns an unsubscribe function.

`UserInfo` is the shape returned by `security/auth.get_user_info`. `Gender`
values are `GenderMale` (`male`), `GenderFemale` (`female`), and
`GenderUnknown` (`unknown`). `UserMenuType` values are
`UserMenuTypeDirectory` (`directory`), `UserMenuTypeMenu` (`menu`),
`UserMenuTypeView` (`view`), `UserMenuTypeDashboard` (`dashboard`), and
`UserMenuTypeReport` (`report`).

`UserInfo` serializes as `id`, `name`, `gender`, `avatar`,
`permissionTokens`, `menus`, and optional `details`. `UserMenu` serializes as
`type`, `path`, `name`, `icon`, optional `meta`, and optional `children`.

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
