---
sidebar_position: 4
---

# 路由

VEF 支持两种路由策略，但它们共享同一个操作模型：

- RPC：通过 `POST /api`
- REST：通过 `/api/<resource>`

RPC 和 REST 都是显式资源类型。RPC 资源不会自动生成 REST 路由，REST 资源也不会复用 RPC 的单端点传输模型。

## 路由策略总览

| 策略 | 入口路径 | 操作标识来源 |
| --- | --- | --- |
| RPC | `POST /api` | 请求体里的 `resource`、`action`、`version` |
| REST | `/api/<resource>` | 资源名 + action 定义的 HTTP method 和子路径 |

## RPC 路由

RPC 请求统一进入：

```text
POST /api
```

RPC 请求形态：

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

这套形态会直接映射成 `api.Request`。

### RPC 命名规则

| 字段 | 规则 | 示例 |
| --- | --- | --- |
| `resource` | 斜杠分段的小写资源路径 | `user`、`sys/user`、`approval/category` |
| `action` | `snake_case` | `find_page`、`get_user_info`、`resolve_challenge` |
| `version` | `v<number>` | `v1`、`v2` |

### RPC 传输形式

RPC 请求解析支持：

| 内容类型 | 请求数据读取方式 |
| --- | --- |
| JSON | 请求体直接解码为 `api.Request` |
| form | `resource` / `action` / `version` 来自表单字段，`params` / `meta` 作为 JSON 字符串解析 |
| multipart form | 与 form 相同，但上传文件会并入 `params` |

## REST 路由

REST 路由统一挂在：

```text
/api/<resource>
```

HTTP method 和可选子路径由 action 字符串决定。

示例：

| Resource | Action | 最终路由 |
| --- | --- | --- |
| `users` | `get` | `GET /api/users` |
| `users` | `post` | `POST /api/users` |
| `users` | `get profile` | `GET /api/users/profile` |
| `users` | `put /:id` | `PUT /api/users/:id` |
| `users` | `delete /many` | `DELETE /api/users/many` |

### REST action 解析

REST action 字符串支持：

| 模式 | 含义 |
| --- | --- |
| `<method>` | 资源根路径 |
| `<method> <sub-path>` | 附加子路径或 Fiber 风格参数路径 |

解析规则：

- method token 会在挂载时被转换为大写 HTTP method
- 如果子路径没有以 `/` 开头，路由器会自动补上
- Fiber 风格参数，如 `/:id`，会被原样保留

### REST 命名规则

| 字段 | 规则 | 示例 |
| --- | --- | --- |
| resource name | 斜杠分段的小写路径，分段内部多词使用 kebab-case | `users`、`sys/user`、`user-profiles` |
| action method token | 小写 HTTP verb | `get`、`post`、`put`、`delete`、`patch` |
| action sub-path | kebab-case 或显式路由模式 | `profile`、`admin`、`/:id`、`/tree/options` |

## `params` 如何收集

### RPC

对 RPC 请求来说：

| 来源 | 最终落点 |
| --- | --- |
| 请求里的 `params` 对象 | `api.Request.Params` |
| multipart 上传文件 | 合并进 `api.Request.Params` |

### REST

对 REST 请求来说，VEF 会把多个来源合并进 `params`：

| 来源 | 最终落点 | 说明 |
| --- | --- | --- |
| path params | `params` | 从 Fiber route params 提取 |
| query string | `params` | 永远视为 params，而不是 meta |
| `POST` / `PUT` / `PATCH` 的 JSON body | `params` | 对象字段合并进 params |
| multipart form 字段 | `params` | 文本表单字段进入 params |
| multipart 上传文件 | `params` | 文件数组进入 params |

## `meta` 如何收集

### RPC

对 RPC 请求来说，`meta` 直接来自请求体。

### REST

对 REST 请求来说，metadata 通过 `X-Meta-` 前缀的 header 收集。

示例：

```http
X-Meta-page: 1
X-Meta-size: 20
X-Meta-format: excel
```

这些值会进入 `api.Meta`。

这里有个关键后果：

- REST query string 仍然属于 `params`
- `page.Pageable` 这类 typed helper 仍然从 `meta` 解码
- 如果某个 REST endpoint 期望 typed metadata，就要明确告诉调用方使用 `X-Meta-*`

## Typed 请求解码含义

| Handler 参数 | 解码来源 |
| --- | --- |
| 嵌入 `api.P` 的 typed struct | `params` |
| 嵌入 `api.M` 的 typed struct | `meta` |
| `page.Pageable` | `meta` |
| `api.Params` | 原始 params |
| `api.Meta` | 原始 meta |

因此，`?page=1&size=20` 并不会自动填充 typed `page.Pageable`，除非你把分页建模成普通 params 字段，而不是 meta。

## 认证解析顺序

运行时，操作的认证来源按以下顺序解析：

1. `spec.Public == true` -> 公开接口
2. 资源级 `Auth()` 配置
3. API 引擎默认认证

默认引擎认证是 Bearer。

## 内置认证策略

VEF 当前内置这些认证策略名：

| 策略 | 含义 |
| --- | --- |
| `none` | 公开接口 |
| `bearer` | Bearer token 认证 |
| `signature` | 签名认证 |

对应 helper：

| Helper | 含义 |
| --- | --- |
| `api.Public()` | 公开接口 |
| `api.BearerAuth()` | Bearer 认证 |
| `api.SignatureAuth()` | 签名认证 |

## 认证输入

### Bearer

Bearer token 可从以下来源读取：

| 来源 | 格式 |
| --- | --- |
| `Authorization` header | `Bearer <token>` |
| query 参数 | `__accessToken=<token>` |

### Signature

Signature 认证读取：

| Header | 含义 |
| --- | --- |
| `X-App-ID` | 外部应用 ID |
| `X-Timestamp` | 请求时间戳 |
| `X-Nonce` | 防重放 nonce |
| `X-Signature` | 签名值 |

## 默认操作行为

除非某个操作显式覆盖，否则 API 引擎默认会应用：

| 属性 | 默认值 |
| --- | --- |
| version | `v1` |
| timeout | `30s` |
| auth strategy | Bearer |
| rate limit | `100` requests per `5 minutes` |

## 响应形态

handler 通常通过 `result.Ok(...)` 或 `result.Err(...)` 返回响应，因此 RPC 和 REST 共享同一套响应结构：

```json
{
  "code": 0,
  "message": "Success",
  "data": {}
}
```

消息文本受当前语言影响。默认语言下通常会看到 `成功`，切到英文后通常会看到 `Success`。

## 实践建议

- 当 API 更偏动作模型时优先用 RPC
- 当 HTTP method 语义和路径结构更重要时优先用 REST
- 明确文档说明分页和排序到底走 `params` 还是 `meta`
- 不要试图掩盖 RPC 与 REST 的差异，而要把它们显式写清楚
- 用“资源 + 操作”来思考接口，而不是只盯着 endpoint URL

## 下一步

继续阅读 [参数与元信息](./params-and-meta)，看 handler 注入使用的精确解码规则。
