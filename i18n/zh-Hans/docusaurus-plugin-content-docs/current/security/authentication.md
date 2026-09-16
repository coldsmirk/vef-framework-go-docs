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
- `api.IPAuth(...)`（白名单解析方式见下文的 [Signature helpers](./authentication-reference#signature-helpers)）；内置策略是 fail-closed 的：缺失、空白或无法解析的白名单会一律以 `security.ErrIPNotAllowed`（HTTP 401）拒绝请求。“空名单允许全部 IP” 属于签名外部应用使用的底层 `IPWhitelistValidator` 的行为，而不是 `api.IPAuth(...)` 的行为
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
`security.Principal.IsReserved()` 报告以下情况：

- `system` 主体类型（`PrincipalTypeSystem`）；
- 主体 `ID` 等于 `orm.OperatorSystem`（`"system"`）时；
- 主体 `ID` 等于 `orm.OperatorCronJob`（`"cron_job"`）时。

`PrincipalAnonymous` 故意**不**属于保留身份：它表示身份缺失，只有 `public`
认证策略才会合法地产出。

框架在每个边界 fail-closed 地强制该不变量：

- **认证边界**：API 认证中间件拒绝任何返回 nil 或保留主体的策略，返回
  `security.ErrReservedPrincipal`（业务码 `1007`，HTTP 401）。
- **令牌签发**：`JWTTokenGenerator` 和 `OpaqueTokenGenerator` 均拒绝为保留主体
  签发令牌。
- **挑战流程**：挑战令牌解析会把保留主体类型和保留 ID 视为
  `ErrTokenInvalid`；挑战提供者解析后，`resolve_challenge` 再拒绝保留结果，
  返回 `ErrReservedPrincipal`。挑战令牌的 subject 只保存 principal ID，不包含用户名、部门等额外信息。

内置密码登录还会额外拒绝以保留标识符（`system`、`cron_job`、`anonymous`）
作为登录 `principal`。自定义 `UserLoader`、`APIKeyLoader` 等身份源绝不能返回
ID 与保留操作者 ID 冲突的主体。

保留身份的拒绝会**记录审计事件，但不计入登录锁定计数**：凭据本身可能是对的，
错误在于认证器或挑战提供者，而非调用方。锁定行为参见
[登录加固](./login-hardening)。

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

`logout` 会立即返回 ok 结果。在 `jwt_token` 下它实际上是 no-op——不会在服务端吊销或拉黑 token，客户端需要自行删除已保存的 token。在 `opaque_token` 下它会
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
| `LoginContext` | Go-only，只读：`AuthType`、`Username`、`Principal`、`Resolved`——挑战所处的那次登录（见[登录上下文](#登录上下文)） |
| `ChallengeState` | challenge token 携带的 Go-only 状态：内嵌的 `LoginContext`，加上 `Pending`——尚未解决的挑战类型，按评估顺序排列（第一个即当前呈现的挑战） |

两种响应形态——token 载荷与 challenge 包络——的逐字段表格及 JSON 示例见
[RPC Resource: `security/auth`](./authentication-reference#rpc-resource-securityauth)。

### 登录上下文

挑战总是发生在某一次登录之中，provider 拿到的就是这次登录，类型为
`*security.LoginContext`：

| 字段 | 内容 |
| --- | --- |
| `AuthType` | 认证该 principal 的登录方式——即提交给 `login` 的 `type`：`security.AuthTypePassword`（`password`）、`security.AuthTypeTrustCode`（`trust_code`），或宿主 `security.Authenticator` 支持的类型。同一次登录的每一步都相同 |
| `Username` | 以 `principal` 提交给 `login` 的标识，挑战之后发布的事件因此报告的仍是最初提交的标识 |
| `Principal` | 经已解决的挑战充实后的身份 |
| `Resolved` | 迄今已解决的挑战类型，按解决顺序排列 |

`ChallengeProvider.Evaluate(ctx, login)` 判断自己的挑战是否适用于这次登录，
不适用时返回 nil；`Resolve(ctx, login, response)` 校验应答，并返回登录继续
使用的 principal——`login.Principal`，或者像部门选择那样返回一个充实后的副本。
challenge token 在每一步 `resolve_challenge` 之间携带这份上下文，因此后续步骤
看到的 `AuthType` 与 `Username` 和第一步完全相同。上下文归框架所有、以指针
传入：请把它当作只读，想改变身份只能通过 `Resolve` 返回一个 principal。

内置 provider 背后的钩子按同一条规则划分：决定挑战**是否**适用的钩子看到整次
登录；**作用于**身份的钩子只看到 principal——用户以何种方式登录，与密码如何
存储、验证码如何校验毫无关系。

| 决定是否适用（接收 `login`） | 作用于身份（接收 `principal`） |
| --- | --- |
| `PasswordChangeChecker.Check`——包括 `ExpiryPasswordChangeChecker` 与 `NewCompositePasswordChangeChecker` | `PasswordChanger`、`PasswordValidator`、`PasswordMetadataLoader` |
| `OTPEvaluator.Evaluate`——包括 `TOTPEvaluator` | `OTPCodeSender`、`OTPCodeVerifier`、`OTPCodeStore`、`OTPCodeDelivery`、`TOTPSecretLoader` |
| `DepartmentLoader.LoadDepartments` | `DepartmentSelector` |

决定类钩子从 `login.Principal` 读取用户。针对单个用户的条件——用户是否配置了
TOTP 密钥、密码是否已过期——属于这些钩子；而一个挑战属于哪些登录方式，更适合
在注册时一次性声明。各方法签名见[认证参考](./authentication-reference#challenge-providers)。

### 按登录方式限定挑战

provider 默认适用于每一次登录，除非它自己另作判断。要在不改动 provider 的前提
下把它限定到部分登录方式，在注册时用
`security.NewFilteredChallengeProvider(provider, filters...)` 包装它。每个
`security.LoginFilter` 以数据形式声明适用范围：

| 字段 | 构造函数 | 匹配条件 |
| --- | --- | --- |
| `AuthTypes` | `security.ForAuthTypes(types...)` | 登录的 `AuthType` 在列表中——允许列表（allow-list） |
| `ExcludedAuthTypes` | `security.ExceptAuthTypes(types...)` | 登录的 `AuthType` 不在列表中——排除列表（deny-list） |

- 在单个 filter 内，空维度不做约束，因此 `LoginFilter{}` 匹配任何登录；两个
  维度都填写的 filter 要求两者同时满足。
- 传给同一个 provider 的多个 filter 必须全部匹配（AND）。
- 对被 filter 拒绝的登录，provider 被跳过，效果与其 `Evaluate` 返回 nil 完全
  相同，挑战链继续评估下一个 provider。
- 不传 filter 时原样返回该 provider。`LoginFilter.Matches(login)` 就是判定本身。

举个例子：某应用支持密码登录和宿主自定义的微信小程序登录——`type: "wechat_mini"`，
由它自己的认证器处理——同时接受信任登录交接。强制改密只属于密码登录。TOTP 第二
因子属于除交接之外的每一种登录，因为这个应用信任发起交接的系统已经执行过自己的
第二因素：

```go
var Module = vef.Module(
	"app:auth",
	vef.ProvideAuthenticator(NewMiniProgramAuthenticator), // Supports("wechat_mini")
	vef.ProvideChallengeProvider(func(
		checker security.PasswordChangeChecker,
		changer security.PasswordChanger,
		validator security.PasswordValidator,
	) security.ChallengeProvider {
		return security.NewFilteredChallengeProvider(
			security.NewPasswordChangeChallengeProvider(checker, changer, validator),
			security.ForAuthTypes(security.AuthTypePassword),
		)
	}),
	vef.ProvideChallengeProvider(func(loader security.TOTPSecretLoader) security.ChallengeProvider {
		return security.NewFilteredChallengeProvider(
			security.NewTOTPChallengeProvider(loader),
			security.ExceptAuthTypes(security.AuthTypeTrustCode),
		)
	}),
)
```

每个构造函数都返回 `security.ChallengeProvider`，也就是该 provider group 声明
的类型：fx 按类型区分 group，返回具体 provider 类型的构造函数——
`security.NewTOTPChallengeProvider` 本身就是如此——会被静默丢弃，不报任何错误。
checker、changer 与 TOTP 密钥加载器由应用提供；`PasswordValidator` 是框架基于
`vef.security.password_policy` 构建的。每种登录遇到的挑战链：

| 登录 `type` | `totp`（order `100`） | `password_change`（order `400`） |
| --- | --- | --- |
| `password` | 评估 | 评估 |
| `wechat_mini` | 评估 | 跳过 |
| `trust_code` | 跳过 | 跳过 |

“评估”表示交给 provider 自己的钩子决定：没有配置密钥的用户仍会跳过 TOTP，无需
改密的用户也不会遇到改密挑战。

:::caution[按挑战所守护的对象选择列表]
与某一种凭据绑定的挑战用允许列表：强制改密针对的是密码，
`ForAuthTypes(security.AuthTypePassword)` 让它远离从未提交过密码的登录。第二
因子用排除列表：只用 `ExceptAuthTypes(...)` 豁免那些已经具备同等保障的登录方式。
如果把第二因素写成允许列表，之后新增的每一种登录方式都会被悄悄豁免；写成排除
列表，新的登录方式在有人明确决定之前都会被挑战。
:::

部门选择是必需的业务输入，通常不加 filter。信任登录的交接如何经过挑战链，见
[信任登录](./trust-login)。

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
