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

VEF ships ready-made `result.Error` values across the framework. Starting from v0.25, module-specific errors live next to the module that owns them — the `result` package now only keeps cross-cutting errors. The codes themselves stay stable; just the import path changes.

### Cross-cutting errors (`result` package)

| Error value | Business code | Default HTTP status |
| --- | --- | --- |
| `result.ErrAccessDenied` | `result.ErrCodeAccessDenied` (1100) | `403` |
| `result.ErrTooManyRequests` | `result.ErrCodeTooManyRequests` (1401) | `429` |
| `result.ErrRequestTimeout` | `result.ErrCodeRequestTimeout` (1402) | `408` |
| `result.ErrUnknown` | `result.ErrCodeUnknown` (1900) | `500` |
| `result.ErrRecordNotFound` | `result.ErrCodeRecordNotFound` (2001) | `200` |
| `result.ErrRecordAlreadyExists` | `result.ErrCodeRecordAlreadyExists` (2002) | `200` |
| `result.ErrForeignKeyViolation` | `result.ErrCodeForeignKeyViolation` (2003) | `200` |
| `result.ErrDangerousSQL` | `result.ErrCodeDangerousSQL` (1600) | `200` |
| `result.ErrNotImplemented(message)` | `result.ErrCodeNotImplemented` (1500) | `501` |

### Security errors (`security` package)

Authentication, signature, and challenge flow errors live in `github.com/coldsmirk/vef-framework-go/security` with their own `ErrCodeXxx` constants (1000-1038 range).

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
| `security.ErrChallengeTokenInvalid` | `security.ErrCodeChallengeTokenInvalid` (1031) | `401` |
| `security.ErrChallengeTypeInvalid` | `security.ErrCodeChallengeTypeInvalid` (1033) | `400` |
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
| `monitor` | `monitor.ErrNotReady` | 2100 |
| `storage` | `storage.ErrInvalidFileKey`, `storage.ErrFileNotFound`, `storage.ErrFailedToGetFile` | 2200-2202 |
| `schema` | `schema.ErrTableNotFound` | 2300 |
| `crud` | `crud.ErrCodeProcessorInvalidReturn` (a `result.Error` code) plus plain sentinels `crud.ErrModelNoPrimaryKey`, `crud.ErrAuditUserCompositePK`, `crud.ErrSearchTypeMismatch`, `crud.ErrColumnNotFound` | 2400 |

> The `approval` package currently exports two `errors.New`-style sentinels (`approval.ErrCrossTenantAccess`, `approval.ErrInvalidBusinessIdentifier`). These are plain Go errors, **not** `result.Error` values, so they have no code/status fields — wrap them in `result.Err(...)` if you need to return them through the API envelope. See the [Approval module](../modules/approval) for the full domain error catalog (most live in the internal sub-packages).

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
| `401` | `security.ErrCodeUnauthenticated` | `security.ErrMessageUnauthenticated` |
| `403` | `result.ErrCodeAccessDenied` | `result.ErrMessageAccessDenied` |
| `404` | `result.ErrCodeNotFound` | `result.ErrMessageNotFound` |
| `415` | `result.ErrCodeUnsupportedMediaType` | `result.ErrMessageUnsupportedMediaType` |
| `408` | `result.ErrCodeRequestTimeout` | `result.ErrMessageRequestTimeout` |

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
