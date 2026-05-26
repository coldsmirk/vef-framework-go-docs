---
sidebar_position: 3
---

# Approval Module

The `approval` module provides a complete workflow engine for building approval-based business processes. It supports visual flow design (React Flow compatible), multi-level approval chains, conditional branching, parallel approval, delegation, rollbacks, and transactional event publishing.

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

> The outbox-related fields that previously lived under `[vef.approval]` (`outbox_relay_interval`, `outbox_max_retries`, `outbox_batch_size`) moved to `[vef.event.transports.outbox]` in v0.21. The approval module now publishes through the framework-wide outbox transport — see [Event Bus](../features/event-bus). Approval's binding listener and outbox publisher both assert routing at boot via `event.RouteInspector`, so a misconfigured route fails the application instead of degrading silently.

See [Configuration Reference](../reference/configuration-reference) for details.

## Binding Modes

| Mode | Constant | Description |
| --- | --- | --- |
| Standalone | `BindingStandalone` | Form data stored in the approval module's own tables |
| Business | `BindingBusiness` | Links to an existing business data table |

Business binding connects the approval flow to your domain tables via `BusinessTable`, `BusinessPkField`, `BusinessTitleField`, and `BusinessStatusField`.

## Node Types

| Node Kind | Constant | Description |
| --- | --- | --- |
| Start | `NodeStart` | Entry point of the workflow |
| Approval | `NodeApproval` | Requires approval action from assignees |
| Handle | `NodeHandle` | Requires processing/handling action |
| Condition | `NodeCondition` | Branches based on conditions |
| CC | `NodeCC` | Sends notifications to specified users |
| End | `NodeEnd` | Terminal point of the workflow |

## Approval Methods

When a node has multiple assignees:

| Method | Constant | Behavior |
| --- | --- | --- |
| Sequential | `ApprovalSequential` | Approvers process one by one in order |
| Parallel | `ApprovalParallel` | Approvers process simultaneously |

### Pass Rules (for Parallel)

| Rule | Constant | Behavior |
| --- | --- | --- |
| All | `PassAll` | All assignees must approve |
| Any | `PassAny` | At least one approval passes |
| Ratio | `PassRatio` | A percentage must approve |
| Any Reject | `PassAnyReject` | Any rejection fails the node |

## Assignee Types

| Kind | Constant | Description |
| --- | --- | --- |
| User | `AssigneeUser` | Specific users |
| Role | `AssigneeRole` | Users with a role |
| Department | `AssigneeDepartment` | Department head |
| Self | `AssigneeSelf` | The applicant |
| Superior | `AssigneeSuperior` | Direct superior |
| Dept Leader | `AssigneeDepartmentLeader` | Multi-level supervisor chain |
| Form Field | `AssigneeFormField` | Determined by a form field value |

## Instance Lifecycle

```
submit → Running → approve/reject → Approved/Rejected
                 → withdraw       → Withdrawn
                 → rollback       → Returned
                 → terminate      → Terminated
                 → resubmit       → Running (again)
```

### Instance Statuses

| Status | Constant | Final? |
| --- | --- | --- |
| Running | `InstanceRunning` | No |
| Approved | `InstanceApproved` | Yes |
| Rejected | `InstanceRejected` | Yes |
| Withdrawn | `InstanceWithdrawn` | No |
| Returned | `InstanceReturned` | No |
| Terminated | `InstanceTerminated` | Yes |

### Task Statuses

| Status | Constant | Final? |
| --- | --- | --- |
| Waiting | `TaskWaiting` | No |
| Pending | `TaskPending` | No |
| Approved | `TaskApproved` | Yes |
| Rejected | `TaskRejected` | Yes |
| Handled | `TaskHandled` | Yes |
| Transferred | `TaskTransferred` | Yes |
| Rolled Back | `TaskRolledBack` | Yes |
| Canceled | `TaskCanceled` | Yes |
| Removed | `TaskRemoved` | Yes |
| Skipped | `TaskSkipped` | Yes |

## Actions

| Action | Constant | Description |
| --- | --- | --- |
| Submit | `ActionSubmit` | Start a new instance |
| Approve | `ActionApprove` | Approve a task |
| Handle | `ActionHandle` | Complete a handle task |
| Reject | `ActionReject` | Reject a task |
| Transfer | `ActionTransfer` | Transfer to another user |
| Withdraw | `ActionWithdraw` | Applicant withdraws |
| Cancel | `ActionCancel` | Cancel a task |
| Rollback | `ActionRollback` | Roll back to a previous node |
| Add Assignee | `ActionAddAssignee` | Dynamically add an assignee |
| Remove Assignee | `ActionRemoveAssignee` | Remove an assignee |
| Resubmit | `ActionResubmit` | Resubmit a returned instance |
| Reassign | `ActionReassign` | Admin reassigns a task |
| Terminate | `ActionTerminate` | Admin force-terminates |

## Rollback Configuration

| Property | Options |
| --- | --- |
| `RollbackType` | `none`, `previous`, `start`, `any`, `specified` |
| `RollbackDataStrategy` | `clear` (reset form), `keep` (preserve data) |

## Empty Assignee Handling

When no assignee is found for a node:

| Action | Constant |
| --- | --- |
| Auto-pass | `EmptyAssigneeAutoPass` |
| Transfer to admin | `EmptyAssigneeTransferAdmin` |
| Transfer to superior | `EmptyAssigneeTransferSuperior` |
| Transfer to applicant | `EmptyAssigneeTransferApplicant` |
| Transfer to specified | `EmptyAssigneeTransferSpecified` |

## Timeout Handling

| Action | Constant | Behavior |
| --- | --- | --- |
| None | `TimeoutActionNone` | Mark timeout only |
| Auto Pass | `TimeoutActionAutoPass` | Automatically approve |
| Auto Reject | `TimeoutActionAutoReject` | Automatically reject |
| Notify | `TimeoutActionNotify` | Send notification only |
| Transfer Admin | `TimeoutActionTransferAdmin` | Transfer to node admin |

## Form Data Storage

| Mode | Constant | Location |
| --- | --- | --- |
| JSON | `StorageJSON` | `apv_instance.form_data` (JSONB column) |
| Table | `StorageTable` | Dynamic table `apv_form_data_{flow_code}` |

## Event Publication

The approval module publishes its domain events through the framework's transactional outbox transport (see [Event Bus](../features/event-bus)). Every approval command writes the event record in the same transaction as the business mutation; the outbox relay then forwards them to the configured sink.

Subscribers must:

1. Attach with `event.WithGroup("...")` because the route resolves to an at-least-once transport.
2. Rely on the Inbox middleware for dedupe (it activates automatically when `event.middleware.inbox = true` and a transport advertises `AtLeastOnce`).

> The standalone `apv_event_outbox` table and module-private `EventOutboxStatus` constants from earlier snapshots have been retired. Approval no longer carries a private outbox — it composes with the framework one.

### Domain Event Types

All approval events implement `event.Event` and embed enough payload to drive integrations without a follow-up read.

Instance lifecycle:

| Type constant | When |
| --- | --- |
| `approval.instance.created` (`InstanceCreatedEvent`) | a new instance was started |
| `approval.instance.completed` (`InstanceCompletedEvent`) | instance reached a terminal status |
| `approval.instance.withdrawn` (`InstanceWithdrawnEvent`) | applicant withdrew the instance |
| `approval.instance.rolled_back` (`InstanceRolledBackEvent`) | instance was rolled back to a previous node |
| `approval.instance.returned` (`InstanceReturnedEvent`) | instance was returned to applicant |
| `approval.instance.resubmitted` (`InstanceResubmittedEvent`) | returned instance was resubmitted |
| `approval.instance.binding_failed` (`InstanceBindingFailedEvent`) | the binding listener could not write back the final status to the business row |

Node lifecycle:

| Type constant | When |
| --- | --- |
| `approval.node.entered` | engine activated a node |
| `approval.node.auto_passed` | a node auto-passed because no assignees were found |

Task lifecycle:

| Type constant | When |
| --- | --- |
| `approval.task.created` | a task was created (emitted on every task-creation path since v0.25) |
| `approval.task.approved` | task approved |
| `approval.task.handled` | handle task completed |
| `approval.task.rejected` | task rejected |
| `approval.task.transferred` | task transferred to another user |
| `approval.task.reassigned` | admin reassigned a task |
| `approval.task.timed_out` | timeout scanner fired the configured timeout action |
| `approval.task.assignees_added` | dynamic assignees added |
| `approval.task.assignees_removed` | dynamic assignees removed |
| `approval.task.deadline_warning` | pre-warning scanner flagged an approaching deadline |
| `approval.task.urged` | applicant urged an assignee |

CC + Flow:

| Type constant | When |
| --- | --- |
| `approval.cc.notified` | a CC node delivered notifications |
| `approval.flow.created` / `updated` / `deployed` / `toggled` / `published` | flow design lifecycle changes |

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

## Instance Number Generation

Implement the `InstanceNoGenerator` interface to customize instance numbering:

```go
type InstanceNoGenerator interface {
    Generate(ctx context.Context, flowCode string) (string, error)
}
```
