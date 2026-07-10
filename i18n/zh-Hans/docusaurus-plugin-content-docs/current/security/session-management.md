---
sidebar_position: 7
---

# 会话管理

VEF 的登录层支持两种相互独立的令牌机制，由 `vef.security.token_type` 选择。本页介绍其中的有状态机制——opaque token——以及随之而来的服务端会话控制能力。如果你的应用完全没有配置过 `token_type`，那么当前运行的就是无状态的默认机制，本页内容暂时都不适用。

## `jwt_token` 与 `opaque_token` 的取舍

| | `jwt_token`（默认） | `opaque_token` |
| --- | --- | --- |
| 令牌形态 | 自包含 JWT，principal 编码在 claims 中 | 随机引用令牌，principal 存放在服务端 |
| 每次请求的服务端查询 | 无 | 每次请求一次 `SessionStore.Lookup` |
| 过期前主动吊销 | 不可能 | `logout`、管理员吊销、强制下线均可 |
| 并发会话数限制 | 不可能 | `max_concurrent` + `on_exceed` |
| “在线设备”列表 / 挤下线 | 不可能 | `ListByUser`、`Revoke`、`RevokeUser` |
| 无共享状态即可跨节点扩展 | 可以 | 只有搭配共享存储（Redis）才可以 |

无状态 JWT 的校验成本几乎为零，也不需要任何共享基础设施，但一旦签发出去，服务端就再也无法收回——没有会话可吊销。opaque token 用“每次请求多一次查询”换来了真正的会话控制：你可以强制某台设备下线、限制单账号同时在线的会话数，并且能看到谁正在登录。

切换机制只需要改一个配置项：

```toml
[vef.security]
token_type = "opaque_token" # 默认："jwt_token"
```

`token_type` 的类型是 `config.TokenType`，仅有的两个合法取值是常量
`config.TokenTypeJWT`（`"jwt_token"`）和 `config.TokenTypeOpaque`
（`"opaque_token"`）——写错值会在配置校验阶段（`SecurityConfig.Validate`）
直接导致启动失败。无论选哪种机制，`AuthResource.Login`、`refresh`、`logout` 在客户端看来行为完全一致；不同的只是令牌形态和服务端具备的能力。`TokenGenerator.Generate(ctx, principal, SessionMeta)` 与具体机制无关，因此自定义登录流程不需要关心当前激活的是哪一种。

## opaque 会话的工作原理

在 `opaque_token` 机制下，登录不会签发 JWT，而是由 `OpaqueTokenGenerator`：

1. 生成一个高熵随机令牌（`security.GenerateOpaqueToken`）
2. 开启一条 `security.Session` 记录——包含 id、用户 id、`*security.Principal` 的快照、客户端 IP、User-Agent 以及时间戳——以该令牌的 SHA-256 哈希（`security.HashOpaqueToken`）为键，而不是以原始令牌为键
3. 将其存入某个 `security.SessionStore`
4. 将原始令牌作为 `AuthTokens.AccessToken` 返回（没有 refresh token——会话本身会在使用时自我续期，因此不存在需要单独刷新的东西）

在之后每一次已认证请求中，`OpaqueTokenAuthenticator` 会对请求携带的 bearer 令牌取哈希，并调用 `SessionStore.Lookup`。命中后直接返回该会话中的 `Principal` 快照——无需再查一次数据库，因为这就是会话开启时（或最近一次续期时）生效的用户数据。

### 滑动空闲超时，受绝对生命周期上限约束

一条会话受两条相互独立的生命周期约束：

- **`idle_ttl`**——会话在没有活动时能存活多久。启用滑动续期时，每一次已认证请求都会把它再向后延长一个 `idle_ttl`。
- **`max_lifetime`**——从 `CreatedAt` 起算的会话总时长硬上限，无论活动多频繁都无法突破。

`OpaqueTokenAuthenticator.renew` 将下一次过期时间计算为 `now + idle_ttl`，如果这会超过 `CreatedAt + max_lifetime`，则将其钳制到该上限。也就是说，即便会话持续活跃，也仍会在登录后最晚 `max_lifetime` 时过期——滑动续期只是延长空闲存活时间，并不会延长账号的绝对会话预算。

```toml
[vef.security.session]
idle_ttl = "30m"    # 默认值
max_lifetime = "168h" # 默认值：7 天
sliding = true       # 默认值；省略或设为 false 可关闭续期
```

- `idle_ttl` 默认值：`30m`（`config.DefaultSessionIdleTTL`）
- `max_lifetime` 默认值：`168h`（7 天，`config.DefaultSessionMaxLifetime`）
- `sliding` 是 `*bool`：省略该键时默认启用；需显式设为 `false` 才会关闭续期，使会话严格在登录后 `idle_ttl` 到期，与活动情况无关

续期是尽力而为（best-effort）：续期时 `SessionStore` 出错只会被记录日志并吞掉，不会作为请求失败反馈给客户端，因此存储层的短暂抖动不会把一个本应有效的用户强制登出。`ExpiresAt` 是两个内置实现共同的唯一权威过期判断字段——Redis 存储会在续期时同步刷新自身键的 TTL，但即便某个键因为任何原因存活超过了 `ExpiresAt`，读取时仍会被判定为已过期。

## 并发控制（挤下线）

`vef.security.session.max_concurrent` 限制单账号可同时持有的会话数。`0`（默认值）表示不限制。当一次登录会使账号超出上限时，`on_exceed` 决定接下来的行为：

```toml
[vef.security.session]
max_concurrent = 3
on_exceed = "evict_oldest" # 默认值；或 "reject"
```

- **`evict_oldest`**（默认，`security.SessionExceedEvictOldest`）——新登录会吊销该账号最旧的若干会话，直到把新会话计入后恰好达到 `max_concurrent`。这就是“挤下线”：新设备登录会悄悄把旧设备踢下线。
- **`reject`**（`security.SessionExceedReject`）——直接拒绝新登录。`Login` 返回 `security.ErrTooManyConcurrentSessions`（业务码 `1024`——即 `security.ErrCodeTooManyConcurrentSessions`，HTTP `403`），账号已有的会话保持不变。

限制逻辑在 `OpaqueTokenGenerator.enforceConcurrency` 中执行，发生在新会话创建之前，并且**在并发登录场景下是尽力而为**：统计现有会话数和创建新会话是两次独立的存储调用，而非一次原子操作，因此同一账号的一波并发登录可能会短暂地超出 `max_concurrent`，超出的数量取决于同时竞争的请求数。这是一条策略性/爆炸半径限制，而非硬安全边界——请据此看待它；`evict_oldest` 策略会在下一次登录时自我修复任何短暂的超额。

这些配置最终会解析为一个 `security.SessionPolicy`——`MaxConcurrent`、
`OnExceed`（类型为 `security.SessionExceedPolicy`）、`IdleTTL`、
`MaxLifetime`、`Sliding`——由 `vef.security.session` 一次性组装完成，供
opaque token 的 generator 和 authenticator 共享使用，而不是在每次请求时
重新读取原始配置。

## 登出与吊销

`security/auth.logout` 会吊销当前 bearer 令牌背后的会话：

```go
func (a *AuthResource) Logout(ctx fiber.Ctx) error {
	a.revokeCurrentSession(ctx)

	return result.Ok().Response(ctx)
}
```

它按令牌哈希查找会话，并调用 `SessionStore.Revoke(ctx, session.ID)`。这是尽力而为的操作，且始终返回 `Ok`——会话不存在（已过期，或者当前是 JWT 令牌、本来就没有会话）不算错误，吊销过程中的存储失败也只会被记录日志。在 `jwt_token` 机制下，`logout` 实际上是空操作：没有会话可吊销，客户端需要自行丢弃已保存的令牌（参见[身份认证](./authentication)）。

## 构建会话管理端点

框架内部使用的 `security.SessionStore` 是一个常规的、通过 DI 暴露出来的依赖，而不是私有实现细节——你可以把它注入自己的资源，构建“我的设备”或后台会话管理功能：

```go
type SessionResource struct {
	api.Resource
	store security.SessionStore
}

func NewSessionResource(store security.SessionStore) api.Resource {
	return &SessionResource{store: store, /* ... */}
}

// ListMyDevices 返回调用者自己的存活会话。
func (r *SessionResource) ListMyDevices(ctx fiber.Ctx, principal *security.Principal) error {
	sessions, err := r.store.ListByUser(ctx.Context(), principal.ID)
	if err != nil {
		return err
	}

	return result.Ok(sessions).Response(ctx)
}
```

`SessionStore` 恰好提供了后台管理界面所需的能力：

- `ListByUser(ctx, userID)`——某个用户自己的存活会话，按最近活动时间倒序——可作为自助“在线设备”列表的基础
- `Revoke(ctx, id)`——按公开的 `Session.ID`（而非令牌）吊销单个会话——“把这台设备踢下线”
- `RevokeUser(ctx, userID)`——一次性吊销某用户的全部会话——强制登出，例如用于密码重置或账号封禁场景

`Session.ID` 是一个刻意与令牌哈希分离的随机公开标识符，因此可以安全地在设备列表中返回给客户端，而不会暴露任何能反推出有效凭证的信息。

构建这些端点及其鉴权完全是你自己的责任——框架不提供管理后台界面，也不会为 `SessionStore` 提供默认的鉴权规则，只提供存储契约。请像对待其他涉及权限敏感操作的端点一样应用[权限控制](./authorization)（通常应把 `RevokeUser` 以及跨用户读取限制给管理员角色）。

### 通过 `SessionInspector` 实现跨用户可见性

如果要构建一个覆盖所有用户的“全部在线会话”看板，可以对 store 做类型断言，判断其是否实现了可选的 `security.SessionInspector` 能力（这与框架其他地方使用的 `event.StreamInspector` 模式一致）：

```go
type SessionInspector interface {
	ListAll(ctx context.Context) ([]Session, error)
}
```

```go
if inspector, ok := r.store.(security.SessionInspector); ok {
	sessions, err := inspector.ListAll(ctx.Context())
	// ...
}
```

两个内置存储实现（`MemorySessionStore`、`RedisSessionStore`）都实现了该接口。`ListAll` 的复杂度是 `O(全部会话数)`——Redis 实现通过 keyspace `SCAN` 完成，而不是维护一个全局索引集合，因此不会为一次性账号或已删除账号积累残留数据。请把它当作低频的管理端读取操作，而不是请求路径上的调用；规模大到需要分页的部署应当在自己的存储之上另行实现。

## 内存 vs. Redis：单节点 vs. 多节点

默认的 `SessionStore` 是 `security.NewMemorySessionStore()`——基于进程内 map 实现，由框架的安全模块自动装配。对单实例部署而言完全可用：会话保存在三张内存索引表中（按 id、按令牌哈希、按用户），过期条目在被访问时惰性回收。

它**不会**跨进程共享状态。在多节点部署中，某个节点上创建的会话在另一个节点上不可见——需要通过 `fx.Decorate` 把存储替换为 `security.NewRedisSessionStore`：

```go
vef.Run(
	// ...
	fx.Decorate(security.NewRedisSessionStore),
)
```

`NewRedisSessionStore(client *redis.Client) security.SessionStore` 只需要 `internal/redis` 在 `vef.redis.enabled = true` 时已经提供的 `*redis.Client`——除了启用 Redis 并 decorate 该存储之外，不需要额外的装配代码。

```toml
[vef.redis]
enabled = true
# host、port 等其他字段——参见 Redis 配置参考

[vef.security]
token_type = "opaque_token"
```

`RedisSessionStore` 与内存存储保持相同的语义，底层用 `vef:security:session:` 前缀（子前缀 `id:`、`token:`、`user:`）的 Redis 键实现：

- 每一次涉及多个键的变更（`Create`、`RevokeUser`、吊销时的删除）都在一个 Redis `MULTI`/`EXEC` 事务中执行，因此读取方永远不会看到一个写了一半或删了一半的会话
- 滑动续期对会话记录执行 `SET ... XX`，只有键仍然存在时才会成功——这样一次与并发 `Revoke` 竞争的续期就永远不可能“复活”一个刚被删除的会话
- 每次读取时，`Session.ExpiresAt`（而不仅仅是 Redis 键的 TTL）仍然是权威的过期判断依据，因此两种存储对 `max_lifetime` 的执行是一致的

两种存储都可安全并发使用，并且通过同一个 `security.SessionStore` 接口互为直接替换——切换存储时，应用的其余部分不需要做任何改动。

## 参见

- [身份认证](./authentication)——每次请求如何校验令牌，以及内置的 `security/auth` 资源
- [登录加固](./login-hardening)——暴力破解锁定、密码强度与密码历史，这些无论 `token_type` 取何值都适用
