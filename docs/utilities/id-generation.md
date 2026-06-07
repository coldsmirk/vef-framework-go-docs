---
sidebar_position: 1
---

# ID Generation

The `id` package provides pluggable unique identifier generation. The framework ships four built-in strategies — XID, UUID v7, Snowflake, and random/Nano-style IDs.

## Reviewed Public Surface

The current source audit for `github.com/coldsmirk/vef-framework-go/id`
covers 15 top-level exported symbols, no exported fields, and 1 exported
method. The reviewed public-surface fingerprint is
`e9c002ee81d48b44c4f3a4dce5ebaf83f0a5c8d9f9dc2aa7885e94e1d325f79f`.

Reviewed APIs:

| API | Contract |
| --- | --- |
| `id.IDGenerator` | Interface implemented by every built-in generator |
| `IDGenerator.Generate()` | Returns the next ID as a string; concrete format depends on the generator |
| `id.Generate()` | Delegates to `DefaultXIDGenerator.Generate()` and returns a 20-character XID |
| `id.GenerateUUID()` | Delegates to `DefaultUUIDGenerator.Generate()` and returns a UUID v7 string |
| `id.DefaultXIDGenerator` | Package-level XID singleton created by `NewXIDGenerator()` |
| `id.DefaultUUIDGenerator` | Package-level UUID v7 singleton created by `NewUUIDGenerator()` |
| `id.DefaultSnowflakeIDGenerator` | Package-level Snowflake singleton initialized from `VEF_NODE_ID`, or node `0` when unset |
| `id.NewXIDGenerator()` | Returns an `IDGenerator` that wraps `xid.New().String()` |
| `id.NewUUIDGenerator()` | Returns an `IDGenerator` that uses `uuid.NewV7()` and panics if UUID creation fails |
| `id.NewSnowflakeIDGenerator(nodeID)` | Returns a Base36 Snowflake generator for node IDs `0..63`; invalid node IDs return an error |
| `id.NewRandomIDGenerator(opts...)` | Returns a random/Nano-style generator with defaults, then applies options in order |
| `id.RandomIDGeneratorOption` | Function option type used by the random generator constructor |
| `id.WithAlphabet(alphabet)` | Sets the random generator alphabet |
| `id.WithLength(length)` | Sets the random generator output length |
| `id.DefaultRandomIDGeneratorAlphabet` | Default random alphabet, `0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ` |
| `id.DefaultRandomIDGeneratorLength` | Default random output length, `32` |

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
// or
xid = id.NewXIDGenerator().Generate()
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
// or
uuid = id.NewUUIDGenerator().Generate()
```

### Snowflake

Twitter-style Snowflake IDs — 64-bit IDs encoded as Base36 strings. Use these
when you need ordered, distributed IDs.

| Property | Value |
| --- | --- |
| Encoding | Base36 string from a custom Snowflake layout: 6 node bits (0-63 nodes), 12 step bits (4096 IDs/ms/node) |
| Epoch | `1754582400000` (custom epoch baked into the package) |
| Node ID | Read from the `VEF_NODE_ID` environment variable at startup; defaults to `0` when unset |
| Default instance | `id.DefaultSnowflakeIDGenerator` |

```go
snow := id.DefaultSnowflakeIDGenerator.Generate()
// → Base36 string
```

For a custom node ID, build a fresh generator:

```go
gen, err := id.NewSnowflakeIDGenerator(int64(42))
if err != nil {
    return err
}
sid := gen.Generate()
```

`NewSnowflakeIDGenerator` returns an error for node IDs outside `0..63`, including
negative values. During package initialization, `VEF_NODE_ID` is parsed as an
integer; invalid values panic at startup.

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

`RandomIDGeneratorOption` is the option type used by `WithAlphabet(...)` and
`WithLength(...)`.

Options are applied in the order passed to `NewRandomIDGenerator(...)`. The
constructor does not validate custom alphabets or lengths; generation uses
`go-nanoid/v2` `MustGenerate`, so an empty alphabet or zero length panics when
`Generate()` is called.

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
