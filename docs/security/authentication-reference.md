---
sidebar_position: 2
---

# Authentication Reference

The public `security` package surface for authentication: principals, JWT, the auth manager, challenge providers and token stores, signature auth, and login events. For the narrative guide — strategies, the built-in auth resource, and the login flow — see [Authentication](./authentication). The wire-level contract of the built-in auth endpoint — every action's request and response fields — is tabulated in [RPC Resource: `security/auth`](#rpc-resource-securityauth) at the end of this page.

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
audience, validates `iat` when present, requires `exp`, and applies a
10-second leeway. The compact claim keys are:

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

`ChallengeTokenStore.Generate(ctx, principal, username, pending, resolved)` and
`Parse(ctx, token)` carry the state between `login` and `resolve_challenge`
(`username` is the original login identifier the applicant supplied at the
first step, preserved across challenge steps for audit events).
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
(`ClaimChallengePending`), and `rsd` (`ClaimChallengeResolved`). Under the
reserved-identity hardening, challenge parsing accepts only `user` and
`external_app` principal types — `system`, empty, and unknown types are all
rejected with `ErrTokenInvalid` (a challenge token carrying the framework's
internal identity has no legitimate origin), and a parsed principal that
reports `IsReserved()` is rejected as well.

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

Security-domain API errors expose `ErrCode*` constants: `1000`–`1029` for
authentication, `1030`–`1039` for challenges, and `1050` for password policy
(every policy violation shares that one code — see
[Login Hardening](./login-hardening)):

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
| `1023` | `ErrCodeAccountLocked` | dynamic account-locked error (see [Login Hardening](./login-hardening)) | `429` |
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

Authentication-related sentinels include `ErrUnauthenticated`,
`ErrTokenExpired`, `ErrTokenInvalid`, `ErrTokenNotValidYet`,
`ErrTokenInvalidIssuer`, `ErrTokenInvalidAudience`, `ErrAppIDRequired`,
`ErrTimestampRequired`, `ErrSignatureRequired`, `ErrTimestampInvalid`,
`ErrSignatureExpired`, `ErrSignatureInvalid`, `ErrExternalAppNotFound`,
`ErrExternalAppDisabled`, `ErrIPNotAllowed`, `ErrNonceRequired`,
`ErrNonceInvalid`, `ErrNonceAlreadyUsed`, `ErrAuthHeaderMissing`,
`ErrAuthHeaderInvalid`, `ErrChallengeTokenInvalid`,
`ErrChallengeTypeInvalid`, `ErrChallengeResolveFailed`, `ErrOTPCodeRequired`,
`ErrOTPCodeInvalid`, `ErrNewPasswordRequired`, `ErrDepartmentRequired`,
`ErrTooManyConcurrentSessions`, `ErrAPIKeyInvalid`,
`ErrBasicCredentialsInvalid`, and `ErrReservedPrincipal` (rejects a
framework-internal identity at every entry point; it rides
`ErrCodePrincipalInvalid`/`1007` with HTTP 401), plus the factory helpers
`ErrCredentialsInvalid(message)` and `ErrPrincipalInvalid(message)`.
`ErrChallengeResolveFailed` is not a reserved placeholder:
`resolve_challenge` normalizes bare errors returned by a `ChallengeProvider`
into it.

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

## RPC Resource: `security/auth`

The security module mounts the built-in authentication resource as an RPC
resource under `/api`, using the standard envelope (`resource`, `action`,
`version`, `params`). Responses ride the standard result envelope — `code`
(`0` on success), `message`, `data` — and the shapes below describe the
`data` payload. Request parameter tables also appear in
[Built-in Resources](../reference/built-in-resources); this section is the
complete wire contract, including every response field.

| Action | Access | Rate limit (`max`) | Input | Output (`data`) |
| --- | --- | --- | --- | --- |
| `login` | Public | `vef.security.login_rate_limit` (default `6`) | `LoginParams` | `LoginResult` — tokens **or** a challenge envelope |
| `refresh` | Public; mounted only under `token_type = "jwt_token"` | `vef.security.refresh_rate_limit` (default `1`) | `RefreshParams` | `AuthTokens` (no `tokens` wrapper) |
| `logout` | Bearer auth | API engine default | none | empty (`data: null`) |
| `resolve_challenge` | Public | `vef.security.login_rate_limit` (default `6`) | `ResolveChallengeParams` | `LoginResult` — next challenge **or** final tokens |
| `get_user_info` | Bearer auth | API engine default | raw `params` map | `UserInfo` |

The custom limits set only `max`; the window falls back to the API engine's
default rate-limit period (`vef.api.rate_limit`, default `5m`). Operations
without a custom limit inherit the engine default entirely (stock `100`
requests per `5m`). Under `token_type = "opaque_token"` the `refresh`
operation is not mounted at all — calling it fails with the
operation-not-found error (HTTP 404), since opaque sessions renew themselves
on use.

### `login`

`LoginParams`:

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `type` | `string` | Yes | credential type. `password` is the only type the framework ships for this endpoint; custom `security.Authenticator` registrations extend the vocabulary. The framework-issued token types (`jwt_token`, `opaque_token`, `refresh`) are refused with code `1001` so an issued token can never be laundered into a fresh token pair |
| `principal` | `string` | Yes | login identifier, typically the username. The built-in password flow rejects the reserved identifiers (`system`, `cron_job`, `anonymous`) with code `1007` |
| `credentials` | `any` | Yes | credential payload. For `type = "password"` this is the password string — transport-encrypted when a `security.PasswordDecryptor` is configured, plaintext otherwise |

The response is a `LoginResult` and takes exactly one of two shapes.
`tokens`, `challengeToken`, and `challenge` are all `omitempty`: whichever
half does not apply is absent, never `null`.

**Shape 1 — tokens.** No challenge provider is registered, or none applies
to this account. `data.tokens` is an `AuthTokens`:

| Field | Type | Description |
| --- | --- | --- |
| `tokens.accessToken` | `string` | the bearer token for subsequent requests. Under `jwt_token` a JWT with a fixed `30m` TTL; under `opaque_token` a random session reference valid per the [session policy](./session-management) (`idle_ttl` / `max_lifetime`) |
| `tokens.refreshToken` | `string` | JWT refresh token, TTL `vef.security.token_expires` (default `168h`). **Omitted under `opaque_token`** — sessions renew themselves, so no refresh token exists |

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

The payload carries no expiry fields — token lifetimes are deployment
configuration, communicated out of band.

**Shape 2 — challenge envelope.** The credential verified, but at least one
[challenge provider](#challenge-providers) requires a second step. No auth
tokens are issued yet:

| Field | Type | Description |
| --- | --- | --- |
| `challengeToken` | `string` | state token carrying the challenge progress (principal, the original login identifier, pending and resolved types) — clients treat it as an opaque value and pass it to `resolve_challenge`. Each token expires after `ChallengeTokenExpires` (`5m`); every successful step issues a fresh one |
| `challenge` | `LoginChallenge` | the first pending challenge to resolve (below) |

`LoginChallenge`:

| Field | Type | Description |
| --- | --- | --- |
| `type` | `string` | challenge type wire value, e.g. `totp`, `sms_otp`, `password_change` (see [the wire-value table](#challenge-providers)) |
| `data` | `any` | provider-specific presentation data; omitted when the provider supplies none. The OTP providers return `{destination, meta?}` (`OTPChallengeData`); department selection returns `{departments, meta?}` |
| `required` | `bool` | whether the challenge must be resolved to finish the login |

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

Behavior notes:

- Providers are evaluated strictly in `Order()` sequence; providers whose
  `Evaluate(...)` returns `nil` are skipped, so the envelope always carries
  the first challenge that actually applies.
- The brute-force guard clears the failure counter as soon as the credential
  verifies — before any second factor. A rejected credential publishes a
  failure `LoginEvent`; the success event is published only when tokens are
  actually issued — immediately when no challenge applies, otherwise at the
  end of the challenge chain — always carrying the submitted identifier.
- Typical failures, all from the [error-code table](#signature-helpers)
  above: `1001` (unsupported/refused `type`, HTTP 400), `1008` (invalid
  credentials — deliberately the same uniform response for an unknown user,
  a nil principal or empty stored hash, and a wrong password, HTTP 401),
  `1007` (reserved or invalid principal, HTTP 401), `1023` (account locked
  by the brute-force guard, HTTP 429), and the generic `1400` validation
  error for missing required fields (HTTP 400).

### `refresh`

Mounted only under the stateless JWT mechanism (`token_type = "jwt_token"`,
the default).

`RefreshParams`:

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `refreshToken` | `string` | Yes | the refresh token issued by `login`, a previous `refresh`, or a token-issuing `resolve_challenge` |

The response `data` is the new `AuthTokens` pair **directly — without the
`tokens` wrapper `login` uses**:

| Field | Type | Description |
| --- | --- | --- |
| `accessToken` | `string` | fresh access token (`30m` TTL) |
| `refreshToken` | `string` | fresh refresh token (`vef.security.token_expires` TTL) |

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

Behavior notes:

- Internally the token is authenticated as type `refresh`: the JWT must
  parse, carry `typ: "refresh"` (an access token is refused), and its
  subject must be the `id@name` form the built-in generator writes.
- A refresh token is not usable before `vef.security.refresh_not_before`
  (default `15m`) has elapsed since issue — an early exchange fails with
  `1004` (`ErrCodeTokenNotValidYet`).
- The user is reloaded through `UserLoader.LoadByID(...)` so deactivated
  accounts stop refreshing; a loader error is returned to the caller as-is.
- Each exchange returns a fresh pair. The presented refresh token is not
  revoked server-side — the mechanism is stateless — it simply ages out.
- Typical failures: `1003` (malformed token, wrong `typ`, wrong subject
  shape, HTTP 401), `1002` (expired, HTTP 401), `1004` (not valid yet, HTTP
  401), `1400` (missing/empty `refreshToken`, HTTP 400).

### `logout`

No parameters. Always returns success with empty `data` — logout is
deliberately not failable from the client's perspective, and clients must
drop their stored tokens either way.

- Under `opaque_token`, the session backing the presented bearer token is
  revoked best-effort: the token is read exactly like bearer auth reads it
  (`Authorization: Bearer` header, case-insensitive scheme, then the
  `__accessToken` query parameter), hashed, looked up, and revoked. A
  successful revocation notifies the registered
  `security.SessionRevocationListener`s, so coupled grants — e.g.
  WebSocket push connections — are torn down immediately. A missing session
  or a store failure is only logged and never fails the call.
- Under `jwt_token`, there is no server-side session: `logout` is
  effectively a no-op and token invalidation is the client discarding its
  copy.

### `resolve_challenge`

`ResolveChallengeParams`:

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `challengeToken` | `string` | Yes | the state token from the previous `login` or `resolve_challenge` response — always the **latest** one; each step re-issues it |
| `type` | `string` | Yes | the challenge type being resolved. Must equal the first pending type (the `challenge.type` just returned); anything else fails with `1033` |
| `response` | `any` | Yes | provider-specific answer, e.g. the OTP code string, the new password payload, or the selected department ID |

The response is a `LoginResult` with the same two shapes as `login`:

- **Another challenge pending** — a fresh `challengeToken` plus the next
  `challenge`. The chain is strictly sequential in provider order;
  providers whose `Evaluate(...)` returns `nil` for this principal are
  skipped. The new token carries the updated pending/resolved lists and the
  original login identifier (for audit continuity), and restarts the `5m`
  expiry window.
- **All challenges resolved** — `data.tokens` with the final `AuthTokens`,
  exactly as in `login` shape 1. Only at this point are auth tokens issued,
  and the success `LoginEvent` is published with the original login
  identifier.

Behavior notes:

- Any challenge-token parse failure — expired, tampered, wrong `typ`, or
  the reserved-identity rejections (`system`/empty/unknown principal
  types, reserved IDs) described in
  [Challenge providers](#challenge-providers) — surfaces uniformly as
  `1031` (`ErrChallengeTokenInvalid`, HTTP 401) on this endpoint.
- A rejected `response` is treated like a failed login: it counts toward
  the brute-force lockout for the original identifier and is audited.
  Providers that return a typed `result.Error` keep their code (`1035`
  `ErrOTPCodeRequired`, `1036` `ErrOTPCodeInvalid`, `1037`
  `ErrNewPasswordRequired`, `1038` `ErrDepartmentRequired`); a bare error is
  normalized to `1034` (`ErrChallengeResolveFailed`, HTTP 401).
- A provider that resolves to a nil or framework-reserved principal is
  refused with `1007` (`ErrReservedPrincipal`); the rejection is audited but
  not counted toward lockout — the second factor was correct, the fault is
  the provider's.
- Wrong-type and invalid-token protocol errors (`1031`, `1033`) are not
  guarded or audited; the lockout check (`1023`, HTTP 429) applies before
  the provider verifies the response.

### `get_user_info`

Requires Bearer auth. The `params` object is not interpreted by the
framework: it is forwarded verbatim to the application's
`security.UserInfoLoader.LoadUserInfo(ctx, principal, params)`. When no
loader is registered, the action fails with the generic not-implemented
error (code `1500`, HTTP 501, message
`security_user_info_loader_not_implemented`); a loader error is returned
as-is.

The response `data` is the loader's `security.UserInfo`:

| Field | Type | Description |
| --- | --- | --- |
| `id` | `string` | user identifier |
| `name` | `string` | display name |
| `gender` | `string` | one of `male`, `female`, `unknown` (`security.Gender`) |
| `avatar` | `string` \| `null` | avatar URL; `null` when unset (the field is always present) |
| `permissionTokens` | `string[]` | permission tokens granted to the user, typically consumed by the frontend to gate UI affordances |
| `menus` | `UserMenu[]` | navigation menu tree (below) |
| `details` | `any` | application-defined extension payload; omitted when absent (`omitempty`) |

`UserMenu` (recursive):

| Field | Type | Description |
| --- | --- | --- |
| `type` | `string` | one of `directory`, `menu`, `view`, `dashboard`, `report` (`security.UserMenuType`) |
| `path` | `string` | route path |
| `name` | `string` | display name |
| `icon` | `string` \| `null` | icon identifier; `null` when unset (always present) |
| `meta` | `object` | optional extension map; omitted when absent |
| `children` | `UserMenu[]` | child entries; omitted when absent |

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

`permissionTokens` and `menus` carry no `omitempty`: return empty slices
(not nil) from your loader so clients see `[]` rather than `null`.

## Next Step

- [Authentication](./authentication) — the narrative guide these APIs back
- [Session Management](./session-management) — token lifetimes, refresh, and revocation
