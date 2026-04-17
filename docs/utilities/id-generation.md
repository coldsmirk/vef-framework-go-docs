---
sidebar_position: 1
---

# ID Generation

The `id` package provides pluggable unique identifier generation with two built-in strategies.

## Quick Start

```go
import "github.com/coldsmirk/vef-framework-go/id"

// Default: XID (recommended for most use cases)
xid := id.Generate()
// → "9m4e2mr0ui3e8a215n4g" (20 chars, base32)

// UUID v7 (when RFC 4122 compliance is needed)
uuid := id.GenerateUUID()
// → "018f4e42-832a-7123-9abc-def012345678" (36 chars)
```

## Built-in Generators

### XID (Default)

XID is the framework's default ID generation strategy, chosen for its balance of performance and uniqueness.

| Property | Value |
| --- | --- |
| Format | 20-character base32 string (`0-9, a-v`) |
| Sortable | ✅ Time-ordered |
| Globally unique | ✅ Machine ID + counter |
| Performance | Best among all strategies |

```go
xid := id.Generate()
```

### UUID v7

UUID v7 provides time-based ordering and follows RFC 4122 standards.

| Property | Value |
| --- | --- |
| Format | 36-character UUID (`xxxxxxxx-xxxx-7xxx-xxxx-xxxxxxxxxxxx`) |
| Sortable | ✅ Time-ordered |
| RFC compliant | ✅ RFC 4122 |
| Use case | When external systems require UUIDs |

```go
uuid := id.GenerateUUID()
```

## IDGenerator Interface

You can implement custom ID generators by implementing the `IDGenerator` interface:

```go
type IDGenerator interface {
    Generate() string
}
```

The framework uses the default XID generator (`id.DefaultXIDGenerator`) for model primary keys. The `orm` package automatically calls `id.Generate()` when inserting records with empty IDs.

## Pre-built Generator Instances

```go
id.DefaultXIDGenerator  // *XIDGenerator singleton
id.DefaultUUIDGenerator // *UUIDGenerator singleton
```

## When to Use Which

| Scenario | Recommendation |
| --- | --- |
| General application IDs | `id.Generate()` (XID) |
| External API integration | `id.GenerateUUID()` (UUID v7) |
| Custom format needed | Implement `IDGenerator` |
