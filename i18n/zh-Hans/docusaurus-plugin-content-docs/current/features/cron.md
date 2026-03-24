---
sidebar_position: 3
---

# Cron Jobs

VEF 在 `gocron` 之上提供了一层调度器抽象。

## 模块提供了什么

cron 模块通过依赖注入提供 `cron.Scheduler`，并通过应用生命周期自动启动底层 scheduler。

## `cron.Scheduler` 接口

公开调度器接口包含：

| 方法 | 作用 |
| --- | --- |
| `Jobs()` | 列出所有已注册 job |
| `NewJob(definition)` | 创建并注册一个新 job |
| `RemoveByTags(tags...)` | 按 tag 批量移除 job |
| `RemoveJob(id)` | 按 ID 移除单个 job |
| `Start()` | 启动调度和执行 |
| `StopJobs()` | 停止 job 执行，但保留定义 |
| `Update(id, definition)` | 替换现有 job 定义 |
| `JobsWaitingInQueue()` | 在 wait 模式下查看等待队列长度 |

## `cron.Job` 接口

每个已注册 job 都支持以下检查方法：

| 方法 | 作用 |
| --- | --- |
| `ID()` | job 唯一标识 |
| `LastRun()` | 最近一次执行时间 |
| `Name()` | 可读名称 |
| `NextRun()` | 下一次执行时间 |
| `NextRuns(count)` | 多个未来执行时间 |
| `RunNow()` | 立即触发一次执行 |
| `Tags()` | 当前 job 的 tag 列表 |

## 支持的 JobDefinition

当前公共包支持以下 job-definition 构造器：

| 构造器 | 含义 |
| --- | --- |
| `cron.NewOneTimeJob(times, options...)` | 立即执行一次，或在一个/多个指定时间点各执行一次 |
| `cron.NewDurationJob(interval, options...)` | 按固定间隔重复执行 |
| `cron.NewDurationRandomJob(minInterval, maxInterval, options...)` | 在一个随机间隔范围内重复执行 |
| `cron.NewCronJob(expression, withSeconds, options...)` | 使用 cron 表达式调度 |

## JobDescriptor 配置项

job 可通过以下 option 定制：

| Option | 作用 |
| --- | --- |
| `cron.WithName(name)` | 设置 job 显示名称 |
| `cron.WithTags(tags...)` | 为 job 设置 tag，便于分组和批量移除 |
| `cron.WithConcurrent()` | 允许同一个 job 并发执行 |
| `cron.WithStartAt(time)` | 指定调度开始时间 |
| `cron.WithStartImmediately()` | 调度器启动后立即开始 |
| `cron.WithStopAt(time)` | 指定调度停止时间 |
| `cron.WithLimitedRuns(count)` | 限制最大执行次数 |
| `cron.WithContext(ctx)` | 关联一个可取消的 context |
| `cron.WithTask(handler, params...)` | 设置实际执行的任务函数与参数 |

## 最小示例

```go
package jobs

import (
  "context"
  "time"

  "github.com/coldsmirk/vef-framework-go/cron"
)

func RegisterCleanupJob(scheduler cron.Scheduler) error {
  _, err := scheduler.NewJob(
    cron.NewDurationJob(
      10*time.Minute,
      cron.WithName("cleanup-expired-sessions"),
      cron.WithTask(func(ctx context.Context) error {
        return nil
      }),
    ),
  )

  return err
}
```

## 运行时默认值

内部 scheduler 创建时带有这些默认行为：

| 默认值 | 含义 |
| --- | --- |
| 本地时区 | 调度使用 `time.Local` |
| 停止超时 `30s` | 优雅关闭窗口 |
| 并发 job 上限 `1000` | 超出后进入等待队列，而不是丢弃 |

这些值是实现默认值，不应视为应用对外策略承诺。

## 实践场景

以下场景适合使用 cron：

- 周期性清理
- 轮询与同步任务
- 定时聚合
- 重试与超时扫描

如果任务本身业务逻辑很重，应把逻辑放进 service，让 scheduler 只负责编排。

## 下一步

继续阅读 [事务](../advanced/transactions)，如果定时任务需要显式事务控制，就会接到那一层。
