---
sidebar_position: 9
---

# Result

The `result` package defines VEF's shared API response envelope and the
structured business-error type that the app layer turns into that envelope.

## Reviewed Public Surface

The current source audit for `github.com/coldsmirk/vef-framework-go/result`
covers 48 top-level exported symbols, 6 exported fields, and 4 exported
methods. The reviewed public-surface fingerprint is
`f91600ccb5960c2a405fb3ec5b2b84b38676c6488f4bf2dd45c8c22544b96892`.

Reviewed APIs:

| API | Contract |
| --- | --- |
| `result.Result` | Client response envelope type. Only this type has JSON field tags for the public wire shape. |
| `Result.Code` | Business result code, serialized as `code`. `0` means success. |
| `Result.Message` | Human-readable or i18n-resolved message, serialized as `message`. |
| `Result.Data` | Optional response payload, serialized as `data`; `nil` data is preserved as JSON `null`. |
| `Result.Response(ctx, status...)` | Sends the result as JSON; HTTP status defaults to `200 OK`, or the first supplied status value. |
| `Result.IsOk()` | Returns true only when `Code == result.OkCode`. |
| `result.Ok(dataOrOptions...)` | Builds a success `Result`; accepts zero args, one data arg, option-only, or data before `OkOption` values. |
| `result.OkOption` | Function type `func(*Result)` used by success-result options. |
| `result.WithMessage(message)` | Sets `Result.Message` exactly to `message`, including an empty string. |
| `result.WithMessagef(format, args...)` | Sets `Result.Message` with `fmt.Sprintf(format, args...)`. |
| `result.Error` | Structured application error with business code, message, and transport status. It is not the public JSON envelope. |
| `Error.Code` | Business error code used in the response envelope and in `errors.Is` comparisons. |
| `Error.Message` | Error message returned by `Error.Error()` and copied into the response envelope. |
| `Error.Status` | HTTP status used by the app error handler when converting the error into a `Result`. |
| `Error.Error()` | Implements `error` by returning `Message`. |
| `Error.Is(target)` | Matches another `result.Error` by `Code` only; `Message` and `Status` do not affect identity. |
| `result.Err(messageOrOptions...)` | Builds an `Error`; optional message string must be first, followed by `ErrOption` values. |
| `result.Errf(format, args...)` | Builds an `Error` with `fmt.Sprintf`; format args must come before any `ErrOption`. |
| `result.ErrOption` | Function type `func(*Error)` used by error options. |
| `result.WithCode(code)` | Sets `Error.Code`; default is `result.ErrCodeDefault` (`2000`). |
| `result.WithStatus(status)` | Sets `Error.Status`; default is `200 OK`. |
| `result.AsErr(err)` | Extracts a `result.Error` from an error chain. |
| `result.IsRecordNotFound(err)` | Uses `errors.Is(err, result.ErrRecordNotFound)`. |
| `result.ErrNotImplemented(message)` | Builds code `1500` with HTTP `501 Not Implemented`. |
| `result.OkCode` / `result.OkMessage` | Success code `0` and success message key `"ok"`. |
| `ErrCode*` family | Cross-cutting error-code constants listed below. |
| `ErrMessage*` family | Cross-cutting i18n message-key constants listed below. |
| `result.ErrAccessDenied`, `result.ErrTooManyRequests`, `result.ErrRequestTimeout`, `result.ErrUnknown`, `result.ErrRecordNotFound`, `result.ErrRecordAlreadyExists`, `result.ErrForeignKeyViolation`, `result.ErrDangerousSQL` | Ready-to-return `result.Error` values with fixed codes, message keys, and default HTTP statuses. |

## Response Shape

Every API response uses the same `result.Result` envelope:

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

`result.Error` intentionally has no JSON tags. Application code returns it as
an `error`; the app error handler converts it to `result.Result{Code, Message}`
and uses `Error.Status` as the HTTP status. Do not serialize `result.Error`
directly when you want the public `code / message / data` envelope.

## Success Responses

```go
import "github.com/coldsmirk/vef-framework-go/result"

// Empty success
return result.Ok().Response(ctx)

// Success with data
return result.Ok(user).Response(ctx)

// Success with a custom message
return result.Ok(result.WithMessage("Created successfully")).Response(ctx)

// Success with both data and a custom message
return result.Ok(user, result.WithMessage("User created")).Response(ctx)

// Success with a custom HTTP status (e.g. 201 Created)
return result.Ok(user).Response(ctx, 201)
```

`result.Ok(...)` accepts at most one data argument. If data is supplied, it must
come before `OkOption` values. Passing more than one data argument, or passing
data after an option, panics.

### Ok Options

| Option | Effect |
| --- | --- |
| `result.WithMessage(msg)` | Overrides the default success message with `msg`. |
| `result.WithMessagef(format, args...)` | Overrides the success message with `fmt.Sprintf(format, args...)`. |

The default success code is `result.OkCode` (`0`), and the default success
message is `i18n.T(result.OkMessage)`.

## Error Responses

`result.Err(...)` builds a `result.Error` value. Just return it; the framework's
error handler renders it into the response envelope:

```go
// Default business error (code 2000, message from i18n catalog)
return result.Err()

// Business error with a custom message
return result.Err("something went wrong")

// Business error with a specific code
return result.Err("not found", result.WithCode(result.ErrCodeRecordNotFound))

// Override the HTTP status while keeping the structured envelope
return result.Err("forbidden",
    result.WithCode(result.ErrCodeAccessDenied),
    result.WithStatus(fiber.StatusForbidden),
)
```

`result.Err(...)` accepts an optional message string only as the first argument.
Any following arguments must be `ErrOption` values; invalid argument types panic.
There is no error-message option: use the first `Err(...)` argument or
`Errf(...)` for the error message.

`result.Errf(format, args...)` is the formatted version:

```go
return result.Errf("user %s not found", username,
    result.WithCode(result.ErrCodeRecordNotFound),
)
```

`Errf` requires at least one format argument. `ErrOption` values must come after
all format arguments; putting an option before or between format arguments
panics.

### Err Options

| Option | Effect |
| --- | --- |
| `result.WithCode(code)` | Sets the business code. |
| `result.WithStatus(status)` | Sets the HTTP status code. |

Default `Err(...)` and `Errf(...)` values use `result.ErrCodeDefault` (`2000`),
message `i18n.T(result.ErrMessage)`, and HTTP status `200 OK`.

## Error Identity

`result.Error` implements `errors.Is` through `Error.Is(target)`. Two
`result.Error` values match when their `Code` values are equal; `Message` and
`Status` are ignored.

```go
err := result.Errf("user %s missing", username,
    result.WithCode(result.ErrCodeRecordNotFound),
)

if errors.Is(err, result.ErrRecordNotFound) {
    // same business error code: 2001
}
```

Use `result.AsErr(err)` when you need to read `Code`, `Message`, or `Status`
from an error chain. `result.IsRecordNotFound(err)` is a convenience wrapper
around `errors.Is(err, result.ErrRecordNotFound)`.

## Error Codes

Codes are organized by range. The `security` and per-module errors live in
their own packages (see [Error Handling](./error-handling) for the cross-module
table).

### Cross-Cutting Codes

| Code | Constant | Meaning |
| --- | --- | --- |
| `0` | `result.OkCode` | Success |
| `1100` | `result.ErrCodeAccessDenied` | Access denied |
| `1200` | `result.ErrCodeNotFound` | Resource not found; standalone constant used by Fiber error mapping or custom `Err(WithCode(...))` |
| `1300` | `result.ErrCodeUnsupportedMediaType` | Unsupported media type; standalone constant used by Fiber error mapping or custom `Err(WithCode(...))` |
| `1400` | `result.ErrCodeBadRequest` | Bad request; standalone constant used by validation/API packages or custom `Err(WithCode(...))` |
| `1401` | `result.ErrCodeTooManyRequests` | Rate limited |
| `1402` | `result.ErrCodeRequestTimeout` | Request timeout |
| `1500` | `result.ErrCodeNotImplemented` | Not implemented |
| `1600` | `result.ErrCodeDangerousSQL` | Dangerous SQL detected |
| `1900` | `result.ErrCodeUnknown` | Unknown or unwrapped error |
| `2000` | `result.ErrCodeDefault` | Default business error |
| `2001` | `result.ErrCodeRecordNotFound` | Record not found |
| `2002` | `result.ErrCodeRecordAlreadyExists` | Duplicate record |
| `2003` | `result.ErrCodeForeignKeyViolation` | Foreign-key constraint violation |

`ErrCodeNotFound`, `ErrCodeUnsupportedMediaType`, and `ErrCodeBadRequest` do not
have predefined `result.Error` values in this package. They are exported
building blocks for app-layer Fiber mappings and package-specific errors.

### Message Keys

Messages are looked up through the `i18n` module at construction or error
handling time. Application code may also pass an already-translated string
directly to `Err("...")`.

| Constant | Key | Used by |
| --- | --- | --- |
| `result.OkMessage` | `"ok"` | default `Ok(...)` message |
| `result.ErrMessage` | `"error"` | default `Err(...)` message |
| `result.ErrMessageUnknown` | `"unknown_error"` | `ErrUnknown` and unmapped errors |
| `result.ErrMessageNotFound` | `"not_found"` | app-layer Fiber `404` mapping |
| `result.ErrMessageTooManyRequests` | `"too_many_requests"` | `ErrTooManyRequests` |
| `result.ErrMessageAccessDenied` | `"access_denied"` | `ErrAccessDenied` and app-layer Fiber `403` mapping |
| `result.ErrMessageUnsupportedMediaType` | `"unsupported_media_type"` | app-layer Fiber `415` mapping |
| `result.ErrMessageRequestTimeout` | `"request_timeout"` | `ErrRequestTimeout` and app-layer Fiber `408` mapping |
| `result.ErrMessageRecordNotFound` | `"record_not_found"` | `ErrRecordNotFound` |
| `result.ErrMessageRecordAlreadyExists` | `"record_already_exists"` | `ErrRecordAlreadyExists` |
| `result.ErrMessageForeignKeyViolation` | `"foreign_key_violation"` | `ErrForeignKeyViolation` |
| `result.ErrMessageDangerousSQL` | `"dangerous_sql"` | `ErrDangerousSQL` |

## Pre-Built Error Sentinels

VEF exposes common errors as ready-to-return `result.Error` values. The database
and SQL-class business failures intentionally keep HTTP `200 OK`; clients should
read the business `code` to distinguish the failure.

| Error value | Business code | Default HTTP status | Message key |
| --- | --- | --- | --- |
| `result.ErrAccessDenied` | `result.ErrCodeAccessDenied` (`1100`) | `403` | `result.ErrMessageAccessDenied` |
| `result.ErrTooManyRequests` | `result.ErrCodeTooManyRequests` (`1401`) | `429` | `result.ErrMessageTooManyRequests` |
| `result.ErrRequestTimeout` | `result.ErrCodeRequestTimeout` (`1402`) | `408` | `result.ErrMessageRequestTimeout` |
| `result.ErrUnknown` | `result.ErrCodeUnknown` (`1900`) | `500` | `result.ErrMessageUnknown` |
| `result.ErrRecordNotFound` | `result.ErrCodeRecordNotFound` (`2001`) | `200` | `result.ErrMessageRecordNotFound` |
| `result.ErrRecordAlreadyExists` | `result.ErrCodeRecordAlreadyExists` (`2002`) | `200` | `result.ErrMessageRecordAlreadyExists` |
| `result.ErrForeignKeyViolation` | `result.ErrCodeForeignKeyViolation` (`2003`) | `200` | `result.ErrMessageForeignKeyViolation` |
| `result.ErrDangerousSQL` | `result.ErrCodeDangerousSQL` (`1600`) | `200` | `result.ErrMessageDangerousSQL` |
| `result.ErrNotImplemented(message)` | `result.ErrCodeNotImplemented` (`1500`) | `501` | caller-supplied message |

Security-domain sentinels live in the `security` package
(`security.ErrUnauthenticated`, `security.ErrTokenExpired`, and others). See
[Error Handling](./error-handling) for the cross-module list.
