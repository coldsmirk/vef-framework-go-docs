---
sidebar_position: 4
---

# Excel

`excel` 包基于 `tabular` 标签系统提供 Excel 导入/导出功能。底层使用 [excelize](https://github.com/xuri/excelize)，并集成了 VEF 的验证系统。

## 概述

| 操作 | 入口 | 说明 |
| --- | --- | --- |
| 导入 | `excel.NewImporterFor[T]()` | 将 Excel 文件解析为类型化结构体切片 |
| 导出 | `excel.NewExporterFor[T]()` | 将结构体切片写入 Excel 文件 |

导入器和导出器都使用 `tabular` 结构体标签来定义列映射、格式化和解析规则。

其他公开构造器：

| 构造器 | 作用 |
| --- | --- |
| `excel.NewImporter(adapter, opts...)` / `excel.NewExporter(adapter, opts...)` | 使用显式 `tabular.RowAdapter` |
| `excel.NewMapImporter(specs, mapOpts, opts...)` / `excel.NewMapExporter(specs, opts...)` | 使用显式 schema 导入/导出 map 行 |
| `excel.NewTypedImporterFor[T]()` / `excel.NewTypedExporterFor[T]()` | 对通用 tabular interface 的 typed wrapper |

## 已审查公开 API 面

本页覆盖 `github.com/coldsmirk/vef-framework-go/excel` 的完整公开面：17 个顶层符号、0 个导出字段、0 个导出方法。Public API fingerprint:
`a449ebeda509ae9b0a2c7bfa083c70b45bc4635bdb49ff1d674400aced129324`。

| 符号 | 契约 |
| --- | --- |
| `excel.ErrSheetIndexOutOfRange` | 配置的工作表索引为负数或超出 workbook sheet 列表时，通过顶层错误返回的 sentinel。 |
| `excel.ExportOption` | Excel exporter 使用的 option function 类型。 |
| `excel.ImportOption` | Excel importer 使用的 option function 类型。 |
| `excel.NewExporter` | 通过显式 `tabular.RowAdapter` 构造 `tabular.Exporter`。 |
| `excel.NewExporterFor` | 使用 `tabular.NewStructAdapterFor[T]()` 构造结构体驱动的 `tabular.Exporter`。 |
| `excel.NewImporter` | 通过显式 `tabular.RowAdapter` 构造 `tabular.Importer`。 |
| `excel.NewImporterFor` | 使用 `tabular.NewStructAdapterFor[T]()` 构造结构体驱动的 `tabular.Importer`。 |
| `excel.NewMapExporter` | 使用 `tabular.NewMapAdapterFromSpecs(specs)` 校验 `[]tabular.ColumnSpec`，成功后返回 map exporter，否则返回校验错误。 |
| `excel.NewMapImporter` | 使用 `tabular.NewMapAdapterFromSpecs(specs, mapOpts...)` 校验 `[]tabular.ColumnSpec`；不需要 row validator 时，`[]tabular.MapOption` 参数传 `nil`。 |
| `excel.NewTypedExporterFor` | 把结构体 exporter 包装成 `tabular.TypedExporter[T]`，让 export 直接接受 `[]T`。 |
| `excel.NewTypedImporterFor` | 把结构体 importer 包装成 `tabular.TypedImporter[T]`，让 import 直接返回 `[]T`。 |
| `excel.WithImportSheetIndex` | 按 0-based index 选择工作表；设置了 `excel.WithImportSheetName` 时会被忽略。 |
| `excel.WithImportSheetName` | 按名称选择工作表，并优先于 index 选择。 |
| `excel.WithSheetName` | 设置导出的工作表名；默认是 `Sheet1`，导出时会重命名默认 sheet，而不是创建第二个 sheet。 |
| `excel.WithSkipRows` | 在 header/data 处理前跳过前导行；负数会归零。 |
| `excel.WithoutHeader` | 把第一个未跳过的行当作数据，并按 schema 顺序做位置映射。 |
| `excel.WithoutTrimSpace` | 关闭默认 trim；这会影响空行检测、header 匹配和单元格解析。 |

## `tabular` 标签

`tabular` 结构体标签控制模型字段如何映射到 Excel 列。

### 标签属性

| 属性 | 含义 | 示例 |
| --- | --- | --- |
| （默认值） | 列标题名称 | `tabular:"用户名"` |
| `name` | 显式列名 | `tabular:"name=用户名"` |
| `width` | 列宽度提示 | `tabular:"width=20"` |
| `order` | 列顺序（从 0 开始） | `tabular:"order=1"` |
| `default` | 导入时空单元格的默认值 | `tabular:"default=N/A"` |
| `format` | 格式化模板（日期/数字） | `tabular:"format=2006-01-02"` |
| `formatter` | 导出时的自定义格式化器名 | `tabular:"formatter=status"` |
| `parser` | 导入时的自定义解析器名 | `tabular:"parser=date"` |
| `dive` | 递归进入嵌入结构体 | `tabular:"dive"` |
| `-` | 忽略该字段 | `tabular:"-"` |

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

## 导出

### 基本导出

```go
import "github.com/coldsmirk/vef-framework-go/excel"

// 创建类型化的导出器
exporter := excel.NewExporterFor[Employee]()

// 导出到文件
err := exporter.ExportToFile(employees, "employees.xlsx")

// 导出到缓冲区（用于 HTTP 响应）
buf, err := exporter.Export(employees)
```

### 导出选项

```go
// 自定义工作表名（默认："Sheet1"）
exporter := excel.NewExporterFor[Employee](
    excel.WithSheetName("员工列表"),
)
```

Excel 导出在列使用默认 formatter，且没有显式 `format`、`formatter` 或
`FormatterFn` 时，会写入 native typed cell。数字、布尔值、`time.Time`、
`timex.Date` 和 `timex.DateTime` 在 Excel 里仍可排序或求和。一旦列设置了格式
字符串或自定义 formatter，导出器会把格式化结果按文本写入。nil pointer 会写成
空单元格，非 nil pointer 会先解引用；`timex.Time` 会被刻意保留为文本，因为它的
zero-date 部分早于 Excel epoch。

### 自定义格式化器

注册格式化器来自定义值在 Excel 中的显示方式：

```go
exporter := excel.NewExporterFor[Employee]()

// 注册名为 "status" 的自定义格式化器
exporter.RegisterFormatter("status", tabular.FormatterFunc(func(value any) (string, error) {
    if active, ok := value.(bool); ok && active {
        return "在职", nil
    }
    return "离职", nil
}))
```

然后在结构体标签中引用：

```go
IsActive bool `tabular:"状态,formatter=status"`
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

### 基本导入

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
        log.Printf("第 %d 行，列 %s：%v", e.Row, e.Column, e.Err)
    }
}

// 类型断言获取结果
employees := data.([]Employee)
```

### 从 io.Reader 导入

```go
// 从上传文件导入（io.Reader）
data, importErrors, err := importer.Import(reader)
```

### 导入选项

```go
// 按工作表名指定
importer := excel.NewImporterFor[Employee](
    excel.WithImportSheetName("员工"),
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

`WithImportSheetName` 优先于 `WithImportSheetIndex`。负数 `WithSkipRows` 会归零。
启用表头时，跳过行会在 header 解析前生效；关闭表头时，第一个未跳过的行会按位置映射作为数据解析。和 CSV 一样，Excel 导入默认会 trim 单元格值，并在解析前把 workbook rows 读入内存。

### 自定义解析器

注册解析器来自定义单元格值的解析方式：

```go
importer := excel.NewImporterFor[Employee]()

importer.RegisterParser("date", tabular.ParserFunc(func(cellValue string, targetType reflect.Type) (any, error) {
    return time.Parse("01/02/2006", cellValue)
}))
```

然后在结构体标签中引用：

```go
JoinDate time.Time `tabular:"入职日期,parser=date"`
```

### 验证

导入的记录会自动使用 `validator.Validate(...)` 进行验证。如果验证失败，该行会被添加到 `importErrors` 并从结果切片中跳过。

```go
type Employee struct {
    Name  string `tabular:"姓名" validate:"required"`
    Email string `tabular:"邮箱" validate:"required,email"`
}
```

## 错误处理

### 导入错误

| 错误 | 含义 |
| --- | --- |
| `excel.ErrSheetIndexOutOfRange` | 配置的工作表索引为负数或超出可用范围 |
| `tabular.ErrNoDataRowsFound` | 经 `WithSkipRows` 与可选表头处理后没有数据行 |
| `tabular.ErrDuplicateHeaderName` | 表头中存在重复的非空列名 |
| `tabular.ErrUnsetField` | 结构体字段无法设置，通常是未导出字段 |

Excel 包还公开 `ImportOption` 和 `ExportOption`，它们是 `WithSheetName`、
`WithImportSheetName`、`WithImportSheetIndex`、`WithSkipRows`、
`WithoutHeader`、`WithoutTrimSpace` 背后的 option function 类型。

顶层 import error 表示致命的文件或工作表错误。解析失败、validator 失败和 adapter
commit 失败会聚合进 `[]tabular.ImportError`；对应行会被跳过，后续行继续处理。

行级错误以 `[]tabular.ImportError` 返回（非致命）：

```go
type ImportError struct {
    Row    int    // 基于 1 的行号（包含表头行）
    Column string // 列标题名称
    Field  string // 结构体字段名
    Err    error  // 底层错误
}
```

### 导出错误

```go
type ExportError struct {
    Row    int    // 基于 0 的数据行索引
    Column string
    Field  string
    Err    error
}
```

## 列映射规则

1. 导入器通过 Excel 表头名称 → `tabular` 标签名（无标签时使用字段名）进行匹配
2. 未匹配的 Excel 列会被静默忽略
3. 缺失的 Excel 列会让结构体字段保持零值（如指定了 `default` 则使用默认值）
4. 空行会被自动跳过

使用 `excel.WithoutHeader()` 时会绕过 header 匹配，改用
`tabular.DefaultPositionalMapping`：源文件第 1 列映射 schema 第 1 列，第 2 列映射
schema 第 2 列，依此类推。

## 默认类型支持

默认解析器自动处理以下类型：

| Go 类型 | 解析方式 |
| --- | --- |
| `string` | 直接赋值 |
| `int`、`int8`–`int64` | 整数解析 |
| `uint`、`uint8`–`uint64` | 无符号整数解析 |
| `float32`、`float64` | 浮点数解析 |
| `bool` | 布尔解析（`true`/`false`、`1`/`0`）|
| `decimal.Decimal` | Decimal 字符串解析 |
| `timex.Date` / `timex.DateTime` | 使用 `format` 属性进行日期/时间解析 |
| `*T`（指针类型） | 空单元格为 nil，否则解析值 |
