---
sidebar_position: 11
---

# Password

`password` 包提供可插拔的密码编码，通过组合编码器同时支持多种算法。

## 编码器接口

```go
type Encoder interface {
    Encode(password string) (string, error)
    Matches(password, encodedPassword string) bool
    UpgradeEncoding(encodedPassword string) bool
}
```

## 组合编码器

组合编码器将算法标识作为前缀存储在编码后的密码中，实现无缝的算法迁移：

```go
import "github.com/coldsmirk/vef-framework-go/password"

encoder := password.NewCompositeEncoder(
    password.EncoderBcrypt, // 新密码默认使用
    map[password.EncoderID]password.Encoder{
        password.EncoderBcrypt: password.NewBcryptEncoder(),
        password.EncoderArgon2: password.NewArgon2Encoder(),
        password.EncoderScrypt: password.NewScryptEncoder(),
        password.EncoderSha256: password.NewSha256Encoder(),
    },
)
```

### 编码

```go
encoded, err := encoder.Encode("my-password")
// → "{bcrypt}$2a$10$..."
```

### 匹配

编码器自动从 `{前缀}` 检测算法：

```go
// 支持所有已注册的算法
ok := encoder.Matches("my-password", "{bcrypt}$2a$10$...")     // true
ok := encoder.Matches("my-password", "{argon2}$argon2id$...")  // true
ok := encoder.Matches("my-password", "{sha256}abc123...")       // true
```

### 升级检测

检查密码是否需要重新编码（例如从 SHA256 迁移到 bcrypt）：

```go
needsUpgrade := encoder.UpgradeEncoding("{sha256}abc123...")
// → true（因为默认是 bcrypt，不是 sha256）
```

## 可用编码器

| 编码器 ID | 构造函数 | 安全级别 |
| --- | --- | --- |
| `password.EncoderBcrypt` | `password.NewBcryptEncoder()` | ⭐⭐⭐⭐ 推荐 |
| `password.EncoderArgon2` | `password.NewArgon2Encoder()` | ⭐⭐⭐⭐⭐ 最强 |
| `password.EncoderScrypt` | `password.NewScryptEncoder()` | ⭐⭐⭐⭐ 强 |
| `password.EncoderPbkdf2` | `password.NewPbkdf2Encoder()` | ⭐⭐⭐ 标准（FIPS 友好） |
| `password.EncoderSha256` | `password.NewSha256Encoder()` | ⭐⭐ 仅用于遗留系统 |
| `password.EncoderMd5` | `password.NewMd5Encoder()` | ⭐ 仅用于遗留 / 互通 |
| `password.EncoderPlaintext` | `password.NewPlaintextEncoder()` | ⭐ 仅用于测试 |

`{前缀}` 字符串由 `EncoderID` 常量值决定，例如 `EncoderBcrypt` → `{bcrypt}`，`EncoderSha256` → `{sha256}`。

## 在编码器外再套 Cipher

如果客户端在传输前先加密 / 哈希密码（常见的"前端加密、后端再哈希"模式），用 `NewCipherEncoder` 包一层：

```go
inner := password.NewBcryptEncoder()
encoder := password.NewCipherEncoder(rsaCipher, inner)

// 收到客户端发来的（已经被 rsaCipher 加密的）密文。
// 编码器先用 rsaCipher 解开，再交给 inner 哈希。
hashed, err := encoder.Encode(encoded)
```

## 密码格式

编码后的密码格式：`{算法}编码值`

```
{bcrypt}$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
{argon2}$argon2id$v=19$m=65536,t=3,p=4$c29tZXNhbHQ$...
{sha256}5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8
```
