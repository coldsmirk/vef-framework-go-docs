---
sidebar_position: 13
---

# Mold

`mold` 包是一个结构体变换引擎，基于结构体标签修改字段值。它在字段和结构体两个层级都可以操作。

## 工作原理

结构体字段上的 `mold` 标签会触发变换函数。框架在查询结果返回时自动运行 mold 变换器来丰富响应数据。

### 内置：用户名翻译

最常见的内置用法是将用户 ID 翻译为显示名称：

```go
type User struct {
    orm.FullAuditedModel
    // CreatedBy     string `mold:"translate=user?"`  ← 从 FullAuditedModel 继承
    // CreatedByName string `bun:",scanonly"`          ← 由 mold 填充
}
```

当查询结果包含 `CreatedBy = "user-123"` 时，mold 变换器查找用户数据字典并设置 `CreatedByName = "Alice"`。

## 接口

### Transformer

```go
type Transformer interface {
    Struct(ctx context.Context, value any) error
    Field(ctx context.Context, value any, tags string) error
}
```

### FieldTransformer

实现自定义字段级变换：

```go
type FieldTransformer interface {
    Tag() string
    Transform(ctx context.Context, fl FieldLevel) error
}
```

### StructTransformer

实现自定义结构体级变换：

```go
type StructTransformer interface {
    Transform(ctx context.Context, sl StructLevel) error
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

## 标签格式

```
mold:"function=param"
```

多个变换：

```
mold:"function1=param1,function2=param2"
```

## 字典翻译

`translate` 变换器通过 `Translator` 接口翻译字段值。框架内置了一个 `DictionaryTranslator`，负责处理形如 `kind = "dict:xxx"` 的翻译（例如 `mold:"translate=dict:gender"`）。

自定义 translator 实现：

```go
type Translator interface {
    Supports(kind string) bool
    Translate(ctx context.Context, kind, value string) (string, error)
}
```

字典风格的 resolver / loader 接口：

```go
type DictionaryResolver interface {
    Resolve(ctx context.Context, key, code string) (string, error)
}

type DictionaryLoader interface {
    Load(ctx context.Context, key string) (map[string]string, error)
}
```

### `?` 的真实含义

`mold:"translate=user?"` 中的 `?` 表示：当**没有任何 translator 声明 `Supports("user") == true`** 时，安静地跳过这次翻译；如果有匹配的 translator，但它的 `Translate` 内部返回了 error，error 仍会向上抛 —— `?` 不是"吞掉一切错误"。

也就是说，`translate=user?` 仍要求你在容器里注册了一个 `Supports("user")` 为真的 `Translator`；没有注册时字段保持原值（不会报错）。

## 带缓存的解析器

`CachedDictionaryResolver` 包装的是 `DictionaryLoader`（不是 `DictionaryResolver`），并通过订阅 `mold.DictionaryChangedEvent` 自动失效：

```go
resolver := mold.NewCachedDictionaryResolver(loader, bus)
```

缓存按 loader 的 `key` 缓存整张字典。当某个字典的底层数据发生变化时，通过事件总线发布 `mold.DictionaryChangedEvent{Key: "..."}` 来让对应缓存失效。
