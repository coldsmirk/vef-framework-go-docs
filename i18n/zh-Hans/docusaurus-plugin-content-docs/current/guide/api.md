---
sidebar_position: 1
---

# API 包

`api` 包是 VEF 请求处理层的基础。它定义了核心抽象 — 资源、操作、请求/响应类型和处理器解析 — 所有其他包都构建在此之上。

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

使用 CRUD builder 时，通常嵌入 `api.Resource` 和 CRUD provider。框架会自动从所有嵌入的 `OperationsProvider` 实现中收集操作。

## OperationSpec

`OperationSpec` 是 API 端点的静态定义：

```go
type OperationSpec struct {
    Action      string             // Action 名（如 "create"、"find_page"）
    EnableAudit bool               // 启用审计日志
    Timeout     time.Duration      // 请求超时
    Public      bool               // 无需认证
    PermToken   string             // 所需权限令牌
    RateLimit   *RateLimitConfig   // 限流配置
    Handler     any                // 业务处理函数
}
```

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
        FindPage: crud.NewFindPage[User, UserSearch]().PermToken("sys:user:query"),
        Create:   crud.NewCreate[User, UserParams]().PermToken("sys:user:create"),
        Update:   crud.NewUpdate[User, UserParams]().PermToken("sys:user:update"),
        Delete:   crud.NewDelete[User]().PermToken("sys:user:delete"),
    }
}
```

### 配合自定义 Handler 使用

对于非 CRUD 操作，使用 `WithOperations` 或实现 `OperationsProvider`：

```go
resource := api.NewRPCResource("sys/user",
    api.WithOperations(
        api.OperationSpec{
            Action:  "reset_password",
            Handler: resetPasswordHandler,
            PermToken: "sys:user:reset_password",
        },
    ),
)
```

## Identifier

每个操作都有唯一的 `Identifier`：

```go
type Identifier struct {
    Resource string  // 如 "sys/user"
    Action   string  // 如 "create"
    Version  string  // 如 "v1"
}

// 字符串格式："sys/user:create:v1"
id.String()
```

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

### Params

`Params` 承载请求的业务数据（"做什么"）：

```go
type Params map[string]any

// 解码到类型化结构体
var userParams UserParams
err := request.Params.Decode(&userParams)

// 访问单个值
value, exists := request.GetParam("username")
```

### Meta

`Meta` 承载请求元数据（"怎么做" — 分页、排序、格式）：

```go
type Meta map[string]any

// 解码到类型化结构体
var pageable page.Pageable
err := request.Meta.Decode(&pageable)

// 访问单个值
value, exists := request.GetMeta("format")
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
    Strategy string         // "none"、"bearer"、"signature" 或自定义
    Options  map[string]any // 策略特定选项
}
```

### 内置认证策略

| 策略 | 常量 | 说明 |
| --- | --- | --- |
| 无认证 | `api.AuthStrategyNone` | 公开访问 |
| Bearer | `api.AuthStrategyBearer` | Bearer 令牌认证 |
| 签名 | `api.AuthStrategySignature` | 请求签名认证 |

### 辅助函数

```go
api.Public()        // 策略为 "none" 的 AuthConfig
api.BearerAuth()    // 策略为 "bearer" 的 AuthConfig
api.SignatureAuth() // 策略为 "signature" 的 AuthConfig
```

### 资源级 vs 操作级认证

```go
// 资源级：所有操作使用签名认证
api.NewRPCResource("external/webhook", api.WithAuth(api.SignatureAuth()))

// 操作级：按操作覆盖
crud.NewCreate[User, UserParams]().Public()                      // 无认证
crud.NewFindPage[User, UserSearch]().PermToken("sys:user:query") // Bearer + 权限
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

## 下一步

- 阅读 [泛型 CRUD](./crud) 了解 CRUD builder 如何自动生成操作
- 阅读 [自定义 Handler](./custom-handlers) 创建非 CRUD 操作
- 阅读 [路由](./routing) 了解 HTTP 路由细节
- 阅读 [Params 与 Meta](./params-and-meta) 了解请求数据契约
