---
sidebar_position: 1
---

# ID Generation

The `id` package provides pluggable unique identifier generation. The framework ships three built-in strategies — XID, UUID v7, and random/Nano-style IDs.

## API Reference

| API | Contract |
| --- | --- |
| `id.IDGenerator` | Interface implemented by every built-in generator |
| `IDGenerator.Generate()` | Returns the next ID as a string; concrete format depends on the generator |
| `id.Generate()` | Delegates to `DefaultXIDGenerator.Generate()` and returns a 20-character XID |
| `id.GenerateUUID()` | Delegates to `DefaultUUIDGenerator.Generate()` and returns a UUID v7 string |
| `id.DefaultXIDGenerator` | Package-level XID singleton created by `NewXIDGenerator()` |
| `id.DefaultUUIDGenerator` | Package-level UUID v7 singleton created by `NewUUIDGenerator()` |
| `id.NewXIDGenerator()` | Returns an `IDGenerator` that wraps `xid.New().String()` |
| `id.NewUUIDGenerator()` | Returns an `IDGenerator` that uses `uuid.NewV7()` and panics if UUID creation fails |
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
| Performance | Best among the three strategies |

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
id.DefaultXIDGenerator   // IDGenerator
id.DefaultUUIDGenerator  // IDGenerator
```

The `orm` package uses `id.Generate()` (XID) automatically when inserting records with empty IDs.

## When to Use Which

| Scenario | Recommendation |
| --- | --- |
| General application IDs (primary keys) | `id.Generate()` (XID) |
| External APIs expecting UUIDs | `id.GenerateUUID()` |
| Short tokens / invites / shareable links | `id.NewRandomIDGenerator(...)` |
| Custom format | Implement `IDGenerator` yourself |
