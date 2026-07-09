---
sidebar_position: 4
---

# Copier

The `copier` package provides struct-to-struct field copying with built-in type converters for common VEF types.

## API Reference

| API | Contract |
| --- | --- |
| `copier.Copy(src, dst, options...)` | Copies from `src` into `dst` by calling `copier.CopyWithOption(dst, src, opt)` from `github.com/jinzhu/copier`; `dst` must be a pointer destination or the underlying copier returns an error |
| `copier.CopyOption` | Function option type that mutates the underlying option struct from `github.com/jinzhu/copier` |
| `copier.TypeConverter` | Alias for `github.com/jinzhu/copier.TypeConverter` |
| `copier.FieldNameMapping` | Alias for `github.com/jinzhu/copier.FieldNameMapping` |
| `copier.WithIgnoreEmpty()` | Sets `IgnoreEmpty = true` |
| `copier.WithDeepCopy()` | Sets `DeepCopy = true` |
| `copier.WithCaseInsensitive()` | Sets `CaseSensitive = false`; default copying remains case-sensitive |
| `copier.WithFieldNameMapping(...)` | Appends mappings to `FieldNameMapping` |
| `copier.WithTypeConverters(...)` | Appends custom converters after the built-in converters rather than replacing them |

## Quick Start

```go
import "github.com/coldsmirk/vef-framework-go/copier"

type UserParams struct {
    Username string
    Email    string
}

type User struct {
    Username string
    Email    string
    IsActive bool
}

params := UserParams{Username: "alice", Email: "alice@example.com"}
user := User{}

err := copier.Copy(params, &user)
// user.Username = "alice", user.Email = "alice@example.com"
```

## Options

All copier options use the public `CopyOption` type.

`Copy(...)` is case-sensitive by default and always starts with the built-in
converter list. Options are applied in the order they are passed.
`WithTypeConverters(...)` appends custom converters instead of replacing those
defaults. `WithFieldNameMapping(...)` also appends mappings.

### WithIgnoreEmpty

Skips copying fields with zero values. This is how the framework's `Update` CRUD operation merges partial updates:

```go
err := copier.Copy(params, &user, copier.WithIgnoreEmpty())
```

### WithDeepCopy

Enables deep copying of nested structures:

```go
err := copier.Copy(src, &dst, copier.WithDeepCopy())
```

### WithCaseInsensitive

Enables case-insensitive field name matching:

```go
err := copier.Copy(src, &dst, copier.WithCaseInsensitive())
```

### WithFieldNameMapping

Adds custom field name mappings for fields with different names:

```go
err := copier.Copy(src, &dst, copier.WithFieldNameMapping(
    copier.FieldNameMapping{
        SrcType: UserParams{},
        DstType: User{},
        Mapping: map[string]string{
            "Name": "Username",
        },
    },
))
```

### WithTypeConverters

Adds custom type converters:

```go
err := copier.Copy(src, &dst, copier.WithTypeConverters(
    copier.TypeConverter{
        SrcType: MyCustomType{},
        DstType: "",
        Fn: func(src interface{}) (interface{}, error) {
            return src.(MyCustomType).String(), nil
        },
    },
))
```

## Built-in Type Converters

The copier includes automatic converters for value ↔ pointer conversions of all common types:

| Type Pair | Direction |
| --- | --- |
| `string` ↔ `*string` | Both ways |
| `bool` ↔ `*bool` | Both ways |
| `int`, `int8`...`int64` ↔ `*int`...`*int64` | Both ways |
| `uint`, `uint8`...`uint64` ↔ `*uint`...`*uint64` | Both ways |
| `float32`, `float64` ↔ `*float32`, `*float64` | Both ways |
| `decimal.Decimal` ↔ `*decimal.Decimal` | Both ways |
| `time.Time` ↔ `*time.Time` | Both ways |
| `timex.DateTime` ↔ `*timex.DateTime` | Both ways |
| `timex.Date` ↔ `*timex.Date` | Both ways |
| `timex.Time` ↔ `*timex.Time` | Both ways |

This means you can freely use pointer types in params structs for optional fields without worrying about type conversion.

Value-to-pointer converters allocate a new local value and return its address.
Pointer-to-value converters dereference non-nil pointers. If the source pointer
is nil, the converter returns the zero value for the destination type, for
example `""`, `false`, `0`, `decimal.Zero`, or a zero `time.Time` /
`timex.DateTime`.

## Framework Integration

The `copier` package is used internally by the CRUD `Create`, `CreateMany`,
`Update`, and `UpdateMany` builders to copy `TParams` fields into `TModel`
instances. The update builders use `WithIgnoreEmpty()` when merging incoming
models into existing models for partial updates.
