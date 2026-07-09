---
sidebar_position: 6
---

# Password

`password` 包提供可插拔的密码编码，通过组合编码器同时支持多种算法。

## 编码器接口

每个算法专用的 encoder 以及组合编码器都实现同一个 `Encoder` 接口：

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
ok = encoder.Matches("my-password", "{argon2}$argon2id$...")  // true
ok = encoder.Matches("my-password", "{sha256}abc123...")       // true
```

如果 encoded password 没有 `{prefix}`，`CompositeEncoder` 会在匹配和升级检查时
fallback 到默认 encoder ID。未知前缀会让 `Matches` 返回 `false`。

### 升级检测

检查密码是否需要重新编码（例如从 SHA256 迁移到 bcrypt）：

```go
needsUpgrade := encoder.UpgradeEncoding("{sha256}abc123...")
// → true（因为默认是 bcrypt，不是 sha256）
```

`UpgradeEncoding` 遇到有效但非默认的前缀时会直接返回 `true`。默认前缀或无前缀
值会继续委托给默认 encoder 自身的 `UpgradeEncoding` 逻辑，例如 bcrypt cost
比较。

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

公开的 encoder ID 常量是 `EncoderBcrypt`、`EncoderArgon2`、
`EncoderScrypt`、`EncoderPbkdf2`、`EncoderMd5`、`EncoderSha256` 和
`EncoderPlaintext`。

`{prefix}` 字符串由 `EncoderID` 常量值决定：`EncoderBcrypt` → `{bcrypt}`、
`EncoderArgon2` → `{argon2}`、`EncoderScrypt` → `{scrypt}`、
`EncoderPbkdf2` → `{pbkdf2}`、`EncoderMd5` → `{md5}`、
`EncoderSha256` → `{sha256}`、`EncoderPlaintext` → `{plaintext}`。

## 编码器选项

| 编码器 | Option 函数 |
| --- | --- |
| `NewBcryptEncoder` | `WithBcryptCost(cost)` |
| `NewArgon2Encoder` | `WithArgon2Memory(memory)`, `WithArgon2Iterations(iterations)`, `WithArgon2Parallelism(parallelism)` |
| `NewScryptEncoder` | `WithScryptN(n)`, `WithScryptR(r)`, `WithScryptP(p)` |
| `NewPbkdf2Encoder` | `WithPbkdf2Iterations(iterations)`, `WithPbkdf2HashFunction(hashFunction)` |
| `NewMd5Encoder` | `WithMd5Salt(salt)`, `WithMd5SaltPosition(position)` |
| `NewSha256Encoder` | `WithSha256Salt(salt)`, `WithSha256SaltPosition(position)` |

对应的 option 类型也公开为 `BcryptOption`、`Argon2Option`、
`ScryptOption`、`Pbkdf2Option`、`Md5Option` 和 `Sha256Option`。

公开的 option 函数是 `WithBcryptCost`、`WithArgon2Memory`、
`WithArgon2Iterations`、`WithArgon2Parallelism`、`WithScryptN`、
`WithScryptR`、`WithScryptP`、`WithPbkdf2Iterations`、
`WithPbkdf2HashFunction`、`WithMd5Salt`、`WithMd5SaltPosition`、
`WithSha256Salt` 和 `WithSha256SaltPosition`。

默认参数：

| Encoder | 默认值 |
| --- | --- |
| bcrypt | `bcrypt.DefaultCost`（`10`）；合法 cost 范围是 `4..31` |
| Argon2id | memory `64 * 1024` KiB、iterations `3`、parallelism `4` |
| scrypt | `N = 32768`、`r = 8`、`p = 1` |
| PBKDF2 | `310000` 次迭代，默认 `sha256`；也支持 `sha512` |
| MD5 / SHA-256 | 默认无 salt；配置 salt 后默认位置是 `suffix`，也支持 `prefix` |

## 密码传输加密

`password.Encoder` 始终是一个纯 KDF，同时用于存储和比对——不存在给 encoder
套一层 cipher 的包装器。如果客户端在发送前先加密了密码（常见的"前端加密、
后端校验"模式），解密发生在更上层的 authenticator：注册一个
`security.PasswordDecryptor`，`PasswordAuthenticator` 会先把收到的凭据解密为
明文，再拿它与存储的哈希比对。完整配置见
[登录加固](./login-hardening#密码传输加密)。

## 密码格式

编码后的密码格式：`{算法}编码值`；源码中的占位描述写作
`{algorithm}encoded_value`。

```
{bcrypt}$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
{argon2}$argon2id$v=19$m=65536,t=3,p=4$c29tZXNhbHQ$...
{sha256}5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8
```

独立的 MD5 和 SHA-256 encoder 在没有 salt 时返回原始 hex digest。配置 salt 后，
内部 hash 格式是 `{algorithm}$salt$hash`；如果再由 `NewCompositeEncoder` 包装，
composite 会在这个 encoded value 前面追加外层 `{algorithm}` 前缀。

## 错误哨兵

| 错误 | 出现时机 |
| --- | --- |
| `ErrInvalidCost` | bcrypt cost 不在 `4..31` 范围内 |
| `ErrInvalidMemory` | Argon2 memory 参数过小 |
| `ErrInvalidIterations` | iteration count 小于 `1` |
| `ErrInvalidParallelism` | Argon2 parallelism 小于 `1` |
| `ErrInvalidEncoderID` | `CompositeEncoder` 找不到前缀对应的 encoder |
| `ErrInvalidHashFormat` | 编码后的密码格式不合法 |
| `ErrDefaultEncoderNotFound` | `CompositeEncoder` 没有注册默认 encoder ID |
