---
sidebar_position: 8
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

### FixedIVDecrypter

`FixedIVDecrypter` is implemented by block-cipher modes that can decrypt
external ciphertext produced with a caller-supplied constant IV:

```go
type FixedIVDecrypter interface {
    DecryptWithFixedIV(ciphertext string) (string, error)
}
```

Native VEF ciphertext carries a fresh random IV and should use
`Cipher.Decrypt`. `DecryptWithFixedIV` is an interop escape hatch for
AES-CBC/SM4-CBC peers that send base64 ciphertext without the prepended IV; the
fixed IV must have been configured with `WithAESIv` or `WithSM4Iv`.

## Supported Algorithms

All constructors return the framework's `Cipher` / `Signer` / `CipherSigner`
interface. Encoding helpers are algorithm-specific: AES and SM4 expose
`*FromHex` / `*FromBase64`, RSA/SM2/ECDSA expose `*FromPEM` / `*FromHex` /
`*FromBase64`, and ECIES exposes byte, hex, and base64 constructors.

### AES (Symmetric Encryption)

```go
import "github.com/coldsmirk/vef-framework-go/cryptox"

// key must be 16 / 24 / 32 bytes. Default mode is GCM (authenticated; IV is generated per call).
cipher, err := cryptox.NewAES(key)

// Use CBC instead. Encrypt generates a fresh IV and prepends it to ciphertext.
cbcCipher, err := cryptox.NewAES(key,
    cryptox.WithAESMode(cryptox.AesModeCbc),
)

encrypted, err := cipher.Encrypt("hello world")
plaintext, err := cipher.Decrypt(encrypted)
```

Variants: `cryptox.NewAESFromHex(keyHex, ...)`, `cryptox.NewAESFromBase64(keyBase64, ...)`.

For AES-CBC interop with peers that send bare ciphertext without a prepended IV,
configure the fixed IV and call `FixedIVDecrypter.DecryptWithFixedIV`:

```go
fixedCipher, err := cryptox.NewAES(key,
    cryptox.WithAESMode(cryptox.AesModeCbc),
    cryptox.WithAESIv(iv), // 16 bytes
)
plaintext, err := fixedCipher.(cryptox.FixedIVDecrypter).DecryptWithFixedIV(peerCiphertext)
```

### RSA (Asymmetric Encryption + Signing)

```go
// Private key comes first.
cipher, err := cryptox.NewRSAFromPEM(privatePEM, publicPEM)

encrypted, err := cipher.Encrypt("sensitive data")
plaintext, err := cipher.Decrypt(encrypted)

signature, err := cipher.Sign("important message")
valid, err := cipher.Verify("important message", signature)
```

Variants: `cryptox.NewRSA(privateKey, publicKey)` (from `*rsa.PrivateKey` / `*rsa.PublicKey`), `cryptox.NewRSAFromHex`, `cryptox.NewRSAFromBase64`.

Default RSA encryption mode is `RsaModeOAEP`; default signing mode is
`RsaSignModePSS`. `RsaModePKCS1v15` and `RsaSignModePKCS1v15` are explicit
legacy-interoperability choices.

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
// key: 16 bytes. Default mode is GCM (authenticated; IV is generated per call).
cipher, err := cryptox.NewSM4(key)

// Use CBC instead. Encrypt generates a fresh IV and prepends it to ciphertext.
cbcCipher, err := cryptox.NewSM4(key,
    cryptox.WithSM4Mode(cryptox.Sm4ModeCbc),
)

encrypted, err := cipher.Encrypt("data")
plaintext, err := cipher.Decrypt(encrypted)
```

Variants: `cryptox.NewSM4FromHex`, `cryptox.NewSM4FromBase64`.

:::caution SM4 defaults to GCM
`NewSM4` defaults to **GCM** (matching AES), not CBC. To decrypt SM4
ciphertext that was produced in CBC mode, construct the cipher with
`cryptox.WithSM4Mode(cryptox.Sm4ModeCbc)` explicitly.
:::

For SM4-CBC interop with peers that send bare ciphertext without a prepended IV,
configure the fixed IV and call `FixedIVDecrypter.DecryptWithFixedIV` (the fixed
IV only affects that interop decrypt path, never `Encrypt`):

```go
fixedCipher, err := cryptox.NewSM4(key,
    cryptox.WithSM4Mode(cryptox.Sm4ModeCbc),
    cryptox.WithSM4Iv(iv), // 16 bytes
)
plaintext, err := fixedCipher.(cryptox.FixedIVDecrypter).DecryptWithFixedIV(peerCiphertext)
```

### ECDSA (Signing Only)

```go
// Private key comes first.
signer, err := cryptox.NewECDSAFromPEM(privatePEM, publicPEM)

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

`GenerateECIESKey(curve)` and ECIES byte parsers use the requested
`ECIESCurve`; an unknown curve value falls back to P-256. `GenerateECDSAKey`
does the same for `ECDSACurve`.

## Algorithm Comparison

| Algorithm | Type | Encrypt | Sign | Standard |
| --- | --- | --- | --- | --- |
| AES | Symmetric | ✅ | ❌ | International |
| RSA | Asymmetric | ✅ | ✅ | International |
| ECDSA | Asymmetric | ❌ | ✅ | International |
| ECIES | Asymmetric | ✅ | ❌ | International |
| SM2 | Asymmetric | ✅ | ✅ | Chinese (国密) |
| SM4 | Symmetric | ✅ | ❌ | Chinese (国密) |

## Options, Constants, and Key Helpers

| Area | Public API |
| --- | --- |
| AES modes/options | `AESMode`, `AesModeGcm`, `AesModeCbc`, `WithAESMode(mode)`, `WithAESIv(iv)` |
| RSA modes/options | `RSAMode`, `RSASignMode`, `RsaModeOAEP`, `RsaModePKCS1v15`, `RsaSignModePSS`, `RsaSignModePKCS1v15`, `WithRSAMode(mode)`, `WithRSASignMode(mode)` |
| SM4 options | `WithSM4Iv(iv)` |
| ECDSA curves | `ECDSACurve`, `EcdsaCurveP224`, `EcdsaCurveP256`, `EcdsaCurveP384`, `EcdsaCurveP521` |
| ECIES curves | `ECIESCurve`, `EciesCurveP256`, `EciesCurveP384`, `EciesCurveP521`, `EciesCurveX25519` |
| key helpers | `GenerateECDSAKey(curve)`, `GenerateECIESKey(curve)` |
| option types | `AESOption`, `RSAOption`, `SM4Option` |

## Error Sentinels

| Group | Errors |
| --- | --- |
| key availability | `ErrAtLeastOneKeyRequired`, `ErrPublicKeyRequiredForEncrypt`, `ErrPrivateKeyRequiredForDecrypt`, `ErrPrivateKeyRequiredForSign`, `ErrPublicKeyRequiredForVerify` |
| key parsing/type | `ErrFailedDecodePEMBlock`, `ErrUnsupportedPEMType`, `ErrNotRSAPrivateKey`, `ErrNotRSAPublicKey`, `ErrNotECDSAPrivateKey`, `ErrNotECDSAPublicKey` |
| symmetric crypto | `ErrInvalidAESKeySize`, `ErrInvalidSM4KeySize`, `ErrInvalidIVSizeCBC`, `ErrCiphertextNotMultipleOfBlock`, `ErrCiphertextTooShort`, `ErrInvalidPadding` |
| input/signature | `ErrDataEmpty`, `ErrInvalidSignature` |
