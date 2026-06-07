---
sidebar_position: 3
---

# Extension Points

Most VEF extension points are explicit FX groups.

Audit note: this page is part of the root `vef` package audit. It covers the DI extension helpers that register into FX groups, replace defaults through `fx.Decorate`, or supply singleton values through `fx.Supply`.

## Helper mechanisms

The helper name prefix is not enough to tell how the value is wired. Use the mechanism:

| Mechanism | Helpers |
| --- | --- |
| `fx.Provide` + `fx.ResultTags` group append | `ProvideAPIResource`, `ProvideMiddleware`, `ProvideSPAConfig`, `ProvideCQRSBehavior`, `ProvideChallengeProvider`, `ProvideMCPTools`, `ProvideMCPResources`, `ProvideMCPResourceTemplates`, `ProvideMCPPrompts`, `ProvideEventTransport`, `ProvideEventPublishMiddleware`, `ProvideEventConsumeMiddleware`, `ProvideApprovalLifecycleHook`, `ProvideDataSourceProvider` |
| `fx.Supply` with group tags | `SupplySPAConfigs` |
| `fx.Decorate` replacement | `SupplyFileACL`, `SupplyURLKeyMapper`, `SupplyBusinessBindingHook`, `ProvideEventMetricsRecorder`, `ProvideEventErrorSink` |
| plain `fx.Supply` value | `SupplyMCPServerInfo` |

Replacement helpers are single-service overrides, not append-only extension
points. Register only one implementation unless you intentionally want a later
FX option to replace an earlier one.

## API and app groups

- `vef:api:resources`
- `vef:app:middlewares`

Helpers:

- `vef.ProvideAPIResource(...)`
- `vef.ProvideMiddleware(...)`

## Minimal module example

```go
var Module = vef.Module(
  "app:user",
  vef.ProvideAPIResource(NewUserResource),
  vef.ProvideMiddleware(NewAuditMiddleware),
)
```

## API parameter injection

- `vef:api:handler_param_resolvers`
- `vef:api:factory_param_resolvers`

These extend request-time and startup-time handler injection.

## CQRS

- `vef:cqrs:behaviors`

Helper:

- `vef.ProvideCQRSBehavior(...)`

## Security

- `vef:security:challenge_providers`

Helper:

- `vef.ProvideChallengeProvider(...)`

## Event bus

- `vef:event:transports`
- `vef:event:publish-middlewares`
- `vef:event:consume-middlewares`

Helpers:

- `vef.ProvideEventTransport(...)`
- `vef.ProvideEventPublishMiddleware(...)`
- `vef.ProvideEventConsumeMiddleware(...)`

The transport helper registers custom `event/transport.Transport` implementations.
Publish middleware runs before a frame is handed to a transport; consume
middleware runs around subscriber handlers.

Two event integrations replace framework defaults instead of appending group
members:

- `vef.ProvideEventMetricsRecorder(...)` decorates the default
  `event.MetricsRecorder`.
- `vef.ProvideEventErrorSink(...)` decorates the async publish error sink.

## Data source providers

- `vef:datasource:providers`

Helper:

- `vef.ProvideDataSourceProvider(...)`

A `datasource.Provider` loads additional data source specs during startup, after
the primary and static TOML sources are already registered. Every returned
`datasource.Spec` is registered into the `datasource.Registry`; a name collision
with TOML or another provider fails boot.

Reviewed public surface for `github.com/coldsmirk/vef-framework-go/datasource`:

- 17 top-level symbols
- 9 exported struct fields
- 13 exported methods
- fingerprint `a8d1f60b94e7300151d3df0025eec3b3e387d732829ecfff0ecaf7a660ba3cc3`

The primary source is reserved under `datasource.PrimaryName` (`"primary"`).
It comes from `vef.data_sources.primary`, is exposed as the framework-wide
`orm.DB`, and cannot be mutated through the dynamic registry API. Static
non-primary TOML entries are seeded before provider specs. Provider order is undefined;
`Provider.Load` errors and any registry conflict abort application
startup, and `Provider.Name` is used in diagnostic messages.

Top-level datasource APIs:

| API | Contract |
| --- | --- |
| `datasource.ConnectionInfo` | Result returned by `Registry.TestConnection`; exported field `datasource.ConnectionInfo.Version` is `string`. |
| `datasource.ErrClosed` | Returned by `datasource.Registry.Register`, `datasource.Registry.Update`, or `datasource.Registry.Unregister` after registry shutdown begins. |
| `datasource.ErrExists` | Returned by `datasource.Registry.Register` when the name is already registered. |
| `datasource.ErrNameInvalid` | Returned by `datasource.Registry.Register` or `datasource.Registry.Update` when the name is empty or contains whitespace/control characters. |
| `datasource.ErrNotFound` | Returned by `datasource.Registry.Get`, `datasource.Registry.Kind`, `datasource.Registry.Update`, or `datasource.Registry.Unregister` when the name is not registered. |
| `datasource.ErrPrimaryReserved` | Returned by `datasource.Registry.Register`, `datasource.Registry.Update`, or `datasource.Registry.Unregister` for `datasource.PrimaryName`. |
| `datasource.PrimaryName` | Constant `"primary"`, the reserved TOML primary source name. |
| `datasource.Provider` | Startup provider interface with `datasource.Provider.Name` and `datasource.Provider.Load`. |
| `datasource.ReconcileOption` | Functional option type for `Registry.Reconcile`. |
| `datasource.ReconcileOptions` | Option state with exported field `datasource.ReconcileOptions.DryRun` (`bool`). |
| `datasource.ReconcileReport` | Reconcile result with exported fields `datasource.ReconcileReport.Added`, `datasource.ReconcileReport.Updated`, `datasource.ReconcileReport.Removed`, and `datasource.ReconcileReport.Errors`. |
| `datasource.RegisterOption` | Functional option type for `Registry.Update` and `Registry.Unregister`. |
| `datasource.RegisterOptions` | Option state with exported field `datasource.RegisterOptions.CloseGrace` (`time.Duration`). |
| `datasource.Registry` | Injectable registry interface for primary lookup, named lookup, mutation, reconcile, probing, and health checks. |
| `datasource.Spec` | Desired or provider-supplied source with exported fields `datasource.Spec.Name` and `datasource.Spec.Config`. |
| `datasource.WithCloseGrace(d)` | Sets `RegisterOptions.CloseGrace` only when `d > 0`; delayed closes are asynchronous and are cut short by shutdown. |
| `datasource.WithReconcileDryRun()` | Sets `ReconcileOptions.DryRun` so `Reconcile` reports the diff without opening or closing connections. |

`datasource.Registry` methods:

| Method | Contract |
| --- | --- |
| `datasource.Registry.Primary` | Returns the primary `orm.DB`; equivalent to `Get(datasource.PrimaryName)` but does not return an error. |
| `datasource.Registry.Get` | Returns the registered `orm.DB`; unknown names return `datasource.ErrNotFound`. |
| `datasource.Registry.Has` | Reports whether a name is currently registered. |
| `datasource.Registry.Names` | Returns all registered names, including `primary`, in stable lexical order. |
| `datasource.Registry.Kind` | Returns the configured `config.DBKind`; unknown names return `datasource.ErrNotFound`. |
| `datasource.Registry.Register` | Opens and pings a new non-primary source before inserting it. Duplicate names return `datasource.ErrExists`; a failed conflict path closes the new pool. |
| `datasource.Registry.Update` | Opens and pings the replacement before swapping it in. Failed updates leave the old entry untouched; successful updates close the old pool asynchronously and honor `datasource.WithCloseGrace`. |
| `datasource.Registry.Unregister` | Removes a non-primary source, then closes the old pool asynchronously and honors `datasource.WithCloseGrace`. |
| `datasource.Registry.Reconcile` | Serializes reconcile calls, ignores empty and `primary` specs, sorts add/update/remove buckets, and records per-name failures in `ReconcileReport.Errors` without returning a top-level error for partial failure. |
| `datasource.Registry.TestConnection` | Opens a throwaway connection, queries the server version, closes it, returns `datasource.ConnectionInfo`, and never mutates the registry. The probe is capped by the internal 5s timeout while still respecting caller cancellation/deadline. |
| `datasource.Registry.HealthCheck` | Pings primary and all non-primary entries in parallel and returns a name-to-error map; nil errors mean reachable sources. |

## SPA

- `vef:spa`

Helpers:

- `vef.ProvideSPAConfig(...)`
- `vef.SupplySPAConfigs(...)`

## Storage integration

Extension groups:

- `vef:api:resources`
- `vef:app:middlewares`

Helpers:

- `vef.SupplyFileACL(...)`
- `vef.SupplyURLKeyMapper(...)`

`SupplyFileACL` replaces the default `storage.FileACL`. The default only grants
read access to keys under `pub/`; applications that store private files should
provide a business-specific ACL.

`SupplyURLKeyMapper` replaces the default `storage.ProxyURLKeyMapper`, which
maps `/storage/files/<key>` proxy URLs back to object keys. Override it when
rich-text or markdown content embeds CDN URLs or any other URL form that must
map back to storage object keys. Use `storage.IdentityURLKeyMapper` only when
content embeds bare object keys directly. This default is for the framework DI
graph; direct `storage.NewFiles(...)` calls normalize a nil mapper to
`storage.IdentityURLKeyMapper`.

## MCP

- `vef:mcp:tools`
- `vef:mcp:resources`
- `vef:mcp:templates`
- `vef:mcp:prompts`

Helpers:

- `vef.ProvideMCPTools(...)`
- `vef.ProvideMCPResources(...)`
- `vef.ProvideMCPResourceTemplates(...)`
- `vef.ProvideMCPPrompts(...)`
- `vef.SupplyMCPServerInfo(...)`

## Approval lifecycle hooks

- `vef:approval:lifecycle_hooks`

Helper:

- `vef.ProvideApprovalLifecycleHook(...)`

`approval.InstanceLifecycleHook` implementations run synchronously inside the
approval engine transaction for lifecycle moments such as instance creation and
completion. Returning an error rolls back the surrounding approval command. Use
event subscriptions for asynchronous integrations that should run after commit.

## Approval business binding

Helper:

- `vef.SupplyBusinessBindingHook(...)`

This replaces the default `approval.BusinessBindingHook` used when
`Flow.BindingMode == BindingBusiness`. Implementations bridge approval instances
to the host application's business rows and must be idempotent on asynchronous
status write-back.

## Logging

Helper:

- `vef.NamedLogger(name)`

Use this root-package convenience function when integration code needs a
framework `logx.Logger` outside dependency injection. It returns
`logx.Logger`; the `logx` package itself exposes only the `Level` constants,
`Level.String()`, and the `Logger` interface contract.

## See also

- [Modules & Dependency Injection](../modules/overview) for how these groups fit into app composition
- [Custom Param Resolvers](../advanced/custom-param-resolvers) for handler injection extension
