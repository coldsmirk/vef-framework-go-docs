---
sidebar_position: 7
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

schema 模块会注册：

| 资源 |
| --- |
| `sys/schema` |

当前 action：

| Action | 输入参数 | 输出类型 | 说明 |
| --- | --- | --- | --- |
| `list_tables` | 无 | `[]schema.Table` | 列出当前表 |
| `get_table_schema` | `name: string` | `schema.TableSchema` | 表不存在时会返回专门的 schema-table-not-found 业务错误 |
| `list_views` | 无 | `[]schema.View` | 列出当前视图 |

源码层面的实现细节：

- 每个 action 当前都单独设置了 `Max = 60` 的限流上限
- `get_table_schema` 会校验 `name` 必填
- 表不存在时会映射到 `schema.ErrCodeTableNotFound`（`ErrCodeTableNotFound`，对应 `schema.ErrTableNotFound` / `ErrTableNotFound` sentinel）
- RPC 响应仍使用标准 result envelope；业务错误写在 body 里，不通过不同 HTTP status 区分

## 错误 API

| API | 含义 |
| --- | --- |
| `schema.ErrTableNotFound` | 请求的表不存在时返回的业务错误 |
| `schema.ErrCodeTableNotFound` | 表不存在对应的数值业务错误码 |

## 公共 Schema 类型

审查说明：本页覆盖 53 public schema entries，其中包括 41 grouped schema field/method entries，分布在 10 schema receiver/type families；成组 DTO / service surface 包含 38 exported schema field entries 和 3 exported schema method entries。

所有公开 schema DTO 的 JSON 字段名就是 wire contract。带 `omitempty` 的字段在
为空时会被省略，例如 `schema`、`comment`、`primaryKey`、`indexes`、
`uniqueKeys`、`foreignKeys`、`checks`、`default`、`isPrimaryKey`、
`isAutoIncrement`、`onUpdate`、`onDelete`，以及 `schema.View` 的 `columns`。

### `schema.Table`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `name` | `string` | 表名 |
| `schema` | `string` | schema 名 |
| `comment` | `string` | 表注释 |

### `schema.TableSchema`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `name` | `string` | 表名 |
| `schema` | `string` | schema 名 |
| `comment` | `string` | 表注释 |
| `columns` | `[]schema.Column` | 全部列 |
| `primaryKey` | `*schema.PrimaryKey` | 主键定义 |
| `indexes` | `[]schema.Index` | 非唯一索引 |
| `uniqueKeys` | `[]schema.UniqueKey` | 唯一约束 |
| `foreignKeys` | `[]schema.ForeignKey` | 外键约束 |
| `checks` | `[]schema.Check` | check 约束 |

### `schema.Column`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `name` | `string` | 列名 |
| `type` | `string` | 原始数据库类型 |
| `nullable` | `bool` | 是否可空 |
| `default` | `string` | 默认表达式 |
| `comment` | `string` | 列注释 |
| `isPrimaryKey` | `bool` | 是否属于主键 |
| `isAutoIncrement` | `bool` | 是否自增 |

`isAutoIncrement` 会识别 MySQL `AUTO_INCREMENT`、SQLite `AUTOINCREMENT`、
PostgreSQL identity column，以及 PostgreSQL `serial`、`bigserial`、
`smallserial` raw type。

### `schema.PrimaryKey`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `name` | `string` | 主键名称 |
| `columns` | `[]string` | 主键列 |

### `schema.Index`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `name` | `string` | 索引名 |
| `columns` | `[]string` | 索引列 |

### `schema.UniqueKey`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `name` | `string` | 唯一键名 |
| `columns` | `[]string` | 唯一列 |

### `schema.ForeignKey`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `name` | `string` | 外键名 |
| `columns` | `[]string` | 本地列 |
| `refTable` | `string` | 引用表 |
| `refColumns` | `[]string` | 引用列 |
| `onUpdate` | `string` | 更新策略 |
| `onDelete` | `string` | 删除策略 |

### `schema.Check`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `name` | `string` | check 约束名称 |
| `expr` | `string` | check 表达式 |

### `schema.View`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `name` | `string` | 视图名 |
| `schema` | `string` | schema 名 |
| `definition` | `string` | 视图定义 SQL |
| `comment` | `string` | 视图注释 |
| `columns` | `[]string` | 输出列 |

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
