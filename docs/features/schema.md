---
sidebar_position: 7
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

The schema module registers:

| Resource |
| --- |
| `sys/schema` |

Current actions:

| Action | Input params | Output type | Notes |
| --- | --- | --- | --- |
| `list_tables` | none | `[]schema.Table` | list current tables |
| `get_table_schema` | `name: string` | `schema.TableSchema` | returns a dedicated schema-table-not-found business error when the table does not exist |
| `list_views` | none | `[]schema.View` | list current views |

Implementation details visible in source:

- each action currently sets a per-operation rate-limit `Max = 60`
- `get_table_schema` validates that `name` is present
- missing tables are mapped to `schema.ErrCodeTableNotFound` (`ErrCodeTableNotFound`) and the dedicated `schema.ErrTableNotFound` (`ErrTableNotFound`) sentinel
- the RPC response still uses the standard result envelope; business errors are returned in the body, not as a different HTTP status

## Error API

| API | Meaning |
| --- | --- |
| `schema.ErrTableNotFound` | business error returned when a requested table does not exist |
| `schema.ErrCodeTableNotFound` | numeric business error code for missing tables |

## Public Schema Types

Audit note: this page covers 53 public schema entries, including 41 grouped schema field/method entries across 10 schema receiver/type families. The grouped DTO/service surface contains 38 exported schema field entries and 3 exported schema method entries.

All public schema DTOs use their JSON field names as the wire contract. Fields
tagged with `omitempty` are omitted when empty: for example, `schema`,
`comment`, `primaryKey`, `indexes`, `uniqueKeys`, `foreignKeys`, `checks`,
`default`, `isPrimaryKey`, `isAutoIncrement`, `onUpdate`, `onDelete`, and
`columns` on `schema.View`.

### `schema.Table`

| Field | Type | Meaning |
| --- | --- | --- |
| `name` | `string` | table name |
| `schema` | `string` | schema name when available |
| `comment` | `string` | table comment |

### `schema.TableSchema`

| Field | Type | Meaning |
| --- | --- | --- |
| `name` | `string` | table name |
| `schema` | `string` | schema name |
| `comment` | `string` | table comment |
| `columns` | `[]schema.Column` | all columns |
| `primaryKey` | `*schema.PrimaryKey` | primary key definition |
| `indexes` | `[]schema.Index` | non-unique indexes |
| `uniqueKeys` | `[]schema.UniqueKey` | unique constraints |
| `foreignKeys` | `[]schema.ForeignKey` | foreign key constraints |
| `checks` | `[]schema.Check` | check constraints |

### `schema.Column`

| Field | Type | Meaning |
| --- | --- | --- |
| `name` | `string` | column name |
| `type` | `string` | raw database type |
| `nullable` | `bool` | whether the column is nullable |
| `default` | `string` | default expression |
| `comment` | `string` | column comment |
| `isPrimaryKey` | `bool` | whether the column participates in the primary key |
| `isAutoIncrement` | `bool` | whether the column is auto-incrementing |

`isAutoIncrement` is detected for MySQL `AUTO_INCREMENT`, SQLite
`AUTOINCREMENT`, PostgreSQL identity columns, and PostgreSQL `serial`,
`bigserial`, or `smallserial` raw types.

### `schema.PrimaryKey`

| Field | Type | Meaning |
| --- | --- | --- |
| `name` | `string` | primary key name |
| `columns` | `[]string` | primary key columns |

### `schema.Index`

| Field | Type | Meaning |
| --- | --- | --- |
| `name` | `string` | index name |
| `columns` | `[]string` | indexed columns |

### `schema.UniqueKey`

| Field | Type | Meaning |
| --- | --- | --- |
| `name` | `string` | unique key name |
| `columns` | `[]string` | unique columns |

### `schema.ForeignKey`

| Field | Type | Meaning |
| --- | --- | --- |
| `name` | `string` | foreign key name |
| `columns` | `[]string` | local columns |
| `refTable` | `string` | referenced table |
| `refColumns` | `[]string` | referenced columns |
| `onUpdate` | `string` | update action |
| `onDelete` | `string` | delete action |

### `schema.Check`

| Field | Type | Meaning |
| --- | --- | --- |
| `name` | `string` | check constraint name |
| `expr` | `string` | check expression |

### `schema.View`

| Field | Type | Meaning |
| --- | --- | --- |
| `name` | `string` | view name |
| `schema` | `string` | schema name |
| `definition` | `string` | view definition SQL |
| `comment` | `string` | view comment |
| `columns` | `[]string` | projected columns |

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
