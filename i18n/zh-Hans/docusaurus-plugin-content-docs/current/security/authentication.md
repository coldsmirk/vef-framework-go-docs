---
sidebar_position: 1
---

# 认证

VEF 的认证发生在 API 操作层。每个操作都有自己的 auth 配置，API 中间件会在 handler 执行前先解析出当前 principal。

Security grouped-family audit 固定了 174 grouped security field/method
entries，覆盖 66 receiver/type families：其中 75 public field entries、99
public method entries。这些 entries 覆盖 auth DTO wire fields、principal/token
helpers、signature 和 challenge providers、data-scope methods、
permission/resolver interfaces 以及 event fields；verifier 会锁定排序后的签名和
receiver/type 分布。

## 默认行为

如果你不做额外配置：

- 操作默认使用 Bearer 认证
- 显式 `Public` 的操作则不要求认证

这个默认值来自 API 引擎，而不是来自你的应用配置文件。

## 内置认证策略

公开的 `api` 包暴露了三种策略 helper：

- `api.Public()`
- `api.BearerAuth()`
- `api.SignatureAuth()`

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
| `refresh` | 是 | `vef.security.refresh_rate_limit` | `refreshToken`；`validate:"required"` |
| `logout` | 否 | 默认 API rate limit | 无 |
| `resolve_challenge` | 是 | `vef.security.login_rate_limit` | `challengeToken`、`type`、`response`；全部是 `validate:"required"` |
| `get_user_info` | 否 | 默认 API rate limit | 任意 `params`，会转发给 `UserInfoLoader.LoadUserInfo(...)` |

内置 authenticator type 字符串是 `password`、`token`、`refresh` 和
`signature`。普通客户端调用里，`security/auth.login` 使用
`type: "password"` 搭配用户名和密码凭证。Bearer 保护的操作会在内部使用
`token` authenticator，`security/auth.refresh` 会在内部使用 `refresh`，
`SignatureAuth` 会把签名 headers 映射到 `signature` authenticator。

`logout` 会立即返回 ok 结果。它不会在服务端吊销或拉黑 token；客户端需要自行删除已保存的
token，如果应用需要服务端吊销策略，需要自己扩展这部分逻辑。

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
| `ChallengeState` | Go-only state: `Principal`、`Pending`、`Resolved` |

## 应用通常还需要提供什么

这里要分场景来看：

- `security.UserLoader` 通常是用户登录和 refresh 流程的前提
- `security.ExternalAppLoader` 只在你使用签名认证的外部应用场景时需要
- challenge provider 是可选项，只有在你启用了挑战式登录流时才相关
- `security.UserInfoLoader` 只在你希望 `security/auth.get_user_info` 返回应用自定义用户信息时需要

框架提供的是认证流程和中间件，而不是你的身份源本身。

## 认证相关公开 Security API

| API 组 | 公开 surface |
| --- | --- |
| principal | `Principal`, `PrincipalType`, `NewUser`, `NewExternalApp`, `PrincipalSystem`, `PrincipalAnonymous`, `SetUserDetailsType`, `SetExternalAppDetailsType` |
| JWT | `JWT`, `JWTConfig`, `JWTClaimsBuilder`, `JWTClaimsAccessor`, `NewJWT`, `GenerateSecret`, token type constants, `DefaultJWTAudience`, `DefaultJWTSecret`, `JWTIssuer` |
| auth manager | `Authentication`, `AuthTokens`, `Authenticator`, `AuthManager`, `TokenGenerator`, `UserLoader`, `ExternalAppLoader`, `ExternalAppConfig`, `PasswordDecryptor` |
| challenge token | `ChallengeProvider`, `ChallengeState`, `ChallengeTokenStore`, `NewMemoryChallengeTokenStore`, `NewRedisChallengeTokenStore`, `NewJWTChallengeTokenStore` |
| OTP/challenge | `OTPEvaluator`, `OTPCodeSender`, `OTPCodeVerifier`, `OTPCodeStore`, `NewOTPChallengeProvider`, `NewDeliveredCodeSender`, `NewDeliveredCodeVerifier`, `NewDeliveredChallengeProvider`, `NewSMSChallengeProvider`, `NewEmailChallengeProvider` |
| TOTP/password/department | `NewTOTPEvaluator`, `NewTOTPVerifier`, `NewTOTPChallengeProvider`, `WithTOTPDestination`, `NewPasswordChangeChallengeProvider`, `NewDepartmentSelectionChallengeProvider` |
| signature auth | `Signature`, `SignatureCredentials`, `SignatureResult`, `SignatureAlgorithm`, `NewSignature`, `WithAlgorithm`, `WithTimestampTolerance`, `WithNonceStore`, `NonceStore`, `NewMemoryNonceStore`, `NewRedisNonceStore` |
| login event | `LoginEvent`, `LoginEventParams`, `NewLoginEvent`, `SubscribeLoginEvent` |

Bearer 相关常量是 `AuthSchemeBearer` 和 `QueryKeyAccessToken`。token type
常量是 `TokenTypeAccess`、`TokenTypeRefresh`、`TokenTypeChallenge`。

### JWT 与 principal

`NewJWT` 要求 `JWTConfig.Secret` 是十六进制 key，并会在 audience 为空时使用
`DefaultJWTAudience`。低层 `NewJWT` 在 secret 为空时仍会回退到公开的
`DefaultJWTSecret`；框架 security 模块在启动图里包了一层更安全的行为：生成进程内临时
key 并输出警告。生产环境应使用 `GenerateSecret()` 生成私有 key，再写入
`vef.security.secret`。

框架内置 token generator 签发的 access token 固定 `30m` 过期。
`vef.security.token_expires` 配置的是 refresh token 生命周期（默认 `168h`），
`vef.security.refresh_not_before` 默认是 `15m`。
同一次生成的 access token 和 refresh token 会共享同一个 `jti`。

JWT 解析只接受 `HS256`，要求 issuer 为 `JWTIssuer`（`vef`），会校验
audience，要求 `iat` 和 `exp`，并使用 10 秒 leeway。紧凑 claim key 如下：

| Claim | Key |
| --- | --- |
| JWT ID | `jti` |
| subject | `sub` |
| issuer | `iss` |
| audience | `aud` |
| issued at | `iat` |
| not before | `nbf` |
| expires at | `exp` |
| token type | `typ` |
| roles | `rls` |
| details | `det` |

内置 access token 和 refresh token generator 会把 subject 写成 `id@name`。
`JWTTokenAuthenticator` 会直接从这个 subject 重建用户 principal，不查数据库。
`JWTRefreshAuthenticator` 也要求 `id@name`，但会取其中的 `id` 调用
`UserLoader.LoadByID(...)` 重新加载用户。

`JWTClaimsBuilder` 用 `WithID`、`WithSubject`、`WithRoles`、`WithDetails`、
`WithType`、`WithClaim` 写入紧凑 token claims。`JWTClaimsAccessor` 用
`ID`、`Subject`、`Roles`、`Details`、`Type`、`Claim` 读取同一份 payload。
可以用 `NewJWTClaimsBuilder()` 和 `NewJWTClaimsAccessor(...)` 直接创建这两个 helper。

`PrincipalTypeUser`、`PrincipalTypeExternalApp`、`PrincipalTypeSystem` 描述
支持的 principal 类型。`SetUserDetailsType[T]()` 和
`SetExternalAppDetailsType[T]()` 会配置进程级 details 反序列化目标，应在启动阶段调用，
不要等服务开始处理请求后再修改。

`Principal` 的 JSON 字段是 `type`、`id`、`name`、`roles` 和 `details`。
`SetUserDetailsType[T]()` 与 `SetExternalAppDetailsType[T]()` 要求 `T` 是 struct
或 struct pointer；否则会分别以 `ErrUserDetailsNotStruct` 或
`ErrExternalAppDetailsNotStruct` panic。它们会修改 package-level 状态，应视为启动期配置。
未知 principal type 会把 `details` 保留为 `map[string]any`；system principal
反序列化后 `details` 为 `nil`。内置特殊 principal 是 `PrincipalSystem`
（`type: "system"`，id `system`，name `系统`）和 `PrincipalAnonymous`
（`type: "user"`，id `anonymous`，name `匿名`）。

### Challenge providers

内置 challenge type 常量包括：

- `ChallengeTypeTOTP`
- `ChallengeTypeSMS`
- `ChallengeTypeEmail`
- `ChallengeTypePasswordChange`
- `ChallengeTypeDepartmentSelection`

它们的 wire value 和默认顺序是：

| Constant | Wire value | Default order |
| --- | --- | --- |
| `ChallengeTypeTOTP` | `totp` | `100` |
| `ChallengeTypeSMS` | `sms_otp` | `200` |
| `ChallengeTypeEmail` | `email_otp` | `300` |
| `ChallengeTypePasswordChange` | `password_change` | `400` |
| `ChallengeTypeDepartmentSelection` | `department_selection` | `500` |

`ChallengeTokenStore.Generate(ctx, principal, pending, resolved)` 和
`Parse(ctx, token)` 负责在 `login` 与 `resolve_challenge` 之间携带状态。
内置登录资源把这个状态字段暴露为 `challengeToken`。
`JWTChallengeTokenStore` 是无状态实现；`MemoryChallengeTokenStore` 适合测试或单实例；
`RedisChallengeTokenStore` 适合分布式部署。challenge token 的有效期是
`ChallengeTokenExpires`。JWT-backed store 使用 `ClaimChallengePrincipalType`、
`ClaimChallengePrincipalName`、`ClaimChallengePending`、`ClaimChallengeResolved`
作为紧凑 claim key。

challenge token store 的 wire/storage 形状不同：

| Store | Token/state contract |
| --- | --- |
| `JWTChallengeTokenStore` | JWT token，`typ: "challenge"`，5 分钟 `ChallengeTokenExpires`，subject 只保存 principal ID |
| `MemoryChallengeTokenStore` | UUID token，按 `ChallengeTokenExpires` 存在进程内存里 |
| `RedisChallengeTokenStore` | UUID token，按 `ChallengeTokenExpires` 存在 `vef:security:challenge:<token>` |

JWT challenge claim key 是 `ptp`（`ClaimChallengePrincipalType`）、`pnm`
（`ClaimChallengePrincipalName`）、`pnd`（`ClaimChallengePending`）和 `rsd`
（`ClaimChallengeResolved`）。challenge 解析会把空 principal type 作为向后兼容的
user principal，也接受 `user`、`external_app`、`system`，未知 principal type 会被拒绝。

challenge provider 会按 `Order()` 升序排序。内置 convenience provider 的顺序是：
TOTP 为 `100`，SMS 为 `200`，email 为 `300`，password change 为 `400`，
department selection 为 `500`。未注册的 provider，或者 `Evaluate(...)`
返回 `nil` 的 provider，会被跳过。执行 `resolve_challenge` 时，提交的
`type` 必须等于第一个 pending challenge type，否则框架返回
`ErrChallengeTypeInvalid`。

`NewOTPChallengeProvider` 是通用构造器。`OTPChallengeProviderConfig` 要求
`ChallengeType`、`Evaluator`、`Verifier`；`ChallengeOrder` 控制评估顺序，
`Sender` 可选，用于 delivered-code 流程。
`OTPChallengeProvider` 在需要 challenge 时会把 `OTPChallengeData` 返回给客户端。
delivered-code helper 会把 `OTPCodeStore` 和 `OTPCodeDelivery` 组合起来：
`DeliveredCodeSender`、`DeliveredCodeVerifier`、`NewDeliveredCodeSender`、
`NewDeliveredCodeVerifier`、`NewDeliveredChallengeProvider`、
`NewSMSChallengeProvider`、`NewEmailChallengeProvider`。

`NewTOTPChallengeProvider` 只需要 `TOTPSecretLoader`；如果 `LoadSecret(...)`
返回空字符串，就跳过 challenge。`TOTPEvaluator`、`TOTPVerifier`、`TOTPOption`
是 convenience provider 背后的低层组件。TOTP 默认显示 `TOTPDefaultDestination`，
也就是 `Authenticator App`；也可以用 `WithTOTPDestination(...)` 覆盖。

`NewPasswordChangeChallengeProvider` 使用 `PasswordChangeChecker` 和
`PasswordChanger`；需要强制改密时会返回 `PasswordChangeChallengeData`。常见 reason
常量是 `PasswordChangeReasonFirstLogin`（`first_login`）和
`PasswordChangeReasonExpired`（`expired`）。具体 provider
类型是 `PasswordChangeChallengeProvider`。`NewDepartmentSelectionChallengeProvider` 使用
`DepartmentLoader` 与 `DepartmentSelector`；空部门列表会跳过 challenge，resolve 时要求
传入 department ID 字符串。`DepartmentSelectionChallengeData` 会序列化为
`departments` 和可选 `meta`；每个 `DepartmentOption` 会序列化为 `id` 和 `name`。

这些 challenge 构造器属于 wiring-time API。`NewOTPChallengeProvider` 在缺少
`ChallengeType`、`Evaluator` 或 `Verifier` 时会 panic。
`NewPasswordChangeChallengeProvider` 在缺少 `PasswordChangeChecker` 或
`PasswordChanger` 时会 panic。`NewDepartmentSelectionChallengeProvider` 在缺少
`DepartmentLoader` 或 `DepartmentSelector` 时会 panic。

### Signature helpers

`NewSignature(secret, ...)` 要求非空十六进制 secret，默认使用
`SignatureAlgHmacSHA256` 和 5 分钟 timestamp tolerance。option 类型是
`SignatureOption`。其他 algorithm 常量是 `SignatureAlgHmacSHA512` 与
`SignatureAlgHmacSM3`。`WithTimestampTolerance` 会调整接受的时间窗口，
`WithNonceStore` 控制 replay protection。低层 `NewSignature(...)` 默认会创建
`MemoryNonceStore`；只有在你明确想关闭这个 helper 的 nonce 存储时，才传入
`WithNonceStore(nil)`。内置 `SignatureAuthenticator` 会在应用提供
`NonceStore` 时注入该 store；否则每次校验使用低层 helper 的进程本地 memory
store。`MemoryNonceStore` 只适合单进程；`RedisNonceStore` 是分布式选项。
nonce 存储 TTL 是 timestamp tolerance 再加 1 分钟 buffer。

签名 payload 精确为：

```text
app_id=<appID>&method=<method>&nonce=<nonce>&path=<path>&timestamp=<timestamp>
```

字段按这个顺序绑定。`request body is not part of the signature payload`。
其中 `method` 字段是服务端看到的 HTTP method。

`NewIPWhitelistValidator` 会从逗号分隔的 IP 和 CIDR 列表创建
`IPWhitelistValidator`。空 whitelist 表示允许全部 IP；非法 whitelist 会 fail-closed，
拒绝全部请求。当 `ExternalAppConfig.IPWhitelist` 非空但请求 IP 无法解析时，
`SignatureAuthenticator` 也会 fail closed，并返回 `ErrIPNotAllowed`。

Signature 相关存储 key 和默认值：

| Contract | Value |
| --- | --- |
| request headers | `X-App-ID`, `X-Timestamp`, `X-Nonce`, `X-Signature` |
| algorithms | `HMAC-SHA256`, `HMAC-SHA512`, `HMAC-SM3` |
| default algorithm | `HMAC-SHA256` |
| default tolerance | `5m` |
| nonce TTL | tolerance + `1m` |
| Redis nonce prefix | `vef:security:nonce:` |
| disable replay checking | `WithNonceStore(nil)` |

security-domain API 错误暴露 `1000-1039` 范围内的 `ErrCode*` 常量：

| Code | Constant | Error | HTTP status |
| --- | --- | --- | --- |
| `1000` | `ErrCodeUnauthenticated` | `ErrUnauthenticated` | `401` |
| `1001` | `ErrCodeUnsupportedAuthenticationType` | unsupported authentication type | `400` |
| `1002` | `ErrCodeTokenExpired` | `ErrTokenExpired` | `401` |
| `1003` | `ErrCodeTokenInvalid` | `ErrTokenInvalid` | `401` |
| `1004` | `ErrCodeTokenNotValidYet` | `ErrTokenNotValidYet` | `401` |
| `1005` | `ErrCodeTokenInvalidIssuer` | `ErrTokenInvalidIssuer` | `401` |
| `1006` | `ErrCodeTokenInvalidAudience` | `ErrTokenInvalidAudience` | `401` |
| `1007` | `ErrCodePrincipalInvalid` | `ErrPrincipalInvalid(message)` | `401` |
| `1008` | `ErrCodeCredentialsInvalid` | `ErrCredentialsInvalid(message)` | `401` |
| `1009` | `ErrCodeAppIDRequired` | `ErrAppIDRequired` | `401` |
| `1010` | `ErrCodeTimestampRequired` | `ErrTimestampRequired` | `401` |
| `1011` | `ErrCodeSignatureRequired` | `ErrSignatureRequired` | `401` |
| `1012` | `ErrCodeTimestampInvalid` | `ErrTimestampInvalid` | `401` |
| `1013` | `ErrCodeSignatureExpired` | `ErrSignatureExpired` | `401` |
| `1014` | `ErrCodeExternalAppNotFound` | `ErrExternalAppNotFound` | `401` |
| `1015` | `ErrCodeExternalAppDisabled` | `ErrExternalAppDisabled` | `401` |
| `1016` | `ErrCodeIPNotAllowed` | `ErrIPNotAllowed` | `401` |
| `1017` | `ErrCodeSignatureInvalid` | `ErrSignatureInvalid` | `401` |
| `1018` | `ErrCodeNonceRequired` | `ErrNonceRequired` | `401` |
| `1019` | `ErrCodeNonceInvalid` | `ErrNonceInvalid` | `401` |
| `1020` | `ErrCodeNonceAlreadyUsed` | `ErrNonceAlreadyUsed` | `401` |
| `1021` | `ErrCodeAuthHeaderMissing` | `ErrAuthHeaderMissing` | `401` |
| `1022` | `ErrCodeAuthHeaderInvalid` | `ErrAuthHeaderInvalid` | `401` |
| `1031` | `ErrCodeChallengeTokenInvalid` | `ErrChallengeTokenInvalid` | `401` |
| `1033` | `ErrCodeChallengeTypeInvalid` | `ErrChallengeTypeInvalid` | `400` |
| `1034` | `ErrCodeChallengeResolveFailed` | challenge resolve failure message ID | reserved |
| `1035` | `ErrCodeOTPCodeRequired` | `ErrOTPCodeRequired` | `400` |
| `1036` | `ErrCodeOTPCodeInvalid` | `ErrOTPCodeInvalid` | `401` |
| `1037` | `ErrCodeNewPasswordRequired` | `ErrNewPasswordRequired` | `400` |
| `1038` | `ErrCodeDepartmentRequired` | `ErrDepartmentRequired` | `400` |

认证相关 sentinel 包括 `ErrUnauthenticated`、`ErrTokenExpired`、
`ErrTokenInvalid`、`ErrTokenNotValidYet`、`ErrTokenInvalidIssuer`、
`ErrTokenInvalidAudience`、`ErrAppIDRequired`、`ErrTimestampRequired`、
`ErrSignatureRequired`、`ErrTimestampInvalid`、`ErrSignatureExpired`、
`ErrSignatureInvalid`、`ErrExternalAppNotFound`、`ErrExternalAppDisabled`、
`ErrIPNotAllowed`、`ErrNonceRequired`、`ErrNonceInvalid`、
`ErrNonceAlreadyUsed`、`ErrAuthHeaderMissing`、`ErrAuthHeaderInvalid`、
`ErrChallengeTokenInvalid`、`ErrChallengeTypeInvalid`、
`ErrOTPCodeRequired`、`ErrOTPCodeInvalid`、`ErrNewPasswordRequired`、
`ErrDepartmentRequired`，以及 factory helper `ErrCredentialsInvalid(message)`
和 `ErrPrincipalInvalid(message)`。`ErrCodeChallengeResolveFailed` 保留给
challenge resolve 失败场景。

低层 secret 解析错误使用 `ErrDecodeJWTSecretFailed`、
`ErrGenerateJWTSecretFailed`、`ErrDecodeSignatureSecretFailed` 和
`ErrSignatureSecretRequired`。details 类型注册错误使用 `ErrUserDetailsNotStruct`
与 `ErrExternalAppDetailsNotStruct`。
公开的 i18n message ID 常量包括 `ErrMessageChallengeResolveFailed`、
`ErrMessageCredentialsFormatInvalid`、`ErrMessageExternalAppLoaderNotImplemented`、
`ErrMessageUnauthenticated`、`ErrMessageUnsupportedAuthenticationType`、
`ErrMessageUserInfoLoaderNotImplemented` 和 `ErrMessageUserLoaderNotImplemented`。

## 真实项目里常见的 auth 模块形态

很多 VEF 项目里的 auth 模块其实都很小，常见模式是：

- 在 `init()` 里设置用户详情类型
- 提供 `UserLoader`
- 提供 `UserInfoLoader`

概念上通常像这样：

```go
func init() {
  security.SetUserDetailsType[*UserDetails]()
}

var Module = vef.Module(
  "app:auth",
  vef.Provide(
    NewUserLoader,
    NewUserInfoLoader,
  ),
)
```

这样认证接入代码就能和业务资源模块保持分离。

## 实践建议

- `Public` 只用于明确需要匿名访问的操作
- 普通用户认证优先保持在 Bearer
- Signature 更适合系统对系统集成，而不是替代普通用户会话

## 下一步

继续阅读 [授权](./authorization)，看认证之后权限检查是如何继续发生的。
