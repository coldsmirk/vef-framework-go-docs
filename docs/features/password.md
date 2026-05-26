---
sidebar_position: 11
---

# Password

The `password` package provides pluggable password encoding with a composite encoder that supports multiple algorithms simultaneously.

## Encoder Interface

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
ok := encoder.Matches("my-password", "{argon2}$argon2id$...") // true
ok := encoder.Matches("my-password", "{sha256}abc123...")      // true
```

### Upgrade Detection

Check if a password needs re-encoding (e.g., when migrating from SHA256 to bcrypt):

```go
needsUpgrade := encoder.UpgradeEncoding("{sha256}abc123...")
// → true (because default is bcrypt, not sha256)
```

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

The `{prefix}` segment is derived from the `EncoderID` constant value, so `EncoderBcrypt` → `{bcrypt}`, `EncoderSha256` → `{sha256}`, etc.

## Wrapping an Encoder With a Cipher

If the client hashes / encrypts the password before sending (a common "front-end encrypts, backend rehashes" pattern), wrap the underlying encoder with `NewCipherEncoder`:

```go
inner := password.NewBcryptEncoder()
encoder := password.NewCipherEncoder(rsaCipher, inner)

// `encoded` is what reaches the server (e.g. RSA-encrypted by the browser).
// The encoder decrypts it through `rsaCipher` first, then hashes the
// plaintext through `inner`.
hashed, err := encoder.Encode(encoded)
```

## Password Format

Encoded passwords follow the format: `{algorithm}encoded_value`

```
{bcrypt}$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
{argon2}$argon2id$v=19$m=65536,t=3,p=4$c29tZXNhbHQ$...
{sha256}5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8
```
