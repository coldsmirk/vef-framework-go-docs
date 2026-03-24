---
sidebar_position: 1
---

# 缓存

VEF 暴露的是一个类型化缓存包，而不是一个隐藏的全局缓存单例。

## 构造器

当前公共缓存包提供以下构造器：

| 构造器 | 后端 | 说明 |
| --- | --- | --- |
| `cache.NewMemory[T](opts...)` | 进程内内存缓存 | 不依赖外部基础设施 |
| `cache.NewRedis[T](client, namespace, opts...)` | Redis 缓存 | 需要非空 Redis client 和非空 namespace |

这两个构造器都会返回 `cache.Cache[T]`。

## `cache.Cache[T]` 接口

完整缓存接口包括：

| 方法 | 作用 |
| --- | --- |
| `Get(ctx, key)` | 读取单个值 |
| `GetOrLoad(ctx, key, loader, ttl...)` | 读取或在 miss 时加载 |
| `Set(ctx, key, value, ttl...)` | 写入单个值 |
| `Contains(ctx, key)` | 判断 key 是否存在 |
| `Delete(ctx, key)` | 删除单个 key |
| `Clear(ctx)` | 清空所有条目 |
| `Keys(ctx, prefix...)` | 枚举 key，可按前缀过滤 |
| `ForEach(ctx, callback, prefix...)` | 遍历条目，可按前缀过滤 |
| `Size(ctx)` | 返回条目数量 |
| `Close()` | 释放资源 |

## `GetOrLoad` 语义

`GetOrLoad` 是请求驱动代码里最实用的方法，因为实现会保证：对同一个 key 的并发 miss 只会真正执行一次 loader。

这意味着你可以直接得到：

- cache miss 自动加载
- 热 key 并发去重
- 统一的 read-through cache 使用方式

## 内存缓存配置项

`cache.NewMemory[T](...)` 支持以下 option：

| Option | 作用 |
| --- | --- |
| `cache.WithMemMaxSize(size)` | 设置最大条目数或容量上限；`<= 0` 表示禁用该限制 |
| `cache.WithMemDefaultTTL(ttl)` | 设置默认 TTL |
| `cache.WithMemEvictionPolicy(policy)` | 在达到上限时选择淘汰策略 |
| `cache.WithMemGCInterval(interval)` | 设置过期条目清理周期 |

支持的内存淘汰策略：

| 策略 | 含义 |
| --- | --- |
| `cache.EvictionPolicyNone` | 不做淘汰跟踪 |
| `cache.EvictionPolicyLRU` | 最近最少使用 |
| `cache.EvictionPolicyLFU` | 最少使用频次 |
| `cache.EvictionPolicyFIFO` | 先进先出 |

## Redis 缓存配置项

`cache.NewRedis[T](client, namespace, ...)` 的要求和 option 如下：

| 要求或 option | 含义 |
| --- | --- |
| 非空 `client` | 必需 |
| 非空 `namespace` | 必需 |
| `cache.WithRdsDefaultTTL(ttl)` | 设置默认 TTL |

Redis 缓存会在内部构造带前缀的 key，从而保证不同 namespace 之间相互隔离。

## Store 级抽象

缓存包还暴露了一个更底层的 `cache.Store` 接口，用于原始字节级存储后端。

如果你要在同一套高层缓存模型下实现自定义后端，这个接口会很有用。

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

缓存包是一个公共工具包，不会自动注入到你的资源里。缓存实例由应用自己决定在哪里构造、生命周期有多长。

缓存应在模块或 service 层创建一次；如果你在每个请求 handler 里都新建一个 cache，那它实际上就失去了缓存意义。

## 相关功能

- [配置参考](../reference/configuration-reference)：Redis 配置字段
- [事件总线](./event-bus)：如果你想做缓存失效事件

## 下一步

继续看 [事件总线](./event-bus)，如果你的缓存失效和异步刷新流程要一起工作，就会用到它。
