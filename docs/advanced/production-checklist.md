---
sidebar_position: 4
---

# Production Checklist

Production-readiness concerns are spread across many pages. This page collects
them into one ordered checklist: each item states what to set, why, and the
config key or API, then links to the page that documents it in depth. All
defaults and failure modes below reflect v0.37.0.

## Security

1. **Set `vef.security.secret`.** When unset, the framework generates an
   ephemeral per-process JWT signing key and logs a warning: tokens do not
   survive a restart and do not work across nodes. Generate a stable value with
   `security.GenerateSecret()` and provision it per deployment; startup also
   warns when the value is the public `security.DefaultJWTSecret`. See
   [Authentication Reference](../security/authentication-reference).
2. **Set `vef.app.trusted_proxies` behind a reverse proxy.** With the list
   empty, `X-Forwarded-For` is ignored and the client IP is the direct
   connection peer — behind a load balancer every request then shares the
   proxy's IP for rate-limit keys and IP whitelists. List only proxy IPs or
   CIDR ranges you control. See
   [Configuration Reference](../reference/configuration-reference).
3. **Enable CORS deliberately.** The CORS middleware is registered but inert
   until `vef.cors.enabled = true`; browser clients then need an explicit
   `allow_origins` list. See
   [Configuration Reference](../reference/configuration-reference).
4. **Review the auth endpoint rate limits.** `vef.security.login_rate_limit`
   defaults to `6` and `refresh_rate_limit` to `1` (per key, 5-minute sliding
   window). Limiter state is in-memory per instance, so a multi-node
   deployment multiplies the effective limit by the node count. See
   [Built-in Resources](../reference/built-in-resources).
5. **Keep `vef.mcp.require_auth` on.** The MCP endpoint requires Bearer auth
   when the key is unset or `true`; only an explicit `false` allows anonymous
   access. See [MCP](../ai-integration/mcp).
6. **Use a Redis nonce store for signature auth on multiple nodes.** Signature
   authentication defaults to an in-memory nonce store, so replay protection
   is per process; supply `security.NewRedisNonceStore` for distributed
   deployments. See [Authentication](../security/authentication).
7. **Size `vef.app.body_limit`.** The request body limit defaults to `32mib`;
   lower it when you do not accept large payloads, raise it deliberately when
   you do. See [Configuration Reference](../reference/configuration-reference).

## Data

8. **Turn on database TLS.** `ssl_mode` defaults to `disable`; set `require`,
   `verify-ca`, or `verify-full` (plus `ssl_root_cert` for a private CA) on
   every network data source. See [Data Sources](../data-access/datasources).
9. **Consider `enable_sql_guard = true`.** Off by default. When enabled,
   dangerous raw-SQL statements (`DROP`, `TRUNCATE`, `DELETE`/`UPDATE` without
   `WHERE`) are blocked unless the query context is whitelisted. See
   [Data Sources](../data-access/datasources).
10. **Decide on Redis.** `vef.redis.enabled` defaults to `false`, which injects
    a nil client. What depends on Redis: the `redis_stream` event transport
    (its constructor returns nil when the client is nil, so enabling the
    transport without enabling Redis silently leaves its routes unserved),
    `cache.NewRedis` caches (panic on a nil client), and the Redis nonce store
    from item 6. Challenge tokens are JWT-based and per-operation rate limits
    are in-memory, so neither needs Redis. See
    [Configuration Reference](../reference/configuration-reference).

## Storage

11. **Set a real storage provider.** With `vef.storage.provider` unset, the
    framework falls back to in-memory storage and logs a warning — objects are
    lost on restart. Use `filesystem` or `minio` for any non-test deployment.
    See [Storage](../infrastructure/storage).
12. **Register a `FileACL` before storing private files.** The default ACL
    grants reads only under `pub/` and denies every other key regardless of
    who asks; override it via `vef.SupplyFileACL(...)` once you serve `priv/*`
    files. See [Storage](../infrastructure/storage#fileacl).
13. **Route storage events through the outbox.** The storage module fails fast
    at startup unless `vef.storage.*` events resolve to a transactional
    transport: enable `vef.event.transports.outbox` and add a routing rule, or
    make `outbox` the default transport. See
    [Storage](../infrastructure/storage).

## Events

14. **Pick a production transport.** The default `memory` transport is neither
    durable nor transactional: events are lost on crash or restart, and the
    default queue-full policy fails the publish with `event.ErrQueueFull`. For
    anything that must survive the process, use `outbox` (transactional,
    durable, at-least-once, publish-only, relays into a sink) and/or
    `redis_stream` (durable, at-least-once, cross-process). See
    [Event Bus](../infrastructure/event-bus).
15. **Plan for at-least-once semantics.** Durable transports may deliver
    duplicates: subscribe with `event.WithGroup(...)` (required on
    at-least-once routes) and keep the Inbox middleware enabled for dedupe.
    See [Event Bus](../infrastructure/event-bus).

## Operations

16. **Verify shutdown grace periods.** `vef.Run` stops through the FX
    lifecycle on SIGINT/SIGTERM: the HTTP server gets a 30-second grace period
    for in-flight requests and the overall stop timeout is 60 seconds. Give
    your orchestrator at least that much termination grace. See
    [Lifecycle](../core-concepts/lifecycle).
17. **Decide who may call `sys/monitor`.** Monitoring endpoints require Bearer
    auth by default (any authenticated principal, per-action rate-limit max
    `60`); add permission checks or network controls if host metrics are
    sensitive in your environment. See [Monitor](../infrastructure/monitor)
    and [Built-in Resources](../reference/built-in-resources).
18. **Set the log level.** `VEF_LOG_LEVEL` accepts `debug|info|warn|error` and
    defaults to `info`. See [logx](../utilities/logx).
19. **Inject build info.** Generate build metadata with
    `vef-cli generate-build-info` and supply it via `vef.Supply(BuildInfo)`;
    without it, `sys/monitor` reports `unknown` for app version, build time,
    and git commit. See [CLI Tools](./cli-tools) and
    [Monitor](../infrastructure/monitor).

## Next step

Read the [Configuration Reference](../reference/configuration-reference) for
every key mentioned here, or [Lifecycle](../core-concepts/lifecycle) for what
actually happens between `vef.Run(...)` and the first request.
