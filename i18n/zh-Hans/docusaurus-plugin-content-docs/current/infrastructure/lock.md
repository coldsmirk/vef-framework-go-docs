---
sidebar_position: 10
---

# 分布式锁

`lock` 包提供基于租约（lease）的分布式锁，服务于多副本部署下
"只允许我们中的一个来做"的需求——单例 cron 任务、一次性迁移、独占资源的
维护操作。注入 `lock.Locker`，框架会按部署拓扑给你正确的实现。

## 快速开始

```go
type CleanupJob struct {
    locker lock.Locker
}

// 首选入口：持锁运行 fn。
func (j *CleanupJob) Run(ctx context.Context) error {
    return lock.WithLock(ctx, j.locker, "cleanup:orders", func(ctx context.Context) error {
        // 独占区。租约丢失时 ctx 会被取消。
        return j.cleanupExpiredOrders(ctx)
    })
}
```

对于"别的副本已经在跑"属于正常结果的 cron 任务，用一次非阻塞尝试做守卫：

```go
held, err := j.locker.TryAcquire(ctx, "cron:daily-report")
if errors.Is(err, lock.ErrNotAcquired) {
    return nil // 这一轮被其他副本抢到了
}
if err != nil {
    return err
}
defer func() { _ = held.Release(context.WithoutCancel(ctx)) }()
```

## 按拓扑选择的默认实现

DI 默认实现由部署拓扑决定：

- **Redis 已启用**（`vef.redis.enabled = true`）：`lock.RedisLocker`——共享
  同一 Redis 的所有节点获得真正的跨副本互斥。
- **Redis 未启用**：`lock.MemoryLocker`，并在启动时打出醒目警告——它只在
  进程内生效，**没有跨副本互斥**。

这有意偏离了框架惯常的"默认 memory、用 `fx.Decorate` 换"约定：应用之所以
需要分布式锁，正是因为要横向扩容，而一把静默退化为本地的锁会在第二个副本
启动的那一刻停止守护不变量。两个实现的语义完全一致（TTL 过期、所有权
token、fencing token、等待），开发与生产环境行为不会漂移。需要自定义后端时
仍可用 `fx.Decorate` 替换。

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

`Acquire` 会持续重试直到 `WithWait` 窗口耗尽（默认不等待）；`TryAcquire`
只做一次非阻塞尝试。锁被他人持有时二者都返回 `lock.ErrNotAcquired`，后端
出错时 fail closed。租约不再归属自己后，`Release` 和 `Refresh` 返回
`lock.ErrNotHeld`——这是"互斥可能已被打破"的信号。`Done()` 在确认租约丢失
后关闭；丢失由自动续约 watchdog 检测，所以不开 `WithAutoRenew` 时该通道
永远不会关闭。

### Options

| Option | 默认值 | 含义 |
| --- | --- | --- |
| `WithTTL(d)` | `lock.DefaultTTL`（30s） | 租约时长；获取或最近一次续约后经过该时长自动过期，限定崩溃的持有者最多阻塞他人多久 |
| `WithWait(d)` | 0（不等待） | `Acquire` 放弃并返回 `ErrNotAcquired` 前的重试窗口 |
| `WithRetryInterval(d)` | `lock.DefaultRetryInterval`（100ms） | 等待型 `Acquire` 的轮询间隔 |
| `WithAutoRenew(on)` | 裸 `Acquire` / `TryAcquire` 默认关；`WithLock` 内默认开 | 后台 watchdog 每 TTL/3 续约一次，健康的持有者不会在工作中途过期，崩溃的持有者仍在一个 TTL 内释放锁 |

自动续约要求 TTL 不低于 `lock.MinAutoRenewTTL`（30ms）；更短的租约在获取时
就以 `lock.ErrAutoRenewTTLTooShort` 失败。

### `WithLock`

`lock.WithLock(ctx, locker, name, fn, opts...)` 是推荐的封装：

- 以自动续约获取（除非显式关闭），`fn` 可以安全地运行超过 TTL；
- 租约一旦丢失立即取消 `fn` 的 context；
- 之后**总是释放**——即使 `fn` panic——且释放使用一个不受请求取消影响的
  context；
- 返回 `fn` 与释放的 join error：即便 `fn` 成功，释放报出
  `lock.ErrNotHeld` 时整体仍是错误，因为这段"独占区"已不可信。

## 租约语义与注意事项

锁是**协作式租约**，不是绝对保证。每次获取都带 TTL，持有者崩溃后自动过期；
只有持有者（由随机所有权 token 标识）能释放或续约。超过 TTL 的进程停顿——
GC、虚拟机冻结、网络分区——可能放进第二个持有者。绝不能被破坏的状态要用
以下手段之一兜底：

- **Fencing token**：`Lock.FencingToken()` 返回单调递增的序号（全局单调，
  因此对每个锁名也有序）。把它传给受保护资源，让拿着过期租约的迟到写入者
  被拒绝（`WHERE fencing_token < ?` 一类的检查）。
- 受保护操作本身的**幂等性**。
- **数据库约束**（唯一键、条件更新）。

`RedisLocker` 面向**单 Redis 实例**（或锁键落在同一分片的集群）。它有意
*不是* Redlock 实现——跨独立 Redis 节点的 quorum 锁不在范围内。

### Redis 实现说明

获取是一段原子 Lua 脚本：在带过期的锁 hash `vef:lock:key:<name>` 里写入
所有权 token，并从唯一的持久计数器 `vef:lock:fencing` 分配 fencing token。
释放与续约是 token 守卫的 Lua 脚本，租约已过期的慢持有者永远删不掉、也
续不了继任者的锁。三种操作对 go-redis 内部重试都幂等：重放的获取返回原来
的 fencing token 而不是假冲突，重放的释放通过短存活的确认键报告成功，不会
碰继任者的锁。

## 错误

| 错误 | 含义 |
| --- | --- |
| `lock.ErrNotAcquired` | 锁被他人持有（且等待窗口已耗尽，如果有的话） |
| `lock.ErrNotHeld` | 对已过期、已释放或已被接管的租约执行释放/续约 |
| `lock.ErrAutoRenewTTLTooShort` | 开启自动续约但 TTL 低于 `lock.MinAutoRenewTTL` |

---

相关：[缓存](./cache) 介绍 Redis client 配置；[序列号](./sequence) 介绍单调
号段分配（它使用自己的存储级协调，不依赖这把锁）。
