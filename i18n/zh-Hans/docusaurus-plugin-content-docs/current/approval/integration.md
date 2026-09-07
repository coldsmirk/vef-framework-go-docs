---
sidebar_position: 5
---

# 事件与集成

## 事件发布

审批模块通过框架统一的事务性 outbox transport 发布领域事件（见 [事件总线](../infrastructure/event-bus.md)）。每条审批命令都把事件记录写在和业务变更相同的事务里，再由 outbox relay 转发到配置的 sink。单次 engine pass 内产生的事件按发生顺序发布。

订阅方必须：

1. 用 `event.WithGroup("...")` 挂上 —— 路由解析到 at-least-once transport。
2. 依赖 Inbox 中间件去重（在 `event.middleware.inbox = true` 且 transport 声明 `AtLeastOnce` 时自动激活）。

同一事务内的事件按它们发生的顺序发布。引擎在状态变更发生的那一刻立即
发出事件，outbox relay 按该发生顺序转发。

审批命令通过 `event.WithTx(db)` 在业务事务内发布领域事件；只有事务提交后事件
才对订阅方可见。所有审批领域事件共享 `approval.*`  topic 命名空间。

> 老版本里的独立 `apv_event_outbox` 表和 `EventOutboxStatus` 常量已经废弃。审批模块不再自维护一个 outbox，而是组合使用框架统一 outbox。

### 领域事件类型

所有审批事件都实现 `DomainEvent`；`AllEventTypes()` 是导出的 topic 字符串完整注册表。实例和任务事件嵌入 `InstanceEventBase`（`instanceId`、`instanceNo`、`tenantId`、`title`、`flowId`、`flowCode`、可选 `businessRef`、`applicant`、`occurredTime`）。任务事件还嵌入 `TaskEventBase`（`taskId`、`nodeId`、`nodeName`）。流程事件嵌入 `FlowEventBase`（`flowId`、`tenantId`、`code`、`name`、`occurredTime`）。下表只列出这些 base 之外的字段。

实例生命周期：

| 类型常量 | Topic | Payload / 构造器 | 除通用字段外的 payload 字段 | 触发场景 |
| --- | --- | --- | --- | --- |
| `EventTypeInstanceCreated` | `approval.instance.created` | `InstanceCreatedEvent`, `NewInstanceCreatedEvent` | 无 | 一个新实例被启动 |
| `EventTypeInstanceCompleted` | `approval.instance.completed` | `InstanceCompletedEvent`, `NewInstanceCompletedEvent` | `finalStatus`、`finishedAt`、`reason`（仅当 `finalStatus` 为 terminated 时才有值） | 实例进入终态 |
| `EventTypeInstanceWithdrawn` | `approval.instance.withdrawn` | `InstanceWithdrawnEvent`, `NewInstanceWithdrawnEvent` | `operator`（`UserInfo`）、`reason` | 申请人撤回实例 |
| `EventTypeInstanceRolledBack` | `approval.instance.rolled_back` | `InstanceRolledBackEvent`, `NewInstanceRolledBackEvent` | `fromNodeId`、`fromNodeName`、`toNodeId`、`toNodeName`、`operator`（`UserInfo`）、`opinion` | 实例回退到之前的节点 |
| `EventTypeInstanceReturned` | `approval.instance.returned` | `InstanceReturnedEvent`, `NewInstanceReturnedEvent` | `fromNodeId`、`fromNodeName`、`toNodeId`、`toNodeName`、`operator`（`UserInfo`）、`opinion` | 实例退回申请人 |
| `EventTypeInstanceResubmitted` | `approval.instance.resubmitted` | `InstanceResubmittedEvent`, `NewInstanceResubmittedEvent` | `operator`（`UserInfo`） | 已退回或已撤回的实例重新提交 |
| `EventTypeInstanceBindingFailed` | `approval.instance.binding_failed` | `InstanceBindingFailedEvent`, `NewInstanceBindingFailedEvent` | `trigger`、`status`、`businessTable`、`errorMessage` | **最终一致**业务投影的一次尝试在期望状态提交后失败；持久化 worker 会持续重试——该事件是运维通知，不是重试机制。同步模式的失败直接回滚审批事务，不会发出此事件 |

节点生命周期：

| 类型常量 | Topic | Payload / 构造器 | 除通用字段外的 payload 字段 | 触发场景 |
| --- | --- | --- | --- | --- |
| `EventTypeNodeAutoPassed` | `approval.node.auto_passed` | `NodeAutoPassedEvent`, `NewNodeAutoPassedEvent` | `nodeId`、`nodeName`、`reason` | 节点因自动通过 execution、空审批人策略或同申请人策略自动通过 |

任务生命周期：

| 类型常量 | Topic | Payload / 构造器 | 除通用字段外的 payload 字段 | 触发场景 |
| --- | --- | --- | --- | --- |
| `EventTypeTaskCreated` | `approval.task.created` | `TaskCreatedEvent`, `NewTaskCreatedEvent` | `assignee`（`UserInfo`）、`deadline`、`status` | 任务创建；顺序审批中的后续任务可能先没有 `deadline`，表示还在等待前置任务 |
| `EventTypeTaskActivated` | `approval.task.activated` | `TaskActivatedEvent`, `NewTaskActivatedEvent` | `assignee`（`UserInfo`）、`deadline`、`reason`（`TaskActivationReason`：`assigned`、`queue_advanced`、`transferred`、`reassigned`） | 任务变为可操作状态——待办通知应订阅此事件。当任务在同一 engine pass 中（如连续审批人连跳）被清除时会被抑制 |
| `EventTypeTaskApproved` | `approval.task.approved` | `TaskApprovedEvent`, `NewTaskApprovedEvent` | `operator`（`UserInfo`）、`opinion` | 任务被批准 |
| `EventTypeTaskHandled` | `approval.task.handled` | `TaskHandledEvent`, `NewTaskHandledEvent` | `operator`（`UserInfo`）、`opinion` | handle 任务完成 |
| `EventTypeTaskRejected` | `approval.task.rejected` | `TaskRejectedEvent`, `NewTaskRejectedEvent` | `operator`（`UserInfo`）、`opinion` | 任务被驳回 |
| `EventTypeTaskCanceled` | `approval.task.canceled` | `TaskCanceledEvent`, `NewTaskCanceledEvent` | `assignee`（`UserInfo`）、`reason` | 引擎取消不再需要决策的任务 |
| `EventTypeTaskTransferred` | `approval.task.transferred` | `TaskTransferredEvent`, `NewTaskTransferredEvent` | `from`（`UserInfo`）、`to`（`UserInfo`）、`reason` | 任务被转交 |
| `EventTypeTaskReassigned` | `approval.task.reassigned` | `TaskReassignedEvent`, `NewTaskReassignedEvent` | `from`（`UserInfo`）、`to`（`UserInfo`）、`reason` | 管理员重新指派任务 |
| `EventTypeTaskTimedOut` | `approval.task.timed_out` | `TaskTimedOutEvent`, `NewTaskTimedOutEvent` | `assignee`（`UserInfo`）、`deadline` | 超时扫描器触发配置的超时动作 |
| `EventTypeAssigneesAdded` | `approval.task.assignees_added` | `AssigneesAddedEvent`, `NewAssigneesAddedEvent` | `addType`、`assignees`（`[]UserInfo`） | 动态新增审批人 |
| `EventTypeAssigneesRemoved` | `approval.task.assignees_removed` | `AssigneesRemovedEvent`, `NewAssigneesRemovedEvent` | `assignees`（`[]UserInfo`） | 动态移除审批人 |
| `EventTypeTaskDeadlineWarning` | `approval.task.deadline_warning` | `TaskDeadlineWarningEvent`, `NewTaskDeadlineWarningEvent` | `assignee`（`UserInfo`）、`deadline`、`hoursLeft` | 预警扫描器命中即将到期任务 |
| `EventTypeTaskUrged` | `approval.task.urged` | `TaskUrgedEvent`, `NewTaskUrgedEvent` | `urger`（`UserInfo`）、`target`（`UserInfo`）、`message` | 申请人或曾持有任务的参与者（办理人/委托人）催办 |

CC + 流程：

| 类型常量 | Topic | Payload / 构造器 | 除通用字段外的 payload 字段 | 触发场景 |
| --- | --- | --- | --- | --- |
| `EventTypeCCNotified` | `approval.cc.notified` | `CCNotifiedEvent`, `NewCCNotifiedEvent` | `nodeId`、`nodeName`、`recipients`（`[]UserInfo`）、`isManual` | CC 节点或手动 CC 动作完成通知 |
| `EventTypeFlowCreated` | `approval.flow.created` | `FlowCreatedEvent`, `NewFlowCreatedEvent` | `categoryId` | 流程创建 |
| `EventTypeFlowUpdated` | `approval.flow.updated` | `FlowUpdatedEvent`, `NewFlowUpdatedEvent` | 无 | 流程更新 |
| `EventTypeFlowDeployed` | `approval.flow.deployed` | `FlowDeployedEvent`, `NewFlowDeployedEvent` | `versionId`、`version` | 流程版本部署 |
| `EventTypeFlowToggled` | `approval.flow.toggled` | `FlowToggledEvent`, `NewFlowToggledEvent` | `isActive` | 流程启停状态变化 |
| `EventTypeFlowPublished` | `approval.flow.published` | `FlowPublishedEvent`, `NewFlowPublishedEvent` | `versionId` | 流程版本发布 |

所有实例和任务事件都从内嵌的 `InstanceEventBase` 携带 `tenantId` 与 `occurredTime`；
`occurredTime` 是业务时间，经 `event.WithOccurredAt` 投影到 `Envelope.OccurredAt`。
`EventTypeTaskCreated` 是 `approval.task.created` 这一 topic 的常量名。

## CallerContext 与多租户安全

资源/命令 handler 在每次请求时解析出一个 `CallerContext`，承载本次调用的租户权能。`TenantID`、`IsSuperAdmin`、`IsSystemInternal` 必须**且只能**有一个为真；零值 `CallerContext` 走 fail-closed 路径（被视为未授权）。

| 字段 | 含义 |
| --- | --- |
| `TenantID` | 调用者所在租户，所有读写都按它过滤 |
| `IsSuperAdmin` | 调用者持有 `approval:super_admin` 角色，允许跨租户访问 |
| `IsSystemInternal` | 框架内代码发起的调用（扫描器、listener） |

跨租户尝试返回 `approval.ErrCrossTenantAccess`。`IsSuperAdmin(p)` 判定 principal 是否带该 override 角色。override 角色字符串本身以 `approval.SuperAdminRole`（`"approval:super_admin"`）形式导出，方便宿主进行权限装配。

## 生命周期扩展点

| 扩展点 | 时机 | 失败语义 |
| --- | --- | --- |
| `approval.InstanceLifecycleHook`（FX group `vef:approval:lifecycle_hooks`） | 同步、引擎事务内部——实例创建与每一次状态迁移 | 返回 error 会回滚整条命令 |
| `approval.BusinessRefProvider` | 同步、start-instance 事务内部 | 返回 error 会回滚实例创建 |
| `approval.BusinessRefResolver` | 同步、start-instance 事务内的目标认领（claim）阶段 | 返回 error 会回滚实例创建 |
| `approval.SubscribeInstance`（`event.SubscribeTyped` 的封装） | 异步、事务提交后 | bus 通过 outbox relay 重试，消费者必须幂等 |
| `approval.BindCommand`（实例事件 → CQRS 命令桥） | 异步、事务提交后 | 重投会重新派发命令，命令 handler 必须幂等 |

### `InstanceLifecycleHook`

```go
type InstanceLifecycleHook interface {
    OnInstanceCreated(ctx context.Context, db orm.DB, instance *Instance) error
    OnInstanceTransition(ctx context.Context, db orm.DB, instance *Instance, from, to InstanceStatus) error
}
```

这个 hook 是面向迁移的通用形态——没有只针对完成的变体。`OnInstanceTransition(instance, from, to)` 在**每一次**实例状态迁移的同一事务内执行——完成（`to.IsFinal()`）、退回、撤回、重新提交、终止。它在引擎拥有的业务投影记录（同步模式下还包括应用）新状态之后运行，所以 hook 观察到的业务表就是这次迁移留下的样子；`instance.Status` 已经是 `to`。返回 error 会回滚整个迁移。

事务内必须成立的不变量（比如分配一个紧耦合的业务行）应该用 lifecycle hook；其他场景都用事件订阅（或下文的 `BindCommand`）。用 `vef.ProvideApprovalLifecycleHook(constructor)` 注册 hook——constructor 必须返回 `approval.InstanceLifecycleHook`，多个 hook 通过 `vef:approval:lifecycle_hooks` group 组合。多个 hook 之间的调用顺序是**未定义的**（FX value group 不携带顺序），所以 hook 必须彼此独立；任何非 nil error 会中止剩余的 hook。`approval.NewFilteredLifecycleHook(hook, filters...)` 用与 `SubscribeInstance`（见下文）相同的 `InstanceFilter` 词汇包装一个 hook，让 hook 可以按 flow code 或 tenant 限定范围，而不必手写谓词。

### `BusinessRefProvider` 和 `BusinessRefResolver`

```go
type BusinessRefProvider interface {
    OnInstanceCreated(ctx context.Context, db orm.DB, flow *Flow, instance *Instance) (businessRef string, err error)
}

type BusinessRefResolver interface {
    ResolveRecordKey(ctx context.Context, flow *Flow, businessRef string) (BusinessRecordKey, error)
}
```

`BusinessRefProvider` 在 `Flow.BindingMode == BindingBusiness` 时提供或分配 opaque 的 `Instance.BusinessRef`；通过 `vef.SupplyBusinessRefProvider` 注入。实例发起时，引擎通过 `BusinessRefResolver` 把 ref 解析成 `approval.BusinessRecordKey`——一个把每个已配置 `KeyColumns` 列映射到取值的 map。默认 resolver 把单列 ref 按原文使用，把多列 ref 按 JSON 对象解码；当 `BusinessRef` 采用别的形状（业务单号、编码元组等）或解析需要宿主查询时，用 `vef.SupplyBusinessRefResolver` 整体替换。无法解析的 ref 会以 `ErrInvalidBusinessRef`（code `40109`）使实例创建失败；业务绑定流程未携带 ref 时返回 `ErrBusinessRefRequired`（code `40107`）。

### 业务状态投影

业务写回是持久化的期望状态投影，而不是逐 trigger 写回。业务绑定的流程配置一份 `BusinessBinding` 文档——`approval.BusinessBindingConfig`，以 jsonb 存在 `apv_flow` 上，并在每次部署时不可变地快照到 `apv_flow_version`；运行期实例只读版本快照：

| 字段 | JSON | 含义 |
| --- | --- | --- |
| `TableName` | `tableName` | 接收审批状态的业务表 |
| `KeyColumns` | `keyColumns` | 记录键列；必须与表上一个非空主键或唯一键完全一致（保存时对照真实 schema 校验——缺表 → `ErrBindingTableMissing`；缺列 → `ErrBindingColumnMissing`；两者均归码 `40018` / `ErrBindingSchemaInvalid`；键列不完整或非唯一 → `ErrBindingKeyNotUnique`） |
| `StatusColumn` | `statusColumn` | 接收（映射后）实例状态的列 |
| `InstanceIDColumn` | `instanceIdColumn` | **必填**——compare-and-set 防护栏，阻止过期实例覆盖新一轮审批拥有的状态 |
| `StartedAtColumn` | `startedAtColumn` | 可选；接收实例发起时间 |
| `FinishedAtColumn` | `finishedAtColumn` | 可选；最终状态时写入结束时间，否则为 `NULL` |
| `StatusMapping` | `statusMapping` | 可选的 `InstanceStatus` → 宿主状态值映射；缺失的条目回退为原始状态字符串（`ErrBindingStatusMappingInvalid` 拒绝未知状态与空白值） |



收敛分三步：

1. **发起时认领。**在 `start_instance` 事务内，引擎解析记录键，并按物理目标认领一行持久投影（`apv_business_projection`，按表 + 记录键唯一）。目标已被一个未终结实例占有时，认领以 `ErrBindingTargetBusy`（code `40108`）失败；已终结的占有者会被顶替，且此前已应用的占有者继续充当业务表的 CAS 防护栏，直到新的期望状态真正写入。实例通过 `Instance.BusinessProjectionID` 关联自己的认领。
2. **每次迁移登记期望状态。**实例的每次状态迁移都会递增投影的 `DesiredRevision` 并记录完整期望状态（状态、开始、结束时间）。投影收敛到*最新*期望状态——中间状态可能被跳过，这正是"绑定到特定生命周期时刻的副作用应放在 `BindCommand` 或 `SubscribeInstance`"的原因。
3. **应用。**`vef.approval.business_binding.consistency` 决定业务行何时写入：
   - `synchronous`（默认）——在审批事务内写入；写入失败会回滚该审批动作。
   - `eventual`——审批动作立即提交；后台 worker（`scan_interval` 默认 `10s`，`batch_size` 默认 `100`）用 `FOR UPDATE SKIP LOCKED` 加租约认领到期投影并应用最新 revision，失败按指数退避重试（`1s << attemptCount`，即 2s、4s……封顶 1 小时）。每次失败尝试都会发布 `InstanceBindingFailedEvent` 作为运维通知——驱动重试的是持久投影行，不是事件。

应用的 `UPDATE` 设置已配置的列,条件为 `WHERE <记录键匹配> [AND <instanceIdColumn> = <已应用占有者>]`，状态值经 `StatusMapping` 翻译。运维侧通过管理端操作 `find_business_projections` / `retry_business_projection` 巡检并强制收敛（见 [RPC 资源](./resources.md#approvaladmin)——每条投影行通过 `BusinessTable` 与 `recordKey` 标识其物理目标），`get_metrics` 报告 `businessProjectionCounts` 与 `pendingBusinessProjections`。

投影状态枚举为 `BindingProjectionStatus`，取值为 `BindingProjectionPending`（已记录新期望状态）、`BindingProjectionProcessing`（worker 已认领并正在应用）、`BindingProjectionApplied`（最新期望状态已写入）和 `BindingProjectionFailed`（最近写入失败，将重试）。

投影归引擎所有；它由配置驱动，不由扩展 hook 替换。宿主用 lifecycle hook、事件订阅和命令桥在它周围叠加行为。

### 宿主订阅：`SubscribeInstance`

```go
func SubscribeInstance[T InstanceEvent](
    bus event.Bus,
    handler func(ctx context.Context, evt T, env event.Envelope) error,
    opts ...InstanceSubscribeOption,
) (event.Unsubscribe, error)
```

`SubscribeInstance` 是 `event.SubscribeTyped` 针对实例事件（任何内嵌 `InstanceEventBase` 的事件类型，例如 `InstanceCompletedEvent`）的声明式封装。路由过滤器是数据，不是谓词：`InstanceSubscribeOption` 接受 `approval.ForFlows(codes...)` / `approval.ForTenants(ids...)`（`InstanceFilter` 值），同一维度内是 OR，传给同一次调用的多个 filter 之间是 AND；未通过某个 filter 的事件会被直接 ack，不会调用 handler。业务谓词（最终状态、表单字段值）应放在 handler 内部，而不是 filter 里。

handler 还会收到投递的 `event.Envelope`：`Envelope.ID` 是 Inbox 去重键，跨重投稳定，因此当路由是 at-least-once 时，它就是构建手动幂等的键。

```go
vef.Invoke(func(bus event.Bus) error {
    _, err := approval.SubscribeInstance(bus,
        func(ctx context.Context, evt *approval.InstanceCompletedEvent, env event.Envelope) error {
            // 幂等的业务副作用；env.ID 跨重投稳定
            return nil
        },
        approval.ForFlows("leave_request"),
        approval.WithGroup("app:leave-writeback"),
    )
    return err
})
```

consumer group 默认由 handler 的方法身份推导（`vef:sub:<相对模块路径包>.<Type>.<method>`，与其他地方使用的 `vef:default:<uuid>` 形状保持一致）。匿名 handler 会因为 `ErrAnonymousSubscriberGroup` 失败，除非提供 `approval.WithGroup(name)`；同一进程内重复推导出相同 group 会因为 `ErrDerivedGroupConflict` 失败。**重命名或移动 handler 会改变它推导出的 group**（会开启一个新的 consumer group，旧的成为孤儿）——在做这类重构之前先用 `WithGroup` 固定 group 名。`approval.WithConcurrency(n)` 设置该订阅的 worker 并发数。

### 命令桥：`BindCommand`

```go
func BindCommand[E InstanceEvent, C cqrs.Action](
    bus event.Bus,
    commands cqrs.Bus,
    mapper func(evt E, env event.Envelope) (cmd C, ok bool),
    opts ...InstanceSubscribeOption,
) (event.Unsubscribe, error)
```

`BindCommand` 订阅实例事件 `E`，并把映射出的命令 `C` 派发到宿主的 CQRS bus——从审批事实到宿主副作用的声明式桥。`mapper` 是纯翻译：用事件加投递 envelope 组装命令，并报告相关性（`ok=false` 直接 ack、不派发）。业务逻辑属于命令 handler，它会经过宿主完整的 behavior 管线（事务、审计、校验）；handler 的返回值被丢弃——派发是 fire-and-record。

```go
vef.Invoke(func(bus event.Bus, commands cqrs.Bus) error {
    _, err := approval.BindCommand(bus, commands,
        func(evt *approval.InstanceCompletedEvent, env event.Envelope) (app.SettleOrderCmd, bool) {
            if evt.FinalStatus != approval.InstanceApproved {
                return app.SettleOrderCmd{}, false
            }
            return app.SettleOrderCmd{InstanceID: evt.InstanceID, DedupeKey: env.ID}, true
        },
        approval.ForFlows("order_settlement"),
    )
    return err
})
```

consumer group 默认取命令类型的身份（`vef:cmd:<相对模块路径包>.<Type>`）——重命名或移动命令会有意地重置订阅键，与重命名 `SubscribeInstance` handler 完全一样；这类重构前先用 `approval.WithGroup` 固定。同一进程内绑定同一命令类型两次会因 `ErrDerivedGroupConflict` 失败。两个前置守卫会当场拒绝非法绑定：`ErrNonCommandAction`（绑定的 action 类型是 query）和 `ErrUnnamedCommandType`（匿名 struct 或接口类型的 `C`）。filter 和 `WithConcurrency` 与 `SubscribeInstance` 用法一致。

投递继承事件路由的语义：在 at-least-once 传输上命令 handler 必须幂等——需要独立去重键时把 `Envelope.ID` 复制进命令。与收敛到最新期望状态、可能跳过中间状态的最终一致业务投影不同,每一次实例迁移都会派发自己的命令，因此 `BindCommand` 是"绑定到特定生命周期时刻的副作用"的通道。

## 主体解析注册表

谁可以审批、谁会被抄送、谁可以发起流程，是三张**开放注册表**，且由同一份描述符驱动。`approval.AssigneeKind`、`CCKind`、`InitiatorKind` 只带着内置常量，并没有 `IsValid`：可部署的种类，恰好就是那些注册了 resolver 的种类。

宿主实现公开契约并注册即可新增一种：

| 契约 | 动作方法 | 注册方式 |
| --- | --- | --- |
| `approval.AssigneeResolver` | `Resolve(ctx, *AssigneeResolveContext) ([]ResolvedAssignee, error)` | `vef.ProvideApprovalAssigneeResolver(...)` |
| `approval.CCResolver` | `Resolve(ctx, *CCResolveContext) ([]string, error)` | `vef.ProvideApprovalCCResolver(...)` |
| `approval.InitiatorResolver` | `Permits(ctx, *InitiatorResolveContext) (bool, error)` | `vef.ProvideApprovalInitiatorResolver(...)` |

三者都还要实现 `Describe() KindDescriptor[K]`。种类与内置项相同的 resolver 会**就地替换**它，并保留设计器里的选项顺序；其他种类则按 kind 升序追加。这个排序是有承重作用的，不是整洁强迫症：fx value group 的到达顺序是随机的，按到达顺序追加会让设计器的下拉框每次重启都重新洗牌。这也是为什么两个宿主 resolver 争夺同一个种类会在启动时被拒绝，而不是"后者胜出"——在随机顺序的 group 里，"后者"是每次启动抛一次硬币决定谁真正在派活。

### 一份描述符同时驱动设计器输入与校验

```go
type KindDescriptor[K ~string] struct {
    Kind      K             `json:"kind"`
    Label     string        `json:"label"`
    Selection SelectionMode `json:"selection"`
}
```

`approval.SelectionMode` 一次性回答了过去要靠硬编码 switch 分别回答的两个问题——这种种类的规则，设计器该渲染哪种输入；保存时校验又该要求哪个配套字段：

| 模式 | 设计器收集 | 校验要求 |
| --- | --- | --- |
| `approval.SelectionNone` | 什么都不收 | 无 |
| `approval.SelectionUser` | 用户 ID | 至少一个 ID |
| `approval.SelectionRole` | 角色 ID | 至少一个 ID |
| `approval.SelectionDepartment` | 部门 ID | 至少一个 ID |
| `approval.SelectionFormField` | 一个表单字段名 | 该表单字段 |
| `approval.SelectionCustom` | 从宿主提供的目录中挑选 ID，按 kind 解析选择器 | 至少一个 ID |

绝大多数"组织性"规则想要的正是 `SelectionNone` 这种形态——"申请人的护士长"、"科主任"——完全从运行时上下文解析出来。`SelectionCustom` 则是专家组那种形态，设计器按种类自身的名字去解析对应的选择器，因此多个自定义种类可以并存而不冲突。

有两条约束是在启动时强制的，而不是留到以后才被发现：描述符有问题（kind 或 label 为空、selection 不在枚举内）会导致注册失败；**发起人**种类永远不允许是 `SelectionFormField`——发起权限是在实例存在之前检查的，因而也在表单存在之前。

### resolver 能看到什么

审批人解析与抄送解析逐字段共用 `approval.NodeResolveContext`，所以为其中一侧写的种类，原封不动就能用在另一侧：

```go
type NodeResolveContext struct {
    Instance     *Instance          // 申请人快照、租户、流程编码、Globals
    Node         *FlowNode          // Key 是设计器指定、跨部署稳定的标识
    FormData     FormData           // 解析时刻的表单数据，而不是数据库行里的那份
    UserResolver UserInfoResolver   // 为裸 ID 补全姓名与部门
}
```

`approval.AssigneeResolveContext` 与 `approval.CCResolveContext` 内嵌它，再加上正在解析的那一条规则——`Kind`、`IDs`、`FormField`。这些模型属于框架，请当作只读。抄送侧刻意**看不到** `CCTiming`：一条规则是在进入时、通过时还是驳回时触发，是引擎的决定，在 resolver 被咨询之前就已经做出；能读到它的 resolver 反而可能与该决定矛盾。

`approval.InitiatorResolveContext` 刻意小得多——只有 `FlowID`、`TenantID`、`Applicant`、`Kind`、`IDs`。这里没有实例也没有表单数据，因为谁可以发起某个流程，是人和流程的属性，绝不是他填了什么的属性。需要看表单的规则应该放进条件分支，那里表单才是真实存在的。另外要注意：申请人的显示**姓名**并不保证有值——发起表单查询会在不解析姓名的情况下检查权限，所以请基于 ID 和部门做判断。

### 三侧的失败语义并不相同

- 审批人 resolver 返回**空**是合法答案——接下来由节点的 `EmptyAssigneeAction` 决定，可选审批环节正是这么表达的。返回 *error* 则会让整个审批动作失败，而不是悄悄丢掉审批人。
- 抄送保持**尽力而为**：无法解析的规则会被记日志并跳过，绝不回滚触发它的那次审批。
- 未注册的**发起人**种类是一个响亮的错误。把它当成"拒绝"，读起来会像权限 bug，而它其实是配置问题；同时发起权限本身是失败关闭的。

### 把目录交给设计器

`approval/flow.list_kind_options` 返回 `approval.KindOptions`——当前运行的应用真正接受的全部审批人、抄送、发起人种类，直接读自启动时注册的 resolver 集合，label 每次调用现算，因此会跟随 `VEF_I18N_LANGUAGE`。

```go
type KindOptions struct {
    Assignees  []KindDescriptor[AssigneeKind]  `json:"assignees"`
    CCs        []KindDescriptor[CCKind]        `json:"ccs"`
    Initiators []KindDescriptor[InitiatorKind] `json:"initiators"`
}
```

它之所以存在，是因为设计器的选项和服务端的校验过去是两份各自手工维护、可能互相打架的列表。请把服务端返回的列表喂给 React 设计器（`EditorPlugins.assigneeKinds` / `ccKinds`、`BasicStep.initiatorKinds`、`FlowValidationContext`）：不传就会回落到内置目录，从而把宿主种类静默地报成未知类型。

## 业务标识符校验

当 `Flow.BindingMode == BindingBusiness` 时，流程会携带 SQL 标识符（`BusinessBindingConfig.TableName`、每个 `KeyColumns` 条目、`StatusColumn`、`InstanceIDColumn` 以及可选的时间戳列），引擎拥有的投影会把它们直接拼到 `UPDATE` 模板里。为防止 SQL 注入，框架按 `^[A-Za-z_][A-Za-z0-9_]{0,62}$` 白名单校验：

```go
if err := approval.ValidateBusinessIdentifier(table); err != nil {
    return err
}
```

空字符串/全空白通过校验——是否把"未配置"算错误由调用方决定。超出白名单的值会返回 `approval.ErrInvalidBusinessIdentifier`。管理端 Flow CRUD 应该把它向上抛，让操作员看到可读的错误。

## 委托

用户可以将审批权限委托给他人：

```go
type Delegation struct {
    DelegatorID    string         // 委托人
    DelegateeID    string         // 被委托人
    FlowCategoryID *string        // 可选：限定分类
    FlowID         *string        // 可选：限定特定流程
    StartsAt       timex.DateTime // 委托窗口开始时间
    EndsAt         timex.DateTime // 委托窗口结束时间
    IsActive       bool
    Reason         *string        // 可选：委托原因
}
```

## 其他公开 API 索引

| 范围 | 公开 API |
| --- | --- |
| caller safety | `CallerContext`, `SystemCaller`, `IsSuperAdmin`, `SuperAdminRole`, `ErrCrossTenantAccess` |
| form data | `FormData`, `NewFormData`, `FormSchemaParser`, `FormFieldDefinition`, `FormSnapshot`, `ValidationRule`, `StorageMode`, `StorageJSON`, `StorageTable`, `FieldKind`, `FieldInput`, `FieldNumber`, `FieldDate`, `FieldTextarea`, `FieldSelect`, `FieldUpload`, `FieldTable`, `FieldOption`, `FieldOptionSource`, `OptionSourceKind`, `OptionSourceRemote`, `RemoteOptionRequest`, `RemoteOptionMapping`, `DynamicParam`, `DynamicParamKind`, `DynamicParamLiteral`, `DynamicParamExpression`, `ColumnDataType`, `ColumnString`, `ColumnText`, `ColumnInteger`, `ColumnDecimal`, `ColumnBoolean`, `ColumnDate`, `ColumnDatetime`, `ColumnJSON` |
| table storage | `FormTable`, `FormTableColumn`（`FormTable.SourceFieldKey` 对主表是 `""`，对明细表格子投影是对应 table 字段的 key） |
| flow models | `FlowCategory`, `Flow`, `FlowVersion`, `FlowNode`, `FlowEdge`, `FlowInitiator`, `FlowNodeAssignee`, `FlowNodeCC`, `VersionStatus`, `VersionDraft`, `VersionPublished`, `VersionArchived`, `ActionLog`, `UrgeRecord`, `DefaultTenantID` |
| node design | `FlowDefinition`, `NodeDefinition`, `EdgeDefinition`, `Position`, `NodeData`, `BaseNodeData`, `StartNodeData`, `ApprovalNodeData`, `HandleNodeData`, `ConditionNodeData`, `CCNodeData`, `EndNodeData`, `ErrUnknownNodeKind`, `ErrNodeDataUnmarshal` |
| conditions | `ConditionKind`, `ConditionField`, `ConditionExpression`, `Condition`, `ConditionGroup`, `ConditionBranch`, `EvaluationContext`, `ConditionEvaluator`, `InstanceGlobalsResolver`, `AggregateKind`, `AggregateSum`, `AggregateCount`, `AggregateAvg`, `Aggregator` |
| initiators and assignees | `InitiatorKind`, `InitiatorUser`, `InitiatorRole`, `InitiatorDepartment`, `AssigneeKind`, `AssigneeDefinition`, `AssigneeService`, `ResolvedAssignee`, `UserInfo`, `UserInfoResolver`, `RoleMembershipChecker`, `AddAssigneeType`, `AddAssigneeBefore`, `AddAssigneeAfter`, `AddAssigneeParallel`, `AssigneeResolver`, `AssigneeResolveContext`, `InitiatorResolver`, `InitiatorResolveContext`, `KindDescriptor`, `KindOptions`, `SelectionMode`, `SelectionNone`, `SelectionUser`, `SelectionRole`, `SelectionDepartment`, `SelectionFormField`, `SelectionCustom`, `NodeResolveContext` |
| CC | `CCKind`, `CCUser`, `CCRole`, `CCDepartment`, `CCFormField`, `CCTiming`, `CCTimingAlways`, `CCTimingOnApprove`, `CCTimingOnReject`, `CCDefinition`, `CCResolver`, `CCResolveContext`, `CCRecord`（`CCRecord.VisitID` 将记录限定在某一次节点 traversal 内，与 `Task.VisitID` 相呼应） |
| business binding & projection | `BusinessBindingConfig`, `BusinessRecordKey`, `BusinessProjection`（模型，表 `apv_business_projection`）, `BindingProjectionStatus`（`BindingProjectionPending` / `Processing` / `Applied` / `Failed`）, `BindingTrigger`, `BindingTriggerStarted`, `BindingTriggerCompleted`, `BindingTriggerReturned`, `BindingTriggerWithdrawn`, `BindingTriggerResubmitted` |
| instance subscriptions | `SubscribeInstance`, `BindCommand`, `InstanceEvent`, `InstanceFilter`, `InstanceSubscribeOption`, `ForFlows`, `ForTenants`, `WithGroup`, `WithConcurrency`, `NewFilteredLifecycleHook`, `ErrAnonymousSubscriberGroup`, `ErrDerivedGroupConflict`, `ErrNonCommandAction`, `ErrUnnamedCommandType` |
| node behavior | `ApprovalMethod`, `TaskNodeData`, `ExecutionType`, `ExecutionManual`, `ExecutionAutoPass`, `ExecutionAutoReject`, `ConsecutiveApproverAction`, `ConsecutiveApproverNone`, `ConsecutiveApproverAutoPass`, `SameApplicantAction`, `SameApplicantSelfApprove`, `SameApplicantAutoPass`, `SameApplicantTransferSuperior`, `SameApplicantExclude`, `Permission`, `PermissionVisible`, `PermissionEditable`, `PermissionRequired`, `PermissionHidden`, `DefaultExecutionType`, `DefaultApprovalMethod`, `DefaultPassRule`, `DefaultEmptyAssigneeAction`, `DefaultSameApplicantAction`, `DefaultConsecutiveApproverAction`, `DefaultRollbackType`, `DefaultRollbackDataStrategy`, `DefaultTimeoutAction`, `DefaultCCTiming`, `DefaultHandleApprovalMethod`, `DefaultHandlePassRule`, `DefaultUrgeCooldownMinutes` |
| rollback and timeouts | `RollbackType`, `RollbackNone`, `RollbackPrevious`, `RollbackStart`, `RollbackAny`, `RollbackSpecified`, `RollbackDataStrategy`, `RollbackDataClear`, `RollbackDataKeep`, `EmptyAssigneeAction`, `EmptyAssigneeAutoPass`, `EmptyAssigneeTransferAdmin`, `EmptyAssigneeTransferSuperior`, `EmptyAssigneeTransferApplicant`, `EmptyAssigneeTransferSpecified`, `TimeoutAction`, `TimeoutActionNone`, `TimeoutActionAutoPass`, `TimeoutActionAutoReject`, `TimeoutActionNotify`, `TimeoutActionTransferAdmin` |
| action and status enums | `ActionType`, `InstanceStatus`, `TaskStatus`, `NodeKind`, `StorageMode`, `VersionStatus` |
| pass rules | `PassRule`, `PassRuleContext`, `PassRuleStrategy`, `PassRuleResult`, `PassRulePending`, `PassRulePassed`, `PassRuleRejected` |
| progress views | `TimelineEntryKind`, `TimelineEntry`, `NodeVisitStatus`, `NodeProgressStatus`, `InstanceFlowGraph`, `FlowGraphNode`, `FlowGraphNodeData`, `FlowGraphEdge`, `NodeParticipant`, `Activity`, `ActivityUrge`, `CCRecipient` |
| events | 所有 `New...Event` 构造器、`DomainEvent`、`InstanceEventBase`、`TaskEventBase`、`FlowEventBase`、`NewInstanceEventBase`、`NewTaskEventBase`、`NewFlowEventBase`、`PayloadOccurredAt`、`AllEventTypes`、`TaskActivationReason`、`TaskActivationAssigned`、`TaskActivationQueueAdvanced`、`TaskActivationTransferred`、`TaskActivationReassigned` 和 `EventType...` 常量 |
| extension interfaces | `InstanceLifecycleHook`, `BusinessRefProvider`, `BusinessRefResolver`, `InstanceNoGenerator`, `ConditionEvaluator`, `InstanceGlobalsResolver`, `PrincipalTenantResolver`, `PrincipalDepartmentResolver`, `RoleMembershipChecker` |
| DI helpers（`vef` 包） | `SupplyBusinessRefProvider`, `SupplyBusinessRefResolver`, `ProvideApprovalLifecycleHook`, `ProvideApprovalAggregator`, `ProvideApprovalFormSchemaParser`, `ProvideApprovalAssigneeResolver`, `ProvideApprovalCCResolver`, `ProvideApprovalInitiatorResolver`, `ProvideApprovalGlobalsResolver` |
| programmatic control | `Service`, `ErrDBRequired`, `StartInstanceInput`, `WithdrawInstanceInput`, `ResubmitInstanceInput`, `TerminateInstanceInput`, `ApproveTaskInput`, `RejectTaskInput`, `TransferTaskInput`, `RollbackTaskInput`, `ReassignTaskInput`, `AddAssigneeInput`, `RemoveAssigneeInput`, `AddCCInput`, `MarkCCReadInput`, `UrgeTaskInput`, `RetryBusinessProjectionInput` —— 参见[编程式流程控制](./service.md) |
| admin DTOs | `approval/admin` 包：`Instance`, `InstanceDetail`, `InstanceDetailInfo`, `Task`, `ActionLog`, `Metrics`, `BusinessProjection` |
| user DTOs | `approval/my` 包：`PendingTask`, `CompletedTask`, `CCRecord`, `InitiatedInstance`, `AvailableFlow`, `InstanceDetail`, `InstanceInfo`, `PendingCounts`, `StartForm`, `ViewerTask`, `RollbackTarget`, `RemovableAssignee` |

---

下一步：回到[概览](./overview.md)查看装配与配置，或到 [RPC 资源](./resources.md) 查看这些事件所来自的 API 面。
