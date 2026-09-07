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
`OperatorIsNotEmpty`. The built-in evaluator compares field conditions
**natively in Go** — no expression templating, no injection surface — with
operator semantics typed per field kind; an operator outside the vocabulary
fails the evaluation with an error rather than silently evaluating to
`false`.

`ConditionExpression` evaluates the raw `Expression` string.
The evaluation environment exposes:

| Name | Value |
| --- | --- |
| `formData` | the instance `FormData` as a map |
| `applicantId` | current applicant ID |
| `applicantDepartmentId` | applicant department ID, or `""` when absent |
| globals | host-resolved `Instance.Globals` values exposed as top-level bindings |

Expression conditions are evaluated through the framework's
`expression.Engine` abstraction (currently backed by `expr-lang`), wired by
DI — the same engine documented in
[Expression Engine](../data-tools/expression).

Host applications can implement `approval.InstanceGlobalsResolver` to resolve
global variables from the authenticated principal at instance start. The
snapshot is persisted on `Instance.Globals`; clients cannot submit it in the
`start` request. Field conditions resolve `Subject` against globals before
`formData`, and expression conditions expose globals as top-level bindings while
the built-in `formData`, `applicantId`, and `applicantDepartmentId` names win
collisions.

### Detail-Table Aggregation

A field condition may fold a detail table's rows instead of comparing a scalar
subject. The condition stays structured — no string DSL: `subject`
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

### Pass Rules

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
| Department | `AssigneeDepartment` | `department` | Leaders of the configured departments |
| Self | `AssigneeSelf` | `self` | The applicant |
| Superior | `AssigneeSuperior` | `superior` | Direct superior |
| Dept Leader | `AssigneeDepartmentLeader` | `department_leader` | Leaders of the applicant's own department (single-level lookup) |
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

An absent key means `visible`. Deploy validation checks the map
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
| `HandleNodeData` | base fields + `TaskNodeData` fields; deploy always sets `approvalMethod` to `sequential` and `passRule` to `any` (handle nodes expose no such wire fields) |
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

The form definition is split in two at deploy:

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
compilable patterns, coherent bounds, single-level tables, coherent option
sources). Each `FormFieldDefinition` entry uses
`key`, `kind`, `label`, `placeholder`, `defaultValue`, `isRequired`,
`options`, `optionSource`, `validation`, `props`, `sortOrder`, `columnType`,
`scale`, and `columns`. Each option uses `label` and `value`. `columns`
defines the row shape of a table field (`kind` is `table`): each entry is
itself a `FormFieldDefinition` and must not declare its own `columns` —
detail tables are single-level. On the table field itself, `validation.minLength` /
`maxLength` bound the row count and `isRequired` means at least one row.

`validation` supports `minLength`, `maxLength`, `min`, `max`, `pattern`, and
`message`. Submitted `formData` is capped by `vef.approval.form_data_max_bytes`
after JSON encoding (default 64 KiB via `config.DefaultFormDataMaxBytes`), even
when the flow has no form schema. When a schema exists, extra form keys are
rejected; required fields reject absent, `null`, blank-string, and empty-array
values. `input`, `textarea`, and `date` fields must be strings and may use
`minLength`, `maxLength`, and `pattern`. `number` fields accept numeric JSON
values and may use `min` and `max`; a `number` field whose `columnType` is
`integer` additionally rejects fractional values. `select` fields validate scalar or array
values against `options` when options are present. `upload` fields accept a
non-blank string, a non-empty `[]string`, or a non-empty array of non-blank
strings. `validation.message` is used as the custom error message for
`pattern` mismatches; other validation failures use the module i18n messages.

#### Selection Options: Enumerated or Remote

A selection field carries its options in exactly one of two shapes, never in
both:

- `options` (`[]FieldOption`) — the option list was enumerable at deploy, so
  the projection wrote it out. A static source is enumerated whether it was
  configured inline on the field or reached through a form-global `ref`.
- `optionSource` (`*FieldOptionSource`) — the options could **not** be
  enumerated, so the projection emits the descriptor a consumer needs to fetch
  them itself. Only remote sources reach this shape.

A field with neither is free-form: nothing constrains its value.

The projection is **post-dereference**. A designer `ref` is resolved against
the form-global sources before the field is written, so a consumer never has to
chase a `dataSourceId`. A `ref` that points at nothing yields neither `options`
nor `optionSource` — exactly like a field with no source at all.

`FieldOptionSource` is:

| Field | Type | Meaning |
| --- | --- | --- |
| `kind` | `OptionSourceKind` | classifies the source; `OptionSourceRemote` (`remote`) is the only value emitted |
| `request` | `*RemoteOptionRequest` | the operation returning the option records |
| `mapping` | `*RemoteOptionMapping` | how to read a label and a value out of each record; `nil` means the defaults |

The designer's own source union also has `static` and `ref`, but neither
survives projection, so `OptionSourceKind` is deliberately narrower than it.

`RemoteOptionRequest` addresses the operation in the framework's own
resource/action addressing: `resource`, `action`, `version` (empty means the
default), and `params`. `RemoteOptionMapping` names the record keys —
`labelKey` (default `label`), `valueKey` (default `value`), `disabledKey`, and
`descriptionKey`; an empty entry falls back to its default.

**The framework never resolves an option source.** Form-data validation reads
`options` only, and a `select` field with no options accepts any submitted
value — so an `optionSource` constrains nothing on the server. It exists for
consumers that must render a *stored* value as its label: a business list view
showing a stored select value needs the label, not the raw code, and the static
case was already covered by `options`.

Deploy validation covers option sources at the top level and inside a detail
table alike: an `optionSource` whose `kind` is outside the vocabulary, or a
remote one whose `request` is absent or names a blank `resource` / `action`,
fails the deploy. Neither would fail anywhere else — the framework never
resolves the source — so the version would deploy cleanly and then render as a
raw value in every consumer that replayed it.

Detail-table columns carry their own: `columns` entries are full
`FormFieldDefinition` values, so a `select` column inside a `table` field has
the same `options` / `optionSource` pair a top-level field has.

#### Remote Request Parameters

`RemoteOptionRequest.Params` is a `map[string]DynamicParam`, carried
**unevaluated**. Each parameter is one of two kinds:

| `DynamicParamKind` | Wire value | Payload | Meaning |
| --- | --- | --- | --- |
| `DynamicParamLiteral` | `literal` | `value` | a fixed value the designer typed, identical for every evaluation |
| `DynamicParamExpression` | `expression` | `source` | an expression the form runtime evaluates against the live form values before issuing the request — what makes a cascading select work |

The backend holds no form values to evaluate an expression against, so it
stores the source text and a consumer replaying the request owns that step.
`value` is deliberately not `omitempty`: a literal `false`, `0`, or `""` is a
value the designer chose, and dropping it would silently change the request.

**Ask `HasBoundParams()` before translating a column.**
`(*RemoteOptionRequest).HasBoundParams()` reports whether any parameter is an
expression. It is the question a consumer must ask before assuming one lookup
covers a whole column:

- **`false`** — the request resolves identically for every row, so **one** call
  builds a value-to-label map for the entire column.
- **`true`** — the option set is **per row** (its parameters depend on that
  row's own form values), so a list view must either evaluate and issue the
  request per row, or leave that column untranslated.

The method is nil-safe: a nil receiver reports `false`.

#### Translating Stored Values in a Host

Inject `approval.FormSchemaParser` — it resolves at root scope whenever
`vef.ApprovalModule` is enabled — and parse the version's form schema:

```go
package options

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/coldsmirk/vef-framework-go/approval"
)

// OptionLabelPlan says how a list view should turn each selection field's
// stored value into a display label.
type OptionLabelPlan struct {
    // Static maps a field key to its value-to-label table, built from the
    // options the definition already carries.
    Static map[string]map[string]string
    // Remote maps a field key to the request a consumer replays ONCE to build
    // that table itself. Only fields whose request has no bound parameters.
    Remote map[string]*approval.RemoteOptionRequest
    // PerRow lists the field keys whose option set depends on the row's own
    // form values, so one lookup cannot translate the whole column.
    PerRow []string
}

type OptionLabelService struct {
    parser approval.FormSchemaParser
}

func NewOptionLabelService(parser approval.FormSchemaParser) *OptionLabelService {
    return &OptionLabelService{parser: parser}
}

func (s *OptionLabelService) Plan(ctx context.Context, schema json.RawMessage) (*OptionLabelPlan, error) {
    fields, err := s.parser.ParseFormFields(ctx, schema)
    if err != nil {
        return nil, err
    }

    plan := &OptionLabelPlan{
        Static: make(map[string]map[string]string),
        Remote: make(map[string]*approval.RemoteOptionRequest),
    }

    for _, field := range fields {
        switch {
        case field.Options != nil:
            // Enumerated at deploy: translate locally, no lookup at all.
            labels := make(map[string]string, len(field.Options))
            for _, option := range field.Options {
                labels[fmt.Sprint(option.Value)] = option.Label
            }

            plan.Static[field.Key] = labels

        case field.OptionSource != nil:
            request := field.OptionSource.Request
            if request.HasBoundParams() {
                plan.PerRow = append(plan.PerRow, field.Key)

                continue
            }

            plan.Remote[field.Key] = request

        default:
            // Free-form field: the stored value is already what to display.
        }
    }

    return plan, nil
}

// labelKey and valueKey apply RemoteOptionMapping's defaults.
func labelKey(mapping *approval.RemoteOptionMapping) string {
    if mapping == nil || mapping.LabelKey == "" {
        return "label"
    }

    return mapping.LabelKey
}

func valueKey(mapping *approval.RemoteOptionMapping) string {
    if mapping == nil || mapping.ValueKey == "" {
        return "value"
    }

    return mapping.ValueKey
}
```

Values are compared by their string form (`fmt.Sprint`), which is what the
framework's own `select` validation does, so `1` and `"1"` address the same
option.

Parsing is not the only way in. The same flat list is persisted on
`apv_flow_version.form_fields` and returned as `formFields` by the flow-version
detail response, so a host holding that row reads `FlowVersion.FormFields`
directly. Parse when all it has is the schema document — the instance-detail
responses return the host document verbatim as `formSchema` and never the
derived fields.

`OptionSource.Request` can be relied on for a version that deployed: the
missing-operation check above rejects a remote source without one, so it is
never `nil` on fields that survived deploy.

#### Upload Fields

The built-in parser projects the designer's `upload` widget onto the `upload`
field kind.

An upload field's `columnType` is inferred from its **file count**, not from
`maxLength`: `maxCount > 1` yields `json` (the value is an array of storage
keys — the same shape `checkbox-group` carries), and anything else yields
`text` (a single storage key). `maxLength` bounds a *string's* length, while an
upload's bound is how many files it accepts, so reading it here would size the
column off the wrong quantity. An explicit `columnType` on the field still
wins, as it does for every kind.

## Flow Validation

### Initiator Rules

`isAllInitiationAllowed` and `initiators` are **mutually exclusive**: when
`isAllInitiationAllowed` is `true`, submitting initiator rules returns
`ErrInitiatorsNotAllowed` (code `40022`). When restricted, at least one
initiator rule with at least one ID is required — an empty rule set or a rule
whose `ids` is empty returns `ErrInitiatorsRequired` (code `40023`), because
a rule that selects nobody leaves the flow just as unstartable as no rule at
all.

### Business Binding Schema

Flows use `BindingMode` to decide how form/instance data is associated with
business data:

- `BindingStandalone` (wire value `standalone`) — form data lives in approval
  tables (the default).
- `BindingBusiness` (wire value `business`) — the flow is bound to an existing
  business row via `BusinessBindingConfig`.

Business-bound flows (`BindingMode = business`) are validated against the live
database schema at save time. When the configured `tableName` does not exist,
`ErrBindingTableMissing(tableName)` is returned (code `40018`,
`ErrCodeBindingSchemaInvalid`). When a configured column does not exist in the
table, `ErrBindingColumnMissing(columnName)` is returned (same code). Both
name the missing object so operators can fix the configuration without
inspecting the schema themselves.

## Flow Models and Designer Enums

Flow design and persistence models exposed by the public package include
`FlowCategory`, `Flow`, `FlowVersion`, `FlowNode`, `FlowEdge`, `FlowInitiator`,
`FlowNodeAssignee`, `FlowNodeCC`, `FormFieldDefinition`, `FieldOption`,
`FieldOptionSource`, `RemoteOptionRequest`, `RemoteOptionMapping`,
`DynamicParam`, `FormSnapshot`, `ActionLog`, `UserInfo`, and `UrgeRecord`
(there is no structured `FormDefinition` wrapper — the host document is opaque,
and `FormFieldDefinition` with the shapes it nests is the only framework form
shape). Flow-version
status uses `VersionStatus`: `VersionDraft` (`draft`), `VersionPublished`
(`published`), and `VersionArchived` (`archived`).

`Flow.Labels` is host-owned selection metadata — validated at save
time by the shared label rule and equality-filterable in list queries; see
[RPC Resources](./resources.md) for the wire shape.

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
| `OptionSourceKind` | `remote` |
| `DynamicParamKind` | `literal`, `expression` |

---

Next: [Instance Runtime](./runtime.md) for what happens after a designed flow starts running.
