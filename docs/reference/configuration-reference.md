---
sidebar_position: 1
---

# Configuration Reference

This page summarizes the main config structs currently exposed by the framework.

Minimal starter block:

```toml
[vef.app]
name = "my-app"
port = 8080

[vef.data_source]
type = "sqlite"
```

## `vef.app`

- `name`
- `port`
- `body_limit`

## `vef.data_source`

- `type`
- `host`
- `port`
- `user`
- `password`
- `database`
- `schema`
- `path`
- `enable_sql_guard`

Runtime note:

- the current runtime provider registry supports `postgres`, `mysql`, and `sqlite`

## `vef.cors`

- `enabled`
- `allow_origins`

## `vef.security`

- `token_expires`
- `refresh_not_before`
- `login_rate_limit`
- `refresh_rate_limit`

## `vef.redis`

- `host`
- `port`
- `user`
- `password`
- `database`
- `network`

Runtime note:

- the default `vef.Run(...)` boot graph includes the Redis module
- Redis only becomes a practical prerequisite when some dependency actually uses `*redis.Client` or another Redis-backed capability
- if these fields are omitted, the client defaults to `127.0.0.1:6379` over `tcp`

## `vef.storage`

- `provider`
- `minio`
- `filesystem`

### `vef.storage.minio`

- `endpoint`
- `access_key`
- `secret_key`
- `bucket`
- `region`
- `use_ssl`

### `vef.storage.filesystem`

- `root`

## `vef.monitor`

- `sample_interval`
- `sample_duration`

## `vef.mcp`

- `enabled`
- `require_auth`

## `vef.approval`

- `auto_migrate`
- `outbox_relay_interval`
- `outbox_max_retries`
- `outbox_batch_size`

## See also

- [Configuration](../getting-started/configuration) for explanations and setup examples
- [Built-in Resources](./built-in-resources) for the modules these settings affect at runtime
