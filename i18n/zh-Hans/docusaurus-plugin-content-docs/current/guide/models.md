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
	orm.Model

	Username string `json:"username" validate:"required,alphanum,max=32" label:"Username"`
	Email    string `json:"email" validate:"omitempty,email,max=128" label:"Email"`
	IsActive bool   `json:"isActive"`
}
```

这里同时组合了两类职责：

- `bun.BaseModel`：Bun 的表元信息
- `orm.Model`：框架标准化的 ID 和审计字段

## 常用基础模型类型

VEF 通过 `orm` 暴露了几种常见基础模型：

- `orm.Model`：ID + 创建/更新审计字段
- `orm.IDModel`：只有 ID
- `orm.CreatedModel`：只有创建相关字段
- `orm.AuditedModel`：只有创建/更新审计字段，不包含主键

按你的实体生命周期选择最小但够用的那个。

这些类型本身就是给匿名嵌入准备的可复用字段片段，不只是几个互斥的“基础父类”。`orm.Model` 是最常见场景下的便利组合，而更小的类型则让你只拼装实体真正需要的字段。

可以这样理解它们：

- `orm.IDModel`：只补一个主键字段
- `orm.CreatedModel`：只补 `created_*` 这一组创建审计字段
- `orm.AuditedModel`：补 `created_*` 和 `updated_*` 两组审计字段，但不带主键
- `orm.Model`：补主键，再加上与 `orm.AuditedModel` 相同的审计字段

从语义上说，`orm.Model` 可以看成是 `orm.IDModel + orm.AuditedModel` 的预组合版本。

## 匿名嵌入与组合

当实体不需要完整的 `orm.Model` 字段集合时，更小的基础模型片段就很有用。

```go
type Tag struct {
	bun.BaseModel `bun:"table:tag,alias:t"`
	orm.IDModel

	Name string `json:"name" bun:"name,notnull"`
}

type EventOutbox struct {
	bun.BaseModel `bun:"table:event_outbox,alias:eo"`
	orm.IDModel
	orm.CreatedModel

	EventType string `json:"eventType" bun:"event_type,notnull"`
}

type Delegation struct {
	bun.BaseModel `bun:"table:delegation,alias:d"`
	orm.Model

	DelegatorID string `json:"delegatorId" bun:"delegator_id,notnull"`
}
```

常见选择大致可以这样分：

- `orm.IDModel`：字典表、关联表，或者根本不需要审计列的记录
- `orm.IDModel` + `orm.CreatedModel`：追加写入型记录、快照、日志、outbox 之类的表
- `orm.Model`：标准可变实体，需要同时跟踪创建和更新信息
- `orm.AuditedModel`：主键字段要自己定义，但仍想复用框架标准审计列的实体

## 审计字段

框架统一约定了这些常见审计列：

- `id`
- `created_at`
- `created_by`
- `created_by_name`
- `updated_at`
- `updated_by`
- `updated_by_name`

并不是每个模型都会带上全部这些字段。`orm.CreatedModel` 只提供 `created_*` 这一组，`orm.AuditedModel` 和 `orm.Model` 才会同时提供 `created_*` 与 `updated_*` 两组。

当你嵌入 `orm.Model` 时，VEF 的上下文与 CRUD 行为会围绕这些字段协同工作。

这里有个容易忽略的细节：`created_by_name` 和 `updated_by_name` 在当前模型定义中是 `scanonly` 字段，更适合视为查询结果辅助字段，而不是框架会直接持久化写入的数据库列。

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
- 只有确实想接入框架标准审计行为时，再嵌入 `orm.Model`
- 数据库标签、验证标签和搜索标签尽量跟字段放在一起，保持规则可见

## 下一步

继续阅读 [泛型 CRUD](./crud)，看这些模型如何接入类型化操作 builder。
