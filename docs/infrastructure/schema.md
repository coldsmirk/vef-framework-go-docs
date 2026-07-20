---
sidebar_position: 8
---

# Schema Inspection

VEF includes a schema inspection service and a built-in resource for reading database structure through the application API.

The built-in implementation inspects the primary data source only. It supports
PostgreSQL, MySQL, and SQLite through Atlas. PostgreSQL uses
`vef.datasource.*.schema` when configured and otherwise defaults to `public`;
MySQL inspects the current `DATABASE()`; SQLite inspects the `main` schema.

## Module Outputs

The schema module provides:

| Output | Meaning |
| --- | --- |
| `schema.Service` | schema inspection service |
| `sys/schema` | built-in RPC resource |

## `schema.Service` Interface

The public schema service exposes:

| Method | Return type | Purpose |
| --- | --- | --- |
| `ListTables(ctx)` | `[]schema.Table` | list tables in the current schema or database |
| `GetTableSchema(ctx, name)` | `*schema.TableSchema` | inspect one table in detail |
| `ListViews(ctx)` | `[]schema.View` | list views |

## Built-In Resource

The schema module registers the `sys/schema` RPC resource, mounted under
`/api` with the standard envelope (`resource`, `action`, `version`,
`params`, `meta`). No operation is public and none declares a dedicated
permission token: every action inherits the API engine's default Bearer
authentication.

Every action sets a custom per-operation rate limit of `max 60`. The window
length is not overridden, so it inherits `vef.api.rate_limit.period`
(default `5m`); the limiter counts per operation + client IP + principal,
in process memory on each node.

| Action | Access | Rate limit | Input | Output |
| --- | --- | --- | --- | --- |
| `list_tables` | Bearer auth | `max 60` | none | `[]schema.Table` |
| `get_table_schema` | Bearer auth | `max 60` | `GetTableSchemaParams` | `schema.TableSchema` |
| `list_views` | Bearer auth | `max 60` | none | `[]schema.View` |

Purpose of each action:

- `list_tables` — enumerates the tables of the inspected schema/database,
  with name, schema, and comment. It takes no framework-defined input
  parameters.
- `get_table_schema` — inspects one table in detail: columns, primary key,
  indexes, unique keys, foreign keys, and check constraints.
- `list_views` — enumerates the views of the inspected schema/database,
  including their definition SQL and projected columns. It takes no
  framework-defined input parameters.

`GetTableSchemaParams` (input of `get_table_schema`):

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `name` | `string` | Yes | name of the table to inspect; a missing value fails validation, an unknown table returns the table-not-found business error |

Behavior visible in source:

- `get_table_schema` validates that `name` is present (standard validation
  error when missing)
- a table that inspection cannot find surfaces as the plain Go sentinel
  `schema.ErrTableMissing` inside the service and is mapped by the resource
  to the `schema.ErrTableNotFound` business error (`ErrCodeTableNotFound`,
  `2300`)
- the RPC response uses the standard result envelope: the HTTP status stays
  `200` and business errors are carried by the body `code`

## Error API

| API | Meaning |
| --- | --- |
| `schema.ErrTableNotFound` | business error returned when a requested table does not exist |
| `schema.ErrCodeTableNotFound` (`2300`) | numeric business error code for missing tables |
| `schema.ErrTableMissing` | plain Go sentinel reported by `Service` when inspection cannot find the requested table — for Go callers using the service directly; the RPC surface still maps to `ErrTableNotFound` |

## Responses by Action

All public schema DTOs use their JSON field names as the wire contract.
Fields tagged with `omitempty` are omitted when empty: `schema`, `comment`,
`primaryKey`, `indexes`, `uniqueKeys`, `foreignKeys`, `checks`, `default`,
`isPrimaryKey`, `isAutoIncrement`, `predicate`, `hasExpressions`,
`onUpdate`, `onDelete`, `name` on `schema.PrimaryKey`, and `schema` /
`comment` / `columns` on `schema.View`.

### `list_tables` — `[]schema.Table`

One entry per table in the inspected schema/database.

| Field | Type | Description |
| --- | --- | --- |
| `name` | `string` | table name |
| `schema` | `string` | schema name when available (PostgreSQL: the inspected schema; MySQL: the current database; SQLite: `main`); omitted when empty |
| `comment` | `string` | table comment; omitted when empty |

### `get_table_schema` — `schema.TableSchema`

The full structure of one table.

| Field | Type | Description |
| --- | --- | --- |
| `name` | `string` | table name |
| `schema` | `string` | schema name; omitted when empty |
| `comment` | `string` | table comment; omitted when empty |
| `columns` | `[]schema.Column` | all columns, in database order |
| `primaryKey` | `*schema.PrimaryKey` | primary key definition; omitted when the table has none |
| `indexes` | `[]schema.Index` | non-unique indexes; omitted when empty |
| `uniqueKeys` | `[]schema.UniqueKey` | unique constraints and unique indexes; omitted when empty |
| `foreignKeys` | `[]schema.ForeignKey` | foreign key constraints; omitted when empty |
| `checks` | `[]schema.Check` | check constraints; omitted when empty |

#### `schema.Column`

| Field | Type | Description |
| --- | --- | --- |
| `name` | `string` | column name |
| `type` | `string` | raw database type |
| `nullable` | `bool` | whether the column is nullable |
| `default` | `string` | default expression exactly as the database reports it; omitted when the column has no default |
| `comment` | `string` | column comment; omitted when empty |
| `isPrimaryKey` | `bool` | whether the column participates in the primary key; omitted when `false` |
| `isAutoIncrement` | `bool` | whether the column auto-generates its value; omitted when `false` |

`isAutoIncrement` is detected for MySQL `AUTO_INCREMENT`, SQLite
`AUTOINCREMENT`, PostgreSQL identity columns, and PostgreSQL `serial`,
`bigserial`, or `smallserial` raw types.

#### `schema.PrimaryKey`

| Field | Type | Description |
| --- | --- | --- |
| `name` | `string` | primary key constraint name; omitted when the database reports none |
| `columns` | `[]string` | primary key columns, in key order |

#### `schema.Index`

| Field | Type | Description |
| --- | --- | --- |
| `name` | `string` | index name |
| `columns` | `[]string` | indexed columns, in index order; an expression part appears as an empty string |

#### `schema.UniqueKey`

| Field | Type | Description |
| --- | --- | --- |
| `name` | `string` | unique key name |
| `columns` | `[]string` | unique columns, in index order; an expression part appears as an empty string |
| `predicate` | `string` | partial-index predicate when the unique index is conditional; omitted when empty |
| `hasExpressions` | `bool` | `true` when the unique index contains expression columns (such keys do not guarantee plain column-tuple uniqueness); omitted when `false` |

#### `schema.ForeignKey`

| Field | Type | Description |
| --- | --- | --- |
| `name` | `string` | foreign key name |
| `columns` | `[]string` | local columns |
| `refTable` | `string` | referenced table |
| `refColumns` | `[]string` | referenced columns |
| `onUpdate` | `string` | update action: `CASCADE`, `SET NULL`, `SET DEFAULT`, `RESTRICT`, or `NO ACTION`; omitted when the database reports none |
| `onDelete` | `string` | delete action, same vocabulary as `onUpdate`; omitted when the database reports none |

#### `schema.Check`

| Field | Type | Description |
| --- | --- | --- |
| `name` | `string` | check constraint name |
| `expr` | `string` | check expression |

### `list_views` — `[]schema.View`

One entry per view in the inspected schema/database.

| Field | Type | Description |
| --- | --- | --- |
| `name` | `string` | view name |
| `schema` | `string` | schema name; omitted when empty |
| `definition` | `string` | view definition SQL |
| `comment` | `string` | view comment; omitted when empty (SQLite reports no view comments) |
| `columns` | `[]string` | projected columns, in ordinal order; omitted when empty |

View listing uses database-specific queries. PostgreSQL and MySQL read from
`information_schema.views`; SQLite reads `sqlite_schema` and excludes internal
`sqlite_%` views.

## Minimal Request Example

```json
{
  "resource": "sys/schema",
  "action": "get_table_schema",
  "version": "v1",
  "params": {
    "name": "sys_user"
  }
}
```

## Intended Use

This feature is useful for:

- admin tooling
- internal developer tooling
- schema-aware integrations
- MCP or prompt workflows that need database metadata

## Next Step

Read [Built-in Resources](../reference/built-in-resources) to see how `sys/schema` relates to the other framework-provided RPC resources.
