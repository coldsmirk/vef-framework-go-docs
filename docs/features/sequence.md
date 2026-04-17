---
sidebar_position: 14
---

# Sequence

The `sequence` package provides configurable serial number generation with customizable formats and reset policies.

## Core Concepts

A sequence consists of:
- **Rule**: defines the format template and reset policy
- **Store**: persists the current counter value (database, Redis, or memory)

## Sequence Rule

```go
import "github.com/coldsmirk/vef-framework-go/sequence"

rule := sequence.Rule{
    Name:        "order-number",
    Format:      "ORD-{year}{month}{day}-{seq:6}",
    ResetPolicy: sequence.ResetDaily,
}
```

### Format Tokens

| Token | Description | Example |
| --- | --- | --- |
| `{year}` | 4-digit year | `2024` |
| `{month}` | 2-digit month | `03` |
| `{day}` | 2-digit day | `15` |
| `{hour}` | 2-digit hour | `14` |
| `{minute}` | 2-digit minute | `30` |
| `{second}` | 2-digit second | `05` |
| `{seq:N}` | Zero-padded sequence number (N digits) | `000001` |

### Reset Policies

| Policy | Constant | Behavior |
| --- | --- | --- |
| Never | `sequence.ResetNever` | Counter grows indefinitely |
| Daily | `sequence.ResetDaily` | Resets at midnight |
| Monthly | `sequence.ResetMonthly` | Resets on 1st of each month |
| Yearly | `sequence.ResetYearly` | Resets on Jan 1st |

## Stores

### Database Store

Uses the ORM database for persistence — best for most applications:

```go
store := sequence.NewDBStore(db)
```

### Redis Store

Uses Redis for persistence — best for high-throughput scenarios:

```go
store := sequence.NewRedisStore(redisClient)
```

### Memory Store

In-memory storage — for testing only:

```go
store := sequence.NewMemoryStore()
```

## Generating Sequences

```go
generator := sequence.New(rule, store)

// Generate next sequence number
number, err := generator.Next(ctx)
// → "ORD-20240315-000001"

number, err = generator.Next(ctx)
// → "ORD-20240315-000002"
```

## Example Formats

| Use Case | Format | Example Output |
| --- | --- | --- |
| Order Number | `ORD-{year}{month}{day}-{seq:6}` | `ORD-20240315-000001` |
| Invoice | `INV{year}{month}-{seq:4}` | `INV202403-0001` |
| Document | `DOC-{year}-{seq:8}` | `DOC-2024-00000001` |
| Simple Counter | `{seq:10}` | `0000000001` |
