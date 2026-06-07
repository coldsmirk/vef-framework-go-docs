---
sidebar_position: 4
---

# Copier

`copier` 包提供结构体之间的字段复制，内置了 VEF 常用类型的类型转换器。

## 已审查公开 Surface

当前源码审计覆盖 `github.com/coldsmirk/vef-framework-go/copier` 的 9 个
top-level exported symbols，没有 exported fields，也没有 exported methods。
已审查 public-surface fingerprint 是
`44b6cf428fb9c642afca0cd25257c8ade57c9ac855b3ecc67cf575c1323fdf58`。

已审查 API：

| API | Contract |
| --- | --- |
| `copier.Copy(src, dst, options...)` | 通过 `github.com/jinzhu/copier` 的 `copier.CopyWithOption(dst, src, opt)` 从 `src` 复制到 `dst`；`dst` 必须是 pointer destination，否则底层 copier 会返回 error |
| `copier.CopyOption` | 修改 `github.com/jinzhu/copier` 底层 option struct 的 function option 类型 |
| `copier.TypeConverter` | `github.com/jinzhu/copier.TypeConverter` 的 alias |
| `copier.FieldNameMapping` | `github.com/jinzhu/copier.FieldNameMapping` 的 alias |
| `copier.WithIgnoreEmpty()` | 设置 `IgnoreEmpty = true` |
| `copier.WithDeepCopy()` | 设置 `DeepCopy = true` |
| `copier.WithCaseInsensitive()` | 设置 `CaseSensitive = false`；默认复制仍然是 case-sensitive |
| `copier.WithFieldNameMapping(...)` | 追加 mappings 到 `FieldNameMapping` |
| `copier.WithTypeConverters(...)` | 在内置 converters 后追加自定义 converters，而不是替换内置 converters |

## 快速开始

```go
import "github.com/coldsmirk/vef-framework-go/copier"

type UserParams struct {
    Username string
    Email    string
}

type User struct {
    Username string
    Email    string
    IsActive bool
}

params := UserParams{Username: "alice", Email: "alice@example.com"}
user := User{}

err := copier.Copy(params, &user)
// user.Username = "alice", user.Email = "alice@example.com"
```

## 选项

所有 copier 选项都使用公开的 `CopyOption` 类型。

`Copy(...)` 默认是 case-sensitive，并且总是先带上内置 converter 列表。
options 会按传入顺序应用。`WithTypeConverters(...)` 会追加自定义 converters，
而不是替换默认 converters。`WithFieldNameMapping(...)` 也会追加 mappings。

### WithIgnoreEmpty

跳过零值字段的复制。这是框架 CRUD `Update` 操作合并部分更新的方式：

```go
err := copier.Copy(params, &user, copier.WithIgnoreEmpty())
```

### WithDeepCopy

启用嵌套结构的深拷贝：

```go
err := copier.Copy(src, &dst, copier.WithDeepCopy())
```

### WithCaseInsensitive

启用不区分大小写的字段名匹配：

```go
err := copier.Copy(src, &dst, copier.WithCaseInsensitive())
```

### WithFieldNameMapping

添加自定义字段名映射：

```go
err := copier.Copy(src, &dst, copier.WithFieldNameMapping(
    copier.FieldNameMapping{
        SrcType: UserParams{},
        DstType: User{},
        Mapping: map[string]string{
            "Name": "Username",
        },
    },
))
```

### WithTypeConverters

添加自定义类型转换器：

```go
err := copier.Copy(src, &dst, copier.WithTypeConverters(
    copier.TypeConverter{
        SrcType: MyCustomType{},
        DstType: "",
        Fn: func(src interface{}) (interface{}, error) {
            return src.(MyCustomType).String(), nil
        },
    },
))
```

## 内置类型转换器

copier 包含所有常用类型的值 ↔ 指针自动转换器：

| 类型对 | 方向 |
| --- | --- |
| `string` ↔ `*string` | 双向 |
| `bool` ↔ `*bool` | 双向 |
| `int`, `int8`...`int64` ↔ `*int`...`*int64` | 双向 |
| `uint`, `uint8`...`uint64` ↔ `*uint`...`*uint64` | 双向 |
| `float32`, `float64` ↔ `*float32`, `*float64` | 双向 |
| `decimal.Decimal` ↔ `*decimal.Decimal` | 双向 |
| `time.Time` ↔ `*time.Time` | 双向 |
| `timex.DateTime` ↔ `*timex.DateTime` | 双向 |
| `timex.Date` ↔ `*timex.Date` | 双向 |
| `timex.Time` ↔ `*timex.Time` | 双向 |

这意味着你可以在 params 结构体中自由使用指针类型表示可选字段，无需担心类型转换问题。

value-to-pointer converter 会分配一个新的局部值并返回它的地址。
pointer-to-value converter 会解引用非 nil 指针。如果 source pointer 为 nil，
converter 会返回目标类型的零值，例如 `""`、`false`、`0`、`decimal.Zero`，
或 zero `time.Time` / `timex.DateTime`。

## 框架集成

`copier` 包被 CRUD 的 `Create`、`CreateMany`、`Update` 和 `UpdateMany`
builders 内部使用，用于将 `TParams` 字段复制到 `TModel` 实例。update builders
在把 incoming model 合并进 existing model 做部分更新时，会使用
`WithIgnoreEmpty()`。
