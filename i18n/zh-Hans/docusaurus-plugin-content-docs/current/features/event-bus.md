---
sidebar_position: 2
---

# 事件总线

VEF 会自动启动一个内存事件总线，并通过公共 `event` 包对外暴露。

## 模块自动提供的内容

事件模块会注册一个内存总线，并以这些接口形式暴露：

| 接口 |
| --- |
| `event.Bus` |
| `event.Publisher` |
| `event.Subscriber` |

这个总线会通过 FX 生命周期自动启动和停止。

## 核心事件接口

### `event.Event`

自定义事件应实现以下方法：

| 方法 | 含义 |
| --- | --- |
| `ID()` | 唯一事件实例 ID |
| `Type()` | 事件类型字符串 |
| `Source()` | 事件来源 |
| `Time()` | 发生时间 |
| `Meta()` | 元数据 map |

### 发布与订阅接口

| 接口 | 方法 |
| --- | --- |
| `event.Publisher` | `Publish(event)` |
| `event.Subscriber` | `Subscribe(eventType, handler)` |
| `event.Bus` | 组合了 `Publisher`、`Subscriber`、`Start()`、`Shutdown(ctx)` |

### 中间件接口

| 接口或类型 | 作用 |
| --- | --- |
| `event.Middleware` | 拦截事件投递 |
| `event.MiddlewareFunc` | 中间件链的下一跳函数 |

## `event.BaseEvent`

绝大多数自定义事件都会嵌入 `event.BaseEvent`：

```go
type UserCreatedEvent struct {
  event.BaseEvent

  UserID string `json:"userId"`
}
```

创建 base 部分的方式：

```go
&UserCreatedEvent{
  BaseEvent: event.NewBaseEvent(
    "user.created",
    event.WithSource("user-service"),
    event.WithMeta("scope", "admin"),
  ),
  UserID: "user-1001",
}
```

BaseEvent 相关 helper：

| Helper | 含义 |
| --- | --- |
| `event.NewBaseEvent(type, opts...)` | 创建基础事件 |
| `event.WithSource(source)` | 设置 source |
| `event.WithMeta(key, value)` | 添加 metadata |

## 发布与订阅示例

```go
package userevents

import (
  "context"

  "github.com/coldsmirk/vef-framework-go/event"
)

func PublishUserCreated(publisher event.Publisher, userID string) {
  publisher.Publish(&UserCreatedEvent{
    BaseEvent: event.NewBaseEvent("user.created"),
    UserID:    userID,
  })
}

func RegisterUserCreatedHandler(subscriber event.Subscriber) event.UnsubscribeFunc {
  return subscriber.Subscribe("user.created", func(ctx context.Context, evt event.Event) {
    _ = evt
  })
}
```

## 事件中间件

事件总线支持通过 FX group 挂载事件中间件：

```text
vef:event:middlewares
```

适合放在这里的横切能力包括：

- 事件日志
- tracing
- 过滤
- 事件投递前的轻量变换

## 框架内置事件类型

框架当前会发布以下核心事件类型：

| 事件类型 | 来源 |
| --- | --- |
| `vef.api.request.audit` | API 审计事件 |
| `vef.security.login` | 登录流程事件 |
| `vef.storage.file.promoted` | storage promoter 事件 |
| `vef.storage.file.deleted` | storage promoter 事件 |
| `vef.security.role_permissions.changed` | 角色权限缓存失效事件 |

如果启用了 approval 模块，它还会额外发布一批审批领域事件，但那一层属于业务域事件，而不是这里的框架核心事件。

## 常见接线模式

在规模稍大的应用里，订阅者通常不是在资源里注册，而是放在一个集成模块里：

```go
var Module = vef.Module(
  "app:event",
  vef.Invoke(registerEventSubscribers),
)
```

这个 `registerEventSubscribers` 函数可以订阅框架事件，并在需要时自行挂生命周期清理逻辑。

## 什么时候使用它

内置事件总线适合以下场景：

- 生产者和消费者都在同一个应用进程内
- 只需要异步解耦，不需要外部消息代理
- 当前还不需要跨进程消息系统

即使未来需要跨进程消息系统，也仍然可以保留这个内置总线作为应用层抽象。

## 下一步

继续阅读 [缓存](./cache)，如果你想把事件发布和缓存失效、异步刷新流程串起来，它会直接用到。
