---
sidebar_position: 3
---

# Cron Jobs

VEF exposes `github.com/coldsmirk/vef-framework-go/cron` as a typed wrapper around `gocron`. The framework DI module provides `cron.Scheduler` automatically; applications that need custom wiring can pass their own `gocron.Scheduler` to `cron.NewScheduler`.

## Public Surface

The package has no exported fields. Its public top-level API is:

| API | Contract |
| --- | --- |
| `cron.Scheduler` | high-level scheduler interface used by DI, custom handlers, and modules |
| `cron.NewScheduler(scheduler gocron.Scheduler) cron.Scheduler` | wraps a caller-provided `gocron.Scheduler`; the caller must pass a usable scheduler |
| `cron.Job` | interface returned by `Scheduler.NewJob`, `Scheduler.Jobs`, and `Scheduler.Update` |
| `cron.JobDefinition` | interface accepted by `Scheduler.NewJob` and `Scheduler.Update` |
| `cron.JobDescriptorOption` | option type accepted by all built-in job-definition constructors |
| `cron.OneTimeJobDefinition` | one-time job definition returned by `NewOneTimeJob` |
| `cron.DurationJobDefinition` | fixed-interval job definition returned by `NewDurationJob` |
| `cron.DurationRandomJobDefinition` | random-interval job definition returned by `NewDurationRandomJob` |
| `cron.CronJobDefinition` | cron-expression job definition returned by `NewCronJob` |
| `cron.NewOneTimeJob(times []time.Time, options ...cron.JobDescriptorOption) *cron.OneTimeJobDefinition` | runs immediately when `times` is nil/empty, once at one time, or once at each supplied time |
| `cron.NewDurationJob(interval time.Duration, options ...cron.JobDescriptorOption) *cron.DurationJobDefinition` | repeats at a fixed duration interval |
| `cron.NewDurationRandomJob(minInterval time.Duration, maxInterval time.Duration, options ...cron.JobDescriptorOption) *cron.DurationRandomJobDefinition` | repeats with a random interval between `minInterval` and `maxInterval` |
| `cron.NewCronJob(expression string, withSeconds bool, options ...cron.JobDescriptorOption) *cron.CronJobDefinition` | uses a cron expression; `withSeconds=true` expects the seconds field, `false` uses the standard 5-field form |
| `cron.WithName(name string)` | sets the required human-readable job name |
| `cron.WithTags(tags ...string)` | assigns tags used by `RemoveByTags` and `Job.Tags` |
| `cron.WithConcurrent()` | allows overlapping executions of the same job; without it, jobs use singleton wait mode |
| `cron.WithStartAt(startAt time.Time)` | starts the schedule at a specific time; this takes precedence over `WithStartImmediately` |
| `cron.WithStartImmediately()` | starts the schedule immediately when no `WithStartAt` time is set |
| `cron.WithStopAt(stopAt time.Time)` | stops the schedule at a specific time |
| `cron.WithLimitedRuns(limitedRuns uint)` | applies a run limit only when `limitedRuns > 0` |
| `cron.WithContext(ctx context.Context)` | forwards a non-nil context to the underlying job option for cancellation support |
| `cron.WithTask(handler any, params ...any)` | sets the required function handler and forwards `params` to it through `gocron.NewTask` |
| `cron.ErrJobNameRequired` | job options were built without a name |
| `cron.ErrJobTaskHandlerRequired` | task build found a nil handler |
| `cron.ErrJobTaskHandlerMustFunc` | task build found a non-function handler |

## `cron.Scheduler`

The public scheduler interface includes:

| Method | Signature | Contract |
| --- | --- | --- |
| `Jobs` | `Jobs() []cron.Job` | returns the jobs currently registered with the wrapped scheduler |
| `NewJob` | `NewJob(definition cron.JobDefinition) (cron.Job, error)` | builds the definition, registers it, and returns a `cron.Job`; validation errors from the definition are returned before registration |
| `RemoveByTags` | `RemoveByTags(tags ...string)` | removes all jobs that have any of the supplied tags |
| `RemoveJob` | `RemoveJob(id string) error` | parses `id` as a UUID and removes that job; invalid non-UUID strings return a parse error before delegation |
| `Start` | `Start()` | starts scheduling and execution; jobs added after start are scheduled by the wrapped scheduler |
| `StopJobs` | `StopJobs() error` | stops job execution without removing definitions; jobs can run again after `Start()` |
| `Update` | `Update(id string, definition cron.JobDefinition) (cron.Job, error)` | parses `id` as a UUID, builds the replacement definition, and preserves the job identifier |
| `JobsWaitingInQueue` | `JobsWaitingInQueue() int` | returns the wrapped scheduler's waiting queue count; it is meaningful when wait mode is in use |

Use IDs returned by `Job.ID()` for `RemoveJob` and `Update`. Arbitrary strings are not accepted.

## `cron.Job`

Each registered job can be inspected or triggered through:

| Method | Signature | Contract |
| --- | --- | --- |
| `ID` | `ID() string` | returns the underlying UUID as a string |
| `LastRun` | `LastRun() (time.Time, error)` | returns the last run start time |
| `Name` | `Name() string` | returns the configured job name |
| `NextRun` | `NextRun() (time.Time, error)` | returns the next scheduled run time |
| `NextRuns` | `NextRuns(count int) ([]time.Time, error)` | returns `count` future scheduled run times |
| `RunNow` | `RunNow() error` | triggers an immediate execution while respecting job/scheduler limits and run limits |
| `Tags` | `Tags() []string` | returns the configured tag list |

## Job Definitions

`cron.JobDefinition` is the common interface accepted by `Scheduler.NewJob(...)` and `Scheduler.Update(...)`. The four exported definition structs are returned by their constructors; callers usually do not instantiate them directly because their scheduling fields are internal.

| Constructor | Schedule |
| --- | --- |
| `cron.NewOneTimeJob(nil, options...)` or `cron.NewOneTimeJob([]time.Time{}, options...)` | one immediate run |
| `cron.NewOneTimeJob([]time.Time{t}, options...)` | one run at `t` |
| `cron.NewOneTimeJob([]time.Time{a, b}, options...)` | one run at each supplied time |
| `cron.NewDurationJob(interval, options...)` | fixed interval |
| `cron.NewDurationRandomJob(minInterval, maxInterval, options...)` | random interval in the supplied range |
| `cron.NewCronJob(expression, true, options...)` | cron expression with seconds |
| `cron.NewCronJob(expression, false, options...)` | standard 5-field cron expression |

All built-in definitions first build the task and then build job options. That means a nil or non-function task handler takes precedence over a missing name if both are invalid.

## Job Options

Every constructor accepts `JobDescriptorOption` values:

| Option | Behavior |
| --- | --- |
| `cron.WithName(name)` | required; without it `NewJob`/`Update` return an error wrapping `ErrJobNameRequired` |
| `cron.WithTask(handler, params...)` | required; `handler` must be a function; `params` are forwarded to `gocron.NewTask(handler, params...)` |
| `cron.WithTags(tags...)` | stores tags for listing and `RemoveByTags` |
| `cron.WithConcurrent()` | disables the default singleton wait mode for this job |
| `cron.WithStartAt(startAt)` | sets a start date/time and wins over `WithStartImmediately` when both are present |
| `cron.WithStartImmediately()` | applies only when `WithStartAt` did not set a non-zero time |
| `cron.WithStopAt(stopAt)` | adds a stop date/time only when `stopAt` is non-zero |
| `cron.WithLimitedRuns(limitedRuns)` | adds a run limit only when `limitedRuns > 0` |
| `cron.WithContext(ctx)` | adds the context option only when `ctx` is non-nil; handlers that accept `context.Context` can observe cancellation |

By default, a job cannot overlap with another execution of itself. The package adds singleton wait mode unless `WithConcurrent()` is supplied.

## Error Sentinels

| Error | Trigger |
| --- | --- |
| `cron.ErrJobTaskHandlerRequired` | `WithTask` was not supplied or received nil |
| `cron.ErrJobTaskHandlerMustFunc` | `WithTask` received a non-function handler |
| `cron.ErrJobNameRequired` | task build succeeded, but `WithName` was missing or empty |

These sentinels are wrapped by the build path, for example `failed to build job task: ...` or `failed to build job options: ...`. Use `errors.Is(err, cron.ErrJobNameRequired)` and the other sentinels instead of direct equality.

## DI Scheduler Defaults

When the framework creates the scheduler through its DI module, the observable defaults are:

| Default | Meaning |
| --- | --- |
| local time zone | schedules use `time.Local` |
| stop timeout `30s` | graceful shutdown window |
| scheduler logger and monitor | job scheduling and completion are logged |
| concurrent job limit `1000` with wait mode | excess executions wait in the queue instead of being dropped |
| lifecycle start/shutdown | app start calls scheduler start; app stop shuts the scheduler down |

These defaults apply to the framework-provided scheduler. A scheduler passed to `cron.NewScheduler` uses whatever options the caller configured.

## Minimal Example

```go
package jobs

import (
  "context"
  "time"

  "github.com/coldsmirk/vef-framework-go/cron"
)

func RegisterCleanupJob(scheduler cron.Scheduler) error {
  _, err := scheduler.NewJob(
    cron.NewDurationJob(
      10*time.Minute,
      cron.WithName("cleanup-expired-sessions"),
      cron.WithTask(func(ctx context.Context) error {
        return nil
      }),
    ),
  )

  return err
}
```

## Practical Usage

Use cron when you need application-owned background work such as:

- periodic cleanup
- polling and sync tasks
- scheduled aggregation
- retry or timeout scanning

If the task is domain-heavy, keep business logic in a service and use the scheduler only for orchestration.

## Next Step

Read [Transactions](../advanced/transactions) if scheduled work needs explicit database transaction control.
