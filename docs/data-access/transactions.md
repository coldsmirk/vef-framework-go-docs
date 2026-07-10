---
sidebar_position: 9
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

## Worked example
```go
// Automatic transaction (recommended)
err := db.RunInTx(ctx, func(ctx context.Context, tx orm.DB) error {
	_, err := tx.NewInsert().Model(order).Exec(ctx)
	if err != nil {
		return err // auto rollback
	}

	_, err = tx.NewUpdate().Model((*Inventory)(nil)).
		Set("quantity", newQty).
		Where(func(cb orm.ConditionBuilder) {
			cb.PKEquals(itemID)
		}).Exec(ctx)

	return err // auto commit if nil
})

// Read-only transaction
err := db.RunInReadOnlyTx(ctx, func(ctx context.Context, tx orm.DB) error {
	return tx.NewSelect().Model(&report).Scan(ctx)
})

// Manual transaction
tx, err := db.BeginTx(ctx, nil)
if err != nil {
	return err
}
defer tx.Rollback()

// ... operations with tx ...

return tx.Commit()
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

## Transactional event publishing

To publish an event atomically with a business write, publish inside the transaction callback and hand the transaction to the bus with `event.WithTx`:

```go
err := db.RunInTx(ctx, func(ctx context.Context, tx orm.DB) error {
	if _, err := tx.NewInsert().Model(order).Exec(ctx); err != nil {
		return err
	}

	return bus.Publish(ctx, orderCreated, event.WithTx(tx))
})
```

- The transaction handle is passed **explicitly** — the bus does not read it from `ctx`. Pass the `tx` you received in the callback; passing the outer `db` would silently write the event outside your transaction.
- With `WithTx`, the bus narrows the route to transports with the `Transactional` capability — in practice the outbox transport — and stores the event as a row in `sys_event_outbox` within your transaction. After commit, the relay forwards it to the sink transport; on rollback, the row disappears along with everything else.
- If the event type routes to no transactional transport, `Publish` fails with `event.ErrTxRequired`. Modules that rely on this pattern can assert their routes at startup via `event.RouteInspector.HasTransactionalRoute`.
- `event.WithTx` and `event.WithAsync` are mutually exclusive (`event.ErrTxAsyncMutex`): a transactional publish must complete before the transaction commits.
- The transaction must come from the primary data source, where the outbox table lives — see [Multiple Data Sources](./datasources).

Outbox configuration, relay retries, and DLQ behavior are transport concerns — see [Event Bus](../infrastructure/event-bus).

## Read-only transactions

When you want consistency for read flows without write intent, use `RunInReadOnlyTx(...)`.

## Manual transactions

If you need lower-level control, `BeginTx(...)` is available and returns a transaction that supports explicit `Commit` and `Rollback`.

Use this only when callback-based transactions are not enough.

## Isolation levels and options

`RunInTx` and `RunInReadOnlyTx` take no options — both run at `READ COMMITTED` isolation, with `RunInReadOnlyTx` additionally marking the transaction read-only. Neither helper has an options variant.

When you need a different isolation level, use `BeginTx(ctx, opts)`. It accepts a standard-library `*sql.TxOptions`:

```go
tx, err := db.BeginTx(ctx, &sql.TxOptions{
	Isolation: sql.LevelSerializable,
})
```

Passing `nil` uses the driver's default isolation level — not `READ COMMITTED` — so pass explicit options when you depend on a specific level.

## Nested transactions

Calling `RunInTx` (or `BeginTx`) on a transaction-scoped `orm.DB` does not open a second database transaction. It creates a **savepoint** inside the current one:

- The inner callback runs inside a `SAVEPOINT`; if it returns an error, only the savepoint is rolled back. The error still propagates to the outer callback, which decides whether the whole transaction continues or aborts.
- If the inner callback returns `nil`, the savepoint is released. Nothing becomes visible to other connections until the outermost transaction commits.
- Transaction options are ignored on nested calls: the savepoint inherits the outer transaction's isolation level, and a nested `RunInReadOnlyTx` does not make the inner scope read-only.

This is why generic CRUD mutations — which wrap themselves in `RunInTx` — are safe to call from inside your own transaction: they join it through a savepoint instead of committing early.

## Context cancellation

The context you pass to `RunInTx` or `BeginTx` governs the whole transaction. If it is cancelled (or its deadline expires) before commit, `database/sql` rolls the transaction back, in-flight queries fail with the context error, and `RunInTx` returns that error to the caller.

## Transactions from background code

Cron jobs, event subscribers, and other code running outside an HTTP request have no request context, so `contextx.DB(ctx)` returns `nil` there (see [Extending Handler Parameters](../advanced/extending-parameters)). Get `orm.DB` through dependency injection instead: any constructor or `vef.Invoke` function can declare an `orm.DB` parameter and receives the primary data source.

```go
vef.Invoke(func(scheduler cron.Scheduler, db orm.DB) error {
	_, err := scheduler.NewJob(cron.NewCronJob("0 3 * * *", false,
		cron.WithName("nightly-rollup"),
		cron.WithTask(func(ctx context.Context) error {
			return db.RunInTx(ctx, func(ctx context.Context, tx orm.DB) error {
				// ... batch work ...
				return nil
			})
		}),
	))

	return err
})
```

The same applies to event handlers: capture the injected `orm.DB` in the closure you register with `event.SubscribeTyped`. See [Cron Jobs](../infrastructure/cron) and [Event Bus](../infrastructure/event-bus) for the registration patterns.
