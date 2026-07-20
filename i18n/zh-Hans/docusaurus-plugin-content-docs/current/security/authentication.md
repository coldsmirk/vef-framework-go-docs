---
sidebar_position: 1
---

# 认证

VEF 的认证发生在 API 操作层。每个操作都有自己的 auth 配置，API 中间件会在 handler 执行前先解析出当前 principal。

## 默认行为

如果你不做额外配置：

- 操作默认使用 Bearer 认证
- 显式 `Public` 的操作则不要求认证

这个默认值来自 API 引擎，而不是来自你的应用配置文件。

## 内置认证策略

公开的 `api` 包暴露了这些策略 helper：

- `api.Public()`
- `api.BearerAuth()`
- `api.SignatureAuth()`
- `api.IPAuth(...)`（白名单解析方式见下文的 [Signature helpers](./authentication-reference#signature-helpers)）
- `api.APIKeyAuth(...)`
- `api.HTTPBasicAuth()`

实际使用中，你通常通过操作配置来控制：

```go
api.OperationSpec{
	Action: "login",
	Public: true,
}
```

或者通过资源级别的 auth 默认值来设置。

## Bearer 认证

Bearer 认证支持两种 token 来源：

- `Authorization: Bearer <token>`
- 查询参数 `__accessToken`

真正的 token 校验逻辑由安全模块中的 auth manager 负责。

## Signature 认证

Signature 认证主要用于外部应用和请求签名场景。

它要求这些 header：

- `X-App-ID`
- `X-Timestamp`
- `X-Nonce`
- `X-Signature`

校验逻辑由安全模块的 signature authenticator 执行。

## API Key 认证

`api.APIKeyAuth()` 以静态密钥认证机器对机器的调用方，默认从 `X-API-Key`
头读取；传入一个头名（`api.APIKeyAuth("X-Custom-Key")`）可改用自定义头。

提交的 key 通过已注册的 `security.APIKeyLoader` 解析。框架内置基于配置的
loader，读取 `vef.security.api_keys`（对全部配置项做常数时间比较）：

```toml
[vef.security.api_keys.reporting]
key = "high-entropy-random-string"
roles = ["reporting"]
```

应用可以提供自己的 `security.APIKeyLoader`，从数据库或配置中心加载：

```go
type APIKeyLoader interface {
    // LoadByKey 把提交的 key 解析为其 Principal；无匹配时返回 nil。
    // 返回 error 表示基础设施故障，而非拒绝。
    LoadByKey(ctx context.Context, key string) (*security.Principal, error)
}
```

扫描候选 key 的实现必须做常数时间比较；按 key 建索引的实现只应服务高熵
随机 key——此时查找时间不泄露任何信息。

缺失或未匹配的 key 统一以 `security.ErrAPIKeyInvalid`（HTTP 401）拒绝。
基于配置的 loader 将命中解析为以条目命名的外部应用主体
（`api_key:<name>`），携带配置的角色。

## HTTP Basic 认证

`api.HTTPBasicAuth()` 认证 RFC 7617 的 `Authorization: Basic` 凭证。这是
机器对机器的服务账号——请存放高熵随机密钥，不要放用户密码。

框架内置基于配置的 loader，读取 `vef.security.basic_accounts`
（map 键即用户名）：

```toml
[vef.security.basic_accounts.metrics-scraper]
password = "high-entropy-random-string"
roles = ["metrics"]
```

应用可以提供自己的 `security.BasicAccountLoader`；loader 返回已存密钥，
常数时间比较由框架执行，因此所有实现共享同一 fail-closed 语义：

```go
type BasicAccountLoader interface {
    // LoadByUsername 按用户名取服务账号，返回 Principal 与已存密钥。
    // nil Principal 或空密钥表示账号未知；error 表示基础设施故障。
    LoadByUsername(ctx context.Context, username string) (*security.Principal, string, error)
}
```

畸形请求头、未知账号与错误密码统一以
`security.ErrBasicCredentialsInvalid`（HTTP 401）拒绝，调用方无法区分
失败的是哪一部分。

## 保留身份

某些身份用于归因框架在请求之外执行的工作——它们是审计作者，绝不是调用方。
`security.Principal.IsReserved()` 报告它们：`system` 主体类型，以及
`orm.OperatorSystem` / `orm.OperatorCronJob` 操作者 ID。

框架在每个边界 fail-closed 地强制该不变量：

- 解析出 nil 或保留主体的认证器（包括自定义 `security.Authenticator`
  实现与挑战提供者）在认证边界即被拒绝；
- 令牌签发与挑战流程拒绝为保留身份签发令牌
  （`security_reserved_principal_forbidden`，HTTP 401）。

内置密码登录还额外拒绝以 `anonymous` 作为登录标识（它表示"身份缺失"，
只有 public 策略可以合法产出）。自定义 `UserLoader`、`APIKeyLoader` 等
实现绝不能返回 ID 与保留操作者 ID 冲突的主体。

## 公开操作

公开操作会得到一个匿名 principal，而不是直接被拒绝。

适合标记为 `Public` 的接口包括：

- 登录
- 刷新 token
- 某些匿名健康检查或回调入口

## 内置认证资源

安全模块会自动注册一个内置 RPC 资源：

```text
security/auth
```

主要 actions 包括：

- `login`
- `refresh`
- `logout`
- `resolve_challenge`
- `get_user_info`

这些请求字段、公开标记和限流来源也是运行时 contract 的一部分：

| Action | Public | Rate limit | 请求字段 |
| --- | --- | --- | --- |
| `login` | 是 | `vef.security.login_rate_limit` | `type`、`principal`、`credentials`；全部是 `validate:"required"` |
| `refresh` | 是 | `vef.security.refresh_rate_limit` | `refreshToken`；`validate:"required"`。仅在 `token_type = "jwt_token"` 下挂载——`opaque_token` 下该操作不存在（会话自行续期） |
| `logout` | 否 | 默认 API rate limit | 无 |
| `resolve_challenge` | 是 | `vef.security.login_rate_limit` | `challengeToken`、`type`、`response`；全部是 `validate:"required"` |
| `get_user_info` | 否 | 默认 API rate limit | 任意 `params`，会转发给 `UserInfoLoader.LoadUserInfo(...)` |

完整的字段级契约——每个 action 的请求参数**和**响应字段，含登录响应
两种形态的 JSON 示例——收录在
[RPC Resource: `security/auth`](./authentication-reference#rpc-resource-securityauth)。

这个资源、所有已注册的 `Authenticator`，以及 `AuthManager` 聚合器，都由
框架的安全模块完成装配——同一个模块还装配了暴力破解锁定、密码强度/历
史/过期（参见[登录加固](./login-hardening)）以及 opaque token 会话控制
（参见[会话管理](./session-management)）。

内置 authenticator type 字符串是 `password`、`jwt_token`、`opaque_token`、
`refresh` 和 `signature`（JWT 认证器名为 `jwt_token`，不叫
`token`）。普通客户端调用里，`security/auth.login` 使用
`type: "password"` 搭配用户名和密码凭证。Bearer 保护的操作会在内部按
`vef.security.token_type` 分派已配置的令牌机制（`jwt_token` 或
`opaque_token`），`security/auth.refresh` 会在内部使用 `refresh`，
`SignatureAuth` 会把签名 headers 映射到 `signature` authenticator。只有
已配置机制的认证器会被注册，且 `login` 拒绝框架签发的令牌类型作为登录
凭据（见[会话管理](./session-management)）。

`logout` 总是返回 ok 结果。在 `jwt_token` 下它实际上是 no-op——服务端没有
可吊销的会话，客户端需要自行删除已保存的 token。在 `opaque_token` 下它会
吊销当前 bearer token 背后的会话，尽力而为（会话不存在或存储出错只记录
日志）。

## 登录流程

内置认证资源支持两阶段模型：

1. 先校验凭证
2. 如有需要，再进入 challenge 流程

如果没有 challenge，`login` 会直接返回 token。

如果 challenge provider 已配置且当前用户需要额外挑战，`login` 会返回：

- challenge token
- 下一步 challenge 描述

客户端之后继续调用 `resolve_challenge`，直到所有挑战都完成。
Go API 层里，这个响应形状由 `LoginResult` 表示；当前步骤由 `LoginChallenge` 表示。

登录响应 DTO 使用这些精确字段：

| DTO | 字段 |
| --- | --- |
| `AuthTokens` | JSON `accessToken`、`refreshToken` |
| `Authentication` | JSON `type`、`principal`、`credentials` |
| `LoginResult` | JSON `tokens`、`challengeToken`、`challenge` |
| `LoginChallenge` | JSON `type`、`data`、`required` |
| `ChallengeState` | Go-only state: `Principal`、`Username`（第一步提交的原始登录标识，跨挑战步骤保留、用于审计事件）、`Pending`、`Resolved` |

两种响应形态——token 载荷与 challenge 包络——的逐字段表格及 JSON 示例见
[RPC Resource: `security/auth`](./authentication-reference#rpc-resource-securityauth)。

## 应用通常还需要提供什么

这里要分场景来看：

- `security.UserLoader` 通常是用户登录和 refresh 流程的前提
- `security.ExternalAppLoader` 只在你使用签名认证的外部应用场景时需要
- challenge provider 是可选项，只有在你启用了挑战式登录流时才相关
- `security.UserInfoLoader` 只在你希望 `security/auth.get_user_info` 返回应用自定义用户信息时需要

框架提供的是认证流程和中间件，而不是你的身份源本身。

## 公开 API 一览

认证相关的完整公开接口面——principal、JWT、认证管理器、挑战提供者与令牌存储、签名认证、登录事件——连同契约说明收录在[认证参考](./authentication-reference)中。

## 一个可运行的登录模块

真实的 VEF 项目里，auth 模块通常都很小：一张用户表、一个实现 loader 接口的包，以及一个把它们提供出去的模块声明。框架已经内置了 `security/auth` 资源、password 和 refresh authenticator，以及默认的 bcrypt `password.Encoder`；应用只需要提供自己的身份源。下面这个 `auth` 包已经完整到可以对着一张真实的表登录。

### 用户模型

```go title="internal/auth/user.go"
package auth

import (
	"github.com/uptrace/bun"

	"github.com/coldsmirk/vef-framework-go/orm"
)

type User struct {
	bun.BaseModel `bun:"table:app_user,alias:au"`
	orm.FullAuditedModel

	Username     string `json:"username" validate:"required,alphanum,max=32" label:"Username"`
	Name         string `json:"name" validate:"required,max=32" label:"Name"`
	PasswordHash string `json:"-" bun:"password_hash,notnull"`
	Role         string `json:"role"`
	IsActive     bool   `json:"isActive"`
}

// UserDetails becomes Principal.Details and travels inside issued access tokens.
type UserDetails struct {
	Username string `json:"username"`
}
```

`PasswordHash` 必须存放与登录流程所用 `password.Encoder` 一致的输出——安全模块默认提供 bcrypt（`password.NewBcryptEncoder`）。在创建或初始化用户的地方注入 `password.Encoder`，存入 `encoder.Encode(plaintext)`；内置的 password authenticator 之后会用 `encoder.Matches(plaintext, storedHash)` 校验登录凭证。

### UserLoader

`security.UserLoader` 只有两个方法：`LoadByUsername` 支撑 `type: "password"` 登录，返回 principal 和已存储的密码哈希；`LoadByID` 支撑 token 刷新。

```go title="internal/auth/user_loader.go"
package auth

import (
	"context"

	"github.com/coldsmirk/vef-framework-go/orm"
	"github.com/coldsmirk/vef-framework-go/security"
)

type userLoader struct {
	db orm.DB
}

func NewUserLoader(db orm.DB) security.UserLoader {
	return &userLoader{db: db}
}

func (l *userLoader) LoadByUsername(ctx context.Context, username string) (*security.Principal, string, error) {
	user, err := l.findActive(ctx, "username", username)
	if err != nil {
		return nil, "", err
	}

	return toPrincipal(user), user.PasswordHash, nil
}

func (l *userLoader) LoadByID(ctx context.Context, id string) (*security.Principal, error) {
	user, err := l.findActive(ctx, "id", id)
	if err != nil {
		return nil, err
	}

	return toPrincipal(user), nil
}

func (l *userLoader) findActive(ctx context.Context, column string, value any) (*User, error) {
	var user User

	err := l.db.NewSelect().Model(&user).
		Where(func(cb orm.ConditionBuilder) {
			cb.Equals(column, value).IsTrue("is_active")
		}).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func toPrincipal(user *User) *security.Principal {
	principal := security.NewUser(user.ID, user.Name, user.Role)
	principal.Details = &UserDetails{Username: user.Username}

	return principal
}
```

这里的错误语义与内置 authenticator 的预期一致：

- `Scan` 已经把“无记录”映射为 `result.ErrRecordNotFound`，所以直接原样返回错误即可。查询里过滤 `is_active`，可以让被禁用的用户与不存在的用户表现完全一致。
- `login` 期间，`LoadByUsername` 返回的任何错误——以及 `nil` principal 或空哈希——都会归并为通用的凭证无效错误（code `1008`），因此无法通过响应枚举用户名。record-not-found 记录 info 级日志，其他错误记录 warn 级。
- `refresh` 期间，`LoadByID` 的错误会原样返回给调用方；refresh authenticator 之所以重新加载用户，正是为了让被停用的账号无法继续刷新。

### 权限与用户信息

```go title="internal/auth/loaders.go"
package auth

import (
	"context"

	"github.com/coldsmirk/vef-framework-go/security"
)

type rolePermissionsLoader struct{}

func NewRolePermissionsLoader() security.RolePermissionsLoader {
	return &rolePermissionsLoader{}
}

func (*rolePermissionsLoader) LoadPermissions(_ context.Context, role string) (map[string]security.DataScope, error) {
	if role == "admin" {
		return map[string]security.DataScope{
			"user:manage": security.NewAllDataScope(),
			"order:read":  security.NewAllDataScope(),
		}, nil
	}

	return map[string]security.DataScope{
		"order:read": security.NewSelfDataScope(""),
	}, nil
}

type userInfoLoader struct{}

func NewUserInfoLoader() security.UserInfoLoader {
	return &userInfoLoader{}
}

func (*userInfoLoader) LoadUserInfo(_ context.Context, principal *security.Principal, _ map[string]any) (*security.UserInfo, error) {
	return &security.UserInfo{
		ID:     principal.ID,
		Name:   principal.Name,
		Gender: security.GenderUnknown,
	}, nil
}
```

生产环境的 `RolePermissionsLoader` 应该读取角色-权限表，而不是写死的 switch；安全模块会自动把你提供的 loader 包上一层由 `RolePermissionsChangedEvent` 触发失效的缓存。权限 token 会进入[授权](./authorization)中描述的 RBAC 检查器。

### 装配

构造函数必须返回接口类型——框架从 DI 图中按 `security.UserLoader`、`security.UserInfoLoader`、`security.RolePermissionsLoader` 这些确切的接口类型消费可选依赖。

```go title="internal/auth/module.go"
package auth

import (
	"github.com/coldsmirk/vef-framework-go"
	"github.com/coldsmirk/vef-framework-go/security"
)

func init() {
	security.SetUserDetailsType[*UserDetails]()
}

var Module = vef.Module(
	"app:auth",
	vef.Provide(
		NewUserLoader,
		NewUserInfoLoader,
		NewRolePermissionsLoader,
	),
)
```

在 `main` 里把 `auth.Module` 传给 `vef.Run(...)`，内置的 `security/auth` 资源就会自动拿到这些 loader——不需要任何额外注册。这样认证接入代码就能和业务资源模块保持分离。

### 登录验证

假设已初始化一个用户（`admin` / `ChangeMe_123`，哈希由 bcrypt encoder 生成），调用内置资源：

```bash
curl http://localhost:8080/api \
  -H 'Content-Type: application/json' \
  -d '{
    "resource": "security/auth",
    "action": "login",
    "version": "v1",
    "params": {
      "type": "password",
      "principal": "admin",
      "credentials": "ChangeMe_123"
    }
  }'
```

在没有注册 challenge provider 的情况下，响应直接携带 token 对：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "tokens": {
      "accessToken": "eyJhbGciOiJIUzI1NiIs...",
      "refreshToken": "eyJhbGciOiJIUzI1NiIs..."
    }
  }
}
```

访问 token 在 30 分钟后过期（框架内固定常量）；刷新 token 的有效期来自 `vef.security.token_expires`（默认 7 天）。用刷新 token 换取新的 token 对——注意 `refresh` 的 `data` 直接就是 token 对，没有外层的 `tokens` 包装：

```bash
curl http://localhost:8080/api \
  -H 'Content-Type: application/json' \
  -d '{
    "resource": "security/auth",
    "action": "refresh",
    "version": "v1",
    "params": { "refreshToken": "eyJhbGciOiJIUzI1NiIs..." }
  }'
```

每个 `security/auth` action 的请求参数在[内置资源](../reference/built-in-resources)中有完整表格；响应字段——包括 challenge 包络与 `get_user_info` 的 `UserInfo` 形状——见 [RPC Resource: `security/auth`](./authentication-reference#rpc-resource-securityauth)。

## 实践建议

- `Public` 只用于明确需要匿名访问的操作
- 普通用户认证优先保持在 Bearer
- Signature 更适合系统对系统集成，而不是替代普通用户会话

## 下一步

- [认证参考](./authentication-reference) — 本指南背后的完整公开 API 面
- [授权](./authorization) — 认证之后权限检查如何继续发生
