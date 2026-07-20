---
sidebar_position: 6
---

# Login Hardening

The `security` package ships five independent, opt-in hardening layers around
the login endpoint. All of them are wired inside the framework's security
module, none of them change how a client calls `security/auth.login`, and you can enable any
subset of them — a fresh application with no extra config gets none of this
protection.

| Layer | Protects against | Config section | App-owned interface |
| --- | --- | --- | --- |
| Encrypted password transport | Credential sniffing between client and server | none (DI wiring only) | `security.PasswordDecryptor` |
| Brute-force lockout | Credential stuffing / password guessing | `vef.security.lockout` | none (swap `security.LoginGuard` for Redis) |
| Password strength | Weak passwords | `vef.security.password_policy` | none |
| Password history | Password reuse | `vef.security.password_policy.history_depth` | `security.PasswordHistoryStore` |
| Password expiry | Stale passwords | `vef.security.password_policy.max_age` | `security.PasswordMetadataLoader` |

This page covers each layer shallow-to-deep: wiring first, then the
config keys, then the extension interfaces an application implements.

## Encrypted Password Transport

If the client encrypts the password before sending it (a common "encrypt in
the browser, hash on the server" pattern for defense against network-level
sniffing), register a `security.PasswordDecryptor`:

```go
type PasswordDecryptor interface {
	Decrypt(encryptedPassword string) (string, error)
}
```

`PasswordAuthenticator` decrypts the transmitted credential to plaintext
*before* verifying it against the stored hash. This keeps `password.Encoder` a
plain KDF used identically for both storage and comparison — decryption is a
concern of the authenticator layer, not the encoder. See
[Password](./password) for the encoder itself and
[Cryptox](./cryptox) for the `Cipher`/`CipherSigner`
implementations you'd typically back a decryptor with.

`cryptox.NewRSA` already satisfies `PasswordDecryptor` — its `Decrypt(ciphertext string) (string, error)` method matches the interface verbatim, so you only need to supply it through DI:

```go
fx.Provide(func() (security.PasswordDecryptor, error) {
	return cryptox.NewRSA(privateKey, publicKey)
})
```

`PasswordDecryptor` is an optional dependency of `NewPasswordAuthenticator`
(along with `UserLoader` and any transport itself) — leaving it unregistered
means the authenticator treats the transmitted credential as plaintext, which
is the zero-config default.

A malformed ciphertext is treated exactly like a wrong password: the
authenticator runs a dummy KDF comparison so decrypt failures cost the same
time as a genuine mismatch, closing a timing side-channel that would
otherwise distinguish "bad ciphertext" from "wrong password" or "unknown
user".

## Brute-Force Lockout

`security.LoginGuard` throttles the login endpoint *before* authentication
runs, based on accumulated failures for an identity:

```go
type LoginGuard interface {
	Check(ctx context.Context, attempt LoginAttempt) (LoginDecision, error)
	RecordFailure(ctx context.Context, attempt LoginAttempt) (LoginDecision, error)
	RecordSuccess(ctx context.Context, attempt LoginAttempt) error
}

type LoginAttempt struct {
	Identity string // the login identifier the client sent
	ClientIP string // the resolved source address
}

type LoginDecision struct {
	Allowed    bool
	RetryAfter time.Duration // zero when Allowed is true
}
```

`AuthResource.Login` calls `Check` before authenticating, `RecordFailure`
when authentication fails, and `RecordSuccess` as soon as the credential
verifies (before any second-factor challenge, since the brute-forced
credential has already succeeded). Failures accumulate per `LockoutPolicy.Key`
and reset on success.

Since v0.38 the same guard also covers `resolve_challenge`: a failed
second-factor guess counts toward the same lockout key, and a tripped lockout
blocks both endpoints — so an attacker who reaches the challenge step cannot
brute-force it outside the lockout budget.

### Enabling and configuring it

Lockout is **on by default** (`max_failures = 10`). Configure it under
`vef.security.lockout`:

```toml
[vef.security.lockout]
enabled = true          # default: true
max_failures = 10       # default: 10
window = "15m"           # default: 15m — a spell with no new failures this long resets the counter
lock_duration = "15m"    # default: 15m — block length under the "lock" strategy
strategy = "lock"        # "lock" or "backoff", default: "lock"
backoff_base = "1s"      # default: 1s — first delay under the "backoff" strategy
backoff_max = "15m"      # default: 15m — cap on the backoff delay
key = "user_ip"          # "user", "ip", or "user_ip", default: "user_ip"
```

Set `enabled = false` to switch lockout off entirely. Every other field
resolves to its default when omitted or zero — read them through
`config.LockoutConfig`'s `Effective*` accessors, not the raw struct, if you
build a policy yourself.

- **`strategy = "lock"`** (Go constant `security.LockoutStrategyLock`) blocks
  all attempts for `lock_duration` once `max_failures` is reached — the
  attacker gets exactly `max_failures` guesses, and the failure that hits the
  threshold triggers the lock.
- **`strategy = "backoff"`** (Go constant `security.LockoutStrategyBackoff`)
  imposes an escalating delay instead: the failure that reaches the threshold
  starts a `backoff_base` wait, and each further failure doubles the wait,
  capped at `backoff_max`. A legitimate user is slowed down but never fully
  locked out, and an attacker cannot use it to lock a victim out indefinitely
  by feeding wrong passwords.
- **`key`** selects the identity dimension failures are counted by:
  `"user"` (`security.LockoutKeyUser` — per login identifier, across all
  source IPs), `"ip"` (`security.LockoutKeyIP` — per source address, across
  all identifiers), or `"user_ip"` (`security.LockoutKeyUserIP` — the
  default, per identifier-and-source pair, which throttles credential
  guessing without letting an attacker lock a victim out from every IP by
  guessing one account).

### Storage backend

The default guard is `security.MemoryLoginGuard`, built by
`security.NewMemoryLoginGuard(policy)` — suitable for a single instance.
Multi-node deployments override it with `security.NewRedisLoginGuard`, which
returns a `security.RedisLoginGuard` backed by Redis-shared counters, via
`fx.Decorate` so failure counters are shared across nodes:

```go
vef.Run(
	// ...
	fx.Decorate(func(client *redis.Client, cfg *config.SecurityConfig) security.LoginGuard {
		l := cfg.Lockout
		return security.NewRedisLoginGuard(client, security.LockoutPolicy{
			MaxFailures:  l.EffectiveMaxFailures(),
			Window:       l.EffectiveWindow(),
			LockDuration: l.EffectiveLockDuration(),
			Strategy:     security.LockoutStrategy(l.EffectiveStrategy()),
			BackoffBase:  l.EffectiveBackoffBase(),
			BackoffMax:   l.EffectiveBackoffMax(),
			Key:          security.LockoutKey(l.EffectiveKey()),
		})
	}),
)
```

This needs `vef.redis.enabled = true` so the `*redis.Client` is available in
DI (see [Session Management](./session-management) for the same pattern
applied to opaque-token session storage).

### Failure handling and the error it raises

A `LoginGuard` backend error (Redis unreachable, etc.) **fails open**: the
guard logs a warning and lets the attempt proceed rather than denying every
login because the counter store is down. The guard is defense-in-depth, not
the authoritative auth result.

A tripped lockout returns `security.ErrAccountLocked(retryAfter)` — HTTP 429,
business code `security.ErrCodeAccountLocked` (`1023`), with the retry window
rounded up to whole minutes (never below one) in the response message. The
message is rendered from the i18n key `security.ErrMessageAccountLocked`
(`"security_account_locked"`).

## Password Strength

`security.PasswordValidator` checks a candidate plaintext password against a
policy:

```go
type PasswordValidator interface {
	Validate(ctx context.Context, principal *Principal, plaintext string) error
}
```

It's built from composable `PasswordRule`s via `NewRuleBasedValidator`; with
no rules registered it accepts every password (the zero-config default):

```go
type PasswordRule interface {
	Check(principal *Principal, plaintext string) error
}
```

Built-in rules:

| Constructor | Rule |
| --- | --- |
| `NewMinLengthRule(minLength)` | at least `minLength` runes |
| `NewMaxLengthRule(maxLength)` | at most `maxLength` runes (also guards against slow-KDF DoS and silent bcrypt truncation) |
| `NewCharacterClassRule(requireUpper, requireLower, requireDigit, requireSymbol, minClasses)` | required character classes, and/or a minimum count of distinct classes present. The symbol class is any non-letter, non-digit, non-space rune; caseless letters such as CJK count toward no class (v0.38 fix) |
| `NewDisallowIdentityRule()` | rejects a password containing the principal's `ID` or `Name` (case-insensitive; tokens shorter than 3 runes are ignored, counted in runes so two-character CJK names do not reject most passwords) |
| `NewBlocklistRule(entries)` | rejects passwords matching a deny list (case-insensitive, trimmed) |

### Configuration

The framework builds a `PasswordValidator` from `vef.security.password_policy`
automatically and injects it wherever the framework needs one (the
forced-password-change challenge, described below). Every field is opt-in — a
zero value disables that rule:

```toml
[vef.security.password_policy]
min_length = 12
max_length = 128
require_upper = true
require_lower = true
require_digit = true
require_symbol = false
min_char_classes = 3
disallow_username = true
blocklist = ["password", "123456", "qwerty"]
```

Inject the same `security.PasswordValidator` into your own registration or
reset flows to reuse the configured policy instead of re-declaring rules.

### Violations

Every strength violation carries business code
`security.ErrCodePasswordPolicyViolation` (`1050`), HTTP 400; the i18n
message identifies which rule broke:

| Error | Trigger |
| --- | --- |
| `ErrPasswordTooShort(minLength)` | below `min_length` |
| `ErrPasswordTooLong(maxLength)` | above `max_length` |
| `ErrPasswordMissingUppercase` / `ErrPasswordMissingLowercase` / `ErrPasswordMissingDigit` / `ErrPasswordMissingSymbol` | a required character class is absent |
| `ErrPasswordTooFewCharClasses(minClasses)` | fewer than `min_char_classes` distinct classes present |
| `ErrPasswordContainsIdentity` | password contains the account `ID` or `Name` |
| `ErrPasswordBlocked` | password matches a blocklist entry |

The three templated messages carry named i18n-key constants for callers that
need the raw key rather than a constructed `result.Error`:
`security.ErrMessagePasswordTooShort`, `security.ErrMessagePasswordTooLong`,
and `security.ErrMessagePasswordTooFewCharClasses`.

## Password History (Reuse Prevention)

Password history rejects a new password that repeats one of the subject's
recent passwords. The framework only reads history to check reuse and
performs the hash comparison itself — the application owns storing it, since
the user database belongs to the application:

```go
type PasswordHistoryStore interface {
	// Recent returns the subject's most recent encoded passwords, newest first,
	// capped at limit.
	Recent(ctx context.Context, principalID string, limit int) ([]string, error)
	// Add records encodedPassword as the subject's newest history entry.
	Add(ctx context.Context, principalID, encodedPassword string) error
}
```

Register your implementation as a plain DI value:

```go
fx.Provide(func(db orm.DB) security.PasswordHistoryStore {
	return myapp.NewPasswordHistoryStore(db)
})
```

Set `vef.security.password_policy.history_depth` to a positive number:

```toml
[vef.security.password_policy]
history_depth = 5
```

When both a `PasswordHistoryStore` is registered *and* `history_depth > 0`,
the framework composes a history check into the injected `PasswordValidator`
via `NewChainValidator` — strength rules run first, then reuse. A password
matching any of the last `history_depth` entries fails with
`security.ErrPasswordReused` (same `ErrCodePasswordPolicyViolation` / 400).

Call `Add` from your own `PasswordChanger.ChangePassword` implementation
right after persisting the new hash — the framework reads and compares
history, the application writes to it. `NewHistoryValidator` can also be used
directly (`NewHistoryValidator(store, encoder, depth)`) if you assemble your
own validator chain instead of relying on the config-driven default.

## Password Expiry

Password expiry forces a change once a password is older than a configured
maximum age. The framework needs to know when a password was last set, which
it doesn't own, so the application implements a loader:

```go
type PasswordMetadataLoader interface {
	// PasswordChangedAt returns when the principal's password was last set. A
	// zero time means "unknown", treated as not-yet-expired rather than forcing
	// a change on incomplete data.
	PasswordChangedAt(ctx context.Context, principal *Principal) (time.Time, error)
}
```

Wrap it in an `ExpiryPasswordChangeChecker`:

```go
checker := security.NewExpiryPasswordChangeChecker(myMetadataLoader, 90*24*time.Hour)
```

The conventional place to declare the age limit is
`vef.security.password_policy.max_age` (zero disables expiry) — but note the
framework does **not** consume this key itself. Unlike `history_depth`
(auto-wired once a `PasswordHistoryStore` is registered), `max_age` is purely
declarative: your wiring code reads it from `config.SecurityConfig` and
passes it to `NewExpiryPasswordChangeChecker`, as in the example above.

```toml
[vef.security.password_policy]
max_age = "2160h" # 90 days
```

`ExpiryPasswordChangeChecker` implements `security.PasswordChangeChecker`, the
same interface used for other forced-change reasons (e.g. first login).
Combine several with `NewCompositePasswordChangeChecker`, which returns the
first reason that applies:

```go
checker := security.NewCompositePasswordChangeChecker(
	firstLoginChecker,
	security.NewExpiryPasswordChangeChecker(myMetadataLoader, 90*24*time.Hour),
)
```

Wire the composed checker, your `PasswordChanger`, and (optionally) a
`PasswordValidator` into `NewPasswordChangeChallengeProvider`, then register
it as a login challenge provider:

```go
vef.ProvideChallengeProvider(func(
	checker security.PasswordChangeChecker,
	changer security.PasswordChanger,
	validator security.PasswordValidator,
) security.ChallengeProvider {
	return security.NewPasswordChangeChallengeProvider(checker, changer, validator)
})
```

When the checker fires, `security/auth.login` returns a `password_change`
challenge (`Reason: "expired"`) instead of tokens; the client resolves it via
`resolve_challenge` with the new password as the response. The provider
validates the new password through the supplied `PasswordValidator` — so
strength and history rules apply here too — before calling
`PasswordChanger.ChangePassword` to persist it. See
[Authentication](./authentication) for the full login/challenge request
shape.

## Putting It Together

A deployment using every layer combines the config sections above under one
`[vef.security]` table:

```toml
[vef.security.lockout]
max_failures = 10
strategy = "lock"
key = "user_ip"

[vef.security.password_policy]
min_length = 12
require_upper = true
require_lower = true
require_digit = true
min_char_classes = 3
disallow_username = true
history_depth = 5
max_age = "2160h"
```

...plus, on the Go side, whatever of `PasswordDecryptor`,
`PasswordHistoryStore`, and `PasswordMetadataLoader` your threat model calls
for — each is independent, and none is required for the others to work.
