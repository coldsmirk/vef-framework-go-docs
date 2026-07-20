---
sidebar_position: 4
---

# Sequence

`sequence` 包用于生成序列号（订单号、发票号等），支持自定义前缀、日期段、零填充计数器和自动重置策略。

## 概念

一个序列号生成器由两部分组成：

- **`sequence.Rule`** —— 格式与重置策略（前缀、日期格式、计数器位宽、重置周期、溢出策略）。
- **`sequence.Store`** —— rule 与当前计数器的存储位置。默认运行时 store 是内存实现；公开包同时提供数据库和 Redis store，适合持久化或分布式部署。如果需要其他持久化模型，也可以实现同一个接口。Store 负责原子地递增计数器。

框架已经把 `sequence.Generator` 装配好了，也会暴露具体的 `*sequence.MemoryStore`，所以业务模块可以在启动期 seed rule，然后调用 `Generate(ctx, key)`。

辅助函数 `sequence.FormatDate(dt, format)` 也是公开 API。它会使用
`Rule.DateFormat` 相同的 `yyyy` / `MM` / `dd` / `HH` / `mm` / `ss` token
渲染日期片段。

## 定义规则

```go
import (
    "github.com/coldsmirk/vef-framework-go/sequence"
)

rule := &sequence.Rule{
    Key:              "order-number",     // 调用 Generate(ctx, key) 时的查询键
    Name:             "订单号",
    Prefix:           "ORD-",
    DateFormat:       "yyyyMMdd-",        // 可选；使用下方的日期 token
    SeqLength:        6,                  // 零填充到 6 位 → 000001
    SeqStep:          1,
    StartValue:       0,                  // 第一个生成的值 = StartValue + SeqStep
    MaxValue:         0,                  // 0 表示不设上限
    OverflowStrategy: sequence.OverflowError,
    ResetCycle:       sequence.ResetDaily,
    IsActive:         true,
}
```

### Rule 字段

| 字段 | 含义 |
| --- | --- |
| `Key` | 调用 `Generate(ctx, key)` 时的查询键。在 store 内唯一。 |
| `Name` | 可读名称（用于管理后台展示）。 |
| `Prefix` / `Suffix` | 日期+计数器前后的固定文字，可选。 |
| `DateFormat` | 可选的日期布局（token 见下）；为空表示不带日期。 |
| `SeqLength` | 零填充后的计数器位宽。 |
| `SeqStep` | 每次生成时的递增步长（通常为 1）。 |
| `StartValue` | 重置后的计数器起点。第一个生成的值 = `StartValue + SeqStep`。 |
| `MaxValue` | 上限。`0` 表示不限制。 |
| `OverflowStrategy` | 达到 `MaxValue` 后的行为，见下。 |
| `ResetCycle` | 何时重置计数器，见下。 |
| `CurrentValue` | 上一次预留后的计数器值。store 会更新这个游标；业务调用方通常不直接修改它。 |
| `LastResetAt` | 用来判断下一次预留是否跨过 reset 边界的时间戳。`nil` 表示该 rule 从未 reset。 |
| `IsActive` | 非活动的 rule 会让 `Generate` 返回 `sequence.ErrRuleNotFound`。 |

`Rule.Clone()` 会返回 rule 快照的深拷贝；如果 `LastResetAt` 不为 nil，
指针指向的时间值也会被复制。

### 日期布局 token

`DateFormat` 使用类 Java / .NET 风格的 token（内部会翻译为 Go layout）：

| Token | 含义 | 示例 |
| --- | --- | --- |
| `yyyy` | 4 位年份 | `2024` |
| `yy` | 2 位年份 | `24` |
| `MM` | 2 位月份 | `03` |
| `dd` | 2 位日期 | `15` |
| `HH` | 2 位小时（24 制） | `14` |
| `mm` | 2 位分钟 | `30` |
| `ss` | 2 位秒 | `05` |

其他字符原样保留，例如 `yyyyMMdd-` 生成 `20240315-`。

### 重置周期

| 常量 | 周期 |
| --- | --- |
| `sequence.ResetNone` | 永不重置 |
| `sequence.ResetDaily` | 每日凌晨重置 |
| `sequence.ResetWeekly` | 每周开始重置 |
| `sequence.ResetMonthly` | 每月 1 日重置 |
| `sequence.ResetQuarterly` | 每季度首日重置 |
| `sequence.ResetYearly` | 每年 1 月 1 日重置 |

空 `ResetCycle` 等同于 `sequence.ResetNone`。除公开常量外的其他值会按防御性
fallback 处理为“不重置”。对非 none 周期，`LastResetAt == nil` 会让下一次预
留触发 reset。

### 溢出策略

| 常量 | 达到 `MaxValue` 后的行为 |
| --- | --- |
| `sequence.OverflowError`（默认） | 返回 `sequence.ErrSequenceOverflow`，直到下次重置都拒绝生成。 |
| `sequence.OverflowReset` | 把计数器重置到 `StartValue` 后继续。 |
| `sequence.OverflowExtend` | 不受 `SeqLength` 限制继续递增（结果会变长）。 |

`MaxValue` 会在应用周期 reset 后再检查。如果 reset 边界已经把计数器回到
`StartValue`，但本次批量预留仍然超过 `MaxValue`，即便
`OverflowReset` 也会返回 `sequence.ErrSequenceOverflow`，因为再次 reset 也
无法让这批值落进上限。除公开 `OverflowStrategy` 常量外的其他值会回退为
`OverflowError`。

## Store

当前内置运行时 store 是内存实现，且不持久化：进程重启后计数器和已注册
rule 都会丢失。它适合测试、开发和单进程部署；需要分布式或持久化计数时，
应提供自定义 `sequence.Store`。

`sequence.Store.Reserve(ctx, key, count, now)` 是自定义 store 的契约边界。
实现必须按 rule key 序列化 read-modify-write 路径，并用一次原子操作预留整
个 `count` 批次。

### 内存

适用于测试、开发、单进程部署。Rule 必须提前 `Register`：

```go
store := sequence.NewMemoryStore()
store.Register(rule)
```

`Register` 会覆盖相同 `Key` 的已有 rule，并存入深拷贝；之后再修改原来的
`*Rule` 不会影响 store 内部状态。

在 VEF 应用中，如果需要 seed rule，可以注入具体的 `*sequence.MemoryStore`：

```go
func SeedSequenceRules(store *sequence.MemoryStore) {
    store.Register(rule)
}
```

### 数据库

`sequence.NewDBStore(db)` 返回基于 `sys_sequence_rule` 表的
`*sequence.DBStore`，表名常量是 `sequence.DBStoreTableName`。`DBStore.Init(ctx)`
会在表不存在时创建它；`Reserve(...)` 会对 rule 行加锁并在数据库事务内原子预留计数。
`sequence.RuleModel` 是该表对应的 ORM model。

### Redis

`sequence.NewRedisStore(client)` 返回适合分布式部署的 `*sequence.RedisStore`。
规则以 Redis hash 存在 `vef:sequence:<key>` 前缀下；`RedisStore.RegisterRule(ctx, rule)`
用于 seed 或替换一条规则，`Reserve(...)` 使用 Redis `WATCH` / transaction retry
原子预留计数。

每个 store 都**要求 rule 先存在再调用 `Generate(...)`**。调一个 store 里没有的 key 会得到 `sequence.ErrRuleNotFound`。

公开的 store API 有意保持很小：

| API | 作用 |
| --- | --- |
| `sequence.Store.Reserve(ctx, key, count, now)` | 为某个 rule 原子预留 `count` 个序号，并返回 rule 快照与本批最终计数值 |
| `sequence.MemoryStore.Register(rules...)` | 使用深拷贝预加载或替换内存规则 |
| `sequence.MemoryStore.Reserve(...)` | `Store.Reserve` 的内存实现；返回克隆后的 rule 快照 |
| `sequence.NewDBStore(db)` / `sequence.DBStore` | 基于 `sequence.DBStoreTableName`（`sys_sequence_rule`）和 `sequence.RuleModel` 的数据库 store |
| `sequence.NewRedisStore(client)` / `sequence.RedisStore` | Redis store；用 `RedisStore.RegisterRule(ctx, rule)` seed hash-backed 规则 |
| `sequence.Rule.Clone()` | 深拷贝一个 rule 快照 |

## 生成序列号

框架已经基于当前生效的 `Store` 装配好 `sequence.Generator`，业务代码直接注入并调用：

```go
type OrderService struct {
    seq sequence.Generator
}

func (s *OrderService) NewOrder(ctx context.Context) (string, error) {
    return s.seq.Generate(ctx, "order-number")
}
```

需要原子地一次取多个时（例如 100 个）：

```go
numbers, err := seq.GenerateN(ctx, "order-number", 100)
```

`GenerateN` 用一次原子操作预定整个区间。返回值按从小到大的预留顺序排列；
当 `SeqStep > 1` 时，值会按该步长间隔递增（例如 `0002`、`0004`、
`0006`）。

`SeqLength` 是最小零填充位宽，不是最大长度。数值位数超过 `SeqLength` 时会
完整输出，不会被截断。

## 示例规则

| 场景 | 配置 | 输出样例 |
| --- | --- | --- |
| 订单号 | `Prefix:"ORD-"`、`DateFormat:"yyyyMMdd-"`、`SeqLength:6`、`ResetDaily` | `ORD-20240315-000001` |
| 发票号 | `Prefix:"INV"`、`DateFormat:"yyyyMM-"`、`SeqLength:4`、`ResetMonthly` | `INV202403-0001` |
| 文档 ID | `Prefix:"DOC-"`、`DateFormat:"yyyy-"`、`SeqLength:8`、`ResetYearly` | `DOC-2024-00000001` |
| 普通计数器 | `SeqLength:10`、`ResetNone` | `0000000001` |

## 错误

| 错误 | 触发 |
| --- | --- |
| `sequence.ErrRuleNotFound` | store 里没有这个 key，或 rule 处于 `IsActive=false`。 |
| `sequence.ErrSequenceOverflow` | 达到 `MaxValue` 且 `OverflowStrategy` 为 `OverflowError`。 |
| `sequence.ErrInvalidCount` | 调用 `GenerateN` 时传入了小于 1 的数量。 |
| `sequence.ErrInvalidStep` | rule 的 `SeqStep` 小于 1。 |

## 下一步

参考 [缓存](./cache) 或 [定时任务](./cron)——序列生成器经常与定时归档或预热任务配套使用。
