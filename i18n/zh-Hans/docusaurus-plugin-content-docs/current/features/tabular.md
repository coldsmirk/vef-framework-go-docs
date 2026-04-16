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

    Name       string          `tabular:"姓名,width:20"`
    Email      string          `tabular:"邮箱,width:30"`
    Department string          `tabular:"name:部门,order:2,width:15"`
    JoinDate   timex.Date      `tabular:"入职日期,format:2006-01-02,width:15"`
    Salary     decimal.Decimal `tabular:"薪资,width:12,format:#,##0.00"`
    Status     string          `tabular:"状态,default:active,formatter:status"`
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

## Schema

`Schema` 类型在初始化时从结构体字段预解析表格元数据：

```go
schema := tabular.NewSchemaFor[Employee]()

columns := schema.Columns()       // []*Column — 所有解析出的列
names := schema.ColumnNames()     // []string{"姓名", "邮箱", ...}
count := schema.ColumnCount()     // 6
```

列会按 `order` 属性自动排序。未显式指定 `order` 的字段使用其声明顺序。

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
tabular.ValueParserFunc(func(cellValue string, targetType reflect.Type) (any, error) { ... })
```

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
