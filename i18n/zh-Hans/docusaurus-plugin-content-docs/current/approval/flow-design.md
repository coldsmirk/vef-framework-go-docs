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

### 通过规则

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
| 部门 | `AssigneeDepartment` | `department` | 所配置部门的负责人 |
| 申请人本人 | `AssigneeSelf` | `self` | 申请人自己 |
| 直接上级 | `AssigneeSuperior` | `superior` | 直接上级 |
| 部门领导链 | `AssigneeDepartmentLeader` | `department_leader` | 申请人所在部门的负责人（单层查询） |
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
| `HandleNodeData` | base 字段 + `TaskNodeData` 字段；部署时无条件设置 `approvalMethod = sequential`、`passRule = any`（handle 节点没有这两个 wire 字段） |
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
明细表单层、选项来源自洽）。每个 `FormFieldDefinition` 条目使用
`key`、`kind`、`label`、`placeholder`、`defaultValue`、`isRequired`、
`options`、`optionSource`、`validation`、`props`、`sortOrder`、`columnType`、
`scale` 和 `columns`。每个 option 使用 `label` 和 `value`。`columns` 定义 table 字段
（`kind` 为 `table`）的行结构：每一列本身也是一个 `FormFieldDefinition`，且不能再声明
自己的 `columns`——明细表格只能是单层。table 字段自身的 `validation.minLength` /
`maxLength` 用于约束行数，`isRequired` 表示至少要有一行。

`validation` 支持 `minLength`、`maxLength`、`min`、`max`、`pattern` 和
`message`。提交的 `formData` 受 `vef.approval.form_data_max_bytes` 限制
（默认值 64 KiB，来自 `config.DefaultFormDataMaxBytes`），即使流程没有表单
schema 也会执行这个大小限制。有 schema 时，额外的表单 key 会被拒绝；必填字段会拒绝
缺失、`null`、空白字符串和空数组。`input`、`textarea`、`date` 字段必须是字符串，
可使用 `minLength`、`maxLength` 和 `pattern`。`number` 字段接受 JSON 数字，
可使用 `min` 和 `max`；`columnType` 为 `integer` 的 `number` 字段还会拒绝小数值。
`select` 字段在配置了 `options` 时会校验标量或数组值是否
存在于选项中。`upload` 字段接受非空白字符串、非空 `[]string`，或非空且每项都是非空白
字符串的数组。`validation.message` 只作为 `pattern` 不匹配时的自定义错误信息；
其他校验失败使用模块 i18n 消息。

#### 选项来源：枚举出来的，还是远程的

选择类字段的选项只会以两种形状之一出现，绝不会同时出现：

- `options`（`[]FieldOption`）——部署时选项可以被枚举，投影就直接把它写出来。
  静态来源无论是内联配置在字段上，还是通过表单全局的 `ref` 引用到的，都会被
  枚举出来。
- `optionSource`（`*FieldOptionSource`）——选项**无法**被枚举，投影于是输出
  消费方自行拉取所需的描述符。只有远程来源会走到这个形状。

两者都没有的字段就是自由输入：没有任何东西约束它的取值。

投影结果是**解引用之后**的。设计器里的 `ref` 会先针对表单全局来源解析，再写入
字段，所以消费方永远不需要自己去追 `dataSourceId`。指向不存在来源的 `ref` 既不
产出 `options` 也不产出 `optionSource`——和完全没有配置来源的字段一样。

`FieldOptionSource` 的结构：

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| `kind` | `OptionSourceKind` | 来源分类；目前只会输出 `OptionSourceRemote`（`remote`） |
| `request` | `*RemoteOptionRequest` | 返回选项记录的那个操作 |
| `mapping` | `*RemoteOptionMapping` | 如何从每条记录里读出 label 和 value；`nil` 表示用默认值 |

设计器自身的来源联合类型还有 `static` 和 `ref`，但两者都活不过投影，所以
`OptionSourceKind` 的词汇表刻意比设计器的更窄。

`RemoteOptionRequest` 用框架自己的 resource/action 寻址方式定位那个操作：
`resource`、`action`、`version`（为空表示默认版本）和 `params`。
`RemoteOptionMapping` 指出记录里的键名——`labelKey`（默认 `label`）、
`valueKey`（默认 `value`）、`disabledKey` 和 `descriptionKey`；留空的条目回落到
各自的默认值。

**框架从不解析选项来源**：表单数据校验只读 `options`，而没有 options 的
`select` 字段接受任意提交值——所以 `optionSource` 在服务端不约束任何东西。它的
存在是为了那些必须把**已存储的值**渲染成 label 的消费方：业务列表页展示一个存下来
的 select 值时需要的是 label 而不是原始编码，而静态那一半早就由 `options` 覆盖了。

部署校验同样覆盖选项来源，顶层字段和明细表格列一视同仁：`kind` 不在词汇表内的
`optionSource`，或者 `request` 缺失、`resource` / `action` 为空白的远程来源，都会
让部署失败。这两种情况在别处都不会失败——框架从不解析来源——于是版本会干干净净地
部署上去，然后在每一个回放它的消费方那里渲染成原始值。

明细表格的列自带这一份：`columns` 里的条目本身就是完整的 `FormFieldDefinition`，
所以 `table` 字段内部的 `select` 列拥有和顶层字段一样的 `options` /
`optionSource` 组合。

#### 远程请求参数

`RemoteOptionRequest.Params` 是 `map[string]DynamicParam`，且**不求值**地携带。
每个参数只会是两种 kind 之一：

| `DynamicParamKind` | Wire value | 载荷 | 含义 |
| --- | --- | --- | --- |
| `DynamicParamLiteral` | `literal` | `value` | 设计器填写的固定值，每次求值都一样 |
| `DynamicParamExpression` | `expression` | `source` | 表单运行时在发起请求前针对当前表单值求值的表达式——级联选择就是靠它 |

后端手里没有表单值可以拿来对表达式求值，因此它只存源文本，求值这一步归回放请求的
消费方所有。`value` 刻意不带 `omitempty`：字面量的 `false`、`0`、`""` 都是设计器
选定的值，丢掉它会静默改变请求。

**翻译整列之前，先问 `HasBoundParams()`。**
`(*RemoteOptionRequest).HasBoundParams()` 报告是否存在表达式参数。消费方在假定
"一次查询就能覆盖整列"之前，必须先问这个问题：

- **`false`**——请求对每一行的解析结果都相同，所以**一次**调用就能为整列建好
  value → label 映射。
- **`true`**——选项集是**按行**的（参数依赖该行自己的表单值），列表页要么逐行求值
  并发起请求，要么让这一列不翻译。

该方法对 nil 安全：nil receiver 返回 `false`。

#### 在宿主里把存储值翻译成 label

注入 `approval.FormSchemaParser`——只要启用了 `vef.ApprovalModule`，它就能在
root scope 解析出来——然后解析版本的表单 schema：

```go
package options

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/coldsmirk/vef-framework-go/approval"
)

// OptionLabelPlan says how a list view should turn each selection field's
// stored value into a display label.
type OptionLabelPlan struct {
    // Static maps a field key to its value-to-label table, built from the
    // options the definition already carries.
    Static map[string]map[string]string
    // Remote maps a field key to the request a consumer replays ONCE to build
    // that table itself. Only fields whose request has no bound parameters.
    Remote map[string]*approval.RemoteOptionRequest
    // PerRow lists the field keys whose option set depends on the row's own
    // form values, so one lookup cannot translate the whole column.
    PerRow []string
}

type OptionLabelService struct {
    parser approval.FormSchemaParser
}

func NewOptionLabelService(parser approval.FormSchemaParser) *OptionLabelService {
    return &OptionLabelService{parser: parser}
}

func (s *OptionLabelService) Plan(ctx context.Context, schema json.RawMessage) (*OptionLabelPlan, error) {
    fields, err := s.parser.ParseFormFields(ctx, schema)
    if err != nil {
        return nil, err
    }

    plan := &OptionLabelPlan{
        Static: make(map[string]map[string]string),
        Remote: make(map[string]*approval.RemoteOptionRequest),
    }

    for _, field := range fields {
        switch {
        case field.Options != nil:
            // Enumerated at deploy: translate locally, no lookup at all.
            labels := make(map[string]string, len(field.Options))
            for _, option := range field.Options {
                labels[fmt.Sprint(option.Value)] = option.Label
            }

            plan.Static[field.Key] = labels

        case field.OptionSource != nil:
            request := field.OptionSource.Request
            if request.HasBoundParams() {
                plan.PerRow = append(plan.PerRow, field.Key)

                continue
            }

            plan.Remote[field.Key] = request

        default:
            // Free-form field: the stored value is already what to display.
        }
    }

    return plan, nil
}

// labelKey and valueKey apply RemoteOptionMapping's defaults.
func labelKey(mapping *approval.RemoteOptionMapping) string {
    if mapping == nil || mapping.LabelKey == "" {
        return "label"
    }

    return mapping.LabelKey
}

func valueKey(mapping *approval.RemoteOptionMapping) string {
    if mapping == nil || mapping.ValueKey == "" {
        return "value"
    }

    return mapping.ValueKey
}
```

值按字符串形态（`fmt.Sprint`）比较——框架自己的 `select` 校验就是这么做的——所以
`1` 和 `"1"` 命中同一个选项。

解析并不是唯一入口。同一份扁平列表已经持久化在 `apv_flow_version.form_fields`
上，并由流程版本详情响应以 `formFields` 返回，所以手里已经有那一行的宿主直接读
`FlowVersion.FormFields` 即可。只有当手里只有 schema 文档时才需要解析——实例详情
响应原样返回宿主文档 `formSchema`，从不返回派生字段。

对于成功部署的版本，`OptionSource.Request` 是可以信赖的：上面那条"缺操作即拒绝"
的校验会挡掉没有 request 的远程来源，所以在活下来的字段上它永远不是 `nil`。

#### 上传字段

内置 parser 会把设计器的 `upload` 组件投影为 `upload` 字段 kind。

上传字段的 `columnType` 由它的**文件数量**推断，而不是由 `maxLength` 推断：
`maxCount > 1` 得到 `json`（值是一组存储 key 的数组——和 `checkbox-group` 的形状
相同），其余情况得到 `text`（单个存储 key）。`maxLength` 约束的是*字符串*长度，
而上传字段的界限是它接受多少个文件，在这里读它就会按错误的量纲给列定尺寸。字段上
显式声明的 `columnType` 依然优先，这一点和所有 kind 一致。

## 流程校验

### 发起人规则

`isAllInitiationAllowed` 与 `initiators` **互斥**：当 `isAllInitiationAllowed`
为 `true` 时提交发起人规则会返回 `ErrInitiatorsNotAllowed`（code `40022`）。
当限定时，至少需要一个发起人规则且至少包含一个 ID——空规则集或 `ids` 为空
的规则会返回 `ErrInitiatorsRequired`（code `40023`），因为一个不选中任何人的规则
和完全没有规则一样，流程都无法被发起。

### 业务绑定 Schema

流程使用 `BindingMode` 决定表单/实例数据如何与业务数据关联：

- `BindingStandalone`（wire 值 `standalone`）—— 表单数据存在审批自己的表中（默认）。
- `BindingBusiness`（wire 值 `business`）—— 流程通过 `BusinessBindingConfig`
  挂接到已有的业务行。

业务绑定流程（`BindingMode = business`）在保存时会对照真实数据库 schema 校验。
当配置的 `tableName` 不存在时，返回 `ErrBindingTableMissing(tableName)`（code
`40018`，`ErrCodeBindingSchemaInvalid`）。当配置的列在表中不存在时，返回
`ErrBindingColumnMissing(columnName)`（同 code）。两者都会命名缺失的对象，让
操作员不必亲自查看 schema 就能修正配置。

## 流程模型与设计器枚举

公开包暴露的流程设计和持久化模型包括 `FlowCategory`、`Flow`、`FlowVersion`、
`FlowNode`、`FlowEdge`、`FlowInitiator`、`FlowNodeAssignee`、`FlowNodeCC`、
`FormFieldDefinition`、`FieldOption`、`FieldOptionSource`、
`RemoteOptionRequest`、`RemoteOptionMapping`、`DynamicParam`、
`FormSnapshot`、`ActionLog`、
`UserInfo` 和 `UrgeRecord`（没有结构化的 `FormDefinition`
包装——宿主文档是 opaque 的，框架侧的表单形状只有 `FormFieldDefinition`
及它内嵌的那些类型）。流程版本状态使用 `VersionStatus`：
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
| `OptionSourceKind` | `remote` |
| `DynamicParamKind` | `literal`、`expression` |

---

下一步：[实例运行时](./runtime.md) 了解设计好的流程启动后会发生什么。
