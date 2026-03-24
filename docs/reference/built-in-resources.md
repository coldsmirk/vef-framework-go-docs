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
| `sys/storage` | `storage` | Bearer auth by default | File upload, presigned URL generation, temporary object cleanup, object metadata, object listing |
| `sys/schema` | `schema` | Bearer auth by default | Database schema inspection |
| `sys/monitor` | `monitor` | Bearer auth by default | Runtime and host monitoring data |

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

Storage resource provided by the storage module.

### Operations

| Action | Access | Rate limit | Purpose | Params |
| --- | --- | --- | --- | --- |
| `upload` | Bearer auth required | API engine default | Uploads one file into storage and returns object metadata | `UploadParams` via multipart form |
| `get_presigned_url` | Bearer auth required | API engine default | Generates a temporary presigned URL for object access | `GetPresignedURLParams` |
| `delete_temp` | Bearer auth required | API engine default | Deletes an object only when its key is under the `temp/` prefix | `DeleteTempParams` |
| `stat` | Bearer auth required | API engine default | Returns metadata for one object | `StatParams` |
| `list` | Bearer auth required | API engine default | Lists objects under a prefix | `ListParams` |

Only the actions above are registered as built-in RPC operations. The underlying service supports more capabilities such as copy and move, but they are not exposed here by default.

Related HTTP route:

- `/storage/files/<key>` is an app-level download proxy route, not an RPC action
- it does not automatically inherit RPC Bearer authentication

### `upload` parameters

This action expects `multipart/form-data`, not a JSON RPC body. Multipart fields are decoded into `params`.

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `file` | `file` | Yes | Uploaded file content. This is required and maps to `*multipart.FileHeader` |
| `contentType` | `string` | No | Explicit content type override. If omitted, the server uses the uploaded file header content type |
| `metadata` | `object<string, string>` | No | Optional storage metadata map passed to the storage service |

Notes:

- object keys are generated server-side under a date-based `temp/YYYY/MM/DD/...` path
- the original filename is automatically stored in metadata
- JSON requests are rejected for this action

### `get_presigned_url` parameters

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `key` | `string` | Yes | Object key to access |
| `expires` | `int` | No | Expiration time in seconds. Defaults to `3600` |
| `method` | `string` | No | HTTP method used when signing the URL. Defaults to `GET` |

### `delete_temp` parameters

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `key` | `string` | Yes | Object key to delete. Must start with `temp/` |

### `stat` parameters

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `key` | `string` | Yes | Object key whose metadata should be returned |

### `list` parameters

| Parameter | Type | Required | Description |
| --- | --- | --- | --- |
| `prefix` | `string` | No | Prefix filter used for object listing |
| `recursive` | `bool` | No | Whether to list recursively instead of only the current prefix level |
| `maxKeys` | `int` | No | Maximum number of objects to return |

Minimal request example:

```json
{
  "resource": "sys/storage",
  "action": "list",
  "version": "v1",
  "params": {
    "prefix": "temp/",
    "recursive": false
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
| `list_triggers` | Bearer auth required | Custom operation max `60` | Returns all triggers in the current database or schema | None |

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

Notes:

- these actions do not accept framework-defined input parameters
- some actions may return a monitor-not-ready error when the underlying data source is unavailable

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

Those resources are intentionally not expanded in this page because they are domain-level workflow resources, not the framework's core general-purpose built-ins.

## See also

- [Authentication](../security/authentication) for the behavior of `security/auth`
- [File Storage](../features/storage) for `sys/storage`
- [Schema](../features/schema) for `sys/schema`
- [Monitor](../features/monitor) for `sys/monitor`
