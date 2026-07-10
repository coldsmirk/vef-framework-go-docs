---
sidebar_position: 3
title: 你的第一个 CRUD API
---

# 你的第一个 CRUD API

[快速开始](./quick-start.md)只提供了一个手写的处理函数。本教程会构建你接下来真正需要的东西：围绕 `Product` 实体的一套完整 CRUD API——模型、数据表、类型化请求参数、通用 CRUD 操作，以及一个自定义钩子——全程都可以用 curl 验证。

## 你将构建什么

- 一个 `app/product` RPC 资源，提供 `create`、`update`、`delete` 和 `find_page`
- 一张遵循[数据库规范](../conventions/database-conventions.md)的 `app_product` 表
- 由 `search` 标签驱动的关键字与状态过滤，以及分页能力
- 一个拒绝重复商品编码的 pre-create 钩子

## 前置条件

- 已完成[快速开始](./quick-start.md)并拥有可运行的应用
- 安装了 `sqlite3` 命令行工具（本教程沿用快速开始中的 SQLite 配置）

完成后的目录结构如下：

```text
my-app/
├── configs/
│   └── application.toml
├── data/
│   └── app.db
├── db/
│   └── app_product.sql
├── internal/
│   └── product/
│       ├── model.go
│       ├── payload.go
│       ├── resource.go
│       └── module.go
└── main.go
```

## 1. 定义模型

创建 `internal/product/model.go`：

```go
package product

import (
	"github.com/uptrace/bun"

	"github.com/coldsmirk/vef-framework-go/orm"
)

type Product struct {
	bun.BaseModel `bun:"table:app_product,alias:ap"`
	orm.FullAuditedModel

	Name     string `json:"name" bun:"name,notnull"`
	Code     string `json:"code" bun:"code,notnull"`
	Stock    int    `json:"stock" bun:"stock,notnull"`
	IsActive bool   `json:"isActive" bun:"is_active,notnull"`
	Remark   string `json:"remark" bun:"remark"`
}
```

两个嵌入类型承担了主要工作：

- `bun.BaseModel` 把结构体绑定到 `app_product` 表，别名为 `ap`
- `orm.FullAuditedModel` 提供 `ID`、`CreatedAt`、`CreatedBy`、`CreatedByName`、`UpdatedAt`、`UpdatedBy` 和 `UpdatedByName`

这些框架托管的字段从不需要你手动赋值。插入时，框架会为空字符串主键生成一个紧凑的字符串 ID，并填充 `created_at` / `created_by`；更新时会依据当前操作者维护 `updated_at` / `updated_by`。基础模型的完整目录见[模型](../data-access/models.md)。

## 2. 创建数据表

VEF 不会根据模型生成表结构。应用项目自行维护 DDL 脚本，并遵循[数据库规范](../conventions/database-conventions.md)：表名带模块前缀（这里是 `app_`）、固定的审计列，以及显式命名的约束。

创建 `db/app_product.sql`：

```sql
BEGIN;

CREATE TABLE IF NOT EXISTS app_product (
    id                       VARCHAR(32) NOT NULL,
    created_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by               VARCHAR(32) NOT NULL DEFAULT 'system',
    updated_by               VARCHAR(32) NOT NULL DEFAULT 'system',
    name                     VARCHAR(32) NOT NULL,
    code                     VARCHAR(32) NOT NULL,
    stock                    INTEGER NOT NULL DEFAULT 0,
    is_active                BOOLEAN NOT NULL DEFAULT FALSE,
    remark                   VARCHAR(512),

    CONSTRAINT pk_app_product PRIMARY KEY (id),
    CONSTRAINT uk_app_product__code UNIQUE (code)
);

COMMIT;
```

把它应用到应用将要使用的 SQLite 数据库文件：

```bash
mkdir -p data
sqlite3 data/app.db < db/app_product.sql
```

这份脚本刻意只保留了可移植的子集。在真实的 PostgreSQL 项目中，规范还要求 `LOCALTIMESTAMP` 默认值、为每张表和每个列编写 `COMMENT ON` 语句，以及让 `created_by` / `updated_by` 外键引用 `sys_user(id)`——完整模板见[数据库规范](../conventions/database-conventions.md)。

## 3. 定义写入参数与查询参数

持久化模型不应兼任请求契约。创建 `internal/product/payload.go`，分别定义一个写入结构体和一个查询结构体：

```go
package product

import (
	"github.com/coldsmirk/vef-framework-go/api"
)

type ProductParams struct {
	api.P

	ID       string `json:"id"`
	Name     string `json:"name" validate:"required,max=32" label:"Name"`
	Code     string `json:"code" validate:"required,max=32" label:"Code"`
	Stock    int    `json:"stock" validate:"gte=0" label:"Stock"`
	IsActive *bool  `json:"isActive"`
	Remark   string `json:"remark" validate:"max=512" label:"Remark"`
}

type ProductSearch struct {
	api.P

	Keyword  *string `json:"keyword" search:"contains,column=name|code"`
	IsActive *bool   `json:"isActive" search:"eq"`
	MinStock *int    `json:"minStock" search:"gte,column=stock"`
}
```

各部分的作用：

- 嵌入的 `api.P` 哨兵类型告诉框架从请求的 `params` 字段解码该结构体并进行校验
- `validate` 标签会在你的操作执行前自动运行；`label` 决定错误消息中的字段名称
- `search` 标签会直接翻译成 `WHERE` 条件：`keyword` 变成对 `name` 或 `code` 的 `LIKE` 匹配，`minStock` 变成 `stock >= ?`
- `ID` 在创建时留空（框架会生成），在更新时必填
- 指针字段用于区分「未提供」和零值——更新只合并非空字段，因此 `IsActive *bool` 才能让客户端显式传 `false`

## 4. 组装 API 资源

创建 `internal/product/resource.go`：

```go
package product

import (
	"github.com/coldsmirk/vef-framework-go/api"
	"github.com/coldsmirk/vef-framework-go/crud"
)

type ProductResource struct {
	api.Resource

	crud.FindPage[Product, ProductSearch]
	crud.Create[Product, ProductParams]
	crud.Update[Product, ProductParams]
	crud.Delete[Product]
}

func NewProductResource() api.Resource {
	return &ProductResource{
		Resource: api.NewRPCResource("app/product"),
		FindPage: crud.NewFindPage[Product, ProductSearch]().Public(),
		Create:   crud.NewCreate[Product, ProductParams]().Public(),
		Update:   crud.NewUpdate[Product, ProductParams]().Public(),
		Delete:   crud.NewDelete[Product]().Public(),
	}
}
```

每个嵌入的构建器都实现了 `api.OperationsProvider`，框架会自动收集它们，并为每个构建器注册一个操作：

| 嵌入的构建器 | 默认 action | 行为 |
| --- | --- | --- |
| `crud.FindPage[Product, ProductSearch]` | `find_page` | 带过滤和总数统计的分页列表 |
| `crud.Create[Product, ProductParams]` | `create` | 把参数拷贝进模型，在事务中插入 |
| `crud.Update[Product, ProductParams]` | `update` | 按 `id` 加载记录，合并非空字段后更新 |
| `crud.Delete[Product]` | `delete` | 按 `id` 加载记录，在事务中删除 |

`Public()` 让本教程无需认证提供者即可运行。在真实应用中应去掉它，改为按操作声明权限：

```go
crud.NewCreate[Product, ProductParams]().RequiredPermission("app:product:create")
```

## 5. 注册模块并运行

创建 `internal/product/module.go`：

```go
package product

import (
	"github.com/coldsmirk/vef-framework-go"
)

var Module = vef.Module(
	"app:product",
	vef.ProvideAPIResource(NewProductResource),
)
```

在 `main.go` 中组合它：

```go
package main

import (
	"github.com/coldsmirk/vef-framework-go"

	"example.com/my-app/internal/product"
)

func main() {
	vef.Run(
		product.Module,
	)
}
```

让 `configs/application.toml` 指向第 2 步创建的数据库文件：

```toml
[vef.app]
name = "my-app"
port = 8080

[vef.data_sources.primary]
type = "sqlite"
path = "data/app.db"

[vef.event.transports.outbox]
enabled = true

[[vef.event.routing]]
pattern = "vef.storage.*"
transports = ["outbox"]
```

相比快速开始，最后两个配置块是新增的。通用写操作会在事务内执行[文件存储](../infrastructure/storage.md)生命周期，而存储模块通过事务性传输发布领域事件——如果 `vef.storage.*` 事件没有这样的路由，框架会在启动时快速失败。启用 outbox 传输（它会自动创建自己的表）并把存储事件路由过去即可通过该检查。

启动应用：

```bash
go run .
```

## 6. 调用 API

四个操作都通过同一个 RPC 端点 `POST /api` 完成，由信封字段 `resource` / `action` / `version` 选择目标。

### 创建

```bash
curl http://localhost:8080/api \
  -H 'Content-Type: application/json' \
  -d '{
    "resource": "app/product",
    "action": "create",
    "version": "v1",
    "params": {
      "name": "Espresso Beans",
      "code": "P-1001",
      "stock": 20,
      "isActive": true
    }
  }'
```

响应会返回生成的主键：

```json
{
  "code": 0,
  "message": "新增成功",
  "data": {
    "id": "d1nbkq2s7kg5jkvvs7lg"
  }
}
```

与快速开始一样，消息文本跟随框架的默认语言；设置 `VEF_I18N_LANGUAGE=en` 后会得到 `Created successfully`。

### 分页查询

过滤条件来自 `params`（即你的 `ProductSearch`），分页参数来自 `meta`：

```bash
curl http://localhost:8080/api \
  -H 'Content-Type: application/json' \
  -d '{
    "resource": "app/product",
    "action": "find_page",
    "version": "v1",
    "params": { "keyword": "Espresso" },
    "meta": { "page": 1, "size": 10 }
  }'
```

`data` 载荷是一个分页对象：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "page": 1,
    "size": 10,
    "total": 1,
    "items": [
      {
        "id": "d1nbkq2s7kg5jkvvs7lg",
        "createdAt": "2026-07-09 10:30:00",
        "createdBy": "anonymous",
        "createdByName": "",
        "updatedAt": "2026-07-09 10:30:00",
        "updatedBy": "anonymous",
        "updatedByName": "",
        "name": "Espresso Beans",
        "code": "P-1001",
        "stock": 20,
        "isActive": true,
        "remark": ""
      }
    ]
  }
}
```

注意 `createdBy: "anonymous"`：审计列是自动填充的，而由于该操作是公开的，此时还没有已认证的操作者。

### 更新

传入 `id` 和要修改的字段；未传的字段保持数据库中的原值：

```bash
curl http://localhost:8080/api \
  -H 'Content-Type: application/json' \
  -d '{
    "resource": "app/product",
    "action": "update",
    "version": "v1",
    "params": {
      "id": "d1nbkq2s7kg5jkvvs7lg",
      "name": "Espresso Beans",
      "code": "P-1001",
      "stock": 35
    }
  }'
```

```json
{ "code": 0, "message": "保存成功", "data": null }
```

### 删除

```bash
curl http://localhost:8080/api \
  -H 'Content-Type: application/json' \
  -d '{
    "resource": "app/product",
    "action": "delete",
    "version": "v1",
    "params": { "id": "d1nbkq2s7kg5jkvvs7lg" }
  }'
```

```json
{ "code": 0, "message": "删除成功", "data": null }
```

信封字段与传输规则的完整说明见[路由](../building-apis/routing.md)，`code` / `message` / `data` 契约见[结果与错误](../building-apis/results-and-errors.md)。

## 7. 添加一个创建钩子

通用构建器支持在写入的同一事务内运行钩子。用 `WithPreCreate` 在 API 层校验商品编码唯一性，让调用方得到结构化的业务错误，而不是底层的约束冲突。

更新 `internal/product/resource.go` 中的 `Create` 构建器：

```go
Create: crud.NewCreate[Product, ProductParams]().
	Public().
	WithPreCreate(func(model *Product, params *ProductParams, query orm.InsertQuery, ctx fiber.Ctx, tx orm.DB) error {
		exists, err := tx.NewSelect().Model((*Product)(nil)).
			Where(func(cb orm.ConditionBuilder) { cb.Equals("code", model.Code) }).
			Exists(ctx.Context())
		if err != nil {
			return err
		}
		if exists {
			return result.Err("product code already exists",
				result.WithCode(result.ErrCodeRecordAlreadyExists))
		}
		return nil
	}),
```

同时在该文件中补充新的导入：`github.com/gofiber/fiber/v3`、`github.com/coldsmirk/vef-framework-go/orm` 和 `github.com/coldsmirk/vef-framework-go/result`。

重启应用，把第 6 步的创建请求重放两次。第二次调用会得到干净的失败响应：

```json
{
  "code": 2002,
  "message": "product code already exists",
  "data": null
}
```

钩子在插入之前、事务内部执行，可以访问待写入的模型、解码后的参数、插入查询以及事务版 `orm.DB`。每个构建器都有对应的一组钩子——`WithPreUpdate` / `WithPostUpdate`、`WithPreDelete` / `WithPostDelete` 等——完整目录见[通用 CRUD](../data-access/crud.md)。

## 框架为你做了什么

你只写了一个模型、两个请求结构体和一个资源。你没有为这四个操作写 SQL，也没有写请求解码、校验接线、事务、ID 生成、审计列维护、分页计数或响应信封。

## 下一步

- [通用 CRUD](../data-access/crud.md)：`crud` 包中的全部构建器、选项与钩子
- [API 资源](../building-apis/api.md)：操作、认证配置，以及 CRUD 之外的自定义 action
- [模型](../data-access/models.md)：基础模型类型、标签与查询结构体模式
