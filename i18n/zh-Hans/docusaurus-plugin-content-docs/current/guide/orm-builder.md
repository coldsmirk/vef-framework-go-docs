---
sidebar_position: 10
---

# ORM 构造器

VEF 将 Bun 封装成类型安全的流式查询构造器 API，提供类型安全的 SQL 构造、自动审计字段处理以及跨数据库方言支持。

本页是在 VEF 项目中构建 SQL 查询的完整参考手册，所有查询构造器都通过 `orm.DB` 访问。

## 概述

| 分类 | 构造器 | 说明 |
| --- | --- | --- |
| 查询 | `db.NewSelect()` | 查询数据 |
| 插入 | `db.NewInsert()` | 创建记录 |
| 更新 | `db.NewUpdate()` | 修改记录 |
| 删除 | `db.NewDelete()` | 删除记录 |
| 合并 | `db.NewMerge()` | Upsert 操作 |
| 原始 SQL | `db.NewRawQuery()` | 执行原始 SQL |

## 快速开始

所有查询操作都从 `orm.DB` 开始，通过依赖注入获取：

```go
func NewUserService(db orm.DB) *UserService {
	return &UserService{db: db}
}
```

## SELECT 子句

### 基本查询

```go
// SELECT 所有列
var users []User
err := db.NewSelect().
	Model(&users).
	Scan(ctx)
```

### 选择特定列

```go
// SELECT su.id, su.username FROM sys_user AS su
var users []User
err := db.NewSelect().
	Model(&users).
	Select("id", "username").
	Scan(ctx)
```

### 列别名

```go
// SELECT su.username AS name FROM sys_user AS su
var users []User
err := db.NewSelect().
	Model(&users).
	SelectAs("username", "name").
	Scan(ctx)
```

### 排除列

```go
// 选择所有列但排除 password
var users []User
err := db.NewSelect().
	Model(&users).
	Exclude("password").
	Scan(ctx)
```

### 表达式列

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

### 选择模型列 / 主键列

```go
// 仅选择模型声明的列（不包含额外表达式）
db.NewSelect().Model(&users).SelectModelColumns()

// 仅选择主键列
db.NewSelect().Model(&users).SelectModelPKs()
```

### DISTINCT

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

## WHERE 条件

`Where` 方法接收一个 `ConditionBuilder` 回调，提供类型安全的条件构造。

### 等于

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

### 比较运算符

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

### NULL 检查

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

### 字符串匹配 (LIKE)

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

// 不区分大小写：WHERE LOWER(su.email) LIKE LOWER('%Test%')
// （PostgreSQL 自动使用 ILIKE）
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.ContainsIgnoreCase("email", "Test")
	}).Scan(ctx)

// 匹配多个值中的任一个
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.ContainsAny("name", []string{"John", "Jane"})
	}).Scan(ctx)
```

### OR 条件

每个条件方法都有 `Or` 前缀的变体：

```go
// WHERE su.status = 'active' OR su.status = 'pending'
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.Equals("status", "active").
			OrEquals("status", "pending")
	}).Scan(ctx)
```

### 分组条件（括号）

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

### 列与列比较

```go
// WHERE su.created_at <> su.updated_at
db.NewSelect().Model(&users).
	Where(func(cb orm.ConditionBuilder) {
		cb.NotEqualsColumn("created_at", "updated_at")
	}).Scan(ctx)
```

### 子查询条件

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

### 表达式条件

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

### 审计条件快捷方法

```go
// WHERE su.created_by = ? （当前上下文用户）
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

### 主键快捷方法

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

## JOIN 操作

VEF 支持多种 JOIN 策略，每种都有 4 种数据源变体：Model、表名、子查询和表达式。

### 通过模型 JOIN

```go
// INNER JOIN sys_department AS sd ON su.department_id = sd.id
db.NewSelect().Model(&users).
	Join((*Department)(nil), func(cb orm.ConditionBuilder) {
		cb.EqualsColumn("su.department_id", "sd.id")
	}).Scan(ctx)
```

### LEFT JOIN

```go
// LEFT JOIN sys_department AS sd ON su.department_id = sd.id
db.NewSelect().Model(&users).
	LeftJoin((*Department)(nil), func(cb orm.ConditionBuilder) {
		cb.EqualsColumn("su.department_id", "sd.id")
	}).Scan(ctx)
```

### 自定义别名的 JOIN

```go
// LEFT JOIN sys_department AS dept ON su.department_id = dept.id
db.NewSelect().Model(&users).
	LeftJoin((*Department)(nil), func(cb orm.ConditionBuilder) {
		cb.EqualsColumn("su.department_id", "dept.id")
	}, "dept").Scan(ctx)
```

### 通过表名 JOIN

```go
// LEFT JOIN departments AS d ON su.department_id = d.id
db.NewSelect().Model(&users).
	LeftJoinTable("departments", func(cb orm.ConditionBuilder) {
		cb.EqualsColumn("su.department_id", "d.id")
	}, "d").Scan(ctx)
```

### 子查询 JOIN

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

### JoinRelations（声明式）

针对常见的外键 JOIN，`JoinRelations` 提供声明式写法：

```go
// 自动解析：LEFT JOIN sys_department AS sd ON su.department_id = sd.id
db.NewSelect().Model(&users).
	JoinRelations(&orm.RelationSpec{
		Model: (*Department)(nil),
		SelectedColumns: []orm.ColumnInfo{
			{Name: "name", Alias: "department_name"},
		},
	}).Scan(ctx)
```

`RelationSpec` 字段说明：
- `Model`：关联模型（必填）
- `Alias`：自定义表别名（默认：模型的默认别名）
- `JoinType`：`orm.JoinLeft`（默认）、`orm.JoinInner`、`orm.JoinRight` 等
- `ForeignColumn`：为空时自动解析为 `{模型名}_{主键}`
- `ReferencedColumn`：为空时自动解析为主键
- `SelectedColumns`：要选择的列及其别名
- `On`：附加 JOIN 条件

### Bun 关联

```go
// 加载 Bun 定义的关联关系
db.NewSelect().Model(&users).
	Relation("Department").
	Relation("Roles", func(sq orm.SelectQuery) {
		sq.Where(func(cb orm.ConditionBuilder) {
			cb.IsTrue("is_active")
		})
	}).Scan(ctx)
```

## 排序

```go
// ORDER BY su.created_at ASC
db.NewSelect().Model(&users).
	OrderBy("created_at").Scan(ctx)

// ORDER BY su.created_at DESC
db.NewSelect().Model(&users).
	OrderByDesc("created_at").Scan(ctx)

// ORDER BY CASE ... END（表达式排序）
db.NewSelect().Model(&users).
	OrderByExpr(func(eb orm.ExprBuilder) any {
		return eb.Case(func(cb orm.CaseBuilder) {
			cb.When("status", "active").Then(1).
				When("status", "pending").Then(2).
				Else(3)
		})
	}).Scan(ctx)
```

## 分页

```go
// 简单的 LIMIT/OFFSET
db.NewSelect().Model(&users).
	Limit(20).Offset(40).Scan(ctx)

// 使用 page.Pageable（框架约定）
p := page.Pageable{Page: 3, Size: 20}
db.NewSelect().Model(&users).
	Paginate(p).Scan(ctx)

// ScanAndCount：一次调用获取分页数据 + 总数
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
		ja.Key("code").Value("name")
	})
}, "code_map")
```

## 窗口函数

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

// 窗口聚合函数：SUM() OVER (...)
eb.WinSum(func(ws orm.WindowSumBuilder) {
	ws.Column("amount").PartitionBy("department_id").OrderBy("created_at")
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

## 行级锁定

```go
// SELECT ... FOR UPDATE
db.NewSelect().Model(&user).
	Where(func(cb orm.ConditionBuilder) { cb.PKEquals(id) }).
	ForUpdate().
	Scan(ctx)

// FOR UPDATE NOWAIT
db.NewSelect().Model(&user).ForUpdateNoWait().Scan(ctx)

// FOR UPDATE SKIP LOCKED（用于任务队列）
db.NewSelect().Model(&tasks).
	Where(func(cb orm.ConditionBuilder) {
		cb.Equals("status", "pending")
	}).
	Limit(10).
	ForUpdateSkipLocked().
	Scan(ctx)

// FOR SHARE
db.NewSelect().Model(&user).ForShare().Scan(ctx)

// FOR KEY SHARE / FOR NO KEY UPDATE（仅 PostgreSQL）
db.NewSelect().Model(&user).ForKeyShare().Scan(ctx)
db.NewSelect().Model(&user).ForNoKeyUpdate().Scan(ctx)
```

> 注意：SQLite 不支持行级锁定——调用会被静默忽略并记录警告。

## 执行方法

| 方法 | 返回值 | 用途 |
| --- | --- | --- |
| `Scan(ctx, dest...)` | `error` | 将行扫描到模型或目标 |
| `Exec(ctx, dest...)` | `sql.Result, error` | 执行但不扫描 |
| `Rows(ctx)` | `*sql.Rows, error` | 获取原始行迭代器 |
| `ScanAndCount(ctx)` | `int64, error` | 扫描 + 计算总数（用于分页）|
| `Count(ctx)` | `int64, error` | 仅计数 |
| `Exists(ctx)` | `bool, error` | 检查是否存在 |

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
	cb.When("status", "active").Then("Active").
		When("status", "inactive").Then("Inactive").
		Else("Unknown")
})
```

### JSON 函数（方言自适应）

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

### 跨数据库方言支持

```go
// 根据数据库方言执行不同的 SQL
eb.ExprByDialect(orm.DialectExprs{
	Postgres: func() any { return eb.Expr("?::JSONB", value) },
	MySQL:    func() any { return eb.Expr("CAST(? AS JSON)", value) },
	SQLite:   func() any { return eb.Expr("JSON(?)", value) },
	Default:  func() any { return eb.Literal(value) },
})
```

## 事务

```go
// 自动事务（推荐）
err := db.RunInTX(ctx, func(ctx context.Context, tx orm.DB) error {
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
err := db.RunInReadOnlyTX(ctx, func(ctx context.Context, tx orm.DB) error {
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

## 原始 SQL 查询

```go
// 带参数绑定的原始 SQL
var result []MyStruct
db.NewRaw("SELECT * FROM users WHERE status = ?", "active").Scan(ctx, &result)
```

## 查询组合与 Apply

`Apply` / `ApplyIf` 模式支持可复用的查询片段：

```go
// 定义可复用条件
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

// 使用它们
db.NewSelect().Model(&users).
	Apply(ActiveOnly, CreatedAfter(lastMonth)).
	Scan(ctx)

// 条件性应用
db.NewSelect().Model(&users).
	ApplyIf(keyword != "", func(q orm.SelectQuery) {
		q.Where(func(cb orm.ConditionBuilder) {
			cb.Contains("username", keyword)
		})
	}).
	Scan(ctx)
```

## DDL 操作

### 创建表

```go
_, err := db.NewCreateTable().
	Model((*User)(nil)).
	IfNotExists().
	Exec(ctx)
```

### 创建索引

```go
_, err := db.NewCreateIndex().
	Model((*User)(nil)).
	Index("idx_user_email").
	Column("email").
	Unique().
	IfNotExists().
	Exec(ctx)
```

### 其他 DDL

```go
db.NewDropTable().Model((*User)(nil)).IfExists().Exec(ctx)
db.NewTruncateTable().Model((*User)(nil)).Exec(ctx)
db.NewAddColumn().Model((*User)(nil)).ColumnExpr("phone VARCHAR(20)").Exec(ctx)
db.NewDropColumn().Model((*User)(nil)).Column("phone").Exec(ctx)
```

## 软删除支持

```go
// 仅查询已软删除的记录
db.NewSelect().Model(&users).WhereDeleted().Scan(ctx)

// 包含已软删除的记录
db.NewSelect().Model(&users).IncludeDeleted().Scan(ctx)
```

## 下一步

- [模型](./models) — 如何定义数据模型
- [泛型 CRUD](./crud) — 构建在查询构造器之上的高级 CRUD 操作
- [查询构造器（搜索标签）](./query-builder) — 基于搜索标签的自动查询构建
