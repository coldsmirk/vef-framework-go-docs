---
sidebar_position: 14
---

# Sequence

`sequence` 包用于生成序列号（订单号、发票号等），支持自定义前缀、日期段、零填充计数器和自动重置策略。

## 概念

一个序列号生成器由两部分组成：

- **`sequence.Rule`** —— 格式与重置策略（前缀、日期格式、计数器位宽、重置周期、溢出策略）。
- **`sequence.Store`** —— rule 与当前计数器的存储位置（数据库、Redis 或内存）。Store 负责原子地递增计数器。

框架已经把 `sequence.Generator` 装配好了；业务代码只需要把 rule 注册到 store，然后调用 `Generate(ctx, key)`。

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
| `IsActive` | 非活动的 rule 会让 `Generate` 返回 `sequence.ErrRuleNotFound`。 |

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

### 溢出策略

| 常量 | 达到 `MaxValue` 后的行为 |
| --- | --- |
| `sequence.OverflowError`（默认） | 返回 `sequence.ErrSequenceOverflow`，直到下次重置都拒绝生成。 |
| `sequence.OverflowReset` | 把计数器重置到 `StartValue` 后继续。 |
| `sequence.OverflowExtend` | 不受 `SeqLength` 限制继续递增（结果会变长）。 |

## Store

按部署形态选择：

### 内存

适用于测试、开发、单进程部署。Rule 必须提前 `Register`：

```go
store := sequence.NewMemoryStore().(*sequence.MemoryStore)
store.Register(rule)
```

### 数据库

把 rule 和计数器持久化到 `sys_sequence_rule`：

```go
store := sequence.NewDBStore(db)
```

每条 rule 占 `sys_sequence_rule` 一行。通过 migration 或管理端接口先插入数据。

### Redis

把计数器存到 Redis，适合低延迟、分布式部署：

```go
store := sequence.NewRedisStore(redisClient)
```

> 三种 store 都**要求 rule 先存在再调用 `Generate(...)`**。调一个 store 里没有的 key 会得到 `sequence.ErrRuleNotFound`。

## 生成序列号

框架已经基于你注入的 `Store` 装配好 `sequence.Generator`，业务代码直接注入并调用：

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

`GenerateN` 用一次原子操作预定整个区间，返回的号一定是连续的。

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
| `sequence.ErrInvalidRule` | rule 配置自相矛盾（例如 `SeqStep <= 0`）。 |

## 下一步

参考 [缓存](./cache) 或 [定时任务](./cron)——序列生成器经常与定时归档或预热任务配套使用。
