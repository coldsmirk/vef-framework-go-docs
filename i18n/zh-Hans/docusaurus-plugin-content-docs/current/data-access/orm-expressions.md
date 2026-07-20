---
sidebar_position: 6
---

# ORM：表达式与聚合

计算型 SQL 的构建单元：聚合与窗口函数、公共表表达式、集合操作以及表达式构造器。本页内容与 [ORM：查询](./orm-querying) 中的查询接口组合使用。

## 聚合函数

`ExprBuilder` 提供了所有标准 SQL 聚合函数：

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

高级聚合：

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

`JSONObjectAggBuilder` 的键来自 `KeyColumn(column)` 或 `KeyExpr(expr)`，
值来自继承的 `Column(...)` / `Expr(...)` 方法。

## 窗口函数

窗口子句始终从 `Over()` 开始：分区与排序方法都挂在 `Over()` 返回的分区
构建器上。

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

// 窗口聚合函数：SUM() OVER (...)
eb.WinSum(func(ws orm.WindowSumBuilder) {
	ws.Column("amount")
	ws.Over().PartitionBy("department_id").OrderBy("created_at")
})
```

## 公共表表达式 (CTE)

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

### 递归 CTE

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

`SelectQuery`、`InsertQuery`、`UpdateQuery`、`DeleteQuery` 和 `MergeQuery`
都公开 `WithValues(name, model)` 用于 VALUES-based CTE，也公开
`WithOrderedValues(name, model)` 用于会追加 ordinal column 的 VALUES CTE，
从而保留 slice 顺序。当前 `WithValues` 签名不再接收排序 flag；需要顺序时请调用
`WithOrderedValues`。

## 集合操作

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

## 表达式构造器 (ExprBuilder)

`ExprBuilder` 是 VEF 类型安全 SQL 表达式构建的核心。它在 `SelectExpr`、`Where` → `Expr`、`SetExpr`、`ColumnExpr` 等方法中可用。

### 字符串函数

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
eb.Contains(eb.Column("name"), "admin")        // LIKE '%admin%'
eb.StartsWith(eb.Column("name"), "prefix")      // LIKE 'prefix%'
eb.EndsWith(eb.Column("name"), "suffix")         // LIKE '%suffix'
eb.ContainsIgnoreCase(eb.Column("name"), "ADM")  // 不区分大小写
```

### 日期和时间函数

```go
eb.CurrentDate()                    // CURRENT_DATE
eb.CurrentTime()                    // CURRENT_TIME
eb.CurrentTimestamp()               // CURRENT_TIMESTAMP
eb.Now()                            // NOW()（方言自适应）
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

### 数学函数

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

### 类型转换（方言自适应）

```go
eb.ToString(eb.Column("id"))         // CAST(id AS TEXT) / ::TEXT
eb.ToInteger(eb.Column("str_col"))    // CAST(... AS INTEGER) / ::INTEGER
eb.ToDecimal(eb.Column("price"))      // CAST(... AS DECIMAL) / ::NUMERIC
eb.ToDecimal(eb.Column("price"), 10, 2) // 带精度
eb.ToFloat(eb.Column("rate"))         // CAST(... AS DOUBLE) / ::DOUBLE PRECISION
```

### 条件函数

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

`CaseBuilder.When` 接受条件构建回调（不存在 `When(column, value)` 简写）；
表达式与子查询形态的分支用 `WhenExpr(expr)` 和 `WhenSubQuery(op, build)`。

### JSON 函数（方言自适应）

```go
eb.JSONExtract(eb.Column("data"), "$.name")            // 路径处的原始 JSON 值
eb.JSONUnquote(eb.JSONExtract(eb.Column("data"), "$.name")) // 去引号文本
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

不存在 `JSONRemove` helper；删除路径请用 `JSONSet` 置 `null` 或方言特定的
原始表达式实现。

### 跨数据库方言支持

优先使用内置的方言自适应 helper，例如 `ToString`、`ToDecimal`、
`JSONExtract`、`JSONObject` 等表达式方法。`ExprByDialect` 因为出现在
公开 method set 上，所以会被生成索引列出；但它的配置类型目前没有从公开 `orm`
包 re-export，应用代码不应该直接构造 dialect map。

## 下一步

- [ORM：查询](./orm-querying) — SELECT 子句、条件、连接与执行
- [ORM：写入操作](./orm-mutations) — INSERT、UPDATE、DELETE、原始 SQL 与软删除行为
