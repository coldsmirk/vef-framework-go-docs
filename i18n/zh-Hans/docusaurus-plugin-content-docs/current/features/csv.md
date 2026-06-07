---
sidebar_position: 5
---

# CSV

`csv` 包使用与 [Excel](./excel) 包相同的 `tabular` 标签系统提供 CSV 导入/导出功能。它共享相同的 `tabular.Importer` 和 `tabular.Exporter` 接口，便于在格式之间切换。

## 概述

| 操作 | 入口 | 说明 |
| --- | --- | --- |
| 导入 | `csv.NewImporterFor[T]()` | 将 CSV 文件/读取器解析为类型化结构体切片 |
| 导出 | `csv.NewExporterFor[T]()` | 将结构体切片写入 CSV 文件/缓冲区 |

其他公开构造器：

| 构造器 | 作用 |
| --- | --- |
| `csv.NewImporter(adapter, opts...)` / `csv.NewExporter(adapter, opts...)` | 使用显式 `tabular.RowAdapter` |
| `csv.NewMapImporter(specs, mapOpts, opts...)` / `csv.NewMapExporter(specs, opts...)` | 使用显式 schema 导入/导出 map 行 |
| `csv.NewTypedImporterFor[T]()` / `csv.NewTypedExporterFor[T]()` | 对通用 tabular interface 的 typed wrapper |

## 已审查公开 API 面

本页覆盖 `github.com/coldsmirk/vef-framework-go/csv` 的完整公开面：18 个顶层符号、0 个导出字段、0 个导出方法。Public API fingerprint:
`625d27224a8fbc9542243e3ffabba202710b5feba0b34d2d4e1ca0c43630f978`。

| 符号 | 契约 |
| --- | --- |
| `csv.ExportOption` | CSV exporter 使用的 option function 类型。 |
| `csv.ImportOption` | CSV importer 使用的 option function 类型。 |
| `csv.NewExporter` | 通过显式 `tabular.RowAdapter` 构造 `tabular.Exporter`。 |
| `csv.NewExporterFor` | 使用 `tabular.NewStructAdapterFor[T]()` 构造结构体驱动的 `tabular.Exporter`。 |
| `csv.NewImporter` | 通过显式 `tabular.RowAdapter` 构造 `tabular.Importer`。 |
| `csv.NewImporterFor` | 使用 `tabular.NewStructAdapterFor[T]()` 构造结构体驱动的 `tabular.Importer`。 |
| `csv.NewMapExporter` | 使用 `tabular.NewMapAdapterFromSpecs(specs)` 校验 `[]tabular.ColumnSpec`，成功后返回 map exporter，否则返回校验错误；它不接受 `mapOpts`。 |
| `csv.NewMapImporter` | 使用 `tabular.NewMapAdapterFromSpecs(specs, mapOpts...)` 校验 `[]tabular.ColumnSpec`；不需要 row validator 时，`[]tabular.MapOption` 参数传 `nil`。 |
| `csv.NewTypedExporterFor` | 把结构体 exporter 包装成 `tabular.TypedExporter[T]`，让 export 直接接受 `[]T`。 |
| `csv.NewTypedImporterFor` | 把结构体 importer 包装成 `tabular.TypedImporter[T]`，让 import 直接返回 `[]T`。 |
| `csv.WithCRLF` | 导出时写 Windows 风格 CRLF 行尾。 |
| `csv.WithComment` | 为导入启用注释字符；默认值 `0` 表示不启用注释处理。 |
| `csv.WithExportDelimiter` | 设置导出分隔符；默认是逗号。 |
| `csv.WithImportDelimiter` | 设置导入分隔符；默认是逗号。 |
| `csv.WithSkipRows` | 在 header/data 处理前跳过前导行；负数会归零。 |
| `csv.WithoutHeader` | 把第一个未跳过的行当作数据，并按 schema 顺序做位置映射。 |
| `csv.WithoutTrimSpace` | 关闭框架层 trim；这会影响空行检测、header 匹配和单元格解析。 |
| `csv.WithoutWriteHeader` | 导出时不写 header 行。 |

## 模型定义

CSV 使用与 Excel 相同的 `tabular` 结构体标签。完整标签参考请查看 [Excel 文档](./excel#tabular-标签)。

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

### 基本导出

```go
import "github.com/coldsmirk/vef-framework-go/csv"

exporter := csv.NewExporterFor[Employee]()

// 导出到文件
err := exporter.ExportToFile(employees, "employees.csv")

// 导出到缓冲区（用于 HTTP 响应）
buf, err := exporter.Export(employees)
```

### 导出选项

| 选项 | 默认值 | 说明 |
| --- | --- | --- |
| `csv.WithExportDelimiter(rune)` | `,` | 字段分隔符 |
| `csv.WithoutWriteHeader()` | 写入表头 | 跳过表头行 |
| `csv.WithCRLF()` | 仅 LF | 使用 Windows 风格 CRLF 换行符 |

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

### 自定义格式化器

```go
exporter := csv.NewExporterFor[Employee]()

exporter.RegisterFormatter("status", tabular.FormatterFunc(func(value any) (string, error) {
    if active, ok := value.(bool); ok && active {
        return "是", nil
    }
    return "否", nil
}))
```

## 导入

### 基本导入

```go
importer := csv.NewImporterFor[Employee]()

// 从文件导入
data, importErrors, err := importer.ImportFromFile("employees.csv")
if err != nil {
    return err // 致命错误
}

// 检查行级错误
for _, e := range importErrors {
    log.Printf("第 %d 行：%v", e.Row, e.Err)
}

employees := data.([]Employee)
```

### 从 io.Reader 导入

```go
data, importErrors, err := importer.Import(reader)
```

### 导入选项

| 选项 | 默认值 | 说明 |
| --- | --- | --- |
| `csv.WithImportDelimiter(rune)` | `,` | 字段分隔符 |
| `csv.WithoutHeader()` | 有表头 | CSV 无表头行；按位置映射列 |
| `csv.WithSkipRows(n)` | `0` | 跳过表头前的前导行 |
| `csv.WithoutTrimSpace()` | 启用修剪 | 禁用自动空白修剪 |
| `csv.WithComment(rune)` | 无 | 注释字符（以此开头的行被跳过）|

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

负数 `WithSkipRows` 会归零。`WithoutTrimSpace()` 也会影响空行检测和 header 匹配。
这层 trim 由 VEF 的 `tabular` 层执行；底层 Go 标准库 CSV reader 不启用
`TrimLeadingSpace`。`Import` 会先调用该 reader 的 `ReadAll` 再解析，因此峰值内存会随文件大小和最终结果切片一起增长。reader 使用
`FieldsPerRecord = -1`，因此允许不等长行，缺失的映射单元格交给 tabular adapter
处理。

### 有表头 vs 无表头模式

| 模式 | 列映射方式 |
| --- | --- |
| 有表头（默认）| 表头名称 → `tabular` 标签名 |
| 无表头 | 列位置 → `tabular` 字段顺序 |

使用 `WithoutHeader()` 时，列按位置匹配。使用 `order` 标签属性控制字段排序：

```go
type Record struct {
    Name  string `tabular:"姓名,order=0"`
    Email string `tabular:"邮箱,order=1"`
    Age   int    `tabular:"年龄,order=2"`
}
```

`WithSkipRows` 会在 header 检测前生效。启用表头时，`rows[skipRows]` 是 header 行，
数据从其后一行开始；使用 `csv.WithoutHeader()` 时，第一个未跳过的行会作为数据解析。
Import error 中的行号是基于 1 的 CSV 文件行号，并包含 skip/header 偏移，方便直接对应文本编辑器里看到的行。

### 自定义解析器

```go
importer := csv.NewImporterFor[Employee]()

importer.RegisterParser("date", tabular.ParserFunc(func(cellValue string, targetType reflect.Type) (any, error) {
    return time.Parse("01/02/2006", cellValue)
}))
```

### 验证

导入的记录会自动使用 `validator.Validate(...)` 进行验证，与 Excel 导入器相同。

## 错误处理

| 错误 | 含义 |
| --- | --- |
| `tabular.ErrDataMustBeSlice` | 导出数据必须是切片 |
| `tabular.ErrNoDataRowsFound` | 经 `csv.WithSkipRows` 与可选表头处理后没有数据行 |
| `tabular.ErrDuplicateHeaderName` | 表头中存在重复的非空列名 |
| `tabular.ErrUnsetField` | 结构体字段无法设置 |

CSV 包还公开 `ImportOption` 与 `ExportOption`，它们是上面 import/export
选项背后的 option function 类型。

顶层错误表示致命的读写或结构错误，包括 `ReadAll` 失败、没有数据行、导出 schema
不匹配，以及最终 writer flush 失败（`flush CSV writer: ...`）。解析失败、validator
失败和 adapter commit 失败会聚合进 `[]tabular.ImportError`；import 可以返回
`err == nil` 同时带有非空行级错误，对应行会被跳过，后续行继续处理。

## CSV vs Excel

| 特性 | CSV | Excel |
| --- | --- | --- |
| 文件格式 | `.csv`（纯文本）| `.xlsx`（二进制）|
| 多工作表 | 不支持 | 支持 |
| 列宽度 | 忽略 | 应用 |
| 分隔符 | 可配置 | 不适用 |
| 注释行 | 支持 | 不适用 |
| 空白修剪 | 可配置 | 不适用 |
| 换行符 | LF 或 CRLF | 不适用 |
| 依赖 | Go 标准库 | excelize |

两个包都实现了 `tabular.Importer` / `tabular.Exporter` 接口，因此可以在不更改模型定义的情况下互换使用。
