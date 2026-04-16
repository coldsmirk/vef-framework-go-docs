---
sidebar_position: 3
---

# 上下文辅助函数

在大多数情况下，应用代码应该优先使用 handler 参数注入，而不是手动从 context 里取值。`contextx` 包主要用于那些只能拿到 `context.Context` 的底层场景。

## 概述

`contextx` 包为请求级数据提供类型安全的 getter 和 setter。它透明地处理 `fiber.Ctx`（通过 `Locals`）和标准 `context.Context`（通过 `context.WithValue`）两种上下文。

## API 参考

### 上下文键

| 键 | 类型 | 说明 |
| --- | --- | --- |
| `KeyRequestID` | `string` | 唯一请求标识符 |
| `KeyRequestIP` | `string` | 客户端 IP 地址 |
| `KeyPrincipal` | `*security.Principal` | 已认证的用户主体 |
| `KeyLogger` | `logx.Logger` | 请求级日志器 |
| `KeyDB` | `orm.DB` | 请求级数据库连接 |
| `KeyDataPermApplier` | `security.DataPermissionApplier` | 数据权限应用器 |

### Request ID

```go
// 获取请求 ID
id := contextx.RequestID(ctx) // 未设置时返回 ""

// 设置请求 ID
ctx = contextx.SetRequestID(ctx, "req-abc-123")
```

### Request IP

```go
// 获取客户端 IP 地址
ip := contextx.RequestIP(ctx) // 未设置时返回 ""

// 设置客户端 IP
ctx = contextx.SetRequestIP(ctx, "192.168.1.1")
```

### Principal（当前用户）

```go
// 获取已认证的主体
principal := contextx.Principal(ctx) // 未认证时返回 nil

// 设置主体
ctx = contextx.SetPrincipal(ctx, principal)

// 常见用法
if p := contextx.Principal(ctx); p != nil {
    userID := p.ID
    tenantID := p.TenantID
}
```

### Logger

```go
// 获取请求级日志器
logger := contextx.Logger(ctx)

// 带回退日志器（返回第一个非 nil 的）
logger := contextx.Logger(ctx, fallbackLogger1, fallbackLogger2)

// 设置日志器
ctx = contextx.SetLogger(ctx, logger)
```

存储在上下文中的日志器已附带请求标识（请求 ID、用户 ID）和操作标识，便于在更深的服务层中进行请求关联日志记录。

### Database (orm.DB)

```go
// 获取请求级 DB
db := contextx.DB(ctx)

// 带回退（当上下文 DB 可能未设置时很有用）
db := contextx.DB(ctx, globalDB)

// 设置 DB
ctx = contextx.SetDB(ctx, db)
```

> **重要**：请求级 `orm.DB` 与全局原始 DB 实例不同。它包含操作者信息（当前用户），用于自动填充审计字段（`created_by`、`updated_by`）。

### Data Permission Applier

```go
// 获取数据权限应用器
applier := contextx.DataPermApplier(ctx) // 未设置时返回 nil

// 设置应用器
ctx = contextx.SetDataPermApplier(ctx, applier)
```

## Fiber / 标准库透明处理

`contextx` 包透明处理 Fiber 和标准上下文：

| 上下文类型 | 读取 | 写入 |
| --- | --- | --- |
| `fiber.Ctx` | `ctx.Value(key)` → 类型断言 | `ctx.Locals(key, value)` |
| `context.Context` | `ctx.Value(key)` → 类型断言 | `context.WithValue(ctx, key, value)` |

这意味着无论你是在 Fiber handler 中还是在使用标准 `context.Context` 的服务层中，相同的 `contextx` 调用都能正常工作。

## 优先用参数注入

对于 API handler，优先使用直接参数注入：

```go
// ✅ 推荐：签名中显式声明依赖
func (r *UserResource) FindPage(ctx fiber.Ctx, db orm.DB, principal *security.Principal) error {
    // ...
}

// ❌ 避免：通过 context 隐藏依赖
func (r *UserResource) FindPage(ctx fiber.Ctx) error {
    db := contextx.DB(ctx)
    principal := contextx.Principal(ctx)
    // ...
}
```

## 什么场景适合 `contextx`

适用场景：

- 被多个入口复用的 service 代码
- handler 层之下的辅助库
- 只能拿到 `context.Context`，而非完整 handler 签名
- 在深层调用栈中需要请求关联日志

## 这些值是谁设置的

框架中间件链自动填充上下文值：

| 值 | 设置者 |
| --- | --- |
| Request ID | App 中间件（路由前）|
| Request IP | App 中间件（路由前）|
| Logger | App 中间件 + 上下文中间件 |
| Principal | 认证中间件 |
| DB | 上下文中间件 |
| DataPermApplier | 数据权限中间件 |

## 下一步

如果你希望避免反复手动从 context 里取值，继续阅读 [自定义参数解析器](./custom-param-resolvers)。
