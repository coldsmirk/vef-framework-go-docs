---
sidebar_position: 4
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
| REST | `/api/<resource>` | resource name + action-defined HTTP method and sub-path |

## RPC Routing

RPC requests go to a single endpoint:

```text
POST /api
```

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

The HTTP method and optional sub-path come from the action string.

Examples:

| Resource | Action | Final route |
| --- | --- | --- |
| `users` | `get` | `GET /api/users` |
| `users` | `post` | `POST /api/users` |
| `users` | `get profile` | `GET /api/users/profile` |
| `users` | `put /:id` | `PUT /api/users/:id` |
| `users` | `delete /many` | `DELETE /api/users/many` |

### REST action parsing

REST action strings support:

| Pattern | Meaning |
| --- | --- |
| `<method>` | root route under the resource path |
| `<method> <sub-path>` | extra path segment or Fiber-style parameter path |

Parsing rules:

- the method token is uppercased when mounted
- if the sub-path does not start with `/`, the router adds it automatically
- Fiber-style params such as `/:id` are preserved

### REST naming rules

| Field | Rule | Examples |
| --- | --- | --- |
| resource name | slash-separated lowercase path segments, kebab-case within segments when needed | `users`, `sys/user`, `user-profiles` |
| action method token | lowercase HTTP verb | `get`, `post`, `put`, `delete`, `patch` |
| action sub-path | kebab-case or explicit route pattern | `profile`, `admin`, `/:id`, `/tree/options` |

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
| multipart form fields | `params` | text form fields go into params |
| multipart uploaded files | `params` | file arrays go into params |

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

Helpers:

| Helper | Meaning |
| --- | --- |
| `api.Public()` | public operation |
| `api.BearerAuth()` | Bearer auth |
| `api.SignatureAuth()` | signature auth |

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

## Default Operation Behavior

Unless an operation overrides them, the API engine applies these defaults:

| Property | Default |
| --- | --- |
| version | `v1` |
| timeout | `30s` |
| auth strategy | Bearer |
| rate limit | `100` requests per `5 minutes` |

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
