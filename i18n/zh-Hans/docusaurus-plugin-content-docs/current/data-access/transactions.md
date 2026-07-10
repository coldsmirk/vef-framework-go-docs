---
sidebar_position: 9
---

# 事务

VEF 通过 `orm.DB` 提供事务能力，不少 CRUD 写操作内部本身就已经在使用事务。

## 主要事务 API

公开的入口是：

- `RunInTx`
- `RunInReadOnlyTx`
- `BeginTx`

> v0.24 把原来大小写为 `RunInTX` / `RunInReadOnlyTX` 的方法改名为 `RunInTx` / `RunInReadOnlyTx`，以保持与框架其他地方一致的大小写风格。

最常见的用法是：

```go
db.RunInTx(ctx, func(ctx context.Context, tx orm.DB) error {
  return nil
})
```

## 完整示例
```go
// 自动事务（推荐）
err := db.RunInTx(ctx, func(ctx context.Context, tx orm.DB) error {
	_, err := tx.NewInsert().Model(order).Exec(ctx)
	if err != nil {
		return err // 自动回滚
	}

	_, err = tx.NewUpdate().Model((*Inventory)(nil)).
		Set("quantity", newQty).
		Where(func(cb orm.ConditionBuilder) {
			cb.PKEquals(itemID)
		}).Exec(ctx)

	return err // 返回 nil 则自动提交
})

// 只读事务
err := db.RunInReadOnlyTx(ctx, func(ctx context.Context, tx orm.DB) error {
	return tx.NewSelect().Model(&report).Scan(ctx)
})

// 手动事务
tx, err := db.BeginTx(ctx, nil)
if err != nil {
	return err
}
defer tx.Rollback()

// ... 使用 tx 执行操作 ...

return tx.Commit()
```

## CRUD 自动做了什么

Create、update、delete、import，以及若干批量变更操作，内部已经使用了 `RunInTx(...)`。

也就是说，除非你在更高的编排层扩展行为，一般**不需要**再给一次泛型 CRUD 写操作额外包一层事务。

## 事务回调里能拿到什么

在事务回调内部，`tx` 仍然是一个 `orm.DB`，所以你使用的查询构造 API 完全一致：

- `NewSelect`
- `NewInsert`
- `NewUpdate`
- `NewDelete`
- `NewMerge`

这让事务内的代码保持可预期，并与框架其余部分风格一致。

## 事务内发布事件

要让事件和业务写入原子地一起生效，在事务回调内部发布事件，并通过 `event.WithTx` 把事务交给 bus：

```go
err := db.RunInTx(ctx, func(ctx context.Context, tx orm.DB) error {
	if _, err := tx.NewInsert().Model(order).Exec(ctx); err != nil {
		return err
	}

	return bus.Publish(ctx, orderCreated, event.WithTx(tx))
})
```

- 事务句柄是**显式**传入的 —— bus 不会从 `ctx` 里读取事务。要传回调里拿到的 `tx`；如果传了外层的 `db`，事件会被悄悄写到你的事务之外。
- 使用 `WithTx` 时，bus 会把路由收窄到具有 `Transactional` capability 的 transport（实践中就是 outbox transport），把事件作为一行记录写进 `sys_event_outbox`，且写入发生在你的事务内。commit 之后 relay 才会把它转发给 sink transport；rollback 时该行随其他改动一起消失。
- 如果事件类型的路由里没有事务性 transport，`Publish` 会以 `event.ErrTxRequired` 失败。依赖这一模式的模块可以在启动时用 `event.RouteInspector.HasTransactionalRoute` 提前断言路由。
- `event.WithTx` 和 `event.WithAsync` 互斥（`event.ErrTxAsyncMutex`）：事务性发布必须在事务提交之前完成。
- 事务必须开在主数据源上 —— outbox 表就在那里，见[多数据源](./datasources)。

Outbox 配置、relay 重试和 DLQ 行为属于 transport 层面的内容 —— 见[事件总线](../infrastructure/event-bus)。

## 只读事务

如果读流程只需要一致性、不涉及写意图，用 `RunInReadOnlyTx(...)`。

## 手动事务

如果需要更底层的控制，可以使用 `BeginTx(...)`，它返回的事务支持显式的 `Commit` 和 `Rollback`。

只有当回调式事务不够用时才应该使用这种方式。

## 隔离级别与选项

`RunInTx` 和 `RunInReadOnlyTx` 不接受任何选项 —— 两者都固定运行在 `READ COMMITTED` 隔离级别，`RunInReadOnlyTx` 额外把事务标记为只读。这两个辅助方法都没有带选项的变体。

需要其他隔离级别时，用 `BeginTx(ctx, opts)`，它接受标准库的 `*sql.TxOptions`：

```go
tx, err := db.BeginTx(ctx, &sql.TxOptions{
	Isolation: sql.LevelSerializable,
})
```

传 `nil` 用的是驱动的默认隔离级别 —— 而不是 `READ COMMITTED` —— 所以依赖特定隔离级别时要显式传选项。

## 嵌套事务

在事务作用域的 `orm.DB` 上再调用 `RunInTx`（或 `BeginTx`）不会开启第二个数据库事务，而是在当前事务内创建一个**保存点（savepoint）**：

- 内层回调运行在一个 `SAVEPOINT` 里；返回错误时只回滚到保存点。错误仍会传播给外层回调，由外层决定整个事务是继续还是中止。
- 内层回调返回 `nil` 时释放保存点。在最外层事务提交之前，其他连接看不到任何改动。
- 嵌套调用会忽略事务选项：保存点继承外层事务的隔离级别，嵌套的 `RunInReadOnlyTx` 也不会让内层作用域变成只读。

这正是内部用 `RunInTx` 包裹自身的泛型 CRUD 写操作可以安全地在你自己的事务里调用的原因：它们通过保存点加入你的事务，而不是提前提交。

## Context 取消

传给 `RunInTx` 或 `BeginTx` 的 context 管辖整个事务。如果它在提交前被取消（或超过 deadline），`database/sql` 会回滚事务，进行中的查询以 context 错误失败，`RunInTx` 把这个错误返回给调用方。

## 后台代码中的事务

Cron job、事件订阅者等运行在 HTTP 请求之外的代码没有请求 context，所以 `contextx.DB(ctx)` 在那里返回 `nil`（见[扩展 Handler 参数](../advanced/extending-parameters)）。应改用依赖注入获取 `orm.DB`：任何构造函数或 `vef.Invoke` 函数都可以声明一个 `orm.DB` 参数，拿到的就是主数据源。

```go
vef.Invoke(func(scheduler cron.Scheduler, db orm.DB) error {
	_, err := scheduler.NewJob(cron.NewCronJob("0 3 * * *", false,
		cron.WithName("nightly-rollup"),
		cron.WithTask(func(ctx context.Context) error {
			return db.RunInTx(ctx, func(ctx context.Context, tx orm.DB) error {
				// ... 批处理工作 ...
				return nil
			})
		}),
	))

	return err
})
```

事件 handler 同理：把注入的 `orm.DB` 捕获进你注册给 `event.SubscribeTyped` 的闭包里。注册模式见 [Cron Jobs](../infrastructure/cron) 和[事件总线](../infrastructure/event-bus)。
