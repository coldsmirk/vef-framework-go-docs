---
sidebar_position: 1
---

# 模块与依赖注入

VEF 基于 Uber FX 构建。公开的 `vef` 包把最常用的 FX helper 重新导出了，所以大多数应用都可以在同一套 API 表面内完成组合。

## 核心思路

你不需要手动启动每个子系统，而是通过 FX option 组合它们：

```go
vef.Run(
  user.Module,
  auth.Module,
  vef.ProvideAPIResource(resources.NewHealthResource),
)
```

在内部，`vef.Run(...)` 会把你的 option 追加到框架自己的模块列表里，再启动整个 FX 应用。

在真实项目中，`main.go` 常常更像这样：

```go
vef.Run(
  ivef.Module,  // 框架侧集成
  mcp.Module,   // 自定义 MCP provider
  web.Module,   // SPA 托管
  auth.Module,  // 认证加载器
  sys.Module,   // 系统/管理资源
  md.Module,    // 主数据资源
  pmr.Module,   // 业务资源
)
```

这反映出一个重要模式：VEF 应用通常是由多块小模块组合出来的，而不是一个超大的总模块。

## `vef` 重新导出的常用 helper

业务侧最常接触的是这些：

- `vef.Run`
- `vef.Module`
- `vef.Provide`
- `vef.Supply`
- `vef.Invoke`
- `vef.Decorate`
- `vef.Replace`

这样大多数应用代码不需要频繁直接 import `fx`。

## 基于 group 的扩展点

很多框架能力都是通过 FX group 串起来的。对应用开发者最重要的是：

- `vef:api:resources`
- `vef:app:middlewares`
- `vef:cqrs:behaviors`
- `vef:security:challenge_providers`
- `vef:mcp:tools`
- `vef:mcp:resources`
- `vef:mcp:templates`
- `vef:mcp:prompts`

`di.go` 里的各种 helper，本质上就是帮你更安全地把值注册进这些 group。

## API 资源注册

最常用的 helper 是：

```go
vef.ProvideAPIResource(NewUserResource)
```

它会把构造函数结果打上 API resource group 标签。启动时，API 模块会从这个 group 收集所有资源，并把其中的操作注册进引擎。

## 其他 provider 也是同一个模式

中间件、CQRS、登录挑战、MCP 也都一样：

```go
vef.ProvideMiddleware(NewAuditTrailMiddleware)
vef.ProvideCQRSBehavior(NewTracingBehavior)
vef.ProvideChallengeProvider(NewTOTPChallengeProvider)
vef.ProvideMCPTools(NewToolProvider)
```

这些 helper 的价值不在“增加新能力”，而在于把 FX group tag 隐藏掉，让应用代码更易读。

## 大型应用里的模块角色

在更接近生产的 VEF 项目里，通常会出现这些固定角色：

- `internal/vef`：build info、共享框架侧 service、事件订阅器
- `internal/auth`：`UserLoader`、`UserInfoLoader`、认证相关初始化
- 若干业务域模块：注册 API 资源
- 可选的 `web` 和 `mcp` 模块：分别负责 SPA 与 MCP 集成

这样职责会比“一个 app 模块包所有东西”清楚得多。

## 用 `vef.Invoke(...)` 做集成型模块

在规模更大的应用里，经常会有一个专门的集成模块，通过 `vef.Invoke(...)` 做启动期 wiring，而不是直接暴露业务资源。

例如：

```go
var Module = vef.Module(
  "app:vef",
  vef.Supply(BuildInfo),
  vef.Provide(NewDataDictLoader, password.NewBcryptEncoder),
  vef.Invoke(registerEventSubscribers),
)
```

这种模块很适合放 build info、共享框架侧 service，以及事件订阅器注册。

## 为什么资源能自动发现

API engine 不要求你手动挂路由。它会按下面的流程工作：

1. 从容器里收集资源
2. 从资源本身和嵌入的 CRUD provider 里收集操作
3. 解析 handler
4. 把 handler 适配成 Fiber handler
5. 挂载到 RPC 或 REST router

所以 VEF 的典型开发方式更像“定义资源 + 注册构造函数”，而不是“声明 router + 绑定 handler + 手挂中间件”。

## 什么时候直接用 `fx`

大多数应用都可以只用 `vef` 包装层，但如果你需要这些高级能力，直接用 `fx` 也完全没问题：

- 更复杂的注解
- 可选依赖
- 直接操作 lifecycle hook
- 测试环境替换实现

VEF 并没有限制你使用 FX，它只是把最常见路径做短了。

## 下一步

继续看 [应用生命周期](./lifecycle)，理解 `vef.Run(...)` 从启动到监听端口到底做了什么。
