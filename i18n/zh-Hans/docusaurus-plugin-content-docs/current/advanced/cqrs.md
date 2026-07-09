---
sidebar_position: 1
---

# CQRS

VEF 内置了一套轻量级 CQRS bus，用来承载 command、query 和统一行为（behavior）中间件。

## 面向用户的公共接口

应用层最常用的是这些：

- `cqrs.BaseCommand`
- `cqrs.BaseQuery`
- `cqrs.Register(...)`
- `cqrs.Send(...)`
- `cqrs.Behavior`
- `cqrs.BehaviorFunc`
- `vef.ProvideCQRSBehavior(...)`

其他公开 API：

| API | 作用 |
| --- | --- |
| `cqrs.NewBus(behaviors)` | 创建独立 bus，主要用于测试或自定义装配 |
| `cqrs.Bus` | handler registry 和 dispatcher 抽象 |
| `cqrs.Action` / `ActionKind` | action 契约；`Command` 和 `Query` 是公开 kind 常量 |
| `cqrs.Handler[TAction, TResult]` / `HandlerFunc[...]` | 类型化 handler 契约 |
| `cqrs.Behavior` / `BehaviorFunc` | command/query 执行管线 |
| `cqrs.Ordered` | 可选的 behavior 排序 hook |
| `cqrs.Unit` | command 常用的空结果类型 |
| `ErrHandlerNotFound` | action 类型没有注册 handler |
| `ErrResultTypeMismatch` | handler 返回值无法转换成请求的结果类型 |

### Action kind 契约

`Action.Kind()` 返回 action discriminator。`BaseCommand.Kind()` 返回
`Command`（`0`），`BaseQuery.Kind()` 返回 `Query`（`1`）。

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

handler 按 action 类型注册，`Send` 也是类型安全的：

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

`Register` 使用具体的 `TAction` 类型作为 registry key；如果同一个 action type
已经注册过 handler，会 panic。`Send` 按同一个 action type 分发。没有 handler
时返回可用 `errors.Is` 匹配 `ErrHandlerNotFound` 的错误，并返回 zero-value
result。handler 自身返回的 error 会原样向外传播。

`HandlerFunc` 是 `Handler` 的函数适配器；它的 `Handle` 方法只是调用被包装的函数。

## Behavior pipeline

CQRS bus 支持在 command/query 执行前后包裹 middleware 风格的 behavior。

使用 `vef.ProvideCQRSBehavior(...)` 把 behavior 注册进运行时：

```go
vef.ProvideCQRSBehavior(NewLoggingBehavior)
```

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

`BehaviorFunc` 是 `Behavior` 的函数适配器；它的 `Handle` 方法只是调用被包装的函数。
behavior 会收到原始 action，也可以不调用 `next` 而直接 short-circuit。如果返回
`nil`，`Send` 会返回 `TResult` 的 zero value。如果 short-circuit 返回了非 nil
但 concrete type 不是 `TResult` 的值，`Send` 会返回可用 `errors.Is` 匹配
`ErrResultTypeMismatch` 的错误，而不是 panic。

bus 创建时会对 behavior 排序一次。`Ordered.Order()` 控制包裹顺序：较小的值包在
较大的值外面。不实现 `Ordered` 的 behavior 默认 order 是 `0`；独立调用
`NewBus` 时相同 order 会保留输入顺序，但 FX value group 的输入顺序不稳定。框架
约定的 order band 是：

| Order band | 用途 |
| --- | --- |
| `0..99` | transactional / contextual setup，需要包住所有逻辑 |
| `100..199` | audit / collector lifecycle |
| `200..299` | event publish / outbox side effects |
| `1000+` | custom host behaviors |

## 什么时候用它

CQRS 是可选的。当你需要以下能力时会比较合适：

- 明确的 command/query 边界
- 类型安全的分发
- 在 HTTP resource 层之外，围绕应用层 action 建立一条 pipeline

## 它不会帮你做什么

CQRS bus 不会自动扫描 handler。你仍然需要在自己的模块图里明确把 handler 注册进共享 bus。

## 下一步

如果你的 behavior 需要把命令执行包进数据库边界，继续阅读 [事务](../data-access/transactions)。
