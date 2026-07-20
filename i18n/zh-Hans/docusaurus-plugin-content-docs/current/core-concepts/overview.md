---
sidebar_position: 1
---

# 模块与依赖注入

VEF 基于 Uber FX 构建。公开的 `vef` 包把最常用的 FX helper 重新导出了，所以大多数应用都可以在同一套 API 表面内完成组合。本页、[扩展点](../reference/extension-points) 和 [应用生命周期](./lifecycle) 三页合起来覆盖了 root `vef` package 的公开表面。

## 核心思路

你不需要手动启动每个子系统。你的代码通过向 FX group 提供构造函数来加入正在运行的应用，框架会自动发现并把它们接入依赖图。API 资源、中间件、CQRS behavior 背后都是同一套机制——具体的资源写法见 [API 资源](../building-apis/api)；本页讲的是这套机制底层的 DI 原理。

```go
vef.Run(
  user.Module,
  auth.Module,
)
```

在内部，`vef.Run(...)` 会把你的 option 追加到框架自己的模块列表里，再启动整个 FX 应用。

在真实项目中，`main.go` 常常更像这样：

```go
vef.Run(
  ivef.Module,  // 框架侧集成
  tools.Module, // 通过 vef.ProvideMCP* 注册自定义 MCP provider
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
- `vef.Annotate`
- `vef.As`
- `vef.From`
- `vef.ParamTags`
- `vef.ResultTags`
- `vef.Self`
- `vef.Invoke`
- `vef.Decorate`
- `vef.Replace`
- `vef.Populate`
- `vef.Private`
- `vef.OnStart`
- `vef.OnStop`

它也重新导出了常用 FX 标记类型：

- `vef.In`
- `vef.Out`
- `vef.Lifecycle`
- `vef.Hook`
- `vef.HookFunc`

生命周期 hook 包装函数包括 `vef.StartHook`、`vef.StopHook` 和
`vef.StartStopHook`。

包装层也公开了 `vef.From`、`vef.Replace` 和 `vef.Populate`，用于更高级的 DI
场景。它们就是同名 FX primitive，只是放在 `vef` 包下，方便框架侧模块通常不
必直接 import `go.uber.org/fx`。

这样大多数应用代码不需要频繁直接 import `fx`。

## 基于 group 的扩展点

很多框架能力都是通过 FX group 串起来的。对应用开发者最重要的是：

- `vef:api:resources`
- `vef:api:auth_strategies`
- `vef:app:middlewares`
- `vef:cqrs:behaviors`
- `vef:security:challenge_providers`
- `vef:mcp:tools`
- `vef:mcp:resources`
- `vef:mcp:templates`
- `vef:mcp:prompts`
- `vef:event:transports`
- `vef:event:publish-middlewares`
- `vef:event:consume-middlewares`
- `vef:datasource:providers`
- `vef:approval:lifecycle_hooks`
- `vef:cron:job_handlers`（v0.39）
- `vef:js:libs`（v0.39）
- `vef:integration:inbound_handlers` / `vef:integration:outbound_auth_schemes` / `vef:integration:inbound_auth_schemes`（v0.39）
- `vef:security:session_revocation_listeners`（v0.39）

`di.go` 里的各种 helper，本质上就是帮你更安全地把值注册进这些 group。helper
名称前缀不一定等于 FX 机制——有的 `Provide*` 是向 group 追加，有的则是替换默认的单个服务。逐个 helper 的权威表格（机制、group、契约）见[扩展点](../reference/extension-points)，本页不再重复维护。

## API 资源注册

最常用的 helper 是：

```go
vef.ProvideAPIResource(NewUserResource)
```

它会把构造函数结果打上 API resource group 标签。启动时，API 模块会从这个 group 收集所有资源，并把其中的操作注册进引擎。资源本身怎么写见 [API 资源](../building-apis/api)。

## 其他 provider 也是同一个模式

中间件、CQRS、登录挑战、MCP 也都一样：

```go
vef.ProvideMiddleware(NewAuditTrailMiddleware)
vef.ProvideAuthStrategy(NewAPIKeyStrategy)
vef.ProvideCQRSBehavior(NewTracingBehavior)
vef.ProvideChallengeProvider(NewTOTPChallengeProvider)
vef.ProvideMCPTools(NewToolProvider)
vef.ProvideMCPResources(NewResourceProvider)
vef.ProvideMCPResourceTemplates(NewTemplateProvider)
vef.ProvideMCPPrompts(NewPromptProvider)
vef.ProvideSPAConfig(NewWebConfig)
vef.ProvideEventTransport(NewKafkaTransport)
vef.ProvideEventPublishMiddleware(NewAuditPublishMiddleware)
vef.ProvideEventConsumeMiddleware(NewRecoverConsumeMiddleware)
vef.ProvideDataSourceProvider(NewTenantDataSourceProvider)
```

这些 helper 的价值不在“增加新能力”，而在于把 FX group tag 隐藏掉，让应用代码更易读。

另一些 helper 是替换框架默认实现，而不是追加 group 成员：

- `vef.ProvideEventMetricsRecorder(...)`
- `vef.ProvideEventErrorSink(...)`
- `vef.SupplyFileACL(...)`
- `vef.SupplyURLKeyMapper(...)`
- `vef.SupplyBusinessRefProvider(...)`
- `vef.SupplyBusinessRefResolver(...)`

`vef.SupplyMCPServerInfo(...)` 不同：它 supply 单个 `mcp.ServerInfo` 值。
`vef.SupplySPAConfigs(...)` 也不同：它把一个或多个 `middleware.SPAConfig` 值
supply 到 `vef:spa` group。

当集成代码需要在 DI 之外拿框架日志接口时，可以使用
`vef.NamedLogger(name)` 创建 `logx.Logger`。

## 可选功能模块

有些框架功能不在默认 boot graph 中。只有应用需要时才显式启用：

```go
vef.Run(
  vef.ApprovalModule,
  user.Module,
)
```

`vef.ApprovalModule` 会开启审批/工作流功能，并注册它的 API resources、
CQRS handlers、engine、业务投影 worker 和 scanners。审批的
`approval.*` 事件需要 transactional route；路由细节
见[审批模块](../approval)。

`vef.IntegrationModule`（v0.39）开启集成引擎——契约、系统、适配器、路由、
入站 HTTP 网关与 `integration/*` 管理资源；见
[集成引擎](../integration/overview)。

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
