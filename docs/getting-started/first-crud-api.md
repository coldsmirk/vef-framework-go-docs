---
sidebar_position: 3
title: Your First CRUD API
---

# Your First CRUD API

The [Quick Start](./quick-start.md) served a hand-written handler. This tutorial builds the next thing you will actually need: a complete CRUD API for a `Product` entity — model, table, typed request params, generic CRUD operations, and one customization hook — all verifiable with curl.

## What you will build

- an `app/product` RPC resource exposing `create`, `update`, `delete`, and `find_page`
- a `app_product` table that follows the [database conventions](../conventions/database-conventions.md)
- keyword and status filtering plus pagination, driven by `search` tags
- a pre-create hook that rejects duplicate product codes

## Prerequisites

- a working application from the [Quick Start](./quick-start.md)
- the `sqlite3` command-line tool (the tutorial keeps the SQLite setup from the quick start)

The finished layout looks like this:

```text
my-app/
├── configs/
│   └── application.toml
├── data/
│   └── app.db
├── db/
│   └── app_product.sql
├── internal/
│   └── product/
│       ├── model.go
│       ├── payload.go
│       ├── resource.go
│       └── module.go
└── main.go
```

## 1. Define the model

Create `internal/product/model.go`:

```go
package product

import (
	"github.com/uptrace/bun"

	"github.com/coldsmirk/vef-framework-go/orm"
)

type Product struct {
	bun.BaseModel `bun:"table:app_product,alias:ap"`
	orm.FullAuditedModel

	Name     string `json:"name" bun:"name,notnull"`
	Code     string `json:"code" bun:"code,notnull"`
	Stock    int    `json:"stock" bun:"stock,notnull"`
	IsActive bool   `json:"isActive" bun:"is_active,notnull"`
	Remark   string `json:"remark" bun:"remark"`
}
```

Two embedded types do the heavy lifting:

- `bun.BaseModel` binds the struct to the `app_product` table with alias `ap`
- `orm.FullAuditedModel` contributes `ID`, `CreatedAt`, `CreatedBy`, `CreatedByName`, `UpdatedAt`, `UpdatedBy`, and `UpdatedByName`

You never assign these framework-owned fields yourself. On insert, the framework generates a compact string ID for an empty string primary key and fills `created_at` / `created_by`; on update it maintains `updated_at` / `updated_by` using the current principal. The base model catalog is covered in [Models](../data-access/models.md).

## 2. Create the table

VEF does not generate schema from models. Application projects own their DDL scripts, following the [database conventions](../conventions/database-conventions.md): a module prefix in the table name (`app_` here), fixed audit columns, and named constraints.

Create `db/app_product.sql`:

```sql
BEGIN;

CREATE TABLE IF NOT EXISTS app_product (
    id                       VARCHAR(32) NOT NULL,
    created_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by               VARCHAR(32) NOT NULL DEFAULT 'system',
    updated_by               VARCHAR(32) NOT NULL DEFAULT 'system',
    name                     VARCHAR(32) NOT NULL,
    code                     VARCHAR(32) NOT NULL,
    stock                    INTEGER NOT NULL DEFAULT 0,
    is_active                BOOLEAN NOT NULL DEFAULT FALSE,
    remark                   VARCHAR(512),

    CONSTRAINT pk_app_product PRIMARY KEY (id),
    CONSTRAINT uk_app_product__code UNIQUE (code)
);

COMMIT;
```

Apply it to the SQLite database file the app will use:

```bash
mkdir -p data
sqlite3 data/app.db < db/app_product.sql
```

This script is intentionally the portable subset. In a real PostgreSQL project the conventions additionally require `LOCALTIMESTAMP` defaults, `COMMENT ON` statements for every table and column, and `created_by` / `updated_by` foreign keys to `sys_user(id)` — see the full templates in [Database Conventions](../conventions/database-conventions.md).

## 3. Define write params and search params

Persistence models should not double as request contracts. Create `internal/product/payload.go` with one struct for writes and one for reads:

```go
package product

import (
	"github.com/coldsmirk/vef-framework-go/api"
)

type ProductParams struct {
	api.P

	ID       string `json:"id"`
	Name     string `json:"name" validate:"required,max=32" label:"Name"`
	Code     string `json:"code" validate:"required,max=32" label:"Code"`
	Stock    int    `json:"stock" validate:"gte=0" label:"Stock"`
	IsActive *bool  `json:"isActive"`
	Remark   string `json:"remark" validate:"max=512" label:"Remark"`
}

type ProductSearch struct {
	api.P

	Keyword  *string `json:"keyword" search:"contains,column=name|code"`
	IsActive *bool   `json:"isActive" search:"eq"`
	MinStock *int    `json:"minStock" search:"gte,column=stock"`
}
```

What each piece does:

- the embedded `api.P` sentinel tells the framework to decode this struct from the request's `params` field and validate it
- `validate` tags run automatically before your operation executes; `label` names the field in error messages
- `search` tags translate directly into `WHERE` clauses: `keyword` becomes a `LIKE` match across `name` OR `code`, `minStock` becomes `stock >= ?`
- `ID` stays empty on create (the framework generates one) and is required on update
- pointer fields distinguish "not provided" from a zero value — update merges non-empty fields only, so `IsActive *bool` is what lets a client explicitly send `false`

## 4. Assemble the API resource

Create `internal/product/resource.go`:

```go
package product

import (
	"github.com/coldsmirk/vef-framework-go/api"
	"github.com/coldsmirk/vef-framework-go/crud"
)

type ProductResource struct {
	api.Resource

	crud.FindPage[Product, ProductSearch]
	crud.Create[Product, ProductParams]
	crud.Update[Product, ProductParams]
	crud.Delete[Product]
}

func NewProductResource() api.Resource {
	return &ProductResource{
		Resource: api.NewRPCResource("app/product"),
		FindPage: crud.NewFindPage[Product, ProductSearch]().Public(),
		Create:   crud.NewCreate[Product, ProductParams]().Public(),
		Update:   crud.NewUpdate[Product, ProductParams]().Public(),
		Delete:   crud.NewDelete[Product]().Public(),
	}
}
```

Each embedded builder implements `api.OperationsProvider`, so the framework collects them automatically and registers one operation per builder:

| Embedded builder | Default action | Behavior |
| --- | --- | --- |
| `crud.FindPage[Product, ProductSearch]` | `find_page` | filtered, paginated list with total count |
| `crud.Create[Product, ProductParams]` | `create` | copies params into a model, inserts it in a transaction |
| `crud.Update[Product, ProductParams]` | `update` | loads the record by `id`, merges non-empty fields, updates |
| `crud.Delete[Product]` | `delete` | loads the record by `id`, deletes it in a transaction |

`Public()` keeps the tutorial runnable without an auth provider. In a real application, drop it and protect each operation instead:

```go
crud.NewCreate[Product, ProductParams]().RequiredPermission("app:product:create")
```

## 5. Register the module and run

Create `internal/product/module.go`:

```go
package product

import (
	"github.com/coldsmirk/vef-framework-go"
)

var Module = vef.Module(
	"app:product",
	vef.ProvideAPIResource(NewProductResource),
)
```

Compose it in `main.go`:

```go
package main

import (
	"github.com/coldsmirk/vef-framework-go"

	"example.com/my-app/internal/product"
)

func main() {
	vef.Run(
		product.Module,
	)
}
```

Point `configs/application.toml` at the database file from step 2:

```toml
[vef.app]
name = "my-app"
port = 8080

[vef.data_sources.primary]
type = "sqlite"
path = "data/app.db"

[vef.event.transports.outbox]
enabled = true

[[vef.event.routing]]
pattern = "vef.storage.*"
transports = ["outbox"]
```

The last two blocks are new compared to the quick start. The generic write operations run the [file storage](../infrastructure/storage.md) lifecycle inside their transactions, and storage publishes its domain events through a transactional transport — the framework fails fast at startup if `vef.storage.*` events have no such route. Enabling the outbox transport (it creates its own table automatically) and routing storage events to it satisfies the check.

Start the app:

```bash
go run .
```

## 6. Call the API

All four operations go through the same RPC endpoint, `POST /api`, selected by the `resource` / `action` / `version` envelope fields.

### Create

```bash
curl http://localhost:8080/api \
  -H 'Content-Type: application/json' \
  -d '{
    "resource": "app/product",
    "action": "create",
    "version": "v1",
    "params": {
      "name": "Espresso Beans",
      "code": "P-1001",
      "stock": 20,
      "isActive": true
    }
  }'
```

The response returns the generated primary key:

```json
{
  "code": 0,
  "message": "新增成功",
  "data": {
    "id": "d1nbkq2s7kg5jkvvs7lg"
  }
}
```

As in the quick start, messages follow the framework's default language; set `VEF_I18N_LANGUAGE=en` to get `Created successfully` instead.

### Find page

Filters come from `params` (your `ProductSearch`), pagination from `meta`:

```bash
curl http://localhost:8080/api \
  -H 'Content-Type: application/json' \
  -d '{
    "resource": "app/product",
    "action": "find_page",
    "version": "v1",
    "params": { "keyword": "Espresso" },
    "meta": { "page": 1, "size": 10 }
  }'
```

The `data` payload is a page object:

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "page": 1,
    "size": 10,
    "total": 1,
    "items": [
      {
        "id": "d1nbkq2s7kg5jkvvs7lg",
        "createdAt": "2026-07-09 10:30:00",
        "createdBy": "anonymous",
        "createdByName": "",
        "updatedAt": "2026-07-09 10:30:00",
        "updatedBy": "anonymous",
        "updatedByName": "",
        "name": "Espresso Beans",
        "code": "P-1001",
        "stock": 20,
        "isActive": true,
        "remark": ""
      }
    ]
  }
}
```

Note `createdBy: "anonymous"`: the audit columns were filled automatically, and because the operation is public there is no authenticated principal yet.

### Update

Send the `id` plus the fields to change; unset fields keep their stored values:

```bash
curl http://localhost:8080/api \
  -H 'Content-Type: application/json' \
  -d '{
    "resource": "app/product",
    "action": "update",
    "version": "v1",
    "params": {
      "id": "d1nbkq2s7kg5jkvvs7lg",
      "name": "Espresso Beans",
      "code": "P-1001",
      "stock": 35
    }
  }'
```

```json
{ "code": 0, "message": "保存成功", "data": null }
```

### Delete

```bash
curl http://localhost:8080/api \
  -H 'Content-Type: application/json' \
  -d '{
    "resource": "app/product",
    "action": "delete",
    "version": "v1",
    "params": { "id": "d1nbkq2s7kg5jkvvs7lg" }
  }'
```

```json
{ "code": 0, "message": "删除成功", "data": null }
```

The envelope fields and transport rules are specified in [Routing](../building-apis/routing.md), and the `code` / `message` / `data` contract in [Results and Errors](../building-apis/results-and-errors.md).

## 7. Add a create hook

The generic builders accept hooks that run inside the same transaction as the write. Enforce the unique product code at the API level with `WithPreCreate`, so callers get a structured business error instead of a raw constraint violation.

Update the `Create` builder in `internal/product/resource.go`:

```go
Create: crud.NewCreate[Product, ProductParams]().
	Public().
	WithPreCreate(func(model *Product, params *ProductParams, query orm.InsertQuery, ctx fiber.Ctx, tx orm.DB) error {
		exists, err := tx.NewSelect().Model((*Product)(nil)).
			Where(func(cb orm.ConditionBuilder) { cb.Equals("code", model.Code) }).
			Exists(ctx.Context())
		if err != nil {
			return err
		}
		if exists {
			return result.Err("product code already exists",
				result.WithCode(result.ErrCodeRecordAlreadyExists))
		}
		return nil
	}),
```

Add the new imports to the file: `github.com/gofiber/fiber/v3`, `github.com/coldsmirk/vef-framework-go/orm`, and `github.com/coldsmirk/vef-framework-go/result`.

Restart the app and replay the create request from step 6 twice. The second call now fails cleanly:

```json
{
  "code": 2002,
  "message": "product code already exists",
  "data": null
}
```

The hook runs before the insert, inside the transaction, with the pending model, the decoded params, the insert query, and the transactional `orm.DB`. Every builder has a matching pair — `WithPreUpdate` / `WithPostUpdate`, `WithPreDelete` / `WithPostDelete`, and more — cataloged in [Generic CRUD](../data-access/crud.md).

## What the framework did for you

You wrote a model, two request structs, and a resource. You did not write SQL for the four operations, request decoding, validation wiring, transactions, ID generation, audit-column maintenance, pagination counting, or response envelopes.

## Where to go next

- [Generic CRUD](../data-access/crud.md): every builder, option, and hook in the `crud` package
- [API Resources](../building-apis/api.md): operations, auth config, and custom actions beyond CRUD
- [Models](../data-access/models.md): base model types, tags, and search struct patterns
