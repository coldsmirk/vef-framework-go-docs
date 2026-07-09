---
sidebar_position: 4
---

# Upgrade Notes to v0.35

This page is the cross-version audit map for the backend commits from
`9e7e009 feat: added multi-data source functionality` through `v0.35.0`
(`9e7e009^..v0.35.0`). Use it when upgrading an application whose docs or
integration assumptions were last checked before the multi-data-source work.

It is not a replacement for the generated indexes. After applying the migration
notes below, verify exact Go symbols and wire fields against the
[Public API Index](../reference/public-api-index) and
[Runtime API Index](../reference/runtime-api-index).

## Immediate Checklist

- Replace any old single-source database config with `[vef.data_sources.primary]`.
  The `primary` entry is mandatory and reserved; additional sources live under
  `[vef.data_sources.<name>]`.
- Update renamed Go symbols: the old env key prefix constant is now
  `config.EnvPrefix`, the old CORS config type is now `config.CORSConfig`,
  `datasource.Spec.Cfg` is now `datasource.Spec.Config`, and the SSL constants are
  `SSLDisable`, `SSLRequire`, `SSLVerifyCA`, and `SSLVerifyFull`.
- If a client signs API requests, update the canonical payload to include the
  HTTP method and path:
  `app_id=<appID>&method=<method>&nonce=<nonce>&path=<path>&timestamp=<timestamp>`.
  The request body is intentionally not part of the HMAC payload.
- Treat `/mcp` as Bearer-protected by default. Set `vef.mcp.require_auth=false`
  only when the MCP surface is deliberately anonymous.
- Enable approval explicitly with `vef.ApprovalModule`; it is no longer part of
  the default boot graph. Approval events must route through a transactional
  transport and any framework subscriber also needs a subscribable sink.
- Review approval designer/runtime contracts: node config is normalized and
  validated at deploy time, `PassRule` is limited to `all`, `any`, and `ratio`,
  instance detail now includes `formSchema`, `timeline`, and `flowGraph`, and
  business binding uses opaque `businessRef` plus engine-owned status write-back.
- Use the storage multipart upload protocol (`init_upload`, `upload_part`,
  `list_parts`, `complete_upload`, `abort_upload`) and the
  `/storage/files/<key>` proxy. Public uploads require
  `vef.storage.allow_public_uploads=true`.
- For `vef-cli generate-model-schema`, regenerate schemas when models use Bun
  table aliases, bare table tags, default table names, or `m2m` relations; the
  generator now follows Bun tag parsing more closely.

## Release-by-Release Audit

| Release | User-facing changes to review |
| --- | --- |
| `v0.27.0` | Multi-data-source config and registry, mandatory `primary`, `datasource.Provider`, `Registry.Register` / `Update` / `Unregister` / `Reconcile` / `TestConnection` / `HealthCheck`, and the removal of the legacy `vef.data_source` fallback. |
| `v0.28.0` | Core expression engine, source-IP security checks, MCP read-only query guard, event middleware ordering and tracing, storage and CRUD API error surfaces, approval opt-in module and route checks, and stricter request error handling. |
| `v0.29.0` / `v0.29.1` | `Effective*()` config defaults, database TLS support, dialect-aware SQL guard hardening, approval tenant fail-closed behavior, signature replay-window fixes, and `security.UserMenu` JSON alignment. |
| `v0.30.0` | Restored sequence DB/Redis stores, `cache.NewPrefixKeyBuilderWithSeparator`, `logx.LoggerConfigurable`, and `reflectx.BreadthFirst`. |
| `v0.31.0` | The expression backend moved to pure-Go `expr-lang`; expression itself no longer requires `CGO_ENABLED=1`. |
| `v0.32.x` | Config API renames (`EnvPrefix`, `CORSConfig`, SSL constants), restored approval module after hardening, and Bun-compatible model-schema generation fixes. |
| `v0.33.x` | Approval detail DTOs gained `formSchema`, `timeline`, flow graph progress, action-log metadata, and persistent flow-node ids for rollback targeting; model-schema parsing became a faithful Bun tag parser. |
| `v0.34.0` | Custom auth strategies via `vef.ProvideAuthStrategy` and source-IP whitelist auth via `api.IPAuth(...)`, `vef.security.ip_whitelists`, and `vef.app.trusted_proxies`. |
| `v0.35.0` | Approval instance detail moved to the engine-recorded visit trail with person snapshots; condition routing snapshots host globals; business binding uses opaque refs; approval domain events were rebuilt as self-describing envelopes. |

## Data Sources and Database

Configuration now uses a named data-source map:

```toml
[vef.data_sources.primary]
type = "postgres"
host = "127.0.0.1"
port = 5432
user = "postgres"
password = "postgres"
database = "app"
schema = "public"

[vef.data_sources.analytics]
type = "sqlite"
path = "./analytics.db"
```

`primary` is both `config.PrimaryDataSourceName` and
`datasource.PrimaryName`. It powers the framework-wide `orm.DB` injection and
cannot be dynamically registered, updated, or unregistered. Runtime providers
currently exist for `postgres`, `mysql`, and `sqlite`; the `oracle` and
`sqlserver` `DBKind` constants are reserved for future providers.

The `datasource.Registry` is the public extension point for named runtime
sources. It exposes `Primary`, `Get`, `Has`, `Names`, `Kind`, `Register`,
`Update`, `Unregister`, `Reconcile`, `TestConnection`, and `HealthCheck`.
`TestConnection` opens a throwaway connection, returns
`datasource.ConnectionInfo.Version`, and does not mutate the registry.

Network database sources can opt into TLS with `ssl_mode` and `ssl_root_cert`.
Use `disable`, `require`, `verify-ca`, or `verify-full`. SQLite ignores these
fields.

## API, Security, and MCP

The auth strategy registry is now extensible through
`vef.ProvideAuthStrategy(...)`. The built-in strategies are `none`, `bearer`,
`signature`, and `ip`; `api.IPAuth()` uses the `default` whitelist and
`api.IPAuth("ops")` targets `vef.security.ip_whitelists.ops`.

Behind a proxy, configure `vef.app.trusted_proxies`; otherwise forwarded headers
from untrusted clients are ignored and IP-based auth sees the direct peer.

Request decoding and auth failures are now more explicit:

- malformed params and meta return `api.ErrInvalidRequestParams` /
  `api.ErrInvalidRequestMeta` with HTTP 400 semantics;
- invalid Bearer tokens return 401 rather than falling through as anonymous;
- public operations get an isolated anonymous principal.

MCP database queries are read-only and dialect-aware, and `/mcp` requires Bearer
auth unless `vef.mcp.require_auth=false`.

## Event Bus

Applications should treat events as the current multi-transport platform:
memory, transactional outbox, Redis Streams, optional Inbox dedupe, deterministic
publish/consume middleware ordering, W3C trace propagation, and route
inspection through `event.RouteInspector`.

Routes that use the transactional outbox for framework subscribers must include
a subscribable sink such as `redis_stream` or `memory`; an `["outbox"]`-only
route is valid only for publisher-only flows. When Inbox dedupe is enabled, its
retention window must outlast the outbox retry horizon so delayed duplicates are
still recognized.

Approval and storage both publish transactional events and fail fast at startup
when required event routes are missing.

## Approval

Approval is now an optional feature module. Add `vef.ApprovalModule` to
`vef.Run(...)` before relying on any `approval/*` resources, CQRS handlers,
binding listeners, timeout scanners, or approval event publication.

Important runtime and API changes:

- deploy validates and normalizes flow node config before publication;
- conditions use approval's own `expr-lang` evaluator, not the public
  `expression.Engine`;
- host-resolved globals come from `approval.InstanceGlobalsResolver` at instance
  start and are snapshotted onto `Instance.Globals`;
- `FlowGraphNode.Kind` is `approval.NodeKind`;
- instance detail DTOs expose `formSchema`, `timeline`, and `flowGraph`;
- `FlowGraphNode.NodeID` is the persistent flow-node id used by action logs and
  rollback targets;
- business binding uses opaque `Instance.BusinessRef`; the engine owns the
  final status write-back and uses `BusinessRefResolver` only to extract the
  record id;
- domain events use `InstanceEventBase`, `TaskEventBase`, and `FlowEventBase`
  so subscribers can act from the event envelope without immediately loading
  approval tables;
- tenant handling is fail-closed for empty or ambiguous caller context.

See [Approval Module](../approval) for the full resource, DTO, event,
and extension-point contract.

## Expression and Mold

`expression.Engine` is part of the core boot graph and is injectable into API
handlers. The current backend is pure-Go `expr-lang`; the older Goja-based `js`
package remains separate. Mold's `expr` transformer evaluates against the
containing struct in declaration order, so derived fields can reference earlier
sibling fields.

## Storage

The storage resource is multipart-only at the HTTP layer. The business-side file
lifecycle uses `storage.Files` / `storage.FilesFor[T]`, claim consumption, and
pending-delete rows so file ownership changes commit or roll back with the
business transaction. Delete failures that exhaust retries emit
`vef.storage.delete.dead_letter` and are removed from the queue.

Default download URLs resolve through `/storage/files/<key>` unless a host
application supplies a custom `URLKeyMapper`. Public keys under `pub/` are
served anonymously; private keys call `storage.FileACL`.

## Other Public Surfaces

- CRUD mutation responses use operation-specific success messages, and CRUD /
  storage / approval outward errors use coded `result.Error` values.
- Result error sentinels freeze their i18n text at initialization time; language
  switches affect later translations but not already-constructed sentinel
  values.
- `i18n.SetLanguage(...)` is race-free with concurrent `T` / `Te` calls.
- `sequence.NewDBStore`, `sequence.NewRedisStore`, and injectable
  `*sequence.MemoryStore` are available for durable, distributed, and startup
  seeding scenarios.
- `cache.Invalidating[T]` is the event-invalidated read-through cache wrapper;
  `cache.NewPrefixKeyBuilderWithSeparator` preserves the separator exactly as
  provided.
- `logx.LoggerConfigurable[T]` and `reflectx.BreadthFirst` are part of the
  public utility surface again.
- The monitor module deduplicates disk devices by unique partitions rather than
  merging independent devices.
- Schema inspection uses the primary data source.
