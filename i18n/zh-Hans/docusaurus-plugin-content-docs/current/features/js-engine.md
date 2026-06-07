---
sidebar_position: 16
---

# JS 引擎

`js` 包提供基于 [goja](https://github.com/dop251/goja) 的嵌入式 JavaScript 运行时，支持在 Go 应用中执行服务端 JavaScript。

## 快速开始

```go
import "github.com/coldsmirk/vef-framework-go/js"

vm, err := js.New()
if err != nil {
    return err
}

result, err := vm.RunString(`1 + 2`)
fmt.Println(result.Export()) // 3
```

## 运行时特性

调用 `js.New()` 时，运行时预配置了：

- 一个新的 `goja.New()` runtime
- `vm.SetParserOptions(parser.WithDisableSourceMaps)`
- `vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))`
- 按顺序执行预编译的 browser bundles：`dayjs`、`Big`、`utils`、`validator`

如果任意 bundled library 在 `RunProgram` 阶段失败，`New()` 会返回 `nil` 和第一
个加载错误。bundled libraries 在 package initialization 阶段以 strict mode
编译。

### 预加载 JavaScript Globals

`js.New()` 通过执行 vendored browser/UMD bundles 暴露这些 globals；它不会
安装 Node-style module loader。

| Global | Bundle | Version | 契约 |
| --- | --- | --- | --- |
| `dayjs` | `libs/day.v1_11_19.js` | Day.js 1.11.19 | bundled Day.js build 暴露的日期时间解析、格式化、运算与比较 helper |
| `Big` | `libs/big.v7_0_1.js` | big.js 7.0.1 | bundled big.js build 暴露的任意精度 decimal constructor 和 methods |
| `utils` | `libs/utils.v12_7_0.js` | utils 12.7.0 | bundled utility bundle 暴露的工具函数；测试覆盖了 `capitalize`、`camel`、`snake`、`unique`、`sum`、`group`、`sort`、`max` 等示例 |
| `validator` | `libs/validator.v13_15_20.js` | validator.js 13.15.20 | bundled validator.js build 暴露的 string validators 和 sanitizers；测试覆盖了 `isEmail`、`isURL`、`isUUID`、`isJSON`、`isNumeric`、`isISO8601` 等示例 |

这些 JavaScript globals 的 member APIs 遵循 vendored JavaScript bundles。VEF
不会逐个包装 library functions。

### Runtime Boundary

VEF runtime setup 不会调用 `require.NewRegistry`，不会注册 native modules，也
不会启用 `console` shim。VEF 也不会安装 `fs`、`net`、timers 等 Node APIs。
应用如果需要额外 globals 或 host functions，应使用 `vm.Set(...)` 或透传的
goja runtime methods。

`js.New()` 本身不是 sandbox policy。时间限制、取消与 interrupt 行为由透传的
goja runtime surface 控制，例如 `Runtime.Interrupt(...)` 和
`Runtime.ClearInterrupt()`。

## Go–JavaScript 互操作

### 传递 Go 值

```go
vm.Set("user", map[string]any{
    "name": "Alice",
    "age":  30,
})

result, _ := vm.RunString(`user.name + " is " + user.age`)
// → "Alice is 30"
```

### 传递 Go 函数

```go
vm.Set("greet", func(name string) string {
    return "Hello, " + name + "!"
})

result, _ := vm.RunString(`greet("World")`)
// → "Hello, World!"
```

### 返回值

```go
result, _ := vm.RunString(`({name: "Alice", score: 95})`)
obj := result.Export().(map[string]any)
// obj["name"] = "Alice", obj["score"] = 95
```

## 编译脚本

对于重复执行的脚本，进行预编译：

```go
program, err := js.Compile("my-script", `
    function calculate(a, b) {
        return a + b;
    }
    calculate(x, y);
`, true) // 第三个参数开启严格模式

vm.Set("x", 10)
vm.Set("y", 20)
result, err := vm.RunProgram(program)
// → 30
```

## 解析 AST

```go
ast, err := js.Parse("my-script", scriptSource)
```

`Parse` 返回 `*js.AstProgram`，并且始终用 `parser.WithDisableSourceMaps`
调用 `goja.Parse`。

## 类型别名

该包重导出了 goja 的关键类型：

```go
js.Runtime    // = goja.Runtime
js.Value      // = goja.Value
js.Object     // = goja.Object
js.Program    // = goja.Program
js.AstProgram // = ast.Program
```

常用 goja helper 也以函数别名形式导出：

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

### 透传策略

goja pass-through surface 是有意暴露的便利层：`js.Runtime`、`js.Value`、
`js.Object`、`js.Program` 和 `js.AstProgram` 遵循上游
[github.com/dop251/goja](https://github.com/dop251/goja) API，版本为 source
dependency 中固定的 `v0.0.0-20260311135729-065cd970411c`。`js.Compile`、
`js.MustCompile` 和 `js.Is*` helpers 是直接的 goja function aliases；它们不会
继承 `js.New()` 的 runtime parser options。VEF 额外提供 `js.New()`、
`js.Parse(...)` 以及上文描述的 preload/configuration behavior；所有透传方法
签名都在 public API index 中审计。

## 线程安全

> **警告**：JavaScript 运行时**不是**线程安全的。每个 goroutine 必须通过 `js.New()` 创建自己的运行时实例。
