---
sidebar_position: 2
---

# RPC 资源

启用审批模块后，框架会注册下面六个 RPC 资源。它们全部挂载在 `/api` 下，
使用 [API](../building-apis/api.md) 中的标准 envelope（`resource`、`action`、
`version`、`params`、`meta`）。这些操作都不是公开接口：调用者必须已认证；
表中列出 `RequiredPermission` 时还会校验对应权限点。生成的
[运行时 API 索引](../reference/runtime-api-index.md) 包含每个请求/响应 DTO 的
完整 JSON 字段清单。

本页使用的约定：

- 命令式操作（`approval/flow`、`approval/instance`、`approval/my`、
  `approval/admin`）的参数结构体嵌入 `api.P`，字段从请求的 `params`
  对象解码。
- CRUD 读操作（`approval/category`、`approval/delegation`）的搜索结构体嵌入
  `crud.Sortable`（一个 `meta` 结构体），因此过滤字段从请求的 `meta` 对象
  解码，与 `meta.page` / `meta.size`（`page.Pageable`）和 `meta.sort` 并列。
- 分页响应使用 `page.Page[T]`：`page`、`size`、`total`、`items`。
- 持久化模型在响应中携带标准审计列（`id`、`createdAt`、`createdBy`、
  `updatedAt`、`updatedBy`），下面的字段表不再重复列出。
- 枚举词汇（`InstanceStatus`、`TaskStatus`、节点语义）的定义见
  [实例运行时](./runtime.md) 与 [流程设计](./flow-design.md)。

## `approval/category`

流程分类管理（`apv_flow_category`）。读操作按租户隔离：super-admin 可见
所有租户，其他调用者只能看到自己租户的数据，无租户时直接拒绝（fail
closed）。

| Action | 权限 | 入参 | 出参 |
| --- | --- | --- | --- |
| `find_tree` | `approval.category.query` | `CategorySearch`（meta） | 嵌套的 `FlowCategory[]`（填充 children） |
| `create` | `approval.category.create` | `CategoryParams` | 创建后的 `FlowCategory` |
| `update` | `approval.category.update` | `CategoryParams` | 更新后的 `FlowCategory` |
| `delete` | `approval.category.delete` | 主键参数（`params.id`） | 成功 |

没有 `find_tree_options` 操作；选项列表请由 `find_tree` 构建。

`CategorySearch`（查询过滤，从 `meta` 解码）：

| 字段 | 类型 | 匹配 | 说明 |
| --- | --- | --- | --- |
| `name` | `string` | contains | 按分类名称片段过滤 |
| `isActive` | `bool` | equals | 按启用状态过滤；省略则两者都匹配 |
| `sort` | `OrderSpec[]` | — | 排序声明（`crud.Sortable`） |

`CategoryParams`（create/update，从 `params` 解码）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `id` | `string` | 仅 update | 要更新记录的主键 |
| `tenantId` | `string` | 是 | 所属租户。create 时非 super-admin 会以调用者自己的租户覆盖写入（提交值被忽略）；update/delete 时调用者必须对记录的租户有权限 |
| `code` | `string` | 是 | 分类业务编码 |
| `name` | `string` | 是 | 显示名称 |
| `icon` | `string` | 否 | 显示图标标识 |
| `parentId` | `string` | 否 | 父分类 id；`null` 表示根分类 |
| `sortOrder` | `int` | 否 | 同级排序权重 |
| `isActive` | `bool` | 否 | 停用分类仍可查询，宿主通常在选择器中隐藏 |
| `remark` | `string` | 否 | 备注 |

`FlowCategory`（响应模型）：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `tenantId` | `string` | 所属租户 |
| `code` | `string` | 分类业务编码 |
| `name` | `string` | 显示名称 |
| `icon` | `string` \| `null` | 显示图标标识 |
| `parentId` | `string` \| `null` | 父分类 id |
| `sortOrder` | `int` | 排序权重 |
| `isActive` | `bool` | 启用标记 |
| `remark` | `string` \| `null` | 备注 |
| `children` | `FlowCategory[]` | 子分类；`find_tree` 填充，其余场景缺省 |

## `approval/delegation`

审批委托管理（`apv_delegation`）。强制所有权：非 super-admin 只能查看、
创建、更新、删除自己作为委托人的记录 —— create 时 `delegatorId` 以调用者
覆盖写入，update 时钉住原委托人，防止把记录转给他人。

| Action | 权限 | 入参 | 出参 |
| --- | --- | --- | --- |
| `find_page` | `approval.delegation.query` | `DelegationSearch` + 分页 meta | `page.Page[Delegation]` |
| `create` | `approval.delegation.create` | `DelegationParams` | 创建后的 `Delegation` |
| `update` | `approval.delegation.update` | `DelegationParams` | 更新后的 `Delegation` |
| `delete` | `approval.delegation.delete` | 主键参数（`params.id`） | 成功 |

`DelegationSearch`（查询过滤，从 `meta` 解码）：

| 字段 | 类型 | 匹配 | 说明 |
| --- | --- | --- | --- |
| `delegatorId` | `string` | equals | 按委托人过滤（仅 super-admin 有意义 —— 其他人始终被限定为本人） |
| `delegateeId` | `string` | equals | 按受托人过滤 |
| `isActive` | `bool` | equals | 按启用状态过滤 |
| `sort` | `OrderSpec[]` | — | 排序声明 |

`DelegationParams`（create/update，从 `params` 解码）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `id` | `string` | 仅 update | 要更新记录的主键 |
| `delegatorId` | `string` | 是 | 委托人。非 super-admin create 时以 principal 覆盖写入，update 时钉住原值 |
| `delegateeId` | `string` | 是 | 接收委托任务的用户 |
| `flowCategoryId` | `string` | 否 | 将委托限定到某个流程分类；`null` 覆盖全部分类 |
| `flowId` | `string` | 否 | 将委托限定到某个流程；`null` 覆盖全部流程 |
| `startTime` | `DateTime` | 是 | 委托生效开始时间 |
| `endTime` | `DateTime` | 是 | 委托生效结束时间 |
| `isActive` | `bool` | 否 | 停用的委托在办理人解析中被忽略 |
| `reason` | `string` | 否 | 委托原因，展示在被委托的任务上 |

`Delegation`（响应模型）：业务字段与入参一致 —— `delegatorId`、
`delegateeId`、`flowCategoryId`、`flowId`、`startTime`、`endTime`、
`isActive`、`reason` —— 外加审计列。经委托产生的任务会把委托人作为独立的
人员快照携带（见下文 `NodeParticipant.delegator`）。

## `approval/flow`

流程定义管理：可变的流程行、不可变的已部署版本、面向设计器的图查询。

| Action | 权限 | 入参 | 出参 | 审计 |
| --- | --- | --- | --- | --- |
| `create` | `approval.flow.create` | `CreateFlowParams` | 创建后的 `Flow` | 是 |
| `deploy` | `approval.flow.deploy` | `DeployFlowParams` | 创建的 `FlowVersion`（draft） | 是 |
| `publish_version` | `approval.flow.publish` | `PublishVersionParams` | 成功 | 是 |
| `update` | `approval.flow.update` | `UpdateFlowParams` | 更新后的 `Flow` | 是 |
| `toggle_active` | `approval.flow.update` | `ToggleActiveParams` | 成功 | 是 |
| `get_graph` | `approval.flow.query` | `GetGraphParams` | `FlowGraph` | — |
| `find_flows` | `approval.flow.query` | `FindFlowsParams` | `page.Page[Flow]` | — |
| `find_versions` | `approval.flow.query` | `FindVersionsParams` | `FlowVersionSummary[]` | — |
| `find_initiators` | `approval.flow.query` | `FindInitiatorsParams` | `FlowInitiator[]` | — |

`CreateFlowParams`（`create`）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `tenantId` | `string` | 是 | 所属租户；空值回退为 `"default"`，调用者必须对最终租户有权限 |
| `code` | `string` | 是 | 全局唯一的流程业务编码；创建后不可变（`start` 以它定位流程） |
| `name` | `string` | 是 | 显示名称 |
| `categoryId` | `string` | 是 | 所属 `FlowCategory` id |
| `icon` | `string` | 否 | 显示图标标识 |
| `description` | `string` | 否 | 描述 |
| `labels` | `object`（string→string） | 否 | 宿主自有的选择元数据；按共享 label 规则校验 —— 见下方说明 |
| `bindingMode` | `string` | 是 | `standalone`（表单数据存审批自己的表）或 `business`（挂接既有业务行） |
| `businessBinding` | `BusinessBindingConfig` | business 模式 | 回写目标描述（见下）；standalone 流程提交会被拒绝（`ErrBindingUnexpected`） |
| `adminUserIds` | `string[]` | 否 | 流程管理员（用于 `transfer_admin` 空办理人处理与管理可见性） |
| `isAllInitiationAllowed` | `bool` | 否 | `true` 允许所有用户发起；`false` 时只有 `initiators` 可发起 |
| `instanceTitleTemplate` | `string` | 否 | 实例标题的 Go `text/template`，如 `{{.applicantName}}的请假申请`；可用绑定：`flowName`、`flowCode`、`instanceNo`、`formData`、`applicantId`、`applicantName`（以及嵌套的 `flow.name` / `flow.code`、`applicant.id` / `applicant.name`）。留空回退为 `流程名-实例编号`；解析失败报 `ErrInvalidTitleTemplate` |
| `initiators` | `CreateInitiatorParams[]` | 否 | `isAllInitiationAllowed` 为 `false` 时允许发起的人（见下） |

`CreateInitiatorParams` 条目：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `kind` | `string` | 是 | `user`、`role` 或 `department` |
| `ids` | `string[]` | 是 | 所选用户 / 角色 / 部门的 id |

`BusinessBindingConfig`：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `tableName` | `string` | 是 | 接收审批状态的业务表；按 SQL 安全标识符校验（`ErrInvalidBusinessIdentifier`）并检查存在性（`ErrBindingSchemaInvalid`） |
| `keyColumns` | `string[]` | 是 | 定位绑定行的列；必须与某个非空主键或唯一键完全一致（`ErrBindingKeyNotUnique`） |
| `statusColumn` | `string` | 是 | 接收（映射后）实例状态的列 |
| `instanceIdColumn` | `string` | 是 | 接收当前实例 id 的列；作为 compare-and-set 栅栏，防止过期实例覆盖新一轮审批的状态 |
| `startedAtColumn` | `string` | 否 | 接收实例开始时间的列 |
| `finishedAtColumn` | `string` | 否 | 接收实例结束时间的列 |
| `statusMapping` | `object`（InstanceStatus→string） | 否 | 把实例状态翻译成宿主业务词汇；缺失条目回退为状态字符串本身（未知键或空值报 `ErrBindingStatusMappingInvalid`） |

两个绑定字段指向同一列会报 `ErrBindingColumnsConflict`。已部署版本会对绑定
做快照，编辑流程的绑定不影响运行在旧版本上的实例。回写生命周期见
[业务集成](./integration.md)。

`params.labels` 是流程上宿主自有的选择元数据 —— 在 `find_flows`
与 `my.find_available_flows` 中可按相等过滤（提交的每一对都必须匹配），
在实例详情视图中透出，引擎从不解释其取值。校验为共享的
`orm.ValidateLabels` 规则：键为字母数字加内部 `-`/`_`（不允许点号），
≤ 63 字符；值 ≤ 256 字符，允许空值。`update` 时 labels 整体替换 ——
省略即清空流程的 labels。

`DeployFlowParams`（`deploy` —— 创建新的 draft 版本）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `flowId` | `string` | 是 | 目标流程 |
| `description` | `string` | 否 | 展示在版本列表的版本描述 |
| `storageMode` | `string` | 否 | `json`（默认；表单数据留在 `apv_instance.form_data`）或 `table`（发布时生成专用物理投影表） |
| `flowDefinition` | `FlowDefinition` | 是 | 设计器图文档 —— 节点、边与每个节点的 `data`；部署时校验（`ErrInvalidFlowDesign`）。线上传输格式见 [流程设计](./flow-design.md) |
| `formSchema` | JSON 文档 | 否 | 宿主自有的表单设计器文档，原样透传并存储；引擎消费的扁平字段清单在部署时从中推导（见 [表单 Schema 与派生字段](./flow-design.md#表单-schema-与派生字段)） |

`PublishVersionParams`（`publish_version` —— 将 draft 版本设为线上版本并归档
前一版本）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `versionId` | `string` | 是 | 要发布的 draft 版本（否则报 `ErrVersionNotDraft`） |

`UpdateFlowParams`（`update` —— 只修改流程行；已部署版本不可变）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `flowId` | `string` | 是 | 要更新的流程 |
| `name` | `string` | 是 | 显示名称 |
| `icon` | `string` | 否 | 显示图标标识 |
| `description` | `string` | 否 | 描述 |
| `labels` | `object`（string→string） | 否 | 整体替换；省略即清空 |
| `bindingMode` | `string` | 是 | 绑定模式（见 `create`）；修改只影响后续部署 |
| `businessBinding` | `BusinessBindingConfig` | business 模式 | 回写目标（见 `create`） |
| `adminUserIds` | `string[]` | 否 | 流程管理员 |
| `isAllInitiationAllowed` | `bool` | 否 | 发起开放性 |
| `instanceTitleTemplate` | `string` | 是 | 实例标题模板 |
| `initiators` | `CreateInitiatorParams[]` | 否 | 发起人配置；整体替换 |

`ToggleActiveParams`（`toggle_active`）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `flowId` | `string` | 是 | 目标流程 |
| `isActive` | `bool` | 否 | 目标状态；停用的流程拒绝发起（`ErrFlowNotActive`），运行中的实例不受影响 |

`GetGraphParams`（`get_graph`）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `flowId` | `string` | 是 | 要加载图的流程 |
| `tenantId` | `string` | 否 | 可选预过滤；真正的跨租户闸门是调用者的租户权限 |
| `versionId` | `string` | 否 | 显式加载某个版本 —— 设计器从最新部署（无论是否已发布）继续编辑的场景；省略则解析最新已发布版本 |

`get_graph` 响应 `FlowGraph`：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `flow` | `Flow` | 可变的流程行（字段见下） |
| `version` | `FlowVersion` | 解析出的版本，含 `flowSchema`（部署的 `FlowDefinition`）、`formSchema`（宿主文档，原样）与 `formFields`（推导出的扁平字段清单） |
| `nodes` | `FlowNode[]` | 该版本持久化的节点行 —— 每个节点一行，携带全部已解析的节点配置（kind、执行类型、审批方式、通过规则、回退 / 加签 / 抄送开关、超时配置、分支）。各字段语义见 [流程设计](./flow-design.md) |
| `edges` | `FlowEdge[]` | 持久化的边行：`key`、`sourceNodeId` / `sourceNodeKey`、`targetNodeId` / `targetNodeKey`、`sourceHandle`（条件分支锚点） |

`FindFlowsParams`（`find_flows`）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `tenantId` | `string` | 否 | 租户过滤；非 super-admin 无论如何都被限定在自己租户 |
| `categoryId` | `string` | 否 | 按分类过滤 |
| `keyword` | `string` | 否 | 对流程名称做 contains 匹配 |
| `isActive` | `bool` | 否 | 按启用状态过滤 |
| `labels` | `object`（string→string） | 否 | label 相等过滤 —— 提交的每一对都必须匹配 |
| `page` | `int` | 否 | 页码（从 1 开始） |
| `pageSize` | `int` | 否 | 每页大小 |

`find_flows` 响应 `page.Page[Flow]`。`Flow`（响应模型）：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `tenantId` | `string` | 所属租户 |
| `categoryId` | `string` | 所属分类 |
| `code` | `string` | 唯一流程业务编码（不可变） |
| `name` | `string` | 显示名称 |
| `icon` | `string` \| `null` | 显示图标标识 |
| `description` | `string` \| `null` | 描述 |
| `labels` | `object` \| 缺省 | 宿主自有选择元数据 |
| `bindingMode` | `string` | `standalone` 或 `business` |
| `businessBinding` | `BusinessBindingConfig` \| 缺省 | 当前回写配置（可变副本；各版本有自己的快照） |
| `adminUserIds` | `string[]` | 流程管理员 |
| `isAllInitiationAllowed` | `bool` | 发起开放性 |
| `instanceTitleTemplate` | `string` | 实例标题模板 |
| `isActive` | `bool` | 启用标记 |
| `currentVersion` | `int` | 最新已发布版本号；首次发布前为 `0` |

`FindVersionsParams` / `FindInitiatorsParams`（`find_versions`、
`find_initiators`）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `flowId` | `string` | 是 | 要查看的流程 |
| `tenantId` | `string` | 否 | 可选预过滤（跨租户闸门是调用者权限） |

`find_versions` 返回 `FlowVersionSummary` 条目 —— 不含图文档
（`flowSchema` / `formSchema` / `formFields`）的版本列表，列表本就不渲染
它们。单个版本的完整定义通过 `get_graph` 携带 `params.versionId` 获取。

| `FlowVersionSummary` 字段 | 类型 | 说明 |
| --- | --- | --- |
| `id` | `string` | 版本 id |
| `flowId` | `string` | 所属流程 |
| `version` | `int` | 单调递增的版本号 |
| `status` | `string` | `draft`、`published` 或 `archived` |
| `description` | `string` \| `null` | 版本描述 |
| `storageMode` | `string` | `json` 或 `table` |
| `publishedAt` | `DateTime` \| `null` | 发布时间 |
| `publishedBy` | `string` \| `null` | 发布人用户 id |
| `createdAt` | `DateTime` | 部署时间 |
| `createdBy` | `string` | 部署人用户 id |

`find_initiators` 返回 `FlowInitiator[]`：每条携带 `flowId`、`kind`
（`user` / `role` / `department`）与 `ids`（配置的 id 列表）。

## `approval/instance`

实例生命周期命令。每次状态变更都记录在 action log 中；标记「审计」的操作
还会额外捕获框架级 IP / UA / request-id 审计条目。

| Action | 权限 | 入参 | 出参 | 审计 |
| --- | --- | --- | --- | --- |
| `start` | `approval.instance.start` | `StartParams` | 创建的 `Instance` | 是 |
| `process_task` | `approval.task.process` | `ProcessTaskParams` | 成功 | 是 |
| `withdraw` | `approval.instance.withdraw` | `WithdrawParams` | 成功 | 是 |
| `resubmit` | `approval.instance.resubmit` | `ResubmitParams` | 成功 | 是 |
| `add_cc` | `approval.instance.cc` | `AddCCParams` | 成功 | 是 |
| `mark_cc_read` | `approval.instance.cc` | `MarkCCReadParams` | 成功 | — |
| `add_assignee` | `approval.task.add_assignee` | `AddAssigneeParams` | 成功 | 是 |
| `remove_assignee` | `approval.task.remove_assignee` | `RemoveAssigneeParams` | 成功 | 是 |
| `urge_task` | `approval.task.urge` | `UrgeTaskParams` | 成功 | 限流：每分钟最多 `10` 次 |

`process_task` 有意把 approve / reject / transfer / rollback / handle 归在
同一个权限（`approval.task.process`）下：设计器的节点级开关
（`isTransferAllowed`、`isRollbackAllowed` 等）已经决定了节点在运行时提供
哪些动作。

`StartParams`（`start`）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `tenantId` | `string` | 是 | 发起所在租户；空值回退为 `"default"`，且调用者必须对流程的租户有权限 |
| `flowCode` | `string` | 是 | 要发起流程的业务编码；解析最新已发布版本（`ErrFlowNotFound` / `ErrFlowNotActive` / `ErrNoPublishedVersion`） |
| `businessRef` | `string`（≤ 512） | business 模式 | 绑定业务行的不透明引用；业务绑定流程必填，除非注册的 `BusinessRefProvider` 提供（`ErrBusinessRefRequired`）。默认形状：单键直接取值、复合键为 JSON 对象 |
| `formData` | `object` | 否 | 按字段 key 组织的表单值；按已发布版本的派生字段清单校验（`40401` 系列），超过 64 KiB 拒绝，并剥除申请人无权编辑的字段 |

申请人身份与条件路由的全局变量在服务端从已认证 principal 解析
（`PrincipalDepartmentResolver`、`InstanceGlobalsResolver`）—— 绝不接受请求体
提交，否则申请人可以伪造它们来操纵流程走向。

`start` 响应创建的 `Instance`：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `tenantId` | `string` | 所属租户 |
| `flowId` / `flowCode` / `flowVersionId` | `string` | 实例运行所依据的流程与不可变版本快照 |
| `title` | `string` | 由流程的 `instanceTitleTemplate` 渲染 |
| `instanceNo` | `string` | 人类可读的实例编号 |
| `applicantId` / `applicantName` | `string` | 发起时的申请人快照 |
| `applicantDepartmentId` / `applicantDepartmentName` | `string` \| `null` | 申请人部门快照 |
| `status` | `string` | `running`、`approved`、`rejected`、`withdrawn`、`returned` 或 `terminated` |
| `currentNodeId` | `string` \| `null` | 实例当前停留的节点 |
| `finishedAt` | `DateTime` \| `null` | 实例到达终态时写入 |
| `businessRef` | `string` \| `null` | 不透明业务引用（业务绑定流程） |
| `formData` | `object` | 提交的表单数据（校验后） |
| `globals` | `object` | 发起时宿主提供的全局变量快照；条件求值读取它，保证路由在重复求值间确定 |
| `businessProjectionId` | `string` \| 缺省 | 发起时认领的持久回写状态（业务绑定流程） |

`ProcessTaskParams`（`process_task`）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `taskId` | `string` | 是 | 要处理的待办任务；调用者必须是其办理人（`ErrNotAssignee`、`ErrTaskNotPending`） |
| `action` | `string` | 是 | `approve`、`reject`、`transfer`、`rollback` 或 `handle`（办理节点以 `handle` 完结；语义同 approve） |
| `opinion` | `string`（≤ 2000） | 视节点 | 处理意见；节点设置 `isOpinionRequired` 时必填（`ErrOpinionRequired`） |
| `formData` | `object` | 否 | 随动作写入的表单更新，按节点字段权限过滤 |
| `attachments` | `string[]`（≤ 20 × ≤ 512） | 否 | 存入 action log 的附件引用 |
| `transferToId` | `string` | `transfer` | 转办目标用户；必须非空且不同于操作者（`ErrInvalidTransferTarget`）；仅节点开启 `isTransferAllowed` 时允许（`ErrTransferNotAllowed`） |
| `targetNodeId` | `string` | `rollback` | 回退目标节点；必须是节点 `rollbackType` 与实例访问轨迹允许的目标之一（`ErrInvalidRollbackTarget`、`ErrRollbackNotAllowed`）。合法目标由 `my.get_instance_detail` → `myTask.rollbackTargets` 提供 |

`WithdrawParams`（`withdraw` —— 申请人撤回运行中的实例）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `instanceId` | `string` | 是 | 要撤回的实例；调用者必须是申请人（`ErrNotApplicant`），状态必须允许（`ErrWithdrawNotAllowed`） |
| `reason` | `string`（≤ 2000） | 否 | 撤回原因，记入 action log |

`ResubmitParams`（`resubmit` —— 重启被退回或已撤回的实例）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `instanceId` | `string` | 是 | 要重新提交的实例（`returned` / `withdrawn` 之外报 `ErrResubmitNotAllowed`） |
| `formData` | `object` | 否 | 替换的表单数据；校验同 `start` |

`AddCCParams` / `MarkCCReadParams`（`add_cc`、`mark_cc_read`）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `instanceId` | `string` | 是 | 目标实例 |
| `ccUserIds` | `string[]`（1–50） | 仅 `add_cc` | 要抄送的用户；仅当前节点开启 `isManualCcAllowed` 时允许（`ErrManualCcNotAllowed`） |

`mark_cc_read` 会把调用者在该实例上的所有未读抄送记录标记已读 ——
自助已读回执，因此不做审计。

`AddAssigneeParams`（`add_assignee` —— 动态加签）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `taskId` | `string` | 是 | 调用者自己的待办任务 |
| `userIds` | `string[]`（1–50） | 是 | 要添加的用户 |
| `addType` | `string` | 是 | `before`（新办理人先办理，原任务等待）、`after`（原办理人完成后再办理）或 `parallel`（并入当前并行组）。必须在节点的 `addAssigneeTypes` 内（`ErrAddAssigneeNotAllowed` / `ErrInvalidAddAssigneeType`） |

`RemoveAssigneeParams`（`remove_assignee`）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `taskId` | `string` | 是 | 要取消的同组任务；必须是调用者本轮访问中仍可办理的同组任务，且不能是最后一个有效办理人（`ErrLastAssigneeRemoval`），节点须允许减签（`ErrRemoveAssigneeNotAllowed`）。可减签对象由 `myTask.removableAssignees` 提供 |

`UrgeTaskParams`（`urge_task`）：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `taskId` | `string` | 是 | 要催办的待办任务 |
| `message` | `string`（≤ 500） | 否 | 随通知投递的催办消息 |

催办遵守节点按任务粒度的 `urgeCooldownMinutes`（过于频繁时报 `40601`；
非正配置默认为 30 分钟），此外该操作还有调用者每分钟 10 次的限流。

## `approval/my`

面向当前用户的自助查询。这些操作不声明 `RequiredPermission` —— 任何已认证
principal 都可调用；每个查询都在服务端锚定到调用者身份。

| Action | 入参 | 出参 |
| --- | --- | --- |
| `find_available_flows` | `FindAvailableFlowsParams` | `page.Page[AvailableFlow]` |
| `get_start_form` | `GetStartFormParams` | `StartForm` |
| `find_initiated` | `FindInitiatedParams` | `page.Page[InitiatedInstance]` |
| `find_pending_tasks` | `FindPendingTasksParams` | `page.Page[PendingTask]` |
| `find_completed_tasks` | `FindCompletedTasksParams` | `page.Page[CompletedTask]` |
| `find_cc_records` | `FindCCRecordsParams` | `page.Page[CCRecord]` |
| `get_pending_counts` | `GetPendingCountsParams` | `PendingCounts` |
| `get_instance_detail` | `GetInstanceDetailParams` | `InstanceDetail` |

请求参数（全部从 `params` 解码）：

| Action | 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- | --- |
| `find_available_flows` | `tenantId` | `string` | 否 | 租户过滤 |
| | `keyword` | `string` | 否 | 对流程名称做 contains 匹配 |
| | `labels` | `object` | 否 | label 相等过滤 —— 每一对都必须匹配 |
| | `page` / `pageSize` | `int` | 否 | 分页 |
| `get_start_form` | `tenantId` | `string` | 是 | 流程所在租户 |
| | `flowCode` | `string` | 是 | 要加载发起表单的流程 |
| `find_initiated` | `tenantId` | `string` | 否 | 租户过滤 |
| | `status` | `string` | 否 | 实例状态过滤（`running` / `approved` / `rejected` / `withdrawn` / `returned` / `terminated`） |
| | `keyword` | `string` | 否 | 对实例标题做 contains 匹配 |
| | `page` / `pageSize` | `int` | 否 | 分页 |
| `find_pending_tasks` | `tenantId` | `string` | 否 | 租户过滤 |
| | `page` / `pageSize` | `int` | 否 | 分页 |
| `find_completed_tasks` | `tenantId` | `string` | 否 | 租户过滤 |
| | `page` / `pageSize` | `int` | 否 | 分页 |
| `find_cc_records` | `tenantId` | `string` | 否 | 租户过滤 |
| | `isRead` | `bool` | 否 | 已读状态过滤 |
| | `page` / `pageSize` | `int` | 否 | 分页 |
| `get_pending_counts` | `tenantId` | `string` | 否 | 租户过滤 |
| `get_instance_detail` | `instanceId` | `string` | 是 | 要加载的实例；调用者必须是参与者 —— 申请人、（曾经的）办理人或抄送对象（`ErrAccessDenied`） |

响应 DTO（`approval/my` 包）：

`AvailableFlow` —— 一条调用者可发起的流程：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `flowId` / `flowCode` / `flowName` | `string` | 流程标识 |
| `flowIcon` | `string` \| 缺省 | 显示图标 |
| `description` | `string` \| 缺省 | 流程描述 |
| `labels` | `object` \| 缺省 | 宿主自有选择元数据 |
| `categoryId` / `categoryName` | `string` | 所属分类标识 |

`StartForm` —— 流程的提交前视图。加载受到与发起实例完全一致的
闸门约束（流程启用、发起权限、存在已发布版本），因此能渲染出的表单必然
对应一个可发起的流程：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `flowId` / `flowCode` / `flowName` | `string` | 渲染发起页头部所需的流程标识 |
| `flowIcon` | `string` \| 缺省 | 显示图标 |
| `description` | `string` \| 缺省 | 流程描述 |
| `versionId` | `string` | 表单所属的已发布版本 |
| `version` | `int` | 已发布版本号 |
| `formSchema` | JSON 文档 \| 缺省 | 宿主表单设计器文档，原样返回 |

`InitiatedInstance` —— 一条调用者提交的实例：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `instanceId` / `instanceNo` / `title` | `string` | 实例标识 |
| `flowName` | `string` | 流程显示名称 |
| `flowIcon` | `string` \| 缺省 | 流程图标 |
| `status` | `string` | 实例状态 |
| `currentNodeName` | `string` \| 缺省 | 当前进行中节点的名称 |
| `createdAt` | `DateTime` | 提交时间 |
| `finishedAt` | `DateTime` \| 缺省 | 完成时间 |

`PendingTask` —— 一条等待调用者处理的任务：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `taskId` | `string` | 提交 `process_task` 时使用的任务 id |
| `instanceId` / `instanceTitle` / `instanceNo` | `string` | 所属实例标识 |
| `flowName` / `flowIcon` | `string` | 流程显示标识 |
| `applicant` | `UserInfo` | 申请人快照 |
| `nodeName` | `string` | 任务所属节点 |
| `createdAt` | `DateTime` | 任务创建时间 |
| `deadline` | `DateTime` \| 缺省 | 节点配置了超时时的截止时间 |
| `isTimeout` | `bool` | 是否已超期 |

`CompletedTask` —— 一条调用者已处理的任务：标识字段与 `PendingTask` 相同，
另有 `status`（处理结果 —— `approved`、`rejected`、`handled`、
`transferred`、`rolled_back` 等）与 `finishedAt`；没有 `deadline` /
`isTimeout`。

`CCRecord` —— 一条发给调用者的抄送通知：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `ccRecordId` | `string` | 抄送记录 id |
| `instanceId` / `instanceTitle` / `instanceNo` | `string` | 所属实例标识 |
| `flowName` / `flowIcon` | `string` | 流程显示标识 |
| `applicant` | `UserInfo` | 申请人快照 |
| `nodeName` | `string` \| 缺省 | 产生抄送的节点；实例级抄送时缺省 |
| `isRead` | `bool` | 已读回执状态 |
| `createdAt` | `DateTime` | 投递时间 |

`PendingCounts` —— 角标计数：`pendingTaskCount`（待办任务数）与
`unreadCcCount`（未读抄送数）。

`InstanceDetail` —— 自助详情视图。每个顶层字段对应一个可渲染的关注点：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `instance` | `InstanceInfo` | 运行时状态（见下） |
| `formSchema` | JSON 文档 \| 缺省 | 版本锁定的宿主表单设计器文档，原样返回 —— 实例提交时所依据的 schema |
| `timeline` | `TimelineEntry[]` | 实例实际走过路径的逐节点记录（见下） |
| `flowGraph` | `InstanceFlowGraph` | React Flow 就绪、标注进度的只读图（见下） |
| `availableActions` | `string[]` | 面向当前查看者的动作提示（见下） |
| `fieldPermissions` | `object`（字段→权限） | 查看者维度的字段交互投影：`visible` / `editable` / `hidden` / `required`，对每个顶层表单字段都物化；客户端原样应用，且 `instance.formData` 已剥除查看者无权看到的字段（见 [节点字段权限](./flow-design.md#节点字段权限)） |
| `myTask` | `ViewerTask` \| `null` | 查看者自己的可操作上下文（见下） |

`InstanceInfo`：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `instanceId` / `instanceNo` / `title` | `string` | 实例标识 |
| `flowName` / `flowIcon` | `string` | 流程显示标识，查询时从可变流程行读取 |
| `labels` | `object` \| 缺省 | 流程的宿主自有选择元数据 —— 与 `flowName` 一样属于显示标识，不是版本锁定的快照 |
| `applicant` | `UserInfo` | 申请人快照 |
| `status` | `string` | 实例状态 |
| `currentNodeId` / `currentNodeName` | `string` \| 缺省 | 当前进行中的节点 |
| `businessRef` | `string` \| 缺省 | 不透明业务引用（业务绑定流程） |
| `formData` | `object` \| 缺省 | 表单数据，已剥除查看者无权看到的字段 |
| `createdAt` / `finishedAt` | `DateTime` | 生命周期时间戳 |

`ViewerTask` —— `process_task` 应指向的待办任务，加上客户端构建
操作 UI 所需的节点级配置，客户端无需重新推导引擎语义。查看者在该实例上
没有待办任务时为 `null`：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `taskId` | `string` | 待办任务 |
| `nodeId` | `string` | 任务所在节点 |
| `isOpinionRequired` | `bool` | 镜像节点配置：设置时 approve / reject 必须携带非空意见 |
| `addAssigneeTypes` | `string[]` | 节点允许的加签位置（`before` / `after` / `parallel`）；不允许加签时为空 |
| `rollbackTargets` | `{nodeId, name}[]` | 合法回退目标，按节点回退配置与实例访问轨迹解析，与回退命令的校验完全一致；不允许回退时为空 |
| `removableAssignees` | `{taskId, assignee, status}[]` | 查看者可减签的同组任务（`status` 为 `pending` / `waiting`），与减签命令的授权完全一致：本轮访问中仍可办理、排除查看者本人的同组任务；节点不允许减签时为空 |

`availableActions` 是查询层的 UI 提示。对申请人：实例可转入 `withdrawn`
时包含 `withdraw`，实例处于退回或已撤回状态时包含 `resubmit`。对待办任务：
办理节点为 `handle`，否则为 `approve`，然后是 `reject`，再加上当前节点允许
时的 `transfer`、`rollback`、`add_assignee`、`add_cc`。实例存在任何待办任务
时还会包含 `urge`。命令处理器仍会做自己的校验。

## `approval/admin`

管理端管理与可观测能力。对所有列表和指标查询：非 super-admin 提交的
`tenantId` 覆盖会被忽略，始终过滤到自己的租户；super-admin 可传 `tenantId`
过滤单个租户，或省略以获得跨租户视图。

| Action | 权限 | 入参 | 出参 | 审计 |
| --- | --- | --- | --- | --- |
| `find_instances` | `approval.instance.query` | `AdminFindInstancesParams` | `page.Page[Instance]` | — |
| `find_tasks` | `approval.task.query` | `AdminFindTasksParams` | `page.Page[Task]` | — |
| `get_instance_detail` | `approval.instance.detail` | `AdminGetInstanceDetailParams` | `InstanceDetail` | — |
| `find_action_logs` | `approval.action_log.query` | `AdminFindActionLogsParams` | `page.Page[ActionLog]` | — |
| `get_metrics` | `approval.metrics.query` | `AdminGetMetricsParams` | `Metrics` | — |
| `find_business_projections` | `approval.binding.query` | `AdminFindBusinessProjectionsParams` | `page.Page[BusinessProjection]` | — |
| `terminate_instance` | `approval.instance.terminate` | `AdminTerminateInstanceParams` | 成功 | 是 |
| `reassign_task` | `approval.task.reassign` | `AdminReassignTaskParams` | 成功 | 是 |
| `retry_business_projection` | `approval.binding.retry` | `AdminRetryBusinessProjectionParams` | 成功 | 是 |

请求参数（全部从 `params` 解码）：

| Action | 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- | --- |
| `find_instances` | `tenantId` | `string` | 否 | 租户过滤（仅 super-admin 生效，见上） |
| | `applicantId` | `string` | 否 | 按申请人过滤 |
| | `status` | `string` | 否 | 实例状态过滤 |
| | `flowId` | `string` | 否 | 按流程过滤 |
| | `keyword` | `string` | 否 | 对实例标题做 contains 匹配 |
| | `page` / `pageSize` | `int` | 否 | 分页 |
| `find_tasks` | `tenantId` | `string` | 否 | 租户过滤 |
| | `assigneeId` | `string` | 否 | 按办理人过滤 |
| | `instanceId` | `string` | 否 | 按所属实例过滤 |
| | `status` | `string` | 否 | 任务状态过滤（`waiting` / `pending` / `approved` / `rejected` / `handled` / `transferred` / `rolled_back` / `canceled` / `removed` / `skipped`） |
| | `page` / `pageSize` | `int` | 否 | 分页 |
| `get_instance_detail` | `instanceId` | `string` | 是 | 要加载的实例 |
| `find_action_logs` | `instanceId` | `string` | 是 | 要分页浏览审计轨迹的实例 |
| | `tenantId` | `string` | 否 | 租户过滤 |
| | `page` / `pageSize` | `int` | 否 | 分页 |
| `get_metrics` | `tenantId` | `string` | 否 | 租户范围（super-admin 可省略以跨租户） |
| `find_business_projections` | `tenantId` | `string` | 否 | 租户过滤 |
| | `status` | `string` | 否 | 投影状态过滤：`pending`、`processing`、`applied`、`failed` |
| | `page` / `pageSize` | `int` | 否 | 分页 |
| `terminate_instance` | `instanceId` | `string` | 是 | 要强制终止的运行中实例（非运行状态报 `ErrTerminateNotAllowed`） |
| | `reason` | `string`（≤ 2000） | 否 | 终止原因，记入 action log |
| `reassign_task` | `taskId` | `string` | 是 | 要改派的待办任务 |
| | `newAssigneeId` | `string` | 是 | 替换的办理人（无效时报 `ErrInvalidTransferTarget`） |
| | `reason` | `string`（≤ 2000） | 否 | 改派原因 |
| `retry_business_projection` | `projectionId` | `string` | 是 | 要立即重试的投影（不存在时报 `ErrBindingProjectionNotFound`） |

响应 DTO（`approval/admin` 包）：

`Instance` —— 管理列表中的一条实例：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `instanceId` / `instanceNo` / `title` | `string` | 实例标识 |
| `tenantId` | `string` | 所属租户 |
| `flowId` / `flowName` | `string` | 流程标识 |
| `applicant` | `UserInfo` | 申请人快照 |
| `status` | `string` | 实例状态 |
| `currentNodeName` | `string` \| 缺省 | 当前进行中的节点 |
| `createdAt` / `finishedAt` | `DateTime` | 生命周期时间戳 |

`Task` —— 管理列表中的一条任务：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `taskId` | `string` | 任务 id |
| `instanceId` / `instanceTitle` | `string` | 所属实例标识 |
| `flowName` | `string` | 流程显示名称 |
| `nodeName` | `string` | 任务所属节点 |
| `assignee` | `UserInfo` | 办理人快照 |
| `status` | `string` | 任务状态 |
| `createdAt` | `DateTime` | 创建时间 |
| `deadline` | `DateTime` \| 缺省 | 超时截止时间 |
| `finishedAt` | `DateTime` \| 缺省 | 完成时间 |

`InstanceDetail` —— `my.get_instance_detail` 的管理端对应视图，但没有查看者
维度的字段（`availableActions` / `fieldPermissions` / `myTask`）：`instance`
（`InstanceDetailInfo`）、`formSchema`（原样宿主文档）、`timeline`
（`TimelineEntry[]`）与 `flowGraph`（`InstanceFlowGraph`）。
`InstanceDetailInfo` 与 `my.InstanceInfo` 一致，另加 `tenantId`、`flowId`
与 `flowVersionId`，且其 `formData` 不做过滤。

`ActionLog` —— 一条审计记录。人员引用统一为动作发生时捕获的 `UserInfo`
快照：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `logId` | `string` | 日志条目 id |
| `action` | `string` | `ActionType` 字符串：`submit`、`approve`、`handle`、`reject`、`transfer`、`withdraw`、`cancel`、`rollback`、`add_assignee`、`remove_assignee`、`execute`、`resubmit`、`reassign`、`terminate`、`add_cc` |
| `nodeId` | `string` \| 缺省 | 动作发生的节点；实例级动作缺省 |
| `taskId` | `string` \| 缺省 | 动作针对的任务 |
| `operator` | `UserInfo` | 操作者快照 |
| `transferTo` | `UserInfo` \| 缺省 | 转办 / 改派接收人 |
| `rollbackToNodeId` | `string` \| 缺省 | 回退目标节点 |
| `addedAssignees` / `removedAssignees` | `UserInfo[]` \| 缺省 | 动态加签 / 减签变更 |
| `ccUsers` | `UserInfo[]` \| 缺省 | 手动抄送的用户 |
| `opinion` | `string` \| 缺省 | 动作意见 / 原因 |
| `attachments` | `string[]` \| 缺省 | 附件引用 |
| `createdAt` | `DateTime` | 动作时间 |

`Metrics` —— 面向仪表盘与运维告警的引擎健康聚合：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `tenantId` | `string` | 快照的租户范围；跨租户快照（仅 super-admin）时为空 |
| `capturedAt` | `DateTime` | 指标物化时刻 |
| `instanceCounts` | `object`（状态→int） | 按 `InstanceStatus` 字符串分组的实例计数 |
| `taskCounts` | `object`（状态→int） | 按 `TaskStatus` 字符串分组的任务计数 |
| `timeoutTaskCount` | `int` | 已超期的待办任务数 |
| `avgCompletionSeconds` | `float` | 全部终态实例的端到端平均时长（`createdAt` → `finishedAt`）；`-1` 表示「尚无完成的实例」 |
| `pendingBindingFailures` | `int` | 最近一次写入失败、已排期重试的投影目标数 |
| `businessProjectionCounts` | `object`（状态→int） | 按收敛状态分组的持久投影行数（`pending` / `processing` / `applied` / `failed`） |
| `pendingBusinessProjections` | `int` | 期望修订尚未应用的最终一致投影数 |

`BusinessProjection` —— 一条被绑定业务记录的运维收敛状态（回写模型见
[业务集成](./integration.md)）：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `projectionId` | `string` | 投影行 id |
| `tenantId` | `string` | 所属租户 |
| `flowId` / `flowVersionId` | `string` | 产生期望状态的流程与版本 |
| `ownerInstanceId` | `string` | 生命周期产生期望状态的实例 |
| `appliedOwnerInstanceId` | `string` \| 缺省 | 状态最近一次成功写入业务行的实例 |
| `businessTable` | `string` | 目标业务表 |
| `recordKey` | JSON 对象 | 定位绑定行的键列取值 |
| `consistency` | `string` | 配置的绑定一致性模式（`transactional` / `eventual`） |
| `desiredStatus` | `string` | 等待回写的实例状态 |
| `desiredStartedAt` / `desiredFinishedAt` | `DateTime` | 等待回写的生命周期时间戳 |
| `desiredRevision` / `appliedRevision` | `int` | 单调修订号；两者相等即已收敛 |
| `status` | `string` | 收敛状态：`pending`、`processing`、`applied`、`failed` |
| `attemptCount` | `int` | 已尝试写入次数 |
| `nextAttemptAt` | `DateTime` \| 缺省 | 下次排期重试时间 |
| `leaseUntil` | `DateTime` \| 缺省 | `processing` 期间的 worker 租约到期时间 |
| `lastError` | `string` \| 缺省 | 最近一次写入失败信息 |
| `appliedAt` | `DateTime` \| 缺省 | 期望状态最近一次应用时间 |
| `updatedAt` | `DateTime` | 最近状态变更时间 |

## 共享投影类型

详情视图（`my.get_instance_detail`、`admin.get_instance_detail`）共享公开
`approval` 包中的以下类型。

`UserInfo` —— 所有出现人员的地方使用的统一人员快照：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `id` | `string` | 用户 id |
| `name` | `string` | 动作发生时的显示名称 |
| `departmentId` / `departmentName` | `string` \| 缺省 | 动作发生时的部门快照 |

`TimelineEntry` —— 实例时间线的一步：实例实际走过路径的按时间、逐节点
记录。条件分支互斥，走过的路径永远是一条线；回退后重新进入的节点会产生
第二条记录。条目终止于当前进行中的节点 —— 不预测未到达的节点：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `kind` | `string` | 节点访问为 `start`、`approval`、`handle`、`cc`；实例级里程碑为 `withdraw`、`terminate`。结构性节点（`condition` / `end`）从不出现 |
| `nodeId` | `string` \| 缺省 | 被访问的节点；里程碑条目缺省 |
| `name` | `string` | 节点显示名称（或里程碑动作名） |
| `status` | `string` | 节点访问状态：`active`、`passed`、`rejected`、`returned`、`canceled` |
| `executionType` | `string` | 节点执行类型（`manual` / `auto_pass` / `auto_reject`） |
| `approvalMethod` | `string` | `sequential` / `parallel`（审批节点） |
| `passRule` | `string` | `all` / `any` / `ratio`（审批节点） |
| `passRatio` | `decimal` \| 缺省 | `passRule` 为 `ratio` 时的阈值 |
| `participants` | `NodeParticipant[]` | 审批 / 办理节点上每个任务一条（见下） |
| `ccRecipients` | `CCRecipient[]` | 已投递的抄送：`user`（`UserInfo`）加 `readAt` 已读回执 |
| `activities` | `Activity[]` | 节点上的旁路动作（见下）；里程碑条目仅含一条描述谁、为何关闭实例的活动 |
| `startedAt` / `finishedAt` | `DateTime` | 访问区间；进行中时 `finishedAt` 缺省 |

`NodeParticipant` —— 一次访问中一位办理人的参与情况：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `taskId` | `string` | 任务标识（任务操作的目标） |
| `user` | `UserInfo` | 办理人快照 |
| `delegator` | `UserInfo` \| 缺省 | 任务经委托到达时的委托人 |
| `status` | `string` | 任务状态原样 |
| `deadline` | `DateTime` \| 缺省 | 任务截止时间 |
| `isTimeout` | `bool` | 任务由超时扫描器裁决或升级 |
| `opinion` / `attachments` / `actionTime` | — | 从完结该任务的 action log 融合出的结果明细 |
| `transferTo` | `UserInfo` \| 缺省 | 任务被转办时的接收人 |

`Activity` —— 节点上记录的旁路动作：`action` 携带 `ActionType` 字符串
（`transfer`、`rollback`、`add_assignee`、`remove_assignee`、`add_cc`、
`reassign`、`execute`、`submit`、`resubmit`、`withdraw`、`terminate`），
外加催办记录的 `urge`。`operator` 是操作者；`opinion` 是动作自由文本
（转办理由、撤回原因、催办消息）；`target` 指向定向动作的对方（被催办的
办理人）；`transferTo`、`rollbackToNodeId` / `rollbackToNodeName`、
`addedAssignees`、`removedAssignees`、`ccUsers`、`attachments` 携带各动作的
明细，`createdAt` 是动作时间。决定本身（approve / handle / reject）不会在
活动中重复 —— 它们记录在做出决定的参与者上。

`InstanceFlowGraph` —— React Flow 就绪、只读的实例流程定义投影，标注运行
进度。`nodes` 与 `edges` 直接映射 React Flow 的节点 / 边形状，唯节点种类
保留在 `kind` 字段（React Flow 的 `type` 属于客户端）：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `nodes[].id` | `string` | React Flow 标识 —— 位置与边引用的设计时节点 key |
| `nodes[].nodeId` | `string` | 持久化流程节点 id —— 即 `actionLog.nodeId` / `rollbackToNodeId` 携带、`process_task` 回退 API 期望的 `targetNodeId`，客户端可据此映射并直接驱动回退 |
| `nodes[].kind` | `string` | 节点种类（`start` / `approval` / `handle` / `condition` / `cc` / `end`） |
| `nodes[].position` | `{x, y}` | 设计器坐标 |
| `nodes[].data` | `FlowGraphNodeData` | 节点标签、审批语义、进度 `status`（`pending` / `active` / `passed` / `rejected` / `returned` / `canceled`），以及按遍历顺序跨访问聚合的 `participants` / `ccRecipients` / `activities`，加 `startedAt` / `finishedAt` 区间 |
| `edges[]` | `{id, source, target, sourceHandle}` | 按 id 连接节点的 React Flow 边 |

## 错误面

可导入的 `approval` 包导出四个普通 Go 哨兵错误。它们可用 `errors.Is`
识别，但不是 `result.Error` 值，本身不携带 API code 或 HTTP 状态。

| 错误 | 来源包 | 含义 |
| --- | --- | --- |
| `approval.ErrCrossTenantAccess` | `approval` | 非 super-admin 调用者尝试跨租户访问 |
| `approval.ErrInvalidBusinessIdentifier` | `approval` | 业务表 / 字段标识符未通过 SQL 标识符白名单 |
| `approval.ErrUnknownNodeKind` | `approval` | `NodeDefinition.ParseData` 遇到不支持的 `kind` |
| `approval.ErrNodeDataUnmarshal` | `approval` | `NodeDefinition.ParseData` 无法解码节点 `data` |

内置审批资源通过标准 API envelope 返回模块自有的 `result.Error`。这些值
位于 internal 包中，宿主应用应把下面的 code/message 对当作公开线上契约，
而不是导入内部 Go 符号。

| Code | Code 常量 | 错误值 | i18n message key | 说明 |
| --- | --- | --- | --- | --- |
| `40001` | `ErrCodeFlowNotFound` | `ErrFlowNotFound` | `approval_flow_not_found` | 流程查找失败 |
| `40002` | `ErrCodeFlowNotActive` | `ErrFlowNotActive` | `approval_flow_not_active` | 流程已停用 |
| `40003` | `ErrCodeNoPublishedVersion` | `ErrNoPublishedVersion` | `approval_no_published_version` | 流程没有已发布版本 |
| `40004` | `ErrCodeVersionNotDraft` | `ErrVersionNotDraft` | `approval_version_not_draft` | 操作要求 draft 版本 |
| `40005` | `ErrCodeInvalidFlowDesign` | `ErrInvalidFlowDesign` | `approval_invalid_flow_design` | 图或节点设计校验失败 |
| `40006` | `ErrCodeFlowCodeExists` | `ErrFlowCodeExists` | `approval_flow_code_exists` | 流程编码重复 |
| `40007` | `ErrCodeVersionNotFound` | `ErrVersionNotFound` | `approval_version_not_found` | 流程版本查找失败 |
| `40008` | `ErrCodeInvalidBusinessIdentifier` | `ErrInvalidBusinessIdentifier` | `approval_invalid_business_identifier` | 业务表 / 字段标识符校验失败 |
| `40009` | `ErrCodeInvalidTitleTemplate` | `ErrInvalidTitleTemplate` | `approval_invalid_title_template` | 实例标题模板解析失败 |
| `40010` | `ErrCodeInvalidFormDesign` | `ErrInvalidFormDesign` | `approval_invalid_form_design` | 表单 schema 设计期校验失败 |
| `40011` | `ErrCodeBindingIncomplete` | `ErrBindingIncomplete` | `approval_binding_incomplete` | 业务绑定缺少必需的表 / 键 / 状态 / 实例 id 字段 |
| `40012` | `ErrCodeInvalidBindingMode` | `ErrInvalidBindingMode` | `approval_invalid_binding_mode` | 流程绑定模式不在枚举内 |
| `40013` | `ErrCodeInvalidInitiatorKind` | `ErrInvalidInitiatorKind` | `approval_invalid_initiator_kind` | 流程发起人类型不在枚举内 |
| `40014` | `ErrCodeInvalidStorageMode` | `ErrInvalidStorageMode` | `approval_invalid_storage_mode` | deploy 请求了 `json` / `table` 之外的存储模式 |
| `40015` | — | — | — | 未使用（原流程绑定锁）；该编码不会被复用 |
| `40016` | `ErrCodeBindingColumnsConflict` | `ErrBindingColumnsConflict` | `approval_binding_columns_conflict` | 两个业务绑定字段指向同一列 |
| `40017` | `ErrCodeBindingUnexpected` | `ErrBindingUnexpected` | `approval_binding_unexpected` | standalone 流程提交了业务绑定 |
| `40018` | `ErrCodeBindingSchemaInvalid` | `ErrBindingSchemaInvalid` | `approval_binding_schema_invalid` | 配置的绑定表或列在主库中不存在 |
| `40019` | `ErrCodeBindingKeyNotUnique` | `ErrBindingKeyNotUnique` | `approval_binding_key_not_unique` | 键列没有对应一个完整的非空主键或唯一键 |
| `40020` | `ErrCodeBindingStatusMappingInvalid` | `ErrBindingStatusMappingInvalid` | `approval_binding_status_mapping_invalid` | 状态映射包含未知状态或映射为空值 |
| `40101` | `ErrCodeInstanceNotFound` | `ErrInstanceNotFound` | `approval_instance_not_found` | 实例查找失败 |
| `40102` | `ErrCodeInstanceCompleted` | `ErrInstanceCompleted` | `approval_instance_completed` | 实例已经完结 |
| `40103` | `ErrCodeNotAllowedInitiate` | `ErrNotAllowedInitiate` | `approval_not_allowed_initiate` | 调用者不能发起该流程 |
| `40104` | `ErrCodeWithdrawNotAllowed` | `ErrWithdrawNotAllowed` | `approval_withdraw_not_allowed` | 当前状态不允许撤回 |
| `40105` | `ErrCodeResubmitNotAllowed` | `ErrResubmitNotAllowed` | `approval_resubmit_not_allowed` | 当前状态不允许重新提交 |
| `40106` | `ErrCodeInvalidInstanceTransition` | `ErrInvalidInstanceTransition` | `approval_invalid_instance_transition` | 实例状态迁移非法 |
| `40107` | `ErrCodeBusinessRefRequired` | `ErrBusinessRefRequired` | `approval_business_ref_required` | 业务绑定流程发起时缺少业务引用 |
| `40108` | `ErrCodeBindingTargetBusy` | `ErrBindingTargetBusy` | `approval_binding_target_busy` | 业务记录已被未完结的审批实例占用 |
| `40109` | `ErrCodeInvalidBusinessRef` | `ErrInvalidBusinessRef` | `approval_invalid_business_ref` | 业务引用无法解析为配置的记录键 |
| `40110` | `ErrCodeBindingProjectionNotFound` | `ErrBindingProjectionNotFound` | `approval_binding_projection_not_found` | 投影查找失败（管理端重试） |
| `40201` | `ErrCodeTaskNotFound` | `ErrTaskNotFound` | `approval_task_not_found` | 任务查找失败 |
| `40202` | `ErrCodeTaskNotPending` | `ErrTaskNotPending` | `approval_task_not_pending` | 任务不是待处理状态 |
| `40203` | `ErrCodeNotAssignee` | `ErrNotAssignee` | `approval_not_assignee` | 调用者不是任务办理人 |
| `40204` | `ErrCodeInvalidTaskTransition` | `ErrInvalidTaskTransition` | `approval_invalid_task_transition` | 任务状态迁移非法 |
| `40205` | `ErrCodeRollbackNotAllowed` | `ErrRollbackNotAllowed` | `approval_rollback_not_allowed` | 回退被禁用或此处不可用 |
| `40206` | `ErrCodeAddAssigneeNotAllowed` | `ErrAddAssigneeNotAllowed` | `approval_add_assignee_not_allowed` | 动态加签被禁用 |
| `40207` | `ErrCodeTransferNotAllowed` | `ErrTransferNotAllowed` | `approval_transfer_not_allowed` | 转办被禁用 |
| `40208` | `ErrCodeOpinionRequired` | `ErrOpinionRequired` | `approval_opinion_required` | 必填意见为空 |
| `40209` | `ErrCodeManualCcNotAllowed` | `ErrManualCcNotAllowed` | `approval_manual_cc_not_allowed` | 手动抄送被禁用 |
| `40210` | `ErrCodeRemoveAssigneeNotAllowed` | `ErrRemoveAssigneeNotAllowed` | `approval_remove_assignee_not_allowed` | 动态减签被禁用 |
| `40211` | `ErrCodeInvalidAddAssigneeType` | `ErrInvalidAddAssigneeType` | `approval_invalid_add_assignee_type` | `addType` 不是 `before`、`after`、`parallel` 之一 |
| `40212` | `ErrCodeNotApplicant` | `ErrNotApplicant` | `approval_not_applicant` | 调用者不是申请人 |
| `40213` | `ErrCodeInvalidRollbackTarget` | `ErrInvalidRollbackTarget` | `approval_invalid_rollback_target` | 回退目标不被允许 |
| `40214` | `ErrCodeLastAssigneeRemoval` | `ErrLastAssigneeRemoval` | `approval_last_assignee_removal` | 减签将导致没有有效办理人 |
| `40215` | `ErrCodeInvalidTransferTarget` | `ErrInvalidTransferTarget` | `approval_invalid_transfer_target` | 转办或改派目标非法 |
| `40216` | `ErrCodeNoUsersSpecified` | `ErrNoUsersSpecified` | `approval_no_users_specified` | 用户列表操作没有收到目标用户 |
| `40301` | `ErrCodeNoAssignee` | `ErrNoAssignee` | `approval_no_assignee` | 无法解析出办理人 |
| `40302` | `ErrCodeAssigneeResolveFailed` | `ErrAssigneeResolveFailed` | `approval_assignee_resolve_failed` | 办理人解析器执行失败 |
| `40401` | `ErrCodeFormValidationFailed` | `ErrFormValidationFailed` | `approval_form_validation_failed` | 通用表单校验失败 |
| `40401` | `ErrCodeFormValidationFailed` | `ErrFormDataTooLarge` | `approval_form_data_too_large` | 同一编码；JSON 编码后的 `formData` 超过 64 KiB |
| `40401` | `ErrCodeFormValidationFailed` | 动态表单校验 `result.Err` | `approval_form_field_not_defined`、`approval_form_field_required`、`approval_form_field_must_be_string`、`approval_form_field_must_be_number`、`approval_form_field_must_be_integer`、`approval_form_field_min_length`、`approval_form_field_max_length`、`approval_form_field_invalid_validation`、`approval_form_field_pattern_mismatch`、`approval_form_field_min_value`、`approval_form_field_max_value`、`approval_form_field_empty`、`approval_form_field_invalid_file_item`、`approval_form_field_must_be_file`、`approval_form_field_invalid_value`、`approval_form_field_must_be_row_list`、`approval_form_field_must_be_row_object`、`approval_form_field_min_rows`、`approval_form_field_max_rows`、`approval_form_field_table_cell` | 字段级校验消息为动态构造 |
| `40601` | `ErrCodeUrgeCooldown` | 动态催办 `result.Err` | `approval_urge_too_frequent` | 无静态哨兵；消息用 `minutes` 渲染；非正的 `urgeCooldownMinutes` 默认 30 分钟 |
| `40701` | `ErrCodeAccessDenied` | `ErrAccessDenied` | `approval_access_denied` | 调用者缺少审批域访问权 |
| `40702` | `ErrCodeTerminateNotAllowed` | `ErrTerminateNotAllowed` | `approval_terminate_not_allowed` | 当前实例状态不允许终止 |

`ErrEventRouteNotTransactional`、`ErrEventRouteNotSubscribable`、
`ErrTenantNotResolved` 等启动与租户解析诊断错误位于
`internal/approval/...` 下；它们不是可导入的公开 Go API，但事件路由或租户
principal 配置有误时，运维人员可能在包装后的报错信息里看到它们。

---

下一步：[流程设计](./flow-design.md) 了解 `deploy` 背后的设计器传输格式，或 [实例运行时](./runtime.md) 了解实例动作背后的生命周期语义。
