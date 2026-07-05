---
sidebar_position: 15
---

# Expression Engine

The `expression` package exposes a backend-agnostic expression contract. The
current framework runtime wires a pure-Go `expr-lang` engine and makes it available to
API handlers and mold field transformations.

The older [`js` package](./js-engine) is still a separate Goja JavaScript
runtime. Use `expression` for business rules and derived fields that should stay
behind the stable `expression.Engine` interface.

## Runtime Backend

VEF currently provides the engine from the core boot graph, backed by
`github.com/expr-lang/expr`.

The current backend is pure Go. The expression module itself does not require
CGO; whether your application enables CGO depends on the database drivers or
other native integrations you choose.

## Core API

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

- `Evaluate` compiles and evaluates a source expression in one step.
- `Compile` returns a reusable `Program`; the current backend parses and
  validates the expression eagerly.
- `env` is the variable environment. It can be a map or a struct with JSON
  tags.
- Context cancellation is best-effort. The current `expr-lang` backend honors
  an already-canceled context before starting synchronous evaluation, but cannot
  interrupt evaluation already in flight.

## Evaluating Values

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

`EvaluateAs[T]` evaluates and decodes the result into `T`. For manual result
handling, use the `Value` helpers:

| API | Contract |
| --- | --- |
| `NewValue(raw)` | wraps a backend result; backend adapters normally construct it |
| `Value.Interface()` | returns the raw backend result |
| `Value.IsNil()` | reports whether the raw backend result is `nil` |
| `Value.Bool()` | returns a boolean result or `expression.ErrUnexpectedType` for non-bool values |
| `Value.Decode(target)` | JSON-decodes the raw value into a non-nil pointer target |
| `DecodeValue[T](value)` | generic helper that decodes into a new `T` |

Decoding goes through JSON, so very large integers can lose precision in the
same way as JSON number conversion. A nil result decodes as JSON `null`; for
non-pointer scalar targets this produces the type's zero value.

## Predicate Matching

Use `Match` for boolean predicates:

```go
ok, err := expression.Match(
    context.Background(),
    engine,
    ">= 5",
    map[string]any{"$": 10},
)
```

`Match` compiles the expression with `expression.AsPredicate()`, then expects a
boolean result.

## Handler Injection

The core boot graph registers a handler parameter resolver for
`expression.Engine`, so API handlers can request it directly:

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

Application code should depend on `expression.Engine`, not the internal expr-lang
adapter.

## Mold `expr` Fields

The framework also registers a mold field transformer named `expr`. It evaluates
the expression against the containing struct and writes the decoded result into
the tagged field:

```go
type LineItem struct {
    Price float64 `json:"price"`
    Qty   float64 `json:"qty"`
    Total float64 `json:"total" mold:"expr=price * qty"`
}
```

Important details:

- The containing struct is the expression environment.
- Fields are transformed in declaration order, so an `expr` field can reference
  sibling fields declared above it, including earlier derived fields.
- A reference to a field declared below the current field sees that field's zero
  value.
- Mold splits tag functions on commas. If an expression contains a comma, escape
  it as `0x2C` in the tag parameter.

## Error Behavior

Backend evaluation failures are wrapped under `expression.ErrEvaluationFailed`.
Type conversion failures surface from `Value.Bool()`, `Value.Decode(...)`, or
`DecodeValue[T](...)`.

Error details:

| Error surface | Contract |
| --- | --- |
| `expression.ErrEvaluationFailed` | API-facing result error for backend compile/evaluation failures |
| `ErrCodeEvaluationFailed` | `2500` |
| i18n key | `expression_evaluation_failed` |
| `expression.ErrUnexpectedType` | returned by `Value.Bool()` for non-bool values; message text is `expression: unexpected result type` |
| empty mold `expr` tag | returns `expression: empty expression in field tag` |
| non-settable mold target field | returns `expression: target field is not settable` |

With the current `expr-lang` backend, malformed expressions fail during
`Compile(...)`; runtime type/value failures can still surface from
`Program.Run(...)`.

Supporting public APIs include `CompileOption`, `CompileOptions`,
`CompileOptions.Predicate`, `AsPredicate()`, `ErrCodeEvaluationFailed`, and
`ErrUnexpectedType`.
