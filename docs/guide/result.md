---
sidebar_position: 9
---

# Result

The `result` package provides the standard API response envelope used throughout the framework.

## Response Structure

Every API response follows this shape:

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

## Creating Success Responses

```go
import "github.com/coldsmirk/vef-framework-go/result"

// Simple success (no data)
return result.Ok().Response(ctx)

// Success with data
return result.Ok(user).Response(ctx)

// Success with custom message
return result.Ok(result.WithMessage("Created successfully")).Response(ctx)

// Success with data and custom message
return result.Ok(user, result.WithMessage("User created")).Response(ctx)

// Success with custom HTTP status
return result.Ok(user).Response(ctx, 201)
```

### Ok Options

| Option | Effect |
| --- | --- |
| `result.WithMessage(msg)` | Override the default success message |
| `result.WithCode(code)` | Override the default success code (0) |

## Creating Error Responses

```go
// Simple business error
return result.Err("something went wrong").Response(ctx)

// Error with custom code
return result.Err("not found", result.WithErrCode(result.ErrCodeRecordNotFound)).Response(ctx)

// Error with custom HTTP status
return result.Err("unauthorized").Response(ctx, 401)
```

### Err Options

| Option | Effect |
| --- | --- |
| `result.WithErrCode(code)` | Set a specific error code |
| `result.WithErrData(data)` | Attach data to the error response |

## Checking Results

```go
r := result.Ok(data)
r.IsOk() // true (code == 0)

r = result.Err("fail")
r.IsOk() // false
```

## Error Codes

### Authentication Errors (1000–1099)

| Code | Constant | Meaning |
| --- | --- | --- |
| 1000 | `ErrCodeUnauthenticated` | Not authenticated |
| 1001 | `ErrCodeUnsupportedAuthenticationType` | Unsupported auth type |
| 1002 | `ErrCodeTokenExpired` | Token expired |
| 1003 | `ErrCodeTokenInvalid` | Token invalid |
| 1004 | `ErrCodeTokenNotValidYet` | Token not valid yet |
| 1010 | `ErrCodePrincipalInvalid` | Principal invalid |
| 1011 | `ErrCodeCredentialsInvalid` | Credentials invalid |
| 1020 | `ErrCodeSignatureInvalid` | Signature invalid |

### Challenge Errors (1030–1039)

| Code | Constant | Meaning |
| --- | --- | --- |
| 1030 | `ErrCodeChallengeRequired` | Challenge required |
| 1035 | `ErrCodeOTPCodeRequired` | OTP code required |
| 1036 | `ErrCodeOTPCodeInvalid` | OTP code invalid |
| 1037 | `ErrCodeNewPasswordRequired` | New password required |

### Authorization Errors (1100–1199)

| Code | Constant | Meaning |
| --- | --- | --- |
| 1100 | `ErrCodeAccessDenied` | Access denied |

### Request Errors (1200–1499)

| Code | Constant | Meaning |
| --- | --- | --- |
| 1200 | `ErrCodeNotFound` | Resource not found |
| 1300 | `ErrCodeUnsupportedMediaType` | Unsupported media type |
| 1400 | `ErrCodeBadRequest` | Bad request |
| 1401 | `ErrCodeTooManyRequests` | Rate limited |
| 1402 | `ErrCodeRequestTimeout` | Request timeout |

### Business Errors (2000+)

| Code | Constant | Meaning |
| --- | --- | --- |
| 2000 | `ErrCodeDefault` | Generic business error |
| 2001 | `ErrCodeRecordNotFound` | Record not found in database |
| 2002 | `ErrCodeRecordAlreadyExists` | Duplicate record |
| 2003 | `ErrCodeForeignKeyViolation` | Foreign key constraint violation |

## I18n Integration

Error and success messages are automatically localized through the `i18n` module. The message keys (e.g., `"record_not_found"`, `"success"`) are looked up in the configured language bundle at runtime.

## Pre-built Error Constructors

The framework provides common error constructors:

```go
result.ErrRecordNotFound()       // code: 2001
result.ErrRecordAlreadyExists()  // code: 2002
result.ErrForeignKeyViolation()  // code: 2003
result.ErrAccessDenied()         // code: 1100
result.ErrUnauthenticated()      // code: 1000
result.ErrBadRequest(msg)        // code: 1400
result.ErrNotFound()             // code: 1200
result.ErrTooManyRequests()      // code: 1401
```
