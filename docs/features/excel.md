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

Other public constructors:

| Constructor | Purpose |
| --- | --- |
| `excel.NewImporter(adapter, opts...)` / `excel.NewExporter(adapter, opts...)` | use an explicit `tabular.RowAdapter` |
| `excel.NewMapImporter(specs, mapOpts, opts...)` / `excel.NewMapExporter(specs, opts...)` | import/export map rows with an explicit schema |
| `excel.NewTypedImporterFor[T]()` / `excel.NewTypedExporterFor[T]()` | typed wrappers over the generic tabular interfaces |

## Reviewed Public Surface

This page covers the complete public surface of
`github.com/coldsmirk/vef-framework-go/excel`: 17 top-level symbols, 0 exported
fields, 0 exported methods. Public API fingerprint:
`a449ebeda509ae9b0a2c7bfa083c70b45bc4635bdb49ff1d674400aced129324`.

| Symbol | Contract |
| --- | --- |
| `excel.ErrSheetIndexOutOfRange` | Sentinel returned through the top-level error when the configured sheet index is negative or outside the workbook's sheet list. |
| `excel.ExportOption` | Option function type consumed by Excel exporters. |
| `excel.ImportOption` | Option function type consumed by Excel importers. |
| `excel.NewExporter` | Builds a `tabular.Exporter` from an explicit `tabular.RowAdapter`. |
| `excel.NewExporterFor` | Builds a struct-backed `tabular.Exporter` using `tabular.NewStructAdapterFor[T]()`. |
| `excel.NewImporter` | Builds a `tabular.Importer` from an explicit `tabular.RowAdapter`. |
| `excel.NewImporterFor` | Builds a struct-backed `tabular.Importer` using `tabular.NewStructAdapterFor[T]()`. |
| `excel.NewMapExporter` | Validates `[]tabular.ColumnSpec` with `tabular.NewMapAdapterFromSpecs(specs)` and returns a map-backed exporter or the validation error. |
| `excel.NewMapImporter` | Validates `[]tabular.ColumnSpec` with `tabular.NewMapAdapterFromSpecs(specs, mapOpts...)`; pass `nil` for the `[]tabular.MapOption` argument when no row validators are needed. |
| `excel.NewTypedExporterFor` | Wraps the struct exporter in `tabular.TypedExporter[T]` so export accepts `[]T`. |
| `excel.NewTypedImporterFor` | Wraps the struct importer in `tabular.TypedImporter[T]` so import returns `[]T`. |
| `excel.WithImportSheetIndex` | Selects the worksheet by 0-based index; ignored when `excel.WithImportSheetName` is set. |
| `excel.WithImportSheetName` | Selects the worksheet by name and takes precedence over index selection. |
| `excel.WithSheetName` | Sets the export worksheet name; default is `Sheet1`, and export renames the default sheet rather than creating a second sheet. |
| `excel.WithSkipRows` | Skips leading rows before header/data processing; negative values are clamped to `0`. |
| `excel.WithoutHeader` | Treats the first non-skipped row as data and maps columns positionally in schema order. |
| `excel.WithoutTrimSpace` | Disables default trimming; this affects empty-row detection, header matching, and cell parsing. |

## `tabular` Tag

The `tabular` struct tag controls how model fields map to Excel columns.

### Tag Attributes

| Attribute | Meaning | Example |
| --- | --- | --- |
| (default value) | Column header name | `tabular:"Username"` |
| `name` | Explicit column name | `tabular:"name=Username"` |
| `width` | Column width hint | `tabular:"width=20"` |
| `order` | Column order (0-based) | `tabular:"order=1"` |
| `default` | Default value for empty cells on import | `tabular:"default=N/A"` |
| `format` | Format template (date/number) | `tabular:"format=2006-01-02"` |
| `formatter` | Custom formatter name for export | `tabular:"formatter=status"` |
| `parser` | Custom parser name for import | `tabular:"parser=date"` |
| `dive` | Recurse into embedded struct | `tabular:"dive"` |
| `-` | Ignore this field | `tabular:"-"` |

### Model Example

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

Excel export writes native typed cells when a column uses the default formatter
and has no explicit `format`, `formatter`, or `FormatterFn`. Numeric, boolean,
`time.Time`, `timex.Date`, and `timex.DateTime` values remain sortable or
summable in Excel. Once a column sets a format string or custom formatter, the
exporter writes the formatted result as text. Nil pointers become empty cells,
non-nil pointers are dereferenced, and `timex.Time` is intentionally left as
text because its zero-date component predates the Excel epoch.

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
IsActive bool `tabular:"Status,formatter=status"`
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

// Treat the first non-skipped row as data, mapped by schema position
importer := excel.NewImporterFor[Employee](
    excel.WithoutHeader(),
)

// Preserve leading/trailing spaces in headers and cells
importer := excel.NewImporterFor[Employee](
    excel.WithoutTrimSpace(),
)
```

`WithImportSheetName` takes precedence over `WithImportSheetIndex`. Negative
`WithSkipRows` values are clamped to `0`. With headers enabled, skipped rows are
ignored before header resolution; without headers, the first non-skipped row is
parsed as data with positional mapping. Like CSV, Excel import trims cell values
by default and reads the workbook rows into memory before parsing.

### Custom Parser

Register a parser to customize how cell values are parsed:

```go
importer := excel.NewImporterFor[Employee]()

importer.RegisterParser("date", tabular.ParserFunc(func(cellValue string, targetType reflect.Type) (any, error) {
    return time.Parse("01/02/2006", cellValue)
}))
```

Then reference it in the struct tag:

```go
JoinDate time.Time `tabular:"Join Date,parser=date"`
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
| `excel.ErrSheetIndexOutOfRange` | Configured sheet index is negative or exceeds available sheets |
| `tabular.ErrNoDataRowsFound` | File has no data rows after `WithSkipRows` and optional header handling |
| `tabular.ErrDuplicateHeaderName` | Duplicate non-empty column names in the header row |
| `tabular.ErrUnsetField` | Struct field cannot be set, typically because it is unexported |

The Excel package exposes `ImportOption` and `ExportOption` as the option
function types behind `WithSheetName`, `WithImportSheetName`,
`WithImportSheetIndex`, `WithSkipRows`, `WithoutHeader`, and
`WithoutTrimSpace`.

Top-level import errors are fatal file or worksheet failures. Parse failures,
validator failures, and adapter commit failures are collected as
`[]tabular.ImportError`; affected rows are skipped while later rows continue.

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

When `excel.WithoutHeader()` is used, header matching is bypassed and the
importer uses `tabular.DefaultPositionalMapping`: the first source column maps
to the first schema column, the second source column maps to the second schema
column, and so on.

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
