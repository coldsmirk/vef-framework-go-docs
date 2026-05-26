---
sidebar_position: 3
---

# 审批模块

`approval` 模块提供完整的工作流引擎，用于构建基于审批的业务流程。支持可视化流程设计（兼容 React Flow）、多级审批链、条件分支、并行审批、委托、回退和事务性事件发布。

## 架构概览

```
流程分类 → 流程 → 流程版本 → 节点 + 边
                               ↓
                           实例 → 任务 → 操作日志
```

| 概念 | 数据表 | 说明 |
| --- | --- | --- |
| 流程分类 | `apv_flow_category` | 流程的层级分组 |
| 流程 | `apv_flow` | 工作流定义（如"请假申请"）|
| 流程版本 | `apv_flow_version` | 版本快照，包含节点、边和表单 schema |
| 流程节点 | `apv_flow_node` | 工作流中的一个步骤 |
| 流程边 | `apv_flow_edge` | 节点之间的有向连接 |
| 实例 | `apv_instance` | 流程的运行实例 |
| 任务 | `apv_task` | 分配给用户的审批/办理任务 |
| 操作日志 | `apv_action_log` | 所有操作的审计追踪 |

## 配置

```toml
[vef.approval]
auto_migrate              = true
timeout_scan_interval     = "1m"
pre_warning_scan_interval = "5m"
cleanup_scan_interval     = "24h"
delegation_max_depth      = 10
form_snapshot_retention   = "2160h"  # 90 天
urge_record_retention     = "720h"   # 30 天
cc_record_retention       = "2160h"  # 90 天
```

> 老版本里归属 `[vef.approval]` 的 `outbox_relay_interval` / `outbox_max_retries` / `outbox_batch_size` 已在 v0.21 迁移至 `[vef.event.transports.outbox]`，由全框架统一的 outbox transport 服务所有模块——参考 [事件总线](../features/event-bus)。审批模块的 binding listener 和 outbox 发布两侧都会在启动时通过 `event.RouteInspector` 断言路由，路由配错时应用直接启动失败，而不是静默降级。

详见[配置参考](../reference/configuration-reference)。

## 绑定模式

| 模式 | 常量 | 说明 |
| --- | --- | --- |
| 独立 | `BindingStandalone` | 表单数据存储在审批模块自有表中 |
| 业务 | `BindingBusiness` | 关联到已有的业务数据表 |

业务绑定通过 `BusinessTable`、`BusinessPkField`、`BusinessTitleField` 和 `BusinessStatusField` 将审批流程与业务表关联。

## 节点类型

| 节点类型 | 常量 | 说明 |
| --- | --- | --- |
| 开始 | `NodeStart` | 工作流入口 |
| 审批 | `NodeApproval` | 需要审批人执行审批动作 |
| 办理 | `NodeHandle` | 需要处理人执行办理动作 |
| 条件 | `NodeCondition` | 基于条件进行分支 |
| 抄送 | `NodeCC` | 向指定用户发送通知 |
| 结束 | `NodeEnd` | 工作流终点 |

## 审批方式

当节点有多个审批人时：

| 方式 | 常量 | 行为 |
| --- | --- | --- |
| 顺序 | `ApprovalSequential` | 审批人按顺序逐个处理 |
| 并行 | `ApprovalParallel` | 审批人同时处理 |

### 通过规则（并行模式）

| 规则 | 常量 | 行为 |
| --- | --- | --- |
| 全部 | `PassAll` | 所有审批人必须同意 |
| 任意 | `PassAny` | 至少一人同意即通过 |
| 比例 | `PassRatio` | 达到一定比例即通过 |
| 一票否决 | `PassAnyReject` | 任何一人拒绝即失败 |

## 审批人类型

| 类型 | 常量 | 说明 |
| --- | --- | --- |
| 指定用户 | `AssigneeUser` | 特定用户 |
| 角色 | `AssigneeRole` | 拥有某角色的用户 |
| 部门 | `AssigneeDepartment` | 部门负责人 |
| 申请人本人 | `AssigneeSelf` | 申请人自己 |
| 直接上级 | `AssigneeSuperior` | 直接上级 |
| 部门领导链 | `AssigneeDepartmentLeader` | 多级主管链 |
| 表单字段 | `AssigneeFormField` | 由表单字段值决定 |

## 实例生命周期

```
提交 → 运行中 → 同意/拒绝 → 已同意/已拒绝
              → 撤回       → 已撤回
              → 回退       → 已退回
              → 终止       → 已终止
              → 重新提交   → 运行中（再次）
```

### 实例状态

| 状态 | 常量 | 是否终态 |
| --- | --- | --- |
| 运行中 | `InstanceRunning` | 否 |
| 已同意 | `InstanceApproved` | 是 |
| 已拒绝 | `InstanceRejected` | 是 |
| 已撤回 | `InstanceWithdrawn` | 否 |
| 已退回 | `InstanceReturned` | 否 |
| 已终止 | `InstanceTerminated` | 是 |

### 任务状态

| 状态 | 常量 | 是否终态 |
| --- | --- | --- |
| 等待中 | `TaskWaiting` | 否 |
| 待处理 | `TaskPending` | 否 |
| 已同意 | `TaskApproved` | 是 |
| 已拒绝 | `TaskRejected` | 是 |
| 已办理 | `TaskHandled` | 是 |
| 已转交 | `TaskTransferred` | 是 |
| 已回退 | `TaskRolledBack` | 是 |
| 已取消 | `TaskCanceled` | 是 |
| 已移除 | `TaskRemoved` | 是 |
| 已跳过 | `TaskSkipped` | 是 |

## 操作类型

| 操作 | 常量 | 说明 |
| --- | --- | --- |
| 提交 | `ActionSubmit` | 发起新实例 |
| 同意 | `ActionApprove` | 审批通过 |
| 办理 | `ActionHandle` | 完成办理任务 |
| 拒绝 | `ActionReject` | 审批拒绝 |
| 转交 | `ActionTransfer` | 转交给其他用户 |
| 撤回 | `ActionWithdraw` | 申请人撤回 |
| 取消 | `ActionCancel` | 取消任务 |
| 回退 | `ActionRollback` | 退回到上一节点 |
| 加签 | `ActionAddAssignee` | 动态添加审批人 |
| 减签 | `ActionRemoveAssignee` | 移除审批人 |
| 重新提交 | `ActionResubmit` | 重新提交已退回的实例 |
| 改派 | `ActionReassign` | 管理员改派任务 |
| 强制终止 | `ActionTerminate` | 管理员强制终止 |

## 回退配置

| 属性 | 可选值 |
| --- | --- |
| `RollbackType` | `none`、`previous`、`start`、`any`、`specified` |
| `RollbackDataStrategy` | `clear`（重置表单）、`keep`（保留数据）|

## 空审批人处理

当节点找不到审批人时：

| 操作 | 常量 |
| --- | --- |
| 自动通过 | `EmptyAssigneeAutoPass` |
| 转交管理员 | `EmptyAssigneeTransferAdmin` |
| 转交上级 | `EmptyAssigneeTransferSuperior` |
| 转交申请人 | `EmptyAssigneeTransferApplicant` |
| 转交指定人 | `EmptyAssigneeTransferSpecified` |

## 超时处理

| 操作 | 常量 | 行为 |
| --- | --- | --- |
| 无操作 | `TimeoutActionNone` | 仅标记超时 |
| 自动通过 | `TimeoutActionAutoPass` | 自动审批通过 |
| 自动拒绝 | `TimeoutActionAutoReject` | 自动审批拒绝 |
| 发送通知 | `TimeoutActionNotify` | 仅发送通知 |
| 转交管理员 | `TimeoutActionTransferAdmin` | 转交给节点管理员 |

## 表单数据存储

| 模式 | 常量 | 存储位置 |
| --- | --- | --- |
| JSON | `StorageJSON` | `apv_instance.form_data`（JSONB 列）|
| 数据表 | `StorageTable` | 动态表 `apv_form_data_{flow_code}` |

## 事件发布

审批模块通过框架统一的事务性 outbox transport 发布领域事件（见 [事件总线](../features/event-bus)）。每条审批命令都把事件记录写在和业务变更相同的事务里，再由 outbox relay 转发到配置的 sink。

订阅方必须：

1. 用 `event.WithGroup("...")` 挂上 —— 路由解析到 at-least-once transport。
2. 依赖 Inbox 中间件去重（在 `event.middleware.inbox = true` 且 transport 声明 `AtLeastOnce` 时自动激活）。

> 老版本里的独立 `apv_event_outbox` 表和 `EventOutboxStatus` 常量已经废弃。审批模块不再自维护一个 outbox，而是组合使用框架统一 outbox。

### 领域事件类型

所有审批事件都实现 `event.Event`，载荷已经包含足够字段来驱动集成而不再回源查询。

实例生命周期：

| 类型常量 | 触发场景 |
| --- | --- |
| `approval.instance.created`（`InstanceCreatedEvent`） | 一个新实例被启动 |
| `approval.instance.completed`（`InstanceCompletedEvent`） | 实例进入终态 |
| `approval.instance.withdrawn`（`InstanceWithdrawnEvent`） | 申请人撤回实例 |
| `approval.instance.rolled_back`（`InstanceRolledBackEvent`） | 实例回退到之前的节点 |
| `approval.instance.returned`（`InstanceReturnedEvent`） | 实例退回申请人 |
| `approval.instance.resubmitted`（`InstanceResubmittedEvent`） | 退回的实例重新提交 |
| `approval.instance.binding_failed`（`InstanceBindingFailedEvent`） | binding listener 未能把终态写回业务行 |

节点生命周期：

| 类型常量 | 触发场景 |
| --- | --- |
| `approval.node.entered` | 引擎激活某节点 |
| `approval.node.auto_passed` | 节点因无审批人自动通过 |

任务生命周期：

| 类型常量 | 触发场景 |
| --- | --- |
| `approval.task.created` | 任务创建（v0.25 起，每条任务创建路径都会发出） |
| `approval.task.approved` | 任务被批准 |
| `approval.task.handled` | handle 任务完成 |
| `approval.task.rejected` | 任务被驳回 |
| `approval.task.transferred` | 任务被转交 |
| `approval.task.reassigned` | 管理员重新指派任务 |
| `approval.task.timed_out` | 超时扫描器触发配置的超时动作 |
| `approval.task.assignees_added` | 动态新增审批人 |
| `approval.task.assignees_removed` | 动态移除审批人 |
| `approval.task.deadline_warning` | 预警扫描器命中即将到期任务 |
| `approval.task.urged` | 申请人催办 |

CC + 流程：

| 类型常量 | 触发场景 |
| --- | --- |
| `approval.cc.notified` | CC 节点完成通知 |
| `approval.flow.created` / `updated` / `deployed` / `toggled` / `published` | 流程设计态生命周期变化 |

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
| `approval.BusinessBindingHook` | 混合：`OnInstanceCreated` 在 start_instance 事务里，`WriteBackStatus` 在 binding listener 异步执行（监听 `InstanceCompletedEvent`） | 同步失败回滚实例创建；异步失败改为发布 `InstanceBindingFailedEvent`，不回滚 |
| 事件订阅（`event.SubscribeTyped`） | 异步、事务提交后 | bus 通过 outbox relay 重试，消费者必须幂等 |

### `InstanceLifecycleHook`

```go
type InstanceLifecycleHook interface {
    OnInstanceCreated(ctx context.Context, db orm.DB, instance *Instance) error
    OnInstanceCompleted(ctx context.Context, db orm.DB, instance *Instance, finalStatus InstanceStatus) error
}
```

事务内必须成立的不变量（比如分配一个紧耦合的业务行）应该用 lifecycle hook；其他场景都用事件订阅。

### `BusinessBindingHook`

```go
type BusinessBindingHook interface {
    OnInstanceCreated(ctx context.Context, db orm.DB, flow *Flow, instance *Instance) (businessRecordID string, err error)
    WriteBackStatus(ctx context.Context, db orm.DB, flow *Flow, instance *Instance, finalStatus InstanceStatus) error
}
```

`Flow.BindingMode == BindingBusiness` 时把审批引擎和宿主业务表桥接起来。通过 `vef.SupplyBusinessBindingHook` 注入。`WriteBackStatus` 在异步路径上由 binding listener 调用——实现必须幂等，因为 outbox relay 可能重试。

## 业务标识符校验

当 `Flow.BindingMode == BindingBusiness` 时，流程会携带 SQL 标识符（`BusinessTable`、`BusinessPkField`、`BusinessTitleField`、`BusinessStatusField`），默认 binding hook 会把它们直接拼到 `UPDATE` 模板里。为防止 SQL 注入，框架按 `^[A-Za-z_][A-Za-z0-9_]{0,62}$` 白名单校验：

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

## 流程定义（兼容 React Flow）

`FlowDefinition` 结构体兼容 React Flow 的 JSON 格式：

```go
type FlowDefinition struct {
    Nodes []NodeDefinition `json:"nodes"`
    Edges []EdgeDefinition `json:"edges"`
}
```

每个 `NodeDefinition` 包含 `Kind` 和强类型的 `Data` 字段，框架会按 `Kind` 把 `Data` 解析成对应结构（`StartNodeData`、`ApprovalNodeData`、`HandleNodeData`、`ConditionNodeData`、`CCNodeData`、`EndNodeData`）。

## 实例编号生成

实现 `InstanceNoGenerator` 接口来自定义实例编号：

```go
type InstanceNoGenerator interface {
    Generate(ctx context.Context, flowCode string) (string, error)
}
```
