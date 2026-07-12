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
| `EventTypeInstanceBindingFailed` | `approval.instance.binding_failed` | `InstanceBindingFailedEvent`, `NewInstanceBindingFailedEvent` | `trigger`, `status`, `businessTable`, `errorMessage` | an **eventual** business projection attempt failed after the desired state committed; the durable worker keeps retrying — this event is an operator notification, not the retry mechanism. Synchronous failures roll back the approval transaction and never emit it |

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
| `approval.InstanceLifecycleHook` (FX group `vef:approval:lifecycle_hooks`) | synchronous, inside the engine transaction — instance creation and every status transition | returning an error rolls back the surrounding command |
| `approval.BusinessRefProvider` | synchronous, inside the start-instance transaction | returning an error rolls back instance creation |
| `approval.BusinessRefResolver` | synchronous, at claim time inside the start-instance transaction | returning an error rolls back instance creation |
| `approval.SubscribeInstance` (`event.SubscribeTyped` wrapper) | asynchronous, after the transaction commits | the bus retries via the outbox relay; consumers must be idempotent |
| `approval.BindCommand` (instance event → CQRS command bridge) | asynchronous, after the transaction commits | redelivery re-dispatches the command; the command handler must be idempotent |

### `InstanceLifecycleHook`

```go
type InstanceLifecycleHook interface {
    OnInstanceCreated(ctx context.Context, db orm.DB, instance *Instance) error
    OnInstanceTransition(ctx context.Context, db orm.DB, instance *Instance, from, to InstanceStatus) error
}
```

v0.38 generalizes the hook: the former `OnInstanceCompleted(instance, finalStatus)` is replaced by `OnInstanceTransition(instance, from, to)`, which runs inside the same transaction as **every** instance status transition — completion (`to.IsFinal()`), return, withdrawal, resubmission, termination. It runs after the engine-owned business projection has recorded (and, in synchronous mode, applied) the new state, so the hook observes the business table as the transition leaves it; `instance.Status` already carries `to`. Returning an error rolls back the whole transition.

Use lifecycle hooks for invariants that must hold inside the transaction (e.g. allocating a tightly-coupled business row). Use event subscriptions (or `BindCommand`, below) for everything else. Register hooks with `vef.ProvideApprovalLifecycleHook(constructor)` — the constructor must return `approval.InstanceLifecycleHook`, and multiple hooks compose via the `vef:approval:lifecycle_hooks` group. The invocation order across hooks is **unspecified** (FX value groups carry no ordering), so hooks must be mutually independent; any non-nil error stops the remaining hooks. `approval.NewFilteredLifecycleHook(hook, filters...)` wraps a hook with the same `InstanceFilter` vocabulary used by `SubscribeInstance` (below), so a hook can be scoped to specific flow codes or tenants without hand-written predicates.

### `BusinessRefProvider` and `BusinessRefResolver`

```go
type BusinessRefProvider interface {
    OnInstanceCreated(ctx context.Context, db orm.DB, flow *Flow, instance *Instance) (businessRef string, err error)
}

type BusinessRefResolver interface {
    ResolveRecordKey(ctx context.Context, flow *Flow, businessRef string) (BusinessRecordKey, error)
}
```

`BusinessRefProvider` supplies or allocates the opaque `Instance.BusinessRef` when `Flow.BindingMode == BindingBusiness`; inject it with `vef.SupplyBusinessRefProvider`. At instance start the engine resolves the ref into an `approval.BusinessRecordKey` — a map from every configured `KeyColumns` entry to its value — through `BusinessRefResolver`. The default resolver treats a single-column ref verbatim and decodes a multi-column ref from a JSON object; register a replacement with `vef.SupplyBusinessRefResolver` when `BusinessRef` uses another shape (business number, encoded tuple, …) or resolving the key requires a host lookup. A ref that cannot be resolved fails instance creation with `ErrInvalidBusinessRef` (code `40109`); a business-bound flow started without a ref fails with `ErrBusinessRefRequired` (code `40107`).

### Business-State Projection

v0.38 replaces the per-trigger write-back (and its binding listener) with a durable desired-state projection. A business-bound flow configures one `BusinessBinding` document — `approval.BusinessBindingConfig`, stored as jsonb on `apv_flow` and snapshotted immutably onto each deployed `apv_flow_version`; runtime instances only ever read the version snapshot:

| Field | JSON | Meaning |
| --- | --- | --- |
| `TableName` | `tableName` | business table receiving approval state |
| `KeyColumns` | `keyColumns` | record key columns; must exactly match a non-null primary or unique key on the table (verified against the live schema at save time — `ErrBindingSchemaInvalid` / `ErrBindingKeyNotUnique`) |
| `StatusColumn` | `statusColumn` | column receiving the (mapped) instance status |
| `InstanceIDColumn` | `instanceIdColumn` | **mandatory** — the compare-and-set fence that stops a stale instance from overwriting state owned by a newer approval round |
| `StartedAtColumn` | `startedAtColumn` | optional; receives the instance start time |
| `FinishedAtColumn` | `finishedAtColumn` | optional; receives the finish time for final statuses, `NULL` otherwise |
| `StatusMapping` | `statusMapping` | optional `InstanceStatus` → host status value map; missing entries fall back to the raw status string (`ErrBindingStatusMappingInvalid` rejects unknown statuses and blank values) |

Convergence works in three steps:

1. **Claim at start.** Inside the `start_instance` transaction the engine resolves the record key and claims one durable projection row per physical target (`apv_business_projection`, unique per table + record key). A non-final instance already owning the target blocks the claim with `ErrBindingTargetBusy` (code `40108`); a final owner is superseded, and the previously applied owner is retained as the business-table CAS fence until the new desired state is actually written. The instance links to its claim via `Instance.BusinessProjectionID`.
2. **Desired state per transition.** Every instance status transition bumps the projection's `DesiredRevision` and records the full desired state (status, started-at, finished-at). The projection converges on the *latest* desired state — intermediate statuses may be skipped, which is exactly why side effects tied to a specific lifecycle moment belong in `BindCommand` or `SubscribeInstance` instead.
3. **Apply.** `vef.approval.business_binding.consistency` selects when the business row is written:
   - `synchronous` (default) — inside the approval transaction; a write failure rolls the approval action back.
   - `eventual` — the approval action commits immediately; a background worker (`scan_interval`, default `10s`; `batch_size`, default `100`) claims due projections with `FOR UPDATE SKIP LOCKED` plus a lease and applies the latest revision, retrying failures with exponential backoff (1s doubling, capped at 1 hour). Each failed attempt publishes `InstanceBindingFailedEvent` as an operator notification — the durable projection row, not the event, drives the retry.

The applied `UPDATE` sets the configured columns `WHERE <key columns match> [AND <instanceIdColumn> = <applied owner>]`, translating the status through `StatusMapping`. Operators inspect and force convergence through the admin operations `find_business_projections` / `retry_business_projection` (see [RPC Resources](./resources.md#approvaladmin)), and `get_metrics` reports `businessProjectionCounts` and `pendingBusinessProjections`.

The engine owns this projection; it is configuration-driven and not replaced by an extension hook. Hosts add behavior around it with lifecycle hooks, event subscriptions, and command bridges.

### Host Subscriptions: `SubscribeInstance`

```go
func SubscribeInstance[T InstanceEvent](
    bus event.Bus,
    handler func(ctx context.Context, evt T, env event.Envelope) error,
    opts ...InstanceSubscribeOption,
) (event.Unsubscribe, error)
```

`SubscribeInstance` is the declarative wrapper over `event.SubscribeTyped` for instance events (any event type embedding `InstanceEventBase`, e.g. `InstanceCompletedEvent`). Routing filters are data, not predicates: `InstanceSubscribeOption` accepts `approval.ForFlows(codes...)` / `approval.ForTenants(ids...)` (`InstanceFilter` values), each OR-ing within its own dimension and AND-ing across filters passed to the same call; events that fail a filter are acknowledged without invoking the handler. Business predicates (final status, form values) belong in the handler body, not in filters.

Since v0.38 the handler also receives the delivery `event.Envelope`: `Envelope.ID` is the Inbox dedupe key, stable across redeliveries, and therefore the key to build manual idempotency on when the route is at-least-once.

```go
vef.Invoke(func(bus event.Bus) error {
    _, err := approval.SubscribeInstance(bus,
        func(ctx context.Context, evt *approval.InstanceCompletedEvent, env event.Envelope) error {
            // idempotent side effect; env.ID is stable across redeliveries
            return nil
        },
        approval.ForFlows("leave_request"),
        approval.WithGroup("app:leave-writeback"),
    )
    return err
})
```

The consumer group defaults to one derived from the handler's method identity (`vef:sub:<module-relative pkg>.<Type>.<method>`, mirroring the `vef:default:<uuid>` shape used elsewhere). Anonymous handlers fail with `ErrAnonymousSubscriberGroup` unless `approval.WithGroup(name)` is given; a duplicate derived group registered twice in one process fails with `ErrDerivedGroupConflict`. **Renaming or moving a handler changes its derived group** (a new consumer group starts and the old one is orphaned) — pin the group with `WithGroup` before such refactors. `approval.WithConcurrency(n)` sets the per-subscription worker count.

### Command Bridge: `BindCommand`

```go
func BindCommand[E InstanceEvent, C cqrs.Action](
    bus event.Bus,
    commands cqrs.Bus,
    mapper func(evt E, env event.Envelope) (cmd C, ok bool),
    opts ...InstanceSubscribeOption,
) (event.Unsubscribe, error)
```

`BindCommand` (v0.38) subscribes to instance event `E` and dispatches the mapped command `C` through the host's CQRS bus — the declarative bridge from approval facts to host side effects. `mapper` is a pure translation: it shapes the command from the event plus its delivery envelope and reports relevance (`ok=false` acknowledges without dispatching). Business logic belongs in the command handler, which runs the host's full behavior pipeline (transaction, audit, validation); the handler's result is discarded — dispatch is fire-and-record.

```go
vef.Invoke(func(bus event.Bus, commands cqrs.Bus) error {
    _, err := approval.BindCommand(bus, commands,
        func(evt *approval.InstanceCompletedEvent, env event.Envelope) (app.SettleOrderCmd, bool) {
            if evt.FinalStatus != approval.InstanceApproved {
                return app.SettleOrderCmd{}, false
            }
            return app.SettleOrderCmd{InstanceID: evt.InstanceID, DedupeKey: env.ID}, true
        },
        approval.ForFlows("order_settlement"),
    )
    return err
})
```

The consumer group defaults to the command type's identity (`vef:cmd:<module-relative pkg>.<Type>`) — renaming or moving the command deliberately re-keys the subscription, exactly like renaming a `SubscribeInstance` handler; pin with `approval.WithGroup` ahead of such refactors. Binding the same command type twice in one process fails with `ErrDerivedGroupConflict`. Two guards reject invalid bindings up front: `ErrNonCommandAction` (the bound action type is a query) and `ErrUnnamedCommandType` (anonymous struct or interface-typed `C`). Filters and `WithConcurrency` apply as in `SubscribeInstance`.

Delivery inherits the event route's semantics: on an at-least-once transport the command handler must be idempotent — copy `Envelope.ID` into the command when it needs its own dedupe key. Unlike the eventual business projection — which converges on the latest desired state and may skip intermediate statuses — every instance transition dispatches its own command, making `BindCommand` the lane for side effects tied to a specific lifecycle moment.

## Business Identifier Validation

`Flow.BindingMode == BindingBusiness` flows carry SQL identifiers (`BusinessBindingConfig.TableName`, every `KeyColumns` entry, `StatusColumn`, `InstanceIDColumn`, and the optional timestamp columns) that the engine-owned projection interpolates directly into an `UPDATE` template. To prevent SQL injection, the framework whitelists identifiers against `^[A-Za-z_][A-Za-z0-9_]{0,62}$`:

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
| form data | `FormData`, `NewFormData`, `FormSchemaParser`, `FormFieldDefinition`, `FormSnapshot`, `ValidationRule`, `StorageMode`, `StorageJSON`, `StorageTable`, `FieldKind`, `FieldInput`, `FieldNumber`, `FieldDate`, `FieldTextarea`, `FieldSelect`, `FieldUpload`, `FieldTable`, `FieldOption`, `ColumnDataType`, `ColumnString`, `ColumnText`, `ColumnInteger`, `ColumnDecimal`, `ColumnBoolean`, `ColumnDate`, `ColumnDatetime`, `ColumnJSON` |
| table storage | `FormTable`, `FormTableColumn` (`FormTable.SourceFieldKey` is `""` for the main table, or the owning table field's key for a detail-table child projection) |
| flow models | `FlowCategory`, `Flow`, `FlowVersion`, `FlowNode`, `FlowEdge`, `FlowInitiator`, `FlowNodeAssignee`, `FlowNodeCC`, `VersionStatus`, `VersionDraft`, `VersionPublished`, `VersionArchived`, `ActionLog`, `OperatorInfo`, `UrgeRecord`, `DefaultTenantID` |
| node design | `FlowDefinition`, `NodeDefinition`, `EdgeDefinition`, `Position`, `NodeData`, `BaseNodeData`, `StartNodeData`, `ApprovalNodeData`, `HandleNodeData`, `ConditionNodeData`, `CCNodeData`, `EndNodeData`, `ErrUnknownNodeKind`, `ErrNodeDataUnmarshal` |
| conditions | `ConditionKind`, `ConditionField`, `ConditionExpression`, `Condition`, `ConditionGroup`, `ConditionBranch`, `EvaluationContext`, `ConditionEvaluator`, `InstanceGlobalsResolver`, `AggregateKind`, `AggregateSum`, `AggregateCount`, `AggregateAvg`, `Aggregator` |
| initiators and assignees | `InitiatorKind`, `InitiatorUser`, `InitiatorRole`, `InitiatorDepartment`, `AssigneeKind`, `AssigneeDefinition`, `AssigneeService`, `ResolvedAssignee`, `UserInfo`, `UserInfoResolver`, `RoleMembershipChecker`, `AddAssigneeType`, `AddAssigneeBefore`, `AddAssigneeAfter`, `AddAssigneeParallel` |
| CC | `CCKind`, `CCUser`, `CCRole`, `CCDepartment`, `CCFormField`, `CCTiming`, `CCTimingAlways`, `CCTimingOnApprove`, `CCTimingOnReject`, `CCDefinition`, `CCRecord` (`CCRecord.VisitID` scopes the record to one node traversal, mirroring `Task.VisitID`) |
| business binding & projection | `BusinessBindingConfig`, `BusinessRecordKey`, `BusinessProjection` (model, table `apv_business_projection`), `BindingProjectionStatus` (`BindingProjectionPending` / `Processing` / `Applied` / `Failed`), `BindingTrigger`, `BindingTriggerStarted`, `BindingTriggerCompleted`, `BindingTriggerReturned`, `BindingTriggerWithdrawn`, `BindingTriggerResubmitted` |
| instance subscriptions | `SubscribeInstance`, `BindCommand`, `InstanceEvent`, `InstanceFilter`, `InstanceSubscribeOption`, `ForFlows`, `ForTenants`, `WithGroup`, `WithConcurrency`, `NewFilteredLifecycleHook`, `ErrAnonymousSubscriberGroup`, `ErrDerivedGroupConflict`, `ErrNonCommandAction`, `ErrUnnamedCommandType` |
| node behavior | `ApprovalMethod`, `TaskNodeData`, `ExecutionType`, `ExecutionManual`, `ExecutionAutoPass`, `ExecutionAutoReject`, `ConsecutiveApproverAction`, `ConsecutiveApproverNone`, `ConsecutiveApproverAutoPass`, `SameApplicantAction`, `SameApplicantSelfApprove`, `SameApplicantAutoPass`, `SameApplicantTransferSuperior`, `Permission`, `PermissionVisible`, `PermissionEditable`, `PermissionRequired`, `PermissionHidden`, `DefaultExecutionType`, `DefaultApprovalMethod`, `DefaultPassRule`, `DefaultEmptyAssigneeAction`, `DefaultSameApplicantAction`, `DefaultConsecutiveApproverAction`, `DefaultRollbackType`, `DefaultRollbackDataStrategy`, `DefaultTimeoutAction`, `DefaultCCTiming`, `DefaultHandleApprovalMethod`, `DefaultHandlePassRule`, `DefaultUrgeCooldownMinutes` |
| rollback and timeouts | `RollbackType`, `RollbackNone`, `RollbackPrevious`, `RollbackStart`, `RollbackAny`, `RollbackSpecified`, `RollbackDataStrategy`, `RollbackDataClear`, `RollbackDataKeep`, `EmptyAssigneeAction`, `EmptyAssigneeAutoPass`, `EmptyAssigneeTransferAdmin`, `EmptyAssigneeTransferSuperior`, `EmptyAssigneeTransferApplicant`, `EmptyAssigneeTransferSpecified`, `TimeoutAction`, `TimeoutActionNone`, `TimeoutActionAutoPass`, `TimeoutActionAutoReject`, `TimeoutActionNotify`, `TimeoutActionTransferAdmin` |
| action and status enums | `ActionType`, `InstanceStatus`, `TaskStatus`, `NodeKind`, `StorageMode`, `VersionStatus` |
| pass rules | `PassRule`, `PassRuleContext`, `PassRuleStrategy`, `PassRuleResult`, `PassRulePending`, `PassRulePassed`, `PassRuleRejected` |
| progress views | `TimelineEntryKind`, `TimelineEntry`, `NodeVisitStatus`, `NodeProgressStatus`, `InstanceFlowGraph`, `FlowGraphNode`, `FlowGraphNodeData`, `FlowGraphEdge`, `NodeParticipant`, `Activity`, `ActivityUrge`, `CCRecipient` |
| events | all `New...Event` constructors, `DomainEvent`, `InstanceEventBase`, `TaskEventBase`, `FlowEventBase`, `NewInstanceEventBase`, `NewTaskEventBase`, `NewFlowEventBase`, `PayloadOccurredAt`, `AllEventTypes`, and the `EventType...` constants |
| extension interfaces | `InstanceLifecycleHook`, `BusinessRefProvider`, `BusinessRefResolver`, `InstanceNoGenerator`, `ConditionEvaluator`, `InstanceGlobalsResolver`, `PrincipalTenantResolver`, `PrincipalDepartmentResolver`, `RoleMembershipChecker` |
| DI helpers (package `vef`) | `SupplyBusinessRefProvider`, `SupplyBusinessRefResolver`, `ProvideApprovalLifecycleHook`, `ProvideApprovalAggregator`, `ProvideApprovalFormSchemaParser` |
| admin DTOs | package `approval/admin`: `Instance`, `InstanceDetail`, `InstanceDetailInfo`, `Task`, `ActionLog`, `Metrics`, `BusinessProjection` |
| user DTOs | package `approval/my`: `PendingTask`, `CompletedTask`, `CCRecord`, `InitiatedInstance`, `AvailableFlow`, `InstanceDetail`, `InstanceInfo`, `PendingCounts` |

---

Next: back to the [Overview](./overview.md) for wiring and configuration, or [RPC Resources](./resources.md) for the API surface these events originate from.
