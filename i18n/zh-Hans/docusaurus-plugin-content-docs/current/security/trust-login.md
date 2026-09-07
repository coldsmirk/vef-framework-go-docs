---
sidebar_position: 8
---

# 信任登录（SSO 交接）

信任登录让第三方系统把用户直接带进你的应用：它跳转到框架自带的网关，URL 上携带一个已签名的用户标识，网关再把这次交接换成一次普通的登录会话。用户看不到登录表单，你的应用也永远拿不到用户密码——由外部系统为“他是谁”背书。

该功能默认关闭，在 `vef.security.trust_login`（`config.TrustLoginConfig`）下开启。

## 为什么要分成两段

一条点开就能登录的链接，本质上是一个写在 URL 里的凭据——它会落进浏览器历史、`Referer` 链，以及每一层反向代理的访问日志。因此信任登录把交接拆成两段，缺一不可：

1. **网关**——`GET /sso/trust`，一条以 `app.Middleware` 形式挂载在 order 460 的真实路由。它校验整次交接的 HMAC，把外部用户解析成本地 principal，将该身份寄存在一个一次性 code 之下，然后 302 跳转到 `<redirect>?app_id=…&code=…`。
2. **兑换**——一次普通的 `security/auth.login` 调用，`type` 为 `"trust_code"`，`principal` 填 app ID，`credentials` 填那个 code。`TrustCodeAuthenticator` 兑换 code 并返回寄存的 principal。

正因为第二段就是普通登录，这次交接原封不动地继承了**完整登录管线**：challenge 链（部门选择、强制 `password_change`）、按当前 `token_type` 签发令牌、并发会话与挤下线、以及登录审计事件。这里刻意**没有“跳过 challenge”的开关**——部门选择是必需的业务输入，强制改密是策略，跳过任何一个都是漏洞；而是否跳过二次验证属于各家自己的判断，不该由框架一刀切。

所以那条签名 URL 本身留着并不危险：它只能换来一个 code，真正能登录任何人的是 code。code 一次性、存活以秒计，并且绑定到收到它的那个浏览器。

## 配置

```toml
[vef.security.trust_login]
enabled = true
# path = "/sso/trust"      # config.DefaultTrustLoginPath
# code_ttl = "60s"         # config.DefaultTrustLoginCodeTTL
# bind_user_agent = true   # 省略即为启用
# bind_client_ip = false

[vef.security.trust_login.rate_limit]
# max = 120                # config.DefaultTrustLoginRateLimitMax
# period = "1m"            # config.DefaultTrustLoginRateLimitPeriod

[vef.security.trust_login.apps.his]
redirect_urls = ["https://portal.example.com/sso/callback"]
```

| 配置项 | 类型 | 含义 |
| --- | --- | --- |
| `enabled` | `bool` | 挂载网关路由并注册 `trust_code` 认证器。关闭时该路由根本不存在，`trust_code` 会被当作不支持的登录类型拒绝。 |
| `path` | `string` | 网关路由。它是签名载荷的一部分，改动它会让外部系统已经生成的所有链接全部失效。 |
| `code_ttl` | `time.Duration` | 已签发的 code 可兑换的时长。要设得短——它是跟着重定向 URL 走的。 |
| `bind_user_agent` | `*bool` | 要求兑换 code 的浏览器提交与网关重定向时相同的 `User-Agent`。省略时解析为启用：一次重定向之内请求头不会变，因此这条绑定不花任何代价。 |
| `bind_client_ip` | `bool` | 额外要求来源地址一致。默认关闭——移动端可能在重定向途中切换网络，由此产生的失败读起来像集成坏了，而不像一道防线。 |
| `apps.<appId>.redirect_urls` | `[]string` | 该 app 发起的交接允许落到哪些地址。`config.TrustLoginAppConfig` 只有这一个字段。 |
| `rate_limit.max` / `rate_limit.period` | `int` / `time.Duration` | `config.TrustLoginRateLimitConfig`，按 (app ID, 客户端 IP) 逐节点计数。 |

默认值以 `Effective*` 访问器的形式写在字段旁边，解析后的配置结构体不会被就地改写。

### 准入是一项独立授权

一个 app 必须出现在 `apps` 里才允许发起交接。这是刻意为之：**能调用你的 API，绝不等于能把某个用户登录进来。** 签名密钥来自 `security.ExternalAppLoader`——与 API 签名认证读取的是同一处，所以密钥不必写进配置文件——但真正授予浏览器单点登录的是 `apps` 这张表，不在其中的 app 一律拒绝。

重定向目标必须与某条白名单的 scheme 和 host **完全一致**，并且位于其 path 之下；比较按路径分段进行，`.`/`..` 会先归一化。

## 校验失败

`SecurityConfig.Validate` 会让启动直接失败，而不是带着一个半配好的网关跑起来：

| 错误 | 原因 |
| --- | --- |
| `config.ErrTrustLoginPathInvalid` | `path` 不以 `/` 开头。 |
| `config.ErrTrustLoginAppsRequired` | 信任登录已开启但 `apps` 为空——没有任何人能发起交接。 |
| `config.ErrTrustLoginRedirectsEmpty` | 某个 app 没有声明 `redirect_urls`，它永远无法完成一次交接。 |
| `config.ErrTrustLoginRedirectInvalid` | 某条重定向不是绝对 URL，或带有 query、fragment。 |

## 如何签名

网关校验的 HMAC 覆盖 app ID、请求方法与路径、时间戳与 nonce，**以及**两个真正关键的调用方参数：`user_id` 和 `redirect`。这两个参数通过 `security.SignatureRequest` 传入：

```go
req := security.SignatureRequest{
    AppID:  "his",
    Method: "GET",
    Path:   "/sso/trust",
    BoundParams: map[string]string{
        "user_id":  "1001",
        "redirect": "https://portal.example.com/sso/callback",
    },
}
```

判断准则：**签名链接上带着、却不在 `BoundParams` 里的参数，就是攻击者可控的参数。** 未签名的 `user_id` 就是一次彻底的认证绕过——任何拿到一条有效链接的人，都可以把别人的标识替换进去。

`BoundParams` 的值是**解码后**的值，绝不是它们在 URL 上的编码形态。键和值会被百分号编码进规范化载荷（RFC 3986 unreserved 集合、大写十六进制、空格编码为 `%20`），原因正是签名的 `redirect` 里通常就带着 `&` 和 `=`；若原样拼接，两组不同的参数可能压平成同一个字符串。与固定字段（`app_id`、`method`、`nonce`、`path`、`timestamp`）同名的键会被 `security.ErrSignatureBoundKeyReserved` 拒绝，而不是悄悄把固定字段盖掉。

最终的 URL 形如：

```text
GET /sso/trust?app_id=his&user_id=1001&redirect=<编码后>&timestamp=<Unix 秒>&nonce=<随机串>&signature=<hex hmac>
```

有两个尖锐的细节，值得明确告诉实现签名那一侧的人：

- **签名里的 `path` 必须完全一致。** Fiber 的 `StrictRouting` 是关闭的，所以 `/sso/trust/` 同样能到达网关，但它参与哈希的结果不同，只会得到一个永久且不解释原因的 401。
- 一个**不带任何绑定参数**的请求，其载荷与普通 API 签名认证逐字节相同。因此已经会签 API 调用的第三方，只需要多加这两个参数，而不必再学一套新方案。

## 解析用户

`security.TrustUserResolver` 负责把外部标识映射成本地 principal：

```go
type TrustUserResolver interface {
    ResolveUser(ctx context.Context, appID, externalUserID string) (*Principal, error)
}
```

注册它是可选的。不注册时网关走 `UserLoader.LoadByID`——只要两个系统用同一套用户标识，这就已经是正确行为。当标识不一致，或映射关系取决于是哪个外部系统发起的交接时，才需要自己实现。

:::caution[解析本身就是一次认证决策]
外部系统背书的是用户**是谁**，而不是你的应用**是否仍然接纳**他。因此无论最终跑的是你的 resolver 还是 `LoadByID`，都必须像密码登录那样，对已禁用、已锁定、已过期的账号返回 nil principal 予以拒绝。很多应用只把这类检查放在 `LoadByUsername` 里，而信任登录根本不会走到那里。
:::

返回 error 表示“解析失败”（例如目录服务不可达）；“这里没有这个用户”应该返回 nil principal，而不是 error。

## 一次性 code

`security.TrustCodeStore` 寄存已验证的身份，并且只能兑换一次：

```go
type TrustCodeStore interface {
    Issue(ctx context.Context, state TrustCodeState, ttl time.Duration) (string, error)
    Consume(ctx context.Context, code string) (*TrustCodeState, error)
}
```

`security.TrustCodeState` 携带解析出的 `Principal`、发起交接的 `AppID`，以及网关重定向的那个浏览器的 `UserAgent` 与 `ClientIP`。两个浏览器字段都是无条件记录的，因此把某条绑定打开之后，从下一次交接起就会生效，而不需要重新签发什么。

兑换时会把记录的 `UserAgent` 与 `contextx.RequestUserAgent(ctx)` 比对——后者由 API 认证中间件通过 `contextx.SetRequestUserAgent` 写在 `contextx.KeyRequestUserAgent` 键下——在 `bind_client_ip` 打开时，还会比对 `contextx.RequestIP(ctx)`。

`Consume` 必须原子地读取并作废：一次性正是这个 code 的全部安全价值所在，躺在浏览器历史里的那条重定向 URL 绝不能被重放。未知的、已兑换的、已过期的 code 一律返回 `security.ErrTrustCodeInvalid`——三者被刻意做成无法区分。

### 存储按部署拓扑选择，而不是靠装饰替换

与 `lock.Locker` 一样，code 存储是**根据部署形态选出来的**，而不是默认内存实现再等人来换：

| Redis 客户端 | 存储 |
| --- | --- |
| 可用 | `security.RedisTrustCodeStore`（`security.NewRedisTrustCodeStore`），依赖 `GETDEL`，需要 Redis 6.2+ |
| 不可用 | `security.MemoryTrustCodeStore`（`security.NewMemoryTrustCodeStore`），并打印启动警告 |

之所以不能是一个安静的内存默认值：内存存储**根本无法**服务第二个副本。一个节点签发的 code 对其他节点来说完全未知，兑换会直接失败，而不只是安全性变弱。共享的 `security.NonceStore` 按同样方式选择，而它的启动警告只在信任登录开启时才发出——按进程存放 nonce 对服务器间的 API 调用只是削弱，但对一个浏览器可见的凭据来说，等于每个节点都送它一次新的兑换机会。

code 由 `GenerateOpaqueToken` 生成，并以 `HashOpaqueToken` 作为键，从不明文存储，因此存储泄露不会泄露任何可用的 code。

## 错误

网关的失败以**纯文本**加对应状态码返回，绝不使用 JSON 信封——这个端点是浏览器直接导航过去的。

| 哨兵 | 业务码 | 状态码 | 含义 |
| --- | --- | --- | --- |
| `security.ErrTrustAuthFailed` | `security.ErrCodeTrustAuthFailed`（1060） | 401 | 所有校验失败：签名错误、nonce 过期或重放、app 未知或已禁用、时间戳格式错误。 |
| `security.ErrTrustRedirectNotAllowed` | `security.ErrCodeTrustRedirectNotAllowed`（1061） | 400 | 交接本身通过了认证，但重定向不匹配任何白名单条目。 |
| `security.ErrTrustUserNotResolved` | `security.ErrCodeTrustUserNotResolved`（1062） | 401 | 外部标识在本地找不到对应用户。 |
| `security.ErrTrustCodeInvalid` | `security.ErrCodeTrustCodeInvalid`（1063） | 401 | code 未知、已兑换、已过期，或由另一个浏览器提交。 |

所有校验失败都塌缩成 `ErrTrustAuthFailed`，这样就无法用它枚举出哪些 app ID 存在——而且调用方在每一种情况下的修法都一样：把交接重新签对。`ErrTrustRedirectNotAllowed` 被刻意区分开，因为它的修法是改配置。

那次 302 会带上 `Cache-Control: no-store` 和 `Referrer-Policy: no-referrer`，让 `Location` 里的 code 和请求 URL 里的签名都不会传得比必要范围更远。

## 限流

网关是公开的，而且在能够拒绝任何东西之前就必须做一次 `ExternalAppLoader` 查询——通常是一次数据库往返——因此它通过 `rate_limit` 按 (app ID, 客户端 IP) 逐节点限流。默认值刻意放得宽松：一整个组织共用一个 NAT 出口，在上班时段集中登录是完全正常的。它约束的是洪水，而不是登录行为本身。

兑换那一段**不受暴力破解登录守卫约束**，这一点与其他所有登录类型都不同。code 是 `GenerateOpaqueToken` 级别的随机数，一次性且存活以秒计，再多次尝试也逼近不了它；而它出示的身份是发起系统的 app ID，不是某个人。统计失败次数只会把该系统的所有用户塞进同一个锁定桶里，而任何能 POST 一个错误 code 的人都能把它填满。它真正的限流是上面的网关限流、code 的一次性，以及 `vef.security.login_rate_limit`。

## 发起方接入清单

1. 通过你的 `ExternalAppLoader` 的供给渠道，取得 app ID 与签名密钥。
2. 请 VEF 应用的运维在 `vef.security.trust_login.apps` 下加入该 app 及其重定向白名单。
3. 对 `app_id + method + path + user_id + nonce + redirect + timestamp` 签名，其中 `user_id` 与 `redirect` 作为绑定参数传入。
4. 把浏览器导航到 `<base>/sso/trust?...`。
5. 在回调页面从 query 中读出 `code`，调用 `security/auth.login`，`type` 为 `"trust_code"`、`principal` 为 app ID、`credentials` 为 code。challenge 响应的处理方式与密码登录完全一致。
