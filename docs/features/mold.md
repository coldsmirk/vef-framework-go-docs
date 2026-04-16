---
sidebar_position: 13
---

# Mold

The `mold` package is a struct transformation engine that modifies field values based on struct tags. It operates at both field and struct levels.

## How It Works

The `mold` tag on struct fields triggers transformation functions. The framework automatically runs the mold transformer on query results to enrich response data.

### Built-in: User Name Translation

The most common built-in usage is translating user IDs into display names:

```go
type User struct {
    orm.FullAuditedModel
    // CreatedBy     string `mold:"translate=user?"`  ← inherited from FullAuditedModel
    // CreatedByName string `bun:",scanonly"`          ← populated by mold
}
```

When a query result contains `CreatedBy = "user-123"`, the mold transformer looks up the user data dictionary and sets `CreatedByName = "Alice"`.

## Interfaces

### Transformer

```go
type Transformer interface {
    Struct(ctx context.Context, value any) error
    Field(ctx context.Context, value any, tags string) error
}
```

### FieldTransformer

Implement custom field-level transformations:

```go
type FieldTransformer interface {
    Tag() string
    Transform(ctx context.Context, fl FieldLevel) error
}
```

### StructTransformer

Implement custom struct-level transformations:

```go
type StructTransformer interface {
    Transform(ctx context.Context, sl StructLevel) error
}
```

### Interceptor

Redirect transformation to inner values (e.g., `sql.NullString` → its inner string):

```go
type Interceptor interface {
    Intercept(current reflect.Value) (inner reflect.Value)
}
```

## FieldLevel API

Inside a field transformer, `FieldLevel` provides:

| Method | Returns | Purpose |
| --- | --- | --- |
| `Transformer()` | `Transformer` | Access the parent transformer |
| `Name()` | `string` | Current field name |
| `Parent()` | `reflect.Value` | Parent struct value |
| `Field()` | `reflect.Value` | Current field value |
| `Param()` | `string` | Parameter from tag (e.g., `user?` in `translate=user?`) |
| `SiblingField(name)` | `reflect.Value, bool` | Access sibling field by name |

## Tag Format

```
mold:"function=param"
```

Multiple transformations:

```
mold:"function1=param1,function2=param2"
```

## Data Dictionary Resolution

The `translate` transformer resolves field values through a data dictionary:

```go
type DataDictResolver interface {
    Resolve(ctx context.Context, dictType string, keys []string) (map[string]string, error)
}
```

The `?` suffix in `translate=user?` means the translation is optional — if the lookup fails, the original value is kept instead of returning an error.

## Cached Resolution

The `CachedDataDictResolver` wraps a `DataDictResolver` with in-request caching to avoid redundant lookups when the same dict type is resolved multiple times in a single request.
