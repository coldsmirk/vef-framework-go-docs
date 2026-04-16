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
auto_migrate = true
outbox_relay_interval = 5
outbox_max_retries = 10
outbox_batch_size = 100
```

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

## Event Outbox

The approval module uses the transactional outbox pattern for reliable event publishing. Events are written to `apv_event_outbox` within the same transaction as the approval action, then relayed asynchronously.

| Status | Constant |
| --- | --- |
| Pending | `EventOutboxPending` |
| Processing | `EventOutboxProcessing` |
| Completed | `EventOutboxCompleted` |
| Failed | `EventOutboxFailed` |

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
