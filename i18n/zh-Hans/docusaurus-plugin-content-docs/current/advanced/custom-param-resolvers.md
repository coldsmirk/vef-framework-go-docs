# 自定义参数解析器

VEF 不只支持注入 `fiber.Ctx` 和请求结构体。如果内置参数集合不够，你可以自己注册 handler 参数解析器和 factory 参数解析器。

## 内置 resolver 已经支持什么

当前内部 API param 模块已经注册了这些 handler resolver：

- `fiber.Ctx`
- `orm.DB`
- `log.Logger`
- `*security.Principal`
- `cron.Scheduler`
- `event.Publisher`
- `mold.Transformer`
- `storage.Service`
- `api.Params`
- `api.Meta`

同时也注册了 handler factory 需要的 factory resolver。

## 什么时候需要扩展

如果你希望写出这种 handler 签名：

```go
func (r *UserResource) Find(ctx fiber.Ctx, currentTenant TenantContext) error
```

或者这种 factory 签名：

```go
func (r *UserResource) CreateHandler(service UserService) func(ctx fiber.Ctx) error
```

那就需要自定义 resolver。

## 注册的 FX group

你需要往这些 group 里注册：

- `group:"vef:api:handler_param_resolvers"`
- `group:"vef:api:factory_param_resolvers"`

## handler resolver 示例

```go
package tenantresolver

import (
  "reflect"

  "github.com/gofiber/fiber/v3"

  "github.com/coldsmirk/vef-framework-go/api"
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

注册时使用 ``fx.Annotate(..., fx.ResultTags(`group:"vef:api:handler_param_resolvers"`))``。

完整一点的注册示例：

```go
fx.Provide(
  fx.Annotate(
    func() api.HandlerParamResolver { return &TenantResolver{} },
    fx.ResultTags(`group:"vef:api:handler_param_resolvers"`),
  ),
)
```

## factory resolver

factory resolver 的思路完全类似，只不过它是在启动期、handler factory 被物化的时候执行一次。

## 扩展前先问自己

在增加自定义 resolver 之前，最好先判断一下下面几种更简单的方式是否够用：

- 把依赖作为 resource 字段注入
- 直接使用已有内置 resolver
- 通过请求结构体传值

只有当某个依赖在大量 handler 中反复出现时，自定义 resolver 才真的值得。

## 下一步

如果你还想先看清楚内置请求解码已经支持什么，再决定要不要扩展，继续阅读 [参数与元信息（Meta）](../guide/params-and-meta)。
