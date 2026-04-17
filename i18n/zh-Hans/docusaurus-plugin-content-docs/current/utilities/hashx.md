---
sidebar_position: 5
---

# Hashx

`hashx` 包提供一行调用的哈希和 HMAC 函数，覆盖常见算法。

## 哈希函数

所有哈希函数接受字符串并返回十六进制编码的哈希值：

```go
import "github.com/coldsmirk/vef-framework-go/hashx"

hashx.MD5("hello")     // "5d41402abc4b2a76b9719d911017c592"
hashx.SHA1("hello")    // "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"
hashx.SHA256("hello")  // "2cf24dba5fb0a30e26e83b2ac5b9e29e..."
hashx.SHA512("hello")  // "9b71d224bd62f3785d96d46ad3ea3d73..."
hashx.SM3("hello")     // SM3（国密标准）
```

### 字节变体

用于原始字节输入：

```go
hashx.MD5Bytes(data)
hashx.SHA1Bytes(data)
hashx.SHA256Bytes(data)
hashx.SHA512Bytes(data)
hashx.SM3Bytes(data)
```

## HMAC 函数

所有 HMAC 函数接受密钥和数据的字节切片：

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
| `MD5` | MD5 | 32 个十六进制字符 |
| `SHA1` | SHA-1 | 40 个十六进制字符 |
| `SHA256` | SHA-256 | 64 个十六进制字符 |
| `SHA512` | SHA-512 | 128 个十六进制字符 |
| `SM3` | SM3（国密）| 64 个十六进制字符 |
