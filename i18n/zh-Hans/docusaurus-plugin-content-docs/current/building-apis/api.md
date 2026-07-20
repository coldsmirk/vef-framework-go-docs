---
sidebar_position: 1
---

# API 包

`api` 包是 VEF 请求处理层的基础。它定义了核心抽象 — 资源、操作、请求/响应类型和处理器解析 — 所有其他包都构建在此之上。

## API 参考

生成的 [公开 API 索引](../reference/public-api-index) 是所有 exported
symbol、exported field 和 exported method 的无遗漏清单。本指南负责说明这些
surface 背后的受支持行为和运行时 contract。

本指南涵盖的 API 组：

| 分组 | 公开 API |
| --- | --- |
| resource 和 kind | `api.Resource`, `api.Kind`, `api.KindRPC`, `api.KindREST`, `api.ValidateActionName(action, kind) error`, `api.NewRPCResource(name, opts...)`, `api.NewRESTResource(name, opts...)`, `api.WithVersion(v)`, `api.WithAuth(config)`, `api.WithOperations(specs...)` |
| engine 和 routing 扩展 | `api.Engine`, `api.RouterStrategy`, `api.Middleware` |
| operations | `api.OperationSpec`, `api.Operation`, `api.RateLimitConfig`, `api.OperationsProvider`, `api.OperationsCollector` |
| request model | `api.Identifier`, `api.Request`, `api.Params`, `api.Meta`, `api.P`, `api.StrictP`（v0.39）, `api.M` |
| auth | `api.AuthConfig`, `api.Public()`, `api.BearerAuth()`, `api.SignatureAuth()`, `api.IPAuth(...)`, `api.APIKeyAuth(...)`（v0.39）, `api.HTTPBasicAuth()`（v0.39）, `api.AuthStrategy`, `api.AuthStrategyRegistry` |
| handler 扩展 | `api.HandlerResolver`, `api.HandlerAdapter`, `api.HandlerParamResolver`, `api.FactoryParamResolver` |
| audit、headers、versions、errors | `api.AuditEvent`, `api.SubscribeAuditEvent`, `api.HeaderXMetaPrefix`, `api.HeaderXTimestamp`, `api.HeaderXNonce`, `api.HeaderXSignature`, `api.HeaderXAppID`, `api.VersionV1`, `api.VersionV9`, `api.ErrInvalidRequestParams`, `api.ErrInvalidRequestMeta`, `api.ErrInvalidParamsType`, `api.ErrInvalidMetaType` |

## 架构

```
api.Engine
├── Register(resources...)       — 注册资源
├── Mount(router)                — 挂载到 Fiber
└── Lookup(id)                   — 运行时查找操作

api.Resource
├── Kind()        — RPC 或 REST
├── Name()        — 资源路径（如 "sys/user"）
├── Version()     — API 版本（如 "v1"）
├── Auth()        — 认证配置
└── Operations()  — 操作规格列表

api.OperationSpec → api.Operation（运行时）
```

## Resource

`Resource` 将相关的 API 操作归组到一个公共路径下。VEF 提供两种资源类型：

### 创建资源

```go
// RPC 资源 — 使用 snake_case 命名 action
resource := api.NewRPCResource("sys/user")

// REST 资源 — 使用 HTTP 动词
resource := api.NewRESTResource("sys/user")

// 带选项
resource := api.NewRPCResource("sys/user",
    api.WithVersion("v2"),
    api.WithAuth(api.BearerAuth()),
    api.WithOperations(
        api.OperationSpec{Action: "create", Handler: createHandler},
        api.OperationSpec{Action: "find_page", Handler: findPageHandler},
    ),
)
```

### 资源类型（Kind）

| 类型 | 常量 | 名称格式 | Action 格式 | 示例 |
| --- | --- | --- | --- | --- |
| RPC | `api.KindRPC` | `snake_case` + `/` 分隔 | `snake_case` | `sys/user` → `create`、`find_page` |
| REST | `api.KindREST` | `kebab-case` + `/` 分隔 | `<动词>` 或 `<动词> <子资源>` | `sys/user` → `get`、`post`、`get user-friends` |

`Kind.String()` 对 `KindRPC` 返回 `rpc`，对 `KindREST` 返回 `rest`，其他值返回
`unknown`。

### 资源名称规则

| 规则 | 合法 | 非法 |
| --- | --- | --- |
| 必须以小写字母开头 | `user`、`sys/user` | `User`、`1user` |
| 不能以斜杠开头/结尾 | `sys/user` | `/sys/user/` |
| 不能有连续斜杠 | `sys/user` | `sys//user` |
| RPC：snake_case 分段 | `sys/data_dict` | `sys/data-dict` |
| REST：kebab-case 分段 | `sys/data-dict` | `sys/data_dict` |

### 资源选项

| 选项 | 说明 |
| --- | --- |
| `api.WithVersion(v)` | 覆盖引擎的默认版本（如 `"v2"`）|
| `api.WithAuth(config)` | 设置资源级认证 |
| `api.WithOperations(specs...)` | 直接提供操作规格 |

`api.NewRPCResource` 和 `api.NewRESTResource` 会在构造期校验 resource
name、version，以及通过直接 `api.WithOperations(...)` 传入的 specs；校验失败会
panic。`api.ValidateActionName(action, kind)` 是公开函数，适合动态构建
resource 的代码在创建 `OperationSpec` 前复用框架同一套 RPC/REST action
校验规则。

REST action 校验接受这些小写 method token：`get`、`post`、`put`、
`delete`、`patch`、`head`、`options`、`trace`、`connect` 和 `all`。
sub-resource 路径可以包含 `/`，但每一段都必须是 kebab-case；动态 Fiber
参数如 `/:id` 不会被公开 validator 接受。

## Resource 接口

```go
type Resource interface {
    Kind() Kind
    Name() string
    Version() string
    Auth() *AuthConfig
    Operations() []OperationSpec
}
```

使用 CRUD builder 时，通常嵌入 `api.Resource` 和 CRUD provider。内置
collector 会读取直接的 `Resource.Operations()` specs，也会读取匿名嵌入的
`OperationsProvider`。直接 `WithOperations(...)` 传入的 specs 会由 resource
constructor 校验。engine 注册期间，收集到的 operation specs 必须有非空
action；自定义 `OperationsProvider` 仍应产出已经满足 `api.ValidateActionName(...)`
的 action 字符串。

## OperationSpec

`OperationSpec` 是 API 端点的静态定义：

```go
type OperationSpec struct {
    Action             string            // Action 名（如 "create"、"find_page"）
    EnableAudit        bool              // 启用审计日志
    Timeout            time.Duration     // 请求超时
    Public             bool              // 无需认证
    RequiredPermission string            // 所需权限令牌（v0.24 起从 PermToken 重命名）
    RateLimit          *RateLimitConfig  // 限流配置
    Handler            any               // 业务处理函数
}
```

运行期 engine 会把它物化为带最终 `Auth` 和 `RateLimit` 指针的 `Operation`，
供 router strategy、middleware、诊断和测试使用。`Operation` 和 `Request` 都嵌入了 `Identifier`，所以
`Identifier.String()` 会提升为 `Operation.String()` 和 `Request.String()`。

Operation 默认化规则：

| 字段 | 运行时行为 |
| --- | --- |
| `Action` | 必填；直接 `WithOperations(...)` 传入的 specs 会由 resource constructor 按 resource kind 校验；engine 注册会拒绝空 action |
| `EnableAudit` | 直接复制到运行时 operation |
| `Timeout` | 非正值使用 engine 默认值；未覆盖时默认 `30s` |
| `Public` | 为 `true` 时先解析为 `api.Public()`，优先于资源级/default auth |
| `RequiredPermission` | 非空时写入 auth options 中的 required permission token |
| `RateLimit` | nil 使用 engine 默认——配置了 `vef.api.rate_limit` 时用配置值，否则 `Max=100`、`Period=5m`（v0.38）；自定义 `RateLimitConfig` 会替换默认值。显式 `Max <= 0` **不会**关闭限流——中间件对非正值回退到 engine/配置默认值 |
| `Handler` | RPC 可在省略时从 action 推断；REST 必须显式提供 handler |

自 v0.39 起，权限声明在注册时校验，违规将导致启动失败：

- `RequiredPermission` 必须是匹配 `^[A-Za-z0-9_]+(\.[A-Za-z0-9_]+)*$` 的
  **点分隔**令牌（如 `sys.user.query`）；冒号、斜杠、连字符、空段与空白
  均被 `ErrPermissionTokenInvalid` 拒绝。
- 在解析后认证策略为 `none` 的操作上（`Public: true`，或资源级
  `api.Public()`）声明 `RequiredPermission` 是自相矛盾的——匿名主体永远
  无法满足它——会被 `ErrPermissionOnPublicOp` 拒绝。

### 配合 CRUD Builder 使用

大多数操作通过 CRUD builder 定义，而非直接使用 `OperationSpec`：

```go
type UserResource struct {
    api.Resource

    crud.FindPage[User, UserSearch]
    crud.Create[User, UserParams]
    crud.Update[User, UserParams]
    crud.Delete[User]
}

func NewUserResource() *UserResource {
    return &UserResource{
        Resource: api.NewRPCResource("sys/user"),
        FindPage: crud.NewFindPage[User, UserSearch]().RequiredPermission("sys.user.query"),
        Create:   crud.NewCreate[User, UserParams]().RequiredPermission("sys.user.create"),
        Update:   crud.NewUpdate[User, UserParams]().RequiredPermission("sys.user.update"),
        Delete:   crud.NewDelete[User]().RequiredPermission("sys.user.delete"),
    }
}
```

### 配合自定义 Handler 使用

对于非 CRUD 操作，使用 `WithOperations` 或实现 `OperationsProvider`：

```go
resource := api.NewRPCResource("sys/user",
    api.WithOperations(
        api.OperationSpec{
            Action:             "reset_password",
            Handler:            resetPasswordHandler,
            RequiredPermission: "sys.user.reset_password",
        },
    ),
)
```

## Identifier

每个操作都有唯一的 `Identifier`：

```go
type Identifier struct {
    Resource string `json:"resource" form:"resource" validate:"required,alphanum_us_slash" label_i18n:"api_request_resource"`
    Action   string `json:"action" form:"action" validate:"required" label_i18n:"api_request_action"`
    Version  string `json:"version" form:"version" validate:"required,alphanum" label_i18n:"api_request_version"`
}

// 字符串格式："sys/user:create:v1"
id.String()
```

JSON RPC body 和 form RPC request 都使用同一组 `Identifier` 字段。
`Identifier.String()` 始终格式化为 `{resource}:{action}:{version}`。

## Request / Params / Meta

### Request

统一的 API 请求结构：

```go
type Request struct {
    Identifier
    Params Params `json:"params"`
    Meta   Meta   `json:"meta"`
}
```

如果希望 handler 注入明确从 `params` 或 `meta` 解码，请分别在请求参数结构体
中嵌入 `api.P`，在 metadata 结构体中嵌入 `api.M`。

### Params

`Params` 承载请求的业务数据（"做什么"）：

```go
type Params map[string]any

// 解码到类型化结构体
var userParams UserParams
err := request.Params.Decode(&userParams)

// 访问单个值
value, exists := request.Params["username"]
```

### Meta

`Meta` 承载请求元数据（"怎么做" — 分页、排序、格式）：

```go
type Meta map[string]any

// 解码到类型化结构体
var pageable page.Pageable
err := request.Meta.Decode(&pageable)

// 访问单个值
value, exists := request.Meta["format"]
```

### Params vs Meta

| 维度 | `Params` | `Meta` |
| --- | --- | --- |
| 用途 | 业务数据 | 请求控制 |
| 解码目标 | `TParams` 或 `TSearch` | `page.Pageable`、`DataOptionConfig` 等 |
| 示例 | `username`、`email`、`departmentId` | `page`、`size`、`sort`、`format` |

## 认证

### AuthConfig

```go
type AuthConfig struct {
    Strategy string         // "none"、"bearer"、"signature"、"ip" 或自定义
    Options  map[string]any // 策略特定选项
}
```

`AuthConfig.Clone()` 会复制 strategy/options。资源或操作需要调整 auth 但不
想修改共享配置时使用它。

### 内置认证策略

| 策略 | 常量 | 说明 |
| --- | --- | --- |
| 无认证 | `api.AuthStrategyNone` | 公开访问 |
| Bearer | `api.AuthStrategyBearer` | Bearer 令牌认证 |
| 签名 | `api.AuthStrategySignature` | 请求签名认证 |
| IP | `api.AuthStrategyIP` | 来源 IP 白名单认证 |
| API key | `api.AuthStrategyAPIKey` | 静态 API key 认证（v0.39） |
| HTTP Basic | `api.AuthStrategyHTTPBasic` | RFC 7617 Basic 认证（v0.39） |

### 辅助函数

```go
api.Public()               // 策略为 "none" 的 AuthConfig
api.BearerAuth()           // 策略为 "bearer" 的 AuthConfig
api.SignatureAuth()        // 策略为 "signature" 的 AuthConfig
api.IPAuth()               // 策略为 "ip"，使用 "default" whitelist
api.IPAuth("ops")          // 策略为 "ip"，使用 "ops" whitelist
api.APIKeyAuth()           // 策略为 "api_key"，读 X-API-Key 头（v0.39）
api.APIKeyAuth("X-My-Key") // 自定义 key 头；传多个名称会 panic
api.HTTPBasicAuth()        // 策略为 "http_basic"（v0.39）
```

`api.IPAuth(...)` 接受 0 或 1 个 whitelist 名称。不传时使用
`api.DefaultIPWhitelist`（`"default"`）；选中的名称会写入
`AuthConfig.Options` 的 `api.AuthOptionWhitelist`。传入多个名称会 panic。
内置 IP strategy 通过 `security.IPWhitelistLoader` 解析命名列表；默认 loader
读取 `vef.security.ip_whitelists`。所有认证失败都返回
`security.ErrIPNotAllowed`，缺失或为空的命名 whitelist 会 fail-closed，而不会
降级为公开访问。位于反向代理之后时，需要配置 `vef.app.trusted_proxies`，让
Fiber 解析到真实客户端 IP。

`api.APIKeyAuth(...)` 默认从 `api.HeaderXAPIKey`（`X-API-Key`）读取密钥；
可以传一个自定义头名，写入 `api.AuthOptionAPIKeyHeader`。密钥经
`security.APIKeyLoader` 解析（默认读 `vef.security.api_keys`）。
`api.HTTPBasicAuth()` 经 `security.BasicAccountLoader`（默认读
`vef.security.basic_accounts`）以常数时间比较验证 `Authorization: Basic`
凭证。两者都统一以 401 拒绝（`ErrAPIKeyInvalid` /
`ErrBasicCredentialsInvalid`）；loader 契约见[认证](../security/authentication)。

自定义认证策略实现 `api.AuthStrategy`，并通过
`vef.ProvideAuthStrategy(...)` 注册到 `vef:api:auth_strategies`。

### 资源级 vs 操作级认证

```go
// 资源级：所有操作使用签名认证
api.NewRPCResource("external/webhook", api.WithAuth(api.SignatureAuth()))

// 操作级：按操作覆盖
crud.NewCreate[User, UserParams]().Public()                      // 无认证
crud.NewFindPage[User, UserSearch]().RequiredPermission("sys.user.query") // Bearer + 权限
```

## 限流

```go
type RateLimitConfig struct {
    Max    int           // 允许的最大请求数
    Period time.Duration // 时间窗口
    Key    string        // 自定义限流 key（可选）
}
```

通过 CRUD builder 使用：

```go
crud.NewCreate[User, UserParams]().RateLimit(100, time.Minute)
```

内置 rate limiter 使用 sliding window，并按节点独立计数。未声明自己
`RateLimit` 的 operation 使用 engine 默认值，v0.38 起可由用户配置：

```toml
[vef.api.rate_limit]
max    = 100   # 省略时的默认值
period = "5m"  # 省略时的默认值
```

显式提供 `RateLimitConfig` 且 `Max <= 0` 时**不会**关闭限流：中间件只采用
正值，其余情况回退到 engine/配置默认值（恒为正）——不存在按操作关闭限流的
开关。框架默认 key 包含 resource、version、action、解析后的客户端 IP 和
principal ID；匿名请求使用 anonymous principal。

Operation 认证配置由 `Operation.Auth` 承载；公开 operation 会解析为
`api.AuthStrategyNone`，受保护 operation 则携带选中的 auth strategy 和 options。

## Engine

`Engine` 管理资源注册和 HTTP 路由：

```go
type Engine interface {
    Register(resources ...Resource) error  // 添加资源
    Lookup(id Identifier) *Operation       // 运行时查找操作
    Mount(router fiber.Router) error       // 挂载到 Fiber 路由器
}
```

## Handler 解析

VEF 通过参数注入支持灵活的 handler 签名：

```go
// 最简 handler
func (r *UserResource) Create(ctx fiber.Ctx) error { ... }

// 自动注入参数
func (r *UserResource) Create(ctx fiber.Ctx, db orm.DB, principal *security.Principal) error { ... }

// 带类型化 params
func (r *UserResource) Create(ctx fiber.Ctx, params UserParams, db orm.DB) error { ... }
```

### 参数解析接口

| 接口 | 用途 |
| --- | --- |
| `HandlerParamResolver` | 运行时从请求上下文解析 handler 参数 |
| `FactoryParamResolver` | 启动时解析 handler 参数（依赖注入）|
| `HandlerAdapter` | 将任意 handler 签名转换为 `fiber.Handler` |
| `HandlerResolver` | 在资源上查找 handler 函数 |

`RouterStrategy` 也是公开扩展点，用于自定义 HTTP 暴露方式。strategy 声明
自己能处理哪些 `api.Kind`，在 `Setup` 中接收 Fiber router，并在 `Route`
中注册每个已解析操作。

## 错误类型

| 错误 | 含义 |
| --- | --- |
| `ErrEmptyResourceName` | 资源名称为空 |
| `ErrInvalidResourceName` | 资源名称不符合命名规则 |
| `ErrResourceNameSlash` | 资源名称以 `/` 开头或结尾 |
| `ErrResourceNameDoubleSlash` | 资源名称包含 `//` |
| `ErrInvalidResourceKind` | 无效的资源类型值 |
| `ErrInvalidVersionFormat` | 版本不匹配 `v\d+` 格式 |
| `ErrEmptyActionName` | Action 名称为空 |
| `ErrInvalidActionName` | Action 不符合类型特定规则 |
| `ErrInvalidParamsType` | Params.Decode 目标不是指向结构体的指针 |
| `ErrInvalidMetaType` | Meta.Decode 目标不是指向结构体的指针 |

`Params.Decode` 和 `Meta.Decode` 都要求传入 struct 指针。其他目标会返回
`ErrInvalidParamsType` 或 `ErrInvalidMetaType`。

`ErrInvalidRequestParams` 和 `ErrInvalidRequestMeta` 使用
`result.ErrCodeBadRequest`（`1400`）和 HTTP status `400`。RPC form 的
`params`/`meta` JSON 或 REST JSON body 解析失败时会返回它们。

`api` 包的其余公开接口面——版本常量、请求 header、审计事件、认证策略
registry、operation 收集类型、请求 helper 与 handler/router 扩展接口——在
[公开 API 索引](../reference/public-api-index)中逐符号收录。

## 下一步

- 阅读 [泛型 CRUD](../data-access/crud) 了解 CRUD builder 如何自动生成操作
- 阅读 [自定义 Handler](./custom-handlers) 创建非 CRUD 操作
- 阅读 [路由](./routing) 了解 HTTP 路由细节
- 阅读 [Params 与 Meta](./params-and-meta) 了解请求数据契约
