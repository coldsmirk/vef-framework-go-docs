# CQRS

VEF 内置了一套轻量级 CQRS bus，用来承载 command、query 和统一行为（behavior）中间件。

## 面向用户的公共接口

应用层最常用的是这些：

- `cqrs.BaseCommand`
- `cqrs.BaseQuery`
- `cqrs.Register(...)`
- `cqrs.Send(...)`
- `vef.ProvideCQRSBehavior(...)`

## 定义 action

command 嵌入 `cqrs.BaseCommand`，query 嵌入 `cqrs.BaseQuery`：

```go
type CreateUser struct {
  cqrs.BaseCommand
  Name string
}

type GetUser struct {
  cqrs.BaseQuery
  ID string
}
```

## 注册 handler

```go
package useractions

import (
  "context"

  "github.com/coldsmirk/vef-framework-go/cqrs"
)

type CreateUser struct {
  cqrs.BaseCommand
  Name string
}

type CreateUserHandler struct{}

func (CreateUserHandler) Handle(ctx context.Context, cmd CreateUser) (cqrs.Unit, error) {
  return cqrs.Unit{}, nil
}

func RegisterHandlers(bus cqrs.Bus) {
  cqrs.Register(bus, CreateUserHandler{})
}

func Run(ctx context.Context, bus cqrs.Bus) error {
  _, err := cqrs.Send[CreateUser, cqrs.Unit](ctx, bus, CreateUser{Name: "alice"})
  return err
}
```

## Behavior

Behavior 会像 middleware 一样包裹所有 command/query。注册方式是：

```go
vef.ProvideCQRSBehavior(NewLoggingBehavior)
```

内部 bus 会倒序包裹 behavior，因此最先注册的 behavior 会成为最外层。

比较适合用 behavior 做的事情：

- 日志
- tracing
- metrics
- 统一校验
- 事务包装

最小 behavior 示例：

```go
func NewLoggingBehavior() cqrs.Behavior {
  return cqrs.BehaviorFunc(func(ctx context.Context, action cqrs.Action, next func(context.Context) (any, error)) (any, error) {
    return next(ctx)
  })
}
```

## 它不会帮你做什么

CQRS bus 不会自动扫描 handler。你仍然需要在自己的模块图里明确把 handler 注册进共享 bus。

## 下一步

如果你的 behavior 需要把命令执行包进数据库边界，继续阅读 [事务](./transactions)。
