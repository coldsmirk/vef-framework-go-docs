---
sidebar_position: 1
---

# 内置资源

当对应模块处于默认启动链中时，VEF 会自动注册一批 RPC 资源。

除非特别说明，本页默认约定如下：

- 这些资源都是挂在 `/api` 下的 RPC 资源
- 请求仍然使用标准 RPC 包装格式：`resource`、`action`、`version`、`params`、`meta`
- 没有标记为公开接口的 action，默认继承 API 引擎的 Bearer 认证
- 没有单独配置限流的 action，默认继承 API 引擎限流
  框架默认值是 `100` 次请求 / `5` 分钟，但应用也可以覆盖这个默认值

## 资源总览

| 资源 | 来源模块 | 默认访问模型 | 说明 |
| --- | --- | --- | --- |
| `security/auth` | `security` | 混合：部分 action 公开，部分需要 Bearer 认证 | 登录、刷新令牌、登出、挑战校验、当前用户信息 |
| `sys/storage` | `storage` | 默认 Bearer 认证 | 多片上传 session 生命周期（init / part / list / complete / abort）。下载走 `/storage/files/<key>` app 代理，不通过 RPC。 |
| `sys/schema` | `schema` | 默认 Bearer 认证 | 数据库结构检查 |
| `sys/monitor` | `monitor` | 默认 Bearer 认证 | 运行时与宿主机监控信息 |
| `approval/*` | `approval` | 需要 Bearer 认证；声明权限的 action 还会校验对应权限点 | 仅在启用 `vef.ApprovalModule` 时注册的可选工作流资源 |

## `security/auth`

由 security 模块提供的认证资源。

### 操作列表

| Action | 访问方式 | 限流 | 作用 | 参数 |
| --- | --- | --- | --- | --- |
| `login` | 公开接口 | `max = vef.security.login_rate_limit`（模块默认值 `6`） | 执行登录，返回最终 token，或者返回当前待处理的登录挑战 | `LoginParams` |
| `refresh` | 公开接口 | `max = vef.security.refresh_rate_limit`（模块默认值 `1`） | 用合法的 refresh token 换取新的 token 对 | `RefreshParams` |
| `logout` | 需要 Bearer 认证 | 继承 API 引擎默认限流 | 立即返回成功，实际 token 失效通常由客户端清理本地凭证实现 | 无 |
| `resolve_challenge` | 公开接口 | `max = vef.security.login_rate_limit`（模块默认值 `6`） | 校验当前登录挑战，返回下一步挑战或最终 token | `ResolveChallengeParams` |
| `get_user_info` | 需要 Bearer 认证 | 继承 API 引擎默认限流 | 通过 `security.UserInfoLoader` 加载当前用户资料、菜单、权限点等信息 | 原始 `params` map，由应用自行定义 |

### `login` 参数

| 参数名 | 类型 | 必填 | 含义 |
| --- | --- | --- | --- |
| `type` | `string` | 是 | 登录方式/认证类型。目前仅支持 `password`，即账号密码登录 |
| `principal` | `string` | 是 | 登录标识，通常就是用户名 |
| `credentials` | `string` | 是 | 登录凭证。在 `type = "password"` 时就是明文密码 |

最小请求示例：

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

### `refresh` 参数

| 参数名 | 类型 | 必填 | 含义 |
| --- | --- | --- | --- |
| `refreshToken` | `string` | 是 | 用于换发新 token 的 refresh token |

### `resolve_challenge` 参数

| 参数名 | 类型 | 必填 | 含义 |
| --- | --- | --- | --- |
| `challengeToken` | `string` | 是 | 前一步 `login` 或 `resolve_challenge` 返回的 challenge 状态 token |
| `type` | `string` | 是 | 当前要处理的挑战类型，例如 `totp` 或其他 provider 自定义类型 |
| `response` | `any` | 是 | 对应挑战 provider 消费的响应载荷 |

### `get_user_info` 参数

这个 action 没有定义固定的 typed params 结构。任何 `params` 对象都会原样传给 `security.UserInfoLoader.LoadUserInfo(...)`。

| 参数名 | 类型 | 必填 | 含义 |
| --- | --- | --- | --- |
| 框架固定参数 | 无 | 否 | 框架本身不在这里保留固定字段 |
| 应用自定义参数 | `object` | 否 | 由你自己的 `security.UserInfoLoader` 实现解释的可选扩展参数 |

补充说明：

- 如果没有注册 `security.UserInfoLoader`，这个 action 会返回 `not implemented`
- 返回体结构由 `security.UserInfo` 决定

## `sys/storage`

由 storage 模块提供的存储资源。v0.21 起单次 PUT 形式的 `upload` 动作已废弃，所有上传都走下面的多片 session 协议。围绕的生命周期（claim、pending-delete、ACL）见 [文件存储](../infrastructure/storage)。

### 操作列表

| Action | 访问方式 | 限流 | 作用 | 参数 |
| --- | --- | --- | --- | --- |
| `init_upload` | 需要 Bearer 认证 | 继承 API 引擎默认限流 | 开启新的多片 session，服务端返回协商好的 part 计划和不透明 `claimId` | `InitUploadParams` |
| `upload_part` | 需要 Bearer 认证 | 继承 API 引擎默认限流 | 上传一片（multipart form） | `UploadPartParams` |
| `list_parts` | 需要 Bearer 认证 | 继承 API 引擎默认限流 | 列出当前 session 已上传的 parts | `ListPartsParams` |
| `complete_upload` | 需要 Bearer 认证 | 继承 API 引擎默认限流 | 完成 session；服务端从已记录的 parts 组装清单 | `CompleteUploadParams` |
| `abort_upload` | 需要 Bearer 认证 | 继承 API 引擎默认限流 | 中止并释放 session | `AbortUploadParams` |

相关 HTTP 路由：

- `/storage/files/<key>` 是 app-level 下载代理路由，不是 RPC action。
- 它不会自动继承 RPC 层 Bearer 认证；`pub/*` 匿名可读，其他 key 由 `storage.FileACL` 决定。

### `init_upload` 参数

| 参数名 | 类型 | 必填 | 含义 |
| --- | --- | --- | --- |
| `filename` | `string` | 是 | 原始文件名（≤ 255 字符），用于推断安全扩展名并写入 upload claim。 |
| `size` | `int` | 是 | 对象总字节数（≥ 1），服务端按 `vef.storage.max_upload_size` 校验上限。 |
| `contentType` | `string` | 否 | 客户端 MIME（≤ 127 字符）。服务端会做白名单清洗——不安全的值会被扩展名探测结果或 `application/octet-stream` 覆盖。 |
| `public` | `bool` | 否 | 把 key 放到 `pub/` 而不是 `priv/`，需要 `vef.storage.allow_public_uploads = true`。 |

请求里 `public = true` 时，必须先启用 `vef.storage.allow_public_uploads`。

只有 `claimId` 是客户端可见 session handle；后端 session handle 是内部值，不会返回。

响应：

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `key` | `string` | 计划中的最终 object key，位于 `priv/` 或 `pub/` 下。 |
| `claimId` | `string` | 后续上传 action 使用的不透明客户端 session handle。 |
| `originalFilename` | `string` | 写入 upload claim 的客户端原始文件名。 |
| `partSize` | `int` | 后端权威 part 字节大小。 |
| `partCount` | `int` | 客户端需要上传的 part 数；小文件也会使用 `partCount = 1`。 |
| `expiresAt` | timestamp | claim 过期时间。 |

### `upload_part` 参数

此 action 要求使用 `multipart/form-data`，不能用 JSON。form 中带标准 RPC 字段（`resource`、`action`、`version`），一个内容为 JSON 的 `params` 字段（例如 `{"claimId":"...","partNumber":1}`），以及名为 `file` 的文件 part。

| 参数名 | 类型 | 必填 | 含义 |
| --- | --- | --- | --- |
| `file` | file | 是 | part 原始字节。 |
| `claimId` | `string` | 是 | `init_upload` 返回的 `claimId`。 |
| `partNumber` | `int` | 是 | 1 起的 part 序号，必须 `≤ partCount`；除最后一片外尺寸必须等于服务端 `partSize`。 |

后端 ETag 不会返回给客户端——服务端自己持久化记录，并在 `complete_upload` 时使用。

响应：

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `partNumber` | `int` | 已接收的 part 序号。 |
| `size` | `int` | 服务端记录的该 part 字节数。 |

### `list_parts` 参数

| 参数名 | 类型 | 必填 | 含义 |
| --- | --- | --- | --- |
| `claimId` | `string` | 是 | 要查询的 active pending session。 |

响应：

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `parts` | `object[]` | 按 `partNumber` 升序排列的已上传 part；每项包含 `partNumber` 和 `size`。part ETag 不会暴露。 |

### `complete_upload` 参数

| 参数名 | 类型 | 必填 | 含义 |
| --- | --- | --- | --- |
| `claimId` | `string` | 是 | 要完成的 session。服务端会从自家 part-store 还原 part 清单——不接受客户端传入 ETag。 |

成功后服务端会把已有 claim 标记为 uploaded、清掉已记录的 parts，并返回 object metadata 和 `originalFilename`。该 uploaded claim 仍等待业务采纳。同一个已 uploaded claim 后续重复调用会走幂等快速路径。
如果组装后的对象大小与 claim 不一致，action 返回
`ErrCodeUploadSizeMismatch`。

响应：

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `bucket` | `string` | provider 返回的后端 bucket 名。 |
| `key` | `string` | 最终 object key。 |
| `eTag` | `string` | 最终 object ETag。它不是 part ETag，也不会作为 `complete_upload` 入参。 |
| `size` | `int` | 最终 object 字节数。 |
| `contentType` | `string` | 对象保存的清洗后 content type。 |
| `lastModified` | timestamp | 后端 last-modified 时间。 |
| `metadata` | `object` | 可选后端 metadata map。HTTP 上传 API 不接收用户自定义 metadata。 |
| `originalFilename` | `string` | `init_upload` 时记录的原始文件名。 |

### `abort_upload` 参数

| 参数名 | 类型 | 必填 | 含义 |
| --- | --- | --- | --- |
| `claimId` | `string` | 是 | 要中止的 session。 |

响应没有 data payload。该 action 可安全重试：缺失 claim 返回成功；同 owner 下已经不是 pending 的 claim 也是 no-op。归属其他 principal 的 claim 仍会被拒绝。

最小请求示例：

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

由 schema 模块提供的结构检查资源。

### 操作列表

| Action | 访问方式 | 限流 | 作用 | 参数 |
| --- | --- | --- | --- | --- |
| `list_tables` | 需要 Bearer 认证 | 自定义 action 限流上限 `60` | 返回当前数据库或 schema 下的所有表 | 无 |
| `get_table_schema` | 需要 Bearer 认证 | 自定义 action 限流上限 `60` | 返回单个表的详细结构信息 | `GetTableSchemaParams` |
| `list_views` | 需要 Bearer 认证 | 自定义 action 限流上限 `60` | 返回当前数据库或 schema 下的所有视图 | 无 |

### `get_table_schema` 参数

| 参数名 | 类型 | 必填 | 含义 |
| --- | --- | --- | --- |
| `name` | `string` | 是 | 要检查的表名 |

## `sys/monitor`

由 monitor 模块提供的监控资源。

### 操作列表

| Action | 访问方式 | 限流 | 作用 | 参数 |
| --- | --- | --- | --- | --- |
| `get_overview` | 需要 Bearer 认证 | 自定义 action 限流上限 `60` | 返回整体系统概览快照 | 无 |
| `get_cpu` | 需要 Bearer 认证 | 自定义 action 限流上限 `60` | 返回 CPU 信息与使用情况 | 无 |
| `get_memory` | 需要 Bearer 认证 | 自定义 action 限流上限 `60` | 返回内存使用情况 | 无 |
| `get_disk` | 需要 Bearer 认证 | 自定义 action 限流上限 `60` | 返回磁盘与分区信息 | 无 |
| `get_network` | 需要 Bearer 认证 | 自定义 action 限流上限 `60` | 返回网络接口与 I/O 统计 | 无 |
| `get_host` | 需要 Bearer 认证 | 自定义 action 限流上限 `60` | 返回宿主机静态信息 | 无 |
| `get_process` | 需要 Bearer 认证 | 自定义 action 限流上限 `60` | 返回当前应用进程信息 | 无 |
| `get_load` | 需要 Bearer 认证 | 自定义 action 限流上限 `60` | 返回系统负载信息 | 无 |
| `get_build_info` | 需要 Bearer 认证 | 自定义 action 限流上限 `60` | 返回应用构建信息 | 无 |
| `get_event_streams` | 需要 Bearer 认证 | 自定义 action 限流上限 `60` | 通过可选的 `event.StreamInspector` 报告每个 redis_stream 流及其消费组状态（consumer、pending、lag、last-delivered），便于运维发现孤儿消费组 | 无 |

补充说明：

- 这些 action 都没有框架定义的输入参数
- 某些监控项在底层数据源不可用时，可能返回 monitor not ready 类错误
- 当没有可用的 `event.StreamInspector`（例如 redis_stream 传输未启用）时，`get_event_streams` 返回一个空的、`enabled = false` 的报告

最小请求示例：

```json
{
  "resource": "sys/monitor",
  "action": "get_overview",
  "version": "v1"
}
```

## Approval 资源

如果你显式引入 approval 模块，框架还会额外挂载一组 `approval/*` 资源。

实际注册的资源包括 `approval/category`、`approval/delegation`、`approval/flow`、`approval/instance`、`approval/my` 和 `approval/admin`。

这些资源已经在 [审批模块](../approval) 中按 action 展开，包括每个 action 名称、权限点、参数类型、租户规则、审计设置和限流。本页只保留索引，因为它们更偏工作流业务域，而不是框架核心通用内置资源。

## 延伸阅读

- [认证](../security/authentication)：`security/auth` 的行为与使用方式
- [文件存储](../infrastructure/storage)：`sys/storage`
- [Schema 结构检查](../infrastructure/schema)：`sys/schema`
- [监控](../infrastructure/monitor)：`sys/monitor`
