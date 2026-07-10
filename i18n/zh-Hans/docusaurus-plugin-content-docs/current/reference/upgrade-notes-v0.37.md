---
sidebar_position: 7
---

# 升级到 v0.36 / v0.37 注意事项

本页是从 `v0.35.0` 到 `v0.37.0` 的跨版本审计地图，覆盖提交范围
`v0.35.0..v0.37.0`。两个版本都以审批模块为中心；v0.37 还新增了事件流
可观测性。如果你的应用文档或集成假设停留在
[升级到 v0.35 注意事项](../reference/upgrade-notes-v0.35)，先按这里过一遍迁移点。

它不是生成索引的替代品。迁移完成后，精确的 Go 符号和 wire 字段仍以
[公开 API 索引](../reference/public-api-index) 与
[运行时 API 索引](../reference/runtime-api-index) 为准。

## 立即检查

- 更新 `approval/flow` 的 RPC 客户端：更新操作的 action 改为 `update`
  （原 `update_flow`），action log 权限改为 `approval.action_log.query`
  （原 `approval.log.query`）。权限种子数据要同步更新。
- 更新列表视图消费者：管理端实例/任务行和我的待办/已办/抄送行的人员字段
  改为嵌套 `UserInfo` 对象（`applicant`、`assignee`），不再是扁平的
  `*Id` / `*Name` 字符串对。
- 委托创建/更新要求必填 `startTime` 和 `endTime`，使用 canonical
  `timex.DateTime` wire 格式（`2006-01-02 15:04:05`），不再是 RFC 3339。
- 审批事件订阅者：`InstanceBindingFailedEvent` 用 `trigger` + `status`
  替代 `finalStatus`；completed/withdrawn 事件新增 `reason`；
  rolled-back/returned 事件新增 `opinion`。
- Go 符号重命名：`Flow.BusinessPkField` 改为 `Flow.BusinessPKField`
  （JSON 仍是 `businessPkField`）。withdrawn、rolled-back、returned、
  binding-failed 事件的构造函数签名有变化。
- 错误码匹配逻辑：`invalid storage mode` 从 40012 移到 40014，
  `flow binding locked` 从 40013 移到 40015（v0.36 新增的 binding-mode /
  initiator-kind 校验错误意外复用了这两个码，v0.37 做了去重）。
- 执行审批 DB 迁移：`apv_cc_record` 新增 `visit_id`，`apv_form_table`
  新增 `source_field_key`，`apv_flow` 新增三个可选的回写联动列。
- 重新测试依赖 rollback 目标、抄送已读确认或加签的流程：三者的语义都收紧了
  （见下方 v0.36 行为修正）。
- 升级到 v0.37 后，实例事件优先使用 `approval.SubscribeInstance` 而不是裸的
  `event.SubscribeTyped`，并考虑启用 `idle_group_retention` 和
  `get_event_streams` 监控 action 来发现并回收孤儿 Redis consumer group。

## 按版本审计

| 版本 | 需要审查的用户可见变化 |
| --- | --- |
| `v0.36.0` | 全部在审批模块：单层明细表表单字段与按字段拆分的子投影表、聚合字段条件（`sum` / `count` / `avg` 加自定义聚合器）、wire 重命名（`update` action、`approval.action_log.query` 权限、嵌套 `UserInfo` 列表行、必填的 canonical 委托时间窗口）、生命周期事件携带操作人 `reason` / `opinion`、保存/发布时校验 binding mode、initiator kind、空条件分支和 sequential + parallel 加签，以及 rollback 边界、按 visit 的抄送作用域和加签插队的行为修正。 |
| `v0.37.0` | 审批五触发点业务回写矩阵（`Flow` 新联动列、重塑的 `InstanceBindingFailedEvent`、去重后的错误码）、声明式 `approval.SubscribeInstance`（派生 consumer group）与 `NewFilteredLifecycleHook`，以及事件流可观测性（`event.StreamInspector`、`get_event_streams`、空闲 consumer group 回收）。 |

## v0.36.0

v0.36 的所有变化都在审批模块。

### 破坏性：明细表表单字段

`approval.FieldKind` 新增 `FieldTable`（`"table"`）：单层明细表，取值是
行的列表。`FormFieldDefinition.Columns` 定义行的形状；每一列本身也是字段
定义，复用现有的 kind 和校验规则，但列不能再是 table。对 table 字段本身，
`Validation.MinLength` / `MaxLength` 约束的是行数，`IsRequired` 表示
"至少一行"。

对 `StorageTable` 的流程，每个明细表字段会投影成独立子表，命名为
`apv_form_<versionID>__<sanitized field key>`，行以 `instance_id` 关联
（有索引、不唯一），按 `row_index` 排序。`approval.FormTable` registry 现在
每张物理表一行——主投影表加上每个明细表字段一张子表——并新增
`SourceFieldKey`（`sourceFieldKey`；主表为空）。唯一约束从 `version_id`
变为 `(version_id, source_field_key)`。表存储现在还会把绑定到文本列的
非字符串标量做字符串化，而不是失败。

### 破坏性：聚合字段条件

字段条件可以把明细表折叠成一个可比较的数值。`approval.Condition` 新增两个
wire 字段：

```json
{
  "kind": "field",
  "subject": "items",
  "aggregate": "sum",
  "column": "amount",
  "operator": "gt",
  "value": 10000
}
```

`approval.AggregateKind` 定义 `sum`、`count`、`avg`。语义遵循 SQL 聚合：
`count` 是行数且 `column` 必须留空，空表的 `sum` 是 0，空表的 `avg` 不匹配
任何比较（NULL 语义）。折叠在 `float64` 中计算，比较小数金额的 sum / avg
时优先用大小比较运算符而不是 `eq` / `ne`。发布校验通过
`AggregateKind.FoldsColumn()` 推导 column 必填/禁止规则。

自定义聚合实现 `approval.Aggregator`（`Kind()` 加
`Fold(values, rowCount) (result, matchable)`），通过
`vef.ProvideApprovalAggregator` 注册：

```go
vef.ProvideApprovalAggregator(func() approval.Aggregator { return myMedian{} })
```

### 破坏性：wire 与权限重命名

- `approval/flow` resource 的 `update_flow` action 改为 `update`。Go 参数
  类型跟随重命名（`UpdateFlowParams` → `UpdateParams`、
  `StartInstanceParams` → `StartParams`）。
- action log 查询权限改为 `approval.action_log.query`
  （原 `approval.log.query`）。
- 列表行嵌套人员快照。以管理端实例行为例，前 / 后：

```json
{ "applicantId": "u1", "applicantName": "Alice" }
```

```json
{ "applicant": { "id": "u1", "name": "Alice", "departmentId": "d1", "departmentName": "Sales" } }
```

  同样的形状适用于管理端任务行和我的待办/已办行的 `assignee`，以及我的
  抄送行的人员字段。
- 委托创建/更新使用 canonical `timex.DateTime` 格式，时间窗口现在必填：

```json
{ "startTime": "2030-01-01T00:00:00Z", "endTime": null }
```

```json
{ "startTime": "2030-01-01 00:00:00", "endTime": "2030-06-01 00:00:00" }
```

### 破坏性：生命周期事件携带 reason 和 opinion

- `InstanceCompletedEvent.Reason`（`reason`，可选）携带管理员终止
  （`terminated`）时填写的原因；approved/rejected 时为 nil，决定性意见在
  task 事件上。
- `InstanceWithdrawnEvent.Reason` 携带申请人的撤回原因。
- `InstanceRolledBackEvent.Opinion` 和 `InstanceReturnedEvent.Opinion`
  携带操作人的回退意见。

对应的 `New*Event` 构造函数增加了尾部参数，直接构造这些事件的 Go 代码
需要更新。

### 破坏性：更严格的保存/发布校验

- 保存流程时拒绝枚举之外的 `BindingMode` 和 `InitiatorKind`
  （`BindingMode.IsValid()` / `InitiatorKind.IsValid()` 已导出）。以前
  写错的 binding mode 会静默按 `standalone` 处理，直接禁用业务回写。
- 发布时拒绝没有 condition group 的条件分支和没有 condition 的 group——
  引擎把结构上为空的条件当作无条件匹配，会静默遮蔽所有更低优先级的分支。
- 发布时拒绝在 sequential 审批节点上配置 `parallel` 加签类型；顺序队列
  没有可供并行加入的东西。

### 破坏性行为修正

- **rollback 目标以 visit trail 为边界。** rollback 目标必须是本实例实际
  走过（有已完结 visit）的决策节点（approval / handle）或 start 节点。
  以前版本内任意节点都可接受，任务持有人可以指向 End 节点强制以通过
  结束实例。
- **抄送记录和已读确认门槛按节点 visit 作用域。** `CCRecord` 新增
  `VisitID`（`visitId`），rollback 重走时获得自己的通知和已读确认周期，
  不会被上一轮的记录静默满足。
- **加签跳过在当前 visit 中已决策的用户**（approved / rejected /
  handled），pass rule 不会强迫同一人在一轮里决策两次。被移除或转办的
  用户可以重新加签。
- **sequential 加签在锚点处插队**，不再追加到队列末尾："before" 加签
  接管锚点位置，"after" 加签紧随其后。

### 其他修正

- 撤回实例后，flow graph 中暂停的驻留节点保持 active，不再显示无当前节点。
- 我的详情 action 列表只对申请人或持有 assignee 任务的用户提供 `urge`，
  与 urge handler 的实际授权一致。
- 审批分页列表使用 `id` tiebreaker 保证排序稳定。
- 新索引加速已办任务列表（`finished_at`）和 `apv_task` 的 `visit_id` 查询。

## v0.37.0

### 破坏性：五触发点业务回写矩阵

engine-owned 业务回写现在覆盖完整实例生命周期。`approval.BindingTrigger`
标识驱动回写的时刻；每个触发点投影固定的列子集（只有流程配置了的列才会
被写入）：

| Trigger | status | instance_id | started_at | finished_at |
| --- | --- | --- | --- | --- |
| `started` | running | 实例 ID | now | NULL |
| `completed` | 最终状态 | — | — | `FinishedAt` |
| `returned` | returned | — | — | — |
| `withdrawn` | withdrawn | — | — | — |
| `resubmitted` | running | — | — | NULL |

- `approval.Flow` 新增三个可选联动列：`BusinessInstanceIDField`、
  `BusinessStartedAtField`、`BusinessFinishedAtField`
  （`businessInstanceIdField` / `businessStartedAtField` /
  `businessFinishedAtField`）。业务绑定流程只有状态列仍是必填。Go 侧
  `Flow.BusinessPkField` 重命名为 `Flow.BusinessPKField`（JSON 不变）。
- `started` 回写在 `start_instance` 事务内同步执行——失败会回滚整个发起。
  其余四个通过 binding listener 异步执行，以
  `InstanceBindingFailedEvent` 补偿，所以 `started` 不会出现在该事件里。
- `InstanceBindingFailedEvent` 重塑；前 / 后：

```json
{ "finalStatus": "approved", "businessTable": "orders", "errorMessage": "..." }
```

```json
{ "trigger": "completed", "status": "approved", "businessTable": "orders", "errorMessage": "..." }
```

- 保存流程时校验绑定字段必须指向互不相同的业务列
  （`approval binding columns conflict`，错误码 40016）；重复列会渲染成
  `SET col = ?, col = ?`，在申请人的发起事务中以运行时 SQL 错误暴露。
  现有的绑定冻结（有实例运行时禁止修改绑定）现在也覆盖新联动列。
- 错误码去重：`invalid storage mode` 改为 40014（原 40012），
  `flow binding locked` 改为 40015（原 40013）；40012 / 40013 现在属于
  v0.36 的 binding-mode / initiator-kind 校验错误。

### 声明式实例订阅

`approval.SubscribeInstance` 用声明式路由过滤器把类型化 handler 订阅到
某个实例事件类型，替代手写的 `event.SubscribeTyped` 接线：

```go
unsubscribe, err := approval.SubscribeInstance(bus,
    svc.OnCompleted, // func(ctx context.Context, evt *approval.InstanceCompletedEvent) error
    approval.ForFlows("expense_claim"),
    approval.WithGroup("mms.expense.completed"),
)
```

- `approval.InstanceFilter`（`ForFlows`、`ForTenants`）把"这个实例是不是
  我的？"表达为数据。单个 filter 内已填充的维度按 OR 匹配；多个 filter
  之间是 AND。不匹配的事件直接 ack，不调用 handler。业务谓词（最终状态、
  表单值）留在 handler 内部。
- consumer group 默认从 handler 的方法身份派生
  （`vef:sub:<pkg>.<Type>.<Method>`，去掉主模块前缀）。匿名函数报
  `ErrAnonymousSubscriberGroup`；同进程内派生出相同 group 报
  `ErrDerivedGroupConflict`。重命名或移动 handler 会改变派生 group 并把
  旧的变成孤儿——生产关键订阅者在这类重构前先用 `approval.WithGroup`
  固定名字。`approval.WithConcurrency` 透传每订阅的 worker 数。
- `approval.NewFilteredLifecycleHook(hook, filters...)` 把相同过滤器应用到
  同步的 `InstanceLifecycleHook`，只服务单一流程的 hook 可以在
  `vef.ProvideApprovalLifecycleHook` 注册处收窄作用域。

### 事件流可观测性

redis_stream transport 现在暴露 consumer group 状态，并能回收被下线
订阅者遗留的孤儿 group：

- `event.StreamInspector` 列出 transport 前缀下的每个 stream 及其
  consumer group（`StreamInfo` / `StreamGroupInfo`：name、consumers、
  pending、lag、`lastDeliveredId`）。它是可选依赖——redis_stream 关闭时
  为 nil。
- `sys/monitor` resource 新增 `get_event_streams` action，返回
  `monitor.EventStreamsInfo`（`enabled` 加 stream 列表）。lag 持续增长而
  consumer 全部空闲的 group 就是孤儿候选。
- 空闲 group 回收是 opt-in：

```toml
[vef.event.transports.redis_stream]
enabled = true
idle_group_retention = "72h"        # 零值（默认）禁用清扫
idle_group_sweep_interval = "10m"   # 默认 10m
```

只有当 group 没有 pending 条目、且每个 consumer 记录的空闲时间都超过
retention 窗口时才会被销毁。

完整契约见[审批模块](../approval)、[事件总线](../infrastructure/event-bus)
和[监控](../infrastructure/monitor)。
