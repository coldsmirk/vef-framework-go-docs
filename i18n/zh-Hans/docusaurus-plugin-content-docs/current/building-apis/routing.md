---
sidebar_position: 2
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

内部 router 常量 `DefaultRPCEndpoint` 是 `/api`。form 和 multipart RPC
传输会分别从 `FormKeyParams` 与 `FormKeyMeta` 表单字段读取 JSON 编码的
`params` 和 `meta`。

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

挂载时，REST 路由会在 action method 转成大写、可选 sub-resource path 补上
`/` 前缀后，形成 `/api/<resource>/<subpath>`。

HTTP method 和可选子路径由 action 字符串决定。

示例：

| Resource | Action | 最终路由 |
| --- | --- | --- |
| `users` | `get` | `GET /api/users` |
| `users` | `post` | `POST /api/users` |
| `users` | `get profile` | `GET /api/users/profile` |
| `users` | `put profile` | `PUT /api/users/profile` |
| `users` | `delete many` | `DELETE /api/users/many` |

### REST action 解析

REST action 字符串支持：

| 模式 | 含义 |
| --- | --- |
| `<method>` | 资源根路径 |
| `<method> <sub-resource>` | 附加 kebab-case 子资源路径 |

解析规则：

- method token 会在挂载时被转换为大写 HTTP method
- 公开资源校验只接受小写 HTTP verb，以及可选的 kebab-case 子资源路径；
  `admin/users` 这类斜杠分段子路径是允许的，但每一段都必须是 kebab-case
- `/:id` 这类动态 Fiber route 参数不会通过 `api.ValidateActionName` /
  `api.NewRESTResource` 的公开校验

### REST 命名规则

| 字段 | 规则 | 示例 |
| --- | --- | --- |
| resource name | 斜杠分段的小写路径，分段内部多词使用 kebab-case | `users`、`sys/user`、`user-profiles` |
| action method token | 小写 HTTP verb | `get`、`post`、`put`、`delete`、`patch` |
| action sub-resource | 可选的斜杠分段 kebab-case 路径 | `profile`、`admin/users`、`user-friends` |

## 框架保留的 Query Key

框架会保留一些 query key 用于自身 plumbing。REST 路由器在构建业务参数空间之前会先把它们剥离，因此它们不会出现在 `params` 中。

| 保留 key | 用途 |
| --- | --- |
| `__accessToken` | 用于无法设置 `Authorization` 头的请求（例如浏览器 `<img src>` 或下载链接）传递 bearer token；认证层直接从 Fiber context 读取。 |

由于这些 key 被框架占用，它们不会进入业务 handler，也不会被报告为未映射键。

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
| `POST` / `PUT` / `PATCH` 的 multipart form 字段 | `params` | 文本表单字段进入 params |
| `POST` / `PUT` / `PATCH` 的 multipart 上传文件 | `params` | 文件数组进入 params |

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

去掉 `X-Meta-` 前缀后，存入 `meta` 的 key 会统一转成小写。例如
`X-Meta-page` 和 `X-Meta-Page` 都会变成 `meta["page"]`。

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
| `ip` | 来源 IP 白名单认证 |
| `api_key` | 静态 API key 认证 |
| `http_basic` | RFC 7617 Basic 认证 |

对应 helper：

| Helper | 含义 |
| --- | --- |
| `api.Public()` | 公开接口 |
| `api.BearerAuth()` | Bearer 认证 |
| `api.SignatureAuth()` | 签名认证 |
| `api.IPAuth(...)` | 来源 IP 白名单认证 |
| `api.APIKeyAuth(...)` | API key 认证；可选自定义头名 |
| `api.HTTPBasicAuth()` | HTTP Basic 认证 |

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

### IP

IP 认证通过 `security.IPWhitelistLoader` 解析命名白名单（默认实现读取
`vef.security.ip_whitelists` 配置）。位于反向代理之后时，需要配置
`vef.app.trusted_proxies`，让 Fiber 解析到真实客户端 IP。所有认证失败都会以
`security.ErrIPNotAllowed` 拒绝；空白名单或未找到的白名单会 fail-closed。

### API key

| 来源 | 格式 |
| --- | --- |
| 请求头（默认 `X-API-Key`，可按操作经 `api.APIKeyAuth("X-My-Key")` 自定义） | 原始密钥值 |

密钥经 `security.APIKeyLoader` 解析（默认由 `vef.security.api_keys` 支撑）；
缺失或未匹配的密钥统一以 `security.ErrAPIKeyInvalid`（401）拒绝。

### HTTP Basic

| 来源 | 格式 |
| --- | --- |
| `Authorization` header | `Basic base64(username:password)` |

账号经 `security.BasicAccountLoader` 解析（默认由
`vef.security.basic_accounts` 支撑），常数时间比较；畸形请求头、未知账号
与错误密码统一以 `security.ErrBasicCredentialsInvalid`（401）拒绝。

## 默认操作行为

除非某个操作显式覆盖，否则 API 引擎默认会应用：

| 属性 | 默认值 |
| --- | --- |
| version | `v1` |
| timeout | `30s` |
| auth strategy | Bearer |
| rate limit | `100` requests per `5 minutes` |

## 请求 Body 传输编码

对于 `/api` 表面的请求，客户端可以选择对请求 body 进行传输编码，让携带脚本或代码形态的 payload（集成适配器脚本、信封/认证脚本、`dry_run` body）能够穿过会对这些内容误报的中间件。编码会在调度器解析之前被解码回原始 JSON。

客户端通过请求头启用：

```http
X-Body-Encoding: base64
```

或：

```http
X-Body-Encoding: gzip+base64
```

支持的值：

| 值 | 含义 |
| --- | --- |
| `base64` | body 是原始 JSON 的 base64 |
| `gzip+base64` | body 是原始 JSON 经 gzip 压缩后再 base64（解码顺序：先 base64，再 gunzip） |

原生的 `Content-Encoding`（如 `gzip`、`br`、`deflate`、`zstd`）由 Fiber 自身解压；这个中间件只覆盖 Fiber 不认识的 base64 形式。

### Protected body transport

设置 `vef.api.body_encoding.enabled = true` 后，`/api` 表面上的 JSON body
会启用带认证的非明文传输：

```toml
[vef.api.body_encoding]
enabled = true
encoding = "aes-gcm+base64" # 或 "sm4-gcm+base64"
key = "<standard-base64-key>"
```

两种编码都使用 GCM，wire body 是裸 standard Base64 文本，没有 JSON 外壳。
解码前的字节布局是 `Base64(12 字节随机 nonce || ciphertext || 16 字节 tag)`。
每条请求和响应都会生成新的 nonce。

启用 protected transport 后：

- JSON 请求必须把 `X-Body-Encoding` 设置为配置值，并把编码后的文本作为
  raw request body 发送，而不是作为 JSON string 发送。
- JSON 响应，包括传输层和下游错误信封，都使用同一编码，并通过
  `X-Body-Encoding` 返回标记；CORS 会 expose 这个头，让浏览器可以读取。
- 明文 JSON 以及较弱的 `base64` / `gzip+base64` 请求会被拒绝，不能把
  protected 模式降级。
- multipart 与二进制 body 保持原生格式，非 `/api` 响应不受影响。
- 解码后的请求 body 仍然受 `vef.app.body_limit` 限制。
- protected `/api` 响应会跳过 HTTP compression，因为 ciphertext 本身已是
  高熵内容。

AES key 可以是 16、24 或 32 字节，SM4 key 必须是 16 字节。standard-base64
key 必须与客户端一致。这个传输层不能替代 HTTPS：浏览器代码持有的对称
key 是可以被提取的。API `signature` 认证与它正交，因为签名覆盖
`appID+method+path+timestamp+nonce`，不包含请求 body。

解码发生在 content-type 检查和调度器之前，且仅在 `/api` 路径上生效。存储、内容哈希缓存和审计看到的都是原始 body。body 限制（`vef.app.body_limit`）作用于解码后的大小，防止解压炸弹超出限制。

失败模式：

| 情况 | 结果 |
| --- | --- |
| 启用 protected transport 后，JSON 请求未携带标记 | `api.ErrBodyEncodingRequired`（HTTP 400） |
| 不支持的 `X-Body-Encoding` 值 | `api.ErrUnsupportedBodyEncoding`（HTTP 400） |
| base64 格式错误或 gzip 流损坏 | `api.ErrBodyDecodeFailed`（HTTP 400） |
| 解码后的 body 超过配置限制 | `api.ErrBodyTooLarge`（HTTP 413） |

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
