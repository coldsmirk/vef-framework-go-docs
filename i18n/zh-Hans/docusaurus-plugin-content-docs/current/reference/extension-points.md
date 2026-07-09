---
sidebar_position: 3
---

# 扩展点

VEF 的大多数扩展点都是显式的 FX group。

本页覆盖通过 FX group 注册、通过 `fx.Decorate` 替换默认实现、以及通过 `fx.Supply` 提供单例值的 DI extension helpers。

## Helper 机制

helper 名称前缀不足以判断实际 wiring 方式，要看具体机制：

| 机制 | Helpers |
| --- | --- |
| `fx.Provide` + `fx.ResultTags` 追加到 group | `ProvideAPIResource`, `ProvideAuthStrategy`, `ProvideMiddleware`, `ProvideSPAConfig`, `ProvideCQRSBehavior`, `ProvideChallengeProvider`, `ProvideMCPTools`, `ProvideMCPResources`, `ProvideMCPResourceTemplates`, `ProvideMCPPrompts`, `ProvideEventTransport`, `ProvideEventPublishMiddleware`, `ProvideEventConsumeMiddleware`, `ProvideApprovalLifecycleHook`, `ProvideApprovalAggregator`, `ProvideDataSourceProvider` |
| 带 group tag 的 `fx.Supply` | `SupplySPAConfigs` |
| `fx.Decorate` 替换默认实现 | `SupplyFileACL`, `SupplyURLKeyMapper`, `SupplyBusinessRefProvider`, `SupplyBusinessRefResolver`, `ProvideEventMetricsRecorder`, `ProvideEventErrorSink`, `ProvideApprovalFormSchemaParser` |
| 普通 `fx.Supply` 值 | `SupplyMCPServerInfo` |

替换型 helper 是单服务 override，不是 append-only extension point。除非你明确希望后面的 FX option 替换前面的实现，否则同一个默认服务只注册一个替代实现。

## API 与应用级 group

- `vef:api:resources`
- `vef:api:auth_strategies`
- `vef:app:middlewares`

Helpers：

- `vef.ProvideAPIResource(...)`
- `vef.ProvideAuthStrategy(...)`
- `vef.ProvideMiddleware(...)`

`ProvideAuthStrategy` 会把自定义 `api.AuthStrategy` 追加到认证策略 group。资源或操作通过 `api.AuthConfig.Strategy` 选择其 `Name()` 返回的策略名；内置策略是 `none`、`bearer`、`signature` 和 `ip`。

## 最小模块示例

```go
var Module = vef.Module(
  "app:user",
  vef.ProvideAPIResource(NewUserResource),
  vef.ProvideMiddleware(NewAuditMiddleware),
)
```

## API 参数注入

- `vef:api:handler_param_resolvers`
- `vef:api:factory_param_resolvers`

这两个 group 分别扩展请求期和启动期的 handler 注入。

## CQRS

- `vef:cqrs:behaviors`

Helper：

- `vef.ProvideCQRSBehavior(...)`

## 安全

- `vef:security:challenge_providers`

Helper：

- `vef.ProvideChallengeProvider(...)`

## 事件总线

- `vef:event:transports`
- `vef:event:publish-middlewares`
- `vef:event:consume-middlewares`

Helpers：

- `vef.ProvideEventTransport(...)`
- `vef.ProvideEventPublishMiddleware(...)`
- `vef.ProvideEventConsumeMiddleware(...)`

transport helper 用于注册自定义 `event/transport.Transport` 实现。publish middleware 在事件 frame 交给 transport 之前运行；consume middleware 包裹在 subscriber handler 周围。

另外两个事件集成点是替换框架默认实现，而不是追加 group 成员：

- `vef.ProvideEventMetricsRecorder(...)` 替换默认的 `event.MetricsRecorder`。
- `vef.ProvideEventErrorSink(...)` 替换异步 publish 的 error sink。

## 数据源 provider

- `vef:datasource:providers`

Helper：

- `vef.ProvideDataSourceProvider(...)`

`datasource.Provider` 会在启动期加载额外的数据源 spec，此时 primary 和静态 TOML 数据源已经注册完成。每个返回的 `datasource.Spec` 都会注册进 `datasource.Registry`；如果和 TOML 或另一个 provider 发生名称冲突，启动会失败。

主数据源固定保留在 `datasource.PrimaryName`（`"primary"`）下。它来自 `vef.data_sources.primary`，作为全框架 `orm.DB` 暴露，不能通过动态 registry API 修改。静态 TOML 非 primary 条目会先于 provider spec 注册。provider 顺序未定义；`Provider.Load` 返回错误或 registry 名称冲突都会让应用启动失败，`Provider.Name` 会出现在诊断信息里。

datasource 顶层 API：

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
| `datasource.WithCloseGrace(d)` | 仅在 `d > 0` 时设置 `RegisterOptions.CloseGrace`；延迟关闭异步执行，并会被 shutdown 提前截断。 |
| `datasource.WithReconcileDryRun()` | 设置 `ReconcileOptions.DryRun`，让 `Reconcile` 只报告 diff，不打开或关闭连接。 |

`datasource.Registry` 的方法：

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

## SPA

- `vef:spa`

Helpers：

- `vef.ProvideSPAConfig(...)`
- `vef.SupplySPAConfigs(...)`

## 存储集成

扩展 group：

- `vef:api:resources`
- `vef:app:middlewares`

Helpers：

- `vef.SupplyFileACL(...)`
- `vef.SupplyURLKeyMapper(...)`

`SupplyFileACL` 替换默认的 `storage.FileACL`。默认实现只允许读取 `pub/` 下的 key；如果应用要存储私有文件，应提供业务自己的 ACL。

`SupplyURLKeyMapper` 替换默认的 `storage.ProxyURLKeyMapper`，它会把 `/storage/files/<key>` 代理 URL 映射回对象 key。当富文本或 markdown 内容嵌入的是 CDN URL 或其他需要映射回对象 key 的 URL 形式时，覆盖它。只有内容直接嵌入裸 object key 时才使用 `storage.IdentityURLKeyMapper`。这个默认值针对的是框架 DI graph；直接调用 `storage.NewFiles(...)` 时，nil mapper 会被规整为 `storage.IdentityURLKeyMapper`。

## MCP

- `vef:mcp:tools`
- `vef:mcp:resources`
- `vef:mcp:templates`
- `vef:mcp:prompts`

Helpers：

- `vef.ProvideMCPTools(...)`
- `vef.ProvideMCPResources(...)`
- `vef.ProvideMCPResourceTemplates(...)`
- `vef.ProvideMCPPrompts(...)`
- `vef.SupplyMCPServerInfo(...)`

## 审批生命周期 hook

- `vef:approval:lifecycle_hooks`

Helper：

- `vef.ProvideApprovalLifecycleHook(...)`

`approval.InstanceLifecycleHook` 的实现会在审批 engine 事务内同步运行，覆盖实例创建、完成等生命周期节点。返回 error 会回滚外层的审批命令。需要在提交后才运行的异步集成应使用事件订阅。

## 审批 aggregator 与表单 schema

- `vef:approval:aggregators`

Helpers：

- `vef.ProvideApprovalAggregator(...)`
- `vef.ProvideApprovalFormSchemaParser(...)`

`ProvideApprovalAggregator` 为审批字段条件注册自定义的 detail-table aggregator，与内置的 sum / count / avg 并存。constructor 必须返回 `approval.Aggregator`；条件求值器按其 `AggregateKind` 选用。如果内置聚合类型未注册对应实现，启动会失败。

`ProvideApprovalFormSchemaParser` 替换框架默认的 `approval.FormSchemaParser`（内置实现是 vef-framework-react 表单编辑器解析器）。这个替换是整体覆盖，不是追加：每一次部署的表单 schema 都会经过它，因此它必须理解宿主提交的每一种设计器文档。解析只在流程部署时运行一次；更早部署的版本仍保留部署时持久化的 `form_fields`。

## 审批业务绑定

Helpers：

- `vef.SupplyBusinessRefProvider(...)`
- `vef.SupplyBusinessRefResolver(...)`

`SupplyBusinessRefProvider` 替换默认的 no-op `approval.BusinessRefProvider`。它在 `Flow.BindingMode == BindingBusiness` 的 `start_instance` 事务内运行，让宿主解析或创建业务行，并返回 opaque 的 `Instance.BusinessRef`。

`SupplyBusinessRefResolver` 替换默认的 identity `approval.BusinessRefResolver`。当 `Instance.BusinessRef` 不是裸主键、而引擎写回需要提取出与 `Flow.BusinessPKField` 匹配的值时，注册一个自定义实现。

审批状态写回本身由 engine 拥有。宿主可以用 `approval.InstanceLifecycleHook` 或事件订阅在其周围扩展，但不再替换写回路径本身。

## 日志

Helper：

- `vef.NamedLogger(name)`

当集成代码需要在依赖注入之外获取框架的 `logx.Logger` 时，使用这个 root-package 便捷函数。它返回 `logx.Logger`；`logx` package 本身公开 `Level` 常量、`Level.String()`、`Logger` interface 契约，以及 `LoggerConfigurable[T]`（供 immutable component 通过 `WithLogger` 返回一个配置了 logger 的副本）。

## 一个简单的判断原则

决定如何扩展 VEF 时，先问自己：

- 框架是否已经为这个概念预留了 group？
- 这个扩展是否应该纳入启动和生命周期管理？
- 你是否希望模块之间的依赖保持显式、可测试？

如果答案是肯定的，通常应该走 FX group，而不是依赖隐式全局状态或手写单例。

## 延伸阅读

- [模块与依赖注入](../core-concepts/overview)：这些 group 如何融入应用装配流程
- [扩展 Handler 参数](../advanced/extending-parameters)：handler 注入扩展的具体做法
