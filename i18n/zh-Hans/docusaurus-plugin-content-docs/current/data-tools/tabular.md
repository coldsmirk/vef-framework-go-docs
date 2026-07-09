---
sidebar_position: 5
---

# Tabular Import & Export

VEF 在 `tabular` 包中提供统一的表格引擎，并通过 `csv` 与 `excel` 两个轻量驱动包暴露具体格式。三者拥有完全对称的工厂函数，所有读写都经由同一个 `RowAdapter` 抽象。这意味着：

- **静态行**（用 `tabular` 标签描述的 Go 结构体）和
- **动态行**（运行期定义列、数据是 `map[string]any`）

共用同一条 importer / exporter 流水线。你只需选择合适的适配器，格式驱动负责剩下的工作。

## 架构

```
tabular/   // schema、列、适配器、formatter / parser、错误
  ├── adapter.go        // RowAdapter, RowReader, RowView, RowWriter, RowBuilder
  ├── schema.go         // Schema, Column
  ├── struct_adapter.go // StructAdapter（结构体 + 框架 validator）
  ├── map_adapter.go    // MapAdapter（map + Required / Validators / RowValidator）
  ├── spec.go           // ColumnSpec, NewSchemaFromSpecs, NewMapAdapterFromSpecs
  ├── resolver.go       // FormatterFn / Formatter / 默认实现的优先级解析
  ├── mapping.go        // Header → schema 列映射
  ├── parse_row.go       // ParseRow, IsEmptyRow
  ├── import_rows.go     // ImportRows（驱动共享的 core）
  ├── typed.go           // TypedImporter[T] / TypedExporter[T] 包装器
  └── errors.go          // 共享错误（ErrRequiredMissing、ErrSchemaMismatch 等）

csv/       // CSV 驱动：NewImporter、NewExporter、NewImporterFor、NewExporterFor、
           //          NewMapImporter、NewMapExporter、NewTyped*For
excel/     // Excel 驱动，与 csv 对称，外加 sheet 等 Excel 专属选项
```

`csv` / `excel` 不持有任何 model 相关的反射或私有错误类型。常用错误统一定义在 `tabular`（`ErrDataMustBeSlice`、`ErrRequiredMissing`、`ErrUnknownColumn`、`ErrSchemaMismatch` 等），仅 `excel.ErrSheetIndexOutOfRange` 是驱动专属的。

## 何时使用哪种用法

| 场景 | 推荐工厂 |
| --- | --- |
| 已有 Go 结构体（例如某个 model）描述了所有列 | `csv.NewImporterFor[T]` / `csv.NewExporterFor[T]` 或 `excel.NewImporterFor[T]` / `excel.NewExporterFor[T]`（以及 `*Typed*` 变体） |
| 列在运行期才确定——多租户表格、用户自定义模板、动态表单 | `csv.NewMapImporter` / `csv.NewMapExporter` 或 `excel.NewMapImporter` / `excel.NewMapExporter`，由 `[]tabular.ColumnSpec` 驱动 |
| 行数据自带特殊来源（channel、自定义业务类型等） | 自行实现 `tabular.RowAdapter`，传入 `csv.NewImporter` / `csv.NewExporter` 或 `excel.NewImporter` / `excel.NewExporter` |

导入返回类型由适配器决定：

- 结构体适配器 → `[]T`
- map 适配器 → `[]map[string]any`

选定适配器之后，再选择具体格式：

| 特性 | CSV | Excel |
| --- | --- | --- |
| 文件格式 | `.csv`（纯文本）| `.xlsx`（二进制）|
| 多工作表 | 不支持 | 支持 |
| 列宽度 | 忽略 | 应用 |
| 分隔符 | 可配置 | 不适用 |
| 注释行 | 支持 | 不适用 |
| 空白修剪 | 可配置 | 不适用 |
| 换行符 | LF 或 CRLF | 不适用 |
| Native typed cell（数字/日期在文件中仍可排序）| 不适用（文本格式）| 默认支持 |
| 依赖 | Go 标准库 | [excelize](https://github.com/xuri/excelize) |

两个包都实现了 `tabular.Importer` / `tabular.Exporter` 接口，因此可以在不更改模型定义的情况下互换使用。

## `tabular` 标签

使用结构体标签定义字段如何映射到列：

```go
type Employee struct {
    orm.FullAuditedModel `tabular:"-"`

    Name       string          `tabular:"姓名,width=20"`
    Email      string          `tabular:"邮箱,width=30"`
    Department string          `tabular:"name=部门,order=2,width=15"`
    JoinDate   timex.Date      `tabular:"入职日期,format=2006-01-02,width=15"`
    Salary     decimal.Decimal `tabular:"薪资,width=12,formatter=money"` // 形如 "#,##0.00" 这种含逗号的 format 无法通过 tag 配置 —— 改注册一个 formatter
    Status     string          `tabular:"状态,default=active,formatter=status"`
}
```

### 标签属性

| 属性 | 类型 | 说明 |
| --- | --- | --- |
| （默认值） | string | 列标题名称 |
| `name` | string | 显式列名（默认值的替代写法）|
| `order` | int | 列显示顺序（从 0 开始，默认：字段声明顺序）|
| `width` | float64 | 列宽度提示（Excel 导出时使用）|
| `default` | string | 导入时空单元格的默认值 |
| `format` | string | 格式化模板（日期格式、数字格式）|
| `formatter` | string | 导出时的自定义格式化器名称 |
| `parser` | string | 导入时的自定义解析器名称 |

### 特殊标签

| 标签 | 含义 |
| --- | --- |
| `tabular:"-"` | 完全忽略该字段 |
| `tabular:"dive"` | 递归进入嵌入结构体字段 |

tag parser 使用逗号分隔的 `key=value` 对。分号不是分隔符；`tabular:"name=ID;order=1"` 会被当作一个 `name` 值。

`dive` 只会递归进入结构体或指向结构体的指针字段；作用在其他 kind 上时，该字段既不会被递归也不会被输出，框架会记录一条警告，而不是悄悄丢弃它。

上面每个属性键和哨兵值都同时以 `tabular` 常量的形式导出，供以编程方式构建或检查标签的调用方使用，而不必手写字符串字面量：`TagTabular`（结构体标签名本身，值为 `"tabular"`）、`IgnoreField`（`"-"` 忽略哨兵值）、`AttrName`、`AttrOrder`、`AttrWidth`、`AttrDefault`、`AttrFormat`、`AttrFormatter`、`AttrParser` 以及 `AttrDive`。

## Schema

`Schema` 在初始化时预解析表格元数据——可以来自结构体类型（`NewSchemaFor[T]` / `NewSchema`），也可以来自动态列描述（`NewSchemaFromSpecs`）：

```go
schema := tabular.NewSchemaFor[Employee]()

columns := schema.Columns()             // []*Column — 所有解析出的列
names := schema.ColumnNames()           // []string{"姓名", "邮箱", ...}
count := schema.ColumnCount()           // 6
col, ok := schema.ColumnByKey("Name")   // 按逻辑 key（结构体字段名）查找
col, ok = schema.ColumnByName("姓名")    // 按 header 名称查找
```

列会按 `order` 属性自动排序。未显式指定 `order` 的字段使用其声明顺序。

`NewSchemaFromSpecs` 会在构造期校验动态 schema：缺 `Key` 返回 `ErrMissingColumnKey`，缺 `Type` 返回 `ErrMissingColumnType`，重复 key 返回 `ErrDuplicateColumnKey`，解析后的 header name 重复返回 `ErrDuplicateHeaderName`。

### `Column`

无论是结构体解析出来的列，还是动态列，都用同一个 `Column` 结构体表示：

| 字段 | 含义 |
| --- | --- |
| `Key` | 逻辑标识：结构体字段名，或动态 schema 中的 map key |
| `Name` | 导出时显示的表头文本，导入时用于匹配表头；默认等于 `Key` |
| `Type` | 解析单元格值所用的 `reflect.Type` |
| `Order` | 列顺序的稳定排序键 |
| `Width` | 列宽提示（Excel）|
| `Default` | 导入时，源单元格为空时使用的默认值 |
| `Format` | 默认 `Formatter`/`ValueParser` 使用的格式模板 |
| `Formatter` / `Parser` | 从导出器 / 导入器注册表按名字查找 |
| `FormatterFn` / `ParserFn` | 直接绑定在列上的 `Formatter` / `ValueParser` 实例（最高优先级）|
| `Required` | 导入时空单元格触发 `ErrRequiredMissing`（动态 schema）|
| `Validators` | 解析后执行的 `[]CellValidator`（动态 schema）|
| `Index` | `StructAdapter` 使用的结构体字段 index path；动态列为 `nil` |

## 接口

### Importer

```go
type Importer interface {
    RegisterParser(name string, parser ValueParser)
    ImportFromFile(filename string) (any, []ImportError, error)
    Import(reader io.Reader) (any, []ImportError, error)
}
```

### Exporter

```go
type Exporter interface {
    RegisterFormatter(name string, formatter Formatter)
    ExportToFile(data any, filename string) error
    Export(data any) (*bytes.Buffer, error)
}
```

### Formatter（导出）

```go
type Formatter interface {
    Format(value any) (string, error)
}

// 便捷适配器
tabular.FormatterFunc(func(value any) (string, error) { ... })
```

### ValueParser（导入）

```go
type ValueParser interface {
    Parse(cellValue string, targetType reflect.Type) (any, error)
}

// 便捷适配器
tabular.ParserFunc(func(cellValue string, targetType reflect.Type) (any, error) { ... })
```

## Header 映射与行导入

两个驱动都通过同一套共享 core 解析 header：

- `BuildHeaderMapping(headerRow, schema, opts)` 按 `Column.Name` 匹配 header 单元格。启用 `MappingOptions.TrimSpace` 时，会先修剪 header name 再匹配。空 header 单元格会跳过；未知 header 单元格会跳过（不会因为多余列失败）。重复的非空 header 是致命错误：`ErrDuplicateHeaderName`。
- Importer 配置为 `WithoutHeader()` 时，driver 会回退到 `DefaultPositionalMapping(schema)`——源文件第 0 列对应 schema 第一列，依此类推。
- `ParseRow(cells, mapping, schema, builder, parsers, rowNumber, opts)` 会先应用 `Column.Default` 再解析单元格；默认值替换后仍为空的 cell 会被跳过；parse 和 `Set` failure 会作为行级 `ImportError` 返回。如果返回了 row error，row builder **不会**提交 partial row。
- `IsEmptyRow(cells, trimSpace)` 判断一行的所有单元格是否都为空（用于自动跳过空行）。
- `tabular.ImportRows(rows, adapter, parsers, opts)` 会通过 `RowAdapter` 解析一个已经 materialized 的 `[][]string` 表格。`ImportRowsOptions` 控制 `SkipRows`、`HasHeader`（是否按 header 映射，还是按位置映射）以及 `TrimSpace`。CSV 和 Excel importer 读取各自格式后，都会委托给这个共享 core。

`BuildHeaderMapping` 和 `DefaultPositionalMapping` 都返回原始的 `map[int]int`（源列 index → schema 列 index）；`ImportRows` 每次导入只调用一次 `NewColumnMapping` 把它包装成 `ColumnMapping`，预先排好源 index，避免 `ParseRow` 每行都重新排序。`ParseRow` 本身接受一个 `ParseRowOptions{TrimSpace bool}`（`ParseRowOptions.TrimSpace`）控制单元格级 trim，它和 header 级的 `MappingOptions` 是两个不同的类型。

## 静态结构体用法

在结构体字段上打 `tabular` 标签，然后用 `csv.NewImporterFor[T]` / `excel.NewImporterFor[T]`（或对应的 exporter）为该结构体创建类型化的 importer/exporter。校验委托给框架的 `validator` 包——继续使用 `validate:"…"` 标签即可，提交每行时会自动执行。

### 导出

```go
exp := csv.NewExporterFor[Employee]()
buf, err := exp.Export(employees) // employees: []Employee 或 []*Employee
// 或者直接写入磁盘：
err = exp.ExportToFile(employees, "employees.csv")
```

Excel：

```go
exp := excel.NewExporterFor[Employee](excel.WithSheetName("Employees"))
buf, err := exp.Export(employees)
```

### 导入

```go
imp := csv.NewImporterFor[Employee]()
result, rowErrors, err := imp.Import(reader)
if err != nil {
    return err // 顶层失败（例如文件损坏）
}
employees := result.([]Employee)
for _, ie := range rowErrors {
    log.Warnf("row %d column %s: %v", ie.Row, ie.Column, ie.Err)
}
```

行级失败（解析错误、结构体校验失败、adapter commit 失败）会聚合到 `[]tabular.ImportError`，**不会中断**整体导入。顶层 `error` 仅在出现致命问题（无法读取文件、Header 行损坏、没有数据行）时返回。

### Typed 包装器

`any` 返回值有时不便使用，两个包都提供了泛型包装器替你做类型断言：

```go
imp := csv.NewTypedImporterFor[Employee]()
employees, rowErrors, err := imp.Import(reader) // employees 直接是 []Employee，无需类型断言

exp := csv.NewTypedExporterFor[Employee]()
buf, err := exp.Export(employees)               // 直接接受 []Employee
```

`TypedImporter[T]` / `TypedExporter[T]` 包裹底层的 `tabular.Importer` / `tabular.Exporter`。如需直接调用 `RegisterParser` / `RegisterFormatter`，可以使用 `TypedImporter.Inner` / `TypedExporter.Inner` 取出内部实例。如果被包装的 importer 返回的行元素类型与 `T` 不匹配，typed 包装器会返回 `ErrTypedRowMismatch`。

`csv.NewTypedImporterFor[T]` / `excel.NewTypedImporterFor[T]`（以及对应的 exporter）是对 `tabular.NewTypedImporter[T](inner)` / `tabular.NewTypedExporter[T](inner)` 的便捷封装。如果你包装的是自己构造的 `Importer`/`Exporter`（例如 `csv.NewImporter(adapter, ...)` 搭配自定义 `RowAdapter`，而不是结构体适配器的 `*For` 变体），直接调用 `tabular` 的构造器即可。

## 动态 Map 用法

动态列允许在运行期构造 schema，无需预先声明结构体。每列由 `tabular.ColumnSpec` 描述：

```go
import (
    "reflect"
    "time"

    "github.com/coldsmirk/vef-framework-go/csv"
    "github.com/coldsmirk/vef-framework-go/excel"
    "github.com/coldsmirk/vef-framework-go/tabular"
)

specs := []tabular.ColumnSpec{
    {Key: "id",       Name: "用户ID", Type: reflect.TypeFor[int](),       Required: true, Order: 1},
    {Key: "name",     Name: "姓名",   Type: reflect.TypeFor[string](),    Required: true, Order: 2},
    {Key: "birthday", Name: "生日",   Type: reflect.TypeFor[time.Time](), Format: "2006-01-02", Order: 3},
    {Key: "active",   Name: "激活",   Type: reflect.TypeFor[bool](),      Default: "false", Order: 4},
}
```

`ColumnSpec` 字段：

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| `Key` | 是 | 逻辑标识，也是读写时使用的 map key，必须唯一 |
| `Type` | 是 | 解析目标类型，使用 `reflect.TypeFor[T]()` |
| `Name` | 否 | Header 名称，默认等于 `Key` |
| `Order` | 否 | 列顺序的稳定排序键 |
| `Width` | 否 | Excel 列宽提示 |
| `Default`（`ColumnSpec.Default`） | 否 | 源单元格为空时的默认值 |
| `Format` | 否 | 默认 formatter / parser 的模板（日期、浮点等）|
| `Formatter` / `Parser` | 否 | 名字，从 importer / exporter 的注册表查找 |
| `FormatterFn` / `ParserFn` | 否 | 直接绑定的 `tabular.Formatter` / `tabular.ValueParser` 实例 |
| `Required` | 否 | 导入时空值会触发 `ErrRequiredMissing` |
| `Validators` | 否 | 解析后执行的 `[]CellValidator` |

`NewSchemaFromSpecs` 会立即校验输入：缺 `Key`、缺 `Type`、`Key` 重复、解析后的 `Name` 重复都会在构造期返回错误（`ErrMissingColumnKey`、`ErrMissingColumnType`、`ErrDuplicateColumnKey`、`ErrDuplicateHeaderName`）。

### 导出

```go
exp, err := excel.NewMapExporter(specs, excel.WithSheetName("Users"))
if err != nil { return err }

buf, err := exp.Export([]map[string]any{
    {"id": 1, "name": "张三", "birthday": time.Now(), "active": true},
    {"id": 2, "name": "李四", "birthday": time.Now(), "active": false},
})
```

CSV 完全对称：

```go
exp, err := csv.NewMapExporter(specs)
buf, err := exp.Export(rows)
```

### 导入

```go
imp, err := csv.NewMapImporter(specs, nil) // 第二参数 nil 表示不附加 MapAdapter 选项
if err != nil { return err }

result, rowErrors, err := imp.Import(reader)
if err != nil { return err }

rows := result.([]map[string]any)
```

行为细节：

- 源文件中未知的 header **直接跳过**——多余的列不会让导入失败。
- schema 中存在但源文件没有的列，**不会**出现在解析结果的 map 中（key 不存在，而不是零值）。这样 `Required` 与行级校验可以区分「缺失」和「显式零值」。
- 单元格经 `TrimSpace`（CSV 与 Excel 默认均开启，可用 `WithoutTrimSpace()` 关闭）+ `Default` 兜底后，仍为空字符串时会跳过 `Set`：结构体保持零值，map 保留 key 缺失。
- 单元格解析错误、行 Commit 错误、校验错误都会聚合到 `[]tabular.ImportError`，文件其余部分继续处理。

### 行级校验

将 `MapOption` 作为 `NewMapImporter` 的第二个参数传入：

```go
imp, err := csv.NewMapImporter(specs,
    []tabular.MapOption{
        tabular.WithRowValidator(func(row map[string]any) error {
            if row["name"] == "" {
                return errors.New("name must not be empty")
            }
            return nil
        }),
    },
)
```

单元格级校验在每列的 `Validators` 中配置：

```go
specs := []tabular.ColumnSpec{
    {
        Key:  "email",
        Name: "邮箱",
        Type: reflect.TypeFor[string](),
        Validators: []tabular.CellValidator{
            func(col *tabular.Column, value any) error {
                s, _ := value.(string)
                if !strings.Contains(s, "@") {
                    return fmt.Errorf("invalid email: %q", s)
                }
                return nil
            },
        },
    },
}
```

`ColumnSpec.Required`、每列的 `Validators`、以及 map 级 `RowValidator` 都会在 map-row commit 时执行。`Required`、`Validators` 与 `RowValidator` 的错误会通过 `errors.Join` 合并成单行的一个错误。要枚举所有叶子错误：

```go
for _, ie := range rowErrors {
    if errors.Is(ie.Err, tabular.ErrRequiredMissing) {
        // …
    }
    if multi, ok := ie.Err.(interface{ Unwrap() []error }); ok {
        for _, leaf := range multi.Unwrap() {
            log.Warn(leaf)
        }
    }
}
```

## 自定义 Formatter 与 Parser

实现 `tabular` 中的小接口，或者用 `tabular.FormatterFunc` / `tabular.ParserFunc` 把普通函数适配进去。每列按以下三级优先级解析（见 `tabular/resolver.go`）：

1. **`Column.FormatterFn` / `Column.ParserFn`**——直接绑定在列上的实例（最高优先级）
2. **`Column.Formatter` / `Column.Parser`**——根据名字到 importer / exporter 注册表查找，使用 `RegisterFormatter` / `RegisterParser` 注册
3. **默认 formatter / parser**——使用 `Column.Format` 处理日期、浮点等

`ResolveFormatter(col, registry)` / `ResolveParser(col, registry)` 对单列执行这套优先级；`ResolveFormatters` / `ResolveParsers` 一次性为每一列都执行一遍，结果与 `schema.Columns()` 对齐，driver 因此不需要在逐个单元格处理时重复查表。`tabular.IsDefaultFormatter(col, registry)` 判断某列是否解析到内置默认实现（没有 `FormatterFn`，也没有注册命名 `Formatter`）；Excel 导出器用它来决定是写 native typed cell 还是写格式化字符串（见 [Excel → Native Typed Cell](#native-typed-cell)）。

直接绑定（最高优先级）示例：

```go
yenFormatter := tabular.FormatterFunc(func(v any) (string, error) {
    return fmt.Sprintf("¥%.2f", v), nil
})

specs := []tabular.ColumnSpec{
    {
        Key:         "price",
        Name:        "Price",
        Type:        reflect.TypeFor[float64](),
        FormatterFn: yenFormatter,
    },
}
```

命名注册表示例（结构体与 map 适配器都适用）：

```go
exp := csv.NewExporterFor[Order]()
exp.RegisterFormatter("currency", currencyFormatter)
// 标签上写了 `formatter=currency` 的列会使用它
```

```go
importer := csv.NewImporterFor[Employee]()

importer.RegisterParser("date", tabular.ParserFunc(func(cellValue string, targetType reflect.Type) (any, error) {
    return time.Parse("01/02/2006", cellValue)
}))
```

## 自定义 RowAdapter

任何数据源都可以通过实现 `tabular.RowAdapter` 接入引擎：

```go
type RowAdapter interface {
    Schema() *Schema
    Reader(data any) (RowReader, error)
    Writer(capacity int) RowWriter
}

type RowReader interface {
    All() iter.Seq2[int, RowView]
}

type RowView interface {
    Get(column *Column) (any, error)
}

type RowWriter interface {
    NewRow() RowBuilder
    Commit(row RowBuilder) error
    Build() any
}

type RowBuilder interface {
    Set(column *Column, value any) error
    Validate() error
    Value() any
}
```

`tabular.NewStructAdapter(typ)` / `tabular.NewStructAdapterFor[T]()` 和 `tabular.NewMapAdapter(schema, opts...)` / `tabular.NewMapAdapterFromSpecs(specs, opts...)` 是两个内置实现。自定义 adapter 可以直接接入：

```go
adapter := myStreamingAdapter()
imp := csv.NewImporter(adapter)
exp := excel.NewExporter(adapter)
```

适合用于 channel 流式数据源、JOIN 视图、或既不是结构体也不是 map 的业务类型。

## 默认类型支持

内置的 `DefaultParser` 和 `DefaultFormatter`（`tabular.NewDefaultParser(format)` / `tabular.NewDefaultFormatter(format)`）会自动处理以下类型，CSV 与 Excel 都一样：

| Go 类型 | 导入（解析） | 导出（格式化） |
| --- | --- | --- |
| `string` | 直接赋值 | 直接输出 |
| `int`、`int8`–`int64` | 整数解析 | 整数格式化 |
| `uint`、`uint8`–`uint64` | 无符号整数解析 | 整数格式化 |
| `float32`、`float64` | 浮点数解析 | 浮点数格式化 |
| `bool` | `true`/`false`、`1`/`0` | 布尔格式化 |
| `decimal.Decimal` | Decimal 字符串解析 | Decimal 格式化 |
| `time.Time` | 使用 `format` 属性（默认 `time.DateTime`）| 使用 `format` 属性（默认 `time.DateTime`）|
| `timex.Date` / `timex.DateTime` / `timex.Time` | 使用 `format` 属性 | 使用 `format` 属性 |
| `*T`（指针类型） | 空值为 nil，否则解析 | 优雅处理 nil |

## 错误类型

`ImportError` 和 `ExportError` 都实现了标准的 `Unwrap() error` 方法，因此标准库的 `errors.Unwrap`、`errors.Is` 和 `errors.As` 都能直接作用于它们，无需任何 `tabular` 专属的解包辅助函数。

### ImportError

```go
type ImportError struct {
    Row    int    // 基于 1 的行号（包含表头行）
    Column string // 列标题名称
    Field  string // 结构体字段名
    Err    error  // 底层错误
}
```

`ImportError` 实现了 `error` 和 `Unwrap() error`（`ImportError.Unwrap`）。当单行产生多个失败（例如多个 `Required` miss，或者 cell validator 与 row validator 同时失败）时，`Err` 本身可能通过 `errors.Join` 携带多个叶子错误——对 `ImportError` 使用 `errors.Is` 可以匹配特定原因；如果要枚举所有叶子错误，把 `Err` 断言为 `interface{ Unwrap() []error }`。

导入错误按行返回，不会中断导入过程。这允许批量处理：有效行被导入，无效行被报告。

### ExportError

```go
type ExportError struct {
    Row    int    // 基于 0 的数据行索引
    Column string // 列标题名称
    Field  string // 结构体字段名
    Err    error  // 底层错误
}
```

`ExportError` 同样实现了 `error` 和 `Unwrap() error`（`ExportError.Unwrap`）。

### 共享错误哨兵

`tabular` 暴露一组两个驱动共享的错误。请使用 `errors.Is` 判断，不要做字符串匹配：

| 错误 | 触发场景 |
| --- | --- |
| `ErrDataMustBeSlice` | 导出参数不是切片 |
| `ErrSchemaMismatch` | 元素类型与适配器 schema 不匹配（结构体 / map 不一致）|
| `ErrUnknownColumn` | 调用方引用了 schema 中不存在的列 |
| `ErrRequiredMissing` | 动态导入时 `Required` 单元格为空 |
| `ErrNoDataRowsFound` | 经 skip-rows 与可选 header 处理后没有数据行 |
| `ErrDuplicateHeaderName` | Header 行存在重复非空名，或动态 schema 中两列解析出同一个 name |
| `ErrDuplicateColumnKey` | 动态 `ColumnSpec` 切片中有两个条目 `Key` 相同 |
| `ErrUnsetField` | 结构体字段不可写（通常是未导出字段）|
| `ErrMissingColumnKey` / `ErrMissingColumnType` | `ColumnSpec` 缺关键字段 |
| `ErrTypedRowMismatch` | `TypedImporter[T]` / `TypedExporter[T]` 收到的元素类型不是 `T` |
| `ErrUnsupportedType` | 默认 parser 被要求解析成一个它不认识的 Go 类型 |

`excel.ErrSheetIndexOutOfRange` 是唯一的驱动专属 sentinel——见 [Excel → 错误处理](#错误处理-1)。

## CRUD 集成

`Export` 和 `Import` CRUD 构建器内部使用 `tabular`：

```go
// 导出构建器
crud.NewExport[Employee, EmployeeSearch]().
    WithDefaultFormat("excel")

// 导入构建器
crud.NewImport[Employee]().
    WithDefaultFormat("excel").
    WithPreImport(func(ctx context.Context, models []Employee) error {
        // 插入前验证或转换
        return nil
    })
```

`WithDefaultFormat` 接受一个 `crud.TabularFormat`（`crud.FormatExcel` / `crud.FormatCsv`，或等价的字符串常量 `"excel"` / `"csv"`），当请求未显式指定格式时使用。完整的构建器 API（包括 `WithExcelOptions`、`WithCsvOptions`、`WithPreExport`、`WithFilenameBuilder`）参见 [CRUD → 导出与导入 Builder](../data-access/crud)。

## CSV

`csv` 包基于上文描述的共享 `tabular` 引擎提供 CSV 导入/导出。它与 `excel` 共享相同的 `tabular.Importer` 和 `tabular.Exporter` 接口，便于在不改动模型定义的情况下切换格式。

### 包结构

```go
csv.NewImporter(adapter, opts...)            // tabular.Importer
csv.NewExporter(adapter, opts...)            // tabular.Exporter
csv.NewImporterFor[T](opts...)               // 结构体快捷方式
csv.NewExporterFor[T](opts...)               // 结构体快捷方式
csv.NewTypedImporterFor[T](opts...)          // 直接返回 []T
csv.NewTypedExporterFor[T](opts...)          // 直接接受 []T
csv.NewMapImporter(specs, mapOpts, opts...)  // 动态 map importer
csv.NewMapExporter(specs, opts...)           // 动态 map exporter
```

CSV 的 option marker 类型是 `csv.ExportOption` 与 `csv.ImportOption`。`csv.NewMapExporter` 用 `tabular.NewMapAdapterFromSpecs(specs)` 校验 `specs`，**不**接受 `mapOpts`；`csv.NewMapImporter` 用 `tabular.NewMapAdapterFromSpecs(specs, mapOpts...)` 校验——不需要 row validator 时 `mapOpts` 传 `nil` 即可。

| 选项 | 默认值 | 用途 |
| --- | --- | --- |
| `csv.WithImportDelimiter(r)` | `,` | 导入时的字段分隔符 |
| `csv.WithoutHeader()` | 含 header | 第一行作为数据，按 schema 顺序按位置映射 |
| `csv.WithSkipRows(n)` | `0` | 读取前跳过 n 行；负数会归零 |
| `csv.WithoutTrimSpace()` | 自动 trim | 关闭单元格首尾空白裁剪（同时影响空行检测与 header 匹配）|
| `csv.WithComment(r)` | 无（`0`）| 以该字符开头的行被忽略 |
| `csv.WithExportDelimiter(r)` | `,` | 导出时的字段分隔符 |
| `csv.WithoutWriteHeader()` | 写 header | 导出时不写 header 行 |
| `csv.WithCRLF()` | LF | 使用 Windows 风格的换行符 |

### 模型定义

CSV 使用与 [`tabular` 标签](#tabular-标签) 相同的结构体标签：

```go
type Employee struct {
    orm.FullAuditedModel `tabular:"-"`

    Name       string          `tabular:"姓名,width=20"`
    Email      string          `tabular:"邮箱,width=30"`
    Department string          `tabular:"部门"`
    JoinDate   timex.Date      `tabular:"入职日期,format=2006-01-02"`
    Salary     decimal.Decimal `tabular:"薪资"`
    IsActive   bool            `tabular:"是否在职"`
}
```

### 导出

```go
import "github.com/coldsmirk/vef-framework-go/csv"

exporter := csv.NewExporterFor[Employee]()

// 导出到文件
err := exporter.ExportToFile(employees, "employees.csv")

// 导出到缓冲区（用于 HTTP 响应）
buf, err := exporter.Export(employees)
```

```go
// TSV 导出，使用 Windows 换行符
exporter := csv.NewExporterFor[Employee](
    csv.WithExportDelimiter('\t'),
    csv.WithCRLF(),
)

// 不写入表头行
exporter := csv.NewExporterFor[Employee](
    csv.WithoutWriteHeader(),
)
```

### 导入

```go
importer := csv.NewImporterFor[Employee]()

// 从文件导入
data, importErrors, err := importer.ImportFromFile("employees.csv")
if err != nil {
    return err // 致命错误
}

// 检查行级错误
for _, e := range importErrors {
    log.Printf("Row %d: %v", e.Row, e.Err)
}

employees := data.([]Employee)

// 或者从 io.Reader 导入（例如上传的文件）
data, importErrors, err = importer.Import(reader)
```

```go
// TSV 文件，有 2 行标题，# 开头为注释行
importer := csv.NewImporterFor[Employee](
    csv.WithImportDelimiter('\t'),
    csv.WithSkipRows(2),
    csv.WithComment('#'),
)

// 无表头的 CSV（按位置/顺序匹配列）
importer := csv.NewImporterFor[Employee](
    csv.WithoutHeader(),
)
```

`Import` 会先调用标准库 CSV reader 的 `ReadAll` 再解析，因此峰值内存会随文件大小和最终结果切片一起增长。reader 使用 `FieldsPerRecord = -1`，因此允许不等长行，缺失的映射单元格由 tabular adapter 处理（空单元格 → `Default` → 跳过）。`WithoutTrimSpace()` 关闭的这层 trim 由 VEF 的 `tabular` 层执行；底层 Go 标准库 CSV reader 本身不启用 `TrimLeadingSpace`。

#### 有表头 vs 无表头模式

| 模式 | 列映射方式 |
| --- | --- |
| 有表头（默认）| 表头名称 → `tabular` 标签名 |
| 无表头 | 列位置 → `tabular` 字段顺序 |

使用 `WithoutHeader()` 时，列按位置匹配。使用 `order` 标签属性控制字段排序：

```go
type Record struct {
    Name  string `tabular:"Name,order=0"`
    Email string `tabular:"Email,order=1"`
    Age   int    `tabular:"Age,order=2"`
}
```

`WithSkipRows` 会在 header 检测前生效。启用表头时，`rows[skipRows]` 是 header 行，数据从其后一行开始；使用 `csv.WithoutHeader()` 时，第一个未跳过的行会作为数据解析。Import error 中的行号是基于 1 的 CSV 文件行号，并包含 skip/header 偏移，方便直接对应文本编辑器里看到的行。

### 自定义 Formatter 与 Parser

```go
exporter := csv.NewExporterFor[Employee]()

exporter.RegisterFormatter("status", tabular.FormatterFunc(func(value any) (string, error) {
    if active, ok := value.(bool); ok && active {
        return "Y", nil
    }
    return "N", nil
}))
```

```go
importer := csv.NewImporterFor[Employee]()

importer.RegisterParser("date", tabular.ParserFunc(func(cellValue string, targetType reflect.Type) (any, error) {
    return time.Parse("01/02/2006", cellValue)
}))
```

### 验证

导入的记录会自动使用 `validator.Validate(...)` 进行验证，与 Excel 导入器相同。

### 错误处理

顶层错误表示致命的读写或结构错误，包括 `ReadAll` 失败、没有数据行（`tabular.ErrNoDataRowsFound`）、导出 schema 不匹配（`tabular.ErrDataMustBeSlice`）、重复 header（`tabular.ErrDuplicateHeaderName`），以及最终 writer flush 失败（`flush CSV writer: ...`）。解析失败、validator 失败和 adapter commit 失败（包括 `tabular.ErrUnsetField`）会聚合进 `[]tabular.ImportError`；import 可以返回 `err == nil` 同时带有非空行级错误，对应行会被跳过，后续行继续处理。

## Excel

`excel` 包基于共享的 `tabular` 引擎提供 Excel 导入/导出。底层使用 [excelize](https://github.com/xuri/excelize)，并集成了 VEF 的验证系统。

### 包结构

```go
excel.NewImporter(adapter, opts...)
excel.NewExporter(adapter, opts...)
excel.NewImporterFor[T](opts...)
excel.NewExporterFor[T](opts...)
excel.NewTypedImporterFor[T](opts...)
excel.NewTypedExporterFor[T](opts...)
excel.NewMapImporter(specs, mapOpts, opts...)
excel.NewMapExporter(specs, opts...)
```

Excel 的 option marker 类型是 `excel.ExportOption` 与 `excel.ImportOption`。`excel.NewMapExporter` 用 `tabular.NewMapAdapterFromSpecs(specs)` 校验 `specs`；`excel.NewMapImporter` 用 `tabular.NewMapAdapterFromSpecs(specs, mapOpts...)` 校验——不需要 row validator 时 `mapOpts` 传 `nil` 即可。

| 选项 | 默认值 | 用途 |
| --- | --- | --- |
| `excel.WithSheetName(name)` | `Sheet1` | 导出工作表名；会重命名默认 sheet，而不是创建第二个 sheet |
| `excel.WithImportSheetName(name)` | 无 | 按名称读取工作表，优先于 `WithImportSheetIndex` |
| `excel.WithImportSheetIndex(i)` | `0` | 按 0-based index 读取工作表（负数或越界返回 `excel.ErrSheetIndexOutOfRange`）；设置了 `WithImportSheetName` 时会被忽略 |
| `excel.WithSkipRows(n)` | `0` | 读取前跳过 n 行；负数会归零 |
| `excel.WithoutHeader()` | 含 header | 第一个未跳过的行是数据，按位置映射 |
| `excel.WithoutTrimSpace()` | 自动 trim | 关闭单元格首尾空白裁剪（同时影响空行检测与 header 匹配）|

`ColumnSpec` 中设置的 `Column.Width`（或结构体标签 `width=…`）会作用到生成的工作表列宽；CSV 会忽略它。

### 工作表

Excel workbook 可以包含多个工作表。导出总是写入单个具名 sheet（`excel.WithSheetName`，默认 `Sheet1`）；导入既可以按名称选择源 sheet（`excel.WithImportSheetName`，设置后优先生效），也可以按 0-based index 选择（`excel.WithImportSheetIndex`，默认 `0`）。

### 模型示例

```go
type Employee struct {
    orm.FullAuditedModel `tabular:"-"`

    Name       string          `json:"name" bun:"name" tabular:"姓名,width=20"`
    Email      string          `json:"email" bun:"email" tabular:"邮箱,width=30"`
    Department string          `json:"department" bun:"department" tabular:"部门,width=15"`
    JoinDate   timex.Date      `json:"joinDate" bun:"join_date" tabular:"入职日期,format=2006-01-02,width=15"`
    Salary     decimal.Decimal `json:"salary" bun:"salary" tabular:"薪资,width=12"`
    IsActive   bool            `json:"isActive" bun:"is_active" tabular:"是否在职,width=10"`
}
```

### 导出

```go
import "github.com/coldsmirk/vef-framework-go/excel"

// 创建类型化的导出器
exporter := excel.NewExporterFor[Employee]()

// 导出到文件
err := exporter.ExportToFile(employees, "employees.xlsx")

// 导出到缓冲区（用于 HTTP 响应）
buf, err := exporter.Export(employees)
```

```go
// 自定义工作表名（默认："Sheet1"）
exporter := excel.NewExporterFor[Employee](
    excel.WithSheetName("Employees"),
)
```

#### Native Typed Cell

Excel 导出在列使用默认 formatter，且没有显式 `format`、`formatter` 或 `FormatterFn`（通过 `tabular.IsDefaultFormatter` 判断）时，会写入 native typed cell。整数、浮点数、布尔值、`time.Time`、`timex.Date` 和 `timex.DateTime` 都以这种方式在 Excel 里保持可排序或可求和。一旦列设置了格式字符串或自定义 formatter，导出器会改为把格式化结果按文本写入。

Native cell 转换的细节：nil pointer 会写成空单元格，非 nil pointer 会先解引用；`timex.Date` / `timex.DateTime` 会被解包成 `time.Time`，让 excelize 存储 native 日期(时间)单元格；`decimal.Decimal` 会转换成 `float64`（在约 15–16 位有效数字内精确，超出则有损——需要完整 decimal 精度的列应改为声明显式 `format`，按精确文本渲染）；`timex.Time` 会被刻意保留为文本，因为它的 zero-date 部分早于 Excel epoch，直接转换会渲染出一个错误的日期。

#### 自定义 Formatter 与 Parser

```go
exporter := excel.NewExporterFor[Employee]()

// 注册名为 "status" 的自定义格式化器
exporter.RegisterFormatter("status", tabular.FormatterFunc(func(value any) (string, error) {
    if active, ok := value.(bool); ok && active {
        return "Active", nil
    }
    return "Inactive", nil
}))
```

然后在结构体标签中引用：

```go
IsActive bool `tabular:"Status,formatter=status"`
```

```go
importer := excel.NewImporterFor[Employee]()

importer.RegisterParser("date", tabular.ParserFunc(func(cellValue string, targetType reflect.Type) (any, error) {
    return time.Parse("01/02/2006", cellValue)
}))
```

然后在结构体标签中引用：

```go
JoinDate time.Time `tabular:"Join Date,parser=date"`
```

#### HTTP Handler 中导出

```go
func (r *EmployeeResource) Export(ctx fiber.Ctx, db orm.DB) error {
    var employees []Employee
    err := db.NewSelect().Model(&employees).Scan(ctx.Context())
    if err != nil {
        return err
    }

    exporter := excel.NewExporterFor[Employee]()
    buf, err := exporter.Export(employees)
    if err != nil {
        return err
    }

    ctx.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
    ctx.Set("Content-Disposition", "attachment; filename=employees.xlsx")
    return ctx.Send(buf.Bytes())
}
```

### 导入

```go
// 创建类型化的导入器
importer := excel.NewImporterFor[Employee]()

// 从文件导入
data, importErrors, err := importer.ImportFromFile("employees.xlsx")
if err != nil {
    // 致命错误（文件不存在等）
    return err
}

// 检查行级错误
if len(importErrors) > 0 {
    for _, e := range importErrors {
        log.Printf("Row %d, Column %s: %v", e.Row, e.Column, e.Err)
    }
}

// 类型断言获取结果
employees := data.([]Employee)

// 或者从上传文件导入（io.Reader）
data, importErrors, err = importer.Import(reader)
```

```go
// 按工作表名指定
importer := excel.NewImporterFor[Employee](
    excel.WithImportSheetName("Staff"),
)

// 按工作表索引指定（默认：0）
importer := excel.NewImporterFor[Employee](
    excel.WithImportSheetIndex(1),
)

// 跳过前导行（如标题行在表头之前）
importer := excel.NewImporterFor[Employee](
    excel.WithSkipRows(2),
)

// 不读取表头：第一个未跳过的行按 schema 位置当作数据解析
importer := excel.NewImporterFor[Employee](
    excel.WithoutHeader(),
)

// 保留 header 与单元格首尾空白
importer := excel.NewImporterFor[Employee](
    excel.WithoutTrimSpace(),
)
```

`WithImportSheetName` 优先于 `WithImportSheetIndex`。负数 `WithSkipRows` 会归零。启用表头时，跳过行会在 header 解析前生效；关闭表头时，第一个未跳过的行会按位置映射作为数据解析。和 CSV 一样，Excel 导入默认会 trim 单元格值，并在解析前把 workbook rows 读入内存——excelize 会一次性加载整个 workbook，因此峰值内存会随文件大小和最终结果切片一起增长。

### 验证

导入的记录会自动使用 `validator.Validate(...)` 进行验证。如果验证失败，该行会被添加到 `importErrors` 并从结果切片中跳过。

```go
type Employee struct {
    Name  string `tabular:"Name" validate:"required"`
    Email string `tabular:"Email" validate:"required,email"`
}
```

### 列映射规则

1. 导入器通过 Excel 表头名称 → `tabular` 标签名（无标签时使用字段名）进行匹配。
2. 未匹配的 Excel 列会被静默忽略。
3. 缺失的 Excel 列会让结构体字段保持零值（如指定了 `default` 则使用默认值）。
4. 空行会被自动跳过。

使用 `excel.WithoutHeader()` 时会绕过 header 匹配，改用 `tabular.DefaultPositionalMapping`：源文件第 1 列映射 schema 第 1 列，第 2 列映射 schema 第 2 列，依此类推。

Excel 导入/导出的默认类型支持与 CSV 完全相同——见上文的[默认类型支持](#默认类型支持)。

### 错误处理

| 错误 | 含义 |
| --- | --- |
| `excel.ErrSheetIndexOutOfRange` | 配置的工作表索引为负数或超出可用范围 |
| `tabular.ErrNoDataRowsFound` | 经 `WithSkipRows` 与可选表头处理后没有数据行 |
| `tabular.ErrDuplicateHeaderName` | 表头中存在重复的非空列名 |
| `tabular.ErrUnsetField` | 结构体字段无法设置，通常是未导出字段 |

顶层 import error 表示致命的文件或工作表错误。解析失败、validator 失败和 adapter commit 失败会聚合进 `[]tabular.ImportError`（`Row`/`Column`/`Field`/`Err`，基于 1 的行号，包含表头行）；对应行会被跳过，后续行继续处理。Export error 使用相同的 `ExportError` 结构，行索引基于 0。

## 速查

```go
// 静态结构体往返
imp := csv.NewTypedImporterFor[User]()
exp := csv.NewTypedExporterFor[User]()
buf, _ := exp.Export(users)
imported, errs, _ := imp.Import(buf)

// 动态 map 往返
specs := []tabular.ColumnSpec{
    {Key: "id",   Name: "ID",   Type: reflect.TypeFor[int](),    Required: true},
    {Key: "name", Name: "Name", Type: reflect.TypeFor[string]()},
}
exp, _ := excel.NewMapExporter(specs, excel.WithSheetName("Data"))
imp, _ := excel.NewMapImporter(specs, nil)
buf, _ := exp.Export([]map[string]any{{"id": 1, "name": "Alice"}})
rows, errs, _ := imp.Import(buf)

// 动态 + 行级校验
imp, _ := csv.NewMapImporter(specs,
    []tabular.MapOption{tabular.WithRowValidator(func(r map[string]any) error {
        if r["id"].(int) <= 0 { return errors.New("id must be positive") }
        return nil
    })},
)
```

`crud.NewExport` / `crud.NewImport` 在以上工厂之上做了进一步封装，参见 [CRUD → 导出与导入 Builder](../data-access/crud)。
