---
sidebar_position: 5
---

# Hashx

The `hashx` package provides one-line hash and HMAC functions for common algorithms.

## Reviewed Public Surface

The current source audit for `github.com/coldsmirk/vef-framework-go/hashx`
covers 15 top-level exported symbols, no exported fields, and no exported
methods. The reviewed public-surface fingerprint is
`22e7f661d37170d375f54592fa00078a3ea92b1b93459672709422aab54d5a01`.

Reviewed APIs:

| API | Contract |
| --- | --- |
| `hashx.MD5(data string)` | Converts `data` to bytes, calls `MD5Bytes`, and returns a 32-character lowercase hex MD5 digest |
| `hashx.MD5Bytes(data []byte)` | Hashes raw bytes with `crypto/md5` and returns a 32-character lowercase hex digest |
| `hashx.SHA1(data string)` | Converts `data` to bytes, calls `SHA1Bytes`, and returns a 40-character lowercase hex SHA-1 digest |
| `hashx.SHA1Bytes(data []byte)` | Hashes raw bytes with `crypto/sha1` and returns a 40-character lowercase hex digest |
| `hashx.SHA256(data string)` | Converts `data` to bytes, calls `SHA256Bytes`, and returns a 64-character lowercase hex SHA-256 digest |
| `hashx.SHA256Bytes(data []byte)` | Hashes raw bytes with `crypto/sha256` and returns a 64-character lowercase hex digest |
| `hashx.SHA512(data string)` | Converts `data` to bytes, calls `SHA512Bytes`, and returns a 128-character lowercase hex SHA-512 digest |
| `hashx.SHA512Bytes(data []byte)` | Hashes raw bytes with `crypto/sha512` and returns a 128-character lowercase hex digest |
| `hashx.SM3(data string)` | Converts `data` to bytes, calls `SM3Bytes`, and returns a 64-character lowercase hex SM3 digest |
| `hashx.SM3Bytes(data []byte)` | Hashes raw bytes with `github.com/tjfoc/gmsm/sm3` and returns a 64-character lowercase hex digest |
| `hashx.HmacMD5(key, data []byte)` | Computes HMAC-MD5 with `key` and `data`, returning a 32-character lowercase hex digest |
| `hashx.HmacSHA1(key, data []byte)` | Computes HMAC-SHA1 with `key` and `data`, returning a 40-character lowercase hex digest |
| `hashx.HmacSHA256(key, data []byte)` | Computes HMAC-SHA256 with `key` and `data`, returning a 64-character lowercase hex digest |
| `hashx.HmacSHA512(key, data []byte)` | Computes HMAC-SHA512 with `key` and `data`, returning a 128-character lowercase hex digest |
| `hashx.HmacSM3(key, data []byte)` | Computes HMAC-SM3 with `key` and `data`, returning a 64-character lowercase hex digest |

## Hash Functions

String hash functions accept a string, convert it with `[]byte(data)`, and
return a lowercase hex-encoded digest:

```go
import "github.com/coldsmirk/vef-framework-go/hashx"

hashx.MD5("hello")     // "5d41402abc4b2a76b9719d911017c592"
hashx.SHA1("hello")    // "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"
hashx.SHA256("hello")  // "2cf24dba5fb0a30e26e83b2ac5b9e29e..."
hashx.SHA512("hello")  // "9b71d224bd62f3785d96d46ad3ea3d73..."
hashx.SM3("hello")     // SM3 (Chinese National Standard)
```

### Byte Variants

For raw byte input, use the `Bytes` variants. They return the digest of the
byte slice exactly as passed; a nil slice is accepted and hashes like an empty
slice.

```go
hashx.MD5Bytes(data)
hashx.SHA1Bytes(data)
hashx.SHA256Bytes(data)
hashx.SHA512Bytes(data)
hashx.SM3Bytes(data)
```

## HMAC Functions

All HMAC functions accept `key` and `data` as byte slices and return lowercase
hex-encoded MACs:

```go
key := []byte("secret")
data := []byte("message")

hashx.HmacMD5(key, data)
hashx.HmacSHA1(key, data)
hashx.HmacSHA256(key, data)
hashx.HmacSHA512(key, data)
hashx.HmacSM3(key, data)    // HMAC with SM3 (Chinese National Standard)
```

## Algorithm Reference

| Function | Algorithm | Output Length |
| --- | --- | --- |
| `MD5`, `MD5Bytes`, `HmacMD5` | MD5 / HMAC-MD5 | 32 hex chars |
| `SHA1`, `SHA1Bytes`, `HmacSHA1` | SHA-1 / HMAC-SHA1 | 40 hex chars |
| `SHA256`, `SHA256Bytes`, `HmacSHA256` | SHA-256 / HMAC-SHA256 | 64 hex chars |
| `SHA512`, `SHA512Bytes`, `HmacSHA512` | SHA-512 / HMAC-SHA512 | 128 hex chars |
| `SM3`, `SM3Bytes`, `HmacSM3` | SM3 / HMAC-SM3 (Chinese National Standard) | 64 hex chars |

:::caution
`MD5` and `SHA1` are kept for compatibility checksums and legacy integrations.
Do not use them for password storage; use the `password` package instead.
:::
