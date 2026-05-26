---
sidebar_position: 14
---

# Sequence

The `sequence` package generates serial numbers (order numbers, invoice numbers, etc.) with configurable prefixes, date parts, zero-padded counters, and automatic reset policies.

## Concepts

A serial number generator consists of two pieces:

- **`sequence.Rule`** — the format and reset policy (prefix, date format, counter width, reset cycle, overflow strategy).
- **`sequence.Store`** — where the rule + current counter live (database, Redis, or in-memory). The store is responsible for atomic counter increments.

The framework wires a `sequence.Generator` for you; business code only needs to register rules in the store and call `Generate(ctx, key)`.

## Defining a Rule

```go
import (
    "github.com/coldsmirk/vef-framework-go/sequence"
)

rule := &sequence.Rule{
    Key:              "order-number",     // unique lookup key used by Generate(...)
    Name:             "Order number",
    Prefix:           "ORD-",
    DateFormat:       "yyyyMMdd-",        // optional; uses the date layout tokens below
    SeqLength:        6,                  // zero-pad to 6 digits → 000001
    SeqStep:          1,
    StartValue:       0,                  // first generated value = StartValue + SeqStep
    MaxValue:         0,                  // 0 = unlimited
    OverflowStrategy: sequence.OverflowError,
    ResetCycle:       sequence.ResetDaily,
    IsActive:         true,
}
```

### Rule fields

| Field | Meaning |
| --- | --- |
| `Key` | Lookup key passed to `Generate(ctx, key)`. Unique per store. |
| `Name` | Human-readable name (for admin UI). |
| `Prefix` / `Suffix` | Optional fixed text before / after the date+counter. |
| `DateFormat` | Optional date layout (token list below); empty = no date part. |
| `SeqLength` | Zero-padded counter width. |
| `SeqStep` | Counter increment per generation (usually 1). |
| `StartValue` | Counter value after a reset. The first generated number is `StartValue + SeqStep`. |
| `MaxValue` | Upper bound. `0` means unlimited. |
| `OverflowStrategy` | What to do when `MaxValue` is reached. See below. |
| `ResetCycle` | When the counter resets. See below. |
| `IsActive` | Inactive rules return `sequence.ErrRuleNotFound` from `Generate`. |

### Date layout tokens

`DateFormat` uses Java/.NET-style date tokens (translated internally to Go layout):

| Token | Meaning | Example |
| --- | --- | --- |
| `yyyy` | 4-digit year | `2024` |
| `yy` | 2-digit year | `24` |
| `MM` | 2-digit month | `03` |
| `dd` | 2-digit day | `15` |
| `HH` | 2-digit hour (24h) | `14` |
| `mm` | 2-digit minute | `30` |
| `ss` | 2-digit second | `05` |

Any other character passes through verbatim, so `yyyyMMdd-` produces `20240315-`.

### Reset cycles

| Constant | Cycle |
| --- | --- |
| `sequence.ResetNone` | Never reset |
| `sequence.ResetDaily` | Reset at the start of each day |
| `sequence.ResetWeekly` | Reset at the start of each week |
| `sequence.ResetMonthly` | Reset on the 1st of each month |
| `sequence.ResetQuarterly` | Reset on the first day of each calendar quarter |
| `sequence.ResetYearly` | Reset on January 1st |

### Overflow strategies

| Constant | Behavior when `MaxValue` is exceeded |
| --- | --- |
| `sequence.OverflowError` *(default)* | Return `sequence.ErrSequenceOverflow` and refuse to generate further numbers until the next reset. |
| `sequence.OverflowReset` | Reset the counter to `StartValue` and continue. |
| `sequence.OverflowExtend` | Keep counting past `SeqLength` (the result simply gets more digits). |

## Stores

Pick one based on your deployment topology:

### In-memory

For tests, dev, single-process deployments. Rules must be registered up front via `Register`:

```go
store := sequence.NewMemoryStore().(*sequence.MemoryStore)
store.Register(rule)
```

### Database

Persists rules and the counter in `sys_sequence_rule`:

```go
store := sequence.NewDBStore(db)
```

Each rule is one row in `sys_sequence_rule`. Seed rules through migrations or admin endpoints.

### Redis

Stores the counter in Redis for low-latency, distributed deployments:

```go
store := sequence.NewRedisStore(redisClient)
```

> Every store **expects rules to exist before `Generate(...)` is called**. Calling `Generate(ctx, "unknown-key")` on a store that doesn't know the rule returns `sequence.ErrRuleNotFound`.

## Generating numbers

The framework already wires a `sequence.Generator` from whichever `Store` you provide to FX. Inject it and call:

```go
type OrderService struct {
    seq sequence.Generator
}

func (s *OrderService) NewOrder(ctx context.Context) (string, error) {
    return s.seq.Generate(ctx, "order-number")
}
```

For batches (e.g. allocating 100 numbers atomically):

```go
numbers, err := seq.GenerateN(ctx, "order-number", 100)
```

`GenerateN` reserves the full range in one atomic store operation, so the returned numbers are guaranteed contiguous.

## Example rules

| Use case | Configuration | Sample output |
| --- | --- | --- |
| Order number | `Prefix:"ORD-"`, `DateFormat:"yyyyMMdd-"`, `SeqLength:6`, `ResetDaily` | `ORD-20240315-000001` |
| Invoice | `Prefix:"INV"`, `DateFormat:"yyyyMM-"`, `SeqLength:4`, `ResetMonthly` | `INV202403-0001` |
| Document ID | `Prefix:"DOC-"`, `DateFormat:"yyyy-"`, `SeqLength:8`, `ResetYearly` | `DOC-2024-00000001` |
| Simple counter | `SeqLength:10`, `ResetNone` | `0000000001` |

## Errors

| Error | Cause |
| --- | --- |
| `sequence.ErrRuleNotFound` | Key not registered in the store, or the rule has `IsActive=false`. |
| `sequence.ErrSequenceOverflow` | Counter reached `MaxValue` and `OverflowStrategy` is `OverflowError`. |
| `sequence.ErrInvalidRule` | Rule config is internally inconsistent (e.g. `SeqStep <= 0`). |

## Next step

Read [Cache](./cache) or [Cron](./cron) — sequence generators are often paired with scheduled archival or warm-up jobs.
