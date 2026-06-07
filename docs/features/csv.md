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

Other public constructors:

| Constructor | Purpose |
| --- | --- |
| `csv.NewImporter(adapter, opts...)` / `csv.NewExporter(adapter, opts...)` | use an explicit `tabular.RowAdapter` |
| `csv.NewMapImporter(specs, mapOpts, opts...)` / `csv.NewMapExporter(specs, opts...)` | import/export map rows with an explicit schema |
| `csv.NewTypedImporterFor[T]()` / `csv.NewTypedExporterFor[T]()` | typed wrappers over the generic tabular interfaces |

## Reviewed Public Surface

This page covers the complete public surface of
`github.com/coldsmirk/vef-framework-go/csv`: 18 top-level symbols, 0 exported
fields, 0 exported methods. Public API fingerprint:
`625d27224a8fbc9542243e3ffabba202710b5feba0b34d2d4e1ca0c43630f978`.

| Symbol | Contract |
| --- | --- |
| `csv.ExportOption` | Option function type consumed by CSV exporters. |
| `csv.ImportOption` | Option function type consumed by CSV importers. |
| `csv.NewExporter` | Builds a `tabular.Exporter` from an explicit `tabular.RowAdapter`. |
| `csv.NewExporterFor` | Builds a struct-backed `tabular.Exporter` using `tabular.NewStructAdapterFor[T]()`. |
| `csv.NewImporter` | Builds a `tabular.Importer` from an explicit `tabular.RowAdapter`. |
| `csv.NewImporterFor` | Builds a struct-backed `tabular.Importer` using `tabular.NewStructAdapterFor[T]()`. |
| `csv.NewMapExporter` | Validates `[]tabular.ColumnSpec` with `tabular.NewMapAdapterFromSpecs(specs)` and returns a map-backed exporter or the validation error; it does not accept `mapOpts`. |
| `csv.NewMapImporter` | Validates `[]tabular.ColumnSpec` with `tabular.NewMapAdapterFromSpecs(specs, mapOpts...)`; pass `nil` for the `[]tabular.MapOption` argument when no row validators are needed. |
| `csv.NewTypedExporterFor` | Wraps the struct exporter in `tabular.TypedExporter[T]` so export accepts `[]T`. |
| `csv.NewTypedImporterFor` | Wraps the struct importer in `tabular.TypedImporter[T]` so import returns `[]T`. |
| `csv.WithCRLF` | Writes Windows-style CRLF line endings on export. |
| `csv.WithComment` | Enables a comment rune for import; default `0` disables comment handling. |
| `csv.WithExportDelimiter` | Sets the export delimiter; default is comma. |
| `csv.WithImportDelimiter` | Sets the import delimiter; default is comma. |
| `csv.WithSkipRows` | Skips leading rows before header/data processing; negative values are clamped to `0`. |
| `csv.WithoutHeader` | Treats the first non-skipped row as data and maps columns positionally in schema order. |
| `csv.WithoutTrimSpace` | Disables framework-level trimming; this affects empty-row detection, header matching, and cell parsing. |
| `csv.WithoutWriteHeader` | Suppresses the header row during export. |

## Model Definition

CSV uses the same `tabular` struct tag as Excel. See the [Excel documentation](./excel#tabular-tag) for the full tag reference.

```go
type Employee struct {
    orm.FullAuditedModel `tabular:"-"`

    Name       string          `tabular:"Name,width=20"`
    Email      string          `tabular:"Email,width=30"`
    Department string          `tabular:"Department"`
    JoinDate   timex.Date      `tabular:"Join Date,format=2006-01-02"`
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
| `csv.WithCRLF()` | LF only | Use Windows-style CRLF line endings |

```go
// TSV export with Windows line endings
exporter := csv.NewExporterFor[Employee](
    csv.WithExportDelimiter('\t'),
    csv.WithCRLF(),
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

Negative `WithSkipRows` values are clamped to `0`. `WithoutTrimSpace()` also
affects empty-row detection, header matching, and cell parsing. This trimming is
performed by VEF's `tabular` layer; the underlying Go standard-library CSV
reader does not use `TrimLeadingSpace`. `Import` calls that reader's `ReadAll`
before parsing, so peak memory scales with the file size plus the materialized
result slice. The reader uses `FieldsPerRecord = -1`, so ragged rows are
accepted and missing mapped cells can be handled by the tabular adapter.

### Header vs No-Header Mode

| Mode | Column Mapping |
| --- | --- |
| With header (default) | Header names → `tabular` tag names |
| Without header | Column position → `tabular` field order |

When using `WithoutHeader()`, columns are matched by position. Use the `order` tag attribute to control field ordering:

```go
type Record struct {
    Name  string `tabular:"Name,order=0"`
    Email string `tabular:"Email,order=1"`
    Age   int    `tabular:"Age,order=2"`
}
```

`WithSkipRows` is applied before header detection. With headers enabled,
`rows[skipRows]` is the header row and data starts after it; with
`csv.WithoutHeader()`, the first non-skipped row is parsed as data. Import error
row numbers are 1-based CSV file line positions after this offset, matching the
line a user sees in a text editor.

### Custom Parser

```go
importer := csv.NewImporterFor[Employee]()

importer.RegisterParser("date", tabular.ParserFunc(func(cellValue string, targetType reflect.Type) (any, error) {
    return time.Parse("01/02/2006", cellValue)
}))
```

### Validation

Imported records are automatically validated using `validator.Validate(...)`, same as the Excel importer.

## Error Handling

| Error | Meaning |
| --- | --- |
| `tabular.ErrDataMustBeSlice` | Export data must be a slice |
| `tabular.ErrNoDataRowsFound` | No data rows remain after `csv.WithSkipRows` and optional header handling |
| `tabular.ErrDuplicateHeaderName` | Duplicate non-empty column names in the header |
| `tabular.ErrUnsetField` | Struct field cannot be set |

The CSV package exposes `ImportOption` and `ExportOption` as the option function
types behind the import/export options above.

Top-level errors are fatal read/write or structural failures, including
`ReadAll` failures, no-data files, export schema mismatches, and final writer
flush failures (`flush CSV writer: ...`). Parse failures, validator failures,
and adapter commit failures are collected as `[]tabular.ImportError`; import can
return `err == nil` with non-empty row-level errors, and the affected rows are
skipped while later rows continue.

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
