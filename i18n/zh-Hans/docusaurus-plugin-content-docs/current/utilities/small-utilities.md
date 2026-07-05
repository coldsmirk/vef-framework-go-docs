---
sidebar_position: 8
---

# 小工具

本页文档记录了提供常用辅助函数的轻量级工具包。

## `ptr` — 指针助手

用于安全指针操作的泛型工具函数。

`ptr` 的 public surface 包含 6 个 exported functions。它没有 exported types、
没有 exported fields，也没有 exported methods。

| API | Signature | 作用 |
| --- | --- | --- |
| `Of` | `ptr.Of[T comparable](v T) *T` | 返回 `v` 的指针；当 `v` 是 zero value 时返回 `nil`。 |
| `Zero` | `ptr.Zero[T any]() T` | 返回 `T` 的 zero value。 |
| `Value` | `ptr.Value[T any](p *T, fallbacks ...*T) T` | 先解引用 `p`，再取第一个非 nil fallback，最后返回 zero。 |
| `ValueOrElse` | `ptr.ValueOrElse[T any](p *T, fn func() T) T` | 解引用 `p`；只有 `p` 为 nil 时才调用 `fn`。 |
| `Equal` | `ptr.Equal[T comparable](a *T, b *T) bool` | 按指向的值比较两个指针，并处理 nil。 |
| `Coalesce` | `ptr.Coalesce[T any](ptrs ...*T) *T` | 返回第一个非 nil 指针本身。 |

`ptr.Of` 要求 `T comparable`，因为它会把 `v` 与 zero value 比较。对于 `""`、
`0`、`false` 这样的 zero values，它返回 `nil`；对非零值，它返回指向 `v`
副本的指针。

`ptr.Value` 会从左到右检查 fallbacks。非 nil 的 primary pointer 优先于所有
fallbacks。指向 zero value 的非 nil fallback 仍然会被使用。当所有指针都是
nil 时，它返回 `ptr.Zero[T]()`。

`ptr.ValueOrElse` 是 lazy 的：`p` 非 nil 时不会调用 `fn`。如果 `p` 为 nil，
`fn` 必须非 nil，因为 helper 会直接调用它。

`ptr.Equal` 在两个指针都为 nil 时返回 true，只有一个为 nil 时返回 false，
否则比较 `*a == *b`。

`ptr.Coalesce` 返回准确的第一个非 nil 指针，不会复制值。没有参数或参数全是
nil 时返回 nil。

示例：

```go
import "github.com/coldsmirk/vef-framework-go/ptr"

p := ptr.Of("hello")   // *string → "hello"
p = ptr.Of("")          // *string → nil（零值）
p = ptr.Of(42)          // *int → 42
p = ptr.Of(0)           // *int → nil

s := ptr.Value(p)                        // 解引用，nil 时返回零值
s = ptr.Value(p, fallback1, fallback2)   // 依次尝试每个回退

s := ptr.ValueOrElse(p, func() string {
    return computeDefault()
})

z := ptr.Zero[string]()  // ""
ptr.Equal(a, b)  // 都为 nil → true，一个 nil → false，否则比较值
result := ptr.Coalesce(p1, p2, p3)  // 第一个非 nil，或 nil
```

## 其他公开工具包

这些包暴露的是小而集中的 helper，不需要各自拆成完整 feature 页。

### `page`

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

Surface count：6 个 top-level exported symbols、6 个 exported fields、5 个
exported methods，没有 exported variables。

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

### `sortx`

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

Surface count：9 个 top-level exported symbols、3 个 exported fields、7 个
exported methods，没有 exported functions。

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

### `monad`

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

Surface count：2 个 top-level exported symbols、2 个 exported fields、6 个
exported methods，没有 exported constants，也没有 exported variables。

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

### `strx`

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

Surface count：11 个 top-level exported symbols，没有 exported fields、
没有 exported methods，也没有 exported variables。

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

### `dbx`

数据库相关小工具。

`dbx` 的 public surface 包含 3 个 exported functions。它没有 exported types、
没有 exported constants、没有 exported variables、没有 exported fields，也没有
exported methods。

| API | Signature | 作用 |
| --- | --- |
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

### `httpx`

Fiber 请求 helper。

`httpx` 的 public surface 包含 3 个 exported functions。它没有 exported types、
没有 exported fields，也没有 exported methods。

| API | Signature | 作用 |
| --- | --- |
| `IsJSON` | `httpx.IsJSON(ctx fiber.Ctx) bool` | 判断 Fiber 是否认为请求 content type 是 JSON。 |
| `IsMultipart` | `httpx.IsMultipart(ctx fiber.Ctx) bool` | 判断 `Content-Type` 是否以 `multipart/form-data` 开头。 |
| `GetIP` | `httpx.GetIP(ctx fiber.Ctx) string` | 返回 Fiber 解析出的客户端 IP。 |

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

### `reflectx`

反射和转换 helper。

`reflectx` 的 public surface 包含 82 个 exported top-level entries 和 12 个
exported fields。它没有 exported methods。

Surface summary: 82 exported top-level entries, 12 exported fields, no exported
methods, fingerprint
`0c7f22e87dd56b6a1e33be78c9bb398abee1b4b776888a4a3a64d94167751bba`。

| API 组 | Audited public APIs |
| --- | --- |
| cast aliases | `reflectx.ToString`, `reflectx.ToStringE`, `reflectx.ToInt`, `reflectx.ToIntE`, `reflectx.ToInt8`, `reflectx.ToInt8E`, `reflectx.ToInt16`, `reflectx.ToInt16E`, `reflectx.ToInt32`, `reflectx.ToInt32E`, `reflectx.ToInt64`, `reflectx.ToInt64E`, `reflectx.ToUint`, `reflectx.ToUintE`, `reflectx.ToUint8`, `reflectx.ToUint8E`, `reflectx.ToUint16`, `reflectx.ToUint16E`, `reflectx.ToUint32`, `reflectx.ToUint32E`, `reflectx.ToUint64`, `reflectx.ToUint64E`, `reflectx.ToFloat32`, `reflectx.ToFloat32E`, `reflectx.ToFloat64`, `reflectx.ToFloat64E`, `reflectx.ToBool`, `reflectx.ToBoolE` |
| decimal conversion | `reflectx.ToDecimal`, `reflectx.ToDecimalE` |
| type compatibility and methods | `reflectx.Indirect`, `reflectx.IsPointerToStruct`, `reflectx.IsSimilarType`, `reflectx.IsTypeCompatible`, `reflectx.ConvertValue`, `reflectx.ErrCannotConvertType`, `reflectx.FindMethod`, `reflectx.CollectMethods` |
| string field helpers | `reflectx.IsStringType`, `reflectx.IsStringSliceType`, `reflectx.IsStringMapType`, `reflectx.GetStringValue`, `reflectx.SetStringValue`, `reflectx.GetStringSliceValue`, `reflectx.SetStringSliceValue`, `reflectx.GetStringMapValue`, `reflectx.SetStringMapValue` |
| value helpers | `reflectx.IsEmpty`, `reflectx.IsNotEmpty`, `reflectx.IsNumeric`, `reflectx.IsInteger`, `reflectx.IsSignedInt`, `reflectx.IsUnsignedInt`, `reflectx.IsFloat`, `reflectx.Equal`, `reflectx.Contains` |
| visitor actions, callbacks, traversal modes, and options | `reflectx.VisitAction`, `reflectx.Continue`, `reflectx.Stop`, `reflectx.SkipChildren`, `reflectx.TraversalMode`, `reflectx.DepthFirst`, `reflectx.BreadthFirst`, `reflectx.TagConfig`, `reflectx.VisitorConfig`, `reflectx.VisitorOption`, `reflectx.Visitor`, `reflectx.StructVisitor`, `reflectx.FieldVisitor`, `reflectx.MethodVisitor`, `reflectx.TypeVisitor`, `reflectx.StructTypeVisitor`, `reflectx.FieldTypeVisitor`, `reflectx.MethodTypeVisitor`, `reflectx.VisitOf`, `reflectx.Visit`, `reflectx.VisitType`, `reflectx.VisitFor`, `reflectx.WithTraversalMode`, `reflectx.WithDisableRecursive`, `reflectx.WithDiveTag`, `reflectx.WithMaxDepth` |

Exported fields：

| Type | Fields |
| --- | --- |
| `reflectx.TagConfig` | `TagConfig.Name`, `TagConfig.Value` |
| `reflectx.VisitorConfig` | `VisitorConfig.TraversalMode`, `VisitorConfig.Recursive`, `VisitorConfig.DiveTag`, `VisitorConfig.MaxDepth` |
| `reflectx.Visitor` | `Visitor.VisitStruct`, `Visitor.VisitField`, `Visitor.VisitMethod` |
| `reflectx.TypeVisitor` | `TypeVisitor.VisitStructType`, `TypeVisitor.VisitFieldType`, `TypeVisitor.VisitMethodType` |

Cast aliases 是来自 `github.com/spf13/cast` 的 pass-through helpers：`To*E`
variants 会返回 conversion errors，non-E variants 在失败时返回目标类型的
zero value。`reflectx.ToDecimalE` 对 nil values 和 nil pointer/interface
values 返回 `decimal.Zero` 且没有 error；其他 pointer/interface values 会递归
dereference，然后委托 `decimal.NewFromAny`。`reflectx.ToDecimal` 会丢弃 error，
失败时返回 `decimal.Zero`。

`reflectx.Indirect` 只 dereference 一层 pointer type。
`reflectx.IsPointerToStruct` 要求传入非 nil `reflect.Type`，其 kind 是 pointer，
且 element kind 是 struct。`reflectx.IsSimilarType` 对相同 type 返回 true；
对 generic instantiations，只有 `PkgPath` 相同且 `[` 之前的 base type name
相同时才返回 true。`reflectx.IsTypeCompatible` 接受 exact/assignable types、
interface targets、pointer-to-pointer element compatibility、value-to-pointer
element assignment 和 pointer-to-value element assignment。
`reflectx.ConvertValue` 对这些 pointer/value conversions 做对应转换，对 nil
pointer inputs 返回 zero target values，对 value-to-pointer 和
pointer-to-pointer conversions 分配新 pointer，并用 `reflectx.ErrCannotConvertType`
包装不支持的转换。`reflectx.FindMethod` 先查 value，再对 non-pointer value
使用 addressable pointer copy。`reflectx.CollectMethods` 会 dereference pointers，
nil 或 non-struct values 返回 empty map，并从 pointer method set 按名称收集
promoted pointer/value methods。

String field helpers 要求调用方传入有效的 `reflect.Type` 或 `reflect.Value`。
`reflectx.IsStringType` 只接受 `string` 和 `*string`；`reflectx.IsStringSliceType`
只接受 `[]string`；`reflectx.IsStringMapType` 只接受 `map[string]string`。
getter helpers 对 incompatible types 和 nil string pointers/slices/maps 返回
`false`。setter helpers 对 incompatible values 是 no-op。
`reflectx.SetStringValue` 会把 `*string` 替换成 fresh `*string` pointer，而不是
修改已有 pointee。

`reflectx.IsEmpty` 把 nil、invalid values、zero scalars、empty
strings/arrays/slices/maps、nil pointer/interface/channel/function values 和
zero structs 视为空。它对 `*string` 有特殊规则：指向 empty string 的非 nil
pointer 也是 empty，但其他非 nil pointers 不算 empty。`reflectx.Equal` 只在
signed integers、unsigned integers、floats 各自 category 内做 numeric 比较；
cross-category numeric comparisons 返回 false。exact-type comparable values
使用 `==`；same-type non-comparable nil-able values 只有在两边都 nil 时才相等。
`reflectx.Contains` 支持 string substring checks（element 必须是 string）、
slices/arrays 通过 `reflectx.Equal` 比较，以及 maps 通过 key lookup，且支持
convertible map keys。

Visitor traversal 默认使用 depth-first order。`reflectx.DepthFirst = 0`、
`reflectx.BreadthFirst = 1`、`reflectx.Continue = 0`、
`reflectx.Stop = 1`、`reflectx.SkipChildren = 2`。默认
`VisitorConfig.TraversalMode` 为 `DepthFirst`，`VisitorConfig.Recursive` 为 true，`VisitorConfig.DiveTag` 是
`TagConfig{Name: "visit", Value: "dive"}`，`VisitorConfig.MaxDepth == 0` 表示
unlimited depth。Anonymous struct fields 会自动 recurse；named struct fields
只有在 tag 匹配 `VisitorConfig.DiveTag` 时才 recurse。
`reflectx.WithTraversalMode` 选择 `DepthFirst` 或 `BreadthFirst`；
`reflectx.WithDisableRecursive` 设置 `VisitorConfig.Recursive = false`；
`reflectx.WithDiveTag` 替换 tag selector；`reflectx.WithMaxDepth` 设置 depth
cap，且 traversal 在 `depth >= MaxDepth` 时停止向下。`reflectx.SkipChildren`
跳过当前 field 的 recursion，`reflectx.Stop` 中止 traversal。`reflectx.Visit`
和 `reflectx.VisitOf` 使用 value callbacks，`reflectx.VisitType` 和
`reflectx.VisitFor` 使用 type-only callbacks；invalid、nil pointer 和
non-struct inputs 会被忽略。Field callbacks 收到的 `StructField.Index` 会被重写为
absolute index path，method callbacks 会通过 pointer method set 访问。

### `validator`

框架验证入口。

| API | 作用 |
| --- | --- |
| `Validate(value)` | 验证值并返回第一个框架验证错误 |
| `RegisterValidationRules(rules...)` | 添加自定义 `ValidationRule` |
| `RegisterTypeFunc(fn, types...)` | 为应用自定义 wrapper 注册类型提取函数 |
| `CustomTypeFunc` | `RegisterTypeFunc` 接收的回调类型 |
| `ValidationRule` | 自定义规则定义，包含 tag、消息、validator 回调、参数解析和 null-call 标记 |

### `logx`

框架使用的日志契约。

`logx` 的 public surface 包含 8 个 exported top-level entries 和 16 个
exported methods。它没有 exported fields。

Surface summary: 8 exported top-level entries, 16 exported methods, no exported fields, fingerprint
`2512ac2e4b928560900dfb481cda5645aeda3afbd4229fafa3a42dedcf19d4a3`。

| API | Contract | 作用 |
| --- | --- |
| `logx.Level` | `type Level int8` | 日志优先级；数值越高优先级越高。 |
| `logx.LevelDebug = 1` | debug level constant | 通常较多，生产环境一般关闭。 |
| `logx.LevelInfo = 2` | info level constant | 默认日志优先级。 |
| `logx.LevelWarn = 3` | warning level constant | 比 info 更重要，但不一定需要逐条人工处理。 |
| `logx.LevelError = 4` | error level constant | 表示非预期应用行为的高优先级日志。 |
| `logx.LevelPanic = 5` | panic level constant | 记录消息后触发 panic。 |
| `logx.Logger` | logging interface | 由框架提供的 logger 和自定义 logger 实现的契约。 |
| `logx.LoggerConfigurable[T]` | generic interface | 由 immutable component 实现，通过 `WithLogger` 返回配置了 logger 的副本。 |

`Level.String() string` 返回 `debug`、`info`、`warn`、`error` 或 `panic`。
未知 level 值，包括 zero value，返回 `unknown`。

| Method | Contract |
| --- | --- |
| `Logger.Named(name string) logx.Logger` | 返回带 namespace 的 child logger |
| `Logger.WithCallerSkip(skip int) logx.Logger` | 返回调整 caller stack-frame reporting 的 logger |
| `Logger.Enabled(level logx.Level) bool` | 报告指定 level 是否启用 |
| `Logger.Sync()` | flush buffered log entries；接口不返回 error |
| `Logger.Debug(message string)` | 记录 debug message |
| `Logger.Debugf(template string, args ...any)` | 记录 formatted debug message |
| `Logger.Info(message string)` | 记录 info message |
| `Logger.Infof(template string, args ...any)` | 记录 formatted info message |
| `Logger.Warn(message string)` | 记录 warning message |
| `Logger.Warnf(template string, args ...any)` | 记录 formatted warning message |
| `Logger.Error(message string)` | 记录 error message |
| `Logger.Errorf(template string, args ...any)` | 记录 formatted error message |
| `Logger.Panic(message string)` | 记录 panic message 后触发 panic |
| `Logger.Panicf(template string, args ...any)` | 记录 formatted panic message 后触发 panic |
| `LoggerConfigurable[T].WithLogger(logger logx.Logger) T` | 返回配置了指定 logger 的 component 副本 |

当集成代码需要在 dependency injection 之外拿框架 logger 时，可以调用
`vef.NamedLogger(name string) logx.Logger`。这个 helper 从 root `vef`
package 导出，不是 `logx` 里的额外 top-level symbol。

### `version`

框架版本常量。

`version` 的 public surface 包含 1 个 exported constant。它没有 exported
functions、没有 exported types、没有 exported fields，也没有 exported methods。

| API | Contract | 作用 |
| --- | --- |
| `VEFVersion` | `version.VEFVersion` 是 untyped string constant，当前值为 `"v0.35.0"`。 | 当前框架版本字符串。 |

源码注释把该值描述为当前 VEF Framework version 的 semver format。公开值包含
前导 `v` prefix。

## 实际用法

`ptr` 包常用于：

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
    IsActive: ptr.Of(true),
    DeptID:   ptr.Of("dept-123"),
}
```
