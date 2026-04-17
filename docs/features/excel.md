---
sidebar_position: 4
---

# Excel

The `excel` package provides Excel import/export capabilities based on the `tabular` tag system. It uses [excelize](https://github.com/xuri/excelize) under the hood and integrates with VEF's validation system.

## Overview

| Operation | Entry Point | Description |
| --- | --- | --- |
| Import | `excel.NewImporterFor[T]()` | Parse Excel file into typed struct slices |
| Export | `excel.NewExporterFor[T]()` | Write struct slices to Excel files |

Both importer and exporter use `tabular` struct tags to define column mapping, formatting, and parsing rules.

## `tabular` Tag

The `tabular` struct tag controls how model fields map to Excel columns.

### Tag Attributes

| Attribute | Meaning | Example |
| --- | --- | --- |
| (default value) | Column header name | `tabular:"Username"` |
| `name` | Explicit column name | `tabular:"name:Username"` |
| `width` | Column width hint | `tabular:"width:20"` |
| `order` | Column order (0-based) | `tabular:"order:1"` |
| `default` | Default value for empty cells on import | `tabular:"default:N/A"` |
| `format` | Format template (date/number) | `tabular:"format:2006-01-02"` |
| `formatter` | Custom formatter name for export | `tabular:"formatter:status"` |
| `parser` | Custom parser name for import | `tabular:"parser:date"` |
| `dive` | Recurse into embedded struct | `tabular:"dive"` |
| `-` | Ignore this field | `tabular:"-"` |

### Model Example

```go
type Employee struct {
    orm.FullAuditedModel `tabular:"-"`

    Name       string          `json:"name" bun:"name" tabular:"姓名,width:20"`
    Email      string          `json:"email" bun:"email" tabular:"邮箱,width:30"`
    Department string          `json:"department" bun:"department" tabular:"部门,width:15"`
    JoinDate   timex.Date      `json:"joinDate" bun:"join_date" tabular:"入职日期,format:2006-01-02,width:15"`
    Salary     decimal.Decimal `json:"salary" bun:"salary" tabular:"薪资,width:12"`
    IsActive   bool            `json:"isActive" bun:"is_active" tabular:"是否在职,width:10"`
}
```

## Exporting

### Basic Export

```go
import "github.com/coldsmirk/vef-framework-go/excel"

// Create a typed exporter
exporter := excel.NewExporterFor[Employee]()

// Export to file
err := exporter.ExportToFile(employees, "employees.xlsx")

// Export to buffer (for HTTP response)
buf, err := exporter.Export(employees)
```

### Export Options

```go
// Custom sheet name (default: "Sheet1")
exporter := excel.NewExporterFor[Employee](
    excel.WithSheetName("Employees"),
)
```

### Custom Formatter

Register a formatter to customize how values are rendered in Excel:

```go
exporter := excel.NewExporterFor[Employee]()

// Register a custom formatter for the "status" name
exporter.RegisterFormatter("status", tabular.FormatterFunc(func(value any) (string, error) {
    if active, ok := value.(bool); ok && active {
        return "Active", nil
    }
    return "Inactive", nil
}))
```

Then reference it in the struct tag:

```go
IsActive bool `tabular:"Status,formatter:status"`
```

### Export in HTTP Handler

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

## Importing

### Basic Import

```go
// Create a typed importer
importer := excel.NewImporterFor[Employee]()

// Import from file
data, importErrors, err := importer.ImportFromFile("employees.xlsx")
if err != nil {
    // Fatal error (file not found, etc.)
    return err
}

// Check row-level errors
if len(importErrors) > 0 {
    for _, e := range importErrors {
        log.Printf("Row %d, Column %s: %v", e.Row, e.Column, e.Err)
    }
}

// Type-assert the result
employees := data.([]Employee)
```

### Import from io.Reader

```go
// Import from an uploaded file (io.Reader)
data, importErrors, err := importer.Import(reader)
```

### Import Options

```go
// Specify sheet by name
importer := excel.NewImporterFor[Employee](
    excel.WithImportSheetName("Staff"),
)

// Specify sheet by index (default: 0)
importer := excel.NewImporterFor[Employee](
    excel.WithImportSheetIndex(1),
)

// Skip leading rows (e.g., title rows before the header)
importer := excel.NewImporterFor[Employee](
    excel.WithSkipRows(2),
)
```

### Custom Parser

Register a parser to customize how cell values are parsed:

```go
importer := excel.NewImporterFor[Employee]()

importer.RegisterParser("date", tabular.ValueParserFunc(func(cellValue string, targetType reflect.Type) (any, error) {
    return time.Parse("01/02/2006", cellValue)
}))
```

Then reference it in the struct tag:

```go
JoinDate time.Time `tabular:"Join Date,parser:date"`
```

### Validation

Imported records are automatically validated using `validator.Validate(...)`.  If validation fails, the row is added to `importErrors` and skipped from the result slice.

```go
type Employee struct {
    Name  string `tabular:"Name" validate:"required"`
    Email string `tabular:"Email" validate:"required,email"`
}
```

## Error Handling

### Import Errors

| Error | Meaning |
| --- | --- |
| `ErrSheetIndexOutOfRange` | Sheet index exceeds available sheets |
| `ErrNoDataRowsFound` | File has no data rows (only header or empty) |
| `ErrDuplicateColumnName` | Duplicate column names in header row |
| `ErrFieldNotSettable` | Struct field cannot be set (unexported) |

Row-level errors are returned as `[]tabular.ImportError` (non-fatal):

```go
type ImportError struct {
    Row    int    // 1-based row number (including header)
    Column string // Column header name
    Field  string // Struct field name
    Err    error  // Underlying error
}
```

### Export Errors

```go
type ExportError struct {
    Row    int    // 0-based data row index
    Column string
    Field  string
    Err    error
}
```

## Column Mapping Rules

1. The importer matches Excel header names → `tabular` tag name (or field name if no tag)
2. Unmatched Excel columns are silently ignored
3. Missing Excel columns leave the struct field at its zero value (or `default` if specified)
4. Empty rows are automatically skipped

## Default Type Support

The default parser automatically handles:

| Go Type | Parsing |
| --- | --- |
| `string` | Direct assignment |
| `int`, `int8`–`int64` | Integer parsing |
| `uint`, `uint8`–`uint64` | Unsigned integer parsing |
| `float32`, `float64` | Float parsing |
| `bool` | Boolean parsing (`true`/`false`, `1`/`0`) |
| `decimal.Decimal` | Decimal string parsing |
| `timex.Date` / `timex.DateTime` | Date/time parsing using `format` attribute |
| `*T` (pointer types) | Nil for empty cells, parsed value otherwise |
