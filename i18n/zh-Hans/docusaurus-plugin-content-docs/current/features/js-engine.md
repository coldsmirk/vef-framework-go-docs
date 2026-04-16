---
sidebar_position: 15
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

- **禁用 Source Maps** 以提升性能
- **JSON 结构体标签映射** — Go 结构体自动使用 `json` 标签映射
- **预加载库**：
  - `dayjs` — 日期时间操作
  - `Big` — 任意精度算术
  - `utils` — 工具函数
  - `validator` — 数据验证

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
`)

vm.Set("x", 10)
vm.Set("y", 20)
result, err := vm.RunProgram(program)
// → 30
```

## 解析 AST

```go
ast, err := js.Parse("my-script", scriptSource)
```

## 类型别名

该包重导出了 goja 的关键类型：

```go
js.Runtime    // = goja.Runtime
js.Value      // = goja.Value
js.Object     // = goja.Object
js.Program    // = goja.Program
```

## 线程安全

> **警告**：JavaScript 运行时**不是**线程安全的。每个 goroutine 必须通过 `js.New()` 创建自己的运行时实例。
