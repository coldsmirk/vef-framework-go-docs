---
sidebar_position: 4
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

These sections map directly to the public config package and internal module constructors.

### `vef.app`

Application-level settings:

```toml
[vef.app]
name = "my-app"
port = 8080
body_limit = "10mib"
```

Key fields:

- `name`: used as the app name and as input to JWT audience generation
- `port`: HTTP port for the Fiber app
- `body_limit`: parsed by Fiber; defaults to `10mib` when omitted

### `vef.data_source`

Database connection settings:

```toml
[vef.data_source]
type = "postgres"
host = "127.0.0.1"
port = 5432
user = "postgres"
password = "postgres"
database = "my_app"
schema = "public"
enable_sql_guard = true
```

Supported `type` values:

- `postgres`
- `mysql`
- `sqlite`
- `oracle`
- `sqlserver`

For SQLite, `path` is optional. When it is omitted, the framework uses a shared in-memory database.

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
token_expires = "2h"
refresh_not_before = "1h"
login_rate_limit = 6
refresh_rate_limit = 1
```

Runtime notes:

- `refresh_not_before` defaults to half of the access-token window when not set, which is `15m` in the current runtime
- login and refresh rate limits are also normalized by the security module

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

### `vef.redis`

The default boot graph includes the Redis module during `vef.Run(...)`.

Redis becomes a practical prerequisite only when something in the dependency graph actually constructs and uses `*redis.Client` or another Redis-backed capability.

When Redis is used and settings are omitted, the client still defaults to:

- host: `127.0.0.1`
- port: `6379`
- network: `tcp`

So in minimal examples, only add `vef.redis` when the application really depends on Redis.

### `vef.monitor`

Monitoring configuration is injected into the monitor module. The module also applies its own defaults internally.

### `vef.mcp`

MCP support is present in the runtime, but the MCP server only activates when enabled in configuration.

### `vef.approval`

Approval workflow engine settings:

```toml
[vef.approval]
auto_migrate = true
outbox_relay_interval = 5
outbox_max_retries = 10
outbox_batch_size = 100
```

Key fields:

- `auto_migrate`: automatically create approval tables on startup
- `outbox_relay_interval`: polling interval in seconds (default: 5)
- `outbox_max_retries`: max retry attempts for outbox events (default: 10)
- `outbox_batch_size`: max events per poll (default: 100)

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
