---
sidebar_position: 8
---

# Upgrade Notes: v0.38

This page is the cross-version audit map for the backend commits from
`v0.37.0` through `v0.38.0` (`v0.37.0..v0.38.0`). The release centers on
security (login hardening and stateful opaque-token sessions), a redesigned
approval business binding (composite keys + durable desired-state
projection), a new distributed lock module, numeric fidelity in request
binding, and container-aware monitoring. Use it when upgrading an application
whose docs or integration assumptions were last checked against
[Upgrade Notes: v0.36 / v0.37](../reference/upgrade-notes-v0.37).

It is not a replacement for the generated indexes. After applying the
migration notes below, verify exact Go symbols and wire fields against the
[Public API Index](../reference/public-api-index) and
[Runtime API Index](../reference/runtime-api-index).

## Immediate Checklist

- **Security — token mechanism**: the `token` authenticator type string is now
  `jwt_token`, and `vef.security.token_type` selects `jwt_token` (default) or
  `opaque_token`. Only the configured mechanism's authenticators register:
  under `opaque_token` a leftover JWT no longer authenticates anywhere
  (including MCP) and the `refresh` operation is **not mounted**. `login` also
  refuses framework-issued token types as credentials.
- **Security — password transport**: `password.NewCipherEncoder` was removed.
  Register a `security.PasswordDecryptor` (e.g. `cryptox.NewRSA`) instead;
  decryption now happens in the authenticator, not the encoder.
- **Security — lockout**: brute-force lockout is on by default
  (`vef.security.lockout`, `max_failures = 10`) and, since v0.38, failed
  `resolve_challenge` guesses count toward the same lockout key.
- **API — numeric fidelity**: JSON `params`/`meta` are parsed with
  `json.Number`. Integer fields now reject fractional/exponent forms
  (`mapx.ErrJSONNumberNotInteger`) and out-of-range values
  (`mapx.ErrJSONNumberOverflow`) instead of silently truncating;
  `json.RawMessage` captures keep exact digits; untyped (`any`) targets still
  see `float64`.
- **API — rate limit config**: the default per-operation limit (100 req / 5
  min) is now configurable via `vef.api.rate_limit.max` / `.period`.
- **Approval — business binding redesign**: the six flat `Flow.Business*`
  columns are replaced by one `Flow.BusinessBinding` jsonb document
  (`approval.BusinessBindingConfig`: `tableName`, composite `keyColumns`,
  `statusColumn`, **mandatory** `instanceIdColumn`, optional timestamp
  columns, optional `statusMapping`), snapshotted per deployed version.
  Write-back is now a durable desired-state projection
  (`apv_business_projection`) with `synchronous` (default) or `eventual`
  consistency (`vef.approval.business_binding`). Update flow create/update
  clients (`params.businessBinding`) and permission seeds
  (`approval.binding.query` / `approval.binding.retry` for the new admin
  operations).
- **Approval — Go interfaces**:
  `InstanceLifecycleHook.OnInstanceCompleted(instance, finalStatus)` →
  `OnInstanceTransition(instance, from, to)` (fires on every transition);
  `BusinessRefResolver.ResolveRecordID` →
  `ResolveRecordKey(...) (BusinessRecordKey, error)`;
  `SubscribeInstance` handlers gain a third `env event.Envelope` parameter;
  `FormSchemaParser.ParseFormFields` gained a leading `ctx`.
- **Approval — form schema**: the structured `approval.FormDefinition` type is
  gone. `deploy` takes the host designer document opaque as
  `params.formSchema` (was `params.formDefinition`); `form_fields` are derived
  once at deploy via `FormSchemaParser`. Detail APIs return `formSchema`
  verbatim, and `my.get_instance_detail` adds a viewer-scoped
  `fieldPermissions` map.
- **Approval — migrations**: run the v0.38 DDL (`auto_migrate` or your own
  pipeline): `apv_flow.business_binding`, `apv_flow_version.business_binding`
  + `form_fields`, `apv_instance.business_projection_id`, and the new
  `apv_business_projection` table.
- **Approval — routing**: the binding listener is gone; the module now only
  requires a transactional route for `approval.*` (a bare `["outbox"]`
  passes). Keep a subscribable sink in the route only when the host itself
  subscribes (`SubscribeInstance` / the new `BindCommand`).
- **Monitor**: the overview disk summary now reports the root filesystem only
  (`partitions` is always `1`) and `vef.monitor.excluded_mounts` was removed;
  CPU summaries gained `effectiveCores` and memory figures respect cgroup
  limits — dashboards that summed partitions or normalized by host cores need
  adjusting.
- **Removals**: the `ptr` package (use `samber/lo` / builtin `new`; see
  [Small Helpers](../utilities/small-helpers)) and the snowflake ID generator
  (`id.DefaultSnowflakeIDGenerator`, `id.NewSnowflakeIDGenerator`, the
  `VEF_NODE_ID` env key) are gone — XID (`id.Generate()`) is the ID mechanism.
- **New**: inject `lock.Locker` for lease-based distributed locks
  (`lock.WithLock`, fencing tokens, Redis/memory backends selected by
  topology) — see [Distributed Lock](../infrastructure/lock).

## Release-by-Release Audit

| Release | User-facing changes to review |
| --- | --- |
| `v0.38.0` | Security: PasswordDecryptor transport, default-on lockout (+ challenge counting), password policy/history/expiry, `jwt_token` rename, opaque-token sessions with concurrency control, strict mechanism gating. API: numeric fidelity, configurable default rate limit. Approval: composite business bindings + durable projections, transition lifecycle hooks, `BindCommand`, delivery envelopes, opaque form schema + field permissions. New `lock` package. Monitor cgroup/root-disk overhaul. `ptr` and snowflake removals; sequence/cron/timex fixes. |

## Security

v0.38 turns the login path into a layered, opt-in hardening stack and adds a
stateful session mechanism. Full documentation:
[Login Hardening](../security/login-hardening) and
[Session Management](../security/session-management).

### Breaking: `jwt_token` rename and strict mechanism gating

The bearer authenticator's type string `token` became `jwt_token`
(`AuthTypeJWTToken`); `opaque_token` (`AuthTypeOpaqueToken`) joins it as the
stateful alternative selected by `vef.security.token_type`. Only the
configured mechanism's authenticators are registered — under `opaque_token`
the `refresh` operation is not mounted and JWTs stop authenticating on every
surface, including MCP. `login` rejects `jwt_token` / `opaque_token` /
`refresh` as login credential types, closing the "launder a stolen access
token into a refresh token" path.

### Breaking: `password.NewCipherEncoder` removed

Encrypted password transport moved from the encoder to the authenticator:
register a `security.PasswordDecryptor` and `PasswordAuthenticator` decrypts
before hash comparison. `password.ErrCipherRequired` / `ErrEncoderRequired`
were removed with the type. A malformed ciphertext costs a dummy KDF
comparison, closing a timing side channel.

### Opaque-token sessions

`OpaqueTokenGenerator` + `SessionStore` (memory default; Redis via
`fx.Decorate`) provide sliding idle TTL under an absolute `max_lifetime` cap
(clamped at issue time too), per-account concurrency limits
(`max_concurrent`, `on_exceed` = `evict_oldest` default / `reject`), logout
revocation, session-admin surfaces (`ListByUser` / `Revoke` / `RevokeUser`)
and the optional `SessionInspector.ListAll`. The memory store is built on the
TTL cache (GC'd, no unbounded growth); the Redis store performs atomic
multi-key mutations, self-heals its per-user set, and lists via keyspace
`SCAN`.

### Login hardening layers

`security.LoginGuard` (default-on lockout, `lock`/`backoff` strategies,
`user`/`ip`/`user_ip` keys, fail-open store errors, HTTP 429
`ErrAccountLocked`), config-driven password strength
(`vef.security.password_policy`, composable `PasswordRule`s), password
history (`PasswordHistoryStore` + `history_depth`), and password expiry
(`PasswordMetadataLoader` + `max_age` + `NewExpiryPasswordChangeChecker`).
The character-class rule counts caseless letters (CJK) toward no class, and
the identity rule measures tokens in runes.

## API

### Breaking: numeric fidelity in request binding

`api.Params` / `api.Meta` unmarshal with `json.Decoder.UseNumber`. Typed
numeric fields get exact-digit parses with `encoding/json`-grade strictness
(new sentinels `mapx.ErrJSONNumberNotInteger`, `mapx.ErrJSONNumberOverflow`);
`json.RawMessage` fields preserve full precision; `any` targets still receive
`float64`. Clients that relied on silent truncation of out-of-range or
fractional values now receive errors — that is the fix, not a regression.

### Configurable default rate limit

`vef.api.rate_limit.max` / `.period` set the default applied to operations
without their own `OperationSpec.RateLimit` (defaults 100 / 5m, keyed per
operation × client, counted per node).

## Approval

Full documentation: [Events & Integration](../approval/integration),
[Flow Design](../approval/flow-design), [RPC Resources](../approval/resources).

### Breaking: composite business bindings and durable projections

One `BusinessBindingConfig` document replaces the six flat `Flow` columns;
`KeyColumns` supports composite keys validated against a live non-null
primary/unique key; `InstanceIDColumn` is mandatory as a CAS fence;
`StatusMapping` translates statuses into the host vocabulary. Bindings are
snapshotted onto each deployed version. State converges through
`apv_business_projection` (claim at start with `ErrBindingTargetBusy` on a
busy target, desired-state revisions per transition, `synchronous` or
`eventual` apply with a leased worker and exponential backoff).
`InstanceBindingFailedEvent` is now an eventual-mode operator notification
only. New admin operations: `find_business_projections`,
`retry_business_projection`; `get_metrics` gains `businessProjectionCounts` /
`pendingBusinessProjections`. New error codes `40017`–`40020` and
`40107`–`40110`.

### Breaking: generalized lifecycle hooks and subscription envelopes

`OnInstanceTransition(instance, from, to)` fires inside every status
transition's transaction (order across hooks unspecified);
`SubscribeInstance` handlers receive the delivery `event.Envelope`
(`Envelope.ID` is the stable dedupe key); `BindCommand[E, C]` bridges
instance events to CQRS commands with command-identity derived consumer
groups (`vef:cmd:...`) and guards `ErrNonCommandAction` /
`ErrUnnamedCommandType`.

### Breaking: opaque form schema and field permissions

`FlowVersion.FormSchema` is the host designer document stored verbatim
(`json.RawMessage`); `FlowVersion.FormFields` is derived once at deploy via
the context-aware `FormSchemaParser` (replace with
`vef.ProvideApprovalFormSchemaParser`). The `FormDefinition` wrapper type was
removed and `deploy` renamed the parameter to `formSchema`. Field permissions
are validated at deploy (keys against fields, CC subset, `required` vs
`auto_pass` timeout), enforced on the write path (`editable`/`required`
merge; `required` checked on approve/handle only; permission-less nodes drop
submitted data), and projected per viewer in `my.get_instance_detail`
(`fieldPermissions`, max-merged lattice, fail-closed).

## New: Distributed Lock

The `lock` package provides lease-based locks with a topology-selected DI
default (`RedisLocker` when Redis is enabled, else `MemoryLocker` with a boot
warning), `WithLock` (auto-renew, lost-lease cancellation, panic-safe
release), fencing tokens, and retry-idempotent Redis scripts. See
[Distributed Lock](../infrastructure/lock).

## Monitor

- **Breaking**: the overview disk summary reports the root filesystem only;
  `vef.monitor.excluded_mounts` was removed.
- CPU summaries/info gained `effectiveCores` (cgroup v1/v2 quota-aware
  utilization normalization); memory headline figures respect cgroup limits;
  the CPU total is derived from the per-core sample and the sampler lifecycle
  was hardened.

## Removals and Smaller Fixes

- **`ptr` package removed** — migration table in
  [Small Helpers](../utilities/small-helpers). (`ptr.Of` briefly changed to
  always return a pointer before the package was dropped entirely.)
- **Snowflake generator removed** (`id.DefaultSnowflakeIDGenerator`,
  `NewSnowflakeIDGenerator`, `config.EnvNodeID` / `VEF_NODE_ID`); `id.Generate()`
  (XID) is the framework's ID mechanism.
- `cron.Scheduler.Update` preserves the job identifier, so the handle you hold
  stays valid across updates.
- `sequence` rejects `SeqStep < 1` with `sequence.ErrInvalidStep`, and
  `MemoryStore.Register` serializes with in-flight `Reserve` calls.
- `timex` lenient `Parse*` fallback now parses zone-less input in local time,
  matching the primary path.
- `cache.NewInvalidating` forwards `MemoryOption`s (bound read-through caches
  with `WithMemMaxSize`); the mold cached dictionary resolver caps its LRU at
  4096 entries.
- `schema.UniqueKey` gained `predicate` / `hasExpressions`; `schema.ErrTableMissing`
  is the new Go-level sentinel for missing tables.
- `expression` passes JSON-native environments through without re-marshaling
  (performance only).
- `event` documentation clarified `Ordered` semantics: order is delivery
  order; `WithConcurrency > 1` interleaves handler executions.
