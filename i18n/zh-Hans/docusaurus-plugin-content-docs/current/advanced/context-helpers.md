---
sidebar_position: 3
---

# 上下文辅助函数

在大多数情况下，应用代码应该优先使用 handler 参数注入，而不是手动从 context 里取值。`contextx` 包主要用于只能拿到 `context.Context` 的底层场景，或者需要跨 context 边界传递框架请求级值的集成代码。

## 概述

`contextx` 的 public surface 包含 9 个 exported constants 和 16 个 exported functions。它没有 exported types、没有 exported fields，也没有 exported methods。

这些 exported key constants 使用一个未导出的 key 类型。应用代码可以把这些常量传给 `ctx.Value(...)` 这类 API，但不能构造同类型的新 key。

## API 参考

### 上下文键

| 键 | 存储值 | 访问函数 |
| --- | --- | --- |
| `KeyRequest` | caller 如果选择存储 request container value，可使用这个 key | `contextx` 中不存在 Request 或 SetRequest accessor；这个 package 自身不读写该 key。 |
| `KeyRequestID` | `string` | `RequestID`, `SetRequestID` |
| `KeyRequestIP` | `string` | `RequestIP`, `SetRequestIP` |
| `KeyPrincipal` | `*security.Principal` | `Principal`, `SetPrincipal` |
| `KeyLogger` | `logx.Logger` | `Logger`, `SetLogger` |
| `KeyDB` | `orm.DB` | `DB`, `SetDB` |
| `KeyDataPermApplier` | `security.DataPermissionApplier` | `DataPermApplier`, `SetDataPermApplier` |
| `KeyRequestMethod` | `string` | `RequestMethod`, `SetRequestMethod` |
| `KeyRequestPath` | `string` | `RequestPath`, `SetRequestPath` |

这些 constant 的值按顺序固定为：`KeyRequest = 0`、`KeyRequestID = 1`、`KeyRequestIP = 2`、`KeyPrincipal = 3`、`KeyLogger = 4`、`KeyDB = 5`、`KeyDataPermApplier = 6`、`KeyRequestMethod = 7`、`KeyRequestPath = 8`。

### Functions

| Function | Signature | 缺失或类型不匹配时 |
| --- | --- | --- |
| `RequestID` | `contextx.RequestID(ctx context.Context) string` | 返回 `""`。 |
| `SetRequestID` | `contextx.SetRequestID(ctx context.Context, requestID string) context.Context` | 存储 `requestID`。 |
| `RequestIP` | `contextx.RequestIP(ctx context.Context) string` | 返回 `""`。 |
| `SetRequestIP` | `contextx.SetRequestIP(ctx context.Context, ip string) context.Context` | 存储 `ip`。 |
| `RequestMethod` | `contextx.RequestMethod(ctx context.Context) string` | 返回 `""`。 |
| `SetRequestMethod` | `contextx.SetRequestMethod(ctx context.Context, method string) context.Context` | 存储 `method`。 |
| `RequestPath` | `contextx.RequestPath(ctx context.Context) string` | 返回 `""`。 |
| `SetRequestPath` | `contextx.SetRequestPath(ctx context.Context, path string) context.Context` | 存储 `path`。 |
| `Principal` | `contextx.Principal(ctx context.Context) *security.Principal` | 返回 `nil`。 |
| `SetPrincipal` | `contextx.SetPrincipal(ctx context.Context, principal *security.Principal) context.Context` | 存储 `principal`。 |
| `Logger` | `contextx.Logger(ctx context.Context, fallbacks ...logx.Logger) logx.Logger` | 使用 fallbacks，然后返回 `nil`。 |
| `SetLogger` | `contextx.SetLogger(ctx context.Context, logger logx.Logger) context.Context` | 存储 `logger`。 |
| `DB` | `contextx.DB(ctx context.Context, fallbacks ...orm.DB) orm.DB` | 使用 fallbacks，然后返回 `nil`。 |
| `SetDB` | `contextx.SetDB(ctx context.Context, db orm.DB) context.Context` | 存储 `db`。 |
| `DataPermApplier` | `contextx.DataPermApplier(ctx context.Context) security.DataPermissionApplier` | 返回 `nil`。 |
| `SetDataPermApplier` | `contextx.SetDataPermApplier(ctx context.Context, applier security.DataPermissionApplier) context.Context` | 存储 `applier`。 |

string getters 在值未设置或存储类型不匹配时返回 zero value `""`。它们无法区分“未设置”和“显式设置为空字符串”。

`Principal` 和 `DataPermApplier` 在值未设置或存储类型不匹配时返回 `nil`。

`Logger` 和 `DB` 会先返回 context 中类型正确的值。只有 context 中没有对应类型时，才会检查 fallbacks。fallbacks 按从左到右扫描，并返回第一个 `reflectx.IsNotEmpty(...)` 的值，所以 nil 和 typed nil fallbacks 会被跳过。这个过滤只适用于 fallbacks：如果 typed nil 已经存进 context，类型断言成功后仍会直接返回。

### 请求标识

```go
id := contextx.RequestID(ctx) // 未设置时返回 ""
ctx = contextx.SetRequestID(ctx, "req-abc-123")

ip := contextx.RequestIP(ctx) // 未设置时返回 ""
ctx = contextx.SetRequestIP(ctx, "192.168.1.1")

method := contextx.RequestMethod(ctx) // 例如 "GET"；未设置时返回 ""
path := contextx.RequestPath(ctx)     // 例如 "/api/users"；未设置时返回 ""

ctx = contextx.SetRequestMethod(ctx, "POST")
ctx = contextx.SetRequestPath(ctx, "/api/orders")
```

### Principal（当前用户）

```go
principal := contextx.Principal(ctx) // 未认证时返回 nil
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

// context 中的值优先于 fallback。
logger := contextx.Logger(ctx, fallbackLogger1, fallbackLogger2)

ctx = contextx.SetLogger(ctx, logger)
```

框架中间件写入的请求级日志器会附带请求标识和操作标识，便于在更深的服务层中进行请求关联日志记录。

### Database (orm.DB)

```go
db := contextx.DB(ctx)

// context 中的值优先于 fallback。
db := contextx.DB(ctx, globalDB)

ctx = contextx.SetDB(ctx, db)
```

> **重要**：请求级 `orm.DB` 与全局原始 DB 实例不同。它包含操作者信息（当前用户），用于自动填充审计字段（`created_by`、`updated_by`）。

### Data Permission Applier

```go
applier := contextx.DataPermApplier(ctx) // 未设置时返回 nil
ctx = contextx.SetDataPermApplier(ctx, applier)
```

## Fiber / 标准库透明处理

所有 setter functions 都返回一个 `context.Context`，但写入行为取决于具体 context 类型：

| 上下文类型 | 读取 | 写入 |
| --- | --- | --- |
| `fiber.Ctx` | `ctx.Value(key)` 读取 Fiber request locals，然后 getter 做类型断言。 | 通过 `ctx.Locals(key, value)` 写入，并返回同一个 Fiber context。 |
| 标准 `context.Context` | `ctx.Value(key)` 读取 context values，然后 getter 做类型断言。 | 返回 `context.WithValue(ctx, key, value)`。原 context 不会被修改。 |

使用标准 context 时，必须保留返回值：

```go
ctx = contextx.SetRequestID(ctx, "req-abc-123")
```

对于 `fiber.Ctx`，setter 会原地修改 request locals，但保留返回值仍然无害，也能让调用点保持一致。

## 优先用参数注入

对于 API handler，优先使用直接参数注入：

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
| Request ID | Logger middleware 同时写入 Fiber locals 和嵌入的标准 context。 |
| Logger | Logger middleware 同时写入两条路径；contextual middleware 可能替换为带 operation scope 的 logger。 |
| Request IP | Auth middleware 在 authenticator 运行前写入嵌入的标准 context。Signature auth 也会用它做 IP whitelist 检查。 |
| Request method/path | Auth middleware 在 authenticator 运行前写入嵌入的标准 context。Signature auth 会把两者绑定进签名校验。 |
| Principal | Auth middleware 同时写入 Fiber locals 和嵌入的标准 context。 |
| DB | Contextual middleware 同时写入 Fiber locals 和嵌入的标准 context。 |
| DataPermApplier | Data permission middleware 同时写入 Fiber locals 和嵌入的标准 context。 |

## 下一步

如果你希望避免反复手动从 context 里取值，继续阅读 [自定义参数解析器](./custom-param-resolvers)。
