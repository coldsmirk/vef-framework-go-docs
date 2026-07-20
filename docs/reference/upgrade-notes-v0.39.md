---
sidebar_position: 9
---

# Upgrade Notes: v0.39

This page is the cross-version audit map for the backend commits from
`v0.38.0` through `v0.39.0` and the fixes that followed it. The release
centers on three new optional modules — the
[integration engine](../integration/overview), the
[durable cron schedule store](../infrastructure/cron-store), and
[WebSocket server push](../infrastructure/push) — plus a rebuilt JS engine,
a new outbound HTTP client, SQL Server / Oracle connectivity, and a set of
breaking renames.

It is not a replacement for the generated indexes. After applying the
migration notes below, verify exact Go symbols and wire fields against the
[Public API Index](./public-api-index).

## Immediate Checklist (Breaking)

- **`httpx` → `fiberx` rename**: the Fiber request helpers (`IsJSON`,
  `IsMultipart`, `GetIP`) moved from
  `github.com/coldsmirk/vef-framework-go/httpx` to `.../fiberx`. The `httpx`
  import path now hosts the new [outbound HTTP client](../utilities/httpx).
  Update imports mechanically.
- **JS engine rebuilt**: `js.New()` is gone. Build a shared `js.Engine` with
  `js.NewEngine(...)` and stamp per-execution runtimes with
  `engine.NewRuntime(...)`; run with `rt.RunString(ctx, src)` /
  `rt.RunProgram(ctx, prog)` (context-first, cancellation-aware). The
  vendored `dayjs`/`Big`/`utils`/`validator` globals were replaced by one
  stdlib bundle: `BigNumber`, `dayjs`, `fxp`, `radashi`, `z`, and `URL` /
  `URLSearchParams`. See [JS Engine](../data-tools/js-engine).
- **JS SQL bindings renamed**: script-side `sql` verbs are now
  `sql.queryList` / `sql.queryOne` / `sql.execute`; the Go option is
  `jssql.WithExecute()` and the sentinel is `jssql.ErrExecuteDisabled`.
  `sql.queryList` now returns `[]` (not `null`) when no rows match.
- **`mold` dictionary → code set rename**: tag prefix `dict:` became
  `codes:`; `Dictionary*` identifiers became `CodeSet*`
  (`CachedCodeSetResolver`, `CodeSetChangedEvent`,
  `PublishCodeSetChangedEvent`, ...). See the
  [rename map](../data-tools/mold#v038--v039-rename-map).
- **`cryptox.NewSM4` defaults to GCM**: SM4 ciphertext from earlier versions
  was CBC; decrypt it with
  `cryptox.NewSM4(key, cryptox.WithSM4Mode(cryptox.Sm4ModeCbc))`. AES was
  already GCM-default; the two ciphers now match.
- **Security — reserved identities enforced**: authenticators, challenge
  providers, and token issuance now refuse principals that claim
  framework-reserved identities (`security.Principal.IsReserved()`: the
  `system` type, `orm.OperatorSystem`, `orm.OperatorCronJob`). Audit any
  custom `Authenticator` / `UserLoader` that could produce such IDs.
- **API — permission declarations validated**: contradictory or malformed
  `RequiredPermission` declarations are rejected at resource registration
  (boot fails fast). Permission tokens follow the dot-separated convention
  (`domain.entity.action`).
- **Approval**:
  - `approval/category.find_tree_options` was removed.
  - `approval/flow.update` params type renamed `UpdateParams` →
    `UpdateFlowParams` (wire shape unchanged).
  - `find_versions` now returns version summaries without graph documents;
    fetch a specific version through `get_graph` with the new
    `params.versionId`.
  - Error code `40015` (flow-binding lock) is retired and never reassigned.

## New Modules (Opt-In)

- **Integration engine** (`vef.IntegrationModule`): contracts, systems,
  adapters (outbound + inbound), routes, per-system code maps, invocation
  logs and statistics, dry-run consoles, and the inbound HTTP gateway
  (`POST /integration/inbound/:systemCode/:contractCode`). Configuration
  lives under `vef.integration`. Start at
  [Integration Engine](../integration/overview).
- **Durable cron schedules** (`vef.cron.store.enabled = true`): persisted
  schedules with cluster-wide single fire, misfire/concurrency policies, a
  run journal, crash recovery, and the `sys/cron/schedule` /
  `sys/cron/run` resources. Register handlers with
  `vef.ProvideCronJobHandler`. See
  [Durable Schedules](../infrastructure/cron-store).
- **Server push** (`vef.push.enabled = true`): WebSocket endpoint (default
  `/ws`), `push.Notifier` for user/role/broadcast targeting, session-revocation
  integration, and an automatic Redis relay across nodes. See
  [Server Push](../infrastructure/push).

## New Capabilities (Non-Breaking)

- **`httpx` outbound HTTP client**: fluent per-call request builder over an
  immutable per-upstream client — base URL, default headers/query, retries
  with jittered backoff (idempotent methods; 429/502/503/504), hooks,
  proxy/TLS/redirect/body-size controls. See
  [httpx](../utilities/httpx).
- **JS capability libraries**: `jshttp` (fetch-style HTTP), `jssql`
  (guarded SQL), `jscache`, `jsevents`, `jscrypto`, `jsconsole` — opt-in per
  runtime through the engine catalog.
- **SQL Server and Oracle datasources**: `type = "sqlserver"` (default port
  1433) and `type = "oracle"` (service name in `database`, default port
  1521) are now real connection providers with read-only SQL-guard parser
  coverage. See [Datasources](../data-access/datasources).
- **Security auth strategies**: `api.APIKeyAuth(...)` (`X-API-Key` header by
  default, config-backed `vef.security.api_keys`) and `api.HTTPBasicAuth()`
  (config-backed `vef.security.basic_accounts`), both with replaceable
  loaders. See [Authentication](../security/authentication).
- **Session revocation seam**: `security.SessionRevocationListener` /
  `SessionRevocationNotifier` (used by push to kick dead sessions
  instantly). See
  [Session Management](../security/session-management#revocation-listeners-v039).
- **Approval**: host-owned flow `labels` with equality filtering
  (`find_flows`, `my.find_available_flows`), labels surfaced in instance
  details, self-service `my.get_start_form`, and the viewer task context
  (`myTask` with rollback targets and removable assignees) in
  `my.get_instance_detail`. See
  [Approval RPC Resources](../approval/resources).
- **Monitor**: `sys/monitor.get_integration_stats` reports per-node
  integration invocation statistics.
- **timex**: `DateTime.AsLocal()` reinterprets naive wall-clock fields in
  the process-local zone before instant arithmetic.

## Behavioral Notes

- Integration secrets at rest: configure `vef.integration.secret_key`
  (AES-GCM default; `secret_algorithm = "sm4"` for SM4-GCM). Switching
  algorithms requires re-entering stored secrets.
- Cron store times: caller-supplied schedule times are normalized to the
  local wall clock before persistence; cron triggers without a timezone
  evaluate in UTC, never the node's process-local zone.
- Push relay: enabled implicitly by `vef.redis.enabled = true`; requires
  `vef.app.name` for channel namespacing.
- The `security/auth` login flow continues to refuse framework-issued token
  types as credentials; the reserved-identity gate is additional hardening
  on top.

## See also

- [Upgrade Notes: v0.38](./upgrade-notes-v0.38) for the previous release
- [Configuration Reference](./configuration-reference) for the new
  `vef.cron`, `vef.integration`, and `vef.push` sections
- [Built-in Resources](./built-in-resources) for the expanded resource index
