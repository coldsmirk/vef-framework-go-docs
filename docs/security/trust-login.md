---
sidebar_position: 8
---

# Trust Login (SSO Handoff)

Trust login lets a third-party system drop a user straight into your
application: it links to a framework-owned gateway with a signed user
identifier, and the gateway trades that handoff for an ordinary login session.
The user never sees a login form, and your application never learns the user's
password — the external system vouches for who they are.

It is off by default. Turn it on under `vef.security.trust_login`
(`config.TrustLoginConfig`).

## Two legs, and why

A single link that logs someone in is a credential in a URL — it lands in
browser history, in `Referer` chains, and in every reverse proxy's access log.
Trust login therefore splits the handoff into two legs, and both are required:

1. **The gateway** — `GET /sso/trust`, a real route mounted as an
   `app.Middleware` at order 460. It verifies an HMAC over the handoff,
   resolves the external user to a local principal, parks that identity under a
   one-time code, and 302s the browser to
   `<redirect>?app_id=…&code=…`.
2. **The exchange** — an ordinary `security/auth.login` call with
   `type: "trust_code"`, `principal` set to the app ID and `credentials` set to
   the code. The `TrustCodeAuthenticator` redeems the code and returns the
   parked principal.

Because the second leg is a normal login, the handoff inherits the **full login
pipeline** unchanged: the challenge chain (department selection, forced
`password_change`), token issuance under whichever `token_type` is configured,
session concurrency and eviction, and the login audit event. There is
deliberately **no "skip challenges" switch** — department selection is a
required business input and a forced password change is policy; skipping either
would be a hole, and skipping a second factor is a per-provider decision rather
than a blanket one.

The signed URL on its own is therefore harmless to keep: it mints a code, and
the code is what logs anyone in. The code is single-use, seconds-lived, and
bound to the browser that received it.

## Configuration

```toml
[vef.security.trust_login]
enabled = true
# path = "/sso/trust"      # config.DefaultTrustLoginPath
# code_ttl = "60s"         # config.DefaultTrustLoginCodeTTL
# bind_user_agent = true   # omitted resolves to enabled
# bind_client_ip = false

[vef.security.trust_login.rate_limit]
# max = 120                # config.DefaultTrustLoginRateLimitMax
# period = "1m"            # config.DefaultTrustLoginRateLimitPeriod

[vef.security.trust_login.apps.his]
redirect_urls = ["https://portal.example.com/sso/callback"]
```

| Key | Type | Meaning |
| --- | --- | --- |
| `enabled` | `bool` | Mounts the gateway route and registers the `trust_code` authenticator. While off, the route does not exist and `trust_code` is refused as an unsupported login type. |
| `path` | `string` | The gateway route. It is part of the signed payload, so changing it invalidates every link the external system has already generated. |
| `code_ttl` | `time.Duration` | How long an issued code stays redeemable. Keep it short — it rides a redirect URL. |
| `bind_user_agent` | `*bool` | Requires the browser redeeming the code to present the same `User-Agent` the gateway redirected. An omitted value resolves to enabled: the header cannot change within one redirect, so the binding costs nothing. |
| `bind_client_ip` | `bool` | Additionally requires the same source address. Off by default — a mobile client can change networks mid-redirect, and the resulting failure reads as a broken integration rather than a defense. |
| `apps.<appId>.redirect_urls` | `[]string` | Where a handoff from that app may land. `config.TrustLoginAppConfig` carries nothing else. |
| `rate_limit.max` / `rate_limit.period` | `int` / `time.Duration` | `config.TrustLoginRateLimitConfig`, counted per (app ID, client IP) per node. |

Defaults live next to the fields as `Effective*` accessors, so the parsed
config is never rewritten in place.

### Participation is a separate grant

An app must appear in `apps` to initiate a handoff at all. That is deliberate:
**being able to call your API never implies being able to log a user in.** The
signing secret comes from `security.ExternalAppLoader` — the same source the
API signature authenticator reads, so no secret enters a config file — but the
`apps` map is what grants browser single sign-on, and an app absent from it is
refused.

A redirect target must match an entry's scheme and host **exactly** and sit
under its path, compared segment-wise with dot segments resolved.

### Validation failures

`SecurityConfig.Validate` fails the boot rather than starting a
half-configured gateway:

| Error | Cause |
| --- | --- |
| `config.ErrTrustLoginPathInvalid` | `path` does not start with `/`. |
| `config.ErrTrustLoginAppsRequired` | Trust login is enabled but `apps` is empty — nobody could initiate a handoff. |
| `config.ErrTrustLoginRedirectsEmpty` | An app declares no `redirect_urls`, so it could never complete a handoff. |
| `config.ErrTrustLoginRedirectInvalid` | A redirect entry is not an absolute URL, or carries a query or fragment. |

## Signing the handoff

The gateway verifies an HMAC over the app ID, the request method and path, the
timestamp and nonce, **and** the two caller-supplied parameters that matter:
`user_id` and `redirect`. Those two ride in `security.SignatureRequest`:

```go
req := security.SignatureRequest{
    AppID:  "his",
    Method: "GET",
    Path:   "/sso/trust",
    BoundParams: map[string]string{
        "user_id":  "1001",
        "redirect": "https://portal.example.com/sso/callback",
    },
}
```

The rule of thumb: **a parameter the signed link carries that is not in
`BoundParams` is attacker-controlled.** An unsigned `user_id` would be a
complete authentication bypass — anyone holding one valid link could swap in
another person's identifier.

`BoundParams` values are the **decoded** ones, never their URL-encoded wire
form. Keys and values are percent-encoded into the canonical payload (RFC 3986
unreserved set, uppercase hex, space as `%20`) precisely because a signed
`redirect` routinely contains `&` and `=`; rendered raw, two different
parameter sets could flatten to the same string. A key colliding with the fixed
payload fields (`app_id`, `method`, `nonce`, `path`, `timestamp`) is refused
with `security.ErrSignatureBoundKeyReserved` rather than silently shadowing
one.

The resulting URL is:

```text
GET /sso/trust?app_id=his&user_id=1001&redirect=<encoded>&timestamp=<unix seconds>&nonce=<random>&signature=<hex hmac>
```

Two sharp edges worth stating to whoever implements the signing side:

- **The signed `path` must match exactly.** Fiber runs with `StrictRouting`
  off, so `/sso/trust/` still reaches the gateway but hashes differently and
  yields a permanent, opaque 401.
- A request with **no bound parameters** produces a payload byte-identical to
  plain API signature auth, so a third party that already signs API calls only
  has to add the two parameters — not learn a second scheme.

## Resolving the user

`security.TrustUserResolver` maps the external identifier onto a local
principal:

```go
type TrustUserResolver interface {
    ResolveUser(ctx context.Context, appID, externalUserID string) (*Principal, error)
}
```

Registering one is optional. With none, the gateway resolves through
`UserLoader.LoadByID`, which is already correct whenever both systems key users
by the same identifier. Implement it when the identifiers differ, or when the
mapping depends on which external system initiated the handoff.

:::caution[Resolution is an authentication decision]
The external system vouched for **who** the user is, not for whether your
application still admits them. Whichever of the two runs — your resolver or
`LoadByID` — must refuse a disabled, locked or expired account by returning a
nil principal, exactly as the password path does. Applications commonly keep
that check in `LoadByUsername` alone, where trust login never reaches it.
:::

Returning an error means "resolution failed" (an unreachable directory);
"no such user here" is a nil principal, not an error.

## The one-time code

`security.TrustCodeStore` parks the verified identity and redeems it exactly
once:

```go
type TrustCodeStore interface {
    Issue(ctx context.Context, state TrustCodeState, ttl time.Duration) (string, error)
    Consume(ctx context.Context, code string) (*TrustCodeState, error)
}
```

`security.TrustCodeState` carries the resolved `Principal`, the `AppID` that
initiated the handoff, and the `UserAgent` and `ClientIP` of the browser the
gateway redirected. Both browser fields are recorded unconditionally, so
turning a binding on takes effect on the next handoff instead of requiring a
reissue.

The exchange compares the recorded `UserAgent` against
`contextx.RequestUserAgent(ctx)` — recorded by the API auth middleware through
`contextx.SetRequestUserAgent` under the `contextx.KeyRequestUserAgent` key —
and, when `bind_client_ip` is on, the client IP against
`contextx.RequestIP(ctx)`.

`Consume` must read and invalidate atomically: single use is the whole security
value of the code, and a redirect URL sitting in browser history must not be
replayable. An unknown, already-redeemed or expired code all yield
`security.ErrTrustCodeInvalid` — the three are deliberately
indistinguishable.

### The store is selected by topology, not decorated

Like `lock.Locker`, the code store is **chosen from the deployment**, not
defaulted to memory and left to be swapped:

| Redis client | Store |
| --- | --- |
| available | `security.RedisTrustCodeStore` (`security.NewRedisTrustCodeStore`), which needs Redis 6.2+ for `GETDEL` |
| absent | `security.MemoryTrustCodeStore` (`security.NewMemoryTrustCodeStore`), with a boot warning |

The reason it is not a silent memory default: a memory store cannot serve a
second replica **at all**. A code issued on one node is unknown to every other,
so the exchange fails outright rather than merely weakening. The shared
`security.NonceStore` is selected the same way, and its boot warning is raised
only while trust login is on — a per-process nonce store merely weakens a
server-to-server API call, but it hands a browser-visible credential a fresh
redemption per node.

Codes are generated with `GenerateOpaqueToken` and keyed by
`HashOpaqueToken`, never stored raw, so a store leak yields no live codes.

## Errors

Gateway failures render as **plain text** with the error's status, never a JSON
envelope — the endpoint is navigated to directly by a browser.

| Sentinel | Code | Status | Meaning |
| --- | --- | --- | --- |
| `security.ErrTrustAuthFailed` | `security.ErrCodeTrustAuthFailed` (1060) | 401 | Every verification failure: bad signature, expired or replayed nonce, unknown or disabled app, malformed timestamp. |
| `security.ErrTrustRedirectNotAllowed` | `security.ErrCodeTrustRedirectNotAllowed` (1061) | 400 | The handoff authenticated, but the redirect matches no allowlist entry. |
| `security.ErrTrustUserNotResolved` | `security.ErrCodeTrustUserNotResolved` (1062) | 401 | No local user corresponds to the external identifier. |
| `security.ErrTrustCodeInvalid` | `security.ErrCodeTrustCodeInvalid` (1063) | 401 | The code is unknown, already redeemed, expired, or presented by a different browser. |

Every verification failure collapses into `ErrTrustAuthFailed` so app IDs
cannot be enumerated — and the caller has the same fix in every case: re-sign
the handoff correctly. `ErrTrustRedirectNotAllowed` is deliberately distinct
because its fix is a configuration change instead.

The 302 carries `Cache-Control: no-store` and `Referrer-Policy: no-referrer`,
so neither the code in `Location` nor the signature in the request URL travels
further than it must.

## Throttling

The gateway is public and performs an `ExternalAppLoader` lookup — typically a
database round trip — before it can reject anything, so it is rate-limited per
(app ID, client IP) per node through `rate_limit`. The default is deliberately
generous: a whole organization behind one NAT legitimately signs in through the
same address at the start of a shift. It bounds a flood; it does not police
logins.

The exchange is **exempt from the brute-force login guard**, unlike every other
login type. A code is `GenerateOpaqueToken` randomness, single-use and
seconds-lived, so no number of attempts approaches it — while the identity it
presents is the initiating system's app ID rather than a person. Counting
failures would put every user of that system into one lockout bucket that
anyone able to POST a wrong code could fill. Its real throttles are the gateway
rate limit above, the code's single use, and `vef.security.login_rate_limit`.

## Checklist for the initiating system

1. Obtain the app ID and signing secret through whatever provisions your
   `ExternalAppLoader`.
2. Ask the VEF application's operator to add the app and its redirect
   allowlist under `vef.security.trust_login.apps`.
3. Sign `app_id + method + path + user_id + redirect + timestamp + nonce`, with
   `user_id` and `redirect` as bound parameters.
4. Link the browser to `<base>/sso/trust?...`.
5. On the callback page, read `code` from the query and call
   `security/auth.login` with `type: "trust_code"`, `principal` = the app ID,
   `credentials` = the code. Handle the challenge response exactly as a
   password login would.
