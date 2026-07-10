---
sidebar_position: 7
---

# Excel

The Excel backend of the [tabular import/export core](./tabular). It implements the shared schema, formatter, and parser contracts for `.xlsx` workbooks.

The `excel` package provides Excel import/export using the shared `tabular` engine. It uses [excelize](https://github.com/xuri/excelize) under the hood and integrates with VEF's validation system.

## Package Surface

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

The Excel option marker types are `excel.ExportOption` and `excel.ImportOption`. `excel.NewMapExporter` validates `specs` with `tabular.NewMapAdapterFromSpecs(specs)`; `excel.NewMapImporter` validates with `tabular.NewMapAdapterFromSpecs(specs, mapOpts...)` — pass `nil` for `mapOpts` when no row validators are needed.

| Option | Default | Purpose |
| --- | --- | --- |
| `excel.WithSheetName(name)` | `Sheet1` | Worksheet name on export; renames the default sheet rather than creating a second one |
| `excel.WithImportSheetName(name)` | none | Read a worksheet by name; takes precedence over `WithImportSheetIndex` |
| `excel.WithImportSheetIndex(i)` | `0` | Read a worksheet by 0-based index (returns `excel.ErrSheetIndexOutOfRange` if negative or out of range); ignored when `WithImportSheetName` is set |
| `excel.WithSkipRows(n)` | `0` | Skip the first `n` rows before reading; negative values are clamped to `0` |
| `excel.WithoutHeader()` | header on | First non-skipped row is data; positional mapping |
| `excel.WithoutTrimSpace()` | trim on | Disable cell trimming (also affects empty-row detection and header matching) |

`Column.Width` set on a `ColumnSpec` (or via the struct tag `width=…`) is applied to the generated worksheet's column width; CSV ignores it.

## Sheets

Excel workbooks can hold multiple sheets. Export always writes to a single named sheet (`excel.WithSheetName`, default `Sheet1`); import selects a source sheet either by name (`excel.WithImportSheetName`, which wins if set) or by 0-based index (`excel.WithImportSheetIndex`, default `0`).

## Model Example

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

```go
import "github.com/coldsmirk/vef-framework-go/excel"

// Create a typed exporter
exporter := excel.NewExporterFor[Employee]()

// Export to file
err := exporter.ExportToFile(employees, "employees.xlsx")

// Export to buffer (for HTTP response)
buf, err := exporter.Export(employees)
```

```go
// Custom sheet name (default: "Sheet1")
exporter := excel.NewExporterFor[Employee](
    excel.WithSheetName("Employees"),
)
```

### Native Typed Cells

Excel export writes native typed cells when a column uses the default formatter and declares no explicit `format`, `formatter`, or `FormatterFn` (checked via `tabular.IsDefaultFormatter`). Integers, floats, booleans, `time.Time`, `timex.Date`, and `timex.DateTime` remain sortable or summable in Excel this way. Once a column sets a format string or a custom formatter, the exporter writes the formatted result as text instead.

Details of the native-cell conversion: nil pointers become empty cells and non-nil pointers are dereferenced; `timex.Date` / `timex.DateTime` are unwrapped to `time.Time` so excelize stores a native date(time) cell; `decimal.Decimal` is converted to `float64` (exact within ~15–16 significant digits, lossy beyond that — a column needing full decimal precision should declare an explicit `format` to render exact text instead); `timex.Time` is intentionally left as text because its zero-date component predates the Excel epoch and would otherwise render a bogus date.

### Custom Formatter and Parser

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

// Or import from an uploaded file (io.Reader)
data, importErrors, err = importer.Import(reader)
```

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

`WithImportSheetName` takes precedence over `WithImportSheetIndex`. Negative `WithSkipRows` values are clamped to `0`. With headers enabled, skipped rows are ignored before header resolution; without headers, the first non-skipped row is parsed as data with positional mapping. Like CSV, Excel import trims cell values by default and reads the workbook rows into memory before parsing — excelize loads the whole workbook up front, so peak memory scales with the file size in addition to the materialized result slice.

## Validation

Imported records are automatically validated using `validator.Validate(...)`. If validation fails, the row is added to `importErrors` and skipped from the result slice.

```go
type Employee struct {
    Name  string `tabular:"Name" validate:"required"`
    Email string `tabular:"Email" validate:"required,email"`
}
```

## Column Mapping Rules

1. The importer matches Excel header names → `tabular` tag name (or field name if no tag).
2. Unmatched Excel columns are silently ignored.
3. Missing Excel columns leave the struct field at its zero value (or `default` if specified).
4. Empty rows are automatically skipped.

When `excel.WithoutHeader()` is used, header matching is bypassed and the importer uses `tabular.DefaultPositionalMapping`: the first source column maps to the first schema column, the second source column maps to the second schema column, and so on.

Default type support for Excel import/export is identical to CSV — see [Default Type Support](./tabular#default-type-support) in the tabular core.

## Error Handling

| Error | Meaning |
| --- | --- |
| `excel.ErrSheetIndexOutOfRange` | Configured sheet index is negative or exceeds available sheets |
| `tabular.ErrNoDataRowsFound` | File has no data rows after `WithSkipRows` and optional header handling |
| `tabular.ErrDuplicateHeaderName` | Duplicate non-empty column names in the header row |
| `tabular.ErrUnsetField` | Struct field cannot be set, typically because it is unexported |

Top-level import errors are fatal file or worksheet failures. Parse failures, validator failures, and adapter commit failures are collected as `[]tabular.ImportError` (`Row`/`Column`/`Field`/`Err`, 1-based row including the header); affected rows are skipped while later rows continue. Export errors use the same `ExportError` shape with a 0-based data row index.

## Next Step

- [Tabular Import & Export](./tabular) — the shared schema, tags, and interfaces
- [CSV](./csv) — the CSV backend
