---
sidebar_position: 10
---

# Small Helpers

This page documents small, focused utility packages that do not need a full
feature page on their own: `page`, `sortx`, `monad`, `strx`, `dbx`, `httpx`,
and `version`, plus a note on the removed `ptr` package and its replacement.

## `ptr` — Removed

:::caution
The `ptr` package has been removed from the framework in favor of Go's
builtin `new` and `github.com/samber/lo`'s pointer helpers. If you are
migrating code that still imports `github.com/coldsmirk/vef-framework-go/ptr`,
use the table below.
:::

| Old `ptr` API | Replacement | Notes |
| --- | --- | --- |
| `ptr.Of(v)` | `new(v)` (builtin) | Go's builtin `new` accepts a value expression and returns a pointer to a copy of it. Unlike the old `ptr.Of`, `new(v)` always returns a non-nil pointer — including for zero values such as `""`, `0`, and `false`. |
| `ptr.Of(v)` (nil-for-zero semantics) | `lo.EmptyableToPtr(v)` | Returns `nil` when `v` is the zero value, matching the old `ptr.Of` behavior exactly. |
| `ptr.Zero[T]()` | `lo.Empty[T]()` | Returns the zero value for `T`. |
| `ptr.Value(p)` | `lo.FromPtr(p)` | Dereferences `p`, or returns zero if `p` is nil. |
| `ptr.Value(p, fallback)` / `ptr.ValueOrElse(p, fn)` | `lo.FromPtrOr(p, fallback)` | Dereferences `p`, or returns the given fallback if `p` is nil. `lo.FromPtrOr` evaluates the fallback eagerly, unlike the old lazy `ptr.ValueOrElse`. |
| `ptr.Coalesce(p1, p2, ...)` | plain `if`/`switch` on nil checks | `lo` has no direct multi-pointer coalesce helper; a short manual check is clearest. |

```go
import "github.com/samber/lo"

p := new(true)                  // *bool → true (builtin new, always non-nil)
p2 := lo.EmptyableToPtr("")     // *string → nil (zero value)
p3 := lo.EmptyableToPtr(42)     // *int → 42

s := lo.FromPtr(p2)             // "" (dereference, or zero if nil)
s = lo.FromPtrOr(p2, "default") // "default"
```

## `page`

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

## `sortx`

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

## `monad`

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

## `strx`

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

## `dbx`

Database-adjacent helpers.

| API | Signature | Purpose |
| --- | --- | --- |
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

## `httpx`

Fiber request helpers.

| API | Signature | Purpose |
| --- | --- | --- |
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

## `version`

Framework version constant.

| API | Contract | Purpose |
| --- | --- | --- |
| `VEFVersion` | `version.VEFVersion` is an untyped string constant. | Current framework version string. |

The source comment describes the value as the current VEF Framework version in
semver format. The published value includes the leading `v` prefix.

## Practical Usage

### Optional fields with pointers

Since the `ptr` package was removed, use the builtin `new` for pointers that
should always be non-nil, or `lo.EmptyableToPtr` when a zero value should
become `nil`:

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
    IsActive: new(true),
    DeptID:   new("dept-123"),
}
```

### Pagination round trip

`page.Pageable` decodes request params, `page.New` wraps the query result:

```go
var pageable page.Pageable
pageable.Normalize() // clamps Page/Size to sane defaults

items, total := queryItems(pageable)
response := page.New(pageable, total, items)
```

### Ordering from request parameters

`sortx.OrderDirection` round-trips through JSON as `"asc"` / `"desc"`, so it
plugs directly into typed search structs alongside `sortx.OrderSpec` for
query-builder ORDER BY generation.

### Detecting constraint violations

```go
if err := db.NewInsert().Model(&user).Exec(ctx); err != nil {
    if dbx.IsDuplicateKeyError(err) {
        return result.Err("username already taken")
    }
    return err
}
```
