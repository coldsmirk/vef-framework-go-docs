---
sidebar_position: 4
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

The exact processor type declarations — single-record, batch, and
export/import — are maintained in one place: the
[Generic CRUD processor reference](./crud#processor-type-signatures).
This page focuses on when to use which hook and what runs inside the
transaction.

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

## CRUD Hooks And File Lifecycle

Create, update, and delete builders integrate with the `storage.Files` / `FilesFor[T]` lifecycle facade (the replacement for the old `Promoter[T]`).

That means:

- `meta`-tagged file fields are reconciled inside the same transaction as the business write
- update reconciliation enqueues asynchronous deletes for replaced file values
- delete reconciliation enqueues every referenced file for asynchronous removal
- on transaction rollback no claim is consumed and no deletion is enqueued — there is nothing to "restore"

Hooks and file lifecycle therefore share one transactional lifecycle; the actual backend deletion happens asynchronously through the storage delete worker (see [Storage](../infrastructure/storage)).

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

Read [Validation](../building-apis/validation) and [Error Handling](../building-apis/results-and-errors) to see how request failures and business errors are surfaced to clients.
