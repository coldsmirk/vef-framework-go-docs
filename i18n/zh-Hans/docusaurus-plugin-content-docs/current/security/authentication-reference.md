---
sidebar_position: 2
---

# 认证参考

`security` 包中认证相关的公开接口面：principal、JWT、认证管理器、挑战提供者与令牌存储、签名认证与登录事件。叙述性指南——认证策略、内置认证资源与登录流程——见[认证](./authentication)。内置认证端点的 wire 层契约——每个 action 的请求与响应字段——收录在本页末尾的 [RPC Resource: `security/auth`](#rpc-resource-securityauth)。

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

## JWT 与 principal

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
audience，校验存在时的 `iat`、强制要求 `exp`，并使用 10 秒 leeway。紧凑 claim key 如下：

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

## Challenge providers

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

`ChallengeTokenStore.Generate(ctx, principal, username, pending, resolved)` 和
`Parse(ctx, token)` 负责在 `login` 与 `resolve_challenge` 之间携带状态
（`username` 是申请人在第一步提交的原始登录标识，跨挑战步骤保留、用于审计
事件）。
内置登录资源把这个状态字段暴露为 `challengeToken`。
`JWTChallengeTokenStore` 是无状态实现；`MemoryChallengeTokenStore` 适合测试或单实例；
`RedisChallengeTokenStore` 适合分布式部署。challenge token 的有效期是
`ChallengeTokenExpires`。JWT-backed store 使用 `ClaimChallengePrincipalType`、
`ClaimChallengePrincipalName`、`ClaimChallengeUsername`、`ClaimChallengePending`、`ClaimChallengeResolved`
作为紧凑 claim key。

challenge token store 的 wire/storage 形状不同：

| Store | Token/state contract |
| --- | --- |
| `JWTChallengeTokenStore` | JWT token，`typ: "challenge"`，5 分钟 `ChallengeTokenExpires`，subject 只保存 principal ID |
| `MemoryChallengeTokenStore` | UUID token，按 `ChallengeTokenExpires` 存在进程内存里 |
| `RedisChallengeTokenStore` | UUID token，按 `ChallengeTokenExpires` 存在 `vef:security:challenge:<token>` |

JWT challenge claim key 是 `ptp`（`ClaimChallengePrincipalType`）、`pnm`
（`ClaimChallengePrincipalName`）、`unm`（`ClaimChallengeUsername`）、`pnd`
（`ClaimChallengePending`）和 `rsd`（`ClaimChallengeResolved`）。在
保留身份加固下，challenge 解析只接受 `user` 与 `external_app` 两种
principal type——`system`、空值与未知类型一律以 `ErrTokenInvalid` 拒绝
（携带框架内部身份的挑战 token 不可能有合法来源），解析出的 principal 若
`IsReserved()` 同样被拒绝。

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

## Signature helpers

`NewSignature(secret, ...)` 要求非空十六进制 secret，默认使用
`SignatureAlgHmacSHA256` 和 5 分钟 timestamp tolerance。option 类型是
`SignatureOption`。其他 algorithm 常量是 `SignatureAlgHmacSHA512` 与
`SignatureAlgHmacSM3`。`WithTimestampTolerance` 会调整接受的时间窗口，
`WithNonceStore` 控制 replay protection。低层 `NewSignature(...)` 默认会创建
`MemoryNonceStore`；只有在你明确想关闭这个 helper 的 nonce 存储时，才传入
`WithNonceStore(nil)`。内置 `SignatureAuthenticator` 会在应用提供
`NonceStore` 时注入该 store；否则每次校验使用低层 helper 的进程本地 memory
store。`MemoryNonceStore` 只适合单进程；`RedisNonceStore` 是分布式选项。
nonce 存储 TTL 是 2 倍 timestamp tolerance 再加 1 分钟 buffer。

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

`security.NewIPWhitelistValidatorFromEntries(entries)` 是基于 slice 的构造器，
内置 `api.IPAuth(...)` strategy 使用它。该 strategy 会通过
`security.IPWhitelistLoader` 解析命名的 `security.IPWhitelist`；默认 loader 读取
`vef.security.ip_whitelists`，应用也可以提供自己的 loader，从数据库或配置中心加载
名单。

Signature 相关存储 key 和默认值：

| Contract | Value |
| --- | --- |
| request headers | `X-App-ID`, `X-Timestamp`, `X-Nonce`, `X-Signature` |
| algorithms | `HMAC-SHA256`, `HMAC-SHA512`, `HMAC-SM3` |
| default algorithm | `HMAC-SHA256` |
| default tolerance | `5m` |
| nonce TTL | `2*tolerance + 1m` |
| Redis nonce prefix | `vef:security:nonce:` |
| disable replay checking | `WithNonceStore(nil)` |

security-domain API 错误暴露的 `ErrCode*` 常量按段划分：`1000`–`1029` 为
认证、`1030`–`1039` 为挑战、`1050` 为密码策略（所有策略违规共享这一个
code——见[登录加固](./login-hardening)）：

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
| `1023` | `ErrCodeAccountLocked` | 动态账号锁定错误（见[登录加固](./login-hardening)） | `429` |
| `1024` | `ErrCodeTooManyConcurrentSessions` | `ErrTooManyConcurrentSessions` | `403` |
| `1025` | `ErrCodeAPIKeyInvalid` | `ErrAPIKeyInvalid` | `401` |
| `1026` | `ErrCodeBasicCredentialsInvalid` | `ErrBasicCredentialsInvalid` | `401` |
| `1031` | `ErrCodeChallengeTokenInvalid` | `ErrChallengeTokenInvalid` | `401` |
| `1033` | `ErrCodeChallengeTypeInvalid` | `ErrChallengeTypeInvalid` | `400` |
| `1034` | `ErrCodeChallengeResolveFailed` | `ErrChallengeResolveFailed` | `401` |
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
`ErrChallengeResolveFailed`、`ErrOTPCodeRequired`、`ErrOTPCodeInvalid`、
`ErrNewPasswordRequired`、`ErrDepartmentRequired`、
`ErrTooManyConcurrentSessions`、`ErrAPIKeyInvalid`、
`ErrBasicCredentialsInvalid`，以及 `ErrReservedPrincipal`（在认证、
挑战解析与令牌签发每个入口拒绝框架内部身份；复用
`ErrCodePrincipalInvalid`/`1007`，HTTP 401），另有 factory helper
`ErrCredentialsInvalid(message)` 和 `ErrPrincipalInvalid(message)`。
`ErrChallengeResolveFailed` 不是保留占位：`resolve_challenge`
会把 `ChallengeProvider` 返回的裸 error 归一化为它。

低层 secret 解析错误使用 `ErrDecodeJWTSecretFailed`、
`ErrGenerateJWTSecretFailed`、`ErrDecodeSignatureSecretFailed` 和
`ErrSignatureSecretRequired`。details 类型注册错误使用 `ErrUserDetailsNotStruct`
与 `ErrExternalAppDetailsNotStruct`。
公开的 i18n message ID 常量包括 `ErrMessageChallengeResolveFailed`、
`ErrMessageCredentialsFormatInvalid`、`ErrMessageExternalAppLoaderNotImplemented`、
`ErrMessageUnauthenticated`、`ErrMessageUnsupportedAuthenticationType`、
`ErrMessageUserInfoLoaderNotImplemented` 和 `ErrMessageUserLoaderNotImplemented`。

## RPC Resource: `security/auth`

安全模块把内置认证资源作为 RPC 资源挂载在 `/api` 下，使用标准信封
（`resource`、`action`、`version`、`params`）。响应遵循标准 result 信封——
`code`（成功为 `0`）、`message`、`data`——下文描述的都是 `data` 载荷的
形状。请求参数表也收录在[内置资源](../reference/built-in-resources)；
本节是完整的 wire 层契约，包含每个响应字段。

| Action | 访问性 | Rate limit（`max`） | 入参 | 出参（`data`） |
| --- | --- | --- | --- | --- |
| `login` | Public | `vef.security.login_rate_limit`（默认 `6`） | `LoginParams` | `LoginResult`——token **或** challenge 包络 |
| `refresh` | Public；仅在 `token_type = "jwt_token"` 下挂载 | `vef.security.refresh_rate_limit`（默认 `1`） | `RefreshParams` | `AuthTokens`（无 `tokens` 包装） |
| `logout` | Bearer 认证 | API 引擎默认 | 无 | 空（`data: null`） |
| `resolve_challenge` | Public | `vef.security.login_rate_limit`（默认 `6`） | `ResolveChallengeParams` | `LoginResult`——下一个 challenge **或** 最终 token |
| `get_user_info` | Bearer 认证 | API 引擎默认 | 原样 `params` map | `UserInfo` |

自定义限流只设置 `max`；时间窗口回退到 API 引擎的默认限流周期
（`vef.api.rate_limit`，默认 `5m`）。未声明自定义限流的操作完整继承引擎
默认值（出厂为每 `5m` `100` 次）。在 `token_type = "opaque_token"` 下，
`refresh` 操作根本不会挂载——调用它会得到操作不存在错误（HTTP 404），
因为 opaque 会话在使用中自行续期。

### `login`

`LoginParams`：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `type` | `string` | 是 | 凭证类型。框架为该端点内置的只有 `password`；注册自定义 `security.Authenticator` 可扩展取值。框架签发的令牌类型（`jwt_token`、`opaque_token`、`refresh`）会以 code `1001` 拒绝，已签发的令牌永远无法在此洗换成新的 token 对 |
| `principal` | `string` | 是 | 登录标识，通常是用户名。内置密码流程拒绝保留标识（`system`、`cron_job`、`anonymous`），返回 code `1007` |
| `credentials` | `any` | 是 | 凭证载荷。`type = "password"` 时是密码字符串——配置了 `security.PasswordDecryptor` 时为传输加密密文，否则为明文 |

响应是一个 `LoginResult`，形态严格二选一。`tokens`、`challengeToken`、
`challenge` 均为 `omitempty`：不适用的那一半直接缺失，绝不会是 `null`。

**形态一——token。** 未注册 challenge provider，或没有任何 provider 适用
于该账号。`data.tokens` 是一个 `AuthTokens`：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `tokens.accessToken` | `string` | 后续请求使用的 bearer token。`jwt_token` 下是固定 `30m` 有效期的 JWT；`opaque_token` 下是随机会话引用，有效期由[会话策略](./session-management)（`idle_ttl` / `max_lifetime`）决定 |
| `tokens.refreshToken` | `string` | JWT refresh token，有效期 `vef.security.token_expires`（默认 `168h`）。**`opaque_token` 下缺失**——会话自行续期，不存在 refresh token |

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

载荷中不携带任何过期时间字段——令牌有效期属于部署配置，需另行告知客户端。

**形态二——challenge 包络。** 凭证已通过校验，但至少一个
[challenge provider](#challenge-providers) 要求第二步。此时尚未签发任何
认证 token：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `challengeToken` | `string` | 携带挑战进度（principal、原始登录标识、pending 与 resolved 类型列表）的状态令牌——客户端将其视为不透明值，传给 `resolve_challenge` 即可。每个令牌在 `ChallengeTokenExpires`（`5m`）后过期；每个成功步骤都会签发新令牌 |
| `challenge` | `LoginChallenge` | 第一个待解的 challenge（见下） |

`LoginChallenge`：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `type` | `string` | challenge type 的 wire value，如 `totp`、`sms_otp`、`password_change`（见 [wire value 表](#challenge-providers)） |
| `data` | `any` | provider 特定的展示数据；provider 不提供时缺失。OTP provider 返回 `{destination, meta?}`（`OTPChallengeData`）；部门选择返回 `{departments, meta?}` |
| `required` | `bool` | 完成登录是否必须解决该 challenge |

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "challengeToken": "eyJhbGciOiJIUzI1NiIs...",
    "challenge": {
      "type": "totp",
      "data": { "destination": "Authenticator App" },
      "required": true
    }
  }
}
```

行为说明：

- provider 严格按 `Order()` 顺序评估；`Evaluate(...)` 返回 `nil` 的
  provider 被跳过，因此包络里始终是第一个真正适用的 challenge。
- 暴力破解 guard 在凭证通过校验的那一刻——早于任何第二因子——就清空失败
  计数。凭证被拒绝会发布失败 `LoginEvent`；成功事件只在 token 真正签发时
  发布——无挑战时立即发布，有挑战时在挑战链末尾发布——始终携带提交的
  登录标识。
- 典型失败（均见上文[错误码表](#signature-helpers)）：`1001`（不支持/被
  拒绝的 `type`，HTTP 400）、`1008`（凭证无效——未知用户、nil principal
  或空存储哈希、密码错误刻意返回同一响应，HTTP 401）、`1007`（保留或非
  法 principal，HTTP 401）、`1023`（触发暴力破解锁定，HTTP 429），以及
  缺少必填字段时的通用校验错误 `1400`（HTTP 400）。

### `refresh`

仅在无状态 JWT 机制（`token_type = "jwt_token"`，即默认值）下挂载。

`RefreshParams`：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `refreshToken` | `string` | 是 | 由 `login`、上一次 `refresh` 或签发 token 的 `resolve_challenge` 返回的 refresh token |

响应 `data` **直接就是**新的 `AuthTokens` 对——**没有 `login` 那层
`tokens` 包装**：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `accessToken` | `string` | 新 access token（`30m` 有效期） |
| `refreshToken` | `string` | 新 refresh token（有效期 `vef.security.token_expires`） |

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIs...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

行为说明：

- 内部按 `refresh` 类型认证：JWT 必须能解析、携带 `typ: "refresh"`
  （access token 会被拒绝），且 subject 必须是内置 generator 写入的
  `id@name` 形式。
- refresh token 在签发后 `vef.security.refresh_not_before`（默认 `15m`）
  内不可用——过早兑换会以 `1004`（`ErrCodeTokenNotValidYet`）失败。
- 会通过 `UserLoader.LoadByID(...)` 重新加载用户，使被停用的账号无法继续
  刷新；loader 的错误原样返回给调用方。
- 每次兑换返回全新 token 对。提交的 refresh token 不会在服务端吊销——该
  机制是无状态的——它只会自然过期。
- 典型失败：`1003`（令牌畸形、`typ` 不对、subject 形状不对，HTTP 401）、
  `1002`（已过期，HTTP 401）、`1004`（尚未生效，HTTP 401）、`1400`
  （`refreshToken` 缺失/为空，HTTP 400）。

### `logout`

无参数。总是返回成功且 `data` 为空——从客户端视角 logout 刻意设计为不可
失败，无论如何客户端都必须丢弃已存储的 token。

- `opaque_token` 下，尽力吊销当前 bearer token 背后的会话：按与 bearer
  认证完全相同的方式读取 token（`Authorization: Bearer` 头，scheme 大小
  写不敏感，其次 `__accessToken` 查询参数），哈希、查找并吊销。吊销成功
  会通知已注册的 `security.SessionRevocationListener`，关联授权
  ——例如 WebSocket 推送连接——随即被拆除。会话不存在或存储故障只记录
  日志，绝不使调用失败。
- `jwt_token` 下没有服务端会话：`logout` 实际上是 no-op，令牌失效等于
  客户端删除自己的副本。

### `resolve_challenge`

`ResolveChallengeParams`：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `challengeToken` | `string` | 是 | 上一次 `login` 或 `resolve_challenge` 响应中的状态令牌——必须是**最新**那个；每一步都会重新签发 |
| `type` | `string` | 是 | 正在解决的 challenge type。必须等于第一个 pending type（即刚返回的 `challenge.type`）；否则以 `1033` 失败 |
| `response` | `any` | 是 | provider 特定的应答，如 OTP 验证码字符串、新密码载荷或所选部门 ID |

响应是与 `login` 相同两种形态的 `LoginResult`：

- **仍有 challenge 待解**——新的 `challengeToken` 加下一个 `challenge`。
  挑战链严格按 provider 顺序推进；对该 principal `Evaluate(...)` 返回
  `nil` 的 provider 被跳过。新令牌携带更新后的 pending/resolved 列表和
  原始登录标识（保证审计连续性），并重新开始 `5m` 过期窗口。
- **全部挑战已解决**——`data.tokens` 携带最终 `AuthTokens`，与 `login`
  形态一完全一致。认证 token 只在此刻签发，成功的 `LoginEvent` 以原始
  登录标识发布。

行为说明：

- challenge token 的任何解析失败——过期、被篡改、`typ` 不对，或
  [Challenge providers](#challenge-providers) 中描述的保留身份拒绝
  （`system`/空/未知 principal type、保留 ID）——在该端点统一表现为
  `1031`（`ErrChallengeTokenInvalid`，HTTP 401）。
- 被拒绝的 `response` 按登录失败对待：计入原始标识的暴力破解锁定并被审
  计。返回类型化 `result.Error` 的 provider 保留自己的 code（`1035`
  `ErrOTPCodeRequired`、`1036` `ErrOTPCodeInvalid`、`1037`
  `ErrNewPasswordRequired`、`1038` `ErrDepartmentRequired`）；裸 error 被
  归一化为 `1034`（`ErrChallengeResolveFailed`，HTTP 401）。
- provider 解析出 nil 或框架保留 principal 会以 `1007`
  （`ErrReservedPrincipal`）拒绝；该拒绝会被审计但不计入锁定——第二因子
  本身是正确的，错在 provider。
- 类型不符与令牌无效这类协议错误（`1031`、`1033`）不经过 guard 也不审
  计；锁定检查（`1023`，HTTP 429）在 provider 校验应答之前进行。

### `get_user_info`

要求 Bearer 认证。`params` 对象不被框架解释：原样转发给应用的
`security.UserInfoLoader.LoadUserInfo(ctx, principal, params)`。未注册
loader 时，该 action 以通用 not-implemented 错误失败（code `1500`，
HTTP 501，消息 `security_user_info_loader_not_implemented`）；loader 的
错误原样返回。

响应 `data` 是 loader 返回的 `security.UserInfo`：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `id` | `string` | 用户标识 |
| `name` | `string` | 显示名 |
| `gender` | `string` | `male`、`female`、`unknown` 之一（`security.Gender`） |
| `avatar` | `string` \| `null` | 头像 URL；未设置时为 `null`（字段总是存在） |
| `permissionTokens` | `string[]` | 授予该用户的权限 token 列表，通常由前端消费以控制 UI 能力 |
| `menus` | `UserMenu[]` | 导航菜单树（见下） |
| `details` | `any` | 应用自定义的扩展载荷；缺省时省略（`omitempty`） |

`UserMenu`（递归）：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `type` | `string` | `directory`、`menu`、`view`、`dashboard`、`report` 之一（`security.UserMenuType`） |
| `path` | `string` | 路由路径 |
| `name` | `string` | 显示名 |
| `icon` | `string` \| `null` | 图标标识；未设置时为 `null`（总是存在） |
| `meta` | `object` | 可选扩展 map；缺省时省略 |
| `children` | `UserMenu[]` | 子节点；缺省时省略 |

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "id": "user001",
    "name": "Alice",
    "gender": "female",
    "avatar": null,
    "permissionTokens": ["user.read", "order.read"],
    "menus": [
      {
        "type": "directory",
        "path": "/system",
        "name": "System Management",
        "icon": "setting",
        "children": [
          { "type": "menu", "path": "/system/users", "name": "User Management", "icon": null }
        ]
      }
    ]
  }
}
```

`permissionTokens` 与 `menus` 没有 `omitempty`：loader 请返回空 slice
（而不是 nil），让客户端拿到 `[]` 而不是 `null`。

## 下一步

- [认证](./authentication) — 这些 API 背后的叙述性指南
- [会话管理](./session-management) — 令牌有效期、刷新与吊销
