---
sidebar_position: 9
---

# Result

`result` 包定义 VEF 共享的 API 响应信封，以及会被应用层转换成该信封的结构化业务错误类型。

## 已审查公开 Surface

当前源码审计覆盖 `github.com/coldsmirk/vef-framework-go/result` 的 48 个
top-level exported symbols、6 个 exported fields、4 个 exported methods。
已审查 public-surface fingerprint 是
`f91600ccb5960c2a405fb3ec5b2b84b38676c6488f4bf2dd45c8c22544b96892`。

已审查 API：

| API | Contract |
| --- | --- |
| `result.Result` | 返回给客户端的响应信封类型。只有这个类型带有 public wire shape 的 JSON field tags。 |
| `Result.Code` | 业务结果码，序列化为 `code`；`0` 表示成功。 |
| `Result.Message` | 面向用户或经 i18n 解析后的消息，序列化为 `message`。 |
| `Result.Data` | 可选响应载荷，序列化为 `data`；`nil` data 会保留为 JSON `null`。 |
| `Result.Response(ctx, status...)` | 以 JSON 发送结果；HTTP status 默认 `200 OK`，如果传入 status 则使用第一个值。 |
| `Result.IsOk()` | 仅当 `Code == result.OkCode` 时返回 true。 |
| `result.Ok(dataOrOptions...)` | 构造成功 `Result`；支持无参数、一个 data 参数、仅 option、或 data 位于 `OkOption` 之前。 |
| `result.OkOption` | 成功结果 option 类型：`func(*Result)`。 |
| `result.WithMessage(message)` | 将 `Result.Message` 精确设置为 `message`，空字符串也会生效。 |
| `result.WithMessagef(format, args...)` | 使用 `fmt.Sprintf(format, args...)` 设置 `Result.Message`。 |
| `result.Error` | 带业务码、消息和传输层 status 的结构化应用错误。它不是 public JSON 信封。 |
| `Error.Code` | 业务错误码，会进入响应信封，也用于 `errors.Is` 比较。 |
| `Error.Message` | `Error.Error()` 返回的错误消息，也会复制进响应信封。 |
| `Error.Status` | 应用错误处理器把错误转换成 `Result` 时使用的 HTTP status。 |
| `Error.Error()` | 通过返回 `Message` 实现 `error` 接口。 |
| `Error.Is(target)` | 只按 `Code` 匹配另一个 `result.Error`；`Message` 和 `Status` 不参与 identity 比较。 |
| `result.Err(messageOrOptions...)` | 构造 `Error`；可选 message string 必须是第一个参数，后面只能跟 `ErrOption`。 |
| `result.Errf(format, args...)` | 使用 `fmt.Sprintf` 构造 `Error`；format args 必须位于所有 `ErrOption` 之前。 |
| `result.ErrOption` | 错误 option 类型：`func(*Error)`。 |
| `result.WithCode(code)` | 设置 `Error.Code`；默认是 `result.ErrCodeDefault`（`2000`）。 |
| `result.WithStatus(status)` | 设置 `Error.Status`；默认是 `200 OK`。 |
| `result.AsErr(err)` | 从 error chain 中提取 `result.Error`。 |
| `result.IsRecordNotFound(err)` | 使用 `errors.Is(err, result.ErrRecordNotFound)` 判断记录不存在错误。 |
| `result.ErrNotImplemented(message)` | 构造 code `1500`、HTTP `501 Not Implemented` 的错误。 |
| `result.OkCode` / `result.OkMessage` | 成功码 `0` 和成功消息 key `"ok"`。 |
| `ErrCode*` family | 下方列出的跨模块错误码常量。 |
| `ErrMessage*` family | 下方列出的跨模块 i18n message key 常量。 |
| `result.ErrAccessDenied`, `result.ErrTooManyRequests`, `result.ErrRequestTimeout`, `result.ErrUnknown`, `result.ErrRecordNotFound`, `result.ErrRecordAlreadyExists`, `result.ErrForeignKeyViolation`, `result.ErrDangerousSQL` | 可直接返回的 `result.Error` 值，带固定 code、message key 和默认 HTTP status。 |

## 响应结构

所有 API 响应都使用同一个 `result.Result` 信封：

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

`result.Error` 刻意没有 JSON tags。应用代码把它作为 `error` 返回；应用错误处理器会把它转换成
`result.Result{Code, Message}`，并使用 `Error.Status` 作为 HTTP status。
如果目标是公开的 `code / message / data` 响应信封，不要直接序列化 `result.Error`。

## 创建成功响应

```go
import "github.com/coldsmirk/vef-framework-go/result"

// 简单成功（无数据）
return result.Ok().Response(ctx)

// 带数据的成功
return result.Ok(user).Response(ctx)

// 自定义消息的成功
return result.Ok(result.WithMessage("创建成功")).Response(ctx)

// 带数据和自定义消息
return result.Ok(user, result.WithMessage("用户已创建")).Response(ctx)

// 自定义 HTTP 状态码
return result.Ok(user).Response(ctx, 201)
```

`result.Ok(...)` 最多接受一个 data 参数。如果提供 data，它必须位于 `OkOption`
之前。传入多个 data 参数，或把 data 放在 option 之后，都会 panic。

### Ok 选项

| 选项 | 效果 |
| --- | --- |
| `result.WithMessage(msg)` | 使用 `msg` 覆盖默认成功消息。 |
| `result.WithMessagef(format, args...)` | 使用 `fmt.Sprintf(format, args...)` 覆盖成功消息。 |

默认成功码是 `result.OkCode`（`0`），默认成功消息是 `i18n.T(result.OkMessage)`。

## 创建错误响应

`result.Err(...)` 会构造 `result.Error`。直接返回它，框架的错误处理器会把它渲染为统一响应信封：

```go
// 默认业务错误（code 2000，message 来自 i18n catalog）
return result.Err()

// 自定义错误消息
return result.Err("出错了")

// 带自定义错误码
return result.Err("未找到", result.WithCode(result.ErrCodeRecordNotFound))

// 自定义 HTTP 状态码，同时保留结构化响应信封
return result.Err("未授权",
    result.WithCode(result.ErrCodeAccessDenied),
    result.WithStatus(fiber.StatusForbidden),
)
```

`result.Err(...)` 的可选 message string 只能作为第一个参数。后续参数必须是
`ErrOption`；其他类型会 panic。错误消息没有 option 形式：要用第一个
`Err(...)` 参数，或使用 `Errf(...)` 设置错误消息。

`result.Errf(format, args...)` 是格式化版本：

```go
return result.Errf("user %s not found", username,
    result.WithCode(result.ErrCodeRecordNotFound),
)
```

`Errf` 至少需要一个 format arg。`ErrOption` 必须放在所有 format args 之后；
把 option 放在 format args 之前或中间都会 panic。

### Err 选项

| 选项 | 效果 |
| --- | --- |
| `result.WithCode(code)` | 设置业务错误码。 |
| `result.WithStatus(status)` | 设置 HTTP status code。 |

默认 `Err(...)` 和 `Errf(...)` 使用 `result.ErrCodeDefault`（`2000`）、
message `i18n.T(result.ErrMessage)`、HTTP status `200 OK`。

## 错误身份判断

`result.Error` 通过 `Error.Is(target)` 实现 `errors.Is`。两个
`result.Error` 只要 `Code` 相同就会匹配；`Message` 和 `Status` 不参与比较。

```go
err := result.Errf("user %s missing", username,
    result.WithCode(result.ErrCodeRecordNotFound),
)

if errors.Is(err, result.ErrRecordNotFound) {
    // 相同业务错误码：2001
}
```

需要从 error chain 读取 `Code`、`Message` 或 `Status` 时，使用
`result.AsErr(err)`。`result.IsRecordNotFound(err)` 是
`errors.Is(err, result.ErrRecordNotFound)` 的便捷封装。

## 错误码

错误码按范围组织。`security` 和各模块专属错误位于各自包中，跨模块列表见[错误处理](./error-handling)。

### 跨模块错误码

| 错误码 | 常量 | 含义 |
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
| `2001` | `result.ErrCodeRecordNotFound` | 数据库记录未找到 |
| `2002` | `result.ErrCodeRecordAlreadyExists` | 记录已存在 |
| `2003` | `result.ErrCodeForeignKeyViolation` | 外键约束冲突 |

`ErrCodeNotFound`、`ErrCodeUnsupportedMediaType` 和 `ErrCodeBadRequest` 在这个包里没有预置
`result.Error` 值。它们是给应用层 Fiber mapping 和各包专属错误使用的 building blocks。

### Message Keys

消息会在构造或错误处理时通过 `i18n` 模块查找。应用代码也可以直接把已经翻译好的字符串传给
`Err("...")`。

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

## 预置错误 Sentinel

VEF 把常用错误暴露为可以直接返回的 `result.Error` 值。数据库和 SQL 类业务失败刻意保持
HTTP `200 OK`；客户端应读取业务 `code` 区分失败类型。

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

安全域 sentinel 位于 `security` 包（如 `security.ErrUnauthenticated`、
`security.ErrTokenExpired` 等）。跨模块列表见[错误处理](./error-handling)。
