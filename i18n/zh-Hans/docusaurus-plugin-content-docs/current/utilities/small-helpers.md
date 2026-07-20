---
sidebar_position: 10
---

# 小助手

本页文档记录了小而集中、不需要单独 feature 页的工具包：`page`、`sortx`、
`monad`、`strx`、`dbx`、`fiberx` 和 `version`。

## `page`

分页请求参数与响应 helper。

公开 surface：

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

行为契约：

- `Pageable.Page` 是 1-based 页码，`Pageable.Size` 是请求页大小
- `Normalize(size...)` 会原地修改 receiver
- `Normalize` 会把 `Page < 1` 重置为 `DefaultPageNumber`
- 当 `Size < 1` 时，`Normalize` 只使用第一个可选 fallback size；没有
  fallback 时使用 `DefaultPageSize`
- 大于 `MaxPageSize` 的值会被截到 `MaxPageSize`
- `Normalize` 不会对负数自定义 fallback 再做一次校正，因此传入的 fallback
  size 应该是正数
- `Offset()` 只是 `(Page - 1) * Size` 计算，调用前应先完成 normalize
- `New` 会把 `pageable.Page`、`pageable.Size` 和 `total` 复制进响应
- `New` 会把 nil `items` 转成非 nil 空 slice；非 nil `items` 按原样使用，
  不会 clone
- `TotalPages()` 在 `Size == 0` 时返回 `0`；否则按 `Total / Size` 向上取整
- `HasNext()` 判断 `Page < TotalPages()`
- `HasPrevious()` 在 `Page > 1` 时返回 true
- 这个 helper 不会校验负数 total、手工构造后的负数 size，或 `Normalize`
  之外不一致的 `Page` 值

## `sortx`

排序表达的基础类型，常用于 query builder 和 API 参数。

公开 surface：

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

行为契约：

- `OrderDirection.String()` 只有在值为 `OrderDesc` 时返回 `DESC`；其他值都会
  渲染为 `ASC`
- `MarshalText()` 会把 `String()` 的结果转成小写，因此正常值会 marshal 成
  `asc` 或 `desc`
- `UnmarshalText(text)` 会 trim 首尾空白，并以大小写不敏感方式接受 `asc` /
  `desc`
- 非法 text 会返回包装了 `ErrInvalidOrderDirection` 的错误
- `MarshalJSON()` 委托 `MarshalText()`，并输出 JSON string
- `UnmarshalJSON(data)` 要求输入是 JSON string；number、boolean 和 `null`
  会在方向解析前返回 `OrderDirection must be a JSON string` 错误
- `NullsDefault.String()` 返回空字符串
- `NullsFirst.String()` 返回 `NULLS FIRST`
- `NullsLast.String()` 返回 `NULLS LAST`
- 其他 `NullsOrder` 值也返回空字符串
- `OrderSpec.IsValid()` 只检查 `Column != ""`；它不会把 column 当 SQL
  identifier 校验，也不会校验 `Direction` / `NullsOrder`

## `monad`

闭区间工具。

公开 surface：

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

行为契约：

- `Range[T]` 约束为 `cmp.Ordered`
- 区间两端都是闭区间：`[Start, End]`
- exported fields 没有 JSON tags；Go 默认 JSON 编码会使用 `Start` 和 `End`
- `NewRange(start, end)` 会按原样保存两个边界，不会重排
- `IsValid()` 和 `IsNotEmpty()` 返回 `Start <= End`
- `IsEmpty()` 返回 `Start > End`
- `Contains(value)` 检查 `Start <= value && value <= End`，因此包含两个端点
- `Overlaps(other)` 检查 `Start <= other.End && other.Start <= End`；共享一个
  端点的相邻区间也算重叠
- `Intersection(other)` 会先调用 `Overlaps(other)`
- 没有重叠时，`Intersection` 返回 `Range[T]{}, false`；返回的 zero range
  没有业务意义
- 有重叠时，`Intersection` 返回从 `max(Start, other.Start)` 到
  `min(End, other.End)` 的闭区间，并返回 `true`

## `strx`

结构体 tag 与紧凑 key/value 字符串解析。

公开 surface：

| API | Signature / value |
| --- | --- |
| `DefaultKey` | untyped string constant `"__default"` |
| `BareValueMode` | `type BareValueMode int` |
| `BareAsValue` | `strx.BareValueMode = 0` |
| `BareAsKey` | `strx.BareValueMode = 1` |
| `ParseOption` | `ParseTag` 接收的 option 类型 |
| `ParseTag` | `strx.ParseTag(input string, opts ...strx.ParseOption) map[string]string` |
| `WithPairDelimiter` | `strx.WithPairDelimiter(delimiter rune) strx.ParseOption` |
| `WithPairDelimiterFunc` | `strx.WithPairDelimiterFunc(fn func(rune) bool) strx.ParseOption` |
| `WithSpacePairDelimiter` | `strx.WithSpacePairDelimiter() strx.ParseOption` |
| `WithValueDelimiter` | `strx.WithValueDelimiter(delimiter rune) strx.ParseOption` |
| `WithBareValueMode` | `strx.WithBareValueMode(mode strx.BareValueMode) strx.ParseOption` |

行为契约：

- 默认解析使用 comma-separated pairs、`=` 作为 key/value separator，并使用
  `BareAsValue`
- `ParseTag` 总是返回非 nil map
- pair token 会在分隔后 trim；空 token 会被跳过
- key/value pair 只在第一个 value delimiter 处分割
- 显式重复 key 会以后出现的值覆盖先出现的值
- 在 value delimiter 处分割后，key 和 value 本身不会再 trim；pair 内部的
  whitespace 会保留在 key 或 value 中
- empty keys 和 empty values 都会被接受（`=value`、`key=` 和 `=` 都能解析）
- 特殊字符和 Unicode 会按 raw strings 保留
- 在 `BareAsValue` 中，第一个 bare token 会写入 `DefaultKey`；后续 bare
  tokens 会被忽略并记录 warnings
- 在 `BareAsKey` 中，每个 bare token 都会成为 value 为空字符串的 key；
  duplicate bare keys 按普通 map overwrite 行为折叠
- `WithPairDelimiter(delimiter)` 会把 pair separator 替换成单个 rune 的相等判断
- `WithPairDelimiterFunc(fn)` 会把 pair separator 替换成 `fn`
- `WithSpacePairDelimiter()` 使用 `unicode.IsSpace`，因此 spaces、tabs 和
  newlines 都会分隔 pairs
- `WithValueDelimiter(delimiter)` 会把 key/value separator 从 `=` 改成指定 rune
- options 按顺序应用；后面的 options 可以覆盖前面的 separator 或 bare-value
  setting
- option callbacks 会被直接调用；nil `ParseOption` 或 nil delimiter function
  在执行到时会 panic

## `dbx`

数据库相关小工具。

| API | Signature | 作用 |
| --- | --- | --- |
| `ColumnWithAlias` | `dbx.ColumnWithAlias(column string, alias ...string) string` | 用第一个非空 alias 参数为 `column` 加前缀。 |
| `IsDuplicateKeyError` | `dbx.IsDuplicateKeyError(err error) bool` | 通过 driver codes 和 message fallback 判断 duplicate-key 错误。 |
| `IsForeignKeyError` | `dbx.IsForeignKeyError(err error) bool` | 通过 driver codes 和 message fallback 判断 foreign-key 错误。 |

`ColumnWithAlias("name", "u")` 返回 `u.name`；没有 alias 或 alias 为空时返回
`name`。只有第一个 alias 参数会被使用：`ColumnWithAlias("email", "user",
"profile")` 返回 `user.email`。这个 helper 不会 validate、quote 或 escape
identifiers；`ColumnWithAlias("", "t")` 会返回 `t.`。

`IsDuplicateKeyError(nil)` 返回 false。它会先检查 wrapped PostgreSQL
`pgdriver.Error` code `23505`，再检查 wrapped MySQL `*mysql.MySQLError`
numbers `1062` 和 `1169`。如果 typed checks 没有命中，它会把 `err.Error()`
转成小写，并按 PostgreSQL、MySQL、SQLite、SQL Server、Oracle 的兼容 message
patterns 兜底，包括 `duplicate key`、`unique violation`、`duplicate entry`、
`unique constraint failed`、`violation of primary key constraint`、
`violation of unique key constraint`、`cannot insert duplicate key`、
`ora-00001`，以及同时包含 `unique constraint` 和 `violated` 的消息。

`IsForeignKeyError(nil)` 返回 false。它会先检查 wrapped PostgreSQL
`pgdriver.Error` code `23503`，再检查 wrapped MySQL `*mysql.MySQLError`
numbers `1451` 和 `1452`。message fallback 覆盖 PostgreSQL、MySQL、SQLite、
SQL Server、Oracle 的 patterns，例如 `violates foreign key constraint`、
`foreign key violation`、`a foreign key constraint fails`、`cannot add or
update a child row`、`cannot delete or update a parent row`、`foreign key
constraint failed`、`sqlite_constraint_foreignkey`、`foreign key mismatch`、
`conflicted with the foreign key constraint`、`statement conflicted with the
foreign key`、`ora-02291` 和 `ora-02292`。它也会匹配包含 `violated`，且包含
`parent key not found` 或 `child record found` 的 Oracle integrity-constraint
消息。

## `fiberx`

Fiber 请求 helper。不要与 `httpx` 混淆——那个名字属于
[出站 HTTP 客户端](./httpx)。

| API | Signature | 作用 |
| --- | --- | --- |
| `IsJSON` | `fiberx.IsJSON(ctx fiber.Ctx) bool` | 判断 Fiber 是否认为请求 content type 是 JSON。 |
| `IsMultipart` | `fiberx.IsMultipart(ctx fiber.Ctx) bool` | 判断 `Content-Type` 是否以 `multipart/form-data` 开头。 |
| `GetIP` | `fiberx.GetIP(ctx fiber.Ctx) string` | 返回 Fiber 解析出的客户端 IP。 |

这个包刻意暴露 helper 名称，而不是要求直接调用 Fiber：`IsJSON`、
`IsMultipart` 和 `GetIP`。

`IsJSON` 委托 Fiber 的 `ctx.Is("json")`，因此会接受标准 JSON content type，
包括带 charset 的形式。`IsMultipart` 检查 `multipart/form-data` 前缀。
它使用 `strings.HasPrefix(...)`，因此会接受
`multipart/form-data; boundary=...` 这样的 boundary 参数。

`GetIP` 委托 `ctx.IP()`。在框架的 app 配置下，这意味着
`vef.app.trusted_proxies` 控制 proxy headers 是否可信：没有 trusted proxies
时，客户端直接伪造的 `X-Forwarded-For` 会被忽略；配置 trusted proxy 后，
Fiber 会按自己的 proxy settings 处理 `X-Forwarded-For`。

## `version`

框架版本常量。

| API | Contract | 作用 |
| --- | --- | --- |
| `VEFVersion` | `version.VEFVersion` 是 untyped string constant。 | 当前框架版本字符串。 |

源码注释把该值描述为当前 VEF Framework version 的 semver format。公开值包含
前导 `v` prefix。

## 实际用法

### 用指针表达可选字段

需要始终非 nil 的指针时使用内置的 `new`；需要 zero-value 变成 `nil` 时
使用 `lo.EmptyableToPtr`：

- **搜索结构体**：可选筛选字段使用 `*bool`、`*string` 等
- **模型字段**：可空数据库列使用指针类型
- **参数结构体**：可选更新字段

```go
type UserSearch struct {
    api.P
    IsActive *bool   `json:"isActive" search:"eq,column=is_active"`
    DeptID   *string `json:"deptId" search:"eq,column=department_id"`
}

// 设置可选搜索筛选条件
search := UserSearch{
    IsActive: new(true),
    DeptID:   new("dept-123"),
}
```

### 分页往返

`page.Pageable` 解码请求参数，`page.New` 包装查询结果：

```go
var pageable page.Pageable
pageable.Normalize() // 把 Page/Size 收敛到合理的默认值

items, total := queryItems(pageable)
response := page.New(pageable, total, items)
```

### 从请求参数构造排序

`sortx.OrderDirection` 通过 JSON 以 `"asc"` / `"desc"` 往返，可以直接用在
typed search 结构体里，配合 `sortx.OrderSpec` 生成 query builder 的
ORDER BY。

### 检测约束冲突

```go
if err := db.NewInsert().Model(&user).Exec(ctx); err != nil {
    if dbx.IsDuplicateKeyError(err) {
        return result.Err("username already taken")
    }
    return err
}
```
