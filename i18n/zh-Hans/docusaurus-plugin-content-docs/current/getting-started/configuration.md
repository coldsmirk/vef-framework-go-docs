---
sidebar_position: 4
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

这些段都直接映射到公开 `config` 包和内部模块构造函数。

### `vef.app`

应用级配置：

```toml
[vef.app]
name = "my-app"
port = 8080
body_limit = "10mib"
```

关键字段：

- `name`：应用名，也会参与 JWT audience 的构造
- `port`：Fiber HTTP 服务监听端口
- `body_limit`：请求体大小限制；未配置时默认是 `10mib`

### `vef.data_source`

数据库配置：

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

当前支持的 `type`：

- `postgres`
- `mysql`
- `sqlite`
- `oracle`
- `sqlserver`

对 SQLite 来说，`path` 可以省略；省略后框架会使用共享内存数据库。

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
token_expires = "2h"
refresh_not_before = "1h"
login_rate_limit = 6
refresh_rate_limit = 1
```

基于当前实现，有两个要点：

- `refresh_not_before` 为空时，会默认取 access token 窗口的一半；按当前运行时常量，这个值是 `15m`
- 登录和刷新限流也会在安全模块里做归一化处理

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

### `vef.redis`

默认启动图在 `vef.Run(...)` 中包含 Redis 模块。

但只有当依赖图里真的有组件消费 `*redis.Client` 或其他 Redis 相关能力时，Redis 才会成为实际前提。

如果此时你不写 Redis 配置，客户端仍然会默认使用：

- host：`127.0.0.1`
- port：`6379`
- network：`tcp`

所以在最小示例里，只有当应用确实依赖 Redis 时才需要补 `vef.redis`。

### `vef.monitor`

监控配置会注入 monitor 模块，而 monitor 模块自身还会补默认值。

### `vef.mcp`

MCP 相关代码默认在运行时里，但 MCP server 只有在配置里显式启用后才会真正生效。

### `vef.approval`

审批工作流引擎配置：

```toml
[vef.approval]
auto_migrate = true
outbox_relay_interval = 5
outbox_max_retries = 10
outbox_batch_size = 100
```

关键字段：

- `auto_migrate`：启动时自动创建审批相关表
- `outbox_relay_interval`：outbox 轮询间隔，单位秒（默认 5）
- `outbox_max_retries`：outbox 事件最大重试次数（默认 10）
- `outbox_batch_size`：单次轮询最大事件数（默认 100）

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
