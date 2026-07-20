---
sidebar_position: 8
---

# Schema 检查

VEF 内置了一个 schema 检查 service，以及一个可通过应用 API 读取数据库结构的内置资源。

内置实现只检查 primary data source。它通过 Atlas 支持 PostgreSQL、MySQL 和
SQLite。PostgreSQL 会使用已配置的 `vef.datasource.*.schema`，未配置时默认
为 `public`；MySQL 检查当前 `DATABASE()`；SQLite 检查 `main` schema。

## 模块输出

schema 模块会提供：

| 输出 | 含义 |
| --- | --- |
| `schema.Service` | schema 检查服务 |
| `sys/schema` | 内置 RPC 资源 |

## `schema.Service` 接口

公开的 schema service 暴露如下方法：

| 方法 | 返回类型 | 作用 |
| --- | --- | --- |
| `ListTables(ctx)` | `[]schema.Table` | 列出当前 schema / database 下的表 |
| `GetTableSchema(ctx, name)` | `*schema.TableSchema` | 检查单张表的详细结构 |
| `ListViews(ctx)` | `[]schema.View` | 列出视图 |

## 内置资源

schema 模块注册 `sys/schema` RPC 资源，挂载在 `/api` 下，使用标准请求
envelope（`resource`、`action`、`version`、`params`、`meta`）。没有任何
操作是公开的，也没有声明专门的权限点：每个 action 都继承 API 引擎默认的
Bearer 认证。

每个 action 都单独设置了 `max 60` 的限流上限。窗口长度未覆写，因此继承
`vef.api.rate_limit.period`（默认 `5m`）；限流按「操作 + 客户端 IP +
principal」计数，每个节点在进程内存中独立执行。

| Action | 访问 | 限流 | 入参 | 出参 |
| --- | --- | --- | --- | --- |
| `list_tables` | Bearer 认证 | `max 60` | 无 | `[]schema.Table` |
| `get_table_schema` | Bearer 认证 | `max 60` | `GetTableSchemaParams` | `schema.TableSchema` |
| `list_views` | Bearer 认证 | `max 60` | 无 | `[]schema.View` |

各 action 的用途：

- `list_tables` —— 枚举被检查 schema / database 下的表，返回表名、schema
  与注释。没有框架定义的入参。
- `get_table_schema` —— 检查单张表的详细结构：列、主键、索引、唯一键、
  外键和 check 约束。
- `list_views` —— 枚举被检查 schema / database 下的视图，包含定义 SQL 与
  输出列。没有框架定义的入参。

`GetTableSchemaParams`（`get_table_schema` 的入参）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `name` | `string` | 是 | 要检查的表名；缺失时校验失败，表不存在时返回 table-not-found 业务错误 |

源码中可见的行为语义：

- `get_table_schema` 会校验 `name` 必填（缺失时返回标准校验错误）
- 检查找不到的表在 service 内部以普通 Go sentinel `schema.ErrTableMissing`
  上抛，资源层将其映射为 `schema.ErrTableNotFound` 业务错误
  （`ErrCodeTableNotFound`，`2300`）
- RPC 响应使用标准 result envelope：HTTP 状态保持 `200`，业务错误通过
  body 的 `code` 传递

## 错误 API

| API | 含义 |
| --- | --- |
| `schema.ErrTableNotFound` | 请求的表不存在时返回的业务错误 |
| `schema.ErrCodeTableNotFound`（`2300`） | 表不存在对应的数值业务错误码 |
| `schema.ErrTableMissing` | 普通 Go sentinel，`Service` 检查找不到请求的表时返回——供直接使用 service 的 Go 调用方；RPC 层仍映射为 `ErrTableNotFound` |

## 按 Action 划分的响应结构

所有公开 schema DTO 的 JSON 字段名就是 wire contract。带 `omitempty` 的
字段在为空时会被省略：`schema`、`comment`、`primaryKey`、`indexes`、
`uniqueKeys`、`foreignKeys`、`checks`、`default`、`isPrimaryKey`、
`isAutoIncrement`、`predicate`、`hasExpressions`、`onUpdate`、`onDelete`、
`schema.PrimaryKey` 的 `name`，以及 `schema.View` 的 `schema` /
`comment` / `columns`。

### `list_tables` — `[]schema.Table`

被检查 schema / database 下每张表一条记录。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `name` | `string` | 表名 |
| `schema` | `string` | schema 名（PostgreSQL：被检查的 schema；MySQL：当前 database；SQLite：`main`）；为空时省略 |
| `comment` | `string` | 表注释；为空时省略 |

### `get_table_schema` — `schema.TableSchema`

单张表的完整结构。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `name` | `string` | 表名 |
| `schema` | `string` | schema 名；为空时省略 |
| `comment` | `string` | 表注释；为空时省略 |
| `columns` | `[]schema.Column` | 全部列，按数据库顺序 |
| `primaryKey` | `*schema.PrimaryKey` | 主键定义；表没有主键时省略 |
| `indexes` | `[]schema.Index` | 非唯一索引；为空时省略 |
| `uniqueKeys` | `[]schema.UniqueKey` | 唯一约束与唯一索引；为空时省略 |
| `foreignKeys` | `[]schema.ForeignKey` | 外键约束；为空时省略 |
| `checks` | `[]schema.Check` | check 约束；为空时省略 |

#### `schema.Column`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `name` | `string` | 列名 |
| `type` | `string` | 原始数据库类型 |
| `nullable` | `bool` | 是否可空 |
| `default` | `string` | 数据库上报的默认表达式原文；列没有默认值时省略 |
| `comment` | `string` | 列注释；为空时省略 |
| `isPrimaryKey` | `bool` | 是否属于主键；为 `false` 时省略 |
| `isAutoIncrement` | `bool` | 是否自动生成取值；为 `false` 时省略 |

`isAutoIncrement` 会识别 MySQL `AUTO_INCREMENT`、SQLite `AUTOINCREMENT`、
PostgreSQL identity column，以及 PostgreSQL `serial`、`bigserial`、
`smallserial` raw type。

#### `schema.PrimaryKey`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `name` | `string` | 主键约束名称；数据库未上报时省略 |
| `columns` | `[]string` | 主键列，按键序 |

#### `schema.Index`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `name` | `string` | 索引名 |
| `columns` | `[]string` | 索引列，按索引序；表达式成员以空字符串出现 |

#### `schema.UniqueKey`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `name` | `string` | 唯一键名 |
| `columns` | `[]string` | 唯一列，按索引序；表达式成员以空字符串出现 |
| `predicate` | `string` | 条件唯一索引（partial index）的谓词；为空时省略 |
| `hasExpressions` | `bool` | 唯一索引含表达式列时为 `true`（这类键不保证纯列组合的唯一性）；为 `false` 时省略 |

#### `schema.ForeignKey`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `name` | `string` | 外键名 |
| `columns` | `[]string` | 本地列 |
| `refTable` | `string` | 引用表 |
| `refColumns` | `[]string` | 引用列 |
| `onUpdate` | `string` | 更新策略：`CASCADE`、`SET NULL`、`SET DEFAULT`、`RESTRICT` 或 `NO ACTION`；数据库未上报时省略 |
| `onDelete` | `string` | 删除策略，词汇与 `onUpdate` 相同；数据库未上报时省略 |

#### `schema.Check`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `name` | `string` | check 约束名称 |
| `expr` | `string` | check 表达式 |

### `list_views` — `[]schema.View`

被检查 schema / database 下每个视图一条记录。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `name` | `string` | 视图名 |
| `schema` | `string` | schema 名；为空时省略 |
| `definition` | `string` | 视图定义 SQL |
| `comment` | `string` | 视图注释；为空时省略（SQLite 不提供视图注释） |
| `columns` | `[]string` | 输出列，按列序；为空时省略 |

视图列表使用数据库方言专属查询：PostgreSQL 和 MySQL 读取
`information_schema.views`；SQLite 读取 `sqlite_schema`，并排除内部
`sqlite_%` 视图。

## 最小请求示例

```json
{
  "resource": "sys/schema",
  "action": "get_table_schema",
  "version": "v1",
  "params": {
    "name": "sys_user"
  }
}
```

## 典型用途

这个功能很适合：

- 后台管理工具
- 内部开发者工具
- 需要 schema 感知的集成能力
- 需要数据库元信息的 MCP / prompt 工作流

## 下一步

继续阅读 [内置资源](../reference/built-in-resources)，看 `sys/schema` 在整个框架内置 RPC 资源体系里所处的位置。
