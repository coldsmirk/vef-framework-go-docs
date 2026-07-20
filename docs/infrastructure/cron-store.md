---
sidebar_position: 3
---

# Durable Schedules

The durable schedule store (v0.39) extends the [cron module](./cron) with
database-persisted schedules: cluster-wide single fire per occurrence,
operator-editable triggers, misfire policies, a run journal, and crash
recovery. The in-memory scheduler keeps serving process-local jobs; the store
is a separate engine for jobs that must survive restarts and coordinate
across nodes.

It is off by default. Enabling it loads schedules from the primary data
source and mounts the `sys/cron/schedule` and `sys/cron/run` resources:

```toml
[vef.cron.store]
enabled = true
auto_migrate = true   # create crn_schedule / crn_fire_request / crn_run on start
```

## Model

Two concepts drive the engine:

- A **job handler** (`cron.JobHandler`) is Go code registered under a unique
  name at boot. Handlers are the *what*.
- A **schedule** (`cron.Schedule`, table `crn_schedule`) is a persisted
  trigger: *when* to fire *which* job, with which params, under which
  policies. Schedules are data — created in code or through the management
  API, edited at runtime.

Every fire is journaled as a **run** (`cron.Run`, table `crn_run`).

### Registering job handlers

```go
vef.ProvideCronJobHandler(func(svc *ReportService) cron.JobHandler {
    return cron.NewTypedJobHandler("daily-report",
        func(ctx context.Context, params ReportParams) error {
            return svc.Generate(ctx, params)
        },
        // Optional: seed a default schedule at boot when none of this
        // name exists yet; operator changes are never overwritten.
        cron.WithDefaultSchedule(cron.ScheduleSpec{
            Trigger: cron.Expr("0 2 * * *", "Asia/Shanghai"),
        }),
    )
})
```

| API | Contract |
| --- | --- |
| `cron.JobHandler` | `Name() string` + `Execute(ctx, execution) error`; exactly one handler per job name |
| `cron.NewJobHandler(name, execute, opts...)` | adapts a function; `execute` receives the full `cron.Execution` |
| `cron.NewTypedJobHandler[P](name, execute, opts...)` | decodes the schedule's params into `P` before the function runs; a decode failure journals the run as failed without invoking it |
| `cron.WithDefaultSchedule(spec)` | ships a default schedule; the store seeds it at boot when absent. The spec's `Name` falls back to the job name |
| `cron.Execution` | read-only view of the run: `RunID`, `ScheduleID`, `ScheduleName`, `JobName`, `ScheduledAt` (logical fire time), `Params` (raw JSON), `BindParams(v)` |

A fire is claimed by at most one node, but a crashed run re-fires when the
schedule sets `Recover` — delivery is at-least-once, so handlers should be
idempotent.

### Triggers

`cron.TriggerSpec` declares when a schedule fires; build specs with the
constructors:

| Constructor | Kind | Semantics |
| --- | --- | --- |
| `cron.Expr(expr, timezone)` | `cron` | cron expression evaluated in an IANA timezone. 5-field, 6-field (leading seconds), and `@`-descriptors (`@daily`, `@every 90m`) are accepted. Empty timezone resolves to `UTC` — durable schedules never depend on a node's process-local zone (`"Local"` is rejected). The embedded tzdata keeps timezones working on zoneinfo-less deployments |
| `cron.Every(duration)` | `interval` | fixed rate, minimum `1s` (`cron.MinInterval`). The rate is anchored to the schedule's start (`StartsAt`, else creation time), keeping the fire phase stable regardless of catch-ups or manual fires |
| `cron.Once(at)` | `once` | a single fire |

Fields not belonging to the selected kind are rejected
(`ErrTriggerFieldsConflict`), as are unparsable expressions, unloadable
timezones, sub-second intervals, and missing fire times.

### Schedule spec

`cron.ScheduleSpec` declares a schedule to create, or the desired state of an
update:

| Field | Meaning |
| --- | --- |
| `Name` | unique management key; on a seeded default schedule it falls back to the job name |
| `JobName` | the registered `JobHandler` to execute |
| `Trigger` | when to fire (above) |
| `Params` | JSON-marshaled and delivered verbatim to the handler on every run |
| `StartsAt` / `EndsAt` | optional fire window; `StartsAt` also anchors the fixed-rate phase of interval triggers |
| `MisfirePolicy` | `fire_now` (default) or `skip` (below) |
| `ConcurrencyPolicy` | `forbid` (default) or `allow` (below) |
| `Recover` | re-fire runs abandoned mid-execution; requires idempotent handlers |
| `Timeout` | per-run bound (whole milliseconds); zero inherits `vef.cron.store.run_timeout` |
| `Enabled` | initial/updated enablement; `nil` means enabled |

Caller-supplied times (`Trigger.At`, `StartsAt`, `EndsAt`) are normalized to
the local wall clock before persistence, so an instant built in another zone
still denotes the same moment when read back.

## Policies

### Misfire

A fire that starts later than `vef.cron.store.misfire_threshold` (default
`1m`) counts as misfired — downtime, a paused schedule, or no free executor.
The schedule's `MisfirePolicy` then applies:

| Policy | Behavior |
| --- | --- |
| `fire_now` (default) | run one catch-up fire immediately and resume the regular sequence from now |
| `skip` | advance to the next future fire without running |

Whichever policy applies, occurrences that will never run are journaled as a
single `missed` run covering the whole gap (`missedCount` carries the
occurrence count).

### Concurrency

| Policy | Behavior |
| --- | --- |
| `forbid` (default) | a fire that would overlap a still-running run of the same schedule is suppressed and journaled as `skipped`. Recovery requests stay pending until the active run ends |
| `allow` | runs of the same schedule may overlap |

### Pause / resume semantics

`Pause` clears the operator-owned `isEnabled` flag; running fires are
unaffected. The fire cursor is deliberately preserved while paused, so
`Resume` hands the paused gap to the misfire policy instead of silently
dropping it: under `fire_now` one catch-up runs immediately, under `skip`
the schedule waits for its next regular fire.

### Manual trigger

`TriggerNow` persists one independent immediate fire request (table
`crn_fire_request`) — single node, journaled, concurrency policy respected —
without moving the regular trigger cursor. Recovery re-fires ride the same
request table. A paused schedule refuses with `ErrScheduleDisabled`. When a
manual fire and a regular occurrence land on the same logical instant, the
regular one wins and both are journaled per policy.

## Execution and Recovery

- Each node polls for due schedules (`poll_interval`, default `5s`, sleeping
  adaptively until the nearest known fire — the interval is the visibility
  latency of schedules created on other nodes, not the fire precision),
  claims up to `batch_size` fires transactionally, and executes them on up to
  `max_concurrent` local slots.
- Executors heartbeat their running journal rows every `heartbeat_interval`
  (default `10s`). A running row whose heartbeat goes stale for
  `abandoned_after` (default `1m`; must be at least twice the heartbeat
  interval) is taken over in one transaction by the recovery sweep and marked
  `abandoned`; schedules with `Recover` re-fire it as a fresh run.
- A run that outlives its timeout is journaled as `failed`; graceful shutdown
  journals interrupted runs as `canceled`.
- Reshaping a schedule (trigger/window changes) recomputes the next fire but
  preserves the fire history; renames keep the journal linked by denormalized
  names.

## The Run Journal

`cron.Run` (table `crn_run`) records every fire. Rows survive schedule
deletion — `scheduleName` and `jobName` are denormalized for that reason.

| Field | Type | Meaning |
| --- | --- | --- |
| `id` | `string` | journal row ID |
| `scheduleId` / `scheduleName` | `string` | the schedule that fired |
| `jobName` | `string` | the executed handler |
| `scheduledAtUnixMs` | `int64` | the logical fire time; a catch-up fire starts later than it. Deliberately not unique: manual and recovery fires may share one instant |
| `claimedAtUnixMs` | `int64` | when a node claimed the fire |
| `status` | `string` | `running`, `succeeded`, `failed`, `missed`, `skipped`, `abandoned`, `canceled` |
| `nodeId` | `string` | executing node; empty on rows that never executed (`missed`, `skipped`) |
| `startedAtUnixMs` / `finishedAtUnixMs` | `int64` | execution window |
| `durationMs` | `int64` | execution duration |
| `heartbeatAtUnixMs` | `int64` | executor liveness signal; a stale heartbeat turns the run `abandoned` |
| `error` | `string` | failure message, truncated; empty on success |
| `missedCount` | `int` | occurrences a `missed` row covers |

`run_retention` prunes terminal journal rows older than the window (hourly
sweep); zero keeps rows forever — deletion of the journal is strictly opt-in.

## Programmatic Management

`cron.ScheduleManager` is available in DI whenever the cron module is loaded;
with the store disabled every method returns `ErrStoreDisabled`. API
mutations and programmatic ones share one validation and wake path.

| Method | Contract |
| --- | --- |
| `Create(ctx, spec)` | validates and persists a new schedule; the job name must be registered on this node; a taken name fails with `ErrScheduleExists` |
| `Update(ctx, name, spec)` | reshapes the named schedule, including a rename when the spec carries a different untaken name; trigger/window changes recompute the next fire |
| `Delete(ctx, name)` | removes the schedule; journaled runs are kept |
| `Pause(ctx, name)` / `Resume(ctx, name)` | see [pause semantics](#pause--resume-semantics) |
| `TriggerNow(ctx, name)` | see [manual trigger](#manual-trigger) |
| `Get(ctx, name)` | returns the named schedule or `ErrScheduleNotFound` |
| `List(ctx, filter)` | schedules matching `ScheduleFilter` (`JobName`, `Enabled *bool`), ordered by name |
| `ListRuns(ctx, filter)` | journal records matching `RunFilter` (`ScheduleName`, `JobName`, `Statuses`, `Since`/`Until` on the logical fire time, `Limit` — zero resolves to 100, capped at 1000), newest first |

## Events

Both topics are best-effort operational notifications published outside any
transaction on the default event route — subscribe for alerting; never drive
correctness from them (the run journal is the durable truth).

| Topic | Event | Fields |
| --- | --- | --- |
| `vef.cron.run.failed` | `cron.RunFailedEvent` | `runId`, `scheduleName`, `jobName`, `scheduledAtUnixMs`, `nodeId`, `error` |
| `vef.cron.run.abandoned` | `cron.RunAbandonedEvent` | `runId`, `scheduleName`, `jobName`, `scheduledAtUnixMs`, `nodeId` |

## RPC Resources

With the store enabled, two management resources mount under `/api`. With it
disabled the resources mount no operations — a feature that is off exposes no
surface. Mutating operations are audited.

### `sys/cron/schedule`

| Action | Permission | Input | Output |
| --- | --- | --- | --- |
| `find_page` | `cron.schedule.query` | `ScheduleSearch` + pageable meta | `page.Page[Schedule]` |
| `get` | `cron.schedule.query` | `ScheduleNameParams` | `ScheduleDetail` |
| `list_jobs` | `cron.schedule.query` | none | `string[]` |
| `preview_fires` | `cron.schedule.query` | `PreviewFiresParams` | `FiresPreview` |
| `create` | `cron.schedule.manage` (audited) | `ScheduleParams` | created `Schedule` |
| `update` | `cron.schedule.manage` (audited) | `ScheduleParams` | updated `Schedule` |
| `delete` | `cron.schedule.manage` (audited) | `ScheduleNameParams` | success |
| `pause` | `cron.schedule.manage` (audited) | `ScheduleNameParams` | success |
| `resume` | `cron.schedule.manage` (audited) | `ScheduleNameParams` | success |
| `trigger_now` | `cron.schedule.manage` (audited) | `ScheduleNameParams` | success |

`ScheduleSearch` (query filters for `find_page`):

| Field | Type | Match | Description |
| --- | --- | --- | --- |
| `name` | `string` | contains | filter by schedule name fragment |
| `jobName` | `string` | equals | filter by job name |
| `kind` | `string` | equals | trigger kind: `cron`, `interval`, or `once` |
| `isEnabled` | `bool` | equals | filter by enablement |

`ScheduleNameParams` (used by `get`, `delete`, `pause`, `resume`,
`trigger_now`):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `name` | `string` | Yes | the schedule's unique name |

`ScheduleParams` (create/update; unknown fields are rejected — the params
struct is strict):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `name` | `string` | Yes | on create: the new schedule's unique name; on update: the schedule being addressed |
| `newName` | `string` | No | update only: renames the schedule after addressing it by `name` |
| `jobName` | `string` | Yes | registered job handler to execute; unregistered names fail with `ErrJobNotRegistered` |
| `trigger` | `TriggerParams` | Yes | the trigger definition (below) |
| `params` | any JSON value | No | delivered verbatim to the handler on every run |
| `startsAtUnixMs` | `int64` (unix ms) | No | fire window start; also anchors the fixed-rate phase of interval triggers |
| `endsAtUnixMs` | `int64` (unix ms) | No | fire window end; must be after `startsAtUnixMs` |
| `misfirePolicy` | `string` | No | `fire_now` (default when omitted) or `skip` |
| `concurrencyPolicy` | `string` | No | `forbid` (default when omitted) or `allow` |
| `recover` | `bool` | No | re-fire abandoned runs; handlers must be idempotent |
| `timeoutMs` | `int64` | No | per-run timeout; zero inherits `vef.cron.store.run_timeout`; negative values are rejected |
| `enabled` | `bool` | No | omitted means enabled |

`TriggerParams` (exactly the fields of the selected `kind`; extra fields
fail with `ErrTriggerInvalid`):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `kind` | `string` | Yes | `cron`, `interval`, or `once` |
| `expr` | `string` | for `cron` | cron expression (5/6-field or `@`-descriptor) |
| `timezone` | `string` | No (`cron` only) | IANA zone the expression is evaluated in; empty means `UTC`; `"Local"` is rejected |
| `everyMs` | `int64` | for `interval` | fixed rate in milliseconds, minimum 1000 |
| `atUnixMs` | `int64` (unix ms) | for `once` | the single fire time |

`ScheduleDetail` (`get` response):

| Field | Type | Description |
| --- | --- | --- |
| `schedule` | `Schedule` | the schedule row (below) |
| `nextFiresUnixMs` | `int64[]` | preview of the next (up to 5) exact fire times from now. An overdue cursor is projected through the schedule's misfire policy; a paused or spent schedule returns an empty list |

`Schedule` (returned by `get`, `create`, `update`, and `find_page` items;
standard audit columns omitted):

| Field | Type | Description |
| --- | --- | --- |
| `name` | `string` | unique management key |
| `jobName` | `string` | the handler the schedule fires |
| `kind` | `string` | trigger kind |
| `expr` | `string` | cron expression (cron kind) |
| `timezone` | `string` | evaluation zone (cron kind) |
| `everyMs` | `int64` | fixed rate (interval kind) |
| `fireAtUnixMs` | `int64` | single fire time (once kind); absent otherwise |
| `startsAtUnixMs` / `endsAtUnixMs` | `int64` | fire window bounds; absent when unbounded |
| `anchorAtUnixMs` | `int64` | fixed-rate phase anchor (creation time, or `startsAtUnixMs` when set) |
| `params` | JSON | handler params, verbatim |
| `misfirePolicy` | `string` | `fire_now` or `skip` |
| `concurrencyPolicy` | `string` | `forbid` or `allow` |
| `recover` | `bool` | abandoned-run re-fire flag |
| `timeoutMs` | `int64` | per-run timeout; `0` inherits the configured default |
| `isEnabled` | `bool` | operator-owned enablement (`pause` clears, `resume` restores) |
| `nextFireAtUnixMs` | `int64` | next fire the engine will claim; absent when the trigger yields no further occurrence (completed one-shot, expired window) — pausing preserves it |
| `lastFireAtUnixMs` | `int64` | most recent claimed fire's logical time; absent before the first fire |

`list_jobs` returns the job names registered **on this node** — the
vocabulary the schedule editor's job picker offers. Heterogeneous
deployments may register different sets per node; the answering node's view
is returned.

`PreviewFiresParams` (`preview_fires` — editor-time validation of an unsaved
trigger against the real parser; rejects exactly what a save would):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `trigger` | `TriggerParams` | Yes | the unsaved trigger to project |
| `startsAtUnixMs` | `int64` (unix ms) | No | window start to project under |
| `endsAtUnixMs` | `int64` (unix ms) | No | window end; must be after the start |

`FiresPreview` response:

| Field | Type | Description |
| --- | --- | --- |
| `nextFiresUnixMs` | `int64[]` | the trigger's upcoming fire times from now (up to 5); empty when it yields no occurrence inside its window |

### `sys/cron/run`

Read-only journal views: the paged view for browsing and the single-record
view for the full error text. Default order is newest claim first.

| Action | Permission | Input | Output |
| --- | --- | --- | --- |
| `find_page` | `cron.run.query` | `RunSearch` + pageable meta | `page.Page[Run]` |
| `find_one` | `cron.run.query` | `RunSearch` | one `Run` |

`RunSearch` (query filters):

| Field | Type | Match | Description |
| --- | --- | --- | --- |
| `id` | `string` | equals | addresses one journal row — `find_one` has no other way to name the record |
| `scheduleName` | `string` | equals | filter by schedule |
| `jobName` | `string` | equals | filter by job |
| `status` | `string` | equals | one of the run statuses |
| `nodeId` | `string` | equals | filter by executing node |
| `scheduledAtFromUnixMs` | `int64` | ≥ | logical fire time lower bound |
| `scheduledAtToUnixMs` | `int64` | ≤ | logical fire time upper bound |

The `Run` response fields are the [run journal columns](#the-run-journal)
plus the creation audit columns.

## Error Codes

Cron API errors use response codes `2700`–`2799` and ride HTTP 200 with the
failure in the body code.

| Code | Error | Meaning |
| --- | --- | --- |
| `2700` | `ErrScheduleNotFound` | schedule lookup failed |
| `2701` | `ErrScheduleExists` | schedule name already taken |
| `2702` | `ErrScheduleDisabled` | manual trigger against a paused schedule |
| `2703` | `ErrTriggerInvalid(reason)` | trigger failed validation (conflicting fields, bad expression, bad timezone, short interval, missing fire time) |
| `2704` | `ErrJobNotRegistered` | schedule references a job name not registered on this node |
| `2705` | `ErrStoreDisabled` | store operation while `vef.cron.store.enabled = false` |
| `2706` | `ErrScheduleInvalid(reason)` | non-trigger spec fault (name, window, timeout, params, policy vocabulary) |

## Configuration

```toml
[vef.cron.store]
enabled = false            # master switch; off touches no tables
auto_migrate = false       # run the cron DDL migration on start
poll_interval = "5s"       # schedule-table re-read bound (visibility latency, not fire precision)
batch_size = 32            # schedules claimed per poll tick
max_concurrent = 16        # concurrent runs per node
misfire_threshold = "1m"   # how late a fire may start before the misfire policy applies
heartbeat_interval = "10s" # executor liveness cadence on running runs
abandoned_after = "1m"     # stale-heartbeat window; must be ≥ 2 × heartbeat_interval
run_timeout = "0s"         # default per-run bound; zero leaves runs unbounded
run_retention = "0s"       # journal retention; zero keeps rows forever
```

Validation rejects negative durations and an `abandoned_after` tighter than
twice the heartbeat interval (healthy executors would be declared dead) at
startup.

## Next Step

The in-memory scheduler for process-local work is documented in
[Cron Jobs](./cron). For alerting on failed or abandoned runs, subscribe to
the events above through the [Event Bus](./event-bus).
