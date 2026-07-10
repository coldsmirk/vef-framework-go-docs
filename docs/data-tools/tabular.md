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

`ResolveFormatter(col, registry)` / `ResolveParser(col, registry)` apply that precedence for a single column; `ResolveFormatters` / `ResolveParsers` call them for every column once up front, aligned with `schema.Columns()`, so drivers don't repeat registry lookups per cell. `tabular.IsDefaultFormatter(col, registry)` reports whether a column resolves to the built-in default (no `FormatterFn` and no registered named `Formatter`); the Excel exporter uses this to decide between a native typed cell and a formatted string (see [Excel → Native Typed Cells](./excel#native-typed-cells)).

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

`excel.ErrSheetIndexOutOfRange` is the one driver-specific sentinel — see [Excel → Error Handling](./excel#error-handling).

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

## Format Backends

The format backends implement this core for concrete file types. Their APIs mirror each other:

- [CSV](./csv) — delimiter-separated files, with and without header rows
- [Excel](./excel) — `.xlsx` workbooks, sheets, and styling options

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
