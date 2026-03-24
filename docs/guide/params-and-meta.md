---
sidebar_position: 5
---

# Parameters And Metadata

VEF separates request input into two sections:

- `params`: business input
- `meta`: request-level control data

That split exists for RPC requests and is preserved internally for REST requests as well.

## Request Model Overview

| Section | Purpose | Typical content |
| --- | --- | --- |
| `params` | business payload | search fields, write payloads, uploaded files, command inputs |
| `meta` | request controls | paging, sorting, export format, option column mapping |

## Supported Typed Targets

The framework supports these request-decoding targets:

| Target type | Decoded from | Validation | Typical use |
| --- | --- | --- | --- |
| typed struct embedding `api.P` | `params` | Yes | business params |
| typed struct embedding `api.M` | `meta` | Yes | typed meta |
| `page.Pageable` | `meta` | Yes | paging |
| `api.Params` | `params` | No typed validation | raw dynamic payload |
| `api.Meta` | `meta` | No typed validation | raw dynamic meta |

## `api.P` Marks Params Structs

Embed `api.P` in structs that should decode from `Request.Params`:

```go
type CreateUserParams struct {
	api.P

	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}
```

When a handler accepts `CreateUserParams` or `*CreateUserParams`, the framework:

1. decodes `params`
2. validates the struct
3. injects the typed value

## `api.M` Marks Meta Structs

Embed `api.M` in structs that should decode from `Request.Meta`:

```go
type PageMeta struct {
	api.M
	page.Pageable
}
```

This is how typed request controls are injected.

## Built-In Meta Helpers

The framework has built-in support for these meta-oriented helper types:

| Type | Meaning | Notes |
| --- | --- | --- |
| `page.Pageable` | page number and page size | directly recognized as a built-in meta type |
| `crud.Sortable` | sort specs | usually embedded inside a typed meta struct |

Important distinction:

- `page.Pageable` is a built-in meta target type
- `crud.Sortable` is not resolved as a standalone built-in meta type, but it works naturally when embedded in a typed `api.M` struct

## Raw Access

If you do not want typed decoding, handlers can accept:

| Type | Meaning |
| --- | --- |
| `api.Params` | raw params map |
| `api.Meta` | raw meta map |

Use raw access for dynamic, proxy-style, or partially unknown payloads. Prefer typed structs for stable business APIs.

## RPC Decoding Rules

For RPC requests, decoding depends on the transport content type:

| RPC request type | `params` source | `meta` source | Notes |
| --- | --- | --- | --- |
| JSON body | request JSON `params` object | request JSON `meta` object | standard RPC shape |
| form request | form field `params` parsed as JSON string | form field `meta` parsed as JSON string | used for form-style clients |
| multipart form | form field `params` parsed as JSON string, plus uploaded files merged into params | form field `meta` parsed as JSON string | file fields are added into params |

## REST Decoding Rules

For REST requests:

| Input source | Lands in | Notes |
| --- | --- | --- |
| query string | `params` | used for read filters and plain request fields |
| JSON body on `POST` / `PUT` / `PATCH` | `params` | write payload |
| multipart fields on `POST` / `PUT` / `PATCH` | `params` | includes uploaded files |
| `X-Meta-*` headers | `meta` | request-level control values |

That means paging and sorting are not automatically pulled from query string into built-in meta helpers. If a handler expects meta-based controls such as `page.Pageable`, the caller should provide them through `X-Meta-*` headers or a typed meta contract.

## Multipart File Support

Multipart uploads can populate params fields such as:

| Shape | Notes |
| --- | --- |
| `*multipart.FileHeader` | standard single-file upload field |
| raw file entries inside `api.Params` | useful for proxy-style or dynamic handlers |

This is how built-in storage and import endpoints receive uploaded files.

## Validation Behavior

Typed params and typed meta values are automatically validated after decoding.

| Target type | Validation |
| --- | --- |
| typed `api.P` struct | yes |
| typed `api.M` struct | yes |
| `page.Pageable` | yes |
| `api.Params` | no typed validation |
| `api.Meta` | no typed validation |

Validation uses `validator.Validate(...)` after decoding. If validation fails, the framework returns a bad-request style result with translated field messages.

## Practical Patterns

### Standard search request

```go
type UserSearch struct {
	api.P
	Keyword string `json:"keyword" search:"contains,column=username|email"`
}

type UserMeta struct {
	api.M
	page.Pageable
	crud.Sortable
}
```

### Dynamic proxy-style request

```go
func (*ProxyResource) Forward(params api.Params, meta api.Meta) error {
	// handle raw data
	return nil
}
```

## Practical Advice

- put business fields in `params`
- put paging, sorting, export mode, and similar request controls in `meta`
- prefer typed structs over raw maps for long-term maintainability
- embed `api.P` and `api.M` explicitly so decoding intent stays obvious
- use raw `api.Params` / `api.Meta` only when the request contract is truly dynamic

## Next Step

Read [Custom Handlers](./custom-handlers) to see how these decoded values are injected into handler signatures.
