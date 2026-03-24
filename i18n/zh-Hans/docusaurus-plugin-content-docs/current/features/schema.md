---
sidebar_position: 7
---

# Schema 检查

VEF 内置了一个 schema 检查 service，以及一个可通过应用 API 读取数据库结构的内置资源。

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
| `ListTriggers(ctx)` | `[]schema.Trigger` | 列出触发器 |

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
| `list_triggers` | 无 | `[]schema.Trigger` | 列出当前触发器 |

源码层面的实现细节：

- 每个 action 当前都单独设置了 `Max = 60` 的限流上限
- `get_table_schema` 会校验 `name` 必填
- 表不存在时会映射到 `result.ErrCodeSchemaTableNotFound`

## 公共 Schema 类型

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
| `materialized` | `bool` | 是否物化视图 |

### `schema.Trigger`

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `name` | `string` | 触发器名 |
| `table` | `string` | 目标表（如适用） |
| `view` | `string` | 目标视图（如适用） |
| `actionTime` | `string` | before / after / instead-of 时机 |
| `events` | `[]string` | 触发事件，如 insert / update |
| `forEachRow` | `bool` | 行级还是语句级触发器 |
| `body` | `string` | 触发器定义体 |

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
