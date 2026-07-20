---
sidebar_position: 6
---

# Results and Errors

VEF 把传输层 HTTP 行为和业务层结果码区分开来，但最终都会通过统一的 `code / message / data` 响应信封返回给客户端。`result` 包（`github.com/coldsmirk/vef-framework-go/result`）定义了这个信封，以及应用代码返回、用于产出该信封的结构化业务错误类型。

## Success & Error Envelope

VEF 主要使用两种紧密相关但职责不同的结果类型：

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

```go
type Result struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    any    `json:"data"`
}
```

| 字段/方法 | 说明 |
| --- | --- |
| `Result.Code` | 业务结果码，序列化为 `code`；`0` 表示成功。 |
| `Result.Message` | 面向用户或经 i18n 解析后的消息，序列化为 `message`。 |
| `Result.Data` | 可选响应载荷，序列化为 `data`；`nil` data 会保留为 JSON `null`。 |
| `Result.Response(ctx, status...)` | 以 JSON 发送结果；HTTP status 默认 `200 OK`，如果传入 status 则使用第一个值。 |
| `Result.IsOk()` | 仅当 `Code == result.OkCode` 时返回 true。 |

`result.Error` 刻意没有 JSON tags —— 它不是公开的 JSON 信封。应用代码把它作为 `error` 返回；应用错误处理器会提取 `Code`、`Message` 和 `Status`，把 `Code` 与 `Message` 放入 `result.Result` 信封，并使用 `Status` 作为 HTTP status。`result.Error` 不会被直接序列化为客户端响应。如果目标是公开的 `code / message / data` 响应信封，不要直接序列化 `result.Error`。

```go
type Error struct {
    Code    int
    Message string
    Status  int
}
```

| 字段/方法 | 说明 |
| --- | --- |
| `Error.Code` | 业务错误码，会进入响应信封，也用于 `errors.Is` 比较。 |
| `Error.Message` | `Error.Error()` 返回的错误消息，也会复制进响应信封。 |
| `Error.Status` | 应用错误处理器把错误转换成 `Result` 时使用的 HTTP status。 |
| `Error.Error()` | 通过返回 `Message` 实现 `error` 接口。 |
| `Error.Is(target)` | 只按 `Code` 匹配另一个 `result.Error`。 |

## Building Results (Ok/Err/Options)

### 创建成功响应

```go
import "github.com/coldsmirk/vef-framework-go/result"

// 简单成功（无数据）
return result.Ok().Response(ctx)

// 带数据的成功
return result.Ok(user).Response(ctx)

// 自定义消息的成功
return result.Ok(result.WithMessage("Created successfully")).Response(ctx)

// 带数据和自定义消息
return result.Ok(user, result.WithMessage("User created")).Response(ctx)

// 自定义 HTTP 状态码（例如 201 Created）
return result.Ok(user).Response(ctx, 201)
```

`result.Ok(dataOrOptions...)` 支持的形式：

| 写法 | 含义 |
| --- | --- |
| `result.Ok()` | 不带数据的成功结果 |
| `result.Ok(data)` | 带数据的成功结果 |
| `result.Ok(result.WithMessage(...))` | 自定义成功消息 |
| `result.Ok(data, result.WithMessage(...))` | 同时自定义消息和数据 |

`result.Ok(...)` 最多接受一个 data 参数。如果提供 data，它必须位于 `OkOption` 之前。传入多个 data 参数，或把 data 放在 option 之后，都会 panic。

`result.OkOption` 是 option 函数类型（`func(*Result)`）：

| 选项 | 效果 |
| --- | --- |
| `result.WithMessage(message)` | 将 `Result.Message` 精确设置为 `message`，空字符串也会生效 |
| `result.WithMessagef(format, args...)` | 使用 `fmt.Sprintf(format, args...)` 设置 `Result.Message` |

默认成功码是 `result.OkCode`（`0`），默认成功消息是 `i18n.T(result.OkMessage)`。`result.OkCode` / `result.OkMessage` 共同定义了默认成功信封。

### 创建错误响应

`result.Err(...)` 会构造 `result.Error`。直接返回它，框架的错误处理器会把它渲染为统一响应信封：

```go
// 默认业务错误（code 2000，message 来自 i18n catalog）
return result.Err()

// 自定义错误消息
return result.Err("something went wrong")

// 带自定义错误码
return result.Err("not found", result.WithCode(result.ErrCodeRecordNotFound))

// 自定义 HTTP 状态码，同时保留结构化响应信封
return result.Err("forbidden",
    result.WithCode(result.ErrCodeAccessDenied),
    result.WithStatus(fiber.StatusForbidden),
)
```

`result.Err(messageOrOptions...)` 支持的形式：

| 写法 | 含义 |
| --- | --- |
| `result.Err()` | 默认业务错误 |
| `result.Err("message")` | 自定义错误消息 |
| `result.Err("message", result.WithCode(...))` | 自定义业务码 |
| `result.Err("message", result.WithStatus(...))` | 自定义 HTTP 状态 |
| `result.Err("message", result.WithCode(...), result.WithStatus(...))` | 同时覆盖业务码和 HTTP 状态 |

可选 message string 必须是 `Err(...)` 的第一个参数。后续参数必须是 `ErrOption`，其他类型会 panic。错误消息没有 option 形式，`result.Error` 没有 message option：要用第一个 `Err(...)` 参数，或使用 `Errf(...)` 设置错误消息。

`result.Errf(format, args...)` 是格式化版本：

```go
return result.Errf("user %s not found", username,
    result.WithCode(result.ErrCodeRecordNotFound),
)
```

`Errf` 至少需要一个 format arg。`ErrOption` 必须放在所有 format args 之后；把 option 放在 format args 之前或中间都会 panic。

`result.ErrOption` 是 option 函数类型（`func(*Error)`）：

| 选项 | 效果 |
| --- | --- |
| `result.WithCode(code)` | 设置 `Error.Code`；默认是 `result.ErrCodeDefault`（`2000`） |
| `result.WithStatus(status)` | 设置 `Error.Status`；默认是 `200 OK` |

默认 `Err(...)` 和 `Errf(...)` 使用 `result.ErrCodeDefault`（`2000`）、message `i18n.T(result.ErrMessage)`、HTTP status `200 OK`。

## Error Identity & Sentinels

`result.Error` 通过只比较 `Code` 实现 `errors.Is`，具体通过 `Error.Is(target)` 完成。两个 `result.Error` 只要 `Code` 相同就会匹配；`Message` 和 `Status` 不参与比较。这样动态格式化的错误也能匹配同一 code 的预置 sentinel。

```go
err := result.Errf("user %s missing", username,
    result.WithCode(result.ErrCodeRecordNotFound),
)

if errors.Is(err, result.ErrRecordNotFound) {
    // 相同业务错误码：2001
}
```

需要从 error chain 读取 `Code`、`Message` 或 `Status` 时，使用 `result.AsErr(err)`。`result.IsRecordNotFound(err)` 是 `errors.Is(err, result.ErrRecordNotFound)` 的便捷封装。

### 预置错误 sentinel

VEF 把常用错误暴露为可以直接返回的 `result.Error` 值：`result.ErrAccessDenied`, `result.ErrTooManyRequests`, `result.ErrRequestTimeout`, `result.ErrUnknown`, `result.ErrRecordNotFound`, `result.ErrRecordAlreadyExists`, `result.ErrForeignKeyViolation`, `result.ErrDangerousSQL`。数据库和 SQL 类业务失败刻意保持 HTTP `200 OK`；客户端应读取业务 `code` 区分失败类型。

| 错误值 | 业务码 | 默认 HTTP status | Message key |
| --- | --- | --- | --- |
| `result.ErrAccessDenied` | `result.ErrCodeAccessDenied`（`1100`） | `403` | `result.ErrMessageAccessDenied` |
| `result.ErrTooManyRequests` | `result.ErrCodeTooManyRequests`（`1401`） | `429` | `result.ErrMessageTooManyRequests` |
| `result.ErrRequestTimeout` | `result.ErrCodeRequestTimeout`（`1402`） | `408` | `result.ErrMessageRequestTimeout` |
| `result.ErrUnknown` | `result.ErrCodeUnknown`（`1900`） | `500` | `result.ErrMessageUnknown` |
| `result.ErrRecordNotFound` | `result.ErrCodeRecordNotFound`（`2001`） | `200` | `result.ErrMessageRecordNotFound` |
| `result.ErrRecordAlreadyExists` | `result.ErrCodeRecordAlreadyExists`（`2002`） | `200` | `result.ErrMessageRecordAlreadyExists` |
| `result.ErrForeignKeyViolation` | `result.ErrCodeForeignKeyViolation`（`2003`） | `200` | `result.ErrMessageForeignKeyViolation` |
| `result.ErrDangerousSQL` | `result.ErrCodeDangerousSQL`（`1600`） | `200` | `result.ErrMessageDangerousSQL` |
| `result.ErrNotImplemented(message)` | `result.ErrCodeNotImplemented`（`1500`） | `501` | caller-supplied message |

安全域 sentinel 位于 `security` 包（如 `security.ErrUnauthenticated`、`security.ErrTokenExpired` 等）——见下方 [Per-Module Error Tables](#per-module-error-tables)。

## The Result Code/Message Catalog

错误码按范围组织。`security` 和各模块专属错误位于各自包中（见下方 [Per-Module Error Tables](#per-module-error-tables)）。

### 跨模块错误码（`result` 包）

`ErrCode*` family 和 `ErrMessage*` family 共同定义了下表的跨模块业务码 / 消息对。

| Code | 常量 | 含义 |
| --- | --- | --- |
| `0` | `result.OkCode` | 成功 |
| `1100` | `result.ErrCodeAccessDenied` | 拒绝访问 |
| `1200` | `result.ErrCodeNotFound` | 资源未找到；standalone constant，用于 Fiber error mapping 或自定义 `Err(WithCode(...))` |
| `1300` | `result.ErrCodeUnsupportedMediaType` | 不支持的媒体类型；standalone constant，用于 Fiber error mapping 或自定义 `Err(WithCode(...))` |
| `1400` | `result.ErrCodeBadRequest` | 错误请求；standalone constant，用于 validation/API 包或自定义 `Err(WithCode(...))` |
| `1401` | `result.ErrCodeTooManyRequests` | 请求频率限制 |
| `1402` | `result.ErrCodeRequestTimeout` | 请求超时 |
| `1500` | `result.ErrCodeNotImplemented` | 未实现 |
| `1600` | `result.ErrCodeDangerousSQL` | 检测到危险 SQL |
| `1900` | `result.ErrCodeUnknown` | 未知或未包装错误 |
| `2000` | `result.ErrCodeDefault` | 默认业务错误 |
| `2001` | `result.ErrCodeRecordNotFound` | 记录未找到 |
| `2002` | `result.ErrCodeRecordAlreadyExists` | 记录已存在 |
| `2003` | `result.ErrCodeForeignKeyViolation` | 外键约束冲突 |

`ErrCodeNotFound`、`ErrCodeUnsupportedMediaType` 和 `ErrCodeBadRequest` 在这个包里没有预置 `result.Error` 值。`result.ErrCodeBadRequest`、`result.ErrCodeNotFound` 和 `result.ErrCodeUnsupportedMediaType` 是 exported building-block constants，供应用层 Fiber mapping 和各包专属错误使用。

### Message Keys

消息会在构造或错误处理时通过 `i18n` 模块查找。应用代码也可以直接把已经翻译好的字符串传给 `Err("...")`。

| 常量 | Key | 使用位置 |
| --- | --- | --- |
| `result.OkMessage` | `"ok"` | 默认 `Ok(...)` 消息 |
| `result.ErrMessage` | `"error"` | 默认 `Err(...)` 消息 |
| `result.ErrMessageUnknown` | `"unknown_error"` | `ErrUnknown` 和未映射错误 |
| `result.ErrMessageNotFound` | `"not_found"` | 应用层 Fiber `404` mapping |
| `result.ErrMessageTooManyRequests` | `"too_many_requests"` | `ErrTooManyRequests` |
| `result.ErrMessageAccessDenied` | `"access_denied"` | `ErrAccessDenied` 和应用层 Fiber `403` mapping |
| `result.ErrMessageUnsupportedMediaType` | `"unsupported_media_type"` | 应用层 Fiber `415` mapping |
| `result.ErrMessageRequestTimeout` | `"request_timeout"` | `ErrRequestTimeout` 和应用层 Fiber `408` mapping |
| `result.ErrMessageRecordNotFound` | `"record_not_found"` | `ErrRecordNotFound` |
| `result.ErrMessageRecordAlreadyExists` | `"record_already_exists"` | `ErrRecordAlreadyExists` |
| `result.ErrMessageForeignKeyViolation` | `"foreign_key_violation"` | `ErrForeignKeyViolation` |
| `result.ErrMessageDangerousSQL` | `"dangerous_sql"` | `ErrDangerousSQL` |

### 业务码范围

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

## Mapping Framework & Fiber Errors

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

### 错误解析顺序

运行时，VEF 会按以下顺序处理错误：

1. `fiber.Error`
2. `result.Error`
3. 其他未识别错误 -> `result.ErrUnknown`

这也是为什么对预期业务失败，应该优先返回显式的 `result.Error`，而不是不透明的普通错误。

## Per-Module Error Tables

VEF 在框架各处预置了一批 `result.Error`。模块专属错误位于各自模块包中，`result` 包只保留跨模块通用错误。

### 安全错误（`security` 包）

认证、签名、会话、challenge 流相关的错误位于 `github.com/coldsmirk/vef-framework-go/security`，有各自的 `ErrCodeXxx` 常量。认证错误使用 `1000-1029`（当前用到 `1026`）；challenge 错误预留 `1030-1039`，当前实际导出 `1031`、`1033-1038`；密码策略违规共用一个错误码 `1050`（完整的密码规则与锁定错误目录见 [登录加固](../security/login-hardening)）。

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
| `security.ErrAccountLocked(retryAfter)` | `security.ErrCodeAccountLocked`（1023） | `429` |
| `security.ErrTooManyConcurrentSessions` | `security.ErrCodeTooManyConcurrentSessions`（1024） | `403` |
| `security.ErrAPIKeyInvalid` | `security.ErrCodeAPIKeyInvalid`（1025） | `401` |
| `security.ErrBasicCredentialsInvalid` | `security.ErrCodeBasicCredentialsInvalid`（1026） | `401` |
| `security.ErrReservedPrincipal` | `security.ErrCodePrincipalInvalid`（1007，共用） | `401` |
| `security.ErrChallengeTokenInvalid` | `security.ErrCodeChallengeTokenInvalid`（1031） | `401` |
| `security.ErrChallengeTypeInvalid` | `security.ErrCodeChallengeTypeInvalid`（1033） | `400` |
| `security.ErrChallengeResolveFailed` | `security.ErrCodeChallengeResolveFailed`（1034） | `401` |
| `security.ErrOTPCodeRequired` | `security.ErrCodeOTPCodeRequired`（1035） | `400` |
| `security.ErrOTPCodeInvalid` | `security.ErrCodeOTPCodeInvalid`（1036） | `401` |
| `security.ErrNewPasswordRequired` | `security.ErrCodeNewPasswordRequired`（1037） | `400` |
| `security.ErrDepartmentRequired` | `security.ErrCodeDepartmentRequired`（1038） | `400` |
| `security.ErrCredentialsInvalid(message)` | `security.ErrCodeCredentialsInvalid`（1008） | `401` |
| `security.ErrPrincipalInvalid(message)` | `security.ErrCodePrincipalInvalid`（1007） | `401` |

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

> 这四个公开 `approval` sentinel 都是普通 Go 错误，**不**是 `result.Error`，没有 code/status 字段。内置审批资源响应使用 internal 的 40xxx result envelope 目录；完整 code 与 message key 见 [Approval 模块](../approval)。

## Practical Patterns

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

## Practical Advice

- 把 `result` 当作对外响应契约来看待
- 只要内置错误已经能表达你的语义，就优先复用
- 当客户端需要按不同失败类型做不同处理时，再定义明确的业务码
- 对预期业务失败，优先返回结构化 `result.Error`，不要返回随意拼接的字符串错误
- 除非你有意绕开框架契约，否则不要手工输出原始 JSON 响应

## Next Step

继续阅读 [Authentication](../security/authentication)，看认证失败是如何进入这一套结果模型的。
