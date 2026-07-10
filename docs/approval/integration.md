---
sidebar_position: 5
---

# Events & Integration

## Event Publication

The approval module publishes its domain events through the framework's transactional outbox transport (see [Event Bus](../infrastructure/event-bus.md)). Every approval command writes the event record in the same transaction as the business mutation; the outbox relay then forwards them to the configured sink.

Subscribers must:

1. Attach with `event.WithGroup("...")` because the route resolves to an at-least-once transport.
2. Rely on the Inbox middleware for dedupe (it activates automatically when `event.middleware.inbox = true` and a transport advertises `AtLeastOnce`).

> The standalone `apv_event_outbox` table and module-private `EventOutboxStatus` constants from earlier snapshots have been retired. Approval no longer carries a private outbox — it composes with the framework one.

### Domain Event Types

All approval events implement `DomainEvent`; `AllEventTypes()` is the exhaustive exported registry of their topic strings. Instance and task events embed `InstanceEventBase` (`instanceId`, `instanceNo`, `tenantId`, `title`, `flowId`, `flowCode`, optional `businessRef`, `applicant`, `occurredTime`). Task events additionally embed `TaskEventBase` (`taskId`, `nodeId`, `nodeName`). Flow events embed `FlowEventBase` (`flowId`, `tenantId`, `code`, `name`, `occurredTime`). The tables below list fields beyond those bases.

Instance lifecycle:

| Type constant | Topic | Payload / constructor | Payload fields beyond common fields | When |
| --- | --- | --- | --- | --- |
| `EventTypeInstanceCreated` | `approval.instance.created` | `InstanceCreatedEvent`, `NewInstanceCreatedEvent` | none | a new instance was started |
| `EventTypeInstanceCompleted` | `approval.instance.completed` | `InstanceCompletedEvent`, `NewInstanceCompletedEvent` | `finalStatus`, `finishedAt`, `reason` (set only when `finalStatus` is terminated) | instance reached a terminal status |
| `EventTypeInstanceWithdrawn` | `approval.instance.withdrawn` | `InstanceWithdrawnEvent`, `NewInstanceWithdrawnEvent` | `operator` (`UserInfo`), `reason` | applicant withdrew the instance |
| `EventTypeInstanceRolledBack` | `approval.instance.rolled_back` | `InstanceRolledBackEvent`, `NewInstanceRolledBackEvent` | `fromNodeId`, `fromNodeName`, `toNodeId`, `toNodeName`, `operator` (`UserInfo`), `opinion` | instance was rolled back to a previous node |
| `EventTypeInstanceReturned` | `approval.instance.returned` | `InstanceReturnedEvent`, `NewInstanceReturnedEvent` | `fromNodeId`, `fromNodeName`, `toNodeId`, `toNodeName`, `operator` (`UserInfo`), `opinion` | instance was returned to applicant |
| `EventTypeInstanceResubmitted` | `approval.instance.resubmitted` | `InstanceResubmittedEvent`, `NewInstanceResubmittedEvent` | `operator` (`UserInfo`) | returned or withdrawn instance was resubmitted |
| `EventTypeInstanceBindingFailed` | `approval.instance.binding_failed` | `InstanceBindingFailedEvent`, `NewInstanceBindingFailedEvent` | `trigger`, `status`, `businessTable`, `errorMessage` | the engine-owned write-back could not persist `status` to the business row after `trigger` |

Node lifecycle:

| Type constant | Topic | Payload / constructor | Payload fields beyond common fields | When |
| --- | --- | --- | --- | --- |
| `EventTypeNodeAutoPassed` | `approval.node.auto_passed` | `NodeAutoPassedEvent`, `NewNodeAutoPassedEvent` | `nodeId`, `nodeName`, `reason` | a node auto-passed because of auto-pass execution, empty-assignee handling, or same-applicant handling |

Task lifecycle:

| Type constant | Topic | Payload / constructor | Payload fields beyond common fields | When |
| --- | --- | --- | --- | --- |
| `EventTypeTaskCreated` | `approval.task.created` | `TaskCreatedEvent`, `NewTaskCreatedEvent` | `assignee` (`UserInfo`), `deadline` | a task was created; sequential follow-up tasks may start with `deadline` omitted while waiting |
| `EventTypeTaskApproved` | `approval.task.approved` | `TaskApprovedEvent`, `NewTaskApprovedEvent` | `operator` (`UserInfo`), `opinion` | task approved |
| `EventTypeTaskHandled` | `approval.task.handled` | `TaskHandledEvent`, `NewTaskHandledEvent` | `operator` (`UserInfo`), `opinion` | handle task completed |
| `EventTypeTaskRejected` | `approval.task.rejected` | `TaskRejectedEvent`, `NewTaskRejectedEvent` | `operator` (`UserInfo`), `opinion` | task rejected |
| `EventTypeTaskCanceled` | `approval.task.canceled` | `TaskCanceledEvent`, `NewTaskCanceledEvent` | `assignee` (`UserInfo`), `reason` | engine canceled a task that no longer needs a decision |
| `EventTypeTaskTransferred` | `approval.task.transferred` | `TaskTransferredEvent`, `NewTaskTransferredEvent` | `from` (`UserInfo`), `to` (`UserInfo`), `reason` | task transferred to another user |
| `EventTypeTaskReassigned` | `approval.task.reassigned` | `TaskReassignedEvent`, `NewTaskReassignedEvent` | `from` (`UserInfo`), `to` (`UserInfo`), `reason` | admin reassigned a task |
| `EventTypeTaskTimedOut` | `approval.task.timed_out` | `TaskTimedOutEvent`, `NewTaskTimedOutEvent` | `assignee` (`UserInfo`), `deadline` | timeout scanner fired the configured timeout action |
| `EventTypeAssigneesAdded` | `approval.task.assignees_added` | `AssigneesAddedEvent`, `NewAssigneesAddedEvent` | `addType`, `assignees` (`[]UserInfo`) | dynamic assignees added |
| `EventTypeAssigneesRemoved` | `approval.task.assignees_removed` | `AssigneesRemovedEvent`, `NewAssigneesRemovedEvent` | `assignees` (`[]UserInfo`) | dynamic assignees removed |
| `EventTypeTaskDeadlineWarning` | `approval.task.deadline_warning` | `TaskDeadlineWarningEvent`, `NewTaskDeadlineWarningEvent` | `assignee` (`UserInfo`), `deadline`, `hoursLeft` | pre-warning scanner flagged an approaching deadline |
| `EventTypeTaskUrged` | `approval.task.urged` | `TaskUrgedEvent`, `NewTaskUrgedEvent` | `urger` (`UserInfo`), `target` (`UserInfo`), `message` | applicant urged an assignee |

CC + Flow:

| Type constant | Topic | Payload / constructor | Payload fields beyond common fields | When |
| --- | --- | --- | --- | --- |
| `EventTypeCCNotified` | `approval.cc.notified` | `CCNotifiedEvent`, `NewCCNotifiedEvent` | `nodeId`, `nodeName`, `recipients` (`[]UserInfo`), `isManual` | a CC node or manual CC action delivered notifications |
| `EventTypeFlowCreated` | `approval.flow.created` | `FlowCreatedEvent`, `NewFlowCreatedEvent` | `categoryId` | flow created |
| `EventTypeFlowUpdated` | `approval.flow.updated` | `FlowUpdatedEvent`, `NewFlowUpdatedEvent` | none | flow updated |
| `EventTypeFlowDeployed` | `approval.flow.deployed` | `FlowDeployedEvent`, `NewFlowDeployedEvent` | `versionId`, `version` | flow version deployed |
| `EventTypeFlowToggled` | `approval.flow.toggled` | `FlowToggledEvent`, `NewFlowToggledEvent` | `isActive` | flow active flag changed |
| `EventTypeFlowPublished` | `approval.flow.published` | `FlowPublishedEvent`, `NewFlowPublishedEvent` | `versionId` | flow version published |

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
| `approval.BusinessRefProvider` | synchronous, inside the start-instance transaction | returning an error rolls back instance creation |
| `approval.BusinessRefResolver` | asynchronous write-back path (`completed` / `returned` / `withdrawn` / `resubmitted` triggers) | returning an error fails the engine-owned write-back, publishes `InstanceBindingFailedEvent`, and is retried by the outbox path |
| `approval.SubscribeInstance` (`event.SubscribeTyped` wrapper) | asynchronous, after the transaction commits | the bus retries via the outbox relay; consumers must be idempotent |

### `InstanceLifecycleHook`

```go
type InstanceLifecycleHook interface {
    OnInstanceCreated(ctx context.Context, db orm.DB, instance *Instance) error
    OnInstanceCompleted(ctx context.Context, db orm.DB, instance *Instance, finalStatus InstanceStatus) error
}
```

Use lifecycle hooks for invariants that must hold inside the transaction (e.g. allocating a tightly-coupled business row). Use event subscriptions for everything else. Register hooks with `vef.ProvideApprovalLifecycleHook(constructor)` — the constructor must return `approval.InstanceLifecycleHook`, and multiple hooks compose via the `vef:approval:lifecycle_hooks` group. `approval.NewFilteredLifecycleHook(hook, filters...)` wraps a hook with the same `InstanceFilter` vocabulary used by `SubscribeInstance` (below), so a hook can be scoped to specific flow codes or tenants without hand-written predicates.

### `BusinessRefProvider` and `BusinessRefResolver`

```go
type BusinessRefProvider interface {
    OnInstanceCreated(ctx context.Context, db orm.DB, flow *Flow, instance *Instance) (businessRef string, err error)
}

type BusinessRefResolver interface {
    ResolveRecordID(ctx context.Context, flow *Flow, businessRef string) (string, error)
}
```

`BusinessRefProvider` supplies or allocates the opaque `Instance.BusinessRef` when `Flow.BindingMode == BindingBusiness`; inject it with `vef.SupplyBusinessRefProvider`. The engine never parses `BusinessRef` directly. During write-back, it asks `BusinessRefResolver` (default: identity) to extract the record id matched against `Flow.BusinessPKField`; inject custom composite-ref handling with `vef.SupplyBusinessRefResolver`.

### Business Write-Back Linkage Matrix

A business-bound flow configures `Flow.BusinessTable` / `BusinessPKField` / `BusinessStatusField` (required) plus three optional columns — `BusinessInstanceIDField`, `BusinessStartedAtField`, `BusinessFinishedAtField` (Go names; JSON keeps `businessPkField` etc.). `nil` means "never touch that column". The engine-owned write-back (`binding.Writer.WriteBack`) projects the instance's current state per `approval.BindingTrigger` (v0.37):

| `BindingTrigger` | status column | instance-id column | started-at column | finished-at column |
| --- | --- | --- | --- | --- |
| `BindingTriggerStarted` (`started`) | `running` | the instance's id | now | cleared (`NULL`) |
| `BindingTriggerCompleted` (`completed`) | final status | — | — | `Instance.FinishedAt` |
| `BindingTriggerReturned` (`returned`) | `returned` | — | — | — |
| `BindingTriggerWithdrawn` (`withdrawn`) | `withdrawn` | — | — | — |
| `BindingTriggerResubmitted` (`resubmitted`) | `running` | — | — | cleared (`NULL`) |

The `started` projection runs synchronously inside the `start_instance` transaction — a failure rolls back the whole initiation. The other four triggers run asynchronously through the binding listener; a failure publishes `InstanceBindingFailedEvent` (carrying `trigger` and `status`) and is retried via the outbox path.

The approval engine owns the status write-back (`UPDATE businessTable SET businessStatusField = ? WHERE businessPkField = ?`, plus the optional columns above). Hosts can add lifecycle hooks or event subscribers around that behavior, but the write-back itself is configuration-driven and not replaced by an extension hook.

### Host Subscriptions: `SubscribeInstance`

```go
func SubscribeInstance[T InstanceEvent](
    bus event.Bus,
    handler func(ctx context.Context, evt T) error,
    opts ...InstanceSubscribeOption,
) (event.Unsubscribe, error)
```

`SubscribeInstance` (v0.37) is the declarative wrapper over `event.SubscribeTyped` for instance events (any event type embedding `InstanceEventBase`, e.g. `InstanceCompletedEvent`). Routing filters are data, not predicates: `InstanceSubscribeOption` accepts `approval.ForFlows(codes...)` / `approval.ForTenants(ids...)` (`InstanceFilter` values), each OR-ing within its own dimension and AND-ing across filters passed to the same call; events that fail a filter are acknowledged without invoking the handler. Business predicates (final status, form values) belong in the handler body, not in filters.

```go
vef.Invoke(func(bus event.Bus) error {
    _, err := approval.SubscribeInstance(bus,
        func(ctx context.Context, evt *approval.InstanceCompletedEvent) error {
            // idempotent side effect; the outbox relay redelivers on failure
            return nil
        },
        approval.ForFlows("leave_request"),
        approval.WithGroup("app:leave-writeback"),
    )
    return err
})
```

The consumer group defaults to one derived from the handler's method identity (`vef:sub:<module-relative pkg>.<Type>.<method>`, mirroring the `vef:default:<uuid>` shape used elsewhere). Anonymous handlers fail with `ErrAnonymousSubscriberGroup` unless `approval.WithGroup(name)` is given; a duplicate derived group registered twice in one process fails with `ErrDerivedGroupConflict`. **Renaming or moving a handler changes its derived group** (a new consumer group starts and the old one is orphaned) — pin the group with `WithGroup` before such refactors. `approval.WithConcurrency(n)` sets the per-subscription worker count.

## Business Identifier Validation

`Flow.BindingMode == BindingBusiness` flows carry SQL identifiers (`BusinessTable`, `BusinessPKField`, `BusinessStatusField`, and the three optional linkage columns above) that the engine-owned write-back interpolates directly into an `UPDATE` template. To prevent SQL injection, the framework whitelists identifiers against `^[A-Za-z_][A-Za-z0-9_]{0,62}$`:

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

## Supporting Public API Map

| Area | Public API |
| --- | --- |
| caller safety | `CallerContext`, `SystemCaller`, `IsSuperAdmin`, `SuperAdminRole`, `ErrCrossTenantAccess` |
| form data | `FormData`, `NewFormData`, `FormDefinition`, `FormFieldDefinition`, `FormSnapshot`, `ValidationRule`, `StorageMode`, `StorageJSON`, `StorageTable`, `FieldKind`, `FieldInput`, `FieldNumber`, `FieldDate`, `FieldTextarea`, `FieldSelect`, `FieldUpload`, `FieldTable`, `FieldOption`, `ColumnDataType`, `ColumnString`, `ColumnText`, `ColumnInteger`, `ColumnDecimal`, `ColumnBoolean`, `ColumnDate`, `ColumnDatetime`, `ColumnJSON` |
| table storage | `FormTable`, `FormTableColumn` (`FormTable.SourceFieldKey` is `""` for the main table, or the owning table field's key for a detail-table child projection) |
| flow models | `FlowCategory`, `Flow`, `FlowVersion`, `FlowNode`, `FlowEdge`, `FlowInitiator`, `FlowNodeAssignee`, `FlowNodeCC`, `VersionStatus`, `VersionDraft`, `VersionPublished`, `VersionArchived`, `ActionLog`, `OperatorInfo`, `UrgeRecord`, `DefaultTenantID` |
| node design | `FlowDefinition`, `NodeDefinition`, `EdgeDefinition`, `Position`, `NodeData`, `BaseNodeData`, `StartNodeData`, `ApprovalNodeData`, `HandleNodeData`, `ConditionNodeData`, `CCNodeData`, `EndNodeData`, `ErrUnknownNodeKind`, `ErrNodeDataUnmarshal` |
| conditions | `ConditionKind`, `ConditionField`, `ConditionExpression`, `Condition`, `ConditionGroup`, `ConditionBranch`, `EvaluationContext`, `ConditionEvaluator`, `InstanceGlobalsResolver`, `AggregateKind`, `AggregateSum`, `AggregateCount`, `AggregateAvg`, `Aggregator` |
| initiators and assignees | `InitiatorKind`, `InitiatorUser`, `InitiatorRole`, `InitiatorDepartment`, `AssigneeKind`, `AssigneeDefinition`, `AssigneeService`, `ResolvedAssignee`, `UserInfo`, `UserInfoResolver`, `RoleMembershipChecker`, `AddAssigneeType`, `AddAssigneeBefore`, `AddAssigneeAfter`, `AddAssigneeParallel` |
| CC | `CCKind`, `CCUser`, `CCRole`, `CCDepartment`, `CCFormField`, `CCTiming`, `CCTimingAlways`, `CCTimingOnApprove`, `CCTimingOnReject`, `CCDefinition`, `CCRecord` (`CCRecord.VisitID` scopes the record to one node traversal, mirroring `Task.VisitID`) |
| business write-back | `BindingTrigger`, `BindingTriggerStarted`, `BindingTriggerCompleted`, `BindingTriggerReturned`, `BindingTriggerWithdrawn`, `BindingTriggerResubmitted` |
| instance subscriptions | `SubscribeInstance`, `InstanceEvent`, `InstanceFilter`, `InstanceSubscribeOption`, `ForFlows`, `ForTenants`, `WithGroup`, `WithConcurrency`, `NewFilteredLifecycleHook`, `ErrAnonymousSubscriberGroup`, `ErrDerivedGroupConflict` |
| node behavior | `ApprovalMethod`, `TaskNodeData`, `ExecutionType`, `ExecutionManual`, `ExecutionAutoPass`, `ExecutionAutoReject`, `ConsecutiveApproverAction`, `ConsecutiveApproverNone`, `ConsecutiveApproverAutoPass`, `SameApplicantAction`, `SameApplicantSelfApprove`, `SameApplicantAutoPass`, `SameApplicantTransferSuperior`, `Permission`, `PermissionVisible`, `PermissionEditable`, `PermissionRequired`, `PermissionHidden`, `DefaultExecutionType`, `DefaultApprovalMethod`, `DefaultPassRule`, `DefaultEmptyAssigneeAction`, `DefaultSameApplicantAction`, `DefaultConsecutiveApproverAction`, `DefaultRollbackType`, `DefaultRollbackDataStrategy`, `DefaultTimeoutAction`, `DefaultCCTiming`, `DefaultHandleApprovalMethod`, `DefaultHandlePassRule`, `DefaultUrgeCooldownMinutes` |
| rollback and timeouts | `RollbackType`, `RollbackNone`, `RollbackPrevious`, `RollbackStart`, `RollbackAny`, `RollbackSpecified`, `RollbackDataStrategy`, `RollbackDataClear`, `RollbackDataKeep`, `EmptyAssigneeAction`, `EmptyAssigneeAutoPass`, `EmptyAssigneeTransferAdmin`, `EmptyAssigneeTransferSuperior`, `EmptyAssigneeTransferApplicant`, `EmptyAssigneeTransferSpecified`, `TimeoutAction`, `TimeoutActionNone`, `TimeoutActionAutoPass`, `TimeoutActionAutoReject`, `TimeoutActionNotify`, `TimeoutActionTransferAdmin` |
| action and status enums | `ActionType`, `InstanceStatus`, `TaskStatus`, `NodeKind`, `StorageMode`, `VersionStatus` |
| pass rules | `PassRule`, `PassRuleContext`, `PassRuleStrategy`, `PassRuleResult`, `PassRulePending`, `PassRulePassed`, `PassRuleRejected` |
| progress views | `TimelineEntryKind`, `TimelineEntry`, `NodeVisitStatus`, `NodeProgressStatus`, `InstanceFlowGraph`, `FlowGraphNode`, `FlowGraphNodeData`, `FlowGraphEdge`, `NodeParticipant`, `Activity`, `ActivityUrge`, `CCRecipient` |
| events | all `New...Event` constructors, `DomainEvent`, `InstanceEventBase`, `TaskEventBase`, `FlowEventBase`, `NewInstanceEventBase`, `NewTaskEventBase`, `NewFlowEventBase`, `PayloadOccurredAt`, `AllEventTypes`, and the `EventType...` constants |
| extension interfaces | `InstanceLifecycleHook`, `BusinessRefProvider`, `BusinessRefResolver`, `InstanceNoGenerator`, `ConditionEvaluator`, `InstanceGlobalsResolver`, `PrincipalTenantResolver`, `PrincipalDepartmentResolver`, `RoleMembershipChecker` |
| DI helpers (package `vef`) | `SupplyBusinessRefProvider`, `SupplyBusinessRefResolver`, `ProvideApprovalLifecycleHook`, `ProvideApprovalAggregator` |
| admin DTOs | package `approval/admin`: `Instance`, `InstanceDetail`, `InstanceDetailInfo`, `Task`, `ActionLog`, `Metrics` |
| user DTOs | package `approval/my`: `PendingTask`, `CompletedTask`, `CCRecord`, `InitiatedInstance`, `AvailableFlow`, `InstanceDetail`, `InstanceInfo`, `PendingCounts` |

---

Next: back to the [Overview](./overview.md) for wiring and configuration, or [RPC Resources](./resources.md) for the API surface these events originate from.
