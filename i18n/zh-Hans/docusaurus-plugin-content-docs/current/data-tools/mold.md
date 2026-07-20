---
sidebar_position: 3
---

# Mold

`mold` 包是一个结构体变换引擎，基于结构体标签修改字段值。它在字段和结构体两个层级都可以操作。

## 工作原理

结构体字段上的 `mold` 标签会触发变换函数。CRUD 查询动作会在
`find_one`、`find_all`、`find_page`、`find_tree` 和 `export` 结果返回前运行
transformer，因此响应模型可以暴露派生字段或翻译字段。

### 内置：码值集翻译

内置 `translate` transformer 会通过已注册的 `Translator` 实现解析源字段，并把
结果写入同级 `<Field>Name` 字段。框架自带一个内置 translator：
`CodeSetTranslator`，它只处理 `mold:"translate=codes:status"` 这样的
`codes:` kind。

> **v0.39 重命名**："字典（dictionary）"词汇统一改为"码值集（code set）"。
> 标签前缀由 `dict:` 变为 `codes:`，所有 `Dictionary*` 标识符变为
> `CodeSet*`（见下方对照表）。

```go
type Order struct {
    Status     string `json:"status" mold:"translate=codes:status"`
    StatusName string `json:"statusName" bun:",scanonly"`
}
```

当查询结果包含 `Status = "active"` 时，transformer 会向码值集 resolver 查询
码值集 `status`、code `active`，然后把显示名写入 `StatusName`。

`orm.FullAuditedModel` 这类审计模型会在 `CreatedBy` 和 `UpdatedBy` 上使用
`mold:"translate=user?"`。这个标签是给自定义用户 translator 预留的可选 hook；
内置码值集 translator 并不提供 `user` 翻译。

## 接口

### Transformer

```go
type Transformer interface {
    Struct(ctx context.Context, value any) error
    Field(ctx context.Context, value any, tags string) error
}
```

`Transformer.Struct` 要求传入非 nil 的 struct 指针。传入 nil、非指针、nil
指针、指向非 struct 的指针，或 `time.Time` 值都会返回 error。
`Transformer.Field` 要求传入非 nil 指针；但 tag 字符串为空或为 `"-"` 时是
no-op，不会检查字段值。

### `FieldTransformer`

实现自定义字段级变换：

```go
type FieldTransformer interface {
    Tag() string
    Transform(ctx context.Context, fl FieldLevel) error
}
```

### `StructTransformer`

实现自定义结构体级变换：

```go
type StructTransformer interface {
    Transform(ctx context.Context, sl StructLevel) error
}
```

### Interceptor

把变换重定向到内部值（例如 `sql.NullString` → 内部 string）：

```go
type Interceptor interface {
    Intercept(current reflect.Value) (inner reflect.Value)
}
```

## FieldLevel API

在字段变换器内部，`FieldLevel` 提供：

| 方法 | 返回值 | 用途 |
| --- | --- | --- |
| `Transformer()` | `Transformer` | 访问父级变换器 |
| `Name()` | `string` | 当前字段名 |
| `Parent()` | `reflect.Value` | 父级结构体值 |
| `Field()` | `reflect.Value` | 当前字段值 |
| `Param()` | `string` | 标签中的参数（如 `translate=user?` 中的 `user?`）|
| `SiblingField(name)` | `reflect.Value, bool` | 按名称访问同级字段 |
| `Struct()` | `reflect.Value` | 当前字段所在的 struct；单独变换字段时可能是 invalid value |

`StructLevel` 提供 `Transformer()`、`Parent()` 和 `Struct()`，供结构体级
transformer 使用。

函数适配器也是公开 API：

| 适配器 | 作用 |
| --- | --- |
| `mold.Func` | 用普通函数实现字段级 transformer |
| `mold.StructLevelFunc` | 用普通函数实现结构体级 transformer |
| `mold.InterceptorFunc` | 用普通函数实现 `Interceptor` |

## 标签格式

```
mold:"function=param"
```

多个变换：

```
mold:"function1=param1,function2=param2"
```

`mold:"-"` 会跳过字段。`dive` 会递归处理 slice、array 或 map。对 map 来说，
`dive,keys,...,endkeys,...` 会把 `keys` 和 `endkeys` 之间的标签应用到 map key，
把后续标签应用到 map value。嵌套 struct 字段会自动遍历；slice 和 map 元素
只有写了 `dive` 才会被 transform。参数中如果包含逗号，需要写成 `0x2C`。

### 内置：表达式派生字段

core runtime 会注册一个由 `expression.Engine` 支持的 `expr` field
transformer。它会把当前结构体作为表达式环境，执行表达式，并把解码后的
结果写入带标签的字段：

```go
type LineItem struct {
    Price float64 `json:"price"`
    Qty   float64 `json:"qty"`
    Total float64 `json:"total" mold:"expr=price * qty"`
}
```

字段按声明顺序执行，所以派生字段可以引用声明在它上方的 sibling field。
如果表达式里包含逗号，需要在 mold tag 中写成 `0x2C`。完整 API 见
[表达式引擎](./expression)。

`expr` tag 由 expression module 通过 `vef:mold:field_transformers` group 提供，
不是 `mold` module 单独提供的默认 tag。`mold` module 自身提供内置
`translate` field transformer 和 `CodeSetTranslator`；其他 field transformer
需要通过同一个 group 注册，或由应用自行构造 custom transformer。

## 码值集翻译

`translate` 变换器通过 `Translator` 接口翻译字段值。框架内置了一个
`CodeSetTranslator`，负责处理 `codes:` 前缀的 kind（例如
`mold:"translate=codes:gender"`）。如果 kind 是 `codes:status?`，内置 translator
仍会认为完整字符串被支持，并解析码值集 key `status?`；它不会自动剥掉 `?` 后缀。

支持的源字段类型包括 `string`、`*string`、有符号和无符号整数、这些整数的指针、
`[]string`，以及经 mold 解引用后的 `*[]string`。标量目标字段必须是 `string`
或 `*string`；slice 目标字段必须是 `[]string` 或 `*[]string`。目标字段始终是
源字段名加 `Name`（`<Field>Name`）。空标量值会跳过，nil 源 slice 会保持目标
字段不变，空源 slice 会写入空目标 slice。

自定义 translator 实现：

```go
type Translator interface {
    Supports(kind string) bool
    Translate(ctx context.Context, kind, value string) (string, error)
}
```

码值集的 resolver / loader 接口：

```go
type CodeSetResolver interface {
    Resolve(ctx context.Context, codeSet, code string) (string, error)
}

type CodeSetLoader interface {
    Load(ctx context.Context, codeSet string) (map[string]string, error)
}
```

`CodeSetLoaderFunc` 可以让普通函数满足 `CodeSetLoader` 接口。

### 可枚举目录（可选）

码值集可枚举的宿主可以在其 loader（或整体替换的 resolver）之上额外实现
`mold.CodeSetInspector`：

```go
type CodeSetInspector interface {
    ListCodeSets(ctx context.Context) ([]CodeSetInfo, error) // {codeSet, name}
    ListCodes(ctx context.Context, codeSet string) ([]CodeInfo, error) // {code, label}
}
```

消费方对其做类型断言，缺失时优雅降级。
[集成模块](../integration/code-maps#宿主码值目录)用它驱动码值映射编辑器的
选择器，并校验码值映射的标识。

### `?` 的真实含义

`mold:"translate=user?"` 中的 `?` 表示：当**没有任何 translator 支持完整 kind
字符串**时，安静地跳过这次翻译；如果有匹配的 translator，但它的 `Translate`
内部返回了 error，error 仍会向上抛 —— `?` 不是“吞掉一切错误”。

也就是说，如果你希望 `translate=user?` 真正执行，需要在容器里注册一个
`Supports("user?")` 为真的自定义 `Translator`；没有注册时字段保持原值（不会
报错）。`translate=user` 这样的必需 kind 在没有 translator 支持时会返回 error。

## 带缓存的解析器

`CachedCodeSetResolver` 包装的是 `CodeSetLoader`（不是 `CodeSetResolver`），并通过订阅 `mold.CodeSetChangedEvent` 自动失效：

```go
resolver := mold.NewCachedCodeSetResolver(loader, bus)
```

`NewCachedCodeSetResolver` 在 `CodeSetLoader` 或 `event.Bus` 为 nil 时会 panic。
缓存按 loader 的码值集缓存整套码值，并会合并同一 key 的并发加载。`Resolve` 在
码值集为空、code 为空，或集合里找不到 code 时，都会无 error 地返回空字符串。

当某个码值集的底层数据发生变化时，通过事件总线发布
`mold.CodeSetChangedEvent{Keys: []string{"..."}}` 来让对应缓存失效。

也可以使用辅助函数发布同一事件：

```go
err := mold.PublishCodeSetChangedEvent(ctx, bus, "gender", "status")
```

`CodeSetChangedEvent.EventType()` 返回缓存失效 subscriber 使用的框架事件类型
（`vef.translate.code_set.changed`）。

调用 `PublishCodeSetChangedEvent(ctx, bus)` 但不传 key 时，表示要求订阅方
清空整个码值集缓存。

这条缓存路径上的公开 API 是 `CachedCodeSetResolver`、
`CodeSetChangedEvent`、`CodeSetChangedEvent.Keys`、
`PublishCodeSetChangedEvent` 和 `CachedCodeSetResolver.Resolve`；其中
`CachedCodeSetResolver.Resolve` 实现 `CodeSetResolver.Resolve`。

### v0.38 → v0.39 重命名对照

| 旧标识符 | 新标识符 |
| --- | --- |
| 标签 `mold:"translate=dict:xxx"` | `mold:"translate=codes:xxx"` |
| `DictionaryTranslator` | `CodeSetTranslator` |
| `DictionaryResolver` | `CodeSetResolver` |
| `DictionaryLoader` / `DictionaryLoaderFunc` | `CodeSetLoader` / `CodeSetLoaderFunc` |
| `CachedDictionaryResolver` / `NewCachedDictionaryResolver` | `CachedCodeSetResolver` / `NewCachedCodeSetResolver` |
| `DictionaryChangedEvent` / `PublishDictionaryChangedEvent` | `CodeSetChangedEvent` / `PublishCodeSetChangedEvent` |
