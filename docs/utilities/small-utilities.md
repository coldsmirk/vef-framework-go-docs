---
sidebar_position: 8
---

# Small Utilities

This page documents lightweight utility packages that provide common helper functions.

## `ptr` — Pointer Helpers

Generic utility functions for safe pointer operations.

The public `ptr` surface has 6 exported functions. It has no exported types,
no exported fields, and no exported methods.

| API | Signature | Purpose |
| --- | --- | --- |
| `Of` | `ptr.Of[T comparable](v T) *T` | Returns a pointer to `v`, or `nil` when `v` is the zero value. |
| `Zero` | `ptr.Zero[T any]() T` | Returns the zero value for `T`. |
| `Value` | `ptr.Value[T any](p *T, fallbacks ...*T) T` | Dereferences `p`, then the first non-nil fallback, then returns zero. |
| `ValueOrElse` | `ptr.ValueOrElse[T any](p *T, fn func() T) T` | Dereferences `p`, or calls `fn` only when `p` is nil. |
| `Equal` | `ptr.Equal[T comparable](a *T, b *T) bool` | Compares pointers by pointed value, with nil-aware behavior. |
| `Coalesce` | `ptr.Coalesce[T any](ptrs ...*T) *T` | Returns the first non-nil pointer itself. |

`ptr.Of` requires `T comparable` because it compares `v` with the zero value.
It returns `nil` for zero values such as `""`, `0`, and `false`; for non-zero
values it returns a pointer to a copy of `v`.

`ptr.Value` checks fallbacks from left to right. A non-nil primary pointer wins
over all fallbacks. A non-nil fallback pointing at a zero value is still used.
When every pointer is nil, it returns `ptr.Zero[T]()`.

`ptr.ValueOrElse` is lazy: `fn` is not called when `p` is non-nil. If `p` is
nil, `fn` must be non-nil because the helper calls it directly.

`ptr.Equal` returns true when both pointers are nil, false when exactly one is
nil, and otherwise compares `*a == *b`.

`ptr.Coalesce` returns the exact first non-nil pointer, not a copy. With no
arguments, or with only nil arguments, it returns nil.

Examples:

```go
import "github.com/coldsmirk/vef-framework-go/ptr"

p := ptr.Of("hello")   // *string → "hello"
p = ptr.Of("")          // *string → nil (zero value)
p = ptr.Of(42)          // *int → 42
p = ptr.Of(0)           // *int → nil

s := ptr.Value(p)                      // Dereference, or zero if nil
s = ptr.Value(p, fallback1, fallback2) // Try each fallback in order

s := ptr.ValueOrElse(p, func() string {
    return computeDefault()
})

z := ptr.Zero[string]() // ""
ptr.Equal(a, b)  // Both nil → true, one nil → false, else compare values
result := ptr.Coalesce(p1, p2, p3) // First non-nil pointer, or nil
```

## Other Public Utility Packages

These packages expose small focused helpers that do not need a full feature page.

### `page`

Pagination helpers for request parameters and responses.

Public surface:

| API | Signature / wire shape |
| --- | --- |
| `DefaultPageNumber` | `int = 1` |
| `DefaultPageSize` | `int = 15` |
| `MaxPageSize` | `int = 1000` |
| `Pageable` | `type Pageable struct { Page int json:"page"; Size int json:"size" }` |
| `Pageable.Page` | `int json:"page"` |
| `Pageable.Size` | `int json:"size"` |
| `Page[T]` | `type Page[T any] struct { Page int json:"page"; Size int json:"size"; Total int64 json:"total"; Items []T json:"items" }` |
| `Page.Page` | `int json:"page"` |
| `Page.Size` | `int json:"size"` |
| `Page.Total` | `int64 json:"total"` |
| `Page.Items` | `[]T json:"items"` |
| `New` | `page.New[T any](pageable page.Pageable, total int64, items []T) page.Page[T]` |
| `Pageable.Normalize` | `func (p *Pageable) Normalize(size ...int)` |
| `Pageable.Offset` | `func (p Pageable) Offset() int` |
| `Page.TotalPages` | `func (page Page[T]) TotalPages() int` |
| `Page.HasNext` | `func (page Page[T]) HasNext() bool` |
| `Page.HasPrevious` | `func (page Page[T]) HasPrevious() bool` |

Surface count: 6 top-level exported symbols, 6 exported fields, 5 exported
methods, no exported variables.

Behavior contract:

- `Pageable.Page` is 1-based and `Pageable.Size` is the requested page size
- `Normalize(size...)` mutates the receiver in place
- `Normalize` resets `Page < 1` to `DefaultPageNumber`
- when `Size < 1`, `Normalize` uses only the first optional fallback size, or
  `DefaultPageSize` when no fallback is provided
- values above `MaxPageSize` are clamped to `MaxPageSize`
- `Normalize` does not re-validate a negative custom fallback, so pass a
  positive fallback size
- `Offset()` is a plain `(Page - 1) * Size` calculation and assumes the value
  has already been normalized
- `New` copies `pageable.Page`, `pageable.Size`, and `total` into the response
- `New` converts nil `items` to a non-nil empty slice; non-nil `items` are used
  as provided and are not cloned
- `TotalPages()` returns `0` when `Size == 0`; otherwise it performs ceiling
  division of `Total / Size`
- `HasNext()` compares `Page < TotalPages()`
- `HasPrevious()` returns true when `Page > 1`
- the helper does not validate negative totals, negative sizes after manual
  construction, or inconsistent `Page` values outside `Normalize`

### `sortx`

Ordering primitives used by query builders and API parameters.

Public surface:

| API | Signature / value |
| --- | --- |
| `OrderAsc` | `sortx.OrderDirection = 0` |
| `OrderDesc` | `sortx.OrderDirection = 1` |
| `OrderDirection` | `type OrderDirection int` |
| `OrderDirection.String` | `func (od OrderDirection) String() string` |
| `OrderDirection.MarshalText` | `func (od OrderDirection) MarshalText() ([]byte, error)` |
| `OrderDirection.UnmarshalText` | `func (od *OrderDirection) UnmarshalText(text []byte) error` |
| `OrderDirection.MarshalJSON` | `func (od OrderDirection) MarshalJSON() ([]byte, error)` |
| `OrderDirection.UnmarshalJSON` | `func (od *OrderDirection) UnmarshalJSON(data []byte) error` |
| `ErrInvalidOrderDirection` | exported `error` var |
| `NullsDefault` | `sortx.NullsOrder = 0` |
| `NullsFirst` | `sortx.NullsOrder = 1` |
| `NullsLast` | `sortx.NullsOrder = 2` |
| `NullsOrder` | `type NullsOrder int` |
| `NullsOrder.String` | `func (no NullsOrder) String() string` |
| `OrderSpec` | `type OrderSpec struct` |
| `OrderSpec.Column` | `string` |
| `OrderSpec.Direction` | `sortx.OrderDirection` |
| `OrderSpec.NullsOrder` | `sortx.NullsOrder` |
| `OrderSpec.IsValid` | `func (spec OrderSpec) IsValid() bool` |

Surface count: 9 top-level exported symbols, 3 exported fields, 7 exported
methods, no exported functions.

Behavior contract:

- `OrderDirection.String()` returns `DESC` only for `OrderDesc`; every other
  value renders as `ASC`
- `MarshalText()` lowercases the result of `String()`, so normal values marshal
  as `asc` or `desc`
- `UnmarshalText(text)` trims surrounding whitespace and accepts `asc` /
  `desc` case-insensitively
- invalid text returns an error wrapping `ErrInvalidOrderDirection`
- `MarshalJSON()` delegates to `MarshalText()` and emits a JSON string
- `UnmarshalJSON(data)` requires a JSON string; numbers, booleans, and `null`
  return an `OrderDirection must be a JSON string` error before direction
  parsing
- `NullsDefault.String()` returns an empty string
- `NullsFirst.String()` returns `NULLS FIRST`
- `NullsLast.String()` returns `NULLS LAST`
- any other `NullsOrder` value also returns an empty string
- `OrderSpec.IsValid()` checks only `Column != ""`; it does not validate the
  column as a SQL identifier or validate `Direction` / `NullsOrder`

### `monad`

Inclusive ordered ranges.

Public surface:

| API | Signature / field |
| --- | --- |
| `NewRange` | `monad.NewRange[T cmp.Ordered](start T, end T) monad.Range[T]` |
| `Range[T]` | `type Range[T cmp.Ordered] struct` |
| `Range.Start` | `T` |
| `Range.End` | `T` |
| `Range.Contains` | `func (r Range[T]) Contains(value T) bool` |
| `Range.IsValid` | `func (r Range[T]) IsValid() bool` |
| `Range.IsEmpty` | `func (r Range[T]) IsEmpty() bool` |
| `Range.IsNotEmpty` | `func (r Range[T]) IsNotEmpty() bool` |
| `Range.Overlaps` | `func (r Range[T]) Overlaps(other monad.Range[T]) bool` |
| `Range.Intersection` | `func (r Range[T]) Intersection(other monad.Range[T]) (monad.Range[T], bool)` |

Surface count: 2 top-level exported symbols, 2 exported fields, 6 exported
methods, no exported constants, and no exported variables.

Behavior contract:

- `Range[T]` is constrained to `cmp.Ordered`
- ranges are inclusive at both ends: `[Start, End]`
- the exported fields have no JSON tags; default Go JSON encoding uses `Start`
  and `End`
- `NewRange(start, end)` stores the two bounds exactly as provided and does not
  reorder them
- `IsValid()` and `IsNotEmpty()` return `Start <= End`
- `IsEmpty()` returns `Start > End`
- `Contains(value)` checks `Start <= value && value <= End`, so both endpoints
  are included
- `Overlaps(other)` checks `Start <= other.End && other.Start <= End`; adjacent
  ranges that share one endpoint overlap
- `Intersection(other)` first calls `Overlaps(other)`
- when there is no overlap, `Intersection` returns `Range[T]{}, false`; the
  returned zero range carries no business meaning
- when there is overlap, `Intersection` returns the inclusive range from
  `max(Start, other.Start)` to `min(End, other.End)` and `true`

### `strx`

Struct-tag and compact key/value string parsing.

Public surface:

| API | Signature / value |
| --- | --- |
| `DefaultKey` | untyped string constant `"__default"` |
| `BareValueMode` | `type BareValueMode int` |
| `BareAsValue` | `strx.BareValueMode = 0` |
| `BareAsKey` | `strx.BareValueMode = 1` |
| `ParseOption` | option type accepted by `ParseTag` |
| `ParseTag` | `strx.ParseTag(input string, opts ...strx.ParseOption) map[string]string` |
| `WithPairDelimiter` | `strx.WithPairDelimiter(delimiter rune) strx.ParseOption` |
| `WithPairDelimiterFunc` | `strx.WithPairDelimiterFunc(fn func(rune) bool) strx.ParseOption` |
| `WithSpacePairDelimiter` | `strx.WithSpacePairDelimiter() strx.ParseOption` |
| `WithValueDelimiter` | `strx.WithValueDelimiter(delimiter rune) strx.ParseOption` |
| `WithBareValueMode` | `strx.WithBareValueMode(mode strx.BareValueMode) strx.ParseOption` |

Surface count: 11 top-level exported symbols, no exported fields, no exported
methods, and no exported variables.

Behavior contract:

- default parsing uses comma-separated pairs, `=` as the key/value separator,
  and `BareAsValue`
- `ParseTag` always returns a non-nil map
- pair tokens are trimmed after splitting; empty tokens are skipped
- key/value pairs are split only at the first value delimiter
- explicit duplicate keys overwrite earlier values
- keys and values themselves are not trimmed after splitting at the value
  delimiter; whitespace inside the pair remains part of the key or value
- empty keys and empty values are accepted (`=value`, `key=`, and `=` all parse)
- special characters and Unicode are preserved as raw strings
- in `BareAsValue`, the first bare token is stored under `DefaultKey`; later
  bare tokens are ignored and logged as warnings
- in `BareAsKey`, every bare token becomes a key with an empty string value;
  duplicate bare keys collapse through normal map overwrite behavior
- `WithPairDelimiter(delimiter)` replaces the pair separator with a single-rune
  equality check
- `WithPairDelimiterFunc(fn)` replaces the pair separator with `fn`
- `WithSpacePairDelimiter()` uses `unicode.IsSpace`, so spaces, tabs, and
  newlines separate pairs
- `WithValueDelimiter(delimiter)` changes the key/value separator from `=`
- options are applied in order; later options can override earlier separator or
  bare-value settings
- option callbacks are called directly; a nil `ParseOption` or nil delimiter
  function will panic when reached

### `dbx`

Database-adjacent helpers.

The public `dbx` surface has 3 exported functions. It has no exported types,
no exported constants, no exported variables, no exported fields, and no
exported methods.

| API | Signature | Purpose |
| --- | --- |
| `ColumnWithAlias` | `dbx.ColumnWithAlias(column string, alias ...string) string` | Prefixes `column` with the first non-empty alias argument. |
| `IsDuplicateKeyError` | `dbx.IsDuplicateKeyError(err error) bool` | Detects duplicate-key errors through driver codes and message fallback. |
| `IsForeignKeyError` | `dbx.IsForeignKeyError(err error) bool` | Detects foreign-key errors through driver codes and message fallback. |

`ColumnWithAlias("name", "u")` returns `u.name`; without an alias, or with an
empty first alias, it returns `name`. Only the first alias argument is used:
`ColumnWithAlias("email", "user", "profile")` returns `user.email`. The helper
does not validate, quote, or escape identifiers; `ColumnWithAlias("", "t")`
returns `t.`.

`IsDuplicateKeyError(nil)` returns false. It first checks wrapped PostgreSQL
`pgdriver.Error` code `23505`, then wrapped MySQL `*mysql.MySQLError` numbers
`1062` and `1169`. If those typed checks do not match, it lowercases
`err.Error()` and searches compatibility patterns for PostgreSQL, MySQL,
SQLite, SQL Server, and Oracle, including `duplicate key`, `unique violation`,
`duplicate entry`, `unique constraint failed`, `violation of primary key
constraint`, `violation of unique key constraint`, `cannot insert duplicate
key`, `ora-00001`, and messages that contain both `unique constraint` and
`violated`.

`IsForeignKeyError(nil)` returns false. It first checks wrapped PostgreSQL
`pgdriver.Error` code `23503`, then wrapped MySQL `*mysql.MySQLError` numbers
`1451` and `1452`. Its message fallback covers PostgreSQL, MySQL, SQLite, SQL
Server, and Oracle patterns such as `violates foreign key constraint`, `foreign
key violation`, `a foreign key constraint fails`, `cannot add or update a child
row`, `cannot delete or update a parent row`, `foreign key constraint failed`,
`sqlite_constraint_foreignkey`, `foreign key mismatch`, `conflicted with the
foreign key constraint`, `statement conflicted with the foreign key`,
`ora-02291`, and `ora-02292`. It also matches Oracle integrity-constraint
messages that contain `violated` plus either `parent key not found` or `child
record found`.

### `httpx`

Fiber request helpers.

The public `httpx` surface has 3 exported functions. It has no exported types,
no exported fields, and no exported methods.

| API | Signature | Purpose |
| --- | --- |
| `IsJSON` | `httpx.IsJSON(ctx fiber.Ctx) bool` | Checks whether Fiber considers the request content type JSON. |
| `IsMultipart` | `httpx.IsMultipart(ctx fiber.Ctx) bool` | Checks whether `Content-Type` starts with `multipart/form-data`. |
| `GetIP` | `httpx.GetIP(ctx fiber.Ctx) string` | Returns the client IP resolved by Fiber. |

The package intentionally exposes helper names rather than raw Fiber calls:
`IsJSON`, `IsMultipart`, and `GetIP`.

`IsJSON` delegates to Fiber's `ctx.Is("json")`, so it accepts standard JSON
content types including charset variants. `IsMultipart` checks for a
`multipart/form-data` prefix with `strings.HasPrefix(...)`, so it accepts
boundary parameters such as `multipart/form-data; boundary=...`.

`GetIP` delegates to `ctx.IP()`. With the framework's app configuration, that
means `vef.app.trusted_proxies` controls whether proxy headers are trusted:
without trusted proxies, a raw client-supplied `X-Forwarded-For` is ignored;
with a trusted proxy configuration, Fiber may honor `X-Forwarded-For` according
to its proxy settings.

### `reflectx`

Reflection and conversion helpers.

The public `reflectx` surface has 78 exported top-level entries and 11 exported
fields. It has no exported methods.

Surface summary: 78 exported top-level entries, 11 exported fields, no exported
methods, fingerprint
`bb62b3bd50f5b54c5af99deb16b7cfb61fa52e69f92e3ab789dc81c744f6d3de`.

| API group | Audited public APIs |
| --- | --- |
| cast aliases | `reflectx.ToString`, `reflectx.ToStringE`, `reflectx.ToInt`, `reflectx.ToIntE`, `reflectx.ToInt8`, `reflectx.ToInt8E`, `reflectx.ToInt16`, `reflectx.ToInt16E`, `reflectx.ToInt32`, `reflectx.ToInt32E`, `reflectx.ToInt64`, `reflectx.ToInt64E`, `reflectx.ToUint`, `reflectx.ToUintE`, `reflectx.ToUint8`, `reflectx.ToUint8E`, `reflectx.ToUint16`, `reflectx.ToUint16E`, `reflectx.ToUint32`, `reflectx.ToUint32E`, `reflectx.ToUint64`, `reflectx.ToUint64E`, `reflectx.ToFloat32`, `reflectx.ToFloat32E`, `reflectx.ToFloat64`, `reflectx.ToFloat64E`, `reflectx.ToBool`, `reflectx.ToBoolE` |
| decimal conversion | `reflectx.ToDecimal`, `reflectx.ToDecimalE` |
| type compatibility and methods | `reflectx.Indirect`, `reflectx.IsPointerToStruct`, `reflectx.IsSimilarType`, `reflectx.IsTypeCompatible`, `reflectx.ConvertValue`, `reflectx.ErrCannotConvertType`, `reflectx.FindMethod`, `reflectx.CollectMethods` |
| string field helpers | `reflectx.IsStringType`, `reflectx.IsStringSliceType`, `reflectx.IsStringMapType`, `reflectx.GetStringValue`, `reflectx.SetStringValue`, `reflectx.GetStringSliceValue`, `reflectx.SetStringSliceValue`, `reflectx.GetStringMapValue`, `reflectx.SetStringMapValue` |
| value helpers | `reflectx.IsEmpty`, `reflectx.IsNotEmpty`, `reflectx.IsNumeric`, `reflectx.IsInteger`, `reflectx.IsSignedInt`, `reflectx.IsUnsignedInt`, `reflectx.IsFloat`, `reflectx.Equal`, `reflectx.Contains` |
| visitor actions, callbacks, and options | `reflectx.VisitAction`, `reflectx.Continue`, `reflectx.Stop`, `reflectx.SkipChildren`, `reflectx.TagConfig`, `reflectx.VisitorConfig`, `reflectx.VisitorOption`, `reflectx.Visitor`, `reflectx.StructVisitor`, `reflectx.FieldVisitor`, `reflectx.MethodVisitor`, `reflectx.TypeVisitor`, `reflectx.StructTypeVisitor`, `reflectx.FieldTypeVisitor`, `reflectx.MethodTypeVisitor`, `reflectx.VisitOf`, `reflectx.Visit`, `reflectx.VisitType`, `reflectx.VisitFor`, `reflectx.WithDisableRecursive`, `reflectx.WithDiveTag`, `reflectx.WithMaxDepth` |

Exported fields:

| Type | Fields |
| --- | --- |
| `reflectx.TagConfig` | `TagConfig.Name`, `TagConfig.Value` |
| `reflectx.VisitorConfig` | `VisitorConfig.Recursive`, `VisitorConfig.DiveTag`, `VisitorConfig.MaxDepth` |
| `reflectx.Visitor` | `Visitor.VisitStruct`, `Visitor.VisitField`, `Visitor.VisitMethod` |
| `reflectx.TypeVisitor` | `TypeVisitor.VisitStructType`, `TypeVisitor.VisitFieldType`, `TypeVisitor.VisitMethodType` |

Cast aliases are pass-through helpers from `github.com/spf13/cast`: `To*E`
variants return conversion errors, while non-E variants return the destination
type's zero value on failure. `reflectx.ToDecimalE` returns `decimal.Zero` with
no error for nil values and nil pointer/interface values; otherwise it
recursively dereferences pointer/interface values and delegates to
`decimal.NewFromAny`. `reflectx.ToDecimal` discards the error and returns
`decimal.Zero` on failure.

`reflectx.Indirect` dereferences one pointer type. `reflectx.IsPointerToStruct`
requires a non-nil `reflect.Type` whose kind is pointer and whose element kind
is struct. `reflectx.IsSimilarType` is true for identical types, or for generic
instantiations with the same `PkgPath` and the same base type name before `[`.
`reflectx.IsTypeCompatible` accepts exact/assignable types, interface targets,
pointer-to-pointer element compatibility, value-to-pointer element assignment,
and pointer-to-value element assignment. `reflectx.ConvertValue` mirrors those
pointer/value conversions, returns zero target values for nil pointer inputs,
allocates new pointers for value-to-pointer and pointer-to-pointer conversions,
and wraps unsupported conversions with `reflectx.ErrCannotConvertType`.
`reflectx.FindMethod` checks the value first, then an addressable pointer copy
for non-pointer values. `reflectx.CollectMethods` dereferences pointers, returns
an empty map for nil or non-struct values, and collects promoted pointer/value
methods by name from the pointer method set.

String field helpers require callers to pass a valid `reflect.Type` or
`reflect.Value`. `reflectx.IsStringType` accepts `string` and `*string` only;
`reflectx.IsStringSliceType` accepts `[]string`; `reflectx.IsStringMapType`
accepts `map[string]string`. Getter helpers return `false` for incompatible
types and nil string pointers/slices/maps. Setter helpers are no-ops for
incompatible values. `reflectx.SetStringValue` replaces a `*string` with a fresh
`*string` pointer instead of mutating an existing pointee.

`reflectx.IsEmpty` treats nil, invalid values, zero scalars, empty
strings/arrays/slices/maps, nil pointer/interface/channel/function values, and
zero structs as empty. It has a special `*string` rule: a non-nil pointer to an
empty string is empty, while other non-nil pointers are not empty.
`reflectx.Equal` compares signed integers only with signed integers, unsigned
integers only with unsigned integers, and floats only with floats; cross-category
numeric comparisons return false. Exact-type comparable values use `==`, and
same-type non-comparable nil-able values are equal only when both are nil.
`reflectx.Contains` supports string substring checks with string elements,
slices/arrays via `reflectx.Equal`, and maps by key lookup with convertible map keys.

Visitor traversal uses depth-first order. `reflectx.Continue = 0`,
`reflectx.Stop = 1`, and `reflectx.SkipChildren = 2`. By default
`VisitorConfig.Recursive` is true, `VisitorConfig.DiveTag` is
`TagConfig{Name: "visit", Value: "dive"}`, and `VisitorConfig.MaxDepth == 0`
means unlimited depth. Anonymous struct fields recurse automatically; named
struct fields recurse only when their tag matches `VisitorConfig.DiveTag`.
`reflectx.WithDisableRecursive` sets `VisitorConfig.Recursive = false`;
`reflectx.WithDiveTag` replaces the tag selector; `reflectx.WithMaxDepth` sets
the depth cap, and traversal stops descending when `depth >= MaxDepth`.
`reflectx.SkipChildren` skips recursion for the current field and
`reflectx.Stop` aborts traversal. `reflectx.Visit` and `reflectx.VisitOf` use
value callbacks, `reflectx.VisitType` and `reflectx.VisitFor` use type-only
callbacks, and invalid, nil pointer, and non-struct inputs are ignored. Field
callbacks receive `StructField.Index` rewritten to the absolute index path, and
method callbacks are visited through the pointer method set.

### `validator`

Framework validation entry points.

| API | Purpose |
| --- | --- |
| `Validate(value)` | validates a value and returns the first framework validation error |
| `RegisterValidationRules(rules...)` | adds custom `ValidationRule` entries |
| `RegisterTypeFunc(fn, types...)` | registers custom type extraction for application-specific wrappers |
| `CustomTypeFunc` | callback type accepted by `RegisterTypeFunc` |
| `ValidationRule` | custom rule definition with tag, messages, validator callback, parameter parser, and null-call flag |

### `logx`

Logging contracts used by the framework.

The public `logx` surface has 7 exported top-level entries and 15 exported
methods. It has no exported fields.

Surface summary: 7 exported top-level entries, 15 exported methods, no exported
fields, fingerprint
`4ff9c19b53d9985911e2985c2763802337d6a69f783962bb73b9ee7424481eaf`.

| API | Contract | Purpose |
| --- | --- |
| `logx.Level` | `type Level int8` | Logging priority; higher levels are more important. |
| `logx.LevelDebug = 1` | debug level constant | Voluminous logs, usually disabled in production. |
| `logx.LevelInfo = 2` | info level constant | Default logging priority. |
| `logx.LevelWarn = 3` | warning level constant | More important than info but not necessarily human-reviewed one by one. |
| `logx.LevelError = 4` | error level constant | High-priority logs for unexpected application behavior. |
| `logx.LevelPanic = 5` | panic level constant | Logs a message and then panics. |
| `logx.Logger` | logging interface | Contract implemented by framework-provided and custom loggers. |

`Level.String() string` returns `debug`, `info`, `warn`, `error`, or `panic`.
Unknown level values, including the zero value, return `unknown`.

| Method | Contract |
| --- | --- |
| `Logger.Named(name string) logx.Logger` | returns a child logger with the given namespace |
| `Logger.WithCallerSkip(skip int) logx.Logger` | returns a logger that adjusts caller stack-frame reporting |
| `Logger.Enabled(level logx.Level) bool` | reports whether the given level is enabled |
| `Logger.Sync()` | flushes buffered log entries; the interface does not return an error |
| `Logger.Debug(message string)` | logs a debug message |
| `Logger.Debugf(template string, args ...any)` | logs a formatted debug message |
| `Logger.Info(message string)` | logs an info message |
| `Logger.Infof(template string, args ...any)` | logs a formatted info message |
| `Logger.Warn(message string)` | logs a warning message |
| `Logger.Warnf(template string, args ...any)` | logs a formatted warning message |
| `Logger.Error(message string)` | logs an error message |
| `Logger.Errorf(template string, args ...any)` | logs a formatted error message |
| `Logger.Panic(message string)` | logs a panic message and then panics |
| `Logger.Panicf(template string, args ...any)` | logs a formatted panic message and then panics |

Application integration code can call `vef.NamedLogger(name string) logx.Logger`
when it needs the framework logger outside dependency injection. That helper is
exported from the root `vef` package; it is not an additional top-level symbol
in `logx`.

### `version`

Framework version constant.

The public `version` surface has 1 exported constant. It has no exported
functions, no exported types, no exported fields, and no exported methods.

| API | Contract | Purpose |
| --- | --- |
| `VEFVersion` | `version.VEFVersion` is an untyped string constant currently equal to `"v0.28.0"`. | Current framework version string. |

The source comment describes the value as the current VEF Framework version in
semver format. The published value includes the leading `v` prefix.

## Practical Usage

The `ptr` package is commonly used in:

- **Search structs**: Optional filter fields use `*bool`, `*string`, etc.
- **Model fields**: Nullable database columns use pointer types
- **Params structs**: Optional update fields

```go
type UserSearch struct {
    api.P
    IsActive *bool   `json:"isActive" search:"eq,column=is_active"`
    DeptID   *string `json:"deptId" search:"eq,column=department_id"`
}

// Setting optional search filters
search := UserSearch{
    IsActive: ptr.Of(true),
    DeptID:   ptr.Of("dept-123"),
}
```
