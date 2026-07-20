---
sidebar_position: 6
---

# ORM: Expressions & Aggregates

Building blocks for computed SQL: aggregate and window functions, common table expressions, set operations, and the expression builder. Everything here composes with the query surface documented in [ORM: Querying](./orm-querying).

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
		ja.KeyColumn("code").Column("name")
	})
}, "code_map")
```

`JSONObjectAggBuilder` takes its keys from `KeyColumn(column)` or
`KeyExpr(expr)` and its values from the inherited `Column(...)` / `Expr(...)`
methods.

## Window Functions

The window clause always starts from `Over()`: partitioning and ordering
methods live on the partition builder that `Over()` returns.

```go
// ROW_NUMBER() OVER (PARTITION BY department_id ORDER BY salary DESC)
db.NewSelect().Model(&users).
	SelectExpr(func(eb orm.ExprBuilder) any {
		return eb.RowNumber(func(rn orm.RowNumberBuilder) {
			rn.Over().
				PartitionBy("department_id").
				OrderByDesc("salary")
		})
	}, "row_num").
	Scan(ctx)

// RANK(), DENSE_RANK(), PERCENT_RANK()
eb.Rank(func(r orm.RankBuilder) { r.Over().OrderByDesc("score") })
eb.DenseRank(func(r orm.DenseRankBuilder) { r.Over().OrderByDesc("score") })

// LAG / LEAD
eb.Lag(func(l orm.LagBuilder) {
	l.Column("salary").Offset(1).DefaultValue(0)
	l.Over().OrderBy("created_at")
})
eb.Lead(func(l orm.LeadBuilder) {
	l.Column("salary").Offset(1)
	l.Over().OrderBy("created_at")
})

// FIRST_VALUE / LAST_VALUE / NTH_VALUE
eb.FirstValue(func(fv orm.FirstValueBuilder) {
	fv.Column("name")
	fv.Over().OrderBy("created_at")
})

// Windowed aggregates: SUM() OVER (...)
eb.WinSum(func(ws orm.WindowSumBuilder) {
	ws.Column("amount")
	ws.Over().PartitionBy("department_id").OrderBy("created_at")
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

`SelectQuery`, `InsertQuery`, `UpdateQuery`, `DeleteQuery`, and `MergeQuery`
all expose `WithValues(name, model)` for VALUES-based CTEs and
`WithOrderedValues(name, model)` for VALUES CTEs that append an ordinal column
so slice order can be preserved. The current `WithValues` signature does not
accept an ordering flag; call `WithOrderedValues` when order matters.

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
	cb.When(func(c orm.ConditionBuilder) { c.Equals("status", "active") }).Then("Active").
		When(func(c orm.ConditionBuilder) { c.Equals("status", "inactive") }).Then("Inactive").
		Else("Unknown")
})
```

`CaseBuilder.When` takes a condition-builder callback (there is no
`When(column, value)` shorthand); `WhenExpr(expr)` and
`WhenSubQuery(op, build)` cover expression- and subquery-shaped branches.

### JSON Functions (Dialect-Aware)

```go
eb.JSONExtract(eb.Column("data"), "$.name")            // raw JSON value at path
eb.JSONUnquote(eb.JSONExtract(eb.Column("data"), "$.name")) // unquoted text
eb.JSONObject("key1", value1, "key2", value2)
eb.JSONArray(value1, value2, value3)
eb.JSONContains(eb.Column("tags"), `"admin"`)
eb.JSONContainsPath(eb.Column("data"), "$.address")
eb.JSONLength(eb.Column("items"))
eb.JSONKeys(eb.Column("data"))
eb.JSONType(eb.Column("data"), "$.age")
eb.JSONValid(eb.Column("payload"))
eb.JSONSet(eb.Column("data"), "$.status", "active")
eb.JSONInsert(eb.Column("data"), "$.new", 1)
eb.JSONReplace(eb.Column("data"), "$.old", 2)
eb.JSONArrayAppend(eb.Column("tags"), "$", "extra")
```

There is no `JSONRemove` helper; model removals as `JSONSet` to `null` or a
dialect-specific raw expression.

### Cross-Database Dialect Support

Prefer the built-in dialect-aware helpers (`ToString`, `ToDecimal`,
`JSONExtract`, `JSONObject`, and similar expression methods). The
low-level `ExprByDialect` hook is part of the generated index because it appears
on the public method set, but its configuration type is not re-exported from the
public `orm` package; application code should not construct dialect maps
directly.

## Next Step

- [ORM: Querying](./orm-querying) — SELECT clauses, conditions, joins, and execution
- [ORM: Mutations](./orm-mutations) — INSERT, UPDATE, DELETE, raw queries, and soft-delete behavior
