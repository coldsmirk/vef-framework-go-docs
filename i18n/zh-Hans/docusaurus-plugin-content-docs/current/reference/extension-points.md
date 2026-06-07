---
sidebar_position: 3
---

# 扩展点

VEF 的扩展核心是 FX group。大多数框架级定制都不是通过修改运行时本身完成，而是通过把组件注册到合适的 group 中完成。

审查说明：本页是 root `vef` package 审查的一部分，覆盖通过 FX group 注册、
通过 `fx.Decorate` 替换默认实现、以及通过 `fx.Supply` supply 单例值的 DI
extension helpers。

## Helper 机制

helper 名称前缀不足以判断实际 wiring 方式，要看具体机制：

| 机制 | Helpers |
| --- | --- |
| `fx.Provide` + `fx.ResultTags` 追加到 group | `ProvideAPIResource`, `ProvideMiddleware`, `ProvideSPAConfig`, `ProvideCQRSBehavior`, `ProvideChallengeProvider`, `ProvideMCPTools`, `ProvideMCPResources`, `ProvideMCPResourceTemplates`, `ProvideMCPPrompts`, `ProvideEventTransport`, `ProvideEventPublishMiddleware`, `ProvideEventConsumeMiddleware`, `ProvideApprovalLifecycleHook`, `ProvideDataSourceProvider` |
| 带 group tag 的 `fx.Supply` | `SupplySPAConfigs` |
| `fx.Decorate` 替换默认实现 | `SupplyFileACL`, `SupplyURLKeyMapper`, `SupplyBusinessBindingHook`, `ProvideEventMetricsRecorder`, `ProvideEventErrorSink` |
| 普通 `fx.Supply` 值 | `SupplyMCPServerInfo` |

替换型 helper 是单服务 override，不是 append-only extension point。除非你明确
希望后面的 FX option 替换前面的实现，否则同一个默认服务只注册一个替代实现。

## API 资源

使用：

```go
vef.ProvideAPIResource(...)
```

对应 FX group：

```text
vef:api:resources
```

## 应用级 middleware

使用：

```go
vef.ProvideMiddleware(...)
```

对应 FX group：

```text
vef:app:middlewares
```

## SPA 配置

使用：

```go
vef.ProvideSPAConfig(...)
vef.SupplySPAConfigs(...)
```

对应 FX group：

```text
vef:spa
```

## CQRS behavior

使用：

```go
vef.ProvideCQRSBehavior(...)
```

对应 FX group：

```text
vef:cqrs:behaviors
```

## 安全 challenge provider

使用：

```go
vef.ProvideChallengeProvider(...)
```

对应 FX group：

```text
vef:security:challenge_providers
```

## 事件总线

使用：

```go
vef.ProvideEventTransport(...)
vef.ProvideEventPublishMiddleware(...)
vef.ProvideEventConsumeMiddleware(...)
```

对应 FX group：

- `vef:event:transports`
- `vef:event:publish-middlewares`
- `vef:event:consume-middlewares`

transport helper 用于注册自定义 `event/transport.Transport`。publish
middleware 在事件 frame 交给 transport 之前运行；consume middleware 包裹
subscriber handler。

另外两个事件扩展点是替换默认实现，而不是追加 group 成员：

- `vef.ProvideEventMetricsRecorder(...)` 替换默认的 `event.MetricsRecorder`
- `vef.ProvideEventErrorSink(...)` 替换异步 publish error sink

## 数据源 provider

使用：

```go
vef.ProvideDataSourceProvider(...)
```

对应 FX group：

```text
vef:datasource:providers
```

`datasource.Provider` 会在启动期加载额外的数据源 spec；此时 primary 和静
态 TOML 数据源已经注册完成。每个返回的 `datasource.Spec` 都会注册进
`datasource.Registry`；如果和 TOML 或另一个 provider 发生名称冲突，启动
会失败。

`github.com/coldsmirk/vef-framework-go/datasource` 已审查公开 surface：

- 17 个 top-level symbols
- 9 个 exported struct fields
- 13 个 exported methods
- fingerprint `a8d1f60b94e7300151d3df0025eec3b3e387d732829ecfff0ecaf7a660ba3cc3`

主数据源固定保留在 `datasource.PrimaryName`（`"primary"`）下。它来自
`vef.data_sources.primary`，作为全框架 `orm.DB` 暴露，不能通过动态
registry API 修改。静态 TOML 非 primary 条目会先于 provider spec 注册。
provider 顺序未定义；`Provider.Load` 返回错误或 registry 名称冲突都会让
应用启动失败，`Provider.Name` 会出现在诊断信息里。

datasource top-level API：

| API | 契约 |
| --- | --- |
| `datasource.ConnectionInfo` | `Registry.TestConnection` 的返回结果；exported field `datasource.ConnectionInfo.Version` 是 `string`。 |
| `datasource.ErrClosed` | registry 开始 shutdown 后，`datasource.Registry.Register`、`datasource.Registry.Update`、`datasource.Registry.Unregister` 返回该错误。 |
| `datasource.ErrExists` | `datasource.Registry.Register` 遇到已注册名称时返回该错误。 |
| `datasource.ErrNameInvalid` | `datasource.Registry.Register` 或 `datasource.Registry.Update` 遇到空名称、包含 whitespace/control characters 的名称时返回该错误。 |
| `datasource.ErrNotFound` | `datasource.Registry.Get`、`datasource.Registry.Kind`、`datasource.Registry.Update`、`datasource.Registry.Unregister` 遇到未注册名称时返回该错误。 |
| `datasource.ErrPrimaryReserved` | `datasource.Registry.Register`、`datasource.Registry.Update`、`datasource.Registry.Unregister` 操作 `datasource.PrimaryName` 时返回该错误。 |
| `datasource.PrimaryName` | 常量 `"primary"`，TOML primary 数据源的保留名称。 |
| `datasource.Provider` | 启动期 provider interface，包含 `datasource.Provider.Name` 和 `datasource.Provider.Load`。 |
| `datasource.ReconcileOption` | `Registry.Reconcile` 的 functional option 类型。 |
| `datasource.ReconcileOptions` | option 状态，exported field 是 `datasource.ReconcileOptions.DryRun`（`bool`）。 |
| `datasource.ReconcileReport` | reconcile 结果，exported fields 是 `datasource.ReconcileReport.Added`、`datasource.ReconcileReport.Updated`、`datasource.ReconcileReport.Removed`、`datasource.ReconcileReport.Errors`。 |
| `datasource.RegisterOption` | `Registry.Update` 和 `Registry.Unregister` 的 functional option 类型。 |
| `datasource.RegisterOptions` | option 状态，exported field 是 `datasource.RegisterOptions.CloseGrace`（`time.Duration`）。 |
| `datasource.Registry` | 可注入的 registry interface，用于 primary lookup、named lookup、mutation、reconcile、probe 和 health check。 |
| `datasource.Spec` | 期望状态或 provider 返回的数据源，exported fields 是 `datasource.Spec.Name` 和 `datasource.Spec.Config`。 |
| `datasource.WithCloseGrace(d)` | 仅在 `d > 0` 时设置 `RegisterOptions.CloseGrace`；延迟关闭异步执行，并会被 shutdown 提前唤醒。 |
| `datasource.WithReconcileDryRun()` | 设置 `ReconcileOptions.DryRun`，让 `Reconcile` 只报告 diff，不打开或关闭连接。 |

`datasource.Registry` methods：

| Method | 契约 |
| --- | --- |
| `datasource.Registry.Primary` | 返回 primary `orm.DB`；等价于 `Get(datasource.PrimaryName)`，但不返回 error。 |
| `datasource.Registry.Get` | 返回已注册 `orm.DB`；未知名称返回 `datasource.ErrNotFound`。 |
| `datasource.Registry.Has` | 判断名称当前是否已注册。 |
| `datasource.Registry.Names` | 返回全部已注册名称，包括 `primary`，顺序为稳定 lexical order。 |
| `datasource.Registry.Kind` | 返回配置的 `config.DBKind`；未知名称返回 `datasource.ErrNotFound`。 |
| `datasource.Registry.Register` | 先打开并 ping 新的非 primary 数据源，再插入 registry。重复名称返回 `datasource.ErrExists`；冲突路径会关闭新连接池。 |
| `datasource.Registry.Update` | 先打开并 ping 替换连接，再交换进去。失败时保留旧 entry；成功后异步关闭旧连接池，并支持 `datasource.WithCloseGrace`。 |
| `datasource.Registry.Unregister` | 移除非 primary 数据源，然后异步关闭旧连接池，并支持 `datasource.WithCloseGrace`。 |
| `datasource.Registry.Reconcile` | 串行化 reconcile 调用，忽略空名称和 `primary` specs，对 add/update/remove bucket 排序；局部失败写入 `ReconcileReport.Errors`，不会因为 partial failure 返回 top-level error。 |
| `datasource.Registry.TestConnection` | 打开临时连接、查询 server version、关闭连接、返回 `datasource.ConnectionInfo`，且绝不修改 registry。probe 有内部 5s timeout 上限，同时仍尊重调用方 cancellation/deadline。 |
| `datasource.Registry.HealthCheck` | 并行 ping primary 和全部非 primary entries，返回 name-to-error map；nil error 表示该数据源可达。 |

## MCP provider

使用：

```go
vef.ProvideMCPTools(...)
vef.ProvideMCPResources(...)
vef.ProvideMCPResourceTemplates(...)
vef.ProvideMCPPrompts(...)
vef.SupplyMCPServerInfo(...)
```

对应 FX group：

- `vef:mcp:tools`
- `vef:mcp:resources`
- `vef:mcp:templates`
- `vef:mcp:prompts`

## 存储集成

Extension groups:

- `vef:api:resources`
- `vef:app:middlewares`

使用：

```go
vef.SupplyFileACL(...)
vef.SupplyURLKeyMapper(...)
```

`SupplyFileACL` 替换默认 `storage.FileACL`。默认 ACL 只允许读取 `pub/`
下的 key；如果应用存储私有文件，应提供业务自己的 ACL。

`SupplyURLKeyMapper` 替换默认的 `storage.ProxyURLKeyMapper`，它会把
`/storage/files/<key>` 代理 URL 映射回对象 key。当富文本或 markdown 中
嵌入的是 CDN URL 或其他需要映射回对象 key 的 URL 形式时，使用它提供映射
规则。只有内容直接嵌入 bare object key 时才使用
`storage.IdentityURLKeyMapper`。这个默认值指的是框架 DI graph；如果直接调用
`storage.NewFiles(...)`，nil mapper 会被规整为 `storage.IdentityURLKeyMapper`。

## 审批生命周期 hooks

使用：

```go
vef.ProvideApprovalLifecycleHook(...)
```

对应 FX group：

```text
vef:approval:lifecycle_hooks
```

`approval.InstanceLifecycleHook` 会在审批 engine transaction 内同步运行，
例如实例创建和实例完成等生命周期点。hook 返回 error 会回滚当前审批命
令。提交后的异步集成应使用事件订阅。

## 审批业务绑定

使用：

```go
vef.SupplyBusinessBindingHook(...)
```

这会替换默认 `approval.BusinessBindingHook`，用于
`Flow.BindingMode == BindingBusiness` 的场景。实现负责把审批实例和宿主
业务表桥接起来；异步状态回写路径必须保持幂等。

## 日志

使用：

```go
vef.NamedLogger(name)
```

当集成代码需要在 DI 之外拿框架 `logx.Logger` 时，可以使用
`vef.NamedLogger(name)` 这个 root-package 便捷函数。它返回
`logx.Logger`；`logx` package 本身只公开 `Level` constants、
`Level.String()` 和 `Logger` interface contract。

## API 参数注入解析器

这是更进阶的扩展点：

- `vef:api:handler_param_resolvers`
- `vef:api:factory_param_resolvers`

当内置 handler 参数集合不够时，可以往这里注册自定义 resolver。

## 一个简单判断原则

当你想扩展 VEF 时，优先判断：

- 框架是否已经为这个概念预留了 group
- 你的扩展是否应该纳入启动和生命周期管理
- 你是否希望模块之间依赖保持显式和可测试

如果答案是肯定的，那通常就应该走 FX group，而不是用隐式全局状态或手写单例。

## 延伸阅读

- [模块与依赖注入](../modules/overview)：这些 group 如何进入应用装配流程
- [自定义参数解析器](../advanced/custom-param-resolvers)：handler 注入扩展的具体做法
