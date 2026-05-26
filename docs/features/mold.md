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

## Dictionary Resolution

The `translate` transformer resolves field values through the `Translator` interface. The framework ships one built-in translator — `DictionaryTranslator` — that handles `kind` strings prefixed with `dict:` (e.g. `mold:"translate=dict:gender"`).

Custom translators implement:

```go
type Translator interface {
    Supports(kind string) bool
    Translate(ctx context.Context, kind, value string) (string, error)
}
```

The dictionary-style resolver and loader interfaces:

```go
type DictionaryResolver interface {
    Resolve(ctx context.Context, key, code string) (string, error)
}

type DictionaryLoader interface {
    Load(ctx context.Context, key string) (map[string]string, error)
}
```

### What `?` actually means

The `?` suffix in `mold:"translate=user?"` makes the lookup **silently skip** when **no translator supports the `kind`** (here, `user`). If a translator matches but its `Translate` call returns an error, the error is still propagated — the `?` is not a "swallow all errors" switch.

So `translate=user?` requires that you register a `Translator` whose `Supports("user")` returns true. Without one, the field is left untouched (no error).

## Cached Resolution

`CachedDictionaryResolver` wraps a `DictionaryLoader` (not a `DictionaryResolver`) with in-process caching, and subscribes to `mold.DictionaryChangedEvent` for invalidation:

```go
resolver := mold.NewCachedDictionaryResolver(loader, bus)
```

The cache holds entire dictionaries keyed by the loader's `key`. When the data underlying a dictionary changes, publish `mold.DictionaryChangedEvent{Key: "..."}` through the event bus to invalidate the matching cache entry.
