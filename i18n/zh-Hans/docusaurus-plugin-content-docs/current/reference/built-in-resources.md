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
| `sys/storage` | `storage` | 默认 Bearer 认证 | 文件上传、预签名 URL、临时文件删除、对象元数据、对象列表 |
| `sys/schema` | `schema` | 默认 Bearer 认证 | 数据库结构检查 |
| `sys/monitor` | `monitor` | 默认 Bearer 认证 | 运行时与宿主机监控信息 |

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

由 storage 模块提供的存储资源。

### 操作列表

| Action | 访问方式 | 限流 | 作用 | 参数 |
| --- | --- | --- | --- | --- |
| `upload` | 需要 Bearer 认证 | 继承 API 引擎默认限流 | 上传一个文件并返回对象元数据 | 通过 multipart form 解码的 `UploadParams` |
| `get_presigned_url` | 需要 Bearer 认证 | 继承 API 引擎默认限流 | 生成对象访问的临时预签名 URL | `GetPresignedURLParams` |
| `delete_temp` | 需要 Bearer 认证 | 继承 API 引擎默认限流 | 仅当对象 key 位于 `temp/` 前缀下时才允许删除 | `DeleteTempParams` |
| `stat` | 需要 Bearer 认证 | 继承 API 引擎默认限流 | 返回单个对象的元数据 | `StatParams` |
| `list` | 需要 Bearer 认证 | 继承 API 引擎默认限流 | 按前缀列出对象 | `ListParams` |

只有上面这些 action 会作为内置 RPC 接口注册出来。底层 `storage.Service` 还有 copy、move 等能力，但默认不会从这个资源直接暴露。

相关 HTTP 路由：

- `/storage/files/<key>` 是 app-level 下载代理路由，不是 RPC action
- 它不会自动继承 RPC 层的 Bearer 认证

### `upload` 参数

这个 action 要求使用 `multipart/form-data`，不能用 JSON RPC body。multipart 字段会被解码到 `params` 中。

| 参数名 | 类型 | 必填 | 含义 |
| --- | --- | --- | --- |
| `file` | `file` | 是 | 上传的文件内容，对应 `*multipart.FileHeader` |
| `contentType` | `string` | 否 | 显式覆盖内容类型；不传时会优先使用上传文件头里的内容类型 |
| `metadata` | `object<string, string>` | 否 | 传给底层存储服务的附加元数据 |

补充说明：

- 对象 key 由服务端自动生成，路径格式类似 `temp/YYYY/MM/DD/...`
- 原始文件名会自动写入 metadata
- JSON 请求会被拒绝

### `get_presigned_url` 参数

| 参数名 | 类型 | 必填 | 含义 |
| --- | --- | --- | --- |
| `key` | `string` | 是 | 目标对象 key |
| `expires` | `int` | 否 | URL 过期时间，单位秒，默认 `3600` |
| `method` | `string` | 否 | 参与签名的 HTTP 方法，默认 `GET` |

### `delete_temp` 参数

| 参数名 | 类型 | 必填 | 含义 |
| --- | --- | --- | --- |
| `key` | `string` | 是 | 要删除的对象 key，且必须以 `temp/` 开头 |

### `stat` 参数

| 参数名 | 类型 | 必填 | 含义 |
| --- | --- | --- | --- |
| `key` | `string` | 是 | 要查询元数据的对象 key |

### `list` 参数

| 参数名 | 类型 | 必填 | 含义 |
| --- | --- | --- | --- |
| `prefix` | `string` | 否 | 列表前缀过滤条件 |
| `recursive` | `bool` | 否 | 是否递归列出，而不是只列当前层级 |
| `maxKeys` | `int` | 否 | 最多返回多少个对象 |

最小请求示例：

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

由 schema 模块提供的结构检查资源。

### 操作列表

| Action | 访问方式 | 限流 | 作用 | 参数 |
| --- | --- | --- | --- | --- |
| `list_tables` | 需要 Bearer 认证 | 自定义 action 限流上限 `60` | 返回当前数据库或 schema 下的所有表 | 无 |
| `get_table_schema` | 需要 Bearer 认证 | 自定义 action 限流上限 `60` | 返回单个表的详细结构信息 | `GetTableSchemaParams` |
| `list_views` | 需要 Bearer 认证 | 自定义 action 限流上限 `60` | 返回当前数据库或 schema 下的所有视图 | 无 |
| `list_triggers` | 需要 Bearer 认证 | 自定义 action 限流上限 `60` | 返回当前数据库或 schema 下的所有触发器 | 无 |

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

补充说明：

- 这些 action 都没有框架定义的输入参数
- 某些监控项在底层数据源不可用时，可能返回 monitor not ready 类错误

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

这类资源更偏工作流业务域，因此本页不展开成逐接口明细；这里聚焦的是框架核心自带的通用内置资源。

## 延伸阅读

- [认证](../security/authentication)：`security/auth` 的行为与使用方式
- [文件存储](../features/storage)：`sys/storage`
- [Schema 结构检查](../features/schema)：`sys/schema`
- [监控](../features/monitor)：`sys/monitor`
