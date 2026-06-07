---
sidebar_position: 6
---

# Mapx

The `mapx` package provides bidirectional conversion between Go structs and `map[string]any`, built on top of `github.com/go-viper/mapstructure/v2`. VEF overrides the upstream default tag — the framework uses `json` tags by default.

## Reviewed Public Surface

The current source audit for `github.com/coldsmirk/vef-framework-go/mapx`
covers 26 top-level exported symbols, no exported fields, and no exported
methods. The reviewed public-surface fingerprint is
`dbf62ac9ffca276278fe475c4abe86ca386784ad1ddbc6fb4a91b3fe67491ee5`.

Reviewed APIs:

| API | Contract |
| --- | --- |
| `mapx.DecoderHook` | Exported composed `mapstructure.DecodeHookFunc` used by default decoders |
| `mapx.DecoderOption` | Function option type that mutates `mapstructure.DecoderConfig` |
| `mapx.Metadata` | Alias for `mapstructure.Metadata` |
| `mapx.NewDecoder(result, options...)` | Creates a `mapstructure.Decoder` with VEF defaults, then applies options in order |
| `mapx.ToMap(value, options...)` | Converts a struct or pointer-to-struct into `map[string]any`; non-struct input returns `ErrInvalidToMapValue` |
| `mapx.FromMap[T](value, options...)` | Converts `map[string]any` into `*T`; non-struct `T` returns `ErrInvalidFromMapType` |
| `mapx.WithTagName(tagName)` | Sets `DecoderConfig.TagName`; default is `json` |
| `mapx.WithIgnoreUntaggedFields(ignore)` | Sets `DecoderConfig.IgnoreUntaggedFields` to the supplied boolean |
| `mapx.WithDecodeHook(decodeHook)` | Replaces `DecoderConfig.DecodeHook`; compose with `mapx.DecoderHook` yourself to preserve defaults |
| `mapx.WithMatchName(matchName)` | Replaces the key/field matcher; the default is `mapKey == lo.CamelCase(fieldName)` |
| `mapx.WithErrorUnused()` | Sets `ErrorUnused = true` |
| `mapx.WithErrorUnset()` | Sets `ErrorUnset = true` |
| `mapx.WithZeroFields()` | Sets `ZeroFields = true` |
| `mapx.WithAllowUnsetPointer()` | Sets `AllowUnsetPointer = true` |
| `mapx.WithMetadata(metadata)` | Stores decode metadata in the supplied `*mapx.Metadata` |
| `mapx.WithWeaklyTypedInput()` | Sets `WeaklyTypedInput = true` |
| `mapx.WithDecodeNil()` | Sets `DecodeNil = true` |
| `mapx.ErrInvalidToMapValue` | Sentinel for `ToMap` input that is not a struct or pointer-to-struct |
| `mapx.ErrInvalidFromMapType` | Sentinel for `FromMap[T]` when `T` is not a struct |
| `mapx.ErrCollectionSetNilElement` | Sentinel for nil elements while decoding into collection sets |
| `mapx.ErrCollectionSetIncompatibleKind` | Sentinel for string/numeric family mismatches in collection set elements |
| `mapx.ErrCollectionSetOverflow` | Sentinel for numeric overflow while converting collection set elements |
| `mapx.ErrCollectionSetNonInteger` | Sentinel for fractional float values targeting integer set elements |
| `mapx.ErrCollectionSetNotFinite` | Sentinel for NaN or infinity targeting integer set elements |
| `mapx.ErrCollectionSetNegative` | Sentinel for negative values targeting unsigned set elements |
| `mapx.ErrCollectionSetUnsupportedTarget` | Sentinel for collection set element kinds without a conversion strategy |

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
| `WithIgnoreUntaggedFields(ignore)` | Set whether fields without the active tag are ignored. |
| `WithDecodeHook(hook)` | Replace the default decode hook. |
| `WithMatchName(fn)` | Custom field-name matcher (default: exact match against `lo.CamelCase(fieldName)`). |
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
- `json.RawMessage` — marshals the source value to JSON bytes
- `*multipart.FileHeader` — picks the only entry when the source is `[]*multipart.FileHeader` with length 1
- `collections.Set` / `SortedSet` / `ConcurrentSet` / `ConcurrentSortedSet` — turns a slice or array into the corresponding set type
- `encoding.TextUnmarshaler` — any type that implements `UnmarshalText`
- string → primitive coercions (int / uint / float / bool)

Collection-set decoding is registered for `string`, signed integers, unsigned
integers, `float32`, and `float64`. It rejects nil elements, string/numeric
family mismatches, numeric overflow, fractional floats targeting integer sets,
NaN or infinity targeting integer sets, and negative values targeting unsigned
sets.

`timex.DateTime` / `timex.Date` / `timex.Time` are defined as named types over `time.Time`; whether they hit the `time.Time` hook depends on mapstructure's underlying-type unwrapping. Verify case by case if you rely on automatic decoding for those types.

`WithDecodeHook(myHook)` replaces the default composed hook. To extend the
defaults, compose your hook with `mapx.DecoderHook` before passing it to
`WithDecodeHook`.

The composed default hook is also exported as `mapx.DecoderHook`, and metadata
collection uses the exported alias `mapx.Metadata`.

## Error Sentinels

| Error | Meaning |
| --- | --- |
| `ErrInvalidToMapValue` | `ToMap` received a non-struct value |
| `ErrInvalidFromMapType` | `FromMap[T]` was instantiated with a non-struct `T` |
| `ErrCollectionSetNilElement` | a nil element cannot be inserted into a collection set |
| `ErrCollectionSetIncompatibleKind` | source value kind does not match the set element kind |
| `ErrCollectionSetOverflow` | numeric source value overflows the target set element type |
| `ErrCollectionSetNonInteger` | fractional float would lose data when decoded into an integer set |
| `ErrCollectionSetNotFinite` | NaN or infinity cannot be decoded into an integer set |
| `ErrCollectionSetNegative` | negative value cannot decode into an unsigned set element |
| `ErrCollectionSetUnsupportedTarget` | target set element kind has no conversion strategy |
