---
sidebar_position: 1
---

# 模型

在 VEF 中，模型本质上还是普通 Go 结构体，但它们通常会同时配合 Bun、验证标签、搜索标签以及框架内置的审计约定一起使用。

## 常见模型写法

大多数持久化模型会长这样：

```go
type User struct {
	bun.BaseModel `bun:"table:sys_user,alias:su"`
	orm.FullAuditedModel

	Username string `json:"username" validate:"required,alphanum,max=32" label:"Username"`
	Email    string `json:"email" validate:"omitempty,email,max=128" label:"Email"`
	IsActive bool   `json:"isActive"`
}
```

这里同时组合了两类职责：

- `bun.BaseModel`：Bun 的表元信息
- `orm.FullAuditedModel`：框架标准化的 ID、创建审计和更新审计字段

## 基础模型类型

VEF 通过 `orm` 暴露了 **五种** 可复用的模型类型。它们是为匿名嵌入设计的可组合字段片段。

### 不含主键的类型

| 类型 | 字段 | 使用场景 |
| --- | --- | --- |
| `orm.CreationTrackedModel` | `CreatedAt`、`CreatedBy`、`CreatedByName` | 需要创建追踪的复合主键表 |
| `orm.FullTrackedModel` | `CreatedAt`、`CreatedBy`、`CreatedByName`、`UpdatedAt`、`UpdatedBy`、`UpdatedByName` | 需要完整审计追踪的复合主键表 |

### 含主键的类型

| 类型 | 字段 | 使用场景 |
| --- | --- | --- |
| `orm.Model` | 仅 `ID` | 字典表、关联表、最简记录 |
| `orm.CreationAuditedModel` | `ID`、`CreatedAt`、`CreatedBy`、`CreatedByName` | 追加写入型记录、日志、outbox 表 |
| `orm.FullAuditedModel` | `ID`、`CreatedAt`、`CreatedBy`、`CreatedByName`、`UpdatedAt`、`UpdatedBy`、`UpdatedByName` | 标准可变实体，需要完整审计追踪 |

### 如何选择

按实体的生命周期选择最小但够用的类型：

- **`orm.Model`**：只需要主键，完全不需要审计追踪
- **`orm.CreationAuditedModel`**：追加写入型记录——写入就不再更新
- **`orm.FullAuditedModel`**：最常见的选择——标准可变实体，需要同时追踪创建和更新信息
- **`orm.CreationTrackedModel`**：和 `CreationAuditedModel` 一样但不含主键——适合复合主键表
- **`orm.FullTrackedModel`**：和 `FullAuditedModel` 一样但不含主键——适合复合主键表

从语义上说，`orm.FullAuditedModel` 可以看成是 `orm.Model` + `orm.FullTrackedModel` 的预组合版本。

### 内部字段定义

以下是每个类型的完整字段定义，包含所有结构体标签：

```go
// orm.Model — 仅主键
type Model struct {
	ID string `json:"id" bun:"id,pk"`
}

// orm.CreationTrackedModel — 创建审计（无主键）
type CreationTrackedModel struct {
	CreatedAt     timex.DateTime `json:"createdAt" bun:",notnull,type:timestamp,default:CURRENT_TIMESTAMP,skipupdate"`
	CreatedBy     string         `json:"createdBy" bun:",notnull,skipupdate" mold:"translate=user?"`
	CreatedByName string         `json:"createdByName" bun:",scanonly"`
}

// orm.FullTrackedModel — 完整审计（无主键）
type FullTrackedModel struct {
	CreatedAt     timex.DateTime `json:"createdAt" bun:",notnull,type:timestamp,default:CURRENT_TIMESTAMP,skipupdate"`
	CreatedBy     string         `json:"createdBy" bun:",notnull,skipupdate" mold:"translate=user?"`
	CreatedByName string         `json:"createdByName" bun:",scanonly"`
	UpdatedAt     timex.DateTime `json:"updatedAt" bun:",notnull,type:timestamp,default:CURRENT_TIMESTAMP"`
	UpdatedBy     string         `json:"updatedBy" bun:",notnull" mold:"translate=user?"`
	UpdatedByName string         `json:"updatedByName" bun:",scanonly"`
}

// orm.CreationAuditedModel — 主键 + 创建审计
type CreationAuditedModel struct {
	ID            string         `json:"id" bun:"id,pk"`
	CreatedAt     timex.DateTime `json:"createdAt" bun:",notnull,type:timestamp,default:CURRENT_TIMESTAMP,skipupdate"`
	CreatedBy     string         `json:"createdBy" bun:",notnull,skipupdate" mold:"translate=user?"`
	CreatedByName string         `json:"createdByName" bun:",scanonly"`
}

// orm.FullAuditedModel — 主键 + 完整审计
type FullAuditedModel struct {
	ID            string         `json:"id" bun:"id,pk"`
	CreatedAt     timex.DateTime `json:"createdAt" bun:",notnull,type:timestamp,default:CURRENT_TIMESTAMP,skipupdate"`
	CreatedBy     string         `json:"createdBy" bun:",notnull,skipupdate" mold:"translate=user?"`
	CreatedByName string         `json:"createdByName" bun:",scanonly"`
	UpdatedAt     timex.DateTime `json:"updatedAt" bun:",notnull,type:timestamp,default:CURRENT_TIMESTAMP"`
	UpdatedBy     string         `json:"updatedBy" bun:",notnull" mold:"translate=user?"`
	UpdatedByName string         `json:"updatedByName" bun:",scanonly"`
}
```

### 关键标签说明

- **`bun:",skipupdate"`**：`created_at` 和 `created_by` 仅在插入时设置，更新操作不会覆盖
- **`bun:",scanonly"`**：`created_by_name` 和 `updated_by_name` 是读侧辅助字段——它们通过 JOIN 查询填充，不作为独立列持久化
- **`mold:"translate=user?"`**：`mold` 变换器通过数据字典将用户 ID 翻译为显示名称，自动填充 `*ByName` 字段
- **`timex.DateTime`**：框架自定义的时间戳类型（参见 [Timex](../utilities/timex)），而不是标准库的 `time.Time`

## 匿名嵌入与组合

当实体不需要完整的 `orm.FullAuditedModel` 字段集合时，更小的基础模型片段就很有用。

```go
// 最简：仅主键
type Tag struct {
	bun.BaseModel `bun:"table:tag,alias:t"`
	orm.Model

	Name string `json:"name" bun:"name,notnull"`
}

// 追加写入：主键 + 创建审计
type ActionLog struct {
	bun.BaseModel `bun:"table:apv_action_log,alias:aal"`
	orm.Model
	orm.CreationTrackedModel

	InstanceID string `json:"instanceId" bun:"instance_id"`
	Action     string `json:"action" bun:"action"`
}

// 标准可变实体：主键 + 完整审计
type Role struct {
	bun.BaseModel `bun:"table:sys_role,alias:sr"`
	orm.FullAuditedModel

	Name     string `json:"name" bun:"name,notnull"`
	IsActive bool   `json:"isActive" bun:"is_active"`
}

// 复合主键：审计字段但主键需要单独定义
type UserRole struct {
	bun.BaseModel `bun:"table:sys_user_role,alias:sur"`
	orm.Model
	orm.CreationTrackedModel

	UserID string `json:"userId" bun:"user_id,notnull"`
	RoleID string `json:"roleId" bun:"role_id,notnull"`
}
```

常见选择：

- `orm.Model`：字典表、关联表，或根本不需要审计列的记录
- `orm.Model` + `orm.CreationTrackedModel`：追加写入型记录、快照、日志、outbox 表
- `orm.FullAuditedModel`：标准可变实体，需要同时追踪创建和更新信息
- `orm.FullTrackedModel`：复合主键实体，但仍需要完整审计追踪

## Bun 模型钩子

VEF 通过 `orm` 重导出了 Bun 的模型生命周期钩子接口：

| 钩子接口 | 调用时机 |
| --- | --- |
| `orm.BeforeSelectHook` | SELECT 查询执行前 |
| `orm.AfterSelectHook` | SELECT 查询执行后 |
| `orm.BeforeInsertHook` | INSERT 查询执行前 |
| `orm.AfterInsertHook` | INSERT 查询执行后 |
| `orm.BeforeUpdateHook` | UPDATE 查询执行前 |
| `orm.AfterUpdateHook` | UPDATE 查询执行后 |
| `orm.BeforeDeleteHook` | DELETE 查询执行前 |
| `orm.AfterDeleteHook` | DELETE 查询执行后 |
| `orm.BeforeScanRowHook` | 扫描行之前 |
| `orm.AfterScanRowHook` | 扫描行之后 |

在模型结构体上实现这些接口即可加入生命周期行为：

```go
func (u *User) BeforeInsert(ctx context.Context, query *orm.BunInsertQuery) error {
	// 在插入前设置默认值、验证或记录日志
	return nil
}
```

## 审计字段

框架统一约定了这些常见审计列：

| 列名 | JSON 名 | 是否持久化 | 说明 |
| --- | --- | --- | --- |
| `id` | `id` | ✅ | 主键 |
| `created_at` | `createdAt` | ✅ | 创建时间 |
| `created_by` | `createdBy` | ✅ | 创建者用户 ID |
| `created_by_name` | `createdByName` | ❌ scanonly | 创建者显示名（通过 mold 或 JOIN 填充）|
| `updated_at` | `updatedAt` | ✅ | 最后更新时间 |
| `updated_by` | `updatedBy` | ✅ | 最后更新者用户 ID |
| `updated_by_name` | `updatedByName` | ❌ scanonly | 更新者显示名（通过 mold 或 JOIN 填充）|

并不是每个模型都会带上全部这些字段。`orm.CreationTrackedModel` 只提供 `created_*` 这一组，`orm.FullTrackedModel` 和 `orm.FullAuditedModel` 才会同时提供 `created_*` 与 `updated_*` 两组。

框架还导出了审计列名和字段名常量：

```go
orm.ColumnID            // "id"
orm.ColumnCreatedAt     // "created_at"
orm.ColumnUpdatedAt     // "updated_at"
orm.ColumnCreatedBy     // "created_by"
orm.ColumnUpdatedBy     // "updated_by"
orm.ColumnCreatedByName // "created_by_name"
orm.ColumnUpdatedByName // "updated_by_name"

orm.FieldID             // "ID"
orm.FieldCreatedAt      // "CreatedAt"
// ... 以此类推
```

系统操作者常量（用于 `created_by` / `updated_by`）：

```go
orm.OperatorSystem    // "system" — 系统初始化使用
orm.OperatorCronJob   // "cron_job" — 定时任务使用
orm.OperatorAnonymous // "anonymous" — 未认证操作使用
```

## 最常见的标签

### `bun`

控制表名、别名、主键、空值行为和关联关系。

### `json`

控制请求和响应中的字段名。实际项目里通常会保持 JSON 使用 camelCase，而数据库列使用 snake_case。

### `validate`

用于请求参数自动验证。

### `label` / `label_i18n`

用于在验证错误里显示更友好的字段名。

### `search`

用于搜索解析器和 CRUD 查询 builder，把搜索字段翻译成 SQL 条件。

### `meta`

用于存储 promoter，识别上传文件字段、富文本字段和 Markdown 字段。

### `mold`

用于结构体变换器进行字段级数据转换。最常见的内置用法是 `*ByName` 字段上的 `mold:"translate=user?"`。

## 搜索模型通常单独定义

不要把搜索语义硬塞进持久化模型里。更推荐单独写搜索结构体：

```go
type UserSearch struct {
	api.P

	Keyword  string `json:"keyword" search:"contains,column=username|email"`
	IsActive *bool  `json:"isActive" search:"eq,column=is_active"`
}
```

这样查询规则就能保持清晰，也不会让写模型变成半个 DSL。

## 分页和排序元信息

分页接口通常还会配一个 metadata 结构体，例如：

```go
type UserSearch struct {
	api.P
	Keyword string `json:"keyword" search:"contains,column=username|email"`
}

type UserMeta struct {
	api.M
	page.Pageable
	crud.Sortable
}
```

这里的 `page.Pageable` 和 `crud.Sortable` 都是元信息辅助类型，不是持久化模型本身。

## 实践建议

- 持久化模型尽量小而明确
- 写操作和读操作分别使用独立 params / search 结构体
- 当实体不需要完整字段集时，优先组合更小的基础嵌入模型
- 只有确实想接入框架标准审计行为时，再嵌入 `orm.FullAuditedModel`
- 数据库标签、验证标签和搜索标签尽量跟字段放在一起，保持规则可见
- 记住 `*ByName` 字段是 scanonly 的——它们永远不会被写入数据库

## 下一步

继续阅读 [泛型 CRUD](./crud) 了解这些模型如何接入类型化操作 builder，或阅读 [ORM SQL 构造器](./orm-builder) 获取 SQL 查询构造的完整参考。
