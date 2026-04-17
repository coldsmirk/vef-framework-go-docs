---
sidebar_position: 10
---

# Cryptox

The `cryptox` package provides a unified interface for encryption/decryption and digital signing across multiple algorithms.

## Interfaces

### Cipher

For encryption and decryption:

```go
type Cipher interface {
    Encrypt(plaintext string) (string, error)
    Decrypt(ciphertext string) (string, error)
}
```

### Signer

For signing and verification:

```go
type Signer interface {
    Sign(data string) (signature string, err error)
    Verify(data, signature string) (bool, error)
}
```

### CipherSigner

Combined encryption and signing:

```go
type CipherSigner interface {
    Cipher
    Signer
}
```

## Supported Algorithms

### AES (Symmetric Encryption)

```go
import "github.com/coldsmirk/vef-framework-go/cryptox"

cipher, err := cryptox.NewAESCipher(key) // key: 16, 24, or 32 bytes

encrypted, err := cipher.Encrypt("hello world")
plaintext, err := cipher.Decrypt(encrypted)
```

### RSA (Asymmetric Encryption + Signing)

```go
// From PEM-encoded keys
cipher, err := cryptox.NewRSACipher(publicKeyPEM, privateKeyPEM)

// Encrypt / Decrypt
encrypted, err := cipher.Encrypt("sensitive data")
plaintext, err := cipher.Decrypt(encrypted)

// Sign / Verify
signature, err := cipher.Sign("important message")
valid, err := cipher.Verify("important message", signature)
```

### SM2 (Chinese National Standard — Asymmetric)

```go
cipher, err := cryptox.NewSM2Cipher(publicKeyPEM, privateKeyPEM)

encrypted, err := cipher.Encrypt("data")
plaintext, err := cipher.Decrypt(encrypted)

signature, err := cipher.Sign("data")
valid, err := cipher.Verify("data", signature)
```

### SM4 (Chinese National Standard — Symmetric)

```go
cipher, err := cryptox.NewSM4Cipher(key) // key: 16 bytes

encrypted, err := cipher.Encrypt("data")
plaintext, err := cipher.Decrypt(encrypted)
```

### ECDSA (Signing Only)

```go
signer, err := cryptox.NewECDSACipher(publicKeyPEM, privateKeyPEM)

signature, err := signer.Sign("data to sign")
valid, err := signer.Verify("data to sign", signature)
```

### ECIES (Encryption Only)

```go
cipher, err := cryptox.NewECIESCipher(publicKeyPEM, privateKeyPEM)

encrypted, err := cipher.Encrypt("secret data")
plaintext, err := cipher.Decrypt(encrypted)
```

## Algorithm Comparison

| Algorithm | Type | Encrypt | Sign | Standard |
| --- | --- | --- | --- | --- |
| AES | Symmetric | ✅ | ❌ | International |
| RSA | Asymmetric | ✅ | ✅ | International |
| ECDSA | Asymmetric | ❌ | ✅ | International |
| ECIES | Asymmetric | ✅ | ❌ | International |
| SM2 | Asymmetric | ✅ | ✅ | Chinese (国密) |
| SM4 | Symmetric | ✅ | ❌ | Chinese (国密) |
