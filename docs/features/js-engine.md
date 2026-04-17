---
sidebar_position: 15
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

- **Source maps disabled** for better performance
- **JSON struct tag mapping** — Go structs are automatically mapped using `json` tags
- **Pre-loaded libraries**:
  - `dayjs` — date/time manipulation
  - `Big` — arbitrary-precision arithmetic
  - `utils` — utility functions
  - `validator` — data validation

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
`)

vm.Set("x", 10)
vm.Set("y", 20)
result, err := vm.RunProgram(program)
// → 30
```

## Parsing AST

```go
ast, err := js.Parse("my-script", scriptSource)
```

## Type Aliases

The package re-exports key goja types for convenience:

```go
js.Runtime    // = goja.Runtime
js.Value      // = goja.Value
js.Object     // = goja.Object
js.Program    // = goja.Program
```

## Thread Safety

> **WARNING**: The JavaScript runtime is **NOT** thread-safe. Each goroutine must create its own runtime instance via `js.New()`.
