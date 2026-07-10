---
sidebar_position: 6
---

# CSV

The CSV backend of the [tabular import/export core](./tabular). It implements the shared schema, formatter, and parser contracts for delimiter-separated files.

The `csv` package provides CSV import/export using the shared `tabular` engine described above. It shares the same `tabular.Importer` and `tabular.Exporter` interfaces as `excel`, making it easy to swap between formats without touching model definitions.

## Package Surface

```go
csv.NewImporter(adapter, opts...)            // tabular.Importer
csv.NewExporter(adapter, opts...)            // tabular.Exporter
csv.NewImporterFor[T](opts...)               // struct shortcut
csv.NewExporterFor[T](opts...)               // struct shortcut
csv.NewTypedImporterFor[T](opts...)          // returns []T directly
csv.NewTypedExporterFor[T](opts...)          // accepts []T directly
csv.NewMapImporter(specs, mapOpts, opts...)  // dynamic map importer
csv.NewMapExporter(specs, opts...)           // dynamic map exporter
```

The CSV option marker types are `csv.ExportOption` and `csv.ImportOption`. `csv.NewMapExporter` validates `specs` with `tabular.NewMapAdapterFromSpecs(specs)` and does **not** accept `mapOpts`; `csv.NewMapImporter` validates with `tabular.NewMapAdapterFromSpecs(specs, mapOpts...)` — pass `nil` for `mapOpts` when no row validators are needed.

| Option | Default | Purpose |
| --- | --- | --- |
| `csv.WithImportDelimiter(r)` | `,` | Field delimiter for import |
| `csv.WithoutHeader()` | header on | Treat first row as data; columns mapped positionally in schema order |
| `csv.WithSkipRows(n)` | `0` | Skip the first `n` rows before reading; negative values are clamped to `0` |
| `csv.WithoutTrimSpace()` | trim on | Disable cell trimming (also affects empty-row detection and header matching) |
| `csv.WithComment(r)` | none (`0`) | Lines starting with this rune are ignored |
| `csv.WithExportDelimiter(r)` | `,` | Field delimiter for export |
| `csv.WithoutWriteHeader()` | header on | Skip the header row on export |
| `csv.WithCRLF()` | LF | Use Windows-style line endings |

## Model Definition

CSV uses the same `tabular` struct tag described in [The `tabular` Tag](./tabular#the-tabular-tag):

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

```go
import "github.com/coldsmirk/vef-framework-go/csv"

exporter := csv.NewExporterFor[Employee]()

// Export to file
err := exporter.ExportToFile(employees, "employees.csv")

// Export to buffer (for HTTP response)
buf, err := exporter.Export(employees)
```

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

## Importing

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

// Or import from an io.Reader (e.g. an uploaded file)
data, importErrors, err = importer.Import(reader)
```

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

`Import` calls the standard-library CSV reader's `ReadAll` before parsing, so peak memory scales with the file size plus the materialized result slice. The reader uses `FieldsPerRecord = -1`, so ragged rows are accepted and missing mapped cells are handled by the tabular adapter (empty cell → `Default` → skip). The trimming performed by `WithoutTrimSpace()` is done by VEF's `tabular` layer; the underlying Go standard-library CSV reader does not use `TrimLeadingSpace`.

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

`WithSkipRows` is applied before header detection. With headers enabled, `rows[skipRows]` is the header row and data starts after it; with `csv.WithoutHeader()`, the first non-skipped row is parsed as data. Import error row numbers are 1-based CSV file line positions after this offset, matching the line a user sees in a text editor.

## Custom Formatter and Parser

```go
exporter := csv.NewExporterFor[Employee]()

exporter.RegisterFormatter("status", tabular.FormatterFunc(func(value any) (string, error) {
    if active, ok := value.(bool); ok && active {
        return "Y", nil
    }
    return "N", nil
}))
```

```go
importer := csv.NewImporterFor[Employee]()

importer.RegisterParser("date", tabular.ParserFunc(func(cellValue string, targetType reflect.Type) (any, error) {
    return time.Parse("01/02/2006", cellValue)
}))
```

## Validation

Imported records are automatically validated using `validator.Validate(...)`, same as the Excel importer.

## Error Handling

Top-level errors are fatal read/write or structural failures, including `ReadAll` failures, no-data files (`tabular.ErrNoDataRowsFound`), export schema mismatches (`tabular.ErrDataMustBeSlice`), duplicate headers (`tabular.ErrDuplicateHeaderName`), and final writer flush failures (`flush CSV writer: ...`). Parse failures, validator failures, and adapter commit failures (including `tabular.ErrUnsetField`) are collected as `[]tabular.ImportError`; import can return `err == nil` with non-empty row-level errors, and the affected rows are skipped while later rows continue.

## Next Step

- [Tabular Import & Export](./tabular) — the shared schema, tags, and interfaces
- [Excel](./excel) — the Excel backend
