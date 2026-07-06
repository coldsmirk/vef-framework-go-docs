---
sidebar_position: 3
---

# 审批模块

`approval` 模块提供完整的工作流引擎，用于构建基于审批的业务流程。支持可视化流程设计（兼容 React Flow）、多级审批链、条件分支、并行审批、委托、回退和事务性事件发布。

## 启用模块

审批是一个可选功能模块。它有意不包含在默认 `vef.Run(...)` boot graph 中，
所以不需要审批工作流的应用不会注册它的 API resources、CQRS handlers、
engine、binding listener 或 timeout scanners。

需要时显式启用：

```go
vef.Run(
    vef.ApprovalModule,
    app.Module,
)
```

审批会通过 `event.WithTx` 发布 `approval.*` 事件，它的 binding listener
也会订阅这些事件。宿主应用必须把 `approval.*` 路由到一个带有可订阅 sink
的 transactional transport，例如 sink 为 Redis Streams 的 outbox 路由。

`InstanceBindingFailedEvent` 是 transactional-route 启动检查的例外：它由异步
binding listener 在审批事务已经提交后发出。`InstanceCompletedEvent` 的路由要求
最严格，因为 binding listener 会订阅它；路由必须在 transactional outbox 之外
同时包含可订阅 sink transport，例如 `memory` 或 `redis_stream`。

## RPC 资源

启用模块后，会注册以下 RPC 资源。这些操作都不是公开接口：调用者必须已认证；
表中列出 `RequiredPermission` 时还会校验对应权限点。

RPC 调用使用 [API](../guide/api) 里的标准 envelope：`resource`、`action`、
`version`、`params` 和 `meta`。下面标成 `params.*` 的字段从 `params` 解码；
标成 `meta.*` 的字段从 `meta` 解码。生成的
[运行时 API 索引](../reference/runtime-api-index) 包含每个请求/响应 DTO 的
完整 JSON 字段 ledger。

Grouped-family audit 还固定了 1059 grouped approval field/method entries：
930 approval package entries、62 approval/admin DTO field entries、67
approval/my DTO field entries。这些 entries 覆盖公开 Go DTO fields、domain
event fields、node-data helpers、lifecycle hooks、resolver interfaces 和
status helper methods；精确签名列在 public API index 中，verifier 会锁定排序后的
签名以及 receiver/type 分布。

### `approval/category`

| Action | 权限 | 参数 | 说明 |
| --- | --- | --- | --- |
| `find_tree` | `approval.category.query` | `CategorySearch` | 按租户过滤的树查询 |
| `find_tree_options` | `approval.category.query` | `CategorySearch` + `DataOptionConfig` | 按租户过滤的树形选项 |
| `create` | `approval.category.create` | `CategoryParams` | 非 super-admin 的租户由调用者覆盖写入 |
| `update` | `approval.category.update` | `CategoryParams` | 非 super-admin 只能修改自己租户的数据 |
| `delete` | `approval.category.delete` | 主键参数 | 非 super-admin 只能删除自己租户的数据 |

| Action | 请求字段 |
| --- | --- |
| `find_tree` | `meta.name`、`meta.isActive`、`meta.sort` |
| `find_tree_options` | `meta.name`、`meta.isActive`、`meta.sort`，以及 option 映射 metadata：`meta.labelColumn`、`meta.valueColumn`、`meta.descriptionColumn`、`meta.metaColumns` |
| `create` | `params.id`、`params.tenantId` 必填、`params.code` 必填、`params.name` 必填、`params.icon`、`params.parentId`、`params.sortOrder`、`params.isActive`、`params.remark` |
| `update` | `params.id`、`params.tenantId` 必填、`params.code` 必填、`params.name` 必填、`params.icon`、`params.parentId`、`params.sortOrder`、`params.isActive`、`params.remark` |
| `delete` | `params.id` 必填 |

### `approval/delegation`

| Action | 权限 | 参数 | 说明 |
| --- | --- | --- | --- |
| `find_page` | `approval.delegation.query` | `DelegationSearch` + pageable meta | 非 super-admin 只查询自己的委托 |
| `create` | `approval.delegation.create` | `DelegationParams` | 非 super-admin 的 delegator 由调用者覆盖写入 |
| `update` | `approval.delegation.update` | `DelegationParams` | 非 super-admin 不能转移委托所有权 |
| `delete` | `approval.delegation.delete` | 主键参数 | 非 super-admin 按委托所有者限制 |

| Action | 请求字段 |
| --- | --- |
| `find_page` | `meta.delegatorId`、`meta.delegateeId`、`meta.isActive`、`meta.sort`、`meta.page`、`meta.size` |
| `create` | `params.id`、`params.delegatorId` 必填、`params.delegateeId` 必填、`params.flowCategoryId`、`params.flowId`、`params.startTime`、`params.endTime`、`params.isActive`、`params.reason` |
| `update` | `params.id`、`params.delegatorId` 必填、`params.delegateeId` 必填、`params.flowCategoryId`、`params.flowId`、`params.startTime`、`params.endTime`、`params.isActive`、`params.reason` |
| `delete` | `params.id` 必填 |

### `approval/flow`

| Action | 权限 | 参数 | 说明 |
| --- | --- | --- | --- |
| `create` | `approval.flow.create` | `CreateFlowParams` | 开启审计 |
| `deploy` | `approval.flow.deploy` | `DeployFlowParams` | 开启审计 |
| `publish_version` | `approval.flow.publish` | `PublishVersionParams` | 开启审计 |
| `update_flow` | `approval.flow.update` | `UpdateFlowParams` | 开启审计 |
| `toggle_active` | `approval.flow.update` | `ToggleActiveParams` | 开启审计 |
| `get_graph` | `approval.flow.query` | `GetGraphParams` | 读取已发布流程图 |
| `find_flows` | `approval.flow.query` | `FindFlowsParams` | 分页查询流程 |
| `find_initiators` | `approval.flow.query` | `FindInitiatorsParams` | 查询发起人配置 |
| `find_versions` | `approval.flow.query` | `FindVersionsParams` | 查询单个流程的版本列表 |

| Action | 请求字段 |
| --- | --- |
| `create` | `params.tenantId` 必填、`params.code` 必填、`params.name` 必填、`params.categoryId` 必填、`params.bindingMode` 必填、`params.icon`、`params.description`、`params.businessTable`、`params.businessPkField`、`params.businessStatusField`、`params.adminUserIds`、`params.isAllInitiationAllowed`、`params.instanceTitleTemplate`、`params.initiators` |
| `deploy` | `params.flowId` 必填、`params.description`、`params.flowDefinition` 必填、`params.formDefinition`、`params.storageMode` |
| `publish_version` | `params.versionId` 必填 |
| `update_flow` | `params.flowId` 必填、`params.name` 必填、`params.bindingMode` 必填、`params.instanceTitleTemplate` 必填、`params.icon`、`params.description`、`params.businessTable`、`params.businessPkField`、`params.businessStatusField`、`params.adminUserIds`、`params.isAllInitiationAllowed`、`params.initiators` |
| `toggle_active` | `params.flowId` 必填、`params.isActive` |
| `get_graph` | `params.flowId` 必填、`params.tenantId` |
| `find_flows` | `params.tenantId`、`params.categoryId`、`params.keyword`、`params.isActive`、`params.page`、`params.pageSize` |
| `find_initiators` | `params.flowId` 必填、`params.tenantId` |
| `find_versions` | `params.flowId` 必填、`params.tenantId` |

`params.initiators` 条目使用 `kind`（`user`、`role` 或 `department`）和
`ids`。业务绑定场景下，`businessTable`、`businessPkField` 和
`businessStatusField` 都是 SQL identifier，会经过
`ValidateBusinessIdentifier` 校验。

### `approval/instance`

| Action | 权限 | 参数 | 说明 |
| --- | --- | --- | --- |
| `start` | `approval.instance.start` | `StartInstanceParams` | 开启审计 |
| `process_task` | `approval.task.process` | `ProcessTaskParams` | 开启审计；`action` 必须是 `approve`、`reject`、`transfer`、`rollback` 或 `handle` |
| `withdraw` | `approval.instance.withdraw` | `WithdrawParams` | 开启审计 |
| `resubmit` | `approval.instance.resubmit` | `ResubmitParams` | 开启审计；接受已退回或已撤回的实例 |
| `add_cc` | `approval.instance.cc` | `AddCCParams` | 开启审计 |
| `mark_cc_read` | `approval.instance.cc` | `MarkCCReadParams` | 读回执，不开启审计 |
| `add_assignee` | `approval.task.add_assignee` | `AddAssigneeParams` | 开启审计；`addType` 为 `before`、`after` 或 `parallel` |
| `remove_assignee` | `approval.task.remove_assignee` | `RemoveAssigneeParams` | 开启审计 |
| `urge_task` | `approval.task.urge` | `UrgeTaskParams` | 额外限流：`1m` 内最多 `10` 次 |

| Action | 请求字段 |
| --- | --- |
| `start` | `params.tenantId` 必填、`params.flowCode` 必填、`params.businessRef` 最多 512 字符、`params.formData` |
| `process_task` | `params.taskId` 必填、`params.action` 必填（`approve`、`reject`、`transfer`、`rollback` 或 `handle`）、`params.opinion` 最多 2000 字符、`params.formData`、`params.attachments` 最多 20 项，每项最多 512 字符、`params.transferToId`、`params.targetNodeId` |
| `withdraw` | `params.instanceId` 必填、`params.reason` 最多 2000 字符 |
| `resubmit` | `params.instanceId` 必填、`params.formData` |
| `add_cc` | `params.instanceId` 必填、`params.ccUserIds` 必填，1-50 个 ID |
| `mark_cc_read` | `params.instanceId` 必填 |
| `add_assignee` | `params.taskId` 必填、`params.userIds` 必填，1-50 个 ID，`params.addType` 必填（`before`、`after` 或 `parallel`） |
| `remove_assignee` | `params.taskId` 必填 |
| `urge_task` | `params.taskId` 必填、`params.message` 最多 500 字符 |

### `approval/my`

自助查询不声明 `RequiredPermission`，但仍要求当前调用者已认证。

| Action | 参数 | 输出 |
| --- | --- | --- |
| `find_available_flows` | `FindAvailableFlowsParams` | `page.Page[my.AvailableFlow]` |
| `find_initiated` | `FindInitiatedParams` | `page.Page[my.InitiatedInstance]` |
| `find_pending_tasks` | `FindPendingTasksParams` | `page.Page[my.PendingTask]` |
| `find_completed_tasks` | `FindCompletedTasksParams` | `page.Page[my.CompletedTask]` |
| `find_cc_records` | `FindCCRecordsParams` | `page.Page[my.CCRecord]` |
| `get_pending_counts` | `GetPendingCountsParams` | `my.PendingCounts` |
| `get_instance_detail` | `GetInstanceDetailParams` | `my.InstanceDetail` |

| Action | 请求字段 |
| --- | --- |
| `find_available_flows` | `params.tenantId`、`params.keyword`、`params.page`、`params.pageSize` |
| `find_initiated` | `params.tenantId`、`params.status`、`params.keyword`、`params.page`、`params.pageSize` |
| `find_pending_tasks` | `params.tenantId`、`params.page`、`params.pageSize` |
| `find_completed_tasks` | `params.tenantId`、`params.page`、`params.pageSize` |
| `find_cc_records` | `params.tenantId`、`params.isRead`、`params.page`、`params.pageSize` |
| `get_pending_counts` | `params.tenantId` |
| `get_instance_detail` | `params.instanceId` 必填 |

`my.InstanceDetail` 的 JSON payload 包含 `instance`、`formSchema`、`timeline`、
`flowGraph` 和 `availableActions`。已提交表单数据在 `instance.formData` 中；
业务绑定流程的 opaque 业务引用在 `instance.businessRef` 中。

`availableActions` 是查询层给 UI 的提示。对申请人来说，实例可以流转到
`withdrawn` 时包含 `withdraw`，实例已退回或已撤回时包含
`resubmit`。对 pending task 来说，handle 节点包含 `handle`，其他节点
包含 `approve`，随后包含 `reject`，并按当前节点开关追加 `transfer`、
`rollback`、`add_assignee` 或 `add_cc`。只要实例存在任何 pending task，
还会包含 `urge`。命令 handler 仍会独立做最终校验。

### `approval/admin`

| Action | 权限 | 参数 | 说明 |
| --- | --- | --- | --- |
| `find_instances` | `approval.instance.query` | `AdminFindInstancesParams` | 非 super-admin 按租户过滤 |
| `find_tasks` | `approval.task.query` | `AdminFindTasksParams` | 非 super-admin 按租户过滤 |
| `get_instance_detail` | `approval.instance.detail` | `AdminGetInstanceDetailParams` | 完整管理端详情 |
| `find_action_logs` | `approval.log.query` | `AdminFindActionLogsParams` | 要求 `instanceId` |
| `get_metrics` | `approval.metrics.query` | `AdminGetMetricsParams` | 聚合指标 |
| `terminate_instance` | `approval.instance.terminate` | `AdminTerminateInstanceParams` | 开启审计 |
| `reassign_task` | `approval.task.reassign` | `AdminReassignTaskParams` | 开启审计 |

| Action | 请求字段 |
| --- | --- |
| `find_instances` | `params.tenantId`、`params.applicantId`、`params.status`、`params.flowId`、`params.keyword`、`params.page`、`params.pageSize` |
| `find_tasks` | `params.tenantId`、`params.assigneeId`、`params.instanceId`、`params.status`、`params.page`、`params.pageSize` |
| `get_instance_detail` | `params.instanceId` 必填 |
| `find_action_logs` | `params.instanceId` 必填、`params.tenantId`、`params.page`、`params.pageSize` |
| `get_metrics` | `params.tenantId` |
| `terminate_instance` | `params.instanceId` 必填、`params.reason` 最多 2000 字符 |
| `reassign_task` | `params.taskId` 必填、`params.newAssigneeId` 必填、`params.reason` 最多 2000 字符 |

管理端列表和指标查询中，非 super-admin 调用者会忽略提交的 `tenantId` override，
并被过滤到自己的租户。super-admin 可以传 `tenantId` 只看一个租户，也可以省略它
获得跨租户可见性。

管理端列表/详情 DTO 在线上的租户字段名保持为 `tenantId`。原始审计日志仍可通过
`find_action_logs` 单独分页查询；详情响应携带 timeline 和 flow graph 投影供 UI 渲染。

### 响应 DTO 字段

管理端响应使用 `approval/admin` 包中的 DTO：

| DTO | JSON 字段 |
| --- | --- |
| `admin.Instance` | `instanceId`、`instanceNo`、`title`、`tenantId`、`flowId`、`flowName`、`applicantId`、`applicantName`、`status`、`currentNodeName`、`createdAt`、`finishedAt` |
| `admin.Task` | `taskId`、`instanceId`、`instanceTitle`、`flowName`、`nodeName`、`assigneeId`、`assigneeName`、`status`、`createdAt`、`deadline`、`finishedAt` |
| `admin.InstanceDetail` | `instance`、`formSchema`、`timeline`、`flowGraph` |
| `admin.InstanceDetailInfo` | `instanceId`、`instanceNo`、`title`、`tenantId`、`flowId`、`flowName`、`flowVersionId`、`applicant`、`status`、`currentNodeId`、`currentNodeName`、`businessRef`、`formData`、`createdAt`、`finishedAt` |
| `admin.ActionLog` | `logId`、`action`、`nodeId`、`taskId`、`operator`、`transferTo`、`rollbackToNodeId`、`addedAssignees`、`removedAssignees`、`ccUsers`、`opinion`、`attachments`、`createdAt` |
| `admin.Metrics` | `tenantId`、`capturedAt`、`instanceCounts`、`taskCounts`、`timeoutTaskCount`、`avgCompletionSeconds`、`pendingBindingFailures` |

自助响应使用 `approval/my` 包中的 DTO：

| DTO | JSON 字段 |
| --- | --- |
| `my.AvailableFlow` | `flowId`、`flowCode`、`flowName`、`flowIcon`、`description`、`categoryId`、`categoryName` |
| `my.InitiatedInstance` | `instanceId`、`instanceNo`、`title`、`flowName`、`flowIcon`、`status`、`currentNodeName`、`createdAt`、`finishedAt` |
| `my.PendingTask` | `taskId`、`instanceId`、`instanceTitle`、`instanceNo`、`flowName`、`flowIcon`、`applicantName`、`nodeName`、`createdAt`、`deadline`、`isTimeout` |
| `my.CompletedTask` | `taskId`、`instanceId`、`instanceTitle`、`instanceNo`、`flowName`、`flowIcon`、`applicantName`、`nodeName`、`status`、`finishedAt` |
| `my.CCRecord` | `ccRecordId`、`instanceId`、`instanceTitle`、`instanceNo`、`flowName`、`flowIcon`、`applicantName`、`nodeName`、`isRead`、`createdAt` |
| `my.PendingCounts` | `pendingTaskCount`、`unreadCcCount` |
| `my.InstanceDetail` | `instance`、`formSchema`、`timeline`、`flowGraph`、`availableActions` |
| `my.InstanceInfo` | `instanceId`、`instanceNo`、`title`、`flowName`、`flowIcon`、`applicant`、`status`、`currentNodeId`、`currentNodeName`、`businessRef`、`formData`、`createdAt`、`finishedAt` |

## 架构概览

```
流程分类 → 流程 → 流程版本 → 节点 + 边
                               ↓
                           实例 → 任务 → 操作日志
```

| 概念 | 数据表 | 说明 |
| --- | --- | --- |
| 流程分类 | `apv_flow_category` | 流程的层级分组 |
| 流程 | `apv_flow` | 工作流定义（如"请假申请"）|
| 流程版本 | `apv_flow_version` | 版本快照，包含节点、边和表单 schema |
| 流程节点 | `apv_flow_node` | 工作流中的一个步骤 |
| 流程边 | `apv_flow_edge` | 节点之间的有向连接 |
| 实例 | `apv_instance` | 流程的运行实例 |
| 任务 | `apv_task` | 分配给用户的审批/办理任务 |
| 操作日志 | `apv_action_log` | 所有操作的审计追踪 |

## 配置

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

`auto_migrate` 是普通 boolean 开关，不会由 `ApprovalConfig.ApplyDefaults()`
自动设为 true；需要启动时执行 approval DDL 时必须显式开启。
`cc_record_retention` 只清理已经读过的 CC 记录。

> 老版本里归属 `[vef.approval]` 的 `outbox_relay_interval` / `outbox_max_retries` / `outbox_batch_size` 已在 v0.21 迁移至 `[vef.event.transports.outbox]`，由全框架统一的 outbox transport 服务所有模块——参考 [事件总线](../features/event-bus)。审批模块的 binding listener 和 outbox 发布两侧都会在启动时通过 `event.RouteInspector` 断言路由，路由配错时应用直接启动失败，而不是静默降级。

详见[配置参考](../reference/configuration-reference)。

## 绑定模式

| 模式 | 常量 | Wire value | 说明 |
| --- | --- | --- | --- |
| 独立 | `BindingStandalone` | `standalone` | 表单数据存储在审批模块自有表中 |
| 业务 | `BindingBusiness` | `business` | 关联到已有的业务数据表 |

业务绑定通过 `BusinessTable`、`BusinessPkField` 和 `BusinessStatusField` 将审批流程与业务表关联。

## 节点类型

| 节点类型 | 常量 | Wire value | 说明 |
| --- | --- | --- | --- |
| 开始 | `NodeStart` | `start` | 工作流入口 |
| 审批 | `NodeApproval` | `approval` | 需要审批人执行审批动作 |
| 办理 | `NodeHandle` | `handle` | 需要处理人执行办理动作 |
| 条件 | `NodeCondition` | `condition` | 基于条件进行分支 |
| 抄送 | `NodeCC` | `cc` | 向指定用户发送通知 |
| 结束 | `NodeEnd` | `end` | 工作流终点 |

## 条件分支

条件节点会按 priority 顺序评估 `ConditionBranch`。每个分支包含一个或多个
`ConditionGroup`：同一个 group 内的条件按 AND 组合，同一分支上的多个
group 按 OR 组合。

`ConditionField` 使用结构化的 `Subject` / `Operator` / `Value` 字段。
`Operator` 类型是 `ConditionOperator`；公开常量包括 `OperatorEquals`、
`OperatorNotEquals`、`OperatorGreater`、`OperatorGreaterOrEq`、`OperatorLess`、
`OperatorLessOrEq`、`OperatorIn`、`OperatorNotIn`、`OperatorContains`、
`OperatorNotContains`、`OperatorStartsWith`、`OperatorEndsWith`、
`OperatorIsEmpty` 和 `OperatorIsNotEmpty`。内置 evaluator 会把 field condition
转换成 `expr-lang` 表达式；未知 operator 会被转换成恒为 `false` 的表达式。

`ConditionExpression` 直接用 `expr-lang` 执行原始 `Expression` 字符串。
评估环境暴露：

| 名称 | 值 |
| --- | --- |
| `formData` | 当前实例的 `FormData` map |
| `applicantId` | 当前申请人 ID |
| `applicantDepartmentId` | 申请人部门 ID；不存在时是 `""` |
| globals | 宿主解析出的 `Instance.Globals`，作为顶层 binding 暴露 |

审批条件目前有意直接使用 `expr-lang`，不走公开的 `expression.Engine`
wrapper。这样工作流条件语义由审批 evaluator 固定，不依赖 expression
module 的装配方式。

宿主应用可以实现 `approval.InstanceGlobalsResolver`，在实例启动时根据已认证
principal 解析全局变量。该快照会持久化到 `Instance.Globals`；客户端不能在
`start` 请求体里提交它。Field condition 会先从 globals 解析 `Subject`，再查
`formData`；expression condition 会把 globals 暴露为顶层 binding，但内置的
`formData`、`applicantId`、`applicantDepartmentId` 名称发生冲突时优先级更高。

## 审批方式

当节点有多个审批人时：

| 方式 | 常量 | Wire value | 行为 |
| --- | --- | --- | --- |
| 顺序 | `ApprovalSequential` | `sequential` | 审批人按顺序逐个处理 |
| 并行 | `ApprovalParallel` | `parallel` | 审批人同时处理 |

枚举类型是 `ApprovalMethod`。

### 通过规则（并行模式）

| 规则 | 常量 | Wire value | 行为 |
| --- | --- | --- | --- |
| 全部 | `PassAll` | `all` | 所有审批人必须同意 |
| 任意 | `PassAny` | `any` | 至少一人同意即通过 |
| 比例 | `PassRatio` | `ratio` | 达到一定比例即通过 |

自定义通过规则实现使用 `PassRuleStrategy`、`PassRuleContext`，并返回
`PassRuleResult`（`PassRulePending`、`PassRulePassed`、`PassRuleRejected`）。

## 审批人类型

| 类型 | 常量 | Wire value | 说明 |
| --- | --- | --- | --- |
| 指定用户 | `AssigneeUser` | `user` | 特定用户 |
| 角色 | `AssigneeRole` | `role` | 拥有某角色的用户 |
| 部门 | `AssigneeDepartment` | `department` | 部门负责人 |
| 申请人本人 | `AssigneeSelf` | `self` | 申请人自己 |
| 直接上级 | `AssigneeSuperior` | `superior` | 直接上级 |
| 部门领导链 | `AssigneeDepartmentLeader` | `department_leader` | 多级主管链 |
| 表单字段 | `AssigneeFormField` | `form_field` | 由表单字段值决定 |

枚举类型是 `AssigneeKind`。动态加签位置使用 `AddAssigneeType`：
`AddAssigneeBefore`（`before`）、`AddAssigneeAfter`（`after`）、
`AddAssigneeParallel`（`parallel`）。

## 实例生命周期

```
提交 → 运行中 → 同意/拒绝 → 已同意/已拒绝
              → 撤回       → 已撤回
              → 回退       → 已退回
              → 终止       → 已终止
已撤回/已退回 → 重新提交   → 运行中（再次）
```

运行时状态机只声明以下实例状态流转：

| From | To |
| --- | --- |
| `running` | `approved` |
| `running` | `rejected` |
| `running` | `withdrawn` |
| `running` | `terminated` |
| `running` | `returned` |
| `returned` | `running` |
| `returned` | `terminated` |
| `returned` | `withdrawn` |
| `withdrawn` | `running` |
| `withdrawn` | `terminated` |

### 实例状态

| 状态 | 常量 | Wire value | 是否终态 |
| --- | --- | --- | --- |
| 运行中 | `InstanceRunning` | `running` | 否 |
| 已同意 | `InstanceApproved` | `approved` | 是 |
| 已拒绝 | `InstanceRejected` | `rejected` | 是 |
| 已撤回 | `InstanceWithdrawn` | `withdrawn` | 否 |
| 已退回 | `InstanceReturned` | `returned` | 否 |
| 已终止 | `InstanceTerminated` | `terminated` | 是 |

枚举类型是 `InstanceStatus`。

### 任务状态

| 状态 | 常量 | Wire value | 是否终态 |
| --- | --- | --- | --- |
| 等待中 | `TaskWaiting` | `waiting` | 否 |
| 待处理 | `TaskPending` | `pending` | 否 |
| 已同意 | `TaskApproved` | `approved` | 是 |
| 已拒绝 | `TaskRejected` | `rejected` | 是 |
| 已办理 | `TaskHandled` | `handled` | 是 |
| 已转交 | `TaskTransferred` | `transferred` | 是 |
| 已回退 | `TaskRolledBack` | `rolled_back` | 是 |
| 已取消 | `TaskCanceled` | `canceled` | 是 |
| 已移除 | `TaskRemoved` | `removed` | 是 |
| 已跳过 | `TaskSkipped` | `skipped` | 是 |

运行时状态机只声明以下任务状态流转：

| From | To |
| --- | --- |
| `waiting` | `pending` |
| `waiting` | `canceled` |
| `waiting` | `skipped` |
| `waiting` | `removed` |
| `pending` | `approved` |
| `pending` | `handled` |
| `pending` | `rejected` |
| `pending` | `transferred` |
| `pending` | `rolled_back` |
| `pending` | `canceled` |
| `pending` | `waiting` |
| `pending` | `removed` |

## 操作类型

| 操作 | 常量 | Wire value | 说明 |
| --- | --- | --- | --- |
| 提交 | `ActionSubmit` | `submit` | 发起新实例 |
| 同意 | `ActionApprove` | `approve` | 审批通过 |
| 办理 | `ActionHandle` | `handle` | 完成办理任务 |
| 拒绝 | `ActionReject` | `reject` | 审批拒绝 |
| 转交 | `ActionTransfer` | `transfer` | 转交给其他用户 |
| 撤回 | `ActionWithdraw` | `withdraw` | 申请人撤回 |
| 取消 | `ActionCancel` | `cancel` | 取消任务 |
| 回退 | `ActionRollback` | `rollback` | 退回到上一节点 |
| 加签 | `ActionAddAssignee` | `add_assignee` | 动态添加审批人 |
| 减签 | `ActionRemoveAssignee` | `remove_assignee` | 移除审批人 |
| 加抄送 | `ActionAddCC` | `add_cc` | 动态添加抄送人 |
| 执行 | `ActionExecute` | `execute` | 自动节点处理使用的内部执行动作 |
| 重新提交 | `ActionResubmit` | `resubmit` | 重新提交已退回或已撤回的实例 |
| 改派 | `ActionReassign` | `reassign` | 管理员改派任务 |
| 强制终止 | `ActionTerminate` | `terminate` | 管理员强制终止 |

## 回退配置

| 属性 | 可选值 |
| --- | --- |
| `RollbackType` | `RollbackNone`（`none`）、`RollbackPrevious`（`previous`）、`RollbackStart`（`start`）、`RollbackAny`（`any`）、`RollbackSpecified`（`specified`） |
| `RollbackDataStrategy` | `RollbackDataClear`（`clear`，重置表单）、`RollbackDataKeep`（`keep`，保留数据）|

同一申请人处理使用 `SameApplicantAction`，取值包括
`SameApplicantSelfApprove`（`self_approve`）、`SameApplicantAutoPass`
（`auto_pass`）、`SameApplicantTransferSuperior`（`transfer_superior`）。
连续审批人处理使用 `ConsecutiveApproverAction`，取值包括
`ConsecutiveApproverNone`（`none`）与 `ConsecutiveApproverAutoPass`
（`auto_pass`）。

## 空审批人处理

当节点找不到审批人时：

| 操作 | 常量 | Wire value |
| --- | --- | --- |
| 自动通过 | `EmptyAssigneeAutoPass` | `auto_pass` |
| 转交管理员 | `EmptyAssigneeTransferAdmin` | `transfer_admin` |
| 转交上级 | `EmptyAssigneeTransferSuperior` | `transfer_superior` |
| 转交申请人 | `EmptyAssigneeTransferApplicant` | `transfer_applicant` |
| 转交指定人 | `EmptyAssigneeTransferSpecified` | `transfer_specified` |

枚举类型是 `EmptyAssigneeAction`。

## 超时处理

| 操作 | 常量 | Wire value | 行为 |
| --- | --- | --- | --- |
| 无操作 | `TimeoutActionNone` | `none` | 仅标记超时 |
| 自动通过 | `TimeoutActionAutoPass` | `auto_pass` | 自动审批通过 |
| 自动拒绝 | `TimeoutActionAutoReject` | `auto_reject` | 自动审批拒绝 |
| 发送通知 | `TimeoutActionNotify` | `notify` | 仅发送通知 |
| 转交管理员 | `TimeoutActionTransferAdmin` | `transfer_admin` | 转交给节点管理员 |

枚举类型是 `TimeoutAction`。

## 表单数据存储

| 模式 | 常量 | Wire value | 存储位置 |
| --- | --- | --- | --- |
| JSON | `StorageJSON` | `json` | `apv_instance.form_data`（JSONB 列）|
| Table | `StorageTable` | `table` | 每个已发布版本一张生成的物理表，同时继续写入 `apv_instance.form_data` 作为 canonical JSON snapshot |

`StorageMode.IsValid()` 接受这两种导出模式。Table 模式通过 `FormTable` 和 `FormTableColumn` 记录生成的 DDL 元数据。

公开包暴露的流程设计和持久化模型包括 `FlowCategory`、`Flow`、`FlowVersion`、
`FlowNode`、`FlowEdge`、`FlowInitiator`、`FlowNodeAssignee`、`FlowNodeCC`、
`FormDefinition`、`FormFieldDefinition`、`FormSnapshot`、`ActionLog`、
`OperatorInfo` 和 `UrgeRecord`。流程版本状态使用 `VersionStatus`：
`VersionDraft`（`draft`）、`VersionPublished`（`published`）、
`VersionArchived`（`archived`）。

其他流程设计器枚举：

| 枚举 | Wire values |
| --- | --- |
| `InitiatorKind` | `user`、`role`、`department` |
| `ExecutionType` | `manual`、`auto_pass`、`auto_reject` |
| `ConditionKind` | `field`、`expression` |
| `CCKind` | `user`、`role`、`department`、`form_field` |
| `CCTiming` | `always`、`on_approve`、`on_reject` |
| `FieldKind` | `input`、`textarea`、`select`、`number`、`date`、`upload` |
| `ColumnDataType` | `string`、`text`、`integer`、`decimal`、`boolean`、`date`、`datetime`、`json` |
| `Permission` | `visible`、`editable`、`hidden`、`required` |

### 表存储元数据

当已发布版本使用 `StorageTable` 时，框架公开两个元数据模型：

| 模型 | 用途 | 关键 JSON 字段 |
| --- | --- | --- |
| `FormTable` | 每个流程版本对应一张生成的物理表 | `flowId`、`versionId`、`physicalTableName` |
| `FormTableColumn` | 每个表单字段或内置列对应一个生成列 | `formTableId`、`columnName`、`columnType`、`isNullable`、`sourceFieldKey`、`sortOrder` |

`ColumnDataType` 是表单定义使用的逻辑字段到列类型词汇，storage 层再把它映射成具体数据库方言的 SQL 类型。

### 设计器默认值

`NodeData.ApplyTo` 会把省略的设计器字段解析为导出的默认值，保证运行时和未触碰过的设计器控件一致：

| 常量 | 值 |
| --- | --- |
| `DefaultExecutionType` | `ExecutionManual` |
| `DefaultApprovalMethod` | `ApprovalParallel` |
| `DefaultPassRule` | `PassAll` |
| `DefaultEmptyAssigneeAction` | `EmptyAssigneeAutoPass` |
| `DefaultSameApplicantAction` | `SameApplicantSelfApprove` |
| `DefaultConsecutiveApproverAction` | `ConsecutiveApproverNone` |
| `DefaultRollbackType` | `RollbackPrevious` |
| `DefaultRollbackDataStrategy` | `RollbackDataKeep` |
| `DefaultTimeoutAction` | `TimeoutActionNone` |
| `DefaultCCTiming` | `CCTimingAlways` |
| `DefaultHandleApprovalMethod` | `ApprovalSequential` |
| `DefaultHandlePassRule` | `PassAny` |
| `DefaultUrgeCooldownMinutes` | `30` |
| `DefaultTenantID` | `"default"` |

### 实例进度投影

管理端和用户端 instance detail 响应会在实例快照旁公开两个只读投影：

| 投影 | 公开类型 | 用途 |
| --- | --- | --- |
| timeline | `TimelineEntryKind`、`TimelineEntry`、`NodeVisitStatus`、`NodeParticipant`、`Activity`、`ActivityUrge`、`CCRecipient` | 实例实际经过路径的时间线 |
| flow graph | `NodeProgressStatus`、`InstanceFlowGraph`、`FlowGraphNode`、`FlowGraphNodeData`、`FlowGraphEdge` | 带运行时进度标注、兼容 React Flow 的流程图 |

`FlowGraphNode.ID` 是 React Flow 设计期节点 id。`FlowGraphNode.NodeID` 是 action log 和 rollback target 使用的持久化 flow-node id。

进度和时间线枚举显式导出：

| 枚举 | 常量 |
| --- | --- |
| `TimelineEntryKind` | `TimelineEntryStart`、`TimelineEntryApproval`、`TimelineEntryHandle`、`TimelineEntryCC`、`TimelineEntryWithdraw`、`TimelineEntryTerminate` |
| `NodeVisitStatus` | `NodeVisitActive`、`NodeVisitPassed`、`NodeVisitRejected`、`NodeVisitReturned`、`NodeVisitCanceled` |
| `NodeProgressStatus` | `NodeProgressPending`、`NodeProgressActive`、`NodeProgressPassed`、`NodeProgressRejected`、`NodeProgressReturned`、`NodeProgressCanceled` |

这些投影背后的持久化 visit 模型是 `NodeVisit`。

## 事件发布

审批模块通过框架统一的事务性 outbox transport 发布领域事件（见 [事件总线](../features/event-bus)）。每条审批命令都把事件记录写在和业务变更相同的事务里，再由 outbox relay 转发到配置的 sink。

订阅方必须：

1. 用 `event.WithGroup("...")` 挂上 —— 路由解析到 at-least-once transport。
2. 依赖 Inbox 中间件去重（在 `event.middleware.inbox = true` 且 transport 声明 `AtLeastOnce` 时自动激活）。

> 老版本里的独立 `apv_event_outbox` 表和 `EventOutboxStatus` 常量已经废弃。审批模块不再自维护一个 outbox，而是组合使用框架统一 outbox。

### 领域事件类型

所有审批事件都实现 `DomainEvent`；`AllEventTypes()` 是导出的 topic 字符串完整注册表。实例和任务事件嵌入 `InstanceEventBase`（`instanceId`、`instanceNo`、`tenantId`、`title`、`flowId`、`flowCode`、可选 `businessRef`、`applicant`、`occurredTime`）。任务事件还嵌入 `TaskEventBase`（`taskId`、`nodeId`、`nodeName`）。流程事件嵌入 `FlowEventBase`（`flowId`、`tenantId`、`code`、`name`、`occurredTime`）。下表只列出这些 base 之外的字段。

实例生命周期：

| 类型常量 | Topic | Payload / 构造器 | 除通用字段外的 payload 字段 | 触发场景 |
| --- | --- | --- | --- | --- |
| `EventTypeInstanceCreated` | `approval.instance.created` | `InstanceCreatedEvent`, `NewInstanceCreatedEvent` | 无 | 一个新实例被启动 |
| `EventTypeInstanceCompleted` | `approval.instance.completed` | `InstanceCompletedEvent`, `NewInstanceCompletedEvent` | `finalStatus`、`finishedAt` | 实例进入终态 |
| `EventTypeInstanceWithdrawn` | `approval.instance.withdrawn` | `InstanceWithdrawnEvent`, `NewInstanceWithdrawnEvent` | `operator`（`UserInfo`） | 申请人撤回实例 |
| `EventTypeInstanceRolledBack` | `approval.instance.rolled_back` | `InstanceRolledBackEvent`, `NewInstanceRolledBackEvent` | `fromNodeId`、`fromNodeName`、`toNodeId`、`toNodeName`、`operator`（`UserInfo`） | 实例回退到之前的节点 |
| `EventTypeInstanceReturned` | `approval.instance.returned` | `InstanceReturnedEvent`, `NewInstanceReturnedEvent` | `fromNodeId`、`fromNodeName`、`toNodeId`、`toNodeName`、`operator`（`UserInfo`） | 实例退回申请人 |
| `EventTypeInstanceResubmitted` | `approval.instance.resubmitted` | `InstanceResubmittedEvent`, `NewInstanceResubmittedEvent` | `operator`（`UserInfo`） | 已退回或已撤回的实例重新提交 |
| `EventTypeInstanceBindingFailed` | `approval.instance.binding_failed` | `InstanceBindingFailedEvent`, `NewInstanceBindingFailedEvent` | `finalStatus`、`businessTable`、`errorMessage` | 引擎拥有的写回未能把终态持久化到业务行 |

节点生命周期：

| 类型常量 | Topic | Payload / 构造器 | 除通用字段外的 payload 字段 | 触发场景 |
| --- | --- | --- | --- | --- |
| `EventTypeNodeAutoPassed` | `approval.node.auto_passed` | `NodeAutoPassedEvent`, `NewNodeAutoPassedEvent` | `nodeId`、`nodeName`、`reason` | 节点因自动通过 execution、空审批人策略或同申请人策略自动通过 |

任务生命周期：

| 类型常量 | Topic | Payload / 构造器 | 除通用字段外的 payload 字段 | 触发场景 |
| --- | --- | --- | --- | --- |
| `EventTypeTaskCreated` | `approval.task.created` | `TaskCreatedEvent`, `NewTaskCreatedEvent` | `assignee`（`UserInfo`）、`deadline` | 任务创建；顺序审批中的后续任务可能先没有 `deadline`，表示还在等待前置任务 |
| `EventTypeTaskApproved` | `approval.task.approved` | `TaskApprovedEvent`, `NewTaskApprovedEvent` | `operator`（`UserInfo`）、`opinion` | 任务被批准 |
| `EventTypeTaskHandled` | `approval.task.handled` | `TaskHandledEvent`, `NewTaskHandledEvent` | `operator`（`UserInfo`）、`opinion` | handle 任务完成 |
| `EventTypeTaskRejected` | `approval.task.rejected` | `TaskRejectedEvent`, `NewTaskRejectedEvent` | `operator`（`UserInfo`）、`opinion` | 任务被驳回 |
| `EventTypeTaskCanceled` | `approval.task.canceled` | `TaskCanceledEvent`, `NewTaskCanceledEvent` | `assignee`（`UserInfo`）、`reason` | 引擎取消不再需要决策的任务 |
| `EventTypeTaskTransferred` | `approval.task.transferred` | `TaskTransferredEvent`, `NewTaskTransferredEvent` | `from`（`UserInfo`）、`to`（`UserInfo`）、`reason` | 任务被转交 |
| `EventTypeTaskReassigned` | `approval.task.reassigned` | `TaskReassignedEvent`, `NewTaskReassignedEvent` | `from`（`UserInfo`）、`to`（`UserInfo`）、`reason` | 管理员重新指派任务 |
| `EventTypeTaskTimedOut` | `approval.task.timed_out` | `TaskTimedOutEvent`, `NewTaskTimedOutEvent` | `assignee`（`UserInfo`）、`deadline` | 超时扫描器触发配置的超时动作 |
| `EventTypeAssigneesAdded` | `approval.task.assignees_added` | `AssigneesAddedEvent`, `NewAssigneesAddedEvent` | `addType`、`assignees`（`[]UserInfo`） | 动态新增审批人 |
| `EventTypeAssigneesRemoved` | `approval.task.assignees_removed` | `AssigneesRemovedEvent`, `NewAssigneesRemovedEvent` | `assignees`（`[]UserInfo`） | 动态移除审批人 |
| `EventTypeTaskDeadlineWarning` | `approval.task.deadline_warning` | `TaskDeadlineWarningEvent`, `NewTaskDeadlineWarningEvent` | `assignee`（`UserInfo`）、`deadline`、`hoursLeft` | 预警扫描器命中即将到期任务 |
| `EventTypeTaskUrged` | `approval.task.urged` | `TaskUrgedEvent`, `NewTaskUrgedEvent` | `urger`（`UserInfo`）、`target`（`UserInfo`）、`message` | 申请人催办 |

CC + 流程：

| 类型常量 | Topic | Payload / 构造器 | 除通用字段外的 payload 字段 | 触发场景 |
| --- | --- | --- | --- | --- |
| `EventTypeCCNotified` | `approval.cc.notified` | `CCNotifiedEvent`, `NewCCNotifiedEvent` | `nodeId`、`nodeName`、`recipients`（`[]UserInfo`）、`isManual` | CC 节点或手动 CC 动作完成通知 |
| `EventTypeFlowCreated` | `approval.flow.created` | `FlowCreatedEvent`, `NewFlowCreatedEvent` | `categoryId` | 流程创建 |
| `EventTypeFlowUpdated` | `approval.flow.updated` | `FlowUpdatedEvent`, `NewFlowUpdatedEvent` | 无 | 流程更新 |
| `EventTypeFlowDeployed` | `approval.flow.deployed` | `FlowDeployedEvent`, `NewFlowDeployedEvent` | `versionId`、`version` | 流程版本部署 |
| `EventTypeFlowToggled` | `approval.flow.toggled` | `FlowToggledEvent`, `NewFlowToggledEvent` | `isActive` | 流程启停状态变化 |
| `EventTypeFlowPublished` | `approval.flow.published` | `FlowPublishedEvent`, `NewFlowPublishedEvent` | `versionId` | 流程版本发布 |

## CallerContext 与多租户安全

资源/命令 handler 在每次请求时解析出一个 `CallerContext`，承载本次调用的租户权能。`TenantID`、`IsSuperAdmin`、`IsSystemInternal` 必须**且只能**有一个为真；零值 `CallerContext` 走 fail-closed 路径（被视为未授权）。

| 字段 | 含义 |
| --- | --- |
| `TenantID` | 调用者所在租户，所有读写都按它过滤 |
| `IsSuperAdmin` | 调用者持有 `approval:super_admin` 角色，允许跨租户访问 |
| `IsSystemInternal` | 框架内代码发起的调用（扫描器、listener） |

跨租户尝试返回 `approval.ErrCrossTenantAccess`。`IsSuperAdmin(p)` 判定 principal 是否带该 override 角色。override 角色字符串本身以 `approval.SuperAdminRole`（`"approval:super_admin"`）形式导出，方便宿主进行权限装配。

## 生命周期扩展点

| 扩展点 | 时机 | 失败语义 |
| --- | --- | --- |
| `approval.InstanceLifecycleHook`（FX group `vef:approval:lifecycle_hooks`） | 同步、业务事务内部 | 返回 error 会回滚整条命令 |
| `approval.BusinessRefProvider` | 同步、start-instance 事务内部 | 返回 error 会回滚实例创建 |
| `approval.BusinessRefResolver` | `InstanceCompletedEvent` 之后的异步写回路径 | 返回 error 会让引擎拥有的写回失败，发布 `InstanceBindingFailedEvent`，并由 outbox 路径重试 |
| 事件订阅（`event.SubscribeTyped`） | 异步、事务提交后 | bus 通过 outbox relay 重试，消费者必须幂等 |

### `InstanceLifecycleHook`

```go
type InstanceLifecycleHook interface {
    OnInstanceCreated(ctx context.Context, db orm.DB, instance *Instance) error
    OnInstanceCompleted(ctx context.Context, db orm.DB, instance *Instance, finalStatus InstanceStatus) error
}
```

事务内必须成立的不变量（比如分配一个紧耦合的业务行）应该用 lifecycle hook；其他场景都用事件订阅。

### `BusinessRefProvider` 和 `BusinessRefResolver`

```go
type BusinessRefProvider interface {
    OnInstanceCreated(ctx context.Context, db orm.DB, flow *Flow, instance *Instance) (businessRef string, err error)
}

type BusinessRefResolver interface {
    ResolveRecordID(ctx context.Context, flow *Flow, businessRef string) (string, error)
}
```

`BusinessRefProvider` 在 `Flow.BindingMode == BindingBusiness` 时提供或分配 opaque 的 `Instance.BusinessRef`；通过 `vef.SupplyBusinessRefProvider` 注入。引擎不会直接解析 `BusinessRef`。最终状态写回时，它会询问 `BusinessRefResolver`（默认是 identity resolver）提取要和 `Flow.BusinessPkField` 匹配的 record id；复合 ref 形状通过 `vef.SupplyBusinessRefResolver` 注入自定义解析。

审批引擎拥有状态写回（`UPDATE businessTable SET businessStatusField = ? WHERE businessPkField = ?`）。宿主可以用 lifecycle hook 或事件订阅围绕它扩展，但写回本身由配置驱动，不再由扩展 hook 替换。

## 业务标识符校验

当 `Flow.BindingMode == BindingBusiness` 时，流程会携带 SQL 标识符（`BusinessTable`、`BusinessPkField`、`BusinessStatusField`），引擎拥有的写回会把它们直接拼到 `UPDATE` 模板里。为防止 SQL 注入，框架按 `^[A-Za-z_][A-Za-z0-9_]{0,62}$` 白名单校验：

```go
if err := approval.ValidateBusinessIdentifier(table); err != nil {
    return err
}
```

空字符串/全空白通过校验——是否把"未配置"算错误由调用方决定。超出白名单的值会返回 `approval.ErrInvalidBusinessIdentifier`。管理端 Flow CRUD 应该把它向上抛，让操作员看到可读的错误。

## 错误面

可 import 的 `approval` 包导出四个普通 Go sentinel。它们可用
`errors.Is` 识别，但它们不是 `result.Error`，本身不带 API code 或 HTTP
status。

| 错误 | 源包 | 含义 |
| --- | --- | --- |
| `approval.ErrCrossTenantAccess` | `approval` | 非 super-admin 调用者尝试跨租户访问 |
| `approval.ErrInvalidBusinessIdentifier` | `approval` | business table / field 标识符未通过 SQL 标识符白名单 |
| `approval.ErrUnknownNodeKind` | `approval` | `NodeDefinition.ParseData` 遇到不支持的 `kind` |
| `approval.ErrNodeDataUnmarshal` | `approval` | `NodeDefinition.ParseData` 无法解码节点 `data` |

内置审批资源通过标准 API envelope 返回模块自己的 `result.Error`。这些值位于
internal 包中，所以宿主应用应把下表的 code/message 组合视为公开 wire surface，
而不是 import internal Go symbol。

| Code | Code constant | Error value | i18n message key | 说明 |
| --- | --- | --- | --- | --- |
| `40001` | `ErrCodeFlowNotFound` | `ErrFlowNotFound` | `approval_flow_not_found` | flow 查找失败 |
| `40002` | `ErrCodeFlowNotActive` | `ErrFlowNotActive` | `approval_flow_not_active` | flow 已禁用 |
| `40003` | `ErrCodeNoPublishedVersion` | `ErrNoPublishedVersion` | `approval_no_published_version` | flow 没有已发布版本 |
| `40004` | `ErrCodeVersionNotDraft` | `ErrVersionNotDraft` | `approval_version_not_draft` | 当前操作要求 draft 版本 |
| `40005` | `ErrCodeInvalidFlowDesign` | `ErrInvalidFlowDesign` | `approval_invalid_flow_design` | 图或节点设计未通过校验 |
| `40006` | `ErrCodeFlowCodeExists` | `ErrFlowCodeExists` | `approval_flow_code_exists` | flow code 重复 |
| `40007` | `ErrCodeVersionNotFound` | `ErrVersionNotFound` | `approval_version_not_found` | flow version 查找失败 |
| `40008` | `ErrCodeInvalidBusinessIdentifier` | `ErrInvalidBusinessIdentifier` | `approval_invalid_business_identifier` | business table / field 标识符未通过校验 |
| `40009` | `ErrCodeInvalidTitleTemplate` | `ErrInvalidTitleTemplate` | `approval_invalid_title_template` | instance title template 无法解析 |
| `40010` | `ErrCodeInvalidFormDesign` | `ErrInvalidFormDesign` | `approval_invalid_form_design` | 表单 schema 未通过设计期校验 |
| `40011` | `ErrCodeBindingIncomplete` | `ErrBindingIncomplete` | `approval_binding_incomplete` | business binding 缺少必需 table / primary-key / status 字段 |
| `40012` | `ErrCodeInvalidStorageMode` | `ErrInvalidStorageMode` | `approval_invalid_storage_mode` | deploy 请求的 storage mode 不是 `json` 或 `table` |
| `40013` | `ErrCodeFlowBindingLocked` | `ErrFlowBindingLocked` | `approval_flow_binding_locked` | flow business-binding 设置在仍有 running instance 时被锁定 |
| `40101` | `ErrCodeInstanceNotFound` | `ErrInstanceNotFound` | `approval_instance_not_found` | instance 查找失败 |
| `40102` | `ErrCodeInstanceCompleted` | `ErrInstanceCompleted` | `approval_instance_completed` | instance 已完成 |
| `40103` | `ErrCodeNotAllowedInitiate` | `ErrNotAllowedInitiate` | `approval_not_allowed_initiate` | 调用者不能发起该 flow |
| `40104` | `ErrCodeWithdrawNotAllowed` | `ErrWithdrawNotAllowed` | `approval_withdraw_not_allowed` | 当前状态不允许撤回 |
| `40105` | `ErrCodeResubmitNotAllowed` | `ErrResubmitNotAllowed` | `approval_resubmit_not_allowed` | 当前状态不允许重新提交 |
| `40106` | `ErrCodeInvalidInstanceTransition` | `ErrInvalidInstanceTransition` | `approval_invalid_instance_transition` | instance 状态流转无效 |
| `40201` | `ErrCodeTaskNotFound` | `ErrTaskNotFound` | `approval_task_not_found` | task 查找失败 |
| `40202` | `ErrCodeTaskNotPending` | `ErrTaskNotPending` | `approval_task_not_pending` | task 不是 pending |
| `40203` | `ErrCodeNotAssignee` | `ErrNotAssignee` | `approval_not_assignee` | 调用者不是该 task 的 assignee |
| `40204` | `ErrCodeInvalidTaskTransition` | `ErrInvalidTaskTransition` | `approval_invalid_task_transition` | task 状态流转无效 |
| `40205` | `ErrCodeRollbackNotAllowed` | `ErrRollbackNotAllowed` | `approval_rollback_not_allowed` | 此处不允许回退 |
| `40206` | `ErrCodeAddAssigneeNotAllowed` | `ErrAddAssigneeNotAllowed` | `approval_add_assignee_not_allowed` | 禁止动态新增 assignee |
| `40207` | `ErrCodeTransferNotAllowed` | `ErrTransferNotAllowed` | `approval_transfer_not_allowed` | 禁止转交 |
| `40208` | `ErrCodeOpinionRequired` | `ErrOpinionRequired` | `approval_opinion_required` | 必填 opinion 为空 |
| `40209` | `ErrCodeManualCcNotAllowed` | `ErrManualCcNotAllowed` | `approval_manual_cc_not_allowed` | 禁止手动 CC |
| `40210` | `ErrCodeRemoveAssigneeNotAllowed` | `ErrRemoveAssigneeNotAllowed` | `approval_remove_assignee_not_allowed` | 禁止动态移除 assignee |
| `40211` | `ErrCodeInvalidAddAssigneeType` | `ErrInvalidAddAssigneeType` | `approval_invalid_add_assignee_type` | `addType` 不是 `before`、`after` 或 `parallel` |
| `40212` | `ErrCodeNotApplicant` | `ErrNotApplicant` | `approval_not_applicant` | 调用者不是申请人 |
| `40213` | `ErrCodeInvalidRollbackTarget` | `ErrInvalidRollbackTarget` | `approval_invalid_rollback_target` | 回退目标不允许 |
| `40214` | `ErrCodeLastAssigneeRemoval` | `ErrLastAssigneeRemoval` | `approval_last_assignee_removal` | 移除后会没有 active assignee |
| `40215` | `ErrCodeInvalidTransferTarget` | `ErrInvalidTransferTarget` | `approval_invalid_transfer_target` | 转交或重新指派目标无效 |
| `40216` | `ErrCodeNoUsersSpecified` | `ErrNoUsersSpecified` | `approval_no_users_specified` | 用户列表操作没有提供目标用户 |
| `40301` | `ErrCodeNoAssignee` | `ErrNoAssignee` | `approval_no_assignee` | 无法解析出 assignee |
| `40302` | `ErrCodeAssigneeResolveFailed` | `ErrAssigneeResolveFailed` | `approval_assignee_resolve_failed` | assignee resolver 失败 |
| `40401` | `ErrCodeFormValidationFailed` | `ErrFormValidationFailed` | `approval_form_validation_failed` | 通用表单校验失败 |
| `40401` | `ErrCodeFormValidationFailed` | `ErrFormDataTooLarge` | `approval_form_data_too_large` | 同一 code；JSON 编码后的 `formData` 超过 64 KiB |
| `40401` | `ErrCodeFormValidationFailed` | 动态表单校验 `result.Err` | `approval_form_field_not_defined`, `approval_form_field_required`, `approval_form_field_must_be_string`, `approval_form_field_must_be_number`, `approval_form_field_must_be_integer`, `approval_form_field_min_length`, `approval_form_field_max_length`, `approval_form_field_invalid_validation`, `approval_form_field_pattern_mismatch`, `approval_form_field_min_value`, `approval_form_field_max_value`, `approval_form_field_empty`, `approval_form_field_invalid_file_item`, `approval_form_field_must_be_file`, `approval_form_field_invalid_value` | 字段级校验消息在运行时构造 |
| `40402` | `ErrCodeFieldNotEditable` | `ErrFieldNotEditable` | `approval_field_not_editable` | 提交字段对当前 task 不可编辑 |
| `40501` | `ErrCodeDelegationNotFound` | `ErrDelegationNotFound` | `approval_delegation_not_found` | delegation 查找失败 |
| `40502` | `ErrCodeDelegationConflict` | `ErrDelegationConflict` | `approval_delegation_conflict` | delegation 时间窗与已有委托冲突 |
| `40601` | `ErrCodeUrgeCooldown` | 动态 urge `result.Err` | `approval_urge_too_frequent` | 没有静态 sentinel；消息会带 `minutes`；非正数 `urgeCooldownMinutes` 默认按 30 分钟处理 |
| `40701` | `ErrCodeAccessDenied` | `ErrAccessDenied` | `approval_access_denied` | 调用者缺少审批域访问权限 |
| `40702` | `ErrCodeTerminateNotAllowed` | `ErrTerminateNotAllowed` | `approval_terminate_not_allowed` | 当前 instance 状态不允许终止 |

`ErrEventRouteNotTransactional`、`ErrEventRouteNotSubscribable`、
`ErrTenantNotResolved` 这类启动和租户解析诊断位于 `internal/approval/...`；
它们不是可 import 的公开 Go API，但 event routing 或 principal tenant
details 配错时，操作员可能会看到它们被包装后的错误消息。

## 其他公开 API 索引

| 范围 | 公开 API |
| --- | --- |
| caller safety | `CallerContext`, `SystemCaller`, `IsSuperAdmin`, `SuperAdminRole`, `ErrCrossTenantAccess` |
| form data | `FormData`, `NewFormData`, `FormDefinition`, `FormFieldDefinition`, `FormSnapshot`, `ValidationRule`, `StorageMode`, `StorageJSON`, `StorageTable`, `FieldKind`, `FieldInput`, `FieldNumber`, `FieldDate`, `FieldTextarea`, `FieldSelect`, `FieldUpload`, `FieldOption`, `ColumnDataType`, `ColumnString`, `ColumnText`, `ColumnInteger`, `ColumnDecimal`, `ColumnBoolean`, `ColumnDate`, `ColumnDatetime`, `ColumnJSON` |
| table storage | `FormTable`, `FormTableColumn` |
| flow models | `FlowCategory`, `Flow`, `FlowVersion`, `FlowNode`, `FlowEdge`, `FlowInitiator`, `FlowNodeAssignee`, `FlowNodeCC`, `VersionStatus`, `VersionDraft`, `VersionPublished`, `VersionArchived`, `ActionLog`, `OperatorInfo`, `UrgeRecord`, `DefaultTenantID` |
| node design | `FlowDefinition`, `NodeDefinition`, `EdgeDefinition`, `Position`, `NodeData`, `BaseNodeData`, `StartNodeData`, `ApprovalNodeData`, `HandleNodeData`, `ConditionNodeData`, `CCNodeData`, `EndNodeData`, `ErrUnknownNodeKind`, `ErrNodeDataUnmarshal` |
| conditions | `ConditionKind`, `ConditionField`, `ConditionExpression`, `Condition`, `ConditionGroup`, `ConditionBranch`, `EvaluationContext`, `ConditionEvaluator`, `InstanceGlobalsResolver` |
| initiators and assignees | `InitiatorKind`, `InitiatorUser`, `InitiatorRole`, `InitiatorDepartment`, `AssigneeKind`, `AssigneeDefinition`, `AssigneeService`, `ResolvedAssignee`, `UserInfo`, `UserInfoResolver`, `RoleMembershipChecker`, `AddAssigneeType`, `AddAssigneeBefore`, `AddAssigneeAfter`, `AddAssigneeParallel` |
| CC | `CCKind`, `CCUser`, `CCRole`, `CCDepartment`, `CCFormField`, `CCTiming`, `CCTimingAlways`, `CCTimingOnApprove`, `CCTimingOnReject`, `CCDefinition`, `CCRecord` |
| node behavior | `ApprovalMethod`, `TaskNodeData`, `ExecutionType`, `ExecutionManual`, `ExecutionAutoPass`, `ExecutionAutoReject`, `ConsecutiveApproverAction`, `ConsecutiveApproverNone`, `ConsecutiveApproverAutoPass`, `SameApplicantAction`, `SameApplicantSelfApprove`, `SameApplicantAutoPass`, `SameApplicantTransferSuperior`, `Permission`, `PermissionVisible`, `PermissionEditable`, `PermissionRequired`, `PermissionHidden`, `DefaultExecutionType`, `DefaultApprovalMethod`, `DefaultPassRule`, `DefaultEmptyAssigneeAction`, `DefaultSameApplicantAction`, `DefaultConsecutiveApproverAction`, `DefaultRollbackType`, `DefaultRollbackDataStrategy`, `DefaultTimeoutAction`, `DefaultCCTiming`, `DefaultHandleApprovalMethod`, `DefaultHandlePassRule`, `DefaultUrgeCooldownMinutes` |
| rollback and timeouts | `RollbackType`, `RollbackNone`, `RollbackPrevious`, `RollbackStart`, `RollbackAny`, `RollbackSpecified`, `RollbackDataStrategy`, `RollbackDataClear`, `RollbackDataKeep`, `EmptyAssigneeAction`, `EmptyAssigneeAutoPass`, `EmptyAssigneeTransferAdmin`, `EmptyAssigneeTransferSuperior`, `EmptyAssigneeTransferApplicant`, `EmptyAssigneeTransferSpecified`, `TimeoutAction`, `TimeoutActionNone`, `TimeoutActionAutoPass`, `TimeoutActionAutoReject`, `TimeoutActionNotify`, `TimeoutActionTransferAdmin` |
| action and status enums | `ActionType`, `InstanceStatus`, `TaskStatus`, `NodeKind`, `StorageMode`, `VersionStatus` |
| pass rules | `PassRule`, `PassRuleContext`, `PassRuleStrategy`, `PassRuleResult`, `PassRulePending`, `PassRulePassed`, `PassRuleRejected` |
| progress views | `TimelineEntryKind`, `TimelineEntry`, `NodeVisitStatus`, `NodeProgressStatus`, `InstanceFlowGraph`, `FlowGraphNode`, `FlowGraphNodeData`, `FlowGraphEdge`, `NodeParticipant`, `Activity`, `ActivityUrge`, `CCRecipient` |
| events | 所有 `New...Event` 构造器、`DomainEvent`、`InstanceEventBase`、`TaskEventBase`、`FlowEventBase`、`NewInstanceEventBase`、`NewTaskEventBase`、`NewFlowEventBase`、`PayloadOccurredAt`、`AllEventTypes` 和 `EventType...` 常量 |
| extension interfaces | `InstanceLifecycleHook`, `BusinessRefProvider`, `BusinessRefResolver`, `InstanceNoGenerator`, `ConditionEvaluator`, `InstanceGlobalsResolver`, `PrincipalTenantResolver`, `PrincipalDepartmentResolver`, `RoleMembershipChecker` |
| admin DTOs | `approval/admin` 包：`Instance`, `InstanceDetail`, `InstanceDetailInfo`, `Task`, `ActionLog`, `Metrics` |
| user DTOs | `approval/my` 包：`PendingTask`, `CompletedTask`, `CCRecord`, `InitiatedInstance`, `AvailableFlow`, `InstanceDetail`, `InstanceInfo`, `PendingCounts` |

## 委托

用户可以将审批权限委托给他人：

```go
type Delegation struct {
    DelegatorID    string         // 委托人
    DelegateeID    string         // 被委托人
    FlowCategoryID *string        // 可选：限定分类
    FlowID         *string        // 可选：限定特定流程
    StartTime      timex.DateTime // 委托开始时间
    EndTime        timex.DateTime // 委托结束时间
    IsActive       bool
}
```

## 流程定义（兼容 React Flow）

`FlowDefinition` 结构体兼容 React Flow 的 JSON 格式：

```go
type FlowDefinition struct {
    Nodes []NodeDefinition `json:"nodes"`
    Edges []EdgeDefinition `json:"edges"`
}
```

每个 `NodeDefinition` 包含 `Kind` 和强类型的 `Data` 字段，框架会按 `Kind` 把 `Data` 解析成对应结构（`StartNodeData`、`ApprovalNodeData`、`HandleNodeData`、`ConditionNodeData`、`CCNodeData`、`EndNodeData`）。

### 流程 JSON Wire Shape

`deploy` 会把流程定义当作完整快照处理。`NodeDefinition.ParseData` 根据
`kind` 选择对应的强类型 `data` 结构；未知 kind 返回 `ErrUnknownNodeKind`，
节点 `data` 的 JSON 解析失败会用 `ErrNodeDataUnmarshal` 包装。

| 类型 | JSON 字段 |
| --- | --- |
| `FlowDefinition` | `nodes`、`edges` |
| `NodeDefinition` | `id`、`kind`、`position`、`data`；`position` 包含 `x` 和 `y` |
| `EdgeDefinition` | `id`、`source`、`target`、`sourceHandle`、`data` |

只有条件节点的出边需要 `sourceHandle`，且它必须匹配某个分支 `id`。非条件节点的出边必须省略 `sourceHandle`。`EdgeDefinition.data` 是设计器元数据，保存在版本 `flowSchema` 中；运行时流转使用 `source`、`target` 和 `sourceHandle`。

节点 `data` 字段如下：

| 节点 data 类型 | JSON 字段 |
| --- | --- |
| `BaseNodeData` | `name`、`description`；每种节点 data 都嵌入它 |
| `StartNodeData` | 只有 base 字段 |
| `EndNodeData` | 只有 base 字段 |
| `TaskNodeData` | `assignees`、`executionType`、`emptyAssigneeAction`、`fallbackUserIds`、`adminUserIds`、`isTransferAllowed`、`isOpinionRequired`、`timeoutHours`、`timeoutAction`、`timeoutNotifyBeforeHours`、`urgeCooldownMinutes`、`ccs`、`fieldPermissions` |
| `ApprovalNodeData` | base 字段 + `TaskNodeData` 字段 + `approvalMethod`、`passRule`、`passRatio`、`sameApplicantAction`、`consecutiveApproverAction`、`rollbackType`、`rollbackDataStrategy`、`rollbackTargetKeys`、`isRollbackAllowed`、`isAddAssigneeAllowed`、`addAssigneeTypes`、`isRemoveAssigneeAllowed`、`isManualCcAllowed` |
| `HandleNodeData` | base 字段 + `TaskNodeData` 字段；未设置时部署会默认 `approvalMethod = sequential`、`passRule = any` |
| `CCNodeData` | base 字段 + `ccs`、`isReadConfirmRequired`、`fieldPermissions` |
| `ConditionNodeData` | base 字段 + `branches` |

`assignees` 条目使用 `kind`、`ids`、`formField` 和 `sortOrder`。`ccs`
条目使用 `kind`、`ids`、`formField` 和 `timing`。部署时这些嵌入数组会额外物化为
`FlowNodeAssignee` 和 `FlowNodeCC` 记录，不只是写入 `FlowNode` 行。

条件分支使用 `id`、`label`、`conditionGroups`、`isDefault` 和 `priority`。
每个 `conditionGroups` 条目包含 `conditions`；每个 condition 使用 `kind`、
`subject`、`operator`、`value` 和 `expression`。

`timeoutHours` 和 `timeoutNotifyBeforeHours` 的单位是小时。
`urgeCooldownMinutes` 的单位是分钟；小于等于 0 时使用 30 分钟运行时默认值。`rollbackTargetKeys` 只在
`rollbackType = specified` 时校验，里面放的是节点 key，不是数据库节点 ID。
任务处理时，提交的 `formData` 只会合并 `fieldPermissions` 中标为 `editable`
或 `required` 的字段；标为 `visible`、`hidden`，或没有出现在 map 中的字段都会在本次任务更新中被忽略。

### 表单 JSON Wire Shape

`FormDefinition` 使用 `fields` 数组。每个 `FormFieldDefinition` 条目使用
`key`、`kind`、`label`、`placeholder`、`defaultValue`、`isRequired`、
`options`、`validation`、`props`、`sortOrder`、`columnType` 和 `scale`。每个 option 使用 `label`
和 `value`。

`validation` 支持 `minLength`、`maxLength`、`min`、`max`、`pattern` 和
`message`。提交的 `formData` 在 JSON 编码后最大 64 KiB，即使流程没有表单
schema 也会执行这个大小限制。有 schema 时，额外的表单 key 会被拒绝；必填字段会拒绝
缺失、`null`、空白字符串和空数组。`input`、`textarea`、`date` 字段必须是字符串，
可使用 `minLength`、`maxLength` 和 `pattern`。`number` 字段接受 JSON 数字，
可使用 `min` 和 `max`。`select` 字段在配置了 `options` 时会校验标量或数组值是否
存在于选项中。`upload` 字段接受非空白字符串、非空 `[]string`，或非空且每项都是非空白
字符串的数组。`validation.message` 只作为 `pattern` 不匹配时的自定义错误信息；
其他校验失败使用模块 i18n 消息。

## 实例编号生成

实现 `InstanceNoGenerator` 接口来自定义实例编号：

```go
type InstanceNoGenerator interface {
    Generate(ctx context.Context, flowCode string) (string, error)
}
```
