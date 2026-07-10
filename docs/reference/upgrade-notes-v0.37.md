---
sidebar_position: 7
---

# Upgrade Notes: v0.36 / v0.37

This page is the cross-version audit map for the backend commits from
`v0.35.0` through `v0.37.0` (`v0.35.0..v0.37.0`). Both releases center on the
approval module; v0.37 also adds event-stream observability. Use it when
upgrading an application whose docs or integration assumptions were last
checked against [Upgrade Notes to v0.35](../reference/upgrade-notes-v0.35).

It is not a replacement for the generated indexes. After applying the migration
notes below, verify exact Go symbols and wire fields against the
[Public API Index](../reference/public-api-index) and
[Runtime API Index](../reference/runtime-api-index).

## Immediate Checklist

- Update RPC clients of `approval/flow`: the update action is now `update`
  (was `update_flow`), and the action-log permission is
  `approval.action_log.query` (was `approval.log.query`). Update permission
  seeds accordingly.
- Update list-view consumers: admin instance/task rows and my
  pending/completed/CC rows expose people as nested `UserInfo` objects
  (`applicant`, `assignee`) instead of flat `*Id` / `*Name` string pairs.
- Delegation create/update requires `startTime` and `endTime`, in the
  canonical `timex.DateTime` wire format (`2006-01-02 15:04:05`), not RFC 3339.
- Approval event subscribers: `InstanceBindingFailedEvent` replaced
  `finalStatus` with `trigger` + `status`; completed/withdrawn events gained
  `reason`; rolled-back/returned events gained `opinion`.
- Go symbol rename: `Flow.BusinessPkField` is now `Flow.BusinessPKField`
  (JSON stays `businessPkField`). Event constructor signatures for withdrawn,
  rolled-back, returned, and binding-failed events changed.
- Error-code matchers: `invalid storage mode` moved from 40012 to 40014 and
  `flow binding locked` from 40013 to 40015 (v0.36 accidentally reused those
  codes for the new binding-mode / initiator-kind validation errors; v0.37
  deduplicated them).
- Run the approval DB migrations: `apv_cc_record` gained `visit_id`,
  `apv_form_table` gained `source_field_key`, and `apv_flow` gained the three
  optional write-back linkage columns.
- Re-test flows that rely on rollback targeting, CC read-confirmation, or
  add-assignee: all three were tightened (see the v0.36 behavior fixes below).
- After upgrading to v0.37, prefer `approval.SubscribeInstance` over raw
  `event.SubscribeTyped` for instance events, and consider enabling
  `idle_group_retention` plus the `get_event_streams` monitor action to detect
  and reclaim orphaned Redis consumer groups.

## Release-by-Release Audit

| Release | User-facing changes to review |
| --- | --- |
| `v0.36.0` | Approval only: single-level detail-table form fields with per-field child projection tables, aggregate field conditions (`sum` / `count` / `avg` plus custom aggregators), wire renames (`update` action, `approval.action_log.query` permission, nested `UserInfo` list rows, required canonical delegation window), operator `reason` / `opinion` on lifecycle events, save/deploy-time validation of binding mode, initiator kind, empty condition branches, and sequential + parallel add-assignee, and behavior fixes for rollback bounding, per-visit CC scoping, and add-assignee splicing. |
| `v0.37.0` | Approval five-trigger business write-back matrix (new `Flow` linkage columns, reshaped `InstanceBindingFailedEvent`, deduplicated error codes), declarative `approval.SubscribeInstance` with derived consumer groups and `NewFilteredLifecycleHook`, and event-stream observability (`event.StreamInspector`, `get_event_streams`, idle consumer-group reclamation). |

## v0.36.0

Every v0.36 change is in the approval module.

### Breaking: detail-table form fields

`approval.FieldKind` gained `FieldTable` (`"table"`): a single-level detail
table whose value is a list of rows. `FormFieldDefinition.Columns` defines the
row shape; each column is itself a field definition reusing the existing kinds
and validation rules, except that a column must not be another table. For the
table field itself, `Validation.MinLength` / `MaxLength` bound the row count
and `IsRequired` means "at least one row".

For flows with `StorageTable`, each detail-table field is projected into its
own child table named `apv_form_<versionID>__<sanitized field key>`, with rows
keyed by `instance_id` (indexed, not unique) and ordered by `row_index`. The
`approval.FormTable` registry now records one row per physical table — the
main projection table plus one child table per detail-table field — and gained
`SourceFieldKey` (`sourceFieldKey`; empty for the main table). Uniqueness moved
from `version_id` to `(version_id, source_field_key)`. Table storage also now
stringifies non-string scalars bound to text columns instead of failing.

### Breaking: aggregate field conditions

Field conditions can fold a detail table into one comparable number.
`approval.Condition` gained two wire fields:

```json
{
  "kind": "field",
  "subject": "items",
  "aggregate": "sum",
  "column": "amount",
  "operator": "gt",
  "value": 10000
}
```

`approval.AggregateKind` defines `sum`, `count`, and `avg`. Semantics follow
SQL aggregates: `count` is the row count and must leave `column` empty,
`sum` over an empty table is 0, and `avg` over an empty table matches no
comparison (NULL semantics). Folds compute in `float64`, so prefer ordering
operators over `eq` / `ne` for fractional amounts. Deploy validation derives
the column-required/forbidden rule from `AggregateKind.FoldsColumn()`.

Custom aggregates implement `approval.Aggregator` (`Kind()` plus
`Fold(values, rowCount) (result, matchable)`) and register through
`vef.ProvideApprovalAggregator`:

```go
vef.ProvideApprovalAggregator(func() approval.Aggregator { return myMedian{} })
```

### Breaking: wire and permission renames

- The `approval/flow` resource action `update_flow` is now `update`. The Go
  params types followed (`UpdateFlowParams` → `UpdateParams`,
  `StartInstanceParams` → `StartParams`).
- The action-log query permission is `approval.action_log.query` (was
  `approval.log.query`).
- List rows nest person snapshots. Before / after for an admin instance row:

```json
{ "applicantId": "u1", "applicantName": "Alice" }
```

```json
{ "applicant": { "id": "u1", "name": "Alice", "departmentId": "d1", "departmentName": "Sales" } }
```

  The same shape applies to `assignee` on admin task rows and my
  pending/completed rows, and to the people fields on my CC rows.
- Delegation create/update takes the canonical `timex.DateTime` format and the
  window is now mandatory:

```json
{ "startTime": "2030-01-01T00:00:00Z", "endTime": null }
```

```json
{ "startTime": "2030-01-01 00:00:00", "endTime": "2030-06-01 00:00:00" }
```

### Breaking: lifecycle events carry reason and opinion

- `InstanceCompletedEvent.Reason` (`reason`, optional) carries the
  administrator's stated reason for `terminated` completions; it is nil for
  approved/rejected, whose deciding opinions live on the task events.
- `InstanceWithdrawnEvent.Reason` carries the applicant's withdrawal reason.
- `InstanceRolledBackEvent.Opinion` and `InstanceReturnedEvent.Opinion` carry
  the operator's rollback opinion.

The corresponding `New*Event` constructors gained a trailing parameter, so Go
code that builds these events directly must be updated.

### Breaking: stricter save- and deploy-time validation

- Flow save rejects out-of-enum `BindingMode` and `InitiatorKind` values
  (`BindingMode.IsValid()` / `InitiatorKind.IsValid()` are exported). A typo'd
  binding mode previously behaved like `standalone` and silently disabled the
  business write-back.
- Deploy rejects condition branches with no condition groups and groups with
  no conditions — the engine treats structurally-empty conditions as an
  unconditional match, which silently shadowed every lower-priority branch.
- Deploy rejects the `parallel` add-assignee type on sequential-approval
  nodes; a sequential queue has nothing for a parallel addition to join.

### Breaking behavior fixes

- **Rollback targets are bounded to the visit trail.** A rollback target must
  be a decision node (approval / handle) or the start node that this instance
  actually traversed (has a concluded visit). Previously any node in the
  version was accepted, which let a task holder target the End node and
  force-complete the instance.
- **CC records and the read-confirm gate are scoped per node visit.**
  `CCRecord` gained `VisitID` (`visitId`), so a rollback redo gets its own
  notification and read-confirm cycle instead of being silently satisfied by a
  prior round's records.
- **Add-assignee skips users who already decided in the open visit**
  (approved / rejected / handled), so pass rules cannot force the same person
  to decide twice in one round. Removed or transferred users may be re-added.
- **Sequential add-assignee splices at the anchor** instead of appending to
  the end of the queue: "before" additions take over the anchor's position,
  "after" additions slot in right behind it.

### Other fixes

- Withdrawing an instance keeps the paused resting node active in the flow
  graph instead of showing no current node.
- The my-detail action list offers `urge` only to the applicant or users with
  an assignee task, matching what the urge handler actually authorizes.
- Paginated approval lists use an `id` tiebreaker for stable ordering.
- New indexes speed up completed-task listing (`finished_at`) and
  `visit_id` lookups on `apv_task`.

## v0.37.0

### Breaking: five-trigger business write-back matrix

The engine-owned business write-back now covers the full instance lifecycle.
`approval.BindingTrigger` identifies which moment drove a write-back; each
trigger projects a fixed column subset (a column is only written when the flow
configures it):

| Trigger | status | instance_id | started_at | finished_at |
| --- | --- | --- | --- | --- |
| `started` | running | instance ID | now | NULL |
| `completed` | final status | — | — | `FinishedAt` |
| `returned` | returned | — | — | — |
| `withdrawn` | withdrawn | — | — | — |
| `resubmitted` | running | — | — | NULL |

- `approval.Flow` gained three optional linkage columns:
  `BusinessInstanceIDField`, `BusinessStartedAtField`,
  `BusinessFinishedAtField` (`businessInstanceIdField` / `businessStartedAtField`
  / `businessFinishedAtField`). Only the status column stays mandatory for a
  business-bound flow. `Flow.BusinessPkField` was renamed to
  `Flow.BusinessPKField` in Go (JSON unchanged).
- The `started` write-back runs synchronously inside the `start_instance`
  transaction — a failure rolls back the whole initiation. The other four run
  asynchronously through the binding listener with
  `InstanceBindingFailedEvent` compensation, so `started` never appears in
  that event.
- `InstanceBindingFailedEvent` was reshaped; before / after:

```json
{ "finalStatus": "approved", "businessTable": "orders", "errorMessage": "..." }
```

```json
{ "trigger": "completed", "status": "approved", "businessTable": "orders", "errorMessage": "..." }
```

- Flow save validates that binding fields name distinct business columns
  (`approval binding columns conflict`, code 40016); a duplicated column would
  render as `SET col = ?, col = ?` and fail at runtime in the applicant's
  start transaction. The existing binding freeze (no binding changes while
  instances are running) now also covers the new linkage columns.
- Error codes were deduplicated: `invalid storage mode` is 40014 (was 40012)
  and `flow binding locked` is 40015 (was 40013); 40012 / 40013 now belong to
  the v0.36 binding-mode / initiator-kind validation errors.

### Declarative instance subscriptions

`approval.SubscribeInstance` subscribes a typed handler to one instance event
type with declarative routing filters, replacing hand-rolled
`event.SubscribeTyped` wiring:

```go
unsubscribe, err := approval.SubscribeInstance(bus,
    svc.OnCompleted, // func(ctx context.Context, evt *approval.InstanceCompletedEvent) error
    approval.ForFlows("expense_claim"),
    approval.WithGroup("mms.expense.completed"),
)
```

- `approval.InstanceFilter` (`ForFlows`, `ForTenants`) captures "is this
  instance mine?" as data. Within one filter a populated dimension matches by
  OR; multiple filters must all match (AND). Non-matching events are
  acknowledged without invoking the handler. Business predicates (final
  status, form values) stay in the handler body.
- The consumer group defaults to a name derived from the handler's method
  identity (`vef:sub:<pkg>.<Type>.<Method>`, with the main module prefix
  stripped). Anonymous functions fail with `ErrAnonymousSubscriberGroup`;
  deriving the same group twice in one process fails with
  `ErrDerivedGroupConflict`. Renaming or moving a handler changes its derived
  group and orphans the old one — pin `approval.WithGroup` for
  production-critical subscribers before such refactors.
  `approval.WithConcurrency` forwards a per-subscription worker count.
- `approval.NewFilteredLifecycleHook(hook, filters...)` applies the same
  filters to a synchronous `InstanceLifecycleHook`, so a hook serving a single
  flow can be scoped at registration via `vef.ProvideApprovalLifecycleHook`.

### Event stream observability

The redis_stream transport now exposes consumer-group state and can reclaim
orphaned groups left behind by decommissioned subscribers:

- `event.StreamInspector` lists every stream under the transport prefix with
  its consumer groups (`StreamInfo` / `StreamGroupInfo`: name, consumers,
  pending, lag, `lastDeliveredId`). It is an optional dependency — nil when
  the redis_stream transport is off.
- The `sys/monitor` resource gained the `get_event_streams` action, returning
  `monitor.EventStreamsInfo` (`enabled` plus the stream list). A group whose
  lag keeps growing while all consumers stay idle is an orphan candidate.
- Idle-group reclamation is opt-in:

```toml
[vef.event.transports.redis_stream]
enabled = true
idle_group_retention = "72h"        # zero (default) disables the sweep
idle_group_sweep_interval = "10m"   # default 10m
```

A group is destroyed only when it has no pending entries and every consumer
record has been idle beyond the retention window.

See [Approval Module](../approval), [Event Bus](../infrastructure/event-bus),
and [Monitor](../infrastructure/monitor) for the full contracts.
