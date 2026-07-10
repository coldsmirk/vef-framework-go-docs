---
sidebar_position: 2
---

# Authentication Reference

The public `security` package surface for authentication: principals, JWT, the auth manager, challenge providers and token stores, signature auth, and login events. For the narrative guide — strategies, the built-in auth resource, and the login flow — see [Authentication](./authentication).

| API group | Public surface |
| --- | --- |
| principals | `Principal`, `PrincipalType`, `NewUser`, `NewExternalApp`, `PrincipalSystem`, `PrincipalAnonymous`, `SetUserDetailsType`, `SetExternalAppDetailsType` |
| JWT | `JWT`, `JWTConfig`, `JWTClaimsBuilder`, `JWTClaimsAccessor`, `NewJWT`, `GenerateSecret`, token type constants, `DefaultJWTAudience`, `DefaultJWTSecret`, `JWTIssuer` |
| auth manager | `Authentication`, `AuthTokens`, `Authenticator`, `AuthManager`, `TokenGenerator`, `UserLoader`, `ExternalAppLoader`, `ExternalAppConfig`, `PasswordDecryptor` |
| challenge tokens | `ChallengeProvider`, `ChallengeState`, `ChallengeTokenStore`, `NewMemoryChallengeTokenStore`, `NewRedisChallengeTokenStore`, `NewJWTChallengeTokenStore` |
| OTP/challenges | `OTPEvaluator`, `OTPCodeSender`, `OTPCodeVerifier`, `OTPCodeStore`, `NewOTPChallengeProvider`, `NewDeliveredCodeSender`, `NewDeliveredCodeVerifier`, `NewDeliveredChallengeProvider`, `NewSMSChallengeProvider`, `NewEmailChallengeProvider` |
| TOTP/password/department | `NewTOTPEvaluator`, `NewTOTPVerifier`, `NewTOTPChallengeProvider`, `WithTOTPDestination`, `NewPasswordChangeChallengeProvider`, `NewDepartmentSelectionChallengeProvider` |
| signature auth | `Signature`, `SignatureCredentials`, `SignatureResult`, `SignatureAlgorithm`, `NewSignature`, `WithAlgorithm`, `WithTimestampTolerance`, `WithNonceStore`, `NonceStore`, `NewMemoryNonceStore`, `NewRedisNonceStore` |
| login events | `LoginEvent`, `LoginEventParams`, `NewLoginEvent`, `SubscribeLoginEvent` |

Bearer constants are `AuthSchemeBearer` and `QueryKeyAccessToken`. The token
type constants are `TokenTypeAccess`, `TokenTypeRefresh`, and
`TokenTypeChallenge`.

## JWT and principals

`NewJWT` expects `JWTConfig.Secret` to be a hex-encoded key and defaults an
empty audience to `DefaultJWTAudience`. Low-level `NewJWT` still falls back to
the public `DefaultJWTSecret` when the secret is empty; the framework security
module wraps this with a safer boot-time behavior that generates an ephemeral
key and warns. Use `GenerateSecret()` to provision a private production key for
`vef.security.secret`.

The built-in framework token generator issues access tokens with a fixed `30m`
TTL. `vef.security.token_expires` configures the refresh-token TTL instead
(default `168h`), and `vef.security.refresh_not_before` defaults to `15m`.
Access and refresh tokens generated together share the same `jti`.

JWT parsing accepts only `HS256`, requires issuer `JWTIssuer` (`vef`), validates
audience, requires `iat` and `exp`, and applies a 10-second leeway. The compact
claim keys are:

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

The built-in access and refresh token generator stores the subject as
`id@name`. `JWTTokenAuthenticator` rebuilds a user principal from that subject
without a database lookup. `JWTRefreshAuthenticator` also expects `id@name`,
but then reloads the user with `UserLoader.LoadByID(...)` using the `id` part.

`JWTClaimsBuilder` writes compact token claims with `WithID`, `WithSubject`,
`WithRoles`, `WithDetails`, `WithType`, and `WithClaim`. `JWTClaimsAccessor`
reads the same payload back with `ID`, `Subject`, `Roles`, `Details`, `Type`,
and `Claim`. Use `NewJWTClaimsBuilder()` and `NewJWTClaimsAccessor(...)` to
create those helpers directly.

`PrincipalTypeUser`, `PrincipalTypeExternalApp`, and `PrincipalTypeSystem`
describe the supported principal kinds. `SetUserDetailsType[T]()` and
`SetExternalAppDetailsType[T]()` configure process-global detail unmarshalling
targets; call them during startup before serving requests.

`Principal` serializes as JSON `type`, `id`, `name`, `roles`, and `details`.
`SetUserDetailsType[T]()` and `SetExternalAppDetailsType[T]()` require `T` to
be a struct or struct pointer and panic with `ErrUserDetailsNotStruct` or
`ErrExternalAppDetailsNotStruct` otherwise. They mutate package-level state and
should be treated as startup-only configuration. Unknown principal types keep
`details` as `map[string]any`; system principals deserialize with `details`
set to `nil`. The built-in special principals are `PrincipalSystem`
(`type: "system"`, id `system`, name `系统`) and `PrincipalAnonymous`
(`type: "user"`, id `anonymous`, name `匿名`).

## Challenge providers

Built-in challenge type constants include:

- `ChallengeTypeTOTP`
- `ChallengeTypeSMS`
- `ChallengeTypeEmail`
- `ChallengeTypePasswordChange`
- `ChallengeTypeDepartmentSelection`

Their wire values and default orders are:

| Constant | Wire value | Default order |
| --- | --- | --- |
| `ChallengeTypeTOTP` | `totp` | `100` |
| `ChallengeTypeSMS` | `sms_otp` | `200` |
| `ChallengeTypeEmail` | `email_otp` | `300` |
| `ChallengeTypePasswordChange` | `password_change` | `400` |
| `ChallengeTypeDepartmentSelection` | `department_selection` | `500` |

`ChallengeTokenStore.Generate(ctx, principal, pending, resolved)` and
`Parse(ctx, token)` carry the state between `login` and `resolve_challenge`.
The built-in login resources expose that state field as `challengeToken`.
`JWTChallengeTokenStore` is stateless; `MemoryChallengeTokenStore` is suitable
for tests or single-instance deployments; `RedisChallengeTokenStore` is for
distributed deployments. Challenge tokens expire after `ChallengeTokenExpires`.
The JWT-backed store uses `ClaimChallengePrincipalType`,
`ClaimChallengePrincipalName`, `ClaimChallengeUsername`,
`ClaimChallengePending`, and `ClaimChallengeResolved` as compact claim keys.

Challenge token stores have different wire/storage shapes:

| Store | Token/state contract |
| --- | --- |
| `JWTChallengeTokenStore` | JWT token, `typ: "challenge"`, 5-minute `ChallengeTokenExpires`, subject is principal ID only |
| `MemoryChallengeTokenStore` | UUID token stored in process memory for `ChallengeTokenExpires` |
| `RedisChallengeTokenStore` | UUID token stored under `vef:security:challenge:<token>` for `ChallengeTokenExpires` |

The JWT challenge claim keys are `ptp` (`ClaimChallengePrincipalType`), `pnm`
(`ClaimChallengePrincipalName`), `unm` (`ClaimChallengeUsername`), `pnd`
(`ClaimChallengePending`), and `rsd` (`ClaimChallengeResolved`). Challenge
parsing accepts empty principal type as a backwards-compatible user principal,
accepts `user`, `external_app`, and `system`, and rejects unknown principal
types.

Challenge providers are sorted by `Order()` in ascending order. The built-in
convenience providers use `100` for TOTP, `200` for SMS, `300` for email, `400`
for password change, and `500` for department selection. Providers that are not
registered, or whose `Evaluate(...)` returns `nil`, are skipped. During
`resolve_challenge`, the submitted `type` must match the first pending
challenge type or the framework returns `ErrChallengeTypeInvalid`.

`NewOTPChallengeProvider` is the generic constructor. Its
`OTPChallengeProviderConfig` requires `ChallengeType`, `Evaluator`, and
`Verifier`; `ChallengeOrder` controls evaluation order, and `Sender` is
optional and is used by delivered-code flows.
`OTPChallengeProvider` returns `OTPChallengeData` to the client when a challenge
is required. The
delivered-code helpers combine `OTPCodeStore` and `OTPCodeDelivery`:
`DeliveredCodeSender`, `DeliveredCodeVerifier`, `NewDeliveredCodeSender`,
`NewDeliveredCodeVerifier`,
`NewDeliveredChallengeProvider`, `NewSMSChallengeProvider`, and
`NewEmailChallengeProvider`.

`NewTOTPChallengeProvider` only needs a `TOTPSecretLoader`; if
`LoadSecret(...)` returns an empty string, the challenge is skipped.
`TOTPEvaluator`, `TOTPVerifier`, and `TOTPOption` are the lower-level pieces
behind the convenience provider. TOTP uses `TOTPDefaultDestination`
(`Authenticator App`) unless `WithTOTPDestination(...)` overrides it.

`NewPasswordChangeChallengeProvider` uses `PasswordChangeChecker` and
`PasswordChanger`; it returns `PasswordChangeChallengeData` when a password
change is required. Common reason constants are `PasswordChangeReasonFirstLogin`
(`first_login`) and `PasswordChangeReasonExpired` (`expired`). The concrete provider type is
`PasswordChangeChallengeProvider`.
`NewDepartmentSelectionChallengeProvider` uses `DepartmentLoader` and
`DepartmentSelector`; an empty department list skips the challenge, while
resolve expects a department ID string. `DepartmentSelectionChallengeData`
serializes as `departments` plus optional `meta`; each `DepartmentOption`
serializes as `id` and `name`.

The challenge constructors are wiring-time APIs. `NewOTPChallengeProvider`
panics when `ChallengeType`, `Evaluator`, or `Verifier` is missing.
`NewPasswordChangeChallengeProvider` panics when `PasswordChangeChecker` or
`PasswordChanger` is missing. `NewDepartmentSelectionChallengeProvider` panics
when `DepartmentLoader` or `DepartmentSelector` is missing.

## Signature helpers

`NewSignature(secret, ...)` requires a non-empty hex-encoded secret and defaults
to `SignatureAlgHmacSHA256` with a 5-minute timestamp tolerance. The option
type is `SignatureOption`. Other algorithm constants are
`SignatureAlgHmacSHA512` and `SignatureAlgHmacSM3`. `WithTimestampTolerance`
changes the accepted timestamp window and
`WithNonceStore` controls replay protection. Low-level `NewSignature(...)`
creates a `MemoryNonceStore` by default; pass `WithNonceStore(nil)` only when
you intentionally want to disable nonce storage for that helper. The built-in
`SignatureAuthenticator` injects the application `NonceStore` when one is
provided, otherwise each verification uses the low-level helper's process-local
memory store. `MemoryNonceStore` is local to one process; `RedisNonceStore` is
the distributed option. Stored nonces use twice the timestamp tolerance plus a
1-minute buffer as TTL.

The signed payload is exactly:

```text
app_id=<appID>&method=<method>&nonce=<nonce>&path=<path>&timestamp=<timestamp>
```

The fields are bound in that order. The `request body is not part of the
signature payload`. The `method` field is the HTTP method observed by the
server.

`NewIPWhitelistValidator` returns an `IPWhitelistValidator` from a
comma-separated list of IPs and CIDR ranges. An empty whitelist allows all IPs;
an invalid whitelist is fail-closed and denies all requests. When an
`ExternalAppConfig.IPWhitelist` is non-empty but the request IP cannot be
resolved, `SignatureAuthenticator` also fails closed with `ErrIPNotAllowed`.

`security.NewIPWhitelistValidatorFromEntries(entries)` is the slice-based
counterpart used by the built-in `api.IPAuth(...)` strategy. The strategy
resolves a named `security.IPWhitelist` through `security.IPWhitelistLoader`;
the default loader reads `vef.security.ip_whitelists`, while applications may
provide their own loader for database or config-center backed lists.

Signature storage keys and defaults:

| Contract | Value |
| --- | --- |
| request headers | `X-App-ID`, `X-Timestamp`, `X-Nonce`, `X-Signature` |
| algorithms | `HMAC-SHA256`, `HMAC-SHA512`, `HMAC-SM3` |
| default algorithm | `HMAC-SHA256` |
| default tolerance | `5m` |
| nonce TTL | `2*tolerance + 1m` |
| Redis nonce prefix | `vef:security:nonce:` |
| disable replay checking | `WithNonceStore(nil)` |

Security-domain API errors expose `ErrCode*` constants in the `1000-1039`
range:

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

Authentication-related sentinels include `ErrUnauthenticated`,
`ErrTokenExpired`, `ErrTokenInvalid`, `ErrTokenNotValidYet`,
`ErrTokenInvalidIssuer`, `ErrTokenInvalidAudience`, `ErrAppIDRequired`,
`ErrTimestampRequired`, `ErrSignatureRequired`, `ErrTimestampInvalid`,
`ErrSignatureExpired`, `ErrSignatureInvalid`, `ErrExternalAppNotFound`,
`ErrExternalAppDisabled`, `ErrIPNotAllowed`, `ErrNonceRequired`,
`ErrNonceInvalid`, `ErrNonceAlreadyUsed`, `ErrAuthHeaderMissing`,
`ErrAuthHeaderInvalid`, `ErrChallengeTokenInvalid`,
`ErrChallengeTypeInvalid`, `ErrOTPCodeRequired`, `ErrOTPCodeInvalid`,
`ErrNewPasswordRequired`, `ErrDepartmentRequired`, plus the factory helpers
`ErrCredentialsInvalid(message)` and `ErrPrincipalInvalid(message)`.
`ErrCodeChallengeResolveFailed` and `ErrChallengeResolveFailed` are reserved
for challenge resolution failures.

Low-level secret parsing errors use `ErrDecodeJWTSecretFailed`,
`ErrGenerateJWTSecretFailed`, `ErrDecodeSignatureSecretFailed`, and
`ErrSignatureSecretRequired`. `ErrUserDetailsNotStruct` and
`ErrExternalAppDetailsNotStruct` are raised when detail-type registration is not
given a struct or struct pointer.
Public i18n message ID constants include `ErrMessageChallengeResolveFailed`,
`ErrMessageCredentialsFormatInvalid`, `ErrMessageExternalAppLoaderNotImplemented`,
`ErrMessageUnauthenticated`, `ErrMessageUnsupportedAuthenticationType`,
`ErrMessageUserInfoLoaderNotImplemented`, and
`ErrMessageUserLoaderNotImplemented`.

## Next Step

- [Authentication](./authentication) — the narrative guide these APIs back
- [Session Management](./session-management) — token lifetimes, refresh, and revocation
