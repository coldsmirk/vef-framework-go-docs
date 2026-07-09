---
sidebar_position: 5
---

# Tabular Import & Export

VEF ships a unified tabular engine in the `tabular` package, with two thin format drivers in `csv` and `excel`. All three packages expose the same factory shape and route every read/write through a single `RowAdapter` abstraction so that:

- **Static rows** (Go structs annotated with `tabular` tags), and
- **Dynamic rows** (runtime-defined columns over `map[string]any`)

share one importer/exporter pipeline. You pick the adapter; the format driver does the rest.

## Architecture

```
tabular/   // schema, columns, adapters, formatter / parser, errors
  ├── adapter.go        // RowAdapter, RowReader, RowView, RowWriter, RowBuilder
  ├── schema.go         // Schema, Column
  ├── struct_adapter.go // StructAdapter (struct rows + framework validator)
  ├── map_adapter.go    // MapAdapter (map rows + Required / Validators / RowValidator)
  ├── spec.go           // ColumnSpec, NewSchemaFromSpecs, NewMapAdapterFromSpecs
  ├── resolver.go       // FormatterFn / Formatter / default precedence resolution
  ├── mapping.go        // Header → schema column mapping
  ├── parse_row.go       // ParseRow, IsEmptyRow
  ├── import_rows.go     // ImportRows (shared driver core)
  ├── typed.go           // TypedImporter[T] / TypedExporter[T] wrappers
  └── errors.go          // Shared error sentinels (ErrRequiredMissing, ErrSchemaMismatch, …)

csv/       // CSV driver: NewImporter, NewExporter, NewImporterFor, NewExporterFor,
           //              NewMapImporter, NewMapExporter, NewTyped*For
excel/     // Excel driver, mirrors csv with Excel-specific options (sheet name, etc.)
```

`csv` / `excel` hold no model-specific reflection or private error types. Common errors live in `tabular` (`ErrDataMustBeSlice`, `ErrRequiredMissing`, `ErrUnknownColumn`, `ErrSchemaMismatch`, …); only `excel.ErrSheetIndexOutOfRange` is driver-specific.

## When to Use Which

| Scenario | Recommended factory |
| --- | --- |
| You already have a Go struct (e.g. a model) describing every column | `csv.NewImporterFor[T]` / `csv.NewExporterFor[T]` or `excel.NewImporterFor[T]` / `excel.NewExporterFor[T]` (and `*Typed*` variants) |
| Columns are decided at runtime — multi-tenant tables, user-defined templates, dynamic forms | `csv.NewMapImporter` / `csv.NewMapExporter` or `excel.NewMapImporter` / `excel.NewMapExporter` driven by `[]tabular.ColumnSpec` |
| You build the row source yourself (channels, custom domain types) | Implement `tabular.RowAdapter` and pass it to `csv.NewImporter` / `csv.NewExporter` or `excel.NewImporter` / `excel.NewExporter` |

The import return type follows the adapter:

- struct adapter → `[]T`
- map adapter → `[]map[string]any`

Once you've picked an adapter, choose the format:

| Feature | CSV | Excel |
| --- | --- | --- |
| File format | `.csv` (plain text) | `.xlsx` (binary) |
| Multiple sheets | No | Yes |
| Column widths | Ignored | Applied |
| Delimiter | Configurable | N/A |
| Comment lines | Supported | N/A |
| Trim whitespace | Configurable | N/A |
| Line endings | LF or CRLF | N/A |
| Native typed cells (numbers/dates stay sortable) | N/A (text format) | Yes, by default |
| Dependencies | Go stdlib | [excelize](https://github.com/xuri/excelize) |

Both packages implement `tabular.Importer` / `tabular.Exporter`, so you can swap between them without changing your model definitions.

## The `tabular` Tag

Struct tags define how fields map to columns:

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

The tag parser uses comma-separated `key=value` pairs. Semicolons are not separators; `tabular:"name=ID;order=1"` is treated as one `name` value.

`dive` only recurses into a struct or pointer-to-struct field; on any other kind the field is neither recursed nor emitted, and the framework logs a warning rather than silently dropping it.

Every attribute key and sentinel value above is also exported as a `tabular` constant, for callers that build or inspect tags programmatically instead of writing them as string literals: `TagTabular` (the struct tag name itself, `"tabular"`), `IgnoreField` (the `"-"` ignore sentinel), `AttrName`, `AttrOrder`, `AttrWidth`, `AttrDefault`, `AttrFormat`, `AttrFormatter`, `AttrParser`, and `AttrDive`.

## Schema

`Schema` pre-parses tabular metadata at initialization time — from a struct type via `NewSchemaFor[T]` / `NewSchema`, or from dynamic column specs via `NewSchemaFromSpecs`:

```go
schema := tabular.NewSchemaFor[Employee]()

columns := schema.Columns()             // []*Column — all parsed columns
names := schema.ColumnNames()           // []string{"姓名", "邮箱", ...}
count := schema.ColumnCount()           // 6
col, ok := schema.ColumnByKey("Name")   // lookup by logical key (struct field name)
col, ok = schema.ColumnByName("姓名")    // lookup by header name
```

Columns are automatically sorted by `order`. Fields without an explicit `order` use their declaration order.

`NewSchemaFromSpecs` validates dynamic schemas at construction time: missing `Key` returns `ErrMissingColumnKey`, missing `Type` returns `ErrMissingColumnType`, duplicate keys return `ErrDuplicateColumnKey`, and duplicate resolved header names return `ErrDuplicateHeaderName`.

### `Column`

Every column — struct-derived or dynamic — is represented by the same `Column` struct:

| Field | Meaning |
| --- | --- |
| `Key` | Logical identifier: struct field name, or the map key for dynamic schemas |
| `Name` | Header text shown on export and matched on import; defaults to `Key` |
| `Type` | `reflect.Type` used to parse cell values |
| `Order` | Stable sort key for column order |
| `Width` | Column width hint (Excel) |
| `Default` | Default value used during import when the source cell is empty |
| `Format` | Format template consumed by the default `Formatter`/`ValueParser` |
| `Formatter` / `Parser` | Named lookup against the exporter / importer registry |
| `FormatterFn` / `ParserFn` | `Formatter` / `ValueParser` instance bound directly on the column (highest priority) |
| `Required` | Empty cells trigger `ErrRequiredMissing` during import (dynamic schemas) |
| `Validators` | `[]CellValidator` run after parsing (dynamic schemas) |
| `Index` | Struct field index path used by `StructAdapter`; `nil` for dynamic columns |

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

### Formatter (export)

```go
type Formatter interface {
    Format(value any) (string, error)
}

// Convenience adapter
tabular.FormatterFunc(func(value any) (string, error) { ... })
```

### ValueParser (import)

```go
type ValueParser interface {
    Parse(cellValue string, targetType reflect.Type) (any, error)
}

// Convenience adapter
tabular.ParserFunc(func(cellValue string, targetType reflect.Type) (any, error) { ... })
```

## Header Mapping and Row Import

Both drivers route header resolution through the same shared core:

- `BuildHeaderMapping(headerRow, schema, opts)` matches header cells against `Column.Name`. When `MappingOptions.TrimSpace` is enabled it trims header names before matching. An empty header cell is skipped. An unknown header cell is skipped (extra columns won't fail the import). A duplicate non-empty header is fatal: `ErrDuplicateHeaderName`.
- When the importer is configured `WithoutHeader()`, the drivers fall back to `DefaultPositionalMapping(schema)` — source column 0 maps to the first schema column, and so on.
- `ParseRow(cells, mapping, schema, builder, parsers, rowNumber, opts)` applies `Column.Default` before parsing, skips cells that are still empty after default substitution, and returns row-level `ImportError` values for parse and `Set` failures. If row errors are returned, the row builder does **not** commit a partial row.
- `IsEmptyRow(cells, trimSpace)` reports whether every cell in a row is empty (used to skip blank rows automatically).
- `tabular.ImportRows(rows, adapter, parsers, opts)` parses a materialized `[][]string` table through a `RowAdapter`. `ImportRowsOptions` controls `SkipRows`, `HasHeader` (header-based versus positional mapping), and `TrimSpace`. CSV and Excel importers delegate to this shared core after reading their source format.

`BuildHeaderMapping` and `DefaultPositionalMapping` both return a raw `map[int]int` (source column index → schema column index); `ImportRows` wraps that map once per import via `NewColumnMapping`, producing a `ColumnMapping` that pre-sorts the source indices so `ParseRow` doesn't re-sort on every row. `ParseRow` itself takes a `ParseRowOptions{TrimSpace bool}` (`ParseRowOptions.TrimSpace`) for cell-level trimming, distinct from the header-level `MappingOptions`.

## Static (Struct-Backed) Usage

Tag your fields with `tabular`, then create a typed importer/exporter for the struct with `csv.NewImporterFor[T]` / `excel.NewImporterFor[T]` (or the exporter equivalents). Validation runs through the framework `validator` package — add `validate:"…"` tags as usual; they're checked automatically when each row is committed.

### Export

```go
exp := csv.NewExporterFor[Employee]()
buf, err := exp.Export(employees) // employees: []Employee or []*Employee
// or write straight to disk:
err = exp.ExportToFile(employees, "employees.csv")
```

For Excel:

```go
exp := excel.NewExporterFor[Employee](excel.WithSheetName("Employees"))
buf, err := exp.Export(employees)
```

### Import

```go
imp := csv.NewImporterFor[Employee]()
result, rowErrors, err := imp.Import(reader)
if err != nil {
    return err // top-level failure (e.g. malformed file)
}
employees := result.([]Employee)
for _, ie := range rowErrors {
    log.Warnf("row %d column %s: %v", ie.Row, ie.Column, ie.Err)
}
```

Per-row failures (parse error, struct validator failure, adapter commit failure) are aggregated into `[]tabular.ImportError` and **do not** abort the import. The top-level `error` is reserved for fatal issues (invalid file, unreadable header, no data rows).

### Typed Wrappers

The any-typed return value is sometimes inconvenient. Both packages expose a generic wrapper that performs the type assertion for you:

```go
imp := csv.NewTypedImporterFor[Employee]()
employees, rowErrors, err := imp.Import(reader) // employees is []Employee, no assertion needed

exp := csv.NewTypedExporterFor[Employee]()
buf, err := exp.Export(employees)               // accepts []Employee directly
```

`TypedImporter[T]` / `TypedExporter[T]` wrap the underlying `tabular.Importer` / `tabular.Exporter`. Use `TypedImporter.Inner` / `TypedExporter.Inner` if you need to call `RegisterParser` / `RegisterFormatter` directly on the inner instance. If the wrapped importer ever returns rows whose element type doesn't match `T`, the typed wrapper returns `ErrTypedRowMismatch`.

`csv.NewTypedImporterFor[T]` / `excel.NewTypedImporterFor[T]` (and the exporter equivalents) are convenience wrappers around `tabular.NewTypedImporter[T](inner)` / `tabular.NewTypedExporter[T](inner)`. Call the `tabular` constructors directly when wrapping a hand-built `Importer`/`Exporter` (e.g. one returned by `csv.NewImporter(adapter, ...)` with a custom `RowAdapter`) instead of a `*For` struct adapter.

## Dynamic (Map-Backed) Usage

Dynamic columns let you build a schema at runtime without declaring a struct. Describe each column with `tabular.ColumnSpec`:

```go
import (
    "reflect"
    "time"

    "github.com/coldsmirk/vef-framework-go/csv"
    "github.com/coldsmirk/vef-framework-go/excel"
    "github.com/coldsmirk/vef-framework-go/tabular"
)

specs := []tabular.ColumnSpec{
    {Key: "id",       Name: "用户ID", Type: reflect.TypeFor[int](),       Required: true, Order: 1},
    {Key: "name",     Name: "姓名",   Type: reflect.TypeFor[string](),    Required: true, Order: 2},
    {Key: "birthday", Name: "生日",   Type: reflect.TypeFor[time.Time](), Format: "2006-01-02", Order: 3},
    {Key: "active",   Name: "激活",   Type: reflect.TypeFor[bool](),      Default: "false", Order: 4},
}
```

`ColumnSpec` fields:

| Field | Required | Notes |
| --- | --- | --- |
| `Key` | yes | Logical id and the map key used to read/write the cell. Must be unique. |
| `Type` | yes | `reflect.Type` of the parsed value. Use `reflect.TypeFor[T]()`. |
| `Name` | no | Header text. Defaults to `Key`. |
| `Order` | no | Stable sort key for column order. |
| `Width` | no | Excel column width hint. |
| `Default` (`ColumnSpec.Default`) | no | Default cell value when the source is empty. |
| `Format` | no | Template for the built-in formatter/parser (dates, floats, …). |
| `Formatter` / `Parser` | no | Names looked up against the importer/exporter registries. |
| `FormatterFn` / `ParserFn` | no | `tabular.Formatter` / `tabular.ValueParser` instances bound directly on the column. |
| `Required` | no | Empty cells are reported as `ErrRequiredMissing` during import. |
| `Validators` | no | `[]CellValidator` run after parsing. |

`NewSchemaFromSpecs` validates the slice eagerly — missing `Key`, missing `Type`, duplicate keys, and duplicate resolved names all surface as construction-time errors (`ErrMissingColumnKey`, `ErrMissingColumnType`, `ErrDuplicateColumnKey`, `ErrDuplicateHeaderName`).

### Export

```go
exp, err := excel.NewMapExporter(specs, excel.WithSheetName("Users"))
if err != nil { return err }

buf, err := exp.Export([]map[string]any{
    {"id": 1, "name": "张三", "birthday": time.Now(), "active": true},
    {"id": 2, "name": "李四", "birthday": time.Now(), "active": false},
})
```

CSV is identical:

```go
exp, err := csv.NewMapExporter(specs)
buf, err := exp.Export(rows)
```

### Import

```go
imp, err := csv.NewMapImporter(specs, nil) // nil → no MapAdapter options
if err != nil { return err }

result, rowErrors, err := imp.Import(reader)
if err != nil { return err }

rows := result.([]map[string]any)
```

Behavior worth noting:

- Unknown headers in the source are **skipped silently** — extra columns won't break the import.
- Schema columns missing from the source simply do not appear in the parsed row map (the key is absent rather than zero-valued). This lets `Required` and row-level validators distinguish "missing" from "explicitly zero".
- After `TrimSpace` (on by default for both CSV and Excel; toggle with `WithoutTrimSpace()`) and `Default` substitution, an empty cell causes `Set` to be skipped: structs keep zero values, maps leave the key absent.
- Per-cell parse errors, per-row commit errors, and validator errors all aggregate into `[]tabular.ImportError`; the rest of the file is still processed.

### Row-Level Validation

Pass `MapOption`s as the second argument to `NewMapImporter`:

```go
imp, err := csv.NewMapImporter(specs,
    []tabular.MapOption{
        tabular.WithRowValidator(func(row map[string]any) error {
            if row["name"] == "" {
                return errors.New("name must not be empty")
            }
            return nil
        }),
    },
)
```

Cell-level validation is set per column via `Validators`:

```go
specs := []tabular.ColumnSpec{
    {
        Key:  "email",
        Name: "邮箱",
        Type: reflect.TypeFor[string](),
        Validators: []tabular.CellValidator{
            func(col *tabular.Column, value any) error {
                s, _ := value.(string)
                if !strings.Contains(s, "@") {
                    return fmt.Errorf("invalid email: %q", s)
                }
                return nil
            },
        },
    },
}
```

`ColumnSpec.Required`, per-column `Validators`, and map-level `RowValidator` all run during map-row commit. Errors from `Required`, `Validators`, and `RowValidator` are joined into a single error per row via `errors.Join`. To enumerate the leaves:

```go
for _, ie := range rowErrors {
    if errors.Is(ie.Err, tabular.ErrRequiredMissing) {
        // …
    }
    if multi, ok := ie.Err.(interface{ Unwrap() []error }); ok {
        for _, leaf := range multi.Unwrap() {
            log.Warn(leaf)
        }
    }
}
```

## Custom Formatters and Parsers

Implement the small interfaces from `tabular`, or adapt a plain function with `tabular.FormatterFunc` / `tabular.ParserFunc`. Three precedence levels are evaluated for every column (see `tabular/resolver.go`):

1. **`Column.FormatterFn` / `Column.ParserFn`** — instances bound directly on the column (highest priority).
2. **`Column.Formatter` / `Column.Parser`** — named lookup against the importer/exporter registry, populated via `RegisterFormatter` / `RegisterParser`.
3. **Default formatter/parser** — uses `Column.Format` for dates, floats, etc.

`ResolveFormatter(col, registry)` / `ResolveParser(col, registry)` apply that precedence for a single column; `ResolveFormatters` / `ResolveParsers` call them for every column once up front, aligned with `schema.Columns()`, so drivers don't repeat registry lookups per cell. `tabular.IsDefaultFormatter(col, registry)` reports whether a column resolves to the built-in default (no `FormatterFn` and no registered named `Formatter`); the Excel exporter uses this to decide between a native typed cell and a formatted string (see [Excel → Native Typed Cells](#native-typed-cells)).

Inline (highest priority) example:

```go
yenFormatter := tabular.FormatterFunc(func(v any) (string, error) {
    return fmt.Sprintf("¥%.2f", v), nil
})

specs := []tabular.ColumnSpec{
    {
        Key:         "price",
        Name:        "Price",
        Type:        reflect.TypeFor[float64](),
        FormatterFn: yenFormatter,
    },
}
```

Named registry example (works for both struct and map adapters):

```go
exp := csv.NewExporterFor[Order]()
exp.RegisterFormatter("currency", currencyFormatter)
// columns whose tag has `formatter=currency` will use it
```

```go
importer := csv.NewImporterFor[Employee]()

importer.RegisterParser("date", tabular.ParserFunc(func(cellValue string, targetType reflect.Type) (any, error) {
    return time.Parse("01/02/2006", cellValue)
}))
```

## Custom RowAdapter

Any data source can be plugged into the engine by implementing `tabular.RowAdapter`:

```go
type RowAdapter interface {
    Schema() *Schema
    Reader(data any) (RowReader, error)
    Writer(capacity int) RowWriter
}

type RowReader interface {
    All() iter.Seq2[int, RowView]
}

type RowView interface {
    Get(column *Column) (any, error)
}

type RowWriter interface {
    NewRow() RowBuilder
    Commit(row RowBuilder) error
    Build() any
}

type RowBuilder interface {
    Set(column *Column, value any) error
    Validate() error
    Value() any
}
```

`tabular.NewStructAdapter(typ)` / `tabular.NewStructAdapterFor[T]()` and `tabular.NewMapAdapter(schema, opts...)` / `tabular.NewMapAdapterFromSpecs(specs, opts...)` are the two built-in implementations. Plug a custom one in directly:

```go
adapter := myStreamingAdapter()
imp := csv.NewImporter(adapter)
exp := excel.NewExporter(adapter)
```

This is useful for channel-driven sources, joined views, or domain types that are neither plain structs nor maps.

## Default Type Support

The built-in `DefaultParser` and `DefaultFormatter` (`tabular.NewDefaultParser(format)` / `tabular.NewDefaultFormatter(format)`) handle these types automatically, for both CSV and Excel:

| Go Type | Import (Parse) | Export (Format) |
| --- | --- | --- |
| `string` | Direct assignment | Direct output |
| `int`, `int8`–`int64` | Integer parsing | Integer formatting |
| `uint`, `uint8`–`uint64` | Unsigned int parsing | Integer formatting |
| `float32`, `float64` | Float parsing | Float formatting |
| `bool` | `true`/`false`, `1`/`0` | Bool formatting |
| `decimal.Decimal` | Decimal string parsing | Decimal formatting |
| `time.Time` | Uses `format` attribute (default `time.DateTime`) | Uses `format` attribute (default `time.DateTime`) |
| `timex.Date` / `timex.DateTime` / `timex.Time` | Uses `format` attribute | Uses `format` attribute |
| `*T` (pointer types) | Nil for empty, parsed otherwise | Handles nil gracefully |

## Error Types

Both `ImportError` and `ExportError` implement the standard `Unwrap() error` method, so the standard library's `errors.Unwrap`, `errors.Is`, and `errors.As` all work against them without any `tabular`-specific unwrapping helper.

### ImportError

```go
type ImportError struct {
    Row    int    // 1-based row number (including header)
    Column string // Column header name
    Field  string // Struct field name
    Err    error  // Underlying error
}
```

`ImportError` implements `error` and `Unwrap() error` (`ImportError.Unwrap`). `Err` may itself carry multiple leaf errors joined via `errors.Join` when a single row produces several failures (e.g. multiple `Required` misses, or both a cell validator and a row validator failing) — use `errors.Is` on the `ImportError` to match a specific cause, or assert `Err` against `interface{ Unwrap() []error }` to enumerate every leaf.

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

`ExportError` also implements `error` and `Unwrap() error` (`ExportError.Unwrap`).

### Shared Error Sentinels

`tabular` exposes a shared error palette used by both drivers. Check them with `errors.Is`, not string matching:

| Error | When it appears |
| --- | --- |
| `ErrDataMustBeSlice` | Export argument is not a slice |
| `ErrSchemaMismatch` | Element type doesn't match the adapter's schema (struct vs map vs other) |
| `ErrUnknownColumn` | Caller addresses a column that isn't in the schema |
| `ErrRequiredMissing` | A `Required` cell is empty during dynamic import |
| `ErrNoDataRowsFound` | No data rows remain after skip-rows and optional header handling |
| `ErrDuplicateHeaderName` | Header row contains a duplicate non-empty header, or a dynamic schema resolves two columns to the same name |
| `ErrDuplicateColumnKey` | A dynamic `ColumnSpec` slice has two entries with the same `Key` |
| `ErrUnsetField` | Struct field cannot be set (typically unexported) |
| `ErrMissingColumnKey` / `ErrMissingColumnType` | `ColumnSpec` is missing required attributes |
| `ErrTypedRowMismatch` | A `TypedImporter[T]` / `TypedExporter[T]` received rows whose element type isn't `T` |
| `ErrUnsupportedType` | The default parser was asked to parse into a Go type it doesn't know how to handle |

`excel.ErrSheetIndexOutOfRange` is the one driver-specific sentinel — see [Excel → Error Handling](#error-handling-1).

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

`WithDefaultFormat` accepts a `crud.TabularFormat` (`crud.FormatExcel` / `crud.FormatCsv`, or the equivalent string constants `"excel"` / `"csv"`) and is used when the request doesn't specify a format explicitly. See [CRUD → Export And Import Builders](../data-access/crud#export-and-import-builders) for the full builder API, including `WithExcelOptions`, `WithCsvOptions`, `WithPreExport`, and `WithFilenameBuilder`.

## CSV

The `csv` package provides CSV import/export using the shared `tabular` engine described above. It shares the same `tabular.Importer` and `tabular.Exporter` interfaces as `excel`, making it easy to swap between formats without touching model definitions.

### Package Surface

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

### Model Definition

CSV uses the same `tabular` struct tag described in [The `tabular` Tag](#the-tabular-tag):

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

### Exporting

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

### Importing

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

#### Header vs No-Header Mode

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

### Custom Formatter and Parser

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

### Validation

Imported records are automatically validated using `validator.Validate(...)`, same as the Excel importer.

### Error Handling

Top-level errors are fatal read/write or structural failures, including `ReadAll` failures, no-data files (`tabular.ErrNoDataRowsFound`), export schema mismatches (`tabular.ErrDataMustBeSlice`), duplicate headers (`tabular.ErrDuplicateHeaderName`), and final writer flush failures (`flush CSV writer: ...`). Parse failures, validator failures, and adapter commit failures (including `tabular.ErrUnsetField`) are collected as `[]tabular.ImportError`; import can return `err == nil` with non-empty row-level errors, and the affected rows are skipped while later rows continue.

## Excel

The `excel` package provides Excel import/export using the shared `tabular` engine. It uses [excelize](https://github.com/xuri/excelize) under the hood and integrates with VEF's validation system.

### Package Surface

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

### Sheets

Excel workbooks can hold multiple sheets. Export always writes to a single named sheet (`excel.WithSheetName`, default `Sheet1`); import selects a source sheet either by name (`excel.WithImportSheetName`, which wins if set) or by 0-based index (`excel.WithImportSheetIndex`, default `0`).

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

### Exporting

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

#### Native Typed Cells

Excel export writes native typed cells when a column uses the default formatter and declares no explicit `format`, `formatter`, or `FormatterFn` (checked via `tabular.IsDefaultFormatter`). Integers, floats, booleans, `time.Time`, `timex.Date`, and `timex.DateTime` remain sortable or summable in Excel this way. Once a column sets a format string or a custom formatter, the exporter writes the formatted result as text instead.

Details of the native-cell conversion: nil pointers become empty cells and non-nil pointers are dereferenced; `timex.Date` / `timex.DateTime` are unwrapped to `time.Time` so excelize stores a native date(time) cell; `decimal.Decimal` is converted to `float64` (exact within ~15–16 significant digits, lossy beyond that — a column needing full decimal precision should declare an explicit `format` to render exact text instead); `timex.Time` is intentionally left as text because its zero-date component predates the Excel epoch and would otherwise render a bogus date.

#### Custom Formatter and Parser

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

#### Export in HTTP Handler

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

### Importing

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

### Validation

Imported records are automatically validated using `validator.Validate(...)`. If validation fails, the row is added to `importErrors` and skipped from the result slice.

```go
type Employee struct {
    Name  string `tabular:"Name" validate:"required"`
    Email string `tabular:"Email" validate:"required,email"`
}
```

### Column Mapping Rules

1. The importer matches Excel header names → `tabular` tag name (or field name if no tag).
2. Unmatched Excel columns are silently ignored.
3. Missing Excel columns leave the struct field at its zero value (or `default` if specified).
4. Empty rows are automatically skipped.

When `excel.WithoutHeader()` is used, header matching is bypassed and the importer uses `tabular.DefaultPositionalMapping`: the first source column maps to the first schema column, the second source column maps to the second schema column, and so on.

Default type support for Excel import/export is identical to CSV — see [Default Type Support](#default-type-support) above.

### Error Handling

| Error | Meaning |
| --- | --- |
| `excel.ErrSheetIndexOutOfRange` | Configured sheet index is negative or exceeds available sheets |
| `tabular.ErrNoDataRowsFound` | File has no data rows after `WithSkipRows` and optional header handling |
| `tabular.ErrDuplicateHeaderName` | Duplicate non-empty column names in the header row |
| `tabular.ErrUnsetField` | Struct field cannot be set, typically because it is unexported |

Top-level import errors are fatal file or worksheet failures. Parse failures, validator failures, and adapter commit failures are collected as `[]tabular.ImportError` (`Row`/`Column`/`Field`/`Err`, 1-based row including the header); affected rows are skipped while later rows continue. Export errors use the same `ExportError` shape with a 0-based data row index.

## Cheat Sheet

```go
// Static struct round-trip
imp := csv.NewTypedImporterFor[User]()
exp := csv.NewTypedExporterFor[User]()
buf, _ := exp.Export(users)
imported, errs, _ := imp.Import(buf)

// Dynamic map round-trip
specs := []tabular.ColumnSpec{
    {Key: "id",   Name: "ID",   Type: reflect.TypeFor[int](),    Required: true},
    {Key: "name", Name: "Name", Type: reflect.TypeFor[string]()},
}
exp, _ := excel.NewMapExporter(specs, excel.WithSheetName("Data"))
imp, _ := excel.NewMapImporter(specs, nil)
buf, _ := exp.Export([]map[string]any{{"id": 1, "name": "Alice"}})
rows, errs, _ := imp.Import(buf)

// Dynamic with row-level validation
imp, _ := csv.NewMapImporter(specs,
    []tabular.MapOption{tabular.WithRowValidator(func(r map[string]any) error {
        if r["id"].(int) <= 0 { return errors.New("id must be positive") }
        return nil
    })},
)
```

For `crud.NewExport` / `crud.NewImport`, which sit on top of these factories, see [CRUD → Export And Import Builders](../data-access/crud#export-and-import-builders).
