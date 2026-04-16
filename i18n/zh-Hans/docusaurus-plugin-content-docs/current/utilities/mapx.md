---
sidebar_position: 6
---

# Mapx

`mapx` 包提供 Go 结构体与 `map[string]any` 之间的双向转换，基于 `mapstructure` 构建。

## 结构体转 Map

```go
import "github.com/coldsmirk/vef-framework-go/mapx"

type User struct {
    Name  string `mapstructure:"name"`
    Email string `mapstructure:"email"`
    Age   int    `mapstructure:"age"`
}

user := User{Name: "Alice", Email: "alice@example.com", Age: 30}
m, err := mapx.ToMap(user)
// m = map[string]any{"name": "Alice", "email": "alice@example.com", "age": 30}
```

## Map 转结构体

```go
data := map[string]any{
    "name":  "Bob",
    "email": "bob@example.com",
    "age":   25,
}

user, err := mapx.FromMap[User](data)
// user.Name = "Bob", user.Email = "bob@example.com", user.Age = 25
```

## 解码器选项

`ToMap` 和 `FromMap` 都接受可选的 `DecoderOption`：

```go
// 使用 JSON 标签替代 mapstructure 标签
m, err := mapx.ToMap(user, mapx.WithTagName("json"))

// 弱类型转换（字符串 "123" → int 123）
user, err := mapx.FromMap[User](data, mapx.WithWeaklyTypedInput())

// 在解码中包含 nil 值
user, err := mapx.FromMap[User](data, mapx.WithDecodeNil())
```

### 可用选项

| 选项 | 效果 |
| --- | --- |
| `WithTagName(tag)` | 使用指定的结构体标签（默认：`mapstructure`）|
| `WithWeaklyTypedInput()` | 启用弱类型转换 |
| `WithDecodeNil()` | 在解码中包含 nil 值 |

## 自定义解码器

高级场景可创建可复用的解码器：

```go
var result User
decoder, err := mapx.NewDecoder(&result, mapx.WithTagName("json"))
if err != nil {
    return err
}
err = decoder.Decode(data)
```

## 解码钩子

该包内置了 VEF 类型的解码钩子：

- `timex.DateTime`、`timex.Date`、`timex.Time` — 自动字符串解析
- `decimal.Decimal` — 自动 decimal 转换

这些钩子会自动注册，因此你可以将包含字符串时间戳的 map 解码为带有 `timex` 字段的结构体，无需任何额外配置。
