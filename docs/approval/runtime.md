---
sidebar_position: 4
---

# Instance Runtime

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
| `returned` | `terminated` |
| `returned` | `withdrawn` |
| `withdrawn` | `running` |
| `withdrawn` | `terminated` |

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
| Table | `StorageTable` | `table` | a generated physical table per published version, with `apv_instance.form_data` still populated as the canonical JSON snapshot |

`StorageMode.IsValid()` accepts both exported modes. Table mode records generated DDL metadata through `FormTable` and `FormTableColumn`.

### Table Storage Metadata

When a published version uses `StorageTable`, the framework exposes two public metadata models:

| Model | Purpose | Key JSON fields |
| --- | --- | --- |
| `FormTable` | one generated physical table (main projection or detail-table child) | `flowId`, `versionId`, `physicalTableName`, `sourceFieldKey` |
| `FormTableColumn` | one generated column per form field or built-in column | `formTableId`, `columnName`, `columnType`, `isNullable`, `sourceFieldKey`, `sortOrder` |

`FormTable` is the single source of truth for the DDL the framework generated:
the engine consults it for idempotency before creating a table, and operators
can map a version to its projection tables through it. `FormTable.SourceFieldKey`
is `""` for the version's main projection table, or the owning table field's
key for a detail-table child projection; `(versionId, sourceFieldKey)` is
unique. `ColumnDataType` is the logical field-to-column vocabulary used by form
definitions before the storage layer maps it to dialect-specific SQL types.

### Generated Table Layout

Publishing a table-mode version provisions one main projection table plus one
child table per detail-table (`table` kind) field:

| Table | Physical name | Built-in columns | Field columns |
| --- | --- | --- | --- |
| main projection | `apv_form_<code>_<versionId>` (sanitized flow code, truncated to the 63-char identifier cap; falls back to `apv_form_<versionId>`) | `id` (PK), `instance_id` (UNIQUE), `created_at` last | one column per scalar field, in declared order |
| detail-table child | `apv_form_<versionId>__<fieldKey suffix>` | `id` (PK), `instance_id` (indexed, NOT unique), `row_index`, `created_at` last | one column per table-field column, in declared order |

Every physical table and column name is validated as a safe SQL identifier
before it reaches a DDL/DML string, field keys must not collide with the
built-in column names, and every form value is bound as a `?` argument — never
interpolated. `CREATE TABLE IF NOT EXISTS` runs outside the publish transaction
(DDL implicitly commits the in-flight transaction on MySQL) and is idempotent;
the `FormTable` / `FormTableColumn` metadata is recorded inside the publish
transaction so it commits or rolls back with the version's published state.

At write time the projection is replace-never-append: the main table holds
exactly one row per instance (backed by the `instance_id` UNIQUE constraint),
a child table holds one row per detail line ordered by `row_index`, and each
projection (at start and at every resubmit) deletes the instance's existing
rows before inserting fresh ones — the tables reflect current form data, never
an accumulating history.

## Instance Progress Projections

Admin and user instance-detail responses expose two read-only projections beside the instance snapshot:

| Projection | Public types | Purpose |
| --- | --- | --- |
| timeline | `TimelineEntryKind`, `TimelineEntry`, `NodeVisitStatus`, `NodeParticipant`, `Activity`, `ActivityUrge`, `CCRecipient` | chronological account of the path the instance actually took |
| flow graph | `NodeProgressStatus`, `InstanceFlowGraph`, `FlowGraphNode`, `FlowGraphNodeData`, `FlowGraphEdge` | React Flow-compatible graph annotated with runtime progress |

`FlowGraphNode.ID` is the React Flow design-time node id. `FlowGraphNode.NodeID` is the persistent flow-node id used by action logs and rollback targets.

Progress and timeline enums are exported explicitly:

| Enum | Constants |
| --- | --- |
| `TimelineEntryKind` | `TimelineEntryStart`, `TimelineEntryApproval`, `TimelineEntryHandle`, `TimelineEntryCC`, `TimelineEntryWithdraw`, `TimelineEntryTerminate` |
| `NodeVisitStatus` | `NodeVisitActive`, `NodeVisitPassed`, `NodeVisitRejected`, `NodeVisitReturned`, `NodeVisitCanceled` |
| `NodeProgressStatus` | `NodeProgressPending`, `NodeProgressActive`, `NodeProgressPassed`, `NodeProgressRejected`, `NodeProgressReturned`, `NodeProgressCanceled` |

The persisted visit model behind these projections is `NodeVisit`.

## Instance Number Generation

Implement the `InstanceNoGenerator` interface to customize instance numbering:

```go
type InstanceNoGenerator interface {
    Generate(ctx context.Context, flowCode string) (string, error)
}
```

---

Next: [Events & Integration](./integration.md) for reacting to lifecycle transitions from host code.
