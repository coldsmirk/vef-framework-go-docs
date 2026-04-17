---
sidebar_position: 5
---

# CSV

The `csv` package provides CSV import/export capabilities using the same `tabular` tag system as the [Excel](./excel) package. It shares the same `tabular.Importer` and `tabular.Exporter` interfaces, making it easy to swap between formats.

## Overview

| Operation | Entry Point | Description |
| --- | --- | --- |
| Import | `csv.NewImporterFor[T]()` | Parse CSV file/reader into typed struct slices |
| Export | `csv.NewExporterFor[T]()` | Write struct slices to CSV files/buffers |

## Model Definition

CSV uses the same `tabular` struct tag as Excel. See the [Excel documentation](./excel#tabular-tag) for the full tag reference.

```go
type Employee struct {
    orm.FullAuditedModel `tabular:"-"`

    Name       string          `tabular:"Name,width:20"`
    Email      string          `tabular:"Email,width:30"`
    Department string          `tabular:"Department"`
    JoinDate   timex.Date      `tabular:"Join Date,format:2006-01-02"`
    Salary     decimal.Decimal `tabular:"Salary"`
    IsActive   bool            `tabular:"Active"`
}
```

## Exporting

### Basic Export

```go
import "github.com/coldsmirk/vef-framework-go/csv"

exporter := csv.NewExporterFor[Employee]()

// Export to file
err := exporter.ExportToFile(employees, "employees.csv")

// Export to buffer (for HTTP response)
buf, err := exporter.Export(employees)
```

### Export Options

| Option | Default | Description |
| --- | --- | --- |
| `csv.WithExportDelimiter(rune)` | `,` | Field delimiter character |
| `csv.WithoutWriteHeader()` | write header | Skip the header row |
| `csv.WithCrlf()` | LF only | Use Windows-style CRLF line endings |

```go
// TSV export with Windows line endings
exporter := csv.NewExporterFor[Employee](
    csv.WithExportDelimiter('\t'),
    csv.WithCrlf(),
)

// Export without header row
exporter := csv.NewExporterFor[Employee](
    csv.WithoutWriteHeader(),
)
```

### Custom Formatter

```go
exporter := csv.NewExporterFor[Employee]()

exporter.RegisterFormatter("status", tabular.FormatterFunc(func(value any) (string, error) {
    if active, ok := value.(bool); ok && active {
        return "Y", nil
    }
    return "N", nil
}))
```

## Importing

### Basic Import

```go
importer := csv.NewImporterFor[Employee]()

// Import from file
data, importErrors, err := importer.ImportFromFile("employees.csv")
if err != nil {
    return err // Fatal error
}

// Check row-level errors
for _, e := range importErrors {
    log.Printf("Row %d: %v", e.Row, e.Err)
}

employees := data.([]Employee)
```

### Import from io.Reader

```go
data, importErrors, err := importer.Import(reader)
```

### Import Options

| Option | Default | Description |
| --- | --- | --- |
| `csv.WithImportDelimiter(rune)` | `,` | Field delimiter character |
| `csv.WithoutHeader()` | has header | CSV has no header row; map columns by position |
| `csv.WithSkipRows(n)` | `0` | Skip leading rows before the header |
| `csv.WithoutTrimSpace()` | trim enabled | Disable automatic whitespace trimming |
| `csv.WithComment(rune)` | none | Comment character (lines starting with this are skipped) |

```go
// TSV file with 2 title rows, comment lines starting with #
importer := csv.NewImporterFor[Employee](
    csv.WithImportDelimiter('\t'),
    csv.WithSkipRows(2),
    csv.WithComment('#'),
)

// CSV without header (columns matched by position/order)
importer := csv.NewImporterFor[Employee](
    csv.WithoutHeader(),
)
```

### Header vs No-Header Mode

| Mode | Column Mapping |
| --- | --- |
| With header (default) | Header names → `tabular` tag names |
| Without header | Column position → `tabular` field order |

When using `WithoutHeader()`, columns are matched by position. Use the `order` tag attribute to control field ordering:

```go
type Record struct {
    Name  string `tabular:"Name,order:0"`
    Email string `tabular:"Email,order:1"`
    Age   int    `tabular:"Age,order:2"`
}
```

### Custom Parser

```go
importer := csv.NewImporterFor[Employee]()

importer.RegisterParser("date", tabular.ValueParserFunc(func(cellValue string, targetType reflect.Type) (any, error) {
    return time.Parse("01/02/2006", cellValue)
}))
```

### Validation

Imported records are automatically validated using `validator.Validate(...)`, same as the Excel importer.

## Error Handling

| Error | Meaning |
| --- | --- |
| `ErrDataMustBeSlice` | Export data must be a slice |
| `ErrNoDataRowsFound` | No data rows in the file |
| `ErrDuplicateColumnName` | Duplicate column names in header |
| `ErrFieldNotSettable` | Struct field cannot be set |

## CSV vs Excel

| Feature | CSV | Excel |
| --- | --- | --- |
| File format | `.csv` (plain text) | `.xlsx` (binary) |
| Multiple sheets | No | Yes |
| Column widths | Ignored | Applied |
| Delimiter | Configurable | N/A |
| Comment lines | Supported | N/A |
| Trim whitespace | Configurable | N/A |
| Line endings | LF or CRLF | N/A |
| Dependencies | Go stdlib | excelize |

Both packages implement `tabular.Importer` / `tabular.Exporter`, so you can swap between them without changing your model definitions.
