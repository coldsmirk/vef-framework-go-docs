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
auto_migrate = true
outbox_relay_interval = 5
outbox_max_retries = 10
outbox_batch_size = 100
```

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

## 事件发件箱

审批模块使用事务性发件箱模式实现可靠的事件发布。事件在审批操作的同一事务中写入 `apv_event_outbox`，然后异步中继。

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

## 实例编号生成

实现 `InstanceNoGenerator` 接口来自定义实例编号：

```go
type InstanceNoGenerator interface {
    Generate(ctx context.Context, flowCode string) (string, error)
}
```
