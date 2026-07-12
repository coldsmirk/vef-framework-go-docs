---
sidebar_position: 2
---

# RPC 资源

启用模块后，会注册以下 RPC 资源。这些操作都不是公开接口：调用者必须已认证；
表中列出 `RequiredPermission` 时还会校验对应权限点。

RPC 调用使用 [API](../building-apis/api.md) 里的标准 envelope：`resource`、`action`、
`version`、`params` 和 `meta`。下面标成 `params.*` 的字段从 `params` 解码；
标成 `meta.*` 的字段从 `meta` 解码。生成的
[运行时 API 索引](../reference/runtime-api-index.md) 包含每个请求/响应 DTO 的
完整 JSON 字段 ledger。

## `approval/category`

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

## `approval/delegation`

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

## `approval/flow`

| Action | 权限 | 参数 | 说明 |
| --- | --- | --- | --- |
| `create` | `approval.flow.create` | `CreateFlowParams` | 开启审计 |
| `deploy` | `approval.flow.deploy` | `DeployFlowParams` | 开启审计 |
| `publish_version` | `approval.flow.publish` | `PublishVersionParams` | 开启审计 |
| `update` | `approval.flow.update` | `UpdateParams` | 开启审计 |
| `toggle_active` | `approval.flow.update` | `ToggleActiveParams` | 开启审计 |
| `get_graph` | `approval.flow.query` | `GetGraphParams` | 读取已发布流程图 |
| `find_flows` | `approval.flow.query` | `FindFlowsParams` | 分页查询流程 |
| `find_initiators` | `approval.flow.query` | `FindInitiatorsParams` | 查询发起人配置 |
| `find_versions` | `approval.flow.query` | `FindVersionsParams` | 查询单个流程的版本列表 |

| Action | 请求字段 |
| --- | --- |
| `create` | `params.tenantId` 必填、`params.code` 必填、`params.name` 必填、`params.categoryId` 必填、`params.bindingMode` 必填、`params.icon`、`params.description`、`params.businessBinding`、`params.adminUserIds`、`params.isAllInitiationAllowed`、`params.instanceTitleTemplate`、`params.initiators` |
| `deploy` | `params.flowId` 必填、`params.description`、`params.flowDefinition` 必填、`params.formSchema`、`params.storageMode` |
| `publish_version` | `params.versionId` 必填 |
| `update` | `params.flowId` 必填、`params.name` 必填、`params.bindingMode` 必填、`params.instanceTitleTemplate` 必填、`params.icon`、`params.description`、`params.businessBinding`、`params.adminUserIds`、`params.isAllInitiationAllowed`、`params.initiators` |
| `toggle_active` | `params.flowId` 必填、`params.isActive` |
| `get_graph` | `params.flowId` 必填、`params.tenantId` |
| `find_flows` | `params.tenantId`、`params.categoryId`、`params.keyword`、`params.isActive`、`params.page`、`params.pageSize` |
| `find_initiators` | `params.flowId` 必填、`params.tenantId` |
| `find_versions` | `params.flowId` 必填、`params.tenantId` |

`params.initiators` 条目使用 `kind`（`user`、`role` 或 `department`）和
`ids`。`params.formSchema` 是可选的宿主自有表单设计器文档，框架原样透传（见
[表单 Schema 与派生字段](./flow-design.md#表单-schema-与派生字段)）；没有表单的
流程直接省略它。业务绑定场景下，`params.businessBinding` 是一个
`approval.BusinessBindingConfig` 对象（`tableName`、`keyColumns`、
`statusColumn`、`instanceIdColumn` 必填，`startedAtColumn` /
`finishedAtColumn` / `statusMapping` 可选）；所有标识符都会经过
`ValidateBusinessIdentifier` 校验，记录键必须与真实表上的非空主键或唯一键
完全一致，且实例运行期间绑定设置保持冻结（`ErrFlowBindingLocked`）。

## `approval/instance`

| Action | 权限 | 参数 | 说明 |
| --- | --- | --- | --- |
| `start` | `approval.instance.start` | `StartParams` | 开启审计 |
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

## `approval/my`

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
`flowGraph`、`availableActions` 和 `fieldPermissions`。已提交表单数据在
`instance.formData` 中；业务绑定流程的 opaque 业务引用在
`instance.businessRef` 中；`formSchema` 是实例提交时所对应版本固定下来的
宿主表单设计器文档，原样返回——框架以语义等价的 JSON 存储且从不解释它。
`fieldPermissions`（v0.38）是按查看者投影的字段交互性映射，为每个顶层表单
字段都物化一个值（`visible` / `editable` / `hidden` / `required`）；客户端
按原样应用，`instance.formData` 已剥离查看者无权看到的字段（见
[节点字段权限](./flow-design.md#节点字段权限)）。

`availableActions` 是查询层给 UI 的提示。对申请人来说，实例可以流转到
`withdrawn` 时包含 `withdraw`，实例已退回或已撤回时包含
`resubmit`。对 pending task 来说，handle 节点包含 `handle`，其他节点
包含 `approve`，随后包含 `reject`，并按当前节点开关追加 `transfer`、
`rollback`、`add_assignee` 或 `add_cc`。只要实例存在任何 pending task，
还会包含 `urge`。命令 handler 仍会独立做最终校验。

## `approval/admin`

| Action | 权限 | 参数 | 说明 |
| --- | --- | --- | --- |
| `find_instances` | `approval.instance.query` | `AdminFindInstancesParams` | 非 super-admin 按租户过滤 |
| `find_tasks` | `approval.task.query` | `AdminFindTasksParams` | 非 super-admin 按租户过滤 |
| `get_instance_detail` | `approval.instance.detail` | `AdminGetInstanceDetailParams` | 完整管理端详情 |
| `find_action_logs` | `approval.action_log.query` | `AdminFindActionLogsParams` | 要求 `instanceId` |
| `get_metrics` | `approval.metrics.query` | `AdminGetMetricsParams` | 聚合指标 |
| `find_business_projections` | `approval.binding.query` | `AdminFindBusinessProjectionsParams` | 持久化绑定收敛状态（v0.38） |
| `terminate_instance` | `approval.instance.terminate` | `AdminTerminateInstanceParams` | 开启审计 |
| `reassign_task` | `approval.task.reassign` | `AdminReassignTaskParams` | 开启审计 |
| `retry_business_projection` | `approval.binding.retry` | `AdminRetryBusinessProjectionParams` | 开启审计；立即重试一条 eventual 投影（v0.38） |

| Action | 请求字段 |
| --- | --- |
| `find_instances` | `params.tenantId`、`params.applicantId`、`params.status`、`params.flowId`、`params.keyword`、`params.page`、`params.pageSize` |
| `find_tasks` | `params.tenantId`、`params.assigneeId`、`params.instanceId`、`params.status`、`params.page`、`params.pageSize` |
| `get_instance_detail` | `params.instanceId` 必填 |
| `find_action_logs` | `params.instanceId` 必填、`params.tenantId`、`params.page`、`params.pageSize` |
| `get_metrics` | `params.tenantId` |
| `find_business_projections` | `params.tenantId`、`params.status`（`pending`、`processing`、`applied`、`failed`）、`params.page`、`params.pageSize` |
| `terminate_instance` | `params.instanceId` 必填、`params.reason` 最多 2000 字符 |
| `reassign_task` | `params.taskId` 必填、`params.newAssigneeId` 必填、`params.reason` 最多 2000 字符 |
| `retry_business_projection` | `params.projectionId` 必填 |

管理端列表和指标查询中，非 super-admin 调用者会忽略提交的 `tenantId` override，
并被过滤到自己的租户。super-admin 可以传 `tenantId` 只看一个租户，也可以省略它
获得跨租户可见性。

管理端列表/详情 DTO 在线上的租户字段名保持为 `tenantId`。原始审计日志仍可通过
`find_action_logs` 单独分页查询；详情响应携带 timeline 和 flow graph 投影供 UI 渲染。

## 响应 DTO 字段

管理端响应使用 `approval/admin` 包中的 DTO：

| DTO | JSON 字段 |
| --- | --- |
| `admin.Instance` | `instanceId`、`instanceNo`、`title`、`tenantId`、`flowId`、`flowName`、`applicant`（`UserInfo`）、`status`、`currentNodeName`、`createdAt`、`finishedAt` |
| `admin.Task` | `taskId`、`instanceId`、`instanceTitle`、`flowName`、`nodeName`、`assignee`（`UserInfo`）、`status`、`createdAt`、`deadline`、`finishedAt` |
| `admin.InstanceDetail` | `instance`、`formSchema`（宿主设计器文档，原样返回）、`timeline`、`flowGraph` |
| `admin.InstanceDetailInfo` | `instanceId`、`instanceNo`、`title`、`tenantId`、`flowId`、`flowName`、`flowVersionId`、`applicant`、`status`、`currentNodeId`、`currentNodeName`、`businessRef`、`formData`、`createdAt`、`finishedAt` |
| `admin.ActionLog` | `logId`、`action`、`nodeId`、`taskId`、`operator`、`transferTo`、`rollbackToNodeId`、`addedAssignees`、`removedAssignees`、`ccUsers`、`opinion`、`attachments`、`createdAt` |
| `admin.Metrics` | `tenantId`、`capturedAt`、`instanceCounts`、`taskCounts`、`timeoutTaskCount`、`avgCompletionSeconds`、`pendingBindingFailures`、`businessProjectionCounts`、`pendingBusinessProjections` |
| `admin.BusinessProjection` | `projectionId`、`tenantId`、`flowId`、`flowVersionId`、`ownerInstanceId`、`appliedOwnerInstanceId`、`businessTable`、`recordKey`、`consistency`、`desiredStatus`、`desiredStartedAt`、`desiredFinishedAt`、`desiredRevision`、`appliedRevision`、`status`、`attemptCount`、`nextAttemptAt`、`leaseUntil`、`lastError`、`appliedAt`、`updatedAt` |

自助响应使用 `approval/my` 包中的 DTO：

| DTO | JSON 字段 |
| --- | --- |
| `my.AvailableFlow` | `flowId`、`flowCode`、`flowName`、`flowIcon`、`description`、`categoryId`、`categoryName` |
| `my.InitiatedInstance` | `instanceId`、`instanceNo`、`title`、`flowName`、`flowIcon`、`status`、`currentNodeName`、`createdAt`、`finishedAt` |
| `my.PendingTask` | `taskId`、`instanceId`、`instanceTitle`、`instanceNo`、`flowName`、`flowIcon`、`applicant`（`UserInfo`）、`nodeName`、`createdAt`、`deadline`、`isTimeout` |
| `my.CompletedTask` | `taskId`、`instanceId`、`instanceTitle`、`instanceNo`、`flowName`、`flowIcon`、`applicant`（`UserInfo`）、`nodeName`、`status`、`finishedAt` |
| `my.CCRecord` | `ccRecordId`、`instanceId`、`instanceTitle`、`instanceNo`、`flowName`、`flowIcon`、`applicant`（`UserInfo`）、`nodeName`、`isRead`、`createdAt` |
| `my.PendingCounts` | `pendingTaskCount`、`unreadCcCount` |
| `my.InstanceDetail` | `instance`、`formSchema`（宿主设计器文档，原样返回）、`timeline`、`flowGraph`、`availableActions`、`fieldPermissions` |
| `my.InstanceInfo` | `instanceId`、`instanceNo`、`title`、`flowName`、`flowIcon`、`applicant`、`status`、`currentNodeId`、`currentNodeName`、`businessRef`、`formData`、`createdAt`、`finishedAt` |

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
| `40011` | `ErrCodeBindingIncomplete` | `ErrBindingIncomplete` | `approval_binding_incomplete` | business binding 缺少必需的 table / key / status / instance-id 字段 |
| `40012` | `ErrCodeInvalidBindingMode` | `ErrInvalidBindingMode` | `approval_invalid_binding_mode` | flow binding mode 超出枚举范围 |
| `40013` | `ErrCodeInvalidInitiatorKind` | `ErrInvalidInitiatorKind` | `approval_invalid_initiator_kind` | flow initiator kind 超出枚举范围 |
| `40014` | `ErrCodeInvalidStorageMode` | `ErrInvalidStorageMode` | `approval_invalid_storage_mode` | deploy 请求的 storage mode 不是 `json` 或 `table` |
| `40015` | `ErrCodeFlowBindingLocked` | `ErrFlowBindingLocked` | `approval_flow_binding_locked` | 作为稳定错误面保留；版本固定的绑定快照意味着当前 flow 命令不再返回它 |
| `40016` | `ErrCodeBindingColumnsConflict` | `ErrBindingColumnsConflict` | `approval_binding_columns_conflict` | 两个 business-binding 字段指向了同一个列 |
| `40017` | `ErrCodeBindingUnexpected` | `ErrBindingUnexpected` | `approval_binding_unexpected` | standalone flow 却提供了 business binding 配置 |
| `40018` | `ErrCodeBindingSchemaInvalid` | `ErrBindingSchemaInvalid` | `approval_binding_schema_invalid` | 配置的绑定表或列在主库中不存在 |
| `40019` | `ErrCodeBindingKeyNotUnique` | `ErrBindingKeyNotUnique` | `approval_binding_key_not_unique` | key 列没有对应一个完整的非空主键或唯一键 |
| `40020` | `ErrCodeBindingStatusMappingInvalid` | `ErrBindingStatusMappingInvalid` | `approval_binding_status_mapping_invalid` | status mapping 含未知状态或映射到空白值 |
| `40101` | `ErrCodeInstanceNotFound` | `ErrInstanceNotFound` | `approval_instance_not_found` | instance 查找失败 |
| `40102` | `ErrCodeInstanceCompleted` | `ErrInstanceCompleted` | `approval_instance_completed` | instance 已完成 |
| `40103` | `ErrCodeNotAllowedInitiate` | `ErrNotAllowedInitiate` | `approval_not_allowed_initiate` | 调用者不能发起该 flow |
| `40104` | `ErrCodeWithdrawNotAllowed` | `ErrWithdrawNotAllowed` | `approval_withdraw_not_allowed` | 当前状态不允许撤回 |
| `40105` | `ErrCodeResubmitNotAllowed` | `ErrResubmitNotAllowed` | `approval_resubmit_not_allowed` | 当前状态不允许重新提交 |
| `40106` | `ErrCodeInvalidInstanceTransition` | `ErrInvalidInstanceTransition` | `approval_invalid_instance_transition` | instance 状态流转无效 |
| `40107` | `ErrCodeBusinessRefRequired` | `ErrBusinessRefRequired` | `approval_business_ref_required` | business 绑定流程发起时缺少业务引用 |
| `40108` | `ErrCodeBindingTargetBusy` | `ErrBindingTargetBusy` | `approval_binding_target_busy` | 业务记录已被一个未终结的审批实例占有 |
| `40109` | `ErrCodeInvalidBusinessRef` | `ErrInvalidBusinessRef` | `approval_invalid_business_ref` | 业务引用无法解析成配置的记录键 |
| `40110` | `ErrCodeBindingProjectionNotFound` | `ErrBindingProjectionNotFound` | `approval_binding_projection_not_found` | 投影查找失败（管理端重试） |
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
| `40401` | `ErrCodeFormValidationFailed` | 动态表单校验 `result.Err` | `approval_form_field_not_defined`, `approval_form_field_required`, `approval_form_field_must_be_string`, `approval_form_field_must_be_number`, `approval_form_field_must_be_integer`, `approval_form_field_min_length`, `approval_form_field_max_length`, `approval_form_field_invalid_validation`, `approval_form_field_pattern_mismatch`, `approval_form_field_min_value`, `approval_form_field_max_value`, `approval_form_field_empty`, `approval_form_field_invalid_file_item`, `approval_form_field_must_be_file`, `approval_form_field_invalid_value`, `approval_form_field_must_be_row_list`, `approval_form_field_must_be_row_object`, `approval_form_field_min_rows`, `approval_form_field_max_rows`, `approval_form_field_table_cell` | 字段级校验消息在运行时构造 |
| `40601` | `ErrCodeUrgeCooldown` | 动态 urge `result.Err` | `approval_urge_too_frequent` | 没有静态 sentinel；消息会带 `minutes`；非正数 `urgeCooldownMinutes` 默认按 30 分钟处理 |
| `40701` | `ErrCodeAccessDenied` | `ErrAccessDenied` | `approval_access_denied` | 调用者缺少审批域访问权限 |
| `40702` | `ErrCodeTerminateNotAllowed` | `ErrTerminateNotAllowed` | `approval_terminate_not_allowed` | 当前 instance 状态不允许终止 |

`ErrEventRouteNotTransactional`、`ErrEventRouteNotSubscribable`、
`ErrTenantNotResolved` 这类启动和租户解析诊断位于 `internal/approval/...`；
它们不是可 import 的公开 Go API，但 event routing 或 principal tenant
details 配错时，操作员可能会看到它们被包装后的错误消息。

---

下一步：[流程设计](./flow-design.md) 了解 `deploy` 背后的设计器 wire shape，或 [实例运行时](./runtime.md) 了解各实例操作背后的生命周期语义。
