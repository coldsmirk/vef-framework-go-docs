---
sidebar_position: 7
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

### FixedIVDecrypter

`FixedIVDecrypter` 由能够使用调用方固定 IV 解密外部 ciphertext 的 block-cipher
模式实现：

```go
type FixedIVDecrypter interface {
    DecryptWithFixedIV(ciphertext string) (string, error)
}
```

VEF 原生 ciphertext 会携带每次生成的新随机 IV，应使用 `Cipher.Decrypt`。
`DecryptWithFixedIV` 是 AES-CBC/SM4-CBC 互操作场景的 escape hatch：外部对端发送
未前置 IV 的 base64 ciphertext 时，解密会使用构造时通过 `WithAESIv` 或
`WithSM4Iv` 配置的固定 IV。

## 支持的算法

所有构造函数返回框架的 `Cipher` / `Signer` / `CipherSigner` 接口。编码辅助构造器按算法提供：AES 和 SM4 提供 `*FromHex` / `*FromBase64`；RSA、SM2、ECDSA 提供 `*FromPEM` / `*FromHex` / `*FromBase64`；ECIES 提供 bytes、hex 和 base64 构造器。

### AES（对称加密）

```go
import "github.com/coldsmirk/vef-framework-go/cryptox"

// key 必须是 16 / 24 / 32 字节，默认 GCM 模式（带认证，IV 自动生成）。
cipher, err := cryptox.NewAES(key)

// 切到 CBC 模式；Encrypt 会生成随机 IV 并前置到 ciphertext。
cbcCipher, err := cryptox.NewAES(key,
    cryptox.WithAESMode(cryptox.AesModeCbc),
)

encrypted, err := cipher.Encrypt("hello world")
plaintext, err := cipher.Decrypt(encrypted)
```

变体：`cryptox.NewAESFromHex(keyHex, ...)`、`cryptox.NewAESFromBase64(keyBase64, ...)`。

如果要和发送 bare ciphertext（未前置 IV）的 AES-CBC 外部系统互操作，
需要配置固定 IV，并调用 `FixedIVDecrypter.DecryptWithFixedIV`：

```go
fixedCipher, err := cryptox.NewAES(key,
    cryptox.WithAESMode(cryptox.AesModeCbc),
    cryptox.WithAESIv(iv), // 16 字节
)
plaintext, err := fixedCipher.(cryptox.FixedIVDecrypter).DecryptWithFixedIV(peerCiphertext)
```

### RSA（非对称加密 + 签名）

```go
// 私钥在前。
cipher, err := cryptox.NewRSAFromPEM(privatePEM, publicPEM)

encrypted, err := cipher.Encrypt("sensitive data")
plaintext, err := cipher.Decrypt(encrypted)

signature, err := cipher.Sign("important message")
valid, err := cipher.Verify("important message", signature)
```

变体：`cryptox.NewRSA(privateKey, publicKey)`（直接接收 `*rsa.PrivateKey` / `*rsa.PublicKey`）、`cryptox.NewRSAFromHex`、`cryptox.NewRSAFromBase64`。

RSA 默认加密模式是 `RsaModeOAEP`；默认签名模式是 `RsaSignModePSS`。
`RsaModePKCS1v15` 和 `RsaSignModePKCS1v15` 是显式选择的 legacy interop 模式。

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
// SM4 使用 CBC；Encrypt 会生成随机 IV 并前置到 ciphertext。
cipher, err := cryptox.NewSM4(key) // key：16 字节

encrypted, err := cipher.Encrypt("data")
plaintext, err := cipher.Decrypt(encrypted)
```

变体：`cryptox.NewSM4FromHex`、`cryptox.NewSM4FromBase64`。

如果要和发送 bare ciphertext（未前置 IV）的 SM4-CBC 外部系统互操作，
需要配置固定 IV，并调用 `FixedIVDecrypter.DecryptWithFixedIV`：

```go
fixedCipher, err := cryptox.NewSM4(key, cryptox.WithSM4Iv(iv))
plaintext, err := fixedCipher.(cryptox.FixedIVDecrypter).DecryptWithFixedIV(peerCiphertext)
```

### ECDSA（仅签名）

```go
// 私钥在前。
signer, err := cryptox.NewECDSAFromPEM(privatePEM, publicPEM)

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

`GenerateECIESKey(curve)` 和 ECIES byte parser 会使用传入的 `ECIESCurve`；
未知 curve 值会 fallback 到 P-256。`GenerateECDSAKey` 对 `ECDSACurve` 也是同样规则。

## 算法对比

| 算法 | 类型 | 加密 | 签名 | 标准 |
| --- | --- | --- | --- | --- |
| AES | 对称 | ✅ | ❌ | 国际 |
| RSA | 非对称 | ✅ | ✅ | 国际 |
| ECDSA | 非对称 | ❌ | ✅ | 国际 |
| ECIES | 非对称 | ✅ | ❌ | 国际 |
| SM2 | 非对称 | ✅ | ✅ | 国密 |
| SM4 | 对称 | ✅ | ❌ | 国密 |

## 选项、常量和密钥辅助函数

| 范围 | 公开 API |
| --- | --- |
| AES modes/options | `AESMode`, `AesModeGcm`, `AesModeCbc`, `WithAESMode(mode)`, `WithAESIv(iv)` |
| RSA modes/options | `RSAMode`, `RSASignMode`, `RsaModeOAEP`, `RsaModePKCS1v15`, `RsaSignModePSS`, `RsaSignModePKCS1v15`, `WithRSAMode(mode)`, `WithRSASignMode(mode)` |
| SM4 options | `WithSM4Iv(iv)` |
| ECDSA curves | `ECDSACurve`, `EcdsaCurveP224`, `EcdsaCurveP256`, `EcdsaCurveP384`, `EcdsaCurveP521` |
| ECIES curves | `ECIESCurve`, `EciesCurveP256`, `EciesCurveP384`, `EciesCurveP521`, `EciesCurveX25519` |
| key helpers | `GenerateECDSAKey(curve)`, `GenerateECIESKey(curve)` |
| option types | `AESOption`, `RSAOption`, `SM4Option` |

## 错误哨兵

| 分组 | 错误 |
| --- | --- |
| key availability | `ErrAtLeastOneKeyRequired`, `ErrPublicKeyRequiredForEncrypt`, `ErrPrivateKeyRequiredForDecrypt`, `ErrPrivateKeyRequiredForSign`, `ErrPublicKeyRequiredForVerify` |
| key parsing/type | `ErrFailedDecodePEMBlock`, `ErrUnsupportedPEMType`, `ErrNotRSAPrivateKey`, `ErrNotRSAPublicKey`, `ErrNotECDSAPrivateKey`, `ErrNotECDSAPublicKey` |
| symmetric crypto | `ErrInvalidAESKeySize`, `ErrInvalidSM4KeySize`, `ErrInvalidIVSizeCBC`, `ErrCiphertextNotMultipleOfBlock`, `ErrCiphertextTooShort`, `ErrInvalidPadding` |
| input/signature | `ErrDataEmpty`, `ErrInvalidSignature` |
