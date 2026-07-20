---
sidebar_position: 3
---

# 流程设计

## 节点类型

| 节点类型 | 常量 | Wire value | 说明 |
| --- | --- | --- | --- |
| 开始 | `NodeStart` | `start` | 工作流入口 |
| 审批 | `NodeApproval` | `approval` | 需要审批人执行审批动作 |
| 办理 | `NodeHandle` | `handle` | 需要处理人执行办理动作 |
| 条件 | `NodeCondition` | `condition` | 基于条件进行分支 |
| 抄送 | `NodeCC` | `cc` | 向指定用户发送通知 |
| 结束 | `NodeEnd` | `end` | 工作流终点 |

## 条件分支

条件节点会按 priority 顺序评估 `ConditionBranch`。每个分支包含一个或多个
`ConditionGroup`：同一个 group 内的条件按 AND 组合，同一分支上的多个
group 按 OR 组合。

`ConditionField` 使用结构化的 `Subject` / `Operator` / `Value` 字段。
`Operator` 类型是 `ConditionOperator`；公开常量包括 `OperatorEquals`、
`OperatorNotEquals`、`OperatorGreater`、`OperatorGreaterOrEq`、`OperatorLess`、
`OperatorLessOrEq`、`OperatorIn`、`OperatorNotIn`、`OperatorContains`、
`OperatorNotContains`、`OperatorStartsWith`、`OperatorEndsWith`、
`OperatorIsEmpty` 和 `OperatorIsNotEmpty`。内置 evaluator 以**原生 Go 比较**
求值 field condition——不做表达式模板化、没有注入面——操作符语义按字段类型
定型；词汇表之外的操作符会以错误终止求值，而不是静默评为 `false`。

`ConditionExpression` 执行原始 `Expression` 字符串。
评估环境暴露：

| 名称 | 值 |
| --- | --- |
| `formData` | 当前实例的 `FormData` map |
| `applicantId` | 当前申请人 ID |
| `applicantDepartmentId` | 申请人部门 ID；不存在时是 `""` |
| globals | 宿主解析出的 `Instance.Globals`，作为顶层 binding 暴露 |

表达式条件经框架的 `expression.Engine` 抽象求值（当前由 `expr-lang`
支撑），由 DI 装配——即[表达式引擎](../data-tools/expression)文档描述的
同一引擎。

宿主应用可以实现 `approval.InstanceGlobalsResolver`，在实例启动时根据已认证
principal 解析全局变量。该快照会持久化到 `Instance.Globals`；客户端不能在
`start` 请求体里提交它。Field condition 会先从 globals 解析 `Subject`，再查
`formData`；expression condition 会把 globals 暴露为顶层 binding，但内置的
`formData`、`applicantId`、`applicantDepartmentId` 名称发生冲突时优先级更高。

### 明细表格聚合

字段条件可以不比较标量 subject，而是对明细表格的行做聚合。条件仍然
是结构化的——没有字符串 DSL：`subject` 指定表格字段，`aggregate` 选择聚合
方式，`column` 指定要聚合的数字列。

| `AggregateKind` | Wire value | `column` | 聚合结果 |
| --- | --- | --- | --- |
| `AggregateSum` | `sum` | 必填 | 指定数字列的求和 |
| `AggregateCount` | `count` | 禁止填写 | 行数 |
| `AggregateAvg` | `avg` | 必填 | 指定数字列的平均值 |

`AggregateKind.FoldsColumn()` 报告某个 kind 聚合的是列（`sum` / `avg`）
还是行（`count`）。聚合本身通过 `approval.Aggregator` 接口可插拔：

```go
type Aggregator interface {
    // Kind 返回该实现所聚合的 aggregate kind。
    Kind() AggregateKind
    // Fold 把提取出的列值（或行数）归约成比较操作数。
    // matchable=false 表示该聚合对输入没有定义值——例如对零行求 avg——
    // 此时条件必须不匹配，与 SQL NULL 比较语义一致。
    Fold(values []float64, rowCount int) (result float64, matchable bool)
}
```

用 `vef.ProvideApprovalAggregator` 在内置 `sum` / `count` / `avg` 之外注册
自定义聚合器；条件 evaluator 按其 `AggregateKind` 自动拾取，无需改动已有代码：

```go
vef.Run(
    vef.ApprovalModule,
    vef.ProvideApprovalAggregator(func() approval.Aggregator { return myMedian{} }),
    app.Module,
)
```

## 审批方式

当节点有多个审批人时：

| 方式 | 常量 | Wire value | 行为 |
| --- | --- | --- | --- |
| 顺序 | `ApprovalSequential` | `sequential` | 审批人按顺序逐个处理 |
| 并行 | `ApprovalParallel` | `parallel` | 审批人同时处理 |

枚举类型是 `ApprovalMethod`。

### 通过规则（并行模式）

| 规则 | 常量 | Wire value | 行为 |
| --- | --- | --- | --- |
| 全部 | `PassAll` | `all` | 所有审批人必须同意 |
| 任意 | `PassAny` | `any` | 至少一人同意即通过 |
| 比例 | `PassRatio` | `ratio` | 达到一定比例即通过 |

自定义通过规则实现使用 `PassRuleStrategy`、`PassRuleContext`，并返回
`PassRuleResult`（`PassRulePending`、`PassRulePassed`、`PassRuleRejected`）。

## 审批人类型

| 类型 | 常量 | Wire value | 说明 |
| --- | --- | --- | --- |
| 指定用户 | `AssigneeUser` | `user` | 特定用户 |
| 角色 | `AssigneeRole` | `role` | 拥有某角色的用户 |
| 部门 | `AssigneeDepartment` | `department` | 部门负责人 |
| 申请人本人 | `AssigneeSelf` | `self` | 申请人自己 |
| 直接上级 | `AssigneeSuperior` | `superior` | 直接上级 |
| 部门领导链 | `AssigneeDepartmentLeader` | `department_leader` | 多级主管链 |
| 表单字段 | `AssigneeFormField` | `form_field` | 由表单字段值决定 |

枚举类型是 `AssigneeKind`。动态加签位置使用 `AddAssigneeType`：
`AddAssigneeBefore`（`before`）、`AddAssigneeAfter`（`after`）、
`AddAssigneeParallel`（`parallel`）。

## 节点字段权限

任务节点（`TaskNodeData`，被审批和办理节点嵌入）和抄送节点（`CCNodeData`）
携带 `fieldPermissions` map：表单字段 key → `Permission`。取值如下：

| 常量 | Wire value | 对该节点参与者的含义 |
| --- | --- | --- |
| `PermissionVisible` | `visible` | 只读 |
| `PermissionEditable` | `editable` | 可以提交新值 |
| `PermissionHidden` | `hidden` | 不展示 |
| `PermissionRequired` | `required` | 可编辑且必须提供 |

缺失的 key 视为 `visible`。部署校验会对照派生的表单字段检查这个
map：每个 key 必须引用一个顶层表单字段，取值必须在枚举内，抄送节点只能使用
`visible` / `hidden` 子集，并且当节点的超时动作解析为 `auto_pass` 时会拒绝
`required` 权限（超时扫描器的 auto-pass 结单时不会执行必填检查）。

该 map 在写路径上强制生效：任务处理时，提交的 `formData` 只会合并
`fieldPermissions` 中标为 `editable` 或 `required` 的字段，且合并的子集会
再按字段定义重新校验。`required` 的必填检查只在 `approve` 和 `handle`
决策时执行——`reject`、`transfer`、`rollback` 保持豁免。没有配置
`fieldPermissions` 的节点不授予任何写权限：提交的表单数据会被丢弃。

在读侧，`my.get_instance_detail` 返回按查看者投影的 `fieldPermissions`：
框架在查看者的参与上下文（本人任务、抄送、申请人）上对格
（`hidden < visible < editable < required`）做 max-merge，只有 Pending
任务或可重新提交的申请人才能获得写强度（`editable` / `required`），所有
只读上下文都被钳制为 `visible`。`hidden` 的值会从返回的 `formData` 中剥离，
没有任何可识别上下文的查看者什么都看不到——解析是 fail-closed 的。

## 设计器默认值

`NodeData.ApplyTo` 会把省略的设计器字段解析为导出的默认值，保证运行时和未触碰过的设计器控件一致：

| 常量 | 值 |
| --- | --- |
| `DefaultExecutionType` | `ExecutionManual` |
| `DefaultApprovalMethod` | `ApprovalParallel` |
| `DefaultPassRule` | `PassAll` |
| `DefaultEmptyAssigneeAction` | `EmptyAssigneeAutoPass` |
| `DefaultSameApplicantAction` | `SameApplicantSelfApprove` |
| `DefaultConsecutiveApproverAction` | `ConsecutiveApproverNone` |
| `DefaultRollbackType` | `RollbackPrevious` |
| `DefaultRollbackDataStrategy` | `RollbackDataKeep` |
| `DefaultTimeoutAction` | `TimeoutActionNone` |
| `DefaultCCTiming` | `CCTimingAlways` |
| `DefaultHandleApprovalMethod` | `ApprovalSequential` |
| `DefaultHandlePassRule` | `PassAny` |
| `DefaultUrgeCooldownMinutes` | `30` |
| `DefaultTenantID` | `"default"` |

## 流程定义（兼容 React Flow）

`FlowDefinition` 结构体兼容 React Flow 的 JSON 格式：

```go
type FlowDefinition struct {
    Nodes []NodeDefinition `json:"nodes"`
    Edges []EdgeDefinition `json:"edges"`
}
```

每个 `NodeDefinition` 包含 `Kind` 和强类型的 `Data` 字段，框架会按 `Kind` 把 `Data` 解析成对应结构（`StartNodeData`、`ApprovalNodeData`、`HandleNodeData`、`ConditionNodeData`、`CCNodeData`、`EndNodeData`）。

### 流程 JSON Wire Shape

`deploy` 会把流程定义当作完整快照处理。`NodeDefinition.ParseData` 根据
`kind` 选择对应的强类型 `data` 结构；未知 kind 返回 `ErrUnknownNodeKind`，
节点 `data` 的 JSON 解析失败会用 `ErrNodeDataUnmarshal` 包装。

| 类型 | JSON 字段 |
| --- | --- |
| `FlowDefinition` | `nodes`、`edges` |
| `NodeDefinition` | `id`、`kind`、`position`、`data`；`position` 包含 `x` 和 `y` |
| `EdgeDefinition` | `id`、`source`、`target`、`sourceHandle`、`data` |

只有条件节点的出边需要 `sourceHandle`，且它必须匹配某个分支 `id`。非条件节点的出边必须省略 `sourceHandle`。`EdgeDefinition.data` 是设计器元数据，保存在版本 `flowSchema` 中；运行时流转使用 `source`、`target` 和 `sourceHandle`。

节点 `data` 字段如下：

| 节点 data 类型 | JSON 字段 |
| --- | --- |
| `BaseNodeData` | `name`、`description`；每种节点 data 都嵌入它 |
| `StartNodeData` | 只有 base 字段 |
| `EndNodeData` | 只有 base 字段 |
| `TaskNodeData` | `assignees`、`executionType`、`emptyAssigneeAction`、`fallbackUserIds`、`adminUserIds`、`isTransferAllowed`、`isOpinionRequired`、`timeoutHours`、`timeoutAction`、`timeoutNotifyBeforeHours`、`urgeCooldownMinutes`、`ccs`、`fieldPermissions` |
| `ApprovalNodeData` | base 字段 + `TaskNodeData` 字段 + `approvalMethod`、`passRule`、`passRatio`、`sameApplicantAction`、`consecutiveApproverAction`、`rollbackType`、`rollbackDataStrategy`、`rollbackTargetKeys`、`isRollbackAllowed`、`isAddAssigneeAllowed`、`addAssigneeTypes`、`isRemoveAssigneeAllowed`、`isManualCcAllowed` |
| `HandleNodeData` | base 字段 + `TaskNodeData` 字段；未设置时部署会默认 `approvalMethod = sequential`、`passRule = any` |
| `CCNodeData` | base 字段 + `ccs`、`isReadConfirmRequired`、`fieldPermissions` |
| `ConditionNodeData` | base 字段 + `branches` |

`assignees` 条目使用 `kind`、`ids`、`formField` 和 `sortOrder`。`ccs`
条目使用 `kind`、`ids`、`formField` 和 `timing`。部署时这些嵌入数组会额外物化为
`FlowNodeAssignee` 和 `FlowNodeCC` 记录，不只是写入 `FlowNode` 行。

条件分支使用 `id`、`label`、`conditionGroups`、`isDefault` 和 `priority`。
每个 `conditionGroups` 条目包含 `conditions`；每个 condition 使用 `kind`、
`subject`、`operator`、`value` 和 `expression`。字段条件也可以改为对明细表格的行做聚合：
`aggregate`（`sum` / `count` / `avg`）对 `subject` 指定的表格字段求值，
`column` 指定要聚合的数字列（`sum` / `avg` 必填，`count` 禁止填写）——见
[明细表格聚合](#明细表格聚合)。

`timeoutHours` 和 `timeoutNotifyBeforeHours` 的单位是小时。
`urgeCooldownMinutes` 的单位是分钟；小于等于 0 时使用 30 分钟运行时默认值。`rollbackTargetKeys` 只在
`rollbackType = specified` 时校验，里面放的是节点 key，不是数据库节点 ID。
`fieldPermissions` 的语义见[节点字段权限](#节点字段权限)。

### 表单 Schema 与派生字段

表单定义在部署时一分为二：

- `FlowVersion.FormSchema`（`formSchema`）是**宿主自有的表单设计器文档**，
  在 `deploy` 时以 `params.formSchema` 提交，以语义等价的 JSON 存储和返回——
  jsonb 列会规范化格式与键序，数字精度端到端保留（全程
  `json.RawMessage`）。框架从不解释它。
- `FlowVersion.FormFields`（`formFields`）是部署时通过注入的
  `approval.FormSchemaParser` 从该文档一次性派生出的扁平
  `[]FormFieldDefinition` 列表，是框架自身消费的**唯一**表单形状——用于
  表单数据校验、storage-table DDL、聚合校验和字段权限解析。parser 升级
  永远不影响已部署的版本。

```go
type FormSchemaParser interface {
    ParseFormFields(ctx context.Context, schema json.RawMessage) ([]FormFieldDefinition, error)
}
```

内置 parser 理解 vef-framework-react form-editor 文档；使用自有设计器的宿主
用 `vef.ProvideApprovalFormSchemaParser(constructor)` 整体替换。nil 或空
schema 产出零字段（没有表单的流程）；parser 报错会中止部署。`ctx` 携带部署
请求的 deadline——执行 I/O 的宿主 parser 必须遵守它。

派生字段会在部署时校验（key 唯一、kind 已知、pattern 可编译、边界一致、
明细表单层）。每个 `FormFieldDefinition` 条目使用
`key`、`kind`、`label`、`placeholder`、`defaultValue`、`isRequired`、
`options`、`validation`、`props`、`sortOrder`、`columnType`、`scale` 和
`columns`。每个 option 使用 `label` 和 `value`。`columns` 定义 table 字段
（`kind` 为 `table`）的行结构：每一列本身也是一个 `FormFieldDefinition`，且不能再声明
自己的 `columns`——明细表格只能是单层。table 字段自身的 `validation.minLength` /
`maxLength` 用于约束行数，`isRequired` 表示至少要有一行。

`validation` 支持 `minLength`、`maxLength`、`min`、`max`、`pattern` 和
`message`。提交的 `formData` 在 JSON 编码后最大 64 KiB，即使流程没有表单
schema 也会执行这个大小限制。有 schema 时，额外的表单 key 会被拒绝；必填字段会拒绝
缺失、`null`、空白字符串和空数组。`input`、`textarea`、`date` 字段必须是字符串，
可使用 `minLength`、`maxLength` 和 `pattern`。`number` 字段接受 JSON 数字，
可使用 `min` 和 `max`。`select` 字段在配置了 `options` 时会校验标量或数组值是否
存在于选项中。`upload` 字段接受非空白字符串、非空 `[]string`，或非空且每项都是非空白
字符串的数组。`validation.message` 只作为 `pattern` 不匹配时的自定义错误信息；
其他校验失败使用模块 i18n 消息。

## 流程模型与设计器枚举

公开包暴露的流程设计和持久化模型包括 `FlowCategory`、`Flow`、`FlowVersion`、
`FlowNode`、`FlowEdge`、`FlowInitiator`、`FlowNodeAssignee`、`FlowNodeCC`、
`FormFieldDefinition`、`FormSnapshot`、`ActionLog`、
`OperatorInfo` 和 `UrgeRecord`（没有结构化的 `FormDefinition`
包装——宿主文档是 opaque 的，框架侧只有 `FormFieldDefinition` 这一种
形状）。流程版本状态使用 `VersionStatus`：
`VersionDraft`（`draft`）、`VersionPublished`（`published`）、
`VersionArchived`（`archived`）。

`Flow.Labels` 是宿主自有的筛选元数据——保存时按共享 label 规则
校验、在列表查询中支持相等过滤；wire 形状见 [RPC 资源](./resources.md)。

其他流程设计器枚举：

| 枚举 | Wire values |
| --- | --- |
| `InitiatorKind` | `user`、`role`、`department` |
| `ExecutionType` | `manual`、`auto_pass`、`auto_reject` |
| `ConditionKind` | `field`、`expression` |
| `CCKind` | `user`、`role`、`department`、`form_field` |
| `CCTiming` | `always`、`on_approve`、`on_reject` |
| `FieldKind` | `input`、`textarea`、`select`、`number`、`date`、`upload`、`table` |
| `ColumnDataType` | `string`、`text`、`integer`、`decimal`、`boolean`、`date`、`datetime`、`json` |
| `Permission` | `visible`、`editable`、`hidden`、`required` |

---

下一步：[实例运行时](./runtime.md) 了解设计好的流程启动后会发生什么。
