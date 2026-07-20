---
sidebar_position: 2
---

# RPC Resources

When the approval module is enabled, the framework registers the six RPC
resources below. All of them are mounted under `/api`, using the standard
envelope (`resource`, `action`, `version`, `params`, `meta`) documented in
[API](../building-apis/api.md). None of the operations are public: callers must
be authenticated, and permissions are enforced wherever a `RequiredPermission`
is listed. The generated [Runtime API Index](../reference/runtime-api-index.md)
contains the exhaustive JSON field ledger for every request and response DTO.

Conventions used on this page:

- Command-style operations (`approval/flow`, `approval/instance`,
  `approval/my`, `approval/admin`) declare params structs embedding `api.P`,
  so their fields decode from the request's `params` object.
- CRUD read operations (`approval/category`, `approval/delegation`) declare
  search structs embedding `crud.Sortable` — a `meta` struct — so their filter
  fields decode from the request's `meta` object, next to `meta.page` /
  `meta.size` (`page.Pageable`) and `meta.sort`.
- Paged responses use `page.Page[T]`: `page`, `size`, `total`, `items`.
- Persisted models carry the standard audited-model columns (`id`,
  `createdAt`, `createdBy`, `updatedAt`, `updatedBy`) in responses; they are
  omitted from the field tables below.
- Enum vocabularies (`InstanceStatus`, `TaskStatus`, node semantics) are
  defined in [Instance Runtime](./runtime.md) and
  [Flow Design](./flow-design.md).

## `approval/category`

Flow category management (`apv_flow_category`). Reads are tenant-scoped:
super-admin callers see all tenants, everyone else is confined to their own
tenant and fails closed without one.

| Action | Permission | Input | Output |
| --- | --- | --- | --- |
| `find_tree` | `approval.category.query` | `CategorySearch` (meta) | nested `FlowCategory[]` (children populated) |
| `create` | `approval.category.create` | `CategoryParams` | created `FlowCategory` |
| `update` | `approval.category.update` | `CategoryParams` | updated `FlowCategory` |
| `delete` | `approval.category.delete` | primary-key params (`params.id`) | success |

There is no `find_tree_options` operation; build option lists from
`find_tree` instead.

`CategorySearch` (query filters, decoded from `meta`):

| Field | Type | Match | Description |
| --- | --- | --- | --- |
| `name` | `string` | contains | filter by category name fragment |
| `isActive` | `bool` | equals | filter by active flag; omit to match both |
| `sort` | `OrderSpec[]` | — | sort specifications (`crud.Sortable`) |

`CategoryParams` (create/update, decoded from `params`):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `id` | `string` | update only | primary key of the row to update |
| `tenantId` | `string` | Yes | owning tenant. On create, non-super-admin callers have this stamped from their own tenant (the submitted value is ignored); on update/delete the caller must be authorized for the row's tenant |
| `code` | `string` | Yes | category business code |
| `name` | `string` | Yes | display name |
| `icon` | `string` | No | display icon identifier |
| `parentId` | `string` | No | parent category id; `null` makes it a root |
| `sortOrder` | `int` | No | ordering weight among siblings |
| `isActive` | `bool` | No | inactive categories stay queryable but hosts typically hide them from pickers |
| `remark` | `string` | No | free-text remark |

`FlowCategory` (response model):

| Field | Type | Description |
| --- | --- | --- |
| `tenantId` | `string` | owning tenant |
| `code` | `string` | category business code |
| `name` | `string` | display name |
| `icon` | `string` \| `null` | display icon identifier |
| `parentId` | `string` \| `null` | parent category id |
| `sortOrder` | `int` | ordering weight |
| `isActive` | `bool` | active flag |
| `remark` | `string` \| `null` | free-text remark |
| `children` | `FlowCategory[]` | child categories; populated by `find_tree`, absent elsewhere |

## `approval/delegation`

Approval delegation management (`apv_delegation`). Ownership is enforced:
non-super-admin callers only see, create, update, and delete delegations where
they are the delegator — on create the `delegatorId` is stamped from the
caller, and on update the original delegator is pinned so the record cannot be
reassigned to another user.

| Action | Permission | Input | Output |
| --- | --- | --- | --- |
| `find_page` | `approval.delegation.query` | `DelegationSearch` + pageable meta | `page.Page[Delegation]` |
| `create` | `approval.delegation.create` | `DelegationParams` | created `Delegation` |
| `update` | `approval.delegation.update` | `DelegationParams` | updated `Delegation` |
| `delete` | `approval.delegation.delete` | primary-key params (`params.id`) | success |

`DelegationSearch` (query filters, decoded from `meta`):

| Field | Type | Match | Description |
| --- | --- | --- | --- |
| `delegatorId` | `string` | equals | filter by delegator (super-admin only — others are always scoped to themselves) |
| `delegateeId` | `string` | equals | filter by delegatee |
| `isActive` | `bool` | equals | filter by active flag |
| `sort` | `OrderSpec[]` | — | sort specifications |

`DelegationParams` (create/update, decoded from `params`):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `id` | `string` | update only | primary key of the row to update |
| `delegatorId` | `string` | Yes | user delegating their tasks. Non-super-admin callers have this stamped from the principal on create and pinned to the original value on update |
| `delegateeId` | `string` | Yes | user receiving the delegated tasks |
| `flowCategoryId` | `string` | No | restrict the delegation to one flow category; `null` covers all categories |
| `flowId` | `string` | No | restrict the delegation to one flow; `null` covers all flows |
| `startTime` | `DateTime` | Yes | delegation window start |
| `endTime` | `DateTime` | Yes | delegation window end |
| `isActive` | `bool` | No | inactive delegations are ignored by assignee resolution |
| `reason` | `string` | No | free-text reason shown in delegated tasks |

`Delegation` (response model): same business fields as the params —
`delegatorId`, `delegateeId`, `flowCategoryId`, `flowId`, `startTime`,
`endTime`, `isActive`, `reason` — plus the audited-model columns. Tasks that
arrive via delegation carry the delegator as a separate person snapshot (see
`NodeParticipant.delegator` below).

## `approval/flow`

Flow definition management: the mutable flow row, its immutable deployed
versions, and the designer-facing graph reads.

| Action | Permission | Input | Output | Audit |
| --- | --- | --- | --- | --- |
| `create` | `approval.flow.create` | `CreateFlowParams` | created `Flow` | Yes |
| `deploy` | `approval.flow.deploy` | `DeployFlowParams` | created `FlowVersion` (draft) | Yes |
| `publish_version` | `approval.flow.publish` | `PublishVersionParams` | success | Yes |
| `update` | `approval.flow.update` | `UpdateFlowParams` | updated `Flow` | Yes |
| `toggle_active` | `approval.flow.update` | `ToggleActiveParams` | success | Yes |
| `get_graph` | `approval.flow.query` | `GetGraphParams` | `FlowGraph` | — |
| `find_flows` | `approval.flow.query` | `FindFlowsParams` | `page.Page[Flow]` | — |
| `find_versions` | `approval.flow.query` | `FindVersionsParams` | `FlowVersionSummary[]` | — |
| `find_initiators` | `approval.flow.query` | `FindInitiatorsParams` | `FlowInitiator[]` | — |

`CreateFlowParams` (`create`):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `tenantId` | `string` | Yes | owning tenant; empty coalesces to `"default"`, and the caller must be authorized for the resulting tenant |
| `code` | `string` | Yes | unique flow business code; immutable after creation (`start` targets it) |
| `name` | `string` | Yes | display name |
| `categoryId` | `string` | Yes | owning `FlowCategory` id |
| `icon` | `string` | No | display icon identifier |
| `description` | `string` | No | free-text description |
| `labels` | `object` (string→string) | No | host-owned selection metadata; validated by the shared label rule — see the note below |
| `bindingMode` | `string` | Yes | `standalone` (form data lives in approval tables) or `business` (links an existing business row) |
| `businessBinding` | `BusinessBindingConfig` | business mode | write-back target description (below); rejected on standalone flows (`ErrBindingUnexpected`) |
| `adminUserIds` | `string[]` | No | flow administrators (used by `transfer_admin` empty-assignee handling and admin visibility) |
| `isAllInitiationAllowed` | `bool` | No | `true` lets every user initiate; `false` restricts initiation to `initiators` |
| `instanceTitleTemplate` | `string` | No | Go `text/template` for instance titles, e.g. `{{.applicantName}}的请假申请`; bindings: `flowName`, `flowCode`, `instanceNo`, `formData`, `applicantId`, `applicantName` (plus nested `flow.name` / `flow.code`, `applicant.id` / `applicant.name`). Empty falls back to `flowName-instanceNo`; parse failure raises `ErrInvalidTitleTemplate` |
| `initiators` | `CreateInitiatorParams[]` | No | who may initiate when `isAllInitiationAllowed` is `false` (below) |

`CreateInitiatorParams` entries:

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `kind` | `string` | Yes | `user`, `role`, or `department` |
| `ids` | `string[]` | Yes | ids of the selected users / roles / departments |

`BusinessBindingConfig`:

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `tableName` | `string` | Yes | business table receiving approval state; validated as a SQL-safe identifier (`ErrInvalidBusinessIdentifier`) and checked to exist (`ErrBindingSchemaInvalid`) |
| `keyColumns` | `string[]` | Yes | columns locating the bound row; must exactly match a non-null primary or unique key (`ErrBindingKeyNotUnique`) |
| `statusColumn` | `string` | Yes | column receiving the (mapped) instance status |
| `instanceIdColumn` | `string` | Yes | column receiving the owning instance id; used as a compare-and-set fence so a stale instance cannot overwrite a newer approval round |
| `startedAtColumn` | `string` | No | column receiving the instance start time |
| `finishedAtColumn` | `string` | No | column receiving the instance finish time |
| `statusMapping` | `object` (InstanceStatus→string) | No | translates instance statuses into host vocabulary; missing entries fall back to the status string itself (`ErrBindingStatusMappingInvalid` for unknown keys or blank values) |

Two binding fields naming the same column fail with
`ErrBindingColumnsConflict`. Deployed versions snapshot their binding, so
editing a flow's binding never affects instances already running under
earlier versions. See [Integration](./integration.md) for the write-back
lifecycle.

`params.labels` is host-owned selection metadata on the flow —
equality-filterable in `find_flows` and `my.find_available_flows` (every
submitted pair must match), surfaced in instance detail views, and never
interpreted by the engine. Validation is the shared `orm.ValidateLabels`
rule: alphanumeric keys with inner `-`/`_` (no dots), ≤ 63 characters;
values ≤ 256 characters, empty values legal. On `update`, labels are
replaced wholesale — omitting them clears the flow's labels.

`DeployFlowParams` (`deploy` — creates a new draft version):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `flowId` | `string` | Yes | flow to deploy under |
| `description` | `string` | No | version description shown in version lists |
| `storageMode` | `string` | No | `json` (default; form data stays in `apv_instance.form_data`) or `table` (a dedicated physical projection table is generated at publish) |
| `flowDefinition` | `FlowDefinition` | Yes | the designer graph document — nodes, edges, and per-node `data`; validated at deploy (`ErrInvalidFlowDesign`). See [Flow Design](./flow-design.md) for the wire shape |
| `formSchema` | JSON document | No | host-owned form designer document, passed through opaque and stored verbatim; the flat field list the engine consumes is derived from it at deploy (see [Form Schema and Derived Fields](./flow-design.md#form-schema-and-derived-fields)) |

`PublishVersionParams` (`publish_version` — makes a draft the live version and
archives the previous one):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `versionId` | `string` | Yes | draft version to publish (`ErrVersionNotDraft` otherwise) |

`UpdateFlowParams` (`update` — mutates the flow row only; deployed versions
are immutable):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `flowId` | `string` | Yes | flow to update |
| `name` | `string` | Yes | display name |
| `icon` | `string` | No | display icon identifier |
| `description` | `string` | No | free-text description |
| `labels` | `object` (string→string) | No | replaced wholesale; omitting clears |
| `bindingMode` | `string` | Yes | binding mode (see `create`); changing it only affects future deployments |
| `businessBinding` | `BusinessBindingConfig` | business mode | write-back target (see `create`) |
| `adminUserIds` | `string[]` | No | flow administrators |
| `isAllInitiationAllowed` | `bool` | No | initiation openness |
| `instanceTitleTemplate` | `string` | Yes | instance title template |
| `initiators` | `CreateInitiatorParams[]` | No | initiator configuration; replaced wholesale |

`ToggleActiveParams` (`toggle_active`):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `flowId` | `string` | Yes | flow to toggle |
| `isActive` | `bool` | No | target state; inactive flows refuse initiation (`ErrFlowNotActive`) while running instances continue |

`GetGraphParams` (`get_graph`):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `flowId` | `string` | Yes | flow whose graph to load |
| `tenantId` | `string` | No | optional pre-filter; the actual cross-tenant gate is the caller's tenant authority |
| `versionId` | `string` | No | explicit version to load — a designer resuming from the newest deployment, published or not; omitted resolves the latest published version |

`get_graph` responds with a `FlowGraph`:

| Field | Type | Description |
| --- | --- | --- |
| `flow` | `Flow` | the mutable flow row (fields below) |
| `version` | `FlowVersion` | the resolved version, including `flowSchema` (the deployed `FlowDefinition`), `formSchema` (host document, verbatim), and `formFields` (the derived flat field list) |
| `nodes` | `FlowNode[]` | persisted node rows of that version — one row per node with all resolved node configuration (kind, execution type, approval method, pass rule, rollback / add-assignee / CC toggles, timeout config, branches). See [Flow Design](./flow-design.md) for each field's semantics |
| `edges` | `FlowEdge[]` | persisted edge rows: `key`, `sourceNodeId` / `sourceNodeKey`, `targetNodeId` / `targetNodeKey`, `sourceHandle` (condition-branch anchor) |

`FindFlowsParams` (`find_flows`):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `tenantId` | `string` | No | tenant filter; non-super-admin callers are constrained to their own tenant regardless |
| `categoryId` | `string` | No | filter by category |
| `keyword` | `string` | No | contains match against the flow name |
| `isActive` | `bool` | No | filter by active flag |
| `labels` | `object` (string→string) | No | label equality filter — every submitted pair must match |
| `page` | `int` | No | page number (1-based) |
| `pageSize` | `int` | No | page size |

`find_flows` responds with `page.Page[Flow]`. `Flow` (response model):

| Field | Type | Description |
| --- | --- | --- |
| `tenantId` | `string` | owning tenant |
| `categoryId` | `string` | owning category |
| `code` | `string` | unique flow business code (immutable) |
| `name` | `string` | display name |
| `icon` | `string` \| `null` | display icon identifier |
| `description` | `string` \| `null` | free-text description |
| `labels` | `object` \| absent | host-owned selection metadata |
| `bindingMode` | `string` | `standalone` or `business` |
| `businessBinding` | `BusinessBindingConfig` \| absent | current write-back configuration (mutable copy; versions snapshot their own) |
| `adminUserIds` | `string[]` | flow administrators |
| `isAllInitiationAllowed` | `bool` | initiation openness |
| `instanceTitleTemplate` | `string` | instance title template |
| `isActive` | `bool` | active flag |
| `currentVersion` | `int` | latest published version number; `0` before the first publish |

`FindVersionsParams` / `FindInitiatorsParams` (`find_versions`,
`find_initiators`):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `flowId` | `string` | Yes | flow to inspect |
| `tenantId` | `string` | No | optional pre-filter (cross-tenant gate is the caller's authority) |

`find_versions` returns `FlowVersionSummary` entries — the
version list without the graph documents (`flowSchema` / `formSchema` /
`formFields`), which a list never renders. Fetch one version's full definition
through `get_graph` with `params.versionId`.

| `FlowVersionSummary` field | Type | Description |
| --- | --- | --- |
| `id` | `string` | version id |
| `flowId` | `string` | owning flow |
| `version` | `int` | monotonically increasing version number |
| `status` | `string` | `draft`, `published`, or `archived` |
| `description` | `string` \| `null` | version description |
| `storageMode` | `string` | `json` or `table` |
| `publishedAt` | `DateTime` \| `null` | publish time |
| `publishedBy` | `string` \| `null` | publisher user id |
| `createdAt` | `DateTime` | deploy time |
| `createdBy` | `string` | deployer user id |

`find_initiators` returns `FlowInitiator[]`: each entry carries `flowId`,
`kind` (`user` / `role` / `department`), and `ids` (the configured id list).

## `approval/instance`

Instance lifecycle commands. Every state change is recorded in the action log;
the operations marked audited additionally capture framework-level IP / UA /
request-id audit entries.

| Action | Permission | Input | Output | Audit |
| --- | --- | --- | --- | --- |
| `start` | `approval.instance.start` | `StartParams` | created `Instance` | Yes |
| `process_task` | `approval.task.process` | `ProcessTaskParams` | success | Yes |
| `withdraw` | `approval.instance.withdraw` | `WithdrawParams` | success | Yes |
| `resubmit` | `approval.instance.resubmit` | `ResubmitParams` | success | Yes |
| `add_cc` | `approval.instance.cc` | `AddCCParams` | success | Yes |
| `mark_cc_read` | `approval.instance.cc` | `MarkCCReadParams` | success | — |
| `add_assignee` | `approval.task.add_assignee` | `AddAssigneeParams` | success | Yes |
| `remove_assignee` | `approval.task.remove_assignee` | `RemoveAssigneeParams` | success | Yes |
| `urge_task` | `approval.task.urge` | `UrgeTaskParams` | success | rate-limited: max `10` per `1m` |

`process_task` deliberately bundles approve / reject / transfer / rollback /
handle under one permission (`approval.task.process`): the designer's
node-level toggles (`isTransferAllowed`, `isRollbackAllowed`, …) already
govern which actions a node offers at runtime.

`StartParams` (`start`):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `tenantId` | `string` | Yes | tenant to start under; empty coalesces to `"default"`, and the caller must be authorized for the flow's tenant |
| `flowCode` | `string` | Yes | business code of the flow to start; resolves the latest published version (`ErrFlowNotFound` / `ErrFlowNotActive` / `ErrNoPublishedVersion`) |
| `businessRef` | `string` (≤ 512) | business mode | opaque reference to the bound business row; required on business-bound flows unless a registered `BusinessRefProvider` supplies it (`ErrBusinessRefRequired`). Default shapes: single key verbatim, composite key as a JSON object |
| `formData` | `object` | No | form values keyed by field key; validated against the published version's derived field list (`40401` family), rejected above 64 KiB, and stripped of fields the applicant may not edit |

The applicant identity and the condition-routing globals are resolved
server-side from the authenticated principal (`PrincipalDepartmentResolver`,
`InstanceGlobalsResolver`) — they are never accepted from the request body,
where an applicant could forge them to steer the flow.

`start` responds with the created `Instance`:

| Field | Type | Description |
| --- | --- | --- |
| `tenantId` | `string` | owning tenant |
| `flowId` / `flowCode` / `flowVersionId` | `string` | the flow and the immutable version snapshot the instance runs under |
| `title` | `string` | rendered from the flow's `instanceTitleTemplate` |
| `instanceNo` | `string` | human-readable instance number |
| `applicantId` / `applicantName` | `string` | applicant snapshot taken at start |
| `applicantDepartmentId` / `applicantDepartmentName` | `string` \| `null` | applicant department snapshot |
| `status` | `string` | `running`, `approved`, `rejected`, `withdrawn`, `returned`, or `terminated` |
| `currentNodeId` | `string` \| `null` | node the instance currently sits on |
| `finishedAt` | `DateTime` \| `null` | set when the instance reaches a final status |
| `businessRef` | `string` \| `null` | opaque business reference (business-bound flows) |
| `formData` | `object` | submitted form data (post-validation) |
| `globals` | `object` | host-supplied global-variable snapshot taken at start; condition evaluation reads it, so routing stays deterministic |
| `businessProjectionId` | `string` \| absent | durable write-back state claimed at start (business-bound flows) |

`ProcessTaskParams` (`process_task`):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `taskId` | `string` | Yes | pending task to act on; the caller must be its assignee (`ErrNotAssignee`, `ErrTaskNotPending`) |
| `action` | `string` | Yes | `approve`, `reject`, `transfer`, `rollback`, or `handle` (handle nodes finish with `handle`; same semantics as approve) |
| `opinion` | `string` (≤ 2000) | conditional | decision comment; required when the node sets `isOpinionRequired` (`ErrOpinionRequired`) |
| `formData` | `object` | No | form updates written with the action, filtered by the node's field permissions |
| `attachments` | `string[]` (≤ 20 × ≤ 512) | No | attachment references stored on the action log |
| `transferToId` | `string` | `transfer` | target user; must be non-empty and different from the operator (`ErrInvalidTransferTarget`); allowed only when the node sets `isTransferAllowed` (`ErrTransferNotAllowed`) |
| `targetNodeId` | `string` | `rollback` | rollback destination node; must be one of the node's valid targets per its `rollbackType` and the instance's visit trail (`ErrInvalidRollbackTarget`, `ErrRollbackNotAllowed`). Valid targets are served in `my.get_instance_detail` → `myTask.rollbackTargets` |

`WithdrawParams` (`withdraw` — applicant pulls a running instance back):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `instanceId` | `string` | Yes | instance to withdraw; caller must be the applicant (`ErrNotApplicant`), state must allow it (`ErrWithdrawNotAllowed`) |
| `reason` | `string` (≤ 2000) | No | withdraw reason recorded in the action log |

`ResubmitParams` (`resubmit` — restart a returned or withdrawn instance):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `instanceId` | `string` | Yes | instance to resubmit (`ErrResubmitNotAllowed` outside `returned` / `withdrawn`) |
| `formData` | `object` | No | replacement form data; validated like `start` |

`AddCCParams` / `MarkCCReadParams` (`add_cc`, `mark_cc_read`):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `instanceId` | `string` | Yes | target instance |
| `ccUserIds` | `string[]` (1–50) | `add_cc` only | users to CC; allowed only when the current node sets `isManualCcAllowed` (`ErrManualCcNotAllowed`) |

`mark_cc_read` stamps the read receipt on all of the caller's unread CC
records for the instance — a self-service read receipt, hence no audit.

`AddAssigneeParams` (`add_assignee` — dynamic assignee insertion):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `taskId` | `string` | Yes | the caller's own pending task |
| `userIds` | `string[]` (1–50) | Yes | users to add |
| `addType` | `string` | Yes | `before` (new assignee first, original waits), `after` (new assignee after the original completes), or `parallel` (joins the current group). Must be one of the node's `addAssigneeTypes` (`ErrAddAssigneeNotAllowed` / `ErrInvalidAddAssigneeType`) |

`RemoveAssigneeParams` (`remove_assignee`):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `taskId` | `string` | Yes | peer task to cancel; must be a still-actionable peer of the caller's own visit, not the last active assignee (`ErrLastAssigneeRemoval`), and the node must allow removal (`ErrRemoveAssigneeNotAllowed`). Eligible peers are served in `myTask.removableAssignees` |

`UrgeTaskParams` (`urge_task`):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `taskId` | `string` | Yes | pending task to urge |
| `message` | `string` (≤ 500) | No | urge message delivered with the notification |

Urges honor the node's `urgeCooldownMinutes` per task (`40601` when urged too
frequently; non-positive config defaults to 30 minutes), and the operation
carries an extra rate limit of 10 calls per minute per caller.

## `approval/my`

Self-service queries for the current user. The operations declare no
`RequiredPermission` — any authenticated principal may call them; every query
is keyed to the caller's identity server-side.

| Action | Input | Output |
| --- | --- | --- |
| `find_available_flows` | `FindAvailableFlowsParams` | `page.Page[AvailableFlow]` |
| `get_start_form` | `GetStartFormParams` | `StartForm` |
| `find_initiated` | `FindInitiatedParams` | `page.Page[InitiatedInstance]` |
| `find_pending_tasks` | `FindPendingTasksParams` | `page.Page[PendingTask]` |
| `find_completed_tasks` | `FindCompletedTasksParams` | `page.Page[CompletedTask]` |
| `find_cc_records` | `FindCCRecordsParams` | `page.Page[CCRecord]` |
| `get_pending_counts` | `GetPendingCountsParams` | `PendingCounts` |
| `get_instance_detail` | `GetInstanceDetailParams` | `InstanceDetail` |

Request parameters (all decoded from `params`):

| Action | Field | Type | Required | Description |
| --- | --- | --- | --- | --- |
| `find_available_flows` | `tenantId` | `string` | No | tenant filter |
| | `keyword` | `string` | No | contains match against the flow name |
| | `labels` | `object` | No | label equality filter — every pair must match |
| | `page` / `pageSize` | `int` | No | pagination |
| `get_start_form` | `tenantId` | `string` | Yes | tenant of the flow |
| | `flowCode` | `string` | Yes | flow to load the start form for |
| `find_initiated` | `tenantId` | `string` | No | tenant filter |
| | `status` | `string` | No | instance status filter (`running` / `approved` / `rejected` / `withdrawn` / `returned` / `terminated`) |
| | `keyword` | `string` | No | contains match against the instance title |
| | `page` / `pageSize` | `int` | No | pagination |
| `find_pending_tasks` | `tenantId` | `string` | No | tenant filter |
| | `page` / `pageSize` | `int` | No | pagination |
| `find_completed_tasks` | `tenantId` | `string` | No | tenant filter |
| | `page` / `pageSize` | `int` | No | pagination |
| `find_cc_records` | `tenantId` | `string` | No | tenant filter |
| | `isRead` | `bool` | No | read-state filter |
| | `page` / `pageSize` | `int` | No | pagination |
| `get_pending_counts` | `tenantId` | `string` | No | tenant filter |
| `get_instance_detail` | `instanceId` | `string` | Yes | instance to load; the caller must be a participant — applicant, assignee, or CC recipient (`ErrAccessDenied`) |

Response DTOs (`approval/my` package):

`AvailableFlow` — one flow the caller may initiate:

| Field | Type | Description |
| --- | --- | --- |
| `flowId` / `flowCode` / `flowName` | `string` | flow identity |
| `flowIcon` | `string` \| absent | display icon |
| `description` | `string` \| absent | flow description |
| `labels` | `object` \| absent | host-owned selection metadata |
| `categoryId` / `categoryName` | `string` | owning category identity |

`StartForm` — the pre-submission view of a flow. Loading it is gated
exactly like starting the instance (active flow, initiation permission,
published version), so a rendered form always implies a startable flow:

| Field | Type | Description |
| --- | --- | --- |
| `flowId` / `flowCode` / `flowName` | `string` | flow identity for the initiation header |
| `flowIcon` | `string` \| absent | display icon |
| `description` | `string` \| absent | flow description |
| `versionId` | `string` | published version the form belongs to |
| `version` | `int` | published version number |
| `formSchema` | JSON document \| absent | host form-designer document, verbatim |

`InitiatedInstance` — one instance the caller submitted:

| Field | Type | Description |
| --- | --- | --- |
| `instanceId` / `instanceNo` / `title` | `string` | instance identity |
| `flowName` | `string` | flow display name |
| `flowIcon` | `string` \| absent | flow icon |
| `status` | `string` | instance status |
| `currentNodeName` | `string` \| absent | name of the node currently in progress |
| `createdAt` | `DateTime` | submission time |
| `finishedAt` | `DateTime` \| absent | completion time |

`PendingTask` — one task awaiting the caller's action:

| Field | Type | Description |
| --- | --- | --- |
| `taskId` | `string` | task to submit `process_task` against |
| `instanceId` / `instanceTitle` / `instanceNo` | `string` | owning instance identity |
| `flowName` / `flowIcon` | `string` | flow display identity |
| `applicant` | `UserInfo` | applicant snapshot |
| `nodeName` | `string` | node the task belongs to |
| `createdAt` | `DateTime` | task creation time |
| `deadline` | `DateTime` \| absent | timeout deadline when the node configures one |
| `isTimeout` | `bool` | whether the task is past its deadline |

`CompletedTask` — one task the caller already processed: same identity fields
as `PendingTask` plus `status` (the outcome — `approved`, `rejected`,
`handled`, `transferred`, `rolled_back`, …) and `finishedAt`; without
`deadline` / `isTimeout`.

`CCRecord` — one CC notification addressed to the caller:

| Field | Type | Description |
| --- | --- | --- |
| `ccRecordId` | `string` | CC record id |
| `instanceId` / `instanceTitle` / `instanceNo` | `string` | owning instance identity |
| `flowName` / `flowIcon` | `string` | flow display identity |
| `applicant` | `UserInfo` | applicant snapshot |
| `nodeName` | `string` \| absent | node that produced the CC; absent for instance-level CCs |
| `isRead` | `bool` | read receipt state |
| `createdAt` | `DateTime` | delivery time |

`PendingCounts` — badge counts: `pendingTaskCount` (tasks awaiting action) and
`unreadCcCount` (unread CC records).

`InstanceDetail` — the self-service detail view. Each top-level field is one
renderable concern:

| Field | Type | Description |
| --- | --- | --- |
| `instance` | `InstanceInfo` | runtime state (below) |
| `formSchema` | JSON document \| absent | version-pinned host form-designer document, verbatim — the schema the instance was submitted under |
| `timeline` | `TimelineEntry[]` | node-by-node account of the path actually taken (below) |
| `flowGraph` | `InstanceFlowGraph` | React Flow–ready read-only graph annotated with progress (below) |
| `availableActions` | `string[]` | viewer-specific action hints (below) |
| `fieldPermissions` | `object` (field→permission) | viewer-scoped field interactivity: `visible` / `editable` / `hidden` / `required`, materialized for every top-level form field; the client applies it verbatim, and `instance.formData` is already stripped of fields the viewer may not see (see [Node Field Permissions](./flow-design.md#node-field-permissions)) |
| `myTask` | `ViewerTask` \| `null` | the viewer's own actionable context (below) |

`InstanceInfo`:

| Field | Type | Description |
| --- | --- | --- |
| `instanceId` / `instanceNo` / `title` | `string` | instance identity |
| `flowName` / `flowIcon` | `string` | flow display identity, read from the mutable flow at query time |
| `labels` | `object` \| absent | the flow's host-owned selection metadata — display identity like `flowName`, not a version-pinned snapshot |
| `applicant` | `UserInfo` | applicant snapshot |
| `status` | `string` | instance status |
| `currentNodeId` / `currentNodeName` | `string` \| absent | node currently in progress |
| `businessRef` | `string` \| absent | opaque business reference (business-bound flows) |
| `formData` | `object` \| absent | form data, stripped of fields the viewer may not see |
| `createdAt` / `finishedAt` | `DateTime` | lifecycle timestamps |

`ViewerTask` — the pending task `process_task` should target plus the
node-level configuration the client needs to build the action UI without
re-deriving engine semantics. `null` when the viewer holds no pending task on
this instance:

| Field | Type | Description |
| --- | --- | --- |
| `taskId` | `string` | the pending task |
| `nodeId` | `string` | its node |
| `isOpinionRequired` | `bool` | mirrors the node config: approve / reject must carry a non-empty opinion when set |
| `addAssigneeTypes` | `string[]` | positions the node allows for dynamic assignee addition (`before` / `after` / `parallel`); empty when adding is not allowed |
| `rollbackTargets` | `{nodeId, name}[]` | valid rollback destinations, resolved from the node's rollback config and the instance's visit trail exactly like the rollback command validates them; empty when rollback is not allowed |
| `removableAssignees` | `{taskId, assignee, status}[]` | peer tasks the viewer may remove (`status` is `pending` / `waiting`), resolved exactly like the remove-assignee command authorizes them: still-actionable peers of the viewer's own visit, excluding the viewer; empty when removal is disallowed |

`availableActions` is a query-layer UI hint. For the applicant it includes
`withdraw` when the instance can transition to `withdrawn`, and `resubmit`
when the instance is returned or withdrawn. For pending tasks it includes
`handle` for handle nodes, otherwise `approve`, then `reject`, plus
`transfer`, `rollback`, `add_assignee`, or `add_cc` when the current node
allows them. If the instance has any pending task, it also includes `urge`.
Command handlers still perform their own validation.

## `approval/admin`

Admin-level management and observability. For every list and metrics query,
non-super-admin callers ignore a submitted `tenantId` override and are
filtered to their own tenant; super-admin callers may pass `tenantId` to
filter one tenant or omit it for cross-tenant visibility.

| Action | Permission | Input | Output | Audit |
| --- | --- | --- | --- | --- |
| `find_instances` | `approval.instance.query` | `AdminFindInstancesParams` | `page.Page[Instance]` | — |
| `find_tasks` | `approval.task.query` | `AdminFindTasksParams` | `page.Page[Task]` | — |
| `get_instance_detail` | `approval.instance.detail` | `AdminGetInstanceDetailParams` | `InstanceDetail` | — |
| `find_action_logs` | `approval.action_log.query` | `AdminFindActionLogsParams` | `page.Page[ActionLog]` | — |
| `get_metrics` | `approval.metrics.query` | `AdminGetMetricsParams` | `Metrics` | — |
| `find_business_projections` | `approval.binding.query` | `AdminFindBusinessProjectionsParams` | `page.Page[BusinessProjection]` | — |
| `terminate_instance` | `approval.instance.terminate` | `AdminTerminateInstanceParams` | success | Yes |
| `reassign_task` | `approval.task.reassign` | `AdminReassignTaskParams` | success | Yes |
| `retry_business_projection` | `approval.binding.retry` | `AdminRetryBusinessProjectionParams` | success | Yes |

Request parameters (all decoded from `params`):

| Action | Field | Type | Required | Description |
| --- | --- | --- | --- | --- |
| `find_instances` | `tenantId` | `string` | No | tenant filter (super-admin only, see above) |
| | `applicantId` | `string` | No | filter by applicant |
| | `status` | `string` | No | instance status filter |
| | `flowId` | `string` | No | filter by flow |
| | `keyword` | `string` | No | contains match against the instance title |
| | `page` / `pageSize` | `int` | No | pagination |
| `find_tasks` | `tenantId` | `string` | No | tenant filter |
| | `assigneeId` | `string` | No | filter by assignee |
| | `instanceId` | `string` | No | filter by owning instance |
| | `status` | `string` | No | task status filter (`waiting` / `pending` / `approved` / `rejected` / `handled` / `transferred` / `rolled_back` / `canceled` / `removed` / `skipped`) |
| | `page` / `pageSize` | `int` | No | pagination |
| `get_instance_detail` | `instanceId` | `string` | Yes | instance to load |
| `find_action_logs` | `instanceId` | `string` | Yes | instance whose audit trail to page through |
| | `tenantId` | `string` | No | tenant filter |
| | `page` / `pageSize` | `int` | No | pagination |
| `get_metrics` | `tenantId` | `string` | No | tenant scope (super-admin may omit for cross-tenant) |
| `find_business_projections` | `tenantId` | `string` | No | tenant filter |
| | `status` | `string` | No | projection status filter: `pending`, `processing`, `applied`, `failed` |
| | `page` / `pageSize` | `int` | No | pagination |
| `terminate_instance` | `instanceId` | `string` | Yes | running instance to force-terminate (`ErrTerminateNotAllowed` outside running states) |
| | `reason` | `string` (≤ 2000) | No | termination reason recorded in the action log |
| `reassign_task` | `taskId` | `string` | Yes | pending task to reassign |
| | `newAssigneeId` | `string` | Yes | replacement assignee (`ErrInvalidTransferTarget` when invalid) |
| | `reason` | `string` (≤ 2000) | No | reassignment reason |
| `retry_business_projection` | `projectionId` | `string` | Yes | projection to retry immediately (`ErrBindingProjectionNotFound` when missing) |

Response DTOs (`approval/admin` package):

`Instance` — one instance in the admin list:

| Field | Type | Description |
| --- | --- | --- |
| `instanceId` / `instanceNo` / `title` | `string` | instance identity |
| `tenantId` | `string` | owning tenant |
| `flowId` / `flowName` | `string` | flow identity |
| `applicant` | `UserInfo` | applicant snapshot |
| `status` | `string` | instance status |
| `currentNodeName` | `string` \| absent | node currently in progress |
| `createdAt` / `finishedAt` | `DateTime` | lifecycle timestamps |

`Task` — one task in the admin list:

| Field | Type | Description |
| --- | --- | --- |
| `taskId` | `string` | task id |
| `instanceId` / `instanceTitle` | `string` | owning instance identity |
| `flowName` | `string` | flow display name |
| `nodeName` | `string` | node the task belongs to |
| `assignee` | `UserInfo` | assignee snapshot |
| `status` | `string` | task status |
| `createdAt` | `DateTime` | creation time |
| `deadline` | `DateTime` \| absent | timeout deadline |
| `finishedAt` | `DateTime` \| absent | completion time |

`InstanceDetail` — the admin counterpart of `my.get_instance_detail`, without
the viewer-specific fields (`availableActions` / `fieldPermissions` /
`myTask`): `instance` (`InstanceDetailInfo`), `formSchema` (verbatim host
document), `timeline` (`TimelineEntry[]`), and `flowGraph`
(`InstanceFlowGraph`). `InstanceDetailInfo` matches `my.InstanceInfo` plus
`tenantId`, `flowId`, and `flowVersionId`, and its `formData` is unfiltered.

`ActionLog` — one audit entry. Person references are uniform `UserInfo`
snapshots captured at action time:

| Field | Type | Description |
| --- | --- | --- |
| `logId` | `string` | log entry id |
| `action` | `string` | `ActionType` string: `submit`, `approve`, `handle`, `reject`, `transfer`, `withdraw`, `cancel`, `rollback`, `add_assignee`, `remove_assignee`, `execute`, `resubmit`, `reassign`, `terminate`, `add_cc` |
| `nodeId` | `string` \| absent | node the action happened at; absent for instance-level actions |
| `taskId` | `string` \| absent | task the action targeted |
| `operator` | `UserInfo` | acting user snapshot |
| `transferTo` | `UserInfo` \| absent | transfer / reassignment recipient |
| `rollbackToNodeId` | `string` \| absent | rollback destination |
| `addedAssignees` / `removedAssignees` | `UserInfo[]` \| absent | dynamic assignee changes |
| `ccUsers` | `UserInfo[]` \| absent | manually CC'd users |
| `opinion` | `string` \| absent | action comment / reason |
| `attachments` | `string[]` \| absent | attachment references |
| `createdAt` | `DateTime` | action time |

`Metrics` — aggregated engine health for dashboards and ops alerting:

| Field | Type | Description |
| --- | --- | --- |
| `tenantId` | `string` | tenant scope of the snapshot; empty for a cross-tenant snapshot (super-admin only) |
| `capturedAt` | `DateTime` | when the metrics were materialized |
| `instanceCounts` | `object` (status→int) | instance counts keyed by `InstanceStatus` string |
| `taskCounts` | `object` (status→int) | task counts keyed by `TaskStatus` string |
| `timeoutTaskCount` | `int` | pending tasks past their deadline |
| `avgCompletionSeconds` | `float` | average end-to-end duration (`createdAt` → `finishedAt`) over all finalized instances; `-1` means "no completed instances yet" |
| `pendingBindingFailures` | `int` | projection targets whose latest write attempt failed and is scheduled for retry |
| `businessProjectionCounts` | `object` (status→int) | durable projection rows by convergence status (`pending` / `processing` / `applied` / `failed`) |
| `pendingBusinessProjections` | `int` | eventual projections whose desired revision has not been applied yet |

`BusinessProjection` — the operator-facing convergence state for one bound
business record (see [Integration](./integration.md) for the write-back
model):

| Field | Type | Description |
| --- | --- | --- |
| `projectionId` | `string` | projection row id |
| `tenantId` | `string` | owning tenant |
| `flowId` / `flowVersionId` | `string` | flow and version that own the desired state |
| `ownerInstanceId` | `string` | instance whose lifecycle produced the desired state |
| `appliedOwnerInstanceId` | `string` \| absent | instance whose state was last successfully written to the business row |
| `businessTable` | `string` | target business table |
| `recordKey` | JSON object | key-column values locating the bound row |
| `consistency` | `string` | binding consistency mode from configuration (`transactional` / `eventual`) |
| `desiredStatus` | `string` | instance status awaiting write-back |
| `desiredStartedAt` / `desiredFinishedAt` | `DateTime` | lifecycle timestamps awaiting write-back |
| `desiredRevision` / `appliedRevision` | `int` | monotonic revisions; the projection has converged when they are equal |
| `status` | `string` | convergence state: `pending`, `processing`, `applied`, `failed` |
| `attemptCount` | `int` | write attempts so far |
| `nextAttemptAt` | `DateTime` \| absent | next scheduled retry |
| `leaseUntil` | `DateTime` \| absent | worker lease expiry while `processing` |
| `lastError` | `string` \| absent | last write failure message |
| `appliedAt` | `DateTime` \| absent | when the desired state was last applied |
| `updatedAt` | `DateTime` | last state change |

## Shared Projection Types

The detail views (`my.get_instance_detail`, `admin.get_instance_detail`) share
these types from the public `approval` package.

`UserInfo` — the uniform person snapshot used everywhere a person appears:

| Field | Type | Description |
| --- | --- | --- |
| `id` | `string` | user id |
| `name` | `string` | display name at action time |
| `departmentId` / `departmentName` | `string` \| absent | department snapshot at action time |

`TimelineEntry` — one step of the instance timeline: the chronological,
node-by-node account of the path an instance actually took. Because condition
branches are exclusive, the traversed path is always a single line; a node
re-entered after a rollback produces a second entry. Entries end at the node
currently in progress — unreached nodes are not predicted:

| Field | Type | Description |
| --- | --- | --- |
| `kind` | `string` | `start`, `approval`, `handle`, `cc` for node visits; `withdraw`, `terminate` for instance-level milestones. Structural kinds (`condition` / `end`) never appear |
| `nodeId` | `string` \| absent | visited node; absent on milestone entries |
| `name` | `string` | node display name (or milestone action name) |
| `status` | `string` | node-visit status: `active`, `passed`, `rejected`, `returned`, `canceled` |
| `executionType` | `string` | node execution type (`manual` / `auto_pass` / `auto_reject`) |
| `approvalMethod` | `string` | `sequential` / `parallel` (approval nodes) |
| `passRule` | `string` | `all` / `any` / `ratio` (approval nodes) |
| `passRatio` | `decimal` \| absent | ratio threshold when `passRule` is `ratio` |
| `participants` | `NodeParticipant[]` | one entry per task at approval / handle nodes (below) |
| `ccRecipients` | `CCRecipient[]` | delivered carbon copies: `user` (`UserInfo`) plus `readAt` read receipt |
| `activities` | `Activity[]` | side actions at the node (below); milestone entries hold a single activity describing who closed the instance and why |
| `startedAt` / `finishedAt` | `DateTime` | visit span; `finishedAt` absent while in progress |

`NodeParticipant` — one assignee's involvement during a single visit:

| Field | Type | Description |
| --- | --- | --- |
| `taskId` | `string` | task identity (what task operations target) |
| `user` | `UserInfo` | assignee snapshot |
| `delegator` | `UserInfo` \| absent | present when the task arrived via delegation |
| `status` | `string` | task status verbatim |
| `deadline` | `DateTime` \| absent | task deadline |
| `isTimeout` | `bool` | task was decided or escalated by the timeout scanner |
| `opinion` / `attachments` / `actionTime` | — | outcome details fused from the action log that finished the task |
| `transferTo` | `UserInfo` \| absent | transfer recipient when the task was transferred |

`Activity` — a side action recorded at a node: `action` carries the
`ActionType` string (`transfer`, `rollback`, `add_assignee`,
`remove_assignee`, `add_cc`, `reassign`, `execute`, `submit`, `resubmit`,
`withdraw`, `terminate`) plus `urge` for urge records. `operator` is the
acting user; `opinion` holds the action's free text (a transfer reason, a
withdraw reason, an urge message); `target` names the counterpart of a
directed action (the urged assignee); `transferTo`, `rollbackToNodeId` /
`rollbackToNodeName`, `addedAssignees`, `removedAssignees`, `ccUsers`, and
`attachments` carry the action-specific details, and `createdAt` the action
time. Decisions themselves (approve / handle / reject) are not repeated as
activities — they live on the participant that made them.

`InstanceFlowGraph` — a React Flow–ready, read-only projection of the
instance's pinned flow definition annotated with runtime progress. `nodes`
and `edges` map directly onto React Flow's shape, except the node kind stays
in `kind` (React Flow's `type` belongs to the client):

| Field | Type | Description |
| --- | --- | --- |
| `nodes[].id` | `string` | React Flow identity — the design-time node key that positions and edges reference |
| `nodes[].nodeId` | `string` | persistent flow-node id — the value `actionLog.nodeId` / `rollbackToNodeId` carry and the `process_task` rollback API expects as `targetNodeId` |
| `nodes[].kind` | `string` | node kind (`start` / `approval` / `handle` / `condition` / `cc` / `end`) |
| `nodes[].position` | `{x, y}` | designer coordinates |
| `nodes[].data` | `FlowGraphNodeData` | node label, approval semantics, progress `status` (`pending` / `active` / `passed` / `rejected` / `returned` / `canceled`), plus `participants` / `ccRecipients` / `activities` aggregated across the node's visits in traversal order, and the `startedAt` / `finishedAt` span |
| `edges[]` | `{id, source, target, sourceHandle}` | React Flow edges connecting nodes by their ids |

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
| `40015` | — | — | — | unused (former flow-binding lock); the code is never reassigned |
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
