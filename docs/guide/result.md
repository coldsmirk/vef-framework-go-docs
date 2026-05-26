---
sidebar_position: 9
---

# Result

The `result` package defines the unified API response envelope used across the framework.

## Response Shape

Every API response uses the same envelope:

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

### Ok options

| Option | Effect |
| --- | --- |
| `result.WithMessage(msg)` | Override the default success message. |
| `result.WithMessagef(format, args...)` | Same as `WithMessage` but with `fmt.Sprintf`-style formatting. |

> `result.Ok(...)` only accepts data **before** options, and at most one data argument. Mixing `WithCode` here would not compile — that option belongs to `Err`.

## Error Responses

`result.Err(...)` builds a `result.Error` value. Just return it — the framework's error handler renders it into the response envelope:

```go
// Default business error (code 2000, message from i18n catalog)
return result.Err()

// Business error with a custom message
return result.Err("something went wrong")

// Business error with a specific code
return result.Err("not found", result.WithCode(result.ErrCodeRecordNotFound))

// Override the HTTP status (still keeps the structured envelope)
return result.Err("forbidden",
    result.WithCode(result.ErrCodeAccessDenied),
    result.WithStatus(fiber.StatusForbidden),
)
```

`result.Errf(format, args...)` is the same with `fmt.Sprintf`-style formatting:

```go
return result.Errf("user %s not found", username,
    result.WithCode(result.ErrCodeRecordNotFound),
)
```

### Err options

| Option | Effect |
| --- | --- |
| `result.WithCode(code)` | Sets the business code (defaults to `ErrCodeDefault` = 2000). |
| `result.WithStatus(status)` | Sets the HTTP status code (defaults to 200). |

> `result.Error` is meant to be returned, not chained with `.Response(ctx)` — that method only exists on `result.Result`. The framework's error handler turns the returned `Error` into the JSON envelope automatically.

## Checking Success

`Result.IsOk()` reports whether the code indicates success:

```go
r := result.Ok(data)
if r.IsOk() {
    // ...
}
```

`result.Error` does not have `IsOk()` — it always represents a failure.

## Error Codes

Codes are organized by range. The `security` and per-module errors live in their own packages (see [Error Handling](./error-handling) for the full table).

### Cross-cutting codes (`result` package)

| Code | Constant | Meaning |
| --- | --- | --- |
| 1100 | `result.ErrCodeAccessDenied` | Access denied |
| 1200 | `result.ErrCodeNotFound` | Resource not found |
| 1300 | `result.ErrCodeUnsupportedMediaType` | Unsupported media type |
| 1400 | `result.ErrCodeBadRequest` | Bad request |
| 1401 | `result.ErrCodeTooManyRequests` | Rate limited |
| 1402 | `result.ErrCodeRequestTimeout` | Request timeout |
| 1500 | `result.ErrCodeNotImplemented` | Not implemented |
| 1600 | `result.ErrCodeDangerousSQL` | Dangerous SQL detected |
| 1900 | `result.ErrCodeUnknown` | Unknown / unwrapped error |
| 2000 | `result.ErrCodeDefault` | Default business error |
| 2001 | `result.ErrCodeRecordNotFound` | Record not found |
| 2002 | `result.ErrCodeRecordAlreadyExists` | Duplicate record |
| 2003 | `result.ErrCodeForeignKeyViolation` | Foreign-key constraint violation |

### Security codes (`security` package, 1000-1038)

Since v0.25 these moved out of `result` and now live in `security`. Examples:

| Code | Constant | Meaning |
| --- | --- | --- |
| 1000 | `security.ErrCodeUnauthenticated` | Not authenticated |
| 1001 | `security.ErrCodeUnsupportedAuthenticationType` | Unsupported auth type |
| 1002 | `security.ErrCodeTokenExpired` | Token expired |
| 1003 | `security.ErrCodeTokenInvalid` | Token invalid |
| 1004 | `security.ErrCodeTokenNotValidYet` | Token not valid yet |
| 1007 | `security.ErrCodePrincipalInvalid` | Principal invalid |
| 1008 | `security.ErrCodeCredentialsInvalid` | Credentials invalid |
| 1017 | `security.ErrCodeSignatureInvalid` | Signature invalid |
| 1030 | `security.ErrCodeChallengeRequired` | Challenge required |
| 1035 | `security.ErrCodeOTPCodeRequired` | OTP code required |
| 1036 | `security.ErrCodeOTPCodeInvalid` | OTP code invalid |
| 1037 | `security.ErrCodeNewPasswordRequired` | New password required |

See [Error Handling](./error-handling) for the full security code list (1000-1038), and the `storage` (2200+), `schema` (2300+), `monitor` (2100+) module codes.

## I18n Integration

Messages are looked up through the `i18n` module — `result.OkMessage` (`"ok"`), `result.ErrMessage` (`"error"`), and the various `ErrMessage*` constants resolve through the configured language bundle at runtime. Application code may pass an already-translated string directly to `Err("...")`.

## Pre-built Error Sentinels

VEF exposes the common errors as ready-to-return `result.Error` values. Return them directly:

```go
return result.ErrRecordNotFound        // code 2001
return result.ErrRecordAlreadyExists   // code 2002
return result.ErrForeignKeyViolation   // code 2003
return result.ErrAccessDenied          // code 1100
return result.ErrTooManyRequests       // code 1401
return result.ErrRequestTimeout        // code 1402
return result.ErrUnknown               // code 1900
return result.ErrDangerousSQL          // code 1600
```

A few sentinels are constructors because the message is application-defined:

```go
return result.ErrNotImplemented("multi-tenant export not yet supported")
```

Security-domain sentinels live in the `security` package (`security.ErrUnauthenticated`, `security.ErrTokenExpired`, …). See [Error Handling](./error-handling) for the complete cross-module list.
