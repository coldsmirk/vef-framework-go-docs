---
sidebar_position: 7
---

# Transactions

VEF exposes transactions through `orm.DB`, and many CRUD write operations already use them internally.

## The main transaction API

The public entry points are:

- `RunInTx`
- `RunInReadOnlyTx`
- `BeginTx`

> v0.24 renamed the previously upper-cased `RunInTX` / `RunInReadOnlyTX` helpers to use `Tx` casing for consistency with the rest of the framework.

The most common one is:

```go
db.RunInTx(ctx, func(ctx context.Context, tx orm.DB) error {
  return nil
})
```

## What CRUD does automatically

Create, update, delete, import, and several batch mutation operations already use `RunInTx(...)` internally.

That means you usually do **not** need to wrap a generic CRUD mutation inside another transaction unless you are extending behavior at a higher orchestration layer.

## What you get inside the transaction

Inside the transaction callback, `tx` is still an `orm.DB`, so you keep the same query-building API:

- `NewSelect`
- `NewInsert`
- `NewUpdate`
- `NewDelete`
- `NewMerge`

This keeps transaction code predictable and consistent with the rest of the framework.

## Read-only transactions

When you want consistency for read flows without write intent, use `RunInReadOnlyTx(...)`.

## Manual transactions

If you need lower-level control, `BeginTx(...)` is available and returns a transaction that supports explicit `Commit` and `Rollback`.

Use this only when callback-based transactions are not enough.
