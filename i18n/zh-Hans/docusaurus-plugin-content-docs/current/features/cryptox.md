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

### AES（对称加密）

```go
cipher, err := cryptox.NewAESCipher(key) // key: 16、24 或 32 字节
encrypted, err := cipher.Encrypt("hello world")
plaintext, err := cipher.Decrypt(encrypted)
```

### RSA（非对称加密 + 签名）

```go
cipher, err := cryptox.NewRSACipher(publicKeyPEM, privateKeyPEM)
encrypted, err := cipher.Encrypt("sensitive data")
plaintext, err := cipher.Decrypt(encrypted)
signature, err := cipher.Sign("important message")
valid, err := cipher.Verify("important message", signature)
```

### SM2（国密 — 非对称）

```go
cipher, err := cryptox.NewSM2Cipher(publicKeyPEM, privateKeyPEM)
encrypted, err := cipher.Encrypt("data")
signature, err := cipher.Sign("data")
valid, err := cipher.Verify("data", signature)
```

### SM4（国密 — 对称）

```go
cipher, err := cryptox.NewSM4Cipher(key) // key: 16 字节
encrypted, err := cipher.Encrypt("data")
plaintext, err := cipher.Decrypt(encrypted)
```

### ECDSA（仅签名）

```go
signer, err := cryptox.NewECDSACipher(publicKeyPEM, privateKeyPEM)
signature, err := signer.Sign("data to sign")
valid, err := signer.Verify("data to sign", signature)
```

### ECIES（仅加密）

```go
cipher, err := cryptox.NewECIESCipher(publicKeyPEM, privateKeyPEM)
encrypted, err := cipher.Encrypt("secret data")
plaintext, err := cipher.Decrypt(encrypted)
```

## 算法对比

| 算法 | 类型 | 加密 | 签名 | 标准 |
| --- | --- | --- | --- | --- |
| AES | 对称 | ✅ | ❌ | 国际 |
| RSA | 非对称 | ✅ | ✅ | 国际 |
| ECDSA | 非对称 | ❌ | ✅ | 国际 |
| ECIES | 非对称 | ✅ | ❌ | 国际 |
| SM2 | 非对称 | ✅ | ✅ | 国密 |
| SM4 | 对称 | ✅ | ❌ | 国密 |
