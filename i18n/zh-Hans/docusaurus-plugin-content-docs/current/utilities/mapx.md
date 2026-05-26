---
sidebar_position: 6
---

# Mapx

`mapx` 包提供 Go 结构体与 `map[string]any` 之间的双向转换，底层基于 `github.com/go-viper/mapstructure/v2`。VEF 覆盖了上游的默认 tag —— 框架默认使用 `json` tag。

## 结构体转 Map

```go
import "github.com/coldsmirk/vef-framework-go/mapx"

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
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
// 切到其他 tag，例如 yaml
m, err := mapx.ToMap(user, mapx.WithTagName("yaml"))

// 弱类型转换（字符串 "123" → int 123）
user, err := mapx.FromMap[User](data, mapx.WithWeaklyTypedInput())

// 把源 map 里有但结构体里没有的字段当错误抛出
user, err := mapx.FromMap[User](data, mapx.WithErrorUnused())
```

### 可用选项

| 选项 | 效果 |
| --- | --- |
| `WithTagName(tag)` | 覆盖 mapx 读取的结构体 tag（**默认 `json`**）。 |
| `WithIgnoreUntaggedFields()` | 跳过没有当前 tag 的字段。 |
| `WithDecodeHook(hooks...)` | 追加额外的 decode hook（内置 hook 仍然生效）。 |
| `WithMatchName(fn)` | 自定义字段名匹配函数（默认大小写不敏感的 camelCase 比较）。 |
| `WithErrorUnused()` | 源 map 里有结构体没有的字段时报错。 |
| `WithErrorUnset()` | 结构体里有源 map 没有的字段时报错。 |
| `WithZeroFields()` | 解码前清零目标结构体字段。 |
| `WithAllowUnsetPointer()` | 允许指针字段保持 nil 而不被初始化。 |
| `WithMetadata(m)` | 把 "unused" / "unset" 键收集到 `mapstructure.Metadata`。 |
| `WithWeaklyTypedInput()` | 常见类型间互转（string ↔ number ↔ bool …）。 |
| `WithDecodeNil()` | 把 `nil` 源值送入解码管线（默认会被跳过）。 |

## 自定义解码器

高级场景可创建可复用的解码器：

```go
var result User
decoder, err := mapx.NewDecoder(&result, mapx.WithTagName("yaml"))
if err != nil {
    return err
}
err = decoder.Decode(data)
```

## 解码钩子

`mapx` 在 `NewDecoder` 内置注册了一整套 decode hook，使得来自 JSON、表单、环境变量的纯字符串 map 可以直接解码进类型化结构体：

- `time.Time` —— 解析 `"2006-01-02 15:04:05"`（Go 的 `time.DateTime` 布局）
- `time.Location` —— 解析 IANA 名称（例如 `"Asia/Shanghai"`）
- `time.Duration` —— 解析 Go 持续时间字符串（例如 `"5m"`）
- `*url.URL` —— 解析 URL
- `net.IP` / `net.IPNet` / `netip.Addr` / `netip.AddrPort` / `netip.Prefix`
- `json.RawMessage` —— 原值透传
- `*multipart.FileHeader` —— 源是 `[]*multipart.FileHeader` 时取第一项（让单文件上传字段自然工作）
- `collections.Set` / `SortedSet` / `ConcurrentSet` / `ConcurrentSortedSet` —— 把切片转为对应的集合类型
- `encoding.TextUnmarshaler` —— 任何实现了 `UnmarshalText` 的类型
- string → 基础类型的隐式转换（int / uint / float / bool）

`timex.DateTime` / `timex.Date` / `timex.Time` 是基于 `time.Time` 的命名类型，能否命中上述 `time.Time` hook 取决于 mapstructure 对底层类型的处理。如果依赖这类自动解码，请按场景实际验证一遍。

若需注册自己的 hook，用 `mapx.WithDecodeHook(myHook)` 追加 —— 框架内置 hook 仍然生效。
