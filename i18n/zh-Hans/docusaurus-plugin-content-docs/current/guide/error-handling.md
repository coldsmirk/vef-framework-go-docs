---
sidebar_position: 9
---

# 错误处理

VEF 把传输层 HTTP 行为和业务层结果码区分开来，但最终都会通过统一的 `code / message / data` 响应包络返回给客户端。

## Result 模型总览

VEF 主要使用两种相关但职责不同的结果类型：

| 类型 | 作用 |
| --- | --- |
| `result.Result` | 最终返回给客户端的响应载荷 |
| `result.Error` | 应用代码内部使用的结构化错误对象 |

`result.Result` 的形态如下：

```json
{
  "code": 0,
  "message": "Success",
  "data": {}
}
```

## 成功响应

成功的 handler 最常见的写法是：

```go
return result.Ok(data).Response(ctx)
```

`result.Ok(...)` 支持的形式：

| 写法 | 含义 |
| --- | --- |
| `result.Ok()` | 不带数据的成功结果 |
| `result.Ok(data)` | 带数据的成功结果 |
| `result.Ok(result.WithMessage(...))` | 自定义成功消息 |
| `result.Ok(data, result.WithMessage(...))` | 同时自定义消息和数据 |

## 结构化错误创建

业务失败时，最常见的写法是：

```go
return result.Err(
  "user already exists",
  result.WithCode(result.ErrCodeRecordAlreadyExists),
)
```

`result.Err(...)` 支持的形式：

| 写法 | 含义 |
| --- | --- |
| `result.Err()` | 默认业务错误 |
| `result.Err("message")` | 自定义错误消息 |
| `result.Err("message", result.WithCode(...))` | 自定义业务码 |
| `result.Err("message", result.WithStatus(...))` | 自定义 HTTP 状态 |
| `result.Err("message", result.WithCode(...), result.WithStatus(...))` | 同时覆盖业务码和 HTTP 状态 |

`result.Errf(...)` 则是格式化版本。

## Result 配置项

可用的 result option：

| Option | 作用对象 | 作用 |
| --- | --- | --- |
| `result.WithCode(code)` | `result.Err(...)` | 设置业务错误码 |
| `result.WithStatus(status)` | `result.Err(...)` | 设置 HTTP 状态码 |
| `result.WithMessage(message)` | `result.Ok(...)` | 覆盖成功消息 |
| `result.WithMessagef(format, ...)` | `result.Ok(...)` | 生成格式化成功消息 |

## 内置错误族

VEF 在 `result` 包里预置了大量常用错误。

### 认证错误

| 错误值 | 业务码 | 默认 HTTP 状态 |
| --- | --- | --- |
| `result.ErrUnauthenticated` | `ErrCodeUnauthenticated` | `401` |
| `result.ErrTokenExpired` | `ErrCodeTokenExpired` | `401` |
| `result.ErrTokenInvalid` | `ErrCodeTokenInvalid` | `401` |
| `result.ErrTokenNotValidYet` | `ErrCodeTokenNotValidYet` | `401` |
| `result.ErrTokenInvalidIssuer` | `ErrCodeTokenInvalidIssuer` | `401` |
| `result.ErrTokenInvalidAudience` | `ErrCodeTokenInvalidAudience` | `401` |
| `result.ErrTokenMissingSubject` | `ErrCodeTokenMissingSubject` | `401` |
| `result.ErrTokenMissingTokenType` | `ErrCodeTokenMissingTokenType` | `401` |

### Signature / 外部应用认证错误

| 错误值 | 业务码 | 默认 HTTP 状态 |
| --- | --- | --- |
| `result.ErrAppIDRequired` | `ErrCodeAppIDRequired` | `401` |
| `result.ErrTimestampRequired` | `ErrCodeTimestampRequired` | `401` |
| `result.ErrSignatureRequired` | `ErrCodeSignatureRequired` | `401` |
| `result.ErrTimestampInvalid` | `ErrCodeTimestampInvalid` | `401` |
| `result.ErrSignatureExpired` | `ErrCodeSignatureExpired` | `401` |
| `result.ErrSignatureInvalid` | `ErrCodeSignatureInvalid` | `401` |
| `result.ErrExternalAppNotFound` | `ErrCodeExternalAppNotFound` | `401` |
| `result.ErrExternalAppDisabled` | `ErrCodeExternalAppDisabled` | `401` |
| `result.ErrIPNotAllowed` | `ErrCodeIPNotAllowed` | `401` |
| `result.ErrNonceRequired` | `ErrCodeNonceRequired` | `401` |
| `result.ErrNonceInvalid` | `ErrCodeNonceInvalid` | `401` |
| `result.ErrNonceAlreadyUsed` | `ErrCodeNonceAlreadyUsed` | `401` |
| `result.ErrAuthHeaderMissing` | `ErrCodeAuthHeaderMissing` | `401` |
| `result.ErrAuthHeaderInvalid` | `ErrCodeAuthHeaderInvalid` | `401` |

### Challenge 流错误

| 错误值 | 业务码 | 默认 HTTP 状态 |
| --- | --- | --- |
| `result.ErrChallengeTokenInvalid` | `ErrCodeChallengeTokenInvalid` | `401` |
| `result.ErrChallengeTypeInvalid` | `ErrCodeChallengeTypeInvalid` | `400` |
| `result.ErrOTPCodeRequired` | `ErrCodeOTPCodeRequired` | `400` |
| `result.ErrOTPCodeInvalid` | `ErrCodeOTPCodeInvalid` | `401` |
| `result.ErrNewPasswordRequired` | `ErrCodeNewPasswordRequired` | `400` |
| `result.ErrDepartmentRequired` | `ErrCodeDepartmentRequired` | `400` |

### 授权与请求错误

| 错误值 | 业务码 | 默认 HTTP 状态 |
| --- | --- | --- |
| `result.ErrAccessDenied` | `ErrCodeAccessDenied` | `403` |
| `result.ErrTooManyRequests` | `ErrCodeTooManyRequests` | `429` |
| `result.ErrRequestTimeout` | `ErrCodeRequestTimeout` | `408` |
| `result.ErrUnknown` | `ErrCodeUnknown` | `500` |

### 业务错误

| 错误值 | 业务码 | 默认 HTTP 状态 |
| --- | --- | --- |
| `result.ErrRecordNotFound` | `ErrCodeRecordNotFound` | `200` |
| `result.ErrRecordAlreadyExists` | `ErrCodeRecordAlreadyExists` | `200` |
| `result.ErrForeignKeyViolation` | `ErrCodeForeignKeyViolation` | `200` |
| `result.ErrDangerousSQL` | `ErrCodeDangerousSQL` | `200` |

### 错误构造器

以下 helper 会生成具有特定语义的结构化错误：

| 构造器 | 典型输出 |
| --- | --- |
| `result.ErrNotImplemented(message)` | `501 Not Implemented` |
| `result.ErrCredentialsInvalid(message)` | `401 Unauthorized`，并带 credentials-invalid 业务码 |
| `result.ErrPrincipalInvalid(message)` | `401 Unauthorized`，并带 principal-invalid 业务码 |

## 业务码范围

常见结果码范围如下：

| 范围 | 含义 |
| --- | --- |
| `0` | 成功 |
| `1000-1099` | 认证与 challenge 错误 |
| `1100-1199` | 授权错误 |
| `1200-1499` | 资源、媒体类型和请求错误 |
| `1500-1699` | 未实现与 SQL 相关错误 |
| `1900-1999` | 未知错误 |
| `2000+` | 业务错误 |

## Fiber 错误映射

应用层会把部分 `fiber.Error` 映射为结构化 result 响应。

当前内置映射如下：

| Fiber HTTP 状态 | Result 业务码 | Message key |
| --- | --- | --- |
| `401` | `ErrCodeUnauthenticated` | `ErrMessageUnauthenticated` |
| `403` | `ErrCodeAccessDenied` | `ErrMessageAccessDenied` |
| `404` | `ErrCodeNotFound` | `ErrMessageNotFound` |
| `415` | `ErrCodeUnsupportedMediaType` | `ErrMessageUnsupportedMediaType` |
| `408` | `ErrCodeRequestTimeout` | `ErrMessageRequestTimeout` |

如果某个 `fiber.Error` 状态码没有映射，VEF 会先记录日志，再回退为通用 unknown error。

## 错误解析顺序

运行时，VEF 会按以下顺序处理错误：

1. `fiber.Error`
2. `result.Error`
3. 其他未识别错误 -> `result.ErrUnknown`

这也是为什么对预期业务失败，应该优先返回显式的 `result.Error`，而不是不透明的普通错误。

## 常见模式

### 成功 + 数据

```go
return result.Ok(user).Response(ctx)
```

### 成功 + 自定义消息

```go
return result.Ok(
  user,
  result.WithMessage("user synced"),
).Response(ctx)
```

### 带业务码的业务错误

```go
return result.Err(
  "user already exists",
  result.WithCode(result.ErrCodeRecordAlreadyExists),
)
```

### 显式覆盖 HTTP 状态

```go
return result.Err(
  "forbidden",
  result.WithCode(result.ErrCodeAccessDenied),
  result.WithStatus(fiber.StatusForbidden),
)
```

## 实践建议

- 把 `result` 当作对外响应契约来看待
- 只要内置错误已经能表达你的语义，就优先复用
- 当客户端需要按不同失败类型做不同处理时，再定义明确的业务码
- 对预期业务失败，优先返回结构化 `result.Error`，不要返回随意拼接的字符串错误
- 除非你有意绕开框架契约，否则不要手工输出原始 JSON 响应

## 下一步

继续阅读 [认证](../security/authentication)，看认证失败是如何进入这一套结果模型的。
