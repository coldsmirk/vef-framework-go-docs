---
sidebar_position: 3
---

# 持久化调度

持久化调度存储（v0.39）在 [cron 模块](./cron)之上扩展出数据库持久化的
调度：集群内每次触发只执行一次、运维可在线编辑触发器、错触（misfire）
策略、运行流水账与崩溃恢复。内存调度器继续服务进程内任务；存储引擎是一套
独立机制，面向必须跨重启存活、跨节点协调的任务。

默认关闭。启用后从主数据源加载调度，并挂载 `sys/cron/schedule` 与
`sys/cron/run` 资源：

```toml
[vef.cron.store]
enabled = true
auto_migrate = true   # 启动时创建 crn_schedule / crn_fire_request / crn_run
```

## 模型

引擎由两个概念驱动：

- **任务处理器**（`cron.JobHandler`）是启动时以唯一名称注册的 Go 代码。
  处理器是"做什么"。
- **调度**（`cron.Schedule`，表 `crn_schedule`）是持久化的触发器：*何时*
  触发*哪个*任务、带什么参数、受哪些策略约束。调度是数据——可在代码中或
  管理 API 中创建，运行期可编辑。

每次触发都以**运行**（`cron.Run`，表 `crn_run`）形式记入流水账。

### 注册任务处理器

```go
vef.ProvideCronJobHandler(func(svc *ReportService) cron.JobHandler {
    return cron.NewTypedJobHandler("daily-report",
        func(ctx context.Context, params ReportParams) error {
            return svc.Generate(ctx, params)
        },
        // 可选：启动时若同名调度不存在则播种默认调度；
        // 运维修改永远不会被覆盖。
        cron.WithDefaultSchedule(cron.ScheduleSpec{
            Trigger: cron.Expr("0 2 * * *", "Asia/Shanghai"),
        }),
    )
})
```

| API | 契约 |
| --- | --- |
| `cron.JobHandler` | `Name() string` + `Execute(ctx, execution) error`；每个任务名恰好一个处理器 |
| `cron.NewJobHandler(name, execute, opts...)` | 适配函数；`execute` 接收完整的 `cron.Execution` |
| `cron.NewTypedJobHandler[P](name, execute, opts...)` | 运行前把调度参数解码为 `P`；解码失败直接把该运行记为 failed 而不调用函数 |
| `cron.WithDefaultSchedule(spec)` | 随处理器携带默认调度；存储在启动时按缺失播种。spec 的 `Name` 回退为任务名 |
| `cron.Execution` | 运行的只读视图：`RunID`、`ScheduleID`、`ScheduleName`、`JobName`、`ScheduledAt`（逻辑触发时间）、`Params`（原始 JSON）、`BindParams(v)` |

每次触发至多被一个节点认领，但当调度设置了 `Recover` 时，崩溃的运行会重新
触发——投递是 at-least-once，处理器应当幂等。

### 触发器

`cron.TriggerSpec` 声明调度何时触发；用构造函数创建：

| 构造函数 | 类型 | 语义 |
| --- | --- | --- |
| `cron.Expr(expr, timezone)` | `cron` | 在 IANA 时区中求值的 cron 表达式。支持 5 段、6 段（前导秒）与 `@` 描述符（`@daily`、`@every 90m`）。时区为空解析为 `UTC` ——持久化调度永远不依赖节点的进程本地时区（`"Local"` 被拒绝）。内嵌 tzdata 保证无 zoneinfo 的部署也能加载时区 |
| `cron.Every(duration)` | `interval` | 固定频率，最小 `1s`（`cron.MinInterval`）。频率锚定在调度起点（`StartsAt`，否则创建时间），触发相位不受补偿触发和手动触发影响 |
| `cron.Once(at)` | `once` | 单次触发 |

不属于所选类型的字段会被拒绝（`ErrTriggerFieldsConflict`），无法解析的
表达式、无法加载的时区、低于 1 秒的间隔、缺失的触发时间同样被拒绝。

### 调度声明

`cron.ScheduleSpec` 声明要创建的调度，或更新的目标状态：

| 字段 | 含义 |
| --- | --- |
| `Name` | 唯一管理键；播种的默认调度回退为任务名 |
| `JobName` | 要执行的已注册 `JobHandler` |
| `Trigger` | 何时触发（见上） |
| `Params` | JSON 序列化后每次运行原样交给处理器 |
| `StartsAt` / `EndsAt` | 可选触发窗口；`StartsAt` 同时是 interval 触发器的相位锚点 |
| `MisfirePolicy` | `fire_now`（默认）或 `skip`（见下） |
| `ConcurrencyPolicy` | `forbid`（默认）或 `allow`（见下） |
| `Recover` | 重新触发执行中途被遗弃的运行；要求处理器幂等 |
| `Timeout` | 单次运行上限（必须是整毫秒）；0 继承 `vef.cron.store.run_timeout` |
| `Enabled` | 初始/更新后的启用状态；`nil` 表示启用 |

调用方提供的时间（`Trigger.At`、`StartsAt`、`EndsAt`）在持久化前统一
归一化到本地墙钟，其他时区构建的时刻读回后仍表示同一瞬间。

## 策略

### 错触（Misfire）

超过 `vef.cron.store.misfire_threshold`（默认 `1m`）才开始的触发计为
错触——停机、暂停中的调度或没有空闲执行槽。此时应用调度的
`MisfirePolicy`：

| 策略 | 行为 |
| --- | --- |
| `fire_now`（默认） | 立即补跑一次，然后从现在恢复常规序列 |
| `skip` | 跳到下一个未来触发，不补跑 |

无论哪种策略，永远不会执行的那些次数会以一条 `missed` 运行记入流水账，
覆盖整个缺口（`missedCount` 记录次数）。

### 并发

| 策略 | 行为 |
| --- | --- |
| `forbid`（默认） | 与同调度仍在执行的运行重叠的触发被抑制并记为 `skipped`。恢复请求保持等待直到活动运行结束 |
| `allow` | 同一调度的运行可以重叠 |

### 暂停 / 恢复语义

`Pause` 清除运维所有的 `isEnabled` 标志；执行中的运行不受影响。触发游标在
暂停期间被刻意保留，因此 `Resume` 会把暂停缺口交给错触策略处理，而不是
静默丢弃：`fire_now` 下立即补跑一次，`skip` 下等待下一个常规触发。

### 手动触发

`TriggerNow` 持久化一条独立的立即触发请求（表 `crn_fire_request`）——单节点、
入流水账、遵守并发策略——且不移动常规触发游标。恢复重触发也走同一张请求表。
暂停中的调度拒绝并返回 `ErrScheduleDisabled`。当手动触发与常规触发落在同一
逻辑时刻时，常规触发优先，两者都按策略入流水账。

## 执行与恢复

- 每个节点轮询到期调度（`poll_interval`，默认 `5s`，并自适应睡眠到最近的
  已知触发——该间隔是其他节点新建调度的可见性延迟，不是触发精度），按事务
  认领至多 `batch_size` 个触发，在本节点至多 `max_concurrent` 个槽位上
  执行。
- 执行器每 `heartbeat_interval`（默认 `10s`）为运行中的流水账行续心跳。
  心跳陈旧超过 `abandoned_after`（默认 `1m`；必须至少是心跳间隔的两倍）的
  running 行由恢复清扫在一个事务内接管并标记为 `abandoned`；设置了
  `Recover` 的调度会将其作为全新运行重新触发。
- 超过超时仍未结束的运行记为 `failed`；优雅停机时被中断的运行记为
  `canceled`。
- 重塑调度（修改触发器/窗口）会重算下次触发但保留触发历史；重命名后流水账
  经反规范化的名称保持关联。

## 运行流水账

`cron.Run`（表 `crn_run`）记录每次触发。行在调度删除后仍保留——
`scheduleName` 与 `jobName` 为此做了反规范化。

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `id` | `string` | 流水账行 ID |
| `scheduleId` / `scheduleName` | `string` | 触发的调度 |
| `jobName` | `string` | 执行的处理器 |
| `scheduledAtUnixMs` | `int64` | 逻辑触发时间；补跑的实际开始晚于它。刻意不唯一：手动与恢复触发可能合法地共享同一时刻 |
| `claimedAtUnixMs` | `int64` | 节点认领触发的时间 |
| `status` | `string` | `running`、`succeeded`、`failed`、`missed`、`skipped`、`abandoned`、`canceled` |
| `nodeId` | `string` | 执行节点；从未执行的行（`missed`、`skipped`）为空 |
| `startedAtUnixMs` / `finishedAtUnixMs` | `int64` | 执行窗口 |
| `durationMs` | `int64` | 执行时长 |
| `heartbeatAtUnixMs` | `int64` | 执行器活性信号；陈旧即转为 `abandoned` |
| `error` | `string` | 失败消息（截断）；成功为空 |
| `missedCount` | `int` | 一条 `missed` 行覆盖的次数 |

`run_retention` 会（每小时清扫）删除超窗的终态流水账行；0 表示永久保留——
删除流水账严格 opt-in。

## 编程式管理

只要加载了 cron 模块，`cron.ScheduleManager` 就在 DI 中可用；存储关闭时
所有方法返回 `ErrStoreDisabled`。API 变更与编程式变更共享同一套校验与唤醒
路径。

| 方法 | 契约 |
| --- | --- |
| `Create(ctx, spec)` | 校验并持久化新调度；任务名必须在本节点注册；名称已占用返回 `ErrScheduleExists` |
| `Update(ctx, name, spec)` | 重塑指定调度，spec 携带不同的未占用名称时同时重命名；触发器/窗口变更重算下次触发 |
| `Delete(ctx, name)` | 删除调度；流水账保留 |
| `Pause(ctx, name)` / `Resume(ctx, name)` | 见[暂停语义](#暂停--恢复语义) |
| `TriggerNow(ctx, name)` | 见[手动触发](#手动触发) |
| `Get(ctx, name)` | 返回指定调度，或 `ErrScheduleNotFound` |
| `List(ctx, filter)` | 匹配 `ScheduleFilter`（`JobName`、`Enabled *bool`）的调度，按名称排序 |
| `ListRuns(ctx, filter)` | 匹配 `RunFilter`（`ScheduleName`、`JobName`、`Statuses`、逻辑触发时间的 `Since`/`Until`、`Limit` —— 0 解析为 100，上限 1000）的流水账，最新在前 |

## 事件

两个主题都是尽力而为的运维通知，在任何事务之外经默认事件路由发布——用于
告警订阅；绝不要用它们驱动正确性（运行流水账才是持久事实）。

| 主题 | 事件 | 字段 |
| --- | --- | --- |
| `vef.cron.run.failed` | `cron.RunFailedEvent` | `runId`、`scheduleName`、`jobName`、`scheduledAtUnixMs`、`nodeId`、`error` |
| `vef.cron.run.abandoned` | `cron.RunAbandonedEvent` | `runId`、`scheduleName`、`jobName`、`scheduledAtUnixMs`、`nodeId` |

## RPC 资源

存储启用后，两个管理资源挂载在 `/api` 下。存储关闭时资源不挂载任何操作——
关闭的特性不暴露任何表面。变更类操作均记入审计。

### `sys/cron/schedule`

| 操作 | 权限 | 输入 | 输出 |
| --- | --- | --- | --- |
| `find_page` | `cron.schedule.query` | `ScheduleSearch` + 分页 meta | `page.Page[Schedule]` |
| `get` | `cron.schedule.query` | `ScheduleNameParams` | `ScheduleDetail` |
| `list_jobs` | `cron.schedule.query` | 无 | `string[]` |
| `preview_fires` | `cron.schedule.query` | `PreviewFiresParams` | `FiresPreview` |
| `create` | `cron.schedule.manage`（审计） | `ScheduleParams` | 创建后的 `Schedule` |
| `update` | `cron.schedule.manage`（审计） | `ScheduleParams` | 更新后的 `Schedule` |
| `delete` | `cron.schedule.manage`（审计） | `ScheduleNameParams` | 成功 |
| `pause` | `cron.schedule.manage`（审计） | `ScheduleNameParams` | 成功 |
| `resume` | `cron.schedule.manage`（审计） | `ScheduleNameParams` | 成功 |
| `trigger_now` | `cron.schedule.manage`（审计） | `ScheduleNameParams` | 成功 |

`ScheduleSearch`（`find_page` 查询过滤）：

| 字段 | 类型 | 匹配 | 说明 |
| --- | --- | --- | --- |
| `name` | `string` | contains | 按调度名片段过滤 |
| `jobName` | `string` | equals | 按任务名过滤 |
| `kind` | `string` | equals | 触发器类型：`cron`、`interval` 或 `once` |
| `isEnabled` | `bool` | equals | 按启用状态过滤 |

`ScheduleNameParams`（`get`、`delete`、`pause`、`resume`、`trigger_now`
使用）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `name` | `string` | 是 | 调度的唯一名称 |

`ScheduleParams`（create/update；未知字段会被拒绝——参数结构体是严格
模式）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `name` | `string` | 是 | create 时：新调度的唯一名称；update 时：被寻址的调度 |
| `newName` | `string` | 否 | 仅 update：先按 `name` 寻址再重命名 |
| `jobName` | `string` | 是 | 要执行的已注册任务处理器；未注册返回 `ErrJobNotRegistered` |
| `trigger` | `TriggerParams` | 是 | 触发器定义（见下） |
| `params` | 任意 JSON 值 | 否 | 每次运行原样交给处理器 |
| `startsAtUnixMs` | `int64`（unix 毫秒） | 否 | 触发窗口起点；同时是 interval 触发器的相位锚点 |
| `endsAtUnixMs` | `int64`（unix 毫秒） | 否 | 触发窗口终点；必须晚于 `startsAtUnixMs` |
| `misfirePolicy` | `string` | 否 | `fire_now`（省略时默认）或 `skip` |
| `concurrencyPolicy` | `string` | 否 | `forbid`（省略时默认）或 `allow` |
| `recover` | `bool` | 否 | 重新触发被遗弃的运行；处理器必须幂等 |
| `timeoutMs` | `int64` | 否 | 单次运行超时；0 继承 `vef.cron.store.run_timeout`；负值被拒绝 |
| `enabled` | `bool` | 否 | 省略表示启用 |

`TriggerParams`（只允许所选 `kind` 的字段；多余字段返回
`ErrTriggerInvalid`）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `kind` | `string` | 是 | `cron`、`interval` 或 `once` |
| `expr` | `string` | `cron` 必填 | cron 表达式（5/6 段或 `@` 描述符） |
| `timezone` | `string` | 否（仅 `cron`） | 表达式求值的 IANA 时区；空表示 `UTC`；`"Local"` 被拒绝 |
| `everyMs` | `int64` | `interval` 必填 | 固定频率（毫秒），最小 1000 |
| `atUnixMs` | `int64`（unix 毫秒） | `once` 必填 | 单次触发时间 |

`ScheduleDetail`（`get` 响应）：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `schedule` | `Schedule` | 调度行（见下） |
| `nextFiresUnixMs` | `int64[]` | 从现在起接下来（至多 5 个）精确触发时间的预览。过期游标按调度错触策略投影；暂停或已耗尽的调度返回空列表 |

`Schedule`（`get`、`create`、`update` 与 `find_page` 条目返回；省略标准
审计字段）：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `name` | `string` | 唯一管理键 |
| `jobName` | `string` | 调度触发的处理器 |
| `kind` | `string` | 触发器类型 |
| `expr` | `string` | cron 表达式（cron 类型） |
| `timezone` | `string` | 求值时区（cron 类型） |
| `everyMs` | `int64` | 固定频率（interval 类型） |
| `fireAtUnixMs` | `int64` | 单次触发时间（once 类型）；其他类型缺省 |
| `startsAtUnixMs` / `endsAtUnixMs` | `int64` | 触发窗口边界；无界时缺省 |
| `anchorAtUnixMs` | `int64` | 固定频率相位锚点（创建时间，或设置了 `startsAtUnixMs` 时为其值） |
| `params` | JSON | 处理器参数，原样 |
| `misfirePolicy` | `string` | `fire_now` 或 `skip` |
| `concurrencyPolicy` | `string` | `forbid` 或 `allow` |
| `recover` | `bool` | 遗弃运行重触发标志 |
| `timeoutMs` | `int64` | 单次运行超时；`0` 继承配置默认 |
| `isEnabled` | `bool` | 运维所有的启用状态（`pause` 清除、`resume` 恢复） |
| `nextFireAtUnixMs` | `int64` | 引擎将认领的下一次触发；触发器不再产生新次数（完成的单次、过期的窗口）时缺省——暂停会保留它 |
| `lastFireAtUnixMs` | `int64` | 最近一次被认领触发的逻辑时间；首次触发前缺省 |

`list_jobs` 返回**本节点**注册的任务名——调度编辑器任务选择器的候选词汇。
异构部署下各节点注册集可能不同；返回的是应答节点的视图。

`PreviewFiresParams`（`preview_fires` —— 在编辑器阶段用真实解析器验证
未保存的触发器；拒绝的恰好是保存会拒绝的）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `trigger` | `TriggerParams` | 是 | 要投影的未保存触发器 |
| `startsAtUnixMs` | `int64`（unix 毫秒） | 否 | 投影所用的窗口起点 |
| `endsAtUnixMs` | `int64`（unix 毫秒） | 否 | 窗口终点；必须晚于起点 |

`FiresPreview` 响应：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `nextFiresUnixMs` | `int64[]` | 触发器从现在起的触发时间（至多 5 个）；窗口内不再产生次数时为空 |

### `sys/cron/run`

只读流水账视图：分页视图用于浏览，单条视图用于查看完整错误文本。默认按
认领时间倒序。

| 操作 | 权限 | 输入 | 输出 |
| --- | --- | --- | --- |
| `find_page` | `cron.run.query` | `RunSearch` + 分页 meta | `page.Page[Run]` |
| `find_one` | `cron.run.query` | `RunSearch` | 一条 `Run` |

`RunSearch`（查询过滤）：

| 字段 | 类型 | 匹配 | 说明 |
| --- | --- | --- | --- |
| `id` | `string` | equals | 寻址单条流水账行——`find_one` 没有其他方式指名记录 |
| `scheduleName` | `string` | equals | 按调度过滤 |
| `jobName` | `string` | equals | 按任务过滤 |
| `status` | `string` | equals | 运行状态之一 |
| `nodeId` | `string` | equals | 按执行节点过滤 |
| `scheduledAtFromUnixMs` | `int64` | ≥ | 逻辑触发时间下界 |
| `scheduledAtToUnixMs` | `int64` | ≤ | 逻辑触发时间上界 |

`Run` 响应字段即[运行流水账字段](#运行流水账)加上创建审计字段。

## 错误码

Cron API 错误使用响应码 `2700`–`2799`，以 HTTP 200 承载、失败由响应体
code 表达。

| 码 | 错误 | 含义 |
| --- | --- | --- |
| `2700` | `ErrScheduleNotFound` | 调度不存在 |
| `2701` | `ErrScheduleExists` | 调度名已占用 |
| `2702` | `ErrScheduleDisabled` | 对暂停中的调度手动触发 |
| `2703` | `ErrTriggerInvalid(reason)` | 触发器校验失败（字段冲突、坏表达式、坏时区、间隔过短、缺触发时间） |
| `2704` | `ErrJobNotRegistered` | 调度引用了本节点未注册的任务名 |
| `2705` | `ErrStoreDisabled` | `vef.cron.store.enabled = false` 时调用存储操作 |
| `2706` | `ErrScheduleInvalid(reason)` | 非触发器的声明故障（名称、窗口、超时、参数、策略词汇） |

## 配置

```toml
[vef.cron.store]
enabled = false            # 总开关；关闭时不触碰任何表
auto_migrate = false       # 启动时执行 cron DDL 迁移
poll_interval = "5s"       # 调度表重读上限（可见性延迟，不是触发精度）
batch_size = 32            # 每个轮询周期认领的调度数
max_concurrent = 16        # 每节点并发运行数
misfire_threshold = "1m"   # 触发晚到多久后应用错触策略
heartbeat_interval = "10s" # 执行器对运行中行的活性节律
abandoned_after = "1m"     # 心跳陈旧窗口；必须 ≥ 2 × heartbeat_interval
run_timeout = "0s"         # 默认单次运行上限；0 表示不限
run_retention = "0s"       # 流水账保留；0 表示永久保留
```

启动校验拒绝负的时长以及比心跳间隔两倍更紧的 `abandoned_after`
（否则健康执行器会被判死亡）。

## 下一步

面向进程内工作的内存调度器见 [Cron 定时任务](./cron)。要对失败或被遗弃的
运行做告警，请通过[事件总线](./event-bus)订阅上述事件。
