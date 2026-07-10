---
sidebar_position: 5
---

# Configuration

VEF reads configuration from `application.toml` through the `config` module, then injects strongly typed config structs into the rest of the runtime.

## File lookup order

At startup, the framework config loader searches for `application.toml` in:

- `./configs`
- `$VEF_CONFIG_PATH`
- `.`
- `../configs`

If the file cannot be read, startup fails immediately.

## Core sections

These sections map directly to the public config package and internal module constructors. For the complete `config` public surface, including exported structs, fields, and methods, see [Configuration Reference](../reference/configuration-reference).

### `vef.app`

Application-level settings:

```toml
[vef.app]
name = "my-app"
port = 8080
body_limit = "32mib"
```

Key fields:

- `name`: used as the app name and as input to JWT audience generation
- `port`: HTTP port for the Fiber app
- `body_limit`: parsed by Fiber; defaults to `32mib` when omitted

### `vef.data_sources`

Database connection settings:

```toml
[vef.data_sources.primary]
type = "postgres"
host = "127.0.0.1"
port = 5432
user = "postgres"
password = "postgres"
database = "my_app"
schema = "public"
enable_sql_guard = true
```

The `primary` entry is mandatory and powers the framework-wide `orm.DB`
injection. Additional named data sources use the same shape:

```toml
[vef.data_sources.analytics]
type = "sqlite"
path = "./analytics.db"
```

Supported `type` values (drivers registered in the framework runtime):

- `postgres`
- `mysql`
- `sqlite`

For SQLite, `path` is optional. When omitted, the framework uses a shared in-memory database.

> The `config.DBKind` enum also declares `oracle` and `sqlserver` constants for future use, but the framework does not currently ship runtime providers for them — configuring those values returns `database.ErrUnsupportedDBKind` at boot.

### `vef.cors`

CORS middleware settings:

```toml
[vef.cors]
enabled = true
allow_origins = ["http://localhost:3000", "https://my-app.com"]
```

Key fields:

- `enabled`: enable CORS middleware
- `allow_origins`: list of allowed origins

### `vef.security`

Security-related runtime settings:

```toml
[vef.security]
secret = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
token_expires = "168h"
refresh_not_before = "15m"
login_rate_limit = 6
refresh_rate_limit = 1
```

Runtime notes:

- `secret` is the hex-encoded JWT signing key. Leave it empty only for local development; the framework then generates an ephemeral per-process key, so tokens do not survive restart or work across nodes. Generate and set a stable private value for production.
- access tokens issued by the built-in JWT token generator expire after `30m`
- `token_expires` controls refresh-token lifetime and defaults to `168h`
- `refresh_not_before` defaults to `15m`, half of the fixed access-token window
- login and refresh rate limits default to `6` and `1` when unset or non-positive

### `vef.storage`

Object storage settings:

```toml
[vef.storage]
provider = "filesystem"

[vef.storage.filesystem]
root = "./data/files"
```

Supported providers:

- `memory`
- `filesystem`
- `minio`

If `provider` is omitted, VEF uses memory storage.
For non-test deployments, set `filesystem` or `minio`; in-memory objects are lost on restart.
The filesystem provider defaults `root` to `./storage`, and MinIO defaults its bucket to `vef.app.name` or `vef-app` when `minio.bucket` is empty.

### `vef.redis`

The default boot graph includes the Redis module during `vef.Run(...)`.

Redis is opt-in. When `vef.redis.enabled` is omitted or false, the framework injects a nil `*redis.Client` and skips startup `PING`; Redis-backed modules that depend on Redis must either stay dormant or require you to enable Redis explicitly.

When `enabled = true` and connection settings are omitted, the client defaults to:

- host: `127.0.0.1`
- port: `6379`
- network: `tcp`

So in minimal examples, leave `vef.redis` out unless the application really depends on Redis. When it does, configure `enabled = true` intentionally.

### `vef.monitor`

Monitoring configuration is injected into the monitor module. The module also applies its own defaults internally.
The default sampling interval is `10s`, with a `2s` sampling window.

### `vef.mcp`

MCP support is present in the runtime, but the MCP server only activates when enabled in configuration.

The `/mcp` endpoint requires Bearer auth by default. If `vef.mcp.require_auth` is omitted or set to `true`, unauthenticated requests are rejected; set it to `false` only for deliberately anonymous MCP surfaces.

### `vef.approval`

Approval workflow engine settings:

```toml
[vef.approval]
auto_migrate              = true
timeout_scan_interval     = "1m"
pre_warning_scan_interval = "5m"
cleanup_scan_interval     = "24h"
delegation_max_depth      = 10
form_snapshot_retention   = "2160h"  # 90 days
urge_record_retention     = "720h"   # 30 days
cc_record_retention       = "2160h"  # 90 days
```

Key fields:

- `auto_migrate`: run the approval DDL migration on startup
- `timeout_scan_interval`: cadence of the timeout scanner (default: 1m)
- `pre_warning_scan_interval`: cadence of the pre-warning scanner (default: 5m)
- `cleanup_scan_interval`: cadence of the retention cleanup job (default: 24h)
- `delegation_max_depth`: maximum delegation chain depth (default: 10)
- `form_snapshot_retention` / `urge_record_retention` / `cc_record_retention`: retention windows for the corresponding tables

`config.ApprovalConfig.ApplyDefaults()` fills the timing and retention defaults above but does not enable `AutoMigrate`; migrations run only when `auto_migrate = true`.

> The outbox-related fields previously lived under `[vef.approval]` (`outbox_relay_interval`, `outbox_max_retries`, `outbox_batch_size`). They have moved to `[vef.event.transports.outbox]` so the framework-wide outbox transport can serve any module — see [Event Bus](../infrastructure/event-bus).

## Environment overrides

VEF uses an environment prefix and dot-to-underscore replacement, so config keys can be overridden with environment variables.

Examples:

- `VEF_CONFIG_PATH`
- `VEF_LOG_LEVEL`
- `VEF_NODE_ID`
- `VEF_I18N_LANGUAGE`

## What configuration does not do

Configuration does **not** replace application composition. You still use code to:

- register resources
- provide services and modules
- register auth loaders and permission resolvers
- register CQRS behaviors
- register MCP providers

Think of configuration as runtime input, not application structure.

## Next step

Once the config file is clear, move to [Project Structure](./project-structure) to organize a real project around modules.
