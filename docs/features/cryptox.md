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

All constructors return the framework's `Cipher` / `Signer` / `CipherSigner` interface. Every algorithm also exposes `*FromPem` / `*FromHex` / `*FromBase64` variants so you can feed keys in whatever encoding you have on hand.

### AES (Symmetric Encryption)

```go
import "github.com/coldsmirk/vef-framework-go/cryptox"

// key must be 16 / 24 / 32 bytes. Default mode is GCM (authenticated; IV is generated per call).
cipher, err := cryptox.NewAES(key)

// Use CBC instead — IV must be supplied explicitly.
cbcCipher, err := cryptox.NewAES(key,
    cryptox.WithAESMode(cryptox.AESModeCBC),
    cryptox.WithAESIv(iv), // 16 bytes
)

encrypted, err := cipher.Encrypt("hello world")
plaintext, err := cipher.Decrypt(encrypted)
```

Variants: `cryptox.NewAESFromHex(keyHex, ...)`, `cryptox.NewAESFromBase64(keyBase64, ...)`.

### RSA (Asymmetric Encryption + Signing)

```go
// Private key comes first.
cipher, err := cryptox.NewRSAFromPem(privatePEM, publicPEM)

encrypted, err := cipher.Encrypt("sensitive data")
plaintext, err := cipher.Decrypt(encrypted)

signature, err := cipher.Sign("important message")
valid, err := cipher.Verify("important message", signature)
```

Variants: `cryptox.NewRSA(privateKey, publicKey)` (from `*rsa.PrivateKey` / `*rsa.PublicKey`), `cryptox.NewRSAFromHex`, `cryptox.NewRSAFromBase64`.

### SM2 (Chinese National Standard — Asymmetric)

```go
// Private key comes first.
cipher, err := cryptox.NewSM2FromPEM(privatePEM, publicPEM)

encrypted, err := cipher.Encrypt("data")
plaintext, err := cipher.Decrypt(encrypted)

signature, err := cipher.Sign("data")
valid, err := cipher.Verify("data", signature)
```

Variants: `cryptox.NewSM2(privateKey, publicKey)`, `cryptox.NewSM2FromHex`, `cryptox.NewSM2FromBase64`.

### SM4 (Chinese National Standard — Symmetric)

```go
// Default mode is CBC — IV is REQUIRED. key: 16 bytes.
cipher, err := cryptox.NewSM4(key, cryptox.WithSM4Iv(iv))

// Or switch to ECB (no IV needed, less secure):
ecbCipher, err := cryptox.NewSM4(key, cryptox.WithSM4Mode(cryptox.SM4ModeECB))

encrypted, err := cipher.Encrypt("data")
plaintext, err := cipher.Decrypt(encrypted)
```

Variants: `cryptox.NewSM4FromHex`, `cryptox.NewSM4FromBase64`.

### ECDSA (Signing Only)

```go
// Private key comes first.
signer, err := cryptox.NewECDSAFromPem(privatePEM, publicPEM)

signature, err := signer.Sign("data to sign")
valid, err := signer.Verify("data to sign", signature)
```

Variants: `cryptox.NewECDSA(privateKey, publicKey)`, `cryptox.NewECDSAFromHex`, `cryptox.NewECDSAFromBase64`.

### ECIES (Encryption Only)

ECIES additionally requires the curve identifier:

```go
// Private key bytes first, then public key bytes, then the curve.
cipher, err := cryptox.NewECIESFromBytes(privateKeyBytes, publicKeyBytes, cryptox.EciesCurveP256)

encrypted, err := cipher.Encrypt("secret data")
plaintext, err := cipher.Decrypt(encrypted)
```

Variants: `cryptox.NewECIES(privateKey, publicKey)` (from `*ecdh.PrivateKey` / `*ecdh.PublicKey`), `cryptox.NewECIESFromHex(..., curve)`, `cryptox.NewECIESFromBase64(..., curve)`. Supported curves: `EciesCurveP256`, `EciesCurveP384`, `EciesCurveP521`, `EciesCurveX25519`.

## Algorithm Comparison

| Algorithm | Type | Encrypt | Sign | Standard |
| --- | --- | --- | --- | --- |
| AES | Symmetric | ✅ | ❌ | International |
| RSA | Asymmetric | ✅ | ✅ | International |
| ECDSA | Asymmetric | ❌ | ✅ | International |
| ECIES | Asymmetric | ✅ | ❌ | International |
| SM2 | Asymmetric | ✅ | ✅ | Chinese (国密) |
| SM4 | Symmetric | ✅ | ❌ | Chinese (国密) |
