---
sidebar_position: 3
---

# Decimal

`decimal` 包提供重新导出的 `shopspring/decimal` 接口，附加了便捷的构造函数和常量，用于任意精度十进制算术。

## 类型

```go
import "github.com/coldsmirk/vef-framework-go/decimal"

// decimal.Decimal 是 shopspring/decimal.Decimal 的别名
var price decimal.Decimal
```

## 常量

预定义的常用十进制值：

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

## 构造函数

重新导出了所有标准 `shopspring/decimal` 构造函数：

```go
decimal.New(value, exp)                   // 原始构造
decimal.NewFromInt(42)                    // 从 int64
decimal.NewFromInt32(42)                  // 从 int32
decimal.NewFromUint64(42)                // 从 uint64
decimal.NewFromFloat(3.14)               // 从 float64
decimal.NewFromFloat32(3.14)             // 从 float32
decimal.NewFromFloatWithExponent(3.14, -2) // 指定指数
decimal.NewFromBigInt(bigInt, exp)        // 从 *big.Int
decimal.NewFromBigRat(bigRat, precision)  // 从 *big.Rat
decimal.NewFromString("123.45")          // 从字符串（返回 error）
decimal.NewFromFormattedString("1,234.56", ",") // 从格式化字符串
decimal.RequireFromString("123.45")      // 从字符串（失败时 panic）
```

## `NewFromAny` — 通用转换器

将任意 Go 值转换为 Decimal：

```go
d, err := decimal.NewFromAny(value)
```

支持的输入类型：

| 类型 | 行为 |
| --- | --- |
| `Decimal`、`*Decimal` | 直接透传 |
| `int`、`int8`–`int64` | 整数转换 |
| `uint`、`uint8`–`uint64` | 无符号整数转换 |
| `float32`、`float64` | 浮点数转换 |
| `string`、`[]byte` | 字符串解析 |
| `bool` | `true` → 1，`false` → 0 |
| `fmt.Stringer` | 使用 `.String()` 方法 |

确定转换不会失败时的 panic 变体：

```go
d := decimal.MustFromAny(value) // 转换失败时 panic
```

## 工具函数

```go
decimal.Max(a, b, c...)       // 多个 decimal 的最大值
decimal.Min(a, b, c...)       // 多个 decimal 的最小值
decimal.Sum(a, b, c...)       // 多个 decimal 的总和
decimal.Avg(a, b, c...)       // 多个 decimal 的平均值
decimal.RescalePair(a, b)     // 将两个 decimal 缩放到相同指数
```

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
