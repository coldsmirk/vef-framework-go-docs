---
sidebar_position: 12
---

# Tabular

`tabular` 包是 VEF 结构化数据导入导出的基础。它提供标签驱动的 schema 系统和通用接口，由 [Excel](./excel) 和 [CSV](./csv) 包实现。

## 架构

```
tabular（核心）
├── Schema      — 标签解析、列元数据
├── Importer    — 导入接口
├── Exporter    — 导出接口
├── Formatter   — 导出值格式化
└── ValueParser — 导入值解析

excel（实现）                    csv（实现）
├── excel.NewImporterFor[T]()   ├── csv.NewImporterFor[T]()
└── excel.NewExporterFor[T]()   └── csv.NewExporterFor[T]()
```

## `tabular` 标签

使用结构体标签定义字段如何映射到 Excel/CSV 列：

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

tag parser 使用逗号分隔的 `key=value` 对。分号不是分隔符；
`tabular:"name=ID;order=1"` 会被当作一个 `name` 值。

## Schema

`Schema` 类型在初始化时从结构体字段预解析表格元数据：

```go
schema := tabular.NewSchemaFor[Employee]()

columns := schema.Columns()       // []*Column — 所有解析出的列
names := schema.ColumnNames()     // []string{"姓名", "邮箱", ...}
count := schema.ColumnCount()     // 6
```

列会按 `order` 属性自动排序。未显式指定 `order` 的字段使用其声明顺序。

`NewSchemaFromSpecs` 会在构造期校验动态 schema：缺 `Key` 返回
`ErrMissingColumnKey`，缺 `Type` 返回 `ErrMissingColumnType`，重复 key 返回
`ErrDuplicateColumnKey`，解析后的 header name 重复返回 `ErrDuplicateHeaderName`。

`ColumnSpec.Required`、每列的 `Validators`、以及 map 级 `RowValidator`
会在 map-row 导入时执行。多个 map-row validation failure 会用
`errors.Join` 合并，因此调用方可以继续用 `errors.Is` 检查合并后的错误。

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

`BuildHeaderMapping` 是 CSV 和 Excel driver 共享的 header resolver。启用
trim 时，它会先修剪 header name 再匹配；空 header 和未知 header 会跳过；
重复的非空 header 会作为 fatal `ErrDuplicateHeaderName` 返回。Importer 配置
`WithoutHeader()` 时，driver 使用 `DefaultPositionalMapping`，按输入列位置依次
映射到 schema column。

`ParseRow` 会先应用默认值再解析单元格；默认值替换后仍为空的 cell 会被跳过；
parse、validation 和 commit failure 会作为行级 `ImportError` 返回。如果返回了
row error，row builder 不会提交 partial row。

## 默认类型支持

内置的 `DefaultParser` 和 `DefaultFormatter` 自动处理以下类型：

| Go 类型 | 导入（解析） | 导出（格式化） |
| --- | --- | --- |
| `string` | 直接赋值 | 直接输出 |
| `int`、`int8`–`int64` | 整数解析 | 整数格式化 |
| `uint`、`uint8`–`uint64` | 无符号整数解析 | 整数格式化 |
| `float32`、`float64` | 浮点数解析 | 浮点数格式化 |
| `bool` | `true`/`false`、`1`/`0` | 布尔格式化 |
| `decimal.Decimal` | Decimal 字符串解析 | Decimal 格式化 |
| `timex.Date` / `timex.DateTime` | 使用 `format` 属性 | 使用 `format` 属性 |
| `*T`（指针类型） | 空值为 nil，否则解析 | 优雅处理 nil |

## 错误类型

### ImportError

```go
type ImportError struct {
    Row    int    // 基于 1 的行号（包含表头行）
    Column string // 列标题名称
    Field  string // 结构体字段名
    Err    error  // 底层错误
}
```

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

## 实现

| 包 | 格式 | 文档 |
| --- | --- | --- |
| `excel` | `.xlsx`（Excel）| [Excel 文档](./excel) |
| `csv` | `.csv`（CSV/TSV）| [CSV 文档](./csv) |

## 公开核心 API

| API 组 | 公开 surface |
| --- | --- |
| schema | `NewSchema`, `NewSchemaFor[T]`, `NewSchemaFromSpecs`, `Column`, `ColumnSpec`，以及 `Schema` lookup 方法 |
| adapter | `NewStructAdapter`, `NewStructAdapterFor[T]`, `NewMapAdapter`, `NewMapAdapterFromSpecs`, `RowAdapter`, `RowReader`, `RowWriter`, `RowView`, `RowBuilder` |
| typed wrapper | `NewTypedImporter[T]`, `NewTypedExporter[T]`, `TypedImporter[T]`, `TypedExporter[T]` |
| mapping/parsing | `BuildHeaderMapping`, `DefaultPositionalMapping`, `ColumnMapping`, `NewColumnMapping`, `ParseRow`, `ParseRowOptions`, `MappingOptions`, `MapOption`, `WithRowValidator`, `RowValidator`, `CellValidator`, `IsEmptyRow` |
| formatter/parser registry | `ResolveFormatter`, `ResolveFormatters`, `ResolveParser`, `ResolveParsers`, `IsDefaultFormatter`, `NewDefaultFormatter`, `NewDefaultParser` |
| 常量/错误 | `TagTabular`, `IgnoreField`, `AttrDive`, `AttrName`, `AttrOrder`, `AttrWidth`, `AttrDefault`, `AttrFormatter`, `AttrParser`, `AttrFormat`，以及 `ErrDataMustBeSlice`, `ErrDuplicateColumnKey`, `ErrDuplicateHeaderName`, `ErrMissingColumnKey`, `ErrMissingColumnType`, `ErrNoDataRowsFound`, `ErrRequiredMissing`, `ErrSchemaMismatch`, `ErrTypedRowMismatch`, `ErrUnknownColumn`, `ErrUnsetField`, `ErrUnsupportedType` 等 tabular sentinel |

当前 tabular 包审计在生成的 API ledger 中锁定 **143 public tabular
entries**。分组 member surface 覆盖 **75 grouped tabular field/method
entries**，分布在 **20 tabular receiver/type families** 中：其中包含 **37
exported tabular field entries** 和 **38 exported tabular method entries**。
生成的公开 API 索引仍是完整签名清单；本页负责说明 schema、mapping、
parser/formatter、adapter 和错误契约家族。

额外已审计字段和 adapter 方法：

| Surface | Public API |
| --- | --- |
| column metadata | `Column.Default`, `Column.Width`, `Column.Order`, `Column.Parser`, `Column.ParserFn`, `Column.FormatterFn`, `Column.Index` |
| dynamic specs | `ColumnSpec.Default`, `ColumnSpec.Width`, `ColumnSpec.Order`, `ColumnSpec.Parser`, `ColumnSpec.ParserFn`, `ColumnSpec.FormatterFn` |
| parsing options | `MappingOptions.TrimSpace`, `ParseRowOptions.TrimSpace` |
| row adapter contract | `RowAdapter.Writer`, `RowReader.All`, `RowView.Get`, `RowBuilder.Set`, `RowBuilder.Validate`, `RowWriter.NewRow`, `RowWriter.Commit`, `RowWriter.Build` |
| schema lookup | `Schema.ColumnByKey`, `Schema.ColumnByName` |
| typed wrappers | `TypedExporter.Inner`, `TypedImporter.Inner` |
| error wrapping | `ImportError.Unwrap`, `ExportError.Unwrap`；调用方可以使用 `errors.Unwrap` / `errors.Is` |

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
