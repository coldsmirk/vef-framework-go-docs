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
	orm.Model

	Username string `json:"username" validate:"required,alphanum,max=32" label:"Username"`
	Email    string `json:"email" validate:"omitempty,email,max=128" label:"Email"`
	IsActive bool   `json:"isActive"`
}
```

This combines two different concerns:

- `bun.BaseModel`: table metadata for Bun
- `orm.Model`: framework-owned base fields such as ID and audit columns

## Base Model Types

VEF exposes a few reusable model shapes through `orm`:

- `orm.Model`: ID plus created/updated audit fields
- `orm.IDModel`: ID only
- `orm.CreatedModel`: creation fields only
- `orm.AuditedModel`: created/updated audit fields without the primary key

Use the smallest one that matches the lifecycle of your entity.

These types are meant to be anonymously embedded as reusable field slices, not just picked as mutually exclusive top-level bases. `orm.Model` is the convenience form for the common case, while the smaller types let you compose only the fields your entity actually needs.

A practical way to think about them:

- `orm.IDModel`: adds only the primary key
- `orm.CreatedModel`: adds only the `created_*` tracking fields
- `orm.AuditedModel`: adds both `created_*` and `updated_*` tracking fields, but no primary key
- `orm.Model`: adds the primary key plus the same audit fields as `orm.AuditedModel`

Conceptually, `orm.Model` is the pre-composed form you would reach for instead of embedding both `orm.IDModel` and `orm.AuditedModel`.

## Embedding And Composition

The smaller base model types are useful when an entity does not need the full `orm.Model` field set.

```go
type Tag struct {
	bun.BaseModel `bun:"table:tag,alias:t"`
	orm.IDModel

	Name string `json:"name" bun:"name,notnull"`
}

type EventOutbox struct {
	bun.BaseModel `bun:"table:event_outbox,alias:eo"`
	orm.IDModel
	orm.CreatedModel

	EventType string `json:"eventType" bun:"event_type,notnull"`
}

type Delegation struct {
	bun.BaseModel `bun:"table:delegation,alias:d"`
	orm.Model

	DelegatorID string `json:"delegatorId" bun:"delegator_id,notnull"`
}
```

Typical choices:

- `orm.IDModel`: reference tables, join tables, or other records with no audit columns
- `orm.IDModel` + `orm.CreatedModel`: append-only records, snapshots, logs, and outbox-style tables
- `orm.Model`: standard mutable entities that track both creation and update metadata
- `orm.AuditedModel`: entities that already define their own primary key field but still want the framework's standard audit columns

## Audit Fields

The framework standardizes common audit columns:

- `id`
- `created_at`
- `created_by`
- `created_by_name`
- `updated_at`
- `updated_by`
- `updated_by_name`

Not every model embeds all of these fields. `orm.CreatedModel` contributes the `created_*` subset, while `orm.AuditedModel` and `orm.Model` contribute both the `created_*` and `updated_*` subsets.

When you embed `orm.Model`, VEF-aware operations and request context can fill and use those fields consistently.

One important detail: `created_by_name` and `updated_by_name` are scan-only convenience fields in the current model definitions. Treat them as read-side fields, not as columns that the framework persists directly.

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
- compose from smaller embedded base models when an entity does not need the full `orm.Model` shape
- embed `orm.Model` only when you want the framework's standard audit behavior
- keep database tags, validation tags, and search tags close to the fields they govern

## Next Step

Read [Generic CRUD](./crud) to see how these models plug into typed operation builders.
