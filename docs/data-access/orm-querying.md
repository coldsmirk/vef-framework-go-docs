---
sidebar_position: 5
---

# ORM: Querying

VEF wraps Bun into a typed, fluent query builder API that provides type-safe SQL construction, automatic audit field handling, and cross-database dialect support. All query builders are accessed through `orm.DB`.

This page covers reading data: SELECT clauses, conditions, joins, ordering, pagination, grouping, locking, execution, and query composition. Sibling pages cover [expressions and aggregates](./orm-expressions), [mutations](./orm-mutations), and [DDL and the public surface map](./orm-ddl).

## API Surface Policy

The ORM package exposes two audited surfaces. VEF-owned ORM method families are
documented by receiver/category on this page and in the model/transaction
guides; every exact method signature is listed in the public API index. These
are the `orm.SelectQuery`, `orm.InsertQuery`, `orm.UpdateQuery`,
`orm.DeleteQuery`, `orm.MergeQuery`, DDL, condition, expression, aggregate, and
window-builder contracts used by application code.

The Bun pass-through surface consists of aliases such as `orm.BunSelectQuery`,
`orm.BunInsertQuery`, `orm.BunUpdateQuery`, and `orm.BunDeleteQuery`; those
methods follow upstream [github.com/uptrace/bun](https://github.com/uptrace/bun)
behavior at the pinned source dependency version `v1.2.18`. Bun/schema aliases
such as `Table`, `Field`, `Relation`, and `Dialect` follow the same
pass-through policy. Do not read VEF query-interface behavior from the Bun
aliases: for example, VEF `orm.SelectQuery.Count` returns `int64`, while the
upstream Bun alias method signatures are tracked separately in the public API
index. Query interfaces that embed `fmt.Stringer` expose `String()` for
SQL/debug rendering; exact signatures are tracked in the public API index.

## Overview

| Category | Builder | Description |
| --- | --- | --- |
| SELECT | `db.NewSelect()` | Query data |
| INSERT | `db.NewInsert()` | Create records |
| UPDATE | `db.NewUpdate()` | Modify records |
| DELETE | `db.NewDelete()` | Delete records |
| MERGE | `db.NewMerge()` | Upsert operations |
| Raw SQL | `db.NewRaw(query, args...)` | Execute raw SQL |

## Getting Started

All query operations start from `orm.DB`, which you receive via dependency injection:

```go
func NewUserService(db orm.DB) *UserService {
	return &UserService{db: db}
}
```

## SELECT Clause

### Basic Select

```go
// SELECT all columns from users
var users []User
err := db.NewSelect().
	Model(&users).
	Scan(ctx)
```

### Selecting Specific Columns

```go
// SELECT su.id, su.username FROM sys_user AS su
var users []User
err := db.NewSelect().
	Model(&users).
	Select("id", "username").
	Scan(ctx)
```

### Column Alias

```go
// SELECT su.username AS name FROM sys_user AS su
var users []User
err := db.NewSelect().
	Model(&users).
	SelectAs("username", "name").
	Scan(ctx)
```

### Excluding Columns

```go
// Select all columns except password
var users []User
err := db.NewSelect().
	Model(&users).
	Exclude("password").
	Scan(ctx)
```

### Expression Columns

```go
// SELECT ..., UPPER(su.username) AS upper_name
var results []struct {
	User
	UpperName string `bun:"upper_name"`
}
err := db.NewSelect().
	Model(&results).
	SelectExpr(func(eb orm.ExprBuilder) any {
		return eb.Upper(eb.Column("username"))
	}, "upper_name").
	Scan(ctx)
```

### Select Model Columns / PKs

```go
// Select only the model's declared columns (no extra expressions)
db.NewSelect().Model(&users).SelectModelColumns()

// Select only primary key columns
db.NewSelect().Model(&users).SelectModelPKs()
```

### Distinct

```go
// SELECT DISTINCT su.department_id FROM sys_user AS su
db.NewSelect().
	Model(&users).
	Select("department_id").
	Distinct().
	Scan(ctx)

// PostgreSQL DISTINCT ON
db.NewSelect().
	Model(&users).
	DistinctOnColumns("department_id").
	Scan(ctx)
```

## WHERE Conditions

The `Where` method takes a callback with a `ConditionBuilder` that provides type-safe condition construction.

### Equality

```go
// WHERE su.is_active = TRUE
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.IsTrue("is_active")
	}).Scan(ctx)

// WHERE su.username = 'admin'
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.Equals("username", "admin")
	}).Scan(ctx)
```

### Comparison Operators

```go
// WHERE su.age > 18 AND su.age <= 65
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.GreaterThan("age", 18).
			LessThanOrEqual("age", 65)
	}).Scan(ctx)
```

### BETWEEN

```go
// WHERE su.created_at BETWEEN '2024-01-01' AND '2024-12-31'
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.Between("created_at", startDate, endDate)
	}).Scan(ctx)
```

### IN / NOT IN

```go
// WHERE su.status IN ('active', 'pending')
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.In("status", []string{"active", "pending"})
	}).Scan(ctx)

// WHERE su.id NOT IN (SELECT ... )
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.NotInSubQuery("id", func(sq orm.SelectQuery) {
			sq.Model((*BlockedUser)(nil)).Select("user_id")
		})
	}).Scan(ctx)
```

### NULL Checks

```go
// WHERE su.deleted_at IS NULL
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.IsNull("deleted_at")
	}).Scan(ctx)

// WHERE su.email IS NOT NULL
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.IsNotNull("email")
	}).Scan(ctx)
```

### String Matching (LIKE)

```go
// WHERE su.username LIKE '%admin%'
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.Contains("username", "admin")
	}).Scan(ctx)

// WHERE su.email LIKE 'test%'
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.StartsWith("email", "test")
	}).Scan(ctx)

// WHERE su.name LIKE '%son'
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.EndsWith("name", "son")
	}).Scan(ctx)

// Case-insensitive: WHERE LOWER(su.email) LIKE LOWER('%Test%')
// (PostgreSQL uses ILIKE automatically)
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.ContainsIgnoreCase("email", "Test")
	}).Scan(ctx)

// Match any of multiple values
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.ContainsAny("name", []string{"John", "Jane"})
	}).Scan(ctx)
```

### OR Conditions

Every condition method has an `Or` prefix variant:

```go
// WHERE su.status = 'active' OR su.status = 'pending'
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.Equals("status", "active").
			OrEquals("status", "pending")
	}).Scan(ctx)
```

### Grouped Conditions (Parentheses)

```go
// WHERE (su.role = 'admin' OR su.role = 'super_admin') AND su.is_active = TRUE
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.Group(func(inner orm.ConditionBuilder) {
			inner.Equals("role", "admin").
				OrEquals("role", "super_admin")
		}).
		IsTrue("is_active")
	}).Scan(ctx)
```

### Column-to-Column Comparison

```go
// WHERE su.created_at <> su.updated_at
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.NotEqualsColumn("created_at", "updated_at")
	}).Scan(ctx)
```

### Subquery Conditions

```go
// WHERE su.department_id = (SELECT id FROM departments WHERE code = 'IT')
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.EqualsSubQuery("department_id", func(sq orm.SelectQuery) {
			sq.Model((*Department)(nil)).
				Select("id").
				Where(func(inner orm.ConditionBuilder) {
					inner.Equals("code", "IT")
				})
		})
	}).Scan(ctx)

// WHERE su.salary > ALL (SELECT salary FROM ...)
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.GreaterThanAll("salary", func(sq orm.SelectQuery) {
			sq.Model((*Employee)(nil)).Select("salary").
				Where(func(inner orm.ConditionBuilder) {
					inner.Equals("department_id", deptID)
				})
		})
	}).Scan(ctx)
```

### Expression Conditions

```go
// WHERE EXTRACT(YEAR FROM su.created_at) = 2024
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.Expr(func(eb orm.ExprBuilder) any {
			return eb.Equals(
				eb.ExtractYear(eb.Column("created_at")),
				2024,
			)
		})
	}).Scan(ctx)
```

### Audit Condition Shortcuts

```go
// WHERE su.created_by = ? (current user from context)
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.CreatedByEqualsCurrent()
	}).Scan(ctx)

// WHERE su.created_at BETWEEN ? AND ?
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.CreatedAtBetween(startTime, endTime)
	}).Scan(ctx)

// WHERE su.updated_by IN ('user1', 'user2')
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.UpdatedByIn([]string{"user1", "user2"})
	}).Scan(ctx)
```

### Primary Key Shortcuts

```go
// WHERE su.id = ?
db.NewSelect().Model(&user).
	Where(func(cb orm.ConditionBuilder) {
		cb.PKEquals(userID)
	}).Scan(ctx)

// WHERE su.id IN (?, ?, ?)
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.PKIn([]string{id1, id2, id3})
	}).Scan(ctx)
```

## JOIN Operations

VEF supports multiple join strategies, each with 4 source variants: Model, Table name, SubQuery, and Expression.

### Join by Model

```go
// INNER JOIN sys_department AS sd ON su.department_id = sd.id
db.NewSelect().Model(&users).
	Join((*Department)(nil), func(cb orm.ConditionBuilder) {
		cb.EqualsColumn("su.department_id", "sd.id")
	}).Scan(ctx)
```

### Left Join

```go
// LEFT JOIN sys_department AS sd ON su.department_id = sd.id
db.NewSelect().Model(&users).
	LeftJoin((*Department)(nil), func(cb orm.ConditionBuilder) {
		cb.EqualsColumn("su.department_id", "sd.id")
	}).Scan(ctx)
```

### Join with Custom Alias

```go
// LEFT JOIN sys_department AS dept ON su.department_id = dept.id
db.NewSelect().Model(&users).
	LeftJoin((*Department)(nil), func(cb orm.ConditionBuilder) {
		cb.EqualsColumn("su.department_id", "dept.id")
	}, "dept").Scan(ctx)
```

### Join by Table Name

```go
// LEFT JOIN departments AS d ON su.department_id = d.id
db.NewSelect().Model(&users).
	LeftJoinTable("departments", func(cb orm.ConditionBuilder) {
		cb.EqualsColumn("su.department_id", "d.id")
	}, "d").Scan(ctx)
```

### Join with SubQuery

```go
// INNER JOIN (SELECT department_id, COUNT(*) AS cnt FROM ...) AS dept_stats
// ON su.department_id = dept_stats.department_id
db.NewSelect().Model(&users).
	JoinSubQuery(
		func(sq orm.SelectQuery) {
			sq.Model((*User)(nil)).
				Select("department_id").
				SelectExpr(func(eb orm.ExprBuilder) any {
					return eb.CountAll()
				}, "cnt").
				GroupBy("department_id")
		},
		func(cb orm.ConditionBuilder) {
			cb.EqualsColumn("su.department_id", "dept_stats.department_id")
		},
		"dept_stats",
	).Scan(ctx)
```

### JoinRelations (Declarative)

For common foreign-key JOINs, `JoinRelations` provides a declarative approach:

```go
// Automatically resolves: LEFT JOIN sys_department AS sd ON su.department_id = sd.id
db.NewSelect().Model(&users).
	JoinRelations(&orm.RelationSpec{
		Model: (*Department)(nil),
		SelectedColumns: []orm.ColumnInfo{
			{Name: "name", Alias: "department_name"},
		},
	}).Scan(ctx)
```

`RelationSpec` fields:
- `Model`: the related model (required)
- `Alias`: custom table alias (default: model's default alias)
- `JoinType`: `orm.JoinLeft` (default), `orm.JoinInner`, `orm.JoinRight`, etc.
- `ForeignColumn`: auto-resolved to `{model_name}_{pk}` if empty
- `ReferencedColumn`: auto-resolved to PK if empty
- `SelectedColumns`: which columns to select with aliasing
- `On`: additional JOIN conditions

### Bun Relations

```go
// Load Bun-defined relations
db.NewSelect().Model(&users).
	Relation("Department").
	Relation("Roles", func(sq orm.SelectQuery) {
		sq.Where(func(cb orm.ConditionBuilder) {
			cb.IsTrue("is_active")
		})
	}).Scan(ctx)
```

## Ordering

```go
// ORDER BY su.created_at ASC
db.NewSelect().Model(&users).
	OrderBy("created_at").Scan(ctx)

// ORDER BY su.created_at DESC
db.NewSelect().Model(&users).
	OrderByDesc("created_at").Scan(ctx)

// ORDER BY CASE ... END (expression-based)
db.NewSelect().Model(&users).
	OrderByExpr(func(eb orm.ExprBuilder) any {
		return eb.Case(func(cb orm.CaseBuilder) {
			cb.When(func(c orm.ConditionBuilder) { c.Equals("status", "active") }).Then(1).
				When(func(c orm.ConditionBuilder) { c.Equals("status", "pending") }).Then(2).
				Else(3)
		})
	}).Scan(ctx)
```

## Pagination

```go
// Simple LIMIT/OFFSET
db.NewSelect().Model(&users).
	Limit(20).Offset(40).Scan(ctx)

// Using page.Pageable (framework convention)
p := page.Pageable{Page: 3, Size: 20}
db.NewSelect().Model(&users).
	Paginate(p).Scan(ctx)

// ScanAndCount: fetch page data + total count in one call
total, err := db.NewSelect().Model(&users).
	Paginate(p).
	ScanAndCount(ctx)
```

## GROUP BY & HAVING

```go
// SELECT department_id, COUNT(*) AS cnt
// FROM sys_user GROUP BY department_id HAVING COUNT(*) > 5
type DeptCount struct {
	DepartmentID string `bun:"department_id"`
	Cnt          int64  `bun:"cnt"`
}
var results []DeptCount
err := db.NewSelect().
	Model((*User)(nil)).
	Select("department_id").
	SelectExpr(func(eb orm.ExprBuilder) any {
		return eb.CountAll()
	}, "cnt").
	GroupBy("department_id").
	Having(func(cb orm.ConditionBuilder) {
		cb.Expr(func(eb orm.ExprBuilder) any {
			return eb.GreaterThan(eb.CountAll(), 5)
		})
	}).Scan(ctx, &results)
```

## Row-Level Locking

```go
// SELECT ... FOR UPDATE
db.NewSelect().Model(&user).
	Where(func(cb orm.ConditionBuilder) { cb.PKEquals(id) }).
	ForUpdate().
	Scan(ctx)

// FOR UPDATE NOWAIT
db.NewSelect().Model(&user).ForUpdateNoWait().Scan(ctx)

// FOR UPDATE SKIP LOCKED (useful for job queues)
db.NewSelect().Model(&tasks).
	Where(func(cb orm.ConditionBuilder) {
		cb.Equals("status", "pending")
	}).
	Limit(10).
	ForUpdateSkipLocked().
	Scan(ctx)

// FOR SHARE
db.NewSelect().Model(&user).ForShare().Scan(ctx)

// FOR KEY SHARE / FOR NO KEY UPDATE (PostgreSQL only)
db.NewSelect().Model(&user).ForKeyShare().Scan(ctx)
db.NewSelect().Model(&user).ForNoKeyUpdate().Scan(ctx)
```

> Note: SQLite does not support row-level locking — calls are silently ignored with a warning.

## Execution Methods

These rows describe VEF `orm.SelectQuery` execution methods, not the upstream
`orm.BunSelectQuery` pass-through alias.

| Method | Returns | Purpose |
| --- | --- | --- |
| `Scan(ctx, dest...)` | `error` | Scan rows into model or dest |
| `Exec(ctx, dest...)` | `sql.Result, error` | Execute without scan |
| `Rows(ctx)` | `*sql.Rows, error` | Get raw rows iterator |
| `ScanAndCount(ctx)` | `int64, error` | Scan + count total (for pagination) |
| `Count(ctx)` | `int64, error` | Count only |
| `Exists(ctx)` | `bool, error` | Check existence |

## Query Composition with Apply

The `Apply` / `ApplyIf` pattern enables reusable query fragments:

```go
// Define reusable conditions
func ActiveOnly(q orm.SelectQuery) {
	q.Where(func(cb orm.ConditionBuilder) {
		cb.IsTrue("is_active")
	})
}

func CreatedAfter(t time.Time) orm.ApplyFunc[orm.SelectQuery] {
	return func(q orm.SelectQuery) {
		q.Where(func(cb orm.ConditionBuilder) {
			cb.CreatedAtGreaterThanOrEqual(t)
		})
	}
}

// Use them
db.NewSelect().Model(&users).
	Apply(ActiveOnly, CreatedAfter(lastMonth)).
	Scan(ctx)

// Conditional application
db.NewSelect().Model(&users).
	ApplyIf(keyword != "", func(q orm.SelectQuery) {
		q.Where(func(cb orm.ConditionBuilder) {
			cb.Contains("username", keyword)
		})
	}).
	Scan(ctx)
```

## Next Step

- [ORM: Expressions & Aggregates](./orm-expressions) — the expression builder, aggregate and window functions, CTEs, and set operations
- [ORM: Mutations](./orm-mutations) — INSERT, UPDATE, DELETE, raw queries, and soft-delete behavior
- [Generic CRUD](./crud) — higher-level CRUD operations built on top of the query builder
