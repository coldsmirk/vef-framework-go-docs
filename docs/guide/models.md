---
sidebar_position: 1
---

# Models

VEF models are regular Go structs, but they are usually designed to cooperate with Bun, validation tags, search tags, and the framework's audit conventions.

## The Common Pattern

Most persistent models look like this:

```go
type User struct {
	bun.BaseModel `bun:"table:sys_user,alias:su"`
	orm.FullAuditedModel

	Username string `json:"username" validate:"required,alphanum,max=32" label:"Username"`
	Email    string `json:"email" validate:"omitempty,email,max=128" label:"Email"`
	IsActive bool   `json:"isActive"`
}
```

This combines two different concerns:

- `bun.BaseModel`: table metadata for Bun
- `orm.FullAuditedModel`: framework-owned base fields including ID, creation and update audit columns

## Base Model Types

VEF exposes **five** reusable model types through `orm`. They are designed for anonymous embedding as composable field slices.

### Types Without Primary Key

| Type | Fields | Use Case |
| --- | --- | --- |
| `orm.CreationTrackedModel` | `CreatedAt`, `CreatedBy`, `CreatedByName` | Composite-PK tables that need creation tracking |
| `orm.FullTrackedModel` | `CreatedAt`, `CreatedBy`, `CreatedByName`, `UpdatedAt`, `UpdatedBy`, `UpdatedByName` | Composite-PK tables that need full audit tracking |

### Types With Primary Key

| Type | Fields | Use Case |
| --- | --- | --- |
| `orm.Model` | `ID` only | Reference tables, join tables, minimal records |
| `orm.CreationAuditedModel` | `ID`, `CreatedAt`, `CreatedBy`, `CreatedByName` | Append-only records, logs, outbox tables |
| `orm.FullAuditedModel` | `ID`, `CreatedAt`, `CreatedBy`, `CreatedByName`, `UpdatedAt`, `UpdatedBy`, `UpdatedByName` | Standard mutable entities with full audit trail |

### Choosing The Right Type

Use the smallest type that matches your entity's lifecycle:

- **`orm.Model`**: you only need a primary key, no audit tracking at all
- **`orm.CreationAuditedModel`**: append-only records — writes once, never updated
- **`orm.FullAuditedModel`**: the most common choice — standard mutable entities that track both creation and update metadata
- **`orm.CreationTrackedModel`**: same as `CreationAuditedModel` but without the primary key — useful for composite-PK tables
- **`orm.FullTrackedModel`**: same as `FullAuditedModel` but without the primary key — useful for composite-PK tables

Conceptually, `orm.FullAuditedModel` is the pre-composed form of `orm.Model` + `orm.FullTrackedModel`.

### Internal Field Definitions

Here is exactly what each type contributes, including all struct tags:

```go
// orm.Model — primary key only
type Model struct {
	ID string `json:"id" bun:"id,pk"`
}

// orm.CreationTrackedModel — creation audit without PK
type CreationTrackedModel struct {
	CreatedAt     timex.DateTime `json:"createdAt" bun:",notnull,type:timestamp,default:CURRENT_TIMESTAMP,skipupdate"`
	CreatedBy     string         `json:"createdBy" bun:",notnull,skipupdate" mold:"translate=user?"`
	CreatedByName string         `json:"createdByName" bun:",scanonly"`
}

// orm.FullTrackedModel — full audit without PK
type FullTrackedModel struct {
	CreatedAt     timex.DateTime `json:"createdAt" bun:",notnull,type:timestamp,default:CURRENT_TIMESTAMP,skipupdate"`
	CreatedBy     string         `json:"createdBy" bun:",notnull,skipupdate" mold:"translate=user?"`
	CreatedByName string         `json:"createdByName" bun:",scanonly"`
	UpdatedAt     timex.DateTime `json:"updatedAt" bun:",notnull,type:timestamp,default:CURRENT_TIMESTAMP"`
	UpdatedBy     string         `json:"updatedBy" bun:",notnull" mold:"translate=user?"`
	UpdatedByName string         `json:"updatedByName" bun:",scanonly"`
}

// orm.CreationAuditedModel — PK + creation audit
type CreationAuditedModel struct {
	ID            string         `json:"id" bun:"id,pk"`
	CreatedAt     timex.DateTime `json:"createdAt" bun:",notnull,type:timestamp,default:CURRENT_TIMESTAMP,skipupdate"`
	CreatedBy     string         `json:"createdBy" bun:",notnull,skipupdate" mold:"translate=user?"`
	CreatedByName string         `json:"createdByName" bun:",scanonly"`
}

// orm.FullAuditedModel — PK + full audit
type FullAuditedModel struct {
	ID            string         `json:"id" bun:"id,pk"`
	CreatedAt     timex.DateTime `json:"createdAt" bun:",notnull,type:timestamp,default:CURRENT_TIMESTAMP,skipupdate"`
	CreatedBy     string         `json:"createdBy" bun:",notnull,skipupdate" mold:"translate=user?"`
	CreatedByName string         `json:"createdByName" bun:",scanonly"`
	UpdatedAt     timex.DateTime `json:"updatedAt" bun:",notnull,type:timestamp,default:CURRENT_TIMESTAMP"`
	UpdatedBy     string         `json:"updatedBy" bun:",notnull" mold:"translate=user?"`
	UpdatedByName string         `json:"updatedByName" bun:",scanonly"`
}
```

### Key Tag Details

- **`bun:",skipupdate"`**: `created_at` and `created_by` are set on insert only and never overwritten on update
- **`bun:",scanonly"`**: `created_by_name` and `updated_by_name` are read-side convenience fields — they are populated from JOIN queries but are not persisted as separate columns
- **`mold:"translate=user?"`**: the `mold` transformer translates user IDs into display names via a data dictionary, populating the `*ByName` fields automatically
- **`timex.DateTime`**: the framework's custom timestamp type (see [Timex](../utilities/timex)) — not `time.Time`

## Embedding And Composition

The smaller base model types are useful when an entity does not need the full `orm.FullAuditedModel` field set.

```go
// Minimal: just a primary key
type Tag struct {
	bun.BaseModel `bun:"table:tag,alias:t"`
	orm.Model

	Name string `json:"name" bun:"name,notnull"`
}

// Append-only: PK + creation audit
type ActionLog struct {
	bun.BaseModel `bun:"table:apv_action_log,alias:aal"`
	orm.Model
	orm.CreationTrackedModel

	InstanceID string `json:"instanceId" bun:"instance_id"`
	Action     string `json:"action" bun:"action"`
}

// Standard mutable entity: PK + full audit
type Role struct {
	bun.BaseModel `bun:"table:sys_role,alias:sr"`
	orm.FullAuditedModel

	Name     string `json:"name" bun:"name,notnull"`
	IsActive bool   `json:"isActive" bun:"is_active"`
}

// Composite PK: audit fields but PK defined separately
type UserRole struct {
	bun.BaseModel `bun:"table:sys_user_role,alias:sur"`
	orm.Model
	orm.CreationTrackedModel

	UserID varchar `json:"userId" bun:"user_id,notnull"`
	RoleID string  `json:"roleId" bun:"role_id,notnull"`
}
```

Typical choices:

- `orm.Model`: reference tables, join tables, or records with no audit columns
- `orm.Model` + `orm.CreationTrackedModel`: append-only records, snapshots, logs, outbox tables
- `orm.FullAuditedModel`: standard mutable entities that track both creation and update metadata
- `orm.FullTrackedModel`: entities with composite primary keys that still want full audit tracking

## Bun Model Hooks

VEF re-exports Bun's model lifecycle hook interfaces through `orm`:

| Hook Interface | When Called |
| --- | --- |
| `orm.BeforeSelectHook` | Before a SELECT query executes |
| `orm.AfterSelectHook` | After a SELECT query executes |
| `orm.BeforeInsertHook` | Before an INSERT query executes |
| `orm.AfterInsertHook` | After an INSERT query executes |
| `orm.BeforeUpdateHook` | Before an UPDATE query executes |
| `orm.AfterUpdateHook` | After an UPDATE query executes |
| `orm.BeforeDeleteHook` | Before a DELETE query executes |
| `orm.AfterDeleteHook` | After a DELETE query executes |
| `orm.BeforeScanRowHook` | Before scanning a row |
| `orm.AfterScanRowHook` | After scanning a row |

Implement any of these on your model struct to add lifecycle behavior:

```go
func (u *User) BeforeInsert(ctx context.Context, query *orm.BunInsertQuery) error {
	// Set defaults, validate, or log before insert
	return nil
}
```

## Audit Fields

The framework standardizes these common audit columns:

| Column | JSON Name | Persisted | Purpose |
| --- | --- | --- | --- |
| `id` | `id` | ✅ | Primary key |
| `created_at` | `createdAt` | ✅ | Creation timestamp |
| `created_by` | `createdBy` | ✅ | Creator user ID |
| `created_by_name` | `createdByName` | ❌ scanonly | Creator display name (populated by mold or JOIN) |
| `updated_at` | `updatedAt` | ✅ | Last update timestamp |
| `updated_by` | `updatedBy` | ✅ | Last updater user ID |
| `updated_by_name` | `updatedByName` | ❌ scanonly | Updater display name (populated by mold or JOIN) |

Not every model embeds all of these fields. `orm.CreationTrackedModel` contributes the `created_*` subset. `orm.FullTrackedModel` and `orm.FullAuditedModel` contribute both subsets.

The framework also exports audit column and field name constants:

```go
orm.ColumnID            // "id"
orm.ColumnCreatedAt     // "created_at"
orm.ColumnUpdatedAt     // "updated_at"
orm.ColumnCreatedBy     // "created_by"
orm.ColumnUpdatedBy     // "updated_by"
orm.ColumnCreatedByName // "created_by_name"
orm.ColumnUpdatedByName // "updated_by_name"

orm.FieldID             // "ID"
orm.FieldCreatedAt      // "CreatedAt"
// ... and so on
```

System operator constants for `created_by` / `updated_by`:

```go
orm.OperatorSystem    // "system" — used by system initialization
orm.OperatorCronJob   // "cron_job" — used by scheduled tasks
orm.OperatorAnonymous // "anonymous" — used by unauthenticated operations
```

## Tags You Will Use Most Often

### `bun`

Controls table name, alias, primary key rules, null behavior, and relations.

### `json`

Controls request and response payload field names. In practice, VEF projects usually use camelCase JSON names even when database columns stay snake_case.

### `validate`

Used by automatic request validation when params or meta are decoded into structs.

### `label` / `label_i18n`

Used by the validator to generate readable field names in error messages.

### `search`

Used by the search parser and CRUD query builders to translate search payloads into SQL conditions.

### `meta`

Used by the storage promoter to detect uploaded file fields, rich text fields, and markdown fields that need temp-file promotion.

### `mold`

Used by the struct transformer for field-level data transformation. The most common built-in usage is `mold:"translate=user?"` on `*ByName` fields.

## Search Models Are Usually Separate

Do not overload your persistence model with search semantics. Instead, define a dedicated search struct:

```go
type UserSearch struct {
	api.P

	Keyword  string `json:"keyword" search:"contains,column=username|email"`
	IsActive *bool  `json:"isActive" search:"eq,column=is_active"`
}
```

This keeps search rules explicit and prevents your write model from becoming a query DSL.

## Pagination And Sorting Metadata

For paging endpoints, metadata often comes through dedicated structs such as:

```go
type UserSearch struct {
	api.P
	Keyword string `json:"keyword" search:"contains,column=username|email"`
}

type UserMeta struct {
	api.M
	page.Pageable
	crud.Sortable
}
```

`page.Pageable` and `crud.Sortable` are metadata helpers, not persistence models.

## Practical Advice

- keep database models small and explicit
- use dedicated params structs for writes
- use dedicated search structs for reads
- compose from smaller embedded base models when an entity does not need the full `orm.FullAuditedModel` shape
- embed `orm.FullAuditedModel` only when you want the framework's standard audit behavior
- keep database tags, validation tags, and search tags close to the fields they govern
- remember that `*ByName` fields are scan-only — they are never written to the database

## Next Step

Read [Generic CRUD](./crud) to see how these models plug into typed operation builders, or read the [ORM SQL Builder](./orm-builder) for a comprehensive reference on constructing SQL queries.
