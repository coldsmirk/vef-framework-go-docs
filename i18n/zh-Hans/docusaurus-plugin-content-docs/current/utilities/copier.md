---
sidebar_position: 4
---

# Copier

`copier` 包提供结构体之间的字段复制，内置了 VEF 常用类型的类型转换器。

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
        DstType: copier.String,
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

## 框架集成

`copier` 包被 CRUD 的 `Create` 和 `Update` 构建器内部使用，用于将 `TParams` 字段复制到 `TModel` 实例。`Update` 构建器特别使用 `WithIgnoreEmpty()` 来支持部分更新。
