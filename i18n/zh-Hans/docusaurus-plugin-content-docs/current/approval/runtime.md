---
sidebar_position: 4
---

# 实例运行时

## 实例生命周期

```
提交 → 运行中 → 同意/拒绝 → 已同意/已拒绝
              → 撤回       → 已撤回
              → 回退       → 已退回
              → 终止       → 已终止
已撤回/已退回 → 重新提交   → 运行中（再次）
```

运行时状态机只声明以下实例状态流转：

| From | To |
| --- | --- |
| `running` | `approved` |
| `running` | `rejected` |
| `running` | `withdrawn` |
| `running` | `terminated` |
| `running` | `returned` |
| `returned` | `running` |
| `returned` | `terminated` |
| `returned` | `withdrawn` |
| `withdrawn` | `running` |
| `withdrawn` | `terminated` |

### 实例状态

| 状态 | 常量 | Wire value | 是否终态 |
| --- | --- | --- | --- |
| 运行中 | `InstanceRunning` | `running` | 否 |
| 已同意 | `InstanceApproved` | `approved` | 是 |
| 已拒绝 | `InstanceRejected` | `rejected` | 是 |
| 已撤回 | `InstanceWithdrawn` | `withdrawn` | 否 |
| 已退回 | `InstanceReturned` | `returned` | 否 |
| 已终止 | `InstanceTerminated` | `terminated` | 是 |

枚举类型是 `InstanceStatus`。

### 任务状态

| 状态 | 常量 | Wire value | 是否终态 |
| --- | --- | --- | --- |
| 等待中 | `TaskWaiting` | `waiting` | 否 |
| 待处理 | `TaskPending` | `pending` | 否 |
| 已同意 | `TaskApproved` | `approved` | 是 |
| 已拒绝 | `TaskRejected` | `rejected` | 是 |
| 已办理 | `TaskHandled` | `handled` | 是 |
| 已转交 | `TaskTransferred` | `transferred` | 是 |
| 已回退 | `TaskRolledBack` | `rolled_back` | 是 |
| 已取消 | `TaskCanceled` | `canceled` | 是 |
| 已移除 | `TaskRemoved` | `removed` | 是 |
| 已跳过 | `TaskSkipped` | `skipped` | 是 |

运行时状态机只声明以下任务状态流转：

| From | To |
| --- | --- |
| `waiting` | `pending` |
| `waiting` | `canceled` |
| `waiting` | `skipped` |
| `waiting` | `removed` |
| `pending` | `approved` |
| `pending` | `handled` |
| `pending` | `rejected` |
| `pending` | `transferred` |
| `pending` | `rolled_back` |
| `pending` | `canceled` |
| `pending` | `waiting` |
| `pending` | `removed` |

## 操作类型

| 操作 | 常量 | Wire value | 说明 |
| --- | --- | --- | --- |
| 提交 | `ActionSubmit` | `submit` | 发起新实例 |
| 同意 | `ActionApprove` | `approve` | 审批通过 |
| 办理 | `ActionHandle` | `handle` | 完成办理任务 |
| 拒绝 | `ActionReject` | `reject` | 审批拒绝 |
| 转交 | `ActionTransfer` | `transfer` | 转交给其他用户 |
| 撤回 | `ActionWithdraw` | `withdraw` | 申请人撤回 |
| 取消 | `ActionCancel` | `cancel` | 取消任务 |
| 回退 | `ActionRollback` | `rollback` | 退回到上一节点 |
| 加签 | `ActionAddAssignee` | `add_assignee` | 动态添加审批人 |
| 减签 | `ActionRemoveAssignee` | `remove_assignee` | 移除审批人 |
| 加抄送 | `ActionAddCC` | `add_cc` | 动态添加抄送人 |
| 执行 | `ActionExecute` | `execute` | 自动节点处理使用的内部执行动作 |
| 重新提交 | `ActionResubmit` | `resubmit` | 重新提交已退回或已撤回的实例 |
| 改派 | `ActionReassign` | `reassign` | 管理员改派任务 |
| 强制终止 | `ActionTerminate` | `terminate` | 管理员强制终止 |

## 回退配置

| 属性 | 可选值 |
| --- | --- |
| `RollbackType` | `RollbackNone`（`none`）、`RollbackPrevious`（`previous`）、`RollbackStart`（`start`）、`RollbackAny`（`any`）、`RollbackSpecified`（`specified`） |
| `RollbackDataStrategy` | `RollbackDataClear`（`clear`，重置表单）、`RollbackDataKeep`（`keep`，保留数据）|

同一申请人处理使用 `SameApplicantAction`，取值包括
`SameApplicantSelfApprove`（`self_approve`）、`SameApplicantAutoPass`
（`auto_pass`）、`SameApplicantTransferSuperior`（`transfer_superior`）。
连续审批人处理使用 `ConsecutiveApproverAction`，取值包括
`ConsecutiveApproverNone`（`none`）与 `ConsecutiveApproverAutoPass`
（`auto_pass`）。

## 空审批人处理

当节点找不到审批人时：

| 操作 | 常量 | Wire value |
| --- | --- | --- |
| 自动通过 | `EmptyAssigneeAutoPass` | `auto_pass` |
| 转交管理员 | `EmptyAssigneeTransferAdmin` | `transfer_admin` |
| 转交上级 | `EmptyAssigneeTransferSuperior` | `transfer_superior` |
| 转交申请人 | `EmptyAssigneeTransferApplicant` | `transfer_applicant` |
| 转交指定人 | `EmptyAssigneeTransferSpecified` | `transfer_specified` |

枚举类型是 `EmptyAssigneeAction`。

## 超时处理

| 操作 | 常量 | Wire value | 行为 |
| --- | --- | --- | --- |
| 无操作 | `TimeoutActionNone` | `none` | 仅标记超时 |
| 自动通过 | `TimeoutActionAutoPass` | `auto_pass` | 自动审批通过 |
| 自动拒绝 | `TimeoutActionAutoReject` | `auto_reject` | 自动审批拒绝 |
| 发送通知 | `TimeoutActionNotify` | `notify` | 仅发送通知 |
| 转交管理员 | `TimeoutActionTransferAdmin` | `transfer_admin` | 转交给节点管理员 |

枚举类型是 `TimeoutAction`。

## 表单数据存储

| 模式 | 常量 | Wire value | 存储位置 |
| --- | --- | --- | --- |
| JSON | `StorageJSON` | `json` | `apv_instance.form_data`（JSONB 列）|
| Table | `StorageTable` | `table` | 每个已发布版本一张生成的物理表，同时继续写入 `apv_instance.form_data` 作为 canonical JSON snapshot |

`StorageMode.IsValid()` 接受这两种导出模式。Table 模式通过 `FormTable` 和 `FormTableColumn` 记录生成的 DDL 元数据。

### 表存储元数据

当已发布版本使用 `StorageTable` 时，框架公开两个元数据模型：

| 模型 | 用途 | 关键 JSON 字段 |
| --- | --- | --- |
| `FormTable` | 一张生成的物理表（主投影表或明细表格子表） | `flowId`、`versionId`、`physicalTableName`、`sourceFieldKey` |
| `FormTableColumn` | 每个表单字段或内置列对应一个生成列 | `formTableId`、`columnName`、`columnType`、`isNullable`、`sourceFieldKey`、`sortOrder` |

`FormTable` 是框架生成 DDL 的单一事实来源：引擎在建表前会先查它保证幂等，
操作员也可以通过它把版本映射到对应的投影表。`FormTable.SourceFieldKey`
对版本的主投影表是 `""`，对明细表格子投影是所属 table 字段的 key；
`(versionId, sourceFieldKey)` 唯一。`ColumnDataType` 是表单定义使用的逻辑
字段到列类型词汇，storage 层再把它映射成具体数据库方言的 SQL 类型。

### 生成的物理表布局

发布 table 模式的版本会创建一张主投影表，外加每个明细表格（`table` 类型）
字段一张子表（v0.36）：

| 表 | 物理表名 | 内置列 | 字段列 |
| --- | --- | --- | --- |
| 主投影表 | `apv_form_<code>_<versionId>`（flow code 经净化处理并按 63 字符标识符上限截断；预算不足时退化为 `apv_form_<versionId>`） | `id`（主键）、`instance_id`（UNIQUE）、`created_at` 排最后 | 每个标量字段一列，按声明顺序 |
| 明细表格子表 | `apv_form_<versionId>__<fieldKey 后缀>` | `id`（主键）、`instance_id`（有索引、非唯一）、`row_index`、`created_at` 排最后 | 每个表格列一列，按声明顺序 |

每个物理表名和列名进入 DDL/DML 字符串前都要先通过安全 SQL 标识符校验，
字段 key 不能和内置列名冲突，所有表单值都以 `?` 参数绑定——从不拼接。
`CREATE TABLE IF NOT EXISTS` 在发布事务之外执行（MySQL 的 DDL 会隐式提交
进行中的事务）且幂等；`FormTable` / `FormTableColumn` 元数据则在发布事务内
写入，随版本的发布状态一起提交或回滚。

写入时投影采用"替换而非追加"策略：主表每个实例恰好一行（由 `instance_id`
UNIQUE 约束保证），子表按 `row_index` 排序、每条明细一行；每次投影（启动时
和每次重新提交时）都会先删除该实例已有的行再插入新行——物理表始终反映当前
表单数据，不积累历史。

## 实例进度投影

管理端和用户端 instance detail 响应会在实例快照旁公开两个只读投影：

| 投影 | 公开类型 | 用途 |
| --- | --- | --- |
| timeline | `TimelineEntryKind`、`TimelineEntry`、`NodeVisitStatus`、`NodeParticipant`、`Activity`、`ActivityUrge`、`CCRecipient` | 实例实际经过路径的时间线 |
| flow graph | `NodeProgressStatus`、`InstanceFlowGraph`、`FlowGraphNode`、`FlowGraphNodeData`、`FlowGraphEdge` | 带运行时进度标注、兼容 React Flow 的流程图 |

`FlowGraphNode.ID` 是 React Flow 设计期节点 id。`FlowGraphNode.NodeID` 是 action log 和 rollback target 使用的持久化 flow-node id。

进度和时间线枚举显式导出：

| 枚举 | 常量 |
| --- | --- |
| `TimelineEntryKind` | `TimelineEntryStart`、`TimelineEntryApproval`、`TimelineEntryHandle`、`TimelineEntryCC`、`TimelineEntryWithdraw`、`TimelineEntryTerminate` |
| `NodeVisitStatus` | `NodeVisitActive`、`NodeVisitPassed`、`NodeVisitRejected`、`NodeVisitReturned`、`NodeVisitCanceled` |
| `NodeProgressStatus` | `NodeProgressPending`、`NodeProgressActive`、`NodeProgressPassed`、`NodeProgressRejected`、`NodeProgressReturned`、`NodeProgressCanceled` |

这些投影背后的持久化 visit 模型是 `NodeVisit`。

## 实例编号生成

实现 `InstanceNoGenerator` 接口来自定义实例编号：

```go
type InstanceNoGenerator interface {
    Generate(ctx context.Context, flowCode string) (string, error)
}
```

---

下一步：[事件与集成](./integration.md) 了解如何在宿主代码中响应生命周期流转。
