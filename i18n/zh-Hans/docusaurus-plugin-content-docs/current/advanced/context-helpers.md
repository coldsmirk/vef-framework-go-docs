# 上下文辅助函数

在大多数情况下，应用代码应该优先使用 handler 参数注入，而不是手动从 context 里取值。`contextx` 包主要用于那些“确实只能拿到 `context.Context`”的低层场景。

## 可以取到什么

`contextx/contextx.go` 暴露了这些请求级数据：

- request ID
- 当前 principal
- 当前 logger
- 当前 `orm.DB`
- 当前 data permission applier

## 优先用参数注入

对于 handler 来说，这种写法通常更好：

```go
func (r *UserResource) FindPage(ctx fiber.Ctx, db orm.DB, principal *security.Principal) error
```

而不是：

```go
db := contextx.DB(ctx)
principal := contextx.Principal(ctx)
```

因为函数签名本身就把依赖关系说明白了。

## 什么场景适合 `contextx`

当你处在这些位置时，`contextx` 会更有价值：

- 被多个入口复用的 service
- 只接收 `context.Context` 的辅助函数
- 不能直接依赖 Fiber 的底层逻辑

## 示例

```go
package auditctx

import (
  "context"

  "github.com/coldsmirk/vef-framework-go/contextx"
)

func AuditContext(ctx context.Context) string {
  return contextx.RequestID(ctx)
}
```

## 这些值是谁放进去的

请求处理过程中，这些值是由框架中间件写入的：

- auth middleware 负责解析 principal
- contextual middleware 负责写入 request-scoped DB 和 logger
- data permission middleware 负责写入 request-scoped applier

而全局 app middleware 还会在 API 分发前先写入 request ID 和 logger。

## 下一步

如果你希望避免反复手动从 context 里取值，继续阅读 [自定义参数解析器](./custom-param-resolvers)。
