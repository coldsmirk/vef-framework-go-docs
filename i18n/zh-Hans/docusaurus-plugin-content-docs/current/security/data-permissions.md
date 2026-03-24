---
sidebar_position: 3
---

# 数据权限

VEF 支持请求级数据权限，让你不仅能控制“用户能不能调用这个操作”，还能控制“用户能访问哪些行数据”。

## 主要组成

数据权限体系围绕这些接口展开：

- `security.DataPermissionResolver`
- `security.DataScope`
- `security.DataPermissionApplier`

最小 scope 示例：

```go
scope := security.NewSelfDataScope(orm.ColumnCreatedBy)
```

请求处理中，API 中间件会先根据当前 principal 和当前 permission token 向 resolver 请求一个数据范围，然后把请求级 applier 放进上下文中。

## CRUD 是如何使用它的

很多 CRUD 操作都可以自动应用数据权限过滤：

- 读操作会在构建查询时套入过滤
- 更新 / 删除操作会在修改前应用限制

这意味着在很多场景里，handler 本身根本不需要知道行级权限细节。

但要注意，自动过滤仍然取决于你的应用是否真的提供了可用的数据权限来源。如果默认 RBAC resolver 背后没有有效的 permission loader，那最终也可能拿不到可应用的 scope。

如果你依赖默认的 RBAC 数据权限解析器，那么 `security.RolePermissionsLoader` 基本可以视为真正让行级过滤生效的前提条件。

## 如何关闭自动数据权限

部分 CRUD builder 暴露了 `DisableDataPerm()`。

只有当你非常清楚为什么不应该自动过滤时才应该使用它。因为一旦关闭，你就要自己保证后续的边界控制仍然正确。

## 为什么这和授权不同

普通授权回答的是：

> 这个 principal 能不能调用这个操作？

数据权限回答的是：

> 这个操作一旦允许执行，它能接触到哪些记录？

VEF 把这两个层次清楚地分开了。

## 常见数据范围场景

典型例子包括：

- 系统管理员可见全部数据
- 普通用户只能看自己创建的数据
- 按部门限制
- 按租户限制

其中 `SelfDataScope` 就是最小也最直观的一个例子；部门级、组织级这类更复杂的数据范围通常由应用自己实现。

具体策略由应用自己的 resolver 和 scope 实现负责。

## 实践建议

- 数据范围规则尽量放在 handler 外部
- 尽量让 CRUD 自动应用行级过滤
- 只有在能明确说明替代方案时才使用 `DisableDataPerm()`

## 下一步

接下来可以进入 [功能特性](../features/cache) 或 [进阶](../advanced/transactions)，看平台能力和更深层的扩展面。
