---
sidebar_position: 3
---

# 自定义处理器

CRUD builder 能覆盖大量标准接口，但真实项目一定会出现通用 CRUD 模型之外的业务动作。VEF 支持为 RPC 和 REST 资源定义自定义处理器，并且可以直接把请求期依赖或启动期依赖注入到 handler 签名里。

## Handler 解析总览

| 资源类型 | Handler 来源 | 解析规则 |
| --- | --- | --- |
| RPC | 省略 `Handler` | 框架把 action name 从 `snake_case` 转成资源上的 PascalCase 方法名 |
| RPC | 显式 `Handler: "MethodName"` | 框架按方法名在资源上查找 |
| RPC | 显式函数值 | 直接使用该函数 |
| REST | 必须显式指定 `Handler` | REST 不会从 action 字符串推导方法名 |

## RPC Handler 解析

对于 RPC 资源，可以省略显式 `Handler`，让框架根据 action name 自动查找方法。

示例：

```go
type UserResource struct {
	api.Resource
}

func NewUserResource() api.Resource {
	return &UserResource{
		Resource: api.NewRPCResource(
			"sys/user",
			api.WithOperations(
				api.OperationSpec{Action: "ping", Public: true},
			),
		),
	}
}

func (*UserResource) Ping(ctx fiber.Ctx) error {
	return result.Ok("pong").Response(ctx)
}
```

RPC 解析规则：

| Action 名 | 解析出的处理器方法 |
| --- | --- |
| `ping` | `Ping` |
| `find_page` | `FindPage` |
| `get_user_info` | `GetUserInfo` |
| `resolve_challenge` | `ResolveChallenge` |

要求：

- RPC action name 必须使用 `snake_case`
- 自动推导出的 handler 方法必须使用 PascalCase
- 如果你不想走自动推导，可以在 `api.OperationSpec` 里显式设置 `Handler`

## REST Handler 解析

对于 REST 资源，必须显式提供 handler。action 字符串只负责定义 HTTP method 和可选子路径。

```go
api.NewRESTResource(
	"users",
	api.WithOperations(
		api.OperationSpec{
			Action:  "get",
			Public:  true,
			Handler: "List",
		},
		api.OperationSpec{
			Action:  "post admin",
			Handler: "CreateAdmin",
		},
	),
)
```

REST action 格式：

| 格式 | 含义 | 示例 |
| --- | --- | --- |
| `<method>` | 根资源路由 | `get`、`post`、`delete` |
| `<method> <sub-resource>` | 子资源路由 | `get profile`、`post admin`、`get user-friends` |

规则：

- HTTP verb 保持小写
- sub-resource 使用 kebab-case
- REST action 字符串不会推导 handler 方法名

## 支持的 Handler 返回形态

直接 handler 只能返回“无返回值”或“一个 `error`”。

| 形态 | 含义 |
| --- | --- |
| `func(...)` | 不显式返回错误 |
| `func(...) error` | 标准 handler 形态 |

其他返回形态都不是合法的直接 handler。

## 支持的 Factory 形态

VEF 还支持启动期 handler factory。Factory 返回一个 handler 闭包，并可选再返回一个 `error`。

| 形态 | 含义 |
| --- | --- |
| `func(...) func(...) error` | factory 返回一个 handler |
| `func(...) (func(...) error, error)` | factory 返回 handler 和启动期错误 |
| `func(...) func(...)` | factory 返回一个无错误返回值的 handler |
| `func(...) (func(...), error)` | factory 返回无错误 handler 和启动期错误 |

当某些依赖希望在应用启动时只解析一次，而不是每次请求都解析时，factory 很合适。

## 内置 Handler 参数注入

框架内置支持以下 handler 参数解析器：

| 参数类型 | 来源 | 常见用途 |
| --- | --- | --- |
| `fiber.Ctx` | 请求上下文 | 直接操作 Fiber 请求/响应 |
| `orm.DB` | 请求上下文 | 查询和事务入口 |
| `log.Logger` | 请求上下文 | 请求级日志 |
| `*security.Principal` | 请求上下文 | 当前登录主体 |
| `api.Params` | 请求体 / query / form | 原始 params map |
| `api.Meta` | 请求 meta | 原始 meta map |
| 嵌入 `api.P` 的 typed struct | 请求 params | 强类型 params 解码 + 自动验证 |
| 嵌入 `api.M` 的 typed struct | 请求 meta | 强类型 meta 解码 + 自动验证 |
| `page.Pageable` | 请求 meta | 分页元信息 |
| `cron.Scheduler` | DI 容器 | 调度器访问 |
| `event.Publisher` | DI 容器 | 事件发布 |
| `mold.Transformer` | DI 容器 | 输出转换 |
| `storage.Service` | DI 容器 | 存储访问 |

重要解码规则：

- typed params 和 typed meta 结构体都会自动执行 `validator.Validate(...)`
- `page.Pageable` 被视为内置 meta helper
- `api.Params` 和 `api.Meta` 不会走 typed 解码和验证

## 字段与资源注入

如果某个 handler 参数类型不在内置解析器里，VEF 会继续尝试从资源结构体本身解析。

解析顺序：

| 步骤 | 规则 |
| --- | --- |
| direct field match | 在资源上查找非匿名、类型兼容的字段 |
| tagged dive field match | 继续搜索 `api:"dive"` 标记的嵌套字段 |
| embedded field match | 搜索匿名嵌入字段 |

这意味着你可以把自定义 service、repository、helper 作为资源字段放进去，再直接在 handler 参数里按类型拿到它们。

## 支持的 Factory 参数

启动期 factory resolver 比请求期 handler resolver 更窄。框架内置支持这些 factory 参数类型：

| 参数类型 | 来源 |
| --- | --- |
| `orm.DB` | DI 容器 |
| `cron.Scheduler` | DI 容器 |
| `event.Publisher` | DI 容器 |
| `mold.Transformer` | DI 容器 |
| `storage.Service` | DI 容器 |

除此之外，factory 参数也支持从资源结构体字段里按兼容类型解析。

## 常见 Handler 形态

### 精简 RPC handler

```go
func (*UserResource) ResetPassword(
	ctx fiber.Ctx,
	db orm.DB,
	principal *security.Principal,
	params *ResetPasswordParams,
) error {
	// business logic
	return result.Ok().Response(ctx)
}
```

### 原始代理风格 handler

```go
func (*DebugResource) Echo(ctx fiber.Ctx, params api.Params, meta api.Meta) error {
	return result.Ok(fiber.Map{
		"params": params,
		"meta":   meta,
	}).Response(ctx)
}
```

### 启动期 factory

```go
func (*UserResource) BuildReport(
	db orm.DB,
) (func(ctx fiber.Ctx, params ReportParams) error, error) {
	return func(ctx fiber.Ctx, params ReportParams) error {
		return result.Ok().Response(ctx)
	}, nil
}
```

## 什么时候用自定义 Handler

以下场景适合自定义 handler：

- 该动作不是标准 CRUD
- 一个接口要编排多个服务
- 响应结构高度业务化
- 这个动作会触发流程、事件、外部集成
- 请求契约更适合直接表达，而不是硬塞进 CRUD builder

## CRUD 与自定义动作混用

真实项目里很常见的模式是：

- 用嵌入 CRUD builder 覆盖标准动作
- 再用显式 `api.OperationSpec` 补那几个不适合 CRUD 的业务动作

这样一个资源仍然保持围绕一个业务域组织，同时也能容纳诸如“保存权限”“发布版本”“查询关联用户”这类额外动作。

## 实践建议

- handler 尽量保持精简，聚焦业务编排
- 优先使用 typed params，而不是原始 `api.Params`
- 原始 map 更适合动态代理类接口，不适合长期维护的业务接口
- 只有在确实需要启动期依赖整形时才使用 factory
- 对于 REST 资源，handler 映射必须显式写清楚
- 如果参数可以建模成 typed `api.P` / `api.M`，优先这么做，不要手工从 `fiber.Ctx` 里拆

## 下一步

继续阅读 [路由](./routing) 和 [参数与元信息](./params-and-meta)，看请求形态和解码规则的细节。
