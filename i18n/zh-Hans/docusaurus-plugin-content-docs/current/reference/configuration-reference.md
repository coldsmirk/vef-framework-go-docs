---
sidebar_position: 1
---

# 配置参考

这一页按区块汇总 VEF 在启动期间会读取的配置项。

最小起步配置：

```toml
[vef.app]
name = "my-app"
port = 8080

[vef.data_source]
type = "sqlite"
```

## 配置文件查找路径

内部配置模块会按以下顺序查找 `application.toml`：

- `./configs`
- `$VEF_CONFIG_PATH`
- `.`
- `../configs`

## 常用环境变量

关键环境变量包括：

- `VEF_CONFIG_PATH`
- `VEF_LOG_LEVEL`
- `VEF_NODE_ID`
- `VEF_I18N_LANGUAGE`

## `vef.app`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `name` | `string` | 应用名称，会影响部分运行时行为，例如 JWT audience 默认值。 |
| `port` | `uint16` | HTTP 服务端口。 |
| `body_limit` | `string` | Fiber body limit，例如 `10mib`。 |

## `vef.data_source`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `type` | `postgres \| mysql \| sqlite` | 当前运行时支持的数据库类型。`oracle` / `sqlserver` 暂未实现。 |
| `host` | `string` | 网络数据库主机。 |
| `port` | `uint16` | 网络数据库端口。 |
| `user` | `string` | 数据库用户名。 |
| `password` | `string` | 数据库密码。 |
| `database` | `string` | 数据库名。 |
| `schema` | `string` | 支持 schema 的驱动下使用的 schema 名。 |
| `path` | `string` | SQLite 文件路径。 |
| `enable_sql_guard` | `bool` | 是否启用 SQL guard。 |

说明：

- 当前运行时注册的 provider 支持 `postgres`、`mysql`、`sqlite`。`oracle` 和 `sqlserver` 是 `DBKind` 常量留作未来扩展，目前未实现，配置后会在启动时报 `database.ErrUnsupportedDBKind`。

## `vef.cors`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `enabled` | `bool` | 是否启用 CORS middleware。 |
| `allow_origins` | `[]string` | 允许的来源列表。 |

## `vef.security`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `token_expires` | `duration` | access token 生命周期。 |
| `refresh_not_before` | `duration` | refresh token 最早可刷新时间。 |
| `login_rate_limit` | `int` | 登录接口限流。 |
| `refresh_rate_limit` | `int` | refresh 接口限流。 |

## `vef.redis`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `host` | `string` | Redis host。 |
| `port` | `uint16` | Redis port。 |
| `user` | `string` | Redis 用户名。 |
| `password` | `string` | Redis 密码。 |
| `database` | `uint8` | Redis database 编号。 |
| `network` | `string` | `tcp` 或 `unix`。 |

说明：

- 默认 `vef.Run(...)` 启动图包含 Redis 模块
- 只有当依赖图里真的有组件使用 `*redis.Client` 或其他 Redis 相关能力时，Redis 才会成为实际前提
- 如果这些字段都不写，客户端仍然会回退到 `127.0.0.1:6379` 与 `tcp`

## `vef.storage`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `provider` | `memory \| minio \| filesystem` | 存储 provider 选择。 |
| `minio.endpoint` | `string` | MinIO endpoint。 |
| `minio.access_key` | `string` | MinIO access key。 |
| `minio.secret_key` | `string` | MinIO secret key。 |
| `minio.bucket` | `string` | bucket 名。 |
| `minio.region` | `string` | region。 |
| `minio.use_ssl` | `bool` | 是否使用 HTTPS。 |
| `filesystem.root` | `string` | filesystem provider 根目录。 |

## `vef.monitor`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `sample_interval` | `duration` | 采样间隔。 |
| `sample_duration` | `duration` | 采样窗口时长。 |

## `vef.mcp`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `enabled` | `bool` | 是否启用 MCP server。 |
| `require_auth` | `bool` | `/mcp` 是否要求 Bearer token。 |

## `vef.approval`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `auto_migrate` | `bool` | 是否在启动时执行 approval DDL 迁移。 |
| `timeout_scan_interval` | `duration` | 超时扫描器轮询节奏，默认 1m。 |
| `pre_warning_scan_interval` | `duration` | 预警扫描器轮询节奏，默认 5m。 |
| `cleanup_scan_interval` | `duration` | 保留期清理任务节奏，默认 24h。 |
| `delegation_max_depth` | `int` | 委托链最大深度，默认 10。 |
| `form_snapshot_retention` | `duration` | apv_form_snapshot 保留期，默认 90 天。 |
| `urge_record_retention` | `duration` | apv_urge_record 保留期，默认 30 天。 |
| `cc_record_retention` | `duration` | apv_cc_record 保留期，默认 90 天。 |

> 原本归属 `[vef.approval]` 的 outbox 配置已在 v0.21 移至 `[vef.event.transports.outbox]`，详见 [事件总线](../features/event-bus)。

## `vef.event`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `default_transport` | `string` | 路由回退使用的 transport 名（默认 `memory`）。 |
| `async_queue_size` | `int` | `WithAsync` 异步队列容量。 |
| `async_workers` | `int` | 异步队列 worker 数量。 |
| `publish_timeout` | `duration` | 单次 Publish 调用上限。 |
| `transports.memory.*` | — | 内存 transport 配置：`queue_size`、`full_policy`、`publish_timeout`。 |
| `transports.outbox.*` | — | outbox transport 配置：`enabled`、`relay_interval`、`max_retries`、`batch_size`、`lease_multiplier`、`min_lease`、`sink`、`cleanup_interval`、`completed_ttl`。 |
| `transports.redis_stream.*` | — | Redis Streams transport 配置：`enabled`、`stream_prefix`、`max_len_approx`、`block_timeout`、`claim_idle`、`claim_interval`、`claim_batch_size`、`consumer_id`、`start_id`。 |
| `middleware.*` | `bool` | 中间件开关：`logging`、`tracing`、`tracing_strict`、`metrics`、`recover`、`inbox`。 |
| `inbox.*` | — | Inbox 去重表配置：`retention`、`processing_lease`、`cleanup_interval`。 |
| `routing` | `[]{pattern, transports}` | 路由规则列表，自顶向下匹配。 |

## 延伸阅读

- [配置](../getting-started/configuration)：配置项的解释与实际示例
- [内置资源](./built-in-resources)：这些配置会影响哪些默认模块
