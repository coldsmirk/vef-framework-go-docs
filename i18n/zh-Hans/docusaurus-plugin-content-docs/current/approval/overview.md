---
sidebar_position: 1
slug: /approval
---

# 审批模块

`approval` 模块提供完整的工作流引擎，用于构建基于审批的业务流程。支持可视化流程设计（兼容 React Flow）、多级审批链、条件分支、并行审批、委托、回退和事务性事件发布。

本分类把模块拆成五页：本页概览（启用、装配、配置）、[RPC 资源](./resources.md)、[流程设计](./flow-design.md)、[实例运行时](./runtime.md) 和 [事件与集成](./integration.md)。

## 启用模块

审批是一个可选功能模块。它有意不包含在默认 `vef.Run(...)` boot graph 中，
所以不需要审批工作流的应用不会注册它的 API resources、CQRS handlers、
engine、业务投影 worker 或 timeout scanners。

需要时显式启用：

```go
vef.Run(
    vef.ApprovalModule,
    app.Module,
)
```

## 事件路由前置条件

审批会通过 `event.WithTx` 发布 `approval.*` 事件，因此除
`approval.instance.binding_failed` 之外的每种事件类型都必须解析到
**transactional** transport。模块会在启动时通过 `event.RouteInspector`
断言这一点——路由配错时应用直接启动失败，而不是静默降级：

```toml
[vef.event]
default_transport = "memory"

[[vef.event.routing]]
pattern    = "approval.*"
transports = ["outbox", "redis_stream"]
```

从 v0.38 起，业务投影不再消费生命周期事件（它从持久化的
`apv_business_projection` 表收敛——见
[业务状态投影](./integration.md#业务状态投影)），所以模块自身不再要求路由带
可订阅 sink：只写 `["outbox"]` 的路由就能通过启动检查。当**宿主**通过
`approval.SubscribeInstance` 或 `approval.BindCommand` 订阅审批事件时，
仍要像上例一样在 `outbox` 之外列出配置的 outbox `sink` transport
（单节点 `memory`，跨节点 `redis_stream`）——outbox transport 是
publish-only 的，订阅方挂在 sink 上。transport、路由语义和 outbox relay 见
[事件总线](../infrastructure/event-bus.md)。

`InstanceBindingFailedEvent` 是 transactional-route 启动检查的例外：它由
最终一致投影 worker 在审批事务已经提交后发出。

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
| 业务投影 | `apv_business_projection` | 业务绑定流程的持久化期望状态收敛（v0.38） |

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

[vef.approval.business_binding]
consistency   = "synchronous"  # 或 "eventual"
scan_interval = "10s"          # eventual worker 的扫描节奏
batch_size    = 100            # 每次扫描处理的投影数
```

`auto_migrate` 是普通 boolean 开关，不会由 `ApprovalConfig.ApplyDefaults()`
自动设为 true；需要启动时执行 approval DDL 时必须显式开启。
`cc_record_retention` 只清理已经读过的 CC 记录。

`business_binding`（v0.38）控制审批状态如何投影到绑定的业务表：
`synchronous`（默认）在审批事务内写业务行，`eventual` 先提交期望状态、由
后台 worker 收敛业务行（见
[业务状态投影](./integration.md#业务状态投影)）。`consistency` 超出枚举或
worker 配置为负值会在启动时的配置校验直接失败
（`config.ErrInvalidApprovalBindingConsistency` /
`ErrInvalidApprovalBusinessBindingWorkerConfig`）。

> 老版本里归属 `[vef.approval]` 的 `outbox_relay_interval` / `outbox_max_retries` / `outbox_batch_size` 已在 v0.21 迁移至 `[vef.event.transports.outbox]`，由全框架统一的 outbox transport 服务所有模块——参考 [事件总线](../infrastructure/event-bus.md)。

详见[配置参考](../reference/configuration-reference.md)。

## 绑定模式

| 模式 | 常量 | Wire value | 说明 |
| --- | --- | --- | --- |
| 独立 | `BindingStandalone` | `standalone` | 表单数据存储在审批模块自有表中 |
| 业务 | `BindingBusiness` | `business` | 关联到已有的业务数据表 |

业务绑定通过单个 `Flow.BusinessBinding` 文档（`approval.BusinessBindingConfig`）
将审批流程与业务表关联：`tableName`、复合 `keyColumns`、`statusColumn`、
必填的 `instanceIdColumn` CAS 防护栏、可选的 `startedAtColumn` /
`finishedAtColumn`，以及可选的 `statusMapping`（见
[业务状态投影](./integration.md#业务状态投影)）。绑定会在每次部署时快照到
流程版本上，其状态通过持久化的 `apv_business_projection` 表收敛。

---

下一步：[RPC 资源](./resources.md) 了解 API 面，或 [流程设计](./flow-design.md) 了解节点类型与设计器 wire shape。
