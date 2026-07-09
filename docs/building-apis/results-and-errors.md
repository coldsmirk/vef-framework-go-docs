---
sidebar_position: 5
---

# Results and Errors

VEF separates transport-level HTTP behavior from business-level result codes, but both are ultimately returned through the same `code / message / data` response envelope. The `result` package (`github.com/coldsmirk/vef-framework-go/result`) defines that envelope and the structured business-error type application code returns to produce it.

## Success & Error Envelope

VEF uses two closely related result types:

| Type | Purpose |
| --- | --- |
| `result.Result` | the final response payload returned to clients |
| `result.Error` | the structured error object used inside application code |

`result.Result` shape:

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

| Field/Method | Description |
| --- | --- |
| `Result.Code` | Business result code, serialized as `code`. `0` means success. |
| `Result.Message` | Human-readable or i18n-resolved message, serialized as `message`. |
| `Result.Data` | Optional response payload, serialized as `data`; `nil` data is preserved as JSON `null`. |
| `Result.Response(ctx, status...)` | Sends the result as JSON; HTTP status defaults to `200 OK`, or the first supplied status value. |
| `Result.IsOk()` | Returns true only when `Code == result.OkCode`. |

`result.Error` intentionally has no JSON tags — it is not the public JSON envelope. Application code returns it as an `error`; the app error handler extracts `Code`, `Message`, and `Status`, sends `Code` and `Message` in the `result.Result` envelope, and uses `Status` as the HTTP status. `result.Error` is not serialized directly as the client response. Do not serialize `result.Error` directly when you want the public `code / message / data` envelope.

```go
type Error struct {
    Code    int
    Message string
    Status  int
}
```

| Field/Method | Description |
| --- | --- |
| `Error.Code` | Business error code used in the response envelope and in `errors.Is` comparisons. |
| `Error.Message` | Error message returned by `Error.Error()` and copied into the response envelope. |
| `Error.Status` | HTTP status used by the app error handler when converting the error into a `Result`. |
| `Error.Error()` | Implements `error` by returning `Message`. |
| `Error.Is(target)` | Matches another `result.Error` by `Code` only. |

## Building Results (Ok/Err/Options)

### Success responses

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

`result.Ok(dataOrOptions...)` supports:

| Pattern | Meaning |
| --- | --- |
| `result.Ok()` | success without payload |
| `result.Ok(data)` | success with payload |
| `result.Ok(result.WithMessage(...))` | success with custom message |
| `result.Ok(data, result.WithMessage(...))` | success with payload and custom message |

`result.Ok(...)` accepts at most one data argument. If data is supplied, it must come before `OkOption` values. Passing more than one data argument, or passing data after an option, panics.

`result.OkOption` is the option function type (`func(*Result)`):

| Option | Effect |
| --- | --- |
| `result.WithMessage(message)` | sets `Result.Message` exactly to `message`, including an empty string |
| `result.WithMessagef(format, args...)` | sets `Result.Message` with `fmt.Sprintf(format, args...)` |

The default success code is `result.OkCode` (`0`), and the default success message is `i18n.T(result.OkMessage)`. Together `result.OkCode` / `result.OkMessage` define the default success envelope.

### Error responses

`result.Err(...)` builds a `result.Error` value. Just return it; the framework's error handler renders it into the response envelope:

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

`result.Err(messageOrOptions...)` supports:

| Pattern | Meaning |
| --- | --- |
| `result.Err()` | default business error |
| `result.Err("message")` | business error with custom message |
| `result.Err("message", result.WithCode(...))` | custom business code |
| `result.Err("message", result.WithStatus(...))` | custom HTTP status |
| `result.Err("message", result.WithCode(...), result.WithStatus(...))` | full override |

The optional message string must be the first `Err(...)` argument. Any following arguments must be `ErrOption` values; invalid argument types panic. There is no error-message option. There is no message option for `result.Error`: use the first `Err(...)` argument or `Errf(...)` for the error message.

`result.Errf(format, args...)` is the formatted version:

```go
return result.Errf("user %s not found", username,
    result.WithCode(result.ErrCodeRecordNotFound),
)
```

`Errf` requires at least one format argument. `ErrOption` values must come after all format arguments; putting an option before or between format arguments panics.

`result.ErrOption` is the option function type (`func(*Error)`):

| Option | Effect |
| --- | --- |
| `result.WithCode(code)` | sets `Error.Code`; default is `result.ErrCodeDefault` (`2000`) |
| `result.WithStatus(status)` | sets `Error.Status`; default is `200 OK` |

Default `Err(...)` and `Errf(...)` values use `result.ErrCodeDefault` (`2000`), message `i18n.T(result.ErrMessage)`, and HTTP status `200 OK`.

## Error Identity & Sentinels

`result.Error` implements `errors.Is` by comparing `Code` only, through `Error.Is(target)`. Two `result.Error` values match when their `Code` values are equal; `Message` and `Status` are ignored. This lets dynamically formatted errors still match a predefined sentinel with the same code.

```go
err := result.Errf("user %s missing", username,
    result.WithCode(result.ErrCodeRecordNotFound),
)

if errors.Is(err, result.ErrRecordNotFound) {
    // same business error code: 2001
}
```

Use `result.AsErr(err)` when you need to read `Code`, `Message`, or `Status` from an error chain. `result.IsRecordNotFound(err)` is a convenience wrapper around `errors.Is(err, result.ErrRecordNotFound)`.

### Pre-built error sentinels

VEF exposes common errors as ready-to-return `result.Error` values: `result.ErrAccessDenied`, `result.ErrTooManyRequests`, `result.ErrRequestTimeout`, `result.ErrUnknown`, `result.ErrRecordNotFound`, `result.ErrRecordAlreadyExists`, `result.ErrForeignKeyViolation`, `result.ErrDangerousSQL`. The database and SQL-class business failures intentionally keep HTTP `200 OK`; clients should read the business `code` to distinguish the failure.

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

Security-domain sentinels live in the `security` package (`security.ErrUnauthenticated`, `security.ErrTokenExpired`, and others) — see [Per-module error tables](#per-module-error-tables) below.

## The Result Code/Message Catalog

Codes are organized by range. The `security` and per-module errors live in their own packages (see [Per-module error tables](#per-module-error-tables) below).

### Cross-cutting codes (`result` package)

The `ErrCode*` family and `ErrMessage*` family together define the cross-cutting business code/message pairs below.

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

`ErrCodeNotFound`, `ErrCodeUnsupportedMediaType`, and `ErrCodeBadRequest` do not have predefined `result.Error` values in this package. `result.ErrCodeBadRequest`, `result.ErrCodeNotFound`, and `result.ErrCodeUnsupportedMediaType` are exported building-block constants for app-layer Fiber mappings and package-specific errors.

### Message keys

Messages are looked up through the `i18n` module at construction or error handling time. Application code may also pass an already-translated string directly to `Err("...")`.

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

### Business code ranges

Selected result code ranges:

| Range | Meaning |
| --- | --- |
| `0` | success |
| `1000-1099` | authentication and challenge errors |
| `1100-1199` | authorization errors |
| `1200-1499` | resource, media type, and request errors |
| `1500-1699` | not implemented and SQL-related errors |
| `1900-1999` | unknown errors |
| `2000+` | business errors |

## Mapping Framework & Fiber Errors

The app layer maps selected `fiber.Error` values into structured result payloads.

Current built-in mappings:

| Fiber HTTP status | Result code | Message key |
| --- | --- | --- |
| `401` | `security.ErrCodeUnauthenticated` | `security.ErrMessageUnauthenticated` |
| `403` | `result.ErrCodeAccessDenied` | `result.ErrMessageAccessDenied` |
| `404` | `result.ErrCodeNotFound` | `result.ErrMessageNotFound` |
| `415` | `result.ErrCodeUnsupportedMediaType` | `result.ErrMessageUnsupportedMediaType` |
| `408` | `result.ErrCodeRequestTimeout` | `result.ErrMessageRequestTimeout` |

If a `fiber.Error` status code is not mapped, VEF logs it and falls back to the generic unknown error result.

### Error resolution order

At runtime, VEF resolves errors in this order:

1. `fiber.Error`
2. `result.Error`
3. unknown or unwrapped error -> `result.ErrUnknown`

This is why returning explicit `result.Error` values is better than returning opaque errors for domain failures.

## Per-Module Error Tables

VEF ships ready-made `result.Error` values across the framework. Starting from v0.25, module-specific errors live next to the module that owns them — the `result` package now only keeps cross-cutting errors. The codes themselves stay stable; just the import path changes.

### Security errors (`security` package)

Authentication, signature, session, and challenge flow errors live in `github.com/coldsmirk/vef-framework-go/security` with their own `ErrCodeXxx` constants. Authentication uses `1000-1024`; challenge errors reserve `1030-1039` and currently export `1031`, `1033-1038`; password-policy violations share a single code, `1050` (see [Login Hardening](../security/login-hardening) for the full password-rule and lockout error catalog).

| Error value | Business code | Default HTTP status |
| --- | --- | --- |
| `security.ErrUnauthenticated` | `security.ErrCodeUnauthenticated` (1000) | `401` |
| `security.ErrTokenExpired` | `security.ErrCodeTokenExpired` (1002) | `401` |
| `security.ErrTokenInvalid` | `security.ErrCodeTokenInvalid` (1003) | `401` |
| `security.ErrTokenNotValidYet` | `security.ErrCodeTokenNotValidYet` (1004) | `401` |
| `security.ErrTokenInvalidIssuer` | `security.ErrCodeTokenInvalidIssuer` (1005) | `401` |
| `security.ErrTokenInvalidAudience` | `security.ErrCodeTokenInvalidAudience` (1006) | `401` |
| `security.ErrAppIDRequired` | `security.ErrCodeAppIDRequired` (1009) | `401` |
| `security.ErrTimestampRequired` | `security.ErrCodeTimestampRequired` (1010) | `401` |
| `security.ErrSignatureRequired` | `security.ErrCodeSignatureRequired` (1011) | `401` |
| `security.ErrTimestampInvalid` | `security.ErrCodeTimestampInvalid` (1012) | `401` |
| `security.ErrSignatureExpired` | `security.ErrCodeSignatureExpired` (1013) | `401` |
| `security.ErrSignatureInvalid` | `security.ErrCodeSignatureInvalid` (1017) | `401` |
| `security.ErrExternalAppNotFound` | `security.ErrCodeExternalAppNotFound` (1014) | `401` |
| `security.ErrExternalAppDisabled` | `security.ErrCodeExternalAppDisabled` (1015) | `401` |
| `security.ErrIPNotAllowed` | `security.ErrCodeIPNotAllowed` (1016) | `401` |
| `security.ErrNonceRequired` | `security.ErrCodeNonceRequired` (1018) | `401` |
| `security.ErrNonceInvalid` | `security.ErrCodeNonceInvalid` (1019) | `401` |
| `security.ErrNonceAlreadyUsed` | `security.ErrCodeNonceAlreadyUsed` (1020) | `401` |
| `security.ErrAuthHeaderMissing` | `security.ErrCodeAuthHeaderMissing` (1021) | `401` |
| `security.ErrAuthHeaderInvalid` | `security.ErrCodeAuthHeaderInvalid` (1022) | `401` |
| `security.ErrAccountLocked(retryAfter)` | `security.ErrCodeAccountLocked` (1023) | `429` |
| `security.ErrTooManyConcurrentSessions` | `security.ErrCodeTooManyConcurrentSessions` (1024) | `403` |
| `security.ErrChallengeTokenInvalid` | `security.ErrCodeChallengeTokenInvalid` (1031) | `401` |
| `security.ErrChallengeTypeInvalid` | `security.ErrCodeChallengeTypeInvalid` (1033) | `400` |
| `security.ErrChallengeResolveFailed` | `security.ErrCodeChallengeResolveFailed` (1034) | `401` |
| `security.ErrOTPCodeRequired` | `security.ErrCodeOTPCodeRequired` (1035) | `400` |
| `security.ErrOTPCodeInvalid` | `security.ErrCodeOTPCodeInvalid` (1036) | `401` |
| `security.ErrNewPasswordRequired` | `security.ErrCodeNewPasswordRequired` (1037) | `400` |
| `security.ErrDepartmentRequired` | `security.ErrCodeDepartmentRequired` (1038) | `400` |
| `security.ErrCredentialsInvalid(message)` | `security.ErrCodeCredentialsInvalid` (1008) | `401` |
| `security.ErrPrincipalInvalid(message)` | `security.ErrCodePrincipalInvalid` (1007) | `401` |

> v0.25.1 dropped the unused `ErrTokenMissingSubject` / `ErrTokenMissingTokenType` sentinels and compacted the surrounding codes. Bumps from older snapshots have no compatibility shim — update call sites to the current names.

### Other module errors

| Module package | Error values | Code range |
| --- | --- | --- |
| `api` | `api.ErrInvalidRequestParams`, `api.ErrInvalidRequestMeta` | 1400 (`result.ErrCodeBadRequest`) |
| `monitor` | `monitor.ErrNotReady`, `monitor.ErrCollectionFailed` | 2100-2101 |
| `storage` | `storage.ErrInvalidFileKey`, `storage.ErrFileNotFound`, `storage.ErrFailedToGetFile`, and multipart upload / claim errors such as `storage.ErrUploadRequiresMultipart`, `storage.ErrUploadPartsIncomplete`, and `storage.ErrAbortFailed` | 2200-2219 |
| `schema` | `schema.ErrTableNotFound` | 2300 |
| `crud` | `crud.ErrCodeProcessorInvalidReturn`, CRUD import/export and primary-key result errors, plus plain sentinels such as `crud.ErrModelNoPrimaryKey` and `crud.ErrAuditUserCompositePK` | 2400-2410 |
| `expression` | `expression.ErrEvaluationFailed` | 2500 |
| `approval` | public plain sentinels: `approval.ErrCrossTenantAccess`, `approval.ErrInvalidBusinessIdentifier`, `approval.ErrUnknownNodeKind`, `approval.ErrNodeDataUnmarshal`; built-in approval resources return internal `result.Error` values | 40001-40702 |

> The four public `approval` sentinels are plain Go errors, **not** `result.Error` values, so they have no code/status fields. Built-in approval resource responses use the internal 40xxx result-envelope catalog instead; see the [Approval module](../approval) for the full code and message-key table.

## Practical Patterns

### Success with payload

```go
return result.Ok(user).Response(ctx)
```

### Success with custom message

```go
return result.Ok(
  user,
  result.WithMessage("user synced"),
).Response(ctx)
```

### Business error with code

```go
return result.Err(
  "user already exists",
  result.WithCode(result.ErrCodeRecordAlreadyExists),
)
```

### Explicit HTTP status override

```go
return result.Err(
  "forbidden",
  result.WithCode(result.ErrCodeAccessDenied),
  result.WithStatus(fiber.StatusForbidden),
)
```

## Practical Advice

- think of `result` as the public response contract
- use predefined result errors when they already match the scenario
- use domain-specific business codes when the client must react differently
- prefer structured `result.Error` values over ad hoc string errors for expected business failures
- avoid manually writing raw JSON responses unless you are intentionally bypassing the result contract

## Next Step

Read [Authentication](../security/authentication) to see how auth failures flow into this result model.
