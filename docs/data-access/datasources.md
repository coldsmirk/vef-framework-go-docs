---
sidebar_position: 10
---

# Multiple Data Sources

Most applications only ever talk to one database: the **primary** data source, injected everywhere as `orm.DB`. This page is for the rest — reporting warehouses, per-tenant databases, legacy systems you read from but do not own — reached through `datasource.Registry`.

## Primary vs Additional Sources

The primary source is declared under `vef.data_sources.primary` in TOML. It is mandatory, it is the source exposed framework-wide as `orm.DB`, and it cannot be mutated through the dynamic registry API — `Register`, `Update`, and `Unregister` all reject `datasource.PrimaryName` (`"primary"`) with `datasource.ErrPrimaryReserved`.

Every other source is "additional": either declared statically in TOML under a different name, or registered dynamically at runtime. Reach an additional source through `datasource.Registry`, injected wherever you need it.

Internal framework modules — CRUD, approval, storage, event inbox/outbox, schema reflection — all operate on the primary source only. Reaching an additional source, and deciding what to do with it, is an application-level concern.

## Static Sources: TOML

The simplest way to add a source is another `vef.data_sources.<name>` table, alongside `primary`:

```toml
[vef.data_sources.primary]
type = "postgres"
host = "127.0.0.1"
port = 5432
user = "postgres"
password = "postgres"
database = "my_app"
schema = "public"

[vef.data_sources.analytics]
type = "sqlite"
path = "./analytics.db"
```

Every entry uses the same `config.DataSourceConfig` shape:

| Field | Type | Meaning |
| --- | --- | --- |
| `type` | `postgres \| mysql \| sqlite` | database kind (`oracle` and `sqlserver` constants exist but are not implemented yet) |
| `host` | `string` | network database host |
| `port` | `uint16` | network database port |
| `user` | `string` | database username |
| `password` | `string` | database password |
| `database` | `string` | database name |
| `schema` | `string` | schema name for drivers that support schemas |
| `path` | `string` | SQLite file path |
| `enable_sql_guard` | `bool` | enables the SQL guard for raw SQL surfaces |
| `ssl_mode` | `disable \| require \| verify-ca \| verify-full` | TLS posture for network dialects; defaults to `disable` |
| `ssl_root_cert` | `string` | optional PEM path for `verify-ca` / `verify-full` |

Every non-primary entry under `vef.data_sources` is registered into `datasource.Registry` under its map key before the application starts serving requests. See [Configuration Reference](../reference/configuration-reference) for the full field list.

## Injecting The Registry

`datasource.Registry` is available in the FX container everywhere. Inject it directly:

```go
package report

import (
	"context"

	"github.com/coldsmirk/vef-framework-go/datasource"
)

type Service struct {
	sources datasource.Registry
}

func NewService(sources datasource.Registry) *Service {
	return &Service{sources: sources}
}

func (s *Service) RunReport(ctx context.Context) error {
	analytics, err := s.sources.Get("analytics")
	if err != nil {
		return err
	}

	var count int
	return analytics.NewSelect().
		ColumnExpr("count(*)").
		Table("events").
		Scan(ctx, &count)
}
```

`Get` returns an `orm.DB`, so it exposes the same query-building API as the primary source (`NewSelect`, `NewInsert`, `NewRaw`, and so on) — see [Query Builder](./query-builder). Unknown or unregistered names return `datasource.ErrNotFound`. `sources.Primary()` returns the same `orm.DB` you would get from injecting `orm.DB` directly; it never errors.

`datasource.Registry` is also a built-in API handler parameter — you can request it directly as a resource handler parameter alongside `orm.DB` and `fiber.Ctx` — see [Custom Handlers](../building-apis/custom-handlers).

## Runtime Sources: `datasource.Provider`

For sources that are not known at deploy time — for example, a tenant table in the primary database whose rows describe additional databases — implement `datasource.Provider`:

```go
type Provider interface {
	Name() string
	Load(ctx context.Context) ([]Spec, error)
}

type Spec struct {
	Name   string
	Config config.DataSourceConfig
}
```

The framework calls `Load` once during startup, **after** the primary and static TOML sources are already registered, and `Register`s every returned `Spec`. A name collision with TOML or another provider fails boot. `Provider.Name` only labels the provider in diagnostics — it is not the data source name.

Register the provider with `vef.ProvideDataSourceProvider`:

```go
func NewTenantSourceProvider(primary orm.DB) datasource.Provider {
	return &tenantSourceProvider{primary: primary}
}

// in your fx.Module:
vef.ProvideDataSourceProvider(NewTenantSourceProvider)
```

Provider order across multiple registered providers is undefined, so specs from different providers must not collide on name.

## Runtime Sources: Direct `Register`

Outside of the startup `Provider` hook — for example, from an admin endpoint that lets an operator add a data source on demand — call `Registry.Register` directly:

```go
db, err := sources.Register(ctx, "tenant-42", config.DataSourceConfig{
	Kind:     config.Postgres,
	Host:     "tenant-42.internal",
	Port:     5432,
	User:     "app",
	Password: pw,
	Database: "tenant_42",
})
```

`Register` opens and pings the new connection before inserting it; on any conflict (`datasource.ErrExists` for a duplicate name, `datasource.ErrPrimaryReserved` for `"primary"`, `datasource.ErrNameInvalid` for an empty or whitespace/control-character name) the freshly opened connection is closed and nothing is inserted. `Register` never closes an existing connection, so it takes no options.

## The Registry Surface

| Method | Contract |
| --- | --- |
| `Primary()` | The primary `orm.DB`. Equivalent to `Get(datasource.PrimaryName)` but never errors. |
| `Get(name)` | The registered `orm.DB`. Returns `datasource.ErrNotFound` for an unregistered or since-unregistered name. |
| `Has(name)` | Reports whether `name` is currently registered and not closed. |
| `Names()` | Every registered name, including `primary`, in stable lexical order. |
| `Kind(name)` | The `config.DBKind` for `name`. Same not-found semantics as `Get`. |
| `Register(ctx, name, cfg)` | Opens and pings a new non-primary source, then inserts it. |
| `Update(ctx, name, cfg, opts...)` | Atomically swaps the connection for an existing source. The old pool is closed asynchronously. |
| `Unregister(ctx, name, opts...)` | Removes a non-primary source. The old pool is closed asynchronously. |
| `Reconcile(ctx, specs, opts...)` | Drives the registry toward a desired set of non-primary sources in one call. |
| `TestConnection(ctx, cfg)` | Opens a throwaway connection, verifies it, closes it. Never mutates the registry. |
| `HealthCheck(ctx)` | Pings every registered source in parallel; returns a `name -> error` map. |

All read methods (`Get`, `Has`, `Names`, `Kind`, `Primary`) are safe for concurrent use. `Register`, `Update`, and `Unregister` mutate the registry atomically.

## Reconciling A Desired Set

When your source list is derived from an external table (the same tenant-table scenario as the `Provider` example), periodic drift between that table and the registry is common: rows get added, updated, or deleted between application restarts. `Reconcile` closes that gap in one call, without you hand-rolling the diff:

```go
report, err := sources.Reconcile(ctx, specs)
```

Given a desired `[]datasource.Spec`, `Reconcile` computes three buckets and drives the registry toward them:

- a spec with no matching registry entry → `Register`
- a spec whose config differs from the current entry → `Update`
- a registry entry with no matching spec → `Unregister`

Specs referencing the primary name are ignored. Per-name failures are collected in `ReconcileReport.Errors` (keyed by name) without aborting the rest of the batch — one bad config in a batch of ten does not block the other nine. Reconciles are serialized: two ticks of a refresher job (typically a cron job calling `Reconcile` on a schedule) can never interleave and race each other. Direct `Register` / `Update` / `Unregister` calls, however, are **not** synchronized against a running `Reconcile`.

Use `datasource.WithReconcileDryRun()` to compute the report without opening or closing anything — useful for previewing what a refresher job would do:

```go
preview, _ := sources.Reconcile(ctx, specs, datasource.WithReconcileDryRun())
```

## Updating And Removing Sources

`Update` and `Unregister` are the two operations that close an existing connection, and both accept `datasource.RegisterOption`. By default the replaced or removed pool closes immediately on a background goroutine; `WithCloseGrace(d)` delays that close so in-flight queries have time to drain:

```go
_, err := sources.Update(ctx, "analytics", newCfg, datasource.WithCloseGrace(10*time.Second))
```

```go
err := sources.Unregister(ctx, "analytics", datasource.WithCloseGrace(10*time.Second))
```

Either way, once the call returns, `Get("analytics")` reflects the new state immediately (the new config for `Update`, `datasource.ErrNotFound` for `Unregister`) — the grace period only affects when the *old* underlying `*sql.DB` closes, so a caller that already holds an `orm.DB` reference from before the swap can finish its in-flight queries.

## Testing And Health-Checking Connections

`TestConnection` is a pure connectivity probe: it opens a throwaway connection from a candidate config, confirms it by querying the server version, and closes it — it never touches the registry. It is the natural backend for a "test connection" button in an admin UI, run before calling `Register` or `Update`:

```go
info, err := sources.TestConnection(ctx, candidateCfg)
if err != nil {
	// unreachable or unusable
}
// info.Version, e.g. "PostgreSQL 16.2 on x86_64-pc-linux-gnu"
```

`HealthCheck` instead pings every currently registered source (including primary) in parallel and returns a `name -> error` map, useful for a liveness/readiness endpoint that needs to see every data source at once.

## Constraints To Keep In Mind

- **Internal modules are primary-only.** CRUD, approval, storage, event inbox/outbox, and schema reflection all read and write the primary source. Additional sources are never touched by framework internals — only by your own code.
- **Cross-source transactions are not supported.** `orm.DB.RunInTx` opens a transaction against a single source. If you publish an event with `event.WithTx(tx)`, `tx` must come from a transaction opened on the primary source — see [Transactions](./transactions).
- **`Reconcile` only manages non-primary sources.** Specs that reference `datasource.PrimaryName` are silently ignored, so a Provider or reconcile job can safely include the primary name in its input without risking a boot failure.

## Next Step

For the model and query-building side of things once you have an `orm.DB` from `Get` or `Primary`, see [Query Builder](./query-builder) and [Transactions](./transactions). For every other framework extension point, including `vef.ProvideDataSourceProvider`, see [Extension Points](../reference/extension-points).
