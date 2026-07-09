---
sidebar_position: 7
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

## 只读事务

如果读流程只需要一致性、不涉及写意图，用 `RunInReadOnlyTx(...)`。

## 手动事务

如果需要更底层的控制，可以使用 `BeginTx(...)`，它返回的事务支持显式的 `Commit` 和 `Rollback`。

只有当回调式事务不够用时才应该使用这种方式。
