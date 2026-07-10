---
sidebar_position: 7
---

# ORM: Mutations

Writing data with the query builder: INSERT, UPDATE, DELETE, raw queries, and soft-delete behavior. For transactional execution of these operations, see [Transactions](./transactions).

## INSERT Clause

```go
// Insert a single record
user := &User{Username: "alice", Email: "alice@example.com"}
_, err := db.NewInsert().Model(user).Exec(ctx)

// Insert with RETURNING (PostgreSQL)
err := db.NewInsert().Model(user).ReturningAll().Scan(ctx)

// Insert specific columns only
_, err := db.NewInsert().Model(user).
	Select("username", "email").
	Exec(ctx)

// Exclude columns
_, err := db.NewInsert().Model(user).
	Exclude("password").
	Exec(ctx)

// Set column values explicitly
_, err := db.NewInsert().Model(user).
	Column("status", "active").
	ColumnExpr("score", func(eb orm.ExprBuilder) any {
		return eb.Literal(100)
	}).
	Exec(ctx)
```

### ON CONFLICT (Upsert)

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

## UPDATE Clause

```go
// Update a model by PK
user.Email = "new@example.com"
_, err := db.NewUpdate().Model(user).WherePK().Exec(ctx)

// Update specific columns
_, err := db.NewUpdate().Model(user).
	Select("email", "updated_at").
	WherePK().Exec(ctx)

// Set values explicitly
_, err := db.NewUpdate().Model((*User)(nil)).
	Set("status", "inactive").
	SetExpr("updated_at", func(eb orm.ExprBuilder) any {
		return eb.Now()
	}).
	Where(func(cb orm.ConditionBuilder) {
		cb.Equals("status", "active").
			CreatedAtLessThan(cutoffTime)
	}).Exec(ctx)

// Omit zero values
_, err := db.NewUpdate().Model(user).OmitZero().WherePK().Exec(ctx)

// Bulk update
_, err := db.NewUpdate().Model(&users).Bulk().Exec(ctx)

// Update with RETURNING
err := db.NewUpdate().Model(user).WherePK().ReturningAll().Scan(ctx)
```

> The framework automatically excludes `created_at` and `created_by` from UPDATE operations to preserve creation audit data.

## DELETE Clause

```go
// Delete by PK
_, err := db.NewDelete().Model(user).WherePK().Exec(ctx)

// Delete with condition
_, err := db.NewDelete().Model((*User)(nil)).
	Where(func(cb orm.ConditionBuilder) {
		cb.Equals("status", "deactivated").
			CreatedAtLessThan(oneYearAgo)
	}).Exec(ctx)

// Force delete (bypass soft delete)
_, err := db.NewDelete().Model(user).WherePK().ForceDelete().Exec(ctx)

// Delete with RETURNING
err := db.NewDelete().Model(user).WherePK().ReturningAll().Scan(ctx)
```

## Raw Queries

```go
// Raw SQL with parameter binding
var result []MyStruct
db.NewRaw("SELECT * FROM users WHERE status = ?", "active").Scan(ctx, &result)
```

## Soft Delete Support

```go
// Query only soft-deleted records
db.NewSelect().Model(&users).WhereDeleted().Scan(ctx)

// Include soft-deleted records
db.NewSelect().Model(&users).IncludeDeleted().Scan(ctx)
```

## Next Step

- [Transactions](./transactions) — running mutations atomically with `RunInTx`
- [ORM: DDL & Surface Map](./orm-ddl) — schema DDL builders and the complete public `orm` surface
