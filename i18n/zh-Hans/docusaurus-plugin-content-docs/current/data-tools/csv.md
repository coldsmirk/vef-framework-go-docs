---
sidebar_position: 6
---

# CSV

[tabular 导入导出核心](./tabular)的 CSV 后端，为分隔符文件实现共享的 schema、格式化器与解析器契约。

`csv` 包基于上文描述的共享 `tabular` 引擎提供 CSV 导入/导出。它与 `excel` 共享相同的 `tabular.Importer` 和 `tabular.Exporter` 接口，便于在不改动模型定义的情况下切换格式。

## 包结构

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

## 模型定义

CSV 使用与 [`tabular` 标签](./tabular#tabular-标签) 相同的结构体标签：

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

## 导出

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

## 导入

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

### 有表头 vs 无表头模式

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

## 自定义 Formatter 与 Parser

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

## 验证

导入的记录会自动使用 `validator.Validate(...)` 进行验证，与 Excel 导入器相同。

## 错误处理

顶层错误表示致命的读写或结构错误，包括 `ReadAll` 失败、没有数据行（`tabular.ErrNoDataRowsFound`）、导出 schema 不匹配（`tabular.ErrDataMustBeSlice`）、重复 header（`tabular.ErrDuplicateHeaderName`），以及最终 writer flush 失败（`flush CSV writer: ...`）。解析失败、validator 失败和 adapter commit 失败（包括 `tabular.ErrUnsetField`）会聚合进 `[]tabular.ImportError`；import 可以返回 `err == nil` 同时带有非空行级错误，对应行会被跳过，后续行继续处理。

## 下一步

- [表格导入导出](./tabular) — 共享的 schema、标签与接口
- [Excel](./excel) — Excel 后端
