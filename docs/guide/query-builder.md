---
sidebar_position: 6
---

# Query Builder

VEF query building is centered around typed search structs, `search` tags, and CRUD find options. The goal is to keep query rules close to the fields they belong to instead of scattering stringly typed conditions across handlers.

## Search Struct Model

The usual shape is:

```go
type UserSearch struct {
	api.P

	ID       string `json:"id" search:"eq"`
	Keyword  string `json:"keyword" search:"contains,column=username|email"`
	IsActive *bool  `json:"isActive" search:"eq,column=is_active"`
}
```

The `search` tag describes how a field becomes one or more SQL conditions.

## Default Behavior Without A `search` Tag

If a field has no `search` tag at all:

- the framework still includes it in the parsed search schema
- the default operator is `eq`
- the default column name is the snake_case form of the field name

That means this field:

```go
Age int
```

behaves like:

```go
Age int `search:"eq,column=age"`
```

## Search Tag Grammar

The `search` tag supports these patterns:

| Pattern | Meaning |
| --- | --- |
| `search:"eq"` | operator only |
| `search:"contains,column=username\|email"` | operator plus explicit target columns |
| `search:"operator=gte,column=price"` | fully explicit key/value form |
| `search:"operator=in,params=delimiter:\|,type:int"` | operator with extra params |
| `search:"dive"` | recurse into nested struct fields |
| `search:"-"` | ignore this field completely |

Supported tag attributes:

| Attribute | Meaning |
| --- | --- |
| default value or `operator` | query operator |
| `column` | one or more target columns, separated by \| |
| `alias` | table alias used when qualifying columns |
| `params` | extra operator parameters |
| `dive` | recurse into nested struct fields |

## Supported Operators

The framework currently supports all of these operators:

### Comparison operators

| Operator | Meaning |
| --- | --- |
| `eq` | equals |
| `neq` | not equals |
| `gt` | greater than |
| `gte` | greater than or equal |
| `lt` | less than |
| `lte` | less than or equal |

### Range operators

| Operator | Meaning |
| --- | --- |
| `between` | inclusive range |
| `notBetween` | outside range |

### Set operators

| Operator | Meaning |
| --- | --- |
| `in` | value is in a set |
| `notIn` | value is not in a set |

### Null operators

| Operator | Meaning |
| --- | --- |
| `isNull` | applies `IS NULL` |
| `isNotNull` | applies `IS NOT NULL` |

### String matching operators

| Operator | Meaning |
| --- | --- |
| `contains` | contains substring |
| `notContains` | does not contain substring |
| `startsWith` | starts with prefix |
| `notStartsWith` | does not start with prefix |
| `endsWith` | ends with suffix |
| `notEndsWith` | does not end with suffix |

### Case-insensitive string operators

| Operator | Meaning |
| --- | --- |
| `iContains` | case-insensitive contains |
| `iNotContains` | case-insensitive not contains |
| `iStartsWith` | case-insensitive starts with |
| `iNotStartsWith` | case-insensitive not starts with |
| `iEndsWith` | case-insensitive ends with |
| `iNotEndsWith` | case-insensitive not ends with |

## Multi-Column Search

One field can target multiple columns by separating column names with `|`.

Example:

```go
Keyword string `search:"contains,column=username|email|mobile"`
```

This is useful for keyword search against multiple text fields.

## Nested Search With `dive`

`dive` is not a query operator. It is a parser directive telling the framework to recurse into nested structs.

Example:

```go
type UserSearch struct {
	Name string `search:"column=user_name,operator=contains"`
}

type OrderSearch struct {
	api.P

	User UserSearch `search:"dive"`
}
```

## Aliases

Use `alias` when the query should qualify columns with a table alias:

```go
Name string `search:"alias=u,column=name,operator=contains"`
```

This is especially useful for joined queries.

## Operator Parameters

Some operators support extra parameters through the `params=...` section.

Currently relevant parameter keys:

| Param key | Meaning |
| --- | --- |
| `delimiter` | custom delimiter for parsing string-based sets or ranges |
| `type` | explicit parsing type such as `int`, `str`, `bool`, `dec`, `date`, `datetime`, or `time` |

## `between` Input Forms

`between` and `notBetween` support multiple input shapes:

| Input shape | Example |
| --- | --- |
| `monad.Range[T]` style struct | `monad.Range[int]{Start: 1, End: 10}` |
| two-item slice | `[]int{1, 10}` |
| delimited string | `"1,10"` |

For string input, parsing can be controlled through `params`.

Examples:

```go
Price string `search:"operator=between,column=price,params=type:int,delimiter:,"`
DateRange string `search:"operator=between,column=created_at,params=type:date,delimiter:|"`
```

## `in` / `notIn` Input Forms

Set operators support:

| Input shape | Example |
| --- | --- |
| slice field | `[]string{"a", "b"}` |
| delimited string | `"a,b,c"` |
| delimited string with custom delimiter | `"1\|\2\|3"` + `params=delimiter:\|,type:int` |

## Sorting

Sorting is usually handled through metadata using `crud.Sortable`:

```go
type QueryMeta struct {
	api.M
	crud.Sortable
}
```

`crud.Sortable` shape:

| Field | Meaning |
| --- | --- |
| `Sort []sortx.OrderSpec` | list of sort specifications |

Each `sortx.OrderSpec` can express:

| Property | Meaning |
| --- | --- |
| `Column` | target column |
| `Direction` | ascending or descending |
| `NullsOrder` | null ordering |

CRUD find builders can apply these sort specs automatically.

## Pagination

Paging uses `page.Pageable`:

```go
type QueryMeta struct {
	api.M
	page.Pageable
}
```

`FindPage` normalizes page and size, applies limits, and returns `page.Page[T]`.

Important detail:

- `page.Pageable` is decoded from `meta`
- for REST handlers, `?page=1&size=20` lands in raw `params`; it does not automatically populate typed `page.Pageable`

## Data Permissions

Many read builders automatically apply request-scoped data-permission filtering through the query layer.

That means:

- search tags and custom conditions are not the only filters in play
- data permission may add additional conditions transparently
- if your query must bypass this behavior, the relevant CRUD builder has to disable it explicitly

## Query Escape Hatches

When search tags are not expressive enough, CRUD find builders support these extension points:

| Method | Use for |
| --- | --- |
| `WithCondition(...)` | additional `WHERE` conditions |
| `WithRelation(...)` | relation joins |
| `WithDefaultSort(...)` | fallback sorting |
| `WithQueryApplier(...)` | arbitrary typed query customization |
| `WithSelect(...)` / `WithSelectAs(...)` | explicit select-list shaping |

For tree APIs, these escape hatches can also be targeted at different query parts such as `QueryBase`, `QueryRecursive`, and `QueryRoot`.

## Practical Patterns

### Simple equality and keyword search

```go
type UserSearch struct {
	api.P

	ID      string `json:"id" search:"eq"`
	Keyword string `json:"keyword" search:"contains,column=username|email"`
}
```

### Range and set filtering

```go
type ProductSearch struct {
	api.P

	PriceRange string `json:"priceRange" search:"operator=between,column=price,params=type:int,delimiter:,"`
	Statuses   string `json:"statuses" search:"operator=in,column=status,params=delimiter:|"`
}
```

### Nested search

```go
type UserSearch struct {
	Name string `search:"column=user_name,operator=contains"`
}

type OrderSearch struct {
	api.P

	User UserSearch `search:"dive"`
}
```

## Practical Advice

- use a dedicated search struct per resource
- use `search` tags for normal filtering and keep query rules next to the field definition
- prefer explicit multi-column tags for keyword search instead of hidden custom SQL
- use metadata for sorting and pagination
- reach for `WithQueryApplier(...)` only when tag-based configuration is no longer expressive enough
- keep the query contract visible in the type definition instead of burying it in handler code

## Next Step

Read [Hooks](./hooks) if your queries or mutations also need lifecycle-aware behavior around CRUD operations.
