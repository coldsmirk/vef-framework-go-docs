---
sidebar_position: 2
---

# 扩展 Handler 参数

`api.Resource` 上的 handler 方法把输入声明为普通的 Go 参数——`ctx fiber.Ctx`、`db orm.DB`、`principal *security.Principal`、一个 `Params`/`Meta` 结构体等等。VEF 在 handler 运行时按类型解析每一个参数。本页覆盖这套机制的完整链路：参数注入是如何解析的、只能拿到 `context.Context` 的代码要用的 `contextx` 包，以及内置注入面不够用时如何注册自定义参数解析器。

## Handler 参数注入是如何工作的

对于每个声明的参数类型，handler 参数解析器 manager 会按以下顺序尝试：

1. **精确类型匹配**——在已注册的 `api.HandlerParamResolver` 中查找（下面列出的内置解析器，加上任何注册进 `group:"vef:api:handler_param_resolvers"` 的解析器）。
2. **`api.Params` / `api.Meta` embedding**——如果该类型嵌入了这两个 sentinel 类型之一（或是被识别的 meta 类型），对应的请求分段会被解码并校验进这个类型。
3. **Resource 字段兜底**——resource 结构体自身同类型的字段，让 handler 可以直接引用 resource 已经持有的依赖。

如果都不匹配，解析会失败，该 handler 也就无法被适配。

Handler factory（在启动期构建 `fiber.Handler` 的函数，例如 `func (r *UserResource) CreateHandler(service UserService) func(ctx fiber.Ctx) error`）走的是另一套启动期解析器——`api.FactoryParamResolver`，注册进 `group:"vef:api:factory_param_resolvers"`。factory resolver 只在装配阶段执行一次，而不是每个请求执行一次。

### 内置 handler 参数解析器

| 类型 | 来源 |
| --- | --- |
| `fiber.Ctx` | 请求上下文本身 |
| `orm.DB` | `contextx.DB(ctx)`——请求级 DB |
| `logx.Logger` | `contextx.Logger(ctx)`——请求级 logger |
| `*security.Principal` | `contextx.Principal(ctx)`——当前已认证的 principal |
| `cron.Scheduler` | 固定值，启动时注入 |
| `event.Bus` | 固定值，启动时注入 |
| `mold.Transformer` | 固定值，启动时注入 |
| `storage.Service` | 固定值，启动时注入 |
| `datasource.Registry` | 固定值，启动时注入 |
| `api.Params` | 从请求的 params 分段解码并校验 |
| `api.Meta` | 从请求的 meta 分段解码并校验 |

### 内置 factory 参数解析器

`orm.DB`、`cron.Scheduler`、`event.Bus`、`mold.Transformer`、`storage.Service`、`storage.Files`、`datasource.Registry`。

## 上下文辅助函数（`contextx`）

在大多数情况下，应用代码应该优先使用 handler 参数注入，而不是手动从 context 里取值。`contextx` 包主要用于只能拿到 `context.Context` 的底层场景，或者需要跨 context 边界传递框架请求级值的集成代码。

### 概述

这些 exported key constants 使用一个未导出的 key 类型。应用代码可以把这些常量传给 `ctx.Value(...)` 这类 API，但不能构造同类型的新 key。

### API 参考

#### 上下文键

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

#### Functions

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

### Fiber / 标准库透明处理

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

### 优先用参数注入

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

### 什么场景适合 `contextx`

适用场景：

- 被多个入口复用的 service 代码
- handler 层之下的辅助库
- 只能拿到 `context.Context`，而非完整 handler 签名
- 在深层调用栈中需要请求关联日志

### 这些值是谁设置的

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

## 自定义参数解析器

如果上一节的内置 handler 注入面不够用，VEF 允许你扩展它。

### 两个扩展 group

你可以添加：

- 请求期参数解析器
- 启动期 factory 参数解析器

对应的 DI group 是：

- `vef:api:handler_param_resolvers`
- `vef:api:factory_param_resolvers`

### resolver 要做什么

一个 handler 参数解析器（`api.HandlerParamResolver`）需要告诉框架：

- 它处理哪个 Go 类型（`Type() reflect.Type`）
- 如何解析出这个类型的值（`Resolve(ctx fiber.Ctx) (reflect.Value, error)`）

factory 参数解析器（`api.FactoryParamResolver`，`Type() reflect.Type` + `Resolve() (reflect.Value, error)`）的思路完全类似，只不过它是在启动期执行，而不是每个请求执行一次，也拿不到 `fiber.Ctx`。

### 什么时候需要一个

自定义 resolver 适用于：

- 需要把某个领域特定的请求级对象直接注入 handler
- 需要注入一个从 context 派生出来的 service wrapper
- 需要在大量 resource 之间复用同一份 handler 契约

### 最小示例

```go
package tenantresolver

import (
  "reflect"

  "github.com/gofiber/fiber/v3"
)

type TenantContext struct {
  ID string
}

type TenantResolver struct{}

func (*TenantResolver) Type() reflect.Type {
  return reflect.TypeFor[TenantContext]()
}

func (*TenantResolver) Resolve(ctx fiber.Ctx) (reflect.Value, error) {
  tenant := TenantContext{ID: ctx.Get("X-Tenant-ID")}
  return reflect.ValueOf(tenant), nil
}
```

这样 handler 就能直接声明这个参数：

```go
func (r *UserResource) Find(ctx fiber.Ctx, currentTenant TenantContext) error {
  // ...
}
```

### 注册示例

```go
fx.Provide(
  fx.Annotate(
    func() api.HandlerParamResolver { return &TenantResolver{} },
    fx.ResultTags(`group:"vef:api:handler_param_resolvers"`),
  ),
)
```

在你的模块中用 `vef.ProvideAPIResource` 同款的 FX 装配方式注册即可——同样的 `fx.Annotate` + `fx.ResultTags` 模式也适用于 `group:"vef:api:factory_param_resolvers"`，用于支持这类 handler factory 签名：

```go
func (r *UserResource) CreateHandler(service UserService) func(ctx fiber.Ctx) error
```

### 建议

在添加自定义 resolver 之前，先确认下面这些更简单的方式是否已经够用：

- 把依赖作为 resource 字段注入
- 直接使用已有的内置 resolver
- 通过 `Params`/`Meta` 请求结构体传值

自定义 resolver 应该用于跨切面的通用约定，而不是一次性的捷径。如果只有一个 handler 需要某个值，直接的函数调用通常更简单。只有当某个依赖在大量 handler 中反复出现时，自定义 resolver 才真的值得。

## 下一步

- 阅读 [参数与元信息（Meta）](../building-apis/params-and-meta)，了解内置请求解码层已经为你注入了什么（`api.Params`、`api.Meta`），再决定要不要上自定义 resolver。
- 阅读 [扩展点](../reference/extension-points)，查看框架里全部 DI 扩展 group 的目录，包括 `vef:api:handler_param_resolvers` 和 `vef:api:factory_param_resolvers`。
