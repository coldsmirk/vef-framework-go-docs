---
sidebar_position: 3
---

# Extension Points

Most VEF extension points are explicit FX groups.

This page covers the DI extension helpers that register into FX groups, replace defaults through `fx.Decorate`, or supply singleton values through `fx.Supply`.

## Helper mechanisms

The helper name prefix is not enough to tell how the value is wired. Use the mechanism:

| Mechanism | Helpers |
| --- | --- |
| `fx.Provide` + `fx.ResultTags` group append | `ProvideAPIResource`, `ProvideAuthStrategy`, `ProvideMiddleware`, `ProvideSPAConfig`, `ProvideCQRSBehavior`, `ProvideChallengeProvider`, `ProvideMCPTools`, `ProvideMCPResources`, `ProvideMCPResourceTemplates`, `ProvideMCPPrompts`, `ProvideEventTransport`, `ProvideEventPublishMiddleware`, `ProvideEventConsumeMiddleware`, `ProvideApprovalLifecycleHook`, `ProvideApprovalAggregator`, `ProvideDataSourceProvider`, `ProvideJSLib`, `ProvideCronJobHandler`, `ProvideIntegrationOutboundAuthScheme`, `ProvideIntegrationInboundAuthScheme`, `ProvideIntegrationInboundHandler`, `ProvideSessionRevocationListener` |
| `fx.Supply` with group tags | `SupplySPAConfigs` |
| `fx.Decorate` replacement | `SupplyFileACL`, `SupplyURLKeyMapper`, `SupplyBusinessRefProvider`, `SupplyBusinessRefResolver`, `ProvideEventMetricsRecorder`, `ProvideEventErrorSink`, `ProvideApprovalFormSchemaParser` |
| plain `fx.Supply` value | `SupplyMCPServerInfo` |

Replacement helpers are single-service overrides, not append-only extension
points. Register only one implementation unless you intentionally want a later
FX option to replace an earlier one.

## API and app groups

- `vef:api:resources`
- `vef:api:auth_strategies`
- `vef:app:middlewares`

Helpers:

- `vef.ProvideAPIResource(...)`
- `vef.ProvideAuthStrategy(...)`
- `vef.ProvideMiddleware(...)`

`ProvideAuthStrategy` appends a custom `api.AuthStrategy` into the auth-strategy
group. The strategy is selected by the name returned from `Name()` through
`api.AuthConfig.Strategy`; built-in strategies are `none`, `bearer`,
`signature`, `ip`, `api_key`, and `http_basic`.

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
- `vef:security:session_revocation_listeners`

Helpers:

- `vef.ProvideChallengeProvider(...)`
- `vef.ProvideSessionRevocationListener(...)` — appends a
  `security.SessionRevocationListener` observing logout, concurrent-login
  eviction, and administrative kicks; see
  [Session Management](../security/session-management#revocation-listeners)

## Cron job handlers

- `vef:cron:job_handlers`

Helper:

- `vef.ProvideCronJobHandler(...)`

Registers a durable `cron.JobHandler` with the schedule store — exactly one
handler per job name; duplicate names fail startup. Handlers optionally ship
a default schedule via `cron.WithDefaultSchedule`. See
[Durable Schedules](../infrastructure/cron-store).

## JS engine libraries

- `vef:js:libs`

Helper:

- `vef.ProvideJSLib(...)`

Contributes a `js.Lib` to the shared `js.Engine`: a lib whose name matches a
built-in replaces it within its tier (always-on vs opt-in catalog), a new
name joins the opt-in catalog. See [JS Engine](../data-tools/js-engine).

## Integration engine (requires `vef.IntegrationModule`)

- `vef:integration:outbound_auth_schemes`
- `vef:integration:inbound_auth_schemes`
- `vef:integration:inbound_handlers`

Helpers:

- `vef.ProvideIntegrationOutboundAuthScheme(...)` — custom
  `integration.OutboundAuthScheme`; a name matching a built-in replaces it
- `vef.ProvideIntegrationInboundAuthScheme(...)` — custom
  `integration.InboundAuthScheme`; same replacement rule
- `vef.ProvideIntegrationInboundHandler(...)` — the business handler serving
  one inbound contract; exactly one handler per contract code

The default `integration.RouteResolver` is replaced with plain
`fx.Decorate` when routing lives outside the `itg_route` table. See
[Integration Engine](../integration/overview).

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
approval engine transaction — `OnInstanceCreated` at instance creation and
`OnInstanceTransition(from, to)` at every instance status transition.
Returning an error rolls back the surrounding approval command. Invocation
order across hooks is unspecified (FX value groups carry no ordering). Use
event subscriptions or `approval.BindCommand` for asynchronous integrations
that should run after commit.

## Approval aggregators and form schema

- `vef:approval:aggregators`

Helpers:

- `vef.ProvideApprovalAggregator(...)`
- `vef.ProvideApprovalFormSchemaParser(...)`

`ProvideApprovalAggregator` registers a custom detail-table aggregator for
approval field conditions alongside the built-in sum / count / avg. The
constructor must return an `approval.Aggregator`; the condition evaluator picks
it up by its `AggregateKind`. Boot fails if a built-in aggregate kind is left
unregistered.

`ProvideApprovalFormSchemaParser` replaces the framework's default
`approval.FormSchemaParser` (the built-in vef-framework-react form-editor
parser). The replacement is wholesale, not additive: every deployed form
schema goes through it, so it must understand every designer document the host
submits. Parsing runs once at flow deploy; versions deployed earlier keep the
`form_fields` they were persisted with.

## Approval business binding

Helpers:

- `vef.SupplyBusinessRefProvider(...)`
- `vef.SupplyBusinessRefResolver(...)`

`SupplyBusinessRefProvider` replaces the default no-op
`approval.BusinessRefProvider`. It runs inside `start_instance` when
`Flow.BindingMode == BindingBusiness` and lets the host resolve or allocate the
business row, returning the opaque `Instance.BusinessRef`.

`SupplyBusinessRefResolver` replaces the default
`approval.BusinessRefResolver`, which resolves `Instance.BusinessRef` into the
`approval.BusinessRecordKey` matched against the flow's configured
`BusinessBindingConfig.KeyColumns` (single-column refs verbatim, composite
refs as a JSON object). Register one when the ref uses another shape or
resolving the key requires a host lookup.

The business-state projection itself is owned by the engine. Hosts extend
around it with `approval.InstanceLifecycleHook`, event subscriptions, or
`approval.BindCommand`; they no longer replace the projection path.

## Logging

Helper:

- `vef.NamedLogger(name)`

Use this root-package convenience function when integration code needs a
framework `logx.Logger` outside dependency injection. It returns
`logx.Logger`; the `logx` package itself exposes the `Level` constants,
`Level.String()`, the `Logger` interface contract, and
`LoggerConfigurable[T]` for immutable components that return a
logger-configured copy from `WithLogger`.

## A simple decision rule

When deciding how to extend VEF, ask:

- has the framework already reserved a group for this concept?
- should the extension participate in startup and lifecycle management?
- do you want dependencies between modules to stay explicit and testable?

If the answer is yes, route through an FX group instead of implicit global
state or a hand-rolled singleton.

## See also

- [Modules & Dependency Injection](../core-concepts/overview) for how these groups fit into app composition
- [Extending Handler Parameters](../advanced/extending-parameters) for handler injection extension
