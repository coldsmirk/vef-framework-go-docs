---
sidebar_position: 3
---

# Approval Module

The `approval` module provides a complete workflow engine for building approval-based business processes. It supports visual flow design (React Flow compatible), multi-level approval chains, conditional branching, parallel approval, delegation, rollbacks, and transactional event publishing.

## Enabling the Module

Approval is an optional feature module. It is intentionally absent from the
default `vef.Run(...)` boot graph, so applications that do not need workflow
support do not register its API resources, CQRS handlers, engine, binding
listener, or timeout scanners.

Enable it explicitly:

```go
vef.Run(
    vef.ApprovalModule,
    app.Module,
)
```

Approval publishes `approval.*` events with `event.WithTx`, and its binding
listener subscribes to those events. The host application must route
`approval.*` to a transactional transport with a subscribable sink, such as an
outbox route whose sink is Redis Streams.

`InstanceBindingFailedEvent` is the exception to the transactional-route
startup check: it is emitted by the asynchronous binding listener after the
approval transaction has already committed. `InstanceCompletedEvent` has the
strictest route requirement because the binding listener subscribes to it; the
route must include a subscribable sink transport such as `memory` or
`redis_stream` alongside the transactional outbox route.

## RPC Resources

When enabled, the module registers these RPC resources. They are not public
operations: callers must be authenticated, and permissions below are enforced
when a `RequiredPermission` is shown.

RPC calls use the standard envelope documented in [API](../guide/api): `resource`,
`action`, `version`, `params`, and `meta`. Fields listed below as `params.*`
are decoded from `params`; fields listed as `meta.*` are decoded from `meta`.
The generated [Runtime API Index](../reference/runtime-api-index) contains the
exhaustive JSON field ledger for every request and response DTO.

The grouped-family audit also pins 766 grouped approval field/method entries:
607 approval package entries, 76 approval/admin DTO field entries, and 83
approval/my DTO field entries. These entries cover the public Go DTO fields,
domain event fields, node-data helpers, lifecycle hooks, resolver interfaces,
and status helper methods whose exact signatures are listed in the public API
index; the verifier locks their sorted signatures and receiver/type
distribution.

### `approval/category`

| Action | Permission | Params | Notes |
| --- | --- | --- | --- |
| `find_tree` | `approval:category:query` | `CategorySearch` | Tenant-scoped tree query |
| `find_tree_options` | `approval:category:query` | `CategorySearch` + `DataOptionConfig` | Tenant-scoped tree options |
| `create` | `approval:category:create` | `CategoryParams` | Non-super-admin tenant is stamped from caller |
| `update` | `approval:category:update` | `CategoryParams` | Non-super-admin can only mutate own tenant |
| `delete` | `approval:category:delete` | Primary-key params | Non-super-admin can only delete own tenant |

| Action | Request fields |
| --- | --- |
| `find_tree` | `meta.name`, `meta.isActive`, `meta.sort` |
| `find_tree_options` | `meta.name`, `meta.isActive`, `meta.sort`, plus option mapping metadata: `meta.labelColumn`, `meta.valueColumn`, `meta.descriptionColumn`, `meta.metaColumns` |
| `create` | `params.id`, `params.tenantId` required, `params.code` required, `params.name` required, `params.icon`, `params.parentId`, `params.sortOrder`, `params.isActive`, `params.remark` |
| `update` | `params.id`, `params.tenantId` required, `params.code` required, `params.name` required, `params.icon`, `params.parentId`, `params.sortOrder`, `params.isActive`, `params.remark` |
| `delete` | `params.id` required |

### `approval/delegation`

| Action | Permission | Params | Notes |
| --- | --- | --- | --- |
| `find_page` | `approval:delegation:query` | `DelegationSearch` + pageable meta | Non-super-admin sees own delegations |
| `create` | `approval:delegation:create` | `DelegationParams` | Non-super-admin delegator is stamped from caller |
| `update` | `approval:delegation:update` | `DelegationParams` | Non-super-admin cannot reassign ownership |
| `delete` | `approval:delegation:delete` | Primary-key params | Owner-scoped for non-super-admin |

| Action | Request fields |
| --- | --- |
| `find_page` | `meta.delegatorId`, `meta.delegateeId`, `meta.isActive`, `meta.sort`, `meta.page`, `meta.size` |
| `create` | `params.id`, `params.delegatorId` required, `params.delegateeId` required, `params.flowCategoryId`, `params.flowId`, `params.startTime`, `params.endTime`, `params.isActive`, `params.reason` |
| `update` | `params.id`, `params.delegatorId` required, `params.delegateeId` required, `params.flowCategoryId`, `params.flowId`, `params.startTime`, `params.endTime`, `params.isActive`, `params.reason` |
| `delete` | `params.id` required |

### `approval/flow`

| Action | Permission | Params | Notes |
| --- | --- | --- | --- |
| `create` | `approval:flow:create` | `CreateFlowParams` | Audited |
| `deploy` | `approval:flow:deploy` | `DeployFlowParams` | Audited |
| `publish_version` | `approval:flow:publish` | `PublishVersionParams` | Audited |
| `update_flow` | `approval:flow:update` | `UpdateFlowParams` | Audited |
| `toggle_active` | `approval:flow:update` | `ToggleActiveParams` | Audited |
| `get_graph` | `approval:flow:query` | `GetGraphParams` | Reads published graph |
| `find_flows` | `approval:flow:query` | `FindFlowsParams` | Paged query |
| `find_versions` | `approval:flow:query` | `FindVersionsParams` | Lists versions for one flow |

| Action | Request fields |
| --- | --- |
| `create` | `params.tenantId` required, `params.code` required, `params.name` required, `params.categoryId` required, `params.bindingMode` required, `params.icon`, `params.description`, `params.businessTable`, `params.businessPkField`, `params.businessTitleField`, `params.businessStatusField`, `params.adminUserIds`, `params.isAllInitiationAllowed`, `params.instanceTitleTemplate`, `params.initiators` |
| `deploy` | `params.flowId` required, `params.description`, `params.flowDefinition` required, `params.formDefinition` |
| `publish_version` | `params.versionId` required |
| `update_flow` | `params.flowId` required, `params.name` required, `params.instanceTitleTemplate` required, `params.icon`, `params.description`, `params.adminUserIds`, `params.isAllInitiationAllowed`, `params.initiators` |
| `toggle_active` | `params.flowId` required, `params.isActive` |
| `get_graph` | `params.flowId` required, `params.tenantId` |
| `find_flows` | `params.tenantId`, `params.categoryId`, `params.keyword`, `params.isActive`, `params.page`, `params.pageSize` |
| `find_versions` | `params.flowId` required, `params.tenantId` |

`params.initiators` entries use `kind` (`user`, `role`, or `department`) and
`ids`. For business binding, `businessTable`, `businessPkField`,
`businessTitleField`, and `businessStatusField` are SQL identifiers validated
by `ValidateBusinessIdentifier`.

### `approval/instance`

| Action | Permission | Params | Notes |
| --- | --- | --- | --- |
| `start` | `approval:instance:start` | `StartInstanceParams` | Audited |
| `process_task` | `approval:task:process` | `ProcessTaskParams` | Audited; `action` must be `approve`, `reject`, `transfer`, `rollback`, or `handle` |
| `withdraw` | `approval:instance:withdraw` | `WithdrawParams` | Audited |
| `resubmit` | `approval:instance:resubmit` | `ResubmitParams` | Audited; accepts returned or withdrawn instances |
| `add_cc` | `approval:instance:cc` | `AddCCParams` | Audited |
| `mark_cc_read` | `approval:instance:cc` | `MarkCCReadParams` | Read receipt, not audited |
| `add_assignee` | `approval:task:add_assignee` | `AddAssigneeParams` | Audited; `addType` is `before`, `after`, or `parallel` |
| `remove_assignee` | `approval:task:remove_assignee` | `RemoveAssigneeParams` | Audited |
| `urge_task` | `approval:task:urge` | `UrgeTaskParams` | Extra rate limit: max `10` per `1m` |

| Action | Request fields |
| --- | --- |
| `start` | `params.tenantId` required, `params.flowCode` required, `params.businessRecordId`, `params.formData` |
| `process_task` | `params.taskId` required, `params.action` required (`approve`, `reject`, `transfer`, `rollback`, or `handle`), `params.opinion` max 2000 chars, `params.formData`, `params.transferToId`, `params.targetNodeId` |
| `withdraw` | `params.instanceId` required, `params.reason` max 2000 chars |
| `resubmit` | `params.instanceId` required, `params.formData` |
| `add_cc` | `params.instanceId` required, `params.ccUserIds` required, 1-50 IDs |
| `mark_cc_read` | `params.instanceId` required |
| `add_assignee` | `params.taskId` required, `params.userIds` required, 1-50 IDs, `params.addType` required (`before`, `after`, or `parallel`) |
| `remove_assignee` | `params.taskId` required |
| `urge_task` | `params.taskId` required, `params.message` max 500 chars |

### `approval/my`

Self-service queries do not declare `RequiredPermission`, but they still require
the authenticated principal.

| Action | Params | Output |
| --- | --- | --- |
| `find_available_flows` | `FindAvailableFlowsParams` | `page.Page[my.AvailableFlow]` |
| `find_initiated` | `FindInitiatedParams` | `page.Page[my.InitiatedInstance]` |
| `find_pending_tasks` | `FindPendingTasksParams` | `page.Page[my.PendingTask]` |
| `find_completed_tasks` | `FindCompletedTasksParams` | `page.Page[my.CompletedTask]` |
| `find_cc_records` | `FindCCRecordsParams` | `page.Page[my.CCRecord]` |
| `get_pending_counts` | `GetPendingCountsParams` | `my.PendingCounts` |
| `get_instance_detail` | `GetInstanceDetailParams` | `my.InstanceDetail` |

| Action | Request fields |
| --- | --- |
| `find_available_flows` | `params.tenantId`, `params.keyword`, `params.page`, `params.pageSize` |
| `find_initiated` | `params.tenantId`, `params.status`, `params.keyword`, `params.page`, `params.pageSize` |
| `find_pending_tasks` | `params.tenantId`, `params.page`, `params.pageSize` |
| `find_completed_tasks` | `params.tenantId`, `params.page`, `params.pageSize` |
| `find_cc_records` | `params.tenantId`, `params.isRead`, `params.page`, `params.pageSize` |
| `get_pending_counts` | `params.tenantId` |
| `get_instance_detail` | `params.instanceId` required |

The `my.InstanceDetail` JSON payload includes `taskId`, `formData`, and
`actionLogs` fields for task references, submitted form data, and action-log
history.

`availableActions` is a query-layer UI hint. For the applicant it includes
`withdraw` when the instance is `running`, and `resubmit` when the instance is
`rejected` or `returned`. For pending tasks it includes `handle` for handle
nodes, otherwise `approve`, then `reject`, plus `transfer`, `rollback`,
`add_assignee`, or `add_cc` when the current node allows them. If the instance
has any pending task, it also includes `urge`. Command handlers still perform
their own validation.

### `approval/admin`

| Action | Permission | Params | Notes |
| --- | --- | --- | --- |
| `find_instances` | `approval:instance:query` | `AdminFindInstancesParams` | Tenant-filtered for non-super-admin |
| `find_tasks` | `approval:task:query` | `AdminFindTasksParams` | Tenant-filtered for non-super-admin |
| `get_instance_detail` | `approval:instance:detail` | `AdminGetInstanceDetailParams` | Full admin detail |
| `find_action_logs` | `approval:log:query` | `AdminFindActionLogsParams` | Requires `instanceId` |
| `get_metrics` | `approval:metrics:query` | `AdminGetMetricsParams` | Aggregated metrics |
| `terminate_instance` | `approval:instance:terminate` | `AdminTerminateInstanceParams` | Audited |
| `reassign_task` | `approval:task:reassign` | `AdminReassignTaskParams` | Audited |

| Action | Request fields |
| --- | --- |
| `find_instances` | `params.tenantId`, `params.applicantId`, `params.status`, `params.flowId`, `params.keyword`, `params.page`, `params.pageSize` |
| `find_tasks` | `params.tenantId`, `params.assigneeId`, `params.instanceId`, `params.status`, `params.page`, `params.pageSize` |
| `get_instance_detail` | `params.instanceId` required |
| `find_action_logs` | `params.instanceId` required, `params.tenantId`, `params.page`, `params.pageSize` |
| `get_metrics` | `params.tenantId` |
| `terminate_instance` | `params.instanceId` required, `params.reason` max 2000 chars |
| `reassign_task` | `params.taskId` required, `params.newAssigneeId` required, `params.reason` max 2000 chars |

For admin list and metrics queries, non-super-admin callers ignore a submitted
`tenantId` override and are filtered to their own tenant. Super-admin callers
may pass `tenantId` to filter one tenant or omit it for cross-tenant visibility.

Admin list/detail DTOs keep the tenant and history fields on the wire as
`tenantId` and `actionLogs`.

### Response DTO Fields

Admin responses use the DTOs from `approval/admin`:

| DTO | JSON fields |
| --- | --- |
| `admin.Instance` | `instanceId`, `instanceNo`, `title`, `tenantId`, `flowId`, `flowName`, `applicantId`, `applicantName`, `status`, `currentNodeName`, `createdAt`, `finishedAt` |
| `admin.Task` | `taskId`, `instanceId`, `instanceTitle`, `flowName`, `nodeName`, `assigneeId`, `assigneeName`, `status`, `createdAt`, `deadline`, `finishedAt` |
| `admin.InstanceDetail` | `instance`, `tasks`, `actionLogs`, `flowNodes` |
| `admin.InstanceDetailInfo` | `instanceId`, `instanceNo`, `title`, `tenantId`, `flowId`, `flowName`, `flowVersionId`, `applicantId`, `applicantName`, `status`, `currentNodeName`, `businessRecordId`, `formData`, `createdAt`, `finishedAt` |
| `admin.TaskDetailInfo` | `taskId`, `nodeId`, `nodeName`, `assigneeId`, `assigneeName`, `delegatorId`, `delegatorName`, `status`, `sortOrder`, `deadline`, `isTimeout`, `createdAt`, `finishedAt` |
| `admin.ActionLog` | `logId`, `action`, `operatorId`, `operatorName`, `operatorDepartmentName`, `transferToId`, `transferToName`, `opinion`, `createdAt` |
| `admin.FlowNodeInfo` | `nodeId`, `key`, `kind`, `name`, `executionType` |
| `admin.Metrics` | `tenantId`, `capturedAt`, `instanceCounts`, `taskCounts`, `timeoutTaskCount`, `avgCompletionSeconds`, `pendingBindingFailures` |

Self-service responses use the DTOs from `approval/my`:

| DTO | JSON fields |
| --- | --- |
| `my.AvailableFlow` | `flowId`, `flowCode`, `flowName`, `flowIcon`, `description`, `categoryId`, `categoryName` |
| `my.InitiatedInstance` | `instanceId`, `instanceNo`, `title`, `flowName`, `flowIcon`, `status`, `currentNodeName`, `createdAt`, `finishedAt` |
| `my.PendingTask` | `taskId`, `instanceId`, `instanceTitle`, `instanceNo`, `flowName`, `flowIcon`, `applicantName`, `nodeName`, `createdAt`, `deadline`, `isTimeout` |
| `my.CompletedTask` | `taskId`, `instanceId`, `instanceTitle`, `instanceNo`, `flowName`, `flowIcon`, `applicantName`, `nodeName`, `status`, `finishedAt` |
| `my.CCRecord` | `ccRecordId`, `instanceId`, `instanceTitle`, `instanceNo`, `flowName`, `flowIcon`, `applicantName`, `nodeName`, `isRead`, `createdAt` |
| `my.PendingCounts` | `pendingTaskCount`, `unreadCcCount` |
| `my.InstanceDetail` | `instance`, `tasks`, `actionLogs`, `flowNodes`, `availableActions` |
| `my.InstanceInfo` | `instanceId`, `instanceNo`, `title`, `flowName`, `flowIcon`, `applicantId`, `applicantName`, `status`, `currentNodeName`, `businessRecordId`, `formData`, `createdAt`, `finishedAt` |
| `my.TaskInfo` | `taskId`, `nodeName`, `assigneeId`, `assigneeName`, `status`, `sortOrder`, `createdAt`, `finishedAt` |
| `my.ActionLogInfo` | `action`, `operatorName`, `opinion`, `createdAt` |
| `my.FlowNodeInfo` | `nodeId`, `key`, `kind`, `name` |

## Architecture Overview

```
Flow Category → Flow → Flow Version → Nodes + Edges
                                        ↓
                                    Instance → Tasks → Action Logs
```

| Concept | Table | Description |
| --- | --- | --- |
| Flow Category | `apv_flow_category` | Hierarchical grouping of flows |
| Flow | `apv_flow` | A workflow definition (e.g., "Leave Request") |
| Flow Version | `apv_flow_version` | Versioned snapshot with nodes, edges, and form schema |
| Flow Node | `apv_flow_node` | A step in the workflow (approval, handle, condition, CC) |
| Flow Edge | `apv_flow_edge` | Directed connection between nodes |
| Instance | `apv_instance` | A running instance of a flow |
| Task | `apv_task` | An individual approval/handle task assigned to a user |
| Action Log | `apv_action_log` | Audit trail of all actions |

## Configuration

```toml
[vef.approval]
auto_migrate              = true
timeout_scan_interval     = "1m"
pre_warning_scan_interval = "5m"
cleanup_scan_interval     = "24h"
delegation_max_depth      = 10
form_snapshot_retention   = "2160h"  # 90 days
urge_record_retention     = "720h"   # 30 days
cc_record_retention       = "2160h"  # 90 days
```

`auto_migrate` is a plain boolean switch and is not set by
`ApprovalConfig.ApplyDefaults()`: enable it explicitly when the app should run
approval DDL on startup. `cc_record_retention` only prunes CC records that have
already been read.

> The outbox-related fields that previously lived under `[vef.approval]` (`outbox_relay_interval`, `outbox_max_retries`, `outbox_batch_size`) moved to `[vef.event.transports.outbox]` in v0.21. The approval module now publishes through the framework-wide outbox transport — see [Event Bus](../features/event-bus). Approval's binding listener and outbox publisher both assert routing at boot via `event.RouteInspector`, so a misconfigured route fails the application instead of degrading silently.

See [Configuration Reference](../reference/configuration-reference) for details.

## Binding Modes

| Mode | Constant | Wire value | Description |
| --- | --- | --- | --- |
| Standalone | `BindingStandalone` | `standalone` | Form data stored in the approval module's own tables |
| Business | `BindingBusiness` | `business` | Links to an existing business data table |

Business binding connects the approval flow to your domain tables via `BusinessTable`, `BusinessPkField`, `BusinessTitleField`, and `BusinessStatusField`.

## Node Types

| Node Kind | Constant | Wire value | Description |
| --- | --- | --- | --- |
| Start | `NodeStart` | `start` | Entry point of the workflow |
| Approval | `NodeApproval` | `approval` | Requires approval action from assignees |
| Handle | `NodeHandle` | `handle` | Requires processing/handling action |
| Condition | `NodeCondition` | `condition` | Branches based on conditions |
| CC | `NodeCC` | `cc` | Sends notifications to specified users |
| End | `NodeEnd` | `end` | Terminal point of the workflow |

## Condition Branching

Condition nodes evaluate `ConditionBranch` entries in priority order. Each
branch contains one or more `ConditionGroup` values: conditions inside a group
are combined with AND logic, while multiple groups on the same branch are
combined with OR logic.

`ConditionField` uses the structured `Subject` / `Operator` / `Value` fields.
The built-in evaluator converts it to an `expr-lang` expression. Supported
operators are `eq`, `ne`, `gt`, `gte`, `lt`, `lte`, `in`, `not_in`,
`contains`, `not_contains`, `starts_with`, `ends_with`, `is_empty`, and
`is_not_empty`; unknown operators evaluate to `false`.

`ConditionExpression` evaluates the raw `Expression` string with `expr-lang`.
The evaluation environment exposes:

| Name | Value |
| --- | --- |
| `formData` | the instance `FormData` as a map |
| `applicantId` | current applicant ID |
| `applicantDepartmentId` | applicant department ID, or `""` when absent |

Approval conditions intentionally use `expr-lang` directly, not the public
`expression.Engine` feature. That keeps approval workflows pure-Go today; the
current public `expression.Engine` backend is Zen and requires CGO.

## Approval Methods

When a node has multiple assignees:

| Method | Constant | Wire value | Behavior |
| --- | --- | --- | --- |
| Sequential | `ApprovalSequential` | `sequential` | Approvers process one by one in order |
| Parallel | `ApprovalParallel` | `parallel` | Approvers process simultaneously |

The enum type is `ApprovalMethod`.

### Pass Rules (for Parallel)

| Rule | Constant | Wire value | Behavior |
| --- | --- | --- | --- |
| All | `PassAll` | `all` | All assignees must approve |
| Any | `PassAny` | `any` | At least one approval passes |
| Ratio | `PassRatio` | `ratio` | A percentage must approve |
| Any Reject | `PassAnyReject` | `any_reject` | Any rejection fails the node |

Custom pass-rule implementations use `PassRuleStrategy`, `PassRuleContext`,
and return a `PassRuleResult` (`PassRulePending`, `PassRulePassed`,
`PassRuleRejected`).

## Assignee Types

| Kind | Constant | Wire value | Description |
| --- | --- | --- | --- |
| User | `AssigneeUser` | `user` | Specific users |
| Role | `AssigneeRole` | `role` | Users with a role |
| Department | `AssigneeDepartment` | `department` | Department head |
| Self | `AssigneeSelf` | `self` | The applicant |
| Superior | `AssigneeSuperior` | `superior` | Direct superior |
| Dept Leader | `AssigneeDepartmentLeader` | `department_leader` | Multi-level supervisor chain |
| Form Field | `AssigneeFormField` | `form_field` | Determined by a form field value |

The enum type is `AssigneeKind`. Dynamic assignee insertion uses
`AddAssigneeType`: `AddAssigneeBefore` (`before`), `AddAssigneeAfter`
(`after`), and `AddAssigneeParallel` (`parallel`).

## Instance Lifecycle

```
submit → Running → approve/reject → Approved/Rejected
                 → withdraw       → Withdrawn
                 → rollback       → Returned
                 → terminate      → Terminated
Withdrawn/Returned → resubmit     → Running (again)
```

The runtime state machine declares only these valid instance transitions:

| From | To |
| --- | --- |
| `running` | `approved` |
| `running` | `rejected` |
| `running` | `withdrawn` |
| `running` | `terminated` |
| `running` | `returned` |
| `returned` | `running` |
| `withdrawn` | `running` |

### Instance Statuses

| Status | Constant | Wire value | Final? |
| --- | --- | --- | --- |
| Running | `InstanceRunning` | `running` | No |
| Approved | `InstanceApproved` | `approved` | Yes |
| Rejected | `InstanceRejected` | `rejected` | Yes |
| Withdrawn | `InstanceWithdrawn` | `withdrawn` | No |
| Returned | `InstanceReturned` | `returned` | No |
| Terminated | `InstanceTerminated` | `terminated` | Yes |

The enum type is `InstanceStatus`.

### Task Statuses

| Status | Constant | Wire value | Final? |
| --- | --- | --- | --- |
| Waiting | `TaskWaiting` | `waiting` | No |
| Pending | `TaskPending` | `pending` | No |
| Approved | `TaskApproved` | `approved` | Yes |
| Rejected | `TaskRejected` | `rejected` | Yes |
| Handled | `TaskHandled` | `handled` | Yes |
| Transferred | `TaskTransferred` | `transferred` | Yes |
| Rolled Back | `TaskRolledBack` | `rolled_back` | Yes |
| Canceled | `TaskCanceled` | `canceled` | Yes |
| Removed | `TaskRemoved` | `removed` | Yes |
| Skipped | `TaskSkipped` | `skipped` | Yes |

The runtime state machine declares only these valid task transitions:

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

## Actions

| Action | Constant | Wire value | Description |
| --- | --- | --- | --- |
| Submit | `ActionSubmit` | `submit` | Start a new instance |
| Approve | `ActionApprove` | `approve` | Approve a task |
| Handle | `ActionHandle` | `handle` | Complete a handle task |
| Reject | `ActionReject` | `reject` | Reject a task |
| Transfer | `ActionTransfer` | `transfer` | Transfer to another user |
| Withdraw | `ActionWithdraw` | `withdraw` | Applicant withdraws |
| Cancel | `ActionCancel` | `cancel` | Cancel a task |
| Rollback | `ActionRollback` | `rollback` | Roll back to a previous node |
| Add Assignee | `ActionAddAssignee` | `add_assignee` | Dynamically add an assignee |
| Remove Assignee | `ActionRemoveAssignee` | `remove_assignee` | Remove an assignee |
| Add CC | `ActionAddCC` | `add_cc` | Dynamically add CC recipients |
| Execute | `ActionExecute` | `execute` | Internal execution action for automatic node handling |
| Resubmit | `ActionResubmit` | `resubmit` | Resubmit a returned or withdrawn instance |
| Reassign | `ActionReassign` | `reassign` | Admin reassigns a task |
| Terminate | `ActionTerminate` | `terminate` | Admin force-terminates |

## Rollback Configuration

| Property | Options |
| --- | --- |
| `RollbackType` | `RollbackNone` (`none`), `RollbackPrevious` (`previous`), `RollbackStart` (`start`), `RollbackAny` (`any`), `RollbackSpecified` (`specified`) |
| `RollbackDataStrategy` | `RollbackDataClear` (`clear`, reset form), `RollbackDataKeep` (`keep`, preserve data) |

Same-applicant handling uses `SameApplicantAction` with
`SameApplicantSelfApprove` (`self_approve`), `SameApplicantAutoPass`
(`auto_pass`), and `SameApplicantTransferSuperior` (`transfer_superior`).
Consecutive-approver handling uses `ConsecutiveApproverAction` with
`ConsecutiveApproverNone` (`none`) and `ConsecutiveApproverAutoPass`
(`auto_pass`).

## Empty Assignee Handling

When no assignee is found for a node:

| Action | Constant | Wire value |
| --- | --- | --- |
| Auto-pass | `EmptyAssigneeAutoPass` | `auto_pass` |
| Transfer to admin | `EmptyAssigneeTransferAdmin` | `transfer_admin` |
| Transfer to superior | `EmptyAssigneeTransferSuperior` | `transfer_superior` |
| Transfer to applicant | `EmptyAssigneeTransferApplicant` | `transfer_applicant` |
| Transfer to specified | `EmptyAssigneeTransferSpecified` | `transfer_specified` |

The enum type is `EmptyAssigneeAction`.

## Timeout Handling

| Action | Constant | Wire value | Behavior |
| --- | --- | --- | --- |
| None | `TimeoutActionNone` | `none` | Mark timeout only |
| Auto Pass | `TimeoutActionAutoPass` | `auto_pass` | Automatically approve |
| Auto Reject | `TimeoutActionAutoReject` | `auto_reject` | Automatically reject |
| Notify | `TimeoutActionNotify` | `notify` | Send notification only |
| Transfer Admin | `TimeoutActionTransferAdmin` | `transfer_admin` | Transfer to node admin |

The enum type is `TimeoutAction`.

## Form Data Storage

| Mode | Constant | Wire value | Location |
| --- | --- | --- | --- |
| JSON | `StorageJSON` | `json` | `apv_instance.form_data` (JSONB column) |

`StorageJSON` is the only storage mode currently exported by the package.

Flow design and persistence models exposed by the public package include
`FlowCategory`, `Flow`, `FlowVersion`, `FlowNode`, `FlowEdge`, `FlowInitiator`,
`FlowNodeAssignee`, `FlowNodeCC`, `FormDefinition`, `FormFieldDefinition`,
`FormSnapshot`, `ActionLog`, `OperatorInfo`, and `UrgeRecord`. Flow-version
status uses `VersionStatus`: `VersionDraft` (`draft`), `VersionPublished`
(`published`), and `VersionArchived` (`archived`).

Additional flow-designer enums:

| Enum | Wire values |
| --- | --- |
| `InitiatorKind` | `user`, `role`, `department` |
| `ExecutionType` | `manual`, `auto`, `auto_pass`, `auto_reject` |
| `ConditionKind` | `field`, `expression` |
| `CCKind` | `user`, `role`, `department`, `form_field` |
| `CCTiming` | `always`, `on_approve`, `on_reject` |
| `FieldKind` | `input`, `textarea`, `select`, `number`, `date`, `upload` |
| `Permission` | `visible`, `editable`, `hidden`, `required` |

## Event Publication

The approval module publishes its domain events through the framework's transactional outbox transport (see [Event Bus](../features/event-bus)). Every approval command writes the event record in the same transaction as the business mutation; the outbox relay then forwards them to the configured sink.

Subscribers must:

1. Attach with `event.WithGroup("...")` because the route resolves to an at-least-once transport.
2. Rely on the Inbox middleware for dedupe (it activates automatically when `event.middleware.inbox = true` and a transport advertises `AtLeastOnce`).

> The standalone `apv_event_outbox` table and module-private `EventOutboxStatus` constants from earlier snapshots have been retired. Approval no longer carries a private outbox — it composes with the framework one.

### Domain Event Types

All approval events implement `event.Event` and embed enough payload to drive
integrations without a follow-up read. Common event JSON fields include
`tenantId` and `occurredTime`; the tables below list the event topic wire value
and the event-specific payload fields.

Instance lifecycle:

| Type constant | Topic | Payload / constructor | Payload fields beyond common fields | When |
| --- | --- | --- | --- | --- |
| `EventTypeInstanceCreated` | `approval.instance.created` | `InstanceCreatedEvent`, `NewInstanceCreatedEvent` | `instanceId`, `flowId`, `title`, `applicantId`, `applicantName` | a new instance was started |
| `EventTypeInstanceCompleted` | `approval.instance.completed` | `InstanceCompletedEvent`, `NewInstanceCompletedEvent` | `instanceId`, `finalStatus`, `finishedAt` | instance reached a terminal status |
| `EventTypeInstanceWithdrawn` | `approval.instance.withdrawn` | `InstanceWithdrawnEvent`, `NewInstanceWithdrawnEvent` | `instanceId`, `operatorId` | applicant withdrew the instance |
| `EventTypeInstanceRolledBack` | `approval.instance.rolled_back` | `InstanceRolledBackEvent`, `NewInstanceRolledBackEvent` | `instanceId`, `fromNodeId`, `toNodeId`, `operatorId` | instance was rolled back to a previous node |
| `EventTypeInstanceReturned` | `approval.instance.returned` | `InstanceReturnedEvent`, `NewInstanceReturnedEvent` | `instanceId`, `fromNodeId`, `toNodeId`, `operatorId` | instance was returned to applicant |
| `EventTypeInstanceResubmitted` | `approval.instance.resubmitted` | `InstanceResubmittedEvent`, `NewInstanceResubmittedEvent` | `instanceId`, `operatorId` | returned or withdrawn instance was resubmitted |
| `EventTypeInstanceBindingFailed` | `approval.instance.binding_failed` | `InstanceBindingFailedEvent`, `NewInstanceBindingFailedEvent` | `instanceId`, `flowId`, `finalStatus`, `businessTable`, `errorMessage` | the binding listener could not write back the final status to the business row |

Node lifecycle:

| Type constant | Topic | Payload / constructor | Payload fields beyond common fields | When |
| --- | --- | --- | --- | --- |
| `EventTypeNodeEntered` | `approval.node.entered` | `NodeEnteredEvent`, `NewNodeEnteredEvent` | `instanceId`, `nodeId`, `nodeName` | engine activated a node |
| `EventTypeNodeAutoPassed` | `approval.node.auto_passed` | `NodeAutoPassedEvent`, `NewNodeAutoPassedEvent` | `instanceId`, `nodeId`, `reason` | a node auto-passed because no assignees were found |

Task lifecycle:

| Type constant | Topic | Payload / constructor | Payload fields beyond common fields | When |
| --- | --- | --- | --- | --- |
| `EventTypeTaskCreated` | `approval.task.created` | `TaskCreatedEvent`, `NewTaskCreatedEvent` | `taskId`, `instanceId`, `nodeId`, `assigneeId`, `assigneeName`, `deadline` | a task was created; sequential follow-up tasks may start with `deadline` omitted while waiting |
| `EventTypeTaskApproved` | `approval.task.approved` | `TaskApprovedEvent`, `NewTaskApprovedEvent` | `taskId`, `instanceId`, `nodeId`, `operatorId`, `opinion` | task approved |
| `EventTypeTaskHandled` | `approval.task.handled` | `TaskHandledEvent`, `NewTaskHandledEvent` | `taskId`, `instanceId`, `nodeId`, `operatorId`, `opinion` | handle task completed |
| `EventTypeTaskRejected` | `approval.task.rejected` | `TaskRejectedEvent`, `NewTaskRejectedEvent` | `taskId`, `instanceId`, `nodeId`, `operatorId`, `opinion` | task rejected |
| `EventTypeTaskTransferred` | `approval.task.transferred` | `TaskTransferredEvent`, `NewTaskTransferredEvent` | `taskId`, `instanceId`, `nodeId`, `fromUserId`, `fromUserName`, `toUserId`, `toUserName`, `reason` | task transferred to another user |
| `EventTypeTaskReassigned` | `approval.task.reassigned` | `TaskReassignedEvent`, `NewTaskReassignedEvent` | `taskId`, `instanceId`, `nodeId`, `fromUserId`, `fromUserName`, `toUserId`, `toUserName`, `reason` | admin reassigned a task |
| `EventTypeTaskTimedOut` | `approval.task.timed_out` | `TaskTimedOutEvent`, `NewTaskTimedOutEvent` | `taskId`, `instanceId`, `nodeId`, `assigneeId`, `assigneeName`, `deadline` | timeout scanner fired the configured timeout action |
| `EventTypeAssigneesAdded` | `approval.task.assignees_added` | `AssigneesAddedEvent`, `NewAssigneesAddedEvent` | `instanceId`, `nodeId`, `taskId`, `addType`, `assigneeIds`, `assigneeNames` | dynamic assignees added |
| `EventTypeAssigneesRemoved` | `approval.task.assignees_removed` | `AssigneesRemovedEvent`, `NewAssigneesRemovedEvent` | `instanceId`, `nodeId`, `taskId`, `assigneeIds`, `assigneeNames` | dynamic assignees removed |
| `EventTypeTaskDeadlineWarning` | `approval.task.deadline_warning` | `TaskDeadlineWarningEvent`, `NewTaskDeadlineWarningEvent` | `taskId`, `instanceId`, `nodeId`, `assigneeId`, `assigneeName`, `deadline`, `hoursLeft` | pre-warning scanner flagged an approaching deadline |
| `EventTypeTaskUrged` | `approval.task.urged` | `TaskUrgedEvent`, `NewTaskUrgedEvent` | `instanceId`, `nodeId`, `taskId`, `urgerId`, `urgerName`, `targetUserId`, `targetUserName`, `message` | applicant urged an assignee |

CC + Flow:

| Type constant | Topic | Payload / constructor | Payload fields beyond common fields | When |
| --- | --- | --- | --- | --- |
| `EventTypeCCNotified` | `approval.cc.notified` | `CCNotifiedEvent`, `NewCCNotifiedEvent` | `instanceId`, `nodeId`, `ccUserIds`, `ccUserNames`, `isManual` | a CC node delivered notifications |
| `EventTypeFlowCreated` | `approval.flow.created` | `FlowCreatedEvent`, `NewFlowCreatedEvent` | `flowId`, `code`, `name`, `categoryId` | flow created |
| `EventTypeFlowUpdated` | `approval.flow.updated` | `FlowUpdatedEvent`, `NewFlowUpdatedEvent` | `flowId` | flow updated |
| `EventTypeFlowDeployed` | `approval.flow.deployed` | `FlowDeployedEvent`, `NewFlowDeployedEvent` | `flowId`, `versionId`, `version` | flow version deployed |
| `EventTypeFlowToggled` | `approval.flow.toggled` | `FlowToggledEvent`, `NewFlowToggledEvent` | `flowId`, `isActive` | flow active flag changed |
| `EventTypeFlowPublished` | `approval.flow.published` | `FlowPublishedEvent`, `NewFlowPublishedEvent` | `flowId`, `versionId` | flow version published |

## Caller Context and Multi-Tenant Safety

Resource and command handlers resolve a `CallerContext` per request that bundles the tenant authority of the call. Exactly one of `TenantID`, `IsSuperAdmin`, or `IsSystemInternal` must be set; a zero `CallerContext` is fail-closed (treated as unauthorized).

| Member | Meaning |
| --- | --- |
| `TenantID` | the tenant this caller acts within; all reads/writes get filtered by it |
| `IsSuperAdmin` | the principal carries `approval:super_admin`; cross-tenant queries allowed |
| `IsSystemInternal` | the call originates from framework-internal code (scanners, listeners) |

Cross-tenant attempts return `approval.ErrCrossTenantAccess`. `IsSuperAdmin(p)` reports whether a principal carries the override role. The override role itself is exposed as `approval.SuperAdminRole` (`"approval:super_admin"`) for host wiring.

## Lifecycle Extension Points

| Extension point | Phase | Failure semantics |
| --- | --- | --- |
| `approval.InstanceLifecycleHook` (FX group `vef:approval:lifecycle_hooks`) | synchronous, inside the business transaction | returning an error rolls back the surrounding command |
| `approval.BusinessBindingHook` | mixed: `OnInstanceCreated` runs inside start_instance transaction; `WriteBackStatus` runs asynchronously from the binding listener after `InstanceCompletedEvent` | sync failure rolls back instance creation; async failure publishes `InstanceBindingFailedEvent` instead of rolling back |
| event subscriptions (`event.SubscribeTyped`) | asynchronous, after the transaction commits | the bus retries via the outbox relay; consumers must be idempotent |

### `InstanceLifecycleHook`

```go
type InstanceLifecycleHook interface {
    OnInstanceCreated(ctx context.Context, db orm.DB, instance *Instance) error
    OnInstanceCompleted(ctx context.Context, db orm.DB, instance *Instance, finalStatus InstanceStatus) error
}
```

Use lifecycle hooks for invariants that must hold inside the transaction (e.g. allocating a tightly-coupled business row). Use event subscriptions for everything else.

### `BusinessBindingHook`

```go
type BusinessBindingHook interface {
    OnInstanceCreated(ctx context.Context, db orm.DB, flow *Flow, instance *Instance) (businessRecordID string, err error)
    WriteBackStatus(ctx context.Context, db orm.DB, flow *Flow, instance *Instance, finalStatus InstanceStatus) error
}
```

Bridge between the approval engine and the host's business tables when `Flow.BindingMode == BindingBusiness`. Inject via `vef.SupplyBusinessBindingHook`. The `WriteBackStatus` async path runs from the binding listener — implementations must be idempotent because the outbox relay may retry.

## Business Identifier Validation

`Flow.BindingMode == BindingBusiness` flows carry SQL identifiers (`BusinessTable`, `BusinessPkField`, `BusinessTitleField`, `BusinessStatusField`) that the default binding hook interpolates directly into an `UPDATE` template. To prevent SQL injection, the framework whitelists identifiers against `^[A-Za-z_][A-Za-z0-9_]{0,62}$`:

```go
if err := approval.ValidateBusinessIdentifier(table); err != nil {
    return err
}
```

Empty / whitespace-only strings pass — the caller decides whether absence itself is an error. Anything outside the whitelist returns `approval.ErrInvalidBusinessIdentifier`. Admin-side Flow CRUD should bubble this up so operators see meaningful errors.

## Error Surface

The importable `approval` package exports four plain Go sentinels. They are
recognized with `errors.Is`, but they are not `result.Error` values and do not
carry an API code or HTTP status by themselves.

| Error | Source package | Meaning |
| --- | --- | --- |
| `approval.ErrCrossTenantAccess` | `approval` | non-super-admin caller attempted cross-tenant access |
| `approval.ErrInvalidBusinessIdentifier` | `approval` | business table / field identifier failed the SQL-identifier whitelist |
| `approval.ErrUnknownNodeKind` | `approval` | `NodeDefinition.ParseData` saw an unsupported `kind` |
| `approval.ErrNodeDataUnmarshal` | `approval` | `NodeDefinition.ParseData` could not decode node `data` |

Built-in approval resources return module-owned `result.Error` values through
the normal API envelope. Those values live under internal packages, so host
applications should treat the code/message pair below as the public wire
surface rather than importing the internal Go symbols.

| Code | Code constant | Error value | i18n message key | Notes |
| --- | --- | --- | --- | --- |
| `40001` | `ErrCodeFlowNotFound` | `ErrFlowNotFound` | `approval_flow_not_found` | flow lookup failed |
| `40002` | `ErrCodeFlowNotActive` | `ErrFlowNotActive` | `approval_flow_not_active` | flow is disabled |
| `40003` | `ErrCodeNoPublishedVersion` | `ErrNoPublishedVersion` | `approval_no_published_version` | flow has no published version |
| `40004` | `ErrCodeVersionNotDraft` | `ErrVersionNotDraft` | `approval_version_not_draft` | operation requires a draft version |
| `40005` | `ErrCodeInvalidFlowDesign` | `ErrInvalidFlowDesign` | `approval_invalid_flow_design` | graph or node design failed validation |
| `40006` | `ErrCodeFlowCodeExists` | `ErrFlowCodeExists` | `approval_flow_code_exists` | duplicate flow code |
| `40007` | `ErrCodeVersionNotFound` | `ErrVersionNotFound` | `approval_version_not_found` | flow version lookup failed |
| `40008` | `ErrCodeInvalidBusinessIdentifier` | `ErrInvalidBusinessIdentifier` | `approval_invalid_business_identifier` | business table / field identifier failed validation |
| `40101` | `ErrCodeInstanceNotFound` | `ErrInstanceNotFound` | `approval_instance_not_found` | instance lookup failed |
| `40102` | `ErrCodeInstanceCompleted` | `ErrInstanceCompleted` | `approval_instance_completed` | instance is already complete |
| `40103` | `ErrCodeNotAllowedInitiate` | `ErrNotAllowedInitiate` | `approval_not_allowed_initiate` | caller cannot initiate this flow |
| `40104` | `ErrCodeWithdrawNotAllowed` | `ErrWithdrawNotAllowed` | `approval_withdraw_not_allowed` | withdraw is not allowed in the current state |
| `40105` | `ErrCodeResubmitNotAllowed` | `ErrResubmitNotAllowed` | `approval_resubmit_not_allowed` | resubmit is not allowed in the current state |
| `40106` | `ErrCodeInvalidInstanceTransition` | `ErrInvalidInstanceTransition` | `approval_invalid_instance_transition` | instance state transition is invalid |
| `40201` | `ErrCodeTaskNotFound` | `ErrTaskNotFound` | `approval_task_not_found` | task lookup failed |
| `40202` | `ErrCodeTaskNotPending` | `ErrTaskNotPending` | `approval_task_not_pending` | task is not pending |
| `40203` | `ErrCodeNotAssignee` | `ErrNotAssignee` | `approval_not_assignee` | caller is not assigned to the task |
| `40204` | `ErrCodeInvalidTaskTransition` | `ErrInvalidTaskTransition` | `approval_invalid_task_transition` | task state transition is invalid |
| `40205` | `ErrCodeRollbackNotAllowed` | `ErrRollbackNotAllowed` | `approval_rollback_not_allowed` | rollback is disabled or not valid here |
| `40206` | `ErrCodeAddAssigneeNotAllowed` | `ErrAddAssigneeNotAllowed` | `approval_add_assignee_not_allowed` | dynamic assignee insertion is disabled |
| `40207` | `ErrCodeTransferNotAllowed` | `ErrTransferNotAllowed` | `approval_transfer_not_allowed` | transfer is disabled |
| `40208` | `ErrCodeOpinionRequired` | `ErrOpinionRequired` | `approval_opinion_required` | required opinion is blank |
| `40209` | `ErrCodeManualCcNotAllowed` | `ErrManualCcNotAllowed` | `approval_manual_cc_not_allowed` | manual CC is disabled |
| `40210` | `ErrCodeRemoveAssigneeNotAllowed` | `ErrRemoveAssigneeNotAllowed` | `approval_remove_assignee_not_allowed` | dynamic assignee removal is disabled |
| `40211` | `ErrCodeInvalidAddAssigneeType` | `ErrInvalidAddAssigneeType` | `approval_invalid_add_assignee_type` | `addType` is not one of `before`, `after`, `parallel` |
| `40212` | `ErrCodeNotApplicant` | `ErrNotApplicant` | `approval_not_applicant` | caller is not the applicant |
| `40213` | `ErrCodeInvalidRollbackTarget` | `ErrInvalidRollbackTarget` | `approval_invalid_rollback_target` | rollback target is not allowed |
| `40214` | `ErrCodeLastAssigneeRemoval` | `ErrLastAssigneeRemoval` | `approval_last_assignee_removal` | removal would leave no active assignee |
| `40215` | `ErrCodeInvalidTransferTarget` | `ErrInvalidTransferTarget` | `approval_invalid_transfer_target` | transfer or reassignment target is invalid |
| `40301` | `ErrCodeNoAssignee` | `ErrNoAssignee` | `approval_no_assignee` | no assignee could be resolved |
| `40302` | `ErrCodeAssigneeResolveFailed` | `ErrAssigneeResolveFailed` | `approval_assignee_resolve_failed` | assignee resolver failed |
| `40401` | `ErrCodeFormValidationFailed` | `ErrFormValidationFailed` | `approval_form_validation_failed` | general form validation failure |
| `40401` | `ErrCodeFormValidationFailed` | `ErrFormDataTooLarge` | `approval_form_data_too_large` | same code; JSON-encoded `formData` exceeded 64 KiB |
| `40401` | `ErrCodeFormValidationFailed` | dynamic form validation `result.Err` | `approval_form_field_not_defined`, `approval_form_field_required`, `approval_form_field_must_be_string`, `approval_form_field_must_be_number`, `approval_form_field_min_length`, `approval_form_field_max_length`, `approval_form_field_invalid_validation`, `approval_form_field_pattern_mismatch`, `approval_form_field_min_value`, `approval_form_field_max_value`, `approval_form_field_empty`, `approval_form_field_invalid_file_item`, `approval_form_field_must_be_file`, `approval_form_field_invalid_value` | field-level validation messages are constructed dynamically |
| `40402` | `ErrCodeFieldNotEditable` | `ErrFieldNotEditable` | `approval_field_not_editable` | submitted field is not editable for this task |
| `40501` | `ErrCodeDelegationNotFound` | `ErrDelegationNotFound` | `approval_delegation_not_found` | delegation lookup failed |
| `40502` | `ErrCodeDelegationConflict` | `ErrDelegationConflict` | `approval_delegation_conflict` | delegation window conflicts with an existing delegation |
| `40601` | `ErrCodeUrgeCooldown` | dynamic urge `result.Err` | `approval_urge_too_frequent` | no static sentinel; message is rendered with `minutes`; non-positive `urgeCooldownMinutes` defaults to 30 minutes |
| `40701` | `ErrCodeAccessDenied` | `ErrAccessDenied` | `approval_access_denied` | caller lacks approval-domain access |
| `40702` | `ErrCodeInstanceNotRunning` | `ErrInstanceNotRunning` | `approval_instance_not_running` | admin action requires a running instance |

Startup and tenant-resolution diagnostics such as
`ErrEventRouteNotTransactional`, `ErrEventRouteNotSubscribable`, and
`ErrTenantNotResolved` live under `internal/approval/...`; they are not
importable public Go API, but operators may see their wrapped messages when
event routing or tenant principal details are misconfigured.

## Supporting Public API Map

| Area | Public API |
| --- | --- |
| caller safety | `CallerContext`, `SystemCaller`, `IsSuperAdmin`, `SuperAdminRole`, `ErrCrossTenantAccess` |
| form data | `FormData`, `NewFormData`, `FormDefinition`, `FormFieldDefinition`, `FormSnapshot`, `ValidationRule`, `StorageMode`, `StorageJSON`, `FieldKind`, `FieldInput`, `FieldNumber`, `FieldDate`, `FieldTextarea`, `FieldSelect`, `FieldUpload`, `FieldOption` |
| flow models | `FlowCategory`, `Flow`, `FlowVersion`, `FlowNode`, `FlowEdge`, `FlowInitiator`, `FlowNodeAssignee`, `FlowNodeCC`, `VersionStatus`, `VersionDraft`, `VersionPublished`, `VersionArchived`, `ActionLog`, `OperatorInfo`, `UrgeRecord` |
| node design | `FlowDefinition`, `NodeDefinition`, `EdgeDefinition`, `Position`, `NodeData`, `BaseNodeData`, `StartNodeData`, `ApprovalNodeData`, `HandleNodeData`, `ConditionNodeData`, `CCNodeData`, `EndNodeData`, `ErrUnknownNodeKind`, `ErrNodeDataUnmarshal` |
| conditions | `ConditionKind`, `ConditionField`, `ConditionExpression`, `Condition`, `ConditionGroup`, `ConditionBranch`, `EvaluationContext`, `ConditionEvaluator` |
| initiators and assignees | `InitiatorKind`, `InitiatorUser`, `InitiatorRole`, `InitiatorDepartment`, `AssigneeKind`, `AssigneeDefinition`, `AssigneeService`, `ResolvedAssignee`, `UserInfo`, `UserInfoResolver`, `AddAssigneeType`, `AddAssigneeBefore`, `AddAssigneeAfter`, `AddAssigneeParallel` |
| CC | `CCKind`, `CCUser`, `CCRole`, `CCDepartment`, `CCFormField`, `CCTiming`, `CCTimingAlways`, `CCTimingOnApprove`, `CCTimingOnReject`, `CCDefinition`, `CCRecord` |
| node behavior | `ApprovalMethod`, `TaskNodeData`, `ExecutionType`, `ExecutionManual`, `ExecutionAuto`, `ExecutionAutoPass`, `ExecutionAutoReject`, `ConsecutiveApproverAction`, `ConsecutiveApproverNone`, `ConsecutiveApproverAutoPass`, `SameApplicantAction`, `SameApplicantSelfApprove`, `SameApplicantAutoPass`, `SameApplicantTransferSuperior`, `Permission`, `PermissionVisible`, `PermissionEditable`, `PermissionRequired`, `PermissionHidden` |
| rollback and timeouts | `RollbackType`, `RollbackNone`, `RollbackPrevious`, `RollbackStart`, `RollbackAny`, `RollbackSpecified`, `RollbackDataStrategy`, `RollbackDataClear`, `RollbackDataKeep`, `EmptyAssigneeAction`, `EmptyAssigneeAutoPass`, `EmptyAssigneeTransferAdmin`, `EmptyAssigneeTransferSuperior`, `EmptyAssigneeTransferApplicant`, `EmptyAssigneeTransferSpecified`, `TimeoutAction`, `TimeoutActionNone`, `TimeoutActionAutoPass`, `TimeoutActionAutoReject`, `TimeoutActionNotify`, `TimeoutActionTransferAdmin` |
| action and status enums | `ActionType`, `InstanceStatus`, `TaskStatus`, `NodeKind`, `StorageMode`, `VersionStatus` |
| pass rules | `PassRule`, `PassRuleContext`, `PassRuleStrategy`, `PassRuleResult`, `PassRulePending`, `PassRulePassed`, `PassRuleRejected` |
| events | all `New...Event` constructors, `DomainEvent`, `PayloadOccurredAt`, and the `EventType...` constants |
| extension interfaces | `InstanceLifecycleHook`, `BusinessBindingHook`, `InstanceNoGenerator`, `ConditionEvaluator`, `PrincipalTenantResolver`, `PrincipalDepartmentResolver` |
| admin DTOs | package `approval/admin`: `Instance`, `InstanceDetail`, `InstanceDetailInfo`, `Task`, `TaskDetailInfo`, `ActionLog`, `FlowNodeInfo`, `Metrics` |
| user DTOs | package `approval/my`: `PendingTask`, `CompletedTask`, `CCRecord`, `InitiatedInstance`, `AvailableFlow`, `InstanceDetail`, `InstanceInfo`, `TaskInfo`, `ActionLogInfo`, `FlowNodeInfo`, `PendingCounts` |

## Delegation

Users can delegate their approval authority to others:

```go
type Delegation struct {
    DelegatorID    string         // Who delegates
    DelegateeID    string         // Who receives delegation
    FlowCategoryID *string        // Optional: limit to category
    FlowID         *string        // Optional: limit to specific flow
    StartTime      timex.DateTime // Delegation start
    EndTime        timex.DateTime // Delegation end
    IsActive       bool
}
```

## Flow Definition (React Flow Compatible)

The `FlowDefinition` struct is compatible with React Flow's JSON format:

```go
type FlowDefinition struct {
    Nodes []NodeDefinition `json:"nodes"`
    Edges []EdgeDefinition `json:"edges"`
}
```

Each `NodeDefinition` contains a `Kind` and typed `Data` that is parsed into the appropriate struct (`StartNodeData`, `ApprovalNodeData`, `HandleNodeData`, `ConditionNodeData`, `CCNodeData`, `EndNodeData`).

### Flow JSON Wire Shape

`deploy` treats the flow definition as a full snapshot. `NodeDefinition.ParseData`
chooses the typed `data` struct from `kind`; an unknown kind returns
`ErrUnknownNodeKind`, and malformed node `data` is wrapped with
`ErrNodeDataUnmarshal`.

| Type | JSON fields |
| --- | --- |
| `FlowDefinition` | `nodes`, `edges` |
| `NodeDefinition` | `id`, `kind`, `position`, `data`; `position` contains `x` and `y` |
| `EdgeDefinition` | `id`, `source`, `target`, `sourceHandle`, `data` |

`sourceHandle` is required only for edges leaving a condition node, where it
must match a branch `id`. Non-condition outgoing edges must omit
`sourceHandle`. `EdgeDefinition.data` is designer metadata stored in the
version `flowSchema`; runtime routing is driven by `source`, `target`, and
`sourceHandle`.

Node `data` fields are:

| Node data type | JSON fields |
| --- | --- |
| `BaseNodeData` | `name`, `description`; embedded by every node data type |
| `StartNodeData` | base fields only |
| `EndNodeData` | base fields only |
| `TaskNodeData` | `assignees`, `executionType`, `emptyAssigneeAction`, `fallbackUserIds`, `adminUserIds`, `isTransferAllowed`, `isOpinionRequired`, `timeoutHours`, `timeoutAction`, `timeoutNotifyBeforeHours`, `urgeCooldownMinutes`, `ccs`, `fieldPermissions` |
| `ApprovalNodeData` | base fields + `TaskNodeData` fields + `approvalMethod`, `passRule`, `passRatio`, `sameApplicantAction`, `consecutiveApproverAction`, `rollbackType`, `rollbackDataStrategy`, `rollbackTargetKeys`, `isRollbackAllowed`, `isAddAssigneeAllowed`, `addAssigneeTypes`, `isRemoveAssigneeAllowed`, `isManualCcAllowed` |
| `HandleNodeData` | base fields + `TaskNodeData` fields; if unset, deploy defaults `approvalMethod` to `sequential` and `passRule` to `any` |
| `CCNodeData` | base fields + `ccs`, `isReadConfirmRequired`, `fieldPermissions` |
| `ConditionNodeData` | base fields + `branches` |

`assignees` entries use `kind`, `ids`, `formField`, and `sortOrder`. `ccs`
entries use `kind`, `ids`, `formField`, and `timing`. During deployment,
these embedded arrays are materialized into `FlowNodeAssignee` and
`FlowNodeCC` records in addition to the `FlowNode` row.

Condition branches use `id`, `label`, `conditionGroups`, `isDefault`, and
`priority`. Each `conditionGroups` entry contains `conditions`; each condition
uses `kind`, `subject`, `operator`, `value`, and `expression`.

`timeoutHours` and `timeoutNotifyBeforeHours` are in hours.
`urgeCooldownMinutes` is in minutes; values less than or equal to 0 use the
runtime default of 30 minutes. `rollbackTargetKeys` is checked when
`rollbackType` is `specified`; it contains node keys, not database node IDs.
During task processing, submitted `formData` is merged only for fields whose
`fieldPermissions` entry is `editable` or `required`; fields marked `visible`,
`hidden`, or omitted from the map are ignored for that task update.

### Form JSON Wire Shape

`FormDefinition` uses a `fields` array. Each `FormFieldDefinition` entry uses
`key`, `kind`, `label`, `placeholder`, `defaultValue`, `isRequired`,
`options`, `validation`, `props`, and `sortOrder`. Each option uses `label`
and `value`.

`validation` supports `minLength`, `maxLength`, `min`, `max`, `pattern`, and
`message`. Submitted `formData` is capped at 64 KiB after JSON encoding, even
when the flow has no form schema. When a schema exists, extra form keys are
rejected; required fields reject absent, `null`, blank-string, and empty-array
values. `input`, `textarea`, and `date` fields must be strings and may use
`minLength`, `maxLength`, and `pattern`. `number` fields accept numeric JSON
values and may use `min` and `max`. `select` fields validate scalar or array
values against `options` when options are present. `upload` fields accept a
non-blank string, a non-empty `[]string`, or a non-empty array of non-blank
strings. `validation.message` is used as the custom error message for
`pattern` mismatches; other validation failures use the module i18n messages.

## Instance Number Generation

Implement the `InstanceNoGenerator` interface to customize instance numbering:

```go
type InstanceNoGenerator interface {
    Generate(ctx context.Context, flowCode string) (string, error)
}
```
