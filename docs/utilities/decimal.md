---
sidebar_position: 3
---

# Decimal

The `decimal` package provides a re-exported `shopspring/decimal` interface with additional convenience constructors and constants for arbitrary-precision decimal arithmetic.

## Type

```go
import "github.com/coldsmirk/vef-framework-go/decimal"

// decimal.Decimal is an alias for shopspring/decimal.Decimal
var price decimal.Decimal
```

## Constants

Pre-defined decimal values for common arithmetic:

```go
decimal.Zero   // 0
decimal.One    // 1
decimal.Two    // 2
decimal.Three  // 3
decimal.Four   // 4
decimal.Five   // 5
decimal.Six    // 6
decimal.Seven  // 7
decimal.Eight  // 8
decimal.Nine   // 9
decimal.Ten    // 10
```

## Constructors

All standard `shopspring/decimal` constructors are re-exported:

```go
decimal.New(value, exp)                   // Raw constructor
decimal.NewFromInt(42)                    // From int64
decimal.NewFromInt32(42)                  // From int32
decimal.NewFromUint64(42)                // From uint64
decimal.NewFromFloat(3.14)               // From float64
decimal.NewFromFloat32(3.14)             // From float32
decimal.NewFromFloatWithExponent(3.14, -2) // With specific exponent
decimal.NewFromBigInt(bigInt, exp)        // From *big.Int
decimal.NewFromBigRat(bigRat, precision)  // From *big.Rat
decimal.NewFromString("123.45")          // From string (returns error)
decimal.NewFromFormattedString("1,234.56", ",") // From formatted string
decimal.RequireFromString("123.45")      // From string (panics on error)
```

## `NewFromAny` — Universal Converter

Convert any Go value to a Decimal:

```go
d, err := decimal.NewFromAny(value)
```

Supported input types:

| Type | Behavior |
| --- | --- |
| `Decimal`, `*Decimal` | Direct pass-through |
| `int`, `int8`–`int64` | Integer conversion |
| `uint`, `uint8`–`uint64` | Unsigned integer conversion |
| `float32`, `float64` | Float conversion |
| `string`, `[]byte` | String parsing |
| `bool` | `true` → 1, `false` → 0 |
| `fmt.Stringer` | Uses `.String()` method |

Panic variant for when you know the conversion will succeed:

```go
d := decimal.MustFromAny(value) // Panics on error
```

## Utility Functions

```go
decimal.Max(a, b, c...)       // Maximum of multiple decimals
decimal.Min(a, b, c...)       // Minimum of multiple decimals
decimal.Sum(a, b, c...)       // Sum of multiple decimals
decimal.Avg(a, b, c...)       // Average of multiple decimals
decimal.RescalePair(a, b)     // Rescale two decimals to same exponent
```

## Usage in Models

```go
type Product struct {
    orm.FullAuditedModel
    Name  string          `json:"name" bun:"name"`
    Price decimal.Decimal `json:"price" bun:"price,type:decimal(10,2)"`
}
```

The `decimal.Decimal` type is fully supported by:
- Bun ORM (database scanning/value)
- JSON marshaling/unmarshaling
- The `copier` package (value ↔ pointer conversion)
- The `mapx` package (decode hooks)
