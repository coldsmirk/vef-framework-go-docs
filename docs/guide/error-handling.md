---
sidebar_position: 9
---

# Error Handling

VEF separates transport-level HTTP behavior from business-level result codes, but both are ultimately returned through the same `code / message / data` response envelope.

## Result Model Overview

VEF uses two closely related result types:

| Type | Purpose |
| --- | --- |
| `result.Result` | final response payload returned to clients |
| `result.Error` | structured error object used inside application code |

`result.Result` shape:

```json
{
  "code": 0,
  "message": "Success",
  "data": {}
}
```

## Successful Responses

Successful handlers usually return:

```go
return result.Ok(data).Response(ctx)
```

`result.Ok(...)` supports:

| Pattern | Meaning |
| --- | --- |
| `result.Ok()` | success without payload |
| `result.Ok(data)` | success with payload |
| `result.Ok(result.WithMessage(...))` | success with custom message |
| `result.Ok(data, result.WithMessage(...))` | success with payload and custom message |

## Structured Error Creation

For business failures, handlers usually return:

```go
return result.Err(
  "user already exists",
  result.WithCode(result.ErrCodeRecordAlreadyExists),
)
```

`result.Err(...)` supports:

| Pattern | Meaning |
| --- | --- |
| `result.Err()` | default business error |
| `result.Err("message")` | business error with custom message |
| `result.Err("message", result.WithCode(...))` | custom business code |
| `result.Err("message", result.WithStatus(...))` | custom HTTP status |
| `result.Err("message", result.WithCode(...), result.WithStatus(...))` | full override |

`result.Errf(...)` provides the same idea with formatted messages.

## Error Options

Available result options:

| Option | Applies to | Effect |
| --- | --- | --- |
| `result.WithCode(code)` | `result.Err(...)` | sets the business error code |
| `result.WithStatus(status)` | `result.Err(...)` | sets the HTTP status code |
| `result.WithMessage(message)` | `result.Ok(...)` | overrides the success message |
| `result.WithMessagef(format, ...)` | `result.Ok(...)` | formats the success message |

## Predefined Error Families

VEF ships a large set of predefined errors in the `result` package.

### Authentication errors

| Error value | Business code | Default HTTP status |
| --- | --- | --- |
| `result.ErrUnauthenticated` | `ErrCodeUnauthenticated` | `401` |
| `result.ErrTokenExpired` | `ErrCodeTokenExpired` | `401` |
| `result.ErrTokenInvalid` | `ErrCodeTokenInvalid` | `401` |
| `result.ErrTokenNotValidYet` | `ErrCodeTokenNotValidYet` | `401` |
| `result.ErrTokenInvalidIssuer` | `ErrCodeTokenInvalidIssuer` | `401` |
| `result.ErrTokenInvalidAudience` | `ErrCodeTokenInvalidAudience` | `401` |
| `result.ErrTokenMissingSubject` | `ErrCodeTokenMissingSubject` | `401` |
| `result.ErrTokenMissingTokenType` | `ErrCodeTokenMissingTokenType` | `401` |

### Signature or external-app auth errors

| Error value | Business code | Default HTTP status |
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

### Challenge flow errors

| Error value | Business code | Default HTTP status |
| --- | --- | --- |
| `result.ErrChallengeTokenInvalid` | `ErrCodeChallengeTokenInvalid` | `401` |
| `result.ErrChallengeTypeInvalid` | `ErrCodeChallengeTypeInvalid` | `400` |
| `result.ErrOTPCodeRequired` | `ErrCodeOTPCodeRequired` | `400` |
| `result.ErrOTPCodeInvalid` | `ErrCodeOTPCodeInvalid` | `401` |
| `result.ErrNewPasswordRequired` | `ErrCodeNewPasswordRequired` | `400` |
| `result.ErrDepartmentRequired` | `ErrCodeDepartmentRequired` | `400` |

### Authorization and request errors

| Error value | Business code | Default HTTP status |
| --- | --- | --- |
| `result.ErrAccessDenied` | `ErrCodeAccessDenied` | `403` |
| `result.ErrTooManyRequests` | `ErrCodeTooManyRequests` | `429` |
| `result.ErrRequestTimeout` | `ErrCodeRequestTimeout` | `408` |
| `result.ErrUnknown` | `ErrCodeUnknown` | `500` |

### Business errors

| Error value | Business code | Default HTTP status |
| --- | --- | --- |
| `result.ErrRecordNotFound` | `ErrCodeRecordNotFound` | `200` |
| `result.ErrRecordAlreadyExists` | `ErrCodeRecordAlreadyExists` | `200` |
| `result.ErrForeignKeyViolation` | `ErrCodeForeignKeyViolation` | `200` |
| `result.ErrDangerousSQL` | `ErrCodeDangerousSQL` | `200` |

### Error constructors

These helpers create structured errors with specific semantics:

| Constructor | Typical output |
| --- | --- |
| `result.ErrNotImplemented(message)` | `501 Not Implemented` |
| `result.ErrCredentialsInvalid(message)` | `401 Unauthorized` with credentials-invalid business code |
| `result.ErrPrincipalInvalid(message)` | `401 Unauthorized` with principal-invalid business code |

## Business Codes

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

## Fiber Error Mapping

The app layer maps selected `fiber.Error` values into structured result payloads.

Current built-in mappings:

| Fiber HTTP status | Result code | Message key |
| --- | --- | --- |
| `401` | `ErrCodeUnauthenticated` | `ErrMessageUnauthenticated` |
| `403` | `ErrCodeAccessDenied` | `ErrMessageAccessDenied` |
| `404` | `ErrCodeNotFound` | `ErrMessageNotFound` |
| `415` | `ErrCodeUnsupportedMediaType` | `ErrMessageUnsupportedMediaType` |
| `408` | `ErrCodeRequestTimeout` | `ErrMessageRequestTimeout` |

If a `fiber.Error` status code is not mapped, VEF logs it and falls back to the generic unknown error result.

## Error Resolution Order

At runtime, VEF resolves errors in this order:

1. `fiber.Error`
2. `result.Error`
3. unknown or unwrapped error -> `result.ErrUnknown`

This is why returning explicit `result.Error` values is better than returning opaque errors for domain failures.

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
