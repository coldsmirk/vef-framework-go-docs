---
sidebar_position: 7
---

# 钩子

VEF 里实际上有两层 hook 面：

- CRUD builder 提供的操作级 hook，用于框架托管的 create / update / delete / import / export 流程
- Bun/ORM 级模型 hook，用于更底层的查询生命周期拦截

这两层 hook 解决的问题不同，不应混为一谈。

## Hook 家族总览

| Hook 家族 | 运行位置 | 作用范围 | 常见用途 |
| --- | --- | --- | --- |
| CRUD hook | CRUD builder 内部 | 单个 CRUD 接口 | 业务约束、事务内副作用、文件提升协同 |
| Bun 模型 hook | ORM 模型生命周期 | 单个模型与单种查询类型 | 底层查询修改、模型生命周期行为、持久化侧检查 |

## CRUD Hook 面

CRUD builder 暴露的 hook API 如下：

| 操作 | 前置 hook | 后置 hook |
| --- | --- | --- |
| `Create` | `WithPreCreate(...)` | `WithPostCreate(...)` |
| `Update` | `WithPreUpdate(...)` | `WithPostUpdate(...)` |
| `Delete` | `WithPreDelete(...)` | `WithPostDelete(...)` |
| `CreateMany` | `WithPreCreateMany(...)` | `WithPostCreateMany(...)` |
| `UpdateMany` | `WithPreUpdateMany(...)` | `WithPostUpdateMany(...)` |
| `DeleteMany` | `WithPreDeleteMany(...)` | `WithPostDeleteMany(...)` |
| `Export` | `WithPreExport(...)` | 没有 post-export hook |
| `Import` | `WithPreImport(...)` | `WithPostImport(...)` |

## CRUD Hook 签名

### 单条 create / update / delete

| Hook | 签名概要 |
| --- | --- |
| `PreCreate` | `func(model *TModel, params *TParams, query orm.InsertQuery, ctx fiber.Ctx, tx orm.DB) error` |
| `PostCreate` | `func(model *TModel, params *TParams, ctx fiber.Ctx, tx orm.DB) error` |
| `PreUpdate` | `func(oldModel, model *TModel, params *TParams, query orm.UpdateQuery, ctx fiber.Ctx, tx orm.DB) error` |
| `PostUpdate` | `func(oldModel, model *TModel, params *TParams, ctx fiber.Ctx, tx orm.DB) error` |
| `PreDelete` | `func(model *TModel, query orm.DeleteQuery, ctx fiber.Ctx, tx orm.DB) error` |
| `PostDelete` | `func(model *TModel, ctx fiber.Ctx, tx orm.DB) error` |

### 批量 create / update / delete

| Hook | 签名概要 |
| --- | --- |
| `PreCreateMany` | `func(models []TModel, paramsList []TParams, query orm.InsertQuery, ctx fiber.Ctx, tx orm.DB) error` |
| `PostCreateMany` | `func(models []TModel, paramsList []TParams, ctx fiber.Ctx, tx orm.DB) error` |
| `PreUpdateMany` | `func(oldModels, models []TModel, paramsList []TParams, query orm.UpdateQuery, ctx fiber.Ctx, tx orm.DB) error` |
| `PostUpdateMany` | `func(oldModels, models []TModel, paramsList []TParams, ctx fiber.Ctx, tx orm.DB) error` |
| `PreDeleteMany` | `func(models []TModel, query orm.DeleteQuery, ctx fiber.Ctx, tx orm.DB) error` |
| `PostDeleteMany` | `func(models []TModel, ctx fiber.Ctx, tx orm.DB) error` |

### 导出与导入

| Hook | 签名概要 |
| --- | --- |
| `PreExport` | `func(models []TModel, search TSearch, ctx fiber.Ctx, db orm.DB) error` |
| `PreImport` | `func(models []TModel, query orm.InsertQuery, ctx fiber.Ctx, tx orm.DB) error` |
| `PostImport` | `func(models []TModel, ctx fiber.Ctx, tx orm.DB) error` |

## CRUD Hook 与事务边界

最重要的一条规则是：CRUD 写操作本身已经在事务中执行。CRUD hook 拿到的 `orm.DB` 就是当前事务作用域内的 `tx`，因此你额外做的数据库操作天然处在同一事务里。

示例：

```go
crud.NewCreate[User, UserParams]().
	WithPostCreate(func(model *User, params *UserParams, ctx fiber.Ctx, tx orm.DB) error {
		_, err := tx.NewInsert().Model(&AuditLog{
			UserID: model.ID,
			Action: "created",
		}).Exec(ctx.Context())
		return err
	})
```

## CRUD Hook 错误行为

如果 CRUD hook 返回错误：

- 当前操作失败
- 外层事务回滚
- 框架按照正常 result 错误链路返回错误

因此，CRUD hook 很适合承载“必须原子化生效”的业务约束。

## CRUD Hook 与文件提升

Create、Update、Delete 这几个 builder 还会和 storage promoter 协同工作。

这意味着：

- 文件提升和 CRUD 操作处于同一生命周期
- Update 回滚时可以恢复被替换文件
- Delete 成功后可以清理已提升文件

因此 hook 与文件提升是共享同一条 CRUD 生命周期的。

## Bun 模型 Hook 面

在 ORM 层，VEF 也暴露了 Bun hook 接口：

| Hook 接口 | 触发时机 |
| --- | --- |
| `orm.BeforeSelectHook` | `SELECT` 前 |
| `orm.AfterSelectHook` | `SELECT` 后 |
| `orm.BeforeInsertHook` | `INSERT` 前 |
| `orm.AfterInsertHook` | `INSERT` 后 |
| `orm.BeforeUpdateHook` | `UPDATE` 前 |
| `orm.AfterUpdateHook` | `UPDATE` 后 |
| `orm.BeforeDeleteHook` | `DELETE` 前 |
| `orm.AfterDeleteHook` | `DELETE` 后 |

这些 hook 是定义在模型类型上的，作用层级是 ORM 生命周期，而不是 API action 生命周期。

## 什么时候用 CRUD Hook

以下场景适合 CRUD hook：

- 额外业务步骤紧贴某一个 CRUD 动作
- 对外接口语义仍然保持 CRUD
- 额外逻辑必须参与同一事务
- hook 需要同时访问 params 和 model 状态

## 什么时候用 Bun 模型 Hook

以下场景更适合 Bun 模型 hook：

- 逻辑本身属于模型生命周期
- 这段行为应该在 API 层之外也生效
- hook 需要直接修改或检查底层 Bun query
- 关注点更偏持久化，而不是 endpoint 编排

## 什么时候不该用 Hook

以下场景 hook 反而不合适：

- 这个动作已经不再是标准 CRUD
- 一个接口在编排多个互不相关的流程
- 该动作语义更适合一个显式业务命令接口
- hook 越堆越多，导致行为拆散后难以理解

这类场景通常应该改用自定义 handler。

## 实践建议

- CRUD hook 保持短小，聚焦单个操作
- Bun 模型 hook 聚焦持久化层行为
- 事务内业务步骤优先放 CRUD hook
- 如果某个行为应该在所有模型访问路径里生效，优先考虑模型 hook
- 如果一个资源开始堆很多 CRUD hook，就该反思是不是应该改成自定义 handler

## 下一步

继续阅读 [验证](./validation) 和 [错误处理](./error-handling)，看请求失败和业务失败是如何向客户端暴露的。
