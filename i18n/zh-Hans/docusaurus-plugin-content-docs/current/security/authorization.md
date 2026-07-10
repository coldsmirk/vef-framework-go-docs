---
sidebar_position: 3
---

# 授权

认证回答的是“调用者是谁”，授权回答的是“调用者能做什么”。

## 操作级权限检查

最常见的授权入口，就是在操作上配置 `RequiredPermission`（v0.24 起从 `PermToken` 重命名）：

```go
crud.NewUpdate[User, UserParams]().
	RequiredPermission("sys:user:update")
```

当操作执行时，API auth 中间件会取出这个 permission token，并交给当前配置的 permission checker 判断当前 principal 是否有权访问。

## 主要接口

应用最关键、最常需要自己提供的其实是：

- `security.RolePermissionsLoader`

安全模块本身已经会构造一个默认的 RBAC 风格 `security.PermissionChecker`。只有当你要完全替换这套默认行为时，才通常需要自己提供 `PermissionChecker`。

因此，应用里最常见的自定义点是：

- `security.RolePermissionsLoader`

内置 checker 依赖 role-permission loader。换句话说：

- 如果你的接口用了 `RequiredPermission(...)`
- 并且你依赖的是默认 RBAC checker
- 那你就必须提供一个可工作的 `security.RolePermissionsLoader`

除非你明确替换掉默认 checker，否则就应该把这个 loader 视为必需项。没有它，默认 RBAC 权限链路背后就没有可靠的权限来源。

## 授权相关公开 API

| API 组 | 公开 surface |
| --- | --- |
| 权限 | `PermissionChecker`, `RolePermissionsLoader`, `CachedRolePermissionsLoader`, `NewCachedRolePermissionsLoader` |
| 缓存失效 | `RolePermissionsChangedEvent`, `PublishRolePermissionsChangedEvent` |
| 用户信息 | `UserInfo`, `UserInfoLoader`, `UserMenu`, `UserMenuType`, `Gender` |
| 登录审计事件 | `LoginEvent`, `LoginEventParams`, `NewLoginEvent`, `SubscribeLoginEvent` |
| 认证/授权失败 | `ErrPrincipalInvalid(...)`, `ErrCredentialsInvalid(...)`, `ErrUnauthenticated`, `ErrCodePrincipalInvalid`, `ErrCodeCredentialsInvalid`，以及 permission check 返回的 access-denied 结果 |

`CachedRolePermissionsLoader` 会监听 `vef.security.role_permissions.changed`
事件。角色权限关系发生变化时发布该事件，默认 RBAC checker 才能刷新缓存授权。
`RolePermissionsChangedEvent` 的 JSON 字段是 `roles`；empty `roles`，也就是空
`roles` 数组，表示全部角色授权缓存都要失效。
如果无法订阅这个失效事件 bus，`NewCachedRolePermissionsLoader` 会 panic。

默认 RBAC permission checker 在 `principal is nil`、`no roles`，或未配置
`RolePermissionsLoader` 时返回 false。它会按顺序加载角色权限，只要任意角色的权限 map
包含当前操作的 permission token，就允许访问。

默认 RBAC data-permission resolver 也按顺序加载角色。当多个角色对同一个 permission token
提供不同 data scope 时，`DataScope.Priority()` 值最高的 scope 会胜出，也就是选择
highest priority scope。

`LoginEvent` 的 event type 是 `vef.security.login`。它的 JSON 字段是
`authType`、`userId`、`username`、`loginIp`、`userAgent`、`traceId`、`isOk`、
`failReason` 和 `errorCode`。`SubscribeLoginEvent` 会注册 typed handler，并返回
unsubscribe function。

`UserInfo` 是 `security/auth.get_user_info` 返回的结构。`Gender` 的取值包括
`GenderMale` (`male`)、`GenderFemale` (`female`)、`GenderUnknown` (`unknown`)。
`UserMenuType` 的取值包括 `UserMenuTypeDirectory` (`directory`)、
`UserMenuTypeMenu` (`menu`)、`UserMenuTypeView` (`view`)、
`UserMenuTypeDashboard` (`dashboard`)、`UserMenuTypeReport` (`report`)。

`UserInfo` 的 JSON 字段是 `id`、`name`、`gender`、`avatar`、
`permissionTokens`、`menus` 和可选 `details`。`UserMenu` 的 JSON 字段是
`type`、`path`、`name`、`icon`、可选 `meta` 和可选 `children`。

## 资源层面的意义

permission token 应该表达的是业务动作本身，而不是传输细节。

好的例子：

- `sys:user:query`
- `sys:user:create`
- `approval:delegation:update`

这样即使请求结构变化，权限语义仍然稳定。

## 权限失败时会怎样

如果权限检查失败，VEF 会返回 access denied 类型的结构化结果，同时保留对应的授权错误码。

## 实践建议

- 按业务动作定义 permission token
- 在操作层配置，而不是把权限判断散落在 handler 里
- 始终把认证和授权看成两个不同层次的问题
- 让 handler 默认假设权限已经校验完成

## 下一步

继续阅读 [数据权限](./data-permissions)，理解行级过滤和请求级数据访问控制。
