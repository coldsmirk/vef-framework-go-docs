---
sidebar_position: 3
---

# Flow Design

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
`Operator` is `ConditionOperator`; the exported constants are
`OperatorEquals`, `OperatorNotEquals`, `OperatorGreater`,
`OperatorGreaterOrEq`, `OperatorLess`, `OperatorLessOrEq`, `OperatorIn`,
`OperatorNotIn`, `OperatorContains`, `OperatorNotContains`,
`OperatorStartsWith`, `OperatorEndsWith`, `OperatorIsEmpty`, and
`OperatorIsNotEmpty`. The built-in evaluator converts field conditions to an
`expr-lang` expression; unknown operators evaluate to `false`.

`ConditionExpression` evaluates the raw `Expression` string with `expr-lang`.
The evaluation environment exposes:

| Name | Value |
| --- | --- |
| `formData` | the instance `FormData` as a map |
| `applicantId` | current applicant ID |
| `applicantDepartmentId` | applicant department ID, or `""` when absent |
| globals | host-resolved `Instance.Globals` values exposed as top-level bindings |

Approval conditions intentionally use `expr-lang` directly, not the public
`expression.Engine` wrapper. That keeps workflow condition semantics tied to
the approval evaluator and independent of expression module wiring.

Host applications can implement `approval.InstanceGlobalsResolver` to resolve
global variables from the authenticated principal at instance start. The
snapshot is persisted on `Instance.Globals`; clients cannot submit it in the
`start` request. Field conditions resolve `Subject` against globals before
`formData`, and expression conditions expose globals as top-level bindings while
the built-in `formData`, `applicantId`, and `applicantDepartmentId` names win
collisions.

### Detail-Table Aggregation

A field condition may fold a detail table's rows instead of comparing a scalar
subject (v0.36). The condition stays structured — no string DSL: `subject`
names the table field, `aggregate` picks the fold, and `column` names the
numeric column to fold.

| `AggregateKind` | Wire value | `column` | Folds |
| --- | --- | --- | --- |
| `AggregateSum` | `sum` | required | sum of the named numeric column |
| `AggregateCount` | `count` | forbidden | row count |
| `AggregateAvg` | `avg` | required | average of the named numeric column |

`AggregateKind.FoldsColumn()` reports whether a kind reduces a column (`sum` /
`avg`) rather than rows (`count`). Folding is pluggable through the
`approval.Aggregator` interface:

```go
type Aggregator interface {
    // Kind returns the aggregate kind this implementation folds.
    Kind() AggregateKind
    // Fold reduces the extracted column values (or the row count) into the
    // comparison operand. matchable=false means the aggregate has no defined
    // value for the input — e.g. avg over zero rows — and the condition must
    // not match, mirroring SQL NULL comparison semantics.
    Fold(values []float64, rowCount int) (result float64, matchable bool)
}
```

Register a custom aggregator alongside the built-in `sum` / `count` / `avg`
with `vef.ProvideApprovalAggregator`; the condition evaluator picks it up by
its `AggregateKind` with no changes to existing code:

```go
vef.Run(
    vef.ApprovalModule,
    vef.ProvideApprovalAggregator(func() approval.Aggregator { return myMedian{} }),
    app.Module,
)
```

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

## Node Field Permissions

Task nodes (`TaskNodeData`, embedded by approval and handle nodes) and CC nodes
(`CCNodeData`) carry a `fieldPermissions` map: form-field key → `Permission`.
The vocabulary is:

| Constant | Wire value | Meaning for the node's participants |
| --- | --- | --- |
| `PermissionVisible` | `visible` | read-only |
| `PermissionEditable` | `editable` | may submit a new value |
| `PermissionHidden` | `hidden` | not shown |
| `PermissionRequired` | `required` | editable and must be provided |

An absent key means `visible`. Deploy validation (v0.38) checks the map
against the derived form fields: every key must reference a top-level form
field, values must be in the enum, CC nodes may only use the `visible` /
`hidden` subset, and a `required` permission is rejected on a node whose
timeout action resolves to `auto_pass` (the timeout scanner's auto-pass
finishes tasks without the required check).

The map is enforced on the write path: during task processing, submitted
`formData` is merged only for fields whose `fieldPermissions` entry is
`editable` or `required`, and the merged subset is re-validated against the
field definitions. The `required` must-fill check applies to `approve` and
`handle` decisions only — `reject`, `transfer`, and `rollback` stay exempt. A
node with no `fieldPermissions` map grants no write access: submitted form
data is dropped.

On the read side, `my.get_instance_detail` returns a viewer-scoped
`fieldPermissions` projection: the framework max-merges the lattice
(`hidden < visible < editable < required`) over the viewer's participation
contexts (own task, CC delivery, applicant), grants write strength
(`editable` / `required`) only from a pending task or a resubmittable
applicant, and clamps every read-only context to `visible`. `hidden` values
are stripped from the returned `formData`, and a viewer with no recognized
context sees nothing — the resolution is fail-closed.

## Designer Defaults

`NodeData.ApplyTo` resolves omitted designer fields to exported defaults so the runtime matches an untouched designer control:

| Constant | Value |
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
uses `kind`, `subject`, `operator`, `value`, and `expression`. A field
condition may instead fold a detail table's rows: `aggregate` (`sum` /
`count` / `avg`) evaluates over the table field named by `subject`, and
`column` names the numeric column to fold (required for `sum` / `avg`,
forbidden for `count`) — see
[Detail-Table Aggregation](#detail-table-aggregation).

`timeoutHours` and `timeoutNotifyBeforeHours` are in hours.
`urgeCooldownMinutes` is in minutes; values less than or equal to 0 use the
runtime default of 30 minutes. `rollbackTargetKeys` is checked when
`rollbackType` is `specified`; it contains node keys, not database node IDs.
`fieldPermissions` semantics are described in
[Node Field Permissions](#node-field-permissions).

### Form Schema and Derived Fields

Since v0.38 the form definition is split in two at deploy:

- `FlowVersion.FormSchema` (`formSchema`) is the **host-owned form-designer
  document**, submitted at `deploy` as `params.formSchema` and stored /
  returned as semantically equal JSON — the jsonb column normalizes formatting
  and key order while numeric precision is preserved (`json.RawMessage` end to
  end). The framework never interprets it.
- `FlowVersion.FormFields` (`formFields`) is the flat `[]FormFieldDefinition`
  list derived from that document exactly once at deploy through the injected
  `approval.FormSchemaParser`, and is the **only** form shape the framework
  consumes — for form-data validation, storage-table DDL, aggregate checks,
  and field-permission resolution. Parser upgrades never affect
  already-deployed versions.

```go
type FormSchemaParser interface {
    ParseFormFields(ctx context.Context, schema json.RawMessage) ([]FormFieldDefinition, error)
}
```

The built-in parser understands the vef-framework-react form-editor document;
hosts with their own designer replace it wholesale with
`vef.ProvideApprovalFormSchemaParser(constructor)`. A nil or empty schema
yields no fields (a flow without a form); parser errors abort the deploy. The
`ctx` carries the deploy request's deadline — a host parser that performs I/O
must honor it.

The derived fields are validated at deploy (unique keys, known kinds,
compilable patterns, coherent bounds, single-level tables). Each
`FormFieldDefinition` entry uses
`key`, `kind`, `label`, `placeholder`, `defaultValue`, `isRequired`,
`options`, `validation`, `props`, `sortOrder`, `columnType`, `scale`, and
`columns`. Each option uses `label` and `value`. `columns` defines the row
shape of a table field (`kind` is `table`): each entry is itself a
`FormFieldDefinition` and must not declare its own `columns` — detail tables
are single-level. On the table field itself, `validation.minLength` /
`maxLength` bound the row count and `isRequired` means at least one row.

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

## Flow Models and Designer Enums

Flow design and persistence models exposed by the public package include
`FlowCategory`, `Flow`, `FlowVersion`, `FlowNode`, `FlowEdge`, `FlowInitiator`,
`FlowNodeAssignee`, `FlowNodeCC`, `FormFieldDefinition`,
`FormSnapshot`, `ActionLog`, `OperatorInfo`, and `UrgeRecord` (the structured
`FormDefinition` wrapper was removed in v0.38 — the host document is opaque
and only `FormFieldDefinition` remains a framework shape). Flow-version
status uses `VersionStatus`: `VersionDraft` (`draft`), `VersionPublished`
(`published`), and `VersionArchived` (`archived`).

Additional flow-designer enums:

| Enum | Wire values |
| --- | --- |
| `InitiatorKind` | `user`, `role`, `department` |
| `ExecutionType` | `manual`, `auto_pass`, `auto_reject` |
| `ConditionKind` | `field`, `expression` |
| `CCKind` | `user`, `role`, `department`, `form_field` |
| `CCTiming` | `always`, `on_approve`, `on_reject` |
| `FieldKind` | `input`, `textarea`, `select`, `number`, `date`, `upload`, `table` |
| `ColumnDataType` | `string`, `text`, `integer`, `decimal`, `boolean`, `date`, `datetime`, `json` |
| `Permission` | `visible`, `editable`, `hidden`, `required` |

---

Next: [Instance Runtime](./runtime.md) for what happens after a designed flow starts running.
