---
sidebar_position: 7
---

# 升级到 v0.35 注意事项

本页是从 `9e7e009 feat: added multi-data source functionality` 到
`v0.35.0` 的跨版本审计地图，覆盖提交范围 `9e7e009^..v0.35.0`。如果你的应用
文档或集成假设停留在多数据源功能之前，先按这里过一遍迁移点。

它不是生成索引的替代品。迁移完成后，精确的 Go 符号和 wire 字段仍以
[公开 API 索引](../reference/public-api-index) 与
[运行时 API 索引](../reference/runtime-api-index) 为准。

## 立即检查

- 把旧的单数据源配置改成 `[vef.data_sources.primary]`。`primary` 是必填且保留
  的主数据源；其他数据源放在 `[vef.data_sources.<name>]` 下。
- 更新重命名的 Go 符号：旧的 env key prefix 常量改为 `config.EnvPrefix`，
  旧的 CORS config 类型改为 `config.CORSConfig`，`datasource.Spec.Cfg` 改为
  `datasource.Spec.Config`，SSL 常量为 `SSLDisable`、`SSLRequire`、
  `SSLVerifyCA`、`SSLVerifyFull`。
- 如果客户端使用签名认证，canonical payload 现在必须包含 HTTP method 和 path：
  `app_id=<appID>&method=<method>&nonce=<nonce>&path=<path>&timestamp=<timestamp>`。
  请求 body 有意不参与 HMAC payload。
- `/mcp` 默认需要 Bearer 认证。只有在明确要开放匿名 MCP surface 时，才设置
  `vef.mcp.require_auth=false`。
- 审批必须通过 `vef.ApprovalModule` 显式启用；它不在默认启动图里。审批事件
  必须路由到 transactional transport，框架订阅者还需要 subscribable sink。
- 审查审批设计器/运行时契约：节点配置在 deploy 时规范化和校验，`PassRule`
  只支持 `all`、`any`、`ratio`，实例详情包含 `formSchema`、`timeline`、
  `flowGraph`，业务绑定改为 opaque `businessRef` 加 engine-owned 状态回写。
- 存储上传使用 multipart 协议（`init_upload`、`upload_part`、`list_parts`、
  `complete_upload`、`abort_upload`）和 `/storage/files/<key>` 代理。公开上传必须
  显式开启 `vef.storage.allow_public_uploads=true`。
- 使用 `vef-cli generate-model-schema` 时，如果模型依赖 Bun 表别名、bare table
  tag、默认表名或 `m2m` 关系，请重新生成 schema；生成器现在更贴近 Bun 的 tag
  解析规则。

## 按版本审计

| 版本 | 需要审查的用户可见变化 |
| --- | --- |
| `v0.27.0` | 多数据源配置与 registry、必填 `primary`、`datasource.Provider`、`Registry.Register` / `Update` / `Unregister` / `Reconcile` / `TestConnection` / `HealthCheck`，以及旧 `vef.data_source` fallback 的移除。 |
| `v0.28.0` | 核心 expression engine、来源 IP 安全检查、MCP read-only query guard、事件 middleware 顺序与 tracing、storage / CRUD API error surface、审批 opt-in module 与路由检查、更严格的请求错误处理。 |
| `v0.29.0` / `v0.29.1` | `Effective*()` 配置默认值、数据库 TLS、按方言加固的 SQL guard、审批空租户 fail-closed、签名 nonce replay-window 修复，以及 `security.UserMenu` JSON 对齐。 |
| `v0.30.0` | 恢复 sequence DB/Redis store、`cache.NewPrefixKeyBuilderWithSeparator`、`logx.LoggerConfigurable`、`reflectx.BreadthFirst`。 |
| `v0.31.0` | expression backend 迁移到 pure-Go `expr-lang`；expression 本身不再要求 `CGO_ENABLED=1`。 |
| `v0.32.x` | 配置 API 重命名（`EnvPrefix`、`CORSConfig`、SSL 常量）、审批模块经加固后恢复，以及兼容 Bun 的 model-schema 生成修正。 |
| `v0.33.x` | 审批详情 DTO 增加 `formSchema`、`timeline`、flow graph progress、action-log metadata，以及 rollback targeting 使用的持久 flow-node id；model-schema 解析改为更忠实的 Bun tag parser。 |
| `v0.34.0` | 通过 `vef.ProvideAuthStrategy` 注册自定义认证策略，以及通过 `api.IPAuth(...)`、`vef.security.ip_whitelists`、`vef.app.trusted_proxies` 使用来源 IP 白名单认证。 |
| `v0.35.0` | 审批实例详情改为基于 engine 记录的 visit trail 和人员快照；条件路由快照宿主 globals；业务绑定使用 opaque ref；审批领域事件改为自描述 envelope。 |

## 数据源与数据库

配置现在使用命名数据源 map：

```toml
[vef.data_sources.primary]
type = "postgres"
host = "127.0.0.1"
port = 5432
user = "postgres"
password = "postgres"
database = "app"
schema = "public"

[vef.data_sources.analytics]
type = "sqlite"
path = "./analytics.db"
```

`primary` 同时是 `config.PrimaryDataSourceName` 和 `datasource.PrimaryName`。
它为全框架 `orm.DB` 注入提供来源，不能被动态 register、update 或 unregister。
当前 runtime provider 支持 `postgres`、`mysql`、`sqlite`；`oracle` 和
`sqlserver` `DBKind` 常量仅为未来 provider 保留。

`datasource.Registry` 是命名运行时数据源的公开扩展点，提供 `Primary`、
`Get`、`Has`、`Names`、`Kind`、`Register`、`Update`、`Unregister`、
`Reconcile`、`TestConnection`、`HealthCheck`。`TestConnection` 会打开一次临时
连接，返回 `datasource.ConnectionInfo.Version`，不会修改 registry。

网络数据库可通过 `ssl_mode` 和 `ssl_root_cert` opt in TLS。可选值是
`disable`、`require`、`verify-ca`、`verify-full`。SQLite 会忽略这些字段。

## API、安全与 MCP

认证策略 registry 现在可通过 `vef.ProvideAuthStrategy(...)` 扩展。内置策略为
`none`、`bearer`、`signature`、`ip`；`api.IPAuth()` 使用 `default` whitelist，
`api.IPAuth("ops")` 对应 `vef.security.ip_whitelists.ops`。

如果服务位于反向代理后面，要配置 `vef.app.trusted_proxies`；否则来自未信任
客户端的 forwarded header 会被忽略，IP 认证看到的是直接连接 peer。

请求解码和认证失败现在更明确：

- malformed params / meta 返回 `api.ErrInvalidRequestParams` /
  `api.ErrInvalidRequestMeta`，语义是 HTTP 400；
- 无效 Bearer token 返回 401，不再退化成 anonymous；
- public operation 使用隔离的 anonymous principal。

MCP database query 是 read-only 且按方言检查；除非设置
`vef.mcp.require_auth=false`，否则 `/mcp` 需要 Bearer 认证。

## 事件总线

应用应按当前 multi-transport 平台理解事件：memory、transactional outbox、
Redis Streams、可选 Inbox dedupe、确定性的 publish/consume middleware 顺序、
W3C trace 传播，以及 `event.RouteInspector` 路由检查。

如果某条路由通过 transactional outbox 服务框架订阅者，还必须包含
`redis_stream` 或 `memory` 这样的 subscribable sink；只有纯发布场景才适合
`["outbox"]`-only 路由。启用 Inbox dedupe 时，retention 必须覆盖 outbox retry
horizon，避免延迟重复投递到来时 dedupe 记录已经被清理。

审批和存储都会发布 transactional events；缺少必需事件路由时会在启动阶段失败。

## 审批

审批现在是可选 feature module。依赖任何 `approval/*` resource、CQRS handler、
binding listener、timeout scanner 或审批事件发布之前，需要把 `vef.ApprovalModule`
加入 `vef.Run(...)`。

重要运行时和 API 变化：

- deploy 阶段会校验并规范化 flow node config；
- 条件使用审批自己的 `expr-lang` evaluator，不使用公开 `expression.Engine`；
- 宿主 globals 由 `approval.InstanceGlobalsResolver` 在实例启动时解析，并快照到
  `Instance.Globals`；
- `FlowGraphNode.Kind` 是 `approval.NodeKind`；
- 实例详情 DTO 暴露 `formSchema`、`timeline`、`flowGraph`；
- `FlowGraphNode.NodeID` 是 action log 和 rollback target 使用的持久 flow-node id；
- 业务绑定使用 opaque `Instance.BusinessRef`；engine 拥有最终状态回写，只通过
  `BusinessRefResolver` 抽取 record id；
- 领域事件使用 `InstanceEventBase`、`TaskEventBase`、`FlowEventBase`，订阅者可以
  直接从 event envelope 做路由和通知，而不必立刻查询审批表；
- tenant 处理对空或歧义 caller context fail-closed。

完整 resource、DTO、事件和扩展点契约见 [审批模块](../modules/approval)。

## Expression 与 Mold

`expression.Engine` 已进入 core boot graph，可注入 API handler。当前 backend 是
pure-Go `expr-lang`；旧的 Goja `js` 包仍然是独立能力。Mold 的 `expr`
transformer 以包含它的 struct 作为环境，并按字段声明顺序执行，所以派生字段可以
引用更早声明的 sibling field。

## 存储

HTTP 层的 storage resource 只支持 multipart 协议。业务侧文件生命周期通过
`storage.Files` / `storage.FilesFor[T]`、claim consumption 和 pending-delete 行
实现，使文件归属变更与业务事务一起提交或回滚。删除重试耗尽后会发出
`vef.storage.delete.dead_letter`，并从队列移除。

默认下载 URL 通过 `/storage/files/<key>` 解析，除非宿主应用提供自定义
`URLKeyMapper`。`pub/` 下的公开 key 匿名可读；私有 key 会调用
`storage.FileACL`。

## 其他公开面

- CRUD mutation response 使用按操作区分的成功消息；CRUD / storage / approval
  对外错误使用带 code 的 `result.Error`。
- result error sentinel 在初始化时冻结 i18n 文本；后续切换语言会影响新的翻译，
  但不会改变已经构造好的 sentinel value。
- `i18n.SetLanguage(...)` 与并发 `T` / `Te` 调用 race-free。
- `sequence.NewDBStore`、`sequence.NewRedisStore` 和可注入的
  `*sequence.MemoryStore` 分别用于持久化、分布式和启动时 seeding。
- `cache.Invalidating[T]` 是 event-invalidated read-through cache wrapper；
  `cache.NewPrefixKeyBuilderWithSeparator` 会按传入值原样保留 separator。
- `logx.LoggerConfigurable[T]` 和 `reflectx.BreadthFirst` 已重新成为公开 utility
  surface 的一部分。
- monitor 模块按唯一分区去重磁盘设备，不会把独立设备错误合并。
- schema inspection 使用 primary data source。
