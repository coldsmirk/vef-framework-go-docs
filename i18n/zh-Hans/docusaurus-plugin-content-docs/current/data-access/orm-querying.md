---
sidebar_position: 5
---

# ORM：查询

VEF 将 Bun 封装成类型安全的流式查询构造器 API，提供类型安全的 SQL 构造、自动审计字段处理以及跨数据库方言支持。所有查询构造器都通过 `orm.DB` 访问。

本页覆盖数据读取：SELECT 子句、条件、连接、排序、分页、分组、锁定、执行与查询组合。相关主题见[表达式与聚合](./orm-expressions)、[写入操作](./orm-mutations)以及 [DDL 与公开接口面](./orm-ddl)。

## API surface 策略

ORM 包暴露两类已审计 surface。VEF-owned ORM method families 按
receiver/category 在本页以及模型、事务文档里说明；每个精确方法签名都列在
public API index 中。这类 surface 包括应用代码使用的 `orm.SelectQuery`、
`orm.InsertQuery`、`orm.UpdateQuery`、`orm.DeleteQuery`、`orm.MergeQuery`、
DDL、condition、expression、aggregate 和 window-builder contracts。

Bun pass-through surface 包括 `orm.BunSelectQuery`、`orm.BunInsertQuery`、
`orm.BunUpdateQuery` 和 `orm.BunDeleteQuery` 等 alias；这些方法遵循上游
[github.com/uptrace/bun](https://github.com/uptrace/bun) 行为，版本为 source
dependency 中固定的 `v1.2.18`。`Table`、`Field`、`Relation`、`Dialect` 等
Bun/schema aliases 遵循同一 pass-through policy。不要用 Bun aliases 推断 VEF
query-interface 行为：例如 VEF `orm.SelectQuery.Count` 返回 `int64`，而上游
Bun alias 的方法签名会在 public API index 中单独记录。嵌入 `fmt.Stringer`
的 query interfaces 公开 `String()` 用于 SQL/debug 渲染；精确签名记录在
public API index 中。

## 概述

| 分类 | 构造器 | 说明 |
| --- | --- | --- |
| 查询 | `db.NewSelect()` | 查询数据 |
| 插入 | `db.NewInsert()` | 创建记录 |
| 更新 | `db.NewUpdate()` | 修改记录 |
| 删除 | `db.NewDelete()` | 删除记录 |
| 合并 | `db.NewMerge()` | Upsert 操作 |
| 原始 SQL | `db.NewRaw(query, args...)` | 执行原始 SQL |

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
			cb.When(func(c orm.ConditionBuilder) { c.Equals("status", "active") }).Then(1).
				When(func(c orm.ConditionBuilder) { c.Equals("status", "pending") }).Then(2).
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

以下行描述的是 VEF `orm.SelectQuery` 执行方法，不是上游
`orm.BunSelectQuery` pass-through alias。

| 方法 | 返回值 | 用途 |
| --- | --- | --- |
| `Scan(ctx, dest...)` | `error` | 将行扫描到模型或目标 |
| `Exec(ctx, dest...)` | `sql.Result, error` | 执行但不扫描 |
| `Rows(ctx)` | `*sql.Rows, error` | 获取原始行迭代器 |
| `ScanAndCount(ctx)` | `int64, error` | 扫描 + 计算总数（用于分页）|
| `Count(ctx)` | `int64, error` | 仅计数 |
| `Exists(ctx)` | `bool, error` | 检查是否存在 |

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

## 下一步

- [ORM：表达式与聚合](./orm-expressions) — 表达式构造器、聚合与窗口函数、CTE 与集合操作
- [ORM：写入操作](./orm-mutations) — INSERT、UPDATE、DELETE、原始 SQL 与软删除行为
- [泛型 CRUD](./crud) — 构建在查询构造器之上的高级 CRUD 操作
