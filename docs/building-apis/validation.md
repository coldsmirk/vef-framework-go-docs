---
sidebar_position: 4
---

# Validation

VEF uses `go-playground/validator` as the base validation engine and layers framework-specific behavior on top:

- translated error messages
- `label` / `label_i18n` support
- custom validation rules
- pointer-based nullable field support through the standard `omitempty` flow

## API Reference

| API | Contract |
| --- | --- |
| `validator.Validate(value)` | Runs the package-level validator against `value`; success returns `nil`, validation failure returns the first translated error as a bad-request `result.Error` |
| `validator.RegisterValidationRules(rules...)` | Registers each `ValidationRule` with the shared validator and both built-in translators; returns the first registration error |
| `validator.RegisterTypeFunc(fn, types...)` | Registers a custom type extractor by delegating to `RegisterCustomTypeFunc` |
| `validator.CustomTypeFunc` | Alias-compatible function type `func(field reflect.Value) any` for custom type extraction |
| `validator.ValidationRule` | Struct used to define custom validation tags and translation behavior |
| `ValidationRule.RuleTag` | Validator tag name passed to `RegisterValidation` and `RegisterTranslation` |
| `ValidationRule.ErrMessageTemplate` | Fallback message template registered with go-playground translators |
| `ValidationRule.ErrMessageI18nKey` | Optional framework i18n key; if it resolves to a real message, that message wins over the fallback template |
| `ValidationRule.Validate` | Actual rule predicate over go-playground `FieldLevel` values |
| `ValidationRule.ParseParam` | Extracts placeholder values from go-playground `FieldError` values for message rendering |
| `ValidationRule.CallValidationEvenIfNull` | Passed through to go-playground `RegisterValidation` |

## Validation Entry Point

Typed params and typed meta structs are validated automatically after decoding. The framework calls:

```go
validator.Validate(value)
```

That means handler parameters using typed request structs get validation by default without manual calls inside the handler.

## Standard Validator Tags

VEF inherits the standard `go-playground/validator` tag set. Common examples include:

| Tag | Meaning | Example |
| --- | --- | --- |
| `required` | field must be present and non-zero | `validate:"required"` |
| `email` | field must be a valid email | `validate:"required,email"` |
| `min` | numeric or length lower bound | `validate:"min=1"` |
| `max` | numeric or length upper bound | `validate:"max=32"` |
| `oneof` | field must match one of several values | `validate:"oneof=admin user guest"` |
| `len` | exact length | `validate:"len=32"` |
| `omitempty` | skip later validations when empty | `validate:"omitempty,email"` |
| `dive` | validate each slice or map item | `validate:"required,dive"` |

VEF does not redefine those upstream rules; it adds framework-specific rules on top of them.

## Label Resolution

Validation errors use field labels instead of raw Go field names when possible.

Resolution order:

| Source | Effect |
| --- | --- |
| `label:"..."` | uses the explicit label text directly |
| `label_i18n:"..."` | resolves the label through `i18n.T`; if the translation is missing, the i18n key may appear in the error message |
| neither tag present | uses the Go field name |

Example:

```go
type UserParams struct {
	api.P

	Username string `json:"username" validate:"required" label:"Username"`
	Phone    string `json:"phone" validate:"phone_number" label_i18n:"user_phone"`
}
```

## Built-In Custom Rules

VEF currently ships these custom validation rules:

| Rule tag | Expected field type | Meaning | Example |
| --- | --- | --- | --- |
| `phone_number` | `string` | validates a Mainland China mobile number | `validate:"phone_number"` |
| `dec_min=<value>` | `decimal.Decimal` | decimal value must be greater than or equal to the threshold | `validate:"dec_min=0"` |
| `dec_max=<value>` | `decimal.Decimal` | decimal value must be less than or equal to the threshold | `validate:"dec_max=999.99"` |
| `alphanum_us` | `string` | letters, numbers, underscore only | `validate:"alphanum_us"` |
| `alphanum_us_slash` | `string` | letters, numbers, underscore, slash only | `validate:"alphanum_us_slash"` |
| `alphanum_us_dot` | `string` | letters, numbers, underscore, dot only | `validate:"alphanum_us_dot"` |

These rules are registered during package initialization. If that registration
fails, startup panics instead of leaving the validator partially configured.

### `phone_number`

| Rule | Details |
| --- | --- |
| accepted values | `^1[3-9]\d{9}$` |
| intended use | mobile phone input validation |
| common error shape | translated message meaning “format is invalid” |

### `dec_min` / `dec_max`

| Rule | Details |
| --- | --- |
| accepted field type | `decimal.Decimal` |
| parameter format | decimal string such as `10.5` |
| behavior | compares the field value against the parsed decimal threshold |
| failure modes | non-`decimal.Decimal` fields or invalid threshold params fail validation |

### `alphanum_us`

| Rule | Allowed characters |
| --- | --- |
| `alphanum_us` | letters, digits, `_` |
| `alphanum_us_slash` | letters, digits, `_`, `/` |
| `alphanum_us_dot` | letters, digits, `_`, `.` |

All three alphanum rules require at least one character; an empty string does
not match their regex.

Typical uses:

| Rule | Typical use |
| --- | --- |
| `alphanum_us` | action names, identifiers, simple codes |
| `alphanum_us_slash` | RPC resource names or slash-separated identifiers |
| `alphanum_us_dot` | file names, module names, dotted identifiers |

## Nullable / Optional Fields

VEF uses pointer types for nullable fields — there is no `null.*` wrapper package.

```go
type UserParams struct {
    // Required: validation runs against the value directly
    Username string `json:"username" validate:"required,min=3"`

    // Optional: validation only runs when the pointer is non-nil
    Phone *string `json:"phone" validate:"omitempty,phone_number"`

    // Optional with explicit "may be empty" semantics
    Bio *string `json:"bio" validate:"omitempty,max=500"`
}
```

The standard `omitempty` validator tag handles nullable cases — when the pointer is `nil`, subsequent rules are skipped; when it points to a value, the rest of the tag chain applies to that value. If you need to register a custom type that should participate in validation, use `validator.RegisterTypeFunc(...)` with a `validator.CustomTypeFunc` to teach the validator how to extract a comparable value from your type.

## Error Behavior

Validation returns the first translated validation error as a framework `result.Error`.

| Case | Outcome |
| --- | --- |
| validation succeeds | handler continues normally |
| one or more validation rules fail | framework returns a bad-request style error |
| non-validation error occurs during validation | framework wraps it as a bad-request style error |

HTTP behavior:

| Property | Value |
| --- | --- |
| business code | `result.ErrCodeBadRequest` |
| HTTP status | `400 Bad Request` |

Built-in go-playground rules and VEF custom rules read the active
`i18n.CurrentLanguage()` at validation time. Chinese (`zh-CN`) uses the Chinese
translator; other supported languages use the English translator.

## Registering Additional Custom Rules

Applications can register their own rules through `validator.RegisterValidationRules(...)`.

A rule is defined by `validator.ValidationRule`, which includes:

| Field | Purpose |
| --- | --- |
| `RuleTag` | validator tag name |
| `ErrMessageTemplate` | fallback translated message template |
| `ErrMessageI18nKey` | i18n message key |
| `Validate` | actual validation function |
| `ParseParam` | parameter extraction for error message placeholders |
| `CallValidationEvenIfNull` | whether the rule should run on null values |

For custom rule messages, `ErrMessageI18nKey` is checked first. When
`i18n.T(key)` returns a value different from the key, placeholders such as
`{0}` and `{1}` are replaced with `ParseParam` values. If no framework i18n
message is found, the go-playground translator renders `ErrMessageTemplate`.

## Practical Patterns

### Simple required fields

```go
type UserParams struct {
	api.P

	Username string `json:"username" validate:"required,alphanum,max=32" label:"Username"`
	Email    string `json:"email" validate:"omitempty,email,max=128" label:"Email"`
}
```

### Decimal range validation

```go
type PriceParams struct {
	api.P

	Amount decimal.Decimal `json:"amount" validate:"dec_min=0,dec_max=999999.99" label:"Amount"`
}
```

### Optional pointer validation

```go
type UserParams struct {
	api.P

	Phone *string `json:"phone" validate:"omitempty,phone_number" label:"Phone"`
}
```

## Custom Validation Rules

Struct-tag validation above (built-in tags plus the framework's own
`phone_number` / `dec_min` / `dec_max` / `alphanum_us*` rules) and custom rule
registration are two views of the same package, not two separate systems: a
`validate:"..."` tag selects which rule runs on a field, and
`validator.RegisterValidationRules` is how a new tag becomes available to
select. `validator.Validate` runs built-in and custom-registered rules
together through one call — there is no separate code path for "custom"
validation.

| API | Purpose |
| --- | --- |
| `validator.Validate(value)` | validates a value and returns the first framework validation error |
| `validator.RegisterValidationRules(rules...)` | adds custom `ValidationRule` entries |
| `validator.RegisterTypeFunc(fn, types...)` | registers custom type extraction for application-specific wrappers |
| `validator.CustomTypeFunc` | callback type accepted by `RegisterTypeFunc` |
| `validator.ValidationRule` | custom rule definition with tag, messages, validator callback, parameter parser, and null-call flag |

See [Registering Additional Custom Rules](#registering-additional-custom-rules)
above for the full `ValidationRule` field reference, message-resolution
order, and a worked example.

## Practical Advice

- put validation rules on typed params and meta structs, not inside handlers
- use `label` or `label_i18n` so error messages remain user-facing
- prefer framework custom rules when they match your contract
- keep custom rules narrow and domain-specific
- use pointer fields plus `omitempty` for nullable input; register a custom type extractor only for application-specific wrapper types

## Next Step

Read [Parameters And Metadata](./params-and-meta) to see where validation is triggered in request decoding.
