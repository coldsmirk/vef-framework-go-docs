---
sidebar_position: 5
---

# Hashx

The `hashx` package provides one-line hash and HMAC functions for common algorithms.

## Hash Functions

All hash functions accept a string and return a hex-encoded hash:

```go
import "github.com/coldsmirk/vef-framework-go/hashx"

hashx.MD5("hello")     // "5d41402abc4b2a76b9719d911017c592"
hashx.SHA1("hello")    // "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"
hashx.SHA256("hello")  // "2cf24dba5fb0a30e26e83b2ac5b9e29e..."
hashx.SHA512("hello")  // "9b71d224bd62f3785d96d46ad3ea3d73..."
hashx.SM3("hello")     // SM3 (Chinese National Standard)
```

### Byte Variants

For raw byte input:

```go
hashx.MD5Bytes(data)
hashx.SHA1Bytes(data)
hashx.SHA256Bytes(data)
hashx.SHA512Bytes(data)
hashx.SM3Bytes(data)
```

## HMAC Functions

All HMAC functions accept a key and data as byte slices:

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
| `MD5` | MD5 | 32 hex chars |
| `SHA1` | SHA-1 | 40 hex chars |
| `SHA256` | SHA-256 | 64 hex chars |
| `SHA512` | SHA-512 | 128 hex chars |
| `SM3` | SM3 (国密) | 64 hex chars |
