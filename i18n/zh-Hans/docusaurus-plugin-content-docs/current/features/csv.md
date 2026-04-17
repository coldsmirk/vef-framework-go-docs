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

## 模型定义

CSV 使用与 Excel 相同的 `tabular` 结构体标签。完整标签参考请查看 [Excel 文档](./excel#tabular-标签)。

```go
type Employee struct {
    orm.FullAuditedModel `tabular:"-"`

    Name       string          `tabular:"姓名,width:20"`
    Email      string          `tabular:"邮箱,width:30"`
    Department string          `tabular:"部门"`
    JoinDate   timex.Date      `tabular:"入职日期,format:2006-01-02"`
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
| `csv.WithCrlf()` | 仅 LF | 使用 Windows 风格 CRLF 换行符 |

```go
// TSV 导出，使用 Windows 换行符
exporter := csv.NewExporterFor[Employee](
    csv.WithExportDelimiter('\t'),
    csv.WithCrlf(),
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

### 有表头 vs 无表头模式

| 模式 | 列映射方式 |
| --- | --- |
| 有表头（默认）| 表头名称 → `tabular` 标签名 |
| 无表头 | 列位置 → `tabular` 字段顺序 |

使用 `WithoutHeader()` 时，列按位置匹配。使用 `order` 标签属性控制字段排序：

```go
type Record struct {
    Name  string `tabular:"姓名,order:0"`
    Email string `tabular:"邮箱,order:1"`
    Age   int    `tabular:"年龄,order:2"`
}
```

### 自定义解析器

```go
importer := csv.NewImporterFor[Employee]()

importer.RegisterParser("date", tabular.ValueParserFunc(func(cellValue string, targetType reflect.Type) (any, error) {
    return time.Parse("01/02/2006", cellValue)
}))
```

### 验证

导入的记录会自动使用 `validator.Validate(...)` 进行验证，与 Excel 导入器相同。

## 错误处理

| 错误 | 含义 |
| --- | --- |
| `ErrDataMustBeSlice` | 导出数据必须是切片 |
| `ErrNoDataRowsFound` | 文件中没有数据行 |
| `ErrDuplicateColumnName` | 表头中存在重复列名 |
| `ErrFieldNotSettable` | 结构体字段无法设置 |

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
