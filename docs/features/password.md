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
    password.Bcrypt, // default for new passwords
    map[password.EncoderID]password.Encoder{
        password.Bcrypt:  password.NewBcryptEncoder(),
        password.Argon2:  password.NewArgon2Encoder(),
        password.Scrypt:  password.NewScryptEncoder(),
        password.SHA256:  password.NewSHA256Encoder(),
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
| `password.Bcrypt` | `NewBcryptEncoder()` | ⭐⭐⭐⭐ Recommended |
| `password.Argon2` | `NewArgon2Encoder()` | ⭐⭐⭐⭐⭐ Strongest |
| `password.Scrypt` | `NewScryptEncoder()` | ⭐⭐⭐⭐ Strong |
| `password.SHA256` | `NewSHA256Encoder()` | ⭐⭐ Legacy only |
| `password.Plaintext` | `NewPlaintextEncoder()` | ⭐ Testing only |

## Password Format

Encoded passwords follow the format: `{algorithm}encoded_value`

```
{bcrypt}$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
{argon2}$argon2id$v=19$m=65536,t=3,p=4$c29tZXNhbHQ$...
{sha256}5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8
```
