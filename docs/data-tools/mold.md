---
sidebar_position: 3
---

# Mold

The `mold` package is a struct transformation engine that modifies field values based on struct tags. It operates at both field and struct levels.

## How It Works

The `mold` tag on struct fields triggers transformation functions. CRUD query
actions run the transformer on `find_one`, `find_all`, `find_page`,
`find_tree`, and `export` results before they are returned, so response models
can expose derived or translated fields.

### Built-in: Code Set Translation

The built-in `translate` transformer resolves a source field through registered
`Translator` implementations and writes the result to a sibling `<Field>Name`
field. The framework ships one built-in translator: `CodeSetTranslator`,
which handles only `codes:` kinds such as `mold:"translate=codes:status"`.

> **Renamed in v0.39**: the "dictionary" vocabulary became "code set". The tag
> prefix changed from `dict:` to `codes:`, and every `Dictionary*` identifier
> became `CodeSet*` (see the mapping table below).

```go
type Order struct {
    Status     string `json:"status" mold:"translate=codes:status"`
    StatusName string `json:"statusName" bun:",scanonly"`
}
```

When a query result contains `Status = "active"`, the transformer asks the
code set resolver for code set `status` and code `active`, then writes the
display name to `StatusName`.

Auditing models such as `orm.FullAuditedModel` use `mold:"translate=user?"` on
`CreatedBy` and `UpdatedBy`. That tag is an optional hook for a custom user
translator; it is not provided by the built-in code set translator.

## Interfaces

### Transformer

```go
type Transformer interface {
    Struct(ctx context.Context, value any) error
    Field(ctx context.Context, value any, tags string) error
}
```

`Transformer.Struct` requires a non-nil pointer to a struct. Passing a nil
value, a non-pointer, a nil pointer, a pointer to a non-struct, or a
`time.Time` value returns an error. `Transformer.Field` requires a non-nil
pointer unless the tag string is empty or `"-"`, in which case it is a no-op.

### `FieldTransformer`

Implement custom field-level transformations:

```go
type FieldTransformer interface {
    Tag() string
    Transform(ctx context.Context, fl FieldLevel) error
}
```

### `StructTransformer`

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
| `Struct()` | `reflect.Value` | Struct that contains the current field; may be invalid when transforming a standalone field |

`StructLevel` exposes `Transformer()`, `Parent()`, and `Struct()` for
struct-level transformers.

Function adapters are also public:

| Adapter | Purpose |
| --- | --- |
| `mold.Func` | use a plain function as a field transformer implementation |
| `mold.StructLevelFunc` | use a plain function for struct-level transformation |
| `mold.InterceptorFunc` | use a plain function as an `Interceptor` |

## Tag Format

```
mold:"function=param"
```

Multiple transformations:

```
mold:"function1=param1,function2=param2"
```

`mold:"-"` skips a field. `dive` recurses into slice, array, or map values.
For maps, `dive,keys,...,endkeys,...` applies the tags between `keys` and
`endkeys` to map keys and the remaining tags to map values. Nested struct
fields are traversed automatically, but slice and map elements are transformed
only when `dive` is present. Commas inside a parameter must be escaped as
`0x2C`.

### Built-in: Expression-Derived Fields

The core runtime registers an `expr` field transformer backed by
`expression.Engine`. It evaluates the expression against the containing struct
and decodes the result into the tagged field:

```go
type LineItem struct {
    Price float64 `json:"price"`
    Qty   float64 `json:"qty"`
    Total float64 `json:"total" mold:"expr=price * qty"`
}
```

Fields are evaluated in declaration order, so derived fields can reference
sibling fields declared above them. If an expression contains a comma, escape it
as `0x2C` inside the mold tag. See [Expression Engine](./expression) for the
full API.

The `expr` tag is provided by the expression module through the
`vef:mold:field_transformers` group. It is not provided by the `mold` module
alone. The `mold` module itself contributes the built-in `translate` field
transformer and the `CodeSetTranslator`; other field transformers must be
registered through the same group or by constructing a custom transformer.

## Code Set Resolution

The `translate` transformer resolves field values through the `Translator`
interface. The framework ships one built-in translator — `CodeSetTranslator`
— that handles `kind` strings prefixed with `codes:` (for example,
`mold:"translate=codes:gender"`). If the kind is `codes:status?`, the built-in
translator still supports the full string and resolves code set key
`status?`; it does not strip the `?` suffix.

Supported source field shapes are `string`, `*string`, signed and unsigned
integer types, pointers to those integer types, `[]string`, and `*[]string`
after mold dereferencing. Scalar targets must be `string` or `*string`; slice
targets must be `[]string` or `*[]string`. The target field is always the
source field name plus `Name` (`<Field>Name`). Empty scalar values are skipped,
nil source slices leave the target untouched, and empty source slices write an
empty target slice.

Custom translators implement:

```go
type Translator interface {
    Supports(kind string) bool
    Translate(ctx context.Context, kind, value string) (string, error)
}
```

The code-set resolver and loader interfaces:

```go
type CodeSetResolver interface {
    Resolve(ctx context.Context, codeSet, code string) (string, error)
}

type CodeSetLoader interface {
    Load(ctx context.Context, codeSet string) (map[string]string, error)
}
```

`CodeSetLoaderFunc` lets a plain function satisfy `CodeSetLoader`.

### Enumerable catalog (optional)

A host whose code sets are enumerable can additionally implement
`mold.CodeSetInspector` alongside its loader (or wholesale-replaced
resolver):

```go
type CodeSetInspector interface {
    ListCodeSets(ctx context.Context) ([]CodeSetInfo, error) // {codeSet, name}
    ListCodes(ctx context.Context, codeSet string) ([]CodeInfo, error) // {code, label}
}
```

Consumers type-assert for it and degrade gracefully when it is absent. The
[integration module](../integration/code-maps#the-host-code-set-catalog) uses
it to power the code-mapping editor's pickers and to validate code map
identifiers.

### What `?` actually means

The `?` suffix in `mold:"translate=user?"` makes the lookup **silently skip**
when **no translator supports the full `kind` string**. If a translator matches
but its `Translate` call returns an error, the error is still propagated — the
`?` is not a "swallow all errors" switch.

So `translate=user?` requires that you register a custom `Translator` whose
`Supports("user?")` returns true if you want it to run. Without one, the field
is left untouched and no error is returned. A required kind such as
`translate=user` returns an error when no translator supports it.

## Cached Resolution

`CachedCodeSetResolver` wraps a `CodeSetLoader` (not a `CodeSetResolver`) with in-process caching, and subscribes to `mold.CodeSetChangedEvent` for invalidation:

```go
resolver := mold.NewCachedCodeSetResolver(loader, bus)
```

`NewCachedCodeSetResolver` panics if the `CodeSetLoader` or `event.Bus`
is nil. The cache holds entire code sets keyed by the loader's code set and
merges concurrent loads for the same key. `Resolve` returns an empty string
without error for an empty code set, an empty code, or a code that is not
present in the loaded set.

When the data underlying a code set changes, publish
`mold.CodeSetChangedEvent{Keys: []string{"..."}}` through the event bus to
invalidate the matching cache entry.

You can publish the same event through the helper:

```go
err := mold.PublishCodeSetChangedEvent(ctx, bus, "gender", "status")
```

`CodeSetChangedEvent.EventType()` returns the framework event type used by
the cache invalidation subscriber (`vef.translate.code_set.changed`).

Calling `PublishCodeSetChangedEvent(ctx, bus)` without keys asks subscribers
to clear their entire code set cache.

The public APIs in this cache path are `CachedCodeSetResolver`,
`CodeSetChangedEvent`, `CodeSetChangedEvent.Keys`,
`PublishCodeSetChangedEvent`, and `CachedCodeSetResolver.Resolve`, which
implements `CodeSetResolver.Resolve`.

### v0.38 → v0.39 rename map

| Old identifier | New identifier |
| --- | --- |
| tag `mold:"translate=dict:xxx"` | `mold:"translate=codes:xxx"` |
| `DictionaryTranslator` | `CodeSetTranslator` |
| `DictionaryResolver` | `CodeSetResolver` |
| `DictionaryLoader` / `DictionaryLoaderFunc` | `CodeSetLoader` / `CodeSetLoaderFunc` |
| `CachedDictionaryResolver` / `NewCachedDictionaryResolver` | `CachedCodeSetResolver` / `NewCachedCodeSetResolver` |
| `DictionaryChangedEvent` / `PublishDictionaryChangedEvent` | `CodeSetChangedEvent` / `PublishCodeSetChangedEvent` |
