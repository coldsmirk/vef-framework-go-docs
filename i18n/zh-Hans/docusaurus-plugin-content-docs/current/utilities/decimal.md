---
sidebar_position: 3
---

# Decimal

`decimal` 包重新导出 `shopspring/decimal v1.4.0`，用于任意精度
十进制算术，并额外提供 `NewFromAny` / `MustFromAny` 转换 helper。

`github.com/coldsmirk/vef-framework-go/decimal` 已审查公开 surface：

- 22 个 top-level symbols
- 0 个 exported fields
- 70 个 exported methods
- fingerprint `ea79b685aa80a0df3929fb69b8a3e0941805ce057e5c7ebf81a853da467d8401`

## Alias 契约

```go
import "github.com/coldsmirk/vef-framework-go/decimal"

var price decimal.Decimal
```

`decimal.Decimal` 是 `shopspring/decimal.Decimal` 的 type alias，不是 wrapper。
所有 `Decimal.*` methods 都是 upstream methods，签名和行为一致。该类型没有
exported fields。

VEF 重新导出了 constructors、constants 和 aggregators，但没有重新导出
upstream package-level knobs，例如 `DivisionPrecision`、
`MarshalJSONWithoutQuotes`、`PowPrecisionNegativeExponent` 或
`ExpMaxIterations`。如果应用必须修改这些全局设置，需要直接 import
`github.com/shopspring/decimal`。

## Top-Level API

常量：

```go
decimal.Zero   // 0
decimal.One    // 1
```

不要用 `==` 比较 `decimal.Zero`；使用 `Decimal.Equal` 或 `Decimal.Cmp`。

重新导出的 constructors：

```go
decimal.New(value, exp)                    // raw value 和 exponent
decimal.NewFromInt(42)                     // int64
decimal.NewFromInt32(42)                   // int32
decimal.NewFromUint64(42)                  // uint64
decimal.NewFromFloat(3.14)                 // float64
decimal.NewFromFloat32(3.14)               // float32
decimal.NewFromFloatWithExponent(3.14, -2) // 指定指数
decimal.NewFromBigInt(bigInt, exp)         // *big.Int
decimal.NewFromBigRat(bigRat, precision)   // *big.Rat
decimal.NewFromString("123.45")            // string，返回 error
decimal.RequireFromString("123.45")        // string，失败时 panic
```

`decimal.NewFromFormattedString` 的第二个参数是 `*regexp.Regexp`，用于在解析前
移除匹配到的格式字符：

```go
cleanup := regexp.MustCompile("[$,]")
amount, err := decimal.NewFromFormattedString("$1,234.56", cleanup)
```

重新导出的 aggregators：

```go
decimal.Max(a, b, c...)
decimal.Min(a, b, c...)
decimal.Sum(a, b, c...)
decimal.Avg(a, b, c...)
decimal.RescalePair(a, b)
```

VEF-specific helpers：

```go
d, err := decimal.NewFromAny(value)
d = decimal.MustFromAny(value)
```

完整 top-level checklist：

| Group | APIs |
| --- | --- |
| Alias and constants | `decimal.Decimal`、`decimal.Zero`、`decimal.One` |
| Constructors | `decimal.New`、`decimal.NewFromInt`、`decimal.NewFromInt32`、`decimal.NewFromUint64`、`decimal.NewFromFloat`、`decimal.NewFromFloat32`、`decimal.NewFromFloatWithExponent`、`decimal.NewFromBigInt`、`decimal.NewFromBigRat`、`decimal.NewFromString`、`decimal.NewFromFormattedString`、`decimal.RequireFromString` |
| Conversion helpers | `decimal.NewFromAny`、`decimal.MustFromAny` |
| Aggregators | `decimal.Max`、`decimal.Min`、`decimal.Sum`、`decimal.Avg`、`decimal.RescalePair` |

## `NewFromAny`

`decimal.NewFromAny` 会把常见 Go 值转换成 `decimal.Decimal`。

支持的输入 family：

| 类型 | 行为 |
| --- | --- |
| `Decimal`、`*Decimal` | 直接透传；nil `*Decimal` 返回 `decimal.Zero` |
| `int`、`int8`–`int64` | 整数转换 |
| `uint`、`uint8`–`uint64` | 无符号整数转换 |
| `float32`、`float64` | 浮点数转换 |
| `string`、`[]byte` | 通过 `NewFromString` 解析 |
| `bool` | `true` 转成 `decimal.One`；`false` 转成 `decimal.Zero` |
| `fmt.Stringer` | 解析其 `String()` 返回值 |

不支持的输入会返回 message 以 `decimal: unsupported type` 开头的 error；这个
sentinel 没有导出。`decimal.MustFromAny` 会在任何转换错误上 panic。

## Method Families

因为 `decimal.Decimal` 是 type alias，method set 继承自
`shopspring/decimal.Decimal`。public API index 列出精确签名；下表是已审查
的 70 个 methods 完整清单。

| Family | Methods |
| --- | --- |
| Arithmetic | `Decimal.Abs`、`Decimal.Neg`、`Decimal.Add`、`Decimal.Sub`、`Decimal.Mul`、`Decimal.Div`、`Decimal.DivRound`、`Decimal.Mod`、`Decimal.QuoRem` |
| Powers and transcendental functions | `Decimal.Pow`、`Decimal.PowBigInt`、`Decimal.PowInt32`、`Decimal.PowWithPrecision`、`Decimal.Sin`、`Decimal.Cos`、`Decimal.Tan`、`Decimal.Atan`、`Decimal.Ln`、`Decimal.ExpTaylor`、`Decimal.ExpHullAbrham` |
| Comparison and sign | `Decimal.Cmp`、`Decimal.Compare`、`Decimal.Equal`、`Decimal.Equals`、`Decimal.GreaterThan`、`Decimal.GreaterThanOrEqual`、`Decimal.LessThan`、`Decimal.LessThanOrEqual`、`Decimal.Sign` |
| Rounding and scale | `Decimal.Ceil`、`Decimal.Floor`、`Decimal.Round`、`Decimal.RoundBank`、`Decimal.RoundCash`、`Decimal.RoundCeil`、`Decimal.RoundDown`、`Decimal.RoundFloor`、`Decimal.RoundUp`、`Decimal.Shift`、`Decimal.Truncate` |
| Inspection and conversion | `Decimal.BigFloat`、`Decimal.BigInt`、`Decimal.Rat`、`Decimal.Float64`、`Decimal.InexactFloat64`、`Decimal.IntPart`、`Decimal.Coefficient`、`Decimal.CoefficientInt64`、`Decimal.Exponent`、`Decimal.NumDigits`、`Decimal.IsInteger`、`Decimal.IsNegative`、`Decimal.IsPositive`、`Decimal.IsZero`、`Decimal.Copy` |
| String formatting | `Decimal.String`、`Decimal.StringFixed`、`Decimal.StringFixedBank`、`Decimal.StringFixedCash`、`Decimal.StringScaled` |
| Encoding, database, and scanner interfaces | `Decimal.MarshalBinary`、`Decimal.UnmarshalBinary`、`Decimal.MarshalJSON`、`Decimal.UnmarshalJSON`、`Decimal.MarshalText`、`Decimal.UnmarshalText`、`Decimal.GobEncode`、`Decimal.GobDecode`、`Decimal.Scan`、`Decimal.Value` |

## 运行时注意点

- `Decimal.Div` 在 quotient 不能整除时使用 upstream `DivisionPrecision`。
  upstream 默认保留小数点后 16 位。
- `Decimal.Div`、`Decimal.DivRound`、`Decimal.Mod`、`Decimal.QuoRem` 遇到
  division by zero 会 panic。
- `Decimal.MarshalJSON` 默认输出 quoted JSON string，例如 `"123.45"`。
  upstream `MarshalJSONWithoutQuotes` global 可以切换为 JSON number，但在
  JavaScript 客户端中可能丢失精度。
- `Decimal.UnmarshalJSON` 同时接受 quoted 和 numeric JSON input。
- `decimal.NewFromFloat` 和 `decimal.NewFromFloat32` 会从 binary floating-point
  value 转换；金额类输入优先使用 `decimal.NewFromString`、整数 constructor
  或 scaled integer storage。
- `Decimal.Float64` 返回 `(value, exact)`；`Decimal.InexactFloat64` 只返回最近
  的 `float64`。
- `decimal.NewFromFloat`、`decimal.NewFromFloat32` 和
  `decimal.NewFromFloatWithExponent` 遇到 NaN 或 infinity 会 panic。

## 在模型中使用

```go
type Product struct {
    orm.FullAuditedModel
    Name  string          `json:"name" bun:"name"`
    Price decimal.Decimal `json:"price" bun:"price,type:decimal(10,2)"`
}
```

`decimal.Decimal` 类型完全支持：
- Bun ORM（数据库扫描/取值）
- JSON 序列化/反序列化
- `copier` 包（值 ↔ 指针转换）
- `mapx` 包（解码钩子）
