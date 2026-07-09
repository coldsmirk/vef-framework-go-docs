---
sidebar_position: 6
---

# Hashx

`hashx` 包提供一行调用的哈希和 HMAC 函数，覆盖常见算法。

## API 参考

| API | Contract |
| --- | --- |
| `hashx.MD5(data string)` | 将 `data` 转成 bytes，调用 `MD5Bytes`，返回 32 个字符的 lowercase hex MD5 digest |
| `hashx.MD5Bytes(data []byte)` | 使用 `crypto/md5` 对原始 bytes 求 hash，返回 32 个字符的 lowercase hex digest |
| `hashx.SHA1(data string)` | 将 `data` 转成 bytes，调用 `SHA1Bytes`，返回 40 个字符的 lowercase hex SHA-1 digest |
| `hashx.SHA1Bytes(data []byte)` | 使用 `crypto/sha1` 对原始 bytes 求 hash，返回 40 个字符的 lowercase hex digest |
| `hashx.SHA256(data string)` | 将 `data` 转成 bytes，调用 `SHA256Bytes`，返回 64 个字符的 lowercase hex SHA-256 digest |
| `hashx.SHA256Bytes(data []byte)` | 使用 `crypto/sha256` 对原始 bytes 求 hash，返回 64 个字符的 lowercase hex digest |
| `hashx.SHA512(data string)` | 将 `data` 转成 bytes，调用 `SHA512Bytes`，返回 128 个字符的 lowercase hex SHA-512 digest |
| `hashx.SHA512Bytes(data []byte)` | 使用 `crypto/sha512` 对原始 bytes 求 hash，返回 128 个字符的 lowercase hex digest |
| `hashx.SM3(data string)` | 将 `data` 转成 bytes，调用 `SM3Bytes`，返回 64 个字符的 lowercase hex SM3 digest |
| `hashx.SM3Bytes(data []byte)` | 使用 `github.com/tjfoc/gmsm/sm3` 对原始 bytes 求 hash，返回 64 个字符的 lowercase hex digest |
| `hashx.HmacMD5(key, data []byte)` | 使用 `key` 和 `data` 计算 HMAC-MD5，返回 32 个字符的 lowercase hex digest |
| `hashx.HmacSHA1(key, data []byte)` | 使用 `key` 和 `data` 计算 HMAC-SHA1，返回 40 个字符的 lowercase hex digest |
| `hashx.HmacSHA256(key, data []byte)` | 使用 `key` 和 `data` 计算 HMAC-SHA256，返回 64 个字符的 lowercase hex digest |
| `hashx.HmacSHA512(key, data []byte)` | 使用 `key` 和 `data` 计算 HMAC-SHA512，返回 128 个字符的 lowercase hex digest |
| `hashx.HmacSM3(key, data []byte)` | 使用 `key` 和 `data` 计算 HMAC-SM3，返回 64 个字符的 lowercase hex digest |

## 哈希函数

字符串哈希函数接受 string，用 `[]byte(data)` 转换后返回 lowercase
hex-encoded digest：

```go
import "github.com/coldsmirk/vef-framework-go/hashx"

hashx.MD5("hello")     // "5d41402abc4b2a76b9719d911017c592"
hashx.SHA1("hello")    // "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"
hashx.SHA256("hello")  // "2cf24dba5fb0a30e26e83b2ac5b9e29e..."
hashx.SHA512("hello")  // "9b71d224bd62f3785d96d46ad3ea3d73..."
hashx.SM3("hello")     // SM3（国密标准）
```

### 字节变体

用于原始字节输入时，使用 `Bytes` 变体。它们会对传入的 byte slice
原样求 digest；nil slice 也会被接受，并按 empty slice 一样求 hash。

```go
hashx.MD5Bytes(data)
hashx.SHA1Bytes(data)
hashx.SHA256Bytes(data)
hashx.SHA512Bytes(data)
hashx.SM3Bytes(data)
```

## HMAC 函数

所有 HMAC 函数接受 `key` 和 `data` 两个 byte slices，并返回 lowercase
hex-encoded MAC：

```go
key := []byte("secret")
data := []byte("message")

hashx.HmacMD5(key, data)
hashx.HmacSHA1(key, data)
hashx.HmacSHA256(key, data)
hashx.HmacSHA512(key, data)
hashx.HmacSM3(key, data)    // 使用 SM3 的 HMAC（国密标准）
```

## 算法参考

| 函数 | 算法 | 输出长度 |
| --- | --- | --- |
| `MD5`, `MD5Bytes`, `HmacMD5` | MD5 / HMAC-MD5 | 32 个十六进制字符 |
| `SHA1`, `SHA1Bytes`, `HmacSHA1` | SHA-1 / HMAC-SHA1 | 40 个十六进制字符 |
| `SHA256`, `SHA256Bytes`, `HmacSHA256` | SHA-256 / HMAC-SHA256 | 64 个十六进制字符 |
| `SHA512`, `SHA512Bytes`, `HmacSHA512` | SHA-512 / HMAC-SHA512 | 128 个十六进制字符 |
| `SM3`, `SM3Bytes`, `HmacSM3` | SM3 / HMAC-SM3（国密标准）| 64 个十六进制字符 |

:::caution
`MD5` 和 `SHA1` 只适合 compatibility checksum 和 legacy integration。
不要把它们用于密码存储；密码场景请使用 `password` 包。
:::
