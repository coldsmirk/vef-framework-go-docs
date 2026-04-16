---
sidebar_position: 10
---

# ORM Builder

VEF wraps Bun into a typed, fluent query builder API that provides type-safe SQL construction, automatic audit field handling, and cross-database dialect support.

This page is a comprehensive reference for building SQL queries in VEF projects. All query builders are accessed through `orm.DB`.

## Overview

| Category | Builder | Description |
| --- | --- | --- |
| SELECT | `db.NewSelect()` | Query data |
| INSERT | `db.NewInsert()` | Create records |
| UPDATE | `db.NewUpdate()` | Modify records |
| DELETE | `db.NewDelete()` | Delete records |
| MERGE | `db.NewMerge()` | Upsert operations |
| Raw SQL | `db.NewRawQuery()` | Execute raw SQL |

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
			cb.When("status", "active").Then(1).
				When("status", "pending").Then(2).
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

## Aggregate Functions

The `ExprBuilder` provides all standard SQL aggregate functions:

```go
db.NewSelect().Model((*User)(nil)).
	SelectExpr(func(eb orm.ExprBuilder) any { return eb.CountAll() }, "total").
	SelectExpr(func(eb orm.ExprBuilder) any { return eb.CountColumn("id", true) }, "distinct_count"). // COUNT(DISTINCT id)
	SelectExpr(func(eb orm.ExprBuilder) any { return eb.SumColumn("salary") }, "total_salary").
	SelectExpr(func(eb orm.ExprBuilder) any { return eb.AvgColumn("salary") }, "avg_salary").
	SelectExpr(func(eb orm.ExprBuilder) any { return eb.MinColumn("salary") }, "min_salary").
	SelectExpr(func(eb orm.ExprBuilder) any { return eb.MaxColumn("salary") }, "max_salary").
	Scan(ctx, &result)
```

Advanced aggregates:

```go
// STRING_AGG / GROUP_CONCAT
SelectExpr(func(eb orm.ExprBuilder) any {
	return eb.StringAgg(func(sa orm.StringAggBuilder) {
		sa.Column("name").Separator(",")
	})
}, "names")

// ARRAY_AGG (PostgreSQL)
SelectExpr(func(eb orm.ExprBuilder) any {
	return eb.ArrayAgg(func(aa orm.ArrayAggBuilder) {
		aa.Column("tag").Distinct()
	})
}, "tags")

// JSON_OBJECT_AGG
SelectExpr(func(eb orm.ExprBuilder) any {
	return eb.JSONObjectAgg(func(ja orm.JSONObjectAggBuilder) {
		ja.Key("code").Value("name")
	})
}, "code_map")
```

## Window Functions

```go
// ROW_NUMBER() OVER (PARTITION BY department_id ORDER BY salary DESC)
db.NewSelect().Model(&users).
	SelectExpr(func(eb orm.ExprBuilder) any {
		return eb.RowNumber(func(rn orm.RowNumberBuilder) {
			rn.PartitionBy("department_id").
				OrderByDesc("salary")
		})
	}, "row_num").
	Scan(ctx)

// RANK(), DENSE_RANK(), PERCENT_RANK()
eb.Rank(func(r orm.RankBuilder) { r.OrderByDesc("score") })
eb.DenseRank(func(r orm.DenseRankBuilder) { r.OrderByDesc("score") })

// LAG / LEAD
eb.Lag(func(l orm.LagBuilder) {
	l.Column("salary").Offset(1).Default(0)
})
eb.Lead(func(l orm.LeadBuilder) {
	l.Column("salary").Offset(1)
})

// FIRST_VALUE / LAST_VALUE / NTH_VALUE
eb.FirstValue(func(fv orm.FirstValueBuilder) {
	fv.Column("name").OrderBy("created_at")
})

// Windowed aggregates: SUM() OVER (...)
eb.WinSum(func(ws orm.WindowSumBuilder) {
	ws.Column("amount").PartitionBy("department_id").OrderBy("created_at")
})
```

## Common Table Expressions (CTE)

```go
// WITH active_users AS (SELECT * FROM sys_user WHERE is_active = TRUE)
// SELECT * FROM active_users WHERE ...
db.NewSelect().Model(&users).
	With("active_users", func(sq orm.SelectQuery) {
		sq.Model((*User)(nil)).
			Where(func(cb orm.ConditionBuilder) {
				cb.IsTrue("is_active")
			})
	}).
	Table("active_users").
	Scan(ctx)
```

### Recursive CTE

```go
// WITH RECURSIVE org_tree AS (
//   SELECT * FROM departments WHERE parent_id IS NULL
//   UNION ALL
//   SELECT d.* FROM departments d JOIN org_tree ot ON d.parent_id = ot.id
// )
db.NewSelect().Model(&departments).
	WithRecursive("org_tree", func(sq orm.SelectQuery) {
		sq.Model((*Department)(nil)).
			Where(func(cb orm.ConditionBuilder) {
				cb.IsNull("parent_id")
			}).
			UnionAll(func(uq orm.SelectQuery) {
				uq.Model((*Department)(nil)).
					JoinTable("org_tree", func(cb orm.ConditionBuilder) {
						cb.EqualsColumn("sd.parent_id", "org_tree.id")
					}, "org_tree")
			})
	}).
	Table("org_tree").
	Scan(ctx)
```

## Set Operations

```go
// UNION / UNION ALL
db.NewSelect().Model(&activeUsers).
	Union(func(sq orm.SelectQuery) {
		sq.Model((*ArchivedUser)(nil))
	}).Scan(ctx)

db.NewSelect().Model(&set1).
	UnionAll(func(sq orm.SelectQuery) { sq.Model((*Set2)(nil)) }).Scan(ctx)

// INTERSECT / EXCEPT
db.NewSelect().Model(&set1).
	Intersect(func(sq orm.SelectQuery) { sq.Model((*Set2)(nil)) }).Scan(ctx)

db.NewSelect().Model(&set1).
	Except(func(sq orm.SelectQuery) { sq.Model((*Set2)(nil)) }).Scan(ctx)
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

| Method | Returns | Purpose |
| --- | --- | --- |
| `Scan(ctx, dest...)` | `error` | Scan rows into model or dest |
| `Exec(ctx, dest...)` | `sql.Result, error` | Execute without scan |
| `Rows(ctx)` | `*sql.Rows, error` | Get raw rows iterator |
| `ScanAndCount(ctx)` | `int64, error` | Scan + count total (for pagination) |
| `Count(ctx)` | `int64, error` | Count only |
| `Exists(ctx)` | `bool, error` | Check existence |

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

## Expression Builder (ExprBuilder)

The `ExprBuilder` is the core of VEF's type-safe SQL expression building. It is available in `SelectExpr`, `Where` → `Expr`, `SetExpr`, `ColumnExpr`, etc.

### String Functions

```go
eb.Concat(eb.Column("first_name"), " ", eb.Column("last_name"))
eb.ConcatWithSep(", ", eb.Column("city"), eb.Column("state"))
eb.Upper(eb.Column("name"))
eb.Lower(eb.Column("email"))
eb.Trim(eb.Column("name"))
eb.TrimLeft(eb.Column("name"))
eb.TrimRight(eb.Column("name"))
eb.SubString(eb.Column("name"), 1, 3)
eb.Length(eb.Column("name"))
eb.CharLength(eb.Column("name"))
eb.Position("@", eb.Column("email"))
eb.Left(eb.Column("name"), 5)
eb.Right(eb.Column("name"), 3)
eb.Replace(eb.Column("name"), "old", "new")
eb.Repeat(eb.Column("char"), 3)
eb.Reverse(eb.Column("name"))
eb.Contains(eb.Column("name"), "admin")       // LIKE '%admin%'
eb.StartsWith(eb.Column("name"), "prefix")     // LIKE 'prefix%'
eb.EndsWith(eb.Column("name"), "suffix")       // LIKE '%suffix'
eb.ContainsIgnoreCase(eb.Column("name"), "ADM") // case-insensitive
```

### Date & Time Functions

```go
eb.CurrentDate()                    // CURRENT_DATE
eb.CurrentTime()                    // CURRENT_TIME
eb.CurrentTimestamp()               // CURRENT_TIMESTAMP
eb.Now()                            // NOW() (dialect-aware)
eb.ExtractYear(eb.Column("created_at"))
eb.ExtractMonth(eb.Column("created_at"))
eb.ExtractDay(eb.Column("created_at"))
eb.ExtractHour(eb.Column("created_at"))
eb.ExtractMinute(eb.Column("created_at"))
eb.ExtractSecond(eb.Column("created_at"))
eb.DateTrunc(orm.UnitMonth, eb.Column("created_at"))
eb.DateAdd(eb.Column("created_at"), 7, orm.UnitDay)
eb.DateSubtract(eb.Column("created_at"), 1, orm.UnitMonth)
eb.DateDiff(eb.Column("start_date"), eb.Column("end_date"), orm.UnitDay)
eb.Age(eb.Column("birth_date"), eb.Now())
```

### Math Functions

```go
eb.Abs(eb.Column("balance"))
eb.Ceil(eb.Column("price"))
eb.Floor(eb.Column("price"))
eb.Round(eb.Column("price"), 2)
eb.Trunc(eb.Column("price"), 2)
eb.Power(eb.Column("base"), 2)
eb.Sqrt(eb.Column("area"))
eb.Mod(eb.Column("id"), 10)
eb.Greatest(eb.Column("a"), eb.Column("b"), eb.Column("c"))
eb.Least(eb.Column("a"), eb.Column("b"))
eb.Random()
eb.Sign(eb.Column("amount"))
```

### Type Casting (Dialect-Aware)

```go
eb.ToString(eb.Column("id"))        // CAST(id AS TEXT) / ::TEXT
eb.ToInteger(eb.Column("str_col"))   // CAST(... AS INTEGER) / ::INTEGER
eb.ToDecimal(eb.Column("price"))     // CAST(... AS DECIMAL) / ::NUMERIC
eb.ToDecimal(eb.Column("price"), 10, 2) // with precision
eb.ToFloat(eb.Column("rate"))        // CAST(... AS DOUBLE) / ::DOUBLE PRECISION
```

### Conditional Functions

```go
eb.Coalesce(eb.Column("nickname"), eb.Column("username"), "Anonymous")
eb.NullIf(eb.Column("value"), 0)
eb.IfNull(eb.Column("name"), "Unknown")
eb.Case(func(cb orm.CaseBuilder) {
	cb.When("status", "active").Then("Active").
		When("status", "inactive").Then("Inactive").
		Else("Unknown")
})
```

### JSON Functions (Dialect-Aware)

```go
eb.JSONExtractText(eb.Column("data"), "$.name")
eb.JSONExtractInt(eb.Column("data"), "$.age")
eb.JSONExtractBool(eb.Column("data"), "$.active")
eb.JSONBuildObject("key1", value1, "key2", value2)
eb.JSONBuildArray(value1, value2, value3)
eb.JSONContains(eb.Column("tags"), "admin")
eb.JSONLength(eb.Column("items"))
eb.JSONKeys(eb.Column("data"))
eb.JSONSet(eb.Column("data"), "$.status", "active")
eb.JSONRemove(eb.Column("data"), "$.temp")
```

### Cross-Database Dialect Support

```go
// Execute different SQL per database dialect
eb.ExprByDialect(orm.DialectExprs{
	Postgres: func() any { return eb.Expr("?::JSONB", value) },
	MySQL:    func() any { return eb.Expr("CAST(? AS JSON)", value) },
	SQLite:   func() any { return eb.Expr("JSON(?)", value) },
	Default:  func() any { return eb.Literal(value) },
})
```

## Transactions

```go
// Automatic transaction (recommended)
err := db.RunInTX(ctx, func(ctx context.Context, tx orm.DB) error {
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
err := db.RunInReadOnlyTX(ctx, func(ctx context.Context, tx orm.DB) error {
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

## Raw Queries

```go
// Raw SQL with parameter binding
var result []MyStruct
db.NewRaw("SELECT * FROM users WHERE status = ?", "active").Scan(ctx, &result)
```

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

## DDL Operations

### Create Table

```go
_, err := db.NewCreateTable().
	Model((*User)(nil)).
	IfNotExists().
	Exec(ctx)
```

### Create Index

```go
_, err := db.NewCreateIndex().
	Model((*User)(nil)).
	Index("idx_user_email").
	Column("email").
	Unique().
	IfNotExists().
	Exec(ctx)
```

### Other DDL

```go
db.NewDropTable().Model((*User)(nil)).IfExists().Exec(ctx)
db.NewTruncateTable().Model((*User)(nil)).Exec(ctx)
db.NewAddColumn().Model((*User)(nil)).ColumnExpr("phone VARCHAR(20)").Exec(ctx)
db.NewDropColumn().Model((*User)(nil)).Column("phone").Exec(ctx)
```

## Soft Delete Support

```go
// Query only soft-deleted records
db.NewSelect().Model(&users).WhereDeleted().Scan(ctx)

// Include soft-deleted records
db.NewSelect().Model(&users).IncludeDeleted().Scan(ctx)
```

## Next Step

- [Models](./models) — how to define your data models
- [Generic CRUD](./crud) — higher-level CRUD operations built on top of the query builder
- [Query Builder (Search Tags)](./query-builder) — automatic query building from search tags
