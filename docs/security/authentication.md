---
sidebar_position: 1
---

# Authentication

VEF authentication happens at the API operation layer. Every operation has an auth configuration, and the API middleware chain resolves a principal before the handler runs.

## Default Behavior

If you do not configure anything special:

- operations are authenticated with the Bearer strategy
- `Public` operations are explicitly unauthenticated

That default comes from the API engine, not from your application config.

## Built-In Strategies

The public `api` package exposes strategy helpers:

- `api.Public()`
- `api.BearerAuth()`
- `api.SignatureAuth()`
- `api.IPAuth(...)` (see [Signature helpers](./authentication-reference#signature-helpers) below for how it resolves whitelists); the built-in strategy is fail-closed: a missing, empty, or unparseable whitelist denies every request with `security.ErrIPNotAllowed` (HTTP 401). That "empty means allow all" behavior belongs to the lower-level `IPWhitelistValidator` used by signature external apps, not to `api.IPAuth(...)`
- `api.APIKeyAuth(...)`
- `api.HTTPBasicAuth()`

In practice, you normally control this through operation settings:

```go
api.OperationSpec{
	Action: "login",
	Public: true,
}
```

or resource-level auth defaults.

## Bearer Authentication

Bearer auth reads tokens from:

- `Authorization: Bearer <token>`
- query parameter `__accessToken`

The API auth strategy delegates actual token validation to the security module's auth manager.

## Signature Authentication

Signature auth is intended for external applications and request signing use cases.

It expects these headers:

- `X-App-ID`
- `X-Timestamp`
- `X-Nonce`
- `X-Signature`

The strategy delegates verification to the security module's signature authenticator.

## API Key Authentication

`api.APIKeyAuth()` authenticates machine-to-machine callers by a static
secret key, read from the `X-API-Key` header by default; pass one header
name (`api.APIKeyAuth("X-Custom-Key")`) to read a custom header.

The presented key is resolved through the registered `security.APIKeyLoader`.
The framework ships a configuration-backed loader over
`vef.security.api_keys` (constant-time comparison across all configured
keys):

```toml
[vef.security.api_keys.reporting]
key = "high-entropy-random-string"
roles = ["reporting"]
```

Applications may provide their own `security.APIKeyLoader` to load keys from
a database or config center:

```go
type APIKeyLoader interface {
    // LoadByKey resolves the presented key to its Principal, or nil when no
    // key matches. An error signals an infrastructure fault, not a rejection.
    LoadByKey(ctx context.Context, key string) (*security.Principal, error)
}
```

Implementations that scan candidate keys must compare in constant time;
implementations that index by key should only serve high-entropy random keys,
where lookup timing reveals nothing.

A missing or unmatched key is rejected uniformly with
`security.ErrAPIKeyInvalid` (HTTP 401). The configuration-backed loader
resolves matches to an external-app principal named after the entry
(`api_key:<name>`) carrying the configured roles.

## HTTP Basic Authentication

`api.HTTPBasicAuth()` authenticates RFC 7617 `Authorization: Basic`
credentials. These are machine-to-machine service accounts — store
high-entropy random secrets, not user passwords.

The framework ships a configuration-backed loader over
`vef.security.basic_accounts` (the map key is the username):

```toml
[vef.security.basic_accounts.metrics-scraper]
password = "high-entropy-random-string"
roles = ["metrics"]
```

Applications may provide their own `security.BasicAccountLoader`; the loader
returns the stored secret and the framework performs the constant-time
comparison, so every implementation shares the same fail-closed semantics:

```go
type BasicAccountLoader interface {
    // LoadByUsername retrieves a service account by username, returning the
    // Principal and its stored secret. A nil Principal or empty secret means
    // the account is unknown; an error signals an infrastructure fault.
    LoadByUsername(ctx context.Context, username string) (*security.Principal, string, error)
}
```

Malformed headers, unknown accounts, and wrong passwords are all rejected
uniformly with `security.ErrBasicCredentialsInvalid` (HTTP 401), so callers
cannot distinguish which part failed.

## Reserved Identities

Certain identities attribute work the framework performs outside any request
— they are audit authors, never callers. `security.Principal.IsReserved()`
reports them:

- the `system` principal type (`PrincipalTypeSystem`);
- a principal whose `ID` equals `orm.OperatorSystem` (`"system"`);
- a principal whose `ID` equals `orm.OperatorCronJob` (`"cron_job"`).

`PrincipalAnonymous` is deliberately **not** reserved: it denotes the absence
of an identity, which the `public` auth strategy produces legitimately.

The framework enforces this invariant fail-closed at every boundary:

- **Authentication boundary**: the API auth middleware rejects any strategy
  that returns a nil or reserved principal with `security.ErrReservedPrincipal`
  (code `1007`, HTTP 401).
- **Token issuance**: both `JWTTokenGenerator` and `OpaqueTokenGenerator`
  refuse to mint tokens for reserved principals.
- **Challenge flow**: challenge-token parsing rejects reserved principal
  types and reserved IDs as `ErrTokenInvalid`; after a challenge provider
  resolves, `resolve_challenge` refuses a reserved result with
  `ErrReservedPrincipal`. Challenge tokens carry only the principal ID in
  their subject; user name, department, and other metadata live in separate
  claims.

The built-in password login additionally refuses the reserved identifiers
(`system`, `cron_job`, and `anonymous`) as a login `principal`. Custom
`UserLoader`, `APIKeyLoader`, and similar identity sources must never return
principals whose ID collides with the reserved operator IDs.

Reserved-principal rejections are **audited but not counted toward login
lockout**: the credential may have been correct, so the fault lies with the
authenticator or challenge provider, not the caller. See
[Login Hardening](./login-hardening) for the lockout behavior.

## Public Operations

Public operations are intentionally anonymous. The auth middleware injects an anonymous principal instead of rejecting the request.

Use `Public` for:

- login
- token refresh
- health-like anonymous endpoints
- public callbacks when appropriate

## Built-In Auth Resource

The security module registers a built-in RPC resource at:

```text
security/auth
```

Its main actions are:

- `login`
- `refresh`
- `logout`
- `resolve_challenge`
- `get_user_info`

The request fields, public flags, and rate-limit sources are part of the
runtime contract:

| Action | Public | Rate limit | Request fields |
| --- | --- | --- | --- |
| `login` | yes | `vef.security.login_rate_limit` | `type`, `principal`, `credentials`; all are `validate:"required"` |
| `refresh` | yes | `vef.security.refresh_rate_limit` | `refreshToken`; `validate:"required"`. Only mounted under `token_type = "jwt_token"` — under `opaque_token` the operation does not exist (sessions renew themselves) |
| `logout` | no | default API rate limit | none |
| `resolve_challenge` | yes | `vef.security.login_rate_limit` | `challengeToken`, `type`, `response`; all are `validate:"required"` |
| `get_user_info` | no | default API rate limit | arbitrary `params`, forwarded to `UserInfoLoader.LoadUserInfo(...)` |

The complete field-level wire contract — every action's request parameters
*and* response fields, with JSON examples for both login response shapes —
is tabulated in
[RPC Resource: `security/auth`](./authentication-reference#rpc-resource-securityauth).

This resource, every registered `Authenticator`, and the `AuthManager`
aggregator are wired by the framework's security module — the same module that layers in
brute-force lockout, password strength/history/expiry (see
[Login Hardening](./login-hardening)) and opaque-token session control (see
[Session Management](./session-management)).

The built-in authenticator type strings are `password`, `jwt_token`,
`opaque_token`, `refresh`, and `signature` (the JWT authenticator is named
`jwt_token`, not `token`). In normal client calls, `security/auth.login`
uses `type: "password"` with username and password credentials.
Bearer-protected operations dispatch the configured token mechanism internally
(`jwt_token` or `opaque_token` per `vef.security.token_type`),
`security/auth.refresh` uses `refresh` internally, and `SignatureAuth` maps
the signature headers to the `signature` authenticator. Only the configured
mechanism's authenticators are registered, and `login` refuses the
framework-issued token types as login credentials (see
[Session Management](./session-management)).

`logout` always returns an ok result. Under `jwt_token` it is effectively a
no-op — there is no server-side session to revoke, clients are expected to
remove their stored tokens. Under `opaque_token` it revokes the session
backing the presented bearer token, best-effort (a missing session or a store
failure is only logged).

## Login Flow

The auth resource supports a two-phase model:

1. authenticate credentials
2. optionally continue through challenge providers

If no challenge is required, `login` returns tokens directly.

If challenges are configured and applicable, `login` returns:

- a challenge token
- the next required challenge

Clients then call `resolve_challenge` until all required challenges are complete.
At the Go API layer, this shape is represented by `LoginResult`; the active
step is a `LoginChallenge`.

The login response DTOs use these exact fields:

| DTO | Fields |
| --- | --- |
| `AuthTokens` | JSON `accessToken`, `refreshToken` |
| `Authentication` | JSON `type`, `principal`, `credentials` |
| `LoginResult` | JSON `tokens`, `challengeToken`, `challenge` |
| `LoginChallenge` | JSON `type`, `data`, `required` |
| `LoginContext` | Go-only, read-only: `AuthType`, `Username`, `Principal`, `Resolved` — the login a challenge runs within (see [The login context](#the-login-context)) |
| `ChallengeState` | Go-only state the challenge token carries: the embedded `LoginContext` plus `Pending`, the challenge types still ahead in evaluation order (the first is the one presented) |

Field-by-field tables for both response shapes — the token payload and the
challenge envelope — with JSON examples live in
[RPC Resource: `security/auth`](./authentication-reference#rpc-resource-securityauth).

### The login context

A challenge always runs within one login, and its provider is handed that login
as a `*security.LoginContext`:

| Field | Holds |
| --- | --- |
| `AuthType` | the login mechanism that authenticated the principal — the `type` sent to `login`: `security.AuthTypePassword` (`password`), `security.AuthTypeTrustCode` (`trust_code`), or a type a host `security.Authenticator` supports. It is the same on every step of one login |
| `Username` | the identifier sent to `login` as `principal`, so events raised after a challenge report the identifier first presented |
| `Principal` | the identity as enriched by the challenges resolved so far |
| `Resolved` | the challenge types resolved so far, in resolution order |

`ChallengeProvider.Evaluate(ctx, login)` decides whether its challenge applies
to the login, returning nil when it does not; `Resolve(ctx, login, response)`
checks the answer and returns the principal the login continues with —
`login.Principal`, or an enriched copy such as the one department selection
produces. The challenge token carries the context across every
`resolve_challenge` step, so a later step sees the same `AuthType` and
`Username` as the first. The framework owns the context and passes it by
pointer: treat it as read-only, and change the identity only by returning a
principal from `Resolve`.

The hooks behind the built-in providers divide along one rule: a hook that
decides *whether* a challenge applies sees the whole login, while a hook that
*acts on* the identity sees only the principal — how the user logged in has no
bearing on how a password is stored or a code is checked.

| Decides whether it applies (takes `login`) | Acts on the identity (takes `principal`) |
| --- | --- |
| `PasswordChangeChecker.Check` — including `ExpiryPasswordChangeChecker` and `NewCompositePasswordChangeChecker` | `PasswordChanger`, `PasswordValidator`, `PasswordMetadataLoader` |
| `OTPEvaluator.Evaluate` — including `TOTPEvaluator` | `OTPCodeSender`, `OTPCodeVerifier`, `OTPCodeStore`, `OTPCodeDelivery`, `TOTPSecretLoader` |
| `DepartmentLoader.LoadDepartments` | `DepartmentSelector` |

A deciding hook reads the user from `login.Principal`. Per-user conditions —
whether the user has a TOTP secret, whether the password has expired — belong
in these hooks; which login mechanisms a challenge belongs to is better
declared once, at registration. The method signatures are tabulated in the
[Authentication Reference](./authentication-reference#challenge-providers).

### Scoping challenges to login mechanisms

A provider applies to every login unless it decides otherwise. To scope one to
some login mechanisms without touching it, wrap it at registration with
`security.NewFilteredChallengeProvider(provider, filters...)`. Each
`security.LoginFilter` states the scope as data:

| Field | Constructor | Matches a login when |
| --- | --- | --- |
| `AuthTypes` | `security.ForAuthTypes(types...)` | its `AuthType` is listed — an allow-list |
| `ExcludedAuthTypes` | `security.ExceptAuthTypes(types...)` | its `AuthType` is not listed — a deny-list |

- Within one filter an empty dimension is unconstrained, so `LoginFilter{}`
  matches every login; a filter populating both dimensions requires both.
- Several filters passed to one provider must all match (AND).
- For a login the filters reject, the provider is skipped exactly as if its
  `Evaluate` had returned nil, and the chain moves on to the next provider.
- No filters returns the provider unchanged. `LoginFilter.Matches(login)` is
  the predicate itself.

Take an application that logs users in with passwords and with a host-defined
WeChat mini-program mechanism — `type: "wechat_mini"`, served by its own
authenticator — and also accepts trust-login handoffs. A forced password change
belongs to password logins. The TOTP second factor belongs to every login except
a handoff, whose initiating system this application trusts to have enforced its
own:

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

Each constructor returns `security.ChallengeProvider`, the type the provider
group is declared with: fx keys a group by type, so a constructor returning a
concrete provider type — as `security.NewTOTPChallengeProvider` itself does —
is dropped without an error. The application supplies the checker, the
changer, and the TOTP secret loader; the `PasswordValidator` is the framework's
own, built from `vef.security.password_policy`. The chain each login meets:

| Login `type` | `totp` (order `100`) | `password_change` (order `400`) |
| --- | --- | --- |
| `password` | evaluated | evaluated |
| `wechat_mini` | evaluated | skipped |
| `trust_code` | skipped | skipped |

"Evaluated" hands the decision to the provider's own hook: TOTP is still
skipped for a user without a secret, and the password change for a user who
need not change.

:::caution[Choose the list by what the challenge guards]
A challenge tied to one credential takes an allow-list: a forced password
change concerns the password, so `ForAuthTypes(security.AuthTypePassword)`
keeps it off logins that never presented one. A second factor takes a
deny-list: exempt with `ExceptAuthTypes(...)` only the mechanisms that already
carry equivalent assurance. Written as an allow-list, a second factor would
silently exempt every login mechanism added later; as a deny-list, a new
mechanism is challenged until someone decides otherwise.
:::

Department selection is a required business input and is normally left
unfiltered. How a trust-login handoff meets the chain is covered in
[Trust Login](./trust-login).

## What Applications Usually Provide

The exact application-owned pieces depend on which auth paths you use:

- `security.UserLoader` is typically required for user login and refresh flows
- `security.ExternalAppLoader` is needed for signature-based external app auth
- challenge providers are optional and only matter if you use challenge-based login flows
- `security.UserInfoLoader` is only needed if you want `security/auth.get_user_info` to return application-defined user data

The framework ships the auth flow and middleware, but application identity sources remain application-owned.

## Public API Surface

The complete public authentication surface — principals, JWT, the auth manager, challenge providers and token stores, signature auth, and login events — is indexed with contract notes in the [Authentication Reference](./authentication-reference).

## A Working Login Module

In real VEF apps, the auth module is often very small: one users table, one package that implements the loader interfaces, and a module declaration that provides them. The framework already ships the `security/auth` resource, the password and refresh authenticators, and a default bcrypt `password.Encoder`; the application only supplies its identity source. The `auth` package below is complete enough to log in against a real table.

### The user model

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

`PasswordHash` must hold output of the same `password.Encoder` the login flow uses — by default the security module provides bcrypt (`password.NewBcryptEncoder`). Wherever you create or seed users, inject `password.Encoder` and store `encoder.Encode(plaintext)`; the built-in password authenticator later verifies the login credential with `encoder.Matches(plaintext, storedHash)`.

### UserLoader

`security.UserLoader` has exactly two methods: `LoadByUsername` powers `type: "password"` login and returns the principal plus the stored hash; `LoadByID` powers token refresh.

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

The error semantics match what the built-in authenticators expect:

- `Scan` already maps "no rows" to `result.ErrRecordNotFound`, so returning the error unchanged is correct. Filtering on `is_active` makes disabled users indistinguishable from missing ones.
- During `login`, any `LoadByUsername` error — and equally a `nil` principal or an empty hash — collapses into the generic invalid-credentials error (code `1008`), so usernames cannot be enumerated. Record-not-found errors are logged at info level, everything else at warn level.
- During `refresh`, a `LoadByID` error is returned to the caller as-is; the refresh authenticator reloads the user precisely so deactivated accounts stop refreshing.

### Permissions and user info

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

A production `RolePermissionsLoader` reads a role-permissions table instead of a switch; the security module automatically wraps whatever you provide in a cache invalidated by `RolePermissionsChangedEvent`. The permission tokens feed the RBAC checker described in [Authorization](./authorization).

### Wiring

Constructors must return the interface types — the framework consumes `security.UserLoader`, `security.UserInfoLoader`, and `security.RolePermissionsLoader` from the DI graph as optional dependencies of exactly those types.

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

Pass `auth.Module` to `vef.Run(...)` in `main` and the built-in `security/auth` resource picks the loaders up — no further registration is needed. This keeps authentication integration isolated from the rest of the application modules.

### Logging in

With a seeded user (`admin` / `ChangeMe_123`, hash produced by the bcrypt encoder), call the built-in resource:

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

With no challenge providers registered, the response carries the token pair directly:

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

Access tokens expire after 30 minutes (a fixed framework constant); the refresh token lifetime comes from `vef.security.token_expires` (default 7 days). Exchange the refresh token for a new pair — note that `refresh` returns the token pair directly in `data`, without the `tokens` wrapper:

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

Request parameters for every `security/auth` action are tabulated in [Built-in Resources](../reference/built-in-resources); the response fields — including the challenge envelope and the `get_user_info` `UserInfo` shape — are covered in [RPC Resource: `security/auth`](./authentication-reference#rpc-resource-securityauth).

## Practical Advice

- use `Public` sparingly and intentionally
- keep browser/API user auth on Bearer unless you have a reason to change it
- use Signature auth for external system integration, not as a replacement for normal user sessions

## Next Step

- [Authentication Reference](./authentication-reference) — the complete public API surface behind this guide
- [Authorization](./authorization) — how authentication leads into permission checks
