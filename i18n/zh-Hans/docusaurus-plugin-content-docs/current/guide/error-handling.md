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

`result.Error` 不会被直接序列化为客户端响应。应用错误处理器会提取
`Code`、`Message` 和 `Status`，把 `Code` 与 `Message` 放入
`result.Result` 信封，并使用 `Status` 作为 HTTP status。

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

`result.Ok(...)` 最多接受一个 data 参数。data 必须位于 `OkOption`
之前；多个 data 参数，或 data 出现在 option 之后，都会 panic。

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

可选 message string 必须是 `Err(...)` 的第一个参数，后续参数必须是
`ErrOption`。`result.Errf(...)` 则是格式化版本；它至少需要一个 format
arg，并且 options 必须放在所有 format args 之后。

## Result 配置项

可用的 result option：

| Option | 作用对象 | 作用 |
| --- | --- | --- |
| `result.WithCode(code)` | `result.Err(...)` | 设置业务错误码 |
| `result.WithStatus(status)` | `result.Err(...)` | 设置 HTTP 状态码 |
| `result.WithMessage(message)` | `result.Ok(...)` | 覆盖成功消息 |
| `result.WithMessagef(format, ...)` | `result.Ok(...)` | 生成格式化成功消息 |

`result.Error` 没有 message option；要用第一个 `Err(...)` 参数或
`Errf(...)` 设置错误消息。默认 `Err(...)` 使用
`result.ErrCodeDefault`（`2000`）、message `i18n.T(result.ErrMessage)`、
HTTP status `200 OK`。

## 错误身份 Helper

`result.Error` 通过只比较 `Code` 实现 `errors.Is`。`Message` 和 `Status`
不影响 identity，因此动态格式化的错误只要 code 相同，也能匹配预置 sentinel。
使用 `result.AsErr(err)` 从 error chain 中提取 `result.Error`，使用
`result.IsRecordNotFound(err)` 判断记录不存在 sentinel。

## 内置错误族

VEF 在框架各处预置了一批 `result.Error`。v0.25 起，模块专属错误从 `result` 包迁出到各自模块包中，`result` 包只保留跨模块通用错误。错误码本身不变，只是 import 路径改了。

### 通用错误（`result` 包）

| 错误值 | 业务码 | 默认 HTTP 状态 |
| --- | --- | --- |
| `result.ErrAccessDenied` | `result.ErrCodeAccessDenied`（1100） | `403` |
| `result.ErrTooManyRequests` | `result.ErrCodeTooManyRequests`（1401） | `429` |
| `result.ErrRequestTimeout` | `result.ErrCodeRequestTimeout`（1402） | `408` |
| `result.ErrUnknown` | `result.ErrCodeUnknown`（1900） | `500` |
| `result.ErrRecordNotFound` | `result.ErrCodeRecordNotFound`（2001） | `200` |
| `result.ErrRecordAlreadyExists` | `result.ErrCodeRecordAlreadyExists`（2002） | `200` |
| `result.ErrForeignKeyViolation` | `result.ErrCodeForeignKeyViolation`（2003） | `200` |
| `result.ErrDangerousSQL` | `result.ErrCodeDangerousSQL`（1600） | `200` |
| `result.ErrNotImplemented(message)` | `result.ErrCodeNotImplemented`（1500） | `501` |

`result.ErrCodeBadRequest`、`result.ErrCodeNotFound` 和
`result.ErrCodeUnsupportedMediaType` 是 exported building-block constants，
不是 `result` 包里的预置 `result.Error` 值。

### 安全错误（`security` 包）

认证、签名、challenge 流相关的错误现在位于 `github.com/coldsmirk/vef-framework-go/security`，对应的 `ErrCodeXxx` 常量也跟着搬过去了。认证错误使用 `1000-1022`；challenge 错误预留 `1030-1039`，当前实际导出 `1031`、`1033` 和 `1034-1038`。

| 错误值 | 业务码 | 默认 HTTP 状态 |
| --- | --- | --- |
| `security.ErrUnauthenticated` | `security.ErrCodeUnauthenticated`（1000） | `401` |
| `security.ErrTokenExpired` | `security.ErrCodeTokenExpired`（1002） | `401` |
| `security.ErrTokenInvalid` | `security.ErrCodeTokenInvalid`（1003） | `401` |
| `security.ErrTokenNotValidYet` | `security.ErrCodeTokenNotValidYet`（1004） | `401` |
| `security.ErrTokenInvalidIssuer` | `security.ErrCodeTokenInvalidIssuer`（1005） | `401` |
| `security.ErrTokenInvalidAudience` | `security.ErrCodeTokenInvalidAudience`（1006） | `401` |
| `security.ErrAppIDRequired` | `security.ErrCodeAppIDRequired`（1009） | `401` |
| `security.ErrTimestampRequired` | `security.ErrCodeTimestampRequired`（1010） | `401` |
| `security.ErrSignatureRequired` | `security.ErrCodeSignatureRequired`（1011） | `401` |
| `security.ErrTimestampInvalid` | `security.ErrCodeTimestampInvalid`（1012） | `401` |
| `security.ErrSignatureExpired` | `security.ErrCodeSignatureExpired`（1013） | `401` |
| `security.ErrSignatureInvalid` | `security.ErrCodeSignatureInvalid`（1017） | `401` |
| `security.ErrExternalAppNotFound` | `security.ErrCodeExternalAppNotFound`（1014） | `401` |
| `security.ErrExternalAppDisabled` | `security.ErrCodeExternalAppDisabled`（1015） | `401` |
| `security.ErrIPNotAllowed` | `security.ErrCodeIPNotAllowed`（1016） | `401` |
| `security.ErrNonceRequired` | `security.ErrCodeNonceRequired`（1018） | `401` |
| `security.ErrNonceInvalid` | `security.ErrCodeNonceInvalid`（1019） | `401` |
| `security.ErrNonceAlreadyUsed` | `security.ErrCodeNonceAlreadyUsed`（1020） | `401` |
| `security.ErrAuthHeaderMissing` | `security.ErrCodeAuthHeaderMissing`（1021） | `401` |
| `security.ErrAuthHeaderInvalid` | `security.ErrCodeAuthHeaderInvalid`（1022） | `401` |
| `security.ErrChallengeTokenInvalid` | `security.ErrCodeChallengeTokenInvalid`（1031） | `401` |
| `security.ErrChallengeTypeInvalid` | `security.ErrCodeChallengeTypeInvalid`（1033） | `400` |
| `security.ErrOTPCodeRequired` | `security.ErrCodeOTPCodeRequired`（1035） | `400` |
| `security.ErrOTPCodeInvalid` | `security.ErrCodeOTPCodeInvalid`（1036） | `401` |
| `security.ErrNewPasswordRequired` | `security.ErrCodeNewPasswordRequired`（1037） | `400` |
| `security.ErrDepartmentRequired` | `security.ErrCodeDepartmentRequired`（1038） | `400` |
| `security.ErrCredentialsInvalid(message)` | `security.ErrCodeCredentialsInvalid`（1008） | `401` |
| `security.ErrPrincipalInvalid(message)` | `security.ErrCodePrincipalInvalid`（1007） | `401` |

> v0.25.1 删除了未使用的 `ErrTokenMissingSubject` / `ErrTokenMissingTokenType` sentinel，并对周围的错误码做了压缩。框架没有为旧版本保留兼容层，请直接更新调用点。

### 其他模块错误

| 模块包 | 错误值 | 编号区间 |
| --- | --- | --- |
| `api` | `api.ErrInvalidRequestParams`、`api.ErrInvalidRequestMeta` | 1400（`result.ErrCodeBadRequest`） |
| `monitor` | `monitor.ErrNotReady`、`monitor.ErrCollectionFailed` | 2100-2101 |
| `storage` | `storage.ErrInvalidFileKey`、`storage.ErrFileNotFound`、`storage.ErrFailedToGetFile`，以及 `storage.ErrUploadRequiresMultipart`、`storage.ErrUploadPartsIncomplete`、`storage.ErrAbortFailed` 等 multipart upload / claim 错误 | 2200-2219 |
| `schema` | `schema.ErrTableNotFound` | 2300 |
| `crud` | `crud.ErrCodeProcessorInvalidReturn`、CRUD import/export 和主键相关 result 错误，以及 `crud.ErrModelNoPrimaryKey`、`crud.ErrAuditUserCompositePK` 等普通 sentinel | 2400-2410 |
| `expression` | `expression.ErrEvaluationFailed` | 2500 |
| `approval` | 公开普通 sentinel：`approval.ErrCrossTenantAccess`、`approval.ErrInvalidBusinessIdentifier`、`approval.ErrUnknownNodeKind`、`approval.ErrNodeDataUnmarshal`；内置审批资源返回 internal `result.Error` | 40001-40702 |

> 这四个公开 `approval` sentinel 都是普通 Go 错误，**不**是 `result.Error`，没有 code/status 字段。内置审批资源响应使用 internal 的 40xxx result envelope 目录；完整 code 与 message key 见 [Approval 模块](../modules/approval)。

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
| `401` | `security.ErrCodeUnauthenticated` | `security.ErrMessageUnauthenticated` |
| `403` | `result.ErrCodeAccessDenied` | `result.ErrMessageAccessDenied` |
| `404` | `result.ErrCodeNotFound` | `result.ErrMessageNotFound` |
| `415` | `result.ErrCodeUnsupportedMediaType` | `result.ErrMessageUnsupportedMediaType` |
| `408` | `result.ErrCodeRequestTimeout` | `result.ErrMessageRequestTimeout` |

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
