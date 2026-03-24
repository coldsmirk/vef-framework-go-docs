---
sidebar_position: 7
---

# Hooks

VEF has two hook surfaces:

- CRUD operation hooks for framework-managed create, update, delete, import, and export flows
- Bun model hooks for lower-level ORM query lifecycle interception

They solve different problems and should not be treated as the same mechanism.

## Hook Families Overview

| Hook family | Where it runs | Scope | Typical use |
| --- | --- | --- | --- |
| CRUD hooks | inside CRUD builders | one CRUD endpoint | business invariants, transactional side effects, file promotion coordination |
| Bun model hooks | inside ORM model lifecycle | one model and query type | low-level query mutation, model lifecycle behavior, persistence-side checks |

## CRUD Hook Surface

CRUD builders expose these hook APIs:

| Operation | Pre hook | Post hook |
| --- | --- | --- |
| `Create` | `WithPreCreate(...)` | `WithPostCreate(...)` |
| `Update` | `WithPreUpdate(...)` | `WithPostUpdate(...)` |
| `Delete` | `WithPreDelete(...)` | `WithPostDelete(...)` |
| `CreateMany` | `WithPreCreateMany(...)` | `WithPostCreateMany(...)` |
| `UpdateMany` | `WithPreUpdateMany(...)` | `WithPostUpdateMany(...)` |
| `DeleteMany` | `WithPreDeleteMany(...)` | `WithPostDeleteMany(...)` |
| `Export` | `WithPreExport(...)` | no post-export hook |
| `Import` | `WithPreImport(...)` | `WithPostImport(...)` |

## CRUD Hook Signatures

### Single-record create/update/delete

| Hook | Signature summary |
| --- | --- |
| `PreCreate` | `func(model *TModel, params *TParams, query orm.InsertQuery, ctx fiber.Ctx, tx orm.DB) error` |
| `PostCreate` | `func(model *TModel, params *TParams, ctx fiber.Ctx, tx orm.DB) error` |
| `PreUpdate` | `func(oldModel, model *TModel, params *TParams, query orm.UpdateQuery, ctx fiber.Ctx, tx orm.DB) error` |
| `PostUpdate` | `func(oldModel, model *TModel, params *TParams, ctx fiber.Ctx, tx orm.DB) error` |
| `PreDelete` | `func(model *TModel, query orm.DeleteQuery, ctx fiber.Ctx, tx orm.DB) error` |
| `PostDelete` | `func(model *TModel, ctx fiber.Ctx, tx orm.DB) error` |

### Batch create/update/delete

| Hook | Signature summary |
| --- | --- |
| `PreCreateMany` | `func(models []TModel, paramsList []TParams, query orm.InsertQuery, ctx fiber.Ctx, tx orm.DB) error` |
| `PostCreateMany` | `func(models []TModel, paramsList []TParams, ctx fiber.Ctx, tx orm.DB) error` |
| `PreUpdateMany` | `func(oldModels, models []TModel, paramsList []TParams, query orm.UpdateQuery, ctx fiber.Ctx, tx orm.DB) error` |
| `PostUpdateMany` | `func(oldModels, models []TModel, paramsList []TParams, ctx fiber.Ctx, tx orm.DB) error` |
| `PreDeleteMany` | `func(models []TModel, query orm.DeleteQuery, ctx fiber.Ctx, tx orm.DB) error` |
| `PostDeleteMany` | `func(models []TModel, ctx fiber.Ctx, tx orm.DB) error` |

### Export and import

| Hook | Signature summary |
| --- | --- |
| `PreExport` | `func(models []TModel, search TSearch, ctx fiber.Ctx, db orm.DB) error` |
| `PreImport` | `func(models []TModel, query orm.InsertQuery, ctx fiber.Ctx, tx orm.DB) error` |
| `PostImport` | `func(models []TModel, ctx fiber.Ctx, tx orm.DB) error` |

## CRUD Transaction Boundary

The important rule is that CRUD write operations already run inside transactions. Your CRUD hook receives the current transaction-scoped `orm.DB`, so additional database work participates in the same transaction automatically.

Example:

```go
crud.NewCreate[User, UserParams]().
	WithPostCreate(func(model *User, params *UserParams, ctx fiber.Ctx, tx orm.DB) error {
		_, err := tx.NewInsert().Model(&AuditLog{
			UserID: model.ID,
			Action: "created",
		}).Exec(ctx.Context())
		return err
	})
```

## CRUD Hook Error Behavior

If a CRUD hook returns an error:

- the operation fails
- the surrounding transaction rolls back
- the framework returns the error through normal result handling

This makes CRUD hooks a good place for business invariants that must be enforced atomically.

## CRUD Hooks And File Promotion

Create, update, and delete builders also integrate with the storage promoter.

That means:

- file promotion happens within the same CRUD flow
- update rollback can restore replaced files
- delete cleanup can remove promoted files after successful deletion

Hooks and file promotion therefore share one transactional lifecycle, even though the storage cleanup itself may happen as a follow-up action inside that flow.

## Bun Model Hook Surface

At the ORM layer, VEF also exposes Bun hook interfaces:

| Hook interface | Trigger |
| --- | --- |
| `orm.BeforeSelectHook` | before `SELECT` |
| `orm.AfterSelectHook` | after `SELECT` |
| `orm.BeforeInsertHook` | before `INSERT` |
| `orm.AfterInsertHook` | after `INSERT` |
| `orm.BeforeUpdateHook` | before `UPDATE` |
| `orm.AfterUpdateHook` | after `UPDATE` |
| `orm.BeforeDeleteHook` | before `DELETE` |
| `orm.AfterDeleteHook` | after `DELETE` |

These hooks are implemented on model types and operate at the ORM lifecycle level, not at the API action level.

## When To Use CRUD Hooks

CRUD hooks are a good fit when:

- the business step belongs tightly to one CRUD action
- the public API should remain CRUD-shaped
- the extra behavior must share the same transaction
- the hook needs access to both params and model state

## When To Use Bun Model Hooks

Bun hooks are a better fit when:

- the behavior belongs to the model lifecycle itself
- the logic should apply outside the API layer too
- the hook needs to mutate or inspect the underlying Bun query
- the concern is persistence-oriented rather than endpoint-oriented

## When Not To Use Hooks

Hooks are a poor fit when:

- the operation is no longer conceptually CRUD
- the endpoint orchestrates multiple unrelated workflows
- the action semantics are clearer as an explicit command endpoint
- the behavior would become hard to understand because it is split across many hook registrations

In those cases, a custom handler is usually clearer.

## Practical Advice

- keep CRUD hook logic short and local to one operation
- keep Bun model hooks focused on persistence behavior
- use CRUD hooks for transactional business steps, not for unrelated side effects
- if you find yourself stacking many CRUD hooks, reconsider whether the resource needs a custom handler
- if the behavior should apply everywhere the model is used, prefer a model hook over an API hook

## Next Step

Read [Validation](./validation) and [Error Handling](./error-handling) to see how request failures and business errors are surfaced to clients.
