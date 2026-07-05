---
sidebar_position: 14
---

# Sequence

The `sequence` package generates serial numbers (order numbers, invoice numbers, etc.) with configurable prefixes, date parts, zero-padded counters, and automatic reset policies.

## Concepts

A serial number generator consists of two pieces:

- **`sequence.Rule`** — the format and reset policy (prefix, date format, counter width, reset cycle, overflow strategy).
- **`sequence.Store`** — where the rule + current counter live. The built-in
  default runtime store is in-memory, and the public package also exposes
  database and Redis store implementations for durable or distributed
  deployments. Custom stores can implement the same interface for another
  persistence model. The store is responsible for atomic counter increments.

The framework wires a `sequence.Generator` for you. It also exposes the concrete
`*sequence.MemoryStore`, so business modules can seed rules during startup and
then call `Generate(ctx, key)`.

The helper `sequence.FormatDate(dt, format)` is public as well. It renders the
same `yyyy` / `MM` / `dd` / `HH` / `mm` / `ss` tokens used by `Rule.DateFormat`.

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
| `CurrentValue` | Last reserved counter value. Stores update this cursor; callers normally do not mutate it directly. |
| `LastResetAt` | Timestamp used to decide whether the next reservation crosses a reset boundary. `nil` means the rule has never reset. |
| `IsActive` | Inactive rules return `sequence.ErrRuleNotFound` from `Generate`. |

`Rule.Clone()` returns a deep copy of the rule snapshot, including a copied
`LastResetAt` pointer when it is set.

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

An empty `ResetCycle` behaves like `sequence.ResetNone`. Any value other than
the exported constants is treated defensively as "do not reset". For non-none
cycles, `LastResetAt == nil` triggers a reset on the next reservation.

### Overflow strategies

| Constant | Behavior when `MaxValue` is exceeded |
| --- | --- |
| `sequence.OverflowError` *(default)* | Return `sequence.ErrSequenceOverflow` and refuse to generate further numbers until the next reset. |
| `sequence.OverflowReset` | Reset the counter to `StartValue` and continue. |
| `sequence.OverflowExtend` | Keep counting past `SeqLength` (the result simply gets more digits). |

`MaxValue` is checked after applying any cycle reset. If a reset boundary moves
the counter back to `StartValue` but the requested batch still exceeds
`MaxValue`, even `OverflowReset` returns `sequence.ErrSequenceOverflow` because
resetting again cannot make the batch fit. Values other than the exported
`OverflowStrategy` constants fall back to `OverflowError`.

## Store

The current built-in runtime store is in-memory and non-durable: counters and
registered rules are lost on process restart. It is suitable for tests, dev, and
single-process deployments. Distributed or durable deployments should provide a
custom `sequence.Store`.

`sequence.Store.Reserve(ctx, key, count, now)` is the contract boundary for
custom stores. Implementations must serialize the read-modify-write path per
rule key and reserve the whole `count` batch atomically.

### In-memory

For tests, dev, single-process deployments. Rules must be registered up front via `Register`:

```go
store := sequence.NewMemoryStore()
store.Register(rule)
```

`Register` overwrites any existing rule with the same `Key` and stores a deep
copy, so later mutations to the original `*Rule` do not change the store.

Inside a VEF app, inject the concrete `*sequence.MemoryStore` when you want to seed rules:

```go
func SeedSequenceRules(store *sequence.MemoryStore) {
    store.Register(rule)
}
```

### Database

`sequence.NewDBStore(db)` returns a `*sequence.DBStore` backed by the
`sys_sequence_rule` table (`sequence.DBStoreTableName`). `DBStore.Init(ctx)`
creates the table when it does not exist, and `Reserve(...)` locks the rule row
for update so each reservation is atomic within the database transaction.
`sequence.RuleModel` is the ORM model for that table.

### Redis

`sequence.NewRedisStore(client)` returns a `*sequence.RedisStore` for
distributed deployments. Rules are stored as Redis hashes under the
`vef:sequence:<key>` prefix; `RedisStore.RegisterRule(ctx, rule)` seeds or
replaces one rule, and `Reserve(...)` uses Redis `WATCH`/transaction retry to
reserve counters atomically.

Every store **expects rules to exist before `Generate(...)` is called**. Calling
`Generate(ctx, "unknown-key")` on a store that doesn't know the rule returns
`sequence.ErrRuleNotFound`.

The public store API is intentionally small:

| API | Purpose |
| --- | --- |
| `sequence.Store.Reserve(ctx, key, count, now)` | atomically reserve `count` numbers for a rule and return the rule snapshot plus the final counter value |
| `sequence.MemoryStore.Register(rules...)` | preload or replace in-memory rules using deep copies |
| `sequence.MemoryStore.Reserve(...)` | in-memory implementation of `Store.Reserve`; returns cloned rule snapshots |
| `sequence.NewDBStore(db)` / `sequence.DBStore` | database-backed store using `sequence.DBStoreTableName` (`sys_sequence_rule`) and `sequence.RuleModel` |
| `sequence.NewRedisStore(client)` / `sequence.RedisStore` | Redis-backed store with `RedisStore.RegisterRule(ctx, rule)` for seeding hash-backed rules |
| `sequence.Rule.Clone()` | deep-copy a rule snapshot |

## Generating numbers

The framework already wires a `sequence.Generator` from the active `Store`.
Inject it and call:

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

`GenerateN` reserves the full range in one atomic store operation. Returned
numbers are ordered from first to last reserved value; when `SeqStep > 1`, the
values are spaced by that step (`0002`, `0004`, `0006`, ...).

`SeqLength` is a minimum zero-pad width, not a maximum. If the numeric value has
more digits than `SeqLength`, it is rendered in full instead of being truncated.

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
| `sequence.ErrInvalidCount` | `GenerateN` was called with a count lower than 1. |

## Next step

Read [Cache](./cache) or [Cron](./cron) — sequence generators are often paired with scheduled archival or warm-up jobs.
