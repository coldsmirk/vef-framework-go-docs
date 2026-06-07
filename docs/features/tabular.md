---
sidebar_position: 12
---

# Tabular

The `tabular` package is the foundation for structured data import/export in VEF. It provides a tag-driven schema system and common interfaces that are implemented by the [Excel](./excel) and [CSV](./csv) packages.

## Architecture

```
tabular (core)
├── Schema      — tag parsing, column metadata
├── Importer    — import interface
├── Exporter    — export interface
├── Formatter   — export value formatting
└── ValueParser — import value parsing

excel (implementation)          csv (implementation)
├── excel.NewImporterFor[T]()   ├── csv.NewImporterFor[T]()
└── excel.NewExporterFor[T]()   └── csv.NewExporterFor[T]()
```

## `tabular` Tag

Use struct tags to define how fields map to columns in Excel/CSV files:

```go
type Employee struct {
    orm.FullAuditedModel `tabular:"-"`

    Name       string          `tabular:"姓名,width=20"`
    Email      string          `tabular:"邮箱,width=30"`
    Department string          `tabular:"name=部门,order=2,width=15"`
    JoinDate   timex.Date      `tabular:"入职日期,format=2006-01-02,width=15"`
    Salary     decimal.Decimal `tabular:"薪资,width=12,formatter=money"` // format strings containing commas (e.g. "#,##0.00") can't be set via tag — register a custom formatter
    Status     string          `tabular:"状态,default=active,formatter=status"`
}
```

### Tag Attributes

| Attribute | Type | Description |
| --- | --- | --- |
| (default value) | string | Column header name |
| `name` | string | Explicit column name (alternative to default) |
| `order` | int | Column display order (0-based, default: field declaration order) |
| `width` | float64 | Column width hint (used by Excel export) |
| `default` | string | Default value for empty cells during import |
| `format` | string | Format template (date format, number format) |
| `formatter` | string | Custom formatter name for export |
| `parser` | string | Custom parser name for import |

### Special Tags

| Tag | Meaning |
| --- | --- |
| `tabular:"-"` | Ignore this field completely |
| `tabular:"dive"` | Recurse into embedded struct fields |

The tag parser uses comma-separated `key=value` pairs. Semicolons are not
separators; `tabular:"name=ID;order=1"` is treated as one `name` value.

## Schema

The `Schema` type pre-parses tabular metadata from struct fields at initialization time:

```go
schema := tabular.NewSchemaFor[Employee]()

columns := schema.Columns()       // []*Column — all parsed columns
names := schema.ColumnNames()     // []string{"姓名", "邮箱", ...}
count := schema.ColumnCount()     // 6
```

Columns are automatically sorted by `order` attribute. Fields without an explicit `order` use their declaration order.

`NewSchemaFromSpecs` validates dynamic schemas at construction time: missing
`Key` returns `ErrMissingColumnKey`, missing `Type` returns
`ErrMissingColumnType`, duplicate keys return `ErrDuplicateColumnKey`, and
duplicate resolved header names return `ErrDuplicateHeaderName`.

`ColumnSpec.Required`, per-column `Validators`, and map-level `RowValidator`
run during map-row import. Multiple map-row validation failures are combined
with `errors.Join`, so callers can inspect the joined error with `errors.Is`.

## Interfaces

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

### Formatter (Export)

```go
type Formatter interface {
    Format(value any) (string, error)
}

// Convenience adapter
tabular.FormatterFunc(func(value any) (string, error) { ... })
```

### ValueParser (Import)

```go
type ValueParser interface {
    Parse(cellValue string, targetType reflect.Type) (any, error)
}

// Convenience adapter
tabular.ParserFunc(func(cellValue string, targetType reflect.Type) (any, error) { ... })
```

## Header Mapping and Row Import

`BuildHeaderMapping` is the shared header resolver used by the CSV and Excel
drivers. When trimming is enabled, it trims header names before matching,
skips empty and unknown headers, and treats duplicate non-empty headers as a
fatal `ErrDuplicateHeaderName`. When an importer is configured with
`WithoutHeader()`, the drivers use `DefaultPositionalMapping`, which maps
input column positions to schema columns in order.

`ParseRow` applies defaults before parsing cell values, skips cells that are
still empty after default substitution, and returns row-level `ImportError`
values for parse, validation, and commit failures. If row errors are returned,
the row builder does not commit a partial row.

## Default Type Support

The built-in `DefaultParser` and `DefaultFormatter` handle these types automatically:

| Go Type | Import (Parse) | Export (Format) |
| --- | --- | --- |
| `string` | Direct assignment | Direct output |
| `int`, `int8`–`int64` | Integer parsing | Integer formatting |
| `uint`, `uint8`–`uint64` | Unsigned int parsing | Integer formatting |
| `float32`, `float64` | Float parsing | Float formatting |
| `bool` | `true`/`false`, `1`/`0` | Bool formatting |
| `decimal.Decimal` | Decimal string parsing | Decimal formatting |
| `timex.Date` / `timex.DateTime` | Uses `format` attribute | Uses `format` attribute |
| `*T` (pointer types) | Nil for empty, parsed otherwise | Handles nil gracefully |

## Error Types

### ImportError

```go
type ImportError struct {
    Row    int    // 1-based row number (including header)
    Column string // Column header name
    Field  string // Struct field name
    Err    error  // Underlying error
}
```

Import errors are returned per-row without stopping the import process. This allows batch processing where valid rows are imported and invalid rows are reported.

### ExportError

```go
type ExportError struct {
    Row    int    // 0-based data row index
    Column string // Column header name
    Field  string // Struct field name
    Err    error  // Underlying error
}
```

## Implementations

| Package | Format | Documentation |
| --- | --- | --- |
| `excel` | `.xlsx` (Excel) | [Excel Documentation](./excel) |
| `csv` | `.csv` (CSV/TSV) | [CSV Documentation](./csv) |

## Public Core APIs

| API group | Public surface |
| --- | --- |
| schema | `NewSchema`, `NewSchemaFor[T]`, `NewSchemaFromSpecs`, `Column`, `ColumnSpec`, and `Schema` lookup methods |
| adapters | `NewStructAdapter`, `NewStructAdapterFor[T]`, `NewMapAdapter`, `NewMapAdapterFromSpecs`, `RowAdapter`, `RowReader`, `RowWriter`, `RowView`, `RowBuilder` |
| typed wrappers | `NewTypedImporter[T]`, `NewTypedExporter[T]`, `TypedImporter[T]`, `TypedExporter[T]` |
| mapping/parsing | `BuildHeaderMapping`, `DefaultPositionalMapping`, `ColumnMapping`, `NewColumnMapping`, `ParseRow`, `ParseRowOptions`, `MappingOptions`, `MapOption`, `WithRowValidator`, `RowValidator`, `CellValidator`, `IsEmptyRow` |
| formatter/parser registry | `ResolveFormatter`, `ResolveFormatters`, `ResolveParser`, `ResolveParsers`, `IsDefaultFormatter`, `NewDefaultFormatter`, `NewDefaultParser` |
| constants/errors | `TagTabular`, `IgnoreField`, `AttrDive`, `AttrName`, `AttrOrder`, `AttrWidth`, `AttrDefault`, `AttrFormatter`, `AttrParser`, `AttrFormat`, and tabular sentinels such as `ErrDataMustBeSlice`, `ErrDuplicateColumnKey`, `ErrDuplicateHeaderName`, `ErrMissingColumnKey`, `ErrMissingColumnType`, `ErrNoDataRowsFound`, `ErrRequiredMissing`, `ErrSchemaMismatch`, `ErrTypedRowMismatch`, `ErrUnknownColumn`, `ErrUnsetField`, `ErrUnsupportedType` |

The tabular package audit currently locks **143 public tabular entries** in the
generated API ledger. The grouped member surface covers **75 grouped tabular
field/method entries** across **20 tabular receiver/type families**: **37
exported tabular field entries** and **38 exported tabular method entries**.
The generated public API index remains the complete signature list; this page
documents the schema, mapping, parser/formatter, adapter, and error-contract
families.

Additional audited fields and adapter methods:

| Surface | Public API |
| --- | --- |
| column metadata | `Column.Default`, `Column.Width`, `Column.Order`, `Column.Parser`, `Column.ParserFn`, `Column.FormatterFn`, `Column.Index` |
| dynamic specs | `ColumnSpec.Default`, `ColumnSpec.Width`, `ColumnSpec.Order`, `ColumnSpec.Parser`, `ColumnSpec.ParserFn`, `ColumnSpec.FormatterFn` |
| parsing options | `MappingOptions.TrimSpace`, `ParseRowOptions.TrimSpace` |
| row adapter contract | `RowAdapter.Writer`, `RowReader.All`, `RowView.Get`, `RowBuilder.Set`, `RowWriter.NewRow`, `RowWriter.Commit`, `RowWriter.Build` |
| schema lookup | `Schema.ColumnByKey`, `Schema.ColumnByName` |
| typed wrappers | `TypedExporter.Inner`, `TypedImporter.Inner` |
| error wrapping | `ImportError.Unwrap`, `ExportError.Unwrap`; callers can use `errors.Unwrap` / `errors.Is` |

## CRUD Integration

The `Export` and `Import` CRUD builders use `tabular` internally:

```go
// Export builder
crud.NewExport[Employee, EmployeeSearch]().
    WithDefaultFormat("excel")

// Import builder
crud.NewImport[Employee]().
    WithDefaultFormat("excel").
    WithPreImport(func(ctx context.Context, models []Employee) error {
        // Validate or transform before insert
        return nil
    })
```
