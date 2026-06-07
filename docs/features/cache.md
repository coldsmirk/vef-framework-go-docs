---
sidebar_position: 1
---

# Cache

VEF exposes `github.com/coldsmirk/vef-framework-go/cache` as a typed utility package. It does not create one hidden global cache, and it does not inject cache instances into resources automatically.

## Public Surface

The package has no exported fields. Its public top-level API is:

| API | Contract |
| --- | --- |
| `cache.Cache[T]` | generic cache interface implemented by memory and Redis backends |
| `cache.NewMemory[T](opts ...cache.MemoryOption) cache.Cache[T]` | creates an in-process memory cache |
| `cache.NewRedis[T](client *redis.Client, namespace string, opts ...cache.RedisOption) cache.Cache[T]` | creates a Redis-backed cache scoped to `namespace`; panics if `client` is nil or `namespace` is empty |
| `cache.MemoryOption` | functional option type accepted by `NewMemory` |
| `cache.RedisOption` | functional option type accepted by `NewRedis` |
| `cache.WithMemMaxSize(size int64)` | caps the memory cache by entry count; `size <= 0` means unlimited |
| `cache.WithMemDefaultTTL(ttl time.Duration)` | default memory TTL used when `Set` or `GetOrLoad` does not receive a positive per-call TTL |
| `cache.WithMemEvictionPolicy(policy cache.EvictionPolicy)` | selects memory eviction policy while a positive max entry count is enforced |
| `cache.WithMemGCInterval(interval time.Duration)` | controls memory expired-entry cleanup; `interval <= 0` falls back to `5m` |
| `cache.WithRdsDefaultTTL(ttl time.Duration)` | default Redis TTL used when a write does not receive a positive per-call TTL |
| `cache.EvictionPolicy` | memory eviction policy enum |
| `cache.EvictionPolicyNone` | no eviction tracking; used automatically for unlimited memory caches |
| `cache.EvictionPolicyLRU` | evicts the least recently used entry; reads and updates refresh recency |
| `cache.EvictionPolicyLFU` | evicts the least frequently used entry; ties are resolved by insertion order within that frequency |
| `cache.EvictionPolicyFIFO` | evicts the oldest inserted entry; reads do not affect order |
| `cache.Key(keyParts ...string) string` | joins parts with `:` using the default key builder; `cache.Key()` returns `""` |
| `cache.KeyBuilder` | interface with `Build(keyParts ...string) string` |
| `cache.PrefixKeyBuilder` | prefix-based implementation of `KeyBuilder` |
| `cache.NewPrefixKeyBuilder(prefix string) *cache.PrefixKeyBuilder` | builds a key builder with `:` as the separator |
| `cache.LoaderFunc[T]` | `func(ctx context.Context) (T, error)` used by `Cache.GetOrLoad` and `SingleflightMixin.GetOrLoad` |
| `cache.KeyedLoaderFunc[T]` | `func(ctx context.Context, key string) (T, error)` used by `Invalidating` |
| `cache.GetFunc[T]` | `func(context.Context, string) (T, bool)` callback used by `SingleflightMixin` |
| `cache.SetFunc[T]` | `func(context.Context, string, T, ...time.Duration) error` callback used by `SingleflightMixin` |
| `cache.SingleflightMixin[T]` | reusable `GetOrLoad` implementation for custom cache implementations |
| `cache.Invalidating[T]` | read-through memory cache wrapper with explicit invalidation |
| `cache.NewInvalidating[T](loader cache.KeyedLoaderFunc[T], logger logx.Logger) *cache.Invalidating[T]` | creates an invalidating cache; callers must pass a non-nil loader and logger |
| `cache.ErrMemoryLimitExceeded` | memory cache hit a positive entry limit and could not find an eviction candidate |
| `cache.ErrCacheClosed` | write path attempted after `Close()` |
| `cache.ErrLoaderRequired` | `GetOrLoad` received a nil loader |
| `cache.ErrTypeAssertionFailed` | `SingleflightMixin` received a result of the wrong generic type |

## `cache.Cache[T]`

Both built-in constructors return `cache.Cache[T]`, which embeds `io.Closer` and exposes these methods:

| Method | Signature | Contract |
| --- | --- | --- |
| `Get` | `Get(ctx context.Context, key string) (T, bool)` | returns `(zero, false)` on miss, expired entry, backend read failure, deserialization failure, or closed cache |
| `GetOrLoad` | `GetOrLoad(ctx context.Context, key string, loader cache.LoaderFunc[T], ttl ...time.Duration) (T, error)` | reads first, loads on miss, suppresses concurrent loads for the same key, and returns `ErrLoaderRequired` when `loader` is nil |
| `Set` | `Set(ctx context.Context, key string, value T, ttl ...time.Duration) error` | stores a value; a positive per-call TTL overrides the backend default TTL; returns `ErrCacheClosed` after `Close()` |
| `Contains` | `Contains(ctx context.Context, key string) bool` | returns false for missing, expired, backend-error, or closed entries |
| `Delete` | `Delete(ctx context.Context, key string) error` | removes one key; deleting from a closed cache is a nil no-op |
| `Clear` | `Clear(ctx context.Context) error` | clears all entries owned by that cache instance; clearing a closed cache is a nil no-op |
| `Keys` | `Keys(ctx context.Context, prefix ...string) ([]string, error)` | returns unexpired/user-facing keys, optionally filtered by one prefix; a closed cache returns `nil, nil` |
| `ForEach` | `ForEach(ctx context.Context, callback func(key string, value T) bool, prefix ...string) error` | iterates unexpired/user-facing entries, optionally filtered by one prefix, and stops when the callback returns false |
| `Size` | `Size(ctx context.Context) (int64, error)` | returns entry count; a closed cache returns `0, nil` |
| `Close` | `Close() error` | marks the cache closed; memory cache also stops its background cleanup loop |

`GetOrLoad` uses singleflight behavior in both built-in backends. It checks the cache before joining the singleflight call, checks it again inside the singleflight function, runs one loader for concurrent misses of the same key, propagates loader errors, and propagates write errors instead of caching a failed load.

## Memory Backend

`cache.NewMemory[T](...)` is process-local and non-durable. Defaults are:

| Setting | Default |
| --- | --- |
| max entry count | `0` (unlimited) |
| default TTL | `0` (no expiration) |
| configured eviction policy | `cache.EvictionPolicyLRU` |
| cleanup interval | `5m` |

When the effective max entry count is `<= 0`, the memory cache forces `EvictionPolicyNone` because it never needs to choose victims. When the max entry count is positive and the configured policy is `EvictionPolicyNone` or any unsupported value, the cache falls back to `EvictionPolicyLRU`.

TTL rules are the same for `Set` and `GetOrLoad`:

- a positive per-call TTL wins;
- otherwise a positive default TTL from `WithMemDefaultTTL` applies;
- `ttl <= 0` does not create an expiration by itself;
- if neither source is positive, the entry does not expire.

Expired memory entries are treated as misses by `Get`, `Contains`, `Keys`, and `ForEach`. `Get` and `Contains` remove expired entries lazily; `Keys` and `ForEach` skip them; the background cleanup loop also removes expired entries.

## Redis Backend

`cache.NewRedis[T](client, namespace, ...)` requires a non-nil `*redis.Client` and a non-empty namespace; either violation panics. The Redis client is owned by the caller: `Close()` only marks the cache closed and does not close the underlying client.

Redis values are serialized with JSON. `Set` returns serialization or Redis write errors. `Get` treats Redis read errors, missing keys, and JSON deserialization failures as misses and logs failures; `ForEach` returns an error when it cannot read or deserialize a scanned value.

The backend stores keys under the observable Redis prefix:

```text
vef:cache:<namespace>
```

`Keys` and `ForEach` strip that internal prefix before returning keys to the caller, so a stored key such as `user:1` is returned as `user:1`, not as the Redis storage key. `Clear` deletes only keys under the cache namespace, and `Size` counts only that namespace. `Keys` and `ForEach` prefix filters use Redis `SCAN MATCH`; fixed namespace and user-prefix pieces escape Redis glob metacharacters `*`, `?`, `[`, `]`, and `\` so user keys are matched literally.

Redis TTL rules match the memory backend: a positive per-call TTL wins, otherwise a positive `WithRdsDefaultTTL` applies, otherwise Redis stores the key without expiration.

## Key Builders

`cache.Key(parts...)` is equivalent to the default prefix builder with an empty prefix:

| Call | Result |
| --- | --- |
| `cache.Key("user", "123")` | `user:123` |
| `cache.Key()` | `""` |
| `cache.NewPrefixKeyBuilder("app").Build()` | `app` |
| `cache.NewPrefixKeyBuilder("app").Build("user", "123")` | `app:user:123` |

`cache.KeyBuilder` is intentionally small:

```go
type KeyBuilder interface {
    Build(keyParts ...string) string
}
```

Custom cache implementations can accept their own `KeyBuilder` while still matching the package's key semantics.

## `cache.SingleflightMixin[T]`

`SingleflightMixin[T]` exposes one public method:

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

Use it when writing a custom `cache.Cache[T]` implementation. The mixin:

- returns `ErrLoaderRequired` for a nil loader;
- calls `getFn` before singleflight coordination;
- calls `getFn` again inside the coordinated function;
- calls `setFn` only after a successful loader result;
- returns loader and setter errors directly;
- returns `ErrTypeAssertionFailed` if the coordinated result cannot be asserted to `T`.

## `cache.Invalidating[T]`

`cache.NewInvalidating[T](loader, logger)` wraps an internal `cache.NewMemory[T]()` instance. It does not accept TTL or Redis options. Pass a non-nil `cache.KeyedLoaderFunc[T]` and a non-nil `logx.Logger`; the constructor does not validate them before later method calls.

`Invalidating[T]` exposes:

| Method | Signature | Contract |
| --- | --- | --- |
| `Get` | `Get(ctx context.Context, key string) (T, error)` | loads `key` with the keyed loader on miss, caches successful results, propagates loader errors, and merges concurrent loads for the same key |
| `Invalidate` | `Invalidate(ctx context.Context, keys ...string) error` | clears the whole cache when `keys` is empty; otherwise deletes exactly the named keys and logs each clear/delete result |

## Error Types

| Error | Trigger |
| --- | --- |
| `cache.ErrMemoryLimitExceeded` | memory `Set` is adding a new key under a positive max entry count and no eviction candidate can be selected |
| `cache.ErrCacheClosed` | memory or Redis write path is called after `Close()`, including `Set` or `GetOrLoad` after a successful loader result |
| `cache.ErrLoaderRequired` | `Cache.GetOrLoad` or `SingleflightMixin.GetOrLoad` is called with nil loader |
| `cache.ErrTypeAssertionFailed` | `SingleflightMixin.GetOrLoad` receives a singleflight result that is not `T` |

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

Create cache instances once at module or service scope. Creating a new cache inside every request handler defeats the point of caching.

## Related Features

- [Configuration Reference](../reference/configuration-reference) for Redis config fields
- [Event Bus](./event-bus) if you want cache invalidation events

## Next Step

Read [Event Bus](./event-bus) if cache invalidation and async refresh flows should work together.
