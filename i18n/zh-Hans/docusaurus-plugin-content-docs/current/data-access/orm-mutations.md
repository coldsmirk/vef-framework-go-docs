---
sidebar_position: 7
---

# ORM：写入操作

使用查询构造器写入数据：INSERT、UPDATE、DELETE、原始 SQL 与软删除行为。这些操作的事务化执行方式见[事务](./transactions)。

## INSERT 子句

```go
// 插入单条记录
user := &User{Username: "alice", Email: "alice@example.com"}
_, err := db.NewInsert().Model(user).Exec(ctx)

// 使用 RETURNING 插入（PostgreSQL）
err := db.NewInsert().Model(user).ReturningAll().Scan(ctx)

// 仅插入指定列
_, err := db.NewInsert().Model(user).
	Select("username", "email").
	Exec(ctx)

// 排除列
_, err := db.NewInsert().Model(user).
	Exclude("password").
	Exec(ctx)

// 显式设置列值
_, err := db.NewInsert().Model(user).
	Column("status", "active").
	ColumnExpr("score", func(eb orm.ExprBuilder) any {
		return eb.Literal(100)
	}).
	Exec(ctx)
```

### ON CONFLICT（Upsert）

```go
// ON CONFLICT (username) DO UPDATE SET email = EXCLUDED.email
_, err := db.NewInsert().Model(user).
	OnConflict(func(cb orm.ConflictBuilder) {
		cb.Columns("username").
			DoUpdate().
			Set("email", user.Email)
	}).Exec(ctx)

// ON CONFLICT DO NOTHING
_, err := db.NewInsert().Model(user).
	OnConflict(func(cb orm.ConflictBuilder) {
		cb.Columns("username").DoNothing()
	}).Exec(ctx)
```

## UPDATE 子句

```go
// 通过主键更新模型
user.Email = "new@example.com"
_, err := db.NewUpdate().Model(user).WherePK().Exec(ctx)

// 更新指定列
_, err := db.NewUpdate().Model(user).
	Select("email", "updated_at").
	WherePK().Exec(ctx)

// 显式设置值
_, err := db.NewUpdate().Model((*User)(nil)).
	Set("status", "inactive").
	SetExpr("updated_at", func(eb orm.ExprBuilder) any {
		return eb.Now()
	}).
	Where(func(cb orm.ConditionBuilder) {
		cb.Equals("status", "active").
			CreatedAtLessThan(cutoffTime)
	}).Exec(ctx)

// 忽略零值
_, err := db.NewUpdate().Model(user).OmitZero().WherePK().Exec(ctx)

// 批量更新
_, err := db.NewUpdate().Model(&users).Bulk().Exec(ctx)

// 带 RETURNING 的更新
err := db.NewUpdate().Model(user).WherePK().ReturningAll().Scan(ctx)
```

> 框架会自动将 `created_at` 和 `created_by` 从 UPDATE 操作中排除，以保护创建审计数据。

## DELETE 子句

```go
// 通过主键删除
_, err := db.NewDelete().Model(user).WherePK().Exec(ctx)

// 条件删除
_, err := db.NewDelete().Model((*User)(nil)).
	Where(func(cb orm.ConditionBuilder) {
		cb.Equals("status", "deactivated").
			CreatedAtLessThan(oneYearAgo)
	}).Exec(ctx)

// 强制删除（跳过软删除）
_, err := db.NewDelete().Model(user).WherePK().ForceDelete().Exec(ctx)

// 带 RETURNING 的删除
err := db.NewDelete().Model(user).WherePK().ReturningAll().Scan(ctx)
```

## 原始 SQL 查询

```go
// 带参数绑定的原始 SQL
var result []MyStruct
db.NewRaw("SELECT * FROM users WHERE status = ?", "active").Scan(ctx, &result)
```

## 软删除支持

```go
// 仅查询已软删除的记录
db.NewSelect().Model(&users).WhereDeleted().Scan(ctx)

// 包含已软删除的记录
db.NewSelect().Model(&users).IncludeDeleted().Scan(ctx)
```

## 下一步

- [事务](./transactions) — 使用 `RunInTx` 原子化地执行写入操作
- [ORM：DDL 与 Surface Map](./orm-ddl) — DDL 构造器与 `orm` 包完整公开接口面
