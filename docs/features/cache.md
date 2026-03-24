---
sidebar_position: 1
---

# Cache

VEF exposes a typed cache package instead of one hidden global cache singleton.

## Constructors

The public package currently provides these cache constructors:

| Constructor | Backend | Notes |
| --- | --- | --- |
| `cache.NewMemory[T](opts...)` | in-process memory cache | no external dependency |
| `cache.NewRedis[T](client, namespace, opts...)` | Redis-backed cache | requires a non-nil Redis client and a non-empty namespace |

Both constructors return `cache.Cache[T]`.

## `cache.Cache[T]` Interface

The full cache interface includes:

| Method | Purpose |
| --- | --- |
| `Get(ctx, key)` | fetch one value |
| `GetOrLoad(ctx, key, loader, ttl...)` | fetch or compute-on-miss |
| `Set(ctx, key, value, ttl...)` | store one value |
| `Contains(ctx, key)` | existence check |
| `Delete(ctx, key)` | remove one key |
| `Clear(ctx)` | clear all entries |
| `Keys(ctx, prefix...)` | list keys, optionally by prefix |
| `ForEach(ctx, callback, prefix...)` | iterate entries, optionally by prefix |
| `Size(ctx)` | return entry count |
| `Close()` | release resources |

## `GetOrLoad` Semantics

`GetOrLoad` is the most practical method in request-driven code because implementations ensure that concurrent calls for the same key only execute one loader.

That gives you:

- cache-miss loading
- duplicate-load suppression
- a single API for “read-through cache” behavior

## Memory Cache Options

`cache.NewMemory[T](...)` supports these options:

| Option | Purpose |
| --- | --- |
| `cache.WithMemMaxSize(size)` | sets maximum entry count or size guard; values `<= 0` disable the limit |
| `cache.WithMemDefaultTTL(ttl)` | sets fallback TTL |
| `cache.WithMemEvictionPolicy(policy)` | selects the eviction policy when max size is enforced |
| `cache.WithMemGCInterval(interval)` | controls how often expired-entry cleanup runs |

Supported memory eviction policies:

| Policy | Meaning |
| --- | --- |
| `cache.EvictionPolicyNone` | no eviction tracking |
| `cache.EvictionPolicyLRU` | least recently used |
| `cache.EvictionPolicyLFU` | least frequently used |
| `cache.EvictionPolicyFIFO` | first in, first out |

## Redis Cache Options

`cache.NewRedis[T](client, namespace, ...)` supports:

| Requirement or option | Meaning |
| --- | --- |
| non-nil `client` | required |
| non-empty `namespace` | required |
| `cache.WithRdsDefaultTTL(ttl)` | sets fallback TTL |

Redis caches build prefixed keys internally so namespaces remain isolated.

## Store-Level Abstraction

The package also exposes a lower-level `cache.Store` interface for raw byte storage backends.

This interface is useful if you want to implement a custom cache backend under the same higher-level cache model.

## Minimal Example

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

## Memory vs Redis

| Backend | Prefer when |
| --- | --- |
| memory | cache can stay process-local, you run one instance, and you want zero infrastructure dependency |
| Redis | multiple instances need shared cache state, cache must survive process restarts, or cross-node coordination matters |

## Scope Of This Feature

The cache package is a public utility package. It is not automatically injected into resources for you. You decide where to construct cache instances and how long they live.

Create cache instances once at module or service scope. Creating a new cache inside every request handler defeats the point of caching.

## Related Features

- [Configuration Reference](../reference/configuration-reference) for Redis config fields
- [Event Bus](./event-bus) if you want cache invalidation events

## Next Step

Read [Event Bus](./event-bus) if cache invalidation and async refresh flows should work together.
