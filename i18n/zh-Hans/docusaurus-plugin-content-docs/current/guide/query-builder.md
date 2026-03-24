---
sidebar_position: 6
---

# 查询构建

VEF 的查询构建主要围绕 typed search 结构体、`search` 标签和 CRUD 的 find 类扩展点展开。目标是把查询规则放在字段定义附近，而不是把一堆字符串条件散落在 handler 里。

## Search 结构体模型

最常见的形态如下：

```go
type UserSearch struct {
	api.P

	ID       string `json:"id" search:"eq"`
	Keyword  string `json:"keyword" search:"contains,column=username|email"`
	IsActive *bool  `json:"isActive" search:"eq,column=is_active"`
}
```

`search` 标签描述的是：一个字段如何被翻译成一个或多个 SQL 条件。

## 不写 `search` 标签时的默认行为

如果一个字段完全没有 `search` 标签：

- 框架仍然会把它纳入搜索结构解析
- 默认操作符是 `eq`
- 默认列名是字段名的 snake_case 形式

也就是说，这个字段：

```go
Age int
```

等价于：

```go
Age int `search:"eq,column=age"`
```

## `search` 标签语法

`search` 标签支持以下写法：

| 模式 | 含义 |
| --- | --- |
| `search:"eq"` | 只写操作符 |
| `search:"contains,column=username|email"` | 写操作符和目标列 |
| `search:"operator=gte,column=price"` | 完整 key/value 形式 |
| `search:"operator=in,params=delimiter:|,type:int"` | 携带额外参数 |
| `search:"dive"` | 递归进入嵌套结构体 |
| `search:"-"` | 完全忽略该字段 |

支持的标签属性：

| 属性 | 含义 |
| --- | --- |
| 默认值或 `operator` | 查询操作符 |
| `column` | 一个或多个目标列，列之间用 `|` 分隔 |
| `alias` | 表别名，用于列限定 |
| `params` | 操作符的额外参数 |
| `dive` | 递归进入嵌套结构 |

## 支持的操作符

框架当前支持以下全部操作符：

### 比较操作符

| 操作符 | 含义 |
| --- | --- |
| `eq` | 等于 |
| `neq` | 不等于 |
| `gt` | 大于 |
| `gte` | 大于等于 |
| `lt` | 小于 |
| `lte` | 小于等于 |

### 区间操作符

| 操作符 | 含义 |
| --- | --- |
| `between` | 落在区间内 |
| `notBetween` | 落在区间外 |

### 集合操作符

| 操作符 | 含义 |
| --- | --- |
| `in` | 属于集合 |
| `notIn` | 不属于集合 |

### Null 操作符

| 操作符 | 含义 |
| --- | --- |
| `isNull` | 生成 `IS NULL` |
| `isNotNull` | 生成 `IS NOT NULL` |

### 字符串匹配操作符

| 操作符 | 含义 |
| --- | --- |
| `contains` | 包含子串 |
| `notContains` | 不包含子串 |
| `startsWith` | 以前缀开头 |
| `notStartsWith` | 不以前缀开头 |
| `endsWith` | 以后缀结尾 |
| `notEndsWith` | 不以后缀结尾 |

### 大小写不敏感字符串操作符

| 操作符 | 含义 |
| --- | --- |
| `iContains` | 大小写不敏感包含 |
| `iNotContains` | 大小写不敏感不包含 |
| `iStartsWith` | 大小写不敏感前缀匹配 |
| `iNotStartsWith` | 大小写不敏感前缀不匹配 |
| `iEndsWith` | 大小写不敏感后缀匹配 |
| `iNotEndsWith` | 大小写不敏感后缀不匹配 |

## 多列搜索

一个字段可以通过 `|` 同时命中多个列。

示例：

```go
Keyword string `search:"contains,column=username|email|mobile"`
```

这非常适合关键词搜索。

## 使用 `dive` 处理嵌套结构

`dive` 不是查询操作符，而是一个解析指令，表示“继续进入嵌套结构体”。

示例：

```go
type UserSearch struct {
	Name string `search:"column=user_name,operator=contains"`
}

type OrderSearch struct {
	api.P

	User UserSearch `search:"dive"`
}
```

## Alias

当查询需要带表别名时，可以使用 `alias`：

```go
Name string `search:"alias=u,column=name,operator=contains"`
```

这在 join 查询中非常有用。

## 操作符参数

有些操作符支持通过 `params=...` 传入额外参数。

当前常见参数键包括：

| 参数键 | 含义 |
| --- | --- |
| `delimiter` | 解析字符串区间或集合时使用的分隔符 |
| `type` | 显式解析类型，例如 `int`、`dec`、`date`、`datetime`、`time` |

## `between` 的输入形态

`between` 和 `notBetween` 支持多种输入形式：

| 输入形态 | 示例 |
| --- | --- |
| `monad.Range[T]` 风格结构体 | `monad.Range[int]{Start: 1, End: 10}` |
| 双元素切片 | `[]int{1, 10}` |
| 分隔字符串 | `"1,10"` |

字符串输入可通过 `params` 控制解析方式。

示例：

```go
Price string `search:"operator=between,column=price,params=type:int,delimiter:,"`
DateRange string `search:"operator=between,column=created_at,params=type:date,delimiter:|"`
```

## `in` / `notIn` 的输入形态

集合操作符支持：

| 输入形态 | 示例 |
| --- | --- |
| slice 字段 | `[]string{"a", "b"}` |
| 分隔字符串 | `"a,b,c"` |
| 自定义分隔符字符串 | `"1|2|3"` + `params=delimiter:|,type:int` |

## 排序

排序通常通过 `crud.Sortable` 放在 `meta` 中：

```go
type QueryMeta struct {
	api.M
	crud.Sortable
}
```

`crud.Sortable` 的结构是：

| 字段 | 含义 |
| --- | --- |
| `Sort []sortx.OrderSpec` | 一组排序规则 |

每个 `sortx.OrderSpec` 可以表达：

| 属性 | 含义 |
| --- | --- |
| `Column` | 目标列 |
| `Direction` | 升序或降序 |
| `NullsOrder` | null 排序位置 |

CRUD 的 find builder 可以自动应用这些排序规则。

## 分页

分页使用 `page.Pageable`：

```go
type QueryMeta struct {
	api.M
	page.Pageable
}
```

`FindPage` 会负责：

- 规范化 page 和 size
- 应用分页限制
- 返回 `page.Page[T]`

需要注意：

- `page.Pageable` 是从 `meta` 解码的
- 对 REST handler 来说，`?page=1&size=20` 只会落到原始 `params`，不会自动填充 typed `page.Pageable`

## 数据权限

很多读 builder 默认会自动应用请求级数据权限过滤。

这意味着：

- 实际查询条件不只来自 `search` 标签和显式条件
- 数据权限可能会透明地再附加一层过滤
- 如果你确实要绕开它，需要显式在对应 builder 上禁用

## 查询逃逸口

当 `search` 标签已经不够表达查询逻辑时，CRUD find builder 还支持这些扩展点：

| 方法 | 适合的用途 |
| --- | --- |
| `WithCondition(...)` | 追加 `WHERE` 条件 |
| `WithRelation(...)` | 增加关联 join |
| `WithDefaultSort(...)` | 设置默认排序 |
| `WithQueryApplier(...)` | 用 typed 方式直接修改查询对象 |
| `WithSelect(...)` / `WithSelectAs(...)` | 明确控制 `SELECT` 列表 |

对于树形 API，这些扩展点还可以定向作用到 `QueryBase`、`QueryRecursive`、`QueryRoot` 等不同查询阶段。

## 常见模式

### 简单等值 + 关键词搜索

```go
type UserSearch struct {
	api.P

	ID      string `json:"id" search:"eq"`
	Keyword string `json:"keyword" search:"contains,column=username|email"`
}
```

### 区间与集合筛选

```go
type ProductSearch struct {
	api.P

	PriceRange string `json:"priceRange" search:"operator=between,column=price,params=type:int,delimiter:,"`
	Statuses   string `json:"statuses" search:"operator=in,column=status,params=delimiter:|"`
}
```

### 嵌套搜索

```go
type UserSearch struct {
	Name string `search:"column=user_name,operator=contains"`
}

type OrderSearch struct {
	api.P

	User UserSearch `search:"dive"`
}
```

## 实践建议

- 每个资源都定义专属 search 结构体
- 普通筛选优先用 `search` 标签表达，把规则放在字段旁边
- 关键词搜索优先显式写多列映射，不要偷偷塞自定义 SQL
- 排序和分页走 `meta`
- 只有当标签已经无法表达时，才使用 `WithQueryApplier(...)`
- 让查询契约留在类型定义里，而不是埋进 handler 代码

## 下一步

继续阅读 [钩子](./hooks)，如果你的查询或变更还需要生命周期感知逻辑，就会接到那一层。
