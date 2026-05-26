---
sidebar_position: 10
---

# Cryptox

`cryptox` 包提供统一的加密/解密和数字签名接口，支持多种算法。

## 接口

### Cipher

用于加密和解密：

```go
type Cipher interface {
    Encrypt(plaintext string) (string, error)
    Decrypt(ciphertext string) (string, error)
}
```

### Signer

用于签名和验证：

```go
type Signer interface {
    Sign(data string) (signature string, err error)
    Verify(data, signature string) (bool, error)
}
```

### CipherSigner

组合加密和签名：

```go
type CipherSigner interface {
    Cipher
    Signer
}
```

## 支持的算法

所有构造函数返回框架的 `Cipher` / `Signer` / `CipherSigner` 接口。每个算法都另外提供 `*FromPem` / `*FromHex` / `*FromBase64` 三种变体，按需选用。

### AES（对称加密）

```go
import "github.com/coldsmirk/vef-framework-go/cryptox"

// key 必须是 16 / 24 / 32 字节，默认 GCM 模式（带认证，IV 自动生成）。
cipher, err := cryptox.NewAES(key)

// 切到 CBC 模式必须显式提供 IV：
cbcCipher, err := cryptox.NewAES(key,
    cryptox.WithAESMode(cryptox.AESModeCBC),
    cryptox.WithAESIv(iv), // 16 字节
)

encrypted, err := cipher.Encrypt("hello world")
plaintext, err := cipher.Decrypt(encrypted)
```

变体：`cryptox.NewAESFromHex(keyHex, ...)`、`cryptox.NewAESFromBase64(keyBase64, ...)`。

### RSA（非对称加密 + 签名）

```go
// 私钥在前。
cipher, err := cryptox.NewRSAFromPem(privatePEM, publicPEM)

encrypted, err := cipher.Encrypt("sensitive data")
plaintext, err := cipher.Decrypt(encrypted)

signature, err := cipher.Sign("important message")
valid, err := cipher.Verify("important message", signature)
```

变体：`cryptox.NewRSA(privateKey, publicKey)`（直接接收 `*rsa.PrivateKey` / `*rsa.PublicKey`）、`cryptox.NewRSAFromHex`、`cryptox.NewRSAFromBase64`。

### SM2（国密 — 非对称）

```go
// 私钥在前。
cipher, err := cryptox.NewSM2FromPEM(privatePEM, publicPEM)

encrypted, err := cipher.Encrypt("data")
plaintext, err := cipher.Decrypt(encrypted)

signature, err := cipher.Sign("data")
valid, err := cipher.Verify("data", signature)
```

变体：`cryptox.NewSM2(privateKey, publicKey)`、`cryptox.NewSM2FromHex`、`cryptox.NewSM2FromBase64`。

### SM4（国密 — 对称）

```go
// 默认 CBC 模式，必须传 IV。key：16 字节。
cipher, err := cryptox.NewSM4(key, cryptox.WithSM4Iv(iv))

// 或切到 ECB（不需要 IV，但更不安全）：
ecbCipher, err := cryptox.NewSM4(key, cryptox.WithSM4Mode(cryptox.SM4ModeECB))

encrypted, err := cipher.Encrypt("data")
plaintext, err := cipher.Decrypt(encrypted)
```

变体：`cryptox.NewSM4FromHex`、`cryptox.NewSM4FromBase64`。

### ECDSA（仅签名）

```go
// 私钥在前。
signer, err := cryptox.NewECDSAFromPem(privatePEM, publicPEM)

signature, err := signer.Sign("data to sign")
valid, err := signer.Verify("data to sign", signature)
```

变体：`cryptox.NewECDSA(privateKey, publicKey)`、`cryptox.NewECDSAFromHex`、`cryptox.NewECDSAFromBase64`。

### ECIES（仅加密）

ECIES 额外需要曲线参数：

```go
// 私钥字节在前，公钥字节在后，再传曲线。
cipher, err := cryptox.NewECIESFromBytes(privateKeyBytes, publicKeyBytes, cryptox.EciesCurveP256)

encrypted, err := cipher.Encrypt("secret data")
plaintext, err := cipher.Decrypt(encrypted)
```

变体：`cryptox.NewECIES(privateKey, publicKey)`（接收 `*ecdh.PrivateKey` / `*ecdh.PublicKey`）、`cryptox.NewECIESFromHex(..., curve)`、`cryptox.NewECIESFromBase64(..., curve)`。支持的曲线：`EciesCurveP256`、`EciesCurveP384`、`EciesCurveP521`、`EciesCurveX25519`。

## 算法对比

| 算法 | 类型 | 加密 | 签名 | 标准 |
| --- | --- | --- | --- | --- |
| AES | 对称 | ✅ | ❌ | 国际 |
| RSA | 非对称 | ✅ | ✅ | 国际 |
| ECDSA | 非对称 | ❌ | ✅ | 国际 |
| ECIES | 非对称 | ✅ | ❌ | 国际 |
| SM2 | 非对称 | ✅ | ✅ | 国密 |
| SM4 | 对称 | ✅ | ❌ | 国密 |
