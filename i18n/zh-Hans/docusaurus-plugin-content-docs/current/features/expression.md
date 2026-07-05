---
sidebar_position: 15
---

# 表达式引擎

`expression` 包提供与后端无关的表达式契约。当前框架运行时会接入一个
pure-Go `expr-lang` engine，并把它提供给 API handler 注入和 mold 字段变换。

旧的 [`js` 包](./js-engine) 仍然是独立的 Goja JavaScript 运行时。业务规
则、派生字段等需要依赖稳定 `expression.Engine` 接口的场景，应该使用
`expression`。

## 运行时后端

VEF 当前在 core boot graph 中提供这个 engine，后端来自
`github.com/expr-lang/expr`。

当前后端是 pure Go。expression 模块本身不要求 CGO；应用是否启用 CGO
取决于你选择的数据库驱动或其他 native integration。

## 核心 API

```go
type Engine interface {
    Evaluate(ctx context.Context, source string, env any) (Value, error)
    Compile(source string, opts ...CompileOption) (Program, error)
}

type Program interface {
    Run(ctx context.Context, env any) (Value, error)
    Source() string
}
```

- `Evaluate` 会在一步内编译并执行 source expression。
- `Compile` 返回可复用的 `Program`；当前后端会在编译阶段解析并校验表达式。
- `env` 是变量环境，可以是 map，也可以是带 JSON tag 的 struct。
- Context cancellation 是 best-effort。当前 `expr-lang` backend 会在同步
  evaluation 开始前尊重已经取消的 context，但无法中断已经开始的 evaluation。

## 求值与解码

```go
import (
    "context"

    "github.com/coldsmirk/vef-framework-go/expression"
)

total, err := expression.EvaluateAs[float64](
    context.Background(),
    engine,
    "price * qty",
    struct {
        Price float64 `json:"price"`
        Qty   float64 `json:"qty"`
    }{Price: 2, Qty: 3},
)
```

`EvaluateAs[T]` 会执行表达式并把结果解码成 `T`。如果要手动处理结果，可
以使用 `Value` helpers：

| API | 契约 |
| --- | --- |
| `NewValue(raw)` | 包装 backend result；通常由 backend adapter 构造 |
| `Value.Interface()` | 返回原始 backend result |
| `Value.IsNil()` | 判断原始 backend result 是否为 `nil` |
| `Value.Bool()` | 返回 boolean result；非 bool 值会返回 `expression.ErrUnexpectedType` |
| `Value.Decode(target)` | 通过 JSON 把原始值解码到 non-nil pointer target |
| `DecodeValue[T](value)` | 把结果解码到一个新的 `T` 的 generic helper |

解码过程会经过 JSON，因此超出 float64 精度的大整数会有 JSON number 转换
本身的精度风险。nil result 会按 JSON `null` 解码；对非 pointer scalar
target 来说，这会得到该类型的零值。

## Predicate 匹配

布尔 predicate 使用 `Match`：

```go
ok, err := expression.Match(
    context.Background(),
    engine,
    ">= 5",
    map[string]any{"$": 10},
)
```

`Match` 会用 `expression.AsPredicate()` 编译表达式，并期望最终结果是
boolean。

## Handler 注入

core boot graph 已经注册了 `expression.Engine` 的 handler parameter
resolver，所以 API handler 可以直接声明这个参数：

```go
func CalculateTotal(
    ctx context.Context,
    engine expression.Engine,
    input CalculateTotalInput,
) (CalculateTotalOutput, error) {
    total, err := expression.EvaluateAs[float64](
        ctx,
        engine,
        "price * qty",
        input,
    )
    if err != nil {
        return CalculateTotalOutput{}, err
    }

    return CalculateTotalOutput{Total: total}, nil
}
```

应用代码应该依赖 `expression.Engine`，不要依赖内部 expr-lang adapter。

## Mold `expr` 字段

框架还注册了名为 `expr` 的 mold field transformer。它会把当前结构体作为
表达式环境，执行表达式，再把解码后的结果写入带标签的字段：

```go
type LineItem struct {
    Price float64 `json:"price"`
    Qty   float64 `json:"qty"`
    Total float64 `json:"total" mold:"expr=price * qty"`
}
```

几个重要细节：

- 当前结构体就是表达式环境。
- 字段按声明顺序变换，所以 `expr` 字段可以引用声明在它上方的 sibling
  field，包括更早的派生字段。
- 如果引用声明在当前字段下方的字段，会读到那个字段的零值。
- Mold 会按逗号拆分 tag function；如果表达式本身包含逗号，需要在 tag
  parameter 中写成 `0x2C`。

## 错误行为

后端 evaluation 错误会包装在 `expression.ErrEvaluationFailed` 下。类型转
换错误会从 `Value.Bool()`、`Value.Decode(...)` 或 `DecodeValue[T](...)`
返回。

错误细节：

| 错误 surface | 契约 |
| --- | --- |
| `expression.ErrEvaluationFailed` | backend compile/evaluation 失败时对 API 暴露的 result error |
| `ErrCodeEvaluationFailed` | `2500` |
| i18n key | `expression_evaluation_failed` |
| `expression.ErrUnexpectedType` | `Value.Bool()` 遇到非 bool 值时返回；message text 是 `expression: unexpected result type` |
| 空 mold `expr` tag | 返回 `expression: empty expression in field tag` |
| 不可设置的 mold target field | 返回 `expression: target field is not settable` |

在当前 `expr-lang` backend 下，格式错误的表达式会在 `Compile(...)`
阶段失败；运行期类型或值相关失败仍可能从 `Program.Run(...)` 返回。

其他公开 API 包括 `CompileOption`, `CompileOptions`, `AsPredicate()`,
`CompileOptions.Predicate`, `ErrCodeEvaluationFailed` 和 `ErrUnexpectedType`。
