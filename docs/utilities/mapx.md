---
sidebar_position: 6
---

# Mapx

The `mapx` package provides bidirectional conversion between Go structs and `map[string]any`, built on top of `mapstructure`.

## Struct to Map

```go
import "github.com/coldsmirk/vef-framework-go/mapx"

type User struct {
    Name  string `mapstructure:"name"`
    Email string `mapstructure:"email"`
    Age   int    `mapstructure:"age"`
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
// Use JSON tags instead of mapstructure tags
m, err := mapx.ToMap(user, mapx.WithTagName("json"))

// Weak type conversion (string "123" → int 123)
user, err := mapx.FromMap[User](data, mapx.WithWeaklyTypedInput())

// Include nil values in decoding
user, err := mapx.FromMap[User](data, mapx.WithDecodeNil())
```

### Available Options

| Option | Effect |
| --- | --- |
| `WithTagName(tag)` | Use a specific struct tag (default: `mapstructure`) |
| `WithWeaklyTypedInput()` | Enable weak type conversion |
| `WithDecodeNil()` | Include nil values in decoding |

## Custom Decoder

For advanced use cases, create a reusable decoder:

```go
var result User
decoder, err := mapx.NewDecoder(&result, mapx.WithTagName("json"))
if err != nil {
    return err
}
err = decoder.Decode(data)
```

## Decode Hooks

The package includes built-in decode hooks for VEF types:

- `timex.DateTime`, `timex.Date`, `timex.Time` — automatic string parsing
- `decimal.Decimal` — automatic decimal conversion

These hooks are registered automatically, so you can decode maps containing string timestamps into structs with `timex` fields without any configuration.
