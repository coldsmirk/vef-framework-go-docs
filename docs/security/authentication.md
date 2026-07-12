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
- `api.IPAuth(...)` (see [Signature helpers](./authentication-reference#signature-helpers) below for how it resolves whitelists)

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

This resource, every registered `Authenticator`, and the `AuthManager`
aggregator are wired by the framework's security module — the same module that layers in
brute-force lockout, password strength/history/expiry (see
[Login Hardening](./login-hardening)) and opaque-token session control (see
[Session Management](./session-management)).

The built-in authenticator type strings are `password`, `jwt_token`,
`opaque_token`, `refresh`, and `signature` (the former `token` string was
renamed to `jwt_token` in v0.38). In normal client calls, `security/auth.login`
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
| `ChallengeState` | Go-only state: `Principal`, `Pending`, `Resolved` |

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

Request parameters for every `security/auth` action are tabulated in [Built-in Resources](../reference/built-in-resources).

## Practical Advice

- use `Public` sparingly and intentionally
- keep browser/API user auth on Bearer unless you have a reason to change it
- use Signature auth for external system integration, not as a replacement for normal user sessions

## Next Step

- [Authentication Reference](./authentication-reference) — the complete public API surface behind this guide
- [Authorization](./authorization) — how authentication leads into permission checks
