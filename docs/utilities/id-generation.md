---
sidebar_position: 1
---

# ID Generation

The `id` package provides pluggable unique identifier generation. The framework ships four built-in strategies — XID, UUID v7, Snowflake, and random/Nano-style IDs.

## Quick Start

```go
import "github.com/coldsmirk/vef-framework-go/id"

// XID (the default for model primary keys)
xid := id.Generate()
// → "9m4e2mr0ui3e8a215n4g" (20 chars, base32)

// UUID v7 (when RFC 4122 compliance is needed)
uuid := id.GenerateUUID()
// → "018f4e42-832a-7123-9abc-def012345678"
```

## Built-in Generators

### XID

XID is the framework's default for model primary keys — best balance of performance and uniqueness.

| Property | Value |
| --- | --- |
| Format | 20-character base32 string (`0-9, a-v`) |
| Sortable | Time-ordered |
| Globally unique | Machine ID + counter |
| Performance | Best among the four strategies |

```go
xid := id.Generate()
// or
xid := id.DefaultXIDGenerator.Generate()
```

### UUID v7

Time-based, RFC 4122-compliant UUIDs — use these when integrating with systems that expect canonical UUIDs.

| Property | Value |
| --- | --- |
| Format | 36-character UUID (`xxxxxxxx-xxxx-7xxx-xxxx-xxxxxxxxxxxx`) |
| Sortable | Time-ordered |
| RFC compliant | RFC 4122 |

```go
uuid := id.GenerateUUID()
// or
uuid := id.DefaultUUIDGenerator.Generate()
```

### Snowflake

Twitter-style Snowflake IDs — 64-bit integers encoded as decimal strings. Use these when you need ordered, distributed integer IDs.

| Property | Value |
| --- | --- |
| Encoding | Custom: 6 node bits (0-63 nodes), 12 step bits (4096 IDs/ms/node) |
| Epoch | `1754582400000` (custom epoch baked into the package) |
| Node ID | Read from the `VEF_NODE_ID` environment variable at startup |
| Default instance | `id.DefaultSnowflakeIDGenerator` |

```go
snow := id.DefaultSnowflakeIDGenerator.Generate()
// → "7234567890123456789"
```

For a custom node ID, build a fresh generator:

```go
gen, err := id.NewSnowflakeIDGenerator(int64(42))
if err != nil {
    return err
}
sid := gen.Generate()
```

> Snowflake supports up to 64 nodes and 4096 IDs per millisecond per node. Pin `VEF_NODE_ID` uniquely per process to avoid collisions.

### Random / Nano-style

Cryptographically random IDs with a configurable alphabet — useful for short, opaque tokens.

| Property | Value |
| --- | --- |
| Default alphabet | `0-9 a-z A-Z` (62 characters; `id.DefaultRandomIDGeneratorAlphabet`) |
| Default length | 32 (`id.DefaultRandomIDGeneratorLength`) |

```go
// Default 32-character alphanumeric token
gen := id.NewRandomIDGenerator()
token := gen.Generate()

// Custom: 16-char numeric-only
gen = id.NewRandomIDGenerator(
    id.WithAlphabet("0123456789"),
    id.WithLength(16),
)
```

## IDGenerator Interface

All built-in generators implement the same interface:

```go
type IDGenerator interface {
    Generate() string
}
```

Pre-built singletons:

```go
id.DefaultXIDGenerator        // IDGenerator
id.DefaultUUIDGenerator       // IDGenerator
id.DefaultSnowflakeIDGenerator // IDGenerator
```

The `orm` package uses `id.Generate()` (XID) automatically when inserting records with empty IDs.

## When to Use Which

| Scenario | Recommendation |
| --- | --- |
| General application IDs (primary keys) | `id.Generate()` (XID) |
| External APIs expecting UUIDs | `id.GenerateUUID()` |
| Distributed ordered integer IDs | `id.DefaultSnowflakeIDGenerator` |
| Short tokens / invites / shareable links | `id.NewRandomIDGenerator(...)` |
| Custom format | Implement `IDGenerator` yourself |
