---
sidebar_position: 5
---

# 配置

VEF 通过 `config` 模块读取 `application.toml`，再把强类型配置结构注入给后续模块。

## 文件查找顺序

启动时，配置加载器会依次查找：

- `./configs`
- `$VEF_CONFIG_PATH`
- `.`
- `../configs`

只要 `application.toml` 没读到，启动就会直接失败。

## 核心配置段

这些段都直接映射到公开 `config` 包和内部模块构造函数。完整的 `config` public surface，包括导出结构、字段和方法，见 [配置参考](../reference/configuration-reference)。

### `vef.app`

应用级配置：

```toml
[vef.app]
name = "my-app"
port = 8080
body_limit = "32mib"
```

关键字段：

- `name`：应用名，也会参与 JWT audience 的构造
- `port`：Fiber HTTP 服务监听端口
- `body_limit`：请求体大小限制；未配置时默认是 `32mib`

### `vef.api`

操作级默认限流（v0.38），作用于所有未声明自己 `OperationSpec.RateLimit`
的操作：

```toml
[vef.api.rate_limit]
max    = 100   # 默认值
period = "5m"  # 默认值
```

### `vef.data_sources`

数据库配置：

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

`primary` 条目是必填项，它为全框架注入的 `orm.DB` 提供来源。其他命名数据源使用同样结构：

```toml
[vef.data_sources.analytics]
type = "sqlite"
path = "./analytics.db"
```

当前支持的 `type`（框架运行时已经注册的驱动）：

- `postgres`
- `mysql`
- `sqlite`

对 SQLite 来说，`path` 可以省略；省略后框架会使用共享内存数据库。

> `config.DBKind` 枚举里还声明了 `oracle` 和 `sqlserver` 两个常量，留作未来扩展，但框架目前**没有**提供对应的运行时 provider —— 实际配置这两个值会在启动时报 `database.ErrUnsupportedDBKind`。

### `vef.cors`

CORS 中间件配置：

```toml
[vef.cors]
enabled = true
allow_origins = ["http://localhost:3000", "https://my-app.com"]
```

关键字段：

- `enabled`：是否启用 CORS 中间件
- `allow_origins`：允许的来源列表

### `vef.security`

安全相关运行时配置：

```toml
[vef.security]
secret = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
token_expires = "168h"
refresh_not_before = "15m"
login_rate_limit = 6
refresh_rate_limit = 1
```

基于当前实现，有两个要点：

- `secret` 是十六进制 JWT signing key。只建议在本地开发时留空；此时框架会生成进程内临时 key，token 无法跨重启或多节点继续使用。生产环境应生成并设置稳定的私有值
- 内置 JWT token generator 签发的 access token 固定 `30m` 过期
- `token_expires` 控制 refresh token 生命周期，默认 `168h`
- `refresh_not_before` 为空时默认 `15m`，也就是固定 access token 窗口的一半
- 登录和刷新限流为空或非正数时默认分别为 `6` 和 `1`

### `vef.storage`

对象存储配置：

```toml
[vef.storage]
provider = "filesystem"

[vef.storage.filesystem]
root = "./data/files"
```

当前支持：

- `memory`
- `filesystem`
- `minio`

如果 `provider` 为空，VEF 会默认使用内存存储。
非测试部署应显式选择 `filesystem` 或 `minio`；内存存储中的对象会在重启后丢失。
filesystem provider 的 `root` 默认是 `./storage`，MinIO 的 bucket 在 `minio.bucket` 为空时会依次回退到 `vef.app.name` 和 `vef-app`。

### `vef.redis`

默认启动图在 `vef.Run(...)` 中包含 Redis 模块。

Redis 是 opt-in。`vef.redis.enabled` 未配置或为 false 时，框架注入 nil `*redis.Client` 并跳过启动 `PING`；依赖 Redis 的模块要么保持 dormant，要么要求你显式启用 Redis。

当 `enabled = true` 且连接参数省略时，客户端会默认使用：

- host：`127.0.0.1`
- port：`6379`
- network：`tcp`

所以在最小示例里，除非应用确实依赖 Redis，否则可以不写 `vef.redis`。一旦需要 Redis，就应明确配置 `enabled = true`。

### `vef.monitor`

监控配置会注入 monitor 模块，而 monitor 模块自身还会补默认值。
默认采样间隔是 `10s`，采样窗口是 `2s`。

### `vef.mcp`

MCP 相关代码默认在运行时里，但 MCP server 只有在配置里显式启用后才会真正生效。

`/mcp` 端点默认要求 Bearer auth。`vef.mcp.require_auth` 未配置或设为 `true` 时，未认证请求会被拒绝；只有明确需要匿名 MCP surface 时才设为 `false`。

### `vef.approval`

审批工作流引擎配置：

```toml
[vef.approval]
auto_migrate              = true
timeout_scan_interval     = "1m"
pre_warning_scan_interval = "5m"
cleanup_scan_interval     = "24h"
delegation_max_depth      = 10
form_snapshot_retention   = "2160h"  # 90 天
urge_record_retention     = "720h"   # 30 天
cc_record_retention       = "2160h"  # 90 天
```

关键字段：

- `auto_migrate`：启动时执行审批模块 DDL 迁移
- `timeout_scan_interval`：超时扫描器节奏（默认 1m）
- `pre_warning_scan_interval`：预警扫描器节奏（默认 5m）
- `cleanup_scan_interval`：保留期清理任务节奏（默认 24h）
- `delegation_max_depth`：委托链最大深度（默认 10）
- `form_snapshot_retention` / `urge_record_retention` / `cc_record_retention`：相应表的保留期窗口

`config.ApprovalConfig.ApplyDefaults()` 会填充上面的节奏和保留期默认值，但不会启用 `AutoMigrate`；只有 `auto_migrate = true` 时才会执行迁移。

> 老版本里的 `outbox_relay_interval` / `outbox_max_retries` / `outbox_batch_size` 已经从 `[vef.approval]` 搬到 `[vef.event.transports.outbox]`，由全框架共享的 outbox transport 服务所有模块——参考 [事件总线](../infrastructure/event-bus)。

## 环境变量覆盖

VEF 使用固定前缀，并把点号替换成下划线，所以配置可以通过环境变量覆盖。

常见示例：

- `VEF_CONFIG_PATH`
- `VEF_LOG_LEVEL`
- `VEF_NODE_ID`
- `VEF_I18N_LANGUAGE`

## 配置不负责什么

配置文件并不会替代应用组合。你依然需要在代码里：

- 注册资源
- 提供服务和模块
- 注册认证加载器和权限解析器
- 注册 CQRS 行为
- 注册 MCP provider

可以把配置理解成“运行时输入”，不要把它当成“应用结构定义”。

## 下一步

配置清楚之后，继续看 [项目结构](./project-structure)，把代码组织成适合 VEF 的模块形态。
