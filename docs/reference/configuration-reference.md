---
sidebar_position: 1
---

# Configuration Reference

This page summarizes the config structs currently exposed by the framework and the runtime defaults applied by built-in modules.

Minimal starter block:

```toml
[vef.app]
name = "my-app"
port = 8080

[vef.data_sources.primary]
type = "sqlite"
```

## File Lookup Order

The internal config module searches for `application.toml` in this order:

- `./configs`
- `$VEF_CONFIG_PATH`
- `.`
- `../configs`

## Common Environment Variables

Common environment keys include:

- `VEF_CONFIG_PATH`
- `VEF_LOG_LEVEL`
- `VEF_I18N_LANGUAGE`

## `vef.app`

| Field | Type | Meaning |
| --- | --- | --- |
| `name` | `string` | application name; also feeds defaults such as JWT audience generation |
| `port` | `uint16` | HTTP server port |
| `body_limit` | `string` | Fiber body limit, for example `10mib`; defaults to `32mib` when omitted |
| `trusted_proxies` | `[]string` | proxy IPs or CIDR ranges trusted to set `X-Forwarded-For`; empty means forwarded headers from untrusted clients are ignored |

## `vef.data_sources`

`vef.data_sources` is a map keyed by data source name. The `primary` entry is
required and powers the framework-wide `orm.DB` injection; other entries are
registered into the data source registry under their map key.

Example:

```toml
[vef.data_sources.primary]
type = "sqlite"

[vef.data_sources.analytics]
type = "sqlite"
path = "./analytics.db"
```

| Field | Type | Meaning |
| --- | --- | --- |
| `type` | `postgres \| mysql \| sqlite` | runtime-supported database kind; `oracle` and `sqlserver` constants exist but are not implemented yet |
| `host` | `string` | network database host |
| `port` | `uint16` | network database port |
| `user` | `string` | database username |
| `password` | `string` | database password |
| `database` | `string` | database name |
| `schema` | `string` | schema name for drivers that support schemas |
| `path` | `string` | SQLite file path |
| `enable_sql_guard` | `bool` | enables the SQL guard for raw SQL surfaces |
| `ssl_mode` | `disable \| require \| verify-ca \| verify-full` | TLS posture for network database dialects; omitted means `disable` |
| `ssl_root_cert` | `string` | optional PEM CA bundle path for `verify-ca` and `verify-full`; empty uses the host system pool |

Runtime note:

- the current runtime provider registry supports `postgres`, `mysql`, and `sqlite`. `oracle` and `sqlserver` are declared as `DBKind` constants for future use but have no runtime provider yet, so configuring them fails at startup with `database.ErrUnsupportedDBKind`.

## `vef.cors`

| Field | Type | Meaning |
| --- | --- | --- |
| `enabled` | `bool` | enables the CORS middleware |
| `allow_origins` | `[]string` | allowed origin list |

## `vef.security`

| Field | Type | Meaning |
| --- | --- | --- |
| `secret` | `string` | hex-encoded JWT signing key. If unset, the framework generates an ephemeral per-process key and warns; tokens do not survive restart or work across nodes. If set to the public `security.DefaultJWTSecret`, startup warns to replace it in production. |
| `token_expires` | `duration` | refresh-token lifetime; default `168h` |
| `refresh_not_before` | `duration` | earliest time a refresh token may be used; default `15m`, half of the fixed `30m` access-token lifetime |
| `login_rate_limit` | `int` | login endpoint rate limit; default `6` |
| `refresh_rate_limit` | `int` | refresh endpoint rate limit; default `1` |
| `ip_whitelists` | `map[string][]string` | named source-IP whitelists (IP or CIDR entries) consumed by the built-in `ip` auth strategy; TOML keys are lowercased, and the no-arg `api.IPAuth()` targets the `default` key |
| `lockout.*` | — | brute-force lockout on the login endpoint: `enabled` default `true`, `max_failures` default `10`, `window` default `15m`, `lock_duration` default `15m`, `strategy` (`lock` \| `backoff`) default `lock`, `backoff_base` default `1s`, `backoff_max` default `15m`, `key` (`user` \| `ip` \| `user_ip`) default `user_ip` |
| `password_policy.*` | — | password strength rules; every field is opt-in (a zero value disables the rule): `min_length`, `max_length`, `require_upper`, `require_lower`, `require_digit`, `require_symbol`, `min_char_classes`, `disallow_username`, `blocklist`, `history_depth` (reuse prevention; requires an app-provided `security.PasswordHistoryStore`), `max_age` (expiry; requires an app-provided `security.PasswordMetadataLoader`) |
| `token_type` | `jwt_token \| opaque_token` | login token mechanism; default `jwt_token`. Session control (concurrency limits, force-offline, renewal) is only available with `opaque_token` |
| `session.*` | — | opaque-token session tuning, no effect under `jwt_token`: `max_concurrent` default `0` (unlimited; enforcement is best-effort under concurrent logins), `on_exceed` (`reject` \| `evict_oldest`) default `evict_oldest`, `idle_ttl` default `30m`, `max_lifetime` default `168h` (7 days), `sliding` default `true` |

Runtime note:

- access tokens issued by the built-in JWT token generator expire after `30m`; `vef.security.token_expires` controls refresh tokens, not access tokens
- lockout is on by default (`max_failures = 10`); a trip returns `security.ErrAccountLocked` (HTTP 429) and guard-store errors fail open
- `history_depth > 0` composes a history validator into the password policy only when a `security.PasswordHistoryStore` is registered; `max_age` only takes effect when the app wires a `security.PasswordMetadataLoader` and `security.NewExpiryPasswordChangeChecker`

## `vef.redis`

| Field | Type | Meaning |
| --- | --- | --- |
| `enabled` | `bool` | constructs the Redis client when true; default `false` |
| `host` | `string` | Redis host |
| `port` | `uint16` | Redis port |
| `user` | `string` | Redis username |
| `password` | `string` | Redis password |
| `database` | `uint8` | Redis database number |
| `network` | `string` | `tcp` or `unix` |

Runtime note:

- the default `vef.Run(...)` boot graph includes the Redis module
- the Redis client is constructed only when `enabled = true`; when `enabled` is false or omitted, the framework provides a nil `*redis.Client` and skips startup `PING`
- when enabled, omitted host/port/network fields default to `127.0.0.1`, `6379`, and `tcp`

## `vef.storage`

| Field | Type | Meaning |
| --- | --- | --- |
| `provider` | `memory \| minio \| filesystem` | storage provider selection |
| `auto_migrate` | `bool` | runs storage DDL migration at startup |
| `minio.endpoint` | `string` | MinIO endpoint |
| `minio.access_key` | `string` | MinIO access key |
| `minio.secret_key` | `string` | MinIO secret key |
| `minio.bucket` | `string` | bucket name |
| `minio.region` | `string` | region |
| `minio.use_ssl` | `bool` | whether to use HTTPS |
| `filesystem.root` | `string` | filesystem provider root directory |
| `max_upload_size` | `int64` | maximum single-object upload size, default 1 GiB |
| `claim_ttl` | `duration` | upload claim lifetime, default 24h |
| `max_pending_claims` | `int` | maximum simultaneous pending claims per principal, default 100 |
| `allow_public_uploads` | `bool` | allows clients to request public uploads; default false |
| `sweep_interval` | `duration` | expired-claim sweep interval, default 5m |
| `sweep_batch_size` | `int` | maximum expired claims processed per sweep, default 200 |
| `delete_worker_interval` | `duration` | pending-delete worker polling interval, default 5m |
| `delete_batch_size` | `int` | rows leased by one delete-worker tick, default 100 |
| `delete_concurrency` | `int` | concurrent object deletions per worker tick, default 8 |
| `delete_max_attempts` | `int` | retry budget before dead-lettering a delete row, default 12 |
| `delete_lease_window` | `duration` | delete-row lease visibility window, default 5m |

Runtime note:

- omitting `provider` selects in-memory storage and logs a warning; objects are lost on restart
- `vef.storage.auto_migrate = true` runs the idempotent storage migration and checks `sys_storage_upload_claim`, `sys_storage_upload_part`, and `sys_storage_pending_delete`
- `filesystem.root` defaults to `./storage`
- `minio.bucket` defaults to `minio.bucket`, then `vef.app.name`, then `vef-app`
- upload-flow and delete-worker tunables have defaults in the framework; use the `StorageConfig` `Effective...` accessors when application code needs the resolved values

## `vef.monitor`

| Field | Type | Meaning |
| --- | --- | --- |
| `sample_interval` | `duration` | interval between samples; default `10s` |
| `sample_duration` | `duration` | sampling window duration; default `2s` |
| `excluded_mounts` | `[]string` | additional mount-point substrings to exclude from disk statistics |

## `vef.mcp`

| Field | Type | Meaning |
| --- | --- | --- |
| `enabled` | `bool` | enables the MCP server and `/mcp` endpoint |
| `require_auth` | `bool` | secure by default: unset or `true` requires Bearer auth; only explicit `false` allows anonymous access. In Go, `MCPConfig.RequireAuth` is `*bool` so the runtime can distinguish unset from false. |

## `vef.approval`

| Field | Type | Meaning |
| --- | --- | --- |
| `auto_migrate` | `bool` | runs approval DDL migration at startup when explicitly enabled; `ApprovalConfig.ApplyDefaults()` does not turn it on |
| `timeout_scan_interval` | `duration` | timeout scanner cadence, default 1m |
| `pre_warning_scan_interval` | `duration` | pre-warning scanner cadence, default 5m |
| `cleanup_scan_interval` | `duration` | retention cleanup cadence, default 24h |
| `delegation_max_depth` | `int` | maximum delegation-chain depth, default 10 |
| `form_snapshot_retention` | `duration` | `apv_form_snapshot` retention, default 90 days |
| `urge_record_retention` | `duration` | `apv_urge_record` retention, default 30 days |
| `cc_record_retention` | `duration` | retention for read `apv_cc_record` rows, default 90 days |

> Outbox-related fields moved to `[vef.event.transports.outbox]` in v0.21; see [Event Bus](../infrastructure/event-bus).

## `vef.event`

| Field | Type | Default / meaning |
| --- | --- | --- |
| `default_transport` | `string` | route fallback, default `memory` |
| `async_queue_size` | `int` | `WithAsync` queue capacity, default `4096` |
| `async_workers` | `int` | async worker count, default `4` |
| `publish_timeout` | `duration` | per-transport publish timeout, default `5s` |
| `transports.memory.*` | — | `queue_size` default `1024`, `full_policy` default `error`, `publish_timeout` default unset/no timeout and only applies when `full_policy = "block"` |
| `transports.outbox.*` | — | `enabled`, `relay_interval` default `10s`, `max_retries` default `10`, `batch_size` default `100`, `lease_multiplier` default `4`, `min_lease` default `15s`, `sink` default `memory`, `cleanup_interval` default `1h`, `completed_ttl` default `168h`; cleanup fields belong to framework config, not `event/transport/outbox.Config` |
| `transports.redis_stream.*` | — | `enabled`, `stream_prefix` default `vef:events:`, `max_len_approx` default `0` (no trimming), `block_timeout` default `5s`, `claim_idle` default `60s`, `claim_interval` default `30s`, `claim_batch_size` default `64`, `reaper_concurrency` default `4`, `handler_timeout` default `30s`, `setup_timeout` default `5s`, `consumer_id` default prefix `vef`, `start_id` default `0` (`"$"` skips backlog for newly created groups), `idle_group_retention` default `0` (disables orphaned consumer-group reclamation), `idle_group_sweep_interval` default `10m` |
| `middleware.*` | `bool` | middleware toggles: `logging`, `tracing`, `tracing_strict`, `metrics`, `recover`, `inbox` |
| `inbox.*` | — | `retention` default `168h`, `processing_lease` default `10m`, `cleanup_interval` default `1h` |
| `routing` | `[]{pattern, transports}` | routing rules, matched top-to-bottom with `path.Match` |

## Config Package API Reference

### Top-Level Public Symbols

| Symbol | Kind | Signature or value |
| --- | --- | --- |
| `config.AppConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.AppConfig` |
| `config.ApprovalConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.ApprovalConfig` |
| `config.Config` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.Config` |
| `config.CORSConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.CORSConfig` |
| `config.DBKind` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.DBKind` |
| `config.DataSourceConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.DataSourceConfig` |
| `config.DataSourcesConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.DataSourcesConfig` |
| `config.DefaultClaimTTL` | `CONST` | `time.Duration = 86400000000000` |
| `config.DefaultDeleteBatchSize` | `CONST` | `int = 100` |
| `config.DefaultDeleteConcurrency` | `CONST` | `int = 8` |
| `config.DefaultDeleteLeaseWindow` | `CONST` | `time.Duration = 300000000000` |
| `config.DefaultDeleteMaxAttempts` | `CONST` | `int = 12` |
| `config.DefaultDeleteWorkerInterval` | `CONST` | `time.Duration = 300000000000` |
| `config.DefaultLockoutBackoffBase` | `CONST` | `time.Duration = 1000000000` |
| `config.DefaultLockoutBackoffMax` | `CONST` | `time.Duration = 900000000000` |
| `config.DefaultLockoutLockDuration` | `CONST` | `time.Duration = 900000000000` |
| `config.DefaultLockoutMaxFailures` | `CONST` | `int = 10` |
| `config.DefaultLockoutWindow` | `CONST` | `time.Duration = 900000000000` |
| `config.DefaultMaxPendingClaims` | `CONST` | `int = 100` |
| `config.DefaultMaxUploadSize` | `CONST` | `int64 = 1073741824` |
| `config.DefaultSessionIdleTTL` | `CONST` | `time.Duration = 1800000000000` |
| `config.DefaultSessionMaxLifetime` | `CONST` | `time.Duration = 604800000000000` |
| `config.DefaultSweepBatchSize` | `CONST` | `int = 200` |
| `config.DefaultSweepInterval` | `CONST` | `time.Duration = 300000000000` |
| `config.EnvConfigPath` | `CONST` | `untyped string = "VEF_CONFIG_PATH"` |
| `config.EnvI18NLanguage` | `CONST` | `untyped string = "VEF_I18N_LANGUAGE"` |
| `config.EnvPrefix` | `CONST` | `untyped string = "VEF"` |
| `config.EnvLogLevel` | `CONST` | `untyped string = "VEF_LOG_LEVEL"` |
| `config.ErrInboxRetentionTooShort` | `VAR` | `error` |
| `config.ErrInvalidLockoutKey` | `VAR` | `error` |
| `config.ErrInvalidLockoutStrategy` | `VAR` | `error` |
| `config.ErrInvalidSessionOnExceed` | `VAR` | `error` |
| `config.ErrInvalidTokenType` | `VAR` | `error` |
| `config.EventConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.EventConfig` |
| `config.EventInboxConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.EventInboxConfig` |
| `config.EventMemoryTransportConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.EventMemoryTransportConfig` |
| `config.EventMiddlewareConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.EventMiddlewareConfig` |
| `config.EventOutboxTransportConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.EventOutboxTransportConfig` |
| `config.EventRedisStreamTransportConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.EventRedisStreamTransportConfig` |
| `config.EventRoutingRule` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.EventRoutingRule` |
| `config.EventTransportsConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.EventTransportsConfig` |
| `config.FilesystemConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.FilesystemConfig` |
| `config.LockoutConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.LockoutConfig` |
| `config.LockoutKey` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.LockoutKey` |
| `config.LockoutKeyIP` | `CONST` | `github.com/coldsmirk/vef-framework-go/config.LockoutKey = "ip"` |
| `config.LockoutKeyUser` | `CONST` | `github.com/coldsmirk/vef-framework-go/config.LockoutKey = "user"` |
| `config.LockoutKeyUserIP` | `CONST` | `github.com/coldsmirk/vef-framework-go/config.LockoutKey = "user_ip"` |
| `config.LockoutStrategy` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.LockoutStrategy` |
| `config.LockoutStrategyBackoff` | `CONST` | `github.com/coldsmirk/vef-framework-go/config.LockoutStrategy = "backoff"` |
| `config.LockoutStrategyLock` | `CONST` | `github.com/coldsmirk/vef-framework-go/config.LockoutStrategy = "lock"` |
| `config.MCPConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.MCPConfig` |
| `config.MinIOConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.MinIOConfig` |
| `config.MonitorConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.MonitorConfig` |
| `config.MySQL` | `CONST` | `github.com/coldsmirk/vef-framework-go/config.DBKind = "mysql"` |
| `config.Oracle` | `CONST` | `github.com/coldsmirk/vef-framework-go/config.DBKind = "oracle"` |
| `config.PasswordPolicyConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.PasswordPolicyConfig` |
| `config.Postgres` | `CONST` | `github.com/coldsmirk/vef-framework-go/config.DBKind = "postgres"` |
| `config.PrimaryDataSourceName` | `CONST` | `untyped string = "primary"` |
| `config.RedisConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.RedisConfig` |
| `config.SQLServer` | `CONST` | `github.com/coldsmirk/vef-framework-go/config.DBKind = "sqlserver"` |
| `config.SQLite` | `CONST` | `github.com/coldsmirk/vef-framework-go/config.DBKind = "sqlite"` |
| `config.SecurityConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.SecurityConfig` |
| `config.SSLDisable` | `CONST` | `github.com/coldsmirk/vef-framework-go/config.SSLMode = "disable"` |
| `config.SSLMode` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.SSLMode` |
| `config.SSLRequire` | `CONST` | `github.com/coldsmirk/vef-framework-go/config.SSLMode = "require"` |
| `config.SSLVerifyCA` | `CONST` | `github.com/coldsmirk/vef-framework-go/config.SSLMode = "verify-ca"` |
| `config.SSLVerifyFull` | `CONST` | `github.com/coldsmirk/vef-framework-go/config.SSLMode = "verify-full"` |
| `config.SessionConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.SessionConfig` |
| `config.SessionExceedEvictOldest` | `CONST` | `github.com/coldsmirk/vef-framework-go/config.SessionExceedPolicy = "evict_oldest"` |
| `config.SessionExceedPolicy` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.SessionExceedPolicy` |
| `config.SessionExceedReject` | `CONST` | `github.com/coldsmirk/vef-framework-go/config.SessionExceedPolicy = "reject"` |
| `config.StorageConfig` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.StorageConfig` |
| `config.StorageFilesystem` | `CONST` | `github.com/coldsmirk/vef-framework-go/config.StorageProvider = "filesystem"` |
| `config.StorageMemory` | `CONST` | `github.com/coldsmirk/vef-framework-go/config.StorageProvider = "memory"` |
| `config.StorageMinIO` | `CONST` | `github.com/coldsmirk/vef-framework-go/config.StorageProvider = "minio"` |
| `config.StorageProvider` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.StorageProvider` |
| `config.TokenType` | `TYPE` | `github.com/coldsmirk/vef-framework-go/config.TokenType` |
| `config.TokenTypeJWT` | `CONST` | `github.com/coldsmirk/vef-framework-go/config.TokenType = "jwt_token"` |
| `config.TokenTypeOpaque` | `CONST` | `github.com/coldsmirk/vef-framework-go/config.TokenType = "opaque_token"` |

### Exported Fields

| Field | Signature and config tag |
| --- | --- |
| `config.AppConfig.Name` | `string [field_order=1 tag="config:\"name\""]` |
| `config.AppConfig.Port` | `uint16 [field_order=2 tag="config:\"port\""]` |
| `config.AppConfig.BodyLimit` | `string [field_order=3 tag="config:\"body_limit\""]` |
| `config.AppConfig.TrustedProxies` | `[]string [field_order=4 tag="config:\"trusted_proxies\""]` |
| `config.ApprovalConfig.AutoMigrate` | `bool [field_order=1 tag="config:\"auto_migrate\""]` |
| `config.ApprovalConfig.TimeoutScanInterval` | `time.Duration [field_order=2 tag="config:\"timeout_scan_interval\""]` |
| `config.ApprovalConfig.PreWarningScanInterval` | `time.Duration [field_order=3 tag="config:\"pre_warning_scan_interval\""]` |
| `config.ApprovalConfig.CleanupScanInterval` | `time.Duration [field_order=4 tag="config:\"cleanup_scan_interval\""]` |
| `config.ApprovalConfig.DelegationMaxDepth` | `int [field_order=5 tag="config:\"delegation_max_depth\""]` |
| `config.ApprovalConfig.FormSnapshotRetention` | `time.Duration [field_order=6 tag="config:\"form_snapshot_retention\""]` |
| `config.ApprovalConfig.UrgeRecordRetention` | `time.Duration [field_order=7 tag="config:\"urge_record_retention\""]` |
| `config.ApprovalConfig.CCRecordRetention` | `time.Duration [field_order=8 tag="config:\"cc_record_retention\""]` |
| `config.CORSConfig.Enabled` | `bool [field_order=1 tag="config:\"enabled\""]` |
| `config.CORSConfig.AllowOrigins` | `[]string [field_order=2 tag="config:\"allow_origins\""]` |
| `config.DataSourceConfig.Kind` | `github.com/coldsmirk/vef-framework-go/config.DBKind [field_order=1 tag="config:\"type\""]` |
| `config.DataSourceConfig.Host` | `string [field_order=2 tag="config:\"host\""]` |
| `config.DataSourceConfig.Port` | `uint16 [field_order=3 tag="config:\"port\""]` |
| `config.DataSourceConfig.User` | `string [field_order=4 tag="config:\"user\""]` |
| `config.DataSourceConfig.Password` | `string [field_order=5 tag="config:\"password\""]` |
| `config.DataSourceConfig.Database` | `string [field_order=6 tag="config:\"database\""]` |
| `config.DataSourceConfig.Schema` | `string [field_order=7 tag="config:\"schema\""]` |
| `config.DataSourceConfig.Path` | `string [field_order=8 tag="config:\"path\""]` |
| `config.DataSourceConfig.EnableSQLGuard` | `bool [field_order=9 tag="config:\"enable_sql_guard\""]` |
| `config.DataSourceConfig.SSLMode` | `github.com/coldsmirk/vef-framework-go/config.SSLMode [field_order=10 tag="config:\"ssl_mode\""]` |
| `config.DataSourceConfig.SSLRootCert` | `string [field_order=11 tag="config:\"ssl_root_cert\""]` |
| `config.DataSourcesConfig.Map` | `map[string]github.com/coldsmirk/vef-framework-go/config.DataSourceConfig [field_order=1 tag=""]` |
| `config.EventConfig.DefaultTransport` | `string [field_order=1 tag="config:\"default_transport\""]` |
| `config.EventConfig.AsyncQueueSize` | `int [field_order=2 tag="config:\"async_queue_size\""]` |
| `config.EventConfig.AsyncWorkers` | `int [field_order=3 tag="config:\"async_workers\""]` |
| `config.EventConfig.PublishTimeout` | `time.Duration [field_order=4 tag="config:\"publish_timeout\""]` |
| `config.EventConfig.Transports` | `github.com/coldsmirk/vef-framework-go/config.EventTransportsConfig [field_order=5 tag="config:\"transports\""]` |
| `config.EventConfig.Middleware` | `github.com/coldsmirk/vef-framework-go/config.EventMiddlewareConfig [field_order=6 tag="config:\"middleware\""]` |
| `config.EventConfig.Inbox` | `github.com/coldsmirk/vef-framework-go/config.EventInboxConfig [field_order=7 tag="config:\"inbox\""]` |
| `config.EventConfig.Routing` | `[]github.com/coldsmirk/vef-framework-go/config.EventRoutingRule [field_order=8 tag="config:\"routing\""]` |
| `config.EventInboxConfig.Retention` | `time.Duration [field_order=1 tag="config:\"retention\""]` |
| `config.EventInboxConfig.ProcessingLease` | `time.Duration [field_order=2 tag="config:\"processing_lease\""]` |
| `config.EventInboxConfig.CleanupInterval` | `time.Duration [field_order=3 tag="config:\"cleanup_interval\""]` |
| `config.EventMemoryTransportConfig.QueueSize` | `int [field_order=1 tag="config:\"queue_size\""]` |
| `config.EventMemoryTransportConfig.FullPolicy` | `string [field_order=2 tag="config:\"full_policy\""]` |
| `config.EventMemoryTransportConfig.PublishTimeout` | `time.Duration [field_order=3 tag="config:\"publish_timeout\""]` |
| `config.EventMiddlewareConfig.Logging` | `bool [field_order=1 tag="config:\"logging\""]` |
| `config.EventMiddlewareConfig.Tracing` | `bool [field_order=2 tag="config:\"tracing\""]` |
| `config.EventMiddlewareConfig.TracingStrict` | `bool [field_order=3 tag="config:\"tracing_strict\""]` |
| `config.EventMiddlewareConfig.Metrics` | `bool [field_order=4 tag="config:\"metrics\""]` |
| `config.EventMiddlewareConfig.Recover` | `bool [field_order=5 tag="config:\"recover\""]` |
| `config.EventMiddlewareConfig.Inbox` | `bool [field_order=6 tag="config:\"inbox\""]` |
| `config.EventOutboxTransportConfig.Enabled` | `bool [field_order=1 tag="config:\"enabled\""]` |
| `config.EventOutboxTransportConfig.RelayInterval` | `time.Duration [field_order=2 tag="config:\"relay_interval\""]` |
| `config.EventOutboxTransportConfig.MaxRetries` | `int [field_order=3 tag="config:\"max_retries\""]` |
| `config.EventOutboxTransportConfig.BatchSize` | `int [field_order=4 tag="config:\"batch_size\""]` |
| `config.EventOutboxTransportConfig.LeaseMultiplier` | `int [field_order=5 tag="config:\"lease_multiplier\""]` |
| `config.EventOutboxTransportConfig.MinLease` | `time.Duration [field_order=6 tag="config:\"min_lease\""]` |
| `config.EventOutboxTransportConfig.SinkName` | `string [field_order=7 tag="config:\"sink\""]` |
| `config.EventOutboxTransportConfig.CleanupInterval` | `time.Duration [field_order=8 tag="config:\"cleanup_interval\""]` |
| `config.EventOutboxTransportConfig.CompletedTTL` | `time.Duration [field_order=9 tag="config:\"completed_ttl\""]` |
| `config.EventRedisStreamTransportConfig.Enabled` | `bool [field_order=1 tag="config:\"enabled\""]` |
| `config.EventRedisStreamTransportConfig.StreamPrefix` | `string [field_order=2 tag="config:\"stream_prefix\""]` |
| `config.EventRedisStreamTransportConfig.MaxLenApprox` | `int64 [field_order=3 tag="config:\"max_len_approx\""]` |
| `config.EventRedisStreamTransportConfig.BlockTimeout` | `time.Duration [field_order=4 tag="config:\"block_timeout\""]` |
| `config.EventRedisStreamTransportConfig.ClaimIdle` | `time.Duration [field_order=5 tag="config:\"claim_idle\""]` |
| `config.EventRedisStreamTransportConfig.ClaimInterval` | `time.Duration [field_order=6 tag="config:\"claim_interval\""]` |
| `config.EventRedisStreamTransportConfig.ClaimBatchSize` | `int64 [field_order=7 tag="config:\"claim_batch_size\""]` |
| `config.EventRedisStreamTransportConfig.ReaperConcurrency` | `int [field_order=8 tag="config:\"reaper_concurrency\""]` |
| `config.EventRedisStreamTransportConfig.HandlerTimeout` | `time.Duration [field_order=9 tag="config:\"handler_timeout\""]` |
| `config.EventRedisStreamTransportConfig.SetupTimeout` | `time.Duration [field_order=10 tag="config:\"setup_timeout\""]` |
| `config.EventRedisStreamTransportConfig.ConsumerID` | `string [field_order=11 tag="config:\"consumer_id\""]` |
| `config.EventRedisStreamTransportConfig.StartID` | `string [field_order=12 tag="config:\"start_id\""]` |
| `config.EventRedisStreamTransportConfig.IdleGroupRetention` | `time.Duration [field_order=13 tag="config:\"idle_group_retention\""]` |
| `config.EventRedisStreamTransportConfig.IdleGroupSweepInterval` | `time.Duration [field_order=14 tag="config:\"idle_group_sweep_interval\""]` |
| `config.EventRoutingRule.Pattern` | `string [field_order=1 tag="config:\"pattern\""]` |
| `config.EventRoutingRule.Transports` | `[]string [field_order=2 tag="config:\"transports\""]` |
| `config.EventTransportsConfig.Memory` | `github.com/coldsmirk/vef-framework-go/config.EventMemoryTransportConfig [field_order=1 tag="config:\"memory\""]` |
| `config.EventTransportsConfig.Outbox` | `github.com/coldsmirk/vef-framework-go/config.EventOutboxTransportConfig [field_order=2 tag="config:\"outbox\""]` |
| `config.EventTransportsConfig.RedisStream` | `github.com/coldsmirk/vef-framework-go/config.EventRedisStreamTransportConfig [field_order=3 tag="config:\"redis_stream\""]` |
| `config.FilesystemConfig.Root` | `string [field_order=1 tag="config:\"root\""]` |
| `config.LockoutConfig.Enabled` | `*bool [field_order=1 tag="config:\"enabled\""]` |
| `config.LockoutConfig.MaxFailures` | `int [field_order=2 tag="config:\"max_failures\""]` |
| `config.LockoutConfig.Window` | `time.Duration [field_order=3 tag="config:\"window\""]` |
| `config.LockoutConfig.LockDuration` | `time.Duration [field_order=4 tag="config:\"lock_duration\""]` |
| `config.LockoutConfig.Strategy` | `github.com/coldsmirk/vef-framework-go/config.LockoutStrategy [field_order=5 tag="config:\"strategy\""]` |
| `config.LockoutConfig.BackoffBase` | `time.Duration [field_order=6 tag="config:\"backoff_base\""]` |
| `config.LockoutConfig.BackoffMax` | `time.Duration [field_order=7 tag="config:\"backoff_max\""]` |
| `config.LockoutConfig.Key` | `github.com/coldsmirk/vef-framework-go/config.LockoutKey [field_order=8 tag="config:\"key\""]` |
| `config.MCPConfig.Enabled` | `bool [field_order=1 tag="config:\"enabled\""]` |
| `config.MCPConfig.RequireAuth` | `*bool [field_order=2 tag="config:\"require_auth\""]` |
| `config.MinIOConfig.Endpoint` | `string [field_order=1 tag="config:\"endpoint\""]` |
| `config.MinIOConfig.AccessKey` | `string [field_order=2 tag="config:\"access_key\""]` |
| `config.MinIOConfig.SecretKey` | `string [field_order=3 tag="config:\"secret_key\""]` |
| `config.MinIOConfig.Bucket` | `string [field_order=4 tag="config:\"bucket\""]` |
| `config.MinIOConfig.Region` | `string [field_order=5 tag="config:\"region\""]` |
| `config.MinIOConfig.UseSSL` | `bool [field_order=6 tag="config:\"use_ssl\""]` |
| `config.MonitorConfig.SampleInterval` | `time.Duration [field_order=1 tag="config:\"sample_interval\""]` |
| `config.MonitorConfig.SampleDuration` | `time.Duration [field_order=2 tag="config:\"sample_duration\""]` |
| `config.MonitorConfig.ExcludedMounts` | `[]string [field_order=3 tag="config:\"excluded_mounts\""]` |
| `config.PasswordPolicyConfig.MinLength` | `int [field_order=1 tag="config:\"min_length\""]` |
| `config.PasswordPolicyConfig.MaxLength` | `int [field_order=2 tag="config:\"max_length\""]` |
| `config.PasswordPolicyConfig.RequireUpper` | `bool [field_order=3 tag="config:\"require_upper\""]` |
| `config.PasswordPolicyConfig.RequireLower` | `bool [field_order=4 tag="config:\"require_lower\""]` |
| `config.PasswordPolicyConfig.RequireDigit` | `bool [field_order=5 tag="config:\"require_digit\""]` |
| `config.PasswordPolicyConfig.RequireSymbol` | `bool [field_order=6 tag="config:\"require_symbol\""]` |
| `config.PasswordPolicyConfig.MinCharClasses` | `int [field_order=7 tag="config:\"min_char_classes\""]` |
| `config.PasswordPolicyConfig.DisallowUsername` | `bool [field_order=8 tag="config:\"disallow_username\""]` |
| `config.PasswordPolicyConfig.Blocklist` | `[]string [field_order=9 tag="config:\"blocklist\""]` |
| `config.PasswordPolicyConfig.HistoryDepth` | `int [field_order=10 tag="config:\"history_depth\""]` |
| `config.PasswordPolicyConfig.MaxAge` | `time.Duration [field_order=11 tag="config:\"max_age\""]` |
| `config.RedisConfig.Enabled` | `bool [field_order=1 tag="config:\"enabled\""]` |
| `config.RedisConfig.Host` | `string [field_order=2 tag="config:\"host\""]` |
| `config.RedisConfig.Port` | `uint16 [field_order=3 tag="config:\"port\""]` |
| `config.RedisConfig.User` | `string [field_order=4 tag="config:\"user\""]` |
| `config.RedisConfig.Password` | `string [field_order=5 tag="config:\"password\""]` |
| `config.RedisConfig.Database` | `uint8 [field_order=6 tag="config:\"database\""]` |
| `config.RedisConfig.Network` | `string [field_order=7 tag="config:\"network\""]` |
| `config.SecurityConfig.Secret` | `string [field_order=1 tag="config:\"secret\""]` |
| `config.SecurityConfig.TokenExpires` | `time.Duration [field_order=2 tag="config:\"token_expires\""]` |
| `config.SecurityConfig.RefreshNotBefore` | `time.Duration [field_order=3 tag="config:\"refresh_not_before\""]` |
| `config.SecurityConfig.LoginRateLimit` | `int [field_order=4 tag="config:\"login_rate_limit\""]` |
| `config.SecurityConfig.RefreshRateLimit` | `int [field_order=5 tag="config:\"refresh_rate_limit\""]` |
| `config.SecurityConfig.IPWhitelists` | `map[string][]string [field_order=6 tag="config:\"ip_whitelists\""]` |
| `config.SecurityConfig.Lockout` | `github.com/coldsmirk/vef-framework-go/config.LockoutConfig [field_order=7 tag="config:\"lockout\""]` |
| `config.SecurityConfig.PasswordPolicy` | `github.com/coldsmirk/vef-framework-go/config.PasswordPolicyConfig [field_order=8 tag="config:\"password_policy\""]` |
| `config.SecurityConfig.TokenType` | `github.com/coldsmirk/vef-framework-go/config.TokenType [field_order=9 tag="config:\"token_type\""]` |
| `config.SecurityConfig.Session` | `github.com/coldsmirk/vef-framework-go/config.SessionConfig [field_order=10 tag="config:\"session\""]` |
| `config.SessionConfig.MaxConcurrent` | `int [field_order=1 tag="config:\"max_concurrent\""]` |
| `config.SessionConfig.OnExceed` | `github.com/coldsmirk/vef-framework-go/config.SessionExceedPolicy [field_order=2 tag="config:\"on_exceed\""]` |
| `config.SessionConfig.IdleTTL` | `time.Duration [field_order=3 tag="config:\"idle_ttl\""]` |
| `config.SessionConfig.MaxLifetime` | `time.Duration [field_order=4 tag="config:\"max_lifetime\""]` |
| `config.SessionConfig.Sliding` | `*bool [field_order=5 tag="config:\"sliding\""]` |
| `config.StorageConfig.Provider` | `github.com/coldsmirk/vef-framework-go/config.StorageProvider [field_order=1 tag="config:\"provider\""]` |
| `config.StorageConfig.AutoMigrate` | `bool [field_order=2 tag="config:\"auto_migrate\""]` |
| `config.StorageConfig.MinIO` | `github.com/coldsmirk/vef-framework-go/config.MinIOConfig [field_order=3 tag="config:\"minio\""]` |
| `config.StorageConfig.Filesystem` | `github.com/coldsmirk/vef-framework-go/config.FilesystemConfig [field_order=4 tag="config:\"filesystem\""]` |
| `config.StorageConfig.MaxUploadSize` | `int64 [field_order=5 tag="config:\"max_upload_size\""]` |
| `config.StorageConfig.ClaimTTL` | `time.Duration [field_order=6 tag="config:\"claim_ttl\""]` |
| `config.StorageConfig.MaxPendingClaims` | `int [field_order=7 tag="config:\"max_pending_claims\""]` |
| `config.StorageConfig.AllowPublicUploads` | `bool [field_order=8 tag="config:\"allow_public_uploads\""]` |
| `config.StorageConfig.SweepInterval` | `time.Duration [field_order=9 tag="config:\"sweep_interval\""]` |
| `config.StorageConfig.SweepBatchSize` | `int [field_order=10 tag="config:\"sweep_batch_size\""]` |
| `config.StorageConfig.DeleteWorkerInterval` | `time.Duration [field_order=11 tag="config:\"delete_worker_interval\""]` |
| `config.StorageConfig.DeleteBatchSize` | `int [field_order=12 tag="config:\"delete_batch_size\""]` |
| `config.StorageConfig.DeleteConcurrency` | `int [field_order=13 tag="config:\"delete_concurrency\""]` |
| `config.StorageConfig.DeleteMaxAttempts` | `int [field_order=14 tag="config:\"delete_max_attempts\""]` |
| `config.StorageConfig.DeleteLeaseWindow` | `time.Duration [field_order=15 tag="config:\"delete_lease_window\""]` |

### Exported Methods

| Method | Signature |
| --- | --- |
| `config.ApprovalConfig.ApplyDefaults` | `func()` |
| `config.Config.Unmarshal` | `func(key string, target any) error` |
| `config.DataSourcesConfig.Primary` | `func() github.com/coldsmirk/vef-framework-go/config.DataSourceConfig` |
| `config.EventConfig.EffectiveAsyncQueueSize` | `func() int` |
| `config.EventConfig.EffectiveAsyncWorkers` | `func() int` |
| `config.EventConfig.EffectiveDefaultTransport` | `func() string` |
| `config.EventConfig.EffectivePublishTimeout` | `func() time.Duration` |
| `config.EventConfig.Validate` | `func() error` |
| `config.EventInboxConfig.EffectiveCleanupInterval` | `func() time.Duration` |
| `config.EventInboxConfig.EffectiveProcessingLease` | `func() time.Duration` |
| `config.EventInboxConfig.EffectiveRetention` | `func() time.Duration` |
| `config.EventOutboxTransportConfig.EffectiveCleanupInterval` | `func() time.Duration` |
| `config.EventOutboxTransportConfig.EffectiveCompletedTTL` | `func() time.Duration` |
| `config.LockoutConfig.IsEnabled` | `func() bool` |
| `config.LockoutConfig.EffectiveMaxFailures` | `func() int` |
| `config.LockoutConfig.EffectiveWindow` | `func() time.Duration` |
| `config.LockoutConfig.EffectiveLockDuration` | `func() time.Duration` |
| `config.LockoutConfig.EffectiveStrategy` | `func() github.com/coldsmirk/vef-framework-go/config.LockoutStrategy` |
| `config.LockoutConfig.EffectiveBackoffBase` | `func() time.Duration` |
| `config.LockoutConfig.EffectiveBackoffMax` | `func() time.Duration` |
| `config.LockoutConfig.EffectiveKey` | `func() github.com/coldsmirk/vef-framework-go/config.LockoutKey` |
| `config.LockoutConfig.Validate` | `func() error` |
| `config.SecurityConfig.EffectiveTokenType` | `func() github.com/coldsmirk/vef-framework-go/config.TokenType` |
| `config.SecurityConfig.Validate` | `func() error` |
| `config.SessionConfig.EffectiveOnExceed` | `func() github.com/coldsmirk/vef-framework-go/config.SessionExceedPolicy` |
| `config.SessionConfig.EffectiveIdleTTL` | `func() time.Duration` |
| `config.SessionConfig.EffectiveMaxLifetime` | `func() time.Duration` |
| `config.SessionConfig.IsSliding` | `func() bool` |
| `config.StorageConfig.EffectiveClaimTTL` | `func() time.Duration` |
| `config.StorageConfig.EffectiveDeleteBatchSize` | `func() int` |
| `config.StorageConfig.EffectiveDeleteConcurrency` | `func() int` |
| `config.StorageConfig.EffectiveDeleteLeaseWindow` | `func() time.Duration` |
| `config.StorageConfig.EffectiveDeleteMaxAttempts` | `func() int` |
| `config.StorageConfig.EffectiveDeleteWorkerInterval` | `func() time.Duration` |
| `config.StorageConfig.EffectiveMaxPendingClaims` | `func() int` |
| `config.StorageConfig.EffectiveMaxUploadSize` | `func() int64` |
| `config.StorageConfig.EffectiveSweepBatchSize` | `func() int` |
| `config.StorageConfig.EffectiveSweepInterval` | `func() time.Duration` |

### Method Semantics

| Method family | Behavior |
| --- | --- |
| `config.Config.Unmarshal(key, target)` | Reads the requested key into `target`; the internal Viper-backed implementation uses `config` struct tags and ignores untagged fields. |
| `config.DataSourcesConfig.Primary()` | Returns `Map[config.PrimaryDataSourceName]`. If the map lacks `primary`, the method returns the zero `config.DataSourceConfig`; framework startup validates the `primary` entry separately. |
| `config.ApprovalConfig.ApplyDefaults()` | Mutates the receiver in place. Non-positive durations and counts become: timeout scan `1m`, pre-warning scan `5m`, cleanup scan `24h`, delegation max depth `10`, form snapshot retention `90d`, urge record retention `30d`, CC record retention `90d`. It does not enable `AutoMigrate`. |
| `config.StorageConfig.Effective...` | Each storage accessor returns the configured value only when it is strictly positive; zero or negative values re-select the exported default constants. Defaults are: max upload size `config.DefaultMaxUploadSize` (`1073741824`, 1 GiB), claim TTL `config.DefaultClaimTTL` (`24h`), max pending claims `100`, sweep interval `5m`, sweep batch size `200`, delete worker interval `5m`, delete batch size `100`, delete concurrency `8`, delete max attempts `12`, delete lease window `5m`. |
| `config.EventConfig.EffectiveDefaultTransport()` | Returns `DefaultTransport` or `"memory"` when unset. |
| `config.EventConfig.EffectiveAsyncQueueSize()` | Returns `AsyncQueueSize` when positive, otherwise `4096`. |
| `config.EventConfig.EffectiveAsyncWorkers()` | Returns `AsyncWorkers` when positive, otherwise `4`. |
| `config.EventConfig.EffectivePublishTimeout()` | Returns `PublishTimeout` when positive, otherwise `5s`. |
| `config.EventOutboxTransportConfig.EffectiveCleanupInterval()` | Returns `CleanupInterval` when positive, otherwise `1h`. |
| `config.EventOutboxTransportConfig.EffectiveCompletedTTL()` | Returns `CompletedTTL` when positive, otherwise `168h`. |
| `config.EventInboxConfig.EffectiveRetention()` | Returns `Retention` when positive, otherwise `168h`. |
| `config.EventInboxConfig.EffectiveProcessingLease()` | Returns `ProcessingLease` when positive, otherwise `10m`. |
| `config.EventInboxConfig.EffectiveCleanupInterval()` | Returns `CleanupInterval` when positive, otherwise `1h`. |
| `config.EventConfig.Validate()` | Runs only when `EventConfig.Middleware.Inbox` is true and `EventConfig.Transports.Outbox.Enabled` is true. It treats `max_retries <= 0` as `10`, computes the worst-case exponential backoff horizon as `sum(2^k seconds)`, saturates overflow fail-closed, and returns an error wrapping `config.ErrInboxRetentionTooShort` when `inbox.retention <= horizon`. |
| `config.SecurityConfig.EffectiveTokenType()` | Returns `TokenType` or `config.TokenTypeJWT` ("jwt_token") when unset. |
| `config.SecurityConfig.Validate()` | Rejects an out-of-enum `TokenType` (wrapping `config.ErrInvalidTokenType`) or `SessionConfig.OnExceed` (wrapping `config.ErrInvalidSessionOnExceed`) so a configuration typo fails fast at boot. |
| `config.LockoutConfig.IsEnabled()` | Returns `true` when `Enabled` is nil (lockout is on by default) or when it points to `true`. |
| `config.LockoutConfig.Effective...()` | Each accessor returns the configured value only when it is strictly positive, otherwise it re-selects the matching default constant: `MaxFailures` -> `config.DefaultLockoutMaxFailures` (`10`), `Window` -> `config.DefaultLockoutWindow` (`15m`), `LockDuration` -> `config.DefaultLockoutLockDuration` (`15m`), `Strategy` -> `config.LockoutStrategyLock` ("lock") when unset, `BackoffBase` -> `config.DefaultLockoutBackoffBase` (`1s`), `BackoffMax` -> `config.DefaultLockoutBackoffMax` (`15m`), `Key` -> `config.LockoutKeyUserIP` ("user_ip") when unset. |
| `config.LockoutConfig.Validate()` | Rejects an out-of-enum `Strategy` (wrapping `config.ErrInvalidLockoutStrategy`) or `Key` (wrapping `config.ErrInvalidLockoutKey`). |
| `config.SessionConfig.EffectiveOnExceed()` | Returns `OnExceed` or `config.SessionExceedEvictOldest` ("evict_oldest") when unset. |
| `config.SessionConfig.EffectiveIdleTTL()` | Returns `IdleTTL` when positive, otherwise `config.DefaultSessionIdleTTL` (`30m`). |
| `config.SessionConfig.EffectiveMaxLifetime()` | Returns `MaxLifetime` when positive, otherwise `config.DefaultSessionMaxLifetime` (`7 * 24h`). |
| `config.SessionConfig.IsSliding()` | Returns `true` when `Sliding` is nil (idle-timeout renewal is on by default) or when it points to `true`. |

`DataSourcesConfig.Map` is intentionally untagged. The internal config module unmarshals `vef.data_sources` into a `map[string]config.DataSourceConfig` first and then wraps it in `DataSourcesConfig{Map: sources}`; this preserves arbitrary data-source names while still reserving `config.PrimaryDataSourceName` (`"primary"`) for the framework-wide `orm.DB`.

## See also

- [Configuration](../getting-started/configuration) for explanations and setup examples
- [Built-in Resources](./built-in-resources) for the modules these settings affect at runtime
