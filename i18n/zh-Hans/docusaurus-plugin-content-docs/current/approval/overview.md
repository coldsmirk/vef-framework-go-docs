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
engine、binding listener 或 timeout scanners。

需要时显式启用：

```go
vef.Run(
    vef.ApprovalModule,
    app.Module,
)
```

## 事件路由前置条件

审批会通过 `event.WithTx` 发布 `approval.*` 事件，它的 binding listener
也会订阅这些事件。宿主应用必须把 `approval.*` 路由到一个带有可订阅 sink
的 transactional transport，例如 sink 为 Redis Streams 的 outbox 路由：

```toml
[vef.event]
default_transport = "memory"

[[vef.event.routing]]
pattern    = "approval.*"
transports = ["outbox", "redis_stream"]
```

路由必须在 `outbox` 之外同时列出配置的 outbox `sink`（这里是
`redis_stream`；进程内的 `memory` transport 同样合格）：outbox transport 是
publish-only 的，订阅方——包括模块自己的 binding listener——都挂在 sink
transport 上。审批模块的 binding listener 和 outbox 发布两侧都会在启动时通过
`event.RouteInspector` 断言路由，路由配错时应用直接启动失败，而不是静默降级。
transport、路由语义和 outbox relay 见 [事件总线](../infrastructure/event-bus.md)。

`InstanceBindingFailedEvent` 是 transactional-route 启动检查的例外：它由异步
binding listener 在审批事务已经提交后发出。`InstanceCompletedEvent` 的路由要求
最严格，因为 binding listener 会订阅它；路由必须在 transactional outbox 之外
同时包含可订阅 sink transport，例如 `memory` 或 `redis_stream`。

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

`auto_migrate` 是普通 boolean 开关，不会由 `ApprovalConfig.ApplyDefaults()`
自动设为 true；需要启动时执行 approval DDL 时必须显式开启。
`cc_record_retention` 只清理已经读过的 CC 记录。

> 老版本里归属 `[vef.approval]` 的 `outbox_relay_interval` / `outbox_max_retries` / `outbox_batch_size` 已在 v0.21 迁移至 `[vef.event.transports.outbox]`，由全框架统一的 outbox transport 服务所有模块——参考 [事件总线](../infrastructure/event-bus.md)。

详见[配置参考](../reference/configuration-reference.md)。

## 绑定模式

| 模式 | 常量 | Wire value | 说明 |
| --- | --- | --- | --- |
| 独立 | `BindingStandalone` | `standalone` | 表单数据存储在审批模块自有表中 |
| 业务 | `BindingBusiness` | `business` | 关联到已有的业务数据表 |

业务绑定通过 `BusinessTable`、`BusinessPKField` 和 `BusinessStatusField` 将审批流程与业务表关联，另外还有可选的
`BusinessInstanceIDField`、`BusinessStartedAtField` 和 `BusinessFinishedAtField` 联动列（见
[业务写回联动矩阵](./integration.md#业务写回联动矩阵)）。

---

下一步：[RPC 资源](./resources.md) 了解 API 面，或 [流程设计](./flow-design.md) 了解节点类型与设计器 wire shape。
