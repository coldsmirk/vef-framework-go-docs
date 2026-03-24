---
sidebar_position: 3
---

# Cron Jobs

VEF includes a scheduler abstraction on top of `gocron`.

## What The Module Provides

The cron module provides `cron.Scheduler` through dependency injection and starts the underlying scheduler through the application lifecycle.

## `cron.Scheduler` Interface

The public scheduler interface includes:

| Method | Purpose |
| --- | --- |
| `Jobs()` | list all registered jobs |
| `NewJob(definition)` | create and register a new job |
| `RemoveByTags(tags...)` | remove jobs by tags |
| `RemoveJob(id)` | remove one job by ID |
| `Start()` | start scheduling and execution |
| `StopJobs()` | stop job execution without removing definitions |
| `Update(id, definition)` | replace an existing job definition |
| `JobsWaitingInQueue()` | inspect the queue length when wait mode is used |

## `cron.Job` Interface

Each scheduled job can be inspected through:

| Method | Purpose |
| --- | --- |
| `ID()` | unique job identifier |
| `LastRun()` | most recent execution time |
| `Name()` | human-readable name |
| `NextRun()` | next scheduled execution time |
| `NextRuns(count)` | multiple future run times |
| `RunNow()` | trigger an immediate execution |
| `Tags()` | list job tags |

## Supported Job Definitions

The public package currently supports these job-definition constructors:

| Constructor | Meaning |
| --- | --- |
| `cron.NewOneTimeJob(times, options...)` | run once immediately, once at a time, or once at multiple specified times |
| `cron.NewDurationJob(interval, options...)` | run repeatedly at a fixed interval |
| `cron.NewDurationRandomJob(minInterval, maxInterval, options...)` | run repeatedly at random intervals inside a range |
| `cron.NewCronJob(expression, withSeconds, options...)` | run using a cron expression |

## Job Descriptor Options

Jobs can be customized with these options:

| Option | Purpose |
| --- | --- |
| `cron.WithName(name)` | sets the job display name |
| `cron.WithTags(tags...)` | attaches tags for grouping and bulk removal |
| `cron.WithConcurrent()` | allows concurrent self-execution |
| `cron.WithStartAt(time)` | delays schedule start until a specific time |
| `cron.WithStartImmediately()` | starts the schedule immediately |
| `cron.WithStopAt(time)` | stops the schedule at a specific time |
| `cron.WithLimitedRuns(count)` | limits how many times the job may run |
| `cron.WithContext(ctx)` | associates a context for cancellation |
| `cron.WithTask(handler, params...)` | sets the actual task function and its parameters |

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

## Runtime Defaults

The internal scheduler is created with these runtime characteristics:

| Default | Meaning |
| --- | --- |
| local time zone | schedules use `time.Local` |
| stop timeout `30s` | graceful shutdown window |
| concurrent job limit `1000` | queueing mode waits instead of dropping |

These are implementation defaults, not application-level policy guarantees.

## Practical Usage

Use cron when you need application-owned background work such as:

- periodic cleanup
- polling and sync tasks
- scheduled aggregation
- retry or timeout scanning

If the task is domain-heavy, keep business logic in a service and use the scheduler only for orchestration.

## Next Step

Read [Transactions](../advanced/transactions) if scheduled work needs explicit database transaction control.
