---
sidebar_position: 9
---

# CSV / Excel 导入导出

VEF 在 [`tabular`](#包结构) 包中提供统一的表格引擎，并通过 [`csv`](#csv-包) 与 [`excel`](#excel-包) 两个轻量驱动包暴露具体格式。两者拥有完全对称的工厂函数，所有读写都经由同一个 `RowAdapter` 抽象。这意味着：

- **静态行**（用 `tabular` 标签描述的 Go 结构体）和
- **动态行**（运行期定义列、数据是 `map[string]any`）

共用同一条 importer / exporter 流水线。你只需选择合适的适配器，格式驱动负责剩下的工作。

## 何时使用哪种用法

| 场景 | 推荐工厂 |
| --- | --- |
| 已有 Go 结构体（例如某个 model）描述了所有列 | `csv.NewImporterFor[T]` / `excel.NewExporterFor[T]`（以及 `*Typed*` 变体） |
| 列在运行期才确定——多租户表格、用户自定义模板、动态表单 | `csv.NewMapImporter` / `excel.NewMapExporter`，由 `[]tabular.ColumnSpec` 驱动 |
| 行数据自带特殊来源（channel、自定义业务类型等） | 自行实现 `tabular.RowAdapter`，传入 `csv.NewImporter` / `excel.NewExporter` |

导入返回类型由适配器决定：

- 结构体适配器 → `[]T`
- map 适配器 → `[]map[string]any`

## 包结构

```
tabular/   // schema、列、适配器、formatter / parser、错误
  ├── adapter.go        // RowAdapter, RowReader, RowView, RowWriter, RowBuilder
  ├── schema.go         // Schema, Column
  ├── struct_adapter.go // StructAdapter（结构体 + 框架 validator）
  ├── map_adapter.go    // MapAdapter（map + Required / Validators / RowValidator）
  ├── spec.go           // ColumnSpec, NewSchemaFromSpecs, NewMapAdapterFromSpecs
  ├── resolver.go       // FormatterFn / Formatter / 默认实现的优先级解析
  ├── mapping.go        // Header → schema 列映射
  ├── typed.go          // TypedImporter[T] / TypedExporter[T] 包装器
  └── errors.go         // 共享错误（ErrRequiredMissing、ErrSchemaMismatch 等）

csv/       // CSV 驱动：NewImporter、NewExporter、NewImporterFor、NewExporterFor、
           //          NewMapImporter、NewMapExporter、NewTyped*For
excel/     // Excel 驱动，与 csv 对称，外加 sheet 等 Excel 专属选项
```

`csv` / `excel` 不再持有任何 model 相关的反射或私有错误类型。常用错误统一定义在 `tabular`（`ErrDataMustBeSlice`、`ErrRequiredMissing`、`ErrUnknownColumn`、`ErrSchemaMismatch` 等），仅 `excel.ErrSheetIndexOutOfRange` 保留在驱动内部。

## 静态结构体用法

在结构体字段上打 `tabular` 标签：

```go
type User struct {
    ID       int       `tabular:"name=用户ID;order=1"`
    Name     string    `tabular:"name=姓名;order=2" validate:"required"`
    Birthday time.Time `tabular:"name=生日;format=2006-01-02;order=3"`
    Active   bool      `tabular:"name=激活;default=false;order=4"`
    Internal string    `tabular:"-"` // 忽略
}
```

支持的标签属性（来自 `tabular/constants.go`）：

| 属性 | 含义 |
| --- | --- |
| `name` | Header 名称（默认使用字段名） |
| `order` | 列顺序（稳定排序，从小到大） |
| `width` | 列宽提示（Excel 使用） |
| `default` | 导入时，源单元格为空使用该默认值 |
| `format` | 默认 formatter / parser 的格式模板，例如 `"2006-01-02"`、`"%.2f"` |
| `formatter` | 已注册 formatter 的名字（导出端） |
| `parser` | 已注册 parser 的名字（导入端） |
| `dive` | 递归展开内嵌结构体 |
| `-` | 忽略该字段 |

校验委托给框架的 `validator` 包——继续使用 `validate:"…"` 标签即可，提交每行时会自动执行。

### 导出

```go
exp := csv.NewExporterFor[User]()
buf, err := exp.Export(users) // users: []User 或 []*User
// 或者直接写入磁盘：
err = exp.ExportToFile(users, "users.csv")
```

Excel：

```go
exp := excel.NewExporterFor[User](excel.WithSheetName("Users"))
buf, err := exp.Export(users)
```

### 导入

```go
imp := csv.NewImporterFor[User]()
result, rowErrors, err := imp.Import(reader)
if err != nil {
    return err // 顶层失败（例如文件损坏）
}
users := result.([]User)
for _, ie := range rowErrors {
    log.Warnf("row %d column %s: %v", ie.Row, ie.Column, ie.Err)
}
```

行级失败（解析错误、结构体校验失败等）会聚合到 `[]tabular.ImportError`，**不会中断**整体导入。顶层 `error` 仅在出现致命问题（无法读取文件、Header 行损坏等）时返回。

### Typed 包装器

`any` 返回值有时不便使用，两个包都提供了泛型包装器替你做类型断言：

```go
imp := csv.NewTypedImporterFor[User]()
users, rowErrors, err := imp.Import(reader) // users 直接是 []User

exp := csv.NewTypedExporterFor[User]()
buf, err := exp.Export(users)               // 直接接受 []User
```

`TypedImporter` / `TypedExporter` 包裹底层的 `tabular.Importer` / `tabular.Exporter`。如需直接调用 `RegisterParser` / `RegisterFormatter`，可以使用 `Inner()` 取出内部实例。

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
| `Default` | 否 | 源单元格为空时的默认值 |
| `Format` | 否 | 默认 formatter / parser 的模板（日期、浮点等） |
| `Formatter` / `Parser` | 否 | 名字，从 importer / exporter 的注册表查找 |
| `FormatterFn` / `ParserFn` | 否 | 直接绑定的 `tabular.Formatter` / `tabular.ValueParser` 实例 |
| `Required` | 否 | 导入时空值会触发 `ErrRequiredMissing` |
| `Validators` | 否 | 解析后执行的 `[]CellValidator` |

`NewSchemaFromSpecs` 会立即校验输入：缺 `Key`、缺 `Type`、`Key` 重复都会在构造期返回错误（`tabular.ErrMissingColumnKey`、`ErrMissingColumnType`、`ErrDuplicateColumnName`）。

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
- 单元格经 `TrimSpace`（CSV 默认开启）+ `Default` 兜底后，仍为空字符串时会跳过 `Set`：结构体保持零值，map 保留 key 缺失。
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

`Required`、`Validators` 与 `RowValidator` 的错误会通过 `errors.Join` 合并成单行的一个错误。要枚举所有叶子错误：

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

实现 `tabular` 中的两个轻量接口：

```go
type Formatter interface {
    Format(value any) (string, error)
}

type ValueParser interface {
    Parse(cellValue string, targetType reflect.Type) (any, error)
}
```

`tabular.FormatterFunc` 与 `tabular.ParserFunc` 可将普通函数适配到这两个接口。

每列按以下三级优先级解析（见 `tabular/resolver.go`）：

1. **`Column.FormatterFn` / `Column.ParserFn`**——直接绑定在列上的实例（最高优先级）
2. **`Column.Formatter` / `Column.Parser`**——根据名字到 importer / exporter 注册表查找，使用 `RegisterFormatter` / `RegisterParser` 注册
3. **默认 formatter / parser**——使用 `Column.Format` 处理日期、浮点等

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

## 自定义 RowAdapter

任何数据源都可以通过实现 `tabular.RowAdapter` 接入引擎：

```go
type RowAdapter interface {
    Schema() *Schema
    Reader(data any) (RowReader, error)
    Writer(capacity int) RowWriter
}
```

之后直接交给驱动即可：

```go
adapter := myStreamingAdapter()
imp := csv.NewImporter(adapter)
exp := excel.NewExporter(adapter)
```

适合用于 channel 流式数据源、JOIN 视图、或既不是结构体也不是 map 的业务类型。

## CSV 包

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

CSV 选项：

| 选项 | 默认值 | 用途 |
| --- | --- | --- |
| `WithImportDelimiter(r)` | `,` | 导入时的字段分隔符 |
| `WithoutHeader()` | 含 header | 第一行作为数据，按 schema 顺序按位置映射 |
| `WithSkipRows(n)` | `0` | 读取前跳过 n 行 |
| `WithoutTrimSpace()` | 自动 trim | 关闭单元格首尾空白裁剪 |
| `WithComment(r)` | 无 | 以该字符开头的行被忽略 |
| `WithExportDelimiter(r)` | `,` | 导出时的字段分隔符 |
| `WithoutWriteHeader()` | 写 header | 导出时不写 header 行 |
| `WithCrlf()` | LF | 使用 Windows 风格的换行符 |

## Excel 包

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

Excel 选项：

| 选项 | 默认值 | 用途 |
| --- | --- | --- |
| `WithSheetName(name)` | `Sheet1` | 导出工作表名 |
| `WithImportSheetName(name)` | 无 | 按名称读取工作表 |
| `WithImportSheetIndex(i)` | `0` | 按索引读取工作表（越界返回 `excel.ErrSheetIndexOutOfRange`） |
| `WithSkipRows(n)` | `0` | 读取前跳过 n 行 |
| `WithoutHeader()` | 含 header | 第一行作为数据，按位置映射 |

`ColumnSpec` 中设置的 `Width`（或结构体标签 `width=…`）会作用到生成的工作表列宽。

## Header → 列映射规则

两个驱动都通过 `tabular.BuildHeaderMapping` 解析 header：

- Header 单元格按 `Column.Name` 匹配。
- 空 Header 单元格被跳过。
- 未知 Header 单元格被跳过（不会因为多余列失败）。
- 重复的非空 Header 是致命错误：`tabular.ErrDuplicateColumnName`。
- 当 importer 配置为 `WithoutHeader()` 时，引擎回退到 `tabular.DefaultPositionalMapping`——第 0 列对应 schema 的第一列，依此类推。

## 错误

`tabular` 暴露一组共享错误。常用于 `errors.Is` 判断的有：

| 错误 | 触发场景 |
| --- | --- |
| `ErrDataMustBeSlice` | 导出参数不是切片 |
| `ErrSchemaMismatch` | 元素类型与适配器 schema 不匹配（结构体 / map 不一致） |
| `ErrUnknownColumn` | 调用方引用了 schema 中不存在的列 |
| `ErrRequiredMissing` | 动态导入时 `Required` 单元格为空 |
| `ErrDuplicateColumnName` | Header 行存在重复非空名 |
| `ErrUnsetField` | 结构体字段不可写（通常是未导出字段） |
| `ErrMissingColumnKey` / `ErrMissingColumnType` | `ColumnSpec` 缺关键字段 |
| `ErrTypedRowMismatch` | `TypedImporter[T]` 收到的元素类型不是 `T` |

判断错误时请使用 `errors.Is`，不要做字符串匹配。

## 旧 API 迁移

如果你之前用的是旧签名，请按下表升级：

| 旧写法 | 新写法 |
| --- | --- |
| `csv.NewImporter(typ, opts...)` | `csv.NewImporterFor[T](opts...)`，或 `csv.NewImporter(tabular.NewStructAdapter(typ), opts...)` |
| `excel.NewExporter(typ, opts...)` | `excel.NewExporterFor[T](opts...)` |
| `csv.ErrDataMustBeSlice` 等 | `tabular.ErrDataMustBeSlice` 等共享错误 |

`excel.ErrSheetIndexOutOfRange` 保持不变——它属于 Excel 驱动专有错误。

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

`crud.NewExport` / `crud.NewImport` 在以上工厂之上做了进一步封装，参见 [CRUD → 导出与导入 Builder](../guide/crud.md#导出与导入-builder)。
