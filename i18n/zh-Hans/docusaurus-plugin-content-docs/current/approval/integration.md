---
sidebar_position: 5
---

# 事件与集成

## 事件发布

审批模块通过框架统一的事务性 outbox transport 发布领域事件（见 [事件总线](../infrastructure/event-bus.md)）。每条审批命令都把事件记录写在和业务变更相同的事务里，再由 outbox relay 转发到配置的 sink。

订阅方必须：

1. 用 `event.WithGroup("...")` 挂上 —— 路由解析到 at-least-once transport。
2. 依赖 Inbox 中间件去重（在 `event.middleware.inbox = true` 且 transport 声明 `AtLeastOnce` 时自动激活）。

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
| `EventTypeInstanceBindingFailed` | `approval.instance.binding_failed` | `InstanceBindingFailedEvent`, `NewInstanceBindingFailedEvent` | `trigger`、`status`、`businessTable`、`errorMessage` | 引擎拥有的写回在 `trigger` 之后未能把 `status` 持久化到业务行 |

节点生命周期：

| 类型常量 | Topic | Payload / 构造器 | 除通用字段外的 payload 字段 | 触发场景 |
| --- | --- | --- | --- | --- |
| `EventTypeNodeAutoPassed` | `approval.node.auto_passed` | `NodeAutoPassedEvent`, `NewNodeAutoPassedEvent` | `nodeId`、`nodeName`、`reason` | 节点因自动通过 execution、空审批人策略或同申请人策略自动通过 |

任务生命周期：

| 类型常量 | Topic | Payload / 构造器 | 除通用字段外的 payload 字段 | 触发场景 |
| --- | --- | --- | --- | --- |
| `EventTypeTaskCreated` | `approval.task.created` | `TaskCreatedEvent`, `NewTaskCreatedEvent` | `assignee`（`UserInfo`）、`deadline` | 任务创建；顺序审批中的后续任务可能先没有 `deadline`，表示还在等待前置任务 |
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
| `EventTypeTaskUrged` | `approval.task.urged` | `TaskUrgedEvent`, `NewTaskUrgedEvent` | `urger`（`UserInfo`）、`target`（`UserInfo`）、`message` | 申请人催办 |

CC + 流程：

| 类型常量 | Topic | Payload / 构造器 | 除通用字段外的 payload 字段 | 触发场景 |
| --- | --- | --- | --- | --- |
| `EventTypeCCNotified` | `approval.cc.notified` | `CCNotifiedEvent`, `NewCCNotifiedEvent` | `nodeId`、`nodeName`、`recipients`（`[]UserInfo`）、`isManual` | CC 节点或手动 CC 动作完成通知 |
| `EventTypeFlowCreated` | `approval.flow.created` | `FlowCreatedEvent`, `NewFlowCreatedEvent` | `categoryId` | 流程创建 |
| `EventTypeFlowUpdated` | `approval.flow.updated` | `FlowUpdatedEvent`, `NewFlowUpdatedEvent` | 无 | 流程更新 |
| `EventTypeFlowDeployed` | `approval.flow.deployed` | `FlowDeployedEvent`, `NewFlowDeployedEvent` | `versionId`、`version` | 流程版本部署 |
| `EventTypeFlowToggled` | `approval.flow.toggled` | `FlowToggledEvent`, `NewFlowToggledEvent` | `isActive` | 流程启停状态变化 |
| `EventTypeFlowPublished` | `approval.flow.published` | `FlowPublishedEvent`, `NewFlowPublishedEvent` | `versionId` | 流程版本发布 |

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
| `approval.InstanceLifecycleHook`（FX group `vef:approval:lifecycle_hooks`） | 同步、业务事务内部 | 返回 error 会回滚整条命令 |
| `approval.BusinessRefProvider` | 同步、start-instance 事务内部 | 返回 error 会回滚实例创建 |
| `approval.BusinessRefResolver` | 异步写回路径（`completed` / `returned` / `withdrawn` / `resubmitted` 四种 trigger） | 返回 error 会让引擎拥有的写回失败，发布 `InstanceBindingFailedEvent`，并由 outbox 路径重试 |
| `approval.SubscribeInstance`（`event.SubscribeTyped` 的封装） | 异步、事务提交后 | bus 通过 outbox relay 重试，消费者必须幂等 |

### `InstanceLifecycleHook`

```go
type InstanceLifecycleHook interface {
    OnInstanceCreated(ctx context.Context, db orm.DB, instance *Instance) error
    OnInstanceCompleted(ctx context.Context, db orm.DB, instance *Instance, finalStatus InstanceStatus) error
}
```

事务内必须成立的不变量（比如分配一个紧耦合的业务行）应该用 lifecycle hook；其他场景都用事件订阅。用 `vef.ProvideApprovalLifecycleHook(constructor)` 注册 hook——constructor 必须返回 `approval.InstanceLifecycleHook`，多个 hook 通过 `vef:approval:lifecycle_hooks` group 组合。`approval.NewFilteredLifecycleHook(hook, filters...)` 用与 `SubscribeInstance`（见下文）相同的 `InstanceFilter` 词汇包装一个 hook，让 hook 可以按 flow code 或 tenant 限定范围，而不必手写谓词。

### `BusinessRefProvider` 和 `BusinessRefResolver`

```go
type BusinessRefProvider interface {
    OnInstanceCreated(ctx context.Context, db orm.DB, flow *Flow, instance *Instance) (businessRef string, err error)
}

type BusinessRefResolver interface {
    ResolveRecordID(ctx context.Context, flow *Flow, businessRef string) (string, error)
}
```

`BusinessRefProvider` 在 `Flow.BindingMode == BindingBusiness` 时提供或分配 opaque 的 `Instance.BusinessRef`；通过 `vef.SupplyBusinessRefProvider` 注入。引擎不会直接解析 `BusinessRef`。写回时，它会询问 `BusinessRefResolver`（默认是 identity resolver）提取要和 `Flow.BusinessPKField` 匹配的 record id；复合 ref 形状通过 `vef.SupplyBusinessRefResolver` 注入自定义解析。

### 业务写回联动矩阵

业务绑定的流程配置 `Flow.BusinessTable` / `BusinessPKField` / `BusinessStatusField`（必填），以及三个可选列——`BusinessInstanceIDField`、`BusinessStartedAtField`、`BusinessFinishedAtField`（Go 字段名；JSON 仍保持 `businessPkField` 等命名）。为 `nil` 表示"永远不动那一列"。引擎拥有的写回（`binding.Writer.WriteBack`）会按 `approval.BindingTrigger` 投影实例的当前状态（v0.37）：

| `BindingTrigger` | 状态列 | 实例 ID 列 | 开始时间列 | 结束时间列 |
| --- | --- | --- | --- | --- |
| `BindingTriggerStarted`（`started`） | `running` | 实例的 id | 当前时间 | 清空（`NULL`） |
| `BindingTriggerCompleted`（`completed`） | 最终状态 | — | — | `Instance.FinishedAt` |
| `BindingTriggerReturned`（`returned`） | `returned` | — | — | — |
| `BindingTriggerWithdrawn`（`withdrawn`） | `withdrawn` | — | — | — |
| `BindingTriggerResubmitted`（`resubmitted`） | `running` | — | — | 清空（`NULL`） |

`started` 投影在 `start_instance` 事务内同步执行——失败会回滚整个发起流程。其余四种 trigger 通过 binding listener 异步执行；失败会发布 `InstanceBindingFailedEvent`（携带 `trigger` 和 `status`），并通过 outbox 路径重试。

审批引擎拥有状态写回（`UPDATE businessTable SET businessStatusField = ? WHERE businessPkField = ?`，加上上面的可选列）。宿主可以用 lifecycle hook 或事件订阅围绕它扩展，但写回本身由配置驱动，不再由扩展 hook 替换。

### 宿主订阅：`SubscribeInstance`

```go
func SubscribeInstance[T InstanceEvent](
    bus event.Bus,
    handler func(ctx context.Context, evt T) error,
    opts ...InstanceSubscribeOption,
) (event.Unsubscribe, error)
```

`SubscribeInstance`（v0.37）是 `event.SubscribeTyped` 针对实例事件（任何内嵌 `InstanceEventBase` 的事件类型，例如 `InstanceCompletedEvent`）的声明式封装。路由过滤器是数据，不是谓词：`InstanceSubscribeOption` 接受 `approval.ForFlows(codes...)` / `approval.ForTenants(ids...)`（`InstanceFilter` 值），同一维度内是 OR，传给同一次调用的多个 filter 之间是 AND；未通过某个 filter 的事件会被直接 ack，不会调用 handler。业务谓词（最终状态、表单字段值）应放在 handler 内部，而不是 filter 里。

```go
vef.Invoke(func(bus event.Bus) error {
    _, err := approval.SubscribeInstance(bus,
        func(ctx context.Context, evt *approval.InstanceCompletedEvent) error {
            // 幂等的业务副作用；失败时由 outbox relay 重投
            return nil
        },
        approval.ForFlows("leave_request"),
        approval.WithGroup("app:leave-writeback"),
    )
    return err
})
```

consumer group 默认由 handler 的方法身份推导（`vef:sub:<相对模块路径包>.<Type>.<method>`，与其他地方使用的 `vef:default:<uuid>` 形状保持一致）。匿名 handler 会因为 `ErrAnonymousSubscriberGroup` 失败，除非提供 `approval.WithGroup(name)`；同一进程内重复推导出相同 group 会因为 `ErrDerivedGroupConflict` 失败。**重命名或移动 handler 会改变它推导出的 group**（会开启一个新的 consumer group，旧的成为孤儿）——在做这类重构之前先用 `WithGroup` 固定 group 名。`approval.WithConcurrency(n)` 设置该订阅的 worker 并发数。

## 业务标识符校验

当 `Flow.BindingMode == BindingBusiness` 时，流程会携带 SQL 标识符（`BusinessTable`、`BusinessPKField`、`BusinessStatusField`，以及上面的三个可选联动列），引擎拥有的写回会把它们直接拼到 `UPDATE` 模板里。为防止 SQL 注入，框架按 `^[A-Za-z_][A-Za-z0-9_]{0,62}$` 白名单校验：

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
    StartTime      timex.DateTime // 委托开始时间
    EndTime        timex.DateTime // 委托结束时间
    IsActive       bool
}
```

## 其他公开 API 索引

| 范围 | 公开 API |
| --- | --- |
| caller safety | `CallerContext`, `SystemCaller`, `IsSuperAdmin`, `SuperAdminRole`, `ErrCrossTenantAccess` |
| form data | `FormData`, `NewFormData`, `FormDefinition`, `FormFieldDefinition`, `FormSnapshot`, `ValidationRule`, `StorageMode`, `StorageJSON`, `StorageTable`, `FieldKind`, `FieldInput`, `FieldNumber`, `FieldDate`, `FieldTextarea`, `FieldSelect`, `FieldUpload`, `FieldTable`, `FieldOption`, `ColumnDataType`, `ColumnString`, `ColumnText`, `ColumnInteger`, `ColumnDecimal`, `ColumnBoolean`, `ColumnDate`, `ColumnDatetime`, `ColumnJSON` |
| table storage | `FormTable`, `FormTableColumn`（`FormTable.SourceFieldKey` 对主表是 `""`，对明细表格子投影是对应 table 字段的 key） |
| flow models | `FlowCategory`, `Flow`, `FlowVersion`, `FlowNode`, `FlowEdge`, `FlowInitiator`, `FlowNodeAssignee`, `FlowNodeCC`, `VersionStatus`, `VersionDraft`, `VersionPublished`, `VersionArchived`, `ActionLog`, `OperatorInfo`, `UrgeRecord`, `DefaultTenantID` |
| node design | `FlowDefinition`, `NodeDefinition`, `EdgeDefinition`, `Position`, `NodeData`, `BaseNodeData`, `StartNodeData`, `ApprovalNodeData`, `HandleNodeData`, `ConditionNodeData`, `CCNodeData`, `EndNodeData`, `ErrUnknownNodeKind`, `ErrNodeDataUnmarshal` |
| conditions | `ConditionKind`, `ConditionField`, `ConditionExpression`, `Condition`, `ConditionGroup`, `ConditionBranch`, `EvaluationContext`, `ConditionEvaluator`, `InstanceGlobalsResolver`, `AggregateKind`, `AggregateSum`, `AggregateCount`, `AggregateAvg`, `Aggregator` |
| initiators and assignees | `InitiatorKind`, `InitiatorUser`, `InitiatorRole`, `InitiatorDepartment`, `AssigneeKind`, `AssigneeDefinition`, `AssigneeService`, `ResolvedAssignee`, `UserInfo`, `UserInfoResolver`, `RoleMembershipChecker`, `AddAssigneeType`, `AddAssigneeBefore`, `AddAssigneeAfter`, `AddAssigneeParallel` |
| CC | `CCKind`, `CCUser`, `CCRole`, `CCDepartment`, `CCFormField`, `CCTiming`, `CCTimingAlways`, `CCTimingOnApprove`, `CCTimingOnReject`, `CCDefinition`, `CCRecord`（`CCRecord.VisitID` 将记录限定在某一次节点 traversal 内，与 `Task.VisitID` 相呼应） |
| business write-back | `BindingTrigger`, `BindingTriggerStarted`, `BindingTriggerCompleted`, `BindingTriggerReturned`, `BindingTriggerWithdrawn`, `BindingTriggerResubmitted` |
| instance subscriptions | `SubscribeInstance`, `InstanceEvent`, `InstanceFilter`, `InstanceSubscribeOption`, `ForFlows`, `ForTenants`, `WithGroup`, `WithConcurrency`, `NewFilteredLifecycleHook`, `ErrAnonymousSubscriberGroup`, `ErrDerivedGroupConflict` |
| node behavior | `ApprovalMethod`, `TaskNodeData`, `ExecutionType`, `ExecutionManual`, `ExecutionAutoPass`, `ExecutionAutoReject`, `ConsecutiveApproverAction`, `ConsecutiveApproverNone`, `ConsecutiveApproverAutoPass`, `SameApplicantAction`, `SameApplicantSelfApprove`, `SameApplicantAutoPass`, `SameApplicantTransferSuperior`, `Permission`, `PermissionVisible`, `PermissionEditable`, `PermissionRequired`, `PermissionHidden`, `DefaultExecutionType`, `DefaultApprovalMethod`, `DefaultPassRule`, `DefaultEmptyAssigneeAction`, `DefaultSameApplicantAction`, `DefaultConsecutiveApproverAction`, `DefaultRollbackType`, `DefaultRollbackDataStrategy`, `DefaultTimeoutAction`, `DefaultCCTiming`, `DefaultHandleApprovalMethod`, `DefaultHandlePassRule`, `DefaultUrgeCooldownMinutes` |
| rollback and timeouts | `RollbackType`, `RollbackNone`, `RollbackPrevious`, `RollbackStart`, `RollbackAny`, `RollbackSpecified`, `RollbackDataStrategy`, `RollbackDataClear`, `RollbackDataKeep`, `EmptyAssigneeAction`, `EmptyAssigneeAutoPass`, `EmptyAssigneeTransferAdmin`, `EmptyAssigneeTransferSuperior`, `EmptyAssigneeTransferApplicant`, `EmptyAssigneeTransferSpecified`, `TimeoutAction`, `TimeoutActionNone`, `TimeoutActionAutoPass`, `TimeoutActionAutoReject`, `TimeoutActionNotify`, `TimeoutActionTransferAdmin` |
| action and status enums | `ActionType`, `InstanceStatus`, `TaskStatus`, `NodeKind`, `StorageMode`, `VersionStatus` |
| pass rules | `PassRule`, `PassRuleContext`, `PassRuleStrategy`, `PassRuleResult`, `PassRulePending`, `PassRulePassed`, `PassRuleRejected` |
| progress views | `TimelineEntryKind`, `TimelineEntry`, `NodeVisitStatus`, `NodeProgressStatus`, `InstanceFlowGraph`, `FlowGraphNode`, `FlowGraphNodeData`, `FlowGraphEdge`, `NodeParticipant`, `Activity`, `ActivityUrge`, `CCRecipient` |
| events | 所有 `New...Event` 构造器、`DomainEvent`、`InstanceEventBase`、`TaskEventBase`、`FlowEventBase`、`NewInstanceEventBase`、`NewTaskEventBase`、`NewFlowEventBase`、`PayloadOccurredAt`、`AllEventTypes` 和 `EventType...` 常量 |
| extension interfaces | `InstanceLifecycleHook`, `BusinessRefProvider`, `BusinessRefResolver`, `InstanceNoGenerator`, `ConditionEvaluator`, `InstanceGlobalsResolver`, `PrincipalTenantResolver`, `PrincipalDepartmentResolver`, `RoleMembershipChecker` |
| DI helpers（`vef` 包） | `SupplyBusinessRefProvider`, `SupplyBusinessRefResolver`, `ProvideApprovalLifecycleHook`, `ProvideApprovalAggregator` |
| admin DTOs | `approval/admin` 包：`Instance`, `InstanceDetail`, `InstanceDetailInfo`, `Task`, `ActionLog`, `Metrics` |
| user DTOs | `approval/my` 包：`PendingTask`, `CompletedTask`, `CCRecord`, `InitiatedInstance`, `AvailableFlow`, `InstanceDetail`, `InstanceInfo`, `PendingCounts` |

---

下一步：回到[概览](./overview.md)查看装配与配置，或到 [RPC 资源](./resources.md) 查看这些事件所来自的 API 面。
