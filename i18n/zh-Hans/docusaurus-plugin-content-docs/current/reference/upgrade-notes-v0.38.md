---
sidebar_position: 8
---

# 升级说明：v0.38

本页是 `v0.37.0` 到 `v0.38.0`（`v0.37.0..v0.38.0`）后端提交的跨版本审计
地图。该版本聚焦安全（登录加固与有状态 opaque token 会话）、重新设计的
审批业务绑定（复合键 + 持久化期望状态投影）、全新的分布式锁模块、请求
绑定的数字保真，以及容器感知的监控。上一次对照文档核对集成假设停留在
[升级说明：v0.36 / v0.37](../reference/upgrade-notes-v0.37) 的应用，升级时
以本页为纲。

本页不能替代生成的索引。按下面的迁移说明处理后，请对照
[Public API Index](../reference/public-api-index) 和
[Runtime API Index](../reference/runtime-api-index) 核对精确的 Go 符号与
wire 字段。

## 立即检查清单

- **安全——令牌机制**：`token` authenticator type 字符串更名为
  `jwt_token`，`vef.security.token_type` 在 `jwt_token`（默认）与
  `opaque_token` 之间选择。只有已配置机制的认证器会注册：`opaque_token`
  下遗留 JWT 在任何地方（含 MCP）都无法通过认证，`refresh` 操作**不再
  挂载**。`login` 还会拒绝框架签发的令牌类型作为登录凭据。
- **安全——密码传输**：`password.NewCipherEncoder` 已移除。改为注册
  `security.PasswordDecryptor`（例如 `cryptox.NewRSA`）；解密现在发生在
  认证器层，而不是 encoder。
- **安全——锁定**：暴力破解锁定默认开启（`vef.security.lockout`，
  `max_failures = 10`），且从 v0.38 起 `resolve_challenge` 的失败猜测计入
  同一个锁定 key。
- **API——数字保真**：JSON `params`/`meta` 按 `json.Number` 解析。整数
  字段现在拒绝小数/指数形式（`mapx.ErrJSONNumberNotInteger`）与越界值
  （`mapx.ErrJSONNumberOverflow`），不再静默截断；`json.RawMessage` 捕获
  保留精确位数；无类型（`any`）目标仍是 `float64`。
- **API——限流配置**：操作级默认限流（100 次 / 5 分钟）可通过
  `vef.api.rate_limit.max` / `.period` 配置。
- **审批——业务绑定重设计**：`Flow` 上六个扁平 `Business*` 列被单个
  `Flow.BusinessBinding` jsonb 文档取代（`approval.BusinessBindingConfig`：
  `tableName`、复合 `keyColumns`、`statusColumn`、**必填**
  `instanceIdColumn`、可选时间戳列、可选 `statusMapping`），并按部署版本
  快照。写回改为持久化期望状态投影（`apv_business_projection`），一致性
  分 `synchronous`（默认）与 `eventual`
  （`vef.approval.business_binding`）。更新 flow create/update 客户端
  （`params.businessBinding`）与权限种子（新管理端操作的
  `approval.binding.query` / `approval.binding.retry`）。
- **审批——Go 接口**：
  `InstanceLifecycleHook.OnInstanceCompleted(instance, finalStatus)` →
  `OnInstanceTransition(instance, from, to)`（每次迁移都触发）；
  `BusinessRefResolver.ResolveRecordID` →
  `ResolveRecordKey(...) (BusinessRecordKey, error)`；
  `SubscribeInstance` handler 新增第三个参数 `env event.Envelope`；
  `FormSchemaParser.ParseFormFields` 增加了前置 `ctx`。
- **审批——表单 schema**：结构化的 `approval.FormDefinition` 类型已删除。
  `deploy` 以 `params.formSchema`（原 `params.formDefinition`）原样接收
  宿主设计器文档；`form_fields` 在部署时经 `FormSchemaParser` 派生一次。
  详情 API 原样返回 `formSchema`，`my.get_instance_detail` 新增按查看者
  投影的 `fieldPermissions`。
- **审批——迁移**：执行 v0.38 DDL（`auto_migrate` 或自有管线）：
  `apv_flow.business_binding`、`apv_flow_version.business_binding` +
  `form_fields`、`apv_instance.business_projection_id`，以及新表
  `apv_business_projection`。
- **审批——路由**：binding listener 已移除；模块现在只要求 `approval.*`
  有 transactional 路由（只写 `["outbox"]` 即可通过）。只有宿主自己订阅
  （`SubscribeInstance` / 新的 `BindCommand`）时才需要在路由里保留可订阅
  sink。
- **监控**：overview 磁盘摘要只报告根文件系统（`partitions` 恒为 `1`），
  `vef.monitor.excluded_mounts` 已移除；CPU 摘要新增 `effectiveCores`，
  内存指标遵循 cgroup 限额——累加分区或按宿主核数归一化的看板需要调整。
- **移除**：`ptr` 包（改用 `samber/lo` / 内建 `new`；见
  [小工具集](../utilities/small-helpers)）与 snowflake ID 生成器
  （`id.DefaultSnowflakeIDGenerator`、`id.NewSnowflakeIDGenerator`、
  `VEF_NODE_ID` 环境变量）——XID（`id.Generate()`）是唯一的 ID 机制。
- **新增**：注入 `lock.Locker` 获得基于租约的分布式锁
  （`lock.WithLock`、fencing token、按拓扑选择 Redis/内存后端）——见
  [分布式锁](../infrastructure/lock)。

## 按版本审计

| 版本 | 需要核查的用户可见变化 |
| --- | --- |
| `v0.38.0` | 安全：PasswordDecryptor 传输、默认开启的锁定（含挑战计数）、密码策略/历史/过期、`jwt_token` 更名、带并发控制的 opaque token 会话、严格机制门控。API：数字保真、可配置默认限流。审批：复合业务绑定 + 持久投影、迁移型生命周期 hook、`BindCommand`、投递 envelope、opaque 表单 schema + 字段权限。新增 `lock` 包。监控 cgroup/根磁盘改造。移除 `ptr` 与 snowflake；sequence/cron/timex 修复。 |

## 安全

v0.38 把登录路径变成分层、可选启用的加固栈，并新增有状态会话机制。完整
文档：[登录加固](../security/login-hardening) 与
[会话管理](../security/session-management)。

### 破坏性：`jwt_token` 更名与严格机制门控

Bearer 认证器的 type 字符串 `token` 变为 `jwt_token`
（`AuthTypeJWTToken`）；`opaque_token`（`AuthTypeOpaqueToken`）作为有状态
备选由 `vef.security.token_type` 选择。只有已配置机制的认证器会注册——
`opaque_token` 下 `refresh` 操作不挂载，JWT 在包括 MCP 在内的所有入口
失效。`login` 拒绝 `jwt_token` / `opaque_token` / `refresh` 作为登录凭据
类型，关闭"把被窃 access token 洗成 refresh token"的通道。

### 破坏性：移除 `password.NewCipherEncoder`

加密密码传输从 encoder 移到认证器：注册 `security.PasswordDecryptor`，
`PasswordAuthenticator` 在哈希比对前解密。`password.ErrCipherRequired` /
`ErrEncoderRequired` 随类型一起移除。畸形密文会付出一次假 KDF 比对的
代价，关闭时序侧信道。

### opaque token 会话

`OpaqueTokenGenerator` + `SessionStore`（默认内存；`fx.Decorate` 换
Redis）提供绝对 `max_lifetime` 上限内的滑动空闲 TTL（签发时同样钳制）、
按账号并发限制（`max_concurrent`，`on_exceed` 默认 `evict_oldest` /
可选 `reject`）、登出吊销、会话管理面（`ListByUser` / `Revoke` /
`RevokeUser`）以及可选的 `SessionInspector.ListAll`。内存存储构建在 TTL
cache 上（后台 GC、不无限增长）；Redis 存储的多键变更原子执行、按用户
集合可自愈、`ListAll` 走 keyspace `SCAN`。

### 登录加固层

`security.LoginGuard`（默认开启的锁定，`lock`/`backoff` 策略，
`user`/`ip`/`user_ip` key，存储出错 fail-open，HTTP 429
`ErrAccountLocked`）、配置驱动的密码强度
（`vef.security.password_policy`，可组合 `PasswordRule`）、密码历史
（`PasswordHistoryStore` + `history_depth`）与密码过期
（`PasswordMetadataLoader` + `max_age` + `NewExpiryPasswordChangeChecker`）。
字符类规则不把无大小写字母（中文）计入任何类别，身份规则按 rune 计数。

## API

### 破坏性：请求绑定的数字保真

`api.Params` / `api.Meta` 以 `json.Decoder.UseNumber` 反序列化。有类型的
数值字段按精确位数解析、严格性对齐 `encoding/json`（新 sentinel
`mapx.ErrJSONNumberNotInteger`、`mapx.ErrJSONNumberOverflow`）；
`json.RawMessage` 字段保留完整精度；`any` 目标仍收到 `float64`。依赖
越界或小数值被静默截断的客户端现在会收到错误——这正是修复本身，不是
回归。

### 可配置的默认限流

`vef.api.rate_limit.max` / `.period` 设置未声明自己
`OperationSpec.RateLimit` 的操作的默认值（默认 100 / 5m，按操作 ×
客户端计 key，按节点计数）。

## 审批

完整文档：[事件与集成](../approval/integration)、
[流程设计](../approval/flow-design)、[RPC 资源](../approval/resources)。

### 破坏性：复合业务绑定与持久投影

单个 `BusinessBindingConfig` 文档取代 `Flow` 上六个扁平列；`KeyColumns`
支持复合键并对照真实的非空主键/唯一键校验；`InstanceIDColumn` 作为 CAS
防护栏必填；`StatusMapping` 把状态翻译为宿主词汇。绑定按部署版本快照。
状态经 `apv_business_projection` 收敛（发起时认领，目标被占返回
`ErrBindingTargetBusy`；每次迁移登记期望 revision；`synchronous` 或
`eventual` 应用，后者由带租约的 worker 指数退避重试）。
`InstanceBindingFailedEvent` 只在 eventual 模式作为运维通知。新管理端
操作：`find_business_projections`、`retry_business_projection`；
`get_metrics` 新增 `businessProjectionCounts` /
`pendingBusinessProjections`。新错误码 `40017`–`40020` 与
`40107`–`40110`。

### 破坏性：泛化生命周期 hook 与订阅 envelope

`OnInstanceTransition(instance, from, to)` 在每次状态迁移的事务内触发
（多个 hook 顺序未定义）；`SubscribeInstance` handler 收到投递
`event.Envelope`（`Envelope.ID` 是稳定去重键）；`BindCommand[E, C]` 把
实例事件桥接为 CQRS 命令，consumer group 由命令身份派生
（`vef:cmd:...`），前置守卫 `ErrNonCommandAction` /
`ErrUnnamedCommandType`。

### 破坏性：opaque 表单 schema 与字段权限

`FlowVersion.FormSchema` 是原样存储的宿主设计器文档
（`json.RawMessage`）；`FlowVersion.FormFields` 在部署时经上下文感知的
`FormSchemaParser` 派生一次（用 `vef.ProvideApprovalFormSchemaParser`
替换）。`FormDefinition` 包装类型已删除，`deploy` 参数更名为
`formSchema`。字段权限在部署时校验（key 对照字段、CC 子集、`required`
与 `auto_pass` 超时互斥）、在写路径强制（只合并 `editable`/`required`；
`required` 仅在 approve/handle 检查；无权限节点丢弃提交数据）、在
`my.get_instance_detail` 按查看者投影（`fieldPermissions`，max-merge
格，fail-closed）。

## 新增：分布式锁

`lock` 包提供基于租约的锁：按拓扑选择的 DI 默认实现（Redis 启用时
`RedisLocker`，否则 `MemoryLocker` + 启动警告）、`WithLock`（自动续约、
租约丢失取消、panic 安全释放）、fencing token，以及对重试幂等的 Redis
脚本。见[分布式锁](../infrastructure/lock)。

## 监控

- **破坏性**：overview 磁盘摘要只报告根文件系统；
  `vef.monitor.excluded_mounts` 已移除。
- CPU 摘要/详情新增 `effectiveCores`（感知 cgroup v1/v2 配额的使用率
  归一化）；内存头部指标遵循 cgroup 限额；CPU 总使用率由每核心采样推导，
  采样器生命周期得到加固。

## 移除与小修复

- **移除 `ptr` 包**——迁移对照表见[小工具集](../utilities/small-helpers)。
  （`ptr.Of` 在包被整体删除前曾短暂改为总是返回指针。）
- **移除 snowflake 生成器**（`id.DefaultSnowflakeIDGenerator`、
  `NewSnowflakeIDGenerator`、`config.EnvNodeID` / `VEF_NODE_ID`）；
  `id.Generate()`（XID）是框架的 ID 机制。
- `cron.Scheduler.Update` 保留 job identifier，你持有的句柄在更新后仍然
  有效。
- `sequence` 以 `sequence.ErrInvalidStep` 拒绝 `SeqStep < 1`，
  `MemoryStore.Register` 与进行中的 `Reserve` 串行化。
- `timex` 宽松 `Parse*` 的回退路径改为按本地时区解析无时区输入，与主
  路径一致。
- `cache.NewInvalidating` 透传 `MemoryOption`（用 `WithMemMaxSize` 给
  read-through 缓存设界）；mold 的缓存字典解析器 LRU 上限 4096。
- `schema.UniqueKey` 新增 `predicate` / `hasExpressions`；
  `schema.ErrTableMissing` 是表缺失的新 Go 级 sentinel。
- `expression` 对 JSON 原生环境直接透传、不再重复序列化（仅性能）。
- `event` 文档澄清 `Ordered` 语义：顺序指交付顺序；`WithConcurrency > 1`
  会使 handler 执行交错。
