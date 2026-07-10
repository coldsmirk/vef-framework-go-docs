---
sidebar_position: 5
---

# Password

The `password` package provides pluggable password encoding with a composite encoder that supports multiple algorithms simultaneously.

## Encoder Interface

Every algorithm-specific encoder and the composite encoder implement the same
`Encoder` interface:

```go
type Encoder interface {
    Encode(password string) (string, error)
    Matches(password, encodedPassword string) bool
    UpgradeEncoding(encodedPassword string) bool
}
```

## Composite Encoder

The composite encoder stores the algorithm identifier as a prefix in the encoded password, enabling seamless algorithm migration:

```go
import "github.com/coldsmirk/vef-framework-go/password"

encoder := password.NewCompositeEncoder(
    password.EncoderBcrypt, // default for new passwords
    map[password.EncoderID]password.Encoder{
        password.EncoderBcrypt: password.NewBcryptEncoder(),
        password.EncoderArgon2: password.NewArgon2Encoder(),
        password.EncoderScrypt: password.NewScryptEncoder(),
        password.EncoderSha256: password.NewSha256Encoder(),
    },
)
```

### Encoding

```go
encoded, err := encoder.Encode("my-password")
// → "{bcrypt}$2a$10$..."
```

### Matching

The encoder automatically detects the algorithm from the `{prefix}`:

```go
// Matches against any supported algorithm
ok := encoder.Matches("my-password", "{bcrypt}$2a$10$...")   // true
ok = encoder.Matches("my-password", "{argon2}$argon2id$...") // true
ok = encoder.Matches("my-password", "{sha256}abc123...")      // true
```

If an encoded password has no `{prefix}`, `CompositeEncoder` falls back to the
default encoder ID for matching and upgrade checks. Unknown prefixes return
`false` from `Matches`.

### Upgrade Detection

Check if a password needs re-encoding (e.g., when migrating from SHA256 to bcrypt):

```go
needsUpgrade := encoder.UpgradeEncoding("{sha256}abc123...")
// → true (because default is bcrypt, not sha256)
```

`UpgradeEncoding` returns `true` immediately for any valid non-default prefix.
For the default prefix or a no-prefix value, it delegates to the default
encoder's own `UpgradeEncoding` logic, such as bcrypt cost comparison.

## Available Encoders

| Encoder ID | Constructor | Security Level |
| --- | --- | --- |
| `password.EncoderBcrypt` | `password.NewBcryptEncoder()` | ⭐⭐⭐⭐ Recommended |
| `password.EncoderArgon2` | `password.NewArgon2Encoder()` | ⭐⭐⭐⭐⭐ Strongest |
| `password.EncoderScrypt` | `password.NewScryptEncoder()` | ⭐⭐⭐⭐ Strong |
| `password.EncoderPbkdf2` | `password.NewPbkdf2Encoder()` | ⭐⭐⭐ Standard (FIPS-friendly) |
| `password.EncoderSha256` | `password.NewSha256Encoder()` | ⭐⭐ Legacy only |
| `password.EncoderMd5` | `password.NewMd5Encoder()` | ⭐ Legacy / interop only |
| `password.EncoderPlaintext` | `password.NewPlaintextEncoder()` | ⭐ Testing only |

The public encoder ID constants are `EncoderBcrypt`, `EncoderArgon2`,
`EncoderScrypt`, `EncoderPbkdf2`, `EncoderMd5`, `EncoderSha256`, and
`EncoderPlaintext`.

The `{prefix}` segment is derived from the `EncoderID` constant value:
`EncoderBcrypt` → `{bcrypt}`, `EncoderArgon2` → `{argon2}`,
`EncoderScrypt` → `{scrypt}`, `EncoderPbkdf2` → `{pbkdf2}`,
`EncoderMd5` → `{md5}`, `EncoderSha256` → `{sha256}`, and
`EncoderPlaintext` → `{plaintext}`.

## Encoder Options

| Encoder | Option functions |
| --- | --- |
| `NewBcryptEncoder` | `WithBcryptCost(cost)` |
| `NewArgon2Encoder` | `WithArgon2Memory(memory)`, `WithArgon2Iterations(iterations)`, `WithArgon2Parallelism(parallelism)` |
| `NewScryptEncoder` | `WithScryptN(n)`, `WithScryptR(r)`, `WithScryptP(p)` |
| `NewPbkdf2Encoder` | `WithPbkdf2Iterations(iterations)`, `WithPbkdf2HashFunction(hashFunction)` |
| `NewMd5Encoder` | `WithMd5Salt(salt)`, `WithMd5SaltPosition(position)` |
| `NewSha256Encoder` | `WithSha256Salt(salt)`, `WithSha256SaltPosition(position)` |

The corresponding option types are exported as `BcryptOption`, `Argon2Option`,
`ScryptOption`, `Pbkdf2Option`, `Md5Option`, and `Sha256Option`.

The exported option functions are `WithBcryptCost`, `WithArgon2Memory`,
`WithArgon2Iterations`, `WithArgon2Parallelism`, `WithScryptN`,
`WithScryptR`, `WithScryptP`, `WithPbkdf2Iterations`,
`WithPbkdf2HashFunction`, `WithMd5Salt`, `WithMd5SaltPosition`,
`WithSha256Salt`, and `WithSha256SaltPosition`.

Default parameters:

| Encoder | Defaults |
| --- | --- |
| bcrypt | `bcrypt.DefaultCost` (`10`); valid cost range is `4..31` |
| Argon2id | memory `64 * 1024` KiB, iterations `3`, parallelism `4` |
| scrypt | `N = 32768`, `r = 8`, `p = 1` |
| PBKDF2 | `310000` iterations with `sha256`; `sha512` is also supported |
| MD5 / SHA-256 | no salt by default; with salt, default salt position is `suffix` and `prefix` is available |

## Encrypted Password Transport

`password.Encoder` stays a plain KDF used for both storage and comparison —
there is no cipher-wrapping encoder. If the client encrypts the password
before sending it (a common "front-end encrypts, backend verifies" pattern),
decryption happens one layer up, at the authenticator: register a
`security.PasswordDecryptor` and `PasswordAuthenticator` decrypts the
transmitted credential to plaintext before verifying it against the stored
hash. See [Login Hardening](./login-hardening#encrypted-password-transport)
for the full setup.

## Password Format

Encoded passwords follow the format: `{algorithm}encoded_value`

```
{bcrypt}$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
{argon2}$argon2id$v=19$m=65536,t=3,p=4$c29tZXNhbHQ$...
{sha256}5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8
```

Standalone MD5 and SHA-256 encoders without salt return the raw hex digest.
When configured with salt, they use the inner hash format
`{algorithm}$salt$hash`; when wrapped by `NewCompositeEncoder`, the composite
adds the outer `{algorithm}` prefix before that encoded value.

## Error Sentinels

| Error | When it appears |
| --- | --- |
| `ErrInvalidCost` | bcrypt cost is outside `4..31` |
| `ErrInvalidMemory` | Argon2 memory is too small |
| `ErrInvalidIterations` | iteration count is lower than `1` |
| `ErrInvalidParallelism` | Argon2 parallelism is lower than `1` |
| `ErrInvalidEncoderID` | `CompositeEncoder` cannot find an encoder for the prefix |
| `ErrInvalidHashFormat` | encoded password format is malformed |
| `ErrDefaultEncoderNotFound` | default encoder ID is not registered in `CompositeEncoder` |
