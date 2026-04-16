---
sidebar_position: 4
---

# Copier

The `copier` package provides struct-to-struct field copying with built-in type converters for common VEF types.

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
        DstType: copier.String,
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

## Framework Integration

The `copier` package is used internally by the CRUD `Create` and `Update` builders to copy `TParams` fields into `TModel` instances. The `Update` builder specifically uses `WithIgnoreEmpty()` to support partial updates.
