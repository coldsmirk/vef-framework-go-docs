---
sidebar_position: 2
---

# 事务

VEF 通过 `orm.DB` 提供事务能力。对大多数用户来说，`RunInTx` 就已经够用了，而且不少 CRUD 写操作本身就已经在事务里执行。

> v0.24 起，`RunInTX` / `RunInReadOnlyTX` 改名为 `RunInTx` / `RunInReadOnlyTx`，保持与框架其他地方一致的大小写。

## 默认模式

当你需要多次写操作一起成功或一起失败时，使用 `RunInTx`：

```go
import (
  "context"

  "github.com/coldsmirk/vef-framework-go/orm"
)

func SaveAll(ctx context.Context, db orm.DB) error {
  return db.RunInTx(ctx, func(txCtx context.Context, tx orm.DB) error {
    if _, err := tx.NewInsert().Model(&User{}).Exec(txCtx); err != nil {
      return err
    }

    if _, err := tx.NewInsert().Model(&Profile{}).Exec(txCtx); err != nil {
      return err
    }

    return nil
  })
}
```

回调返回 error 就回滚，返回 `nil` 就提交。

## 只读事务

如果你只想要只读事务，可以用：

```go
db.RunInReadOnlyTx(...)
```

## 手动事务

你也可以使用 `BeginTx(...)` 自己控制 `Commit` / `Rollback`，但这应该是少数情况，不应作为默认起手式。

## CRUD 自带事务

当前不少泛型 CRUD 写路径内部已经包了 `RunInTx`，包括：

- create
- update
- delete
- import
- 多条写操作的 batch 变体

所以在这些 handler 外面再套一层事务之前，最好先明确你真正想要的数据库行为。

## 请求级 DB

在 API 请求处理中，框架注入的 `orm.DB` 已经带有当前 operator 上下文。进入事务回调后，应该继续使用回调参数里的 `tx orm.DB`，而不是再退回外层的 DB。
