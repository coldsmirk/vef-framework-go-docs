---
sidebar_position: 14
---

# Sequence

`sequence` 包提供可配置的流水号生成，支持自定义格式和重置策略。

## 核心概念

一个序列由以下部分组成：
- **Rule（规则）**：定义格式模板和重置策略
- **Store（存储）**：持久化当前计数器值（数据库、Redis 或内存）

## 序列规则

```go
rule := sequence.Rule{
    Name:        "order-number",
    Format:      "ORD-{year}{month}{day}-{seq:6}",
    ResetPolicy: sequence.ResetDaily,
}
```

### 格式令牌

| 令牌 | 说明 | 示例 |
| --- | --- | --- |
| `{year}` | 4 位年份 | `2024` |
| `{month}` | 2 位月份 | `03` |
| `{day}` | 2 位日期 | `15` |
| `{hour}` | 2 位小时 | `14` |
| `{minute}` | 2 位分钟 | `30` |
| `{second}` | 2 位秒 | `05` |
| `{seq:N}` | 零填充序号（N 位）| `000001` |

### 重置策略

| 策略 | 常量 | 行为 |
| --- | --- | --- |
| 不重置 | `sequence.ResetNever` | 计数器无限增长 |
| 每日 | `sequence.ResetDaily` | 每天午夜重置 |
| 每月 | `sequence.ResetMonthly` | 每月 1 日重置 |
| 每年 | `sequence.ResetYearly` | 每年 1 月 1 日重置 |

## 存储

### 数据库存储

使用 ORM 数据库持久化——大多数应用推荐：

```go
store := sequence.NewDBStore(db)
```

### Redis 存储

使用 Redis 持久化——高吞吐场景推荐：

```go
store := sequence.NewRedisStore(redisClient)
```

### 内存存储

内存存储——仅用于测试：

```go
store := sequence.NewMemoryStore()
```

## 生成序列

```go
generator := sequence.New(rule, store)

number, err := generator.Next(ctx)
// → "ORD-20240315-000001"

number, err = generator.Next(ctx)
// → "ORD-20240315-000002"
```

## 格式示例

| 场景 | 格式 | 输出示例 |
| --- | --- | --- |
| 订单号 | `ORD-{year}{month}{day}-{seq:6}` | `ORD-20240315-000001` |
| 发票号 | `INV{year}{month}-{seq:4}` | `INV202403-0001` |
| 文档编号 | `DOC-{year}-{seq:8}` | `DOC-2024-00000001` |
| 简单计数器 | `{seq:10}` | `0000000001` |
