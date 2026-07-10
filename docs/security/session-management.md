---
sidebar_position: 7
---

# Session Management

VEF's login layer supports two independent token mechanisms, selected by
`vef.security.token_type`. This page covers the stateful one — opaque tokens —
and the server-side session control that comes with it. If you have not
configured `token_type` at all, your application is running the stateless
default and none of this applies yet.

## `jwt_token` vs. `opaque_token`

| | `jwt_token` (default) | `opaque_token` |
| --- | --- | --- |
| Token shape | Self-contained JWT, principal encoded in the claims | Random reference token, principal held server-side |
| Server lookup per request | None | One `SessionStore.Lookup` per request |
| Revocation before expiry | Not possible | `logout`, admin revoke, force-logout |
| Concurrent-session limits | Not possible | `max_concurrent` + `on_exceed` |
| "Active devices" / kick-offline | Not possible | `ListByUser`, `Revoke`, `RevokeUser` |
| Scales across nodes without shared state | Yes | Only with a shared store (Redis) |

A stateless JWT costs nothing to verify and needs no shared infrastructure, but
once it is issued the server cannot take it back — there is no session to
revoke. An opaque token trades that for a lookup on every request in exchange
for real session control: you can force a device offline, cap how many
sessions one account may hold at once, and see who is logged in.

Switch mechanisms with a single config key:

```toml
[vef.security]
token_type = "opaque_token" # default: "jwt_token"
```

`token_type` is a `config.TokenType`; its only accepted values are the
constants `config.TokenTypeJWT` (`"jwt_token"`) and `config.TokenTypeOpaque`
(`"opaque_token"`) — an unrecognized value fails config validation at boot
(`SecurityConfig.Validate`). `AuthResource.Login`, `refresh`, and `logout` behave
identically from the client's perspective either way; only the token shape and
the server-side capabilities differ. `TokenGenerator.Generate(ctx, principal,
SessionMeta)` is mechanism-agnostic, so custom login flows do not need to know
which one is active.

## How Opaque Sessions Work

Under `opaque_token`, login does not mint a JWT. Instead `OpaqueTokenGenerator`:

1. generates a high-entropy random token (`security.GenerateOpaqueToken`)
2. opens a `security.Session` record — id, user id, a snapshot of the
   `*security.Principal`, client IP, user agent, and timestamps — keyed by the
   token's SHA-256 hash (`security.HashOpaqueToken`), never by the raw token
3. stores it in a `security.SessionStore`
4. returns the raw token as `AuthTokens.AccessToken` (there is no refresh
   token — a session renews itself on use, so there is nothing separate to
   refresh)

On each authenticated request, `OpaqueTokenAuthenticator` hashes the presented
bearer token and calls `SessionStore.Lookup`. A hit returns the session's
`Principal` snapshot directly — no database round trip for user data, since the
principal is what was current when the session started (or last renewed).

### Sliding idle timeout, capped by an absolute lifetime

Two independent lifetimes govern a session:

- **`idle_ttl`** — how long a session survives without activity. Every
  authenticated request slides it forward by another `idle_ttl` (when sliding
  is enabled).
- **`max_lifetime`** — a hard ceiling on total session age, measured from
  `CreatedAt`, that no amount of activity can push past.

`OpaqueTokenAuthenticator.renew` computes the next expiry as `now + idle_ttl`,
then clamps it to `CreatedAt + max_lifetime` if that would exceed it. So a
continuously active session still expires at the latest `max_lifetime` after
login — sliding extends idle survival, it does not extend the account's
absolute session budget.

```toml
[vef.security.session]
idle_ttl = "30m"    # default
max_lifetime = "168h" # default: 7 days
sliding = true       # default; omit or set false to disable renewal
```

- `idle_ttl` default: `30m` (`config.DefaultSessionIdleTTL`)
- `max_lifetime` default: `168h` / 7 days (`config.DefaultSessionMaxLifetime`)
- `sliding` is a `*bool`: an omitted key resolves to enabled; set it to `false`
  explicitly to make sessions expire strictly `idle_ttl` after login regardless
  of activity

Renewal is best-effort: a `SessionStore` error during renewal is logged and
swallowed, never surfaced as a failed request, so a transient store hiccup
never logs an otherwise-valid user out. `ExpiresAt` is the single authoritative
field for both bundled stores — the Redis store's own key TTL is refreshed
alongside it, but a stale key that outlives `ExpiresAt` for any reason is still
treated as expired on read.

## Concurrency Control (Kicking Devices Offline)

`vef.security.session.max_concurrent` bounds how many live sessions one
account may hold. `0` (the default) means unlimited. When a login would push
the account over the limit, `on_exceed` decides what happens:

```toml
[vef.security.session]
max_concurrent = 3
on_exceed = "evict_oldest" # default; or "reject"
```

- **`evict_oldest`** (default, `security.SessionExceedEvictOldest`) — the new
  login revokes the account's oldest sessions until admitting it would exactly
  reach `max_concurrent`. This is "kick the earliest device offline": a new
  device login silently signs an old one out.
- **`reject`** (`security.SessionExceedReject`) — the new login is denied
  instead. `Login` returns `security.ErrTooManyConcurrentSessions` (business
  code `1024` — `security.ErrCodeTooManyConcurrentSessions`, HTTP `403`), and
  the account keeps its existing sessions.

Enforcement runs in `OpaqueTokenGenerator.enforceConcurrency`, before the new
session is created, and is **best-effort under concurrent logins**: counting
existing sessions and creating the new one are separate store calls, not one
atomic operation, so a burst of simultaneous logins for the same account can
briefly overshoot `max_concurrent` by the number of racing requests. This is a
policy / blast-radius limit, not a hard security boundary — treat it
accordingly, and note that `evict_oldest` self-heals any transient overshoot on
the next login.

These settings resolve into a single `security.SessionPolicy` —
`MaxConcurrent`, `OnExceed` (a `security.SessionExceedPolicy`), `IdleTTL`,
`MaxLifetime`, `Sliding` — assembled once from `vef.security.session` and
shared by the opaque-token generator and authenticator, rather than re-read
from raw config on every request.

## Logout and Revocation

`security/auth.logout` revokes the session backing the presented bearer token:

```go
func (a *AuthResource) Logout(ctx fiber.Ctx) error {
	a.revokeCurrentSession(ctx)

	return result.Ok().Response(ctx)
}
```

It looks the token up by its hash and calls `SessionStore.Revoke(ctx,
session.ID)`. This is best-effort and always returns `Ok` — a missing session
(already expired, or a JWT token under which no session ever existed) is not
an error, and a store failure during revoke is only logged. Under `jwt_token`,
`logout` is effectively a no-op: there is no session to revoke, and the client
is expected to discard its stored token (see
[Authentication](./authentication)).

## Building Session-Admin Endpoints

The `security.SessionStore` used by the framework is a regular DI-exposed
dependency, not a private implementation detail — inject it into your own
resources to build "my devices" or admin session management:

```go
type SessionResource struct {
	api.Resource
	store security.SessionStore
}

func NewSessionResource(store security.SessionStore) api.Resource {
	return &SessionResource{store: store, /* ... */}
}

// ListMyDevices returns the caller's own live sessions.
func (r *SessionResource) ListMyDevices(ctx fiber.Ctx, principal *security.Principal) error {
	sessions, err := r.store.ListByUser(ctx.Context(), principal.ID)
	if err != nil {
		return err
	}

	return result.Ok(sessions).Response(ctx)
}
```

`SessionStore` exposes exactly what an admin surface needs:

- `ListByUser(ctx, userID)` — a user's own live sessions, newest activity
  first — the basis for a self-service "active devices" list
- `Revoke(ctx, id)` — revoke one session by its public `Session.ID` (never by
  token) — "sign this device out"
- `RevokeUser(ctx, userID)` — revoke every session for a user in one call —
  force-logout, e.g. on password reset or account suspension

`Session.ID` is a random public identifier deliberately separate from the
token hash, so it is safe to return to clients in a device list without
exposing anything that could re-derive a live credential.

Building these endpoints is entirely your responsibility, including their
authorization — the framework does not ship an admin UI or default
authorization rule over `SessionStore`, only the storage contract. Apply
[Authorization](./authorization) (and typically restrict `RevokeUser` /
cross-user reads to admin roles) the same way you would any other
privilege-sensitive endpoint.

### Cross-user visibility with `SessionInspector`

For an "all online sessions" dashboard spanning every user, type-assert the
store for the optional `security.SessionInspector` capability (mirroring the
`event.StreamInspector` pattern used elsewhere in the framework):

```go
type SessionInspector interface {
	ListAll(ctx context.Context) ([]Session, error)
}
```

```go
if inspector, ok := r.store.(security.SessionInspector); ok {
	sessions, err := inspector.ListAll(ctx.Context())
	// ...
}
```

Both bundled stores (`MemorySessionStore`, `RedisSessionStore`) implement it.
`ListAll` is `O(all sessions)` — the Redis implementation performs a keyspace
`SCAN` rather than maintaining a global index set, so it never accumulates
tombstones for one-off or deleted accounts. Treat it as an infrequent
administrative read, not a request-path call; a deployment large enough to
need pagination should build that on top of its own store.

## Memory vs. Redis: Single-Node vs. Multi-Node

The default `SessionStore` is `security.NewMemorySessionStore()` — in-process
maps, wired automatically by the framework's security module. It works correctly for a
single instance, with sessions kept in three in-memory indexes (by id, by
token hash, by user) and expired entries reclaimed lazily on access.

It does **not** share state across processes. In a multi-node deployment, a
session created on one node is invisible to requests landing on another node —
override the store with `security.NewRedisSessionStore` via `fx.Decorate`:

```go
vef.Run(
	// ...
	fx.Decorate(security.NewRedisSessionStore),
)
```

`NewRedisSessionStore(client *redis.Client) security.SessionStore` only needs
the `*redis.Client` already provided by `internal/redis` when
`vef.redis.enabled = true` — no extra wiring beyond enabling Redis and
decorating the store.

```toml
[vef.redis]
enabled = true
# host, port, ... — see the Redis configuration reference

[vef.security]
token_type = "opaque_token"
```

`RedisSessionStore` keeps the same semantics as the memory store but backed by
Redis keys under `vef:security:session:` (`id:`, `token:`, `user:` sub-prefixes):

- every multi-key mutation (`Create`, `RevokeUser`, revoke-on-delete) runs
  inside a Redis `MULTI`/`EXEC` transaction, so a reader never observes a
  half-written or half-deleted session
- sliding renewal issues `SET ... XX` on the session record, which only
  succeeds if the key still exists — a renewal racing a concurrent `Revoke`
  can never resurrect a just-deleted session
- `Session.ExpiresAt` (not the Redis key TTL alone) remains the authoritative
  expiry check on every read, so both stores enforce `max_lifetime` identically

Both stores are safe for concurrent use and are drop-in replacements for each
other through the same `security.SessionStore` interface — nothing else in
your application needs to change when you switch.

## See Also

- [Authentication](./authentication) — how tokens are validated per request and the built-in `security/auth` resource
- [Login Hardening](./login-hardening) — brute-force lockout, password strength, and password history, which apply regardless of `token_type`
