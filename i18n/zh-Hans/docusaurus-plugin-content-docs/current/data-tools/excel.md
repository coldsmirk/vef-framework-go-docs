---
sidebar_position: 7
---

# Excel

[tabular 导入导出核心](./tabular)的 Excel 后端，为 `.xlsx` 工作簿实现共享的 schema、格式化器与解析器契约。

`excel` 包基于共享的 `tabular` 引擎提供 Excel 导入/导出。底层使用 [excelize](https://github.com/xuri/excelize)，并集成了 VEF 的验证系统。

## 包结构

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

## 工作表

Excel workbook 可以包含多个工作表。导出总是写入单个具名 sheet（`excel.WithSheetName`，默认 `Sheet1`）；导入既可以按名称选择源 sheet（`excel.WithImportSheetName`，设置后优先生效），也可以按 0-based index 选择（`excel.WithImportSheetIndex`，默认 `0`）。

## 模型示例

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

## 导出

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

### Native Typed Cell

Excel 导出在列使用默认 formatter，且没有显式 `format`、`formatter` 或 `FormatterFn`（通过 `tabular.IsDefaultFormatter` 判断）时，会写入 native typed cell。整数、浮点数、布尔值、`time.Time`、`timex.Date` 和 `timex.DateTime` 都以这种方式在 Excel 里保持可排序或可求和。一旦列设置了格式字符串或自定义 formatter，导出器会改为把格式化结果按文本写入。

Native cell 转换的细节：nil pointer 会写成空单元格，非 nil pointer 会先解引用；`timex.Date` / `timex.DateTime` 会被解包成 `time.Time`，让 excelize 存储 native 日期(时间)单元格；`decimal.Decimal` 会转换成 `float64`（在约 15–16 位有效数字内精确，超出则有损——需要完整 decimal 精度的列应改为声明显式 `format`，按精确文本渲染）；`timex.Time` 会被刻意保留为文本，因为它的 zero-date 部分早于 Excel epoch，直接转换会渲染出一个错误的日期。

### 自定义 Formatter 与 Parser

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

### HTTP Handler 中导出

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

## 导入

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

## 验证

导入的记录会自动使用 `validator.Validate(...)` 进行验证。如果验证失败，该行会被添加到 `importErrors` 并从结果切片中跳过。

```go
type Employee struct {
    Name  string `tabular:"Name" validate:"required"`
    Email string `tabular:"Email" validate:"required,email"`
}
```

## 列映射规则

1. 导入器通过 Excel 表头名称 → `tabular` 标签名（无标签时使用字段名）进行匹配。
2. 未匹配的 Excel 列会被静默忽略。
3. 缺失的 Excel 列会让结构体字段保持零值（如指定了 `default` 则使用默认值）。
4. 空行会被自动跳过。

使用 `excel.WithoutHeader()` 时会绕过 header 匹配，改用 `tabular.DefaultPositionalMapping`：源文件第 1 列映射 schema 第 1 列，第 2 列映射 schema 第 2 列，依此类推。

Excel 导入/导出的默认类型支持与 CSV 完全相同——见 tabular 核心页的[默认类型支持](./tabular#默认类型支持)。

## 错误处理

| 错误 | 含义 |
| --- | --- |
| `excel.ErrSheetIndexOutOfRange` | 配置的工作表索引为负数或超出可用范围 |
| `tabular.ErrNoDataRowsFound` | 经 `WithSkipRows` 与可选表头处理后没有数据行 |
| `tabular.ErrDuplicateHeaderName` | 表头中存在重复的非空列名 |
| `tabular.ErrUnsetField` | 结构体字段无法设置，通常是未导出字段 |

顶层 import error 表示致命的文件或工作表错误。解析失败、validator 失败和 adapter commit 失败会聚合进 `[]tabular.ImportError`（`Row`/`Column`/`Field`/`Err`，基于 1 的行号，包含表头行）；对应行会被跳过，后续行继续处理。Export error 使用相同的 `ExportError` 结构，行索引基于 0。

## 下一步

- [表格导入导出](./tabular) — 共享的 schema、标签与接口
- [CSV](./csv) — CSV 后端
