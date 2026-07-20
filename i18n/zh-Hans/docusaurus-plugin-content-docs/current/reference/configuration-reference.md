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

[vef.data_sources.primary]
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
- `VEF_I18N_LANGUAGE`

## `vef.app`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `name` | `string` | 应用名称，会影响部分运行时行为，例如 JWT audience 默认值。 |
| `port` | `uint16` | HTTP 服务端口。 |
| `body_limit` | `string` | Fiber body limit，例如 `10mib`；未配置时默认 `32mib`。 |
| `trusted_proxies` | `[]string` | 允许设置 `X-Forwarded-For` 的代理 IP 或 CIDR 列表；为空时只信任直接连接来源。 |

## `vef.api`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `rate_limit.max` | `int` | 未声明自己 `OperationSpec.RateLimit` 的操作使用的默认限流；默认 `100`。 |
| `rate_limit.period` | `duration` | 默认限流的时间窗口；默认 `5m`。 |

限流按操作 × 客户端计数（resource、version、action、客户端 IP、principal
ID），且按节点独立统计。每个端点自己的 `OperationSpec.RateLimit` 仍然优先；
见 [API](../building-apis/api#限流)。

## `vef.data_sources`

`vef.data_sources` 是以数据源名称为 key 的 map。`primary` 条目必填，并为全框架注入的 `orm.DB` 提供来源；其他条目会用各自 map key 注册到数据源 registry。

示例：

```toml
[vef.data_sources.primary]
type = "sqlite"

[vef.data_sources.analytics]
type = "sqlite"
path = "./analytics.db"
```

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
| `ssl_mode` | `disable \| require \| verify-ca \| verify-full` | 网络数据库 dialect 的 TLS 模式；省略时等价于 `disable`。 |
| `ssl_root_cert` | `string` | `verify-ca` 和 `verify-full` 使用的可选 PEM CA bundle 路径；为空时使用主机系统证书池。 |

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
| `secret` | `string` | 十六进制 JWT signing key。为空时框架会生成进程内临时 key 并输出警告；token 无法跨重启或多节点继续使用。如果设置为公开的 `security.DefaultJWTSecret`，启动时也会警告生产环境必须替换。 |
| `token_expires` | `duration` | refresh token 生命周期；默认 `168h`。 |
| `refresh_not_before` | `duration` | refresh token 最早可使用时间；默认 `15m`，也就是固定 `30m` access token 生命周期的一半。 |
| `login_rate_limit` | `int` | 登录接口限流；默认 `6`。 |
| `refresh_rate_limit` | `int` | refresh 接口限流；默认 `1`。 |
| `ip_whitelists` | `map[string][]string` | 内置 `ip` 认证策略使用的命名 IP 白名单（IP 或 CIDR 条目）；TOML key 会被转成小写，无参数的 `api.IPAuth()` 指向 `default`。 |
| `api_keys` | `map[string]{key, roles}` | 默认 `security.APIKeyLoader` 为 `api_key` 认证策略提供的静态 API key；每项携带密钥 `key`（高熵随机串）与授予主体的 `roles`。TOML key 会被转成小写 |
| `basic_accounts` | `map[string]{password, roles}` | 默认 `security.BasicAccountLoader` 为 `http_basic` 认证策略提供的静态服务账号；map 键即用户名。这些是机器对机器凭证，不是用户密码 |
| `lockout.*` | — | 登录接口的暴力破解防护：`enabled` 默认 `true`、`max_failures` 默认 `10`、`window` 默认 `15m`、`lock_duration` 默认 `15m`、`strategy`（`lock` \| `backoff`）默认 `lock`、`backoff_base` 默认 `1s`、`backoff_max` 默认 `15m`、`key`（`user` \| `ip` \| `user_ip`）默认 `user_ip`。 |
| `password_policy.*` | — | 密码强度规则；每个字段都是可选项（零值表示不启用该规则）：`min_length`、`max_length`、`require_upper`、`require_lower`、`require_digit`、`require_symbol`、`min_char_classes`、`disallow_username`、`blocklist`、`history_depth`（防重用，需要应用自行实现 `security.PasswordHistoryStore`）、`max_age`（过期策略，需要应用自行实现 `security.PasswordMetadataLoader`）。 |
| `token_type` | `jwt_token \| opaque_token` | 登录 token 机制；默认 `jwt_token`。会话控制（并发数限制、强制下线、续期）只在 `opaque_token` 下可用。 |
| `session.*` | — | opaque token 的会话调优项，在 `jwt_token` 下不生效：`max_concurrent` 默认 `0`（不限制；并发登录场景下是 best-effort 强制）、`on_exceed`（`reject` \| `evict_oldest`）默认 `evict_oldest`、`idle_ttl` 默认 `30m`、`max_lifetime` 默认 `168h`（7 天）、`sliding` 默认 `true`。 |

说明：

- 内置 JWT token generator 签发的 access token 固定 `30m` 过期；`vef.security.token_expires` 控制的是 refresh token，不是 access token。
- 锁定功能默认开启（`max_failures = 10`）；触发后返回 `security.ErrAccountLocked`（HTTP 429），guard 存储出错时按 fail open 处理。
- 只有注册了 `security.PasswordHistoryStore` 时，`history_depth > 0` 才会把历史密码校验组合进密码策略；只有应用同时接入 `security.PasswordMetadataLoader` 和 `security.NewExpiryPasswordChangeChecker` 时，`max_age` 才会生效。

## `vef.redis`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `enabled` | `bool` | 是否构造 Redis client；默认 `false`。 |
| `host` | `string` | Redis host。 |
| `port` | `uint16` | Redis port。 |
| `user` | `string` | Redis 用户名。 |
| `password` | `string` | Redis 密码。 |
| `database` | `uint8` | Redis database 编号。 |
| `network` | `string` | `tcp` 或 `unix`。 |

说明：

- 默认 `vef.Run(...)` 启动图包含 Redis 模块
- 只有 `enabled = true` 时才会构造 Redis client；`enabled` 为 `false` 或未配置时，框架注入的是 nil `*redis.Client`，并跳过启动 `PING`
- 启用后如果不写 host / port / network，会回退到 `127.0.0.1`、`6379`、`tcp`

## `vef.storage`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `provider` | `memory \| minio \| filesystem` | 存储 provider 选择。 |
| `auto_migrate` | `bool` | 是否在启动时执行 storage DDL 迁移。 |
| `minio.endpoint` | `string` | MinIO endpoint。 |
| `minio.access_key` | `string` | MinIO access key。 |
| `minio.secret_key` | `string` | MinIO secret key。 |
| `minio.bucket` | `string` | bucket 名。 |
| `minio.region` | `string` | region。 |
| `minio.use_ssl` | `bool` | 是否使用 HTTPS。 |
| `filesystem.root` | `string` | filesystem provider 根目录。 |
| `max_upload_size` | `int64` | 单个对象上传大小上限，默认 1 GiB。 |
| `claim_ttl` | `duration` | upload claim 有效期，默认 24h。 |
| `max_pending_claims` | `int` | 单个 principal 可持有的 pending claim 上限，默认 100。 |
| `allow_public_uploads` | `bool` | 是否允许客户端请求 public upload；默认关闭。 |
| `sweep_interval` | `duration` | 过期 claim 扫描间隔，默认 5m。 |
| `sweep_batch_size` | `int` | 单次 claim sweep 处理上限，默认 200。 |
| `delete_worker_interval` | `duration` | pending-delete worker 轮询间隔，默认 5m。 |
| `delete_batch_size` | `int` | 单次 delete worker 租约批量，默认 100。 |
| `delete_concurrency` | `int` | 单轮对象删除并发上限，默认 8。 |
| `delete_max_attempts` | `int` | 删除重试预算，默认 12。 |
| `delete_lease_window` | `duration` | 删除任务租约窗口，默认 5m。 |

说明：

- `provider` 为空时会选择内存存储并输出警告；对象会在进程重启后丢失。
- `vef.storage.auto_migrate = true` 会执行幂等 storage 迁移，并检查 `sys_storage_upload_claim`、`sys_storage_upload_part` 和 `sys_storage_pending_delete`。
- `filesystem.root` 默认是 `./storage`。
- `minio.bucket` 默认依次取 `minio.bucket`、`vef.app.name`、`vef-app`。
- 上传流程和删除 worker 的调优项在框架内有默认值；应用代码需要解析后的有效值时，应使用 `StorageConfig` 的 `Effective...` 方法，而不是直接读取原始字段。

## `vef.monitor`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `sample_interval` | `duration` | 采样间隔；默认 `10s`。 |
| `sample_duration` | `duration` | 采样窗口时长；默认 `2s`。 |

> 没有 `excluded_mounts` 字段：overview 的磁盘摘要只报告根文件系统，
> 不存在挂载点排除（完整挂载清单仍可通过磁盘详情查询获取）。

## `vef.mcp`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `enabled` | `bool` | 是否启用 MCP server。 |
| `require_auth` | `bool` | 默认安全：未配置或为 `true` 时，`/mcp` 要求 Bearer token；只有显式设为 `false` 才允许匿名访问。Go 结构体中 `MCPConfig.RequireAuth` 是 `*bool`，用于区分未配置和 false。 |

## `vef.approval`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `auto_migrate` | `bool` | 显式开启时才会在启动时执行 approval DDL 迁移；`ApprovalConfig.ApplyDefaults()` 不会自动打开它。 |
| `timeout_scan_interval` | `duration` | 超时扫描器轮询节奏，默认 1m。 |
| `pre_warning_scan_interval` | `duration` | 预警扫描器轮询节奏，默认 5m。 |
| `cleanup_scan_interval` | `duration` | 保留期清理任务节奏，默认 24h。 |
| `delegation_max_depth` | `int` | 委托链最大深度，默认 10。 |
| `form_snapshot_retention` | `duration` | apv_form_snapshot 保留期，默认 90 天。 |
| `urge_record_retention` | `duration` | apv_urge_record 保留期，默认 30 天。 |
| `cc_record_retention` | `duration` | 已读 apv_cc_record 记录保留期，默认 90 天。 |
| `business_binding.consistency` | `synchronous \| eventual` | 业务表投影模式；默认 `synchronous`（失败回滚审批动作）。`eventual` 先提交期望状态、由 worker 收敛。 |
| `business_binding.scan_interval` | `duration` | eventual 投影 worker 的扫描节奏，默认 `10s`。 |
| `business_binding.batch_size` | `int` | 每次扫描认领的投影数，默认 `100`。 |

> outbox 配置归属 `[vef.event.transports.outbox]`，而不是 `[vef.approval]`，详见 [事件总线](../infrastructure/event-bus)。

## `vef.cron`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `store.enabled` | `bool` | 打开持久化调度存储；默认 `false`（不加载调度、不触碰任何表，内存调度器不受影响） |
| `store.auto_migrate` | `bool` | 启动时执行 cron DDL 迁移（`crn_schedule`、`crn_fire_request`、`crn_run`） |
| `store.poll_interval` | `duration` | 节点重读调度表的等待上限；默认 `5s`。这是其他节点新建调度的可见性延迟，不是触发精度 |
| `store.batch_size` | `int` | 每个轮询周期认领的调度数；默认 `32` |
| `store.max_concurrent` | `int` | 每节点并发执行的运行数；默认 `16` |
| `store.misfire_threshold` | `duration` | 触发晚到多久后应用调度的错触策略；默认 `1m` |
| `store.heartbeat_interval` | `duration` | 执行器对运行中行的活性节律；默认 `10s` |
| `store.abandoned_after` | `duration` | 心跳陈旧多久后恢复清扫将运行标记为 abandoned；默认 `1m`，必须至少是 `heartbeat_interval` 的两倍 |
| `store.run_timeout` | `duration` | 调度未自带超时时的默认单次运行上限；默认 `0`（不限） |
| `store.run_retention` | `duration` | （每小时清扫）删除超窗的终态流水账行；默认 `0`（永久保留） |

启动校验拒绝负的时长以及比心跳间隔两倍更紧的 `abandoned_after`。见
[持久化调度](../infrastructure/cron-store)。

## `vef.integration`

由可选的集成模块（`vef.IntegrationModule`）读取。

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `auto_migrate` | `bool` | 启动时执行集成模块 DDL 迁移（`itg_*` 表） |
| `secret_key` | `string` | base64 编码密钥，静态加密敏感认证参数与数据源密码，长度按算法（AES：16/24/32 字节；SM4：16 字节）。不设置则明文存储并在启动时警告 |
| `secret_algorithm` | `aes \| sm4` | `secret_key` 使用的算法：`aes`（AES-GCM，默认）或 `sm4`（SM4-GCM）。一种算法封存的值不能在另一种下读取 |
| `run_timeout` | `duration` | 单次适配器脚本执行上限，含线上调用；默认 `30s` |
| `max_response_body` | `int64` | 脚本读取的单个 HTTP 响应体上限；默认 8 MiB |
| `log.mode` | `off \| errors \| all` | 记录哪些调用到 `itg_invocation_log`；默认 `errors` |
| `log.capture_limit` | `int` | 每个捕获负载（输入、输出、线上报文）的字节上限；默认 `4096` |
| `log.mask_fields` | `[]string` | 捕获中额外掩码的 JSON 字段名（不区分大小写），凭证头始终掩码 |
| `log.retention` | `duration` | （每小时清扫）删除超窗的调用日志行；默认 `0`（永久保留） |
| `inbound.rate_limit.max` | `int` | 每窗口每（系统、客户端 IP）允许的入站投递数；默认 `120` |
| `inbound.rate_limit.period` | `duration` | 滑动窗口长度；默认 `1m`。限流器在进程内存中——多节点各自独立计数 |

见[集成引擎](../integration/overview)。

## `vef.push`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `enabled` | `bool` | 打开 WebSocket 推送端点；默认 `false`。端点关闭时 `push.Notifier` 依然可用，投递被静默丢弃 |
| `path` | `string` | 端点路径；默认 `/ws` |
| `allowed_origins` | `[]string` | 握手的浏览器来源白名单；空表示允许所有来源（握手本身经令牌认证，这里是纵深防御） |
| `ping_interval` | `duration` | 服务端心跳周期；连续错过两次 pong 即断开；默认 `30s` |
| `write_timeout` | `duration` | 单个出站帧写入上限；默认 `10s` |
| `send_buffer` | `int` | 每连接出站队列长度；消费过慢的客户端被断开；默认 `32` |
| `max_connections_per_user` | `int` | 每节点每用户并发 socket 上限；`0` 为不限 |
| `session_recheck_interval` | `duration` | 不透明令牌会话复核节律；默认 `60s` |

运行时说明：

- 启用 Redis 后，推送与吊销踢出经 pub/sub 频道
  `vef:push:relay:<redis-db>:<app-name>` 跨节点中继；未设置 `vef.app.name`
  时中继拒绝启动。见[服务端推送](../infrastructure/push)

## `vef.event`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `default_transport` | `string` | 路由回退使用的 transport 名（默认 `memory`）。 |
| `async_queue_size` | `int` | `WithAsync` 异步队列容量，默认 `4096`。 |
| `async_workers` | `int` | 异步队列 worker 数量，默认 `4`。 |
| `publish_timeout` | `duration` | 单次 transport Publish 调用上限，默认 `5s`。 |
| `transports.memory.*` | — | 内存 transport 配置：`queue_size` 默认 `1024`，`full_policy` 默认 `error`，`publish_timeout` 默认不设上限，且只在 `full_policy = "block"` 时生效。 |
| `transports.outbox.*` | — | outbox transport 配置：`enabled`、`relay_interval` 默认 `10s`、`max_retries` 默认 `10`、`batch_size` 默认 `100`、`lease_multiplier` 默认 `4`、`min_lease` 默认 `15s`、`sink` 默认 `memory`、`cleanup_interval` 默认 `1h`、`completed_ttl` 默认 `168h`；cleanup 字段属于框架配置，不是 `event/transport/outbox.Config` 字段。 |
| `transports.redis_stream.*` | — | Redis Streams transport 配置：`enabled`、`stream_prefix` 默认 `vef:events:`、`max_len_approx` 默认 `0`（不裁剪）、`block_timeout` 默认 `5s`、`claim_idle` 默认 `60s`、`claim_interval` 默认 `30s`、`claim_batch_size` 默认 `64`、`reaper_concurrency` 默认 `4`、`handler_timeout` 默认 `30s`、`setup_timeout` 默认 `5s`、`consumer_id` 默认前缀 `vef`、`start_id` 默认 `0`（`"$"` 表示新建 group 跳过 backlog）、`idle_group_retention` 默认 `0`（关闭孤儿消费组回收）、`idle_group_sweep_interval` 默认 `10m`。 |
| `middleware.*` | `bool` | 中间件开关：`logging`、`tracing`、`tracing_strict`、`metrics`、`recover`、`inbox`。 |
| `inbox.*` | — | Inbox 去重表配置：`retention` 默认 `168h`、`processing_lease` 默认 `10m`、`cleanup_interval` 默认 `1h`。 |
| `routing` | `[]{pattern, transports}` | 路由规则列表，按 `path.Match` 语义自顶向下匹配。 |

## config 包 API 参考

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
| `config.Config.Unmarshal(key, target)` | 读取指定 key 并写入 `target`；内部 Viper 实现使用 `config` struct tags，并忽略没有 tag 的字段。 |
| `config.DataSourcesConfig.Primary()` | 返回 `Map[config.PrimaryDataSourceName]`。如果 map 中没有 `primary`，方法会返回零值 `config.DataSourceConfig`；框架启动阶段会另行校验 `primary` 是否存在。 |
| `config.ApprovalConfig.ApplyDefaults()` | 原地修改 receiver。非正数 duration 或 count 会变成：timeout scan `1m`、pre-warning scan `5m`、cleanup scan `24h`、delegation max depth `10`、form snapshot retention `90d`、urge record retention `30d`、CC record retention `90d`。它不会启用 `AutoMigrate`。 |
| `config.StorageConfig.Effective...` | 每个 storage accessor 只有在配置值严格为正时才返回配置值；零值或负值都会重新选择导出的默认常量。默认值包括：max upload size `config.DefaultMaxUploadSize`（`1073741824`，1 GiB）、claim TTL `config.DefaultClaimTTL`（`24h`）、max pending claims `100`、sweep interval `5m`、sweep batch size `200`、delete worker interval `5m`、delete batch size `100`、delete concurrency `8`、delete max attempts `12`、delete lease window `5m`。 |
| `config.EventConfig.EffectiveDefaultTransport()` | 返回 `DefaultTransport`；未配置时返回 `"memory"`。 |
| `config.EventConfig.EffectiveAsyncQueueSize()` | `AsyncQueueSize` 为正时返回配置值，否则返回 `4096`。 |
| `config.EventConfig.EffectiveAsyncWorkers()` | `AsyncWorkers` 为正时返回配置值，否则返回 `4`。 |
| `config.EventConfig.EffectivePublishTimeout()` | `PublishTimeout` 为正时返回配置值，否则返回 `5s`。 |
| `config.EventOutboxTransportConfig.EffectiveCleanupInterval()` | `CleanupInterval` 为正时返回配置值，否则返回 `1h`。 |
| `config.EventOutboxTransportConfig.EffectiveCompletedTTL()` | `CompletedTTL` 为正时返回配置值，否则返回 `168h`。 |
| `config.EventInboxConfig.EffectiveRetention()` | `Retention` 为正时返回配置值，否则返回 `168h`。 |
| `config.EventInboxConfig.EffectiveProcessingLease()` | `ProcessingLease` 为正时返回配置值，否则返回 `10m`。 |
| `config.EventInboxConfig.EffectiveCleanupInterval()` | `CleanupInterval` 为正时返回配置值，否则返回 `1h`。 |
| `config.EventConfig.Validate()` | 只在 `EventConfig.Middleware.Inbox` 为 true 且 `EventConfig.Transports.Outbox.Enabled` 为 true 时运行。它把 `max_retries <= 0` 当作 `10`，按 `sum(2^k seconds)` 计算最坏 exponential backoff horizon，溢出时饱和并 fail-closed；当 `inbox.retention <= horizon` 时返回包装 `config.ErrInboxRetentionTooShort` 的错误。 |
| `config.SecurityConfig.EffectiveTokenType()` | 返回 `TokenType`；未配置时返回 `config.TokenTypeJWT`（`"jwt_token"`）。 |
| `config.SecurityConfig.Validate()` | 拒绝超出枚举范围的 `TokenType`（包装 `config.ErrInvalidTokenType`）或 `SessionConfig.OnExceed`（包装 `config.ErrInvalidSessionOnExceed`），让配置笔误在启动时立即失败。 |
| `config.LockoutConfig.IsEnabled()` | `Enabled` 为 nil（锁定默认开启）或指向 `true` 时返回 `true`。 |
| `config.LockoutConfig.Effective...()` | 每个 accessor 只有在配置值严格为正时才返回配置值，否则重新选择对应的默认常量：`MaxFailures` -> `config.DefaultLockoutMaxFailures`（`10`）、`Window` -> `config.DefaultLockoutWindow`（`15m`）、`LockDuration` -> `config.DefaultLockoutLockDuration`（`15m`）、`Strategy` 未配置时 -> `config.LockoutStrategyLock`（`"lock"`）、`BackoffBase` -> `config.DefaultLockoutBackoffBase`（`1s`）、`BackoffMax` -> `config.DefaultLockoutBackoffMax`（`15m`）、`Key` 未配置时 -> `config.LockoutKeyUserIP`（`"user_ip"`）。 |
| `config.LockoutConfig.Validate()` | 拒绝超出枚举范围的 `Strategy`（包装 `config.ErrInvalidLockoutStrategy`）或 `Key`（包装 `config.ErrInvalidLockoutKey`）。 |
| `config.SessionConfig.EffectiveOnExceed()` | 返回 `OnExceed`；未配置时返回 `config.SessionExceedEvictOldest`（`"evict_oldest"`）。 |
| `config.SessionConfig.EffectiveIdleTTL()` | `IdleTTL` 为正时返回配置值，否则返回 `config.DefaultSessionIdleTTL`（`30m`）。 |
| `config.SessionConfig.EffectiveMaxLifetime()` | `MaxLifetime` 为正时返回配置值，否则返回 `config.DefaultSessionMaxLifetime`（`7 * 24h`）。 |
| `config.SessionConfig.IsSliding()` | `Sliding` 为 nil（默认开启滑动续期）或指向 `true` 时返回 `true`。 |

`DataSourcesConfig.Map` 有意不带 tag。内部配置模块会先把 `vef.data_sources` unmarshal 成 `map[string]config.DataSourceConfig`，再包装成 `DataSourcesConfig{Map: sources}`；这样既保留任意数据源名称，也用 `config.PrimaryDataSourceName`（`"primary"`）为全框架 `orm.DB` 保留主数据源。

## 延伸阅读

- [配置](../getting-started/configuration)：配置项的解释与实际示例
- [内置资源](./built-in-resources)：这些配置会影响哪些默认模块
