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

- the current runtime provider registry supports `postgres`, `mysql`, and `sqlite`. `oracle` and `sqlserver` are declared as `DBKind` constants for future use but have no runtime provider yet.

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
- `timeout_scan_interval`
- `pre_warning_scan_interval`
- `cleanup_scan_interval`
- `delegation_max_depth`
- `form_snapshot_retention`
- `urge_record_retention`
- `cc_record_retention`

> Outbox-related fields moved to `[vef.event.transports.outbox]` in v0.21; see [Event Bus](../features/event-bus).

## `vef.event`

- `default_transport`
- `async_queue_size`
- `async_workers`
- `publish_timeout`
- `transports.memory.queue_size` / `full_policy` / `publish_timeout`
- `transports.outbox.enabled` / `relay_interval` / `max_retries` / `batch_size` / `lease_multiplier` / `min_lease` / `sink` / `cleanup_interval` / `completed_ttl`
- `transports.redis_stream.enabled` / `stream_prefix` / `max_len_approx` / `block_timeout` / `claim_idle` / `claim_interval` / `claim_batch_size` / `consumer_id` / `start_id`
- `middleware.logging` / `tracing` / `tracing_strict` / `metrics` / `recover` / `inbox`
- `inbox.retention` / `processing_lease` / `cleanup_interval`
- `routing` (list of `{pattern, transports}`)

## See also

- [Configuration](../getting-started/configuration) for explanations and setup examples
- [Built-in Resources](./built-in-resources) for the modules these settings affect at runtime
