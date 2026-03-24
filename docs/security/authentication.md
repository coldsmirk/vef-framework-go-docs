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

The public `api` package exposes three strategy helpers:

- `api.Public()`
- `api.BearerAuth()`
- `api.SignatureAuth()`

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

## Login Flow

The auth resource supports a two-phase model:

1. authenticate credentials
2. optionally continue through challenge providers

If no challenge is required, `login` returns tokens directly.

If challenges are configured and applicable, `login` returns:

- a challenge token
- the next required challenge

Clients then call `resolve_challenge` until all required challenges are complete.

## What Applications Usually Provide

The exact application-owned pieces depend on which auth paths you use:

- `security.UserLoader` is typically required for user login and refresh flows
- `security.ExternalAppLoader` is needed for signature-based external app auth
- challenge providers are optional and only matter if you use challenge-based login flows
- `security.UserInfoLoader` is only needed if you want `security/auth.get_user_info` to return application-defined user data

The framework ships the auth flow and middleware, but application identity sources remain application-owned.

## A common application pattern

In real VEF apps, the auth module is often very small. A common setup is:

- set the user details type during package init
- provide `UserLoader`
- provide `UserInfoLoader`

For example, an auth module often looks conceptually like this:

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

This keeps authentication integration isolated from the rest of the application modules.

## Practical Advice

- use `Public` sparingly and intentionally
- keep browser/API user auth on Bearer unless you have a reason to change it
- use Signature auth for external system integration, not as a replacement for normal user sessions

## Next Step

Read [Authorization](./authorization) to see how authentication leads into permission checks.
