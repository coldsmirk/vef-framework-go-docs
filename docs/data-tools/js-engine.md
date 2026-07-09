---
sidebar_position: 2
---

# JS Engine

The `js` package provides an embedded JavaScript runtime powered by [goja](https://github.com/dop251/goja), enabling server-side JavaScript execution within Go applications.

## Quick Start

```go
import "github.com/coldsmirk/vef-framework-go/js"

vm, err := js.New()
if err != nil {
    return err
}

result, err := vm.RunString(`1 + 2`)
fmt.Println(result.Export()) // 3
```

## Runtime Features

When you call `js.New()`, the runtime is pre-configured with:

- a fresh `goja.New()` runtime
- `vm.SetParserOptions(parser.WithDisableSourceMaps)`
- `vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))`
- pre-compiled browser bundles executed in this order: `dayjs`, `Big`,
  `utils`, then `validator`

If any bundled library fails during `RunProgram`, `New()` returns `nil` and the
first load error. The bundled libraries are compiled in strict mode at package
initialization.

### Preloaded JavaScript Globals

`js.New()` exposes these globals by executing vendored browser/UMD bundles; it
does not install a Node-style module loader.

| Global | Bundle | Version | Contract |
| --- | --- | --- | --- |
| `dayjs` | `libs/day.v1_11_19.js` | Day.js 1.11.19 | Date/time parsing, formatting, arithmetic, and comparison helpers exposed by the bundled Day.js build |
| `Big` | `libs/big.v7_0_1.js` | big.js 7.0.1 | Arbitrary-precision decimal constructor and methods exposed by the bundled big.js build |
| `utils` | `libs/utils.v12_7_0.js` | utils 12.7.0 | Utility helpers exposed by the bundled utility bundle, including examples covered by tests such as `capitalize`, `camel`, `snake`, `unique`, `sum`, `group`, `sort`, and `max` |
| `validator` | `libs/validator.v13_15_20.js` | validator.js 13.15.20 | String validators and sanitizers exposed by the bundled validator.js build, including examples covered by tests such as `isEmail`, `isURL`, `isUUID`, `isJSON`, `isNumeric`, and `isISO8601` |

The member APIs of these JavaScript globals follow the vendored JavaScript
bundles. VEF does not wrap individual library functions.

### Runtime Boundary

The VEF runtime setup does not call `require.NewRegistry`, does not register
native modules, and does not enable a `console` shim. VEF also does not install
Node APIs such as `fs`, `net`, or timers. Use `vm.Set(...)` or the pass-through
goja runtime methods when an application needs additional globals or host
functions.

`js.New()` is not a sandbox policy by itself. Time limits, cancellation, and
interrupt behavior are controlled through the pass-through goja runtime surface,
for example `Runtime.Interrupt(...)` and `Runtime.ClearInterrupt()`.

## Go–JavaScript Interop

### Passing Go Values

```go
vm.Set("user", map[string]any{
    "name": "Alice",
    "age":  30,
})

result, _ := vm.RunString(`user.name + " is " + user.age`)
// → "Alice is 30"
```

### Passing Go Functions

```go
vm.Set("greet", func(name string) string {
    return "Hello, " + name + "!"
})

result, _ := vm.RunString(`greet("World")`)
// → "Hello, World!"
```

### Returning Values

```go
result, _ := vm.RunString(`({name: "Alice", score: 95})`)
obj := result.Export().(map[string]any)
// obj["name"] = "Alice", obj["score"] = 95
```

## Compiling Scripts

For repeated execution, pre-compile scripts:

```go
program, err := js.Compile("my-script", `
    function calculate(a, b) {
        return a + b;
    }
    calculate(x, y);
`, true) // third argument enables strict mode

vm.Set("x", 10)
vm.Set("y", 20)
result, err := vm.RunProgram(program)
// → 30
```

## Parsing AST

```go
ast, err := js.Parse("my-script", scriptSource)
```

`Parse` returns `*js.AstProgram` and always calls `goja.Parse` with
`parser.WithDisableSourceMaps`.

## Type Aliases

The package re-exports key goja types for convenience:

```go
js.Runtime    // = goja.Runtime
js.Value      // = goja.Value
js.Object     // = goja.Object
js.Program    // = goja.Program
js.AstProgram // = ast.Program
```

Function aliases are also exported for common goja helpers:

```go
js.Compile
js.MustCompile
js.IsNaN
js.IsString
js.IsBigInt
js.IsNumber
js.IsInfinity
js.IsUndefined
js.IsNull
```

### Pass-through Policy

The goja pass-through surface is intentionally exposed for convenience:
`js.Runtime`, `js.Value`, `js.Object`, `js.Program`, and `js.AstProgram`
follow the upstream [github.com/dop251/goja](https://github.com/dop251/goja)
API at the pinned source dependency version
`v0.0.0-20260311135729-065cd970411c`. `js.Compile`, `js.MustCompile`, and the
`js.Is*` helpers are direct goja function aliases; they do not inherit
`js.New()` runtime parser options. VEF adds `js.New()`, `js.Parse(...)`, and the
preload/configuration behavior described above. The public API index lists the
exact signatures for every exported alias.

## Thread Safety

> **WARNING**: The JavaScript runtime is **NOT** thread-safe. Each goroutine must create its own runtime instance via `js.New()`.
