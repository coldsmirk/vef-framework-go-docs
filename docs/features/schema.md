---
sidebar_position: 7
---

# Schema Inspection

VEF includes a schema inspection service and a built-in resource for reading database structure through the application API.

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
| `ListTriggers(ctx)` | `[]schema.Trigger` | list triggers |

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
| `list_triggers` | none | `[]schema.Trigger` | list current triggers |

Implementation details visible in source:

- each action currently sets a per-operation rate-limit max of `60`
- `get_table_schema` validates that `name` is present
- missing tables are mapped to `result.ErrCodeSchemaTableNotFound`

## Public Schema Types

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
| `materialized` | `bool` | whether the view is materialized |

### `schema.Trigger`

| Field | Type | Meaning |
| --- | --- | --- |
| `name` | `string` | trigger name |
| `table` | `string` | target table when applicable |
| `view` | `string` | target view when applicable |
| `actionTime` | `string` | before/after/instead-of timing |
| `events` | `[]string` | trigger events such as insert or update |
| `forEachRow` | `bool` | row-level or statement-level trigger |
| `body` | `string` | trigger body definition |

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
