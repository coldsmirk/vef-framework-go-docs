---
sidebar_position: 3
---

# Decimal

The `decimal` package re-exports `shopspring/decimal v1.4.0` for
arbitrary-precision decimal arithmetic and adds the `NewFromAny` /
`MustFromAny` conversion helpers.

Reviewed public surface for `github.com/coldsmirk/vef-framework-go/decimal`:

- 22 top-level symbols
- 0 exported fields
- 70 exported methods
- fingerprint `ea79b685aa80a0df3929fb69b8a3e0941805ce057e5c7ebf81a853da467d8401`

## Alias Contract

```go
import "github.com/coldsmirk/vef-framework-go/decimal"

var price decimal.Decimal
```

`decimal.Decimal` is a type alias for `shopspring/decimal.Decimal`, not a
wrapper. All `Decimal.*` methods are the upstream methods with the same
signatures and behavior. The type has no exported fields.

VEF re-exports constructors, constants, and aggregators, but it does not
re-export upstream package-level knobs such as `DivisionPrecision`,
`MarshalJSONWithoutQuotes`, `PowPrecisionNegativeExponent`, or
`ExpMaxIterations`. Import `github.com/shopspring/decimal` directly if an app
must change those globals.

## Top-Level APIs

Constants:

```go
decimal.Zero   // 0
decimal.One    // 1
```

Do not compare `decimal.Zero` with `==`; use `Decimal.Equal` or `Decimal.Cmp`.

Re-exported constructors:

```go
decimal.New(value, exp)                    // raw value and exponent
decimal.NewFromInt(42)                     // int64
decimal.NewFromInt32(42)                   // int32
decimal.NewFromUint64(42)                  // uint64
decimal.NewFromFloat(3.14)                 // float64
decimal.NewFromFloat32(3.14)               // float32
decimal.NewFromFloatWithExponent(3.14, -2) // With specific exponent
decimal.NewFromBigInt(bigInt, exp)         // *big.Int
decimal.NewFromBigRat(bigRat, precision)   // *big.Rat
decimal.NewFromString("123.45")            // string, returns error
decimal.RequireFromString("123.45")        // string, panics on error
```

`decimal.NewFromFormattedString` takes a `*regexp.Regexp` that removes matched
formatting characters before parsing:

```go
cleanup := regexp.MustCompile("[$,]")
amount, err := decimal.NewFromFormattedString("$1,234.56", cleanup)
```

Re-exported aggregators:

```go
decimal.Max(a, b, c...)
decimal.Min(a, b, c...)
decimal.Sum(a, b, c...)
decimal.Avg(a, b, c...)
decimal.RescalePair(a, b)
```

VEF-specific helpers:

```go
d, err := decimal.NewFromAny(value)
d = decimal.MustFromAny(value)
```

Complete top-level checklist:

| Group | APIs |
| --- | --- |
| Alias and constants | `decimal.Decimal`, `decimal.Zero`, `decimal.One` |
| Constructors | `decimal.New`, `decimal.NewFromInt`, `decimal.NewFromInt32`, `decimal.NewFromUint64`, `decimal.NewFromFloat`, `decimal.NewFromFloat32`, `decimal.NewFromFloatWithExponent`, `decimal.NewFromBigInt`, `decimal.NewFromBigRat`, `decimal.NewFromString`, `decimal.NewFromFormattedString`, `decimal.RequireFromString` |
| Conversion helpers | `decimal.NewFromAny`, `decimal.MustFromAny` |
| Aggregators | `decimal.Max`, `decimal.Min`, `decimal.Sum`, `decimal.Avg`, `decimal.RescalePair` |

## `NewFromAny`

`decimal.NewFromAny` converts common Go values into `decimal.Decimal`.

Supported input families:

| Type | Behavior |
| --- | --- |
| `Decimal`, `*Decimal` | Direct pass-through; nil `*Decimal` returns `decimal.Zero` |
| `int`, `int8`–`int64` | Integer conversion |
| `uint`, `uint8`–`uint64` | Unsigned integer conversion |
| `float32`, `float64` | Float conversion |
| `string`, `[]byte` | Parsed through `NewFromString` |
| `bool` | `true` becomes `decimal.One`; `false` becomes `decimal.Zero` |
| `fmt.Stringer` | Parses the returned `String()` value |

Unsupported inputs return an error whose message starts with
`decimal: unsupported type`; the sentinel is not exported. `decimal.MustFromAny`
panics on any conversion error.

## Method Families

Since `decimal.Decimal` is a type alias, the method set is inherited from
`shopspring/decimal.Decimal`. The public API index lists exact signatures; this
table is the reviewed completeness checklist for all 70 methods.

| Family | Methods |
| --- | --- |
| Arithmetic | `Decimal.Abs`, `Decimal.Neg`, `Decimal.Add`, `Decimal.Sub`, `Decimal.Mul`, `Decimal.Div`, `Decimal.DivRound`, `Decimal.Mod`, `Decimal.QuoRem` |
| Powers and transcendental functions | `Decimal.Pow`, `Decimal.PowBigInt`, `Decimal.PowInt32`, `Decimal.PowWithPrecision`, `Decimal.Sin`, `Decimal.Cos`, `Decimal.Tan`, `Decimal.Atan`, `Decimal.Ln`, `Decimal.ExpTaylor`, `Decimal.ExpHullAbrham` |
| Comparison and sign | `Decimal.Cmp`, `Decimal.Compare`, `Decimal.Equal`, `Decimal.Equals`, `Decimal.GreaterThan`, `Decimal.GreaterThanOrEqual`, `Decimal.LessThan`, `Decimal.LessThanOrEqual`, `Decimal.Sign` |
| Rounding and scale | `Decimal.Ceil`, `Decimal.Floor`, `Decimal.Round`, `Decimal.RoundBank`, `Decimal.RoundCash`, `Decimal.RoundCeil`, `Decimal.RoundDown`, `Decimal.RoundFloor`, `Decimal.RoundUp`, `Decimal.Shift`, `Decimal.Truncate` |
| Inspection and conversion | `Decimal.BigFloat`, `Decimal.BigInt`, `Decimal.Rat`, `Decimal.Float64`, `Decimal.InexactFloat64`, `Decimal.IntPart`, `Decimal.Coefficient`, `Decimal.CoefficientInt64`, `Decimal.Exponent`, `Decimal.NumDigits`, `Decimal.IsInteger`, `Decimal.IsNegative`, `Decimal.IsPositive`, `Decimal.IsZero`, `Decimal.Copy` |
| String formatting | `Decimal.String`, `Decimal.StringFixed`, `Decimal.StringFixedBank`, `Decimal.StringFixedCash`, `Decimal.StringScaled` |
| Encoding, database, and scanner interfaces | `Decimal.MarshalBinary`, `Decimal.UnmarshalBinary`, `Decimal.MarshalJSON`, `Decimal.UnmarshalJSON`, `Decimal.MarshalText`, `Decimal.UnmarshalText`, `Decimal.GobEncode`, `Decimal.GobDecode`, `Decimal.Scan`, `Decimal.Value` |

## Runtime Notes

- `Decimal.Div` uses upstream `DivisionPrecision` when the quotient does not
  divide exactly. The upstream default is 16 digits after the decimal point.
- `Decimal.Div`, `Decimal.DivRound`, `Decimal.Mod`, and `Decimal.QuoRem` panic
  on division by zero.
- `Decimal.MarshalJSON` emits a quoted JSON string by default, for example
  `"123.45"`. The upstream `MarshalJSONWithoutQuotes` global switches it to a
  JSON number, but that can lose precision in JavaScript clients.
- `Decimal.UnmarshalJSON` accepts both quoted and numeric JSON input.
- `decimal.NewFromFloat` and `decimal.NewFromFloat32` convert from binary
  floating-point values; use `decimal.NewFromString`, integer constructors, or
  scaled integer storage for money-like inputs.
- `Decimal.Float64` returns `(value, exact)`; `Decimal.InexactFloat64` returns
  only the nearest `float64`.
- `decimal.NewFromFloat`, `decimal.NewFromFloat32`, and
  `decimal.NewFromFloatWithExponent` panic for NaN or infinity.

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
