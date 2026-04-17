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

## 数据字典解析

`translate` 变换器通过数据字典解析字段值：

```go
type DataDictResolver interface {
    Resolve(ctx context.Context, dictType string, keys []string) (map[string]string, error)
}
```

`translate=user?` 中的 `?` 后缀表示翻译是可选的——如果查找失败，保留原始值而不是返回错误。
