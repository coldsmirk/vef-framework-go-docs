---
sidebar_position: 2
---

# RPC Resources

When enabled, the module registers these RPC resources. They are not public
operations: callers must be authenticated, and permissions below are enforced
when a `RequiredPermission` is shown.

RPC calls use the standard envelope documented in [API](../building-apis/api.md): `resource`,
`action`, `version`, `params`, and `meta`. Fields listed below as `params.*`
are decoded from `params`; fields listed as `meta.*` are decoded from `meta`.
The generated [Runtime API Index](../reference/runtime-api-index.md) contains the
exhaustive JSON field ledger for every request and response DTO.

## `approval/category`

| Action | Permission | Params | Notes |
| --- | --- | --- | --- |
| `find_tree` | `approval.category.query` | `CategorySearch` | Tenant-scoped tree query |
| `find_tree_options` | `approval.category.query` | `CategorySearch` + `DataOptionConfig` | Tenant-scoped tree options |
| `create` | `approval.category.create` | `CategoryParams` | Non-super-admin tenant is stamped from caller |
| `update` | `approval.category.update` | `CategoryParams` | Non-super-admin can only mutate own tenant |
| `delete` | `approval.category.delete` | Primary-key params | Non-super-admin can only delete own tenant |

| Action | Request fields |
| --- | --- |
| `find_tree` | `meta.name`, `meta.isActive`, `meta.sort` |
| `find_tree_options` | `meta.name`, `meta.isActive`, `meta.sort`, plus option mapping metadata: `meta.labelColumn`, `meta.valueColumn`, `meta.descriptionColumn`, `meta.metaColumns` |
| `create` | `params.id`, `params.tenantId` required, `params.code` required, `params.name` required, `params.icon`, `params.parentId`, `params.sortOrder`, `params.isActive`, `params.remark` |
| `update` | `params.id`, `params.tenantId` required, `params.code` required, `params.name` required, `params.icon`, `params.parentId`, `params.sortOrder`, `params.isActive`, `params.remark` |
| `delete` | `params.id` required |

## `approval/delegation`

| Action | Permission | Params | Notes |
| --- | --- | --- | --- |
| `find_page` | `approval.delegation.query` | `DelegationSearch` + pageable meta | Non-super-admin sees own delegations |
| `create` | `approval.delegation.create` | `DelegationParams` | Non-super-admin delegator is stamped from caller |
| `update` | `approval.delegation.update` | `DelegationParams` | Non-super-admin cannot reassign ownership |
| `delete` | `approval.delegation.delete` | Primary-key params | Owner-scoped for non-super-admin |

| Action | Request fields |
| --- | --- |
| `find_page` | `meta.delegatorId`, `meta.delegateeId`, `meta.isActive`, `meta.sort`, `meta.page`, `meta.size` |
| `create` | `params.id`, `params.delegatorId` required, `params.delegateeId` required, `params.flowCategoryId`, `params.flowId`, `params.startTime`, `params.endTime`, `params.isActive`, `params.reason` |
| `update` | `params.id`, `params.delegatorId` required, `params.delegateeId` required, `params.flowCategoryId`, `params.flowId`, `params.startTime`, `params.endTime`, `params.isActive`, `params.reason` |
| `delete` | `params.id` required |

## `approval/flow`

| Action | Permission | Params | Notes |
| --- | --- | --- | --- |
| `create` | `approval.flow.create` | `CreateFlowParams` | Audited |
| `deploy` | `approval.flow.deploy` | `DeployFlowParams` | Audited |
| `publish_version` | `approval.flow.publish` | `PublishVersionParams` | Audited |
| `update` | `approval.flow.update` | `UpdateParams` | Audited |
| `toggle_active` | `approval.flow.update` | `ToggleActiveParams` | Audited |
| `get_graph` | `approval.flow.query` | `GetGraphParams` | Reads published graph |
| `find_flows` | `approval.flow.query` | `FindFlowsParams` | Paged query |
| `find_initiators` | `approval.flow.query` | `FindInitiatorsParams` | Lists initiator configuration |
| `find_versions` | `approval.flow.query` | `FindVersionsParams` | Lists versions for one flow |

| Action | Request fields |
| --- | --- |
| `create` | `params.tenantId` required, `params.code` required, `params.name` required, `params.categoryId` required, `params.bindingMode` required, `params.icon`, `params.description`, `params.businessBinding`, `params.adminUserIds`, `params.isAllInitiationAllowed`, `params.instanceTitleTemplate`, `params.initiators` |
| `deploy` | `params.flowId` required, `params.description`, `params.flowDefinition` required, `params.formSchema`, `params.storageMode` |
| `publish_version` | `params.versionId` required |
| `update` | `params.flowId` required, `params.name` required, `params.bindingMode` required, `params.instanceTitleTemplate` required, `params.icon`, `params.description`, `params.businessBinding`, `params.adminUserIds`, `params.isAllInitiationAllowed`, `params.initiators` |
| `toggle_active` | `params.flowId` required, `params.isActive` |
| `get_graph` | `params.flowId` required, `params.tenantId` |
| `find_flows` | `params.tenantId`, `params.categoryId`, `params.keyword`, `params.isActive`, `params.page`, `params.pageSize` |
| `find_initiators` | `params.flowId` required, `params.tenantId` |
| `find_versions` | `params.flowId` required, `params.tenantId` |

`params.initiators` entries use `kind` (`user`, `role`, or `department`) and
`ids`. `params.formSchema` is the optional host-owned form-designer document,
passed through opaque (see
[Form Schema and Derived Fields](./flow-design.md#form-schema-and-derived-fields));
flows without forms omit it. For business binding, `params.businessBinding` is
an `approval.BusinessBindingConfig` object (`tableName`, `keyColumns`,
`statusColumn`, `instanceIdColumn` required, optional `startedAtColumn` /
`finishedAtColumn` / `statusMapping`); all identifiers are validated by
`ValidateBusinessIdentifier`, the key must match a non-null primary or unique
key on the live table, and binding settings stay frozen while instances are
running (`ErrFlowBindingLocked`).

## `approval/instance`

| Action | Permission | Params | Notes |
| --- | --- | --- | --- |
| `start` | `approval.instance.start` | `StartParams` | Audited |
| `process_task` | `approval.task.process` | `ProcessTaskParams` | Audited; `action` must be `approve`, `reject`, `transfer`, `rollback`, or `handle` |
| `withdraw` | `approval.instance.withdraw` | `WithdrawParams` | Audited |
| `resubmit` | `approval.instance.resubmit` | `ResubmitParams` | Audited; accepts returned or withdrawn instances |
| `add_cc` | `approval.instance.cc` | `AddCCParams` | Audited |
| `mark_cc_read` | `approval.instance.cc` | `MarkCCReadParams` | Read receipt, not audited |
| `add_assignee` | `approval.task.add_assignee` | `AddAssigneeParams` | Audited; `addType` is `before`, `after`, or `parallel` |
| `remove_assignee` | `approval.task.remove_assignee` | `RemoveAssigneeParams` | Audited |
| `urge_task` | `approval.task.urge` | `UrgeTaskParams` | Extra rate limit: max `10` per `1m` |

| Action | Request fields |
| --- | --- |
| `start` | `params.tenantId` required, `params.flowCode` required, `params.businessRef` max 512 chars, `params.formData` |
| `process_task` | `params.taskId` required, `params.action` required (`approve`, `reject`, `transfer`, `rollback`, or `handle`), `params.opinion` max 2000 chars, `params.formData`, `params.attachments` max 20 entries, each max 512 chars, `params.transferToId`, `params.targetNodeId` |
| `withdraw` | `params.instanceId` required, `params.reason` max 2000 chars |
| `resubmit` | `params.instanceId` required, `params.formData` |
| `add_cc` | `params.instanceId` required, `params.ccUserIds` required, 1-50 IDs |
| `mark_cc_read` | `params.instanceId` required |
| `add_assignee` | `params.taskId` required, `params.userIds` required, 1-50 IDs, `params.addType` required (`before`, `after`, or `parallel`) |
| `remove_assignee` | `params.taskId` required |
| `urge_task` | `params.taskId` required, `params.message` max 500 chars |

## `approval/my`

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

The `my.InstanceDetail` JSON payload includes `instance`, `formSchema`,
`timeline`, `flowGraph`, `availableActions`, and `fieldPermissions`.
`instance.formData` carries the submitted form data, `instance.businessRef` is
the opaque business reference when the flow is business-bound, and
`formSchema` is the version-pinned host form-designer document returned
verbatim — the framework stores it as semantically equal JSON and never
interprets it. `fieldPermissions` (v0.38) is the viewer-scoped field
interactivity projection, materialized for every top-level form field
(`visible` / `editable` / `hidden` / `required`); the client applies it
verbatim, and `instance.formData` is already stripped of the fields the viewer
may not see (see [Node Field Permissions](./flow-design.md#node-field-permissions)).

`availableActions` is a query-layer UI hint. For the applicant it includes
`withdraw` when the instance can transition to `withdrawn`, and `resubmit` when
the instance is returned or withdrawn. For pending tasks it includes `handle` for handle
nodes, otherwise `approve`, then `reject`, plus `transfer`, `rollback`,
`add_assignee`, or `add_cc` when the current node allows them. If the instance
has any pending task, it also includes `urge`. Command handlers still perform
their own validation.

## `approval/admin`

| Action | Permission | Params | Notes |
| --- | --- | --- | --- |
| `find_instances` | `approval.instance.query` | `AdminFindInstancesParams` | Tenant-filtered for non-super-admin |
| `find_tasks` | `approval.task.query` | `AdminFindTasksParams` | Tenant-filtered for non-super-admin |
| `get_instance_detail` | `approval.instance.detail` | `AdminGetInstanceDetailParams` | Full admin detail |
| `find_action_logs` | `approval.action_log.query` | `AdminFindActionLogsParams` | Requires `instanceId` |
| `get_metrics` | `approval.metrics.query` | `AdminGetMetricsParams` | Aggregated metrics |
| `find_business_projections` | `approval.binding.query` | `AdminFindBusinessProjectionsParams` | Durable binding convergence state (v0.38) |
| `terminate_instance` | `approval.instance.terminate` | `AdminTerminateInstanceParams` | Audited |
| `reassign_task` | `approval.task.reassign` | `AdminReassignTaskParams` | Audited |
| `retry_business_projection` | `approval.binding.retry` | `AdminRetryBusinessProjectionParams` | Audited; immediately retries one eventual projection (v0.38) |

| Action | Request fields |
| --- | --- |
| `find_instances` | `params.tenantId`, `params.applicantId`, `params.status`, `params.flowId`, `params.keyword`, `params.page`, `params.pageSize` |
| `find_tasks` | `params.tenantId`, `params.assigneeId`, `params.instanceId`, `params.status`, `params.page`, `params.pageSize` |
| `get_instance_detail` | `params.instanceId` required |
| `find_action_logs` | `params.instanceId` required, `params.tenantId`, `params.page`, `params.pageSize` |
| `get_metrics` | `params.tenantId` |
| `find_business_projections` | `params.tenantId`, `params.status` (`pending`, `processing`, `applied`, `failed`), `params.page`, `params.pageSize` |
| `terminate_instance` | `params.instanceId` required, `params.reason` max 2000 chars |
| `reassign_task` | `params.taskId` required, `params.newAssigneeId` required, `params.reason` max 2000 chars |
| `retry_business_projection` | `params.projectionId` required |

For admin list and metrics queries, non-super-admin callers ignore a submitted
`tenantId` override and are filtered to their own tenant. Super-admin callers
may pass `tenantId` to filter one tenant or omit it for cross-tenant visibility.

Admin list/detail DTOs keep tenant scope on the wire as `tenantId`. The raw
audit trail is still available through `find_action_logs`; detail responses
carry timeline and flow-graph projections for rendering.

## Response DTO Fields

Admin responses use the DTOs from `approval/admin`:

| DTO | JSON fields |
| --- | --- |
| `admin.Instance` | `instanceId`, `instanceNo`, `title`, `tenantId`, `flowId`, `flowName`, `applicant` (`UserInfo`), `status`, `currentNodeName`, `createdAt`, `finishedAt` |
| `admin.Task` | `taskId`, `instanceId`, `instanceTitle`, `flowName`, `nodeName`, `assignee` (`UserInfo`), `status`, `createdAt`, `deadline`, `finishedAt` |
| `admin.InstanceDetail` | `instance`, `formSchema` (host designer document, verbatim), `timeline`, `flowGraph` |
| `admin.InstanceDetailInfo` | `instanceId`, `instanceNo`, `title`, `tenantId`, `flowId`, `flowName`, `flowVersionId`, `applicant`, `status`, `currentNodeId`, `currentNodeName`, `businessRef`, `formData`, `createdAt`, `finishedAt` |
| `admin.ActionLog` | `logId`, `action`, `nodeId`, `taskId`, `operator`, `transferTo`, `rollbackToNodeId`, `addedAssignees`, `removedAssignees`, `ccUsers`, `opinion`, `attachments`, `createdAt` |
| `admin.Metrics` | `tenantId`, `capturedAt`, `instanceCounts`, `taskCounts`, `timeoutTaskCount`, `avgCompletionSeconds`, `pendingBindingFailures`, `businessProjectionCounts`, `pendingBusinessProjections` |
| `admin.BusinessProjection` | `projectionId`, `tenantId`, `flowId`, `flowVersionId`, `ownerInstanceId`, `appliedOwnerInstanceId`, `businessTable`, `recordKey`, `consistency`, `desiredStatus`, `desiredStartedAt`, `desiredFinishedAt`, `desiredRevision`, `appliedRevision`, `status`, `attemptCount`, `nextAttemptAt`, `leaseUntil`, `lastError`, `appliedAt`, `updatedAt` |

Self-service responses use the DTOs from `approval/my`:

| DTO | JSON fields |
| --- | --- |
| `my.AvailableFlow` | `flowId`, `flowCode`, `flowName`, `flowIcon`, `description`, `categoryId`, `categoryName` |
| `my.InitiatedInstance` | `instanceId`, `instanceNo`, `title`, `flowName`, `flowIcon`, `status`, `currentNodeName`, `createdAt`, `finishedAt` |
| `my.PendingTask` | `taskId`, `instanceId`, `instanceTitle`, `instanceNo`, `flowName`, `flowIcon`, `applicant` (`UserInfo`), `nodeName`, `createdAt`, `deadline`, `isTimeout` |
| `my.CompletedTask` | `taskId`, `instanceId`, `instanceTitle`, `instanceNo`, `flowName`, `flowIcon`, `applicant` (`UserInfo`), `nodeName`, `status`, `finishedAt` |
| `my.CCRecord` | `ccRecordId`, `instanceId`, `instanceTitle`, `instanceNo`, `flowName`, `flowIcon`, `applicant` (`UserInfo`), `nodeName`, `isRead`, `createdAt` |
| `my.PendingCounts` | `pendingTaskCount`, `unreadCcCount` |
| `my.InstanceDetail` | `instance`, `formSchema` (host designer document, verbatim), `timeline`, `flowGraph`, `availableActions`, `fieldPermissions` |
| `my.InstanceInfo` | `instanceId`, `instanceNo`, `title`, `flowName`, `flowIcon`, `applicant`, `status`, `currentNodeId`, `currentNodeName`, `businessRef`, `formData`, `createdAt`, `finishedAt` |

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
| `40009` | `ErrCodeInvalidTitleTemplate` | `ErrInvalidTitleTemplate` | `approval_invalid_title_template` | instance title template failed parsing |
| `40010` | `ErrCodeInvalidFormDesign` | `ErrInvalidFormDesign` | `approval_invalid_form_design` | form schema failed design-time validation |
| `40011` | `ErrCodeBindingIncomplete` | `ErrBindingIncomplete` | `approval_binding_incomplete` | business binding is missing required table / key / status / instance-id fields |
| `40012` | `ErrCodeInvalidBindingMode` | `ErrInvalidBindingMode` | `approval_invalid_binding_mode` | flow binding mode is out of enum |
| `40013` | `ErrCodeInvalidInitiatorKind` | `ErrInvalidInitiatorKind` | `approval_invalid_initiator_kind` | flow initiator kind is out of enum |
| `40014` | `ErrCodeInvalidStorageMode` | `ErrInvalidStorageMode` | `approval_invalid_storage_mode` | deploy requested a storage mode other than `json` or `table` |
| `40015` | `ErrCodeFlowBindingLocked` | `ErrFlowBindingLocked` | `approval_flow_binding_locked` | retained as a stable error surface; version-pinned binding snapshots mean current flow commands no longer return it |
| `40016` | `ErrCodeBindingColumnsConflict` | `ErrBindingColumnsConflict` | `approval_binding_columns_conflict` | two business-binding fields name the same column |
| `40017` | `ErrCodeBindingUnexpected` | `ErrBindingUnexpected` | `approval_binding_unexpected` | business binding supplied on a standalone flow |
| `40018` | `ErrCodeBindingSchemaInvalid` | `ErrBindingSchemaInvalid` | `approval_binding_schema_invalid` | configured binding table or columns do not exist in the primary database |
| `40019` | `ErrCodeBindingKeyNotUnique` | `ErrBindingKeyNotUnique` | `approval_binding_key_not_unique` | key columns are not backed by one complete, non-null primary or unique key |
| `40020` | `ErrCodeBindingStatusMappingInvalid` | `ErrBindingStatusMappingInvalid` | `approval_binding_status_mapping_invalid` | status mapping names an unknown status or maps to a blank value |
| `40101` | `ErrCodeInstanceNotFound` | `ErrInstanceNotFound` | `approval_instance_not_found` | instance lookup failed |
| `40102` | `ErrCodeInstanceCompleted` | `ErrInstanceCompleted` | `approval_instance_completed` | instance is already complete |
| `40103` | `ErrCodeNotAllowedInitiate` | `ErrNotAllowedInitiate` | `approval_not_allowed_initiate` | caller cannot initiate this flow |
| `40104` | `ErrCodeWithdrawNotAllowed` | `ErrWithdrawNotAllowed` | `approval_withdraw_not_allowed` | withdraw is not allowed in the current state |
| `40105` | `ErrCodeResubmitNotAllowed` | `ErrResubmitNotAllowed` | `approval_resubmit_not_allowed` | resubmit is not allowed in the current state |
| `40106` | `ErrCodeInvalidInstanceTransition` | `ErrInvalidInstanceTransition` | `approval_invalid_instance_transition` | instance state transition is invalid |
| `40107` | `ErrCodeBusinessRefRequired` | `ErrBusinessRefRequired` | `approval_business_ref_required` | business-bound flow started without a business reference |
| `40108` | `ErrCodeBindingTargetBusy` | `ErrBindingTargetBusy` | `approval_binding_target_busy` | the business record is already claimed by a non-final approval instance |
| `40109` | `ErrCodeInvalidBusinessRef` | `ErrInvalidBusinessRef` | `approval_invalid_business_ref` | the business reference could not be resolved into the configured record key |
| `40110` | `ErrCodeBindingProjectionNotFound` | `ErrBindingProjectionNotFound` | `approval_binding_projection_not_found` | projection lookup failed (admin retry) |
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
| `40216` | `ErrCodeNoUsersSpecified` | `ErrNoUsersSpecified` | `approval_no_users_specified` | user-list operation received no target users |
| `40301` | `ErrCodeNoAssignee` | `ErrNoAssignee` | `approval_no_assignee` | no assignee could be resolved |
| `40302` | `ErrCodeAssigneeResolveFailed` | `ErrAssigneeResolveFailed` | `approval_assignee_resolve_failed` | assignee resolver failed |
| `40401` | `ErrCodeFormValidationFailed` | `ErrFormValidationFailed` | `approval_form_validation_failed` | general form validation failure |
| `40401` | `ErrCodeFormValidationFailed` | `ErrFormDataTooLarge` | `approval_form_data_too_large` | same code; JSON-encoded `formData` exceeded 64 KiB |
| `40401` | `ErrCodeFormValidationFailed` | dynamic form validation `result.Err` | `approval_form_field_not_defined`, `approval_form_field_required`, `approval_form_field_must_be_string`, `approval_form_field_must_be_number`, `approval_form_field_must_be_integer`, `approval_form_field_min_length`, `approval_form_field_max_length`, `approval_form_field_invalid_validation`, `approval_form_field_pattern_mismatch`, `approval_form_field_min_value`, `approval_form_field_max_value`, `approval_form_field_empty`, `approval_form_field_invalid_file_item`, `approval_form_field_must_be_file`, `approval_form_field_invalid_value`, `approval_form_field_must_be_row_list`, `approval_form_field_must_be_row_object`, `approval_form_field_min_rows`, `approval_form_field_max_rows`, `approval_form_field_table_cell` | field-level validation messages are constructed dynamically |
| `40601` | `ErrCodeUrgeCooldown` | dynamic urge `result.Err` | `approval_urge_too_frequent` | no static sentinel; message is rendered with `minutes`; non-positive `urgeCooldownMinutes` defaults to 30 minutes |
| `40701` | `ErrCodeAccessDenied` | `ErrAccessDenied` | `approval_access_denied` | caller lacks approval-domain access |
| `40702` | `ErrCodeTerminateNotAllowed` | `ErrTerminateNotAllowed` | `approval_terminate_not_allowed` | terminate is not allowed from the current instance state |

Startup and tenant-resolution diagnostics such as
`ErrEventRouteNotTransactional`, `ErrEventRouteNotSubscribable`, and
`ErrTenantNotResolved` live under `internal/approval/...`; they are not
importable public Go API, but operators may see their wrapped messages when
event routing or tenant principal details are misconfigured.

---

Next: [Flow Design](./flow-design.md) for the designer wire shapes behind `deploy`, or [Instance Runtime](./runtime.md) for lifecycle semantics behind the instance actions.
