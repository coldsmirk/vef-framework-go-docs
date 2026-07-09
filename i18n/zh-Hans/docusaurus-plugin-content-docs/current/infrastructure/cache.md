---
sidebar_position: 1
---

# 缓存

VEF 把 `github.com/coldsmirk/vef-framework-go/cache` 暴露为一个类型化工具包。它不会创建隐藏的全局缓存，也不会把 cache 实例自动注入到 resource 中。

## 公开 API 面

该包没有 exported fields。公开 top-level API 如下：

| API | 契约 |
| --- | --- |
| `cache.Cache[T]` | generic cache interface，由内存和 Redis 后端实现 |
| `cache.NewMemory[T](opts ...cache.MemoryOption) cache.Cache[T]` | 创建进程内内存缓存 |
| `cache.NewRedis[T](client *redis.Client, namespace string, opts ...cache.RedisOption) cache.Cache[T]` | 创建按 `namespace` 隔离的 Redis 缓存；`client` 为 nil 或 `namespace` 为空时 panic |
| `cache.MemoryOption` | `NewMemory` 接收的 functional option 类型 |
| `cache.RedisOption` | `NewRedis` 接收的 functional option 类型 |
| `cache.WithMemMaxSize(size int64)` | 按条目数量限制内存缓存；`size <= 0` 表示无限制 |
| `cache.WithMemDefaultTTL(ttl time.Duration)` | 当 `Set` 或 `GetOrLoad` 没有传入正数 per-call TTL 时使用的默认内存 TTL |
| `cache.WithMemEvictionPolicy(policy cache.EvictionPolicy)` | 在正数最大条目数生效时选择内存淘汰策略 |
| `cache.WithMemGCInterval(interval time.Duration)` | 控制内存过期条目清理周期；`interval <= 0` 回退到 `5m` |
| `cache.WithRdsDefaultTTL(ttl time.Duration)` | Redis 写入没有正数 per-call TTL 时使用的默认 TTL |
| `cache.EvictionPolicy` | 内存淘汰策略枚举 |
| `cache.EvictionPolicyNone` | 不做淘汰跟踪；无限制内存缓存会自动使用它 |
| `cache.EvictionPolicyLRU` | 淘汰最近最少使用的条目；读取和更新会刷新 recency |
| `cache.EvictionPolicyLFU` | 淘汰使用频率最低的条目；频率相同按该频率桶内的插入顺序处理 |
| `cache.EvictionPolicyFIFO` | 淘汰最早插入的条目；读取不会改变顺序 |
| `cache.Key(keyParts ...string) string` | 用默认 key builder 以 `:` 拼接片段；`cache.Key()` 返回 `""` |
| `cache.KeyBuilder` | 只包含 `Build(keyParts ...string) string` 的接口 |
| `cache.PrefixKeyBuilder` | 基于 prefix 的 `KeyBuilder` 实现 |
| `cache.NewPrefixKeyBuilder(prefix string) *cache.PrefixKeyBuilder` | 创建使用 `:` 作为分隔符的 key builder |
| `cache.NewPrefixKeyBuilderWithSeparator(prefix, separator string) *cache.PrefixKeyBuilder` | 创建使用自定义 separator 的 prefix key builder；separator 会按传入值原样使用 |
| `cache.LoaderFunc[T]` | `func(ctx context.Context) (T, error)`，供 `Cache.GetOrLoad` 和 `SingleflightMixin.GetOrLoad` 使用 |
| `cache.KeyedLoaderFunc[T]` | `func(ctx context.Context, key string) (T, error)`，供 `Invalidating` 使用 |
| `cache.GetFunc[T]` | `func(context.Context, string) (T, bool)`，供 `SingleflightMixin` 调用 |
| `cache.SetFunc[T]` | `func(context.Context, string, T, ...time.Duration) error`，供 `SingleflightMixin` 调用 |
| `cache.SingleflightMixin[T]` | 给自定义 cache 实现复用的 `GetOrLoad` 逻辑 |
| `cache.Invalidating[T]` | 支持显式失效的 read-through 内存缓存包装器 |
| `cache.NewInvalidating[T](loader cache.KeyedLoaderFunc[T], logger logx.Logger, opts ...cache.MemoryOption) *cache.Invalidating[T]` | 创建 invalidating cache；调用方必须传入非 nil loader 和 logger |
| `cache.ErrMemoryLimitExceeded` | 内存缓存达到正数条目上限且找不到可淘汰候选 |
| `cache.ErrCacheClosed` | cache 关闭后触发写入路径 |
| `cache.ErrLoaderRequired` | `GetOrLoad` 收到 nil loader |
| `cache.ErrTypeAssertionFailed` | `SingleflightMixin` 收到的结果不是当前 generic 类型 `T` |

## `cache.Cache[T]`

两个内置构造器都返回 `cache.Cache[T]`。该接口嵌入 `io.Closer`，并暴露以下方法：

| 方法 | 签名 | 契约 |
| --- | --- | --- |
| `Get` | `Get(ctx context.Context, key string) (T, bool)` | miss、过期、后端读取失败、反序列化失败或 cache 已关闭时返回 `(zero, false)` |
| `GetOrLoad` | `GetOrLoad(ctx context.Context, key string, loader cache.LoaderFunc[T], ttl ...time.Duration) (T, error)` | 先读 cache，miss 时加载；同一 key 的并发加载会合并；`loader` 为 nil 时返回 `ErrLoaderRequired` |
| `Set` | `Set(ctx context.Context, key string, value T, ttl ...time.Duration) error` | 写入值；正数 per-call TTL 会覆盖后端默认 TTL；`Close()` 后返回 `ErrCacheClosed` |
| `Contains` | `Contains(ctx context.Context, key string) bool` | key 不存在、已过期、后端错误或 cache 已关闭时返回 false |
| `Delete` | `Delete(ctx context.Context, key string) error` | 删除单个 key；已关闭 cache 上调用是 nil no-op |
| `Clear` | `Clear(ctx context.Context) error` | 清空该 cache 实例拥有的条目；已关闭 cache 上调用是 nil no-op |
| `Keys` | `Keys(ctx context.Context, prefix ...string) ([]string, error)` | 返回未过期/user-facing key，可按一个 prefix 过滤；已关闭 cache 返回 `nil, nil` |
| `ForEach` | `ForEach(ctx context.Context, callback func(key string, value T) bool, prefix ...string) error` | 遍历未过期/user-facing 条目，可按一个 prefix 过滤；callback 返回 false 时停止 |
| `Size` | `Size(ctx context.Context) (int64, error)` | 返回条目数量；已关闭 cache 返回 `0, nil` |
| `Close` | `Close() error` | 标记 cache 已关闭；内存缓存还会停止后台清理 loop |

两个内置后端的 `GetOrLoad` 都使用 singleflight 行为。它会先在加入 singleflight 前读一次 cache，再在 singleflight 函数内部读一次；同一 key 的并发 miss 只执行一个 loader；loader error 和写入 error 都会原样返回，不会缓存失败结果。

## 内存后端

`cache.NewMemory[T](...)` 是进程内、本地且不持久化的缓存。默认值如下：

| 设置 | 默认值 |
| --- | --- |
| 最大条目数 | `0`（无限制） |
| 默认 TTL | `0`（不过期） |
| 配置的淘汰策略 | `cache.EvictionPolicyLRU` |
| 清理周期 | `5m` |

当有效最大条目数 `<= 0` 时，内存缓存会强制使用 `EvictionPolicyNone`，因为它不需要选择淘汰对象。当最大条目数为正数，而配置的策略是 `EvictionPolicyNone` 或任何不支持的值时，缓存会回退到 `EvictionPolicyLRU`。

`Set` 和 `GetOrLoad` 使用相同 TTL 规则：

- 正数 per-call TTL 优先；
- 否则使用 `WithMemDefaultTTL` 提供的正数默认 TTL；
- `ttl <= 0` 本身不会创建过期时间；
- 如果两者都不是正数，条目不会过期。

过期内存条目会被 `Get`、`Contains`、`Keys` 和 `ForEach` 当成 miss。`Get` 和 `Contains` 会懒删除过期条目；`Keys` 和 `ForEach` 会跳过它们；后台清理 loop 也会删除过期条目。

## Redis 后端

`cache.NewRedis[T](client, namespace, ...)` 要求 `*redis.Client` 非 nil 且 namespace 非空；任一条件不满足都会 panic。Redis client 的生命周期由调用方管理：`Close()` 只会标记 cache 已关闭，不会关闭底层 client。

Redis value 使用 JSON 序列化。`Set` 会返回序列化错误或 Redis 写入错误。`Get` 会把 Redis 读取错误、key 不存在和 JSON 反序列化失败都当成 miss 并记录日志；`ForEach` 在无法读取或反序列化扫描到的 value 时会返回 error。

后端会把 key 存在可观察的 Redis prefix 下：

```text
vef:cache:<namespace>
```

`Keys` 和 `ForEach` 返回前会剥离这个内部 prefix，所以写入的 `user:1` 会以 `user:1` 返回，而不是 Redis 里的存储 key。`Clear` 只删除该 cache namespace 下的 key，`Size` 也只统计该 namespace。`Keys` 和 `ForEach` 的 prefix filter 基于 Redis `SCAN MATCH`；固定 namespace 和用户 prefix 片段会转义 Redis glob 元字符 `*`、`?`、`[`、`]` 和 `\`，因此用户 key 会按字面量匹配。

Redis TTL 规则与内存后端一致：正数 per-call TTL 优先，否则使用正数 `WithRdsDefaultTTL`，否则 Redis key 不设置过期时间。

## Key builders

`cache.Key(parts...)` 等价于空 prefix 的默认 builder：

| 调用 | 结果 |
| --- | --- |
| `cache.Key("user", "123")` | `user:123` |
| `cache.Key()` | `""` |
| `cache.NewPrefixKeyBuilder("app").Build()` | `app` |
| `cache.NewPrefixKeyBuilder("app").Build("user", "123")` | `app:user:123` |
| `cache.NewPrefixKeyBuilderWithSeparator("app", "/").Build("user", "123")` | `app/user/123` |

`cache.KeyBuilder` 有意保持很小：

```go
type KeyBuilder interface {
    Build(keyParts ...string) string
}
```

自定义 cache 实现可以接收自己的 `KeyBuilder`，同时保持和包内 key 语义一致。

## `cache.SingleflightMixin[T]`

`SingleflightMixin[T]` 暴露一个公开方法：

```go
func (m *SingleflightMixin[T]) GetOrLoad(
    ctx context.Context,
    cacheKey string,
    loader cache.LoaderFunc[T],
    ttl []time.Duration,
    getFn cache.GetFunc[T],
    setFn cache.SetFunc[T],
) (value T, err error)
```

编写自定义 `cache.Cache[T]` 实现时可以复用它。该 mixin 会：

- loader 为 nil 时返回 `ErrLoaderRequired`；
- 在 singleflight 协调前调用 `getFn`；
- 在协调后的函数内部再次调用 `getFn`；
- 只有 loader 成功后才调用 `setFn`；
- 直接返回 loader 和 setter error；
- 如果协调结果不能断言为 `T`，返回 `ErrTypeAssertionFailed`。

## `cache.Invalidating[T]`

`cache.NewInvalidating[T](loader, logger, opts...)` 包装的是内部 `cache.NewMemory[T]()` 实例。它接收 `cache.MemoryOption`（例如 `cache.WithMemMaxSize` 或 `cache.WithMemDefaultTTL`），但不接收 Redis option。请传入非 nil 的 `cache.KeyedLoaderFunc[T]` 和非 nil 的 `logx.Logger`；构造器不会在后续方法调用前替你校验它们。

`Invalidating[T]` 暴露：

| 方法 | 签名 | 契约 |
| --- | --- | --- |
| `Get` | `Get(ctx context.Context, key string) (T, error)` | miss 时用 keyed loader 加载 `key`，缓存成功结果，透传 loader error，并合并同一 key 的并发加载 |
| `Invalidate` | `Invalidate(ctx context.Context, keys ...string) error` | `keys` 为空时清空整个缓存；否则只删除指定 key，并记录每次 clear/delete 的结果 |

## 错误类型

| 错误 | 触发条件 |
| --- | --- |
| `cache.ErrMemoryLimitExceeded` | memory `Set` 在正数最大条目数下新增 key，且无法选择淘汰候选 |
| `cache.ErrCacheClosed` | memory 或 Redis 写入路径在 `Close()` 之后被调用，包括 `Set`，或 loader 成功后的 `GetOrLoad` |
| `cache.ErrLoaderRequired` | `Cache.GetOrLoad` 或 `SingleflightMixin.GetOrLoad` 收到 nil loader |
| `cache.ErrTypeAssertionFailed` | `SingleflightMixin.GetOrLoad` 收到的 singleflight 结果不是 `T` |

## 最小示例

```go
package usercache

import (
  "context"
  "time"

  "github.com/coldsmirk/vef-framework-go/cache"
)

type UserCacheService struct {
  users cache.Cache[string]
}

func NewUserCacheService() *UserCacheService {
  return &UserCacheService{
    users: cache.NewMemory[string](
      cache.WithMemDefaultTTL(10*time.Minute),
      cache.WithMemMaxSize(10_000),
      cache.WithMemEvictionPolicy(cache.EvictionPolicyLRU),
    ),
  }
}

func (s *UserCacheService) LoadUserName(ctx context.Context) (string, error) {
  return s.users.GetOrLoad(ctx, "user:1001", func(ctx context.Context) (string, error) {
    return "alice", nil
  }, 5*time.Minute)
}
```

## 内存与 Redis 如何选

| 后端 | 更适合的场景 |
| --- | --- |
| memory | 缓存只需要留在当前进程、服务只有单实例、希望零基础设施依赖 |
| Redis | 多实例共享缓存状态、缓存需要跨进程重启保留、跨节点协调很重要 |

## 这个功能的边界

缓存应在模块或 service 层创建一次；如果你在每个请求 handler 里都新建一个 cache，那它实际上就失去了缓存意义。

## 相关功能

- [配置参考](../reference/configuration-reference)：Redis 配置字段
- [事件总线](./event-bus)：如果你想做缓存失效事件

## 下一步

继续看 [事件总线](./event-bus)，如果你的缓存失效和异步刷新流程要一起工作，就会用到它。
