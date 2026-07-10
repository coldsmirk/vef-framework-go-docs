---
sidebar_position: 4
---

# 生产环境检查清单

生产就绪相关的注意事项分散在许多页面中。本页把它们汇总成一份有序清单：
每一项都说明要设置什么、为什么，以及对应的配置键或 API，并链接到详细
文档所在的页面。以下所有默认值与失效行为均以 v0.37.0 为准。

## 安全（Security）

1. **设置 `vef.security.secret`。** 未设置时，框架会生成一个临时的进程级
   JWT 签名密钥并记录警告：令牌在重启后失效，也无法跨节点使用。使用
   `security.GenerateSecret()` 生成稳定值并按部署环境配置；当值等于公开的
   `security.DefaultJWTSecret` 时启动同样会警告。参见
   [认证参考](../security/authentication-reference)。
2. **反向代理之后要设置 `vef.app.trusted_proxies`。** 列表为空时
   `X-Forwarded-For` 会被忽略，客户端 IP 就是直接连接的对端——在负载均衡
   之后，所有请求的限流 key 和 IP 白名单都会共享代理的 IP。只填写你可控的
   代理 IP 或 CIDR 网段。参见
   [配置参考](../reference/configuration-reference)。
3. **有意识地启用 CORS。** CORS 中间件已注册但在 `vef.cors.enabled = true`
   之前不生效；浏览器客户端还需要显式的 `allow_origins` 列表。参见
   [配置参考](../reference/configuration-reference)。
4. **检查认证端点的限流。** `vef.security.login_rate_limit` 默认 `6`、
   `refresh_rate_limit` 默认 `1`（按 key、5 分钟滑动窗口）。限流状态保存在
   各实例的内存中，多节点部署会按节点数放大实际限额。参见
   [内置资源](../reference/built-in-resources)。
5. **保持 `vef.mcp.require_auth` 开启。** 该键未设置或为 `true` 时 MCP 端点
   要求 Bearer 认证；只有显式的 `false` 才允许匿名访问。参见
   [MCP](../ai-integration/mcp)。
6. **多节点下签名认证要使用 Redis nonce 存储。** 签名认证默认使用内存
   nonce 存储，重放保护仅限单进程；分布式部署请提供
   `security.NewRedisNonceStore`。参见 [认证](../security/authentication)。
7. **调整 `vef.app.body_limit`。** 请求体大小限制默认 `32mib`；不接受大
   载荷时调低，需要时再有意识地调高。参见
   [配置参考](../reference/configuration-reference)。

## 数据（Data）

8. **开启数据库 TLS。** `ssl_mode` 默认 `disable`；每个网络数据源都应设置
   `require`、`verify-ca` 或 `verify-full`（私有 CA 还需 `ssl_root_cert`）。
   参见 [多数据源](../data-access/datasources)。
9. **考虑 `enable_sql_guard = true`。** 默认关闭。开启后，危险的原生 SQL
   语句（`DROP`、`TRUNCATE`、不带 `WHERE` 的 `DELETE`/`UPDATE`）会被拦截，
   除非查询上下文在白名单中。参见 [多数据源](../data-access/datasources)。
10. **决定是否启用 Redis。** `vef.redis.enabled` 默认 `false`，此时注入的是
    nil 客户端。依赖 Redis 的能力有：`redis_stream` 事件 transport（客户端
    为 nil 时其构造函数返回 nil，只开 transport 不开 Redis 会静默导致相关
    路由无人服务）、`cache.NewRedis` 缓存（nil 客户端会 panic），以及第 6
    项的 Redis nonce 存储。质询令牌基于 JWT、按操作限流基于内存，两者都
    不需要 Redis。参见 [配置参考](../reference/configuration-reference)。

## 存储（Storage）

11. **设置真正的存储 provider。** `vef.storage.provider` 未设置时，框架会
    回退到内存存储并记录警告——对象在重启后丢失。任何非测试部署都应使用
    `filesystem` 或 `minio`。参见 [文件存储](../infrastructure/storage)。
12. **存放私有文件前先注册 `FileACL`。** 默认 ACL 只授予 `pub/` 前缀的读
    权限，其余 key 无论调用者是谁一律拒绝；一旦要提供 `priv/*` 文件，就
    通过 `vef.SupplyFileACL(...)` 覆盖。参见
    [文件存储](../infrastructure/storage#fileacl)。
13. **让存储事件走 outbox。** 除非 `vef.storage.*` 事件解析到事务性
    transport，存储模块会在启动时快速失败：启用
    `vef.event.transports.outbox` 并添加路由规则，或把 `outbox` 设为默认
    transport。参见 [文件存储](../infrastructure/storage)。

## 事件（Events）

14. **选择生产用 transport。** 默认的 `memory` transport 既不持久也不支持
    事务：进程崩溃或重启时事件丢失，队列满时的默认策略是让发布失败并返回
    `event.ErrQueueFull`。任何需要跨进程存活的事件都应使用 `outbox`
    （事务性、持久、at-least-once、仅发布、由 relay 转入 sink）和/或
    `redis_stream`（持久、at-least-once、跨进程）。参见
    [事件总线](../infrastructure/event-bus)。
15. **为 at-least-once 语义做准备。** 持久 transport 可能重复投递：订阅时
    使用 `event.WithGroup(...)`（at-least-once 路由上是必需的），并保持
    Inbox 中间件开启以去重。参见 [事件总线](../infrastructure/event-bus)。

## 运维（Operations）

16. **确认优雅停机的宽限期。** `vef.Run` 在收到 SIGINT/SIGTERM 时通过 FX
    生命周期停机：HTTP 服务器有 30 秒宽限期处理进行中的请求，整体停止
    超时为 60 秒。编排系统的终止宽限期至少要给到这个量级。参见
    [应用生命周期](../core-concepts/lifecycle)。
17. **决定谁可以调用 `sys/monitor`。** 监控端点默认要求 Bearer 认证（任何
    已认证主体均可访问，每个操作的限流上限为 `60`）；若主机指标在你的环境
    中属于敏感信息，请补充权限检查或网络层控制。参见
    [监控](../infrastructure/monitor) 和
    [内置资源](../reference/built-in-resources)。
18. **设置日志级别。** `VEF_LOG_LEVEL` 接受 `debug|info|warn|error`，默认
    `info`。参见 [logx](../utilities/logx)。
19. **注入构建信息。** 用 `vef-cli generate-build-info` 生成构建元数据并通过
    `vef.Supply(BuildInfo)` 提供；否则 `sys/monitor` 的应用版本、构建时间和
    git commit 都显示为 `unknown`。参见 [CLI 工具](./cli-tools) 和
    [监控](../infrastructure/monitor)。

## 下一步

阅读[配置参考](../reference/configuration-reference)查看本页提到的所有配置
键，或阅读[应用生命周期](../core-concepts/lifecycle)了解从 `vef.Run(...)`
到第一个请求之间到底发生了什么。
