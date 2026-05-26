---
sidebar_position: 6
---

# Mapx

The `mapx` package provides bidirectional conversion between Go structs and `map[string]any`, built on top of `github.com/go-viper/mapstructure/v2`. VEF overrides the upstream default tag — the framework uses `json` tags by default.

## Struct to Map

```go
import "github.com/coldsmirk/vef-framework-go/mapx"

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

user := User{Name: "Alice", Email: "alice@example.com", Age: 30}
m, err := mapx.ToMap(user)
// m = map[string]any{"name": "Alice", "email": "alice@example.com", "age": 30}
```

## Map to Struct

```go
data := map[string]any{
    "name":  "Bob",
    "email": "bob@example.com",
    "age":   25,
}

user, err := mapx.FromMap[User](data)
// user.Name = "Bob", user.Email = "bob@example.com", user.Age = 25
```

## Decoder Options

Both `ToMap` and `FromMap` accept optional `DecoderOption` values:

```go
// Switch to a different tag, e.g. yaml
m, err := mapx.ToMap(user, mapx.WithTagName("yaml"))

// Weak type conversion (string "123" → int 123)
user, err := mapx.FromMap[User](data, mapx.WithWeaklyTypedInput())

// Surface fields present in the source map but absent from the struct
user, err := mapx.FromMap[User](data, mapx.WithErrorUnused())
```

### Available Options

| Option | Effect |
| --- | --- |
| `WithTagName(tag)` | Override the struct tag mapx reads (**default: `json`**). |
| `WithIgnoreUntaggedFields()` | Skip fields that don't carry the active tag. |
| `WithDecodeHook(hooks...)` | Append extra decode hooks (defaults stay in place). |
| `WithMatchName(fn)` | Custom field-name matcher (default: case-insensitive camelCase compare). |
| `WithErrorUnused()` | Fail when the source map carries keys not present on the struct. |
| `WithErrorUnset()` | Fail when the struct has fields the source map didn't populate. |
| `WithZeroFields()` | Zero out target struct fields before decoding. |
| `WithAllowUnsetPointer()` | Allow pointer fields to remain nil instead of being initialized. |
| `WithMetadata(m)` | Collect "unused" / "unset" key lists into a `mapstructure.Metadata` value. |
| `WithWeaklyTypedInput()` | Coerce common type mismatches (string ↔ number ↔ bool …). |
| `WithDecodeNil()` | Pass `nil` source values into the decode pipeline instead of skipping them. |

## Custom Decoder

For advanced use cases, create a reusable decoder:

```go
var result User
decoder, err := mapx.NewDecoder(&result, mapx.WithTagName("yaml"))
if err != nil {
    return err
}
err = decoder.Decode(data)
```

## Decode Hooks

`mapx` ships a rich set of decode hooks pre-registered on `NewDecoder`, so plain-text maps coming from JSON, form data, or environment configs decode into typed structs without per-field wiring:

- `time.Time` — parses `"2006-01-02 15:04:05"` (Go's `time.DateTime` layout)
- `time.Location` — parses IANA names (e.g. `"Asia/Shanghai"`)
- `time.Duration` — parses Go duration strings (e.g. `"5m"`)
- `*url.URL` — parses URLs
- `net.IP` / `net.IPNet` / `netip.Addr` / `netip.AddrPort` / `netip.Prefix`
- `json.RawMessage` — passes the raw value through verbatim
- `*multipart.FileHeader` — picks the first entry when the source is `[]*multipart.FileHeader` (so single-file uploaded fields work seamlessly)
- `collections.Set` / `SortedSet` / `ConcurrentSet` / `ConcurrentSortedSet` — turns a slice into the corresponding set type
- `encoding.TextUnmarshaler` — any type that implements `UnmarshalText`
- string → primitive coercions (int / uint / float / bool)

`timex.DateTime` / `timex.Date` / `timex.Time` are defined as named types over `time.Time`; whether they hit the `time.Time` hook depends on mapstructure's underlying-type unwrapping. Verify case by case if you rely on automatic decoding for those types.

If you need to register your own hooks, append them with `mapx.WithDecodeHook(myHook)` — VEF's built-in hooks stay in place.
