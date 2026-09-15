---
sidebar_position: 2
---

# Routing

VEF supports two routing strategies with one shared operation model:

- RPC through `POST /api`
- REST through `/api/<resource>`

RPC and REST are both explicit resource kinds. An RPC resource does not automatically generate REST routes, and a REST resource does not reuse the RPC single-endpoint transport.

## Routing Strategy Overview

| Strategy | Entry path | Operation identity source |
| --- | --- | --- |
| RPC | `POST /api` | request body fields `resource`, `action`, `version` |
| REST | `/api/<resource>` | resource name + action-defined HTTP method and sub-resource path |

## RPC Routing

RPC requests go to a single endpoint:

```text
POST /api
```

The internal router constant `DefaultRPCEndpoint` is `/api`. Form and multipart
RPC transports read JSON-encoded `params` and `meta` from the `FormKeyParams`
and `FormKeyMeta` form fields.

RPC request shape:

```json
{
  "resource": "sys/user",
  "action": "find_page",
  "version": "v1",
  "params": {
    "keyword": "tom"
  },
  "meta": {
    "page": 1,
    "size": 20
  }
}
```

This shape maps directly to `api.Request`.

### RPC naming rules

| Field | Rule | Examples |
| --- | --- | --- |
| `resource` | slash-separated lowercase resource path | `user`, `sys/user`, `approval/category` |
| `action` | `snake_case` | `find_page`, `get_user_info`, `resolve_challenge` |
| `version` | `v<number>` | `v1`, `v2` |

### RPC transport forms

RPC parsing supports:

| Content type | How request data is read |
| --- | --- |
| JSON | request body decoded directly into `api.Request` |
| form | `resource`, `action`, `version` from form fields, `params` and `meta` parsed from JSON strings |
| multipart form | same as form, plus uploaded files merged into `params` |

## REST Routing

REST routes are mounted under:

```text
/api/<resource>
```

At mount time, REST routes are built as `/api/<resource>/<subpath>` after the
action method is uppercased and the optional sub-resource path is prefixed with
`/`.

The HTTP method and optional sub-resource path come from the action string.

Examples:

| Resource | Action | Final route |
| --- | --- | --- |
| `users` | `get` | `GET /api/users` |
| `users` | `post` | `POST /api/users` |
| `users` | `get profile` | `GET /api/users/profile` |
| `users` | `put profile` | `PUT /api/users/profile` |
| `users` | `delete many` | `DELETE /api/users/many` |

### REST action parsing

REST action strings support:

| Pattern | Meaning |
| --- | --- |
| `<method>` | root route under the resource path |
| `<method> <sub-resource>` | extra kebab-case sub-resource path |

Parsing rules:

- the method token is uppercased when mounted
- the public resource validator accepts only lowercase HTTP verbs and optional
  kebab-case sub-resource paths; slash-separated sub-resources such as
  `admin/users` are allowed, but each segment must be kebab-case
- dynamic Fiber-style route params such as `/:id` are not accepted by
  `api.ValidateActionName` / `api.NewRESTResource`

### REST naming rules

| Field | Rule | Examples |
| --- | --- | --- |
| resource name | slash-separated lowercase path segments, kebab-case within segments when needed | `users`, `sys/user`, `user-profiles` |
| action method token | lowercase HTTP verb | `get`, `post`, `put`, `delete`, `patch` |
| action sub-resource | optional slash-separated kebab-case path | `profile`, `admin/users`, `user-friends` |

## Framework-Reserved Query Keys

The framework reserves some query keys for its own plumbing. The REST router
strips them from the request before the business parameter space is built, so
they never appear in `params`.

| Reserved key | Purpose |
| --- | --- |
| `__accessToken` | Carries the bearer token for requests that cannot set the `Authorization` header (for example, a browser `<img src>` or download link). The auth layer reads it from the Fiber context. |

Because these keys are consumed by the framework, they are not available to
business handlers and will not be reported as unmapped keys.

## How `params` Are Collected

### RPC

For RPC requests:

| Source | Lands in |
| --- | --- |
| request `params` object | `api.Request.Params` |
| multipart uploaded files | merged into `api.Request.Params` |

### REST

For REST requests, VEF merges multiple sources into `params`:

| Source | Lands in | Notes |
| --- | --- | --- |
| path params | `params` | extracted from Fiber route params |
| query string | `params` | always treated as params, not meta |
| JSON body on `POST` / `PUT` / `PATCH` | `params` | object body fields are merged into params |
| multipart form fields on `POST` / `PUT` / `PATCH` | `params` | text form fields go into params |
| multipart uploaded files on `POST` / `PUT` / `PATCH` | `params` | file arrays go into params |

## How `meta` Is Collected

### RPC

For RPC requests, `meta` comes directly from the request payload.

### REST

For REST requests, metadata is collected from headers using the `X-Meta-` prefix.

Example:

```http
X-Meta-page: 1
X-Meta-size: 20
X-Meta-format: excel
```

These values are stored in `api.Meta`.

The stored meta key is lowercased after removing the `X-Meta-` prefix. For
example, both `X-Meta-page` and `X-Meta-Page` become `meta["page"]`.

Important consequence:

- REST query strings are still `params`
- built-in typed helpers such as `page.Pageable` are still decoded from `meta`
- if a REST endpoint expects typed metadata, document the `X-Meta-*` headers explicitly

## Typed Request Decoding Implications

| Handler parameter | Decoded from |
| --- | --- |
| typed struct embedding `api.P` | `params` |
| typed struct embedding `api.M` | `meta` |
| `page.Pageable` | `meta` |
| `api.Params` | raw params |
| `api.Meta` | raw meta |

That means `?page=1&size=20` is not enough to populate a typed `page.Pageable` on REST endpoints unless you model paging as ordinary params instead of meta.

## Authentication Resolution Order

At operation runtime, authentication is resolved in this order:

1. `spec.Public == true` -> public endpoint
2. resource-level `Auth()` config when present
3. API engine default auth

Default engine auth is Bearer.

## Built-In Auth Strategies

VEF currently has these built-in auth strategy names:

| Strategy | Meaning |
| --- | --- |
| `none` | public endpoint |
| `bearer` | Bearer token authentication |
| `signature` | signature-based authentication |
| `ip` | source-IP whitelist authentication |
| `api_key` | static API key authentication |
| `http_basic` | RFC 7617 Basic authentication |

Helpers:

| Helper | Meaning |
| --- | --- |
| `api.Public()` | public operation |
| `api.BearerAuth()` | Bearer auth |
| `api.SignatureAuth()` | signature auth |
| `api.IPAuth(...)` | source-IP whitelist auth |
| `api.APIKeyAuth(...)` | API key auth; optional custom header name |
| `api.HTTPBasicAuth()` | HTTP Basic auth |

## Authentication Inputs

### Bearer

Bearer tokens are read from:

| Source | Format |
| --- | --- |
| `Authorization` header | `Bearer <token>` |
| query parameter | `__accessToken=<token>` |

### Signature

Signature auth reads:

| Header | Meaning |
| --- | --- |
| `X-App-ID` | external application ID |
| `X-Timestamp` | request timestamp |
| `X-Nonce` | replay-protection nonce |
| `X-Signature` | signature value |

### IP

IP auth resolves the client IP against a named whitelist loaded through
`security.IPWhitelistLoader` (default: config-backed, reading
`vef.security.ip_whitelists`). Behind a reverse proxy, configure
`vef.app.trusted_proxies` so Fiber resolves the real client IP. All failures
deny with `security.ErrIPNotAllowed`; an empty or missing whitelist is
fail-closed.

### API key

| Source | Format |
| --- | --- |
| request header (default `X-API-Key`, configurable per operation via `api.APIKeyAuth("X-My-Key")`) | the raw key value |

Keys resolve through `security.APIKeyLoader` (default backed by
`vef.security.api_keys`); missing or unmatched keys deny uniformly with
`security.ErrAPIKeyInvalid` (401).

### HTTP Basic

| Source | Format |
| --- | --- |
| `Authorization` header | `Basic base64(username:password)` |

Accounts resolve through `security.BasicAccountLoader` (default backed by
`vef.security.basic_accounts`) with constant-time comparison; malformed
headers, unknown accounts, and wrong passwords all deny uniformly with
`security.ErrBasicCredentialsInvalid` (401).

## Default Operation Behavior

Unless an operation overrides them, the API engine applies these defaults:

| Property | Default |
| --- | --- |
| version | `v1` |
| timeout | `30s` |
| auth strategy | Bearer |
| rate limit | `100` requests per `5 minutes` |

## Request Body Transport Encoding

For requests on the `/api` surface, a client can opt into transport encoding the
request body so that code-shaped payloads (integration adapter scripts, envelope
or auth scripts, `dry_run` bodies) survive middleboxes that false-positive on
them. The encoding is decoded back to raw JSON before the dispatcher parses it.

A client opts in by setting the header:

```http
X-Body-Encoding: base64
```

or:

```http
X-Body-Encoding: gzip+base64
```

Supported values:

| Value | Meaning |
| --- | --- |
| `base64` | Body is base64 of the raw JSON |
| `gzip+base64` | Body is base64 of a gzip stream of the raw JSON (decode order: base64, then gunzip) |

Native `Content-Encoding` values such as `gzip`, `br`, `deflate`, and `zstd` are
decompressed by Fiber itself; this middleware covers only the base64 forms Fiber
does not know.

### Protected body transport

Set `vef.api.body_encoding.enabled = true` to require authenticated,
non-plaintext JSON bodies on the `/api` surface:

```toml
[vef.api.body_encoding]
enabled = true
encoding = "aes-gcm+base64" # or "sm4-gcm+base64"
key = "<standard-base64-key>"
```

Both encodings use GCM and transport standard padded Base64 as the raw body,
with no JSON wrapper. The wire bytes are
`Base64(12-byte random nonce || ciphertext || 16-byte tag)`. Every request and
response gets a fresh nonce.

When protected transport is enabled:

- JSON requests must set `X-Body-Encoding` to the configured value and send the
  encoded text as the raw request body — not as a JSON string.
- JSON responses, including transport and downstream error envelopes, use the
  same encoding and return the marker in `X-Body-Encoding`. CORS exposes the
  header so browsers can read it.
- Plaintext JSON and weaker `base64` / `gzip+base64` requests are rejected;
  they cannot downgrade protected mode.
- Multipart and binary bodies keep their native formats, and non-`/api`
  responses remain untouched.
- The decoded request body is still bounded by `vef.app.body_limit`.
- HTTP compression is skipped for protected `/api` responses because the
  ciphertext is already high entropy.

AES accepts a 16-, 24-, or 32-byte key; SM4 accepts a 16-byte key. The
standard-base64 key must match the client. This transport does not replace
HTTPS: a symmetric key held by browser code is extractable. API `signature`
authentication is orthogonal because it signs `appID+method+path+timestamp+nonce`,
not the request body.

Decoding happens before the content-type guard and dispatcher run, on `/api`
paths only. Storage, content-hash caches, and audit all see the raw body. The
body-limit guard (`vef.app.body_limit`) is applied to the decoded size, so a
decompression bomb cannot outgrow it.

Failure modes:

| Case | Result |
| --- | --- |
| JSON request without the marker while protected transport is enabled | `api.ErrBodyEncodingRequired` (HTTP 400) |
| Unsupported `X-Body-Encoding` value | `api.ErrUnsupportedBodyEncoding` (HTTP 400) |
| Malformed base64 or corrupt gzip stream | `api.ErrBodyDecodeFailed` (HTTP 400) |
| Decoded body larger than the configured limit | `api.ErrBodyTooLarge` (HTTP 413) |

## Response Shape

Handlers normally return responses through `result.Ok(...)` or `result.Err(...)`, so both RPC and REST share the same response structure:

```json
{
  "code": 0,
  "message": "Success",
  "data": {}
}
```

The exact message text is language-dependent. With the framework default language you will usually see `成功`; with English selected you will usually see `Success`.

## Practical Advice

- use RPC when your API is action-oriented
- use REST when HTTP method semantics and path structure matter
- document whether paging and sorting are expected in `params` or `meta`
- keep request semantics explicit; do not try to hide RPC and REST differences from yourself
- think in terms of resources and operations, not only endpoints

## Next Step

Read [Parameters And Metadata](./params-and-meta) for the exact decoding behavior used by handler injection.
