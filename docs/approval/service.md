---
sidebar_position: 6
---

# Programmatic Flow Control

`approval.Service` is the approval engine's programmatic control surface: every
runtime operation the `approval/instance` and `approval/admin` resources
expose, callable from host code with no HTTP request in sight. Inject it
wherever `vef.ApprovalModule` is enabled.

```go
type OrderService struct {
    approval approval.Service
}
```

The API resources are themselves callers of this interface, so an operation
behaves **identically** whether a request or a host routine triggered it — the
same validation, the same domain events in the same order, the same audit rows,
lifecycle hooks and business write-back. There is no second code path to keep
in step.

Do not reach for `cqrs.Bus` plus the `internal/approval/command` types instead.
They are unreachable from another module anyway, and dispatched directly they
would bypass the collector behaviors that publish events and write action logs.

## The operations

| Method | Input | Notes |
| --- | --- | --- |
| `StartInstance` | `approval.StartInstanceInput` | Returns the `*Instance` after the engine has traversed the start node — a flow that reaches an end node without stopping comes back already completed. |
| `WithdrawInstance` | `approval.WithdrawInstanceInput` | Operator must be the applicant. |
| `ResubmitInstance` | `approval.ResubmitInstanceInput` | Operator must be the applicant; `FormData` replaces the instance's form data. |
| `TerminateInstance` | `approval.TerminateInstanceInput` | Administrative. Running, returned and withdrawn instances can all be terminated — the instance state machine is the single authority. |
| `ApproveTask` | `approval.ApproveTaskInput` | Approves an approval task **or** finishes a handle task; the two share one command. |
| `RejectTask` | `approval.RejectTaskInput` | Operator must hold the task. |
| `TransferTask` | `approval.TransferTaskInput` | 转办. The node must allow transfer. |
| `RollbackTask` | `approval.RollbackTaskInput` | 退回. `TargetNodeID` must be a configured rollback target when the rollback type is `specified`. |
| `ReassignTask` | `approval.ReassignTaskInput` | Administrative. The operator need not hold the task. |
| `AddAssignee` | `approval.AddAssigneeInput` | 加签. `parallel` is rejected on sequential nodes. |
| `RemoveAssignee` | `approval.RemoveAssigneeInput` | 减签. `TaskID` is the removed assignee's task, not the operator's; a node's last remaining assignee cannot be removed. |
| `AddCC` | `approval.AddCCInput` | 手动抄送. The current node must allow manual CC. |
| `MarkCCRead` | `approval.MarkCCReadInput` | Records a CC read receipt. |
| `UrgeTask` | `approval.UrgeTaskInput` | 催办, subject to the node's per-(task, urger) cooldown. |
| `RetryBusinessProjection` | `approval.RetryBusinessProjectionInput` | Administrative. Re-applies one eventual business projection immediately instead of waiting for the worker's next backoff window. |

Flow-definition management — create, deploy, publish, update, toggle — is
deliberately **not** part of this contract. It is design-time administration
with a different trust posture, served by the `approval/flow` resource.

## The handle is the transaction boundary

Every method takes the `orm.DB` handle the operation runs on, and that handle
decides the transaction:

- pass the handle of an **open `RunInTx` scope** and the operation joins that
  transaction, so a business write and the approval action it triggers commit
  or roll back together;
- pass a **plain handle** and the operation opens and commits a transaction of
  its own.

The handle is explicit rather than read off the context because
`orm.DB.RunInTx` does not attach the transaction handle to the context it hands
the callback — a method that only took a context would silently open a second,
independent transaction inside the caller's. A nil handle is rejected with
`approval.ErrDBRequired` rather than defaulted to the framework's own, and the
handle must belong to the **primary** data source: the approval tables live
there and the domain events publish with `event.WithTx` against it.

A `fiber.Ctx` is unwrapped before binding. That matters because
`contextx.SetDB` mutates a `fiber.Ctx`'s Locals **in place**: binding a
transaction onto one would leave the whole request pointing at a transaction
that dies at the caller's commit.

### Pass the context `RunInTx` hands the callback

Not the one you called `RunInTx` with. `RunInTx` puts one thing on that
context: the commit-hook registry `orm.OnCommit` reads, through which the
`tx_memory` event transport defers delivery. Under a route that uses
`tx_memory`, the outer context therefore fails the publish with
`orm.ErrNoCommitScope` and rolls the whole action back.

```go
db.RunInTx(ctx, func(txCtx context.Context, tx orm.DB) error {
    if err := writeBusinessRow(txCtx, tx); err != nil {
        return err
    }

    return svc.ApproveTask(txCtx, tx, approval.ApproveTaskInput{
        TaskID:   taskID,
        Operator: operator,
        Opinion:  "approved by the settlement routine",
        Caller:   approval.SystemCaller,
    })
})
```

### Binding the audit operator

Audit columns (`created_by` / `updated_by`) render the operator bound on the
handle, and `orm.DB.WithNamedArg` is pool-scoped only — it **panics** on a
transaction handle. Bind the person before opening the transaction and pass the
tx that handle yields:

```go
acting := db.WithNamedArg(orm.PlaceholderKeyOperator, userID)
acting.RunInTx(ctx, func(txCtx context.Context, tx orm.DB) error { … })
```

A handle carrying no operator — the injected primary `orm.DB` — writes
`orm.OperatorSystem` into those columns. The action log records the input's
`Operator` either way, so an unbound handle costs row-level attribution, not
the audit trail.

## The Service does not authorize

:::caution[Authorization is the caller's]
The RBAC tokens that gate the HTTP endpoints (`approval.instance.terminate`,
`approval.task.reassign`, `approval.binding.retry`) live on the API
**operations**, not in the commands behind them. `TerminateInstance`,
`ReassignTask` and `RetryBusinessProjection` enforce only the tenant scope
carried by `Caller`. A host exposing any of them through its own surface owns
that check — the call site looks identical to the framework's own, so nothing
flags the omission.
:::

Every input carries a `Caller` (`approval.CallerContext`), and it is required:
a zero value fails closed — `StartInstanceInput` reports `ErrFlowNotFound`
rather than a permission error, so a caller cannot probe another tenant's
flows. Use `approval.SystemCaller` for trusted automation, or
`approval.CallerContext{TenantID: …}` to pin the operation to one tenant.

`StartInstanceInput.Globals` deserves the same care. On the request path the
framework resolves globals server-side through `InstanceGlobalsResolver`
precisely because they steer condition branches; here the host passes them
directly, so supply values you trust rather than forwarding anything a client
sent.

## Errors

Errors are the public sentinels in the `approval` package
(`ErrFlowNotActive`, `ErrNotAllowedInitiate`, `ErrTaskNotPending`, …),
matchable with `errors.Is` — `result.Error` compares by code, so a sentinel
matches regardless of the rendered message language.

```go
if errors.Is(err, approval.ErrTaskNotPending) {
    // someone else already acted on it
}
```

## Re-entrancy

`approval.Service` is the first host-callable way to re-enter the command
pipeline from **inside** a running command: an `InstanceLifecycleHook` fires
in-transaction and may itself act on another task. The CQRS collectors handle
that — an inner dispatch reuses the collector already on the context and
declines ownership, so the outermost dispatch still owns the flush and events
stay in occurrence order. A failed inner dispatch unwinds only what it
appended, so an enclosing command that swallows the inner error keeps its own
items and drops the ones describing writes that did not happen.

---

Next: [Events & Integration](./integration.md) for the extension points these
operations fire, or [RPC Resources](./resources.md) for the HTTP surface that
calls the same interface.
