---
sidebar_position: 9
---

# 升级说明：v0.39

本页是 `v0.38.0` 到 `v0.39.0` 及其后续修复的跨版本审计地图。本次发布的
核心是三个新的可选模块——[集成引擎](../integration/overview)、
[持久化调度存储](../infrastructure/cron-store)与
[WebSocket 服务端推送](../infrastructure/push)——外加重建的 JS 引擎、新的
出站 HTTP 客户端、SQL Server / Oracle 连接支持，以及一组破坏性重命名。

它不能替代生成的索引页。按下方迁移说明操作后，请对照
[公开 API 索引](./public-api-index)核对精确的 Go 符号与线上字段。

## 立查清单（破坏性）

- **`httpx` → `fiberx` 重命名**：Fiber 请求辅助（`IsJSON`、`IsMultipart`、
  `GetIP`）从 `github.com/coldsmirk/vef-framework-go/httpx` 移到
  `.../fiberx`。`httpx` 导入路径现在承载新的
  [出站 HTTP 客户端](../utilities/httpx)。机械替换导入即可。
- **JS 引擎重建**：`js.New()` 已删除。用 `js.NewEngine(...)` 构建共享
  `js.Engine`，用 `engine.NewRuntime(...)` 产出按执行的运行时；以
  `rt.RunString(ctx, src)` / `rt.RunProgram(ctx, prog)` 执行（context
  先行、支持取消）。vendored 的 `dayjs`/`Big`/`utils`/`validator` 全局被
  单一 stdlib bundle 取代：`BigNumber`、`dayjs`、`fxp`、`radashi`、`z`、
  `URL` / `URLSearchParams`。见 [JS 引擎](../data-tools/js-engine)。
- **JS SQL 绑定重命名**：脚本侧 `sql` 动词现为 `sql.queryList` /
  `sql.queryOne` / `sql.execute`；Go 选项为 `jssql.WithExecute()`，哨兵为
  `jssql.ErrExecuteDisabled`。`sql.queryList` 无匹配行时现在返回 `[]`
  （不再是 `null`）。
- **`mold` dictionary → code set 重命名**：标签前缀 `dict:` 变为
  `codes:`；`Dictionary*` 标识符变为 `CodeSet*`
  （`CachedCodeSetResolver`、`CodeSetChangedEvent`、
  `PublishCodeSetChangedEvent` 等）。见
  [重命名对照表](../data-tools/mold#v038--v039-重命名对照)。
- **`cryptox.NewSM4` 默认 GCM**：早期版本的 SM4 密文是 CBC 的；解密请显式
  使用 `cryptox.NewSM4(key, cryptox.WithSM4Mode(cryptox.Sm4ModeCbc))`。
  AES 此前已默认 GCM；两个算法现在对齐。
- **安全——保留身份强制**：认证器、挑战提供者与令牌签发现在拒绝声称框架
  保留身份的主体（`security.Principal.IsReserved()`：`system` 类型、
  `orm.OperatorSystem`、`orm.OperatorCronJob`）。请审计可能产生此类 ID 的
  自定义 `Authenticator` / `UserLoader`。
- **API——权限声明校验**：矛盾或畸形的 `RequiredPermission` 声明在资源
  注册时被拒绝（启动即失败）。权限令牌遵循点分约定
  （`domain.entity.action`）。
- **审批**：
  - `approval/category.find_tree_options` 已删除。
  - `approval/flow.update` 参数类型 `UpdateParams` → `UpdateFlowParams`
    （线上格式不变）。
  - `find_versions` 现在返回不含图文档的版本摘要；具体版本经 `get_graph`
    的新参数 `params.versionId` 获取。
  - 错误码 `40015`（流程绑定锁）退役且永不复用。

## 新模块（Opt-In）

- **集成引擎**（`vef.IntegrationModule`）：契约、系统、适配器（出站 +
  入站）、路由、按系统码值映射、调用日志与统计、dry-run 测试台，以及入站
  HTTP 网关（`POST /integration/inbound/:systemCode/:contractCode`）。
  配置位于 `vef.integration`。入口见
  [集成引擎](../integration/overview)。
- **持久化调度**（`vef.cron.store.enabled = true`）：持久化调度、集群内
  单次触发、错触/并发策略、运行流水账、崩溃恢复，以及
  `sys/cron/schedule` / `sys/cron/run` 资源。用
  `vef.ProvideCronJobHandler` 注册处理器。见
  [持久化调度](../infrastructure/cron-store)。
- **服务端推送**（`vef.push.enabled = true`）：WebSocket 端点（默认
  `/ws`）、面向用户/角色/广播的 `push.Notifier`、会话吊销集成、Redis
  多节点自动中继。见[服务端推送](../infrastructure/push)。

## 新能力（非破坏性）

- **`httpx` 出站 HTTP 客户端**：不可变的按上游客户端 + 流式按调用请求
  构建——base URL、默认头/查询、带抖动退避的重试（幂等方法；
  429/502/503/504）、钩子、代理/TLS/重定向/响应体控制。见
  [httpx](../utilities/httpx)。
- **JS 能力库**：`jshttp`（fetch 风格 HTTP）、`jssql`（受守卫的 SQL）、
  `jscache`、`jsevents`、`jscrypto`、`jsconsole` —— 经引擎目录按运行时
  opt-in。
- **SQL Server 与 Oracle 数据源**：`type = "sqlserver"`（默认端口 1433）
  与 `type = "oracle"`（`database` 填服务名，默认端口 1521）成为真实连接
  提供者，只读 SQL 守卫解析器同步覆盖。见
  [数据源](../data-access/datasources)。
- **安全认证策略**：`api.APIKeyAuth(...)`（默认读 `X-API-Key` 头，配置于
  `vef.security.api_keys`）与 `api.HTTPBasicAuth()`（配置于
  `vef.security.basic_accounts`），加载器均可替换。见
  [认证](../security/authentication)。
- **会话吊销接缝**：`security.SessionRevocationListener` /
  `SessionRevocationNotifier`（推送模块用它即时踢掉死会话）。见
  [会话管理](../security/session-management)。
- **审批**：宿主自有的流程 `labels` 及相等过滤（`find_flows`、
  `my.find_available_flows`）、实例详情披露 labels、自助
  `my.get_start_form`，以及 `my.get_instance_detail` 中的办理人任务上下文
  （含回退目标与可移除办理人的 `myTask`）。见
  [审批 RPC 资源](../approval/resources)。
- **监控**：`sys/monitor.get_integration_stats` 报告本节点集成调用统计。
- **timex**：`DateTime.AsLocal()` 在做与 `Now` 的时刻运算前，把朴素墙钟
  字段重解释到进程本地时区。

## 行为说明

- 集成密文静态存储：配置 `vef.integration.secret_key`（默认 AES-GCM；
  `secret_algorithm = "sm4"` 用 SM4-GCM）。切换算法需重新录入已存密文。
- 调度存储时间：调用方提供的调度时间在持久化前归一化到本地墙钟；未指定
  时区的 cron 触发器在 UTC 求值，绝不使用节点进程本地时区。
- 推送中继：由 `vef.redis.enabled = true` 隐式启用；要求设置
  `vef.app.name` 作为频道命名空间。
- `security/auth` 登录流程继续拒绝框架签发的令牌类型作为凭证；保留身份
  门禁是其上的追加加固。

## 另请参阅

- [升级说明：v0.38](./upgrade-notes-v0.38) —— 上一个版本
- [配置参考](./configuration-reference) —— 新增的 `vef.cron`、
  `vef.integration`、`vef.push` 章节
- [内置资源](./built-in-resources) —— 扩充后的资源索引
