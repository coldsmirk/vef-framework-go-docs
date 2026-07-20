---
sidebar_position: 10
---

# Distributed Lock

The `lock` package provides lease-based distributed locks for
applications that deploy multiple replicas and need "only one of us does
this" guarantees — singleton cron jobs, one-off migrations, resource-exclusive
maintenance work. Inject `lock.Locker` and the framework hands you the right
implementation for your topology.

## Quick Start

```go
type CleanupJob struct {
    locker lock.Locker
}

// The primary entry point: run fn while holding the lock.
func (j *CleanupJob) Run(ctx context.Context) error {
    return lock.WithLock(ctx, j.locker, "cleanup:orders", func(ctx context.Context) error {
        // Exclusive section. ctx is canceled if the lease is lost.
        return j.cleanupExpiredOrders(ctx)
    })
}
```

For cron jobs where "someone else is already running it" is a normal outcome,
guard with a single non-blocking attempt:

```go
held, err := j.locker.TryAcquire(ctx, "cron:daily-report")
if errors.Is(err, lock.ErrNotAcquired) {
    return nil // another replica won this tick
}
if err != nil {
    return err
}
defer func() { _ = held.Release(context.WithoutCancel(ctx)) }()
```

## Topology-Selected Default

The DI default is selected by deployment topology:

- **Redis enabled** (`vef.redis.enabled = true`): `lock.RedisLocker` — real
  cross-replica mutual exclusion for every node sharing the Redis instance.
- **Redis disabled**: `lock.MemoryLocker` with a loud boot warning — it is
  in-process only and provides **no cross-replica exclusion**.

This deliberately diverges from the framework's usual "memory default, swap
via `fx.Decorate`" convention: applications reach for a distributed lock
precisely because they scale out, and a silently-local lock stops guarding
invariants the moment a second replica starts. Both implementations share
identical semantics (TTL expiry, ownership tokens, fencing tokens, waiting),
so behavior does not change between development and production. Swap in a
custom backend with `fx.Decorate` if needed.

## API

```go
type Locker interface {
    Acquire(ctx context.Context, name string, opts ...Option) (Lock, error)
    TryAcquire(ctx context.Context, name string, opts ...Option) (Lock, error)
}

type Lock interface {
    Release(ctx context.Context) error
    Refresh(ctx context.Context) error
    FencingToken() int64
    Done() <-chan struct{}
}
```

`Acquire` retries until the `WithWait` window is exhausted (no waiting by
default); `TryAcquire` is a single non-blocking attempt. Both return
`lock.ErrNotAcquired` when the lock stays held by someone else, and fail
closed on backend errors. `Release` and `Refresh` return `lock.ErrNotHeld`
once the lease is no longer owned — a signal that mutual exclusion may have
been violated in the meantime. `Done()` closes once the lease is known to be
lost; loss is detected by the auto-renewal watchdog, so without
`WithAutoRenew` the channel never closes.

### Options

| Option | Default | Meaning |
| --- | --- | --- |
| `WithTTL(d)` | `lock.DefaultTTL` (30s) | lease duration; auto-expires this long after acquisition or the last refresh, bounding how long a crashed holder can block others |
| `WithWait(d)` | 0 (no waiting) | how long `Acquire` keeps retrying before giving up with `ErrNotAcquired` |
| `WithRetryInterval(d)` | `lock.DefaultRetryInterval` (100ms) | polling cadence of a waiting `Acquire` |
| `WithAutoRenew(on)` | off for bare `Acquire` / `TryAcquire`; on inside `WithLock` | background watchdog refreshes the lease at TTL/3, so a healthy holder never expires mid-work while a crashed one still frees the lock within one TTL |

Auto-renewal requires a TTL of at least `lock.MinAutoRenewTTL` (30ms);
shorter leases fail the acquisition with `lock.ErrAutoRenewTTLTooShort`.

### `WithLock`

`lock.WithLock(ctx, locker, name, fn, opts...)` is the recommended wrapper:

- acquires with auto-renewal on (unless explicitly disabled), so `fn` may
  safely outlive the TTL;
- cancels `fn`'s context as soon as the lease is lost;
- **always releases** afterwards — even when `fn` panics — on a context that
  survives request cancellation;
- returns the joined error of `fn` and the release: a successful `fn` still
  yields an error when the release reports `lock.ErrNotHeld`, because the
  exclusive section can no longer be trusted to have been exclusive.

## Lease Semantics and Caveats

Locks are **cooperative leases**, not absolute guarantees. Every acquisition
carries a TTL that auto-expires if the holder crashes, and only the holder
(identified by a random ownership token) can release or extend its lease. A
process pause that outlives the TTL — GC, VM freeze, network partition — can
let a second holder in. Guard state that must never be corrupted with one of:

- **Fencing tokens**: `Lock.FencingToken()` returns a monotonically
  increasing sequence (globally monotonic, therefore ordered per lock name).
  Pass it to the protected resource so a delayed writer holding a stale lease
  can be rejected (`WHERE fencing_token < ?`-style checks).
- **Idempotency** of the protected operation.
- **Database constraints** (unique keys, conditional updates).

`RedisLocker` targets a **single Redis instance** (or a cluster where the
lock key hashes to one shard). It is intentionally *not* a Redlock
implementation — quorum locking over independent Redis nodes is out of scope.

### Redis implementation notes

Acquisition is one atomic Lua script: it installs the ownership token and
allocates the fencing token (from the single persistent counter
`vef:lock:fencing`) in the expiring lock hash `vef:lock:key:<name>`. Release
and refresh are token-guarded Lua scripts, so a slow holder whose lease
expired can never delete or extend a successor's lock. All three operations
are idempotent across go-redis's internal retries: a replayed acquisition
returns the original fencing token instead of a false conflict, and a
replayed release reports success through a short-lived acknowledgement key
without touching a successor's lock.

## Errors

| Error | Meaning |
| --- | --- |
| `lock.ErrNotAcquired` | the lock is held by someone else (and the wait window, if any, ran out) |
| `lock.ErrNotHeld` | release/refresh on a lease that expired, was already released, or was taken over |
| `lock.ErrAutoRenewTTLTooShort` | auto-renewal requested with a TTL below `lock.MinAutoRenewTTL` |

---

Related: [Cache](./cache) for the Redis client configuration, and
[Sequence](./sequence) for monotonic number allocation (which uses its own
storage-level coordination, not this lock).
