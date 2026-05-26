---
sidebar_position: 9
---

# CSV / Excel Import & Export

VEF ships a unified tabular engine in the [`tabular`](#package-layout) package, with two thin format drivers in [`csv`](#csv-package) and [`excel`](#excel-package). Both packages expose the same factory shape and route every read / write through a single `RowAdapter` abstraction so that:

- **Static rows** (Go structs annotated with `tabular` tags) and
- **Dynamic rows** (runtime-defined columns over `map[string]any`)

share one importer / exporter pipeline. You pick the adapter; the format driver does the rest.

## When to Use Which

| Scenario | Recommended factory |
| --- | --- |
| You already have a Go struct (e.g. a model) describing every column | `csv.NewImporterFor[T]` / `excel.NewExporterFor[T]` (and `*Typed*` variants) |
| Columns are decided at runtime — multi-tenant tables, user-defined templates, dynamic forms | `csv.NewMapImporter` / `excel.NewMapExporter` driven by `[]tabular.ColumnSpec` |
| You build the row source yourself (channels, custom domain types) | Implement `tabular.RowAdapter` and pass it to `csv.NewImporter` / `excel.NewExporter` |

The import return type follows the adapter:

- struct adapter → `[]T`
- map adapter → `[]map[string]any`

## Package Layout

```
tabular/   // schema, columns, adapters, formatter / parser, errors
  ├── adapter.go        // RowAdapter, RowReader, RowView, RowWriter, RowBuilder
  ├── schema.go         // Schema, Column
  ├── struct_adapter.go // StructAdapter (struct rows + framework validator)
  ├── map_adapter.go    // MapAdapter (map rows + Required / Validators / RowValidator)
  ├── spec.go           // ColumnSpec, NewSchemaFromSpecs, NewMapAdapterFromSpecs
  ├── resolver.go       // FormatterFn / Formatter / default precedence resolution
  ├── mapping.go        // Header → schema column mapping
  ├── typed.go          // TypedImporter[T] / TypedExporter[T] wrappers
  └── errors.go         // Shared error sentinels (ErrRequiredMissing, ErrSchemaMismatch, …)

csv/       // CSV driver: NewImporter, NewExporter, NewImporterFor, NewExporterFor,
           //              NewMapImporter, NewMapExporter, NewTyped*For
excel/     // Excel driver, mirrors csv with Excel-specific options (sheet name, etc.)
```

`csv` / `excel` no longer hold model-specific reflection or private error types. Common errors live in `tabular` (`ErrDataMustBeSlice`, `ErrRequiredMissing`, `ErrUnknownColumn`, `ErrSchemaMismatch`, …); only `excel.ErrSheetIndexOutOfRange` remains driver-specific.

## Static (Struct-Backed) Usage

Tag your fields with `tabular`:

```go
type User struct {
    ID       int       `tabular:"name=用户ID;order=1"`
    Name     string    `tabular:"name=姓名;order=2" validate:"required"`
    Birthday time.Time `tabular:"name=生日;format=2006-01-02;order=3"`
    Active   bool      `tabular:"name=激活;default=false;order=4"`
    Internal string    `tabular:"-"` // ignored
}
```

Recognised tag attributes (see `tabular/constants.go`):

| Attribute | Meaning |
| --- | --- |
| `name` | Header text (defaults to field name) |
| `order` | Column order (stable sort, lower first) |
| `width` | Column width hint (used by Excel) |
| `default` | Default cell value used during import when source cell is empty |
| `format` | Format template consumed by the default formatter / parser, e.g. `"2006-01-02"`, `"%.2f"` |
| `formatter` | Name of a registered formatter (export side) |
| `parser` | Name of a registered parser (import side) |
| `dive` | Recursively visit an embedded struct |
| `-` | Skip this field entirely |

Validation is delegated to the framework `validator` package — add `validate:"…"` tags as usual; they run automatically when each row is committed.

### Export

```go
exp := csv.NewExporterFor[User]()
buf, err := exp.Export(users) // users: []User or []*User
// or write straight to disk:
err = exp.ExportToFile(users, "users.csv")
```

For Excel:

```go
exp := excel.NewExporterFor[User](excel.WithSheetName("Users"))
buf, err := exp.Export(users)
```

### Import

```go
imp := csv.NewImporterFor[User]()
result, rowErrors, err := imp.Import(reader)
if err != nil {
    return err // top-level failure (e.g. malformed file)
}
users := result.([]User)
for _, ie := range rowErrors {
    log.Warnf("row %d column %s: %v", ie.Row, ie.Column, ie.Err)
}
```

Per-row failures (parse error, struct validator failure, etc.) are aggregated into `[]tabular.ImportError` and **do not** abort the import. The top-level `error` is reserved for fatal issues (invalid file, unreadable header).

### Typed Wrappers

The any-typed return value is sometimes inconvenient. Both packages expose a generic wrapper that performs the type assertion for you:

```go
imp := csv.NewTypedImporterFor[User]()
users, rowErrors, err := imp.Import(reader) // users is []User, no assertion needed

exp := csv.NewTypedExporterFor[User]()
buf, err := exp.Export(users)               // accepts []User directly
```

`TypedImporter` / `TypedExporter` wrap the underlying `tabular.Importer` / `tabular.Exporter`. Use `Inner()` if you need to call `RegisterParser` / `RegisterFormatter` directly on the inner instance.

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
| `Key` | yes | Logical id and the map key used to read / write the cell. Must be unique. |
| `Type` | yes | `reflect.Type` of the parsed value. Use `reflect.TypeFor[T]()`. |
| `Name` | no | Header text. Defaults to `Key`. |
| `Order` | no | Stable sort key for column order. |
| `Width` | no | Excel column width hint. |
| `Default` | no | Default cell value when the source is empty. |
| `Format` | no | Template for the built-in formatter / parser (dates, floats, …). |
| `Formatter` / `Parser` | no | Names looked up against the importer / exporter registries. |
| `FormatterFn` / `ParserFn` | no | `tabular.Formatter` / `tabular.ValueParser` instances bound directly on the column. |
| `Required` | no | Empty cells are reported as `ErrRequiredMissing` during import. |
| `Validators` | no | `[]CellValidator` run after parsing. |

`NewSchemaFromSpecs` validates the slice eagerly — missing `Key`, missing `Type` and duplicate keys all surface as construction-time errors (`tabular.ErrMissingColumnKey`, `ErrMissingColumnType`, `ErrDuplicateHeaderName`).

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

Errors from `Required`, `Validators` and the `RowValidator` are joined into a single error per row via `errors.Join`. To enumerate the leaves:

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

## Custom Formatters & Parsers

Implement the small interfaces from `tabular`:

```go
type Formatter interface {
    Format(value any) (string, error)
}

type ValueParser interface {
    Parse(cellValue string, targetType reflect.Type) (any, error)
}
```

`tabular.FormatterFunc` and `tabular.ParserFunc` adapt plain functions to these interfaces.

Three precedence levels are evaluated for every column (see `tabular/resolver.go`):

1. **`Column.FormatterFn` / `Column.ParserFn`** — instances bound directly on the column (highest priority).
2. **`Column.Formatter` / `Column.Parser`** — named lookup against the importer / exporter registry, populated via `RegisterFormatter` / `RegisterParser`.
3. **Default formatter / parser** — uses `Column.Format` for dates, floats, etc.

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

## Custom RowAdapter

You can plug any data source into the engine by implementing `tabular.RowAdapter`:

```go
type RowAdapter interface {
    Schema() *Schema
    Reader(data any) (RowReader, error)
    Writer(capacity int) RowWriter
}
```

Then pass the adapter directly:

```go
adapter := myStreamingAdapter()
imp := csv.NewImporter(adapter)
exp := excel.NewExporter(adapter)
```

This is useful for channel-driven sources, joined views, or domain types that are neither plain structs nor maps.

## CSV Package

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

CSV options:

| Option | Default | Purpose |
| --- | --- | --- |
| `WithImportDelimiter(r)` | `,` | Field delimiter for import |
| `WithoutHeader()` | header on | Treat first row as data; columns mapped positionally in schema order |
| `WithSkipRows(n)` | `0` | Skip the first `n` rows before reading |
| `WithoutTrimSpace()` | trim on | Disable cell trimming (also affects empty-row detection and header matching) |
| `WithComment(r)` | none | Lines starting with this rune are ignored |
| `WithExportDelimiter(r)` | `,` | Field delimiter for export |
| `WithoutWriteHeader()` | header on | Skip the header row on export |
| `WithCRLF()` | LF | Use Windows-style line endings |

## Excel Package

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

Excel options:

| Option | Default | Purpose |
| --- | --- | --- |
| `WithSheetName(name)` | `Sheet1` | Worksheet name on export |
| `WithImportSheetName(name)` | none | Read a worksheet by name |
| `WithImportSheetIndex(i)` | `0` | Read a worksheet by index (returns `excel.ErrSheetIndexOutOfRange` if negative or out of range) |
| `WithSkipRows(n)` | `0` | Skip the first `n` rows before reading |
| `WithoutHeader()` | header on | First non-skipped row is data; positional mapping |
| `WithoutTrimSpace()` | trim on | Disable cell trimming (also affects empty-row detection and header matching) |

`Column.Width` set on a `ColumnSpec` (or via the struct tag `width=…`) is applied to the generated worksheet.

## Header → Column Mapping Rules

Both drivers route header resolution through `tabular.BuildHeaderMapping`:

- Header cells are matched against `Column.Name`.
- An empty header cell is skipped.
- An unknown header cell is skipped (extra columns won't fail the import).
- A duplicate non-empty header is fatal: `tabular.ErrDuplicateHeaderName`.
- When the importer is configured `WithoutHeader()`, the engine falls back to `tabular.DefaultPositionalMapping` — column 0 ↔ first schema column, etc.

## Errors

`tabular` exposes a shared error palette. The most commonly checked sentinels are:

| Error | When it appears |
| --- | --- |
| `ErrDataMustBeSlice` | Export argument is not a slice |
| `ErrSchemaMismatch` | Element type doesn't match the adapter's schema (struct vs map vs other) |
| `ErrUnknownColumn` | Caller addresses a column that isn't in the schema |
| `ErrRequiredMissing` | A `Required` cell is empty during dynamic import |
| `ErrDuplicateHeaderName` | Header row contains a duplicate non-empty header |
| `ErrUnsetField` | Struct field cannot be set (typically unexported) |
| `ErrMissingColumnKey` / `ErrMissingColumnType` | `ColumnSpec` is missing required attributes |
| `ErrTypedRowMismatch` | A `TypedImporter[T]` received rows whose element type isn't `T` |

Use `errors.Is` against these sentinels rather than string matching.

## Migration from Pre-Refactor APIs

If you used the older signatures, update as follows:

| Before | After |
| --- | --- |
| `csv.NewImporter(typ, opts...)` | `csv.NewImporterFor[T](opts...)` _or_ `csv.NewImporter(tabular.NewStructAdapter(typ), opts...)` |
| `excel.NewExporter(typ, opts...)` | `excel.NewExporterFor[T](opts...)` |
| `csv.ErrDataMustBeSlice`, etc. | `tabular.ErrDataMustBeSlice` (and other shared sentinels) |

`excel.ErrSheetIndexOutOfRange` is unchanged — it stays in the `excel` package because it is Excel-specific.

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

For `crud.NewExport` / `crud.NewImport`, which sit on top of these factories, see [CRUD → Export And Import Builders](../guide/crud.md#export-and-import-builders).
