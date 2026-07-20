---
sidebar_position: 2
---

# Cron Jobs

VEF 将 `github.com/coldsmirk/vef-framework-go/cron` 暴露为 `gocron` 之上的类型化包装。框架 DI 模块会自动提供 `cron.Scheduler`；需要自定义装配时，可以把自己的 `gocron.Scheduler` 传给 `cron.NewScheduler`。

本页介绍**内存调度器**：任务在代码中定义、按进程调度、重启即失。若需要
数据库持久化的调度——集群内单次触发、错触策略、运行流水账与管理 API——见
[持久化调度](./cron-store)；二者同属 `cron` 包，可在一个应用中并存。

## 公开 API 面

该包没有 exported fields。公开 top-level API 如下：

| API | 契约 |
| --- | --- |
| `cron.Scheduler` | DI、自定义 handler 和模块使用的高层 scheduler interface |
| `cron.NewScheduler(scheduler gocron.Scheduler) cron.Scheduler` | 包装调用方提供的 `gocron.Scheduler`；调用方必须传入可用 scheduler |
| `cron.Job` | `Scheduler.NewJob`、`Scheduler.Jobs` 和 `Scheduler.Update` 返回的接口 |
| `cron.JobDefinition` | `Scheduler.NewJob` 和 `Scheduler.Update` 接收的接口 |
| `cron.JobDescriptorOption` | 所有内置 job-definition 构造器接收的 option 类型 |
| `cron.OneTimeJobDefinition` | `NewOneTimeJob` 返回的一次性 job 定义 |
| `cron.DurationJobDefinition` | `NewDurationJob` 返回的固定间隔 job 定义 |
| `cron.DurationRandomJobDefinition` | `NewDurationRandomJob` 返回的随机间隔 job 定义 |
| `cron.CronJobDefinition` | `NewCronJob` 返回的 cron 表达式 job 定义 |
| `cron.NewOneTimeJob(times []time.Time, options ...cron.JobDescriptorOption) *cron.OneTimeJobDefinition` | `times` 为 nil/空时立即执行一次；一个时间点执行一次；多个时间点各执行一次 |
| `cron.NewDurationJob(interval time.Duration, options ...cron.JobDescriptorOption) *cron.DurationJobDefinition` | 按固定 duration 间隔重复执行 |
| `cron.NewDurationRandomJob(minInterval time.Duration, maxInterval time.Duration, options ...cron.JobDescriptorOption) *cron.DurationRandomJobDefinition` | 在 `minInterval` 与 `maxInterval` 之间随机选择间隔重复执行 |
| `cron.NewCronJob(expression string, withSeconds bool, options ...cron.JobDescriptorOption) *cron.CronJobDefinition` | 使用 cron 表达式；`withSeconds=true` 表示表达式带 seconds 字段，`false` 表示标准 5-field 格式 |
| `cron.WithName(name string)` | 设置必需的人类可读 job name |
| `cron.WithTags(tags ...string)` | 设置 `RemoveByTags` 与 `Job.Tags` 使用的 tags |
| `cron.WithConcurrent()` | 允许同一个 job 重叠执行；未设置时使用 singleton wait mode |
| `cron.WithStartAt(startAt time.Time)` | 在指定时间开始调度；优先级高于 `WithStartImmediately` |
| `cron.WithStartImmediately()` | 仅当没有设置 `WithStartAt` 的非零时间时立即开始 |
| `cron.WithStopAt(stopAt time.Time)` | 在指定时间停止调度 |
| `cron.WithLimitedRuns(limitedRuns uint)` | 仅当 `limitedRuns > 0` 时设置运行次数限制 |
| `cron.WithContext(ctx context.Context)` | 将非 nil context 传给底层 job option，用于取消支持 |
| `cron.WithTask(handler any, params ...any)` | 设置必需的函数 handler，并通过 `gocron.NewTask` 转发 `params` |
| `cron.ErrJobNameRequired` | 构建 job option 时没有 name |
| `cron.ErrJobTaskHandlerRequired` | 构建 task 时 handler 为 nil |
| `cron.ErrJobTaskHandlerMustFunc` | 构建 task 时 handler 不是函数 |

## `cron.Scheduler`

公开 scheduler interface 包含：

| 方法 | 签名 | 契约 |
| --- | --- | --- |
| `Jobs` | `Jobs() []cron.Job` | 返回 wrapped scheduler 当前注册的 jobs |
| `NewJob` | `NewJob(definition cron.JobDefinition) (cron.Job, error)` | 构建 definition、注册 job 并返回 `cron.Job`；definition validation error 会先于注册返回 |
| `RemoveByTags` | `RemoveByTags(tags ...string)` | 删除拥有任一指定 tag 的 job |
| `RemoveJob` | `RemoveJob(id string) error` | 先把 `id` 解析为 UUID，再删除该 job；非 UUID 字符串会在 delegation 前返回 parse error |
| `Start` | `Start()` | 启动调度和执行；启动后新增的 job 由 wrapped scheduler 调度 |
| `StopJobs` | `StopJobs() error` | 停止 job 执行但不删除定义；之后可再次 `Start()` |
| `Update` | `Update(id string, definition cron.JobDefinition) (cron.Job, error)` | 先把 `id` 解析为 UUID，再构建替换 definition，并保留 job identifier |
| `JobsWaitingInQueue` | `JobsWaitingInQueue() int` | 返回 wrapped scheduler 的等待队列数量；在 wait mode 下才有实际意义 |

`RemoveJob` 和 `Update` 应使用 `Job.ID()` 返回的 ID。任意字符串不会被接受。

## `cron.Job`

每个已注册 job 都支持以下检查或触发方法：

| 方法 | 签名 | 契约 |
| --- | --- | --- |
| `ID` | `ID() string` | 返回底层 UUID 字符串 |
| `LastRun` | `LastRun() (time.Time, error)` | 返回最近一次 run 的开始时间 |
| `Name` | `Name() string` | 返回配置的 job name |
| `NextRun` | `NextRun() (time.Time, error)` | 返回下一次计划执行时间 |
| `NextRuns` | `NextRuns(count int) ([]time.Time, error)` | 返回 `count` 个未来计划执行时间 |
| `RunNow` | `RunNow() error` | 立即触发一次执行，同时仍遵守 job/scheduler 限制和运行次数限制 |
| `Tags` | `Tags() []string` | 返回配置的 tag 列表 |

## Job Definitions

`cron.JobDefinition` 是 `Scheduler.NewJob(...)` 和 `Scheduler.Update(...)` 接收的公共接口。四个 exported definition structs 由对应构造器返回；调用方通常不直接实例化它们，因为调度字段是内部字段。

| 构造器 | 调度方式 |
| --- | --- |
| `cron.NewOneTimeJob(nil, options...)` 或 `cron.NewOneTimeJob([]time.Time{}, options...)` | 立即执行一次 |
| `cron.NewOneTimeJob([]time.Time{t}, options...)` | 在 `t` 执行一次 |
| `cron.NewOneTimeJob([]time.Time{a, b}, options...)` | 在每个给定时间点各执行一次 |
| `cron.NewDurationJob(interval, options...)` | 固定间隔 |
| `cron.NewDurationRandomJob(minInterval, maxInterval, options...)` | 给定范围内的随机间隔 |
| `cron.NewCronJob(expression, true, options...)` | 带 seconds 字段的 cron 表达式 |
| `cron.NewCronJob(expression, false, options...)` | 标准 5-field cron 表达式 |

所有内置 definition 都先构建 task，再构建 job options。因此如果 task handler 为 nil 或不是函数，同时 name 也缺失，调用方会先看到 task handler 相关错误。

## Job Options

每个构造器都接收 `JobDescriptorOption`：

| Option | 行为 |
| --- | --- |
| `cron.WithName(name)` | 必需；缺失时 `NewJob`/`Update` 返回包裹 `ErrJobNameRequired` 的 error |
| `cron.WithTask(handler, params...)` | 必需；`handler` 必须是函数；`params` 会转发给 `gocron.NewTask(handler, params...)` |
| `cron.WithTags(tags...)` | 保存 tags，用于列表展示和 `RemoveByTags` |
| `cron.WithConcurrent()` | 对该 job 禁用默认 singleton wait mode |
| `cron.WithStartAt(startAt)` | 设置开始 date/time；与 `WithStartImmediately` 同时出现时它优先 |
| `cron.WithStartImmediately()` | 只有 `WithStartAt` 没有设置非零时间时才生效 |
| `cron.WithStopAt(stopAt)` | 仅当 `stopAt` 非零时设置停止 date/time |
| `cron.WithLimitedRuns(limitedRuns)` | 仅当 `limitedRuns > 0` 时设置运行次数限制 |
| `cron.WithContext(ctx)` | 仅当 `ctx` 非 nil 时添加 context option；接收 `context.Context` 的 handler 可感知取消 |

默认情况下，同一个 job 不会和自己的另一次执行重叠。除非提供 `WithConcurrent()`，否则包会添加 singleton wait mode。

## 错误哨兵

| 错误 | 触发条件 |
| --- | --- |
| `cron.ErrJobTaskHandlerRequired` | 未提供 `WithTask`，或 `WithTask` 收到 nil |
| `cron.ErrJobTaskHandlerMustFunc` | `WithTask` 收到非函数 handler |
| `cron.ErrJobNameRequired` | task 构建成功，但缺少 `WithName` 或 name 为空 |

这些 sentinel 会被 build path 包裹，例如 `failed to build job task: ...` 或 `failed to build job options: ...`。判断时请使用 `errors.Is(err, cron.ErrJobNameRequired)` 等形式，不要直接用 `==`。

## DI Scheduler 默认值

框架通过 DI module 创建 scheduler 时，可观察默认值如下：

| 默认值 | 含义 |
| --- | --- |
| 本地时区 | 调度使用 `time.Local` |
| 停止超时 `30s` | 优雅关闭窗口 |
| scheduler logger 和 monitor | job 调度与完成情况会记录日志 |
| 并发 job 上限 `1000` 且使用 wait mode | 超出的执行会进入队列等待，而不是被丢弃 |
| lifecycle start/shutdown | 应用启动时 start scheduler；应用停止时 shutdown scheduler |

这些默认值只适用于框架提供的 scheduler。传给 `cron.NewScheduler` 的自定义 scheduler 使用调用方自己的配置。

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

## 实践场景

以下场景适合使用内存调度器：

- 周期性清理
- 轮询与同步任务
- 定时聚合
- 重试与超时扫描

如果任务本身业务逻辑很重，应把逻辑放进 service，让 scheduler 只负责编排。

当调度必须跨重启存活、在多副本间恰好执行一次、可被运维在线编辑或需要可审计
的运行流水账时，请改用[持久化调度](./cron-store)。

## 下一步

阅读[持久化调度](./cron-store)了解持久化调度引擎；如果定时任务需要显式事务控制，请继续阅读[事务](../data-access/transactions)。
