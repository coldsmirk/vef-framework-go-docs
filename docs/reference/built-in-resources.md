---
sidebar_position: 2
---

# Built-in Resources

VEF registers several RPC resources for you when the corresponding modules are enabled in the default boot chain.

Unless noted otherwise:

- resources in this page are RPC resources mounted under `/api`
- operations use the standard RPC request envelope: `resource`, `action`, `version`, `params`, and `meta`
- non-public operations inherit the API engine's default Bearer authentication
- operations without a custom rate limit inherit the API engine default rate limit
  The stock engine default is `100` requests per `5 minutes`, but applications may override it

## Resource Overview

| Resource | Module | Default access model | Notes |
| --- | --- | --- | --- |
| `security/auth` | `security` | Mixed: some actions are public, some require Bearer auth | Login flow, token refresh, logout, challenge resolution, current-user info |
| `sys/storage` | `storage` | Bearer auth by default | Multipart upload session lifecycle (init / part / list / complete / abort). Downloads are served via the `/storage/files/<key>` app proxy, not via RPC. |
| `sys/schema` | `schema` | Bearer auth by default | Database schema inspection |
| `sys/monitor` | `monitor` | Bearer auth by default | Runtime and host monitoring data |
| `approval/*` | `approval` | Bearer auth required; per-action permissions where declared | Optional workflow resources registered only when `vef.ApprovalModule` is enabled |

## `security/auth`

Authentication resource provided by the security module.

### Operations

| Action | Access | Rate limit | Purpose | Params |
| --- | --- | --- | --- | --- |
| `login` | Public | `max = vef.security.login_rate_limit` (module default `6`) | Authenticates a user or external app and returns either tokens or the first pending login challenge | `LoginParams` |
| `refresh` | Public | `max = vef.security.refresh_rate_limit` (module default `1`) | Exchanges a valid refresh token for a fresh token pair | `RefreshParams` |
| `logout` | Bearer auth required | API engine default | Returns success immediately; token invalidation is expected to happen on the client side | None |
| `resolve_challenge` | Public | `max = vef.security.login_rate_limit` (module default `6`) | Resolves the current login challenge and returns either the next challenge or final tokens | `ResolveChallengeParams` |
| `get_user_info` | Bearer auth required | API engine default | Loads current-user profile, menus, permission tokens, and other session data through `security.UserInfoLoader` | Raw `params` map, application-defined |

### `login` parameters

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `type` | `string` | Yes | Login type. The built-in login flow currently supports `password` only |
| `principal` | `string` | Yes | Login identifier, typically the username |
| `credentials` | `string` | Yes | Login credential. For `type = "password"`, this is the plaintext password |

Minimal request example:

```json
{
  "resource": "security/auth",
  "action": "login",
  "version": "v1",
  "params": {
    "type": "password",
    "principal": "alice",
    "credentials": "secret"
  }
}
```

### `refresh` parameters

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `refreshToken` | `string` | Yes | Refresh token that will be validated and exchanged for a new token pair |

### `resolve_challenge` parameters

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `challengeToken` | `string` | Yes | Challenge-state token returned by a previous `login` or `resolve_challenge` call |
| `type` | `string` | Yes | Challenge type currently being resolved, such as `totp` or another provider-specific challenge identifier |
| `response` | `any` | Yes | Challenge response payload consumed by the matching `security.ChallengeProvider` |

### `get_user_info` parameters

This action does not define a typed params struct. Any `params` object is forwarded to `security.UserInfoLoader.LoadUserInfo(...)`.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| Framework-defined parameters | None | No | The framework does not reserve fixed keys here |
| Application-defined parameters | `object` | No | Optional extension data interpreted by your own `security.UserInfoLoader` implementation |

Notes:

- if no `security.UserInfoLoader` is registered, this action returns `not implemented`
- response shape is defined by `security.UserInfo`

## `sys/storage`

Storage resource provided by the storage module. The single-PUT `upload` action was retired in v0.21; every upload now goes through the multipart session lifecycle below. See [File Storage](../infrastructure/storage) for the surrounding lifecycle (claim, pending-delete, ACL).

### Operations

| Action | Access | Rate limit | Purpose | Params |
| --- | --- | --- | --- | --- |
| `init_upload` | Bearer auth required | API engine default | Open a new multipart session. Server returns the negotiated part plan and an opaque `claimId`. | `InitUploadParams` |
| `upload_part` | Bearer auth required | API engine default | Upload one part of an open session (multipart form). | `UploadPartParams` |
| `list_parts` | Bearer auth required | API engine default | List parts already uploaded for a session. | `ListPartsParams` |
| `complete_upload` | Bearer auth required | API engine default | Seal a session; the server assembles the part manifest from recorded parts. | `CompleteUploadParams` |
| `abort_upload` | Bearer auth required | API engine default | Abort and release a session. | `AbortUploadParams` |

Related HTTP route:

- `/storage/files/<key>` is an app-level download proxy route, not an RPC action.
- It does not automatically inherit RPC Bearer authentication; `pub/*` is served anonymously and all other keys are governed by `storage.FileACL`.

### `init_upload` parameters

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `filename` | `string` | Yes | Original filename (≤ 255 chars). Used to derive the safe extension and stored on the upload claim. |
| `size` | `int` | Yes | Total object size in bytes (≥ 1). The server validates against `vef.storage.max_upload_size`. |
| `contentType` | `string` | No | Client-supplied MIME (≤ 127 chars). Sanitized server-side — unsafe values are overridden by extension-based detection or fall back to `application/octet-stream`. |
| `public` | `bool` | No | Place the key under `pub/` instead of `priv/`. Requires `vef.storage.allow_public_uploads = true`. |

Requests with `public = true` are rejected unless `vef.storage.allow_public_uploads`
is enabled.

Only `claimId` is client-facing; the backend session handle is internal and is never returned.

Response:

| Field | Type | Description |
| --- | --- | --- |
| `key` | `string` | Planned final object key under `priv/` or `pub/`. |
| `claimId` | `string` | Opaque client-facing session handle for the remaining upload actions. |
| `originalFilename` | `string` | Client-supplied filename stored on the upload claim. |
| `partSize` | `int` | Backend-authoritative part size in bytes. |
| `partCount` | `int` | Number of parts the client must upload. Small files still use `partCount = 1`. |
| `expiresAt` | timestamp | Claim expiration time. |

### `upload_part` parameters

This action expects `multipart/form-data`, not JSON. The form carries the normal RPC fields (`resource`, `action`, `version`), a `params` field containing JSON such as `{"claimId":"...","partNumber":1}`, and a file part named `file`.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `file` | file | Yes | Raw part bytes. |
| `claimId` | `string` | Yes | The `claimId` returned by `init_upload`. |
| `partNumber` | `int` | Yes | 1-indexed part position. Must be `≤ partCount` and the size must equal the server's `partSize` (the final part may be smaller). |

The backend ETag is intentionally not returned to the client — the server records it server-side and reuses it during `complete_upload`.

Response:

| Field | Type | Description |
| --- | --- | --- |
| `partNumber` | `int` | Accepted part number. |
| `size` | `int` | Recorded byte size for that part. |

### `list_parts` parameters

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `claimId` | `string` | Yes | Active pending session to inspect. |

Response:

| Field | Type | Description |
| --- | --- | --- |
| `parts` | `object[]` | Uploaded parts ordered by `partNumber`; each entry contains `partNumber` and `size`. Part ETags are not exposed. |

### `complete_upload` parameters

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `claimId` | `string` | Yes | Session to seal. The server reassembles the manifest from its own part-store records — no client-supplied ETags are accepted. |

On success the server marks the existing claim as uploaded, clears its recorded parts, and returns object metadata plus `originalFilename`. The uploaded claim is still pending business adoption. Subsequent calls against the same uploaded claim are idempotent fast-paths.
If the assembled object size does not match the claim, the action returns
`ErrCodeUploadSizeMismatch`.

Response:

| Field | Type | Description |
| --- | --- | --- |
| `bucket` | `string` | Backend bucket name when the provider reports one. |
| `key` | `string` | Final object key. |
| `eTag` | `string` | Final object ETag. This is not a part ETag and is not supplied to `complete_upload`. |
| `size` | `int` | Final object size in bytes. |
| `contentType` | `string` | Sanitized content type stored for the object. |
| `lastModified` | timestamp | Backend last-modified time. |
| `metadata` | `object` | Optional backend metadata map. The HTTP upload API does not accept user-supplied metadata. |
| `originalFilename` | `string` | Filename captured during `init_upload`. |

### `abort_upload` parameters

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `claimId` | `string` | Yes | Session to abort. |

Response has no data payload. The action is retry-safe: a missing claim returns success, and a claim already marked non-pending for the same owner is a no-op. An existing claim owned by a different principal is still rejected.

Minimal request example:

```json
{
  "resource": "sys/storage",
  "action": "init_upload",
  "version": "v1",
  "params": {
    "filename": "report.pdf",
    "size": 25600000,
    "contentType": "application/pdf",
    "public": false
  }
}
```

## `sys/schema`

Schema inspection resource provided by the schema module.

### Operations

| Action | Access | Rate limit | Purpose | Params |
| --- | --- | --- | --- | --- |
| `list_tables` | Bearer auth required | Custom operation max `60` | Returns all tables in the current database or schema | None |
| `get_table_schema` | Bearer auth required | Custom operation max `60` | Returns detailed schema information for one table | `GetTableSchemaParams` |
| `list_views` | Bearer auth required | Custom operation max `60` | Returns all views in the current database or schema | None |

### `get_table_schema` parameters

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `name` | `string` | Yes | Table name to inspect |

## `sys/monitor`

Monitoring resource provided by the monitor module.

### Operations

| Action | Access | Rate limit | Purpose | Params |
| --- | --- | --- | --- | --- |
| `get_overview` | Bearer auth required | Custom operation max `60` | Returns a combined system overview snapshot | None |
| `get_cpu` | Bearer auth required | Custom operation max `60` | Returns CPU information and usage data | None |
| `get_memory` | Bearer auth required | Custom operation max `60` | Returns memory usage information | None |
| `get_disk` | Bearer auth required | Custom operation max `60` | Returns disk and partition information | None |
| `get_network` | Bearer auth required | Custom operation max `60` | Returns network interface and I/O statistics | None |
| `get_host` | Bearer auth required | Custom operation max `60` | Returns static host information | None |
| `get_process` | Bearer auth required | Custom operation max `60` | Returns information about the current application process | None |
| `get_load` | Bearer auth required | Custom operation max `60` | Returns system load averages | None |
| `get_build_info` | Bearer auth required | Custom operation max `60` | Returns application build metadata | None |
| `get_event_streams` | Bearer auth required | Custom operation max `60` | Reports every redis_stream stream and consumer group (consumers / pending / lag / last-delivered) via the optional `event.StreamInspector`, so operators can spot orphaned groups | None |

Notes:

- these actions do not accept framework-defined input parameters
- some actions may return a monitor-not-ready error when the underlying data source is unavailable
- `get_event_streams` returns an empty, disabled report when no `event.StreamInspector` is available (for example when the redis_stream transport is off)

Minimal request example:

```json
{
  "resource": "sys/monitor",
  "action": "get_overview",
  "version": "v1"
}
```

## Approval resources

If you explicitly include the approval module, the framework also registers additional `approval/*` resources.

The registered resources are `approval/category`, `approval/delegation`, `approval/flow`, `approval/instance`, `approval/my`, and `approval/admin`.

They are expanded in [Approval Module](../approval), including each action name, required permission, params type, tenancy rule, audit setting, and rate limit. This page keeps them as an index because they are domain-level workflow resources rather than the framework's core general-purpose built-ins.

## See also

- [Authentication](../security/authentication) for the behavior of `security/auth`
- [File Storage](../infrastructure/storage) for `sys/storage`
- [Schema](../infrastructure/schema) for `sys/schema`
- [Monitor](../infrastructure/monitor) for `sys/monitor`
