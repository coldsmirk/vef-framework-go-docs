---
sidebar_position: 8
---

# Validation

VEF uses `go-playground/validator` as the base validation engine and layers framework-specific behavior on top:

- translated error messages
- `label` / `label_i18n` support
- custom validation rules
- null-type support for framework null wrappers

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
| `label_i18n:"..."` | resolves the label through i18n first, then falls back to the field name |
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

### `phone_number`

| Rule | Details |
| --- | --- |
| accepted values | `1[3-9]\\d{9}` |
| intended use | mobile phone input validation |
| common error shape | translated message meaning “format is invalid” |

### `dec_min` / `dec_max`

| Rule | Details |
| --- | --- |
| accepted field type | `decimal.Decimal` |
| parameter format | decimal string such as `10.5` |
| behavior | compares the field value against the parsed decimal threshold |

### `alphanum_us`

| Rule | Allowed characters |
| --- | --- |
| `alphanum_us` | letters, digits, `_` |
| `alphanum_us_slash` | letters, digits, `_`, `/` |
| `alphanum_us_dot` | letters, digits, `_`, `.` |

Typical uses:

| Rule | Typical use |
| --- | --- |
| `alphanum_us` | action names, identifiers, simple codes |
| `alphanum_us_slash` | RPC resource names or slash-separated identifiers |
| `alphanum_us_dot` | file names, module names, dotted identifiers |

## Supported Null Types

VEF registers custom type functions so that these null wrappers participate correctly in validation:

| Supported null type |
| --- |
| `null.String` |
| `null.Int` |
| `null.Int16` |
| `null.Int32` |
| `null.Float` |
| `null.Bool` |
| `null.Byte` |
| `null.DateTime` |
| `null.Date` |
| `null.Time` |
| `null.Decimal` |

Practical effect:

- when a null wrapper is valid, validation runs against the wrapped value
- when it is invalid, the value is treated as `nil`

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

### Null wrapper validation

```go
type UserParams struct {
	api.P

	Phone null.String `json:"phone" validate:"omitempty,phone_number" label:"Phone"`
}
```

## Practical Advice

- put validation rules on typed params and meta structs, not inside handlers
- use `label` or `label_i18n` so error messages remain user-facing
- prefer framework custom rules when they match your contract
- keep custom rules narrow and domain-specific
- if the field type is a framework null wrapper, rely on built-in null-type support instead of manual unwrapping

## Next Step

Read [Parameters And Metadata](./params-and-meta) to see where validation is triggered in request decoding.
