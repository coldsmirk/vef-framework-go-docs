---
sidebar_position: 5
---

# Mapx

`mapx` 包提供 Go 结构体与 `map[string]any` 之间的双向转换，底层基于 `github.com/go-viper/mapstructure/v2`。VEF 覆盖了上游的默认 tag —— 框架默认使用 `json` tag。

## API 参考

| API | Contract |
| --- | --- |
| `mapx.DecoderHook` | 默认 decoder 使用的公开 composed `mapstructure.DecodeHookFunc` |
| `mapx.DecoderOption` | 修改 `mapstructure.DecoderConfig` 的 function option 类型 |
| `mapx.Metadata` | `mapstructure.Metadata` 的 alias |
| `mapx.NewDecoder(result, options...)` | 使用 VEF 默认值创建 `mapstructure.Decoder`，然后按顺序应用 options |
| `mapx.ToMap(value, options...)` | 将 struct 或 pointer-to-struct 转成 `map[string]any`；非 struct 输入返回 `ErrInvalidToMapValue` |
| `mapx.FromMap[T](value, options...)` | 将 `map[string]any` 转成 `*T`；非 struct `T` 返回 `ErrInvalidFromMapType` |
| `mapx.WithTagName(tagName)` | 设置 `DecoderConfig.TagName`；默认是 `json` |
| `mapx.WithIgnoreUntaggedFields(ignore)` | 把 `DecoderConfig.IgnoreUntaggedFields` 设为传入的 boolean |
| `mapx.WithDecodeHook(decodeHook)` | 替换 `DecoderConfig.DecodeHook`；如需保留默认 hook，需要自行和 `mapx.DecoderHook` compose |
| `mapx.WithMatchName(matchName)` | 替换 key/field matcher；默认是 `mapKey == lo.CamelCase(fieldName)` |
| `mapx.WithErrorUnused()` | 设置 `ErrorUnused = true` |
| `mapx.WithErrorUnset()` | 设置 `ErrorUnset = true` |
| `mapx.WithZeroFields()` | 设置 `ZeroFields = true` |
| `mapx.WithAllowUnsetPointer()` | 设置 `AllowUnsetPointer = true` |
| `mapx.WithMetadata(metadata)` | 将 decode metadata 写入传入的 `*mapx.Metadata` |
| `mapx.WithWeaklyTypedInput()` | 设置 `WeaklyTypedInput = true` |
| `mapx.WithDecodeNil()` | 设置 `DecodeNil = true` |
| `mapx.ErrInvalidToMapValue` | `ToMap` 输入不是 struct 或 pointer-to-struct 时的 sentinel |
| `mapx.ErrInvalidFromMapType` | `FromMap[T]` 的 `T` 不是 struct 时的 sentinel |
| `mapx.ErrCollectionSetNilElement` | 解码到 collection set 时遇到 nil element 的 sentinel |
| `mapx.ErrCollectionSetIncompatibleKind` | collection set element 发生 string/numeric family mismatch 的 sentinel |
| `mapx.ErrCollectionSetOverflow` | collection set element 数值转换溢出的 sentinel |
| `mapx.ErrCollectionSetNonInteger` | fractional float 目标是 integer set element 时的 sentinel |
| `mapx.ErrCollectionSetNotFinite` | NaN 或 infinity 目标是 integer set element 时的 sentinel |
| `mapx.ErrCollectionSetNegative` | 负数目标是 unsigned set element 时的 sentinel |
| `mapx.ErrCollectionSetUnsupportedTarget` | collection set element kind 没有转换策略时的 sentinel |
| `mapx.ErrJSONNumberNotInteger` | 小数或指数形式的 `json.Number` 落到整数字段时的 sentinel |
| `mapx.ErrJSONNumberOverflow` | `json.Number` 超出数值目标类型范围时的 sentinel |

decoder hook 链会翻译数字保真 JSON 解析产生的 `json.Number` 值：
数值目标按精确位数解析、严格性对齐 `encoding/json`，`json.Number` /
`json.RawMessage` 目标保留字面量，其余目标——最重要的是 `any`——看到
`float64`，动态消费者的既有运行时契约不变。

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
| `WithIgnoreUntaggedFields(ignore)` | 设置是否跳过没有当前 tag 的字段。 |
| `WithDecodeHook(hook)` | 替换默认 decode hook。 |
| `WithMatchName(fn)` | 自定义字段名匹配函数（默认与 `lo.CamelCase(fieldName)` 精确比较）。 |
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
- `json.RawMessage` —— 将源值 marshal 成 JSON bytes
- `*multipart.FileHeader` —— 源是长度为 1 的 `[]*multipart.FileHeader` 时取唯一一项
- `collections.Set` / `SortedSet` / `ConcurrentSet` / `ConcurrentSortedSet` —— 把 slice 或 array 转为对应的集合类型
- `encoding.TextUnmarshaler` —— 任何实现了 `UnmarshalText` 的类型
- string → 基础类型的隐式转换（int / uint / float / bool）

collection-set 解码为 `string`、有符号整数、无符号整数、`float32` 和
`float64` 注册。它会拒绝 nil element、string/numeric family mismatch、
numeric overflow、fractional float 转 integer set、NaN 或 infinity 转
integer set，以及负数转 unsigned set。

`timex.DateTime` / `timex.Date` / `timex.Time` 是基于 `time.Time` 的命名类型，能否命中上述 `time.Time` hook 取决于 mapstructure 对底层类型的处理。如果依赖这类自动解码，请按场景实际验证一遍。

`WithDecodeHook(myHook)` 会替换默认 composed hook。若要扩展默认行为，请先把
自定义 hook 和 `mapx.DecoderHook` compose，再传给 `WithDecodeHook`。

组合后的默认 hook 也公开为 `mapx.DecoderHook`，metadata 收集使用公开别名
`mapx.Metadata`。

## 错误哨兵

| 错误 | 含义 |
| --- | --- |
| `ErrInvalidToMapValue` | `ToMap` 收到的不是 struct |
| `ErrInvalidFromMapType` | `FromMap[T]` 的 `T` 不是 struct |
| `ErrCollectionSetNilElement` | nil 元素不能插入 collection set |
| `ErrCollectionSetIncompatibleKind` | 源值 kind 与 set 元素 kind 不匹配 |
| `ErrCollectionSetOverflow` | 数值源会溢出目标 set 元素类型 |
| `ErrCollectionSetNonInteger` | 带小数的 float 解码到整数 set 会丢失信息 |
| `ErrCollectionSetNotFinite` | NaN 或 infinity 不能解码到整数 set |
| `ErrCollectionSetNegative` | 负数不能解码到无符号 set 元素 |
| `ErrCollectionSetUnsupportedTarget` | 目标 set 元素 kind 没有转换策略 |
