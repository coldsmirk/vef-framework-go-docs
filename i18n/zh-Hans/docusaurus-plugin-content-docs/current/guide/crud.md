---
sidebar_position: 2
---

# 泛型 CRUD

`crud` 包是 VEF 最重要的用户层能力之一。它把类型化模型和类型化请求结构体变成可复用的 API 操作，并内置事务、验证、数据权限、文件提升和结果格式化。

## 基本模式

最常见的写法，是把 CRUD provider 直接嵌入资源结构体：

```go
type UserResource struct {
	api.Resource

	crud.FindPage[User, UserSearch]
	crud.Create[User, UserParams]
	crud.Update[User, UserParams]
	crud.Delete[User]
}

func NewUserResource() api.Resource {
	return &UserResource{
		Resource: api.NewRPCResource("sys/user"),
		FindPage: crud.NewFindPage[User, UserSearch]().RequiredPermission("sys:user:query"),
		Create:   crud.NewCreate[User, UserParams]().RequiredPermission("sys:user:create"),
		Update:   crud.NewUpdate[User, UserParams]().RequiredPermission("sys:user:update"),
		Delete:   crud.NewDelete[User]().RequiredPermission("sys:user:delete"),
	}
}
```

框架之所以能自动收集这些 CRUD builder，是因为它们本身实现了 `api.OperationsProvider`。

Grouped-family audit 固定了 315 grouped CRUD builder entries，覆盖 27
receiver families：其中 36 public field entries、279 public method entries。
这些 entries 覆盖预置 builder families、通用 builder controls、query-shaping
helpers、批量 params、tree/data-option DTOs、import/export customization 和
processor hooks；verifier 会锁定排序后的签名和 receiver/type 分布。

### 完整的 Model / Params / Search 定义

```go
// Model — 持久化层
type User struct {
	orm.FullAuditedModel

	Username     string `json:"username" bun:"username"`
	Email        string `json:"email" bun:"email"`
	DepartmentID string `json:"departmentId" bun:"department_id"`
	IsActive     bool   `json:"isActive" bun:"is_active"`
	Avatar       string `json:"avatar" bun:"avatar" meta:"uploaded_file"`
}

// Params — 写操作请求体
type UserParams struct {
	Username     string `json:"username" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	DepartmentID string `json:"departmentId"`
	IsActive     *bool  `json:"isActive"`
	Avatar       string `json:"avatar"`
}

// Search — 读操作查询条件
type UserSearch struct {
	Keyword      string  `json:"keyword" search:"contains,column=username|email"`
	DepartmentID *string `json:"departmentId" search:"eq"`
	IsActive     *bool   `json:"isActive" search:"eq"`
}
```

## 泛型参数含义

大多数 CRUD builder 只会用到下面这几种泛型形状：

| 泛型 | 含义 | 常见类型 |
| --- | --- | --- |
| `TModel` | 持久化模型，最终会被查询或写入数据库 | `User`、`Role`、`Flow` |
| `TParams` | 写操作参数类型，从 `Request.Params` 解码 | `UserParams`、`CreateUserParams` |
| `TSearch` | 读操作搜索参数类型，从 `Request.Params` 解码 | `UserSearch`、`RoleSearch` |

各类 builder 对这些泛型的使用方式如下：

| Builder 家族 | 泛型形状 | 含义 |
| --- | --- | --- |
| 单条写操作 | `Create[TModel, TParams]`、`Update[TModel, TParams]` | `TParams` 会被拷贝进模型再持久化 |
| 批量写操作 | `CreateMany[TModel, TParams]`、`UpdateMany[TModel, TParams]` | 框架会把 `TParams` 包装进批量 params 类型 |
| 读操作 | `FindOne[TModel, TSearch]`、`FindPage[TModel, TSearch]` 等 | `TModel` 定义查询目标，`TSearch` 定义查询条件 |
| 删除操作 | `Delete[TModel]`、`DeleteMany[TModel]` | 删除根据主键输入工作，不需要额外的 `TParams` |
| 导出 | `Export[TModel, TSearch]` | 导出先执行读查询，再把结果渲染成文件 |
| 导入 | `Import[TModel]` | 文件中的行会直接解码成模型 |

## 预置 Builder 总览

| Builder | 默认 RPC action | 可用 REST override 示例 | 输入契约 | 输出契约 | 典型用途 |
| --- | --- | --- | --- | --- | --- |
| `NewCreate[TModel, TParams]` | `create` | `post` | `params` 中的 `TParams` | 主键 map | 创建单条记录 |
| `NewUpdate[TModel, TParams]` | `update` | `put` | `params` 中的 `TParams`，且需要包含主键字段 | 成功结果 | 更新单条记录 |
| `NewDelete[TModel]` | `delete` | `delete` | `params` 中的原始主键值 | 成功结果 | 删除单条记录 |
| `NewCreateMany[TModel, TParams]` | `create_many` | `post many` | `CreateManyParams[TParams]`，包含 `list` | 主键 map 列表 | 批量创建 |
| `NewUpdateMany[TModel, TParams]` | `update_many` | `put many` | `UpdateManyParams[TParams]`，包含 `list` | 成功结果 | 批量更新 |
| `NewDeleteMany[TModel]` | `delete_many` | `delete many` | `DeleteManyParams`，包含 `pks` | 成功结果 | 批量删除 |
| `NewFindOne[TModel, TSearch]` | `find_one` | `get one` | `params` 中的 `TSearch` | 单个模型 | 单条查询 |
| `NewFindAll[TModel, TSearch]` | `find_all` | `get` | `params` 中的 `TSearch` | `[]TModel` | 不带分页元数据的列表查询 |
| `NewFindPage[TModel, TSearch]` | `find_page` | `get page` | `params` 中的 `TSearch` + `meta` 中的 `page.Pageable` | `page.Page[T]` | 后台分页列表 |
| `NewFindOptions[TModel, TSearch]` | `find_options` | `get options` | `params` 中的 `TSearch` + `meta` 中的 `DataOptionConfig` | `[]DataOption` | 下拉选项 |
| `NewFindTree[TModel, TSearch](treeBuilder)` | `find_tree` | `get tree` | `params` 中的 `TSearch` | 分层 `[]TModel` | 树形数据 |
| `NewFindTreeOptions[TModel, TSearch]` | `find_tree_options` | `get tree-options` | `params` 中的 `TSearch` + `meta` 中的 `DataOptionConfig` | `[]TreeDataOption` | 树形选项 |
| `NewExport[TModel, TSearch]` | `export` | `get export` | `params` 中的 `TSearch` + `meta` 中的导出格式 | 文件下载 | Excel / CSV 导出 |
| `NewImport[TModel]` | `import` | `post import` | multipart 上传文件 + `meta` 中的导入格式 | `{total: n}` | Excel / CSV 导入 |

当前源码中导出的 `RESTAction*` constants 包含 `post /`、`put /:id`、
`get /page` 这类 slash route pattern。这些常量会出现在公开 API 索引中，
但 public `api.ValidateActionName` / `api.NewRESTResource` 校验只接受小写
HTTP verb，以及可选的 kebab-case 子资源。因此在当前源码下，直接把
`api.KindREST` 传给预置 CRUD 构造器，会在默认 action 校验时 panic。若要
构建 REST 风格 CRUD 操作，应先用普通构造器创建 builder，再显式设置一个
可通过公开 REST 语法的 action：

```go
crud.NewFindPage[User, UserSearch]().
	ResourceKind(api.KindREST).
	Action("get page")
```

导出的 REST action 常量值如下：

| 常量 | 值 |
| --- | --- |
| `RESTActionCreate` | `post /` |
| `RESTActionUpdate` | `put /:id` |
| `RESTActionDelete` | `delete /:id` |
| `RESTActionCreateMany` | `post /many` |
| `RESTActionUpdateMany` | `put /many` |
| `RESTActionDeleteMany` | `delete /many` |
| `RESTActionFindOne` | `get /:id` |
| `RESTActionFindAll` | `get /` |
| `RESTActionFindPage` | `get /page` |
| `RESTActionFindOptions` | `get /options` |
| `RESTActionFindTree` | `get /tree` |
| `RESTActionFindTreeOptions` | `get /tree/options` |
| `RESTActionImport` | `post /import` |
| `RESTActionExport` | `get /export` |

## 共享 Builder 配置

每个 CRUD builder 都继承了 `Builder[T]` 的通用控制项：

| 方法 | 作用 |
| --- | --- |
| `ResourceKind(kind)` | 切换为 RPC 或 REST 命名/校验规则 |
| `Action(action)` | 覆盖默认 action 名 |
| `Public()` | 将该操作标记为无需认证 |
| `RequiredPermission(token)` | 要求调用方具备某个权限点 |
| `Timeout(duration)` | 设置请求超时 |
| `EnableAudit()` | 为该操作启用审计记录 |
| `RateLimit(max, period)` | 对该操作单独配置限流 |

一个容易忽略的点：

- `Action(...)` 会按当前 `ResourceKind(...)` 进行校验
- 如果你要覆盖 REST action，应该先调用 `ResourceKind(api.KindREST)`，
  并使用公开 REST 语法：`method` 或 `method kebab-sub-resource`

## Find 系列共享配置

所有读操作 builder 都建立在 `Find[...]` 之上，因此它们共享更丰富的查询塑形能力：

| 方法 | 作用 |
| --- | --- |
| `WithProcessor(...)` | 在响应序列化前，对查询结果做后处理 |
| `WithOptions(...)` | 追加底层 `FindOperationOption` |
| `WithSelect(column)` | 向 `SELECT` 增加一列 |
| `WithSelectAs(column, alias)` | 向 `SELECT` 增加一列并起别名 |
| `WithDefaultSort(...)` | 当请求里没有动态排序时，设置默认排序 |
| `WithCondition(...)` | 用 `orm.ConditionBuilder` 增加 `WHERE` 条件 |
| `DisableDataPerm()` | 禁用默认的数据权限过滤 |
| `WithRelation(...)` | 通过 `orm.RelationSpec` 增加关联 join |
| `WithAuditUserNames(userModel, nameColumn...)` | 自动 join 审计用户，补齐创建人/更新人名字 |
| `WithQueryApplier(...)` | 用类型安全的方式直接修改查询对象 |

绝大多数 Find builder 的运行时默认行为：

- 会自动应用 `TSearch` 上的 `search:"..."` 标签
- 默认启用数据权限过滤
- 如果模型有单主键，默认按主键倒序排序
- 如果没有单主键，但有 `created_at`，则回退为 `created_at DESC`

### Find 控制示例

#### WithCondition

```go
crud.NewFindPage[User, UserSearch]().
	WithCondition(func(cb orm.ConditionBuilder) {
		cb.IsTrue("is_active")
	})
```

#### WithQueryApplier

```go
crud.NewFindPage[User, UserSearch]().
	WithQueryApplier(func(q orm.SelectQuery, search UserSearch, ctx fiber.Ctx) error {
		if search.DepartmentID != nil {
			q.Where(func(cb orm.ConditionBuilder) {
				cb.Equals("department_id", *search.DepartmentID)
			})
		}
		return nil
	})
```

#### WithRelation

```go
crud.NewFindPage[User, UserSearch]().
	WithRelation(&orm.RelationSpec{
		Model: (*Department)(nil),
		SelectedColumns: []orm.ColumnInfo{
			{Name: "name", Alias: "department_name"},
		},
	})
```

#### WithAuditUserNames

```go
// 自动 join sys_user 填充 created_by_name 和 updated_by_name
crud.NewFindPage[User, UserSearch]().
	WithAuditUserNames((*User)(nil), "username")
```

#### WithProcessor

```go
crud.NewFindPage[User, UserSearch]().
	WithProcessor(func(users []User, search UserSearch, ctx fiber.Ctx) any {
		// 在序列化前转换模型
		result := make([]UserDTO, len(users))
		for i, u := range users {
			result[i] = toDTO(u)
		}
		return result
	})
```

#### WithDefaultSort

```go
crud.NewFindPage[User, UserSearch]().
	WithDefaultSort(&sortx.OrderSpec{Column: "created_at", Direction: sortx.OrderDesc})
```

### Tree Builder 的 QueryPart

树形 builder 使用递归 CTE，因此部分配置项可以作用在不同查询阶段：

| QueryPart | 含义 |
| --- | --- |
| `QueryRoot` | 最外层最终查询 |
| `QueryBase` | 递归 CTE 的起始查询 |
| `QueryRecursive` | 递归 CTE 的递归分支 |
| `QueryAll` | 作用于所有查询部分 |

对于 `FindTree` 和 `FindTreeOptions`，一些方法的默认作用范围和普通查询不同：

- `WithCondition(...)` 默认作用于 `QueryBase`
- `WithQueryApplier(...)` 默认作用于 `QueryBase`
- `WithSelect(...)`、`WithSelectAs(...)`、`WithRelation(...)` 默认同时作用于 `QueryBase` 和 `QueryRecursive`

## 读操作 Builder

### `FindOne[TModel, TSearch]`

当资源应返回一条记录时使用 `FindOne`。

| 维度 | 说明 |
| --- | --- |
| 泛型 | `TModel` 是查询目标模型，`TSearch` 定义筛选条件 |
| 输入 | `params` 中的 `TSearch`，以及 `meta` 中的原始 `api.Meta` |
| 输出 | 单个 `TModel`，或 `WithProcessor(...)` 转换后的结果 |
| 默认行为 | 查询模型列，并自动加上 `LIMIT 1` |
| 常见配置 | `WithCondition`、`WithRelation`、`WithQueryApplier`、`WithAuditUserNames` |

当这个接口本质上仍然是“查询”，而不是“固定上下文的元数据读取”时，优先使用它。

### `FindAll[TModel, TSearch]`

当你需要一个不带分页元数据的列表时使用 `FindAll`。

| 维度 | 说明 |
| --- | --- |
| 泛型 | `TModel` 是结果模型，`TSearch` 定义筛选条件 |
| 输入 | `params` 中的 `TSearch`，`meta` 中的 `api.Meta` |
| 输出 | `[]TModel`，或者 `WithProcessor(...)` 返回的结果切片 |
| 默认行为 | 会应用安全上限 `maxQueryLimit`，并在空结果时返回空切片而不是 `nil` |
| 常见配置 | 共享 Find 配置，尤其是 `WithDefaultSort`、`WithCondition`、`WithRelation`、`WithQueryApplier` |

### `FindPage[TModel, TSearch]`

绝大多数后台列表接口都应该从 `FindPage` 开始。

| 维度 | 说明 |
| --- | --- |
| 泛型 | `TModel` 是列表项模型，`TSearch` 定义查询条件 |
| 输入 | `params` 中的 `TSearch`，`meta` 中的 `page.Pageable`，以及额外的 `api.Meta` |
| 输出 | `page.Page[T]` |
| 默认行为 | 自动分页、统计总数，并规范化分页参数 |
| 特有配置 | `WithDefaultPageSize(size)` 用于设置默认页大小 |

当调用方需要 `total`、页码、页大小和列表项一起返回时，优先使用它。

### `FindOptions[TModel, TSearch]`

当你需要轻量选项列表，例如下拉框时，使用 `FindOptions`。

| 维度 | 说明 |
| --- | --- |
| 泛型 | `TModel` 是来源模型，`TSearch` 定义筛选条件 |
| 输入 | `params` 中的 `TSearch`，`meta` 中的 `DataOptionConfig` |
| 输出 | `[]DataOption` |
| 默认行为 | 将结果映射为 `label`、`value`、`description` 和可选 `meta` |
| 特有配置 | `WithDefaultColumnMapping(mapping)` 用于设置 label / value / description / meta 列映射的默认值 |

`DataOptionConfig` 来自 `meta`，可配置字段包括：

| 字段 | 作用 |
| --- | --- |
| `labelColumn` | `label` 的来源列 |
| `valueColumn` | `value` 的来源列 |
| `descriptionColumn` | 可选的 `description` 来源列 |
| `metaColumns` | 额外塞进 `meta` 对象的列 |

默认值：

- `labelColumn` 默认是 `name`
- `valueColumn` 默认是 `id`

### `FindTree[TModel, TSearch]`

当领域数据本身是树形结构，且你希望返回嵌套模型时，使用 `FindTree`。

构造器形态和其他 Find builder 不同：

```go
func buildCategoryTree(flat []Category) []Category {
	adapter := tree.Adapter[Category]{
		GetID: func(c Category) string {
			return c.ID
		},
		GetParentID: func(c Category) *string {
			return c.ParentID
		},
		GetChildren: func(c Category) []Category {
			return c.Children
		},
		SetChildren: func(c *Category, children []Category) {
			c.Children = children
		},
	}

	return tree.Build(flat, adapter)
}

crud.NewFindTree[Category, CategorySearch](buildCategoryTree)
```

`tree.Build` 的签名是 `func([]T, tree.Adapter[T]) []T`，而 `NewFindTree`
需要 `func([]T) []T`；应传入闭包里包含模型 adapter 的 wrapper，不要直接传
`tree.Build`。

| 维度 | 说明 |
| --- | --- |
| 泛型 | `TModel` 是树节点模型，`TSearch` 定义筛选条件 |
| 输入 | `params` 中的 `TSearch`，`meta` 中的 `api.Meta` |
| 输出 | 分层 `[]TModel` |
| 默认行为 | 先构建递归 CTE 拉平查询结果，再交给 `treeBuilder` 组装成树 |
| 特有配置 | `WithIDColumn(name)`、`WithParentIDColumn(name)` |

默认值：

- 节点 ID 列默认为 `id`
- 父节点列默认为 `parent_id`

### `FindTreeOptions[TModel, TSearch]`

当你需要树形选项，而不是完整模型结构时，使用 `FindTreeOptions`。

| 维度 | 说明 |
| --- | --- |
| 泛型 | `TModel` 是来源模型，`TSearch` 定义筛选条件 |
| 输入 | `params` 中的 `TSearch`，`meta` 中的 `DataOptionConfig` |
| 输出 | `[]TreeDataOption` |
| 默认行为 | 通过递归 CTE 取到树节点，再转换成 `TreeDataOption` 嵌套结构 |
| 特有配置 | `WithDefaultColumnMapping(...)`、`WithIDColumn(...)`、`WithParentIDColumn(...)` |

当客户端只需要 `label` / `value` / `children` 这类树形选项载荷时，优先用它，而不是直接暴露完整模型。

## 写操作 Builder

### `Create[TModel, TParams]`

单条创建时使用 `Create`。

| 维度 | 说明 |
| --- | --- |
| 泛型 | `TModel` 是持久化模型，`TParams` 是写入参数类型 |
| 输入 | `params` 中的 `TParams` |
| 输出 | 新建记录的主键 map |
| 默认行为 | 把 params 拷贝进新模型，处理文件提升，在事务内插入记录 |
| 特有配置 | `WithPreCreate(...)`、`WithPostCreate(...)` |

Hook 的作用分别是：

| 方法 | 执行时机 | 常见用途 |
| --- | --- | --- |
| `WithPreCreate` | 插入前，且仍在同一事务内 | 归一化、校验、补默认值、追加查询控制 |
| `WithPostCreate` | 插入后，且仍在同一事务内 | 同事务内的业务联动、副作用 |

#### Create Hook 示例

```go
crud.NewCreate[User, UserParams]().
	WithPreCreate(func(model *User, params *UserParams, query orm.InsertQuery, ctx fiber.Ctx, tx orm.DB) error {
		// 插入前设置派生字段
		model.Username = strings.ToLower(model.Username)

		// 唯一性校验
		exists, err := tx.NewSelect().Model((*User)(nil)).
			Where(func(cb orm.ConditionBuilder) { cb.Equals("email", model.Email) }).
			Exists(ctx.Context())
		if err != nil {
			return err
		}
		if exists {
			return result.Err("邮箱已存在")
		}
		return nil
	}).
	WithPostCreate(func(model *User, params *UserParams, ctx fiber.Ctx, tx orm.DB) error {
		// 在同一事务中创建关联记录
		role := &UserRole{UserID: model.ID, RoleID: "default"}
		_, err := tx.NewInsert().Model(role).Exec(ctx.Context())
		return err
	})
```

### `Update[TModel, TParams]`

单条更新时使用 `Update`。

| 维度 | 说明 |
| --- | --- |
| 泛型 | `TModel` 是持久化模型，`TParams` 是写入参数类型 |
| 输入 | `params` 中的 `TParams`，且必须包含主键字段 |
| 输出 | 成功结果 |
| 默认行为 | 先把 params 拷贝到临时模型，校验主键，加载旧模型，应用数据权限，合并非空字段，再在事务中更新 |
| 特有配置 | `WithPreUpdate(...)`、`WithPostUpdate(...)`、`DisableDataPerm()` |

一个重要细节：

- `Update` 合并新值时使用的是 `copier.WithIgnoreEmpty()`，也就是空值字段不会覆盖旧值

#### Update Hook 示例

```go
crud.NewUpdate[User, UserParams]().
	WithPreUpdate(func(oldModel, model *User, params *UserParams, query orm.UpdateQuery, ctx fiber.Ctx, tx orm.DB) error {
		// 比较新旧值来执行业务规则
		if oldModel.IsActive && !model.IsActive {
			// 停用：检查是否有待办任务
			count, err := tx.NewSelect().Model((*Task)(nil)).
				Where(func(cb orm.ConditionBuilder) {
					cb.Equals("assignee_id", model.ID).
						Equals("status", "pending")
				}).Count(ctx.Context())
			if err != nil {
				return err
			}
			if count > 0 {
				return result.Err("无法停用：用户有待办任务")
			}
		}
		return nil
	}).
	WithPostUpdate(func(oldModel, model *User, params *UserParams, ctx fiber.Ctx, tx orm.DB) error {
		// 记录变更日志
		return nil
	})
```

### `Delete[TModel]`

单条删除时使用 `Delete`。

| 维度 | 说明 |
| --- | --- |
| 泛型 | `TModel` 是持久化模型 |
| 输入 | `params` 中的原始主键值 |
| 输出 | 成功结果 |
| 默认行为 | 校验主键输入，加载模型，应用数据权限，在事务中删除，并在成功后清理已提升文件 |
| 特有配置 | `WithPreDelete(...)`、`WithPostDelete(...)`、`DisableDataPerm()` |

#### Delete Hook 示例

```go
crud.NewDelete[User]().
	WithPreDelete(func(model *User, query orm.DeleteQuery, ctx fiber.Ctx, tx orm.DB) error {
		// 禁止删除管理员
		if model.Username == "admin" {
			return result.Err("无法删除管理员账户")
		}
		// 级联：删除关联记录
		_, err := tx.NewDelete().Model((*UserRole)(nil)).
			Where(func(cb orm.ConditionBuilder) {
				cb.Equals("user_id", model.ID)
			}).Exec(ctx.Context())
		return err
	})
```

### 批量 Builder

#### `CreateMany[TModel, TParams]`

| 维度 | 说明 |
| --- | --- |
| 输入契约 | `CreateManyParams[TParams]`，字段名是 `list` |
| 输出 | 主键 map 列表 |
| 特有配置 | `WithPreCreateMany(...)`、`WithPostCreateMany(...)` |
| 默认行为 | 把每个 params 项拷贝为模型，并在一个事务里完成批量插入 |

#### `UpdateMany[TModel, TParams]`

| 维度 | 说明 |
| --- | --- |
| 输入契约 | `UpdateManyParams[TParams]`，字段名是 `list` |
| 输出 | 成功结果 |
| 特有配置 | `WithPreUpdateMany(...)`、`WithPostUpdateMany(...)`、`DisableDataPerm()` |
| 默认行为 | 校验每项主键、加载所有旧模型、合并更新值，并在一个事务里执行批量更新 |

#### `DeleteMany[TModel]`

| 维度 | 说明 |
| --- | --- |
| 输入契约 | `DeleteManyParams`，字段名是 `pks` |
| 输出 | 成功结果 |
| 特有配置 | `WithPreDeleteMany(...)`、`WithPostDeleteMany(...)`、`DisableDataPerm()` |
| 默认行为 | 对单主键模型接受标量值，对复合主键模型接受 map，并在一个事务里完成批量删除 |

`DeleteManyParams.pks` 的输入规则：

| 主键形态 | 允许的 payload 形态 |
| --- | --- |
| 单主键 | `["id1", "id2"]` |
| 复合主键 | `[{"user_id":"u1","role_id":"r1"}]` |

`Sortable.Sort` 从 `meta.sort` 解码。`TreeDataOption.ID` 和
`TreeDataOption.ParentID` 是树构建内部字段：它们从 `id` 和 `parent_id`
选出，但因为 JSON tag 是 `json:"-"`，不会作为 JSON 字段输出。

## 导出与导入 Builder

### `Export[TModel, TSearch]`

当你需要把查询结果下载成 Excel 或 CSV 文件时，使用 `Export`。

| 维度 | 说明 |
| --- | --- |
| 泛型 | `TModel` 是导出行模型，`TSearch` 是查询条件类型 |
| 输入 | `params` 中的 `TSearch`，`meta` 中的 `format` |
| 输出 | 文件下载 |
| 默认行为 | 先执行一个 Find 风格查询，再在可选预处理后导出 Excel / CSV |
| 特有配置 | `WithDefaultFormat(...)`、`WithExcelOptions(...)`、`WithCsvOptions(...)`、`WithPreExport(...)`、`WithFilenameBuilder(...)` |

`format` 可选值：

| 格式 | 值 |
| --- | --- |
| Excel | `excel` |
| CSV | `csv` |

默认值：

- 导出格式默认是 `excel`
- 默认文件名分别是 `data.xlsx` 和 `data.csv`

### `Import[TModel]`

当调用方上传 CSV 或 Excel 文件，并希望把行数据解析后插入数据库时，使用 `Import`。

| 维度 | 说明 |
| --- | --- |
| 泛型 | `TModel` 是导入后要持久化的模型类型 |
| 输入 | `multipart` 上传文件 `params.file`，以及 `meta` 中可选的 `format` |
| 输出 | 成功时返回 `{total: n}` |
| 默认行为 | 强制要求 multipart 请求，把文件解析成模型，校验导入结果，并在事务中写入 |
| 特有配置 | `WithDefaultFormat(...)`、`WithExcelOptions(...)`、`WithCsvOptions(...)`、`WithPreImport(...)`、`WithPostImport(...)` |

重要细节：

- Import 不接受 JSON 请求
- 如果导入校验失败，响应会返回 `errors` 载荷，而不是部分写入
- 导入格式默认是 `excel`

## Processor 类型签名

所有 hook / processor 类型都定义在 `crud` 包中：

### 读操作 Processor

```go
// 在响应序列化前转换查询结果
type Processor[TIn, TSearch any] func(input TIn, search TSearch, ctx fiber.Ctx) any
```

### 写操作 Processor（单条）

```go
type PreCreateProcessor[TModel, TParams any]  func(model *TModel, params *TParams, query orm.InsertQuery, ctx fiber.Ctx, tx orm.DB) error
type PostCreateProcessor[TModel, TParams any] func(model *TModel, params *TParams, ctx fiber.Ctx, tx orm.DB) error

type PreUpdateProcessor[TModel, TParams any]  func(oldModel, model *TModel, params *TParams, query orm.UpdateQuery, ctx fiber.Ctx, tx orm.DB) error
type PostUpdateProcessor[TModel, TParams any] func(oldModel, model *TModel, params *TParams, ctx fiber.Ctx, tx orm.DB) error

type PreDeleteProcessor[TModel any]  func(model *TModel, query orm.DeleteQuery, ctx fiber.Ctx, tx orm.DB) error
type PostDeleteProcessor[TModel any] func(model *TModel, ctx fiber.Ctx, tx orm.DB) error
```

### 写操作 Processor（批量）

```go
type PreCreateManyProcessor[TModel, TParams any]  func(models []TModel, paramsList []TParams, query orm.InsertQuery, ctx fiber.Ctx, tx orm.DB) error
type PostCreateManyProcessor[TModel, TParams any] func(models []TModel, paramsList []TParams, ctx fiber.Ctx, tx orm.DB) error

type PreUpdateManyProcessor[TModel, TParams any]  func(oldModels, models []TModel, paramsList []TParams, query orm.UpdateQuery, ctx fiber.Ctx, tx orm.DB) error
type PostUpdateManyProcessor[TModel, TParams any] func(oldModels, models []TModel, paramsList []TParams, ctx fiber.Ctx, tx orm.DB) error

type PreDeleteManyProcessor[TModel any]  func(models []TModel, query orm.DeleteQuery, ctx fiber.Ctx, tx orm.DB) error
type PostDeleteManyProcessor[TModel any] func(models []TModel, ctx fiber.Ctx, tx orm.DB) error
```

### 导出 / 导入 Processor

```go
type PreExportProcessor[TModel, TSearch any] func(models []TModel, search TSearch, ctx fiber.Ctx, db orm.DB) error
type FilenameBuilder[TSearch any] func(search TSearch, ctx fiber.Ctx) string

type PreImportProcessor[TModel any]  func(models []TModel, query orm.InsertQuery, ctx fiber.Ctx, tx orm.DB) error
type PostImportProcessor[TModel any] func(models []TModel, ctx fiber.Ctx, tx orm.DB) error
```

## Supporting 公开 API

| API 组 | 公开 surface |
| --- | --- |
| 直接构造器 | `NewBuilder`、`NewFind`，以及上文列出的 `NewCreate` / `NewUpdate` / `NewDelete` / read / import / export 构造器 |
| action 常量 | RPC action：`RPCActionCreate`, `RPCActionUpdate`, `RPCActionDelete`, `RPCActionCreateMany`, `RPCActionUpdateMany`, `RPCActionDeleteMany`, `RPCActionFindOne`, `RPCActionFindAll`, `RPCActionFindPage`, `RPCActionFindOptions`, `RPCActionFindTree`, `RPCActionFindTreeOptions`, `RPCActionImport`, `RPCActionExport`；REST action：`RESTActionCreate`, `RESTActionUpdate`, `RESTActionDelete`, `RESTActionCreateMany`, `RESTActionUpdateMany`, `RESTActionDeleteMany`, `RESTActionFindOne`, `RESTActionFindAll`, `RESTActionFindPage`, `RESTActionFindOptions`, `RESTActionFindTree`, `RESTActionFindTreeOptions`, `RESTActionImport`, `RESTActionExport` |
| 格式 | `FormatExcel`, `FormatCsv`, `TabularFormat` |
| option payload | `CreateManyParams`, `UpdateManyParams`, `DeleteManyParams`, `DataOption`, `TreeDataOption`, `DataOptionConfig`, `DataOptionColumnMapping` |
| 查询塑形 | `FindOperationConfig`, `FindOperationOption`, `QueryPartsConfig`, `QueryPart`, `Sortable`, `ApplyDataPermission` |
| 审计/树 helper | `GetAuditUserNameRelations`，以及 `IDColumn`, `ParentIDColumn`, `LabelColumn`, `ValueColumn`, `DescriptionColumn` |
| import/export hook | `FilenameBuilder`, `PreExportProcessor`, `PreImportProcessor`, `PostImportProcessor` |
| 错误 helper | `ErrPrimaryKeyRequired(...)`、`ErrFieldNotExistInModel(...)`，以及 `ErrModelNoPrimaryKey`、`ErrCompositePrimaryKeyRequiresMap`、`ErrUnsupportedExportFormat`、`ErrUnsupportedImportFormat`、`ErrImportRequiresFile`、`ErrImportRequiresMultipart`、`ErrFileOpenFailed`、`ErrImportTypeAssertionFailed`、`ErrAuditUserCompositePK` 等 sentinel |
| 错误码 | `ErrCodeProcessorInvalidReturn`, `ErrCodeFieldNotExistInModel`, `ErrCodePrimaryKeyRequired`, `ErrCodeCompositePrimaryKeyRequiresMap`, `ErrCodeUnsupportedExportFormat`, `ErrCodeImportRequiresMultipart`, `ErrCodeImportRequiresFile`, `ErrCodeUnsupportedImportFormat`, `ErrCodeFileOpenFailed`, `ErrCodeImportTypeAssertionFailed`, `ErrCodeImportValidationFailed` |

CRUD 错误码值如下：

| 常量 | 值 |
| --- | --- |
| `ErrCodeProcessorInvalidReturn` | `2400` |
| `ErrCodeFieldNotExistInModel` | `2401` |
| `ErrCodePrimaryKeyRequired` | `2402` |
| `ErrCodeCompositePrimaryKeyRequiresMap` | `2403` |
| `ErrCodeUnsupportedExportFormat` | `2404` |
| `ErrCodeImportRequiresMultipart` | `2405` |
| `ErrCodeImportRequiresFile` | `2406` |
| `ErrCodeUnsupportedImportFormat` | `2407` |
| `ErrCodeFileOpenFailed` | `2408` |
| `ErrCodeImportTypeAssertionFailed` | `2409` |
| `ErrCodeImportValidationFailed` | `2410` |

## 实践建议

- 后台资源优先从 `FindPage + Create + Update + Delete` 开始
- 写 params 和 search params 必须分开
- 权限控制尽量加在 builder 层，而不是散落到 handler 内部
- 除非你非常清楚为什么要关闭，否则优先保留默认数据权限过滤
- 下拉框或轻量树选项优先使用 `FindOptions` / `FindTreeOptions`，不要复用完整模型查询
- 除非业务动作有更强的领域语义，否则优先沿用标准 CRUD 词汇
- 一个资源里混合 CRUD builder 和少量自定义 action，是很常见也很合理的做法

## 下一步

继续看 [自定义 Handler](./custom-handlers)，当你的业务动作超出通用 CRUD 模型时，就该切到那里了。
